# Loop Engineering Role / Boundary Matrix v0.2

## 1. Purpose

This matrix translates Loop Engineering Workflow Specification v0.3 into an explicit ownership and boundary reference.

It is normative for role ownership and boundary separation.

Where this matrix and the Workflow Specification differ, the Workflow Specification is normative; this matrix is its ownership/boundary projection.

It does not create a new state model, workflow runtime, authority system, database, event store, or UI.

## 2. Terminology

### Integration

“Integration” denotes an implementation boundary or adapter belonging to an existing role boundary. It is not an independent Loop Engineering role.

An integration MUST NOT acquire independent coordination, authority, or execution ownership.

An integration MAY be implemented as:

- native wrapper
- adapter
- proxy
- gateway
- sidecar
- library
- direct callback
- equivalent mechanism

These are implementation choices, not additional Loop Engineering roles.

### Domain

Domain is an external semantic authority, not a seventh infrastructure role.

It may be represented by a human reviewer, domain-specific verifier, domain system, scientific method, external service, or other domain-appropriate mechanism.

## 3. Role Model

| Participant | Primary responsibility | Authoritative for | Not authoritative for |
|---|---|---|---|
| Agent | Agency, reasoning, proposal, work, interpretation | Its own reasoning/proposals and produced work artifacts | Coordination state, authority, external effect, domain truth |
| Conductor | Coordination and work lifecycle | Projects/tasks/assignment/claim/activity/lifecycle within its existing model | Authority, external effect, domain truth, Agent reasoning |
| Solvent | Authority decision and authorization | Authorization, authority state, revocation/claim semantics within its existing model | Coordination, execution, domain truth |
| Executor | External effect and execution outcome | Whether/how the external operation was attempted and what actually occurred, subject to external system-of-record semantics | Authority policy, task coordination, domain truth |
| Human | Intervention, review, and optional authority participation through the applicable boundary | Human decision/review within the boundary being exercised | A parallel authority system outside the applicable owner |
| Domain | Domain meaning, truth, acceptance | Domain-specific correctness/acceptance | Infrastructure coordination/authority/execution mechanics |

## 4. Core Ownership Invariants

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
11. Integration is an adapter/boundary, not an independent role.

## 5. Boundary Matrix

| Boundary | Initiating participant | Receiving participant | Purpose | Required evidence | Human intervention | Semantic / authority owner | Enforcement / implementation responsibility |
|---|---|---|---|---|---|---|---|
| Discover | Agent | Conductor/client interface | Identify available work | Task/work reference | Review/redirect where supported | Conductor for coordination | Conductor/client implementation |
| Formulate | Agent | Agent/client context | Formulate proposed action | Proposal/request reference | Review/revise/redirect | Agent/client for proposal | Agent/client implementation |
| Assign / Claim | Agent / Human | Conductor | Establish coordination responsibility | Task + assignment/claim evidence | Reassign/release/reject where supported | Conductor | Conductor implementation |
| Ordinary Work — non-effect capability | Agent | Declared non-effect tool/service | Perform non-consequential work | Artifact/intermediate evidence | Review/redirect/stop where supported | Agent/tool/service as applicable | Tool/service implementation |
| Ordinary Work — effect-capable capability | Agent/client | Effect-capable integration / authority boundary | Invoke an externally consequential operation | Proposal + operation identity | Review/revise/reject before effect where supported | Solvent for authority; Executor/external system for effect | Effect-capable integration + Executor |
| Consequential Proposal | Agent/client | Authority boundary | Request externally consequential action | Proposal + operation identity/version | Review/revise/reject | Solvent for authority | Client/integration implementation |
| Authorization | Solvent | Effect-capable integration / Executor | Decide and expose authority | Authorization evidence + exact operation binding | Human may act through applicable Solvent boundary | Solvent | Integration verifies evidence |
| Execution Eligibility | Effect-capable integration | Executor | Verify required authorization before effect | Authorization evidence + trust basis + operation identity | Intervention if boundary supports it | Solvent | Integration / Executor |
| Execution | Executor | External system/world | Produce external effect | Execution reference + outcome evidence | Cancel/intervene when supported | Executor/external system for effect | Executor/external system |
| Result | Executor / external system | Agent/client / Conductor | Return execution outcome | Execution ID + outcome | Review/challenge/request reconciliation | Executor/external system for outcome | Executor/integration |
| Domain Review | Domain verifier / Human | Domain context | Determine domain acceptance/truth | Domain evidence/review record | Domain-defined | Domain | Domain verifier/system |
| Next Work | Agent / Human / Conductor | Conductor | Continue, revise, redirect, or conclude work | Task/work reference + review outcome | Human may intervene | Conductor for coordination | Conductor implementation |

