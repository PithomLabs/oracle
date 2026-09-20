# ARGUS POC Simplification Pivot — Implementation Plan (Revised)

Revised per 12 consolidated amendments from adversarial review (prompt_review2.md + follow-up).

Status: PLAN ONLY — no implementation code.

==================================================
EXECUTIVE SUMMARY
==================================================

The POC infrastructure has drifted from the thesis it demonstrates. Recent fixes are increasingly about distributed plumbing (DB orchestration, auth propagation, artifact handoff, readiness probes) rather than ARGUS itself.

This plan collapses the multi-service stack (Solvent, Conductor, Coordinator, MCP, Trust UI) into:

    one Go binary (argus)
    one CockroachDB database
    two process instances (argus serve + argus mcp) + OpenCode

while preserving the architectural invariants:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

The pivot is reconnaissance-gated. Phase 0 must produce a GO / NO-GO / REVISE decision before any code changes.

Amendments incorporated:

1. Real Phase 0 gate with NO-GO/REVISE triggers
2. Solvent module boundary resolved via public projection contract
3. Agent packets never perform debt discharge
4. Agent-authored edges remain non-authoritative proposals
5. Verifier trust binding (id, version, input_hash, tolerance)
6. Separate-transaction model (Solvent authority mutations stay atomic within Solvent)
7. DB-enforced idempotency (content-hash unique constraint)
8. Typed PackDefinition for domain-pack contract
9. Explicit Conductor→CRDB migration with type mapping
10. Authentication mechanism + refusal test suite
11. Taskfile health polling, honest complexity metrics
12. Migration ownership: Solvent owns epistemic schema, ARGUS owns work schema

==================================================
COMPLEXITY CLASSIFICATION
==================================================

Before inspecting code, classify each current pain point:

| Pain point | Classification | Pivot fixes? |
|------------|---------------|-------------|
| Cross-service reference validation | API/integration | Yes — same-DB FK |
| Verifier artifact handoff | deployment | Yes — in-process |
| Auth propagation across services | deployment | Yes — no network |
| RCP projecting from two backends | deployment | Yes — one DB |
| Solvent HTTP server overhead | deployment | Yes — removed |
| Conductor SQLite + separate process | deployment | Yes — merged |
| Coordinator HTTP routing layer | deployment | Yes — direct calls |
| Trust UI separate process | deployment | Yes — embedded |
| Packet validation complexity | inherent domain | No — preserved |
| Epistemic authority model | inherent domain | No — preserved |
| Domain pack vocabulary rules | inherent domain | No — preserved |

Verdict: the majority are deployment complexity. The pivot is justified.

==================================================
RECONNAISSANCE FINDINGS
==================================================

## Solvent — reusable in-process

The kernel is a single `kernel.Store` struct wrapping `*sql.DB`:

- 22-method Contract interface (EnterBelief, AddEvidence, RetireDebt, Promote, RetractCascade, IntentOnPromoted, Discharge, CreatePrincipal, CreateTarget, Approve, Authorize, etc.)
- All writes go through `crdb.ExecuteTx` for CockroachDB serialization retry
- SQL is DB-independent (uses `:1`-style placeholders)
- No HTTP dependency — pure Go library

Services that wrap the kernel (all take `*sql.DB`):
- `service/ledger.Service` — API-facing thin passthrough + audit
- `service/audit.Service` — append-only activity ledger
- `service/authority.Service` — execution boundary (kernel + policy + audit + executor)
- `service/policy.Service` — constraint evaluation
- `service/executor.Registry` — in-memory action registry

Read-only projections (currently internal/):
- `internal/view.GetSnapshot` — beliefs + evidence + intents + audit (takes *sql.DB)
- `internal/view.ExplainSnapshot` — human-readable "why" projection (pure function of Snapshot)

Discardable (HTTP/wizard):
- `api/` — all HTTP handlers
- `internal/wizard/` — interactive HTTP wizard
- `cmd/solvent-api/`, `cmd/solvent-mcp/` — standalone binaries

Schema: 10 migration files, 20+ tables. Core tables: belief, belief_edge, evidence, action_intent, principal, authority_target, target_snapshot, target_activation, target_revocation, justification, debt_discharge, audit_activity.

## Conductor — extract minimal domain model

4 domain types:
- `Task` — 6 states (proposed/active/review/accepted/blocked/cancelled), 4 priorities
- `Project` — container for tasks
- `Dependency` — task-to-task blocking
- `Activity` — append-only event log

State machine: atomic compare-and-swap in SQL. Claim is `WHERE status='proposed' AND current_agent IS NULL`. Transition is `WHERE status=<from>`.

Storage: SQLite via `modernc.org/sqlite` (pure Go, no CGO). 4 tables: conductor_project, conductor_task, conductor_activity, conductor_dependency.

