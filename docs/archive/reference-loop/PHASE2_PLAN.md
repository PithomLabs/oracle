# Phase 2 Plan — First Real External-Effect Experiment

**Date:** 2026-09-12
**Status:** Pre-registered plan (planning only, no implementation)
**Frozen Design:** Loop Engineering Workflow Design v1.0
**Frozen Solvent:** commit 7602699947e5fe5bafbeb0709d0a16b17d595318

---

## A. First Real Consequential Operation

**Operation:** GitHub Actions workflow dispatch trigger on `pithomlabs/reference-loop-test`

**Why this operation:**
- Already proven through Phase 1.5 happy path and recording executor
- Exists in the frozen Solvent executor registry (`github_trigger_workflow`)
- Uses the Solvent `actionExecutorMap` with `"deploy" -> "github_trigger_workflow"`
- The test repository and workflow are designed for this exact purpose

**Executor:** Existing `github.NewHTTPProvider` adapter. No new executor.

---

## B. External Source of Record (SOR)

**Authoritative system:** GitHub Actions API

**Specific endpoint:**
```
GET /repos/{owner}/{repo}/actions/runs/{github_run_id}
```

**Authoritative fields:**
- `status`: "completed" / "in_progress" / "queued"
- `conclusion`: "success" / "failure" / "cancelled" / null
- `id`: github_run_id (numeric workflow run ID assigned by GitHub)

**Effect confirmation criteria:**
- `status == "completed" && conclusion == "success"` -> EFFECT_CONFIRMED
- `status == "completed" && conclusion != "success"` -> execution failed at external layer
- `status == "in_progress"` -> AMBIGUOUS (wait/retry)
- HTTP error or empty response -> AMBIGUOUS (retry)

**The SOR is never the Executor.** The Executor invokes; GitHub observes.

**Correlation:** `execution_run_id` <-> `github_run_id` is evidence correlation, not operation identity. GitHub does not know the operation identity when the dispatch is issued.

---

## C. Capability Declaration

File: `declaration_github_deploy_v1.json`

```json
{
  "capability_ref": "github:pithomlabs/reference-loop-test:deploy",
  "declaration_owner": "reference-loop",
  "declaration_version": "v1.0.0",
  "content_hash": "sha256:<computed from this file canonical JSON>",
  "effective_reference": ".github/workflows/ref-loop.yml@main",
  "operation_class": "DeployWorkflow",
  "effect_capable": true,
  "operation_identity_definition": {
    "format": "deploy:repo:workflow:ref:execution_run_id",
    "fields": {
      "operation": { "effect_relevant": true, "fixed": "deploy" },
      "repo": { "effect_relevant": true },
      "workflow": { "effect_relevant": true },
      "ref": { "effect_relevant": true },
      "execution_run_id": { "effect_relevant": true, "description": "immutable Reference Loop correlation identity" }
    },
    "canonicalization": "colon-separated string equality",
    "immaterial_parameters": []
  },
  "authorization_evidence_requirements": {
    "source": "Solvent kernel.AuthorizeAndCreateIntent",
    "fields": ["target_id", "snapshot_id", "intent_id", "intent_state"],
    "minimum_validity": "intent_state == 'live' at claim time",
    "freshness_requirement": "re-evaluated at execution boundary"
  },
  "validity_model": {
    "type": "single-use-intent",
    "intent_transitions": ["live", "executing", "executed"],
    "invalidation": ["belief_retraction", "target_revocation"]
  },
  "replay_idempotency": {
    "behavior": "single-use",
    "duplicate_rejection": "ClaimIntent CAS rejects when state != 'live'",
    "declare_exactly_once": true
  },
  "execution_outcome_model": {
    "source": "GitHub Actions API",
    "states": ["EXECUTION_ATTEMPTED", "EFFECT_CONFIRMED", "RESULT_OBSERVED", "AMBIGUOUS"],
    "ambiguous_handling": "do not retry automatically; log and leave for operator"
  },
  "reconciliation": {
    "owner": "human operator",
    "mechanism": "manual inspection of GitHub Actions runs + Solvent audit log",
    "automatic_reconciliation": false
  }
}
```

---

## D. Plan Scope

```text
plan_id:              phase2-first-effect
plan_version:         1
operation_class:      DeployWorkflow
approved_target_scope:
  repository:         pithomlabs/reference-loop-test
  workflow:           ref-loop.yml
  ref:                main
  effect_type:        workflow_dispatch trigger
  safety_boundary:    isolated test repository, no production impact
```

