# PHASE 5 PLAN — CORPUS MANIFEST

## Files to Create (3 total)

```
oracle/corpus/
├── manifest.json
├── manifest.go
└── manifest_test.go
```

## Step 1: `manifest.json` — Corpus Manifest

```json
{
  "manifest_version": "1.0.0",
  "corpus_id": "bmist-poc-seed-v1",
  "artifacts": [
    {
      "id": "bmist-seed-001",
      "path": "docs/BMIST_POC_SEED.md",
      "sha256": "a684cf93b043d0542eba5a59b059a18f15c4ec30ebbf9a39f0c5f03b33ffe4f7",
      "provenance_class": "operator_asserted",
      "description": "BM-IST POC seed document",
      "stable_locator": "oracle/docs/BMIST_POC_SEED.md"
    },
    {
      "id": "qwen-payoff-001",
      "path": "docs/payoff_qwen.md",
      "sha256": "1da86096ea0995bc9c3c848d31e062aabb7bb61fb6e729b4a7afbf8f4bce4dc0",
      "provenance_class": "agent_derived",
      "description": "Qwen-derived payoff analysis artifact",
      "stable_locator": "oracle/docs/payoff_qwen.md"
    }
  ]
}
```

## Step 2: `manifest.go` — Go Loader + Validator

```go
package corpus

type Manifest struct {
    ManifestVersion string           `json:"manifest_version"`
    CorpusID        string           `json:"corpus_id"`
    Artifacts       []ArtifactRecord `json:"artifacts"`
}

type ArtifactRecord struct {
    ID               string `json:"id"`
    Path             string `json:"path"`
    SHA256           string `json:"sha256"`
    ProvenanceClass  string `json:"provenance_class"`
    Description      string `json:"description"`
    StableLocator    string `json:"stable_locator"`
}

func Load(path string) (*Manifest, error)
func (m *Manifest) Validate() error
func (m *Manifest) Verify(root string) error
```

Key functions:
- `Load()` — reads and parses manifest.json
- `Validate()` — structural validation (unique IDs, valid SHA-256, supported provenance)
- `Verify()` — recomputes hashes, checks file existence, path safety

## Step 3: `manifest_test.go` — ~12 tests

- valid manifest loads
- valid hash verified
- missing file detected
- hash mismatch detected
- duplicate artifact ID rejected
- malformed SHA-256 rejected
- unsupported provenance class rejected
- empty artifact list rejected
- absolute path rejected
- path traversal rejected
- multiple artifacts
- stable repository-relative locators

## Verification

```bash
cd /home/chaschel/Documents/go/oracle
go test ./...
go test -race ./...
go vet ./...
```
