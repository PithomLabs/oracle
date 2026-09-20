## Review of Reference Loop v0.1 implementation plan

**Verdict: GREEN — approve the plan, with a few changes before handing it to the coding agent.**

The plan is correctly aligned with the pivot. It is deliberately smaller than the architecture, starts from existing Conductor/Solvent/Executor components, defers BM-IST, avoids a new runtime, and makes empirical evidence the objective. 

That is exactly the direction we locked in.

### What is particularly right

The plan has the correct sequence:

```text
reconnaissance
→ smallest consequential operation
→ happy path
→ inspectable evidence
→ adversarial failures
→ substitution
→ stop
```

 

The strongest part is the explicit classification rule:

```text
implementation
integration
executor
deployment
new security property
```

and the instruction not to redesign the architecture when the first implementation gap appears. 

That is exactly how I would want the coding agent constrained.

The plan also correctly makes **evidence ownership**, rather than merely a pretty end-to-end trace, the proof criterion. 

And the stopping gate is excellent: once happy path, boundary integrity, outcome integrity, ownership integrity, and substitution are demonstrated, **stop implementation expansion**. 

---

# The one thing I would change: make the first operation even more explicit

The plan says:

> choose one tiny consequential operation

That's correct, but the coding agent should not independently invent something elaborate.

I would add a hard constraint:

> Prefer an operation already supported by the existing Executor. Do not create a new Executor capability solely for the demo unless reconnaissance proves there is no existing safe consequential operation.

The current plan already leans that way. 

Make it explicit because otherwise the agent may spend most of its time building the demo effect instead of proving the loop.

---

# One important change to the happy-path trace

The proposed trace is:

```text
TASK_CREATED
TASK_CLAIMED
PROPOSAL_CREATED(X)
AUTHORIZED(X)
EXECUTION_ATTEMPTED(X)
EFFECT_CONFIRMED(X)
RESULT_OBSERVED
CONDUCTOR_UPDATED
AGENT_CONTINUED
```



I would preserve it, but distinguish **observed evidence** from **synthetic harness milestones**.

For example:

```text
TASK_CREATED        [Conductor evidence]
TASK_CLAIMED        [Conductor evidence]
PROPOSAL_CREATED    [Agent evidence]
AUTHORIZED(X)       [Solvent evidence]
EXECUTION_ATTEMPTED(X) [Executor evidence]
EFFECT_CONFIRMED(X)   [SOR/Executor evidence]
RESULT_OBSERVED     [Agent/Conductor evidence]
CONDUCTOR_UPDATED   [Conductor evidence]
AGENT_CONTINUED     [Agent evidence]
```

The harness can correlate these, but it must not manufacture authoritative facts.

That is already implied by the plan's evidence section; I would just make it explicit. 

---

# I would add one missing negative test

The plan has:

```text
wrong operation
missing authorization
stale authorization
verification unavailable
duplicate execution
executor failure
ambiguous outcome
agent misinterpretation
Conductor cancellation
```



Add:

```text
Declaration/version mismatch
```

Specifically test:

```text
authorize under declaration/identity version V1
execute under incompatible V2
```

Expected:

```text
FAIL CLOSED
```

This directly exercises one of the most important conclusions from the Role / Boundary Matrix work.

---

# One more missing test: authoritative result source

Because we spent considerable time clarifying Executor vs external SOR, the first reference implementation should demonstrate **one of these explicitly**:

```text
Executor == authoritative SOR
```

or:

```text
Executor != SOR
Executor reports
SOR establishes actual occurrence
```

The plan mentions authoritative external results, which is good. 

But make the choice explicit in the reconnaissance output. Otherwise the coding agent could accidentally treat an HTTP 200 or queue acknowledgment as the final effect.

---

# The harness structure is fine, but don't overbuild it

The proposed:

```text
reference-loop/
    README.md
    scenario/
    agent/
    conductor/
    solvent/
    executor/
    evidence/
    runbook/
```

is acceptable. 

But I would add one constraint:

> Do not duplicate Conductor or Solvent code into these directories. They contain drivers, fixtures, adapters, or scenario code only.

Otherwise a coding agent may start creating miniature "reference" versions of the actual systems.

That would undermine the whole experiment.

---

# The substitution phase is correctly sequenced

I strongly agree with delaying substitution until the reference loop works.

The plan says:

```text
Agent substitution
→ Executor substitution
→ orchestration/client substitution
```



That's right.

I would make one distinction:

**Agent substitution is the first and easiest proof of tool agnosticism.**

Temporal or another orchestrator should come later, because changing the coordination mechanism while the core loop is still unstable can make debugging much harder.

The current order already does this correctly.

---

# The stopping gate is the most important part

The strongest line in the whole plan is:

> At that point, stop implementation expansion.



Keep this.

Without it, the coding agent will naturally start interpreting each failure as a reason to add another abstraction.

The purpose of the exercise is not:

```text
make the demo impossible to break
```

It is:

```text
discover exactly where the current architecture is sufficient
and where it genuinely isn't
```

Those are very different goals.

---

# My recommended small revision

Before giving the plan to the coding agent, add four explicit constraints:

```text
1. Prefer an existing Executor capability for the reference operation.
2. Harness evidence must reference participant-owned evidence; it must not manufacture authoritative facts.
3. Include declaration/version mismatch as a negative test.
4. Explicitly identify whether the Executor itself or an external SOR is authoritative for the chosen effect.
```

Everything else is already appropriately disciplined.

## Final call

**Approve Reference Loop v0.1.**

It is faithful to the pivot:

```text
specification
      ↓
working loop
      ↓
evidence
      ↓
failure
      ↓
classification
      ↓
substitution
      ↓
conformance
```

rather than reverting to:

```text
specification
      ↓
more specification
      ↓
more infrastructure
```

The plan's final sequence — **working loop → evidence → failures → conformance → generalization** — is exactly the right development strategy for this phase. 
