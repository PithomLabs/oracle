I would rate the attached plan **AMBER**. The overall direction is correct, and I like the discipline of forcing the work through requirements → contract → conformance → scenarios → POC. But there are several places where the plan could accidentally reintroduce exactly the infrastructure sprawl and semantic ambiguity that we just eliminated.

The good part first: the plan correctly keeps the work to six artifacts, explicitly says no new runtime/database/UI/SDK initially, and applies the Growth Gate only after a real POC failure.  That is the right foundation.

## Finding 1 — HIGH: “Human intervention at any stage” is too strong

The requirements say:

> R6 Human intervention at any stage
> R7 Human review at any stage 

That is not presently a safe architectural requirement.

A human can potentially intervene before a consequential action, reject work, alter coordination, or review an outcome. But **during an external effect**, the system may not have a reversible control at all.

For example:

```text
Executor → sends payment → external bank
```

A human cannot necessarily "pause" that effect merely because the workflow specification says humans may intervene at any stage.

Likewise, the current Conductor lifecycle does not have a generic `paused` state.

The plan later proposes:

> Agent working → human intervenes → task redirected/rejected/paused → agent resumes 

That "paused" concept is not obviously supported by the current Conductor model.

### Fix

Change the requirement to something like:

```text
R6 Human intervention is supported at defined workflow boundaries.
R7 Human review may occur wherever the applicable role contract permits it.
R8 Intervention cannot retroactively invalidate an already-produced external effect.
```

Then explicitly classify intervention points:

```text
before work
during coordination
before consequential authorization
after authorization but before execution, when the executor supports it
after execution
```

That is much more rigorous than “any stage.”

---

## Finding 2 — HIGH: The Workflow Contract risks becoming a second state machine

Step 2 says:

> "For each transition we record current state, actor, action, required evidence, destination, who may intervene, what intervention is allowed." 

This is potentially dangerous.

Conductor already owns a lifecycle state machine.

Solvent already owns authoritative state around target/snapshot/intent/authorization/revocation.

If the Workflow Contract starts defining another persistent state machine with its own canonical states, you may end up with:

```text
Workflow state
Conductor task state
Solvent authority state
Executor state
```

and eventually need synchronization between them.

That would be exactly the sort of fifth system we decided not to build.

### Fix

Explicitly state:

> **Workflow phases are a specification-level vocabulary, not a new persisted state model.**

For example:

```text
DISCOVER
FORMULATE
WORK
REVIEW
```

are conceptual phases.

They are **mapped onto existing Conductor lifecycle and external events**, not persisted as a new `workflow_state`.

The requirements should explicitly forbid:

```text
workflow_state column
workflow state machine
workflow event store
workflow authority state
```

unless a later Growth Gate proves one is necessary.

This is probably the most important correction to make before drafting the Workflow Contract.

---

## Finding 3 — HIGH: The role matrix has ambiguous authority wording

The matrix says:

> Human — "May exercise through Solvent"

and:

> Executor — "Enforces" under Authority. 

The first phrase is ambiguous about who actually owns authority.

We already have a clean invariant:

```text
Solvent = authority
Executor = effect
Human = actor/principal capable of requesting/reviewing/intervening
```

The human may be the principal whose authority Solvent evaluates, or a person operating an agent/executor, but **Human should not appear to be an alternate authority engine**.

Likewise, Executor should not be described as "enforcing authority." It should enforce the supplied execution constraint, while Solvent remains the authority source.

### Fix

Make the matrix more precise:

```text
Authority:
Agent       requests
Conductor   coordinates
Solvent     decides
Executor    requires/enforces supplied authorization
Human       may request/review/approve only through the applicable authority boundary
```

That preserves the ownership model.

---

## Finding 4 — MEDIUM: T06 is underspecified and could force a workflow runtime

The conformance tests include:

> T06 long-running execution does not block workflow globally 

That sounds reasonable, but it silently assumes a particular execution model.

What does "does not block workflow globally" mean?

Could mean:

```text
other tasks continue
```

or:

```text
the Agent can continue reasoning
```

or:

```text
Conductor remains responsive
```

or:

```text
execution is asynchronous
```

Those are very different requirements.

The first three can probably be demonstrated using current infrastructure. The last one may force an execution orchestration design.

### Fix

Define it in terms of coordination, not runtime mechanics:

> A long-running external effect must not prevent unrelated Conductor work from progressing.

Then explicitly state:

> This does not require Conductor to own scheduling, execution, polling, retries, or workflow runtime semantics.

That keeps Temporal and other execution frameworks outside the core.

---

## Finding 5 — MEDIUM: T10 is also too implementation-dependent

The test says:

> T10 human can intervene during execution when supported. 

"When supported" makes the test non-deterministic.

A conformance test should establish a concrete contract.

The correct question is not:

> Can a human always interrupt execution?

It is:

> When an Executor exposes an interruption/control boundary, does the workflow preserve the role boundaries and correctly record the resulting state?

That allows executors with different capabilities without making the workflow contract dishonest.

