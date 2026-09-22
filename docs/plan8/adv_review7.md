# ARGUS — Scoped Adversarial Code Review: Production MCP Validation Remediation

**Reviewer:** Independent adversarial software architect  
**Repository:** `/home/chaschel/Documents/go/oracle`  
**Mode:** Read-only. No code modified.  
**Test environment:** CockroachDB running on `localhost:26257`

---

## 1. EXECUTIVE VERDICT

**FAIL — P0 REMAINS**

The remediation established correct validation plumbing on the production MCP path, but introduced a **fatal foreign-key ordering defect** in `Persist()` that prevents ANY packet from being persisted. Every DB-backed test that exercises the full MCP → Compile → Validate → ValidatePacket → Persist → Commit path fails with:

```
ERROR: insert on table "belief" violates foreign key constraint 
"belief_origin_packet_id_fkey" (SQLSTATE 23503)
```

The system cannot reach `Commit()` for any well-formed packet. The validation layer is trustworthy; the persistence layer is not.

---

## 2. PRIMARY QUESTION

> Can a buggy or malicious AI agent submit a structurally malformed EBP packet through `argus.submit_packet` and cause invalid research state to reach `Persist()`?

**Answer for malformed packets:** NO. The production MCP path now invokes `packet/v1.Validate` via `app.ValidatePacket`, which rejects structurally malformed packets before any database write.

**Answer for well-formed packets:** The system cannot persist ANY packet, well-formed or malformed, because `Persist()` violates a foreign-key constraint. This is a P0 correctness failure.

---

## 3. PRODUCTION MCP PATH TRACE

Current ordering in `internal/mcp/adapter.go:handleSubmitPacket`:

```
JSON unmarshal
    ↓
app.Compile(ctx, &pkt)          // schema_version, packet_id, scenario_id, pack_ref, non-empty
    ↓
app.Validate(ctx, &pkt)         // debt vocab, evidence classes, verifier specs
    ↓
app.ValidatePacket(ctx, &pkt)   // canonical packet/v1.Validate
    ↓
tx.BeginTx(ctx, nil)
    ↓
app.Persist(ctx, tx, &pkt)
    ↓
tx.Commit()
```

**Validation is now correctly placed before `BeginTx`.** Malformed packets are rejected before any transaction is opened. This is the correct design.

---

## 4. VALIDATION COVERAGE

### 4.1 Production MCP validation (via `packet/v1.Validate`)

| Check | Result | Function |
|-------|--------|----------|
| A. Missing `agent.id` | PASS | `validateAgent` |
| B. Missing `agent.role` | PASS | `validateAgent` |
| C. Missing `agent.harness` | PASS | `validateAgent` |
| D. Missing `agent.model` | PASS | `validateAgent` |
| E. `agent.role != packet.role` | PASS | `validateAgent` |
| F. Invalid packet role | PASS | `Validate` (line 29) |
| G. Duplicate belief `local_id` | PASS | `validateBeliefs` |
| H. Duplicate evidence `local_id` | PASS | `validateEvidence` |
| I. Duplicate edge `local_id` | PASS | `validateEdges` |
| J. Duplicate task `local_id` | PASS | `validateTasks` |
| K. Cross-type `local_id` collision | PASS | `validateReferences` (global namespace) |
| L. Empty belief claim | PASS | `validateBeliefs` |
| M. Invalid `claim_type` | PASS | `validateBeliefs` |
| N. Invalid edge kind | PASS | `validateEdges` |
| O. Malformed `local:` reference | PASS | `ParseReference` in `validateReferences` |
| P. Malformed `canonical:belief:` reference | PASS | `ParseReference` in `validateReferences` |
| Q. Invalid `pack_ref` | PASS | `parsePackRef` + `registry.Get` |
| R. Nonexistent pack | PASS | `registry.Get` returns error |
| S. Malformed task structures | PASS | `validateTasks` |

### 4.2 Defense-in-depth in `Persist()`

