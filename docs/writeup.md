# From Architectural Decomposition to Empirical Proof: Why We Pivoted the Loop Engineering Effort

## Introduction

The Loop Engineering effort began with a deceptively simple architectural question:

> How can autonomous or semi-autonomous software perform consequential work while keeping agency, coordination, authority, and execution separate?

The initial architecture was intentionally small:

```text
Agent       = Agency
Conductor   = Coordination
Solvent     = Authority
Executor    = External Effect
Domain      = Meaning / Truth
```

The governing invariant was:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

The difficulty was not in naming these responsibilities. It emerged when we tried to make every boundary explicit enough that independent implementations could not accidentally collapse those responsibilities.

What followed was a series of adversarial review rounds that repeatedly uncovered subtle ambiguities: exact operation binding, authorization evidence, declaration lifecycle, external system-of-record semantics, long-running execution, ambiguous outcomes, human intervention, cancellation, substitution, conformance evidence, and the distinction between semantic ownership and enforcement responsibility.

That review process was valuable. It also exposed a second problem: the project could continue refining the specification indefinitely.

We therefore reached an important engineering decision:

> **Stop expanding the abstract architecture and prove it empirically with the smallest real end-to-end loop.**

This writeup explains how the complexity arose, what the review rounds taught us, why the pivot is necessary, and how the work now proceeds.

---

## 1. The Original Problem Was Not "Build a Workflow Engine"

The architecture was never intended to create another workflow product.

The initial separation was deliberate:

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
```

But the arrows do not mean that these components form one tightly coupled pipeline. They represent different responsibilities.

### Agent: agency

The Agent reasons, proposes, performs work, interprets results, and decides what to do next.

It does not become authoritative merely because it proposed something.

### Conductor: coordination

Conductor answers questions such as:

- What work exists?
- Who is responsible?
- What depends on what?
- What is the current coordination state?
- What happened in the project workflow?

It does not decide whether a consequential operation is authorized.

### Solvent: authority

Solvent answers a much narrower and more security-critical question:

> Is this exact consequential operation authorized against the applicable authoritative state?

Solvent is intended to remain a small trusted authority kernel.

### Executor: external effect

The Executor actually attempts the external operation and reports what happened.

It does not create authority.

### Domain: meaning and truth

The Domain determines what the result means and whether the result satisfies domain-specific requirements.

This distinction matters because software systems often collapse these roles into a single object or state transition. A common shortcut is effectively:

```text
task.status = approved
```

and then allowing "approved" to become synonymous with assigned, authorized, executed, and successful.

The architecture explicitly rejects that collapse.

---

# 2. Why the Review Process Became Complex

The complexity came from a simple fact:

> Once consequential software actions are involved, several different kinds of truth must remain separate at the same time.

A consequential action raises questions such as:

```text
What does the Agent want to do?
What work was assigned?
What exact operation was proposed?
Was that exact operation authorized?
What exact operation was attempted?
What actually happened externally?
Who knows what happened?
What does the result mean?
What should happen next?
```

Those are not the same question.

The review process therefore had to defend boundaries between:

```text
Agency
Coordination
Authority
Execution
External effect
Domain truth
Evidence
```

That produced a large number of apparently small but technically important issues.

---

# 3. The Review Rounds Revealed the Hidden Complexity

The Role / Boundary Matrix became the main focus of repeated adversarial review. Each round improved it, but each improvement exposed another edge case.

The pattern was instructive.

## 3.1 The ordinary-work bypass problem

One early version contained a broad "Ordinary Work" path to tools and external services.

The adversarial question was simple:

> Could an Agent invoke something through "ordinary work" that actually produces an external effect and therefore bypass authorization?

This forced an explicit distinction between:

```text
Non-effect capability
```

and:

```text
Effect-capable capability
```

The resulting rule became:

> Any operation capable of producing an external effect is consequential regardless of workflow phase or invoking client and must use the applicable consequential authorization boundary.

This was a key lesson:

> **Security boundaries cannot depend on the Agent merely labeling an action correctly.**

---

## 3.2 "Integration" accidentally became a seventh role

As the matrix became more precise, the word "integration" began appearing as though it were another participant:

```text
Agent
Conductor
Solvent
Executor
Human
Domain
Integration
```

That was an architectural smell.

The resolution was important:

> Integration is an implementation boundary or adapter belonging to an existing role boundary. It is not an independent Loop Engineering role.

An integration can be a wrapper, adapter, proxy, sidecar, gateway, library, or callback.

But it cannot silently acquire its own:

```text
coordination
authority
execution
```

ownership.

This prevented the specification from accidentally growing a new infrastructure layer simply because adapters exist in real implementations.

---

## 3.3 Ownership and enforcement were not the same thing

Another repeated problem was the temptation to put multiple owners into a single table cell:

```text
Solvent + integration
Executor + integration
Conductor + integration
```

That created ambiguity.

The solution was to distinguish:

```text
Semantic / authority owner
Effect / outcome owner
Enforcement / implementation responsibility
```

This produced a more durable rule:

> One semantic owner; possibly another enforcement mechanism.

For example:

```text
Solvent
  = authority owner

