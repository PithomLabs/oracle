IMPLEMENTATION PLAN ONLY — PHYSICS VERIFIER POC + DOMAIN-AGNOSTIC TRUST UI

You are the senior implementation planner for the Solvent + Conductor ecosystem.

DO NOT implement code yet.
DO NOT modify files.
DO NOT commit.
DO NOT tag.
DO NOT push.
DO NOT rewrite Git history.

Produce a repository-grounded, phase-by-phase implementation plan for the
Physics Verifier POC and its dedicated Trust UI.

The repository is the source of truth for current implementation details.
Inspect the actual Solvent, Conductor, Coordinator, EBP, BM-IST verifier,
tests, schemas, migrations, and current UI capabilities before proposing
APIs, file paths, or abstractions.

==================================================
ARCHITECTURAL OBJECTIVE
==================================================

Demonstrate that BM-IST can run as the FIRST DOMAIN PACK on a
domain-agnostic trust-verification substrate.

The claim is NOT:

    Solvent understands physics.

The claim is:

    The generic trust substrate can govern BM-IST research
    without physics-specific logic entering Solvent or Conductor.

Locked architecture:

    Domain Pack
        ↓
    Trust Verifier / Coordinator
        ↓
    Solvent = epistemic + authority state
    Conductor = operational workflow
        ↓
    Executor / Ferryman
        ↓
    External SOR = effect truth

Locked invariants:

    DEBT IS OPAQUE TO SOLVENT
    WORKFLOW STATE != EPISTEMIC TRUTH
    EPISTEMIC STATE != EXTERNAL EFFECT TRUTH
    CONSEQUENTIAL EPISTEMIC WRITE != ORDINARY OPERATIONAL WRITE

BM-IST is the first Domain Pack, not the kernel.

==================================================
CRITICAL UI ARCHITECTURE DECISION
==================================================

DO NOT MODIFY CONDUCTOR JUST TO OBTAIN A UI.

Conductor remains frozen/minimal and is treated as an optional operational
backend/projection.

Build a separate Trust UI.

The Trust UI consumes canonical state through Solvent/Coordinator APIs.

The UI must NOT create a second source of truth.

The UI must NOT embed physics logic.

The UI must NOT add epistemic primitives to Conductor.

Conductor's existing UI may remain available as an operational view, but it
is NOT the primary trust-system UI.

==================================================
PRODUCT SYMBOLISM
==================================================

Use Greek mythology as the information architecture.

Core metaphor:

    ORACLE
        What does the system know?

    SPHINX
        Is the system entitled to act?

    FERRYMAN
        What actually happened in the external world?

    CHRONICLER
        What happened along the way?

    CONDUCTOR
        Who is doing the work?

The metaphor must clarify architecture, not obscure it.

Use the mythology as UX terminology while keeping underlying technical names
explicit and inspectable.

==================================================
CORE UI SURFACES
==================================================

The initial POC Trust UI should contain six primary surfaces:

1. ORACLE
2. SPHINX / RIDDLE
3. PASSAGE LEDGER
4. CHALLENGE ROOM
5. DECISION ROOM
6. REFUSAL ARCHIVE

BRIEFING SCROLL and CHRONICLE should initially be views/subsections rather than
additional top-level products unless repository evidence shows they need to be
separate.

==================================================
ORACLE — EPISTEMIC QUERY SURFACE
==================================================

Oracle answers ONLY from canonical Solvent state.

Primary question:

    What does the system currently know/claim?

A selected belief should expose:

    CLAIM
    STATUS
    EVIDENCE
    PROVENANCE
    DEBT / OPEN OBLIGATIONS
    CHALLENGES
    RELATIONS
    HISTORY
    UNKNOWN / UNRESOLVED

Rules:

    Oracle never invents facts.
    Oracle never upgrades a claim.
    Oracle never infers authority.
    Missing information must remain UNKNOWN.

A user or agent should be able to "consult the Oracle" before acting.

==================================================
SPHINX — CONSEQUENTIAL ACTION GATE
==================================================

Primary question:

    May this request pass?

Represent a consequential request as a RIDDLE.

