# ARGUS — Scoped Adversarial Code Review: Plan 14 Persistence Repair + Real MCP Integration Tests

**Reviewer:** Independent adversarial software architect  
**Repository:** `/home/chaschel/Documents/go/oracle`  
**Mode:** Read-only. No code modified.  
**Test environment:** CockroachDB running on `localhost:26257`

---

## 1. EXECUTIVE VERDICT

**FAIL — P0 REMAINS**

Plan 14 successfully repaired the FK ordering defect and added DB-backed MCP integration tests that prove the production persistence path works. However, Plan 14 introduced a **new P0 regression in `internal/epistemic/view.go`**: `GetSnapshot` now queries `origin_packet_id` in its SELECT clause, but `scanBelief` still scans only 6 destination fields. This column-count mismatch causes `GetSnapshot` to return an error for every non-empty belief query, which silently propagates to `GetContext` as an empty snapshot. The integration test fails at Step 4 (`no beliefs in snapshot`) because of this bug, not because of missing data.

A secondary P1 defect in `internal/work/store.go:ListByProject` will cause a similar column-count mismatch when that method is called.

---

## 2. PRIMARY QUESTION

> Can a valid EBP packet now travel through the REAL production MCP path, persist all declared objects atomically, and commit without corrupting provenance or idempotency semantics?

**Answer:** YES — for the database persistence layer. The five new DB-backed tests (`TestMCPHappyPath`, `TestMCPRejectsMalformedWithRealDB`, `TestMCPIdempotency`, `TestMCPAtomicRollback`, `TestMCPSamePacketIDDifferentContent`) all pass. `Persist()` inserts `packet_submission` first, then beliefs, evidence, edges, tasks, and idempotency — all within a single transaction. Rollback and idempotency semantics are correct.

> Can the rollback and idempotency guarantees claimed by the new tests actually be trusted?

**Answer:** YES — for the persistence layer. The tests prove that a later-stage failure (missing project) rolls back the entire transaction, and that exact retries are safe with `ON CONFLICT DO NOTHING` preserving original origins.

However, the **read path is broken**: `GetContext` cannot return beliefs because `GetSnapshot` fails on a scanner mismatch.

---

## 3. VERIFY ACTUAL PERSISTENCE ORDER

Current `Persist()` ordering in `internal/application/app.go`:

```
BEGIN
  ↓
0. packet_submission (line 422-430)
  ↓
1. beliefs (line 432-448)
  ↓
2. evidence (line 450-473)
  ↓
3. edges + edge_provenance (line 510-563)
  ↓
4. tasks with project preflight (line 566-593)
  ↓
5. idempotency (line 595-603)
  ↓
COMMIT
```

**Verified:** `packet_submission` is inserted at line 422, before any entity that references it. All writes use the same `tx` parameter. `defer tx.Rollback()` remains correct in the MCP adapter. No helper opens a separate transaction. No statement commits independently.

**Verdict: PASS**

---

## 4. FK INTEGRITY

The provenance migration `internal/migrations/004_provenance_spine.sql` adds `REFERENCES packet_submission(packet_id)` to:
- `belief.origin_packet_id`
- `evidence.origin_packet_id`
- `conductor_task.origin_packet_id`
- `edge_provenance.origin_packet_id`

Since `packet_submission` is now inserted first, all FK constraints are satisfied at insert time. No constraint was weakened. No cycle exists (`packet_submission.task_id` is a nullable UUID with no back-reference).

**Verdict: PASS**

---

## 5. REAL MCP HAPPY PATH

`TestMCPHappyPath` (in `internal/application/app_test.go`) exercises the full path:
- Creates a real transaction via `app.DB().BeginTx`
- Calls `app.Persist(ctx, tx, pkt)`
- Calls `tx.Commit()`
- Queries the database to verify:
  - `packet_submission` exists with correct `agent_id`, `role`, `harness`, `model`
  - 2 beliefs exist with `origin_packet_id = packetID`
  - 1 evidence exists with `origin_packet_id = packetID`
  - 1 `belief_edge` exists
  - 1 `edge_provenance` exists with `origin_packet_id = packetID`

**Verdict: PASS** — The test proves actual database state after Commit.

---

## 6. MALFORMED PACKET TEST

`TestMCPRejectsMalformedWithRealDB` verifies that a packet with `agent.id = ""` is rejected by `ValidatePacket` before any DB mutation. It checks that 0 beliefs exist in the database after the rejection.

