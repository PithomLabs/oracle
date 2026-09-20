# Loop Engineering Role / Boundary Matrix v0.4

## 1. Purpose

This matrix translates Loop Engineering Workflow Specification v0.3 into an explicit ownership and boundary reference.

It is the normative ownership/boundary projection of the Workflow Specification.

Where this matrix refines the parent specification, the refinement is permitted. Where a conflict appears, the Workflow Specification v0.3 is normative and this matrix MUST NOT contradict it.

This matrix does not create a new state model, workflow runtime, authority system, database, event store, or UI.

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
| Executor | External effect and execution reporting | Effect attempt, execution status, and execution reporting; subject to external system-of-record semantics for actual effect occurrence | Authority policy, task coordination, domain truth |
| Human | Intervention, review, and optional authority participation through the applicable boundary | Human decision/review within the boundary being exercised | A parallel authority system outside the applicable owner |
| Domain | Domain meaning, truth, acceptance | Domain-specific correctness/acceptance | Infrastructure coordination/authority/execution mechanics |

### Execution vs External System of Record

The Executor owns the attempted operation and execution reporting.

Where an external system of record exists, that system is authoritative for whether the external effect actually occurred.

The Executor MUST NOT override authoritative external-system outcome evidence.

## 4. Core Ownership Invariants

1. Capability ≠ Work ≠ Authority ≠ Execution.
2. Conductor owns coordination; it does not authorize.
3. Solvent owns authority; it does not execute.
4. Executor owns effect execution/reporting; it does not create authority.
5. External system of record owns actual external-effect occurrence where applicable.
6. Agent exercises agency; capability does not imply authority.
7. Human participation does not create a second authority path.
8. Domain meaning remains outside infrastructure authority.
9. Coordination activity is not proof of external execution.
10. Correlation is not authority, truth, or execution proof.
11. Workflow phases are vocabulary, not an additional authoritative state model.
12. Integration is an adapter/boundary, not an independent role.

## 5. External Boundary Matrix

The following are external participant/system boundaries. Internal Agent phases such as Formulate and Interpret are intentionally excluded.

| Boundary | Initiating participant | Receiving participant | Purpose | Required evidence | Human intervention | Semantic / authority owner | Effect / outcome owner | Enforcement / implementation responsibility |
|---|---|---|---|---|---|---|---|---|
| Discover | Agent | Conductor/client interface | Identify available work | Task/work reference | Review/redirect where supported | Conductor for coordination | N/A | Conductor/client implementation |
| Assign / Claim | Agent / Human | Conductor | Establish coordination responsibility | Task + assignment/claim evidence | Reassign/release/reject where supported | Conductor | N/A | Conductor implementation |
| Non-Effect Capability Invocation | Agent | Declared non-effect capability | Perform work with no declared external effect | Artifact/intermediate evidence | Review/redirect/stop where supported | Agent/tool/service as applicable | Tool/service-specific | Tool/service implementation |
| Effect-Capable Operation Invocation | Agent/client | Effect-capable integration / authority boundary | Invoke an externally consequential operation | Proposal + operation identity | Review/revise/reject before effect where supported | Solvent for authority | Pending | Effect-capable integration + Executor |
| Consequential Proposal | Agent/client | Authority boundary | Request externally consequential action | Proposal/version + operation identity | Review/revise/reject | Solvent for authority | Pending | Client/integration |
| Authorization | Solvent | Effect-capable integration / Executor | Decide and expose authority | Authorization evidence + exact operation binding + declaration/identity version | Human may act through applicable Solvent boundary | Solvent | Pending | Integration verifies evidence |
| Execution Eligibility | Effect-capable integration | Executor | Validate required authorization before effect | Authorization evidence + trust basis + operation identity/version | Intervention if boundary supports it | Solvent | Executor | Integration / Executor |
| Execution | Executor | External system/world | Produce external effect | Execution reference + outcome evidence | Cancel/intervene when supported | Solvent for authority fact | Executor for execution reporting; external SOR for actual effect occurrence | Executor/external system |
| Result | Executor / external system | Agent/client / Conductor | Return execution outcome | Execution ID + outcome | Review/challenge/request reconciliation | Prior authority fact remains Solvent-owned | Executor/external SOR | Executor/integration |
| Domain Review | Domain verifier / Human | Domain context | Determine domain acceptance/truth | Domain evidence/review record | Domain-defined | Domain | Domain | Domain verifier/system |
| Next Work | Agent / Human / Conductor | Conductor | Continue, revise, redirect, or conclude work | Task/work reference + review outcome | Inspect / redirect / revise / conclude where supported | Conductor for coordination | N/A | Conductor implementation |

