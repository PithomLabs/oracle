This revision is substantially better and incorporates the consolidated amendments correctly in most places. I would still **hold implementation for one more small plan revision**. There are four concrete blockers and a few cleanup items.

### 1. Migration ownership is still technically wrong

The plan says Solvent owns canonical migrations, which is correct, but then says ARGUS consumes them “at build time (e.g., via `go:embed` of Solvent's migration SQL files).” 

ARGUS cannot `go:embed` files located in another Go module's repository.

Use this instead:

```text id="6x6bq4"
Solvent owns migrations
        ↓
Solvent exposes versioned embedded migration FS/package
        ↓
ARGUS imports Solvent migration package
        ↓
ARGUS applies Solvent migrations
        ↓
ARGUS applies its own work/idempotency migrations
```

For example, Solvent can expose:

```go
package migrations

var FS embed.FS
```

or an `Apply(db)` function.

That preserves the decision you made—**Solvent remains canonical**—without duplicating SQL.

### 2. The idempotency model still has a failure hole

The plan inserts the idempotency row **before** persistence and then says:

> “On any failure after step 2, the idempotency row remains — caller must retry with same packet_id.” 

But a retry will see the unique row and may conclude “already submitted,” even though the first attempt only partially persisted.

That is not actually retry-safe.

The POC needs a tiny state machine:

```text id="w4h7fa"
submission
    NEW
    IN_PROGRESS
    COMPLETED
    FAILED
```

Or, even simpler, the idempotency record needs:

```text
content_hash
scenario_id
packet_id
status
result/ref
error
created_at
```

Then:

```text id="bb0x0e"
first request → IN_PROGRESS
success       → COMPLETED
failure       → FAILED
retry FAILED  → may retry
retry COMPLETED → return prior result
```

This preserves the reviewers' valid insight—**DB-enforced identity rather than an in-memory cache**—without creating a poisoned idempotency record.

### 3. The verifier binding is still missing `tolerance`

The amendment explicitly called for:

```text
verifier_id
verifier_version
input_hash
tolerance
```

But the proposed `VerificationArtifact` contains no tolerance field. 

Add something like:

```go
Tolerance string `json:"tolerance,omitempty"`
```

or a typed tolerance structure if the actual verifier needs more than one numeric parameter.

Also define exactly how `ArtifactHash` is computed. A self-referential hash must exclude the `artifact_hash` field from the canonical hash input; otherwise the schema is ambiguous.

### 4. The human authentication flow still stops one step short

The plan says:

```text
operator credential
    ↓
server-side principal
    ↓
HttpOnly session cookie / CSRF token
```

but doesn't define how the browser obtains that session. 

For a POC, keep it extremely simple:

```text id="w8p4k2"
GET /login
    ↓
operator token
    ↓
server validates configured token
    ↓
HttpOnly session cookie
    ↓
consequential POST
    ↓
CSRF check + server-derived principal
```

No users table, OAuth, identity provider, or RBAC system.

The key property is that:

```text
discharged_by
```

never comes from the request body.

---

## One important developer-experience correction

`task dev` still resets the entire database on every run. 

That contradicts the earlier objective of making the developer environment safe to repeat.

I would make:

```text
task dev       # non-destructive
task fresh     # destructive reset
task test
task verify
```

`argus reset` remains explicit and destructive.

This is not an architectural blocker, but it is worth fixing before implementation because the whole pivot is partly about reducing operational friction.

---

## What is now solid

The revised plan gets the major architectural decisions right:

**Solvent stays conceptually authoritative while losing its HTTP deployment boundary.** 

**Conductor is reduced to the actual work semantics rather than preserved as a service.** 

**Agent submission is explicitly prevented from performing discharge/promotion/retraction/authorization.** 

**Agent-created edges remain non-authoritative proposals.** 

**The transaction limitation is now stated honestly instead of pretending the pivot creates global atomicity.** 

**The refusal suite is now part of acceptance, which is crucial for this architecture.** 

## Verdict

I would classify this as:

```text
Architecture:        APPROVED
Pivot:               APPROVED
Plan structure:      APPROVED
Implementation:      HOLD briefly
```

The remaining required corrections are small:

```text
1. Solvent migration bundle mechanism
2. Idempotency failure/retry state
3. Verifier tolerance + canonical artifact hashing
4. Minimal browser session bootstrap
5. Non-destructive task dev
```

After those, I would consider the plan **implementation-ready** rather than continuing another large review cycle.
