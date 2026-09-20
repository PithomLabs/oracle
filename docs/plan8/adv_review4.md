# ARGUS FINAL ADVERSARIAL CODE REVIEW

**Date:** 2026-09-18  
**Reviewer:** Independent adversarial architect  
**Scope:** Complete Phase 8 POC implementation in `/home/chaschel/Documents/go/oracle`, `/home/chaschel/Documents/go/solvent-main`  
**Method:** Line-by-line inspection of production code, Taskfile, migrations, tests, and runtime behavior. All test suites executed (`go test -count=1 ./...`, `go vet ./...`). `golangci-lint` v2.1.6 is incompatible with the project's Go 1.25.0 toolchain and could not be run.

---

## P0 — BLOCKERS

### P0-1: `contentHash` Silently Collapses Semantically Distinct Packets

**File:** `internal/application/app.go:78-88`

```go
func contentHash(pkt *packetv1.Packet) string {
    h := sha256.New()
    h.Write([]byte(pkt.ScenarioID))
    h.Write([]byte{0})
    h.Write([]byte(pkt.PacketID))
    for _, b := range pkt.Beliefs {
        h.Write([]byte(b.Claim))
        h.Write([]byte{0})
    }
    return hex.EncodeToString(h.Sum(nil))
}
```

The canonical idempotency hash includes only `ScenarioID`, `PacketID`, and belief `Claim` strings. It **excludes** `Evidence`, `Edges`, and `Tasks`. Two packets with identical beliefs but different evidence, edges, or tasks produce the same hash. The `ON CONFLICT (content_hash, scenario_id) DO NOTHING` on `submission_idempotency` then causes the second packet's distinct entities to be silently discarded.

The Phase 8 plan §18 explicitly states the canonical hash must include sorted claims, evidence, edges, and packet content.

**Impact:** An agent submitting a follow-up packet with new evidence or adversarial edges for an already-submitted belief gets silent no-op success. New evidence never enters the ledger. Adversarial `contradicts` edges are dropped if the target belief was already entered by a work agent.

---

### P0-2: Application Layer Writes Directly to Solvent Tables, Bypassing Solvent's REST Boundary

**File:** `internal/application/app.go:186-266`

`App.Persist` executes raw SQL against `belief`, `evidence`, `belief_edge` — tables owned by Solvent — using the same `*sql.DB` handle. It never invokes Solvent's REST API (`POST /v1/beliefs`, `POST /v1/evidence`, `POST /v1/beliefs/{id}/edges`).

The Phase 8 plan §3 and §4 state:
> Solvent is a separate service/repository with a REST API.
> Coordinator must NOT write directly to Solvent DB.

**Impact:** The main binary does not exercise Solvent's REST endpoints, validation, or authority logic. Solvent's edge-creation validation (cross-scenario parent/child checks, self-edge prevention, kind enumeration) is bypassed. The `belief_edge` insert uses `ON CONFLICT (parent_id, child_id) DO NOTHING`, which silently drops a `contradicts` edge if a `derives` edge already exists for the same pair (see P0-4).

---

### P0-3: `cmdReset` Does Not Apply Solvent Migrations

**File:** `cmd/argus/main.go:182-183`

```go
// TODO: call solventmigrations.Apply(ctx, db) when Solvent migration FS is ready
```

A fresh developer running `task fresh` or `argus reset` gets a database with only ARGUS tables (`conductor_*`, `submission_idempotency`, `belief_retirement_proposal`). Solvent tables (`belief`, `evidence`, `belief_edge`, `action_intent`, `principal`, `debt_discharge`, etc.) are missing. The first `argus serve` then fails on any Solvent operation.

**Impact:** `task fresh` is not reproducible for a first-time developer. The integration tests work around this with a hand-rolled `integApplySolventMigrations` helper that is not available in the production reset path.

---

### P0-4: `ON CONFLICT (parent_id, child_id)` Silently Drops Conflicting Edge Kinds