Governance: read-only `GovernanceReader` interface with null fallback.

Discardable:
- `internal/api/` — HTTP REST API
- `internal/mcp/` — MCP server
- `internal/web/` — HTML Kanban
- `internal/adapter/solvent/` — HTTP client to Solvent
- `cmd/conductor/` — standalone binary

Work semantics inventory:

| Semantic | Used by POC? | Disposition |
|----------|-------------|-------------|
| Task lifecycle (proposed→active→review→accepted/cancelled) | Yes | Preserve |
| Dependencies (blocking) | Minimal | Simplify — keep type, defer propagation |
| Blocked/ready behavior | Minimal | Simplify — basic blocking only |
| Cancellation | Yes | Preserve |
| Dead-end semantics (cancelled → linked to retracted belief) | Yes | Preserve — cross-system link via governance_ref |
| Lineage / REOPEN | Yes | Preserve — creates new task linked to old |
| Claim/release | Yes (agent claims work) | Preserve |
| Activity/history | Yes | Preserve |
| Priority | Yes | Preserve |
| Project container | Yes | Preserve |

Nothing is discarded. Conductor's domain model is small and entirely load-bearing.

## Oracle — library, no standalone binary

Root module: `github.com/PithomLabs/oracle`, Go 1.25, zero external dependencies.

Components:
- `coordinator/` — Go library (compile, validate, persist, context, human decisions). HTTP handler is thin routing. No standalone binary exists.
- `mcp/adapter/` — 2 tools (get_context, submit_packet). Library, no binary.
- `packet/v1/` — pure Go library (EBP packets, 10-step validation, reference resolution).
- `domain-pack/` — Go library + JSON (BM-IST: 6 debt items, retirement rules, evidence classes).
- `verifier/` — Go library (physics verification, artifact registry).
- `corpus/` — Go library + JSON (artifact manifests).
- `trust-ui/` — SEPARATE Go module, standalone binary on :8081. Own go.mod.
- `reference-loop/` — SEPARATE Go module, integration test harness. Depends on solvent + conductor via replace directives.

## Database topology

Current:
- CockroachDB — used by Solvent (20+ tables via 10 migrations)
- SQLite — used by Conductor (4 tables via 2 migrations), embedded in conductor.db file

Target:
- One CockroachDB instance, one database
- Solvent tables: unchanged (Solvent owns canonical migrations)
- Conductor tables: migrate to CockroachDB with `conductor_` prefix (ARGUS owns work-only migrations)
- One connection string

==================================================
PHASE 0 — RECONNAISSANCE + GO/NO-GO
==================================================

Goal: Determine whether the pivot is feasible without losing load-bearing semantics.

Output: GO / NO-GO / REVISE decision with evidence.

## 0.1 Complexity classification

Done above. Majority of pain points are deployment complexity. Pivot justified.

## 0.2 Work-semantics inventory

Done above. All Conductor semantics are load-bearing and small. Nothing discarded.

## 0.3 Solvent kernel reuse analysis

The kernel can be imported as a Go module from its own repository. The `kernel.Store` takes `*sql.DB` directly. No HTTP coupling in the kernel or service layer. The `api/` package is the only HTTP-dependent code and is fully discardable.

The `service/ledger.Service` already exposes most methods ARGUS needs. The gap is the read-only projection queries (`internal/view.GetSnapshot` and `internal/view.ExplainSnapshot`), which are currently internal.

Decision: import Solvent as a Go module. Solvent exposes a narrow public projection contract for ARGUS (see §1.2a). This preserves the standalone-product option and enforces the kernel boundary via module separation.

## 0.4 Component disposition table

| Component | Disposition | Rationale |
|-----------|-------------|-----------|
| Solvent kernel (`kernel/`) | KEEP — import as module | Core epistemic model, reusable in-process |
| Solvent services (`service/`) | KEEP — import as module | Thin wrappers over kernel, take `*sql.DB` |
| Solvent view projections | KEEP — public contract needed | GetSnapshot, ExplainSnapshot must be public |
| Solvent API (`api/`) | DELETE | HTTP handlers, replaced by in-process calls |
| Solvent wizard (`internal/wizard/`) | DELETE | HTTP-based interactive wizard |
| Solvent cmd binaries | DELETE | Standalone server binaries |
| Conductor domain types | KEEP — extract to `internal/work/` | 4 small types, all load-bearing |
| Conductor store | CONVERT — rewrite for CockroachDB | SQLite→CRDB, same semantics |
| Conductor HTTP/MCP/Web | DELETE | Transport layers only |
| Conductor cmd binary | DELETE | Standalone binary |
| Coordinator library | CONVERT — `internal/application/` | Orchestration logic, no HTTP |
| Coordinator HTTP handler | DELETE | Thin routing, replaced by argus serve |
| MCP adapter | CONVERT — `internal/mcp/` | Library, wraps application layer |
| Trust UI | CONVERT — embed in argus process | Same templates, direct function calls |
| Verifier | KEEP — library, same location | Already clean |
| Packet v1 | KEEP — library, same location | Already clean |
| Domain pack | KEEP — library, same location | Already clean |
| Corpus | KEEP — library, same location | Already clean |
| Reference loop | DELETE | Integration test, replaced by in-process tests |

