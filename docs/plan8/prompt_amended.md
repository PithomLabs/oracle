You are the senior architect for a major POC simplification pivot.

DO NOT IMPLEMENT CODE YET.

Your task is to inspect the current ARGUS / Oracle repository and write a NEW, implementation-ready plan for the simplified POC architecture described below.

The goal is to step back from the current multi-service design and remove accidental infrastructure complexity while preserving the architectural thesis.

This prompt is the amended version of the original prompt.md, incorporating 9 net-valid amendments from adversarial review (prompt_review2.md). The amendments are marked with [AMENDMENT] tags.

==================================================
1. PIVOT OBJECTIVE
==================================================

We are pivoting the POC from a distributed multi-process stack into a single-process Go application with one database.

[AMENDMENT] The pivot is reconnaissance-gated, not predetermined. Coordinator and Conductor are provisionally slated for collapse. Repository reconnaissance must determine which semantics are actually load-bearing. The pivot proceeds only if those semantics can be preserved in-process without weakening the thesis.

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

[AMENDMENT] Import-direction rule. The epistemic package is the innermost layer. It imports nothing from application, work, api, mcp, ui, or any domain pack. Domain vocabulary (claim types, debt item names, evidence classes, retirement rules) is forbidden inside the epistemic kernel. Debt remains opaque to Solvent.

The enforcement boundary must be mechanical, not aspirational. Use:

- A depguard lint rule (golangci-lint) that forbids specific import paths per package.
- A compile-time negative test per invariant (e.g., work cannot promote beliefs; epistemic cannot execute tasks; verifier cannot authorize; mcp cannot reach privileged decision APIs).

These are the few invariants that matter. Do not introduce a general architecture-linting framework.

==================================================
4. COMPONENT DECISIONS
==================================================

[AMENDMENT] The following decisions are provisional. Reconnaissance must verify that the semantics being preserved are actually load-bearing for the POC thesis. If reconnaissance reveals that a component slated for removal contains semantics that cannot be preserved in-process without weakening the thesis, the plan must revise the decision.

SOLVENT
- KEEP conceptually.
- Prefer in-process usage.
- Investigate whether existing Solvent Go code can be reused directly.
- Do not automatically preserve the Solvent HTTP server just because it already exists.
- Preserve its epistemic/authority semantics and canonical data model.

CONDUCTOR
- REMOVE as a standalone service.
- Replace with the smallest possible in-process operational task model.
- [AMENDMENT] Before replacing Conductor, inventory the actual semantics the POC uses:
    task lifecycle
    dependencies
    blocked/ready behavior
    cancellation
    dead-end semantics
    lineage / REOPEN
    claim/release, if actually used
    activity/history
- Then explicitly classify each as: preserve, simplify, or discard.
- This inventory is more important than preserving Conductor itself.

COORDINATOR
- REMOVE as a standalone service.
- Convert its useful logic into an in-process application/service layer.
- Preserve mechanical orchestration and contract validation.
- Do not allow it to become a scientific reasoning layer.

VERIFIER
- KEEP as a library.
- [AMENDMENT] Use one verifier implementation with two entry points:
    - `argus verify` CLI subcommand: developer/reproduction runs.
    - In-process invocation: trust-chain-bound verification during packet validation.
  Same library underneath. No second verifier implementation. No file-based verifier-to-Coordinator handoff.

DOMAIN PACK
- KEEP.
- This is critical to proving domain portability.
- BM-IST must remain a domain pack/reference implementation.
- [AMENDMENT] Core owns the mechanism (beliefs, evidence, debt, edges, lifecycle, authority gates). Pack owns the vocabulary (claim vocabulary, evidence classes, retirement rules, falsifiers, domain verifier). The evidence-class mechanism remains generic and opaque to Solvent; do not hard-code BM-IST concepts into the core.

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
4.5. AUTHORIZATION MODEL
==================================================

[AMENDMENT] The collapsed process must distinguish agent calls from human calls, even though they share the same HTTP server.

Agent operations (unprivileged):
    get_context
    submit_packet

Human operations (consequential, require authenticated attribution):
    discharge
    promote
    retract
    reopen
    authorize

Requirements:

- Agent requests must not be able to reach human consequential operations merely because everything shares the same process.
- Human attribution must come from authenticated server-side identity, never from a request-body principal_id or equivalent.
- Minimal CSRF / Origin / Host protection for browser-triggered consequential actions.
- Do NOT introduce a full identity system. Token-derived attribution is sufficient for POC.

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

[AMENDMENT] State the target schema organization explicitly: one CockroachDB database with tables organized by logical subsystem (epistemic tables from Solvent schema, work tables from Conductor schema). Namespace prefixes or schema-level separation are acceptable if needed for clarity, but one connection string is the target.

==================================================
6. RECONNAISSANCE AND DECISION GATE
==================================================

[AMENDMENT] Before proposing the plan, inspect the repository thoroughly. Reconnaissance is not optional; it is the foundation of the plan.

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

[AMENDMENT] For each complexity line item in the current POC, classify it as:

    deployment complexity
    API/integration complexity
    inherent domain complexity

This prevents the pivot from merely relocating distributed complexity into a monolith.

