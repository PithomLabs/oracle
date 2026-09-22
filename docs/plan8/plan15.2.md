# Plan 15.2 — Final Freeze-Gate Items

## Scope

Four freeze-gate items from review12.1.md. No new bookkeeping, no new MCP tools, no authority model changes.

---

## Item 1: DB-backed MCP `argus.get_context` transport test

**Why:** The application-layer `TestGetContextReadsBackPersistedState` proves the read path, but no test verifies the agent-facing JSON serialization through `adapter.HandleTool("argus.get_context")`. An `omitempty` field drop or marshaling error would evade all existing tests.

**What:** Add `TestMCPGetContextReadsBackPersistedState` to `internal/mcp/adapter_test.go`.

**Steps:**

1. Create a new test function that uses a real CockroachDB (via the same `testDB` pattern as `app_test.go`, but importing `solventmigrations` and `migrations`).

2. Persist the same 3-packet state as `TestGetContextReadsBackPersistedState`:
   - P1 (work): 1 belief + 1 evidence
   - P2 (adversarial): 1 belief + 1 evidence
   - P3 (work): 1 edge (contradicts) + 1 witness belief

3. Call `adapter.HandleTool(ctx, "argus.get_context", json.Marshal(map[string]{"task_id": taskID}))`.

4. JSON-decode the returned `interface{}` into a `map[string]interface{}` (the raw JSON-RPC result, not a typed struct).

5. Assert exact fields through the decoded map:
   - `task.id`, `task.project_id`
   - `snapshot.beliefs` — count=3, each has `id`, `claim`, `claim_type`, `status`, `debt`, `origin_packet_id`
   - `snapshot.evidence` — count=2, each has `belief_id`, `provenance_class`, `content_sha256`
   - `snapshot.edges` — count=1, has `parent_id`, `child_id`, `kind`
   - `availability.task.available`, `availability.snapshot.available`

6. Verify that a non-existent task_id returns `availability.task.available = false`.

**Key concern:** The test must exercise the JSON marshal path in `handleMCPRequest` (cmd/argus/mcp_stdio.go), not just `app.GetContext`. The `adapter.HandleTool` returns a typed `*application.Context` — the JSON serialization happens when `cmd/argus/mcp_stdio.go` marshals it. So the test should:
- Call `adapter.HandleTool` to get the result
- Marshal the result to JSON (`json.Marshal(result)`)
- Unmarshal back to `map[string]interface{}`
- Assert fields through the decoded map

This proves the `Context` struct marshals correctly through JSON.

---

## Item 2: Fix the seven `view_test.go` failures

**Why:** Seven tests fail with `relation "belief" does not exist` because `testDB` connects to `defaultdb` without applying any migrations. The freeze requires "no unexplained failing test."

**What:** Fix `internal/epistemic/view_test.go:testDB` to create an isolated test database with the Solvent schema applied, matching the pattern used by `app_test.go` and `ui_test.go`.

**Steps:**

1. Rewrite `testDB` in `view_test.go` to:
   - Open admin connection to CockroachDB
   - Create a unique database `epistemic_<testname>_test`
   - Open a connection to the new database
   - Call `solventmigrations.Apply(ctx, db)` to create the Solvent tables
   - Register cleanup (drop the database)
   - Return the new connection

2. Add imports: `fmt`, `sort`, `strings`, `path/filepath`, `solventmigrations "github.com/PithomLabs/oracle/internal/solventmigrations"`.

3. The existing `applySolventMigrations` and `splitStatements` helper functions in `ui_test.go` that read from `solvent-main/db/` are the old path. The `view_test.go` fix should use `solventmigrations.Apply` directly (the single canonical source).

4. After the fix, all 7 tests should pass. If any still fail due to schema differences (e.g., `origin_packet_id` not in the base Solvent schema), the `solventmigrations` package already handles this.

---

## Item 3: Assert P2 initial debt through Persist → GetContext

**Why:** Review12.1.md's §3 says: "confirm P2's belief carries the pack initial_debt set through Persist→GetContext — or record it as untested-at-freeze with the fail-open trigger."

**Decision:** Document as limitation. Do NOT wire `CompileDebt` into `Persist()`.

**Current production semantics:**

```
packet debt → Persist() → stored belief debt
```

NOT:

```
packet debt + pack initial_debt → Persist()
```

`CompileDebt` (coordinator/validate.go) unions `packInitialDebt` with `beliefDebt`, but is never called from the production MCP path. `Persist` stores `b.Debt` directly from the packet.

**Action:**

1. The existing assertion in `TestGetContextReadsBackPersistedState` already verifies P2's round-trip:
   - P2 submits `Debt: nil`
   - Persisted debt is empty `[]`
   - GetContext returns empty debt
   This proves the current path faithfully round-trips the packet. It does NOT say every adversarial belief should be debt-free.

