package coordinator

import (
	"fmt"
	"testing"
)

func TestSubmitDecisionPromote(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestSubmitDecisionRefuse(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionRefuse,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
		Reason:     "test refusal",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "refused" {
		t.Errorf("expected refused, got %s", record.Result)
	}
}

func TestSubmitDecisionMissingOperator(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	_, err := c.SubmitDecision(req)
	if err == nil {
		t.Fatal("expected error for missing operator")
	}
}

func TestPacketStatus(t *testing.T) {
	c := newTestCoordinator(t)
	status, err := c.PacketStatus("test-001")
	if err != nil {
		t.Fatalf("PacketStatus failed: %v", err)
	}

	if status.PacketID != "test-001" {
		t.Errorf("wrong packet_id: %s", status.PacketID)
	}
}

func TestDecisionContext(t *testing.T) {
	c := newTestCoordinator(t)
	ctx, err := c.DecisionContext("belief-001")
	if err != nil {
		t.Fatalf("DecisionContext failed: %v", err)
	}

	if ctx.BeliefID != "belief-001" {
		t.Errorf("wrong belief_id: %s", ctx.BeliefID)
	}
}

// Phase 7 Tests

func TestPromoteWithDebt(t *testing.T) {
	solvent := &MockSolventClient{
		PromoteBeliefFn: func(beliefID, scenarioID string) error {
			// Simulate Solvent returning a Verdict (debt blocks promotion)
			return fmt.Errorf("promotion blocked: belief has open debt")
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-with-debt",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "refused" {
		t.Errorf("expected refused, got %s", record.Result)
	}

	if record.RefusalReason == "" {
		t.Error("expected refusal reason")
	}
}

func TestPromoteWithNoDebt(t *testing.T) {
	solvent := &MockSolventClient{
		PromoteBeliefFn: func(beliefID, scenarioID string) error {
			// Simulate successful promotion
			return nil
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-no-debt",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRetractPromotedBeliefWithIntent(t *testing.T) {
	solvent := &MockSolventClient{
		RetractBeliefFn: func(beliefID, scenarioID string) error {
			// Simulate Solvent cascading: retracts belief and cancels live intent
			return nil
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionRetract,
		BeliefID:   "belief-promoted-with-intent",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRefuseDecision(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionRefuse,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
		Reason:     "operator refused to promote",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "refused" {
		t.Errorf("expected refused, got %s", record.Result)
	}

	if record.RefusalReason != "operator refused to promote" {
		t.Errorf("wrong refusal reason: %s", record.RefusalReason)
	}
}

func TestDecisionWithoutPrincipal(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	_, err := c.SubmitDecision(req)
	if err == nil {
		t.Fatal("expected error for missing principal")
	}

	expected := "operator principal ID is required"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestDecisionDeduplication(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	// Submit same decision twice
	record1, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("first SubmitDecision failed: %v", err)
	}

	record2, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("second SubmitDecision failed: %v", err)
	}

	// Both should succeed (Coordinator doesn't persist, so no dedup needed)
	if record1.Result != "executed" {
		t.Errorf("first decision: expected executed, got %s", record1.Result)
	}

	if record2.Result != "executed" {
		t.Errorf("second decision: expected executed, got %s", record2.Result)
	}
}

func TestDecisionHistoryReconstructable(t *testing.T) {
	solvent := &MockSolventClient{
		PromoteBeliefFn: func(beliefID, scenarioID string) error {
			// Simulate Solvent creating audit_activity entry
			return nil
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionPromote,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	// Submit decision
	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	// Verify record contains necessary information for reconstruction
	if record.ActorID != "test-operator" {
		t.Errorf("wrong actor_id: %s", record.ActorID)
	}

	if record.BeliefID != "belief-001" {
		t.Errorf("wrong belief_id: %s", record.BeliefID)
	}

	if record.ScenarioID != "scenario-001" {
		t.Errorf("wrong scenario_id: %s", record.ScenarioID)
	}
}

// Phase 7 Additional Tests

func TestRetireDebt(t *testing.T) {
	solvent := &MockSolventClient{
		DischargeFn: func(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error {
			return nil
		},
		ListEvidenceForBeliefFn: func(beliefID, scenarioID string) ([]map[string]interface{}, error) {
			return []map[string]interface{}{
				{"id": "ev-001", "belief_id": beliefID, "provenance_class": "reproducible_artifact"},
			}, nil
		},
	}
	conductor := &MockConductorClient{}
	registry := newTestPackRegistry(t)
	c := NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-001",
		ScenarioID:    "scenario-001",
		DebtItem:      "needMap",
		EvidenceClass: "reproducible_artifact",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
	if record.SolventAuditID == "" {
		t.Error("expected SolventAuditID (InstrumentRef) to be set")
	}
}

func TestReopenCreatesDerivesLineage(t *testing.T) {
	solvent := &MockSolventClient{
		CreateBeliefFn: func(parentID, claim, scenarioID string) (string, error) {
			// Simulate creating new belief with derives edge
			return "belief-reopened-001", nil
		},
		CreateEdgeFn: func(parentID, childID, kind string) error {
			// Simulate creating derives edge
			if kind != "derives" {
				t.Errorf("expected derives edge, got %s", kind)
			}
			return nil
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionReopen,
		BeliefID:   "belief-retracted-001",
		ScenarioID: "scenario-001",
		Reason:     "reopening with new evidence",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestAuthorize(t *testing.T) {
	solvent := &MockSolventClient{
		ApproveTargetFn: func(targetID string) error {
			// Simulate successful authorization
			return nil
		},
	}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

	req := DecisionRequest{
		Type:       DecisionAuthorize,
		BeliefID:   "target-001",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}
