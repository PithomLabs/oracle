ARGUS POC — PHASE 4 IMPLEMENTATION
PHYSICS VERIFIER RUNNER

You are now implementing ONLY Phase 4 of the frozen ARGUS POC.

AUTHORITATIVE PLAN:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

COMPLETED:

    Phase 0 — PASS
    Phase 1 — PASS
    Phase 2 — PASS
    Phase 3 — PASS

Phase 3 deliverable:

    oracle/plan/PHASE3_EBP_PACKET.md

Do NOT begin Phase 5.
Do NOT implement the Corpus Manifest.
Do NOT implement the Coordinator.
Do NOT implement Trust UI.
Do NOT implement Work/Adversarial agents.
Do NOT modify Solvent, except for explicitly forbidden work: NO Solvent changes
are permitted in this phase.
Do NOT modify Conductor.
Do NOT modify reference-loop.

==================================================
PHASE 4 OBJECTIVE
==================================================

Implement the deterministic BM-IST Physics Verifier as the FIRST Domain-specific
verifier running above the generic ARGUS substrate.

The verifier must:

    accept explicitly defined BM-IST proof inputs
    perform deterministic checks
    produce reproducible verification artifacts
    distinguish machine-verified from human-attested steps
    register only trusted verifier output
    make no LLM calls
    never mutate Solvent or Conductor

The verifier is NOT a general-purpose CAS.

The verifier is NOT a scientific truth engine.

The verifier demonstrates bounded, reproducible verification.

==================================================
IMPLEMENTATION SCOPE
==================================================

Create ONLY:

    oracle/verifier/registry.go
    oracle/verifier/registry_test.go

    oracle/verifier/physics/v1/verifier.go
    oracle/verifier/physics/v1/verifier_test.go
    oracle/verifier/physics/v1/artifact.go
    oracle/verifier/physics/v1/artifact_test.go
    oracle/verifier/physics/v1/proof.go
    oracle/verifier/physics/v1/proof_test.go

You may create tightly scoped supporting test fixtures under the verifier
package if necessary.

Do NOT create Phase 5+ packages.

==================================================
TRUST BOUNDARY
==================================================

The verifier package owns trusted artifact registration.

Required architecture:

    Physics Verifier
        ↓
    registerTrusted(...)
        ↓
    ArtifactRegistry
        ↑
    Coordinator later reads only

The registry MUST NOT allow arbitrary external callers to register artifacts.

Use:

    registerTrusted(...)

as an unexported/package-private method.

The Coordinator, when implemented later, must only have:

    Resolve(...)
    VerifyHash(...)

Do not export registration.

This is a hard trust boundary.

==================================================
ARTIFACT REGISTRY
==================================================

Implement:

    oracle/verifier/registry.go

Shape may follow:

    type ArtifactRegistry struct {
        mu        sync.RWMutex
        artifacts map[string]VerificationArtifact
    }

Required operations:

    NewArtifactRegistry()

    registerTrusted(artifact)

    Resolve(ctx, evidenceRef)

    VerifyHash(ctx, evidenceRef, expectedSHA256)

Behavior:

    Register invalid artifact → error
    Duplicate evidence reference → error
    Unknown evidence reference → error
    Hash mismatch → error
    Correct artifact + correct hash → success

The registry is:

    in-memory
    process-local
    non-persistent

Do NOT add:

    database
    SQLite
    Redis
    DSSE
    SLSA
    signing
    external artifact service

==================================================
VERIFICATION ARTIFACT
==================================================

Create:

    oracle/verifier/physics/v1/artifact.go

Use the frozen artifact contract:

    RunID
    VerifierVersion
    VerifierHash
    ClaimHash
    InputHash
    ArtifactHash
    Timestamp
    Steps
    Result
    EvidenceRef

Step contains:

    Index
    Description
    Input
    Rule
    Output
    Verified
    CheckType

CheckType must be:

    machine_verified
    human_attested

Result must be:

    confirmed
    refuted
    inconclusive

Do not introduce additional result semantics.

==================================================
HASH DETERMINISM
==================================================

ArtifactHash must be:

    SHA-256 of canonical serialized artifact content

EXCLUDING:

    ArtifactHash itself
    Timestamp

Therefore:

    same inputs
        →
    identical artifact hash

Timestamp remains metadata but does not affect ArtifactHash.

Test reproducibility across at least 100 repeated runs.

Do not use map iteration order or other nondeterministic serialization.

If canonical JSON serialization is needed, implement a deterministic representation.