### Boundary Rule

An operation is consequential because the invoked capability/integration is declared effect-capable, not because Conductor interprets the semantic payload or because a particular workflow phase labels it “consequential.”

Any operation capable of producing an external effect is consequential regardless of workflow phase or invoking client and MUST use the applicable consequential authorization boundary.

Effect-capable integrations are responsible for complete effect-surface declarations.

Loop Engineering does not provide a universal oracle capable of discovering undeclared side effects inside arbitrary third-party software.

## 6. Internal Participant Phases

The following are workflow phases, not external participant boundaries:

- Formulate
- Interpret
- Review preparation
- Agent-side reasoning

They MUST NOT be treated as additional roles or cross-role authority transitions.

## 7. Consequential Boundary Rules

### 7.1 Classification

Agent/client classification of consequence is advisory.

The effect-capable integration declaration defines which operations are effect-capable.

Conductor does not become a domain-specific policy router.

### 7.2 Fail-Closed Enforcement

Every effect-capable integration MUST reject an external effect when required authorization evidence is:

- absent
- invalid
- stale/expired
- operation-mismatched

The Executor MUST NOT decide which operations ought to be authorized.

### 7.3 Exact Operation Binding

The same declared operation-identity definition MUST be used at:

- proposal
- authorization
- execution

All effect-relevant parameters MUST be included in the identity or explicitly declared immaterial.

Operation identity MUST be defined over semantic effect inputs, not incidental transport metadata.

Any excluded field MUST be explicitly declared immaterial to the external effect.

Examples of normally non-semantic transport metadata may include:

- trace_id
- retry_count
- transport timestamp
- HTTP connection metadata

The declaration, not an implicit implementation convention, determines whether a field is immaterial.

Correlation does not establish authorization binding.

### 7.4 Authorization Evidence

Authorization evidence MUST carry or directly and verifiably expose the authorized operation identity/context.

The trust basis MUST be independently verifiable against the authority owner and MUST NOT rely solely on an unverified client claim.

The integration declaration MUST identify how the evidence is authenticated or validated against the authority owner.

Possible mechanisms include:

- signed artifact
- direct authority lookup
- trusted authority-controlled reference
- equivalent mechanism

Cryptography is not mandated at the Loop Engineering layer.

## 8. Ownership vs Enforcement

A semantic/authority owner and an enforcement mechanism may differ.

The distinction is:

> One semantic owner; possibly another enforcement mechanism.

| Concern | Semantic / authority owner | Enforcement / implementation responsibility |
|---|---|---|
| Authorization | Solvent | Executor-side / integration verification |
| Revocation | Solvent | Applicable integration honors declared revocation model |
| Coordination | Conductor | Conductor implementation |
| Execution outcome | Executor / external system of record | Executor reports/reconciles |
| Domain acceptance | Domain | Domain verifier/system |
| Human intervention | Human within applicable boundary | Owning component implements supported controls |

Enforcement responsibility MUST NOT be interpreted as semantic ownership.

## 9. Human Intervention Matrix

