The pivot changes **what we are delivering first**, not the underlying architecture.

## The deliverables now

I would define the pivot as producing **five concrete artifacts**, in this order.

### 1. Reference Loop — the working end-to-end proof

This is the primary deliverable.

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
External system
  ↓
Result
  ↓
Conductor
  ↓
Agent
```

It should be a thin, runnable demonstration using the **existing** Conductor, Solvent, and an existing/minimal Executor integration—not a new framework.

The goal is to prove the happy path and show that the four responsibilities actually work together.

---

### 2. Evidence Package — proof that the loop really happened

The reference loop must leave enough participant-owned evidence to reconstruct:

```text
task created
→ task claimed
→ operation proposed
→ exact operation authorized
→ exact operation attempted
→ external effect occurred
→ result observed
→ Conductor updated
→ Agent continued
```

This is deliberately **not** a new workflow ledger.

The deliverable is the evidence and a way to inspect/correlate it.

This is arguably the most important artifact because it distinguishes:

> "we have a demo"

from:

> "we proved the architecture."

---

### 3. Adversarial Reference Scenarios — prove the boundaries fail correctly

After the happy path:

```text
authorized X → execute Y
missing authorization
stale authorization
verification unavailable
declaration/version mismatch
duplicate/replay
execution failure
ambiguous outcome
termination/revocation
Conductor cancellation
Agent misinterpretation
```

The deliverable is **testable behavior + evidence**, not another abstract specification.

Every failure gets classified as:

```text
implementation
integration
executor
deployment
new security property
```

Only the last category can trigger architectural reconsideration.

---

### 4. Substitution Proof — demonstrate tool/framework agnosticism

Once the reference loop works:

```text
Agent A → Agent B
```

then:

```text
Executor A → Executor B
```

then:

```text
one orchestration/client mechanism → another
```

The deliverable is evidence that the **semantic contract survives substitution**.

This is where we actually substantiate the broader Loop Engineering claim.

---

### 5. Conformance Test Matrix — derived from reality

This comes **after** the reference loop and adversarial scenarios.

The matrix should now be based on:

```text
observed implementation
+
observed failure modes
+
locked boundary contract
```

rather than continuing to speculate about every conceivable failure before implementation.

The existing Role / Boundary Matrix and Workflow Specification become **inputs to the test design**, not the next thing to endlessly expand.

---

# So the pivot produces this stack

```text
1. Reference Loop
        ↓
2. Evidence Package
        ↓
3. Adversarial Scenarios
        ↓
4. Substitution Proof
        ↓
5. Conformance Test Matrix
```

Then, after those are solid:

```text
6. BM-IST Validation Scenario
7. Agent Skill
8. Conformance Harness
9. POC Runbook / Demo
```

Those later artifacts still exist; they are simply **downstream of empirical proof**.

## What is explicitly NOT a deliverable of this pivot

Not now:

```text
Loop Engineering runtime
Workflow engine
Workflow database
Workflow event bus
New scheduler
New gateway
Merged Conductor/Solvent
Second authority engine
Generic reconciliation daemon
New Solvent kernel functionality
BM-IST infrastructure
```

That is consistent with the implementation plan, whose objective is one thin reference-loop demo using the existing components. 

## The simplest way to think about it

Before the pivot, the work was primarily:

```text
"Define what this architecture means."
```

Now it is:

```text
"Prove that this architecture works."
```

So the **Reference Loop is the centerpiece**, and everything else in this phase exists to answer one question:

> **Can Agent → Conductor → Solvent → Executor → external system → Result → Conductor → Agent work in reality, with the boundaries surviving failure and substitution?**

Once that is demonstrated, we have earned the right to generalize.
