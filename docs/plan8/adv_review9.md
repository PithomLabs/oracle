# ARGUS — ADVERSARIAL REVIEW 9
## Final Pre-Freeze Adversarial Code Review

**Date:** 2026-09-22
**Reviewer:** Senior Go Architect (adversarial)
**Scope:** Final read-path repair, edge validation, authority-surface cleanup, freeze-critical read-back test, MCP fast validation test, UI semantic repair, test/database corrections, UUID normalization, `validateEdgeScenario` UUID comparison
**Constraint:** DO NOT modify code. DO NOT redesign ARGUS. DO NOT add features. DO NOT reopen deferred bookkeeping unless a concrete integrity failure makes the freeze impossible.

---

## 1. CURRENT CHANGE SET UNDER REVIEW

Review the actual current repository, especially these latest changes:

### A. Read-path repairs
- `internal/epistemic/view.go`
  - `scanBelief()` now scans nullable `origin_packet_id`
  - `GetAllBeliefs()` now selects `origin_packet_id`

- `internal/work/store.go`
  - `ListByProject()` now selects `origin_packet_id`

### B. Edge validation repair
- `internal/application/app.go`
  - `validateEdgeScenario()` UUID comparison changed to case-insensitive comparison via `strings.EqualFold`

### C. Authority-surface test cleanup
- `TestAgentCannotPromote`
- `TestAgentCannotDischargeDebt`
- `TestAgentCannotRetractBelief`

  were removed/replaced because they tested the wrong application-layer boundary.

- `TestMCPDoesNotExposeAuthorityTools`
  now verifies the production MCP surface.

### D. Freeze-critical read-back test
- `TestGetContextReadsBackPersistedState`
  now uses:
    - Packet P1
    - Packet P2
    - Packet P3
    - two packet origins / agents
    - beliefs
    - evidence
    - cross-packet contradicts edge
    - exact content assertions
    - per-object origin assertions
    - GetContext reconstruction

### E. MCP fast validation test
- renamed to:
  `TestMCPValidationRejectsMalformed_FastUnit`
- explicitly documented as nil-DB unit coverage, not production integration coverage.

### F. UI semantic repair
- Tasks table:
    `Agent` → `Claimant`
- Recent Agent Submissions keeps `Agent`.

### G. Test/database corrections
- UI test DB setup now uses internal Solvent migrations so provenance schema exists.
- A pre-existing `==` comparison bug in `app_test.go` was fixed.
- `normalizedUUID` helper was added for EntityID↔CRDB UUID format comparison.

### H. Earlier already-implemented guarantees remain in scope:
- production MCP invokes canonical packet validation
- agent identity is required
- role consistency enforced
- local_id global uniqueness
- claim_type validation
- edge kind validation
- reference validation
- pack_ref uses `@`
- packet_submission inserted before provenance-dependent rows
- packet persistence is atomic
- origin_packet_id provenance
- task project preflight
- exact retry idempotency
- GetContext exposes beliefs/evidence/edges
- exactly two MCP tools

---

## 2. PRIMARY FREEZE QUESTION

Answer:

> Can a fresh AI agent submit valid research through the production MCP path,
> have that research committed atomically with correct provenance, and then
> reconstruct the exact committed research state through the production
> `argus.get_context` interface?

The desired answer is:

    YES

Anything preventing that is a freeze blocker.

---

## 3. WRITE → COMMIT → READ-BACK

Trace and verify the complete path:

    argus.submit_packet
      →
    validation
      →
    BeginTx
      →
    Persist
      →
    Commit
      →
    argus.get_context
      →
    GetContext
      →
    GetSnapshot
      →
    beliefs/evidence/edges/task

Verify actual production code, not only tests.

Specifically verify:

- all writes use one transaction
- provenance is committed before dependent FK rows
- rollback is complete
- GetContext uses committed state
- read errors are not silently converted into an apparently empty snapshot
- availability metadata correctly distinguishes unavailable from empty
- no in-process cache is masking DB/read-path problems

IMPORTANT:

A read failure must not appear to the agent as:

    beliefs = []

when the real state exists.

`UNKNOWN != EMPTY` remains a hard invariant.

---

## 4. scanBelief REVIEW

Inspect the current `scanBelief()` implementation.

Verify:

SELECT:

    id
    claim
    claim_type
    status
    debt
    final_truth
    origin_packet_id

matches Scan destination count and order EXACTLY.

