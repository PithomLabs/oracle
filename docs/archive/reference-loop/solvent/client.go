package solvent

import (
	"context"
	"database/sql"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Client is a REST client for the Solvent API
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	DB         *sql.DB
}

// NewClient creates a new Solvent REST client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Principal represents a Solvent principal
type Principal struct {
	PrincipalID   string `json:"principal_id"`
	PrincipalType string `json:"principal_type"`
	Issuer        string `json:"issuer"`
	RevokedAt     string `json:"revoked_at,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// CreatePrincipal creates a new principal
func (c *Client) CreatePrincipal(ctx context.Context, principalType, issuer string) (*Principal, error) {
	reqBody := map[string]interface{}{
		"principal_type": principalType,
		"issuer":         issuer,
	}
	resp, err := c.doRequest(ctx, "POST", "/v1/principals", reqBody)
	if err != nil {
		return nil, err
	}
	var p Principal
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Target represents an authority target
type Target struct {
	TargetID              string `json:"target_id"`
	PrincipalID           string `json:"principal_id"`
	ResourceType          string `json:"resource_type"`
	ResourceID            string `json:"resource_id"`
	Scope                 string `json:"scope"`
	ActionNamespace       string `json:"action_namespace"`
	ActionName            string `json:"action_name"`
	ConsequenceType       string `json:"consequence_type"`
	ConsequenceParameters json.RawMessage `json:"consequence_parameters"`
	CreatedBy             string `json:"created_by"`
	State                 string `json:"state"`
	CreatedAt             string `json:"created_at"`
	RequestedAt           string `json:"requested_at,omitempty"`
	RequestedBy           string `json:"requested_by,omitempty"`
}

// CreateTarget creates a new authority target
func (c *Client) CreateTarget(ctx context.Context, req CreateTargetRequest) (*Target, error) {
	resp, err := c.doRequest(ctx, "POST", "/v1/targets", req)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTargetRequest is the request body for creating a target
type CreateTargetRequest struct {
	PrincipalID           string          `json:"principal_id"`
	ResourceType          string          `json:"resource_type"`
	ResourceID            string          `json:"resource_id"`
	Scope                 string          `json:"scope"`
	ActionNamespace       string          `json:"action_namespace"`
	ActionName            string          `json:"action_name"`
	ConsequenceType       string          `json:"consequence_type"`
	ConsequenceParameters json.RawMessage `json:"consequence_parameters"`
	CreatedBy             string          `json:"created_by"`
}

// AttachJustification attaches a justification to a target
func (c *Client) AttachJustification(ctx context.Context, targetID, beliefID, beliefStatus, attachedBy string) (*Target, error) {
	reqBody := map[string]interface{}{
		"instrument_ref": beliefID,
	}
	url := fmt.Sprintf("/v1/targets/%s/justifications?belief_id=%s", targetID, beliefID)
	resp, err := c.doRequestWithQuery(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// RequestAuthorization requests authorization for a target
func (c *Client) RequestAuthorization(ctx context.Context, targetID, requestedBy string) (*Target, error) {
	reqBody := map[string]interface{}{
		"requested_by": requestedBy,
	}
	resp, err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/targets/%s/request", targetID), reqBody)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// ApproveTarget approves a target
func (c *Client) ApproveTarget(ctx context.Context, targetID, approvedBy, approvalPin string) (*Target, error) {
	reqBody := map[string]interface{}{
		"approval_pin": approvalPin,
	}
	resp, err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/targets/%s/approve", targetID), reqBody)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// GetTarget gets a target by ID
func (c *Client) GetTarget(ctx context.Context, targetID string) (*Target, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/targets/%s", targetID), nil)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// Belief represents a Solvent belief
type Belief struct {
	BeliefID      string   `json:"belief_id"`
	ScenarioID    string   `json:"scenario_id"`
	Claim         string   `json:"claim"`
	ClaimType     string   `json:"claim_type"`
	Status        string   `json:"status"`
	Debt          []string `json:"debt"`
	FinalTruth    bool     `json:"final_truth"`
	EvidenceCount int      `json:"evidence_count"`
	IntentCount   int      `json:"intent_count"`
	CreatedAt     string   `json:"created_at"`
}

// EnterBelief enters a new belief
func (c *Client) EnterBelief(ctx context.Context, scenarioID, claim, claimType string) (*Belief, error) {
	reqBody := map[string]interface{}{
		"scenario_id": scenarioID,
		"claim":       claim,
		"claim_type":  claimType,
	}
	resp, err := c.doRequest(ctx, "POST", "/v1/beliefs", reqBody)
	if err != nil {
		return nil, err
	}
	var b Belief
	if err := json.Unmarshal(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// GetBelief gets a belief by ID
func (c *Client) GetBelief(ctx context.Context, scenarioID, beliefID string) (*Belief, error) {
	url := fmt.Sprintf("/v1/beliefs/%s?scenario_id=%s", beliefID, scenarioID)
	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var b Belief
	if err := json.Unmarshal(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// PromoteBelief promotes a belief
func (c *Client) PromoteBelief(ctx context.Context, scenarioID, beliefID string) (*Belief, error) {
	url := fmt.Sprintf("/v1/beliefs/%s/promote?scenario_id=%s", beliefID, scenarioID)
	resp, err := c.doRequest(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	var b Belief
	if err := json.Unmarshal(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// RetireDebt retires a debt item
func (c *Client) RetireDebt(ctx context.Context, scenarioID, beliefID, debtItem string) (*Belief, error) {
	url := fmt.Sprintf("/v1/beliefs/%s/debt/retire?scenario_id=%s", beliefID, scenarioID)
	reqBody := map[string]interface{}{
		"debt_item": debtItem,
	}
	resp, err := c.doRequest(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	var b Belief
	if err := json.Unmarshal(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// AuthorizeActionResult represents the result of authorize action
type AuthorizeActionResult struct {
	BeliefID    string `json:"belief_id"`
	IntentID    string `json:"intent_id,omitempty"`
	IntentState string `json:"intent_state"`
	Action      string `json:"action"`
	Authority   struct {
		TargetID string `json:"target_id"`
		Allowed  bool   `json:"allowed"`
		Reason   string `json:"reason"`
	} `json:"authority"`
}

// AuthorizeAction authorizes an action and creates an intent
func (c *Client) AuthorizeAction(ctx context.Context, req AuthorizeActionRequest) (*AuthorizeActionResult, error) {
	resp, err := c.doRequest(ctx, "POST", "/v1/authorizations/action", req)
	if err != nil {
		return nil, err
	}
	var result AuthorizeActionResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// AuthorizeActionRequest is the request body for authorize action
type AuthorizeActionRequest struct {
	ScenarioID   string `json:"scenario_id"`
	BeliefID     string `json:"belief_id"`
	Action       string `json:"action"`
	ActionSource string `json:"action_source"` // "user_typed"
	TargetID     string `json:"target_id"`
	ActorID      string `json:"actor_id"`
}

// ExecuteActionResult represents the result of execute action
type ExecuteActionResult struct {
	BeliefID    string `json:"belief_id"`
	IntentID    string `json:"intent_id,omitempty"`
	IntentState string `json:"intent_state,omitempty"`
	Action      string `json:"action"`
	Allowed     bool   `json:"allowed"`
	Success     bool   `json:"success"`
	Output      string `json:"output,omitempty"`
	Error       string `json:"error,omitempty"`
	Reason      string `json:"reason,omitempty"`
	ExecutedAt  string `json:"executed_at"`
}

// ExecuteAction executes an authorized action
func (c *Client) ExecuteAction(ctx context.Context, req ExecuteActionRequest) (*ExecuteActionResult, error) {
	resp, err := c.doRequest(ctx, "POST", "/v1/authorizations/execute", req)
	if err != nil {
		return nil, err
	}
	var result ExecuteActionResult
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecuteActionRequest is the request body for execute action
type ExecuteActionRequest struct {
	ScenarioID      string `json:"scenario_id"`
	BeliefID        string `json:"belief_id"`
	Action          string `json:"action"`
	TargetID        string `json:"target_id"`
	IntentID        string `json:"intent_id"`
	ConsequenceType string `json:"consequence_type"`
}

// ReconcileIntent reconciles an intent
func (c *Client) ReconcileIntent(ctx context.Context, scenarioID, intentID, outcome, operatorID string) error {
	reqBody := map[string]interface{}{
		"scenario_id": scenarioID,
		"intent_id":   intentID,
		"outcome":     outcome,
		"operator_id": operatorID,
	}
	_, err := c.doRequest(ctx, "POST", "/v1/authorizations/reconcile", reqBody)
	return err
}

// ActivityEntry represents an audit activity entry
type ActivityEntry struct {
	ID             string                 `json:"id"`
	ScenarioID     string                 `json:"scenario_id"`
	Type           string                 `json:"type"`
	ActorID        string                 `json:"actor_id"`
	SubjectID      string                 `json:"subject_id"`
	Details        map[string]interface{} `json:"details"`
	SQLState       string                 `json:"sqlstate,omitempty"`
	ConstraintName string                 `json:"constraint_name,omitempty"`
	Refusal        bool                   `json:"refusal"`
	CreatedAt      string                 `json:"created_at"`
}

// GetActivity gets audit activity for a scenario
func (c *Client) GetActivity(ctx context.Context, scenarioID, activityType string, limit int) ([]ActivityEntry, error) {
	url := fmt.Sprintf("/v1/activity?scenario_id=%s&limit=%d", scenarioID, limit)
	if activityType != "" {
		url += "&type=" + activityType
	}
	resp, err := c.doRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Activities []ActivityEntry `json:"activities"`
		Total      int             `json:"total"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, err
	}
	return result.Activities, nil
}

