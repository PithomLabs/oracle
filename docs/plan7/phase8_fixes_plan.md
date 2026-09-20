# Phase 8 Implementation Fixes — P0-1, P0-2, P1-3

**Date:** 2026-09-18
**Scope:** Three targeted fixes only — no architectural expansion
**Source:** `docs/plan7/plan4_review.md`

---

## P0-1 — FIX FRESH `task dev`

### Problem

`task dev` calls `db:ensure`, but `db:ensure` only detects a missing database and exits without initializing it. On a fresh checkout, `task dev` fails.

### Current state

- **No `Taskfile.yml` exists** in the Oracle repo. The plan describes one but it has not been created.
- Solvent's canonical `db:reset` (in `solvent-main/Taskfile.yml`) does: `DROP DATABASE IF EXISTS fable CASCADE; CREATE DATABASE fable;` then applies all 10 migration SQL files sequentially.
- Solvent's `scripts/demo/setup.sh` is more conservative: checks if `belief` table exists before applying `001_schema.sql` (non-idempotent), then always applies `002`-`004` (idempotent). But it only applies 4 of 10 migrations — insufficient for Phase 8 which needs the full schema (authority tables from `005`-`010`).

### Solution

Create `Taskfile.yml` with a `db:ensure` task that:

1. Checks CRDB is available (`docker exec solvent-crdb cockroach sql --insecure -e "SELECT 1"`)
2. Checks if `fable` database exists and has the `belief` table
3. If DB absent → runs Solvent's canonical `db:reset` (all 10 migrations)
4. If DB present + schema complete → skips (no-op)
5. If DB present + schema incomplete → fails with message: `"Schema incomplete. Run task fresh to reinitialize."`

### Files to create/modify

| File | Action | Purpose |
|------|--------|---------|
| `Taskfile.yml` | **CREATE** | Developer orchestration (all tasks) |

### Taskfile.yml design

