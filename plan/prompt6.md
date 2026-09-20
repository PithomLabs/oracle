ARGUS POC — PHASE 2 IMPLEMENTATION
DOMAIN PACK SPECIFICATION

You are now implementing ONLY Phase 2 of the frozen ARGUS POC.

Authoritative plan:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

Phase 0:
    PASS

Phase 1:
    PASS

Phase 1 reconciliation record:

    oracle/plan/freeze_reconciliation.md

DO NOT begin Phase 3.
DO NOT implement the Research Packet.
DO NOT implement the Verifier.
DO NOT implement the Coordinator.
DO NOT implement the Trust UI.
DO NOT implement the Solvent Growth Gate edge endpoint in this phase.

==================================================
PHASE 2 OBJECTIVE
==================================================

Define and implement the BM-IST Domain Pack as a declarative, validated,
versioned artifact consumed by the future Coordinator.

The Domain Pack owns:

    claim types
    evidence classes
    debt vocabulary
    initial debt
    debt retirement rules
    falsifier/challenge vocabulary
    consequential actions
    human-gated transitions

The Domain Pack MUST NOT introduce physics-specific logic into:

    Solvent
    Conductor
    generic Coordinator infrastructure

The Domain Pack is configuration + schema, not an execution engine.

==================================================
IMPLEMENTATION BOUNDARY
==================================================

Create ONLY the Phase 2 components:

    oracle/go.mod
    oracle/domain-pack/registry.go
    oracle/domain-pack/registry_test.go

    oracle/domain-pack/bmist/v1/pack.json
    oracle/domain-pack/bmist/v1/types.go
    oracle/domain-pack/bmist/v1/validate.go
    oracle/domain-pack/bmist/v1/validate_test.go

    oracle/domain-pack/README.md

    oracle/testdata/packs/valid/
    oracle/testdata/packs/invalid/

Do not create Phase 3+ directories:

    packet/
    verifier/
    coordinator/
    run/
    neutrality/
    adversarial/
    evidence/

==================================================
GO MODULE
==================================================

Create:

    oracle/go.mod

with:

    module github.com/PithomLabs/oracle

    go 1.25

Use the minimum dependencies required.

The frozen plan currently specifies:

    github.com/google/uuid v1.6.0

Do not add unnecessary dependencies.

==================================================
BM-IST DOMAIN PACK
==================================================

Create:

    oracle/domain-pack/bmist/v1/pack.json

Use the BM-IST Pack definition from the frozen execution plan.

Required top-level fields:

    pack_id
    version
    name
    description
    claim_types
    evidence_classes
    debt_vocabulary
    initial_debt
    retirement_rules
    falsifiers
    consequential_actions
    human_gated_transitions

The BM-IST debt vocabulary is exactly:

    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview

Do not add extra debt vocabulary.

Do not rename these identifiers.

==================================================
PACK SEMANTICS
==================================================

The pack must declare:

CLAIM TYPES:

    derived
    accommodated
    postulated

EVIDENCE CLASSES:

    reproducible_artifact
    operator_asserted

DEBT:

    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview

INITIAL DEBT:

    exactly the six BM-IST debt identifiers above

RETIREMENT RULES:

    needMap
        → reproducible_artifact:map_check

    needInvariant
        → reproducible_artifact:invariant_check

    needToyCheck
        → reproducible_artifact:toy_model_check

    needNullModel
        → operator_asserted:scope_clarification

    needObstruction
        → reproducible_artifact:obstruction_construction

    needFaithfulnessReview
        → operator_asserted:faithfulness_review

FALSIFIERS:

    counterexample
    contradiction
    missing_evidence
    alternative_explanation

CONSEQUENTIAL ACTION:

    publish_claim

    requires:
        promoted

    gates:
        faithfulness_review

HUMAN-GATED TRANSITIONS:

The frozen structural floor requires:

    faithfulness_review
    scope_clarification
    obstruction_assessment

Additional pack-specific human gates may be declared only if supported by
the frozen plan. Do not invent new ones.

==================================================
PACK TYPES
==================================================

Create:

    oracle/domain-pack/bmist/v1/types.go

Define strongly typed Go structures corresponding to pack.json.