## 0.5 Boundary enforcement design

Import-direction rules:

```
epistemic  → may NOT import: application, work, api, mcp, ui, domainpack
work       → may NOT import: epistemic (may read via application layer)
verifier   → may NOT import: epistemic authority operations
mcp        → may NOT import: decision/discharge/promote operations directly
```

Enforcement:
1. `depguard` rules in `.golangci.yml` for dev feedback
2. Compile-time negative test in `internal/boundary/boundary_test.go` that fails if forbidden imports exist

## 0.6 Auth model design

Two caller classes, one process:

| Caller | Entry point | Operations allowed |
|--------|-------------|-------------------|
| Agent (MCP) | `argus mcp` | get_context, submit_packet |
| Human (HTTP) | `argus serve` | discharge, promote, retract, reopen, authorize, plus read-only views |

Credential mechanism:

```
single configured operator credential
        ↓
authenticated server-side principal (from env/config, not request body)
        ↓
HttpOnly session cookie / minimal CSRF token
        ↓
consequential browser action
```

Enforcement:
- MCP mode: only exposes 2 tools. No HTTP server started.
- HTTP mode: `/api/retire`, `/api/promote`, `/api/retract` require Origin/Host header check + confirmation token. Attribution derived from server-side principal, never from request body.
- No full identity system. Minimal CSRF for browser actions.

## 0.7 Domain Pack contract

The current 3-method interface is insufficient. The pack contract must be typed and data-driven:

```go
type Pack interface {
    ID() string
    Version() string
    Definition() PackDefinition
    Validate() error
}

type PackDefinition struct {
    ClaimTypes      []ClaimType
    DebtVocabulary  []DebtItem
    EvidenceClasses []EvidenceClass
    RetirementRules []RetirementRule
    Falsifiers      []Falsifier
    Verifiers       []VerifierSpec
}
```

The core knows:
- debt exists
- evidence exists
- rules can be evaluated

The core does NOT know:
- needMap, needInvariant, BM-IST claim semantics
- specific evidence class names
- specific debt item names

The pack supplies those meanings.

## 0.8 NO-GO / REVISE / GO triggers

NO-GO if:
- Solvent cannot be reused without importing domain/application concerns into the kernel.
- Required Conductor semantics cannot be preserved without recreating the old service architecture.
- Human/agent authority separation cannot be mechanically preserved.
- The pivot materially weakens a thesis invariant.

REVISE if:
- Semantics can be preserved, but only with a bounded Solvent module change (e.g., adding public projection methods).
- Schema migration or transaction boundaries require a design change.
- The Domain Pack contract needs adjustment.

GO only when:
- load-bearing semantics are identified,
- their replacements are explicit,
- and no thesis invariant is weakened.

## 0.9 GO / NO-GO / REVISE decision

Questions answered:

1. Can the Solvent kernel be imported without pulling in HTTP dependencies? YES — kernel and services have zero HTTP imports.
2. Can Conductor's domain model be extracted with minimal changes? YES — 4 types, 4 tables, pure Go.
3. Are all Conductor semantics load-bearing? YES — inventory confirms nothing can be discarded.
4. Can the Coordinator logic become in-process function calls? YES — it already is a library; HTTP is thin routing.
5. Can the Trust UI embed without its own module? YES — templates + one function call to application layer.
6. Is the Solvent module boundary solvable? YES — narrow public projection contract (GetSnapshot, ExplainSnapshot).
7. Can agent/human authority be separated? YES — MCP mode vs HTTP mode, server-side attribution.

**Decision: GO.**

The pivot can proceed. All load-bearing semantics are preserved in-process. The Solvent module boundary requires a bounded public API addition (not a redesign). All other blockers are resolvable within the plan.

==================================================
PHASE 1 — CORE COLLAPSE
==================================================

Goal: Create the single-process binary with all subsystems as internal packages.

## 1.1 Create `cmd/argus/main.go`

Subcommands:
- `argus serve` — starts HTTP server with Trust UI, API endpoints, decision handlers
- `argus mcp` — starts MCP stdio mode (2 tools only)
- `argus verify` — runs physics verifier as CLI
- `argus reset` — drops and recreates schema (development only)

## 1.2 Extract Solvent kernel → `internal/epistemic/`

Import from `github.com/PithomLabs/solvent` module.

