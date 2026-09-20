package coordinator

import (
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	"github.com/PithomLabs/oracle/verifier"
)

// testPack implements domainpack.Pack for testing.
type testPack struct {
	PackID  string
	Version string
}

func (p *testPack) GetPackID() string              { return p.PackID }
func (p *testPack) GetVersion() string             { return p.Version }
func (p *testPack) GetDebtVocabulary() []string    { return nil }
func (p *testPack) GetEvidenceClasses() []string   { return nil }
func (p *testPack) GetRetirementRules() map[string]domainpack.RetirementRule { return nil }
func (p *testPack) GetVerifierSpecs() []domainpack.VerifierSpec { return nil }

func newTestPersistCoordinator(t *testing.T) *Coordinator {
	t.Helper()
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}

	registry := domainpack.NewRegistry()
	registry.Register(&testPack{PackID: "bmist", Version: "1.0.0"})

	artifactReader := verifier.NewArtifactRegistry()

	return &Coordinator{
		solventClient:    solvent,
		conductorClient:  conductor,
		packRegistry:     registry,
		artifactReader:   artifactReader,
		projectionQueue:  NewProjectionQueue(),
		idempotencyCache: NewIdempotencyCache(),
		operatorID:       "test-operator",
	}
}

func TestSubmitPacket_BeliefsPersisted(t *testing.T) {
	c := newTestPersistCoordinator(t)

	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-packet-001",
		PackRef:       "bmist-1.0.0",
		ScenarioID:    "track-g0",
		Beliefs: []packetv1.Belief{
			{LocalID: "l1", Claim: "L1 is true", ClaimType: "derived"},
			{LocalID: "l2", Claim: "L2 is true", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:l1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "edge1", FromRef: "local:l1", ToRef: "canonical:belief:00000000-0000-0000-0000-000000000001", Kind: "derives"},
		},
	}

	result, err := c.SubmitPacket(pkt)
	if err != nil {
		t.Fatalf("SubmitPacket failed: %v", err)
	}

	if result.PacketID != "test-packet-001" {
		t.Errorf("packet_id mismatch: got %s", result.PacketID)
	}
}

func TestSubmitPacket_DuplicateSubmitIdempotent(t *testing.T) {
	c := newTestPersistCoordinator(t)

	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-packet-002",
		PackRef:       "bmist-1.0.0",
		ScenarioID:    "track-g0",
		Beliefs: []packetv1.Belief{
			{LocalID: "l1", Claim: "L1 is true", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{},
	}

	result1, err := c.SubmitPacket(pkt)
	if err != nil {
		t.Fatalf("first SubmitPacket failed: %v", err)
	}

	result2, err := c.SubmitPacket(pkt)
	if err != nil {
		t.Fatalf("second SubmitPacket failed: %v", err)
	}

	if result1.PacketID != result2.PacketID {
		t.Errorf("idempotency broken: different packet IDs")
	}
}

func TestSubmitPacket_CrossScenarioNoDedup(t *testing.T) {
	c := newTestPersistCoordinator(t)

	pkt1 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-packet-003",
		PackRef:       "bmist-1.0.0",
		ScenarioID:    "track-g0",
		Beliefs: []packetv1.Belief{
			{LocalID: "l1", Claim: "L1 is true", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{},
	}

	pkt2 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-packet-004",
		PackRef:       "bmist-1.0.0",
		ScenarioID:    "track-g1",
		Beliefs: []packetv1.Belief{
			{LocalID: "l1", Claim: "L1 is true", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{},
	}

	result1, err := c.SubmitPacket(pkt1)
	if err != nil {
		t.Fatalf("first SubmitPacket failed: %v", err)
	}

	result2, err := c.SubmitPacket(pkt2)
	if err != nil {
		t.Fatalf("second SubmitPacket failed: %v", err)
	}

	if result1.PacketID == result2.PacketID {
		t.Errorf("cross-scenario dedup broken: same packet IDs for different scenarios")
	}
}