Keep the types declarative.

Do NOT put verifier behavior, Coordinator behavior, or Solvent API calls
inside these types.

Do NOT make Pack types depend on Solvent.

==================================================
PACK VALIDATOR
==================================================

Create:

    oracle/domain-pack/bmist/v1/validate.go

The validator must enforce the frozen Phase 2 rules.

Required structural validation:

1. pack_id exists and is non-empty.

2. version exists and is valid semver.

3. claim_types exists and is non-empty.

4. claim_types contains at least one supported generic claim type:
       derived
       accommodated
       postulated

5. evidence_classes exists and is non-empty.

6. debt_vocabulary exists and is non-empty.

7. initial_debt exists and is a subset of debt_vocabulary.

8. retirement_rules keys are a subset of debt_vocabulary.

9. Every retirement-rule evidence reference refers to a declared
   evidence class.

10. falsifiers exists and is non-empty.

11. consequential_actions are structurally valid.

12. human_gated_transitions contains the structural human-gate floor:

       faithfulness_review
       scope_clarification
       obstruction_assessment

13. No duplicate values in any declared vocabulary array.

14. No empty string values in arrays or required identifiers.

15. Pack version and pack identity are internally consistent.

IMPORTANT:

This validator checks PACK STRUCTURE.

It does NOT prove scientific correctness.

It does NOT evaluate BM-IST physics.

It does NOT evaluate evidence.

It does NOT decide whether debt is actually retired.

It does NOT promote beliefs.

==================================================
STRUCTURAL HUMAN-GATE FLOOR
==================================================

The validator must make the following invariant explicit:

    Domain Pack may ADD human-gated transitions,
    but may not remove the structural human-gate floor.

This is a substrate safety constraint expressed through pack validation.

Test specifically:

    remove faithfulness_review → reject

    remove scope_clarification → reject

    remove obstruction_assessment → reject

==================================================
PACK REGISTRY
==================================================

Create:

    oracle/domain-pack/registry.go
    oracle/domain-pack/registry_test.go

The registry is:

    in-memory
    startup-loaded
    validated before registration
    keyed by pack_id + version
    non-persistent

Suggested shape:

    type PackRegistry struct {
        mu    sync.RWMutex
        packs map[string]*Pack
    }

Required methods:

    NewRegistry()

    Register(pack *Pack)

    Get(packID, version string)

    LoadFromDisk(root string)

You may improve signatures where appropriate, but preserve the architecture.

==================================================
REGISTRY RULES
==================================================

Register():

    validate pack first
    reject invalid pack
    reject duplicate pack_id + version

Get():

    deterministic lookup
    unknown pack → typed error
    wrong version → typed error

LoadFromDisk():

    scan the Domain Pack directory tree

    expected BM-IST location:

        domain-pack/bmist/v1/pack.json

    load pack.json
    decode
    validate
    register

Do not silently accept malformed packs.

For invalid packs:

    return/report an error

Do not silently skip invalid configuration.

The frozen plan says the registry is startup-loaded and validated.

Prefer fail-closed behavior:

    invalid registered pack → startup/configuration failure

==================================================
TEST FIXTURES
==================================================

Create representative fixtures under:

    oracle/testdata/packs/valid/
    oracle/testdata/packs/invalid/

Include at minimum:

VALID:

    valid_bmist_v1.json

INVALID:

    missing_pack_id.json
    missing_version.json
    invalid_version.json
    empty_claim_types.json
    missing_debt_vocabulary.json
    invalid_initial_debt.json
    invalid_retirement_rule.json
    missing_human_gate.json
    duplicate_debt.json

Tests should load these fixtures where useful rather than constructing every
case entirely inline.

==================================================
UNIT TESTS
==================================================

Create comprehensive validation tests.

At minimum test:

    valid BM-IST pack → PASS

    missing pack_id → FAIL

    empty pack_id → FAIL

    missing version → FAIL

    invalid version → FAIL

    empty claim_types → FAIL

    unsupported claim type → FAIL

    missing evidence_classes → FAIL

    missing debt_vocabulary → FAIL

    initial debt contains undeclared item → FAIL

    retirement rule has unknown debt key → FAIL

    retirement rule references unknown evidence class → FAIL

    missing falsifier list → FAIL

    missing faithfulness_review → FAIL

    missing scope_clarification → FAIL

    missing obstruction_assessment → FAIL

    duplicate debt vocabulary → FAIL

    duplicate claim type → FAIL

    empty vocabulary string → FAIL

