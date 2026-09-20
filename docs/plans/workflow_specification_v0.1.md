# Loop Engineering Workflow Specification v0.1

## 1. Purpose

This specification defines the implementation-independent workflow contract connecting Agent, Conductor, Solvent, Executor, Human, and the external Domain.

It operationalizes Loop Engineering Requirements v0.3.

It does not create a workflow runtime, workflow database, second authority system, second task system, or new UI.

## 2. Architectural Model

Infrastructure roles:

- Agent — agency, reasoning, proposal, work, interpretation
- Conductor — coordination and work lifecycle
- Solvent — authority decision and authorization
- Executor — external effect

Human is a participant/intervener.

Domain is an external semantic authority for domain meaning, truth, and acceptance.

Core invariant:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Additional invariants:

- Authorization ≠ execution.
- Coordination activity ≠ proof of external execution.
- Task completion ≠ domain truth.
- Human presence ≠ alternate authority path.
- Correlation ≠ authority.
- Workflow phase ≠ persisted authoritative state.

## 3. Workflow Vocabulary

The canonical vocabulary is:

DISCOVER
→ FORMULATE
→ ASSIGN / CLAIM
→ WORK
→ REVIEW
→ NEXT WORK

A consequential action creates an optional branch from WORK:

WORK
→ CONSEQUENTIAL PROPOSAL
→ AUTHORIZATION BOUNDARY
→ EXECUTION
→ RESULT
→ INTERPRET
→ REVIEW
→ NEXT WORK

These are conceptual workflow phases. They are not additional lifecycle/state authorities.

A phase may be represented by existing Conductor lifecycle/activity, Solvent records, Executor records, agent messages, artifacts, or domain records, depending on ownership.

## 4. General Loop Contract

### 4.1 Discover

The Agent discovers available work through the applicable client/interface.

Expected result:

- candidate work identified
- sufficient context available to formulate an action

Conductor remains authoritative for task/work coordination.

### 4.2 Formulate

The Agent interprets the available work and formulates a proposed course of action.

A formulation may result in:

- ordinary work
- a consequential proposal
- a need for human review/intervention
- a request for additional information

The Agent's formulation does not create authority.

### 4.3 Assign / Claim

Where Conductor assignment/claim semantics apply:

- work is assigned or claimed through Conductor
- Conductor remains authoritative for coordination state
- Agent capability does not imply assignment
- assignment does not imply authority

### 4.4 Work

The Agent performs work using ordinary capabilities and applicable tools.

Ordinary work SHOULD proceed without Solvent involvement when no consequential external effect is proposed.

Work may create:

- artifacts
- intermediate evidence
- computation results
- proposals
- review requests

These are owned by the relevant existing system/domain and do not create a workflow-owned state model.

## 5. Consequential Branch

### 5.1 Proposal

The Agent proposes an action that may create an external effect.

The Agent's classification of consequence is advisory.

### 5.2 Effect-Capable Boundary

Before an effect-capable execution path produces an external effect, it MUST require the authorization evidence defined by its integration contract.

If required authorization evidence is absent or invalid:

- execution MUST be rejected
- no external effect may be produced

The execution integration defines how authorization evidence is checked.

It MUST NOT redefine Solvent's ownership of authority or create an alternative authority system.

### 5.3 Authorization

Solvent determines whether the proposed consequential action is authorized according to its own authority model.

Loop Engineering does not reproduce Solvent policy, target/state semantics, revocation semantics, or authority ledgers.

The workflow records/observes the cross-role interaction; Solvent remains authoritative for the authority decision.

### 5.4 Execution

After valid authorization evidence is available, the Executor may attempt the external effect.

Executor responsibilities:

- perform the external operation
- enforce required supplied authorization evidence
- expose execution outcome/evidence
- expose interruption/control capabilities when supported

Executor MUST NOT create authority or determine authorization policy.

### 5.5 Result

Execution result is distinct from authorization.

Possible execution outcome categories:

- not attempted
- attempted
- succeeded
- failed
- ambiguous / unknown

An ambiguous outcome means the occurrence of an external effect cannot be conclusively established.

Conductor activity alone is never sufficient evidence for claiming that the external effect occurred.

The Executor or external system of record is authoritative for execution occurrence.

## 6. Denial and Revision Loop