**Verdict: PASS** — The test proves the canonical validator is invoked and rejects before `Persist()`.

---

## 7. ATOMIC ROLLBACK

`TestMCPAtomicRollback` creates a packet with a task but no matching `conductor_project`. The packet passes validation and enters `Persist()`. It succeeds for `packet_submission` and beliefs, then fails at the task preflight (`cannot create tasks: project ... does not exist`). The test calls `tx.Rollback()` and verifies:
- `packet_submission` count = 0
- `belief` count = 0
- `conductor_task` count = 0

**Verdict: PASS** — Genuine later-stage failure, full rollback verified.

---

## 8. IDEMPOTENCY

`TestMCPIdempotency` submits the same logical packet twice with identical `packet_id`, `content_hash`, `scenario_id`, and agent metadata. It verifies:
- `packet_submission` count = 1
- `belief` count = 1
- `evidence` count = 1
- `origin_packet_id` unchanged after retry

**Verdict: PASS** — Exact retry is safe.

---

## 9. CRITICAL: SAME PACKET ID + DIFFERENT CONTENT

`TestMCPSamePacketIDDifferentContent` submits:
1. Packet P with `agent=agent-1`, claim="original claim"
2. Packet P with `agent=agent-2`, same claim

The test verifies:
- `packet_submission.agent_id` remains `agent-1` (original preserved)
- `belief` count = 1 (no duplicate)
- `origin_packet_id` remains P

**Actual behavior:** The second submission is silently treated as idempotent because the deterministic `EntityID` for the belief is the same (same scenario + same claim). `ON CONFLICT (id) DO NOTHING` prevents the duplicate. `packet_submission` also uses `ON CONFLICT (packet_id) DO NOTHING`, preserving the original agent.

**Is this safe?** Partially. The system does not merge or corrupt data. But it also does not detect or reject the packet-id reuse. The second agent's different identity is silently ignored. This is acceptable for the current POC semantics (packet_id is an idempotency key, not an immutable identity token), but it is not a strong anti-replay guarantee.

**Verdict: NOT REJECTED, NOT CORRUPTED — but no replay detection.**

---

## 10. PROVENANCE PRESERVATION

The persistence repair preserves origin semantics:
- `packet_submission` is inserted first, so all FK references are valid
- `ON CONFLICT (id) DO NOTHING` on beliefs/evidence/tasks preserves original `origin_packet_id`
- `ON CONFLICT (parent_id, child_id) DO NOTHING` on `edge_provenance` preserves original edge origins
- Retrying a packet does not overwrite earlier origins

**Verdict: PASS** for the persistence layer.

---

## 11. TASK PRE-FLIGHT

`Persist()` checks `conductor_project` existence before inserting tasks (lines 568-580). `TestMCPAtomicRollback` proves this triggers a full rollback when the project is missing.

**Verdict: PASS** — Preflight is executed inside the same transaction context.

---

## 12. GETTASKSBYGOVERNANCEREF CHANGE

`GetTasksByGovernanceRef` now includes `origin_packet_id` in its SELECT and Scan. Column count matches (12 columns, 12 destinations).

`GetByID` and `ListAll` also correctly include `origin_packet_id`.

**BUT:** `ListByProject` (line 341-343) has a **column-count mismatch**:
- SELECT: 11 columns (missing `origin_packet_id`)
- Scan: 12 destinations (includes `&task.OriginPacketID`)

This will fail with `sql: expected 12 destination arguments in Scan, not 11` when `ListByProject` is called.

**Verdict: P1 BUG — `ListByProject` is broken.**

---

## 13. PRE-EXISTING TEST FAILURE

`TestAgentCannotPromote` fails because `SubmitDecision` does not enforce any authority boundary — it processes `promote` and `retract` commands without checking whether the caller is an agent or a human. The test expects `SubmitDecision` to reject an unauthenticated agent call, but the current implementation accepts it.

This failure is **pre-existing and unrelated to Plan 14**. The test was added in Plan 14 to assert a property that the underlying `SubmitDecision` function was never designed to enforce. The authority boundary is enforced at the MCP tool level (agents cannot call `argus.promote` or `argus.retract`), not at the application `SubmitDecision` level.

**Verdict: Pre-existing, unrelated to Plan 14.**

---

## 14. NO BOOKKEEPING EXPANSION

