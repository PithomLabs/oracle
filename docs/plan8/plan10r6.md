# Plan 10r5 — Minimal Implementation Checklist (FOUNDATION FREEZE)

**Date:** 2026-09-18 (revision 6 — foundation freeze, implementation begins)
**Source:** Foundation freeze directive — stop review cycle, build the running system
**Architecture:** One process. One database. Two MCP tools. Real authority gates. Real human decision path. Real Trust UI.

---

## Architecture (frozen)

```
OpenCode
   │
   │ MCP stdio (2 tools)
   ▼
┌──────────────────────────────┐
│ ARGUS — one Go binary        │
│                              │
│ serve  = HTTP + Trust UI     │
│ mcp    = MCP stdio adapter   │
│ verify = physics verifier    │
│ migrate = idempotent DB init │
│ reset  = destructive DB init │
└──────────────┬───────────────┘
               │
          CockroachDB
```

Fixed local ports:

```
CRDB SQL     26257
CRDB Admin   8081
ARGUS HTTP   8080
```

Fail clearly when occupied. **No dynamic port allocator.**

---

## Implementation items (ordered by dependency)

### Fix 1 — contentHash covers full canonical packet

**File:** `internal/application/app.go:78-88`

Hash the complete canonical packet: `ScenarioID`, `PacketID`, `Role`, `PackRef`, then beliefs (sorted by LocalID), evidence (sorted by LocalID), edges (sorted by LocalID), tasks (sorted by Title+Description+GovernanceRef). Use `h.Write([]byte{0})` as field separator, `h.Write([]byte{1})` as entity separator.

Tests: `TestContentHashDifferentEvidenceProducesDifferentHash`, `TestContentHashDifferentEdgesProducesDifferentHash`, `TestContentHashCanonicalOrdering`.

### Fix 8 — Strict semantic version comparison

**File:** `internal/application/app.go:163-166`

Replace lexicographic `versionGTE` with exact integer semver parsing. Reject prerelease versions. Use `strconv.Atoi` (not `fmt.Sscanf`) for strict integer parsing of exactly 3 components.

Tests: `TestVersionGTEEqual`, `TestVersionGTEMajorGreater`, `TestVersionGTEMinorGreater`, `TestVersionGTEMinorLess`, `TestVersionGTEPatchGreater`, `TestVersionGTERejectsPrerelease`.

### Fix 3 — Edge kind conflict detection (concurrency-safe)

**File:** `internal/application/app.go:257-264`

After `ON CONFLICT DO NOTHING`, check `RowsAffected()`. If 0, re-read existing edge kind within same transaction. Different kind → error. Same kind → idempotent.

Tests: `TestEdgeKindConflictReturnsError`, `TestEdgeSameKindIsIdempotent`, `TestEdgeDifferentPairsSucceed`.

### Fix 9 — Server-verify evidence hashes + operator_asserted authority

**File:** `internal/application/app.go:198-214`

In `Persist`, for `reproducible_artifact`: verify `artifact_ref` exists, resolve from registry, compare `ArtifactHash` to `ContentSHA256`, enforce verifier version. For `operator_asserted`: **reject from agent packets** — return error. Evidence creation for `operator_asserted` is deferred to trusted demo setup path (not built in this pass).

Tests: `TestReproducibleArtifactRequiresArtifactRef`, `TestReproducibleArtifactRequiresTrustedArtifact`, `TestReproducibleArtifactHashMismatch`, `TestReproducibleArtifactVerifierNotAuthorized`, `TestReproducibleArtifactHappyPath`, `TestOperatorAssertedRejectedFromAgent`.

### Fix 2 — Solvent migration export

**Files:** `solvent-main/migrations/` (new package), `cmd/argus/main.go`

1. Create `solvent-main/migrations/db/` with 10 SQL files (copy from `db/`)
2. Add `IF NOT EXISTS` to bare statements in `001_schema.sql` (4) and `007_service_tables.sql` (8)
3. Create `solvent-main/migrations/migrations.go` with `Apply(embed.FS)` using single `splitStatements` implementation
4. Update `internal/testdb/testdb.go` to call `migrations.Apply()` instead of duplicating splitter
5. Create `cmdMigrate` (idempotent, non-destructive) and update `cmdReset` (destructive, drops tables first)
6. Wire `cmdMigrate` into `main()` switch
7. Replace filesystem-based test helpers with package calls

Tests: `TestCmdMigrateIsIdempotent`, `TestCmdResetDropsAndReapplies`.

### Fix 4 — UNKNOWN ≠ EMPTY in GetContext (per-section availability)

**Files:** `internal/application/app.go`, `internal/epistemic/view.go`, `internal/mcp/adapter.go`

Add `RCPAvailability` with per-section `SectionAvailability` (task, dependencies, snapshot). Each section independently tracks `Available` + `Reason`. `GetContext` never returns error for availability issues — always returns Context with availability metadata.

Tests: `TestGetContextWithUnavailableWorkStore`, `TestGetContextWithUnavailableSolvent`, `TestGetContextWithEmptyScenario`, `TestGetContextNormal`, `TestGetContextAvailabilityIndependence`.

