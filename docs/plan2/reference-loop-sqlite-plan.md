# Reference Loop Continuation — Implementation Plan

## Forensic Gap/Status Matrix (Agent A Baseline)

| Component | Status | Gap |
|---|---|---|
| Compilation | PASS | Compiles cleanly |
| Tests | FAIL | Zero test files exist |
| Conductor interface | FAIL | `conductor/client.go` uses raw SQLite DB calls; must use MCP `CallTool` |
| Solvent interface | FAIL | Custom auth middleware in `main.go`; must use real `api.AuthMiddleware` |
| Solvent DB driver | FAIL | `initSolventDB` opens `"pgx"` driver; must use `"sqlite"` |
| Solvent kernel compatibility | BLOCKED | Solvent kernel (`kernel/sql.go`) uses CockroachDB-specific SQL (`::UUID`, `JSONB`, `TEXT[]`, `gen_random_uuid()`, `crdb.ExecuteTx`, `FOR UPDATE`, `ON UPDATE CASCADE`, `DEFERRABLE FK`); narrow SQLite compatibility patch is required |
| Solvent SQLite schema | PASS | `solvent_schema.sql` already uses SQLite-compatible syntax |
| Operation identity | FAIL | Missing `run_id`; must be `deploy:owner/repo:workflow.yml:ref:run_id` |
| Evidence — EFFECT_CONFIRMED | FAIL | Agent records it without verifying GitHub workflow run state |
| Evidence — RecordingFunc label | FAIL | Not labeled as non-external-effect proof |
| Contract pin artifact | FAIL | No `contract.json` or equivalent |
| Config | FAIL | Default Solvent DB path is Cockroach DSN `postgresql://root@localhost:26257/solvent_ref_loop` |
| Conductor schema | PASS | `conductor_schema.sql` matches Conductor's `migrations/001_initial.sql` and `002_dependencies.sql` |

---

## Decisions

1. **Conductor interface**: Replace direct SQLite calls in `conductor/client.go` with MCP stdio JSON-RPC to Conductor's `internal/mcp.Server.CallTool`. The Conductor MCP server is invoked in-process (same binary) via `Server.CallTool(name, args)` rather than via stdio pipe, because the reference-loop binary owns both DBs.

2. **Solvent interface**: Run a Solvent REST API server in-process using real `api.AuthMiddleware`, `api.NewServer`, `audit.New`, `policy.New`, `authority.New`, `ledger.New`. Replace the custom auth wrapper in `main.go:startSolventServer` with `api.AuthMiddleware(keyToPrincipal, handler)`.

3. **Solvent SQLite driver**: Change `initSolventDB` to `sql.Open("sqlite", dsn)` and apply `solvent_schema.sql`.

4. **Solvent kernel SQLite compatibility patch**: Modify Solvent's `kernel/sql.go`, `kernel/kernel.go`, `kernel/authority.go`, `service/audit.go`, `service/policy.go`, `service/authority.go`, `service/ledger.go`, `api/auth.go`, and `api/reads.go` to use SQLite-compatible SQL. This is narrow storage compatibility — same domain logic, different SQL dialect and transaction mechanism. No generic DB abstraction introduced. Stop and report BLOCKED if the scope expands.

5. **Operation identity**: Declare as `deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>`. Bind proposal, authorization, execution, and evidence to this identity. Include `run_id` as an executor input.

6. **Evidence separation**: Participant tables are authoritative. Harness `EvidenceCollector` reads participant tables after execution. `EFFECT_CONFIRMED` is only recorded after querying GitHub API for workflow run state in GitHub mode. RecordingFunc mode records `EXECUTION_ATTEMPTED` labeled `non_external_effect_proof=true`.

7. **Contract pin artifact**: Add `contract.json` in the reference-loop root pinning workflow spec v0.3, role/boundary matrix v0.5, operation identity rule, declaration owner/version, pass criteria, sandbox rule, failure taxonomy, finding arbiter, and BM-IST deferral.

---

## Ordered Task List

