package coordinator

import (
	"fmt"
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
)

// ============================================================================
// F7 VERIFICATION: Retraction cancels linked Conductor tasks
// ============================================================================

func TestRetractCancelsLinkedTasks(t *testing.T) {
	var cancelledIDs []string
	conductor := &testConductorClient{
		tasks: []map[string]interface{}{
			{"id": "task-A", "status": "active", "governance_ref": `{"provider":"solvent","reference_id":"scenario-001"}`},
			{"id": "task-B", "status": "proposed", "governance_ref": `{"provider":"solvent","reference_id":"scenario-001"}`},
			{"id": "task-C", "status": "active", "governance_ref": `{"provider":"solvent","reference_id":"scenario-002"}`},
			{"id": "task-D", "status": "accepted", "governance_ref": `{"provider":"solvent","reference_id":"scenario-001"}`},
		},
		onCancel: func(taskID string) error {
			cancelledIDs = append(cancelledIDs, taskID)
			return nil
		},
	}
	solvent := &MockSolventClient{
		RetractBeliefFn: func(beliefID, scenarioID string) error {
			return nil
		},
	}
	registry := newTestPackRegistry(t)
	c := NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)

	req := DecisionRequest{
		Type:       DecisionRetract,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}
	if record.Result != "executed" {
		t.Fatalf("expected executed, got %s", record.Result)
	}

	foundA, foundB, foundC, foundD := false, false, false, false
	for _, id := range cancelledIDs {
		switch id {
		case "task-A":
			foundA = true
		case "task-B":
			foundB = true
		case "task-C":
			foundC = true
		case "task-D":
			foundD = true
		}
	}

	if !foundA {
		t.Errorf("task-A (active, same scenario) was NOT cancelled; cancelled=%v", cancelledIDs)
	}
	if !foundB {
		t.Errorf("task-B (proposed, same scenario) was NOT cancelled; cancelled=%v", cancelledIDs)
	}
	if foundC {
		t.Errorf("task-C (different scenario) was incorrectly cancelled; cancelled=%v", cancelledIDs)
	}
	if foundD {
		t.Errorf("task-D (accepted, same scenario) was incorrectly cancelled; cancelled=%v", cancelledIDs)
	}
}

func TestRetractCancelFailureNotSilentlySwallowed(t *testing.T) {
	conductor := &testConductorClient{
		tasks: []map[string]interface{}{
			{"id": "task-A", "status": "active", "governance_ref": `{"provider":"solvent","reference_id":"scenario-001"}`},
		},
		onCancel: func(taskID string) error {
			return fmt.Errorf("cancel failed: database error")
		},
	}
	solvent := &MockSolventClient{
		RetractBeliefFn: func(beliefID, scenarioID string) error {
			return nil
		},
	}
	registry := newTestPackRegistry(t)
	c := NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)

	req := DecisionRequest{
		Type:       DecisionRetract,
		BeliefID:   "belief-001",
		ScenarioID: "scenario-001",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}
	// Cancel failure must NOT be silently swallowed: retraction should be refused
	if record.Result != "refused" {
		t.Fatalf("expected refused when cancel fails, got %s", record.Result)
	}
	if record.RefusalReason == "" {
		t.Error("expected refusal reason when cancel fails")
	}
}

// ============================================================================
// F8 VERIFICATION: Pack rules endpoint returns actual rules
// ============================================================================

