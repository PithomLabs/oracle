This audit is strong. It found the **real next layer of integrity problems**, and I would proceed toward implementation—but I would make **four corrections before coding**.

The most important conclusion is that the current system has an **incomplete provenance spine**: `packet_submission` knows who produced a packet, but the resulting beliefs, evidence, edges, and tasks do not retain a durable link back to that packet. 

## 1. Provenance plan: approved, with one integrity improvement

The proposed model is sound:

```text
packet_submission
       |
       +--> belief
       +--> evidence
       +--> edge
       +--> task
```

with `origin_packet_id` on belief/evidence/task and an `edge_provenance` relation for edges. 

But make `origin_packet_id` a **foreign key to `packet_submission(packet_id)`**.

Since `packet_id` is now correctly a `STRING`, use:

```sql
origin_packet_id STRING
    REFERENCES packet_submission(packet_id)
```

nullable for objects that did not originate from an agent packet, such as seed or human-created state.

This gives you:

```text
packet_submission P1
        |
        +---- belief B1
        +---- evidence E1
        +---- edge X1
        +---- task T1
```

without allowing a dangling origin reference.

For packet-created objects, `origin_packet_id` should be non-null. For seed/human lifecycle objects, `NULL` has a legitimate meaning.

## 2. The task persistence bug is more urgent than the audit makes it sound

This is the one item I would elevate.

The audit says packet-created tasks use:

```text
project_id = ''
```

and that the insertion likely fails while the packet otherwise succeeds. 

That would violate one of your core expectations:

```text
submit_packet accepted
    ≠
partial packet persistence
```

A packet should not be able to "succeed" while some of its declared objects disappear.

The implementation must ensure the whole persistence operation is atomic:

```text
packet
  ├─ beliefs
  ├─ evidence
  ├─ edges
  └─ tasks
        |
        v
transaction
        |
   ALL SUCCESS
        |
        v
commit
```

or:

```text
ANY FAILURE
    ↓
ROLLBACK
    ↓
submit_packet fails
```

So I would explicitly add:

> **Any failure persisting any packet object must fail the packet submission and roll back the transaction. No partial success.**

And before using `pkt.ScenarioID` as `project_id`, verify that the current schema actually defines scenario ID and project ID as the same identifier. The audit strongly suggests this is the intended mapping, but the implementation should establish it rather than assume it.

## 3. Edge validation: good, but keep it narrow

The audit found:

* cross-scenario edges allowed;
* local self-challenges allowed;
* self-edges allowed;
* packet-level edge validation is incomplete on the MCP path. 

I agree with fixing these, but I would make the rules exactly:

```text
All edge endpoints:
    must resolve
    must belong to same scenario
    must not equal each other

contradicts:
    target must be a canonical existing belief
```

I would **not** require the target to be "unresolved" or "have open debt." A promoted claim can legitimately be challenged later, and a retracted claim can remain a valid historical target.

The key rule is structural:

```text
adversarial new claim
        |
        | contradicts
        v
existing canonical belief
```

That matches the lifecycle you're building.

The plan's proposed validation is therefore good with the self-edge rule made explicit. 

## 4. Put the EBP calibration rule in the shared protocol

Definitely keep this:

> Open debt alone is not an EBP violation when a claim is merely entered. 

This was the most instructive failure from the actual adversarial run.

The distinction should be:

```text
entered + debt
    = normal candidate state

unresolved debt
    = blocks promotion

rule explicitly violated
    = process violation

contradiction
    = research challenge

retraction
    = authorized lifecycle transition
```

That will improve both Work and Adversarial agents.

## 5. The coverage rule is worth adding now

I agree with the audit's decision not to add a persistence subsystem yet.

Add to the common protocol:

```text
For every belief visible through get_context, record one disposition:
- challenged
- no material objection
- insufficient basis to assess
```

That gives you a measurable adversarial review without forcing the database to carry review-state yet. 

This distinction is important:

```text
not challenged
    ≠
reviewed and no objection
```

## 6. The declared identity model is correct

No changes needed there.

The audit correctly treats:

```text
agent_id
harness
model
```

as **declared provenance, not attested runtime identity**. 

That is appropriate for this POC.

## 7. Context snapshot should remain deferred

Agreed.

The absence of a snapshot/version ID doesn't undermine this sequential dry run. It becomes important later when multiple agents operate concurrently and you want a packet to say:

```text
"I reviewed ARGUS state S17."
```

Keep it as named future debt, not current implementation.

---

# One change I would make to the implementation plan

The current Phase 1 says:

> add `origin_packet_id STRING` to belief/evidence/conductor_task. 

Make that:

```text
origin_packet_id STRING NULL
REFERENCES packet_submission(packet_id)
```

and explicitly enforce:

```text
packet-created object → origin_packet_id NOT NULL
seed/human-created object → origin_packet_id NULL
```

Then make the packet submission transaction atomic.

That gives you the provenance invariant:

```text
No orphaned research objects.
No orphaned provenance.
No silently lost packet components.
```

### Final disposition

I would **approve the plan for implementation after these corrections**:

1. `origin_packet_id` references `packet_submission.packet_id`.
2. Packet persistence is all-or-nothing; task insertion failure cannot produce partial success.
3. Verify `ScenarioID → project_id` mapping before using it.
4. Edge rules: same scenario, no self-edge, `contradicts` target canonical existing belief.

Everything else in the audit is properly scoped. The proposed implementation remains small—migration, persistence, validation, protocol, UI, and tests—without adding another service or authority mechanism. 

The key outcome is that after this pass, you should be able to trace:

```text
Agent
  ↓
Packet
  ↓
Belief / Evidence / Edge / Task
  ↓
Adversarial relationship
  ↓
Human adjudication
```

That is the provenance spine ARGUS has been missing.