| Check | Result | Location |
|-------|--------|----------|
| Global `local_id` uniqueness | PASS | `Persist` lines 362-393 |
| `claim_type` enum guard | PASS | `Persist` lines 396-403 |
| Edge kind guard | PASS | `Persist` lines 406-410 |
| Unresolved references | PASS | `resolveRef` returns `""` |
| Cross-scenario edges | PASS | `validateEdgeScenario` |
| Self-edges | PASS | `Persist` lines 498-500 |
| Local `contradicts` target | PASS | `Persist` lines 502-504 |
| Task project preflight | PASS | `Persist` lines 548-560 |
| `operator_asserted` rejection | PASS | `Persist` lines 433-435 |

### 4.3 `pack_ref` grammar consistency

- `packet/v1/validate.go:parsePackRef` — NOW uses `@` (fixed)
- `internal/application/app.go:parsePackRef` — uses `@`
- `domain-pack/registry.go:registryKey` — uses `@`
- All production code and tests use `bmist@1.1.0` or `bmist@1.0.0`

**PASS** — canonical grammar is consistent.

---

## 5. CRITICAL P0 FINDING: FOREIGN-KEY ORDERING DEFECT

### Finding 1: `Persist()` inserts entities before `packet_submission`, violating FK constraints

- **File/Function:** `internal/application/app.go:Persist` (lines 353-606)
- **Exact failure:** `belief.origin_packet_id`, `evidence.origin_packet_id`, `conductor_task.origin_packet_id`, and `edge_provenance.origin_packet_id` all have `REFERENCES packet_submission(packet_id)`. But `Persist()` inserts beliefs (line 419), evidence (line 445), edges (line 509), tasks (line 564), and edge provenance (line 535) BEFORE inserting `packet_submission` (line 594). CockroachDB enforces FK constraints immediately at statement execution time. The first belief insert fails with SQLSTATE 23503.
- **Why current tests do not catch it (prior to this review):** The MCP adapter tests (`internal/mcp/adapter_test.go`) use `application.New(nil)` — a nil DB — so `Persist()` is never reached. The fault-injection test (`TestPartialFailureRetrySucceeds`) fails on the first `ExecContext`, before any belief insert. No DB-backed test exercising the full MCP → Persist → Commit path with the new provenance schema existed.
- **Test evidence of failure:**
  - `TestPartialFailureRetrySucceeds`: `insert belief: ERROR: insert on table "belief" violates foreign key constraint "belief_origin_packet_id_fkey"`
  - `TestIdempotencyDuplicatePacket`: same FK violation
  - `TestAgentIdentityPersistedInSubmission`: same FK violation
  - `TestFullIntegration` (Step 3): same FK violation
- **Why this is P0:** The entire packet persistence path is non-functional. No well-formed packet can be committed. The validation remediation is irrelevant because the database layer rejects everything.

### Recommended minimal fix

Insert `packet_submission` BEFORE any entity that references it:

```go
// In Persist(), move packet_submission insert to BEFORE the belief loop:
// 0. Persist packet submission provenance FIRST.
_, err = tx.ExecContext(ctx,
    `INSERT INTO packet_submission
         (packet_id, scenario_id, task_id, agent_id, role, harness, model, content_sha256)
     VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
     ON CONFLICT (packet_id) DO NOTHING`,
    pkt.PacketID, pkt.ScenarioID, taskRef,
    pkt.Agent.ID, pkt.Agent.Role, pkt.Agent.Harness, pkt.Agent.Model, hash)
if err != nil {
    return nil, fmt.Errorf("insert packet submission: %w", err)
}

// 1. Persist beliefs...
// 2. Persist evidence...
// 3. Persist edges...
// 4. Persist tasks...
// 5. Write idempotency row LAST.
```

Alternatively, make the FK constraints `DEFERRABLE INITIALLY DEFERRED` in `004_provenance_spine.sql`.

---

## 6. SECONDARY FINDINGS

### Finding 2: MCP regression tests give false confidence about production path

