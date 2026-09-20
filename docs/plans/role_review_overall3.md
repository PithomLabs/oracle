## Consolidated adversarial review — Role / Boundary Matrix v0.3

I consolidated the three new reviews with my prior review of v0.3.

### Final verdict

**GREEN- / two substantive blockers remain before the Conformance Test Matrix should be started.**

The macro-architecture is now stable. The remaining issues are no longer about introducing another runtime, authority engine, gateway, or orchestration layer. They are about making two parts of the contract mechanically precise:

1. **Operation-identity comparison semantics are still underspecified.**
2. **The declaration model still lacks a guaranteed owner for generic capabilities that have no distinct integration owner.**

The rest are secondary consistency/testability refinements.

This is materially different from the earlier rounds: there is **no longer a licensed consequential bypass** in the written model. The remaining problems are under-specification at the edges, not a fundamentally broken architecture. 

---

# 1. What v0.3 got right

All three reviews agree that the major v0.2 findings are closed.

The matrix now has:

* explicit external vs internal boundaries,
* non-effect vs effect-capable capability declarations,
* declaration lifecycle,
* fail-closed behavior including authority-verification failure,
* semantic/effect/enforcement ownership separation,
* SOR ownership of actual effect occurrence,
* explicit rejection/failure/termination/ambiguity outcomes,
* versioned declaration evidence,
* review disposition,
* and an explicit refusal to create a new Loop Engineering runtime.   

The prior Gemini review also confirms that the earlier v0.2 weaknesses — oxymoron naming, network-partition fail-open risk, terminated/revoked outcome, and declaration drift — were successfully addressed. 

So the architecture itself is **not reopened**.

---

# 2. HIGH — Operation identity defines *what*, but not *how equality is determined*

This is the strongest new finding, and I agree with it.

The matrix now requires the same operation-identity definition at proposal, authorization, and execution, and it defines semantic effect inputs versus transport metadata. 

But that still leaves a subtle gap:

```text
Identity definition:
    target + operation + amount

Comparison:
    ??? 
```

Two implementations could legitimately interpret the same declared identity differently:

```text
Executor A:
    canonical structured comparison

Executor B:
    serialized byte comparison

Executor C:
    case-insensitive string comparison
```

That can produce either false rejection or, worse, false acceptance.

Z identifies this as the missing final link in the binding chain. 

### Required fix

Add to the declaration:

> **operation-identity comparison rule**

For example:

> The declaration MUST define the equality/comparison semantics used to determine whether two operation identities are the same. The comparison rule MUST be applied consistently at proposal, authorization, and execution.

This does **not** mean mandating RFC 8785.

The review itself correctly recognizes that rejecting mandatory RFC 8785 was the right architectural decision. 

The contract should require the **security property**, not a canonicalization technology.

### Conformance test

Add:

```text
same semantic operation
different serialization / field ordering / representation
        ↓
same binding result
```

and the inverse:

```text
different semantic operation
different effect-relevant value
        ↓
binding MUST fail
```

This should become one of the core Test Matrix cases.

---

# 3. HIGH — generic capabilities still need a declaration owner

This is the other blocker.

The matrix correctly says every capability used as a declared boundary must have an owned declaration. 

But consider:

```text
shell
web fetch
code execution
generic HTTP client
```

There may be no dedicated "integration owner" in the deployment.

Z correctly observes that the rule can become ownerless in precisely the environments relevant to BM-IST. 

### Required fix

Add:

> Where a declared capability has no distinct integration owner, the deployment operator is the default declaration owner and bears responsibility for declaration completeness and lifecycle.

Then:

```text
native integration exists
    → integration owner

no distinct integration owner
    → deployment owner
```

That keeps the model generic without adding a seventh role.

This should also resolve the current `Agent/tool/service as applicable` owner ambiguity in the non-effect row.

---

# 4. Medium — distinguish authority-verification unavailability from cached authorization

This is subtle but important.

The matrix says:

> if required authority verification is unavailable, reject. 

The broader workflow model also allows pre-existing authorization validity to be exercised according to its declared validity model.

