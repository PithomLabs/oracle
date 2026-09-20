## Consolidated adversarial review — Loop Engineering Role / Boundary Matrix v0.4

I consolidated both attached v0.4 reviews with my prior v0.3 review and the accumulated findings from the earlier rounds.

### Verdict

**YELLOW-GREEN / not lockable yet.**

The reviews agree that v0.4 closes several real v0.3 findings, but they also identify **two mechanical security gaps that should be treated as P0**, plus several structural clarifications before the Conformance Test Matrix begins.

The important distinction is that **the macro-architecture remains sound**. The remaining work is primarily to make the protocol mechanically unambiguous across independent implementations and to clarify how coordination observes the consequential branch.

The biggest issue is now unmistakable:

```text
operation identity
     ↓
operation equality/comparison
     ↓
authorization evidence
     ↓
execution
```

The declaration has to control every link in that chain.

---

# 1. What v0.4 successfully closed

The reviews agree that the revision materially improved the model.

Gemini notes that v0.4 incorporated the major v0.3 fixes: asynchronous handoff is no longer automatically treated as success, terminated/revoked is explicit, network verification failure fails closed, and review disposition exists. 

That means the architecture still preserves:

```text
Agent       → agency
Conductor   → coordination
Solvent     → authority
Executor    → effect/reporting
SOR         → actual effect occurrence where applicable
Domain      → meaning/truth
```

and still avoids creating a Loop Engineering runtime.

The hostile Gemini review's proposed remedies — mandatory live Solvent checks, mandatory SOR callback infrastructure, or automatic Conductor-triggered revocation — should **not** automatically be accepted. 

Those are implementation/policy choices, not demonstrated reasons to change the fundamental architecture.

---

# 2. P0 — Authorization evidence still does not explicitly carry declaration version

This is the clearest mechanical defect.

The review identifies the inconsistency:

* authorization evidence contains operation identity/context,
* execution evidence contains declaration/identity version,
* but the authorization evidence itself does not clearly expose the version used for authorization. 

That creates exactly the situation the declaration-drift rules are supposed to prevent:

```text id="yqk5x7"
Authorization:
    operation X
    declaration v1

Execution:
    obtains v2 independently
    validates X under v2
```

That is not a sufficiently closed binding.

### Required fix

Authorization evidence must contain or verifiably expose:

```text id="5kxjlf"
authorized operation identity
+
declaration / identity version
```

The Evidence Matrix authorization row should contain the same version.

This should be **P0**.

---

# 3. P0 — Operation identity needs cross-implementation equality semantics

This is the other genuine security issue.

The matrix apparently now requires deterministic identity comparison, but deterministic **within one implementation** is not enough. The two reviewers correctly distinguish:

```text id="w5gnyj"
deterministic
```

from:

```text id="s9z3jk"
interoperably equivalent across implementations
```



For example:

```text
Executor A:
{"amount":100,"target":"ABC"}

Executor B:
{"target":"ABC","amount":100}
```

or:

```text
"ABC"
vs
"abc"
```

Depending on the declared semantics, those might be:

```text
same operation
```

or:

```text
different operation
```

The specification cannot leave that to each Executor.

### Required fix

The declaration must specify the **comparison/equality rule**, not merely the identity fields.

Suggested rule:

> The operation-identity definition MUST include or reference the comparison rule used to determine equality. All participants evaluating operation binding for the same declared operation MUST apply the same declared comparison rule.

This does **not** require RFC 8785.

It requires the property:

```text
same declaration
+
same semantic inputs
+
same comparison rule
=
same binding result
```

That is the important invariant.

---

# 4. The third major issue: generic capability ownership

One prior concern that remains relevant is the ownerless capability problem.

v0.3 introduced the deployment operator as default owner for a capability without a distinct integration owner. The new review challenges whether this simply transfers semantic responsibility to operations personnel who may not understand the tool's effect surface. 

I agree that this requires a better formulation, but I would **not remove the default owner**.

Without a default owner:

```text
generic shell/web/HTTP/code tool
    ↓
no owner
    ↓
no declaration accountability
```

is worse.

Instead:

> Where no distinct integration owner exists, the deployment operator is the accountable declaration owner for deployment-level conformance, while the tool/provider owner remains responsible for the correctness of any published capability semantics it supplies.

This splits:

```text
deployment accountability
vs
software semantic expertise
```

without creating another role.

The declaration owner has accountability for the boundary; that does not mean the operator must personally reverse-engineer arbitrary software.

---

# 5. High — the two-row consequential boundary issue

The two reviews disagree slightly with my prior position here.

Gemini recommends merging:

```text
Effect-Capable Operation Invocation
Consequential Proposal
```

because §7 says they are sequential aspects of one boundary. 

I still would **not merge them**, because they describe two meaningful interaction moments:

```text
capability invocation
      ↓
proposal establishment
```

