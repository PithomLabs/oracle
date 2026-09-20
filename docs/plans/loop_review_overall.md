I consolidated the three attached reviews with my own review of `Loop Engineering Requirements v0.2`. The combined result is more useful than simply merging their issue lists, because some findings are genuinely necessary while others would accidentally push the specification back toward the infrastructure layer it was designed to avoid.

## Consolidated verdict: **AMBER**

The architecture is sound. The document successfully established the central discipline:

```text
Agent       = agency
Conductor   = coordination
Solvent     = authority
Executor    = effect
Human       = participant/intervener
Domain      = meaning/truth outside infrastructure
```

It also correctly bans a workflow runtime, persisted workflow state, a second authority system, a second task system, and workflow-specific UI/database. 

The remaining issue is that **v0.2 is conceptually rigorous but not yet fully conformance-ready**. The strongest external review identified four residual enforcement gaps; another identified two particularly important extensions—tool substitution and denial/revision; the most adversarial review identified deeper protocol-level gaps around classification, observability, correlation, failure, replay, intervention, and versioning.   

The good news is that **most of these can be closed with specification vocabulary and test contracts, not new infrastructure.**

---

# 1. The three most important issues

## A. Consequential boundary is not yet defined strongly enough

This is the most important finding from the adversarial review.

The current requirement says an Agent may propose consequential action and that effect-capable execution paths must enforce applicable Solvent authorization. 

But this leaves an uncomfortable question:

> Who decides that something is consequential?

If the Agent decides, a buggy/adversarial Agent can simply classify an effect as ordinary.

The strongest review correctly points out this circularity. 

### Consolidated fix

Do **not** create a workflow-level classifier.

Instead establish:

> **An Agent's classification of an action is advisory. Any effect-capable execution path must fail closed unless the required authorization evidence is present and valid.**

That gives you:

```text
Agent says "ordinary"
        ↓
Executor sees external effect
        ↓
authorization required by effect boundary
        ↓
missing/invalid authorization
        ↓
REJECT
```

This also answers the Gemini concern that there must be an actual enforcement point. 

However, I would **not** adopt the wording "authorization token" as a normative requirement. That is too implementation-specific for a tool-agnostic specification. Say **authorization evidence/reference**, leaving the actual mechanism to Solvent/Executor integration.

---

# 2. Conformance needs an observable contract

This is the strongest criticism from the most adversarial review.

The requirements promise that Loop Engineering will be "observable and testable," and Section 13 defines behaviors such as authorization being distinct from execution.  

But the document does not yet define **what a conformance harness actually observes**.

The review correctly calls this a blocking problem. 

### Consolidated fix

Do **not** add a workflow event bus.

Instead define a minimal abstract **conformance evidence model**:

```text
proposal/request reference
authorization reference
execution reference
result/outcome reference
correlation reference
```

The specification does not need to dictate where these are persisted.

It only needs to establish:

> Cross-role interactions must expose enough evidence for a conformance harness to determine which interaction occurred, in what relationship, and with what outcome.

That can use existing Conductor activity, Solvent records, Executor records, and provider-specific adapters.

This preserves the no-workflow-ledger rule.

---

# 3. Correlation is missing and is genuinely necessary

I agree strongly with this finding.

Without some common correlation reference, the following cannot be demonstrated reliably:

```text
this proposal
    ↓
this authorization
    ↓
this execution
    ↓
this result
```

The review correctly notes that otherwise the system cannot convincingly prove ordering or distinguish activity from effect. 

### Consolidated fix

Add a lightweight requirement:

> **Cross-role consequential interactions MUST carry sufficient correlation references to associate the proposal, authorization, execution attempt, and resulting outcome without introducing a workflow-owned ledger.**

Important: this does **not** mean creating `workflow_id`, `workflow_event_id`, etc. as a new workflow authority system.

Use existing/native identifiers where possible:

```text
task_id
action_intent_id
authorization reference
execution_id
external operation ID
```

plus a conformance-run correlation reference where necessary.

---

# 4. Domain needs clarification, not another role

The strongest review correctly found an internal inconsistency.

Section 1 lists Domain alongside the four roles and Human, but the ownership matrix doesn't include Domain.  

This makes the document fail its own "role boundaries are unambiguous" criterion.

### Recommended resolution

Do **not** make Domain a seventh infrastructure role.

Instead define:

```text
Infrastructure roles:
Agent
Conductor
Solvent
Executor

Participant:
Human

External semantic authority:
Domain
```

The domain may be represented by a human, domain system, verifier, scientific method, external service, etc.

That is more faithful to the original architectural principle that **domain meaning stays outside infrastructure**. 

---

# 5. Human intervention is now good, but intervention semantics need one more layer

