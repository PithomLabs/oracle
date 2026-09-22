# ARGUS — Final Read-Path Repair + Freeze-Gate Verification

**Status:** Plan (revised after plan15_review.md)
**Scope:** Fix confirmed read-path regressions, verify write→read loop, freeze the harness
**Frozen constraints:** No new MCP tools, no new accounting infrastructure, no deferred feature implementation, no authority model changes, no scope broadening

---

## §1. Fix Confirmed Read-Path Bugs

### §1A. `internal/epistemic/view.go:scanBelief` (P0)

GetSnapshot SELECT returns 7 columns:

```
id, claim, claim_type, status, debt::STRING, final_truth, origin_packet_id
```

scanBelief scans only 6 destinations. This breaks GetContext for every scenario with beliefs.

**Fix:** Add `&b.OriginPacketID` to the Scan call:

```go
func scanBelief(s scanner) (*BeliefView, error) {
    var b BeliefView
    var debtRaw string
    if err := s.Scan(&b.ID, &b.Claim, &b.ClaimType, &b.Status, &debtRaw, &b.FinalTruth, &b.OriginPacketID); err != nil {
        return nil, err
    }
    b.Debt = parsePGArray(debtRaw)
    return &b, nil
}
```

### §1B. `internal/work/store.go:ListByProject` (P1)

**Verified column ordering** across all three task queries before changing:

| Method | SELECT columns | Scan destinations |
|--------|---------------|-------------------|
| `GetByID` (line 91) | 12 cols: `id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at` | 12 dests |
| `ListAll` (line 318) | 12 cols: same order as GetByID | 12 dests |
| `ListByProject` (line 342) | 11 cols: **missing `origin_packet_id`** | 12 dests ← MISMATCH |

**Fix:** Add `origin_packet_id` after `governance_ref`, before `reopened_from_task_id` (position 8, matching GetByID/ListAll):

```sql
SELECT id, project_id, title, description, status, priority, current_agent,
       governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
FROM conductor_task WHERE project_id = $1 ORDER BY created_at DESC
```

### §1C. `internal/epistemic/view.go:GetAllBeliefs` (P2)

SELECT omits `origin_packet_id` although BeliefView contains it and scanBelief now scans it.

**Fix:** Add `origin_packet_id` to the SELECT:

```sql
SELECT id, claim, claim_type, status, debt::STRING, final_truth, origin_packet_id
FROM belief ORDER BY claim
```

---

## §2. UI Semantic Repair

**Grep results for "Agent" in `internal/ui/templates/`:**

| File | Line | Context | Action |
|------|------|---------|--------|
| `insights.html` | 7 | `<th>Agent</th>` in **Tasks table** | Rename to `<th>Claimant</th>` |
| `insights.html` | 25 | `<th>Agent</th>` in **Recent Agent Submissions** | Keep — correct: submitting agent identity |
| `debts.html` | — | No "Agent" column | No change |

**Change:** Only `insights.html:7`: `Agent` → `Claimant` because `current_agent` means task claimant, not packet provenance.

No new UI systems. No graph visualization.

---

## §3. Disposition `TestAgentCannotPromote`

**Current test (app_test.go:1274-1315):** Calls `SubmitDecision` directly (application layer) with an unauthenticated command and expects rejection. Fails because `SubmitDecision` is an internal function used by authenticated human/control paths — it does not enforce caller identity.

**Root cause:** The test asserts a boundary that does not exist at the application layer. The actual authority boundary is at the MCP tool surface: agents only see `get_context` and `submit_packet`. Agents never reach `SubmitDecision`.

**Decision:** Replace with `TestMCPDoesNotExposeAuthorityTools` that verifies the actual architectural invariant.

**Also remove:** `TestAgentCannotDischargeDebt` (line 1317) and `TestAgentCannotRetractBelief` (line 1387) — same misplaced boundary. The existing `TestMCPToolsCount` and `TestAgentCannotRetract` in `adapter_test.go` already cover these invariants at the correct layer.

**Architectural decision documented:**
- Agent authority is bounded by MCP tool-surface exposure (tested)
- Human/control authority uses the existing operator/control path
- No authentication added to `SubmitDecision` merely to satisfy a misplaced test

