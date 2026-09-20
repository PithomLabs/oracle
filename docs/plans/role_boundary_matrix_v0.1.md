# Loop Engineering Role / Boundary Matrix v0.1

## 1. Purpose

This matrix translates Loop Engineering Workflow Specification v0.3 into an explicit ownership and boundary reference.

It is normative for role ownership and boundary separation.

It does not create a new state model, workflow runtime, authority system, database, event store, or UI.

## 2. Role Model

| Participant | Primary responsibility | Authoritative for | Not authoritative for |
|---|---|---|---|
| Agent | Agency, reasoning, proposal, work, interpretation | Its own reasoning/proposals and produced work artifacts | Coordination state, authority, external effect, domain truth |
| Conductor | Coordination and work lifecycle | Projects/tasks/assignment/claim/activity/lifecycle within its existing model | Authority, external effect, domain truth, Agent reasoning |
| Solvent | Authority decision and authorization | Authorization, authority state, revocation/claim semantics within its existing model | Coordination, execution, domain truth |
| Executor | External effect and execution outcome | Whether/how the external operation was attempted and what actually occurred, subject to the external system of record | Authority policy, task coordination, domain truth |
| Human | Intervention, review, domain judgment, optional authority participation through applicable boundary | Human decision/review within the boundary being exercised | A parallel authority system outside the applicable owner |
| Domain | Domain meaning, truth, acceptance | Domain-specific correctness/acceptance | Infrastructure coordination/authority/execution mechanics |

## 3. Core Ownership Invariants

1. Capability ≠ Work ≠ Authority ≠ Execution.
2. Conductor owns coordination; it does not authorize.
3. Solvent owns authority; it does not execute.
4. Executor owns external effect; it does not create authority.
5. Agent exercises agency; capability does not imply authority.
6. Human participation does not create a second authority path.
7. Domain meaning remains outside infrastructure authority.
8. Coordination activity is not proof of external execution.
9. Correlation is not authority, truth, or execution proof.
10. Workflow phases are vocabulary, not an additional authoritative state model.

## 4. Boundary Matrix

| Boundary | Initiating participant | Receiving participant | Purpose | Required evidence | Human intervention | Authority owner | Effect owner |
|---|---|---|---|---|---|---|---|
| Discover | Agent | Conductor/client interface | Identify available work | Task/work reference | Review/redirect where supported | N/A | N/A |
| Formulate | Agent | Agent/client context | Convert work into proposed action | Proposal/request reference | Review/revise/redirect | N/A | N/A |
| Assign / Claim | Agent / Human | Conductor | Establish coordination responsibility | Task + assignment/claim evidence | Reassign/release/reject where supported | N/A | N/A |
| Ordinary Work | Agent | Tools / external services | Perform non-consequential work | Artifact/intermediate evidence | Review/redirect/stop where supported | N/A | Tool/service-specific |
| Consequential Proposal | Agent/client | Authority boundary | Request externally consequential action | Proposal + operation identity | Review/revise/reject | Solvent | Pending |
| Authorization | Solvent | Effect-capable integration / Executor | Decide and expose authority | Authorization evidence + exact operation binding | Human may act through applicable authority boundary | Solvent | Pending |
| Execution Eligibility | Effect-capable integration | Executor | Validate required authorization before effect | Authorization evidence, trust basis, operation identity | Intervention if boundary supports it | Solvent | Executor |
| Execution | Executor | External system/world | Produce external effect | Execution reference + outcome evidence | Cancel/intervene when supported | Solvent remains authority owner | Executor/external system |
| Result | Executor / external system | Agent/client / Conductor | Return execution outcome | Execution ID + outcome | Review/challenge/request further work | Solvent for prior authority fact | Executor/external system |
| Interpretation | Agent | Agent/domain context | Interpret result | Result + supporting evidence | Review/challenge/redirect | Domain for truth | N/A |
| Domain Review | Domain verifier / Human | Domain context | Determine domain acceptance/truth | Domain evidence/review record | Domain-defined | Domain | N/A |
| Next Work | Agent / Human / Conductor | Conductor | Continue, revise, redirect, or conclude work | Task/work reference + review outcome | Human may intervene | N/A | N/A |

## 5. Consequential Boundary Rules

### 5.1 Classification

Agent/client classification of consequence is advisory.

The effect-capable integration declaration defines which operations are effect-capable.

Conductor does not become a domain-specific policy router.

### 5.2 Fail-Closed Enforcement

Every effect-capable integration MUST reject an external effect when required authorization evidence is:

- absent
- invalid
- stale/expired
- operation-mismatched

The Executor MUST NOT decide which operations ought to be authorized.

### 5.3 Exact Operation Binding

The same declared operation-identity definition MUST be used at:

- proposal
- authorization
- execution

All effect-relevant parameters MUST be included in the identity or explicitly declared immaterial.

Correlation does not establish authorization binding.

### 5.4 Authorization Evidence

