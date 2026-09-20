Yes. I read the attached execution plan, the forensic baseline, and the accompanying reviews. The most important issue is now clear:

> **The current plan still violates one of our strongest locked decisions: Solvent is frozen, and the Reference Loop must not modify the Solvent codebase.**

The plan is otherwise quite disciplined. The forensic-first approach, preservation of Agent A's valid work, real Conductor MCP, real Solvent REST/authentication, participant-owned evidence, and hard stop rules are all good.  

But the SQLite portion needs to be removed.

# The fundamental conflict

The plan says:

```text
Solvent repo is clean, but its current kernel/service/API path is pgx/Cockroach-shaped.
...
Implement strict Solvent SQLite support
...
modify Solvent's kernel/sql.go, kernel/kernel.go,
kernel/authority.go, service/*.go, api/*.go
```

and even defines a hard-stop threshold of more than eight modified kernel files.  

That is already incompatible with:

> **Solvent is frozen.**

The stop rule doesn't solve the problem. It merely says:

> "Modify the frozen kernel, but stop if we modify too much."

That is not the locked architecture.

The frozen boundary is:

```text
Solvent
    ↓
existing authority semantics
    ↓
existing CockroachDB/pgx implementation
```

The Reference Loop is supposed to test Solvent, not redesign its storage layer.

---

# Why this matters

The Reference Loop is supposed to answer:

> Does the existing Agent → Conductor → Solvent → Executor architecture work?

If we modify Solvent first, the experiment becomes:

> Does a modified Solvent adapted to our test environment work?

That weakens the empirical proof.

Worse, SQLite is not merely a storage-driver substitution here. The forensic report correctly identifies that the existing Solvent kernel relies on Cockroach-specific behavior:

```text
::UUID
JSONB
TEXT[]
gen_random_uuid()
crdb.ExecuteTx
FOR UPDATE
ON UPDATE CASCADE
DEFERRABLE FK
```



The proposed changes then alter transaction semantics, SQL semantics, locking semantics, array representation, and cascade behavior. 

Calling that "narrow storage compatibility" understates the risk.

Especially:

```text
crdb.ExecuteTx
        →
sql.Tx

SERIALIZABLE retries
        →
removed

FOR UPDATE
        →
removed

deferred constraints
        →
plain FK
```

Those are not just cosmetic database changes.

Even if domain-level behavior appears unchanged, **the authority kernel's concurrency and transactional semantics have changed**.

That is exactly what we should not be doing during the first empirical proof.

---

# What the plan should do instead

The architecture should become:

```text
Reference Loop
    │
    ├── Conductor
    │      └── SQLite / modernc.org/sqlite
    │
    ├── Solvent
    │      └── EXISTING pgx/CockroachDB path
    │
    └── Executor
           └── RecordingFunc / GitHub
```

Different database technologies are completely acceptable.

In fact, this is arguably a **better architectural test**.

We are proving that:

```text
Conductor
    ≠
Solvent
```

including at the persistence boundary.

There is no reason the Reference Loop needs one database technology everywhere.

---

# Use ephemeral CockroachDB for Solvent

The solution is not to weaken the frozen Solvent requirement.

The solution is to make the **test environment accommodate Solvent**.

For the Reference Loop:

```text
Conductor
    → local file-backed SQLite

Solvent
    → ephemeral/local CockroachDB

Executor
    → RecordingFunc or GitHub
```

The harness can provision a temporary CockroachDB instance for the reference run.

The exact mechanism can be:

```text
local CockroachDB
test container
temporary development instance
existing test fixture
```

depending on what already exists in the Solvent repository and development environment.

The important invariant is:

> **Zero Solvent source modification.**

That is much cleaner than porting Solvent to SQLite.

---

# The plan should therefore change in four places

## 1. Delete Task 2 entirely

This section:

> `Fix Solvent Kernel for SQLite`

should disappear.

There should be no:

```text
kernel/sql.go modifications
kernel/kernel.go modifications
authority.go modifications
service modifications
api SQL modifications
```

No compatibility fork.

No temporary branch.

No "just for the reference loop."

No exceptions.

---

## 2. Replace it with "Provision existing Solvent environment"

Something like:

```text
### Task 2: Provision Frozen Solvent

- Do not modify solvent-main.
- Use the existing Solvent pgx/CockroachDB implementation unchanged.
- Provision an isolated temporary CockroachDB instance for the Reference Loop.
- Apply the existing Solvent schema/migrations exactly as the repository expects.
- Start Solvent using its existing REST/API path.
- Use Solvent's real api.AuthMiddleware.
- Verify the reference loop talks to Solvent exactly as an external client would.
- Record the exact Solvent commit/version used by the experiment.
- If Solvent cannot be run unchanged in the reference environment, STOP and report BLOCKED.
```

This is consistent with the frozen-kernel rule.

---

# 3. The forensic matrix needs a different classification for the SQLite finding

The current forensic matrix says:

```text
Solvent DB driver     FAIL
Solvent kernel compat BLOCKED
```



That diagnosis is actually useful, but the disposition should now be:

```text
Solvent SQLite support
    = OUT OF SCOPE / ARCHITECTURALLY PROHIBITED
```

