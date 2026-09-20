# Plan 11: Self-Bootstrapping `argus serve`

**Status:** DRAFT — awaiting implementation instruction
**Date:** 2026-09-19
**Objective:** Make `argus serve` zero-fuss local POC startup

---

## Current Failure

Running `./argus serve` fails with:

```
ping db: failed to connect to `user=root database=argus`:
127.0.0.1:26257: connect: connection refused
```

Current behavior assumes CockroachDB is already running.
The desired POC experience is: `./argus serve` and nothing else.

## Current State

- `cockroach` binary at `/home/chaschel/.local/bin/cockroach` (v26.2.0)
- `.cockroach-data/` exists with prior data
- No CRDB process running; port 26257 free
- `cmd/argus/main.go` does `sql.Open` + `db.Ping` — fails hard if CRDB absent
- `Taskfile.yml` has CRDB orchestration in shell (start, wait, migrate)
- All migrations use `CREATE TABLE IF NOT EXISTS` — fully idempotent
- No `os/exec` in `cmd/argus/` — only in `reference-loop/conductor/client.go`
- Tests assume running CRDB on `:26257`

## Architecture Constraints

FROZEN — do not change:

- One ARGUS Go process
- One CockroachDB
- Embedded Trust UI
- Two MCP tools
- Fixed local ports (CRDB SQL 26257, CRDB Admin 8081, ARGUS HTTP 8080)
- `CAPABILITY != WORK != AUTHORITY != EXECUTION`
- Server-derived PrincipalID
- Server-selected PackRef
- All authority/security invariants

---

## Files Changed

| File | Action | Purpose |
|------|--------|---------|
| `cmd/argus/bootstrap.go` | **NEW** | ~180 lines: `bootstrap()`, process detection, readiness wait, DB creation, migration |
| `cmd/argus/main.go` | EDIT | Replace `sql.Open`+`Ping` with `bootstrap()` call; add CRDB cleanup on shutdown |
| `Taskfile.yml` | EDIT | Simplify `task dev` -> thin wrapper; keep `task fresh`, `task down` |
| `README.md` | EDIT | Update quick-start to single-command startup |
| `cmd/argus/bootstrap_test.go` | **NEW** | Unit tests for bootstrap logic |

---

## Part 1: `cmd/argus/bootstrap.go`

### `resolveDBURL(flagValue string) string`

Priority chain:
1. `--db` flag (if explicitly set)
2. `ARGUS_DB_URL` env var
3. Default: `postgres://root@localhost:26257/argus?sslmode=disable`

Detect explicit flag vs default by using a second flag variable (`dbURLSet bool`).

### `isLocalURL(dbURL string) bool`

Returns true only if URL targets `localhost:26257` or `127.0.0.1:26257`.
Guards auto-start: only bootstrap CRDB if user hasn't overridden to a remote DB.

### `findCockroachBinary() (string, error)`

Search order:
1. `$PATH`
2. Common locations: `/usr/local/bin/cockroach`, `~/.local/bin/cockroach`

If not found: single clear error message.

### `isCockroachRunning(host string) bool`

TCP dial `host:26257` with 500ms timeout.
- Connectable -> CRDB already running -> reuse
- Not connectable -> need to start one

### `startLocalCockroach() (*exec.Cmd, error)`

```go
exec.Command("cockroach", "start-single-node", "--insecure",
    "--listen-addr", ":26257",
    "--http-addr", ":8081",
    "--store=.cockroach-data")
```

- Attach `Stdout`/`Stderr` to `os.Stdout`/`os.Stderr` (CRDB logs visible)
- Write PID to `.argus-pids/crdb.pid`
- Return `*exec.Cmd` for caller cleanup

### `waitForSQL(host string, timeout time.Duration) error`

Poll loop: TCP dial `host:26257` every 500ms, up to `timeout` (default 30s).
Secondary check: `sql.Open` + `db.Ping`.
Return error if timeout exceeded.

