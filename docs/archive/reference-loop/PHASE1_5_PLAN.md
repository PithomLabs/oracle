# Phase 1.5 Validation Battery — Plan

**Status:** Execution plan
**Frozen Design:** Loop Engineering Workflow Design v1.0
**Frozen Solvent:** commit 7602699
**Date:** 2026-09-12

---

## 1. Objective

Prove the minimum safety/binding properties required before Phase 2 (real external-effect validation).

Run exactly five negative scenarios plus three verification checks against the frozen Solvent implementation at commit 7602699.

---

## 2. Frozen Baseline

```text
Workflow Design:  Loop Engineering Workflow Design v1.0
Solvent:          /home/chaschel/Documents/go/solvent-main
Frozen commit:    7602699
```

Before starting and after ending:

```text
cd /home/chaschel/Documents/go/solvent-main
git rev-parse HEAD          -> 7602699
git status --short          -> clean
git diff --name-only        -> empty
git diff --stat             -> empty
```

---

## 3. Test Environment

```text
Conductor:     existing MCP client (SQLite)
Solvent:       frozen at 7602699, in-process server on :18080
Executor:      recording executor (deterministic, no real effect)
Database:      CockroachDB at localhost:26257
Package:       package main (reference-loop)
```

Real external effects MUST NOT occur in Phase 1.5.

---

## 4. Scenarios

### Test A: Operation Mismatch

Authorize operation X, attempt execution of operation Y where X != Y.

Setup:
1. Initialize Conductor DB, Solvent DB, Solvent server
2. Create belief, retire debt, promote belief
3. Create target, attach justification, request authorization, approve
4. Call AuthorizeAction(action="deploy") -> obtain intent_id

Execute:
- Call ExecuteAction(intent_id, action="destroy")

Expected: Execution denied. Executor not invoked. No external effect.

Pass criteria: result.Allowed == false AND executorCallCount == 0.

Evidence: authorized operation identity, requested operation identity, rejection reason, executor invocation count.

### Test B: Duplicate Delivery

Deliver the same execution request twice.

Setup:
1. Full Solvent authorization for "deploy" -> obtain intent_id

Execute:
1. First: ExecuteAction(intent_id, action="deploy")
2. Second: ExecuteAction(intent_id, action="deploy")

Expected: First call succeeds. Second call rejected (intent state != 'live').

Pass criteria: First Allowed=true, second Allowed=false, executorCallCount == 1.

Evidence: both responses, executor invocation count, intent state transition.

### Test C: Stale/Invalid Authorization

Revoke target after authorization, then attempt execution.

Setup:
1. Full Solvent authorization for "deploy" -> obtain intent_id, target_id

Execute:
1. RevokeTarget(target_id, reason="phase-1.5-stale-test")
2. ExecuteAction(intent_id, action="deploy")

Expected: Execution denied (revocation detected in fresh authority re-evaluation).

Pass criteria: result.Allowed == false AND executorCallCount == 0.

Evidence: revocation evidence, denial reason, executor invocation count.

### Test D: Unknown/Unavailable Authorization

No authoritative authorization decision can be established.

Setup:
1. Create a target but do NOT approve it (no activation)
2. Also test with non-existent target_id

Execute:
1. AuthorizeAction with unapproved target
2. AuthorizeAction with non-existent target

Expected: Allowed=false in both cases.

Document:
- Current Solvent returns Allowed=false for both UNKNOWN and DENIED
- Safety: PASS (fail-closed)
- UNKNOWN != DENIED semantic: NOT YET PROVEN

Pass criteria: Safety PASS. Semantic distinction NOT YET PROVEN.

### Test E: Stale Claim x Expired Intent

Agent claims work, valid authorization exists, intent becomes invalid, agent attempts execution.

Setup:
1. Full Solvent authorization for "deploy" -> obtain intent_id
2. Create and claim task in Conductor

Execute:
1. RetractBelief(scenario_id, belief_id) -> cancels live intents via cascade
2. ExecuteAction(intent_id, action="deploy")

Expected: Execution denied (intent state is 'cancelled').

Document:
- Conductor claim state: still active
- Solvent intent state: cancelled
- Boundary: Conductor claim != Solvent authorization

Pass criteria: Execution denied. Coordination/authorization boundary characterized.

### Test F: Declaration Resolution

Verify capability_ref resolves to declaration owner/version/hash.

Verify:
1. contract.json exists with required fields
2. target_snapshot exists with snapshot_hash, justification_set
3. capability_ref -> declaration resolution mechanism: NOT IMPLEMENTED

Document:
- target_snapshot != capability declaration (different semantic objects)
- verdict: NOT YET PROVEN

Pass criteria: Gap documented. NOT YET PROVEN retained.

### Test G: Operation Identity Equality

Same operation-identity definition used at proposal/authorization/execution.

Sub-tests:
1. Determinism: operationID(x,x,x,x) called twice -> equal
2. Different repo -> different ID
3. Different workflow -> different ID
4. Different ref -> different ID
5. Different runID -> different ID
6. Solvent 8-field tuple: same authorization -> same tuple
7. Solvent tuple: different action_name -> mismatch rejection

Expected: All effect-relevant changes produce different operations. Comparison is deterministic.

Pass criteria: All sub-tests pass.

### Test H: Solvent Freeze Integrity

At beginning and end of battery:

```text
cd /home/chaschel/Documents/go/solvent-main
git rev-parse HEAD    -> 7602699
git status --short    -> empty
git diff --name-only  -> empty
```

Pass criteria: HEAD = 7602699, working tree clean, no source modifications.

---

## 5. Pass/Fail Criteria

| Verdict | Meaning |
|---------|---------|
| PASS | Expected boundary behavior observed through real participant boundary |
| NOT_YET_PROVEN | Current controlled implementation cannot establish required behavior without new architecture |
| BLOCKED | Test cannot run because environment/frozen participant cannot support contract |

Do not convert NOT_YET_PROVEN into PASS.

---

## 6. Failure Classification

For any unexpected result:

```text
IMPLEMENTATION_DEFECT
INTEGRATION_DEFECT
EXECUTOR_DEFECT
DEPLOYMENT_ENVIRONMENT_DEFECT
SPECIFICATION_DEFECT
NEW_SECURITY_PROPERTY
```

Do not pre-classify. Do not modify architecture to make tests pass.

---

## 7. Evidence Requirements

For every scenario record:

```text
scenario_id
run_id
task_id
plan_id / plan_version (where applicable)
operation_id
actor
capability_ref
declaration version/hash
setup
request
expected result
actual result
authorization state
executor invocation
effect classification
authoritative evidence source
verdict
```

---

## 8. Out of Scope

Do NOT implement:

- claim TTL / heartbeat
- scheduler
- proof-token subsystem
- plan-scope enforcement engine
- Conductor policy engine
- cryptographic Conductor<->Solvent ledger
- workflow templates / molecules / DAG runtime
- knowledge graph
- new event bus
- Agent planning engine
- new Solvent capabilities
- real external effect
- BM-IST integration

If one appears necessary -> STOP, report Growth Gate candidate.

---

## 9. Deliverables

```text
PHASE1_5_PLAN.md          <- this file
phase1_5_battery_test.go  <- automated test battery
PHASE1_5_RESULTS.md       <- after test run
PHASE1_5_DISPOSITION.md   <- after analysis
```
