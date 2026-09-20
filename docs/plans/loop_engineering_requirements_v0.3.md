# Loop Engineering Requirements v0.3

## 1. Purpose

Loop Engineering defines the cross-role interaction contract for systems composed of:

- Agent — agency, reasoning, proposal, interpretation
- Conductor — coordination and work lifecycle
- Solvent — authority and authorization
- Executor — external effect
- Human — participant, reviewer, and intervener
- Domain — external semantic authority for domain meaning, truth, and acceptance

The specification makes the workflow domain-agnostic, tool/runtime-agnostic, observable, correlatable, and testable without introducing another workflow runtime or authority system.

## 2. Scope

Loop Engineering SHALL define:

- cross-role handoffs
- workflow vocabulary and phase semantics
- consequential-action classification and enforcement expectations
- human intervention and review semantics
- failure, denial, ambiguous-outcome, and recovery expectations
- replay/idempotency expectations at effect-capable boundaries
- asynchronous external-effect expectations
- evidence and correlation requirements for conformance
- tool/runtime-agnostic conformance requirements
- specification versioning and Growth Gate expectations

Loop Engineering SHALL NOT define:

- a new workflow engine
- a new orchestration runtime
- a new authorization system
- a second task system
- a workflow-specific database
- a workflow event store
- a workflow-specific UI
- a persisted workflow `Intent`
- domain-specific truth or policy
- replacement semantics for Conductor, Solvent, or Executor
- mandatory Temporal, LLM, or other execution technology

## 3. Operating Modes

### R3.1 Autonomous

An autonomous operation allows the Agent to progress through permitted workflow boundaries without requiring human intervention, subject to all applicable coordination, authority, execution, and domain-review requirements.

Autonomous does NOT mean unrestricted authority.

### R3.2 Semi-Autonomous

A semi-autonomous operation allows the Agent to progress through permitted workflow boundaries while providing defined opportunities for human review or intervention.

The workflow MUST define, for each boundary, whether human intervention is possible and what operations are supported.

Autonomous and semi-autonomous operation MUST use the same workflow semantics and role boundaries. They are behavioral operating modes, not separate runtime architectures.

## 4. Core Role Ownership

| Concern | Agent | Human | Conductor | Solvent | Executor |
|---|---|---|---|---|---|
| Agency / reasoning | Owns | May intervene | — | — | — |
| Coordination / work lifecycle | Participates | May intervene | Owns | — | — |
| Authority decision | Requests | May request/review/act through applicable authority boundary | — | Owns | Requires applicable authorization |
| External effect | Proposes | May intervene through supported boundary | — | Authorizes | Owns |
| Domain interpretation | Produces / assists | Reviews / verifies | — | — | Reports effect |
| Coordination activity | Produces / participates | Reviews | Records | — | Reports execution outcome |

Human presence MUST NOT create a second authority path.

Solvent remains the authority owner. Executor produces the external effect and enforces the supplied execution constraint/authorization evidence.

### 4.1 Domain ownership

Domain is NOT a seventh infrastructure role.

Domain is an external semantic authority. It may be represented by a human reviewer, domain-specific verifier, domain system, scientific method, external service, or other appropriate mechanism.

Domain acceptance and truth remain outside infrastructure authority. Infrastructure MAY coordinate and record relevant facts but MUST NOT become the authority for domain truth.

## 5. Workflow Vocabulary

The workflow MAY use conceptual phases such as:

DISCOVER → FORMULATE → ASSIGN / CLAIM → WORK → REVIEW → NEXT WORK

A consequential proposal MAY branch from WORK through the applicable authority and execution boundaries:

WORK → CONSEQUENTIAL PROPOSAL → SOLVENT → EXECUTOR → RESULT → INTERPRET → REVIEW

These phases are specification-level vocabulary only.

They MUST NOT become:

- a Conductor lifecycle state
- a Solvent authority state
- a new persisted workflow state
- an independently authoritative event history

No `workflow_state`, workflow state machine, workflow event store, or workflow authority ledger SHALL be introduced unless a later Growth Gate explicitly demonstrates that existing systems cannot satisfy a proven requirement.

## 6. Boundary Catalog

The Workflow Specification SHALL describe each defined boundary using:

- boundary
- trigger
- owner
- human intervention capability
- persistence location, if any
- evidence available to a conformance harness

The catalog is descriptive and non-authoritative.

“Persistence location” identifies where an existing owning system records a relevant fact. It MUST NOT imply semantic ownership and MUST NOT authorize creation of workflow-specific persistence.

For domain acceptance, the owner is the domain verifier/human or other domain authority. Infrastructure only coordinates/records.

## 7. Consequential Classification and Fail-Closed Enforcement

An Agent MAY classify a proposed action as ordinary or consequential, but that classification is advisory.

Any execution path capable of producing an external effect MUST fail closed unless the applicable authorization evidence/reference required by its integration contract is present and valid.

Conceptually:

Agent proposal
→ effect-capable execution boundary
→ authorization evidence required
→ absent/invalid evidence
→ reject without external effect