The current document correctly says intervention depends on the capabilities of the underlying component and explicitly acknowledges that effects cannot be retroactively undone. 

That's good.

The adversarial review correctly identifies the remaining question: what actually is an intervention?

The current scenarios mention:

```text
redirect
reject
revise
```

but do not define their semantics. 

### Fix

The Workflow Specification—not necessarily Requirements v0.3—should define a small vocabulary:

```text
review
reject
revise
redirect
cancel-before-effect
cancel-in-flight-when-supported
```

Do not introduce `pause` unless an actual component supports it.

For concurrency, the rule should be:

> **A coordination intervention never overrides an already-established authority decision or external effect by itself.**

That prevents human intervention from becoming a shadow authority path.

---

# 6. Do NOT adopt the proposed automatic cancellation → Solvent revocation

This is where I disagree with the Gemini review.

It proposes:

> Conductor task cancellation/redirection MUST trigger revocation of linked Solvent Action Intent. 

I would **not put this into the Requirements**.

That would create:

```text
Conductor
   ↓
must mutate/control Solvent authority state
```

which conflicts with the architecture we've deliberately established.

Instead:

```text
Conductor cancellation
        ↓
coordination fact

Solvent authorization
        ↓
authority fact

Executor
        ↓
must enforce the applicable authorization/revocation semantics
```

If a workflow later demonstrates a real requirement for coordinated cancellation, the integration contract can define how it works without making Conductor the owner of authority.

This is an important example of why simply accepting every adversarial recommendation would be counterproductive.

---

# 7. Long-running execution: accept the concern, reject mandatory heartbeats

The Gemini review correctly spots a real semantic question: what happens when execution lasts hours? 

But I would **not** mandate periodic Solvent heartbeat/re-authorization.

That would be premature infrastructure.

The right requirement is:

> **Each long-running effect must declare its authority-validity model at the execution boundary.**

For example:

```text
Mode A:
authorize → start → authority applies to execution attempt

Mode B:
authorize → start → executor performs defined authority re-checks
```

The workflow specification defines the distinction.

The Executor/integration determines which model applies.

This preserves:

```text
authorization
≠
execution
```

while avoiding a workflow-owned heartbeat architecture.

---

# 8. Replay and idempotency are real cross-boundary concerns

This finding should be accepted.

The adversarial review correctly points out that once authorization and execution are separated, there is potential for:

```text
authorized
   ↓
delivery duplicated
   ↓
executor runs twice
```

That is not merely an internal Solvent or Executor detail—it is a boundary concern. 

### Requirement to add

Something like:

> **Effect-capable integrations MUST define their replay/idempotency behavior for repeated delivery of the same authorized execution request.**

But don't mandate one technical mechanism.

It could be:

```text
idempotency key
single-use authorization reference
executor deduplication
external transaction key
```

The implementation belongs to the Executor/integration.

---

# 9. Ambiguous outcome needs to be explicitly modeled

This is another very good finding.

The requirements currently distinguish:

```text
authorization granted
execution started
execution completed
result accepted
```

which is excellent. 

But there is another state:

```text
execution attempted
    ↓
connection lost
    ↓
UNKNOWN whether effect happened
```

The workflow must not convert that into either success or failure.

### Add:

> **An executor may report an ambiguous outcome when external effect occurrence cannot be conclusively established. Ambiguous outcome MUST remain distinct from execution failure and success.**

And ownership:

> **The Executor or external system of record is authoritative for whether the external effect occurred; Conductor activity is never sufficient evidence.**

That fits the existing invariant perfectly. 

---

# 10. The three scenarios should become five

Qwen's recommendation is good here.

The current A/B/C scenarios are necessary but don't fully demonstrate the two stated BM-IST goals. 

I'd use:

### A — Ordinary Research

No Solvent.

### B — Consequential Shared Computation

Agent → Conductor → Solvent → Executor → result.

### C — Human Intervention

Defined intervention boundary.

### D — Tool Substitution

Same scenario B executed with:

```text
LLM agent
scripted client
human/curl
```

The semantics remain the same.

### E — Denial → Revision

Agent proposes → Solvent denies → Agent revises proposal → resubmits.

That last scenario is particularly valuable because it demonstrates the **loop** rather than merely a pipeline.

---

# 11. Scenario definitions need to become actual test cases

The adversarial review is right that A/B/C are presently sketches rather than conformance tests. 

For the Workflow Specification / Conformance Matrix, each scenario should have:

```text
Setup
Actors
Inputs
Steps
Expected observations
Pass conditions
Failure conditions
```

And deterministic protocol conformance should use scripted clients.

LLM behavior can be measured separately as advisory behavioral adherence, which the current requirements already correctly distinguish. 

---

# 12. Boundary persistence needs tighter wording

