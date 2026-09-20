## prompt

argus serve works now, I am pivoting again for our POC seed, it is now the combined Bohmian Mechanics + Invariant Set Theory + Asymptotic Safety but rethought from first principles (read attached files, those files are the collective POC seed), my goals are:

- ingest attached files as initial corpus of our research (so if someone tries our research program, they already have something to work on as baseline) and use native vector search of CockroachDB just like in our demo
- v6_1.md is the initial research program, revised to v6_2_adv.md when adversarial review was incorporated and now on version v6_3.md as our latest baseline version
- when I run argus serve and any AI agent (e.g. OpenCode), for our initial dry run the AI agent should work on this question:

The phrase **"one flow" is deliberately no longer used as an established architectural fact**. There is currently no constructed flow. The correct status is **one candidate flow architecture whose existence is itself a central theorem target**.

The central question is correspondingly sharp:

$$
\boxed{
\text{Can one explicit microscopic dynamical system generate the structures of quantum theory without assuming them?}
}
$$

The research program lives or dies on that question.


- the goal of the coding agent is to research on that question and break down to manageable tasks per EBP 2.1 "ideas enter free, promotion costs debt", hence the deliverable is to gather evidence, form the logic chain, form the belief lifecycle and list the debts per the agreed upon protocol (domain pack, EBP packet, EBP 2.1)

- when I access the /insights endpoint the AI agent should be able to populate that insights via Argus interface (MCP, API, whatever)

- when I access the debt endpoint, I should be able to discharge debts just like in the original hackathon demo

- when a separate AI agent does its adversarial role, it should be able to look up the current CockroachDB schema (or wherever it is stored), be able to reconstruct the context from the previous AI agent and start its adversarial review (attack the logic, debt list etc of the previous AI agent)

- the goal here is seamless developer experience (trying out Argus and Solvent should be frictionless as much as possible), the goal is to apply EBP 2.1 as physics research program and that this methodology is applicable to other domains)

- suggest any workflow to improve the experience

give prompt to the coding agent if the above requirements can now be served by Argus and Solvent, if not give prompt to draft implementation plan



## chatgpt

Yes, but **not completely with the current surface**.

The existing ARGUS/Solvent foundation already gives us most of the epistemic workflow: beliefs, evidence, debt, edges, promotion/retraction, `get_context`, `submit_packet`, Trust UI, and reconstruction of prior agent work. The original POC design explicitly uses `get_context` for a fresh adversarial agent to reconstruct beliefs, evidence, edges and debt, then submit an adversarial packet. 

Solvent also already has the important **corpus substrate**: `corpus_issue` stores text plus `VECTOR(1024)`, and corpus retrieval is deliberately separate from the authoritative belief ledger.  The prior corpus work also established that corpus membership does **not** confer epistemic authority. 

What is missing is the thin research interface connecting those pieces to the new BM+IST+AS seed:

```text
attached research corpus
        ↓
CockroachDB vector retrieval
        ↓
fresh-agent context
        ↓
EBP research packet
        ↓
belief / evidence / debt / graph
        ↓
Insights + Debts
        ↓
fresh adversarial agent
```

In particular, we need to establish how the new corpus is ingested, embedded and searched; how v6.3 becomes the initial research-program seed; how the current BM-IST pack evolves to cover the tripartite program; and how corpus hits become context without becoming authority.

I would **not implement this blindly yet**. Have the coding agent perform one bounded gap audit and produce the smallest implementation plan.

### Prompt for the coding agent

