# ARGUS POC Simplification Pivot — Implementation Plan (Frozen)

Frozen after three rounds of adversarial review. 12 amendments + 10 fixes + 12 final corrections incorporated.

Status: PLAN ONLY — no implementation code.

This is the final plan revision. Next step: Phase 0/1 implementation checklist execution.

==================================================
EXECUTIVE SUMMARY
==================================================

The POC infrastructure has drifted from the thesis it demonstrates. This plan collapses the multi-service stack (Solvent, Conductor, Coordinator, MCP, Trust UI) into:

    one Go binary (argus)
    one CockroachDB database
    two process instances (argus serve + argus mcp) + OpenCode

while preserving:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Solvent change budget: 2 bounded backward-compatible additions
  1. Public projection contract (GetSnapshot, ExplainSnapshot on ledger.Service)
  2. Migration export package (embed.FS or Apply function)

Any third Solvent change beyond this budget automatically reopens Phase 0 as REVISE.

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

==================================================
RECONNAISSANCE FINDINGS
==================================================

## Solvent — reusable in-process

Kernel: `kernel.Store` wrapping `*sql.DB`, 22-method Contract. All writes via `crdb.ExecuteTx`. No HTTP.

Services (all take `*sql.DB`):
- `service/ledger.Service` — thin passthrough + audit
- `service/audit.Service` — append-only activity ledger
- `service/authority.Service` — execution boundary
- `service/policy.Service` — constraint evaluation

Read-only projections (currently internal/):
- `internal/view.GetSnapshot` — beliefs + evidence + intents
- `internal/view.ExplainSnapshot` — human-readable "why" projection

Discardable: `api/`, `internal/wizard/`, cmd binaries.

Schema: 10 migrations, 20+ tables. Solvent-owned.

## Conductor — extract minimal domain model

4 types: Task (6 states), Project, Dependency, Activity.
Storage: SQLite, 4 tables.
All semantics load-bearing. Nothing discarded.

## Oracle — library, no standalone binary

Root module: zero external deps. Components: coordinator (library), mcp/adapter (2 tools), packet/v1, domain-pack, verifier, corpus, trust-ui (separate module), reference-loop (separate module).

## Database topology

Current: CockroachDB (Solvent) + SQLite (Conductor).
Target: One CockroachDB instance, one database.

==================================================
PHASE 0 — RECONNAISSANCE + GO/NO-GO
==================================================

## 0.1 Complexity classification

Done. Majority deployment complexity. Pivot justified.

## 0.2 Work-semantics inventory

All Conductor semantics load-bearing and small.

## 0.3 Solvent kernel reuse analysis

Kernel importable as Go module. `kernel.Store` takes `*sql.DB`. No HTTP.

Gap: `internal/view.GetSnapshot` and `internal/view.ExplainSnapshot` are internal.

Solvent change budget (2 bounded additions):
  1. Public projection contract on ledger.Service
  2. Migration export package (embed.FS or Apply)

Both backward-compatible, SemVer-minor.

**Phase 0 verification step:** Before Phase 1, verify that Solvent's existing uniqueness constraints cover the entity-level idempotency pattern required by §1.4a. Specifically check:
- belief table: unique constraint on (scenario_id, claim) or equivalent
- evidence table: unique constraint on (scenario_id, content_sha256) or equivalent
- belief_edge table: unique constraint on (scenario_id, parent_id, child_id, kind) or equivalent

If these constraints do not exist, that becomes a third Solvent change-budget item. If the budget is exceeded, reopen Phase 0 as REVISE. Do not silently add constraints to Solvent-owned tables.

Decision: import Solvent as a Go module.

## 0.4 Component disposition table