### Fix 5 — Retirement-rule enforcement before Discharge

**Files:** `domain-pack/registry.go`, `domain-pack/bmist/v1/types.go`, `coordinator/validate.go`, `internal/application/app.go`, `cmd/argus/main.go`

1. Add generic types (`RetirementRule`, `VerifierSpec`) to `domain-pack/registry.go`
2. Add `DecoderFunc` registration to `PackRegistry`; `decodePack` uses registered decoder
3. Extend `Pack` interface with typed accessors
4. Update `bmistv1.Pack` to use generic types; remove old `bmistv1.RetirementRule`/`VerifierSpec`
5. Register BM-IST decoder in composition root (`cmd/argus/main.go`)
6. Simplify `coordinator/validate.go` — remove double type-assertions
7. Simplify `coordinator/coordinator.go:GetPackRules` — no JSON round-trip
8. Split `DecisionRequest` into `ExternalDecisionRequest` (no PrincipalID) + `AuthenticatedDecisionCommand` (server-derived)
9. `SubmitDecision` takes `AuthenticatedDecisionCommand`, resolves pack from scenario mapping (`scenarioPackMapping()`)
10. Full enforcement: debt exists → retirement rule matches → qualifying evidence exists → discharge
11. `parsePackRef` uses `@` separator: `pack_id@version`

Tests: `TestDischargeRequiresKnownDebtItem`, `TestDischargeRequiresEvidenceClass`, `TestDischargeWithMismatchedEvidenceClass`, `TestDischargeRequiresPersistedEvidence`, `TestDischargeWithQualifyingEvidence`, `TestDischargeWithRegistryUnavailable`, `TestSubmitDecisionRejectsExternalPrincipalID`.

### Fix 6 — Trust UI authentication

**Files:** `internal/ui/ui.go`, `cmd/argus/main.go`

1. `NewServer` requires `operatorToken string` — panic if empty (fail-closed)
2. `HandleLogin` — validates token, sets HttpOnly cookie
3. `requireAuth` — checks cookie OR Bearer header
4. `requireOrigin` — validates Origin/Host headers for localhost
5. Protect write endpoints with `requireAuth(requireOrigin(...))`
6. Wire `ARGUS_OPERATOR_TOKEN` in `cmdServe` — `log.Fatal` if empty
7. UI constructs `AuthenticatedDecisionCommand` with server-derived principal, never accepts body principal

Tests: `TestDischargeRequiresAuth`, `TestDischargeWithValidToken`, `TestDischargeWithValidCookie`, `TestDischargeRejectsBodyPrincipal`, `TestDischargeRejectsInvalidOrigin`, `TestNewServerRequiresToken`.

### Fix 7 — MCP stdio transport (official SDK)

**Files:** `cmd/argus/main.go`, `internal/mcp/adapter.go`, `go.mod`

1. `go get github.com/modelcontextprotocol/go-sdk/mcp`
2. `cmdMCP` creates `mcp.NewServer`, registers tools via `adapter.RegisterTools(server)`, runs `server.Run(ctx, &mcp.StdioTransport{})`
3. Adapter exposes `RegisterTools(*mcp.Server)` that registers `argus.get_context` and `argus.submit_packet`
4. Remove `select {}` blocking, remove `fmt.Fprintln` message

Tests: `TestMCPToolsRegistered`, `TestMCPGetContextViaSDK`, `TestMCPSubmitPacketViaSDK`.

### Fix 10 — Developer lifecycle (fixed ports, atomic lock)

**Files:** `cmd/argus/main.go`, `Taskfile.yml`

**Fixed ports:** CRDB SQL `26257`, CRDB Admin `8081`, ARGUS HTTP `8080`. All commands default to these. `--db` flag and `ARGUS_TEST_DSN` override.

**`cmdServe` uses retained listener** (no probe-then-close TOCTOU):

```go
ln, err := net.Listen("tcp", *listen)
if err != nil {
    log.Fatalf("serve: listen on %s: %v (port must be free; use task down first)", *listen, err)
}
// ... setup handler ...
srv.Serve(ln)  // not srv.ListenAndServe()
```

**Atomic startup lock** (Taskfile):

```sh
mkdir -p .task
if ! mkdir .task/dev.lock 2>/dev/null; then
  echo "FATAL: another task dev is running. Use 'task down' first."
  exit 1
fi
echo $$ > .task/dev.lock/pid
trap 'rm -rf .task/dev.lock; ...' EXIT INT TERM
```

`mkdir` is atomic — two processes cannot both succeed.

**Stale lock recovery:** Before refusing, check if recorded PID is still alive. If not, `rm -rf .task/dev.lock` and proceed.

**`task dev`** (non-destructive): resolve ports → start CRDB → wait readiness → `argus migrate` → `argus serve` (retained listener) → cleanup on exit.

**`task fresh`** (destructive): resolve ports → start CRDB → wait readiness → `argus reset` → `argus serve` → cleanup.

**`task down`**: kill PID files → remove `.task/dev.lock`.

