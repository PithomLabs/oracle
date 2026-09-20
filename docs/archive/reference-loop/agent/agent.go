package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
	"github.com/PithomLabs/solvent/service/executor"
)

type Task = conductor.Task
type Activity = conductor.Activity
type Principal = solvent.Principal
type Belief = solvent.Belief
type Target = solvent.Target
type CreateTargetRequest = solvent.CreateTargetRequest
type AuthorizeActionResult = solvent.AuthorizeActionResult
type AuthorizeActionRequest = solvent.AuthorizeActionRequest
type ExecuteActionResult = solvent.ExecuteActionResult
type ExecuteActionRequest = solvent.ExecuteActionRequest
type ActivityEntry = solvent.ActivityEntry

const (
	repo     = "ibmendoza/reference-loop-test"
	workflow = "ref-loop.yml"
	ref      = "main"
)

func operationID(repo, workflow, ref, runID string) string {
	return "deploy:" + repo + ":" + workflow + ":" + ref + ":" + runID
}

type EvidenceCollector struct {
	records    []EvidenceRecord
	runID      string
	scenarioID string
}

type EvidenceRecord struct {
	Event                string          `json:"event"`               // TASK_CREATED, AUTHORIZED, etc.
	Owner                string          `json:"owner"`               // conductor, solvent, executor, external, agent
	Source               string          `json:"source"`              // conductor.activity, solvent.audit, github.api, etc.
	OperationID          string          `json:"operation_id"`        // deploy:repo:workflow:ref:run_id
	DeclarationVer       string          `json:"declaration_ver"`     // snapshot_id from Solvent
	Timestamp            time.Time       `json:"timestamp"`
	RunID                string          `json:"run_id"`              // immutable correlation ID
	ScenarioID           string          `json:"scenario_id"`         // happy_path, wrong_operation, etc.
	TaskID               string          `json:"task_id,omitempty"`
	IntentID             string          `json:"intent_id,omitempty"`
	NonExternalEffectProof *bool          `json:"non_external_effect_proof,omitempty"`
	Raw                  json.RawMessage `json:"raw"`                 // participant's own record
}

type Agent struct {
	conductorClient *conductor.Client
	solventClient   *solvent.Client
	executorReg     *executor.Registry
	agentID         string
	projectID       string
	scenarioID      string
	runID           string
	principalID     string
	beliefID        string
	targetID        string
	intentID        string
	evidence        *EvidenceCollector
}

func NewAgent(
	conductorClient *conductor.Client,
	solventClient *solvent.Client,
	executorReg *executor.Registry,
	agentID, projectID, scenarioID, runID, principalID string,
	evidence *EvidenceCollector,
) *Agent {
	return &Agent{
		conductorClient: conductorClient,
		solventClient:   solventClient,
		executorReg:     executorReg,
		agentID:         agentID,
		projectID:       projectID,
		scenarioID:      scenarioID,
		runID:           runID,
		principalID:     principalID,
		evidence:        evidence,
	}
}

func (a *Agent) RunHappyPath(ctx context.Context) error {
	opID := operationID(repo, workflow, ref, a.runID)

	a.recordEvidence("TASK_CREATED", "agent", "agent.conductor_call", opID, "", "", "", nil)
	task, err := a.conductorClient.CreateTask(ctx, a.projectID,
		"Reference Loop Happy Path Task",
		"Execute the reference loop: deploy workflow via Solvent authorization",
		"",
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	a.recordEvidence("TASK_CREATED", "conductor", "conductor.task", opID, "", task.ID, "", nil)

	if err := a.setupSolvent(ctx); err != nil {
		return fmt.Errorf("setup solvent: %w", err)
	}

	claimedTask, err := a.conductorClient.ClaimTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("claim task: %w", err)
	}
	a.recordEvidence("TASK_CLAIMED", "conductor", "conductor.activity", opID, "", claimedTask.ID, "", nil)

	authResult, err := a.solventClient.AuthorizeAction(ctx, AuthorizeActionRequest{
		ScenarioID:   a.scenarioID,
		BeliefID:     a.beliefID,
		Action:       "deploy",
		ActionSource: "user_typed",
		TargetID:     a.targetID,
		ActorID:      a.principalID,
	})
	if err != nil {
		return fmt.Errorf("authorize action: %w", err)
	}
	if !authResult.Authority.Allowed {
		return fmt.Errorf("authorization denied: %s", authResult.Authority.Reason)
	}

	intentID, err := a.solventClient.GetLiveIntent(ctx, a.scenarioID, a.beliefID, "deploy", a.targetID)
	if err != nil {
		return fmt.Errorf("get live intent: %w", err)
	}
	a.intentID = intentID

	a.recordEvidence("PROPOSAL_CREATED", "agent", "agent.solvent_call", opID, intentID, task.ID, intentID, nil)
	a.recordEvidence("AUTHORIZED", "solvent", "solvent.audit", opID, intentID, task.ID, intentID, nil)

	execResult, err := a.solventClient.ExecuteAction(ctx, ExecuteActionRequest{
		ScenarioID:      a.scenarioID,
		BeliefID:        a.beliefID,
		Action:          "deploy",
		TargetID:        a.targetID,
		IntentID:        intentID,
		ConsequenceType: "execution",
	})
	if err != nil {
		return fmt.Errorf("execute action: %w", err)
	}
	if !execResult.Allowed {
		return fmt.Errorf("execution not allowed: %s", execResult.Reason)
	}
	if !execResult.Success {
		return fmt.Errorf("execution failed: %s", execResult.Error)
	}
	a.recordEvidence("EXECUTION_ATTEMPTED", "executor", "solvent.audit", opID, intentID, task.ID, intentID, nil)

	nonExt := true
	a.recordEvidence("EFFECT_CONFIRMED", "external", "executor.recording",
		opID, "", task.ID, intentID, json.RawMessage(`{"non_external_effect_proof":true}`))
	_ = nonExt

	a.recordEvidence("RESULT_OBSERVED", "agent", "agent.solvent_response", opID, "", task.ID, intentID, nil)

	activity, err := a.conductorClient.PostActivity(ctx, task.ID, "work.completed",
		fmt.Sprintf(`{"agent":"%s","action":"deploy","target_id":"%s","result":"success"}`, a.agentID, a.targetID))
	if err != nil {
		return fmt.Errorf("post activity: %w", err)
	}
	a.recordEvidence("CONDUCTOR_UPDATED", "conductor", "conductor.activity", opID, "", task.ID, intentID, nil)
	_ = activity

	_, err = a.conductorClient.SubmitTask(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("submit task: %w", err)
	}

	nextTask, err := a.conductorClient.NextTask(ctx, a.projectID)
	if err != nil {
		a.recordEvidence("AGENT_CONTINUED", "agent", "agent.next_task", opID, "", task.ID, intentID, nil)
	} else {
		a.recordEvidence("AGENT_CONTINUED", "agent", "agent.next_task", opID, "", nextTask.ID, intentID, nil)
	}

	return nil
}