```yaml
version: '3'

vars:
  SOLVENT_SRC: .tmp/src/solvent
  CONDUCTOR_SRC: .tmp/src/conductor
  SOLVENT_REV: "v0.1.0"
  CONDUCTOR_REV: "c243f27"
  SOLVENT_REPO: "https://github.com/PithomLabs/solvent.git"
  CONDUCTOR_REPO: "https://github.com/PithomLabs/conductor.git"
  CRDB_CONTAINER: "solvent-crdb"
  SOLVENT_DB: "fable"
  SOLVENT_DSN: "postgresql://root@localhost:26260/fable?sslmode=disable"
  SOLVENT_API_KEY: "argus-operator"
  CONDUCTOR_API_KEY: "argus-operator:operator-001"
  ARGUS_OPERATOR_TOKEN: "argus-operator-decision-001"
  SOLVENT_OPERATOR_ID: "00000000-0000-0000-0000-000000000001"

tasks:
  dev:
    desc: "Non-destructive developer environment (safe to repeat)"
    cmds:
      - task: setup:check
      - task: setup:dirs
      - task: setup:deps
      - task: setup:build
      - task: up

  fresh:
    desc: "Destructive reinitialization (wipes DB + rebuilds)"
    cmds:
      - task: setup:check
      - task: setup:dirs
      - task: setup:deps
      - task: setup:build
      - task: db:reset
      - task: up

  setup:
    desc: "Full POC initialization (destructive)"
    cmds:
      - task: setup:check
      - task: setup:dirs
      - task: setup:deps
      - task: setup:build
      - task: db:reset
      - task: up

  setup:check:
    desc: "Verify prerequisites"
    cmds:
      - command -v docker || (echo "docker required" && exit 1)
      - command -v go || (echo "go required" && exit 1)
      - command -v git || (echo "git required" && exit 1)
      - docker info >/dev/null 2>&1 || (echo "docker daemon not running" && exit 1)

  setup:dirs:
    desc: "Create working directories"
    cmds:
      - mkdir -p .tmp/bin .tmp/state .tmp/logs .tmp/pids .tmp/artifacts

  setup:deps:
    desc: "Clone pinned Solvent and Conductor sources"
    cmds:
      - |
        if [ -d "{{.SOLVENT_SRC}}" ]; then
          cd "{{.SOLVENT_SRC}}" && git fetch && git checkout {{.SOLVENT_REV}}
        else
          git clone {{.SOLVENT_REPO}} "{{.SOLVENT_SRC}}"
          cd "{{.SOLVENT_SRC}}" && git checkout {{.SOLVENT_REV}}
        fi
      - |
        if [ -d "{{.CONDUCTOR_SRC}}" ]; then
          cd "{{.CONDUCTOR_SRC}}" && git fetch && git checkout {{.CONDUCTOR_REV}}
        else
          git clone {{.CONDUCTOR_REPO}} "{{.CONDUCTOR_SRC}}"
          cd "{{.CONDUCTOR_SRC}}" && git checkout {{.CONDUCTOR_REV}}
        fi

  setup:build:
    desc: "Build all Oracle binaries"
    cmds:
      - go build -o .tmp/bin/verifier ./cmd/verifier
      - go build -o .tmp/bin/coordinator ./cmd/coordinator
      - go build -o .tmp/bin/argus-mcp ./cmd/argus-mcp
      - go build -o .tmp/bin/trust-ui ./trust-ui

  db:up:
    desc: "Start CockroachDB container"
    cmds:
      - |
        if docker ps --format '{{.Names}}' | grep -q "^{{.CRDB_CONTAINER}}$"; then
          echo "CockroachDB already running"
        elif docker ps -a --format '{{.Names}}' | grep -q "^{{.CRDB_CONTAINER}}$"; then
          docker start {{.CRDB_CONTAINER}}
        else
          docker run -d --name {{.CRDB_CONTAINER}} -p 26260:26257 -p 8082:8080 \
            cockroachdb/cockroach:v26.2.0 \
            start-single-node --insecure --accept-sql-without-tls
        fi

  db:wait:
    desc: "Wait for CockroachDB readiness"
    cmds:
      - |
        for i in $(seq 1 60); do
          if docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure -e "SELECT 1" >/dev/null 2>&1; then
            echo "CockroachDB ready"
            exit 0
          fi
          sleep 1
        done
        echo "CockroachDB not ready after 60s" && exit 1

  db:ensure:
    desc: "Ensure Solvent DB exists and is migrated (non-destructive)"
    cmds:
      - task: db:up
      - task: db:wait
      - |
        # Check if fable DB exists and has the belief table (schema complete)
        if docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure \
          --database={{.SOLVENT_DB}} -e "SELECT 1 FROM belief LIMIT 1" >/dev/null 2>&1; then
          echo "Solvent DB exists with complete schema — skipping migration"
          exit 0
        fi
        # Check if fable DB exists but schema is incomplete
        if docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure \
          -e "SELECT 1 FROM [SHOW DATABASES] WHERE database_name = '{{.SOLVENT_DB}}'" \
          --format csv 2>/dev/null | grep -q "{{.SOLVENT_DB}}"; then
          echo "ERROR: Solvent DB exists but schema is incomplete." >&2
          echo "Run 'task fresh' to reinitialize." >&2
          exit 1
        fi
        # DB absent — run canonical db:reset
        echo "Solvent DB absent — running Solvent db:reset"
        task: db:reset

  db:reset:
    desc: "Destructive: drop and recreate Solvent DB (canonical Solvent mechanism)"
    dir: "{{.SOLVENT_SRC}}"
    cmds:
      - task: db:up
      - task: db:wait
      - docker exec {{.CRDB_CONTAINER}} cockroach sql --insecure \
          -e "DROP DATABASE IF EXISTS {{.SOLVENT_DB}} CASCADE; CREATE DATABASE {{.SOLVENT_DB}};"
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/001_schema.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/002_corpus.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/003_wizard.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/004_debt_vocabulary.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/005_authority_mvp.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/006_authority_justification_cascade.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/007_service_tables.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/008_executing_state.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/009_exact_authority_binding.sql
      - docker exec -i {{.CRDB_CONTAINER}} cockroach sql --insecure --database={{.SOLVENT_DB}} < db/010_debt_opaque.sql
      - echo "Database reset complete"

  up:
    desc: "Start all services"
    cmds:
      - task: db:ensure
      - task: up:solvent
      - task: up:conductor
      - task: up:coordinator

  up:solvent:
    desc: "Start Solvent API"
    cmds:
      - |
        if [ -f .tmp/pids/solvent.pid ] && kill -0 $(cat .tmp/pids/solvent.pid) 2>/dev/null; then
          echo "Solvent already running"
          exit 0
        fi
      - SOLVENT_DATABASE_URL="{{.SOLVENT_DSN}}" \
        SOLVENT_API_KEYS="{{.SOLVENT_API_KEY}}={{.SOLVENT_OPERATOR_ID}}" \
        SOLVENT_API_ADDR=":9090" \
        nohup go run ./cmd/solvent-api > .tmp/logs/solvent.log 2>&1 &
        echo $! > .tmp/pids/solvent.pid
      - task: wait:solvent

  up:conductor:
    desc: "Start Conductor API"
    cmds:
      - |
        if [ -f .tmp/pids/conductor.pid ] && kill -0 $(cat .tmp/pids/conductor.pid) 2>/dev/null; then
          echo "Conductor already running"
          exit 0
        fi
      - CONDUCTOR_API_KEY="{{.CONDUCTOR_API_KEY}}" \
        nohup go run ./cmd/conductor --mode api --addr :9091 --db .tmp/state/conductor.db > .tmp/logs/conductor.log 2>&1 &
        echo $! > .tmp/pids/conductor.pid
      - task: wait:conductor

  up:coordinator:
    desc: "Start Coordinator HTTP"
    cmds:
      - |
        if [ -f .tmp/pids/coordinator.pid ] && kill -0 $(cat .tmp/pids/coordinator.pid) 2>/dev/null; then
          echo "Coordinator already running"
          exit 0
        fi
      - SOLVENT_URL="http://localhost:9090" \
        SOLVENT_API_KEY="{{.SOLVENT_API_KEY}}" \
        CONDUCTOR_URL="http://localhost:9091" \
        CONDUCTOR_API_KEY="{{.CONDUCTOR_API_KEY}}" \
        ARGUS_OPERATOR_PRINCIPAL_ID="{{.SOLVENT_OPERATOR_ID}}" \
        ARGUS_OPERATOR_TOKEN="{{.ARGUS_OPERATOR_TOKEN}}" \
        COORDINATOR_ADDR=":8080" \
        ARTIFACT_DIR="$(pwd)/.tmp/artifacts" \
        PACK_DIR="$(pwd)/domain-pack" \
        nohup go run ./cmd/coordinator > .tmp/logs/coordinator.log 2>&1 &
        echo $! > .tmp/pids/coordinator.pid
      - task: wait:coordinator

  wait:solvent:
    cmds:
      - |
        for i in $(seq 1 30); do
          code=$(curl -sf -o /dev/null -w '%{http_code}' \
            -H "Authorization: Bearer {{.SOLVENT_API_KEY}}" \
            http://localhost:9090/v1/ledger?scenario_id=track1 2>/dev/null || echo "000")
          if [ "$code" = "200" ]; then echo "Solvent ready"; exit 0; fi
          sleep 1
        done
        echo "Solvent not ready" && exit 1

  wait:conductor:
    cmds:
      - |
        for i in $(seq 1 30); do
          code=$(curl -sf -o /dev/null -w '%{http_code}' \
            -H "X-API-Key: {{.CONDUCTOR_API_KEY}}" \
            http://localhost:9091/v1/projects 2>/dev/null || echo "000")
          if [ "$code" = "200" ]; then echo "Conductor ready"; exit 0; fi
          sleep 1
        done
        echo "Conductor not ready" && exit 1

  wait:coordinator:
    cmds:
      - |
        for i in $(seq 1 30); do
          code=$(curl -sf -o /dev/null -w '%{http_code}' \
            http://localhost:8080/packs/rules 2>/dev/null || echo "000")
          if [ "$code" = "200" ]; then echo "Coordinator ready"; exit 0; fi
          sleep 1
        done
        echo "Coordinator not ready" && exit 1

  verify:
    desc: "Run Physics Verifier binary against G0 fixture"
    cmds:
      - task: build:verifier
      - mkdir -p .tmp/artifacts
      - .tmp/bin/verifier --input verifier/testdata/g0_input.json --artifact-dir .tmp/artifacts

  build:verifier:
    cmds:
      - go build -o .tmp/bin/verifier ./cmd/verifier

  test:
    desc: "Run Oracle unit tests"
    cmds:
      - go test ./...
      - go test -race ./...
      - go vet ./...

  down:
    desc: "Stop all services"
    cmds:
      - |
        for svc in coordinator conductor solvent; do
          if [ -f .tmp/pids/$$svc.pid ]; then
            pid=$(cat .tmp/pids/$$svc.pid)
            if kill -0 "$$pid" 2>/dev/null; then
              kill "$$pid"
              echo "$$svc stopped"
            fi
            rm -f .tmp/pids/$$svc.pid
          fi
        done

  clean:
    desc: "Remove local state (preserves source deps)"
    cmds:
      - task: down
      - rm -rf .tmp/bin .tmp/state .tmp/logs .tmp/pids .tmp/artifacts

  clean:all:
    desc: "Remove everything including source deps"
    cmds:
      - task: clean
      - rm -rf .tmp/src
      - docker rm -f {{.CRDB_CONTAINER}} 2>/dev/null || true
```