### Boundary Rule

Any operation capable of producing an external effect is consequential regardless of workflow phase or invoking client and MUST use the applicable consequential authorization boundary.

A tool or service used during ordinary work MUST still be covered by a capability declaration.

Non-effect and effect-capable capability declarations are distinct but use the same declaration lifecycle.

## 6. Internal Participant Phases

The following are workflow phases, not external participant boundaries:

- Formulate
- Interpret
- Review preparation
- Agent-side reasoning

They MUST NOT be treated as additional roles, persisted workflow states, or cross-role authority transitions.

## 7. Effect-Capable Operation Invocation → Consequential Proposal

These are sequential aspects of one boundary, not alternative paths.

The required relationship is:

Effect-Capable Operation Invocation
→ Consequential Proposal
→ Authorization
→ Execution Eligibility
→ Execution

An Effect-Capable Operation Invocation MUST enter the applicable Consequential Proposal and Authorization boundaries before producing an external effect.

Conductor does not need to classify the semantic payload.

## 8. Capability Declaration Lifecycle

Any capability used as a declared Loop Engineering boundary MUST have an owned declaration identifying at minimum:

- declaration owner
- declaration version
- effective version/date or equivalent activation reference
- retrievable declaration content
- effect classification
- capability scope
- operation-identity definition where relevant

For effect-capable capabilities, the declaration MUST additionally identify:

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

A change that can alter:

- effect classification
- effect surface
- operation identity
- authorization evidence requirements
- authority-validity model
- replay/idempotency behavior
- execution outcome model
- reconciliation behavior
- material human-intervention capability

MUST change the applicable declaration version and invalidate prior conformance for the affected boundary until re-conformance completes.

A capability changing from non-effect to effect-capable therefore requires:

non-effect declaration
→ declaration version change
→ consequential boundary coverage
→ re-conformance

Undeclared effect-capable behavior discovered during review/testing is a conformance failure.

Where a declared capability has no distinct integration owner, the deployment operator is the default declaration owner and bears responsibility for declaration completeness and lifecycle.

Where a distinct integration owner exists, that owner is responsible for completeness of its declared effect surface.

Deployment/conformance ownership is responsible for ensuring the declaration under test is the one actually governing the tested boundary.

Loop Engineering does not provide a universal oracle capable of discovering hidden undeclared side effects inside arbitrary third-party software.

## 9. Consequential Boundary Rules

### 9.1 Classification

Agent/client classification of consequence is advisory.

The capability declaration defines the effect classification of the invoked operation.

Conductor does not become a domain-specific policy router.

### 9.2 Fail-Closed Enforcement

Every effect-capable integration MUST reject an external effect when required authorization evidence is:

- absent
- invalid
- stale/expired
- operation-mismatched
- unverifiable due to unavailable/failed authority verification

If required authority verification is unavailable, times out, or cannot establish the required trust basis, the evidence MUST be treated as unavailable/invalid and the external effect MUST be rejected.

Authority verification being unavailable means the required verification mechanism cannot establish authorization. Independently verifiable authorization evidence that remains valid under its declared trust/validity model is not treated as unavailable merely because the authority service is unreachable.

### 9.3 Exact Operation Binding

The same declared operation-identity definition MUST be used at proposal, authorization, and execution.

