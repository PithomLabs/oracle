# ARGUS — Adversarial Code Review: Current Provenance / Agent Workflow Changes

**Reviewer:** Independent adversarial software architect  
**Repository:** `/home/chaschel/Documents/go/oracle`  
**Commit range inspected:** `05ca3ea` → `173c71b` (initial commit through `feat: expose belief edge graph through get_context`)  
**Mode:** Read-only. No code modified.

---

## 1. EXECUTIVE SUMMARY

The recent changes introduced explicit agent identity fields, `origin_packet_id` provenance columns, `edge_provenance`, and `packet_submission`. The **transactional atomicity of the production MCP path is sound**: `internal/mcp/adapter.go` wraps `app.Persist` in a single `BeginTx`/`Commit`/`defer Rollback`, and every write in `Persist` uses that same `tx`.

However, a **critical validation bypass** exists in the production MCP path: `packet/v1.Validate()` is never invoked. The production adapter (`internal/mcp/adapter.go`) calls only `app.Compile` + `app.Validate`, which perform minimal header/debt/evidence-class checks. The comprehensive structural validator — role validation, agent identity consistency, duplicate local_id detection, claim-type enforcement, edge-kind validation, reference format validation — exists in `packet/v1/validate.go` but is **dead code on the production path**. Compounding this, `packet/v1/validate.go:parsePackRef` uses `-` as separator while the actual wire format and all production code use `@`, making the dead validator additionally **broken for real packets**.

Provenance attribution for beliefs, evidence, and tasks is **correct and preserves origin** via `ON CONFLICT (id) DO NOTHING`. Edge provenance is also correct for same-edge retries. But **evidence and edge provenance are invisible to the UI**, and `GetDashboard` leaks all tasks across all scenarios.

---

## 2. FINDINGS (SEVERITY-RANKED)

### P0 — Authority / Integrity Failure

#### Finding 1: `packet/v1.Validate()` is completely bypassed on the production MCP path

- **File/Function:** `internal/mcp/adapter.go:handleSubmitPacket` (lines 57–87)
- **Concrete failure scenario:** An adversarial or buggy agent submits a packet with:
  - `role: "work"` but `agent.role: "adversarial"` (or any invalid role string)
  - Missing `agent.id`, `agent.harness`, or `agent.model`
  - Duplicate `local_id` values within beliefs, evidence, edges, or tasks
  - Empty `claim` on a belief
  - Invalid `claim_type` (e.g., `"derrrived"`)
  - Edge `kind: "supercedes"` (not `"derives"` or `"contradicts"`)
  - Evidence `provenance_class: "fabricated"`
  - `belief_ref: "local:nonexistent"`
  - `from_ref`/`to_ref` with malformed references
- **Why current tests do not catch it:** Tests for `internal/mcp/adapter_test.go` only verify tool count, tool routing, and missing-task-id handling. They never submit a structurally invalid packet through `HandleTool` and assert rejection. `packet/v1/validate_test.go` tests the validator in isolation, but no test verifies the production MCP adapter invokes it.
- **Recommended minimal fix:** In `internal/mcp/adapter.go:handleSubmitPacket`, after unmarshaling, call `packetv1.Validate(&pkt, a.app.PackRegistry())` before `app.Compile`/`app.Validate`. Expose the pack registry on `App` (add a getter or pass it to the adapter at construction).

#### Finding 2: `packet/v1/validate.go:parsePackRef` uses wrong separator

- **File/Function:** `packet/v1/validate.go:parsePackRef` (lines 95–101)
- **Concrete failure scenario:** Even if Finding 1 were fixed, `packet/v1.Validate` would reject every legitimate production packet because it splits `pack_ref` on `"-"` while the actual wire format and `domain-pack/registry.go:registryKey` use `"@"`. A packet with `pack_ref: "bmist@1.1.0"` would parse to `packID="bmist@1.1.0"`, `version=""`, and the registry lookup would fail with `"pack bmist@1.1.0@ not found"`.
- **Why current tests do not catch it:** `packet/v1/validate_test.go` uses `PackRef: "bmist-1.0.0"` (with `-`). The integration test and all production code use `@`. The test suite validates the wrong grammar.
- **Recommended minimal fix:** Align `packet/v1/validate.go:parsePackRef` with `app.go:parsePackRef` and `domain-pack/registry.go:registryKey`: split on `"@"`.

