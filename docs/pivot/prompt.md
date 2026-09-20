# Coding Agent Prompt — Implement Reference Loop v0.1

## Objective

Implement the **smallest real end-to-end Reference Loop** that empirically proves the existing architecture works:

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
External System
  ↓
Authoritative Result
  ↓
Conductor
  ↓
Agent
```

This is an implementation/proof exercise, **not an architecture-expansion exercise**.

The primary deliverable is a thin, runnable Reference Loop using the **existing Conductor, Solvent, and an existing or minimal Executor integration**. Do not create a new framework or workflow runtime. 

The objective is to prove:

1. the happy path works;
2. authority and execution remain separate;
3. exact operation binding survives the complete path;
4. externally authoritative results are preserved;
5. failures remain correctly classified;
6. the loop remains valid when participants are substituted.

---

# 1. Hard constraints

Do **not** introduce:

```text
Loop Engineering runtime
workflow engine
workflow database
workflow event bus
new scheduler
new gateway
merged Conductor/Solvent
second authority engine
generic reconciliation daemon
new Solvent kernel functionality
BM-IST infrastructure
```

These are explicitly outside this implementation phase. 

Do not turn the harness into a new architectural layer.

Do not copy or recreate Conductor or Solvent inside the reference-loop project. The reference-loop directories may contain only drivers, adapters, fixtures, scenarios, test code, evidence collection, and run instructions. 

Do not modify architecture merely because the first implementation gap appears.

Every discovered problem must first be classified as:

```text
implementation
integration
executor
deployment
new security property
```

Only the final category is grounds for architectural reconsideration. 

---

# 2. Phase 0 — Repository reconnaissance

Before changing code, inspect the existing implementations.

Determine concretely:

### Conductor

* How a task is created.
* How an agent discovers a task.
* How an agent claims a task.
* How activity/progress is recorded.
* How task completion/submission works.
* Which API/MCP/client surface should be used.

### Solvent

* How a consequential operation is represented.
* How authorization is requested.
* What constitutes authorization evidence.
* How declaration/version information is represented.
* How exact operation identity is compared.
* How unavailable verification is represented.
* How stale/invalid authorization is rejected.

### Executor

* What existing consequential capabilities already exist.
* What operation they perform.
* How authorization evidence reaches the executor.
* How operation identity is checked.
* How execution outcomes are represented.
* Whether there is already an external system/SOR.

### External result

Determine explicitly whether:

```text
Executor == authoritative SOR
```

or:

```text
Executor != authoritative SOR
Executor reports execution
external SOR establishes actual effect occurrence
```

Do not assume that an HTTP success, queue acknowledgement, or executor acknowledgement is equivalent to the actual external effect.

The authoritative result source must be explicitly recorded in the implementation/reconnaissance report. This distinction is required by the review. 

### Important selection rule

**Prefer an operation already supported by the existing Executor.**

Do not create a new Executor capability solely to make the demo interesting. Only create the minimum integration required if reconnaissance proves that no suitable existing safe consequential operation is available. 

---

# 3. Select the reference operation

Choose exactly **one** small consequential operation that is:

* safe;
* observable;
* reversible/resettable where practical;
* easy to execute repeatedly in a test environment;
* already supported by the existing Executor if possible.

Represent the operation precisely:

```text
operation identity
target
parameters
declaration/version
```

The exact same operation identity must survive:

```text
proposal
→ authorization
→ execution
```

Do not introduce a second interpretation of the operation in the harness.

---

# 4. Implement the happy path

Implement the smallest path that produces this sequence:

```text
TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED(X)
→ AUTHORIZED(X)
→ EXECUTION_ATTEMPTED(X)
→ EFFECT_CONFIRMED(X)
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED
```

This is the central Reference Loop deliverable. 

The operation must actually cross the real authority/effect boundary.

Do not simulate Solvent authorization and do not fake Executor execution merely to create a successful trace.

---

# 5. Evidence ownership is mandatory

The reference loop must produce enough evidence to reconstruct:

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

This is an **evidence package**, not a new workflow ledger. 

Distinguish authoritative participant evidence from harness-level correlation.

For example:

```text
TASK_CREATED
    [Conductor evidence]