not:

```text
BLOCKED pending implementation
```

Why?

Because the incompatibility is not something we intend to solve.

It's simply:

> Solvent uses CockroachDB. Preserve that.

That distinction is important for Agent B.

---

# 4. The Reference Loop should not force one DB technology

This assumption:

```text
strict SQLite for the Reference Loop
```

should be removed.

The actual architectural requirement is:

> Conductor and Solvent have separate persistence boundaries.

That is already achieved with:

```text
Conductor → SQLite
Solvent → CockroachDB
```

The evidence collector can inspect both afterward.

The review itself already correctly emphasizes that participant-owned evidence should remain in participant-owned stores rather than becoming a centralized ledger. 

---

# There is another subtle benefit

This makes the Reference Loop a **stronger** experiment.

We are not testing:

> Can we make two components use the same convenient local database?

We are testing:

> Can two independently owned systems with different persistence technologies cooperate through their actual interfaces?

That is much closer to the architecture we claim to support.

---

# The interface corrections remain correct

The plan's other Agent A findings are good and should remain:

### Conductor

Replace direct SQLite access with the real MCP surface:

```text
Agent/harness
    ↓
Conductor MCP
    ↓
Conductor
```

not:

```text
harness
    ↓
Conductor DB
```

That is exactly the right correction. 

### Solvent

Use its real REST path and real `api.AuthMiddleware`.

That is also correct. 

### Evidence

The plan is also right that the harness must not manufacture:

```text
EFFECT_CONFIRMED
```

In GitHub mode, actual GitHub workflow state must establish the result. RecordingFunc must be explicitly labeled as simulated/local execution evidence. 

### Operation identity

The GitHub deploy operation is a good reference operation:

```text
deploy
repo = pithomlabs/reference-loop-test
workflow = .github/workflows/ref-loop.yml
ref = main
run_id = included when effect-relevant
```

and the same identity must flow through proposal → authorization → execution → evidence. 

---

# One more thing I would change: do not create a test-only Solvent fork

There is a temptation here:

> "We can just keep the frozen Solvent repo untouched and create a small SQLite-compatible fork under `reference-loop/`."

Don't do that either.

That creates:

```text
Solvent production semantics
        ≠
Reference Loop Solvent semantics
```

and now we no longer know whether we are testing Solvent.

The Reference Loop should exercise the **actual Solvent implementation**.

That is the point of the empirical pivot.

---

# Revised architecture under test

The clean environment should be:

```text
                         ┌──────────────────┐
                         │      Agent       │
                         └───────┬──────────┘
                                 │
                              MCP│
                                 ▼
                         ┌──────────────────┐
                         │    Conductor     │
                         │     SQLite       │
                         └───────┬──────────┘
                                 │
                                 │ proposal / coordination
                                 ▼
                         ┌──────────────────┐
                         │     Solvent      │
                         │ EXISTING CODE    │
                         │ pgx/CockroachDB  │
                         └───────┬──────────┘
                                 │
                              REST
                                 │
                                 ▼
                         ┌──────────────────┐
                         │     Executor     │
                         └───────┬──────────┘
                                 │
                                 ▼
                         ┌──────────────────┐
                         │ GitHub Actions   │
                         │   authoritative  │
                         │       SOR        │
                         └──────────────────┘
```

This is a much better Reference Loop than forcing everything onto SQLite.

---

# Revised stop rules

I would make the Agent B rules extremely explicit:

```text
CONDUCTOR:
    no source modification

SOLVENT:
    no source modification

EXECUTOR:
    no semantic modification unless existing implementation is demonstrably defective

Reference Loop:
    adapters/drivers/tests/configuration are allowed

Database:
    Conductor may use SQLite
    Solvent must use its existing CockroachDB path

If Solvent cannot run unchanged:
    STOP
    report BLOCKED
    do not port Solvent to SQLite
```

This is the important change.

---

# One other issue in the current plan

The plan says:

> `Run existing Solvent unit tests`

That is fine.

But the test should be:

> establish that the frozen Solvent repository still passes its own existing tests.

Not:

> change Solvent until tests pass on SQLite.

That's a completely different activity.

---

# Final disposition

The current plan is **very close**, but I would not hand it to Agent B in its current form because the SQLite requirement contradicts the frozen-Solvent decision.

### Keep

```text
forensic assessment first
contract pin
real Conductor MCP
real Solvent REST/auth
GitHub deploy
participant-owned evidence
RecordingFunc distinction
operation identity
happy path only
hard stop / report
BM-IST deferred
```

### Remove

```text
Solvent SQLite port
Cockroach → SQLite kernel changes
SQL-dialect translation
crdb.ExecuteTx → sql.Tx
removal of FOR UPDATE
application emulation of DB semantics
```

### Replace with

```text
Provision isolated existing CockroachDB
        ↓
Run frozen Solvent unchanged
        ↓
Exercise its real REST/API surface
```

This is a much cleaner interpretation of the pivot:

> **We are modifying the Reference Loop to accommodate the frozen architecture—not modifying the frozen architecture to accommodate the Reference Loop.**

That principle should be made the first non-negotiable rule at the top of Agent B's plan.