**File:** `internal/application/app.go:257-264`

```go
_, err := tx.ExecContext(ctx,
    `INSERT INTO belief_edge (parent_id, child_id, kind)
     VALUES ($1::UUID, $2::UUID, $3)
     ON CONFLICT (parent_id, child_id) DO NOTHING`,
    fromID, toID, edge.Kind)
```

The Solvent schema PK is `(parent_id, child_id)` with no `kind` column. If a `derives` edge exists between beliefs A→B, a subsequent `contradicts` edge for the same A→B pair is silently discarded. The agent receives no error, no audit entry, and the adversarial challenge never reaches the ledger.

**Impact:** The dry-run Branch A ("adversarial contradicts → human retracts") can silently fail if the work agent previously inserted a `derives` edge on the same pair.

---

### P0-5: RCP Violates `UNKNOWN ≠ EMPTY` Invariant

**File:** `internal/application/app.go:299-322`

`GetContext` returns a plain `error` when Solvent or Conductor is unreachable. It does **not** return a `Context` struct with `epistemic.available=false` and a `reason` field. The MCP adapter propagates this as a tool-call error.

The Phase 8 plan §15 and §16 require:
> When a backend is unavailable: `epistemic.available=false`, `reason="solvent_unavailable"`.

**Impact:** A fresh agent cannot distinguish "no research exists" from "backend is down." It may infer an empty state and redundantly re-enter beliefs, violating the "ideas enter free" doctrine by creating duplicates that the idempotency layer (already broken per P0-1) may or may not catch.

---

### P0-6: No Retirement-Rule Enforcement Before Solvent Discharge

**File:** `internal/application/app.go:327-343`

```go
case "discharge":
    return a.kern.Discharge(ctx, req.ScenarioID, req.BeliefID, req.ObligationKey, req.InstrumentRef, req.PrincipalID)
```

`SubmitDecision("discharge")` calls Solvent's `Discharge` directly. It does **not** look up the pack's `RetirementRules` to verify that the offered evidence class matches the rule's required class.

The Phase 8 plan §10 and §11 explicitly require:
> Coordinator mechanically validates the offered evidence class against the pack rule before invoking Solvent discharge.

**Impact:** Any caller who can reach `SubmitDecision` can retire any debt item with any instrument, bypassing the pack's retirement rules entirely. The UI passes `evidence_class` in the request body but `SubmitDecision` ignores it.

---

### P0-7: Trust UI Discharge Endpoint Bypasses Coordinator Authentication

**File:** `internal/ui/ui.go:104-140`

```go
principalID := "00000000-0000-0000-0000-000000000001"
```

The UI hardcodes the principal ID and does not authenticate the caller. The Coordinator HTTP layer (`coordinator/http/handler.go:42-59`) requires a Bearer token for `/decisions`. The Trust UI calls the Coordinator directly over HTTP but never sends a Bearer token. When Coordinator auth is enabled, Trust UI discharge will return 401. The current `main.go` does not enable the Coordinator HTTP layer at all, so the UI's `HandleDischarge` is the only path — and it calls `app.SubmitDecision` directly with the hardcoded principal.

**Impact:** The human identity boundary is not enforced. Any browser that can reach the UI can discharge debt with the hardcoded operator identity. There is no authentication, no CSRF protection, and no audit of the browser session.

---

### P0-8: Application Layer Exposes Solvent Kernel Mutation Methods Directly

**File:** `internal/application/app.go:327-343`

`App.SubmitDecision` is a public method that accepts arbitrary `DecisionRequest` values and dispatches to `kern.Discharge`, `kern.Promote`, and `kern.RetractCascade`. It is callable from any package that imports `internal/application`, including the UI and any future HTTP handler.

The Phase 8 plan §5 states:
> Coordinator must NOT write directly to Solvent DB.

