# Plan 10r — Corrective Implementation Pass (rev)

**Date:** 2026-09-18 (revision)
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
- `TestContentHashDifferentEvidenceProducesDifferentHash`
- `TestContentHashDifferentEdgesProducesDifferentHash`
- `TestContentHashCanonicalOrdering`

---

## Fix 2: Solvent migration export

**Problem:** `argus reset` has `// TODO: call solventmigrations.Apply(ctx, db)`. The `solventmigrations` package does not exist. Solvent's `internal/` boundary prevents Oracle from importing internal packages.

**Fix (Solvent Change #1 of 2):** Create a canonical migration package in Solvent.

### Step A: Move canonical SQL to `solvent-main/migrations/db/`

Move (not copy) the existing `solvent-main/db/*.sql` files to `solvent-main/migrations/db/`. Update any internal references (test harness, m0) to use the new path. One authoritative location.

### Step B: Create `solvent-main/migrations/migrations.go`

```go
package migrations

import (
    "context"
    "database/sql"
    "embed"
    "fmt"
    "sort"
    "strings"
)

//go:embed db/*.sql
var FS embed.FS

func Apply(ctx context.Context, db *sql.DB) error {
    entries, _ := FS.ReadDir("db")
    var sqlFiles []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
            sqlFiles = append(sqlFiles, e.Name())
        }
    }
    sort.Strings(sqlFiles)
    for _, name := range sqlFiles {
        data, _ := FS.ReadFile("db/" + name)
        for _, stmt := range splitStatements(string(data)) {
            if _, err := db.ExecContext(ctx, stmt); err != nil {
                return fmt.Errorf("apply %s: %w", name, err)
            }
        }
    }
    return nil
}

func splitStatements(sql string) []string { /* existing logic */ }
```

**Verification:** This is Change #1 of 2 in the Solvent change budget. No new tables, columns, or constraints. Packaging change only.

### Step C: Wire into Oracle's `cmdReset`

**File:** `cmd/argus/main.go:182-183`

```go
import solventmigrations "github.com/PithomLabs/solvent/migrations"

// in cmdReset:
log.Println("argus reset: applying Solvent migrations...")
if err := solventmigrations.Apply(ctx, db); err != nil {
    log.Fatalf("solvent migrations: %v", err)
}
```

The `go.mod` replace directive already points to `../solvent-main`.

**Tests:**
- `TestCmdResetAppliesSolventMigrations` — run against empty CRDB, verify Solvent tables exist
- `TestSolventMigrationsAreIdempotent` — apply twice, no errors

---

## Fix 3: Edge kind conflict detection

**File:** `internal/application/app.go:257-264`

**Problem:** `ON CONFLICT (parent_id, child_id) DO NOTHING` silently drops a `contradicts` edge when a `derives` edge already exists for the same pair.

**Fix:** Replace `DO NOTHING` with explicit conflict detection:

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

**Tests:**
- `TestEdgeKindConflictReturnsError` — insert derives, then contradicts on same pair → error
- `TestEdgeSameKindIsIdempotent` — insert derives twice → no error
- `TestEdgeDifferentPairsSucceed` — different parent/child pairs → both stored

---

## Fix 4: UNKNOWN ≠ EMPTY in GetContext

**Files:**
- `internal/application/app.go:298-322` — `GetContext`
- `internal/epistemic/view.go` — `Snapshot`, `Context` types
- `internal/mcp/adapter.go` — `handleGetContext`

**Problem:** `GetContext` returns `nil, error` when a backend is unavailable. The MCP adapter propagates this as a tool-call error. An agent cannot distinguish "no research exists" from "backend is down."

**Fix:** Add availability metadata to `Context`:

1. Add to `internal/epistemic/view.go`:
```go
type Availability struct {
    WorkAvailable      bool   `json:"work_available"`
    EpistemicAvailable bool   `json:"epistemic_available"`
    Reason             string `json:"reason,omitempty"`
}
```

2. Add to `Context` struct:
```go
type Context struct {
    Task         *work.Task          `json:"task"`
    Dependencies []*work.Dependency  `json:"dependencies"`
    Snapshot     *epistemic.Snapshot `json:"snapshot"`
    Availability Availability       `json:"availability"`
}
```

3. In `GetContext`, catch errors per sub-call:
```go
func (a *App) GetContext(ctx context.Context, taskID string) (*Context, error) {
    result := &Context{
        Availability: Availability{WorkAvailable: true, EpistemicAvailable: true},
        Snapshot:     &epistemic.Snapshot{},
    }

    task, err := a.workStore.GetByID(ctx, taskID)
    if err != nil {
        result.Availability.WorkAvailable = false
        result.Availability.Reason = "work_unavailable"
    } else {
        result.Task = task
    }

    // ... similar for dependencies and snapshot
    // If snapshot query fails:
    //   result.Availability.EpistemicAvailable = false
    //   result.Availability.Reason = "solvent_unavailable"

    return result, nil // never return error for availability issues
}
```

**Tests:**
- `TestGetContextWithUnavailableBackend` — mock DB failure → returns Context with `epistemic_available=false`
- `TestGetContextWithEmptyScenario` — valid DB, no beliefs → returns Context with `epistemic_available=true` and empty snapshot
- `TestGetContextNormal` — valid DB, beliefs exist → returns full Context

---

## Fix 5: Retirement-rule enforcement before Discharge (BLOCKER FIX)

**Files:**
- `internal/application/app.go` — `SubmitDecision`, `App` struct, `DecisionRequest`
- `domain-pack/registry.go` — `decodePack` fix
- `internal/application/app_test.go` — new tests

**Problem:** `SubmitDecision("discharge")` calls `kern.Discharge()` directly without:
1. Checking that the debt item exists in the Pack's retirement rules
2. Checking that the required evidence class matches the rule
3. Verifying that qualifying persisted evidence actually exists for the belief/scenario

The current check `req.EvidenceClass != "" && req.EvidenceClass != rule.EvidenceClass` passes when `EvidenceClass` is empty.

### Step A: Fix `decodePack` to return typed packs

**File:** `domain-pack/registry.go:170-193`

The `decodePack` function always returns `genericPack` which discards retirement rules, debt vocabulary, and verifier specs. Fix to detect BM-IST packs and return `*bmistv1.Pack`:

```go
func decodePack(data []byte, expectedPackID string) (Pack, error) {
    var probe struct {
        PackID  string `json:"pack_id"`
        Version string `json:"version"`
    }
    if err := json.Unmarshal(data, &probe); err != nil {
        return nil, err
    }
    if strings.EqualFold(probe.PackID, "bmist") {
        return bmistv1.ParsePack(data)
    }
    return &genericPack{PackID: probe.PackID, Version: probe.Version}, nil
}
```

Add a `retirementRulesProvider` interface to `domain-pack/registry.go`:
```go
type retirementRulesProvider interface {
    GetRetirementRules() map[string]RetirementRule
    GetDebtVocabulary() []string
}
```

### Step B: Give App access to pack registry

Add to `App` struct:
```go
type App struct {
    // ... existing fields ...
    packRegistry *domainpack.PackRegistry
}
```

Add a constructor option or setter:
```go
func (a *App) SetPackRegistry(r *domainpack.PackRegistry) {
    a.packRegistry = r
}
```

### Step C: Add EvidenceClass to DecisionRequest

```go
type DecisionRequest struct {
    Type          string `json:"type"`
    ScenarioID    string `json:"scenario_id"`
    BeliefID      string `json:"belief_id"`
    ObligationKey string `json:"obligation_key,omitempty"`
    InstrumentRef string `json:"instrument_ref,omitempty"`
    PrincipalID   string `json:"principal_id"`
    EvidenceClass string `json:"evidence_class,omitempty"` // NEW
}
```

### Step D: Full retirement-rule enforcement in SubmitDecision

Replace the `discharge` case with the complete chain:

```go
case "discharge":
    // 1. Pack registry must exist (fail-closed)
    if a.packRegistry == nil {
        return fmt.Errorf("pack registry unavailable: cannot validate retirement")
    }

    // 2. Resolve pack for scenario
    pack, err := a.resolvePack(req.ScenarioID)
    if err != nil {
        return fmt.Errorf("pack resolution failed: %w", err)
    }

    // 3. Debt item must exist in pack vocabulary
    rulesProv, ok := pack.(interface{ GetRetirementRules() map[string]interface{} })
    if !ok {
        return fmt.Errorf("pack does not support retirement rules")
    }
    rules := rulesProv.GetRetirementRules()
    ruleRaw, ok := rules[req.ObligationKey]
    if !ok {
        return fmt.Errorf("no retirement rule for debt item: %s", req.ObligationKey)
    }

    // 4. Extract typed rule
    rule, ok := ruleRaw.(interface{ GetEvidenceClass() string })
    if !ok {
        return fmt.Errorf("invalid retirement rule type for: %s", req.ObligationKey)
    }
    requiredClass := rule.GetEvidenceClass()

    // 5. EvidenceClass is REQUIRED — empty is not allowed
    if req.EvidenceClass == "" {
        return fmt.Errorf("evidence_class is required for debt discharge of: %s", req.ObligationKey)
    }
    if req.EvidenceClass != requiredClass {
        return fmt.Errorf("evidence class mismatch: debt %q requires %q, got %q",
            req.ObligationKey, requiredClass, req.EvidenceClass)
    }

    // 6. Verify qualifying persisted evidence exists
    evidenceIDs, err := a.verifyPersistedEvidence(ctx, req.BeliefID, req.ScenarioID, req.EvidenceClass)
    if err != nil {
        return fmt.Errorf("evidence verification failed: %w", err)
    }
    if len(evidenceIDs) == 0 {
        return fmt.Errorf("no qualifying evidence of class %q for belief %s in scenario %s",
            req.EvidenceClass, req.BeliefID, req.ScenarioID)
    }

    // 7. Build instrument ref from verified evidence
    instrumentRef := buildInstrumentRef(evidenceIDs)

    return a.kern.Discharge(ctx, req.ScenarioID, req.BeliefID, req.ObligationKey, instrumentRef, req.PrincipalID)
```

### Step E: Query persisted evidence

Add a helper that queries the evidence table directly (single-process architecture, same DB):

```go
func (a *App) verifyPersistedEvidence(ctx context.Context, beliefID, scenarioID, evidenceClass string) ([]string, error) {
    rows, err := a.db.QueryContext(ctx,
        `SELECT id FROM evidence
         WHERE belief_id = $1::UUID AND scenario_id = $2::UUID AND provenance_class = $3`,
        beliefID, scenarioID, evidenceClass)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }
    return ids, rows.Err()
}

func buildInstrumentRef(evidenceIDs []string) string {
    if len(evidenceIDs) == 0 { return "" }
    ref := "evidence:"
    for i, id := range evidenceIDs {
        if i > 0 { ref += "," }
        ref += id
    }
    return ref
}
```

### Step F: Wire pack registry into App construction

**File:** `cmd/argus/main.go` — `cmdServe` and `cmdMCP`

```go
// In cmdServe:
packReg := domainpack.NewRegistry()
if err := packReg.LoadFromDisk("domain-pack"); err != nil {
    log.Printf("warning: pack load failed: %v", err)
}
app := application.NewWithRegistry(db, registry, specs)
app.SetPackRegistry(packReg)
```

### Step G: Update Trust UI to pass EvidenceClass

**File:** `internal/ui/ui.go` — `HandleDischarge`

Add `EvidenceClass` to the request struct and pass it through:

```go
var req struct {
    ScenarioID    string `json:"scenario_id"`
    BeliefID      string `json:"belief_id"`
    ObligationKey string `json:"obligation_key"`
    EvidenceClass string `json:"evidence_class"` // NEW
    InstrumentRef string `json:"instrument_ref"`
}
```

**Tests:**
- `TestDischargeRequiresKnownDebtItem` — discharge unknown debt item → error
- `TestDischargeRequiresEvidenceClass` — discharge with empty evidence class → error
- `TestDischargeWithMismatchedEvidenceClass` — discharge with wrong class → error
- `TestDischargeRequiresPersistedEvidence` — discharge with correct class but no evidence in DB → error
- `TestDischargeWithQualifyingEvidence` — all checks pass → success
- `TestDischargeWithRegistryUnavailable` — no pack registry → error (fail-closed)

---

## Fix 6: Trust UI authentication

**Files:**
- `internal/ui/ui.go`
- `cmd/argus/main.go`

**Problem:** UI hardcodes principal ID, has no authentication, no CSRF protection.

**Fix:** Static token + Bearer authentication with browser-compatible flow:

### Design

```
ARGUS_OPERATOR_TOKEN env var
        ↓
Browser: user enters token in login form (or via dev tools)
        ↓
POST /ui/api/login { token: "..." }
        ↓
Server validates token, sets HttpOnly cookie
        ↓
Subsequent requests carry cookie automatically
        ↓
Write endpoints check cookie (or Bearer header for API clients)
        ↓
Server-derived principal, never from request body
```

### Implementation

1. Add to `Server`:
```go
type Server struct {
    app           *application.App
    tmpl          *template.Template
    operatorToken string
}
```

2. Add login endpoint (sets HttpOnly cookie):
```go
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req struct { Token string `json:"token"` }
    json.NewDecoder(r.Body).Decode(&req)
    if req.Token != s.operatorToken {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
    }
    http.SetCookie(w, &http.Cookie{
        Name:     "argus_session",
        Value:    s.operatorToken,
        Path:     "/",
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    })
    w.WriteHeader(http.StatusOK)
}
```

3. Add `requireAuth` — checks cookie OR Bearer header:
```go
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check cookie first (browser path)
        cookie, err := r.Cookie("argus_session")
        if err == nil && cookie.Value == s.operatorToken {
            next(w, r)
            return
        }
        // Fallback: Bearer header (API/CLI path)
        authHeader := r.Header.Get("Authorization")
        if authHeader == "Bearer "+s.operatorToken {
            next(w, r)
            return
        }
        http.Error(w, "unauthorized", http.StatusUnauthorized)
    }
}
```

4. Add Origin/CSRF check for write endpoints:
```go
func (s *Server) requireOrigin(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin != "" && origin != "http://localhost:"+getPort() {
            http.Error(w, "forbidden: invalid origin", http.StatusForbidden)
            return
        }
        next(w, r)
    }
}
```

5. Protect write endpoints:
```go
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/ui", s.HandleIndex)
    mux.HandleFunc("/ui/insights", s.HandleInsights)
    mux.HandleFunc("/ui/debts", s.HandleDebts)
    mux.HandleFunc("/ui/api/login", s.HandleLogin)
    mux.HandleFunc("/ui/api/discharge", s.requireAuth(s.requireOrigin(s.HandleDischarge)))
    mux.HandleFunc("/ui/api/promote", s.requireAuth(s.requireOrigin(s.HandlePromote)))
}
```

6. Server-derived principal from token validation (never from request body):
```go
func (s *Server) HandleDischarge(w http.ResponseWriter, r *http.Request) {
    // ... decode request ...
    // Principal is always "operator" — derived from authenticated session
    principalID := "operator"
    // ... rest of handler
}
```

7. Wire token from `cmdServe`:
```go
operatorToken := os.Getenv("ARGUS_OPERATOR_TOKEN")
if operatorToken == "" {
    operatorToken = "dev-token-not-for-production"
}
uiServer := ui.NewServer(app, operatorToken)
```

**Tests:**
- `TestDischargeRequiresAuth` — POST without token/cookie → 401
- `TestDischargeWithValidToken` — POST with Bearer token → 200
- `TestDischargeWithValidCookie` — POST with session cookie → 200
- `TestDischargeRejectsBodyPrincipal` — body principal_id is ignored
- `TestPromoteRequiresAuth` — POST without token → 401

---

## Fix 7: MCP stdio transport (official SDK)

**Files:**
- `cmd/argus/main.go:102-124`
- `internal/mcp/adapter.go`
- `go.mod` (new dependency)

**Problem:** `argus mcp` prints a message and blocks forever. No MCP protocol framing.

**Fix:** Use the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk/mcp`).

### Verified SDK API (from pkg.go.dev, v1.8.0)

```go
// Server creation
server := mcp.NewServer(&mcp.Implementation{Name: "argus", Version: "v0.1.0"}, nil)

// Tool registration — top-level generic function auto-populates schemas
mcp.AddTool(server, &mcp.Tool{
    Name:        "argus.get_context",
    Description: "Get the RCP/v1 context for a task",
}, handlerFunc)

// Raw schema alternative (for existing adapter)
server.AddTool(mcp.Tool{
    Name:        "argus.submit_packet",
    Description: "Submit an EBP research packet",
    InputSchema: map[string]interface{}{...},
}, rawHandler)

// Stdio transport
server.Run(ctx, &mcp.StdioTransport{})
```

### Implementation

```go
func cmdMCP(args []string) {
    fs := flag.NewFlagSet("mcp", flag.ExitOnError)
    dbURL := fs.String("db", "postgres://root@localhost:26260/argus?sslmode=disable", "database URL")
    fs.Parse(args)

    db, err := sql.Open("pgx", *dbURL)
    if err != nil { log.Fatalf("open db: %v", err) }
    defer db.Close()
    if err := db.Ping(); err != nil { log.Fatalf("ping db: %v", err) }

    app := application.NewWithRegistry(db, verifier.NewArtifactRegistry(), []application.VerifierSpec{
        {VerifierID: "physics-v1", MinVersion: "0.1.0"},
    })

    server := mcp.NewServer(&mcp.Implementation{
        Name:    "argus",
        Version: "v0.1.0",
    }, nil)

    // Register exactly two tools
    adapter := mcpinternal.NewAdapter(app)
    adapter.RegisterTools(server)

    // Run over stdio
    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatalf("mcp: %v", err)
    }
}
```

### Adapter registration

```go
// internal/mcp/adapter.go
func (a *Adapter) RegisterTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "argus.get_context",
        Description: "Get the RCP/v1 context for a task: task, dependencies, epistemic state, and activity.",
    }, a.handleGetContextMCP)

    server.AddTool(mcp.Tool{
        Name:        "argus.submit_packet",
        Description: "Submit an EBP research packet. Validates, persists beliefs/evidence/edges/tasks. Idempotent.",
        InputSchema: map[string]interface{}{
            "type": "object",
            "properties": map[string]interface{}{
                "schema_version": map[string]interface{}{"type": "string"},
                "role":           map[string]interface{}{"type": "string", "enum": []string{"work", "adversarial"}},
                "packet_id":      map[string]interface{}{"type": "string"},
                "pack_ref":       map[string]interface{}{"type": "string"},
                "scenario_id":    map[string]interface{}{"type": "string"},
                "beliefs":        map[string]interface{}{"type": "array"},
                "evidence":       map[string]interface{}{"type": "array"},
                "edges":          map[string]interface{}{"type": "array"},
                "tasks":          map[string]interface{}{"type": "array"},
            },
            "required": []string{"schema_version", "role", "packet_id", "pack_ref", "beliefs", "evidence"},
        },
    }, a.handleSubmitPacketRaw)
}
```

### Architecture boundary
```
OpenCode → stdio/MCP → official MCP SDK → internal/mcp → internal/application
```

No MCP business logic, authority logic, or Solvent knowledge inside the SDK adapter. The adapter remains thin.

**Note:** The `handleSubmitPacket` handler needs to accept `json.RawMessage` args (raw schema mode) since the packet structure is complex. Use `server.AddTool` (non-generic) for this tool.

**Tests:**
- `TestMCPToolsRegistered` — verify server has exactly 2 tools
- `TestMCPGetContextViaSDK` — in-memory transport test
- `TestMCPSubmitPacketViaSDK` — in-memory transport test

---

## Fix 8: Semantic version comparison

**File:** `internal/application/app.go:163-166`

**Problem:** `versionGTE` uses lexicographic string comparison. `"0.10.0" < "0.2.0"` lexicographically.

**Fix:** Implement proper semver comparison. Reject prerelease versions (constrain VerifierSpec to release versions):

```go
func versionGTE(version, minVersion string) bool {
    v, err1 := parseSemver(version)
    m, err2 := parseSemver(minVersion)
    if err1 != nil || err2 != nil {
        return false // malformed versions are never >= anything
    }
    if v[0] != m[0] { return v[0] > m[0] }
    if v[1] != m[1] { return v[1] > m[1] }
    return v[2] >= m[2]
}

