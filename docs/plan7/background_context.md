BACKGROUND CONTEXT — ARGUS / ORACLE TRUST VERIFICATION POC

You are reviewing an engineering plan for a project called ARGUS
(implemented primarily in the Oracle repository). Assume you know
nothing about the project's dependencies before reading this briefing.

Your role is NOT to implement the system. You are an independent
adversarial reviewer. Your job is to understand the architecture,
identify hidden coupling, authority-boundary failures, missing
prerequisites, incorrect assumptions, and places where a plan claims
something that the implementation cannot actually guarantee.

==================================================
1. WHAT ARGUS IS
==================================================

ARGUS is a proof-of-concept trust-verification/control architecture for
agent-driven research.

Its central invariant is:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

The basic idea is that an AI agent may perform substantial research work
without being allowed to declare that work authoritative.

ARGUS therefore separates:

    Agent / OpenCode
        = agency and work

    Conductor
        = operational workflow/task state

    Solvent
        = epistemic authority and consequential state

    Coordinator
        = orchestration and boundary between agents, Conductor and Solvent

    Trust UI
        = human observation and adjudication

    EBP
        = epistemic operating discipline

ARGUS itself is NOT intended to prove scientific truth.

It is intended to make the distinction between:
- work
- claims
- evidence
- unresolved obligations/debt
- human adjudication
- authority transitions

explicit and enforceable.

==================================================
2. EBP v2.1 — THE GOVERNING RESEARCH DISCIPLINE
==================================================

ARGUS follows EBP v2.1:

    Ideas enter free.
    Promotion costs debt.
    Debt does not kill.
    Debt is forever payable.
    New evidence creates new debt.
    No final-truth claim may be promoted.
    Accounting must never become the work.

The important operational interpretation is:

    Agents produce work.
    Humans discharge epistemic debt.
    Solvent records/enforces authority.

An agent may propose:
- a claim
- evidence
- an artifact
- a counterexample
- an adversarial finding
- a proposed debt retirement

but the agent cannot directly:
- retire debt
- promote a belief
- authorize a consequential action
- mutate Solvent authority state
- mutate Conductor workflow state

Phase 8 deliberately defers the operational ADD_DEBT / automatic
"new evidence creates new debt" loop. Do not assume it is implemented.

==================================================
3. SOLVENT
==================================================

Solvent is the epistemic authority service.

Think of it as the system of record for:

- beliefs/claims
- evidence
- debt
- belief relationships/edges
- action intents
- promotion
- retraction
- authorization
- audit history

Solvent determines whether a consequential epistemic transition is
actually permitted.

For example:

    belief has open debt
        → promotion is refused

The important architectural property is:

    Coordinator may request a consequential transition,
    but Solvent is the final authority.

Solvent is a separate service/repository, not merely a library embedded
inside Oracle.

It has a REST API.

A minimal Phase 8 example is:

    POST /v1/beliefs
        enter belief

    POST /v1/evidence
        add evidence

    POST /v1/beliefs/{id}/edges
        create derives/contradicts edge

    POST /v1/discharge
        record attributed human debt discharge

    POST /v1/beliefs/{id}/promote
        request promotion

    POST /v1/beliefs/{id}/retract
        retract with cascade

The Phase 8 edge endpoint is a deliberate minimal Growth Gate extension:
the belief-edge table already exists and is already used by Solvent's
retraction logic; the missing piece was a lawful API for creating edges.

Agents do NOT receive Solvent's MCP mutation tools.

==================================================
4. CONDUCTOR
==================================================

Conductor is the operational workflow/task system.

It owns things such as:

- projects
- tasks
- dependencies
- task status/lifecycle
- activities
- governance references
- task priority

It answers:

    "What work exists, what depends on what,
     and what is the operational state of that work?"

It does NOT own epistemic authority.

For example:

    Conductor: task is accepted
does NOT mean:
    Solvent: claim is promoted

Likewise:

    Conductor: task is cancelled
is operational state, not scientific truth.

Conductor is also a separate service/repository with a REST API.

==================================================
5. COORDINATOR
==================================================

Coordinator lives in the Oracle repository.

It is the orchestration boundary between:
- agents
- Conductor
- Solvent
- Domain Pack

It is deliberately NOT another authority store.