So the contract needs to distinguish:

```text
Solvent unavailable
AND
verification requires live Solvent
```

from:

```text
Solvent unavailable
BUT
authorization evidence is independently verifiable under its declared trust/validity model
```

The fix is one sentence:

> "Authority verification unavailable" means the required verification mechanism itself cannot establish authorization. Independently verifiable authorization evidence that remains valid under the declared trust model is not treated as unavailable merely because the authority service is unreachable.

Otherwise Scenario F can be interpreted in two different ways. Z correctly flags this boundary. 

---

# 5. Medium — one semantic owner per fact, not necessarily per row

The matrix now mostly handles ownership correctly, but there are still compound rows such as:

> Solvent for authority fact; Executor for execution fact. 

That is actually correct architecture.

The acceptance criterion should therefore be understood as:

> Every **fact** has one semantic/authority owner.

rather than:

> Every table row has one participant.

This is mostly a wording improvement, but it makes the model mathematically cleaner.

I would change criterion 1 to:

> Every authoritative fact has one semantic/authority owner.

That aligns with the existing `Execution vs External System of Record` distinction. 

---

# 6. Medium — carried review findings need explicit disposition

The v0.3 review-disposition section is good, but it primarily records **architectural proposals**:

* mandatory runtime → rejected
* Solvent registry → rejected
* gateway → rejected
* shared library → deferred/optional
* oracle → accepted limitation



Z correctly points out that several recurring *findings* remain outside this log. 

That matters because otherwise the next artifact author cannot distinguish:

```text
not addressed
```

from:

```text
deliberately rejected
```

### Fix

Add a compact subsection:

| Finding                            | Status                               | Rationale                                          |
| ---------------------------------- | ------------------------------------ | -------------------------------------------------- |
| Identity comparison semantics      | Adopt                                | Security-critical                                  |
| Generic capability owner           | Adopt                                | Prevent ownerless declarations                     |
| Concurrency ownership              | Deferred                             | Test Matrix concern unless implementation requires |
| Client-internal state              | Rejected as Loop concern             | Client responsibility                              |
| Production observation persistence | Deferred to conformance requirements | Not Loop-owned state                               |
| Cached-auth evidence marker        | Adopt/define in Test Matrix          | Needed for deterministic evidence                  |

No need to reopen the architecture for every carry.

---

# 7. Medium — terminated/revoked has two edge cases worth pinning down

The new `terminated / revoked` category is correct. 

Two details should be explicit.

### Halt outcome itself can be ambiguous

If:

```text
operation started
cancel requested
executor disappears
```

then we do **not** know whether termination succeeded.

Therefore:

```text
termination known successful
    → terminated/revoked

termination outcome unknown
    → ambiguous
```

That follows the existing ambiguity model and should be written down.

### Terminated does not imply zero external effect

An operation may have produced partial external effects before being halted.

The SOR remains authoritative for occurrence.

This is consistent with the matrix's current SOR model. 

---

# 8. Medium — `Effect-Capable Operation Invocation` and `Consequential Proposal`

Gemini argues these should be physically merged because §7 says they are sequential aspects of one boundary. 

I **do not recommend merging them**.

The current structure is useful because it separates:

```text
capability invocation
        ↓
proposal establishment
```

while §7 already states they are not alternative execution paths. 

The right fix is **not structural merging**.

Add one phrase to the boundary table:

> These rows represent sequential aspects of the same consequential interaction and MUST NOT be implemented as independent alternative paths.

That preserves the conceptual distinction without suggesting two different authority systems.

---

# 9. Medium — don't make "succeeded" imply SOR confirmation universally

Gemini raises a good distributed-systems concern:

> `202 Accepted` or queue handoff does not necessarily mean the external effect occurred. 

The deeper architecture already handles this by giving the external SOR authority over actual occurrence.

I would refine the outcome rule rather than ban `succeeded` unless the SOR is directly queried.

Add:

> An Executor MUST NOT report `succeeded` merely because an asynchronous handoff was accepted. `succeeded` requires evidence sufficient under the declared execution/outcome model to establish successful effect occurrence; otherwise the outcome remains `attempted` or `ambiguous`.

This is better than a blanket rule requiring synchronous SOR confirmation.

---

# 10. Medium — under-declaration is a conformance failure, but the generic harness cannot always discover it

This is a useful distinction from Gemini.

The matrix says omission of an effect-relevant parameter is a conformance failure. 

It also explicitly says Loop Engineering cannot provide a universal hidden-side-effect oracle. 

These are not contradictory, but the scope should be explicit:

> Incomplete operation-identity declaration is an integration/design-time conformance failure. The generic protocol harness verifies declared identity behavior; completeness of the declaration may require integration-specific tests, domain review, or deployment review.

That makes the division of labor precise.

Do not weaken the conformance obligation; clarify **who/what can establish it**.

---

# 11. Medium — Autonomous mode and indefinite UNKNOWN

Gemini argues that autonomous mode cannot use operator-controlled indefinite wait. 

I agree with the underlying principle, but I would not hard-code "autonomous = terminal" into the role matrix unless Workflow Specification v0.3 explicitly requires that.

The safer rule is:

> An operator-controlled indefinite wait is valid only when an active operator exists to control it. A deployment operating without such an operator MUST use a finite/remediation policy.

That preserves both:

```text
semi-autonomous + human hold
```

and:

```text
autonomous + no indefinite wait
```

without creating a new state model.

---

# 12. Medium — declaration version should explicitly flow through authorization evidence

The matrix already records declaration/identity version in authorization at the boundary level and in execution evidence.  

But the `Authorization Evidence` rule itself could be more mechanical.

Gemini's proposed fix is good:

> authorization evidence MUST expose the specific declaration/identity version under which it was evaluated. 

I would adopt that.

This is especially useful for queued and long-running operations.

---

# 13. Low — declaration requirements are duplicated

Sections 8 and 18 both describe declaration contents.

Gemini recommends consolidating them. 

I agree.

Keep §8 as the **canonical declaration contract/lifecycle definition**.

Make §18 say:

> Every declared capability boundary MUST satisfy the declaration requirements defined in §8.

Then §18 focuses on conformance/application rather than repeating the field list.

This reduces future drift.

---

# 14. The "zombie authorization" objection does not justify Conductor → Solvent revocation

Gemini3 strongly argues that Conductor cancellation can leave a still-valid authorization token. 

This is a **real system-design concern**, but its proposed architectural fix is wrong for this model.

Do **not** establish:

```text
Conductor cancellation
        ↓
Solvent revocation
```

as a universal rule.

That would make coordination state an implicit authority mutation mechanism, violating the ownership model.

Instead the architecture should require the consequential integration to declare the relationship between:

```text
work cancellation
authorization validity
execution cancellation/revocation
```

The authority owner remains Solvent.

A deployment that needs cancellation to revoke authority can explicitly invoke the Solvent revocation mechanism as an **authority operation**, initiated by an actor with authority to do so. It should not happen because the Conductor lifecycle itself secretly changed Solvent state.

That distinction matters:

```text
Conductor says: task cancelled
```

does not imply:

```text
Solvent says: authorization revoked
```

unless an actual authority operation occurs.

So this remains **implementation/integration policy**, not a Loop Engineering primitive.

---

# 15. The SOR reconciliation objection also does not justify a Loop runtime

Gemini3 argues that a crashed Executor can leave an ambiguous outcome without a daemon to reconcile it. 

Again, the operational concern is legitimate.

The architecture already requires a declared reconciliation owner/mechanism and persistence/recovery expectations. 

The correct requirement is:

> Every effect-capable integration whose operation can produce an ambiguous outcome MUST declare a viable reconciliation mechanism and recovery ownership.

That mechanism can be:

* an external provider status API,
* a durable job record,
* Executor persistence,
* a retry/reconciliation process,
* operator review,
* or another domain-appropriate mechanism.

It does **not** become:

```text
Loop Engineering Reconciliation Runtime
```

