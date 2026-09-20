After consolidating both reviews with our design rule — **subtract until there is nothing left to subtract** — I think the two-mode model is the right resolution, but one important distinction must be made:

> **Plan approval governs what work the human has agreed to. Solvent governs whether a specific consequential operation may actually occur.**

The two gates operate at different grains and should never be merged. The attached review captures this cleanly. 

## 1. Lock the two modes

There are exactly two Agent modes:

```text
PLAN
WORK
```

Nothing else.

### PLAN

Agent may:

```text
inspect
reason
decompose
spawn subagents
iterate
revise
propose
```

The Agent presents the plan to the human.

The human decides:

```text
APPROVE
REJECT
```

Conductor does **not** orchestrate plan iteration.

That remains entirely:

```text
Human ↔ Agent
```

This cleanly prevents Conductor from becoming a planner/orchestration brain. 

### WORK

After approval:

```text
Agent enters WORK mode
```

The Agent may execute ordinary work and interact with Conductor.

Conductor now owns the shared durable project state:

```text
current state
history
tasks
dependencies
claims
provenance
```

---

# 2. Human approval is not authority

This should be written into the protocol as an invariant:

> **Plan approval authorizes the Agent to enter Work Mode; it does not authorize any consequential external operation.**

That is essential because otherwise Solvent becomes ceremonial.

The attached review explicitly identifies this distinction. 

So:

```text
Human APPROVES PLAN
        ≠
Solvent AUTHORIZES OPERATION
```

The human decides:

> “Yes, this is the work I want pursued.”

Solvent decides:

> “Yes, this exact consequential operation is authorized now.”

Those are different authorities at different levels.

---

# 3. The task graph remains simple

Conductor does not need an ordinary/consequential task fork.

Delete:

```text
kind: ordinary
kind: consequential
```

and:

```text
READY_FOR_AUTHORIZATION
```

The work graph is simply:

```text
Task
├── objective
├── capability_ref
├── inputs
├── outputs
├── verification_ref
├── dependencies
├── claim
└── lifecycle
```

`READY` is derived.

The Agent creates/decomposes the graph.

Conductor persists and coordinates it.

---

# 4. Consequentiality belongs at the operation boundary

This is the most important semantic refinement.

The rule becomes:

> **Anything whose successful completion intentionally mutates authoritative state outside the workflow is consequential and must pass through Solvent.**

That means:

```text
ordinary internal work
    ↓
Agent
```

and:

```text
external mutation
    ↓
Solvent
    ↓
Executor
```

The task itself is not "consequential."

The **operation** is.

This resolves the category problem identified in the review. A single task can contain many ordinary actions and eventually produce one consequential operation. 

---

# 5. No Conductor policy router

This is where I would reject the more complicated proposals from the adversarial review.

Conductor must not ask:

> “Is this task dangerous?”

It must not enforce:

> “This task is consequential, therefore route to Solvent.”

It must not become a second policy engine.

Instead:

```text
Agent
  ↓
attempts external-effect operation
  ↓
Solvent boundary
  ↓
deny or authorize
```

That preserves the clean ownership:

```text
Agent     = intelligence
Conductor = shared work state
Solvent   = authority
Executor  = effect
```

---

# 6. But the consequence boundary MUST actually be enforced

This is the strongest point from the other review that I would retain.

A rule that says "all external effects go through Solvent" is useless if the Agent can simply call a raw external API.

So the **actual external-effect-capable tools available to Work Mode must terminate at an enforcement boundary that requires valid Solvent authorization**.

Conceptually:

```text
Agent
  ↓
external-effect capability
  ↓
Solvent checkpoint
  ↓
Executor
  ↓
External SOR
```

The enforcement belongs at the capability/executor boundary, not in Conductor.

That satisfies the security concern without contaminating Conductor. The review's tool-interception proposal is therefore directionally correct, but it should remain an **execution/deployment boundary**, not become another Conductor subsystem. 

---

# 7. Do we need plan versions and operation templates?

Here I would take a **minimal** position.

We should record:

```text
plan_id
plan_version
approval decision
approver
```

because otherwise durable project history cannot tell which approved plan led to the work.

But I would **not require a giant operation-inventory system yet**.

The plan should contain enough information for a human to understand the intended scope and consequential actions, but exact operation identity is still formed at runtime and checked by Solvent.

