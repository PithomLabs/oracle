# Plan 10r3 — Corrective Implementation Pass (rev3)

**Date:** 2026-09-18 (revision 3)
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
| 5 | Retirement-rule enforcement before Discharge | Medium | 4 (registry.go, bmist types.go, app.go, validate.go) + tests |
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

**Scope classification:** This is **one bounded Solvent migration-packaging change** that also makes the existing migration scripts replay-safe. It includes:
1. Moving canonical SQL to `solvent-main/migrations/db/`
2. Making bare `CREATE TABLE`/`CREATE INDEX` statements idempotent (`IF NOT EXISTS`)
3. Creating the `solvent-main/migrations/` package with `Apply(embed.FS)`
4. Reusing Solvent's existing `splitStatements` logic (not maintaining two implementations)

The semantic behavior of the migrations is unchanged — the same tables, columns, constraints, and defaults are created. The only behavioral difference is that `Apply()` is now safe to call multiple times without errors.

### Step A: Create `solvent-main/migrations/db/` with SQL files

`go:embed` cannot embed files from `../db/` — the path must stay within the package directory. Therefore, create `solvent-main/migrations/db/` and copy the 10 SQL files there. This becomes the single canonical location for the migration package.

### Step B: Make existing SQL statements idempotent

Add `IF NOT EXISTS` to all bare statements in these files:

**`001_schema.sql`** (4 changes):
- `CREATE TABLE belief` → `CREATE TABLE IF NOT EXISTS belief`
- `CREATE TABLE belief_edge` → `CREATE TABLE IF NOT EXISTS belief_edge`
- `CREATE INDEX belief_edge_child` → `CREATE INDEX IF NOT EXISTS belief_edge_child`
- `CREATE INDEX live_intents` → `CREATE INDEX IF NOT EXISTS live_intents`

**`007_service_tables.sql`** (8 changes):
- `CREATE TABLE workflow_token` → `CREATE TABLE IF NOT EXISTS workflow_token`
- `CREATE TABLE policy_tool` → `CREATE TABLE IF NOT EXISTS policy_tool`
- `CREATE TABLE policy_actor` → `CREATE TABLE IF NOT EXISTS policy_actor`
- `CREATE TABLE audit_activity` → `CREATE TABLE IF NOT EXISTS audit_activity`
- `CREATE INDEX workflow_token_scenario` → `CREATE INDEX IF NOT EXISTS workflow_token_scenario`
- `CREATE INDEX workflow_token_state` → `CREATE INDEX IF NOT EXISTS workflow_token_state`
- `CREATE INDEX audit_activity_scenario` → `CREATE INDEX IF NOT EXISTS audit_activity_scenario`
- `CREATE INDEX audit_activity_refusal` → `CREATE INDEX IF NOT EXISTS audit_activity_refusal`

All other files (002-006, 008-010) already use `IF NOT EXISTS` or equivalent.

### Step C: Create `solvent-main/migrations/migrations.go`

Reuse the existing `splitStatements` logic from `internal/testdb/testdb.go` (same algorithm: strip `--` comments, split on `;`):

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

// splitStatements strips line comments and splits on semicolons.
// Reuses the same algorithm as internal/testdb/testdb.go.
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

### Step D: Update internal Solvent references

Update the following to use `migrations/db/` instead of `db/`:
- `internal/testdb/testdb.go` — callers of `ApplySchema`
- `internal/m0/schema.go` — DDL applier references
- `cmd/solvent/main.go` — `resolveSchemaPaths()`
- All test suites that hardcode schema paths
- `Taskfile.yml` — `db:reset` task
- `demo/cloud/init/main.go` — hardcoded paths

### Step E: Wire into Oracle's `cmdReset`

**File:** `cmd/argus/main.go:182-183`

```go
import solventmigrations "github.com/PithomLabs/solvent/migrations"

// in cmdReset:
log.Println("argus reset: applying Solvent migrations...")
if err := solventmigrations.Apply(ctx, db); err != nil {
    log.Fatalf("solvent migrations: %v", err)
}
```

### Step F: Replace test helpers

