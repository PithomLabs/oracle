# Loop Engineering Workflow Design v1.0

**Status:** Freeze Candidate  
**Parent:** Loop Engineering Workflow Specification v0.3  
**Purpose:** Minimal workflow/control-plane design derived from the Reference Loop, subsequent adversarial review, and the agreed "less is more" constraint.

---

## 1. Scope and Lineage

This document defines the minimal workflow design for:

- Agent Plan/Work participation
- Human plan approval and plan iteration
- Conductor project state and work coordination
- consequential-operation handoff to Solvent
- Executor effect production
- durable project state and history

It does not create:

- a workflow runtime
- a planning engine
- a workflow database
- a second authority system
- a second task system
- a policy engine inside Conductor
- a new UI
- a new infrastructure role
- a workflow event store

### 1.1 Relationship to Workflow Specification v0.3

Workflow Design v1.0 **supersedes only the workflow-mode and autonomy-mode portions** of Workflow Specification v0.3.

Workflow Specification v0.3 remains the parent normative source for:

- role ownership and semantic boundaries
- exact operation binding
- fail-closed execution
- authorization decision semantics
- execution outcome semantics
- evidence and correlation
- replay/idempotency
- long-running authority models
- reconciliation
- declaration lifecycle
- conformance principles
- Growth Gate

This document does not silently replace those semantics. Where this document changes the workflow representation, the v1.0 rule is the applicable workflow design and the retained v0.3 rules remain normative for the underlying authority/execution contract.

---

## 2. Role Ownership

| Role | Owns |
|---|---|
| **Human** | plan approval and plan iteration |
| **Agent** | intelligence, planning, decomposition, work, interpretation |
| **Conductor** | durable project state, work coordination, lifecycle, dependencies, claims, history |
| **Declaration / Integration Owner** | effect-capability classification record for declared operation classes |
| **Solvent** | consequential authorization |
| **Executor** | external effect and execution outcome |
| **External SOR** | whether an external effect occurred |
| **Domain** | whether the resulting outcome is correct / acceptable |

The declaration/integration owner is **not** a new infrastructure role. It is the owner of an integration/capability declaration.

Conductor does not:

- plan
- decompose
- authorize
- execute
- determine domain truth
- classify domain consequences
- become a policy router

---

## 3. Locked Invariants

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
PLAN APPROVAL ≠ AUTHORITY
AUTHORIZATION ≠ EXECUTION
CORRELATION ≠ AUTHORITY
TASK COMPLETION ≠ DOMAIN TRUTH
HUMAN PRESENCE ≠ ALTERNATE AUTHORITY PATH
COORDINATION ACTIVITY ≠ PROOF OF EXTERNAL EXECUTION
DECISION ≠ AUTHORIZATION
```

Approval permits Work Mode. It never substitutes for Solvent authorization.

Conductor coordination records do not become authoritative authorization, execution, or external-effect facts.

---

## 4. Agent Modes

There are exactly two Agent modes:

```text
PLAN
WORK
```

No additional Agent mode is normative.

### 4.1 PLAN

The Agent may:

- inspect
- reason
- decompose
- use subagents
- iterate
- revise
- propose

The Agent presents the plan to the Human.

The Human may:

```text
APPROVE
REJECT
```

Plan iteration remains a Human/Agent concern.

Conductor does not orchestrate or manage plan iteration.

### 4.2 WORK

After approval, the Agent performs work against the approved project state.

The Agent may make tactical adjustments within approved scope.

A structural plan change returns to PLAN and requires a new Human approval.

The tactical/structural distinction is advisory. External-effect enforcement does not depend on the Agent classifying its own behavior correctly.

---

## 5. Plan Record

A plan is durable project state.

Minimum plan record:

```text
plan_id
version
scope
approval decision
approved_by
approved_at
```

### 5.1 Plan Scope

Plan scope is expressed at **operation-class grain**, not exact operation identity.

Example:

```text
DeployWorkflow
target scope:
  repository = pithomlabs/reference-loop-test