A denied consequential proposal does not terminate Loop Engineering by default.

Canonical path:

Agent proposes
→ Solvent denies
→ Agent receives denial outcome
→ Agent revises/reformulates
→ new proposal
→ authorization boundary
→ execution if authorized
→ result
→ review

The original denial remains distinct from the revised proposal.

No Conductor authority over Solvent revocation is implied by task cancellation, rejection, or redirection.

## 7. Human Participation

Human participation is cross-cutting, not a separate workflow phase.

At each defined boundary, the workflow contract specifies:

- whether human intervention is possible
- which intervention operations are supported
- what information/evidence is exposed
- what component owns the resulting state

Intervention vocabulary:

- review
- reject
- revise
- redirect
- cancel-before-effect
- cancel-in-flight-when-supported

`pause` is not assumed unless an actual component exposes it.

### 7.1 Before work

Human may:

- approve/reject a proposed work direction
- redirect the work
- request revision
- stop progression

### 7.2 During coordination

Human may, where supported:

- review assignment/claim
- redirect or reassign coordination
- reject proposed coordination changes

Conductor remains the coordination authority.

### 7.3 Before consequential authorization

Human may review or modify the proposal before it reaches the authority boundary.

If the human is exercising actual authority, that authority is evaluated through the applicable Solvent boundary.

Human presence never creates a second authority path.

### 7.4 After authorization / before execution

Where the integration exposes a controllable boundary, human intervention may prevent or alter progression before the external effect.

Any resulting authorization-state change remains owned by Solvent or the applicable authority system.

### 7.5 During execution

Only intervention operations explicitly supported by the Executor may be used.

Examples:

- cancel in flight
- request safe termination
- inspect progress

Absence of an interruption mechanism MUST be represented explicitly.

### 7.6 After execution

Human may:

- inspect execution evidence
- challenge the result
- request further work
- initiate domain review

Already-produced external effects cannot be retroactively negated by workflow intervention.

### 7.7 Review

Human/domain review may:

- accept
- reject
- request revision
- request additional evidence
- redirect next work

Review outcome is not automatically an authority decision.

## 8. Long-Running Execution

A long-running external operation MUST NOT prevent unrelated Conductor work from progressing.

The workflow does not prescribe scheduling, polling, queues, retries, or a workflow runtime.

Each long-running effect MUST expose a declared authority-validity model at the execution boundary.

Permitted abstract models include:

1. Authorization applies to the execution attempt once started.
2. Executor performs defined authorization re-checks during execution.

The selected model belongs to the Executor/integration and applicable authority model.

Workflow semantics distinguish:

authorization granted
→ execution started
→ execution completed/terminated
→ outcome
→ review

## 9. Replay and Idempotency

Effect-capable integrations MUST define behavior for duplicate delivery of the same authorized execution request.

The integration may use:

- idempotency key
- single-use authorization reference
- Executor deduplication
- external transaction identifier
- equivalent mechanism

Loop Engineering does not mandate a mechanism.

Repeated delivery MUST NOT silently create an unintended second external effect when the integration's defined replay policy forbids it.

## 10. Evidence and Correlation

A conformance harness must be able to establish:

- which proposal/request occurred
- which authorization decision/reference applied
- which execution attempt occurred
- which result/outcome corresponds to it
- relevant ordering/relationship

Preferred existing/native references include:

- task identifier
- Action Intent identifier
- authorization reference
- executor execution identifier
- external operation identifier

A conformance-run correlation identifier MAY be introduced for test purposes.

Correlation is an association mechanism.

Correlation MUST NOT be treated as:

- authority
- truth
- semantic ownership
- proof of execution by itself

No workflow-owned event ledger is required.

## 11. Persistence and Ownership

The Workflow Specification does not define a workflow database.

Persistence location means:

> the existing owning system that records the fact.

Examples:

- Conductor — task/assignment/activity/coordination records
- Solvent — authority/authorization records
- Executor/external system — execution records
- Domain — domain evidence/acceptance records

Persistence location does not imply that Loop Engineering owns the semantics of the persisted fact.

## 12. Failure Semantics

At each boundary, failures remain attributable to the owning responsibility.

Examples:

- coordination failure
- authorization denial
- authorization unavailable/unknown
- execution failure
- execution ambiguous/unknown
- domain rejection
- human rejection
- client/tool failure

