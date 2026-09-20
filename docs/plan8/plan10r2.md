# Plan 10r2 — Corrective Implementation Pass (rev2)

**Date:** 2026-09-18 (revision 2)
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
| 5 | Retirement-rule enforcement before Discharge | Medium | 3 (app.go, registry, bmist types) + tests |
| 6 | Trust UI authentication | Medium | 2 (ui.go, main.go) + tests |
| 7 | MCP stdio transport (official SDK) | High | 3 (main.go, adapter.go, go.mod) + tests |
| 8 | Semantic version comparison | Low | 1 (app.go) + tests |
| 9 | Server-verify evidence hashes | Low | 1 (app.go) + tests |
| 10 | Port preflight and lifecycle | High | 3 (main.go, Taskfile, tests) |
| 11 | task down + PID cleanup | Medium | 1 (Taskfile) |
| 12 | Integration test CRDB config | Low | 1 (integration_test.go) |
| 13 | Tests for every fix | Low | multiple test files |

---

## Semantic clarification: retirement-rule enforcement

The POC mechanically enforces the Pack's declared evidence-class retirement gate. It does **not** interpret the human-readable `Rule` string. The enforcement chain is:

1. Debt item exists in Pack's `retirement_rules` map
2. The rule's `EvidenceClass` matches the caller's declared evidence class
3. Qualifying persisted evidence exists for the belief/scenario with that `provenance_class`

The `Rule` field is metadata for human review — the POC does not parse or evaluate it.

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

### Step A: Create `solvent-main/migrations/db/` with SQL files

`go:embed` cannot embed files from `../db/` — the path must stay within the package directory. Therefore, create `solvent-main/migrations/db/` and copy the 10 SQL files there. This becomes the single canonical location for the migration package.

The existing `solvent-main/db/` directory is referenced by internal test harness (`internal/testdb`), m0, and `cmd/solvent`. Update those references to point to `migrations/db/` instead. The old `db/` directory is then deleted.

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

