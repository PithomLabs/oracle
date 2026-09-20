# Phase 8 Implementation Plan — ARGUS Trust Verification POC

## 1. Phase 8 Objective

Phase 8 is a **thin first dry run** of the full ARGUS trust-verification loop. It exercises the complete architecture against one real BM–IST research slice (Gate G0), proving that:

1. A human creates a research task in Conductor, anchored to a Solvent scenario.
2. A Work OpenCode agent reconstructs canonical context via RCP (MCP → Coordinator → Conductor + Solvent).
3. The Work agent performs bounded research and submits an EBP packet.
4. The Coordinator persists claims, evidence, edges, and tasks through existing Conductor/Solvent REST clients.
5. A fresh Adversarial OpenCode agent (no conversation history) reconstructs the same canonical state via RCP.
6. The Adversarial agent attacks unresolved work, creating a `contradicts` edge against the governing belief.
7. The Trust UI shows the evolving research state (Insights) and human debt-discharge surface (Debts).
8. A human retracts the contradicted belief through the UI; Solvent cascades the retraction.
9. The Conductor task becomes terminal `cancelled`; Insights structurally derives a dead end.
10. A human discharges remaining debt through the UI; the Coordinator validates the applicable Domain Pack retirement rule before invoking Solvent.
11. Solvent authoritatively records attributed debt retirement with operator identity.
12. A human requests promotion; Solvent — not the UI or agent — decides whether the gate passes.

**What this proves:** An autonomous agent can do substantial research work, another autonomous agent can attack that work, and a human can adjudicate the resulting epistemic state through a thin control surface — without either agent being able to declare its own work authoritative.

## 2. Frozen Architecture

```
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Agent / OpenCode    = agency / work capability
Conductor           = operational workflow / task state
Solvent             = epistemic authority, belief/evidence/debt ledger,
                      authorization and consequential promotion gates
Trust UI            = human observation and adjudication
EBP v2.1            = epistemic operating discipline
RCP                 = thin read-only context protocol for agents
MCP                 = agent-facing transport
HTTP/API            = canonical integration contract
```

### EBP v2.1 doctrine (mandatory)

- Ideas enter free.
- Promotion costs debt.
- Debt does not kill.
- Debt is forever payable.
- New evidence creates new debt.
- No final-truth claim may be promoted.
- Accounting must never become the work.
- Human beings discharge epistemic debt. Agents do NOT directly retire debt, promote claims, authorize consequential actions, or mutate Solvent/Conductor directly.

### Data ownership

| Owner | State |
|-------|-------|
| Conductor | project, task, dependency, activity, task status, governance_ref, priority |
| Solvent | belief, belief_edge, evidence, debt, action_intent, audit_activity, promotion/retraction/authorization lifecycle |

No new persistence layer. No new Solvent schema. No new Conductor schema. No direct DB access from UI or agents.

## 3. Existing Interfaces

### Solvent REST API

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/v1/beliefs` | Enter a belief (ideas enter free) |
| GET | `/v1/beliefs/{id}` | Get belief |
| GET | `/v1/beliefs` | List beliefs (filter by status/claimType) |
| POST | `/v1/beliefs/{id}/debt/retire` | Retire one debt item (mechanical primitive) |
| POST | `/v1/beliefs/{id}/promote` | Promote (DB refuses if debt remains) |
| POST | `/v1/beliefs/{id}/retract` | Retract with cascade |
| GET | `/v1/beliefs/{id}/explain` | Explain promotion/authorization readiness |
| GET | `/v1/beliefs/{id}/evidence` | List evidence for a belief |
| POST | `/v1/evidence` | Add evidence to a belief |
| GET | `/v1/evidence/{id}` | Get evidence |
| POST | `/v1/discharge` | Record debt discharge with attribution |
| GET | `/v1/activity` | Read audit activity entries |
| GET | `/v1/ledger` | Ledger summary (6 aggregate counts) |
| **POST** | **`/v1/beliefs/{id}/edges`** | **Create edge (Phase 8 addition — see §4)** |

**Key types:**
- `EnterBeliefRequest{ScenarioID, Claim, ClaimType, Debt []string}`
- `AddEvidenceRequest{ScenarioID, BeliefID, ProvenanceClass, SourceURL, ContentSHA256}`
- `RetireDebtRequest{DebtItem string}`
- `DischargeRequest{ScenarioID, BeliefID, ObligationKey, InstrumentRef, DischargedBy}`
- `CreateEdgeRequest{ChildID, Kind}` — `Kind` is `"derives"` or `"contradicts"`

**Promotion gate:** DB CHECK constraint `promoted_is_debt_free` — `array_length(debt,1) > 0` OR `final_truth = true` → promotion refused (SQLSTATE 23514).

**Solvent MCP tools (18 total):** Read (ledger, explain, activity), mutation (retire_debt, promote, falsify, discharge, ingest_evidence), authority lifecycle (create_principal, revoke_principal, create_target, attach_justification, request_authorization, approve, authorize, revoke_target), execution (authorize_action, execute). No edge-creation MCP tool — agents never own edge mutations.

### Conductor REST API

| Method | Path | Purpose |
|--------|------|---------|
| GET/POST | `/v1/projects` | List/create projects |
| GET/PUT/DELETE | `/v1/projects/{id}` | CRUD project |
| GET/POST | `/v1/projects/{id}/tasks` | List/create tasks |
| GET/PATCH | `/v1/tasks/{id}` | Get/update task (title/description/priority only) |
| POST | `/v1/tasks/{id}/claim` | Claim (proposed → active) |
| POST | `/v1/tasks/{id}/submit` | Submit (active → review) |
| POST | `/v1/tasks/{id}/accept` | Accept (review → accepted) |
| POST | `/v1/tasks/{id}/reject` | Reject (review → active) |
| POST | `/v1/tasks/{id}/blocker` | Report blocker (active → blocked) |
| POST | `/v1/tasks/{id}/resolve` | Resolve blocker (blocked → active) |
| POST | `/v1/tasks/{id}/release` | Release (active → proposed) |
| POST | `/v1/tasks/{id}/activity` | Post activity |
| GET | `/v1/tasks/{id}/governance` | Get governance state (read-only) |

**Task state machine:**
```
proposed → active (claim) → review (submit) → accepted (accept) [terminal]
                ↓ blocked (report blocker) → active (resolve)
                ↓ cancelled [terminal]
