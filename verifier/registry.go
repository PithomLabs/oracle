package verifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/PithomLabs/oracle/verifier/model"
)

// ArtifactReader is the read-only interface exposed to external consumers (e.g., Coordinator).
type ArtifactReader interface {
	Resolve(ctx context.Context, ref string) (model.VerificationArtifact, error)
	VerifyHash(ref string, expectedSHA256 string) error
}

// ArtifactRegistry is an in-memory, process-local registry of verification artifacts.
type ArtifactRegistry struct {
	mu        sync.RWMutex
	artifacts map[string]model.VerificationArtifact // key = EvidenceRef
}

// NewArtifactRegistry creates an empty ArtifactRegistry.
func NewArtifactRegistry() *ArtifactRegistry {
	return &ArtifactRegistry{
		artifacts: make(map[string]model.VerificationArtifact),
	}
}

// registerTrusted adds a verified artifact to the registry.
// This method is unexported — only callable from within the verifier/ package.
// External packages cannot register artifacts.
func (r *ArtifactRegistry) registerTrusted(artifact model.VerificationArtifact) error {
	if artifact.EvidenceRef == "" {
		return fmt.Errorf("artifact evidence_ref is required")
	}

	if artifact.ArtifactHash == "" {
		return fmt.Errorf("artifact hash is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.artifacts[artifact.EvidenceRef]; exists {
		return fmt.Errorf("artifact with evidence_ref %q already registered", artifact.EvidenceRef)
	}

	r.artifacts[artifact.EvidenceRef] = artifact
	return nil
}

// Resolve retrieves an artifact by its evidence reference.
func (r *ArtifactRegistry) Resolve(ctx context.Context, ref string) (model.VerificationArtifact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artifact, ok := r.artifacts[ref]
	if !ok {
		return model.VerificationArtifact{}, fmt.Errorf("artifact with evidence_ref %q not found", ref)
	}
	return artifact, nil
}

// VerifyHash checks that the artifact's hash matches the expected value.
func (r *ArtifactRegistry) VerifyHash(ref string, expectedSHA256 string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	artifact, ok := r.artifacts[ref]
	if !ok {
		return fmt.Errorf("artifact with evidence_ref %q not found", ref)
	}

	if artifact.ArtifactHash != expectedSHA256 {
		return fmt.Errorf("hash mismatch for evidence_ref %q: expected %s, got %s", ref, expectedSHA256, artifact.ArtifactHash)
	}

	return nil
}

// Register adds an artifact to the registry. Intended for test use
// where RunPhysicsVerifier is not available.
func (r *ArtifactRegistry) Register(artifact model.VerificationArtifact) error {
	return r.registerTrusted(artifact)
}

// AsReader returns a read-only view of the registry for external consumers.
func (r *ArtifactRegistry) AsReader() ArtifactReader {
	return r
}
