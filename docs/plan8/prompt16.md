Plan 10r4 is **not quite final**. It fixed the previous BM-IST coupling and evidence-authority problems, but I found **five remaining blockers**.

### 1. P0 — `PackRef` is still client-controlled for human decisions

The plan now resolves the pack from:

```go
req.PackRef
```

and the Trust UI sends `PackRef` in the request body. 

That violates the earlier authority rule: **PackRef must be derived server-side from authoritative scenario metadata**, not trusted from the UI.

Change the flow to:

```text
authenticated human request
        ↓
scenario_id
        ↓
server resolves authoritative scenario.pack_ref
        ↓
PackRegistry.Get(...)
```

The request must not be able to select another domain pack.

### 2. P0 — `PrincipalID` is still externally representable

`DecisionRequest` still contains:

```go
PrincipalID string
```

even though the plan says the principal is server-derived. 

Separate:

```text
ExternalDecisionRequest
    no principal

AuthenticatedDecisionCommand
    principal supplied by trusted auth boundary
```

Otherwise a future caller can construct `DecisionRequest{PrincipalID:"someone-else"}` directly against `App`.

### 3. P0 — `task dev` still invokes `argus reset`

The Taskfile says:

```text
CRDB ready
→ argus reset
→ argus serve
```



Earlier we explicitly established:

```text
task dev   = non-destructive
task fresh = destructive
```

If `argus reset` retains destructive semantics, this violates the developer contract.

Separate:

```text
argus migrate / ensure-db   # non-destructive
argus reset                 # explicitly destructive
```

Then `task dev` uses migrate/ensure, while `task fresh` uses reset.

### 4. P0 — ARGUS HTTP startup retry is claimed but not implemented

The plan says both CRDB **and ARGUS** retry on `EADDRINUSE`, but the actual Taskfile only shows CRDB retry; ARGUS is simply launched once.  

Implement the same bounded retry/re-resolution for ARGUS.

Otherwise the exact resilience requirement you requested is still incomplete.

### 5. P1 — runtime port resolution is vulnerable to concurrent `task dev`

`.argus-runtime.env` and `.argus-pids/` are shared global state. Two concurrent invocations can overwrite each other's topology and PID files.

At minimum add a startup lock:

```text
.task/dev.lock
```

so only one local ARGUS dev stack owns the shared runtime directory.

Also handle an already-running ARGUS/CRDB instance explicitly rather than starting another process against the same `.cockroach-data`.

---

### Two smaller corrections

`parsePackRef()` is unsafe for pack IDs containing `-`; splitting on the first hyphen makes something like `market-risk-1.0.0` ambiguous. Prefer `packID@version` as the canonical format. 

The `GetContext` availability model has only one `Reason`. If both work and epistemic sources fail, one reason overwrites the other. The original RCP requirement was per-source availability, so use per-section status rather than one shared reason. 

## Prompt to coding agent

```text
Revise Plan 10r4 one final time. Do NOT implement yet.

Fix these remaining issues:

1. Human decision PackRef:
   Never trust PackRef from the request body. Resolve it server-side from
   authoritative scenario metadata using scenario_id. The UI must not select
   the domain pack.

2. Human principal:
   Remove PrincipalID from externally supplied DecisionRequest data.
   Separate external request DTO from an internal authenticated command carrying
   the server-derived principal.

3. Developer lifecycle:
   `task dev` MUST be non-destructive. Separate `argus migrate/ensure-db`
   from destructive `argus reset`. `task dev` applies idempotent migrations;
   `task fresh` invokes destructive reset.

4. ARGUS port resilience:
   Implement the same bounded EADDRINUSE re-resolution/retry for ARGUS that
   exists for CRDB. Do not claim this behavior unless it is implemented and tested.

5. Concurrent dev startup:
   Protect `.argus-runtime.env`, `.argus-pids/`, and `.cockroach-data` with
   a local startup lock or equivalent ownership mechanism. Handle an already
   running ARGUS/CRDB instance explicitly; never start a second process against
   the same store accidentally.

6. PackRef syntax:
   Use an unambiguous canonical form such as `pack_id@version`; do not split
   arbitrary pack IDs on the first hyphen.

7. RCP availability:
   Make availability per section/source so work and epistemic failures cannot
   overwrite each other's reason. Preserve UNKNOWN != EMPTY.

Keep the architecture frozen. No third Solvent change, no REST resurrection,
no new services/databases.

After these corrections, freeze Plan 10r4 and proceed to implementation.
```

The core architecture itself is now sound; these remaining issues are mostly **authority provenance and developer-runtime correctness**, not reasons to reopen the pivot. 
