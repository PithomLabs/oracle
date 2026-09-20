ARGUS POC — PHASE 3 IMPLEMENTATION
EBP RESEARCH PACKET v1 CONTRACT

You are now implementing ONLY Phase 3 of the frozen ARGUS POC.

AUTHORITATIVE PLAN:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

PHASE STATUS:

    Phase 0 — PASS
    Phase 1 — PASS
    Phase 2 — PASS

PHASE REPORTS:

    oracle/plan/PHASE0_RECONNAISSANCE.md
    oracle/plan/freeze_reconciliation.md
    oracle/plan/PHASE2_DOMAIN_PACK.md

Do NOT begin Phase 4.
Do NOT implement the Physics Verifier.
Do NOT implement the Coordinator.
Do NOT implement Trust UI.
Do NOT modify Solvent.
Do NOT modify Conductor.
Do NOT modify reference-loop.

==================================================
PHASE 3 OBJECTIVE
==================================================

Define EBP Research Packet v1 as the canonical machine-readable interchange
between research-capable agents and the future ARGUS Coordinator.

The packet is:

    agent output
        ↓
    Research Packet
        ↓
    validation + deterministic reference resolution
        ↓
    Coordinator

The packet is NOT:

    canonical epistemic state
    authority
    truth
    policy evaluation
    execution authorization

Agents produce candidate material.

Coordinator/Solvent own canonical state transitions.

==================================================
IMPLEMENTATION SCOPE
==================================================

Create ONLY:

    oracle/packet/v1/schema.json
    oracle/packet/v1/types.go
    oracle/packet/v1/resolve.go
    oracle/packet/v1/validate.go
    oracle/packet/v1/validate_test.go

    oracle/packet/v1/examples/valid_packet.json

    oracle/packet/v1/examples/malformed_missing_packet_id.json
    oracle/packet/v1/examples/malformed_empty_claim.json
    oracle/packet/v1/examples/malformed_unknown_claim_type.json
    oracle/packet/v1/examples/malformed_free_text_edge.json
    oracle/packet/v1/examples/malformed_bad_reference.json
    oracle/packet/v1/examples/malformed_bad_debt.json

Do NOT create:

    verifier/
    coordinator/
    run/
    neutrality/
    adversarial/
    evidence/

==================================================
PACK DEPENDENCY
==================================================

Phase 2 already provides:

    oracle/domain-pack/

Use the existing generic PackRegistry and BM-IST Pack.

Do NOT duplicate the Domain Pack vocabulary in packet code.

Packet validation must resolve:

    pack_ref
        ↓
    PackRegistry
        ↓
    Pack

The packet validator consumes the pack's declared:

    claim_types
    debt_vocabulary
    evidence_classes
    falsifiers

Do not hardcode BM-IST values into the packet package.

==================================================
CANONICAL PACKET MODEL
==================================================

Use the frozen EBP Research Packet v1 shape, generalized where required:

    packet_id
    pack_ref
    scenario_id
    run_id
    task_ref
    project_ref
    scenario_ref
    scope
    corpus_ref
    agent
    beliefs
    evidence
    debts
    edges
    tasks

The exact shape must remain compatible with the architectural contract already
established for EBP Research Packet v1.

Do NOT introduce:

    intent
    unknown object
    conflict object
    falsifier object as a separate canonical object
    confidence score
    truth score
    authorization state

Represent challenge/falsification through edges, evidence, debt, and packet
content rather than inventing parallel epistemic primitives.

==================================================
OBJECT OWNERSHIP
==================================================

Packet objects map to:

    beliefs  → Solvent belief
    evidence → Solvent evidence
    debts    → Solvent opaque debt identifiers
    edges    → Solvent belief_edge
    tasks    → Conductor operational tasks

The packet MUST NOT contain commands that directly mutate Solvent or Conductor.

Examples of prohibited packet fields:

    "promote": true
    "authorize": true
    "retire_debt": true
    "execute": true
    "status": "promoted"

