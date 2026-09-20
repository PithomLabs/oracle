# PHASE 8 ADVERSARIAL CODE REVIEW

## 1. Executive Verdict

**BLOCKED**

The implementation does not satisfy the frozen Phase 8 architecture. A malicious browser, an adversarial agent, or a malformed packet can cause ARGUS to record consequential epistemic transitions that the architecture reserves for humans under Solvent's attributed authority. The retirement-rule enforcement is non-functional in the production HTTP path, unknown debt items are silently accepted, and the human identity boundary is absent on the Coordinator decision endpoint.

---

## 2. Findings

### CRITICAL — F1: Unknown debt items silently execute

**Evidence:**
`coordinator/validate.go:130-145` — `ValidateRetirementRule` returns `nil` (allow) when `getRetirementRule` returns `nil` (no rule found for the debt item).
`coordinator/human.go:104-131` — `handleRetireDebt` does **not** call `ValidateDebtMembership`; it only checks the retirement rule.
`coordinator/retirement_test.go:225-245` — The test `TestRetireDebt_UnknownDebtItem` explicitly asserts that `unknownDebt` with `reproducible_artifact` produces `"executed"`.

**Failure:**
A client can submit an arbitrary debt string (e.g., `"arbitraryDebt"`). Because no retirement rule exists for it, `ValidateRetirementRule` returns `nil`, and the Coordinator proceeds to call `SolventClient.RetireDebt`, which executes `array_remove(debt, $2)` on the belief. The debt vocabulary is Pack-defined by frozen design; arbitrary strings must be rejected.

**Expected:**
Unknown debt items must be refused. Debt vocabulary lookup is mandatory before any retirement.

**Disposition:** fix now

---

### CRITICAL — F2: Evidence class is caller-supplied, not verified against persisted evidence

**Evidence:**
`coordinator/validate.go:130-145` — `ValidateRetirementRule` compares the submitted `evidence_class` string against the rule's expected class. It never queries Solvent to determine whether qualifying evidence of that class actually exists for the belief.
`trust-ui/templates/debts.html:40-43` — The UI provides a hardcoded dropdown of two classes.
`trust-ui/server.go:94-112` — The UI forwards the selected class as `evidence_class` in the JSON body.

**Failure:**
A malicious browser can send:
```json
{"debt_item":"needMap","evidence_class":"reproducible_artifact"}
```
The Coordinator trusts the string and retires the debt without verifying that a `reproducible_artifact` evidence row is actually persisted in Solvent for that belief. The browser can satisfy any rule without there being qualifying evidence.

**Expected:**
The Coordinator must query Solvent for actual evidence rows bound to the belief, confirm at least one exists with the required provenance class, and only then allow retirement.

**Disposition:** fix now

---

### CRITICAL — F3: Coordinator `/decisions` endpoint is unauthenticated

**Evidence:**
`coordinator/http/handler.go:22-37` — `ServeHTTP` routes `POST /decisions` to `handleSubmitDecision` with no authentication middleware.
`coordinator/http/server.go:14-24` — `NewServer` returns the raw handler; there is no `AuthMiddleware`.
`trust-ui/server.go:35` — Trust UI calls `coordinatorURL+"/decisions"` directly. There is no token, session, or HMAC in that call.

**Failure:**
Any network client that can reach the Coordinator can invoke `submitDecision` directly:
```bash
curl -X POST http://coordinator:9090/decisions \
  -H "Content-Type: application/json" \
  -d '{"type":"RETIRE_DEBT","belief_id":"...","scenario_id":"...","debt_item":"needMap"}'
```
Because the endpoint is unauthenticated, a browser script, an adversarial agent, or a `curl` command can trigger human-gated transitions without any human identity check.

**Expected:**
The `/decisions` endpoint must require an authenticated operator identity derived from server-side configuration (e.g., API key or mTLS), not from the request body.

**Disposition:** fix now

---

### CRITICAL — F4: HTTP handler drops `EvidenceClass`, breaking production retirement path

**Evidence:**
`coordinator/http/types.go:3-10` — `SubmitDecisionRequest` has no `EvidenceClass` field.
`coordinator/http/handler.go:39-67` — `handleSubmitDecision` decodes into `SubmitDecisionRequest` and constructs `coordinator.DecisionRequest` without `EvidenceClass`.
`trust-ui/server.go:106-112` — Trust UI includes `"evidence_class": req.EvidenceClass` in the body sent to `/decisions`.