func parseSemver(v string) ([3]int, error) {
    var parts [3]int
    v = strings.TrimPrefix(v, "v")
    // Reject prerelease — only release versions allowed
    if strings.Contains(v, "-") {
        return parts, fmt.Errorf("prerelease versions not supported: %s", v)
    }
    segs := strings.SplitN(v, ".", 3)
    if len(segs) != 3 {
        return parts, fmt.Errorf("invalid semver: %s", v)
    }
    for i, s := range segs {
        if _, err := fmt.Sscanf(s, "%d", &parts[i]); err != nil {
            return parts, fmt.Errorf("invalid semver component: %s", err)
        }
    }
    return parts, nil
}
```

**Tests:**
- `TestVersionGTEEqual` — "0.1.0" >= "0.1.0" → true
- `TestVersionGTEMajorGreater` — "1.0.0" >= "0.9.9" → true
- `TestVersionGTEMinorGreater` — "0.10.0" >= "0.2.0" → true
- `TestVersionGTEMinorLess` — "0.2.0" >= "0.10.0" → false
- `TestVersionGTEPatchGreater` — "0.1.2" >= "0.1.1" → true
- `TestVersionGTERejectsPrerelease` — "0.2.0-rc1" → error

---

## Fix 9: Server-verify evidence hashes

**File:** `internal/application/app.go:198-214`

**Problem:** `content_sha256` is client-supplied and never verified.

**Fix:** For `reproducible_artifact` evidence with an `ArtifactRef`, verify the hash against the artifact registry's stored hash. For other evidence classes, accept the agent-declared hash (agent-attested provenance is expected for `operator_asserted`):

```go
// In Persist, after resolving beliefID for each evidence item:
if e.ArtifactRef != "" && a.artifactReg != nil {
    artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
    if err == nil {
        if artifact.ArtifactHash != e.ContentSHA256 {
            return nil, fmt.Errorf("evidence[%d]: content hash mismatch: agent says %s, artifact has %s",
                i, e.ContentSHA256, artifact.ArtifactHash)
        }
    }
}
```

**Constraint:** Only trusted verifier execution (via `verifier.RunPhysicsVerifier`) can produce qualifying `reproducible_artifact` evidence. Agent-declared verifier/hash fields are never sufficient by themselves — the artifact must be registered via the trusted path.

**Tests:**
- `TestEvidenceHashVerifiedAgainstArtifact` — agent declares wrong hash → error
- `TestEvidenceHashMatchesArtifact` — agent declares correct hash → success
- `TestEvidenceHashAgentAttested` — no artifact, agent-supplied hash → accepted

---

## Fix 10: Port preflight and lifecycle (BLOCKER FIX)

**Files:**
- `cmd/argus/main.go`
- `Taskfile.yml`
- Test files

**Problem:** No port preflight, no fallback, no cleanup. `:8080` CRDB admin/ARGUS collision. Probe-then-close race condition.

**Design principle:** Single runtime-resolved topology. All ports resolved before any process starts.

### Architecture

```
port preflight
    ↓