Files to create:
- `internal/epistemic/store.go` — re-export `kernel.Store` constructor
- `internal/epistemic/ledger.go` — re-export `service/ledger.Service` methods
- `internal/epistemic/authority.go` — re-export `service/authority.Service`
- `internal/epistemic/audit.go` — re-export `service/audit.Service`
- `internal/epistemic/policy.go` — re-export `service/policy.Service`

These are thin wrappers. The actual code lives in the Solvent module.

## 1.2a Solvent public projection contract

Solvent must expose `GetSnapshot` and `ExplainSnapshot` as public APIs. The recommended approach: add methods to `service/ledger.Service` that delegate to the internal view functions.

Solvent public:
    ledger.Service.GetSnapshot(ctx, scenarioID, opts) (*view.Snapshot, error)
    ledger.Service.ExplainSnapshot(scenario, scenarioID, snap) (*view.ExplainResult, error)
    view.Snapshot, view.SnapshotOpts, view.Belief, view.Evidence, view.Intent types
    view.ExplainResult, view.BeliefExplain types

Solvent private:
    internal/view (implementation stays internal)
    internal/belief (pipeline wiring stays internal)
    internal/pipeline (fixture processing stays internal)

ARGUS imports only the public surface. The internal implementation details remain Solvent's private concern.

## 1.2b Extract Domain Pack contract

Create `internal/domainpack/pack.go` in ARGUS with the Pack interface and PackDefinition types. The BM-IST pack in `domain-pack/bmist/` implements this interface.

```go
type Pack interface {
    ID() string
    Version() string
    Definition() PackDefinition
    Validate() error
}

type PackDefinition struct {
    ClaimTypes      []ClaimType       // e.g., "derived", "accommodated", "postulated"
    DebtVocabulary  []DebtItem        // e.g., "needMap", "needInvariant", ...
    EvidenceClasses []EvidenceClass   // e.g., "reproducible_artifact", "operator_asserted"
    RetirementRules []RetirementRule  // maps debt items to evidence classes + rules
    Falsifiers      []Falsifier       // what would kill a claim
    Verifiers       []VerifierSpec    // allowed verifier identities
}
```

The core consumes PackDefinition generically. Domain-specific vocabulary enters only through the pack.

## 1.3 Extract Conductor domain types → `internal/work/`

Files to create:
- `internal/work/task.go` — Task struct + state machine (copy from `conductor/internal/domain/task.go`)
- `internal/work/project.go` — Project struct
- `internal/work/dependency.go` — Dependency struct
- `internal/work/activity.go` — Activity struct
- `internal/work/store.go` — CockroachDB repository (rewrite SQLite→CRDB, see §1.5)

## 1.4 Extract Coordinator → `internal/application/`

Files to create:
- `internal/application/app.go` — App struct (replaces Coordinator struct, holds epistemic + work references directly)
- `internal/application/compile.go` — packet compilation (from `coordinator/compiler.go`)
- `internal/application/validate.go` — validation rules (from `coordinator/validate.go`)
- `internal/application/persist.go` — submit packet (from `coordinator/persist.go`)
- `internal/application/context.go` — RCP context assembly (from `coordinator/context.go`)
- `internal/application/human.go` — decision handling (from `coordinator/human.go`)
- `internal/application/idempotency.go` — DB-enforced dedup (replaces in-memory cache, see §1.4a)

Key change: all REST client calls (SolventClient, ConductorClient) become direct function calls to `internal/epistemic` and `internal/work`.

## 1.4a DB-enforced idempotency

Replace the in-memory `IdempotencyCache` with a database-enforced identity:

```sql
CREATE TABLE IF NOT EXISTS submission_idempotency (
    content_hash STRING NOT NULL,
    scenario_id  STRING NOT NULL,
    packet_id    STRING NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(content_hash, scenario_id)
);
```

The deterministic write order:
1. Compute content hash from packet (deterministic, includes scenario_id)
2. INSERT into submission_idempotency (fails on duplicate = already submitted)
3. Persist beliefs, evidence, edges, tasks
4. On any failure after step 2, the idempotency row remains — caller must retry with same packet_id

This remains correct across restart (unlike in-memory cache).

## 1.5 CockroachDB schema

### Migration ownership

Solvent owns all epistemic/authority schema migrations. ARGUS consumes Solvent's migration bundle at build time (e.g., via `go:embed` of Solvent's migration SQL files).

ARGUS owns only work-specific schema additions (conductor tables).

One database. Tables organized by source:

```
-- Epistemic (Solvent-owned, applied from Solvent migration bundle)
belief, belief_edge, evidence, action_intent,
principal, authority_target, target_snapshot, target_activation,
target_revocation, justification, debt_discharge,
audit_activity

-- Work (ARGUS-owned, new migration)
conductor_project, conductor_task, conductor_activity, conductor_dependency

-- Idempotency (ARGUS-owned, new migration)
submission_idempotency
```