| Component | Disposition | Rationale |
|-----------|-------------|-----------|
| Solvent kernel (`kernel/`) | KEEP — import as module | Core epistemic model |
| Solvent services (`service/`) | KEEP — import as module | Thin wrappers |
| Solvent view projections | KEEP — public contract | GetSnapshot, ExplainSnapshot |
| Solvent API (`api/`) | DELETE | HTTP handlers |
| Solvent wizard | DELETE | HTTP wizard |
| Solvent cmd binaries | DELETE | Standalone binaries |
| Conductor domain types | KEEP — extract to `internal/work/` | All load-bearing |
| Conductor store | CONVERT — rewrite for CRDB | SQLite→CRDB |
| Conductor HTTP/MCP/Web | DELETE | Transport layers |
| Coordinator library | CONVERT — `internal/application/` | Orchestration logic |
| Coordinator HTTP handler | DELETE | Thin routing |
| MCP adapter | CONVERT — `internal/mcp/` | Library |
| Trust UI | CONVERT — embed in argus process | Templates + direct calls |
| Verifier | KEEP — library | Already clean |
| Packet v1 | KEEP — library | Already clean |
| Domain pack | KEEP — library | Already clean |
| Corpus | KEEP — library | Already clean |
| Reference loop | DELETE | Replaced by in-process tests |
| Solvent authority.Service | DEFER — not exercised by POC | No execution path in demo |
| Solvent executor.Registry | DEFER — not exercised by POC | No external execution in demo |

## 0.5 Boundary enforcement design

Import-direction rules:

```
epistemic     → may NOT import: application, work, api, mcp, ui, domainpack
work          → may NOT import: epistemic (may read via application layer)
verifier      → may NOT import: epistemic authority operations
mcp           → may NOT import: decision/discharge/promote/retract
application   → may import: epistemic, work, domainpack (interface), verifier, packet
application   → may NOT import: mcp, ui
```

Enforcement:
1. `depguard` rules in `.golangci.yml`
2. Compile-time negative test in `internal/boundary/boundary_test.go`
3. Application boundary row in matrix

## 0.6 Auth model design

| Caller | Entry point | Operations |
|--------|-------------|-----------|
| Agent (MCP) | `argus mcp` | get_context, submit_packet |
| Human (HTTP) | `argus serve` | discharge, promote, retract, reopen, authorize, views |

No agent HTTP surface. Agents are MCP-only.

Credential mechanism:

```
single configured operator token (env/config)
        ↓
GET /login validates token → HttpOnly session cookie (HMAC-signed, in-memory)
        ↓
consequential POST requires valid session + CSRF (Origin/Host)
        ↓
server-derived principal (never from request body)
```

Session: expires after configurable timeout (default 1h). Logout clears session. Rotation: change env var, restart.

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
type ClaimType string   // "derived", "accommodated", "postulated"
type DebtItem string    // "needMap", "needInvariant", etc.
type EvidenceClass string // "reproducible_artifact", "operator_asserted"

type RetirementRule struct {
    DebtItem      DebtItem      `json:"debt_item"`
    RequiredClass EvidenceClass `json:"required_class"`
    Rule          string        `json:"rule"`
    RequiresHuman bool          `json:"requires_human"`
}

type Falsifier struct {
    ID          string `json:"id"`
    Description string `json:"description"`
}

