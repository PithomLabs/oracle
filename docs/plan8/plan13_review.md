Plan 13 is **almost right**, but I would make four corrections before implementation. The important part is that the plan correctly targets the real blocker: the production MCP path is bypassing the canonical packet validator. 

### 1. `ValidatePacket()` must fail closed

This is the most important correction.

The plan proposes:

```go
func (a *App) ValidatePacket(ctx context.Context, pkt *packetv1.Packet) error {
    if a.packRegistry == nil {
        return nil
    }
    return packetv1.Validate(pkt, a.packRegistry)
}
```

That reintroduces the exact class of problem you are trying to remove.

If the registry is missing, the validator is silently disabled.

It must be:

```go
func (a *App) ValidatePacket(ctx context.Context, pkt *packetv1.Packet) error {
    if a.packRegistry == nil {
        return fmt.Errorf("packet validation unavailable: pack registry is not configured")
    }
    return packetv1.Validate(pkt, a.packRegistry)
}
```

The agent-facing boundary should be:

```text
validator available → validate
validator unavailable → reject
```

not:

```text
validator unavailable → accept
```

### 2. Add one real MCP-path validation test

The current plan has unit tests for `packet/v1` and `Persist()` integration tests. 

But the whole P0 finding is:

> the **production MCP adapter** does not invoke the validator.

So there must be at least one test through:

```text
argus.submit_packet
    → HandleTool
    → production MCP adapter
    → rejection
```

For example:

```text
TestMCPRejectsMissingAgentID
TestMCPRejectsRoleMismatch
TestMCPRejectsDuplicateLocalID
TestMCPRejectsInvalidEdgeKind
```

You don't necessarily need four separate tests; a focused table-driven MCP validation test is enough.

Without that, you could accidentally reintroduce the bypass later and the tests would still pass.

### 3. Clarify `local_id` uniqueness scope

Phase 4 says:

> "Build combined local_id set across all beliefs/evidence/edges/tasks." 

I actually prefer **global uniqueness within a packet** because your reference grammar is:

```text
local:<id>
```

not:

```text
belief:<id>
evidence:<id>
```

So `b1` should not simultaneously identify a belief and an evidence object.

But then the tests need to cover cross-type collisions too:

```text
belief local_id = "x"
evidence local_id = "x"
→ reject
```

The current plan only lists duplicate IDs within each collection. Add at least one cross-type collision test.

### 4. Add `claim_type` validation to the JSON schema too

The plan adds the enum in Go validation:

```text
derived
accommodated
postulated
```

but the packet schema should express the same contract.

That gives you two useful layers:

```text
JSON schema
    → structural contract

Go validator
    → executable domain validation
```

Otherwise the schema tells consumers less than the actual wire contract.

---

## One thing the plan gets exactly right

The task preflight is appropriately conservative:

```text
packet with tasks
    ↓
verify project exists
    ↓
missing → clear error
```

and:

> **do not auto-create projects**. 

That is the right behavior for a research packet. `submit_packet` should not invent operational infrastructure.

### Revised implementation instruction

Give the coding agent:

```text
Proceed with Plan 13, with these four corrections:

1. ValidatePacket() must fail closed.
   If packRegistry is nil/unavailable, return an error. Never silently skip
   canonical packet validation.

2. Add at least one regression test through the REAL production MCP path
   (`HandleTool` / `argus.submit_packet`) proving malformed packets are
   rejected before any persistence occurs.

3. Make local_id uniqueness explicitly global within a packet because
   references use `local:<id>`. Add at least one cross-type collision test
   (e.g. belief local_id == evidence local_id).

4. Add claim_type enum constraints to packet/v1/schema.json as well as the Go
   validator:
       derived
       accommodated
       postulated

Keep the rest of Plan 13 unchanged.

After implementation:
- go vet ./...
- go test ./...
- go build ./...
- run the production MCP validation tests
- confirm no malformed packet reaches Persist()
```

After this pass, the most important proof is no longer merely `77/77` unit tests. It is:

```text
malformed AI packet
      ↓
argus.submit_packet
      ↓
REJECTED
      ↓
zero database mutation
```

That is the actual trust boundary you're trying to establish.
