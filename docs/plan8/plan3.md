# Plan 10 — Corrective Implementation Pass (adv_review4)

**Date:** 2026-09-18
**Source:** `docs/plan8/adv_review4.md` — net-valid defects only
**Architecture constraint:** Preserve the frozen single-process ARGUS architecture. Do NOT restore Solvent/Conductor/Coordinator REST services, separate Trust UI, or internal service-to-service HTTP.

---

## Scope

13 corrective items. Each is a standalone fix with tests. Ordered by dependency — earlier items are prerequisites for later ones.

| # | Fix | Risk | Est. files changed |
|---|-----|------|--------------------|
| 1 | contentHash covers full packet | Low | 1 (app.go) + tests |
| 2 | Solvent migration export | Medium | 2 (solvent + oracle) |
| 3 | Edge kind conflict detection | Low | 1 (app.go) + tests |
| 4 | UNKNOWN ≠ EMPTY in GetContext | Medium | 3 (app.go, view.go, adapter.go) + tests |
| 5 | Retirement-rule enforcement before Discharge | Medium | 2 (app.go, registry fix) + tests |
| 6 | Trust UI authentication | Medium | 2 (ui.go, main.go) + tests |
| 7 | MCP stdio transport (official SDK) | High | 3 (main.go, adapter.go, go.mod) + tests |
| 8 | Semantic version comparison | Low | 1 (app.go) + tests |
| 9 | Server-verify evidence hashes | Low | 1 (app.go) + tests |
| 10 | Port preflight and lifecycle | High | 3 (main.go, Taskfile, tests) |
| 11 | task down + PID cleanup | Medium | 1 (Taskfile) |
| 12 | Integration test CRDB config | Low | 1 (integration_test.go) |
| 13 | Tests for every fix | Low | multiple test files |

---

## Fix 1: contentHash covers full canonical packet

**File:** `internal/application/app.go:78-88`

**Problem:** `contentHash` only hashes `ScenarioID + PacketID + Belief.Claims`. Two packets with identical beliefs but different evidence, edges, or tasks produce the same hash. The `ON CONFLICT (content_hash, scenario_id) DO NOTHING` then silently discards the second packet's distinct entities.

**Fix:** Rewrite `contentHash` to hash the complete canonical packet:
1. `ScenarioID`, `PacketID`, `Role`, `PackRef`
2. Beliefs (sorted by LocalID): `LocalID`, `Claim`, `ClaimType`, sorted `Debt`, `InputSpec`
3. Evidence (sorted by LocalID): `LocalID`, `BeliefRef`, `ProvenanceClass`, `ContentSHA256`, `SourceURL`, `ArtifactRef`
4. Edges (sorted by LocalID): `LocalID`, `FromRef`, `ToRef`, `Kind`
5. Tasks (sorted by Title): `Title`, `Description`, `GovernanceRef`

**Canonical ordering:** Sort each entity list by a deterministic key (LocalID for beliefs/evidence/edges, Title for tasks) before hashing. Use `h.Write([]byte{0})` as field separator and `h.Write([]byte{1})` as entity separator.

**Schema change:** The `submission_idempotency.content_hash` column is `STRING NOT NULL` — no schema change needed.

**Tests:** `TestIdempotencyDuplicatePacket` already exists. Add:
- `TestContentHashDifferentEvidenceProducesDifferentHash` — same beliefs, different evidence → different hashes
- `TestContentHashDifferentEdgesProducesDifferentHash` — same beliefs, different edges → different hashes
- `TestContentHashCanonicalOrdering` — entity order does not affect hash

---

## Fix 2: Solvent migration export

**Problem:** `argus reset` has `// TODO: call solventmigrations.Apply(ctx, db)`. The `solventmigrations` package does not exist in Solvent. Solvent's `internal/` boundary prevents Oracle from importing internal packages.

**Fix (two bounded Solvent changes):**

### Change A: Create `solvent-main/migrations/` package (new top-level package)