| Boundary | Human may inspect | Human may intervene | Limit |
|---|---|---|---|
| Discover | Work/context | Redirect/review | Must use applicable coordination boundary |
| Formulate | Proposal/context | Revise/reject/redirect | Does not create authority |
| Assign/Claim | Assignment/lifecycle | Reassign/release/reject where supported | Conductor remains coordination owner |
| Ordinary non-effect Work | Work/artifacts | Redirect/stop/review where supported | Capability is component-dependent |
| Before Authorization | Proposal + evidence | Revise/reject/redirect | Material change creates distinct proposal/version |
| Authorization | Authority decision | Exercise authority only through applicable Solvent boundary | Human presence does not bypass Solvent |
| Before Execution | Authorization + operation | Cancel/change progression where supported | Cannot invent authority |
| During Execution | Execution status | Cancel/intervene when Executor supports it | Already-produced effects are not retroactively undone |
| After Execution | Execution evidence | Challenge/request reconciliation/review | Executor/external system remains outcome authority |
| Domain Review | Domain evidence | Accept/reject/revise | Domain owns truth/acceptance |

## 10. Authorization UNKNOWN Temporal Policy

Authorization UNKNOWN/UNAVAILABLE is distinct from DENIED.

An applicable integration/authority declaration MUST define a maximum unresolved-wait/age policy.

The policy MUST define one of:

- explicit operator-controlled indefinite wait, or
- a terminal/remediation action such as retry, escalation, revision, or termination through an existing applicable boundary.

The policy MUST be actionable rather than merely documenting “wait.”

Loop Engineering does not introduce a generic `Abandoned`, timeout, escalation, or human-escalation state.

## 11. Long-Running Execution and Authority Precedence

Each long-running effect MUST declare:

- authority-validity model
- validity scope
- revalidation model, if any
- mid-flight revocation behavior, if any

Authority-owner constraints are authoritative.

An integration MAY impose a stricter validity constraint but MUST NOT extend, weaken, or override an authority constraint.

Conceptually:

- Solvent validity 1h + integration validity 24h → effective validity ≤ 1h
- Solvent validity 24h + integration validity 1h → effective validity ≤ 1h

The conformance scenario MUST exercise the declared model.

## 12. Execution Outcomes

Execution outcomes MUST distinguish:

- not attempted
- rejected before effect
- attempted
- succeeded
- failed
- ambiguous / unknown

### Execution Rejected

Means:

- required authorization/evidence validation failed
- the effect-capable boundary prevented the external effect from being attempted

It MUST NOT be reinterpreted as execution failure.

### Execution Failed

Means:

- the effect was attempted
- the operation did not successfully complete

It MUST NOT be reinterpreted as authorization denial.

### Execution Ambiguous

Means:

- effect occurrence cannot be conclusively established

It MUST remain distinct from success and failure.

## 13. Failure Ownership Matrix

| Failure/outcome | Semantic / authority owner | Meaning | Must not be reinterpreted as |
|---|---|---|---|
| Coordination failure | Conductor/integration | Work coordination could not proceed | Authorization denial |
| Authorization denied | Solvent | Authority decision exists and rejects operation | Unknown/unavailable |
| Authorization unknown/unavailable | Solvent/integration | No authoritative decision currently established | Denial |
| Authorization invalid/stale | Solvent/integration boundary | Required authorization evidence cannot be accepted | Authorized |
| Operation mismatch | Effect-capable integration / Executor boundary | Requested operation differs from authorized operation | Authorized |
| Execution rejection | Effect-capable integration | Required authorization/evidence failed; external effect was not attempted | Execution failure or authorization denial |
| Execution failure | Executor/external system | Effect was attempted but did not successfully complete | Authorization denial |
| Execution ambiguous | Executor/integration | Effect occurrence cannot be conclusively established | Success or failure |
| Domain rejection | Domain verifier/Human | Domain acceptance failed | Infrastructure authorization denial |
| Human rejection | Human within applicable boundary | Human declined continuation/review/action | Automatic Solvent revocation |

## 14. Evidence / Correlation Matrix

