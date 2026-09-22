# ARGUS — Final Read-Path Repair + Freeze-Gate Verification

**Status:** Plan  
**Scope:** Fix confirmed read-path regressions, verify write→read loop, freeze the harness  
**Frozen constraints:** No new MCP tools, no new accounting infrastructure, no deferred feature implementation, no authority model changes

---

## 1. Fix Confirmed Read-Path Bugs

### 1A. `internal/epistemic/view.go:scanBelief` (P0)

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

### 1B. `internal/work/store.go:ListByProject` (P1)

SELECT returns 11 columns (missing `origin_packet_id`). Scan expects 12 destinations (includes `&task.OriginPacketID`).

**Fix:** Add `origin_packet_id` to the SELECT in the position expected by Scan (after `governance_ref`, before `reopened_from_task_id`):

```sql
SELECT id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
FROM conductor_task WHERE project_id = $1 ORDER BY created_at DESC
```

This matches the column order used by `GetByID`, `ListAll`, and `GetTasksByGovernanceRef`.

### 1C. `internal/epistemic/view.go:GetAllBeliefs` (P2)

SELECT omits `origin_packet_id` although BeliefView contains it and scanBelief now scans it.

**Fix:** Add `origin_packet_id` to the SELECT:

```sql
SELECT id, claim, claim_type, status, debt::STRING, final_truth, origin_packet_id
FROM belief ORDER BY claim
```

---

## 2. Rename "Agent" → "Claimant" in UI

**File:** `internal/ui/templates/insights.html:7`

The task column header says "Agent" but `current_agent` means task claimant, not packet provenance. This is a one-word semantic repair.

**Change:** `<th>Agent</th>` → `<th>Claimant</th>`

No new claimant/provenance subsystem. No graph visualization.

---

## 3. Disposition `TestAgentCannotPromote`

**File:** `internal/application/app_test.go:1274-1315`

**Current behavior:** The test calls `SubmitDecision` directly (application layer) with an unauthenticated command and expects rejection. The test fails because `SubmitDecision` does not enforce caller identity — it is an internal function used by authenticated human/control paths (Trust UI).

**Root cause:** The test asserts a boundary that does not exist at the application layer. The actual authority boundary is at the MCP tool surface: agents only see `get_context` and `submit_packet`. Agents never reach `SubmitDecision`.

**Decision:** Reclassify the test to verify the actual architectural invariant.

**Replace with:** `TestMCPDoesNotExposeAuthorityTools`

This test verifies that the MCP adapter tool list contains exactly:
- `argus.get_context`
- `argus.submit_packet`

And does NOT contain:
- `argus.promote`
- `argus.retract`
- `argus.discharge`

**Rationale:** Adding authentication checks to `SubmitDecision` merely to satisfy a misplaced test would be scope creep under the bookkeeping freeze. The authority boundary is enforced by tool-surface configuration, not by application-level auth checks.

The existing `TestMCPToolsCount` in `internal/mcp/adapter_test.go:17` already verifies this partially. The new test replaces `TestAgentCannotPromote` with a clear comment documenting the decision.

---

## 4. Add End-to-End Read-Back Test

**File:** `internal/application/app_test.go` (new test)

**Test: `TestGetContextReadsBackPersistedState`**

This is the freeze-critical invariant: `WRITE → COMMIT → READ-BACK`.

**Steps:**
1. Create a project in `conductor_project`
2. Create a task linked to that project
3. Build a valid packet with beliefs, evidence, edges
4. Persist via `app.Persist(ctx, tx, pkt)` + `tx.Commit()`
5. Call `app.GetContext(ctx, taskID)`
6. Assert:
   - `len(snapshot.Beliefs) == expected`
   - Each belief has correct `ID`, `Claim`, `Status`, `Debt`
   - Each belief has correct `OriginPacketID == packetID`
   - `len(snapshot.Evidence) == expected`
   - Evidence has correct `BeliefID`, `ProvenanceClass`, `ContentSHA256`
   - `len(snapshot.Edges) == expected`
   - Edges have correct `ParentID`, `ChildID`, `Kind`
   - `snapshot.Availability.Snapshot.Available == true`
   - `snapshot.Availability.Task.Available == true`

