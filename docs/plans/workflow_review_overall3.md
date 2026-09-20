I consolidated the three new adversarial reviews with my prior review of `Workflow Specification v0.2`. The result is now quite clear: **the architecture is effectively closed; the remaining attack surface is primarily the certification machinery and specification hygiene.**

The three reviewers converge strongly on that point. One calls the current version essentially locked, while the other two identify residual issues around operation identity, reconciliation, declaration authority, observability, and certification integrity.  

# Consolidated verdict: **AMBER — very close to lock**

I would **not redesign the architecture**.

I would make one final specification-cleanup pass before generating the Role/Boundary Matrix and Conformance Test Matrix.

## 1. The core binding chain is now fundamentally sound

This is the biggest success.

The current specification establishes:

```text
Operation identity
      ↓
Authorization evidence
      ↓
Exact operation binding
      ↓
Executor comparison
      ↓
Fail closed
```

The latest review explicitly recognizes that the previous "identity/class subsumption" problem has been substantially closed. 

That is the right architecture.

The remaining issue is not that binding is missing. It is that the **identity function itself must be the same at proposal, authorization, and execution**.

Add a simple rule:

> **The same declared operation-identity definition MUST be used when the operation is proposed, authorized, and executed.**

That closes the implementation-divergence attack without introducing a new runtime.

Also, I agree with the Gemini review that you should **not over-engineer canonicalization here**. RFC 8785 or cryptographic signing may be useful implementation techniques, but they should not become mandatory Loop Engineering machinery. 

---

# 2. Declaration lifecycle is now the real control surface

The declaration mechanism is now extremely important:

```text
effect-capable operations
authorization boundary
operation identity
validity model
replay policy
outcome model
reconciliation
```

The reviews correctly observe that this has become a load-bearing contract. 

The current document already requires versioning and re-conformance when capabilities change. That's good.

What remains:

### Add declaration provenance

Require:

```text
declaration owner
declaration version
effective version
```

and:

> The declaration content MUST be retrievable by the conformance harness and inspectable by the applicable authority/domain owner.

The latest Z review correctly points out that merely exposing the declaration **version** is insufficient if the harness cannot read the declaration itself. 

This is a small but important correction.

---

# 3. Production conformance is now the main certification risk

The two-layer model is excellent:

```text
Protocol Fixture Conformance
        +
Integration Conformance
```



But the new Z review finds a serious issue with **how integration conformance is performed**:

> A sandbox/dry-run could short-circuit before the actual authorization enforcement point, producing a false green result. 

That's a very good catch.

### Required refinement

Add:

> **Integration conformance MUST preserve the production authorization/enforcement path. Neutralization of the external side effect MUST occur at or after the enforcement boundary, not before it.**

So:

```text
production path
    ↓
authorization verification
    ↓
operation binding check
    ↓
fail/allow decision
    ↓
SAFE/SANDBOX EFFECT
```

not:

```text
test harness
    ↓
intercept before security logic
    ↓
"pretend" execution
```

This is probably the most important new conformance-process finding.

---

# 4. Trust basis for authorization evidence still needs one sentence

The current contract allows:

* trusted artifact
* Solvent callback
* signed data
* deployment-appropriate mechanism



That's intentionally implementation-neutral.

But the latest review correctly notices the remaining phrase:

> "directly and verifiably available"

doesn't say **why the Executor trusts the source**. 

Add:

> **The integration declaration MUST identify the trust basis by which the Executor authenticates authorization evidence to the declared authority owner.**

That's enough.

Do not mandate cryptography.

---

# 5. The harness is not allowed to become a hidden workflow engine

Gemini's concern here is understandable:

> a sophisticated harness could itself become a workflow runtime. 

I don't think the architecture actually requires that.

But the specification should explicitly state:

> **The Conformance Harness may orchestrate test execution, but its orchestration behavior has no normative production semantics and MUST NOT be treated as an implementation reference for Loop Engineering runtime behavior.**

That cleanly separates:

```text
test orchestration
≠
workflow runtime
```

The current distinction that the harness is test-only is already good. 

---

# 6. Ambiguous outcomes: the architecture is right, but the review exposes an operational limitation

The current model:

```text
ambiguous
    ↓
Executor/integration reconciliation
```

is architecturally correct.

The Gemini review argues that stateless Executors may be unable to reconcile after they disappear. 

That's a legitimate deployment concern, but **do not solve it by adding a mandatory background worker**.

Instead the declaration should state:

```text
reconciliation owner
reconciliation mechanism
persistence/recovery expectation
```

Possible mechanisms:

```text
external status query
external system of record
persistent executor
manual investigation
human escalation
idempotent retry
```

The important thing is that an integration declaring an ambiguous-capable operation must have a **real reconciliation story**.

---

# 7. Operation identity: keep exact binding normative, but don't mandate canonical JSON

Gemini recommends RFC 8785 canonicalization and an ignore-list for transport metadata. 

I would reject making that normative.

The specification should instead say:

> The declared operation identity MUST distinguish every effect-relevant parameter and MUST exclude declared non-semantic transport metadata.

That leaves implementation choices open:

```text
canonical JSON
typed protobuf
hash of normalized structure
domain-specific stable identifier
```

This is exactly where tool/runtime agnosticism matters.

---

# 8. Missing conformance coverage remains a real issue

The new review correctly observes that the requirements contain obligations that aren't fully represented in the scenario set. 

The final scenario set should explicitly include:

```text
A  Ordinary research
B  Consequential computation
C  Human intervention
D  Client substitution
D2 Orchestration substitution
E  Denial → revision
F  Unknown / reconciliation
G  Long-running execution
H  Persistence / phase inspection
```

and the harness should have explicit support for G/H.

The current document has these scenarios; the remaining job is to ensure the **fixtures can actually execute them**.

In particular:

* TestExecutor needs a held/long-running mode.
* TestConductor or equivalent test observation needs to establish unrelated-work progress.
* Phase/persistence inspection needs a defined inspection mechanism.

The Z review is right that a scenario without a fixture is not yet a conformance scenario. 

---

# 9. Ordering semantics need one final definition

This remains the longest-running unresolved medium issue.

The spec asks the harness to establish:

> relevant ordering. 

But "ordering" can mean different things.

I recommend defining it as:

> **Conformance ordering is causal/happens-before ordering established through the scenario's declared correlation chain; wall-clock timestamps alone are insufficient.**

Then a scenario can establish:

```text
proposal
happens-before
authorization
happens-before
execution attempt
```

without requiring a global clock.

That closes M-B cleanly. 

---

# 10. Production conformance drift still needs a stance

One review points out the difference between:

```text
passes conformance today
```

and:

```text
remains conformant after integration changes
```



The declaration lifecycle already helps.

I would add:

> **Conformance is version-scoped. A change to a declared effect surface invalidates the prior conformance result for the affected boundary until re-conformance completes.**

That is much cleaner than trying to continuously monitor production.

---

# 11. Migrate the mode semantics exactly once

One review still flags autonomous vs semi-autonomous semantics.

The current document now defines:

```text
Autonomous:
no human action is required at declared boundaries.

Semi-autonomous:
at least one declared boundary requires or permits human action.
```

That is good. 

I would not reopen this.

---

# 12. Minor but important document integrity issue: version and numbering

This is the one place where I strongly agree with Z and Qwen.

The attached specification has duplicate section numbers:

* two `§22`
* two `§28`

The review explicitly flags this. 

That should be fixed **before any downstream artifact cites section numbers**.

Also, because this document is materially different from the previous v0.2 content, its own versioning policy should be respected. The review correctly points out that a semantic change should increment the specification version. 

So I would rename this:

> **Loop Engineering Workflow Specification v0.3**

rather than continuing to call it v0.2.

That is not cosmetic. The document explicitly mandates normative version increments.

---

# 13. Merge residue should be removed

There are duplicated normative blocks around:

* integration declarations
* authorization taxonomy

The Z review correctly calls this out. 

For a specification intended to become a conformance oracle, **one normative statement per rule** is important.

Otherwise a later revision can modify one copy and not the other.

Consolidate them before locking.

---

# 14. One point I would reject from the reviews

I would **not** introduce a universal integration helper module merely because 20 Executors might otherwise duplicate verification logic.

Gemini recommends a reusable helper. 

That could eventually be useful as an ecosystem convenience, but it is **not part of the specification**.

It should only emerge as an implementation artifact after the POC demonstrates real duplication.

The same applies to a generic Egress Gateway.

---

# Final consolidated disposition

### Architecture

**GREEN.**

The role separation, authority ownership, non-persistence rule, optional consequential branch, tool/runtime agnosticism, and Growth Gate are all holding.

### Specification mechanics

**GREEN- / nearly closed.**

Exact binding, declarations, evidence, conformance layers, ambiguity, long-running execution, and substitution are now substantially specified.

### Remaining blockers before lock

I would reduce the entire combined review set to **six concrete actions**:

```text
1. Rename this document to Workflow Specification v0.3.
   Fix duplicate section numbers and merge duplicated normative blocks.

2. Add one rule requiring the same declared operation-identity function
   at proposal, authorization, and execution.

3. Add declaration provenance/content accessibility:
   owner + version + retrievable declaration content.

4. Add trust-basis requirement for authorization evidence.

5. Strengthen integration conformance:
   production enforcement path MUST be exercised;
   safe/dry-run neutralization cannot bypass the enforcement point.

6. Define ordering as causal/happens-before evidence rather than
   wall-clock ordering, and ensure G/H scenarios have actual fixtures.
```

The rest should be handled while deriving the Conformance Test Matrix rather than by continuing to grow the Workflow Specification.

## My recommended next state

```text
Requirements v0.3
        ↓
Workflow Specification v0.3   ← finish six items
        ↓
LOCK
        ↓
Role / Boundary Matrix
        ↓
Conformance Test Matrix
        ↓
Agent Skill
        ↓
Conformance Harness
        ↓
BM-IST POC
```

This is now an important transition point. The reviews have stopped finding architectural defects and are increasingly finding **certification-process and specification-integrity issues**. That is a strong signal that we should make this final cleanup, lock the contract, and start proving it against BM-IST rather than continuing theoretical expansion.
