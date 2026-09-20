# Plan 10r4 — Corrective Implementation Pass (rev4, FINAL)

**Date:** 2026-09-18 (revision 5 — final before implementation)
**Source:** `docs/plan8/adv_review4.md` — net-valid defects only
**Architecture constraint:** Preserve the frozen single-process ARGUS architecture. Do NOT restore Solvent/Conductor/Coordinator REST services, separate Trust UI, or internal service-to-service HTTP.

---

## Scope

13 corrective items. Each is a standalone fix with tests. Ordered by dependency.

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
4. Reusing Solvent's existing `splitStatements` algorithm (not maintaining two implementations)

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

Reuse the existing `splitStatements` logic from `internal/testdb/testdb.go`:

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
    entries, err := FS.ReadDir("db")
    if err != nil {
        return fmt.Errorf("read embedded migrations: %w", err)
    }
    var sqlFiles []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
            sqlFiles = append(sqlFiles, e.Name())
        }
    }
    sort.Strings(sqlFiles)
    for _, name := range sqlFiles {
        data, err := FS.ReadFile("db/" + name)
        if err != nil {
            return fmt.Errorf("read migration %s: %w", name, err)
        }
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

**Note:** `ReadDir` and `ReadFile` errors are now handled explicitly — a broken embedded migration must fail closed.

### Step D: Update internal Solvent references

Update the following to use `migrations/db/` instead of `db/`:
- `internal/testdb/testdb.go` — callers of `ApplySchema`
- `internal/m0/schema.go` — DDL applier references
- `cmd/solvent/main.go` — `resolveSchemaPaths()`
- All test suites that hardcode schema paths
- `Taskfile.yml` — `db:reset` task
- `demo/cloud/init/main.go` — hardcoded paths

### Step E: Wire into Oracle — split `cmdReset` into `cmdMigrate` + `cmdReset`

**File:** `cmd/argus/main.go`

```go
func cmdMigrate(args []string) {
    // Non-destructive: applies Solvent + ARGUS migrations idempotently.
    // Safe to run on every startup. No data loss.
    fs := flag.NewFlagSet("migrate", flag.ExitOnError)
    dbURL := fs.String("db", loadRuntimeDBURL(), "database URL")
    fs.Parse(args)
    if *dbURL == "" { log.Fatal("migrate: --db flag or .argus-runtime.env required") }

    db, err := sql.Open("pgx", *dbURL)
    if err != nil { log.Fatalf("open db: %v", err) }
    defer db.Close()
    if err := db.Ping(); err != nil { log.Fatalf("ping db: %v", err) }

    ctx := context.Background()
    log.Println("argus migrate: applying Solvent migrations...")
    if err := solventmigrations.Apply(ctx, db); err != nil {
        log.Fatalf("solvent migrations: %v", err)
    }
    log.Println("argus migrate: applying ARGUS migrations...")
    if err := migrations.Apply(ctx, db); err != nil {
        log.Fatalf("argus migrations: %v", err)
    }
    log.Println("argus migrate: done")
}

func cmdReset(args []string) {
    // Destructive: drops all tables and re-applies. Use only for fresh start.
    fs := flag.NewFlagSet("reset", flag.ExitOnError)
    dbURL := fs.String("db", loadRuntimeDBURL(), "database URL")
    fs.Parse(args)
    if *dbURL == "" { log.Fatal("reset: --db flag or .argus-runtime.env required") }

    db, err := sql.Open("pgx", *dbURL)
    if err != nil { log.Fatalf("open db: %v", err) }
    defer db.Close()
    if err := db.Ping(); err != nil { log.Fatalf("ping db: %v", err) }

    ctx := context.Background()
    log.Println("argus reset: dropping all tables...")
    // ... explicit DROP TABLE IF EXISTS for every table ...
    log.Println("argus reset: applying Solvent migrations...")
    if err := solventmigrations.Apply(ctx, db); err != nil {
        log.Fatalf("solvent migrations: %v", err)
    }
    log.Println("argus reset: applying ARGUS migrations...")
    if err := migrations.Apply(ctx, db); err != nil {
        log.Fatalf("argus migrations: %v", err)
    }
    log.Println("argus reset: done")
}
```

`main()` switch gains `"migrate"`:
```go
case "migrate":
    cmdMigrate(os.Args[2:])
```

### Step F: Replace test helpers

Replace filesystem-based test helpers with the package:
- `integration_test.go:integApplySolventMigrations` → `solventmigrations.Apply(ctx, db)`
- `internal/ui/ui_test.go:applySolventMigrations` → `solventmigrations.Apply(ctx, db)`
- `internal/application/app_test.go:applySolventMigrations` → `solventmigrations.Apply(ctx, db)`

**Tests:**
- `TestCmdMigrateIsIdempotent` — run twice, no errors, no data loss
- `TestCmdResetDropsAndReapplies` — verify tables recreated
- `TestCmdMigrateRequiresDB` — no flag and no .argus-runtime.env → error

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

## Fix 4: UNKNOWN ≠ EMPTY in GetContext (per-section availability)

**Files:**
- `internal/application/app.go:298-322` — `GetContext`
- `internal/epistemic/view.go` — `Snapshot`, `Context` types
- `internal/mcp/adapter.go` — `handleGetContext`

**Problem:** `GetContext` returns `nil, error` when a backend is unavailable. The MCP adapter propagates this as a tool-call error. An agent cannot distinguish "no research exists" from "backend is down." Additionally, a single shared `Reason` field means one failure overwrites the other.

**Fix:** Per-section availability metadata:

1. Add to `internal/epistemic/view.go`:
```go
// SectionAvailability tracks whether a specific data source is reachable.
type SectionAvailability struct {
    Available bool   `json:"available"`
    Reason    string `json:"reason,omitempty"`
}

// RCPAvailability reports the status of each data source independently.
// An agent seeing task.available=false, snapshot.available=true knows
// that epistemic data is trustworthy but task data is not.
type RCPAvailability struct {
    Task           SectionAvailability `json:"task"`
    Dependencies   SectionAvailability `json:"dependencies"`
    Snapshot       SectionAvailability `json:"snapshot"`
}
```

2. Add to `Context` struct (in `internal/application/app.go`):
```go
type Context struct {
    Task         *work.Task                  `json:"task"`
    Dependencies []*work.Dependency          `json:"dependencies"`
    Snapshot     *epistemic.Snapshot         `json:"snapshot"`
    Availability epistemic.RCPAvailability   `json:"availability"`
}
```

3. In `GetContext`, catch errors per sub-call independently:
```go
func (a *App) GetContext(ctx context.Context, taskID string) (*Context, error) {
    result := &Context{
        Availability: epistemic.RCPAvailability{
            Task:         epistemic.SectionAvailability{Available: true},
            Dependencies: epistemic.SectionAvailability{Available: true},
            Snapshot:     epistemic.SectionAvailability{Available: true},
        },
        Snapshot: &epistemic.Snapshot{},
    }

    task, err := a.workStore.GetByID(ctx, taskID)
    if err != nil {
        result.Availability.Task.Available = false
        result.Availability.Task.Reason = "work_unavailable"
    } else {
        result.Task = task
    }

    deps, err := a.listDependencies(ctx, taskID)
    if err != nil {
        result.Availability.Dependencies.Available = false
        result.Availability.Dependencies.Reason = "dependencies_unavailable"
    } else {
        result.Dependencies = deps
    }

    if result.Task != nil {
        snap, err := epistemic.GetSnapshot(ctx, a.db, result.Task.ProjectID, epistemic.SnapshotOpts{
            IncludeEvidence: true,
        })
        if err != nil {
            result.Availability.Snapshot.Available = false
            result.Availability.Snapshot.Reason = "solvent_unavailable"
        } else {
            result.Snapshot = snap
        }
    } else {
        result.Availability.Snapshot.Available = false
        result.Availability.Snapshot.Reason = "skipped: task unavailable"
    }

    return result, nil // never return error for availability issues
}
```

**Tests:**
- `TestGetContextWithUnavailableWorkStore` — work failure → `task.available=false`, `snapshot.available=false`
- `TestGetContextWithUnavailableSolvent` — snapshot failure → `task.available=true`, `snapshot.available=false`
- `TestGetContextWithEmptyScenario` — valid DB, no beliefs → all sections available, empty snapshot
- `TestGetContextNormal` — valid DB, beliefs exist → full Context, all available
- `TestGetContextAvailabilityIndependence` — work succeeds but snapshot fails → only snapshot unavailable

---

## Fix 5: Retirement-rule enforcement before Discharge (BLOCKER FIX)

**Files:**
- `domain-pack/registry.go` — extend `Pack` interface, add decoder registration, generic types
- `domain-pack/bmist/v1/types.go` — use generic types, implement Validator
- `coordinator/validate.go` — simplify type assertions
- `internal/application/app.go` — `SubmitDecision`, `App` struct, `DecisionRequest`
- `cmd/argus/main.go` — composition root: register decoder, load packs

### Type architecture

The `domain-pack/registry.go` package is the generic layer. It must NOT import `bmistv1`. Instead, the composition root (`cmd/argus/main.go`) registers a decoder function for the `"bmist"` pack ID. The generic registry calls the decoder when loading pack JSON.

**Import cycle prevention:**
```
domain-pack (generic) ← does NOT import bmistv1
bmistv1 (domain) → imports domain-pack (for generic types)
cmd/argus (composition root) → imports both domain-pack and bmistv1
```

### Step A: Add generic types to `domain-pack/registry.go`

```go
// RetirementRule describes the evidence-class gate for retiring a debt item.
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

### Step B: Add decoder registration to `PackRegistry`

```go
// DecoderFunc parses raw JSON into a Pack. Registered per pack_id.
type DecoderFunc func(data []byte) (Pack, error)

type PackRegistry struct {
    mu      sync.RWMutex
    packs   map[string]Pack         // key = "pack_id@version"
    decoders map[string]DecoderFunc // key = pack_id
}

func NewRegistry() *PackRegistry {
    return &PackRegistry{
        packs:    make(map[string]Pack),
        decoders: make(map[string]DecoderFunc),
    }
}

// RegisterDecoder registers a domain-specific decoder for a pack_id.
// The composition root calls this before LoadFromDisk.
func (r *PackRegistry) RegisterDecoder(packID string, fn DecoderFunc) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.decoders[packID] = fn
}
```

### Step C: Fix `decodePack` to use registered decoders

```go
func (r *PackRegistry) decodePack(data []byte, expectedPackID string) (Pack, error) {
    var probe struct {
        PackID  string `json:"pack_id"`
        Version string `json:"version"`
    }
    if err := json.Unmarshal(data, &probe); err != nil {
        return nil, err
    }

    if expectedPackID != "" && !strings.EqualFold(probe.PackID, expectedPackID) {
        return nil, fmt.Errorf("pack_id mismatch: expected %q, got %q", expectedPackID, probe.PackID)
    }

    // Use registered decoder if available
    r.mu.RLock()
    decoder, ok := r.decoders[probe.PackID]
    r.mu.RUnlock()
    if ok {
        return decoder(data)
    }

    // Fallback to generic pack
    return &genericPack{PackID: probe.PackID, Version: probe.Version}, nil
}
```

Note: `decodePack` becomes a method on `PackRegistry` (not a standalone function) to access the decoder map.

### Step D: Extend `Pack` interface

```go
type Pack interface {
    GetPackID() string
    GetVersion() string
    GetDebtVocabulary() []string
    GetEvidenceClasses() []string
    GetRetirementRules() map[string]RetirementRule
    GetVerifierSpecs() []VerifierSpec
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
func (g *genericPack) GetRetirementRules() map[string]RetirementRule { return nil }
func (g *genericPack) GetVerifierSpecs() []VerifierSpec { return nil }
```

### Step F: Update `bmistv1.Pack` to use generic types

**File:** `domain-pack/bmist/v1/types.go`

```go
import "github.com/PithomLabs/oracle/domain-pack"

type Pack struct {
    PackID                string                              `json:"pack_id"`
    Version               string                              `json:"version"`
    Name                  string                              `json:"name"`
    Description           string                              `json:"description"`
    ClaimTypes            []string                            `json:"claim_types"`
    EvidenceClasses       []string                            `json:"evidence_classes"`
    DebtVocabulary        []string                            `json:"debt_vocabulary"`
    InitialDebt           []string                            `json:"initial_debt"`
    RetirementRules       map[string]domainpack.RetirementRule `json:"retirement_rules"`
    Falsifiers            []string                            `json:"falsifiers"`
    ConsequentialActions  []ConsequentialAction               `json:"consequential_actions"`
    HumanGatedTransitions []string                            `json:"human_gated_transitions"`
    VerifierSpecs         []domainpack.VerifierSpec           `json:"verifier_specs"`
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

Remove the old `bmistv1.RetirementRule` and `bmistv1.VerifierSpec` types and their methods. `ConsequentialAction` stays domain-specific (not consumed by the application layer).

### Step G: Make `bmistv1.Pack` implement `Validator`

The existing `Validate` is a free function. Wrap it as a method:

```go
func (p *Pack) Validate() error {
    return validatePack(p)
}
```

Rename the existing `Validate(p *Pack) error` to `validatePack(p *Pack) error` (unexported).

### Step H: Register BM-IST decoder in composition root

**File:** `cmd/argus/main.go`

```go
import (
    domainpack "github.com/PithomLabs/oracle/domain-pack"
    bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
)

func cmdServe(args []string) {
    // ...
    packReg := domainpack.NewRegistry()
    packReg.RegisterDecoder("bmist", func(data []byte) (domainpack.Pack, error) {
        return bmistv1.ParsePack(data)
    })
    if err := packReg.LoadFromDisk("domain-pack"); err != nil {
        log.Printf("warning: pack load failed: %v", err)
    }
    app := application.NewWithRegistry(db, registry, specs)
    app.SetPackRegistry(packReg)
    // ...
}
```

`cmdMCP` does the same. Generic `domain-pack` never imports `bmistv1`.

### Step I: Simplify coordinator validation

**File:** `coordinator/validate.go`

```go
func packHasDebt(pack domainpack.Pack, debt string) bool {
    for _, d := range pack.GetDebtVocabulary() {
        if d == debt { return true }
    }
    return false
}

func getRetirementRule(debtItem string, pack domainpack.Pack) *domainpack.RetirementRule {
    rules := pack.GetRetirementRules()
    if rules == nil { return nil }
    rule, exists := rules[debtItem]
    if !exists { return nil }
    return &rule
}

func ValidateRetirementRule(debtItem string, evidenceClass string, pack domainpack.Pack) error {
    rule := getRetirementRule(debtItem, pack)
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

Remove `RetirementRuleInfo` struct, `RetirementRuleProvider` interface, and all double type-assertions.

### Step J: Simplify `GetPackRules`

**File:** `coordinator/coordinator.go:110-138`

```go
func (c *Coordinator) GetPackRules() map[string]interface{} {
    result := make(map[string]interface{})
    if c.packRegistry == nil { return result }
    pack, err := c.packRegistry.Get("bmist", "1.0.0")
    if err != nil { return result }
    rules := pack.GetRetirementRules()
    if rules == nil { return result }
    result["rules"] = rules // domainpack.RetirementRule marshals to JSON directly
    return result
}
```

No more JSON round-trip conversion.

### Step K: Human decision — pack resolved from authoritative scenario metadata

**Authority model:** `PackRef` is **never** trusted from external request bodies. For human decisions, the pack is resolved server-side from scenario metadata. The `scenarioPackMapping()` function is the authoritative source of truth for scenario→pack mapping.

**PackRef canonical format:** `pack_id@version` (unambiguous — `@` separator, never `-`). The existing `domain-pack/registry.go` already uses `@` as the registry key separator (`registryKey` at line 195). The `@` format is unambiguous for pack IDs containing hyphens (e.g., `market-risk@1.0.0`).

**File:** `internal/application/app.go`

Remove `PackRef` from external request DTOs. The application's `resolvePack` derives from scenario metadata:

```go
// resolvePackFromScenario derives the pack for a scenario.
// Uses the same authoritative mapping as the coordinator.
// The scenario→pack relationship is established at packet ingestion time
// and is the ONLY legitimate source for retirement-rule resolution.
func (a *App) resolvePackFromScenario(ctx context.Context, scenarioID string) (domainpack.Pack, error) {
    if a.packRegistry == nil {
        return nil, fmt.Errorf("pack registry unavailable")
    }
    // Authoritative scenario→pack mapping (POC: hardcoded, same as coordinator).
    packID, packVersion := scenarioPackMapping(scenarioID)
    return a.packRegistry.Get(packID, packVersion)
}

func scenarioPackMapping(scenarioID string) (string, string) {
    return "bmist", "1.0.0"
}

func parsePackRef(ref string) (string, string) {
    // Canonical format: pack_id@version (e.g., "bmist@1.0.0")
    parts := strings.SplitN(ref, "@", 2)
    if len(parts) == 2 {
        return parts[0], parts[1]
    }
    return ref, ""
}
```

**Design note:** The `belief` table currently has no `pack_ref` column. For the POC, the pack is derived from the scenario via `scenarioPackMapping()` (which the coordinator already uses at `coordinator/human.go:193`). The application layer calls the coordinator's pattern. If the POC needs runtime pack resolution, add a `pack_ref` column to `belief` in a follow-up migration.

### Step L: Full retirement-rule enforcement in SubmitDecision (authenticated path)

**Separation of DTOs:**

```go
// ExternalDecisionRequest is the externally-supplied request from the Trust UI.
// It does NOT carry PrincipalID — the principal is derived at the auth boundary.
type ExternalDecisionRequest struct {
    Type          string `json:"type"`
    ScenarioID    string `json:"scenario_id"`
    BeliefID      string `json:"belief_id"`
    ObligationKey string `json:"obligation_key,omitempty"`
    InstrumentRef string `json:"instrument_ref,omitempty"`
    EvidenceClass string `json:"evidence_class,omitempty"`
}

// AuthenticatedDecisionCommand is the internally constructed command
// after authentication. PrincipalID is server-derived, never from request body.
type AuthenticatedDecisionCommand struct {
    Type          string
    ScenarioID    string
    BeliefID      string
    ObligationKey string
    InstrumentRef string
    PrincipalID   string // server-derived from token validation
    EvidenceClass string
}
```

**SubmitDecision takes the internal command:**

```go
func (a *App) SubmitDecision(ctx context.Context, cmd *AuthenticatedDecisionCommand) error {
    switch cmd.Type {
    case "discharge":
        if a.packRegistry == nil {
            return fmt.Errorf("pack registry unavailable: cannot validate retirement")
        }

        // Resolve pack from authoritative scenario metadata — NEVER from request body
        pack, err := a.resolvePackFromScenario(ctx, cmd.ScenarioID)
        if err != nil {
            return fmt.Errorf("pack resolution failed: %w", err)
        }

        rules := pack.GetRetirementRules()
        if rules == nil {
            return fmt.Errorf("pack does not support retirement rules")
        }
        rule, exists := rules[cmd.ObligationKey]
        if !exists {
            return fmt.Errorf("no retirement rule for debt item: %s", cmd.ObligationKey)
        }

        if cmd.EvidenceClass == "" {
            return fmt.Errorf("evidence_class is required for debt discharge of: %s", cmd.ObligationKey)
        }
        if cmd.EvidenceClass != rule.EvidenceClass {
            return fmt.Errorf("evidence class mismatch: debt %q requires %q, got %q",
                cmd.ObligationKey, rule.EvidenceClass, cmd.EvidenceClass)
        }

        evidenceIDs, err := a.verifyPersistedEvidence(ctx, cmd.BeliefID, cmd.ScenarioID, cmd.EvidenceClass)
        if err != nil {
            return fmt.Errorf("evidence verification failed: %w", err)
        }
        if len(evidenceIDs) == 0 {
            return fmt.Errorf("no qualifying evidence of class %q for belief %s in scenario %s",
                cmd.EvidenceClass, cmd.BeliefID, cmd.ScenarioID)
        }

        instrumentRef := buildInstrumentRef(evidenceIDs)
        return a.kern.Discharge(ctx, cmd.ScenarioID, cmd.BeliefID, cmd.ObligationKey, instrumentRef, cmd.PrincipalID)

    case "promote":
        return a.kern.Promote(ctx, cmd.ScenarioID, cmd.BeliefID)

    case "retract":
        retracted, err := a.kern.RetractCascade(ctx, cmd.ScenarioID, cmd.BeliefID)
        if err != nil { return err }
        _ = retracted
        return a.cancelLinkedTasks(ctx, cmd.BeliefID)

    default:
        return fmt.Errorf("unknown decision type: %s", cmd.Type)
    }
}
```

### Step M: Query persisted evidence

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

### Step N: Trust UI constructs AuthenticatedDecisionCommand

**File:** `internal/ui/ui.go`

The UI handler receives `ExternalDecisionRequest`, authenticates, then constructs `AuthenticatedDecisionCommand` with server-derived principal. The UI does NOT send PackRef — it is resolved server-side:

```go
func (s *Server) HandleDischarge(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var extReq application.ExternalDecisionRequest
    if err := json.NewDecoder(r.Body).Decode(&extReq); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    // Server-derived principal — ignore any principal in request body
    principalID := "operator" // derived from token validation (Fix 6)

    cmd := &application.AuthenticatedDecisionCommand{
        Type:          "discharge",
        ScenarioID:    extReq.ScenarioID,
        BeliefID:      extReq.BeliefID,
        ObligationKey: extReq.ObligationKey,
        InstrumentRef: extReq.InstrumentRef,
        PrincipalID:   principalID,
        EvidenceClass: extReq.EvidenceClass,
    }

    err := s.app.SubmitDecision(r.Context(), cmd)
    if err != nil {
        http.Error(w, fmt.Sprintf("discharge failed: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintln(w, `{"status":"ok"}`)
}
```

**Tests:**
- `TestDischargeRequiresPackResolution`
- `TestDischargeRequiresKnownDebtItem`
- `TestDischargeRequiresEvidenceClass`
- `TestDischargeWithMismatchedEvidenceClass`
- `TestDischargeRequiresPersistedEvidence`
- `TestDischargeWithQualifyingEvidence`
- `TestDischargeWithRegistryUnavailable`
- `TestSubmitDecisionRejectsExternalPrincipalID`
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

4. `requireOrigin` — explicit Origin/Host validation for consequential POSTs:
```go
func (s *Server) requireOrigin(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        if origin != "" {
            u, err := url.Parse(origin)
            if err != nil || (u.Host != "localhost" && u.Host != "127.0.0.1") {
                http.Error(w, "forbidden: invalid origin", http.StatusForbidden)
                return
            }
        }
        host := r.Host
        if host != "" {
            hostname := strings.Split(host, ":")[0]
            if hostname != "localhost" && hostname != "127.0.0.1" {
                http.Error(w, "forbidden: invalid host", http.StatusForbidden)
                return
            }
        }
        next(w, r)
    }
}
```

5. Protect write endpoints with BOTH checks:
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

6. Server-derived principal from token validation (never from request body).

7. Wire token — **FAIL CLOSED**:
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
- `TestDischargeRejectsInvalidOrigin`
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

## Fix 9: Server-verify evidence hashes + operator_asserted authority

**File:** `internal/application/app.go:198-214`

**Problem:** `content_sha256` is client-supplied and never verified. Additionally, `operator_asserted` evidence can be submitted by agents without human attestation.

**Fix:** Enforce evidence-class authority at the packet submission boundary:

### Evidence-class enforcement in Persist

```go
for i, e := range pkt.Evidence {
    // operator_asserted requires human attestation — agents cannot create it
    if e.ProvenanceClass == "operator_asserted" {
        return nil, fmt.Errorf("evidence[%d]: operator_asserted evidence requires human attestation; agent-submitted packets cannot create it", i)
    }

    // reproducible_artifact requires trusted verifier artifact
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

    // external_feed, live_scan: accept agent-supplied hash (no registry lookup)
}
```

### How operator_asserted evidence enters the system

The `operator_asserted` evidence path is:
1. Human types artifact text in Solvent wizard
2. Wizard computes SHA256 server-side, creates evidence with `provenance_class = "operator_asserted"`
3. Trust UI references this existing evidence when discharging debt

Agents never create `operator_asserted` evidence. The `Persist` method rejects it.

### Tests
- `TestReproducibleArtifactRequiresArtifactRef`
- `TestReproducibleArtifactRequiresTrustedArtifact`
- `TestReproducibleArtifactHashMismatch`
- `TestReproducibleArtifactVerifierNotAuthorized`
- `TestReproducibleArtifactHappyPath`
- `TestOperatorAssertedRejectedFromAgent`
- `TestOperatorAssertedAcceptedFromWizard` (via direct SQL insert)

---

## Fix 10: Port preflight and lifecycle (BLOCKER FIX)

**Files:**
- `cmd/argus/main.go`
- `Taskfile.yml`
- Test files

**Problem:** No port preflight, no fallback, no cleanup. `:8080` CRDB admin/ARGUS collision. Port mismatch between Taskfile (`26260`) and tests (`26257`). No persistent runtime configuration.

### Architecture

```
argus resolve-ports → writes .argus-runtime.env (FAILS if no port available)
    ↓
task dev / task test / argus mcp / integration tests
    ↓
source .argus-runtime.env → all commands use resolved ports
    ↓
CRDB startup retries on EADDRINUSE → re-resolves ports → FAILS after max attempts
    ↓
ARGUS startup retries on EADDRINUSE → re-resolves ports → FAILS after max attempts
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

    crdbSQL, err := resolvePort(*crdbSQLPreferred)
    if err != nil { log.Fatalf("resolve CRDB SQL port: %v", err) }

    crdbAdmin, err := resolvePortDifferent(*crdbAdminPreferred, crdbSQL)
    if err != nil { log.Fatalf("resolve CRDB admin port: %v", err) }

    argus, err := resolvePortDifferent(*argusPreferred, crdbSQL, crdbAdmin)
    if err != nil { log.Fatalf("resolve ARGUS port: %v", err) }

    env := fmt.Sprintf("CRDB_SQL_PORT=%d\nCRDB_ADMIN_PORT=%d\nARGUS_PORT=%d\nARGUS_DB_URL=postgres://root@localhost:%d/argus?sslmode=disable\n",
        crdbSQL, crdbAdmin, argus, crdbSQL)
    if err := os.WriteFile(".argus-runtime.env", []byte(env), 0644); err != nil {
        log.Fatalf("write runtime env: %v", err)
    }
    fmt.Print(env)
}
```

### Port resolution functions — FAIL CLOSED

```go
func resolvePort(preferred int) (int, error) {
    addr := fmt.Sprintf(":%d", preferred)
    ln, err := net.ListenTimeout("tcp", addr, 500*time.Millisecond)
    if err == nil {
        ln.Close()
        return preferred, nil
    }
    for p := preferred + 1; p <= preferred+10; p++ {
        ln, err := net.ListenTimeout("tcp", fmt.Sprintf(":%d", p), 500*time.Millisecond)
        if err == nil {
            ln.Close()
            return p, nil
        }
    }
    return 0, fmt.Errorf("no available port near %d (all %d-%d occupied)", preferred, preferred, preferred+10)
}

func resolvePortDifferent(preferred int, occupied ...int) (int, error) {
    occupiedSet := make(map[int]bool, len(occupied))
    for _, p := range occupied { occupiedSet[p] = true }
    for {
        port, err := resolvePort(preferred)
        if err != nil { return 0, err }
        if !occupiedSet[port] { return port, nil }
        preferred = port + 1
    }
}
```

### CRDB + ARGUS startup with EADDRINUSE retry — FAIL CLOSED

Both CRDB and ARGUS share the same retry pattern. The `argus serve` command itself has built-in EADDRINUSE retry:

```go
func cmdServe(args []string) {
    fs := flag.NewFlagSet("serve", flag.ExitOnError)
    dbURL := fs.String("db", loadRuntimeDBURL(), "database URL")
    listen := fs.String("listen", ":8080", "HTTP listen address")
    fs.Parse(args)
    if *dbURL == "" { log.Fatal("serve: --db flag or .argus-runtime.env required") }

    for attempt := 1; attempt <= 3; attempt++ {
        ln, err := net.Listen("tcp", *listen)
        if err == nil {
            ln.Close()
            break // port available
        }
        if !strings.Contains(err.Error(), "address already in use") || attempt == 3 {
            log.Fatalf("serve: listen on %s: %v", *listen, err)
        }
        log.Printf("serve: port %s in use (attempt %d/3), re-resolving...", *listen, attempt)
        if newPort := resolveArgusPort(); newPort > 0 {
            *listen = fmt.Sprintf(":%d", newPort)
        }
        time.Sleep(2 * time.Second)
    }
    // ... rest of cmdServe
}

func resolveArgusPort() int {
    data, err := os.ReadFile(".argus-runtime.env")
    if err != nil { return 0 }
    for _, line := range strings.Split(string(data), "\n") {
        if strings.HasPrefix(line, "ARGUS_PORT=") {
            var port int
            fmt.Sscanf(strings.TrimPrefix(line, "ARGUS_PORT="), "%d", &port)
            return port
        }
    }
    return 0
}
```

The Taskfile applies the same pattern for CRDB, with **concurrent dev lock**:

```yaml
  dev:
    desc: Start ARGUS in development mode (non-destructive)
    cmds:
      - mkdir -p .argus-pids .cockroach-data
      - |
        # Acquire startup lock — only one dev stack per workspace
        LOCKFILE=".task/dev.lock"
        if [ -f "$LOCKFILE" ]; then
          LOCK_PID=$(cat "$LOCKFILE")
          if kill -0 "$LOCK_PID" 2>/dev/null; then
            echo "FATAL: another task dev is running (PID $LOCK_PID). Use 'task down' first."
            exit 1
          else
            echo "removing stale lock (PID $LOCK_PID no longer running)"
            rm -f "$LOCKFILE"
          fi
        fi
        echo $$ > "$LOCKFILE"
        trap 'rm -f "$LOCKFILE"; for f in .argus-pids/*.pid; do [ -f "$f" ] && kill $(cat "$f") 2>/dev/null; rm -f "$f"; done' EXIT INT TERM
        
        # Resolve ports and write runtime config
        go run ./cmd/argus resolve-ports || { echo "FATAL: cannot resolve ports"; exit 1; }
        source .argus-runtime.env
        
        # Start CRDB with retry on EADDRINUSE (max 3 attempts)
        CRDB_STARTED=false
        for attempt in 1 2 3; do
          cockroach start-single-node --insecure \
            --listen-addr :$CRDB_SQL_PORT \
            --http-addr :$CRDB_ADMIN_PORT \
            --store=.cockroach-data 2>/dev/null &
          CRDB_PID=$!
          echo $CRDB_PID > .argus-pids/crdb.pid
          
          READY=false
          for i in $(seq 1 30); do
            if cockroach node status --host :$CRDB_SQL_PORT --insecure 2>/dev/null; then
              READY=true
              break
            fi
            if ! kill -0 $CRDB_PID 2>/dev/null; then
              break
            fi
            sleep 1
          done
          
          if [ "$READY" = "true" ]; then
            CRDB_STARTED=true
            break
          fi
          
          echo "CRDB not ready (attempt $attempt/3), re-resolving ports..."
          kill $CRDB_PID 2>/dev/null; wait $CRDB_PID 2>/dev/null
          rm -f .argus-pids/crdb.pid
          go run ./cmd/argus resolve-ports || { echo "FATAL: cannot resolve ports after retry"; exit 1; }
          source .argus-runtime.env
        done
        
        if [ "$CRDB_STARTED" != "true" ]; then
          echo "FATAL: CRDB failed to start after 3 attempts"
          exit 1
        fi
        
        # Non-destructive migration (idempotent, safe on every startup)
        go run ./cmd/argus migrate --db "$ARGUS_DB_URL"
        
        # Start ARGUS (has its own EADDRINUSE retry built-in)
        go run ./cmd/argus serve --db "$ARGUS_DB_URL" --listen ":$ARGUS_PORT" &
        echo $! > .argus-pids/argus.pid
        
        wait
```

**Lock semantics:** The lock file `.task/dev.lock` contains the PID of the owning `task dev` process. If the PID is still running, the second invocation refuses to start. If the PID is stale (process crashed or was killed), the lock is removed and the new invocation proceeds. The `trap EXIT` ensures the lock is always cleaned up.

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

### Remove hardcoded CRDB port from all commands

`cmdMCP`, `cmdMigrate`, `cmdReset`, `cmdVerify` all default to reading from `.argus-runtime.env`:

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

### Test DSN derivation

**File:** `integration_test.go`, `app_test.go`, `ui_test.go`

Parse admin DSN to extract host:port for test database construction:

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
    u, _ := url.Parse(dsn)
    hostPort := u.Host // e.g., "localhost:26257"
    // ... create test DB using hostPort ...
    testDSN := fmt.Sprintf("postgres://%s/%s?sslmode=disable", hostPort, dbName)
}
```

### Port allocation tests

- `TestResolvePortPreferredAvailable`
- `TestResolvePortFallback`
- `TestResolvePortExhaustion` — all ports occupied → error
- `TestResolvePortDifferent`
- `TestCRDBAdminARGUSCollision` — CRDB admin on 8080, ARGUS tries 8080 → resolves to 8081
- `TestServeRetriesOnEADDRINUSE` — mock listener → retry then succeed
- `TestServeFailsAfterMaxAttempts` — all attempts fail → fatal
- `TestDevLockPreventsConcurrentStartup` — two `task dev` invocations → second refuses
- `TestDevLockStalePIDRecovery` — stale lock file → removed and new invocation proceeds

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
4. Fix 9 (evidence hashes + operator_asserted) — no dependencies
5. Fix 2 (Solvent migrations) — needs Solvent package creation
6. Fix 4 (UNKNOWN ≠ EMPTY) — needs view.go changes
7. Fix 5 (retirement rules) — needs registry/interface/bmist type changes + decoder registration
8. Fix 6 (UI auth) — needs token wiring + Origin/Host check
9. Fix 7 (MCP SDK) — needs `go get` + adapter rewrite
10. Fix 10 (port preflight) — needs port package + runtime env + CRDB retry
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
- **`domain-pack` does NOT import `bmistv1`.** Decoder registration from composition root.
- **`internal/application` does NOT hardcode "bmist".** Pack resolved from authoritative scenario metadata via `scenarioPackMapping()`.
- **`operator_asserted` evidence rejected from agent packets.** Requires human attestation.
- **`RetirementRule` is a plain struct in `domainpack`.** No methods, no interface, no type assertions.
- **Port resolution FAILS on exhaustion.** Never returns occupied port.
- **Migration FS errors handled explicitly.** `Apply()` fails closed.
- **`.argus-runtime.env` is the persistent runtime config.** All commands consume it.
- **Origin/Host validation on consequential UI POSTs.** In addition to token/cookie auth.
- **`PackRef` is NEVER trusted from external request bodies.** Server-derived from scenario metadata.
- **`PrincipalID` is NEVER externally representable.** Split external DTO from authenticated command.
- **`task dev` is NON-DESTRUCTIVE.** Uses `argus migrate` (idempotent). `task fresh` uses destructive `argus reset`.
- **Both CRDB and ARGUS have EADDRINUSE retry.** `cmdServe` retries listen with re-resolution; Taskfile retries CRDB startup.
- **Concurrent `task dev` prevented.** `.task/dev.lock` with PID-based ownership check.
- **PackRef canonical format: `pack_id@version`.** The `@` separator is unambiguous for pack IDs containing hyphens.