```

Exact operation identity is constructed at runtime from the effective declaration and effect-relevant inputs.

Plan scope identifiers and declaration-defined operation-class identifiers MUST use the same declared vocabulary.

Every task created under an approved plan MUST carry:

```text
plan_id
plan_version
```

Consequential work records/outcomes MUST retain the same plan linkage for later conformance and audit.

Plan/operation mismatch is observable in v1.0.

Automatic prevention of out-of-plan operations is not part of this design and requires a future Growth Gate decision.

---

## 6. Conductor Project State

Conductor is the durable control plane for project coordination.

It stores:

```text
current project state
durable coordination history
plans and decisions
tasks
dependencies
claims
activity
provenance
artifact/evidence references
```

Large raw corpora and external artifacts remain in their owning systems.

### 6.1 History Fence

Conductor history contains attributed coordination facts and references to externally owned evidence.

It MUST NOT become a shadow authority/execution ledger.

In particular:

```text
DECISION ≠ AUTHORIZATION
COORDINATION RECORD ≠ EXECUTION PROOF
```

History is an append-oriented coordination record with current-state projection. It is not a new workflow event store.

---

## 7. Minimal Task Model

```text
Task
├── task_id
├── objective
├── capability_ref
├── inputs
├── outputs
├── verification_ref
├── dependencies
├── claim
├── lifecycle
├── plan_id
├── plan_version
├── created_by
└── created_at
```

### 7.1 Semantics

- `objective` — Agent-authored description of the work.
- `capability_ref` — reference to the applicable capability/integration declaration.
- `inputs` — required inputs/references.
- `outputs` — expected outputs/artifact references.
- `verification_ref` — reference to the authority/system that decides acceptance.
- `dependencies` — hard prerequisites.
- `claim` — exclusive coordination ownership.
- `lifecycle` — Conductor-only work lifecycle.
- `plan_id`, `plan_version` — approved plan linkage.
- `created_by`, `created_at` — provenance.

Task text is data, not executable instruction.

Task creation policy for untrusted or arbitrary agents is a deployment decision; it is not inferred as a new Conductor authority model.

---

## 8. Capability Declaration

Effect-capable integrations MUST expose a versioned declaration.

At minimum, the declaration identifies:

```text
operation / operation class
effect-capable classification
required authorization boundary
authorization evidence/reference requirements
operation-identity definition
authority-validity model
replay/idempotency behavior
execution outcome model
reconciliation model / owner
trust basis for authorization evidence
UNKNOWN / unresolved-wait age policy
human intervention capability, where applicable
persistence / recovery expectation, where applicable
```

### 8.1 Resolution and Integrity

A `capability_ref` used by a task MUST resolve to:

```text
declaration owner
declaration version
effective reference
content hash
retrievable declaration content
```

The declaration content MUST be integrity-pinned by version + content hash.

A task whose capability reference cannot resolve to an effective declaration MUST NOT become actionable.

The declaration is the authoritative classification record. The declaration owner owns that classification.

The declaration does not transfer authority ownership from Solvent.

### 8.2 Declaration Changes

A declaration change that affects any of the following requires re-conformance of the affected boundary:

```text
effect-capable operation set
operation identity
authorization evidence requirements
authority-validity model
replay/idempotency behavior
execution outcome model
reconciliation behavior
human intervention capability where material
```

Prior conformance is invalid for the affected boundary until required re-conformance completes.

---

## 9. Consequential Operation Definition

> **An operation is consequential when the operation class it instantiates is declared effect-capable—capable of producing an observable state change outside the workflow—regardless of agent intent, workflow phase, or invoking client.**

Classification is at **operation-class grain**.

Examples:

```text
shell / python capability
  local-file operation class
      → non-effect

shell / python capability
  network / push / deploy operation class
      → effect-capable
