Use this as the handoff prompt for **Coding Agent B**. The key is that B is a **forensic continuation agent**, not a fresh implementer: first establish exactly what Agent A has already changed, then reconcile that work against the latest advisory, and only then implement missing pieces.

```text
You are Coding Agent B taking over an in-progress Reference Loop implementation that was started by Coding Agent A.

Your job is NOT to restart the work, redesign the architecture, or blindly continue from the original plan.

Your job is to:

1. inspect exactly what Agent A has already done;
2. determine what is complete, partially complete, incorrect, or still missing;
3. reconcile the implementation against the latest Reference Loop advisory below;
4. implement only the missing/corrective work that is justified;
5. preserve good work already completed by Agent A;
6. produce a precise handoff report with evidence.

============================================================
AUTHORITATIVE ARCHITECTURAL DIRECTION
============================================================

The project has deliberately pivoted from specification expansion to empirical proof.

The architectural north star remains:

    Agent
      ↓
    Conductor
      ↓
    Solvent
      ↓
    Executor
      ↓
    External System
      ↓
    Result
      ↓
    Conductor
      ↓
    Agent

Invariant:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Conductor = coordination.
Solvent = authority.
Executor = external effect.
External SOR = authoritative external result where one exists.
Agent = agency.
Domain = meaning/truth.

Do NOT create a new workflow runtime, workflow engine, workflow DB,
event bus, scheduler, gateway, second authority engine, generic
reconciliation daemon, or merged Conductor/Solvent architecture.

Integration is an implementation boundary/adapter, not a new role.

============================================================
LATEST IMPLEMENTATION DECISIONS
============================================================

The current Reference Loop decisions are:

1. Conductor interface:
   MCP

2. Solvent interface:
   REST

3. Reference operation:
   existing GitHub "deploy" action through the existing
   github_trigger_workflow Executor capability

4. Test repository:
   pithomlabs/reference-loop-test

5. GitHub workflow:
   .github/workflows/ref-loop.yml

6. Executor:
   - RecordingFunc for deterministic local/CI testing
   - real GitHub executor for external-effect integration proof

7. External SOR:
   GitHub workflow run state / GitHub API is authoritative for
   actual external effect occurrence.

8. Conductor DB:
   plain SQLite using modernc.org/sqlite
   file-backed test DB
   NO CockroachDB

9. Solvent DB:
   separate plain SQLite using modernc.org/sqlite
   file-backed test DB
   NO CockroachDB

10. Correlation:
    immutable run_id
    plus separate:
      scenario_id
      task_id
      intent_id
      operation_id

11. Evidence collector:
    runs after the execution path for reconstruction/correlation.
    It must NOT become part of the execution semantics.

12. First implementation scenario:
    happy_path first.
    Do not build all adversarial cases before the happy path works.

13. BM-IST:
    BM-IST is NOT to be silently treated as validated by the GitHub
    demo.
    If BM-IST infrastructure is not already available, explicitly
    record BM-IST validation as deferred rather than creating a
    second engineering project during this phase.

14. workflow.md / five persistent primitives:
    deferred architecture hypothesis, NOT active implementation.
    Do not introduce a persistent cross-role Intent/Work/Authorization/
    Effect/Outcome framework merely because workflow.md proposed it.

============================================================
LATEST ADVISORY THAT YOU MUST RECONCILE AGAINST
============================================================

Before broad implementation, the following minimum contract issues
must be explicitly closed or verified in the actual artifacts:

A. Operation identity comparison/equality semantics
B. Declaration/version binding
C. Generic capability declaration ownership
D. Pinned contract versions
E. Reference Loop pass criteria
F. Sandbox enforcement-path rule
G. "specification defect" failure classification
H. An explicit arbiter/owner for material finding classification

The intended methodology is:

    Phase 0 — contract closure
      ↓
    Phase 1 — pin experiment
      ↓
    Phase 2 — reference loop
      ↓
    Phase 3 — participant-owned evidence
      ↓
    Phase 4 — adversarial scenarios
      ↓
    Phase 5 — classify findings
      ↓
    Phase 6 — substitution
      ↓
    Phase 7 — conformance

Do NOT reopen broad specification design.

Close only the minimum security-critical contract needed to make the
experiment well-defined.

============================================================
IMPORTANT ADDITIONS FROM THE LATEST REVIEW
============================================================

1. Operation identity is not sufficiently tested merely because
   "X authorized, Y executed" fails.

   There must eventually be a named substitution scenario:

       cross_implementation_operation_identity

   Test whether independently implemented construction of the same
   semantic operation produces the same binding/equality result.

   Consider representation differences such as:
       JSON key ordering
       numeric representation
       string encoding
       casing
       normalization
       serialization

   This belongs in substitution testing, NOT in the first happy path.

2. The first consequential operation should not be a meaningless toy.

   The preferred criterion is:

   "smallest safe consequential operation that exercises at least
   one non-trivial architectural property while remaining
   independently observable and safely resettable."

   The selected GitHub deploy operation is acceptable because it
   exercises:
       - exact authorization binding
       - real external effect
       - independent observation
       - distinct external SOR
       - asynchronous workflow execution
       - an existing Executor capability

   Do NOT replace it merely for theoretical reasons unless actual
   repository inspection shows it cannot provide the required proof.

3. Sandbox rule:

   A sandbox may neutralize the final effect.

   A sandbox MUST NOT bypass:
       - authorization verification
       - exact operation binding
       - fail-closed enforcement

4. Failure taxonomy must include:

       implementation
       integration
       executor
       deployment
       specification defect
       new security property

   "specification defect" means the contract itself is wrong,
   incomplete, or ambiguous. It is not automatically an architecture
   expansion.

5. The coding agent must not unilaterally classify material findings
   to avoid architecture review.

   Each material finding must have:

       finding
       expected
       observed
       classification
       evidence
       owner/decision arbiter
       disposition
       fix
       rerun result

6. Version pinning matters.

   Determine the actual current versions/revisions of the normative
   Workflow Specification and Role/Boundary Matrix present in the
   repository.

   Do NOT assume their versions from memory.

   Record exactly what the Reference Loop is testing against.

   If a normative change affects the experiment, the affected
   scenario(s) must be rerun.

============================================================
PHASE 1 — FORENSIC INSPECTION OF AGENT A'S WORK
============================================================

Before changing anything, inspect the repository and working tree.

First locate the actual repositories/components involved. Do not assume
paths if the repository layout differs.

Inspect at minimum:

    git status
    git diff
    git diff --stat
    git log --oneline -n <reasonable range>
    untracked files
    staged files
    generated files
    test output
    README/runbook changes

Then inspect the actual implementation and relevant architecture docs.

At minimum look for:

    AGENTS.md
    Workflow Specification
    Role / Boundary Matrix
    Reference Loop plan
    reference-loop implementation
    Conductor implementation
    Solvent implementation
    Executor implementation
    MCP interfaces
    Solvent REST API
    SQLite setup/migrations
    GitHub Executor
    RecordingFunc
    evidence collector
    scenario code
    runbooks

Also inspect the latest review/advisory materials available in the
workspace if they exist.

Do not trust a retrospective document blindly.

Determine what the code ACTUALLY does.

============================================================
PHASE 2 — PRODUCE A GAP/STATUS MATRIX BEFORE MODIFYING CODE
============================================================

Create an internal or repository-local assessment with:

    Area
    Status
    Evidence
    Agent A changes
    Advisory requirement
    Gap
    Required action
    Risk

Use statuses:

    COMPLETE
    PARTIAL
    INCORRECT
    MISSING
    BLOCKED
    NOT APPLICABLE

At minimum assess:

    1. Contract closure
    2. Contract version pinning
    3. Conductor MCP integration
    4. Solvent REST integration
    5. SQLite/modernc setup
    6. Separate Conductor/Solvent DBs
    7. GitHub test repository integration
    8. GitHub deploy operation
    9. RecordingFunc path
   10. Real GitHub path
   11. Exact operation identity
   12. declaration/version binding
   13. authorization evidence
   14. fail-closed enforcement
   15. evidence ownership
   16. correlation model
   17. authoritative SOR semantics
   18. evidence collector separation
   19. happy path
   20. scenario pass criteria
   21. adversarial scenarios
   22. substitution
   23. cross-implementation identity test
   24. BM-IST scope statement
   25. specification-defect classification
   26. finding arbiter/disposition
   27. no-new-infrastructure constraint

DO NOT start by rewriting everything.

============================================================
PHASE 3 — VERIFY THE MINIMUM CONTRACT CLOSURE
============================================================

For each of the following, determine whether the repository already
contains an explicit, testable decision:

------------------------------------------------------------
A. Operation identity comparison semantics
------------------------------------------------------------

Find the actual identity definition and actual equality/comparison rule.

Questions:

    What fields define deploy identity?
    How are they serialized/canonicalized?
    What counts as equal?
    Is equality deterministic?
    Is the rule visible to independent implementations?

If this is already explicitly defined and implemented correctly:
    preserve it.

If missing or ambiguous:
    make the smallest contract-level correction required.

Do NOT invent a new generic canonicalization framework.

------------------------------------------------------------
B. Declaration/version binding
------------------------------------------------------------

Verify that authorization evidence and execution enforcement bind to
the correct declaration/identity version.

Specifically test/check:

    authorized under V1
    execute under incompatible V2

Expected:

    FAIL CLOSED

------------------------------------------------------------
C. Generic capability declaration ownership
------------------------------------------------------------

Verify who owns the deployed declaration when there is no distinct
integration/provider owner.

Do not invent a new role.

If the existing contract already resolves this, preserve it.

------------------------------------------------------------
D. Contract version pinning
------------------------------------------------------------

Determine the actual versions/revisions in use and record them.

The Reference Loop must have an explicit baseline.

------------------------------------------------------------
E. Reference Loop pass criteria
------------------------------------------------------------

Verify that the happy path has explicit, finite acceptance criteria.

At minimum:

    task created
    task claimed
    proposal created
    exact authorization
    exact execution attempt
    authoritative external effect
    result observed
    Conductor updated
    Agent continues
    participant-owned evidence reconstructs the chain

------------------------------------------------------------
F. Sandbox enforcement
------------------------------------------------------------

Verify that RecordingFunc does not bypass the actual authority/enforcement
boundary being tested.

If the fake executor is positioned in a way that makes tests pass without
testing authorization verification, exact binding, and fail-closed
behavior, fix that.

Do NOT make the sandbox "more realistic" by creating unnecessary
infrastructure.

------------------------------------------------------------
G. Specification defect classification
------------------------------------------------------------

Add the category if missing.

------------------------------------------------------------
H. Finding arbiter
------------------------------------------------------------

Identify who/what is responsible for disposition of material findings.

The coding agent may document findings but must not silently decide
that a potentially architectural finding is merely "integration."

============================================================
PHASE 4 — RECONCILE AGENT A'S IMPLEMENTATION
============================================================

After inspection, preserve everything that is valid.

For each current implementation piece, answer:

    Is it already correct?
    Is it compatible with the latest advisory?
    Does it need adjustment?
    Can it be left untouched?

Examples:

If Agent A already implemented:

    MCP Conductor client
    REST Solvent client
    SQLite modernc
    GitHub executor
    RecordingFunc
    happy path

do NOT rewrite them merely for stylistic reasons.

If something is structurally correct but lacks a small required piece,
add the piece.

If something conflicts with the architecture, explain exactly why before
changing it.

============================================================
PHASE 5 — COMPLETE THE HAPPY PATH
============================================================

Only after the investigation/gap analysis:

Ensure the happy path can actually execute:

    Agent
      ↓ MCP
    Conductor
      ↓ coordination
    Agent
      ↓ REST
    Solvent
      ↓ authorization evidence
    Executor
      ↓
    GitHub workflow
      ↓
    GitHub authoritative SOR
      ↓
    result
      ↓
    Conductor
      ↓
    Agent

Do not make the evidence collector part of the execution path.

============================================================
PHASE 6 — EVIDENCE
============================================================

Participant-owned evidence must remain authoritative.

Expected evidence ownership should remain approximately:

    TASK_CREATED
        Conductor

    TASK_CLAIMED
        Conductor

    PROPOSAL_CREATED
        Agent / actual API interaction

    AUTHORIZED(X)
        Solvent

    EXECUTION_ATTEMPTED(X)
        Executor

    EFFECT_CONFIRMED(X)
        GitHub SOR

    RESULT_OBSERVED
        Agent / appropriate coordination record

    CONDUCTOR_UPDATED
        Conductor

    AGENT_CONTINUED
        Agent

The harness may correlate these.

It must not manufacture authoritative facts.

============================================================
PHASE 7 — DO NOT PREMATURELY BUILD ALL ADVERSARIAL SCENARIOS
============================================================

The immediate goal is:

    inspect A
    fix contract prerequisites
    make happy path demonstrably correct
    stop and report

Do NOT automatically implement the entire adversarial suite unless the
repository state and happy path already justify doing so.

If the happy path is incomplete, finish only what is necessary for the
first proof.

Adversarial scenarios should come afterward, one controlled step at a
time.

============================================================
PHASE 8 — FUTURE TESTS THAT MUST BE TRACKED
============================================================

Even if not implemented yet, record these as explicit pending work:

    wrong operation
    missing authorization
    stale authorization
    verification unavailable
    declaration/version mismatch
    replay/duplicate
    execution failure
    ambiguous outcome
    termination/revocation
    Conductor cancellation
    Agent misinterpretation
    cross-implementation operation identity consistency

The last one is NEW and mandatory for the substitution phase.

============================================================
PHASE 9 — BM-IST
============================================================

Do not silently claim BM-IST validation.

Determine from the actual workspace whether a safe, bounded,
reproducible BM-IST operation already exists and can realistically be
integrated without creating a second engineering project.

If yes:
    report whether it should replace or follow the current GitHub
    reference loop.

If no:
    preserve the GitHub Reference Loop and explicitly record:

        BM-IST validation = deferred

Do not build BM-IST infrastructure merely to satisfy this review.

============================================================
PHASE 10 — ARCHITECTURAL CHANGE RULE
============================================================

Do NOT change the architecture merely because something is difficult.

Classify findings as:

    implementation
    integration
    executor
    deployment
    specification defect
    new security property

Interpretation:

    implementation:
        bug in the code

    integration:
        incorrect wiring/boundary interaction

    executor:
        executor capability/behavior problem

    deployment:
        environment/configuration problem

    specification defect:
        contract is ambiguous, contradictory, or wrong

    new security property:
        genuinely new responsibility/security invariant not covered
        by the current architecture

Only the last category should automatically raise the architecture
question.

A specification defect reopens the contract, not necessarily the
architecture.

============================================================
PHASE 11 — PRESERVE THE ARCHITECTURAL NON-GOALS
============================================================

Do not introduce:

    workflow engine
    workflow runtime
    workflow DB
    event bus
    scheduler
    heartbeat service
    generic gateway
    central effect registry
    mandatory verification SDK
    generic reconciliation daemon
    merged Conductor/Solvent
    second authority engine
    new Solvent kernel functionality

Do not modify Conductor/Solvent/Executor source merely to make the demo
prettier.

If an existing component has a real defect, make the smallest justified
change and document it.

============================================================
PHASE 12 — FINAL REPORT TO THE HUMAN REVIEWER
============================================================

At the end, provide a precise report with these sections:

1. EXECUTIVE STATUS
   - what Agent A had completed
   - what was missing
   - what you changed
   - current overall state

2. AGENT A FORENSIC FINDINGS
   - files inspected
   - existing implementation
   - valid work preserved
   - problematic work identified

3. CONTRACT STATUS
   For each:
       operation identity comparison
       declaration/version
       capability ownership
       version pin
       pass criteria
       sandbox enforcement
       specification-defect taxonomy
       finding arbiter

   Mark:
       CLOSED
       OPEN
       BLOCKED

4. IMPLEMENTATION STATUS
   Conductor MCP
   Solvent REST
   SQLite modernc
   GitHub executor
   RecordingFunc
   evidence collector
   correlation
   happy path

5. CHANGES MADE
   exact files and reasons

6. TESTS RUN
   exact commands
   exact results

7. REMAINING WORK
   ordered by priority

8. ARCHITECTURAL FINDINGS
   Only actual findings.
   Explicitly distinguish:
       implementation
       integration
       executor
       deployment
       specification defect
       new security property

9. BM-IST STATUS
   proven / deferred / blocked
   with reason

10. RECOMMENDATION
    State exactly what should happen next.

============================================================
MOST IMPORTANT RULE
============================================================

Do not assume Agent A is wrong.

Do not assume Agent A is right.

Establish the truth from:

    repository state
    source code
    tests
    runtime behavior
    existing artifacts
    participant-owned evidence

Your first responsibility is forensic verification.

Your second responsibility is the smallest corrective implementation.

Your third responsibility is to leave the repository in a state where
the human reviewer can clearly see:

    DONE
    PARTIAL
    MISSING
    INCORRECT
    BLOCKED
    NEXT

Do not restart the project.
Do not overwrite good work.
Do not create a second implementation beside Agent A's implementation.
Do not turn this handoff into another architecture-design exercise.

The goal is to recover the actual state of the Reference Loop, close the
minimal security-critical gaps identified by the latest advisory, and
get to a trustworthy empirical proof.
```

This framing deliberately makes **Agent B an auditor-plus-continuation agent** rather than another greenfield builder. It should first establish whether Agent A already solved any of the issues the later review identified, because we should not "fix" problems that are already resolved in the actual repository. The latest review specifically warns that some retrospective claims may have outrun the recorded artifacts, so repository evidence should win over narrative. 

The most important handoff distinction is: **B should close only the minimum contract gaps, then continue the empirical loop; B should not restart the specification phase.** 

