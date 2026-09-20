# From Architectural Decomposition to Empirical Proof
## Technical Writeup: What the Reference Loop Has Taught Us, What Remains Hard, and Where the Architecture Goes Next

**Status:** Phase 1 Reference Loop complete  
**Current verdict:** PASS  
**Frozen Solvent baseline:** `7602699`  
**Scope:** Agent, Conductor, Solvent, Executor, and the Reference Loop protocol boundary

---

## Executive Summary

The most important result so far is not that a happy-path demo works. It is that the architecture has survived contact with an actual end-to-end implementation without requiring the roles to collapse into one another.

The Reference Loop has now demonstrated, with a frozen Solvent implementation, that a consequential operation can move through distinct agency, coordination, authority, and execution boundaries:

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
result / evidence
  ↓
Conductor
  ↓
Agent
```

The core invariant remains intact:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

The experiments have also empirically confirmed:

```text
AUTHORIZE ≠ EXECUTE
```

A mismatched actor is denied. Execution without a valid authorization is denied. A valid authorization with a valid `intent_id` reaches the executor exactly once in the controlled execution experiment.

Just as important, the experiments have **not** yet justified adding a new cross-role primitive, a unified workflow runtime, a second authority layer, or a large protocol framework.

That is a major architectural finding.

The project has therefore moved from a period of conceptual decomposition toward a more disciplined empirical program:

```text
architecture
    ↓
minimal runnable loop
    ↓
boundary tests
    ↓
evidence
    ↓
formalize only what the evidence proves necessary
```

The next stage is not "build everything." It is to close the remaining empirical gaps—especially real external-effect confirmation—and then formalize the minimal protocols and skills that the observed system actually requires.

---

# 1. The Central Insight: Separation Is Testable

The original architecture separated four concerns:

| Role | Responsibility |
|---|---|
| **Agent** | Agency: reasoning, decomposition, planning, replanning, work creation, execution of ordinary work, reporting |
| **Conductor** | Coordination: work items, dependencies, lifecycle, claim/release, activity, coordination context |
| **Solvent** | Authority: determine whether an exact consequential action is authorized against an exact state |
| **Executor** | Effect: perform the authorized external operation and return execution/result evidence |

The temptation in systems like this is to collapse responsibilities because a single component can technically perform multiple functions.

For example:

```text
Agent + planner + task manager + policy engine + executor
```

may look simpler initially.

The Reference Loop is demonstrating why that simplification is dangerous.

The system can remain useful without requiring each layer to understand every other layer's semantics.

The resulting boundary is:

```text
Agent:
    What work should exist?
    How should complex work be decomposed?
    What should happen next based on results?

Conductor:
    What work exists?
    Who is doing it?
    What depends on what?
    What is the current lifecycle state?

Solvent:
    Is this exact consequential action authorized?

Executor:
    Perform this already-authorized operation.
```

This is not merely a conceptual taxonomy anymore. The implementation has begun to exercise those distinctions directly.

---

# 2. Task Decomposition Belongs to Agency, Not Conductor

One of the most useful clarifications was resolving the relationship between a complex task and Conductor.

Suppose a user gives an Agent:

> Implement GitHub OAuth.

The Agent—not Conductor—should reason about the technical structure:

```text
Task A: Implement GitHub OAuth

    B: configuration
    C: authorization endpoint
    D: callback handling
    E: token persistence
    F: authentication middleware
    G: tests
```

The Agent can discover and revise this decomposition as repository understanding improves.

Conductor does not need an LLM planner.

Instead, Conductor stores and coordinates whatever work the Agent has decided should exist:

```text
B ────────┐
C ────────┤
D ────────┼──→ G
E ────────┘
```

Conductor answers:

- What work exists?
- Who owns it?
- What depends on what?
- What is its state?
- What happened?

The Agent answers:

- What work should exist?
- How should it be decomposed?
- What should be done next?

This preserves a crucial distinction:

```text
DECOMPOSITION ≠ COORDINATION
```

It also allows progressive decomposition: a supposedly simple task B can later become B1/B2/B3 without turning Conductor into a hierarchical planning engine.

---

# 3. The Reference Loop Changed the Development Method

A major strategic insight was that the project was approaching diminishing returns from architecture-first reasoning.

There is a point at which additional role diagrams, specifications, invariants, and hypothetical edge cases stop increasing confidence.

The answer was to construct the smallest real loop that could falsify the architecture.

That became the Reference Loop:

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
external effect / result
  ↓
Conductor
  ↓
Agent
```

