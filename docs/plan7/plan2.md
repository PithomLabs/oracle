# PHASE 8 DEV ENV IMPLEMENTATION PLAN (REVISED)

**Date:** 2026-09-18
**Status:** Draft v2 — Revised per architectural review
**Scope:** Developer experience reproducibility for ARGUS Phase 8 dry run
**Review applied:** docs/plan7/plan_review.md (all 16 correction points)

---

## Changelog from v1

| # | Issue | v1 Defect | v2 Fix |
|---|-------|-----------|--------|
| 1 | MCP architecture | `cmd/argus-mcp` constructed in-process Coordinator | MCP is thin HTTP client; sole downstream dependency is Coordinator |
| 2 | Startup ordering | `up:solvent` waited before starting | Every service: start → wait → next |
| 3 | Sibling repos | `../solvent-main`, `../conductor` assumed | Auto-clone pinned revisions into `.tmp/src/` |
| 4 | Solvent migrations | Inlined 10 migration filenames | Call Solvent's canonical `db:reset` via `dir:` override |
| 5 | Conductor project | Created before Conductor running | Moved to post-readiness in `task up` |
| 6 | Verifier execution | `go test -run TestPhysicsVerifier` (library only) | Execute `cmd/verifier` binary against checked-in fixture |
| 7 | Test scope | Oracle-only | Clarified: `task test` = Oracle+TrustUI, `task test:all` = full stack |
| 8 | Process idempotency | Unconditional `nohup` | Check-before-start, stale PID handling |
| 9 | MCP env vars | SOLVENT_URL, CONDUCTOR_URL, operator ID exposed | Only `ARGUS_COORDINATOR_URL` |
| 10 | Coordinator launcher | Already correct | Kept as-is |
| 11 | PACK_DIR | CWD-relative `./domain-pack` | Absolute repo-root path |
| 12 | Port map | Frozen | Unchanged |

---

## 1. Objective

```text
git clone <oracle-repo>
    ↓
task setup        # obtain deps, build, init DB
    ↓
task up           # start services in order
    ↓
task status       # all READY
    ↓
task test         # Oracle + Trust UI unit tests
task verify       # Physics Verifier binary + corpus
    ↓
Two fresh OpenCode processes + Trust UI → Phase 8 dry run
```

---

## 2. Corrected Runtime Topology

```text
task setup
│
├─ verify prerequisites (docker, go, task, git)
├─ clone pinned Solvent → .tmp/src/solvent
├─ clone pinned Conductor → .tmp/src/conductor
├─ build binaries → .tmp/bin/
├─ start CockroachDB (Docker)
├─ wait DB ready
├─ invoke Solvent db:reset (canonical, all 10 migrations)
├─ seed operator principal
└─ generate .opencode/config.json

task up
│
├─ start Solvent :9090
├─ wait Solvent (GET /v1/ledger → 200)
├─ start Conductor :9091
├─ wait Conductor (GET /v1/projects → 200)
├─ create default project (idempotent POST)
├─ start Coordinator :8080
├─ wait Coordinator (GET /packs/rules → 200)
├─ start Trust UI :8081
└─ wait Trust UI (GET /insights → 200)

OpenCode (Work Agent)
│
└─ stdio MCP
      ↓
  argus-mcp  (thin HTTP client)
      ↓ HTTP
  Coordinator :8080
      ├── Solvent :9090
      └── Conductor :9091

OpenCode (Adversarial Agent)
│
└─ same MCP path, fresh process, no shared state
```

---

## 3. New Files

| File | Repository | Purpose | Dependencies |
|------|-----------|---------|-------------|
| `cmd/coordinator/main.go` | oracle | Standalone Coordinator HTTP server | coordinator, coordinator/http, domain-pack, verifier |
| `cmd/argus-mcp/main.go` | oracle | MCP stdio entry point | mcp/server, mcp/adapter (HTTP client only) |
| `cmd/verifier/main.go` | oracle | Physics Verifier CLI | verifier, verifier/physics/v1 |
| `mcp/server.go` | oracle | MCP stdio JSON-RPC transport | github.com/modelcontextprotocol/go-sdk |
| `Taskfile.yml` | oracle | Developer orchestration | All above + .tmp/src deps |
| `.opencode/config.json` (generated) | oracle | OpenCode MCP config | — |
| `.gitignore` (update) | oracle | Exclude `.tmp/` | — |

