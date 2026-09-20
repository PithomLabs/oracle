ARGUS POC — PHASE 0 IMPLEMENTATION

You are now the implementation agent for the frozen ARGUS POC.

The authoritative implementation baseline is:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

Status:

    FROZEN — IMPLEMENTATION BASELINE

DO NOT reinterpret the architecture.
DO NOT redesign the system.
DO NOT add features.
DO NOT modify frozen components unless the execution plan explicitly identifies
the required Growth Gate change.

Your first task is ONLY PHASE 0 — REPOSITORY RECONNAISSANCE.

==================================================
PHASE 0 OBJECTIVE
==================================================

Establish and verify the exact starting state of:

    Solvent
    Conductor
    Oracle / ARGUS
    reference-loop

before any implementation begins.

Phase 0 is READ-ONLY.

DO NOT create production code.
DO NOT create new packages.
DO NOT create migrations.
DO NOT modify schemas.
DO NOT modify Conductor.
DO NOT modify reference-loop.
DO NOT modify Solvent.
DO NOT modify the frozen execution plan.

==================================================
AUTHORITATIVE PLAN
==================================================

Read:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

Treat it as the binding implementation contract.

Do not substitute assumptions from previous versions of the plan.

==================================================
REPOSITORIES TO INSPECT
==================================================

Inspect:

    /home/chaschel/Documents/go/solvent-main
    /home/chaschel/Documents/go/conductor
    /home/chaschel/Documents/go/oracle

Also inspect:

    /home/chaschel/Documents/go/oracle/reference-loop

Reference-loop is a separate Go module and MUST remain untouched.

==================================================
REQUIRED VERIFICATION
==================================================

Run exactly the relevant baseline checks from the frozen plan.

SOLVENT:

    cd /home/chaschel/Documents/go/solvent-main

    git rev-parse HEAD
    git status --short
    go test ./...

Verify:

    db/010_debt_opaque.sql exists
    debt default is ARRAY[]::TEXT[]
    Plan 11.1 debt-opaque implementation is present
    FullDebt is absent from kernel/
    working tree is clean

Also inspect the current Solvent API/schema relevant to the POC:

    belief
    evidence
    belief_edge
    action_intent
    audit_activity
    refusal behavior / Verdict
    authority target lifecycle

IMPORTANT:

The frozen plan includes one explicit Growth Gate exception:

    POST /v1/beliefs/{parent_id}/edges

This endpoint is NOT expected to exist yet.

Record whether it currently exists.
Do NOT implement it in Phase 0.

CONDUCTOR:

    cd /home/chaschel/Documents/go/conductor

    git rev-parse HEAD
    git status --short
    go test ./...

Verify:

    Conductor is domain-agnostic
    existing Solvent integration remains intact
    existing UI is still the current operational UI
    no ARGUS modifications have already been introduced

ORACLE:

    cd /home/chaschel/Documents/go/oracle

    git rev-parse HEAD
    git status --short

Inspect:

    current module layout
    existing oracle/go.mod state
    existing reference-loop/go.mod
    existing plans
    existing domain/research artifacts
    existing code that may overlap with the frozen ARGUS plan

Verify whether any Phase 2+ implementation already exists.

REFERENCE-LOOP:

Inspect:

    oracle/reference-loop/

Verify:

    separate go.mod
    current build/test state if practical
    no changes are made

==================================================
STATE COMPARISON
==================================================

Compare actual repository state against the frozen Phase 0 assumptions.

Create a reconciliation table:

    AREA
    EXPECTED BY PLAN
    ACTUAL STATE
    MATCH / MISMATCH
    EVIDENCE
    ACTION

At minimum include:

    Solvent freeze state
    Solvent migration state
    Solvent API state
    belief_edge API availability
    Solvent Verdict behavior
    Solvent audit/refusal behavior
    Conductor state
    Oracle module state
    reference-loop state
    existing ARGUS-related code
    working-tree cleanliness

==================================================
CRITICAL RULES
==================================================

1. Phase 0 is observational only.

2. If something differs from the plan:
       STOP and report it.
   Do not silently repair it.

3. Do not modify frozen systems merely to make Phase 0 pass.

4. Do not create the Growth Gate edge endpoint yet.

5. Do not create PackRegistry yet.

6. Do not create ArtifactRegistry yet.

7. Do not create Coordinator yet.

8. Do not create Trust UI yet.

9. Do not modify Conductor.

10. Do not modify reference-loop.

11. Do not commit anything unless the repository already has an established
    Phase 0 reporting convention that requires a documentation-only commit.
    Prefer no commit.

==================================================
PHASE 0 DELIVERABLE
==================================================

Produce:

    oracle/plan/PHASE0_RECONNAISSANCE.md

This is a read-only reconnaissance report.

It must contain:

    1. Execution timestamp
    2. Repository paths inspected
    3. Git HEAD for each repository
    4. Working-tree status
    5. Test results
    6. Solvent state verification
    7. Conductor state verification
    8. Oracle state verification
    9. reference-loop state verification
    10. API/schema findings
    11. Expected-vs-actual reconciliation table
    12. Any mismatches
    13. Explicit recommendation:
            PROCEED TO PHASE 1
        or
            STOP — RECONCILIATION REQUIRED

Do not claim anything passed unless you actually verified it.

==================================================
PHASE 0 ACCEPTANCE CRITERIA
==================================================

PASS only when:

    - Solvent tests pass
    - Conductor tests pass
    - required working trees are clean
    - Plan 11.1 state is actually present
    - FullDebt is absent from Solvent kernel
    - reference-loop remains untouched
    - current belief_edge API availability is explicitly documented
    - current Verdict/refusal behavior is explicitly documented
    - all known deviations from the frozen plan are recorded
    - PHASE0_RECONNAISSANCE.md is complete

If any critical prerequisite is missing:

    DO NOT PROCEED TO PHASE 1.

==================================================
FINAL RESPONSE FORMAT
==================================================

After completing Phase 0, report only:

    PHASE 0 RESULT: PASS / STOP

    Repositories inspected:
    ...

    Tests:
    ...

    Critical findings:
    ...

    Mismatches:
    ...

    Deliverable:
    oracle/plan/PHASE0_RECONNAISSANCE.md

    Next phase:
    PHASE 1 — FREEZE RECONCILIATION

Do not begin Phase 1 in this run.