Registry tests:

    valid pack registers

    invalid pack rejected

    duplicate pack_id + version rejected

    Get existing pack succeeds

    Get unknown pack fails

    Get wrong version fails

    LoadFromDisk loads BM-IST pack

    LoadFromDisk rejects malformed pack

    concurrent Register/Get is race-safe

Run:

    go test ./...

and:

    go test -race ./...

==================================================
DOMAIN NEUTRALITY
==================================================

The generic registry must remain domain-neutral.

The generic registry package must NOT contain:

    Fisher
    quantum
    Hamiltonian
    BM-IST debt names
    physics-specific validation rules

BM-IST-specific semantics belong only under:

    oracle/domain-pack/bmist/v1/

Do not put BM-IST vocabulary into:

    registry.go
    generic registry types
    generic validation helpers

The registry loads generic Pack structures.

==================================================
README
==================================================

Create:

    oracle/domain-pack/README.md

Document:

    what a Domain Pack is
    what BM-IST declares
    what the registry does
    what the registry does NOT do
    relationship to Solvent
    relationship to the future Coordinator
    relationship to EBP

Explicitly state:

    Domain Pack semantics are declarative.
    Solvent does not interpret them.
    The Coordinator consumes them.
    Agents do not gain authority from the Pack.

==================================================
DO NOT IMPLEMENT
==================================================

Do NOT implement:

    Research Packet v1
    Physics Verifier
    ArtifactRegistry
    Coordinator
    Coordinator HTTP API
    Solvent edge endpoint
    Trust UI
    Work agent
    Adversarial stub
    Neutrality scanner
    Evidence package

Do NOT modify:

    Solvent
    Conductor
    reference-loop

The Growth Gate edge endpoint remains planned but NOT implemented in Phase 2.

==================================================
PHASE 2 ACCEPTANCE CRITERIA
==================================================

PASS only when:

    1. oracle/go.mod exists and builds.

    2. BM-IST pack.json exists and validates.

    3. BM-IST vocabulary exactly matches the frozen plan.

    4. Structural human-gate floor is enforced.

    5. Retirement rules reference only declared evidence classes.

    6. Invalid pack fixtures are rejected.

    7. PackRegistry registers and retrieves valid packs.

    8. PackRegistry rejects invalid/duplicate packs.

    9. LoadFromDisk works.

    10. Registry concurrency tests pass.

    11. go test ./... passes.

    12. go test -race ./... passes.

    13. No Solvent changes.

    14. No Conductor changes.

    15. No reference-loop changes.

    16. No Phase 3+ implementation exists.

    17. Generic registry code contains no BM-IST/physics vocabulary.

==================================================
PHASE 2 DELIVERABLE
==================================================

Create:

    oracle/plan/PHASE2_DOMAIN_PACK.md

The report must contain:

    Phase 2 status
    files created
    pack structure
    validator rules
    registry behavior
    test results
    neutrality verification
    repositories changed
    repositories explicitly unchanged
    acceptance results
    any deviations

Do not create a speculative report.

Only report what was actually implemented and tested.

==================================================
GIT SAFETY
==================================================

Do not modify Conductor or Solvent.

Do not rewrite history.

Do not reset unrelated work.

Do not remove pre-existing untracked user files.

At completion:

    git status --short

Report all Oracle changes explicitly.

Do NOT commit unless explicitly instructed after review.

==================================================
FINAL RESPONSE
==================================================

Return only:

    PHASE 2 RESULT: PASS / STOP

    Pack:
    ...

    Registry:
    ...

    Tests:
    ...

    Neutrality:
    ...

    Files created:
    ...

    Other repositories modified:
    ...

    Deliverable:
    oracle/plan/PHASE2_DOMAIN_PACK.md

    Next phase:
    PHASE 3 — EBP RESEARCH PACKET CONTRACT

Do NOT begin Phase 3 in this run.