**Failure:**
The `EvidenceClass` sent by the Trust UI is silently discarded by the HTTP handler. In production, `DecisionRequest.EvidenceClass` is always `""`. Consequently:
- For known debt items with rules, `ValidateRetirementRule` compares the rule's expected class against `""` and refuses. **All legitimate retirements via Trust UI are refused.**
- For unknown debt items, `ValidateRetirementRule` returns `nil` (allow), so **arbitrary debt strings are retired**.

The 149 passing tests all call `coordinator.SubmitDecision` directly, bypassing this handler. The production path is broken.

**Expected:**
`EvidenceClass` must be passed through the HTTP layer to the Coordinator decision logic.

**Disposition:** fix now

---

### HIGH — F5: Human path uses bare `RetireDebt`, not attributed `Discharge`

**Evidence:**
`coordinator/human.go:122-131` — `handleRetireDebt` calls `c.solventClient.RetireDebt(req.BeliefID, req.DebtItem, req.ScenarioID)`.
`coordinator/client.go:64-70` — `RetireDebt` maps to `POST /v1/beliefs/{id}/debt/retire`.
`solvent-main/api/discharge.go:10-79` — `POST /v1/discharge` exists, requires `instrument_ref`, `obligation_key`, and `discharged_by`, and inserts into `debt_discharge` with attribution.

**Failure:**
The frozen design mandates the human path use `POST /v1/discharge` so that `DischargedBy` and `InstrumentRef` are persisted in the `debt_discharge` audit table. The Coordinator instead calls the lower-level bare retire primitive, which produces no attributed discharge record.

**Expected:**
`handleRetireDebt` must call `c.solventClient.Discharge(...)` with `instrument_ref` and `discharged_by` (server-side principal), not `RetireDebt`.

**Disposition:** fix now

---

### HIGH — F6: Pack lookup failure silently bypasses validation

**Evidence:**
`coordinator/human.go:106-120` — Validation is gated on `c.packRegistry != nil`, `extractPackRef(req.ScenarioID)` returning non-empty `packID`, and `packRegistry.Get` succeeding. If any step fails, the code falls through to `RetireDebt` without validation.
`coordinator/human.go:134-139` — `extractPackRef` hardcodes `("bmist", "1.0.0")` for every scenario.

**Failure:**
If the Pack registry is uninitialized, the pack is not registered, or the hardcoded mapping is wrong, retirement proceeds with zero mechanical checks. The Pack-defined vocabulary and retirement rules are bypassed entirely.

**Expected:**
Pack resolution failure must be a hard error. The scenario-to-pack mapping must come from the scenario's actual metadata, not a hardcoded constant.

**Disposition:** fix now

---

### HIGH — F7: Retraction does not cascade to Conductor tasks

**Evidence:**
`coordinator/human.go:141-151` — `handleRetract` calls only `c.solventClient.RetractBelief(req.BeliefID, req.ScenarioID)`.
`solvent-main/kernel/kernel.go:186-209` — `RetractCascade` cancels live `action_intent` rows and updates `belief.status` to `retracted` within Solvent only.
`conductor/internal/store/task_repo.go:164-201` — `Transition` can cancel tasks, but nothing in the Coordinator calls it when a belief is retracted.

**Failure:**
The frozen design requires that a `contradicts` edge leading to retraction must make the linked Conductor task terminal-cancelled. Solvent retraction and Conductor task lifecycle are currently disconnected. A retracted belief leaves its Conductor tasks in whatever state they were in.

**Expected:**
`handleRetract` must also cancel linked Conductor tasks (or the architecture must define a binding between belief retraction and task terminal state).

**Disposition:** fix now

---

### MEDIUM — F8: Trust UI shows static text instead of the actual enforced rule

**Evidence:**
`trust-ui/templates/debts.html:36` — `<div class="debt-rule">Applicable retirement rules apply</div>` is hardcoded. No rule name, evidence class, or requirement is displayed.

**Failure:**
The UI tells the operator "applicable retirement rules apply" without showing which rule is enforced or what evidence class is required. This is a UI/backend semantics mismatch; the operator cannot adjudicate against the actual contract.

**Expected:**
The UI must render the actual retirement rule and required evidence class retrieved from the Pack for each debt item.

**Disposition:** fix now

