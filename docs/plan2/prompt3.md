You are implementing the Reference Loop continuation plan.

THIS IS A FROZEN-SOLVENT TASK.

The Solvent codebase is intended to be frozen. Your first responsibility is to verify that freeze state and protect it.

======================================================================
0. ABSOLUTE SOLVENT RULE
======================================================================

DO NOT MODIFY SOLVENT.

The Solvent repository is:

    /home/chaschel/Documents/go/solvent-main

You may inspect, build, test, run, and invoke Solvent.

You may NOT:
- edit Solvent source files
- edit Solvent migrations
- edit Solvent schemas
- edit Solvent API behavior
- edit Solvent authentication
- edit Solvent kernel/service logic
- port Solvent to SQLite
- change pgx/CockroachDB persistence
- create a Solvent compatibility branch
- create a Solvent fork with behavioral changes
- make a temporary patch and revert it later
- modify generated Solvent code
- modify dependencies specifically to alter Solvent behavior

The Reference Loop must adapt to Solvent.

Solvent must not be adapted to the Reference Loop.

======================================================================
1. FIRST ACTION: FORENSICALLY VERIFY SOLVENT FREEZE
======================================================================

Before implementing anything else, inspect the Solvent repository.

Determine:

- repository path
- current branch
- current HEAD
- working-tree state
- staged changes
- untracked files
- frozen baseline commit, if documented
- commits after the frozen baseline
- source-file differences from the frozen baseline
- dependency/configuration differences
- whether any changes appear related to Reference Loop work

Run appropriate Git commands and inspect the relevant diffs.

Produce a report in the Reference Loop repository:

    FORENSIC.md

The report MUST contain:

    SOLVENT FREEZE CHECK
    --------------------
    Repository:
    Frozen baseline:
    Current HEAD:
    Current branch:
    Working tree clean:
    Source diff from frozen baseline:
    Untracked files:
    Staged changes:
    Post-freeze commits:
    Behavior-affecting dependency/config changes:
    Verdict:

Then provide:

    EXACT POST-FREEZE SOURCE CHANGES

For every changed Solvent source file:
- path
- commit
- lines changed
- behavioral or non-behavioral
- relevance to Reference Loop

Do not merely report that files changed.
Inspect the actual diff.

======================================================================
2. CRITICAL: DO NOT MAKE THE FREEZE DECISION YOURSELF
======================================================================

A critical distinction:

"Solvent is currently modified relative to the named freeze"

is a forensic fact.

"Accept current HEAD as the new freeze"

or

"Revert current HEAD to the old freeze"

is an ownership/governance decision.

YOU MUST NOT make that decision yourself.

Therefore:

If the current Solvent HEAD differs from the documented frozen baseline:

    DO NOT revert it.
    DO NOT accept it as the new baseline.
    DO NOT rewrite history.
    DO NOT tag it.
    DO NOT modify it.

Instead:

    STOP.

Report:

    BLOCKED: SOLVENT FREEZE BOUNDARY REQUIRES HUMAN DECISION

Provide the exact evidence needed for the owner to choose the baseline.

This is true even if the changes appear harmless.

Do not rationalize a new freeze because the changes are "minor".

Do not revert commits merely because they appear unnecessary.

======================================================================
3. ONLY PROCEED WHEN THE FREEZE BOUNDARY IS EXPLICIT
======================================================================

Proceed only when one of these is already established by the repository/workflow:

A. The documented frozen baseline is the current HEAD.

OR

B. The repository has already been intentionally restored to the documented frozen commit.

Do not perform the transition yourself unless explicitly instructed to do so.

Once the freeze is established:

    Solvent source = READ ONLY

From that point forward, verify before and after the implementation that no Solvent files changed.

======================================================================
4. PROTECT AGAINST ACCIDENTAL SOLVENT MODIFICATION
======================================================================

Before implementation, capture:

    git rev-parse HEAD
    git status --short
    git diff --stat
    git diff --name-only

After implementation, repeat the same checks.

Also compare the final tree against the approved frozen baseline.

The final report MUST explicitly state:

    Solvent source modified during this task: NO

If any Solvent source file changes during the task:

    STOP

Do not revert the change automatically.

Report the exact file and diff.

======================================================================
5. SOLVENT DATABASE
======================================================================

When implementation is authorized to proceed:

Solvent MUST use its existing persistence architecture:

    pgx + CockroachDB

Provision an isolated/local/ephemeral CockroachDB instance for the Reference Loop.

Use Solvent's existing migrations exactly as they are.

Do not modify those migrations.

Do not port them to SQLite.

Do not create a compatibility layer inside Solvent.

If frozen Solvent cannot run unchanged against the existing CockroachDB path:

    STOP
    report BLOCKED

Do not solve the problem by modifying Solvent.

======================================================================
6. CONDUCTOR DATABASE
======================================================================

Conductor remains independent:

    Conductor → file-backed SQLite

This is not a reason to unify persistence.

The Reference Loop should preserve the real architectural boundary:

    Conductor → SQLite
    Solvent   → CockroachDB

======================================================================
7. REFERENCE LOOP INTERFACES
======================================================================

Use the actual participant boundaries.

Conductor:
- interact through MCP
- no direct database access from the Reference Loop client
- use MCP CallTool
- preserve the existing Conductor implementation

