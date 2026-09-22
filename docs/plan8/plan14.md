# Plan 14 — P0 Persistence Repair + Real MCP Integration Test

## Source

Remediation of P0 from `docs/plan8/adv_review7.md` and `docs/plan8/review9.md`.
The validation boundary is correct; the persistence layer is broken because
`origin_packet_id` FK constraints reference `packet_submission` before it exists.

## FK Dependency Analysis

### Discovered Schema

**`packet_submission`** (migration 003):
```sql
CREATE TABLE packet_submission (
    packet_id      STRING PRIMARY KEY,
    scenario_id    UUID NOT NULL,
    task_id        UUID,          -- NULLABLE, NO FK CONSTRAINT
    agent_id       STRING NOT NULL,
    ...
);
```

**`belief`, `evidence`, `conductor_task`, `edge_provenance`** (migration 004):
All have `origin_packet_id` with `REFERENCES packet_submission(packet_id)`.

### Critical Finding: No Circular Dependency

- `packet_submission.task_id` has **NO FK** to `conductor_task`
- `packet_submission` references nothing via FK
- Dependency is strictly one-directional: `belief/evidence/task/edge_provenance → packet_submission`
- **Safe to insert `packet_submission` first**

### Current Broken Ordering

```
Step 1: belief INSERT         ← origin_packet_id references packet_submission (NOT YET INSERTED)
Step 2: evidence INSERT       ← same FK problem
Step 3: edge INSERT           ← edge_provenance.origin_packet_id same FK problem
Step 4: task INSERT           ← origin_packet_id same FK problem
Step 5: idempotency INSERT
Step 6: packet_submission INSERT  ← TOO LATE
```

### Required Ordering

```
BEGIN
  packet_submission    ← INSERT FIRST (no FK deps)
  beliefs
  evidence
  edges + edge_provenance
  tasks
  idempotency
COMMIT
```

---

## Implementation Phases

### Phase 1: Reorder Persist()

**File:** `internal/application/app.go` — `Persist()` method

Move the `packet_submission` INSERT (currently lines 585-603) to BEFORE the belief loop (before line 412).

The `taskRef` variable (currently computed at lines 570-578) must also move earlier — it's needed for the `packet_submission` INSERT.

New ordering:
1. Compute `taskRef` (moved from step 6)
2. Insert `packet_submission` (moved from step 6)
3. Insert beliefs (existing step 1)
4. Insert evidence (existing step 2)
5. Input-to-claim binding check (existing step 2.5)
6. Insert edges + edge_provenance (existing step 3)
7. Insert tasks (existing step 4)
8. Insert idempotency row (existing step 5)

All within the same transaction. No DEFERRABLE FKs needed.

### Phase 2: DB-backed MCP Happy-Path Test

**File:** `internal/mcp/adapter_test.go`

Add `TestMCPSubmitPacketHappyPath`:
- Uses `testDB(t)` helper for real CockroachDB
- Creates adapter with real `application.New(db)` + pack registry
- Submits valid packet through `adapter.HandleTool("argus.submit_packet", ...)`
- Asserts no error
- Queries DB to verify rows in:
  - `packet_submission` (1 row, correct agent fields)
  - `belief` (N rows, correct `origin_packet_id`)
  - `evidence` (if included, correct `origin_packet_id`)
  - `belief_edge` (if included)
  - `edge_provenance` (if included, correct `origin_packet_id`)
  - `conductor_task` (if included, correct `origin_packet_id`)

### Phase 3: DB-backed Malformed-Packet Test

**File:** `internal/mcp/adapter_test.go`

Add `TestMCPRejectsMalformedWithRealDB`:
- Uses `testDB(t)` for real CockroachDB
- Submits malformed packet (missing agent.id) through `adapter.HandleTool`
- Asserts error returned
- Queries ALL tables to verify zero rows created:
  - `packet_submission` = 0
  - `belief` = 0
  - `evidence` = 0
  - `belief_edge` = 0
  - `edge_provenance` = 0
  - `conductor_task` = 0

### Phase 4: Idempotency Test

**File:** `internal/mcp/adapter_test.go`

Add `TestMCPIdempotency`:
- Submits same valid packet twice through `adapter.HandleTool`
- Verifies no duplicates:
  - `packet_submission` = 1 row
  - `belief` = N rows (not 2N)
  - `evidence` = M rows (not 2M)
  - `origin_packet_id` values unchanged

### Phase 5: Atomic Rollback Test

**File:** `internal/mcp/adapter_test.go`

Add `TestMCPAtomicRollback`:
- Submits a packet that will fail during edge insertion (e.g., contradicts edge targeting a local reference)
- Verifies error returned
- Queries ALL tables to verify zero rows remain

### Phase 6: Existing Test Repair

**File:** `internal/application/app_test.go`

The existing DB-backed tests (`TestPartialFailureRetrySucceeds`, `TestIdempotencyDuplicatePacket`, `TestAgentIdentityPersistedInSubmission`, `TestFullIntegration`) should now pass with the FK ordering fix. Verify they pass.

### Phase 7: Verify

```bash
go vet ./...
go build ./...
go test ./...    # with CockroachDB running
```

---

## Files Changed (estimated)

| File | Change |
|------|--------|
| `internal/application/app.go` | Move `packet_submission` INSERT + `taskRef` computation before belief loop |
| `internal/mcp/adapter_test.go` | Add 4 DB-backed integration tests + import testdb helper |

## Invariants Preserved

1. Append-only `packet_submission` — `ON CONFLICT (packet_id) DO NOTHING`
2. `origin_packet_id` immutability — same as before
3. Full rollback on any failure — same `defer tx.Rollback()`
4. `ON CONFLICT DO NOTHING` semantics — unchanged
5. All writes in same transaction — unchanged
