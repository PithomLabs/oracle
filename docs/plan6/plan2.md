# PHASE 8 ADVERSARIAL FIX — IMPLEMENTATION PLAN

## Codebase State Summary

| Component | Key Finding |
|---|---|
| `coordinator/http/types.go` | `SubmitDecisionRequest` missing `EvidenceClass` field |
| `coordinator/http/handler.go` | No auth middleware; drops `EvidenceClass` in mapping |
| `coordinator/http/server.go` | No `AuthMiddleware` wrapper |
| `coordinator/human.go` | Uses `RetireDebt` not `Discharge`; hardcoded pack; no edge on reopen; no conductor cancel on retract |
| `coordinator/validate.go` | `ValidateRetirementRule` returns nil (allow) for unknown debt |
| `coordinator/client.go` | `Discharge` method exists; `ListEvidence` is stub returning `[]`; no `ListEvidenceForBelief` method |
| `coordinator/mock.go` | `Discharge` has no injectable fn; `ListEvidence` returns empty |
| `coordinator/context.go` | `mergeActivity` concatenates without deterministic ordering |
| Solvent `GET /v1/beliefs/{id}/evidence` | **Exists** — returns `[]EvidenceResponse` with `evidence_id`, `provenance_class` |
| Solvent auth | API key `Authorization: Bearer <key>` → `AuthFromContext` → `PrincipalID` |
| Conductor `POST /v1/tasks/{id}/cancel` | **Exists** — transitions any non-terminal state to `cancelled` |
| Conductor task→scenario link | `governance_ref` JSON with `reference_id` = scenario_id; **no reverse lookup** |

## Implementation Order (10 phases)

---

### Phase 1: HTTP EvidenceClass Propagation (F4)

**Files:** `coordinator/http/types.go`, `coordinator/http/handler.go`

1. Add `EvidenceClass string` field to `SubmitDecisionRequest` in `types.go`
2. In `handleSubmitDecision`, map `req.EvidenceClass` → `decisionReq.EvidenceClass`

**Test:** `coordinator/http/handler_test.go` — new file
- `TesthandleSubmitDecision_PreservesEvidenceClass`: POST JSON with `evidence_class: "reproducible_artifact"`, verify `DecisionRequest.EvidenceClass` is `"reproducible_artifact"` (not `""`)

---

### Phase 2: Pack Fail-Closed + Unknown Debt Rejection (F1, F6)

**Files:** `coordinator/validate.go`, `coordinator/human.go`

1. Add `ValidateDebtItemMembership(debtItem string, pack domainpack.Pack) error` — checks a single debt item against `GetDebtVocabulary()`; returns error if unknown
2. In `handleRetireDebt`, **before** calling `ValidateRetirementRule`:
   - Resolve pack (must succeed, or refuse — **no fallback**)
   - Call `ValidateDebtItemMembership(req.DebtItem, pack)` — refuse if unknown
3. Remove `extractPackRef` hardcoding; require pack resolution from scenario metadata or fail closed. For POC, accept a `PackRef` field on `DecisionRequest` that the Trust UI must supply (validated against known packs).
4. `ValidateRetirementRule`: change `info == nil` case from `return nil` to `return nil` only if the debt item is in vocabulary but has no rule; if debt item is unknown, it should never reach this function (already refused above).

**Tests:**
- `TestRetireDebt_UnknownDebtItem` → expect **refused**, verify `RetireDebt` NOT called
- `TestRetireDebt_PackLookupFails` → expect refused
- `TestRetireDebt_PackRegistryNil` → expect refused

---

### Phase 3: Evidence Verification Against Persisted Evidence (F2, D3)

**Files:** `coordinator/client.go`, `coordinator/mock.go`, `coordinator/human.go`

1. Add `ListEvidenceForBelief(beliefID, scenarioID string) ([]map[string]interface{}, error)` to `SolventClientInterface`, `SolventClient`, and `MockSolventClient`
   - `SolventClient`: calls `GET /v1/beliefs/{id}/evidence?scenario_id=` (endpoint already exists in Solvent)
   - `MockSolventClient`: injectable `ListEvidenceForBeliefFn`

2. In `handleRetireDebt`, after rule validation:
   - If rule requires `reproducible_artifact` (or any non-`operator_asserted` class):
     - Call `c.solventClient.ListEvidenceForBelief(req.BeliefID, req.ScenarioID)`
     - Verify at least one evidence row with `provenance_class == requiredClass` exists
     - If none found → refuse
   - If rule requires `operator_asserted`:
     - The human operator's attestation IS the qualifying evidence (EBP v2.1 meaning)
     - No agent evidence row is needed
     - Proceed directly

3. Build `InstrumentRef` from verified evidence IDs (D3):
   - Collect evidence IDs that matched the required class
   - Format: `"evidence:<id1>,<id2>,..."` or `"operator_asserted:<belief_id>"` for operator-asserted
   - Pass to `Discharge` as `instrumentRef`

