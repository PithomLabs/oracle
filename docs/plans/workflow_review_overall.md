I consolidated the three attached adversarial reviews with my own review of `Workflow Specification v0.1`. The combined conclusion is materially sharper than any one review alone.

## Consolidated verdict: **RED for conformance-readiness, GREEN for architecture**

The underlying architecture is still sound. The document successfully preserves:

```text
Agent      = agency
Conductor  = coordination
Solvent    = authority
Executor   = effect
Human      = intervention/review
Domain     = semantic authority
```

and it explicitly avoids a workflow runtime, workflow database, second authority system, and second task system. 

But the combined reviews identify a more fundamental problem:

> **The specification now says enough about what must be true, but not enough about how a conforming implementation proves those facts without making hidden architectural decisions.**

The most serious findings are C1/C2/C3 from the Z review, which are reinforced by Qwen's routing/integration critique and Gemini's enforcement/state-reconciliation concerns.   

I would **not proceed to implement the conformance harness yet**. The Workflow Specification needs one focused revision first.

---

# 1. Critical consensus finding: the authority boundary is not mechanically closed

This is the strongest finding across the reviews.

The specification says:

> "The execution integration defines how authorization evidence is checked." 

That's too open-ended.

Z correctly identifies the sharp attack:

> Solvent authorizes X, but the Executor receives X's authorization evidence and performs Y. The current invariants do not explicitly require `authorized operation == executed operation`. 

That's a **real security gap**, not merely documentation detail.

### Required invariant

Add:

> **Authorization evidence MUST be bound to the exact operation it authorizes. An effect-capable execution boundary MUST reject an execution request when the requested operation does not correspond to the authorized operation.**

Conceptually:

```text
Proposal X
   ↓
Solvent authorizes X
   ↓
authorization evidence bound to X
   ↓
Executor receives Y
   ↓
Y ≠ X
   ↓
REJECT
```

This is essential.

---

# 2. Critical consensus finding: "consequential" still has a routing problem

Qwen's critique is also legitimate: if Agent classification is advisory, **how does the system know to invoke Solvent before an effect?** 

But I would **not** solve this with a new centralized router.

The right resolution is a distinction the current specification does not yet make explicitly:

### Detection vs enforcement

```text
Detection / routing:
    Agent or domain-specific client proposes a consequential action.

Coordination:
    Conductor records/co-ordinates the proposal.

Authority:
    Solvent evaluates it.

Enforcement:
    Executor refuses execution without matching authorization.
```

Conductor does **not** need to understand BM-IST semantics.

The effect-capable integration knows that a particular operation requires authority and requests/validates the corresponding authorization evidence.

So the next spec revision should explicitly say:

> **Conductor does not classify domain consequences. A client/integration may declare that an operation requires the consequential boundary; the effect-capable boundary remains the final fail-closed enforcement point.**

That resolves the routing paradox without creating a policy router.

---

# 3. Critical consensus finding: the conformance harness still lacks a defined observation mechanism

This is the other major issue.

The specification says the harness must establish:

* proposal
* authorization
* execution attempt
* outcome
* ordering/correlation 

But it doesn't yet specify **how the harness gets access to those facts**.

Z calls this the observability gap. 

Qwen similarly points out that "existing/native references" are insufficient if they are inaccessible to the conformance harness. 

### Required fix

Do not build a workflow event store.

Instead define a **Conformance Observation Interface**, which can simply mean:

```text
Conductor API / activity
Solvent API / authority record
Executor test interface
```

The requirement should be:

> **A conforming test deployment MUST expose the evidence required by the conformance scenarios through the existing owning interfaces or explicit test adapters.**

This lets the harness inspect facts without turning them into a new authoritative state system.

That also gives us a legitimate place for deterministic fault-injecting test doubles.

---

# 4. Critical consensus finding: correlation is not enough unless identity binding is explicit

The spec currently does something good:

> correlation ≠ authority

and lists native identifiers such as task ID, Action Intent ID, execution ID, and external operation ID. 

But the reviews correctly identify that we need **two separate concepts**:

```text
correlation
= these records belong to the same workflow episode

authorization binding
= this execution is the specific action Solvent authorized
```

The second cannot be satisfied merely by sharing a correlation ID.

So I would add:

> **Correlation identifies related records; authorization binding establishes whether the execution request is the exact action authorized by the authority system.**

This is probably the single most important vocabulary refinement in the whole specification.

---

# 5. High: conformance needs deterministic fault-injection actors

Z's C3 is correct. The Requirements demand things like:

> invalid authorization evidence prevents execution

but the specification doesn't define how the test harness can deliberately produce invalid evidence, ambiguous outcomes, replay, etc. 

