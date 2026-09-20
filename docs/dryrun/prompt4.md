Yes. This is now the **right implementation plan** for the first POC. It has successfully collapsed the scope around the actual experiment rather than building retrieval infrastructure prematurely. The locked decisions are internally coherent: static repository corpus, exactly two MCP tools, dedicated `argus mcp`, `bmist@1.1.0`, minimal seed, no Solvent kernel changes, and no new services/databases. 

I would make only **three final edits before handing it to the coding agent**.

### 1. Make pack directory/version handling explicit

There is one ambiguity here:

> `domain-pack/bmist/v1/pack.json` — updated content (or new version directory)

while immediately afterward:

> Do NOT modify `domain-pack/bmist/v1/pack.json` content for 1.0.0. 

The coding agent should not have to decide this.

Change it to:

```text
Create a new immutable pack artifact for bmist@1.1.0.
Do not modify the existing bmist@1.0.0 artifact.

Use whatever directory/layout the current registry actually requires to represent multiple versions. First inspect the registry implementation; then add 1.1.0 alongside 1.0.0.

Both bmist@1.0.0 and bmist@1.1.0 must remain loadable.
```

That preserves historical reproducibility without prescribing a directory convention that may not match the repository.

### 2. Do not blindly update all tests to 1.1.0

The current wording says:

> Existing tests updated to reference 1.1.0. 

That could accidentally eliminate the compatibility test for `1.0.0`.

Replace it with:

```text
Update current-path tests to use bmist@1.1.0 where they test current behavior.

Retain or add a compatibility test proving bmist@1.0.0 remains loadable and semantically unchanged.

Do not mechanically replace every 1.0.0 reference.
```

That is important because the whole reason for versioning is preserved historical state.

### 3. Make the MCP requirement slightly less implementation-prescriptive

The plan says:

> newline-delimited JSON-RPC, ~100–150 lines, no SDK. 

The architectural requirement is correct, but the exact implementation should remain subordinate to the actual repository.

I'd change it to:

```text
Inspect the existing MCP adapter and any protocol implementation already present.

Implement the minimum standards-compliant stdio transport required by the current adapter.

Prefer the standard library if no existing MCP transport dependency is present.

Do not add an MCP SDK merely for convenience.

stdout is protocol-only; logs go to stderr.
```

This prevents the agent from writing a second protocol abstraction unnecessarily.

---

## One thing I especially approve

The revised plan correctly makes repository filesystem access an **agent capability**, not an ARGUS capability:

```text
OpenCode
  ├── argus.get_context
  ├── argus.submit_packet
  └── repository filesystem
        └── docs/corpus/*.md
```

while ARGUS remains the authority boundary. That is exactly the separation you want.

The role cards explicitly establish that distinction. 

And the final dry run is now wonderfully small:

```text
argus serve
    ↓
minimal seed
    ↓
fresh Work Agent
    ↓
get_context
    ↓
read 7 files
    ↓
research
    ↓
submit_packet
    ↓
Insights / Debts
    ↓
fresh Adversarial Agent
    ↓
get_context
    ↓
read 7 files
    ↓
attack
    ↓
submit_packet
    ↓
human adjudication
    ↓
promotion gate
```

That is the experiment—not corpus engineering. The plan explicitly preserves that focus. 

### Final assessment

**Approved for implementation after those three wording fixes.**

The most important property is that the plan now has a built-in brake:

> “What concrete failure in the current dry run does this solve?”
> If none, don't build it. 

That is the correct governing rule for this POC.