```text
You are now designing the NEXT minimal ARGUS/Solvent POC slice.

DO NOT IMPLEMENT CODE YET.

First inspect the actual current repository and determine whether the existing frozen ARGUS + Solvent foundation can support the new BM–IST–Asymptotic Safety research dry run with only thin interface-layer additions.

The architecture remains frozen:

    OpenCode
       |
       | MCP stdio
       v
    ARGUS — one Go process
       |
       v
    CockroachDB

Core invariant:

    CAPABILITY != WORK != AUTHORITY != EXECUTION

Do NOT:
- expand the Solvent authority kernel unless a genuine missing invariant is proven
- restore old Solvent/Conductor/Coordinator REST services
- create another database
- create a workflow engine
- create a separate vector database
- add unnecessary MCP tools
- turn the domain pack into the research database
- make corpus retrieval an authority mechanism
- implement code during this task

==================================================
CURRENT RESEARCH SEED
==================================================

The attached files constitute the initial research corpus:

- v6_1.md
- v6_2_adv.md
- v6_3.md
- adv_review(8).md
- adv_review2(3).md
- adv_review3(2).md
- writeup_v1.md

Interpret them as follows:

v6_3.md
    = current master research-program baseline

v6_2_adv.md
    = immediately preceding adversarially revised baseline

v6_1.md
    = prior research-program revision

adv_review*.md
    = adversarial research/review material

writeup_v1.md
    = explanatory research-program write-up

The central research question is:

    Can one explicit microscopic dynamical system generate
    the structures of quantum theory without assuming them?

The phrase "one flow" must NOT be treated as an established fact.
There is only a candidate flow architecture whose existence is itself
a central theorem target.

==================================================
EBP METHOD
==================================================

The research workflow follows EBP v2.1:

    Ideas enter free.
    Promotion costs debt.
    Debt does not kill.
    Debt is forever payable.
    New evidence creates new debt.
    Promotion never means truth.

Current debt classes include:

    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview
    needInitialCondition
    needRegularity

The system must preserve the distinction:

    corpus
        !=
    evidence
        !=
    belief
        !=
    authority

Corpus material is research context.
It does not become promoted truth merely because it was ingested.

==================================================
DESIRED USER EXPERIENCE
==================================================

A developer should be able to:

    argus serve

then start a fresh OpenCode research agent.

The agent should immediately have a seeded research program and something
real to work on.

The initial work should be centered on:

    "Can one explicit microscopic dynamical system generate
     the structures of quantum theory without assuming them?"

The agent should:

1. reconstruct existing ARGUS state;
2. retrieve relevant seed-corpus material using CockroachDB native vector search;
3. identify the current research frontier;
4. break the question into manageable EBP research tasks;
5. create candidate beliefs;
6. attach explicit debt;
7. collect evidence;
8. create derives relationships;
9. submit the resulting EBP packet;
10. make the resulting state visible through /insights and /debts.

Then a SECOND fresh OpenCode process should be able to:

1. call the same ARGUS context mechanism;
2. reconstruct what the first agent did;
3. see existing claims, evidence, debt and graph relationships;
4. retrieve the relevant corpus;
5. independently attack the logic, evidence and debt list;
6. submit adversarial findings;
7. create contradictions/challenges through the existing EBP path.

The human then reviews/discharges debt through the existing Trust UI.

==================================================
CRITICAL DESIGN PREFERENCE
==================================================

Prefer preserving exactly TWO MCP tools:

    argus.get_context
    argus.submit_packet

Do NOT add a third:

    argus.search_corpus

unless the audit proves that the existing two-tool surface cannot support
retrieval cleanly.

Preferred design:

    argus.get_context(task_id, optional research query)

returns:

    task
    current beliefs
    evidence
    edges
    debt
    activity
    availability
    relevant corpus hits

This keeps the agent-facing surface minimal.

==================================================
CORPUS / VECTOR REQUIREMENT
==================================================

Audit the existing Solvent corpus infrastructure.

There is already a corpus layer with:

    corpus_issue
    VECTOR(1024)
    belief_corpus_citation

and native CockroachDB vector retrieval.

Determine exactly what already exists for:

- corpus ingestion
- embedding generation
- embedding storage
- vector similarity search
- top-K retrieval
- provenance
- citation back to corpus records
- scenario isolation

Determine whether the existing infrastructure can be reused unchanged.

Prefer:

    existing Solvent corpus schema
    +
    thin ARGUS research adapter

over:

    new Solvent schema
    +
    new vector subsystem

Do not expand the Solvent kernel merely to support retrieval.

==================================================
CORPUS INGESTION
==================================================

Design the smallest bootstrap mechanism that can ingest the seven attached
research documents into the existing corpus.

Requirements:

- deterministic source identity
- SHA-256
- source path/locator
- title
- full text
- corpus membership
- embedding
- idempotent re-ingestion
- source attribution
- no authority elevation

Determine whether the existing corpus manifest package can be reused.

Do not make agents manually upload the documents.

The initial POC should bootstrap them automatically.

==================================================
DOMAIN PACK
==================================================

Audit the current BM-IST Domain Pack against v6.3.

Determine whether the current pack can express the new tripartite program.

The domain pack should own:

- claim types
- evidence classes
- debt vocabulary
- retirement rules
- verifier declarations
- domain methodology

It should NOT own:

- the research corpus
- current research state
- individual claims
- research tasks
- agent history

If the current BM-IST pack is too narrow, propose the smallest successor pack
needed for the combined BM–IST–AS research program.

Do not make core ARGUS/Solvent code physics-specific.

==================================================
INITIAL RESEARCH SEED
==================================================

Determine the minimum state that should be seeded automatically.

Prefer something like:

    one research scenario
    one root research task
    one root candidate belief/question
    initial EBP debt
    attached corpus

rather than pre-populating dozens of claims.

The seed should create room for the agent to do the research rather than
telling the agent what conclusions to reach.

The seed must preserve v6.3's current status:

    research program
    not established physical theory

==================================================
INSIGHTS / DEBTS
==================================================

Verify whether the existing /insights and /debts surfaces already expose
everything required.

Insights should make it possible to see:

- active research work
- candidate claims
- open debt
- supporting evidence
- adversarial challenges
- graph structure
- current research activity

Debts should make it possible to see:

- claim
- debt item
- why it exists
- relevant evidence
- relevant corpus references
- retirement rule
- evidence class
- human discharge action

Do not turn the UI into a scientific authoring system.

The agent should do the research.
The UI should make the epistemic state legible and let humans adjudicate.

==================================================
ADVERSARIAL WORKFLOW
==================================================

Verify that a fresh adversarial agent can reconstruct the complete state
without conversational memory.

Minimum context:

    seed research program
    current task
    candidate claims
    evidence
    open debt
    derives edges
    contradicts edges
    previous activity
    relevant corpus hits

The adversarial agent must be able to attack:

- logical validity
- hidden assumptions
- missing evidence
- debt completeness
- theorem applicability
- category/type errors
- counterexamples
- alternative explanations
- unsupported promotion pressure

It must NOT directly discharge debt or promote claims.

==================================================
DOMAIN-AGNOSTIC TEST
==================================================

The implementation should demonstrate:

    Solvent/ARGUS core
        =
    generic epistemic substrate

    Domain Pack
        =
    domain methodology

    Research corpus
        =
    domain knowledge

    Verifier
        =
    domain computation

A future domain should be able to reuse the same core without changing the
belief/evidence/debt graph model.

Do not require implementation of a second domain.

==================================================
MINIMAL WORKFLOW RECOMMENDATION
==================================================

Evaluate this target loop:

    argus serve
        ↓
    bootstrap research seed
        ↓
    Work Agent
        ↓
    get_context(query)
        ↓
    vector retrieval + epistemic context
        ↓
    research
        ↓
    submit_packet
        ↓
    Insights / Debts
        ↓
    fresh Adversarial Agent
        ↓
    get_context(query)
        ↓
    attack
        ↓
    submit_packet
        ↓
    Human adjudication
        ↓
    discharge / retract / promote

Determine whether this can be achieved without adding new agent tools.

==================================================
OUTPUT: GAP AUDIT + IMPLEMENTATION PLAN
==================================================

Produce:

# ARGUS BM–IST–AS POC GAP AUDIT

## 1. What already works

Map each requested capability to the actual code.

Classify:

    READY
    PARTIAL
    MISSING

Do not infer from old plans if the current code differs.

## 2. Missing capabilities

Only list genuine gaps.

For each:

- current behavior
- required behavior
- why existing code is insufficient
- whether the gap belongs in:
    ARGUS interface
    Domain Pack
    corpus/retrieval adapter
    research seed
    Solvent kernel
- minimal solution

## 3. Solvent kernel impact

Explicitly state:

    No kernel change required

or identify the exact invariant that truly cannot be expressed outside the
kernel.

Be ruthless here.

The preferred answer is zero kernel changes.

## 4. Minimal implementation architecture

Show the smallest architecture needed.

Prefer:

    existing Solvent schema
    +
    thin research/corpus adapter
    +
    existing ARGUS application
    +
    existing MCP tools
    +
    existing Trust UI

Avoid new services.

## 5. Seed design

Define:

- research scenario
- root task
- root belief/question
- domain pack
- initial debt
- corpus manifest
- corpus provenance

## 6. Agent workflow

Specify exactly what the Work Agent and Adversarial Agent do.

Keep both agent roles within:

    get_context
    submit_packet

wherever possible.

## 7. Vector retrieval contract

Define the minimum retrieval contract:

    query
    top_k
    source/document identity
    similarity
    excerpt
    locator
    provenance

Clarify where retrieval appears in `get_context`.

## 8. UI acceptance criteria

Define the smallest observable workflow:

    Work Agent creates research state
        ↓
    /insights shows it
        ↓
    /debts shows obligations
        ↓
    Adversarial Agent reconstructs it
        ↓
    adversarial findings appear
        ↓
    human discharges a debt
        ↓
    promotion gate behaves correctly

## 9. Developer experience

Target exactly:

    argus serve
    → seeded corpus
    → seeded research question
    → OpenCode can start work immediately

No manual SQL.
No manual corpus ingestion.
No manual vector-store setup.
No manual database configuration.
No extra MCP server.

## 10. Acceptance tests

Define the smallest end-to-end tests proving:

1. clean startup
2. corpus bootstrap
3. vector retrieval
4. Work Agent context reconstruction
5. EBP packet submission
6. belief/debt persistence
7. Insights projection
8. Debts projection
9. fresh adversarial reconstruction
10. adversarial challenge persistence
11. human debt discharge
12. promotion gate

## 11. Implementation order

Provide a dependency-ordered implementation sequence.

IMPORTANT:

Do NOT implement anything in this task.

The purpose of this task is to determine whether the current ARGUS/Solvent
foundation is sufficient and, where it is not, produce the smallest
implementation plan that closes only the real gaps.

The governing design rule remains:

    The optimal design is reached when there is nothing left to subtract.
```