type VerifierSpec struct {
    VerifierID string `json:"verifier_id"`
    MinVersion string `json:"min_version"`
}
```

BM-IST retirement rules:

| Debt Item | Required Class | Rule | Human |
|-----------|---------------|------|-------|
| needMap | reproducible_artifact | Faithful map reproduction | yes |
| needInvariant | reproducible_artifact | Invariant preservation | yes |
| needToyCheck | reproducible_artifact | Toy-model consistency | yes |
| needNullModel | reproducible_artifact | Null-model ruled out | yes |
| needObstruction | reproducible_artifact | Obstruction criterion | yes |
| needFaithfulnessReview | reproducible_artifact | Reference result within tolerance | yes |

Validate() checks: every debt item has retirement rule, every rule references declared class, no duplicates, valid semver in VerifierSpec.

Core knows: debt exists, evidence exists, rules evaluable.
Core does NOT know: needMap, needInvariant, BM-IST specifics.

## 0.8 Discharge vs RetireDebt

- `RetireDebt(ctx, scenario, beliefID, item)` — kernel method. ARGUS-internal only. Not reachable from MCP or HTTP. Application calls it within the Discharge path.
- `Discharge(ctx, scenario, beliefID, obligationKey, instrumentRef, dischargedBy)` — higher-level method with replay protection. Human-facing API.

Boundary enforcement: depguard rule prevents MCP and UI packages from importing `RetireDebt`. Only `internal/application/` may call it.

## 0.9 Dependency semantics

Dependency = "task with unfinished blocker cannot become active."

Formal rule:
```sql
-- A task can only transition from proposed to active if:
--   no conductor_dependency exists where task_id = this_task AND blocked_by_id references a task
--   whose status NOT IN ('accepted', 'cancelled')
```

Implemented in `work/store.go` Transition method: before claiming, check for unresolved dependencies.

## 0.10 NO-GO / REVISE / GO triggers

NO-GO if:
- Solvent kernel cannot be reused without domain/application imports.
- Conductor semantics cannot be preserved without recreating old service.
- Human/agent authority separation cannot be mechanically preserved.
- Thesis invariant weakened.

REVISE only for:
- Material scope expansion.
- Weakened thesis invariants.
- Solvent changes requiring API redesign.

**Governance rule:** Any third Solvent change beyond the declared two-item budget automatically reopens Phase 0 as REVISE.

GO when:
- Load-bearing semantics identified.
- Replacements explicit.
- Solvent changes within budget (2 bounded backward-compatible additions).
- No thesis invariant weakened.

## 0.11 GO decision

Solvent change budget: 2 items (public projections + migration FS). Both backward-compatible.

Phase 0 verification step: check Solvent uniqueness constraints before Phase 1.

**Decision: GO.**

==================================================
PHASE 1 — CORE COLLAPSE
==================================================

## 1.1 Create `cmd/argus/main.go`

Subcommands: `argus serve`, `argus mcp`, `argus verify`, `argus reset`.

## 1.2 Extract Solvent kernel → `internal/epistemic/`

Import from `github.com/PithomLabs/solvent`. Thin wrappers re-exporting kernel, ledger, audit, authority, policy.

## 1.2a Solvent public projection contract (Change 1)

Solvent exposes on `service/ledger.Service`:
- `GetSnapshot(ctx, scenarioID, opts) (*view.Snapshot, error)`
- `ExplainSnapshot(scenario, scenarioID, snap) *view.ExplainResult`

Public types: Snapshot, SnapshotOpts, Belief, Evidence, Intent, ExplainResult, BeliefExplain.

Stability: documented as "unstable, may change without notice" for POC.

## 1.2b Migration export package (Change 2)

Solvent exports:
```go
package solventmigrations
//go:embed *.sql
var FS embed.FS
func Apply(db *sql.DB) error { ... }
```

ARGUS calls `solventmigrations.Apply(db)` then applies own migrations.

Fallback if Solvent declines: copy ~200 lines of view logic into ARGUS.

## 1.2c Domain Pack contract

`internal/domainpack/pack.go` with Pack interface and PackDefinition types (§0.7).

## 1.3 Extract Conductor → `internal/work/`

Task, Project, Dependency, Activity types + CockroachDB store.

## 1.4 Extract Coordinator → `internal/application/`

App struct, compile, validate, persist, context, human decisions, proposed retirement handling.

## 1.4a Entity-level idempotency

Entity IDs are deterministic: derived from scenario_id + canonical content hash. Retries return the same identifiers rather than silently discarding newly generated UUIDs.

```go
func entityID(scenarioID, entityType, content string) string {
    h := sha256.Sum256([]byte(scenarioID + ":" + entityType + ":" + content))
    return hex.EncodeToString(h[:16]) // first 16 bytes = 32 hex chars
}
```

Write order:
1. Compute content hash from packet
2. Persist beliefs — `INSERT ... ON CONFLICT (scenario_id, claim_hash) DO NOTHING` (using deterministic entity IDs, ON CONFLICT becomes safe no-op)
3. Persist evidence — `INSERT ... ON CONFLICT (scenario_id, content_sha256) DO NOTHING`
4. Persist edges — `INSERT ... ON CONFLICT (scenario_id, from_id, to_id, kind) DO NOTHING`
5. Persist tasks — deterministic ID, `INSERT ... ON CONFLICT DO NOTHING`
6. Write `submission_idempotency` row LAST — row existence = completed

On partial failure:
- Succeeded entity writes are idempotent (ON CONFLICT DO NOTHING + deterministic IDs)
- Retry resumes from step 2 — already-written entities are no-ops
- No idempotency row → packet not complete
- Duplicate submission of completed packet → idempotency row exists, reconstruct result

Duplicate-result behavior: reconstruct result from persisted entities by packet_id. Query beliefs, evidence, edges, tasks where packet_id matches. Return reconstructed RCP context. No separate cache needed.

Schema:
```sql
CREATE TABLE IF NOT EXISTS submission_idempotency (
    content_hash  STRING NOT NULL,
    scenario_id   STRING NOT NULL,
    packet_id     STRING NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(content_hash, scenario_id)
);
```

No status column. Row existence = completed.

## 1.5 CockroachDB schema

### Migration ownership

Solvent owns epistemic/authority migrations. ARGUS applies via `solventmigrations.Apply(db)`.

ARGUS owns: work tables, idempotency table, proposed_retirement side table.

### Schema layout

```
-- Epistemic (Solvent-owned, applied via solventmigrations.Apply)
belief, belief_edge, evidence, action_intent,
principal, authority_target, target_snapshot, target_activation,
target_revocation, justification, debt_discharge,
audit_activity

