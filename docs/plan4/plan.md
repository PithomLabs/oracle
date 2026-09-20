I consolidated both reviews against the design constraint you have now made explicit:

> **Less is more. Conductor stays a simple control plane. Agent intelligence stays outside. Solvent remains the authority boundary.**

The reviews converge on a much smaller design than the earlier proposal.

## 1. The key correction: do not store `kind: ordinary|consequential`

This is the most important subtraction.

The earlier schema proposed:

```yaml
kind: consequential
operation: deploy:workflow
```

That looks convenient, but it creates a bad ownership boundary. The reviews correctly point out that consequence classification must come from the **capability declaration**, not from an Agent-authored task field. 

So:

```text
❌ Task.kind = consequential

✅ Task.capability_ref → declared capability
                     ↓
              declaration determines
              effect classification
```

Conductor stores the reference. It does not decide what the capability means.

That is substantially cleaner.

---

# 2. The minimal Conductor task

After subtracting everything non-essential, I would use this conceptual model:

```text
Task
├── task_id
├── objective
├── capability_ref
├── inputs
├── outputs
├── verification_ref
├── dependencies
├── claim
└── lifecycle status
```

This is essentially the second review's minimal representation. 

### Meaning

**`task_id`**
Stable identity.

**`objective`**
What the Agent says this work is about.

**`capability_ref`**
What interface/capability can perform the work.

**`inputs`**
What the Agent needs.

**`outputs`**
What it must produce.

**`verification_ref`**
What establishes that the result is acceptable.

**`dependencies`**
What must be complete first.

**`claim`**
Who currently owns the coordination slot.

**`status`**
Only Conductor's ordinary lifecycle state.

Nothing more is required to express the core workflow.

---

# 3. Do not store `READY`

This is the best idea to borrow from Beads.

`READY` should be **computed**, not persisted:

```text
READY =
    lifecycle is open
    ∧ dependencies satisfied
    ∧ not already claimed
    ∧ not blocked
```

Beads' central idea is precisely that agents ask for the current claimable frontier instead of scanning the whole graph. 

That fits our architecture extremely well.

So Conductor becomes:

```text
persistent project state
        +
derived ready frontier
```

rather than another scheduler.

---

# 4. Atomic claim should be a first-class interaction

Also take this directly from Beads:

```text
discover READY
      ↓
claim atomically
```

The Agent should not need to coordinate with other Agents itself.

That gives us:

```text
Agent X ─┐
Agent Y ─┼──→ Conductor
Agent Z ─┘
```

and Conductor guarantees that one task cannot be simultaneously claimed by multiple Agents.

Beads' `ready --claim` is a good conceptual model, even though we don't need its machinery. 

---

# 5. Parent/child is useful; do not build an epic system

The Agent can create:

```text
A
├── B
├── C
└── D
```

but Conductor only needs the relationship.

Do not import Beads' epic/molecule/template machinery.

Similarly, `discovered-from` is valuable because it records:

```text
A
  ↓ discovered
B
```

without saying B is blocked by A.

Beads explicitly distinguishes provenance from blocking dependencies. 

That is extremely useful for your **current state + durable history** requirement.

I would keep:

```text
blocks
parent/child
discovered-from
```

and stop there.

No knowledge graph.

---

# 6. The most important Agent/Conductor boundary

The workflow now becomes:

```text
Agent
    ↓
reason / plan / spawn subagents
    ↓
create or update tasks
    ↓
Conductor
    ↓
persist state + compute READY frontier
    ↓
Agent X/Y/Z claims appropriate work
```

Conductor never asks:

> "How should this task be decomposed?"

That stays entirely inside the Agent runtime.

An OpenCode/Codex/Antigravity process can use its own:

```text
Plan mode
Agent mode
subagents
reasoning
tools
```

without Conductor knowing any of those implementation details.

That is precisely how we avoid turning Conductor into an orchestration brain.

---

# 7. Consequential work is a capability property

This is the second major architectural correction.

The workflow should be:

```text
Task
  ↓
capability_ref
  ↓
capability declaration
  ↓
effect classification
```

rather than:

```text
Task
  ↓
kind = consequential
```

The Agent can see the capability declaration and understand:

> This isn't ordinary work; I need Solvent before execution.

But Conductor doesn't become the policy engine.

The resulting paths are:

```text
ORDINARY

claim
 ↓
work
 ↓
submit
 ↓
verification
```