==================================================
VERIFIER SCOPE
==================================================

Implement EXACTLY four bounded proof obligations.

1. MACHINE VERIFIED

    Δ√ρ/√ρ
      =
    ½(Δρ/ρ)
      −
    ¼(|∇ρ|²/ρ²)

2. MACHINE VERIFIED

    coefficient matching:

    a(ρ) = κ²/(8mρ)

3. HUMAN ATTESTED

    |∇ρ|² coefficient consistency

4. HUMAN ATTESTED

    b′(ρ) = 0

Do NOT expand the mathematical scope.

Do NOT attempt to create a general variational calculus engine.

Do NOT implement arbitrary symbolic integration.

Do NOT implement a full CAS.

==================================================
SYMBOLIC ENGINE — STEP 1
==================================================

Use a deliberately small symbolic representation.

Supported expression forms:

    Var
    Const
    Add
    Mul
    Pow
    Div
    Sqrt
    Grad
    Laplacian

Implement only the transformations required to verify the specific identity.

Required capabilities:

    constant folding
    basic power rules
    multiplication/distribution where necessary
    derivative/gradient handling necessary for the identity
    deterministic simplification

The implementation must be bounded to the proof obligation.

Avoid generic "solve anything" abstractions.

==================================================
STEP 1 REQUIREMENT
==================================================

The verifier must actually symbolically construct and check the identity.

Do NOT merely compare two hardcoded strings.

Do NOT hardcode the final expected expression and call it symbolic verification.

The test should demonstrate:

    symbolic input
        ↓
    symbolic transformation
        ↓
    deterministic simplification
        ↓
    equality check

A known-good identity must produce:

    Verified = true
    CheckType = machine_verified

A deliberately altered identity must produce:

    Verified = false
    Result = refuted or inconclusive according to the verifier contract

Do not invent unsupported mathematical semantics.

==================================================
SYMBOLIC ENGINE — STEP 2
==================================================

For coefficient matching, implement only the symbolic manipulation necessary
to verify:

    a(ρ) = κ²/(8mρ)

Do NOT build a general Euler-Lagrange system.

Represent the required coefficient relationship as a bounded symbolic
verification problem.

Again:

    symbolic input
        ↓
    deterministic transformation
        ↓
    equality / coefficient comparison

No hardcoded "true" result.

==================================================
STEPS 3–4
==================================================

These remain HUMAN_ATTESTED.

Implement deterministic assertions over the expected proof-obligation inputs.

They may use explicit expected-value constants because the POC deliberately
does not claim that these two steps were independently derived by the verifier.

Record:

    CheckType = human_attested

The artifact must make this distinction obvious.

IMPORTANT:

An overall artifact result of "confirmed" MUST NOT imply that every step was
machine-derived.

The artifact must preserve the per-step distinction.

==================================================
RESULT SEMANTICS
==================================================

Define clear result aggregation.

Required principles:

    any machine-verified contradiction
        → refuted

    unresolved/inconclusive machine verification
        → inconclusive

    all required steps satisfied/attested
        → confirmed

Do NOT infer scientific truth from missing information.

Document the aggregation rule in code comments and tests.

==================================================
VERIFIER METADATA
==================================================

Populate:

    VerifierVersion
    VerifierHash
    ClaimHash
    InputHash
    RunID
    Timestamp

Do not fabricate values.

VerifierHash must be deterministic for the verifier build/source representation
used by the POC, or use an explicitly documented POC method if executable
self-hashing is impractical.

ClaimHash:

    SHA-256 of canonical claim text/input representation.

InputHash:

    SHA-256 of canonical verifier input.

EvidenceRef:

    stable identifier for the generated artifact.

Keep all metadata reproducible except Timestamp.

==================================================
ARTIFACT REGISTRATION
==================================================

The Physics Verifier runner must:

    construct artifact
    calculate ArtifactHash
    validate artifact
    register via registerTrusted(...)

No other package should have an exported registration path.

Tests MUST verify the trust boundary.

At minimum:

    forged artifact object not produced by verifier
        cannot simply call an exported Register method

    registry empty + unknown ref
        → reject

    registered trusted artifact + correct hash
        → accept

    registered artifact + wrong hash
        → reject

==================================================
TEST MATRIX
==================================================

Write deterministic tests for:

REGISTRY:

    valid artifact registration
    duplicate reference rejection
    unknown reference rejection
    correct hash accepted
    wrong hash rejected
    concurrent reads safe
    concurrent trusted registrations safe
    race detector clean

