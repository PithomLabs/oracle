# Phase 1.5 Validation Battery — Results

**Date:** 2026-09-12
**Frozen Solvent:** commit 7602699
**Test Package:** package main (reference-loop)
**Test Duration:** 3.78s (all 7 scenarios)

---

## Summary

| Scenario | Verdict |
|----------|---------|
| A. Operation Mismatch | PASS |
| B. Duplicate Delivery | PASS |
| C. Stale/Invalid Authorization | PASS |
| D. Unknown/Unavailable Authorization | SAFETY=PASS SEMANTICS=NOT_YET_PROVEN |
| E. Stale Claim x Expired Intent | PASS |
| F. Declaration Resolution | NOT_YET_PROVEN |
| G. Operation Identity Equality | PASS |
| H. Solvent Freeze Integrity | YES (beginning + end) |

---

## Test A — Operation Mismatch

**scenario_id:** e1f0371b-188a-44a0-95c5-69911f3aa81b
**run_id:** phase15-20260912093702

**Setup:**
- Created belief, retired debt, promoted belief
- Created target, attached justification, requested authorization, approved
- Authorized action="deploy" -> obtained intent_id=23874240-09f5-48d9-a2ea-f94f47319edc
- Target: 9fa7b6ac-0c80-4652-ba6b-5fdc8c21019d

**Request:**
- ExecuteAction with action="destroy" (mismatched with authorized "deploy")

**Result:**
- Allowed=false
- Reason="action_name mismatch"
- Executor invocation count: 0

**Authorization State:**
- Solvent 8-field tuple comparison rejected action name mismatch
- ClaimIntent CAS never reached
- Fresh authority re-evaluation in PrepareForAction detected mismatch

**Evidence Source:** Solvent kernel.authorizeWithinTx field-by-field comparison

**Verdict:** PASS

---

## Test B — Duplicate Delivery

**scenario_id:** 9e578eee-078d-4737-8495-e24674c8a21c
**run_id:** phase15-20260912093703

**Setup:**
- Full Solvent authorization for "deploy" -> intent_id=26ae18a8-841c-451c-8fa5-1df726b60f4d

**Request:**
1. First: ExecuteAction(intent_id, action="deploy")
2. Second: ExecuteAction(intent_id, action="deploy")

**Result:**
- First: Allowed=true, Success=true, IntentState="executed"
- Second: Allowed=false, Error="claim intent: intent is not in live state"
- Executor invocation count: 1

**Authorization State:**
- First call: intent transitioned live -> executing -> executed
- Second call: ClaimIntent CAS failed (WHERE state='live' predicate not met)
- Intent state is 'executed' after first call

**Evidence Source:** Solvent kernel.ClaimIntent atomic CAS

**Verdict:** PASS — duplicate delivery rejected; single-use authorization enforced

---

## Test C — Stale/Invalid Authorization

**scenario_id:** 7b4d6187-dcf3-403d-b3bd-ddb5fa826f19
**run_id:** phase15-20260912093703

**Setup:**
- Full Solvent authorization for "deploy" -> intent_id=2bfecedc-3ff0-4467-a4da-fd8cc66c3f82
- Target: 76e944f8-833c-4c57-82f5-14d68e03e243

**Request:**
1. RevokeTarget(target_id, reason="phase-1.5-stale-test")
2. ExecuteAction(intent_id, action="deploy")

**Result:**
- Allowed=false
- Reason="no activation or revocation exists"
- Executor invocation count: 0

**Authorization State:**
- Target revocation recorded
- PrepareForAction re-reads current state
- kernel.Authorize detects revocation via target_revocation table
- Fresh evaluation returns Allowed=false

**Evidence Source:** Solvent kernel.authorizeWithinTx + target_revocation check

**Verdict:** PASS — stale authorization rejected by fresh authority re-evaluation

---

## Test D — Unknown/Unavailable Authorization

**scenario_id:** c29cddff-858e-4529-a7c3-3992afeb4334
**run_id:** phase15-20260912093704

**Sub-test 1: Unapproved target (no activation)**
- Created target but did NOT approve (no activation)
- AuthorizeAction -> Allowed=false, Reason="no activation or revocation exists"

**Sub-test 2: Non-existent target_id**
- Used random UUID as target_id
- AuthorizeAction -> HTTP 404, "target not found"

**Result:**
- Both sub-tests: authorization denied
- Executor invocation count: 0

**Safety Result:** PASS (fail-closed)

**Semantic Documentation:**
- Current Solvent returns Allowed=false for both UNKNOWN and DENIED
- The interface does not distinguish UNKNOWN/UNAVAILABLE from DENIED
- UNKNOWN != DENIED semantic is NOT YET PROVEN against the frozen Solvent interface
- Classification: conformance/specification gap, not implementation defect