The goal was not to build a new runtime or orchestration system.

The goal was to answer:

> Can the existing responsibilities cooperate on a consequential operation without collapsing their boundaries?

That shift matters because evidence now has priority over architectural speculation.

The Reference Loop is therefore an empirical instrument, not another infrastructure product.

---

# 4. The Frozen Solvent Rule Became Part of the Experimental Method

One of the most important operational discoveries was that the Reference Loop must adapt to Solvent—not modify Solvent to fit the Reference Loop.

The approved frozen boundary was deliberately established at:

```text
7602699  "✨ solvent kernel freeze"
```

The current repository initially contained later commits. The forensic process identified that the active branch was ahead of the freeze and that source files differed from the baseline. The repository was then deliberately restored to the approved frozen commit.

The final invariant became:

```text
Solvent HEAD = 7602699
working tree = clean
Solvent source modifications = none
```

This matters scientifically as much as it matters operationally.

The Reference Loop is now testing against a stable authority implementation.

That makes claims about the integration more meaningful because the integration experiment is not silently changing the authority system it is supposed to evaluate.

The rule is now:

```text
Reference Loop adapts to frozen Solvent.
Frozen Solvent does not adapt to Reference Loop.
```

---

# 5. Persistence Boundaries Are Allowed to Differ

Another useful lesson was that infrastructure uniformity is not the same thing as architectural consistency.

Conductor uses:

```text
SQLite
```

Solvent uses:

```text
pgx + CockroachDB
```

These are separate persistence boundaries.

Trying to force both into SQLite would have made the integration experiment easier to set up but less faithful to the actual system architecture.

The better model is:

```text
Conductor → its own persistence
Solvent   → its own persistence
```

The Reference Loop connects participants through their interfaces rather than flattening their implementation differences.

That is an important principle for future protocol work as well:

> Interoperability should occur at the contract boundary, not by erasing implementation differences.

---

# 6. The Agent–Conductor Boundary Is Becoming Concrete

The Reference Loop observed the intended Agent → Conductor interaction through MCP.

The Agent:

- discovers work,
- claims work,
- updates work,
- submits work,
- continues based on returned state.

The Agent does not directly manipulate Conductor storage.

This is significant because it provides the beginning of what will later become an explicit Agent Skill or work protocol.

The minimum observed future Agent Skill is already visible:

```text
discover work
claim work
report progress
submit
report blockers
continue / replan
```

But the experiment does **not** yet justify a complete Agent Skill specification.

That is deliberately deferred.

---

# 7. Authorization Is a Different Kind of State Than Work

The most important security insight is that authorization is not merely another workflow state.

Conductor can say:

```text
Task = active
```

but that does not mean:

```text
Action X is authorized
```

Likewise, an Agent can have a capability to perform an operation, but that does not create authority for a particular target and state.

The Reference Loop has now directly exercised this distinction.

Solvent is the authority boundary.

The critical rule remains:

```text
AUTHORIZE ≠ EXECUTE
```

Authorization establishes that a precise operation is eligible to execute.

The Executor performs the operation.

That distinction prevents several failure modes:

- a task being mistaken for permission;
- an executor becoming an authority engine;
- a Conductor becoming a policy engine;
- a stale cached authorization becoming equivalent to current authority.

---

# 8. Exact Operation Identity Is More Important Than Generic Permission

The Reference Loop has also reinforced that authorization must be bound to the exact operation.

The canonical operation identity used in the experiment is:

```text
deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>
```

The important property is not the string itself.

The important property is that the same intended operation survives across:

```text
proposal
authorization
intent
execution request
executor invocation
result/evidence correlation
```

The experiment uses distinct identifiers for different concepts:

```text
run_id
scenario_id
task_id
intent_id
operation_id
```

Those identifiers must not be collapsed merely because they can technically be represented by one value.

The practical invariant is:

```text
No silent operation drift.
```

If authorization is for X and execution is attempted for Y, the system must not accidentally treat the authorization as valid.

---

# 9. The Negative-Path Tests Are as Important as the Happy Path

The happy path by itself proves very little.

The first two focused negative tests added real value.

## Wrong actor

Authorization for actor A was followed by an authorization attempt by actor B.

The real Solvent REST authorization boundary rejected it.

```text
authorized actor ≠ attempted actor
            ↓
          DENY
```

The executor was never invoked.

## No authorization

Execution was attempted with no valid authorization intent.

The real Solvent execution boundary rejected it.