A failure MUST NOT be reinterpreted as success merely because a later workflow phase requires progress.

A failure MAY produce a revision/retry/new proposal according to the owning system's semantics.

## 13. Autonomous and Semi-Autonomous Modes

The same workflow contract applies to both modes.

### Autonomous

Agent progresses through permitted transitions without required human intervention.

Human intervention may still occur where external requirements or explicit policy require it.

### Semi-Autonomous

Agent progresses through permitted transitions while defined human review/intervention boundaries are available and may stop, redirect, revise, or approve continuation.

Changing mode MUST NOT change role ownership or authority semantics.

## 14. Tool / Runtime Substitution

A client may be:

- autonomous Agent
- alternate Agent
- deterministic scripted client
- human-operated client
- client running under Temporal or another orchestration framework

Substitution MUST preserve:

- workflow vocabulary
- role ownership
- Conductor semantics
- Solvent semantics
- Executor contract
- human intervention semantics
- evidence/correlation expectations

Only the implementation of the client or orchestration mechanism changes.

## 15. Workflow-to-System Mapping

The following mapping is normative at the ownership level:

| Workflow concept | Existing owner / representation |
|---|---|
| Discover / formulate | Agent/client behavior |
| Assign / claim | Conductor |
| Coordination lifecycle | Conductor |
| Consequential authorization | Solvent |
| External effect | Executor |
| Execution evidence | Executor/external system of record |
| Domain truth/acceptance | Domain authority |
| Human intervention | Human acting through applicable owning boundary |
| Cross-role observation | Existing records/interfaces |
| Correlation | Existing/native identifiers where available |

No row creates a new infrastructure owner.

## 16. Conformance Observation Contract

Each conformance scenario defines:

- setup
- actors
- inputs
- steps
- expected evidence
- expected correlation relationships
- pass conditions
- failure conditions

The harness verifies cross-role interaction semantics, not internal component correctness.

Examples:

### Ordinary work

Expected evidence:

- work identified
- work coordinated through Conductor
- artifact/result produced
- review recorded where applicable
- no unnecessary Solvent call

### Consequential work

Expected evidence:

- proposal
- authorization request/decision
- valid/invalid authorization evidence
- execution attempt or rejection
- execution outcome
- correlation across these facts

### Denial/revision

Expected evidence:

- initial proposal
- denial
- revised proposal
- separate authorization decision
- subsequent execution/outcome when authorized

### Human intervention

Expected evidence:

- intervention point
- human action
- resulting owning-system state/evidence
- continuation/termination outcome

## 17. Specification Versioning

Core architectural invariants are stable design constraints.

Normative semantic changes require a specification version increment and re-running relevant conformance scenarios.

Additive clarifications that do not alter semantics may be documented without changing the major contract.

This is v0.1; no mature long-term compatibility regime is implied.

## 18. Growth Gate

If implementation reveals a missing capability, record:

1. observed limitation
2. evidence
3. required behavior
4. why existing capability is insufficient
5. why documentation or adapter/integration cannot solve it
6. proposed change
7. decision
8. owner/date

A new primitive is not justified merely by convenience.

## 19. Initial Conformance Scenarios

The initial suite SHALL implement:

A. Ordinary Research
- no consequential effect
- Conductor coordination
- ordinary Agent work
- result/review

B. Consequential Shared Computation
- Agent proposal
- Solvent authorization
- Executor execution
- result and correlation

C. Human Intervention
- defined intervention boundary
- human review/intervention
- continued/revised/rejected work

D. Tool Substitution
- same consequential scenario
- LLM-driven client
- alternate client
- deterministic/scripted client
- optionally human-operated client
- unchanged workflow semantics

E. Denial → Revision
- proposal
- Solvent denial
- Agent revision
- resubmission
- authorization
- execution if authorized
- review

These scenarios are validation cases, not a replacement for domain-specific BM-IST scientific validation.

## 20. Deliverable Boundary

This specification is complete enough to derive:

- Role / Boundary Matrix v0.1
- Conformance Test Matrix v0.1
- BM-IST Validation Scenarios v0.1
- Agent Skill
- Conformance Harness
- POC Runbook

It does not authorize implementation of a new workflow runtime or infrastructure component.
