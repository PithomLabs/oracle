You are revising the previously drafted PHASE8_IMPLEMENTATION_PLAN.md for the ARGUS Trust Verification POC.

Do NOT implement code yet.

Produce the corrected final implementation plan only.

Treat the existing Phase 8 plan plus the consolidated adversarial review below as the authoritative review baseline. Preserve everything already correct. Do not broaden the architecture. Do not invent new persistence or metadata systems.

==================================================
FROZEN ARCHITECTURE
==================================================

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

OpenCode
  = agent/work capability

Conductor
  = operational workflow/task state

Solvent
  = epistemic authority, evidence, debt, promotion,
    retraction, authorization and audit state

Trust UI
  = human observation and adjudication

EBP v2.1
  = epistemic operating law

RCP
  = thin read-only agent context projection

MCP
  = agent-facing transport

HTTP/API
  = canonical integration contract

No new research database.
No new Solvent schema.
No new Conductor schema.
No direct DB access from UI or agents.

==================================================
EBP v2.1 — MANDATORY
==================================================

Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
No final-truth claim may be promoted.
Accounting must never become the work.

Human beings discharge epistemic debt.

Agents may produce candidate work, evidence, counterexamples,
artifacts and proposed findings.

Agents MUST NOT:
- retire debt directly
- promote beliefs
- authorize consequential actions
- mutate Solvent directly
- mutate Conductor directly

The authoritative human transition must be recorded by Solvent.

==================================================
PHASE 8 PURPOSE
==================================================

Phase 8 is a thin first dry run against one real BM–IST research
slice, preferably Gate G0.

Required loop:

human task
→ Work OpenCode
→ RCP context
→ bounded research
→ EBP work packet
→ ARGUS
→ Conductor + Solvent persistence
→ fresh Adversarial OpenCode
→ RCP reconstruction
→ adversarial attack
→ EBP adversarial packet
→ ARGUS
→ human reviews Insights / Debts
→ human discharges debt
→ Solvent records attributed retirement
→ human requests promotion
→ Solvent decides gate

The dry run must prove that a completely fresh adversarial OpenCode
process can reconstruct prior work without conversational memory.

==================================================
C1 — HUMAN DISCHARGE MUST BE ATTRIBUTED
==================================================

The existing Solvent interfaces include:

POST /v1/beliefs/{id}/debt/retire

and:

POST /v1/discharge

The human adjudication path MUST use the attributed discharge path.

Use conceptually:

Trust UI
→ Coordinator
→ Solvent POST /v1/discharge

with:
- scenario_id
- belief_id
- obligation/debt key
- operator identity
- instrument/decision reference

The bare debt-retire endpoint may remain a lower-level primitive,
but the human adjudication workflow must use attributed discharge.

The authoritative Solvent audit must contain the operator identity.

==================================================
C2 — TASK / BELIEF LINKAGE
==================================================

Conductor governance_ref is write-once and immutable.

Therefore do NOT create a human task with an empty belief_id and
later attempt to mutate governance_ref.

For Phase 8:
- anchor the task to the scenario
- use the scenario as the governing context
- use a dedicated one-governing-belief scenario for the G0 dry run
- do not modify Conductor schema

Example:

governance_ref:
{
  "provider": "solvent",
  "reference_id": "track-g0"
}

Do NOT depend on adding a future belief_id into immutable metadata.

Make explicit in the plan that Phase 8 uses one governing belief per
dry-run scenario to keep the task↔epistemic relationship deterministic.

Do not create duplicate "Formalize Gate G0" tasks merely to obtain a
later governance reference.

==================================================
C3 — HUMAN ACTOR IDENTITY
==================================================

Human authority must not be inferred from browser content.

Use the established server-side operator principal pattern:

ARGUS_OPERATOR_PRINCIPAL_ID

Requirements:
- Trust UI injects the configured operator identity
- Coordinator validates it
- Coordinator uses it for human decisions
- attributed Solvent discharge receives this identity
- missing operator identity fails closed

DecisionRequest should NOT trust an arbitrary actor supplied by the browser.

The authoritative path is:

browser
→ Trust UI
→ Coordinator
→ configured operator principal
→ Solvent

Record the actor in:
- Coordinator DecisionRecord
- Solvent discharge/audit record

==================================================
H1 — ADVERSARIAL CHALLENGES
==================================================

Do not claim that adversarial findings can be queried from
provenance_class.

`provenance_class` is not the packet role.

For Phase 8:
- rename the Insights statistic to "Adversarial Challenges"
- derive it structurally from `contradicts` edges

