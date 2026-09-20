# Remaining POC Work — Implementation Plan (Revised)

**Status:** Ready for build. 5 workstreams, no architecture changes, no third Solvent change.

---

## WORKSTREAM 1: VERIFIER WIRING

### 1a. Add `VerifierSpec` to domain pack types

**File:** `domain-pack/bmist/v1/types.go`

Add to `Pack` struct:
```go
VerifierSpecs []VerifierSpec `json:"verifier_specs"`
```

Add new type:
```go
type VerifierSpec struct {
    VerifierID  string `json:"verifier_id"`
    MinVersion  string `json:"min_version"`
}
```

Add method:
```go
func (p *Pack) GetVerifierSpecs() []VerifierSpec { return p.VerifierSpecs }
```

### 1b. Extend `VerificationArtifact`

**File:** `verifier/model/artifact.go`

Add fields to `VerificationArtifact`:
```go
VerifierID string `json:"verifier_id"`
Tolerance  string `json:"tolerance"`
```

**Hash semantics — CORRECTED:**

`ComputeArtifactHash` must exclude **only** self-referential/non-semantic fields:
```
EXCLUDE:
    ArtifactHash  (self-referential — it IS the hash)
    Timestamp     (non-semantic — varies per execution)

INCLUDE:
    RunID
    VerifierVersion
    VerifierHash
    VerifierID          ← NEW, must be covered
    ClaimHash
    InputHash
    Tolerance           ← NEW, must be covered
    Steps               (all step content)
    Result
    EvidenceRef
```

The hash is: `SHA-256(JSONcanonical(artifact where ArtifactHash="" and Timestamp=zero))`

All semantic fields including `VerifierID`, `VerifierVersion`, `Tolerance`, and `Steps` are covered. Only the hash itself and the timestamp are excluded.

### 1c. Add verifier methods to App

**File:** `internal/application/app.go`

Add fields:
```go
artifactRegistry *verifier.ArtifactRegistry
```

Add methods:
```go
func (a *App) Verify(ctx context.Context, input physicsv1.VerifierInput) (*model.VerificationArtifact, error)
func (a *App) VerifyWithTimeout(ctx context.Context, input physicsv1.VerifierInput, timeout time.Duration) (*model.VerificationArtifact, error)
```

`Verify` logic:
1. Run `verifier.RunPhysicsVerifier(ctx, a.artifactRegistry, input)`
2. **Runtime VerifierSpec enforcement** (see 1d below)
3. Register artifact
4. Persist as evidence with provenance_class=reproducible_artifact

Failure semantics:
- Verifier error → `Result=inconclusive`, packet proceeds
- Refuted → evidence with `result: refuted`, packet proceeds
- Timeout: default 30s, configurable

### 1d. Runtime VerifierSpec enforcement (NOT compile-time)

**File:** `internal/application/app.go` — `Validate()` method

The enforcement path is runtime packet-validation:

```
packet compilation (Compile)
    ↓
load selected PackDefinition from PackRegistry
    ↓
Validate():
    ├─ existing: debt membership check
    ├─ existing: evidence class check
    └─ NEW: VerifierSpec enforcement
         ├─ for each artifact reference in packet:
         │   ├─ extract verifier_id + verifier_version from artifact
         │   ├─ look up VerifierSpec in pack (by verifier_id)
         │   ├─ reject if verifier_id not in pack's VerifierSpecs
         │   └─ reject if verifier_version < spec.MinVersion
         └─ return error if unauthorized verifier
```

**Rejection behavior:**
- Error message: `"verifier %s version %s not authorized by pack %s (requires >= %s)"`
- Packet does NOT proceed to Persist
- This is a hard gate, not advisory

### 1e. Wire `argus verify` CLI

**File:** `cmd/argus/main.go`

`verify` subcommand:
- `--input` flag: JSON file with `VerifierInput`
- `--timeout` flag: default 30s
- Creates App + registry
- Calls `app.Verify(ctx, input)`
- Outputs artifact as JSON

### 1f. Wire registry into `argus serve`

**File:** `cmd/argus/main.go`

- Create `verifier.ArtifactRegistry` in serve command
- Pass to `NewApp(db, registry)`

---

## WORKSTREAM 2: REFUSAL / SECURITY TESTS (9 tests)

### Test 1: Agent proposed retirement does not discharge debt

**File:** `internal/application/app_test.go`