This distinction is central to keeping the architecture minimal.

---

# 16. Agent hallucination does not require a protocol authority layer

Gemini3 argues that an LLM can still misinterpret `FAILED` as success and propose another consequential action. 

That's true.

But the matrix correctly prevents the Agent from rewriting the authoritative outcome. 

The protocol cannot prevent an agent from **thinking incorrectly**.

What it can ensure is:

```text
Executor says FAILURE
        ↓
Agent may misunderstand it
        ↓
Agent proposes new action
        ↓
new action must independently satisfy consequential authorization
```

That is the correct security boundary.

Do not add an infrastructure layer attempting to police Agent cognition.

---

# 17. Security SDK / gateway proposals remain optional implementation artifacts

Gemini3 again recommends a shared verification library because many adapters could otherwise diverge. 

This is a **very reasonable engineering recommendation** but should remain exactly where v0.3 puts it:

> shared verification library — deferred/optional.



I would preserve that decision.

After the POC, empirical integration experience may justify:

```text
Pithom Loop Verification Library
```

as a reusable implementation package.

That would be an **ecosystem convenience**, not a new normative runtime or authority component.

---

# Consolidated patch list

I would make the following final edits before declaring v0.3 frozen:

### P0 — required

**1. Add operation-identity comparison semantics**

Declaration must specify not just identity fields but how equality is determined.

**2. Define default declaration owner**

Where no distinct integration owner exists, deployment operator becomes the declaration owner.

These are the two genuine remaining blockers.

### P1 — strongly recommended

**3. Distinguish live authority-verification failure from independently verifiable cached evidence.**

**4. State that every authoritative fact, rather than every table row, has one semantic owner.**

**5. Make `terminated` vs `ambiguous termination` explicit.**

**6. Clarify that `terminated/revoked` does not imply zero external effect.**

**7. Require `succeeded` to be supported by evidence sufficient under the declared outcome model, not merely a handoff receipt.**

**8. Make declaration/identity version part of authorization evidence.**

### P2 — cleanup

**9. Clarify that under-declaration is a design/integration conformance failure and generic harnesses cannot universally discover hidden side effects.**

**10. Require no indefinite UNKNOWN hold without an actual operator.**

**11. Deduplicate capability-declaration field lists between lifecycle and application sections.**

**12. Extend review disposition with recurring finding status, not only architectural proposal decisions.**

---

# What I would explicitly *not* change

After six rounds, this is now important.

Do **not** respond to these reviews by introducing:

```text
mandatory workflow runtime
mandatory reconciliation daemon
mandatory Operations Gateway
central Solvent effect registry
Conductor→Solvent automatic revocation
mandatory SDK
mandatory RFC 8785
new workflow state engine
new authority service
```

Those would solve implementation problems by changing the architecture rather than improving the contract.

The review dispositions already correctly reject those proposals. 

---

# Final consolidated verdict

```text
Macro architecture:                 GREEN
Ownership model:                    GREEN
Authority/effect separation:        GREEN
Ordinary-work bypass:               CLOSED
Declaration lifecycle:              GREEN
Fail-closed boundary:               GREEN
Execution outcome model:            GREEN
SOR model:                           GREEN
Human intervention:                 GREEN
Runtime neutrality:                 GREEN

Remaining:
  Operation identity equality:      HIGH — FIX
  Generic capability owner:         HIGH — FIX

Everything else:                    P1/P2 cleanup
```

The architecture has **converged**.

The right next move is not another broad redesign. Apply the two P0 fixes, fold in the small P1/P2 clarifications, then **freeze the Role / Boundary Matrix**.

The resulting artifact should be treated as the boundary contract that the Conformance Test Matrix must test—not as a document to keep expanding.

One caveat remains from the Z review: these reviews are assessing the matrix's relationship to **Workflow Specification v0.3**, but the parent specification itself has not been supplied in these review files.  That means the parent-level consistency judgment remains necessarily provisional; the matrix's own internal consistency, however, is now strong enough to proceed once the two P0 items land.