```

A universal capability is not itself necessarily consequential.

Agent/client classification is advisory.

Conductor does not classify consequences.

The declaration/integration contract deterministically evaluates operation-class membership.

---

## 10. Consequential Enforcement

Any effect-capable execution boundary MUST enforce:

```text
exact requested operation
+
valid Solvent authorization evidence
```

before producing the external effect.

The boundary MUST reject when authorization evidence is:

```text
absent
invalid
stale
expired
mismatched
unverifiable
```

This enforcement is independent of whether the Agent requested authorization.

The enforcement boundary belongs to the effect-capable deployment/integration path, not to Conductor.

An unmediated external-effect path is a conformance finding.

A sandbox, dry run, reversible target, or equivalent safety mechanism MUST NOT bypass:

- authorization verification
- exact operation binding
- fail-closed enforcement

---

## 11. Consequential Flow

```text
Agent
  ↓
exact operation proposal
  ↓
Solvent authorization
  ↓
authorization evidence
  ↓
Executor
  ↓
external effect
  ↓
result / outcome
```

Conductor records only coordination-relevant state and references.

Conductor does not route, authorize, or execute the consequence path.

---

## 12. Ordinary Work

Ordinary work follows:

```text
claim
→ work
→ submit
→ verification / acceptance
```

No Solvent involvement is required unless an effect-capable operation is reached.

An ordinary task may contain many internal operations; the consequential boundary is determined at the effect-capable operation/class boundary.

---

## 13. Exact Operation Identity

The same declared operation-identity definition MUST be used at:

```text
proposal
authorization
execution
```

Operation identity MUST be total over all effect-relevant parameters.

For each effect-relevant parameter, exactly one is true:

```text
included in operation identity
OR
explicitly declared immaterial
```

Undeclared effect-relevant differences constitute a different operation.

The Executor MUST reject:

```text
authorized X
requested Y
```

when:

```text
X ≠ Y
```

Correlation identifiers do not establish authorization binding.

Operation-identity construction and equality/canonicalization MUST be deterministic across implementations.

---

## 14. Authorization and Acceptance

### 14.1 Authorization

Owned by Solvent.

Solvent determines whether the exact consequential operation is authorized.

### 14.2 Acceptance

Owned by the applicable verification, domain, or human authority.

Conductor records:

```text
ACCEPT
REJECT
```

but does not generate the decision.

Acceptance/rejection records MUST include decider attribution.

Conductor MUST only record the decision when the decider is consistent with `verification_ref`.

This prevents Agent self-assertion from becoming Conductor authority.

---

## 15. Conductor Lifecycle

Conductor owns only work lifecycle.

Minimal lifecycle:

```text
OPEN
ACTIVE
BLOCKED
SUBMITTED
DONE
REJECTED
```

Authorization, execution, and external-effect states are not Conductor lifecycle states.

---

## 16. READY Frontier

`READY` is derived, not stored.

Conceptually:

```text
READY =
  eligible lifecycle
  ∧ dependencies satisfied
  ∧ unclaimed
  ∧ unblocked
  ∧ effective capability declaration resolvable
```

Agents discover READY work through Conductor.

Claim is atomic and exclusive.

No scheduler is introduced.

### 16.1 Blocked

A blocked task MUST reference the relevant blocker:

```text
task dependency
OR
declared external/waiting reference
```

No second blocking model is introduced.

---

## 17. Provenance

Keep only the minimum provenance relationships needed for coordination and durable project history:

```text
parent / child
discovered-from
```

Dependency/blocking semantics remain scheduling-relevant.

`discovered-from` is historical provenance and does not itself block readiness.

---

## 18. Multi-Agent Participation

Any compatible Agent may join later:

```text
Agent X
Agent Y
Agent Z
```

All interact through the shared Conductor project state.

Agents do not require inherited context windows.

A new Agent can reconstruct the project from:

```text
current state
+
durable history
+
plans
+
tasks
+
dependencies
+
READY frontier
```

Conductor does not understand the Agent's internal:

- Plan mode implementation
- Work mode implementation
- subagent mechanism
- model
- prompt strategy

---

## 19. Human Plan Approval

Human approval means:

```text
approved to enter/use WORK mode within the approved plan scope
```

It does not mean:

```text
authorized to perform any consequential operation
```

Every consequential operation remains subject to Solvent.

A materially changed plan becomes a new plan version and requires new Human approval.

Human presence never creates a second authority path.

---

## 20. Tactical vs Structural Change

This is an Agent-side advisory distinction.

```text
tactical adjustment
  → may continue in WORK