resolve CRDB SQL port (default 26257, fallback 26258-26267)
    ↓
resolve CRDB admin port (default 8081, fallback 8082-8091)
    ↓
resolve ARGUS HTTP port (default 8080, fallback 8080-8090)
    ↓
start CRDB with resolved ports
    ↓
write runtime config (.argus-runtime.env)
    ↓
ARGUS + tests consume resolved ports
    ↓
task down / Ctrl-C → kill only tracked PIDs
```

### Port package

Create `internal/port/preflight.go`:

```go
package port

import (
    "fmt"
    "net"
    "time"
)

// Reserve holds a listener until the actual service starts.
type Reserve struct {
    ln     net.Listener
    Port   int
    addr   string
}

// Reserve attempts to bind a port. If preferred is occupied, tries next 10.
// Returns a Reserve that must be Released or the port stays bound.
func Reserve(preferred int) (*Reserve, error) {
    for p := preferred; p <= preferred+10; p++ {
        addr := fmt.Sprintf(":%d", p)
        ln, err := net.ListenTimeout("tcp", addr, 1*time.Second)
        if err == nil {
            return &Reserve{ln: ln, Port: p, addr: addr}, nil
        }
    }
    return nil, fmt.Errorf("no available port near %d", preferred)
}

// Release closes the listener, freeing the port for the actual service.
func (r *Reserve) Release() error {
    return r.ln.Close()
}
```

**Key:** The `Reserve` holds the listener open until `Release()` is called just before the service binds. This eliminates the probe-then-close race. The actual startup must happen immediately after release.

### Preflight in cmdServe

```go
func cmdServe(args []string) {
    // 1. Preflight all ports
    crdbSQLReserve, err := port.Reserve(26257)
    if err != nil { log.Fatalf("crdb sql port: %v", err) }
    defer crdbSQLReserve.Release()

    crdbAdminReserve, err := port.Reserve(8081)
    if err != nil { log.Fatalf("crdb admin port: %v", err) }
    defer crdbAdminReserve.Release()

    argusReserve, err := port.Reserve(8080)
    if err != nil { log.Fatalf("argus port: %v", err) }
    defer argusReserve.Release()

    // 2. Print resolved topology
    log.Printf("resolved ports: crdb-sql=%d crdb-admin=%d argus=%d",
        crdbSQLReserve.Port, crdbAdminReserve.Port, argusReserve.Port)

    // 3. Release ports just before starting services
    // ... release and start CRDB ...
    // ... release and start ARGUS ...
}
```

### Taskfile changes

Replace hardcoded ports with runtime-resolved values:

```yaml
  dev:
    cmds:
      - mkdir -p .argus-pids
      - |
        # Resolve ports
        source <(go run ./cmd/argus resolve-ports 2>/dev/null)
        echo "ports: crdb-sql=$CRDB_SQL_PORT crdb-admin=$CRDB_ADMIN_PORT argus=$ARGUS_PORT"
        cockroach start-single-node --insecure \
          --listen-addr :$CRDB_SQL_PORT \
          --http-addr :$CRDB_ADMIN_PORT \
          --store=.cockroach-data &
        echo $! > .argus-pids/crdb.pid
        # Wait for SQL readiness
        until cockroach node status --host :$CRDB_SQL_PORT --insecure 2>/dev/null; do
          sleep 1
        done
        go run ./cmd/argus reset --db postgres://root@localhost:$CRDB_SQL_PORT/argus?sslmode=disable
        go run ./cmd/argus serve \
          --db postgres://root@localhost:$CRDB_SQL_PORT/argus?sslmode=disable \
          --listen :$ARGUS_PORT &
        echo $! > .argus-pids/argus.pid
