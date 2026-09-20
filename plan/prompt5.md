ARGUS POC — PHASE 1 IMPLEMENTATION
FREEZE RECONCILIATION

You are now implementing ONLY Phase 1 of the frozen ARGUS POC.

Authoritative plan:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

Phase 0 has passed.

Phase 0 report:

    oracle/plan/PHASE0_RECONNAISSANCE.md

Do NOT begin Phase 2.
Do NOT implement Domain Pack, Packet, Verifier, Coordinator, Trust UI,
or the Solvent edge endpoint in this phase.

==================================================
PHASE 1 OBJECTIVE
==================================================

Establish the immutable baseline for ARGUS implementation and create the
strategy-level freeze reconciliation record.

This phase is primarily documentation and verification.

The purpose is to make explicit:

    what is frozen
    what was superseded
    what was amended
    what remains deferred
    what explicit Growth Gate exceptions exist

==================================================
CRITICAL RULE
==================================================

The repository state discovered in Phase 0 is authoritative.

Do NOT invent a new freeze hash.

Do NOT assume a commit hash from the plan is current.

Determine the current Solvent HEAD directly.

==================================================
STEP 1 — VERIFY SOLVENT FREEZE
==================================================

Run:

    cd /home/chaschel/Documents/go/solvent-main

    git rev-parse HEAD
    git status --short
    go test ./...
    grep -r "FullDebt" kernel/

Record:

    NEW_FREEZE_HASH = actual current HEAD

Working tree must contain no modified tracked files.

Untracked files discovered in Phase 0 must not be deleted merely to obtain
a clean baseline.

The Phase 0 environment limitation remains valid:

    Solvent integration tests may fail if CockroachDB is unavailable.

Do NOT claim all Solvent tests pass when the database is unavailable.

Record:

    total passed
    total failed
    exact environmental failure cause

If failures occur for any reason other than the previously documented
CockroachDB environment limitation, STOP.

==================================================
STEP 2 — VERIFY PLAN 11.1 STATE
==================================================

Verify:

    db/010_debt_opaque.sql exists

    debt default is:

        ARRAY[]::TEXT[]

    FullDebt is absent from kernel/

    wizardDebt exists in:

        internal/belief/debt.go

    EnterBelief accepts caller-supplied debt

    existing debt vocabulary remains outside the generic Solvent kernel

Do not modify any of these files.

==================================================
STEP 3 — VERIFY FROZEN SCHEMA
==================================================

Inspect and record the current relevant frozen schema:

    belief
    belief_edge
    evidence
    action_intent
    audit_activity
    refusal_log

Record the important architectural facts:

    belief.debt is TEXT[] and opaque to Solvent

    belief_edge exists as frozen schema

    belief_edge currently has no API endpoint

    action_intent remains downstream of promoted authority state

    audit_activity is canonical activity history

    refusal behavior uses existing Solvent mechanisms

Do NOT implement the Growth Gate edge endpoint yet.

==================================================
STEP 4 — CREATE FREEZE RECONCILIATION
==================================================

Create:

    oracle/plan/freeze_reconciliation.md

This document is the strategy-level freeze record.

Include the following sections exactly or with equivalent structure:

    1. Baseline
    2. Prior Freeze
    3. Current Freeze
    4. Superseded Designs
    5. Amended Designs
    6. Deferred Capabilities
    7. Growth Gate Exceptions
    8. Repository Boundaries
    9. Verification Results
    10. Phase 1 Acceptance

==================================================
BASELINE
==================================================

Record:

    ARGUS POC Execution Plan v1.1
    Status: FROZEN — IMPLEMENTATION BASELINE

Reference:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

Record the Phase 0 report:

    oracle/plan/PHASE0_RECONNAISSANCE.md

==================================================
PRIOR FREEZE
==================================================

Record:

    Prior freeze hash:
        7602699

State explicitly that it was superseded by the Plan 11.1 debt-opaque design.