Create `/home/chaschel/Documents/go/solvent-main/migrations/migrations.go`:
```go
package solventmigrations

import (
    "context"
    "database/sql"
    "embed"
    "fmt"
    "io/fs"
    "sort"
    "strings"
)

//go:embed db/*.sql
var fs embed.FS

func Apply(ctx context.Context, db *sql.DB) error {
    entries, _ := fs.ReadDir("db")
    var sqlFiles []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
            sqlFiles = append(sqlFiles, e.Name())
        }
    }
    sort.Strings(sqlFiles)
    for _, name := range sqlFiles {
        data, _ := fs.ReadFile("db/" + name)
        for _, stmt := range splitStatements(string(data)) {
            if _, err := db.ExecContext(ctx, stmt); err != nil {
                return fmt.Errorf("apply %s: %w", name, err)
            }
        }
    }
    return nil
}

func splitStatements(sql string) []string { /* same logic as Oracle's */ }
```

Also copy (or symlink) the `db/*.sql` files into `solvent-main/migrations/db/` so they are co-located with the embed directive. Since these are Solvent's own schema files, this is not a schema change — it is a packaging change.

**Verification:** This is Change #1 of 2 in the Solvent change budget. It adds no new tables, no new columns, no new constraints. It only creates a new Go package that embeds existing SQL.

### Change B: Wire into Oracle's `cmdReset`

**File:** `cmd/argus/main.go:182-183`

Replace the TODO with:
```go
import solventmigrations "github.com/PithomLabs/solvent/migrations"

// in cmdReset:
log.Println("argus reset: applying Solvent migrations...")
if err := solventmigrations.Apply(ctx, db); err != nil {
    log.Fatalf("solvent migrations: %v", err)
}
```

The `go.mod` replace directive (`replace github.com/PithomLabs/solvent => ../solvent-main`) already points to the local Solvent directory.

**Tests:**
- `TestCmdResetAppliesSolventMigrations` — run `cmdReset` against an empty CRDB, verify `belief`, `evidence`, `belief_edge` tables exist
- `TestSolventMigrationsAreIdempotent` — apply twice, no errors

---

## Fix 3: Edge kind conflict detection

**File:** `internal/application/app.go:257-264`

**Problem:** `ON CONFLICT (parent_id, child_id) DO NOTHING` silently drops a `contradicts` edge when a `derives` edge already exists for the same pair.

**Fix:** Replace `DO NOTHING` with explicit conflict detection:
```go
// Try insert. On PK conflict, check if kind differs.
_, err := tx.ExecContext(ctx,
    `INSERT INTO belief_edge (parent_id, child_id, kind)
     VALUES ($1::UUID, $2::UUID, $3)
     ON CONFLICT (parent_id, child_id) DO UPDATE
     SET kind = EXCLUDED.kind
     WHERE belief_edge.kind <> EXCLUDED.kind`,
    fromID, toID, edge.Kind)
if err != nil {
    return nil, fmt.Errorf("insert edge: %w", err)
}
```

Wait — this would silently update the kind. Instead, use a two-step approach:

**Option A (preferred):** Query first, then insert or refuse:
```go
var existingKind string
err := tx.QueryRowContext(ctx,
    `SELECT kind FROM belief_edge WHERE parent_id = $1::UUID AND child_id = $2::UUID`,
    fromID, toID).Scan(&existingKind)
if err == nil && existingKind != edge.Kind {
    return nil, fmt.Errorf("edge conflict: %s->%s already has kind=%s, cannot add kind=%s",
        fromID, toID, existingKind, edge.Kind)
}
// Insert (first time or same kind — idempotent)
_, err = tx.ExecContext(ctx,
    `INSERT INTO belief_edge (parent_id, child_id, kind)
     VALUES ($1::UUID, $2::UUID, $3)
     ON CONFLICT (parent_id, child_id) DO NOTHING`,
    fromID, toID, edge.Kind)
```

**Option B (simpler, single statement):** Use `ON CONFLICT ... DO UPDATE` with a `WHERE` guard that returns an error via a check constraint. This is harder in raw SQL.