**No changes to Solvent or Conductor source code.**

---

## 4. Coordinator Launcher

### File: `oracle/cmd/coordinator/main.go`

Thin launcher. Owns: Solvent/Conductor URLs, PackRegistry, ArtifactReader,
operator config, HTTP serving. Does NOT duplicate Coordinator logic.

```go
func main() {
    solventURL := envOr("SOLVENT_URL", "http://localhost:9090")
    conductorURL := envOr("CONDUCTOR_URL", "http://localhost:9091")
    operatorID := envRequired("ARGUS_OPERATOR_PRINCIPAL_ID")
    operatorToken := envRequired("ARGUS_OPERATOR_TOKEN")
    addr := envOr("COORDINATOR_ADDR", ":8080")
    // Absolute path to avoid CWD sensitivity
    packDir := envOr("PACK_DIR", "")

    packRegistry := domainpack.NewRegistry()
    if packDir != "" {
        packRegistry.LoadFromDisk(packDir)
    }

    artifactRegistry := verifier.NewArtifactRegistry()

    coord, err := coordinator.New(coordinator.Config{
        SolventBaseURL:   solventURL,
        ConductorBaseURL: conductorURL,
        PackRegistry:     packRegistry,
        ArtifactReader:   artifactRegistry.AsReader(),
        OperatorID:       operatorID,
    })
    coord.SetOperatorToken(operatorToken)

    srv := coordinatorhttp.NewServer(coord)
    httpSrv := &http.Server{Addr: addr, Handler: srv.Handler()}
    // graceful shutdown on SIGINT/SIGTERM
}
```

**Env vars:**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SOLVENT_URL` | No | `http://localhost:9090` | Solvent API base URL |
| `CONDUCTOR_URL` | No | `http://localhost:9091` | Conductor API base URL |
| `ARGUS_OPERATOR_PRINCIPAL_ID` | Yes | — | Operator principal UUID |
| `ARGUS_OPERATOR_TOKEN` | Yes | — | Bearer token for /decisions |
| `COORDINATOR_ADDR` | No | `:8080` | Listen address |
| `PACK_DIR` | No | (empty = skip loading) | Absolute path to domain-pack dir |

**Readiness probe:** `GET /packs/rules` (no auth, confirms packs loaded).

**Tests:** Existing `coordinator/http/handler_test.go`. Add `TestCoordinatorServerStartsAndStops`.

---

## 5. MCP Architecture (Corrected)

### Principle

```text
OpenCode → stdio MCP → cmd/argus-mcp → HTTP → Coordinator :8080
```

`cmd/argus-mcp` is a **thin stdio transport + HTTP client**. It:
- Reads JSON-RPC from stdin, writes to stdout
- Translates tool calls to Coordinator HTTP requests
- Logs to stderr only
- Has exactly one downstream dependency: Coordinator

It does NOT:
- Construct a Coordinator
- Load PackRegistry or ArtifactRegistry
- Know Solvent URL, Conductor URL, or operator identity
- Have any write access to Solvent or Conductor

### File: `oracle/mcp/server.go`

MCP stdio JSON-RPC server using `github.com/modelcontextprotocol/go-sdk`.

```go
func Run(ctx context.Context, handler http.Handler, coordinatorURL string) error {
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "argus",
        Version: "v0.1.0",
    }, nil)

    server.AddTool(&mcp.Tool{
        Name:        "argus.get_context",
        Description: "Read RCP/v1 context for a task from Coordinator",
        InputSchema: map[string]any{
            "type": "object",
            "properties": map[string]any{
                "task_id": map[string]any{"type": "string"},
            },
            "required": []string{"task_id"},
        },
    }, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // HTTP GET to Coordinator
    })

    server.AddTool(&mcp.Tool{
        Name:        "argus.submit_packet",
        Description: "Submit an EBP research packet to Coordinator",
        InputSchema: /* full packet schema */,
    }, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // HTTP POST to Coordinator
    })

    return server.Run(ctx, &mcp.StdioTransport{})
}
```

