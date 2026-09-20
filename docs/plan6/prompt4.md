You are the implementation agent fixing the Phase 8 ARGUS Trust Verification POC
according to the attached PHASE 8 ADVERSARIAL CODE REVIEW.

Do NOT redesign the architecture.
Do NOT add a new database.
Do NOT weaken the authority boundaries.
Do NOT mark Phase 8 complete until the production-path defects below are fixed
and the required tests pass.

The adversarial review verdict is currently BLOCKED.

Your job is to implement the required fixes, add/repair tests, run the complete
suite, and produce a concise implementation report.

==================================================
FROZEN ARCHITECTURE
==================================================

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

OpenCode
  = work / agency

Conductor
  = operational task/workflow state

Solvent
  = epistemic authority:
    beliefs, evidence, debt, edges,
    promotion, retraction, authorization, audit

Trust UI
  = human observation + adjudication

EBP v2.1
  = epistemic operating law

RCP
  = thin read-only context projection

MCP
  = agent-facing transport

HTTP/API
  = canonical integration contract

Agents must NOT directly:
- retire debt
- promote
- authorize
- mutate Solvent
- mutate Conductor

Human consequential decisions must flow:

Trust UI
  → authenticated Coordinator
  → Solvent authoritative transition

==================================================
PRIMARY OBJECTIVE
==================================================

Close the exact defects identified by the adversarial review.

The review found that:

1. Unknown debt items can currently execute.
2. Evidence class is caller-supplied and not verified against persisted
   evidence.
3. Coordinator POST /decisions is unauthenticated.
4. HTTP drops EvidenceClass before production decision handling.
5. Human path uses bare RetireDebt instead of attributed Discharge.
6. Pack lookup failure bypasses validation.
7. Retraction does not cancel linked Conductor tasks.
8. UI does not show the actual enforced retirement rule.
9. RCP activity ordering is nondeterministic.
10. REOPEN does not create derives lineage.

The authoritative review is in the attached adv_review(7).md.
Use it as the defect baseline.

==================================================
STEP 0 — READ THE ACTUAL CODE FIRST
==================================================

Before editing, inspect:

oracle/coordinator/
oracle/domain-pack/
oracle/packet/
oracle/trust-ui/
solvent-main/
conductor/

Specifically inspect:

- coordinator/http/handler.go
- coordinator/http/types.go
- coordinator/human.go
- coordinator/validate.go
- coordinator/client.go
- coordinator/context.go
- coordinator/persist.go
- coordinator/coordinator.go
- trust-ui/server.go
- trust-ui/templates/debts.html
- Solvent discharge implementation
- Solvent debt-retire implementation
- Solvent belief/evidence APIs
- Solvent belief_edge APIs
- Conductor task lifecycle/store/API
- PackRegistry and BM-IST pack
- existing tests and mocks

Do not blindly follow the review's file names if the current repository differs.
Verify actual interfaces before editing.

==================================================
F1 — UNKNOWN DEBT MUST NEVER EXECUTE
==================================================

Current defect:

ValidateRetirementRule may return allow when no rule exists, and
handleRetireDebt does not enforce Domain Pack debt vocabulary first.

Fix:

1. Resolve the applicable Domain Pack.
2. Validate debt item membership in the Pack vocabulary.
3. If debt item is unknown:
   - return refusal/error
   - do NOT call Solvent
4. Only known debt items may proceed to retirement-rule validation.

Required invariant:

unknown debt
→ refusal
→ zero Solvent mutation

Do NOT interpret:
"no retirement rule"
as:
"retirement is allowed."

Fail closed.

Update the existing incorrect test:

TestRetireDebt_UnknownDebtItem

Expected:
refused / error
and verify Solvent retirement was NOT called.

==================================================
F2 — EVIDENCE CLASS MUST BE VERIFIED, NOT TRUSTED
==================================================

Current defect:

The browser sends:

{
  "debt_item": "...",
  "evidence_class": "reproducible_artifact"
}

