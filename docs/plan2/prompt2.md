You are implementing the Reference Loop continuation work described in the attached plans/reviews.

CRITICAL ARCHITECTURAL CONSTRAINT — FROZEN SOLVENT

Solvent is a frozen codebase.

You MUST NOT modify the Solvent source code under any circumstances during this task.

This prohibition is absolute:
- no SQLite port
- no CockroachDB → SQLite compatibility changes
- no schema changes
- no API changes
- no auth middleware changes
- no kernel/service changes
- no “small compatibility fix”
- no fork of Solvent with modified source
- no temporary patch that is later reverted
- no generated-code modification that changes behavior
- no test-only source modification inside Solvent
- no vendor/dependency modification that effectively changes Solvent behavior

The Reference Loop must adapt its environment and adapters to the frozen Solvent implementation, not the other way around.

======================================================================
1. FIRST: FORENSICALLY VERIFY WHETHER SOLVENT WAS MODIFIED
======================================================================

Before implementing anything, determine whether the current Solvent codebase differs from its frozen baseline.

You must explicitly verify this rather than assuming it is unchanged.

Locate:
- the Solvent repository
- the frozen baseline commit/tag/hash if documented
- git status
- current branch
- recent commits
- diff against the frozen baseline
- untracked files
- staged changes
- submodule/vendor changes if applicable
- generated files that may have changed
- dependency/configuration changes that materially alter Solvent behavior

Use Git history and repository metadata where available.

Produce a forensic report:

SOLVENT FREEZE CHECK
--------------------
Repository:
Frozen baseline:
Current HEAD:
Current branch:
Working tree clean: YES/NO
Source diff from frozen baseline: YES/NO
Untracked files: ...
Staged changes: ...
Relevant commits after freeze: ...
Behavior-affecting dependency/config changes: ...
Verdict: FROZEN / MODIFIED / BASELINE UNKNOWN

IMPORTANT:
If you find modifications relative to the frozen baseline, DO NOT overwrite, revert, or “fix” them automatically.

Instead:
1. preserve the current state;
2. identify exactly what changed;
3. determine whether the changes are part of the existing working tree/history;
4. report whether they appear related to the Reference Loop work;
5. STOP before modifying Solvent.

A modified Solvent repository is a BLOCKER until the human owner decides what to do.

If the frozen baseline itself cannot be established with sufficient confidence, STOP and report:
BLOCKED: FROZEN SOLVENT BASELINE CANNOT BE VERIFIED.

Do not infer that “no obvious changes” means frozen.

======================================================================
2. HARD STOP RULE FOR SOLVENT
======================================================================

From this point onward:

SOLVENT SOURCE = READ-ONLY.

You may:
- inspect it
- compile it
- run it
- test it
- configure its external environment
- provision its existing database
- invoke its existing APIs
- write Reference Loop adapters outside the Solvent repository
- write integration tests outside the Solvent repository

You may NOT:
- edit Solvent files
- patch Solvent SQL
- change Solvent schemas
- change Solvent authentication
- change Solvent REST behavior
- change Solvent kernel semantics
- change Solvent database driver
- create a compatibility branch/fork
- alter Solvent to accommodate SQLite

If frozen Solvent cannot run in the intended environment unchanged, STOP and report BLOCKED.

Do not solve the problem by changing Solvent.

======================================================================
3. USE THE EXISTING SOLVENT PERSISTENCE PATH
======================================================================

The Reference Loop must run frozen Solvent using its existing persistence architecture.

Therefore:

- Solvent uses its existing pgx/CockroachDB path.
- Provision an isolated local/ephemeral CockroachDB instance for the Reference Loop.
- Do not substitute SQLite for Solvent.
- Do not port Solvent SQL to SQLite.
- Do not introduce a second Solvent persistence implementation.

Conductor may continue using its own file-backed SQLite implementation.

It is acceptable, and desirable, for:

    Conductor → SQLite
    Solvent   → CockroachDB

because they are separate persistence boundaries.

The Reference Loop must prove the actual architecture rather than flattening infrastructure differences for convenience.

======================================================================
4. DO NOT MODIFY CONDUCTOR OR EXECUTOR SEMANTICS UNNECESSARILY
======================================================================

Use the real participant interfaces.

Conductor:
- interact through its MCP interface;
- do not bypass MCP with direct database manipulation in the Reference Loop;
- use external adapters only where necessary for setup/observation.

Solvent:
- interact through its existing REST/API surface;
- use the existing authentication middleware and authorization path;
- do not recreate Solvent authorization logic in the harness.

Executor:
- use the existing Executor capability where possible;
- RecordingFunc is acceptable for deterministic/local proof;
- real GitHub execution may be used for external-effect proof;
- do not manufacture authoritative execution/effect evidence.

Do not make semantic changes to participant implementations merely to make the demo work.

======================================================================
5. REFERENCE LOOP SCOPE
======================================================================

Implement only the smallest real end-to-end happy path.

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

First scenario:

happy_path

Do NOT implement the full adversarial suite yet.

Do NOT build the Conformance Test Matrix yet.

Do NOT expand the architecture.

======================================================================
6. OPERATION IDENTITY
======================================================================

Preserve exact operation identity across:

- proposal
- authorization
- execution
- result/evidence
- reconstruction

Use a single immutable top-level:

run_id

Keep these distinct:

run_id
scenario_id
task_id
intent_id
operation_id

Do not collapse them into one identifier.

Authorization and execution must refer to the exact operation/target/snapshot required by the existing Solvent semantics.

======================================================================
7. AUTHORIZATION BOUNDARY
======================================================================

The Reference Loop must demonstrate:

AUTHORIZE ≠ EXECUTE

The executor must not be treated as an authority engine.

The harness must not bypass authorization.

A successful call to an executor does not itself prove authorization.

The execution path must fail closed when the required authorization evidence is absent, invalid, stale, mismatched, or unverifiable.

Sandboxing may neutralize the final external effect, but it must NOT bypass:

- authorization verification
- exact operation binding
- fail-closed enforcement

======================================================================
8. EVIDENCE
======================================================================

Evidence must be reconstructed from participant-owned records.

Do not manufacture authoritative events from harness assumptions.

The evidence collector must NOT be on the critical execution path.

It may run after or alongside execution to reconstruct what happened.

Clearly distinguish:

- real authoritative participant evidence
- harness correlation evidence
- simulated/local execution evidence
- real external-effect evidence

Never label a RecordingFunc-only event as real external EFFECT_CONFIRMED.

For example:

RecordingFunc:
    local execution attempt / simulated effect

GitHub executor:
    external effect evidence, subject to the declared source of record

The GitHub workflow run state should be treated as the source of record for the real external-effect portion where applicable.

======================================================================
9. HAPPY-PATH TRACE
======================================================================

The first proof should establish a trace equivalent to:

TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED(X)
→ AUTHORIZED(X)
→ EXECUTION_ATTEMPTED(X)
→ EFFECT_CONFIRMED(X) / correctly marked simulated
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED

The exact event representation must come from the actual participants and existing contracts.

Do not invent a new event system merely to display this sequence.

======================================================================
10. TEST REPOSITORY
======================================================================

Use:

pithomlabs/reference-loop-test

The test repository/workspace may contain:

- adapters
- harness code
- environment configuration
- test fixtures
- evidence reconstruction
- integration tests
- orchestration scripts
- documentation

It must NOT contain a modified Solvent implementation.

======================================================================
11. DATABASE SETUP
======================================================================

Provision separate databases:

Conductor:
    file-backed SQLite using its existing implementation.

Solvent:
    isolated/ephemeral CockroachDB using its existing frozen implementation.

Do not attempt to unify these databases.

The database setup is part of the Reference Loop environment, not a reason to alter either participant.

======================================================================
12. FAILURE CLASSIFICATION
======================================================================

When something fails, classify it before changing anything:

1. implementation defect
2. integration defect
3. executor defect
4. deployment/environment defect
5. specification defect
6. genuinely new security property

Do NOT classify “Solvent currently uses CockroachDB” as a defect.

Do NOT classify “SQLite would be more convenient” as a defect.

Do NOT change architecture merely to eliminate environmental inconvenience.

A genuinely frozen participant incompatibility is:

BLOCKED

unless it can be solved entirely outside the frozen participant.

======================================================================
13. REQUIRED FORENSIC / IMPLEMENTATION REPORT
======================================================================

Before declaring success, produce:

A. Solvent Freeze Verification
- frozen baseline
- current HEAD
- git status
- exact diff status
- whether any post-freeze modifications exist
- whether any changes appear related to this work

B. Environment
- Conductor version/commit
- Solvent version/commit
- CockroachDB version
- test repository commit
- relevant configuration

C. Interface Path
- Agent → Conductor interface used
- Conductor → Solvent interface used
- Solvent → Executor interface used
- external source of record

D. Run Identity
- run_id
- scenario_id
- task_id
- intent_id
- operation_id

E. Happy-Path Evidence
- participant-owned evidence
- authorization evidence
- execution evidence
- result evidence
- Conductor observation

F. Boundary Verification
Explicitly confirm:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

and:

AUTHORIZE ≠ EXECUTE

G. Known limitations
Anything not proven by the first loop must be explicitly listed.

======================================================================
14. STOP CONDITIONS
======================================================================

STOP immediately and report rather than improvising if:

- Solvent frozen baseline cannot be established.
- Solvent has been modified and ownership of those changes is unclear.
- Solvent requires source changes to run.
- Solvent cannot run unchanged against its existing CockroachDB path.
- The Reference Loop requires bypassing the Solvent authorization path.
- Conductor can only be made to work by bypassing its MCP interface.
- authoritative evidence would have to be fabricated.
- the proposed fix requires changing Solvent semantics.
- the task starts expanding into a new runtime, workflow engine, event bus, gateway, or authority layer.

Do not “make progress” by violating a stop condition.

======================================================================
15. FINAL VERDICT FORMAT
======================================================================

End with exactly one of:

PASS
BLOCKED
PASS WITH LIMITATIONS

And provide:

- Solvent frozen: YES/NO
- Solvent source modified during this task: YES/NO
- Frozen baseline verified: YES/NO
- Conductor interface used: ...
- Solvent interface used: ...
- Solvent DB: CockroachDB / ...
- Happy path proven: YES/NO
- Real external effect proven: YES/NO
- Outstanding blockers: ...

The most important invariant is:

THE REFERENCE LOOP MUST ADAPT TO FROZEN SOLVENT.
FROZEN SOLVENT MUST NOT BE MODIFIED TO FIT THE REFERENCE LOOP.