The Coordinator currently exists primarily as a Go library plus HTTP
handlers; a standalone launcher is being considered/added to make the
developer experience reproducible.

Coordinator responsibilities include:

- validate EBP packets
- resolve Domain Pack
- perform idempotency
- persist packet objects through Conductor/Solvent REST APIs
- expose RCP context
- route human decisions
- mechanically validate retirement rules
- preserve capability/authority boundaries

Coordinator must NOT:
- write directly to Solvent DB
- write directly to Conductor DB
- perform scientific adjudication
- infer truth
- become a research reasoning engine

A critical architectural distinction:

    Coordinator = orchestration
    Solvent     = authority

==================================================
6. RCP — RESEARCH CONTEXT PROTOCOL
==================================================

RCP is a very thin read protocol for agents.

It does not store anything.

It is essentially a projection of existing Conductor + Solvent state so
that a completely fresh agent can reconstruct the current research
context without relying on inherited conversation memory.

Canonical concept:

    GET /v1/context/{task_id}

The context contains roughly:

    task
    dependencies
    beliefs
    evidence
    edges
    debt
    intents
    activity

RCP is intentionally simple.

For Phase 8 it uses a broad scenario projection rather than sophisticated
graph traversal.

An important rule is:

    UNKNOWN ≠ EMPTY

If Solvent is unavailable, RCP must not make the agent believe that
"there are no beliefs."

RCP is deterministic in presentation, but it does NOT invent
cross-system happens-before semantics between independently clocked
Conductor and Solvent events.

==================================================
7. EBP PACKETS
==================================================

Agents communicate work to ARGUS using EBP research packets.

There are two relevant roles:

    work
    adversarial

A packet can contain:

- beliefs
- evidence
- edges
- tasks

The packet is a submission of work/results.

It is NOT authority.

The canonical reference grammar is:

    local:<id>
        = object inside the submitted packet

    canonical:belief:<uuid>
        = reference to an existing Solvent belief

The Coordinator resolves/validates packet references and persists the
resulting objects into the appropriate authoritative systems.

==================================================
8. OPENCODE AS AGENT HARNESS
==================================================

OpenCode is the agent harness.

There are two independent OpenCode processes in the intended workflow.

### Work Agent

A human gives the initial task.

The Work Agent:

    calls argus.get_context
        ↓
    reconstructs current state
        ↓
    performs bounded research
        ↓
    produces EBP packet
        ↓
    calls argus.submit_packet

### Adversarial Agent

A completely separate fresh OpenCode process is started.

It has NO conversational memory of the Work Agent.

It:

    calls argus.get_context
        ↓
    reconstructs what has already been done
        ↓
    identifies unresolved work/debt
        ↓
    attacks the current work
        ↓
    produces an adversarial EBP packet
        ↓
    calls argus.submit_packet

This is important:

    Agents reconstruct context from system state,
    not from inherited conversation memory.

The adversarial agent is intentionally supposed to challenge the work,
not simply improve it.

==================================================
9. MCP CAPABILITY BOUNDARY
==================================================

The agent-facing MCP surface is intentionally tiny.

Agents should receive only:

    argus.get_context
    argus.submit_packet

They should NOT receive:
- Solvent MCP mutation tools
- Conductor write tools
- database access
- direct edge-mutation tools
- promotion/discharge/retraction tools

This is a capability boundary, not merely a prompt convention.

If an agent cannot see a privileged tool, it cannot invoke it.

The intended path is:

    OpenCode
        ↓
    ARGUS MCP adapter
        ↓
    Coordinator
        ↓
    Conductor + Solvent

==================================================
10. TRUST UI
==================================================

The Trust UI is a very thin human control surface.

Phase 8 intentionally has only two primary views:

    Insights
    Debts

Insights answers:

    "What has happened and what needs attention?"

It uses Conductor task priority/state plus Solvent epistemic state.

Debts answers:

    "What does this claim still owe before promotion?"

The human can inspect:
- debt
- evidence
- adversarial challenges
- retirement rule
- relevant history

Then the human may request debt discharge.

The important path is:

    Human
      ↓
    Trust UI
      ↓
    authenticated Coordinator
      ↓
    mechanical Pack/evidence validation
      ↓
    Solvent attributed discharge

The UI must never imply:

    "AI verified debt"

or:

    "Agent retired debt"