**Tests:**
- `TestEdgeKindConflictReturnsError` — insert derives, then contradicts on same pair → error
- `TestEdgeSameKindIsIdempotent` — insert derives twice → no error, no duplicate
- `TestEdgeDifferentPairsSucceed` — different parent/child pairs → both stored

---

## Fix 4: UNKNOWN ≠ EMPTY in GetContext

**Files:**
- `internal/application/app.go:298-322` — `GetContext`
- `internal/epistemic/view.go` — `Snapshot`, `Context` types
- `internal/mcp/adapter.go` — `handleGetContext`

**Problem:** `GetContext` returns `nil, error` when a backend is unavailable. The MCP adapter propagates this as a tool-call error. An agent cannot distinguish "no research exists" from "backend is down."

**Fix:** Change `GetContext` to return a `Context` with availability metadata instead of an error for backend failures:

1. Add availability fields to `Context`:
```go
type Context struct {
    Task         *work.Task          `json:"task"`
    Dependencies []*work.Dependency  `json:"dependencies"`
    Snapshot     *epistemic.Snapshot `json:"snapshot"`
    Availability Availability       `json:"availability"`
}

type Availability struct {
    WorkAvailable     bool   `json:"work_available"`
    EpistemicAvailable bool  `json:"epistemic_available"`
    Reason            string `json:"reason,omitempty"`
}
```

2. In `GetContext`, catch errors from each sub-call and set availability:
```go
func (a *App) GetContext(ctx context.Context, taskID string) (*Context, error) {
    result := &Context{Availability: Availability{
        WorkAvailable: true, EpistemicAvailable: true,
    }}

    task, err := a.workStore.GetByID(ctx, taskID)
    if err != nil {
        result.Availability.WorkAvailable = false
        result.Availability.Reason = "work_unavailable"
        // Still return the context — task is nil, but snapshot may work
    } else {
        result.Task = task
    }

    // ... similar for dependencies and snapshot
    // If snapshot query fails:
    //   result.Availability.EpistemicAvailable = false
    //   result.Availability.Reason = "solvent_unavailable"
    //   result.Snapshot = &epistemic.Snapshot{} // empty, not nil

    return result, nil // never return error for availability issues
}
```

3. Only return `error` for truly unrecoverable situations (e.g., invalid taskID format).

**Tests:**
- `TestGetContextWithUnavailableBackend` — mock DB failure → returns Context with `epistemic_available=false`
- `TestGetContextWithEmptyScenario` — valid DB, no beliefs → returns Context with `epistemic_available=true` and empty snapshot
- `TestGetContextNormal` — valid DB, beliefs exist → returns full Context

---

## Fix 5: Retirement-rule enforcement before Discharge

**Files:**
- `internal/application/app.go` — `SubmitDecision`, `App` struct
- `domain-pack/registry.go` — `decodePack` fix

**Problem:** `SubmitDecision("discharge")` calls `kern.Discharge()` directly without checking the pack's `RetirementRules`. The App has no access to packs.

**Fix (two parts):**

### Part A: Fix `decodePack` to return typed packs

The `decodePack` function in `domain-pack/registry.go` always returns `genericPack` which discards all metadata. Fix it to detect BM-IST packs and return `*bmistv1.Pack`:

```go
func decodePack(data []byte, expectedPackID string) (Pack, error) {
    var probe struct {
        PackID  string `json:"pack_id"`
        Version string `json:"version"`
    }
    if err := json.Unmarshal(data, &probe); err != nil {
        return nil, err
    }

    // Try typed decode for known packs
    if strings.EqualFold(probe.PackID, "bmist") {
        return bmistv1.ParsePack(data)
    }

    // Fallback to generic
    return &genericPack{PackID: probe.PackID, Version: probe.Version}, nil
}
```

Add a `retirementRulesProvider` interface:
```go
type retirementRulesProvider interface {
    GetRetirementRules() map[string]interface{}
}
```

### Part B: Give App access to pack registry and enforce rules

Add a `PackRegistry` field to `App` (or pass it to `SubmitDecision`):