Executor-side adapter
  = enforcement mechanism
```

The enforcement component does not become the semantic authority simply because it checks the authorization.

---

## 3.4 Exact operation binding became a security problem

A major thread through the reviews was:

> How do we prove that the operation authorized is the operation executed?

It is not enough to carry a correlation ID.

For example:

```text
Authorized:
    transfer(A, B, 100)

Executed:
    transfer(A, B, 1000)
```

must not be accepted merely because both requests have the same task ID, request ID, or correlation ID.

The matrix therefore evolved toward an exact operation identity:

```text
proposal
  ↓
authorization
  ↓
execution
```

all using the same operation-identity definition.

Eventually the reviews uncovered another subtlety:

> Defining the identity fields is not enough. Independent implementations also need the same equality/comparison semantics.

Thus the declaration must define not only *what fields constitute identity*, but also *how equality is determined*.

This was one of the most important late-stage refinements.

---

# 4. Authorization Evidence Became Its Own Problem

The next question was:

> Even if Solvent authorizes something, how does the Executor know that the authorization is real?

A simple client-supplied claim such as:

```json
{"authorized": true}
```

obviously cannot be sufficient.

The contract therefore evolved toward:

```text
authorization evidence
+
exact operation identity
+
declaration / identity version
+
declared trust basis
```

The mechanism remained intentionally implementation-neutral.

The specification did not require one particular cryptographic scheme or canonicalization library.

Instead it required the security property:

> The evidence must be independently verifiable against the authority owner and must not rely solely on an unverified client claim.

That preserved tool and runtime neutrality without making the authority model vague.

---

# 5. Capability Declaration Became a Lifecycle

As the review progressed, another issue emerged:

> Who defines which operations are effect-capable?

If that information is missing, the boundary cannot be evaluated consistently.

This led to capability declarations containing concepts such as:

```text
owner
version
effective reference
effect classification
scope
operation identity
authorization boundary
trust basis
validity model
replay/idempotency behavior
outcome model
reconciliation mechanism
```

The declaration became versioned and subject to re-conformance when material semantics changed.

This also solved the non-effect → effect-capable transition:

```text
non-effect declaration
        ↓
effect surface changes
        ↓
new declaration version
        ↓
consequential coverage
        ↓
re-conformance
```

But again, this generated more questions: what if no distinct integration owner exists? The resolution was to make deployment ownership accountable for the deployed declaration without inventing another architectural role.

---

# 6. Long-Running Execution Exposed the Difference Between "Started" and "Finished"

A normal synchronous demo makes execution look easy:

```text
call()
→ response
→ done
```

Real systems are not always like that.

An operation can be:

- long-running,
- asynchronous,
- queued,
- interrupted,
- externally reconciled,
- or temporarily impossible to classify.

This forced us to define:

```text
not attempted
rejected before effect
attempted
succeeded
failed
terminated / revoked
ambiguous / unknown
```

The distinction between `failed` and `ambiguous` became especially important.

If we cannot establish whether the external effect occurred, we must not casually call it success or failure.

The review process also clarified that a queue acceptance or HTTP `202 Accepted` is not necessarily proof that the external effect occurred. Conversely, when the Executor itself is the authoritative system of record, its own evidence may be sufficient to establish success.

This is a real distributed-systems problem, not just terminology.

---

# 7. External Systems of Record Complicated the Result Model

A critical refinement was distinguishing:

```text
Executor
```

from:

```text
External System of Record
```

Where a separate system of record exists:

```text
Executor
  = attempted operation / execution reporting

External SOR
  = authoritative evidence of whether the effect actually occurred
```

Where there is no separate SOR:

```text
Executor
  = both execution reporter and effect authority
