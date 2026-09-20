You are the senior architect for a major POC simplification pivot.

DO NOT IMPLEMENT CODE YET.

Your task is to inspect the current ARGUS / Oracle repository and write a NEW, implementation-ready plan for the simplified POC architecture described below.

The goal is to step back from the current multi-service design and remove accidental infrastructure complexity while preserving the architectural thesis.

==================================================
1. PIVOT OBJECTIVE
==================================================

We are pivoting the POC from a distributed multi-process stack into a single-process Go application with one database.

The POC must demonstrate:

- domain-agnostic epistemic workflow
- provenance-aware beliefs/evidence
- unresolved debt
- contradiction and lineage
- human-gated consequential transitions
- operational work separate from epistemic authority
- verification as a domain adapter
- agent interaction through MCP
- auditable decisions/history
- a reusable domain-pack model

Physics / BM-IST is ONLY the reference domain implementation.

The POC must make it clear that the reusable platform is the domain-agnostic core, while BM-IST / physics is replaceable domain-specific methodology.

==================================================
2. TARGET ARCHITECTURE
==================================================

Target runtime:

    OpenCode
       |
       | MCP stdio
       v
    ARGUS
    one Go process
       |
       v
    CockroachDB

ARGUS should expose:

    argus serve
    argus mcp
    argus verify
    argus reset   (only if useful)

Prefer one binary with subcommands over multiple independently deployed binaries.

HTTP/API + Trust UI should be served by the same ARGUS process.

MCP should be another mode of the same binary.

Physics verification should be a Go library/package, with a CLI/subcommand for running the reference verifier.

There should be no requirement for:

- separate Coordinator service
- separate Conductor service
- separate Solvent API service
- separate Trust UI process
- separate verifier daemon
- artifact bootstrap between processes
- internal HTTP calls merely to connect components within ARGUS

==================================================
3. WHAT MUST REMAIN CONCEPTUALLY SEPARATE
==================================================

Do NOT collapse the architecture semantically just because deployment is collapsed.

Preserve these invariants:

    CAPABILITY != WORK != AUTHORITY != EXECUTION

Preserve these logical planes:

    Operational state     -> work subsystem
    Epistemic state       -> Solvent subsystem
    Authorization state   -> decision/authority boundary
    External effect truth -> executor / external SOR

The following should remain explicit package/module boundaries:

    epistemic / Solvent
    work / Conductor replacement
    application orchestration / Coordinator replacement
    verifier
    domain-pack
    decision / authority
    MCP
    HTTP/API
    UI / projections

The distinction must be enforced through APIs and package ownership, not process boundaries.

==================================================
4. COMPONENT DECISIONS
==================================================

Treat the following as the baseline hypothesis for the new plan:

SOLVENT
- KEEP conceptually.
- Prefer in-process usage.
- Investigate whether existing Solvent Go code can be reused directly.
- Do not automatically preserve the Solvent HTTP server just because it already exists.
- Preserve its epistemic/authority semantics and canonical data model.

CONDUCTOR
- REMOVE as a standalone service.
- Replace with the smallest possible in-process operational task model.
- Preserve only the work semantics actually required by the POC:
  task
  dependency
  state
  priority
  activity
  cancellation/dead-end linkage

COORDINATOR
- REMOVE as a standalone service.
- Convert its useful logic into an in-process application/service layer.
- Preserve mechanical orchestration and contract validation.
- Do not allow it to become a scientific reasoning layer.

VERIFIER
- KEEP as a library.
- Physics/BM-IST is the first reference implementation.
- No separate verifier process-to-Coordinator artifact handoff.
- Verification results should enter the same in-process evidence/artifact pathway.

DOMAIN PACK
- KEEP.
- This is critical to proving domain portability.
- BM-IST must remain a domain pack/reference implementation.
- Domain-specific vocabulary and retirement rules must remain outside the generic epistemic kernel.

TRUST UI
- KEEP.
- Serve from the same ARGUS process.
- Do not create a separate Trust UI server.

