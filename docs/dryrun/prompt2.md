## prompt

I am reconsidering the initial corpus ingestion by just requiring the AI agent to just read the 7 files as background context in the initial prompt, what do you think?

## answer

Yes. For the **initial POC, I think this is a better design**.

The corpus/vector-ingestion work in `plan2` is technically coherent, but it is solving a retrieval problem before we have demonstrated that retrieval is actually necessary. The current plan adds a corpus schema, embeddings, vector search, bootstrap ingestion, provenance persistence, and `get_context(query)` purely to make seven documents discoverable. 

For the first dry run, I would subtract all of that.

### Initial POC

Keep the seven documents as ordinary repository files:

```text
docs/corpus/
    v6_3.md
    v6_2_adv.md
    v6_1.md
    adv_review.md
    adv_review2.md
    adv_review3.md
    writeup_v1.md
```

Then make the Work Agent's initial instruction explicit:

```text
Before doing research:

1. Call argus.get_context with the assigned task_id.
2. Read all seven documents under docs/corpus/ as background research context.
3. Treat v6_3.md as the current research-program baseline.
4. Treat v6_1.md and v6_2_adv.md as historical/superseded context.
5. Treat adv_review*.md as adversarial/agent-derived material, not authoritative evidence.
6. Treat writeup_v1.md as explanatory context.
7. Do not treat any document as authoritative Solvent state.
8. Reconstruct the actual current state from argus.get_context.
9. Work only from the combination of current ARGUS state and the seven background documents.
```

That gives you something much more important for this experiment: **a clean test of whether an agent can perform useful research inside the EBP loop without needing a retrieval subsystem.**

### The architectural advantage

It makes the boundary extremely clear:

```text
Repository documents
        ↓
   agent background
        ↓
ARGUS get_context
        ↓
current epistemic/work state
        ↓
agent reasoning
        ↓
submit_packet
        ↓
Solvent
```

The seven files are **context**, not a second authority system.

That is actually cleaner than putting them into `corpus_issue`, because otherwise you immediately have to answer questions about corpus identity, chunk identity, embeddings, provenance, retrieval quality, citation semantics, and whether retrieved text can influence authority.

### And this fits the agent model better

The important experiment is not:

> "Can ARGUS perform vector retrieval?"

It is:

> **"Can a fresh research agent reconstruct the state of an ongoing research program and make useful, auditable progress without inheriting conversational memory?"**

You already have `get_context` for the authoritative state and `submit_packet` for controlled mutation. Those are the two important MCP primitives. The seven documents are static intellectual background.

That gives you a much smaller first milestone:

```text
argus serve
    ↓
seed one task + one hypothesis + eight debts
    ↓
fresh Work Agent
    ↓
get_context
    ↓
read seven research documents
    ↓
research
    ↓
submit_packet
    ↓
Insights / Debts
    ↓
fresh Adversarial Agent
    ↓
get_context
    ↓
read seven documents
    ↓
attack prior work
    ↓
submit_packet
```

No embedding model.
No vector index.
No corpus migration.
No corpus bootstrap.
No `get_context(query)`.
No third retrieval pathway.

### One important caveat

Do **not** put all seven documents literally into the system prompt.

Have the prompt instruct the agent to **read the files**. That preserves the distinction between:

```text
prompt = operating instructions
files  = research material
ARGUS  = authoritative system state
```

It also avoids making the prompt enormous and makes document updates trivial.

For the fresh adversarial agent, use the same seven files. That is actually valuable: both agents receive the same external background, while only ARGUS state carries the accumulated research history.

### When should vector search return?

Only after the dry run demonstrates a concrete failure such as:

> "The agent cannot efficiently find relevant material across the research corpus."

Then add retrieval to solve that demonstrated problem.

That gives you a much stronger sequence:

```text
Phase 1
Static background files
        ↓
Prove EBP agent loop

Phase 2
Measure context/retrieval failure
        ↓
Add targeted retrieval

Phase 3
Only then consider embeddings/vector search
```

I would therefore **remove Phases 1–3 of the current corpus implementation from the POC entirely** and make the seven files repository-level background context. The current seven-file seed itself can remain exactly as the research corpus; only the database ingestion machinery disappears. The existing plan's core seed, agent roles, and end-to-end verification remain useful.  

This is much closer to your governing principle:

> **Don't build the retrieval infrastructure until the research workflow proves it needs retrieval.**


## prompt

