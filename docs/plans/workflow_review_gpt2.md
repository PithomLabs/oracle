I reviewed the attached `Loop Engineering Requirements v0.3` adversarially, treating it as the current baseline rather than assuming the previous review findings were already resolved.

## Verdict: **AMBER — strong architecture, but not yet ready to freeze as the conformance baseline**

The document is substantially sound. It preserves the central separation and successfully avoids creating a fifth infrastructure role. 

But I would **not sign off v0.3 yet**, because several issues remain in the actual text. Most are fixable with a few normative clarifications; they do not require changing the architecture.

---

# 1. HIGH — exact authorization-to-operation binding is still missing

This remains the sharpest gap.

The document says:

> effect-capable execution fails closed without required authorization evidence

and:

> the effect-capable integration determines ... what authorization evidence is required. 

But it never explicitly requires:

```text
authorized operation == executed operation
```

That leaves the classic substitution attack:

```text
Agent proposes X
    ↓
Solvent authorizes X
    ↓
Executor receives Y + authorization evidence for X
    ↓
Executor performs Y
```

The current document can observe that an authorization exists, an execution exists, and the records correlate, without establishing that **Y is the operation Solvent authorized**.

### Required fix

Add a hard invariant:

> **Authorization evidence MUST be bound to the specific operation, target, and other authority-defining context it authorizes. An effect-capable execution boundary MUST reject execution when the requested operation does not correspond to the authorized operation/context.**

Do not reproduce Solvent's internal target/snapshot semantics here; just establish the cross-role binding invariant.

---

# 2. HIGH — the integration still effectively defines its own consequential boundary

Section 7 says:

> "The effect-capable integration determines what operations constitute externally consequential effects..." 

This is understandable for domain/tool agnosticism, but it creates a loophole:

```text
Integration author:
"Operation Z isn't consequential."
```

Therefore no authorization is required.

The specification has not established an objective baseline for what an externally consequential effect is.

### Better formulation

Define an external effect generically:

> An external effect is an operation that can cause a state change observable outside the requesting process or execution context.

Then require every effect-capable integration to declare:

```text
operation
→ produces external effect?
→ required authorization boundary
→ required evidence
```

This is **declaration metadata**, not a new policy engine.

The conformance suite can then inspect the declaration and test behavior against it.

---

# 3. HIGH — conformance observation is still not an actual access contract

The document says:

> Cross-role consequential interactions MUST expose enough evidence for a conformance harness... 

Good.

But **where does the harness get that evidence?**

You allow:

* Conductor task ID
* Solvent action-intent ID
* authorization reference
* executor execution ID
* external operation ID

but never require that a conforming test deployment make those records accessible to the harness.

This means an implementation could satisfy the semantic language while the actual test suite cannot observe the required facts.

### Fix

Add:

> **A conforming test deployment MUST expose the evidence required by each conformance scenario through existing owning-system interfaces or explicit test adapters.**

That is enough.

Do not create a workflow event bus or workflow database.

---

# 4. HIGH — correlation is defined, but proposal → authorization → execution binding is still ambiguous

You correctly say:

> Correlation is an association mechanism.
> Correlation MUST NOT be treated as authority/truth/semantic ownership. 

Excellent.

But correlation alone is insufficient for the security property above.

You need two different concepts:

```text
Correlation
= these records belong to the same workflow episode

Authorization binding
= this execution is the exact action authorized
```

The latter must be explicit.

I would add it rather than expanding the correlation mechanism.

---

# 5. HIGH — denial and Solvent unavailability still need different recovery semantics

Your failure model lists:

> authorization denial
> authorization unavailable/unknown 

but Section 17 only gives a denial → revision path.

These are semantically different:

```text
DENIED
→ authoritative decision exists
→ revise / abandon / request human review

UNKNOWN / UNAVAILABLE
→ no authoritative decision exists
→ retry / wait / escalate
→ NOT equivalent to denial
```

This distinction should be explicit before the conformance matrix is written.

Otherwise an Agent can learn to interpret Solvent outage as rejection.

---

# 6. MEDIUM — denial → revision has no bounded conformance semantics

Scenario E is good:

```text
proposal
→ denial
→ revision
→ resubmission
```



But autonomous mode could theoretically loop forever.

Do not add a generic workflow retry counter.

Instead specify for **conformance scenarios**:

> Revision loops MUST be bounded by the scenario's declared iteration limit.

Production retry/review policy remains external to Loop Engineering.

---

# 7. MEDIUM — ambiguous execution outcome needs a reconciliation contract

You correctly define:

```text
not attempted
attempted
succeeded
failed
ambiguous
```

and explicitly prevent `ambiguous` from collapsing into success/failure. 