Verify nullable `origin_packet_id` handling is correct.

Test/read both:

A. non-null origin_packet_id
B. NULL origin_packet_id for legitimate seed/human-created objects

Do not accept a scanner that works only for packet-created beliefs.

---

## 5. GetAllBeliefs / ListByProject REVIEW

Verify every SELECT/Scan pair modified by the latest change.

For:

- GetByID
- ListAll
- ListByProject
- GetTasksByGovernanceRef
- GetAllBeliefs
- GetSnapshot

verify:

    SELECT column count == Scan destination count
    SELECT order == Scan order
    nullable fields handled correctly

Pay special attention to the exact position of `origin_packet_id`.

Do not trust the implementation report.

---

## 6. GetContext END-TO-END TEST QUALITY

Inspect `TestGetContextReadsBackPersistedState`.

Verify that it actually proves:

P1:
    work agent
    belief
    evidence

P2:
    adversarial agent
    different belief
    evidence

P3:
    edge between P1 and P2 beliefs

Then verify:

- exact belief count
- exact belief IDs
- exact claim text
- exact status
- exact debt
- exact origin_packet_id
- P1 origin != P2 origin
- exact evidence fields
- exact edge parent
- exact edge child
- exact edge kind
- task ID
- project/scenario identity
- availability metadata

Verify that the test calls the production `argus.get_context` MCP path if
the freeze contract requires agent-facing verification.

If it only calls `app.GetContext`, classify that as a freeze-test gap.

Also inspect P2 debt expectations.

Do not assume empty debt unless actual Domain Pack/compiler behavior produces
empty debt. Verify what the current implementation actually does.

---

## 7. READ ERROR PROPAGATION

This is a critical regression check.

Trace:

    GetSnapshot error
        →
    GetContext

Verify the system does NOT transform:

    SQL scan/query failure

into:

    successful context with empty beliefs/evidence/edges

The agent must be told that context is unavailable when the read fails.

Check `Availability` semantics.

---

## 8. MCP AUTHORITY SURFACE

Inspect `TestMCPDoesNotExposeAuthorityTools`.

Verify it asserts the EXACT production MCP tool set is:

    argus.get_context
    argus.submit_packet

and nothing else.

Do not rely merely on:

    "promote not found"

The complete tool set must be exact.

Verify this is the actual production adapter, not a mock.

Also inspect the old removed tests and confirm they were genuinely
application-layer misplacements.

Do not reintroduce application-level authentication into `SubmitDecision`
unless an actual current authority-path failure is demonstrated.

---

## 9. AUTHORITY THREAT MODEL

Verify current code and docs make this distinction explicit:

Agent authority:
    bounded by MCP tool surface

Human/control authority:
    existing operator/control path

Network authority:
    POC assumes no untrusted direct network access to internal authority
    endpoints

Do NOT implement network hardening.

This is a known limitation with a trigger:

    non-loopback deployment
    OR
    multiple untrusted actors

Verify the limitation is documented rather than silently assumed.

---

## 10. UUID NORMALIZATION CHANGE

Review the new `normalizedUUID` / UUID comparison changes.

Determine:

- why the normalization was required
- whether it changes identity semantics
- whether it changes EntityID determinism
- whether it changes equality only, not stored identity
- whether case-insensitive UUID comparison is sufficient
- whether it can accidentally make two semantically different identifiers equal

Verify tests cover uppercase/lowercase UUID representations.

Do NOT convert deterministic EntityID to UUIDv7.

Do NOT reopen UUID-generation design.

---

## 11. validateEdgeScenario UUID COMPARISON

Inspect the `strings.EqualFold` change.

Verify it only addresses textual UUID representation and does not weaken:

- scenario isolation
- endpoint existence
- self-edge rejection
- canonical contradicts target requirement

Test or reason through:

same UUID, different case
different UUID
missing belief
cross-scenario belief
self-edge

---

## 12. UI CHANGES

Verify:

Tasks table:

    Claimant

means `current_agent`.

Recent Agent Submissions:

    Agent

means packet submission provenance.

Make sure the labels are not accidentally swapped.

Verify `Agent` still appears correctly where packet provenance is being shown.

Do not add new provenance UI.

---

## 13. DATABASE / MIGRATION TEST FIX

The report says UI tests were corrected to use internal Solvent migrations
because the external schema lacked `origin_packet_id`.

Verify this did NOT introduce:

- duplicated schema ownership
- production-only migration behavior differing from tests
- migration ordering problems
- hidden test-only tables
- stale external schema assumptions

The current rule must remain:

    Solvent owns Solvent migrations.

Do not move ownership.

---

## 14. ENTITY ID TEST FIX

Review the `normalizedUUID` helper and the fixed `==` bug.

Verify this was actually a test representation problem rather than a production
identity change.

Confirm:

    EntityID(same inputs) == same UUID

still holds.

Confirm:

    EntityID(different claim) != EntityID(original claim)

still holds.

This must remain deterministic.

---

## 15. FULL CURRENT TEST SUITE

Run or inspect:

    go vet ./...
    go build ./...
    go test ./...

If the application suite is slow, report actual results rather than skipping
it.

For any failure:

- determine whether it is pre-existing
- determine whether the latest changes caused it
- do not accept "pre-existing" without checking

The previously failing authority tests should no longer be present if Plan
15.1 was executed as specified.

---

## 16. MANUAL SMOKE READINESS

Verify that the code is now ready for the final manual smoke test:

    ./bin/argus serve

Fresh OpenCode process:

    argus.get_context
    argus.submit_packet

The fresh process must not rely on conversational history or in-process
state from the Work Agent.

Do NOT perform the manual smoke test unless requested; determine whether the
code is ready for it.

---

## 17. BOOKKEEPING FREEZE

The following are OUT OF SCOPE and must remain deferred:

- packet-ID replay detection
- context snapshot IDs
- review-coverage persistence
- epistemic-kind field
- agent attestation
- additional provenance UI
- graph visualization
- corpus/vector search
- new MCP tools
- additional authority layers
- network hardening

Do not recommend implementing them unless a concrete current failure makes
one necessary.

The governing rule is:

    concrete failure → smallest fix
    no concrete failure → defer

Research methodology and EBP scientific obligations remain completely
unfrozen.

---

## 18. REQUIRED OUTPUT

Return:

### EXECUTIVE VERDICT

One of:

    READY FOR FINAL SMOKE TEST
    NOT READY — P1/P0 REMAINS
    NOT READY — TEST/VERIFICATION GAP ONLY

### WRITE PATH

Can a valid packet commit?

### READ PATH

Can the committed packet be reconstructed exactly through production
`argus.get_context`?

### AUTHORITY

Is the agent MCP surface exactly the two intended tools?

### PROVENANCE

Does every packet-created belief/evidence/edge/task retain correct
origin_packet_id?

### TEST QUALITY

Identify any tests that can pass while the actual production path is broken.

### REMAINING FINDINGS

Severity:

P0 = trust/integrity blocker
P1 = correctness blocker
P2 = meaningful but non-blocking
P3 = documentation/maintainability

For every finding give:

- exact file/function
- concrete scenario
- why the current tests miss it
- smallest correction

### BOOKKEEPING CHECK

Explicitly state whether this review found any justified reason to add new
accounting infrastructure.

Default answer should be NO.

### FINAL QUESTION

> Is ARGUS now boring enough to freeze so the project can return to
> substantive BM–IST–AS research?

Do not modify code.
Do not propose architecture expansion.
Prefer concrete runtime evidence over speculation.

---

# ADVERSARIAL REVIEW 9 — FINDINGS

## EXECUTIVE VERDICT

**NOT READY — TEST/VERIFICATION GAP ONLY**

---

## WRITE PATH

**YES** — a valid packet commits atomically with correct provenance.

Verified in production code (`internal/application/app.go:353-606`):

- `handleSubmitPacket` opens a single `BeginTx`, calls `app.Persist`, then `tx.Commit()` with `defer tx.Rollback()`.
- `Persist` writes in dependency order: `packet_submission` first (FK anchor), then beliefs, evidence, edges with `edge_provenance`, tasks, then `submission_idempotency` last.
- `ON CONFLICT (id) DO NOTHING` on beliefs/evidence/tasks makes retries idempotent.
- `origin_packet_id` is committed on every packet-created object before dependent rows.
- Rollback is complete on any `Persist` error (the caller rolls back).
- DB-backed tests (`TestPartialFailureRetrySucceeds`, `TestMCPAtomicRollback`, `TestIdempotencyDuplicatePacket`, `TestMCPSamePacketIDDifferentContent`) confirm zero-object survival on rollback and no duplicates on retry.

No write-path integrity failure found.

---

