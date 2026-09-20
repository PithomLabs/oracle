## Consolidated adversarial review — Role / Boundary Matrix v0.2

After combining the three new reviews with my prior review of v0.2, the disposition is:

**GREEN- / not quite ready to freeze.**

The important architectural holes from v0.1 are genuinely closed. The remaining findings are mostly **consistency and lifecycle precision**, but two are important enough to fix before the Conformance Test Matrix is derived:

1. **The non-effect declaration class is underspecified.**
2. **The matrix's own tables and acceptance criteria are no longer perfectly aligned on ownership/effect ownership.**

The newer reviews also identify several worthwhile micro-patches. I would incorporate most of those, but **not** their proposals to introduce a mandatory verification SDK, Conductor→Solvent cancellation revocation, or Solvent-owned approval of every immaterial-parameter list. Those would push the architecture back toward centralized infrastructure.

---

# 1. What all reviews agree v0.2 successfully fixed

The major v0.1 problems are closed.

### Consequential bypass

The two-row ordinary-work split plus the Boundary Rule closes the original bypass:

> Any operation capable of producing an external effect is consequential regardless of workflow phase or invoking client. 

That was the most important correction.

### Integration is not a new role

The terminology is now explicit:

> Integration is an implementation boundary or adapter belonging to an existing role boundary. 

All three reviews accept this architectural move. 

### Trust basis is properly constrained

The evidence must be independently verifiable and cannot merely be an Agent/client assertion. The mechanism remains implementation-neutral. 

### UNKNOWN is no longer an undefined wait

The matrix requires an actionable maximum unresolved-wait/age policy and deliberately does not invent a new global `Abandoned` or escalation state. 

### Long-running precedence is correct

Solvent's authority constraints are authoritative; integrations may tighten but cannot weaken them. 

### Rejection vs failure is now explicit

The distinction between:

```text
rejected before effect
```

and:

```text
attempted, then failed
```

is now first-class and testable. 

### Evidence coverage materially improved

Proposal version, human intervention, declaration, rejection, long-running progress, and reconciliation now have explicit evidence rows. 

So the macro-architecture is sound.

---

# 2. Critical remaining issue: the non-effect declaration class

This is the strongest new finding.

The matrix now depends on:

> `Declared non-effect tool/service`

for the ordinary path. 

But the formal declaration requirements in §16 are written specifically for **effect-capable integrations**. 

That leaves an awkward lifecycle:

```text
Tool initially declared non-effect
        ↓
tool gains a mutating operation
        ↓
what forces re-declaration?
```

The Boundary Rule itself says the operation must be consequential, but the conformance lifecycle does not yet have an explicit artifact/change obligation for the non-effect class.

Qwen/Z correctly identifies this as the remaining foundation-level ambiguity. 

### Fix

Extend the declaration lifecycle to **all declared capability boundaries**, not only effect-capable ones.

For example:

> Any capability used as a declared Loop Engineering boundary MUST have an owned declaration identifying its effect classification, owner, version, scope, and effective reference. A change that can alter effect classification or operation identity invalidates prior conformance for the affected boundary.

Then:

```text
Non-effect declaration
    ↓
changes to effect-capable
    ↓
declaration version changes
    ↓
re-conformance required
```

This is better than creating a separate elaborate "non-effect registry."

---

# 3. Ordinary Work — effect-capable should be renamed

This is a good micro-level finding from Qwen. 

Even though the Boundary Rule makes the semantics correct, the phrase:

> `Ordinary Work — effect-capable capability`

is internally awkward.

An engineer scanning the matrix could interpret "Ordinary Work" as the bypass path.

Rename it to:

> **Effect-Capable Operation Invocation**

Then make the rule explicit:

> Any invocation of an effect-capable capability is a consequential action and MUST satisfy the consequential boundary rules, regardless of the workflow phase from which it originates.

This also makes the relationship to `Consequential Proposal` easier to understand.

---

# 4. The relationship between "Effect-Capable Operation Invocation" and "Consequential Proposal" needs one explicit sentence

Z correctly notes that these rows now overlap. 

Today the matrix has:

```text
Effect-Capable Operation Invocation
Consequential Proposal
Authorization
Execution Eligibility
Execution
```

An implementer could wonder whether these represent two alternative paths.

They should not.

The intended relationship is:

```text
Effect-Capable Operation Invocation
        ↓
Consequential Proposal
        ↓
Authorization
        ↓
Execution Eligibility
        ↓
Execution
```

So add:

> An Effect-Capable Operation Invocation MUST enter the applicable Consequential Proposal and Authorization boundaries before effect.

That closes the ambiguity without creating another state.

---

# 5. Ownership refactoring introduced a new internal inconsistency

This is probably the most important structural cleanup remaining.

The matrix correctly replaced slash-separated owners with:

> Semantic / authority owner
> Enforcement / implementation responsibility



But the acceptance criteria still says every consequential boundary must identify:

> authority owner **and effect owner**



Those are no longer represented consistently in the table.

Z identifies this as a direct criterion violation. 

### Fix

Do not revert to slash-duals.