Do NOT imply the prior commit is the current baseline.

==================================================
CURRENT FREEZE
==================================================

Record:

    NEW_FREEZE_HASH = actual output of git rev-parse HEAD

Also record:

    repository path
    verification timestamp
    working-tree state

==================================================
SUPERSEDED DESIGNS
==================================================

Explicitly document at minimum:

    old FullDebt kernel constant
    old generic Solvent coupling to wizard deployment-review vocabulary

Do not remove historical records.

==================================================
AMENDED DESIGNS
==================================================

Document at minimum:

    debt is opaque TEXT[]
    Solvent default debt is empty
    callers supply domain-specific debt
    wizardDebt owns wizard-specific vocabulary
    EBP/domain pack owns domain debt vocabulary
    Solvent does not interpret debt strings

==================================================
DEFERRED CAPABILITIES
==================================================

Explicitly retain the frozen non-goals, including:

    DSSE/SLSA
    cryptographic attestation
    second real Domain Pack
    probabilistic inference
    source reputation
    W3C PROV
    persistent idempotency infrastructure
    persistent ArtifactRegistry
    crash-recoverable projection journal
    other explicitly deferred capabilities in v1.1

Do not introduce new capabilities here.

==================================================
GROWTH GATE EXCEPTION
==================================================

Record the single explicit planned kernel/API exception:

    POST /v1/beliefs/{parent_id}/edges

Purpose:

    canonical creation of belief_edge relationships

Required structural checks:

    parent exists
    child exists
    parent != child
    kind ∈ {derives, contradicts}
    uniqueness enforced

Explicit architectural rule:

    Coordinator MUST NOT write Solvent SQL directly.

The future Coordinator will call this Solvent endpoint.

This is a minimal Growth Gate exception, not a general reopening of the
Solvent freeze.

Do NOT implement this endpoint in Phase 1.

==================================================
REPOSITORY BOUNDARIES
==================================================

Record the locked boundaries:

    Solvent:
        canonical epistemic + authority state

    Conductor:
        operational workflow
        unchanged

    Oracle:
        ARGUS backend
        separate Go module

    Trust UI:
        separate Go module

    reference-loop:
        separate module
        untouched
        historical/reference implementation

==================================================
VERIFICATION RESULTS
==================================================

Record actual results from Phase 1.

Include:

    Solvent HEAD
    git status
    test result
    FullDebt search result
    migration verification

Explicitly preserve the Phase 0 CockroachDB limitation if still present.

Never convert an environment failure into a PASS.

==================================================
PHASE 1 ACCEPTANCE
==================================================

PASS only if:

    1. Actual current Solvent HEAD is recorded.
    2. Working tree has no modified tracked files.
    3. Plan 11.1 state is verified.
    4. FullDebt is absent from kernel.
    5. debt-opaque migration is present.
    6. Freeze reconciliation document is complete.
    7. Growth Gate exception is explicitly documented.
    8. No production code was changed.
    9. No Solvent endpoint was added.
    10. No Conductor code was changed.
    11. reference-loop was not changed.

If any architectural inconsistency is discovered:

    STOP.

Do not silently repair it.

==================================================
NO-CODE RULE
==================================================

Phase 1 may create/update ONLY:

    oracle/plan/freeze_reconciliation.md

Do not modify:

    Solvent source
    Solvent schema
    Conductor
    reference-loop
    oracle production code
    Trust UI

==================================================
FINAL RESPONSE
==================================================

Return:

    PHASE 1 RESULT: PASS / STOP

    Current freeze hash:
    ...

    Solvent verification:
    ...

    CockroachDB limitation:
    ...

    Freeze reconciliation:
    oracle/plan/freeze_reconciliation.md

    Production code changes:
    NONE

    Next phase:
    PHASE 2 — DOMAIN PACK SPECIFICATION

Do NOT begin Phase 2 in this run.