Use this prompt:

```text
IMPLEMENTATION PLAN ONLY — PHYSICS VERIFIER POC AS FIRST DOMAIN PACK

You are the senior implementation planner for the Solvent + Conductor repository.

Do NOT implement code yet.
Do NOT modify files.
Do NOT commit, tag, push, or rewrite Git history.

Produce a repository-grounded, phase-by-phase implementation plan for the
Physics Verifier POC described below.

The repository is the source of truth for current implementation details.
Inspect the actual Conductor, Solvent, Coordinator, EBP, test, migration,
and existing BM-IST code before proposing file paths or APIs.

==================================================
ARCHITECTURAL OBJECTIVE
==================================================

The POC must demonstrate that BM-IST physics can operate as a Domain Pack
on a domain-neutral trust-verification substrate.

The key claim is NOT:

    "Solvent understands physics."

The claim is:

    "The same generic verification substrate can govern BM-IST research
     without physics-specific logic entering Solvent or Conductor."

Locked separation:

    Domain Pack
        ↓
    Trust Verifier / Coordinator
        ↓
    Solvent = epistemic + authority state
    Conductor = operational workflow
        ↓
    Executor
        ↓
    External SOR / effect truth

Locked invariants:

    DEBT IS OPAQUE TO SOLVENT
    WORKFLOW STATE != EPISTEMIC TRUTH
    EPISTEMIC STATE != EXTERNAL EFFECT TRUTH
    CONSEQUENTIAL EPISTEMIC WRITE != ORDINARY OPERATIONAL WRITE

BM-IST is the first Domain Pack, not the kernel.

==================================================
POC DELIVERABLES
==================================================

The final plan must cover these deliverables:

A. Prerequisites
1. Plan 11.1 executed and kernel re-frozen.
2. New freeze hash pinned.
3. Freeze reconciliation/decision record completed, including deferred H5 items.

B. Domain-Pack Contract
4. Physics Domain Pack specification:
   - claim types;
   - evidence classes;
   - EBP six debt vocabulary;
   - discharge semantics;
   - falsifier/challenge declarations;
   - consequential-action declarations;
   - structural gate map;
   - explicit human-gated transitions.

5. EBP Research Packet v1:
   - belief;
   - evidence;
   - debt;
   - edge;
   - task;
   - local vs canonical references;
   - deterministic ownership mapping to Solvent/Conductor.

Do NOT invent additional canonical packet objects unless repository evidence
requires them.

C. Deterministic Physics Verification
6. Physics Verifier Runner:
   - symbolic/algebraic/toy-model checks;
   - deterministic inputs/outputs;
   - reproducible artifacts;
   - verifier/version/input metadata;
   - artifact hashes;
   - no LLM authority.

D. Corpus
7. Corpus manifest:
   - BM-IST seed;
   - Qwen artifact;
   - hashes;
   - provenance class;
   - stable locators/anchors.

E. Coordinator
8. Deterministic Go Coordinator:
   - packet validation;
   - mandatory ebpInitialDebt attachment;
   - negative regression test proving the initial debt cannot silently disappear;
   - compilation to Solvent;
   - compilation to Conductor;
   - packet_id idempotency;
   - Solvent-before-Conductor ordering;
   - Conductor task projection using governance_ref;
   - human decision endpoint;
   - human decisions represented using operator_asserted provenance where
     compatible with the existing Solvent schema.

Coordinator must NOT:
   - reason scientifically;
   - infer truth;
   - score confidence;
   - automatically promote;
   - automatically retire debt;
   - adjudicate disputes.

F. Human Adjudication
9. Human adjudication loop:
   - review;
   - retire/open debt;
   - merge where explicitly supported;
   - retract;
   - consequential branch pruning through Solvent semantics;
   - action_intent involved where the decision changes consequential state.

G. End-to-End Research Run
10. One WORK agent + one ADVERSARIAL agent:
    - fresh processes;
    - distinct model families;
    - bounded task scope;
    - packet-only system interface;
    - agents cannot directly mutate protected state.

11. Pre-registered acceptance criteria:
    - malformed packet rejected;
    - incomplete adversarial coverage detected;
    - zero agent-caused protected-state mutations;
    - human can reconstruct belief state from persisted projections alone;
    - promotion gate works;
    - consequential action intent requires promoted state;
    - contradiction/retraction invalidates downstream authority/intents.

H. Domain-Neutrality Proof
12. Neutrality artifact:
    - no physics/domain vocabulary in Solvent;
    - no physics/domain vocabulary in Conductor;
    - pack interacts with substrate only through generic contracts/API;
    - automated conformance test;
    - repository-wide static checks as an additional guard.

Important:
    grep is not sufficient proof.
The conformance test is the primary proof mechanism.

I. Evidence Package
13. Final POC evidence package:
    - corpus manifest;
    - packet examples;
    - verifier artifacts;
    - Solvent state transitions;
    - Conductor task/dependency projections;
    - adversarial results;
    - human adjudication decisions;
    - provenance trail;
    - consequential-action demonstration;
    - failure cases;
    - architecture/documentation report.

==================================================
EXPLICIT NON-GOALS
==================================================

Do NOT include implementation of these in this POC unless the repository
already requires them for existing functionality:

- cryptographic evidence signing / DSSE;
- SLSA/in-toto integration;
- W3C PROV export;
- second Domain Pack;
- probabilistic inference;
- source reputation engine;
- workflow primitives inside Solvent;
- epistemic primitives inside Conductor;
- pack-interface freeze;
- new domain-specific tables in Solvent;
- LLM inside the Go Coordinator.

These remain future roadmap items.

==================================================
PHASED PLAN
==================================================

Structure the plan into explicit phases.

For every phase provide:

    Objective
    Preconditions
    Repository areas/files to inspect
    Files likely to change
    API/schema impact
    Tests
    Acceptance criteria
    Rollback/escape condition

Do not leave design choices to the future coding agent.

Recommended phase structure:

PHASE 0 — Repository reconnaissance
    Map actual current implementation.
    Identify existing BM-IST verifier code.
    Identify Coordinator entry points.
    Identify exact Conductor/Solvent integration points.
    Identify existing packet/EBP artifacts.
    Identify Plan 11.1/freeze dependency.

PHASE 1 — Freeze reconciliation
    Execute/verify Plan 11.1 prerequisites.
    Record exact kernel commit/hash and decision record.
    Establish the immutable baseline for the POC.

PHASE 2 — Domain Pack specification
    Define the declarative Physics Domain Pack.
    Explicitly separate:
        domain semantics
        generic verifier mechanics
        Solvent structural constraints
        Conductor workflow projection.

PHASE 3 — Packet contract
    Implement/reconcile EBP Research Packet v1.
    Define:
        packet-local IDs
        canonical IDs
        reference resolution
        validation rules
        malformed-packet behavior
        idempotency semantics.

PHASE 4 — Physics verifier runner
    Isolate the actual BM-IST verification logic.
    Define deterministic artifact outputs.
    Define hashes and run metadata.
    Make verifier outputs consumable as evidence without giving the verifier
    authority over promotion.

PHASE 5 — Corpus manifest
    Hash-pin the supplied research corpus.
    Define stable source locators.
    Make provenance explicit.

PHASE 6 — Coordinator compiler
    Build deterministic translation:
        packet → Solvent
        packet → Conductor
    Implement:
        initial debt enforcement
        idempotency
        ordering
        governance_ref projection
        failure/refusal handling.

PHASE 7 — Human adjudication
    Implement only the minimum human-decision path needed for the POC.
    Ensure consequential pruning/retraction traverses Solvent authority semantics.
    No autonomous adjudication.

PHASE 8 — Work + adversarial run
    Run one WORK and one ADVERSARIAL fresh agent process.
    Capture packet outputs.
    Validate mechanically.
    Persist only through the Coordinator.

PHASE 9 — End-to-end consequential-state demonstration
    Demonstrate at least:

        candidate claim
        → evidence
        → unresolved debt
        → verification
        → human review
        → debt retirement
        → promotion
        → action intent

    And separately:

        promoted belief
        → contradictory evidence
        → retraction
        → downstream invalidation
        → intent invalidation/refusal

PHASE 10 — Domain-neutrality proof
    Prove substrate does not depend on BM-IST vocabulary.
    Add conformance tests.
    Run static repository scans.
    Demonstrate that physics semantics are entirely outside Solvent/Conductor.

PHASE 11 — POC evidence and report
    Produce the complete auditable evidence package.
    Record:
        successful paths
        failed paths
        adversarial findings
        unresolved gaps
        architectural observations.

PHASE 12 — Final adversarial review
    Attempt to falsify the actual POC claim:

        "BM-IST is a Domain Pack, not a kernel customization."

    Look specifically for:
        hidden physics assumptions in Coordinator;
        domain leakage into Conductor;
        domain leakage into Solvent;
        packet fields carrying hidden physics semantics;
        automatic scientific adjudication;
        accidental agent authority;
        duplicate sources of truth.

PHASE 13 — Acceptance / exit
    Define the exact conditions under which the POC is considered successful,
    partially successful, or failed.

==================================================
IMPORTANT REPOSITORY CONSTRAINTS
==================================================

Before proposing implementation changes:

- inspect the actual Solvent schema;
- inspect actual Conductor schema;
- inspect actual Coordinator code;
- inspect current BM-IST code;
- inspect existing EBP packet/skill artifacts;
- inspect Plan 11.1 and freeze records.

Do not invent a third canonical database.

Do not add epistemic concepts to Conductor.

Do not add BM-IST concepts to Solvent.

Do not silently alter frozen kernel invariants.

If the current repository cannot support a proposed requirement without changing
a frozen invariant, identify that explicitly as a Growth Gate rather than
silently designing around it.

==================================================
SUCCESS CLAIM
==================================================

The POC must prove:

    BM-IST semantics live entirely in the Domain Pack/verifier layer.

    Solvent remains domain-agnostic.

    Conductor remains domain-agnostic.

    The Coordinator is deterministic and does not become a reasoning engine.

    Agents produce candidate research material, not authority.

    Humans control consequential epistemic transitions.

    The same substrate can represent verification, debt, promotion,
    retraction, and consequential authorization without knowing anything
    about physics.

Do NOT claim that one domain proves general domain-agnosticism.

The honest claim after this POC is:

    "The substrate has successfully hosted one complete Domain Pack."

A second domain is required to empirically demonstrate cross-domain
generality.

==================================================
FINAL OUTPUT
==================================================

Return ONLY the implementation plan.

Include:

1. Phase list and dependency graph.
2. Exact repository components affected by each phase.
3. Exact deliverables per phase.
4. Required schema/API changes, if any.
5. Tests and acceptance criteria.
6. Explicit non-goals.
7. Risks and failure modes.
8. Final POC acceptance criteria.
9. Estimated implementation order:
       sequential vs parallel work.
10. A final "definition of done".

Do not implement anything.
Do not commit.
Do not push.
```


## prompt

can we include that UI considerations in the implementation plan prompt or make it separate? brief but concise reply only

Include them **in the implementation plan**, but as a **separate UI/UX phase and boundary**, not mixed into the kernel/domain-pack implementation.

Recommended structure:

```text
Phase 0–10   Trust Verifier / Domain Pack / Solvent / Coordinator
Phase 11     UI architecture + information model
Phase 12     POC UI implementation
Phase 13     adversarial UX review + acceptance
```

The plan should lock:

* **Oracle** → epistemic state/query surface
* **Sphinx** → consequential-action gate
* **Passage Ledger** → debt/obligations
* **Challenge Room** → adversarial verification
* **Decision Room** → human consequential decisions
* **Refusal Archive** → refusals/failed crossings
* **Chronicle** → historical activity
* **Conductor remains unchanged** and is merely an optional operational backend/projection.

This is preferable to a separate plan because the UI must be derived from the **same state/ownership contracts** established by the POC, while keeping implementation concerns isolated.