### The workflow I would target

The strongest simplification is to make **`get_context` the agent's single research window**:

```text
get_context(
    task_id,
    query?
)
```

It returns both the epistemic graph and the relevant corpus. That avoids adding `search_corpus`, `get_beliefs`, `get_debts`, `get_activity`, etc. The existing design already treats the context response as a full scenario projection rather than a sophisticated graph API. 

That gives you a very clean loop:

```text
Corpus
  ↓
Research context
  ↓
Agent work
  ↓
EBP packet
  ↓
Solvent belief graph
  ↓
Human debt adjudication
  ↓
Adversarial context
  ↓
Attack
```

And importantly, **the corpus remains memory, not authority**. That separation is already present in the Solvent design and is exactly what you want for a physics research program. 

The v6.3 program itself reinforces why this structure is appropriate: it explicitly defines the work as a research program, not established theory, and makes “construct the arrows, or abandon the claim” the governing posture. 

One final design choice I would preserve: **do not seed the database with the conclusions of v6.3. Seed it with the question, corpus, and obligations.** Let the first agent construct the belief graph from the program. That is much closer to EBP 2.1 than turning the research charter into pre-promoted knowledge. EBP explicitly treats ideas as free to enter while promotion requires debt retirement. 



## Z

Assessment first, then the prompt. Verdict on "can Argus + Solvent serve this now": **the core loop — yes, conditionally; three bounded gaps — no.** So the prompt below is hybrid: verify the fixes, plan-and-implement the gaps, then run the dry run. It's one prompt with a plan-first stage, consistent with the project's own Growth Gate discipline.

