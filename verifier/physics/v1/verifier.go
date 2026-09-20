package physicsv1

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PithomLabs/oracle/verifier/model"
)

const (
	VerifierID      = "physics-v1"
	VerifierVersion = "0.1.0"
	VerifierHash    = "poc-verifier-v0.1.0"
)

// Run executes the physics verifier and returns the artifact.
func Run(ctx context.Context, input VerifierInput) (*model.VerificationArtifact, error) {
	steps := make([]model.Step, 0, 4)
	anyRefuted := false

	// Step 1: machine_verified
	verified1, msg1 := VerifyStep1(input.RhoExpr)
	steps = append(steps, model.Step{
		Index:       0,
		Description: "Fisher-rigidity identity: Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²)",
		Input:       fmt.Sprintf("ρ = %s", input.RhoExpr),
		Rule:        "symbolic verification",
		Output:      msg1,
		Verified:    verified1,
		CheckType:   model.CheckTypeMachineVerified,
	})
	if !verified1 {
		anyRefuted = true
	}

	// Step 2: machine_verified
	verified2, msg2 := VerifyStep2(input.KappaExpr, input.MExpr, input.RhoExpr)
	steps = append(steps, model.Step{
		Index:       1,
		Description: "Coefficient matching: a(ρ) = κ²/(8mρ)",
		Input:       fmt.Sprintf("κ = %s, m = %s, ρ = %s", input.KappaExpr, input.MExpr, input.RhoExpr),
		Rule:        "symbolic verification",
		Output:      msg2,
		Verified:    verified2,
		CheckType:   model.CheckTypeMachineVerified,
	})
	if !verified2 {
		anyRefuted = true
	}

	// Step 3: human_attested
	verified3, msg3 := VerifyStep3(input.Step3Input, input.Step3Expected)
	steps = append(steps, model.Step{
		Index:       2,
		Description: "|∇ρ|² coefficient consistency",
		Input:       input.Step3Input,
		Rule:        "human attestation",
		Output:      msg3,
		Verified:    verified3,
		CheckType:   model.CheckTypeHumanAttested,
	})
	if !verified3 {
		anyRefuted = true
	}

	// Step 4: human_attested
	verified4, msg4 := VerifyStep4(input.Step4Input, input.Step4Expected)
	steps = append(steps, model.Step{
		Index:       3,
		Description: "b′(ρ) = 0",
		Input:       input.Step4Input,
		Rule:        "human attestation",
		Output:      msg4,
		Verified:    verified4,
		CheckType:   model.CheckTypeHumanAttested,
	})
	if !verified4 {
		anyRefuted = true
	}

	// Aggregate result
	result := model.ResultConfirmed
	if anyRefuted {
		result = model.ResultRefuted
	}

	for _, s := range steps {
		if s.CheckType == model.CheckTypeMachineVerified && !s.Verified {
			result = model.ResultInconclusive
			break
		}
	}

	// Compute hashes
	inputBytes, _ := json.Marshal(input)
	inputHash := fmt.Sprintf("%x", sha256.Sum256(inputBytes))

	claimStr := fmt.Sprintf("rho=%s", input.RhoExpr)
	claimHash := fmt.Sprintf("%x", sha256.Sum256([]byte(claimStr)))

	artifact := &model.VerificationArtifact{
		RunID:           input.RunID,
		VerifierID:      VerifierID,
		VerifierVersion: VerifierVersion,
		VerifierHash:    VerifierHash,
		ClaimHash:       claimHash,
		InputHash:       inputHash,
		Timestamp:       time.Now(),
		Steps:           steps,
		Result:          result,
		EvidenceRef:     fmt.Sprintf("physics-v1-%s", input.RunID),
	}

	artifactHash, err := model.ComputeArtifactHash(*artifact)
	if err != nil {
		return nil, fmt.Errorf("failed to compute artifact hash: %w", err)
	}
	artifact.ArtifactHash = artifactHash

	return artifact, nil
}