-- Proposed retirement (ARGUS-owned, new migration)
CREATE TABLE belief_retirement_proposal (
    id              STRING PRIMARY KEY,
    belief_id       STRING NOT NULL REFERENCES belief(id),
    scenario_id     STRING NOT NULL,
    debt_items      JSONB NOT NULL,           -- ["needMap", "needInvariant"]
    proposer        STRING NOT NULL,           -- agent ID from packet
    proposed_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(belief_id, scenario_id)
);

-- Work (ARGUS-owned, new migration)
conductor_project, conductor_task, conductor_activity, conductor_dependency

-- Idempotency (ARGUS-owned, new migration)
submission_idempotency
```

`proposed_retirement` is NOT on the Solvent-owned `belief` table. It is an ARGUS-owned side table with FK to belief.

### Conductor→CRDB type mapping

| SQLite | CRDB | Notes |
|--------|------|-------|
| `id TEXT PRIMARY KEY` | `id STRING PRIMARY KEY` | UUID from Go |
| `TEXT NOT NULL` | `STRING NOT NULL` | Direct |
| `TEXT DEFAULT '...' CHECK (...)` | `STRING DEFAULT '...' CHECK (...)` | Direct |
| `TEXT` (nullable) | `STRING` (nullable) | Direct |
| `TEXT DEFAULT (datetime('now'))` | `TIMESTAMPTZ DEFAULT now()` | Function change |
| `TEXT CHECK (actor_type IN (...))` | `STRING CHECK (actor_type IN (...))` | Direct |
| `TEXT` (JSON) | `STRING` (JSON) | JSON as text |

`?` → `$N`. `datetime('now')` → `now()`.

`governance_ref TEXT` → `governance_ref STRING REFERENCES belief(id)` (FK upgrade, one DB makes this possible).

### REOPEN lineage

Add `reopened_from_task_id STRING REFERENCES conductor_task(id)` to `conductor_task`.

### Dead-end task lookup

```sql
SELECT id FROM conductor_task WHERE governance_ref = :belief_id AND status NOT IN ('cancelled', 'accepted')
```
Cancel each in its own transaction.

### CRDB retry wrapper

```go
func (db *DB) Transaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
    const maxAttempts = 5
    for attempt := 0; attempt < maxAttempts; attempt++ {
        tx, err := db.BeginTx(ctx, nil)
        if err != nil { return err }
        if err := fn(tx); err != nil {
            tx.Rollback()
            if isRetryable(err) && attempt < maxAttempts-1 { continue }
            return err
        }
        if err := tx.Commit(); err != nil {
            if isRetryable(err) && attempt < maxAttempts-1 { continue }
            return err
        }
        return nil
    }
    return fmt.Errorf("transaction failed after %d attempts", maxAttempts)
}
```

Retryable = SQLSTATE 40001 anywhere in closure (statement or commit), not just commit.

### Release() bug fix

Wrap Release() UPDATE + activity INSERT in a transaction.

### Timestamp handling

Go types store timestamps as `string`. CRDB `TIMESTAMPTZ` via `pgx` formats as RFC 3339. Go-side uses `time.RFC3339`. Format-consistent. Update test.

## 1.6 Package boundary enforcement

```
epistemic     → may NOT import: application, work, api, mcp, ui, domainpack
work          → may NOT import: epistemic
verifier      → may NOT import: epistemic authority operations
mcp           → may NOT import: decision/discharge/promote/retract
application   → may import: epistemic, work, domainpack (interface), verifier, packet
application   → may NOT import: mcp, ui
```

Enforcement: depguard + compile-time negative test.

RetireDebt is kernel-internal. Depguard prevents MCP and UI packages from importing it. Only `internal/application/` may call it within the Discharge path.

## 1.7 Auth boundary

MCP mode: only get_context, submit_packet. No decision imports.
HTTP mode: session + CSRF. No agent HTTP surface.

==================================================
PHASE 2 — INTEGRATION / END-TO-END FLOW
==================================================

## 2.1 Wire MCP adapter

`HandleTool("argus.get_context", ...)` → `app.GetContext(ctx, taskID)`
`HandleTool("argus.submit_packet", ...)` → `app.SubmitPacket(ctx, packet)`

## 2.2 Wire Trust UI

Embed templates. Direct `app.GetContext()` and `app.SubmitDecision()` calls.
Debts view renders proposed retirement from `belief_retirement_proposal` table.

UI import boundary: `internal/ui/` may import `internal/application/` (read-only queries + decision submission). May NOT import `internal/epistemic/` directly. Must go through application layer.

## 2.3 Wire verifier + trust binding

VerificationArtifact schema:
```go
type VerificationArtifact struct {
    RunID             string `json:"run_id"`
    VerifierID        string `json:"verifier_id"`
    VerifierVersion   string `json:"verifier_version"`
    VerificationInput string `json:"verification_input_hash"`
    ArtifactHash      string `json:"artifact_hash"`
    Tolerance         string `json:"tolerance,omitempty"`
    Steps             []Step `json:"steps"`
    Result            string `json:"result"`
    EvidenceRef       string `json:"evidence_ref,omitempty"`
}
```

Canonical hash: SHA-256 of JSON excluding `artifact_hash` and `run_id`.

Artifact persistence: evidence rows with `provenance_class = 'reproducible_artifact'`.

Verifier trust path:
- `argus verify --submit`: goes through application validation path (same as packet submission)
- Packet-supplied artifacts: revalidated against pack's VerifierSpec + input binding
- Never trust agent-claimed verifier results
- VerifierSpec enforcement: compile-time check that artifact's verifier_id + version are in pack's allowed list. If not, packet rejected.

Verifier failure semantics:
- Error/crash → artifact recorded `inconclusive`, packet proceeds
- Refuted → evidence with `result: refuted`, packet proceeds
- Timeout: configurable, default 30s

## 2.4 EBP packet submission flow

Agent packets NEVER perform Discharge, Promote, Retract, or Authorize.

```
Agent → MCP → app.SubmitPacket(packet)
  → compile(packet)
     → validate structure
     → resolve references (local:, canonical:belief:<uuid>)
     → scenario containment check (all refs in same scenario)
     → if verification evidence: invoke verifier, check VerifierSpec, produce artifact
     → compute content hash
  → persist(packet)
     → epistemic.EnterBelief     [Solvent, own tx, ON CONFLICT DO NOTHING]
     → epistemic.AddEvidence     [Solvent, own tx, ON CONFLICT DO NOTHING]
     → work.CreateTask           [Work, own tx, ON CONFLICT DO NOTHING]
     → if proposed_retirement: INSERT INTO belief_retirement_proposal
     → audit.Log
  → write submission_idempotency [AFTER all entities succeed]
  → return result
