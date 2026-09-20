package verifier

import (
	"context"
	"testing"
	"time"

	"github.com/PithomLabs/oracle/verifier/model"
)

func testArtifact(ref string) model.VerificationArtifact {
	return model.VerificationArtifact{
		RunID:           "test-run-001",
		VerifierVersion: "0.1.0",
		VerifierHash:    "abc123",
		ClaimHash:       "claim-hash-001",
		InputHash:       "input-hash-001",
		ArtifactHash:    "artifact-hash-001",
		Timestamp:       time.Now(),
		Steps: []model.Step{
			{Index: 0, Description: "test", Verified: true, CheckType: model.CheckTypeMachineVerified},
		},
		Result:      model.ResultConfirmed,
		EvidenceRef: ref,
	}
}

func TestRegisterTrustedValid(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-001")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Errorf("valid registration failed: %v", err)
	}
}

func TestRegisterTrustedEmptyRef(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("")
	if err := registry.registerTrusted(artifact); err == nil {
		t.Error("empty evidence_ref accepted")
	}
}

func TestRegisterTrustedEmptyHash(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-002")
	artifact.ArtifactHash = ""
	if err := registry.registerTrusted(artifact); err == nil {
		t.Error("empty artifact hash accepted")
	}
}

func TestRegisterTrustedDuplicate(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-003")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Fatalf("first registration failed: %v", err)
	}
	if err := registry.registerTrusted(artifact); err == nil {
		t.Error("duplicate registration accepted")
	}
}

func TestResolveValid(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-004")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	resolved, err := registry.Resolve(context.Background(), "evidence-004")
	if err != nil {
		t.Errorf("resolve failed: %v", err)
	}
	if resolved.EvidenceRef != "evidence-004" {
		t.Errorf("wrong evidence_ref: got %q", resolved.EvidenceRef)
	}
}

func TestResolveUnknown(t *testing.T) {
	registry := NewArtifactRegistry()
	_, err := registry.Resolve(context.Background(), "nonexistent")
	if err == nil {
		t.Error("unknown reference resolved")
	}
}

func TestVerifyHashCorrect(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-005")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if err := registry.VerifyHash("evidence-005", "artifact-hash-001"); err != nil {
		t.Errorf("correct hash rejected: %v", err)
	}
}

func TestVerifyHashWrong(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-006")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	if err := registry.VerifyHash("evidence-006", "wrong-hash"); err == nil {
		t.Error("wrong hash accepted")
	}
}

func TestVerifyHashUnknown(t *testing.T) {
	registry := NewArtifactRegistry()
	if err := registry.VerifyHash("nonexistent", "hash"); err == nil {
		t.Error("unknown reference accepted")
	}
}

func TestAsReader(t *testing.T) {
	registry := NewArtifactRegistry()
	artifact := testArtifact("evidence-007")
	if err := registry.registerTrusted(artifact); err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	reader := registry.AsReader()
	resolved, err := reader.Resolve(context.Background(), "evidence-007")
	if err != nil {
		t.Errorf("reader resolve failed: %v", err)
	}
	if resolved.EvidenceRef != "evidence-007" {
		t.Errorf("wrong evidence_ref: got %q", resolved.EvidenceRef)
	}
}

func TestConcurrentReads(t *testing.T) {
	registry := NewArtifactRegistry()
	for i := 0; i < 10; i++ {
		artifact := testArtifact("evidence-concurrent-" + string(rune('a'+i)))
		if err := registry.registerTrusted(artifact); err != nil {
			t.Fatalf("registration failed: %v", err)
		}
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			ref := "evidence-concurrent-" + string(rune('a'+idx))
			_, _ = registry.Resolve(context.Background(), ref)
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestConcurrentRegistrations(t *testing.T) {
	registry := NewArtifactRegistry()
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			artifact := testArtifact("evidence-reg-" + string(rune('a'+idx)))
			_ = registry.registerTrusted(artifact)
			done <- true
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