TASK_CLAIMED
    [Conductor evidence]

PROPOSAL_CREATED(X)
    [Agent evidence]

AUTHORIZED(X)
    [Solvent evidence]

EXECUTION_ATTEMPTED(X)
    [Executor evidence]

EFFECT_CONFIRMED(X)
    [SOR/Executor authoritative evidence]

RESULT_OBSERVED
    [Agent/Conductor evidence]

CONDUCTOR_UPDATED
    [Conductor evidence]

AGENT_CONTINUED
    [Agent evidence]
```

The harness may correlate these facts.

**The harness must not manufacture authoritative facts.**

This distinction is explicitly required by the review. 

For every evidence item, make it possible to determine:

```text
owner
source
operation identity
declaration/version where applicable
timestamp or ordering information where available
correlation identifier
```

Do not create a centralized authority ledger.

---

# 6. Create the minimal reference-loop harness

Use a structure approximately like:

```text
reference-loop/
    README.md
    scenario/
        happy_path
        wrong_operation
        denied
        unavailable
        declaration_version_mismatch
        replay
        execution_failure
        ambiguous
        termination_revocation
        cancellation
        misinterpretation
    agent/
    conductor/
    solvent/
    executor/
    evidence/
    runbook/
```

Adapt the structure to the actual repository.

Do not mechanically create empty directories if the repository does not need them.

Again, these are test/integration artifacts, **not duplicate implementations of the infrastructure components**.

---

# 7. Provide one runnable entry point

Create a simple command equivalent to:

```text
run-reference-loop
```

It should run the happy path end-to-end and produce inspectable evidence.

It should not become a new workflow runtime.

The output should make the important transitions visible, for example:

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

However, the displayed trace is only a **correlation/view** over participant-owned evidence. It is not itself authoritative.

---

# 8. Verify exact binding

This is one of the most important tests.

Establish:

```text
Authorized(X)
Execute(Y)
```

where `Y != X`.

Expected result:

```text
FAIL CLOSED
NO UNAUTHORIZED EFFECT
```

Prove from evidence that:

```text
authorization identity != execution identity
```

and that the execution boundary rejected the mismatch.

Do not merely test that the calling client behaved correctly.

---

# 9. Adversarial Reference Scenarios

After the happy path works, implement these scenarios:

### A. Wrong operation

```text
authorized X
→ execute Y
```

Must fail.

### B. Missing authorization

```text
execute X
without valid authorization
```

Must fail closed.

### C. Stale authorization

Use expired or otherwise stale authorization.

Must fail according to the declared validity model.

### D. Verification unavailable

Make required authorization verification unavailable.

Must fail closed unless the existing declared semantics provide independently verifiable valid evidence.

### E. Declaration/version mismatch

Explicitly test:

```text
authorize under declaration/identity version V1
→ execute under incompatible V2
```

Expected:

```text
FAIL CLOSED
```

This test is mandatory; it directly exercises the previously established declaration/version boundary. 

### F. Duplicate/replay

Submit the same authorized consequential operation more than once.

Verify the existing replay/idempotency semantics.

Do not invent a new replay service.

### G. Execution failure

Force the Executor to fail.

Result must remain:

```text
FAILURE
```

It must not be converted into success by Conductor, Agent, or harness logic.

### H. Ambiguous outcome

Create a case where the final external outcome cannot immediately be known.

Result must remain:

```text
AMBIGUOUS
```

Do not collapse ambiguity into success or failure.

### I. Termination/revocation

Test the existing semantics for termination/revocation.

Distinguish:

```text
known terminated/revoked
```

from:

```text
ambiguous termination
```

Do not infer zero external effect merely because an execution was terminated.

### J. Conductor cancellation

Cancel or stop the coordinated task through Conductor.

Verify that this does **not silently become Solvent authority revocation**.

### K. Agent misinterpretation

Give the Agent a deliberately incorrect interpretation of a returned result.

Verify that the authoritative result itself remains unchanged.

These scenarios are the concrete adversarial proof phase. 

---

# 10. Evidence for every adversarial scenario

Every scenario must record:

```text
scenario name
setup
operation identity
authorization state
execution attempt
external outcome
expected outcome
actual outcome
authoritative evidence
classification
```

For failures, use exactly:

```text
implementation
integration
executor
deployment
new security property
```

Do not conceal failures.

The purpose is to discover where the architecture is already sufficient and where it is genuinely insufficient. 

---

# 11. Substitution proof

Only begin substitution **after the reference loop and adversarial tests are stable**.

Do this in order.

## A. Agent substitution

Replace:

```text
Agent A
```

with:

```text
Agent B
```

while keeping:

```text
Conductor
Solvent
Executor
External System
```

unchanged.

The semantic contract must remain unchanged.

## B. Executor substitution

Replace:

```text
Executor A
```

with:

```text
Executor B
```

while preserving the authority semantics.

## C. Orchestration/client substitution

Replace the mechanism used to drive the interaction with another mechanism.

The requirement is not:

> "it still somehow works."

The requirement is:

> **The semantic contract remains unchanged.**

This is the actual evidence for tool/framework agnosticism. 

---

# 12. Do not implement the Conformance Test Matrix yet

Do **not** make the Conformance Test Matrix the first deliverable.

First produce:

```text
1. Reference Loop
2. Evidence Package
3. Adversarial Scenarios
4. Substitution Proof
```

Only then derive:

```text
5. Conformance Test Matrix
```

The matrix must be derived from:

```text
observed implementation
+
observed failure modes
+
locked boundary contract
```

rather than from additional theoretical speculation. 

---

# 13. Explicit non-expansion rule

During implementation, when something fails, **do not immediately add abstractions**.

Ask:

```text
Is this an implementation defect?
Is this an integration defect?
Is this an Executor limitation?
Is this a deployment/environment problem?
Or have we discovered a genuinely new security property?
```

If it is one of the first four:

```text
fix the implementation
```

If it is the fifth:

```text
STOP
document the evidence
do not redesign silently
```

Return the finding for architectural review.

---

# 14. Stopping gate

Stop expanding the implementation once all five conditions are demonstrated:

### Gate 1 — Happy path

The real loop works:

```text
Agent
→ Conductor
→ Solvent
→ Executor
→ External System
→ Result
→ Conductor
→ Agent
```

### Gate 2 — Boundary integrity

Unauthorized, mismatched, stale, unavailable, or incompatible operations do not cross the effect boundary.

### Gate 3 — Outcome integrity

Success, failure, ambiguity, and termination remain correctly distinguished.

### Gate 4 — Ownership integrity

Each authoritative fact remains owned by the correct participant.

### Gate 5 — Substitution

Agent substitution plus Executor/orchestration substitution preserve the semantic contract.

At this point:

**STOP.**

Do not add features merely because the demo could be made more elaborate.

The purpose is empirical architectural proof, not an increasingly sophisticated demo. 

---

# 15. Final report required from the coding agent

When finished, return a concise implementation report containing:

```text
1. Reference operation selected
2. Why it was selected
3. Existing components reused
4. Files changed
5. Happy-path result
6. Evidence produced
7. Authoritative result source
8. Adversarial scenario results
9. Substitution results
10. Failures discovered
11. Classification of each failure
12. Any genuinely new security property discovered
13. Architectural changes required, if any
14. What was deliberately NOT changed
15. Exact command(s) used to reproduce the proof
```

Also provide a compact final trace showing the actual participant-owned evidence corresponding to:

```text
TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED
→ AUTHORIZED
→ EXECUTION_ATTEMPTED
→ EFFECT_CONFIRMED
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED
```

The final report must clearly distinguish:

```text
implemented
verified
observed
assumed
not yet proven
```

Do not claim architectural validation for anything that was not actually demonstrated.

---

## Success criterion

The implementation is successful when we can honestly say:

> **The existing Agent → Conductor → Solvent → Executor → external system → Result → Conductor → Agent loop works in reality, its authority/execution boundaries survive adversarial cases, its authoritative evidence can be reconstructed, and its semantic contract survives substitution.**

That is the entire purpose of this phase. The five deliverables are intentionally ordered **Reference Loop → Evidence → Adversarial Scenarios → Substitution → Conformance**, with BM-IST and the later Agent Skill/Harness/Runbook work deliberately downstream. 

**Do not expand scope beyond this prompt without stopping and reporting the reason first.**
