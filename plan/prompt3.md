UPDATE THE ARGUS EXECUTION PLAN — FINAL PRE-IMPLEMENTATION REVISION

Do NOT implement the ARGUS POC yet.

Modify the current execution plan only. Do not change Solvent, Conductor, reference-loop, or any production code.

SOURCE OF TRUTH:
- Current ARGUS Execution Plan v1.0: execution_plan.md
- Review-consolidated ARGUS POC plan v1.1
- All previously locked architectural decisions in this conversation

OBJECTIVE

Bring the execution plan to final implementation-ready status by incorporating
the four remaining blocking clarifications below.

The resulting document must be internally consistent, explicit enough for a
coding agent to implement without inventing missing infrastructure, and remain
minimal for the POC.

==================================================
LOCKED ARCHITECTURE
==================================================

Project:

    ARGUS — Trust Verification POC

Modules:

    oracle/
        single Go module

    trust-ui/
        separate Go module

    oracle/reference-loop/
        remains untouched, separate go.mod

Coordinator:

    oracle/coordinator/
        in-process Go library

    oracle/coordinator/http/
        thin HTTP wrapper

Trust UI:

    calls Coordinator HTTP API
    never imports Coordinator Go package directly
    never performs direct consequential writes to Solvent

Conductor:

    unchanged

Solvent:

    canonical epistemic + authority state

Agents:

    packet producers only
    never directly mutate protected state

==================================================
BLOCKING FIX 1 — VERIFIER ARTIFACT REGISTRY
==================================================

Problem:

Phase 6 requires:

    reproducible_artifact
        → registered verifier output
        → artifact hash validation

but the execution plan does not define the registry.

Add an explicit POC Artifact Registry.

Design:

    oracle/verifier/registry.go
    oracle/verifier/registry_test.go

The registry is:

    in-memory
    owned by the Coordinator / run harness
    populated only by the trusted Physics Verifier runner
    read by the Coordinator during packet compilation

Suggested abstraction:

    type ArtifactRegistry interface {
        Register(ctx context.Context, artifact VerificationArtifact) error
        Resolve(ctx context.Context, ref string) (VerificationArtifact, error)
        VerifyHash(ref string, expectedSHA256 string) error
    }

Exact API may differ if repository conventions justify it.

Required invariants:

1. A packet cannot admit `reproducible_artifact` evidence merely because
   the packet contains a forged JSON artifact.

2. The Coordinator must resolve the artifact reference through the registry.

3. The resolved artifact hash must equal the packet's declared content_sha256.

4. The artifact must have been produced by the trusted verifier runner
   during the current POC process.

5. Unknown artifact reference → reject.

6. Hash mismatch → reject.

7. Registry contents are not persisted.

8. Restart persistence is explicitly OUT OF SCOPE for the POC.

9. Do NOT add cryptographic signing, DSSE, SLSA, or a database.

Update:

    Phase 4
    Phase 6
    Phase 10
    risks
    tests
    acceptance criteria
    evidence package

Add explicit test:

    forged verifier artifact → rejected
    unknown artifact_ref → rejected
    correct artifact + correct hash → admitted

==================================================
BLOCKING FIX 2 — DECISION RECORD PERSISTENCE
==================================================

Problem:

Phase 6 defines DecisionRecord but does not specify where it persists.

Do NOT create a Coordinator database.

Locked design:

    Solvent state
        +
    Solvent audit/activity
        =
    authoritative persisted decision history

Coordinator DecisionRecord is only:

    request/result envelope
    normalized representation for the UI/API
    not a second system of record

Update the plan so that:

1. Every consequential decision results in the existing Solvent state
   transition where applicable.

2. The consequential action is represented in Solvent's existing audit/activity
   trail.

3. REFUSE must use the existing Solvent refusal-log mechanism where applicable.

4. The Coordinator may return a DecisionRecord to the UI, but does not persist
   a second authoritative copy.

5. Trust UI decision history must be reconstructable from canonical Solvent
   state + audit/refusal records.

6. Do NOT create:
       coordinator.sqlite
       decision database
       decision table
       UI decision store

7. Add tests proving that a Coordinator restart does not create divergent
   decision history, because history is reconstructed from Solvent.

Update:

    Phase 6
    Phase 7
    Phase 8
    Phase 9
    Phase 11
    Phase 14 evidence package

==================================================
BLOCKING FIX 3 — SPHINX AUTHORIZATION CONTEXT
==================================================

Problem:

The Trust UI includes the Sphinx authorization surface, but the current
Coordinator HTTP API has no read endpoint that provides a complete
authorization context.

Do NOT add new authority semantics to Solvent.

Do NOT add a new Solvent kernel endpoint solely for the UI.

Add exactly one thin Coordinator read projection:

    GET /authorization-context/:target_id

Purpose:

    Assemble existing canonical authority state into a UI-readable
    authorization/riddle context.

