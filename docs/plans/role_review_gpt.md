## Adversarial review — Loop Engineering Role / Boundary Matrix v0.1

**Verdict: GREEN / ready to become the basis for the Conformance Test Matrix, with a few targeted corrections before locking v0.1.**

The matrix successfully does what it should: it translates Workflow Specification v0.3 into explicit ownership without accidentally introducing another runtime, state machine, authority layer, or coordination system. The core ownership split is particularly clean: Agent → agency, Conductor → coordination, Solvent → authority, Executor → effect, Domain → meaning/truth. 

The strongest part is that it preserves the central invariant rather than merely repeating it. The matrix explicitly separates coordination, authority, execution, correlation, and domain truth. 

### What is already very strong

**1. There is no hidden fifth runtime role.**

The document explicitly states that it does not create a state model, workflow runtime, authority system, database, event store, or UI. 

That is exactly the correct containment boundary.

**2. The consequential path is properly split.**

This sequence is architecturally sound:

```text
Agent
  ↓ proposal
Solvent
  ↓ authorization
Effect-capable integration / Executor
  ↓ execution
External system
  ↓ outcome
Agent / Domain
```

The matrix does not collapse authorization and execution, and it makes the exact operation binding explicit across proposal, authorization, and execution.  

That is probably the single most important thing this artifact needed to establish.

**3. Fail-closed behavior is correctly located at the effect-capable boundary.**

The requirement that missing, invalid, stale, or mismatched authorization must cause rejection is clear, while the Executor is prevented from deciding policy. 

This avoids the dangerous architecture:

```text
Conductor → Executor
Executor decides whether action is okay
```

and preserves:

```text
Solvent → authority
Executor → effect
```

**4. Human intervention is handled without creating shadow authority.**

The human matrix is particularly good because it distinguishes *review/intervention* from *authority creation*. The explicit statement that human presence does not bypass Solvent is important. 

The distinction between:

```text
Human rejects continuation
```

and

```text
Human creates an independent authorization path
```

is exactly right.

**5. Failure ownership is unusually clean.**

The distinction among:

```text
DENIED
UNKNOWN
INVALID/STALE
MISMATCH
EXECUTION FAILURE
EXECUTION AMBIGUOUS
DOMAIN REJECTION
```

is strong and should translate directly into conformance cases. 

Especially important:

> Execution failure ≠ authorization denial

and

> Execution ambiguous ≠ success/failure

That prevents a large class of bad orchestration logic.

**6. Persistence ownership is correctly separated from semantic ownership.**

This is a subtle but important section. The fact that something is recorded in Conductor, an Executor, or a test harness does not make that component the semantic authority for the fact. 

That is exactly the distinction needed to prevent accidental authority migration through database design.

**7. The production-conformance rule is excellent.**

The requirement that integration conformance exercise the *real authorization/enforcement path*, even when the final effect is safely neutralized, is one of the strongest parts of the artifact. 

Without this, a conformance suite could easily certify a fake architecture.

**8. Runtime/tool substitution is now explicit.**

The matrix makes the important claim testable:

```text
GPT → Claude
Claude → scripted client
API → Temporal
Executor A → Executor B
```

must not change semantic ownership. 

That directly supports your broader goal of proving that the workflow is not tied to a particular AI model or orchestration framework.

---

# The remaining issues

These are not architectural failures. They are **precision gaps that the Conformance Test Matrix should close**.

### 1. "Authorization" row slightly conflates Solvent with integration

This row says:

> `Solvent → Effect-capable integration / Executor`

That is conceptually understandable, but it could leave one ambiguity:

**Does Solvent actually expose authority directly to the Executor, or does an integration retrieve/verify it through some defined authority boundary?**

The Workflow Specification already addresses the trust basis, so this is not a conceptual problem. But the matrix should avoid implying a mandatory transport topology.

I would change the receiving participant wording to:

```text
Authorization:
Solvent → Effect-capable integration
```

and make the Executor consumption of authorization explicit in the next boundary.

That keeps:

```text
Solvent = authority decision
Integration = enforcement boundary
Executor = execution
```

cleanly separated.

### 2. "Execution Eligibility" needs one more explicit sentence

The current row is good:

> Validate required authorization before effect. 

But there is a potentially important ambiguity:

**Who is responsible if the integration validates authorization correctly but the Executor subsequently receives a materially different operation?**

