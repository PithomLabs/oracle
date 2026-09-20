package coordinator

import (
	"fmt"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	"github.com/PithomLabs/oracle/verifier"
)

// Config holds the Coordinator configuration.
type Config struct {
	SolventBaseURL   string
	ConductorBaseURL string
	PackRegistry     *domainpack.PackRegistry
	ArtifactReader   verifier.ArtifactReader
	OperatorID       string
}

// Coordinator is the deterministic packet compiler.
type Coordinator struct {
	solventClient    SolventClientInterface
	conductorClient  ConductorClientInterface
	packRegistry     *domainpack.PackRegistry
	artifactReader   verifier.ArtifactReader
	projectionQueue  *ProjectionQueue
	idempotencyCache *IdempotencyCache
	operatorID       string
	operatorToken    string
}

// New creates a new Coordinator.
func New(cfg Config) (*Coordinator, error) {
	if cfg.SolventBaseURL == "" {
		return nil, fmt.Errorf("solvent_base_url is required")
	}
	if cfg.ConductorBaseURL == "" {
		return nil, fmt.Errorf("conductor_base_url is required")
	}
	if cfg.PackRegistry == nil {
		return nil, fmt.Errorf("pack_registry is required")
	}
	if cfg.ArtifactReader == nil {
		return nil, fmt.Errorf("artifact_reader is required")
	}
	if cfg.OperatorID == "" {
		return nil, fmt.Errorf("operator_id is required (ARGUS_OPERATOR_PRINCIPAL_ID)")
	}

	return &Coordinator{
		solventClient:    NewSolventClient(cfg.SolventBaseURL),
		conductorClient:  NewConductorClient(cfg.ConductorBaseURL),
		packRegistry:     cfg.PackRegistry,
		artifactReader:   cfg.ArtifactReader,
		projectionQueue:  NewProjectionQueue(),
		idempotencyCache: NewIdempotencyCache(),
		operatorID:       cfg.OperatorID,
	}, nil
}

// NewWithMockClients creates a new Coordinator with mock clients for testing.
func NewWithMockClients(solvent SolventClientInterface, conductor ConductorClientInterface, operatorID string) *Coordinator {
	return &Coordinator{
		solventClient:    solvent,
		conductorClient:  conductor,
		projectionQueue:  NewProjectionQueue(),
		idempotencyCache: NewIdempotencyCache(),
		operatorID:       operatorID,
	}
}

// NewWithMockClientsAndPack creates a new Coordinator with mock clients and pack registry for testing.
func NewWithMockClientsAndPack(solvent SolventClientInterface, conductor ConductorClientInterface, operatorID string, packRegistry *domainpack.PackRegistry) *Coordinator {
	return &Coordinator{
		solventClient:    solvent,
		conductorClient:  conductor,
		packRegistry:     packRegistry,
		projectionQueue:  NewProjectionQueue(),
		idempotencyCache: NewIdempotencyCache(),
		operatorID:       operatorID,
	}
}

// CompilePacket is the core compilation entry point.
func (c *Coordinator) CompilePacket(pkt *packetv1.Packet) (*CompilationResult, error) {
	return c.compile(pkt)
}

// SubmitDecision handles human decisions.
func (c *Coordinator) SubmitDecision(req DecisionRequest) (*DecisionRecord, error) {
	return c.submitDecision(req)
}

// PacketStatus returns the status of a compiled packet.
func (c *Coordinator) PacketStatus(packetID string) (*PacketStatus, error) {
	return c.packetStatus(packetID)
}

// DecisionContext returns decision context for a belief.
func (c *Coordinator) DecisionContext(beliefID string) (*DecisionContext, error) {
	return c.decisionContext(beliefID)
}

// AuthorizationContext returns the Sphinx projection.
func (c *Coordinator) AuthorizationContext(targetID string) (*SphinxProjection, error) {
	return c.authorizationContext(targetID)
}

// GetPackRules returns retirement rules from all registered packs.
func (c *Coordinator) GetPackRules() map[string]interface{} {
	result := make(map[string]interface{})
	if c.packRegistry == nil {
		return result
	}
	pack, err := c.packRegistry.Get("bmist", "1.0.0")
	if err != nil {
		return result
	}
	rules := pack.GetRetirementRules()
	if rules == nil {
		return result
	}
	result["rules"] = rules
	return result
}

// RetryPending retries failed Conductor projections.
func (c *Coordinator) RetryPending() error {
	return c.projectionQueue.RetryPending()
}

// SetOperatorToken configures the bearer token for /decisions authentication.
func (c *Coordinator) SetOperatorToken(token string) {
	c.operatorToken = token
}

// AuthenticateOperator validates a bearer token against the configured operator token.
// Returns the operator ID if valid, error otherwise.
func (c *Coordinator) AuthenticateOperator(bearerToken string) (string, error) {
	if c.operatorToken == "" {
		return "", fmt.Errorf("operator token not configured")
	}
	if bearerToken != c.operatorToken {
		return "", fmt.Errorf("invalid operator token")
	}
	return c.operatorID, nil
}
