# ARGUS Documentation + Repository Cleanup — Authoritative Current State

The ARGUS implementation has now reached a working POC state.

The current source code and current behavior of the `argus` binary are the
authoritative reference for this repository.

The existing AGENTS.md and README.md contain substantial historical material
from earlier architectural iterations and must now be reconciled with the
actual implementation.

Your job is documentation correction + repository cleanup only.

Do NOT add new product features.

Do NOT change the ARGUS architecture.

Do NOT add retrieval/vector infrastructure.

Do NOT add new MCP tools.

Do NOT add agent prompts; the existing Work Agent and Adversarial Agent
prompts are already present.

The objective is:

> Make this repository tell the truth about the ARGUS that actually exists
> today, and remove obsolete artifacts that could mislead a human developer
> or AI coding agent.

==================================================
1. SOURCE OF TRUTH
==================================================

Treat the current repository implementation as authoritative.

In order of authority:

1. actual compiled/runnable ARGUS code
2. actual tests
3. actual migrations/schema
4. actual embedded UI/templates
5. existing working commands/runtime behavior
6. documentation

Do NOT preserve an old architectural statement merely because it appears in
AGENTS.md, README.md, docs/, plan/, or another historical artifact.

Before editing documentation, inspect the current implementation.

At minimum inspect:

- cmd/argus/
- internal/application/
- internal/ui/
- internal/mcp/
- internal/work/
- internal/epistemic/
- domain-pack/
- seed/
- migrations / solventmigrations
- existing agent-facing adapter code
- existing tests
- existing prompts/opencode-work.md
- existing prompts/opencode-adversarial.md

Also inspect the actual available CLI:

    argus serve
    argus mcp
    argus verify
    argus migrate
    argus reset

Document only commands that actually exist.

==================================================
2. CURRENT ARCHITECTURAL MODEL
==================================================

The documentation must reflect the current architecture, not the earlier
multi-service design.

Current conceptual architecture:

OpenCode / other AI agent harness
            |
            | MCP stdio
            v
        argus mcp
            |
            v
     ARGUS application
       |     |      |
       |     |      +-- domain pack
       |     +--------- Solvent epistemic kernel/subsystem
       +--------------- work / coordination
            |
            v
      CockroachDB

`argus serve` provides the integrated local runtime and Trust UI.

Do NOT describe ARGUS as requiring separately launched Solvent, Conductor,
Coordinator, or MCP services unless the current source code actually proves
that requirement.

Do NOT describe internal HTTP/REST service boundaries that no longer exist.

Do NOT resurrect the previous multi-process/multi-service architecture in
documentation.

Preserve the conceptual authority separation:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Conceptually:

Agent / Agent Runtime = agency
ARGUS / application     = coordination
Solvent                = epistemic authority
Trust UI               = human observation/adjudication
External systems        = effect truth where applicable

The conceptual separation remains important even though the current POC is
implemented as one ARGUS process/binary architecture.

==================================================
3. AGENTS.md — REWRITE AS AI-AGENT OPERATING CONTRACT
==================================================

AGENTS.md is primarily for AI coding agents working on this repository.

It should be concise enough to be useful, but authoritative enough that a
fresh OpenCode/Hermes/other coding agent does not accidentally resurrect
obsolete architecture.

Rewrite it around:

A. What ARGUS is

B. Non-negotiable architecture invariants

C. Current repository structure

D. How to build/test/run ARGUS

E. How AI agents interface with ARGUS

F. Capability boundaries

G. Where authority lives

H. What must not be changed casually

I. Current dry-run workflow

J. Current deferred scope

The most important new section is:

# AI AGENT INTERFACE

Document the actual interfaces by inspecting the code.

Primary interface:

    argus mcp

over MCP stdio.

The agent-facing MCP surface is exactly:

    argus.get_context
    argus.submit_packet

Document for each:

- purpose
- input schema
- output shape
- whether read/write
- what state it can affect
- what it cannot do

Make the capability boundary explicit:

Agents do NOT receive:

- direct database access
- Solvent privileged mutation tools
- Conductor privileged mutation tools
- promotion tools
- discharge tools
- retraction tools
- direct edge mutation tools
- authority-changing tools

The intended path is:

    AI Agent
      |
      v
    MCP stdio
      |
      v
    ARGUS MCP adapter
      |
      v
    ARGUS application
      |
      +--> work state
      +--> Solvent authority
      +--> Conductor/work state where applicable

Do not describe this from memory. Verify the exact implementation.

==================================================
4. DOCUMENT OTHER REAL INTERFACES
==================================================

AI agents may use more than MCP in future or through development tooling.

Inspect the current repository and document every REAL externally usable
interface that exists today, including where applicable:

- MCP stdio
- CLI subcommands
- HTTP endpoints
- JSON APIs
- local files intentionally exposed to agents
- verification commands

But do not invent APIs.

For every documented API endpoint, verify it exists in the current source.

Clearly distinguish:

Agent-facing interface
Human-facing interface
Developer/operator interface
Internal implementation

