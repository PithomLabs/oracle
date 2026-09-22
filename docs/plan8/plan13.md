# Plan 13 — MCP Production Path Validation Remediation

## Source

Remediation of adversarial code-review findings from `docs/plan8/review7.md`.
Wire the canonical `packetv1.Validate()` into the production MCP submission path,
add defense-in-depth checks in `Persist()`, fix the `parsePackRef` delimiter bug,
and add task FK preflight.

## Reconciliation of 12 Review Items

### Confirmed gaps (need implementation)

| # | Finding | Where now | Fix |
|---|---------|-----------|-----|
| 1-3 | Agent identity + role consistency not required | `Compile()`/`Validate()` skip agent check | Wire `packetv1.validateAgent()` |
| 4-6 | Duplicate local_ids in beliefs/evidence/edges/tasks | No validation, map corruption | `packetv1` checks + Persist defense |
| 7 | Invalid claim_type silently inserted | No enum check anywhere | Add enum in `validateBeliefs` + Persist |
| 8 | Invalid edge kind — only DB catches | DB CHECK constraint only | Add Go validation before INSERT |
| 9 | Reference format not validated | Generic "unresolved" error | `packetv1.validateReferences()` |
| — | `packetv1.parsePackRef` uses `-` not `@` | Incompatible with production | Unify to `@` |
| — | Task FK: missing conductor_project | Opaque FK error | Preflight query |

### Already enforced (no action)

| # | Finding | Already enforced by |
|---|---------|-------------------|
| 10 | local contradicts target | `Persist()` lines 441-444 |
| 11 | cross-scenario edge | `Persist()` `validateEdgeScenario` |
| 12 | self-edge | `Persist()` lines 437-439 |

### False/duplicate/inconsistent in review

- The review's edge validation descriptions were internally contradictory. Confirmed: Persist already catches #10, #11, #12.
- FK ordering hypothesis not a demonstrated failure — CockroachDB handles deferred FK within transactions.
- `contentHash` excluding agent identity is a legitimate semantic choice (dedup research, not agents).
- Evidence/edge provenance missing from UI is projection gap, not authority failure.

---

## Implementation Phases

### Phase 1: Fix `parsePackRef` delimiter

**File:** `packet/v1/validate.go` line 96
- Change `strings.SplitN(ref, "-", 2)` to `strings.SplitN(ref, "@", 2)`
- Update doc comment to say `"bmist@1.1.0"`

**File:** `packet/v1/schema.json` line 33
- Change description from `"pack_id-version"` to `"pack_id@version"`

**File:** `packet/v1/validate_test.go` lines 34, 96
- Change `"bmist-1.0.0"` to `"bmist@1.0.0"` (2 occurrences)

### Phase 2: Add claim_type enum validation

**File:** `packet/v1/validate.go` — in `validateBeliefs()` (~line 119)
- After non-empty check, add:
  ```go
  if b.ClaimType != "derived" && b.ClaimType != "accommodated" && b.ClaimType != "postulated" {
      return fmt.Errorf("belief[%d]: invalid claim_type %q", i, b.ClaimType)
  }
  ```

### Phase 3: Wire `packetv1.Validate()` into production path

**File:** `internal/application/app.go` — new method
```go
func (a *App) ValidatePacket(ctx context.Context, pkt *packetv1.Packet) error {
    if a.packRegistry == nil {
        return nil
    }
    return packetv1.Validate(pkt, a.packRegistry)
}
```

**File:** `internal/mcp/adapter.go` — `handleSubmitPacket`
- Add `a.app.ValidatePacket(ctx, &pkt)` call between `app.Validate()` and `BeginTx`

Pipeline: `json.Unmarshal` → `Compile()` → `Validate()` → **`ValidatePacket()`** → `BeginTx` → `Persist()` → `Commit`

### Phase 4: Defense-in-depth in `Persist()`

**File:** `internal/application/app.go` — at top of `Persist()` after existing checks
- Build combined local_id set across all beliefs/evidence/edges/tasks
- Reject on duplicates with clear error message identifying which entity type and local_id
- Before belief INSERT: validate `b.ClaimType` against allowed enum values
- Before edge INSERT: validate `edge.Kind` against `"derives"`/`"contradicts"`

### Phase 5: Task integrity preflight