**Evidence Source:** Solvent REST API + kernel.Authorize

**Verdict:** SAFETY=PASS SEMANTICS=NOT_YET_PROVEN

---

## Test E — Stale Claim x Expired Intent

**scenario_id:** f5f1d3e7-a168-447c-81c6-5099c0484385
**run_id:** phase15-20260912093705

**Setup:**
- Full Solvent authorization for "deploy" -> intent_id=9a4b77b9-3063-4308-9648-75af9c8c7bfa
- Conductor task created: 0c6db83a-15d7-436a-b29c-ade49e6d00c8
- Task claimed: status=active

**Request:**
1. RetractBelief(scenario_id, belief_id) -> cancels live intents via cascade
2. ExecuteAction(intent_id, action="deploy")

**Result:**
- Allowed=false
- Reason="belief 4caf1a69-5db4-4961-a26d-2a78cccde89d not promoted"
- Executor invocation count: 0

**Coordination Claim State:**
- Conductor task status: active (still claimed)
- Conductor claim is coordination state, not authorization

**Solvent Authorization State:**
- Belief retracted -> RetractCascade cancelled live intents
- Intent state: cancelled (cascaded from belief retraction)
- ExecuteAction denied because belief is no longer promoted

**Boundary Characterization:**
- Conductor claim != Solvent authorization
- Coordination claim does not establish authorization validity
- Belief retraction is the primary mechanism for invalidating stale authority

**Evidence Source:** Solvent kernel.RetrapCascade + FK cascade + authority re-evaluation

**Verdict:** PASS — stale claim x expired intent boundary characterized

---

## Test F — Declaration Resolution

**scenario_id:** d2da54ad-e67a-44bf-9f24-09f5e8df0903
**run_id:** phase15-20260912093705

**Verification:**
- contract.json exists with all 10 required fields (verified present)
- target_snapshot has snapshot_hash, justification_set, integrity-pinned FK constraints

**Gap Identified:**
- Workflow Design v1.0 Section 8 requires capability_ref resolution to:
  - declaration owner
  - declaration version
  - content hash
  - effective reference
  - retrievable declaration content
- The frozen Solvent uses target_snapshot (authority construct)
- target_snapshot != capability declaration (different semantic objects)
- capability_ref -> declaration resolution mechanism: NOT IMPLEMENTED
- The reference-loop does not resolve capability_ref to declaration owner/version/hash

**Classification:** integration/design coverage gap

**Verdict:** NOT_YET_PROVEN

---

## Test G — Operation Identity Equality

**scenario_id:** operation-identity
**run_id:** 20260912093706

**Sub-tests:**
| Sub-test | Description | Result |
|----------|-------------|--------|
| a | Determinism: same inputs -> same output | PASS |
| b | Different repo -> different ID | PASS |
| c | Different workflow -> different ID | PASS |
| d | Different ref -> different ID | PASS |
| e | Different runID -> different ID | PASS |
| f | Canonicalization rule: deploy:repo:workflow:ref:run_id | PASS |
| g | Colon-separated 5-part structure | PASS |

**Canonicalization Rule:**
- Format: `deploy:repo:workflow:ref:run_id`
- All 5 components are effect-relevant
- No declared-immaterial parameters in current implementation
- Comparison is deterministic (string equality)

**Evidence Source:** operationID function + Solvent 8-field tuple comparison

**Verdict:** PASS

---

## Test H — Solvent Freeze Integrity

**Beginning Check:**
```text
HEAD: 7602699947e5fe5bafbeb0709d0a16b17d595318
Status: (clean)
Diff names: (empty)
Diff stat: (empty)
```

**End Check:**
- HEAD = 7602699
- Working tree clean
- No source modifications

**Verdict:** YES (Solvent frozen)

---

## Evidence Classification

| Classification | Count | Scenarios |
|---------------|-------|-----------|
| PASS | 5 | A, B, C, E, G |
| NOT_YET_PROVEN | 2 | D (semantics), F |
| BLOCKED | 0 | — |

---

## Executor Invocation Summary

| Scenario | Expected | Actual | Pass |
|----------|----------|--------|------|
| A (Mismatch) | 0 | 0 | YES |
| B (Duplicate) | 1 | 1 | YES |
| C (Stale) | 0 | 0 | YES |
| D (Unknown) | 0 | 0 | YES |
| E (Stale Claim) | 0 | 0 | YES |

---

## Effect Classification

All scenarios produce simulated/non-external effects only.
No real external effects occurred during Phase 1.5.
