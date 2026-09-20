# ARGUS POC Simplification Pivot — Implementation Plan (Final)

Final revision per three rounds of adversarial review. Incorporates 12 amendments + 10 concrete fixes.

Status: PLAN ONLY — no implementation code.

This is the last plan revision. The remaining defects are implementation semantics, not architectural problems.

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

The pivot is reconnaissance-gated. Phase 0 produces a GO / NO-GO / REVISE decision before any code changes.

Solvent change budget: 2 bounded backward-compatible additions
  1. Public projection contract (GetSnapshot, ExplainSnapshot on ledger.Service)
  2. Migration export FS (embed.FS or Apply function)

Both are SemVer-minor-class changes. No API redesign required.

==================================================
COMPLEXITY CLASSIFICATION
==================================================

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
- `service/authority.Service` — execution boundary
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

State machine: atomic compare-and-swap in SQL.

Storage: SQLite via `modernc.org/sqlite`. 4 tables.

Work semantics inventory:

| Semantic | Used by POC? | Disposition |
|----------|-------------|-------------|
| Task lifecycle | Yes | Preserve |
| Dependencies | Minimal | Simplify — keep type, defer propagation |
| Blocked/ready | Minimal | Simplify — basic blocking only |
| Cancellation | Yes | Preserve |
| Dead-end semantics | Yes | Preserve — governance_ref FK to belief |
| Lineage / REOPEN | Yes | Preserve — reopened_from_task_id column |
| Claim/release | Yes | Preserve |
| Activity/history | Yes | Preserve |
| Priority | Yes | Preserve |
| Project container | Yes | Preserve |

Nothing is discarded. Conductor's domain model is small and entirely load-bearing.

## Oracle — library, no standalone binary

Root module: `github.com/PithomLabs/oracle`, Go 1.25, zero external dependencies.

Components:
- `coordinator/` — Go library. No standalone binary.
- `mcp/adapter/` — 2 tools. Library, no binary.
- `packet/v1/` — pure Go library.
- `domain-pack/` — Go library + JSON.
- `verifier/` — Go library.
- `corpus/` — Go library + JSON.
- `trust-ui/` — SEPARATE Go module, standalone binary.
- `reference-loop/` — SEPARATE Go module, integration test harness.

## Database topology

Current: CockroachDB (Solvent) + SQLite (Conductor).
Target: One CockroachDB instance, one database. Solvent tables unchanged. Conductor tables migrated to CRDB with `conductor_` prefix.

==================================================
PHASE 0 — RECONNAISSANCE + GO/NO-GO
==================================================

Goal: Determine whether the pivot is feasible without losing load-bearing semantics.

## 0.1 Complexity classification

Done above. Pivot justified.

## 0.2 Work-semantics inventory

Done above. All Conductor semantics are load-bearing and small.

## 0.3 Solvent kernel reuse analysis

The kernel can be imported as a Go module. `kernel.Store` takes `*sql.DB` directly. No HTTP coupling.

The `service/ledger.Service` exposes most methods ARGUS needs. The gap is the read-only projections (`internal/view.GetSnapshot` and `internal/view.ExplainSnapshot`), which are currently internal.

Solvent change budget (2 bounded additions):
  1. Public projection contract on ledger.Service (GetSnapshot, ExplainSnapshot + types)
  2. Migration export package (embed.FS or Apply function)

Both are backward-compatible, SemVer-minor additions. No API redesign.

Decision: import Solvent as a Go module.

## 0.4 Component disposition table

| Component | Disposition | Rationale |
|-----------|-------------|-----------|
| Solvent kernel (`kernel/`) | KEEP — import as module | Core epistemic model |
| Solvent services (`service/`) | KEEP — import as module | Thin wrappers, take `*sql.DB` |
| Solvent view projections | KEEP — public contract needed | GetSnapshot, ExplainSnapshot |
| Solvent API (`api/`) | DELETE | HTTP handlers |
| Solvent wizard (`internal/wizard/`) | DELETE | HTTP wizard |
| Solvent cmd binaries | DELETE | Standalone binaries |
| Conductor domain types | KEEP — extract to `internal/work/` | All load-bearing |
| Conductor store | CONVERT — rewrite for CRDB | SQLite→CRDB |
| Conductor HTTP/MCP/Web | DELETE | Transport layers |
| Conductor cmd binary | DELETE | Standalone binary |
| Coordinator library | CONVERT — `internal/application/` | Orchestration logic |
| Coordinator HTTP handler | DELETE | Thin routing |
| MCP adapter | CONVERT — `internal/mcp/` | Library |
| Trust UI | CONVERT — embed in argus process | Templates + direct calls |
| Verifier | KEEP — library | Already clean |
| Packet v1 | KEEP — library | Already clean |
| Domain pack | KEEP — library | Already clean |
| Corpus | KEEP — library | Already clean |
| Reference loop | DELETE | Replaced by in-process tests |

