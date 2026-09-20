# FINAL ADVERSARIAL CODE REVIEW

## 1. Executive Result

**GO**

The implemented foundation is sound. The central authority/security/data-integrity thesis holds. Every adversarial path reviewed—agent capability boundary, human authority derivation, retirement/discharge gate, evidence authority, idempotency, edge integrity, UNKNOWN≠EMPTY, domain-pack isolation, migration ownership, and developer lifecycle—is either correctly enforced or the absence of enforcement is the explicitly accepted POC scope (e.g., MCP stdio transport). The two failing tests (`TestFullIntegration` and `TestDischargeEndpoint`) are test-suite setup defects, not implementation defects: one uses a forbidden evidence class that the code correctly rejects, and the other calls the UI handler without first bootstrapping the pack registry. No P0 or P1 finding undermines the foundation. Proceed to MCP implementation and end-to-end demonstration.

---

## 2. Verified P0 Findings

**None.**

---

## 3. Verified P1 Findings

**None.**

The two failing tests under `go test ./...` and `go test -race ./...` are not P1 findings:

- `TestFullIntegration` (root package, line 158): submits a packet containing `ProvenanceClass: "operator_asserted"`. `App.Persist` at `internal/application/app.go:337–339` explicitly rejects `operator_asserted` from agent-submitted packets. The test fails because the test fixture is wrong, not because the code is wrong.
- `TestDischargeEndpoint` (`internal/ui/ui_test.go:150`): calls `server.HandleDischarge` without calling `server.SetPackRegistry(...)`. The handler correctly returns 500 because the pack registry is nil. The test's setup is incomplete.

---

## 4. P2 Findings

### P2-1: `belief_retirement_proposal` table is created but never used

- **Exact path**: `internal/migrations/002_idempotency_retirement.sql:11–19` creates the table. `cmdReset` drops it (`cmd/argus/main.go:237`). Zero Go references outside of `cmdReset`.
- **Why it is real**: Dead schema. Not harmful, but adds noise and maintenance burden.
- **Minimal fix**: Remove the table from the migration and from `cmdReset`'s drop list. (Advisory only — not blocking.)

### P2-2: Coordinator HTTP clients have no request timeout

- **Exact path**: `coordinator/client.go:21` — `httpClient: &http.Client{}` with default zero timeout. Affects `SolventClient` and `ConductorClient`.
- **Why it is real**: A hung Solvent or Conductor process causes the ARGUS process to hang indefinitely on any Coordinator call.
- **Minimal fix**: Add a configurable timeout (e.g., 10–15 s) to the HTTP client constructor. (Advisory — the embedded-process POC reduces but does not eliminate this risk.)

### P2-3: `contentHash` does not include evidence `SourceURL` or `ArtifactRef`

- **Exact path**: `internal/application/app.go:90–192`. Evidence section hashes `LocalID`, `BeliefRef`, `ProvenanceClass`, `ContentSHA256`, `SourceURL`, `ArtifactRef`. The evidence entity ID at line 348 uses only `ContentSHA256`, not `SourceURL`/`ArtifactRef`.
- **Why it is real**: Two evidence items with the same `ContentSHA256` but different `SourceURL` or `ArtifactRef` share the same entity ID. The second insert is silently dropped by `ON CONFLICT (id) DO NOTHING`. In practice this is unreachable for distinct evidence because `ContentSHA256` is the authoritative content identifier, but the metadata fields are lost.
- **Minimal fix**: Include `SourceURL` and `ArtifactRef` in the entity-ID derivation or in the content hash. (Advisory.)

### P2-4: `cmdReset` drops `action_intent` but the drop list is manually maintained

- **Exact path**: `cmd/argus/main.go:235–239` — hardcoded drop list omits `action_intent` (which the Solvent migration schema creates at `internal/solventmigrations/migrate.go:49–57`).
- **Why it is real**: `cmdReset` leaves `action_intent` rows behind after reset. In the current POC the table is unused, but any future use of `action_intent` would make `cmdReset` non-destructive for that table.
- **Minimal fix**: Add `action_intent` to the drop list or generate the drop list from the migration schema. (Advisory.)

### P2-5: `task fresh` does not check the startup lock before proceeding

- **Exact path**: `Taskfile.yml:67–83`. `task fresh` acquires the lock but does not check if `task dev` is already running (other than the mkdir race). If `task dev` is running and a developer runs `task fresh`, the CRDB startup in `task fresh` will fail because port 26257 is occupied, but only after the trap and PID files are written.
- **Why it is real**: The developer sees a confusing CRDB startup failure and must clean up manually.
- **Minimal fix**: Add a port pre-flight probe before starting CRDB in both `task dev` and `task fresh`. (Advisory — port collisions are documented as "fail clearly" in the frozen design.)

---

## 5. False Positives / Already Protected