## Requirement-by-requirement assessment

| Requirement | Served? | Notes |
|---|---|---|
| v6.3 as baseline, full doc set as corpus | Partially | Solvent corpus layer exists (002/003 migrations, `belief_corpus_citation`), but **vector search does not** — new capability, needs additive migration + embedding pipeline + search surface |
| Agent works the central question → tasks/debts/belief lifecycle per EBP 2.1 | Yes | Phase 8 loop covers it; pack needs v6.3 update (two new debt classes: `needInitialCondition`, `needRegularity`) — pack is data, so this is a legitimate update |
| `/insights` populated by agent | Yes, with one clarification | Agents never write Insights; they submit packets → state changes → Insights derives. That *is* the lawful population path |
| Debt discharge via UI | Yes, pending verification | Attributed discharge + rule enforcement were BLOCKED pending F1–F7; "argus serve works" ≠ "fixes verified" — pre-flight required |
| Adversarial agent reconstructs context | Yes — with one correction below | Via RCP, not schema inspection |
| Frictionless DX | No | No one-command bring-up, no seeded first-boot state, no quickstart |

**Two corrections to your framing before the prompt**, because they're doctrine, not nitpicks:

1. **"OpenCode runs Argus, Argus saves the objects"** — inverted. The agent submits an EBP packet to the Coordinator; *the Coordinator* persists. The agent never saves objects. That's the entire authority model.
2. **"The adversarial agent should look up the CockroachDB schema"** — no. Schema inspection is developer work; the agent's interface is RCP (`argus.get_context`). If an agent needs the schema to reconstruct context, RCP has failed and that's a bug to fix — not a workflow to adopt. The schema-lookup instinct is exactly what RCP exists to make unnecessary.