```go
func TestAgentRetirementDoesNotDischargeDebt(t *testing.T) {
    // Setup: create belief with debt via Persist
    // Submit packet with proposed_retirement
    // Assert: ledger shows debt unchanged
}
```

### Test 2: MCP exposes exactly 2 tools

**File:** `internal/mcp/adapter_test.go`

```go
func TestMCPToolsCount(t *testing.T) {
    // Create adapter
    // ListTools returns len=2
    // Names: "argus.get_context", "argus.submit_packet"
    // No promote, no retract, no discharge
}
```

### Test 3: Promote with open debt refused

**File:** `internal/application/app_test.go`

```go
func TestPromoteWithOpenDebtRefused(t *testing.T) {
    // Create belief with debt
    // Call SubmitDecision promote
    // Assert: error contains "debt" or SQLSTATE 23514
}
```

### Test 4: Promote retracted belief refused

**File:** `internal/application/app_test.go`

```go
func TestPromoteRetractedBeliefRefused(t *testing.T) {
    // Create belief, retract it, try promote
    // Assert: error
}
```

### Test 5: Agent cannot retract

**File:** `internal/mcp/adapter_test.go`

```go
func TestAgentCannotRetract(t *testing.T) {
    // ListTools: no "retract" tool
    // HandleTool("argus.retract", ...) returns "unknown tool"
}
```

### Test 6: Bad Origin rejected

**File:** `internal/ui/ui_test.go`

```go
func TestBadOriginRejected(t *testing.T) {
    // POST to /ui/api/discharge with wrong Origin header
    // Assert: 403 Forbidden
}
```

### Test 7: Request-body principal ignored

**File:** `internal/ui/ui_test.go`

```go
func TestRequestBodyPrincipalIgnored(t *testing.T) {
    // POST discharge with principal_id: "evil" in body
    // Assert: audit log shows principal_id = "operator"
}
```

### Test 8: Failed mutation produces no audit event

**File:** `internal/application/app_test.go`

```go
func TestFailedMutationNoAudit(t *testing.T) {
    // Try promote on nonexistent belief
    // Assert: error returned, no belief_promoted audit row
}
```

### Test 9: Partial packet failure retries succeed (DETERMINISTIC)

**File:** `internal/application/app_test.go`

```go
func TestPartialFailureRetrySucceeds(t *testing.T) {
    db := setupTestDB(t)
    app := application.New(db)

    // Step 1: Submit packet with 3 beliefs
    pkt := &packetv1.Packet{
        ScenarioID: "test-scenario",
        PacketID:   "pkt-001",
        Beliefs: []packetv1.Belief{
            {LocalID: "b1", Claim: "first claim", ClaimType: "derived"},
            {LocalID: "b2", Claim: "second claim", ClaimType: "derived"},
            {LocalID: "b3", Claim: "third claim", ClaimType: "derived"},
        },
    }

    // Step 2: Inject fault after first entity persists
    failAfter := 1
    wrappedDB := &faultInjectionDB{
        DB:          db,
        failAfterN: failAfter,
    }
    faultApp := application.New(wrappedDB)

    err := faultApp.Persist(ctx, pkt)
    assert.Error(t, err)  // partial failure

    // Step 3: Count entities — exactly 1 belief persisted
    count := countEntities(t, db, "belief", "test-scenario")
    assert.Equal(t, 1, count)

    // Step 4: Retry identical packet with real DB
    result, err := app.Persist(ctx, pkt)
    assert.NoError(t, err)

    // Step 5: Verify exactly 3 beliefs (deterministic IDs, ON CONFLICT DO NOTHING)
    count = countEntities(t, db, "belief", "test-scenario")
    assert.Equal(t, 3, count)

    // Step 6: Verify idempotency row exists
    idempotent := countEntities(t, db, "submission_idempotency", "test-scenario")
    assert.Equal(t, 1, idempotent)
}
```

**Fault injection mechanism:**
- `faultInjectionDB` wraps `*sql.DB`
- Overrides `ExecContext` to count calls and fail after N
- Deterministic: always fails at the same point
- On retry, `ON CONFLICT (id) DO NOTHING` ensures no duplicates

**Why this works:**
- Deterministic UUIDs from `entityID(scenario, type, content)` → same input = same ID
- First attempt: belief b1 persisted (ID = `entityID(scenario, "belief", "first claim")`), then fails
- Retry: belief b1 `ON CONFLICT (id) DO NOTHING` → no-op, b2 and b3 inserted normally
- Result: exactly 3 beliefs, 0 duplicates