==================================================
LOCAL IDS AND CANONICAL REFERENCES
==================================================

This is a critical architectural boundary.

Every packet object that can be referenced MUST have:

    local_id

Use two explicit reference forms:

    local:<id>

for objects declared inside the same packet.

And:

    canonical:belief:<uuid>

for already-existing canonical Solvent beliefs.

Do NOT support free-text claim matching.

Do NOT infer identity from claim text.

Do NOT fuzzy-match beliefs.

Do NOT normalize/paraphrase claims to resolve references.

Reference resolution must be deterministic.

Examples:

    local:b1
    local:e1
    canonical:belief:550e8400-e29b-41d4-a716-446655440000

A malformed reference must fail validation.

==================================================
PACKET STRUCTURE
==================================================

Use a versioned schema.

Recommended top-level structure:

{
  "schema_version": "ebp-research-packet/v1",
  "role": "work | adversarial",
  "packet_id": "uuid",
  "run_id": "uuid-or-run-reference",
  "task_ref": "string-or-reference",
  "project_ref": "string-or-reference",
  "scenario_ref": "string-or-reference",
  "scope": "string",
  "corpus_ref": "string",
  "agent": {},
  "beliefs": [],
  "evidence": [],
  "debts": [],
  "edges": [],
  "tasks": []
}

Preserve compatibility with the previously locked packet contract.

Do not add fields merely because they seem useful.

==================================================
BELIEFS
==================================================

Each belief must have:

    local_id
    claim
    claim_type

Optional:

    debt
    metadata only if already supported by the frozen contract

Rules:

    claim is non-empty
    claim_type must be declared by the resolved Domain Pack
    local_id must be unique within the packet
    debt items must be declared by the Domain Pack

A packet's debt is additive candidate debt.

The Coordinator later unions:

    pack.initial_debt
    +
    packet-supplied debt

Do NOT allow a packet to override/remove pack initial debt.

==================================================
EVIDENCE
==================================================

Each evidence object must have:

    local_id
    belief_ref
    provenance_class
    content_sha256

Optional fields may include:

    source_url
    artifact_ref
    snapshot
    source_observed_at
    artifact_type

Rules:

    belief_ref must resolve deterministically
    content_sha256 is required
    content_sha256 must be valid hexadecimal SHA-256
    provenance_class must be declared by the Domain Pack

Important:

A packet may REFER to a reproducible artifact.

It must NOT embed an artifact and thereby make it authoritative.

Future Phase 4/6 will resolve verifier artifacts through ArtifactRegistry.

Phase 3 only validates the packet reference structure.

==================================================
DEBTS
==================================================

The frozen architecture treats debt as opaque to Solvent.

Packet validation may verify:

    debt identifier exists in the selected Domain Pack vocabulary

but must NOT attach semantic interpretation to the string.

No generic packet code may hardcode:

    needMap
    needInvariant
    etc.

Domain Pack supplies those identifiers.

==================================================
EDGES
==================================================

Each edge must contain:

    local_id
    from_ref
    to_ref
    kind

where:

    kind ∈ {
        derives,
        contradicts
    }

Reference rules:

    from_ref and to_ref must use:
        local:<id>
        or
        canonical:belief:<uuid>

Reject:

    free-text claim strings
    fuzzy references
    inferred references

The packet validator must detect:

    malformed reference
    unknown local_id
    invalid canonical UUID
    invalid kind

It is acceptable for an edge to connect:

    local belief → canonical belief
    canonical belief → local belief
    local belief → local belief
    canonical belief → canonical belief

Do NOT enforce graph semantics beyond structural validation in Phase 3.

Do NOT create belief edges here.

==================================================
TASKS
==================================================

Each task must contain:

    local_id
    title
    description

Optional:

    governance_ref

Tasks submitted through the packet are always compiled later as:

    proposed

