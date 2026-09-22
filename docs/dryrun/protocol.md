# Implementation Plan: Agent Identity + Common Agent Protocol

**Version:** 1.0
**Date:** 2026-09-20
**Status:** PLAN — NOT IMPLEMENTED

---

## Executive Summary

Implement two tightly related improvements:

1. Explicitly identify the AI agent/harness/model performing Work or Adversarial research, and expose that identity in ARGUS state/UI.
2. Create one shared protocol of agent Do's and Don'ts that both Work and Adversarial agents must follow.

**Key audit finding:** The `Agent` struct already exists in `packet/v1/types.go` and `schema.json`, but `App.Persist()` silently discards it. The implementation must complete the persistence path.

---

## Audit Results

| Question | Finding |
|----------|---------|
| Does packet support agent object? | **Yes** — `Agent{ID, Model, Role}` exists in types.go and schema.json |
| Where is agent identity persisted? | **Nowhere** — `App.Persist()` discards it |
| What is `conductor_task.current_agent`? | Agent that **claimed** a task (different concept) |
| Why `<nil>` in UI? | `CurrentAgent` is `*string`, nil by default, template renders raw |
| Is identity accepted but not supplied? | **Yes** — prompts don't include it |

---

## Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Agent struct fields | Add `Harness` | Distinguishes runtime (OpenCode vs Hermes) using same model |
| Persistence | New `packet_submission` table | Clean separation from task claiming; append-only provenance |
| Validation | All fields required (fail-closed) | Agent provenance is part of the audit trail |
| Identity source | Human-supplied in prompt | Agent must not self-declare provenance |
| Schema enforcement | `agent.role == packet.role` | Prevents audit-trail contradiction |

---

## Canonical Agent Identity

```json
{
  "id": "work-001",
  "role": "work",
  "harness": "OpenCode",
  "model": "MiMo-V2.5"
}
```

**Invariants:**
- Agent identity is provenance/observability metadata only
- Agent identity is NOT authority, NOT a principal, NOT an authorization credential
- One packet_id maps to one immutable submission record
- Existing `conductor_task.current_agent` semantics remain unchanged

---

## Phase 1: Packet Schema Changes

### 1.1 Add Harness field to Agent struct

**File:** `packet/v1/types.go`

```go
// Agent identifies the creator of the packet.
type Agent struct {
    ID      string `json:"id"`
    Role    string `json:"role"`
    Harness string `json:"harness"`
    Model   string `json:"model"`
}
```

Remove `omitempty` tags — all fields are now required.

### 1.2 Update JSON schema

**File:** `packet/v1/schema.json`

Change the `agent` property to:

```json
"agent": {
  "type": "object",
  "required": ["id", "role", "harness", "model"],
  "properties": {
    "id": { "type": "string", "minLength": 1 },
    "role": { "type": "string", "enum": ["work", "adversarial"] },
    "harness": { "type": "string", "minLength": 1 },
    "model": { "type": "string", "minLength": 1 }
  }
}
```

### 1.3 Add agent validation to packet validator

**File:** `packet/v1/validate.go`

Add `validateAgent(pkt *Packet) error`:

- `agent.id` must be non-empty
- `agent.role` must be non-empty
- `agent.harness` must be non-empty
- `agent.model` must be non-empty
- `agent.role` must equal `pkt.Role`

Call from `Validate()` after role check (step 3).

---

## Phase 2: Persistence

### 2.1 New migration

**File:** `internal/migrations/003_packet_submission.sql`

```sql
-- Packet-level agent provenance (append-only).
-- Records which agent/harness/model produced each submitted research packet.
-- Distinct from conductor_task.current_agent (task claiming).

CREATE TABLE IF NOT EXISTS packet_submission (
    packet_id      UUID PRIMARY KEY,
    scenario_id    UUID NOT NULL,
    task_id        UUID,
    agent_id       STRING NOT NULL,
    role           STRING NOT NULL,
    harness        STRING NOT NULL,
    model          STRING NOT NULL,
    content_sha256 STRING NOT NULL,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

**Invariant:** Append-only. Once written, never updated. Same packet_id = same record (idempotent via PK).

### 2.2 Persist agent identity in App.Persist

**File:** `internal/application/app.go`

After writing idempotency row (step 5), add step 6:

```go
// 6. Persist packet submission provenance (agent identity).
taskRef := sql.NullString{String: pkt.TaskRef, Valid: pkt.TaskRef != ""}
_, err = tx.ExecContext(ctx,
    `INSERT INTO packet_submission
         (packet_id, scenario_id, task_id, agent_id, role, harness, model, content_sha256)
     VALUES ($1::UUID, $2::UUID, $3, $4, $5, $6, $7, $8)
     ON CONFLICT (packet_id) DO NOTHING`,
    pkt.PacketID, pkt.ScenarioID, taskRef,
    pkt.Agent.ID, pkt.Agent.Role, pkt.Agent.Harness, pkt.Agent.Model, hash)