The agent may supply evidence.
The human adjudicates.
Solvent records the authoritative transition.

==================================================
11. BM-IST / PHYSICS VERIFIER
==================================================

The current domain used to exercise ARGUS is BM-IST
(Bounded Mathematics — Invariant Structure Theory).

Phase 8 uses one narrow research slice, Gate G0.

The relevant claims are:

    L1 — Orbit Rigidity
    L2 — One-Parameter Triviality
    G0 — countable substrate + literal continuous unitary
         evolution → structural incompatibility

These are being used as research objects for the POC.

ARGUS is NOT independently claiming that these scientific claims are
true.

The Physics Verifier is a bounded verification component that can
produce trusted verification artifacts for specific mathematical/physics
obligations.

It is deliberately NOT:
- a general theorem prover
- a full computer algebra system
- a scientific authority

Its output becomes evidence/artifact material that can enter the ARGUS
workflow.

Solvent remains the authority layer.

==================================================
12. PHASE 8 DRY RUN
==================================================

The Phase 8 dry run is designed around two complementary branches.

### Branch A — adversarial / dead end

    Work Agent
        ↓
    research packet
        ↓
    Adversarial Agent
        ↓
    contradicts edge
        ↓
    human retracts
        ↓
    Solvent RetractCascade
        ↓
    linked Conductor task cancelled
        ↓
    Insights derives Dead End

### Branch B — human debt discharge / promotion

    surviving belief or successor
        ↓
    human reviews debt
        ↓
    qualifying evidence
        ↓
    Coordinator validates retirement rule
        ↓
    authenticated attributed discharge
        ↓
    Solvent
        ↓
    promotion request
        ↓
    Solvent promotion gate

A retracted belief is NOT later promoted.

==================================================
13. WHY THE CURRENT WORK IS ABOUT DEVELOPER EXPERIENCE
==================================================

Phase 8 implementation is currently ready for its first live dry run.

The remaining engineering goal is to make the developer experience
reproducible.

Today, several services are separate:
- CockroachDB
- Solvent
- Conductor
- Coordinator
- Trust UI
- ARGUS MCP
- optional Physics Verifier tooling

The current effort is to create a simple developer orchestration layer,
primarily through Taskfile, so a developer can move toward:

    task setup
    task up
    task status
    task test
    task verify
    task dry-run
    task down

The goal is not to create another application layer.

The Taskfile should only orchestrate existing services/binaries,
database readiness, configuration and tests.

==================================================
14. CURRENT ARCHITECTURAL QUESTIONS FOR REVIEW
==================================================

When reviewing implementation plans, pay particular attention to:

1. Are Solvent and Conductor being reused through their existing APIs,
   rather than bypassed with direct database access?

2. Does Coordinator remain orchestration rather than becoming a new
   authority layer?

3. Are agent capabilities actually constrained by tools/configuration?

4. Can a malicious browser/client bypass human adjudication?

5. Does the system distinguish caller claims from actual persisted
   evidence?

6. Does RCP reconstruct context without becoming another database?

7. Does a failed backend appear as UNKNOWN rather than EMPTY?

8. Are cross-ledger events presented deterministically without inventing
   false causal ordering?

9. Are EBP claims distinguished between implemented and deferred?

10. Does the proposed developer tooling reduce complexity rather than
    create another orchestration framework?

11. Can one command reliably start the complete local environment?

12. Can a fresh Work Agent and fresh Adversarial Agent actually operate
    against that environment?

13. Are Physics Verifier outputs treated as evidence/artifacts rather
    than scientific authority?

==================================================
15. YOUR REVIEW ROLE
==================================================

You are an independent adversarial architect.

Assume the implementers are intelligent and the plan is well intentioned.

Do not review for style.

Try to find:
- hidden authority bypasses
- fail-open behavior
- missing runtime dependencies
- incorrect service assumptions
- duplicated state
- missing readiness checks
- accidental coupling
- misleading documentation
- tests that only prove mocks rather than production paths
- developer workflows that appear simple but are not actually reproducible

Prefer concrete failures over speculative architecture expansion.

The objective is:

    Make the smallest system that can honestly demonstrate:

    agent work
        +
    adversarial review
        +
    human debt adjudication
        +
    Solvent authority
        +
    reproducible developer execution

without allowing the infrastructure itself to become the complexity.