```text
no valid intent
      ↓
    DENY
```

Again, the executor was never invoked.

These tests demonstrate fail-closed behavior at the authority boundary rather than simply testing that the happy path works.

---

# 10. Exactly-Once Invocation Is a Separate Claim From External Effect

The positive execution experiment added an important refinement.

We instrumented the fake/RecordingFunc executor with a thread-safe call counter and proved:

```text
valid authorization
      ↓
valid intent_id
      ↓
ExecuteAction
      ↓
executor invoked exactly once
```

The test also verified:

- one `AUTHORIZED` event;
- one execution attempt;
- one executor invocation;
- consistent operation identity;
- simulation/non-external-effect labeling.

This establishes a much stronger execution-path invariant than "the request returned success."

However:

```text
exactly-one executor invocation
        ≠
real external effect
```

The test uses a fake/RecordingFunc executor.

Therefore the system has **not yet proven** that an external service actually changed state exactly once.

That distinction must remain explicit.

---

# 11. Evidence Ownership Is Emerging as a Fundamental Principle

The experiments have clarified a useful evidence model.

### Conductor owns

- work lifecycle
- activity
- coordination state

### Solvent owns

- authority
- intent
- authority/audit records

### Executor owns

- execution attempt evidence within the declared execution model

### External SOR owns

- actual external system state, when a real external operation is performed

### Harness owns

- post-run aggregation and correlation evidence

The harness must not manufacture authoritative operational state.

This yields a strong principle:

> Evidence must remain attached to the system that actually knows the fact.

That becomes especially important when real external effects are introduced.

---

# 12. What the Reference Loop Has Not Proven

The current evidence is strong but deliberately incomplete.

The Reference Loop has **not yet demonstrated**:

### Real external effect

The executor is still simulated.

There is no verified real GitHub workflow result bound to the operation in the current Phase 1 evidence.

### Concurrency

There has been no meaningful multi-agent contention experiment.

### Cancellation

There is no validated cancellation protocol.

### Replay / duplicate handling

Exactly-once invocation was tested in one controlled positive path, but broader replay/idempotency semantics remain unproven.

### Revocation

The project has not yet established how long-running execution interacts with revoked authority.

### Failure recovery

Executor failures and partial failures remain future work.

### Long-running authorization

The current tests do not establish behavior across longer-lived authority or state transitions.

### Substitution

We have not yet demonstrated that replacing an Agent implementation, Conductor client, Executor, or other participant preserves the same protocol semantics.

### Full protocol formalization

The work has identified protocol boundaries, but the final protocols have not yet been formalized.

---

# 13. The Agent, Conductor, Solvent, and Executor Skills Should Be Derived From Evidence

A major design decision now becomes clearer.

There should eventually be explicit participant skills/contracts:

```text
Agent Skill
Conductor Skill / operational contract
Solvent Authority Interface
Executor Capability Contract
```

But these should be formalized after the Reference Loop exposes the minimum required behavior.

The Agent Skill should describe how an autonomous participant:

```text
receive work
→ reason
→ decompose
→ create/claim work
→ perform work
→ report
→ replan
```

The Conductor contract should describe how an Agent interacts with coordination:

```text
discover
→ claim
→ update
→ report blockers
→ submit
→ continue
```

Solvent's contract should describe:

```text
exact action
+
exact target/state
→
authorization decision
```

Executor's contract should describe:

```text
authorized operation
→
execution
→
execution/result evidence
```

The key is that these should be **boundary contracts**, not one giant protocol runtime.

---

# 14. The Protocol Should Not Be One Monolith

The emerging protocol model is better represented as several contracts:

```text
Agent ↔ Conductor
    Work Protocol

Conductor ↔ Solvent
    Consequential-Action / Authorization Protocol

Solvent ↔ Executor
    Authorization Handoff / Execution Protocol

Executor ↔ External SOR
    Effect Protocol
```

Across all of them there is a shared semantic layer:

```text
operation identity
correlation
authorization evidence
result/evidence semantics
```

This is a cleaner model than creating a single "mega-protocol" that understands every participant.

---

# 15. The Five-Primitive Workflow Idea Is Still Only a Hypothesis

A possible future abstraction has been discussed around:

```text
Intent
Work
Authorization
Effect
Outcome
```

The Phase 1 evidence does **not** justify turning those into a new runtime primitive.

This is important.

The experiments so far have been expressible using existing role boundaries.

Therefore:

```text
New primitive justified = NO
New security property discovered = NO
```