---

## Finding 6 — MEDIUM: Named GPT/Claude tests are the wrong level of abstraction

The agent tests are:

```text
T13 GPT behavioral adherence
T14 Claude behavioral adherence
T15 scripted-client adherence
```



I would change those.

The whole point of tool agnosticism is that the conformance contract should not care whether the agent is GPT, Claude, or something else.

Otherwise the conformance suite slowly becomes:

```text
provider compatibility test suite
```

instead of a workflow conformance suite.

### Better

Use:

```text
T13 autonomous agent client
T14 alternate agent client
T15 non-agent deterministic client
```

Then instantiate them with:

```text
GPT
Claude
script
human/curl
```

as test implementations.

The test asserts:

> different clients can satisfy the same protocol.

It does not assert that a particular model "behaves correctly" in the abstract.

The plan itself correctly says stochastic behavior should be measured rather than treated as deterministic conformance; I'd push that principle further and remove provider names from the normative test definitions. 

---

## Finding 7 — MEDIUM: The conformance suite risks duplicating Solvent's own test suite

T02–T05 are sensible:

```text
consequential request cannot bypass authority
denied authorization cannot execute
authorization != execution
execution failure != authorization failure
```



But the Loop Engineering suite should test the **cross-role boundary**, not re-test the Solvent kernel internally.

For example:

```text
Loop conformance:
Agent requests consequential action
→ Solvent decides
→ denied
→ Executor is not invoked
```

Good.

But:

```text
Loop conformance:
Does Solvent correctly bind target_id + snapshot_id?
```

That belongs to Solvent's own conformance/security suite.

Otherwise Loop Engineering becomes a second Solvent test suite.

Add an explicit rule:

> Loop Engineering tests role-boundary behavior and handoffs; component-specific correctness remains the responsibility of the component's own test suite.

---

## Finding 8 — LOW: “Long-running execution supported” should be qualified

R12 says:

> Long-running execution supported. 

That's reasonable for BM-IST, but "supported" should not be interpreted as:

```text
Conductor must orchestrate long-running processes.
```

It should mean:

> The workflow can represent and continue around an external operation whose duration exceeds a normal request/response interaction.

The Executor or external workflow engine can own the actual long-running machinery.

That distinction should be explicit in the requirements.

---

# What I would keep unchanged

The core requirements are excellent:

```text
R1  Domain-agnostic
R2  Tool/runtime-agnostic
R3  Four roles remain distinct
R9  Authorization ≠ execution
R10 Activity ≠ execution proof
R11 Domain truth remains outside infrastructure
R14 No new persistent workflow authority state
R15 Existing Conductor UI preferred
R16 No new runtime/framework unless demonstrated necessary
```



And the non-goals are exactly the right guardrails, especially:

```text
No workflow engine
No second task system
No workflow-specific DB
No workflow-specific UI
No persisted workflow Intent
```



The last one should absolutely stay.

---

# The three BM-IST scenarios are well chosen

These are a good initial validation set:

**A — ordinary research**, explicitly without Solvent.
**B — consequential shared computation**, exercising the authority boundary.
**C — human intervention**, exercising semi-autonomous operation. 

That gives good coverage without trying to build "the entire BM-IST system."

I would just rewrite Scenario C to avoid promising `pause` unless we define what that means in the actual existing Conductor lifecycle.

---

# Most important change to the plan

The plan currently has this implicit model:

```text
Workflow Contract
    ↓
new state/transition model
```

I would change it to:

```text
Workflow Specification
    ↓
cross-role interaction contract
    ↓
maps onto existing systems
```

So:

```text
Loop phase
    ↓
Conductor task/activity
    ↓
optional Solvent authorization
    ↓
Executor effect
    ↓
result
    ↓
Agent interpretation
```

There is **no workflow database and no workflow state authority**.

---

# Adversarial verdict

### **AMBER — good plan, but tighten before Deliverable #1**

The plan's strategic direction is right. The risks are mainly **semantic creep**, not architectural failure.

Before writing `Loop Engineering Requirements v0.1`, I would amend five things:

1. Replace **"human intervention at any stage"** with bounded intervention points.
2. Explicitly declare Workflow phases **non-persistent and non-authoritative**.
3. Remove ambiguity around Human/Executor authority in the role matrix.
4. Reframe T06/T10 around observable coordination behavior rather than requiring a runtime capability.
5. Replace GPT/Claude-specific normative tests with **client-class conformance tests**, using GPT/Claude as implementations.

With those corrections, the proposed sequence:

```text
Requirements
    ↓
Workflow Contract
    ↓
Role Matrix
    ↓
Conformance
    ↓
3 BM-IST scenarios
    ↓
POC
    ↓
Growth Gate
```

is exactly the right way to approach this. 

Most importantly, **do not start coding after Deliverable #1**. The plan is right to establish the specification, then derive conformance scenarios, then prove the current Conductor/Solvent architecture can carry them before earning any new implementation. That is the right application of the Growth Gate to the BM-IST experiment itself.