- **File/Function:** `internal/mcp/adapter_test.go:TestMCPValidationRejectsMalformed`
- **Failure scenario:** The test creates the adapter with `application.New(nil)` — a nil database handle. When `HandleTool` calls `app.ValidatePacket`, it passes. But if the test ever reached `Persist()`, it would panic on nil DB. More importantly, the test never proves that a valid packet can survive the full path to `Commit()`. The test suite has zero DB-backed MCP integration tests.
- **Why tests miss it:** The test is entirely in-memory. It proves the validator is called, but not that the validator's output leads to successful persistence.
- **Recommended minimal fix:** Add a DB-backed test that submits a valid packet through `adapter.HandleTool` and asserts that all expected rows exist after commit. Add a DB-backed test that submits an invalid packet and asserts zero rows were inserted.

### Finding 3: `GetDashboard` still leaks all tasks across all scenarios

- **File/Function:** `internal/application/app.go:GetDashboard` (line 619)
- **Failure scenario:** `a.workStore.ListAll(ctx)` returns every task in every project. A human viewing `/ui/insights?scenario_id=A` sees tasks from scenarios B, C, and D.
- **Why tests miss it:** `integration_test.go:TestFullIntegration` only checks that the page renders, not that it is scenario-scoped.
- **Recommended minimal fix:** When `scenarioID != ""`, call `a.workStore.ListByProject(ctx, scenarioID)` instead of `ListAll`.

### Finding 4: Evidence and edge provenance invisible in UI

- **File/Function:** `internal/epistemic/view.go:EvidenceView`, `EdgeView`, `GetSnapshot`
- **Failure scenario:** The Trust UI Insights page shows belief origins but not evidence origins or edge origins. `edge_provenance` table is never queried. `evidence.origin_packet_id` is not selected.
- **Why tests miss it:** UI tests only check for page-rendering strings, not provenance columns.
- **Recommended minimal fix:** Add `origin_packet_id` to `EvidenceView` and `EdgeView`, join `edge_provenance` in `GetSnapshot`, and update `insights.html`.

---

## 7. ATOMICITY VERDICT

**NO — the system cannot persist any packet.**

The transaction boundary in `internal/mcp/adapter.go` is correct: `BeginTx` → `Persist` → `Commit` with `defer tx.Rollback()`. Every write in `Persist()` uses the same `tx`. If any write fails, the transaction rolls back.

However, because the FK ordering defect causes EVERY packet insert to fail, the system is currently in a state where:

```
any well-formed packet
    ->
    Persist()
    ->
    belief insert fails (FK violation)
    ->
    tx.Rollback()
    ->
    ZERO objects persisted
```

This is technically "all-or-nothing" (no partial persistence), but it is a **complete functional failure**, not a correct atomicity guarantee.

---

## 8. TASK INTEGRITY VERDICT

**N/A — no packet can reach task insertion.**

Because the first belief insert fails, no packet ever reaches the task insertion stage. The task preflight (`conductor_project` existence check) is correctly implemented but is never exercised.

---

## 9. PROVENANCE VERDICT

**Partially preserved, but unverifiable.**

The provenance data model is correct:
- `packet_submission` records agent identity
- `origin_packet_id` on beliefs, evidence, tasks, and `edge_provenance` is populated on creation
- `ON CONFLICT (id) DO NOTHING` preserves original origins on retry

However, because NO packet can be persisted, these invariants cannot be observed in production. The previous review's provenance concerns remain valid, but they are now secondary to the fundamental persistence failure.

---

## 10. PACK REF GRAMMAR

**PASS.** The entire repository now uses `@` as the canonical separator:
- `packet/v1/validate.go:parsePackRef` — `@`
- `internal/application/app.go:parsePackRef` — `@`
- `domain-pack/registry.go:registryKey` — `@`
- All test packets use `bmist@1.1.0` or `bmist@1.0.0`

No stale `-` separator usage found in production paths.

---

## 11. AGENT IDENTITY