---

## §4. Freeze-Critical End-to-End Test

**New test in `internal/application/app_test.go`:** `TestGetContextReadsBackPersistedState`

**Proves:** `WRITE → COMMIT → READ-BACK`

**Uses real DB.**

### Steps

1. Create project row in `conductor_project`
2. Create task via `work.Store.Create` linked to that project
3. **Packet P1** (work agent "agent-p1", packet_id "pkt-p1"):
   - 1 belief: claim="first claim from P1", debt=["needMap"]
   - 1 evidence: provenance_class="external_feed", content_sha256="hash-p1"
4. Persist P1 + commit
5. **Packet P2** (adversarial agent "agent-p2", packet_id "pkt-p2"):
   - 1 belief: claim="contradiction from P2"
   - 1 evidence: provenance_class="reproducible_artifact", content_sha256="hash-p2"
6. Persist P2 + commit
7. **Packet P3** (work agent "agent-p1", packet_id "pkt-p3"):
   - 1 edge: from P1's belief → P2's belief, kind="contradicts"
8. Persist P3 + commit
9. Call `app.GetContext(ctx, taskID)`

### Assertions

**Task identity:**
- `ctxResult.Task.ID == taskID`
- `ctxResult.Task.ProjectID == projectID`

**Availability:**
- `ctxResult.Availability.Snapshot.Available == true`
- `ctxResult.Availability.Task.Available == true`

**Beliefs (exact content, not just count):**
- `len(ctxResult.Snapshot.Beliefs) == 2`
- Belief from P1: `Claim == "first claim from P1"`, `OriginPacketID == "pkt-p1"`, `Debt == ["needMap"]`
- Belief from P2: `Claim == "contradiction from P2"`, `OriginPacketID == "pkt-p2"`
- `P1 != P2` — provenance is per-object, not copied

**Evidence (exact fields):**
- `len(ctxResult.Snapshot.Evidence) == 2`
- Evidence from P1: `BeliefID == belief1ID`, `ProvenanceClass == "external_feed"`, `ContentSHA256 == "hash-p1"`
- Evidence from P2: `BeliefID == belief2ID`, `ProvenanceClass == "reproducible_artifact"`, `ContentSHA256 == "hash-p2"`

**Edges (exact fields):**
- `len(ctxResult.Snapshot.Edges) == 1`
- Edge: `ParentID == belief1ID`, `ChildID == belief2ID`, `Kind == "contradicts"`

### Invariants proven

- `scanBelief` correctly scans all 7 fields (P0 regression caught)
- Provenance is per-object and survives read-back
- Two distinct origins are preserved correctly (not accidentally uniform)
- Edges reconstruct through GetContext
- GetSnapshot returns non-empty snapshot after commit

---

## §5. `ListByProject` Regression Test

**New test in `internal/work/store_test.go`:** `TestListByProjectReturnsOriginPacketID`

**Steps:**
1. Create project
2. Create task with `OriginPacketID` set to a known value
3. Call `store.ListByProject(ctx, projectID)`
4. Assert:
   - Query succeeds (no scan mismatch)
   - `task.OriginPacketID` is populated correctly
   - `task.Title`, `task.Status`, `task.Priority` remain correct

---

## §6. Fast MCP Test Clarification

**File:** `internal/mcp/adapter_test.go`

Rename `TestMCPValidationRejectsMalformed` → `TestMCPValidationRejectsMalformed_FastUnit`

Add comment:
```go
// Fast validation-path unit test using nil DB.
// Not production integration coverage.
// The DB-backed MCP tests (in app_test.go) are the production-path proof.
```

Keep the test. Do not delete a useful cheap test.

---

## §7. Authority Threat-Model Scope Note

Add to freeze record (short paragraph, not a new subsystem):

```
Agent authority is bounded by MCP tool-surface exposure (tested).
Human/control authority uses the existing operator/control path.

This POC freeze does not claim protection against an arbitrary network
actor that can directly reach an internal authority endpoint.

Revisit before non-loopback or multi-actor deployment.
```

Do not implement network hardening now unless an existing POC test proves it is required for the current localhost threat model.