// GetLiveIntent queries the database for a live intent matching the given criteria.
// This is a workaround for the frozen Solvent REST API not returning intent_id
// from AuthorizeAction.
func (c *Client) GetLiveIntent(ctx context.Context, scenarioID, beliefID, action, targetID string) (string, error) {
	if c.DB == nil {
		return "", fmt.Errorf("client DB not configured")
	}
	var intentID string
	err := c.DB.QueryRowContext(ctx, `
		SELECT id FROM action_intent
		WHERE scenario_id = $1::UUID
		  AND belief_id = $2::UUID
		  AND action = $3
		  AND target_id = $4::UUID
		  AND state = 'live'
		
		LIMIT 1
	`, scenarioID, beliefID, action, targetID).Scan(&intentID)
	if err != nil {
		return "", fmt.Errorf("query live intent: %w", err)
	}
	return intentID, nil
}

// doRequest performs an HTTP request to the Solvent API
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	return c.doRequestWithQuery(ctx, method, path, body)
}

// doRequestWithQuery performs an HTTP request with query parameters in path
func (c *Client) doRequestWithQuery(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("solvent API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// RevokeTarget revokes an authority target
func (c *Client) RevokeTarget(ctx context.Context, targetID, revokedBy, reason string) (*Target, error) {
	reqBody := map[string]interface{}{
		"revoked_by": revokedBy,
		"reason":     reason,
	}
	resp, err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/targets/%s/revoke", targetID), reqBody)
	if err != nil {
		return nil, err
	}
	var t Target
	if err := json.Unmarshal(resp, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// RetractBelief retracts a belief and cascades to cancel live intents
func (c *Client) RetractBelief(ctx context.Context, scenarioID, beliefID string) (*Belief, error) {
	url := fmt.Sprintf("/v1/beliefs/%s/retract?scenario_id=%s", beliefID, scenarioID)
	resp, err := c.doRequest(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	var b Belief
	if err := json.Unmarshal(resp, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