```

Transaction model: Solvent mutations atomic within Solvent. Cross-subsystem not globally atomic. Partial failures explicit, retryable, idempotent.

## 2.4a Agent-authored edge behavior

Agent edges stored as `belief_edge` rows. Epistemic proposals, not authority transitions.

Enforceable invariant: no agent HTTP surface (MCP-only, 2 tools). MCP has no retract/promote/discharge. Application never acts on agent edges for authority transitions without human confirmation.

No schema state change needed. Application-layer discipline sufficient.

## 2.5 Context assembly

```
app.GetContext(taskID)
  → work.GetTask(taskID)
  → work.ListDependencies(taskID)
  → epistemic.GetSnapshot(scenario)
  → epistemic.GetActivity(scenario)
  → assemble RCP context
  → set truncated flag if results exceed limit
  → return
```

`truncated` field: true if any result set exceeds configurable limit (default 1000). `TRUNCATED ≠ COMPLETE`.

## 2.6 Human decision flow

```
Human → Login (token → session) → UI
  → app.SubmitDecision(type, params)
  → validate session + CSRF
  → validate retirement rule [domain pack]
  → epistemic.Discharge(...)   [RETIRE_DEBT — own Solvent tx]
  → epistemic.Promote(...)     [PROMOTE — own Solvent tx]
  → epistemic.Retract(...)     [RETRACT — own Solvent tx]
     → query conductor_task WHERE governance_ref = belief_id
     → cancel each matching task
  → audit.Log(...)
