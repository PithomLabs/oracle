The coding agent’s observations are good, and I would **address them now before issuing the final implementation prompt**. The important point is that none of them invalidate the simplified architecture; they mostly tighten the contracts.

### 1. Pack identity: use `bmist@1.1.0`

I agree with the agent here.

My earlier recommendation of `bmist-as@0.1.0` is unnecessarily disruptive for this POC. The existing code already has `bmist@1.0.0` embedded in multiple places. The new debt vocabulary is additive, and the current methodology is still the same BM–IST–AS program.

Use:

```text
bmist@1.1.0
```

But make one invariant explicit:

> `bmist@1.0.0` remains immutable and valid for historical records; `bmist@1.1.0` is the current pack.

Do **not** modify the contents of `1.0.0` in place.

The fact that the pack is embedded does not, by itself, solve semantic versioning; an old binary can reproduce old behavior, but a new binary must not reinterpret an old pack reference.

So the coding agent should update the hardcoded current references from `1.0.0` → `1.1.0`, while retaining compatibility with existing `1.0.0` records.

### 2. Provenance-boundary test: don't overbuild it

The agent's observation is correct: the existing discharge machinery already prevents a plain document read from becoming retirement evidence.

The key distinction is:

```text
read file
    ≠
persist evidence
    ≠
satisfy retirement rule
    ≠
human discharge
```

I would **not create a new provenance subsystem** just to test this.

But I would add one focused test or dry-run assertion:

```text
An agent may read adv_review*.md and reason from it,
but cannot use the document itself as retirement evidence.
```

The existing evidence/promotion machinery should be the enforcement mechanism.

So this is a **small verification enhancement**, not a new feature.

### 3. `governance_ref` is genuinely Phase-0 blocking

This is the one observation I would elevate.

The agent says the schema/application mismatch could cause seed insertion to fail. That means the coding agent must resolve this **before writing the seed implementation**.

The prompt should say:

> Inspect the actual SQL type, application model, serialization path, and persistence path for `governance_ref`. Establish the canonical representation and make the seed use the existing application contract. Do not work around the mismatch with ad hoc JSON/UUID conversions.

Better still, prefer an existing application-level constructor/persistence path over raw SQL if one already exists. That prevents the seed from creating records that the normal application path could not understand.

### 4. OpenCode filesystem access is an agent capability, not ARGUS functionality

Correct.

We should explicitly state in the role cards:

```text
You have repository filesystem access.
Read the seven research documents directly from docs/corpus/.
Do not attempt to obtain them through ARGUS.
```

And the acceptance criterion should change from:

> "ARGUS can provide all seven documents"

to:

> "In the intended local OpenCode runtime, the Work and Adversarial agents can read the seven repository documents."

That's an important architectural distinction.

ARGUS does not need to become a file server.

### 5. ADD_DEBT remains deferred

Agreed completely.

There is no evidence from this architecture that ADD_DEBT is needed for the first experiment. The agent already has a valid mechanism for submitting new beliefs and associated obligations, while contradictions are represented through the existing packet/edge mechanism.

Do not add it.

---

## The revised decision set

I would now lock these decisions:

| Area                   | Decision                                |
| ---------------------- | --------------------------------------- |
| Research corpus        | Static repository files                 |
| Corpus DB ingestion    | **Removed**                             |
| Embeddings             | **Removed**                             |
| Vector search          | **Removed**                             |
| `get_context(query)`   | **Removed**                             |
| MCP tools              | Exactly 2                               |
| MCP transport          | Dedicated `argus mcp`                   |
| Domain pack            | **`bmist@1.1.0`**                       |
| `bmist@1.0.0`          | Immutable historical version            |
| Seed                   | 1 project + 1 task + 1 belief + 8 debts |
| Debt source            | Domain pack, not hardcoded seed list    |
| `governance_ref`       | Resolve before seed implementation      |
| Provenance enforcement | Existing evidence/retirement machinery  |
| Provenance test        | One focused verification, no subsystem  |
| ADD_DEBT               | Deferred                                |
| Solvent kernel         | No changes                              |
| New service            | None                                    |
| New DB                 | None                                    |

This makes the POC substantially cleaner.

The actual experiment becomes:

```text
             STATIC BACKGROUND
          ┌─────────────────────┐
          │ 7 research files    │
          └──────────┬──────────┘
                     │
                     v
Fresh Work Agent ──> ARGUS get_context
                     │
                     v
                research work
                     │
                     v
                submit_packet
                     │
                     v
                  Solvent
                     │
                     v
             Insights / Debts
                     │
                     v
Fresh Adversarial Agent
                     │
              get_context
                     │
              read same 7 files
                     │
                     v
              challenge work
                     │
                     v
                submit_packet
```

That is the right first experiment. **Do not let the retrieval subsystem sneak back in before this loop has been demonstrated.**

The only substantive changes I would make to the prompt I gave you are therefore: change the pack decision to `bmist@1.1.0`, make `governance_ref` a hard Phase-0 blocker, explicitly describe repository filesystem access as an OpenCode capability, and downgrade the provenance requirement to a focused verification rather than a new enforcement mechanism.