### File: `oracle/cmd/argus-mcp/main.go`

```go
func main() {
    coordinatorURL := envOr("ARGUS_COORDINATOR_URL", "http://localhost:8080")

    // Verify Coordinator reachable on startup
    resp, err := http.Get(coordinatorURL + "/packs/rules")
    if err != nil || resp.StatusCode != 200 {
        log.Fatalf("Coordinator unreachable at %s", coordinatorURL)
    }

    // Create HTTP client for Coordinator
    client := &http.Client{Timeout: 30 * time.Second}

    mcpserver.Run(context.Background(), client, coordinatorURL)
}
```

**Env vars (ONLY):**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ARGUS_COORDINATOR_URL` | No | `http://localhost:8080` | Coordinator base URL |

No SOLVENT_URL, no CONDUCTOR_URL, no operator identity, no pack loading.

### OpenCode configuration (generated by `task config:generate`)

```json
{
  "mcpServers": {
    "argus": {
      "command": ".tmp/bin/argus-mcp",
      "env": {
        "ARGUS_COORDINATOR_URL": "http://localhost:8080"
      }
    }
  }
}
```

**Tests:**

| Test | What it verifies |
|------|-----------------|
| Binary builds | `go build ./cmd/argus-mcp` succeeds |
| Exactly 2 tools | `tools/list` returns `argus.get_context` + `argus.submit_packet` |
| No authority tools | Solvent/Conductor mutation tools absent |
| Startup check | Fails fast if Coordinator unreachable |
| HTTP forwarding | Tool calls forwarded to Coordinator correctly |

---

## 6. Physics Verifier (Corrected)

### File: `oracle/cmd/verifier/main.go`

Standalone CLI. Two modes:
- `--input <file>` — read VerifierInput JSON from file
- Default (no flag) — use built-in G0 test input

```go
func main() {
    inputFile := flag.String("input", "", "Path to VerifierInput JSON (optional)")
    flag.Parse()

    var input physicsv1.VerifierInput
    if *inputFile != "" {
        data, _ := os.ReadFile(*inputFile)
        json.Unmarshal(data, &input)
    } else {
        input = physicsv1.DefaultG0Input()  // built-in fixture
    }

    registry := verifier.NewArtifactRegistry()
    artifact, err := verifier.RunPhysicsVerifier(context.Background(), registry, input)

    // Print artifact JSON to stdout
    // Print summary to stderr
    // Exit 0 on confirmed, 1 on refuted/inconclusive
}
```

### Fixture: `oracle/verifier/testdata/g0_input.json`

Checked-in deterministic input for `task verify`. Versioned, inspectable,
independently reusable.

### `task verify` (corrected)

```yaml
verify:
  desc: "Run Physics Verifier binary + corpus checks"
  cmds:
    - task: verify:corpus
    - task: verify:verifier

verify:verifier:
  desc: "Execute Physics Verifier against G0 fixture"
  cmds:
    - task: build:verifier
    - .tmp/bin/verifier --input verifier/testdata/g0_input.json
    # Exit code 0 = confirmed, non-zero = refuted/inconclusive
```

This actually exercises the standalone binary, not just the library tests.

---

## 7. Dependency Acquisition (Fresh Checkout)

### Variables

```yaml
vars:
  SOLVENT_REPO_URL: "https://github.com/PithomLabs/solvent.git"
  SOLVENT_REV: "v0.1.0"          # pinned tag or commit SHA
  CONDUCTOR_REPO_URL: "https://github.com/PithomLabs/conductor.git"
  CONDUCTOR_REV: "v0.1.0"        # pinned tag or commit SHA
  SRC: .tmp/src
```

### Task: `setup:deps`