### Task 1: Contract Pin Artifact
- Create `contract.json` in `/home/chaschel/Documents/go/oracle/reference-loop/contract.json`
- Fields: `workflow_spec_version`, `role_boundary_matrix_version`, `operation_identity_rule`, `declaration_owner`, `declaration_version`, `pass_criteria`, `sandbox_rule`, `failure_taxonomy`, `finding_arbiter`, `bmist_status`
- Add test: `contract_test.go` — assert file exists, parseable, required fields present and non-empty

### Task 2: Fix Solvent Kernel for SQLite (in `/home/chaschel/Documents/go/solvent-main`)
This is a bounded storage-compatibility patch. No authority semantics change.

**2a. kernel/sql.go** — Replace all CockroachDB-specific SQL:
- Replace `$1::UUID` with `$1` (TEXT placeholders); UUIDs stored as `lower(hex(randomblob(16)))`
- Replace `gen_random_uuid()` with `lower(hex(randomblob(16)))`
- Replace `::STRING`, `::UUID`, `::JSONB`, `::STRING[]`, `::UUID[]` casts with untyped `$N`
- Replace `JSONB` columns with `TEXT`; application code marshals/unmarshals JSON
- Replace `TEXT[]` debt column with `TEXT` storing JSON array; replace `array_remove()` with `json_remove()` or application-level update; replace `array_length()` with `json_array_length()`
- Replace `now()` with `datetime('now')`
- Replace `FOR UPDATE` with no-op (remove the clause)
- Replace `ON UPDATE CASCADE` on composite FK with application-level cascade in kernel
- Replace `DEFERRABLE` FK with plain FK (SQLite does not defer)
- Remove `VECTOR(1024)` and vector index from `db/002_corpus.sql` (not in reference-loop path)
- Replace `::UUID[]` + `ANY($1::UUID[])` with `json_each` or multiple `OR` clauses

**2b. kernel/kernel.go** — Replace `crdb.ExecuteTx` with `sql.Tx` (Begin → Commit/Rollback). Remove SERIALIZABLE retry logic. Keep same function signatures and domain behavior.

**2c. kernel/authority.go** — Same transaction change. Remove `FOR UPDATE` usage. Handle `ON UPDATE CASCADE` behavior in application code where needed.

**2d. service/audit.go, service/policy.go, service/authority.go, service/ledger.go** — Ensure they work with the patched kernel. No semantic changes.

**2e. api/auth.go** — Replace `$1::UUID` with `$1` in `verifyPrincipalActive`.

**2f. api/reads.go** — Replace `parsePGArray()` with JSON-based array parsing for SQLite.

**2g. Validation**: Run Solvent's existing unit tests. Document any failures. If a Solvent test fails because of removed CockroachDB features (VECTOR, FOR UPDATE serialization guarantees), mark as deferred with evidence.

**STOP RULE**: If the SQLite patch requires touching >8 files in kernel/ or introduces a generic DB interface, report BLOCKED.

### Task 3: Fix Solvent REST Server in reference-loop/main.go
- In `startSolventServer`: build `keyToPrincipal` map; wrap handler with `api.AuthMiddleware(keyToPrincipal, handler)` instead of custom middleware
- In `initSolventDB`: change driver from `"pgx"` to `"sqlite"`; apply `solvent_schema.sql`
- In `createTestPrincipal`: insert the test principal with UUID generated via `lower(hex(randomblob(16)))`
- Ensure `keyToPrincipal["test-api-key"]` maps to `testPrincipalID`

### Task 4: Fix Conductor Client to Use MCP
- Replace `conductor/client.go` direct SQLite operations with MCP tool calls
- Wire: `mcp.NewServer(store.NewDB(conductorDBPath), governance.NullReader{}, agentID)` then use `server.CallTool(name, args)`
- Map `CreateTask`, `ClaimTask`, `SubmitTask`, `PostActivity`, `NextTask`, `GetTask` to `conductor_create_task`, `conductor_claim_task`, `conductor_submit_task`, `conductor_post_activity`, `conductor_next_task`, `conductor_get_task`
- Remove raw SQL from `conductor/client.go`

### Task 5: Fix Config
- Change `DefaultConfig.SolventDBPath` from Cockroach DSN to `filepath.Join(".tmp", "solvent.db")`
- Ensure `LoadFromEnv` still supports override

