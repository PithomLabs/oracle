  # Reference Loop Continuation Plan — Strict SQLite

  ## Summary

  Continue as Coding Agent B from the existing Agent A work in /home/chaschel/Documents/go/oracle/reference-loop. Do not restart the implementation. Preserve the current harness shape where valid, but correct it against the latest
  advisory.

  Observed baseline:

  - oracle/reference-loop already exists and compiles, but has no tests.
  - Conductor repo is clean and already uses SQLite/modernc.org/sqlite.
  - Solvent repo is clean, but its current kernel/service/API path is pgx/Cockroach-shaped.
  - Normative pins present in Oracle are:
      - Workflow Specification: plans/workflow_specification_v0.3_final.md
      - Role / Boundary Matrix: plans/role_boundary_matrix_v0.5.md

  - Agent A implemented a monolithic happy-path harness, but it currently uses direct Conductor SQLite access instead of MCP, attempts Solvent auth context injection with an incompatible private key type, and has conflicting Solvent
    SQLite/Cockroach setup.

  ## Key Changes

  - First produce a repository-local forensic assessment for Agent A’s work with the required gap/status matrix. Treat oracle/reference-loop, solvent-main, and conductor as the evidence sources.
  - Add a small Reference Loop contract closure artifact pinning:
      - Workflow Spec v0.3 final
      - Role / Boundary Matrix v0.5
      - Reference Loop v0.1
      - operation identity rule
      - declaration owner/version
      - pass criteria
      - sandbox enforcement rule
      - failure taxonomy including specification defect
      - finding arbiter/disposition owner
      - BM-IST validation deferred

  - Implement strict Solvent SQLite support only as a narrow storage compatibility boundary needed for the Reference Loop:
      - Use modernc.org/sqlite.
      - Reuse existing Solvent authority semantics.
      - Do not redesign the kernel, add a generic DB abstraction, or build a second authority engine.
      - If preserving semantics on SQLite requires broad Solvent changes, stop and report BLOCKED.

  - Correct the harness to use real component interfaces:
      - Conductor via MCP tools, not direct DB client calls.
      - Solvent via REST with Solvent’s real api.AuthMiddleware.
      - Separate file-backed SQLite DBs for Conductor and Solvent.
      - RecordingFunc path remains deterministic; GitHub path uses github_trigger_workflow.

  - Correct operation identity for the GitHub deploy:
      - operation: deploy
      - repo: pithomlabs/reference-loop-test
      - workflow: .github/workflows/ref-loop.yml
      - ref: main
      - include all effect-relevant inputs, especially run_id, if the workflow requires it.
      - bind proposal, authorization, execution, and evidence to the same declared identity/version.

  - Preserve evidence collector separation:
      - Participant records are authoritative.
      - Harness correlation is test evidence only.
      - Do not manufacture EFFECT_CONFIRMED; real GitHub mode must query GitHub workflow run state/API.
      - RecordingFunc mode may prove local execution attempt, but must be labeled non-external-effect proof.

  ## Test Plan

  - Run existing repo checks:
      - rtk go test ./... in oracle/reference-loop
      - focused Conductor tests for MCP/task lifecycle
      - focused Solvent tests for authority, REST authorization, execution, and SQLite compatibility

  - Add focused Reference Loop tests for:
      - contract version pin artifact is present and parseable
      - operation identity equality is deterministic for the declared GitHub deploy operation
      - authorized under declaration V1 then executed under incompatible V2 fails closed
      - RecordingFunc does not bypass authorization verification, exact binding, or fail-closed enforcement
      - happy path produces finite required evidence:
        TASK_CREATED, TASK_CLAIMED, PROPOSAL_CREATED, AUTHORIZED, EXECUTION_ATTEMPTED, EFFECT_CONFIRMED or correctly marked simulated, RESULT_OBSERVED, CONDUCTOR_UPDATED, AGENT_CONTINUED

  - Do not implement the full adversarial suite now. Record pending scenarios:
      - wrong operation
      - missing authorization
      - stale authorization
      - verification unavailable
      - declaration/version mismatch
      - replay/duplicate
      - execution failure
      - ambiguous outcome
      - termination/revocation
      - Conductor cancellation
      - Agent misinterpretation
      - cross_implementation_operation_identity

  ## Assumptions And Stop Rules

  - Use strict SQLite for the Reference Loop even though Solvent is currently Cockroach-shaped.
  - Place compatibility work in the narrowest Solvent/reference-loop boundary that preserves existing semantics.
  - Do not modify Conductor/Solvent core behavior unless a verified implementation defect blocks the Reference Loop.
  - Do not introduce a workflow runtime, workflow DB, event bus, scheduler, gateway, reconciliation daemon, merged Conductor/Solvent, or new Solvent kernel functionality.
  - If SQLite compatibility cannot be achieved without broad architectural surgery, stop and report the incompatibility as BLOCKED, with evidence.
  - Final report must include executive status, Agent A forensic findings, contract status, implementation status, exact changes, tests run, remaining work, architectural findings, BM-IST status, and next recommendation.