That is not a lack of progress.

It is evidence that the current architecture is sufficient for the behaviors tested.

A new cross-role abstraction should only appear if the next experiments demonstrate a durable problem that cannot be expressed cleanly with the current contracts.

---

# 16. The Hardest Challenges Ahead

The next challenges are not primarily about adding components.

They are about proving semantics under conditions where systems normally become ambiguous.

## 16.1 Real external effects

The largest immediate gap is:

```text
Executor invocation
        ↓
actual external state transition
        ↓
authoritative confirmation
```

This requires a real external source of record.

The system must be able to distinguish:

```text
attempted
succeeded
failed
ambiguous
```

without turning a local assumption into authoritative effect evidence.

The external system may itself have retries, delayed state, duplicated requests, or partial visibility.

This is where the architecture's evidence model gets seriously tested.

---

## 16.2 Ambiguous outcomes

A particularly difficult case is:

```text
request sent
    ↓
network failure
    ↓
unknown whether effect occurred
```

The correct state may be:

```text
AMBIGUOUS
```

not:

```text
FAILED
```

and definitely not:

```text
SUCCEEDED
```

unless the declared source of record proves success.

This is likely to be one of the most important future tests.

---

## 16.3 Revocation and long-running operations

What happens when:

```text
authorization granted
        ↓
operation starts
        ↓
authority is revoked
        ↓
execution continues
```

The system needs a defined precedence rule.

This is not necessarily a new role.

It is a protocol semantics question that must be answered empirically and then formalized.

---

## 16.4 Concurrency and contention

Multiple Agents may eventually interact with the same work.

Questions include:

- Who gets to claim a task?
- Can two agents act on the same task?
- What happens when one Agent is slow?
- What happens when a task changes while an Agent is acting?
- What happens when two valid operations race?

These are Conductor and coordination questions, not reasons to move authorization into Conductor.

---

## 16.5 Replay and idempotency

Exactly-once invocation in a controlled test does not prove exactly-once behavior across distributed failures.

Future experiments will need to distinguish:

```text
duplicate request
duplicate execution
duplicate effect
duplicate observation
duplicate evidence
```

These are different phenomena.

---

## 16.6 Substitution

A healthy protocol should survive implementation substitution.

For example:

```text
Agent A → Conductor → Solvent → Executor X
Agent B → Conductor → Solvent → Executor X
Agent A → Conductor → Solvent → Executor Y
```

provided the substitutions preserve the required contracts.

This is one of the reasons protocol formalization should be driven by observed semantics rather than internal implementation details.

---

# 17. Roadmap

The roadmap should remain evidence-driven.

## Phase 1 — Reference Loop

**Status: complete for the current controlled experiment.**

Proven:

```text
happy path
wrong actor denied
missing authorization denied
valid authorization → exactly-one execution
operation identity preserved
evidence boundaries preserved
Solvent remains frozen
```

Still limited by:

```text
fake executor
no real external effect
no concurrency/cancellation/replay/revocation
```

---

## Phase 2 — Real External-Effect Validation

The most natural next experiment is to exercise a real consequential operation against a real external source of record.

The goal is to prove:

```text
authorization
→ execution
→ actual external effect
→ authoritative result observation
→ correct classification
```

This phase should establish the difference between:

```text
EXECUTION_ATTEMPTED
EFFECT_CONFIRMED
RESULT_OBSERVED
AMBIGUOUS
```

without collapsing those states.

The external SOR should be treated as authoritative for the portion of reality it actually controls.

---

## Phase 3 — Focused Adversarial Scenarios

Only after a real external effect is demonstrated should the broader failure matrix be expanded.

Potential scenarios:

```text
wrong operation
missing authorization
stale authorization
verification unavailable
declaration/version mismatch
replay/duplicate
execution failure
ambiguous outcome
termination/revocation
Conductor cancellation
Agent misinterpretation
cross-implementation operation identity mismatch
```

The purpose is not to produce an enormous test suite.

Each scenario should test a specific architectural invariant.

---

## Phase 4 — Protocol Formalization

Now formalize what the experiments have actually shown to be necessary.

Produce explicit boundary contracts for:

```text
Agent ↔ Conductor
Conductor ↔ Solvent
Solvent ↔ Executor
Executor ↔ External SOR
```

Also formalize shared semantics for:

```text
operation identity
correlation
authorization evidence
result/evidence
failure classification
```

This is where the previously identified Agent Skill and Conductor Skill become concrete specifications.