**Tests:**
- `TestRetireDebt_CorrectPersistedEvidence` → executed
- `TestRetireDebt_WrongEvidenceClass` → refused
- `TestRetireDebt_NoEvidence` → refused
- `TestRetireDebt_EvidenceFromWrongBelief` → refused
- `TestRetireDebt_EvidenceFromWrongScenario` → refused
- `TestRetireDebt_OperatorAssertedPath` → executed (no evidence row needed)
- `TestRetireDebt_InstrumentRefBindsToEvidence` → verify InstrumentRef format

---

### Phase 4: Attributed Discharge (F5, D5 — BLOCKING)

**Files:** `coordinator/human.go`

1. Replace `c.solventClient.RetireDebt(...)` with `c.solventClient.Discharge(...)`:
   ```go
   c.solventClient.Discharge(
       req.ScenarioID,
       req.BeliefID,
       req.DebtItem,           // obligation_key
       instrumentRef,          // from Phase 3 evidence binding
       c.operatorID,           // discharged_by = configured operator
   )
   ```

2. Verify `Discharge` result contains attribution (audit trail)

**Tests:**
- `TestRetireDebt_UsesDischargeNotRetireDebt` — mock verifies `Discharge` called, `RetireDebt` NOT called
- `TestRetireDebt_DischargeContainsOperatorIdentity` — verify `discharged_by` = operator
- `TestRetireDebt_DischargeContainsInstrumentRef` — verify instrument_ref != ""
- `TestRetireDebt_BrowserCannotOverrideOperator` — request body `actor_id` is ignored

---

### Phase 5: Authenticate `/decisions` (F3, D4)

**Files:** `coordinator/http/server.go`, `coordinator/http/handler.go`

1. Implement minimal API-key auth middleware matching Solvent's pattern:
   - Read `Authorization: Bearer <key>` header
   - Look up key in configured `map[string]string` (key → principal UUID)
   - If missing/invalid → HTTP 401
   - On success, inject authenticated principal into context
   - Server-side configured operator principal (`ARGUS_OPERATOR_PRINCIPAL_ID`) is used for decision handling, NOT the browser-supplied identity

2. `NewServer` takes an `authKeys map[string]string` parameter
3. `ServeHTTP` wraps handler with auth middleware for `/decisions` routes only
4. RCP `/context/{task_id}` remains unauthenticated (read-only)

**Tests:**
- `TestDecisions_UnauthenticatedRefused` → 401
- `TestDecisions_InvalidKeyRefused` → 401
- `TestDecisions_ValidKeyAllowed` → reaches handler
- `TestDecisions_BrowserActorCannotOverride` → configured operator used regardless of request body
- `TestContext_NoAuthRequired` → RCP still accessible

---

### Phase 6: Retraction → Conductor Task Cancellation (F7)

**Files:** `coordinator/human.go`, `coordinator/client.go`, `coordinator/mock.go`

**Problem:** Conductor has no `ListByScenario` endpoint. Tasks link to scenarios via `governance_ref` JSON.

**Solution (POC-minimal):** Since the Coordinator created the tasks and knows the project, and all tasks in a project share the same scenario (via `governance_ref`), the Coordinator can:
1. Use an existing `ListTasks(projectID)` call (needs projectID from task metadata)
2. Filter client-side by parsing `governance_ref` to find tasks matching the scenario
3. Cancel matching active/proposed tasks

**Implementation:**
1. Add `ListTasksByGovernanceRef(projectID, scenarioID string) ([]map[string]interface{}, error)` to `ConductorClientInterface`
   - Implementation: `ListTasks(projectID)` then filter by `governance_ref.reference_id == scenarioID`
2. In `handleRetract`, after `RetractBelief`:
   - Determine `projectID` and `scenarioID` from the request/task context
   - Call `ListTasksByGovernanceRef` to find linked tasks
   - For each non-terminal task: call `CancelTask(taskID)`

**Tests:**
- `TestRetract_CancelsLinkedTasks` — verify `CancelTask` called for matching tasks
- `TestRetract_DoesNotCancelUnlinkedTasks` — verify unrelated tasks untouched
- `TestRetract_DoesNotCancelTerminalTasks` — verify accepted/cancelled tasks untouched
- `TestRetract_TaskCancelledIsTerminal` — verify `cancelled` is terminal in Conductor

---

### Phase 7: REOPEN Derives Edge (F10)

**Files:** `coordinator/human.go`

1. In `handleReopen`, after `CreateBelief`:
   - Call `c.solventClient.CreateEdge(req.BeliefID, newBeliefID, "derives")`
   - The original retracted belief is the parent; the new belief is the child

**Tests:**
- `TestReopen_CreatesDerivesEdge` — verify `CreateEdge` called with `kind="derives"`
- `TestReopen_OriginalRetained` — original belief not deleted
- `TestReopen_LineageVisibleViaRCP` — edge appears in RCP context

---

### Phase 8: Trust UI Rule Rendering (F8)

**Files:** `trust-ui/templates/debts.html`, `trust-ui/server.go`, `coordinator/http/handler.go`, `coordinator/http/types.go`