### Conductor→CRDB type mapping

| SQLite Column | SQLite Constraint | CRDB Column | CRDB Constraint |
|--------------|-------------------|-------------|-----------------|
| `id TEXT` | `PRIMARY KEY` | `id STRING` | `PRIMARY KEY` |
| `name TEXT` | `NOT NULL` | `name STRING` | `NOT NULL` |
| `description TEXT` | nullable | `description STRING` | nullable |
| `status TEXT` | `NOT NULL DEFAULT '...' CHECK (...)` | `status STRING` | `NOT NULL DEFAULT '...' CHECK (...)` |
| `priority TEXT` | `NOT NULL DEFAULT 'medium' CHECK (...)` | `priority STRING` | `NOT NULL DEFAULT 'medium' CHECK (...)` |
| `current_agent TEXT` | nullable | `current_agent STRING` | nullable |
| `governance_ref TEXT` | nullable | `governance_ref STRING` | nullable |
| `created_at TEXT` | `NOT NULL DEFAULT (datetime('now'))` | `created_at TIMESTAMPTZ` | `NOT NULL DEFAULT now()` |
| `updated_at TEXT` | `NOT NULL DEFAULT (datetime('now'))` | `updated_at TIMESTAMPTZ` | `NOT NULL DEFAULT now()` |
| `actor_type TEXT` | `NOT NULL CHECK (...)` | `actor_type STRING` | `NOT NULL CHECK (...)` |
| `actor_id TEXT` | `NOT NULL` | `actor_id STRING` | `NOT NULL` |
| `action TEXT` | `NOT NULL` | `action STRING` | `NOT NULL` |
| `details TEXT` | nullable (JSON) | `details STRING` | nullable (JSON) |

### SQLite function mapping

| SQLite | CRDB | Where |
|--------|------|-------|
| `datetime('now')` | `now()` | Schema DEFAULT clauses, inline SQL |
| `?` placeholder | `$N` positional | All queries |

### CRDB retry requirement

The Conductor `Transaction()` wrapper must add CRDB serialization retry:

```go
func (db *DB) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
    for {
        tx, err := db.BeginTx(ctx, nil)
        if err != nil {
            return err
        }
        if err := fn(tx); err != nil {
            tx.Rollback()
            return err
        }
        if err := tx.Commit(); err != nil {
            if isRetryable(err) { // SQLSTATE 40001
                continue
            }
            return err
        }
        return nil
    }
}
```

This affects: `Transition()`, `TransitionWithAgent()`, `Claim()`.

### Pre-existing bug: Release() is not transactional

The `Release()` UPDATE and subsequent activity INSERT are separate statements. Wrap in a transaction.

### Timestamp handling

Go domain types store timestamps as `string`. Under CRDB with `TIMESTAMPTZ`, the `pgx` driver formats as RFC 3339 (`"2006-01-02T15:04:05Z07:00"`). The Go-side generation uses `time.RFC3339` too. Both are format-consistent. The test at `store_test.go:542` uses a non-RFC format and must be updated.

### Proposed CRDB migration (011_work_schema.sql)

```sql
CREATE TABLE IF NOT EXISTS conductor_project (
    id          STRING PRIMARY KEY,
    name        STRING NOT NULL,
    description STRING,
    status      STRING NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'completed', 'archived')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conductor_task (
    id              STRING PRIMARY KEY,
    project_id      STRING NOT NULL REFERENCES conductor_project(id),
    title           STRING NOT NULL,
    description     STRING,
    status          STRING NOT NULL DEFAULT 'proposed'
                    CHECK (status IN ('proposed','active','review','accepted','blocked','cancelled')),
    priority        STRING NOT NULL DEFAULT 'medium'
                    CHECK (priority IN ('low','medium','high','critical')),
    current_agent   STRING,
    governance_ref  STRING,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_task_project ON conductor_task(project_id);
CREATE INDEX IF NOT EXISTS idx_task_status ON conductor_task(status);
CREATE INDEX IF NOT EXISTS idx_task_agent ON conductor_task(current_agent);

CREATE TABLE IF NOT EXISTS conductor_activity (
    id          STRING PRIMARY KEY,
    task_id     STRING NOT NULL REFERENCES conductor_task(id),
    actor_type  STRING NOT NULL CHECK (actor_type IN ('human','agent','system')),
    actor_id    STRING NOT NULL,
    action      STRING NOT NULL,
    details     STRING,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_activity_task ON conductor_activity(task_id);
CREATE INDEX IF NOT EXISTS idx_activity_created ON conductor_activity(task_id, created_at DESC);

CREATE TABLE IF NOT EXISTS conductor_dependency (
    id              STRING PRIMARY KEY,
    task_id         STRING NOT NULL REFERENCES conductor_task(id),
    blocked_by_id   STRING NOT NULL REFERENCES conductor_task(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(task_id, blocked_by_id),
    CHECK (task_id != blocked_by_id)
);

CREATE INDEX IF NOT EXISTS idx_dependency_task ON conductor_dependency(task_id);
CREATE INDEX IF NOT EXISTS idx_dependency_blocked_by ON conductor_dependency(blocked_by_id);

CREATE TABLE IF NOT EXISTS submission_idempotency (
    content_hash STRING NOT NULL,
    scenario_id  STRING NOT NULL,
    packet_id    STRING NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(content_hash, scenario_id)
);
```