```yaml
setup:deps:
  desc: "Clone pinned Solvent and Conductor sources"
  cmds:
    # Solvent
    - |
      if [ -d "{{.SRC}}/solvent" ]; then
        cd "{{.SRC}}/solvent" && git fetch && git checkout {{.SOLVENT_REV}}
      else
        git clone {{.SOLVENT_REPO_URL}} "{{.SRC}}/solvent"
        cd "{{.SRC}}/solvent" && git checkout {{.SOLVENT_REV}}
      fi
    # Conductor
    - |
      if [ -d "{{.SRC}}/conductor" ]; then
        cd "{{.SRC}}/conductor" && git fetch && git checkout {{.CONDUCTOR_REV}}
      else
        git clone {{.CONDUCTOR_REPO_URL}} "{{.SRC}}/conductor"
        cd "{{.SRC}}/conductor" && git checkout {{.CONDUCTOR_REV}}
      fi
```

**Idempotent:** If source exists at correct revision, reuses it.

**`task clean`** does NOT remove `.tmp/src/` (source dependencies persist).
**`task clean:all`** removes everything including `.tmp/src/`.

### All references updated

| Old (v1) | New (v2) |
|----------|----------|
| `SOLVENT_REPO: ../solvent-main` | `SOLVENT_SRC: .tmp/src/solvent` |
| `CONDUCTOR_REPO: ../conductor` | `CONDUCTOR_SRC: .tmp/src/conductor` |
| `go build ... {{.SOLVENT_REPO}}/cmd/solvent-api` | `go build ... {{.SOLVENT_SRC}}/cmd/solvent-api` |
| `go build ... {{.CONDUCTOR_REPO}}/cmd/conductor` | `go build ... {{.CONDUCTOR_SRC}}/cmd/conductor` |

---

## 8. Solvent Migrations (Corrected)

### Principle

> Oracle orchestrates. Solvent owns its migrations.

Oracle does NOT inline migration filenames. It calls Solvent's canonical
`db:reset` via `dir:` override.

### Task: `setup:db`

```yaml
setup:db:
  desc: "Initialize CockroachDB via Solvent's canonical db:reset"
  cmds:
    - task: db:up
    - task: db:wait
    - task: db:solvent-reset

db:solvent-reset:
  desc: "Invoke Solvent's canonical db:reset (all 10 migrations)"
  dir: "{{.SOLVENT_SRC}}"
  cmds:
    - task db:reset
```

**Why `dir:` works:** Task's `dir:` changes the working directory for the
task's `cmds`. Solvent's `db:reset` uses `< db/NNN_name.sql` which resolves
relative to the working directory. With `dir: "{{.SOLVENT_SRC}}"`, the SQL
files resolve correctly.

**Caveat:** Solvent's `db:reset` hardcodes `solvent-crdb` as the container
name. This matches our standardized container name. If Solvent changes this
in the future, the Oracle Taskfile would need updating — this is acceptable
for a POC with pinned revisions.

---

## 9. Full `task setup` Sequence

```yaml
setup:
  desc: "Full POC initialization (destructive/reproducible)"
  cmds:
    - task: setup:check      # verify prerequisites
    - task: setup:dirs       # create .tmp/bin, .tmp/state, .tmp/logs, .tmp/pids
    - task: setup:deps       # clone pinned Solvent + Conductor
    - task: setup:build      # build all binaries
    - task: setup:db         # start CRDB, wait, Solvent db:reset
    - task: setup:seed       # seed principal (after DB ready)
    - task: setup:config     # generate .opencode/config.json
```

**Removed from setup:** `conductor:create-project` (moved to `task up`).

---

## 10. Full `task up` Sequence (Corrected)

```yaml
up:
  desc: "Start all services in dependency order"
  cmds:
    - task: db:up                        # ensure CRDB running
    - task: db:wait                      # wait CRDB ready
    - task: up:solvent                   # start → wait
    - task: up:conductor                 # start → wait
    - task: up:conductor:create-project  # idempotent project creation
    - task: up:coordinator               # start → wait
    - task: up:trust-ui                  # start → wait
```

### Each service: start → wait (idempotent)

