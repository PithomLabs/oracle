This audit answers the key question very clearly:

**Yes, the current `argus.get_context(task_id)` exposes every belief in the current research scenario to the adversarial agent.** It is not merely returning beliefs attached to the task. The production SQL is explicitly:

```sql
SELECT id, claim, claim_type, status, debt::STRING, final_truth
FROM belief
WHERE scenario_id = $1::UUID
ORDER BY claim
```

with no task, project, edge, or status filtering.  

So your adversarial workflow is currently:

```text
Fresh Adversarial Agent
        ↓
argus.get_context(task_id)
        ↓
ALL beliefs in scenario
ALL evidence
ALL intents
        ↓
read docs/corpus/*
        ↓
analyze/challenge claims
        ↓
submit adversarial packet
```

That means the agent **can see and attack every existing belief ID in the current scenario**. 

### But there is one significant limitation

The audit found that **edges are missing from the returned context**.

The system persists `belief_edge` records, but `GetSnapshot` does not retrieve them, and `Snapshot` has no `Edges` field. 

So the adversarial agent currently sees:

```text
Beliefs       ✓
Evidence      ✓
Intents       ✓
Debt          ✓
Edges         ✗
```

This matters a lot for the lifecycle you're building.

Suppose your Work Agent produced:

```text
B1 = foundational claim
B2 = supporting proposition
B3 = derived proposition
B4 = prior objection
```

with:

```text
B2 --derives----> B1
B3 --derives----> B2
B4 --contradicts-> B1
```

The adversarial agent currently sees B1–B4, but **not the relationships**. So it has to reconstruct the logical structure from the prose rather than from the actual ARGUS graph. 

### My recommendation for the dry run

I would **not stop the dry run because of this**, but I also would not consider the adversarial agent's result a full test of ARGUS's graph-aware adversarial capability.

You have two choices:

**Run now.**

This tests:

> Can a fresh adversarial agent enumerate all current claims and independently challenge them?

That is already a valuable experiment.

**Or add `edges` to `get_context` first.**

That would test the stronger proposition:

> Can a fresh adversarial agent reconstruct the actual epistemic graph, identify which claims support or contradict which other claims, and attack the structure rather than merely individual prose claims?

Given that **your entire lifecycle is explicitly graph-based**, I lean toward the second.

Importantly, this does **not** require a third MCP tool or a new subsystem. It is simply completing the existing RCP projection. The audit already confirms that edges are already persisted; they are just not projected. 

### What I would expect from the adversarial agent

With the current implementation, I would expect it to produce something like:

```text
B1 — G0
    Challenge: YES
    ...

B2 — substrate definition
    Challenge: YES
    ...

B3 — encoding
    Challenge: NO MATERIAL OBJECTION
    ...

B4 — derived consequence
    Challenge: YES
    ...
```

and submit new adversarial beliefs with:

```text
B5 --contradicts--> B1
B6 --contradicts--> B3
```

using the existing belief IDs as targets.

That part is supported by the current capability model: the agent sees all beliefs, and the packet grammar supports `canonical:belief:<uuid>` and `contradicts`. 

What it **cannot currently do well** is say:

> "B3 depends on B2, which depends on B1, so attacking B1 invalidates B3."

because it does not receive the `derives` edges.

### I would make one small change before the adversarial run

Have the coding agent audit the exact effort required to add `Edges` to `get_context`. Not to redesign RCP—just to expose the already-persisted graph.

Use this:

```text
Audit the existing ARGUS RCP/get_context implementation for exposing the
persisted belief edge graph.

Do not change code yet.

We have verified:
- get_context returns ALL beliefs in the current scenario
- get_context returns ALL evidence
- get_context returns ALL intents
- belief_edge records are persisted
- but Snapshot currently has no Edges field and GetSnapshot does not query
  belief_edge

Determine the minimum implementation needed to expose the existing edge
graph through argus.get_context.

Answer:

1. Exact current belief_edge schema.
2. Exact existing edge kinds.
3. Exact write path used by submit_packet.
4. Whether edge data can be added to Snapshot without changing the MCP tool
   surface.
5. Exact response shape you recommend.
6. Whether this is a small projection-only change or requires architectural
   changes.
7. Tests needed to prove:
   - derives edges are returned
   - contradicts edges are returned
   - edges reference the correct immutable belief IDs
   - no authority capability is added

Do not implement yet.
```

### Bottom line

Your most important question is now answered:

**Yes, the adversarial agent sees all existing claims in the current scenario.** 

The remaining question is whether you want your first adversarial dry run to test **claim-level adversarial review** or **full graph-aware adversarial review**.

Given the lifecycle you've designed, I would make `edges` visible before the serious adversarial run. The missing edge projection is the one remaining gap between “the agent can see the research state” and “the agent can actually reconstruct the research argument.”
