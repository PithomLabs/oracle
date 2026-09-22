I reviewed Plan 1.0. The architecture is right, but I would **not execute it unchanged**. There are two important implementation issues and a couple of smaller refinements.

### 1. Critical: `packet_submission.packet_id` must not be UUID

This plan currently proposes:

```sql
packet_id UUID PRIMARY KEY
```

and:

```sql
VALUES ($1::UUID, ...)
```

But your actual packet IDs are strings such as:

```text
deliverable-A-substrate-v1-001
```

The dry run already demonstrated this exact style of packet ID.

So this migration would fail for the very packets you are submitting.

Change it to:

```sql
packet_id STRING PRIMARY KEY
```

and:

```sql
VALUES ($1, ...)
```

The packet schema already treats `packet_id` as a string, so the persistence layer should preserve that contract. 

This is the most important correction.

---

### 2. `task_id UUID` needs an explicit contract

The plan has:

```sql
task_id UUID
```

but then does:

```go
taskRef := sql.NullString{String: pkt.TaskRef, Valid: pkt.TaskRef != ""}
```

and passes that string directly into the UUID column. 

Before implementing, the agent needs to inspect the actual type/semantics of `Packet.TaskRef`.

The clean rule should be:

```text
TaskRef is a UUID
    → store as UUID

TaskRef absent
    → NULL

TaskRef is some other identifier
    → do not silently cast it
```

I would explicitly tell the agent to resolve this before migration/code.

---

### 3. Scope `GetAllSubmissions` to the current research scope

The plan currently says:

```sql
SELECT ...
FROM packet_submission
ORDER BY submitted_at DESC
```

which means the new dashboard would display submissions from **all scenarios**.

That conflicts with the dashboard semantics you established earlier: Insights is a dashboard for the current research scope, not a global ledger browser.

Make it:

```sql
WHERE scenario_id = ?
ORDER BY submitted_at DESC
```

using the same current-scenario mechanism as the existing dashboard.

The table can remain globally stored; the projection should be scoped.

---

### 4. Keep the packet submission record truly append-only

The conceptual design is correct: `packet_submission` is provenance, not authority. 

I'd make the invariant explicit:

```text
packet_id → one immutable provenance record
```

and ensure no application code exposes:

```text
UPDATE packet_submission
DELETE packet_submission
```

The `ON CONFLICT DO NOTHING` is good for idempotency, but it does not by itself establish database-level immutability.

You don't need triggers or another elaborate subsystem for this POC; simply ensure there is no application update/delete path.

---

### 5. The protocol design is good

The four-field identity is exactly right:

```json
{
  "id": "work-001",
  "role": "work",
  "harness": "OpenCode",
  "model": "MiMo-V2.5"
}
```

The plan correctly makes all four required, human-supplied, and non-authoritative. 

And this is particularly good:

```text
agent.role == packet.role
```

because otherwise you could have an adversarial packet claiming `role=work` in its embedded identity.

---

### 6. One UI change I would make

Don't make the submission table the primary representation of the task's `Agent` column.

The task's:

```text
current_agent
```

still means task claiming.

The new:

```text
packet_submission.agent_id / harness / model
```

means packet provenance.

So I'd actually label the new UI section:

> **Recent Agent Submissions**

rather than trying to repurpose `CurrentAgent`.

That preserves the semantic distinction the audit uncovered.

You can leave the task table's `Agent` column as:

```text
—
```

until a task has actually been claimed.

Then separately show:

```text
Recent Agent Submissions

Packet | Agent | Harness | Model | Role | Time
```

This avoids conflating two different concepts.

---

## The corrected persistence shape

I would use:

```sql
CREATE TABLE IF NOT EXISTS packet_submission (
    packet_id      STRING PRIMARY KEY,
    scenario_id    UUID NOT NULL,
    task_id        UUID,
    agent_id       STRING NOT NULL,
    role           STRING NOT NULL,
    harness        STRING NOT NULL,
    model          STRING NOT NULL,
    content_sha256 STRING NOT NULL,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

Then:

```text
packet
   |
   +-- packet_id
   +-- agent.id
   +-- agent.role
   +-- agent.harness
   +-- agent.model
   |
   v
packet_submission
```

And:

```text
conductor_task.current_agent
    = operational task claimant

packet_submission.agent_*
    = immutable provenance of submitted research
```

That's the clean distinction.

### Final disposition

**Approve the plan after these three required corrections:**

1. `packet_submission.packet_id` → `STRING`, not UUID.
2. Resolve and explicitly validate `TaskRef` → `task_id UUID`.
3. Scope submission projection to the current research scenario.

Everything else in the plan is aligned with the architecture you're building, particularly the separation between packet provenance and task claiming, and the shared protocol for Work and Adversarial agents. 

I would have the coding agent make those corrections **before implementation**, then proceed.