```yaml
up:solvent:
  desc: "Start Solvent API (idempotent)"
  cmds:
    # Idempotency: skip if PID file exists and process alive
    - |
      if [ -f {{.PIDS}}/solvent.pid ] && kill -0 $$(cat {{.PIDS}}/solvent.pid) 2>/dev/null; then
        echo "Solvent already running (PID: $$(cat {{.PIDS}}/solvent.pid))"
        exit 0
      fi
      rm -f {{.PIDS}}/solvent.pid
    # Start
    - SOLVENT_DATABASE_URL="{{.SOLVENT_DSN}}" \
      SOLVENT_API_KEYS="{{.SOLVENT_API_KEY}}={{.SOLVENT_OPERATOR_ID}}" \
      SOLVENT_API_ADDR=":9090" \
      nohup {{.BIN}}/solvent-api > {{.LOGS}}/solvent.log 2>&1 &
      echo $! > {{.PIDS}}/solvent.pid
      echo "Solvent started (PID: $$(cat {{.PIDS}}/solvent.pid))"
    # Wait
    - task: wait:solvent
```

```yaml
up:conductor:
  desc: "Start Conductor API (idempotent)"
  cmds:
    - |
      if [ -f {{.PIDS}}/conductor.pid ] && kill -0 $$(cat {{.PIDS}}/conductor.pid) 2>/dev/null; then
        echo "Conductor already running (PID: $$(cat {{.PIDS}}/conductor.pid))"
        exit 0
      fi
      rm -f {{.PIDS}}/conductor.pid
    - CONDUCTOR_API_KEY="{{.SOLVENT_API_KEY}}:operator-001" \
      nohup {{.BIN}}/conductor --mode api --addr :9091 --db {{.STATE}}/conductor.db > {{.LOGS}}/conductor.log 2>&1 &
      echo $! > {{.PIDS}}/conductor.pid
      echo "Conductor started (PID: $$(cat {{.PIDS}}/conductor.pid))"
    - task: wait:conductor
```

```yaml
up:conductor:create-project:
  desc: "Create default project in Conductor (idempotent)"
  cmds:
    - |
      # Check if project already exists
      code=$$(curl -sf -o /dev/null -w '%{http_code}' http://{{.CONDUCTOR_ADDR}}/v1/projects 2>/dev/null || true)
      if [ "$code" != "200" ]; then
        echo "Conductor not ready, skipping project creation"
        exit 0
      fi
      # Attempt creation (idempotent if project already exists)
      curl -sf -X POST http://{{.CONDUCTOR_ADDR}}/v1/projects \
        -H "Content-Type: application/json" \
        -H "X-API-Key: {{.SOLVENT_API_KEY}}" \
        -d '{"name": "default", "description": "ARGUS Phase 8 POC"}' \
        2>/dev/null || true
      echo "Default project ensured"
```

```yaml
up:coordinator:
  desc: "Start Coordinator HTTP (idempotent)"
  cmds:
    - |
      if [ -f {{.PIDS}}/coordinator.pid ] && kill -0 $$(cat {{.PIDS}}/coordinator.pid) 2>/dev/null; then
        echo "Coordinator already running (PID: $$(cat {{.PIDS}}/coordinator.pid))"
        exit 0
      fi
      rm -f {{.PIDS}}/coordinator.pid
    - SOLVENT_URL="http://{{.SOLVENT_ADDR}}" \
      CONDUCTOR_URL="http://{{.CONDUCTOR_ADDR}}" \
      ARGUS_OPERATOR_PRINCIPAL_ID="{{.SOLVENT_OPERATOR_ID}}" \
      ARGUS_OPERATOR_TOKEN="{{.SOLVENT_API_KEY}}" \
      COORDINATOR_ADDR=":8080" \
      PACK_DIR="$(pwd)/domain-pack" \
      nohup {{.BIN}}/coordinator > {{.LOGS}}/coordinator.log 2>&1 &
      echo $! > {{.PIDS}}/coordinator.pid
      echo "Coordinator started (PID: $$(cat {{.PIDS}}/coordinator.pid))"
    - task: wait:coordinator
```

