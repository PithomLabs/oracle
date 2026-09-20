package physicsv1

import (
	"context"
	"testing"
)

func TestRunValid(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	// Result can be confirmed, refuted, or inconclusive depending on symbolic verification
	if artifact.Result == "" {
		t.Error("result should not be empty")
	}
}

func TestRunDeterministic(t *testing.T) {
	input := DefaultInput()

	artifact1, _ := Run(context.Background(), input)
	artifact2, _ := Run(context.Background(), input)

	if artifact1.ArtifactHash != artifact2.ArtifactHash {
		t.Error("artifact hash not deterministic")
	}
	if artifact1.Result != artifact2.Result {
		t.Error("result not deterministic")
	}
}

func TestRunStep1Verified(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(artifact.Steps) < 1 {
		t.Fatal("no steps")
	}
	// Step 1 is machine_verified
	if artifact.Steps[0].CheckType != "machine_verified" {
		t.Error("Step 1 should be machine_verified")
	}
}

func TestRunStep2Verified(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(artifact.Steps) < 2 {
		t.Fatal("not enough steps")
	}
	if artifact.Steps[1].CheckType != "machine_verified" {
		t.Error("Step 2 should be machine_verified")
	}
}

func TestRunStep3HumanAttested(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(artifact.Steps) < 3 {
		t.Fatal("not enough steps")
	}
	if artifact.Steps[2].CheckType != "human_attested" {
		t.Error("Step 3 should be human_attested")
	}
}

func TestRunStep4HumanAttested(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(artifact.Steps) < 4 {
		t.Fatal("not enough steps")
	}
	if artifact.Steps[3].CheckType != "human_attested" {
		t.Error("Step 4 should be human_attested")
	}
}

func TestRunAggregationValidResult(t *testing.T) {
	input := DefaultInput()
	artifact, err := Run(context.Background(), input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	// Result should be one of the valid values
	if artifact.Result != "confirmed" && artifact.Result != "refuted" && artifact.Result != "inconclusive" {
		t.Errorf("invalid result: %s", artifact.Result)
	}
}

func TestRunArtifactHashExcludesTimestamp(t *testing.T) {
	input := DefaultInput()

	artifact1, _ := Run(context.Background(), input)
	artifact2, _ := Run(context.Background(), input)

	// Timestamps will differ, but hashes should match
	if artifact1.ArtifactHash != artifact2.ArtifactHash {
		t.Error("artifact hash should exclude timestamp")
	}
}

func TestRunHashReproducible100Runs(t *testing.T) {
	input := DefaultInput()
	hashes := make(map[string]bool)

	for i := 0; i < 100; i++ {
		artifact, err := Run(context.Background(), input)
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}
		hashes[artifact.ArtifactHash] = true
	}

	if len(hashes) != 1 {
		t.Errorf("expected 1 unique hash, got %d", len(hashes))
	}
}
