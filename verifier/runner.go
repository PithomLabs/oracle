package verifier

import (
	"context"
	"fmt"

	"github.com/PithomLabs/oracle/verifier/model"
	physicsv1 "github.com/PithomLabs/oracle/verifier/physics/v1"
)

// RunPhysicsVerifier invokes the physics verifier and registers the artifact.
// This is the ONLY trusted registration path.
func RunPhysicsVerifier(ctx context.Context, registry *ArtifactRegistry, input physicsv1.VerifierInput) (*model.VerificationArtifact, error) {
	artifact, err := physicsv1.Run(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("physics verifier failed: %w", err)
	}

	if err := registry.registerTrusted(*artifact); err != nil {
		return nil, fmt.Errorf("failed to register artifact: %w", err)
	}

	return artifact, nil
}
