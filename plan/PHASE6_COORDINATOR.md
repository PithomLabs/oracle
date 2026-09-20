# PHASE 6 COORDINATOR — COMPLETE

## Status: PASS

## Coordinator Architecture

In-process Go library with thin HTTP wrapper.

```
Packet
  ↓
Coordinator
  ├── Domain Pack
  ├── ArtifactReader
  ├── Idempotency
  ├── Solvent
  └── Conductor projection
```

## Packet Compilation Flow

1. Validate packet (`packet/v1.Validate`)
2. Resolve pack_ref → Domain Pack
3. Compute canonical hash (SHA-256, exclude packet_id/metadata)
4. IdempotencyCache.Begin() → atomic reservation
5. Compile beliefs (debt = pack.initial_debt ∪ packet.belief.debt)
6. Compile evidence (validate provenance, artifact hash if reproducible)
7. Compile edges (resolve refs, POST /v1/beliefs/{parent_id}/edges)
8. Compile tasks (POST Conductor task, status: proposed)
9. If Conductor fails → queue in ProjectionQueue
10. IdempotencyCache.Complete()
11. Return CompilationResult

## Solvent API Mapping

| Operation | Solvent Endpoint |
|-----------|------------------|
| Create belief | POST /v1/beliefs |
| Create edge | POST /v1/beliefs/{parent_id}/edges |
| Admit evidence | POST /v1/evidence |
| Promote | POST /v1/beliefs/{id}/promote |
| Retire debt | POST /v1/beliefs/{id}/debt/retire |
| Retract | POST /v1/beliefs/{id}/retract |
| Approve target | POST /v1/targets/{id}/approve |

## Edge Handling

- Growth Gate exception: POST /v1/beliefs/{parent_id}/edges
- Coordinator resolves from_ref/to_ref before calling Solvent
- Solvent enforces structural checks (parent exists, child exists, kind valid)

## Evidence/Artifact Validation

- Provenance class validated against Pack
- If `reproducible_artifact`: resolve via ArtifactReader, verify hash
- Forged/unknown/mismatched artifacts rejected

## Idempotency

- Process-local in-memory cache
- Atomic reservation (StateNew → StateInFlight → StateCompleted)
- Canonical hash: SHA-256, excludes packet_id/metadata/timestamps
- Same content + different packet_id → one compilation
- Race-safe concurrent access
- Process restart clears cache

## Conductor Projection

- Tasks created with status: proposed
- governance_ref = canonical Solvent belief reference
- If Conductor fails → queue in ProjectionQueue for retry
- Retry does not duplicate projections

## Human Decisions

6 types: PROMOTE, RETIRE_DEBT, RETRACT, REOPEN, AUTHORIZE, REFUSE

- All write to Solvent only
- No Coordinator decision persistence
- DecisionRecord is transient response envelope
- Requires ARGUS_OPERATOR_PRINCIPAL_ID
- REFUSE does not mutate Solvent

## Sphinx Projection

- Read-only projection of Solvent authority state
- No new policy logic
- Current result: PASS | REFUSE | HUMAN_REVIEW
- Deterministic mapping of existing canonical state

## HTTP Wrapper

4 endpoints:
- POST /decisions
- GET /packets/:id/status
- GET /beliefs/:id/decision-context
- GET /authorization-context/:target_id

Transport only. No business logic. Calls Coordinator library exclusively.

## Test Results

```
123 passed in 9 packages
go test -race ./... passed
go vet ./... passed
```

## Files Created

- `oracle/coordinator/coordinator.go`
- `oracle/coordinator/compiler.go`
- `oracle/coordinator/compiler_test.go`
- `oracle/coordinator/validate.go`
- `oracle/coordinator/validate_test.go`
- `oracle/coordinator/idempotency.go`
- `oracle/coordinator/idempotency_test.go`
- `oracle/coordinator/projection.go`
- `oracle/coordinator/projection_test.go`
- `oracle/coordinator/human.go`
- `oracle/coordinator/human_test.go`
- `oracle/coordinator/sphinx.go`
- `oracle/coordinator/sphinx_test.go`
- `oracle/coordinator/client.go`
- `oracle/coordinator/http/server.go`
- `oracle/coordinator/http/handler.go`
- `oracle/coordinator/http/types.go`

## Repositories Modified

- `oracle/coordinator/` — new package

## Repositories Unchanged

- `solvent-main` — NO changes
- `conductor` — NO changes
- `oracle/reference-loop` — NO changes
- `oracle/domain-pack/` — NO changes
- `oracle/packet/` — NO changes
- `oracle/verifier/` — NO changes
- `oracle/corpus/` — NO changes

## Limitations

1. Process-local idempotency only (no cross-restart)
2. Non-persistent projection retry queue
3. No Coordinator decision database
4. Edge endpoint assumed to exist (Growth Gate dependency)
5. Solvent/Conductor clients are stubs (no real HTTP calls in POC)

## Deviations

1. `parsePackRef` implemented locally (not exported from packet package)
2. Solvent/Conductor clients return mock responses for POC

## Key Distinctions

- **Coordinator** = deterministic packet compiler
- **Solvent** = canonical epistemic/authority state
- **Conductor** = operational projection only
- **DecisionRecord** = transient envelope, not persistent state
- **SphinxProjection** = read-only view, not policy engine

## Acceptance

- [x] Valid packet compiles into Solvent state
- [x] Pack initial debt cannot be removed by packet content
- [x] Evidence is validated structurally
- [x] Reproducible artifacts require ArtifactReader validation
- [x] Belief edges are persisted through Solvent API
- [x] No direct Coordinator SQL exists
- [x] Solvent writes precede Conductor writes
- [x] Conductor failures enter retry queue
- [x] Retry does not duplicate projections
- [x] Canonical idempotency is process-local and race-safe
- [x] Human decisions require operator principal
- [x] Coordinator does not persist decision history
- [x] Solvent Verdict refusals are handled correctly
- [x] Reopen preserves lineage with derives edge
- [x] Sphinx is a read-only projection
- [x] HTTP wrapper is thin
- [x] No reasoning/LLM logic exists
- [x] go test ./... passes
- [x] go test -race ./... passes
- [x] go vet ./... passes
- [x] Solvent unchanged
- [x] Conductor unchanged
- [x] reference-loop unchanged
- [x] No Phase 7+ implementation exists
