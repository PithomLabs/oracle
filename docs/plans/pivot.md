# BACKGROUND CONTEXT — LOOP ENGINEERING PIVOT

## Why this work exists

We are developing a domain-agnostic architecture for autonomous/semi-autonomous software work in which four responsibilities remain explicitly separate:

    Agent       = Agency
    Conductor   = Coordination
    Solvent     = Authority
    Executor    = External Effect

The governing invariant is:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Conductor and Solvent are intentionally separate systems.

Do NOT merge them.

Do NOT create a new Loop Engineering runtime.

Do NOT add a workflow engine, workflow database, workflow event store, generic scheduler, second authority engine, or mandatory gateway.

Loop Engineering is a protocol/contract describing how the existing components interact. It is not another service.

---

## The development strategy has now pivoted

We spent substantial effort formalizing Loop Engineering through:

    Workflow Specification
    → Role / Boundary Matrix
    → Conformance concepts

That work was useful because it exposed and clarified the security and ownership boundaries.

However, we are now deliberately pivoting from:

    "keep refining the abstract specification"

to:

    "prove the architecture with one small real end-to-end workflow"

The specification should now serve as a constraint on the implementation, not become the implementation project itself.

The next objective is empirical proof.

We want to determine whether the architecture actually works when real components interact, fail, and are substituted.

---

# THE REFERENCE LOOP WE WANT TO PROVE

Build and exercise the smallest meaningful loop:

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

This is the reference interaction model.

It should be possible to trace one piece of work through this loop and answer:

    1. What work exists?
    2. Who is responsible for it?
    3. What exact consequential operation is proposed?
    4. Was that exact operation authorized?
    5. What exact operation was actually attempted?
    6. What actually happened externally?
    7. How did Conductor observe the coordination result?
    8. How did the Agent interpret the result?

---

# RESPONSIBILITY OF EACH COMPONENT

## Agent

The Agent owns:

    agency
    reasoning
    proposal
    work
    interpretation

The Agent does NOT become authoritative merely because it proposes an action.

Agent claims are not authority.

Agent interpretation cannot rewrite an authoritative execution outcome.

The Agent may misunderstand a result; that is an agency problem, not a reason to create another authority layer.

---

## Conductor

Conductor owns coordination:

    projects
    tasks
    assignment
    claiming
    dependencies
    lifecycle
    activity
    coordination visibility

Conductor answers:

    "What work exists?"
    "Who is doing it?"
    "What depends on what?"
    "What is happening in the coordination lifecycle?"

Conductor MUST NOT:

    authorize consequential actions
    decide whether an operation is allowed
    execute external effects
    become a policy engine
    interpret domain semantics as an authority decision
    become a Solvent replacement

Conductor MAY observe coordination-relevant facts such as:

    awaiting authorization
    awaiting execution result
    awaiting reconciliation

Those are coordination facts only.

They are NOT authority facts.

---

## Solvent

Solvent is the small trusted authority kernel.

Solvent answers:

    "Is this exact consequential operation authorized against the applicable authoritative state?"

Solvent owns:

    authorization
    authority state
    revocation
    exact operation binding

Solvent MUST NOT become:

    project manager
    task manager
    workflow runtime
    executor
    domain reasoning engine
    global effect registry

Keep Solvent small.

---

## Executor

Executor owns the external effect and execution reporting.

Executor answers:

    "What operation did we actually attempt?"
    "What happened?"

Where an external system of record exists:

    external SOR = authoritative for actual effect occurrence

Where no distinct SOR exists:

    Executor may itself be authoritative for actual effect occurrence.

Executor MUST NOT:

    create authority
    decide policy
    reinterpret Solvent authorization
    invent operation semantics

---

# IMPORTANT: AUTHORIZE ≠ EXECUTE

This distinction must remain visible in the implementation.

The reference flow is conceptually:

    Agent proposes X
        ↓
    Solvent authorizes exact X
        ↓
    Executor is permitted to attempt exact X
        ↓
    External system produces actual outcome

Do not collapse this into:

    approve()
        ↓
    execute()

as though those were one operation.

The point of the proof is to demonstrate that they remain distinct.

---

# WHAT WE ARE PROVING

The first implementation should answer whether the basic architecture is real.

We are NOT trying to build the final Loop Engineering platform.

We are NOT trying to solve every distributed-systems problem generically.

We are trying to prove the minimal reference loop.

The first proof should establish:

## A. Coordination works

Conductor can create/assign/claim/advance work without knowing domain-specific authority semantics.

## B. Authority works

Solvent can authorize an exact consequential operation.