---

### MEDIUM — F9: RCP activity merge is not deterministically ordered

**Evidence:**
`coordinator/context.go:126-140` — `mergeActivity` concatenates Conductor activity and Solvent activity without sorting by `created_at`.

**Failure:**
An adversarial agent or a partial failure can observe activity in an order that does not reflect actual temporal sequence. The frozen design requires deterministic ordering.

**Expected:**
Activity from both sources must be merged and sorted by `created_at` before inclusion in the RCP response.

**Disposition:** fix now

---

### MEDIUM — F10: REOPEN does not preserve `derives` lineage

**Evidence:**
`coordinator/human.go:153-164` — `handleReopen` calls `c.solventClient.CreateBelief(req.BeliefID, req.Reason, req.ScenarioID)`. It does not create an edge from the retracted belief to the new belief.

**Failure:**
The original retracted belief is not linked to the reopened belief via a `derives` edge. The lineage is broken in the graph.

**Expected:**
`handleReopen` must create a `derives` edge from the original belief to the new belief.

**Disposition:** fix now

---

### LOW — F11: Solvent bare retire endpoint accepts arbitrary debt strings

**Evidence:**
`solvent-main/kernel/sql.go:22-24` — `sqlRetireDebt` runs `array_remove(debt, $2::STRING)` with no vocabulary check.
`solvent-main/api/belief.go:108-170` — `handleRetireDebt` validates only that the debt item is non-empty; it does not check the Pack vocabulary.

**Failure:**
If an agent or client obtains Solvent API credentials, it can retire arbitrary debt strings via the lower-level primitive. This is a lower-level primitive issue; the higher-level fault is that the Coordinator calls this primitive instead of `Discharge`.

**Expected:**
The bare retire endpoint should remain internal or enforce vocabulary at the kernel layer. The Coordinator must not expose it as the human path.

**Disposition:** acceptable POC limitation if the endpoint is firewalled, but must be replaced by `Discharge` at the Coordinator layer.

---

## 3. Authority Boundary Audit

### Agent → MCP → Coordinator → Solvent

- **MCP Adapter** (`oracle/mcp/adapter/adapter.go:22-38`): Exposes exactly two tools: `argus.get_context` and `argus.submit_packet`.
  **Verdict:** PASS. Agents cannot issue Solvent mutations or retirement commands through MCP.

- **submit_packet** (`oracle/packet/v1/types.go`, `oracle/packet/v1/validate.go`): Packet schema contains beliefs, evidence, edges, and tasks. No authority transition, retirement, promotion, or discharge fields are present.
  **Verdict:** PASS. Packet submission cannot smuggle an authority transition.

- **Coordinator `SubmitPacket`** (`oracle/coordinator/persist.go:13-69`): Persists beliefs, edges, evidence, and tasks to Solvent/Conductor. No direct DB access.
  **Verdict:** PASS.

- **RCP** (`oracle/coordinator/context.go:36-89`): Read-only projection. Does not invent state. Reports `available: false` with reasons when backends are unavailable.
  **Verdict:** PASS with reservation (F9: ordering).

### Human → Trust UI → Coordinator → attributed discharge → Solvent

- **Trust UI** (`oracle/trust-ui/server.go:87-131`): Forwards `/api/retire` to Coordinator `/decisions`. No direct Solvent writes.
  **Verdict:** PASS.

- **Coordinator `/decisions`** (`oracle/coordinator/http/handler.go:39-67`): **No authentication.** Any client can invoke it.
  **Verdict:** FAIL. Human identity is not enforced at the boundary.

- **Coordinator `submitDecision`** (`oracle/coordinator/human.go:60-90`): Sets `ActorID` to server-side `c.operatorID`. Good. But `handleRetireDebt` calls `RetireDebt` instead of `Discharge`.
  **Verdict:** FAIL. Attributed discharge is not recorded.

- **Solvent `Discharge`** (`solvent-main/api/discharge.go:15-79`): Requires authenticated principal, cross-scenario guard, `discharged_by` matching principal, and inserts into `debt_discharge`.
  **Verdict:** PASS. This is the correct primitive, but it is not used by the Coordinator.

---

## 4. EBP v2.1 Compliance

