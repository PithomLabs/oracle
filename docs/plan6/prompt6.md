Use this prompt for the coding agent:

```text
Draft the top-level README.md for the ARGUS / Oracle Physics Trust Verifier repository.

Do NOT change implementation code.
Do NOT invent capabilities that are not implemented.
Do NOT claim Phase 8 is fully successful yet.
The README should document the system as it exists through Phase 8 implementation, immediately before / ready for the first live Phase 8 dry run.

==================================================
PURPOSE
==================================================

ARGUS is a trust-verification proof of concept demonstrating a control
architecture for agent-driven research where:

CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

The README should explain the architecture, why the system exists,
how the components interact, what has been implemented through Phase 8,
how the Work and Adversarial OpenCode agents interact with ARGUS, and
how the human remains the authority for epistemic debt and consequential
decisions.

The README is for:
- engineers evaluating the architecture
- researchers interested in the workflow
- reviewers evaluating the POC
- someone approaching the repository for the first time

Keep it technically precise and relatively concise.
Prefer diagrams and concrete examples over marketing language.

==================================================
SYSTEM IDENTITY
==================================================

Use these terms consistently:

ARGUS
  Trust Verification POC / control surface

Oracle
  The epistemic query/projection side of ARGUS:
  "What does the system currently know?"

Solvent
  Epistemic authority ledger:
  beliefs, evidence, debt, edges, promotion,
  retraction, authorization, audit

Conductor
  Operational workflow/task coordination

Trust UI
  Human observation and adjudication surface

EBP v2.1
  Epistemic operating discipline

RCP
  Research Context Protocol

MCP
  Agent-facing transport

OpenCode
  Agent harness

Work Agent
  OpenCode used to perform bounded research work

Adversarial Agent
  Separate OpenCode process used to challenge prior work

BM-IST
  Domain Pack / research program used for the POC

Gate G0
  Current Phase 8 dry-run research slice

==================================================
CORE DOCTRINE
==================================================

State prominently:

Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
No final-truth claim may be promoted.
Accounting must never become the work.

Also state:

Human beings discharge epistemic debt.

Agents may produce work, evidence, artifacts, challenges and proposed
findings, but they do not directly retire debt or promote claims.

==================================================
ARCHITECTURE
==================================================

Include a clear architecture diagram similar to:

                BM–IST Research Program
                         │
                         ▼
                  Domain Pack / EBP
                         │
              ┌──────────┴──────────┐
              ▼                     ▼
        Work OpenCode        Adversarial OpenCode
              │                     │
              └──────────┬──────────┘
                         ▼
                    ARGUS / Oracle
                         │
                    Coordinator
                   /             \
                  ▼               ▼
             Conductor          Solvent
          operational state   epistemic authority
                  ▲               ▲
                  │               │
                RCP/MCP       authoritative state
                  ▲
                  │
              OpenCode agents

And the human control path:

Human
  ↓
Trust UI
  ↓
Coordinator
  ↓
Solvent

Make the boundaries explicit:

OpenCode = agency/work
Conductor = operational state
Solvent = epistemic authority
Trust UI = human adjudication
External systems = effect truth where applicable

==================================================
WHY ARGUS EXISTS
==================================================

Explain the central problem:

Agent/workflow completion is not equivalent to epistemic authority.

A work agent may produce:
- a claim
- evidence
- proof material
- an experiment
- a counterexample
- a proposed retirement of debt

But the system must distinguish:
- work being performed
- claims being recorded
- evidence being accepted
- debt being discharged
- authority being granted

The README should explain that Solvent is the authority boundary and
Coordinator is the orchestration boundary.

Avoid claiming that ARGUS proves scientific truth.
It does not.

==================================================
EBP WORKFLOW
==================================================

Show the workflow:

Idea
→ Enter
→ Work
→ Evidence
→ Adversarial Challenge
→ Human Adjudication
→ Debt Discharge
→ Promotion Request
→ Solvent Gate

Explain:

- ideas may enter without promotion
- debt is attached to promoted/candidate beliefs
- debt blocks promotion
- human discharges debt
- Solvent records authoritative state
- agents cannot declare their own work authoritative

Clarify that Phase 8 intentionally defers ADD_DEBT / automatic
"new evidence creates new debt" behavior.

This must be described as:
DEFERRED / NOT EXERCISED

not as fully implemented behavior.

==================================================
OPEN CODE WORKFLOW
==================================================

Explain the two-agent workflow.

### Work Agent

1. Human selects/creates a Conductor task.
2. Fresh OpenCode requests context through:
   MCP → ARGUS → RCP → Conductor + Solvent.
3. Agent performs bounded research.
4. Agent produces an EBP work packet.
5. `argus.submit_packet` submits it.
6. Coordinator persists work to Solvent/Conductor.

### Adversarial Agent

1. Fresh separate OpenCode process.
2. No conversational history from the Work Agent.
3. Calls `argus.get_context`.
4. Reconstructs current state from RCP.
5. Sees existing work, evidence, debt, dependencies,
   activity and prior challenges.
6. Performs adversarial review.
7. Produces an adversarial EBP packet.
8. `argus.submit_packet` records the challenge.
9. Adversarial challenges are represented structurally through
   `contradicts` edges.

Emphasize:

> Agents reconstruct context from ARGUS state, not from inherited
> conversation memory.

==================================================
RCP
==================================================

Document RCP as:

> A thin read-only context projection over existing Conductor and
> Solvent state. It is a protocol, not a persistence layer.

Canonical endpoint:

GET /v1/context/{task_id}

Canonical protocol identity:

RCP/v1

Describe the conceptual response:

- task
- dependencies
- epistemic state
  - beliefs
  - evidence
  - edges
  - debt
  - intents
- activity

Explain:
- full scenario projection for Phase 8
- no new research database
- no sophisticated graph traversal yet
- source attribution
- eventual consistency
- UNKNOWN ≠ EMPTY
- backend failure must not look like empty research state
- deterministic presentation without inventing cross-ledger
  happens-before semantics

==================================================
MCP CAPABILITY BOUNDARY
==================================================

The Work and Adversarial OpenCode processes expose exactly:

argus.get_context
argus.submit_packet

They do NOT receive:
- Solvent MCP
- Conductor write MCP tools
- direct DB access
- edge mutation MCP tools

Explain why this matters:

> Capability restriction is enforced by tool-surface configuration, not
> by telling the agent to behave.

Include the conceptual path:

OpenCode
→ MCP
→ ARGUS adapter
→ Coordinator

==================================================
EBP PACKETS
==================================================

Explain that EBP packets are the common write grammar.

Packet categories:
- work
- adversarial

Core packet objects include:
- beliefs
- evidence
- edges
- tasks

Reference grammar:

local:<id>
  = intra-packet reference

canonical:belief:<uuid>
  = existing Solvent belief reference

Explain that packets contain proposed work/results, not authority.

==================================================
PHASE 8 IMPLEMENTATION
==================================================

Summarize implemented Phase 8 capabilities:

### RCP
- context endpoint
- Conductor + Solvent projection
- availability metadata
- deterministic activity presentation

### Packet persistence
- beliefs persisted to Solvent
- evidence persisted to Solvent
- derives/contradicts edges persisted to Solvent
- tasks persisted to Conductor
- idempotency is scenario-scoped

### Solvent Growth Gate
Document the minimal Phase 8 addition:

POST /v1/beliefs/{parent_id}/edges

Accepts:
- child_id
- kind = derives | contradicts

Solvent validates:
- parent exists
- child exists
- parent != child
- valid kind
- uniqueness

No edge-creation MCP tool is exposed to agents.

### Human debt discharge
Explain the enforced path:

Trust UI
→ authenticated Coordinator
→ Pack membership
→ retirement rule
→ actual evidence verification
→ attributed Solvent /v1/discharge

Mention that:
- unknown debt fails closed
- Pack lookup failure fails closed
- caller-supplied evidence class alone is insufficient
- human operator identity is server-side
- discharge records attribution

### Retraction / Dead End
Explain:

contradicts edge
→ human retraction
→ Solvent RetractCascade
→ linked Conductor tasks cancelled
→ Insights derives Dead End structurally

Dead End requires:
- terminal cancelled task
- governing belief retracted/contradicted

Cancelled alone is not a dead end.
Rejected is not terminal.

### REOPEN
Explain:

retracted belief
→ new belief
→ derives edge
→ original remains in history

### Trust UI
Only two surfaces:

Insights
Debts

Insights:
- research statistics
- adversarial challenges
- open claims
- open debt
- promoted
- refusals
- dead ends
- Conductor priority ordering

Debts:
- debt item
- why it exists
- evidence
- adversarial challenges
- actual Pack retirement rule
- evidence class
- human discharge action

Explicitly say:

The UI never presents:
"AI verified debt"
"Agent retired debt"

==================================================
BM-IST G0 DRY RUN
==================================================

Introduce Gate G0 as the current Phase 8 research slice.

Include only the documented POC-level description:

L1:
Orbit Rigidity

L2:
One-Parameter Triviality

G0:
countable substrate + literal continuous unitary evolution
→ structural incompatibility

Dependencies:
L1 → G0
L2 → G0

Do not turn the README into a physics paper.
Do not independently validate or embellish the scientific claims.
Describe them as the research claims selected for the POC.

==================================================
DRY-RUN BRANCHES
==================================================

Explain that Phase 8 deliberately has two complementary branches.

### Branch A — adversarial / dead end

Work
→ adversarial challenge
→ contradicts edge
→ human retracts
→ Solvent cascade
→ Conductor task cancelled
→ dead end visible

### Branch B — debt / promotion

surviving belief or successor
→ human reviews debt
→ qualifying evidence
→ Coordinator mechanical validation
→ attributed Solvent discharge
→ promotion request
→ Solvent gate

IMPORTANT:
Do not imply that a retracted belief is later promoted.

==================================================
CURRENT STATUS
==================================================

State precisely:

Phase 8 implementation:
READY FOR DRY RUN

Do NOT state:
"Phase 8 fully successful"

The actual end-to-end dry run has not yet been executed.

Known verification status:
- targeted adversarial implementation fixes pass
- race detector clean
- go vet clean
- Oracle verification complete
- Conductor verification complete
- Solvent DB-dependent integration tests require live CockroachDB
- actual Work-Agent / Adversarial-Agent OpenCode run remains outstanding
- Trust UI currently has source-level verification rather than a full automated UI suite

Make clear that the next milestone is the actual live dry run.

==================================================
PROJECT STRUCTURE
==================================================

Inspect the actual repository and document its real current structure.

At minimum explain the purpose of:
- oracle/
- trust-ui/
- domain-pack/
- packet/
- verifier/
- corpus/
- coordinator/
- mcp/
- Solvent
- Conductor

Do NOT invent directories that do not exist.

Use the actual repository layout discovered from the codebase.

==================================================
RUNNING / DEVELOPMENT
==================================================

Inspect the actual repository scripts, binaries, module files and current
README conventions before writing this section.

Document only commands that actually work in the current repository.

Include, as applicable:
- prerequisites
- Go version/module expectations
- Conductor startup
- Solvent/CockroachDB requirement
- Oracle/Coordinator startup
- Trust UI startup
- ARGUS MCP startup
- OpenCode configuration
- test commands

Do not fabricate commands.

Explicitly identify anything required from a live CockroachDB instance.

==================================================
SECURITY / AUTHORITY NOTES
==================================================

State clearly:

- Trust UI decisions require authenticated access to Coordinator.
- Operator identity is configuration-trusted for the Phase 8 POC.
- `ARGUS_OPERATOR_PRINCIPAL_ID` is server-side.
- Browser-provided actor identity is not trusted.
- Agents do not have Solvent or Conductor mutation tools.
- Direct DB access is prohibited from agents/UI/Coordinator.
- Full multi-user human authentication is beyond Phase 8 scope.

==================================================
LIMITATIONS / DEFERRED
==================================================

Include an explicit list of Phase 8 limitations, grounded in actual
implementation.

Known examples:
- ADD_DEBT deferred
- automatic "new evidence creates debt" deferred
- full multi-user authentication deferred
- project selection currently POC-scoped if still true
- RCP full scenario projection rather than sophisticated graph closure
- live DB needed for complete integration run
- actual agent-driven dry run pending
- UI automated test coverage limited if still true

Do not turn limitations into complaints.
They are boundaries of the POC.

==================================================
DESIGN PRINCIPLES
==================================================

End with the most important principles:

1. Agents do work; they do not own authority.
2. Conductor remembers operational work.
3. Solvent records epistemic authority.
4. EBP governs debt and promotion.
5. RCP provides agent context without becoming a database.
6. MCP constrains agent capability.
7. Humans discharge epistemic debt.
8. Solvent decides consequential gates.
9. Refusals are first-class evidence.
10. Unknown must not silently become empty or allowed.

==================================================
STYLE
==================================================

Write a serious engineering README.

Avoid:
- marketing claims
- "revolutionary"
- "AI-powered" hype
- unsupported claims of scientific validity
- unsupported claims of production readiness

Prefer:
- concise sections
- ASCII diagrams
- tables where useful
- concrete request/response examples
- explicit boundaries
- current status and limitations

The README should make a new engineer understand within five minutes:

WHAT ARGUS IS
WHY IT EXISTS
HOW WORK AGENTS AND ADVERSARIAL AGENTS USE IT
WHERE AUTHORITY LIVES
HOW EBP DEBT WORKS
WHAT PHASE 8 HAS ACTUALLY IMPLEMENTED
WHAT HAS NOT YET BEEN PROVEN
HOW TO RUN THE POC

==================================================
DELIVERABLE
==================================================

Create or update:

README.md

Do not modify application behavior.

Before finalizing:
1. inspect the actual repository
2. verify every command/example in README against the code
3. verify terminology against the Phase 8 implementation
4. verify status claims against the current verification reports
5. do not claim the dry run occurred

Then report:
- README created/updated
- sections included
- any repository facts you could not verify
- no code behavior changed
```