The Agent Skill should focus on agency and work behavior.

The Conductor contract should focus on coordination semantics.

Solvent remains authority.

Executor remains effect.

---

## Phase 5 — Conformance and Substitution

Build a Conformance Matrix from observed behavior.

Then test:

```text
participant substitution
implementation substitution
client/tool substitution
executor substitution
```

The goal is semantic interoperability, not implementation uniformity.

A conforming substitute should preserve the contract even if its internal architecture differs.

---

## Phase 6 — Domain Validation

Only after the generic loop and contracts are stable should domain-specific workloads such as BM-IST become a serious validation target.

BM-IST should be treated as a workload that tests whether the generic architecture is useful—not as a reason to contaminate Conductor or Solvent with domain-specific semantics.

The domain belongs at the workload/extension boundary.

---

# 18. What Should Not Be Built

Several tempting directions should remain explicitly out of scope unless evidence forces them.

Do not prematurely build:

- a universal workflow engine;
- an event bus;
- a second authority engine;
- a domain ontology inside Conductor;
- a planner inside Conductor;
- a policy engine inside Conductor;
- a Solvent replacement;
- a generic execution broker;
- a new cross-role primitive simply because it is elegant;
- a giant Agent/Conductor/Solvent/Executor protocol before real failure modes have been observed.

The guiding principle is:

> Grow the ecosystem, not the kernel.

Every new abstraction should earn its place through a demonstrated requirement.

---

# 19. The Most Important Architectural Lessons

### 19.1 Agency and coordination are orthogonal

The Agent reasons about work.

Conductor coordinates the resulting work.

Neither responsibility needs to absorb the other.

### 19.2 Authority should be an explicit boundary

Authorization is not a task status.

It is not a capability.

It is not an executor result.

It is a distinct security fact.

### 19.3 Execution is not authority

The Executor should never be able to create authority merely because it can perform an action.

### 19.4 Evidence has an owner

A system should not declare reality merely because its own process returned success.

The authoritative source of a fact should be explicit.

### 19.5 Exact identity is stronger than generic permission

"Can this agent deploy?" is different from:

> "Is this exact deployment operation authorized against this exact target and state?"

The latter is what consequential execution needs.

### 19.6 Frozen components improve experimental quality

A frozen Solvent made it possible to ask:

> Does the integration work against an authority implementation without changing the authority implementation?

That is a much stronger experiment than modifying the authority system until the demo works.

### 19.7 Evidence beats architectural speculation

The strongest progress so far came from small experiments:

```text
happy path
→ wrong actor
→ missing authorization
→ exactly-one execution
```

Each experiment eliminated a class of uncertainty.

---

# 20. Current State

The current empirical state can be summarized as:

```text
                    ┌─────────────┐
                    │    Agent    │
                    │   Agency    │
                    └──────┬──────┘
                           │
                     Work Protocol
                           │
                           ▼
                    ┌─────────────┐
                    │  Conductor  │
                    │ Coordination│
                    └──────┬──────┘
                           │
             Authorization Protocol
                           │
                           ▼
                    ┌─────────────┐
                    │   Solvent   │
                    │  Authority  │
                    └──────┬──────┘
                           │
              Execution Protocol
                           │
                           ▼
                    ┌─────────────┐
                    │  Executor   │
                    │    Effect   │
                    └──────┬──────┘
                           │
                    External SOR
```

The current evidence supports the model:

```text
Agent        = agency
Conductor    = coordination
Solvent      = authority
Executor     = execution
External SOR = authoritative external state
```

The next major unknown is:

```text
Can the same boundaries remain correct when the system produces
a real external effect and that effect is observed through an
authoritative external source of record?
```

That is the next experiment worth running.

---

# 21. Final Perspective

The deepest lesson so far is not that the architecture is finished.

It is that the architecture has begun to become **falsifiable**.

That is a much better position than having a beautiful diagram that has never been forced through a real execution path.

The Reference Loop has shown that:

```text
agency
coordination
authority
execution
```

can remain separate while still forming a useful system.

It has also shown that the project does not yet need a larger architecture simply because the system is conceptually rich.

The next stage should therefore continue the same discipline:

```text
prove
    ↓
observe
    ↓
classify
    ↓
formalize
    ↓
conform
```

not:

```text
imagine
    ↓
abstract
    ↓
generalize
    ↓
build
    ↓
hope
```

The architecture should earn its next layer.

And the protocol should be something discovered and formalized from working boundaries—not something invented in the abstract and then forced onto the implementation.
