Plan 10r2 is improved, but I would **not start implementation yet**. There are four substantive issues remaining.

### 1. Fix 5 has a Go type/design contradiction

This is the biggest one.

The plan proposes:

```go
type RetirementRule interface {
    GetEvidenceClass() string
    GetRule() string
}

type RetirementRulesProvider interface {
    GetRetirementRules() map[string]RetirementRule
}
```

but the BM-IST implementation currently owns:

```go
map[string]bmistv1.RetirementRule
```

In Go, `map[string]bmistv1.RetirementRule` does **not** implement `map[string]domainpack.RetirementRule`, even when the element type implements the interface.

Worse, the plan says `domain-pack/registry.go` should import `bmistv1` while also having BM-IST implement a generic `domainpack` contract; that can easily create the very module dependency cycle the plan says it avoids. 

The cleaner solution is to use the **already-approved generic `PackDefinition` contract**:

```text
Pack
  ↓
PackDefinition
  ↓
RetirementRules
```

The application consumes the generic typed definition. The generic registry should not special-case BM-IST.

### 2. Fix 3 has a concurrency hole

The plan does:

```text
SELECT existing edge
        ↓
if different → error
        ↓
INSERT ... ON CONFLICT DO NOTHING
```

Two concurrent transactions can both observe no row, then one inserts and the other gets `DO NOTHING`. The second request can therefore still silently lose a conflicting edge. 

After `RowsAffected() == 0`, re-read the existing edge kind and return an explicit conflict when it differs.

### 3. Port propagation is still incomplete

The resolved ports are shell variables, but the plan does not establish a persistent runtime configuration that later commands/tests can consume. Also `argus mcp` still has a hardcoded CRDB port. 

Use:

```text
argus resolve-ports
        ↓
.argus-runtime.env
        ↓
task dev / task test / argus mcp / integration tests
```

and make **CRDB startup itself** retry/re-resolve on `EADDRINUSE`, not just ARGUS.

Otherwise the exact failure you experienced can still occur during the CRDB bind phase.

### 4. The Solvent migration change is larger than the plan claims

Moving the canonical SQL and then modifying those SQL files to add `IF NOT EXISTS` is no longer purely a packaging change. 

That can still be acceptable, but the plan should explicitly classify it as:

> **one bounded Solvent migration-packaging change that also makes the existing migration scripts replay-safe**

and verify that the semantic behavior of the migrations is unchanged.

Also reuse the existing Solvent migration statement-splitting logic rather than maintaining two implementations.

### Verdict

```text
Architecture:          ✅
Direction:             ✅
Plan:                  ⚠️ revise before build

Required:
1. Resolve Pack/RetirementRule type architecture
2. Make edge conflict handling concurrency-safe
3. Complete runtime port propagation + CRDB retry
4. Clarify Solvent migration change scope/idempotency
```

Everything else in 10r2 is substantially aligned with the pivot, including fail-closed operator authentication, trusted reproducible artifacts, MCP SDK adoption, and `UNKNOWN ≠ EMPTY`. 

### Prompt to the coding agent

```text
Revise Plan 10r2 only; do not implement yet.

Fix these four issues:

1. Retirement rules:
   Do NOT special-case BM-IST in the generic PackRegistry and do NOT use
   map[string]interface{} or incompatible map[string]interface-type> interfaces.
   Use the existing generic PackDefinition/typed contract. The application must
   consume generic typed retirement rules without importing BM-IST.

2. Edge conflicts:
   Make conflict detection concurrency-safe. After INSERT ... ON CONFLICT DO NOTHING
   reports no row inserted, re-read the existing edge kind and return an explicit
   conflict when it differs. Same-kind retry remains idempotent.

3. Port lifecycle:
   Persist resolved ports in `.argus-runtime.env` (or equivalent runtime config).
   All task commands, MCP, readiness probes, and integration tests consume the
   resolved topology. Remove hardcoded runtime CRDB ports. CRDB startup must also
   retry/re-resolve on EADDRINUSE. Test the exact CRDB-admin/ARGUS :8080 collision.

4. Solvent migrations:
   Explicitly classify moving the canonical SQL + making existing migration
   statements replay-safe as the single bounded migration-export change.
   Preserve one authoritative SQL location and reuse Solvent's existing
   statement-splitting logic.

Do not add services, databases, REST boundaries, or a third Solvent change.
After revision, freeze the plan.
```

The first item is the one I would resolve most carefully; it is a real package-level design issue, not merely a type cleanup.
