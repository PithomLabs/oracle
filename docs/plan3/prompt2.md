We have completed Phase 1 of the Reference Loop.

Phase 1 result:

    PASS WITH LIMITATIONS

The happy path is reproducible and the core boundaries are holding.

Do NOT expand the architecture yet.

The next objective is a very small boundary-validation experiment:

1. mismatched actor_id authorization denial
2. execution attempt without valid authorization denial

Keep Solvent completely frozen.

======================================================================
1. SOLVENT FREEZE — ABSOLUTE
======================================================================

Before doing anything:

    cd /home/chaschel/Documents/go/solvent-main

Verify:

    git rev-parse HEAD
    git status --short
    git diff --name-only

Expected:

    HEAD = 7602699
    working tree = clean

If this fails:

    STOP
    report BLOCKED

Do NOT modify Solvent.

No source edits.
No schema edits.
No migrations.
No API changes.
No auth changes.
No compatibility patches.

======================================================================
2. TEST ONLY TWO NEGATIVE CASES
======================================================================

Add exactly these focused tests.

----------------------------------------------------------------------
TEST A — MISMATCHED ACTOR AUTHORIZATION DENIAL
----------------------------------------------------------------------

Establish a valid consequential operation.

Authorize it under actor A.

Then attempt to use that authorization as actor B.

Expected:

    authorization is denied

The test must demonstrate that authorization is bound to the correct actor
identity.

Verify that the denial happens through the real authorization/enforcement
path rather than through a harness-side precondition.

Document:

    authorized actor
    attempted actor
    operation identity
    authorization evidence
    observed denial

----------------------------------------------------------------------
TEST B — EXECUTION WITHOUT AUTHORIZATION DENIAL
----------------------------------------------------------------------

Attempt execution of a consequential operation without a valid authorization
intent.

Possible forms include:

    - missing intent_id
    - nonexistent intent_id
    - expired/invalid authorization, if already supported by the frozen API

Use whichever case is already naturally supported by the existing interface.

Expected:

    execution is denied

The executor must NOT execute the operation.

The important property is:

    AUTHORIZE ≠ EXECUTE

and:

    execution requires valid authorization evidence

Do not create a new authorization mechanism merely to implement this test.

======================================================================
3. NO NEW ARCHITECTURE
======================================================================

Do NOT implement:

- full adversarial suite
- cancellation tests
- concurrency
- replay tests
- revocation tests
- long-running execution
- BM-IST
- Agent Skill
- Conductor Skill
- full protocol specification
- Conformance Matrix
- new runtime
- new event bus
- new workflow abstraction
- Intent/Work/Authorization/Effect/Outcome primitive

The purpose is to test whether the existing architecture fails safely at two
critical boundaries.

======================================================================
4. USE REAL PARTICIPANT BOUNDARIES
======================================================================

Do not mock away the properties being tested.

Use:

    Agent/reference harness
        ↓
    Conductor where applicable
        ↓
    frozen Solvent REST/auth path
        ↓
    existing Executor path

Do not bypass Solvent authorization.

Do not implement authorization logic in the test harness.

RecordingFunc may remain the execution implementation, because this phase is
testing authorization denial, not proving external effect.

======================================================================
5. VERIFY NO SIDE EFFECT ON DENIAL
======================================================================

For each negative case verify:

    authorization denied
    execution not performed
    no false EFFECT_CONFIRMED
    no fabricated result
    no misleading Conductor state

Where appropriate, confirm that the executor's invocation count remains zero.

======================================================================
6. DOCUMENT THE ACTUAL OBSERVATION
======================================================================

Create/update:

    PHASE1_NEGATIVE_PATH.md

For each test record:

    Scenario:
    Setup:
    Operation identity:
    Actor:
    Authorization state:
    Execution request:
    Expected result:
    Actual result:
    Executor invoked:
    Effect observed:
    Evidence source:
    Verdict:

Also classify any failure as exactly one of:

    IMPLEMENTATION_DEFECT
    INTEGRATION_DEFECT
    EXECUTOR_DEFECT
    DEPLOYMENT_ENVIRONMENT_DEFECT
    SPECIFICATION_DEFECT
    NEW_SECURITY_PROPERTY

Do not classify a failure as a new architectural requirement without evidence.

======================================================================
7. SPECIAL ATTENTION: FROZEN SOLVENT REST API
======================================================================

The Phase 1 review found that the frozen Solvent REST API does not return
intent_id from AuthorizeAction and the Reference Loop currently uses a DB
fallback.

Do NOT modify Solvent to fix this.

Keep the existing workaround if it is required for these tests.

However, document whether that fallback affects either negative-path test.

If the fallback becomes unsafe or prevents the test from faithfully exercising
the real authorization boundary:

    STOP
    report BLOCKED

Do not patch Solvent.

======================================================================
8. RUN VALIDATION
======================================================================

Run the complete Reference Loop suite:

    cd /home/chaschel/Documents/go/oracle/reference-loop
    go test ./...
    go vet ./...

Then explicitly run the new negative tests.

Afterward, re-check Solvent:

    cd /home/chaschel/Documents/go/solvent-main
    git rev-parse HEAD
    git status --short
    git diff --name-only

Expected:

    HEAD = 7602699
    working tree = clean
    no Solvent modifications

======================================================================
9. FINAL ASSESSMENT
======================================================================

Answer these questions:

    1. Does mismatched actor authorization fail closed?
    2. Does execution without authorization fail closed?
    3. Is the executor prevented from executing in both cases?
    4. Is any false effect confirmation produced?
    5. Does Conductor remain a coordination component rather than an authority?
    6. Does Solvent remain the sole authority?
    7. Did these tests reveal a specification defect?
    8. Did these tests reveal a genuinely new security property?
    9. Is a new cross-role protocol primitive now justified?

For question 9, default to:

    NO

unless the actual experiment demonstrates that the existing contracts cannot
represent the required behavior.

======================================================================
10. STOP AFTER THIS EXPERIMENT
======================================================================

Do not automatically proceed to the next phase.

Return one of:

    PASS
    PASS WITH LIMITATIONS
    BLOCKED

If PASS or PASS WITH LIMITATIONS, provide one recommended next experiment,
but do not implement it.

The goal is:

    happy path
        ↓
    boundary denial tests
        ↓
    evaluate what the evidence actually requires

Do not design the next architecture before the evidence requires it.