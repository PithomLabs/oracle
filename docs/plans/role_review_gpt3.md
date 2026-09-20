## Role / Boundary Matrix v0.3 — consolidated review

**Verdict: GREEN / lock it.**

This revision closes the remaining substantive findings from the previous adversarial rounds. The document is now internally coherent enough to serve as the **normative ownership/boundary projection** of Workflow Specification v0.3 and as the direct input to the Conformance Test Matrix.

The important point is that v0.3 solved the remaining problems **without introducing the architecture the adversarial reviews repeatedly tried to smuggle in**.

### What is now convincingly closed

**Non-effect declaration lifecycle**

This was the strongest remaining v0.2 issue. v0.3 now requires every declared capability boundary to have an owner, version, effective reference, effect classification, scope, and identity definition, while effect-capable capabilities get the additional authority/execution requirements. A non-effect → effect transition now explicitly requires version change, consequential-boundary coverage, and re-conformance. 

That closes the lifecycle hole.

**Ordinary-work ambiguity**

The old `Ordinary Work — effect-capable` wording is gone. The matrix now has `Non-Effect Capability Invocation` and `Effect-Capable Operation Invocation`, and explicitly states that the latter must proceed through:

```text
Effect-Capable Operation Invocation
→ Consequential Proposal
→ Authorization
→ Execution Eligibility
→ Execution
```

 

That is substantially cleaner.

**Ownership/effect-owner consistency**

The boundary matrix now separately represents:

```text
Semantic / authority owner
Effect / outcome owner
Enforcement / implementation responsibility
```



And the acceptance criteria explicitly require both authority and effect/outcome ownership. 

The failure table also received the same treatment, so the earlier "fix §5 but leave §13 broken" problem is gone. 

**Formulate / Interpret boundary consistency**

Formulate and Interpret are now explicitly internal phases and are excluded from the external boundary matrix.  

Next Work now has its corresponding human-intervention entry. 

Good cleanup.

**Execution vs external system of record**

This is now very clear:

> Executor owns attempted operation and execution reporting.

> External system of record is authoritative for whether the external effect actually occurred.

> Executor cannot override that evidence. 

That is an important distinction for ambiguous and reconciliation scenarios.

**Fail-closed network behavior**

The matrix now explicitly treats unavailable/failed authority verification as unavailable/invalid and requires rejection. 

That closes the classic distributed-systems "availability wins over authorization" hole.

**Under-declaration**

The integration owner is explicitly responsible for completeness of operation identity, and omission of an effect-relevant parameter is a conformance failure. 

That is the right accountability mechanism without turning Solvent into a global schema registry.

**Terminated / revoked**

The outcome taxonomy now distinguishes:

```text
not attempted
rejected
attempted
succeeded
failed
terminated / revoked
ambiguous
```

with explicit semantics for termination/revocation. 

This is useful and materially improves the conformance model.

**Declaration version at execution**

Execution evidence now carries declaration/identity version, and the anti-pattern section explicitly rejects using a non-current effective declaration.  

That closes the declaration-drift problem.

**Review disposition**

This is a surprisingly valuable addition. The matrix now records rejected, accepted, deferred, and implementation-option decisions, including explicit rejection of:

* mandatory workflow runtime
* central Solvent effect registry
* mandatory operations gateway
* automatic Conductor → Solvent revocation
* mandatory verification SDK
* mandatory cryptographic mechanism/RFC 8785



That should stop the same architectural proposals from being rediscovered in every subsequent review.

---

# Two very minor observations

These are not blockers.

### 1. Declaration version wording

The anti-pattern says an effect-capable integration must use the **currently effective declared version**. 

For long-running authorized operations, the more precise principle is actually:

```text
authorization binds to declaration/identity version V
execution validates against applicable version V
```

rather than simply "latest version."

The matrix elsewhere already moves in that direction through explicit declaration/identity-version evidence. 

I would not reopen the document over this unless the Conformance Test Matrix exposes a concrete conflict.

### 2. `Effect / outcome owner` is intentionally contextual

Some rows say:

```text
Executor
```

while execution with an external SOR says:

```text
Executor for reporting; external SOR for actual occurrence
```

That's correct. It is no longer the ambiguous slash-owner problem from v0.2 because the semantic distinction is explicitly defined in §3. 

---

# The architecture now passes the important adversarial questions

The new §26 is exactly what we want as a bridge to testing.

It now tests:

* non-effect → effect drift
* ordinary-work bypass
* X→Y operation substitution
* identity consistency
* Solvent bypass
* authority verification failure
* UNKNOWN policy
* ambiguous execution
* rejection vs failure
* termination/revocation
* human shadow authority
* Conductor cancellation
* declaration drift
* enforcement-path bypass
* undeclared effects
* workflow-semantic absorption
* Agent overriding authoritative outcomes



That is now a much better test inventory than another round of abstract architectural review.

# Final disposition

**Lock Role / Boundary Matrix v0.3.**

I would not make another substantive architecture revision before the next artifact.

The current chain is now:

```text
Workflow Specification v0.3
        ↓
Role / Boundary Matrix v0.3   ← LOCK
        ↓
Conformance Test Matrix v0.1
        ↓
BM-IST Validation Scenarios
        ↓
Agent Skill
        ↓
Conformance Harness
        ↓
POC Runbook
```

The document itself correctly states that it does **not** authorize a new workflow runtime or infrastructure component and is now the basis for those next deliverables. 

**The architecture phase for this artifact is done. The next challenge should be empirical: can the Conformance Test Matrix actually make these boundaries fail visibly when an implementation violates them?**