But the current terminology apparently creates a genuine matrix/prose contradiction.

So the right fix is:

### Keep both rows, explicitly define them as

> two sequential mechanical steps within one logical consequential boundary.

Then remove any wording that implies they are alternative boundaries.

That preserves the useful decomposition without causing implementers to infer two independent workflow protocols.

---

# 6. High — Executor-as-SOR needs explicit semantics

This is a good new finding.

The current model recognizes:

```text
Executor
+
external SOR
```

but not sufficiently the case:

```text
Executor == system of record
```

Examples include:

```text
direct database mutation
synchronous filesystem operation
synchronous API returning authoritative resource state
```

Gemini correctly flags this. 

### Required clarification

Add:

> Where no distinct external system of record exists, the Executor is authoritative for actual effect occurrence and may report `succeeded` when its own authoritative execution evidence establishes successful effect occurrence.

That gives a clean model:

```text
External SOR exists
    → SOR owns actual occurrence

No external SOR
    → Executor owns actual occurrence
```

This also makes the `succeeded` rule much cleaner.

---

# 7. High — Conductor disappears from the consequential branch

This is probably the most important **coordination**, rather than security, finding.

The external boundary sequence is apparently:

```text
Agent
→ Integration
→ Solvent
→ Executor
→ Result
```

with Conductor absent.

Yet Conductor owns task lifecycle.

The review correctly asks:

> how does Conductor know a task is awaiting authorization or long-running execution?



This does **not** mean Conductor needs a new authorization state.

The correct fix is an observational/coordination statement:

> Conductor MAY record coordination-relevant facts about a consequential operation, including that work is awaiting an external authorization or execution result, using its existing lifecycle/activity model. Such records are coordination facts and do not constitute authority decisions.

That gives:

```text
Conductor knows:
    work is waiting

Solvent knows:
    whether it is authorized

Executor knows:
    whether it happened
```

No owner is stolen.

This also supports Scenario G's non-blocking test.

---

# 8. High — Result routing is ambiguous

The current `Agent/client / Conductor` receiver is too loose if the specification is intended to be normative.

Gemini correctly identifies the ambiguity. 

The clean statement is:

> The execution result MUST be available to both the Agent/client for interpretation and Conductor for coordination/lifecycle observation, through applicable interfaces. Initial delivery may be direct or mediated.

That keeps transport flexible while making the required observability invariant explicit.

---

# 9. Medium — "succeeded" must mean confirmed effect, not accepted handoff

This is a good refinement from Gemini. 

The correct rule is not simply:

```text
202 ≠ succeeded
```

because an operation may legitimately have an intermediate handoff model.

Instead:

> `succeeded` requires evidence sufficient under the declared execution/outcome model to establish actual effect occurrence. A queue acceptance, message handoff, or equivalent intermediate acknowledgment alone is insufficient unless that artifact is itself authoritative for the declared operation model.

Therefore:

```text
202 Accepted
    → attempted / pending / ambiguous depending on declaration

authoritative SOR confirmation
    → succeeded
```

and:

```text
Executor is authoritative SOR
    → its authoritative confirmation can establish succeeded
```

That is much more general.

---

# 10. Medium — generic harness cannot detect omitted parameters universally

The reviewers are right that:

> "omission is a conformance failure"

and:

> "the generic harness cannot discover hidden undeclared effects"

can coexist, but the scope must be stated.



Add:

> Completeness of an operation-identity declaration is an integration/design-time conformance responsibility. The generic protocol harness verifies that declared identity inputs are correctly bound; it does not provide a universal oracle for discovering hidden effect-relevant inputs.

This is an important epistemic boundary.

Do not weaken the requirement that the declaration be complete.

---

# 11. Medium — Autonomous mode cannot have a fake "operator-controlled indefinite wait"

This concern is valid.

Merely recording human-intervention evidence does not create a human.

The rule should be:

> An operator-controlled indefinite wait is valid only in an operating mode in which an active human operator can act on the unresolved hold. Autonomous operation MUST use a terminal or remediation policy.



No new lifecycle state is required.

---

# 12. Medium — cached-auth evidence disposition is currently dangling

The Review Disposition claims that a cached-auth evidence marker was adopted, but the normative sections apparently do not actually contain it. 

This is a documentation integrity issue.

Either:

```text id="t0x8k2"
actually add the marker
```

to the authority evidence requirements,

or:

```text id="4z7u2o"
change disposition to deferred
```

Do not leave "adopted" without a normative location.

---

# 13. Medium — Sections 8 and 18 duplicate the declaration contract

This finding has now survived another round. 

That should simply be cleaned up.

One section should be canonical.

Recommended:

```text
§8 = canonical capability declaration + lifecycle
§18 = conformance application/reference to §8
```

This eliminates specification drift.

---

