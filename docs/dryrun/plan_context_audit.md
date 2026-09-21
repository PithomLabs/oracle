# Plan: Expose belief edge graph through argus.get_context

Date: 2026-09-21
Status: Plan-only, no code changes yet

## 1. Exact current belief_edge schema

```sql
CREATE TABLE IF NOT EXISTS belief_edge (
    parent_id  UUID NOT NULL REFERENCES belief(id),
    child_id   UUID NOT NULL REFERENCES belief(id),
    kind       STRING NOT NULL CHECK (kind IN ('derives','contradicts')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (parent_id, child_id)
)
```

Index: `belief_edge_child ON belief_edge (child_id)`
No `scenario_id` column — scenario scoping is via JOIN to `belief`.

## 2. Exact existing edge kinds

Two kinds, enforced by CHECK constraint:
- `derives` — child proposition follows from parent
- `contradicts` — child challenges/contradicts parent

Defined in `packet/v1/types.go`:
```go
const (
    EdgeDerives     = "derives"
    EdgeContradicts = "contradicts"
)
```

## 3. Exact write path used by submit_packet

`internal/application/app.go:430-462`:

```go
for _, edge := range pkt.Edges {
    fromID := resolveRef(edge.FromRef, result.BeliefIDs)
    toID := resolveRef(edge.ToRef, result.BeliefIDs)
    // ... validation ...
    tx.ExecContext(ctx,
        `INSERT INTO belief_edge (parent_id, child_id, kind)
         VALUES ($1::UUID, $2::UUID, $3)
         ON CONFLICT (parent_id, child_id) DO NOTHING`,
        fromID, toID, edge.Kind)
}
```

Edges are persisted with deterministic IDs resolved from packet local refs.
`ON CONFLICT DO NOTHING` — idempotent, no error on duplicate.
`result.EdgeCount++` tracks how many were inserted.

## 4. Can edge data be added to Snapshot without changing MCP tool surface?

**YES.** The MCP adapter (`internal/mcp/adapter.go:53`) returns `a.app.GetContext(ctx, req.TaskID)` directly. The response shape is determined by `application.Context` → `epistemic.Snapshot`. Adding an `Edges` field to `Snapshot` automatically includes it in the JSON response. No MCP adapter changes needed. No tool schema changes needed. The agent sees edges in the existing `get_context` response.

## 5. Exact response shape

Add to `internal/epistemic/view.go`:

```go
type EdgeView struct {
    ParentID string `json:"parent_id"`
    ChildID  string `json:"child_id"`
    Kind     string `json:"kind"`
}
```

Add to `Snapshot` struct:

```go
type Snapshot struct {
    Beliefs                []BeliefView   `json:"beliefs"`
    Evidence               []EvidenceView `json:"evidence,omitempty"`
    Edges                  []EdgeView     `json:"edges,omitempty"`     // NEW
    Intents                []IntentView   `json:"intents"`
    AuditLiveOnNonPromoted int            `json:"audit_live_on_nonpromoted"`
}
```

JSON returned to agent:
```json
{
  "task": {...},
  "dependencies": [...],
  "snapshot": {
    "beliefs": [...],
    "evidence": [...],
    "edges": [
      {"parent_id": "uuid", "child_id": "uuid", "kind": "derives"},
      {"parent_id": "uuid", "child_id": "uuid", "kind": "contradicts"}
    ],
    "intents": [...]
  }
}
```

## 6. Small projection-only or architectural?

**Small projection-only change.** Three files, ~15 lines of code:

1. `internal/epistemic/view.go` — add `EdgeView` struct, add `Edges` field to `Snapshot`, add query in `GetSnapshot`
2. No changes to `app.go`, `adapter.go`, `types.go`, `validate.go`, or any other file
3. No new MCP tools, no new subsystems, no schema migrations

The query joins `belief_edge` through `belief` for scenario scoping:
```sql
SELECT be.parent_id, be.child_id, be.kind
FROM belief_edge be
JOIN belief b ON b.id = be.parent_id
WHERE b.scenario_id = $1::UUID
ORDER BY be.parent_id, be.child_id
```

## 7. Tests needed

### In `internal/epistemic/view_test.go` (new file or extend existing):

1. **Derives edges returned** — persist belief B1, B2, edge B1→B2 kind=derives, call GetSnapshot, verify edge present with correct parent_id, child_id, kind
2. **Contradicts edges returned** — persist belief B1, B3, edge B3→B1 kind=contradicts, call GetSnapshot, verify edge present
3. **Edges reference correct immutable belief IDs** — verify parent_id and child_id match the persisted belief UUIDs (deterministic IDs from EntityID)
4. **No authority capability added** — verify Snapshot has no Promote/Discharge/Retract methods, edges are read-only projection
5. **Scenario scoping** — persist edges in scenario A, query scenario B, verify no edges returned
6. **Empty edges** — scenario with beliefs but no edges returns `edges: []` (not null)

### In `internal/application/app_test.go`:

7. **End-to-end via Persist** — submit packet with edges, call GetContext, verify edges in response

## Summary

| Question | Answer |
|---|---|
| belief_edge schema | parent_id, child_id, kind, created_at; PK(parent_id, child_id) |
| Edge kinds | `derives`, `contradicts` |
| Write path | INSERT with ON CONFLICT DO NOTHING in Persist() |
| MCP surface change needed? | NO — add to Snapshot struct, JSON auto-includes |
| Response shape | `EdgeView{ParentID, ChildID, Kind}` in snapshot.edges |
| Change size | Small projection-only (~15 lines across 3 files) |
| Tests | 7 tests covering derives, contradicts, ID correctness, authority isolation, scoping, empty state, end-to-end |