```

This distinction prevents a transport acknowledgment from being mistaken for actual effect.

It also forces the architecture to answer a deeper question:

> What does "succeeded" actually mean?

The final answer was:

> `succeeded` requires evidence sufficient under the declared execution/outcome model to establish successful effect occurrence.

That is significantly more robust than equating HTTP success with real-world success.

---

# 8. Human Intervention Turned Out to Be Another Boundary Problem

Human involvement introduced another potential architecture trap.

A human might:

- review a proposal,
- revise it,
- reject it,
- cancel progression,
- participate in authorization,
- interrupt execution,
- request reconciliation.

But none of that should automatically create a second authority path.

The rule became:

```text
Human participation
≠
second authority system
```

A material revision becomes a new proposal/version and must use the applicable authority boundary.

A human cancelling coordination does not automatically become an instruction to mutate Solvent authority.

Again, the complexity came from keeping semantics separated rather than allowing convenient implicit coupling.

---

# 9. The Conductor/Solvent Cancellation Problem Illustrated the Architectural Tension

One recurring adversarial argument was:

> If Conductor cancels a task, why shouldn't that automatically revoke the corresponding Solvent authorization?

This sounds intuitive.

But doing that creates a hidden dependency:

```text
Conductor lifecycle change
        ↓
Solvent authority mutation
```

That means Conductor has acquired an authority capability.

The architecture deliberately rejects that.

The correct model is:

```text
Conductor cancellation
    = coordination fact

Solvent revocation
    = authority operation
```

A deployment can explicitly perform a Solvent revocation if required by its policy.

But the coordination system does not gain implicit authority simply because the two systems are integrated.

This is one of the clearest examples of why the architecture is harder than a single unified workflow system.

---

# 10. The Result: We Were Getting Very Good at Finding Problems

By the later rounds, the review process was highly productive.

We were finding things like:

- operation identity comparison semantics,
- declaration version propagation,
- non-effect capability ownership,
- Executor-as-SOR semantics,
- autonomous UNKNOWN behavior,
- Conductor's coordination visibility,
- result routing,
- ambiguous termination,
- evidence requirements,
- declaration drift,
- conformance scope,
- and the difference between semantic ownership and enforcement responsibility.

This was valuable work.

But it also exposed the danger.

We could continue indefinitely.

Every new precision rule creates another edge where someone can ask:

> "What if this happens?"

Then another clause is added.

Then the clause introduces an interaction with another clause.

Then the next adversarial reviewer finds that interaction.

At some point:

```text
architecture review
→ specification
→ edge case
→ new rule
→ new edge case
→ new rule
→ ...
```

becomes an infinite loop.

That is the phenomenon we now call **analysis paralysis** in this effort.

---

# 11. Why the Pivot Became Necessary

The pivot is not an admission that the architecture was wrong.

It is the opposite.

The architecture is now sufficiently mature that we have reached diminishing returns from purely abstract refinement.

We therefore changed the development question.

### Before

> "Can we specify every important boundary before implementation?"

### Now

> "Can the existing boundaries survive a real implementation?"

This is a major change in method.

We stop trying to anticipate every failure.

Instead, we deliberately build the smallest working system and let actual behavior tell us which assumptions matter.

---

# 12. The Reference Loop Is the Experiment

The central experimental artifact is:

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

The purpose is not to create a new product.

It is to prove that the four main responsibilities can operate together without collapsing into one.

The first milestone is deliberately boring:

> one small consequential operation, one real or safely sandboxed external effect, one end-to-end path, and enough evidence to reconstruct what happened.

That gives us a concrete laboratory in which the architecture can be challenged.

---

# 13. Why This Curbs Analysis Paralysis

The pivot changes the feedback loop.

### Old feedback loop

```text
Question
  ↓
specification discussion
  ↓
adversarial review
  ↓
new clause
  ↓
another question
```

There is no natural stopping point.

### New feedback loop

```text
Hypothesis
  ↓
implement
  ↓
run
  ↓
observe
  ↓
failure
  ↓
classify
  ↓
