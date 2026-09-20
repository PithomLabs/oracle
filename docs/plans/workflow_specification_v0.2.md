# Loop Engineering Workflow Specification v0.2

## 1. Purpose

This specification defines the implementation-independent workflow contract connecting Agent, Conductor, Solvent, Executor, Human, and the external Domain.

It operationalizes Loop Engineering Requirements v0.3.

It does not create a workflow runtime, workflow database, second authority system, second task system, or new UI.

## 2. Three Contract Layers

This specification explicitly separates three kinds of contracts.

### A. Semantic Contract

Defines what each participant means and owns:

- Agent = agency, reasoning, proposal, work, interpretation
- Conductor = coordination and work lifecycle
- Solvent = authority decision and authorization
- Executor = external effect
- Human = participant, reviewer, and intervener
- Domain = external semantic authority for domain meaning, truth, and acceptance

### B. Boundary Contract

Defines what must happen when information or action crosses role boundaries:

proposal
→ authorization
→ execution
→ result

The boundary contract includes exact operation binding, authorization evidence, execution eligibility, and outcome semantics.

### C. Conformance Evidence Contract

Defines how semantic and boundary behavior can be demonstrated:

- identifiers
- correlation
- observable evidence
- test fixtures/actors
- pass/fail observations

These layers do not create a new state authority.

## 3. Architectural Model

Core invariant:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Additional invariants:

- Authorization ≠ execution.
- Coordination activity ≠ proof of external execution.
- Task completion ≠ domain truth.
- Human presence ≠ alternate authority path.
- Correlation ≠ authority.
- Workflow phase ≠ persisted authoritative state.
- Correlation establishes association; authorization binding establishes authorized operation identity.

## 4. Workflow Vocabulary

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

These are conceptual workflow phases, not additional lifecycle/state authorities.

A phase may be represented through existing Conductor lifecycle/activity, Solvent records, Executor records, agent messages, artifacts, or domain records according to ownership.

## 5. General Loop Contract

### 5.1 Discover

The Agent discovers available work through the applicable client/interface.

Expected result:

- candidate work identified
- sufficient context available to formulate an action

Conductor remains authoritative for task/work coordination.

### 5.2 Formulate

The Agent interprets available work and formulates a proposed course of action.

A formulation may result in:

- ordinary work
- consequential proposal
- request for human review/intervention
- request for additional information

Formulation does not create authority.

### 5.3 Assign / Claim

Where Conductor assignment/claim semantics apply:

- work is assigned or claimed through Conductor
- Conductor remains authoritative for coordination state
- Agent capability does not imply assignment
- assignment does not imply authority

### 5.4 Work

The Agent performs work using ordinary capabilities and applicable tools.

Ordinary work SHOULD proceed without Solvent involvement when no consequential external effect is proposed.

Work may create artifacts, intermediate evidence, computation results, proposals, or review requests.

These remain owned by existing systems/domain and do not create workflow-owned state.

## 6. External Effect Definition and Consequential Operation Declaration

Each effect-capable integration MUST declare which operations can produce an external effect.

The declaration is an integration property, not a new domain-policy engine.

The declaration MUST identify:

- the operation or operation class
- whether the operation is effect-capable
- the required authorization boundary
- the authorization evidence/reference expected by that boundary
- the authority-validity model
- replay/idempotency behavior
- execution outcome model

The declaration MUST NOT transfer authority ownership from Solvent to the integration.

Agent/client classification of consequence remains advisory.

A client/integration MAY identify that an operation requires the consequential boundary.

The final effect-capable enforcement point MUST fail closed.

Conductor does not classify domain consequences and does not become a policy router.

## 7. Consequential Branch

### 7.1 Proposal

The Agent or client proposes an action that may create an external effect.

The proposal MUST identify the operation sufficiently for the authority boundary to determine what exact action is being requested.

### 7.2 Exact Operation Binding

Authorization evidence MUST be bound to the exact operation it authorizes.

The bound operation SHOULD be representable using a canonical operation identity or equivalent stable description containing the relevant target/action parameters.

Conceptually:

Proposal X
→ Solvent authorizes X
→ authorization evidence bound to X
→ Executor receives Y
→ if X ≠ Y, reject

An effect-capable execution boundary MUST reject execution when the requested operation does not correspond to the authorized operation.

Correlation identifiers do not establish operation binding.