func splitStatements(sqlText string) []string {
    var cleaned strings.Builder
    for _, line := range strings.Split(sqlText, "\n") {
        if idx := strings.Index(line, "--"); idx >= 0 {
            line = line[:idx]
        }
        cleaned.WriteString(line)
        cleaned.WriteString("\n")
    }
    var out []string
    for _, part := range strings.Split(cleaned.String(), ";") {
        if s := strings.TrimSpace(part); s != "" {
            out = append(out, s)
        }
    }
    return out
}
```

### Step B.1: Update internal Solvent references

Update the following to use `migrations/db/` instead of `db/`:
- `internal/testdb/testdb.go` — `ApplySchema` callers
- `internal/m0/schema.go` — DDL applier references
- `cmd/solvent/main.go` — `resolveSchemaPaths()`
- All test suites that hardcode schema paths (kernel, authority, etc.)
- `Taskfile.yml` — `db:reset` task
- `demo/cloud/init/main.go` — hardcoded paths

### Step C: Verify idempotency of existing SQL

Before creating the package, verify each file's statements are idempotent:

| File | Statement Type | Idempotent? | Fix needed? |
|------|---------------|-------------|-------------|
| `001_schema.sql` | `CREATE TABLE belief` | **NO** — bare `CREATE TABLE` | Add `IF NOT EXISTS` |
| `001_schema.sql` | `CREATE TABLE belief_edge` | **NO** — bare `CREATE TABLE` | Add `IF NOT EXISTS` |
| `001_schema.sql` | `CREATE INDEX belief_edge_child` | **NO** — bare `CREATE INDEX` | Add `IF NOT EXISTS` |
| `001_schema.sql` | `CREATE INDEX live_intents` | **NO** — bare `CREATE INDEX` | Add `IF NOT EXISTS` |
| `002_corpus.sql` | All statements | YES — uses `IF NOT EXISTS` | No |
| `003_wizard.sql` | All statements | YES — uses `IF NOT EXISTS` | No |
| `004_debt_vocabulary.sql` | `ALTER TABLE ... SET DEFAULT` | YES — idempotent | No |
| `005_authority_mvp.sql` | All statements | YES — uses `IF NOT EXISTS` | No |
| `006_authority_justification_cascade.sql` | `DROP CONSTRAINT IF EXISTS` + re-add | YES | No |
| `007_service_tables.sql` | `CREATE TABLE workflow_token` | **NO** — bare `CREATE TABLE` | Add `IF NOT EXISTS` |
| `007_service_tables.sql` | `CREATE TABLE policy_tool` | **NO** | Add `IF NOT EXISTS` |
| `007_service_tables.sql` | `CREATE TABLE policy_actor` | **NO** | Add `IF NOT EXISTS` |
| `007_service_tables.sql` | `CREATE TABLE audit_activity` | **NO** | Add `IF NOT EXISTS` |
| `007_service_tables.sql` | 4 bare `CREATE INDEX` | **NO** | Add `IF NOT EXISTS` |
| `008_executing_state.sql` | `DROP CONSTRAINT IF EXISTS` + re-add + index | YES | No |
| `009_exact_authority_binding.sql` | `ADD COLUMN IF NOT EXISTS` + index | YES | No |
| `010_debt_opaque.sql` | `ALTER TABLE ... SET DEFAULT` | YES | No |

**Action:** Add `IF NOT EXISTS` to all bare `CREATE TABLE` and `CREATE INDEX` statements in `001_schema.sql` and `007_service_tables.sql`. This makes `Apply()` genuinely idempotent — running twice produces no errors. No migration version table needed.

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

### Step D: Replace test helpers with package

After the migrations package exists, replace the filesystem-based test helpers:
- `integration_test.go:integApplySolventMigrations` → `solventmigrations.Apply(ctx, db)`
- `internal/ui/ui_test.go` — `solvent-main/db` reference → `solventmigrations.Apply(ctx, db)`
- `internal/application/app_test.go:testDB` — `solvent-main/db` reference → `solventmigrations.Apply(ctx, db)`

These currently read from `../solvent-main/db/` at runtime. After the package, they use the embedded FS.

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

2. Add to `Context` struct (in `internal/application/app.go`):
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
- `domain-pack/bmist/v1/types.go` — `GetRetirementRules()` return type fix
- `coordinator/validate.go` — `getRetirementRule` type assertions
- `internal/application/app_test.go` — new tests

**Problem:** `SubmitDecision("discharge")` calls `kern.Discharge()` directly without:
1. Checking that the debt item exists in the Pack's retirement rules
2. Checking that the required evidence class matches the rule
3. Verifying that qualifying persisted evidence actually exists for the belief/scenario

Additionally, the `GetRetirementRules()` method returns `map[string]interface{}` which forces double type-assertions downstream. The concrete field is `map[string]RetirementRule` but the method erases the type.

### Step A: Fix `GetRetirementRules()` to return typed map

**File:** `domain-pack/bmist/v1/types.go:41-47`

Change from:
```go
func (p *Pack) GetRetirementRules() map[string]interface{} {
    result := make(map[string]interface{})
    for k, v := range p.RetirementRules {
        result[k] = v
    }
    return result
}
```

To:
```go
func (p *Pack) GetRetirementRules() map[string]RetirementRule {
    return p.RetirementRules
}
```

This is a direct accessor — no copy needed since `RetirementRules` is already `map[string]RetirementRule`.

### Step B: Add `RetirementRulesProvider` interface to `domain-pack/registry.go`

```go
// RetirementRulesProvider is implemented by Pack types that declare retirement rules.
type RetirementRulesProvider interface {
    GetRetirementRules() map[string]RetirementRule
}

// RetirementRule mirrors the domain-specific retirement rule type.
// This is the canonical type used by the application layer.
type RetirementRule struct {
    EvidenceClass string
    Rule          string
}
```

Wait — `RetirementRule` is defined in `domain-pack/bmist/v1/types.go`. The application layer cannot import `bmistv1` directly (that would couple application to a specific domain). Instead, define a minimal interface in `domain-pack/`:

```go
// RetirementRule describes the evidence-class gate for retiring a debt item.
type RetirementRule interface {
    GetEvidenceClass() string
    GetRule() string
}
```

Then `bmistv1.RetirementRule` already satisfies this interface (it has `GetEvidenceClass()` and `GetRule()` methods).

**File:** `domain-pack/registry.go`
```go
// RetirementRule describes the evidence-class gate for retiring a debt item.
type RetirementRule interface {
    GetEvidenceClass() string
    GetRule() string
}