The answer should be that the **effect-capable boundary remains responsible for ensuring that the exact operation handed to the Executor is the one that was authorized**, not merely that some authorization object existed.

The existing exact-binding section strongly implies this, but one explicit sentence would eliminate the ambiguity.

Suggested addition:

> The effect-capable integration MUST bind the operation presented to the Executor to the same exact operation identity that was authorized.

### 3. Human "cancel" should not imply authority revocation

The human intervention row says:

> Before Execution → Cancel/change progression where supported. 

This is potentially dangerous terminology because elsewhere the architecture explicitly says Conductor cancellation must not silently mutate Solvent authority. 

The document should distinguish:

```text
cancel execution progression
```

from:

```text
revoke authority
```

Suggested clarification:

> Cancellation or intervention does not by itself revoke an existing authorization; any revocation follows the authority owner's applicable mechanism.

That keeps the previous architectural decision intact.

### 4. `N/A` under authority/effect ownership could be misunderstood

Rows such as Discover, Formulate, Assign/Claim and Next Work use `N/A`. 

This is fine, but a hostile reader could interpret `N/A` as "no owner exists."

The intended meaning is:

```text
No consequential authority/effect exists at this boundary.
```

A tiny note in the table header or preceding text would make that unambiguous.

### 5. "Domain" as a participant needs one defensive qualifier

The role model currently lists Domain as a participant and authority for domain-specific correctness/acceptance. 

Architecturally this is correct, but there is one possible future failure mode:

A developer could interpret "Domain" as a service/component that must be added to the infrastructure.

The spec already repeatedly says domain semantics remain external to infrastructure. I would add:

> "Domain" is a semantic authority category, not a required infrastructure component.

That makes the boundary harder to misread.

### 6. Scenario 7.2-style exact binding needs to be reflected in the acceptance test vocabulary

The matrix correctly requires:

> same declared operation identity definition at proposal, authorization, execution. 

The Conformance Test Matrix should therefore test mutation attacks explicitly:

```text
Authorized:
    DELETE account=123

Attempt:
    DELETE account=124
```

and:

```text
Authorized:
    transfer(source=A, destination=B, amount=100)

Attempt:
    transfer(source=A, destination=B, amount=1000)
```

The critical point is that **correlation success must not compensate for operation identity failure**.

The matrix already asks the right question in §19:

> Can a proposal be authorized for X but executed as Y? 

The next artifact should turn that directly into executable tests rather than expanding this document further.

---

# One issue I would explicitly *not* add

Do **not** add:

```text
Workflow Coordinator
Workflow Runtime
Loop Controller
Workflow State Manager
Workflow Event Bus
Workflow Scheduler
Workflow Policy Engine
Workflow Authority
```

as additional participants.

The matrix's strongest architectural feature is precisely that Loop Engineering is a **contract across existing components**, not a new component. The document explicitly protects that boundary. 

Likewise, do not add a generic `WorkflowState` row merely to make the matrix look more formal. The document correctly says workflow phases are vocabulary rather than an authoritative state model. 

---

# Final assessment

| Dimension                        | Verdict    |
| -------------------------------- | ---------- |
| Role ownership                   | **GREEN**  |
| Authority boundary               | **GREEN**  |
| Execution boundary               | **GREEN**  |
| Exact operation binding          | **GREEN**  |
| Human intervention               | **GREEN-** |
| Failure semantics                | **GREEN**  |
| Persistence ownership            | **GREEN**  |
| Production conformance boundary  | **GREEN**  |
| Tool/runtime substitution        | **GREEN**  |
| Risk of creating another runtime | **GREEN**  |
| Domain leakage risk              | **GREEN-** |
| Ready for next artifact          | **YES**    |

### Recommended disposition

**Do not conduct another broad architecture rewrite.**

Make the few wording-level corrections above, then **freeze Role / Boundary Matrix v0.1**.

The next artifact should now be the **Conformance Test Matrix v0.1**, and it should derive its cases directly from §19's adversarial questions rather than inventing a new conceptual framework. The matrix already ends at exactly the right boundary: it is the basis for the conformance tests, BM-IST scenarios, Agent Skill, harness, and runbook, while explicitly not authorizing new infrastructure. 

The architecture is no longer the bottleneck. **Evidence is.** The next phase should prove that these boundaries survive adversarial execution, denial, mutation, substitution, ambiguity, and failure.
