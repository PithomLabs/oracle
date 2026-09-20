Agreed. **We should stop reviewing the castle and build the foundation.**

Plan 10r5 has reached the point where continued adversarial refinement is producing diminishing returns. The document already freezes the single-process architecture and enumerates 13 corrective items; the remaining debate is largely about making local infrastructure more theoretically robust rather than making the ARGUS thesis demonstrably sound. 

## Foundation freeze

I would now reduce the design to these invariants and **treat everything else as implementation detail or future work**.

### The system

```text
OpenCode
   │
   │ MCP stdio
   ▼
┌──────────────────────────────┐
│ ARGUS — one Go process       │
│                              │
│ MCP                          │
│ Application                  │
│ Solvent / epistemic kernel   │
│ Work state                   │
│ Domain Pack                  │
│ Verifier                     │
│ Trust UI                     │
└──────────────┬───────────────┘
               │
               ▼
          CockroachDB
```

One process. One database. One HTTP port. MCP is stdio.

That is the architecture we build.

---

# The eight things that actually matter

### 1. Agent capability boundary

Exactly:

```text
argus.get_context
argus.submit_packet
```

Agents can produce work.

Agents cannot:

```text
promote
retract
discharge
authorize
```

This is the core agency/authority separation.

### 2. Solvent remains the authority ledger

Keep:

```text
belief
evidence
debt
edges
promotion
retraction
action intent
```

The application may call the kernel internally. We do **not** recreate REST boundaries merely to preserve an old architecture.

### 3. Human authority is authenticated and server-derived

Trust UI writes require:

```text
operator token
→ authenticated request
→ server-derived principal
→ application
→ authority operation
```

Never accept `PrincipalID` from the browser.

Never accept authorization from the agent.

That is foundational.

### 4. Discharge has one mechanical gate

For the POC:

```text
debt exists
      ↓
Pack defines retirement rule
      ↓
evidence class matches
      ↓
qualifying persisted evidence exists
      ↓
human discharge
```

Do not build a rule interpreter.

The natural-language `Rule` remains explanatory metadata, exactly as the plan already specifies. 

### 5. Idempotency must be real

Keep the full-packet content hash and deterministic entity IDs.

That protects the basic property:

```text
same packet twice
    ≠ duplicate state

different packet
    ≠ silently discarded
```

This is worth fixing now.

### 6. UNKNOWN must remain UNKNOWN

`get_context` must distinguish:

```text
empty database
```

from:

```text
database unavailable
```

The per-section availability representation is sufficient. Do not elaborate it further.

### 7. Domain-pack boundary stays clean

Keep:

```text
generic domain-pack
        ↑
      BM-IST
```

and registration at the composition root.

But **do not build dynamic domain infrastructure yet**.

For the POC there is one pack:

```text
bmist@1.0.0
```

The server chooses it.

The client does not.

We do **not** need a `scenario_pack` database table yet. When a second real domain exists, that becomes a justified requirement rather than speculative infrastructure.

### 8. Developer lifecycle must simply work

The essential lifecycle is:

```text
task dev
    ↓
start CRDB
    ↓
apply migrations
    ↓
start ARGUS
```

and:

```text
task fresh
    ↓
destructive reset
    ↓
start ARGUS
```

and:

```text
task down
    ↓
stop ARGUS + CRDB
```

Use fixed local ports:

```text
CRDB SQL     26257
CRDB Admin   8081
ARGUS HTTP   8080
```

Fail clearly when occupied.

**Do not build a port-allocation subsystem.**

There is almost no product value in solving arbitrary multi-stack port allocation for this POC.

A simple workspace startup lock is enough to prevent accidentally launching two stacks against the same `.cockroach-data`.

---

# What we explicitly do NOT build now

This is the important part.

### No dynamic port allocator

Remove:

```text
resolve-ports
runtime port negotiation
port re-resolution
EADDRINUSE retry machinery
```

Fixed ports + clear failure is sufficient.

### No elaborate stale-process recovery

PID files + `task down` + startup lock.

Enough.

### No operator-attestation wizard

The boundary matters:

```text
agents cannot create operator_asserted evidence
```

For the POC, operator evidence can be seeded through a trusted human/demo setup path.

We don't need to build an entire evidence-authoring product to prove the authority model.

### No real scenario→pack persistence yet

One pack exists.

Use the server-selected BM-IST pack.

When a second domain actually arrives, introduce persistent scenario binding.

### No sophisticated semver

There is one pack/verifier version.

Exact version comparison is sufficient for the POC.

### No pagination

The dataset is tiny.

### No elaborate audit framework

Existing authoritative state and activity/refusal records are enough to demonstrate the thesis.

### No second domain

The portability claim is:

> A second domain should require a new pack rather than rewriting the epistemic core.

We do not need to build that second domain merely to prove the sentence.

---

# The resulting foundation

The architecture becomes:

```text
                 AGENCY
                OpenCode
                   │
                   │ 2 MCP tools
                   ▼
        ┌───────────────────────┐
        │        ARGUS          │
        │                       │
        │ Work ───────┐        │
        │             │        │
        │ Epistemic ──┼─ Solvent
        │             │        │
        │ Domain Pack │        │
        │ Verifier    │        │
        │ Trust UI    │        │
        └─────────────┼────────┘
                      │
                  CockroachDB
```

And the governing invariant remains:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

That is the thing we are actually trying to prove.

Everything else is scaffolding.

## Decision

**Stop the review cycle. Freeze the architecture at this level and implement.**

Plan 10r5 should not receive another comprehensive adversarial rewrite. It should be converted into a **minimal implementation checklist**, with the nonessential infrastructure explicitly deferred rather than perfected.

The correct next milestone is no longer “a flawless plan.”

It is:

```text
one process
+
one database
+
two agent tools
+
real authority gates
+
real human decision path
+
real verifier artifact
+
real Trust UI
+
passing end-to-end demonstration
```

Then we learn from the running system rather than continuing to optimize a document about a system that does not yet exist.
