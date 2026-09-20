package coordinator

import (
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
)

func newTestRetirementCoordinator(t *testing.T) *Coordinator {
	t.Helper()
	solvent := &MockSolventClient{
		ListEvidenceForBeliefFn: func(beliefID, scenarioID string) ([]map[string]interface{}, error) {
			// Return evidence matching the class each test expects
			switch beliefID {
			case "belief-001": // needMap + reproducible_artifact
				return []map[string]interface{}{
					{"id": "ev-001", "belief_id": beliefID, "provenance_class": "reproducible_artifact"},
				}, nil
			case "belief-003": // needNullModel + operator_asserted
				return []map[string]interface{}{
					{"id": "ev-003", "belief_id": beliefID, "provenance_class": "operator_asserted"},
				}, nil
			case "belief-005": // needInvariant + reproducible_artifact
				return []map[string]interface{}{
					{"id": "ev-005", "belief_id": beliefID, "provenance_class": "reproducible_artifact"},
				}, nil
			case "belief-007": // needObstruction + reproducible_artifact
				return []map[string]interface{}{
					{"id": "ev-007", "belief_id": beliefID, "provenance_class": "reproducible_artifact"},
				}, nil
			case "belief-008": // needFaithfulnessReview + operator_asserted
				return []map[string]interface{}{
					{"id": "ev-008", "belief_id": beliefID, "provenance_class": "operator_asserted"},
				}, nil
			default:
				// belief-002, -004, -006 (wrong class tests), -009 (unknown debt)
				return []map[string]interface{}{}, nil
			}
		},
	}
	conductor := &MockConductorClient{}

	// Load the bmist pack
	packData := []byte(`{
		"pack_id": "bmist",
		"version": "1.0.0",
		"name": "BM-IST Fisher-Rigidity Domain Pack",
		"description": "Domain Pack for the BM-IST Fisher-rigidity claim verification POC",
		"claim_types": ["derived", "accommodated", "postulated"],
		"evidence_classes": ["reproducible_artifact", "operator_asserted"],
		"debt_vocabulary": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
		"initial_debt": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
		"retirement_rules": {
			"needMap": {"evidence_class": "reproducible_artifact", "rule": "map_check"},
			"needInvariant": {"evidence_class": "reproducible_artifact", "rule": "invariant_check"},
			"needToyCheck": {"evidence_class": "reproducible_artifact", "rule": "toy_model_check"},
			"needNullModel": {"evidence_class": "operator_asserted", "rule": "scope_clarification"},
			"needObstruction": {"evidence_class": "reproducible_artifact", "rule": "obstruction_construction"},
			"needFaithfulnessReview": {"evidence_class": "operator_asserted", "rule": "faithfulness_review"}
		},
		"falsifiers": ["counterexample", "contradiction", "missing_evidence", "alternative_explanation"],
		"consequential_actions": [],
		"human_gated_transitions": ["faithfulness_review", "scope_clarification", "obstruction_assessment"]
	}`)

	pack, err := bmistv1.ParsePack(packData)
	if err != nil {
		t.Fatalf("failed to parse pack: %v", err)
	}

	registry := domainpack.NewRegistry()
	if err := registry.Register(pack); err != nil {
		t.Fatalf("failed to register pack: %v", err)
	}

	return NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)
}

func TestRetireDebt_NeedMapReproducibleArtifact(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-001",
		ScenarioID:    "track-g0",
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
}

func TestRetireDebt_NeedMapOperatorAsserted(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-002",
		ScenarioID:    "track-g0",
		DebtItem:      "needMap",
		EvidenceClass: "operator_asserted",
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

func TestRetireDebt_NeedNullModelOperatorAsserted(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-003",
		ScenarioID:    "track-g0",
		DebtItem:      "needNullModel",
		EvidenceClass: "operator_asserted",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRetireDebt_NeedNullModelReproducibleArtifact(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-004",
		ScenarioID:    "track-g0",
		DebtItem:      "needNullModel",
		EvidenceClass: "reproducible_artifact",
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

func TestRetireDebt_NeedInvariantReproducibleArtifact(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-005",
		ScenarioID:    "track-g0",
		DebtItem:      "needInvariant",
		EvidenceClass: "reproducible_artifact",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRetireDebt_NeedToyCheckOperatorAsserted(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-006",
		ScenarioID:    "track-g0",
		DebtItem:      "needToyCheck",
		EvidenceClass: "operator_asserted",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "refused" {
		t.Errorf("expected refused, got %s", record.Result)
	}
}

func TestRetireDebt_NeedObstructionReproducibleArtifact(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-007",
		ScenarioID:    "track-g0",
		DebtItem:      "needObstruction",
		EvidenceClass: "reproducible_artifact",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRetireDebt_NeedFaithfulnessReviewOperatorAsserted(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-008",
		ScenarioID:    "track-g0",
		DebtItem:      "needFaithfulnessReview",
		EvidenceClass: "operator_asserted",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	if record.Result != "executed" {
		t.Errorf("expected executed, got %s", record.Result)
	}
}

func TestRetireDebt_UnknownDebtItem(t *testing.T) {
	c := newTestRetirementCoordinator(t)

	req := DecisionRequest{
		Type:          DecisionRetireDebt,
		BeliefID:      "belief-009",
		ScenarioID:    "track-g0",
		DebtItem:      "unknownDebt",
		EvidenceClass: "reproducible_artifact",
	}

	record, err := c.SubmitDecision(req)
	if err != nil {
		t.Fatalf("SubmitDecision failed: %v", err)
	}

	// Unknown debt items MUST be refused (fail-closed)
	if record.Result != "refused" {
		t.Errorf("expected refused for unknown debt, got %s", record.Result)
	}
	if record.RefusalReason == "" {
		t.Error("expected refusal reason for unknown debt")
	}
}
