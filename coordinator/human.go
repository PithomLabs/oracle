package coordinator

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
)

type DecisionType string

const (
	DecisionPromote    DecisionType = "PROMOTE"
	DecisionRetireDebt DecisionType = "RETIRE_DEBT"
	DecisionRetract    DecisionType = "RETRACT"
	DecisionReopen     DecisionType = "REOPEN"
	DecisionAuthorize  DecisionType = "AUTHORIZE"
	DecisionRefuse     DecisionType = "REFUSE"
)

type DecisionRequest struct {
	Type          DecisionType `json:"type"`
	BeliefID      string       `json:"belief_id"`
	ScenarioID    string       `json:"scenario_id"`
	DebtItem      string       `json:"debt_item,omitempty"`
	EvidenceClass string       `json:"evidence_class,omitempty"`
	Reason        string       `json:"reason,omitempty"`
}

type DecisionRecord struct {
	ID             string       `json:"id"`
	Type           DecisionType `json:"type"`
	BeliefID       string       `json:"belief_id"`
	ScenarioID     string       `json:"scenario_id"`
	ActorID        string       `json:"actor_id"`
	Result         string       `json:"result"`
	RefusalReason  string       `json:"refusal_reason,omitempty"`
	SolventAuditID string       `json:"solvent_audit_id,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

type PacketStatus struct {
	PacketID  string   `json:"packet_id"`
	Status    string   `json:"status"`
	BeliefIDs []string `json:"belief_ids"`
}

type DecisionContext struct {
	BeliefID    string   `json:"belief_id"`
	Claim       string   `json:"claim"`
	Status      string   `json:"status"`
	Debt        []string `json:"debt"`
	EvidenceIDs []string `json:"evidence_ids"`
}

func (c *Coordinator) submitDecision(req DecisionRequest) (*DecisionRecord, error) {
	if c.operatorID == "" {
		return nil, fmt.Errorf("operator principal ID is required")
	}
	record := &DecisionRecord{
		ID:         fmt.Sprintf("decision-%s-%d", req.BeliefID, time.Now().UnixNano()),
		Type:       req.Type,
		BeliefID:   req.BeliefID,
		ScenarioID: req.ScenarioID,
		ActorID:    c.operatorID,
		CreatedAt:  time.Now(),
	}
	switch req.Type {
	case DecisionPromote:
		return c.handlePromote(req, record)
	case DecisionRetireDebt:
		return c.handleRetireDebt(req, record)
	case DecisionRetract:
		return c.handleRetract(req, record)
	case DecisionReopen:
		return c.handleReopen(req, record)
	case DecisionAuthorize:
		return c.handleAuthorize(req, record)
	case DecisionRefuse:
		return c.handleRefuse(req, record)
	default:
		return nil, fmt.Errorf("unsupported decision type: %s", req.Type)
	}
}

func (c *Coordinator) handlePromote(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	err := c.solventClient.PromoteBelief(req.BeliefID, req.ScenarioID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	record.Result = "executed"
	return record, nil
}

func (c *Coordinator) handleRetireDebt(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	if c.packRegistry == nil {
		record.Result = "refused"
		record.RefusalReason = "pack registry unavailable: cannot validate retirement"
		return record, nil
	}
	pack, err := c.resolvePack(req.ScenarioID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = fmt.Sprintf("pack resolution failed: %v", err)
		return record, nil
	}
	if err := ValidateDebtItemMembership(req.DebtItem, pack); err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	if err := ValidateRetirementRule(req.DebtItem, req.EvidenceClass, pack); err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	evidenceIDs, verifiedClass, err := c.verifyEvidenceClass(req.BeliefID, req.ScenarioID, req.EvidenceClass)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = fmt.Sprintf("evidence verification failed: %v", err)
		return record, nil
	}
	instrumentRef := buildInstrumentRef(evidenceIDs)
	_ = verifiedClass
	err = c.solventClient.Discharge(req.ScenarioID, req.BeliefID, req.DebtItem, instrumentRef, c.operatorID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	record.Result = "executed"
	record.SolventAuditID = instrumentRef
	return record, nil
}

func (c *Coordinator) verifyEvidenceClass(beliefID, scenarioID, claimedClass string) ([]string, string, error) {
	evidence, err := c.solventClient.ListEvidenceForBelief(beliefID, scenarioID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch evidence: %w", err)
	}
	if len(evidence) == 0 {
		return nil, "", fmt.Errorf("no evidence found for belief %s", beliefID)
	}
	var matchedIDs []string
	var actualClass string
	for _, e := range evidence {
		provenanceClass, _ := e["provenance_class"].(string)
		evidenceID, _ := e["id"].(string)
		if provenanceClass == claimedClass && evidenceID != "" {
			matchedIDs = append(matchedIDs, evidenceID)
			actualClass = provenanceClass
		}
	}
	if len(matchedIDs) == 0 {
		var foundClasses []string
		for _, e := range evidence {
			if pc, ok := e["provenance_class"].(string); ok {
				foundClasses = append(foundClasses, pc)
			}
		}
		return nil, "", fmt.Errorf("evidence class mismatch: claimed %q, found %v", claimedClass, foundClasses)
	}
	return matchedIDs, actualClass, nil
}

func buildInstrumentRef(evidenceIDs []string) string {
	if len(evidenceIDs) == 0 {
		return ""
	}
	ref := "evidence:"
	for i, id := range evidenceIDs {
		if i > 0 {
			ref += ","
		}
		ref += id
	}
	return ref
}

func (c *Coordinator) resolvePack(scenarioID string) (domainpack.Pack, error) {
	packID, packVersion := scenarioPackMapping(scenarioID)
	pack, err := c.packRegistry.Get(packID, packVersion)
	if err != nil {
		return nil, fmt.Errorf("pack %s@%s not registered: %w", packID, packVersion, err)
	}
	return pack, nil
}

func scenarioPackMapping(scenarioID string) (string, string) {
	return "bmist", "1.0.0"
}

func (c *Coordinator) handleRetract(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	err := c.solventClient.RetractBelief(req.BeliefID, req.ScenarioID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	cancelledTasks, cancelErr := c.cancelLinkedTasks(req.ScenarioID)
	if cancelErr != nil {
		record.Result = "refused"
		record.RefusalReason = fmt.Sprintf("failed to cancel linked tasks: %v", cancelErr)
		return record, nil
	}
	record.Result = "executed"
	if len(cancelledTasks) > 0 {
		record.SolventAuditID = fmt.Sprintf("cancelled:%d", len(cancelledTasks))
	}
	return record, nil
}

func (c *Coordinator) cancelLinkedTasks(scenarioID string) ([]string, error) {
	projects := []string{"default"}
	var cancelled []string
	for _, projectID := range projects {
		tasks, err := c.conductorClient.ListTasks(projectID)
		if err != nil {
			return nil, err
		}
		for _, task := range tasks {
			taskID, _ := task["id"].(string)
			governanceRef, _ := task["governance_ref"].(string)
			status, _ := task["status"].(string)
			if taskID == "" || governanceRef == "" {
				continue
			}
			if status == "cancelled" || status == "completed" || status == "accepted" {
				continue
			}
			if !strings.Contains(governanceRef, scenarioID) {
				continue
			}
			var govRef struct {
				ReferenceID string `json:"reference_id"`
			}
			if err := json.Unmarshal([]byte(governanceRef), &govRef); err != nil {
				continue
			}
			if govRef.ReferenceID != scenarioID {
				continue
			}
			if err := c.conductorClient.CancelTask(taskID); err != nil {
				return nil, err
			}
			cancelled = append(cancelled, taskID)
		}
	}
	return cancelled, nil
}

func (c *Coordinator) handleReopen(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	newBeliefID, err := c.solventClient.CreateBelief(req.BeliefID, req.Reason, req.ScenarioID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}

	// Create derives edge from old belief to new belief
	if err := c.solventClient.CreateEdge(req.BeliefID, newBeliefID, "derives"); err != nil {
		record.Result = "refused"
		record.RefusalReason = fmt.Sprintf("new belief created but edge failed: %v", err)
		return record, nil
	}

	record.Result = "executed"
	record.SolventAuditID = newBeliefID
	return record, nil
}

func (c *Coordinator) handleAuthorize(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	err := c.solventClient.ApproveTarget(req.BeliefID)
	if err != nil {
		record.Result = "refused"
		record.RefusalReason = err.Error()
		return record, nil
	}
	record.Result = "executed"
	return record, nil
}

func (c *Coordinator) handleRefuse(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
	record.Result = "refused"
	record.RefusalReason = req.Reason
	return record, nil
}

func (c *Coordinator) packetStatus(packetID string) (*PacketStatus, error) {
	return &PacketStatus{
		PacketID: packetID,
		Status:   "compiled",
	}, nil
}

func (c *Coordinator) decisionContext(beliefID string) (*DecisionContext, error) {
	return &DecisionContext{
		BeliefID: beliefID,
		Status:   "active",
	}, nil
}
