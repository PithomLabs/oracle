ARGUS POC — PHASE 5 IMPLEMENTATION
CORPUS MANIFEST

You are now implementing ONLY Phase 5 of the frozen ARGUS POC.

AUTHORITATIVE PLAN:

    oracle/plan/ARGUS_POC_EXECUTION_PLAN_v1.1.md

COMPLETED:

    Phase 0 — PASS
    Phase 1 — PASS
    Phase 2 — PASS
    Phase 3 — PASS
    Phase 4 — PASS

Phase reports:

    oracle/plan/PHASE0_RECONNAISSANCE.md
    oracle/plan/freeze_reconciliation.md
    oracle/plan/PHASE2_DOMAIN_PACK.md
    oracle/plan/PHASE3_EBP_PACKET.md
    oracle/plan/PHASE4_PHYSICS_VERIFIER.md

Do NOT begin Phase 6.
Do NOT implement the Coordinator.
Do NOT implement Coordinator HTTP.
Do NOT implement Trust UI.
Do NOT implement the Work agent.
Do NOT implement the Adversarial stub.
Do NOT modify Solvent.
Do NOT modify Conductor.
Do NOT modify reference-loop.

==================================================
PHASE 5 OBJECTIVE
==================================================

Create a deterministic, hash-pinned manifest for the BM-IST POC corpus.

The manifest establishes:

    what artifacts belong to the corpus
    where they are located
    what their content hash is
    what provenance class they have
    how they can be located deterministically

The manifest is provenance metadata.

It is NOT:

    epistemic authority
    source reputation
    truth
    a second evidence ledger
    a knowledge graph

==================================================
IMPLEMENTATION SCOPE
==================================================

Create ONLY:

    oracle/corpus/manifest.json
    oracle/corpus/manifest.go
    oracle/corpus/manifest_test.go

Do not create Phase 6+ code.

==================================================
CORPUS CONTENT
==================================================

Include at minimum:

    1. BM-IST POC seed:
       oracle/docs/BMIST_POC_SEED.md

    2. Qwen-derived artifact used by the POC.

The exact Qwen artifact path MUST be discovered from the actual repository
during implementation.

Do NOT invent a filename or path.

If the artifact cannot be located:

    STOP and report the missing artifact.

Do not silently omit it.

==================================================
PROVENANCE
==================================================

Use provenance carefully.

For the BM-IST seed:

    provenance_class = operator_asserted

For the Qwen-derived artifact:

    provenance_class = agent_derived

Important:

    agent_derived does NOT confer authority.

The Qwen artifact is contextual research material, not source truth, unless
explicitly admitted elsewhere by the frozen architecture.

Only actual verifier-produced artifacts may later enter Solvent as:

    reproducible_artifact

Do not classify the Qwen artifact as:

    reproducible_artifact

==================================================
MANIFEST STRUCTURE
==================================================

Create:

    oracle/corpus/manifest.json

Use a versioned structure such as:

{
  "manifest_version": "1.0.0",
  "corpus_id": "bmist-poc-seed-v1",
  "artifacts": [
    {
      "id": "bmist-seed-001",
      "path": "docs/BMIST_POC_SEED.md",
      "sha256": "...",
      "provenance_class": "operator_asserted",
      "description": "...",
      "stable_locator": "oracle/docs/BMIST_POC_SEED.md"
    }
  ]
}

Include a manifest timestamp only if it is explicitly treated as metadata and
does NOT enter the artifact hash model.

Do not introduce unnecessary manifest fields.

==================================================
ARTIFACT RECORD RULES
==================================================

Each artifact MUST have:

    id
    path
    sha256
    provenance_class
    description
    stable_locator

Rules:

    id must be unique
    path must point to an existing repository file
    sha256 must be valid SHA-256 hex
    provenance_class must be one of the supported corpus provenance values
    stable_locator must be deterministic and repository-relative

Use repository-relative paths.

Do not use absolute local filesystem paths.

==================================================
HASHING
==================================================

Compute SHA-256 directly from the exact file bytes.

Do NOT hash:

    parsed JSON
    normalized Markdown
    trimmed whitespace
    line endings after transformation

Hash the actual file contents.

The manifest must record the exact content hash.

==================================================
MANIFEST LOADER
==================================================

Implement:

    oracle/corpus/manifest.go

Provide typed Go structures and validation.

Required behavior:

    Load manifest
    Validate structure
    Verify referenced paths exist
    Recompute hashes
    Compare recorded hashes
    Return typed errors on mismatch

Suggested API:

    Load(path string) (*Manifest, error)

    (m *Manifest) Validate() error

    (m *Manifest) Verify(root string) error

Use project conventions if a cleaner API is appropriate.

==================================================
VALIDATION RULES
==================================================

Validate:

    manifest_version is supported
    corpus_id is non-empty
    artifacts is non-empty
    artifact IDs are unique
    paths are non-empty
    stable locators are non-empty
    SHA-256 values are exactly valid hexadecimal
    provenance classes are supported
    referenced files exist
    recorded hash matches actual file bytes

Reject:

    missing file
    hash mismatch
    duplicate artifact ID
    malformed SHA-256
    empty path
    empty stable locator
    unsupported provenance class