Loop Engineering SHALL NOT define a workflow-level classifier or duplicate Solvent policy semantics.

The effect-capable integration determines what operations constitute externally consequential effects for that execution boundary and what authorization evidence is required.

## 8. Evidence and Correlation

Loop Engineering is testable through observable cross-role evidence.

Cross-role consequential interactions MUST expose enough evidence for a conformance harness to determine:

- which proposal/request occurred
- which authorization decision/reference applied
- which execution attempt occurred
- which result/outcome corresponds to that attempt
- the relationship and ordering among these facts

Consequential interactions MUST carry sufficient correlation references to associate:

proposal/request
→ authorization
→ execution attempt
→ resulting outcome

The specification SHOULD reuse existing/native identifiers where available, such as:

- Conductor task identifier
- Solvent action-intent identifier
- authorization reference
- executor execution identifier
- external operation identifier

A conformance-run correlation reference MAY be added for test harness purposes.

Loop Engineering MUST NOT introduce a workflow-owned event ledger merely to provide correlation.

## 9. Human Intervention and Review

Human intervention is a first-class workflow capability.

The workflow MUST define whether and how human intervention is possible at each defined boundary.

Where the underlying component exposes no intervention mechanism, that limitation MUST be represented explicitly.

Potential intervention boundaries include:

- before work
- during coordination
- before consequential authorization
- after authorization but before execution, where supported
- during execution, where the Executor exposes a control boundary
- after execution
- during interpretation/review
- before next work

The workflow vocabulary MAY include:

- review
- reject
- revise
- redirect
- cancel-before-effect
- cancel-in-flight-when-supported

`pause` MUST NOT be assumed as a generic capability unless an actual component exposes it.

Human review, human intervention, and human authorization are distinct concepts.

A coordination intervention MUST NOT by itself override an established authority decision or already-produced external effect.

No automatic requirement exists for Conductor cancellation to revoke Solvent authority. Any future coordinated cancellation contract MUST be defined at the integration boundary without transferring authority ownership to Conductor.

## 10. Long-Running External Operations

A long-running external operation MUST NOT prevent unrelated Conductor work from progressing.

This is an observable coordination requirement. It does NOT require Conductor to own scheduling, polling, retries, execution orchestration, or workflow runtime semantics.

Long-running execution MAY be owned by an Executor or external execution framework.

Each long-running effect MUST have a declared authority-validity model at its execution boundary.

Examples include:

- authorize → start → authorization applies to the execution attempt
- authorize → start → executor performs defined authority re-checks

Loop Engineering defines the distinction, not the implementation mechanism.

Authorization granted
≠
execution started
≠
execution completed
≠
result accepted

## 11. Replay and Idempotency

Effect-capable integrations MUST define their replay/idempotency behavior for repeated delivery of the same authorized execution request.

The specification does not mandate a particular mechanism.

Possible mechanisms include:

- idempotency key
- single-use authorization reference
- executor deduplication
- external transaction key
- equivalent mechanism

The mechanism remains owned by the applicable Executor/integration and SHALL NOT create workflow-specific authority state.

## 12. Execution Outcomes

Execution outcomes MUST distinguish at least:

- not attempted
- attempted
- succeeded
- failed
- ambiguous / unknown

An Executor MAY report an ambiguous outcome when external effect occurrence cannot be conclusively established.

An ambiguous outcome MUST remain distinct from both execution success and execution failure.

The Executor or external system of record is authoritative for whether the external effect occurred. Conductor activity alone is never sufficient evidence of external execution.

## 13. Cross-Role Invariants

The following are hard invariants:

1. Capability ≠ Work ≠ Authority ≠ Execution.
2. Authorization ≠ Execution.
3. Coordination activity ≠ proof of external execution.
4. Task completion ≠ domain truth.
5. Human presence ≠ alternate authority engine.
6. Workflow vocabulary ≠ persisted authoritative state.
7. Domain meaning remains outside infrastructure.
8. Conductor does not authorize.
9. Solvent does not execute.
10. Executor does not become the authority engine.
11. Agent capability does not imply authority.
12. Agent classification of consequence is advisory.
13. Effect-capable execution fails closed without required authorization evidence.
14. Ambiguous execution outcome ≠ success and ≠ failure.
15. Human intervention cannot retroactively negate an already-produced external effect.

## 14. Tool / Runtime Agnosticism

Tool/runtime agnosticism SHALL be defined behaviorally.

A conforming client MAY be:

- an autonomous agent
- an alternate agent
- a deterministic/scripted client
- a human-operated client
- a client running under another orchestration framework

Replacing the client MUST NOT require changes to:

- Conductor semantics
- Solvent semantics
- Executor contract
- workflow-phase semantics
- role ownership

The implementation technology may change; the cross-role contract may not.

Provider names such as GPT or Claude SHALL NOT appear in normative protocol conformance definitions.

## 15. Conformance Scope

Loop Engineering conformance tests interaction contracts and cross-role behavior.

Component-specific correctness remains the responsibility of each component's own test suite.

Therefore:

- Loop Engineering tests Agent ↔ Conductor ↔ Solvent ↔ Executor handoffs.
- Conductor tests verify Conductor correctness.
- Solvent tests verify Solvent correctness.
- Executor tests verify execution correctness.
- Domain validation verifies domain truth and acceptance.

Loop Engineering SHALL NOT duplicate internal component test suites.

Behavioral adherence metrics for stochastic agents are advisory and non-gating for v0.x protocol conformance.

## 16. Minimum Conformance Requirements

The first conformance suite SHALL establish at least:

- ordinary work can proceed without Solvent when no consequential action is involved
- Agent classification alone cannot bypass an effect-capable authorization boundary
- missing/invalid authorization evidence prevents the corresponding external effect
- denied authorization prevents the corresponding Executor action
- authorization and execution remain distinct
- execution failure remains distinct from authorization failure
- ambiguous execution remains distinct from success and failure
- long-running external work does not block unrelated Conductor progress
- human intervention behavior matches the declared capability of the boundary
- human review does not create a second authority path
- workflow phases do not become persisted workflow state
- cross-role consequential interactions are observable and correlatable
- repeated authorized requests have defined replay/idempotency behavior
- alternate client implementations can use the same contract
- component-specific correctness remains outside Loop Engineering conformance

## 17. Initial BM-IST Validation Scenarios

### Scenario A — Ordinary Research

Agent → Conductor → local research/computation → artifact → review

No Solvent involvement unless a consequential action is actually proposed.

### Scenario B — Consequential Shared Computation

Agent → Conductor → consequential proposal → Solvent → Executor → result → interpretation → review

This validates authority/effect separation, evidence, and correlation.

### Scenario C — Human Intervention

Agent working → defined human intervention boundary → human redirects/rejects/revises → work continues → review

The exact intervention operation MUST match the capabilities of the existing component.

### Scenario D — Tool Substitution

Run Scenario B with:

- an LLM-driven client
- an alternate agent client
- a deterministic scripted client
- optionally a human-operated client

The semantics, role ownership, and cross-role contracts MUST remain unchanged.

### Scenario E — Denial → Revision

Agent proposes consequential action
→ Solvent denies
→ Agent revises proposal
→ resubmits
→ authorization decision
→ execution if authorized
→ result → review

This validates the loop rather than a one-way pipeline.

## 18. Scenario Test Structure

Every conformance scenario SHALL define:

- setup
- actors
- inputs
- steps
- expected observations/evidence
- correlation expectations
- pass conditions
- failure conditions

Deterministic protocol conformance SHALL use scripted/deterministic clients.

LLM behavior may be measured separately as advisory behavioral adherence.

## 19. Growth Gate

Before introducing any new infrastructure primitive, record:

1. observed POC limitation
2. evidence
3. proposed change
4. why documentation, existing capability, or an adapter/integration boundary cannot solve it
5. decision
6. owner/date

No new primitive is justified merely because it is more convenient.

The Growth Gate is the same discipline applied to Loop Engineering itself.

## 20. Specification Versioning

Core architectural invariants are intended to remain stable.

Normative semantic changes require a specification version increment and re-running conformance.

Additive clarifications that do not change semantics may be documented without changing the major contract.

Because this is v0.x, no mature long-term compatibility regime is implied.

## 21. Acceptance Criteria

Requirements v0.3 is satisfied when:

1. Agent, Conductor, Solvent, Executor, Human, and Domain ownership are unambiguous.
2. Domain is explicitly outside infrastructure authority.
3. Autonomous and semi-autonomous operation use the same contract.
4. Every defined boundary declares its human intervention capability or explicitly states that none exists.
5. Consequential effect-capable paths fail closed without required authorization evidence.
6. Cross-role consequential interactions expose sufficient observable, correlatable evidence for conformance.
7. Proposal, authorization, execution, and outcome can be distinguished.
8. Ambiguous execution outcomes are explicitly representable.
9. Replay/idempotency behavior is defined at effect-capable boundaries.
10. Workflow phases are explicitly non-persistent and non-authoritative.
11. Tool/runtime replacement does not require semantic changes to Conductor, Solvent, Executor, workflow-phase semantics, or role ownership.
12. Loop conformance is limited to cross-role interaction behavior and does not duplicate component correctness suites.
13. Human intervention/review does not create a second authority path.
14. The five BM-IST validation scenarios can be expressed using the baseline systems.
15. Baseline implementation is explicitly:
    - existing Conductor
    - existing Solvent
    - a selected Executor reference implementation
16. No new runtime, database, UI, SDK, scheduler, heartbeat engine, or authority control loop is required by the requirements.
17. Any proposed new primitive must pass the Growth Gate based on an observed POC limitation.

## 22. Required Deliverables After Requirements v0.3

Once this document is accepted, proceed to:

1. Workflow Specification v0.1
2. Role / Boundary Matrix v0.1
3. Conformance Test Matrix v0.1
4. BM-IST Validation Scenarios v0.1
5. POC Runbook

Do not implement a new runtime or infrastructure component before these artifacts expose a demonstrated need.
