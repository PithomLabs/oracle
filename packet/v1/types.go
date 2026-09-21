package packetv1

// Packet is the top-level EBP Research Packet v1.
type Packet struct {
	SchemaVersion string     `json:"schema_version"`
	Role          string     `json:"role"`
	PacketID      string     `json:"packet_id"`
	PackRef       string     `json:"pack_ref"`
	ScenarioID    string     `json:"scenario_id,omitempty"`
	RunID         string     `json:"run_id,omitempty"`
	TaskRef       string     `json:"task_ref,omitempty"`
	ProjectRef    string     `json:"project_ref,omitempty"`
	ScenarioRef   string     `json:"scenario_ref,omitempty"`
	Scope         string     `json:"scope,omitempty"`
	CorpusRef     string     `json:"corpus_ref,omitempty"`
	Agent         Agent      `json:"agent"`
	Beliefs       []Belief   `json:"beliefs"`
	Evidence      []Evidence `json:"evidence"`
	Edges         []Edge     `json:"edges,omitempty"`
	Tasks         []Task     `json:"tasks,omitempty"`
}

// Agent identifies the creator of the packet.
// All fields are required. Agent identity is provenance/observability metadata only,
// never an authority credential.
type Agent struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Harness string `json:"harness"`
	Model   string `json:"model"`
}

// Belief is an atomic epistemic object within the packet.
type Belief struct {
	LocalID   string   `json:"local_id"`
	Claim     string   `json:"claim"`
	ClaimType string   `json:"claim_type"`
	Debt      []string `json:"debt,omitempty"`
	InputSpec string   `json:"input_spec,omitempty"` // expected input hash for verification binding
}

// Evidence is a content-addressed artifact bound to a belief.
type Evidence struct {
	LocalID         string `json:"local_id"`
	BeliefRef       string `json:"belief_ref"`
	ProvenanceClass string `json:"provenance_class"`
	ContentSHA256   string `json:"content_sha256"`
	SourceURL       string `json:"source_url,omitempty"`
	ArtifactRef     string `json:"artifact_ref,omitempty"`
}

// Edge represents a relationship between two beliefs.
type Edge struct {
	LocalID string `json:"local_id"`
	FromRef string `json:"from_ref"`
	ToRef   string `json:"to_ref"`
	Kind    string `json:"kind"`
}

// Task is a proposed operational unit. Cannot carry authoritative state.
type Task struct {
	LocalID       string `json:"local_id"`
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	GovernanceRef string `json:"governance_ref,omitempty"`
}

// Reference prefixes
const (
	RefPrefixLocal     = "local:"
	RefPrefixCanonical = "canonical:belief:"
)

// Edge kinds
const (
	EdgeDerives     = "derives"
	EdgeContradicts = "contradicts"
)

// Roles
const (
	RoleWork        = "work"
	RoleAdversarial = "adversarial"
)

// Schema version
const SchemaVersion = "ebp-research-packet/v1"