The priority for AI-agent documentation is:

1. MCP
2. repository filesystem/context
3. CLI/developer tooling where relevant
4. any actual machine-readable API

==================================================
5. DOCUMENT THE ACTUAL AGENT WORKFLOW
==================================================

The canonical agent workflow should be:

    Fresh Work Agent
        |
        +--> read current task instructions
        |
        +--> call argus.get_context
        |
        +--> read docs/corpus/*.md as background research
        |
        +--> perform bounded research
        |
        +--> construct EBP packet
        |
        +--> argus.submit_packet
        |
        v
    persisted ARGUS state

Then:

    Fresh Adversarial Agent
        |
        +--> call argus.get_context
        |
        +--> read same research background
        |
        +--> independently reconstruct state
        |
        +--> attack claims/evidence/debt/logic
        |
        +--> submit adversarial EBP packet

Explicitly state:

Agents reconstruct research state from ARGUS state, not inherited
conversation memory.

The existing Work Agent and Adversarial Agent prompt files already exist.

Do not create replacement prompts.

==================================================
6. STATIC RESEARCH CORPUS
==================================================

The seven research documents under docs/corpus/ are static repository
background material for AI agents.

They are NOT:

- a database corpus
- vector search
- embeddings
- an authority store
- a second epistemic ledger
- an MCP search service

Document this distinction clearly.

Agents can read them through repository filesystem access in the intended
local development workflow.

Do not add instructions to implement corpus ingestion, vector search, or
semantic retrieval.

==================================================
7. DOMAIN PACK
==================================================

Document the current pack as:

    bmist@1.1.0

if and only if that is what the current registry/source actually reports.

State that:

- `bmist@1.0.0` is historical/immutable
- `bmist@1.1.0` is current

Do not describe an older pack as current.

Document the actual debt vocabulary only from the current pack.

Do not duplicate methodology definitions unnecessarily.

==================================================
8. TRUST UI
==================================================

Document the current human-facing Trust UI accurately.

Current primary surfaces:

    /ui/insights
    /ui/debts

Explain:

Insights
- operational/research state
- current claims
- tasks
- challenges
- debt
- other currently implemented projections

Debts
- debt obligations
- evidence
- retirement rules
- human discharge path

Important:

The UI is a human control surface.

It must not be described as AI authority.

Do not write "AI verified debt" or "agent retired debt."

Human adjudication remains distinct from agent work.

==================================================
9. DEVELOPER EXPERIENCE
==================================================

The primary local entrypoint is:

    argus serve

Document what it actually does.

Based on the current implementation, verify and document:

- CockroachDB bootstrap behavior
- database creation
- migration behavior
- seed behavior
- HTTP port
- Trust UI
- default local operator behavior
- relevant environment variables
- MCP startup behavior
- whether `argus mcp` requires a separate process

Do not describe an imaginary Taskfile workflow.

Do not require Docker Compose unless the current repository actually requires
it.

The README should make the first successful run obvious:

    build
    argus serve
    open Trust UI
    connect OpenCode using argus mcp

Use exact current commands from the repository.

==================================================
10. DRY RUN STATUS
==================================================

Do NOT claim the full Work-Agent + Adversarial-Agent research dry run has
already occurred unless it actually has.

The current implementation has already demonstrated:

- clean compilation
- test suite passing
- ARGUS bootstrap
- seed creation
- get_context
- submit_packet
- work belief persistence
- adversarial contradiction persistence
- scenario-specific seed idempotency
- Trust UI Insights rendering after the recent dashboard fix
- Trust UI Debts rendering

The actual human-run OpenCode research dry run is the next operational
experiment.

Document this distinction precisely.

Do not confuse:

"system integration dry run"

with:

"actual AI research dry run"

The latter happens when the human actually runs fresh OpenCode Work and
Adversarial agents against ARGUS.

==================================================
11. CURRENT TEST STATUS
==================================================

Inspect the current repository before documenting test status.

Report only tests that actually pass in the current tree.

Do not preserve stale claims from older README material.

Include the commands developers should actually run, e.g. where applicable:

    go test ./...
    go vet ./...
    go test -race ./...

But verify them first.

Do not claim race/vet success unless currently verified.

==================================================
12. SWEEPING REPOSITORY CLEANUP
==================================================

Now perform a repository-wide documentation/obsolete-artifact audit.

Do NOT blindly delete files.

First inventory likely obsolete artifacts and classify each:

KEEP
UPDATE
ARCHIVE
DELETE

Look especially for remnants of previous assumptions:

- old multi-service startup instructions
- Docker Compose instructions
- standalone Coordinator server instructions
- standalone Solvent/Conductor runtime instructions
- old MCP architecture documentation
- vector-search/corpus-ingestion plans
- embedding infrastructure plans
- superseded implementation plans
- stale phase plans
- duplicate architectural descriptions
- old pack-version assumptions
- obsolete dry-run instructions
- generated planning artifacts
- abandoned prototypes
- duplicate prompts
- obsolete scripts
- dead configuration files

IMPORTANT:

Do not delete real code simply because it came from an earlier phase.

Before deleting anything, prove that it is unused or superseded.

Use code search/reference analysis.

For documents that still have historical value but must not influence current
development, prefer:

    docs/archive/

over leaving them beside active instructions.

But do not create an archive dumping ground unnecessarily.

Delete genuinely redundant material.

The cleanup criterion is:

> A fresh AI coding agent should be able to inspect the repository and
> understand the current ARGUS architecture without being misled by obsolete
> plans.

==================================================
13. SPECIAL CLEANUP TARGETS
==================================================

Pay particular attention to the old documentation that claims:

- Solvent is a separate service
- Conductor is a separate service
- Coordinator needs a standalone server
- Trust UI expects another Coordinator process
- Phase 8 requires multiple separately launched services
- bmist@1.0.0 is current
- a separate RCP server must be launched
- the actual OpenCode run has not been integrated
- corpus/vector infrastructure is required
- Taskfile orchestration is required

Verify every such claim against current code.

Correct, archive, or delete as appropriate.

==================================================
14. README.md
==================================================

Rewrite README.md for a HUMAN developer/operator.

It should answer, quickly:

1. What is ARGUS?
2. Why does it exist?
3. What is the architecture?
4. How do I build it?
5. How do I run it?
6. How do I open the UI?
7. How does an AI agent connect?
8. What can the AI agent do?
9. What can the AI agent NOT do?
10. How do I run the tests?
11. How do I run the actual Work/Adversarial dry run?
12. What is deliberately not implemented yet?

Lead with working commands, not historical architecture.

Make the AI interface section prominent.

A human should understand that the key developer experience is:

    argus serve

and the key AI-agent interface is:

    argus mcp

Do not bury MCP beneath research background.

==================================================
15. OPENAI / OPENCODE / HERMES / OTHER AGENT HARNESS GUIDANCE
==================================================

Do not hardcode documentation to OpenCode alone.

Use terminology such as:

"AI coding agent / agent harness"

and give OpenCode as the current concrete example.

Document the generic model:

    Agent Harness
        |
        | MCP stdio
        v
    argus mcp
        |
        v
    argus.get_context
    argus.submit_packet

Explain that another MCP-capable agent harness such as Hermes can use the
same ARGUS MCP surface.

Do not claim compatibility with a specific harness unless verified.

The objective is that ARGUS exposes a stable capability boundary independent
of which agent harness is driving it.

==================================================
16. DO NOT ADD PROMPTS
==================================================

There are already:

    prompts/opencode-work.md
    prompts/opencode-adversarial.md

Keep them if they are current.

Update them only if they contain demonstrably stale architectural
instructions necessary to make the current system work.

Do not create additional role-card/prompt variants.

==================================================
17. DOCUMENTATION STYLE
==================================================

Prefer:

- concrete commands
- current paths
- actual interfaces
- short architectural diagrams
- explicit capability boundaries
- exact limitations

Avoid:

- speculative future architecture
- long historical narratives
- obsolete phase descriptions
- duplicated explanations
- claims unsupported by current code
- "will eventually" statements unless specifically marked as deferred

Current code beats historical documentation.

==================================================
18. ACCEPTANCE CRITERIA
==================================================

The task is complete only when:

AC1
AGENTS.md accurately describes the current ARGUS architecture.

AC2
README.md accurately describes the current ARGUS architecture.

AC3
Neither document claims that separate Solvent/Conductor/Coordinator
services are required unless the current implementation actually requires
them.

AC4
`argus serve` is documented as the primary local entrypoint.

AC5
`argus mcp` is documented as the primary AI-agent transport.

AC6
The exact two current MCP tools are documented.

AC7
The agent capability boundary is explicitly documented.

AC8
The distinction between agent work and human authority is documented.

AC9
The static seven-file research background is documented correctly.

AC10
The current `bmist@1.1.0` pack is documented correctly and 1.0.0 is not
described as current.

AC11
The current Trust UI routes are documented correctly.

AC12
The actual developer/test commands are verified against the repository.

AC13
The upcoming live Work-Agent / Adversarial-Agent dry run is accurately
described as an experiment rather than retroactively claimed as completed.

AC14
Stale architecture and obsolete planning material has been removed,
updated, or clearly archived.

AC15
No live code required by the current ARGUS runtime is deleted.

AC16
No new services, databases, MCP tools, agent prompts, or product features
are added.

AC17
A fresh AI coding agent reading only AGENTS.md and the repository can
understand how to interface with ARGUS without consulting historical plans.

AC18
A human developer reading README.md can get ARGUS running and understand how
to connect an MCP-capable AI agent.

==================================================
19. FINAL REPORT
==================================================

At completion report:

1. AGENTS.md changes
2. README.md changes
3. files deleted
4. files archived
5. files retained but explicitly marked historical
6. stale architectural claims removed
7. actual AI-agent interfaces documented
8. tests run and results
9. any remaining ambiguity

Do not make additional architectural changes.

The goal is repository truthfulness.

The current ARGUS implementation is the authority.