## 0.5 Boundary enforcement design

Import-direction rules:

```
epistemic     → may NOT import: application, work, api, mcp, ui, domainpack
work          → may NOT import: epistemic (may read via application layer)
verifier      → may NOT import: epistemic authority operations
mcp           → may NOT import: decision/discharge/promote/retract operations
application   → may import: epistemic, work, domainpack (interface only), verifier, packet
application   → may NOT import: mcp, ui
```

Enforcement:
1. `depguard` rules in `.golangci.yml`
2. Compile-time negative test in `internal/boundary/boundary_test.go`
3. Application boundary row in matrix

## 0.6 Auth model design

Two caller classes, one process:

| Caller | Entry point | Operations allowed |
|--------|-------------|-------------------|
| Agent (MCP) | `argus mcp` | get_context, submit_packet |
| Human (HTTP) | `argus serve` | discharge, promote, retract, reopen, authorize, read-only views |

No agent HTTP surface exists in the POC. Agents are MCP-only.

Credential mechanism:

```
single configured operator token (env/config, static secret)
        ↓
GET /login validates token
        ↓
HttpOnly session cookie (HMAC-signed, in-memory session store)
        ↓
consequential POST requires valid session
        ↓
CSRF check (Origin/Host header + confirmation token)
        ↓
server-derived principal (never from request body)
```

No users table, OAuth, identity provider, or RBAC. `discharged_by` never comes from the request body.

Session lifecycle:
- Token validated at login, session cookie issued
- Session expires after configurable timeout (default: 1 hour)
- Logout clears session
- Credential rotation: change env var, restart (POC only)

## 0.7 Domain Pack contract

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

Field schemas:

```go
type ClaimType string
// Opaque classification label. Core uses as string, no semantic evaluation.
// BM-IST: "derived", "accommodated", "postulated"

type DebtItem string
// Opaque string identifier for a domain-specific debt obligation.
// BM-IST: "needMap", "needInvariant", "needToyCheck",
//         "needNullModel", "needObstruction", "needFaithfulnessReview"

type EvidenceClass string
// Opaque identifier for a class of evidence.
// BM-IST: "reproducible_artifact", "operator_asserted"

type RetirementRule struct {
    DebtItem      DebtItem      `json:"debt_item"`
    RequiredClass EvidenceClass `json:"required_class"`
    Rule          string        `json:"rule"`           // human-readable
    RequiresHuman bool          `json:"requires_human"` // human discharge required
}

type Falsifier struct {
    ID          string `json:"id"`
    Description string `json:"description"`
}

type VerifierSpec struct {
    VerifierID string `json:"verifier_id"`
    MinVersion string `json:"min_version"` // semver
}
```

BM-IST retirement rules (one per debt item):

| Debt Item | Required Class | Rule | Human |
|-----------|---------------|------|-------|
| needMap | reproducible_artifact | Artifact must demonstrate faithful map reproduction | yes |
| needInvariant | reproducible_artifact | Artifact must verify invariant preservation | yes |
| needToyCheck | reproducible_artifact | Artifact must pass toy-model consistency check | yes |
| needNullModel | reproducible_artifact | Artifact must rule out null-model hypothesis | yes |
| needObstruction | reproducible_artifact | Artifact must demonstrate obstruction criterion | yes |
| needFaithfulnessReview | reproducible_artifact | Artifact must show faithful reproduction of reference result within declared tolerance | yes |

Pack.Validate() checks:
- Every debt item has at least one retirement rule
- Every retirement rule references a declared evidence class
- No duplicate debt items or evidence classes
- VerifierSpec versions are valid semver ranges
- At least one VerifierSpec is declared if any retirement rule RequiresHuman

The core knows: debt exists, evidence exists, rules can be evaluated.
The core does NOT know: needMap, needInvariant, BM-IST claim semantics, specific evidence class names.

## 0.8 Discharge vs RetireDebt

Two kernel methods with different semantics:

- `RetireDebt(ctx, scenario, beliefID, item)` — kernel method, retires exactly one debt item. Internal-use-only. Called by the application layer's discharge path. Not exposed via MCP or HTTP directly.
- `Discharge(ctx, scenario, beliefID, obligationKey, instrumentRef, dischargedBy)` — higher-level method that records attributed discharge with replay protection (UNIQUE constraint on belief+obligation+instrument). This is the method the human flow calls.

In ARGUS, `RetireDebt` is internal to the application layer. `Discharge` is the human-facing API. The boundary enforcement test verifies MCP cannot reach `RetireDebt`.

## 0.9 NO-GO / REVISE / GO triggers

NO-GO if:
- Solvent cannot be reused without importing domain/application concerns into the kernel.
- Required Conductor semantics cannot be preserved without recreating the old service architecture.
- Human/agent authority separation cannot be mechanically preserved.
- The pivot materially weakens a thesis invariant.

REVISE only for:
- Material scope expansion (new services, new databases, new invariants).
- Weakened thesis invariants.
- Solvent changes requiring API redesign (not backward-compatible additions).

GO when:
- Load-bearing semantics identified.
- Replacements explicit.
- All Solvent changes are backward-compatible additions (public projections + migration FS).
- No thesis invariant weakened.

## 0.10 GO / NO-GO / REVISE decision

Questions answered:

1. Solvent kernel reusable without HTTP? YES.
2. Conductor domain model extractable? YES.
3. All Conductor semantics load-bearing? YES.
4. Coordinator logic becomes in-process calls? YES.
5. Trust UI embeddable? YES.
6. Solvent module boundary solvable? YES — 2 bounded backward-compatible additions.
7. Agent/human authority separable? YES — MCP-only agents, server-side attribution.

Solvent change ledger:
  Change 1: Public projection contract (GetSnapshot, ExplainSnapshot + types on ledger.Service)
  Change 2: Migration export package (embed.FS or Apply function)
  Both: backward-compatible, SemVer-minor, no redesign.

**Decision: GO.**

The gate fired REVISE on the Solvent module boundary. Resolved by the 2 bounded additions within the approved change budget. No material scope expansion. No thesis invariant weakened. GO granted.

==================================================
PHASE 1 — CORE COLLAPSE
==================================================

Goal: Create the single-process binary with all subsystems as internal packages.

## 1.1 Create `cmd/argus/main.go`

Subcommands:
- `argus serve` — HTTP server with Trust UI, API, decision handlers
- `argus mcp` — MCP stdio mode (2 tools only)
- `argus verify` — physics verifier CLI
- `argus reset` — drops and recreates schema (development only)

## 1.2 Extract Solvent kernel → `internal/epistemic/`

Import from `github.com/PithomLabs/solvent` module.

Files to create:
- `internal/epistemic/store.go` — re-export `kernel.Store` constructor
- `internal/epistemic/ledger.go` — re-export `service/ledger.Service` methods
- `internal/epistemic/authority.go` — re-export `service/authority.Service`
- `internal/epistemic/audit.go` — re-export `service/audit.Service`
- `internal/epistemic/policy.go` — re-export `service/policy.Service`

Thin wrappers. Actual code lives in Solvent module.

## 1.2a Solvent public projection contract (Change 1)

Solvent exposes on `service/ledger.Service`:

```go
func (s *Service) GetSnapshot(ctx context.Context, scenarioID string, opts view.SnapshotOpts) (*view.Snapshot, error)
func (s *Service) ExplainSnapshot(scenario, scenarioID string, snap *view.Snapshot) *view.ExplainResult
```

Public types promoted from internal/view:
- `view.Snapshot`, `view.SnapshotOpts`
- `view.Belief`, `view.Evidence`, `view.Intent`
- `view.ExplainResult`, `view.BeliefExplain`

Stability: documented as "unstable, may change without notice" for this POC. SemVer-stable commitment deferred to post-POC.

Internal implementation stays in `internal/view/`. ARGUS imports only the public surface.

## 1.2b Migration export package (Change 2)

Solvent exports:

```go
package solventmigrations

import "embed"

//go:embed *.sql
var FS embed.FS

func Apply(db *sql.DB) error {
    // reads *.sql from FS, applies in order
}
```

ARGUS calls `solventmigrations.Apply(db)` then applies its own work/idempotency migrations.

If Solvent maintainers decline: copy ~200 lines of view logic into ARGUS's `internal/epistemic/view.go`. Accept drift risk for POC.

## 1.2c Extract Domain Pack contract

Create `internal/domainpack/pack.go` in ARGUS with the Pack interface and PackDefinition types (§0.7).

BM-IST pack in `domain-pack/bmist/` implements this interface.