| Doctrine | Status | Evidence |
|---|---|---|
| Ideas enter free | PASS | `kernel.go:35-55` — `EnterBelief` is ungated. |
| Promotion costs debt | PASS | `kernel.go:107-121` — `Promote` relies on DB CHECK `promoted_is_debt_free`. |
| Debt does not kill | PASS | Debt retirement exists; promotion is refused, not fatal. |
| Debt is forever payable | PASS | `array_remove` on retire; debt can be re-added. |
| New evidence creates new debt | PASS | Evidence addition does not mutate belief debt. |
| No final-truth promotion | PASS | `api/reads.go:473-477` — `final_truth=true` blocks promotion. |
| Human debt discharge | **FAIL** | `human.go:122-131` uses `RetireDebt`; no `Discharge` call, no `instrument_ref`, no `debt_discharge` row. |
| Accounting does not become work | PASS | Coordinator performs mechanical validation only. |

---

## 5. Test Adequacy

| Invariant | Test Exists? | Tests Real Path? | Verdict |
|---|---|---|---|
| Unknown debt rejected | Yes (`TestRetireDebt_UnknownDebtItem`) | No — asserts **executed** | **Wrong expectation** |
| Evidence class matched | Yes (8 rule tests) | No — calls `SubmitDecision` directly, bypassing HTTP handler where `EvidenceClass` is dropped | **Mocks pass, production fails** |
| Actual qualifying evidence exists | No | — | **Missing** |
| Cross-scenario isolation | No | — | **Missing** |
| Operator identity enforced | No (HTTP path) | — | **Missing** |
| Pack resolution failure | No | — | **Missing** |
| Duplicate discharge / replay | No | — | **Missing** |
| Retraction cancels Conductor tasks | No | — | **Missing** |

All retirement tests use `NewWithMockClientsAndPack` and call `SubmitDecision` directly. They prove the validator works in isolation. They do **not** prove the Trust UI → HTTP → Coordinator production path enforces the rule. The mock-based evidence is necessary but insufficient.

---

## 6. Required Fixes

1. **Wire `EvidenceClass` through the HTTP layer.** Add `EvidenceClass` to `coordinator/http/types.go` `SubmitDecisionRequest` and pass it to `coordinator.DecisionRequest` in `handleSubmitDecision`.

2. **Reject unknown debt items.** In `handleRetireDebt`, call `ValidateDebtMembership` before `ValidateRetirementRule`. Refuse if the debt item is not in the Pack vocabulary.

3. **Verify actual evidence existence.** In `handleRetireDebt`, after confirming the rule's evidence class matches, query Solvent (`ListEvidenceForBelief` or equivalent) to confirm at least one evidence row with that class exists for the belief.

4. **Use `Discharge` instead of `RetireDebt`.** Replace the `RetireDebt` call in `handleRetireDebt` with `Discharge`, passing `instrument_ref` (from request or generated) and `discharged_by` (from `c.operatorID`).

5. **Authenticate `/decisions`.** Wrap the Coordinator HTTP handler with an authentication middleware that validates the operator identity against server-side configuration.

6. **Make pack lookup failures hard errors.** If `extractPackRef` cannot resolve a pack or `packRegistry.Get` fails, refuse the retirement. Do not fall through to Solvent.

7. **Link retraction to Conductor task cancellation.** After `RetractCascade`, identify linked Conductor tasks (via `governance_ref` or belief→task mapping) and transition them to `cancelled`.

8. **Create `derives` edge on reopen.** `handleReopen` must call `CreateEdge` with `kind="derives"` from the original belief to the new belief.

9. **Show actual rule in Trust UI.** Render the Pack's retirement rule and required evidence class for each debt item instead of the hardcoded sentence.

10. **Sort RCP activity by timestamp.** In `mergeActivity`, sort the merged slice by `created_at` before returning.

---

## 7. Final Verdict

**Phase 8 may not proceed to the actual dry run.**

The retirement-rule enforcement is wired only in the unit-test path. The production HTTP path drops the evidence class, rejects all legitimate retirements, and allows arbitrary debt strings. The human identity is not enforced at the Coordinator boundary. The human path records a bare retire instead of an attributed discharge. These defects violate the frozen authority boundary: a fresh adversarial agent or a malicious browser can cause Solvent to retire debt without a qualifying evidence object, without a Pack-defined obligation, and without a discharged-by attribution.

**Block until fixes F1–F5 are implemented and the retirement test suite is expanded to cover the full HTTP production path with evidence-existence and identity checks.**
