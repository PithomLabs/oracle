# PHASE 2 — DOMAIN PACK SPECIFICATION REPORT

**Timestamp:** 2026-09-15
**Status:** PASS

## Files Created

| File | Purpose |
|------|---------|
| `oracle/go.mod` | Go module definition (zero dependencies) |
| `oracle/domain-pack/README.md` | Domain Pack documentation |
| `oracle/domain-pack/registry.go` | In-memory pack registry |
| `oracle/domain-pack/registry_test.go` | Registry tests |
| `oracle/domain-pack/bmist/v1/pack.json` | BM-IST Domain Pack descriptor |
| `oracle/domain-pack/bmist/v1/types.go` | Go types matching JSON schema |
| `oracle/domain-pack/bmist/v1/validate.go` | Structural validator (15 rules) |
| `oracle/domain-pack/bmist/v1/validate_test.go` | Validator tests |
| `oracle/domain-pack/testdata/bmist/v1/pack.json` | Test fixture for LoadFromDisk |
| `oracle/testdata/packs/valid/valid_bmist_v1.json` | Valid test fixture |
| `oracle/testdata/packs/invalid/missing_pack_id.json` | Invalid: missing pack_id |
| `oracle/testdata/packs/invalid/missing_version.json` | Invalid: missing version |
| `oracle/testdata/packs/invalid/invalid_version.json` | Invalid: v prefix |
| `oracle/testdata/packs/invalid/empty_claim_types.json` | Invalid: empty claim_types |
| `oracle/testdata/packs/invalid/missing_debt_vocabulary.json` | Invalid: empty debt_vocabulary |
| `oracle/testdata/packs/invalid/invalid_initial_debt.json` | Invalid: undeclared debt item |
| `oracle/testdata/packs/invalid/invalid_retirement_rule.json` | Invalid: unknown evidence class |
| `oracle/testdata/packs/invalid/missing_human_gate.json` | Invalid: missing faithfulness_review |
| `oracle/testdata/packs/invalid/duplicate_debt.json` | Invalid: duplicate debt vocabulary |

## Pack Structure

### BM-IST Pack (`bmist/v1/pack.json`)

```json
{
  "pack_id": "bmist",
  "version": "1.0.0",
  "name": "BM-IST Fisher-Rigidity Domain Pack",
  "description": "Domain Pack for the BM-IST Fisher-rigidity claim verification POC",
  "claim_types": ["derived", "accommodated", "postulated"],
  "evidence_classes": ["reproducible_artifact", "operator_asserted"],
  "debt_vocabulary": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
  "initial_debt": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
  "retirement_rules": {
    "needMap": {"evidence_class": "reproducible_artifact", "rule": "map_check"},
    "needInvariant": {"evidence_class": "reproducible_artifact", "rule": "invariant_check"},
    "needToyCheck": {"evidence_class": "reproducible_artifact", "rule": "toy_model_check"},
    "needNullModel": {"evidence_class": "operator_asserted", "rule": "scope_clarification"},
    "needObstruction": {"evidence_class": "reproducible_artifact", "rule": "obstruction_construction"},
    "needFaithfulnessReview": {"evidence_class": "operator_asserted", "rule": "faithfulness_review"}
  },
  "falsifiers": ["counterexample", "contradiction", "missing_evidence", "alternative_explanation"],
  "consequential_actions": [
    {"action": "publish_claim", "requires": "promoted", "gates": ["faithfulness_review"]}
  ],
  "human_gated_transitions": ["faithfulness_review", "scope_clarification", "obstruction_assessment"]
}
```

## Validator Rules