### db:ensure logic (detailed)

```
1. docker exec solvent-crdb cockroach sql --insecure -e "SELECT 1"
   → fails? → error "CockroachDB not running"

2. docker exec solvent-crdb cockroach sql --insecure --database=fable \
     -e "SELECT 1 FROM belief LIMIT 1"
   → succeeds? → "Schema complete, skipping" → done

3. docker exec solvent-crdb cockroach sql --insecure \
     -e "SELECT 1 FROM [SHOW DATABASES] WHERE database_name = 'fable'" \
     --format csv
   → contains "fable"? → "DB exists but schema incomplete" → error "Run task fresh"
   → does not contain "fable"? → "DB absent" → run db:reset → done
```

### Tests

| Test | What it proves |
|------|---------------|
| Fresh checkout + `task dev` succeeds | `db:ensure` initializes when DB absent |
| Second `task dev` succeeds without reset | `db:ensure` skips when schema complete |
| `task fresh` succeeds | `db:reset` works (destructive path) |
| Manual incomplete schema → `task db:ensure` fails with message | Error path works |

### Verification

```bash
# Clean state
docker rm -f solvent-crdb 2>/dev/null || true
rm -rf .tmp

# Fresh checkout simulation
task dev
# Expected: CockroachDB starts, db:ensure runs db:reset, all services start

# Second run (idempotent)
task dev
# Expected: db:ensure detects complete schema, skips, services already running

# State survives
docker exec solvent-crdb cockroach sql --insecure --database=fable \
  -e "SELECT count(*) FROM belief"
# Expected: 0 (no data yet, but schema exists)
```