// RetirementRulesProvider is implemented by Pack types that declare retirement rules.
type RetirementRulesProvider interface {
    GetRetirementRules() map[string]RetirementRule
}
```

### Step C: Fix `decodePack` to return typed packs

**File:** `domain-pack/registry.go:170-193`

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

This requires importing `bmistv1` from `domain-pack/registry.go`. Check for import cycle: `domain-pack/` → `domain-pack/bmist/v1/` — no cycle, `bmistv1` does not import `domain-pack`.

### Step D: Update coordinator validation to use typed interface

**File:** `coordinator/validate.go:103-137`

Replace `getRetirementRule` with the typed version:
```go
func getRetirementRule(debtItem string, pack domainpack.Pack) (*RetirementRuleInfo, error) {
    if p, ok := pack.(domainpack.RetirementRulesProvider); ok {
        rules := p.GetRetirementRules()
        if rule, exists := rules[debtItem]; exists {
            return &RetirementRuleInfo{
                EvidenceClass: rule.GetEvidenceClass(),
                Rule:          rule.GetRule(),
            }, nil
        }
    }
    return nil, nil
}
```

No more `interface{}` assertions. The `bmistv1.RetirementRule` struct satisfies `domainpack.RetirementRule` interface via its `GetEvidenceClass()` and `GetRule()` methods.

Also update `RetirementRuleProvider` interface in `coordinator/validate.go`:
```go
type RetirementRuleProvider interface {
    GetRetirementRules() map[string]domainpack.RetirementRule
}
```

### Step E: Give App access to pack registry

Add to `App` struct:
```go
type App struct {
    // ... existing fields ...
    packRegistry *domainpack.PackRegistry
}
```

Add a setter:
```go
func (a *App) SetPackRegistry(r *domainpack.PackRegistry) {
    a.packRegistry = r
}
```

Add a `resolvePack` helper:
```go
func (a *App) resolvePack(scenarioID string) (domainpack.Pack, error) {
    if a.packRegistry == nil {
        return nil, fmt.Errorf("pack registry unavailable")
    }
    // For POC: all scenarios use bmist@v1.0.0
    return a.packRegistry.Get("bmist", "v1.0.0")
}
```

### Step F: Add EvidenceClass to DecisionRequest

```go
type DecisionRequest struct {
    Type          string `json:"type"`
    ScenarioID    string `json:"scenario_id"`
    BeliefID      string `json:"belief_id"`
    ObligationKey string `json:"obligation_key,omitempty"`
    InstrumentRef string `json:"instrument_ref,omitempty"`
    PrincipalID   string `json:"principal_id"`
    EvidenceClass string `json:"evidence_class,omitempty"`
}
```

### Step G: Full retirement-rule enforcement in SubmitDecision

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

    // 3. Debt item must have a retirement rule in the pack
    rulesProv, ok := pack.(domainpack.RetirementRulesProvider)
    if !ok {
        return fmt.Errorf("pack does not support retirement rules")
    }
    rules := rulesProv.GetRetirementRules()
    rule, exists := rules[req.ObligationKey]
    if !exists {
        return fmt.Errorf("no retirement rule for debt item: %s", req.ObligationKey)
    }

    // 4. EvidenceClass is REQUIRED — empty is not allowed
    requiredClass := rule.GetEvidenceClass()
    if req.EvidenceClass == "" {
        return fmt.Errorf("evidence_class is required for debt discharge of: %s", req.ObligationKey)
    }
    if req.EvidenceClass != requiredClass {
        return fmt.Errorf("evidence class mismatch: debt %q requires %q, got %q",
            req.ObligationKey, requiredClass, req.EvidenceClass)
    }

    // 5. Verify qualifying persisted evidence exists
    evidenceIDs, err := a.verifyPersistedEvidence(ctx, req.BeliefID, req.ScenarioID, req.EvidenceClass)
    if err != nil {
        return fmt.Errorf("evidence verification failed: %w", err)
    }
    if len(evidenceIDs) == 0 {
        return fmt.Errorf("no qualifying evidence of class %q for belief %s in scenario %s",
            req.EvidenceClass, req.BeliefID, req.ScenarioID)
    }

    // 6. Build instrument ref from verified evidence
    instrumentRef := buildInstrumentRef(evidenceIDs)

    return a.kern.Discharge(ctx, req.ScenarioID, req.BeliefID, req.ObligationKey, instrumentRef, req.PrincipalID)
```

### Step H: Query persisted evidence

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

### Step I: Wire pack registry into App construction

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

### Step J: Update Trust UI to pass EvidenceClass

**File:** `internal/ui/ui.go` — `HandleDischarge`

Add `EvidenceClass` to the request struct and pass it through:

```go
var req struct {
    ScenarioID    string `json:"scenario_id"`
    BeliefID      string `json:"belief_id"`
    ObligationKey string `json:"obligation_key"`
    EvidenceClass string `json:"evidence_class"`
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

**Fix:** Static token + Bearer authentication with browser-compatible flow.

**CRITICAL CONSTRAINT:** `ARGUS_OPERATOR_TOKEN` must NOT have a default. An absent credential must make consequential actions unavailable, not silently install a known credential.

### Design

```
ARGUS_OPERATOR_TOKEN env var (REQUIRED for UI)
        ↓
Missing → startup error: "ARGUS_OPERATOR_TOKEN required for UI"
        ↓
Present → Server stores token, exposes login endpoint
        ↓
Browser: POST /ui/api/login { token: "..." }
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

2. `NewServer` requires token:
```go
func NewServer(app *application.App, operatorToken string) *Server {
    if operatorToken == "" {
        panic("ARGUS_OPERATOR_TOKEN is required for Trust UI")  // fail-closed at construction
    }
    tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))
    return &Server{app: app, tmpl: tmpl, operatorToken: operatorToken}
}
```

3. Add login endpoint (sets HttpOnly cookie):
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

4. Add `requireAuth` — checks cookie OR Bearer header:
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

5. Add Origin/CSRF check for write endpoints:
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

6. Protect write endpoints:
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

7. Server-derived principal from token validation (never from request body):
```go
func (s *Server) HandleDischarge(w http.ResponseWriter, r *http.Request) {
    // ... decode request ...
    // Principal is always "operator" — derived from authenticated session
    principalID := "operator"
    // ... rest of handler
}
```

8. Wire token from `cmdServe` — **FAIL CLOSED**:
```go
operatorToken := os.Getenv("ARGUS_OPERATOR_TOKEN")
if operatorToken == "" {
    log.Fatal("ARGUS_OPERATOR_TOKEN is required for Trust UI; set it before starting")
}
uiServer := ui.NewServer(app, operatorToken)
```

For `cmdMCP` (which does not serve UI), `NewServer` is not called, so the token is irrelevant.

**Tests:**
- `TestDischargeRequiresAuth` — POST without token/cookie → 401
- `TestDischargeWithValidToken` — POST with Bearer token → 200
- `TestDischargeWithValidCookie` — POST with session cookie → 200
- `TestDischargeRejectsBodyPrincipal` — body principal_id is ignored
- `TestPromoteRequiresAuth` — POST without token → 401
- `TestNewServerRequiresToken` — empty token → panic

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

    // Load pack registry
    packReg := domainpack.NewRegistry()
    if err := packReg.LoadFromDisk("domain-pack"); err != nil {
        log.Printf("warning: pack load failed: %v", err)
    }
    app.SetPackRegistry(packReg)

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

**Fix:** For `reproducible_artifact` evidence, the rules are:
- `ArtifactRef` is **REQUIRED** (no agent-supplied hash without a trusted artifact)
- Trusted registry artifact **REQUIRED** (artifact must exist in the registry, registered via trusted path only)
- Hash must match the registry's `ArtifactHash`
- Verifier must be authorized by the pack

For `operator_asserted` evidence, the agent-attested hash is accepted (no artifact registry lookup).

```go
// In Persist, after resolving beliefID for each evidence item:
if e.ProvenanceClass == "reproducible_artifact" {
    // ArtifactRef is REQUIRED for reproducible_artifact
    if e.ArtifactRef == "" {
        return nil, fmt.Errorf("evidence[%d]: reproducible_artifact requires artifact_ref", i)
    }
    if a.artifactReg == nil {
        return nil, fmt.Errorf("evidence[%d]: artifact registry unavailable for reproducible_artifact", i)
    }
    artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
    if err != nil {
        return nil, fmt.Errorf("evidence[%d]: trusted artifact not found: %w", i, err)
    }
    // Hash must match
    if artifact.ArtifactHash != e.ContentSHA256 {
        return nil, fmt.Errorf("evidence[%d]: content hash mismatch: agent says %s, trusted artifact has %s",
            i, e.ContentSHA256, artifact.ArtifactHash)
    }
    // Verifier must be authorized
    spec := a.findVerifierSpec(artifact.VerifierID)
    if spec == nil {
        return nil, fmt.Errorf("evidence[%d]: verifier %s not authorized by pack", i, artifact.VerifierID)
    }
    if !versionGTE(artifact.VerifierVersion, spec.MinVersion) {
        return nil, fmt.Errorf("evidence[%d]: verifier %s version %s below minimum %s",
            i, artifact.VerifierID, artifact.VerifierVersion, spec.MinVersion)
    }
}
// operator_asserted: accept agent-supplied hash (no registry lookup)
```

