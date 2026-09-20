# Reference Loop Continuation — FROZEN SOLVENT Implementation Plan

## Solvent Freeze Forensic Report

```
SOLVENT FREEZE CHECK
--------------------
Repository:          /home/chaschel/Documents/go/solvent-main
Frozen baseline:     7602699 ("✨ solvent kernel freeze", 2026-09-08 18:51 +0800)
Current HEAD:        7e5ca2d ("📝 AGENTS.md")
Current branch:      main
Working tree clean:  YES
Source diff from frozen baseline:  YES — 2 source files changed
Untracked files:     none
Staged changes:      none
Relevant commits after freeze:
    7e5ca2d 📝 AGENTS.md
    1ddb045 🐛 adv_review11
Behavior-affecting dependency/config changes:  none
Verdict: MODIFIED (2 commits after freeze; 2 source files changed)
```

### Exact Post-Freeze Source Changes

**`api/auth.go`** (+29 lines):
- Added `verifyPrincipalActive(ctx, db, principalID)` — a best-effort liveness check that queries `principal` table for revoked status.
- Imports added: `database/sql`, `errors`, `github.com/PithomLabs/solvent/kernel`

**`kernel/kernel.go`** (+6 lines):
- Comment tweak in `EnterBelief`: clarified nil-to-empty guard rationale.
- No behavioral change.

**All other changes** are documentation, plans, images, or test files.

### Assessment

The two post-freeze commits do **not** appear related to the Reference Loop work. Commit messages and file changes indicate general Solvent development. No Reference Loop adapter, schema, or integration code is present.

**However**, the frozen baseline commit `7602699` exists and is well-defined. The current HEAD contains source modifications beyond that baseline. Per the frozen-Solvent rule, this state requires human decision before proceeding.

---

## Required Decision

**Choose one:**

**Option A — Accept current HEAD as the frozen baseline**
- Treat `7e5ca2d` as the effective frozen state.
- The 2 source changes are minor and not Reference Loop related.
- Proceed with implementation using current HEAD.

**Option B — Revert to `7602699` as the frozen baseline**
- Reset Solvent to the explicit "solvent kernel freeze" commit.
- Guarantees a clean, named freeze point.
- Requires discarding the 2 post-freeze commits.

**Option C — STOP / BLOCKED**
- Do not proceed until the freeze boundary is formally defined.

---

## Revised Implementation Plan (Frozen Solvent)

### Architectural Constraint

```
THE REFERENCE LOOP MUST ADAPT TO FROZEN SOLVENT.
FROZEN SOLVENT MUST NOT BE MODIFIED TO FIT THE REFERENCE LOOP.
```

**Non-negotiable rules:**
- Solvent source = READ-ONLY. Zero edits.
- Solvent must run unchanged against its existing CockroachDB/pgx path.
- No SQLite port, no schema changes, no API changes, no auth middleware changes.
- No fork, no compatibility branch, no temporary patch.

### Environment

```
Conductor  →  file-backed SQLite (modernc.org/sqlite)   [existing Conductor impl]
Solvent    →  isolated ephemeral CockroachDB             [existing Solvent impl, unchanged]
Executor   →  RecordingFunc / GitHub (unchanged)
```

### Ordered Task List

#### Task 1: Forensic Assessment (reference-loop repo)
- Produce the gap/status matrix for Agent A's work.
- Document findings in `FORENSIC.md` under `reference-loop/`.

#### Task 2: Contract Pin Artifact
- Create `contract.json` in `reference-loop/` root.
- Pin: workflow spec v0.3, role/boundary matrix v0.5, operation identity rule, declaration owner/version, pass criteria, sandbox rule, failure taxonomy, finding arbiter, BM-IST deferral.
- Add `contract_test.go`: assert file exists, parseable, required fields non-empty.

#### Task 3: Provision CockroachDB for Solvent
- Provision an isolated/local/ephemeral CockroachDB instance.
- Apply Solvent's existing schema migrations exactly as the repository expects.
- Do NOT modify Solvent's schema or migrations.
- Record the exact Solvent commit/version used.
- If Solvent cannot run unchanged against CockroachDB in this environment, STOP and report BLOCKED.

#### Task 4: Start Frozen Solvent
- Start Solvent using its existing REST/API path (`cmd/solvent-api/main.go`).
- Use real `api.AuthMiddleware` with static `keyToPrincipal` mapping.
- Wire `audit.New`, `policy.New`, `authority.New`, `ledger.New`.
- Register GitHub executor if `GITHUB_TOKEN` is provided.
- No Solvent source modifications.