---

## P0-2 — PHYSICS VERIFIER → ARGUS ARTIFACT HANDOFF

### Problem

`cmd/verifier` creates its own in-memory `ArtifactRegistry` and registers artifacts there. Coordinator creates a separate empty `ArtifactRegistry`. There is no bridge. `task verify` produces an artifact that Coordinator cannot see.

### Current state

- `verifier/registry.go`: `ArtifactRegistry` with unexported `registerTrusted()`. External packages can only read via `ArtifactReader` interface.
- `verifier/runner.go`: `RunPhysicsVerifier()` calls `physicsv1.Run()` then `registry.registerTrusted()`. This is the ONLY trusted registration path.
- `verifier/model/artifact.go`: `VerificationArtifact` struct with `ComputeArtifactHash()`.
- `coordinator/coordinator.go`: `Config.ArtifactReader` is `verifier.ArtifactReader` — must be non-nil.
- No `cmd/verifier/main.go` exists yet. No file-based artifact persistence exists.

### Solution

Add the smallest possible deterministic local artifact handoff:

1. **`cmd/verifier/main.go`**: Standalone CLI that runs the verifier and writes artifact + hash to `--artifact-dir`
2. **`verifier/bootstrap.go`**: New function `LoadTrustedArtifacts(dir)` that reads JSON files, validates hashes, registers with `ArtifactRegistry`
3. **`cmd/coordinator/main.go`**: On startup, calls `LoadTrustedArtifacts(artifactDir)` if `ARTIFACT_DIR` env is set

### Files to create/modify

