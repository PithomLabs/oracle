# Physics Verifier + Solvent — Architecture Discussion So Far

## 1. Goal

- **Objective:** Explore a Solvent demo in which AI swarm agents work on physics verification while Solvent provides the durable authority and evidence discipline around consequential decisions.
- **Core idea:** AI agents can reason, propose, attack, verify, and generate evidence, but they should not silently become the authority.
- **Important boundary:** Physics/math meaning stays outside the Solvent kernel. Solvent should enforce generic structural facts without understanding the content of the physics.

---

## 2. The Basic Mental Model

- **Agent/model:** Does reasoning, planning, mathematical exploration, proof search, counterexample search, and proposal generation.
- **Harness/runtime:** Controls agents, tools, credentials, files, network, sandbox, identity, and activity.
- **Policy/verifier layer:** Decides domain-specific semantic questions such as whether a claim is applicable, whether an attestation is required, and whether a review obligation is human-gated.
- **Solvent kernel:** Enforces generic durable facts such as belief identity, current promotion state, debt lifecycle, authority lifecycle, exact authority binding, and atomic transitions.
- **CockroachDB:** Provides the durable database invariants and concurrency guarantees.
- **Executor:** Performs the actual consequential external action only after the Solvent checks succeed.

Conceptually:

```text
Physics / domain semantics
        ↓
Policy / verifier / human review
        ↓
Solvent kernel
        ↓
CockroachDB invariants
        ↓
Execution
```

Important nuance: this is primarily an **authority flow**, not necessarily a literal dependency stack.

---

## 3. Single Canonical Mapping from Physics DAG to Solvent

The methodological rule is:

> Every distinct DAG concept should have exactly one canonical Solvent representation.

This is intended to prevent agents from choosing different valid-looking representations for the same concept and silently losing meaning.

### Mapping

1. **Node / claim**
   - One `Belief`, scoped to a `Scenario`.
   - Axiom, definition, result, or hypothesis is still represented as a belief rather than separate kernel tables.

2. **Epistemic kind**
   - `Belief.claim_type`.
   - The design discussed a many-to-few mapping such as:
     - axiom → Postulated
     - externally established result → Accommodated
     - internally proved result → Derived

3. **Mathematical / verification state**
   - Current Solvent state is represented through belief status plus debt rather than adding a separate `proved` state.
   - Roughly:
     - open → `entered` with unresolved debt
     - proved → `promoted`
     - refuted → `retracted`

4. **Review obligations**
   - Represented as debt / obligations.
   - Important correction: the kernel should not own a domain-specific vocabulary of what each debt *means*.

5. **Supporting artifacts**
   - `Evidence` rows linked to beliefs.
   - Evidence does not automatically discharge debt.

6. **Downstream consequential reliance**
   - `Action Intent` cites a belief.
   - The kernel checks whether that cited belief is currently promoted.

7. **Version revision**
   - Retract old belief → enter new belief → link with a `supersedes` relationship.
   - The old version remains readable and immutable.

8. **Staleness**
   - Prefer a computed query over stored staleness flags where the state can already be derived from belief relationships.
   - This avoids duplicated state that can drift.

9. **Separate research programs**
   - Use separate `Scenario`s.

10. **Actor capabilities**
    - Explicitly define which actors/tools can enter beliefs, add evidence, retire obligations, promote, etc.
    - Do not rely on descriptions or convention alone.

---

## 4. The Restricted-Belief Problem

### The problem

A belief may be true only under a narrower set of assumptions.

Example:

```text
Claim Y is true under {A, B, C}
```

A downstream action may operate under:

```text
{A, B}
```

The missing semantic question is:

```text
Is {A, B} a valid narrowing of {A, B, C}?
```

A foreign key can prove:

```text
intent.belief_id = belief.id
```

but it cannot prove:

```text
citing_context ⊆ belief.validity_conditions
```

That is a semantic-compatibility problem, not an identity problem.

### Key architectural conclusion

- The **kernel must not learn** about:
  - assumptions A/B/C
  - logical implication
  - subset relationships
  - mathematical meaning
  - domain-specific applicability rules

- Those belong in the **physics/policy/verifier layer**.

---

## 5. The Attestation Pattern

A proposed solution is to create an independent attestation belief.

Instead of:

```text
Restricted claim Y
       ↓
Action Intent directly cites Y
```

use:

```text
Restricted claim Y
       ↓
Attestation belief A
       ↓
Action Intent cites A
```

The attestation would mean, in effect:

> “Context X satisfies the applicability conditions for claim Y.”

