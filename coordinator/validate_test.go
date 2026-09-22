package coordinator

import (
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

func newTestPackRegistry(t *testing.T) *domainpack.PackRegistry {
	t.Helper()
	registry := domainpack.NewRegistry()
	pack := &bmistv1.Pack{
		PackID:          "bmist",
		Version:         "1.0.0",
		ClaimTypes:      []string{"derived", "accommodated", "postulated"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		DebtVocabulary:  []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		InitialDebt:     []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		RetirementRules: map[string]domainpack.RetirementRule{
			"needMap":                {EvidenceClass: "reproducible_artifact", Rule: "map_check"},
			"needInvariant":          {EvidenceClass: "reproducible_artifact", Rule: "invariant_check"},
			"needToyCheck":           {EvidenceClass: "reproducible_artifact", Rule: "toy_model_check"},
			"needNullModel":          {EvidenceClass: "operator_asserted", Rule: "scope_clarification"},
			"needObstruction":        {EvidenceClass: "reproducible_artifact", Rule: "obstruction_construction"},
			"needFaithfulnessReview": {EvidenceClass: "operator_asserted", Rule: "faithfulness_review"},
		},
		Falsifiers:            []string{"counterexample", "contradiction"},
		HumanGatedTransitions: []string{"faithfulness_review", "scope_clarification", "obstruction_assessment"},
		ConsequentialActions:  []bmistv1.ConsequentialAction{{Action: "publish_claim", Requires: "promoted", Gates: []string{"faithfulness_review"}}},
	}
	if err := registry.Register(pack); err != nil {
		t.Fatalf("failed to register test pack: %v", err)
	}
	return registry
}

func TestValidatePacketValid(t *testing.T) {
	registry := newTestPackRegistry(t)
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-001",
		PackRef:       "bmist@1.0.0",
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "test claim", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
	}

	if err := ValidatePacket(pkt, registry); err != nil {
		t.Errorf("ValidatePacket failed: %v", err)
	}
}

func TestCompileDebt(t *testing.T) {
	packInitial := []string{"needMap", "needInvariant"}
	beliefDebt := []string{"needToyCheck"}

	compiled := CompileDebt(beliefDebt, packInitial)

	if len(compiled) != 3 {
		t.Errorf("expected 3 debt items, got %d", len(compiled))
	}
}

func TestCompileDebtDuplicates(t *testing.T) {
	packInitial := []string{"needMap", "needInvariant"}
	beliefDebt := []string{"needMap", "needToyCheck"}

	compiled := CompileDebt(beliefDebt, packInitial)

	if len(compiled) != 3 {
		t.Errorf("expected 3 debt items, got %d", len(compiled))
	}
}