Replace filesystem-based test helpers with the package:
- `integration_test.go:integApplySolventMigrations` → `solventmigrations.Apply(ctx, db)`
- `internal/ui/ui_test.go:applySolventMigrations` → `solventmigrations.Apply(ctx, db)`
- `internal/application/app_test.go:applySolventMigrations` → `solventmigrations.Apply(ctx, db)`

**Tests:**
- `TestCmdResetAppliesSolventMigrations` — run against empty CRDB, verify Solvent tables exist
- `TestSolventMigrationsAreIdempotent` — apply twice, no errors

---

## Fix 3: Edge kind conflict detection (concurrency-safe)

**File:** `internal/application/app.go:257-264`

**Problem:** `ON CONFLICT (parent_id, child_id) DO NOTHING` silently drops a `contradicts` edge when a `derives` edge already exists for the same pair. Two concurrent transactions can both observe no row, then one inserts and the other gets `DO NOTHING`.

**Fix:** Use `RowsAffected()` to detect the race, then re-read to confirm:

```go
for _, edge := range pkt.Edges {
    fromID := resolveRef(edge.FromRef, result.BeliefIDs)
    toID := resolveRef(edge.ToRef, result.BeliefIDs)
    if fromID == "" || toID == "" {
        return nil, fmt.Errorf("edge references unresolved: from=%s to=%s", edge.FromRef, edge.ToRef)
    }

    result2, err := tx.ExecContext(ctx,
        `INSERT INTO belief_edge (parent_id, child_id, kind)
         VALUES ($1::UUID, $2::UUID, $3)
         ON CONFLICT (parent_id, child_id) DO NOTHING`,
        fromID, toID, edge.Kind)
    if err != nil {
        return nil, fmt.Errorf("insert edge: %w", err)
    }

    rowsAffected, err := result2.RowsAffected()
    if err != nil {
        return nil, fmt.Errorf("edge rows affected: %w", err)
    }

    if rowsAffected == 0 {
        // Insert was a no-op: an edge already exists for this (parent, child) pair.
        // Re-read the existing kind to distinguish idempotent same-kind from conflict.
        var existingKind string
        err := tx.QueryRowContext(ctx,
            `SELECT kind FROM belief_edge WHERE parent_id = $1::UUID AND child_id = $2::UUID`,
            fromID, toID).Scan(&existingKind)
        if err != nil {
            return nil, fmt.Errorf("edge conflict re-read failed: %w", err)
        }
        if existingKind != edge.Kind {
            return nil, fmt.Errorf("edge conflict: %s->%s already has kind=%s, cannot add kind=%s",
                fromID, toID, existingKind, edge.Kind)
        }
        // Same kind — idempotent, no error
    }

    result.EdgeCount++
}
```

**Why this is safe:** After `RowsAffected() == 0`, we re-read within the same transaction. If another transaction committed a different kind, we see it and return an error. If the same kind was already there, we treat it as idempotent. The re-read is within the same `tx`, so we see our own transaction's snapshot.

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
- `domain-pack/registry.go` — extend `Pack` interface, add generic types
- `domain-pack/bmist/v1/types.go` — use generic types
- `coordinator/validate.go` — simplify type assertions
- `internal/application/app.go` — `SubmitDecision`, `App` struct, `DecisionRequest`

**Problem:** `SubmitDecision("discharge")` calls `kern.Discharge()` directly without retirement-rule enforcement. Additionally, the `GetRetirementRules()` method returns `map[string]interface{}` which forces fragile double type-assertions downstream. The `genericPack` discards all domain-specific data.

### Type architecture

The existing `domainpack.Pack` interface is intentionally thin (`GetPackID()`, `GetVersion()`). All other access uses ad-hoc type assertions. The fix extends the interface with typed accessors that return generic domain types, eliminating the need for type assertions in the application layer.

**Key constraint:** `map[string]bmistv1.RetirementRule` does NOT implement `map[string]domainpack.RetirementRule` in Go (map value type invariance). Therefore, the BM-IST pack struct must use `domainpack.RetirementRule` as its field type directly.

### Step A: Add generic types to `domain-pack/registry.go`