All effect-relevant parameters MUST be included in the identity or explicitly declared immaterial.

Operation identity MUST be defined over semantic effect inputs, not incidental transport metadata.

Any excluded field MUST be explicitly declared immaterial to the external effect.

The declaration MUST define the equality/comparison semantics used to determine whether two operation identities are the same.

The same comparison rule MUST be applied at proposal, authorization, and execution.

The comparison rule MUST be deterministic and MUST operate over the declared semantic effect inputs rather than incidental transport metadata.

The integration owner is responsible for completeness of this definition.

Omitting an effect-relevant parameter is a conformance failure.

Correlation does not establish authorization binding.

### 9.4 Authorization Evidence

Authorization evidence MUST carry or directly and verifiably expose:

- the authorized operation identity/context
- the declaration version/effective reference
- the operation-identity version/effective reference, where separately versioned

The trust basis MUST be independently verifiable against the authority owner and MUST NOT rely solely on an unverified client claim.

The integration declaration MUST identify how evidence is authenticated or validated against the authority owner.

Cryptography is not mandated at the Loop Engineering layer.

## 10. Ownership vs Enforcement

A semantic/authority owner and an enforcement mechanism may differ.

> One semantic owner; possibly another enforcement mechanism.

| Concern | Semantic / authority owner | Enforcement / implementation responsibility |
|---|---|---|
| Authorization | Solvent | Executor-side / integration verification |
| Revocation | Solvent | Applicable integration honors declared revocation model |
| Coordination | Conductor | Conductor implementation |
| Execution reporting | Executor | Executor implementation |
| Actual external-effect occurrence | External system of record, where applicable | Executor reports/reconciles |
| Domain acceptance | Domain | Domain verifier/system |
| Human intervention | Human within applicable boundary | Owning component implements supported controls |

Enforcement responsibility MUST NOT be interpreted as semantic ownership.

## 11. Human Intervention Matrix

| Boundary | Human may inspect | Human may intervene | Limit |
|---|---|---|---|
| Discover | Work/context | Redirect/review | Must use applicable coordination boundary |
| Assign/Claim | Assignment/lifecycle | Reassign/release/reject where supported | Conductor remains coordination owner |
| Non-Effect Work | Work/artifacts | Redirect/stop/review where supported | Capability is component-dependent |
| Effect-Capable Operation Invocation | Proposed operation/context | Review/revise/reject before effect where supported | Material change creates distinct proposal/version |
| Before Authorization | Proposal + evidence | Revise/reject/redirect | Does not create authority |
| Authorization | Authority decision | Exercise authority only through applicable Solvent boundary | Human presence does not bypass Solvent |
| Before Execution | Authorization + operation | Cancel/change progression where supported | Cannot invent authority |
| During Execution | Execution status | Cancel/intervene when Executor supports it | Already-produced effects are not retroactively undone |
| After Execution | Execution evidence | Challenge/request reconciliation/review | Executor/SOR remains outcome authority |
| Domain Review | Domain evidence | Accept/reject/revise | Domain owns truth/acceptance |
| Next Work | Current work + review context | Redirect/revise/conclude where supported | Conductor remains coordination owner |

## 12. Authorization UNKNOWN Temporal Policy

Authorization UNKNOWN/UNAVAILABLE is distinct from DENIED.

An applicable capability/authority declaration MUST define a maximum unresolved-wait/age policy.

The policy MUST define one of:

- explicit operator-controlled indefinite wait, or
- a terminal/remediation action such as retry, escalation, revision, or termination through an existing applicable boundary.

An operator-controlled unresolved hold is a human intervention and MUST produce intervention evidence.

Loop Engineering does not introduce a generic Abandoned, timeout, escalation, or human-escalation state.

## 13. Long-Running Execution and Authority Precedence

Each long-running effect MUST declare:

- authority-validity model
- validity scope
- revalidation model, if any
- mid-flight revocation behavior, if any

