# Phase 1.5 Validation Battery — Disposition

**Date:** 2026-09-12
**Frozen Solvent:** commit 7602699
**Frozen Design:** Loop Engineering Workflow Design v1.0

---

## Final Gate Summary

```text
Phase 1.5 Verdict:
    PASS WITH LIMITATIONS

Operation mismatch:
    PASS

Duplicate delivery:
    PASS

Stale/invalid authorization:
    PASS

UNKNOWN/unavailable:
    PASS (safety) / NOT YET PROVEN (semantics)

Stale claim x expired intent:
    PASS

Declaration resolution:
    NOT YET PROVEN

Operation identity equality:
    PASS

Solvent frozen:
    YES

New architecture required:
    NO

Growth Gate candidate:
    NO
```

---

## Findings

### Finding 1: Operation Mismatch Rejected (PASS)

**Finding:** Solvent's 8-field tuple comparison in kernel.authorizeWithinTx rejects execution when the requested action name does not match the authorized action name.

**Classification:** Implementation correct — no defect.

**Disposition:** PASS. No action required.

**Evidence:** Test A — action_name mismatch detected, Allowed=false, executor not invoked.

### Finding 2: Duplicate Delivery Rejected (PASS)

**Finding:** Solvent's ClaimIntent atomic CAS (WHERE state='live') prevents duplicate execution. After first execution, intent state transitions to 'executed', and subsequent attempts fail with "intent is not in live state".

**Classification:** Implementation correct — single-use authorization enforced.

**Disposition:** PASS. No action required.

**Evidence:** Test B — first call Allowed=true, second call Allowed=false, executor invoked exactly once.

### Finding 3: Stale Authorization Rejected (PASS)

**Finding:** Solvent's fresh authority re-evaluation in PrepareForAction detects target revocation. After revocation, execution is denied with "no activation or revocation exists".

**Classification:** Implementation correct — no cached authority.

**Disposition:** PASS. No action required.

**Evidence:** Test C — revocation detected, Allowed=false, executor not invoked.

### Finding 4: UNKNOWN/UNAVAILABLE Treated as DENIED (SAFETY PASS, SEMANTICS NOT YET PROVEN)

**Finding:** Solvent returns Allowed=false when no authoritative authorization decision can be available (unapproved target, non-existent target). The interface does not distinguish UNKNOWN/UNAVAILABLE from DENIED.

**Classification:** Conformance/specification gap. The workflow design requires UNKNOWN != DENIED semantic. The frozen Solvent implements fail-closed behavior (both result in Allowed=false).

**Disposition:** Safety PASS. Semantic distinction NOT YET PROVEN. Do not modify Solvent to create the distinction. Record as known gap for Phase 2 consideration.

**Evidence:** Test D — both unapproved and non-existent targets return Allowed=false.

### Finding 5: Stale Claim x Expired Intent (PASS)

**Finding:** Belief retraction cascades to cancel live intents via FK cascade. Conductor task claim remains active (coordination state). The boundary between Conductor claim and Solvent authorization is correctly characterized.

**Classification:** Implementation correct — coordination claim does not establish authorization validity.

**Disposition:** PASS. No action required.

**Evidence:** Test E — belief retracted, intent cancelled, Conductor task still active, execution denied.

### Finding 6: Declaration Resolution NOT YET PROVEN

**Finding:** The Workflow Design v1.0 Section 8 requires capability_ref resolution to declaration owner/version/content-hash/effective-reference/retrievable-content. The frozen Solvent uses target_snapshot (authority construct), which is a different semantic object. The capability_ref -> declaration resolution mechanism is not implemented.

**Classification:** Integration/design coverage gap.

**Disposition:** NOT YET PROVEN. Do NOT reinterpret target_snapshot as a capability declaration. Do NOT build a capability registry in this phase. Record as required work before Phase 2.

**Evidence:** Test F — contract.json present, target_snapshot present, capability_ref resolution mechanism absent.

### Finding 7: Operation Identity Equality (PASS)

**Finding:** The operation identity format is `deploy:repo:workflow:ref:run_id`. All 5 components are effect-relevant. Comparison is deterministic (string equality). No declared-immaterial parameters exist in the current implementation.

**Classification:** Implementation correct — deterministic canonicalization.

**Disposition:** PASS. No action required.

**Evidence:** Test G — all 7 sub-tests passed.

---

## Architecture Change Assessment

**Is architecture change justified?** NO.

The frozen design holds for the tested scenarios. The two NOT_YET_PROVEN findings are:
1. UNKNOWN != DENIED semantic (conformance gap, not safety defect)
2. capability_ref resolution (integration gap, not architecture defect)

Neither requires architectural modification in Phase 1.5.

---

## Growth Gate Assessment

**Is a Growth Gate candidate triggered?** NO.

The NOT_YET_PROVEN findings are expected outcomes of the Phase 1.5 battery. They identify work required before Phase 2, not architectural deficiencies that require design changes.

The Phase 2 pre-conditions are:
- pinned declaration
- pinned component versions
- pre-registered pass criteria
- pre-registered failure/falsification criteria
- verified enforcement path

The NOT_YET_PROVEN findings inform the Phase 2 preparation but do not trigger a Growth Gate.

---

## Phase 2 Readiness

**What must be done before Phase 2:**

1. Implement capability_ref -> declaration resolution (Test F)
2. Document UNKNOWN != DENIED semantic requirement (Test D)
3. Pin declaration version and content hash
4. Pre-register pass/fail/falsification criteria
5. Verify enforcement path with pinned declaration

**What is NOT required:**
- Architecture changes
- New Solvent capabilities
- Claim TTL/heartbeat
- Scheduler
- Proof-token subsystem
- Cryptographic ledger

---

## Disposition Log

| Finding | Classification | Disposition | Evidence |
|---------|---------------|-------------|----------|
| Operation mismatch rejected | Implementation correct | PASS | Test A |
| Duplicate delivery rejected | Implementation correct | PASS | Test B |
| Stale authorization rejected | Implementation correct | PASS | Test C |
| UNKNOWN treated as DENIED | Conformance gap | SAFETY PASS / SEMANTICS NOT YET PROVEN | Test D |
| Stale claim boundary | Implementation correct | PASS | Test E |
| Declaration resolution absent | Integration gap | NOT YET PROVEN | Test F |
| Operation identity equality | Implementation correct | PASS | Test G |
| Solvent frozen | Integrity verified | YES | Test H |

---

## Evidence Classification

| Category | Count |
|----------|-------|
| PASS | 5 |
| NOT_YET_PROVEN | 2 |
| BLOCKED | 0 |
| IMPLEMENTATION_DEFECT | 0 |
| INTEGRATION_DEFECT | 0 |
| EXECUTOR_DEFECT | 0 |
| SPECIFICATION_DEFECT | 0 |
| NEW_SECURITY_PROPERTY | 0 |

---

## Recommendation

**Phase 1.5 evidence is sufficient to proceed to Phase 2 preparation.**

The frozen design holds. The safety/binding properties are proven for the tested scenarios. The NOT_YET_PROVEN findings are known gaps that must be addressed before Phase 2 real external-effect validation, but they do not indicate architectural deficiencies.

The next human decision should be whether the NOT_YET_PROVEN findings are acceptable for Phase 2 entry, or whether additional work is required first.
