# PHASE 6 PLAN — DETERMINISTIC COORDINATOR

## Files to Create (15 total)

```
oracle/coordinator/
├── coordinator.go
├── compiler.go
├── compiler_test.go
├── validate.go
├── validate_test.go
├── idempotency.go
├── idempotency_test.go
├── projection.go
├── projection_test.go
├── human.go
├── human_test.go
├── sphinx.go
├── sphinx_test.go
├── client.go
└── http/
    ├── server.go
    ├── handler.go
    └── types.go
```

## Step 1: `coordinator.go` — Core Types + Config

```go
package coordinator

type Config struct {
    SolventBaseURL   string
    ConductorBaseURL string
    PackRegistry     *domainpack.PackRegistry
    ArtifactReader   verifier.ArtifactReader
    OperatorID       string // ARGUS_OPERATOR_PRINCIPAL_ID
}

type Coordinator struct { ... }
func New(cfg Config) *Coordinator
```

## Step 2: `compiler.go` — Packet Compilation

Flow:
1. Validate packet (`packet/v1.Validate`)
2. Resolve pack_ref → Domain Pack
3. Check idempotency cache → return cached if hit
4. For each belief: compile debt, POST /v1/beliefs
5. For each evidence: validate, POST /v1/evidence
6. For each edge: resolve refs, POST /v1/beliefs/{parent_id}/edges
7. For each task: POST Conductor task (status: proposed)
8. If Conductor fails → queue in ProjectionQueue
9. Store in idempotency cache
10. Return CompilationResult

## Step 3: `validate.go` — Packet + Debt Validation

- Structural validation via `packet/v1.Validate`
- Debt membership check against Pack vocabulary
- Evidence class validation against Pack
- Reference resolution (local/canonical)

## Step 4: `idempotency.go` — Process-Local Cache

```go
type IdempotencyCache struct {
    mu      sync.RWMutex
    entries map[string]*CompilationResult
}

func (c *IdempotencyCache) Begin(canonicalHash string) (*CompilationResult, bool)
func (c *IdempotencyCache) Complete(canonicalHash string, result *CompilationResult)
```

## Step 5: `projection.go` — Retry Queue

```go
type ProjectionQueue struct { ... }
func (q *ProjectionQueue) Enqueue(record ProjectionRecord)
func (q *ProjectionQueue) RetryPending(ctx context.Context) error
```

## Step 6: `human.go` — Decision Logic

6 decision types: PROMOTE, RETIRE_DEBT, RETRACT, REOPEN, AUTHORIZE, REFUSE

All write to Solvent only. No Coordinator persistence.

## Step 7: `sphinx.go` — Authorization Context

Read-only projection of Solvent authority state.

## Step 8: `client.go` — Solvent + Conductor HTTP Clients

Typed REST clients for all required endpoints.

## Step 9: `http/` — Thin HTTP Wrapper

4 endpoints: POST /decisions, GET /packets/:id/status, GET /beliefs/:id/decision-context, GET /authorization-context/:target_id

## Verification

```bash
cd /home/chaschel/Documents/go/oracle
go test ./...
go test -race ./...
go vet ./...
```
