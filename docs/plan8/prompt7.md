I reviewed the **Remaining POC Work** plan. It is close, but I would **not execute it unchanged**. There are four concrete corrections.

### 1. Blocker: Workstream 4 violates the frozen migration-ownership decision

The plan says to create:

```text
internal/epistemic/migrations/
  001_schema.sql
  ...
  010_debt_opaque.sql
```

by copying all Solvent migrations into ARGUS. 

That contradicts the approved architecture:

> Solvent owns canonical epistemic migrations; ARGUS consumes them through the Solvent migration export.

The plan should instead call the already-approved Solvent migration export:

```go
solventmigrations.Apply(db)
```

or its exported `embed.FS`.

**Do not copy the ten SQL files into ARGUS.**

### 2. Blocker: verifier hash specification is unsafe

The plan says `ComputeArtifactHash` excludes:

```text
ArtifactHash
Timestamp
VerifierID
Tolerance
```



That means someone could alter `VerifierID` or `Tolerance` without changing the artifact hash. Those are precisely part of the trust binding.

The canonical hash should exclude only self-referential/non-semantic fields, e.g.:

```text
exclude:
    ArtifactHash
    Timestamp, only if deliberately non-semantic

include:
    VerifierID
    VerifierVersion
    VerificationInput
    Tolerance
    Steps
    Result
```

The verifier binding must be cryptographically bound to the artifact content.

### 3. Incorrect wording: `VerifierSpec` enforcement is not compile-time

The risk section says:

> “VerifierSpec enforcement is compile-time in Validate().” 

That is incorrect.

`Validate()` is runtime validation. The actual enforcement should be:

```text
packet compilation
    ↓
load selected PackDefinition
    ↓
check verifier_id + version against VerifierSpec
    ↓
reject if not allowed
```

Rename the assertion to **runtime packet-validation enforcement**.

### 4. Integration environment needs precise wording

The plan says:

> “Uses real in-memory CRDB.” 

Use **ephemeral CockroachDB test instance** unless the implementation specifically uses CockroachDB's memory store mode. The test must state exactly how the database is started and isolated.

---

### One additional concern

The idempotency test says:

> “Submit packet, fail mid-transaction, retry same packet.” 

But the architecture intentionally uses **separate Solvent/work transactions**, so the test needs a deterministic fault injection point that causes a partial write and then proves entity-level deterministic IDs make the retry converge.

Otherwise this is only a happy-path duplicate test.

## Prompt to coding agent

```text
Revise the Remaining POC Work plan before execution. Do not implement yet.

Fix only these issues:

1. Remove the copied Solvent migration files from ARGUS. Use the approved
   Solvent-owned migration export (`solventmigrations.Apply(db)` / exported FS).
   ARGUS owns only work + idempotency/proposal migrations.

2. Fix VerificationArtifact hashing: VerifierID, VerifierVersion, input hash,
   tolerance, steps, and result must be covered by the canonical artifact hash.
   Only self-referential/non-semantic fields may be excluded.

3. Change "VerifierSpec enforcement is compile-time" to "runtime packet-validation
   enforcement" and specify the exact rejection path for an unauthorized verifier.

4. Replace "in-memory CRDB" with an explicit ephemeral CockroachDB test strategy.

5. Make the partial-failure retry test deterministic: inject failure after one
   entity persists, retry the identical packet, and assert deterministic IDs +
   ON CONFLICT produce exactly one logical entity of each type.

Do not change the architecture, add services, or expand the Solvent change budget.
Then freeze the plan and proceed to implementation.
```
