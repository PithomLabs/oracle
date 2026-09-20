# Loop Engineering Workflow Design v1.0

**Status:** Proposed freeze artifact  
**Parent:** Loop Engineering Workflow Specification v0.3  
**Purpose:** Minimal workflow/control-plane design derived from the Phase 1 Reference Loop and subsequent review.

---

## 1. Scope

This document defines the minimal workflow contract for:

- Agent Plan/Work participation
- Human plan approval
- Conductor project/work coordination
- consequential-operation handoff to Solvent
- Executor effect production
- durable project state and history

It does not create a workflow runtime, planning engine, second authority system, event store, or new infrastructure role.

### Lineage

This document supersedes the workflow-mode and autonomy-mode portions of Workflow Specification v0.3.

Workflow Specification v0.3 remains the parent source for authority, exact operation binding, fail-closed enforcement, authorization/result semantics, evidence, replay/idempotency, reconciliation, declarations, and conformance principles.

---

## 2. Role Ownership

| Role | Owns |
|---|---|
| Human | plan approval and plan iteration |
| Agent | intelligence, planning, decomposition, work, interpretation |
| Conductor | durable project state, work coordination, lifecycle, dependencies, claims, history |
| Declaration Owner | operation-class/effect classification record |
| Solvent | consequential authorization |
| Executor | external effect and execution outcome |
| External/Domain authority | external truth and domain acceptance |

Conductor does not plan, authorize, execute, or interpret domain truth.

---

## 3. Locked Invariants

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
PLAN APPROVAL ≠ AUTHORITY
AUTHORIZATION ≠ EXECUTION
CORRELATION ≠ AUTHORITY
TASK COMPLETION ≠ DOMAIN TRUTH
HUMAN PRESENCE ≠ ALTERNATE AUTHORITY PATH
```

Approval permits Work Mode. It never substitutes for Solvent authorization.

---

## 4. Agent Modes

There are exactly two Agent modes:

```text
PLAN
WORK
```

### PLAN

The Agent may:

- inspect
- reason
- decompose
- use subagents
- iterate
- revise
- propose

The Agent presents the plan to the Human.

The Human may approve or reject it.

Plan iteration is a Human/Agent concern. Conductor does not orchestrate it.

### WORK

After approval, the Agent performs work against the approved project state.

The Agent may make tactical adjustments within the approved work.

A structural plan change returns to PLAN and requires new Human approval.

The distinction is advisory; external-effect enforcement does not depend on the Agent classifying its own behavior correctly.

---

## 5. Plan Record

A plan is durable project state.

Minimum:

```text
plan_id
version
scope
approval decision
approved_by
approved_at
```

### Plan scope

Scope is expressed at **operation-class grain**, not exact operation identity.

Example:

```text
DeployWorkflow
 target scope:
   repository = pithomlabs/reference-loop-test
```

Exact operation identity is constructed at runtime from the effective declaration and effect-relevant inputs.

Plan scope and declaration class identifiers MUST use the same declared vocabulary.

The system records `plan_id`/version against consequential work and outcomes for later conformance and audit.

Plan/operation mismatch is observable in v1.0; automatic prevention is not part of this design.

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

### History fence

Conductor history contains attributed coordination facts and references to externally owned evidence.

It MUST NOT become a shadow authority/execution ledger.

```text
DECISION ≠ AUTHORIZATION
```

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
├── created_by
└── created_at
```

### Semantics

- `objective` — Agent-authored description of the work.
- `capability_ref` — reference to the capability/integration contract.
- `inputs` — required inputs/references.
- `outputs` — expected outputs/artifact references.
- `verification_ref` — reference to the authority/system that decides acceptance.
- `dependencies` — hard prerequisites.
- `claim` — exclusive coordination ownership.
- `lifecycle` — Conductor-only work lifecycle.
- `created_by`, `created_at` — provenance.

Task text is data, not executable instruction.

---

## 8. Capability Declaration

Effect-capable integrations MUST expose a versioned declaration.

The declaration identifies, at minimum:

```text
operation / operation class
effect-capable classification
required authorization boundary
operation-identity definition
authorization evidence requirements
authority-validity model
replay/idempotency behavior
execution outcome model
reconciliation model/owner
```

### Resolution

A `capability_ref` used by a task MUST resolve to:

```text
declaration owner
declaration version
content hash
effective reference
retrievable declaration content
```

The declaration content MUST be integrity-pinned by version + content hash.

A task whose capability reference cannot resolve to an effective declaration MUST NOT become actionable.

The declaration is the authoritative classification record; the declaration owner owns that classification.

---

## 9. Consequential Operation Definition

> **An operation is consequential when the operation class it instantiates is declared effect-capable—capable of producing an observable state change outside the workflow—regardless of agent intent, workflow phase, or invoking client.**

The classification is at **operation-class grain**.

Examples:

```text
shell / python
  local-file operation → non-effect class
  network/push/deploy operation → effect-capable class
```

A universal capability is not itself necessarily consequential.

Agent classification is advisory.

Conductor does not classify consequences.

---

## 10. Consequential Enforcement

