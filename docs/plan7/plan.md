# PHASE 8 DEV ENV IMPLEMENTATION PLAN

**Date:** 2026-09-18
**Status:** Draft — Awaiting approval before implementation
**Scope:** Developer experience reproducibility for ARGUS Phase 8 dry run

---

## 1. Objective

Make the entire ARGUS trust-verification developer experience seamless
and reproducible:

```text
git clone
    ↓
task setup
    ↓
task up
    ↓
task status
    → CockroachDB READY
    → Solvent READY
    → Conductor READY
    → Coordinator READY
    → Trust UI READY
    ↓
task test
task verify
    ↓
Two fresh OpenCode processes + Trust UI → Phase 8 dry run
```

The developer must NOT need to:
- discover undocumented ports
- manually construct service commands
- manually run database migrations
- manually build Coordinator
- manually discover MCP configuration
- manually wire Solvent and Conductor URLs
- manually create runtime directories
- manually guess readiness
- manually clean up orphan processes

---

## 2. Actual Current Runtime Architecture

### What exists today

| Component | Repository | Binary? | Status |
|-----------|-----------|---------|--------|
| CockroachDB | External (Docker) | Docker image | Taskfile in Solvent starts it |
| Solvent API | solvent-main | `cmd/solvent-api` | Standalone binary, :8080 default |
| Solvent MCP | solvent-main | `cmd/solvent-mcp` | Standalone binary, stdio, 18 tools |
| Conductor | conductor | `cmd/conductor` | Standalone binary, api/mcp/web modes |
| Coordinator | oracle | **NONE** | Library + HTTP handler only |
| Trust UI | oracle/trust-ui | `trust-ui/server.go` | Standalone binary, :8081 |
| MCP Adapter | oracle/mcp/adapter | **NONE** | Library only, no transport |
| Physics Verifier | oracle/verifier | **NONE** | Library only, no binary |
| Reference Loop | oracle/reference-loop | `reference-loop/main.go` | Integration test harness |

### What must be created

| Component | Path | Type | Reason |
|-----------|------|------|--------|
| Coordinator launcher | `oracle/cmd/coordinator/main.go` | New binary | No standalone binary exists |
| MCP stdio server | `oracle/mcp/server.go` | New file | No transport layer exists |
| MCP binary | `oracle/cmd/argus-mcp/main.go` | New binary | No entry point exists |
| Verifier binary | `oracle/cmd/verifier/main.go` | New binary | No standalone binary exists |
| Taskfile | `oracle/Taskfile.yml` | New file | No orchestration exists |
| Config | `oracle/.env.example` | New file | No config exists |

---

## 3. Repository/Service Inventory

### Oracle (`github.com/PithomLabs/oracle`)

| Package | Path | Purpose |
|---------|------|---------|
| coordinator | `coordinator/` | Core coordinator logic |
| coordinator/http | `coordinator/http/` | HTTP handler + server |
| packet/v1 | `packet/v1/` | EBP packet types + validation |
| domain-pack | `domain-pack/` | Pack registry + BM-IST pack |
| verifier | `verifier/` | Artifact registry + physics verifier |
| corpus | `corpus/` | Manifest management |
| mcp/adapter | `mcp/adapter/` | Agent-facing tool logic |
| trust-ui | `trust-ui/` | Human adjudication web UI |
| reference-loop | `reference-loop/` | Integration test harness |

### Solvent (`github.com/PithomLabs/solvent`)

| Component | Path | Purpose |
|-----------|------|---------|
| API server | `cmd/solvent-api/` | REST API for beliefs, evidence, authority |
| MCP server | `cmd/solvent-mcp/` | Agent-facing MCP (18 tools) |
| Schema | `db/001..010_schema.sql` | 10 sequential migrations |
| Taskfile | `Taskfile.yml` | Existing dev orchestration |

### Conductor (`github.com/PithomLabs/conductor`)

| Component | Path | Purpose |
|-----------|------|---------|
| Binary | `cmd/conductor/` | api/mcp/web modes |
| Schema | `migrations/001..002.sql` | 2 embedded migrations |
| Taskfile | `Taskfile.yml` | Build/test/lint |

---

## 4. Required Binaries

### Binaries to create

| Binary | Source | Build Command | Purpose |
|--------|--------|---------------|---------|
| `oracle/cmd/coordinator/main.go` | New | `go build -o .tmp/bin/coordinator ./cmd/coordinator` | Coordinator HTTP server |
| `oracle/cmd/argus-mcp/main.go` | New | `go build -o .tmp/bin/argus-mcp ./cmd/argus-mcp` | MCP stdio adapter |
| `oracle/cmd/verifier/main.go` | New | `go build -o .tmp/bin/verifier ./cmd/verifier` | Physics verifier CLI |

### Binaries to reuse

| Binary | Source | Build Command | Purpose |
|--------|--------|---------------|---------|
| `solvent-api` | solvent-main/cmd/solvent-api | `go build -o .tmp/bin/solvent-api ./cmd/solvent-api` | Solvent REST API |
| `conductor` | conductor/cmd/conductor | `go build -o .tmp/bin/conductor ./cmd/conductor` | Conductor API server |
| `trust-ui` | oracle/trust-ui | `cd trust-ui && go build -o ../.tmp/bin/trust-ui .` | Trust UI web |

---

## 5. Required Databases

| Database | Engine | Location | Purpose |
|----------|--------|----------|---------|
| `fable` | CockroachDB | Docker `solvent-crdb` | Solvent backend |
| `conductor.db` | SQLite | `.tmp/state/conductor.db` | Conductor tasks |

### CockroachDB

- Docker image: `cockroachdb/cockroach:v26.2.0`
- Container name: `solvent-crdb`
- SQL port: `26260` (host) → `26257` (container)
- HTTP port: `8082` (host) → `8080` (container)
- Mode: single-node, insecure, no TLS
- DSN: `postgresql://root@localhost:26260/fable?sslmode=disable`

### SQLite (Conductor)

- File: `.tmp/state/conductor.db`
- Mode: WAL, foreign keys enabled
- Migrations: embedded in binary (001_initial.sql, 002_dependencies.sql)

---