---

## E. Exact Operation Identity

```text
Format:     deploy:repo:workflow:ref:execution_run_id
Example:    deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:run-phase2-001

Effect-relevant fields (ALL):
  operation   = "deploy" (fixed)
  repo        = "pithomlabs/reference-loop-test"
  workflow    = "ref-loop.yml"
  ref         = "main"
  execution_run_id = unique immutable Reference Loop correlation identity

Immaterial fields: NONE

Canonicalization: string equality (deterministic, proven in Phase 1.5 Test G)
```

---

## F. Authorization Evidence

**What Solvent returns:**
```json
{
  "belief_id": "<uuid>",
  "intent_id": "<uuid>",
  "intent_state": "live",
  "action": "deploy",
  "authority": {
    "target_id": "<uuid>",
    "allowed": true,
    "reason": ""
  }
}
```

**How the Executor obtains it:**
- Calls `POST /v1/authorizations/action` with `action="deploy"`
- Receives authorization result with `intent_state: "live"`
- Calls `POST /v1/authorizations/execute` with `intent_id`

**How the Executor verifies it:**
- `PrepareForAction` re-reads current state from DB
- `kernel.Authorize` performs 8-field tuple comparison against approved snapshot
- `ClaimIntent` CAS atomically transitions `live -> executing` (sole authority gate)

**Trust basis:**
- Solvent kernel is the sole authority oracle
- Authorization is re-evaluated fresh at execution time (no cached authority)
- The Executor never constructs its own authorization evidence

**Validity requirements:**
- `intent_state == 'live'` at claim time
- No revocation exists for the target
- All justification beliefs are currently promoted
- Exact tuple match on all 8 authority fields

**Mismatch behavior:**
- Tuple mismatch -> `Allowed=false`, execution denied
- Intent not live -> `ErrIntentNotLive`, execution denied
- Revocation detected -> `Allowed=false`, execution denied

---

## G. External Effect Definitions

| Event | Definition | Source |
|-------|-----------|--------|
| EXECUTION_ATTEMPTED | `ClaimIntent` CAS succeeds, executor function invoked | Solvent audit + executor log |
| EFFECT_CONFIRMED | GitHub API returns `status=completed && conclusion=success` for `github_run_id` | GitHub Actions API |
| RESULT_OBSERVED | Agent receives `ExecuteAction` response from Solvent | Solvent REST response |
| AMBIGUOUS | GitHub API returns `status=in_progress` or HTTP error | GitHub Actions API |

**Do NOT equate executor invocation with external effect.** The executor may succeed locally but GitHub may reject or delay the workflow.

---

## H. Safety Model

| Property | Mechanism |
|----------|-----------|
| Reversible | GitHub workflow can be re-triggered; artifacts can be deleted |
| Sandboxed | `pithomlabs/reference-loop-test` is an isolated test repository |
| Dry-run | workflow_dispatch generates result.txt + artifact; no production deployment |
| Isolated | Separate from any production repository or service |

**Safety neutralization order:**
1. Authorization verification (Solvent kernel)
2. Exact operation binding (8-field tuple comparison)
3. Fail-closed enforcement (CAS gate)
4. THEN external effect (GitHub workflow)

**The production authorization/enforcement path is NOT bypassed.**

---

## 2. Pre-Registered Pass Criteria

```text
PASS-1:  Human-approved plan is recorded in Conductor
PASS-2:  Exact operation is constructed: deploy:repo:workflow:ref:run_id
PASS-3:  Solvent authorizes the exact operation (Allowed=true, intent_state=live)
PASS-4:  Authorization evidence reaches the effect-capable boundary (Executor)
PASS-5:  The effect-capable execution boundary performs fresh Solvent
        authorization verification before invoking the external-effect Executor.
        Executor != authority.
PASS-6:  Authorization binding is exact (ClaimIntent CAS succeeds)
PASS-7:  Executor produces the intended external effect (workflow_dispatch)
PASS-8:  GitHub API confirms EFFECT_CONFIRMED (status=completed, conclusion=success)
PASS-9:  Operation identity correlates end-to-end (same ID in proposal/authorization/execution/effect)
PASS-10: Conductor records coordination state only (no authority/effect ledger)
PASS-11: Solvent source unmodified (HEAD=7602699, clean)
```