**PASS.** All four fields are required and `agent.role == packet.role` is enforced in `packet/v1/validate.go:validateAgent`. The production MCP path invokes this via `app.ValidatePacket`. Agent identity is provenance metadata only — it does not grant authorization, substitute for a human principal, or affect promotion/discharge authority.

---

## 12. LOCAL-ID INTEGRITY

**PASS.** `packet/v1/validate.go:validateReferences` enforces global `local_id` uniqueness across all entity types within a single packet. `Persist()` repeats this check as defense-in-depth. Cross-type collisions (e.g., belief `local_id="x"` and evidence `local_id="x"`) are rejected before any DB write.

---

## 13. CLAIM TYPE

**PASS.** Allowed values (`derived`, `accommodated`, `postulated`) are enforced in `packet/v1/validate.go:validateBeliefs` and `Persist()`. No divergent enum definitions found.

---

## 14. EDGE KIND

**PASS.** Allowed kinds (`derives`, `contradicts`) are enforced in `packet/v1/validate.go:validateEdges` and `Persist()`. DB-level CHECK constraints provide final protection.

---

## 15. EDGE TARGET RULES

**PASS.** `Persist()` enforces:
- Endpoints exist (via `validateEdgeScenario` which queries `belief.scenario_id`)
- Both endpoints belong to the same scenario (`validateEdgeScenario`)
- `parent != child` (self-edge check)
- `contradicts` targets canonical existing belief only (local target rejected)

---

## 16. TASK PRE-FLIGHT

**Code present but unreachable.** `Persist()` correctly checks `conductor_project` existence before inserting tasks (lines 548-560). However, because the FK ordering defect causes the first belief insert to fail, this preflight is never exercised in practice.

---

## 17. MCP REGRESSION TEST QUALITY

**PARTIAL.** `TestMCPValidationRejectsMalformed` proves the production MCP adapter invokes validators for:
- missing `agent.id`
- role mismatch
- duplicate belief `local_id`
- invalid edge kind
- invalid `claim_type`
- cross-type `local_id` collision

These are real production-path tests (they call `adapter.HandleTool`). However:

1. **No DB-backed MCP integration test exists.** No test proves a valid packet survives the full path to `Commit()`.
2. **The tests pass only because they use `application.New(nil)`** (nil DB). If the adapter stopped calling `ValidatePacket` but kept calling `Compile` and `Validate`, the tests would still pass for the cases they cover, because `Compile` + `Validate` do not check agent identity, edge kinds, or `local_id` uniqueness.
3. **No test verifies the fix for the original bypass.** There is no test that asserts `packet/v1.Validate` is called for every malformed-packet case (e.g., invalid `pack_ref`, malformed `canonical:belief:` reference, empty `claim`).

---

## 18. DOCUMENTATION CONSISTENCY

The recent changes do not introduce incorrect documentation claims. `prompts/argus-agent-protocol.md` correctly states that agent identity is required and that validation failures are informative.

---

## 19. FINAL VERDICTS

### Executive Verdict

**FAIL — P0 REMAINS**

The validation remediation is architecturally correct but functionally broken. `Persist()` cannot insert any entity because `origin_packet_id` foreign keys reference `packet_submission` before it exists in the database.

### Production MCP Boundary

> Can a malformed AI packet reach `Persist()` through `argus.submit_packet`?

**NO** — malformed packets are rejected by `packet/v1.Validate` before any transaction begins.

**BUT** — well-formed packets also cannot reach `Commit()` because `Persist()` violates an FK constraint on the first entity insert.

### Validation Coverage

| Check | Production MCP | Persist defense | DB constraint | Test |
|-------|----------------|-----------------|---------------|------|
| Agent identity | PASS | PASS (packet_submission) | NOT NULL | PASS (nil-DB) |
| Role consistency | PASS | PASS (packet_submission) | NOT NULL | PASS (nil-DB) |
| Global local_id uniqueness | PASS | PASS | UNIQUE (implicit via PK) | PASS (nil-DB) |
| Claim type | PASS | PASS | — | PASS (nil-DB) |
| Edge kind | PASS | PASS | — | PASS (nil-DB) |
| Reference format | PASS | PASS | — | PASS (nil-DB) |
| Pack ref grammar | PASS | — | — | PASS (nil-DB) |
| Task project preflight | — | PASS | FK | NOT REACHED |
| Origin packet FK | — | BROKEN | ENFORCED | FAILS |