The packet MUST NOT specify:

    active
    accepted
    completed
    cancelled
    blocked

Do not allow packet content to grant operational authority.

Task scope governance remains a Coordinator concern in Phase 6.

==================================================
AGENT OBJECT
==================================================

Agent metadata may identify:

    agent id
    model/provider
    role
    run metadata

It must NOT contain:

    authority
    approval
    promotion
    truth designation

The packet's:

    role

must be one of:

    work
    adversarial

No third role is introduced in Phase 3.

==================================================
SCHEMA VALIDATION
==================================================

Create:

    oracle/packet/v1/schema.json

The JSON Schema must enforce structural constraints such as:

    required fields
    types
    enums
    non-empty strings
    local_id presence
    unique local IDs within objects where mechanically expressible
    packet_id format
    content_sha256 structure

However, cross-object semantic checks such as:

    local reference exists
    pack_ref resolves
    debt belongs to selected Pack
    evidence belongs to referenced belief

should be implemented in Go validation code.

Do not pretend JSON Schema alone enforces graph semantics it cannot reliably
enforce.

==================================================
GO TYPES
==================================================

Create:

    oracle/packet/v1/types.go

Use typed Go structures matching schema.json.

Avoid:

    interface{}
    map[string]any everywhere
    loosely typed packet objects

Use explicit types.

Do not import Solvent kernel types.

Do not import Conductor domain types.

Packet package must remain independently testable.

==================================================
REFERENCE RESOLVER
==================================================

Create:

    oracle/packet/v1/resolve.go

Implement deterministic reference resolution.

Suggested abstractions:

    type ReferenceKind string

    const (
        ReferenceLocal ReferenceKind = "local"
        ReferenceCanonical ReferenceKind = "canonical"
    )

    type Reference struct {
        Kind string
        ID   string
    }

You may choose a cleaner structure.

Required behavior:

    local:b1
        → packet-local belief b1

    canonical:belief:<uuid>
        → canonical belief reference

The resolver must NOT call Solvent during Phase 3.

It only resolves packet-local references and validates canonical reference shape.

The future Coordinator will resolve canonical IDs against actual Solvent state.

This is important:

    Phase 3 resolves REFERENCE FORM.

    Phase 6 resolves CANONICAL STATE.

Do not collapse those responsibilities.

==================================================
RESOLVED PACKET
==================================================

Define a useful internal representation for the future Coordinator, such as:

    ResolvedPacket

It may contain normalized internal references, but must not perform external
state mutation.

A valid packet should be deterministic to transform into this representation.

==================================================
VALIDATION ORDER
==================================================

Use a fail-closed validation sequence:

    1. Decode JSON
    2. Validate JSON Schema
    3. Validate packet_id
    4. Resolve pack_ref
    5. Validate belief structure
    6. Validate evidence structure
    7. Validate debt vocabulary membership
    8. Validate local IDs
    9. Validate all references
    10. Validate edges
    11. Validate tasks
    12. Produce validated packet

Return typed errors where practical.

Do not partially accept a malformed packet.

==================================================
ERRORS
==================================================

Define clear validation errors.

Examples:

    ErrMissingPacketID
    ErrInvalidPacketID
    ErrUnknownPack
    ErrMissingLocalID
    ErrDuplicateLocalID
    ErrInvalidReference
    ErrUnknownLocalReference
    ErrInvalidCanonicalReference
    ErrUnknownClaimType
    ErrUnknownDebt
    ErrInvalidEvidenceClass
    ErrInvalidEdgeKind

You may consolidate errors if a clean existing project convention suggests it.

Do not overengineer the error taxonomy.

==================================================
NEGATIVE TESTS
==================================================