---

### P1 — Correctness / Provenance / Atomicity Failure

#### Finding 3: Invalid edge kinds are persisted without rejection

- **File/Function:** `internal/application/app.go:Persist` (lines 429–483)
- **Concrete failure scenario:** Because `packet/v1.Validate` is bypassed, an agent submits an edge with `kind: "invalid_kind"`. `app.Persist` does not check `edge.Kind` against the allowed vocabulary. The edge is inserted into `belief_edge` with the invalid kind. No error is returned.
- **Why current tests do not catch it:** `TestAdversarialContradictionVisible` and `TestAgentEdgesDoNotRetract` only test valid `"derives"` and `"contradicts"` kinds. No test submits an invalid kind through the production MCP path.
- **Recommended minimal fix:** Add `edge.Kind` validation in `app.Persist` before the insert:
  ```go
  if edge.Kind != packetv1.EdgeDerives && edge.Kind != packetv1.EdgeContradicts {
      return nil, fmt.Errorf("edge[%s]: invalid kind %q", edge.LocalID, edge.Kind)
  }
  ```

#### Finding 4: Duplicate `local_id` causes silent data corruption

- **File/Function:** `internal/application/app.go:Persist` (lines 352–367)
- **Concrete failure scenario:** A packet contains two beliefs with `local_id: "b1"` but different claims. `EntityID` computes different deterministic UUIDs. Both inserts succeed (different IDs). But `result.BeliefIDs["b1"]` is overwritten by the second belief. Subsequent evidence referencing `local:b1` resolves to the **second** belief's ID, not the first. The first belief is orphaned — persisted but unreferenced by the packet's own evidence.
- **Why current tests do not catch it:** No test submits duplicate local_ids through the production path.
- **Recommended minimal fix:** Add duplicate local_id detection in `app.Compile` or at the top of `app.Persist`:
  ```go
  seen := make(map[string]bool)
  for _, b := range pkt.Beliefs {
      if seen[b.LocalID] { return nil, fmt.Errorf("duplicate belief local_id: %s", b.LocalID) }
      seen[b.LocalID] = true
  }
  // Repeat for Evidence, Edges, Tasks.
  ```

#### Finding 5: Agent identity / role consistency is not enforced in production

- **File/Function:** `internal/application/app.go:Compile` (lines 203–220) and `internal/application/app.go:Validate` (lines 225–291)
- **Concrete failure scenario:** An agent submits a packet with `role: "work"` but `agent.role: "adversarial"`, or with empty `agent.id`. `app.Compile` and `app.Validate` do not check agent fields at all. The packet is persisted with inconsistent provenance. The `packet_submission` table records the raw agent fields, making the inconsistency permanent and visible in the Trust UI.
- **Why current tests do not catch it:** `TestAgentIdentityPersistedInSubmission` verifies that valid agent fields are persisted, but never tests invalid/missing fields or mismatched roles.
- **Recommended minimal fix:** Add agent identity validation to `app.Compile` (or `app.Validate`):
  ```go
  if pkt.Agent.ID == "" { return fmt.Errorf("agent.id is required") }
  if pkt.Agent.Role == "" { return fmt.Errorf("agent.role is required") }
  if pkt.Agent.Harness == "" { return fmt.Errorf("agent.harness is required") }
  if pkt.Agent.Model == "" { return fmt.Errorf("agent.model is required") }
  if pkt.Agent.Role != pkt.Role { return fmt.Errorf("agent.role %q must equal packet.role %q", pkt.Agent.Role, pkt.Role) }
  ```

#### Finding 6: Task `project_id` FK failure for scenarios without a project

- **File/Function:** `internal/application/app.go:Persist` (lines 486–498)
- **Concrete failure scenario:** An agent submits a packet for a new `scenario_id` (e.g., `"00000000-0000-0000-0000-000000000200"`). `Persist` creates tasks with `project_id = pkt.ScenarioID`. But `conductor_task.project_id` FK-references `conductor_project.id`. If no project with that UUID exists, the task insert fails with a FK violation, rolling back the entire packet — including beliefs and evidence that were already inserted.
- **Why current tests do not catch it:** All tests explicitly pre-create a `conductor_project` row. `cmdServe` seeds only one hardcoded project (`ProjectID = "00000000-0000-0000-0000-000000000100"`). No test verifies that a packet can be submitted for a scenario without a pre-existing project.
- **Recommended minimal fix:** Either:
  - Auto-create a project row when none exists for the scenario (with a default name), OR
  - Decouple `project_id` from `scenario_id` and require an explicit `project_id` in the packet, OR
  - Document and enforce that every scenario must have a corresponding project before packet submission.