active → proposed (release)
```

`rejected` transitions review → active (NOT terminal). Only `accepted` and `cancelled` are terminal.

**Governance ref:** `TEXT` column on task. Write-once at creation, immutable after. JSON: `{provider, reference_id}`.

**Priority:** `low | medium | high | critical`, default `medium`. POC-scoped manual assignment.

### Oracle Coordinator

**Existing HTTP endpoints:**
| Method | Path | Purpose |
|--------|------|---------|
| POST | `/decisions` | Submit human decision |
| GET | `/packets/{id}/status` | Get packet compilation status |
| GET | `/beliefs/{id}/decision-context` | Get decision context for a belief |
| GET | `/authorization-context/{id}` | Get Sphinx authorization projection |

**Existing Solvent client methods:** `CreateBelief`, `CreateEdge`, `CreateEvidence`, `PromoteBelief`, `RetireDebt`, `RetractBelief`, `ApproveTarget`, `GetBelief`

**Existing Conductor client methods:** `CreateTask`

**Critical gap:** `CompilePacket` creates in-memory IDs but does NOT persist to Conductor/Solvent. Phase 8 must close this gap.

### BM-IST Domain Pack

```json
{
  "pack_id": "bmist",
  "version": "1.0.0",
  "debt_vocabulary": [
    "needMap", "needInvariant", "needToyCheck",
    "needNullModel", "needObstruction", "needFaithfulnessReview"
  ],
  "initial_debt": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
  "evidence_classes": ["reproducible_artifact", "operator_asserted"],
  "retirement_rules": {
    "needMap": {"evidence_class": "reproducible_artifact", "rule": "map_check"},
    "needInvariant": {"evidence_class": "reproducible_artifact", "rule": "invariant_check"},
    "needToyCheck": {"evidence_class": "reproducible_artifact", "rule": "toy_model_check"},
    "needNullModel": {"evidence_class": "operator_asserted", "rule": "scope_clarification"},
    "needObstruction": {"evidence_class": "reproducible_artifact", "rule": "obstruction_construction"},
    "needFaithfulnessReview": {"evidence_class": "operator_asserted", "rule": "faithfulness_review"}
  }
}
```

### Gate G0 (dry-run subject)

From BM–IST Synthesis v5.1 §3:

- **L1 (Orbit Rigidity):** A strongly continuous unitary group cannot act nontrivially on a countable discrete/totally disconnected defined-state set. [THEOREM-LEVEL]
- **L2 (One-Parameter Triviality):** A finite signed-permutation operator group cannot itself support a nontrivial continuous one-parameter unitary flow. [THEOREM-LEVEL]
- **Gate G0:** countable substrate + literal continuous unitary evolution → structural incompatibility.
- **Required consequence:** discrete substrate time must be primitive; continuous Schrödinger evolution must be emergent.
- **Failure condition:** any construction that directly places literal continuous Schrödinger dynamics on the countable substrate fails G0.

Claims for dry run: G0 (main), L1, L2 (supporting). Dependencies: L1→G0, L2→G0.

## 4. Minimal Changes Required

### Solvent — One endpoint addition (no schema change)

| Change | Reason |
|--------|--------|
| Add `POST /v1/beliefs/{id}/edges` | Edge creation required for contradicts/derives edges. The `belief_edge` table already exists. No schema change. |

The endpoint accepts `{"child_id": "...", "kind": "derives|contradicts"}`. Solvent validates: parent exists, child exists, parent≠child, kind is valid, edge is unique. Then inserts into `belief_edge`.

**No MCP exposure.** Agents never call this endpoint directly. The Coordinator calls it when persisting EBP packets.

Add Solvent API tests: create derives edge, create contradicts edge, reject unknown parent, reject unknown child, reject self-edge, reject invalid kind, reject duplicate edge.

### Conductor — No changes

The existing schema and API are sufficient. Task lifecycle, governance_ref (write-once), and priority field all work as-is.

### Coordinator — Read methods + persistence + context endpoint

| Change | Reason |
|--------|--------|
| `coordinator/client.go` — add Solvent read methods: `ListBeliefs`, `ListEvidence`, `ListEdges`, `ListActivities`, `GetLedger`, `ListIntents` | RCP needs to read Solvent state |
| `coordinator/client.go` — add Conductor read methods: `GetTask`, `ListTasks`, `ListDependencies`, `ListActivity`, `GetGovernance` | RCP needs to read Conductor state |
| `coordinator/compiler.go` — add persistence stage after validation | `submit_packet` must actually persist beliefs/evidence/edges/tasks through REST clients |
| `coordinator/context.go` (new file) — `GetContext` method | Assembles RCP response from Conductor + Solvent reads |
| `coordinator/http/handler.go` — add `GET /v1/context/{task_id}` route | RCP endpoint |
| `coordinator/http/types.go` — add RCP request/response types | Protocol types |

### Summary

| Component | Change | Reason |
|-----------|--------|--------|
| `solvent/api/api.go` | Add route `POST /v1/beliefs/{id}/edges` | Edge creation |
| `solvent/api/edges.go` (new) | Handler + validation for edge creation | Edge creation |
| `solvent/api/types.go` | Add `CreateEdgeRequest` type | Edge creation |
| `coordinator/client.go` | Add Solvent + Conductor read methods | RCP projection |
| `coordinator/compiler.go` | Add persistence stage | submit_packet must persist |
| `coordinator/context.go` (new) | GetContext method | RCP endpoint logic |
| `coordinator/http/handler.go` | Add `/v1/context/{task_id}` route | RCP endpoint |
| `coordinator/http/types.go` | Add RCP types | Protocol types |

## 5. RCP Protocol Design

### Protocol identity

```json
{"protocol": "RCP/v1"}
```

### Canonical endpoint

```
GET /v1/context/{task_id}
```

### Conceptual response

```json
{
  "protocol": "RCP/v1",
  "task": {
    "id": "...",
    "project_id": "...",
    "title": "...",
    "description": "...",
    "status": "active",
    "priority": "high",
    "governance_ref": {
      "provider": "solvent",
      "reference_id": "track-g0"
    },
    "current_agent": "...",
    "created_at": "...",
    "updated_at": "..."
  },
  "dependencies": [
    {
      "task_id": "...",
      "blocked_by_id": "...",
      "blocked_by_status": "accepted",
      "source": "conductor"
    }
  ],
  "epistemic": {
    "available": true,
    "beliefs": [
      {
        "id": "...",
        "claim": "...",
        "claim_type": "derived",
        "status": "entered",
        "debt": ["needMap", "needInvariant"],
        "final_truth": false,
        "source": "solvent"
      }
    ],
    "evidence": [
      {
        "id": "...",
        "belief_id": "...",
        "provenance_class": "reproducible_artifact",
        "source_url": "...",
        "content_sha256": "...",
        "artifact_ref": "...",
        "source": "solvent"
      }
    ],
    "edges": [
      {
        "parent_id": "...",
        "child_id": "...",
        "kind": "derives",
        "source": "solvent"
      }
    ],
    "debt": [
      {
        "belief_id": "...",
        "items": ["needMap", "needInvariant"],
        "source": "solvent"
      }
    ],
    "intents": [
      {
        "belief_id": "...",
        "action": "...",
        "state": "live",
        "source": "solvent"
      }
    ]
  },
  "activity": [
    {
      "source": "conductor",
      "type": "task.claimed",
      "actor_id": "...",
      "details": "...",
      "at": "..."
    },
    {
      "source": "solvent",
      "type": "belief_entered",
      "actor_id": "...",
      "details": "...",
      "at": "..."
    }
  ]
}
```

### Design constraints

1. **v1 scope:** Full scenario projection. No sophisticated task-anchored graph closure.
2. **Cross-ledger consistency:** Eventually consistent. Each fact tags its source. Gaps are projection lag, not synthesized state.
3. **Source attribution:** Every returned object identifies whether it came from Conductor or Solvent.
4. **Availability metadata:** Each section includes `available` boolean and optional `reason` when a backend is unreachable. `UNKNOWN ≠ EMPTY`.
5. **No prescription:** RCP returns context. It does not prescribe scientific next steps.
6. **No new schema:** RCP objects are projections of existing Conductor + Solvent objects.

### Inclusion algorithm (v1)

RCP v1 returns the **full scenario projection** for the given task's scenario. The task-to-epistemic relationship is broad: all beliefs, evidence, edges, debt, and intents in the task's scenario are returned.

The task's `governance_ref.reference_id` provides the scenario identity. Dependencies show upstream task status. Activity merges Conductor and Solvent audit trails.

## 6. Coordinator Implementation

### RCP endpoint

```
GET /v1/context/{task_id}
```

Implementation in `coordinator/context.go`:

```go
func (c *Coordinator) GetContext(ctx context.Context, taskID string) (*RCPContext, error) {
    // 1. Get task from Conductor
    task, err := c.conductorClient.GetTask(ctx, taskID)
    if err != nil {
        return nil, fmt.Errorf("conductor: task not found: %w", err)
    }

    // 2. Extract scenario from governance_ref
    scenarioID := extractScenarioFromGovernance(task.GovernanceRef)

    // 3. Read from Conductor (fail-safe: unavailable ≠ empty)
    deps, depErr := c.conductorClient.ListDependencies(ctx, taskID)
    activity, actErr := c.conductorClient.ListActivity(ctx, taskID)

    // 4. Read from Solvent (fail-safe: unavailable ≠ empty)
    beliefs, belErr := c.solventClient.ListBeliefs(ctx, scenarioID)
    evidence, evErr := c.solventClient.ListEvidence(ctx, scenarioID)
    edges, edgeErr := c.solventClient.ListEdges(ctx, scenarioID)
    debt, debtErr := c.solventClient.GetDebt(ctx, scenarioID)
    intents, intErr := c.solventClient.ListIntents(ctx, scenarioID)
    solventActivity, solActErr := c.solventClient.ListActivities(ctx, scenarioID)

    // 5. Assemble response with availability metadata
    // Each section has available=true/false and reason when unavailable
    // NEVER silently convert errors to empty arrays
}
```

### Persistence stage

After `CompilePacket` validates and canonicalizes, `submit_packet` persists through REST:

```
EBP packet
  → Coordinator validation (rule check, idempotency)
  → Solvent POST /v1/beliefs (beliefs with initial debt)
  → Solvent POST /v1/beliefs/{id}/edges (derives/contradicts edges)
  → Solvent POST /v1/evidence (evidence items)
  → Conductor POST /v1/projects/{id}/tasks (tasks from packet)