| Fact | Minimum evidence | Semantic / authority owner | Correlation role |
|---|---|---|---|
| Proposal Vn | Proposal/version reference + operation identity | Agent/client | Starting record for action lineage |
| Human intervention | Human action/reference + affected proposal/work reference | Human within applicable boundary | Links intervention to affected work/proposal |
| Authorization | Decision + authorization reference + bound operation | Solvent | Links proposal to authority |
| Declaration | Owner + version + effective reference + retrievable content | Integration owner | Establishes applicable contract version |
| Execution attempt | Execution reference + requested operation | Executor | Links authorized operation to effect attempt |
| Execution rejection | Rejection reason/evidence + operation reference | Effect-capable boundary | Demonstrates fail-closed behavior |
| Execution outcome | Outcome + execution reference | Executor/external system of record | Links attempt to result |
| Long-running progress | Start/reference + unrelated Conductor progress evidence | Executor + Conductor | Demonstrates non-blocking coordination |
| Reconciliation determination | Reconciliation reference + authoritative determination | Executor/integration | Resolves ambiguous outcome |
| Domain acceptance | Domain review/evidence reference | Domain authority | Links result to domain interpretation |
| Coordination activity | Conductor activity/task evidence | Conductor | Operational history only |

Correlation identifiers establish association.

They do not establish:

- authority
- truth
- semantic ownership
- proof of execution by themselves

Conformance ordering is causal/happens-before ordering, not merely timestamp ordering.

## 15. Persistence / State Ownership

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

## 16. Integration Declaration Boundary

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

## 17. Production Conformance Boundary

Conformance has two layers:

1. Protocol Fixture Conformance
2. Integration Conformance

Integration conformance MUST exercise the real authorization/enforcement path.

A safe target, sandbox, dry-run, or reversible environment may neutralize the final external effect, but it MUST NOT bypass:

- authorization verification
- exact operation binding
- fail-closed enforcement

Passing a TestExecutor/TestSolvent scenario does not prove the production integration behaves correctly.

## 18. Client / Runtime Substitution Boundary

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

## 19. Agent Skill Boundary

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

## 20. Conformance Harness Boundary

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

## 21. Boundary Anti-Patterns

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
- integration enforcement being treated as independent semantic ownership

## 22. Role Substitution Tests

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

## 23. Matrix-Level Acceptance Criteria

The Role / Boundary Matrix is acceptable when:

1. Every responsibility has one semantic/authority owner.
2. No role silently owns another role's semantics.
3. Every consequential boundary identifies the authority owner and effect owner.
4. Human intervention is defined without creating a second authority path.
5. Domain remains outside infrastructure authority.
6. Exact operation binding is distinguishable from correlation.
7. Failure/outcome ownership is explicit, including rejection vs failure.
8. Persistence locations are distinguished from semantic ownership.
9. Production conformance is distinguishable from fixture conformance.
10. Client and orchestration substitution preserve the contract.
11. No matrix row requires a new runtime, database, event store, UI, or authority engine.
12. Effect-capable integration completeness is treated as an explicit integration/deployment responsibility.
13. Any proposed exception can be evaluated through the Growth Gate.

## 24. Adversarial Questions for the Conformance Matrix

Before accepting the Conformance Test Matrix, verify:

1. Can a tool invoked during “ordinary work” produce an external effect without entering the consequential boundary?
2. Can a proposal be authorized for X but executed as Y?
3. Is the same operation-identity definition used at proposal, authorization, and execution?
4. Can an effect-capable integration bypass Solvent?
5. Can a client call ordinary execution when the invoked capability is actually effect-capable?
6. Can UNKNOWN authorization be silently treated as denial or wait indefinitely without a declared policy?
7. Can an ambiguous execution outcome be silently treated as success?
8. Can a fail-closed rejection be confused with execution failure?
9. Can a human intervene without creating shadow authority?
10. Can Conductor cancellation mutate Solvent authority?
11. Can integration declaration drift make an old conformance result appear valid?
12. Can the conformance harness pass while bypassing the real enforcement path?
13. Can an integration hide an undeclared effect-capable operation and still claim conformance?
14. Can a new component absorb workflow semantics merely by renaming them?

If any answer is unclear, the Conformance Test Matrix MUST expose that ambiguity as a test.

## 25. Precedence and Deliverable Boundary

The Workflow Specification v0.3 is the normative parent artifact.

This matrix is its ownership/boundary projection.

This matrix is the basis for:

- Conformance Test Matrix v0.1
- BM-IST Validation Scenarios v0.1
- Agent Skill
- Conformance Harness
- POC Runbook

It does not authorize implementation of a new workflow runtime or infrastructure component.