That's excellent.

But what happens next?

Add:

> **An ambiguous outcome MUST follow the Executor/integration's declared reconciliation path and MUST NOT be treated as successful or failed merely to allow workflow progress.**

Possible reconciliation mechanisms can remain implementation-specific:

```text
status query
external reconciliation
human review
idempotent retry
manual investigation
```

No new workflow state machine.

---

# 8. MEDIUM — long-running authority model needs declared scope

The two models are useful:

```text
authorize → start → authorization applies

authorize → start → executor re-checks
```



But they don't specify:

* authority duration
* re-check interval, if applicable
* what happens after authority becomes stale
* whether already-running execution is allowed to finish

The requirements should require each integration to declare:

```text
authority validity model
recheck model, if any
mid-flight behavior
```

Don't prescribe heartbeat infrastructure.

---

# 9. MEDIUM — "Human authorization" remains slightly ambiguous

The role matrix says:

> Human may request/review/act through applicable authority boundary. 

This is mostly correct, but I would make one point explicit:

> **Human approval does not itself constitute authority; when authorization is required, the authority decision remains owned by Solvent or the explicitly designated authority system.**

That removes any possible interpretation that a human can directly authorize the Executor.

---

# 10. MEDIUM — human modification needs proposal identity/version semantics

The document permits:

> Human may review or modify the proposal before it reaches the authority boundary. 

That creates an important traceability requirement.

Conceptually:

```text
Proposal V1
    ↓ human revision
Proposal V2
    ↓
new authority decision
```

V2 should not silently mutate V1.

Add a requirement in the Workflow Specification later that a materially modified proposal is a distinct proposal/request instance with attributable provenance.

---

# 11. MEDIUM — tool agnosticism should distinguish client substitution from orchestration substitution

The document says a client may run under Temporal or another framework. 

But there are actually two proof points:

```text
Client substitution
GPT / Claude / scripted / human

Orchestration substitution
direct client / Temporal / another framework
```

Scenario D currently mainly exercises the first. 

For the BM-IST thesis, those should be measured separately.

---

# 12. LOW — ordinary work says SHOULD, conformance implies MUST

Section 4 isn't shown in this attachment's line range? Actually the document says in prior version `Ordinary work SHOULD proceed without Solvent`, while §16 now says:

> ordinary work can proceed without Solvent when no consequential action is involved. 

The intent is clear, but I'd make it internally normative:

> **Ordinary work that does not produce an external consequential effect MUST NOT require Solvent merely because the task exists.**

That captures the architecture more strongly and prevents accidental “everything goes through Solvent” adapters.

---

# 13. LOW — the document's acceptance criterion still says "Requirements v0.3 is satisfied"

Section 21 is in the requirements document and says:

> "Requirements v0.3 is satisfied when..." 

That's slightly awkward now that this is the requirements baseline itself.

I'd phrase it as:

> **Requirements v0.3 acceptance criteria**

Not a substantive issue.

---

# What I would *not* do

I would reject several possible "fixes" suggested by more aggressive reviews:

### Do not add a generic workflow router

Conductor should not inspect domain proposals to decide whether something is consequential.

### Do not add a mandatory authorization proxy/gateway for every Executor

That would risk creating a fifth infrastructure role.

### Do not make Conductor cancellation automatically revoke Solvent authority

That would make Conductor mutate an authority system it does not own.

### Do not add a generic `Abandoned` workflow state

Use existing Conductor lifecycle semantics unless BM-IST proves a real missing coordination primitive.

### Do not add an event store

The specification already correctly establishes that correlation can use existing/native identifiers. 

---

# Final disposition

### **AMBER — v0.3 is very close, but I would not freeze it yet.**

The four fixes I would consider **blocking** are:

```text
1. Exact authorization ↔ executed-operation binding.
2. Mandatory observability access for conformance deployments.
3. Explicit consequential-operation declaration / effect definition.
4. Denial vs authority-unavailable recovery semantics.
```

Then these should be tightened before deriving the conformance matrix:

```text
5. Ambiguous-outcome reconciliation.
6. Long-running authority validity declaration.
7. Human authorization wording.
8. Proposal revision/version attribution.
9. Client vs orchestration substitution distinction.
10. Bounded conformance revision loops.
```

Crucially, **none of those require a new runtime, database, authority engine, or UI**.

The architecture remains sound. The issue is that the specification still needs to close the **mechanical contract between authorization evidence and actual effect**, and make its promised conformance evidence genuinely observable.

Once those are fixed, I would move directly to the **Role / Boundary Matrix and Conformance Test Matrix**, rather than continuing an open-ended requirements cycle.