#### Task 5: Fix Conductor Client to Use MCP
- Replace `conductor/client.go` direct SQLite calls with MCP `CallTool` invocations.
- Wire `mcp.NewServer` with Conductor's SQLite store.
- Map operations: `conductor_create_task`, `conductor_claim_task`, `conductor_submit_task`, `conductor_post_activity`, `conductor_next_task`, `conductor_get_task`.
- Remove raw SQL from `conductor/client.go`.

#### Task 6: Fix Config
- `SolventDBPath` should be a CockroachDB DSN, not SQLite path.
- `SolventRESTAddr` should point to the provisioned Solvent server.
- `ConductorDBPath` remains file-backed SQLite.

#### Task 7: Fix Operation Identity
- Declare: `deploy:pithomlabs/reference-loop-test:ref-loop.yml:main:<run_id>`.
- Bind proposal, authorization, execution, and evidence to this identity.
- Include `run_id` as executor input.

#### Task 8: Fix Evidence Collection
- `EFFECT_CONFIRMED` only recorded after querying GitHub API for workflow run state (GitHub mode).
- `RecordingFunc` events labeled `non_external_effect_proof=true`.
- Evidence collector runs AFTER execution (not on critical path).
- Participant tables are authoritative; harness correlation is test evidence only.

#### Task 9: Add Tests
- `contract_test.go` — contract pin present and parseable.
- `operation_identity_test.go` — deterministic equality, includes run_id.
- `authorization_test.go` — RecordingFunc does not bypass auth; declaration/version mismatch fails closed.
- `happy_path_test.go` — evidence sequence: TASK_CREATED, TASK_CLAIMED, PROPOSAL_CREATED, AUTHORIZED, EXECUTION_ATTEMPTED, EFFECT_CONFIRMED/simulated, RESULT_OBSERVED, CONDUCTOR_UPDATED, AGENT_CONTINUED.

#### Task 10: Run Test Suite
```bash
cd /home/chaschel/Documents/go/oracle/reference-loop
go test ./...
go vet ./...

cd /home/chaschel/Documents/go/solvent-main
go test ./...   # establish frozen baseline passes its own tests
go vet ./...
```

#### Task 11: Record Pending Work
Document in `contract.json` under `pending_work` (do not implement):
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

#### Task 12: BM-IST Deferral
- Mark BM-IST validation as deferred in `contract.json`.
- No BM-IST implementation in this pass.

---

## Risk Register

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| CockroachDB not available in environment | Medium | High | Use ephemeral/local instance; stop and report BLOCKED if unavailable |
| Solvent post-freeze commits cause incompatibility | Low | Medium | Review diff; if changes are behavioral, escalate to owner |
| GitHub API rate limits | Low | Medium | Use RecordingFunc for CI; real GitHub for integration proof |
| Conductor MCP in-process wiring | Low | Low | Use `CallTool` directly |
| Evidence collector timing | Low | Low | Run after execution completes |

---

## Stop Conditions

STOP immediately and report rather than improvising if:

- Solvent frozen baseline cannot be established.
- Solvent has been modified and ownership of those changes is unclear.
- Solvent requires source changes to run.
- Solvent cannot run unchanged against its existing CockroachDB path.
- The Reference Loop requires bypassing the Solvent authorization path.
- Conductor can only be made to work by bypassing its MCP interface.
- Authoritative evidence would have to be fabricated.
- The task starts expanding into a new runtime, workflow engine, event bus, gateway, or authority layer.

---

## Validation Commands

```bash
# Solvent: establish frozen baseline passes its own tests
cd /home/chaschel/Documents/go/solvent-main
go test ./...
go vet ./...

# Reference Loop
cd /home/chaschel/Documents/go/oracle/reference-loop
go test ./...
go vet ./...
```

---

## Final Verdict Format (To Be Completed After Implementation)

```
PASS / BLOCKED / PASS WITH LIMITATIONS

- Solvent frozen: YES/NO
- Solvent source modified during this task: NO
- Frozen baseline verified: YES/NO
- Conductor interface used: MCP
- Solvent interface used: REST / api.AuthMiddleware
- Solvent DB: CockroachDB / ...
- Happy path proven: YES/NO
- Real external effect proven: YES/NO
- Outstanding blockers: ...
```

---

## Open Question

**Which frozen baseline should be used for Solvent?**

- Option A: Accept current HEAD `7e5ca2d` as frozen.
- Option B: Revert to `7602699` ("solvent kernel freeze") as frozen.
- Option C: STOP until formal decision is made.

My recommendation is **Option B**: revert to `7602699`. It provides a named, unambiguous freeze point, and the two subsequent commits (`adv_review11`, `AGENTS.md`) are not required for the Reference Loop. This preserves the strongest possible empirical claim that the Reference Loop exercises a frozen Solvent implementation unchanged.
