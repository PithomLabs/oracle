This review found a **real P0**, and unlike the earlier speculative FK-ordering concern, this one is now demonstrated by actual execution: every well-formed packet reaches `Persist()` and then fails on `belief.origin_packet_id → packet_submission(packet_id)` because `packet_submission` has not yet been inserted. The report records the failure in multiple DB-backed tests. 

So the current state is:

```text
Malformed packet
    ↓
correctly rejected
```

but:

```text
Valid packet
    ↓
correct validation
    ↓
Persist
    ↓
FK violation
    ↓
rollback
```

That means the **validation boundary is now good, but the actual research system is not operational**. The review correctly calls it not ready.  

## Fix I would choose

I would choose **move `packet_submission` insertion to the beginning of `Persist()`**, not introduce deferred foreign keys.

The desired transaction becomes:

```text
BEGIN
  ↓
packet_submission
  ↓
beliefs
  ↓
evidence
  ↓
edges
  ↓
edge_provenance
  ↓
tasks
  ↓
idempotency
  ↓
COMMIT
```

with:

```text
ANY FAILURE
    ↓
ROLLBACK EVERYTHING
```

This preserves your existing append-only provenance model and avoids introducing more complicated FK semantics.

There is one thing the coding agent must verify before moving it: **whether `packet_submission.task_id` itself has any FK dependency on `conductor_task`**. If it does, inserting `packet_submission` first creates a cycle. The agent should inspect the actual migration rather than assume. If `task_id` is only informational/nullable, the simple reorder is correct.

## The second finding is equally important

The code-reviewer correctly caught that the new MCP tests are not sufficient. They prove:

```text
HandleTool → validation → rejection
```

but not:

```text
HandleTool → validation → Persist → Commit
```

because they use `application.New(nil)`. 

That is exactly the kind of false confidence your process is supposed to catch.

The minimum real test should be:

```text
valid packet
  ↓
real HandleTool("argus.submit_packet")
  ↓
real DB
  ↓
Commit
  ↓
query:
  packet_submission
  belief
  evidence
  edge
  edge_provenance
  task
```

Then one malformed packet test should verify:

```text
malformed packet
  ↓
HandleTool
  ↓
reject
  ↓
zero DB mutations
```

The reviewer explicitly identifies this as the minimum readiness bar. 

## The P2 findings should wait

I agree with leaving these alone until persistence is working:

```text
GetDashboard scenario leakage
Evidence/edge provenance missing from UI
```

They are legitimate, but they are not blockers to restoring the core packet lifecycle. 

And this is exactly where your concern about bookkeeping becomes relevant: **don't chase those while the ledger cannot accept a valid packet. Fix the P0, prove the real path works, then return to domain research.**

## Give the coding agent this prompt

```text
# ARGUS — P0 Persistence Repair + Real MCP Integration Test

Fix the confirmed P0 from the latest adversarial code review.

The current production MCP validation boundary is correct, but valid packets
cannot persist because origin_packet_id foreign keys reference
packet_submission before packet_submission is inserted.

Observed failure:

    insert on table "belief" violates foreign key constraint
    "belief_origin_packet_id_fkey"
    SQLSTATE 23503

This occurs in the real production path:

    argus.submit_packet
        -> Compile
        -> Validate
        -> ValidatePacket
        -> Persist
        -> belief insert
        -> FK failure
        -> rollback

Do not broaden the architecture.

## 1. Fix Persist() ordering

Inspect the actual migration/schema first.

Prefer this ordering:

    BEGIN
      packet_submission
      beliefs
      evidence
      edges
      edge_provenance
      tasks
      idempotency
    COMMIT

All writes remain in the SAME transaction.

Before moving packet_submission earlier, verify whether
packet_submission.task_id itself has any FK dependency on conductor_task.

- If task_id has no FK dependency, insert packet_submission first.
- If it does, do NOT create a circular dependency. Report the exact schema
  and use the smallest safe alternative.

Do NOT introduce DEFERRABLE foreign keys unless the schema requires it and
the concrete transaction dependency proves that reordering is impossible.

Preserve:
- append-only packet provenance
- ON CONFLICT DO NOTHING semantics
- origin_packet_id immutability
- full rollback on any failure

## 2. Add a REAL DB-backed MCP integration test

The current MCP validation tests use application.New(nil) and therefore
cannot prove that a valid packet reaches Commit().

Add a DB-backed test through the actual production path:

    adapter.HandleTool("argus.submit_packet", ...)

using a real test CockroachDB.

The valid packet must be accepted and committed.

Verify persisted rows for:

- packet_submission
- belief
- evidence where included
- belief_edge where included
- edge_provenance where included
- task where included

Also verify origin_packet_id values.

## 3. Add malformed-packet integration test

Through the SAME real HandleTool path, submit a structurally malformed packet.

Verify:

- validation rejects it
- Persist is never reached
- zero rows are created
- no packet_submission row exists
- no belief/evidence/edge/task/provenance rows exist

## 4. Verify idempotency

Submit the same valid packet twice.

Verify:

- no duplicate packet_submission
- no duplicate beliefs
- no duplicate evidence
- no duplicate edges
- no duplicate tasks
- original origin_packet_id values remain unchanged

## 5. Verify atomic rollback

Create a packet that fails during a later persistence stage after earlier
objects would have been inserted.

Verify:

    zero packet objects remain

after rollback.

Do not consider "everything rolled back because the first belief insert
always fails" a passing atomicity test. The valid happy path must work first.

## 6. Do NOT address these yet

Do not implement:
- context snapshots
- new MCP tools
- agent attestation
- review-coverage persistence
- graph visualization
- additional provenance UI
- scenario dashboard cleanup

Those are secondary.

## 7. Verification

Run with CockroachDB actually running:

    go vet ./...
    go build ./...
    go test ./...

Then run the real local runtime:

    ./bin/argus serve

and manually submit one valid packet through MCP.

Final report must state:

1. Exact FK dependency discovered.
2. Exact ordering fix.
3. Whether packet_submission.task_id participates in an FK.
4. DB-backed MCP happy-path result.
5. DB-backed malformed-packet result.
6. Idempotency result.
7. Atomic rollback result.
8. Full test results.

Do not claim readiness until a valid packet has successfully reached Commit().
```

Once this passes, I would **stop infrastructure work again**. The next thing should be a real research packet, not another review of ARGUS mechanics.
