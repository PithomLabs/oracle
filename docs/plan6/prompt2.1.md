**Choose option 1: add a minimal Solvent REST endpoint.**

This is exactly the kind of **Growth Gate exception** we already anticipated: extend Solvent's interface, not its schema or authority model. The `belief_edge` table already exists and `RetractCascade` already depends on it, so the missing piece is simply the lawful write interface.

### Freeze the decision

Add:

```http
POST /v1/beliefs/{parent_id}/edges
```

Request:

```json
{
  "child_id": "<uuid>",
  "kind": "derives"
}
```

Allowed:

```text
derives
contradicts
```

Solvent performs the structural checks:

```text
parent exists
child exists
parent != child
kind is valid
edge is unique
```

Then inserts into `belief_edge`.

No schema change.

No Coordinator SQL.

No Coordinator import of Solvent kernel.

No general edge-management subsystem.

### Why this is the correct boundary

The architecture is:

```text
OpenCode
   ↓
argus.submit_packet
   ↓
Coordinator
   ↓ HTTP
Solvent
   ↓
belief_edge
```

The Coordinator remains orchestration.

Solvent remains the authority over epistemic state.

That is materially better than importing the Solvent kernel directly into Coordinator, which would collapse the repo/interface boundary.

### Do not add an edge MCP tool

This is important.

Agents still get only:

```text
argus.get_context
argus.submit_packet
```

They do **not** get:

```text
solvent_create_edge
```

The Work/Adversarial packet expresses the intended edge; Coordinator validates the packet and calls the Solvent REST endpoint.

So:

```text
Agent proposes edge
        ↓
EBP packet
        ↓
Coordinator
        ↓
Solvent validates + persists edge
```

The agent never owns the edge mutation.

### This also fixes the Phase 8 dry run

Now we can genuinely exercise:

```text
Work Agent
  L1 ──derives──> G0
  L2 ──derives──> G0

Adversarial Agent
  A ──contradicts──> G0

Human
  RETRACT G0

Solvent
  RetractCascade()
      ↓
  dependent state affected

Conductor
  terminal cancelled task

Insights
  structurally derives Dead End
```

So H1 and H3 are no longer decorative plan items.

### Required tests

Add Solvent API tests:

```text
Create derives edge
Create contradicts edge
Reject unknown parent
Reject unknown child
Reject self-edge
Reject invalid kind
Reject duplicate edge
```

Coordinator integration tests:

```text
work packet → derives edges persisted
adversarial packet → contradicts edge persisted
```

And the end-to-end Phase 8 test:

```text
contradicts edge
→ retract
→ cascade
→ dead-end derivation
```

### One important implementation rule

Before coding, **verify the existing Growth Gate contract against the actual Solvent code**. The endpoint should be added as the smallest interface necessary, matching the already frozen edge semantics. Do not invent a broader generic graph API.

So the coding-agent answer should effectively be:

> **Select #1. Add `POST /v1/beliefs/{parent_id}/edges` to Solvent as the Phase 8 Growth Gate exception. Keep edge creation inside Solvent, expose it through REST only, do not expose it through MCP to agents, and do not allow Coordinator direct DB/kernel access.**