The attestation gets its own evidence/debt/promotion lifecycle.

### Why this is attractive

- The semantic subset test remains outside the kernel.
- The final input to Solvent is a normal promoted belief.
- The kernel does not need to understand the mathematics.

### The critical gap

The attestation is useless as a security mechanism if an Action Intent can simply cite the original restricted belief directly.

That created the central **bypass question**:

> What actually prevents direct citation of a restricted belief?

This cannot be answered by saying “agents should use an attestation.” That is policy convention, not enforcement.

---

## 6. Case 1 vs Case 2 for Attestation Enforcement

Two cases were identified.

### Case 1 — Service-mediated Action Intent creation

- All Action Intent creation goes through a policy/service boundary.
- That layer can enforce:
  - restricted belief → attestation required
  - direct citation → rejected
- The Solvent kernel remains generic.
- No kernel understanding of “restricted” is required.

This is the preferred architecture **unless the threat model requires untrusted direct kernel access**.

### Case 2 — Untrusted direct kernel access

If an untrusted actor can directly invoke the kernel's authority-changing capabilities:

- The kernel does not understand “restricted.”
- Therefore it cannot enforce “attestation required” without importing domain semantics.
- That would be the wrong architectural move.

The resulting prerequisite is:

> Determine whether any untrusted actor can bypass the policy/service boundary and directly exercise authority-changing kernel capabilities.

---

## 7. The Human-Only Debt Retirement Problem Is the Same Trust-Boundary Problem

A similar proposed rule was:

> `semantic_applicability_unresolved` may only be retired through a human-facing interface.

The correction was that “human-only” should be treated as an **actor-capability policy**, not as domain semantics in the kernel.

But it introduces the same bypass question:

```text
Agent
  ↓
policy boundary
  ↓
human-gated retirement rule
```

What if an actor bypasses the policy layer and calls the kernel directly?

Then the “human-only” rule becomes decorative unless the trust boundary prevents the bypass.

### Therefore

The attestation-bypass problem and the human-only-retirement problem are really one problem:

> **Which actors are permitted to reach the authority-changing kernel capabilities, and what policy obligations must already be satisfied before they do?**

---

## 8. Debt Vocabulary: Important Correction

An earlier design treated the kernel's fixed debt enum as the canonical vocabulary.

That does **not** generalize cleanly.

### Correct separation

- **Kernel**
  - Stores generic obligation/debt identifiers.
  - Enforces debt lifecycle.
  - Enforces the generic promotion gate:
    - unresolved debt → promotion blocked
    - no unresolved debt → debt no longer blocks promotion
- **Policy/domain layer**
  - Defines what an obligation means.
  - Defines which obligations are relevant to a domain.
  - Defines when a particular obligation can be retired.

This avoids making the kernel domain-aware merely because a new domain needs a new obligation type.

### Important nuance

A generic debt key can still be opaque and structured.

The issue is not whether the kernel stores a string or identifier.

The issue is whether the kernel **understands the meaning** of that identifier.

---

## 9. Should Debt Have a Maximum Count?

The discussion rejected the idea of imposing an arbitrary semantic limit just for “discipline.”

Two separate concepts must be distinguished.

### Semantic limit

Example:

```text
A belief may have at most 10 obligations.
```

This would encode assumptions about how a good workflow should look.

That does not belong in the kernel merely because 100 feels excessive.

### Technical/resource limit

Example:

```text
Reject pathologically large obligation collections
because they create unacceptable request/transaction/resource cost.
```

That can be legitimate.

### Design rule

- Do **not** assume an array needs a limit before verifying the actual representation.
- First inspect whether debt is:
  - a JSON/JSONB array
  - SQL array
  - normalized rows
  - another representation
- If a limit is necessary, determine which layer should own it.
- Prefer a resource/input boundary where that is sufficient.
- Do not invent an arbitrary number such as 6, 10, or 100 without evidence.

The question is:

> What security or resource property does the number enforce?

Not:

> What number looks disciplined?

---

## 10. Trust Boundary: Five Distinct Attack Vectors

The trust-boundary review was expanded into five separate vectors because “can a caller reach the kernel?” is too compressed.

### Vector 1 — Direct kernel function calls

Problem:

- An untrusted actor might directly invoke authority-changing kernel functions.

Controls discussed:

- Runtime authentication/authorization at the process/API boundary.
- I-7 remains useful but is **not** runtime access control.

I-7 means:

> Developers should not accidentally create raw DB write paths outside the kernel.

It does **not** mean:

> An untrusted actor cannot invoke a legitimate kernel function.

Those are different properties.