### Atomicity

> Can a packet partially persist?

**NO — but not by design.** Any failure in `Persist()` rolls back the transaction. However, the first insert always fails due to FK violation, so the system currently achieves "zero persistence" rather than "all-or-nothing."

### Task Integrity

> Can an accepted packet lose a task while preserving other packet objects?

**N/A — no packet reaches task insertion.** The current code fails before any task is processed.

### Provenance

> Does the remediation preserve the existing provenance invariants?

**Unverifiable.** The data model is correct, but no data can be written. The previous review's provenance concerns (evidence/edge origins invisible in UI, dashboard scenario leakage) remain valid and are now secondary to the persistence failure.

---

## 20. REMAINING FINDINGS

### Finding 1: P0 — FK ordering makes `Persist()` non-functional

- **Severity:** P0
- **File/Function:** `internal/application/app.go:Persist` (lines 412-600)
- **Failure scenario:** Every packet insert fails at the first belief with `SQLSTATE 23503` because `origin_packet_id` references `packet_submission` before it exists.
- **Why tests miss it:** Prior tests either used nil DB or failed before reaching belief inserts. No DB-backed test exercised the full production MCP path with the new provenance schema.
- **Smallest fix:** Move the `packet_submission` insert to BEFORE the belief/evidence/edge/task loops. Alternatively, add `DEFERRABLE INITIALLY DEFERRED` to the FK constraints in `004_provenance_spine.sql`.

### Finding 2: P1 — No DB-backed MCP integration test

- **Severity:** P1
- **File/Function:** `internal/mcp/adapter_test.go`
- **Failure scenario:** The regression tests pass with a nil DB, giving false confidence that the production path works. There is no test that submits a valid packet through `HandleTool` with a real database and asserts successful commit.
- **Why tests miss it:** Tests were written with mock/nil DB to avoid database dependency, but this skips the exact layer that is now broken.
- **Smallest fix:** Add one DB-backed test: submit a minimal valid packet through `adapter.HandleTool`, then query the database to confirm belief, evidence, edge, task, and `packet_submission` rows exist.

### Finding 3: P2 — `GetDashboard` returns all tasks across all scenarios

- **Severity:** P2
- **File/Function:** `internal/application/app.go:GetDashboard` (line 619)
- **Failure scenario:** Human reviewing scenario A sees tasks from all scenarios.
- **Why tests miss it:** Integration test only checks page renders.
- **Smallest fix:** Use `ListByProject(ctx, scenarioID)` when `scenarioID != ""`.

### Finding 4: P2 — Evidence and edge provenance not visible in UI

- **Severity:** P2
- **File/Function:** `internal/epistemic/view.go:EvidenceView`, `EdgeView`, `GetSnapshot`
- **Failure scenario:** Humans cannot trace which packet created a specific evidence or edge.
- **Why tests miss it:** UI tests check for page rendering, not provenance columns.
- **Smallest fix:** Add `origin_packet_id` to view structs, query `edge_provenance` in `GetSnapshot`, update `insights.html`.

---

## 21. DRY-RUN READINESS

**NOT READY.**

The system cannot persist any packet through the production MCP path. A Work Agent submitting a valid packet will receive an error. An Adversarial Agent cannot challenge anything because no work exists to challenge.

**Minimum bar for dry-run readiness:**
1. Fix FK ordering in `Persist()` (move `packet_submission` insert before entity inserts, or make FKs deferrable).
2. Add a DB-backed MCP integration test that proves a valid packet can survive the full path to `Commit()`.
3. Verify all existing DB-backed tests pass after the FK fix.