```go
type App struct {
    // ... existing fields ...
    packRegistry *domainpack.PackRegistry
}
```

In `SubmitDecision("discharge")`, before calling `kern.Discharge`:
```go
case "discharge":
    // Look up pack retirement rule
    pack, err := a.resolvePack(req.ScenarioID)
    if err != nil {
        return fmt.Errorf("resolve pack: %w", err)
    }
    rules := getRetirementRules(pack)
    rule, ok := rules[req.ObligationKey]
    if !ok {
        return fmt.Errorf("no retirement rule for debt item: %s", req.ObligationKey)
    }
    // Validate evidence class matches rule
    if req.EvidenceClass != "" && req.EvidenceClass != rule.EvidenceClass {
        return fmt.Errorf("evidence class %q does not match required %q for %s",
            req.EvidenceClass, rule.EvidenceClass, req.ObligationKey)
    }
    return a.kern.Discharge(ctx, req.ScenarioID, req.BeliefID, req.ObligationKey, req.InstrumentRef, req.PrincipalID)
```

Add `EvidenceClass` to the `DecisionRequest` struct.

**Tests:**
- `TestDischargeRequiresMatchingEvidenceClass` — discharge with wrong evidence class → error
- `TestDischargeWithMatchingEvidenceClass` — discharge with correct class → success
- `TestDischargeUnknownDebtItem` — discharge unknown debt item → error

---

## Fix 6: Trust UI authentication

**Files:**
- `internal/ui/ui.go`
- `cmd/argus/main.go`

**Problem:** UI hardcodes principal ID, has no authentication, no CSRF protection.

**Fix:** Static token + Bearer authentication (POC scope — no session management):

1. Add `operatorToken` to `Server`:
```go
type Server struct {
    app           *application.App
    tmpl          *template.Template
    operatorToken string // from ARGUS_OPERATOR_TOKEN env
}
```

2. Add `requireAuth` middleware — checks `Authorization: Bearer <token>` header:
```go
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader != "Bearer "+s.operatorToken {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}
```

3. Protect write endpoints:
```go
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/ui", s.HandleIndex)
    mux.HandleFunc("/ui/insights", s.HandleInsights)
    mux.HandleFunc("/ui/debts", s.HandleDebts)
    mux.HandleFunc("/ui/api/discharge", s.requireAuth(s.HandleDischarge))
    mux.HandleFunc("/ui/api/promote", s.requireAuth(s.HandlePromote))
}
```

4. Server-derived principal from token (never from request body):
```go
func (s *Server) HandleDischarge(w http.ResponseWriter, r *http.Request) {
    // ... decode request ...
    // Principal is always "operator" — derived from authenticated session
    principalID := "operator"
    // ... rest of handler
}
```

5. Wire token from `cmdServe`:
```go
operatorToken := os.Getenv("ARGUS_OPERATOR_TOKEN")
if operatorToken == "" {
    operatorToken = "dev-token-not-for-production"
}
uiServer := ui.NewServer(app, operatorToken)
```

**Tests:**
- `TestDischargeRequiresAuth` — POST without token → 401
- `TestDischargeWithValidToken` — POST with Bearer token → 200
- `TestPromoteRequiresAuth` — POST without token → 401

---

## Fix 7: MCP stdio transport

**Files:**
- `cmd/argus/main.go:102-124`
- `internal/mcp/adapter.go`
- `go.mod` (new dependency)

**Problem:** `argus mcp` prints a message and blocks forever. No MCP protocol framing.