### Vector 2 — Other internal authority-mutating APIs

Problem:

- A second internal API or tool may expose authority-changing behavior without the same policy controls.

Need to enumerate all such entry points.

Possible controls:

- Narrow authority-changing interfaces.
- Explicit capability surfaces.
- Runtime access control.
- Structural call-path discipline.

### Vector 3 — Lower-level service path bypassing policy

This was considered the most suspicious evolving-codebase risk.

Examples:

- future admin tool
- batch script
- debugging endpoint
- internal service
- migration utility

Each may accidentally bypass the policy layer.

The open design question is whether a capability/session-context mechanism is actually needed.

Important restraint:

> Do not add capability tokens preemptively.

First determine from the real system:

- how many authority-mutating entry points exist
- who can call them
- which are public/internal
- whether a bypass exists
- whether the permitted call graph can be made structurally obvious

Only then decide whether a capability mechanism is justified.

### Vector 4 — Direct CockroachDB writes

This is fundamentally a database-credentials problem.

The important question is:

> Which processes have write credentials to `belief`, `debt`, `action_intent`, and related authority state?

Ideal boundary:

```text
policy/service + kernel
       ↓
DB write credentials

reporting / analytics / ops tools
       ↓
no write access
```

If other processes have write credentials, they can bypass policy/kernel controls regardless of how elegant the application design is.

This is therefore a deployment configuration fact.

### Vector 5 — Deployment/configuration exposure

Examples:

- overly permissive IAM
- debug flag
- raw SQL console
- exposed internal endpoint
- unsafe runtime configuration

This is operational security.

No amount of clean physics-to-Solvent modeling automatically protects against it.

It still belongs on the explicit review checklist.

---

## 11. I-7 Is Not the Whole Trust Boundary

The architecture now distinguishes:

```text
I-7
→ build-time code-hygiene invariant

authentication/authorization
→ runtime actor control

DB credentials
→ database trust boundary

deployment configuration
→ operational trust boundary
```

A clean Solvent codebase can demonstrate I-7 and still be insecure if the deployment gives an attacker authority-changing access.

---

## 12. The Trusted Computing Base (TCB)

A more useful definition emerged:

> **A layer is part of the trusted computing base whenever bypassing or corrupting it can cause an otherwise unauthorized consequential action.**

This is stronger than:

> “The kernel is small, therefore trusted.”

Smallness is useful because it makes auditing easier, but **consequence** is the real criterion.

### Consequence

The physics policy/verifier layer can become part of the TCB because it may decide:

- whether a restricted belief can be relied upon
- whether an attestation is required
- whether an obligation can be retired
- whether a human review is required

Even though these are outside the kernel.

---

## 13. The Physics / Policy Layer Cannot Be Sloppy

Once the policy layer makes security-relevant decisions, it needs similar engineering discipline to the kernel.

Desired properties include:

- explicit rules
- versioning
- reviewability
- auditable decisions
- non-bypassable enforcement
- explicit actor capabilities
- clear provenance
- controlled changes

The kernel stays small by refusing domain semantics.

That does **not** mean the policy layer is allowed to be informal.

---

## 14. The Policy Versioning Question

A proposed policy model looked like:

```text
Policy P17
Attestation A42
A42 is valid under P17
```

This created a deeper question.

The statement:

```text
A42 is valid under P17
```

is itself a claim.

Under the general Solvent worldview, claims are beliefs with:

- their own evidence
- their own debt
- their own promotion
- their own retraction/supersession behavior

Therefore a policy version should not merely be a bare label attached to an attestation if the validity of the attestation depends on that policy.

### Desired conceptual behavior

```text
Policy P17
   ↓
used to evaluate attestation A42
```

If P17 is later retracted or superseded:

```text
P17 invalid
   ↓
A42 may become stale/invalid
   ↓
downstream reliance must be discoverable/rejected
```

This mirrors the existing identity-drift discipline used for ordinary beliefs.

---

## 15. The Kernel Need Not Understand Policy Meaning

Even if policy versions become tracked, the kernel does not need to know what P17 means.

A higher-level verifier/policy system can establish a fact such as:

```text
A42 was evaluated under P17
and the applicability test passed
```

while Solvent's kernel continues to enforce only generic facts such as:

```text
A42 is a promoted belief
```

The domain meaning remains outside the kernel.

---

## 16. A Subtle Layering Correction

The clean diagram:

```text
Physics/domain semantics
        ↓
Policy
        ↓
Kernel
```

is an **authority diagram**, not a pure dependency diagram.

The policy layer may itself depend on Solvent's ledger to track:

- policy beliefs
- policy attestations
- retractions
- supersession
- debt
- audit facts

So both statements can be true:

1. The kernel has no opinion about policy meaning.
2. The policy layer may rely on the kernel's durability and atomicity to keep its own security-relevant history coherent.

This distinction should be explicit in architecture documentation.

---

## 17. The Phase 0–5 Investigation Loop

The proposed sequencing was refined from a simple one-way order into an iterative loop.

### Phase 0 — Provisional model

Draft the debt/policy model enough to expose important future call signatures and security properties.

### Phase 1 — Trust-boundary investigation

Run the five-vector investigation against the provisional model.

For every authority-changing entry point, ask:

- Who can call it?
- Can it bypass policy?
- What does it accept?
- Which fields does it trust?
- What verifies those fields?

### Phase 2 — Refine debt/policy model

Resolve:

- vocabulary ownership
- representation
- auditability
- resource limits
- retirement semantics
- actor capabilities
- policy-version representation
- whether a policy version is itself represented as a belief

### Phase 3 — Trust-boundary re-check

Re-run the critical V3 checks against the refined model.

This matters because Phase 2 can change what the authority-mutating calls actually need to enforce.

### Phase 4 — Relationship/schema verification

Verify things such as:

- `belief_edge` structure
- edge typing/labels
- representation of `supersedes`
- attestation relationships

### Phase 5 — Freeze decision

Only after the model and trust boundary stabilize should the design be considered complete enough for implementation/freeze decisions.

---

## 18. The Bidirectional Dependency

The important correction is:

```text
A/V3 ↔ B
```

not simply:

```text
A/V3 → B
```

Why?

- V3 constrains what policy-mediated operations are allowed to look like.
- B can change what V3 itself needs to check.

Example:

If debt retirement becomes:

```text
RetireDebt(
    belief,
    obligation,
    policy_version,
    review_context,
    actor,
    evidence
)
```

then each new field creates another potential truthfulness problem.

The trust review must ask not only:

> Can an unauthorized actor call this?

but also:

> For each security-relevant field, what verifies that its content is truthful?

---

## 19. Every New Policy Field Is Another Place a Lie Can Enter

Example fields:

- `actor`
  - Is it truly a human or merely a script claiming to be human?

- `policy_version`
  - Is it the correct/current applicable policy rather than an old permissive one?

- `evidence`
  - Does the referenced artifact actually establish what retirement claims it establishes?

This is distinct from path authorization.

Therefore the trust-boundary review must examine both:

```text
Who can call?
```

and:

```text
Why should we believe the arguments?
```

---

## 20. Domain Semantics as Part of the TCB

The same consequence test should be applied one level above policy.

If changing domain semantics can cause an otherwise unauthorized consequential action, then the domain-semantics mechanism is also part of the TCB.

This does not mean domain semantics must move into Solvent.

It means the system should explicitly govern the mechanism carrying those semantics.

Examples:

- versioned policy repository
- reviewed verifier configuration
- controlled mathematical specification
- auditable semantic rules

The physical/semantic layer can remain outside the Solvent kernel while still being recognized as security-critical for the demo.

---

## 21. Canonical Trust Architecture

The emerging architecture is:

```text
                TRUSTED COMPUTING BASE

    ┌───────────────────────────────────────────┐
    │ Domain semantics                          │
    │                                           │
    │ What claims mean                          │
    │ What assumptions imply                    │
    │ What obligations mean                     │
    └───────────────────────────────────────────┘
                      ↓
    ┌───────────────────────────────────────────┐
    │ Policy / verifier / human review          │
    │                                           │
    │ Applicability                             │
    │ Attestation requirements                  │
    │ Human-gated obligations                   │
    │ Actor capabilities                        │
    │ Policy versions                           │
    └───────────────────────────────────────────┘
                      ↓
    ┌───────────────────────────────────────────┐
    │ Solvent kernel                            │
    │                                           │
    │ Identity                                  │
    │ Currency                                  │
    │ Debt lifecycle                            │
    │ Promotion gate                            │
    │ Authority binding                         │
    │ Atomic transitions                        │
    └───────────────────────────────────────────┘
                      ↓
                  Execution
```

This is primarily an authority decomposition.

The kernel may also serve as the durable substrate on which higher-level policy history is recorded.

---

## 22. What Must NOT Enter the Solvent Kernel

The negative boundary is intentionally explicit.

The kernel should never need code that inspects:

- what a physics assumption means
- why an obligation exists
- whether one condition implies another
- whether a context is mathematically valid
- what a domain-specific debt item means
- whether a citation is semantically applicable