1. Add `GET /v1/packs/{packId}/retirement-rules` endpoint to Coordinator HTTP (or embed rules in RCP context)
2. Alternative (simpler): Include retirement rules in the RCP context response under `epistemic.retirement_rules`
3. In `debts.html`, for each debt item:
   - Display the actual required evidence class and rule name from the Pack
   - Replace generic "Applicable retirement rules apply" with:
     ```
     needMap
     Required evidence: reproducible_artifact
     Rule: map_check
     ```

**Test:**
- `TestDebtsPage_ShowsActualRule` — HTTP GET /debts, verify response contains actual rule text

---

### Phase 9: Deterministic RCP Activity Ordering (F9, D1)

**Files:** `coordinator/context.go`

1. Replace `mergeActivity` with source-grouped deterministic ordering:
   - Group by source: Conductor entries first (in their own order), then Solvent entries (in their own order)
   - Do NOT sort across sources by timestamp
   - Each entry already has `source: "conductor"` or `source: "solvent"` tag
   - Preserve each source's authoritative ordering

2. Document: "RCP guarantees deterministic presentation, not cross-ledger happens-before semantics."

**Tests:**
- `TestMergeActivity_Deterministic` — same inputs → same output order
- `TestMergeActivity_PreservesSourceOrder` — Conductor events stay in Conductor order; Solvent in Solvent order
- `TestMergeActivity_NoCrossLedgerChronology` — verify adjacency doesn't imply causality

---

### Phase 10: EBP Compliance + Documentation (D2, D6, D7)

**Files:** `PHASE8_ADVERSARIAL_FIX_REPORT.md`, plan update

1. Mark "New evidence creates new debt" as **DEFERRED / NOT EXERCISED**
2. State explicitly: "ADD_DEBT is not exercised by the Phase 8 dry run"
3. Add to spec: "A debt item MUST belong to the active Domain Pack debt vocabulary before retirement. Absence from vocabulary is fail-closed."
4. No implementation of ADD_DEBT

---

## Test Matrix (Production-Path)

| # | Test | Defect | Path |
|---|---|---|---|
| 1 | HTTP EvidenceClass survives transport | F4 | Trust UI → HTTP → Coordinator |
| 2 | Unknown debt refused | F1 | HTTP → Coordinator → refusal |
| 3 | Pack lookup failure refused | F6 | HTTP → Coordinator → refusal |
| 4 | Correct evidence exists → executed | F2 | HTTP → Coordinator → Solvent evidence check → Discharge |
| 5 | Wrong evidence class → refused | F2 | HTTP → Coordinator → Solvent evidence check → refusal |
| 6 | No evidence → refused | F2 | HTTP → Coordinator → Solvent evidence check → refusal |
| 7 | Evidence from wrong belief → refused | F2,D3 | HTTP → Coordinator → Solvent evidence check → refusal |
| 8 | Operator asserted path → executed | F2 | HTTP → Coordinator → Discharge (no evidence row) |
| 9 | InstrumentRef binds to evidence | D3 | HTTP → Coordinator → Discharge with evidence IDs |
| 10 | Unauthenticated → refused | F3 | curl without auth → 401 |
| 11 | Invalid auth → refused | F3 | curl with bad key → 401 |
| 12 | Valid auth → allowed | F3 | curl with valid key → reaches handler |
| 13 | Browser actor cannot override | F3,D4 | HTTP with actor_id in body → configured operator used |
| 14 | Discharge called, RetireDebt NOT called | F5 | Mock verification |
| 15 | Discharge contains operator identity | F5 | Mock verification of discharged_by |
| 16 | Retraction cancels linked tasks | F7 | Mock verification of CancelTask calls |
| 17 | Reopen creates derives edge | F10 | Mock verification of CreateEdge call |
| 18 | Activity ordering deterministic | F9,D1 | Repeated reads → identical output |
| 19 | Activity preserves source order | D1 | No cross-ledger chronology |
| 20 | UI shows actual retirement rule | F8 | HTTP GET /debts → response contains rule text |

## Remaining Limitations

- **ADD_DEBT deferred**: Adversarial contradictions create `contradicts` edges but do not reinstate debt. This is a future human-gated extension.
- **Conductor reverse lookup**: POC uses client-side filtering of `ListTasks(projectID)` by `governance_ref`. Production would need a proper index and query endpoint.
- **Pack resolution**: For POC, `DecisionRequest.PackRef` is supplied by Trust UI and validated against registered packs. Production would derive from scenario metadata.
- **Network binding**: localhost-only deployment reduces exposure but does not replace the auth boundary (documented).

## Acceptance Criteria Checklist

| Criterion | Status |
|---|---|
| Unknown debt fails closed | Phase 2 |
| Evidence actually verified | Phase 3 |
| `/decisions` authenticated | Phase 5 |
| EvidenceClass survives HTTP | Phase 1 |
| Human discharge attributed | Phase 4 (BLOCKING) |
| Pack failure fails closed | Phase 2 |
| Retraction cancels linked tasks | Phase 6 |
| Reopen preserves derives lineage | Phase 7 |
| RCP ordering deterministic | Phase 9 |
| Agent capability surface unchanged | Verified (2 tools only) |
| ADD_DEBT deferred | Phase 10 |