Below is the revised coding-agent prompt. It removes corpus ingestion, embeddings, vector search, and `get_context(query)` from the initial POC, while preserving the rest of the EBP dry run. It also incorporates the architectural corrections identified from `plan2`: immutable pack versioning, pack-driven debt vocabulary, dedicated `argus mcp`, and executable provenance boundaries. The original plan's audit showed that the genuine gaps are concentrated in MCP transport, seed state, the two missing debt classes, and first-run developer flow.  

````text
# ARGUS BM-IST-AS v6.3 — Minimal POC Implementation Prompt

You are implementing the initial ARGUS research POC for the BM–IST–Asymptotic Safety research program.

The governing principle is:

> The optimal design is reached when there is nothing left to subtract.

Do not expand the architecture. Do not build infrastructure for problems the first dry run has not demonstrated.

The POC must prove one thing:

> Can a fresh AI research agent reconstruct the current research state, read the supplied research corpus as background context, perform useful research, submit an auditable EBP packet, and then be challenged by a completely fresh adversarial agent?

The authoritative state remains ARGUS/Solvent.

The seven research documents are static background context for the agents. They are NOT a second authority system, are NOT ingested into CockroachDB, and are NOT searched through embeddings in this phase.

---

# 1. ARCHITECTURE — FROZEN

Keep this architecture:

OpenCode
   |
   | MCP stdio
   v
ARGUS — one Go binary
   |
   +-- application/orchestration
   +-- Solvent epistemic kernel
   +-- minimal work subsystem
   +-- domain pack
   +-- verifier
   +-- MCP adapter
   +-- HTTP/Trust UI
   |
   v
CockroachDB

Capability ≠ Work ≠ Authority ≠ Execution

Agent / Agent Runtime = AGENCY
ARGUS application     = COORDINATION
Solvent                = AUTHORITY / canonical epistemic state
Executor               = EFFECT
External SOR           = EFFECT TRUTH

Do not introduce:
- another service
- another database
- another persistence layer
- REST between internal components
- a third MCP tool
- a vector database
- a corpus retrieval service
- a background worker architecture
- a new orchestration framework

No Solvent kernel changes unless the audit proves there is a genuine invariant-level deficiency.

---

# 2. FIRST RULE: AUDIT THE CURRENT REPOSITORY

Before changing code, inspect the actual current implementation.

Do not trust this prompt, old plans, or previous assumptions over the repository.

Produce an internal audit first and classify each item:

READY
PARTIAL
MISSING
INCORRECT

Inspect at minimum:

1. domain-pack loading and validation
2. debt vocabulary and retirement rules
3. hardcoded debt/evidence validation
4. MCP adapter
5. MCP stdio transport
6. `argus serve`
7. `argus mcp`
8. database bootstrap
9. migrations
10. task/project creation
11. belief creation
12. `governance_ref`
13. packet compilation
14. packet persistence
15. adversarial packet handling
16. Insights UI
17. Debts UI
18. authenticated human discharge
19. promotion gate
20. the seven research documents
21. role/prompt files
22. existing tests

Do not implement something merely because this prompt says to implement it if the current repository already satisfies it.

Do not remove working functionality merely to make the implementation look like this prompt.

---

# 3. MAJOR SCOPE DECISION: NO DATABASE CORPUS IN THIS POC

The seven research documents are intentionally kept as ordinary repository files.

They must NOT be:

- inserted into CockroachDB
- chunked into database records
- embedded
- vector-indexed
- retrieved through cosine similarity
- added to `get_context`
- exposed through a new MCP search tool
- treated as authoritative evidence
- treated as Solvent ledger state

Therefore DO NOT implement:

- `corpus_issue`
- corpus vector migrations
- `VECTOR(...)` columns
- vector indexes
- hash-based embedding
- Bedrock embedding calls
- `corpus.Search`
- `corpus.BootstrapCorpus`
- `get_context(query)`
- `argus.search_corpus`
- corpus retrieval APIs
- corpus citation persistence for these seven files

If there is already a reusable corpus implementation elsewhere in the repository, do not delete it blindly. But do not wire it into the Oracle/ARGUS POC unless the existing implementation is already part of the working architecture.

The first dry run must prove whether retrieval infrastructure is actually necessary.

---

# 4. THE SEVEN FILES ARE STATIC BACKGROUND CONTEXT

The seven files must remain available in the repository under a stable location such as:

docs/corpus/

Expected files:

1. `v6_3.md`
2. `v6_2_adv.md`
3. `v6_1.md`
4. `adv_review.md`
5. `adv_review2.md`
6. `adv_review3.md`
7. `writeup_v1.md`

Before implementation, verify the actual paths and filenames in the repository. Do not invent filenames if they differ.

The agent must be explicitly instructed to READ these files as background research context.

Important distinction:

ARGUS/Solvent state
    = authoritative current state

Seven documents
    = background research material

Agent reasoning
    = candidate work

Submitted packet
    = controlled state mutation subject to validation

The agent must never treat the documents as authoritative merely because they are in the repository.

---

# 5. DOCUMENT ROLE CLASSIFICATION

The role prompt must establish:

`v6_3.md`
- current BM–IST–AS research-program baseline
- primary background reference

`v6_2_adv.md`
- superseded historical baseline
- useful for understanding prior revisions and objections

`v6_1.md`
- earlier historical research-program baseline

`adv_review.md`
`adv_review2.md`
`adv_review3.md`
- adversarial / AI-derived review material
- potentially useful criticism
- not authoritative evidence
- not equivalent to human adjudication

`writeup_v1.md`
- explanatory/supporting research-program material

Do not silently upgrade any of these classifications into epistemic authority.

---

# 6. DOMAIN PACK

The current research methodology needs eight debt classes:

- needMap
- needInvariant
- needToyCheck
- needNullModel
- needObstruction
- needFaithfulnessReview
- needInitialCondition
- needRegularity

Do not leave validation hardcoded to the old six-item list.

Replace hardcoded debt validation with domain-pack registry lookup.

Conceptually:

```go
func (a *App) isValidDebtItem(item string, packRef string) bool
````

must resolve the referenced pack and validate the item against the pack's actual debt vocabulary.

Do the same for evidence-class validation wherever current code hardcodes methodology-specific values.

The seed must obtain initial debt from the loaded pack rather than duplicating the eight names in application code.

For example:

```go
initialDebt := pack.GetInitialDebt()
```

The domain pack is the source of truth for methodology.

---

# 7. PACK VERSIONING — IMMUTABLE SEMANTICS

Do NOT silently mutate the semantics of an already-referenced pack version.

Audit the current identity first.

Because the methodology now explicitly covers:

BM + IST + Asymptotic Safety

use a new immutable pack identity/version if the current repository does not already have an appropriate one.

Preferred:

`bmist-as@0.1.0`

An acceptable alternative is an incremented immutable version such as:

`bmist@1.1.0`

Do NOT change the meaning of an existing `bmist@1.0.0` in place merely to avoid updating references.

Past epistemic records must remain reproducible.

Do not create cascading compatibility work unless required by the current code.

---

# 8. MCP — EXACTLY TWO TOOLS

Preserve exactly:

1. `argus.get_context`
2. `argus.submit_packet`

Do not add:

* `argus.search_corpus`
* `argus.read_document`
* `argus.add_debt`
* any other MCP tool

The two existing tools are sufficient for this POC.

`get_context` remains responsible for deterministic reconstruction of ARGUS state.

Do NOT add a `query` argument for corpus retrieval in this phase.

---

# 9. MCP STDIO TRANSPORT

This is a genuine current gap and must be implemented.

Keep:

```text
argus serve
argus mcp
argus verify
argus migrate
argus reset
```

`argus mcp` should be the dedicated stdio MCP process.

Do NOT merge HTTP server stdout and MCP protocol traffic.

The architecture is:

same binary
same codebase
same database
same Solvent kernel
different process modes

Implement the minimum MCP stdio transport necessary for OpenCode.

Support at minimum:

1. `initialize`
2. `tools/list`
3. `tools/call`

Use the existing MCP adapter rather than duplicating application logic.

Protocol output must go to stdout.

Logs/debugging must NOT corrupt stdout.

Prefer standard library JSON/stdio handling unless the repository already has an MCP SDK that is actually functioning and appropriate.

Do not add an MCP dependency merely for convenience.

Acceptance:

```text
argus mcp
```

can communicate with OpenCode over stdin/stdout and successfully call both existing tools.

---

# 10. FIRST-RUN SEED

A fresh `argus serve` must no longer produce an empty research environment.

Seed only the minimum needed to begin the experiment.

Create:

1. one project
2. one root research task
3. one incumbent hypothesis/belief
4. eight initial debts

Do NOT create a prewritten belief graph.

Do NOT create dozens of tasks or beliefs.

The first Work Agent must earn the decomposition.

Root research question:

"Can one explicit microscopic dynamical system generate the structures of quantum theory without assuming them?"

Current framing:

This is a research program, not an established physical theory.

The incumbent hypothesis should remain explicitly hypothetical, for example:

"A discrete arithmetic microscopic dynamical system is a candidate substrate from which structures associated with quantum theory may emerge without being assumed."

Do not use final-truth language.

The belief should carry the pack-derived initial debt list.

The root task should be connected to the belief through the existing governance mechanism.

---

# 11. GOVERNANCE_REF

Before implementing seed logic, inspect the actual current schema and application contract for `governance_ref`.

Do not assume the representation from an earlier plan is correct.

Resolve the actual contract between:

* conductor task
* application layer
* database schema
* belief
* governance reference

Use one canonical representation.

Do not introduce compatibility hacks unless the existing architecture genuinely requires them.

---

# 12. DATABASE BOOTSTRAP

Preserve the existing zero-friction developer experience:

```text
argus serve
```

must be the canonical first-run entrypoint.

It should:

1. start/bootstrap local CockroachDB as currently supported
2. create the database if necessary
3. apply migrations
4. ensure the required domain pack is available
5. create the minimal research seed if the target environment is empty
6. start the Trust UI

Do not introduce Docker Compose.

Do not introduce another database.

Do not introduce manual setup steps if the existing bootstrap architecture can perform them.

Seeding must be idempotent.

Second boot must not duplicate the project/task/belief.

Use existing atomic/bootstrap mechanisms where available.

---

# 13. ADD_DEBT REMAINS DEFERRED

Do not implement ADD_DEBT.

The audit must first demonstrate that current packet submission cannot represent new obligations.

For this POC, new claims and derived beliefs can carry their own applicable debt.

The objective is to demonstrate the living EBP loop, not to expand the protocol.

Only revisit ADD_DEBT if the actual dry run proves the current representation cannot express a necessary research obligation.

---

# 14. ROLE CARDS

Create short reusable prompts:

`prompts/opencode-work.md`

`prompts/opencode-adversarial.md`

Do not make them enormous.

## Work Agent prompt

The Work Agent must be instructed to:

1. call `argus.get_context` first
2. reconstruct current ARGUS state
3. read all seven research files as background context
4. distinguish current ARGUS state from background documents
5. treat v6_3 as current baseline
6. treat historical/adversarial documents according to their stated role
7. identify the actual research frontier
8. inspect unresolved debt
9. decompose the central question into manageable candidate claims
10. gather evidence
11. make logic chains explicit
12. create derives relationships where appropriate
13. submit an EBP packet
14. never mutate protected state directly

The agent MUST NOT:

* write to the database
* inspect/reconstruct the research state by querying database tables directly
* retire debt
* promote claims
* retract beliefs
* authorize actions
* treat its own conclusions as authoritative
* treat the seven documents as authoritative Solvent evidence

The role card should preserve EBP v2.1 principles:

Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
Promotion never means truth.

---

# 15. ADVERSARIAL AGENT

The adversarial agent must be a genuinely fresh process.

It has:

* no conversational memory
* no prior agent transcript
* no database/schema inspection privileges
* the same two MCP tools
* access to the same seven background documents

Its first action must be:

```text
argus.get_context(task_id)
```

Then it reads the seven files as background context.

It must independently reconstruct the current state.

It must challenge:

* logical validity
* hidden assumptions
* evidence completeness
* debt completeness
* theorem applicability
* category errors
* counterexamples
* alternative explanations
* unsupported leaps
* mismatch between formalization and intended claim

It must never assume that prior work is correct merely because it appears in ARGUS.

It submits an adversarial EBP packet through the existing packet path.

It does NOT:

* retire debt
* promote claims
* retract beliefs
* directly edit the DB
* inspect the schema to reconstruct context
* declare its own challenge authoritative

---

# 16. PROVENANCE BOUNDARY

The seven static research files are background context.

The distinction must remain explicit:

```text
background document
    ≠