**Test DSN derivation:** Parse `ARGUS_TEST_DSN` or default to `postgres://root@localhost:26257/defaultdb?sslmode=disable`.

Tests: `TestResolvePortPreferredAvailable`, `TestResolvePortExhaustion`, `TestDevLockPreventsConcurrentStartup`, `TestDevLockStalePIDRecovery`.

### Fix 11 — task down + PID cleanup

**File:** `Taskfile.yml`

```yaml
down:
  desc: Stop background CRDB and ARGUS processes
  cmds:
    - |
      for pidfile in .argus-pids/*.pid; do
        [ -f "$pidfile" ] || continue
        pid=$(cat "$pidfile")
        kill "$pid" 2>/dev/null || true
        rm "$pidfile"
      done
      rm -rf .task/dev.lock
```

### Fix 12 — Integration test CRDB config

**File:** `integration_test.go`, `app_test.go`, `ui_test.go`

Tests default to `postgres://root@localhost:26257/defaultdb?sslmode=disable`. `ARGUS_TEST_DSN` env var overrides. Replace filesystem-based `applySolventMigrations` with `solventmigrations.Apply(ctx, db)`.

### Fix 13 — Tests for every fix

Full suite regression: `go test ./...`, `go test -race ./...`, `go vet ./...`. Acceptance criteria re-evaluation against 25 criteria from `PIVOT_POC_IMPLEMENTATION_PLAN.md`.

---

## What is explicitly deferred

| Item | Reason |
|------|--------|
| Dynamic port allocator | Fixed ports + clear failure sufficient for POC |
| EADDRINUSE retry/re-resolution | Fixed ports: occupied = fatal |
| Scenario→pack persistence (`scenario_pack` table) | One pack exists; server-selected is sufficient |
| Operator-attestation wizard endpoint | Seeded through trusted demo setup path |
| Sophisticated semver | One pack/verifier version; exact comparison sufficient |
| Pagination | Dataset is tiny |
| Elaborate audit framework | Existing state + activity records sufficient |
| Second domain | Portability claim proven by architecture, not by building second domain |
| Stale CRDB process recovery beyond PID check | `task down` + startup lock sufficient |

---

## Acceptance criteria impact

| AC | Criterion | Before | After |
|----|-----------|--------|-------|
| 1 | Fresh-agent reconstruction | FAIL | PASS (Fix 4 + Fix 7) |
| 5 | Retirement-rule enforcement | FAIL | PASS (Fix 5) |
| 6 | Adversarial challenge (contradicts edge) | PARTIAL | PASS (Fix 3) |
| 10 | UNKNOWN ≠ EMPTY | FAIL | PASS (Fix 4) |
| 11 | Operator identity enforced | FAIL | PASS (Fix 6) |
| 12 | Idempotency isolation | PARTIAL | PASS (Fix 1) |
| 15 | All tests pass | PASS | PASS (Fix 13) |
| 21 | Migration ownership | FAIL | PASS (Fix 2) |

---

## Execution order

1. Fix 1 (contentHash) — no dependencies
2. Fix 8 (semver) — no dependencies
3. Fix 3 (edge conflicts) — no dependencies
4. Fix 9 (evidence hashes + operator_asserted) — no dependencies
5. Fix 2 (Solvent migrations + cmdMigrate) — needs Solvent package
6. Fix 4 (UNKNOWN ≠ EMPTY) — needs view.go changes
7. Fix 5 (retirement rules + DTO split) — needs registry/bmist type changes
8. Fix 6 (UI auth) — needs token wiring
9. Fix 7 (MCP SDK) — needs `go get` + adapter rewrite
10. Fix 10 (lifecycle: fixed ports, retained listener, atomic lock) — needs Taskfile + main.go
11. Fix 11 (task down) — needs Taskfile
12. Fix 12 (test CRDB config) — after Fix 10
13. Fix 13 (full test pass) — after all fixes

---

## Constraints

- **One process, one database, two MCP tools.** No REST service resurrection.
- **`CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION`.** Enforced via depguard + compile-time tests.
- **`ARGUS_OPERATOR_TOKEN` has no default.** Missing token → startup failure.
- **`domain-pack` does NOT import `bmistv1`.** Decoder registration from composition root.
- **`operator_asserted` rejected from agent packets.** Human path deferred to trusted setup.
- **`PackRef` never from external request bodies.** Server-derived from scenario mapping.
- **`PrincipalID` never externally representable.** Split DTO + authenticated command.
- **Fixed ports: 26257, 8081, 8080.** Fail when occupied. No dynamic allocation.
- **Retained listener, not probe-then-close.** No TOCTOU race on HTTP bind.
- **Atomic `mkdir` lock for `task dev`.** Stale PID check before refusal.
- **`task dev` is non-destructive.** Uses `argus migrate`. `task fresh` uses `argus reset`.
- **Single `splitStatements` implementation.** In migrations package; test helpers call it.
- **Strict integer semver.** `strconv.Atoi`, exactly 3 components, reject prerelease.
- **Full canonical sort key for contentHash.** Title+Description+GovernanceRef for tasks.