```

Failure handling: if any step fails, return the error. Partial state is observable. Retry with same idempotency key does not duplicate completed objects. No distributed transactions.

### Retirement rule enforcement

The Coordinator mechanically validates the offered evidence class against the pack rule before invoking Solvent discharge:

```
human selects debt item
  → Coordinator identifies applicable retirement rule from pack
  → Coordinator identifies evidence offered for discharge
  → Coordinator compares evidence class against rule's required class
  → MISMATCH → refusal ("evidence class X does not satisfy rule requiring Y")
  → MATCH → invoke Solvent POST /v1/discharge with attributed identity
```

The Coordinator does NOT judge scientific correctness. It enforces mechanical class matching only.

## 7. MCP Capability Boundary

### Agent tool surface

OpenCode Work and Adversarial processes receive exactly:

```
argus.get_context
argus.submit_packet
```

They MUST NOT have:
- `solvent-mcp` (no `solvent_*` tools)
- Conductor write MCP tools (no `conductor_*` tools)
- Direct DB access

The capability boundary is enforced by tool-surface configuration, not prompt language.

### MCP adapter (new file: `mcp/adapter.go`)

```go
// Tool: argus.get_context
// Input: {"task_id": "..."}
// Calls: Coordinator.GetContext
// Returns: RCP context response

