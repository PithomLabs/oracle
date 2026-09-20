TASK: ARGUS / Solvent — BM + IST + Asymptotic Safety v6.3 research POC seed

This is the next POC slice.

The architecture is already frozen. Do not redesign it.

==================================================
FROZEN ARCHITECTURE
==================================================

    OpenCode
       |
       | MCP stdio
       v
    ARGUS — one Go process
       |
       v
    CockroachDB

ARGUS contains:

- application/orchestration
- Solvent epistemic subsystem
- domain-pack
- verifier
- Trust UI
- MCP adapter
- local database bootstrap

Exactly two MCP tools remain:

    argus.get_context
    argus.submit_packet

Do NOT add a third MCP tool unless the current implementation proves that
the required research workflow cannot be supported through the existing
two-tool surface.

Preferred solution:

    extend get_context with an optional focused corpus query

rather than adding `argus.search_corpus`.

Core invariant:

    CAPABILITY != WORK != AUTHORITY != EXECUTION

Agents:

    may read context
    may perform research
    may submit EBP packets

Agents may NOT:

    retire debt
    promote claims
    retract beliefs
    authorize
    execute consequences
    write Solvent directly
    write workflow state directly

Human authority remains:

    Trust UI
      -> authenticated application path
      -> Solvent authority boundary

==================================================
RESEARCH CORPUS
==================================================

These files are the initial research corpus under docs/corpus folder

    v6_1.md
    v6_2_adv.md
    v6_3.md
    adv_review.md
    adv_review2.md
    adv_review3.md
    writeup_v1.md

Interpretation:

    v6_3.md
        current master research-program baseline

    v6_2_adv.md
    v6_1.md
        superseded research-program history

    adv_review*.md
        adversarial / agent-derived research material

    writeup_v1.md
        explanatory research material

Do NOT flatten these into one document.

Preserve document identity and provenance.

The current baseline remains explicitly:

    research program
    not established physical theory

The central research question is:

    Can one explicit microscopic dynamical system generate
    the structures of quantum theory without assuming them?

The phrase "one flow" MUST NOT be treated as an established fact.

The current status is:

    one candidate flow architecture
    whose existence is itself a central theorem target.

==================================================
EBP v2.1
==================================================

Use:

    Ideas enter free.
    Promotion costs debt.
    Debt does not kill.
    Debt is forever payable.
    New evidence creates new debt.
    Promotion never means truth.

Current debt vocabulary from v6.3:

    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview
    needInitialCondition
    needRegularity

Do not silently invent additional debt classes.

The initial research seed must NOT pre-promote conclusions.

The seed should establish a governing hypothesis under test.

A suitable initial claim is:

    "A discrete arithmetic microscopic dynamical system is a
     candidate substrate from which structures associated with
     quantum theory may emerge without being assumed."

This is a hypothesis under test, not an asserted physical fact.

==================================================
FIRST: AUDIT CURRENT IMPLEMENTATION
==================================================

Before changing code, inspect the actual current repository.

Do not trust old implementation plans.

Determine whether each capability already works:

1. corpus storage
2. corpus ingestion
3. embedding generation
4. CockroachDB native vector search
5. corpus provenance
6. corpus citations
7. get_context
8. submit_packet
9. Insights
10. Debts
11. retirement-rule enforcement
12. human discharge
13. promotion gate
14. adversarial packet handling
15. domain-pack loading
16. first-run database bootstrap
17. embedded BM-IST pack
18. current developer startup

For each, classify:

    READY
    PARTIAL
    MISSING

Do not implement something that already works.

==================================================
CORPUS + VECTOR SEARCH
==================================================

Reuse the existing Solvent corpus layer whenever possible.

The corpus is research memory.

It is NOT authoritative epistemic state.

The system must preserve:

    corpus != evidence != belief != authority

Audit the existing CockroachDB vector implementation.

If native vector search already exists:

    reuse it.

Do NOT create:

    vector database
    external search service
    separate retrieval store

If vector search is incomplete, add only the smallest corpus-layer changes
needed to support:

    document ingestion
    embedding storage
    similarity search
    provenance
    locator/hash preservation

The authoritative Solvent ledger must remain meaningful without vectors.

Vector retrieval must never decide:

    promotion
    debt discharge
    truth
    authority

Vector search is advisory context retrieval only.

==================================================
CORPUS INGESTION
==================================================

The seven files must be available as the initial corpus automatically.

Requirements:

- deterministic identity
- SHA-256
- document identity
- provenance
- embedding
- idempotent re-ingestion
- locator
- source document/version
- embedding model identity and version if applicable

Use the existing corpus manifest mechanism if present.

Do not require a developer to manually upload or ingest the files.

On a fresh local ARGUS database:

    argus serve

should make the corpus available automatically.

On subsequent starts:

    do not duplicate the corpus.