### `ensureDatabase(dbURL string) error`

1. Parse database name from URL (default: `argus`)
2. Connect to `postgres://root@localhost:26257/defaultdb?sslmode=disable`
3. `CREATE DATABASE IF NOT EXISTS argus`
4. Close admin connection

### `applyMigrations(dbURL string) error`

1. `sql.Open("pgx", dbURL)`
2. `solventmigrations.Apply(ctx, db)`
3. `migrations.Apply(ctx, db)`
4. Close

All idempotent — safe every startup.

### `bootstrap(dbURL string) (string, *exec.Cmd, error)`

Orchestrates full sequence:

1. `resolveDBURL(dbURL)`
2. `sql.Open` + `db.Ping` — if succeeds -> return `(dbURL, nil)` (existing CRDB, no process to manage)
3. If ping fails AND `isLocalURL(dbURL)`:
   - `findCockroachBinary()`
   - `isCockroachRunning()` — if true -> error (port open but SQL not responding)
   - `startLocalCockroach()`
   - `waitForSQL()`
   - `ensureDatabase(dbURL)`
   - `applyMigrations(dbURL)`
   - Return `(dbURL, cmd)` — caller must clean up process
4. If ping fails AND NOT `isLocalURL(dbURL)` -> single actionable error

---

## Part 2: `cmd/argus/main.go` Changes

### `cmdServe` rewrite

```go
func cmdServe(args []string) {
    fs := flag.NewFlagSet("serve", flag.ExitOnError)
    dbURL := fs.String("db", "", "database URL (default: local CRDB)")
    listen := fs.String("listen", ":8080", "HTTP listen address")
    fs.Parse(args)

    resolvedURL, managedCRDB, err := bootstrap(*dbURL)
    if err != nil {
        log.Fatalf("ARGUS: %v", err)
    }
    if managedCRDB != nil {
        defer func() {
            log.Println("stopping managed CockroachDB...")
            managedCRDB.Process.Signal(syscall.SIGTERM)
            managedCRDB.Wait()
            os.Remove(".argus-pids/crdb.pid")
        }()
    }

    db, err := sql.Open("pgx", resolvedURL)
    if err != nil {
        log.Fatalf("ARGUS: open db: %v", err)
    }
    defer db.Close()

    // ... rest unchanged (app setup, pack load, UI, HTTP server) ...
}
```

### Same pattern for `cmdMigrate`, `cmdReset`, `cmdMCP`

All entry points call `bootstrap()`.

### Shutdown cleanup

```go
go func() {
    <-ctx.Done()
    log.Println("shutting down...")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(shutdownCtx)
    if managedCRDB != nil {
        managedCRDB.Process.Signal(syscall.SIGTERM)
        managedCRDB.Wait()
        os.Remove(".argus-pids/crdb.pid")
    }
}()
```

---

## Part 3: Taskfile Simplification

### `task dev`

```yaml
dev:
  desc: Start ARGUS (non-destructive)
  cmds:
    - go run ./cmd/argus serve
```

### `task fresh`

```yaml
fresh:
  desc: Reset database and start fresh (destructive)
  cmds:
    - go run ./cmd/argus reset
    - go run ./cmd/argus serve
```

### `task down`

```yaml
down:
  desc: Stop ARGUS and managed CRDB
  cmds:
    - |
      if [ -f .argus-pids/argus.pid ]; then
        kill $(cat .argus-pids/argus.pid) 2>/dev/null || true
        rm -f .argus-pids/argus.pid
      fi
      if [ -f .argus-pids/crdb.pid ]; then
        kill $(cat .argus-pids/crdb.pid) 2>/dev/null || true
        rm -f .argus-pids/crdb.pid
      fi
```

---

## Part 4: Double-CRDB Prevention

Critical invariant: **never start a second CRDB against `.cockroach-data`**.

