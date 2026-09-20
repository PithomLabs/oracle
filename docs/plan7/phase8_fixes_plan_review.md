I reviewed the fixes plan against the three clarification questions. The decisions should be:

### 1. `setup:build`: do **not** create stubs

**Decision: build only binaries that actually exist and are required by this fixes scope. Do not create placeholder `main.go` files.**

The fixes plan explicitly says this is limited to “three targeted fixes only — no architectural expansion.” 

For this plan, the newly required binaries are:

```text
cmd/verifier/main.go
cmd/coordinator/main.go
```

Those are explicitly part of P0-2. 

`argus-mcp` and `trust-ui` should **not** be fabricated merely to satisfy `setup:build`. In fact, the exact-file list for this fixes plan does not include either as files to create. 

So change:

```yaml
setup:build:
  cmds:
    - go build ... ./cmd/verifier
    - go build ... ./cmd/coordinator
    - go build ... ./cmd/argus-mcp
    - go build ... ./trust-ui
```

to the binaries actually implemented by this phase:

```yaml
setup:build:
  cmds:
    - go build -o .tmp/bin/verifier ./cmd/verifier
    - go build -o .tmp/bin/coordinator ./cmd/coordinator
```

Do not leave dead build references as TODOs in the main `task dev` path. MCP can have a separate build task later; Trust UI remains outside these three fixes.

---

### 2. Coordinator authentication: **add API keys to `Coordinator.Config`, but only as transport configuration**

**Decision: add `SolventAPIKey` and `ConductorAPIKey` to `Coordinator.Config`.**

The launcher reads the environment variables:

```text
SOLVENT_API_KEY
CONDUCTOR_API_KEY
```

and passes them into `Coordinator.Config`. The Coordinator's HTTP client/adapters then attach:

```http
Authorization: Bearer <SOLVENT_API_KEY>
X-API-Key: <CONDUCTOR_API_KEY>
```

This is preferable to trying to inject headers externally from `cmd/coordinator/main.go`, because the configured clients are the components actually making the HTTP requests.

The important architectural constraint is that these fields are **transport credentials**, not business-policy inputs. They must never participate in promotion, discharge, authorization, belief evaluation, etc.

So structurally:

```go
type Config struct {
    SolventURL        string
    SolventAPIKey     string

    ConductorURL      string
    ConductorAPIKey   string

    // existing fields...
}
```

and:

```text
cmd/coordinator/main.go
        ↓
Coordinator.Config
        ↓
Solvent/Conductor HTTP clients
        ↓
auth headers
```

The environment variables are already specified in the plan's coordinator startup configuration. 

This keeps the launcher thin while ensuring the production HTTP path actually authenticates.

---

### 3. P1-3 N+1 lookups: **acceptable for the POC, but fix the implementation to deduplicate**

An N+1 HTTP lookup is acceptable for this POC because the purpose of P1-3 is correctness—preventing cross-scenario contamination—not optimization.

However, **do not implement the pseudocode literally**. The proposed structure in the plan nests the edge scan inside the belief loop, which can accidentally become:

```text
number of packet beliefs × number of packet edges
```

instead of the intended one lookup per referenced belief.

Use a unique-reference set:

```text
collect all canonical belief IDs from:
    edge.from_ref
    edge.to_ref
    evidence.belief_ref

deduplicate IDs

for each unique ID:
    GetBelief(id)
    verify scenario_id == packet.ScenarioID
```

So the actual cost becomes:

```text
O(U) HTTP calls
```

where `U` is the number of unique canonical belief references.

That stays within the plan's fail-closed requirement: validation happens before any Solvent/Conductor persistence, and any lookup failure or scenario mismatch rejects the packet. 

No batching endpoint, caching layer, or broader optimization is justified for this fixes phase.

### Final answers for the coding agent

```text
1. Do NOT create stubs.
   Modify setup:build to build only verifier and coordinator, which are actually
   created by this fixes plan. Do not add dead argus-mcp/trust-ui build references.

2. YES, add SolventAPIKey and ConductorAPIKey to Coordinator.Config.
   Treat them strictly as HTTP transport credentials. cmd/coordinator/main.go reads
   the environment and supplies them; the Solvent/Conductor clients attach the
   appropriate auth headers.

3. YES, N+1 is acceptable for the POC.
   But deduplicate canonical belief UUIDs first. Implement O(U) HTTP lookups,
   not the nested-loop O(B×E) shape shown in the draft pseudocode.
```

That resolves all three questions without expanding the architecture or scope of the fixes.