[AMENDMENT] Phase 0 of the implementation plan must produce a GO / NO-GO / REVISE decision gate. The gate answer is:

"Can the POC be materially simplified without losing load-bearing semantics?"

If the answer is NO or REVISE, the plan must state what would need to change for a future GO.

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

[AMENDMENT] The portability test is argued by construction, not demonstrated by execution. The plan must list the specific interfaces the domain pack must satisfy, so the argument has concrete form:

    Pack interface:
        GetPackID() string
        GetVersion() string
        Validate() error

    Debt vocabulary:
        Pack declares its own debt item names.

    Evidence classes:
        Pack declares its own evidence class names.

    Retirement rules:
        Pack maps debt items to evidence classes and rules.

[AMENDMENT] Add a null domain-pack CI compile test: the core builds against a minimal stub domain pack that satisfies the interface with empty/null implementations. This proves at compile time that the core does not import BM-IST vocabulary. This stub is a test fixture, not a second domain pack.

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
10.5. MIGRATION AND TEST DISPOSAL
==================================================

[AMENDMENT] State the migration model explicitly:

- In-place pivot on a branch (preferred for POC).
- Rollback criterion: if Phase 0 gate is REVISE, abandon the branch.
- If Phase 1+ encounters unrecoverable issues, revert to main.

[AMENDMENT] Classify existing tests:

    PRESERVE — tests architectural invariants (e.g., belief lifecycle, debt promotion gate, edge semantics)
    REWRITE  — same behavior, new location (e.g., coordinator compiler test moved to application package)
    DELETE   — tests old service topology (e.g., HTTP handler tests for separate Solvent/Conductor services)
    ADD      — tests new boundaries (e.g., import-direction enforcement, agent/human auth separation)

This is important because otherwise the old architecture will be unconsciously preserved through its test suite.

==================================================
11. IMPLEMENTATION PLAN FORMAT
==================================================

Write a new document:

    docs/plan8/PIVOT_POC_IMPLEMENTATION_PLAN.md

The plan must be lean. No more than 4 phases. Each phase must be executable and have a clear output.

==================================================
12. PHASE DESIGN
==================================================

[AMENDMENT] Use exactly 4 phases:

Phase 0 — Reconnaissance + GO/NO-GO
    Complexity classification (deployment vs API vs domain)
    Work-semantics inventory (Conductor actual usage)
    Solvent kernel reuse analysis
    Component disposition table
    Boundary enforcement design
    Auth model design
    Output: GO / NO-GO / REVISE decision with evidence

Phase 1 — Core collapse
    Create cmd/argus/main.go (serve, mcp, verify, reset subcommands)
    Extract Solvent kernel into internal/epistemic/
    Extract Conductor domain types into internal/work/
    Extract Coordinator logic into internal/application/
    CockroachDB schema (merge or namespace Conductor tables)
    Package boundary enforcement (depguard + compile-time test)
    Auth boundary (agent vs human)

Phase 2 — Integration / end-to-end flow
    Wire MCP adapter to application layer (no HTTP)
    Wire Trust UI to application layer (no HTTP)
    Wire verifier to application layer (in-process)
    EBP packet submission via in-process persistence
    Context assembly via in-process reads
    Human decision flow (discharge, promote, retract)

Phase 3 — Verification, portability proof + cleanup
    End-to-end demo walkthrough
    Null domain-pack CI compile test
    Simplification ledger verification
    Delete old service topology code + tests
    Taskfile.yml (task dev, task test, task verify)

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
- [AMENDMENT] add new packages, interfaces, or abstractions under the guise of simplification (e.g., a "clean architecture" refactor that introduces more layers)

Do not solve hypothetical future scale problems.

==================================================
14. CRITICAL DESIGN QUESTION
==================================================

The plan must directly answer:

"What architectural property would be lost if Coordinator and Conductor stopped being separate services?"

The honest answer: nothing architectural — only operational isolation (independent deployment, ownership, trust-scoping), which no POC needs. Their conceptual roles survive as internal modules while their process boundaries are removed.

"What makes Solvent worth preserving as a distinct conceptual subsystem?"

Solvent is worth preserving because it is the only subsystem whose state is authoritative rather than operational — the authority model is the thesis; everything else is machinery. It is also the candidate reusable product, hence the kernel discipline.

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

[AMENDMENT] Replace vague criteria with a concrete before/after simplification ledger:

CURRENT POC:
    Processes:           ~7 (+ OpenCode)
    Ports:               ~6
    Databases:           1-2 (CockroachDB + SQLite for Conductor)
    Go modules:          3 (oracle, trust-ui, reference-loop)
    Credential sets:     multiple
    Startup steps:       orchestrated readiness chain
    Internal HTTP calls: many
    Cross-process handoffs: yes

PIVOTED POC:
    Processes:           1 (+ OpenCode)
    Ports:               1 (HTTP+UI; MCP is stdio)
    Databases:           1 (CockroachDB)
    Go modules:          1 (+ Solvent module import)
    Credential sets:     1
    Startup steps:       `task dev`
    Internal HTTP calls: 0
    Cross-process handoffs: none

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
- the simplification ledger above is achieved
- package boundaries are mechanically enforced (depguard + compile-time test)
- the core compiles against a null domain pack (portability proof)

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