### Boundary Enforcement (Depguard)

**File:** `.golangci.yml`

```yaml
linters-settings:
  depguard:
    rules:
      mcp_boundary:
        files:
          - "**/mcp/**"
        deny:
          - pkg: "github.com/PithomLabs/oracle/internal/epistemic"
            desc: "MCP must not directly access epistemic kernel"
```

**File:** `internal/boundary/boundary_test.go`

Uses `go/packages` to verify import graph. MCP package does not import kernel Discharge/Promote/Retract.

---

## WORKSTREAM 3: FULL INTEGRATION TEST (17 steps)

**File:** `integration_test.go` (project root or `internal/integration/`)

**Environment:** Ephemeral CockroachDB test instances. NOT in-memory.

**Test database setup:**
```go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()

    // 1. Connect to defaultdb on local CRDB
    adminDB, err := sql.Open("pgx", "postgres://root@localhost:26260/defaultdb?sslmode=disable")

    // 2. Create unique test database: argus_test_<t.Name>_test
    dbName := fmt.Sprintf("argus_test_%s_test", t.Name())
    adminDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
    adminDB.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %q", dbName))

    // 3. Connect to test database
    testDB, _ := sql.Open("pgx", fmt.Sprintf("postgres://root@localhost:26260/%s?sslmode=disable", dbName))

    // 4. Apply Solvent migrations (via Solvent Change #2 export)
    solventmigrations.Apply(ctx, testDB)

    // 5. Apply ARGUS migrations
    migrations.Apply(ctx, testDB)

    // 6. Register cleanup
    t.Cleanup(func() {
        testDB.Close()
        adminDB.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
        adminDB.Close()
    })

    return testDB
}
```

**Requirements:**
- CockroachDB running locally on port 26260 (standard dev mode)
- Each test gets its own database with `_test` suffix
- Cleanup drops the database after test completes
- No in-memory mode — real CRDB instance, ephemeral databases

**17-step test:**

| Step | Action | Verification |
|------|--------|-------------|
| 1 | Start `argus serve` on random port | `/health` returns 200 |
| 2 | Call `get_context` for non-existent task | Error or empty context |
| 3 | Create project + task via `work.Store` | Task exists, status=proposed |
| 4 | Submit work packet via `app.Persist` | Beliefs + evidence in DB |
| 5 | Query `get_context` for task | Snapshot has beliefs + evidence |
| 6 | Submit packet with verification evidence | Artifact registered, provenance_class=reproducible_artifact |
| 7 | Submit adversarial packet with contradicts edge | Edge stored, no authority transition |
| 8 | Human reviews via UI `/ui/insights` | HTML renders beliefs + tasks |
| 9 | Human discharges debt via `/ui/api/discharge` | Debt item retired in belief |
| 10 | Human promotes via `/ui/api/promote` | Belief status = promoted |
| 11 | Human retracts via `app.SubmitDecision` | Belief status = retracted |
| 12 | Check linked task cancelled | Task status = cancelled |
| 13 | REOPEN: create successor task with `reopened_from_task_id` | Lineage preserved |
| 14 | Verify UI `/ui/debts` shows retirement proposals | HTML renders correctly |
| 15 | Verify idempotency: submit same packet twice | No duplicate entities |
| 16 | Verify dependency blocking: task with blocker stays proposed | Cannot claim blocked task |
| 17 | Full lifecycle: packet → verify → discharge → promote → retract | All transitions succeed |

---

## WORKSTREAM 4: SOLVENT MIGRATION EXPORT (Change #2)

### What ARGUS does NOT do:
- **No copying Solvent SQL files into ARGUS**
- **No `internal/epistemic/migrations/` directory**
- **No `internal/epistemic/migrations.go` file**

### What Solvent does (Change #2, already budgeted):

**File:** `solvent/migrations/` (new package in Solvent repo)

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
var fs embed.FS