Show:

    WHAT?
        exact action

    WHY?
        supporting belief

    EVIDENCE?
        evidence/provenance

    AUTHORITY?
        requesting/approving principal

    RESULT:
        PASS
        REFUSE
        HUMAN REVIEW

The Sphinx must surface structural gate failures rather than generate vague
policy explanations.

The Sphinx is the presentation of existing Solvent authority semantics, not
a second policy engine.

==================================================
PASSAGE LEDGER — DEBT / OBLIGATIONS
==================================================

Represent debt as an explicit passage requirement.

Example:

    Belief: BM-IST B1

    PASSAGE REQUIREMENTS
      ✓ needMap
      ✓ needInvariant
      ✓ needToyCheck
      ○ needNullModel
      ○ needObstruction
      ○ needFaithfulnessReview

    Promotion: BLOCKED

Each obligation should expose:
    current state
    discharge evidence
    reviewer/actor
    timestamp
    history

The UI must not reinterpret debt strings.
It renders the Domain Pack's semantics.

==================================================
CHALLENGE ROOM — ADVERSARIAL VERIFICATION
==================================================

Represent the Sphinx's second role: attack the claim.

Show:

    CURRENT CLAIM
        ↓
    ATTACKS
        ↓
    counterexample
    contradiction
    missing evidence
    alternative explanation
    unknown

Possible outcomes:

    NO ISSUE FOUND
    NEW DEBT
    FALSIFIER FOUND
    CONTRADICTION FOUND
    UNKNOWN

A promoted belief must remain challengeable.

Do not turn this into an autonomous truth engine.

==================================================
DECISION ROOM — HUMAN CONSEQUENCE CONTROL
==================================================

Provide explicit human handling for consequential epistemic transitions:

    PROMOTE
    RETIRE DEBT
    RETRACT
    REOPEN
    AUTHORIZE
    REFUSE

Every consequential transition must display:

    current state
    requested transition
    evidence
    open obligations
    downstream consequences
    acting human/principal
    resulting Solvent transition

No silent consequential state mutation.

==================================================
REFUSAL ARCHIVE — "THE BONES"
==================================================

Expose refusal_log / refusal evidence as a first-class audit surface.

Show:

    time
    requester
    action
    supporting belief
    failed gate
    evidence/context
    decision

Treat refusal as useful audit evidence, not as an error to hide.

==================================================
BRIEFING SCROLL
==================================================

The Trust UI must provide a pre-decision briefing view:

    current beliefs
    promoted vs entered
    open debts
    active challenges
    pending human decisions
    active action intents
    consequences at stake

This should be a projection of canonical state.

Do NOT create a separate dashboard database.

==================================================
CHRONICLE
==================================================

Provide a historical activity view that can correlate:

    evidence ingestion
    verification
    debt retirement
    promotion
    retraction
    authorization
    refusal
    execution

Conductor activity may be displayed here, but remains non-authoritative.

The Chronicle is historical context, not epistemic truth.

==================================================
DOMAIN PACK VS UI
==================================================

The UI must be domain-neutral.

For BM-IST it renders:
    BM-IST claim types
    BM-IST debt vocabulary
    BM-IST challenges

without embedding those semantics in the UI kernel.

The UI should consume a Domain Pack descriptor/configuration or generic
Coordinator projection where appropriate.

The UI must NOT contain:
    physics equations
    BM-IST-specific decision logic
    hardcoded EBP debt semantics
    domain-specific promotion rules

==================================================
POC DELIVERABLES
==================================================

A. Prerequisites
1. Plan 11.1 executed and kernel re-frozen.
2. New freeze hash pinned.
3. Freeze reconciliation/decision record completed.

B. Domain-Pack Contract
4. Physics Domain Pack specification:
   - claim types
   - evidence classes
   - EBP six debt vocabulary
   - discharge semantics
   - falsifier/challenge declarations
   - consequential-action declarations
   - structural gate map
   - human-gated transitions

C. Packet Contract
5. EBP Research Packet v1:
   - belief
   - evidence
   - debt
   - edge
   - task
   - local/canonical references
   - ownership mapping to Solvent/Conductor

