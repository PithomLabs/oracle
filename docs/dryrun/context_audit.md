# ARGUS Audit — Adversarial Context Coverage

Date: 2026-09-21
Status: Audit-only, no code changes

## Complete Production Path

```
MCP argus.get_context (adapter.go:42)
    → application.GetContext (app.go:551)
        → workStore.GetByID (store.go:80)          -- task lookup
        → listDependencies (app.go:650)             -- dependency lookup
        → epistemic.GetSnapshot (view.go:65)        -- belief/evidence/intent query
            → SQL: belief WHERE scenario_id=$1      -- ALL beliefs, no filter
            → SQL: evidence WHERE scenario_id=$1    -- ALL evidence
            → SQL: action_intent WHERE scenario_id=$1 -- ALL intents
    → returns Context{Task, Dependencies, Snapshot, Availability}
```

## A. Belief Coverage

Given scenario S with beliefs B1, B2, ..., Bn and task T:

- **ALL scenario beliefs are returned** — no task, project, edge, or status filtering
- `entered`, `promoted`, and `retracted` beliefs are ALL returned
- No filtering by `claim_type`
- No edge traversal — flat query on `belief` table
- Scenario resolution uses `task.ProjectID` as the filter

## B. Exact SQL/Filtering

```sql
-- Beliefs (ALL, no status filter)
SELECT id, claim, claim_type, status, debt::STRING, final_truth
FROM belief WHERE scenario_id=$1::UUID ORDER BY claim

-- Evidence (ALL)
SELECT belief_id, source_url, provenance_class, content_sha256
FROM evidence WHERE scenario_id=$1::UUID ORDER BY belief_id, ingested_at

-- Intents (ALL)
SELECT belief_id, action, state
FROM action_intent WHERE scenario_id=$1::UUID ORDER BY belief_id
```

Inclusion: `scenario_id = task.ProjectID`. That's it. No other filters.

## C. Adversarial Visibility

**YES** — a fresh adversarial agent using ONLY `argus.get_context(task_id)` can see every belief in the current research scenario.

## D. Dry-Run Example

| Metric | Value |
|---|---|
| Beliefs seeded | 1 (G0) |
| Beliefs after Work Agent submits N beliefs | N+1 |
| Beliefs returned by get_context | N+1 (all) |
| Visible to adversarial agent | YES |

Seed belief ID: `00000000-0000-0000-0000-000000000102`

## E. Coverage Gap

**CRITICAL GAP: Edges are NOT returned.**

The `belief_edge` table is written during `submit_packet` (`app.go:437`) but **NOT read** in `GetSnapshot` (`view.go:65-140`). The `Snapshot` struct has no `Edges` field.

The adversarial agent can see:
- All beliefs ✓
- All evidence ✓
- All intents ✓
- **Edge graph (derives/contradicts) ✗**

This means the adversarial agent cannot see which beliefs derive from or contradict which other beliefs through the graph. It must infer relationships from claim text alone.

## F. Recommendation

**PARTIAL:**

`get_context` provides the full belief surface (all beliefs, all evidence, all intents) but does NOT provide the edge graph. The adversarial workflow must operate without explicit derivation/contradiction links.

**Impact on dry run:** The adversarial agent can still challenge any belief by ID, but cannot trace derivation chains or identify which beliefs are linked by `derives`/`contradicts` edges. This is workable but limits the adversarial agent's ability to find structural weaknesses in the argument graph.

**No code change proposed** — this is an audit-only finding.
