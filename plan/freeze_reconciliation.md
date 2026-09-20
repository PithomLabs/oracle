# ARGUS POC — Freeze Reconciliation Record

## 1. Baseline

- **Plan:** ARGUS POC Execution Plan v1.1
- **Status:** FROZEN — IMPLEMENTATION BASELINE
- **Reference:** `oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md`
- **Phase 0 report:** `oracle/plan/PHASE0_RECONNAISSANCE.md`

## 2. Prior Freeze

- **Prior freeze hash:** `7602699`
- **Superseded by:** Plan 11.1 debt-opaque design
- **Note:** The prior commit is NOT the current baseline. It is recorded for historical reference only.

## 3. Current Freeze

- **Current freeze hash:** `6e23e41311be57004a008a84189c4298d2f0efb0`
- **Repository:** `/home/chaschel/Documents/go/solvent-main`
- **Verification timestamp:** 2026-09-15
- **Working-tree state:** Clean (no modified tracked files)

## 4. Superseded Designs

| Design | Status | Notes |
|--------|--------|-------|
| `FullDebt` kernel constant | Superseded | Absent from `kernel/` as of current freeze |
| Generic Solvent coupling to wizard deployment-review vocabulary | Superseded | Debt is now opaque; domain packs own vocabulary |

## 5. Amended Designs

| Design | Amendment | Evidence |
|--------|-----------|----------|
| `belief.debt` | Now `TEXT[]` with empty default | `db/010_debt_opaque.sql:11` |
| Solvent default debt | Empty `ARRAY[]::TEXT[]` | `db/010_debt_opaque.sql:11` |
| Caller-supplied debt | `EnterBelief` accepts `initialDebt []string` | `kernel/kernel.go:35` |
| Wizard vocabulary | `wizardDebt` in `internal/belief/debt.go` | Separate from kernel |
| EBP/Domain Pack vocabulary | Domain packs own debt vocabulary | Not in generic Solvent |
| Debt interpretation | Solvent does not interpret debt strings | Opacity invariant maintained |

## 6. Deferred Capabilities

| Capability | Status | Notes |
|------------|--------|-------|
| DSSE/SLSA | Deferred | Out of scope for POC |
| Cryptographic attestation | Deferred | Trusted path relies on process-level trust |
| Second real Domain Pack | Deferred | Only conformance test with mock pack |
| Probabilistic inference | Deferred | Not in POC scope |
| Source reputation | Deferred | Not in POC scope |
| W3C PROV | Deferred | Not in POC scope |
| Persistent idempotency infrastructure | Deferred | Process-local only |
| Persistent ArtifactRegistry | Deferred | In-memory only |
| Crash-recoverable projection journal | Deferred | In-memory only |

## 7. Growth Gate Exceptions

| Endpoint | Purpose | Structural Checks | Status |
|----------|---------|-------------------|--------|
| `POST /v1/beliefs/{parent_id}/edges` | Canonical creation of `belief_edge` relationships | parent exists, child exists, parent ≠ child, kind ∈ {derives, contradicts}, uniqueness enforced | **NOT YET IMPLEMENTED** |

**Architectural rule:** Coordinator MUST NOT write Solvent SQL directly. The future Coordinator will call this Solvent endpoint. This is a minimal Growth Gate exception, not a general reopening of the Solvent freeze.

## 8. Repository Boundaries

| Repository | Role | Boundary |
|------------|------|----------|
| Solvent | Canonical epistemic + authority state | Frozen kernel, Growth Gate exceptions only |
| Conductor | Operational workflow | Unchanged, domain-agnostic |
| Oracle | ARGUS backend | Separate Go module |
| Trust UI | Read/projection surface | Separate Go module, calls Coordinator HTTP API |
| reference-loop | Historical/reference implementation | Separate module, untouched |

## 9. Verification Results

| Check | Result | Evidence |
|-------|--------|----------|
| Solvent HEAD | `6e23e41311be57004a008a84189c4298d2f0efb0` | `git rev-parse HEAD` |
| Working tree clean | PASS | `git status --short` empty |
| Tests | 75 passed, 15 failed | CockroachDB unavailable (environment limitation) |
| FullDebt absent | PASS | `grep -r "FullDebt" kernel/` returns empty |
| `010_debt_opaque.sql` exists | PASS | File exists with correct ALTER |
| `wizardDebt` exists | PASS | `internal/belief/debt.go:21` |
| `EnterBelief` accepts debt | PASS | `kernel/kernel.go:35` |
| Edge endpoint absent | PASS | Not in `api/api.go` routes |
| Verdict behavior | PASS | `api/belief.go:194-201` |
| audit_activity | PASS | `api/reads.go:352` |

## 10. Phase 1 Acceptance

**PASS** — All acceptance criteria met:

1. ✅ Actual current Solvent HEAD recorded
2. ✅ Working tree has no modified tracked files
3. ✅ Plan 11.1 state verified
4. ✅ FullDebt absent from kernel
5. ✅ debt-opaque migration present
6. ✅ Freeze reconciliation document complete (this document)
7. ✅ Growth Gate exception explicitly documented
8. ✅ No production code changed
9. ✅ No Solvent endpoint added
10. ✅ No Conductor code changed
11. ✅ reference-loop not changed