fix or revise architecture
```

The most important new rule is:

> **Do not change architecture just because implementation is difficult.**

When something fails, classify it:

```text
implementation
integration
executor
deployment
new security property
```

Only the last category automatically raises the architecture question.

This prevents ordinary engineering friction from masquerading as an architectural defect.

---

# 14. The Reference Loop Also Changes the Meaning of the Specification

The specification is no longer the thing we are trying to finish.

It becomes a **testable contract**.

For example, the specification says:

```text
Conductor does not authorize.
```

The implementation now gets asked:

> Can we demonstrate that?

The specification says:

```text
Exact operation X must not become Y.
```

The implementation gets asked:

> Can we deliberately attempt Y and prove that it fails?

The specification says:

```text
Ambiguous ≠ success.
```

The implementation gets asked:

> Can we produce an ambiguous state and observe that classification?

This is a much healthier relationship between architecture and code.

---

# 15. What We Are Delivering in the Pivot

The pivot has a small, concrete set of outputs.

## 1. Reference Loop

A thin runnable demonstration using the existing:

```text
Agent
Conductor
Solvent
Executor
External system
```

No new workflow runtime.

## 2. Evidence Package

Enough participant-owned evidence to reconstruct:

```text
task
→ proposal
→ authorization
→ execution
→ external result
→ coordination update
→ agent continuation
```

The proof comes from the real components, not a fabricated workflow ledger.

## 3. Adversarial Scenarios

Deliberately break the loop:

```text
X authorized → Y executed
missing authorization
stale authorization
verification unavailable
version mismatch
duplicate execution
execution failure
ambiguous result
termination/revocation
cancellation
agent misinterpretation
```

## 4. Substitution Proof

Replace:

```text
Agent A → Agent B
Executor A → Executor B
orchestration/client mechanism A → B
```

and test whether the semantics remain unchanged.

## 5. Conformance Test Matrix

Once the implementation reveals what is actually observable, derive the conformance matrix from those real behaviors.

This ordering matters.

We are no longer writing a test matrix primarily from imagination.

We are deriving it from:

```text
architecture
+
working implementation
+
observed failures
```

---

# 16. What We Explicitly Do Not Build

The pivot also introduces a strong negative scope.

We are not building:

```text
Loop Engineering runtime
workflow engine
workflow database
workflow event bus
generic scheduler
generic heartbeat service
second authority engine
mandatory operations gateway
central effect registry
mandatory verification SDK
generic reconciliation daemon
merged Conductor/Solvent
```

A gateway, SDK, callback, polling mechanism, or sidecar may still be a sensible implementation choice for a particular integration.

But it does not become a new Loop Engineering role or mandatory platform component.

This is crucial.

---

# 17. Why Conductor and Solvent Stay Separate

The review process made the architectural argument stronger rather than weaker.

Solvent's value comes partly from being small.

It is supposed to be:

```text
trusted
narrow
authority-focused
formally reasoned
```

If coordination is merged into it, Solvent starts accumulating:

```text
tasks
assignments
dependencies
projects
workflow lifecycle
agents
activity
scheduling
```

The trusted authority kernel becomes a general workflow platform.

That increases:

- trusted computing base size,
- state-model complexity,
- attack surface,
- verification burden,
- coupling,
- and future feature pressure.

Conductor has the opposite property: it can evolve as a coordination system without becoming the authority kernel.

Keeping them separate therefore becomes more strategically valuable as the architecture matures.

---

# 18. The New Development Philosophy

The pivot can be summarized as five rules.

### Build the smallest thing that can falsify the architecture

Do not build a polished platform first.

### Prefer evidence over abstraction

A real trace is more valuable than another theoretical clause.

### Treat failure as information

A broken test tells us where our assumptions were wrong.

### Do not fix architectural problems with arbitrary new architecture

Only introduce a new architectural component when the failure reveals a genuinely new responsibility or security property.

### Stop when the proof is sufficient

Once the reference loop, adversarial cases, and substitution tests work, stop expanding it and move to generalization.

---

# 19. What Happens Next

The next stage is intentionally practical:

```text
1. Inspect existing Conductor/Solvent/Executor interfaces
2. Select the smallest safe consequential operation
3. Implement the happy path
4. Capture participant-owned evidence
5. Prove exact operation binding
6. Attack the loop
7. Classify failures
8. Fix implementation/integration/executor/deployment problems
9. Test substitution
10. Stop and report what reality revealed
```

Only afterward do we move toward:

```text
Conformance Test Matrix
      ↓
BM-IST validation
      ↓
Agent Skill
      ↓
Conformance Harness
      ↓
POC Runbook / demonstration
```

The implementation plan establishes this progression as:

> **working loop → evidence → failures → conformance → generalization**

---

# 20. The Deeper Lesson

The most important insight from the review rounds is methodological.

There is a point where architecture becomes more precise by being **tested** rather than further **described**.

We reached that point.

The repeated reviews were useful because they forced us to make implicit assumptions explicit. They helped us distinguish:

```text
coordination
from
authority

authority
from
execution

execution
from
truth

correlation
from
evidence
```

But continuing indefinitely would eventually produce a specification whose complexity exceeded the complexity of the problem it was intended to describe.

The Reference Loop is the mechanism for escaping that trap.

It changes the project from:

> **"Can we reason our way to a complete architecture?"**

to:

> **"Can reality prove or falsify the architecture we have?"**

That is the point where engineering becomes more valuable than additional abstraction.

---

## Conclusion

The pivot is not a retreat from rigor.

It is a change in where rigor happens.

Before the pivot, rigor was primarily:

```text
adversarial analysis
formal boundaries
specification
```

After the pivot, rigor becomes:

```text
real execution
participant-owned evidence
failure injection
substitution
conformance
```

The architectural north star remains unchanged:

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

with:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

The objective now is simple:

> **Prove that this separation works in reality before building anything more abstract around it.**

That is how the project moves forward without becoming trapped in its own analysis.