```yaml
up:trust-ui:
  desc: "Start Trust UI (idempotent)"
  cmds:
    - |
      if [ -f {{.PIDS}}/trust-ui.pid ] && kill -0 $$(cat {{.PIDS}}/trust-ui.pid) 2>/dev/null; then
        echo "Trust UI already running (PID: $$(cat {{.PIDS}}/trust-ui.pid))"
        exit 0
      fi
      rm -f {{.PIDS}}/trust-ui.pid
    - COORDINATOR_URL="http://{{.COORDINATOR_ADDR}}" \
      ARGUS_OPERATOR_PRINCIPAL_ID="{{.SOLVENT_OPERATOR_ID}}" \
      ARGUS_OPERATOR_TOKEN="{{.SOLVENT_API_KEY}}" \
      nohup {{.BIN}}/trust-ui > {{.LOGS}}/trust-ui.log 2>&1 &
      echo $! > {{.PIDS}}/trust-ui.pid
      echo "Trust UI started (PID: $$(cat {{.PIDS}}/trust-ui.pid))"
    - task: wait:trust-ui
```

---

## 11. Readiness Probes

| Service | Probe | Method | Timeout | Retry |
|---------|-------|--------|---------|-------|
| CockroachDB | `docker exec ... cockroach sql -e "SELECT 1"` | Shell | 60s | 1s × 60 |
| Solvent | `GET /v1/ledger?scenario_id=track1` + Bearer token | HTTP | 30s | 1s × 30 |
| Conductor | `GET /v1/projects` | HTTP | 30s | 1s × 30 |
| Coordinator | `GET /packs/rules` | HTTP | 30s | 1s × 30 |
| Trust UI | `GET /insights` | HTTP | 10s | 1s × 10 |

**Trust UI readiness is now part of `task up`** (was missing in v1).

---

## 12. Port Map (Frozen)

```text
CockroachDB SQL  localhost:26260
CockroachDB HTTP localhost:8082   (admin UI, not used)
Solvent API      localhost:9090
Conductor API    localhost:9091
Coordinator      localhost:8080
Trust UI         localhost:8081
MCP adapter      stdio (no port)
```

---

## 13. Testing Strategy (Clarified)

### `task test` — Oracle + Trust UI (no DB required)

```yaml
test:
  desc: "Run Oracle + Trust UI unit tests"
  cmds:
    - go test ./...
    - go test -race ./...
    - go vet ./...
    - cd trust-ui && go test ./...
```

### `task test:all` — Full stack (requires running stack)

```yaml
test:all:
  desc: "Run all tests including Solvent and Conductor"
  cmds:
    - task: test
    - task: test:solvent
    - task: test:conductor
    - task: test:integration

test:solvent:
  desc: "Run Solvent unit tests"
  dir: "{{.SOLVENT_SRC}}"
  cmds:
    - go test ./...

test:conductor:
  desc: "Run Conductor unit tests"
  dir: "{{.CONDUCTOR_SRC}}"
  cmds:
    - go test ./...

test:integration:
  desc: "Run integration tests (requires running stack)"
  cmds:
    - go test -tags=integration ./...
```

### `task test:integration` — Oracle integration only

```yaml
test:integration:
  desc: "Run Oracle integration tests (requires running stack)"
  cmds:
    - go test -tags=integration ./...
```

---

## 14. Process/State Management

### Idempotent startup

Every `up:*` target checks PID file + `kill -0` before starting.
Re-running `task up` is safe.

### Stale PID handling

```yaml
down:solvent:
  cmds:
    - |
      if [ -f {{.PIDS}}/solvent.pid ]; then
        pid=$$(cat {{.PIDS}}/solvent.pid)
        if kill -0 "$$pid" 2>/dev/null; then
          kill "$$pid"
          echo "Solvent stopped (PID: $$pid)"
        else
          echo "Solvent PID $$pid already exited"
        fi
        rm -f {{.PIDS}}/solvent.pid
      fi
```

### `task clean` — safe

```yaml
clean:
  desc: "Remove local state (safe — preserves source deps)"
  cmds:
    - task: down
    - rm -rf {{.BIN}} {{.STATE}} {{.LOGS}} {{.PIDS}}
```