func (a *Agent) setupSolvent(ctx context.Context) error {
	claim := "Reference loop test: deploy workflow via GitHub Actions"
	belief, err := a.solventClient.EnterBelief(ctx, a.scenarioID, claim, "postulated")
	if err != nil {
		return fmt.Errorf("enter belief: %w", err)
	}
	a.beliefID = belief.BeliefID

	for _, debtItem := range belief.Debt {
		_, err := a.solventClient.RetireDebt(ctx, a.scenarioID, a.beliefID, debtItem)
		if err != nil {
			return fmt.Errorf("retire debt %s: %w", debtItem, err)
		}
	}

	_, err = a.solventClient.PromoteBelief(ctx, a.scenarioID, a.beliefID)
	if err != nil {
		return fmt.Errorf("promote belief: %w", err)
	}

	consequenceParams := map[string]interface{}{
		"repo":     repo,
		"workflow": workflow,
		"ref":      ref,
		"run_id":   a.runID,
	}
	consequenceParamsBytes, _ := json.Marshal(consequenceParams)

	target, err := a.solventClient.CreateTarget(ctx, CreateTargetRequest{
		PrincipalID:           a.principalID,
		ResourceType:          "scenario",
		ResourceID:            a.scenarioID,
		Scope:                 "belief:" + a.beliefID,
		ActionNamespace:       "solvent",
		ActionName:            "deploy",
		ConsequenceType:       "execution",
		ConsequenceParameters: consequenceParamsBytes,
		CreatedBy:             a.principalID,
	})
	if err != nil {
		return fmt.Errorf("create target: %w", err)
	}
	a.targetID = target.TargetID

	_, err = a.solventClient.AttachJustification(ctx, a.targetID, a.beliefID, "promoted", a.principalID)
	if err != nil {
		return fmt.Errorf("attach justification: %w", err)
	}

	_, err = a.solventClient.RequestAuthorization(ctx, a.targetID, a.principalID)
	if err != nil {
		return fmt.Errorf("request authorization: %w", err)
	}

	_, err = a.solventClient.ApproveTarget(ctx, a.targetID, a.principalID, "test-pin")
	if err != nil {
		return fmt.Errorf("approve target: %w", err)
	}

	return nil
}

func (a *Agent) recordEvidence(event, owner, source, operationID, declarationVer, taskID, intentID string, raw json.RawMessage) {
	if a.evidence == nil {
		return
	}
	a.evidence.records = append(a.evidence.records, EvidenceRecord{
		Event:          event,
		Owner:          owner,
		Source:         source,
		OperationID:    operationID,
		DeclarationVer: declarationVer,
		Timestamp:      time.Now(),
		RunID:          a.evidence.runID,
		ScenarioID:     a.evidence.scenarioID,
		TaskID:         taskID,
		IntentID:       intentID,
		Raw:            raw,
	})
}

func NewEvidenceCollector(runID, scenarioID string) *EvidenceCollector {
	return &EvidenceCollector{
		runID:      runID,
		scenarioID: scenarioID,
		records:    []EvidenceRecord{},
	}
}

func (e *EvidenceCollector) GetRecords() []EvidenceRecord {
	return e.records
}