## READ PATH

**YES for the application layer; GAP for the MCP transport layer.**

Verified in production code:

- `app.GetContext` (`internal/application/app.go:647-689`) never returns `nil, error` for availability issues. It always returns a `*Context` with per-section `Availability` metadata.
- `GetSnapshot` (`internal/epistemic/view.go:74-172`) returns errors on query/scan failure; `GetContext` catches them and sets `available=false` with a reason string.
- No in-process cache exists in `App`. Every `GetContext` call hits the database.
- `UNKNOWN != EMPTY` is preserved: a failed snapshot returns `Snapshot: &Snapshot{}` with `Availability.Snapshot.Available=false`, not `beliefs=[]` masquerading as a successful empty state.

**Freeze-test gap:** `TestGetContextReadsBackPersistedState` (`internal/application/app_test.go:1794-2022`) calls `app.GetContext(ctx, taskID)` directly. It does **not** call the production MCP path (`adapter.HandleTool("argus.get_context", ...)`) through the stdio JSON-RPC transport. If the MCP adapter's JSON marshaling of `Context`/`Snapshot` were broken (e.g., missing fields, type mismatches, `omitempty` dropping evidence), this test would still pass. The freeze contract requires agent-facing verification through the production MCP interface.

---

## AUTHORITY

**YES** — the agent MCP surface is exactly the two intended tools.

Verified in production code (`internal/mcp/adapter.go:29-38`):

```go
switch name {
case ToolGetContext:   // "argus.get_context"
case ToolSubmitPacket: // "argus.submit_packet"
default:
    return nil, fmt.Errorf("unknown tool: %s", name)
}
```

`TestMCPDoesNotExposeAuthorityTools` (`internal/mcp/adapter_test.go:236-271`) asserts the exact tool set and verifies forbidden tools (`argus.promote`, `argus.retract`, `argus.discharge`, `argus.authorize`) return errors. The old `TestAgentCannotPromote`/`TestAgentCannotDischargeDebt`/`TestAgentCannotRetractBelief` were correctly removed — they tested an application-layer boundary that does not exist. The authority boundary is at the MCP tool surface.

---

## PROVENANCE

**YES** — every packet-created belief/evidence/edge/task retains correct `origin_packet_id`.

Verified:

- `Persist` inserts `packet_submission` before any provenance-dependent rows.
- Beliefs: `INSERT INTO belief (...) origin_packet_id VALUES (... pkt.PacketID)`.
- Evidence: `INSERT INTO evidence (...) origin_packet_id VALUES (... pkt.PacketID)`.
- Tasks: `INSERT INTO conductor_task (...) origin_packet_id VALUES (... pkt.PacketID)`.
- Edges: `INSERT INTO edge_provenance (parent_id, child_id, origin_packet_id)` records the packet origin for every edge.
- `scanBelief` uses `sql.NullString` and only sets `OriginPacketID` when `Valid` is true, correctly preserving NULL for seed/human-created objects.
- `GetAllBeliefs` and `GetSnapshot` both select `origin_packet_id`.
- `ListByProject`, `GetByID`, `ListAll`, `GetTasksByGovernanceRef` all select `origin_packet_id`.

`TestGetContextReadsBackPersistedState` verifies per-object distinct origins: `P1.origin != P2.origin`, and `TestMCPHappyPath` verifies `origin_packet_id` on beliefs, evidence, and `edge_provenance`.

---

## TEST QUALITY

Tests that can pass while the actual production path is broken:

1. **`TestGetContextReadsBackPersistedState`** (P1 gap) — passes if the application-layer `GetContext` works but the MCP stdio JSON transport marshal/unmarshal is broken. No test traverses `tools/call` → `HandleTool("argus.get_context")` → `json.Marshal(result)` with a real DB.
2. **`TestMCPValidationRejectsMalformed_FastUnit`** (by design) — uses nil DB. If `Persist` had a DB-specific bug (wrong column name, FK violation), this test would still pass. The DB-backed tests (`TestMCPRejectsMalformedWithRealDB`, `TestMCPHappyPath`) cover production validation.
3. **`internal/epistemic/view_test.go`** (pre-existing) — all 7 tests fail with `ERROR: relation "belief" does not exist` because they connect to `defaultdb` without applying Solvent migrations. They do not exercise the production read path.

---

## REMAINING FINDINGS

### P1 — correctness blocker (test gap, not code bug)