// Tool: argus.submit_packet
// Input: EBP packet JSON
// Calls: Coordinator.SubmitPacket
// Returns: compilation + persistence result
```

### Tests

- Only two ARGUS tools exposed
- `solvent_*` tools absent
- `conductor_*` tools absent
- Attempted direct Solvent mutation fails (tool not available)
- Attempted direct Conductor mutation fails (tool not available)

## 8. Work Agent Workflow

```
Human task (in Conductor)
  → argus.get_context (RCP)
  → fresh OpenCode sees task/dependencies/beliefs/evidence/debt/activity
  → bounded research on Gate G0
  → EBP work packet:
      claims: L1 (derives→G0), L2 (derives→G0), G0 (main)
      evidence: reproducible_artifact
      edges: L1→G0 (derives), L2→G0 (derives)
  → argus.submit_packet
  → Coordinator validates, persists to Conductor + Solvent
```

### Packet reference grammar (frozen)

```
local:<id>              = intra-packet reference (within the submitted packet)
canonical:belief:<uuid> = reference to an already-existing Solvent belief
```

All examples in the plan and implementation must use these conventions consistently.

## 9. Adversarial Agent Workflow

```
fresh OpenCode process (no conversation history)
  → argus.get_context (RCP)
  → sees current task/dependencies/beliefs/evidence/debt/activity
  → attacks unresolved work on G0
  → creates contradiction:
      claim: "G0 fails because..." (adversarial finding)
      edge: contradicts → canonical:belief:<G0-belief-id>
      evidence: reproducible_artifact
  → EBP adversarial packet
  → argus.submit_packet
  → Coordinator validates, persists contradicts edge to Solvent
```

The adversarial packet creates an actual `contradicts` relationship against the target belief. This is structurally visible in the Solvent ledger and derivable in Insights as "Adversarial Challenges".

## 10. Human Identity (C3)

Human authority is a server-side capability boundary, not a browser claim.

```
ARGUS_OPERATOR_PRINCIPAL_ID
```

The path:

```
Browser
  → Trust UI (injects configured operator identity)
  → Coordinator (validates operator identity)
  → Solvent POST /v1/discharge (DischargedBy = configured operator)
  → audit_activity.actor_id = operator identity
```

**Requirements:**
- Trust UI injects the configured operator identity from server-side configuration
- Coordinator validates that operator identity is present and non-empty
- Coordinator uses configured identity for all human decisions (not browser-supplied)
- Solvent discharge receives configured identity as `DischargedBy`
- Missing operator identity fails closed (HTTP 400 or 500, never proceeds)
- `DecisionRequest` does NOT trust an arbitrary `actor` field from the browser

**Recorded in:**
- Coordinator `DecisionRecord`
- Solvent discharge audit trail (`audit_activity.actor_id`)

## 11. EBP / Debt Workflow

### Initial debt

BM-IST pack declares 6 initial debt items. When the Work agent enters beliefs via `argus.submit_packet`, the Coordinator calls Solvent `POST /v1/beliefs` with:

```json
{
  "scenario_id": "track-g0",
  "claim": "...",
  "claim_type": "derived",
  "debt": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"]
}
```

Beliefs enter free with pack initial debt. Promotion is blocked until all debt is discharged.

### Debt discharge flow

```
Trust UI (Debts screen)
  → human selects debt item, reviews evidence
  → clicks RETIRE DEBT
  → Trust UI sends decision to Coordinator
  → Coordinator:
      1. Validates operator identity is present
      2. Loads pack retirement rule for this debt item
      3. Validates offered evidence class matches rule's required class
      4. If mismatch → refusal response ("evidence class X does not satisfy rule requiring Y")
      5. If match → calls Solvent POST /v1/discharge:
           - scenario_id
           - belief_id
           - obligation_key (the debt item)
           - instrument_ref (Coordinator decision reference)
           - discharged_by (configured operator principal)
  → Solvent records discharge in audit_activity
  → UI updates: "Debt retired by operator"
```

### Language rules

The UI MUST NEVER display:
- "AI verified debt"
- "Agent retired debt"
- "Automatically discharged"

The UI MUST display:
- "Adversarial agent supplied evidence relevant to this debt."
- "Human decision required."
- "Debt retired by operator."

### Promotion

```
human requests promotion
  → Coordinator
  → Solvent POST /v1/beliefs/{id}/promote
  → Solvent gate (DB CHECK: promoted_is_debt_free)
  → promotion or refusal