Important provenance rule:

    v6_3.md
        = source_of_record / current baseline

    v6_1.md / v6_2_adv.md
        = historical / superseded

    adversarial reviews
        = agent_derived or external-review material

    writeup_v1.md
        = supporting research material

AI-generated material must never become authoritative merely by ingestion.

==================================================
EMBEDDING / CREDENTIAL AUDIT
==================================================

Audit the actual embedding implementation.

Determine:

- provider
- required credential
- model
- vector dimension
- current demo setup
- local defaults
- failure behavior

For the POC, optimize for zero-friction startup.

Prefer an existing working embedding path over introducing anything new.

Do NOT fabricate semantic embeddings merely to avoid credentials.

If the current demo already has a sane embedding bootstrap, reuse it.

If a real external credential is genuinely unavoidable, make that explicit
as the sole exceptional prerequisite and produce one clear startup message.

Do not create a credential-management subsystem.

==================================================
DOMAIN PACK
==================================================

Audit the current BM-IST domain pack against v6.3.

v6.3 requires the two additional debt classes:

    needInitialCondition
    needRegularity

Determine the cleanest minimal pack representation.

Do not expand core ARGUS/Solvent code with physics-specific vocabulary.

The pack should own:

- claim types
- debt vocabulary
- initial debt
- evidence classes
- retirement rules
- verifier declarations
- domain methodology

The pack must NOT own:

- research corpus
- current beliefs
- tasks
- agent history
- authoritative research conclusions

Because the research scope is now explicitly BM + IST + AS, do not blindly
pretend this is unchanged BM-IST semantics.

Choose the smallest honest pack identity/versioning scheme.

Preferred:

    dedicated tripartite pack identity

if the semantic change is large enough to constitute a different domain
methodology.

Otherwise use a compatible version increment.

Explain the choice in the audit.

==================================================
INITIAL RESEARCH SEED
==================================================

On an empty local database, seed only the minimum state necessary to begin.

Seed:

    one research scenario
    one root research task
    one incumbent hypothesis
    eight initial debts
    the seven-document corpus

Do NOT seed a huge pre-written belief graph.

The purpose is to let the first Work Agent construct the research graph.

The root task should point to:

    "Can one explicit microscopic dynamical system generate
     the structures of quantum theory without assuming them?"

The initial state should make clear:

    current baseline = v6.3
    status = research program
    no claim is promoted
    no claim is final truth

==================================================
WORK AGENT
==================================================

A fresh OpenCode Work Agent should:

1. call argus.get_context
2. receive the current research state
3. optionally retrieve focused corpus passages
4. identify the actual current research frontier
5. decompose the central question into manageable claims/tasks
6. gather evidence
7. form logic chains
8. create EBP beliefs
9. attach appropriate debt
10. create derives edges
11. submit the resulting EBP packet

The system, not the agent, persists the resulting objects.

Agent output must become:

    packet
      -> application validation
      -> Solvent / workflow persistence
      -> derived Insights / Debts state

==================================================
GET_CONTEXT
==================================================

Prefer keeping `argus.get_context` as the single research-context tool.

It should reconstruct state from authoritative ARGUS/Solvent data.

Minimum returned state:

    task
    beliefs
    evidence
    debt
    edges
    activity
    availability

Add corpus retrieval only as supplementary context.

Preferred shape:

    get_context(
        task_id,
        query optional
    )

When query is supplied, return:

    top-k corpus passages
    document identity
    locator
    similarity
    content hash
    provenance

Do NOT let semantic retrieval replace deterministic epistemic
reconstruction.

The agent must be able to reconstruct research state even with no
corpus query.

==================================================
INSIGHTS
==================================================

Agents do not write Insights directly.

Insights is a projection of authoritative/application state.

Verify that the fresh Work Agent's packet submission causes Insights to
show:

- active work
- current claims
- open debt
- evidence
- adversarial challenges
- graph relationships
- activity

Do not turn Insights into an AI chat interface.

==================================================
DEBTS
==================================================

The Debts surface must show:

- belief/claim
- debt item
- why it exists
- relevant evidence
- relevant corpus material where applicable
- retirement rule
- evidence class
- human discharge action

Human discharge remains the authority event.

Verify the complete chain:

    human
      -> authenticated Trust UI
      -> application
      -> pack retirement rule
      -> qualifying persisted evidence
      -> Solvent discharge
      -> audit attribution

Do NOT allow agents to discharge debt.

==================================================
ADVERSARIAL AGENT
==================================================

The adversarial agent is a FRESH OpenCode process.

It must NOT rely on:

    conversation history
    previous-agent memory
    direct CockroachDB access
    repository schema inspection

The agent reconstructs context through:

    argus.get_context

That is the whole purpose of RCP/context.

If the adversarial agent cannot reconstruct the prior research state through
ARGUS, that is a system defect to fix.

The adversarial agent should attack:

- logical validity
- hidden assumptions
- missing evidence
- debt completeness
- theorem applicability
- category errors
- counterexamples
- alternative explanations
- unsupported leaps
- whether the claim actually says what the formalization says

It submits an adversarial EBP packet.

It must not promote, retract, or discharge.

==================================================
ADD-DEBT
==================================================

Do NOT automatically implement an ADD_DEBT subsystem merely to satisfy this
dry run.

First determine whether the current EBP/ARGUS implementation can represent
new obligations without it.

For this POC, it is acceptable for newly identified obligations to be
recorded as new entered claims/tasks or otherwise remain explicitly open.

Only introduce a dedicated ADD_DEBT transition if the audit proves that the
current workflow cannot express the required research semantics.

==================================================
DEVELOPER EXPERIENCE
==================================================

`argus serve` is the canonical local entrypoint.

The target experience is:

    argus serve
        ↓
    local CRDB bootstrapped if needed
        ↓
    migrations
        ↓
    embedded domain pack
        ↓
    corpus available
        ↓
    initial research seed available on first run
        ↓
    Trust UI ready
        ↓
    OpenCode can begin work immediately

Do NOT replace this with:

    docker compose
    multiple services
    manual Solvent startup
    manual Conductor startup

The POC is already intentionally a single process.

Seed only when the local database is genuinely uninitialized.

Subsequent `argus serve` calls must remain idempotent.

==================================================
ROLE CARDS
==================================================

Create minimal reusable prompts if useful:

    prompts/opencode-work.md
    prompts/opencode-adversarial.md

Work Agent:

    reconstruct context
    research
    submit packet

Adversarial Agent:

    reconstruct context
    attack unresolved claims
    submit adversarial packet

Both must state:

    do not mutate system directly
    do not declare authority
    submit packets only

Keep these prompts short.

==================================================
END-TO-END DRY RUN
==================================================

The real POC acceptance path is:

    1. argus serve

    2. Fresh Work Agent
       -> get_context
       -> corpus-assisted research
       -> submit_packet

    3. Inspect /insights
       -> claims
       -> work
       -> evidence
       -> debt

    4. Inspect /debts
       -> obligations
       -> evidence
       -> retirement rules

    5. Fresh Adversarial Agent
       -> get_context
       -> reconstruct previous work
       -> attack claim/debt/evidence
       -> submit_packet

    6. Inspect /insights
       -> contradiction/challenge visible

    7. Human discharges one legitimate debt

    8. Attempt mismatched debt discharge
       -> refused

    9. Attempt promotion while other debt remains
       -> refused by Solvent

This is the core demonstration.

==================================================
ACCEPTANCE CRITERIA
==================================================

AC1
`argus serve` starts from a clean local state with no manual setup.

AC2
The v6.3 corpus is available automatically and ingestion is idempotent.

AC3
Native CockroachDB vector retrieval returns relevant v6.3 passages.

AC4
A fresh Work Agent can reconstruct the seeded research context.

AC5
The Work Agent can submit an EBP packet containing claims, evidence,
edges, and debt.

AC6
The resulting state appears correctly in Insights and Debts.

AC7
A fresh Adversarial Agent can reconstruct the same state without DB/schema
access.

AC8
The adversarial packet can introduce a contradiction/challenge.

AC9
Human debt discharge is authenticated, pack-checked, evidence-checked,
and attributed.

AC10
Promotion with remaining debt is refused by the authoritative gate.

AC11
Existing foundation tests remain green.

AC12
No physics-specific vocabulary leaks into generic core packages.

AC13
No new persistence layer is introduced.

AC14
No third MCP tool is introduced unless the audit proves the existing two
tools are insufficient.

AC15
No Solvent kernel changes are made unless a genuine missing invariant is
demonstrated.

==================================================
IMPLEMENTATION DECISION
==================================================

After the audit, choose ONE:

A. READY

If the current implementation already supports the requirements:
perform the actual dry run and produce evidence.

B. MINIMAL IMPLEMENTATION PLAN

If genuine gaps remain:
produce a concise dependency-ordered implementation plan first.

Do not implement until the plan is complete.

The plan must identify:

- files changed
- schema changes, if genuinely required
- interfaces changed
- acceptance tests
- developer-experience impact
- whether Solvent kernel remains unchanged

The preferred result is:

    zero Solvent kernel changes
    minimal corpus-layer/interface changes
    no new services
    no new database
    no unnecessary MCP tools

==================================================
FINAL DESIGN PRINCIPLE
==================================================

Be ruthless.

The system is successful when a researcher can:

    argus serve
        ↓
    ask the question
        ↓
    research
        ↓
    submit
        ↓
    see the epistemic state
        ↓
    hand it to a fresh adversarial agent
        ↓
    attack it
        ↓
    adjudicate debt

with as few concepts, tools, screens, and moving parts as possible.

The optimal design is reached when there is nothing left to subtract.