## 6. Coordinator Launcher Plan

### File: `oracle/cmd/coordinator/main.go`

**What it does:** Thin launcher around existing Coordinator library.

**Construction:**

```go
func main() {
    // 1. Read config from env
    solventURL := envOr("SOLVENT_URL", "http://localhost:9090")
    conductorURL := envOr("CONDUCTOR_URL", "http://localhost:9091")
    operatorID := envRequired("ARGUS_OPERATOR_PRINCIPAL_ID")
    operatorToken := envRequired("ARGUS_OPERATOR_TOKEN")
    addr := envOr("COORDINATOR_ADDR", ":8080")
    packDir := envOr("PACK_DIR", "./domain-pack")

    // 2. Load domain packs
    packRegistry := domainpack.NewRegistry()
    packRegistry.LoadFromDisk(packDir)

    // 3. Create artifact registry (in-memory, trusted registration)
    artifactRegistry := verifier.NewArtifactRegistry()

    // 4. Construct Coordinator
    coord, err := coordinator.New(coordinator.Config{
        SolventBaseURL:   solventURL,
        ConductorBaseURL: conductorURL,
        PackRegistry:     packRegistry,
        ArtifactReader:   artifactRegistry.AsReader(),
        OperatorID:       operatorID,
    })

    // 5. Set operator token for /decisions auth
    coord.SetOperatorToken(operatorToken)

    // 6. Create HTTP server
    srv := coordinatorhttp.NewServer(coord)

    // 7. Listen + graceful shutdown
    httpSrv := &http.Server{Addr: addr, Handler: srv.Handler()}
    // ... signal handling, shutdown
}
```

**Environment variables:**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `SOLVENT_URL` | No | `http://localhost:9090` | Solvent API base URL |
| `CONDUCTOR_URL` | No | `http://localhost:9091` | Conductor API base URL |
| `ARGUS_OPERATOR_PRINCIPAL_ID` | Yes | — | Operator principal UUID |
| `ARGUS_OPERATOR_TOKEN` | Yes | — | Bearer token for /decisions |
| `COORDINATOR_ADDR` | No | `:8080` | Listen address |
| `PACK_DIR` | No | `./domain-pack` | Domain pack directory |

**Graceful shutdown:** Listen for `SIGINT`/`SIGTERM`, call `httpSrv.Shutdown(ctx)` with 5s timeout.

**Health:** The `/packs/rules` endpoint (GET, no auth) serves as a de facto readiness probe — it requires a loaded PackRegistry.

**Tests:** The existing `coordinator/http/handler_test.go` tests cover the HTTP handler. Add a `TestCoordinatorServerStartsAndStops` integration test.

---

## 7. Physics Verifier Launcher Plan

### File: `oracle/cmd/verifier/main.go`

**What it does:** Runs the physics verifier independently and prints the artifact.

**Construction:**

```go
func main() {
    // 1. Parse flags
    inputJSON := flag.String("input", "", "Path to VerifierInput JSON file")
    flag.Parse()

    // 2. Read input
    inputBytes, _ := os.ReadFile(*inputJSON)
    var input physicsv1.VerifierInput
    json.Unmarshal(inputBytes, &input)

    // 3. Create artifact registry
    registry := verifier.NewArtifactRegistry()

    // 4. Run verifier
    artifact, err := verifier.RunPhysicsVerifier(ctx, registry, input)

    // 5. Print result
    // Print artifact JSON to stdout
    // Print summary to stderr
    // Exit 0 on confirmed, exit 1 on refuted/inconclusive
}
```

**Expected input:** A JSON file containing `physicsv1.VerifierInput` (coefficient expressions, domain parameters).

**Expected output:** `VerificationArtifact` JSON to stdout, human-readable summary to stderr.

**Tests:** Existing `verifier/physics/v1/verifier_test.go` covers the verification logic. Add an integration test that runs the binary.

---

## 8. Solvent Startup Plan

### Step 1: Start CockroachDB

```bash
docker run -d \
  --name solvent-crdb \
  -p 26260:26257 \
  -p 8082:8080 \
  cockroachdb/cockroach:v26.2.0 \
  start-single-node --insecure --accept-sql-without-tls
```

### Step 2: Wait for database readiness

Poll `docker exec solvent-crdb cockroach sql --insecure -e "SELECT 1"` with 60 retries, 1s interval.

### Step 3: Apply all 10 migrations (db:reset equivalent)

```bash
docker exec solvent-crdb cockroach sql --insecure \
  -e "DROP DATABASE IF EXISTS fable CASCADE; CREATE DATABASE fable;"

for f in 001_schema 002_corpus 003_wizard 004_debt_vocabulary \
         005_authority_mvp 006_authority_justification_cascade \
         007_service_tables 008_executing_state \
         009_exact_authority_binding 010_debt_opaque; do
  docker exec -i solvent-crdb cockroach sql --insecure --database=fable \
    < solvent-main/db/${f}.sql
done
```

### Step 4: Seed Phase 8 required data

Insert the operator principal (required for discharge attribution):

```bash
docker exec solvent-crdb cockroach sql --insecure --database=fable -e "
  INSERT INTO principal (principal_id, principal_type, issuer, created_at)
  VALUES ('00000000-0000-0000-0000-000000000001', 'service', 'argus-poc', now())
  ON CONFLICT DO NOTHING;
"
```

### Step 5: Start Solvent API

```bash
SOLVENT_DATABASE_URL="postgresql://root@localhost:26260/fable?sslmode=disable" \
SOLVENT_API_KEYS="argus-operator=00000000-0000-0000-0000-000000000001" \
SOLVENT_API_ADDR=":9090" \
.tmp/bin/solvent-api
```

### Step 6: Verify Solvent readiness

Hit any authenticated endpoint:

```bash
curl -sf -H "Authorization: Bearer argus-operator" \
  http://localhost:9090/v1/ledger?scenario_id=track1
```

Retry with 30 attempts, 1s interval. HTTP 200 = ready.