MCP
- KEEP.
- One binary/subcommand.
- Exactly preserve the thin agent boundary:
    argus.get_context
    argus.submit_packet

ORACLE / SPHINX / CHRONICLER / PASSAGE LEDGER / etc.
- Treat these as views, projections, queries, or application modules.
- They are not separate deployable components.

==================================================
5. DATABASE
==================================================

Prefer one CockroachDB instance/database for the POC.

Do not introduce:

- second database
- second persistence system
- local artifact database
- service-local databases

Re-use the existing Solvent schema where technically sensible.

If the current Solvent schema cannot cleanly support the simplified architecture, identify the smallest required change instead of reproducing the current service topology.

==================================================
6. FIRST TASK: INSPECT REAL CODE
==================================================

Before proposing the plan, inspect the repository thoroughly.

Determine:

1. Which Solvent code is reusable in-process.
2. Which Coordinator logic is worth keeping.
3. Which Conductor functionality is actually used by the POC.
4. Which verifier code is already cleanly reusable.
5. Which Trust UI code can be embedded/served directly.
6. Which MCP code can become a subcommand.
7. Which database schema elements are truly required.
8. Which current code exists only because of the multi-service deployment model.
9. Which current tests encode real architectural requirements versus deployment mechanics.
10. Which existing code should be deleted rather than adapted.

Do not assume the current implementation matches old planning documents.

Code is the source of truth for implementation details.

==================================================
7. REQUIRED ARCHITECTURAL OUTCOME
==================================================

Produce a proposed repository structure approximately like:

    cmd/
      argus/

    internal/
      application/
      epistemic/
      work/
      decision/
      verifier/
      domainpack/
        bmist/
      api/
      mcp/
      ui/

    migrations/
    testdata/
    domain-pack/
    Taskfile.yml

This is illustrative, not mandatory.

Choose the smallest structure that preserves the conceptual boundaries.

==================================================
8. REQUIRED END-TO-END DEMONSTRATION
==================================================

The new POC must support this complete narrative:

1. Start ARGUS.
2. OpenCode connects through MCP.
3. Agent requests context.
4. Agent performs bounded domain work.
5. Agent submits an EBP/claim packet.
6. ARGUS validates and persists epistemic proposals/evidence.
7. Debt is visible.
8. Physics verifier provides reproducible verification evidence.
9. An adversarial agent introduces a contradiction.
10. Contradiction/lineage becomes visible.
11. Human reviews the consequential transition.
12. Human discharges appropriate debt where permitted by the BM-IST pack.
13. Claim is promoted subject to existing Solvent rules.
14. A later contradiction can cause human retraction.
15. Associated operational work can become a dead end.
16. REOPEN creates new lineage without erasing history.
17. UI shows the full state and history.

The demonstration should prove the architecture, not the infrastructure.

==================================================
9. DOMAIN-AGNOSTICITY TEST
==================================================

The plan must explicitly separate:

GENERIC CORE

from:

BM-IST / PHYSICS DOMAIN PACK

Create a portability test in the plan.

The plan should explain exactly which parts would remain unchanged if the market were to replace BM-IST with another domain such as:

- financial risk
- software security
- compliance
- engineering verification
- supply-chain assurance

Do NOT implement those additional domains now.

Instead, demonstrate that replacing the domain pack does not require changing the generic epistemic/authority workflow.

==================================================
10. DELETE-FIRST REQUIREMENT
==================================================

This pivot is specifically intended to reduce complexity.

For every existing component/service, classify it as:

    KEEP
    CONVERT TO LIBRARY
    CONVERT TO INTERNAL PACKAGE
    DELETE
    DEFER

Explicitly identify infrastructure that exists only because of the old topology.

Do not preserve something merely because it already exists.

The default should be subtraction.

==================================================
11. IMPLEMENTATION PLAN FORMAT
==================================================

Write a new document:

    docs/plan8/PIVOT_POC_IMPLEMENTATION_PLAN.md

The plan must contain:

1. Executive summary
2. Current-state findings from the actual repository
3. Target architecture
4. Component-by-component disposition
5. Package/module boundaries
6. Database strategy
7. MCP strategy
8. Trust UI strategy
9. Verifier/domain-pack strategy
10. Migration/deletion strategy
11. Phased implementation plan
12. Test strategy
13. End-to-end demo scenario
14. Domain-portability demonstration
15. Developer experience
16. Risks and deliberate non-goals
17. Acceptance criteria
18. Explicit list of files/directories expected to be deleted, moved, created, or modified

==================================================
12. PHASE DESIGN
==================================================

Keep the phases small and executable.

Prefer something like:

Phase 0 — repository reconnaissance
Phase 1 — target architecture / interfaces
Phase 2 — collapse runtime topology
Phase 3 — in-process epistemic subsystem
Phase 4 — minimal work subsystem
Phase 5 — application/orchestration layer
Phase 6 — verifier + BM-IST pack
Phase 7 — MCP + context/packet flow
Phase 8 — Trust UI
Phase 9 — end-to-end adversarial demonstration
Phase 10 — portability / domain-boundary proof
Phase 11 — cleanup, deletion, acceptance

Adjust the phases based on actual code.

Do not create phases merely to make the plan look comprehensive.

==================================================
13. HARD CONSTRAINTS
==================================================

Do NOT:

- add microservices
- introduce Kubernetes
- introduce service meshes
- add message brokers
- add event buses
- introduce distributed transactions
- add cryptographic attestation infrastructure
- introduce a second database
- add autonomous promotion
- add autonomous adjudication
- put LLM reasoning inside the epistemic kernel
- create a second domain pack during this POC
- rebuild the current architecture simply under different names
- retain Coordinator or Conductor as services for compatibility alone

Do not solve hypothetical future scale problems.

==================================================
14. CRITICAL DESIGN QUESTION
==================================================

The plan must directly answer:

"What architectural property would be lost if Coordinator and Conductor stopped being separate services?"

If the answer is "none", the plan must explicitly state that their conceptual roles survive as internal modules while their process boundaries are removed.

Likewise answer:

"What makes Solvent worth preserving as a distinct conceptual subsystem?"

The answer must focus on the epistemic/authority model, not deployment.

==================================================
15. DEVELOPER EXPERIENCE
==================================================

The target developer experience should be approximately:

    git clone
    task dev

Then:

    ARGUS ready
    Trust UI ready
    MCP ready
    database ready

And:

    task test
    task verify

should exercise the actual POC rather than merely checking service plumbing.

Avoid multi-service startup scripts, PID management, readiness orchestration, API credential propagation, and artifact handoff unless the repository inspection proves one is still necessary.

==================================================
16. SUCCESS CRITERIA
==================================================

The proposed POC is successful when:

- one Go application can run the entire demonstration
- one database is sufficient
- the agent boundary remains explicit
- epistemic authority remains distinct from operational work
- human consequential authority remains explicit
- physics verification remains domain-specific
- domain-independent workflow remains reusable
- Trust UI makes the state understandable
- the complete adversarial workflow can be demonstrated locally
- a hypothetical second domain would require a new domain pack, not a rewrite of the core
- the resulting system is materially simpler than the current multi-service architecture

The plan must explicitly compare:

CURRENT POC COMPLEXITY
versus
PIVOTED POC COMPLEXITY

using concrete counts where possible:

- processes
- ports
- databases
- repositories/modules
- credentials
- startup steps
- inter-service HTTP boundaries
- persistence boundaries

==================================================
17. IMPORTANT
==================================================

This is a PLAN-WRITING task only.

Do not modify implementation code.

Do not create the new architecture.

Do not start the migration.

Do not implement any of the phases.

Inspect the repository, reason from the actual code, and produce:

    docs/plan8/PIVOT_POC_IMPLEMENTATION_PLAN.md

The plan must be decisive.

Do not preserve old architecture merely because it already exists.

The purpose of this exercise is to determine the smallest system that can convincingly prove the ARGUS thesis.