```

Attribution from server-side principal, never request body.

==================================================
PHASE 3 — VERIFICATION, PORTABILITY PROOF + CLEANUP
==================================================

## 3.1 End-to-end demo walkthrough

17-step narrative with adversarial-agent independence as procedural (separate MCP session, shared system state intentionally, no conversational memory).

## 3.2 Refusal test suite

| # | Test | Entry | Input | Expected | Assertion |
|---|------|-------|-------|----------|-----------|
| 1 | Agent cannot discharge | MCP submit_packet | Packet with proposed_retirement | Metadata recorded, debt NOT retired | belief.debt unchanged |
| 2 | Agent cannot promote | MCP (no tool) | N/A | Tool not in ListTools | len=2, no "promote" |
| 3 | Open debt → promote | HTTP /api/promote | belief with debt | 409 "open debt" | status unchanged |
| 4 | Retracted → promote | HTTP /api/promote | retracted belief | 409 "retracted" | status unchanged |
| 5 | Agent cannot retract | MCP (no tool) | N/A | Tool not in ListTools | len=2, no "retract" |
| 6 | Bad Origin → REFUSED | HTTP /api/retire | Origin: evil.com | 403 | No audit row |
| 7 | Body principal → IGNORED | HTTP /api/retire | discharged_by in body | Success, body ignored | audit.discharged_by = server principal |
| 8 | Failed mutation → no audit | HTTP /api/promote | belief with debt | 409 | No promotion_granted audit |
| 9 | Partial retry completes | Fail after first entity | Retry same packet | Completes, no duplicates | Entity count correct; idempotency row exists |

## 3.3 BM-IST import isolation test

Uses `go/packages` to analyze import graph. Verifies no core package imports BM-IST symbols. Combined with depguard.

Behavioral smoke test: register stub pack, verify core queries its definition without domain knowledge.

## 3.4 Simplification ledger

| Dimension | Before | After | Verification |
|-----------|--------|-------|-------------|
| Binaries | ~5 | 1 | count main.go in cmd/ |
| Process instances | ~7 | 2 + OpenCode | count net.Listen + MCP |
| Ports | ~6 | 1 | count net.Listen |
| Databases | 1-2 | 1 | single connection string |
| Go modules | 3 | 1 + import | count go.mod |
| Credentials | multiple | 1 | single env var |
| Internal HTTP | many | 0 | grep -r "http.Client" internal/ → nothing |
| Cross-process handoffs | yes | none | no file artifact passing |

## 3.5 Delete old code + tests

Delete: coordinator/http/, coordinator/client.go, trust-ui/ as module, reference-loop/, deployment-mechanic tests.

Preserve: coordinator logic → internal/application/, packet tests, domain pack tests, verifier tests.

Test disposal:

| Category | Tests | Action |
|----------|-------|--------|
| PRESERVE | coordinator/*_test.go | Move to internal/application/ |
| PRESERVE | packet/v1/*_test.go | Keep |
| PRESERVE | domain-pack/*_test.go | Keep |
| PRESERVE | verifier/*_test.go | Keep |
| REWRITE | coordinator/http/handler_test.go | Application unit tests |
| REWRITE | mcp/adapter/adapter_test.go | In-process App |
| DELETE | trust-ui/ tests | Embedded |
| DELETE | reference-loop/*_test.go | In-process |
| NOT APPLICABLE | Solvent/Conductor tests | Not imported |
| ADD | boundary_test.go | Import enforcement |
| ADD | isolation_test.go | BM-IST import isolation |
| ADD | contract_test.go | Pack contract smoke |
| ADD | Refusal suite (9 cases) | §3.2 |
| ADD | Partial-retry test | #16b |
| ADD | End-to-end integration | 17-step narrative |

## 3.6 Taskfile.yml

```yaml
version: '3'
tasks:
  dev:
    desc: Start ARGUS in development mode (non-destructive)
    cmds:
      - |
        # Start persistent single-node CockroachDB
        mkdir -p .cockroach-data
        cockroach start-single-node --insecure --listen-addr :26260 --store=.cockroach-data &
        CRDB_PID=$!
        until cockroach node status --host :26260 --insecure 2>/dev/null; do
          sleep 1
        done
        # Apply migrations idempotently
        go run ./cmd/argus reset --db postgres://root@localhost:26260/argus?sslmode=disable
        # Start ARGUS
        go run ./cmd/argus serve --db postgres://root@localhost:26260/argus?sslmode=disable

  fresh:
    desc: Reset database and start fresh (destructive)
    cmds:
      - rm -rf .cockroach-data
      - mkdir -p .cockroach-data
      - cockroach start-single-node --insecure --listen-addr :26260 --store=.cockroach-data &
      - sleep 3
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

