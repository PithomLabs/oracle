I reviewed the attached **Workflow Specification v0.3** against the design we have now converged on. My conclusion is:

> **Do not modify the implementation yet. Do not freeze v0.3 as the new workflow design. Use v0.3 as the parent normative lineage, then create the much smaller Workflow Design v1.0 that supersedes the obsolete workflow-mode material while preserving the security semantics that remain valid.**

The important point is that v0.3 is **not wrong**. It contains many of the hard-won security and conformance semantics we need, especially exact operation binding, fail-closed execution, authorization taxonomy, replay/idempotency, ambiguity, evidence ownership, declaration lifecycle, and Growth Gate.   

But its **workflow model is now too large for the design we have chosen**.

## What v0.3 should retain

These remain foundational:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
AUTHORIZATION ≠ EXECUTION
CORRELATION ≠ AUTHORITY
TASK COMPLETION ≠ DOMAIN TRUTH
HUMAN PRESENCE ≠ ALTERNATE AUTHORITY PATH
```

Those are already explicit in v0.3. 

Retain the entire exact-operation-binding discipline:

```text
proposal
→ authorization
→ execution

X authorized
Y requested
→ reject
```

including total identity, effect-relevant parameters, shared identity definition, and fail-closed behavior. 

Retain:

```text
AUTHORIZED
DENIED
UNKNOWN / UNAVAILABLE
```

and the rule that UNKNOWN is not silently interpreted as DENIED or success. 

Retain replay/idempotency, ambiguous-outcome reconciliation, declaration lifecycle, and Growth Gate. These are useful precisely because they constrain future expansion.  

## What v0.3 now supersedes

The biggest mismatch is section 4:

```text
DISCOVER
→ FORMULATE
→ ASSIGN / CLAIM
→ WORK
→ REVIEW
→ NEXT WORK
```

with a consequential branch. 

Our locked design is now simpler:

```text
PLAN
→ HUMAN APPROVAL
→ WORK
```

There is no third mode, no separate workflow phase authority, and no Conductor-controlled plan iteration.

Likewise, v0.3 still contains **Autonomous vs Semi-Autonomous modes**. 

Those should no longer be normative workflow modes. They were superseded by the explicit two-mode model:

```text
PLAN
WORK
```

Human involvement is no longer a mode taxonomy. It is a plan-approval/interaction rule.

## The biggest semantic change

v0.3 currently says:

> "Agent/client classification of consequence remains advisory." 

That survives only with the newer grain correction:

> **An operation is consequential when the operation class it instantiates is declared effect-capable—capable of producing an observable state change outside the workflow—regardless of intent, phase, or invoking client.**

The declaration owner owns that classification record; the Agent does not.

This resolves the `shell`/`python` problem without making every use of a universal capability consequential.

## What is missing from v0.3

The new design adds something important that v0.3 does not really model:

```text
Conductor
    =
durable current project state
+
durable project history
+
persistent work graph
```

That is now a core design requirement.

v0.3 deliberately says it does not define a workflow database, and its persistence section keeps each participant's own records separate. 

That remains correct from the Loop Engineering perspective, but it needs a **new Conductor-specific capability definition**:

> Conductor remains the authority for coordination/project state, not for authorization, execution, or domain truth.

This is the one legitimate architectural delta that needs to be documented rather than quietly inferred.

## What I would do next

### Step 1 — Do not edit v0.3

Treat:

```text
Workflow Specification v0.3
```

as the **parent security/conformance lineage**.

Do not endlessly revise it to accommodate the new simplified workflow.

### Step 2 — Create `Workflow Design v1.0`

This becomes the actual frozen workflow design for the next implementation phase.

It should contain only:

```text
1. Scope
2. Role ownership
3. PLAN / WORK modes
4. Human approval semantics
5. Conductor project-state/history model
6. Minimal task model
7. READY frontier / atomic claim
8. Provenance + attribution
9. Plan/version + operation-class scope
10. Consequential-operation boundary
11. Solvent checkpoint
12. Executor boundary
13. Exact operation identity
14. Evidence / acceptance ownership
15. Deferred items
16. Known open items
17. Disposition log
18. Lineage / supersession statement
```

### Step 3 — Explicitly declare what v1.0 supersedes

The document should say something like:

```text
Workflow Design v1.0 supersedes the workflow-phase and
autonomy-mode portions of Workflow Specification v0.3.

Workflow Specification v0.3 remains the parent source for:
- authority semantics
- exact operation binding
- fail-closed enforcement
- authorization taxonomy
- evidence semantics
- replay/idempotency
- reconciliation
- declaration lifecycle
- conformance principles
```

This avoids two competing normative documents.

### Step 4 — Add the three remaining lock-level semantics

Before freezing v1.0:

```text
A. operation-class definition
B. declaration reference + version + content-hash + resolution check
C. effect-boundary enforcement is fail-closed independent of Agent behavior
```

The declaration structure in v0.3 is already strong enough to support this; it explicitly requires operation/class, effect capability, identity, validity, replay, outcome, and reconciliation semantics. 

### Step 5 — Freeze with a decision record

Record:

```text
Workflow Design v1.0
Owner
Date
Approval criteria
Lineage decision
Disposition log
Known open items
Decision = FROZEN
```

### Step 6 — Only then run Phase 1.5

Use the exact remaining battery:

```text
operation X authorized → operation Y executed
duplicate delivery
stale/invalid authorization
UNKNOWN / unavailable
stale claim × expired intent
```

The attached v0.3 already gives us the semantic basis for most of these tests.  

### Step 7 — Then Phase 2

Before the first real external effect:

```text
declaration pinned
operation identity pinned
pass/fail criteria preregistered
falsification conditions defined
real enforcement path verified
```

Then exercise the real external SOR.

---

## Bottom line

I would **not throw away Workflow Specification v0.3**.

I would treat it as the **security/conformance foundation** and build a smaller `Workflow Design v1.0` above it that changes only the workflow/control-plane model:

```text
v0.3
security + authority + execution semantics
                │
                ▼
v1.0
minimal workflow:
PLAN → HUMAN APPROVAL → WORK
                │
                ▼
Conductor = durable project state/history
Agent     = intelligence
Solvent   = authority
Executor  = effect
```

That is the cleanest next move, and it finally ends the review loop: **one final artifact, one final review, one freeze, then experiments.**