D. Deterministic Physics Verification
6. Physics Verifier Runner:
   - symbolic/algebraic/toy checks
   - deterministic artifacts
   - reproducible outputs
   - verifier/version/input metadata
   - content hashes
   - no LLM authority

E. Corpus
7. Corpus manifest:
   - BM-IST seed
   - Qwen artifact
   - hashes
   - provenance classes
   - stable locators

F. Coordinator
8. Deterministic Go Coordinator:
   - packet validation
   - mandatory ebpInitialDebt attachment
   - negative regression test
   - Solvent compilation
   - Conductor compilation
   - packet_id idempotency
   - Solvent-before-Conductor ordering
   - governance_ref projection
   - human decision endpoint
   - operator_asserted provenance where applicable

Coordinator must NOT:
   - reason scientifically
   - infer truth
   - score confidence
   - auto-promote
   - auto-retire debt
   - adjudicate

G. Human Adjudication
9. Human loop:
   - review
   - retire/open debt
   - retract
   - branch pruning
   - consequential action-intent path

H. End-to-End Research Run
10. One WORK agent + one ADVERSARIAL agent:
    - fresh processes
    - distinct model families
    - bounded scope
    - packet-only system interface
    - no protected-state mutation by agents

11. Pre-registered acceptance criteria:
    - malformed packet rejected
    - incomplete adversarial coverage detected
    - zero agent-caused protected-state mutations
    - human can reconstruct belief state from persisted projections
    - promotion gate works
    - action intent requires promotion
    - retraction invalidates dependent authority/intents

I. Domain-Neutrality Proof
12. Neutrality artifact:
    - no physics/domain vocabulary in Solvent
    - no physics/domain vocabulary in Conductor
    - pack interacts only through generic schema/API
    - automated pack conformance test
    - static repository scan as secondary evidence

J. Trust UI
13. Separate Trust UI implementation:
    - Oracle
    - Sphinx/Riddle
    - Passage Ledger
    - Challenge Room
    - Decision Room
    - Refusal Archive
    - Briefing Scroll view
    - Chronicle view

UI MUST consume canonical Solvent/Coordinator state.

UI MUST NOT modify Conductor.

==================================================
PHASE STRUCTURE
==================================================

For EVERY phase provide:

    Objective
    Preconditions
    Repository areas/files to inspect
    Files/components likely to change
    Data/API impact
    UI impact
    Tests
    Acceptance criteria
    Rollback/escape condition
    Dependencies on previous/future phases

PHASE 0 — REPOSITORY RECONNAISSANCE
    Inspect current:
    - Solvent schema/API
    - Conductor schema/API/UI
    - Coordinator
    - EBP artifacts
    - BM-IST verifier
    - existing trust/research UI
    - Plan 11.1/freeze state

    Explicitly identify which current Conductor capabilities can be reused
    without changing Conductor.

PHASE 1 — FREEZE RECONCILIATION
    Establish immutable baseline.

PHASE 2 — DOMAIN PACK SPECIFICATION
    Define declarative Physics Pack.

PHASE 3 — EBP PACKET CONTRACT
    Define/implement packet validation and mapping.

PHASE 4 — PHYSICS VERIFIER RUNNER
    Isolate deterministic BM-IST verification and artifact production.

PHASE 5 — CORPUS MANIFEST
    Hash-pin the research corpus.

PHASE 6 — DETERMINISTIC COORDINATOR
    Implement packet → Solvent + Conductor compilation and idempotency.

PHASE 7 — HUMAN ADJUDICATION
    Implement minimum consequential decision path.

PHASE 8 — TRUST UI INFORMATION ARCHITECTURE
    Define the UI's:
    - information model
    - navigation
    - state projections
    - domain-pack rendering contract
    - Oracle/Sphinx/Passage/Challenge/Decision/Refusal semantics
    - empty/unknown/error states
    - read-vs-write boundaries

    Do NOT implement visual polish yet.

PHASE 9 — TRUST UI IMPLEMENTATION
    Build the six primary surfaces plus Briefing/Chronicle views.

    IMPORTANT:
    - do not modify Conductor;
    - UI writes only through approved Coordinator/Solvent APIs;
    - no direct database writes from UI;
    - no domain logic in presentation code.