**Impact:** The `internal/application` package is not a boundary enforcement point; it is a direct Solvent kernel wrapper. Any code path that obtains an `*App` can mutate Solvent authority state without going through the Coordinator's retirement-rule validation or authentication.

---

## P1 — IMPORTANT

### P1-1: No Deduplication of Evidence by Content Hash

**File:** `internal/application/app.go:199-214`

Evidence is inserted with a deterministic ID based on `EntityID(scenarioID, "evidence", contentSHA256)`. However, the `evidence` table in Solvent's schema has no unique constraint on `(scenario_id, belief_id, content_sha256)`. The `ON CONFLICT (id) DO NOTHING` only prevents duplicate `id` values, but if an agent submits the same evidence with a different `local_id` (which does not affect `EntityID` because it hashes `content_sha256`), the second insert is a no-op. If an agent submits different content that happens to map to a different deterministic ID, both are stored. The real risk is that `content_sha256` is client-supplied and never verified server-side.

**Impact:** An agent can flood the ledger with near-duplicate evidence entries by varying non-semantic fields. The `content_sha256` field is trusted as a unique key but is not actually unique in the schema.

---

### P1-2: MCP Stdio Transport Not Wired

**File:** `cmd/argus/main.go:102-124`

```go
fmt.Fprintln(os.Stderr, "argus mcp: starting MCP stdio adapter")
// TODO: wire MCP stdio transport
select {}
```

The `argus mcp` command starts, prints a message, and blocks forever. No MCP stdio transport is implemented. The `mcp/adapter` package exists but is never wired to stdin/stdout in production.

**Impact:** The primary agent-facing interface does not function. Agents cannot actually call `argus.get_context` or `argus.submit_packet` via MCP in the current binary.

---

### P1-3: Verifier Version Check Uses Lexicographic String Comparison

**File:** `internal/application/app.go:163-166`

```go
func versionGTE(version, minVersion string) bool {
    return version >= minVersion
}
```

Semantic version comparison is implemented as lexicographic string comparison. Version `"0.10.0"` compares less than `"0.2.0"` because `'1' < '2'`.

**Impact:** A verifier with version `"0.10.0"` would be rejected when the minimum is `"0.2.0"`. Conversely, `"0.2.0"` would be accepted against a minimum of `"0.10.0"`. The `physics-v1` verifier is hardcoded at `"0.1.0"` so this is not exploitable in the POC, but it is a latent defect in the enforcement logic.

---

### P1-4: Integration Tests Bypass Coordinator and MCP Layers

**File:** `integration_test.go`, `internal/application/app_test.go`

The integration and unit tests instantiate `application.App` and `ui.Server` directly, bypassing the Coordinator HTTP layer, MCP adapter, and authentication. The `coordinator/` package has its own tests with mock clients, but there are no tests that exercise the full path: MCP adapter → Coordinator → Solvent/Conductor REST.

**Impact:** Tests prove that direct DB access works. They do not prove that the documented architecture (agents → MCP → Coordinator → REST → Solvent/Conductor) is intact.

---

### P1-5: Trust UI Is a Separate Module But Is Not Used by `argus serve`

**File:** `cmd/argus/main.go:68-76`, `trust-ui/go.mod`

`cmd/argus` imports `internal/ui`, not `trust-ui`. The `trust-ui/` directory contains a separate Go module with its own `go.mod` and a compiled binary. The `internal/ui` package is a simplified embedded UI. The separate `trust-ui` module is listed in the architecture but is not built, run, or tested by any Taskfile target.

**Impact:** The Trust UI described in the Phase 8 plan (separate module, operator identity from env, Coordinator auth) is not the UI that `argus serve` actually serves. The architectural boundary exists in documentation but not in the developer workflow.

---

### P1-6: No Timeout Propagation to Solvent/Conductor HTTP Clients

**File:** `coordinator/client.go:18-23`

```go
func NewSolventClient(baseURL string) *SolventClient {
    return &SolventClient{
        baseURL:    baseURL,
        httpClient: &http.Client{},
    }
}
```