## 1.3 Extract Conductor domain types → `internal/work/`

Files to create:
- `internal/work/task.go` — Task struct + state machine
- `internal/work/project.go` — Project struct
- `internal/work/dependency.go` — Dependency struct
- `internal/work/activity.go` — Activity struct
- `internal/work/store.go` — CockroachDB repository (§1.5)

## 1.4 Extract Coordinator → `internal/application/`

Files to create:
- `internal/application/app.go` — App struct (holds epistemic + work references)
- `internal/application/compile.go` — packet compilation
- `internal/application/validate.go` — validation rules
- `internal/application/persist.go` — submit packet
- `internal/application/context.go` — RCP context assembly
- `internal/application/human.go` — decision handling
- `internal/application/proposed_retirement.go` — proposed retirement metadata handling

Key change: all REST client calls become direct function calls.

## 1.4a Entity-level idempotency (replaces broken insert-first design)

The previous design inserted an idempotency row before persistence, blocking retry on partial failure. The fix moves the guarantee to entity-level uniqueness:

Write order:
1. Compute content hash from packet (deterministic, includes scenario_id)
2. Persist beliefs — each belief write uses `INSERT ... ON CONFLICT (scenario_id, claim_hash) DO NOTHING` (deterministic unique constraint per scenario)
3. Persist evidence — each evidence write uses `INSERT ... ON CONFLICT (scenario_id, content_sha256) DO NOTHING`
4. Persist edges — each edge write uses `INSERT ... ON CONFLICT (scenario_id, from_id, to_id, kind) DO NOTHING`
5. Persist tasks — each task write uses deterministic ID, `INSERT ... ON CONFLICT DO NOTHING`
6. Write `submission_idempotency` row LAST with status `COMPLETED` — its presence certifies completion

On partial failure:
- Entity writes that succeeded are naturally idempotent (ON CONFLICT DO NOTHING)
- Retry resumes from step 2 — already-written entities are no-ops
- No idempotency row written → caller knows packet is not complete
- Duplicate submission of completed packet → idempotency row exists, return cached result

Schema:

```sql
CREATE TABLE IF NOT EXISTS submission_idempotency (
    content_hash  STRING NOT NULL,
    scenario_id   STRING NOT NULL,
    packet_id     STRING NOT NULL,
    status        STRING NOT NULL DEFAULT 'completed'
                  CHECK (status IN ('completed')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(content_hash, scenario_id)
);
```

Acceptance test #16b: partial-failure retry completes the packet with no duplicate state.

## 1.5 CockroachDB schema

### Migration ownership

Solvent owns all epistemic/authority schema migrations. ARGUS consumes Solvent's migration via `solventmigrations.Apply(db)`.

ARGUS owns only work-specific and idempotency schema.

### Table layout

```
-- Epistemic (Solvent-owned, applied via solventmigrations.Apply)
belief, belief_edge, evidence, action_intent,
principal, authority_target, target_snapshot, target_activation,
target_revocation, justification, debt_discharge,
audit_activity

-- Proposed retirement (ARGUS-owned, new migration)
ALTER TABLE belief ADD COLUMN proposed_retirement JSONB;

-- Work (ARGUS-owned, new migration)
conductor_project, conductor_task, conductor_activity, conductor_dependency

-- Idempotency (ARGUS-owned, new migration)
submission_idempotency
```

### Conductor→CRDB type mapping

| SQLite Column | CRDB Column | Notes |
|--------------|-------------|-------|
| `id TEXT PRIMARY KEY` | `id STRING PRIMARY KEY` | UUID from Go |
| `name TEXT NOT NULL` | `name STRING NOT NULL` | Direct |
| `status TEXT DEFAULT '...' CHECK (...)` | `status STRING DEFAULT '...' CHECK (...)` | Direct |
| `priority TEXT DEFAULT 'medium' CHECK (...)` | `priority STRING DEFAULT 'medium' CHECK (...)` | Direct |
| `current_agent TEXT` | `current_agent STRING` | Nullable |
| `governance_ref TEXT` | `governance_ref STRING REFERENCES belief(id)` | FK upgrade |
| `created_at TEXT DEFAULT (datetime('now'))` | `created_at TIMESTAMPTZ DEFAULT now()` | Function change |
| `updated_at TEXT DEFAULT (datetime('now'))` | `updated_at TIMESTAMPTZ DEFAULT now()` | Function change |
| `actor_type TEXT CHECK (...)` | `actor_type STRING CHECK (...)` | Direct |
| `details TEXT` | `details STRING` | JSON as text |