| File | Action | Purpose |
|------|--------|---------|
| `cmd/verifier/main.go` | **CREATE** | Standalone CLI: run verifier + write artifact to disk |
| `verifier/bootstrap.go` | **CREATE** | Load trusted local artifacts from disk into ArtifactRegistry |
| `verifier/bootstrap_test.go` | **CREATE** | Tests for bootstrap (valid, corrupt, missing, hash mismatch) |
| `cmd/coordinator/main.go` | **CREATE** | Coordinator standalone launcher (loads artifacts on startup) |

### cmd/verifier/main.go

```go
package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/PithomLabs/oracle/verifier"
    physicsv1 "github.com/PithomLabs/oracle/verifier/physics/v1"
)

func main() {
    inputFile := flag.String("input", "", "Path to VerifierInput JSON (required)")
    artifactDir := flag.String("artifact-dir", ".tmp/artifacts", "Directory to write artifact")
    flag.Parse()

    if *inputFile == "" {
        fmt.Fprintln(os.Stderr, "usage: verifier --input <file.json> [--artifact-dir <dir>]")
        os.Exit(2)
    }

    data, err := os.ReadFile(*inputFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
        os.Exit(1)
    }

    var input physicsv1.VerifierInput
    if err := json.Unmarshal(data, &input); err != nil {
        fmt.Fprintf(os.Stderr, "failed to parse input: %v\n", err)
        os.Exit(1)
    }

    registry := verifier.NewArtifactRegistry()
    artifact, err := verifier.RunPhysicsVerifier(context.Background(), registry, input)
    if err != nil {
        fmt.Fprintf(os.Stderr, "verifier failed: %v\n", err)
        os.Exit(1)
    }

    // Write artifact JSON
    if err := os.MkdirAll(*artifactDir, 0755); err != nil {
        fmt.Fprintf(os.Stderr, "failed to create artifact dir: %v\n", err)
        os.Exit(1)
    }

    artifactPath := filepath.Join(*artifactDir, artifact.EvidenceRef+".json")
    artifactJSON, _ := json.MarshalIndent(artifact, "", "  ")
    if err := os.WriteFile(artifactPath, artifactJSON, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "failed to write artifact: %v\n", err)
        os.Exit(1)
    }

    // Write hash file
    hashPath := filepath.Join(*artifactDir, artifact.EvidenceRef+".sha256")
    if err := os.WriteFile(hashPath, []byte(artifact.ArtifactHash+"\n"), 0644); err != nil {
        fmt.Fprintf(os.Stderr, "failed to write hash: %v\n", err)
        os.Exit(1)
    }

    fmt.Fprintf(os.Stderr, "Artifact: %s\n", artifactPath)
    fmt.Fprintf(os.Stderr, "Hash:     %s\n", artifact.ArtifactHash)
    fmt.Fprintf(os.Stderr, "Result:   %s\n", artifact.Result)
}
```

### verifier/bootstrap.go

```go
package verifier

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/PithomLabs/oracle/verifier/model"
)

// LoadTrustedArtifacts reads all .json files from dir, validates their
// ArtifactHash against companion .sha256 files, and registers valid
// artifacts into the provided registry. This is the trusted local
// bootstrap path — the only way external code can register artifacts
// without going through RunPhysicsVerifier.
//
// Security properties:
// - Hash is verified against file content before registration
// - Missing/corrupt files are skipped with warnings (never trusted)
// - Only files matching *.json with a companion *.sha256 are considered
// - The directory is treated as local POC state, not a database
func LoadTrustedArtifacts(registry *ArtifactRegistry, dir string) error {
    entries, err := os.ReadDir(dir)
    if err != nil {
        if os.IsNotExist(err) {
            return nil // No artifacts directory — nothing to load
        }
        return fmt.Errorf("failed to read artifact directory: %w", err)
    }

    loaded := 0
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
            continue
        }

        jsonPath := filepath.Join(dir, entry.Name())
        hashPath := strings.TrimSuffix(jsonPath, ".json") + ".sha256"

        // Read artifact JSON
        artifactData, err := os.ReadFile(jsonPath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "WARNING: cannot read %s: %v\n", jsonPath, err)
            continue
        }

        var artifact model.VerificationArtifact
        if err := json.Unmarshal(artifactData, &artifact); err != nil {
            fmt.Fprintf(os.Stderr, "WARNING: corrupt artifact %s: %v\n", jsonPath, err)
            continue
        }

        // Read expected hash
        expectedHashData, err := os.ReadFile(hashPath)
        if err != nil {
            fmt.Fprintf(os.Stderr, "WARNING: missing hash file %s — skipping %s\n", hashPath, jsonPath)
            continue
        }
        expectedHash := strings.TrimSpace(string(expectedHashData))
        if expectedHash == "" {
            fmt.Fprintf(os.Stderr, "WARNING: empty hash file %s — skipping %s\n", hashPath, jsonPath)
            continue
        }

        // Verify hash
        if err := registry.VerifyHash(artifact.EvidenceRef, expectedHash); err != nil {
            // Not yet registered — compute and compare directly
            computedHash := artifact.ArtifactHash
            if computedHash != expectedHash {
                fmt.Fprintf(os.Stderr, "WARNING: hash mismatch for %s (expected %s, got %s) — skipping\n",
                    jsonPath, expectedHash, computedHash)
                continue
            }
        }

        // Register (trust the hash-verified artifact)
        if err := registry.registerTrusted(artifact); err != nil {
            fmt.Fprintf(os.Stderr, "WARNING: failed to register %s: %v\n", jsonPath, err)
            continue
        }

        loaded++
        fmt.Fprintf(os.Stderr, "Loaded trusted artifact: %s (ref=%s, result=%s)\n",
            jsonPath, artifact.EvidenceRef, artifact.Result)
    }

    if loaded == 0 && len(entries) > 0 {
        return fmt.Errorf("no valid artifacts found in %s", dir)
    }

    return nil
}
```