**Fix:** Use the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk/mcp`) for proper MCP protocol compliance. This is the one place where a new dependency is warranted — MCP is an external protocol boundary and the official SDK tracks the current protocol revision.

### Add SDK dependency
```bash
go get github.com/modelcontextprotocol/go-sdk/mcp
```

### Rewire `cmdMCP` to use SDK stdio transport

```go
func cmdMCP(args []string) {
    fs := flag.NewFlagSet("mcp", flag.ExitOnError)
    dbURL := fs.String("db", "postgres://root@localhost:26260/argus?sslmode=disable", "database URL")
    fs.Parse(args)

    db, err := sql.Open("pgx", *dbURL)
    if err != nil { log.Fatalf("open db: %v", err) }
    defer db.Close()
    if err := db.Ping(); err != nil { log.Fatalf("ping db: %v", err) }

    app := application.New(db)
    adapter := mcp.NewAdapter(app)

    // Create MCP server with stdio transport
    server := mcp.NewServer(&mcp.ServerConfig{
        Name:    "argus",
        Version: "0.1.0",
    }, nil)

    // Register tools
    server.AddTool(mcp.Tool{
        Name:        "argus.get_context",
        Description: "Get the RCP/v1 context for a task",
        InputSchema: mcp.ToolInputSchema{...},
    }, adapter.HandleGetContextMCP)

    server.AddTool(mcp.Tool{
        Name:        "argus.submit_packet",
        Description: "Submit an EBP research packet",
        InputSchema: mcp.ToolInputSchema{...},
    }, adapter.handleSubmitPacketMCP)

    // Run over stdio
    if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
        log.Fatalf("mcp: %v", err)
    }
}
```

### Adapt the existing adapter

The existing `internal/mcp/adapter.go` already has `HandleTool` dispatch. Adapt the handler signatures to match the SDK's expected `ToolHandlerFunc` type:

```go
// HandleGetContextMCP adapts the existing handler for the SDK
func (a *Adapter) HandleGetContextMCP(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := req.Params.Arguments.(map[string]interface{})
    taskID, _ := args["task_id"].(string)
    // ... call a.app.GetContext, return mcp.CallToolResult
}
```

The `ListTools` function is no longer needed — tools are registered with the SDK server.

### Architecture boundary
```
OpenCode → stdio/MCP → official MCP SDK → internal/mcp → internal/application
```

No MCP business logic, authority logic, or Solvent knowledge inside the SDK adapter. The adapter remains thin.

**Tests:**
- `TestMCPStdioInitialize` — send initialize request → get capabilities response
- `TestMCPStdioToolsList` — send tools/list → get two tools
- `TestMCPStdioGetContext` — send tools/call argus.get_context → get context
- `TestMCPStdioSubmitPacket` — send tools/call argus.submit_packet → get result
- `TestMCPStdioUnknownMethod` → error response

---

## Fix 8: Semantic version comparison

**File:** `internal/application/app.go:163-166`

**Problem:** `versionGTE` uses lexicographic string comparison. `"0.10.0" < "0.2.0"` lexicographically.

**Fix:** Implement proper semver comparison without adding a dependency:

```go
func versionGTE(version, minVersion string) bool {
    vParts := parseSemver(version)
    mParts := parseSemver(minVersion)
    for i := 0; i < 3; i++ {
        if vParts[i] > mParts[i] { return true }
        if vParts[i] < mParts[i] { return false }
    }
    return true // equal
}

