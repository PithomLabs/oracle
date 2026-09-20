package packetv1

import (
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
)

func newTestRegistry(t *testing.T) *domainpack.PackRegistry {
	t.Helper()
	registry := domainpack.NewRegistry()
	pack := &bmistv1.Pack{
		PackID:          "bmist",
		Version:         "1.0.0",
		ClaimTypes:      []string{"derived", "accommodated", "postulated"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		DebtVocabulary:  []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		Falsifiers:              []string{"counterexample", "contradiction"},
		HumanGatedTransitions:   []string{"faithfulness_review", "scope_clarification", "obstruction_assessment"},
		ConsequentialActions:    []bmistv1.ConsequentialAction{{Action: "publish_claim", Requires: "promoted", Gates: []string{"faithfulness_review"}}},
	}
	if err := registry.Register(pack); err != nil {
		t.Fatalf("failed to register test pack: %v", err)
	}
	return registry
}

func validPacket() *Packet {
	return &Packet{
		SchemaVersion: SchemaVersion,
		Role:          RoleWork,
		PacketID:      "test-packet-001",
		PackRef:       "bmist-1.0.0",
		Agent:         Agent{ID: "test-agent"},
		Beliefs: []Belief{
			{LocalID: "b1", Claim: "Fisher-rigidity holds for 2D", ClaimType: "derived", Debt: []string{"needMap"}},
		},
		Evidence: []Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
		Edges: []Edge{
			{LocalID: "edge1", FromRef: "local:b1", ToRef: "local:b1", Kind: EdgeDerives},
		},
		Tasks: []Task{
			{LocalID: "t1", Title: "Run toy model check"},
		},
	}
}

func TestValidateValidPacket(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	if err := Validate(pkt, registry); err != nil {
		t.Errorf("valid packet rejected: %v", err)
	}
}

func TestValidateNilPacket(t *testing.T) {
	registry := newTestRegistry(t)
	if err := Validate(nil, registry); err == nil {
		t.Error("nil packet accepted")
	}
}

func TestValidateBadSchemaVersion(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.SchemaVersion = "wrong-version"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad schema version accepted")
	}
}

func TestValidateMissingPacketID(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.PacketID = ""
	if err := Validate(pkt, registry); err == nil {
		t.Error("missing packet_id accepted")
	}
}

func TestValidateBadRole(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Role = "invalid"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad role accepted")
	}
}

func TestValidateBadPackRef(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.PackRef = "nonexistent-1.0.0"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad pack_ref accepted")
	}
}

func TestValidateEmptyBeliefClaim(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Beliefs[0].Claim = ""
	if err := Validate(pkt, registry); err == nil {
		t.Error("empty claim accepted")
	}
}

func TestValidateDuplicateBeliefLocalID(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Beliefs = append(pkt.Beliefs, Belief{LocalID: "b1", Claim: "dup", ClaimType: "derived"})
	if err := Validate(pkt, registry); err == nil {
		t.Error("duplicate local_id accepted")
	}
}

func TestValidateBadEvidenceHash(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Evidence[0].ContentSHA256 = "not-a-hash"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad hash accepted")
	}
}

func TestValidateBadEdgeKind(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Edges[0].Kind = "invalid"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad edge kind accepted")
	}
}

func TestValidateEmptyTaskTitle(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Tasks[0].Title = ""
	if err := Validate(pkt, registry); err == nil {
		t.Error("empty task title accepted")
	}
}

func TestValidateFreeTextEdgeRef(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Edges[0].FromRef = "free text claim"
	if err := Validate(pkt, registry); err == nil {
		t.Error("free-text edge ref accepted")
	}
}

func TestValidateBadBeliefRef(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Evidence[0].BeliefRef = "bad-ref"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad belief_ref accepted")
	}
}

func TestValidateBadCanonicalUUID(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Evidence[0].BeliefRef = "canonical:belief:not-a-uuid"
	if err := Validate(pkt, registry); err == nil {
		t.Error("bad canonical UUID accepted")
	}
}

func TestValidateEmptyLocalRef(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Evidence[0].BeliefRef = "local:"
	if err := Validate(pkt, registry); err == nil {
		t.Error("empty local ref accepted")
	}
}

func TestValidateAdversarialRole(t *testing.T) {
	registry := newTestRegistry(t)
	pkt := validPacket()
	pkt.Role = RoleAdversarial
	if err := Validate(pkt, registry); err != nil {
		t.Errorf("adversarial role rejected: %v", err)
	}
}
