# Reference Loop v0.1 — Implementation Plan

## 1. Architecture Summary (from reconnaissance)

**Conductor (MCP interface):**
- Task lifecycle: `proposed` → `active` → `review` → `accepted`
- `ClaimTask` (atomic), `SubmitTask`, `PostActivity`, `GetGovernance` via MCP
- Tasks carry `GovernanceRef` (opaque JSON) pointing to Solvent belief

**Solvent (REST API):**
- Authority flow: `CreateTarget` → `RequestAuthorization` → `Approve` → `AuthorizeAction` (creates intent) → `ExecuteAction` (claims intent, invokes executor)
- Executor registry: `github_trigger_workflow` maps to `deploy` action
- GitHub executor triggers workflow with `repo`, `workflow`, `ref`, `inputs` from snapshot
- Authority evidence: intent created, claim recorded, executor outcome logged in audit

**External Result Source:**
- **GitHub = authoritative SOR** — workflow run state on GitHub is the external truth
- Executor reports `RunID`; GitHub API confirms actual execution outcome

---

## 2. Reference Operation Selection

| Property | Value |
|----------|-------|
| **Action** | `deploy` (existing) |
| **Executor** | `github_trigger_workflow` (registered in `actionExecutorMap`) |
| **Target** | Dedicated test repo (e.g., `pithomlabs/reference-loop-test`) |
| **Workflow** | Simple workflow (e.g., `.github/workflows/ref-loop.yml` — creates artifact or writes file) |
| **Ref** | `main` branch |
| **Operation Identity** | `deploy:repo:workflow:ref` — survives proposal→authorization→execution |

---

## 3. Repository Structure

```
/home/chaschel/Documents/go/oracle/reference-loop/
├── README.md
├── go.mod
├── main.go                    # run-reference-loop entry point
├── config/
│   └── config.go              # Executor selection, endpoints, test repo
├── scenario/
│   ├── happy_path/
│   │   └── scenario.go
│   ├── wrong_operation/
│   │   └── scenario.go
│   ├── denied/
│   ├── stale_auth/
│   ├── unavailable/
│   ├── decl_version_mismatch/
│   ├── replay/
│   ├── execution_failure/
│   ├── ambiguous/
│   ├── termination_revocation/
│   ├── cancellation/
│   └── misinterpretation/
├── agent/
│   └── agent.go               # MCP-based agent driver
├── conductor/
│   └── client.go              # MCP client wrapper
├── solvent/
│   └── client.go              # REST client wrapper
├── executor/
│   ├── registry.go            # Configurable executor registry
│   ├── recording.go           # RecordingFunc (default)
│   └── github.go              # GitHub provider wrapper
├── evidence/
│   ├── collector.go           # Correlates participant-owned evidence
│   ├── types.go               # Evidence record types
│   └── printer.go             # Human-readable trace output
└── runbook/
    ├── local.md               # Run with RecordingFunc
    └── integration.md         # Run with real GitHub executor
```

---

## 4. Phase 0 — Setup (Before Implementation)

1. **Create test GitHub repository** with simple workflow:
   - `.github/workflows/ref-loop.yml` — runs `echo "ref-loop-$(date +%s)" > result.txt` and uploads artifact
   - Public or accessible via `GITHUB_TOKEN`

2. **Configure Solvent** with target for this operation:
   - `CreateTarget` with consequence_params: `{"repo": "...", "workflow": "ref-loop.yml", "ref": "main"}`
   - `RequestAuthorization` → `Approve` → produces snapshot

3. **Seed Conductor** with task referencing the Solvent belief

---

## 5. Happy Path Implementation

**Sequence:**
```text
TASK_CREATED           [Conductor evidence: task row + activity]
→ TASK_CLAIMED         [Conductor evidence: claim transition]
→ PROPOSAL_CREATED(X)  [Agent evidence: authorize_action call + returned intent]
→ AUTHORIZED(X)        [Solvent evidence: intent row + audit log]
→ EXECUTION_ATTEMPTED(X) [Executor evidence: ClaimIntent + adapter_invoked]
→ EFFECT_CONFIRMED(X)  [SOR evidence: GitHub workflow run + artifact]
→ RESULT_OBSERVED      [Agent evidence: ExecuteAction response]
→ CONDUCTOR_UPDATED    [Conductor evidence: activity + status transition]
→ AGENT_CONTINUED      [Agent evidence: next task discovery]
```

**Key implementation points:**
- Single binary `run-reference-loop` with flags: `-scenario=happy_path`, `-executor=recording|github`
- Agent driver uses Conductor MCP `CallTool` for all interactions
- Solvent REST client calls `/v1/authorizations/action` then `/v1/authorizations/execute`
- Evidence collector polls each system's audit/ledger and correlates by correlation ID
- No centralized ledger — each participant's evidence remains in their store