## 1.6 Package boundary enforcement

Create:
- `.golangci.yml` with depguard rules
- `internal/boundary/boundary_test.go` with compile-time negative tests

## 1.7 Auth boundary

- MCP mode (`internal/mcp/`): only exposes get_context and submit_packet. No import of decision/discharge packages.
- HTTP mode (`internal/ui/` or `internal/api/`): human endpoints require Origin check + confirmation token. Agent endpoints (if any HTTP) are read-only.

==================================================
PHASE 2 — INTEGRATION / END-TO-END FLOW
==================================================

Goal: Wire all components together and verify the complete EBP packet flow.

## 2.1 Wire MCP adapter

Convert `mcp/adapter/adapter.go` to use `internal/application.App` directly instead of HTTP client.

- `HandleTool("argus.get_context", ...)` → `app.GetContext(ctx, taskID)`
- `HandleTool("argus.submit_packet", ...)` → `app.SubmitPacket(ctx, packet)`

No HTTP in the path.

## 2.2 Wire Trust UI

Convert `trust-ui/server.go` from separate module to `internal/ui/server.go`.

- Embed templates via `//go:embed`
- Replace HTTP call to Coordinator with direct `app.GetContext()` call
- Replace HTTP call to Coordinator `/decisions` with direct `app.SubmitDecision()` call
- Serve from `argus serve` on `:8080`

## 2.3 Wire verifier + trust binding

Keep `verifier/` as-is. The `argus verify` subcommand calls `verifier.RunPhysicsVerifier()` directly.

For in-process verification during packet validation: call the verifier library from `internal/application/compile.go` when a packet includes verification evidence.

Verifier trust binding — the VerificationArtifact schema must include:

```go
type VerificationArtifact struct {
    RunID             string   `json:"run_id"`
    VerifierID        string   `json:"verifier_id"`        // e.g., "bm-ist-physics-v1"
    VerifierVersion   string   `json:"verifier_version"`   // e.g., "0.1.0"
    VerificationInput string   `json:"verification_input_hash"` // SHA-256 of inputs
    ArtifactHash      string   `json:"artifact_hash"`      // SHA-256 of this artifact
    Steps             []Step   `json:"steps"`
    Result            string   `json:"result"`             // confirmed/refuted/inconclusive
    EvidenceRef       string   `json:"evidence_ref,omitempty"`
}
```

The Domain Pack selects the allowed verifier identity via `VerifierSpec`. The agent must not choose arbitrary verification logic. `argus verify` and in-process verification use the same verifier library and produce the same artifact schema.

## 2.4 EBP packet submission flow (revised)

Agent packets must NEVER perform debt discharge, promotion, retraction, or authorization.

```
Agent → MCP → app.SubmitPacket(packet)
  → compile(packet)              [validate, resolve refs, compute hash]
  → persist(packet)
     → epistemic.EnterBelief     [Solvent kernel — own transaction]
     → epistemic.AddEvidence     [Solvent kernel — own transaction]
     → work.CreateTask           [Work store — own transaction]
     → audit.Log                 [Solvent audit]
  → return result
```

Transaction model: Solvent authority mutations remain transactionally atomic within Solvent. Cross-subsystem packet persistence is not globally atomic in the POC. Partial failures are explicit, retryable, and idempotent.

What submit_packet does NOT do:
- RetireDebt (human decision only)
- Promote (human decision only)
- Retract (human decision only)
- Authorize (human decision only)

Proposed retirement: the packet can declare a `proposed_retirement` field indicating which debt items the agent believes should be discharged. This is recorded as metadata for human review, not executed.

## 2.4a Agent-authored edge behavior

Agent-created contradicts/derives edges enter as epistemic proposals/evidence, NOT as authority transitions.

```
Agent packet edge
    → stored as belief_edge row
    → visible in UI as "proposed by agent ⟨id⟩"
    → does NOT trigger retraction
    → does NOT block promotion
    → human RETRACT remains the only mechanism for authority state change
```

