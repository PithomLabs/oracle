package coordinator

import "fmt"

// SolventClientInterface is an interface for the Solvent client.
type SolventClientInterface interface {
	CreateBelief(parentID, claim, scenarioID string) (string, error)
	CreateEdge(parentID, childID, kind string) error
	CreateEvidence(beliefID, provenanceClass, contentSHA256 string) (string, error)
	PromoteBelief(beliefID, scenarioID string) error
	RetireDebt(beliefID, debtItem, scenarioID string) error
	RetractBelief(beliefID, scenarioID string) error
	ApproveTarget(targetID string) error
	GetBelief(beliefID string) (map[string]interface{}, error)
	ListBeliefs(scenarioID string) ([]map[string]interface{}, error)
	ListEvidence(scenarioID string) ([]map[string]interface{}, error)
	ListEdges(scenarioID string) ([]map[string]interface{}, error)
	ListIntents(scenarioID string) ([]map[string]interface{}, error)
	ListActivities(scenarioID string) ([]map[string]interface{}, error)
	Discharge(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error
	ListEvidenceForBelief(beliefID, scenarioID string) ([]map[string]interface{}, error)
}

// ConductorClientInterface is an interface for the Conductor client.
type ConductorClientInterface interface {
	CreateTask(projectID, title, description, governanceRef string) (string, error)
	GetTask(taskID string) (map[string]interface{}, error)
	ListTasks(projectID string) ([]map[string]interface{}, error)
	ListDependencies(taskID string) ([]map[string]interface{}, error)
	ListActivity(taskID string) ([]map[string]interface{}, error)
	CancelTask(taskID string) error
}

// MockSolventClient is a mock implementation of the Solvent client.
type MockSolventClient struct {
	CreateBeliefFn       func(parentID, claim, scenarioID string) (string, error)
	CreateEdgeFn         func(parentID, childID, kind string) error
	CreateEvidenceFn     func(beliefID, provenanceClass, contentSHA256 string) (string, error)
	PromoteBeliefFn      func(beliefID, scenarioID string) error
	RetireDebtFn         func(beliefID, debtItem, scenarioID string) error
	RetractBeliefFn      func(beliefID, scenarioID string) error
	ApproveTargetFn      func(targetID string) error
	GetBeliefFn          func(beliefID string) (map[string]interface{}, error)
	ListEvidenceForBeliefFn func(beliefID, scenarioID string) ([]map[string]interface{}, error)
	DischargeFn             func(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error

	// evidenceTracker records created evidence for ListEvidenceForBelief
	evidenceTracker []map[string]interface{}
	evidenceCounter int
}

// CreateBelief creates a new belief in Solvent.
func (m *MockSolventClient) CreateBelief(parentID, claim, scenarioID string) (string, error) {
	if m.CreateBeliefFn != nil {
		return m.CreateBeliefFn(parentID, claim, scenarioID)
	}
	return "belief-mock-001", nil
}

// CreateEdge creates an edge between beliefs in Solvent.
func (m *MockSolventClient) CreateEdge(parentID, childID, kind string) error {
	if m.CreateEdgeFn != nil {
		return m.CreateEdgeFn(parentID, childID, kind)
	}
	return nil
}

// CreateEvidence admits evidence to Solvent.
func (m *MockSolventClient) CreateEvidence(beliefID, provenanceClass, contentSHA256 string) (string, error) {
	if m.CreateEvidenceFn != nil {
		return m.CreateEvidenceFn(beliefID, provenanceClass, contentSHA256)
	}
	m.evidenceCounter++
	evidenceID := fmt.Sprintf("evidence-mock-%03d", m.evidenceCounter)
	m.evidenceTracker = append(m.evidenceTracker, map[string]interface{}{
		"id":               evidenceID,
		"belief_id":        beliefID,
		"provenance_class": provenanceClass,
		"content_sha256":   contentSHA256,
	})
	return evidenceID, nil
}

// PromoteBelief promotes a belief in Solvent.
func (m *MockSolventClient) PromoteBelief(beliefID, scenarioID string) error {
	if m.PromoteBeliefFn != nil {
		return m.PromoteBeliefFn(beliefID, scenarioID)
	}
	return nil
}

// RetireDebt retires debt on a belief in Solvent.
func (m *MockSolventClient) RetireDebt(beliefID, debtItem, scenarioID string) error {
	if m.RetireDebtFn != nil {
		return m.RetireDebtFn(beliefID, debtItem, scenarioID)
	}
	return nil
}

// RetractBelief retracts a belief in Solvent.
func (m *MockSolventClient) RetractBelief(beliefID, scenarioID string) error {
	if m.RetractBeliefFn != nil {
		return m.RetractBeliefFn(beliefID, scenarioID)
	}
	return nil
}

// ApproveTarget approves an authority target in Solvent.
func (m *MockSolventClient) ApproveTarget(targetID string) error {
	if m.ApproveTargetFn != nil {
		return m.ApproveTargetFn(targetID)
	}
	return nil
}

// GetBelief retrieves a belief from Solvent.
func (m *MockSolventClient) GetBelief(beliefID string) (map[string]interface{}, error) {
	if m.GetBeliefFn != nil {
		return m.GetBeliefFn(beliefID)
	}
	return map[string]interface{}{
		"id":      beliefID,
		"status":  "active",
		"debt":    []string{"needMap"},
	}, nil
}

// ListBeliefs returns all beliefs in a scenario.
func (m *MockSolventClient) ListBeliefs(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListEvidenceForBelief returns evidence for a specific belief.
func (m *MockSolventClient) ListEvidenceForBelief(beliefID, scenarioID string) ([]map[string]interface{}, error) {
	if m.ListEvidenceForBeliefFn != nil {
		return m.ListEvidenceForBeliefFn(beliefID, scenarioID)
	}
	// Return tracked evidence matching this belief
	var result []map[string]interface{}
	for _, e := range m.evidenceTracker {
		if e["belief_id"] == beliefID {
			result = append(result, e)
		}
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result, nil
}

// ListEvidence returns all evidence in a scenario.
func (m *MockSolventClient) ListEvidence(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListEdges returns all edges in a scenario.
func (m *MockSolventClient) ListEdges(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListIntents returns all action intents in a scenario.
func (m *MockSolventClient) ListIntents(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListActivities returns audit activity entries for a scenario.
func (m *MockSolventClient) ListActivities(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// Discharge records attributed debt discharge in Solvent.
func (m *MockSolventClient) Discharge(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error {
	if m.DischargeFn != nil {
		return m.DischargeFn(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy)
	}
	return nil
}

// MockConductorClient is a mock implementation of the Conductor client.
type MockConductorClient struct {
	CreateTaskFn func(projectID, title, description, governanceRef string) (string, error)
}

// CreateTask creates a new task in Conductor.
func (m *MockConductorClient) CreateTask(projectID, title, description, governanceRef string) (string, error) {
	if m.CreateTaskFn != nil {
		return m.CreateTaskFn(projectID, title, description, governanceRef)
	}
	return "task-mock-001", nil
}

// GetTask retrieves a task from Conductor.
func (m *MockConductorClient) GetTask(taskID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"id":     taskID,
		"status": "active",
	}, nil
}

// ListTasks returns all tasks in a project.
func (m *MockConductorClient) ListTasks(projectID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListDependencies returns dependencies for a task.
func (m *MockConductorClient) ListDependencies(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// ListActivity returns activity entries for a task.
func (m *MockConductorClient) ListActivity(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// CancelTask cancels a task in Conductor.
func (m *MockConductorClient) CancelTask(taskID string) error {
	return nil
}