PHASE 10 — WORK + ADVERSARIAL RUN
    Execute the first complete BM-IST research cycle.

PHASE 11 — CONSEQUENTIAL STATE DEMONSTRATION
    Demonstrate:

        claim
        → evidence
        → debt
        → verification
        → human decision
        → promotion
        → authorized intent

    Then:

        promoted belief
        → contradiction
        → retraction
        → dependent invalidation
        → intent invalidation/refusal

    Every transition must be visible in the Trust UI.

PHASE 12 — DOMAIN-NEUTRALITY PROOF
    Prove BM-IST is data/configuration, not substrate logic.

    Test that replacing BM-IST-specific pack data does not require
    Solvent/Conductor/UI kernel changes.

PHASE 13 — ADVERSARIAL SYSTEM REVIEW
    Attempt to falsify:

        "BM-IST is a Domain Pack, not a kernel customization."

    Look for:
        physics leakage
        UI semantic leakage
        hidden domain state
        agent authority
        duplicated state
        direct UI database writes
        Conductor modifications
        Coordinator reasoning
        incorrect authority projections

PHASE 14 — POC EVIDENCE PACKAGE
    Produce:
        corpus manifest
        packet examples
        verifier artifacts
        Solvent state traces
        Conductor task projections
        UI screenshots/flows
        adversarial results
        human decisions
        refusal evidence
        consequential transition traces
        failure cases
        architecture report

PHASE 15 — FINAL ACCEPTANCE
    Decide:
        PASS
        PARTIAL PASS
        FAIL

    Based on explicit acceptance criteria, not aesthetics.

==================================================
UI DESIGN PRINCIPLES
==================================================

Lock these principles into the plan:

1. Oracle never guesses.
2. Sphinx never silently passes.
3. Refusal is first-class evidence.
4. Unknown is a valid state.
5. Every consequential state change is visible.
6. Human authority is explicit.
7. The UI presents canonical state; it does not invent state.
8. Domain Pack semantics are rendered dynamically.
9. Conductor remains unchanged.
10. The mythology is a semantic aid, not a replacement for technical labels.

Technical names should remain available in UI details/tooltips.

==================================================
EXPLICIT NON-GOALS
==================================================

Do NOT implement in this POC:

- modifications to Conductor core or UI;
- cryptographic DSSE/SLSA attestation;
- W3C PROV export;
- second Domain Pack;
- probabilistic inference;
- source reputation engine;
- ontology/knowledge graph engine;
- workflow primitives in Solvent;
- epistemic primitives in Conductor;
- embedded LLM reasoning inside Solvent;
- embedded LLM reasoning inside Coordinator;
- autonomous promotion;
- autonomous adjudication;
- pack-interface freeze.

These remain future roadmap items.

==================================================
SUCCESS CLAIM
==================================================

The POC proves:

    BM-IST semantics live in the Domain Pack/verifier layer.

    Solvent remains domain-agnostic.

    Conductor remains unchanged and optional as an operational backend.

    The Trust UI is a separate presentation layer.

    Coordinator is deterministic and non-reasoning.

    Agents produce candidate research material, not authority.

    Humans control consequential epistemic transitions.

    The same substrate supports:
        verification
        evidence
        debt
        promotion
        retraction
        consequential authorization

Do NOT claim that one Domain Pack proves complete domain agnosticism.

The honest claim is:

    "The substrate successfully hosts one complete Domain Pack through
     a domain-neutral verification and trust UI."

A second Domain Pack is required for empirical cross-domain validation.

==================================================
FINAL OUTPUT
==================================================

Return ONLY the implementation plan.

Include:

1. Executive decision
2. Phase dependency graph
3. Exact deliverables by phase
4. Repository components/files affected
5. API/schema impact
6. Trust UI information architecture
7. Trust UI implementation plan
8. Test and acceptance matrix
9. Explicit non-goals
10. Risks/failure modes
11. Parallel vs sequential work
12. Definition of Done

Do NOT implement anything.
Do NOT modify files.
Do NOT commit.
Do NOT push.