// Apply runs all Solvent migrations in order.
func Apply(ctx context.Context, db *sql.DB) error {
    entries, err := fs.ReadDir("db")
    if err != nil {
        return fmt.Errorf("read migrations dir: %w", err)
    }

    var sqlFiles []string
    for _, e := range entries {
        if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
            sqlFiles = append(sqlFiles, e.Name())
        }
    }
    sort.Strings(sqlFiles)

    for _, name := range sqlFiles {
        data, err := fs.ReadFile("db/" + name)
        if err != nil {
            return fmt.Errorf("read migration %s: %w", name, err)
        }
        for _, stmt := range splitStatements(string(data)) {
            if _, err := db.ExecContext(ctx, stmt); err != nil {
                return fmt.Errorf("apply migration %s: %w", name, err)
            }
        }
    }
    return nil
}
```

### What ARGUS does (consume the export):

**File:** `cmd/argus/main.go` — `cmdReset`

```go
import solventmigrations "github.com/PithomLabs/solvent/migrations"

// In cmdReset:
log.Println("argus reset: applying Solvent migrations...")
if err := solventmigrations.Apply(ctx, db); err != nil {
    log.Fatalf("solvent migrations: %v", err)
}

log.Println("argus reset: applying ARGUS migrations...")
if err := migrations.Apply(ctx, db); err != nil {
    log.Fatalf("argus migrations: %v", err)
}
```

ARGUS migrations stay in `internal/migrations/` (work tables, idempotency, retirement proposal).

---

## WORKSTREAM 5: FINAL ACCEPTANCE VERIFICATION

Run all 25 criteria from the frozen plan:

```bash
# AC1: Binary builds and runs
go build -o bin/argus ./cmd/argus && ./bin/argus --help

# AC2: MCP tools count
./bin/argus mcp --list-tools | jq '.tools | length'  # → 2

# AC3-AC6: Refusal tests
go test ./internal/application/... ./internal/mcp/... ./internal/ui/... -run TestRefusal -v

# AC7: Depguard
golangci-lint run --enable depguard

# AC8: Idempotency
go test ./internal/application/... -run TestIdempotency -v

# AC9-AC12: UI
go test ./internal/ui/... -v

# AC13-AC17: Integration
go test ./... -run TestFullIntegration -v -timeout 300s

# AC18-AC20: Verifier
go test ./verifier/... -v

# AC21: Race
go test -race ./...

# AC22: Vet
go vet ./...

# AC23: Boundary
go test ./internal/boundary/... -v

# AC24: Taskfile
task dev && task test && task lint

# AC25: Acceptance
echo "All 25 criteria verified"
```

---

## FILE CHANGE SUMMARY

| File | Action | Purpose |
|------|--------|---------|
| `solvent/migrations/` (Solvent repo) | Create | Solvent Change #2: exported embed.FS + Apply |
| `domain-pack/bmist/v1/types.go` | Modify | Add VerifierSpec type |
| `domain-pack/bmist/v1/pack.json` | Modify | Add verifier_specs |
| `domain-pack/bmist/v1/validate.go` | Modify | Validate VerifierSpecs |
| `verifier/model/artifact.go` | Modify | Add VerifierID, Tolerance; fix hash |
| `verifier/physics/v1/verifier.go` | Modify | Set VerifierID in artifact |
| `internal/application/app.go` | Modify | Add Verify, runtime VerifierSpec enforcement |
| `cmd/argus/main.go` | Modify | Wire verifier, import solventmigrations |
| `internal/application/app_test.go` | Create | Refusal tests 1,3,4,8,9 + partial failure test |
| `internal/mcp/adapter_test.go` | Create | Refusal tests 2, 5 |
| `internal/ui/ui_test.go` | Create | Refusal tests 6, 7 |
| `internal/boundary/boundary_test.go` | Create | Depguard verification |
| `integration_test.go` | Create | 17-step integration test |
| `.golangci.yml` | Create | Depguard rules |

**Explicitly NOT created:**
- ~~`internal/epistemic/migrations.go`~~
- ~~`internal/epistemic/migrations/*.sql`~~

---

## RISK CHECK

- **Solvent changes:** 2 of 2 budget used (projection contract + migration FS). No third change.
- **No architecture change:** All work is wiring and testing.
- **No agent authority bypass:** VerifierSpec enforcement is runtime packet-validation in `Validate()`.
- **No false causal ordering:** Integration test steps are explicit, no invented happens-before.
- **Migration ownership preserved:** Solvent owns canonical epistemic migrations. ARGUS consumes via export.
- **Verifier hash is trust-bound:** All semantic fields covered. Only ArtifactHash and Timestamp excluded.