and Coordinator compares the string against the Pack rule without
verifying that qualifying evidence actually exists.

This is an authority-boundary defect.

Implement actual evidence validation.

For retirement rules requiring:

reproducible_artifact

the Coordinator must verify that:
- evidence exists
- evidence belongs to the target belief
- evidence belongs to the correct scenario
- evidence.provenance_class matches the Pack-required class

Do NOT accept a caller-declared evidence class as sufficient.

For:
operator_asserted

respect the EBP meaning:
the human operator's explicit attestation is the qualifying human
assertion. Do not fabricate an agent-produced evidence record merely to
satisfy the type. Use the existing Solvent discharge semantics if they
already represent this attribution.

The important invariant is:

caller says evidence_class=X
≠
evidence actually satisfies retirement rule

The actual persisted evidence / operator attestation must satisfy the Pack.

Add tests for:
- correct persisted evidence class
- wrong evidence class
- no evidence
- evidence belonging to another belief
- evidence belonging to another scenario
- malformed/missing evidence class
- operator_asserted human attestation
- agent output is NOT automatically operator_asserted

==================================================
F3 — AUTHENTICATE COORDINATOR /decisions
==================================================

Current defect:

POST /decisions has no authentication boundary.

Fix with the smallest POC mechanism compatible with the existing
architecture.

Requirements:

- authentication must be server-to-server / server-configured
- identity must NOT come from the browser request body
- authenticated identity resolves to the configured operator principal
- missing/invalid credentials fail closed
- browser-supplied actor cannot override authenticated identity

Prefer an existing repository auth mechanism if one exists.

If none exists, implement the smallest explicit POC mechanism, such as:
- configured API key/header between Trust UI and Coordinator

Do NOT build a full user identity system.

Use the existing:
ARGUS_OPERATOR_PRINCIPAL_ID

for the authoritative operator identity after authentication.

Expected flow:

Browser
→ Trust UI server
→ authenticated Coordinator request
→ configured operator principal
→ decision handler

Add tests for:
- missing auth
- invalid auth
- valid auth
- forged browser actor
- configured operator identity reaches decision handling

==================================================
F4 — PRESERVE EvidenceClass THROUGH HTTP
==================================================

Current defect:

Trust UI sends EvidenceClass,
but coordinator/http/types.go does not carry it,
so production HTTP silently drops it.

Fix:

Add EvidenceClass to:

coordinator/http/types.go
SubmitDecisionRequest

Then explicitly map it in:

handleSubmitDecision
→ coordinator.DecisionRequest

Add a REAL HTTP handler test:

HTTP POST /decisions
→ JSON contains evidence_class
→ Coordinator receives same evidence class
→ retirement validation sees it

Do NOT rely only on direct SubmitDecision unit tests.

This test must exercise the production transport path.

==================================================
F5 — HUMAN DISCHARGE MUST USE ATTRIBUTED SOLVENT DISCHARGE
==================================================

Current defect:

handleRetireDebt calls:

POST /v1/beliefs/{id}/debt/retire

That is the lower-level unattributed primitive.

The human path MUST use:

POST /v1/discharge

Use the existing Solvent DischargeRequest contract:

- ScenarioID
- BeliefID
- ObligationKey
- InstrumentRef
- DischargedBy

Use:
- DischargedBy = authenticated/configured operator principal
- InstrumentRef = stable Coordinator decision/reference identifier

Do NOT remove the lower-level RetireDebt primitive unless clearly
unnecessary.

But human adjudication MUST NOT use it.

Verify the authoritative Solvent record contains:
- actor/operator identity
- obligation/debt
- belief
- scenario
- instrument/decision reference

Add tests proving:
- Discharge endpoint called
- RetireDebt endpoint NOT called on human path
- audit attribution exists
- operator identity cannot be overridden

==================================================
F6 — PACK LOOKUP FAILURE MUST FAIL CLOSED
==================================================

Current defect:

Pack validation can be skipped when:
- registry absent
- pack unresolved
- pack reference missing
- hardcoded extraction fails

Fix:

Pack resolution failure is a hard refusal/error.

Never:

pack lookup failed
→ continue
→ Solvent retire

Instead:

pack resolution failed
→ refusal/error
→ no Solvent mutation

Also remove any hardcoded:

("bmist", "1.0.0")

lookup used as the general mechanism.

Resolve Pack from actual scenario/task metadata already available in the
architecture.

Do NOT let browser input arbitrarily select a Pack.

For Phase 8 BM-IST:
- scenario/task metadata must establish the pack
- Coordinator reads that authoritative metadata
- PackRegistry resolves it

Add tests:
- valid pack
- missing pack
- unknown pack
- unregistered pack
- malformed pack reference
- Pack registry unavailable

==================================================
F7 — RETRACTION MUST UPDATE LINKED CONDUCTOR TASKS
==================================================

Current defect:

Solvent RetractCascade only changes Solvent state.
Conductor tasks remain active/proposed/etc.

The approved Phase 8 behavior is:

contradiction
→ human retract
→ Solvent RetractCascade
→ linked task becomes terminal cancelled
→ Insights derives Dead End

Inspect the actual Conductor API.

If a cancellation endpoint already exists:
- use it.

If it does NOT exist:
- add the smallest Conductor REST cancellation interface necessary
- no new schema
- no direct DB access
- document this as a Growth Gate interface extension

Do NOT create hidden DB writes.

Coordinator retraction flow must explicitly:
1. invoke Solvent RetractCascade
2. determine scenario-linked / belief-linked Conductor tasks
3. cancel the relevant task(s)
4. allow Insights to derive Dead End structurally

Add integration tests:
- belief retracted
- task cancelled
- cancelled is terminal
- active task with retracted belief alone is not considered dead end
- unrelated tasks are not cancelled

==================================================
F8 — UI MUST SHOW ACTUAL RETIREMENT RULE
==================================================

Current UI uses generic text:
"Applicable retirement rules apply"

Replace this with actual Pack-derived information.

For each debt item show:
- debt item
- required evidence class
- retirement rule identifier/name
- human-readable rule requirement where available

Example:

needMap
Required evidence:
reproducible_artifact
Rule:
map_check

Do NOT hardcode these values in the template.

Data must come from Coordinator/Pack state.

The same rule displayed by the UI must be the rule enforced by
Coordinator.

Add test:
UI displays actual rule for the current debt item.

==================================================
F9 — SORT RCP ACTIVITY DETERMINISTICALLY
==================================================

Current defect:

mergeActivity concatenates activity arrays without deterministic sort.

Fix:

merge Conductor + Solvent activity
→ sort by created_at
→ stable deterministic tie-breaker

Define the ordering explicitly.

Recommended:
1. created_at ascending
2. source deterministic order
3. stable ID as final tie-breaker

Or another deterministic equivalent if the repository already has a
canonical timestamp ordering rule.

Add test with interleaved timestamps and ties.

Repeated identical reads must return identical activity ordering.

==================================================
F10 — REOPEN MUST PRESERVE DERIVES LINEAGE
==================================================

Current defect:

handleReopen creates successor belief but no edge.

Fix:

retracted original belief
→ create successor belief
→ create edge:

parent = original/retracted belief
child = new belief
kind = derives

Original remains retracted.

Add tests:
- successor created
- derives edge exists
- original retained
- lineage visible through RCP

Use the newly implemented Solvent REST edge endpoint.
Do NOT import Solvent kernel into Coordinator.
Do NOT write directly to DB.

==================================================
REGRESSION / EXISTING EDGE GROWTH GATE
==================================================

The Phase 8 Solvent Growth Gate adds:

POST /v1/beliefs/{parent_id}/edges

Keep:
- parent exists validation
- child exists validation
- parent != child
- kind derives|contradicts
- uniqueness
- no MCP edge tool for agents
- Coordinator → REST → Solvent
- no direct DB

