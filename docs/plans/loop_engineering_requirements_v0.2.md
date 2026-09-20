# Loop Engineering Requirements v0.2

## 1. Purpose

Loop Engineering defines the cross-role interaction contract for systems composed of:

- Agent — agency, reasoning, proposal, interpretation
- Conductor — coordination and work lifecycle
- Solvent — authority and authorization
- Executor — external effect
- Human — participant, reviewer, and intervener
- Domain — domain meaning, truth, and acceptance

The specification exists to make this workflow domain-agnostic, tool/runtime-agnostic, observable, and testable without introducing another workflow runtime or authority system.

## 2. Scope

Loop Engineering SHALL define:
- cross-role handoffs
- workflow vocabulary and phase semantics
- consequential-action branching
- human intervention and review semantics
- failure, denial, and recovery expectations
- asynchronous external-effect expectations
- tool/runtime-agnostic conformance requirements
- cross-role invariants

Loop Engineering SHALL NOT define:
- a new workflow engine
- a new orchestration runtime
- a new authorization system
- a second task system
- a workflow-specific database
- a workflow-specific UI
- a persisted workflow `Intent`
- domain-specific truth or policy
- replacement semantics for Conductor, Solvent, or Executor

## 3. Operating Modes

### R3.1 Autonomous

An autonomous operation allows the Agent to progress through permitted workflow boundaries without requiring human intervention, subject to all existing coordination, authority, execution, and domain-review requirements.

Autonomous does NOT mean unrestricted authority.

### R3.2 Semi-Autonomous

A semi-autonomous operation allows the Agent to progress through permitted workflow boundaries while providing defined opportunities for human review or intervention.

The workflow MUST define, for each boundary, whether human intervention is available and what operations are supported.

Autonomous and semi-autonomous operation MUST use the same workflow semantics and role boundaries. They are behavioral operating modes, not separate runtime architectures.

## 4. Core Role Ownership

| Concern | Agent | Human | Conductor | Solvent | Executor |
|---|---|---|---|---|---|
| Agency / reasoning | Owns | May intervene | — | — | — |
| Coordination / work lifecycle | Participates | May intervene | Owns | — | — |
| Authority decision | Requests | May request/review/act through applicable authority boundary | — | Owns | Requires supplied authorization |
| External effect | Proposes | May intervene through supported boundary | — | Authorizes | Owns |
| Domain interpretation | Produces / assists | Reviews / verifies | — | — | Reports effect |
| Coordination activity | Produces / participates | Reviews | Records | — | Reports execution outcome |

Human participation MUST NOT create a second authority path.

Solvent remains the authority owner. Executor produces the external effect.

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

No `workflow_state`, workflow state machine, workflow event store, or workflow authority ledger SHALL be introduced unless a later Growth Gate explicitly demonstrates that an existing system cannot satisfy a proven requirement.

## 6. Boundary Catalog Requirement

The Workflow Specification SHALL describe each defined boundary using:
- boundary
- trigger
- owner
- human intervention capability
- persistence location, if any

The catalog is descriptive and non-authoritative.

Domain acceptance remains owned by the domain verifier/human. Infrastructure may coordinate and record relevant workflow history but SHALL NOT become the authority for domain truth.

## 7. Human Intervention and Review

Human intervention is a first-class workflow capability.

The workflow MUST define whether and how human intervention is possible at each defined boundary.

Where the underlying component exposes no intervention mechanism, that limitation MUST be represented explicitly rather than assumed away.

Potential intervention boundaries include:
- before work
- during coordination
- before consequential authorization
- after authorization but before execution, where supported
- during execution, where the Executor exposes a control boundary
- after execution
- during interpretation/review
- before next work

Human review and human intervention are distinct concepts.

Human authorization is also distinct from both and occurs through the applicable authority boundary.

An intervention cannot retroactively invalidate an external effect that has already occurred.

## 8. Consequential Boundary

An Agent MAY propose a consequential action but MUST NOT establish its own authority to execute it.

Effect-capable execution paths MUST enforce the applicable Solvent authorization before producing the external effect.

