# Plan 3: Fix `/ui/insights` Nil Pointer Crash

**Status:** AWAITING APPROVAL — no implementation yet
**Date:** 2026-09-20
**Scope:** Trust UI Insights/Debts dashboard + template error handling

---

## Observed Failure

```
GET /ui/insights  (browser, no query params)
        ↓
template: insights.html:10:6:
executing "insights.html" at <.ID>:
nil pointer evaluating interface {}.ID
        ↓
server log: http: superfluous response.WriteHeader call
from internal/ui.(*Server).HandleInsights (ui.go:131)
```

Server starts successfully: CockroachDB ready, migrations applied, seed inserted, HTTP listening on :8080.

---

## Root Cause

### Primary: Wrong query model for dashboard

`HandleInsights` (ui.go:113) defaults `taskID` to `"default"` when no `task_id` query param is provided:

```go
func (s *Server) HandleInsights(w http.ResponseWriter, r *http.Request) {
    taskID := r.URL.Query().Get("task_id")
    if taskID == "" {
        taskID = "default"
    }
    ctx, err := s.app.GetContext(r.Context(), taskID)
    ...
```

`GetContext("default")` calls `workStore.GetByID(ctx, "default")`. No task has ID `"default"` — the seed creates task `00000000-0000-0000-0000-000000000101`.

When `GetByID` returns `ErrNotFound`:
- `result.Task` stays nil (app.go:507-512)
- Snapshot is skipped because `result.Task == nil` (app.go:522-534)
- Returns `Context{Task: nil, Snapshot: &epistemic.Snapshot{Beliefs: nil}}`

The handler then wraps nil in a slice:

```go
data := map[string]interface{}{
    "Tasks":   []interface{}{ctx.Task},    // []interface{}{nil}
    "Beliefs": ctx.Snapshot.Beliefs,        // nil
}
```

Template iterates `{{range .Tasks}}` → one element (nil) → `{{.ID}}` → nil pointer dereference.

### Why existing tests pass

`TestInsightsPage` (ui_test.go:235) explicitly creates a task and passes `?task_id=<real-uuid>`. The bug only manifests when visiting `/ui/insights` without a `task_id` — the default browser experience.

### Secondary: Unsafe error handling

Template execution writes to `ResponseWriter` directly via `s.tmpl.ExecuteTemplate(w, ...)`. If the template fails mid-render (after partial output), `http.Error` attempts a second `WriteHeader`, producing the superfluous warning.

### Tertiary: Same bug in HandleDebts

`HandleDebts` (ui.go:135) has the identical `"default"` fallback. It doesn't crash because its template doesn't range over Tasks, but it shows empty data — the beliefs are nil because the snapshot was skipped.

---

## Design Analysis

The `GetContext(taskID)` method is designed for the **RCP protocol** — a single-task context view for agents. The Insights page is a **human dashboard** that should show ALL tasks and ALL beliefs. These are different query patterns:

| Concern | GetContext (RCP) | Insights (Dashboard) |
|---------|------------------|----------------------|
| Scope | One task | All tasks |
| Beliefs | Filtered by task's project | All beliefs |
| Dependencies | Task-specific | N/A |
| Input | Required task_id | None |

Forcing a dashboard through `GetContext` with a fallback task ID is the architectural mismatch.

---

## Proposed Fix

### 1. Add `ListAll` to `work.Store`

New method returning all tasks without a project filter.

```go
// store.go
func (s *Store) ListAll(ctx context.Context) ([]*Task, error) {
    rows, err := s.db.QueryContext(ctx,
        `SELECT id, project_id, title, description, status, priority,
                current_agent, governance_ref, reopened_from_task_id,
                created_at, updated_at
         FROM conductor_task ORDER BY created_at DESC`)
    ...
}
```

### 2. Add `GetAllBeliefs` to `epistemic/view.go`

New function returning all beliefs across all scenarios. The existing `GetSnapshot` is scenario-scoped.

```go
// view.go
func GetAllBeliefs(ctx context.Context, db *sql.DB) ([]BeliefView, error) {
    rows, err := db.QueryContext(ctx,
        `SELECT id, claim, claim_type, status, debt::STRING, final_truth
         FROM belief ORDER BY claim`)
    ...
}
```