The solution is not to expose Solvent internals.

Add to the conformance model:

```text
Conformance actors may include:
- real component
- deterministic test double
- fault-injecting test executor
- controlled authority fixture
```

For example:

```text
TestSolvent
    deny
    return unknown
    return authorization for X
    return authorization for stale X
```

```text
TestExecutor
    accept
    reject invalid evidence
    produce failure
    produce ambiguous outcome
    simulate duplicate delivery
```

These are **test fixtures**, not production infrastructure.

That distinction should be explicit.

---

# 6. High: denial ≠ unavailability

This was correctly raised by Z. 

Current failure semantics list:

```text
authorization denial
authorization unavailable/unknown
```

but then largely group them under generic recovery.

They need different semantics:

```text id="n9h5j5"
DENIED
→ authority decision exists
→ Agent may revise/reformulate

UNKNOWN / UNAVAILABLE
→ no authoritative decision exists
→ do not reinterpret as denial
→ retry/escalate according to integration policy
```

An autonomous Agent must not learn:

```text "Solvent timed out, therefore my proposal was denied."
```

That is a genuine semantic distinction.

---

# 7. High: denial→revision needs termination semantics

Scenario E is a good addition:

```text
proposal
→ deny
→ revise
→ resubmit
```



But as Z points out, nothing says when autonomous revision stops. 

I would **not** introduce a workflow-level retry counter.

Instead require:

> **Conformance scenarios MUST define a bounded number of revision iterations; production systems MAY impose their own retry/review policies.**

And:

> **Denial count is evidence, not workflow authority.**

That keeps policy outside Loop Engineering.

---

# 8. High: long-running authority validity remains underspecified

The current spec offers two models:

```text
authorize → start → authorization applies
authorize → start → executor re-checks
```



That is good.

But Gemini correctly points out the missing question:

> What happens if authority is revoked or the relevant state changes during the execution? 

The right answer is **not** "mandate heartbeats."

Instead the Workflow Specification should require each long-running integration to declare:

```text
authority validity model
revalidation model, if any
mid-flight revocation behavior, if any
```

Then the conformance matrix tests that declared model.

---

# 9. Medium: ambiguous execution outcome is correctly modeled, but recovery must remain external

The current spec properly defines:

```text
not attempted
attempted
succeeded
failed
ambiguous
```

and says the Executor/external system of record owns execution truth. 

That is excellent.

Gemini's complaint that an ambiguous result needs retry/reconciliation is valid, but its proposed mandatory human escalation is too prescriptive. 

The correct requirement is:

> **An ambiguous outcome MUST enter the Executor/integration's declared reconciliation path and MUST NOT be automatically interpreted as success or failure.**

The reconciliation mechanism could be:

```text
status query
human escalation
idempotent retry
external reconciliation
manual investigation
```

depending on the Executor.

Loop Engineering defines the distinction; it does not own the recovery engine.

---

# 10. Medium: do not add a mandatory Conductor → Solvent cancellation/revocation path

Gemini recommends explicit revocation when Conductor cancels a task. 

I reject this recommendation.

The current spec explicitly says:

> "No Conductor authority over Solvent revocation is implied..." 

That is correct.

Otherwise we'd create:

```text
Conductor cancellation
     ↓
mutate Solvent authority
```

which weakens the authority ownership boundary.

Instead, define integration semantics around **execution eligibility**, without making Conductor a revocation controller.

---

# 11. Medium: Agent reconciliation must not become hidden orchestration

Qwen argues that the current design pushes state reconciliation onto the Agent. 

This is a useful warning, but I would not accept the proposed "aggregated status webhook" as a requirement.

The workflow contract should simply say:

> **Clients MUST have access to the state/evidence necessary to continue the workflow through the owning interfaces.**

How that happens is implementation-specific:

```text
poll
MCP request
HTTP request
subscription
Temporal activity
human observation
```

The Agent isn't necessarily a distributed state reconciliation engine; it is simply responsible for interpreting the available evidence.

That distinction should be clarified.

---

# 12. Medium: no need for an "Abandoned" workflow state

Qwen proposes an `Abandoned` state because the Executor may continue after a human stops caring about the task. 

I would **not** add this to the Workflow Specification as a new generic state.

The existing Conductor already has lifecycle semantics, including cancellation/release, while Executor execution state is separate.

The workflow specification can express:

```text
coordination stops progressing
≠
external execution stops
```

without inventing another canonical state.

If BM-IST proves that an explicit "superseded" concept is necessary, that should be evaluated through the Growth Gate against existing Conductor lifecycle.

---

# 13. Medium: human-modified proposals need attribution/binding