Conductor SHALL NOT become a domain-specific policy engine merely to determine authority.

The authority boundary SHALL remain outside Conductor.

## 9. Long-Running External Operations

A long-running external operation MUST NOT prevent unrelated Conductor work from progressing.

This requirement is about observable coordination behavior.

It does NOT require:
- Conductor-owned scheduling
- Conductor-owned polling
- Conductor-owned retries
- Conductor-owned execution orchestration
- a new workflow runtime

Long-running execution MAY be owned by an Executor or an external execution framework.

The workflow MUST distinguish:

authorization granted ≠ execution started ≠ execution completed ≠ execution result accepted

## 10. Cross-Role Invariants

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
12. Agent proposals do not bypass consequential-action controls.

## 11. Tool / Runtime Agnosticism

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

## 12. Conformance Scope

Loop Engineering conformance tests interaction contracts and cross-role behavior.

Component-specific correctness remains the responsibility of each component's own test suite.

Therefore:
- Loop Engineering tests Agent ↔ Conductor ↔ Solvent ↔ Executor handoffs.
- Conductor tests verify Conductor correctness.
- Solvent tests verify Solvent correctness.
- Executor tests verify execution correctness.
- Domain validation verifies domain truth and acceptance.

The Loop Engineering suite SHALL NOT duplicate internal component test suites.

## 13. Minimum Conformance Requirements

The first conformance suite SHALL establish at least:
- ordinary work can proceed without Solvent when no consequential action is involved
- consequential proposals cannot bypass the authority boundary
- denied authorization prevents the corresponding Executor action
- authorization and execution remain distinct
- execution failure remains distinct from authorization failure
- long-running external work does not block unrelated Conductor progress
- human intervention behavior matches the declared capability of the boundary
- human review does not create a second authority path
- workflow phases do not become persisted workflow state
- alternate client implementations can use the same contract

Agent behavioral adherence MAY be measured statistically. It SHALL NOT redefine deterministic protocol conformance.

## 14. Initial BM-IST Validation Scenarios

### Scenario A — Ordinary Research

Agent → Conductor → local research/computation → artifact → review

No Solvent involvement unless a consequential action is actually proposed.

### Scenario B — Consequential Shared Computation

Agent → Conductor → consequential proposal → Solvent → Executor → result → interpretation → review

This validates the authority/effect boundary.

### Scenario C — Human Intervention

Agent working → defined human intervention boundary → human redirects/rejects/revises → work continues → review

The exact intervention operation MUST match the capabilities of the existing Conductor/component involved.

## 15. Non-Goals

The first implementation SHALL NOT introduce:
- Loop Engine runtime
- workflow server
- workflow database
- workflow-specific frontend
- workflow SDK
- provider-specific agent integrations
- duplicate Solvent state
- duplicate Conductor lifecycle state
- mandatory Temporal dependency
- mandatory LLM dependency

## 16. Acceptance Criteria

Requirements v0.2 is satisfied when:

1. The role boundaries are unambiguous.
2. Autonomous and semi-autonomous operation use the same contract.
3. Every defined workflow boundary declares its human intervention capability or explicitly states that none exists.
4. Consequential actions cannot bypass the applicable authority boundary.
5. Workflow phases are explicitly non-persistent and non-authoritative.
6. Tool/runtime replacement does not require semantic changes to Conductor, Solvent, Executor, or role ownership.
7. Loop conformance is limited to cross-role interaction behavior.
8. Domain acceptance remains outside infrastructure authority.
9. The three initial BM-IST scenarios can be expressed using the existing systems without introducing new infrastructure.
10. Any proposed new primitive must pass the Growth Gate based on an observed POC limitation.

## 17. Growth Gate

No new infrastructure primitive is authorized merely because it would make the workflow more convenient.

Before introducing any new primitive, demonstrate:
- an existing component cannot express the proven requirement;
- an adapter/integration/documentation boundary is insufficient;
- the deficiency appears in an actual POC scenario;
- the proposed addition does not duplicate another role's authority or state.

Only then may implementation scope expand.