structural plan change
  → return to PLAN
  → new plan version
  → Human approval
```

Misclassification does not weaken Solvent enforcement.

---

## 21. Protocol Surface

### 21.1 Agent ↔ Conductor

```text
CREATE
DISCOVER
CLAIM
RELEASE
UPDATE
BLOCK
SUBMIT
ACCEPT
REJECT
```

### 21.2 Agent ↔ Solvent / Executor

Consequence operations are not Conductor verbs.

Conceptually:

```text
PROPOSE
AUTHORIZE
EXECUTE
RESULT
```

These belong to their owning boundaries.

---

## 22. Deferred

Explicitly deferred:

```text
claim TTL / heartbeat
automatic plan-scope enforcement
proof-token subsystem
Conductor-side policy engine
cryptographic Conductor↔Solvent ledger
workflow templates / molecules
scheduler
knowledge graph
DAG runtime
```

These require evidence and a Growth Gate before addition.

---

## 23. Known Open Items

These are deliberately carried forward rather than hidden by the freeze:

```text
operation-identity canonicalization across implementations
first real capability declarations
real enforcement-boundary integrity in Agent runtimes
Phase 1.5 negative battery:
  - operation mismatch
  - duplicate delivery
  - stale/invalid authorization
  - UNKNOWN/unavailable
  - stale-claim × expired-intent

Phase 2:
  - pre-registered pass criteria
  - pre-registered failure/falsification criteria
  - pinned declaration
  - pinned components
  - real external-effect confirmation

Further conformance items:
  - concurrency/interference between concurrent consequential proposals
  - client-internal state under orchestration substitution
  - cached-authorization execution under UNKNOWN
  - production observation persistence / evidence decay
  - production revision-loop bound ownership
