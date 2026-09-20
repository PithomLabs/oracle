package agent

import (
	"strings"
	"testing"
)

func TestOperationIdentityDeterministic(t *testing.T) {
	repo := "ibmendoza/reference-loop-test"
	workflow := "ref-loop.yml"
	ref := "main"
	runID := "run-abc-123"

	first := operationID(repo, workflow, ref, runID)
	second := operationID(repo, workflow, ref, runID)

	if first != second {
		t.Errorf("operation identity not deterministic: %q != %q", first, second)
	}
	if first != "deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:run-abc-123" {
		t.Errorf("unexpected operation identity: %q", first)
	}
}

func TestOperationIdentityIncludesRunID(t *testing.T) {
	runID := "run-xyz-789"
	opID := operationID("ibmendoza/reference-loop-test", "ref-loop.yml", "main", runID)

	parts := strings.Split(opID, ":")
	if len(parts) != 5 {
		t.Fatalf("expected 5 colon-separated parts, got %d: %q", len(parts), opID)
	}
	if parts[0] != "deploy" {
		t.Errorf("expected operation 'deploy', got %q", parts[0])
	}
	if parts[1] != "ibmendoza/reference-loop-test" {
		t.Errorf("expected repo 'ibmendoza/reference-loop-test', got %q", parts[1])
	}
	if parts[2] != "ref-loop.yml" {
		t.Errorf("expected workflow 'ref-loop.yml', got %q", parts[2])
	}
	if parts[3] != "main" {
		t.Errorf("expected ref 'main', got %q", parts[3])
	}
	if parts[4] != runID {
		t.Errorf("expected run_id %q at end, got %q", runID, parts[4])
	}
}