```

The UI/agent NEVER declares promotion authoritative. Solvent decides.

### New evidence creates new debt

When new evidence enters the system (via agent submission or human action), the Coordinator ensures the associated belief's debt is updated per the pack's rules. The belief re-enters the debt-discharge cycle.

## 12. Retraction / Reopen / Dead End

### Retraction path

```
human retracts belief
  → Trust UI
  → Coordinator
  → Solvent POST /v1/beliefs/{id}/retract (RetractCascade)
  → Solvent:
      1. Cancels dependent live intents (via descendantsCTE + belief_edge)
      2. Retracts the belief and its descendants
  → dependent Conductor tasks become terminal cancelled
```

### REOPEN

```
retracted belief
  → human requests REOPEN
  → Coordinator
  → Solvent: new belief created
  → Solvent: derives edge from original
  → original remains in history
  → lineage preserved
```

### Dead-end definition (H8)

```
Dead End =
  terminal cancelled Conductor task
  + governance/scenario linkage to a retracted or contradicted governing belief
```

**Not dead end:**
- `cancelled` alone (task cancelled for other reasons)
- `rejected` (returns to active, not terminal)
- `active` + retracted belief (task not terminal)

**Structurally derived:** The dead-end status is computed by the Coordinator from Conductor task state + Solvent belief state. It is NOT agent-written metadata.

### Dry-run branch (H3)

The Phase 8 dry run must exercise:

```
Work Agent creates G0 belief
  → Adversarial Agent creates contradicts edge against G0
  → human retracts G0
  → Solvent RetractCascade (using belief_edge for transitive closure)
  → Conductor task becomes terminal cancelled
  → Insights derives Dead End
```

This branch proves that dead-end machinery is not merely implemented-but-unused.

## 13. Trust UI

### Module boundary (H10)

Trust UI is a **separate Go module** at `trust-ui/`, NOT inside `oracle/`.

```
trust-ui/
  go.mod
  server.go
  handlers.go
  templates/
    shell.html
    insights.html
    debts.html
  static/
    style.css
    app.js
