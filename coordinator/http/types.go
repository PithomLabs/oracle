package http

// SubmitDecisionRequest is the HTTP request for submitting a decision.
type SubmitDecisionRequest struct {
	Type          string `json:"type"`
	BeliefID      string `json:"belief_id"`
	ScenarioID    string `json:"scenario_id"`
	DebtItem      string `json:"debt_item,omitempty"`
	EvidenceClass string `json:"evidence_class,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

// SubmitDecisionResponse is the HTTP response for submitting a decision.
type SubmitDecisionResponse struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	BeliefID       string `json:"belief_id"`
	Result         string `json:"result"`
	RefusalReason  string `json:"refusal_reason,omitempty"`
	CreatedAt      string `json:"created_at"`
}

// PacketStatusResponse is the HTTP response for packet status.
type PacketStatusResponse struct {
	PacketID  string   `json:"packet_id"`
	Status    string   `json:"status"`
	BeliefIDs []string `json:"belief_ids"`
}

// DecisionContextResponse is the HTTP response for decision context.
type DecisionContextResponse struct {
	BeliefID    string   `json:"belief_id"`
	Claim       string   `json:"claim"`
	Status      string   `json:"status"`
	Debt        []string `json:"debt"`
	EvidenceIDs []string `json:"evidence_ids"`
}

// AuthorizationContextResponse is the HTTP response for authorization context.
type AuthorizationContextResponse struct {
	Target            interface{} `json:"target"`
	BeliefStatus      string      `json:"belief_status"`
	AuthorityState    string      `json:"authority_state"`
	CurrentResult     string      `json:"current_result"`
	RefusalReason     string      `json:"refusal_reason,omitempty"`
}

// ErrorResponse is the HTTP error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