Instead add explicit conceptual columns:

```text
Semantic / authority owner
Effect / outcome owner
Enforcement / implementation responsibility
```

For example:

| Boundary               | Authority owner                  | Effect/outcome owner    | Enforcement          |
| ---------------------- | -------------------------------- | ----------------------- | -------------------- |
| Consequential Proposal | Solvent                          | Pending                 | Client/integration   |
| Authorization          | Solvent                          | Pending                 | Integration          |
| Execution Eligibility  | Solvent                          | Executor                | Integration/Executor |
| Execution              | Solvent                          | Executor / external SOR | Executor/SOR         |
| Result                 | Solvent for prior authority fact | Executor / external SOR | Executor/integration |

That preserves the one-owner principle while satisfying the document's own criterion.

---

# 6. The Failure Ownership Matrix still contains slash-dual ownership

Same issue, different table.

The matrix has entries such as:

```text
Conductor/integration
Solvent/integration
Effect-capable integration / Executor boundary
Executor/integration
```



That partially undoes §8.

The semantic owner should be singular.

For example:

```text
Authorization UNKNOWN → Solvent
Execution rejection  → Effect-capable boundary
Execution ambiguous  → Executor / external SOR
```

with enforcement/reporting responsibility moved into a separate column or explanatory text.

This is exactly the sort of consistency issue that becomes painful once the Conformance Test Matrix starts citing these rows.

---

# 7. Formulate is simultaneously a boundary and not a boundary

This is a real editorial inconsistency.

§5 contains:

> Formulate

as a boundary row. 

§6 then says:

> Formulate is not an external participant boundary. 

Qwen/Z correctly flags that this is contradictory. 

### Fix

Treat Formulate as an **internal phase**, just as §6 says.

Therefore remove Formulate from the external Boundary Matrix and keep it in the internal-phase section.

Likewise, the intervention matrix should not treat it as an external boundary.

---

# 8. Next Work needs an intervention entry

The opposite inconsistency also exists.

`Next Work` appears in the boundary matrix with human intervention capability, but it does not have a corresponding dedicated intervention row.  

Add:

```text
Next Work
    inspect: current work/review context
    intervene: redirect/revise/conclude where supported
    limit: Conductor remains coordination owner
```

Simple consistency fix.

---

# 9. Under-declaration should be an explicit conformance failure

This is a strong security observation from Qwen.

The identity model says effect-relevant parameters must be included or explicitly declared immaterial. 

But there is not yet a direct statement about **incomplete declaration**.

Add:

> The integration owner is responsible for the completeness of the operation-identity definition. Omitting an effect-relevant parameter is a conformance failure.

I would **not** call this "strict liability" in a legal sense. Architecture specifications do not need that terminology.

The important thing is:

```text
missing parameter
      ↓
conformance failure
      ↓
cannot claim a valid operation boundary
```

---

# 10. Network failure should explicitly fail closed

This is a legitimate micro-level improvement.

The matrix permits direct Solvent lookup as a trust mechanism. 

Therefore:

```text
Executor
   ↓
verify with Solvent
   ↓
network partition
```

must not become:

```text
"probably still authorized"
```

Add:

> If required authority verification is unavailable, times out, or cannot establish the required trust basis, the evidence MUST be treated as unavailable/invalid and the external effect MUST be rejected.

That makes fail-closed behavior explicit.

---

# 11. "Terminated / revoked" is worth adding

Qwen identifies a legitimate outcome-model gap. 

The matrix already supports mid-flight revocation. 

So an operation that:

```text
started
→ was explicitly halted/revoked
→ halt was known to succeed
```

should not be mislabeled `failed`.

Add:

```text
terminated / revoked
```

with:

> The effect was attempted and was explicitly halted or revoked through an authorized intervention. It is distinct from system failure.

This is a useful test category.

---

# 12. Declaration version should be bound to the applicable execution contract

This is a good finding, but the prescription needs some care.

The declaration is already observable and versioned.  

The missing rule is:

```text
Which declaration semantics apply to this authorization?
```

The clean solution is not necessarily a cryptographic binding.

Add:

> Authorization and execution MUST identify the declaration/operation-identity version under which the authorization was established. The Executor MUST validate the operation against that applicable version.

That prevents:

```text
authorize under v1
execute under silently substituted v2
```

without mandating a particular transport mechanism.

---

# 13. Operator-controlled indefinite UNKNOWN needs evidence

This is subtle and correct.

The matrix says indefinite operator-controlled wait is allowed. 

But an operator hold with no recorded evidence becomes:

```text
UNKNOWN
... forever
```

Add:

> An operator-controlled unresolved hold is a human intervention and MUST produce intervention evidence.

That makes the exception observable.

---

# 14. Effect-detection limitation should remain, but ownership should be explicit

The current declaration is honest:

> Loop Engineering does not provide a universal oracle capable of discovering undeclared side effects. 

Good.

But the acceptance criteria say completeness is an integration/deployment responsibility without naming the owner. Z calls this out. 

Add:

> The integration owner is responsible for completeness of its declared effect surface; deployment/conformance ownership is responsible for ensuring the declaration is the one under test.