if err != nil {
    return nil, fmt.Errorf("insert packet submission: %w", err)
}
```

**Note:** `task_id` uses `sql.NullString` because packets may exist without a direct task reference.

---

## Phase 3: UI Projection

### 3.1 Add PacketSubmissionView and query

**File:** `internal/epistemic/view.go`

Add struct:

```go
type PacketSubmissionView struct {
    PacketID    string `json:"packet_id"`
    ScenarioID  string `json:"scenario_id"`
    AgentID     string `json:"agent_id"`
    Role        string `json:"role"`
    Harness     string `json:"harness"`
    Model       string `json:"model"`
    SubmittedAt string `json:"submitted_at"`
}
```

Add function:

```go
func GetAllSubmissions(ctx context.Context, db *sql.DB) ([]PacketSubmissionView, error)
```

Query: `SELECT packet_id, scenario_id, agent_id, role, harness, model, submitted_at FROM packet_submission ORDER BY submitted_at DESC`

### 3.2 Extend Dashboard type

**File:** `internal/application/app.go`

```go
type Dashboard struct {
    Tasks       []*work.Task
    Beliefs     []epistemic.BeliefView
    Submissions []epistemic.PacketSubmissionView
}
```

Update `GetDashboard()` to populate `Submissions`.

### 3.3 Update Insights template

**File:** `internal/ui/templates/insights.html`

Fix `<nil>` display — replace raw CurrentAgent with:

```html
<td>{{if .CurrentAgent}}{{.CurrentAgent}}{{else}}—{{end}}</td>
```

Add Recent Submissions table after Epistemic State:

```html
<h2>Recent Submissions</h2>
<table border="1">
<tr><th>Packet</th><th>Agent</th><th>Harness</th><th>Model</th><th>Role</th><th>Time</th></tr>
{{range .Submissions}}
<tr>
<td>{{.PacketID}}</td><td>{{.AgentID}}</td><td>{{.Harness}}</td><td>{{.Model}}</td><td>{{.Role}}</td><td>{{.SubmittedAt}}</td>
</tr>
{{end}}
</table>
```

---

## Phase 4: Common Agent Protocol

### 4.1 Create shared protocol file

**File:** `prompts/argus-agent-protocol.md`

Contains the full common agent protocol:

- Identity block (AGENT_ID, AGENT_ROLE, AGENT_HARNESS, AGENT_MODEL)
- Required first actions (read task, get_context, reconstruct state, read corpus)
- Agents MUST (use get_context, use submit_packet, preserve belief IDs, create new IDs for new propositions, etc.)
- Agents MUST NOT (access DB directly, promote, discharge, retract, fabricate evidence, etc.)
- Belief lifecycle rule (new proposition = new belief ID + graph edge)
- Adversarial rule (challenge, don't replace)
- Human authority rule (agents propose, humans decide)
- Provenance rule (background docs != authoritative evidence)
- Unknown rule (never convert UNKNOWN to TRUE)

### 4.2 Update Work Agent prompt

**File:** `prompts/opencode-work.md`

- Add identity block at top with placeholder values the human fills in
- Add reference: "Read and follow prompts/argus-agent-protocol.md"
- Update packet example to include `agent` object
- Retain role-specific behavior (decomposition, evidence gathering, etc.)

### 4.3 Update Adversarial Agent prompt

**File:** `prompts/opencode-adversarial.md`

- Add identity block at top with placeholder values the human fills in
- Add reference: "Read and follow prompts/argus-agent-protocol.md"
- Update packet example to include `agent` object
- Retain role-specific behavior (independent challenge, counterexamples, etc.)

---

## Phase 5: Tests

### 5.1 Agent validation tests

**File:** `packet/v1/validate_test.go`

- Test: valid packet with all agent fields → passes
- Test: missing agent.id → rejected
- Test: missing agent.role → rejected
- Test: missing agent.harness → rejected
- Test: missing agent.model → rejected
- Test: agent.role != packet.role → rejected

### 5.2 Agent identity propagation test

**File:** `internal/application/app_test.go`

- Test: submit work packet with agent identity → verify packet_submission row persisted with correct fields
- Test: submit adversarial packet with different identity → verify distinct agent_id/role preserved

### 5.3 Authority isolation test

**File:** `internal/application/app_test.go`

- Test: agent identity cannot be used as human operator principal
- Test: agent identity does not grant discharge/promote/retract permissions

### 5.4 Existing tests remain green

Run `go test ./...` — all existing tests must pass without modification (except updating packet construction to include required agent fields where needed).

### 5.5 Integration test update

**File:** `integration_test.go`

Update packet construction in `TestFullIntegration` to include agent identity. Verify submission appears in dashboard.

---

## File Change Summary

| File | Change |
|------|--------|
| `packet/v1/types.go` | Add Harness field, remove omitempty |
| `packet/v1/schema.json` | Make all agent fields required |
| `packet/v1/validate.go` | Add validateAgent() |
| `packet/v1/validate_test.go` | Add agent validation tests |
| `internal/migrations/003_packet_submission.sql` | New migration |
| `internal/application/app.go` | Persist agent identity, extend Dashboard |
| `internal/application/app_test.go` | Agent propagation + authority isolation tests |
| `internal/epistemic/view.go` | Add PacketSubmissionView + GetAllSubmissions |
| `internal/ui/templates/insights.html` | Fix nil display, add submissions table |
| `prompts/argus-agent-protocol.md` | New shared protocol file |
| `prompts/opencode-work.md` | Add identity block + protocol reference |
| `prompts/opencode-adversarial.md` | Add identity block + protocol reference |
| `integration_test.go` | Update packet construction |

---

## Flow Demonstration

After implementation:

```
Work Agent
  identity (work-001 / OpenCode / MiMo-V2.5)
    |
    v
argus.submit_packet (agent object in packet)
    |
    v
ARGUS persists to packet_submission table
    |
    v
Insights shows: work-001 / OpenCode / MiMo-V2.5

Adversarial Agent
  different identity (adversarial-001 / OpenCode / MiMo-V2.5)
    |
    v
argus.submit_packet (agent object in packet)
    |
    v
ARGUS persists to packet_submission table
    |
    v
Insights shows: adversarial-001 / OpenCode / MiMo-V2.5
```

The architectural invariant remains:

```
CAPABILITY != WORK != AUTHORITY != EXECUTION
```

Agent identity improves provenance and observability.
It must never become authority.