Only introduce PROPOSED/ACTIVE edge persistence if repository inspection shows an agent-created edge can currently trigger a consequential transition without it. For this POC, the existing edge model should be sufficient if the application layer does not act on agent edges without human confirmation.

## 2.5 Context assembly flow

```
Agent → MCP → app.GetContext(taskID)
  → work.GetTask(taskID)              [CRDB read]
  → work.ListDependencies(taskID)     [CRDB read]
  → epistemic.GetSnapshot(scenario)   [Solvent view, CRDB read]
  → epistemic.GetActivity(scenario)   [Solvent audit, CRDB read]
  → assemble RCP context
  → return
```

Single-DB reads. No cross-service projection.

## 2.6 Human decision flow

```
Human → Browser → UI → app.SubmitDecision(type, params)
  → validate retirement rule  [domain pack]
  → epistemic.Discharge(...)   [if type=RETIRE_DEBT — own Solvent transaction]
  → epistemic.Promote(...)     [if type=PROMOTE — own Solvent transaction]
  → epistemic.Retract(...)     [if type=RETRACT — own Solvent transaction]
  → work.CancelTask(...)       [if retraction causes dead end — own work transaction]
  → audit.Log(...)
```

All in-process. Attribution from server-side principal, never from request body.

==================================================
PHASE 3 — VERIFICATION, PORTABILITY PROOF + CLEANUP
==================================================

Goal: Prove the pivot works, prove domain portability, delete old code.

## 3.1 End-to-end demo walkthrough

Execute the full 17-step narrative:

1. `argus serve` starts
2. OpenCode connects via MCP
3. Agent calls get_context
4. Agent performs bounded work
5. Agent calls submit_packet
6. Beliefs, evidence, debt recorded
7. Debt visible in UI
8. `argus verify` produces verification artifact
9. Adversarial agent submits contradiction packet
10. Contradiction visible in UI
11. Human reviews in Trust UI
12. Human discharges debt (via UI)
13. Claim promoted (via Solvent kernel gate)
14. Human retracts (via UI)
15. Linked task cancelled (dead end)
16. REOPEN creates new lineage
17. UI shows full history

## 3.2 Refusal test suite (negative tests)

At minimum:

```
agent → discharge           REFUSED
agent → promote             REFUSED
open debt → promote         REFUSED
retracted belief → promote  REFUSED
unconfirmed contradiction → retract REFUSED
bad Origin/Host → REFUSED
body-supplied principal → IGNORED
failed mutation → no success audit event
```

These are acceptance criteria, not nice-to-haves. Refusal is part of the thesis.

## 3.3 BM-IST import isolation test

Rename from "null domain-pack CI compile test" to what it actually proves: BM-IST import isolation.

Create `internal/domainpack/isolation_test.go`:

```go
// Build-tagged test that compiles the core against a minimal stub domain pack.
// Proves core does not import BM-IST vocabulary.
func TestBMISTImportIsolation(t *testing.T) {
    // stub pack satisfies the Pack interface with empty implementations
    // core packages are imported — if they reference BM-IST types, compilation fails
}
```

Add a small behavioral smoke test for the generic pack contract:

```go
func TestGenericPackContract(t *testing.T) {
    // verify that a minimal pack can be registered, queried, and its definition used
    // by the core without any domain-specific knowledge
}
```

## 3.4 Simplification ledger verification

| Dimension | Before | After | Verified? |
|-----------|--------|-------|-----------|
| Binaries | ~5 | 1 | |
| Process instances | ~7 (+ OpenCode) | 2 (argus serve + argus mcp) + OpenCode | |
| Ports | ~6 | 1 (HTTP+UI; MCP is stdio) | |
| Databases | 1-2 (CockroachDB + SQLite) | 1 (CockroachDB) | |
| Go modules | 3 (oracle, trust-ui, reference-loop) | 1 (+ Solvent import) | |
| Credential sets | multiple | 1 | |
| Internal HTTP calls | many | 0 | |
| Cross-process handoffs | yes | none | |

## 3.5 Delete old code + tests

Delete:
- `coordinator/http/` — HTTP handler (replaced by internal application calls)
- `coordinator/client.go` — REST clients (replaced by direct function calls)
- `coordinator/mock.go` — may need rewrite for in-process mocks
- `trust-ui/` as separate module — embedded in main binary
- `reference-loop/` — replaced by in-process integration tests
- `trust-ui/go.mod` — no longer separate module
- HTTP handler tests for old service topology
- Tests that encode deployment mechanics

Preserve:
- `coordinator/compiler.go` logic → moved to `internal/application/compile.go`
- `coordinator/validate.go` logic → moved to `internal/application/validate.go`
- `coordinator/persist.go` logic → moved to `internal/application/persist.go`
- `coordinator/context.go` logic → moved to `internal/application/context.go`
- `coordinator/human.go` logic → moved to `internal/application/human.go`
- All packet validation tests
- All domain pack validation tests
- All verifier tests

