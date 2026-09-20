You are performing the FINAL VERIFICATION PASS for Phase 8 of the ARGUS Trust Verification POC.

Do NOT redesign anything.
Do NOT add new architecture.
Do NOT add new persistence.
Do NOT expand Phase 8 scope.

The implementation has already passed the adversarial fix pass:

- 50/50 fix tests passing
- race detector clean
- go vet clean

The remaining purpose of this pass is to establish whether Phase 8 is
actually ready for the end-to-end dry run.

==================================================
CURRENT STATUS
==================================================

The prior adversarial review identified F1–F10.

The implementation report states:

F1  unknown debt fail-closed                 DONE
F2  evidence verification                    DONE
F3  /decisions authentication                DONE
F4  HTTP EvidenceClass propagation           DONE
F5  attributed Discharge                     DONE
F6  pack resolution fail-closed              DONE
F7  retraction → Conductor cancellation      DONE
F8  Trust UI retirement rule rendering       DONE
F9  deterministic RCP activity ordering      DONE
F10 REOPEN derives lineage                   DONE

EBP:
- ADD_DEBT = DEFERRED
- "new evidence creates debt" = DEFERRED / NOT EXERCISED

The final verification pass must specifically validate the remaining
evidence gaps around F7, F8, and the complete Phase 8 workflow.

==================================================
PRIMARY QUESTION
==================================================

Can Phase 8 now proceed to the real dry run?

Answer only after inspecting the actual code paths and executing the
required tests.

==================================================
1. VERIFY F7 — RETRACTION → CONDUCTOR CANCELLATION
==================================================

Do not accept the existing test merely because a mock Conductor returns
an empty task list.

Inspect the actual implementation of:

- coordinator/human.go
- cancelLinkedTasks()
- Conductor client methods
- Conductor task lifecycle / cancel endpoint
- governance_ref parsing

Then create or execute a real targeted test with at least:

Task A:
  same scenario
  active

Task B:
  same scenario
  proposed

Task C:
  different scenario
  active

Task D:
  same scenario
  accepted

Retract the governing belief.

Required result:

A → cancelled
B → cancelled
C → unchanged
D → unchanged

Also verify:
- cancelled is actually terminal
- no unrelated task is cancelled
- failure to cancel a linked task is not silently swallowed
- Solvent retraction succeeds before task cancellation is attempted
- the operation is deterministic

Prefer an integration test against the actual Conductor HTTP API if the
existing test infrastructure makes that practical.

If only mock testing is possible, state that limitation explicitly.

==================================================
2. VERIFY F8 — ACTUAL RETIREMENT RULE RENDERING
==================================================

Inspect the real Trust UI implementation.

Do not accept "the debt item name is shown" as sufficient.

For at least:

needMap
needNullModel
needFaithfulnessReview

verify that the Debts UI obtains and renders the actual Pack data:

- debt item
- required evidence class
- retirement rule identifier/name

Expected examples:

needMap
Required evidence: reproducible_artifact
Rule: map_check

needNullModel
Required evidence: operator_asserted
Rule: scope_clarification

The displayed values MUST originate from Pack/Coordinator state,
not hardcoded template strings.

Verify the rule shown in the UI is the same rule actually enforced by
Coordinator.

Add or execute an HTTP/UI test proving this.

==================================================
3. VERIFY F2 THROUGH THE REAL HTTP PATH
==================================================

The previous defect was especially important because direct Coordinator
tests passed while the production HTTP path was broken.

Test:

Trust UI-shaped JSON
→ authenticated /decisions
→ Coordinator
→ evidence lookup
→ rule validation
→ attributed Solvent discharge

At minimum execute:

A. valid evidence → discharge succeeds
B. wrong evidence class → refused
C. no evidence → refused
D. evidence belongs to another belief → refused
E. evidence belongs to another scenario → refused
F. unknown debt → refused
G. unknown/unavailable Pack → refused

Do not accept direct method invocation as the sole proof.

==================================================
4. VERIFY ATTRIBUTED DISCHARGE
==================================================

Trace the real production path.

Confirm:

- human path calls Solvent /v1/discharge
- bare /v1/beliefs/{id}/debt/retire is NOT used
- DischargedBy = configured operator principal
- InstrumentRef references the verified evidence or human-attestation
  justification
- the Solvent audit/discharge record actually contains attribution

Verify that browser-supplied actor information cannot replace the
configured operator identity.

==================================================
5. VERIFY /decisions AUTHENTICATION
==================================================

Execute:

A. no Authorization header
   → 401

B. invalid Authorization header
   → 401

C. valid configured credential
   → accepted

D. request contains forged actor_id
   → ignored

E. operator identity is taken from server-side configuration

Also verify that the read-only RCP endpoint remains accessible according
to the intended POC boundary.

==================================================
6. VERIFY RCP DETERMINISM HONESTLY
==================================================

Inspect the current mergeActivity implementation.

Confirm:

- Conductor activity preserves Conductor's own order
- Solvent activity preserves Solvent's own order
- source grouping is deterministic
- no global cross-store timestamp sort is being used to imply
  happens-before
- repeated identical reads produce identical ordering

The contract is:

"RCP guarantees deterministic presentation, not cross-ledger
happens-before semantics."

Verify that the implementation and tests agree with that statement.

==================================================
7. VERIFY REOPEN LINEAGE
==================================================

Execute:

retracted belief
→ REOPEN
→ new belief
→ derives edge old → new

Confirm:
- original remains retracted
- new belief exists
- edge kind = derives
- edge parent = original
- edge child = successor
- RCP exposes the lineage

==================================================
8. VERIFY PACK RESOLUTION
==================================================

Confirm the browser does NOT determine the active Pack.

Verify:

scenario
→ authoritative server-side Pack resolution
→ PackRegistry
→ retirement vocabulary/rule

Test:
- valid Pack
- missing Pack
- invalid Pack
- registry unavailable

All invalid cases must fail closed.

==================================================
9. VERIFY UNKNOWN DEBT
==================================================

Execute the real production path with:

debt_item = "arbitraryDebt"

Required:

- request refused
- no Discharge call
- no RetireDebt call
- no mutation to belief debt
- refusal is auditable where the existing mechanism supports it

==================================================
10. VERIFY CAPABILITY BOUNDARY
==================================================

Confirm actual OpenCode MCP configuration contains exactly:

argus.get_context
argus.submit_packet

and does NOT contain:
- solvent-mcp
- Conductor write tools
- edge mutation tool
- direct database capability

Execute or inspect the negative tests.

Do not rely on prompt instructions.

==================================================
11. VERIFY FULL PHASE 8 DRY-RUN PATH
==================================================

This is the most important test.

Run the actual thin workflow, preferably with Gate G0:

--------------------------------------------------
BRANCH A — ADVERSARIAL / DEAD-END
--------------------------------------------------

1. Human creates/selects a Conductor task.
2. Task is created with the correct scenario governance_ref.
3. Fresh Work OpenCode calls argus.get_context.
4. Work agent performs bounded research.
5. Work agent submits EBP packet.
6. Coordinator persists:
   - beliefs
   - evidence
   - derives edges
   - tasks
7. Fresh Adversarial OpenCode process starts with NO conversation history.
8. It calls argus.get_context.
9. It sees the prior work, evidence, debt and activity.
10. It submits an adversarial packet.
11. Coordinator persists a real contradicts edge.
12. Insights shows the adversarial challenge.
13. Human retracts the governing belief.
14. Solvent RetractCascade succeeds.
15. Linked Conductor task(s) become terminal cancelled.
16. Insights derives Dead End structurally.

--------------------------------------------------
BRANCH B — HUMAN DISCHARGE / PROMOTION
--------------------------------------------------

Use a surviving belief or successor, not the retracted belief.

1. Human reviews open debt in Debts UI.
2. UI shows actual Pack retirement rule.
3. Human supplies/selects qualifying evidence.
4. Authenticated decision reaches Coordinator.
5. Coordinator verifies:
   - Pack membership
   - retirement rule
   - actual qualifying evidence
6. Coordinator calls attributed /v1/discharge.
7. Solvent records operator identity and InstrumentRef.
8. UI refreshes debt state.
9. Human requests promotion.
10. Solvent refuses if any debt remains.
11. After all required debt is actually discharged,
    promotion succeeds if all other gates pass.
12. Audit trail remains reconstructible.

Do NOT promote the retracted belief.

==================================================
12. VERIFY EBP COMPLIANCE CLAIMS
==================================================

Final matrix:

Ideas enter free
  → verify PASS

Promotion costs debt
  → verify PASS

Debt does not kill
  → verify PASS

Debt remains payable
  → verify PASS

Human discharges debt
  → verify PASS

New evidence creates new debt
  → remain DEFERRED / NOT EXERCISED

ADD_DEBT
  → remain DEFERRED / NOT EXERCISED

No final-truth promotion
  → verify PASS

Accounting never becomes the work
  → verify PASS

Do not upgrade deferred EBP items to PASS merely because no test fails.

==================================================
13. COMPLETE TEST EXECUTION
==================================================

Run:

- go test ./...
- go test -race ./...
- go vet ./...

Run tests for:
- Solvent
- Conductor
- Coordinator
- MCP
- Trust UI
- packet
- domain-pack
- verifier
- corpus

Where repositories have separate Go modules, run the appropriate command
per module.

Do not report only the 50 fix tests.

Report the total repository/module test results.

==================================================
14. SOURCE-LEVEL FINAL AUDIT
==================================================

After tests pass, manually inspect the final code for these invariants:

1. No unknown debt can reach Solvent discharge.
2. No caller-supplied evidence class can substitute for actual evidence.
3. No unauthenticated decision invocation is possible.
4. No browser actor can choose operator identity.
5. Human retirement always uses attributed Discharge.
6. Pack resolution cannot silently fall through.
7. Retraction actually reaches Conductor cancellation.
8. REOPEN creates derives lineage.
9. RCP activity presentation is deterministic and honest.
10. Agents cannot access authority mutation tools.

==================================================
OUTPUT
==================================================

Produce:

PHASE8_FINAL_VERIFICATION_REPORT.md

Structure:

# Executive Verdict

One of:

PASS
PASS WITH LIMITATIONS
BLOCKED

Then:

## 1. F7 Verification
Actual test and result.

## 2. F8 Verification
Actual UI/rule-rendering test and result.

## 3. Production HTTP Retirement Path
Results for valid/invalid cases.

## 4. Authentication / Attribution
Results.

## 5. RCP Verification
Results.

## 6. REOPEN / Retraction / Dead-End
Results.

## 7. Full Dry-Run Result
Branch A and Branch B.

## 8. EBP Compliance
PASS vs DEFERRED.

## 9. Complete Test Results
Exact counts, race, vet.

## 10. Remaining Limitations
Only real limitations. Do not invent future work.

## 11. Final Recommendation
State explicitly whether the implementation is ready for Phase 8
dry-run acceptance.

==================================================
STOP CONDITION
==================================================

Do NOT modify architecture during this pass.

If a test exposes a defect:
- make the smallest implementation fix necessary
- add the regression test
- rerun the affected and full suites
- document the fix

If the complete workflow passes and no blocking issue remains:

Declare:

PHASE 8 IMPLEMENTATION VERIFIED
READY FOR DRY RUN

Do not declare Phase 8 fully successful until the actual dry run itself
has been executed.