```

---

## 24. Disposition Log

| Finding | Disposition | Reason | Source |
|---|---|---|---|
| Task-level ordinary/consequential kind | Rejected | Classification belongs to operation class | Workflow design reviews |
| Conductor policy routing | Rejected | Conductor remains coordination-only | Workflow design reviews |
| Formula/molecule workflow machinery | Deferred/rejected for v1.0 | Agent owns planning; no evidence requires templates | Beads review |
| Dolt synchronization | Rejected for v1.0 | Conductor is the shared control plane | Beads review |
| Hash task IDs | Rejected for v1.0 | Existing semantic identity model sufficient | Beads review |
| Claim TTL / heartbeat | Deferred | Liveness issue not yet experimentally earned | Adversarial reviews |
| Proof-token subsystem | Rejected for v1.0 | Would turn Conductor into verification infrastructure | Adversarial reviews |
| Exact-operation inventory in plans | Rejected/replaced | Operation-class scope + runtime exact identity is sufficient | Workflow design reviews |
| Automatic plan-scope enforcement | Deferred | Detection first; prevention requires Growth Gate | Workflow design reviews |
| Durable project state + history | Adopted | Required for cross-session and multi-Agent continuity | Design decision |
| READY frontier | Adopted | Minimal claimable frontier; derived state | Beads review |
| Atomic claim | Adopted | Prevents concurrent claim collision | Beads review |
| Parent/child provenance | Adopted | Minimal work structure | Beads review |
| Discovered-from provenance | Adopted | Durable discovery history without blocking | Beads review |
| Approval ≠ authority | Adopted | Human plan approval cannot replace Solvent | Workflow design reviews |
| Operation-class consequentiality | Adopted | Prevents universal-tool grain error | Workflow design reviews |
| Declaration version + content hash | Adopted | Integrity of effect classification record | Workflow design reviews |
| Acceptance ownership | Adopted | Conductor records decisions; does not create them | Workflow design reviews |
| History fence | Adopted | Prevents shadow authority/execution ledger | Workflow design reviews |
| Attribution | Adopted | Required for multi-Agent durable history | Workflow design reviews |
| Tool-boundary enforcement | Adopted | External-effect paths must fail closed | Workflow design reviews |
| Stale-claim × expired-intent scenario | Deferred to Phase 1.5 | Semantic interaction must be tested before liveness machinery | Workflow design reviews |
| Autonomous / semi-autonomous mode taxonomy | Replaced | Replaced by exactly two Agent modes | Workflow design reviews |
| Conductor plan iteration | Rejected | Human/Agent concern | Workflow design reviews |

---

## 25. Freeze Boundary

This document becomes the normative Workflow Design v1.0 only after:

1. one final consistency review;
2. explicit reconciliation with Workflow Specification v0.3;
3. completion of the disposition log;
4. a separate freeze decision record naming owner, date, and criteria.

After freeze, design changes require the Growth Gate.

---

## 26. Phase 1.5 Gate

Run exactly the pre-agreed battery:

```text
1. operation X authorized → operation Y attempted
2. duplicate delivery
3. stale / invalid authorization
4. UNKNOWN / unavailable
5. stale claim × expired intent
```

Also verify:

```text
declaration resolution
operation-identity equality
Solvent freeze integrity
```

Do not add architecture merely because a test is difficult.

---

## 27. Phase 2 Gate

Before the first real external effect:

```text
pinned declaration
pinned declaration version + content hash
pinned component versions
pre-registered pass criteria
pre-registered failure criteria
pre-registered falsification conditions
verified effect-enforcement path
verified external source of record
```

The Phase 2 experiment MUST establish:

```text
approved plan scope
  ↓
exact operation
  ↓
Solvent authorization
  ↓
effect-boundary verification
  ↓
Executor
  ↓
real external effect
  ↓
external SOR
  ↓
authoritative observation
```

A sandbox may neutralize the final effect, but MUST NOT bypass:

```text
authorization verification
exact operation binding
fail-closed enforcement
```

---

## 28. Growth Gate

If implementation reveals a missing capability, record:

1. observed limitation
2. evidence
3. required behavior
4. why existing capability is insufficient
5. why documentation or an adapter/integration cannot solve it
6. proposed change
7. decision
8. owner/date

A new primitive is not justified merely by convenience.

A proposed change MUST NOT merely relocate prohibited workflow functionality into Conductor, Solvent, Executor, an integration, or another component under a different name.

---

## 29. Design Lock Statement

> **The Agent has exactly two modes: PLAN and WORK. Humans approve plans and plan iteration. Conductor owns durable project state and work coordination. Declaration owners classify operation classes through effective declarations. Solvent authorizes exact consequential operations. The Executor produces external effects.**

Conductor remains a control plane, not an intelligence layer, planner, policy engine, authority system, execution broker, or conformance harness.

---

## 30. Freeze Decision Record Template

```text
WORKFLOW DESIGN FREEZE

Artifact:
    Loop Engineering Workflow Design v1.0

Parent:
    Loop Engineering Workflow Specification v0.3

Owner:
    ____________________

Date:
    ____________________

Criteria:
    [ ] final consistency review complete
    [ ] parent lineage reconciled
    [ ] disposition log complete
    [ ] known open items recorded
    [ ] Phase 1 evidence incorporated
    [ ] no unresolved design-layer objection

Decision:
    FROZEN / NOT FROZEN

Notes:
    ____________________
```