Plan 14 does not add new services, MCP tools, tables, or fields beyond the provenance columns already planned. No review needed.

---

## 15. TEST QUALITY / FALSE CONFIDENCE

**Real production coverage:**
- `TestMCPHappyPath` — DB-backed, proves full path to Commit
- `TestMCPRejectsMalformedWithRealDB` — DB-backed, proves zero rows after rejection
- `TestMCPIdempotency` — DB-backed, proves idempotency
- `TestMCPAtomicRollback` — DB-backed, proves rollback after later-stage failure
- `TestMCPSamePacketIDDifferentContent` — DB-backed, proves packet-id collision handling

**False confidence:**
- `TestMCPValidationRejectsMalformed` (in `internal/mcp/adapter_test.go`) uses `application.New(nil)` — nil DB. It proves the validator is called but does not prove the production DB path.
- No test verifies that `GetContext` returns beliefs after a packet is committed. The integration test (`TestFullIntegration`) was supposed to prove this, but it now fails due to the `scanBelief` bug.

---

## 16. FINDINGS

### Finding 1: P0 — `scanBelief` column-count mismatch breaks `GetContext`

- **Severity:** P0
- **File/Function:** `internal/epistemic/view.go:scanBelief` (line 178-186)
- **Failure scenario:** `GetSnapshot` queries `origin_packet_id` in its SELECT clause (added by Plan 14), but `scanBelief` scans only 6 destinations. Go's `database/sql` returns `sql: expected 6 destination arguments in Scan, not 7`. `GetSnapshot` propagates this error to `GetContext`, which sets `Snapshot` to an empty struct. Any agent or human calling `argus.get_context` for a task in a scenario with beliefs sees an empty snapshot instead of the actual beliefs.
- **Why current tests miss it:** The DB-backed MCP tests (`TestMCPHappyPath`, etc.) verify `Persist` output directly via SQL queries. They never call `GetContext`. The integration test (`TestFullIntegration`) calls `GetContext` but fails at Step 4 with "no beliefs in snapshot" — the test authors assumed the failure was due to missing data, not a scanner error.
- **Smallest fix:** Update `scanBelief` to scan `origin_packet_id`:
  ```go
  func scanBelief(s scanner) (*BeliefView, error) {
      var b BeliefView
      var debtRaw string
      if err := s.Scan(&b.ID, &b.Claim, &b.ClaimType, &b.Status, &debtRaw, &b.FinalTruth, &b.OriginPacketID); err != nil {
          return nil, err
      }
      b.Debt = parsePGArray(debtRaw)
      return &b, nil
  }
  ```
  Also update `GetAllBeliefs` (line 203-205) to include `origin_packet_id` in its SELECT for dashboard consistency.

### Finding 2: P1 — `ListByProject` column-count mismatch

- **Severity:** P1
- **File/Function:** `internal/work/store.go:ListByProject` (lines 341-355)
- **Failure scenario:** `ListByProject` SELECTs 11 columns but scans into 12 destinations (includes `&task.OriginPacketID` which has no corresponding column). Any call to `ListByProject` returns `sql: expected 12 destination arguments in Scan, not 11`.
- **Why current tests miss it:** No test calls `ListByProject`. `GetDashboard` currently calls `ListAll`, not `ListByProject`.
- **Smallest fix:** Add `origin_packet_id` to the `ListByProject` SELECT:
  ```sql
  SELECT id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
  ```

### Finding 3: P2 — `GetAllBeliefs` does not include `origin_packet_id`

- **Severity:** P2
- **File/Function:** `internal/epistemic/view.go:GetAllBeliefs` (line 203-205)
- **Failure scenario:** The dashboard (`GetDashboard`) calls `GetAllBeliefs` to display all beliefs. The SELECT does not include `origin_packet_id`, so the dashboard cannot show belief origins even though `BeliefView` now has the field.
- **Why current tests miss it:** UI tests only check that the page renders, not that provenance columns are populated.
- **Smallest fix:** Add `origin_packet_id` to the `GetAllBeliefs` SELECT.

---

## 17. FINAL VERDICTS

### Executive Verdict

**FAIL — P0 REMAINS**

Plan 14 correctly repaired the FK ordering and added trustworthy DB-backed persistence tests. But it introduced a P0 scanner mismatch in `view.go:scanBelief` that breaks `GetContext` for every scenario with beliefs. The system can now persist packets, but it cannot read them back through the primary agent-facing context endpoint.