`task dev`: persistent CRDB, idempotent migrations, non-destructive.
`task fresh`: destructive reset (deletes data directory).
`argus reset`: explicit, destructive.

==================================================
SHARED VS CONVERSATIONAL STATE
==================================================

Agents share system state intentionally (beliefs, evidence, tasks in CRDB are visible to all agents via MCP context). Agents do NOT share conversational memory. Each MCP session is fresh; context is reconstructed from system state via `get_context`.

==================================================
MIGRATION MODEL
==================================================

- In-place pivot on feature branch.
- Rollback: REVISE → abandon branch. Phase 1+ failure → revert to main.
- No live state (POC, fresh DB each run).
- Solvent via `go.mod` replace directive during development.

==================================================
RISKS AND NON-GOALS
==================================================

Risks:
- Solvent public API change requires Solvent repo (mitigated: bounded, backward-compatible)
- Solvent migration export requires Solvent repo (mitigated: bounded, backward-compatible)
- If Solvent uniqueness constraints missing → third budget item → REVISE gate
- Separate-transaction model → partial failures (mitigated: entity-level idempotency)
- If Solvent declines: copy ~200 lines of view logic (accepted drift risk)

Non-goals:
- Production-grade auth
- Container-per-agent isolation
- Cryptographic attestation
- Second domain pack
- Distributed deployment
- Performance optimization
- Global atomicity
- Full identity infrastructure
- Generalized architecture framework
- External execution (authority.Service / executor.Registry deferred)

==================================================
GOVERNANCE RULE
==================================================

Any third Solvent change beyond the declared two-item budget (public projections + migration FS) automatically reopens Phase 0 as REVISE. The gate re-evaluates whether the pivot remains feasible within the expanded Solvent change cost.

==================================================
ACCEPTANCE CRITERIA
==================================================

1. `argus serve` starts and serves Trust UI on :8080
2. `argus mcp` starts and exposes exactly 2 tools
3. `argus verify` produces artifact with verifier_id, version, input_hash, tolerance
4. Agent submit_packet → beliefs + evidence + debt recorded in CRDB
5. Agent submit_packet does NOT perform Discharge, Promote, Retract, or Authorize
6. Human discharge debt via UI → debt retired (Discharge method, server-side attribution)
7. Human promote claim → belief promoted (only if debt empty)
8. Adversarial contradiction → visible in UI
9. Retraction → linked task cancelled (governance_ref FK lookup)
10. Agent edges do NOT trigger retraction or block promotion
11. `go test ./...` passes
12. `golangci-lint run` passes (depguard rules)
13. BM-IST import isolation test passes (go/packages + depguard)
14. Simplification ledger verified (0 internal HTTP, 2 instances, 1 DB)
15. Refusal suite passes (9 executable tests)
16. Idempotency: duplicate submission returns same result
16b. Partial-failure retry completes with no duplicate state
17. Solvent owns canonical migrations; ARGUS owns work/idempotency/retirement-proposal
18. Proposed retirement visible in Trust UI from belief_retirement_proposal table
19. VerifierSpec enforcement: unknown verifier_id rejected at compile
20. REOPEN creates task linked via reopened_from_task_id
21. `task dev` non-destructive; `task fresh` destructive
22. Dependency blocking: task with unfinished blocker cannot become active
23. Entity IDs deterministic: same content → same ID across retries
24. RetireDebt unreachable from MCP and UI (depguard enforced)
25. Third Solvent change reopens Phase 0 as REVISE
