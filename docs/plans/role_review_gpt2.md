## Review of Role / Boundary Matrix v0.2

**Verdict: GREEN / ready to lock, with only two minor wording issues worth fixing before treating it as final.**

The v0.2 revision successfully addresses essentially all substantive findings from the three adversarial reviews plus my consolidated review. The important changes are not cosmetic; they close the actual architectural gaps.

### What v0.2 fixed correctly

**1. The Ordinary Work bypass is closed.**

This was the most important defect in v0.1. v0.2 now explicitly distinguishes:

* ordinary work over a declared **non-effect capability**, and
* ordinary work invoking an **effect-capable capability**.

More importantly, the normative rule says an operation capable of producing an external effect is consequential **regardless of workflow phase or invoking client** and must enter the applicable authorization boundary. 

That directly resolves the strongest adversarial finding.

The new explicit limitation is also correct:

> Loop Engineering does not provide a universal oracle for discovering undeclared side effects inside arbitrary third-party software.



That is the right architectural answer. It defines the conformance boundary instead of inventing an impossible universal effect detector.

---

**2. "Integration" is no longer a shadow role.**

The new terminology section is exactly what was needed:

> Integration denotes an implementation boundary or adapter belonging to an existing role boundary. It is not an independent Loop Engineering role.

It may be a wrapper, adapter, proxy, gateway, sidecar, library, callback, etc. 

That is excellent because it permits practical implementation techniques without turning any of them into Loop Engineering architecture.

---

**3. Ownership vs enforcement is now much cleaner.**

The new separation:

> **One semantic owner; possibly another enforcement mechanism.**

is probably the most important structural improvement after the ordinary-work fix. 

The matrix can now say:

```text
Solvent       = authority owner
Integration   = enforcement mechanism
Executor      = effect/outcome owner
```

without pretending those are the same responsibility.

This also makes the acceptance criterion internally consistent: "one owner" now means **one semantic/authority owner**, rather than one implementation component.

---

**4. Exact operation identity is now substantially stronger.**

The revision explicitly separates semantic effect inputs from transport metadata and requires excluded fields to be declared immaterial. 

That closes the parameter-normalization ambiguity without introducing a mandatory canonicalization or cryptographic scheme.

This is the right level of abstraction.

---

**5. Trust-basis language is now properly security-oriented without over-specifying implementation.**

The new wording requires independently verifiable evidence and explicitly prohibits relying solely on an unverified client claim. It also gives valid mechanism classes while stating cryptography is not mandated. 

That is much better than either extreme:

```text
"trust whatever the client says"
```

or

```text
"all integrations must use protocol X / cryptographic mechanism Y"
```

The contract remains mechanism-neutral while preserving the security property.

---

**6. UNKNOWN is now operationally defined.**

v0.2 explicitly requires a maximum unresolved-wait/age policy and requires that policy to be actionable. It allows either an explicit operator-controlled indefinite wait or an existing-boundary remediation/termination mechanism. 

Most importantly, it explicitly refuses to invent generic:

```text
Abandoned
Timeout
Escalation
Human Escalation
```

states into Loop Engineering.

That is exactly the correct architectural restraint.

---

**7. Long-running validity precedence is now unambiguous.**

The new rule is clean:

> Authority-owner constraints are authoritative.
> An integration may tighten them, but cannot extend, weaken, or override them.

The examples make the intended effective validity obvious. 

That closes the Solvent-versus-integration precedence problem.

---

**8. Execution rejection is now explicitly distinct from failure.**

This is an important conformance primitive:

```text
not attempted
rejected before effect
attempted
succeeded
failed
ambiguous
```



The distinction:

```text
rejected = no external effect attempted
failed   = external effect was attempted but failed
```

is now explicit and testable.

---

**9. The evidence matrix is substantially more mature.**

v0.2 now has explicit evidence for:

* proposal version
* human intervention
* declaration version
* execution rejection
* long-running progress
* reconciliation determination



That is exactly what the Conformance Test Matrix will need.

---

**10. The matrix now explicitly preserves Workflow Specification precedence.**

This was another good correction:

> Workflow Specification v0.3 is the normative parent artifact.



That prevents the projection from quietly becoming a competing normative specification.

---

# Two residual issues

These are **not blockers**, but I would clean them up now because the artifact is about to become a test-design basis.

### 1. "Executor / external system" still appears as a semantic owner in a few places

For example, the execution outcome is:

> Executor / external system of record

and the execution row similarly treats them jointly. 

This is defensible because the Executor may rely on an external system of record. But the matrix's own ownership discipline would be slightly cleaner with:

```text
Semantic owner: external system of record
Reporting/reconciliation responsibility: Executor
```

where that distinction actually applies.

That preserves the important principle:

```text
Executor reports what the external SOR says happened;
Executor does not become a substitute for the actual external truth.
```

This is especially useful for ambiguous outcomes.

### 2. "Ordinary Work — effect-capable capability" could still be phrased more sharply

The row is now correct, but:

> Agent/client → Effect-capable integration / authority boundary

could be misread as saying the ordinary-work phase itself has an authority transition.

The later boundary rule fixes this semantically, so this is not an architectural flaw. For maximum clarity, I would label it:

```text
Ordinary Work — invokes effect-capable capability
```

That emphasizes that the *phase* remains ordinary work while the *operation* crosses the consequential boundary.

---

# What I would not change

At this point, I would resist all further proposals to add infrastructure.

Specifically, **do not add**:

```text
Loop workflow runtime
Operations Gateway as mandatory component
Solvent effect-schema registry
Conductor authorization router
Generic timeout daemon
Generic reconciliation service
Workflow event store
Workflow heartbeat
New lifecycle states
Second human authority mechanism
```

The v0.2 document already explicitly protects that boundary.  

The first adversarial review's concern about ephemeral Executors and reconciliation remains a legitimate **integration implementation concern**, but v0.2 correctly handles it by requiring declaration of reconciliation and persistence/recovery expectations rather than creating Loop Engineering infrastructure. 

---

# The important architectural test now passes

The document's adversarial questions are much stronger now:

1. Can ordinary work invoke an effect-capable operation without entering the consequential boundary?
2. Can X be authorized and Y executed?
3. Is operation identity identical across proposal/authorization/execution?
4. Can an effect-capable integration bypass Solvent?
5. Can a client disguise an effect-capable operation as ordinary work?
6. Can UNKNOWN become indefinite undocumented waiting?
7. Can ambiguity become success?
8. Can rejection become failure?
9. Can human intervention create shadow authority?
10. Can Conductor cancellation alter Solvent authority?
11. Can declaration drift invalidate old conformance?
12. Can the harness bypass real enforcement?
13. Can an undeclared effect surface escape the declaration model?
14. Can a new component absorb workflow semantics under another name?



That is now an excellent handoff to the Conformance Test Matrix.

## Final disposition

**Role / Boundary Matrix v0.2: GREEN.**

I would make the two minor wording refinements above, then **freeze it**.

The important thing is that v0.2 did not "solve" the adversarial pressure by adding more architecture. It solved it by making the contracts, ownership, evidence, and boundaries more precise. That is exactly the direction Loop Engineering should take.

**Next artifact: Conformance Test Matrix v0.1.**

Its job is now straightforward: take the 14 adversarial questions in §24 and turn each into concrete positive, negative, substitution, long-running, ambiguity, and production-enforcement tests. 
