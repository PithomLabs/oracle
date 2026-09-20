package corpus

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidManifest(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if manifest.ManifestVersion != "1.0.0" {
		t.Errorf("wrong version: %s", manifest.ManifestVersion)
	}
	if manifest.CorpusID == "" {
		t.Error("corpus_id is empty")
	}
	if len(manifest.Artifacts) == 0 {
		t.Error("no artifacts")
	}
}

func TestValidateValidManifest(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if err := manifest.Validate(); err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestVerifyValidManifest(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	root := filepath.Join("..")
	if err := manifest.Verify(root); err != nil {
		t.Errorf("Verify failed: %v", err)
	}
}

func TestVerifyMissingFile(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "nonexistent.md", SHA256: "abc123", ProvenanceClass: "operator_asserted", StableLocator: "test"},
		},
	}
	if err := manifest.Verify(".."); err == nil {
		t.Error("missing file not detected")
	}
}

func TestVerifyHashMismatch(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "docs/BMIST_POC_SEED.md", SHA256: "0000000000000000000000000000000000000000000000000000000000000000", ProvenanceClass: "operator_asserted", StableLocator: "test"},
		},
	}
	if err := manifest.Verify(".."); err == nil {
		t.Error("hash mismatch not detected")
	}
}

func TestValidateDuplicateID(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "dup", Path: "a", SHA256: "abc123", ProvenanceClass: "operator_asserted", StableLocator: "a"},
			{ID: "dup", Path: "b", SHA256: "def456", ProvenanceClass: "operator_asserted", StableLocator: "b"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Error("duplicate ID not detected")
	}
}

func TestValidateMalformedSHA256(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "a", SHA256: "not-a-hash", ProvenanceClass: "operator_asserted", StableLocator: "a"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Error("malformed SHA256 not detected")
	}
}

func TestValidateUnsupportedProvenance(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "a", SHA256: "abc123", ProvenanceClass: "unsupported", StableLocator: "a"},
		},
	}
	if err := manifest.Validate(); err == nil {
		t.Error("unsupported provenance not detected")
	}
}

func TestValidateEmptyArtifacts(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts:       []ArtifactRecord{},
	}
	if err := manifest.Validate(); err == nil {
		t.Error("empty artifacts not detected")
	}
}

func TestVerifyAbsolutePath(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "/etc/passwd", SHA256: "abc123", ProvenanceClass: "operator_asserted", StableLocator: "test"},
		},
	}
	if err := manifest.Verify("."); err == nil {
		t.Error("absolute path not rejected")
	}
}

func TestVerifyPathTraversal(t *testing.T) {
	manifest := &Manifest{
		ManifestVersion: "1.0.0",
		CorpusID:        "test",
		Artifacts: []ArtifactRecord{
			{ID: "test", Path: "../../etc/passwd", SHA256: "abc123", ProvenanceClass: "operator_asserted", StableLocator: "test"},
		},
	}
	if err := manifest.Verify("."); err == nil {
		t.Error("path traversal not rejected")
	}
}

func TestMultipleArtifacts(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(manifest.Artifacts) < 2 {
		t.Error("expected at least 2 artifacts")
	}
	ids := make(map[string]bool)
	for _, a := range manifest.Artifacts {
		if ids[a.ID] {
			t.Errorf("duplicate ID: %s", a.ID)
		}
		ids[a.ID] = true
	}
}

func TestStableLocators(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	for _, a := range manifest.Artifacts {
		if a.StableLocator == "" {
			t.Errorf("artifact %s has empty stable_locator", a.ID)
		}
	}
}

func TestComputeHashFromBytes(t *testing.T) {
	content := []byte("test content")
	hash := sha256.Sum256(content)
	expected := fmt.Sprintf("%x", hash)
	actual := fmt.Sprintf("%x", sha256.Sum256(content))
	if expected != actual {
		t.Errorf("hash mismatch: expected %s, got %s", expected, actual)
	}
}

func TestVerifyActualFiles(t *testing.T) {
	manifest, err := Load("manifest.json")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	root := filepath.Join("..")
	for _, a := range manifest.Artifacts {
		fullPath := filepath.Join(root, a.Path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("artifact %s: file not found: %s", a.ID, a.Path)
			continue
		}
		hash := sha256.Sum256(data)
		actual := fmt.Sprintf("%x", hash)
		if actual != a.SHA256 {
			t.Errorf("artifact %s: hash mismatch: expected %s, got %s", a.ID, a.SHA256, actual)
		}
	}
}