```go
// RetirementRule describes the evidence-class gate for retiring a debt item.
// Plain struct — no methods, no interface. The application reads fields directly.
type RetirementRule struct {
    EvidenceClass string `json:"evidence_class"`
    Rule          string `json:"rule"`
}

// VerifierSpec declares which verifier is authorized by the pack.
type VerifierSpec struct {
    VerifierID string `json:"verifier_id"`
    MinVersion string `json:"min_version"`
}
```

### Step B: Extend `Pack` interface

```go
// Pack is the interface that all Domain Pack types must implement.
type Pack interface {
    GetPackID() string
    GetVersion() string
    GetDebtVocabulary() []string
    GetEvidenceClasses() []string
    GetRetirementRules() map[string]RetirementRule
    GetVerifierSpecs() []VerifierSpec
}
```

### Step C: Update `bmistv1.Pack` to use generic types

**File:** `domain-pack/bmist/v1/types.go`

```go
import "github.com/PithomLabs/oracle/domain-pack"

type Pack struct {
    PackID                string                           `json:"pack_id"`
    Version               string                           `json:"version"`
    // ... other fields ...
    DebtVocabulary        []string                         `json:"debt_vocabulary"`
    EvidenceClasses       []string                         `json:"evidence_classes"`
    RetirementRules       map[string]domainpack.RetirementRule `json:"retirement_rules"`
    VerifierSpecs         []domainpack.VerifierSpec        `json:"verifier_specs"`
}

func (p *Pack) GetPackID() string              { return p.PackID }
func (p *Pack) GetVersion() string             { return p.Version }
func (p *Pack) GetDebtVocabulary() []string    { return p.DebtVocabulary }
func (p *Pack) GetEvidenceClasses() []string   { return p.EvidenceClasses }
func (p *Pack) GetRetirementRules() map[string]domainpack.RetirementRule {
    return p.RetirementRules
}
func (p *Pack) GetVerifierSpecs() []domainpack.VerifierSpec {
    return p.VerifierSpecs
}
```

Remove the old `bmistv1.RetirementRule` and `bmistv1.VerifierSpec` types. The BM-IST pack's `ConsequentialAction` type stays domain-specific (it's not consumed by the application layer).

**Import cycle check:** `bmistv1` imports `domainpack` — safe, no reverse dependency.

### Step D: Fix `decodePack` to return typed packs

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

### Step E: Update `genericPack` to satisfy extended interface

```go
type genericPack struct {
    PackID  string `json:"pack_id"`
    Version string `json:"version"`
}

func (g *genericPack) GetPackID() string              { return g.PackID }
func (g *genericPack) GetVersion() string             { return g.Version }
func (g *genericPack) GetDebtVocabulary() []string    { return nil }
func (g *genericPack) GetEvidenceClasses() []string   { return nil }
func (g *genericPack) GetRetirementRules() map[string]domainpack.RetirementRule { return nil }
func (g *genericPack) GetVerifierSpecs() []domainpack.VerifierSpec { return nil }
```

### Step F: Simplify coordinator validation

**File:** `coordinator/validate.go`

Replace the double type-assertion `getRetirementRule` with direct typed access:

```go
// ValidateDebtMembership — use interface method directly
func packHasDebt(pack domainpack.Pack, debt string) bool {
    for _, d := range pack.GetDebtVocabulary() {
        if d == debt {
            return true
        }
    }
    return false
}

// getRetirementRule — direct typed access, no assertions
func getRetirementRule(debtItem string, pack domainpack.Pack) (*domainpack.RetirementRule, error) {
    rules := pack.GetRetirementRules()
    if rules == nil {
        return nil, nil
    }
    rule, exists := rules[debtItem]
    if !exists {
        return nil, nil
    }
    return &rule, nil
}

// ValidateRetirementRule — use rule.EvidenceClass directly
func ValidateRetirementRule(debtItem string, evidenceClass string, pack domainpack.Pack) error {
    rule, err := getRetirementRule(debtItem, pack)
    if err != nil { return err }
    if rule == nil {
        return fmt.Errorf("no retirement rule defined for debt item %q", debtItem)
    }
    if rule.EvidenceClass != evidenceClass {
        return fmt.Errorf("retirement rule mismatch: debt %q requires evidence class %q, got %q",
            debtItem, rule.EvidenceClass, evidenceClass)
    }
    return nil
}
```

