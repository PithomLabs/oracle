package coordinator

// SphinxProjection is a read-only projection of Solvent authority state.
type SphinxProjection struct {
	Target            TargetInfo     `json:"target"`
	SupportingBelief  *BeliefInfo    `json:"supporting_belief,omitempty"`
	BeliefStatus      string         `json:"belief_status"`
	EvidenceSummary   []EvidenceInfo `json:"evidence_summary"`
	OpenDebt          []string       `json:"open_debt"`
	AuthorityState    string         `json:"authority_state"`
	JustificationState string        `json:"justification_state"`
	RequestState      string         `json:"request_state"`
	ApprovalState     string         `json:"approval_state"`
	IntentState       string         `json:"intent_state"`
	CurrentResult     string         `json:"current_result"`
	RefusalReason     string         `json:"refusal_reason,omitempty"`
}

// TargetInfo contains target information.
type TargetInfo struct {
	ID          string `json:"id"`
	ResourceType string `json:"resource_type"`
	ResourceID  string `json:"resource_id"`
}

// BeliefInfo contains belief information.
type BeliefInfo struct {
	ID     string `json:"id"`
	Claim  string `json:"claim"`
	Status string `json:"status"`
}

// EvidenceInfo contains evidence information.
type EvidenceInfo struct {
	ID               string `json:"id"`
	ProvenanceClass  string `json:"provenance_class"`
	ContentSHA256    string `json:"content_sha256"`
}

// authorizationContext returns the Sphinx projection.
func (c *Coordinator) authorizationContext(targetID string) (*SphinxProjection, error) {
	// Read-only projection of existing Solvent state
	projection := &SphinxProjection{
		Target: TargetInfo{
			ID: targetID,
		},
		BeliefStatus:      "unknown",
		AuthorityState:    "pending",
		JustificationState: "none",
		RequestState:      "none",
		ApprovalState:     "none",
		IntentState:       "none",
		CurrentResult:     "HUMAN_REVIEW",
	}

	// In production, this would read from Solvent
	// For POC, return a basic projection
	return projection, nil
}