**Failure modes caught:**
- scanBelief column mismatch → GetSnapshot error → empty snapshot
- provenance disappearing during projection
- GetContext returning stale or missing data after commit

---

## 5. Add `ListByProject` Regression Test

**File:** `internal/work/store_test.go` (new file or addition to existing)

**Test: `TestListByProjectReturnsOriginPacketID`**

**Steps:**
1. Create a project
2. Create a task with `OriginPacketID` set
3. Call `store.ListByProject(ctx, projectID)`
4. Assert:
   - Query succeeds (no scan mismatch)
   - `task.OriginPacketID` is populated correctly
   - Other fields (`Title`, `Status`, `Priority`) remain correct

---

## 6. Rename nil-DB MCP Test

**File:** `internal/mcp/adapter_test.go:121`

**Current name:** `TestMCPValidationRejectsMalformed`  
**New name:** `TestMCPValidationRejectsMalformed_FastUnit`

Add comment: `// Fast validation-path unit test using nil DB. Not production integration coverage.`

This prevents false confidence without deleting a useful cheap test. The DB-backed MCP tests remain the production-path proof.

---

## 7. Full Verification

### Automated

```bash
go vet ./...
go build ./...
go test ./...
```

### Manual smoke

1. `./bin/argus serve`
2. Submit one valid packet through the production MCP path
3. Call `argus.get_context` — verify new beliefs are visible
4. Verify `/ui/insights` — verify "Claimant" column, provenance visible, tasks render
5. Verify `/ui/debts` — verify debt obligations render
6. Verify provenance is visible where currently supported

---

## 8. Files Changed Summary

| File | Change |
|------|--------|
| `internal/epistemic/view.go` | Fix `scanBelief` (add `OriginPacketID`), fix `GetAllBeliefs` SELECT |
| `internal/work/store.go` | Fix `ListByProject` SELECT (add `origin_packet_id`) |
| `internal/ui/templates/insights.html` | Rename "Agent" → "Claimant" |
| `internal/application/app_test.go` | Replace `TestAgentCannotPromote` with `TestMCPDoesNotExposeAuthorityTools`, add `TestGetContextReadsBackPersistedState` |
| `internal/work/store_test.go` | Add `TestListByProjectReturnsOriginPacketID` |
| `internal/mcp/adapter_test.go` | Rename nil-DB test, add scoping comment |

---

## 9. Read-Path Defects Fixed

| # | Severity | File | Defect | Fix |
|---|----------|------|--------|-----|
| 1 | P0 | `internal/epistemic/view.go:scanBelief` | Scans 6 fields, SELECT returns 7 | Add `&b.OriginPacketID` to Scan |
| 2 | P1 | `internal/work/store.go:ListByProject` | SELECT 11 cols, Scan 12 dests | Add `origin_packet_id` to SELECT |
| 3 | P2 | `internal/epistemic/view.go:GetAllBeliefs` | SELECT omits `origin_packet_id` | Add to SELECT |

---

## 10. Freeze Gate

### PASS / FAIL

> Can a fresh agent submit research and then reconstruct exactly that committed research state through get_context?

**Answer after implementation:** [to be filled after tests pass]

### Freeze checklist

- [ ] packet validation on production MCP path
- [ ] packet persistence atomic
- [ ] provenance persisted
- [ ] valid packet reaches Commit
- [ ] malformed packet rejected before persistence
- [ ] exact retry idempotent
- [ ] GetContext read-back works
- [ ] Insights renders current state
- [ ] known TestAgentCannotPromote failure dispositioned
- [ ] full test suite green, or every exception explicitly dispositioned
- [ ] no new deferred accounting feature introduced

### Freeze recommendation

    READY TO FREEZE
    or
    NOT READY TO FREEZE

[to be determined after verification]

---

## 11. What This Plan Does NOT Do

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

All remain deferred under the bookkeeping freeze.

Research methodology and EBP scientific obligations are NOT frozen.