Also one EBP note on the seed: the boxed question is a *question*, and EBP enters claims. The correct seed belief is the governing hypothesis phrased as claim-under-test — and it must pass the final-truth language check ("candidate/hypothesis under test," never "can generate" as fact).

---

## Prompt for the coding agent

```text
TASK: ARGUS/Solvent POC pivot — v6.3 corpus ingestion, vector search, domain pack update, seed dry run

CONTEXT
Argus (Coordinator + Trust UI), Solvent (epistemic authority kernel, frozen), and Conductor
(workflow) implement a packet-based research trust loop: agents submit EBP packets; the
Coordinator validates and persists to Solvent/Conductor; humans adjudicate via the Trust UI;
promotion is gated by Solvent's database constraints; RCP is the read-only agent context
protocol exposed over MCP. A prior code review (BLOCKED verdict, findings F1–F7) was issued
against this loop; you are told "argus serve works" — do NOT assume the fixes landed. Verify.

Do not modify the Solvent kernel. All work is additive (corpus layer, coordinator, UI, pack
data, migrations 011+). Agents get read + packet-submission capability only.

STAGE 0 — PRE-FLIGHT VERIFICATION (blocking)
Verify each fix from code review adv_review.md with a test that exercises the PRODUCTION
HTTP path, not mocks:
  F1: unknown debt item → refused (not executed)
  F2: evidence_class verified against ACTUAL persisted Solvent evidence rows, not caller claim
  F3: /decisions requires operator identity; unauthenticated call fails closed
  F4: EvidenceClass passes through the HTTP layer (SubmitDecisionRequest → DecisionRequest)
  F5: human retirement path calls POST /v1/discharge (attributed: DischargedBy, InstrumentRef),
      not bare /debt/retire
  F6: pack resolution failure → hard error, never fall-through
  F7: belief retraction cascades to terminal-cancelled linked Conductor tasks
If any fails: fix it first, with the test, before proceeding. Record results in a preflight.md.

STAGE 1 — IMPLEMENTATION PLAN (short, before code)
Draft a one-page plan for the three gaps below, per the project's plan conventions
(phases, acceptance criteria, rollback, files-changed). Get it approved before implementing.

GAP A — Corpus ingestion with native CockroachDB vector search:
  - Additive migration (db/011_corpus_vector.sql or corpus-layer equivalent): embedding
    vector column + vector index on corpus artifacts, following the pattern used in the
    prior CockroachDB vector-search demo.
  - Ingestion command/endpoint: chunk + embed the corpus documents; embeddings are derived
    artifacts — record the embedding model identity + version in the corpus manifest
    (provenance discipline).
  - Corpus manifest (hash-pinned, per Phase 5 plan):
      * BM_IST_AS_v6_3.md          → provenance: source_of_record (current baseline)
      * v6.2, v6.1, writeup_v1     → provenance: superseded/history (retained, never deleted)
      * adv_review*.md             → provenance: agent_derived / external review artifacts
    NOTE: pinned is not authoritative. The Qwen artifact lesson applies: AI-generated
    documents are agent_derived, never source evidence.
  - Semantic search read endpoint on the Coordinator: query → top-k corpus passages with
    locators + hashes.
  - MCP surface: add ONE third read-only tool argus.search_corpus (query → passages).
    This changes the frozen two-tool boundary — update the capability table, the
    boundary tests, and record the change as a read-only surface extension. Agents still
    get ZERO mutation tools.

GAP B — Domain Pack v6.3 (pack data, version 1.1.0):
  - debt_vocabulary and initial_debt gain: needInitialCondition, needRegularity
    (8 items total, per v6.3 §2.3).
  - Retirement rules for the two new items:
      needInitialCondition → operator_asserted / preparation_accounting
      needRegularity       → reproducible_artifact / regularity_audit
  - Verify the Debts UI renders the new items dynamically (it should — pack-driven
    rendering was a locked constraint; if it doesn't, that's a bug to fix).

GAP C — Seed packet + initial belief lifecycle (the dry run's starting state):
  Central claim (entered belief, phrased to pass the final-truth check):
    "A discrete arithmetic microscopic dynamical system is a candidate substrate from which
     the structures of quantum theory (Born statistics, linear Schrödinger evolution,
     guidance law) may emerge without being assumed." — status: INCUMBENT HYPOTHESIS
  - initial_debt: all 8 pack items.
  - Supporting claim candidates extracted from v6.3 (entered free, each with debt):
      e.g., arithmetic substrate supports canonical-height definability (§8.3);
             tame-factor/Kronecker route for the phase interface (§9.3);
             dFRG as AS-methodology-on-substrate (§14.2).
  - Edges: derives edges from sub-claims to the central claim.
  - Conductor tasks (per v6.3 §21 Phase 0 + Phase 1 and §28 First 30 Days):
      T1 constraint compilation + pre-registration templates (Phase 0)
      T2 construct smallest F_L: X_L, encoding, F_L, L, placewise action (Deliverable A)
      T3 canonical-height applicability report (Deliverable B)
      T4 Markov/symbolic applicability report on the SAME map (Deliverable C)
      T5 Koopman/Kronecker-spectrum tame-factor test (Deliverable D)
      T6 prime-universality scan + chi_p closure test (Deliverable E)
      T7 non-arithmetic null substrate of comparable complexity (Deliverable F)
    Each task carries governance_ref anchored to its belief (scenario grain).
  - IMPORTANT: agents propose this structure via packets; you may seed it directly via the
    Coordinator path for bootstrapping, but the seeded state must be indistinguishable from
    packet-compiled state (same attribution and audit trail).

STAGE 2 — DEVELOPER EXPERIENCE
  - One-command bring-up: docker compose (CockroachDB + Solvent + Conductor + Coordinator +
    Trust UI), idempotent first-boot seed (corpus ingest + pack load + seed packet), health
    checks. `argus up` (or compose equivalent) must take a new developer from zero to the
    Insights screen showing the seeded v6.3 state.
  - README quickstart: the two-agent walkthrough (work agent prompt, adversarial agent
    prompt, human debt-discharge step) with exact commands.
  - Agent role cards as files: prompts/opencode-work.md, prompts/opencode-adversarial.md —
    each stating: reconstruct context via argus.get_context; treat returned state as
    authoritative; agents never mutate state directly; submit packets only; adversarial
    agents attack unresolved claims and may create contradicts edges via packet.

STAGE 3 — DRY RUN (pre-registered acceptance criteria)
  1. Corpus: all artifacts ingested, hashes verified; semantic query "Kronecker factor
     amplitude" returns the §9.3 passage; query "canonical height preperiodic" returns §8.3.
  2. Seed state: central belief entered with 8 debts; 7 tasks created with governance_refs;
     Insights shows correct counts (Open Claims, Open Debt=8, Active Tasks, Promoted=0).
  3. Work agent (fresh OpenCode, argus MCP only): get_context → produces a packet that
     decomposes the central question into sub-claims + tasks per EBP → Coordinator persists.
  4. Adversarial agent (fresh process, no shared context): get_context reconstructs prior
     state WITHOUT reading the DB schema or repository; attacks an unresolved claim via
     packet (contradicts edge and/or add-debt proposal).
  5. Human: discharges one debt via Debts UI → attributed discharge recorded in Solvent
     audit trail with operator identity; rule-class validation fires (refuse a mismatched
     class attempt deliberately).
  6. Promotion attempted while debt remains → Solvent refuses (23514); UI shows the refusal
     with outstanding items. This is a required demonstration, not an accident.
  7. Capability boundary: agents expose exactly three read/packet tools; zero mutation
     tools; attempted direct Solvent mutation fails (tool absent).
  8. add-debt path: an adversarial contradicts edge leads to a human ADD_DEBT decision that
     reinstates debt on the target belief (the living-ledger loop). If the ADD_DEBT
     decision type is not implemented, implement it now — it was a known gap.
  9. Regression: all existing Solvent/Conductor/Coordinator tests pass; kernel untouched
     (grep sweep: no physics vocabulary in kernel; debt strings only in pack/UI/wizard).

BOUNDARIES (all standing, non-negotiable)
  - No Solvent kernel changes; vector search lives in the corpus layer.
  - Agents: read context, search corpus, submit packets. Never: retire, promote, retract,
    authorize, execute, write Conductor.
  - Human decisions flow UI → Coordinator → attributed Solvent discharge.
  - Solvent-before-Conductor write ordering; packet_id idempotency.
  - Promotion/refusal decided by Solvent's DB gate only.
  - Provenance classes honest: v6.3 is source_of_record; AI reviews are agent_derived.

DELIVERABLES
  preflight.md, implementation plan, migrations, corpus manifest + loader, search endpoint
  + MCP tool, pack v1.1.0, seed script, compose file, README quickstart, role cards,
  test suite (per acceptance criteria), dry-run evidence (traces + screenshots).

Report PASS/FAIL per acceptance criterion with evidence. Do not convert failures into
successes; a clean FAIL with evidence is a valid outcome.
```