### 3. Add `GetDashboard` to App

Orchestrates dashboard data without requiring a task ID.

```go
// app.go
type Dashboard struct {
    Tasks   []*work.Task
    Beliefs []epistemic.BeliefView
}

func (a *App) GetDashboard(ctx context.Context) (*Dashboard, error) {
    tasks, err := a.workStore.ListAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("list tasks: %w", err)
    }
    beliefs, err := epistemic.GetAllBeliefs(ctx, a.db)
    if err != nil {
        return nil, fmt.Errorf("get beliefs: %w", err)
    }
    return &Dashboard{
        Tasks:   tasks,
        Beliefs: beliefs,
    }, nil
}
```

Returns `Dashboard` with empty slices (never nil) when no data exists.

### 4. Rewrite `HandleInsights`

Remove `task_id` → `GetContext` path. Use dashboard.

```go
func (s *Server) HandleInsights(w http.ResponseWriter, r *http.Request) {
    dash, err := s.app.GetDashboard(r.Context())
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to load dashboard: %v", err),
            http.StatusInternalServerError)
        return
    }
    data := map[string]interface{}{
        "Tasks":   dash.Tasks,
        "Beliefs": dash.Beliefs,
    }
    var buf bytes.Buffer
    if err := s.tmpl.ExecuteTemplate(&buf, "insights.html", data); err != nil {
        http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    buf.WriteTo(w)
}
```

### 5. Rewrite `HandleDebts`

Same pattern — remove `task_id` → `GetContext("default")`, use dashboard for beliefs.

### 6. Buffer template execution (both handlers)

Render to `bytes.Buffer` first. Only write to `ResponseWriter` if `Execute` succeeds. This prevents the superfluous `WriteHeader` bug.

### 7. Add regression test

New test in `ui_test.go`:

```go
func TestInsightsPage_NoTaskID(t *testing.T) {
    // seed project + task + belief (reuse existing test helpers)
    // GET /ui/insights  (NO task_id param)
    // assert HTTP 200
    // assert body contains seeded task title
    // assert body contains seeded belief claim
    // assert no template error
}
```

### 8. Update existing test

`TestInsightsPage` currently passes `?task_id=...`. Update to verify dashboard behavior (all tasks visible without passing task_id).

---

## Files Changed

| File | Action | Purpose |
|------|--------|---------|
| `internal/work/store.go` | EDIT | Add `ListAll` method |
| `internal/epistemic/view.go` | EDIT | Add `GetAllBeliefs` function |
| `internal/application/app.go` | EDIT | Add `Dashboard` type + `GetDashboard` method |
| `internal/ui/ui.go` | EDIT | Rewrite `HandleInsights` + `HandleDebts`, add buffer rendering |
| `internal/ui/ui_test.go` | EDIT | Add regression test, update existing test |

No schema changes. No migration changes. No changes to `GetContext` (still used by MCP/RCP). No changes to seed, domain pack, Solvent kernel, or research architecture.

---

## Verification

```bash
go test ./internal/ui/...
go test ./internal/work/...
go test ./internal/epistemic/...
go test ./...
go vet ./...
```

Then manually: `argus serve` → `http://localhost:8080/ui/insights` (no task_id) → should show seeded task and belief.

---

## What This Does NOT Change

- `GetContext` remains unchanged — still used by MCP `get_context` and RCP
- No new services, databases, or ports
- Template HTML structure unchanged (already expects `Tasks` and `Beliefs` as lists)
- No authority model changes
- No Trust UI auth changes

---

## Open Question

**Should the Debts page also become a full dashboard (showing all beliefs regardless of task), or should it remain task-scoped with a proper fallback when task_id is missing?**

The Debts page shows belief-level debt info, which is naturally a dashboard concern. The current template already ranges over `.Beliefs` as a flat list, not per-task. I'd lean toward making it a dashboard too — same pattern as Insights — but want to confirm intent before implementation.

---

## Risk Assessment

- **Low risk:** Changes are confined to the Trust UI presentation layer + two small store methods
- **No authority boundary changes:** Solvent, Coordinator, MCP unchanged
- **Backward compatible:** `GetContext` untouched; agents continue using RCP as before
- **Testable:** Existing tests provide the test infrastructure; new test covers the exact regression