The HTTP clients used by the Coordinator have no timeout. A slow or hung Solvent/Conductor service will block the Coordinator indefinitely.

**Impact:** Under partial failure, the Coordinator can hang rather than returning an error. The MCP adapter and UI will appear frozen.

---

### P1-7: Evidence Content Hash Is Client-Supplied, Not Server-Computed

**File:** `packet/v1/types.go:44`, `internal/application/app.go:199-214`

Agents supply `ContentSHA256` in the packet. The server stores it verbatim. There is no server-side hash computation or verification that the stored hash matches the actual evidence content.

**Impact:** An agent can claim any hash for any content. The binding between evidence content and its hash is not trusted.

---

## P2 — HARDENING

### P2-1: No CSRF Protection on UI Write Endpoints

`internal/ui/ui.go` exposes `/ui/api/discharge` and `/ui/api/promote` as POST endpoints with no CSRF token, no Origin check, and no SameSite cookie. Any browser page can form-post to these endpoints.

---

### P2-2: Hardcoded Ports in Trust UI and `cmdServe`

`trust-ui/server.go:38` listens on `:8081` hardcoded. `cmd/argus/main.go:50` listens on `:8080` hardcoded. No dynamic allocation or preflight.

---

### P2-3: Physics Verifier Uses Hardcoded Constants

`verifier/physics/v1/verifier.go:13-17` defines `VerifierID`, `VerifierVersion`, and `VerifierHash` as package-level constants. The `VerifierHash` is a static string `"poc-verifier-v0.1.0"` and is never computed from the verifier binary or source.

---

### P2-4: No Pagination on RCP Context

`internal/application/app.go:64-82` returns all beliefs, evidence, and intents for a scenario in a single response. For large scenarios, this will produce unbounded response sizes.

---

### P2-5: `integration_test.go` Uses Port 26257, `Taskfile` Uses 26260

`integrationDB` defaults to `postgres://root@localhost:26257/defaultdb?sslmode=disable`. The Taskfile starts CRDB on `:26260`. A developer running `task fresh` must manually start CRDB on `:26257` for integration tests to pass, or set `ARGUS_TEST_DSN`.

---

### P2-6: No Audit Trail for Failed Discharge Attempts

`kernel.Discharge` returns an error on failure but writes no audit record. An operator cannot distinguish "never attempted" from "attempted and refused."

---

### P2-7: `belief_retirement_proposal` Table Is Created but Unused

`internal/migrations/002_idempotency_retirement.sql` creates `belief_retirement_proposal`. No code in the repository reads from or writes to this table.

---

### P2-8: `Taskfile` Lacks `task down`

There is no `task down` to stop background CRDB and ARGUS processes. `task dev` backgrounds CRDB with `&` and never records the PID. Ctrl-C leaves orphaned CRDB processes.

---

## PORTABILITY / DEVELOPER EXPERIENCE

### Port-Conflict Defect (Critical)

The `task dev` Taskfile starts CRDB SQL on `:26260` and then starts `argus serve` on `:8080`. CRDB's HTTP admin port also defaults to `:8080`. If CRDB's admin port binds before ARGUS, `argus serve` fails to start. There is no preflight, no fallback, and no cleanup.

Additional defects:
- No PID tracking for background CRDB process.
- No signal propagation; Ctrl-C does not kill the background CRDB.
- `task fresh` uses `sleep 3` instead of a readiness probe.
- `task dev` uses `until cockroach node status` but that checks the SQL port; the HTTP/admin port may still be initializing when `argus serve` starts.

### Required Behavior Not Implemented

The review requirements demand:
1. Port preflight before starting any process.
2. Automatic fallback to next available port from a deterministic local range.
3. Reserve all ports before starting services.
4. Propagate selected ports to DSN, listen address, readiness probes, and integration tests.
5. Print resolved topology.
6. Clean up only processes started by this invocation.
7. Handle Ctrl-C / task interruption without orphaned processes.
8. Handle stale PID/data state.