The adversarial workflow MUST therefore create an actual
`contradicts` relationship against a target belief in at least one
dry-run branch.

Do not add a new Solvent role field.

==================================================
H2 — RCP MUST DISTINGUISH EMPTY FROM UNAVAILABLE
==================================================

The current sample implementation incorrectly ignores client errors.

NEVER do:

deps, _ := ...
beliefs, _ := ...

A Solvent or Conductor outage must not look like empty state.

RCP response shape remains deterministic, but sections need
availability/degradation information.

Example:

"epistemic": {
  "available": false,
  "reason": "solvent_unavailable",
  "beliefs": [],
  "evidence": [],
  "edges": [],
  "debt": [],
  "intents": []
}

Likewise for Conductor where appropriate.

Rule:

UNKNOWN / UNAVAILABLE ≠ EMPTY

The agent must never infer "no research exists" from a failed
backend read.

==================================================
H3 — ACTUALLY EXERCISE RETRACTION / DEAD-END
==================================================

The dry run must not only exercise:

entry
→ evidence
→ debt
→ refusal
→ discharge
→ promotion

It must also exercise one adversarial branch:

work
→ adversarial contradiction
→ `contradicts` edge
→ human retraction decision
→ Solvent cascade
→ terminal cancelled task / dead-end derivation
→ Insights shows the structurally derived dead end

Keep this branch thin.

Do not add a new UI surface.

Use the existing Coordinator/Solvent lifecycle.

The implementation must prove that dead-end machinery is not merely
implemented-but-unused.

==================================================
H4 — RETIREMENT RULES MUST ACTUALLY BE ENFORCED
==================================================

The BM-IST Domain Pack declares retirement rules such as:

needMap
  → reproducible_artifact / map_check

needInvariant
  → reproducible_artifact / invariant_check

needToyCheck
  → reproducible_artifact / toy_model_check

needNullModel
  → operator_asserted / scope_clarification

needObstruction
  → reproducible_artifact / obstruction_construction

needFaithfulnessReview
  → operator_asserted / faithfulness_review

The Coordinator MUST mechanically validate the offered evidence
against the required evidence class before invoking Solvent.

Flow:

human selects debt
→ identify applicable retirement rule
→ identify evidence offered for discharge
→ compare evidence class
→ mismatch = refusal
→ match = invoke attributed Solvent discharge

The Coordinator does NOT judge scientific correctness.

Do not describe the retirement rules as normative if they are merely
displayed.

Do not allow a human click alone to satisfy a rule requiring a
specific evidence class.

==================================================
AGENT EVIDENCE PROVENANCE
==================================================

Do NOT treat autonomous agent output as automatically
`operator_asserted`.

An agent is not the human operator.

For Phase 8:

agent analysis
→ candidate evidence

Candidate evidence becomes authoritative retirement material only
through:
- a reproducible artifact satisfying the pack rule, OR
- explicit human attestation where the pack permits operator_asserted

Do not overload `operator_asserted` to mean "agent acted on behalf
of operator."

==================================================
H5 — VERIFY EDGE API + LISTINTENTS
==================================================

The plan must explicitly verify the actual Solvent endpoint for
belief-edge creation before implementing the persistence path.

Expected shape is conceptually:

POST /v1/beliefs/{parent_id}/edges

If the endpoint exists:
- use it

If it does not:
- STOP
- invoke the existing Growth Gate
- do NOT use direct SQL as a workaround

Also add the missing Coordinator Solvent client method required by
RCP:

ListIntents

The RCP implementation must not describe an interface that the
client layer does not actually provide.

==================================================
H6 — PACKET REFERENCE GRAMMAR
==================================================

Freeze the existing EBP reference grammar:

local:<id>
  = intra-packet reference

canonical:belief:<uuid>
  = reference to an already-existing Solvent belief

Use one consistent field naming convention.

Fix all plan examples so work and adversarial packets use the same
reference rules.

Do not invent free-form semantic claim matching.

==================================================
H7 — IDEMPOTENCY SCOPING
==================================================

The canonical packet idempotency hash MUST preserve scenario
isolation.

Include scenario identity in the canonical hash.

Continue to exclude runtime-only values such as:
- packet_id
- run_id
- timestamps
- runtime metadata

Requirement:

same content + same scenario
→ deduplicates

same content + different scenario
→ does NOT deduplicate across scenario boundary

Preserve the existing atomic NEW / IN_FLIGHT / COMPLETED / FAILED
semantics from Phase 6.