Both the current requirements and Gemini review touch this.

The current catalog says:

> persistence location, if any

while also banning a workflow-specific database. 

This is potentially ambiguous.

### Fix

Add:

> **Persistence location identifies where an existing owning system records the relevant fact; it does not authorize creation of workflow-specific persistence. Persistence location does not imply semantic ownership.**

That is enough.

I would **not** restrict persistence only to Conductor activity or Solvent audit tables, because the Executor/external system may legitimately own execution records.

---

# 13. Growth Gate process should be lightweight

The adversarial review is right that the Growth Gate currently describes criteria but not how the decision is recorded. 

Don't create a governance bureaucracy.

Require only:

```text
Observed POC limitation
→ evidence
→ proposed change
→ why adapter/integration/docs cannot solve it
→ decision
→ owner/date
```

That's enough.

---

# 14. Versioning needs clarification

This is also valid.

The document says:

> The implementation technology may change; the cross-role contract may not. 

But the document itself is `v0.2`.

The intended meaning should be:

```text
Core invariants = stable
Specification versions = may evolve
```

Add a lightweight policy:

> **Changes to normative semantics require a specification version increment and re-running conformance. Additive clarifications that do not change semantics may be documented without changing the major contract.**

Since this is v0.x, don't pretend to have a mature long-term compatibility regime yet.

---

# 15. Low-level findings

These are worth incorporating without allowing them to dominate the design.

### Statistical adherence

State explicitly:

> Behavioral adherence metrics are advisory and non-gating for conformance v0.1.

That resolves the dangling "MAY be measured." 

### Naming

I would **not rename Conductor or Solvent**. Those are already established architectural names. A small terminology note is sufficient.

### Existing systems

Acceptance criterion 9 should name the baseline explicitly:

```text
Existing Conductor
Existing Solvent
A selected Executor reference implementation
```

rather than saying merely "existing systems."

---

# Consolidated disposition

| Finding                                | Consensus | Disposition                           |
| -------------------------------------- | --------- | ------------------------------------- |
| Consequential classification           | Strong    | **Fix**                               |
| Executor/effect-boundary enforcement   | Strong    | **Fix**, but don't mandate "token"    |
| Observable conformance surface         | Strong    | **Fix**                               |
| Correlation references                 | Strong    | **Fix**                               |
| Domain ownership ambiguity             | Strong    | **Fix**                               |
| Intervention vocabulary/races          | Valid     | **Fix in Workflow Spec**              |
| Ambiguous execution outcome            | Valid     | **Fix**                               |
| Replay/idempotency                     | Valid     | **Fix at boundary-contract level**    |
| Long-running authority                 | Valid     | **Clarify**, don't mandate heartbeats |
| Automatic Conductor→Solvent revocation | Overreach | **Reject**                            |
| Boundary persistence ambiguity         | Valid     | **Fix wording**                       |
| Scenarios A–C too thin                 | Valid     | **Expand**                            |
| Scenario D tool substitution           | Strong    | **Add**                               |
| Scenario E denial/revision             | Strong    | **Add**                               |
| Growth Gate process                    | Valid     | **Lightweight fix**                   |
| Versioning policy                      | Valid     | **Lightweight fix**                   |
| Statistical adherence                  | Valid     | **Clarify advisory**                  |
| Renaming Conductor/Solvent             | Low-value | **Reject**                            |

## The consolidated architectural principle

The four reviews collectively point to one important distinction:

> **Loop Engineering must specify enough interaction semantics to make the four-role architecture testable, but not so much that it becomes a workflow runtime.**

That means the next version should gain:

```text
classification semantics
boundary enforcement
correlation
observable evidence
failure/outcome semantics
replay/idempotency expectations
intervention semantics
versioning
```

but **not** gain:

```text
workflow database
workflow event store
scheduler
heartbeat engine
new authority state
Conductor→Solvent authority control loop
workflow SDK
new UI
```

That is the line we should hold.

### Final verdict: **AMBER → GREEN after a focused v0.3 revision**

I would **not move to Workflow Specification v0.1 yet**. The three reviews collectively uncovered enough genuinely important conformance gaps—especially consequential enforcement, observability, correlation, and ambiguous outcomes—that they should be resolved first.

But I also would **not restart Requirements from scratch**. The current v0.2 foundation is good. Add the narrowly defined corrections above, especially the four hard ones:

```text
1. Fail-closed effect boundary
2. Observable/correlatable cross-role evidence
3. Explicit domain ownership boundary
4. Explicit unknown/ambiguous execution outcome
```

Then add the D/E BM-IST scenarios and the lightweight replay/intervention/versioning semantics.

At that point, **Requirements v0.3 becomes a proper conformance baseline rather than merely a strong architectural statement**.