Why?

Because forcing the Agent to pre-enumerate every exact future operation can turn Plan Mode into another workflow engine.

The better minimal model is:

```text
Human approves Plan vN
        ↓
Agent works
        ↓
exact consequential operation arises
        ↓
Solvent checks exact operation
```

The relationship is recorded:

```text
plan_id + operation_id
```

so later evidence can answer:

> Which approved plan produced this consequential operation?

The review's proposal to make plan/operation conformance observable is valuable; I would treat automatic prevention as **not yet earned**. 

---

# 8. This yields a much cleaner workflow

```text
                    HUMAN
                      │
                 gives objective
                      │
                      ▼
               ┌────────────┐
               │   AGENT    │
               │ PLAN MODE  │
               └──────┬─────┘
                      │
                 plan iteration
                 stays here
                      │
                 HUMAN APPROVES
                      │
                      ▼
               ┌────────────┐
               │   AGENT    │
               │  WORK MODE │
               └──────┬─────┘
                      │
              ordinary work
                      │
                      ▼
                CONDUCTOR
          durable project state
          + history + work graph
                      │
                READY frontier
                      │
             Agent X / Y / Z
                      
                  consequential
                      │
                      ▼
                   SOLVENT
                      │
                  authorize
                      │
                      ▼
                  EXECUTOR
                      │
                      ▼
                 External SOR
```

This is much simpler than the earlier designs.

---

# 9. What Conductor actually does

Conductor remains boring.

That is a feature.

It does:

```text
persist
discover
claim
release
update
block
submit
record history
compute READY
```

It does **not**:

```text
plan
reason
decompose
authorize
execute
route policy
evaluate domain semantics
interpret Solvent
```

The attached review is right that the whole architecture becomes cleaner once these responsibilities are kept separate. 

---

# 10. What happens when Work Mode discovers something new?

This is an important edge case.

Suppose:

```text
Plan A approved
```

and during Work Mode the Agent discovers:

```text
Task E
```

The Agent may create/update work state in Conductor, but the **human remains the authority over plan iteration**.

So there are two cases:

```text
Tactical adjustment
    → Agent may perform within approved scope

Structural plan change
    → return to PLAN
    → revise
    → human approves
```

I would keep that rule deliberately simple.

Do not create a six-level approval taxonomy.

The review correctly identifies approval fatigue as the danger; the answer is simply separating tactical execution from structural plan changes. 

---

# 11. Multi-Agent participation now becomes almost automatic

Agent A creates the project plan.

After approval:

```text
Conductor

Task B READY
Task C BLOCKED by B
Task D READY
```

Agent X can claim B.

Agent Y can claim D.

Agent Z can join later and discover whatever becomes READY.

They don't need to inherit the previous Agent's context window.

They inherit:

```text
current state
+
durable history
+
work graph
```

This was the whole reason to make Conductor the durable control plane.

---

# 12. Final subtraction

After consolidating everything, I would explicitly **reject**:

```text
❌ task consequence flags
❌ Conductor policy routing
❌ Conductor authorization
❌ cryptographic Conductor↔Solvent ledger
❌ plan workflow engine
❌ autonomous Conductor replanning
❌ READY_FOR_AUTHORIZATION state
❌ formulas/molecules/DAG runtime
❌ claim leases for now
❌ Conductor verification engine
❌ operation inventory enforcement for now
```

And keep only:

```text
✓ PLAN / WORK
✓ human plan approval
✓ durable project state
✓ durable history
✓ task graph
✓ dependencies
✓ provenance
✓ READY frontier
✓ atomic claim
✓ capability reference
✓ Solvent checkpoint for consequential operations
✓ Executor for effects
✓ plan_id ↔ operation_id correlation
```

## The final architecture sentence

I would lock this:

> **The Agent owns intelligence and has exactly two modes—Plan and Work. Humans approve plans and their iteration. Conductor owns durable project state and coordination. Any operation that mutates authoritative state outside the workflow is consequential and must pass through Solvent immediately before execution. Executor alone produces the external effect.**

That is the smallest design I've seen emerge from all of these reviews that still preserves the security and control properties we care about.

And most importantly, **Conductor remains Conductor**. We are adding durable project memory and a claimable work frontier—not turning it into another AI, workflow engine, policy engine, or harness.