### Step 7: Verify required Phase 8 endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/v1/beliefs` | POST | Enter belief |
| `/v1/beliefs/{id}/edges` | POST | Create edge (growth gate) |
| `/v1/discharge` | POST | Human debt discharge |
| `/v1/beliefs/{id}/promote` | POST | Promotion request |
| `/v1/beliefs/{id}/retract` | POST | Retraction with cascade |
| `/v1/ledger` | GET | Ledger summary |
| `/v1/activity` | GET | Audit activity |

All endpoints exist in the current Solvent API (confirmed by route inspection).

---

## 9. Conductor Startup Plan

### Step 1: Build binary

```bash
go build -o .tmp/bin/conductor ./cmd/conductor
```

(from conductor repo)

### Step 2: Start Conductor API

```bash
CONDUCTOR_API_KEY="argus-operator:operator-001" \
.tmp/bin/conductor --mode api --addr :9091 --db .tmp/state/conductor.db
```

Migrations are embedded and auto-applied on startup.

### Step 3: Verify Conductor readiness

```bash
curl -sf http://localhost:9091/v1/projects
```

Retry with 30 attempts, 1s interval. HTTP 200 = ready.

### Step 4: Create default project

```bash
curl -sf -X POST http://localhost:9091/v1/projects \
  -H "Content-Type: application/json" \
  -d '{"name": "default", "description": "ARGUS Phase 8 POC"}'
```

### Required Conductor features (all confirmed present)

| Feature | Endpoint | Status |
|---------|----------|--------|
| Task lifecycle | POST /v1/tasks/{id}/claim | Implemented |
| Cancellation | POST /v1/tasks/{id}/cancel | Implemented |
| governance_ref | POST /v1/projects/{id}/tasks | Implemented |
| Priority | PATCH /v1/tasks/{id} | Implemented |
| Dependencies | Store layer (no HTTP) | Store only |
| Activity | POST /v1/tasks/{id}/activity | Implemented |

---

## 10. Trust UI Startup Plan

### Build

```bash
cd trust-ui && go build -o ../.tmp/bin/trust-ui .
```

### Run

```bash
COORDINATOR_URL=http://localhost:8080 \
ARGUS_OPERATOR_PRINCIPAL_ID=00000000-0000-0000-0000-000000000001 \
ARGUS_OPERATOR_TOKEN=argus-operator \
.tmp/bin/trust-ui
```

### Verify

```bash
curl -sf http://localhost:8081/insights
```

### Routes

| Route | Purpose |
|-------|---------|
| `/` | Redirect → `/insights` |
| `/insights` | Epistemic state dashboard |
| `/debts` | Debt obligation view + retirement |
| `/api/retire` | POST debt retirement decision |

### Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `COORDINATOR_URL` | No | `http://localhost:8080` | Coordinator base URL |
| `ARGUS_OPERATOR_PRINCIPAL_ID` | Yes | — | Operator identity |
| `ARGUS_OPERATOR_TOKEN` | No | — | Bearer token for /decisions |

---

## 11. MCP/OpenCode Startup Plan

### Architecture

```text
OpenCode
   │
   │ MCP / JSON-RPC over stdio
   ▼
cmd/argus-mcp (stdio process)
   │
   │ HTTP
   ▼
Coordinator :8080
   │
   ├── Solvent :9090
   └── Conductor :9091
```

### File: `oracle/mcp/server.go`

**What it does:** Thin MCP stdio JSON-RPC server using `github.com/modelcontextprotocol/go-sdk`.

