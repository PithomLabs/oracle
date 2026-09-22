This is the point where I would say **the production agent boundary has been repaired sufficiently to move back to the research experiment**, with one remaining verification gate: the DB-backed test suite.

The important result is this:

```text
AI packet
   ↓
argus.submit_packet
   ↓
canonical packet validation
   ↓
malformed packet
   ↓
REJECT
   ↓
Persist() never reached
   ↓
ZERO database mutation
```

The report explicitly shows six malformed packets exercised through the real `HandleTool("argus.submit_packet")` path, all rejected before persistence. That directly closes the most serious review finding.

### What is now solid

The report says the following invariants are enforced:

```text
agent.id / role / harness / model required
agent.role == packet.role
packet.role ∈ {work, adversarial}
global local_id uniqueness
valid claim types
valid edge kinds
valid references
pack_ref uses @
required project exists
packet persistence atomic
```

That is a significant improvement because these are now enforced at the **actual agent boundary**, not just by dead or isolated validators.

The specific regression result:

```text
91/91 non-DB tests PASS
```

plus:

```text
6 malformed packets
0 reached Persist()
0 DB mutation
```

is particularly strong evidence.

### But don't call the whole system green yet

The report still says:

```text
DB-dependent tests require CockroachDB
```

So the final verification should be:

```bash
./bin/argus serve
```

then in another terminal:

```bash
go test ./...
```

That is the last technical gate before the next research run.

I'd also manually verify one legitimate packet after startup:

```text
valid Work packet
    → accepted
    → beliefs/evidence/edges/tasks persist
    → provenance persists
```

and one invalid packet:

```text
invalid Work packet
    → rejected
    → no state change
```

### I would now stop infrastructure expansion

This is where your earlier EBP concern becomes important.

We've spent enough effort making the ledger trustworthy. The recent work has closed concrete failures:

```text
validation bypass
duplicate IDs
invalid claim/edge types
provenance loss
partial task persistence
edge integrity
```

Those were justified.

I would now **freeze the bookkeeping architecture temporarily**.

Don't immediately add:

* context snapshots
* review-coverage tables
* epistemic-kind taxonomy
* agent attestation
* more UI lineage
* more provenance columns

unless the next actual research cycle demonstrates a concrete failure requiring them.

### Then do the real experiment again

The next sequence should be:

```text
Clean ARGUS state
      ↓
Work Agent
      ↓
Deliverable A / further research
      ↓
Fresh Adversarial Agent
      ↓
comprehensive challenge
      ↓
Fresh Work Agent
      ↓
respond to challenges
      ↓
human adjudication
```

And now, unlike the first run, you can trust that a malformed packet won't silently contaminate the ledger.

One thing I would watch closely in that next cycle is **whether the agents spend their time doing domain reasoning or bookkeeping**. At this stage, that is the more important architectural test.