SQLite `?` placeholders → CRDB `$N` positional.
SQLite `datetime('now')` → CRDB `now()`.

### CRDB retry wrapper

```go
func (db *DB) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
    const maxAttempts = 5
    for attempt := 0; attempt < maxAttempts; attempt++ {
        tx, err := db.BeginTx(ctx, nil)
        if err != nil {
            return err
        }
        if err := fn(tx); err != nil {
            tx.Rollback()
            if isRetryable(err) && attempt < maxAttempts-1 {
                continue
            }
            return err
        }
        if err := tx.Commit(); err != nil {
            if isRetryable(err) && attempt < maxAttempts-1 {
                continue
            }
            return err
        }
        return nil
    }
    return fmt.Errorf("transaction failed after %d attempts", maxAttempts)
}

func isRetryable(err error) bool {
    // SQLSTATE 40001 = serialization failure
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return pgErr.SQLState() == "40001"
    }
    return false
}
```

Retryable errors caught anywhere in the closure (statement or commit), not just commit.

### REOPEN lineage

Add to `conductor_task`:

```sql
reopened_from_task_id STRING REFERENCES conductor_task(id)
```

On REOPEN: create new task with `reopened_from_task_id = old_task.id`. Old task remains cancelled. New task enters proposed state.

### Dead-end task lookup

On retraction:
```sql
SELECT id FROM conductor_task WHERE governance_ref = :belief_id AND status != 'cancelled'
```
Cancel each matching task in its own transaction.

### Proposed retirement persistence

Add to `belief` table:
```sql
proposed_retirement JSONB
```

When agent submits packet with proposed retirement, store as metadata on the belief. Trust UI Debts view renders: "Agent proposes retiring ⟨debt item⟩ — awaiting human decision."

### Pre-existing bug: Release() not transactional

The UPDATE and subsequent activity INSERT are separate statements. Wrap in a transaction.

### Timestamp handling

Go domain types store timestamps as `string`. Under CRDB with `TIMESTAMPTZ`, `pgx` driver formats as RFC 3339. Go-side generation uses `time.RFC3339`. Both format-consistent. Update test at `store_test.go:542` to use RFC 3339 format.

## 1.6 Package boundary enforcement

Create:
- `.golangci.yml` with depguard rules
- `internal/boundary/boundary_test.go` with compile-time negative tests

Boundary matrix:

```
epistemic     → may NOT import: application, work, api, mcp, ui, domainpack
work          → may NOT import: epistemic
verifier      → may NOT import: epistemic authority operations
mcp           → may NOT import: decision/discharge/promote/retract
application   → may import: epistemic, work, domainpack (interface), verifier, packet
application   → may NOT import: mcp, ui
```

## 1.7 Auth boundary

- MCP mode (`internal/mcp/`): only exposes get_context and submit_packet. No import of decision/discharge packages.
- HTTP mode (`internal/ui/`): human endpoints require session + CSRF. Agent endpoints: none (no agent HTTP surface).

==================================================
PHASE 2 — INTEGRATION / END-TO-END FLOW
==================================================

Goal: Wire all components together and verify the complete EBP packet flow.

## 2.1 Wire MCP adapter

Convert `mcp/adapter/adapter.go` to use `internal/application.App` directly.

- `HandleTool("argus.get_context", ...)` → `app.GetContext(ctx, taskID)`
- `HandleTool("argus.submit_packet", ...)` → `app.SubmitPacket(ctx, packet)`

No HTTP in the path.

## 2.2 Wire Trust UI

Convert `trust-ui/server.go` to `internal/ui/server.go`.

- Embed templates via `//go:embed`
- Replace HTTP calls with direct `app.GetContext()` and `app.SubmitDecision()` calls
- Serve from `argus serve` on `:8080`
- Debts view shows proposed retirement metadata from belief row

## 2.3 Wire verifier + trust binding

Keep `verifier/` as-is. `argus verify` calls `verifier.RunPhysicsVerifier()` directly.

For in-process verification during packet validation: call the verifier library from `internal/application/compile.go`.

Verifier trust binding — VerificationArtifact schema:

```go
type VerificationArtifact struct {
    RunID             string `json:"run_id"`
    VerifierID        string `json:"verifier_id"`
    VerifierVersion   string `json:"verifier_version"`
    VerificationInput string `json:"verification_input_hash"` // SHA-256 of inputs
    ArtifactHash      string `json:"artifact_hash"`            // SHA-256 excluding this field + run_id
    Tolerance         string `json:"tolerance,omitempty"`      // e.g., "1e-10"
    Steps             []Step `json:"steps"`
    Result            string `json:"result"` // confirmed/refuted/inconclusive
    EvidenceRef       string `json:"evidence_ref,omitempty"`
}
```