**Pattern (reuses Solvent's MCP SDK):**

```go
func Run(ctx context.Context, adapter *adapter.Adapter) error {
    server := mcp.NewServer(&mcp.Implementation{
        Name:    "argus",
        Version: "v0.1.0",
    }, nil)

    for _, tool := range adapter.ListTools() {
        toolCopy := tool
        server.AddTool(&mcp.Tool{
            Name:        toolCopy.Name,
            Description: toolCopy.Description,
            InputSchema: toolCopy.InputSchema,
        }, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            args, _ := json.Marshal(req.Params.Arguments)
            result, err := adapter.HandleTool(ctx, toolCopy.Name, args)
            if err != nil {
                return &mcp.CallToolResult{
                    Content: []mcp.Content{&mcp.TextContent{
                        Text: fmt.Sprintf(`{"error":"%s"}`, err),
                    }}, nil
            }
            data, _ := json.Marshal(result)
            return &mcp.CallToolResult{
                Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
            }, nil
        })
    }

    return server.Run(ctx, &mcp.StdioTransport{})
}
```

**Key design decisions:**
- Log to stderr, never stdout (MCP stdio requirement)
- Use `github.com/modelcontextprotocol/go-sdk` v1.7.0 (same as Solvent)
- Fail fast if Coordinator unreachable (check on startup)
- Tool errors return in MCP result, not as Go errors

### File: `oracle/cmd/argus-mcp/main.go`

```go
func main() {
    coordinatorURL := envOr("ARGUS_COORDINATOR_URL", "http://localhost:8080")

    // Verify Coordinator is reachable
    resp, err := http.Get(coordinatorURL + "/packs/rules")
    if err != nil || resp.StatusCode != 200 {
        log.Fatalf("Coordinator unreachable at %s", coordinatorURL)
    }

    // Construct Coordinator client (for adapter)
    coord, _ := coordinator.New(coordinator.Config{
        SolventBaseURL:   envOr("SOLVENT_URL", "http://localhost:9090"),
        ConductorBaseURL: envOr("CONDUCTOR_URL", "http://localhost:9091"),
        PackRegistry:     loadPacks(),
        ArtifactReader:   verifier.NewArtifactRegistry().AsReader(),
        OperatorID:       envOr("ARGUS_OPERATOR_PRINCIPAL_ID", "operator-001"),
    })

    adapter := adapterpkg.NewAdapter(coord)
    mcpserver.Run(context.Background(), adapter)
}
```

**Environment variables:**

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ARGUS_COORDINATOR_URL` | No | `http://localhost:8080` | Coordinator base URL |
| `SOLVENT_URL` | No | `http://localhost:9090` | Solvent URL (for adapter construction) |
| `CONDUCTOR_URL` | No | `http://localhost:9091` | Conductor URL (for adapter construction) |
| `ARGUS_OPERATOR_PRINCIPAL_ID` | No | `operator-001` | Operator identity |

### OpenCode configuration

Generate `.opencode/config.json` (or equivalent) with:

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

This gives OpenCode exactly two tools: `argus.get_context` and `argus.submit_packet`. No Solvent MCP, no Conductor write tools, no DB access.

---

## 12. Port/Configuration Map

```text
┌─────────────────────────────────────────────────────────┐
│ ARGUS Local Development Port Map                        │
├──────────────────┬──────────────────────────────────────┤
│ CockroachDB SQL  │ localhost:26260                      │
│ CockroachDB HTTP │ localhost:8082 (admin UI)            │
│ Solvent API      │ localhost:9090                       │
│ Conductor API    │ localhost:9091                       │
│ Coordinator      │ localhost:8080                       │
│ Trust UI         │ localhost:8081                       │
│ MCP adapter      │ stdio (no port)                      │
└──────────────────┴──────────────────────────────────────┘
```

**Note:** CockroachDB HTTP admin UI is mapped to 8082 to avoid conflict
with Trust UI on 8081. The admin UI is not used during normal development.

---

## 13. Taskfile Design

### File: `oracle/Taskfile.yml`

```yaml
version: '3'

vars:
  BIN: .tmp/bin
  STATE: .tmp/state
  LOGS: .tmp/logs
  PIDS: .tmp/pids
  SOLVENT_REPO: ../solvent-main
  CONDUCTOR_REPO: ../conductor
  SOLVENT_DSN: "postgresql://root@localhost:26260/fable?sslmode=disable"
  SOLVENT_API_KEY: "argus-operator"
  SOLVENT_OPERATOR_ID: "00000000-0000-0000-0000-000000000001"
  CONDUCTOR_ADDR: "localhost:9091"
  SOLVENT_ADDR: "localhost:9090"
  COORDINATOR_ADDR: "localhost:8080"
  TRUST_UI_PORT: "8081"
  CRDB_CONTAINER: "solvent-crdb"

tasks:
  # ── SETUP ──────────────────────────────────────────
  setup:
    desc: "Full POC initialization (destructive/reproducible)"
    cmds:
      - task: setup:check
      - task: setup:dirs
      - task: setup:build
      - task: setup:db
      - task: setup:seed
      - task: setup:config

  setup:check:
    desc: "Verify prerequisites"
    cmds:
      - which docker || (echo "docker required" && exit 1)
      - which go || (echo "go required" && exit 1)
      - which task || (echo "task required" && exit 1)

  setup:dirs:
    desc: "Create local state directories"
    cmds:
      - mkdir -p {{.BIN}} {{.STATE}} {{.LOGS}} {{.PIDS}}

  setup:build:
    desc: "Build all required binaries"
    cmds:
      - task: build:coordinator
      - task: build:argus-mcp
      - task: build:verifier
      - task: build:trust-ui
      - task: build:solvent-api
      - task: build:conductor

  setup:db:
    desc: "Initialize CockroachDB + apply all migrations"
    cmds:
      - task: db:up
      - task: db:wait
      - task: db:reset

  setup:seed:
    desc: "Seed Phase 8 required data"
    cmds:
      - task: db:seed-principal
      - task: conductor:create-project

  setup:config:
    desc: "Generate local configuration"
    cmds:
      - task: config:generate

  # ── BUILD ──────────────────────────────────────────
  build:coordinator:
    desc: "Build Coordinator binary"
    cmds:
      - go build -o {{.BIN}}/coordinator ./cmd/coordinator

  build:argus-mcp:
    desc: "Build ARGUS MCP adapter binary"
    cmds:
      - go build -o {{.BIN}}/argus-mcp ./cmd/argus-mcp

  build:verifier:
    desc: "Build Physics Verifier binary"
    cmds:
      - go build -o {{.BIN}}/verifier ./cmd/verifier

  build:trust-ui:
    desc: "Build Trust UI binary"
    cmds:
      - cd trust-ui && go build -o ../{{.BIN}}/trust-ui .

  build:solvent-api:
    desc: "Build Solvent API binary"
    dir: "{{.SOLVENT_REPO}}"
    cmds:
      - go build -o ../oracle/{{.BIN}}/solvent-api ./cmd/solvent-api

  build:conductor:
    desc: "Build Conductor binary"
    dir: "{{.CONDUCTOR_REPO}}"
    cmds:
      - go build -o ../oracle/{{.BIN}}/conductor ./cmd/conductor

  # ── DATABASE ───────────────────────────────────────
  db:up:
    desc: "Start CockroachDB container"
    cmds:
      - |
        if docker ps --format '{{ "{{" }}.Names{{ "}}"' | grep -q '^{{.CRDB_CONTAINER}}$$'; then
          echo "CockroachDB already running"
        elif docker ps -a --format '{{ "{{" }}.Names{{ "}}"' | grep -q '^{{.CRDB_CONTAINER}}$$'; then
          docker start {{.CRDB_CONTAINER}}
        else
          docker run -d \
            --name {{.CRDB_CONTAINER}} \
            -p 26260:26257 \
            -p 8082:8080 \
            cockroachdb/cockroach:v26.2.0 \
            start-single-node --insecure --accept-sql-without-tls
        fi

  db:wait:
    desc: "Wait for CockroachDB readiness"
    cmds:
      - |
        echo -n "Waiting for CockroachDB..."
        ready=0
        for i in $(seq 1 60); do
          if docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure -e "SELECT 1" &>/dev/null; then
            ready=1; break
          fi
          sleep 1; echo -n "."
        done
        echo
        [ "$ready" -eq 1 ] || (echo "FAILED" && docker logs {{.CRDB_CONTAINER}} --tail 20 && exit 1)
        echo "CockroachDB: READY"

  db:reset:
    desc: "Drop and recreate fable database with all 10 migrations"
    cmds:
      - docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure -e "DROP DATABASE IF EXISTS fable CASCADE; CREATE DATABASE fable;"
      - for f in 001_schema 002_corpus 003_wizard 004_debt_vocabulary 005_authority_mvp 006_authority_justification_cascade 007_service_tables 008_executing_state 009_exact_authority_binding 010_debt_opaque; do \
          docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database=fable < {{.SOLVENT_REPO}}/db/$${f}.sql; \
        done
        echo "Database reset complete (10 migrations)"

  db:seed-principal:
    desc: "Insert operator principal for Phase 8"
    cmds:
      - docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure --database=fable -e \
          "INSERT INTO principal (principal_id, principal_type, issuer, created_at) VALUES ('{{.SOLVENT_OPERATOR_ID}}', 'service', 'argus-poc', now()) ON CONFLICT DO NOTHING;"

  db:down:
    desc: "Stop CockroachDB container"
    cmds:
      - docker stop {{.CRDB_CONTAINER}} 2>/dev/null || true

  # ── UP/DOWN ────────────────────────────────────────
  up:
    desc: "Start all services in dependency order"
    cmds:
      - task: db:up
      - task: db:wait
      - task: up:solvent
      - task: up:conductor
      - task: up:coordinator
      - task: up:trust-ui

  up:solvent:
    desc: "Start Solvent API"
    cmds:
      - task: wait:solvent
        # If already running, skip
      - SOLVENT_DATABASE_URL="{{.SOLVENT_DSN}}" \
        SOLVENT_API_KEYS="{{.SOLVENT_API_KEY}}={{.SOLVENT_OPERATOR_ID}}" \
        SOLVENT_API_ADDR=":9090" \
        nohup {{.BIN}}/solvent-api > {{.LOGS}}/solvent.log 2>&1 &
        echo $! > {{.PIDS}}/solvent.pid
        echo "Solvent started (PID: $$(cat {{.PIDS}}/solvent.pid))"

  up:conductor:
    desc: "Start Conductor API"
    cmds:
      - CONDUCTOR_API_KEY="{{.SOLVENT_API_KEY}}:operator-001" \
        nohup {{.BIN}}/conductor --mode api --addr :9091 --db {{.STATE}}/conductor.db > {{.LOGS}}/conductor.log 2>&1 &
        echo $! > {{.PIDS}}/conductor.pid
        echo "Conductor started (PID: $$(cat {{.PIDS}}/conductor.pid))"

  up:coordinator:
    desc: "Start Coordinator HTTP"
    cmds:
      - SOLVENT_URL="http://{{.SOLVENT_ADDR}}" \
        CONDUCTOR_URL="http://{{.CONDUCTOR_ADDR}}" \
        ARGUS_OPERATOR_PRINCIPAL_ID="{{.SOLVENT_OPERATOR_ID}}" \
        ARGUS_OPERATOR_TOKEN="{{.SOLVENT_API_KEY}}" \
        COORDINATOR_ADDR=":8080" \
        PACK_DIR="./domain-pack" \
        nohup {{.BIN}}/coordinator > {{.LOGS}}/coordinator.log 2>&1 &
        echo $! > {{.PIDS}}/coordinator.pid
        echo "Coordinator started (PID: $$(cat {{.PIDS}}/coordinator.pid))"

  up:trust-ui:
    desc: "Start Trust UI"
    cmds:
      - COORDINATOR_URL="http://{{.COORDINATOR_ADDR}}" \
        ARGUS_OPERATOR_PRINCIPAL_ID="{{.SOLVENT_OPERATOR_ID}}" \
        ARGUS_OPERATOR_TOKEN="{{.SOLVENT_API_KEY}}" \
        nohup {{.BIN}}/trust-ui > {{.LOGS}}/trust-ui.log 2>&1 &
        echo $! > {{.PIDS}}/trust-ui.pid
        echo "Trust UI started (PID: $$(cat {{.PIDS}}/trust-ui.pid))"

  down:
    desc: "Graceful shutdown in reverse dependency order"
    cmds:
      - task: down:trust-ui
      - task: down:coordinator
      - task: down:conductor
      - task: down:solvent
      - task: db:down

  down:trust-ui:
    cmds:
      - test -f {{.PIDS}}/trust-ui && kill $$(cat {{.PIDS}}/trust-ui) 2>/dev/null && rm {{.PIDS}}/trust-ui || true

  down:coordinator:
    cmds:
      - test -f {{.PIDS}}/coordinator && kill $$(cat {{.PIDS}}/coordinator) 2>/dev/null && rm {{.PIDS}}/coordinator || true

  down:conductor:
    cmds:
      - test -f {{.PIDS}}/conductor && kill $$(cat {{.PIDS}}/conductor) 2>/dev/null && rm {{.PIDS}}/conductor || true

  down:solvent:
    cmds:
      - test -f {{.PIDS}}/solvent && kill $$(cat {{.PIDS}}/solvent) 2>/dev/null && rm {{.PIDS}}/solvent || true

  # ── STATUS ─────────────────────────────────────────
  status:
    desc: "Show process status, ports, and readiness"
    cmds:
      - task: status:crdb
      - task: status:solvent
      - task: status:conductor
      - task: status:coordinator
      - task: status:trust-ui

  status:crdb:
    cmds:
      - docker ps --filter name={{.CRDB_CONTAINER}} --format "CockroachDB  {{.Status}}"

  status:solvent:
    cmds:
      - test -f {{.PIDS}}/solvent && kill -0 $$(cat {{.PIDS}}/solvent) 2>/dev/null && echo "Solvent     RUNNING (port 9090)" || echo "Solvent     STOPPED"

  status:conductor:
    cmds:
      - test -f {{.PIDS}}/conductor && kill -0 $$(cat {{.PIDS}}/conductor) 2>/dev/null && echo "Conductor   RUNNING (port 9091)" || echo "Conductor   STOPPED"

  status:coordinator:
    cmds:
      - test -f {{.PIDS}}/coordinator && kill -0 $$(cat {{.PIDS}}/coordinator) 2>/dev/null && echo "Coordinator RUNNING (port 8080)" || echo "Coordinator STOPPED"

  status:trust-ui:
    cmds:
      - test -f {{.PIDS}}/trust-ui && kill -0 $$(cat {{.PIDS}}/trust-ui) 2>/dev/null && echo "Trust UI    RUNNING (port 8081)" || echo "Trust UI    STOPPED"

  # ── WAIT / HEALTH ──────────────────────────────────
  wait:solvent:
    desc: "Wait for Solvent API readiness"
    cmds:
      - |
        echo -n "Waiting for Solvent..."
        ready=0
        for i in $(seq 1 30); do
          code=$$(curl -sf -o /dev/null -w '%{http_code}' -H "Authorization: Bearer {{.SOLVENT_API_KEY}}" http://{{.SOLVENT_ADDR}}/v1/ledger?scenario_id=track1 2>/dev/null || true)
          if [ "$code" = "200" ]; then ready=1; break; fi
          sleep 1; echo -n "."
        done
        echo
        [ "$ready" -eq 1 ] || (echo "Solvent: FAILED" && exit 1)
        echo "Solvent: READY"

  wait:conductor:
    desc: "Wait for Conductor API readiness"
    cmds:
      - |
        echo -n "Waiting for Conductor..."
        ready=0
        for i in $(seq 1 30); do
          code=$$(curl -sf -o /dev/null -w '%{http_code}' http://{{.CONDUCTOR_ADDR}}/v1/projects 2>/dev/null || true)
          if [ "$code" = "200" ]; then ready=1; break; fi
          sleep 1; echo -n "."
        done
        echo
        [ "$ready" -eq 1 ] || (echo "Conductor: FAILED" && exit 1)
        echo "Conductor: READY"

  wait:coordinator:
    desc: "Wait for Coordinator readiness"
    cmds:
      - |
        echo -n "Waiting for Coordinator..."
        ready=0
        for i in $(seq 1 30); do
          code=$$(curl -sf -o /dev/null -w '%{http_code}' http://{{.COORDINATOR_ADDR}}/packs/rules 2>/dev/null || true)
          if [ "$code" = "200" ]; then ready=1; break; fi
          sleep 1; echo -n "."
        done
        echo
        [ "$ready" -eq 1 ] || (echo "Coordinator: FAILED" && exit 1)
        echo "Coordinator: READY"

  # ── LOGS ───────────────────────────────────────────
  logs:
    desc: "Tail all service logs"
    cmds:
      - tail -f {{.LOGS}}/*.log

  logs:solvent:
    cmds:
      - tail -f {{.LOGS}}/solvent.log

  logs:conductor:
    cmds:
      - tail -f {{.LOGS}}/conductor.log

  logs:coordinator:
    cmds:
      - tail -f {{.LOGS}}/coordinator.log

  logs:trust-ui:
    cmds:
      - tail -f {{.LOGS}}/trust-ui.log

  # ── TEST ───────────────────────────────────────────
  test:
    desc: "Run all unit tests (no DB required)"
    cmds:
      - go test ./...
      - go test -race ./...
      - go vet ./...
      - cd trust-ui && go test ./...

  test:integration:
    desc: "Run integration tests (requires running stack)"
    cmds:
      - go test -tags=integration ./...

  # ── VERIFY ─────────────────────────────────────────
  verify:
    desc: "Run Physics Verifier and corpus checks"
    cmds:
      - task: verify:corpus
      - task: verify:verifier

  verify:corpus:
    desc: "Verify corpus manifest integrity"
    cmds:
      - go test -run TestManifest ./corpus/...

  verify:verifier:
    desc: "Run Physics Verifier"
    cmds:
      - go test -run TestPhysicsVerifier ./verifier/...

  # ── DRY RUN ────────────────────────────────────────
  "dry-run":
    desc: "Execute Phase 8 dry run procedure"
    cmds:
      - task: status
      - task: dry-run:check-stack
      - echo ""
      - echo "Phase 8 Dry Run Procedure:"
      - echo "=========================="
      - echo ""
      - echo "1. Open Terminal A — Work Agent:"
      - echo "   opencode --mcp-argus"
      - echo "   Task: 'Research BM-IST Gate G0...'"
      - echo ""
      - echo "2. Wait for Work Agent to submit packet"
      - echo ""
      - echo "3. Open Terminal B — Adversarial Agent:"
      - echo "   opencode --mcp-argus"
      - echo "   Task: 'Review existing G0 work and challenge any weak claims'"
      - echo ""
      - echo "4. Open Browser — Trust UI:"
      - echo "   http://localhost:8081"
      - echo ""
      - echo "5. Observe:"
      - echo "   - /insights: work agent claims appear"
      - echo "   - /insights: adversarial challenge appears"
      - echo "   - Exercise Branch A: RETRACT -> Dead End"
      - echo "   - Exercise Branch B: Discharge debt -> Promotion gate"
      - echo ""
      - echo "Note: Human adjudication is required."

  dry-run:check-stack:
    desc: "Verify all services are running before dry run"
    cmds:
      - task: status:crdb
      - task: status:solvent
      - task: status:conductor
      - task: status:coordinator
      - task: status:trust-ui

  # ── CONFIG ─────────────────────────────────────────
  config:generate:
    desc: "Generate local OpenCode MCP configuration"
    cmds:
      - mkdir -p .opencode
      - |
        cat > .opencode/config.json << 'CONFEOF'
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
        CONFEOF
        echo "Generated .opencode/config.json"

  # ── CLEAN ──────────────────────────────────────────
  clean:
    desc: "Remove POC-local state (safe — never removes source)"
    cmds:
      - task: down
      - rm -rf {{.BIN}} {{.STATE}} {{.LOGS}} {{.PIDS}}

  clean:all:
    desc: "Remove all local state including Docker container"
    cmds:
      - task: clean
      - docker rm -f {{.CRDB_CONTAINER}} 2>/dev/null || true
```

---

## 14. Readiness/Health Strategy

### Per-service readiness probes

| Service | Probe | Method | Timeout | Retry |
|---------|-------|--------|---------|-------|
| CockroachDB | `docker exec ... cockroach sql -e "SELECT 1"` | Shell | 60s | 1s x 60 |
| Solvent | `GET /v1/ledger?scenario_id=track1` with Bearer token | HTTP | 30s | 1s x 30 |
| Conductor | `GET /v1/projects` | HTTP | 30s | 1s x 30 |
| Coordinator | `GET /packs/rules` | HTTP | 30s | 1s x 30 |
| Trust UI | `GET /insights` | HTTP | 10s | 1s x 10 |

### Why these probes

- **CockroachDB:** The only database-level check. Docker exec is reliable for single-node.
- **Solvent:** Every route is behind auth middleware. `GET /v1/ledger` is the lightest authenticated read that touches the DB.
- **Conductor:** Unauthenticated `GET /v1/projects` confirms the API is up and SQLite is initialized.
- **Coordinator:** `GET /packs/rules` confirms the Coordinator is constructed and packs are loaded. No auth required.
- **Trust UI:** `GET /insights` confirms the server is up and can reach Coordinator.

### What is NOT used

- No arbitrary `sleep` durations as the sole readiness mechanism
- No `process.exists()` checks (a process can be up but not ready)
- No custom health endpoints needed (reuse existing endpoints)

---

## 15. Local State/Log Management

### Directory structure

```text
.tmp/
├── bin/                    # Built binaries
│   ├── coordinator
│   ├── argus-mcp
│   ├── verifier
│   ├── trust-ui
│   ├── solvent-api
│   └── conductor
├── state/                  # Persistent local state
│   └── conductor.db        # SQLite database
├── logs/                   # Service logs
│   ├── solvent.log
│   ├── conductor.log
│   ├── coordinator.log
│   └── trust-ui.log
└── pids/                   # Process IDs for lifecycle
    ├── solvent.pid
    ├── conductor.pid
    ├── coordinator.pid
    └── trust-ui.pid
```

### Convention

- `.tmp/` is the single root for all local state
- `task clean` removes `.tmp/` entirely
- `task clean:all` also removes the Docker container
- `.gitignore` should include `.tmp/`
- Never store source code in `.tmp/`

---

## 16. Testing Strategy

### Unit tests (no DB)

```bash
task test
```

Runs `go test ./...` for Oracle core, Trust UI. These are already passing.

### Integration tests (requires running stack)

```bash
task test:integration
```

Requires `task up` to have been run. Tests the Coordinator HTTP handler against live Solvent and Conductor.

### Coordinator standalone binary

| Test | What it verifies |
|------|-----------------|
| Binary builds | `go build ./cmd/coordinator` succeeds |
| Connects to Solvent | Solvent URL resolves, API key works |
| Connects to Conductor | Conductor URL resolves |
| `/context/{task_id}` | RCP context assembly works |
| `/decisions` | Auth middleware works |
| `/packs/rules` | Pack registry loaded |

### Solvent

| Test | What it verifies |
|------|-----------------|
| DB migration | All 10 migrations apply cleanly |
| API startup | `solvent-api` starts on :9090 |
| Phase 8 endpoints | Edge creation, discharge, promotion all work |
| Auth | Bearer token required, invalid token rejected |

### Conductor

| Test | What it verifies |
|------|-----------------|
| Binary builds | `go build ./cmd/conductor` succeeds |
| API startup | Starts on :9091 |
| Task lifecycle | Create -> claim -> submit -> accept works |
| Cancellation | Task cancellation works |
| governance_ref | Set at creation, immutable after |
| Priority | Updateable via PATCH |

### MCP

| Test | What it verifies |
|------|-----------------|
| Binary builds | `go build ./cmd/argus-mcp` succeeds |
| Exactly 2 tools | `tools/list` returns exactly `argus.get_context` and `argus.submit_packet` |
| No authority tools | Solvent/Conductor mutation tools absent |
| OpenCode can connect | JSON-RPC over stdio works |

### Physics Verifier

| Test | What it verifies |
|------|-----------------|
| Binary builds | `go build ./cmd/verifier` succeeds |
| Artifact generated | RunPhysicsVerifier returns valid artifact |
| Hash verified | ArtifactHash matches computed content |
| Exit status | 0 for confirmed, 1 for refuted/inconclusive |

### End-to-end

| Test | What it verifies |
|------|-----------------|
| Work Agent context | `argus.get_context` returns RCP/v1 with task + beliefs |
| Packet persists | `argus.submit_packet` creates beliefs in Solvent, tasks in Conductor |
| Adversarial context | Fresh agent sees existing work via `argus.get_context` |
| Contradicts edge | Adversarial packet creates contradicts edge |
| Retraction path | Human RETRACT -> Solvent cascade -> task cancelled |
| Debt discharge | Human discharge -> Solvent promotion gate |

---

## 17. CI Strategy

### What should run in CI

```yaml
- task test                    # Unit tests (no DB)
- go build ./cmd/coordinator   # Binary builds
- go build ./cmd/argus-mcp
- go build ./cmd/verifier
- go vet ./...
```

### DB-dependent tests in CI

For integration tests requiring CockroachDB:

```yaml
services:
  cockroachdb:
    image: cockroachdb/cockroach:v26.2.0
    ports: ["26260:26257"]
    options: >-
      --name solvent-crdb
      start-single-node --insecure --accept-sql-without-tls

steps:
  - task: db:reset           # Apply all migrations
  - task: db:seed-principal  # Seed required data
  - task: test:integration   # Run DB-backed tests
```

CI should use the same Taskfile targets as developers. No separate CI-only setup.

### Physics Verifier in CI

- Deterministic fixtures (no network dependency)
- Fixed input -> fixed artifact hash
- No CockroachDB required for verifier-only tests

---

## 18. End-to-End Phase 8 Dry-Run Procedure

### Pre-conditions

```bash
task setup    # One-time initialization
task up       # Start all services
task status   # Verify all READY
task verify   # Run verifier + corpus checks
```

### Procedure

**Terminal A — Work Agent:**

```bash
opencode --mcp-argus
```

Prompt:
```text
Research BM-IST Gate G0. Specifically investigate:
- L1: Orbit Rigidity
- L2: One-Parameter Triviality
- Whether countable substrate + literal continuous unitary evolution
  implies structural incompatibility (G0)

First call argus.get_context with the current task ID to understand
what work already exists. Then perform bounded research and submit
your findings via argus.submit_packet.
```

Expected behavior:
1. Agent calls `argus.get_context` -> sees task state
2. Agent performs research
3. Agent calls `argus.submit_packet` with work packet
4. Coordinator persists beliefs, evidence, tasks to Solvent/Conductor

**Inspect Trust UI:**
- Open `http://localhost:8081/insights`
- Observe: work agent beliefs, evidence, and tasks appear
- Open `http://localhost:8081/debts`
- Observe: debt items attached to beliefs

**Terminal B — Adversarial Agent:**

```bash
opencode --mcp-argus
```

Prompt:
```text
Review all existing G0 research in the system. First call
argus.get_context to reconstruct the current state. Then
perform an adversarial review: identify weak claims, missing
evidence, or logical gaps. Submit your challenges via
argus.submit_packet with role "adversarial".
```

Expected behavior:
1. Agent calls `argus.get_context` -> sees ALL existing work
2. Agent identifies challenges
3. Agent calls `argus.submit_packet` with adversarial packet
4. Coordinator creates contradicts edges in Solvent

**Inspect Trust UI:**
- `/insights` now shows adversarial challenges
- Contradicts edges visible

**Exercise Branch A (Dead End):**

In Trust UI `/debts`:
1. Review the contradicted belief
2. Select RETRACT decision
3. Confirm retraction

Expected:
- Solvent RetractCascade fires
- Linked Conductor task cancelled
- `/insights` shows Dead End

**Exercise Branch B (Debt Discharge -> Promotion):**

In Trust UI `/debts`:
1. Review a surviving (non-retracted) belief's debt
2. Select RETIRE_DEBT decision
3. Provide qualifying evidence
4. Confirm discharge

Expected:
- Coordinator validates Pack retirement rule
- Solvent `/v1/discharge` records attribution
- Promotion request sent
- Solvent gate evaluates (debt-free -> promoted)

**The developer observes the full EBP cycle without any component being faked.**

---

## 19. Developer Experience Acceptance Criteria

### Fresh checkout

```bash
git clone <oracle-repo>
cd oracle
task setup
task up
task status
```

Output:

```text
CockroachDB  Up (running)
Solvent      RUNNING (port 9090)
Conductor    RUNNING (port 9091)
Coordinator  RUNNING (port 8080)
Trust UI     RUNNING (port 8081)
```

### Tests pass

```bash
task test       # Unit tests pass
task verify     # Verifier + corpus pass
```

### Dry run is executable

```bash
task dry-run    # Prints procedure with clear instructions
```

Developer opens 3 terminals:
1. Work Agent (OpenCode)
2. Adversarial Agent (OpenCode)
3. Trust UI (browser)

### The developer does NOT need to:

- discover undocumented ports
- manually construct service commands
- manually run database migrations
- manually build Coordinator
- manually discover MCP configuration
- manually wire Solvent and Conductor URLs
- manually create runtime directories
- manually guess readiness
- manually clean up orphan processes

---

## 20. Risks / Known Limitations

### Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Solvent API has no unauthenticated health endpoint | Readiness check requires API key | Use authenticated probe (GET /v1/ledger with Bearer token) |
| Conductor has no health endpoint | Cannot distinguish "starting" from "crashed" | Use GET /v1/projects (unauthenticated, touches DB) |
| Coordinator HTTP handler may fail silently if Solvent/Conductor are down | Partial state observable | Startup checks in Coordinator binary; readiness probe catches this |
| MCP SDK version drift between Solvent and Oracle | Potential protocol incompatibility | Pin to same `github.com/modelcontextprotocol/go-sdk` version |
| CockroachDB Docker image pull on first run | Slow setup | Document in README; cache image |
| PID-based process management is fragile | Orphan processes on crash | `task down` + `task clean` cleanup; could upgrade to process manager later |
| Port 8081 conflict (CockroachDB HTTP + Trust UI) | Confusion | Map CockroachDB HTTP to 8082 in Taskfile |

### Known limitations

| Limitation | Status |
|------------|--------|
| `ADD_DEBT` / automatic new-evidence-creates-debt | Deferred (Phase 8 scope) |
| Full multi-user authentication | Deferred |
| Trust UI automated test coverage | Source-level verification only |
| Cross-ledger happens-before semantics | Not implemented |
| Dependencies have no HTTP API (store layer only) | Conductor limitation |
| Solvent `GET /v1/authorizations/action` doesn't return intent_id | Known workaround in reference-loop |

---

## 21. Explicit Non-Goals

The following are **NOT** part of this plan:

- Merging Solvent into Oracle
- Merging Conductor into Oracle
- Embedding Solvent database access into Coordinator
- Giving agents Solvent MCP tools
- Giving agents Conductor mutation tools
- Creating a second research database
- Adding direct database writes to the Taskfile
- Making Taskfile bypass service APIs
- Making Trust UI bypass Coordinator
- Turning Coordinator into a scientific reasoning engine
- Turning Physics Verifier into a general CAS
- Making the dry run entirely autonomous (human adjudication required)
- Creating a process manager beyond Taskfile
- Production deployment
- Multi-developer collaboration
- Persistent shared state between developers

---

## 22. Proposed File Changes Summary

### New files in `oracle/`

| File | Repository | Reason | Dependencies |
|------|-----------|--------|-------------|
| `cmd/coordinator/main.go` | oracle | Standalone Coordinator binary | coordinator, coordinator/http, domain-pack, verifier |
| `cmd/argus-mcp/main.go` | oracle | MCP stdio adapter entry point | mcp/server, mcp/adapter, coordinator |
| `cmd/verifier/main.go` | oracle | Physics Verifier CLI | verifier, verifier/physics/v1 |
| `mcp/server.go` | oracle | MCP stdio JSON-RPC transport | github.com/modelcontextprotocol/go-sdk |
| `Taskfile.yml` | oracle | Developer orchestration | All above |
| `.env.example` | oracle | Configuration template | — |
| `.gitignore` (update) | oracle | Exclude .tmp/ | — |

### No changes to Solvent or Conductor

The plan reuses existing binaries from Solvent and Conductor as-is. The only Solvent/Conductor dependency is:
- Building `solvent-api` from `solvent-main/cmd/solvent-api`
- Building `conductor` from `conductor/cmd/conductor`
- Applying Solvent schema migrations from `solvent-main/db/`

These are all read-only operations against the existing repositories.
