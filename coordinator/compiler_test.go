package coordinator

import (
	"testing"
)

func newTestCoordinator(t *testing.T) *Coordinator {
	t.Helper()
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}

	c := NewWithMockClients(solvent, conductor, "test-operator")
	return c
}

func TestNewCoordinatorValid(t *testing.T) {
	c := newTestCoordinator(t)
	if c == nil {
		t.Fatal("coordinator is nil")
	}
}

func TestNewCoordinatorMissingOperator(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}

	c := NewWithMockClients(solvent, conductor, "")
	if c == nil {
		t.Fatal("coordinator is nil")
	}
	// Note: operatorID validation happens in submitDecision, not in constructor
}

func TestNewCoordinatorWithMockClients(t *testing.T) {
	solvent := &MockSolventClient{}
	conductor := &MockConductorClient{}

	c := NewWithMockClients(solvent, conductor, "test-operator")
	if c.operatorID != "test-operator" {
		t.Errorf("expected test-operator, got %s", c.operatorID)
	}
}