### 7.3 Authorization Evidence Verification Contract

The execution integration MUST define how it obtains and verifies authorization evidence.

The contract MUST identify:

- evidence/reference supplied
- authority issuer/owner
- operation identity to which it is bound
- validity conditions applicable to the integration
- required handling of absent, invalid, stale, or mismatched evidence

The integration defines the mechanism of verification.

It MUST NOT create authority or redefine Solvent policy.

### 7.4 Fail-Closed Enforcement

Before an effect-capable execution path produces an external effect, required authorization evidence MUST be present and valid.

If absent, invalid, stale, or mismatched:

- execution MUST be rejected
- no external effect may be produced

### 7.5 Authorization

Solvent determines whether the proposed consequential action is authorized according to its own authority model.

Loop Engineering does not reproduce Solvent policy, target/state semantics, revocation semantics, or authority ledgers.

The workflow records/observes the cross-role interaction.

### 7.6 Execution

After valid authorization evidence is available, the Executor may attempt the external effect.

Executor responsibilities:

- perform the external operation
- reject execution when required authorization evidence is absent or invalid
- enforce the supplied execution constraint
- expose execution outcome/evidence
- expose interruption/control capabilities when supported

Executor MUST NOT determine authorization policy or create authority.

### 7.7 Result

Execution outcome is distinct from authorization.

Possible outcome categories:

- not attempted
- attempted
- succeeded
- failed
- ambiguous / unknown

An ambiguous outcome means external effect occurrence cannot be conclusively established.

The Executor or external system of record is authoritative for execution occurrence.

Conductor activity alone is never sufficient evidence of external execution.

## 8. Authorization Decision Taxonomy

Authorization outcomes MUST distinguish at least:

### DENIED

An authoritative authority decision exists and the requested operation is not authorized.

Expected response:

- do not execute
- agent/client MAY revise/reformulate
- a revised proposal is a distinct proposal instance/version

### UNKNOWN / UNAVAILABLE

No authoritative authorization decision is available or the authority result cannot currently be established.

Expected response:

- do not execute unless an independently valid existing authorization already covers the requested operation
- MUST NOT reinterpret unknown/unavailable as denial
- MAY retry, wait, or escalate according to the integration policy

The workflow specification does not prescribe a retry mechanism.

## 9. Denial → Revision

Canonical path:

Agent proposes V1
→ Solvent denies V1
→ Agent receives denial
→ Agent revises
→ Proposal V2
→ new authorization decision
→ execution if authorized
→ result
→ review

V2 MUST have its own proposal/reference identity.

A revision MUST NOT silently mutate the identity of V1.

Conformance scenarios SHALL use a bounded number of revision iterations.

The bound applies to the test scenario, not as a normative production workflow retry counter.

Denial count is evidence, not workflow authority.

## 10. Human Participation

Human participation is cross-cutting, not a separate workflow phase.

At each defined boundary, the workflow contract specifies:

- whether intervention is possible
- supported intervention operations
- evidence exposed
- owning component for resulting state

Intervention vocabulary:

- review
- reject
- revise
- redirect
- cancel-before-effect
- cancel-in-flight-when-supported

`pause` is not assumed unless an actual component exposes it.

Human review, intervention, and authorization are distinct.

### 10.1 Human-Modifiable Proposal

When a human materially modifies a consequential proposal:

Proposal V1
→ human revision
→ Proposal V2
→ new authorization decision

V2 MUST carry its own proposal/reference identity.

Human modification does not inherit V1 authorization automatically.

### 10.2 Human Authority

If a human is exercising actual authority, the action remains subject to the applicable Solvent authority boundary.

Human presence never creates a second authority path.

### 10.3 Coordination Intervention

A coordination intervention does not automatically revoke or mutate Solvent authority.

Conductor cancellation/redirection is a coordination fact.

Solvent authorization/revocation remains a Solvent authority fact.

Any future coordinated cancellation contract MUST be defined at the integration boundary without transferring authority ownership to Conductor.

## 11. Long-Running Execution

A long-running external operation MUST NOT prevent unrelated Conductor work from progressing.

Each long-running effect MUST declare:

- authority-validity model
- revalidation model, if any
- mid-flight revocation behavior, if any

Permitted abstract models include:

1. authorize → start → authorization applies to the execution attempt
2. authorize → start → executor performs defined authority re-checks

The selected model belongs to the Executor/integration and applicable authority model.