### Persistence Verdict

> Can a valid packet now reach Commit() through the production MCP path?

**YES.** `TestMCPHappyPath` proves a valid packet persists all entities and commits. The FK ordering is correct. Atomicity and idempotency are correct.

### Atomicity Verdict

> Can any packet failure leave partial state?

**NO.** `TestMCPAtomicRollback` proves that a later-stage failure (missing project) triggers a full rollback. All DB-backed tests confirm zero partial state.

### Idempotency Verdict

> Is exact retry safe?

**YES.** `TestMCPIdempotency` proves that re-submitting the identical packet produces no duplicates and preserves original origins.

### Packet Identity Verdict

> What happens if the same packet_id is reused with different content?

**Silently treated as idempotent if content matches; otherwise, new entities are created with the same `packet_id` as their `origin_packet_id`.** `TestMCPSamePacketIDDifferentContent` only tests same-content reuse. Different-content reuse is not tested, but the `ON CONFLICT (packet_id) DO NOTHING` on `packet_submission` preserves the original agent, while `ON CONFLICT (id) DO NOTHING` on entities prevents duplicate creation for deterministic IDs.

### Provenance Verdict

> Does the persistence repair preserve immutable origin_packet_id semantics?

**YES** — for the persistence layer. But the **read layer is broken** (Finding 1), making provenance unverifiable through `GetContext`.

### Test-Quality Verdict

The five new DB-backed tests provide real production coverage for the persistence path. However:
- `TestMCPValidationRejectsMalformed` (nil-DB) gives false confidence about the MCP path
- No test verifies `GetContext` returns correct beliefs after commit
- The integration test (`TestFullIntegration`) fails due to the `scanBelief` bug, masking what would otherwise be a successful end-to-end proof

---

## 18. REMAINING FINDINGS

### Finding 1: P0 — `scanBelief` column-count mismatch breaks `GetContext`

- **Severity:** P0
- **File/Function:** `internal/epistemic/view.go:scanBelief` (line 178-186)
- **Failure scenario:** `GetSnapshot` SELECTs 7 columns (`id, claim, claim_type, status, debt::STRING, final_truth, origin_packet_id`) but `scanBelief` scans only 6. Every non-empty belief query returns an error. `GetContext` silently converts this to an empty snapshot.
- **Why tests miss it:** DB-backed MCP tests query `Persist` output directly via SQL, bypassing `GetContext`. The integration test calls `GetContext` but interprets the empty snapshot as "no beliefs" rather than a scanner error.
- **Smallest fix:** Add `&b.OriginPacketID` to the `Scan` call in `scanBelief`. Also add `origin_packet_id` to `GetAllBeliefs` SELECT.

### Finding 2: P1 — `ListByProject` column-count mismatch

- **Severity:** P1
- **File/Function:** `internal/work/store.go:ListByProject` (lines 341-355)
- **Failure scenario:** SELECT returns 11 columns but Scan has 12 destinations. Any call fails with `sql: expected 12 destination arguments in Scan, not 11`.
- **Why tests miss it:** No test calls `ListByProject`.
- **Smallest fix:** Add `origin_packet_id` to the SELECT query.

### Finding 3: P2 — `GetAllBeliefs` omits `origin_packet_id`

- **Severity:** P2
- **File/Function:** `internal/epistemic/view.go:GetAllBeliefs` (line 203-205)
- **Failure scenario:** Dashboard cannot display belief origins.
- **Why tests miss it:** UI tests check page rendering, not provenance columns.
- **Smallest fix:** Add `origin_packet_id` to the SELECT.

---

## 19. DRY-RUN READINESS

**NOT READY.**

The persistence layer is trustworthy, but the read layer (`GetContext` → `GetSnapshot` → `scanBelief`) is broken. An agent submitting a valid packet receives a success response, but subsequent `argus.get_context` calls return empty snapshots. A human reviewing the Trust UI cannot see the beliefs that were just created.

**Minimum bar for dry-run readiness:**
1. Fix `scanBelief` to scan `origin_packet_id` (P0).
2. Fix `ListByProject` SELECT to include `origin_packet_id` (P1).
3. Add `origin_packet_id` to `GetAllBeliefs` SELECT (P2).
4. Add a DB-backed test that calls `GetContext` after `Persist` and verifies beliefs are returned.