Remove `RetirementRuleInfo` struct and `RetirementRuleProvider` interface — no longer needed.

### Step G: Simplify `GetPackRules`

**File:** `coordinator/coordinator.go:110-138`

```go
func (c *Coordinator) GetPackRules() map[string]interface{} {
    result := make(map[string]interface{})
    if c.packRegistry == nil { return result }
    pack, err := c.packRegistry.Get("bmist", "1.0.0")
    if err != nil { return result }
    rules := pack.GetRetirementRules()
    if rules == nil { return result }
    result["rules"] = rules // domainpack.RetirementRule marshals directly to JSON
    return result
}
```

No more JSON round-trip conversion.

### Step H: Give App access to pack registry

Add to `App` struct:
```go
type App struct {
    // ... existing fields ...
    packRegistry *domainpack.PackRegistry
}
```

Add a setter and resolver:
```go
func (a *App) SetPackRegistry(r *domainpack.PackRegistry) { a.packRegistry = r }

func (a *App) resolvePack(scenarioID string) (domainpack.Pack, error) {
    if a.packRegistry == nil {
        return nil, fmt.Errorf("pack registry unavailable")
    }
    return a.packRegistry.Get("bmist", "v1.0.0")
}
```

### Step I: Add EvidenceClass to DecisionRequest

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

### Step J: Full retirement-rule enforcement in SubmitDecision

Replace the `discharge` case:

```go
case "discharge":
    if a.packRegistry == nil {
        return fmt.Errorf("pack registry unavailable: cannot validate retirement")
    }
    pack, err := a.resolvePack(req.ScenarioID)
    if err != nil {
        return fmt.Errorf("pack resolution failed: %w", err)
    }

    // Debt item must have a retirement rule (fail-closed)
    rules := pack.GetRetirementRules()
    if rules == nil {
        return fmt.Errorf("pack does not support retirement rules")
    }
    rule, exists := rules[req.ObligationKey]
    if !exists {
        return fmt.Errorf("no retirement rule for debt item: %s", req.ObligationKey)
    }

    // EvidenceClass is REQUIRED
    if req.EvidenceClass == "" {
        return fmt.Errorf("evidence_class is required for debt discharge of: %s", req.ObligationKey)
    }
    if req.EvidenceClass != rule.EvidenceClass {
        return fmt.Errorf("evidence class mismatch: debt %q requires %q, got %q",
            req.ObligationKey, rule.EvidenceClass, req.EvidenceClass)
    }

    // Verify qualifying persisted evidence exists
    evidenceIDs, err := a.verifyPersistedEvidence(ctx, req.BeliefID, req.ScenarioID, req.EvidenceClass)
    if err != nil {
        return fmt.Errorf("evidence verification failed: %w", err)
    }
    if len(evidenceIDs) == 0 {
        return fmt.Errorf("no qualifying evidence of class %q for belief %s in scenario %s",
            req.EvidenceClass, req.BeliefID, req.ScenarioID)
    }

    instrumentRef := buildInstrumentRef(evidenceIDs)
    return a.kern.Discharge(ctx, req.ScenarioID, req.BeliefID, req.ObligationKey, instrumentRef, req.PrincipalID)
```

### Step K: Query persisted evidence

