package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

// VerificationArtifact is the generic verification artifact produced by any verifier.
type VerificationArtifact struct {
	RunID           string    `json:"run_id"`
	VerifierID      string    `json:"verifier_id"`
	VerifierVersion string    `json:"verifier_version"`
	VerifierHash    string    `json:"verifier_hash"`
	ClaimHash       string    `json:"claim_hash"`
	InputHash       string    `json:"input_hash"`
	Tolerance       string    `json:"tolerance,omitempty"`
	ArtifactHash    string    `json:"artifact_hash"`
	Timestamp       time.Time `json:"timestamp,omitempty"`
	Steps           []Step    `json:"steps"`
	Result          string    `json:"result"`
	EvidenceRef     string    `json:"evidence_ref"`
}

// Step represents a single verification step within an artifact.
type Step struct {
	Index       int    `json:"index"`
	Description string `json:"description"`
	Input       string `json:"input"`
	Rule        string `json:"rule"`
	Output      string `json:"output"`
	Verified    bool   `json:"verified"`
	CheckType   string `json:"check_type"` // "machine_verified" | "human_attested"
}

// Result constants
const (
	ResultConfirmed    = "confirmed"
	ResultRefuted      = "refuted"
	ResultInconclusive = "inconclusive"
)

// CheckType constants
const (
	CheckTypeMachineVerified = "machine_verified"
	CheckTypeHumanAttested  = "human_attested"
)

// ComputeArtifactHash computes SHA-256 of canonical artifact content,
// excluding only ArtifactHash (self-referential) and Timestamp (non-semantic).
// All semantic fields including VerifierID, VerifierVersion, Tolerance, and Steps
// are covered by the hash.
func ComputeArtifactHash(artifact VerificationArtifact) (string, error) {
	snapshot := artifact
	snapshot.ArtifactHash = ""
	snapshot.Timestamp = time.Time{}

	canonical, err := json.Marshal(snapshot)
	if err != nil {
		return "", fmt.Errorf("failed to marshal artifact for hashing: %w", err)
	}

	hash := sha256.Sum256(canonical)
	return fmt.Sprintf("%x", hash), nil
}
