Plan 14 is **correct and appropriately narrow**, and the FK diagnosis is now proven rather than speculative. The schema analysis shows `packet_submission.task_id` has no FK, so inserting `packet_submission` first is safe. 

I would approve it, with **three additions before implementation**.

### 1. Add a packet-ID collision test

This is the most important missing test.

`packet_submission.packet_id` is the primary key and uses:

```sql
ON CONFLICT (packet_id) DO NOTHING
```

That is fine for an **exact retry**, but the system must not silently accept:

```text
P1 + content A
P1 + different content B
```

and then partially merge B into existing state while keeping P1's original provenance.

Add:

```text
TestMCPRejectsSamePacketIDDifferentContent
```

Expected behavior:

```text
first P1 → succeeds

second P1, different content/agent/role/hash
    → REJECT
    → zero new objects
    → original P1 provenance remains unchanged
```

A packet ID should identify one immutable packet, not merely be a best-effort deduplication key.

### 2. The atomic rollback test needs a failure that occurs inside `Persist()`

Plan 14 proposes:

> fail during edge insertion, e.g. local contradicts target. 

But a local `contradicts` target is now rejected **before persistence** by validation. That would test validation, not transaction rollback.

Use a genuine persistence-stage failure instead.

A good candidate is the existing missing-project case:

```text
valid beliefs/evidence
+
task creation
+
missing conductor_project
```

Because the task preflight occurs later in `Persist()`, the test can verify:

```text
beliefs inserted
evidence inserted
preflight fails
→ transaction rollback
→ zero rows remain
```

Or use a deliberate DB executor fault injection at a later `ExecContext`.

Keep the distinction:

```text
validation failure
    → no transaction / no writes

persistence failure
    → transaction rolls back all writes
```

### 3. The happy-path test should verify the exact provenance graph

Plan 14 already proposes checking:

```text
packet_submission
belief.origin_packet_id
evidence.origin_packet_id
edge_provenance.origin_packet_id
task.origin_packet_id
```

That's good. 

Add one assertion that all those origins equal the **same packet ID**, and that a repeated identical packet does not change them.

That directly proves:

```text
packet P
   ↓
belief B
evidence E
edge X
task T
```

all belong to P.

### One other thing I like about Plan 14

The ordering is clean:

```text
BEGIN
  packet_submission
  beliefs
  evidence
  edges + edge_provenance
  tasks
  idempotency
COMMIT
```

with `packet_submission` having no outgoing FK dependency, while all provenance-bearing objects depend on it. 

No deferred FK machinery is necessary.

### Give the coding agent this final instruction

```text
Proceed with Plan 14, with three additions:

1. Add `TestMCPRejectsSamePacketIDDifferentContent`.

Behavior:
- first packet with packet_id P succeeds
- second packet with same packet_id but different content/hash or identity
  must be rejected
- no new beliefs/evidence/edges/tasks are persisted
- original packet_submission and origin provenance remain unchanged

Do not treat `ON CONFLICT DO NOTHING` as sufficient protection against
conflicting reuse of packet_id.

2. Replace/augment `TestMCPAtomicRollback`.

Do not use a malformed `contradicts local:` reference because that now fails
during pre-persistence validation.

Use a genuine persistence-stage failure, such as a missing project for a
packet containing tasks or an explicit later-stage DB fault.

Verify:
- packet submission initially becomes part of the transaction
- beliefs/evidence/edges inserted before the failure are rolled back
- zero packet objects remain after failure

3. In the happy-path test, verify all created objects have:
    origin_packet_id == submitted packet_id

and verify an exact retry preserves those original origins.

Keep the rest of Plan 14 unchanged.

Do not add new MCP tools or architecture.

After implementation run:
    go vet ./...
    go build ./...
    go test ./...

with CockroachDB running.
```

Once these pass, **I would stop modifying infrastructure**. At that point the next test should be the substantive Work → Adversarial → Work cycle, not another bookkeeping iteration.