## Workflow suggestions (beyond the prompt)

1. **The session rhythm, documented:** work-agent run → human reads Insights → adversarial run → human adjudicates (discharge/add-debt/retract). Write this loop on the quickstart's first page — it's the product's actual usage pattern, and right now it exists only in review threads.
2. **Scheduled red-team = the adversarial agent on a cadence.** v6.3 §2.4 requires external red-team at milestones; your Argus adversarial agent *is* that mechanism for AI-generated claims. Wire a `argus redteam` command later that spawns the adversarial run against all promoted beliefs — the governance section of the research program and the tool then reinforce each other.
3. **Version-pin the pack and corpus in every packet** (`pack_ref` + `corpus_manifest_ref` already exist in the schema) — so any future "which corpus did this agent see?" question is answerable from the packet alone. Cheap now, invaluable after the first re-ingestion.
4. **Resist making search the default context.** RCP's scenario projection is deterministic; vector search is *supplementary* retrieval. If agents start building context primarily from semantic search, reconstruction stops being reproducible — keep search for targeted lookup ("find the passage about X"), not for state reconstruction.

The honest bottom line for your pivot: this is the first seed where the *research program itself* was produced by the discipline the tooling enforces — v6.3 carries its own self-audit, subtraction record, and debt classes, and it enters Argus the same way every future claim will: free, in debt, and falsifiable. The dry run's real question isn't whether the UI renders — it's whether a fresh agent can pick up a 44-page foundational-physics charter and produce a *lawful* decomposition of it. That's the demo worth shipping.