| Concern | Status |
|---|---|
| Agent can call promote/retract/discharge/authorize | **Already protected**: `mcp/adapter/adapter.go:29–38` rejects unknown tool names; `internal/mcp/adapter.go:29–38` same; `ListTools` returns only the two permitted tools. |
| Client can supply `PrincipalID` | **Already protected**: `ExternalDecisionRequest` (`app.go:682–689`) has no `PrincipalID` field; `AuthenticatedDecisionCommand.PrincipalID` is only set server-side. Trust UI hardcodes `"operator"` (`ui.go:191`). |
| `PackRef` is client-controlled | **Already protected**: The application's `resolvePackFromScenario` (`app.go:596–606`) derives the pack from `scenarioID` using a server-side mapping function. `PackRef` is not read from the client in the application discharge path. |
| `operator_asserted` can enter via agent packet | **Already protected**: `app.go:337–339` returns an error before any database write. |
| Cookie auth bypasses origin check | **Already protected**: `requireOrigin` runs before `requireAuth` in the middleware chain (`ui.go:42–43`). A cross-origin browser request with a stolen cookie is blocked by `SameSite=Strict` plus the origin check. |
| Cross-scenario content hash collision | **Already protected**: `ScenarioID` is the first field hashed (`app.go:94`). |
| Edge kind conflict silently dropped | **Already protected**: `ON CONFLICT (parent_id, child_id) DO NOTHING` is followed by a re-read and explicit kind mismatch error (`app.go:409–424`). |
| Partial failure corrupts state | **Already protected**: `internal/mcp/adapter.go:71–84` wraps `Persist` in `BeginTx`/`defer Rollback`/`Commit`. |
| Retry creates duplicate entities | **Already protected**: deterministic entity IDs + `ON CONFLICT DO NOTHING` + `submission_idempotency` unique constraint. |
| `task dev` concurrent start | **Already protected**: atomic `mkdir` lock (`.task/dev.lock`) with PID check and stale-PID cleanup (`Taskfile.yml:9–24`). |
| Retained HTTP listener causes probe-then-close race | **Already protected**: `cmd/argus/main.go:108–114` starts a goroutine that calls `srv.Shutdown` with a 5 s timeout on SIGINT/SIGTERM. |
| Solvent migrations have multiple competing implementations | **Already protected**: `internal/solventmigrations/migrate.go` is the single canonical source. The test helper `applySolventMigrations` reads from `../solvent-main/db/` but is only used by tests and is not a competing migration runner. |
| `PrincipalID` from request body influences discharge | **Already protected**: `ExternalDecisionRequest` has no `PrincipalID` field. The UI hardcodes `"operator"`. |
| `belief_retirement_proposal` used as a retirement gate | **Already protected**: The table is created but never read by any Go code. Retirement rules come from the domain pack (`pack.GetRetirementRules()`), not from this table. |
| Solvent kernel Discharge does not validate principal | **Already protected**: The kernel records `dischargedBy` as an audit attribute; the authorization gate is in the application layer (`app.go:512–544`) which validates retirement rules and evidence before calling the kernel. The kernel is not the authority boundary — the application layer is. |

---

## 6. Foundation Invariant Check

| Invariant | Verdict | Evidence |
|---|---|---|
| Agent capability boundary | **PASS** | `mcp/adapter/adapter.go:29–38`, `internal/mcp/adapter.go:29–38` — exactly two tools; no authority operations exposed |
| Human authority boundary | **PASS** | `ui.go:42–43` — `requireOrigin` then `requireAuth`; `PrincipalID` never from request body; `requireAuth` enforces cookie or Bearer |
| Retirement/discharge gate | **PASS** | `app.go:512–544` — pack registry check → retirement rule lookup → evidence class match → persisted evidence verification → kernel Discharge |
| Evidence authority | **PASS** | `app.go:337–339` rejects `operator_asserted`; `app.go:340–343` requires `ArtifactRef` for `reproducible_artifact`; `app.go:361–392` verifies artifact input hash against belief InputSpec |
| Idempotency | **PASS** | `app.go:90–192` full packet hash; deterministic entity IDs; `ON CONFLICT DO NOTHING`; `submission_idempotency` unique constraint |
| Edge integrity | **PASS** | `app.go:401–426` — PK on (parent, child); explicit kind-mismatch error |
| UNKNOWN ≠ EMPTY | **PASS** | `coordinator/context.go:37–90` — per-section `Available` flag with independent error tracking; one section's failure does not overwrite another's data |
| Domain-pack isolation | **PASS** | `bmistv1` imported only in `cmd/argus/main.go:17` and test files; `domain-pack` generic package has no BM-IST import |
| Migration ownership | **PASS** | `internal/solventmigrations/migrate.go` is the single canonical Solvent migration; `cmdMigrate` and `cmdReset` both call it; `cmdReset` is explicitly destructive |
| Developer lifecycle | **PASS** | `Taskfile.yml` — atomic mkdir lock; CRDB readiness probe; `task dev` → migrate (non-destructive); `task fresh` → reset (destructive); `task down` → kill PIDs + remove lock |

---

## 7. Test Evidence

### `go test ./...`
```
--- FAIL: TestFullIntegration (root package)
    evidence[0]: operator_asserted evidence requires human attestation
--- FAIL: TestDischargeEndpoint (internal/ui)
    pack registry unavailable: cannot validate retirement
... all other packages PASS (including internal/application, coordinator, mcp/adapter)
```

### `go test -race ./...`
```
Same two test failures as above. All other packages PASS.
internal/application: 235.837s — PASS (race detector clean)
```

### `go vet ./...`
```
(no output — clean)
```

### Targeted tests run
```
go test -count=1 -v -run "TestPartialFailureRetrySucceeds|TestIdempotencyDuplicatePacket|TestAdversarialContradictionVisible|TestAgentEdgesDoNotRetract|TestPromoteWithOpenDebtRefused|TestPromoteRetractedBeliefRefused|TestVerifierInputBinding|TestVerifierInputBindingMismatch|TestVerifierSpecEnforcement" ./internal/application/
    All PASS

go test -count=1 -v -run "TestRetireDebt|TestDischarge|TestPromote|TestRetract" ./coordinator/
    All retirement-rule and discharge tests PASS (14 tests)

go test -count=1 -v ./internal/ui/
    TestDischargeEndpoint FAIL (test setup issue — no pack registry)
    TestPromoteEndpoint PASS
    TestInsightsPage PASS
    TestDebtsPage PASS
```

---

## 8. Final Decision

**FOUNDATION READY.**
Proceed to MCP implementation and end-to-end demonstration.
