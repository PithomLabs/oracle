  # Reference Loop Continuation Plan — Strict SQLite

  ## Summary

  Treat this document as the execution order. Do not implement anything until the forensic gap/status matrix has been produced. The matrix is the baseline for deciding what to preserve, fix, defer, or block.

  Continue as Coding Agent B from the existing Agent A work in /home/chaschel/Documents/go/oracle/reference-loop. Do not restart or rebuild the harness from scratch. Preserve valid work, document what was wrong, and make the smallest
  corrective changes.

  Observed baseline:

  - oracle/reference-loop already exists and compiles, but has no tests.
  - Conductor repo is clean and already uses SQLite/modernc.org/sqlite.
  - Solvent repo is clean, but its current kernel/service/API path is pgx/Cockroach-shaped.
  - Normative pins present in Oracle are:
      - Workflow Specification: plans/workflow_specification_v0.3_final.md
      - Role / Boundary Matrix: plans/role_boundary_matrix_v0.5.md

  ## Key Changes

  - First produce a repository-local forensic assessment for Agent A’s work with the required gap/status matrix.
  - Add a small Reference Loop contract closure artifact pinning the normative versions, operation identity rule, declaration owner/version, pass criteria, sandbox rule, failure taxonomy including specification defect, finding arbiter,
    and BM-IST deferral.

  - Implement strict Solvent SQLite support only as narrow storage compatibility:
      - Use modernc.org/sqlite.
      - Preserve existing Solvent authority semantics.
      - Do not modify Conductor behavior.
      - For Solvent, permit only the storage compatibility changes required to run existing authority semantics on SQLite.
      - Do not alter authority semantics, security invariants, API contracts, operation identity semantics, or introduce a generic database abstraction.
      - Any change beyond that narrow compatibility boundary requires stopping and reporting BLOCKED.

  - Correct the harness to use real component interfaces:
      - Conductor via MCP tools, not direct DB calls.
      - Solvent via REST with Solvent’s real api.AuthMiddleware.
      - Separate file-backed SQLite DBs for Conductor and Solvent.
      - RecordingFunc path for deterministic proof; GitHub path for real external-effect proof.

  - Correct GitHub deploy operation identity:
      - operation: deploy
      - repo: pithomlabs/reference-loop-test
      - workflow: .github/workflows/ref-loop.yml
      - ref: main
      - include all effect-relevant inputs, especially run_id, if required by the workflow.

  - Preserve evidence collector separation:
      - Participant records are authoritative.
      - Harness correlation is test evidence only.
      - Do not manufacture EFFECT_CONFIRMED; real GitHub mode must query GitHub workflow run state/API.
      - RecordingFunc mode may prove local execution attempt, but must be labeled non-external-effect proof.

  ## Test Plan

  - Run:
      - rtk go test ./... in oracle/reference-loop
      - focused Conductor MCP/task lifecycle tests
      - focused Solvent REST/authority/execution/SQLite compatibility tests

  - Add focused tests for:
      - contract pin artifact present and parseable
      - deterministic declared operation identity equality
      - declaration/version mismatch fails closed
      - RecordingFunc does not bypass authorization verification, exact binding, or fail-closed enforcement
      - happy path produces the required evidence sequence

  Do not implement the full adversarial suite now. Record these as pending work: wrong operation, missing authorization, stale authorization, verification unavailable, declaration/version mismatch, replay/duplicate, execution failure,
  ambiguous outcome, termination/revocation, Conductor cancellation, Agent misinterpretation, and cross_implementation_operation_identity.

  ## Assumptions And Stop Rules

  - The Reference Loop must use strict SQLite/modernc.org/sqlite with separate Conductor and Solvent DBs.
  - The GitHub deploy operation remains the first consequential operation unless repository evidence proves it cannot satisfy the proof.
  - BM-IST validation is deferred.
  - No workflow runtime, workflow DB, event bus, scheduler, gateway, reconciliation daemon, merged Conductor/Solvent, or new Solvent kernel functionality.
  - If SQLite compatibility requires broad architectural surgery, stop and report BLOCKED with evidence.
  - Final report must include executive status, Agent A forensic findings, contract status, implementation status, exact changes, tests run, remaining work, architectural findings, BM-IST status, and next recommendation.