---

## §8. Freeze Limitations Register

Add a short section to the freeze record. No new database table.

```
Deferred:
- packet-ID replay detection — trigger: concrete integrity incident from reuse
- context snapshot IDs — trigger: concurrent/replay dispute requiring exact reviewed state
- review-coverage persistence — trigger: human cannot distinguish reviewed-clean from never-reviewed
- epistemic-kind field — trigger: repeated cross-stratum misclassification
- agent attestation — trigger: declared identity cannot resolve a provenance dispute
- network authority hardening — trigger: deployment beyond loopback or multiple untrusted actors
```

These remain DEFERRED, not implemented.

---

## §9. Bookkeeping Freeze

The current freeze rule remains:

```
concrete failure → smallest fix
no concrete failure → defer
```

Do NOT add:
- context snapshots
- review-coverage persistence
- epistemic-kind fields
- agent attestation
- graph visualization
- corpus/vector search
- new MCP tools
- additional provenance infrastructure

Research methodology and EBP scientific obligations remain unfrozen.

---

## §10. Files Changed Summary

| File | Change |
|------|--------|
| `internal/epistemic/view.go` | Fix `scanBelief` (+1 field), fix `GetAllBeliefs` SELECT |
| `internal/work/store.go` | Fix `ListByProject` SELECT (+1 column) |
| `internal/ui/templates/insights.html` | Rename "Agent" → "Claimant" (line 7 only) |
| `internal/application/app_test.go` | Remove `TestAgentCannotPromote` + `TestAgentCannotDischargeDebt` + `TestAgentCannotRetractBelief`, add `TestGetContextReadsBackPersistedState` |
| `internal/mcp/adapter_test.go` | Rename nil-DB test, add scoping comment |
| `internal/work/store_test.go` | Add `TestListByProjectReturnsOriginPacketID` |

---

## §11. Read-Path Defects Fixed

| # | Severity | File | Defect | Fix |
|---|----------|------|--------|-----|
| 1 | P0 | `internal/epistemic/view.go:scanBelief` | Scans 6 fields, SELECT returns 7 | Add `&b.OriginPacketID` to Scan |
| 2 | P1 | `internal/work/store.go:ListByProject` | SELECT 11 cols, Scan 12 dests | Add `origin_packet_id` to SELECT |
| 3 | P2 | `internal/epistemic/view.go:GetAllBeliefs` | SELECT omits `origin_packet_id` | Add to SELECT |

---

## §12. Verification

### Automated

```bash
go vet ./...
go build ./...
go test ./...
```

### Manual smoke

1. `./bin/argus serve`
2. Submit one valid packet through production MCP
3. Call `argus.get_context` — verify newly persisted beliefs are returned exactly
4. Verify `/ui/insights` — "Claimant" column, provenance visible, tasks render
5. Verify `/ui/debts` — debt obligations render
6. Record actual results

---

## §13. Freeze Gate

### Freeze checklist

- [ ] packet validation works on production MCP path
- [ ] packet persistence is atomic
- [ ] provenance persists
- [ ] valid packet reaches Commit
- [ ] malformed packet rejected before persistence
- [ ] exact retry is idempotent
- [ ] GetContext reconstructs committed state exactly
- [ ] Insights renders current state
- [ ] TestAgentCannotPromote dispositioned/reclassified
- [ ] full suite green OR every remaining failure explicitly dispositioned
- [ ] limitations register exists
- [ ] no new accounting capability introduced
- [ ] manual smoke performed and results recorded

### Final question

> Can a fresh agent submit research and then reconstruct exactly that
> committed research state through get_context?

### Freeze recommendation

    READY TO FREEZE
    or
    NOT READY TO FREEZE

[to be determined after verification]

---

## §14. What This Plan Does NOT Do

- packet-id replay policy
- context snapshot IDs
- review coverage persistence
- epistemic-kind field
- agent attestation
- additional provenance UI
- graph visualization
- corpus/vector search
- new MCP tools
- additional authority layers
- network hardening

All remain deferred under the bookkeeping freeze.

Research methodology and EBP scientific obligations are NOT frozen.
