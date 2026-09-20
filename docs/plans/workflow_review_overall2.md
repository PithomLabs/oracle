I consolidated the three adversarial reviews with my own review of the current `Loop Engineering Workflow Specification v0.2`. The result is quite clear:

# Consolidated verdict: **AMBER — architecture is sound, but v0.2 is not yet a safe conformance baseline**

The encouraging part is that the previous major holes are actually being closed. The document now has explicit operation binding, declarations, fault-injection concepts, correlation, reconciliation, tool substitution, and clearer ownership. 

But the newest reviews correctly identify that the remaining risks have **moved from architecture to the exact semantics of the new mechanisms**. The most important unresolved issue is the operation-binding model, followed by declaration lifecycle and completeness of the conformance suite. 

## 1. CRITICAL — operation binding is still not sufficiently defined

This is the strongest finding across the new reviews, and I agree.

The v0.2 design apparently now binds authorization to an operation, but the binding can still be weakened by two things:

### A. Operation identity is not necessarily total

A stable operation identity that omits a materially relevant parameter can still authorize one thing and execute another.

For example:

```text id="m65s2g"
Authorized:
    delete(file=/a, recursive=false)

Executed:
    delete(file=/a, recursive=true)
```

If `recursive` wasn't part of the identity, the binding technically passes.

### B. Operation classes are dangerous

If Solvent authorizes:

```text
class = "shared-compute"
```

and the Executor decides which operations belong to that class, the Executor has effectively acquired part of the authorization semantics.

The Z review identifies this precisely. 

### Required rule

The next revision should say:

> **Operation identity MUST be total over all effect-relevant parameters. A parameter is either included in the operation identity or explicitly declared immaterial. Any undeclared effect-relevant difference constitutes a different operation.**

And, if operation classes remain allowed:

> **Class membership MUST be explicitly declared and deterministically evaluated according to the declaration; the Executor MUST NOT invent class semantics.**

I would strongly prefer **exact operation binding for v1** and allow operation classes only as a later demonstrated need.

---

# 2. HIGH — declarations are now load-bearing, but their lifecycle is undefined

This is another strong finding.

The new architecture relies heavily on integration declarations:

```text
what operations are effect-capable
what evidence is required
what authority model applies
what replay model applies
```

But what happens when an integration changes?

The review correctly identifies the risk:

```text id="s6qzpy"
Executor v1:
    operation A requires auth

Executor v2:
    operation B added, but declaration unchanged

Conformance suite:
    still passes
```

That is certification drift.

### Add a declaration lifecycle rule

Something like:

> **Effect-capable integration declarations are versioned artifacts. A capability change MUST update the declaration and trigger re-conformance of the affected boundary. Undeclared effect-capable behavior discovered during review is a conformance finding. The applicable authority owner MAY inspect the declaration.**

This doesn't create a policy engine. It simply makes the declared contract auditable.

That is one of the best findings in the new review. 

---

# 3. HIGH — exact-operation binding must be carried by the authorization evidence itself

The Qwen review makes an excellent refinement:

The Executor can't reliably compare X and Y if X is only implicit in Solvent's state.

The authorization evidence should carry the binding.

Conceptually:

```text id="2k1h0u"
Solvent
  ↓
Authorization Evidence {
    authorized_operation: X
    ...
  }
  ↓
Executor
  ↓
requested_operation = Y

X != Y
    ↓
REJECT
```

That is materially stronger than:

```text
Executor somehow looks up X
```

because it eliminates a second lookup/race/mapping problem.

The recommendation to make the evidence structurally carry the authorized operation is sound. 

One qualification: **do not mandate cryptography yet**. The contract should require verifiable binding, while the concrete mechanism can be a signed artifact, Solvent validation callback, or another trusted mechanism appropriate to the deployment.

---

# 4. HIGH — conformance suite needs production negative testing

The review correctly distinguishes:

```text
TestExecutor works
```

from:

```text
production Executor actually rejects bad authorization
```

Those are not the same proof.

The conformance system therefore needs two layers:

### Protocol fixture conformance

Controlled `TestSolvent`, `TestExecutor`, fault injection, etc.