Authorization evidence MUST carry or directly and verifiably expose the authorized operation identity/context.

The integration declaration MUST identify the trust basis used to authenticate the evidence to the authority owner.

The mechanism may vary by implementation.

## 6. Human Intervention Matrix

| Boundary | Human may inspect | Human may intervene | Limit |
|---|---|---|---|
| Discover | Work/context | Redirect/review | Must use applicable coordination boundary |
| Formulate | Proposal/context | Revise/reject/redirect | Does not create authority |
| Assign/Claim | Assignment/lifecycle | Reassign/release/reject where supported | Conductor remains coordination owner |
| Ordinary Work | Work/artifacts | Redirect/stop/review where supported | Capability is component-dependent |
| Before Authorization | Proposal + evidence | Revise/reject/redirect | Material change creates distinct proposal/version |
| Authorization | Authority decision | Exercise authority only through applicable Solvent boundary | Human presence does not bypass Solvent |
| Before Execution | Authorization + operation | Cancel/change progression where supported | Cannot invent authority |
| During Execution | Execution status | Cancel/intervene when Executor supports it | Already-produced effects are not retroactively undone |
| After Execution | Execution evidence | Challenge/request reconciliation/review | Executor/external system remains execution authority |
| Domain Review | Domain evidence | Accept/reject/revise | Domain owns truth/acceptance |

## 7. Long-Running Execution Boundary

| Concern | Owner | Requirement |
|---|---|---|
| Coordination while operation runs | Conductor | Unrelated work must remain able to progress |
| Authority validity model | Solvent + integration declaration | Declared explicitly |
| Revalidation, if any | Executor/integration | Follow declared model |
| Mid-flight revocation behavior | Solvent + integration contract | Declared where applicable |
| Actual execution | Executor | Owns effect and execution outcome |
| Reconciliation | Executor/integration | Owns ambiguous-effect reconciliation |
| Interpretation | Agent | Interprets available evidence |
| Domain acceptance | Domain/Human | Owns domain meaning/truth |

No global lock or Loop Engineering heartbeat mechanism is implied.

## 8. Failure Ownership Matrix

| Failure/outcome | Primary owner | Meaning | Must not be reinterpreted as |
|---|---|---|---|
| Coordination failure | Conductor/integration | Work coordination could not proceed | Authorization denial |
| Authorization denied | Solvent | Authority decision exists and rejects operation | Unknown/unavailable |
| Authorization unknown/unavailable | Solvent/integration | No authoritative decision currently established | Denial |
| Authorization invalid/stale | Integration/Solvent boundary | Required authorization evidence cannot be accepted | Authorized |
| Operation mismatch | Executor/integration | Requested operation differs from authorized operation | Authorized |
| Execution failure | Executor/external system | Effect attempt did not successfully complete | Authorization denial |
| Execution ambiguous | Executor/integration | Effect occurrence cannot be conclusively established | Success or failure |
| Domain rejection | Domain verifier/Human | Domain acceptance failed | Infrastructure authorization denial |
| Human rejection | Human within applicable boundary | Human declined continuation/review/action | Automatic Solvent revocation |

## 9. Evidence / Correlation Matrix

| Fact | Minimum evidence | Authoritative owner | Correlation role |
|---|---|---|---|
| Proposal | Proposal/request reference + operation identity | Agent/client/domain workflow context | Starting record |
| Authorization | Decision + authorization reference + bound operation | Solvent | Links proposal to authority |
| Execution attempt | Execution reference + requested operation | Executor | Links authorized operation to effect attempt |
| Execution outcome | Outcome + execution reference | Executor/external system of record | Links attempt to result |
| Domain acceptance | Domain review/evidence reference | Domain authority | Links result to domain interpretation |
| Coordination activity | Conductor activity/task evidence | Conductor | Operational history only |

Correlation identifiers establish association.

They do not establish:

- authority
- truth
- semantic ownership
- proof of execution by themselves

Conformance ordering is causal/happens-before ordering, not merely timestamp ordering.

## 10. Persistence / State Ownership

No workflow-owned persistence is introduced.

| Fact | Where it may be recorded | Semantic owner |
|---|---|---|
| Task/lifecycle/assignment/activity | Conductor | Conductor |
| Authority/authorization/revocation | Solvent | Solvent |
| Execution/outcome | Executor/external system | Executor/external system |
| Domain evidence/acceptance | Domain system/record | Domain |
| Agent artifact/reasoning evidence | Agent/task artifact system as applicable | Agent/domain context |
| Conformance observations | Test-only harness/adapter | Conformance test system |

Persistence location does not imply Loop Engineering ownership.

No `workflow_state`, workflow event store, workflow authority ledger, or equivalent shadow state is permitted without a Growth Gate-approved architectural change.

## 11. Integration Declaration Boundary

Each effect-capable integration declaration MUST identify at minimum:

- declaration owner
- declaration version
- effective version/date or equivalent activation reference
- retrievable declaration content
- effect-capable operations
- operation-identity definition
- authorization boundary
- authorization evidence/reference
- trust basis
- authority-validity model
- replay/idempotency behavior
- execution outcome model
- reconciliation owner
- reconciliation mechanism
- persistence/recovery expectation where applicable
- human intervention capabilities where applicable

Declaration changes that alter effect surface or any listed semantic property invalidate prior conformance for the affected boundary until re-conformance completes.

## 12. Production Conformance Boundary

Conformance has two layers:

1. Protocol Fixture Conformance
2. Integration Conformance

Integration conformance MUST exercise the real authorization/enforcement path.

A safe target, sandbox, dry-run, or reversible environment may neutralize the final external effect, but it MUST NOT bypass:

- authorization verification
- exact operation binding
- fail-closed enforcement

Passing a TestExecutor/TestSolvent scenario does not prove the production integration behaves correctly.

## 13. Client / Runtime Substitution Boundary

### Client substitution

Examples:

- GPT-driven client
- Claude-driven client
- alternate Agent
- deterministic script
- human-operated client

### Orchestration substitution

Examples:

- direct API/MCP interaction
- Temporal-driven client
- another orchestration framework

Substitution MUST NOT change:

- semantic role ownership
- workflow vocabulary
- authorization semantics
- Executor contract
- evidence/correlation requirements
- operation binding
- human intervention semantics

The implementation may change; the contract does not.

## 14. Agent Skill Boundary

The Agent Skill is:

- a client-side behavioral instruction artifact
- derived from Loop Engineering Specification
- guidance for participating in the contract

It is NOT:

- an authority component
- a workflow runtime
- an SDK
- a workflow server
- a required infrastructure dependency

## 15. Conformance Harness Boundary

The Conformance Harness is test-only infrastructure.

It may:

- orchestrate test execution
- provide deterministic fixtures
- inject faults
- inspect evidence
- run assertions
- generate reports

Its orchestration has no normative production semantics.

It MUST NOT become:

- a production workflow runtime
- an authority engine
- a workflow state store
- a reference production scheduler

## 16. Boundary Anti-Patterns

The following constitute architecture violations unless explicitly approved by the Growth Gate:

- Conductor authorizing an action
- Solvent executing an action
- Executor creating or deciding authority
- Agent bypassing required effect authorization
- Human bypassing the applicable authority boundary
- correlation ID treated as authority
- activity record treated as execution proof
- workflow phase persisted as authoritative state
- domain truth persisted as infrastructure truth
- Conductor cancellation silently mutating Solvent authority
- sandbox testing bypassing the real authorization/enforcement path
- integration inventing undeclared effect-capable operations
- Executor inventing operation-class semantics

## 17. Role Substitution Tests

The matrix itself SHOULD be challenged with at least these substitutions:

| Substitution | Expected result |
|---|---|
| GPT → Claude | No semantic change |
| Claude → scripted client | No semantic change |
| scripted client → human/curl | No semantic change |
| direct API → Temporal | No semantic change |
| Executor implementation A → B | Boundary declaration changes only where semantics actually differ |
| domain verifier human → domain system | Domain ownership remains external |
| Human present → absent in autonomous mode | Same contract; different required intervention behavior |

## 18. Matrix-Level Acceptance Criteria

The Role / Boundary Matrix is acceptable when:

1. Every responsibility has one authoritative owner.
2. No role silently owns another role's semantics.
3. Every consequential boundary identifies the authority owner and effect owner.
4. Human intervention is defined without creating a second authority path.
5. Domain remains outside infrastructure authority.
6. Exact operation binding is distinguishable from correlation.
7. Failure/outcome ownership is explicit.
8. Persistence locations are distinguished from semantic ownership.
9. Production conformance is distinguishable from fixture conformance.
10. Client and orchestration substitution preserve the contract.
11. No matrix row requires a new runtime, database, event store, UI, or authority engine.
12. Any proposed exception can be evaluated through the Growth Gate.

## 19. Adversarial Questions for the Next Artifact

Before accepting the Conformance Test Matrix, verify:

1. Can a proposal be authorized for X but executed as Y?
2. Can an effect-capable integration bypass Solvent?
3. Can a client call ordinary execution when the operation is actually effect-capable?
4. Can UNKNOWN authorization be silently treated as denial?
5. Can an ambiguous execution outcome be silently treated as success?
6. Can a human intervene without creating shadow authority?
7. Can Conductor cancellation mutate Solvent authority?
8. Can integration declaration drift make an old conformance result appear valid?
9. Can the conformance harness pass while bypassing the real enforcement path?
10. Can a new component absorb workflow semantics merely by renaming them?

If any answer is unclear, the next Conformance Test Matrix must expose that ambiguity as a test.

## 20. Deliverable Boundary

This matrix is the basis for:

- Conformance Test Matrix v0.1
- BM-IST Validation Scenarios v0.1
- Agent Skill
- Conformance Harness
- POC Runbook

It does not authorize implementation of a new workflow runtime or infrastructure component.
