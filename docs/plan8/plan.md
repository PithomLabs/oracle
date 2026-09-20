# Remaining POC Work — Implementation Plan

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

`ComputeArtifactHash` must exclude `ArtifactHash`, `Timestamp`, `VerifierID`, `Tolerance`.

### 1c. Add verifier methods to App

**File:** `internal/application/app.go`

Add fields:
```go
artifactRegistry *verifier.ArtifactRegistry
verifierSpecs    []bmistv1.VerifierSpec  // injected from pack
```

Add methods:
```go
func (a *App) Verify(ctx context.Context, input verifier.VerifierInput) (*model.VerificationArtifact, error)
func (a *App) VerifyWithTimeout(ctx context.Context, input verifier.VerifierInput, timeout time.Duration) (*model.VerificationArtifact, error)
```

`Verify` logic:
1. Run `verifier.RunPhysicsVerifier(ctx, a.artifactRegistry, input)`
2. Enforce VerifierSpec: check artifact.VerifierID + artifact.VerifierVersion against pack's VerifierSpecs
3. Register artifact
4. Persist as evidence with provenance_class=reproducible_artifact

Failure semantics:
- Verifier error → `Result=inconclusive`, packet proceeds
- Refuted → evidence with `result: refuted`, packet proceeds
- Timeout: default 30s, configurable

### 1d. Wire `argus verify` CLI

**File:** `cmd/argus/main.go`

`verify` subcommand:
- `--input` flag: JSON file with `VerifierInput`
- `--timeout` flag: default 30s
- Creates App + registry
- Calls `app.Verify(ctx, input)`
- Outputs artifact as JSON

### 1e. Wire registry into `argus serve`

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

### Test 9: Partial packet failure retries succeed

**File:** `internal/application/app_test.go`

```go
func TestPartialFailureRetrySucceeds(t *testing.T) {
    // Submit packet, fail mid-transaction, retry same packet
    // Assert: no duplicate entities
}
```

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

Uses real in-memory CRDB. No mocks.

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

## WORKSTREAM 4: DEVELOPER ENVIRONMENT

### 4a. Solvent migration FS

**File:** `internal/epistemic/migrations.go`

```go
package epistemic

import "embed"

//go:embed migrations/*.sql
var solventMigrations embed.FS

func ApplySolventMigrations(ctx context.Context, db *sql.DB) error {
    // Read and execute each .sql file in order
}
```

**Directory:** `internal/epistemic/migrations/`
- Copy all 10 Solvent `.sql` files from `solvent-main/db/`
- Keep naming consistent (001_schema.sql through 010_debt_opaque.sql)

### 4b. Wire into reset command

**File:** `cmd/argus/main.go`

Replace TODO placeholder:
```go
// Before (placeholder):
// solventmigrations.Apply(db)

// After (real):
epistemic.ApplySolventMigrations(ctx, db)
```

### 4c. Verify taskfile idempotency

- `task dev` → runs twice without error
- `task fresh` → wipes `.cockroach-data`
- `task test` → all tests pass
- `task lint` → clean

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
| `domain-pack/bmist/v1/types.go` | Modify | Add VerifierSpec type |
| `domain-pack/bmist/v1/pack.json` | Modify | Add verifier_specs |
| `domain-pack/bmist/v1/validate.go` | Modify | Validate VerifierSpecs |
| `verifier/model/artifact.go` | Modify | Add VerifierID, Tolerance |
| `verifier/physics/v1/verifier.go` | Modify | Set VerifierID in artifact |
| `internal/application/app.go` | Modify | Add Verify, VerifierSpec enforcement |
| `internal/epistemic/migrations.go` | Create | Solvent migration FS |
| `internal/epistemic/migrations/*.sql` | Create | Embed Solvent migrations |
| `cmd/argus/main.go` | Modify | Wire verifier, registry, migrations |
| `internal/application/app_test.go` | Create | Refusal tests 1,3,4,8,9 |
| `internal/mcp/adapter_test.go` | Create | Refusal tests 2,5 |
| `internal/ui/ui_test.go` | Create | Refusal tests 6,7 |
| `internal/boundary/boundary_test.go` | Create | Depguard verification |
| `integration_test.go` | Create | 17-step integration test |
| `.golangci.yml` | Create | Depguard rules |

---

## RISK CHECK

- **Solvent changes:** 2 of 2 budget used (projection contract + migration FS). No third change.
- **No architecture change:** All work is wiring and testing.
- **No agent authority bypass:** VerifierSpec enforcement is compile-time in `Validate()`.
- **No false causal ordering:** Integration test steps are explicit, no invented happens-before.