### Integration conformance

Run the same negative scenarios against the actual integration:

```text
missing authorization
invalid authorization
stale authorization
mismatched operation
duplicate delivery
```

and observe that the **real external effect does not happen**.

The Qwen recommendation is exactly right here. 

This does not require a special production runtime. A test/dry-run environment or reversible target is sufficient.

---

# 5. HIGH — the conformance suite still undercovers the Requirements

The Z review identified an important cross-document consistency problem: the Requirements demand minimum properties that the Scenario set does not yet actually exercise.

Specifically:

```text id="6p9t81"
long-running work does not block unrelated Conductor work
workflow phases do not become persisted state
```

are requirements, but the current scenarios do not explicitly test them. 

That means the next Conformance Matrix could silently omit requirements from its parent document.

### Required addition

Add at least:

```text id="g7d9ag"
Scenario F — Long-running execution
Scenario G — Persistence/phase inspection
```

For F:

```text
start long-running Executor
→ Conductor remains able to progress unrelated work
→ long-running work eventually reconciles
```

For G:

```text inspect declared persistence
→ no workflow phase is authoritative persisted state
```

The second can be an inspection/declaration test rather than a black-box runtime test.

---

# 6. MEDIUM — UNKNOWN / UNAVAILABLE authorization needs temporal semantics

The Qwen review points out that `UNKNOWN` can persist forever. 

I agree with the problem, but I do **not** agree that the universal answer should be:

```text unknown after T
→ automatically deny
```

That is too prescriptive.

Instead:

> **An integration MUST declare the maximum age or waiting policy for unresolved authorization before the request requires explicit retry, escalation, revision, or abandonment handling.**

This lets different deployments choose their semantics while preventing infinite silent limbo.

---

# 7. MEDIUM — long-running authorization-validity needs declaration + evidence

The current model permits:

```text id="43w8wq"
authorize → start → authorization covers execution attempt
```

or:

```text id="s0ff4r"
authorize → start → executor re-checks
```



That's reasonable.

What remains is to require declaration of:

```text
validity scope
recheck policy, if any
mid-flight revocation semantics
```

The conformance scenario should then exercise the declared model.

This avoids inventing a heartbeat engine while still making TOCTOU behavior explicit.

---

# 8. MEDIUM — ambiguous outcome needs an actual reconciliation owner

The specification now says an Executor has a declared reconciliation path, which is good.

But the actual ownership should be explicit:

```text id="2i6pyn"
Executor/integration
    owns ambiguous-effect reconciliation
Conductor
    records/co-ordinates
Agent
    interprets available evidence
Human/domain
    may review
```

The Qwen recommendation to force automatic human escalation is too rigid. The architecture should permit:

```text
status query
idempotent reconciliation
human review
external system lookup
```

depending on the Executor.

The important invariant is simply:

> **Ambiguous MUST NOT silently become success/failure.**

---

# 9. MEDIUM — agent-as-orchestrator remains a legitimate risk

Gemini and the earlier reviews both identify this.

If the client must independently poll:

```text
Conductor
Solvent
Executor
```

and reconstruct everything in its own context, we risk making the Agent a hidden workflow engine. 

I would **not solve this with a mandatory workflow aggregator**.

Instead define a minimal observation rule:

> **A conforming client must be able to obtain the evidence required by the current workflow step through the owning interfaces. The mechanism may be polling, request/response, subscription, MCP, orchestration-framework activity, or human observation.**

That preserves tool neutrality.

---

# 10. MEDIUM — do not mandate an Egress Gateway

Gemini proposes a generic gateway/sidecar to avoid every Executor implementing authorization checks. 

I would **not make that normative**.

A shared gateway can absolutely be a practical implementation, but making it a required Loop Engineering component creates exactly the infrastructure layer we're trying to avoid.

The correct contract is:

```text
Every effect-capable integration
    MUST enforce the applicable authority boundary
```

A deployment is free to implement that with:

```text
native wrapper
adapter
proxy
gateway
sidecar
library
direct callback
```

That preserves implementation agnosticism.

---

# 11. MEDIUM — deterministic replay is not a universal workflow requirement

