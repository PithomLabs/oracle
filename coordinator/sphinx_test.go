package coordinator

import (
	"testing"
)

func TestAuthorizationContext(t *testing.T) {
	c := newTestCoordinator(t)
	projection, err := c.AuthorizationContext("target-001")
	if err != nil {
		t.Fatalf("AuthorizationContext failed: %v", err)
	}

	if projection.Target.ID != "target-001" {
		t.Errorf("wrong target ID: %s", projection.Target.ID)
	}

	if projection.CurrentResult != "HUMAN_REVIEW" {
		t.Errorf("expected HUMAN_REVIEW, got %s", projection.CurrentResult)
	}
}

func TestSphinxProjectionReadOnly(t *testing.T) {
	c := newTestCoordinator(t)

	// Call authorization context twice
	p1, _ := c.AuthorizationContext("target-001")
	p2, _ := c.AuthorizationContext("target-001")

	// Should return same projection (read-only)
	if p1.CurrentResult != p2.CurrentResult {
		t.Error("projection should be deterministic")
	}
}