## C. Exact binding works

The operation authorized is the operation presented to the Executor.

For example:

    authorized:
        transfer(A, B, 100)

    attempted:
        transfer(A, B, 100)

must succeed.

But:

    authorized:
        transfer(A, B, 100)

    attempted:
        transfer(A, B, 1000)

must not be treated as authorized.

## D. Execution works

Executor records what was attempted and what actually happened.

## E. Coordination observes the result

Conductor can observe that the work progressed through:

    authorization
    execution
    result

without becoming the authority owner.

## F. Agent can continue the loop

The Agent receives the result and can decide what work to propose next.

---

# FAILURE CASES ARE MORE IMPORTANT THAN THE HAPPY PATH

After the happy path works, deliberately attack it.

At minimum test:

    1. wrong operation
    2. missing authorization
    3. stale/invalid authorization
    4. authority verification unavailable
    5. execution timeout
    6. execution failure
    7. ambiguous external outcome
    8. duplicate execution/replay
    9. human intervention
    10. cancellation
    11. Agent misinterpreting the result
    12. declaration/version mismatch

The objective is not merely that errors occur.

The objective is that each error remains owned by the correct component.

Examples:

    Solvent denial
        ≠
    Executor failure

    Executor failure
        ≠
    ambiguous external effect

    Conductor cancellation
        ≠
    automatic Solvent revocation

    Agent interpretation
        ≠
    authoritative execution outcome

---

# DO NOT SOLVE FAILURE BY MERGING SYSTEMS

If an integration problem appears, do NOT immediately respond by moving the responsibility into Conductor or Solvent.

For example:

    "Conductor needs to know more about authorization"

does NOT imply:

    "Put authorization logic into Conductor."

Likewise:

    "Executor reconciliation is difficult"

does NOT imply:

    "Build a Loop reconciliation runtime."

And:

    "Cancellation and revocation sometimes interact"

does NOT imply:

    "Conductor cancellation automatically mutates Solvent authority."

First determine whether the problem is:

    implementation concern
    integration contract
    executor responsibility
    deployment responsibility
    genuinely new security property

Only the last category should trigger consideration of kernel/architecture changes.

---

# TOOL / CLIENT SUBSTITUTION IS PART OF THE PROOF

The architecture is supposed to be tool-agnostic.

Therefore, once the reference loop works, substitute participants.

At minimum:

    Agent A → Agent B

For example:

    GPT-driven client
    → deterministic scripted client

and:

    direct API/MCP interaction
    → another orchestration client such as Temporal

The semantic contract should not change.

We are trying to prove:

    implementation may change
    contract does not

---

# IMPORTANT SCOPE LIMIT

Do not turn the proof into a general-purpose workflow platform.

The desired implementation is deliberately small.

Do NOT introduce:

    workflow runtime
    workflow scheduler
    workflow event bus
    workflow database
    workflow state machine
    generic heartbeat service
    global effect registry
    mandatory gateway
    mandatory verification SDK

A gateway, SDK, polling mechanism, callback, sidecar, or other implementation mechanism may be used where a particular integration needs one.

But these are implementation choices, not new Loop Engineering roles.

---

# WHAT SUCCESS LOOKS LIKE

Success is NOT:

    "we built another workflow platform."

Success is:

    "an Agent can perform a complete consequential work loop through
     Conductor + Solvent + Executor, and the architectural boundaries
     remain intact under failure and substitution."

The strongest proof would be:

    GPT Agent
       ↓
    Conductor
       ↓
    Solvent
       ↓
    Executor A
       ↓
    External system

and then:

    scripted Agent
       ↓
    Conductor
       ↓
    Solvent
       ↓
    Executor B
       ↓
    same external system

with the same semantic outcome.

---

# IMPLEMENTATION PRIORITY

Prioritize in this order:

    1. smallest happy-path loop
    2. exact authorization/execution binding
    3. authoritative result handling
    4. Conductor coordination observation
    5. failure cases
    6. substitution tests

Do not spend implementation time building speculative abstractions before the reference loop exists.

When there is a choice between:

    more framework
    vs.
    more evidence

choose evidence.

When there is a choice between:

    broader abstraction
    vs.
    smaller working reference path

choose the smaller working reference path.

---

# ARCHITECTURAL NORTH STAR

Keep this visible throughout implementation:

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

with the invariant:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

and the architectural split:

    Conductor = coordination
    Solvent   = authority
    Executor  = effect
    Agent     = agency
    Domain    = meaning/truth

The purpose of this POC is to prove that these boundaries work together without requiring the systems to merge.