The specification does not mandate heartbeats or periodic reauthorization.

If relevant authority/state changes during execution, the integration MUST follow its declared model. Workflow semantics do not silently assume that execution remains valid forever.

## 12. Replay and Idempotency

Effect-capable integrations MUST define behavior for duplicate delivery of the same authorized execution request.

Possible mechanisms include:

- idempotency key
- single-use authorization reference
- executor deduplication
- external transaction identifier
- equivalent mechanism

Repeated delivery MUST follow the declared integration policy and MUST NOT silently create an unintended duplicate external effect where the policy forbids it.

## 13. Execution Reconciliation

An ambiguous execution outcome MUST enter the Executor/integration's declared reconciliation path.

Reconciliation MAY use:

- status query
- human escalation
- idempotent retry
- external reconciliation
- manual investigation
- equivalent mechanism

Loop Engineering defines the distinction and requirement for reconciliation, but does not own the reconciliation engine.

Ambiguous MUST NOT automatically become success or failure.

## 14. Evidence and Correlation

A conforming deployment MUST expose sufficient evidence through existing owning interfaces or explicit test adapters for the conformance harness to establish:

- proposal/request
- authorization decision/reference
- execution attempt
- result/outcome
- relationship/order
- operation identity binding

Possible evidence sources:

- Conductor API/activity
- Solvent API/authority record
- Executor test interface/records
- external system of record
- explicit test adapter

No workflow event store is required.

### 14.1 Correlation

Correlation identifiers associate related records.

Examples:

- task identifier
- Action Intent identifier
- authorization reference
- execution identifier
- external operation identifier
- conformance-run correlation reference

Correlation MUST NOT be treated as:

- authority
- truth
- semantic ownership
- proof of execution by itself

### 14.2 Authorization Binding

Authorization binding separately establishes whether the specific execution request is the exact operation authorized.

Correlation and authorization binding are independent concepts.

## 15. Conformance Observation Interface

A conformance deployment MUST expose the evidence required by its test scenarios through:

- existing owning interfaces, OR
- explicit test-only adapters.

The observation mechanism MUST permit deterministic inspection of:

- request/proposal
- authority decision
- execution attempt
- outcome
- correlations
- relevant ordering
- operation binding

The observation interface is read/test infrastructure.

It does NOT become an authoritative workflow ledger.

## 16. Conformance Test Actors

Conformance tests MAY use:

- real components
- deterministic test doubles
- controlled authority fixtures
- fault-injecting test executors
- deterministic/scripted clients

Test fixtures may deliberately produce:

### TestSolvent

- authorize X
- deny X
- return unknown/unavailable
- return stale/invalid authorization evidence
- authorize X while execution requests Y

### TestExecutor

- accept valid authorization for exact operation
- reject missing evidence
- reject invalid evidence
- reject stale evidence
- reject operation mismatch
- return execution failure
- return ambiguous outcome
- simulate duplicate delivery
- expose or omit interruption boundary

These are test fixtures only. They do not define or replace production component semantics.

## 17. Failure Semantics

Failures remain attributable to the owning responsibility.

Examples:

- coordination failure
- authorization denial
- authorization unknown/unavailable
- execution rejection
- execution failure
- execution ambiguous/unknown
- domain rejection
- human rejection
- client/tool failure

A failure MUST NOT be reinterpreted as success merely because a later phase expects progress.

Clients MUST have access to the state/evidence necessary to continue through the owning interfaces.

How the client obtains it is implementation-specific:

- polling
- MCP request
- HTTP request
- subscription
- external orchestration activity
- human observation

No mandatory webhook, event bus, or reconciliation runtime is implied.

## 18. Autonomous and Semi-Autonomous Modes

The same contract applies to both modes.

### Autonomous

Agent progresses through permitted transitions without required human intervention.

External requirements, authority decisions, or explicit review requirements may still require non-agent participation.

### Semi-Autonomous

Agent progresses through permitted transitions while defined human review/intervention boundaries are available and may stop, redirect, revise, or approve continuation.

Changing mode MUST NOT change role ownership or authority semantics.

## 19. Tool and Orchestration Substitution

Two substitution axes are explicitly distinguished.

### 19.1 Client substitution

Examples:

- GPT-driven client
- Claude-driven client
- alternate agent
- deterministic scripted client
- human-operated client