1. `isCockroachRunning("localhost:26257")` checks TCP before `start-single-node`
2. TCP connectable -> CRDB already running -> reuse it
3. `startLocalCockroach()` writes PID to `.argus-pids/crdb.pid` atomically
4. `cmdServe` registers `SIGTERM` handler -> sends `SIGTERM` to managed CRDB on exit
5. If `.argus-pids/crdb.pid` exists and PID alive -> `task dev` refuses to start

---

## Part 5: Error Messages

Each failure produces exactly one message:

| Failure | Message |
|---------|---------|
| `cockroach` not found | `ARGUS: cockroach executable not found in PATH. Install CockroachDB or provide --db/ARGUS_DB_URL.` |
| CRDB start fails | `ARGUS: failed to start CockroachDB: <error>. Check .cockroach-data permissions.` |
| Readiness timeout | `ARGUS: CockroachDB did not become ready within 30s. Check .cockroach-data for corruption.` |
| DB creation fails | `ARGUS: failed to create 'argus' database: <error>.` |
| Migration fails | `ARGUS: migration failed: <error>.` |
| Custom DB unreachable | `ARGUS: cannot connect to database at <url>. Start CockroachDB or check the URL.` |
| HTTP port conflict | `ARGUS: listen tcp :8080: bind: address already in use.` |

---

## Part 6: Tests (`cmd/argus/bootstrap_test.go`)

### Unit tests (no CRDB needed)

1. `TestResolveDBURL_Default` — no flag, no env -> default URL
2. `TestResolveDBURL_FlagOverride` — `--db` flag -> flag value
3. `TestResolveDBURL_EnvOverride` — `ARGUS_DB_URL` set -> env value
4. `TestIsLocalURL_Localhost` — `postgres://root@localhost:26257/argus` -> true
5. `TestIsLocalURL_127` — `postgres://root@127.0.0.1:26257/argus` -> true
6. `TestIsLocalURL_Remote` — `postgres://root@otherhost:26257/argus` -> false
7. `TestFindCockroachBinary_NotFound` — empty PATH -> error

### Integration tests (need CRDB, skip with `-short`)

8. `TestEnsureDatabase_CreateIfNotExists`
9. `TestMigrationsIdempotent` — apply twice -> no error
10. `TestBootstrap_ExistingDB` — CRDB running -> reuse, no managed process
11. `TestBootstrap_LocalAbsent` — no CRDB -> start, create DB, migrate

---

## Part 7: `README.md` Update

### Quick Start (replace current)

```
## Quick Start

    git clone <repo-url> && cd oracle
    go build ./cmd/argus
    ./argus serve

ARGUS automatically starts CockroachDB, creates the `argus` database,
applies migrations, and starts the Trust UI on :8080.

Login token: `argus-local-operator` (default).
```

### Overrides (brief)

- `--db <url>` or `ARGUS_DB_URL` — connect to external CRDB
- `ARGUS_OPERATOR_TOKEN` — override login token
- `--listen <addr>` — override HTTP listen address

---

## Verification

```bash
# Unit tests (no CRDB needed)
go test -short ./cmd/argus/...
go test -short ./...

# Full tests (bootstrap handles CRDB)
go test ./...
go test -race ./...
go vet ./...

# Smoke test from clean state
pkill cockroach 2>/dev/null
rm -rf .cockroach-data .argus-pids .task
go build -o argus ./cmd/argus
./argus serve
# -> starts CRDB, creates DB, migrate, listen on :8080
# Ctrl+C
./argus serve
# -> detects existing CRDB, reuses, no double-start
```

---

## What This Does NOT Change

- No new services, databases, or ports
- No dynamic port allocation
- No process supervisor framework
- No authority model changes
- No migration content changes
- No Trust UI, MCP, or Coordinator code changes
- Tests continue to assume CRDB on `:26257`

---

## Stop Condition

Once `./argus serve` works from a clean local state with sane defaults, STOP.
The target is not a perfect process supervisor. The target is ZERO-FUSS LOCAL POC STARTUP.