At minimum test:

    valid packet → PASS

    missing packet_id → FAIL

    malformed packet_id → FAIL

    empty claim → FAIL

    unknown claim_type → FAIL

    unknown debt → FAIL

    unknown evidence class → FAIL

    missing local_id → FAIL

    duplicate local_id → FAIL

    malformed local reference → FAIL

    unknown local reference → FAIL

    malformed canonical belief UUID → FAIL

    free-text edge target → FAIL

    invalid edge kind → FAIL

    evidence references nonexistent belief → FAIL

    task missing title → FAIL

    task with forbidden operational status → FAIL

    unsupported role → FAIL

    packet references unknown pack → FAIL

Also verify:

    canonical belief references are structurally accepted

even though actual Solvent lookup occurs only later.

==================================================
EXAMPLE PACKET
==================================================

Create:

    oracle/packet/v1/examples/valid_packet.json

It should contain:

    one work belief
    one evidence item
    representative debt
    at least one edge
    one proposed task
    valid agent metadata
    valid pack reference

Make the example demonstrate local references correctly.

Also create malformed examples matching the required negative tests.

Do not use placeholders that pass superficial validation but do not represent
a realistic packet.

==================================================
DOMAIN NEUTRALITY
==================================================

The generic packet code MUST NOT contain BM-IST-specific identifiers.

Allowed:

    oracle/domain-pack/bmist/v1/
    packet/v1/examples/

Generic packet code must not contain:

    Fisher
    rigidity
    quantum
    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview

The packet package learns vocabulary from the Domain Pack.

==================================================
TESTING
==================================================

Run:

    go test ./...

    go test -race ./...

    go vet ./...

If practical, also validate:

    schema.json itself is valid JSON Schema

Do not add a new external schema-validation dependency unless actually needed.

Prefer the smallest dependency set consistent with the project.

If an implementation choice requires a dependency, explain why before adding
it.

==================================================
REPOSITORY SAFETY
==================================================

Do NOT modify:

    solvent-main
    conductor
    oracle/reference-loop

Do not modify:

    oracle/domain-pack/

unless a strictly necessary compatibility correction is discovered and is
already implied by the frozen Phase 2 contract. If such a correction is
needed, STOP and report it rather than silently changing Phase 2.

Do not modify the frozen ARGUS plan.

Do not reset unrelated work.

Do not remove pre-existing untracked files.

Do not commit unless explicitly instructed.

==================================================
PHASE 3 DELIVERABLE
==================================================

Create:

    oracle/plan/PHASE3_EBP_PACKET.md

Include:

    Phase 3 status
    packet schema summary
    reference model
    validation rules
    error handling
    test results
    example packet
    malformed packet coverage
    domain-neutrality check
    repositories changed
    repositories unchanged
    acceptance results
    deviations, if any

Only report what was actually implemented and tested.

==================================================
PHASE 3 ACCEPTANCE
==================================================

PASS only when:

    1. schema.json exists and is valid.
    2. Go types match schema.
    3. pack_ref resolves through PackRegistry.
    4. local/canonical references are deterministic.
    5. free-text references are rejected.
    6. claim types come from Domain Pack.
    7. debt vocabulary comes from Domain Pack.
    8. evidence classes come from Domain Pack.
    9. tasks cannot specify authoritative operational states.
    10. malformed packets are rejected.
    11. valid packet example passes.
    12. go test ./... passes.
    13. go test -race ./... passes.
    14. go vet ./... is clean.
    15. generic packet code contains no BM-IST vocabulary.
    16. Solvent remains unchanged.
    17. Conductor remains unchanged.
    18. reference-loop remains unchanged.
    19. No Phase 4+ implementation exists.
    20. PHASE3_EBP_PACKET.md is complete.

==================================================
FINAL RESPONSE
==================================================

Return only:

    PHASE 3 RESULT: PASS / STOP

    Schema:
    ...

    Validation:
    ...

    Reference resolution:
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
    oracle/plan/PHASE3_EBP_PACKET.md

    Next phase:
    PHASE 4 — PHYSICS VERIFIER RUNNER

Do NOT begin Phase 4 in this run.