```

### Port allocation tests

- `TestPortReservePreferredAvailable` — preferred port free → gets preferred
- `TestPortReserveFallback` — preferred occupied → gets next available
- `TestPortReserveExhaustion` — all ports occupied → error
- `TestPortReserveRace` — two concurrent reserves → no conflict

---

## Fix 11: task down + PID cleanup

**File:** `Taskfile.yml`

**Add:**
```yaml
  down:
    desc: Stop background CRDB and ARGUS processes started by this invocation
    cmds:
      - |
        for pidfile in .argus-pids/*.pid; do
          [ -f "$pidfile" ] || continue
          pid=$(cat "$pidfile")
          kill "$pid" 2>/dev/null || true
          rm "$pidfile"
        done
```

**Modify `dev` and `fresh`** to record PIDs in `.argus-pids/`.

**Add signal handling** via a wrapper or `trap`:
```yaml
  dev:
    cmds:
      - mkdir -p .argus-pids
      - |
        trap 'for f in .argus-pids/*.pid; do [ -f "$f" ] && kill $(cat "$f") 2>/dev/null; rm -f "$f"; done' EXIT INT TERM
        # ... existing startup ...
```

---

## Fix 12: Integration test CRDB config

**Problem:** Tests default to port 26257. Taskfile starts CRDB on 26260. Port mismatch.

**Fix:** After Fix 10 (dynamic ports), integration tests consume the runtime-resolved port via `ARGUS_CRDB_PORT` env var. The test helper becomes:

```go
func testDB(t *testing.T) *sql.DB {
    t.Helper()
    dsn := os.Getenv("ARGUS_TEST_DSN")
    if dsn == "" {
        port := os.Getenv("ARGUS_CRDB_PORT")
        if port == "" { port = "26257" }
        dsn = fmt.Sprintf("postgres://root@localhost:%s/defaultdb?sslmode=disable", port)
    }
    // ... rest of test DB setup
}
```

The `reference-loop` test files that hardcode port 26257 without env override also need updating.

**Tests:** No new tests — existing tests validate the config.

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
7. Fix 5 (retirement rules) — needs registry fix + evidence query
8. Fix 6 (UI auth) — needs token wiring
9. Fix 7 (MCP SDK) — needs `go get` + adapter rewrite
10. Fix 10 (port preflight) — needs port package
11. Fix 11 (task down) — needs Taskfile changes
12. Fix 12 (test CRDB config) — after Fix 10
13. Fix 13 (full test pass) — after all fixes

---

## Constraints

- **No third Solvent change.** Changes A (migration export) and B (decodePack) are the two allowed.
- **No Coordinator REST resurrection.** All fixes work within the single-process architecture.
- **No separate Trust UI process.** Auth is embedded in `internal/ui`.
- **No new architecture layers.** Fixes are within existing packages.
- **Preserve all existing passing tests.** No test deletions.