Canonical artifact hash: SHA-256 of JSON serialization of all fields EXCEPT `artifact_hash` and `run_id`. Document the exclusion.

Artifact persistence: artifacts are stored as evidence rows with `provenance_class = 'reproducible_artifact'`. The CLI (`argus verify`) does NOT write to DB directly. It produces a JSON artifact that enters via `submit_packet` (agent) or `argus verify --submit` (application layer, same validation path).

VerifierSpec enforcement: at packet compile time, the application checks the artifact's `verifier_id` + `verifier_version` against the pack's `VerifierSpec` list. If not in the allowed list, the packet is rejected.

Verifier failure semantics:
- Verifier error/crash → artifact recorded as `inconclusive`, packet proceeds
- Verifier refuted → evidence recorded with `result: refuted`, packet proceeds
- Add timeout to synchronous compile path (configurable, default 30s)

## 2.4 EBP packet submission flow (revised)

Agent packets NEVER perform debt discharge, promotion, retraction, or authorization.

```
Agent → MCP → app.SubmitPacket(packet)
  → compile(packet)
     → validate packet structure
     → resolve references (local:, canonical:belief:<uuid>)
     → verify scenario containment (all refs in same scenario)
     → if verification evidence: invoke verifier, check VerifierSpec, produce artifact
     → compute content hash
  → persist(packet)
     → epistemic.EnterBelief     [Solvent — own transaction, ON CONFLICT DO NOTHING]
     → epistemic.AddEvidence     [Solvent — own transaction, ON CONFLICT DO NOTHING]
     → work.CreateTask           [Work store — own transaction, ON CONFLICT DO NOTHING]
     → if proposed_retirement: store metadata on belief
     → audit.Log                 [Solvent audit]
  → write submission_idempotency [AFTER all entities succeed]
  → return result
```

Transaction model: Solvent authority mutations remain transactionally atomic within Solvent. Cross-subsystem packet persistence is not globally atomic. Partial failures are explicit, retryable, and idempotent via entity-level ON CONFLICT.

What submit_packet does NOT do:
- RetireDebt (human decision only — Discharge path)
- Promote (human decision only)
- Retract (human decision only)
- Authorize (human decision only)

## 2.4a Agent-authored edge behavior

Agent-created contradicts/derives edges enter as `belief_edge` rows. They are epistemic proposals/evidence, NOT authority transitions.

Enforceable invariant:
- No agent HTTP surface exists (agents are MCP-only, 2 tools)
- MCP does not expose retract, promote, or discharge tools
- Application layer never acts on agent edges for retract/promote/authorize without human confirmation
- Boundary enforcement: `mcp` package cannot import retract/promote/discharge packages

No schema state change needed. Application-layer discipline is sufficient for this POC.

## 2.5 Context assembly flow

```
Agent → MCP → app.GetContext(taskID)
  → work.GetTask(taskID)              [CRDB read]
  → work.ListDependencies(taskID)     [CRDB read]
  → epistemic.GetSnapshot(scenario)   [Solvent view, CRDB read]
  → epistemic.GetActivity(scenario)   [Solvent audit, CRDB read]
  → assemble RCP context
  → set truncated flag if results exceed limit
  → return
```

Single-DB reads. `truncated` field in response: true if any result set exceeds configurable limit (default: 1000 rows). `TRUNCATED ≠ COMPLETE`.

## 2.6 Human decision flow