```

Trust UI communicates with Coordinator over HTTP. It never writes directly to Solvent or Conductor. This preserves the frozen two-Go-module architecture.

### Visual pattern

Uses the existing Solvent wizard style (vanilla HTML/CSS/JS, `html/template`). Same pattern as `solvent-main/internal/wizard/`.

### Two surfaces

**Sidebar:** Insights | Debts

**No other surfaces.** The UI is not a general-purpose dashboard. It is a debt-discharge and situational-awareness surface.

### Insights

Compact program statistics:

| Field | Source |
|-------|--------|
| Active Tasks | Conductor tasks with status=active |
| Adversarial Challenges | Count of `contradicts` edges in scenario (Solvent) |
| Open Claims | Solvent beliefs with status=entered |
| Open Debt | Solvent beliefs with debt items remaining |
| Promoted | Solvent beliefs with status=promoted |
| Refusals | Solvent promotion refusals (from audit_activity) |
| Dead Ends | Derived: terminal cancelled Conductor tasks + retracted/contradicted governing belief |

**Priority ordering:** Tasks ordered by Conductor `priority` (critical → high → medium → low), then `created_at` deterministic secondary. Client-side sorting, no SQL.

### Debts

For each belief with open debt:

| Field | Source |
|-------|--------|
| Claim | Solvent belief.claim |
| Debt Item | Solvent belief.debt[] |
| Why Debt Exists | Pack initial_debt + context |
| Evidence | Solvent evidence for this belief |
| Adversarial Challenges | Contradicts edges targeting this belief |
| Applicable Retirement Rule | Pack retirement_rules[item] |
| Proposed Discharge | Coordinator-computed rule match |
| History | Solvent audit_activity for this belief |
| RETIRE DEBT | Button → Coordinator → Solvent discharge |
| KEEP OPEN | Button → no action |

**Language:** "Debt retired by operator" (never "AI verified" or "Agent retired").

## 14. Artifact Retrieval

Evidence includes `artifact_ref` where available. Hash-only is insufficient for fresh agents that need to inspect actual evidence content.

The Coordinator resolves `artifact_ref` to a retrievable path when assembling the RCP context. Where the artifact is a local file, the `artifact_ref` is a path relative to a known base directory.

For Phase 8 dry run, artifacts are small text files in the domain pack's evidence directory.

## 15. Eventual Consistency

RCP is eventually consistent across Conductor and Solvent stores. Each fact tags its source (`conductor` or `solvent`). Response sections include availability metadata:

```json
{
  "epistemic": {
    "available": true,
    "beliefs": [...],
    ...
  }
}
```

When a backend is unavailable:

```json
{
  "epistemic": {
    "available": false,
    "reason": "solvent_unavailable",
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  }
}
```

**UNKNOWN ≠ EMPTY.** A Solvent outage must never look like "no research exists." The agent must never infer "no work has happened" from a failed backend read.

Projection gaps are visible as lag/degradation, not silently converted into empty state.

## 16. Error / Refusal Semantics

| Error | Response |
|-------|----------|
| Unknown task_id | HTTP 404, RCP context not available |
| Solvent unavailable | RCP returns `epistemic.available=false`, `reason="solvent_unavailable"` |
| Conductor unavailable | RCP returns task section unavailable |
| Retirement rule mismatch | Coordinator refuses, returns "evidence class X does not satisfy rule requiring Y" |
| Missing operator identity | HTTP 400, fail closed |
| Promotion with open debt | Solvent returns SQLSTATE 23514 refusal |
| Final-truth promotion | Solvent returns SQLSTATE 23514 refusal |
| Duplicate edge | Solvent returns HTTP 409 Conflict |
| Self-edge | Solvent returns HTTP 400 Bad Request |
| Invalid edge kind | Solvent returns HTTP 400 Bad Request |

## 17. Security / Capability Boundary

### Enforcement

The capability boundary is enforced by **tool-surface configuration**, not prompt language.

OpenCode Work and Adversarial processes are configured with exactly two MCP tools:
- `argus.get_context`
- `argus.submit_packet`

Solvent MCP tools (`solvent_*`) and Conductor MCP tools are NOT in the agent configuration. The agent cannot call what it cannot see.

### Tests

- Verify only two ARGUS tools are exposed in agent config
- Verify `solvent_*` tools are absent from agent config
- Verify `conductor_*` tools are absent from agent config
- Attempted direct Solvent mutation fails (no tool available)
- Attempted direct Conductor mutation fails (no tool available)

## 18. Idempotency Scoping (H7)

The canonical packet idempotency hash **includes scenario identity**:

```
canonical_hash = hash(
  scenario_id
  + sorted_claims
  + sorted_evidence
  + sorted_edges
  + packet_content
)
```

**Excluded from hash:**
- packet_id
- run_id
- timestamps
- runtime metadata

**Rule:**
- Same content + same scenario → deduplicates
- Same content + different scenario → does NOT deduplicate

Preserves the existing atomic NEW / IN_FLIGHT / COMPLETED / FAILED semantics from Phase 6.

## 19. Implementation Plan Structure

The revised document includes:

1. Objective
2. Frozen architecture
3. Existing interfaces
4. Minimal changes
5. RCP
6. Coordinator
7. MCP
8. Work Agent
9. Adversarial Agent
10. Persistence
11. EBP/debt workflow
12. Human identity
13. Retraction/reopen/dead-end
14. Trust UI
15. Insights
16. Debts
17. Artifact retrieval
18. Eventual consistency
19. Error/refusal semantics
20. Security/capability boundary
21. Idempotency scoping
22. Test strategy
23. Detailed tests
24. Acceptance criteria
25. Non-goals
26. Freeze/reconciliation notes
27. Risks

## 20. Test Strategy

### Unit tests

**Solvent edge endpoint:**
- Create derives edge (success)
- Create contradicts edge (success)
- Reject unknown parent (404)
- Reject unknown child (404)
- Reject self-edge (400)
- Reject invalid kind (400)
- Reject duplicate edge (409)

**Coordinator persistence:**
- Work packet → beliefs persisted
- Work packet → derives edges persisted
- Adversarial packet → contradicts edge persisted
- Evidence persisted with correct belief linkage
- Tasks persisted with correct governance_ref
- Canonical ID mapping (local→canonical references)
- Duplicate submit → idempotent (no duplicate objects)
- Cross-scenario idempotency isolation (same content, different scenario → no dedup)
- Partial persistence failure (belief created, evidence fails → error returned, partial state observable)
- Retry with same idempotency key → no duplicates

**Coordinator context (RCP):**
- Valid context (task exists, scenario has beliefs)
- Unknown task → 404
- Empty state (task exists, no beliefs yet)
- Dependency projection
- Source attribution (each object tagged conductor/solvent)
- Deterministic response (same inputs → same shape)
- Solvent unavailable → `epistemic.available=false`, `reason="solvent_unavailable"`
- Conductor unavailable → task section unavailable
- Projection lag visible (tags indicate freshness)

**Retirement rule enforcement:**
- NeedMap + reproducible_artifact → allowed
- NeedMap + operator_asserted → refused
- NeedNullModel + operator_asserted → allowed
- NeedNullModel + reproducible_artifact → refused
- Unknown debt item → error

**Human identity:**
- Missing operator identity → fail closed (HTTP 400)
- Configured operator identity → reaches Solvent audit_activity.actor_id
- Browser-supplied actor cannot override configured identity

### Integration tests

**Full dry-run path:**
- Work agent → get_context → research → submit_packet → beliefs/edges/evidence/tasks persisted
- Fresh adversarial agent → get_context → sees previous work → submit_packet → contradicts edge persisted
- Human retracts G0 → Solvent cascade → Conductor task cancelled → dead end derived
- Human discharges debt → Coordinator validates rule → Solvent records attributed discharge
- Human promotes → Solvent gate → promotion or refusal (depending on remaining debt)

**Adversarial reconstruction:**
- Fresh agent sees previous work via RCP
- Fresh agent sees previous evidence via RCP
- Fresh agent sees open debt via RCP
- Fresh agent sees prior adversarial findings (contradicts edges) via RCP
- Adversarial packet creates contradicts edge

**Dead end derivation:**
- Cancelled task + retracted governing belief → dead end
- Cancelled task + active belief → not dead end
- Rejected task + active belief → not dead end (rejected is not terminal)
- Contradicted governing belief → dead end

### Regression tests

- Existing Coordinator tests pass
- Existing Solvent tests pass
- Existing Conductor tests pass
- Reference-loop unchanged
- Domain-pack validation unchanged
- Packet validation unchanged
- Verifier unchanged
- Corpus unchanged

### Concurrency tests

- RCP concurrent reads (no races)
- Submit_packet duplicate concurrency (idempotent)
- Concurrent discharge remains safe/idempotent

### Capability boundary tests

- Only two ARGUS tools exposed in agent config
- Solvent tools absent from agent config
- Conductor tools absent from agent config
- Attempted direct Solvent mutation fails
- Attempted direct Conductor mutation fails

## 21. Detailed Tests

### Solvent edge endpoint tests

```go
func TestCreateEdge_Derives(t *testing.T)       // Success
func TestCreateEdge_Contradicts(t *testing.T)    // Success
func TestCreateEdge_UnknownParent(t *testing.T)  // 404
func TestCreateEdge_UnknownChild(t *testing.T)   // 404
func TestCreateEdge_SelfEdge(t *testing.T)       // 400
func TestCreateEdge_InvalidKind(t *testing.T)    // 400
func TestCreateEdge_Duplicate(t *testing.T)      // 409
func TestCreateEdge_Uniqueness(t *testing.T)     // Same edge twice → 409
```

### Coordinator persistence tests

```go
func TestSubmitPacket_WorkBeliefsPersisted(t *testing.T)
func TestSubmitPacket_DerivesEdgesPersisted(t *testing.T)
func TestSubmitPacket_ContradictsEdgePersisted(t *testing.T)
func TestSubmitPacket_EvidencePersisted(t *testing.T)
func TestSubmitPacket_TasksPersisted(t *testing.T)
func TestSubmitPacket_LocalRefMapping(t *testing.T)
func TestSubmitPacket_DuplicateSubmitIdempotent(t *testing.T)
func TestSubmitPacket_CrossScenarioNoDedup(t *testing.T)
func TestSubmitPacket_PartialFailureObservable(t *testing.T)
func TestSubmitPacket_RetryNoDuplicates(t *testing.T)
```

### RCP context tests

```go
func TestGetContext_ValidTask(t *testing.T)
func TestGetContext_UnknownTask(t *testing.T)
func TestGetContext_EmptyState(t *testing.T)
func TestGetContext_DependencyProjection(t *testing.T)
func TestGetContext_SourceAttribution(t *testing.T)
func TestGetContext_Deterministic(t *testing.T)
func TestGetContext_SolventUnavailable(t *testing.T)
func TestGetContext_ConductorUnavailable(t *testing.T)
func TestGetContext_ProjectionLagVisible(t *testing.T)
func TestGetContext_ArtifactRefPresent(t *testing.T)
```

### Retirement rule enforcement tests

```go
func TestRetireDebt_NeedMapReproducibleArtifact(t *testing.T)
func TestRetireDebt_NeedMapOperatorAsserted(t *testing.T)
func TestRetireDebt_NeedNullModelOperatorAsserted(t *testing.T)
func TestRetireDebt_NeedNullModelReproducibleArtifact(t *testing.T)
func TestRetireDebt_UnknownDebtItem(t *testing.T)
```

### Human identity tests

```go
func TestDischarge_MissingOperatorFailsClosed(t *testing.T)
func TestDischarge_ConfiguredOperatorReachesAudit(t *testing.T)
func TestDischarge_BrowserActorCannotOverride(t *testing.T)
```

### Adversarial reconstruction tests

```go
func TestAdversarial_FreshAgentSeesPreviousWork(t *testing.T)
func TestAdversarial_FreshAgentSeesPreviousEvidence(t *testing.T)
func TestAdversarial_FreshAgentSeesOpenDebt(t *testing.T)
func TestAdversarial_FreshAgentSeesPriorFindings(t *testing.T)
func TestAdversarial_PacketCreatesContradictsEdge(t *testing.T)
```

### Retraction and dead-end tests

```go
func TestRetract_CascadeEffect(t *testing.T)
func TestRetract_DependentIntentsCancelled(t *testing.T)
func TestRetract_OriginalBeliefRetained(t *testing.T)
func TestReopen_CreatesDerivesLineage(t *testing.T)
func TestDeadEnd_CancelledPlusRetracted(t *testing.T)
func TestDeadEnd_CancelledPlusActive(t *testing.T)
func TestDeadEnd_RejectedPlusActive(t *testing.T)
func TestDeadEnd_ContradictedGoverningBelief(t *testing.T)
```

### UI tests

```go
func TestInsights_Statistics(t *testing.T)
func TestInsights_ConductorPriorityOrdering(t *testing.T)
func TestInsights_DeadEndCount(t *testing.T)
func TestInsights_AdversarialChallengeCount(t *testing.T)
func TestDebts_Rendering(t *testing.T)
func TestDebts_EvidenceAndRuleDisplay(t *testing.T)
func TestDebts_NoAILanguage(t *testing.T)
func TestDebts_DischargePath(t *testing.T) // UI → Coordinator → Solvent
func TestDebts_NoDirectDBWrites(t *testing.T)
```

### Regression tests

```go
func TestRegression_CoordinatorExistingTests(t *testing.T)
func TestRegression_SolventExistingTests(t *testing.T)
func TestRegression_ConductorExistingTests(t *testing.T)
func TestRegression_ReferenceLoopUnchanged(t *testing.T)
func TestRegression_DomainPackValidation(t *testing.T)
func TestRegression_PacketValidation(t *testing.T)
func TestRegression_Verifier(t *testing.T)
func TestRegression_Corpus(t *testing.T)
```

### Concurrency tests

```go
func TestConcurrency_RCPConcurrentReads(t *testing.T)
func TestConcurrency_SubmitPacketDuplicate(t *testing.T)
func TestConcurrency_ConcurrentDischargeSafe(t *testing.T)
```

### Capability boundary tests

```go
func TestCapability_TwoARGUSToolsOnly(t *testing.T)
func TestCapability_SolventToolsAbsent(t *testing.T)
func TestCapability_ConductorToolsAbsent(t *testing.T)
func TestCapability_DirectSolventMutationFails(t *testing.T)
func TestCapability_DirectConductorMutationFails(t *testing.T)
```

## 22. Acceptance Criteria

Every mechanism promised below is described in implementation sections above:

1. **Fresh-agent reconstruction:** A completely fresh adversarial OpenCode process (no conversation history) can reconstruct prior work via `argus.get_context`. RCP returns full scenario projection with beliefs, evidence, edges, debt, and activity.

2. **Refusal-before-promotion:** Solvent refuses promotion when debt remains (SQLSTATE 23514). The agent or UI cannot override this.

3. **Solvent-decides-the-gate:** Only Solvent can authorize promotion. The UI and agents can request but never declare.

4. **Human-discharged-debt:** Debt retirement uses attributed `POST /v1/discharge` with `DischargedBy = configured operator principal`. The Solvent audit trail contains the operator identity.

5. **Retirement-rule enforcement:** The Coordinator mechanically validates the offered evidence class against the pack rule before invoking Solvent. Mismatched classes are refused.

6. **Adversarial challenge:** The adversarial agent creates an actual `contradicts` edge against the target belief. "Adversarial Challenges" in Insights is derived from `contradicts` edges.

7. **Dead-end derivation:** Dead end is structurally derived from terminal cancelled Conductor task + retracted/contradicted governing belief. Not agent-written metadata.

8. **Retraction exercise:** The dry run exercises contradiction → retraction → Solvent cascade → cancelled task → dead end.

9. **Capability boundary:** Only `argus.get_context` and `argus.submit_packet` are available to agents. Solvent and Conductor MCP tools are absent. Attempted direct mutation fails.

10. **No silent failures:** RCP distinguishes unavailable from empty. Backend errors surface as `available=false` with `reason`.

11. **Operator identity:** Missing operator identity fails closed. Configured identity reaches Solvent audit. Browser cannot override.

12. **Idempotency isolation:** Same packet content in different scenarios does not deduplicate.

13. **No UI bypass:** No Trust UI write path bypasses Coordinator. No direct DB writes from UI.

14. **No AI language:** UI never displays "AI verified debt," "Agent retired debt," or "Automatically discharged."

15. **All existing tests pass:** Coordinator, Solvent, Conductor, domain-pack, packet validation, verifier, corpus tests all pass unchanged.

## 23. Non-Goals

- No new Solvent schema
- No new Conductor schema
- No new persistence layer
- No general-purpose graph API
- No sophisticated task-anchored graph traversal (deferred to v1.1)
- No third UI surface (retraction uses existing Coordinator decision path)
- No MCP edge-creation tool for agents
- No direct DB access from Coordinator, UI, or agents
- No distributed transactions
- No new research database
- No priority field in Solvent
- No agent-inferred dead-end metadata

## 24. Freeze / Reconciliation Notes

### Conductor state machine reconciliation

The actual Conductor lifecycle (verified against `conductor/internal/domain/task.go`):

```
proposed → active → review → accepted [terminal]
   ↓          ↓        ↓