| # | Rule | Test |
|---|------|------|
| 1 | pack_id non-empty | `TestMissingPackID` |
| 2 | version valid SemVer 2.0.0 (no v prefix) | `TestInvalidVersion`, `TestInvalidVersionShort` |
| 3 | claim_types non-empty | `TestEmptyClaimTypes` |
| 4 | claim_types ⊆ {derived, accommodated, postulated} | `TestUnsupportedClaimType` |
| 5 | evidence_classes non-empty | `TestMissingEvidenceClasses` |
| 6 | debt_vocabulary non-empty | `TestMissingDebtVocabulary` |
| 7 | initial_debt ⊆ debt_vocabulary | `TestInvalidInitialDebt` |
| 8 | retirement_rules keys ⊆ debt_vocabulary | `TestInvalidRetirementRuleKey` |
| 9 | retirement_rules evidence_class ∈ declared evidence_classes | `TestInvalidRetirementRuleEvidenceClass` |
| 10 | falsifiers non-empty | `TestMissingFaithfulnessReview` (indirect) |
| 11 | consequential_actions structurally valid | Passes with valid pack |
| 12 | human_gated_transitions ⊇ {faithfulness_review, scope_clarification, obstruction_assessment} | `TestMissingFaithfulnessReview`, `TestMissingScopeClarification`, `TestMissingObstructionAssessment` |
| 13 | no duplicate values in arrays | `TestDuplicateDebtVocabulary`, `TestDuplicateClaimType` |
| 14 | no empty strings in arrays or required identifiers | `TestEmptyVocabularyString` |
| 15 | pack version and identity internally consistent | Passes with valid pack |

## Registry Behavior

| Operation | Behavior |
|-----------|----------|
| `Register(pack)` | Validate first, reject invalid, reject duplicate. Returns error, never panics. |
| `Get(packID, version)` | Deterministic lookup. Returns typed error for not-found. |
| `LoadFromDisk(root)` | Scan `<root>/<pack-id>/<version>/pack.json`. Fail fast on first invalid pack. |

## Test Results

```
go test ./...     → 32 passed in 2 packages
go test -race ./... → 32 passed in 2 packages
go vet ./...      → clean
```

### Test Breakdown

- `domain-pack/bmist/v1`: 24 tests passed (validation + fixture tests)
- `domain-pack`: 8 tests passed (registry + concurrency tests)

## Neutrality Verification

- `registry.go` contains ZERO BM-IST/physics vocabulary
- `types.go` contains ZERO BM-IST/physics vocabulary
- All BM-IST specifics live under `bmist/v1/`
- Generic registry loads `Pack` interface, not concrete types

**PASS** — Generic registry code contains no BM-IST/physics vocabulary.

## Repositories Changed

| Repository | Changes |
|------------|---------|
| Oracle | `go.mod`, `domain-pack/` (new), `testdata/` (new) |

## Repositories Explicitly Unchanged

| Repository | Status |
|------------|--------|
| Solvent | Unchanged |
| Conductor | Unchanged |
| reference-loop | Unchanged |

## Acceptance Results

| # | Criterion | Result |
|---|-----------|--------|
| 1 | `oracle/go.mod` exists and builds | ✅ PASS |
| 2 | BM-IST pack.json exists and validates | ✅ PASS |
| 3 | BM-IST vocabulary exactly matches frozen plan | ✅ PASS |
| 4 | Structural human-gate floor is enforced | ✅ PASS |
| 5 | Retirement rules reference only declared evidence classes | ✅ PASS |
| 6 | Invalid pack fixtures are rejected | ✅ PASS |
| 7 | PackRegistry registers and retrieves valid packs | ✅ PASS |
| 8 | PackRegistry rejects invalid/duplicate packs | ✅ PASS |
| 9 | LoadFromDisk works | ✅ PASS |
| 10 | Registry concurrency tests pass | ✅ PASS |
| 11 | `go test ./...` passes | ✅ PASS |
| 12 | `go test -race ./...` passes | ✅ PASS |
| 13 | No Solvent changes | ✅ PASS |
| 14 | No Conductor changes | ✅ PASS |
| 15 | No reference-loop changes | ✅ PASS |
| 16 | No Phase 3+ implementation exists | ✅ PASS |
| 17 | Generic registry code contains no BM-IST/physics vocabulary | ✅ PASS |

## Deviations

None. All acceptance criteria met.

## Next Phase

PHASE 3 — EBP RESEARCH PACKET CONTRACT