#### Finding 7: Evidence and edge `origin_packet_id` are not projected to the UI

- **File/Function:** `internal/epistemic/view.go:EvidenceView` (lines 22–27), `internal/epistemic/view.go:EdgeView` (lines 36–41)
- **Concrete failure scenario:** A human opens the Trust UI Insights page. They can see task origin and belief origin, but cannot determine which packet created a specific piece of evidence or which packet introduced a specific edge. The `edge_provenance` table exists but is never queried by `GetSnapshot`. The `evidence.origin_packet_id` column exists but is not selected in `EvidenceView`.
- **Why current tests do not catch it:** `ui_test.go` only checks that the page renders and contains expected strings. It never asserts that evidence or edge provenance is displayed.
- **Recommended minimal fix:** Add `origin_packet_id` to `EvidenceView` and `EdgeView`, and join `edge_provenance` in `GetSnapshot` to populate `EdgeView.OriginPacketID`. Update `insights.html` to display these columns.

#### Finding 8: `GetDashboard` returns all tasks across all scenarios

- **File/Function:** `internal/application/app.go:GetDashboard` (lines 543–567)
- **Concrete failure scenario:** A human views the Insights page for scenario A. The Tasks table shows tasks from every project/scenario, not just scenario A. This conflates work from different research threads and makes it impossible to isolate a scenario's task state from the dashboard.
- **Why current tests do not catch it:** `integration_test.go:TestFullIntegration` only checks that the page renders, not that it is scenario-scoped.
- **Recommended minimal fix:** Filter tasks by `project_id = scenarioID` in `GetDashboard`, or add a `ListByProject` call to the work store and use it when `scenarioID != ""`.

---

### P2 — Significant Workflow or Observability Gap

#### Finding 9: Two divergent MCP adapters create maintenance confusion

- **File/Function:** `internal/mcp/adapter.go` (production) vs `mcp/adapter/adapter.go` (unused)
- **Concrete failure scenario:** A developer fixing a validation bug in the coordinator-based adapter (`mcp/adapter/adapter.go`) assumes the fix applies to production, but `cmdMCP` uses `internal/mcp/adapter.go` which has a completely separate path. The old adapter calls `coordinator.SubmitPacket` (which invokes `packet/v1.Validate`), while the production adapter calls `app.Compile`/`app.Validate` (which does not).
- **Why current tests do not catch it:** Both adapters have their own test packages. No test verifies that both adapters enforce the same validation rules.
- **Recommended minimal fix:** Delete `mcp/adapter/adapter.go` and its tests, or redirect `cmdMCP` to use the coordinator path if it is the intended canonical path. Ensure there is exactly one production MCP adapter.

#### Finding 10: `contentHash` excludes agent identity

- **File/Function:** `internal/application/app.go:contentHash` (lines 96–198)
- **Concrete failure scenario:** Two different agents submit the exact same belief/evidence/edge/task content. They get the same `content_hash`. The `submission_idempotency` table prevents the second packet's entities from being re-inserted (correct), but `packet_submission` records both packets with different `packet_id`s. The system cannot distinguish "same research, different author" from "same research, same author retrying".
- **Why current tests do not catch it:** `TestContentHash*` tests verify ordering and component differences, but none test that agent identity is included.
- **Recommended minimal fix:** Decide on semantics. If agent identity is part of packet identity, include `pkt.Agent.ID` + `pkt.Agent.Harness` + `pkt.Agent.Model` in the hash. If it is not, document this explicitly so agents understand that identical content from different agents is deduplicated.

#### Finding 11: `origin_packet_id` FK ordering is fragile