**File:** `internal/application/app.go` — in `Persist()`, before task INSERTs
```go
if len(pkt.Tasks) > 0 {
    var exists bool
    err := tx.QueryRowContext(ctx,
        `SELECT EXISTS(SELECT 1 FROM conductor_project WHERE id = $1::UUID)`,
        pkt.ScenarioID,
    ).Scan(&exists)
    if err != nil {
        return nil, fmt.Errorf("verify project for scenario: %w", err)
    }
    if !exists {
        return nil, fmt.Errorf("cannot create tasks: project %s does not exist", pkt.ScenarioID)
    }
}
```
- Preserve existing DB FK as final safeguard
- Do NOT auto-create projects

### Phase 6: Regression tests

**File:** `packet/v1/validate_test.go` — unit tests (no DB)

| Test | Malformed input | Expected error |
|------|----------------|---------------|
| `TestRejectMissingAgentID` | Agent with empty ID | `agent.id is required` |
| `TestRejectMissingAgentHarness` | Agent with empty harness | `agent.harness is required` |
| `TestRejectMissingAgentModel` | Agent with empty model | `agent.model is required` |
| `TestRejectRoleMismatch` | agent.role="work", packet.role="adversarial" | `agent.role does not match packet.role` |
| `TestRejectDuplicateBeliefLocalID` | Two beliefs with local_id="b1" | `duplicate local_id` |
| `TestRejectDuplicateEvidenceLocalID` | Two evidence with local_id="e1" | `duplicate local_id` |
| `TestRejectDuplicateEdgeLocalID` | Two edges with local_id="ed1" | `duplicate local_id` |
| `TestRejectDuplicateTaskLocalID` | Two tasks with local_id="t1" | `duplicate local_id` |
| `TestRejectInvalidClaimType` | claim_type="pizza" | `invalid claim_type` |
| `TestRejectInvalidEdgeKind` | kind="invalid" | `invalid edge kind` |
| `TestRejectMalformedRef` | belief_ref="garbage" | `invalid reference format` |
| `TestRejectBadPackRef` | pack_ref="bmist-1.0.0" | `pack_ref resolution failed` |

**File:** `internal/application/app_test.go` — integration tests (real DB)

| Test | What it verifies |
|------|-----------------|
| `TestRejectDuplicateLocalIDThroughPersist` | Duplicate belief local_id → error, transaction rolled back, no partial state |
| `TestRejectInvalidClaimTypeThroughPersist` | Invalid claim_type → error |
| `TestRejectInvalidEdgeKindThroughPersist` | Invalid edge kind → error (before DB CHECK) |
| `TestRejectMissingProjectForTasks` | Packet with tasks + non-existent project → clear error, no partial state |
| `TestPacketWithTasksAndExistingProject` | Packet with tasks + existing project → succeeds |

### Phase 7: Verify

```bash
go vet ./...
go test ./...
go build ./...
```

---

## Files Changed (estimated)

| File | Change type |
|------|------------|
| `packet/v1/validate.go` | Fix parsePackRef delimiter, add claim_type enum |
| `packet/v1/schema.json` | Update pack_ref description |
| `packet/v1/validate_test.go` | Fix pack_ref format, add 12 unit tests |
| `internal/application/app.go` | Add ValidatePacket(), defense-in-depth, task FK preflight |
| `internal/application/app_test.go` | Add 5 integration tests |
| `internal/mcp/adapter.go` | Wire ValidatePacket() into pipeline |

## Invariants Now Enforced

1. Agent identity required: `agent.id`, `agent.role`, `agent.harness`, `agent.model` non-empty
2. Role consistency: `agent.role == packet.role`
3. Packet role valid: must be `"work"` or `"adversarial"`
4. Belief `local_id` unique within packet
5. Evidence `local_id` unique within packet
6. Edge `local_id` unique within packet
7. Task `local_id` unique within packet
8. Claim type valid: must be `"derived"`, `"accommodated"`, or `"postulated"`
9. Edge kind valid: must be `"derives"` or `"contradicts"` (Go check, not just DB)
10. References well-formed: `local:` or `canonical:belief:` prefix required
11. Pack ref format: `@` delimiter, resolves via pack registry
12. Task project exists: `conductor_project` row must exist for `pkt.ScenarioID`
13. Atomicity: any failure rolls back entire packet (unchanged — `tx.Rollback()` on defer)
