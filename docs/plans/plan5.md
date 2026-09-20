Now we stop designing and **run the smallest real proof**.

## Next move: build the reference loop

Use the existing components as they are:

```text
Agent
  ↓
Conductor
  ↓
Solvent
  ↓
Executor
  ↓
External system
  ↓
Result
  ↓
Conductor
  ↓
Agent
```

### Phase 1 — Happy path

Pick **one tiny consequential operation** that is easy to observe and safely reversible.

The goal is not BM-IST yet. The goal is proving the plumbing.

The test should demonstrate:

```text
1. Agent discovers a Conductor task
2. Agent claims it
3. Agent proposes exact operation X
4. Solvent authorizes X
5. Executor receives authorization for X
6. Executor performs X against a real/sandbox external system
7. External system produces authoritative result
8. Result reaches Agent
9. Conductor records coordination-relevant progress
10. Agent can propose the next piece of work
```

At the end, we should be able to reconstruct the entire loop from evidence.

## Phase 2 — Attack the loop

Immediately after the happy path, test the invariants that matter most:

```text
Authorized X → Execute Y        MUST FAIL
Missing authorization           MUST FAIL
Stale authorization             MUST FAIL
Verification unavailable       MUST FAIL CLOSED
Duplicate execution             MUST be controlled
Execution failure               MUST stay failure
Ambiguous result                MUST stay ambiguous
Agent misinterprets result      MUST NOT change authoritative result
Conductor cancellation         MUST NOT silently revoke Solvent authority
```

Do not build new infrastructure to make these tests pass.

When one fails, classify it:

```text
implementation
integration
executor
deployment
genuinely new security property
```

Only the last category is grounds for architectural expansion.

## Phase 3 — Substitution

Once the loop is solid, replace one participant at a time.

First:

```text
Agent A → Agent B
```

Then:

```text
Executor A → Executor B
```

Then:

```text
direct API/client → alternate orchestration mechanism
```

The test is not "does it still work somehow?"

The test is:

> **Does the semantic contract remain unchanged?**

That is the actual Loop Engineering proof.

## What I would build first

I would create a **single reference demo harness**, not a framework.

Something approximately like:

```text
reference-loop/
    README.md
    scenario/
        happy_path
        wrong_operation
        denied
        unavailable
        execution_failure
        ambiguous
        cancellation
    agent/
    conductor/
    solvent/
    executor/
    evidence/
    runbook/
```

But keep the implementation extremely thin. Reuse the existing Conductor and Solvent APIs/MCP surfaces rather than introducing a new abstraction over them.

### The first concrete milestone

One command should be able to run:

```text
run-reference-loop
```

and produce a trace equivalent to:

```text
TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED
→ AUTHORIZED(X)
→ EXECUTION_ATTEMPTED(X)
→ EFFECT_CONFIRMED(X)
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED
```

with enough evidence to prove that each fact belongs to the correct owner.

### What comes after that

**Only after the reference loop works** should we implement the Conformance Test Matrix. At that point, the matrix stops being mostly theoretical: each test maps to an actual observable failure mode in the running system.

That is the pivot in practice:

```text
Before:
specification → architecture → more specification

Now:
working loop → evidence → failures → conformance → generalization
```

The next deliverable I would hand the coding agent is therefore a **reference-loop implementation prompt**, not another specification.
