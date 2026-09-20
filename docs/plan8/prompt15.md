Plan 10r3 is much closer, but I would **not freeze it yet**. I see **five important issues**, two of them architectural.

### 1. P0 — `domain-pack` ↔ `bmistv1` creates an import-cycle risk and breaks domain isolation

The plan makes generic `domain-pack/registry.go` import `bmistv1` so `decodePack()` can special-case `"bmist"`, while `bmistv1` imports the parent `domain-pack` package for the generic types. 

That creates:

```text
domain-pack
   ↓ imports
domain-pack/bmist/v1
   ↓ imports
domain-pack
```

Go will reject that.

More importantly, it makes the generic registry know BM-IST, which violates the core portability thesis.

**Fix:** use registration/injection:

```text
PackRegistry
  └── generic decoder map
        "bmist" → bmistv1.ParsePack
```

The composition root registers the BM-IST decoder. Generic `domain-pack` never imports BM-IST.

### 2. P0 — `resolvePack()` hardcodes BM-IST into the application core

The plan explicitly says:

```go
return a.packRegistry.Get("bmist", "v1.0.0")
```



That is directly contrary to the goal that physics is merely the reference domain.

The application should resolve the pack from **authoritative scenario/task metadata or packet `PackRef`**, not from `"bmist"` embedded in application code.

Otherwise replacing BM-IST with a market pack requires modifying core application code.

### 3. P0 — `operator_asserted` evidence is still incorrectly agent-attested

The plan says:

> “For `operator_asserted` evidence, agent-attested hash is accepted.” 

That conflicts with the established EBP doctrine: `operator_asserted` requires **human attestation**, not an agent merely declaring that provenance class.

The safe POC rule is:

```text
Agent submit_packet
    → reproducible_artifact: trusted verifier required
    → operator_asserted: proposal only / rejected as authoritative evidence

Human decision path
    → may create operator_asserted evidence with authenticated principal
```

This should be corrected before implementation.

### 4. P1 — port resolver still has a dangerous failure path

`resolvePort()` returns the preferred port after exhausting all fallback ports:

```text
all occupied
    ↓
return preferred
```



That defeats the whole resilience mechanism.

It must return an error:

```text
no available port
    → fail before starting anything
```

Also, the Taskfile's CRDB retry needs an explicit **success flag**. After three failed attempts, it must not proceed to migration/start ARGUS anyway. The current flow can fall through. 

### 5. P1 — cookie authentication lost the explicit Origin/CSRF check

The previous plan had `requireOrigin`; 10r3 now only has cookie/Bearer authentication. 

`SameSite=Strict` helps, but for consequential browser actions I would retain the explicit localhost Origin/Host validation we previously agreed on.

```text
authenticated
+
expected Origin/Host
=
allowed consequential UI mutation
```

No full CSRF framework is necessary.

### Smaller correction

The migration `Apply()` implementation ignores `ReadDir` and `ReadFile` errors:

```go
entries, _ := FS.ReadDir("db")
data, _ := FS.ReadFile(...)
```

Those must be handled. A broken embedded migration must fail closed, not potentially panic or silently continue. 

## Verdict

```text
Architecture:        ✅
Direction:            ✅
Plan:                 ⚠️ not frozen

Must fix:
1. generic PackRegistry must not import BM-IST
2. application must not hardcode "bmist"
3. operator_asserted cannot be agent-authoritative
4. port exhaustion must fail closed
5. restore Origin/Host protection
```

### Prompt to coding agent

```text
Revise Plan 10r3 one final time. DO NOT implement yet.

Fix only these issues:

1. Remove the domain-pack import cycle and BM-IST leakage.
   - generic domain-pack/registry must NOT import bmistv1
   - use decoder/factory registration from the composition root
   - generic registry remains domain-agnostic

2. Remove hardcoded "bmist@v1.0.0" from internal/application.
   Resolve PackRef from authoritative scenario/task metadata or packet metadata.
   Missing/unknown pack must fail closed.

3. Fix evidence authority:
   - reproducible_artifact requires trusted verifier artifact
   - operator_asserted cannot be authoritative from agent submit_packet
   - operator_asserted requires authenticated human attestation

4. Fix port resolution:
   - exhausted port range returns error, never the occupied preferred port
   - CRDB startup retries must stop with failure after max attempts
   - ARGUS startup must have equivalent retry/failure semantics
   - no service starts if port topology cannot be resolved

5. Restore explicit Origin/Host validation on consequential UI POSTs in addition
   to token/cookie authentication.

6. Handle all migration FS errors explicitly; Apply must fail closed.

Preserve the single-process architecture and do not add a third Solvent change.
Then freeze the plan.
```

The biggest thing to protect here is the **domain-agnostic core**. The moment `internal/application` or generic `domain-pack` starts knowing `"bmist"`, the pivot begins recreating the very coupling we were trying to eliminate.