cancelled  cancelled  cancelled
           ↓
         proposed (release)
           ↓
         blocked → active (resolve)
```

`rejected` transitions review → active (NOT terminal). The only terminal states are `accepted` and `cancelled`.

**Dead-end definition uses `cancelled` as the terminal state.** Any older documentation referencing `rejected` as a terminal dead-end state is superseded by this reconciliation.

### Edge creation reconciliation

The `belief_edge` table exists in Solvent's schema and is used by `RetractCascade` for transitive closure. However, no REST endpoint or MCP tool creates edges. The only edge creation in existing code is direct SQL in test setup.

Phase 8 adds `POST /v1/beliefs/{id}/edges` as the lawful write interface. This is a Growth Gate exception: extending Solvent's interface, not its schema or authority model.

### Discharge path reconciliation

Solvent exposes two debt-retirement paths:
- `POST /v1/beliefs/{id}/debt/retire` — mechanical primitive, no attribution
- `POST /v1/discharge` — attributed discharge with `DischargedBy` and `InstrumentRef`

Phase 8 human adjudication uses the attributed path (`POST /v1/discharge`). The bare retire endpoint remains available as a lower-level primitive but is not used by the human workflow.

### Priority reconciliation

Conductor's `priority` field is `low | medium | high | critical`, default `medium`. It is manually assigned and POC-scoped. Insights ordering is client-side: `critical > high > medium > low`, then `created_at` as deterministic secondary. No SQL derivation. No priority field in Solvent.

### Trust UI module reconciliation

The frozen architecture places Trust UI as a separate Go module at `trust-ui/`, communicating with Coordinator over HTTP. Any plan text placing Trust UI inside `oracle/` is superseded by this reconciliation.

### Agent evidence provenance reconciliation

An autonomous agent's output is NOT automatically `operator_asserted`. Agent analysis produces candidate evidence. Candidate evidence becomes authoritative retirement material only through:
- A reproducible artifact satisfying the pack rule, OR
- Explicit human attestation where the pack permits `operator_asserted`

Do not overload `operator_asserted` to mean "agent acted on behalf of operator."

## 25. Risks

| Risk | Mitigation |
|------|------------|
| Solvent edge endpoint not implemented in time | Verify endpoint exists before coding persistence. If not, stop and add it first. |
| `RetractCascade` depends on edges that don't exist yet | Create edges during persistence before any retraction is attempted. |
| RCP availability metadata adds complexity | Keep it simple: each section has `available` bool + optional `reason`. Default `true`. |
| Retirement rule enforcement requires pack loading | Coordinator loads pack at startup. Rules are in-memory. |
| Trust UI module separation creates build complexity | Use `go.work` workspace. Trust UI imports Coordinator as HTTP client. |
| Partial persistence leaves inconsistent state | Partial state is observable. No distributed transactions. Idempotency preserves safety on retry. |
| Fresh adversarial agent misinterprets RCP context | RCP is explicitly non-prescriptive. Agent reasoning is outside Phase 8 scope. |

---

**Status:** Revised implementation plan. Awaiting implementation approval.
