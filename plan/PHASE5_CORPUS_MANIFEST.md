# PHASE 5 CORPUS MANIFEST — COMPLETE

## Status: PASS

## Corpus Identity

- `corpus_id`: bmist-poc-seed-v1
- `manifest_version`: 1.0.0

## Artifacts Included

| ID | Path | Provenance | SHA-256 |
|----|------|------------|---------|
| bmist-seed-001 | docs/BMIST_POC_SEED.md | operator_asserted | a684cf93b043d0542eba5a59b059a18f15c4ec30ebbf9a39f0c5f03b33ffe4f7 |
| qwen-payoff-001 | docs/payoff_qwen.md | agent_derived | 1da86096ea0995bc9c3c848d31e062aabb7bb61fb6e729b4a7afbf8f4bce4dc0 |

## Provenance Classes

- `operator_asserted`: BM-IST seed document — operator-provided research input
- `agent_derived`: Qwen artifact — contextual research material, NOT epistemic authority

## Hash Verification

All hashes recomputed from actual file bytes before declaring PASS.

## Path Safety

- Repository-relative paths only
- Absolute paths rejected
- Path traversal (`../`) rejected

## Test Results

```
98 passed in 7 packages
go test -race ./... passed
go vet ./... passed
```

Tests include:
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
- actual file verification

## Files Created

- `oracle/corpus/manifest.json`
- `oracle/corpus/manifest.go`
- `oracle/corpus/manifest_test.go`

## Repositories Modified

- `oracle/corpus/` — new package

## Repositories Unchanged

- `solvent-main` — NO changes
- `conductor` — NO changes
- `oracle/reference-loop` — NO changes
- `oracle/domain-pack/` — NO changes
- `oracle/packet/` — NO changes
- `oracle/verifier/` — NO changes

## Limitations

1. Manifest is in-memory, non-persistent
2. Hashes are computed at verification time, not cached
3. No automatic manifest regeneration

## Deviations

None

## Key Distinctions

- **Corpus Manifest** = what research inputs/artifacts comprise the corpus
- **ArtifactRegistry** = trusted verifier outputs admitted during the running POC
- Qwen artifact is `agent_derived`/contextual unless otherwise admitted
- Corpus membership does NOT confer epistemic authority
- Verifier artifacts remain governed by ArtifactRegistry

## Acceptance

- [x] manifest.json exists
- [x] BM-IST seed is included
- [x] Qwen-derived artifact is included
- [x] Every artifact has a correct SHA-256
- [x] Provenance classes are correct
- [x] Repository-relative paths are used
- [x] Path traversal is rejected
- [x] Missing files are detected
- [x] Hash mismatches are detected
- [x] Duplicate IDs are rejected
- [x] go test ./... passes
- [x] go test -race ./... passes
- [x] go vet ./... passes
- [x] Generic corpus code contains no BM-IST-specific logic
- [x] Solvent unchanged
- [x] Conductor unchanged
- [x] reference-loop unchanged
- [x] Domain Pack unchanged
- [x] Packet unchanged
- [x] Verifier unchanged
- [x] No Phase 6+ implementation exists
