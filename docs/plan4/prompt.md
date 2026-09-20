Use this exact prompt. It treats the design as frozen and moves directly to implementation of Phase 1.5. No further architecture review.

```text
You are now implementing Phase 1.5 of the Reference Loop.

IMPORTANT:
The Workflow Design v1.0 is frozen for this phase.

DO NOT redesign it.
DO NOT reopen architectural questions.
DO NOT add new workflow machinery.
DO NOT modify the frozen Solvent codebase.

Your job is to implement and execute the Phase 1.5 validation battery against the frozen design.

======================================================================
1. FROZEN BASELINE
======================================================================

Treat these as immutable for this phase:

Workflow Design:
    Loop Engineering Workflow Design v1.0

Solvent:
    /home/chaschel/Documents/go/solvent-main
    frozen commit = 7602699

Before starting:

    cd /home/chaschel/Documents/go/solvent-main
    git rev-parse HEAD
    git status --short
    git diff --name-only

Expected:

    HEAD = 7602699
    working tree = clean
    no source changes

If not:

    STOP
    report BLOCKED
    do not modify Solvent

======================================================================
2. PHASE 1.5 OBJECTIVE
======================================================================

Prove the minimum safety/binding properties required before Phase 2
(real external-effect validation).

Run exactly these five scenarios:

    A. operation mismatch
    B. duplicate delivery
    C. stale/invalid authorization
    D. UNKNOWN / unavailable authorization
    E. stale claim × expired intent

Also verify:

    F. declaration resolution
    G. operation-identity equality
    H. frozen-Solvent integrity

Do not add unrelated scenarios.

======================================================================
3. TEST A — OPERATION MISMATCH
======================================================================

Prove exact operation binding.

Scenario:

    authorize operation X
    attempt execution of operation Y

where:

    X != Y

The difference MUST be effect-relevant.

Expected:

    execution denied
    executor not invoked
    no external effect
    no false success/effect confirmation

Verify the mismatch is rejected by the actual authorization/enforcement path,
not by harness-side validation.

Record:

    authorized operation identity
    requested operation identity
    authorization evidence
    rejection
    executor invocation count

This is the critical Phase 1.5 binding test.

======================================================================
4. TEST B — DUPLICATE DELIVERY
======================================================================

Use an already-authorized exact operation.

Deliver the same execution request twice.

Use the declared integration behavior supported by the current executor/test
fixture.

Determine whether the implementation provides:

    deduplication
    idempotency
    single-use authorization
    explicit duplicate rejection
    or another declared behavior

Do NOT invent a new exactly-once mechanism.

The test must verify that repeated delivery does not create an unintended
duplicate effect when the declared integration policy forbids it.

If the current controlled executor cannot exercise the relevant declared
behavior, report:

    NOT YET PROVEN

Do not build new infrastructure.

======================================================================
5. TEST C — STALE / INVALID AUTHORIZATION
======================================================================

Exercise the smallest existing stale/invalid authorization case supported by
the current Solvent/integration path.

Examples:

    expired evidence
    stale intent
    invalid authorization reference

Expected:

    execution rejected
    executor not invoked
    no external effect

Do not modify Solvent to manufacture the condition.

Use an existing test fixture, controlled authority, or available mechanism.

======================================================================
6. TEST D — UNKNOWN / UNAVAILABLE
======================================================================

Create a scenario where no authoritative authorization decision can currently
be established.

Expected semantics:

    UNKNOWN / UNAVAILABLE
        ≠
    DENIED

and:

    UNKNOWN / UNAVAILABLE
        ≠
    AUTHORIZED

Unless an independently valid existing authorization already covers the exact
operation, the effect-capable path must not execute.

Verify that the Agent/client may:

    retry
    wait
    escalate
    revise

according to the existing integration behavior, but do not create a retry
engine.

Record the declared unresolved-wait/age behavior where available.

======================================================================
7. TEST E — STALE CLAIM × EXPIRED INTENT
======================================================================

This is a semantic interaction test.

Scenario:

    Agent claims work
        ↓
    valid authorization intent exists
        ↓
    intent becomes expired/invalid
        ↓
    Agent submits / attempts execution after expiry

Determine exactly which boundary rejects the operation and what state is
returned.

The goal is to characterize:

    coordination claim state
    versus
    Solvent authorization validity

Do NOT implement claim TTL/heartbeat machinery in this phase.

This test only establishes the semantics.

Expected outcome must be documented, not guessed.

======================================================================
8. TEST F — DECLARATION RESOLUTION
======================================================================

Verify that a task's:

    capability_ref

resolves to:

    declaration owner
    declaration version
    effective reference
    content hash
    retrievable declaration content

Verify that the integrity pin is actually checked.

A task referencing an unresolved/invalid declaration must not become actionable.

Do not build a new capability registry.

Use the existing declaration mechanism or the smallest test fixture necessary.

======================================================================
9. TEST G — OPERATION IDENTITY EQUALITY
======================================================================

Verify that the same operation-identity definition is used at:

    proposal
    authorization
    execution

Test at least:

    exact equality
    one effect-relevant parameter changed
    declared-immaterial parameter changed

Expected:

    effect-relevant change → different operation
    declared-immaterial change → same operation

The comparison must be deterministic.

Document the canonicalization/equality rule actually implemented.

Do not invent cryptography unless the existing declaration requires it.

======================================================================
10. TEST H — SOLVENT FREEZE INTEGRITY
======================================================================

At the beginning and end of Phase 1.5 verify:

    cd /home/chaschel/Documents/go/solvent-main

    git rev-parse HEAD
    git status --short
    git diff --name-only
    git diff --stat

Expected:

    HEAD = 7602699
    working tree clean
    no source modifications

If any Solvent modification occurs:

    STOP
    report BLOCKED

Do not automatically revert it.

======================================================================
11. TEST ENVIRONMENT

Use:

    Conductor
        existing implementation

    Solvent
        frozen implementation at 7602699

    Executor
        existing deterministic/fake executor where appropriate

    External effect
        MUST NOT be real in Phase 1.5

Phase 1.5 is a safety/binding gate.

Real external-effect validation belongs to Phase 2.

======================================================================
12. EVIDENCE

For every scenario record:

    scenario_id
    run_id
    task_id
    plan_id / plan_version where applicable
    operation_id
    actor
    capability_ref
    declaration version/hash
    setup
    request
    expected result
    actual result
    authorization state
    executor invocation
    effect classification
    authoritative evidence source
    verdict

Do not manufacture authoritative facts.

The harness aggregates evidence but does not become an authority ledger.

======================================================================
13. PASS / FAIL CRITERIA

A scenario passes only when the expected boundary behavior is observed
through the real relevant participant boundary.

PASS examples:

    X authorized, Y executed
        → denied by enforcement boundary

    expired authorization
        → denied

    UNKNOWN authorization
        → not treated as AUTHORIZED

    unresolved declaration
        → not actionable

    effect-relevant identity change
        → different operation

A scenario is:

    NOT YET PROVEN

when the current controlled implementation cannot establish the required
behavior without adding new architecture.

A scenario is:

    BLOCKED

when the test cannot be run because the environment or frozen participant
cannot support the existing contract unchanged.

Do not convert NOT YET PROVEN into PASS.

======================================================================
14. FAILURE CLASSIFICATION

For any unexpected result classify after inspecting evidence:

    IMPLEMENTATION_DEFECT
    INTEGRATION_DEFECT
    EXECUTOR_DEFECT
    DEPLOYMENT_ENVIRONMENT_DEFECT
    SPECIFICATION_DEFECT
    NEW_SECURITY_PROPERTY

Do not pre-classify.

Do not modify architecture merely to make the test pass.

======================================================================
15. DELIVERABLES

Create:

    PHASE1_5_PLAN.md
    PHASE1_5_RESULTS.md
    PHASE1_5_DISPOSITION.md

PHASE1_5_PLAN.md:
    - scenarios
    - setup
    - pass/fail criteria
    - required evidence
    - explicit out-of-scope items

PHASE1_5_RESULTS.md:
    - one section per scenario
    - exact observed evidence
    - verdict
    - Solvent integrity check

PHASE1_5_DISPOSITION.md:
    - finding
    - classification
    - disposition
    - evidence
    - whether architecture change is justified

Do not create a new framework.

======================================================================
16. OUT OF SCOPE

Do NOT implement:

    claim TTL
    heartbeat
    scheduler
    proof-token subsystem
    plan-scope enforcement engine
    Conductor policy engine
    cryptographic Conductor↔Solvent ledger
    workflow templates
    molecules
    DAG runtime
    knowledge graph
    new event bus
    Agent planning engine
    new Solvent capabilities
    real external effect
    BM-IST integration

If one of these appears necessary, STOP and report a Growth Gate candidate.

======================================================================
17. FINAL GATE

At the end, produce this summary:

    Phase 1.5 Verdict:
        PASS
        PASS WITH LIMITATIONS
        BLOCKED

    Operation mismatch:
        PASS / NOT YET PROVEN / BLOCKED

    Duplicate delivery:
        PASS / NOT YET PROVEN / BLOCKED

    Stale/invalid authorization:
        PASS / NOT YET PROVEN / BLOCKED

    UNKNOWN/unavailable:
        PASS / NOT YET PROVEN / BLOCKED

    Stale claim × expired intent:
        PASS / NOT YET PROVEN / BLOCKED

    Declaration resolution:
        PASS / NOT YET PROVEN / BLOCKED

    Operation identity equality:
        PASS / NOT YET PROVEN / BLOCKED

    Solvent frozen:
        YES / NO

    New architecture required:
        YES / NO

    Growth Gate candidate:
        YES / NO

======================================================================
18. STOP AFTER PHASE 1.5

Do not proceed automatically to Phase 2.

Do not redesign the workflow.

The output of Phase 1.5 is evidence.

The next human decision will be whether the evidence is sufficient to authorize
the Phase 2 real external-effect experiment.
```

This is now an **execution plan, not another design exercise**. The Phase 1.5 work should deliberately expose where the frozen design holds, where behavior remains unproven, and where reality actually earns a Growth Gate.