==================================================
H8 — DEAD-END DEFINITION
==================================================

Use the actual Conductor lifecycle.

The real terminal task state is `cancelled`.

`rejected` is not a terminal dead-end state if the existing Conductor
state machine returns rejected work to active.

Therefore:

Dead End =
terminal cancelled Conductor task
+
governance/scenario linkage to
a retracted or contradicted governing belief

Do not use:
- arbitrary "dead_end" text
- agent-written dead-end metadata
- cancelled alone

Dead-end must be structurally derived.

Add/update the reconciliation record documenting the actual Conductor
lifecycle versus any older frozen design documentation.

==================================================
H9 — INSIGHTS PRIORITY
==================================================

Use the EXISTING Conductor `priority` field:

low | medium | high | critical

It is manually assigned and POC-scoped.

Insights ordering is client-side:

critical
→ high
→ medium
→ low
→ deterministic secondary ordering such as created_at

Do not add priority to Solvent.
Do not add a new priority schema.

Any SQL shown in the plan is conceptual only; Trust UI and Coordinator
must not access Conductor DB directly.

==================================================
H10 — TRUST UI MUST REMAIN A SEPARATE MODULE
==================================================

Preserve the frozen two-Go-module architecture.

Trust UI belongs in:

trust-ui/
  go.mod

Do NOT place the Trust UI inside oracle/ unless the frozen architecture
has explicitly been changed.

Trust UI:
- Go html/template
- vanilla HTML
- vanilla CSS
- vanilla JS
- no React
- no Vue
- no HTMX
- communicates with Coordinator via HTTP
- never writes directly to Solvent or Conductor

Keep the existing Solvent wizard style as the visual/technical pattern.

==================================================
RCP V1
==================================================

Canonical identity:

RCP/v1

Canonical endpoint:

GET /v1/context/{task_id}

RCP is a read projection, not a store.

Scope:
- full scenario projection for Phase 8
- no sophisticated graph traversal
- no new research graph database

Conceptual response:

{
  "protocol": "RCP/v1",
  "task": {},
  "dependencies": [],
  "epistemic": {
    "available": true,
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  },
  "activity": []
}

Objects are projections of Conductor/Solvent objects.

Tag source where useful:
- conductor
- solvent

Cross-ledger consistency is eventually consistent.

Projection gaps must be visible as lag/degradation, not silently
converted into empty state.

Artifact evidence must expose a usable artifact_ref where available.

==================================================
MCP CAPABILITY BOUNDARY
==================================================

OpenCode Work and Adversarial processes receive exactly:

argus.get_context
argus.submit_packet

They MUST NOT have:
- solvent-mcp
- Conductor write MCP tools
- direct DB access

The capability boundary must be enforced by tool-surface configuration,
not prompt language.

Add tests:
- only two ARGUS tools exposed
- solvent_* tools absent
- conductor_* tools absent
- attempted direct Solvent mutation fails
- attempted direct Conductor mutation fails

==================================================
EBP WORKFLOW
==================================================

WORK AGENT:

Human task
→ argus.get_context
→ research
→ EBP work packet
→ argus.submit_packet

ADVERSARIAL AGENT:

fresh OpenCode process
→ argus.get_context
→ sees current task/dependencies/beliefs/evidence/debt/activity
→ attacks unresolved work
→ creates contradiction/finding/evidence
→ EBP adversarial packet
→ argus.submit_packet

HUMAN:

Insights
→ Debts
→ inspect evidence/findings/rule
→ discharge debt
→ Coordinator validates mechanical rule
→ Solvent records attributed discharge

PROMOTION:

human requests promotion
→ Coordinator
→ Solvent gate
→ promotion or refusal

The UI/agent NEVER declares promotion authoritative.

==================================================
UI
==================================================

Sidebar exactly:

Insights
Debts

Insights:
- compact program statistics
- Adversarial Challenges
- open claims
- open debt
- active tasks
- promoted
- refusals
- structurally derived dead ends
- ordered task Kanban based on Conductor priority

Debts:
- claim
- debt item
- why debt exists
- evidence
- adversarial challenges
- applicable retirement rule
- proposed discharge
- history
- RETIRE DEBT
- KEEP OPEN

Never display:
"AI verified debt"
"Agent retired debt"
"Automatically discharged"

Use language such as:
"Adversarial agent supplied evidence relevant to this debt."
"Human decision required."
"Debt retired by operator."

==================================================
PERSISTENCE
==================================================