None of this is implemented.

---

## ACCEPTANCE AUDIT

| # | Criterion | Status | Evidence |
|---|-----------|--------|----------|
| 1 | Fresh-agent reconstruction via `argus.get_context` | **FAIL** | MCP stdio transport is not wired (`cmd/argus/main.go:121-123`). RCP returns errors instead of `available=false` (P0-5). |
| 2 | Refusal-before-promotion (Solvent SQLSTATE 23514) | **PASS** | `kernel.Promote` wraps SQLSTATE 23514 as `ErrPromotionBlocked`. Solvent schema enforces `promoted_is_debt_free`. |
| 3 | Solvent-decides-the-gate | **PASS** | UI and MCP call `SubmitDecision("promote")` which calls `kern.Promote`. Solvent DB enforces the gate. |
| 4 | Human-discharged-debt with attributed identity | **PARTIAL** | `Discharge` records `discharged_by` in Solvent. But UI hardcodes the principal ID (P0-7), and no authentication gates the UI path. |
| 5 | Retirement-rule enforcement | **FAIL** | `SubmitDecision("discharge")` calls `kern.Discharge` without checking the pack's `RetirementRules` (P0-6). |
| 6 | Adversarial challenge (`contradicts` edge) | **PARTIAL** | Edge creation works when called directly. But `App.Persist` bypasses Solvent's edge endpoint, and `ON CONFLICT (parent_id, child_id)` silently drops conflicting kinds (P0-2, P0-4). |
| 7 | Dead-end derivation | **PARTIAL** | `cancelLinkedTasks` cancels governance-linked tasks on retraction. But dead-end derivation is not implemented in `GetContext` or the UI; `Insights` and `Debts` pages do not compute dead-end status. |
| 8 | Retraction exercise (Branch A) | **PARTIAL** | `RetractCascade` is called and cancels linked tasks. But the contradicts edge may have been silently dropped (P0-4). |
| 9 | Capability boundary (two tools only) | **PASS** | `mcp/adapter/adapter.go` exposes exactly two tools. `boundary_test.go` verifies no `internal/epistemic` import in MCP. |
| 10 | No silent failures (`UNKNOWN ≠ EMPTY`) | **FAIL** | `GetContext` returns an error on backend failure (P0-5). No `available=false` projection. |
| 11 | Operator identity (missing fails closed, configured reaches audit) | **FAIL** | UI hardcodes principal ID. No authentication. Any browser caller discharges debt as the operator (P0-7). |
| 12 | Idempotency isolation (same content, different scenario) | **PARTIAL** | `contentHash` includes `ScenarioID`, so cross-scenario deduplication is correct. But within-scenario deduplication is broken because the hash excludes evidence/edges/tasks (P0-1). |
| 13 | No UI bypass (all writes through Coordinator) | **FAIL** | `internal/ui/ui.go` calls `app.SubmitDecision` directly, bypassing Coordinator HTTP layer, auth, and retirement-rule validation (P0-8). |
| 14 | No AI language in UI | **PASS** | `internal/ui/ui.go` and `trust-ui/server.go` do not contain the forbidden phrases. |
| 15 | All existing tests pass | **PASS** | `go test -count=1 ./...` passes. `go vet ./...` passes. |
| 16 | Solvent schema ownership intact | **PASS** | `solvent-main/db/` owns the schema. `internal/migrations/` contains only ARGUS-owned tables. `cmdReset` does not apply Solvent migrations (P0-3), but ownership is correct. |
| 17 | No third Solvent change | **PARTIAL** | The implementation does not add a Solvent schema change. But it bypasses the Solvent REST API (P0-2), which is the architectural boundary the plan intended to enforce. |
| 18 | Domain-pack isolation enforced | **PASS** | `domain-pack/bmist/v1` contains all BM-IST vocabulary. `internal/application/app.go` imports `domain-pack` only through the registry. `boundary_test.go` and `.golangci.yml` enforce import boundaries. |
| 19 | Entity-level idempotency | **PARTIAL** | Beliefs, evidence, and edges use deterministic IDs with `ON CONFLICT DO NOTHING`. But the idempotency marker hash is incomplete (P0-1), and evidence deduplication by content hash is not schema-enforced (P1-1). |
| 20 | Separate Solvent/work transactions | **PARTIAL** | `App.Persist` writes to both Solvent and Conductor tables in a single transaction. The plan intended separate Solvent and Conductor transactions via REST calls. |
| 21 | Migration ownership | **FAIL** | `cmdReset` does not apply Solvent migrations (P0-3). |
| 22 | Refusal-first acceptance | **PARTIAL** | Promotion with open debt is refused. But discharge bypasses retirement rules (P0-6), and the UI has no auth (P0-7). |
| 23 | Verifier trust binding | **PASS** | `ArtifactRegistry` is unexported for registration. `RunPhysicsVerifier` is the only trusted path. Input-to-claim binding is checked in `App.Persist` (lines 216-248). Verifier version check exists (P1-3 is a latent defect, not a bypass). |
| 24 | No agent-selected arbitrary verifier | **PASS** | Verifier specs are loaded from the pack at startup. `App.Validate` checks artifact verifier against authorized specs. |
| 25 | CLI and in-process verifier use same implementation | **PASS** | `cmdVerify` calls `verifier.RunPhysicsVerifier` which calls `physicsv1.Run`. The same `Run` function is used in tests. |

