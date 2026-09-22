# Plan: Convert UUID generation to UUID v7 monotonic

Date: 2026-09-21
Status: Plan-only

## Scope

Convert nondeterministic/random UUID generation to UUID v7. Keep `EntityID()` deterministic and unchanged. Do not rewrite existing persisted UUIDs.

## What changes

### 1. `internal/work/store.go` — `generateID()` function

Current:
```go
func generateID() string { return uuid.New().String() }
```

New:
```go
func generateID() string {
    id, err := uuid.NewV7()
    if err != nil {
        // Fallback to v4 if v7 fails (should not happen in practice)
        return uuid.New().String()
    }
    return id.String()
}
```

This affects 5 call sites:
- `Create()` — fallback task ID
- `Transition()` — activity event ID
- `Claim()` — activity event ID
- `Release()` — activity event ID
- `AddDependency()` — dependency ID

### 2. SQL `gen_random_uuid()` defaults — NO CHANGE

The SQL defaults are fallbacks for:
- `conductor_project.id`
- `conductor_task.id`
- `conductor_dependency.id`
- `conductor_activity.id`
- `belief_retirement_proposal.id`
- `principal.principal_id`
- `debt_discharge.discharge_id`
- `belief.id`
- `evidence.id`
- `action_intent.id`

In practice, the Go code always provides explicit IDs for beliefs, evidence, and tasks via `EntityID()`. The SQL defaults only fire for seed data or direct SQL inserts. Changing them would require:
- A new migration to ALTER DEFAULT
- CockroachDB version-specific UUID v7 function (uncertain availability)
- Risk of hotspot issues per CockroachDB docs

**Decision: Keep SQL defaults as `gen_random_uuid()` (UUID v4).** The Go-side change covers all active ID generation.

## Files modified

| File | Change |
|------|--------|
| `internal/work/store.go` | `generateID()` uses `uuid.NewV7()` with v4 fallback |

## What does NOT change

- `EntityID()` in `app.go` — deterministic content-derived IDs, unchanged
- SQL `gen_random_uuid()` defaults — kept as UUID v4
- Seed data UUIDs — fixed constants, unchanged
- Existing persisted UUIDs — not rewritten
- `uuid.Parse()` validation — unchanged

## Verification

1. `go vet ./...`
2. `go test ./...`
3. Verify generated IDs are UUID v7 format (version nibble = 7)

## Risk assessment

- **Low risk**: Only one function changes (`generateID()`)
- **Backward compatible**: UUID v7 is still a valid UUID; existing code that parses UUIDs will work
- **No schema changes**: Column types remain `UUID`
- **No data migration**: Existing rows keep their UUIDs