`CompilePacket` may remain internally responsible for:
- validation
- canonicalization
- idempotency
- mapping local refs

But `submit_packet` MUST persist real state.

Persistence order must be explicitly documented.

Expected path:

EBP packet
→ Coordinator validation
→ Solvent beliefs
→ Solvent edges
→ Solvent evidence
→ Conductor tasks

Use existing REST clients.

No direct DB writes.

Add failure-injection tests for partial persistence:
- belief created, evidence persistence fails
- error is returned
- partial state is observable
- retry does not silently duplicate completed objects
- existing idempotency semantics are preserved

Do not introduce distributed transactions in Phase 8.

==================================================
REOPEN
==================================================

Preserve the existing architecture:

retracted belief
→ REOPEN
→ new belief
→ derives edge from original

Original remains in history.

Lineage must be preserved.

==================================================
TEST SUITE
==================================================

Preserve the existing detailed tests and ADD/UPDATE the following:

RCP:
- valid context
- unknown task
- empty state
- dependency projection
- source attribution
- deterministic response
- Solvent unavailable
- Conductor unavailable
- projection lag visible
- artifact_ref present

MCP:
- exactly two ARGUS tools
- no Solvent tools
- no Conductor write tools
- direct mutation attempts fail

Persistence:
- belief persistence
- evidence persistence
- edge persistence
- task persistence
- canonical ID mapping
- duplicate submit
- cross-scenario idempotency isolation
- partial persistence failure/retry

EBP:
- initial six debt items
- debt blocks promotion
- attributed human discharge
- wrong evidence class refused
- matching evidence class accepted
- agent cannot discharge
- final-truth blocks promotion
- new evidence/debt reopens the gate

Adversarial:
- fresh agent sees previous work
- fresh agent sees previous evidence
- fresh agent sees open debt
- fresh agent sees prior adversarial findings
- adversarial packet creates contradicts edge

Retraction:
- human retracts
- Solvent cascade
- dependent intents cancelled where applicable
- original belief retained
- REOPEN creates derives lineage

Dead end:
- cancelled + retracted belief => dead end
- cancelled + active belief => not dead end
- rejected + active => not dead end
- contradicted governing belief => dead end

Human identity:
- missing operator identity fails closed
- configured operator identity reaches Solvent audit
- browser-supplied actor cannot override configured identity

UI:
- Insights statistics
- Conductor priority ordering
- dead-end count
- adversarial challenge count
- Debts rendering
- evidence/rule display
- no "AI verified" language
- discharge goes UI → Coordinator → Solvent
- no direct DB writes

Regression:
- existing Coordinator tests pass
- existing Solvent tests pass
- existing Conductor tests pass
- reference-loop unchanged
- domain-pack validation unchanged
- packet validation unchanged
- verifier unchanged
- corpus unchanged

Concurrency:
- RCP concurrent reads
- submit_packet duplicate concurrency
- concurrent discharge remains safe/idempotent

==================================================
IMPLEMENTATION PLAN STRUCTURE
==================================================

Revise the document while preserving its useful current structure.

Explicitly include:

1. Objective
2. Frozen architecture
3. Existing interfaces
4. Minimal changes
5. RCP
6. Coordinator
7. MCP
8. Work Agent
9. Adversarial Agent
10. Persistence
11. EBP/debt workflow
12. Human identity
13. Retraction/reopen/dead-end
14. Trust UI
15. Insights
16. Debts
17. Artifact retrieval
18. Eventual consistency
19. Error/refusal semantics
20. Security/capability boundary
21. Test strategy
22. Detailed tests
23. Acceptance criteria
24. Non-goals
25. Freeze/reconciliation notes
26. Risks

==================================================
FINAL REVIEW REQUIREMENT
==================================================

Before finishing the revised plan, perform an internal consistency
check.

Specifically verify that:

- every mechanism promised in acceptance criteria is actually described
  in implementation
- every displayed UI field has a derivable source
- every normative retirement rule is actually enforced
- every human-authority claim has an attributable actor
- every agent capability is explicitly constrained
- no UI/write path bypasses Coordinator
- no RCP failure is represented as empty state
- no dead-end semantics contradict Conductor's actual state machine
- no example uses inconsistent packet references
- no idempotency boundary crosses scenarios
- no new database/schema has been invented

Do NOT implement code.

Return the revised implementation plan as:

PHASE8_IMPLEMENTATION_PLAN.md

This is the final review candidate. After this revision, stop architectural expansion and wait for implementation approval.