Don't put this responsibility into Solvent. Solvent should not become a global capability registry.

---

# 15. Executor vs external system of record should be separated cleanly

This remains somewhat conditional on the exact Workflow Specification v0.3 language, but architecturally the clean model is:

```text
External system of record
    = whether the effect actually occurred

Executor
    = attempted operation + execution report/reconciliation
```

The matrix already hints at this with "subject to external system-of-record semantics." 

That distinction should be made explicit in the role model and outcome table.

---

# 16. Do not adopt these three proposed "hardening" measures

Gemini2 proposes:

> standardized verification SDK / RFC 8785 canonicalization
> Conductor cancellation → Solvent revocation
> Solvent approval of immaterial parameter lists



I would **reject all three as architecture requirements**.

### Standard verification SDK

A shared library could be a perfectly good implementation later.

But:

```text
Loop Engineering contract
        ≠
mandatory SDK
```

The conformance harness should verify the security property, not mandate a particular implementation artifact.

RFC 8785 is likewise unnecessary at this specification layer.

### Conductor-triggered authority revocation

This would violate the established semantic ownership boundary.

Conductor cancellation means:

```text
coordination changed
```

not:

```text
authority changed
```

Automatic coupling would make Conductor a hidden authority mutator.

The correct relationship remains explicit:

> Conductor cancellation does not itself mutate Solvent authority.

Any revocation must occur through the applicable authority boundary.

### Central Solvent approval of "immaterial" fields

This would move Solvent toward becoming an effect-schema registry.

The architecture should instead impose:

```text
integration owner
    → accountable for operation identity definition
    → conformance verifies it
```

Solvent remains the authority engine, not global semantic ownership of every domain/tool payload.

---

# 17. One process issue is emerging: repeated silent carries

Z highlights a process-level problem that is now worth addressing. 

Some findings have appeared across several review rounds without explicit:

```text
adopted
rejected
deferred
```

status.

Before the Conformance Test Matrix, create a tiny **Review Disposition** section or appendix:

| Finding                         | Disposition         | Reason                              |
| ------------------------------- | ------------------- | ----------------------------------- |
| Mandatory workflow runtime      | Rejected            | Violates architecture               |
| Central Solvent effect registry | Rejected            | Expands authority kernel            |
| Effect oracle                   | Accepted limitation | Impossible to guarantee generically |
| Shared verification library     | Deferred/optional   | Implementation mechanism            |
| Conductor→Solvent cancellation  | Rejected            | Cross-owner authority mutation      |

This is not another architecture layer. It is documentation hygiene.

It will prevent the next reviews from repeatedly rediscovering decisions that have already been made.

---

# Consolidated patch order

I would now implement **these 10 patches**:

| Priority | Patch                                                                               |
| -------- | ----------------------------------------------------------------------------------- |
| **P0**   | Define lifecycle/ownership for **non-effect declarations**                          |
| **P0**   | Reconcile semantic/effect owner columns and failure-table ownership                 |
| **P0**   | Clarify parent-spec rule: matrix may **refine/tighten**, not contradict             |
| **P1**   | Rename `Ordinary Work — effect-capable`                                             |
| **P1**   | Remove Formulate from external boundary matrix; add/complete Next Work intervention |
| **P1**   | Explicitly require consequential invocation → proposal → authorization              |
| **P1**   | Add under-declaration as conformance failure                                        |
| **P1**   | Explicit network-failure fail-closed rule                                           |
| **P1**   | Add `terminated / revoked` outcome                                                  |
| **P1**   | Bind authorization/execution to declaration/identity version                        |
| **P2**   | Evidence for operator-controlled UNKNOWN wait                                       |
| **P2**   | Explicit effect-surface completeness owner                                          |
| **P2**   | Clarify Executor vs external SOR                                                    |
| **P2**   | Add review-disposition record                                                       |

---

# Final consolidated verdict

### Macro architecture

**GREEN.**

The central model is holding:

```text
AGENT       = AGENCY
CONDUCTOR   = COORDINATION
SOLVENT     = AUTHORITY
EXECUTOR    = EFFECT
DOMAIN      = MEANING / TRUTH
HUMAN       = INTERVENTION / REVIEW
```

and:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

The v0.2 fixes genuinely closed the major v0.1 bypasses. 

### Remaining state

**GREEN- / patch before freeze.**

The biggest remaining issue is no longer "architecture is wrong." It is:

> **the projection must be internally self-consistent before it becomes a test oracle.**

The highest-priority corrections are therefore not new infrastructure. They are **declaration lifecycle, ownership consistency, boundary classification, and version applicability**.

Once those are patched, I would **freeze Role / Boundary Matrix v0.2 as v0.3 and stop broad architectural adversarial review.**

Then the next artifact should be the **Conformance Test Matrix**, with every remaining rule converted into executable positive/negative/substitution tests. The current §24 questions are already an excellent test inventory. 

The important threshold has been crossed: further useful work is now much more likely to come from **testing the contract against real integrations** than from adding more conceptual machinery.