Authority-owner constraints are authoritative.

An integration MAY impose a stricter validity constraint but MUST NOT extend, weaken, or override an authority constraint.

Examples:

- Solvent validity 1h + integration validity 24h → effective validity ≤ 1h
- Solvent validity 24h + integration validity 1h → effective validity ≤ 1h

The conformance scenario MUST exercise the declared model.

## 14. Execution Outcomes

Execution outcomes MUST distinguish:

- not attempted
- rejected before effect
- attempted
- succeeded
- failed
- terminated / revoked
- ambiguous / unknown

### Execution Rejected

- Required authorization/evidence validation failed.
- The effect-capable boundary prevented the external effect from being attempted.

It MUST NOT be reinterpreted as execution failure.

### Execution Failed

- The effect was attempted.
- The operation did not successfully complete.

It MUST NOT be reinterpreted as authorization denial.

### Execution Terminated / Revoked

- The effect was attempted.
- The operation was explicitly halted or revoked through an authorized intervention.
- The halt/revocation was known to succeed.

It is distinct from ordinary execution failure.

If termination is requested but the termination outcome cannot be established, the execution outcome is `ambiguous`, not `terminated / revoked`.

`Terminated / revoked` does not imply that zero external effect occurred. Partial external effects may have occurred before or during termination; the external system of record remains authoritative for actual effect occurrence.

### Execution Succeeded

`Succeeded` requires evidence sufficient under the declared execution/outcome model to establish successful effect occurrence.

An asynchronous handoff receipt, queue acceptance, or request acknowledgement alone MUST NOT be reported as `succeeded` unless the declared outcome model explicitly defines that evidence as sufficient to establish the effect.

### Execution Ambiguous

- Effect occurrence cannot be conclusively established.

It MUST remain distinct from success and failure.

Agent interpretation MUST NOT mutate or override the authoritative execution outcome.

A new action may be proposed using the outcome as evidence, but the Agent cannot rewrite the outcome.

## 15. Failure Ownership Matrix

| Failure/outcome | Semantic / authority owner | Enforcement / reporting responsibility | Meaning | Must not be reinterpreted as |
|---|---|---|---|---|
| Coordination failure | Conductor | Conductor implementation | Work coordination could not proceed | Authorization denial |
| Authorization denied | Solvent | Solvent/integration | Authority decision exists and rejects operation | Unknown/unavailable |
| Authorization unknown/unavailable | Solvent | Solvent/integration | No authoritative decision currently established | Denial |
| Authorization invalid/stale | Solvent | Solvent/integration | Required authorization evidence cannot be accepted | Authorized |
| Operation mismatch | Solvent for authority fact | Effect-capable integration / Executor boundary | Requested operation differs from authorized operation | Authorized |
| Execution rejection | Solvent for authorization fact | Effect-capable integration | Required authorization/evidence failed; external effect was not attempted | Execution failure |
| Execution failure | Executor | Executor/external system | Effect was attempted but did not successfully complete | Authorization denial |
| Execution terminated/revoked | Solvent for authority fact; Executor for execution fact | Authorized control boundary / Executor | Effect was attempted and explicitly halted/revoked | Ordinary execution failure |
| Execution ambiguous | Executor | Executor/integration | Effect occurrence cannot be conclusively established | Success or failure |
| Domain rejection | Domain | Domain verifier/Human | Domain acceptance failed | Infrastructure authorization denial |
| Human rejection | Human within applicable boundary | Owning component | Human declined continuation/review/action | Automatic Solvent revocation |

## 16. Evidence / Correlation Matrix