Instead, those become higher-level facts.

Solvent receives something structurally simple:

```text
"This belief is promoted."
```

and enforces the consequences of that fact.

---

## 23. What Likely Belongs in the Physics Demo

Once the kernel is frozen, the physics-verifier work should live outside the kernel.

A likely structure is:

```text
demos/
  physics-verifier/
    README.md
    swarm/
    verifier/
    policy/
    attestations/
    scenarios/
    evidence/
    adapters/
    scripts/
```

The demo should show how a swarm/harness becomes **Solvent-ready** without changing Solvent's kernel.

The swarm can perform:

- theorem exploration
- proof attempts
- counterexample search
- numerical checks
- adversarial attacks
- evidence collection
- semantic applicability checks
- attestation generation

Solvent remains the authority boundary.

---

## 24. Important Test Philosophy for the Physics Verifier

The same “prove it structurally” discipline should apply to the demo.

Do not rely on:

> “Agents are instructed to cite the attestation.”

Prefer tests that prove:

```text
restricted belief + no attestation
    → intent rejected

restricted belief + valid promoted attestation
    → intent allowed

attestation based on retracted policy/claim
    → reliance becomes invalid/discoverable

human-only obligation + agent caller
    → rejected

human review + valid policy context
    → obligation can be retired

T1-bound intent + T2 authority
    → rejected
```

The final physics demo should make the same principle visible at the application level:

> **The swarm can propose; the verifier can reason; the policy layer can interpret; Solvent decides whether the resulting authority state is structurally valid.**

---

## 25. Current Open Questions Before Implementation

1. **Trust boundary**
   - Can any untrusted actor directly reach the kernel?
   - Can any internal service bypass the policy layer?
   - What processes have direct DB write credentials?
   - What deployment/configuration paths expose authority-changing capabilities?

2. **Debt representation**
   - Is debt an array, rows, or another structure?
   - Does it actually need a technical cardinality bound?
   - Which layer should enforce that bound?

3. **Debt vocabulary**
   - Can domain vocabulary live entirely outside the kernel?
   - Are current fixed debt enums actually kernel semantics or merely one deployment's policy vocabulary?

4. **Attestation enforcement**
   - Is Action Intent creation always service-mediated?
   - If so, can the policy layer reliably enforce restricted-belief → attestation-required?

5. **Human-gated obligations**
   - Is the human-only rule enforced by the same trusted policy boundary?

6. **Policy versioning**
   - Should a policy version itself be a Belief?
   - If a policy version is retracted, how does invalidity propagate to attestations?

7. **Relationship model**
   - Does `belief_edge` support typed relationships?
   - Is `supersedes` distinguishable from dependency/attestation relationships?

8. **Argument truthfulness**
   - For every security-relevant argument added to policy operations, what verifies its content rather than merely its presence?

---

## 26. Current Architectural Position

The strongest conclusions reached so far are:

- Solvent should remain a **small trusted authority core**, not a physics reasoning engine.
- The physics verifier should not cause new kernel features simply because semantic verification is difficult.
- Restricted-claim applicability is a **semantic policy problem**, not a kernel mathematics problem.
- Attestation only works if the policy/service boundary prevents direct bypass.
- Human-only retirement has the same trust-boundary dependency.
- Debt vocabulary should be **policy-owned**, while the kernel owns generic debt lifecycle/invariants.
- A debt-count cap should only exist if there is a concrete technical/resource reason; an arbitrary discipline number does not belong in the kernel.
- I-7 is a useful build-time invariant but is not runtime authentication or access control.
- The trust boundary has five distinct attack vectors and needs five distinct answers.
- The policy/verifier layer may be part of the TCB even though it is outside the kernel.
- Domain semantics may also be part of the TCB when their corruption can change authorization outcomes.
- Policy versioning should be treated with the same identity, provenance, revision, and stale-propagation discipline as other security-relevant claims.
- The Phase 0–5 process should iterate because the trust-boundary design and debt/policy design constrain each other.
- The kernel-freeze rule remains intact: after freeze, new capabilities should default to service, adapter, executor, deployment, policy, demo, or documentation layers.

---

## 27. The Core Principle

The emerging design can be reduced to one sentence:

> **The physics system is allowed to understand meaning; the policy layer is allowed to decide semantic applicability; Solvent is allowed to enforce durable structural authority—but the security of the whole system depends on preventing untrusted actors from bypassing the layers above the kernel that carry meaning.**

That is the main architectural lesson driving the current physics-verifier exploration.