func parseSemver(v string) [3]int {
    var parts [3]int
    v = strings.TrimPrefix(v, "v")
    segs := strings.SplitN(v, ".", 3)
    for i, s := range segs {
        s = strings.Split(s, "-")[0] // strip pre-release
        fmt.Sscanf(s, "%d", &parts[i])
    }
    return parts
}
```

**Tests:**
- `TestVersionGTEEqual` — "0.1.0" >= "0.1.0" → true
- `TestVersionGTEMajorGreater` — "1.0.0" >= "0.9.9" → true
- `TestVersionGTEMajorLess` — "0.9.9" >= "1.0.0" → false
- `TestVersionGTEMinorGreater` — "0.10.0" >= "0.2.0" → true
- `TestVersionGTEMinorLess` — "0.2.0" >= "0.10.0" → false
- `TestVersionGTEPatchGreater` — "0.1.2" >= "0.1.1" → true
- `TestVersionGTEWithPrerelease` — "0.2.0-rc1" >= "0.1.0" → true

---

## Fix 9: Server-verify evidence hashes

**File:** `internal/application/app.go:198-214`

**Problem:** `content_sha256` is client-supplied and never verified. An agent can claim any hash.

**Fix:** For `reproducible_artifact` evidence with an `ArtifactRef`, verify the hash against the artifact registry's stored hash. For other evidence classes, log a warning but accept (agent-attested provenance):

```go
// In Persist, after resolving beliefID for each evidence item:
if e.ArtifactRef != "" && a.artifactReg != nil {
    artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
    if err == nil {
        // Server-computed hash must match agent-declared hash
        if artifact.ArtifactHash != e.ContentSHA256 {
            return nil, fmt.Errorf("evidence[%d]: content hash mismatch: agent says %s, artifact has %s",
                i, e.ContentSHA256, artifact.ArtifactHash)
        }
    }
}
```

For evidence without artifacts, the hash remains agent-attested (this is the expected behavior for `operator_asserted` provenance).

**Tests:**
- `TestEvidenceHashVerified AgainstArtifact` — agent declares wrong hash → error
- `TestEvidenceHashMatchesArtifact` — agent declares correct hash → success
- `TestEvidenceHashAgentAttested` — no artifact, agent-supplied hash → accepted

---

## Fix 10: Port preflight and lifecycle

**Files:**
- `cmd/argus/main.go`
- `Taskfile.yml`

**Problem:** No port preflight, no fallback, no cleanup. `:8080` CRDB admin/ARGUS collision.

**Fix:**

### Port allocation function
```go
// internal/port/preflight.go
package port

func FindAvailable(preferred int) (int, error) {
    ln, err := net.Listen("tcp", fmt.Sprintf(":%d", preferred))
    if err != nil {
        // Try next 10 ports
        for p := preferred + 1; p <= preferred+10; p++ {
            ln2, err2 := net.Listen("tcp", fmt.Sprintf(":%d", p))
            if err2 == nil {
                ln2.Close()
                return p, nil
            }
        }
        return 0, fmt.Errorf("no available port near %d", preferred)
    }
    ln.Close()
    return preferred, nil
}
```

### Preflight in cmdServe
```go
func cmdServe(args []string) {
    // Preflight ports
    crdbSQLPort, _ := port.FindAvailable(26257)
    crdbHTTPPort, _ := port.FindAvailable(8081) // avoid 8080 collision
    argusPort, _ := port.FindAvailable(8080)

    log.Printf("ports: crdb-sql=%d crdb-http=%d argus=%d", crdbSQLPort, crdbHTTPPort, argusPort)
    // ... use these ports in DSN and listen address
}
```

### Taskfile changes
- Add `task down` that kills only processes started by `task dev/fresh`
- Use PID files in `.argus-pids/`
- Add signal propagation (trap SIGINT/SIGTERM)
- Replace `sleep 3` with readiness probe
- Use `ARGUS_CRDB_PORT` env var for test consistency

### Test port configuration
- Define `const TestCRDBPort = 26257` in a shared test helper
- Add `TestPortFallbackBehavior` — occupy preferred port, verify fallback
- Add `TestPortPreflightFindsAvailable` — verify port discovery

---

## Fix 11: task down + PID cleanup

**File:** `Taskfile.yml`

**Add:**
```yaml
  down:
    desc: Stop background CRDB and ARGUS processes
    cmds:
      - |
        if [ -f .argus-pids/crdb.pid ]; then
          kill $(cat .argus-pids/crdb.pid) 2>/dev/null || true
          rm .argus-pids/crdb.pid
        fi
        if [ -f .argus-pids/argus.pid ]; then
          kill $(cat .argus-pids/argus.pid) 2>/dev/null || true
          rm .argus-pids/argus.pid
        fi
```

**Modify `dev` and `fresh`** to record PIDs:
```yaml
  dev:
    cmds:
      - mkdir -p .argus-pids
      - cockroach start-single-node --insecure --listen-addr :26260 --store=.cockroach-data &
      - echo $! > .argus-pids/crdb.pid
      - until cockroach node status --host :26260 --insecure 2>/dev/null; do sleep 1; done
      - go run ./cmd/argus reset --db postgres://root@localhost:26260/argus?sslmode=disable
      - go run ./cmd/argus serve --db postgres://root@localhost:26260/argus?sslmode=disable &
      - echo $! > .argus-pids/argus.pid