```go
func (a *App) verifyPersistedEvidence(ctx context.Context, beliefID, scenarioID, evidenceClass string) ([]string, error) {
    rows, err := a.db.QueryContext(ctx,
        `SELECT id FROM evidence
         WHERE belief_id = $1::UUID AND scenario_id = $2::UUID AND provenance_class = $3`,
        beliefID, scenarioID, evidenceClass)
    if err != nil { return nil, err }
    defer rows.Close()
    var ids []string
    for rows.Next() {
        var id string
        if err := rows.Scan(&id); err != nil { return nil, err }
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

### Step L: Wire pack registry and update Trust UI

**File:** `cmd/argus/main.go` — `cmdServe` and `cmdMCP`

```go
packReg := domainpack.NewRegistry()
if err := packReg.LoadFromDisk("domain-pack"); err != nil {
    log.Printf("warning: pack load failed: %v", err)
}
app := application.NewWithRegistry(db, registry, specs)
app.SetPackRegistry(packReg)
```

**File:** `internal/ui/ui.go` — `HandleDischarge` adds `EvidenceClass` field.

**Tests:**
- `TestDischargeRequiresKnownDebtItem`
- `TestDischargeRequiresEvidenceClass`
- `TestDischargeWithMismatchedEvidenceClass`
- `TestDischargeRequiresPersistedEvidence`
- `TestDischargeWithQualifyingEvidence`
- `TestDischargeWithRegistryUnavailable`

---

## Fix 6: Trust UI authentication

**Files:**
- `internal/ui/ui.go`
- `cmd/argus/main.go`

**Problem:** UI hardcodes principal ID, has no authentication, no CSRF protection.

**CRITICAL CONSTRAINT:** `ARGUS_OPERATOR_TOKEN` must NOT have a default. An absent credential must make consequential actions unavailable.

### Implementation

1. `NewServer` requires token — fail-closed at construction:
```go
func NewServer(app *application.App, operatorToken string) *Server {
    if operatorToken == "" {
        panic("ARGUS_OPERATOR_TOKEN is required for Trust UI")
    }
    tmpl := template.Must(template.ParseFS(templateFS, "templates/*.html"))
    return &Server{app: app, tmpl: tmpl, operatorToken: operatorToken}
}
```

2. Login endpoint (sets HttpOnly cookie):
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

3. `requireAuth` — checks cookie OR Bearer header:
```go
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("argus_session")
        if err == nil && cookie.Value == s.operatorToken {
            next(w, r)
            return
        }
        authHeader := r.Header.Get("Authorization")
        if authHeader == "Bearer "+s.operatorToken {
            next(w, r)
            return
        }
        http.Error(w, "unauthorized", http.StatusUnauthorized)
    }
}
```

4. Protect write endpoints; server-derived principal.

5. Wire token — **FAIL CLOSED**:
```go
operatorToken := os.Getenv("ARGUS_OPERATOR_TOKEN")
if operatorToken == "" {
    log.Fatal("ARGUS_OPERATOR_TOKEN is required for Trust UI; set it before starting")
}
uiServer := ui.NewServer(app, operatorToken)
```

**Tests:**
- `TestDischargeRequiresAuth`
- `TestDischargeWithValidToken`
- `TestDischargeWithValidCookie`
- `TestDischargeRejectsBodyPrincipal`
- `TestPromoteRequiresAuth`
- `TestNewServerRequiresToken`

---

## Fix 7: MCP stdio transport (official SDK)

**Files:**
- `cmd/argus/main.go:102-124`
- `internal/mcp/adapter.go`
- `go.mod` (new dependency)

**Problem:** `argus mcp` prints a message and blocks forever. No MCP protocol framing.

**Fix:** Use the official MCP Go SDK (`github.com/modelcontextprotocol/go-sdk/mcp`).

### Verified SDK API

```go
server := mcp.NewServer(&mcp.Implementation{Name: "argus", Version: "v0.1.0"}, nil)
mcp.AddTool(server, &mcp.Tool{Name: "argus.get_context", Description: "..."}, handlerFunc)
server.AddTool(mcp.Tool{Name: "argus.submit_packet", Description: "...", InputSchema: map[string]interface{}{...}}, rawHandler)
server.Run(ctx, &mcp.StdioTransport{})
```

### Implementation

```go
func cmdMCP(args []string) {
    fs := flag.NewFlagSet("mcp", flag.ExitOnError)
    dbURL := fs.String("db", "", "database URL (required)")
    fs.Parse(args)

    if *dbURL == "" {
        *dbURL = loadRuntimeDBURL() // reads from .argus-runtime.env
    }
    if *dbURL == "" {
        log.Fatal("mcp: --db flag or .argus-runtime.env required")
    }

    db, err := sql.Open("pgx", *dbURL)
    // ... ping, create app, load packs ...

    server := mcp.NewServer(&mcp.Implementation{Name: "argus", Version: "v0.1.0"}, nil)
    adapter := mcpinternal.NewAdapter(app)
    adapter.RegisterTools(server)

    if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
        log.Fatalf("mcp: %v", err)
    }
}
```

### Adapter registration

```go
func (a *Adapter) RegisterTools(server *mcp.Server) {
    mcp.AddTool(server, &mcp.Tool{
        Name:        "argus.get_context",
        Description: "Get the RCP/v1 context for a task: task, dependencies, epistemic state, and activity.",
    }, a.handleGetContextMCP)

    server.AddTool(mcp.Tool{
        Name:        "argus.submit_packet",
        Description: "Submit an EBP research packet. Validates, persists beliefs/evidence/edges/tasks. Idempotent.",
        InputSchema: map[string]interface{}{...},
    }, a.handleSubmitPacketRaw)
}
```

**Tests:**
- `TestMCPToolsRegistered`
- `TestMCPGetContextViaSDK`
- `TestMCPSubmitPacketViaSDK`

---

## Fix 8: Semantic version comparison

**File:** `internal/application/app.go:163-166`

**Problem:** `versionGTE` uses lexicographic string comparison. `"0.10.0" < "0.2.0"` lexicographically.

**Fix:** Implement proper semver comparison. Reject prerelease versions:

```go
func versionGTE(version, minVersion string) bool {
    v, err1 := parseSemver(version)
    m, err2 := parseSemver(minVersion)
    if err1 != nil || err2 != nil {
        return false
    }
    if v[0] != m[0] { return v[0] > m[0] }
    if v[1] != m[1] { return v[1] > m[1] }
    return v[2] >= m[2]
}