- **File/Function:** `internal/application/app.go:Persist` (lines 352–528)
- **Concrete failure scenario:** `belief.origin_packet_id`, `evidence.origin_packet_id`, `edge_provenance.origin_packet_id`, and `conductor_task.origin_packet_id` all FK-reference `packet_submission(packet_id)`. But `Persist` inserts the entities **before** inserting `packet_submission`. In standard SQL with immediate FK checking, this would fail on the first belief insert because `packet_submission` does not yet contain the packet ID.
- **Why current tests do not catch it:** The fault-injection test (`TestPartialFailureRetrySucceeds`) fails before reaching the belief insert. The success-path tests pass, which empirically suggests CockroachDB either defers FK checks or does not enforce them for this ALTER TABLE pattern, but this is not explicitly verified.
- **Recommended minimal fix:** Either:
  - Insert `packet_submission` first (before beliefs/evidence/edges/tasks), OR
  - Add the FK constraints as `DEFERRABLE INITIALLY DEFERRED`, OR
  - Explicitly verify CockroachDB's FK-checking behavior for ALTER TABLE–added constraints and document it.

#### Finding 12: `GetContext` assumes `Task.ProjectID == ScenarioID`

- **File/Function:** `internal/application/app.go:GetContext` (lines 572–614)
- **Concrete failure scenario:** `GetContext` uses `result.Task.ProjectID` as the scenario ID for `epistemic.GetSnapshot`. If a task is ever created with a `project_id` that does not equal its scenario's UUID, the snapshot returns beliefs/evidence/edges for the wrong scenario. This is currently true by convention but not enforced.
- **Why current tests do not catch it:** No test creates a task with `project_id != scenario_id` and calls `GetContext`.
- **Recommended minimal fix:** Add an explicit `scenario_id` column to `conductor_task` (or validate at insert time that `project_id == scenario_id`), and use that explicit field in `GetContext`.

---

### P3 — Documentation / Maintainability Issue

#### Finding 13: Tests give false confidence about validation coverage

- **File/Function:** `internal/application/app_test.go`, `internal/mcp/adapter_test.go`
- **Concrete failure scenario:** `TestPartialFailureRetrySucceeds` proves that a fault-injected transaction rolls back, but the fault is injected on the **second** `ExecContext`, meaning the first belief insert succeeds. The test never exercises the case where `packet_submission` insert fails after all entities succeed. `TestIdempotencyDuplicatePacket` proves entity-level idempotency but never submits a packet with structural errors.
- **Why current tests do not catch it:** No test submits a packet with an invalid role, missing agent fields, duplicate local_ids, or invalid edge kind through the production MCP adapter and asserts that it is rejected.
- **Recommended minimal fix:** Add adversarial tests that submit structurally invalid packets via `internal/mcp/adapter.HandleTool` and assert that `handleSubmitPacket` returns an error before any DB write occurs.

#### Finding 14: Old coordinator path is dead code with partial-persistence bugs

- **File/Function:** `coordinator/persist.go:SubmitPacket` (lines 13–69)
- **Concrete failure scenario:** The coordinator path calls `c.solventClient.CreateBelief`, `CreateEdge`, `CreateEvidence`, and `c.conductorClient.CreateTask` with **no transaction**. If `CreateEvidence` fails after `CreateBelief` and `CreateEdge` succeed, partial state is observable. This path is not used by `cmdMCP`, but it is still compiled and tested.
- **Why current tests do not catch it:** `coordinator/persist_test.go` uses mock clients that never fail mid-sequence.
- **Recommended minimal fix:** Either remove the coordinator persist path (if it is truly deprecated), or wrap it in a transaction. Add a test that mocks a mid-sequence failure and asserts no partial state.

---

## 3. ARCHITECTURE VERDICT

**FAIL — P1/P0 REMAINS**

The system cannot pass adversarial review in its current state. Two P0 findings (`packet/v1.Validate` bypass + broken pack-ref separator) mean that the production MCP path accepts structurally invalid packets that the documented validator would reject. P1 findings (invalid edge kinds, duplicate local_ids, agent identity inconsistency, task FK failure for new scenarios) further degrade correctness. Provenance attribution for entities is correct, but the UI cannot display evidence or edge provenance, and the dashboard leaks data across scenarios.

---

## 4. PROVENANCE VERDICT

**Partial pass with gaps.**

Can the system reliably trace:

```
agent
  ↓
packet
  ↓
belief/evidence/edge/task
  ↓
adversarial challenge
  ↓
human adjudication
```