```
Human → Browser → Login (token → session cookie)
  → UI → app.SubmitDecision(type, params)
  → validate session + CSRF
  → validate retirement rule [domain pack]
  → epistemic.Discharge(...)   [if type=RETIRE_DEBT — own Solvent transaction]
  → epistemic.Promote(...)     [if type=PROMOTE — own Solvent transaction]
  → epistemic.Retract(...)     [if type=RETRACT — own Solvent transaction]
     → query conductor_task WHERE governance_ref = belief_id
     → cancel each matching task [own work transaction]
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
7. Debt visible in UI (with proposed retirement metadata if applicable)
8. `argus verify` produces verification artifact
9. Adversarial agent submits contradiction packet
10. Contradiction visible in UI
11. Human reviews in Trust UI
12. Human discharges debt (via UI, using Discharge method)
13. Claim promoted (via Solvent kernel gate — debt must be empty)
14. Human retracts (via UI)
15. Linked task cancelled (dead end — governance_ref lookup)
16. REOPEN creates new lineage (reopened_from_task_id)
17. UI shows full history

Adversarial-agent independence is procedural in this demo (separate MCP session, no shared state).

## 3.2 Refusal test suite (executable specifications)

| # | Test | Entry | Input | Expected | Assertion |
|---|------|-------|-------|----------|-----------|
| 1 | Agent cannot discharge | MCP submit_packet | Packet with `proposed_retirement: ["needMap"]` | Proposal recorded as metadata, debt NOT retired | `SELECT debt FROM belief WHERE id = X` unchanged |
| 2 | Agent cannot promote | MCP (no promote tool) | N/A | Tool not in ListTools | `len(tools) == 2`, no "promote" |
| 3 | Open debt → promote | HTTP POST /api/promote | belief_id with non-empty debt | HTTP 409 "promotion blocked: open debt" | `SELECT status FROM belief` remains "entered" |
| 4 | Retracted belief → promote | HTTP POST /api/promote | retracted belief_id | HTTP 409 "belief retracted" | `SELECT status` remains "retracted" |
| 5 | Agent cannot retract | MCP (no retract tool) | N/A | Tool not in ListTools | `len(tools) == 2`, no "retract" |
| 6 | Bad Origin → REFUSED | HTTP POST /api/retire | `Origin: http://evil.com` | HTTP 403 | No audit row for request |
| 7 | Body-supplied principal → IGNORED | HTTP POST /api/retire | Body: `{"discharged_by": "attacker"}` | Request succeeds, body field ignored | `discharged_by` in audit = server principal |
| 8 | Failed mutation → no success audit | HTTP POST /api/promote | belief_id with open debt | HTTP 409 | `SELECT * FROM audit_activity WHERE type = 'promotion_granted'` returns 0 |
| 9 | Partial retry completes | Simulate failure after first entity write | Packet where second entity write fails, then retry | Retry completes, no duplicates | Entity count matches; idempotency row status = completed |

## 3.3 BM-IST import isolation test

Rename: "BM-IST import isolation" (not "portability proof").

Create `internal/domainpack/isolation_test.go`:

Uses `go/packages` to analyze import graph of core packages. Verifies no core package imports BM-IST-specific symbols. Combined with depguard for runtime enforcement.

Add behavioral smoke test: register a minimal stub pack, verify core can query its definition, retirement rules, and evidence classes without domain-specific knowledge.

## 3.4 Simplification ledger verification

| Dimension | Before | After | Verification method |
|-----------|--------|-------|-------------------|
| Binaries | ~5 | 1 | count `main.go` in `cmd/` |
| Process instances | ~7 (+ OpenCode) | 2 (argus serve + argus mcp) + OpenCode | count `net.Listen` + MCP stdio |
| Ports | ~6 | 1 (HTTP+UI; MCP stdio) | count `net.Listen` calls |
| Databases | 1-2 (CockroachDB + SQLite) | 1 (CockroachDB) | single connection string |
| Go modules | 3 | 1 (+ Solvent import) | count `go.mod` files |
| Credential sets | multiple | 1 | single env var |
| Internal HTTP calls | many | 0 | `grep -r "http.Client" internal/` returns nothing |
| Cross-process handoffs | yes | none | no file-based artifact passing |

## 3.5 Delete old code + tests

Delete:
- `coordinator/http/` — HTTP handler
- `coordinator/client.go` — REST clients
- `coordinator/mock.go` — rewrite for in-process mocks
- `trust-ui/` as separate module
- `reference-loop/` — replaced by in-process tests
- HTTP handler tests for old service topology
- Tests that encode deployment mechanics

Preserve:
- `coordinator/compiler.go` logic → `internal/application/compile.go`
- `coordinator/validate.go` logic → `internal/application/validate.go`
- `coordinator/persist.go` logic → `internal/application/persist.go`
- `coordinator/context.go` logic → `internal/application/context.go`
- `coordinator/human.go` logic → `internal/application/human.go`
- All packet validation tests
- All domain pack validation tests
- All verifier tests

Test disposal:

| Category | Tests | Action |
|----------|-------|--------|
| PRESERVE | coordinator/*_test.go (compiler, retirement rules) | Move to internal/application/ |
| PRESERVE | packet/v1/*_test.go | Keep as-is |
| PRESERVE | domain-pack/*_test.go | Keep as-is |
| PRESERVE | verifier/*_test.go | Keep as-is |
| REWRITE | coordinator/http/handler_test.go | Rewrite as application unit tests |
| REWRITE | mcp/adapter/adapter_test.go | Rewrite to use in-process App |
| DELETE | trust-ui/ tests (separate module) | Embedded, tested via integration |
| DELETE | reference-loop/*_test.go | Replaced by in-process tests |
| NOT APPLICABLE | Solvent API handler tests (solvent repo) | Not imported |
| NOT APPLICABLE | Conductor HTTP/MCP/Web tests (conductor repo) | Not imported |
| ADD | internal/boundary/boundary_test.go | Import-direction enforcement |
| ADD | internal/domainpack/isolation_test.go | BM-IST import isolation |
| ADD | internal/domainpack/contract_test.go | Generic pack contract smoke |
| ADD | Refusal test suite (9 cases) | §3.2 |
| ADD | Partial-failure retry test | #16b |
| ADD | End-to-end integration test | Full 17-step narrative |

## 3.6 Taskfile.yml

```yaml
version: '3'
tasks:
  dev:
    desc: Start ARGUS in development mode (non-destructive)
    cmds:
      - |
        # Start CockroachDB with health polling
        cockroach demo --no-example-database --listen-addr :26260 &
        COCKROACH_PID=$!
        until cockroach node status --host :26260 --insecure 2>/dev/null; do
          sleep 1
        done
        # Start ARGUS (assumes schema already applied)
        go run ./cmd/argus serve --db postgres://root@localhost:26260/argus?sslmode=disable

  fresh:
    desc: Reset database and start fresh (destructive)
    cmds:
      - go run ./cmd/argus reset --db postgres://root@localhost:26260/argus?sslmode=disable
      - go run ./cmd/argus serve --db postgres://root@localhost:26260/argus?sslmode=disable

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

`task dev` is non-destructive. `task fresh` is destructive. `argus reset` remains explicit.

==================================================
MIGRATION MODEL
==================================================

- In-place pivot on a feature branch.
- Rollback: if Phase 0 gate is REVISE, abandon branch. If Phase 1+ fails, revert to main.
- No live state to protect (POC, fresh DB each run).
- Solvent module imported via `go.mod` replace directive during development.

==================================================
RISKS AND NON-GOALS
==================================================

Risks:
- Solvent public projection contract requires Solvent repo change (mitigated: bounded, backward-compatible)
- Solvent migration export requires Solvent repo change (mitigated: bounded, backward-compatible)
- Separate-transaction model means partial failures possible (mitigated: entity-level idempotency)
- If Solvent declines public API: copy ~200 lines of view logic (accepted drift risk for POC)

Non-goals (confirmed unchanged):
- Production-grade auth system
- Container-per-agent isolation
- Elaborate verifier attestation
- Second domain pack
- Distributed deployment support
- Performance optimization
- Global atomicity across subsystems
- Full identity infrastructure
- Cryptographic attestation
- Generalized architecture framework

==================================================
ACCEPTANCE CRITERIA
==================================================

1. `argus serve` starts and serves Trust UI on :8080
2. `argus mcp` starts and exposes exactly 2 tools
3. `argus verify` runs physics verifier and produces artifact with verifier_id, version, input_hash, tolerance
4. Agent can submit EBP packet via MCP → beliefs + evidence + debt recorded in CRDB
5. Agent submit_packet does NOT perform RetireDebt, Promote, Retract, or Authorize
6. Human can discharge debt via Trust UI → debt retired in CRDB (attribution from server-side principal)
7. Human can promote claim → belief promoted in CRDB (only if debt empty)
8. Adversarial packet introduces contradiction → visible in UI
9. Retraction causes linked task cancellation (governance_ref lookup)
10. Agent-created edges do NOT trigger retraction or block promotion
11. `go test ./...` passes
12. `golangci-lint run` passes (including depguard rules)
13. BM-IST import isolation test passes (go/packages + depguard)
14. Simplification ledger verified (0 internal HTTP calls, 2 instances, 1 DB)
15. Refusal test suite passes (9 executable test cases, §3.2)
16. Idempotency: duplicate packet submission returns same result, no duplicate state
16b. Partial-failure retry completes the packet with no duplicate state
17. Solvent owns canonical migrations; ARGUS owns only work/idempotency migrations
18. Proposed retirement metadata visible in Trust UI Debts view
19. VerifierSpec enforcement: packet with unknown verifier_id is rejected at compile time
20. REOPEN creates new task linked to old via reopened_from_task_id
21. `task dev` is non-destructive; `task fresh` is destructive
