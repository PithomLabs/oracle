You are drafting the IMPLEMENTATION PLAN for Phase 8 of the ARGUS Trust Verification POC.

Do NOT implement code yet.

Your deliverable is a concrete, reviewable Phase 8 implementation plan, including architecture, files/modules to create or modify, interfaces, API/MCP contracts, workflow, UI behavior, security/capability boundaries, test strategy, acceptance criteria, and explicit non-goals.

==================================================
BACKGROUND
==================================================

ARGUS is a trust-verification POC built around these frozen architectural boundaries:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Agent / OpenCode
  = agency / work

Conductor
  = operational coordination and task state

Solvent
  = epistemic authority, belief/evidence/debt ledger,
    authorization and consequential promotion gates

Trust UI
  = human control surface

EBP v2.1
  = epistemic operating discipline

Research Context Protocol (RCP)
  = thin read-only context protocol for agents

The current research seed is BM–IST Synthesis v5.1, an organized research program with claims, dependencies, evidence, debt, adversarial review, gates, and explicit stop/no-go conditions.

EBP v2.1 doctrine is mandatory:

Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
No final-truth claim may be promoted.
Accounting must never become the work.

Human beings discharge epistemic debt. Agents may produce evidence and proposed work, but agents do NOT directly retire debt or promote claims.

==================================================
PHASE 8 GOAL
==================================================

Phase 8 is a THIN FIRST DRY RUN, not a full research-management platform.

The first run should demonstrate:

human creates/selects a research task
→ Work OpenCode reconstructs context
→ Work OpenCode performs bounded research
→ submits an EBP packet
→ ARGUS records it through Coordinator
→ fresh Adversarial OpenCode reconstructs current state
→ adversarial review occurs
→ findings/evidence/debt are recorded
→ human sees the current situation
→ human reviews and discharges debt
→ Solvent authoritatively records debt retirement
→ human may request promotion
→ Solvent decides whether promotion is permitted

Use one BM–IST slice for the dry run, preferably Gate G0, because it is concrete and has claims, dependencies, evidence and proof obligations.

==================================================
AUTHORITATIVE EXISTING ARCHITECTURE
==================================================

Do NOT invent a new metadata database or duplicate existing state.

Conductor already owns operational state:
- project
- task
- dependency
- activity
- task status
- governance_ref

Solvent already owns epistemic/authority state:
- belief
- belief_edge
- evidence
- debt
- action_intent
- audit_activity
- promotion / retraction / authorization lifecycle

Use the existing Conductor and Solvent APIs/interfaces wherever possible.

Extend Conductor or Solvent ONLY when an existing interface is genuinely insufficient for the Phase 8 dry run.

Do not add direct UI→database writes.
Do not add direct agent→Solvent writes.
Do not modify Conductor's frozen domain model merely to support UI convenience.
Do not create an RCP database/schema.

==================================================
RCP — FROZEN DIRECTION
==================================================

RCP = Research Context Protocol.

RCP is NOT another persistence layer and NOT another research ontology.

RCP is a thin read-only context projection assembled by Coordinator from Conductor + Solvent.

Canonical API contract:

GET /v1/context/{task_id}

Conceptual response:

{
  "protocol": "RCP/v1",
  "task": {},
  "dependencies": [],
  "epistemic": {
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  },
  "activity": []
}

Important constraints:

1. RCP v1 scope is a full scenario/task projection; do NOT invent sophisticated task-anchored graph closure yet.
2. Cross-ledger consistency is eventually consistent.
3. Each returned fact/source should identify whether it came from Conductor or Solvent where useful.
4. Projection lag must be visible rather than silently synthesized.
5. Existing artifact hashes must have a usable artifact resolution path where one exists; hash-only context is insufficient for a fresh agent.
6. RCP returns context. It does not prescribe scientific next steps.
7. "Next smallest useful move" is guidance for the UI/agent/domain pack, not Solvent state.

==================================================
OPEN CODE / MCP BOUNDARY — CRITICAL
==================================================

OpenCode is the agent harness.

There are two separate roles:

WORK AGENT
ADVERSARIAL AGENT

Both agents must use the same narrow ARGUS-facing MCP surface.

The intended MCP surface is conceptually:

- argus.get_context
- argus.submit_packet

The agents MUST NOT have:
- Solvent MCP mutation tools
- direct Solvent MCP lifecycle tools
- Conductor write tools
- direct DB access

In particular, the Solvent repository already exposes MCP functionality such as:
- retire_debt
- promote
- authority lifecycle operations

These capabilities MUST NOT be present in the Work or Adversarial OpenCode tool surface.

The design must make the capability boundary enforceable, not merely a prompt instruction.

Phase 8 must include a negative test proving:

"An agent attempting direct Solvent mutation cannot do so because the capability is absent from its tool surface."

MCP is the agent-facing transport/adapter.
HTTP/API is the canonical integration contract.
Coordinator is the authority-aware orchestration boundary.

Expected architecture:

OpenCode
  ↓ MCP