**What works:**
- `packet_submission` records `agent_id`, `role`, `harness`, `model`, `content_sha256`, `scenario_id`, `task_id`.
- `belief.origin_packet_id`, `evidence.origin_packet_id`, `conductor_task.origin_packet_id` are populated on creation and preserved via `ON CONFLICT (id) DO NOTHING`.
- `edge_provenance` records the creating packet for each `(parent_id, child_id)` edge and is also preserved via `ON CONFLICT`.
- Retrying a packet does not overwrite earlier origins.

**What fails:**
- The Trust UI **cannot display** evidence or edge origins because `EvidenceView` and `EdgeView` omit `origin_packet_id` and `GetSnapshot` does not query `edge_provenance`.
- A human cannot determine from the UI which agent created a specific edge or piece of evidence.
- `GetDashboard` shows all tasks across all scenarios, so task provenance is mixed with unrelated work.

---

## 5. ATOMICITY VERDICT

**Pass for the production MCP path. Potential FK-ordering concern.**

`internal/mcp/adapter.go:handleSubmitPacket` uses a single `BeginTx`/`Commit`/`defer Rollback`. Every write in `app.Persist` uses the same `tx`. If any insert fails, the transaction rolls back and no caller-visible state changes.

**Caveat:** `origin_packet_id` columns FK-reference `packet_submission`, but entities are inserted before `packet_submission`. Current tests pass, suggesting CockroachDB either defers FK checks or does not enforce ALTER TABLE–added FKs for this pattern, but this should be explicitly verified and documented. If FK checks are immediate, the current code would fail on the first belief insert.

---

## 6. ADVERSARIAL-BOUNDARY VERDICT

**Fail — an adversarial packet can fabricate an apparent independent challenge.**

Because `packet/v1.Validate` is bypassed:
1. An adversarial packet can claim `role: "work"` while carrying `agent.role: "adversarial"`, making it appear as work in the UI (`packet_submission.role`) while actually being adversarial.
2. An adversarial packet can submit a `contradicts` edge targeting a **local** belief by using a canonical UUID that doesn't exist — `validateEdgeScenario` returns `"to-belief not found"`, which is a clear error. However, if the adversary uses a valid canonical belief in the same scenario, the challenge is correctly recorded. The boundary is partially enforced by `app.Persist`'s runtime checks (self-edge, cross-scenario, local-contradicts), but the **structural** validation that would catch malformed packets before they touch the DB is absent.

---

## 7. MCP VERDICT

**Fail — the agent-facing MCP path bypasses `packet/v1.Validate`.**

`internal/mcp/adapter.go:handleSubmitPacket` calls:
1. `json.Unmarshal` → raw struct
2. `app.Compile` → minimal header checks
3. `app.Validate` → debt vocab + evidence classes + verifier specs
4. `app.Persist` → DB writes

`packet/v1.Validate` (role validation, agent identity, duplicate local_ids, claim non-empty, claim_type non-empty, edge kind, reference format, task structure) is **never invoked**. A validator that exists but is bypassed is not a passing control.

---

## 8. DRY-RUN READINESS

**Not ready for a fresh Work Agent → Adversarial Agent → corrective Work Agent cycle.**

Reasons:
1. **Validation bypass (P0):** A buggy or malicious agent can submit structurally invalid packets that corrupt the ledger (duplicate local_ids, invalid edge kinds, mismatched roles).
2. **Task FK failure (P1):** A Work Agent submitting a packet for a scenario without a pre-existing project will see the entire packet roll back, including valid beliefs and evidence.
3. **UI observability gaps (P1/P2):** Humans cannot trace evidence or edge origins in the Trust UI, making adjudication of adversarial challenges incomplete.
4. **Dashboard data leakage (P2):** The Insights page mixes tasks from all scenarios, so a human reviewing an adversarial challenge cannot isolate the relevant work.

**Minimum bar for dry-run readiness:**
- Wire `packet/v1.Validate` into the production MCP path.
- Fix `parsePackRef` separator in `packet/v1/validate.go`.
- Add edge-kind and duplicate-local_id guards in `app.Persist`.
- Ensure task creation either auto-provisions a project or clearly documents the prerequisite.
- Expose evidence and edge `origin_packet_id` in the UI.