Any effect-capable execution boundary MUST enforce:

```text
exact operation
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

Unmediated external-effect paths are a conformance finding.

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

Conductor does not route or authorize the consequence path.

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

---

## 13. Exact Operation Identity

The operation identity definition MUST be shared across:

```text
proposal
authorization
execution
```

It MUST be total over all effect-relevant parameters.

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

Operation identity equality/canonicalization MUST be deterministic across implementations.

---

## 14. Authorization vs Acceptance

### Authorization

Owned by Solvent.

### Acceptance

Owned by the applicable verification/domain/human authority.

Conductor records:

```text
ACCEPT
REJECT
```

but does not generate the decision.

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

---

## 17. Provenance

Keep only the minimum provenance relationships needed for coordination/history:

```text
parent/child
discovered-from
```

`blocks`/dependency semantics remain scheduling-relevant.

`discovered-from` is historical provenance and does not itself block readiness.

---

## 18. Multi-Agent Participation

Any compatible Agent may join later:

```text
Agent X
Agent Y
Agent Z
```

All interact through Conductor's shared project state.

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

Conductor does not understand the Agent's internal Plan/Agent mode, subagent mechanism, or model.

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

---

## 20. Tactical vs Structural Change

This is an Agent-side advisory distinction.

```text
tactical adjustment
  → may continue in WORK

structural plan change
  → return to PLAN
  → Human approval
```

Misclassification does not weaken Solvent enforcement.

---

## 21. Protocol Surface

### Agent ↔ Conductor

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

### Agent ↔ Solvent / Executor

Consequence operations are not Conductor verbs:

```text
PROPOSE
AUTHORIZE
EXECUTE
RESULT
```

Those belong to their owning boundaries.

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

These require evidence/Growth Gate justification before addition.

---

## 23. Known Open Items

These are deliberately carried forward rather than hidden by the freeze:

```text
operation-identity canonicalization across implementations
declarations for the first real capabilities
real enforcement-boundary integrity in Agent runtimes
Phase 1.5 negative battery:
  - operation mismatch
  - duplicate delivery
  - stale/invalid authorization
  - UNKNOWN/unavailable
  - stale-claim × expired-intent
Phase 2 pre-registered pass/fail/falsification criteria
real external-effect confirmation
```

---

## 24. Disposition Log

| Finding | Disposition | Reason | Source |
|---|---|---|---|
| Task-level ordinary/consequential kind | Rejected | Classification belongs to operation class | Consolidated review |
| Conductor policy routing | Rejected | Conductor remains coordination-only | Consolidated review |
| Formula/molecule workflow machinery | Deferred/rejected for v1.0 | Agent owns planning; no evidence requires templates | Beads review |
| Dolt synchronization | Rejected for v1.0 | Conductor is the shared control plane | Beads review |
| Hash task IDs | Rejected for v1.0 | Existing semantic identity model sufficient | Beads review |
| Claim TTL / heartbeat | Deferred | Liveness issue not yet experimentally earned | Adversarial review |
| Proof-token subsystem | Rejected for v1.0 | Would turn Conductor into verification infrastructure | Adversarial review |
| Exact operation inventory in plans | Deferred/replaced | Operation-class scope + runtime exact identity is sufficient | Workflow review |
| Automatic plan-scope enforcement | Deferred | Detection first; prevention requires Growth Gate | Workflow review |
| Durable project state + history | Adopted | Required for cross-session/multi-Agent continuity | Design decision |
| READY frontier | Adopted | Minimal claimable frontier; derived state | Beads review |
| Atomic claim | Adopted | Prevents concurrent claim collision | Beads review |
| parent/child provenance | Adopted | Minimal work structure | Beads review |
| discovered-from provenance | Adopted | Durable discovery history without blocking | Beads review |

---

## 25. Freeze Boundary

This document is intended to become the normative Workflow Design v1.0 only after:

1. One final consistency review.
2. Workflow Specification v0.3 lineage is explicitly reconciled.
3. Disposition log is complete.
4. Freeze decision is recorded separately.

After freeze, design changes require the Growth Gate.

---

## 26. Phase Gate After Freeze

### Phase 1.5

Run:

```text
operation X authorized → operation Y attempted
duplicate delivery
stale/invalid authorization
UNKNOWN / unavailable
stale claim × expired intent
```

Also verify:

```text
declaration resolution
operation identity equality
Solvent freeze integrity
```

### Phase 2

Before the first real external effect:

```text
pinned declaration
pinned component versions
pre-registered pass criteria
pre-registered failure/falsification criteria
verified enforcement path
```

Then perform one reversible/sandboxed real external-effect experiment.

No additional architecture is added merely because Phase 2 is harder.

---

## 27. Design Lock Statement

> **The Agent has exactly two modes: PLAN and WORK. Humans approve plans and plan iteration. Conductor owns durable project state and work coordination. The declaration owner defines operation-class effect capability. Solvent authorizes exact consequential operations. The Executor produces external effects.**

Conductor remains a control plane, not an intelligence layer, planner, policy engine, authority system, execution broker, or harness.