### Task 6: Fix Operation Identity
- In `agent/agent.go` `setupSolvent`: set `ConsequenceParameters` to JSON including `repo`, `workflow`, `ref`, `run_id`
- In `agent/agent.go` `RunHappyPath`: pass `run_id` as executor input
- In `evidence/collector.go`: reconstruct operation ID as `deploy:repo:workflow:ref:run_id`
- In `main.go:VerifyOperationBinding`: update expected op ID format

### Task 7: Fix Evidence Collection
- In `evidence/collector.go` `collectGitHubEvidence`: query GitHub API (`GET /repos/{owner}/{repo}/actions/workflows/{workflow}/runs?event=workflow_dispatch&per_page=1`) and only record `EFFECT_CONFIRMED` if a run exists with matching conclusion
- In `evidence/collector.go` `collectExecutorEvidence`: add field `NonExternalEffectProof: true` when recording from `RecordingFunc`
- In `agent/agent.go`: do not pre-emptively record `EFFECT_CONFIRMED`; leave that to the evidence collector

### Task 8: Add Tests
Create `reference-loop` tests:

**8a. contract_test.go**
- `TestContractPinPresent` — `contract.json` exists and required fields are non-empty

**8b. operation_identity_test.go**
- `TestOperationIdentityDeterministic` — same declared inputs produce identical op ID string
- `TestOperationIdentityIncludesRunID` — op ID format contains run_id

**8c. authorization_test.go**
- `TestRecordingFuncDoesNotBypassAuthorization` — RecordingFunc mode still requires Solvent authorize + execute calls; verify no intent created without authorization
- `TestDeclarationVersionMismatchFailsClosed` — authorize under V1, attempt execute under incompatible V2; assert failure

**8d. happy_path_test.go**
- `TestHappyPathEvidenceSequence` — run the loop and assert evidence package contains all 9 required events in order: TASK_CREATED, TASK_CLAIMED, PROPOSAL_CREATED, AUTHORIZED, EXECUTION_ATTEMPTED, EFFECT_CONFIRMED (or correctly labeled simulated), RESULT_OBSERVED, CONDUCTOR_UPDATED, AGENT_CONTINUED
- `TestHappyPathOperationBinding` — all effect-relevant evidence records share the same operation ID

### Task 9: Run Test Suite
- `cd /home/chaschel/Documents/go/oracle/reference-loop && go test ./...`
- Fix any compilation or test failures
- Run `go vet ./...`

### Task 10: Record Pending Work (Do Not Implement)
Document in `contract.json` under `pending_work`:
- wrong operation
- missing authorization
- stale authorization
- verification unavailable
- declaration/version mismatch
- replay/duplicate
- execution failure
- ambiguous outcome
- termination/revocation
- Conductor cancellation
- Agent misinterpretation
- cross_implementation_operation_identity

### Task 11: BM-IST Deferral
- Mark BM-IST validation as deferred in `contract.json`
- No BM-IST implementation in this pass

---

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Solvent SQLite patch scope creep | Medium | High | Stop rule at >8 files; report BLOCKED |
| CockroachDB-specific features absent in SQLite (VECTOR, FOR UPDATE) | High | Medium | Remove unused features; defer corpus search |
| GitHub API rate limits during test | Low | Medium | Use RecordingFunc for CI; GitHub mode needs token |
| Conductor MCP in-process vs stdio | Low | Low | Use `CallTool` directly; no stdio pipe needed |
| `ON UPDATE CASCADE` composite FK behavior differs | Medium | Medium | Implement cascade in kernel application code |

---

## Validation Commands

```bash
cd /home/chaschel/Documents/go/oracle/reference-loop
go test ./...
go vet ./...

cd /home/chaschel/Documents/go/solvent-main
go test ./...
go vet ./...
```

---

## Next Recommendation

After happy path passes:
1. Implement adversarial scenarios one at a time
2. Add conformance test matrix
3. Validate BM-IST deferred work separately
4. Review and upstream Solvent SQLite compatibility patch if viable