## 3.6 Taskfile.yml

```yaml
version: '3'
tasks:
  dev:
    desc: Start ARGUS in development mode
    cmds:
      - |
        # Start CockroachDB with health polling
        cockroach demo --no-example-database --listen-addr :26260 &
        COCKROACH_PID=$!
        until cockroach node status --host :26260 --insecure 2>/dev/null; do
          sleep 1
        done
        # Reset schema
        go run ./cmd/argus reset --db postgres://root@localhost:26260/argus?sslmode=disable
        # Start ARGUS
        go run ./cmd/argus serve --db postgres://root@localhost:26260/argus?sslmode=disable

  test:
    desc: Run all tests
    cmds:
      - go test ./...
      - go test -race ./...
      - go vet ./...

  verify:
    desc: Run physics verifier
    cmds:
      - go run ./cmd/argus verify

  lint:
    desc: Run linter with boundary checks
    cmds:
      - golangci-lint run
```

==================================================
MIGRATION MODEL
==================================================

- In-place pivot on a feature branch.
- Rollback: if Phase 0 gate is REVISE, abandon branch. If Phase 1+ fails, revert to main.
- No live state to protect (POC, fresh DB each run).
- Solvent module imported via `go.mod` replace directive during development, switched to released version for merge.

==================================================
TEST DISPOSAL
==================================================

| Category | Tests | Action |
|----------|-------|--------|
| PRESERVE | `coordinator/*_test.go` (compiler, idempotency, retirement rules) | Move to `internal/application/` |
| PRESERVE | `packet/v1/*_test.go` (validation, resolution) | Keep as-is |
| PRESERVE | `domain-pack/*_test.go` (pack validation) | Keep as-is |
| PRESERVE | `verifier/*_test.go` (physics verification) | Keep as-is |
| REWRITE | `coordinator/http/handler_test.go` | Rewrite as `internal/application/` unit tests |
| REWRITE | `mcp/adapter/adapter_test.go` | Rewrite to use in-process App |
| DELETE | `trust-ui/` tests (separate module) | Embedded, tested via integration |
| DELETE | `reference-loop/*_test.go` | Replaced by in-process integration tests |
| DELETE | Solvent API handler tests (in solvent repo, not imported) | N/A |
| DELETE | Conductor HTTP/MCP/Web tests (in conductor repo, not imported) | N/A |
| ADD | `internal/boundary/boundary_test.go` | Import-direction enforcement |
| ADD | `internal/domainpack/isolation_test.go` | BM-IST import isolation proof |
| ADD | `internal/domainpack/contract_test.go` | Generic pack contract smoke test |
| ADD | Refusal test suite | 8 negative test cases |
| ADD | End-to-end integration test | Full 17-step narrative |

==================================================
RISKS AND NON-GOALS
==================================================

Risks:
- Solvent module import may pull unexpected transitive dependencies (mitigated: kernel/services have zero HTTP imports)
- CockroachDB migration from SQLite requires type/function adaptation (mitigated: explicit mapping in §1.5)
- Solvent public projection contract requires Solvent repo change (mitigated: narrow, backward-compatible addition)
- Separate-transaction model means partial failures are possible (mitigated: idempotent writes, explicit failure states)
- Trust UI template embedding may require path adjustments (mitigated: standard go:embed)

Non-goals:
- Production-grade auth system
- Container-per-agent isolation
- Elaborate verifier attestation
- Second domain pack
- Distributed deployment support
- Performance optimization
- Global atomicity across subsystems

==================================================
ACCEPTANCE CRITERIA
==================================================

1. `argus serve` starts and serves Trust UI on :8080
2. `argus mcp` starts and exposes exactly 2 tools
3. `argus verify` runs physics verifier and produces artifact with verifier_id, version, input_hash
4. Agent can submit EBP packet via MCP → beliefs + evidence + debt recorded in CRDB
5. Agent submit_packet does NOT perform RetireDebt, Promote, Retract, or Authorize
6. Human can discharge debt via Trust UI → debt retired in CRDB (attribution from server-side principal)
7. Human can promote claim → belief promoted in CRDB (only if debt empty)
8. Adversarial packet introduces contradiction → visible in UI
9. Retraction causes linked task cancellation (dead end)
10. Agent-created edges do NOT trigger retraction or block promotion
11. `go test ./...` passes
12. `golangci-lint run` passes (including depguard rules)
13. BM-IST import isolation test passes (core compiles against stub pack)
14. Simplification ledger verified (0 internal HTTP calls, 2 instances, 1 DB)
15. Refusal test suite passes (8 negative test cases)
16. Idempotency: duplicate packet submission returns same result, no duplicate state
17. Solvent owns canonical migrations; ARGUS owns only work migrations
