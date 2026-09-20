I reviewed v4 specifically against the two external adversarial reviews and the POC constraint. 

**Most of those findings can wait for hardening. I see only three things I would add now because they can prevent the Phase 8 POC from honestly demonstrating what it claims.**

### 1. `task dev` has a real first-run defect — fix now

The plan says `task dev` is safe and can run from a fresh checkout, but `task up` calls `db:ensure`; `db:ensure` only detects a missing database and then `exit 1`, with a comment saying this "triggers fallback to setup:db." There is no actual fallback mechanism in the plan.  

So on a genuinely fresh machine:

```text
task dev
→ DB absent
→ db:ensure
→ exit 1
→ task dev fails
```

That violates the core developer-experience objective.

Make `db:ensure`:

```text
DB absent
  → invoke canonical Solvent db:reset
  → success

DB exists + complete schema
  → continue

DB exists + incomplete/invalid schema
  → fail with:
     "Run task fresh to reinitialize"
```

Do not silently reset an existing database.

This is **P0 for the dev-environment plan**.

---

### 2. Physics Verifier artifact handoff is missing — fix now

This is the most important net-new architectural integration finding.

The Coordinator launcher creates a **new empty in-memory ArtifactRegistry**:

```text
coordinator
→ verifier.NewArtifactRegistry()
```



Meanwhile the standalone verifier creates its **own separate registry**, runs the verification, and outputs the artifact. 

There is no bridge.

Therefore:

```text
task verify
→ verifier produces artifact
```

does **not** imply:

```text
ARGUS Coordinator
→ can resolve/read that artifact
```

For a system whose BM-IST retirement rules can require `reproducible_artifact`, that is a real POC integration gap.

Add the smallest possible handoff:

```text
task verify
  ↓
.tmp/artifacts/g0.json
.tmp/artifacts/g0.sha256
```

Then define a **trusted local bootstrap path** for Coordinator:

```text
Coordinator startup
→ load verified local artifacts from .tmp/artifacts
→ register through existing trusted registration mechanism
→ expose only via ArtifactReader
```

Do not create an agent-facing registration API.

The important invariant is:

> `task verify` must produce an artifact that the running ARGUS instance can actually consume as evidence.

This is **P0 for the "test ARGUS + Physics Verifier" objective**.

---

### 3. Canonical references need scenario containment — fix now

The external reviewers caught one genuinely important epistemic-integrity issue: `canonical:belief:<uuid>` needs to be validated against the packet's scenario.

Otherwise:

```text
Packet scenario A
    ↓
canonical:belief:<uuid from scenario B>
    ↓
evidence/edge attached across scenario boundary
```

That violates the isolation implied by scenario-based RCP/idempotency.

The fix is tiny:

```text
packet.scenario_id == referenced_belief.scenario_id
```

and equivalent ownership checks for canonical evidence where applicable.

Add this as a Phase 8 **pre-dry-run invariant/test**, not a new schema.

This is **P1**, but cheap enough that I would do it now rather than knowingly enter the dry run with cross-scenario contamination possible.

---

## What I would deliberately defer

These are valid adversarial observations, but they are **hardening**, not reasons to delay this POC:

```text
CSRF / DNS rebinding hardening
full agent container/network sandbox
PROPOSED/ACTIVE edge lifecycle
cryptographic verifier attestation
environment/host provenance for artifacts
quotas / anti-DoS
saga/distributed transaction machinery
partial debt retirement
multi-round adversarial testing
full multi-user authentication
```

The current capability boundary is sufficient for the POC so long as the OpenCode runtime is not given Solvent/Conductor credentials or privileged MCP servers.

---

## One clarification worth recording, not implementing

Z's concern about agent-authored `contradicts` edges is worth documenting:

```text
Agent creates contradicts edge
        ≠
Agent retracts belief
```

The consequential transition still requires the human `RETRACT` path.

So **do not add `PROPOSED/ACTIVE` edge state now**.

---

# Final disposition

I would amend the attached plan with exactly these three:

```text
P0  Fix fresh `task dev` DB initialization / schema detection.
P0  Add Physics Verifier → ARGUS trusted artifact handoff.
P1  Enforce canonical reference scenario containment.
```

Everything else from DeepSeek/Z can go into a later **Phase 8 hardening backlog**.

After those three are added, **I would approve the plan for implementation**. The POC remains small, but its core demonstration will actually be executable and internally honest.