| Fact | Minimum evidence | Semantic / authority owner | Correlation role |
|---|---|---|---|
| Proposal Vn | Proposal/version reference + operation identity | Agent/client | Starting record for action lineage |
| Human intervention | Human action/reference + affected proposal/work reference | Human within applicable boundary | Links intervention to affected work/proposal |
| Authorization | Decision + authorization reference + bound operation | Solvent | Links proposal to authority |
| Declaration | Owner + version + effective reference + retrievable content | Integration owner | Establishes applicable contract version |
| Execution attempt | Execution reference + requested operation + declaration/identity version | Executor | Links authorized operation to effect attempt |
| Execution rejection | Rejection reason/evidence + operation reference | Effect-capable boundary | Demonstrates fail-closed behavior |
| Long-running progress | Start/reference + unrelated Conductor progress evidence | Executor + Conductor | Demonstrates non-blocking coordination |
| Execution outcome | Outcome + execution reference | Executor/external SOR | Links attempt to result |
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

## 17. Persistence / State Ownership

No workflow-owned persistence is introduced.

| Fact | Where it may be recorded | Semantic owner |
|---|---|---|
| Task/lifecycle/assignment/activity | Conductor | Conductor |
| Authority/authorization/revocation | Solvent | Solvent |
| Execution/reporting | Executor | Executor |
| Actual external-effect occurrence | External system of record, where applicable | External system |
| Domain evidence/acceptance | Domain system/record | Domain |
| Agent artifact/reasoning evidence | Agent/task artifact system as applicable | Agent/domain context |
| Conformance observations | Test-only harness/adapter | Conformance test system |

Persistence location does not imply Loop Engineering ownership.

No workflow_state, workflow event store, workflow authority ledger, or equivalent shadow state is permitted without a Growth Gate-approved architectural change.

## 18. Integration / Capability Declaration

Every declared capability boundary MUST satisfy the declaration requirements defined in §8.

This includes:

- declaration ownership
- declaration version/effective reference
- retrievable declaration content
- effect classification
- capability scope
- operation-identity definition and comparison rule where relevant

Effect-capable capabilities additionally require the authority, trust, validity, replay, outcome, reconciliation, recovery, and human-intervention declarations specified in §8.

Declaration changes that alter effect classification, effect surface, operation identity, or any other listed semantic property invalidate prior conformance for the affected boundary until re-conformance completes.

## 19. Production Conformance Boundary

Conformance has two layers:

1. Protocol Fixture Conformance
2. Integration Conformance

Integration conformance MUST exercise the real authorization/enforcement path.

A safe target, sandbox, dry-run, or reversible environment may neutralize the final external effect, but it MUST NOT bypass:

- authorization verification
- exact operation binding
- fail-closed enforcement

Passing a TestExecutor/TestSolvent scenario does not prove the production integration behaves correctly.

## 20. Client / Runtime Substitution Boundary

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

## 21. Agent Skill Boundary

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

## 22. Conformance Harness Boundary

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

## 23. Boundary Anti-Patterns

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
- an effect-capable integration using a declaration that is not the currently effective declared version

## 24. Role Substitution Tests

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

## 25. Matrix-Level Acceptance Criteria

The Role / Boundary Matrix is acceptable when:

1. Every authoritative fact has one semantic/authority owner.
2. No role silently owns another role's semantics.
3. Every consequential boundary identifies the authority owner and effect/outcome owner.
4. Human intervention is defined without creating a second authority path.
5. Domain remains outside infrastructure authority.
6. Exact operation binding is distinguishable from correlation.
7. Failure/outcome ownership is explicit, including rejection, failure, termination/revocation, and ambiguity.
8. Persistence locations are distinguished from semantic ownership.
9. Production conformance is distinguishable from fixture conformance.
10. Client and orchestration substitution preserve the contract.
11. No matrix row requires a new runtime, database, event store, UI, or authority engine.
12. Effect-capable integration completeness is treated as an explicit integration/deployment responsibility.
13. Non-effect and effect-capable declarations share a clear lifecycle.
14. An unavailable authority verification path fails closed.
15. Applicable declaration/operation-identity version is observable and respected at execution.
16. Any proposed exception can be evaluated through the Growth Gate.

## 26. Adversarial Questions for the Conformance Matrix

Before accepting the Conformance Test Matrix, verify:

1. Can a non-effect capability silently become effect-capable without declaration/re-conformance?
2. Can any tool invoked during ordinary work produce an external effect without entering the consequential boundary?
3. Can a proposal be authorized for X but executed as Y?
4. Is the same operation-identity definition used at proposal, authorization, and execution?
5. Can an effect-capable integration bypass Solvent?
6. Can authority verification failure result in execution?
7. Can UNKNOWN authorization remain unresolved without following the declared temporal policy and evidence requirement?
8. Can an ambiguous execution outcome be silently treated as success?
9. Can a fail-closed rejection be confused with execution failure?
10. Can a terminated/revoked execution be mislabeled ordinary failure?
11. Can a human intervene without creating shadow authority?
12. Can Conductor cancellation mutate Solvent authority?
13. Can integration declaration drift make an old conformance result appear valid?
14. Can the conformance harness pass while bypassing the real enforcement path?
15. Can an integration hide an undeclared effect-capable operation and still claim conformance?
16. Can a new component absorb workflow semantics merely by renaming them?
17. Can Agent interpretation override an authoritative Executor/external-SOR outcome?
18. Can two conforming implementations apply different operation-identity equality rules to the same declared operation?
19. Can a declared generic capability exist without a declaration owner?
20. Can valid independently verifiable authorization evidence be rejected solely because the authority service is temporarily unreachable?

If any answer is unclear, the Conformance Test Matrix MUST expose that ambiguity as a test.

## 27. Review Disposition

### Recurring Finding Status

| Finding | Status | Rationale |
|---|---|---|
| Identity comparison semantics | Adopted | Security-critical conformance property |
| Generic capability owner | Adopted | Prevents ownerless declarations |
| Concurrency ownership | Deferred | Primarily a Conformance Test Matrix / implementation concern unless POC proves otherwise |
| Client-internal state | Rejected as Loop concern | Client responsibility |
| Production observation persistence | Deferred to conformance requirements | Must remain test/integration concern, not Loop-owned state |
| Cached-auth evidence marker/version | Adopted in boundary requirements | Needed for deterministic evidence and queued/long-running execution |


| Finding / proposal | Disposition | Reason |
|---|---|---|
| Mandatory workflow runtime | Rejected | Violates minimal/no-runtime architecture |
| Central Solvent effect registry | Rejected | Expands Solvent from authority kernel toward global effect registry |
| Generic operations gateway as mandatory architecture | Rejected | Adds a de facto fifth infrastructure role |
| Generic Conductor escalation state | Rejected | Not required; existing boundaries can express policy/remediation |
| Conductor → Solvent automatic revocation | Rejected | Cross-owner authority mutation |
| New human-authoring authority path | Rejected | Human revisions are new proposals through the same authority boundary |
| Mandatory verification SDK | Rejected as specification requirement | Implementation option; conformance tests verify the property |
| Mandatory cryptography / RFC 8785 | Rejected as specification requirement | Tool/runtime-mechanism neutrality |
| Universal effect-detection oracle | Accepted limitation | Impossible to guarantee generically for arbitrary third-party software |
| Shared verification library | Deferred / optional | May be justified by implementation evidence after POC |
| Egress gateway / proxy | Allowed implementation option | Not a Loop Engineering role |
| Test-only deterministic fixtures | Accepted | Required for protocol conformance |
| Production-path integration conformance | Accepted | Prevents certification theater |
| Declaration lifecycle for all capability boundaries | Accepted | Prevents non-effect → effect drift |

## 28. Precedence and Deliverable Boundary

The Workflow Specification v0.3 is the normative parent artifact.

This matrix may refine and operationalize ownership/boundary distinctions but MUST NOT contradict the parent specification.

This matrix is the basis for:

- Conformance Test Matrix v0.1
- BM-IST Validation Scenarios v0.1
- Agent Skill
- Conformance Harness
- POC Runbook

It does not authorize implementation of a new workflow runtime or infrastructure component.
