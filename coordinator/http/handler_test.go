package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	coordinator "github.com/PithomLabs/oracle/coordinator"
)

// testPackImpl implements domainpack.Pack for testing
type testPackImpl struct {
	PackID              string
	Version             string
	DebtVocab           []string
	EvidenceClasses     []string
	RetirementRules     map[string]domainpack.RetirementRule
}

func (p *testPackImpl) GetPackID() string              { return p.PackID }
func (p *testPackImpl) GetVersion() string             { return p.Version }
func (p *testPackImpl) GetDebtVocabulary() []string    { return p.DebtVocab }
func (p *testPackImpl) GetEvidenceClasses() []string   { return p.EvidenceClasses }
func (p *testPackImpl) GetRetirementRules() map[string]domainpack.RetirementRule {
	return p.RetirementRules
}
func (p *testPackImpl) GetVerifierSpecs() []domainpack.VerifierSpec { return nil }

// recordingSolvent captures calls
type recordingSolvent struct {
	dischargeCalls []dischargeRecord
}

type dischargeRecord struct {
	scenarioID    string
	beliefID      string
	obligationKey string
	instrumentRef string
	dischargedBy  string
}

func (s *recordingSolvent) CreateBelief(parentID, claim, scenarioID string) (string, error) {
	return "belief-mock", nil
}
func (s *recordingSolvent) CreateEdge(parentID, childID, kind string) error { return nil }
func (s *recordingSolvent) CreateEvidence(beliefID, provenanceClass, contentSHA256 string) (string, error) {
	return "evidence-mock", nil
}
func (s *recordingSolvent) PromoteBelief(beliefID, scenarioID string) error { return nil }
func (s *recordingSolvent) RetireDebt(beliefID, debtItem, scenarioID string) error {
	return nil
}
func (s *recordingSolvent) RetractBelief(beliefID, scenarioID string) error { return nil }
func (s *recordingSolvent) ApproveTarget(targetID string) error            { return nil }
func (s *recordingSolvent) GetBelief(beliefID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": beliefID, "status": "active", "debt": []string{"needMap"}}, nil
}
func (s *recordingSolvent) ListBeliefs(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) ListEvidence(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) ListEdges(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) ListIntents(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) ListActivities(scenarioID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) ListEvidenceForBelief(beliefID, scenarioID string) ([]map[string]interface{}, error) {
	// Return evidence matching the class for correct tests
	if beliefID == "belief-001" {
		return []map[string]interface{}{
			{"id": "ev-001", "belief_id": beliefID, "provenance_class": "reproducible_artifact"},
		}, nil
	}
	return []map[string]interface{}{}, nil
}
func (s *recordingSolvent) Discharge(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error {
	s.dischargeCalls = append(s.dischargeCalls, dischargeRecord{
		scenarioID: scenarioID, beliefID: beliefID,
		obligationKey: obligationKey, instrumentRef: instrumentRef, dischargedBy: dischargedBy,
	})
	return nil
}

type recordingConductor struct{}

func (c *recordingConductor) CreateTask(projectID, title, description, governanceRef string) (string, error) {
	return "task-mock", nil
}
func (c *recordingConductor) GetTask(taskID string) (map[string]interface{}, error) {
	return map[string]interface{}{"id": taskID, "status": "active"}, nil
}
func (c *recordingConductor) ListTasks(projectID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (c *recordingConductor) ListDependencies(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (c *recordingConductor) ListActivity(taskID string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (c *recordingConductor) CancelTask(taskID string) error { return nil }

func newTestPackRegistry(t *testing.T) *domainpack.PackRegistry {
	t.Helper()
	pack := &testPackImpl{
		PackID:          "bmist",
		Version:         "1.0.0",
		DebtVocab:       []string{"needMap", "needInvariant", "needNullModel"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		RetirementRules: map[string]domainpack.RetirementRule{
			"needMap":      {EvidenceClass: "reproducible_artifact", Rule: "map_check"},
			"needInvariant": {EvidenceClass: "reproducible_artifact", Rule: "invariant_check"},
			"needNullModel": {EvidenceClass: "operator_asserted", Rule: "scope_clarification"},
		},
	}
	registry := domainpack.NewRegistry()
	if err := registry.Register(pack); err != nil {
		t.Fatalf("register pack: %v", err)
	}
	return registry
}

func TestHandleSubmitDecision_EvidenceClassSurvivesHTTP(t *testing.T) {
	registry := newTestPackRegistry(t)
	solvent := &recordingSolvent{}
	conductor := &recordingConductor{}
	c := coordinator.NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)
	c.SetOperatorToken("test-token-123")
	handler := NewHandler(c)

	// Send retirement with WRONG evidence class for needMap (requires reproducible_artifact)
	body := map[string]interface{}{
		"type":           "RETIRE_DEBT",
		"belief_id":      "belief-001",
		"scenario_id":    "track-g0",
		"debt_item":      "needMap",
		"evidence_class": "operator_asserted", // WRONG: needMap requires reproducible_artifact
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/decisions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp SubmitDecisionResponse
	json.NewDecoder(w.Body).Decode(&resp)

	// If EvidenceClass survived the HTTP layer, the retirement should be REFUSED
	// because needMap requires reproducible_artifact, not operator_asserted
	if resp.Result != "refused" {
		t.Errorf("expected refused (wrong evidence class), got %s — EvidenceClass may have been dropped by HTTP handler", resp.Result)
	}

	// Verify Discharge was NOT called
	if len(solvent.dischargeCalls) != 0 {
		t.Errorf("Discharge should not be called when retirement is refused, got %d calls", len(solvent.dischargeCalls))
	}
}

func TestHandleSubmitDecision_CorrectEvidenceClassAllowed(t *testing.T) {
	registry := newTestPackRegistry(t)
	solvent := &recordingSolvent{}
	conductor := &recordingConductor{}
	c := coordinator.NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)
	c.SetOperatorToken("test-token-123")
	handler := NewHandler(c)

	// Send retirement with CORRECT evidence class for needMap
	body := map[string]interface{}{
		"type":           "RETIRE_DEBT",
		"belief_id":      "belief-001",
		"scenario_id":    "track-g0",
		"debt_item":      "needMap",
		"evidence_class": "reproducible_artifact", // CORRECT
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/decisions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp SubmitDecisionResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Result != "executed" {
		t.Errorf("expected executed (correct evidence class), got %s", resp.Result)
	}
}

func TestHandleSubmitDecision_MissingEvidenceClassRefused(t *testing.T) {
	registry := newTestPackRegistry(t)
	solvent := &recordingSolvent{}
	conductor := &recordingConductor{}
	c := coordinator.NewWithMockClientsAndPack(solvent, conductor, "test-operator", registry)
	c.SetOperatorToken("test-token-123")
	handler := NewHandler(c)

	// Send retirement with NO evidence class (empty string)
	body := map[string]interface{}{
		"type":        "RETIRE_DEBT",
		"belief_id":   "belief-001",
		"scenario_id": "track-g0",
		"debt_item":   "needMap",
		// No evidence_class field
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/decisions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token-123")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp SubmitDecisionResponse
	json.NewDecoder(w.Body).Decode(&resp)

	// Missing evidence class should be refused (mismatch with required class)
	if resp.Result != "refused" {
		t.Errorf("expected refused (missing evidence class), got %s", resp.Result)
	}
}
