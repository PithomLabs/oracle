Implement ONLY the three Phase 8 fixes identified in the final plan review.

Do NOT implement the broader hardening backlog.
Do NOT redesign ARGUS, Solvent, Conductor, RCP, MCP, or the authority model.
Do NOT modify the existing Phase 8 workflow beyond what is required below.

The three fixes are:

P0-1  Fresh `task dev` must actually work on a fresh checkout.
P0-2  Physics Verifier output must be consumable by the running ARGUS/Coordinator.
P1-3  Canonical packet references must respect scenario isolation.

==================================================
P0-1 — FIX FRESH `task dev`
==================================================

Current problem:

`task dev` is intended to be safe and repeatable, but `task up` calls
`db:ensure`, and the current `db:ensure` exits when the `fable` database
does not exist without actually invoking the initialization path.

Required behavior:

Fresh checkout:

    task dev

must result in:

    CockroachDB READY
    Solvent READY
    Conductor READY
    Coordinator READY
    Trust UI READY

WITHOUT requiring the developer to run `task fresh` first.

Preserve the important distinction:

    task dev
      = safe / non-destructive / repeatable

    task fresh
      = destructive / reproducible clean initialization

Implement `db:ensure` as:

1. CockroachDB unavailable
   → fail clearly.

2. `fable` database absent
   → invoke Solvent's canonical `db:reset`
   → continue.

3. `fable` database exists and schema is complete
   → continue without reset.

4. Database exists but required schema is incomplete/invalid
   → FAIL clearly and instruct:
       task fresh

Do NOT silently destroy an existing database.

Use Solvent's canonical migration/reset mechanism.
Do NOT duplicate Solvent migration knowledge in Oracle.

Add tests or executable Taskfile checks proving:

- fresh checkout + `task dev` succeeds
- second `task dev` does not reset DB
- existing DB state survives repeated `task dev`
- missing/incomplete schema fails with a clear recovery message
- `task fresh` still performs the explicit destructive reset

==================================================
P0-2 — PHYSICS VERIFIER → ARGUS ARTIFACT HANDOFF
==================================================

Current architectural gap:

`cmd/verifier` creates its own in-memory ArtifactRegistry.

Coordinator also creates a separate in-memory ArtifactRegistry.

Therefore:

    task verify
      → verifier produces artifact

does NOT imply:

    Coordinator / ArtifactReader
      → can resolve that artifact

This must be fixed with the smallest POC mechanism.

Do NOT add:
- a database
- an artifact-management service
- an agent-facing artifact registration tool
- a new authority layer
- cryptographic attestation

Required POC flow:

    task verify
        ↓
    deterministic VerificationArtifact
        ↓
    `.tmp/artifacts/`
        ↓
    Coordinator startup
        ↓
    trusted local artifact registration
        ↓
    ArtifactReader
        ↓
    ARGUS evidence/retirement workflow

Implement a deterministic local artifact handoff.

Recommended shape:

    .tmp/artifacts/g0.json
    .tmp/artifacts/g0.sha256

The exact filenames may differ if a better existing convention exists.

Requirements:

1. `task verify` executes the standalone verifier binary.

2. The verifier writes a canonical VerificationArtifact to the configured
   artifact directory.

3. The artifact must have a deterministic content hash.

4. The artifact output must identify enough information to bind it to:
   - verification identity/version
   - input/fixture identity or hash
   - verification result
   - artifact hash

5. Coordinator startup loads ONLY trusted local artifacts from the
   configured artifact directory.

6. Coordinator registers them through the EXISTING trusted registration
   path.

7. Coordinator exposes them only through its existing ArtifactReader.

8. Agents receive NO artifact-registration capability.

9. The artifact directory is local POC state, not a second database.

10. Restarting Coordinator must reproduce the same ArtifactReader state
    from the same `.tmp/artifacts` contents.

11. Missing/corrupt artifact must fail clearly or remain unavailable;
    never silently become trusted evidence.

12. `task verify` must be able to produce an artifact that the live
    Phase 8 environment can actually consume.

Keep the verifier bounded. Do not turn it into a CAS or research engine.

Tests MUST prove:

- verifier produces artifact
- hash is correct
- artifact is written to expected location
- Coordinator loads it
- ArtifactReader can retrieve it
- artifact survives Coordinator restart
- modified/corrupt artifact is rejected
- missing artifact is not presented as trusted
- no agent MCP tool can register artifacts
- the existing trusted registration path remains the only registration path

Also update the developer workflow so:

    task verify
    task dev

is sufficient to make the verified artifact available to ARGUS.

==================================================
P1-3 — SCENARIO-CONTAINMENT FOR CANONICAL REFERENCES
==================================================

Current packet grammar permits references such as:

    canonical:belief:<uuid>

That is correct for referencing existing Solvent state, but the validator
must enforce scenario isolation.

Required invariant:

    packet.scenario_ref
        ==
    referenced belief.scenario_id

If they differ:

    refuse packet / return validation error
    perform ZERO persistence

Apply the same ownership principle to canonical evidence references
where the evidence object has scenario ownership.

Also ensure that:
- local:<id> references remain packet-local
- canonical references never bypass scenario isolation
- edges cannot connect beliefs across unrelated scenarios unless an
  explicitly supported relation exists (none is required for Phase 8)

Do not add new schema.

Do not create a new cross-scenario policy subsystem.

Add tests:

1. canonical belief in same scenario → accepted
2. canonical belief in different scenario → refused
3. canonical evidence in same scenario → accepted
4. canonical evidence in different scenario → refused
5. malformed canonical reference → refused
6. dangling canonical reference → refused
7. failed validation causes no partial persistence

==================================================
REGRESSION REQUIREMENTS
==================================================

After implementing all three fixes, run:

    go test ./...
    go test -race ./...
    go vet ./...

Also run the relevant Solvent/Conductor tests if the Taskfile supports
them.

Run the developer environment checks:

    task dev
    task status
    task verify

Then verify:

    .tmp/artifacts/
        contains the trusted verifier artifact

    Coordinator
        can read the artifact

    Trust UI
        can operate against the same Coordinator state

Do NOT run the full OpenCode dry run yet unless the repository is already
ready for it. This task is to close the three implementation gaps first.

==================================================
DOCUMENTATION
==================================================

Update the implementation plan/report only where necessary to reflect:

P0-1:
    `task dev` works from fresh checkout without destructive reset.

P0-2:
    verifier artifacts have a deterministic local POC handoff into
    Coordinator's trusted ArtifactReader.

P1-3:
    canonical references are scenario-contained.

Do NOT add the broader hardening findings to this implementation.

==================================================
DELIVERABLE
==================================================

Create:

PHASE8_DEV_ENV_FIX_REPORT.md

Include:

1. P0-1 implementation + tests
2. P0-2 implementation + tests
3. P1-3 implementation + tests
4. exact files changed
5. test results
6. Taskfile verification results
7. any remaining limitation specific to these fixes

Final acceptance:

    task dev
      → environment ready

    task verify
      → trusted verifier artifact exists

    Coordinator
      → consumes that artifact

    canonical references
      → cannot cross scenario boundaries

No architectural expansion beyond these three fixes.