### 19.2 Orchestration substitution

Examples:

- direct client/API interaction
- Temporal-driven client
- another orchestration framework

Substitution along either axis MUST preserve:

- workflow vocabulary
- role ownership
- Conductor semantics
- Solvent semantics
- Executor contract
- human intervention semantics
- evidence/correlation expectations
- exact operation/authorization binding

The implementation technology may change; the cross-role contract may not.

## 20. Workflow-to-System Mapping

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
| Cross-role observation | Existing owning interfaces / test adapters |
| Correlation | Existing/native identifiers where available |

No row creates a new infrastructure owner.

## 21. Conformance Contract Types

Each conformance scenario SHALL distinguish:

### Semantic contract tests

Verify that roles retain their intended ownership and do not cross boundaries.

### Boundary contract tests

Verify required handoff behavior, including:

- consequential routing/declaration
- exact operation binding
- authorization evidence
- fail-closed enforcement
- execution outcome semantics

### Conformance evidence tests

Verify that required evidence, correlation, and operation binding can be observed deterministically.

These are interaction tests, not internal component correctness tests.

## 22. Scenario Test Structure

Every conformance scenario SHALL define:

- setup
- actors
- inputs
- steps
- expected evidence
- expected correlations
- operation/authorization binding expectations
- pass conditions
- failure conditions

Deterministic protocol conformance SHALL use deterministic clients/fixtures.

LLM behavior may be measured separately as advisory behavioral adherence.

## 23. Initial Validation Scenarios

### A — Ordinary Research

Agent → Conductor → local research/computation → artifact → review

No Solvent involvement unless a consequential effect is proposed.

### B — Consequential Shared Computation

Agent → Conductor → declared consequential operation → Solvent → Executor → result → review

Tests:

- exact operation binding
- fail-closed enforcement
- observation/correlation
- execution outcome

### C — Human Intervention

Agent working → defined intervention boundary → human redirects/rejects/revises → work continues → review

### D — Tool Substitution

Run Scenario B with:

- LLM-driven client
- alternate agent client
- deterministic scripted client
- optionally human-operated client

### D2 — Orchestration Substitution

Run Scenario B through:

- direct client/API path
- Temporal-driven client
- optionally another orchestration framework

Workflow semantics MUST remain unchanged.

### E — Denial → Revision

Proposal V1
→ Solvent denial
→ Proposal V2
→ new authorization decision
→ execution if authorized
→ result → review

The conformance test uses a bounded revision loop.

### F — Unknown / Reconciliation

Proposal
→ authorization unavailable/unknown OR execution ambiguous
→ declared reconciliation path
→ subsequent authoritative determination
→ continue/revise/review

## 24. Agent Skill

The Agent Skill is a client-side behavioral instruction artifact derived from this specification.

It MAY contain:

- workflow rules
- decision guidance
- handoff conventions
- proposal format
- evidence expectations
- examples

It is NOT:

- a runtime
- an SDK
- an authority component
- a workflow server
- a required infrastructure dependency

The Agent Skill teaches an Agent how to participate in the contract; it does not enforce authority.

## 25. Conformance Harness

The Conformance Harness is test-only infrastructure.

It MAY contain:

- scenario definitions
- deterministic clients
- test doubles
- fault injectors
- observation adapters
- assertions
- reports

It has no production workflow responsibility and does not become a fifth infrastructure role.

## 26. Persistence and Ownership

The Workflow Specification does not define a workflow database.

Persistence location means the existing owning system that records the relevant fact.

Examples:

- Conductor — task/assignment/activity/coordination records
- Solvent — authority/authorization records
- Executor/external system — execution records
- Domain — domain evidence/acceptance records

Persistence location does not imply Loop Engineering ownership of semantics.

## 27. Specification Versioning

Core architectural invariants remain stable design constraints.

Normative semantic changes require a specification version increment and re-running relevant conformance scenarios.

Additive clarifications that do not alter semantics may be documented without changing the major contract.

This is v0.2; no mature long-term compatibility regime is implied.

## 28. Growth Gate

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

## 29. Deliverable Boundary

This specification is complete enough to derive:

- Role / Boundary Matrix v0.1
- Conformance Test Matrix v0.1
- BM-IST Validation Scenarios v0.1
- Agent Skill
- Conformance Harness
- POC Runbook

It does not authorize implementation of a new workflow runtime or infrastructure component.