Qwen suggests observation interfaces should be deterministic/idempotent for Temporal. 

The underlying concern is valid, but the proposed requirement is too broad.

Temporal's deterministic replay is an **orchestration-framework-specific property**.

The correct requirement is:

> **A client/orchestration integration MUST adapt observation semantics to the determinism requirements of its chosen orchestration framework without changing the Loop Engineering contract.**

Otherwise Loop Engineering starts prescribing Temporal semantics to everyone.

---

# 12. LOW — autonomous/semi-autonomous mode distinction still needs tightening

The current mode definitions allow explicit policy to introduce human intervention into autonomous operation.

That makes the terms somewhat fuzzy.

I would define:

```text id="2rt5qg"
Autonomous:
No human action is required at the declared workflow boundaries.

Semi-autonomous:
At least one declared boundary requires or permits human action before continuation.
```

Human intervention can still happen unexpectedly in autonomous mode; the distinction is about **required workflow progression**, not whether humans are physically capable of intervening.

---

# 13. LOW — Growth Gate needs the "smuggling" test

This came up repeatedly and remains useful.

Add:

> **The proposed change MUST NOT merely relocate prohibited workflow functionality into Conductor, Solvent, Executor, or an integration under another name.**

That protects against:

```text
"No workflow engine"
    +
"we'll just put workflow state inside the Executor"
```

The principle is already present; make it explicit.

---

# The most important combined conclusion

Across all three reviews, the architecture itself is **not** under attack anymore.

The attack surface has moved to:

```text id="g51z8u"
DECLARATIONS
BINDING
OBSERVABILITY
TIME
RECONCILIATION
```

That is exactly where we should now focus.

And there is an important architectural lesson from these reviews:

> **The answer is not to add more infrastructure. The answer is to make the cross-role contracts more explicit.**

This is consistent with the strongest parts of v0.2. 

---

# Consolidated action list

I would make these changes before declaring v0.2 closed:

### Blocking

**1. Total operation identity**

Every effect-relevant parameter is either part of the identity or explicitly declared immaterial.

**2. Exact authorization binding**

Authorization evidence itself carries the authorized operation/context; Executor verifies requested operation against it.

**3. Declaration lifecycle**

Declarations are versioned, re-conformance triggered by capability changes, and undeclared effect-capable behavior is a conformance finding.

**4. Production negative conformance**

Real effect-capable integrations must be testable against missing/invalid/stale/mismatched authority.

**5. Missing conformance scenarios**

Add explicit long-running and persistence/state inspection scenarios.

### Important but non-blocking

**6. UNKNOWN temporal policy**

Prevent indefinite unresolved authorization.

**7. Ambiguous reconciliation ownership**

Executor/integration owns reconciliation; workflow doesn't invent a universal recovery engine.

**8. Authority-validity declaration**

Long-running execution must declare validity/re-check behavior.

**9. Autonomous/semi-autonomous semantics**

Define the distinction in terms of required human intervention.

**10. Growth Gate anti-smuggling rule**

Prevent prohibited functionality being relocated into an existing component.

---

## What I would explicitly reject

```text
❌ mandatory Egress Gateway
❌ mandatory workflow event bus
❌ mandatory workflow aggregator
❌ automatic Conductor→Solvent revocation
❌ mandatory human escalation for every ambiguous result
❌ Temporal-specific deterministic observation semantics
❌ new workflow state machine
❌ new workflow database
❌ new Loop Engine runtime
```

Those are implementation choices or architecture expansions, not necessary corrections to the contract.

## Final verdict

**AMBER — one focused mechanical revision remains.**

The right next move is **not another broad architecture redesign**. Tighten the contract around:

```text
operation identity
      ↓
authorization evidence
      ↓
exact execution binding
      ↓
observable evidence
      ↓
negative conformance
```

Then close the declaration lifecycle and the missing conformance scenarios.

After that, the specification should be stable enough to derive the **Role/Boundary Matrix + Conformance Test Matrix + Agent Skill + Conformance Harness + POC Runbook** without inventing hidden architecture.

That is the point where I'd stop adversarially expanding the specification and start testing the actual BM-IST workflow.
