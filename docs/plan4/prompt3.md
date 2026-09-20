Phase 2 is NOT yet passed.

The latest run exercised the in-process Solvent server with a simulated
executor:

    output="simulated_success"

Therefore:

    PASS-7 = NOT PROVEN
    PASS-8 = NOT PROVEN

Do NOT mark the Phase 2 experiment PASS.

This is an integration/test-wiring issue, not an architecture issue.

======================================================================
OBJECTIVE
======================================================================

Make TestPhase2FirstRealEffect exercise the ACTUAL frozen Solvent execution
path with the real GitHub executor:

    github.NewHTTPProvider
    executor = github_trigger_workflow
    GITHUB_TOKEN
    standalone frozen Solvent server

The test must NOT use the in-process simulated executor for the Phase 2
real-effect experiment.

======================================================================
CONSTRAINTS
======================================================================

DO NOT:

- modify frozen Solvent
- change Solvent source
- add a new Executor
- modify Conductor architecture
- add workflow infrastructure
- weaken authorization
- bypass Solvent
- fabricate EFFECT_CONFIRMED
- treat executor invocation as external effect

Solvent must remain:

    HEAD = 7602699
    working tree clean

======================================================================
REQUIRED EXECUTION PATH
======================================================================

The actual path must be:

    Reference Loop
        ↓
    frozen Solvent REST server
        ↓
    frozen Solvent authorization
        ↓
    frozen Solvent execution boundary
        ↓
    github_trigger_workflow
        ↓
    GitHub Actions API
        ↓
    GitHub workflow run
        ↓
    GitHub SOR

The frozen Solvent source must remain unchanged.

======================================================================
1. START REAL SOLVENT SERVER
======================================================================

Use the frozen Solvent server in its existing standalone configuration
with the real GitHub executor.

Verify the registered executor is:

    github_trigger_workflow

and NOT:

    simulated_success
    RecordingFunc
    fake executor

Do not modify Solvent to accomplish this.

======================================================================
2. GITHUB CREDENTIAL
======================================================================

Require:

    GITHUB_TOKEN

from the environment.

Never print or persist the token.

If absent:

    BLOCKED

Do not substitute a fake executor.

======================================================================
3. RUN THE EXISTING PHASE 2 SCENARIO
======================================================================

Keep:

    plan_id
    plan_version
    capability_ref
    declaration version/hash
    execution_run_id
    operation_id

unchanged.

The exact operation remains:

    deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<execution_run_id>

After dispatch, capture:

    github_run_id

as correlation evidence only.

Do NOT put github_run_id into operation identity.

======================================================================
4. POLL THE REAL GITHUB SOR
======================================================================

After workflow dispatch:

    locate the created workflow run

Then poll the GitHub Actions API until:

    status == completed

or the declared timeout expires.

Classify:

    completed + success
        → EFFECT_CONFIRMED

    completed + failure/cancelled
        → external execution failure

    still running at timeout
        → AMBIGUOUS

    API failure with effect state unresolved
        → AMBIGUOUS

Do not infer effect from HTTP 204 dispatch acceptance.

======================================================================
5. UPDATE PASS CRITERIA

The final report must use:

    PASS-7:
        real GitHub executor produced the external operation attempt

    PASS-8:
        GitHub Actions SOR confirmed the resulting external state

Do not report PASS-7 merely because ExecuteAction returned success.

======================================================================
6. EVIDENCE

Capture:

    execution_run_id
    operation_id
    intent_id
    github_run_id
    Solvent authorization evidence
    executor invocation
    GitHub dispatch response
    GitHub workflow-run evidence
    final status
    final conclusion

The distinction must remain:

    executor invocation
        !=
    external effect
        !=
    result observation

======================================================================
7. FINAL VERDICT

Possible final verdicts:

    PASS
    PASS WITH LIMITATIONS
    BLOCKED

PASS requires:

    PASS-1 through PASS-11

and specifically:

    real GitHub execution
    +
    authoritative GitHub SOR confirmation

If the real GitHub executor cannot be reached without modifying frozen
Solvent:

    BLOCKED

Do not modify Solvent.

======================================================================
8. STOP

After the real-effect experiment:

    update PHASE2_RESULTS.md
    create/update PHASE2_EVIDENCE.md
    verify Solvent remains frozen
    STOP

Do not proceed to Phase 3 automatically.