# 14. Zombie authorization: legitimate concern, wrong proposed solution

Gemini4 again raises the concern that Conductor cancellation does not automatically revoke Solvent authority. 

This is a real **deployment policy** issue.

It is not sufficient reason to violate:

```text
coordination ≠ authority
```

The architecture should instead make one thing explicit:

> A coordination cancellation does not itself revoke authority. A deployment MAY define an explicit authority-revocation action that is separately invoked through Solvent's authority boundary.

So:

```text
Conductor:
    cancelled

does not imply:

Solvent:
    revoked
```

but a policy can intentionally perform:

```text
authorized revocation operation
        ↓
Solvent
```

That is still clean.

The important thing is that **Conductor is not silently granted authority to mutate Solvent**.

---

# 15. Asynchronous reconciliation does not justify a Loop runtime

Gemini4 argues that strict outcome semantics push state-machine work into every adapter. 

That is a valid engineering tradeoff, but the architecture correctly keeps reconciliation integration-owned.

The declaration already requires:

```text
reconciliation owner
reconciliation mechanism
persistence/recovery expectation
```

So the next implementation artifact should demonstrate at least one realistic long-running integration.

Do not add a Loop daemon merely to centralize provider-specific state.

A provider may offer:

```text
callback
polling API
durable job ID
webhook
database state
operator reconciliation
```

Those are all acceptable mechanisms.

---

# 16. Deployment operator as declaration owner is not ideal, but still preferable

Gemini4 calls this an accountability trap. 

I would refine rather than remove it.

The declaration owner should be defined as:

```text
responsible for declaring and maintaining
```

not:

```text
person who must understand every implementation detail
```

The provider/tool may supply declaration semantics; deployment ownership verifies what actually exists in the deployed boundary.

This is analogous to configuration ownership, not authorship of software semantics.

---

# Consolidated priority order

### P0 — must fix

**1. Authorization evidence must carry/expose declaration + identity version.**

**2. Operation-identity comparison semantics must be explicitly declared and shared by all participants.**

Those are the genuine security-critical mechanical gaps.

### P1 — fix before lock

**3. Define Executor-as-SOR.**

**4. Clarify Conductor's coordination observation during consequential work without giving it authority.**

**5. Clarify Result availability to both Agent and Conductor.**

**6. Clarify two consequential rows as two mechanical steps within one logical boundary.**

**7. Scope parameter-omission conformance to integration/design-time review plus targeted tests.**

**8. Explicitly ban operator-controlled indefinite UNKNOWN in autonomous mode.**

### P2 — cleanup

**9. Resolve cached-auth disposition.**

**10. Deduplicate declaration requirements.**

**11. Clarify deployment operator as default declaration owner.**

**12. Clarify authoritative `succeeded` for Executor-as-SOR / asynchronous handoff models.**

**13. Extend disposition log to carried findings.**

---

# What remains rejected

I would **not** adopt:

```text
mandatory live Solvent call for every high-impact action
Conductor → Solvent automatic revocation
central Loop reconciliation daemon
mandatory SOR callback infrastructure
mandatory verification SDK
central effect registry
mandatory RFC 8785
```

The reviews demonstrate real implementation tradeoffs, but they do not establish that these require a new Loop Engineering architectural component. The Review Disposition already appropriately treats gateways/libraries as optional implementation techniques rather than Loop roles. 

---

# Final consolidated verdict

| Area                                   | Status           |
| -------------------------------------- | ---------------- |
| Role separation                        | **GREEN**        |
| Authority/effect ownership             | **GREEN**        |
| Consequential boundary                 | **GREEN-**       |
| Declaration lifecycle                  | **GREEN**        |
| Fail-closed enforcement                | **GREEN**        |
| SOR model                              | **GREEN-**       |
| Long-running execution                 | **GREEN-**       |
| Human intervention                     | **GREEN**        |
| Operation identity                     | **YELLOW — P0**  |
| Authorization evidence/version binding | **YELLOW — P0**  |
| Conductor visibility                   | **YELLOW — P1**  |
| Result routing                         | **YELLOW — P1**  |
| Documentation consistency              | **YELLOW-GREEN** |
| Architectural scope discipline         | **GREEN**        |

## Bottom line

**Do not lock v0.4 yet.**

But this is no longer another broad architecture cycle. The architecture has converged; the remaining work is a small set of **mechanical protocol clarifications**.

The two P0 fixes are especially important because they turn:

```text
"same operation identity"
```

into an actually enforceable interoperability/security contract:

```text
identity definition
      +
comparison semantics
      +
declaration version
      +
authorization evidence
      +
execution verification
```

Once those are explicit, I would make the P1/P2 consistency edits and **freeze the Role / Boundary Matrix**.

At that point the Conformance Test Matrix should become the primary adversarial instrument rather than continuing to enlarge the Role / Boundary Matrix.