```

**Add signal handling** via a wrapper script or `trap` in shell cmds.

---

## Fix 12: Integration test CRDB config

**File:** `integration_test.go`, `internal/application/app_test.go`, `internal/ui/ui_test.go`

**Problem:** Tests default to port 26257. Taskfile starts CRDB on 26260. Mismatch.

**Fix:** Standardize on a single test port via environment:

1. Define a shared constant in a test helper package:
```go
// internal/testutil/testdb.go
package testutil

const TestCRDBPort = 26257

func TestDSN(dbName string) string {
    port := os.Getenv("ARGUS_CRDB_PORT")
    if port == "" {
        port = fmt.Sprintf("%d", TestCRDBPort)
    }
    return fmt.Sprintf("postgres://root@localhost:%s/%s?sslmode=disable", port, dbName)
}
```

2. Update all test files to use `testutil.TestDSN()`.

3. Update Taskfile to start CRDB on the same port:
```yaml
  dev:
    cmds:
      - cockroach start-single-node --insecure --listen-addr :26257 --store=.cockroach-data &
```

**Tests:** No new tests needed — existing tests validate the config.

---

## Fix 13: Tests for every fix

Each fix above includes its own tests. Additionally:

1. **Full suite regression:** `go test ./...`, `go test -race ./...`, `go vet ./...`
2. **Acceptance criteria re-evaluation** against the 25 criteria from `PIVOT_POC_IMPLEMENTATION_PLAN.md`

---

## Acceptance Criteria Impact

| AC | Criterion | Before | After |
|----|-----------|--------|-------|
| 1 | Fresh-agent reconstruction via `argus.get_context` | FAIL | **PASS** (Fix 4 + Fix 7) |
| 5 | Retirement-rule enforcement | FAIL | **PASS** (Fix 5) |
| 6 | Adversarial challenge (contradicts edge) | PARTIAL | **PASS** (Fix 3) |
| 10 | UNKNOWN ≠ EMPTY | FAIL | **PASS** (Fix 4) |
| 11 | Operator identity enforced | FAIL | **PASS** (Fix 6) |
| 12 | Idempotency isolation | PARTIAL | **PASS** (Fix 1) |
| 13 | No UI bypass | FAIL | **PARTIAL** (Fix 6 reduces blast radius) |
| 15 | All tests pass | PASS | **PASS** (Fix 13) |
| 21 | Migration ownership | FAIL | **PASS** (Fix 2) |
| P2-5 | Port mismatch | FAIL | **PASS** (Fix 10 + Fix 12) |

---

## Execution Order

1. Fix 1 (contentHash) — no dependencies
2. Fix 8 (semver) — no dependencies
3. Fix 3 (edge conflicts) — no dependencies
4. Fix 9 (evidence hashes) — no dependencies
5. Fix 2 (Solvent migrations) — needs Solvent package creation
6. Fix 4 (UNKNOWN ≠ EMPTY) — needs view.go changes
7. Fix 5 (retirement rules) — needs registry fix
8. Fix 6 (UI auth) — needs token wiring
9. Fix 7 (MCP stdio) — needs JSON-RPC types
10. Fix 10 (port preflight) — needs port package
11. Fix 11 (task down) — needs Taskfile changes
12. Fix 12 (test CRDB config) — needs shared testutil
13. Fix 13 (full test pass) — after all fixes

---

## Constraints

- **No third Solvent change.** Changes A (migration export) and B (registry decodePack) are the two allowed.
- **No Coordinator REST resurrection.** All fixes work within the single-process architecture.
- **No separate Trust UI process.** Auth is embedded in `internal/ui`.
- **No new architecture layers.** Fixes are within existing packages.
- **Preserve all existing passing tests.** No test deletions.
