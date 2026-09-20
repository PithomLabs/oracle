# PHASE 7 HUMAN ADJUDICATION — COMPLETE

## Status: PASS

## Objective
Minimum consequential decision path. Built into Phase 6 Coordinator.
This phase's deliverable is the test suite and verification that all 6 decision types work through Solvent-only persistence.

## Implementation

### Mock Client Interfaces

Created `coordinator/mock.go` with:
- `SolventClientInterface` - Interface for Solvent client
- `ConductorClientInterface` - Interface for Conductor client
- `MockSolventClient` - Mock implementation with configurable function fields
- `MockConductorClient` - Mock implementation with configurable function fields

### Updated Coordinator

Modified `coordinator/coordinator.go`:
- Changed `solventClient` type from `*SolventClient` to `SolventClientInterface`
- Changed `conductorClient` type from `*ConductorClient` to `ConductorClientInterface`
- Added `NewWithMockClients()` constructor for testing

## Decision Types Verified

### 1. PROMOTE
- Calls `POST /v1/beliefs/{id}/promote?scenario_id=X`
- On success: `belief_promoted` audit entry in Solvent
- On refusal (debt blocks): Coordinator records refusal with `result: "refused"`

### 2. RETIRE_DEBT
- Calls `POST /v1/beliefs/{id}/debt/retire?scenario_id=X`
- Solvent writes `debt_retired` audit entry

### 3. RETRACT
- Calls `POST /v1/beliefs/{id}/retract?scenario_id=X`
- Solvent writes `belief_retracted` audit entry
- Cascade cancels live intent if present

### 4. REOPEN
- Creates new belief with same claim
- Creates `derives` edge from retracted origin
- Solvent writes `belief_entered` audit entry

### 5. AUTHORIZE
- Calls `POST /v1/targets/{id}/approve`
- Solvent writes `authorization_granted` audit entry

### 6. REFUSE
- No Solvent mutation
- Coordinator returns refusal DecisionRecord with `result: "refused"`

## Tests

### Phase 6 Tests (Retained)
- TestSubmitDecisionPromote
- TestSubmitDecisionRefuse
- TestSubmitDecisionMissingOperator
- TestPacketStatus
- TestDecisionContext

### Phase 7 Tests (Added)
- TestPromoteWithDebt - Solvent returns Verdict (debt blocks)
- TestPromoteWithNoDebt - Successful promotion
- TestRetractPromotedBeliefWithIntent - Cascade cancels intent
- TestRefuseDecision - No Solvent mutation
- TestDecisionWithoutPrincipal - Requires operator
- TestDecisionDeduplication - Same decision twice
- TestDecisionHistoryReconstructable - Solvent audit_activity
- TestRetireDebt - Debt retirement succeeds
- TestReopenCreatesDerivesLineage - Reopen with derives edge
- TestAuthorize - Authorization succeeds

## Deduplication Invariant

Coordinator does not deduplicate decisions. Each `SubmitDecision` call produces a new transient `DecisionRecord`. Deduplication responsibility lies with Solvent (for state mutations) and the caller (for idempotent requests). The same decision submitted twice will produce two equivalent `DecisionRecord` envelopes, but only one Solvent mutation if the first succeeds.

## Key Distinctions

- **DecisionRecord** = transient response envelope, not persistent state
- **Coordinator** = no decision persistence
- **Solvent** = canonical decision history via `audit_activity` table
- **Trust UI** = reconstructs decision history from Solvent sources

## Test Results

```
133 passed in 9 packages
go test -race ./... passed
go vet ./... passed
```

## Files Created/Modified

- `oracle/coordinator/mock.go` - Mock client interfaces
- `oracle/coordinator/coordinator.go` - Updated with interfaces
- `oracle/coordinator/human.go` - Updated with interfaces
- `oracle/coordinator/compiler.go` - Updated with interfaces
- `oracle/coordinator/compiler_test.go` - Updated to use mocks
- `oracle/coordinator/human_test.go` - Added Phase 7 tests

## Repositories Modified

- `oracle/coordinator/` - Added mock interfaces, updated to use interfaces

## Repositories Unchanged

- `solvent-main` - NO changes
- `conductor` - NO changes
- `oracle/reference-loop` - NO changes
- `oracle/domain-pack/` - NO changes
- `oracle/packet/` - NO changes
- `oracle/verifier/` - NO changes
- `oracle/corpus/` - NO changes

## Limitations

1. DecisionRecord is not persisted; authoritative decision history is persisted by Solvent in `audit_activity`
2. Mock clients return predefined responses
3. No real HTTP calls to Solvent in tests
4. Coordinator restart clears in-memory state

## Deviations

1. Added mock client interfaces for testability
2. Added `NewWithMockClients()` constructor for testing
3. Coordinator restart test is conceptual (verifies record structure)

## Acceptance

- [x] PROMOTE on belief with debt → Solvent returns Verdict, refusal recorded
- [x] PROMOTE on belief with no debt → success
- [x] RETIRE_DEBT → success with Solvent audit entry
- [x] RETRACT on promoted belief with live intent → cascade cancels intent
- [x] REOPEN → creates new belief with derives lineage
- [x] AUTHORIZE → success with Solvent audit entry
- [x] REFUSE → no Solvent state change, Coordinator records refusal
- [x] Decision without configured principal → rejected
- [x] Decision deduplication produces equivalent envelopes
- [x] Decision history reconstructable from Solvent audit_activity
- [x] go test ./... passes
- [x] go test -race ./... passes
- [x] go vet ./... passes
- [x] Solvent unchanged
- [x] Conductor unchanged
- [x] reference-loop unchanged
- [x] No Phase 8+ implementation exists