The projection should expose, as applicable:

    target
    requested action
    supporting belief
    belief status
    evidence summary
    open debt
    authority state
    justification state
    request state
    approval state
    action intent state
    current result:
        PASS
        REFUSE
        HUMAN_REVIEW

Important:

The Coordinator does NOT decide new policy here.

It only projects existing Solvent authority state.

The Sphinx remains a VIEW of canonical state, not a policy engine.

Update:

    Phase 6:
        endpoint implementation

    Phase 8:
        state projection table

    Phase 9:
        Coordinator client
        Sphinx handler

    API contract section

    tests

Required tests:

1. Existing approved authority state → PASS projection.

2. Missing required state → HUMAN_REVIEW or REFUSE according to the
   existing Solvent state.

3. Retracted supporting belief → REFUSE.

4. No new authority decision logic appears in Coordinator.

==================================================
BLOCKING FIX 4 — IDEMPOTENCY LIFETIME
==================================================

Problem:

Phase 6 uses:

    canonical content hash → compilation result

but does not specify lifetime.

Lock:

    in-memory for lifetime of Coordinator process

Therefore:

1. Canonical content hash excludes:
       packet_id
       runtime metadata
       timestamps

2. Identical canonical content under different packet IDs is a duplicate.

3. Duplicate submission during the same Coordinator process is a no-op.

4. The mapping is held in memory only.

5. Coordinator restart clears the mapping.

6. Cross-restart duplicate detection is OUT OF SCOPE for this POC.

7. Do NOT add:
       persistent idempotency DB
       Redis
       SQLite
       distributed locking

8. Concurrent submissions with the same canonical hash must be race-safe.

Add explicit tests:

    same packet twice → one compilation

    same content + different packet_id → one compilation

    concurrent duplicate submissions → one compilation

    Coordinator restart → idempotency cache empty

Document restart behavior as an explicit POC limitation.

==================================================
CONSISTENCY PASS
==================================================

After applying the four changes, inspect the ENTIRE plan for contradictions.

Specifically verify consistency among:

    Phase 4 artifact production
    Phase 5 corpus provenance
    Phase 6 Coordinator
    Phase 7 human adjudication
    Phase 8 UI architecture
    Phase 9 UI implementation
    Phase 10 end-to-end run
    Phase 11 consequential demonstration
    Phase 12 neutrality
    Phase 13 adversarial review
    Phase 14 evidence
    Phase 15 acceptance

Ensure the following remain true:

    Solvent is canonical epistemic/authority state.
    Conductor remains unchanged.
    Coordinator does not reason.
    UI does not create truth.
    Agents do not mutate protected state.
    Verifier output must be independently registered and hash-validated.
    Decision history is reconstructed from Solvent.
    Sphinx is a projection, not a policy engine.
    Idempotency is process-local only.
    No new persistent infrastructure is introduced.

==================================================
UPDATE ACCEPTANCE MATRIX
==================================================

Add explicit acceptance criteria:

Artifact integrity:
    forged artifact rejected
    unknown artifact rejected
    hash mismatch rejected
    valid trusted artifact accepted

Decision persistence:
    decisions recoverable from Solvent audit/activity
    no second decision database exists

Sphinx:
    authorization-context projection works
    projection contains required fields
    no new policy logic

Idempotency:
    canonical duplicate suppressed
    concurrent duplicates race-safe
    different packet IDs with same content deduplicated
    restart clears in-memory cache

==================================================
UPDATE RISKS
==================================================

Add:

1. Artifact registry lost on restart
   → accepted POC limitation

2. Coordinator restart loses idempotency cache
   → accepted POC limitation

3. Solvent audit projection insufficient for a UI field
   → document exact missing canonical source; do not invent persistence

4. Authorization projection accidentally gains policy logic
   → keep Sphinx projection read-only and test state-to-view mapping

==================================================
UPDATE EVIDENCE PACKAGE
==================================================

Phase 14 must include:

    trusted verifier artifact
    artifact registry admission evidence
    forged artifact rejection evidence
    decision reconstruction from Solvent audit
    Sphinx authorization-context response
    idempotency duplicate test result
    concurrent idempotency test result
    restart limitation record

==================================================
FINAL PLAN REQUIREMENTS
==================================================

The revised document must be named:

    ARGUS_POC_EXECUTION_PLAN_v1.1.md

Status:

    PLAN — IMPLEMENTATION READY

Do NOT mark implementation complete.

Do NOT claim tests passed unless they were actually run.

Do NOT modify source code.

Do NOT modify Solvent.

Do NOT modify Conductor.

Do NOT modify reference-loop.

Do NOT introduce new persistent infrastructure.

At the end, provide a concise:

    "Final Implementation Lock"

containing the final locked decisions, including the four changes above.

The resulting plan is the baseline for the subsequent implementation phase.