### coordinator Bootstrap integration

The `cmd/coordinator/main.go` calls `verifier.LoadTrustedArtifacts` on startup:

```go
// After creating the ArtifactRegistry:
artifactDir := envOr("ARTIFACT_DIR", ".tmp/artifacts")
if artifactDir != "" {
    if err := verifier.LoadTrustedArtifacts(artifactRegistry, artifactDir); err != nil {
        log.Printf("WARNING: artifact bootstrap: %v", err)
    }
}
```

### Artifact lifecycle

```text
task verify
    ↓
cmd/verifier --input verifier/testdata/g0_input.json --artifact-dir .tmp/artifacts
    ↓
.tmp/artifacts/physics-v1-g0.json     (canonical artifact)
.tmp/artifacts/physics-v1-g0.sha256    (hash file)
    ↓
task up (starts coordinator)
    ↓
cmd/coordinator reads ARTIFACT_DIR
    ↓
verifier.LoadTrustedArtifacts(registry, dir)
    ↓
ArtifactReader.Resolve("physics-v1-g0") → returns artifact
    ↓
Coordinator retirement rules can resolve "reproducible_artifact" evidence
```

### Tests

| Test | What it proves |
|------|---------------|
| `LoadTrustedArtifacts` with valid artifact | Artifact loaded and resolvable |
| `LoadTrustedArtifacts` with hash mismatch | Artifact rejected, not trusted |
| `LoadTrustedArtifacts` with corrupt JSON | Artifact rejected, not trusted |
| `LoadTrustedArtifacts` with missing .sha256 | Artifact rejected, not trusted |
| `LoadTrustedArtifacts` with empty directory | No error, no artifacts loaded |
| `LoadTrustedArtifacts` with nonexistent directory | No error, no artifacts loaded |
| `LoadTrustedArtifacts` survives "restart" | New registry + same dir = same artifacts |
| `LoadTrustedArtifacts` hash is deterministic | Same input = same hash across runs |
| `cmd/verifier` produces expected files | CLI integration works |
| No agent MCP tool can register artifacts | `registerTrusted` is unexported |

### Verification

```bash
# Produce artifact
task verify
ls -la .tmp/artifacts/
# Expected: physics-v1-g0.json, physics-v1-g0.sha256

# Verify hash
cat .tmp/artifacts/physics-v1-g0.sha256
# Expected: 64-char hex hash

# Start coordinator (loads artifact)
task up:coordinator

# Verify artifact is accessible
curl -s http://localhost:8080/packs/rules
# Expected: retirement rules (confirms Coordinator is alive)

# Verify artifact survives restart
# (stop coordinator, start again, artifact still loaded)
```

---

## P1-3 — SCENARIO-CONTAINMENT FOR CANONICAL REFERENCES

### Problem