2. Add a focused test `TestPersistDoesNotApplyPackInitialDebt` in `app_test.go` that:
   - Creates a packet with `Debt: nil` and `PackRef: "bmist@1.0.0"`
   - Persists it
   - Reads back via GetSnapshot
   - Asserts debt is `[]` (empty)
   - Documents: "If CompileDebt is wired in the future, this test should change to expect the bmist@1.0.0 initial_debt (6 items)."

3. Update the P2 debt comment in `TestGetContextReadsBackPersistedState` to match the user's approved wording:

```go
// P2 submits debt=nil → persisted debt is empty → GetContext returns empty debt.
// This proves the current path round-trips the packet faithfully.
// CompileDebt / pack initial_debt is not applied during production persist.
// This is a known semantic boundary, not a bug.
// Revisit only if a concrete research-cycle failure demonstrates that a belief
// can enter the epistemic ledger without mandatory debt that the active pack requires.
if len(b2.Debt) != 0 {
    t.Errorf("P2 debt = %v, want [] (CompileDebt not in persist path)", b2.Debt)
}
```

---

## Item 4: Record deferred limitations (no implementation)

Add a section to a new file `docs/plan8/freeze-record.md` documenting:

### 4a. Root legacy MCP adapter

```text
Production MCP adapter: internal/mcp/ (internal/mcp/adapter.go).
Root mcp/adapter/ is legacy/non-production, depends on coordinator.Coordinator
instead of application.App. Does not pass context.Context to handlers.
Trigger: any contributor importing mcp/adapter or any test referencing it as
the active adapter.
```

### 4b. Duplicated test migration helpers

```text
Known maintenance risk: test migration setup is duplicated across packages:
  - app_test.go: uses solventmigrations.Apply (internal canonical)
  - ui_test.go: uses solventmigrations.Apply (fixed this cycle)
  - view_test.go: uses solventmigrations.Apply (fixed this cycle)
  - ui_test.go also retains applySolventMigrations reading from solvent-main/db/
Centralizing on solventmigrations.Apply is correct; the old applySolventMigrations
in ui_test.go is dead code but not removed to avoid scope creep.
Trigger: migration divergence or repeated test-schema failure.
```

### 4c. Network exposure limitation

```text
Network threat model: This POC binds to localhost only. No TLS, no
authentication on the Trust UI or MCP stdio transport. The MCP stdio
transport is process-local (stdin/stdout). The Trust UI operator token
(ARGUS_OPERATOR_TOKEN) is the only access control.
Trigger: any deployment beyond local single-user loopback.
```

### 4d. Other deferred limitations

```text
- CompileDebt not wired into production persist path (debt behavior documented above)
- Pack.InitialDebt not exposed via domainpack.Pack interface
- No packet-ID replay detection (idempotency by content hash only)
- No context snapshots (one snapshot per GetContext call)
- No review-coverage persistence
- No epistemic-kind field
- No agent attestation
```

---

## Verification

After all changes:

```
go vet ./...
go build ./...
go test ./...
```

Expected results:
- `go vet`: clean
- `go build`: clean
- `go test`: all green (the oracle integration test at `oracle/TestFullIntegration` may still fail — it's a pre-existing timeout issue unrelated to these changes; disposition it if needed)

## Item 5: Manual smoke instructions in freeze record

The manual smoke is a human-run final acceptance check, not an automated test. Include instructions in `docs/plan8/freeze-record.md`:

```text
MANUAL FRESH-AGENT SMOKE TEST

Prerequisites:
- CockroachDB available
- ARGUS binary built (go build ./cmd/argus)
- OpenCode configured with ARGUS MCP

Steps:

1. Start:
   ./argus serve

2. Start a fresh OpenCode process with no prior ARGUS conversation state.

3. Connect OpenCode to the production ARGUS MCP server.

4. Call:
   argus.get_context
   using the seeded/current task ID.

5. Verify:
   - task is returned
   - existing beliefs are returned
   - exact claim text is present
   - debt is present
   - evidence is present
   - edges are present
   - origin_packet_id is present where applicable
   - availability reports sections as available

6. Submit one small valid Work packet through:
   argus.submit_packet

7. Call:
   argus.get_context
   again from the fresh agent.

8. Verify the newly submitted belief is visible with its correct
   origin_packet_id.

9. Verify /ui/insights and /ui/debts render the resulting state.

Record:
- date
- ARGUS build/version
- agent harness/model
- PASS/FAIL
- any observed discrepancy
```

---

## Freeze checklist

After these items complete:

```
FREEZE GATE
[ ] P1: MCP-transport read-back test (HandleTool + JSON marshal, exact fields)
[ ] P1b: manual fresh-agent smoke (human-run, instructions in freeze record)
[ ] Disposition the 7 view_test.go failures → FIXED (testDB migration setup)
[ ] Confirm P2 belief debt behavior (CompileDebt not wired, empty debt asserted)
[ ] Legacy adapter identified as non-production (freeze record)
[ ] Duplicated migration helpers recorded as maintenance risk (freeze record)
[ ] Network exposure limitation documented (freeze record)
[ ] Known limitations register (freeze record)
[ ] Manual smoke instructions recorded (freeze record)
```
