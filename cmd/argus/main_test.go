package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/PithomLabs/oracle/verifier"
	"github.com/PithomLabs/oracle/verifier/physics/v1"
)

func TestVerifyCLI(t *testing.T) {
	// Create a temp input file
	input := physicsv1.DefaultInput()
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}

	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "verify_input.json")
	if err := os.WriteFile(inputFile, data, 0644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	// Run the verifier directly (simulating CLI behavior)
	registry := verifier.NewArtifactRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 30)
	defer cancel()

	artifact, err := verifier.RunPhysicsVerifier(ctx, registry, input)
	if err != nil {
		t.Fatalf("run verifier: %v", err)
	}

	// Encode to JSON (simulating CLI output)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(artifact); err != nil {
		t.Fatalf("encode artifact: %v", err)
	}

	// Verify output contains required fields
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("parse output: %v", err)
	}

	requiredFields := []string{"verifier_id", "verifier_version", "input_hash", "result", "steps"}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("missing required field: %s", field)
		}
	}

	// Verify verifier_id is correct
	if vid, ok := result["verifier_id"]; ok && vid != "physics-v1" {
		t.Errorf("expected verifier_id=physics-v1, got %v", vid)
	}

	t.Logf("verify CLI output contains all required fields: %v", requiredFields)
}
