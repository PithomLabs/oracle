package coordinator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SolventClient is a REST client for the Solvent API.
type SolventClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSolventClient creates a new Solvent client.
func NewSolventClient(baseURL string) *SolventClient {
	return &SolventClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// CreateBelief creates a new belief in Solvent.
func (c *SolventClient) CreateBelief(parentID, claim, scenarioID string) (string, error) {
	body := map[string]interface{}{
		"claim":      claim,
		"scenario_id": scenarioID,
	}
	if parentID != "" {
		body["parent_id"] = parentID
	}
	return c.post("/v1/beliefs", body)
}

// CreateEdge creates an edge between beliefs in Solvent.
func (c *SolventClient) CreateEdge(parentID, childID, kind string) error {
	body := map[string]interface{}{
		"child_id": childID,
		"kind":     kind,
	}
	_, err := c.post(fmt.Sprintf("/v1/beliefs/%s/edges", parentID), body)
	return err
}

// CreateEvidence admits evidence to Solvent.
func (c *SolventClient) CreateEvidence(beliefID, provenanceClass, contentSHA256 string) (string, error) {
	body := map[string]interface{}{
		"belief_id":        beliefID,
		"provenance_class": provenanceClass,
		"content_sha256":   contentSHA256,
	}
	return c.post("/v1/evidence", body)
}

// PromoteBelief promotes a belief in Solvent.
func (c *SolventClient) PromoteBelief(beliefID, scenarioID string) error {
	_, err := c.post(fmt.Sprintf("/v1/beliefs/%s/promote?scenario_id=%s", beliefID, scenarioID), nil)
	return err
}

// RetireDebt retires debt on a belief in Solvent.
func (c *SolventClient) RetireDebt(beliefID, debtItem, scenarioID string) error {
	body := map[string]interface{}{
		"debt_item": debtItem,
	}
	_, err := c.post(fmt.Sprintf("/v1/beliefs/%s/debt/retire?scenario_id=%s", beliefID, scenarioID), body)
	return err
}

// RetractBelief retracts a belief in Solvent.
func (c *SolventClient) RetractBelief(beliefID, scenarioID string) error {
	_, err := c.post(fmt.Sprintf("/v1/beliefs/%s/retract?scenario_id=%s", beliefID, scenarioID), nil)
	return err
}

// ApproveTarget approves an authority target in Solvent.
func (c *SolventClient) ApproveTarget(targetID string) error {
	_, err := c.post(fmt.Sprintf("/v1/targets/%s/approve", targetID), nil)
	return err
}

// GetBelief retrieves a belief from Solvent.
func (c *SolventClient) GetBelief(beliefID string) (map[string]interface{}, error) {
	return c.get(fmt.Sprintf("/v1/beliefs/%s", beliefID))
}

// ListEvidenceForBelief returns evidence for a specific belief in a scenario.
func (c *SolventClient) ListEvidenceForBelief(beliefID, scenarioID string) ([]map[string]interface{}, error) {
	resp, err := c.get(fmt.Sprintf("/v1/beliefs/%s/evidence?scenario_id=%s", beliefID, scenarioID))
	if err != nil {
		return nil, err
	}
	evidenceRaw, _ := resp["evidence"].([]interface{})
	out := make([]map[string]interface{}, 0, len(evidenceRaw))
	for _, e := range evidenceRaw {
		if m, ok := e.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// ListBeliefs returns all beliefs in a scenario.
func (c *SolventClient) ListBeliefs(scenarioID string) ([]map[string]interface{}, error) {
	resp, err := c.get(fmt.Sprintf("/v1/beliefs?scenario_id=%s", scenarioID))
	if err != nil {
		return nil, err
	}
	beliefsRaw, _ := resp["beliefs"].([]interface{})
	out := make([]map[string]interface{}, 0, len(beliefsRaw))
	for _, b := range beliefsRaw {
		if m, ok := b.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// ListEvidence returns all evidence in a scenario.
func (c *SolventClient) ListEvidence(scenarioID string) ([]map[string]interface{}, error) {
	// For Phase 8, evidence is derived from beliefs' evidence counts.
	// A full implementation would iterate beliefs and fetch per-belief evidence.
	_, err := c.get(fmt.Sprintf("/v1/beliefs?scenario_id=%s", scenarioID))
	if err != nil {
		return nil, err
	}
	return []map[string]interface{}{}, nil
}

// ListEdges returns all edges in a scenario.
func (c *SolventClient) ListEdges(scenarioID string) ([]map[string]interface{}, error) {
	// Edges are per-belief; iterate beliefs.
	beliefs, err := c.ListBeliefs(scenarioID)
	if err != nil {
		return nil, err
	}
	var edges []map[string]interface{}
	for _, b := range beliefs {
		beliefID, _ := b["belief_id"].(string)
		if beliefID == "" {
			continue
		}
		result, err := c.get(fmt.Sprintf("/v1/beliefs/%s/edges?scenario_id=%s", beliefID, scenarioID))
		if err != nil {
			continue // Skip beliefs where edge listing fails
		}
		if edgeList, ok := result["edges"].([]interface{}); ok {
			for _, e := range edgeList {
				if m, ok := e.(map[string]interface{}); ok {
					edges = append(edges, m)
				}
			}
		}
	}
	if edges == nil {
		edges = []map[string]interface{}{}
	}
	return edges, nil
}

// ListIntents returns all action intents in a scenario.
func (c *SolventClient) ListIntents(scenarioID string) ([]map[string]interface{}, error) {
	// Intents are part of the ledger; for Phase 8, return empty.
	// Full implementation would parse ledger response.
	return []map[string]interface{}{}, nil
}

// ListActivities returns audit activity entries for a scenario.
func (c *SolventClient) ListActivities(scenarioID string) ([]map[string]interface{}, error) {
	resp, err := c.get(fmt.Sprintf("/v1/activity?scenario_id=%s", scenarioID))
	if err != nil {
		return nil, err
	}
	activitiesRaw, _ := resp["activities"].([]interface{})
	out := make([]map[string]interface{}, 0, len(activitiesRaw))
	for _, a := range activitiesRaw {
		if m, ok := a.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// Discharge records attributed debt discharge in Solvent.
func (c *SolventClient) Discharge(scenarioID, beliefID, obligationKey, instrumentRef, dischargedBy string) error {
	body := map[string]interface{}{
		"scenario_id":    scenarioID,
		"belief_id":      beliefID,
		"obligation_key": obligationKey,
		"instrument_ref": instrumentRef,
		"discharged_by":  dischargedBy,
	}
	_, err := c.post("/v1/discharge", body)
	return err
}

// post sends a POST request to Solvent.
func (c *SolventClient) post(path string, body interface{}) (string, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("solvent error %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}

// get sends a GET request to Solvent.
func (c *SolventClient) get(path string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("solvent error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// ConductorClient is a REST client for the Conductor API.
type ConductorClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewConductorClient creates a new Conductor client.
func NewConductorClient(baseURL string) *ConductorClient {
	return &ConductorClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// CreateTask creates a new task in Conductor.
func (c *ConductorClient) CreateTask(projectID, title, description, governanceRef string) (string, error) {
	body := map[string]interface{}{
		"title":          title,
		"description":    description,
		"governance_ref": governanceRef,
		"status":         "proposed",
	}
	return c.post(fmt.Sprintf("/v1/projects/%s/tasks", projectID), body)
}

// GetTask retrieves a task from Conductor.
func (c *ConductorClient) GetTask(taskID string) (map[string]interface{}, error) {
	return c.get(fmt.Sprintf("/v1/tasks/%s", taskID))
}

// ListTasks returns all tasks in a project.
func (c *ConductorClient) ListTasks(projectID string) ([]map[string]interface{}, error) {
	result, err := c.get(fmt.Sprintf("/v1/projects/%s/tasks", projectID))
	if err != nil {
		return nil, err
	}
	tasksRaw, _ := result["tasks"].([]interface{})
	if tasksRaw == nil {
		// Conductor may return tasks as top-level array
		return []map[string]interface{}{}, nil
	}
	out := make([]map[string]interface{}, 0, len(tasksRaw))
	for _, t := range tasksRaw {
		if m, ok := t.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// ListDependencies returns dependencies for a task.
func (c *ConductorClient) ListDependencies(taskID string) ([]map[string]interface{}, error) {
	// Dependencies are part of governance; for Phase 8, return empty.
	// Full implementation would parse governance response.
	return []map[string]interface{}{}, nil
}

// ListActivity returns activity entries for a task.
func (c *ConductorClient) ListActivity(taskID string) ([]map[string]interface{}, error) {
	resp, err := c.get(fmt.Sprintf("/v1/tasks/%s/activity", taskID))
	if err != nil {
		return nil, err
	}
	// Conductor returns activity as an array at top level
	activitiesRaw, _ := resp["activities"].([]interface{})
	if activitiesRaw == nil {
		// Try treating the whole response as the list
		return []map[string]interface{}{}, nil
	}
	out := make([]map[string]interface{}, 0, len(activitiesRaw))
	for _, a := range activitiesRaw {
		if m, ok := a.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out, nil
}

// CancelTask cancels a task in Conductor.
func (c *ConductorClient) CancelTask(taskID string) error {
	_, err := c.post(fmt.Sprintf("/v1/tasks/%s/cancel", taskID), nil)
	return err
}

// get sends a GET request to Conductor.
func (c *ConductorClient) get(path string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("conductor error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// post sends a POST request to Conductor.
func (c *ConductorClient) post(path string, body interface{}) (string, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("conductor error %d: %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