Do not change the public agent MCP surface.

==================================================
PRODUCTION-PATH TESTS — MANDATORY
==================================================

The adversarial review found that the previous 149 tests were
insufficient because they primarily called Coordinator methods directly.

Add real HTTP integration tests covering:

Trust UI-shaped JSON
→ Coordinator HTTP handler
→ Coordinator decision logic
→ Solvent client
→ expected result

At minimum:

1. valid needMap retirement
2. wrong evidence class
3. missing evidence
4. unknown debt
5. unknown pack
6. missing operator authentication
7. forged browser actor
8. valid authenticated operator
9. attributed /v1/discharge called
10. bare RetireDebt NOT called
11. promotion while debt remains
12. retraction path
13. reopen path

Do not merely update unit-test expectations.

==================================================
TEST THE ACTUAL AUTHORITY INVARIANTS
==================================================

After implementation, verify:

### Agent authority
Agent has only:
- argus.get_context
- argus.submit_packet

Agent cannot:
- retire
- promote
- discharge
- authorize
- mutate Conductor

### Human authority
Only authenticated/configured operator can call the consequential
decision boundary.

### Solvent authority
Only Solvent decides:
- debt state
- promotion
- retraction
- authorization

### Evidence authority
The client cannot invent evidence class by sending a string.

==================================================
EBP v2.1
==================================================

Confirm after the fixes:

Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
No final-truth claim promotes.
Humans discharge debt.
Accounting remains subordinate to research.

Do not introduce automatic scientific adjudication.

==================================================
WHAT NOT TO DO
==================================================

Do NOT:
- add a new database
- add a new authority layer
- add Solvent tools to OpenCode MCP
- allow browser actor identity
- trust submitted evidence_class without evidence verification
- bypass Pack validation
- call Solvent kernel directly from Coordinator
- write directly to Conductor DB
- write directly to Solvent DB
- convert agent findings into operator_asserted automatically
- redesign RCP
- add UI surfaces unrelated to Phase 8
- weaken tests merely to make them pass

==================================================
EXECUTION ORDER
==================================================

Implement in this order:

1. Inspect existing auth/discharge/pack/evidence interfaces.
2. Fix HTTP EvidenceClass propagation.
3. Fix unknown debt + pack fail-closed behavior.
4. Implement actual qualifying evidence verification.
5. Switch human path to attributed Discharge.
6. Authenticate /decisions.
7. Fix retraction → Conductor cancellation.
8. Fix REOPEN derives edge.
9. Fix UI rule rendering.
10. Fix deterministic RCP activity ordering.
11. Add production-path integration tests.
12. Run all tests with race detector where applicable.
13. Run go vet.
14. Perform a final source-level audit against every F1–F10.

==================================================
FINAL REPORT
==================================================

Produce a report:

PHASE8_ADVERSARIAL_FIX_REPORT.md

Include:

1. Files changed
2. F1–F10 disposition
3. Production-path tests added
4. Test totals
5. go test result
6. go test -race result
7. go vet result
8. Any remaining limitations
9. Explicit confirmation that:
   - unknown debt fails closed
   - evidence is actually verified
   - /decisions is authenticated
   - EvidenceClass survives HTTP
   - human discharge is attributed
   - pack failure fails closed
   - retraction cancels linked tasks
   - reopen preserves derives lineage
   - RCP ordering is deterministic
   - agent capability surface remains unchanged

If any F1–F10 remains unresolved, report the exact blocker instead of
claiming success.

Do not stop at "tests pass."
Trace the production authority path.

Final acceptance target:

Human
→ authenticated Trust UI
→ Coordinator
→ Pack validation
→ actual evidence verification
→ attributed Solvent discharge
→ Solvent audit

and:

Fresh OpenCode
→ ARGUS MCP
→ Coordinator
→ no direct authority mutation

Only after that is Phase 8 eligible for the actual dry run.