---

## 6. Adversarial Scenarios (after happy path)

| Scenario | Test | Expected | Classification if fails |
|----------|------|----------|------------------------|
| Wrong operation | Authorize `deploy` on `repo=A`, execute on `repo=B` | FAIL CLOSED | integration/executor |
| Missing auth | Execute without valid intent | FAIL CLOSED | integration |
| Stale auth | Revoke target, then execute | FAIL CLOSED | integration |
| Unavailable | Stop Solvent, attempt execute | FAIL CLOSED | deployment |
| Decl/version mismatch | Authorize v1 snapshot, execute against v2 | FAIL CLOSED | new security property* |
| Replay | Execute same intent twice | Controlled by kernel | implementation |
| Executor failure | RecordingFunc returns error | FAILURE stays | implementation |
| Ambiguous | Executor returns ambiguous outcome | AMBIGUOUS | executor |
| Termination | Cancel intent via ReconcileIntent | TERMINATED ≠ AMBIGUOUS | implementation |
| Conductor cancel | Cancel task in Conductor | NOT revoke Solvent auth | integration |
| Misinterpretation | Agent reads wrong result | Authoritative result unchanged | implementation |

*Declaration/version mismatch is the only potential new security property — if the kernel already handles this, it's integration.

---

## 7. Substitution Proof (after adversarial tests pass)

| Substitution | Method | Verification |
|--------------|--------|--------------|
| Agent A → Agent B | Different MCP client implementation | Same semantic contract |
| RecordingFunc → GitHub executor | Config flag change | External effect proof |
| MCP → REST for Conductor | Use Conductor REST API instead of MCP | Same task lifecycle |

---

## 8. Evidence Ownership Design

Each evidence item records:
```go
type EvidenceRecord struct {
    Event           string                 // TASK_CREATED, AUTHORIZED, etc.
    Owner           string                 // "conductor", "solvent", "executor", "external", "agent"
    Source          string                 // "conductor.activity", "solvent.audit", "github.api", etc.
    OperationID     string                 // deploy:repo:workflow:ref
    DeclarationVer  string                 // snapshot_id from Solvent
    Timestamp       time.Time
    CorrelationID   string                 // task_id + scenario_id + intent_id
    Raw             json.RawMessage        // Participant's own record
}
```

**Collector** reads from:
- Conductor: `activity` table, `task` transitions
- Solvent: `audit` table, `intent` state, `authorization_decision`
- Executor: audit log entries (`adapter_invoked`, `executor_completed`, `executor_failed`)
- GitHub: workflow run API (external SOR)

---

## 9. Executable Command

```bash
# Local deterministic run (default)
cd /home/chaschel/Documents/go/oracle/reference-loop
go run . -scenario=happy_path -executor=recording

# Integration run (real external effect)
GITHUB_TOKEN=xxx go run . -scenario=happy_path -executor=github -test-repo=owner/repo

# Adversarial
go run . -scenario=wrong_operation -executor=recording
```

---

## 10. Key Implementation Files to Create

| File | Purpose |
|------|---------|
| `main.go` | CLI entry, scenario dispatch, evidence collection |
| `config/config.go` | All endpoints, executor selection, test repo config |
| `agent/agent.go` | MCP-based task lifecycle driver |
| `conductor/client.go` | MCP `CallTool` wrapper for task operations |
| `solvent/client.go` | REST client for authorize/execute |
| `executor/registry.go` | Switchable registry (RecordingFunc ↔ GitHub) |
| `evidence/collector.go` | Correlates evidence from all participants |
| `evidence/printer.go` | Prints the required trace format |
| `scenario/*/scenario.go` | Each scenario's setup/execute/verify |

---

## 11. Files to NOT Modify

- Conductor, Solvent, Executor source code
- No new workflow engine, event bus, scheduler, gateway
- No BM-IST infrastructure

---

## 12. Success Criteria (Gates)

1. **Gate 1** — Happy path produces full trace with participant-owned evidence
2. **Gate 2** — All adversarial scenarios fail closed with correct classification
3. **Gate 3** — Failure/ambiguity/termination remain distinguished
4. **Gate 4** — Evidence ownership traceable to correct participant
5. **Gate 5** — Agent substitution + Executor substitution preserve contract

---

## 13. Questions Before Starting

1. **Test repository**: Should I create `pithomlabs/reference-loop-test` with the workflow, or do you have one?
2. **Conductor DB**: Should the reference loop use an in-memory SQLite or connect to existing CockroachDB?
3. **Solvent DB**: Same question — existing CockroachDB or test database?
4. **Correlation ID scheme**: Use `task_id:scenario_id:intent_id` as composite key?
5. **First scenario**: Start with `happy_path` only, or scaffold all scenario directories first?