---

## 3. Pre-Registered Failure / Falsification Criteria

```text
FALSIFY-1:  Authorized operation X cannot be distinguished from operation Y
            -> Architecture fails: operation binding is broken

FALSIFY-2:  External effect occurs without valid Solvent authorization
            -> Architecture fails: enforcement boundary is bypassed

FALSIFY-3:  Execution boundary cannot perform fresh Solvent authorization verification
            -> Architecture fails: trust boundary is not enforced

FALSIFY-4:  GitHub API cannot distinguish success from ambiguity
            -> Architecture fails: external SOR cannot provide authoritative evidence

FALSIFY-5:  Conductor must become an authority engine to make the loop safe
            -> Architecture fails: role ownership is violated

FALSIFY-6:  Agent-only behavior is required to preserve a consequential boundary
            -> Architecture fails: enforcement depends on Agent correctness

FALSIFY-7:  Plan scope cannot be recorded without turning Conductor into a policy engine
            -> Architecture fails: Conductor role is overloaded
```

**Note:** Ordinary implementation bugs (e.g., network timeout, API rate limit) are NOT architecture failures.

---

## 4. Handling NOT_YET_PROVEN Items

### A. UNKNOWN != DENIED

```text
Current frozen Solvent behavior:
    fail closed with Allowed=false for both UNKNOWN and DENIED

Required workflow semantic:
    UNKNOWN != DENIED (distinct states)

Decision for Phase 2:
    ACCEPTED FOR PHASE 2

Rationale:
    The frozen Solvent implements fail-closed behavior, which is the safer
    direction. UNKNOWN returning Allowed=false is conservative — it prevents
    execution when no authorization can be established.

    Phase 2 tests a positive path (authorized operation executes successfully).
    The UNKNOWN semantic distinction is relevant for negative paths where the
    Agent needs to decide between retry/wait/escalate. This is not exercised
    in Phase 2.

    The semantic gap is recorded as a known limitation. It does not block
    Phase 2 because:
    1. Phase 2 is a positive-path experiment
    2. Fail-closed behavior is the correct safety default
    3. The UNKNOWN distinction can be implemented in a future phase
```

### B. capability_ref -> Declaration Resolution

```text
Required mechanism:
    A declaration file (declaration_github_deploy_v1.json) that provides:
    - owner: reference-loop
    - version: v1.0.0
    - content hash: sha256 of canonical JSON
    - effective reference: .github/workflows/ref-loop.yml@main
    - retrievable content: the declaration file itself

Implementation:
    The declaration file is stored in the reference-loop repository.
    The Agent reads it to construct the operation identity.
    The Executor uses it to verify the operation class.
    This is NOT a general capability registry — it is a single file
    for the first experiment only.

    Implementation:
    1. Load declaration_github_deploy_v1.json
    2. Verify declaration_version == "v1.0.0"
    3. Compute sha256 of canonical JSON content
    4. Verify content_hash matches computed hash
    5. Verify effective_reference == ".github/workflows/ref-loop.yml@main"
    6. Only then does the task become READY/actionable

    This is NOT a general capability registry — it is a single file
    for the first experiment only, with an integrity gate.

    capability_ref resolution: READY (with integrity verification)
```

---

## 5. Pinned Components

```text
Workflow Design v1.0:
    file: docs/plan4/workflow_design_v1.0.md
    status: frozen

Conductor:
    commit: c243f27d9fb3e90ce9230e3e7cac03ad5e7bd7fa
    path: /home/chaschel/Documents/go/conductor

Reference Loop:
    commit: ebd4b9e782bd374a3e54e19ded0d16212eb50a52
    path: /home/chaschel/Documents/go/oracle/reference-loop

Solvent:
    commit: 7602699947e5fe5bafbeb0709d0a16b17d595318
    path: /home/chaschel/Documents/go/solvent-main
    status: FROZEN

Executor/integration:
    adapter: github.NewHTTPProvider
    executor_name: github_trigger_workflow
    backend: github (requires GITHUB_TOKEN env var)

Declaration:
    version: v1.0.0
    file: declaration_github_deploy_v1.json

External test repository:
    repo: pithomlabs/reference-loop-test
    workflow: ref-loop.yml
    ref: main
    trigger: workflow_dispatch
```