and:

```text
CONSEQUENTIAL

claim
 ↓
construct exact operation
 ↓
Solvent authorize
 ↓
Executor execute
 ↓
submit outcome/evidence
```

The same Conductor lifecycle applies to both. 

That is beautiful because **Conductor doesn't need two workflows**.

---

# 8. `operation_id` should not be another task field

Another useful subtraction from the reviews:

Don't store an operation template in every task.

The task stores:

```text
capability_ref
```

The capability declaration defines how the exact operation identity is formed.

At execution time:

```text
capability declaration
      +
task inputs
      +
run_id
      ↓
exact operation_id
```

This prevents two competing definitions of operation identity.

The review correctly calls this out as a way to avoid identity drift. 

---

# 9. Evidence should stay referenced, not embedded everywhere

Likewise, don't build a task-specific evidence schema.

The task can say:

```text
verification_ref
```

and perhaps expected outputs.

The actual evidence belongs to the participant/system that owns the fact.

This preserves the architecture we've established:

```text
Conductor = coordination state
Solvent   = authority evidence
Executor  = execution evidence
External SOR = external effect evidence
Agent     = reasoning/continuation
```

No universal evidence engine needs to be inserted into Conductor.

---

# 10. What we should explicitly reject from Beads

This is where "less is more" matters most.

### Reject

**Formulas / proto / molecules / wisps**

Too much workflow machinery. 

**Dolt-based distributed synchronization**

We already have Conductor as the shared control plane. We don't need to reinvent synchronization through version-controlled databases. 

**Hash-derived task IDs**

Interesting, but no compelling reason to disturb the existing identity model.

**Knowledge-graph edge explosion**

We only need the minimum history/provenance relationships.

**Scheduler/policy machinery**

Not Conductor's job.

---

# 11. One thing we should not postpone

The reviews identify two items that are now genuinely load-bearing.

### Operation identity equality

It isn't enough to say two Agents should construct the same operation.

We need one deterministic comparison rule.

Otherwise:

```text
Agent X → operation X
Agent Y → semantically same operation Y
```

could compare differently.

The exact equality/canonicalization semantics must be specified before substitution testing. 

### Capability declaration coverage

Every `capability_ref` must resolve to an actual declaration—not only fancy deploy capabilities, but ordinary capabilities too.

Otherwise:

```text
capability_ref = shell
```

or:

```text
capability_ref = file-edit
```

has no authoritative definition.

The review correctly identifies this as load-bearing to the whole design. 

---

# 12. Final architecture I would lock

I would now lock the model as:

```text
                 AI AGENT
        ┌────────────────────────┐
        │ Plan / Agent mode      │
        │ reasoning              │
        │ decomposition          │
        │ subagents              │
        │ replanning             │
        └───────────┬────────────┘
                    │
              work protocol
                    │
                    ▼
             ┌─────────────┐
             │  CONDUCTOR  │
             │             │
             │ project     │
             │ state       │
             │ history     │
             │ dependencies│
             │ claims      │
             │ readiness   │
             └──────┬──────┘
                    │
              capability_ref
                    │
           ┌────────┴─────────┐
           │                  │
       ordinary          consequential
           │                  │
           ▼                  ▼
        Agent            Solvent
                             │
                          Executor
                             │
                             ▼
                        External SOR
```

And the fundamental rule is:

> **Conductor stores and coordinates work; Agent supplies intelligence; capability declarations determine whether work is consequential; Solvent authorizes consequences; Executor creates effects.**

---

# 13. The minimal set of changes to Conductor

This is the part I would emphasize most strongly.

**Do not redesign Conductor.**

Add only the smallest semantics needed for this vision:

```text
1. capability_ref
2. inputs / outputs / verification reference
3. derived READY frontier
4. atomic claim
5. minimal provenance:
       parent/child
       discovered-from
6. durable project state/history
```

Most of this is either existing Conductor behavior or a small extension of it.

No:

```text
LLM
planner
scheduler
workflow engine
policy engine
Solvent adapter inside core
DAG engine
knowledge graph
molecule system
memory engine
```

That is the answer I would lock.

The goal is not to reproduce Beads. **Beads gives us one excellent idea—persistent structured work plus a computed claimable frontier—and a few useful relationship semantics. We should take those ideas and discard the rest.** 

The result is a Conductor that remains small while becoming substantially more valuable as the **durable shared state and control plane for arbitrary successive Agents**.