Solvent:
- interact through its existing REST/API surface
- use the real Solvent auth middleware/path
- do not reimplement authorization in the Reference Loop

Executor:
- use the existing executor capability
- RecordingFunc may be used for deterministic local proof
- GitHub may be used for real external-effect proof
- do not create a new executor solely for the demonstration

======================================================================
8. FIRST SCOPE: HAPPY PATH ONLY
======================================================================

Implement only the smallest useful end-to-end path.

Do NOT implement the entire adversarial matrix yet.

Do NOT implement the Conformance Test Matrix yet.

Do NOT add new infrastructure roles.

Reference path:

Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
external effect
  ↓
result
  ↓
Conductor
  ↓
Agent

Initial scenario:

    happy_path

======================================================================
9. CONTRACT PIN
======================================================================

Create:

    reference-loop/contract.json

Pin the currently approved versions/decisions for:

- Workflow Specification v0.3
- Role / Boundary Matrix v0.5
- operation identity rule
- declaration ownership/version
- Reference Loop pass criteria
- sandbox enforcement-path rule
- failure taxonomy
- finding arbiter
- BM-IST deferral

Add a test proving the contract file exists, parses, and has all required fields.

Do not silently change normative contracts during implementation.

======================================================================
10. OPERATION IDENTITY
======================================================================

Use:

    run_id
    scenario_id
    task_id
    intent_id
    operation_id

These are distinct identifiers.

The operation identity must be deterministic and consistent across:

- proposal
- authorization
- execution
- evidence
- result

Initial operation:

    deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>

Include run_id in the executor input where required.

Do not collapse all identifiers into one field.

======================================================================
11. AUTHORIZATION
======================================================================

Prove:

    AUTHORIZE ≠ EXECUTE

The Reference Loop MUST use the real Solvent authorization path.

Do not bypass authorization because the executor is a test executor.

The executor must not become an authority engine.

The system must fail closed if authorization is absent, invalid, stale, mismatched, or unverifiable.

======================================================================
12. EVIDENCE
======================================================================

Participant-owned evidence is authoritative.

Harness-generated correlation data is not authoritative.

The evidence collector must NOT be on the critical execution path.

It reconstructs evidence afterward.

RecordingFunc evidence MUST be labeled clearly as simulated/local/non-external-effect evidence.

Do NOT claim:

    EFFECT_CONFIRMED

from RecordingFunc alone.

For GitHub mode, establish external effect from the GitHub source of record.

======================================================================
13. HAPPY-PATH EVIDENCE
======================================================================

Prove a trace equivalent to:

TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED(X)
→ AUTHORIZED(X)
→ EXECUTION_ATTEMPTED(X)
→ EFFECT_CONFIRMED(X) / simulated
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED

Use actual participant records.

Do not invent an event bus or new event system simply to represent this sequence.

======================================================================
14. IMPLEMENTATION TASKS
======================================================================

After the freeze boundary is formally established, perform:

1. Reference Loop forensic/gap assessment.
2. contract.json + contract test.
3. CockroachDB provisioning for frozen Solvent.
4. Start frozen Solvent using its existing REST/API path.
5. Use the real Solvent auth middleware.
6. Change the Reference Loop Conductor client to MCP if necessary.
7. Keep Conductor on SQLite.
8. Keep Solvent on CockroachDB.
9. Add run_id and exact operation identity.
10. Add evidence reconstruction.
11. Add happy-path authorization/execution tests.
12. Run all required tests and static checks.
13. Record pending adversarial work only; do not implement it yet.
14. Record BM-IST as deferred.

All code changes should be confined to the Reference Loop/test environment unless an existing participant implementation has a demonstrated defect.

Solvent is excluded from modification.

======================================================================
15. TESTING
======================================================================

Run:

    cd /home/chaschel/Documents/go/oracle/reference-loop
    go test ./...
    go vet ./...

For Solvent, ONLY inspect/test the frozen tree:

    cd /home/chaschel/Documents/go/solvent-main
    go test ./...
    go vet ./...

Do not alter Solvent to make its tests pass.

======================================================================
16. FINAL SOLVENT INTEGRITY CHECK
======================================================================

Before declaring success:

Verify again:

    git rev-parse HEAD
    git status --short
    git diff --name-only
    git diff --stat

Compare against the approved frozen baseline.

Explicitly report:

    Solvent frozen: YES/NO
    Frozen baseline:
    Solvent source modified during this task: YES/NO
    Solvent source files changed:
    Solvent behavior changed:
    Reference Loop source changes:
    Final working tree status:

======================================================================
17. FINAL VERDICT
======================================================================

Return exactly one:

    PASS
    BLOCKED
    PASS WITH LIMITATIONS

Use BLOCKED when the freeze boundary is unresolved.

Do not report PASS merely because the Reference Loop itself compiles.

A valid PASS requires:

- Solvent freeze boundary explicitly established
- Solvent unchanged throughout the task
- frozen Solvent runs on its existing CockroachDB path
- Conductor uses its intended MCP boundary
- Solvent uses its intended REST/auth boundary
- authorization is not bypassed
- operation identity is preserved
- participant-owned evidence is preserved
- happy path is demonstrated

======================================================================
MOST IMPORTANT RULE
======================================================================

DO NOT "FIX" SOLVENT.

If something in Solvent prevents the Reference Loop from working:

    STOP
    REPORT THE BLOCKER
    DO NOT MODIFY SOLVENT