---

## ARCHITECTURE VERDICT

**FAIL — architecture/invariant compromised**

The implementation does not preserve the frozen architecture's core invariant:

> CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Specific violations:
1. The `internal/application` package is not orchestration; it is a direct Solvent kernel wrapper that writes to Solvent tables via raw SQL (P0-2, P0-8).
2. The Trust UI bypasses the Coordinator's authentication and retirement-rule enforcement (P0-7).
3. The idempotency contract is broken: distinct packets are silently collapsed (P0-1).
4. RCP does not honor `UNKNOWN ≠ EMPTY` (P0-5).
5. The developer experience is not reproducible: `task fresh` does not produce a working database (P0-3), and `task dev` has a known port-conflict defect with no mitigation.

The Phase 8 dry run cannot honestly demonstrate:
- agent work + adversarial review + human debt adjudication + Solvent authority

because the Solvent authority boundary is not actually exercised by the main binary, and the UI/adjudication path is not authenticated.

---

## IMPLEMENTATION RULE

The following fixes can be implemented without changing the frozen architecture or adding a third Solvent change:

1. **Fix `contentHash`** to include evidence, edges, and tasks (P0-1).
2. **Route `App.Persist` through Solvent/Conductor REST clients** instead of raw SQL (P0-2). This is the intended Coordinator path; the `coordinator/` package already has the clients.
3. **Apply Solvent migrations in `cmdReset`** by embedding or shelling out to `solvent-main/db/` (P0-3).
4. **Change edge uniqueness conflict handling** from `DO NOTHING` to an explicit check that returns an error when `kind` differs from an existing edge (P0-4).
5. **Return `Context` with `available=false`** on backend failure instead of a plain error (P0-5).
6. **Enforce retirement rules in `SubmitDecision("discharge")`** by looking up the pack rule and checking evidence class (P0-6).
7. **Authenticate the Trust UI** by requiring the Coordinator Bearer token or removing the hardcoded principal and requiring the operator identity to be injected server-side with no browser override (P0-7).
8. **Restrict `SubmitDecision`** to be unexported or package-internal, and expose only the Coordinator HTTP handler as the public decision path (P0-8).

No proposed fix requires a third Solvent schema change. The Solvent REST API already exists; the fix is to use it.
