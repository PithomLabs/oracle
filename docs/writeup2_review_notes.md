This critique is materially useful, and I would **change our roadmap based on it**. It catches several places where we were beginning to declare more than the evidence actually established.

The strongest finding is **C2: operation mismatch**. The Phase 1 tests proved wrong-actor rejection and missing-authorization rejection, but they did not prove that authorization for **operation X** is rejected when execution attempts **operation Y**. That is the most direct test of the exact-binding invariant. The critique is right that this belongs **before any real external effect**. 

I also agree with the broader methodological criticism around the PASS. We should not continue treating a verdict as meaningful unless the pass/fail criteria and falsification conditions are explicit. The critique correctly identifies the missing pre-registration as a real methodological gap, not paperwork. 

### One correction to our previous roadmap

I would now replace:

```text
Phase 1 complete
        ↓
Real external effect
```

with:

```text
Phase 1 controlled proof
        ↓
Phase 1.5 safety/binding gate
        ↓
Phase 2 real external effect
```

The **Phase 1.5 gate** should include:

```text
✓ wrong actor → DENY
✓ missing authorization → DENY
✓ operation X authorized, operation Y executed → DENY
✓ duplicate delivery → behavior explicitly characterized
✓ stale/expired authorization → DENY or declared policy
✓ verification unavailable → correct fail-closed/UNKNOWN behavior
```

The first two already passed. The latter four need evidence. The critique correctly notes that the current rejection-cause coverage is only partial. 

### H1 is also important

The Conductor ↔ Solvent edge in the diagram needs clarification.

Our intended architecture is still:

```text
Agent
  ├──→ Conductor          (work protocol)
  │
  └──→ Solvent            (authority interaction, where required)
```

not:

```text
Agent
  ↓
Conductor
  ↓
Solvent
```

unless Conductor is explicitly acting as a **transport-only relay** and carries no authority semantics.

The critique is right that the diagram currently makes that ambiguous. 

So we should correct the diagram before the protocol phase.

### H2 should become a gate

The critique's freeze/conformance point is also important. Freezing Solvent proves stability, but it does **not automatically prove conformity to every version of the specification**. 

Before Phase 2, we therefore need:

```text
Frozen Solvent 7602699
        ↓
Conformance delta against the specific contract under test
        ↓
either:
    compliant
or:
    explicit accepted deviation
```

No need to modify Solvent. The point is to know exactly what the experiment is claiming.

### Declaration becomes mandatory before real effect

I agree with H2b too. The real external-effect experiment should not run against an undeclared operation.

Before GitHub becomes a real effect:

```text
operation declaration
    ↓
declared executor/capability
    ↓
validity / trust basis
    ↓
effect classification
    ↓
idempotency semantics
    ↓
pinned version
    ↓
real execution
```

Otherwise we would be testing an implementation path without testing the declaration contract that the specification says governs that path. 

### Exactly-once needs narrower wording

The critique is also right about M2.

We can legitimately claim:

> "The controlled positive execution test invoked the fake executor exactly once."

We cannot yet claim:

> "The system provides exactly-once execution."

Those are different claims.

A retry/duplicate-delivery experiment is therefore needed before making the stronger statement. 

## Revised roadmap

I would lock the roadmap now as:

```text
PHASE 1
Controlled Reference Loop
    ✓ happy path
    ✓ wrong actor denial
    ✓ no-authorization denial
    ✓ valid authorization → one controlled executor invocation

        ↓

PHASE 1.5
Safety / Binding Gate
    → operation mismatch
    → duplicate delivery
    → stale authorization
    → verification unavailable
    → explicit pass/fail criteria
    → frozen-Solvent conformance delta
    → capability declaration

        ↓

PHASE 2
Real External Effect
    → real GitHub operation
    → external SOR observation
    → authoritative effect/result
    → ambiguous outcome handling

        ↓

PHASE 3
Broader Adversarial Validation
    → cancellation
    → revocation
    → concurrency
    → replay
    → recovery
    → substitution

        ↓

PHASE 4
Protocol + Skills Formalization
    → Agent Skill
    → Conductor Skill
    → Solvent authority contract
    → Executor capability contract

        ↓

PHASE 5
Conformance / Substitution

        ↓

PHASE 6
BM-IST / Domain Validation
```

That is stronger than our previous roadmap because it puts the **highest-value safety property—the exact operation binding—before the first real effect**.

The critique's closing recommendation is therefore the right one: **run the mismatch test, clarify the Conductor/Solvent edge, audit the frozen baseline against the actual contract, establish the declaration, then run Phase 2.** 

I would treat this critique as an architecture-review gate rather than as a rejection of the Reference Loop. The core result still stands; the claim simply becomes narrower and more rigorous.
