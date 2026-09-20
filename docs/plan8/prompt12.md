I reviewed Plan 10. It is materially better, but **I would not execute it unchanged**. There are five important corrections, with two being blockers.

### Blockers

**1. Fix 5 does not actually enforce retirement rules.**

The proposed check only rejects a *wrong non-empty* evidence class:

```go
if req.EvidenceClass != "" && req.EvidenceClass != rule.EvidenceClass
```

An empty evidence class therefore passes. More importantly, it does not verify that qualifying persisted evidence actually exists for that belief/scenario. The architecture requires mechanical rule enforcement before `Discharge`. 

The implementation should require:

```text
debt item exists in Pack
        ↓
retirement rule exists
        ↓
required evidence class known
        ↓
qualifying persisted evidence exists for belief/scenario
        ↓
evidence satisfies rule
        ↓
Discharge
```

Also use the typed `PackDefinition`; do not introduce a parallel `map[string]interface{}` retirement-rule model.

**2. Fix 10/11/12 port handling is internally inconsistent.**

The plan proposes dynamic port fallback, but then Fix 12 standardizes tests on a fixed `26257`.  

That will fail precisely when fallback is exercised.

Use one runtime-resolved topology:

```text
port preflight
    ↓
resolve SQL/admin/ARGUS ports
    ↓
start CRDB with those ports
    ↓
write runtime env/config
    ↓
ARGUS + integration tests consume resolved ports
```

Do not hard-code the integration test port once dynamic allocation exists.

There is also a race in:

```go
net.Listen()
close()
return port
```

Another process can claim the port between the probe and the actual bind. Either reserve listeners until startup or make startup retry atomically on bind failure.

### Important corrections

**3. Solvent migration packaging should not duplicate SQL unnecessarily.**

The plan proposes copying/symlinking the existing Solvent SQL into `migrations/db/` solely to satisfy `go:embed`. 

That is acceptable technically, but only if the existing canonical files are moved so there is **one authoritative copy**. Do not create two copies that can drift.

Better:

```text
Solvent
  migrations/db/*.sql   ← canonical location
  migrations.Apply()
```

with existing migration paths updated as necessary.

**4. Static Bearer authentication is incomplete for browser UI.**

The plan protects `/ui/api/*` with:

```text
Authorization: Bearer <token>
```

but ordinary browser form submissions/navigation do not automatically provide that header. 

For the POC, either make the UI explicitly use `fetch()` with a user-provided token, or use a tiny login-to-HttpOnly-cookie mechanism. Do not put the token in URLs. The security boundary must be executable, not just defined.

**5. Official MCP SDK integration must be verified against the actual SDK API.**

The architecture choice is right, but the proposed code is specific enough that the agent should inspect the installed SDK version rather than assume constructor/types such as `mcp.NewServer(&mcp.ServerConfig{...})`. 

Keep:

```text
OpenCode → official MCP SDK → internal/mcp → application
```

but verify the actual SDK API before implementation.

### Other two worthwhile corrections

The custom semver parser is still not fully SemVer-compliant because it strips prerelease information.  Either implement proper prerelease ordering or constrain `VerifierSpec` to normalized release versions and reject prereleases.

The verifier artifact path should remain strict: `argus verify --submit` and packet verification must go through the same validation path, with the application's trusted verifier producing/validating the artifact. The plan should not allow an agent to submit an arbitrary artifact and have it become trusted merely because its declared hash matches. 

## Recommended prompt

```text
Revise Plan 10 before implementation. Do not change the architecture.

Fix these items:

1. Retirement-rule enforcement:
   - require a known debt item + Pack retirement rule
   - require exact required evidence class
   - require qualifying persisted evidence for the belief/scenario
   - reject missing/empty/mismatched evidence class
   - use the typed PackDefinition; no map[string]interface{} duplicate model

2. Port lifecycle:
   - resolve CRDB SQL, CRDB admin, and ARGUS HTTP ports before startup
   - propagate resolved ports to CRDB, ARGUS, readiness, and integration tests
   - eliminate fixed 26257/26260 assumptions when fallback is active
   - avoid probe-then-close races; retry atomically on bind failure
   - add conflict/fallback/concurrent-start tests
   - ensure task down/Ctrl-C cleans only processes started by this invocation

3. Solvent migrations:
   - one canonical SQL location in Solvent
   - exported Apply/embedded FS consumes that canonical copy
   - no duplicate SQL trees

4. UI authentication:
   - make the browser acquisition/use of the static operator credential explicit
   - no token in URL, no body-supplied principal
   - consequential actions must be authenticated and CSRF/Origin protected

5. MCP:
   - use the official MCP Go SDK
   - inspect the actual SDK version/API before coding; do not assume signatures

6. SemVer:
   - implement real SemVer comparison or explicitly validate/reject unsupported prerelease forms.

7. Verifier:
   - only trusted verifier execution/application validation can create qualifying reproducible-artifact evidence
   - agent-declared verifier/hash fields are never sufficient by themselves.

Keep scope frozen. No third Solvent change, no REST resurrection, no new service/database. After revision, freeze the plan.
```

The overall direction remains right; the remaining work is implementation correctness, especially **authority enforcement and runtime topology**, not another architectural pivot.
