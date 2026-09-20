Phase 1.5 is complete and accepted as:

    PASS WITH LIMITATIONS

Do not reopen the workflow design.

Do not modify frozen Solvent.

Do not implement Phase 2 yet.

Your next task is ONLY to prepare Phase 2 for implementation.

======================================================================
CURRENT BASELINE
======================================================================

Frozen Workflow Design:

    Workflow Design v1.0

Frozen Solvent:

    /home/chaschel/Documents/go/solvent-main
    HEAD = 7602699

Phase 1.5 result:

    PASS
        operation mismatch
        duplicate delivery
        stale/invalid authorization
        stale claim × expired intent
        operation identity equality
        Solvent freeze integrity

    NOT_YET_PROVEN
        UNKNOWN != DENIED semantics
        capability_ref → declaration resolution

    New architecture required:
        NO

Growth Gate:
    NO

No real external effect has yet been performed.

======================================================================
OBJECTIVE
======================================================================

Produce a concrete, minimal, pre-registered Phase 2 plan for the first
real external-effect experiment.

Do not implement the effect yet.

The purpose is to make the Phase 2 experiment sufficiently specified that
implementation becomes mechanical rather than another design exercise.

======================================================================
1. PHASE2_PLAN.md
======================================================================

Create:

    PHASE2_PLAN.md

It must define exactly:

A. First real consequential operation

Choose ONE operation already available through the existing integration path.

Prefer the existing GitHub workflow/deploy operation used by the Reference Loop.

Do not invent a new Executor.

B. External source of record

Define exactly which external system is authoritative for determining whether
the effect occurred.

For GitHub, identify the authoritative workflow/run state that will be used.

C. Capability declaration

Define the minimum declaration required for this operation:

    capability_ref
    declaration owner
    declaration version
    content hash
    effective reference
    operation class
    effect-capable = true
    operation identity definition
    authorization evidence requirements
    validity model
    replay/idempotency behavior
    execution outcome model
    reconciliation owner/mechanism

Do not create a general capability registry.

Implement only the declaration needed for the first experiment.

D. Plan scope

Define the Human-approved operation class/scope for the experiment.

Use:

    plan_id
    plan_version
    operation class
    approved target scope

Do NOT enumerate every future exact operation.

E. Exact operation identity

Define the exact canonical operation identity used by the first experiment.

Use the already-proven form where applicable:

    deploy:repo:workflow:ref:run_id

State exactly which fields are effect-relevant.

F. Authorization evidence

Define:

    what Solvent returns/provides
    how the Executor obtains it
    how the Executor verifies it
    trust basis
    validity requirements
    mismatch behavior

Do NOT add cryptography unless the actual integration requires it.

G. External effect

Define precisely what constitutes:

    EXECUTION_ATTEMPTED
    EFFECT_CONFIRMED
    RESULT_OBSERVED
    AMBIGUOUS

Do not equate executor invocation with external effect.

H. Safety model

The real external-effect experiment MUST be:

    reversible
    sandboxed
    dry-run
    isolated test repository

or otherwise safely bounded.

Safety neutralization MUST occur only at or after:

    authorization verification
    exact operation binding
    fail-closed enforcement

Do not bypass the production authorization/enforcement path.

======================================================================
2. PRE-REGISTER PASS CRITERIA
======================================================================

Before implementation, write explicit criteria.

Phase 2 PASS requires at least:

1. Human-approved plan is recorded.
2. Exact operation is constructed from the declared operation class.
3. Solvent authorizes the exact operation.
4. Authorization evidence reaches the effect-capable boundary.
5. Executor independently verifies required authorization evidence.
6. Authorization binding is exact.
7. Executor produces the intended external effect.
8. External SOR provides authoritative evidence of effect occurrence.
9. Operation identity remains correlated end-to-end.
10. Conductor records coordination state without becoming the authority/effect
    ledger.
11. No frozen Solvent source is modified.

======================================================================
3. PRE-REGISTER FAILURE / FALSIFICATION CRITERIA
======================================================================

Explicitly define what would falsify the current architecture.

Examples:

    - exact authorized operation cannot be safely distinguished from a
      different operation;
    - external effect can occur without valid Solvent authorization;
    - Executor cannot independently enforce authorization evidence;
    - external SOR cannot distinguish success from ambiguity under the
      declared model;
    - Conductor must become an authority engine to make the loop safe;
    - Agent-only behavior is required to preserve a consequential boundary;
    - plan scope cannot be recorded without turning Conductor into a policy
      engine.

Do not predeclare ordinary implementation bugs as architecture failures.

======================================================================
4. HANDLE CURRENT NOT_YET_PROVEN ITEMS
======================================================================

Address the two Phase 1.5 findings explicitly.

A. UNKNOWN != DENIED

Do NOT modify frozen Solvent.

State exactly:

    Current frozen Solvent behavior:
        fail closed with Allowed=false

    Required workflow semantic:
        UNKNOWN != DENIED

Decide whether Phase 2 can safely proceed without this distinction.

If it cannot, mark Phase 2 BLOCKED.

Do not silently reinterpret the current behavior.

B. capability_ref → declaration resolution

Define the minimum integration-level mechanism required for the first real
operation.

It must provide:

    owner
    version
    content hash
    effective reference
    retrievable content

It must not become a general Conductor capability registry.

======================================================================
5. PIN COMPONENTS
======================================================================

Record exact versions/commits for:

    Workflow Design v1.0
    Conductor
    Reference Loop
    Solvent = 7602699
    Executor/integration
    declaration version
    external test repository/target

Do not begin the real effect until the versions are recorded.

======================================================================
6. EVIDENCE PACKAGE
======================================================================

Define the minimum participant-owned evidence required:

    Human:
        approved plan

    Conductor:
        plan/task/claim/coordination state

    Solvent:
        authorization decision/evidence

    Executor:
        execution attempt/outcome

    External SOR:
        authoritative effect evidence

    Agent:
        proposal/interpretation/continuation

The harness only correlates observations.

It does not manufacture authoritative state.

======================================================================
7. PHASE 2 SCENARIO
======================================================================

Define one concrete scenario.

Preferred shape:

    Human approves Plan v1
        ↓
    Agent enters WORK
        ↓
    Agent claims READY task
        ↓
    Agent constructs exact GitHub operation
        ↓
    Solvent authorizes exact operation
        ↓
    Executor verifies authorization
        ↓
    Executor triggers GitHub workflow
        ↓
    GitHub becomes external SOR
        ↓
    Result is observed
        ↓
    Conductor records coordination result
        ↓
    Agent continues

Do not add additional scenarios yet.

======================================================================
8. PHASE 2 STOP RULE
======================================================================

After producing:

    PHASE2_PLAN.md

STOP.

Do not:

    implement the real external effect
    modify Solvent
    modify Conductor architecture
    build a new capability registry
    build a new workflow engine
    create a new Executor
    build a conformance harness
    integrate BM-IST

This phase is planning only.

======================================================================
9. FINAL OUTPUT
======================================================================

Return:

    Phase 2 preparation: READY / BLOCKED

    First operation:
    External SOR:
    Declaration pinned:
    Plan scope defined:
    Operation identity defined:
    Authorization evidence defined:
    Safety mechanism:
    Pass criteria registered:
    Failure/falsification criteria registered:

    UNKNOWN != DENIED:
        ACCEPTED FOR PHASE 2 / BLOCKING

    capability_ref resolution:
        READY / BLOCKING

    Solvent modified:
        NO

    Architecture change required:
        NO

Do not declare Phase 2 READY unless every required prerequisite above is
explicitly resolved.