package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Supported manifest versions
const SupportedManifestVersion = "1.0.0"

// Supported provenance classes
var SupportedProvenanceClasses = map[string]bool{
	"operator_asserted": true,
	"agent_derived":     true,
	"reproducible_artifact": true,
}

// Manifest is the top-level corpus manifest structure.
type Manifest struct {
	ManifestVersion string           `json:"manifest_version"`
	CorpusID        string           `json:"corpus_id"`
	Artifacts       []ArtifactRecord `json:"artifacts"`
}

// ArtifactRecord describes a single corpus artifact.
type ArtifactRecord struct {
	ID              string `json:"id"`
	Path            string `json:"path"`
	SHA256          string `json:"sha256"`
	ProvenanceClass string `json:"provenance_class"`
	Description     string `json:"description"`
	StableLocator   string `json:"stable_locator"`
}

// Load reads and parses a manifest file.
func Load(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// Validate performs structural validation of the manifest.
func (m *Manifest) Validate() error {
	if m.ManifestVersion != SupportedManifestVersion {
		return fmt.Errorf("unsupported manifest version: %s", m.ManifestVersion)
	}

	if m.CorpusID == "" {
		return fmt.Errorf("corpus_id is required")
	}

	if len(m.Artifacts) == 0 {
		return fmt.Errorf("artifacts list is empty")
	}

	seen := make(map[string]bool)
	for i, a := range m.Artifacts {
		if a.ID == "" {
			return fmt.Errorf("artifact[%d]: id is required", i)
		}
		if seen[a.ID] {
			return fmt.Errorf("artifact[%d]: duplicate id %q", i, a.ID)
		}
		seen[a.ID] = true

		if a.Path == "" {
			return fmt.Errorf("artifact[%d]: path is required", i)
		}
		if a.SHA256 == "" {
			return fmt.Errorf("artifact[%d]: sha256 is required", i)
		}
		if len(a.SHA256) != 64 {
			return fmt.Errorf("artifact[%d]: invalid sha256 length", i)
		}
		if _, err := hex.DecodeString(a.SHA256); err != nil {
			return fmt.Errorf("artifact[%d]: invalid sha256 hex: %w", i, err)
		}
		if !SupportedProvenanceClasses[a.ProvenanceClass] {
			return fmt.Errorf("artifact[%d]: unsupported provenance_class %q", i, a.ProvenanceClass)
		}
		if a.StableLocator == "" {
			return fmt.Errorf("artifact[%d]: stable_locator is required", i)
		}
	}

	return nil
}

// Verify checks that all referenced files exist and their hashes match.
func (m *Manifest) Verify(root string) error {
	for i, a := range m.Artifacts {
		// Path safety: reject absolute paths and traversal
		if filepath.IsAbs(a.Path) {
			return fmt.Errorf("artifact[%d]: absolute path not allowed: %s", i, a.Path)
		}
		if strings.Contains(a.Path, "..") {
			return fmt.Errorf("artifact[%d]: path traversal not allowed: %s", i, a.Path)
		}

		fullPath := filepath.Join(root, a.Path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("artifact[%d]: file not found: %s", i, a.Path)
		}

		hash := sha256.Sum256(data)
		actual := fmt.Sprintf("%x", hash)

		if actual != a.SHA256 {
			return fmt.Errorf("artifact[%d]: hash mismatch for %s: expected %s, got %s", i, a.Path, a.SHA256, actual)
		}
	}

	return nil
}