ARTIFACT:

    canonical hashing deterministic
    timestamp excluded from hash
    artifact hash reproducible across 100 runs
    malformed artifact rejected
    invalid result rejected
    invalid check type rejected

STEP 1:

    known-good identity → verified
    deliberately modified identity → rejected/refuted/inconclusive
    output is deterministic

STEP 2:

    known-good coefficient relationship → verified
    deliberately modified coefficient → rejected/refuted/inconclusive
    output is deterministic

STEP 3:

    expected human-attested assertion → pass
    malformed assertion → fail

STEP 4:

    expected human-attested assertion → pass
    malformed assertion → fail

AGGREGATION:

    all required steps satisfied → confirmed

    any refuted step → refuted

    unresolved machine check → inconclusive

TRUST BOUNDARY:

    registration API not exported
    Coordinator-style external package cannot call registerTrusted

==================================================
DEPENDENCY CONSTRAINT
==================================================

Keep the POC dependency footprint minimal.

Before adding any external symbolic algebra dependency, STOP and assess whether
the bounded verifier can be implemented with the current lightweight symbolic
representation.

The frozen architecture explicitly prefers:

    bounded symbolic verification

over:

    general-purpose CAS dependency

Do not pull in a large CAS merely for convenience.

==================================================
DOMAIN BOUNDARY
==================================================

Physics-specific semantics are allowed ONLY under:

    oracle/verifier/physics/v1/

Generic registry code:

    oracle/verifier/registry.go

must contain ZERO BM-IST/physics vocabulary.

The generic registry must operate on generic verification artifacts.

Do not leak:

    Fisher
    rigidity
    quantum
    rho
    Hamiltonian

into generic registry code.

==================================================
REPOSITORY SAFETY
==================================================

Do NOT modify:

    solvent-main
    conductor
    oracle/reference-loop

Do NOT modify:

    oracle/domain-pack/

unless a strict compile compatibility issue exists. If one does exist,
STOP and report it rather than silently changing Phase 2.

Do not modify the frozen execution plan.

Do not remove unrelated user files.

Do not reset unrelated work.

Do not commit without explicit instruction.

==================================================
RUN VERIFICATION
==================================================

From:

    /home/chaschel/Documents/go/oracle

Run:

    go test ./...
    go test -race ./...
    go vet ./...

All must pass.

Also verify that no Phase 5+ package has been created.

==================================================
PHASE 4 DELIVERABLE
==================================================

Create:

    oracle/plan/PHASE4_PHYSICS_VERIFIER.md

Include:

    Phase 4 status
    verifier scope
    proof obligations
    symbolic approach
    human-attested boundary
    artifact structure
    hashing approach
    ArtifactRegistry design
    trust boundary
    test results
    repositories changed
    repositories unchanged
    limitations
    deviations

Be precise about what is machine-verified versus human-attested.

Do not describe the verifier as proving the entire Fisher-rigidity theorem.

The accurate claim is:

    "The POC machine-verifies two bounded symbolic proof obligations and
     records two additional obligations as human-attested."

==================================================
PHASE 4 ACCEPTANCE
==================================================

PASS only when:

    1. Four proof obligations are implemented.
    2. Steps 1–2 perform real symbolic checks.
    3. Steps 3–4 are explicitly human_attested.
    4. Artifact output is deterministic.
    5. Artifact hash is reproducible.
    6. Registry is process-local and in-memory.
    7. Trusted registration is package-private.
    8. Forged/unknown/wrong-hash artifacts are rejected.
    9. No LLM calls exist in verifier.
    10. No external CAS is introduced unnecessarily.
    11. Generic registry contains no physics vocabulary.
    12. go test ./... passes.
    13. go test -race ./... passes.
    14. go vet ./... passes.
    15. Solvent unchanged.
    16. Conductor unchanged.
    17. reference-loop unchanged.
    18. No Phase 5+ implementation exists.
    19. PHASE4_PHYSICS_VERIFIER.md is complete.

==================================================
FINAL RESPONSE
==================================================

Return only:

    PHASE 4 RESULT: PASS / STOP

    Verifier:
    ...

    Machine-verified:
    ...

    Human-attested:
    ...

    Artifact registry:
    ...

    Determinism:
    ...

    Tests:
    ...

    Trust boundary:
    ...

    Files created:
    ...

    Other repositories modified:
    ...

    Deliverable:
    oracle/plan/PHASE4_PHYSICS_VERIFIER.md

    Next phase:
    PHASE 5 — CORPUS MANIFEST

Do NOT begin Phase 5 in this run.