Z correctly points out that Section 7 allows humans to modify proposals but doesn't say how that change is represented. 

The solution should be simple:

```text
Proposal V1
    ↓
human revision
    ↓
Proposal V2
    ↓
new authorization decision
```

And V2 must carry its own proposal/reference identity.

Do not mutate V1 into something ambiguous.

This naturally fits the denial/revision loop already defined.

---

# 14. Low/Medium: the Agent Skill question needs to be removed from this spec

Z raises a legitimate concern: "Agent Skill" appears only at the end as a future deliverable, while Requirements v0.3 explicitly bans a workflow SDK. 

I don't think this means the Agent Skill itself is wrong. It means the Workflow Specification should define it precisely as:

> **a client-side behavioral instruction artifact derived from this specification; not a runtime, SDK, authority component, or required implementation dependency.**

That resolves the concern.

The conformance harness should likewise be defined as:

> **test-only infrastructure with no production runtime obligations.**

That's important and should be explicit.

---

# 15. Medium: Scenario D needs two axes of substitution

The current specification says a client can run under Temporal or another orchestration framework. 

But Scenario D currently primarily demonstrates client substitution. 

For the BM-IST thesis, distinguish:

```text
Client substitution
GPT / Claude / script / human

Orchestration substitution
direct client / Temporal / another framework
```

That is important because **tool agnosticism** and **orchestration-runtime agnosticism** are related but not identical claims.

---

# 16. The most important combined correction

All three reviews and my review converge on one deeper point:

## The Workflow Specification currently has three kinds of contracts, but does not distinguish them explicitly enough.

We should formalize:

### A. Semantic contract

What the roles mean.

```text
Agent = agency
Conductor = coordination
Solvent = authority
Executor = effect
```

### B. Boundary contract

What must happen when information/action crosses between roles.

```text
proposal
authorization
execution
result
```

with:

```text
authorization ↔ exact operation binding
```

### C. Conformance evidence contract

How we prove A and B happened.

```text
identifiers
correlation
observable evidence
test doubles
pass/fail observations
```

That is the missing conceptual layer.

---

# What I would change in v0.2

I would **not rewrite the whole document**.

Add these specific requirements:

```text id="5gwgk4"
1. Objective external-effect definition.
2. Mandatory declaration of consequential operations by each integration.
3. Exact proposal/action ↔ authorization binding.
4. Defined authorization-evidence verification contract.
5. Conformance observation access mechanism.
6. Deterministic fault-injecting test actors.
7. Denial vs unavailability recovery taxonomy.
8. Bounded conformance revision loops.
9. Declared long-running authority validity/recheck model.
10. Ambiguous-outcome reconciliation requirement.
11. Human proposal modification creates a new proposal identity/version.
12. Explicit distinction between client substitution and orchestration substitution.
13. Define Agent Skill as a client-side behavioral artifact, not SDK/runtime.
14. Define Conformance Harness as test-only infrastructure.
```

And explicitly **do not add**:

```text id="dyiyup"
❌ workflow event store
❌ workflow runtime
❌ Conductor→Solvent revocation engine
❌ mandatory heartbeat engine
❌ generic workflow router
❌ new workflow state machine
❌ mandatory proxy/gateway for every Executor
❌ mandatory human escalation for ambiguous outcomes
```

---

# One particularly important correction to the external reviews

Both Gemini and Qwen are correct that heterogeneous Executors create integration burden. But their proposed answer—mandatory proxy/gateway—would be dangerous because it would quietly make the proxy the missing fifth infrastructure role.  

The better rule is:

> **Every effect-capable integration must declare and enforce its authority boundary; shared adapters/proxies may be used where useful, but they are an implementation choice, not a required Loop Engineering role.**

That preserves tool agnosticism.

---

# Final consolidated verdict

### **AMBER / RED for v0.1 conformance-readiness**

Not because the architecture is wrong.

Quite the opposite: the architecture has survived the reviews.

The problem is that **the specification is now strong enough to expose the mechanical contracts it still hasn't defined**.

The highest-priority fixes are:

```text
          SEMANTIC
             │
             ▼
       exact operation
             │
             ▼
      AUTHORIZATION BINDING
             │
             ▼
       EXECUTION CHECK
             │
             ▼
        ACTUAL EFFECT
             │
             ▼
      OBSERVABLE EVIDENCE
             │
             ▼
       CONFORMANCE TEST
```

Once those contracts are explicit, we can derive the Role/Boundary Matrix and Conformance Test Matrix without silently inventing architecture.

**I would therefore revise Workflow Specification v0.1 to v0.2 before writing any code.** The revision should remain small and mechanical; the architectural model itself should not change.