**Tests:**
- `TestReproducibleArtifactRequiresArtifactRef` — reproducible_artifact without artifact_ref → error
- `TestReproducibleArtifactRequiresTrustedArtifact` — reproducible_artifact with unknown ref → error
- `TestReproducibleArtifactHashMismatch` — wrong hash → error
- `TestReproducibleArtifactVerifierNotAuthorized` — verifier not in pack → error
- `TestReproducibleArtifactHappyPath` — all checks pass → success
- `TestOperatorAssertedAcceptsAgentHash` — operator_asserted with agent-supplied hash → accepted

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
argus resolve-ports (new subcommand)
    ↓
print "CRDB_SQL_PORT=NNNN CRDB_ADMIN_PORT=NNNN ARGUS_PORT=NNNN"
    ↓
Taskfile source <(argus resolve-ports)
    ↓
start CRDB with --listen-addr :$CRDB_SQL_PORT --http-addr :$CRDB_ADMIN_PORT
    ↓
wait for SQL readiness
    ↓
argus reset --db postgres://root@localhost:$CRDB_SQL_PORT/argus
    ↓
argus serve --db ... --listen :$ARGUS_PORT &
    ↓
EADDRINUSE → re-resolve ports, retry
```

### New subcommand: `resolve-ports`

**File:** `cmd/argus/main.go`

```go
case "resolve-ports":
    cmdResolvePorts(os.Args[2:])
```

```go
func cmdResolvePorts(args []string) {
    fs := flag.NewFlagSet("resolve-ports", flag.ExitOnError)
    crdbSQLPreferred := fs.Int("crdb-sql", 26257, "preferred CRDB SQL port")
    crdbAdminPreferred := fs.Int("crdb-admin", 8081, "preferred CRDB admin port")
    argusPreferred := fs.Int("argus", 8080, "preferred ARGUS HTTP port")
    fs.Parse(args)

    // Resolve CRDB SQL port
    crdbSQL := resolvePort(*crdbSQLPreferred)
    // Resolve CRDB admin port (must differ from CRDB SQL)
    crdbAdmin := resolvePortDifferent(*crdbAdminPreferred, crdbSQL)
    // Resolve ARGUS HTTP port
    argus := resolvePortDifferent(*argusPreferred, crdbSQL, crdbAdmin)

    // Print shell-compatible assignments
    fmt.Printf("CRDB_SQL_PORT=%d\n", crdbSQL)
    fmt.Printf("CRDB_ADMIN_PORT=%d\n", crdbAdmin)
    fmt.Printf("ARGUS_PORT=%d\n", argus)
}

func resolvePort(preferred int) int {
    addr := fmt.Sprintf(":%d", preferred)
    ln, err := net.ListenTimeout("tcp", addr, 500*time.Millisecond)
    if err == nil {
        ln.Close()
        return preferred
    }
    // Fallback: try next 10 ports
    for p := preferred + 1; p <= preferred+10; p++ {
        ln, err := net.ListenTimeout("tcp", fmt.Sprintf(":%d", p), 500*time.Millisecond)
        if err == nil {
            ln.Close()
            return p
        }
    }
    return preferred // last resort — let the service fail with a clear error
}

