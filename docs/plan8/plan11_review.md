Plan 11 is **close, but I would make three corrections before implementation**. They are implementation-level, not architecture changes.

### 1. Fix the bootstrap sequence — this is the important one

There is a logic hole at lines 129–138. If CockroachDB is already running but the `argus` database does not exist, `db.Ping()` fails, `isCockroachRunning()` returns true, and the plan says to **return an error**. It never reaches `ensureDatabase()` or migrations. 

That defeats the goal:

```text
CRDB running
+
argus DB absent
        ↓
argus serve
        ↓
should create DB + migrate
```

Use this sequence instead:

```text
resolve DB URL
    ↓
can reach CRDB?
    ├─ no + local → start CRDB
    └─ yes         → reuse
    ↓
wait until CRDB ready
    ↓
ensure target database exists
    ↓
connect to target database
    ↓
apply migrations
    ↓
start ARGUS
```

In other words, **database creation and migrations must happen whether ARGUS started CRDB or found an existing CRDB**.

### 2. `task down` cannot currently stop ARGUS

The plan says `task down` stops “ARGUS and managed CRDB,” but the proposed `cmdServe` changes only write:

```text
.argus-pids/crdb.pid
```

There is no `argus.pid` being written by `serve`. Yet `task down` expects one at lines 223–234. 

Because `task dev` runs `go run ./cmd/argus serve` in the foreground, this is not catastrophic, but the specification is inconsistent.

For the ruthless/simple POC, I would actually simplify further:

```text
task dev  → argus serve
```

`Ctrl+C` stops ARGUS, and ARGUS stops a CRDB it personally started. `task down` only needs to kill a managed CRDB left behind by an abnormal termination. No ARGUS PID bookkeeping is necessary.

### 3. Do not use `pkill cockroach` in the smoke test

This:

```bash
pkill cockroach
```

is too destructive for a developer environment. It can kill an unrelated CockroachDB instance. 

Use the managed PID or `task down` instead.

---

## After those three corrections: **GO**

The core of Plan 11 is exactly right for the UX goal:

```text
./argus serve
    ↓
auto-start local CRDB when necessary
    ↓
create argus DB
    ↓
migrate
    ↓
start UI
```

It also correctly keeps external databases explicit and only auto-bootstraps the standard local database, while preserving the frozen architecture. 

I would **not** add more abstractions beyond this. The desired foundation is now:

> **ARGUS owns the local POC bootstrap; the developer owns nothing except running `argus serve`.**

That is exactly the right level of simplicity.