ARGUS adapter
  ↓ HTTP
Coordinator
  ↓
Conductor + Solvent

==================================================
EBP WORKFLOW RULES
==================================================

Agents may:
- inspect context
- perform research
- create code/notebooks/artifacts
- produce evidence
- propose claim updates
- submit EBP packets
- identify obstructions/counterexamples

Agents may NOT:
- directly retire debt
- directly promote beliefs
- directly authorize
- directly retract authority state
- decide scientific faithfulness
- declare a claim "AI verified"

Human:
- reviews evidence
- decides whether debt is discharged
- may request promotion
- is the epistemic adjudicator

Solvent:
- records authoritative debt retirement
- enforces promotion gates
- remains canonical authority state

The UI must NEVER present:
"AI verified debt"

Instead present language such as:
"Adversarial agent supplied evidence relevant to this debt."
"Human decision required."
"Debt retired by operator."

==================================================
DEBT RETIREMENT
==================================================

For Phase 8 choose this rule:

Coordinator mechanically validates that the offered evidence satisfies the applicable Domain Pack retirement rule before invoking Solvent RETIRE_DEBT.

Coordinator does NOT decide whether the underlying scientific claim is correct.

The human still performs the actual adjudication.

Therefore:

Agent evidence
→ Human review
→ Coordinator checks mechanical retirement rule
→ Solvent records retirement

Do NOT implement partial-debt semantics.

Debt items are atomic in Phase 8:
present / retired.

Record this as an explicit Phase 8 implementation decision.

==================================================
REOPEN / DEAD-END / ATTENTION
==================================================

REOPEN:
- preserve lineage
- new belief derives from the retracted origin
- do not create an unrelated fresh belief with no lineage

Dead end:
DO NOT invent a free-form dead_end field or trust agent-written dead_end prose.

For Phase 8 define dead end structurally as:

Rejected Conductor task whose governance/decision record points at a retracted or contradicted belief.

Do NOT use merely "cancelled" as the semantic dead-end definition.

Attention:
- owned by Conductor/workflow, not Solvent
- do not add an epistemic attention field to Solvent
- Phase 8 UI priority is POC-scoped manual ordering in Conductor unless an existing deterministic mechanism already exists
- do not invent HIGH/MEDIUM/LOW metadata in Solvent

==================================================
UI — FROZEN TWO-SURFACE DESIGN
==================================================

Sidebar has exactly:

Insights
Debts

Do not build the larger mythology/navigation surface yet.

------------------------------------------
1. INSIGHTS
------------------------------------------

Purpose:
situational awareness of the research program.

Show compact program statistics, for example:
- open claims
- open debt
- active tasks
- adversarial findings
- promoted claims
- refused transitions
- dead-end tasks

Then an ordered Kanban/task attention view.

IMPORTANT:
Priority must not be invented as a new database field.

For Phase 8, use the existing Conductor task ordering/state as the attention source, and document that this is a POC-scoped manual ordering unless an existing deterministic field is already available.

Dead-end statistics must use the structural definition above.

------------------------------------------
2. DEBTS
------------------------------------------

Purpose:
human epistemic adjudication.

For each open debt show:
- claim
- debt item
- why the debt exists
- supporting evidence
- adversarial findings
- applicable retirement rule/evidence class
- proposed discharge
- prior history
- current state

Human action:
- RETIRE DEBT
- keep open / return

Promotion is NOT performed by the Debts screen unless the existing Coordinator/decision workflow already makes that a separate explicit action.

The UI is a projection and control surface, not an authority store.

==================================================
REPOSITORY / SOURCE MATERIAL TO READ
==================================================

Before drafting the implementation plan, inspect the actual repositories and interfaces in the working tree.

Read as necessary:

1. Solvent repository
   - existing belief/debt/evidence APIs
   - promotion gate
   - RETIRE_DEBT endpoint
   - refusal behavior
   - audit_activity
   - existing solvent-mcp capabilities
   - existing artifact/evidence interfaces

2. Conductor repository
   - task/project/dependency/activity APIs
   - task lifecycle/state machine
   - governance_ref
   - current HTTP API
   - existing task ordering/metadata
   - determine whether any existing priority/order mechanism can support the POC without schema changes

3. oracle/coordinator
   - current Phase 6/7 implementation
   - Coordinator interfaces
   - HTTP wrapper
   - packet compilation
   - human decision paths
   - PackRegistry
   - artifact reader
   - current mocks and tests

4. oracle/domain-pack
   - BM-IST pack
   - debt vocabulary
   - retirement rules
   - evidence classes

5. oracle/packet
   - EBP packet schema/validation
   - work/adversarial role semantics

6. oracle/verifier
   - trusted artifact registration/read path

7. oracle/corpus
   - current BM-IST corpus manifest
   - current authoritative seed artifact

8. BM–IST Synthesis v5.1
   - use attached research document as domain content
   - use Gate G0 as thin dry-run candidate unless the repository reveals a materially better already-supported slice