### `task clean:all` — explicitly destructive

```yaml
clean:all:
  desc: "Remove everything including source dependencies and Docker"
  cmds:
    - task: clean
    - rm -rf {{.SRC}}
    - docker rm -f {{.CRDB_CONTAINER}} 2>/dev/null || true
```

---

## 15. Dry-Run Procedure (Corrected MCP Path)

```text
task setup && task up && task status && task verify
```

**Terminal A — Work Agent:**

```bash
opencode  # uses .opencode/config.json with argus MCP
```

Prompt:
```text
Research BM-IST Gate G0. First call argus.get_context with
the current task ID. Then perform bounded research and submit
your findings via argus.submit_packet.
```

Flow:
```
OpenCode → argus-mcp (stdio) → HTTP → Coordinator :8080 → Solvent + Conductor
```

**Terminal B — Adversarial Agent (after Work Agent submits):**

```bash
opencode  # fresh process, no shared state
```

Prompt:
```text
Review all existing G0 research. Call argus.get_context to
reconstruct current state. Challenge weak claims via
argus.submit_packet with role "adversarial".
```

**Browser — Trust UI:** `http://localhost:8081`

Observe Branch A (Dead End) and Branch B (Debt Discharge → Promotion).

---

## 16. Acceptance Criteria

### Must pass

- [ ] Fresh checkout: `git clone → task setup → task up → task status` produces all READY
- [ ] Dependencies obtained reproducibly (pinned Solvent/Conductor in `.tmp/src/`)
- [ ] All services start in correct order (start → wait → next)
- [ ] All readiness checks pass (no sleep-only)
- [ ] Re-running `task up` is idempotent (no duplicate processes)
- [ ] `task test` passes (Oracle + Trust UI unit tests)
- [ ] `task verify` executes `cmd/verifier` binary against G0 fixture
- [ ] `cmd/argus-mcp` builds, connects to Coordinator over HTTP
- [ ] MCP exposes exactly 2 tools: `argus.get_context`, `argus.submit_packet`
- [ ] MCP does NOT expose Solvent/Conductor mutation tools
- [ ] `cmd/argus-mcp` does NOT construct an in-process Coordinator
- [ ] `cmd/argus-mcp` env vars are ONLY `ARGUS_COORDINATOR_URL`
- [ ] Dry-run instructions are reproducible (two OpenCode processes + Trust UI)
- [ ] `task down` shuts down in reverse order, no orphan processes
- [ ] `task clean` removes `.tmp/` state, `task clean:all` also removes `.tmp/src/`

### Must NOT exist

- [ ] No `coordinator.New()` call in `cmd/argus-mcp`
- [ ] No SOLVENT_URL or CONDUCTOR_URL in MCP env vars
- [ ] No `../solvent-main` or `../conductor` references
- [ ] No wait-before-start circularity
- [ ] No project creation before Conductor readiness
- [ ] No architecture boundary weakened

---

## 17. Risks / Known Limitations

| Risk | Mitigation |
|------|-----------|
| Solvent's `db:reset` hardcodes container name `solvent-crdb` | Our standardized container uses the same name |
| Solvent's `db:reset` uses CWD-relative SQL paths | `dir:` override in Taskfile resolves this |
| Solvent API has no unauthenticated health endpoint | Authenticated probe (GET /v1/ledger with Bearer token) |
| Conductor has no health endpoint | GET /v1/projects (unauthenticated, touches DB) |
| PID-based process management is fragile | Stale PID handling in `down:*` targets |
| Solvent/Conductor pinned revisions may need updating | Pin in Taskfile vars, explicit `git fetch` on re-run |

---

## 18. Explicit Non-Goals

- Merging Solvent or Conductor into Oracle
- Embedding Solvent DB access into Coordinator
- Giving agents Solvent/Conductor mutation tools
- Creating a second research database
- Making Taskfile bypass service APIs
- Making dry run autonomous (human adjudication required)
- Process management beyond Taskfile
- Production deployment
