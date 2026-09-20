I read `reference-loop-plan.md`. The architecture and implementation direction are sound, and the five open questions can now be resolved without expanding scope. 

## Answers to the five questions

### 1. Test repository

**Create `pithomlabs/reference-loop-test`.**

Use a dedicated repository specifically for the reference-loop integration proof.

It should contain the minimal workflow:

```text
.github/workflows/ref-loop.yml
```

The workflow should perform a harmless, deterministic external action, such as creating/updating an artifact or writing `result.txt`, as already proposed in the plan. 

Keep it completely separate from production repositories.

---

### 2. Conductor database

**Use plain SQLite with `modernc.org/sqlite`, file-backed for the test run.**

Do **not** use CockroachDB.

Prefer:

```text
reference-loop/.tmp/conductor.db
```

or an automatically generated temporary database path per run.

I would avoid `:memory:` for the reference loop because we want the evidence to remain inspectable after the process exits. A temporary file gives us deterministic local execution while keeping the environment trivial.

The reference harness should own the test database lifecycle:

```text
create
→ initialize schema using existing Conductor schema/migrations
→ run scenario
→ inspect evidence
→ optionally preserve database
→ cleanup
```

Do not modify Conductor to add SQLite-specific application behavior merely for the harness.

---

### 3. Solvent database

**Same decision: plain SQLite using `modernc.org/sqlite`.**

Use a separate database from Conductor:

```text
reference-loop/.tmp/solvent.db
```

Do **not** use CockroachDB.

The separation is important because we are testing the architectural boundary. The harness should not accidentally make Conductor and Solvent share a database.

Conceptually:

```text
Conductor
   │
   └── conductor.db

Solvent
   │
   └── solvent.db
```

The test harness can inspect both stores afterward, but neither system gets a cross-owned database.

This is particularly consistent with the evidence requirement that participant-owned evidence remain in the participant's own store rather than becoming a new centralized ledger. 

---

### 4. Correlation ID scheme

I would **not** use only:

```text
task_id:scenario_id:intent_id
```

because `intent_id` does not exist at the beginning of the loop.

Use a harness-generated immutable:

```text
run_id
```

as the top-level correlation identifier.

Then keep the participant identifiers as separate fields:

```text
run_id
scenario_id
task_id
intent_id
operation_id
```

For example:

```text
run_id      = 01J...
scenario_id = happy_path
task_id     = conductor-task-123
intent_id   = solvent-intent-456
operation_id = deploy:owner/repo:ref-loop.yml:main
```

The effective correlation model becomes:

```text
run_id
 ├── scenario_id
 ├── task_id
 ├── intent_id
 └── operation_id
```

This is cleaner than constructing a mutable composite string whose components appear at different lifecycle stages.

The existing `EvidenceRecord` can therefore be adjusted to contain these identifiers separately. The `CorrelationID` field can remain as a stable `run_id`, rather than encoding every relationship into one string. The plan already requires operation ID, declaration version, timestamp, and correlation information. 

---

### 5. First scenario

**Start with `happy_path` only. Do not scaffold all scenarios as implementation stubs yet.**

First prove:

```text
TASK_CREATED
→ TASK_CLAIMED
→ PROPOSAL_CREATED(X)
→ AUTHORIZED(X)
→ EXECUTION_ATTEMPTED(X)
→ EFFECT_CONFIRMED(X)
→ RESULT_OBSERVED
→ CONDUCTOR_UPDATED
→ AGENT_CONTINUED
```

using the real component interfaces and SQLite test databases.

Once that works and the evidence can be reconstructed, add the adversarial scenarios **one at a time**.

This preserves the pivot's core discipline: first establish a real working loop, then attack it with observed behavior rather than building a large scaffold in advance. The plan itself puts adversarial scenarios after the happy path. 

---

# One additional implementation decision

There is one change I would make to the attached plan before handing it to the coding agent.

The plan currently says:

```text
Evidence collector polls each system's audit/ledger
```

That's fine for evidence inspection, but **the reference loop itself should not depend on polling to function**.

The execution path should be:

```text
Agent
→ Conductor MCP
→ Solvent REST
→ Executor
→ GitHub
```

and the evidence collector should operate **afterward/asynchronously for reconstruction**.

That keeps:

```text
execution semantics
```

separate from:

```text
evidence collection
```

and prevents the harness from becoming an accidental orchestration layer.

---

# Final implementation decisions

So I would hand the coding agent these explicit answers:

```text
Test repository:
    pithomlabs/reference-loop-test

Conductor DB:
    plain SQLite
    modernc.org/sqlite
    file-backed temporary database

Solvent DB:
    plain SQLite
    modernc.org/sqlite
    separate file-backed temporary database

Correlation:
    immutable run_id
    plus separate scenario_id/task_id/intent_id/operation_id

First scenario:
    happy_path only

Executor:
    RecordingFunc for deterministic local/CI testing
    real GitHub executor for integration proof

Reference operation:
    existing deploy / github_trigger_workflow

Conductor interface:
    MCP

Solvent interface:
    REST

Authoritative external result:
    GitHub workflow run state / GitHub API
```

The final architecture under test is therefore:

```text
             MCP
Agent ───────────────→ Conductor
                         │
                         │ task/activity
                         │
Agent ───────────────→ Solvent
             REST         │
                          ↓
                       Executor
                          │
                          ↓
                    GitHub Actions
                          │
                          ↓
                 GitHub = SOR
                          │
                          ↓
                    Result/evidence
```

with **two isolated SQLite databases** used only to make the local reference environment reproducible and inspectable.

The coding agent should implement **only `happy_path` first**, get that demonstrably working, and stop for review before moving to the adversarial scenarios.
