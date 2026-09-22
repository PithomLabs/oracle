This plan is exactly the kind of change I would make **before the serious adversarial dry run**.

The audit establishes that `get_context` already exposes **all beliefs, evidence, and intents in the scenario**, so the only material missing piece is the graph structure.  

Plan Context Audit confirms that adding `edges` is a small projection-only change: no new tool, no migration, and no change to the MCP adapter.  

### I would approve the plan

The proposed projection is appropriate:

```json
"edges": [
  {
    "parent_id": "uuid",
    "child_id": "uuid",
    "kind": "derives"
  },
  {
    "parent_id": "uuid",
    "child_id": "uuid",
    "kind": "contradicts"
  }
]
```

The existing graph already has exactly those two edge kinds, and `submit_packet` already persists them.  

This means the full adversarial context becomes:

```text
task
dependencies
beliefs
evidence
edges
intents
debt
availability
```

That is much closer to a genuine **research-state reconstruction protocol** than the current flat projection.

### One tiny improvement

The plan calls the SQL columns `parent_id` and `child_id`, which is fine because that is the existing schema. But your agent protocol should explicitly teach the semantics:

```text
derives:
    child proposition follows from parent proposition

contradicts:
    child proposition challenges/contradicts parent proposition
```

The plan already has those definitions. 

That will reduce the chance that an adversarial agent reverses the relationship when constructing a packet.

### After this change, I would expect the adversarial agent to operate like this

```text
get_context
   |
   +-- B1 original claim
   +-- B2 supporting claim
   +-- B3 derived claim
   +-- B4 prior objection
   |
   +-- E1, E2, ...
   |
   +-- B1 -> B2 derives
   +-- B2 -> B3 derives
   +-- B4 -> B1 contradicts
   |
   v
build mental claim/graph inventory
   |
   +--> challenge B1
   +--> challenge B3
   +--> inspect whether B4 already addresses an issue
   +--> identify unsupported transitions
   |
   v
submit new adversarial beliefs
   |
   +-- B5 --contradicts--> B1
   +-- B6 --contradicts--> B3
```

That is the workflow you described earlier, and it now has the information required to perform it.

### I would not add anything else

Do not turn this into a graph-query API.

Do not add:

```text
argus.list_beliefs
argus.get_edges
argus.search_graph
```

`get_context` should remain the single reconstruction call. The plan correctly preserves that. 

And the test suite proposed in the plan is appropriate, especially scenario scoping and end-to-end `Persist → GetContext`. 

### Then run the adversarial dry run

After this small change passes, I would freeze the implementation and let the fresh adversarial agent run **without further prompting it about what it should attack**.

That is important. You want to observe whether, given:

```text
all claims
+
all evidence
+
all debt
+
all edges
+
all background documents
```

the agent naturally performs comprehensive adversarial review.

That will be a much more meaningful experiment of ARGUS itself.

The plan is therefore **approved: implement it, test it, then run the adversarial agent uncoached**.