authoritative evidence
    ≠
promoted belief
    ≠
human adjudication
```

In particular, the `adv_review*.md` files are AI-derived background material.

An agent may use them to discover objections or generate hypotheses.

But reading an AI-generated review must never by itself satisfy a debt retirement rule.

Add a test demonstrating:

1. an agent can read an adversarial review file
2. the agent can cite or reason from it
3. the material cannot silently become authoritative evidence
4. human-gated retirement still requires the existing retirement rule

Do not create a special corpus subsystem merely to enforce this.

---

# 17. TRUST UI

Do not redesign the existing UI.

Use the existing:

`/ui/insights`

`/ui/debts`

Insights should be derived from actual persisted state.

The agent must never write directly into Insights.

Debts must continue to show:

* claim
* debt item
* why the debt exists
* evidence/state relevant to it
* retirement rule
* required evidence class
* human discharge controls

Do not add dashboards merely for the sake of the research demo.

---

# 18. END-TO-END DRY RUN

The POC is complete only when the following sequence works.

## Step 1 — Clean startup

From a clean local environment:

```text
argus serve
```

starts successfully.

No manual database preparation.

## Step 2 — Seed appears

Trust UI shows:

* one research project
* one active root task
* one incumbent hypothesis
* eight open debts

## Step 3 — Work Agent

Start a completely fresh OpenCode process.

The agent:

1. calls `argus.get_context`
2. receives authoritative current state
3. reads the seven background files
4. reconstructs the research situation
5. identifies the actual frontier
6. performs useful research
7. submits an EBP packet

## Step 4 — State projection

The resulting claims/evidence/edges/tasks appear correctly in:

* Insights
* Debts

## Step 5 — Fresh adversarial agent

Start another completely fresh OpenCode process with no history.

It:

1. calls `argus.get_context`
2. reads the same seven background files
3. reconstructs the state
4. independently attacks the current work
5. submits an adversarial packet

## Step 6 — Human adjudication

Human reviews the resulting debt/challenge state.

Human discharge remains authenticated and subject to the existing:

* pack retirement rule
* evidence-class rule
* Solvent authority boundary

## Step 7 — Promotion gate

Attempt promotion while debt remains.

Solvent must refuse promotion.

Do not weaken the existing promotion gate.

---

# 19. ACCEPTANCE CRITERIA

All must pass:

AC1
`argus serve` starts from clean local state with no manual setup.

AC2
Fresh startup creates exactly the minimal research seed.

AC3
Eight debt classes are accepted by the loaded domain pack.

AC4
Debt validation comes from the pack registry rather than a hardcoded six-item list.

AC5
`argus mcp` communicates over stdio.

AC6
OpenCode can call `argus.get_context`.

AC7
OpenCode can call `argus.submit_packet`.

AC8
A fresh Work Agent can reconstruct current research state through `get_context`.

AC9
The Work Agent can read all seven research documents directly from the repository.

AC10
No corpus tables, embeddings, vector indexes, or retrieval service are required for the dry run.

AC11
No `query` parameter is added to `get_context`.

AC12
No third MCP tool exists.

AC13
The seven documents are never treated as authoritative Solvent state.

AC14
AI-derived adversarial documents cannot independently discharge debt.

AC15
Work Agent packets persist correctly.

AC16
Adversarial packets persist correctly.

AC17
Contradictions/challenges become visible through existing state/UI projections.

AC18
Human discharge remains authenticated and rule-checked.

AC19
Promotion with unresolved debt remains blocked.

AC20
Existing foundation tests remain green.

AC21
No physics-specific vocabulary leaks into generic core packages.

AC22
No new persistence layer is introduced.

AC23
No new service is introduced.

AC24
No Solvent kernel change is introduced unless a genuine invariant gap is demonstrated.

---

# 20. WHAT NOT TO BUILD

Do NOT build:

* corpus ingestion
* corpus migrations
* vector search
* embeddings
* hash embeddings
* Bedrock embedding integration
* vector indexes
* retrieval ranking
* corpus API
* `argus.search_corpus`
* `get_context(query)`
* ADD_DEBT
* Docker Compose
* new services
* new DBs
* new MCP tools
* a prebuilt research graph
* a large workflow engine
* a new UI
* a second authority mechanism

Do NOT solve hypothetical future scaling problems.

---

# 21. IMPLEMENTATION ORDER

Use this dependency order:

Phase 0
Audit repository and identify READY / PARTIAL / MISSING / INCORRECT.

Phase 1
Fix/upgrade domain pack and remove hardcoded six-item assumptions.

Phase 2
Fix MCP stdio transport.

Phase 3
Implement minimal idempotent research seed.

Phase 4
Create/update short Work Agent and Adversarial Agent role cards.

Phase 5
Run clean end-to-end dry run.

Phase 6
Fix only failures revealed by the dry run.

Do not implement future retrieval infrastructure during this phase.

---

# 22. FINAL IMPLEMENTATION RULE

Before adding any component, ask:

> What concrete failure in the current dry run does this solve?

If the answer is "none yet", do not add it.

The first experiment is intentionally:

```text
static research files
        +
authoritative ARGUS context
        +
two MCP tools
        +
EBP packet submission
        +
human adjudication
        +
adversarial replay
```

That is the POC.

Only after this works should retrieval infrastructure be reconsidered.

At the end, report:

1. audit findings
2. files changed
3. tests run
4. dry-run result
5. any remaining genuine blocker

Do not broaden scope without evidence from the implementation or dry run.

```

The key change is that the **research corpus remains in the repository but disappears from the persistence architecture**. That keeps the first experiment focused on the thing that matters: whether the ARGUS/EBP loop actually works with fresh agents.
```