func parseSemver(v string) ([3]int, error) {
    var parts [3]int
    v = strings.TrimPrefix(v, "v")
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
- `TestVersionGTEEqual`, `TestVersionGTEMajorGreater`, `TestVersionGTEMinorGreater`
- `TestVersionGTEMinorLess`, `TestVersionGTEPatchGreater`, `TestVersionGTERejectsPrerelease`

---

## Fix 9: Server-verify evidence hashes

**File:** `internal/application/app.go:198-214`

**Problem:** `content_sha256` is client-supplied and never verified.

**Fix:** For `reproducible_artifact` evidence, the rules are:
- `ArtifactRef` is **REQUIRED**
- Trusted registry artifact **REQUIRED**
- Hash must match the registry's `ArtifactHash`
- Verifier must be authorized by the pack

For `operator_asserted` evidence, agent-attested hash is accepted.

```go
if e.ProvenanceClass == "reproducible_artifact" {
    if e.ArtifactRef == "" {
        return nil, fmt.Errorf("evidence[%d]: reproducible_artifact requires artifact_ref", i)
    }
    if a.artifactReg == nil {
        return nil, fmt.Errorf("evidence[%d]: artifact registry unavailable", i)
    }
    artifact, err := a.artifactReg.Resolve(ctx, e.ArtifactRef)
    if err != nil {
        return nil, fmt.Errorf("evidence[%d]: trusted artifact not found: %w", i, err)
    }
    if artifact.ArtifactHash != e.ContentSHA256 {
        return nil, fmt.Errorf("evidence[%d]: hash mismatch: agent says %s, artifact has %s",
            i, e.ContentSHA256, artifact.ArtifactHash)
    }
    spec := a.findVerifierSpec(artifact.VerifierID)
    if spec == nil {
        return nil, fmt.Errorf("evidence[%d]: verifier %s not authorized", i, artifact.VerifierID)
    }
    if !versionGTE(artifact.VerifierVersion, spec.MinVersion) {
        return nil, fmt.Errorf("evidence[%d]: verifier %s version %s below minimum %s",
            i, artifact.VerifierID, artifact.VerifierVersion, spec.MinVersion)
    }
}
```

**Tests:**
- `TestReproducibleArtifactRequiresArtifactRef`
- `TestReproducibleArtifactRequiresTrustedArtifact`
- `TestReproducibleArtifactHashMismatch`
- `TestReproducibleArtifactVerifierNotAuthorized`
- `TestReproducibleArtifactHappyPath`
- `TestOperatorAssertedAcceptsAgentHash`

---

## Fix 10: Port preflight and lifecycle (BLOCKER FIX)

**Files:**
- `cmd/argus/main.go`
- `Taskfile.yml`
- Test files

**Problem:** No port preflight, no fallback, no cleanup. `:8080` CRDB admin/ARGUS collision. Port mismatch between Taskfile (`26260`) and tests (`26257`). No persistent runtime configuration.

### Architecture

```
argus resolve-ports → writes .argus-runtime.env
    ↓
task dev / task test / argus mcp / integration tests
    ↓
source .argus-runtime.env → all commands use resolved ports
    ↓
CRDB startup retries on EADDRINUSE → re-resolves ports
    ↓
ARGUS startup retries on EADDRINUSE → re-resolves ports
```

### Runtime configuration file

`argus resolve-ports` writes `.argus-runtime.env`:

```shell
CRDB_SQL_PORT=26257
CRDB_ADMIN_PORT=8081
ARGUS_PORT=8080
ARGUS_DB_URL=postgres://root@localhost:26257/argus?sslmode=disable
```

### New subcommand: `resolve-ports`

```go
func cmdResolvePorts(args []string) {
    fs := flag.NewFlagSet("resolve-ports", flag.ExitOnError)
    crdbSQLPreferred := fs.Int("crdb-sql", 26257, "preferred CRDB SQL port")
    crdbAdminPreferred := fs.Int("crdb-admin", 8081, "preferred CRDB admin port")
    argusPreferred := fs.Int("argus", 8080, "preferred ARGUS HTTP port")
    fs.Parse(args)

    crdbSQL := resolvePort(*crdbSQLPreferred)
    crdbAdmin := resolvePortDifferent(*crdbAdminPreferred, crdbSQL)
    argus := resolvePortDifferent(*argusPreferred, crdbSQL, crdbAdmin)

    // Write .argus-runtime.env
    env := fmt.Sprintf("CRDB_SQL_PORT=%d\nCRDB_ADMIN_PORT=%d\nARGUS_PORT=%d\nARGUS_DB_URL=postgres://root@localhost:%d/argus?sslmode=disable\n",
        crdbSQL, crdbAdmin, argus, crdbSQL)
    os.WriteFile(".argus-runtime.env", []byte(env), 0644)

    // Also print to stdout for shell sourcing
    fmt.Print(env)
}
```

### Port resolution functions

```go
func resolvePort(preferred int) int {
    addr := fmt.Sprintf(":%d", preferred)
    ln, err := net.ListenTimeout("tcp", addr, 500*time.Millisecond)
    if err == nil {
        ln.Close()
        return preferred
    }
    for p := preferred + 1; p <= preferred+10; p++ {
        ln, err := net.ListenTimeout("tcp", fmt.Sprintf(":%d", p), 500*time.Millisecond)
        if err == nil {
            ln.Close()
            return p
        }
    }
    return preferred
}

func resolvePortDifferent(preferred int, occupied ...int) int {
    occupiedSet := make(map[int]bool, len(occupied))
    for _, p := range occupied { occupiedSet[p] = true }
    for {
        port := resolvePort(preferred)
        if !occupiedSet[port] { return port }
        preferred = port + 1
    }
}
```

### CRDB startup with EADDRINUSE retry

```yaml
  dev:
    desc: Start ARGUS in development mode
    cmds:
      - mkdir -p .argus-pids
      - |
        trap 'for f in .argus-pids/*.pid; do [ -f "$f" ] && kill $(cat "$f") 2>/dev/null; rm -f "$f"; done' EXIT INT TERM
        
        # Resolve ports and write runtime config
        go run ./cmd/argus resolve-ports
        source .argus-runtime.env
        
        # Start CRDB with retry on EADDRINUSE
        for attempt in 1 2 3; do
          cockroach start-single-node --insecure \
            --listen-addr :$CRDB_SQL_PORT \
            --http-addr :$CRDB_ADMIN_PORT \
            --store=.cockroach-data 2>/dev/null &
          CRDB_PID=$!
          echo $CRDB_PID > .argus-pids/crdb.pid
          
          # Wait for SQL readiness
          if until cockroach node status --host :$CRDB_SQL_PORT --insecure 2>/dev/null; do
            if ! kill -0 $CRDB_PID 2>/dev/null; then
              # CRDB exited — likely port conflict
              echo "CRDB exited (attempt $attempt/3), re-resolving ports..."
              rm -f .argus-pids/crdb.pid
              go run ./cmd/argus resolve-ports
              source .argus-runtime.env
              break
            fi
            sleep 1
          done; then
            break
          fi
        done
        
        # Apply migrations
        go run ./cmd/argus reset --db "$ARGUS_DB_URL"
        
        # Start ARGUS
        go run ./cmd/argus serve --db "$ARGUS_DB_URL" --listen ":$ARGUS_PORT" &
        echo $! > .argus-pids/argus.pid
        
        wait
```

### `loadRuntimeDBURL` helper

```go
func loadRuntimeDBURL() string {
    data, err := os.ReadFile(".argus-runtime.env")
    if err != nil { return "" }
    for _, line := range strings.Split(string(data), "\n") {
        if strings.HasPrefix(line, "ARGUS_DB_URL=") {
            return strings.TrimPrefix(line, "ARGUS_DB_URL=")
        }
    }
    return ""
}
```

### Test DSN derivation

**File:** `integration_test.go`, `app_test.go`, `ui_test.go`

Replace hardcoded `26257` with derived port:

```go
func testDB(t *testing.T) *sql.DB {
    t.Helper()
    dsn := os.Getenv("ARGUS_TEST_DSN")
    if dsn == "" {
        dsn = loadRuntimeDBURL()
    }
    if dsn == "" {
        dsn = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
    }
    // Parse DSN to extract host:port for test database construction
    u, _ := url.Parse(dsn)
    hostPort := u.Host // e.g., "localhost:26257"
    // ... create test DB using hostPort ...
    testDSN := fmt.Sprintf("postgres://%s/%s?sslmode=disable", hostPort, dbName)
}
```

### Remove hardcoded CRDB port from all commands

`cmdMCP`, `cmdReset`, `cmdVerify` all default to `26260`. Change defaults to read from `.argus-runtime.env`:

```go
func cmdMCP(args []string) {
    fs := flag.NewFlagSet("mcp", flag.ExitOnError)
    dbURL := fs.String("db", loadRuntimeDBURL(), "database URL")
    fs.Parse(args)
    if *dbURL == "" {
        log.Fatal("mcp: --db flag or .argus-runtime.env required")
    }
    // ...
}
```

### Port allocation tests

- `TestResolvePortPreferredAvailable`
- `TestResolvePortFallback`
- `TestResolvePortDifferent`
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

---

## Fix 12: Integration test CRDB config

**Problem:** Tests default to port 26257. Taskfile starts CRDB on 26260. Port mismatch.

**Fix:** After Fix 10, integration tests consume `.argus-runtime.env` via `loadRuntimeDBURL()`. The `ARGUS_TEST_DSN` env var remains as an override.

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
7. Fix 5 (retirement rules) — needs registry/interface/bmist type changes
8. Fix 6 (UI auth) — needs token wiring
9. Fix 7 (MCP SDK) — needs `go get` + adapter rewrite
10. Fix 10 (port preflight) — needs port package + runtime env
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
- **`RetirementRule` is a plain struct in `domainpack`.** No methods, no interface, no type assertions.
- **`Pack` interface includes typed accessors.** `GetRetirementRules()`, `GetVerifierSpecs()`, `GetDebtVocabulary()`, `GetEvidenceClasses()`.
- **`reproducible_artifact` requires trusted artifact.** ArtifactRef mandatory, hash verified, verifier authorized.
- **Migrations are genuinely idempotent.** All SQL uses `IF NOT EXISTS`.
- **`.argus-runtime.env` is the persistent runtime config.** All commands consume it.
- **Port resolution retries on EADDRINUSE.** CRDB and ARGUS both retry.