func TestGetPackRules_ReturnsActualPackRules(t *testing.T) {
	registry := newTestPackRegistry(t)
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)

	rules := c.GetPackRules()
	if rules == nil {
		t.Fatal("expected non-nil rules map")
	}

	rulesMap, ok := rules["rules"].(map[string]domainpack.RetirementRule)
	if !ok {
		t.Fatalf("expected rules.rules to be a map[string]domainpack.RetirementRule, got %T", rules["rules"])
	}

	needMapRule, ok := rulesMap["needMap"]
	if !ok {
		t.Fatalf("expected needMap rule to exist")
	}
	if needMapRule.EvidenceClass != "reproducible_artifact" {
		t.Errorf("expected needMap evidence_class=reproducible_artifact, got %v", needMapRule.EvidenceClass)
	}
	if needMapRule.Rule != "map_check" {
		t.Errorf("expected needMap rule=map_check, got %v", needMapRule.Rule)
	}

	needNullModelRule := rulesMap["needNullModel"]
	if needNullModelRule.EvidenceClass != "operator_asserted" {
		t.Errorf("expected needNullModel evidence_class=operator_asserted, got %v", needNullModelRule.EvidenceClass)
	}
	if needNullModelRule.Rule != "scope_clarification" {
		t.Errorf("expected needNullModel rule=scope_clarification, got %v", needNullModelRule.Rule)
	}

	needFaithfulnessRule := rulesMap["needFaithfulnessReview"]
	if needFaithfulnessRule.EvidenceClass != "operator_asserted" {
		t.Errorf("expected needFaithfulnessReview evidence_class=operator_asserted, got %v", needFaithfulnessRule.EvidenceClass)
	}
	if needFaithfulnessRule.Rule != "faithfulness_review" {
		t.Errorf("expected needFaithfulnessReview rule=faithfulness_review, got %v", needFaithfulnessRule.Rule)
	}
}

// ============================================================================
// RCP DETERMINISM: mergeActivity sorts by created_at
// ============================================================================

func TestMergeActivity_SortsByCreatedAt(t *testing.T) {
	conductor := []map[string]interface{}{
		{"id": "c2", "created_at": "2026-01-02T00:00:00Z", "action": "late"},
		{"id": "c1", "created_at": "2026-01-01T00:00:00Z", "action": "early"},
	}
	solvent := []map[string]interface{}{
		{"id": "s2", "created_at": "2026-01-04T00:00:00Z", "action": "late"},
		{"id": "s1", "created_at": "2026-01-03T00:00:00Z", "action": "early"},
	}

	merged := mergeActivity(conductor, solvent)

	if merged[0]["id"] != "c1" {
		t.Errorf("expected c1 first, got %s", merged[0]["id"])
	}
	if merged[1]["id"] != "c2" {
		t.Errorf("expected c2 second, got %s", merged[1]["id"])
	}
	if merged[2]["id"] != "s1" {
		t.Errorf("expected s1 third, got %s", merged[2]["id"])
	}
	if merged[3]["id"] != "s2" {
		t.Errorf("expected s2 fourth, got %s", merged[3]["id"])
	}
}

// ============================================================================
// PACK RESOLUTION: fail-closed when pack unavailable
// ============================================================================

func TestRetireDebt_PackRegistryUnavailable_Refused(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClients(solvent, conductor, "test-operator")

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
	if record.Result != "refused" {
		t.Errorf("expected refused when pack registry unavailable, got %s", record.Result)
	}
}

func TestRetireDebt_UnknownScenario_Refused(t *testing.T) {
	registry := newTestPackRegistry(t)
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}
	c := NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-001",
		ScenarioID:    "unknown-scenario",
		DebtItem:      "needMap",
		EvidenceClass: "reproducible_artifact",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}
	if record.Result != "refused" {
		t.Errorf("expected refused for unknown scenario, got %s", record.Result)
	}
}

// ============================================================================
// HELPER: test conductor client with task list and cancel
// ============================================================================

type testConductorClient struct {
	tasks    []map[string]interface{}
	onCancel func(taskID string) error
}

func (c *testConductorClient) CreateTask(projectID, title, description, governanceRef string) (string, error) {
	return "task-mock", nil
}
func (c *testConductorClient) GetTask(taskID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": taskID, "status": "active"}, nil
}
func (c *testConductorClient) ListTasks(projectID string) ([]map[string]interface{}, error) {
	return c.tasks, nil
}
func (c *testConductorClient) ListDependencies(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (c *testConductorClient) ListActivity(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (c *testConductorClient) CancelTask(taskID string) error {
	if c.onCancel != nil {
		return c.onCancel(taskID)
	}
	return nil
}