- **File/function:** `internal/application/app_test.go:1794` — `TestGetContextReadsBackPersistedState`
- **Concrete scenario:** The MCP stdio transport (`cmd/argus/mcp_stdio.go:115-147`) marshals the `*Context` result to JSON via `json.Marshal(result)`. If `Context`, `Snapshot`, or any nested view struct had an incorrect `json` tag, an `omitempty` that drops a required field, or a type that marshals unexpectedly, the agent-facing `argus.get_context` response would be wrong while all application-layer tests pass.
- **Why current tests miss it:** No test calls `adapter.HandleTool("argus.get_context", ...)` with a real DB and inspects the marshaled JSON response through the stdio transport.
- **Smallest correction:** Add a DB-backed test in `internal/application/app_test.go` or `internal/mcp/adapter_test.go` that: (1) persists P1/P2/P3 packets, (2) calls `adapter.HandleTool("argus.get_context", args)` with the task_id, (3) unmarshals the result JSON, and (4) asserts exact belief count, IDs, claims, origins, evidence count, edge kind/parent/child, and availability metadata.

### P2 — meaningful but non-blocking

- **File:** `mcp/adapter/adapter.go` (legacy package at repo root)
- **Concrete scenario:** Two MCP adapters coexist: `internal/mcp/adapter.go` (production, uses `application.App`) and `mcp/adapter/adapter.go` (legacy, uses `coordinator.Coordinator`). The production binary (`cmd/argus/main.go:196`) uses only `internal/mcp.NewAdapter`. The legacy package is dead code for the production path but could mislead contributors into importing or testing the wrong one.
- **Why current tests miss it:** Both packages have passing tests; neither test references the other.
- **Smallest correction:** Add a file-level comment to `mcp/adapter/adapter.go` marking it deprecated/legacy, or remove it if no consumer exists.

- **File:** `internal/epistemic/view_test.go`
- **Concrete scenario:** All tests connect to `defaultdb` without applying Solvent migrations. They fail immediately with `ERROR: relation "belief" does not exist`. This is a pre-existing failure unrelated to the current change set.
- **Why current tests miss it:** They never reach the code under test.
- **Smallest correction:** Apply `solventmigrations.Apply` (or read from `solvent-main/db/`) in `testDB` before running queries.

- **File:** `internal/application/app_test.go` and `internal/ui/ui_test.go`
- **Concrete scenario:** Test helpers duplicate solvent migration application logic (`applySolventMigrations` in `app_test.go:928-955`, `integApplySolventMigrations` in `integration_test.go:69-110`, `applySolventMigrations` in `ui_test.go:68-91`). Each reads from `solvent-main/db/` via a different relative path. If the Solvent schema evolves, these three copies can drift.
- **Why current tests miss it:** They happen to work because `solvent-main/db/` exists and contains the right files.
- **Smallest correction:** Centralize on `solventmigrations.Apply` (the internal package already used by `ui_test.go:58`) or on a single shared helper.

### P3 — documentation/maintainability

- **File:** `internal/solventmigrations/migrate.go`
- **Issue:** The inline schema is documented as "a minimal POC subset of the full Solvent migration set." If the full `solvent-main/db/` schema diverges from this subset, production (which uses `solventmigrations.Apply` during bootstrap) and tests (which may read from `solvent-main/db/`) will see different table structures. This is a known POC trade-off but should be explicitly called out in the freeze decision record.

---

## BOOKKEEPING CHECK

**NO** — this review found no justified reason to add new accounting infrastructure.

The current provenance spine (`origin_packet_id` on belief/evidence/task, `edge_provenance` table), idempotency (`submission_idempotency`), and atomicity guarantees are sufficient. The `BOOKKEEPING FREEZE` criteria are not triggered: no concrete failure in an actual work/review cycle was observed that would require review-coverage persistence, context snapshot IDs, epistemic-kind fields, or agent attestation.

---

## FINAL QUESTION

> Is ARGUS now boring enough to freeze so the project can return to substantive BM–IST–AS research?

**Almost.** The write path, read path, authority surface, and provenance are all correct in the production code. The only remaining item is a **test/verification gap**: no automated test proves the full MCP stdio transport for `argus.get_context`. The code is ready for the final manual smoke test (`./bin/argus serve` + fresh `argus mcp` process), and the manual smoke test should close the verification gap. Once the MCP transport test from the P1 finding above is added, the harness can be frozen.