==================================================
PATH SAFETY
==================================================

The manifest is repository-relative.

Reject paths that escape the repository root.

At minimum reject:

    absolute paths
    ../ traversal

Do not permit:

    ../../secret-file
    /etc/...
    arbitrary host filesystem paths

Stable locators must remain repository-relative.

==================================================
TEST FIXTURES
==================================================

Do NOT modify real corpus files merely to create tests.

For unit tests, create temporary test files under a temporary directory.

Test:

    valid manifest
    valid hash
    missing file
    hash mismatch
    duplicate artifact ID
    malformed SHA-256
    unsupported provenance class
    empty artifact list
    absolute path rejection
    path traversal rejection

Also test:

    multiple artifacts
    stable repository-relative locators

==================================================
ACTUAL CORPUS VERIFICATION
==================================================

After implementing the loader/tests:

1. Discover the actual BM-IST seed path.

2. Discover the actual Qwen-derived artifact path.

3. Compute hashes from exact file bytes.

4. Populate manifest.json.

5. Load manifest using the Go loader.

6. Verify the manifest against the repository.

The checked-in manifest MUST correspond to the current repository contents.

Do not hardcode hashes copied from an earlier session if the files differ.

==================================================
DETERMINISM
==================================================

The same repository contents must produce the same manifest artifact hashes.

No nondeterministic values should be used for identity.

Artifact order in manifest.json should be deterministic.

Use stable ordering, preferably lexicographic by artifact ID or path.

Do not include runtime-specific filesystem metadata.

==================================================
DOMAIN NEUTRALITY
==================================================

The generic corpus loader must contain ZERO BM-IST-specific semantics.

This is allowed:

    oracle/corpus/manifest.json
    actual corpus entries

This is NOT allowed in generic Go code:

    Fisher-specific logic
    Qwen-specific handling
    physics-specific conditionals

The loader must operate on generic corpus artifact records.

==================================================
RELATION TO PHASE 4
==================================================

Do not register verifier outputs into this corpus manifest merely because
they are produced by the Physics Verifier.

Phase 4 verifier artifacts have their own verification/artifact path.

The corpus manifest describes the POC research corpus.

It does not replace ArtifactRegistry.

Keep the two concepts separate:

    Corpus Manifest
        = what research inputs/artifacts comprise the corpus

    ArtifactRegistry
        = trusted verifier outputs admitted during the running POC

==================================================
TESTING
==================================================

Run:

    cd /home/chaschel/Documents/go/oracle

    go test ./...
    go test -race ./...
    go vet ./...

Also manually verify:

    manifest.json loads
    all referenced files exist
    every recorded hash matches

Do not claim PASS unless actually verified.

==================================================
REPOSITORY SAFETY
==================================================

Do NOT modify:

    solvent-main
    conductor
    oracle/reference-loop
    oracle/domain-pack
    oracle/packet
    oracle/verifier

Do not modify the frozen execution plan.

Do not remove unrelated files.

Do not reset unrelated work.

Do not commit unless explicitly instructed.

==================================================
PHASE 5 DELIVERABLE
==================================================

Create:

    oracle/plan/PHASE5_CORPUS_MANIFEST.md

Include:

    Phase 5 status
    corpus identity
    artifacts included
    exact repository-relative paths
    provenance classes
    recorded SHA-256 values
    verification results
    path-safety behavior
    test results
    repositories changed
    repositories unchanged
    limitations
    deviations

Explicitly state:

    Qwen artifact is agent-derived/contextual unless otherwise admitted.
    Corpus membership does not confer epistemic authority.
    Verifier artifacts remain governed by ArtifactRegistry.

==================================================
PHASE 5 ACCEPTANCE
==================================================

PASS only when:

    1. manifest.json exists.
    2. BM-IST seed is included.
    3. Qwen-derived artifact is included, if present and locatable.
    4. Every artifact has a correct SHA-256.
    5. Provenance classes are correct.
    6. Repository-relative paths are used.
    7. Path traversal is rejected.
    8. Missing files are detected.
    9. Hash mismatches are detected.
    10. Duplicate IDs are rejected.
    11. go test ./... passes.
    12. go test -race ./... passes.
    13. go vet ./... passes.
    14. Generic corpus code contains no BM-IST-specific logic.
    15. Solvent unchanged.
    16. Conductor unchanged.
    17. reference-loop unchanged.
    18. Domain Pack unchanged.
    19. Packet unchanged.
    20. Verifier unchanged.
    21. No Phase 6+ implementation exists.
    22. PHASE5_CORPUS_MANIFEST.md is complete.

==================================================
FINAL RESPONSE
==================================================

Return only:

    PHASE 5 RESULT: PASS / STOP

    Corpus:
    ...

    Artifacts:
    ...

    Provenance:
    ...

    Hash verification:
    ...

    Path safety:
    ...

    Tests:
    ...

    Files created:
    ...

    Other repositories modified:
    ...

    Deliverable:
    oracle/plan/PHASE5_CORPUS_MANIFEST.md

    Next phase:
    PHASE 6 — DETERMINISTIC COORDINATOR

Do NOT begin Phase 6 in this run.