---

## 6. Evidence Package

```text
Human evidence:
    approved plan (plan_id, plan_version, approved_by, approved_at)

Conductor evidence:
    project state
    task lifecycle (created -> claimed -> submitted)
    activity log (actor, action, details, timestamp)
    dependency graph

Solvent evidence:
    authorization decision (allowed, reason, target_id, snapshot_id)
    intent lifecycle (live -> executing -> executed)
    audit_activity entries (authorization_granted, adapter_invoked, executor_completed)
    target_snapshot (authority tuple, justification_set, snapshot_hash)

Executor evidence:
    execution attempt (timestamp, params, provider response)
    execution outcome (accepted/rejected/ambiguous)
    intent state transition

External SOR evidence:
    GitHub Actions API response (run_id, status, conclusion, created_at)
    workflow_run details (id, name, head_branch, event)

Agent evidence:
    proposal construction (operation identity, capability_ref)
    interpretation of results
    continuation decision

Harness role:
    Correlates observations across participants
    Does NOT manufacture authoritative state
    Does NOT become an authority ledger
```

---

## 7. Phase 2 Scenario

```text
1. Human approves Plan v1
   - plan_id: phase2-first-effect
   - plan_version: 1
   - approved_by: human
   - approved_at: <timestamp>

2. Agent enters WORK mode
   - Human approval recorded in Conductor

3. Agent claims READY task
   - Task created in Conductor with capability_ref
   - Agent claims task exclusively

4. Agent reads and verifies declaration file
   - Load declaration_github_deploy_v1.json
   - Verify version == v1.0.0
   - Verify content_hash matches file content
   - Verify effective_reference == .github/workflows/ref-loop.yml@main
   - Task becomes READY/actionable only after verification

5. Agent constructs exact GitHub operation
   - operation_id: deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<execution_run_id>
   - All 5 fields populated

6. Agent sets up Solvent authorization
   - EnterBelief -> RetireDebt -> PromoteBelief
   - CreateTarget -> AttachJustification -> RequestAuthorization -> ApproveTarget

7. Solvent authorizes exact operation
   - AuthorizeAction(action="deploy", target_id=<approved>)
   - intent_state: "live"

8. Executor verifies authorization
   - PrepareForAction re-evaluates kernel.Authorize
   - 8-field tuple match against approved snapshot
   - ClaimIntent CAS: live -> executing

9. Executor triggers GitHub workflow
   - POST /repos/pithomlabs/reference-loop-test/actions/workflows/ref-loop.yml/dispatches
   - Body: { "ref": "main", "inputs": { "execution_run_id": "<execution_run_id>" } }

10. GitHub becomes external SOR
    - workflow_dispatch accepted (HTTP 204)
    - workflow run created

11. Effect is observed
    - Poll GitHub API: GET /repos/.../actions/runs?event=workflow_dispatch&per_page=1
    - Check status/conclusion

12. Conductor records coordination result
    - Task submitted
    - Activity logged

13. Agent continues
    - Next task discovered or loop complete
```

---

## 8. Phase 2 Stop Rule

After producing this document:

```text
STOP.

Do not:
  - implement the real external effect
  - modify Solvent
  - modify Conductor architecture
  - build a new capability registry
  - build a new workflow engine
  - create a new Executor
  - build a conformance harness
  - integrate BM-IST

This phase is planning only.
```

---

## 9. Final Output

```text
Phase 2 preparation: READY

First operation:     deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<execution_run_id>
External SOR:        GitHub Actions API (workflow run status/conclusion)
Declaration pinned:  YES (declaration_github_deploy_v1.json, v1.0.0)
Plan scope defined:  YES (phase2-first-effect, DeployWorkflow, isolated test repo)
Operation identity:  YES (deploy:repo:workflow:ref:execution_run_id, 5-part, all effect-relevant)
Authorization evidence: YES (Solvent kernel, 8-field tuple, CAS gate, fresh evaluation)
Safety mechanism:    Reversible, sandboxed, dry-run, isolated test repository
Pass criteria:       REGISTERED (11 criteria)
Failure criteria:    REGISTERED (7 falsification conditions)

UNKNOWN != DENIED:   ACCEPTED FOR PHASE 2 (fail-closed is safe for positive path)
capability_ref:      READY (declaration file, not a general registry)
Solvent modified:    NO
Architecture change: NO
```