9. EBP v2.1
   - especially:
     Ideas enter free
     Promotion costs debt
     Debt does not kill
     Debt is forever payable
     New evidence creates new debt
     Human debt discharge
     No final-truth promotion
     Accounting must never become the work

10. Previously reviewed architecture / Phase 8 baseline
   - use the current conversation decisions as frozen constraints

Do not substitute a new architecture because another design looks more elegant.

==================================================
WHAT THE IMPLEMENTATION PLAN MUST CONTAIN
==================================================

Draft a plan with these sections:

1. Phase 8 objective
2. Current architecture and boundaries
3. Existing interfaces discovered in Conductor/Solvent
4. Minimal changes required in Solvent/Conductor, if any
5. RCP protocol design
6. Coordinator RCP endpoint design
7. MCP adapter/tool surface
8. Work Agent workflow
9. Adversarial Agent workflow
10. EBP packet lifecycle
11. Human debt-discharge workflow
12. Solvent promotion/refusal workflow
13. UI architecture
14. Insights view
15. Debts view
16. Artifact/evidence retrieval
17. Dead-end derivation
18. Cross-ledger consistency handling
19. Security/capability boundary
20. Error/refusal semantics
21. Observability/audit requirements
22. Files/modules to create or modify
23. Test strategy
24. Detailed test matrix
25. Acceptance criteria
26. Explicit non-goals / deferred work
27. Risks and remaining decisions

==================================================
TEST SUITE — MUST BE PART OF THE PLAN
==================================================

The plan must include actual test cases, not merely "add tests."

At minimum cover:

RCP:
- valid context retrieval
- unknown task
- no epistemic state yet
- task with multiple dependencies
- task with existing beliefs/evidence/debt
- source attribution
- eventual-consistency/projection lag behavior
- artifact reference resolution
- deterministic response shape

MCP:
- get_context exposed
- submit_packet exposed
- Solvent mutation tools NOT exposed
- Conductor write tools NOT exposed
- negative test: agent cannot directly mutate Solvent

Agent reconstruction:
- fresh Work Agent can reconstruct context with no prior conversational history
- fresh Adversarial Agent can reconstruct context with no prior conversational history
- completed work is visible
- previous adversarial findings are visible
- dead-end work is visible and structurally identifiable
- unresolved debt is visible
- artifact can be resolved

EBP:
- ideas/claims enter without promotion
- initial debt appears
- debt blocks promotion
- human discharge retires debt
- agent cannot retire debt directly
- new evidence can create/reinstate debt
- final-truth claim cannot be promoted
- promotion remains Solvent-controlled

Retirement rules:
- compliant evidence class succeeds
- wrong evidence class is refused
- missing retirement rule is handled explicitly
- scientific correctness is NOT inferred automatically by Coordinator

REOPEN:
- retracted belief creates successor
- derives lineage is recorded
- original remains in history

Dead end:
- structural derivation works
- cancelled task alone does not count as dead end

UI:
- Insights statistics render from existing state
- attention ordering uses Conductor state, not new Solvent priority
- Debts shows why debt exists
- Debts shows evidence/adversarial findings
- UI never presents "AI verified debt"
- human discharge invokes Coordinator path
- UI does not write directly to Solvent

Regression:
- existing coordinator tests remain green
- existing Solvent tests remain green
- existing Conductor tests remain green
- no modifications to unrelated reference-loop/domain-pack/packet/verifier/corpus unless explicitly justified

Concurrency:
- RCP retrieval race-safe where applicable
- existing coordinator idempotency/decision semantics preserved

==================================================
NON-GOALS
==================================================

Do NOT implement:
- new research database
- generalized research graph database
- task priority model in Solvent
- partial debt
- automatic debt retirement from natural language
- autonomous promotion
- autonomous scientific adjudication
- general lesson/knowledge graph
- full BM-IST research dashboard
- all mythology screens
- direct OpenCode access to Solvent/Conductor writes
- crash-recoverable persistent projection system
- broad RCP graph traversal beyond the agreed v1 scope

==================================================
IMPORTANT ARCHITECTURAL TEST
==================================================

The plan must answer this question clearly:

"If I start a completely fresh OpenCode process as an adversarial agent, with no conversation history, how exactly does it learn what has already been done?"

The answer should be mechanically traceable:

OpenCode
→ MCP get_context
→ ARGUS RCP endpoint
→ Coordinator
→ Conductor + Solvent
→ canonical context response
→ adversarial agent

Then answer:

"If the agent tries to bypass this and call Solvent directly, what prevents it?"

The answer must be capability/tool-surface absence, not merely prompt instructions.

==================================================
DELIVERABLE
==================================================

Produce a document named:

PHASE8_IMPLEMENTATION_PLAN.md

Do not implement the plan.

Do not modify repositories.

Do not create speculative interfaces that are not grounded in the existing repositories.

First inspect the existing APIs/interfaces and then produce the smallest implementation plan that satisfies the above.

The plan will be reviewed before any Phase 8 coding begins.