func resolvePortDifferent(preferred int, occupied ...int) int {
    occupiedSet := make(map[int]bool, len(occupied))
    for _, p := range occupied {
        occupiedSet[p] = true
    }
    for {
        port := resolvePort(preferred)
        if !occupiedSet[port] {
            return port
        }
        preferred = port + 1
    }
}
```

**Key correction:** Unlike the `Reserve` pattern from plan10r, this approach does NOT hold the listener open. It resolves the port, releases, and relies on `EADDRINUSE` retry at startup time. This eliminates the complexity of coordinating reserve-release-startup timing. The tradeoff is a small race window between resolution and bind, but the service retries on `EADDRINUSE`.

### Retry on EADDRINUSE in cmdServe

```go
func cmdServe(args []string) {
    // ... parse flags, open DB ...

    for attempt := 0; attempt < 3; attempt++ {
        srv := &http.Server{
            Addr:         *listen,
            Handler:      mux,
            ReadTimeout:  15 * time.Second,
            WriteTimeout: 15 * time.Second,
            IdleTimeout:  60 * time.Second,
        }

        err := srv.ListenAndServe()
        if err == http.ErrServerClosed {
            break // graceful shutdown
        }
        if strings.Contains(err.Error(), "address already in use") {
            log.Printf("port %s occupied, attempting to resolve new port...", *listen)
            // Re-resolve: try listen+1 through listen+10
            base := extractPort(*listen)
            *listen = fmt.Sprintf(":%d", resolvePort(base+1))
            log.Printf("retrying on %s (attempt %d/3)", *listen, attempt+1)
            continue
        }
        log.Fatalf("serve: %v", err)
    }
}
```

### Taskfile changes

```yaml
  dev:
    desc: Start ARGUS in development mode (non-destructive)
    cmds:
      - mkdir -p .argus-pids
      - |
        # Resolve ports
        eval "$(go run ./cmd/argus resolve-ports)"
        echo "resolved: crdb-sql=$CRDB_SQL_PORT crdb-admin=$CRDB_ADMIN_PORT argus=$ARGUS_PORT"
        
        # Start CRDB
        cockroach start-single-node --insecure \
          --listen-addr :$CRDB_SQL_PORT \
          --http-addr :$CRDB_ADMIN_PORT \
          --store=.cockroach-data &
        echo $! > .argus-pids/crdb.pid
        
        # Wait for SQL readiness
        until cockroach node status --host :$CRDB_SQL_PORT --insecure 2>/dev/null; do
          sleep 1
        done
        
        # Apply migrations
        go run ./cmd/argus reset --db postgres://root@localhost:$CRDB_SQL_PORT/argus?sslmode=disable
        
        # Start ARGUS (background)
        go run ./cmd/argus serve \
          --db "postgres://root@localhost:$CRDB_SQL_PORT/argus?sslmode=disable" \
          --listen ":$ARGUS_PORT" &
        echo $! > .argus-pids/argus.pid
        
        # Trap for cleanup
        trap 'for f in .argus-pids/*.pid; do [ -f "$f" ] && kill $(cat "$f") 2>/dev/null; rm -f "$f"; done' EXIT INT TERM
        
        # Keep script alive for trap
        wait
```

### Port allocation tests

- `TestResolvePortPreferredAvailable` — preferred port free → gets preferred
- `TestResolvePortFallback` — preferred occupied → gets next available
- `TestResolvePortDifferent` — avoids occupied ports
- `TestCRDBAdminARGUSCollision` — CRDB admin on 8080, ARGUS tries 8080 → resolves to 8081

---

## Fix 11: task down + PID cleanup

**File:** `Taskfile.yml`

**Add:**
```yaml
  down:
    desc: Stop background CRDB and ARGUS processes
    cmds:
      - |
        for pidfile in .argus-pids/*.pid; do
          [ -f "$pidfile" ] || continue
          pid=$(cat "$pidfile")
          kill "$pid" 2>/dev/null || true
          rm "$pidfile"
        done
```

**Modify `dev` and `fresh`** to record PIDs in `.argus-pids/` and add trap for cleanup.

---

## Fix 12: Integration test CRDB config

**Problem:** Tests default to port 26257. Taskfile starts CRDB on 26260. Port mismatch.

**Fix:** Integration tests consume the runtime-resolved port via `ARGUS_CRDB_PORT` env var:

```go
func integrationDB(t *testing.T) *sql.DB {
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

Also update the `testDSN` construction within the function to use the same resolved port.

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

- **No third Solvent change.** Changes A (migration export) and B (decodePack/typed rules) are the two allowed.
- **No Coordinator REST resurrection.** All fixes work within the single-process architecture.
- **No separate Trust UI process.** Auth is embedded in `internal/ui`.
- **No new architecture layers.** Fixes are within existing packages.
- **Preserve all existing passing tests.** No test deletions.
- **`ARGUS_OPERATOR_TOKEN` has no default.** Missing token → startup failure.
- **`GetRetirementRules()` returns typed map.** No `interface{}` detours.
- **`reproducible_artifact` requires trusted artifact.** ArtifactRef mandatory, hash verified, verifier authorized.
- **Migrations are genuinely idempotent.** All SQL uses `IF NOT EXISTS`.
- **Port resolution retries on EADDRINUSE.** No probe-then-close races.