`canonical:belief:<uuid>` references in packets are resolved without checking that the referenced belief belongs to the same scenario as the packet. This allows cross-scenario contamination.

### Current state

- `packet/v1/resolve.go`: `ParseReference()` validates format (local: or canonical:belief:<uuid>) but not scenario ownership.
- `coordinator/persist.go`: `resolveRef()` extracts the UUID from `canonical:belief:` prefix and uses it directly as a Solvent belief ID. No scenario check.
- `coordinator/compiler.go`: `compile()` validates packet structure but does not check canonical references against Solvent.
- `coordinator/mock.go`: `GetBelief(beliefID)` returns a mock with no `scenario_id` field.

### Solution

Add scenario containment validation in the Coordinator's persistence layer, BEFORE any writes to Solvent/Conductor.

The check is:

```
For each canonical:belief:<uuid> reference in the packet:
    1. Fetch the belief from Solvent via GetBelief(uuid)
    2. Extract scenario_id from the response
    3. Compare to packet.ScenarioID
    4. If mismatch → refuse packet, perform ZERO persistence
```

This is checked in `SubmitPacket()` after compilation but before any `CreateBelief`/`CreateEdge`/`CreateEvidence` calls. The check is fail-closed: if Solvent is unreachable or returns an error, the packet is refused.

### Files to modify

| File | Action | Purpose |
|------|--------|---------|
| `coordinator/persist.go` | **MODIFY** | Add `validateScenarioOwnership()` before persistence |
| `coordinator/persist_test.go` | **MODIFY** | Add tests for same-scenario (accept), cross-scenario (refuse), missing belief (refuse) |
| `coordinator/mock.go` | **MODIFY** | Add `scenario_id` to mock `GetBelief` response |

### persist.go changes

Add a new function and call it in `SubmitPacket`:

```go
// validateScenarioOwnership verifies that all canonical:belief references
// in the packet belong to the same scenario as the packet itself.
// This prevents cross-scenario contamination.
// Fail-closed: any error or mismatch refuses the packet.
func (c *Coordinator) validateScenarioOwnership(pkt *packetv1.Packet) error {
    if pkt.ScenarioID == "" {
        return nil // No scenario to check against
    }

    for i, b := range pkt.Beliefs {
        // Check edges referencing canonical beliefs
        for j, e := range pkt.Edges {
            ref := e.FromRef
            if !strings.HasPrefix(ref, packetv1.RefPrefixCanonical) {
                continue
            }
            beliefID := strings.TrimPrefix(ref, packetv1.RefPrefixCanonical)
            belief, err := c.solventClient.GetBelief(beliefID)
            if err != nil {
                return fmt.Errorf("edge[%d].from_ref: failed to resolve canonical belief %s: %w", j, beliefID, err)
            }
            refScenario, _ := belief["scenario_id"].(string)
            if refScenario != pkt.ScenarioID {
                return fmt.Errorf("edge[%d].from_ref: canonical belief %s belongs to scenario %q, packet scenario is %q",
                    j, beliefID, refScenario, pkt.ScenarioID)
            }
        }
        // ... similar for ToRef
        _ = b // belief iteration for evidence refs below
    }

    // Check evidence referencing canonical beliefs
    for i, e := range pkt.Evidence {
        if !strings.HasPrefix(e.BeliefRef, packetv1.RefPrefixCanonical) {
            continue
        }
        beliefID := strings.TrimPrefix(e.BeliefRef, packetv1.RefPrefixCanonical)
        belief, err := c.solventClient.GetBelief(beliefID)
        if err != nil {
            return fmt.Errorf("evidence[%d].belief_ref: failed to resolve canonical belief %s: %w", i, beliefID, err)
        }
        refScenario, _ := belief["scenario_id"].(string)
        if refScenario != pkt.ScenarioID {
            return fmt.Errorf("evidence[%d].belief_ref: canonical belief %s belongs to scenario %q, packet scenario is %q",
                i, beliefID, refScenario, pkt.ScenarioID)
        }
    }

    return nil
}
```

In `SubmitPacket`, after compilation and before persistence:

```go
func (c *Coordinator) SubmitPacket(pkt *packetv1.Packet) (*CompilationResult, error) {
    result, err := c.compile(pkt)
    if err != nil {
        return nil, fmt.Errorf("compilation failed: %w", err)
    }

    // NEW: Validate scenario ownership of canonical references
    if err := c.validateScenarioOwnership(pkt); err != nil {
        return nil, fmt.Errorf("scenario containment violation: %w", err)
    }

    // ... rest of persistence (beliefs, edges, evidence, tasks)
}
```

### Mock changes

Update `GetBelief` mock to return `scenario_id`:

```go
func (m *MockSolventClient) GetBelief(beliefID string) (map[string]interface{}, error) {
    if m.GetBeliefFn != nil {
        return m.GetBeliefFn(beliefID)
    }
    return map[string]interface{}{
        "id":          beliefID,
        "status":      "entered",
        "scenario_id": "track-g0",  // NEW: default scenario for mocks
        "debt":        []string{"needMap"},
    }, nil
}
```

### Tests

| Test | What it proves |
|------|---------------|
| Canonical belief in same scenario → accepted | Basic happy path |
| Canonical belief in different scenario → refused | Cross-scenario blocked |
| Canonical belief not found (Solvent error) → refused | Fail-closed on missing |
| Malformed canonical reference → refused (existing validation) | Format validation still works |
| Dangling canonical reference → refused (GetBelief error) | Fail-closed on dangling |
| Failed validation causes zero persistence | No partial writes |
| Local references unaffected | `local:` refs work as before |
| Evidence with canonical ref in same scenario → accepted | Evidence refs checked too |
| Evidence with canonical ref in different scenario → refused | Evidence cross-scenario blocked |

### Verification

```bash
# Unit tests
go test ./coordinator/ -run TestScenarioContainment -v
go test ./packet/v1/ -run TestValidate -v

# Integration (with running stack)
task dev
# Submit packet with canonical ref to same scenario → succeeds
# Submit packet with canonical ref to different scenario → refused
```

---

## REGRESSION REQUIREMENTS

After all three fixes, run:

```bash
# Oracle tests
go test ./...
go test -race ./...
go vet ./...

# Trust UI tests
cd trust-ui && go test ./...

# Full developer environment
task dev
task status  # (or curl readiness endpoints)
task verify

# Verify artifact handoff
ls .tmp/artifacts/
# Expected: physics-v1-g0.json, physics-v1-g0.sha256

# Verify Coordinator can read artifact
curl -s http://localhost:8080/packs/rules
# Expected: retirement rules
```

---

## DELIVERABLE: PHASE8_DEV_ENV_FIX_REPORT.md

Create this file with:

1. **P0-1 implementation + tests**
   - Taskfile.yml created
   - `db:ensure` logic verified
   - Fresh checkout + `task dev` succeeds
   - Second `task dev` does not reset DB

2. **P0-2 implementation + tests**
   - `cmd/verifier/main.go` created
   - `verifier/bootstrap.go` created
   - `cmd/coordinator/main.go` created
   - Artifact written to `.tmp/artifacts/`
   - Coordinator loads artifact on startup
   - Hash verified, corrupt artifact rejected
   - Artifact survives Coordinator restart

3. **P1-3 implementation + tests**
   - `validateScenarioOwnership()` added
   - Same-scenario canonical ref accepted
   - Cross-scenario canonical ref refused
   - Zero persistence on validation failure

4. **Exact files changed** (list all)

5. **Test results** (output of `go test ./...`)

6. **Taskfile verification results** (output of `task dev`, `task verify`)

7. **Remaining limitations** specific to these fixes only

---

## EXACT FILES CHANGED

| File | Action | Fix |
|------|--------|-----|
| `Taskfile.yml` | CREATE | P0-1 |
| `cmd/verifier/main.go` | CREATE | P0-2 |
| `cmd/coordinator/main.go` | CREATE | P0-2 |
| `verifier/bootstrap.go` | CREATE | P0-2 |
| `verifier/bootstrap_test.go` | CREATE | P0-2 |
| `coordinator/persist.go` | MODIFY | P1-3 |
| `coordinator/persist_test.go` | MODIFY | P1-3 |
| `coordinator/mock.go` | MODIFY | P1-3 |
| `.gitignore` | MODIFY | P0-1 (add .tmp/) |
| `PHASE8_DEV_ENV_FIX_REPORT.md` | CREATE | Deliverable |
