Agreed. We should stop treating this as a writing exercise and turn it into a **small engineering/design work package**.

## Deliverables

I would keep it to **6 concrete artifacts**:

| # | Deliverable                            | Purpose                                                                          |
| - | -------------------------------------- | -------------------------------------------------------------------------------- |
| 1 | **Loop Engineering Requirements v0.1** | Freeze what the system must/must not do                                          |
| 2 | **Workflow Contract v0.1**             | Define the actual loop, branches, handoffs, and human intervention               |
| 3 | **Role/Boundary Matrix**               | Make Agent/Conductor/Solvent/Executor/Human responsibilities unambiguous         |
| 4 | **Conformance Test Plan**              | Define deterministic protocol tests + agent behavioral tests                     |
| 5 | **Reference BM-IST Scenario Set**      | 3–5 realistic workflows used to validate the contract                            |
| 6 | **POC Implementation Checklist**       | Exact work required to run the first experiment using existing Conductor/Solvent |

That's it.

**No new runtime. No new database. No new UI. No SDK initially.**

---

# How we start

### Step 1 — Freeze requirements

Before writing more prose, create a **one-page requirements document**.

The initial requirements are already largely known:

```text
REQUIRED

R1  Domain-agnostic
R2  Tool/runtime-agnostic
R3  Four roles remain distinct
R4  Autonomous mode
R5  Semi-autonomous mode
R6  Human intervention at any stage
R7  Human review at any stage
R8  Consequential actions cannot bypass authority
R9  Authorization ≠ execution
R10 Activity ≠ execution proof
R11 Domain truth remains outside infrastructure
R12 Long-running execution supported
R13 Failure/denial/retry supported
R14 No new persistent workflow authority state
R15 Existing Conductor UI is preferred
R16 No new runtime/framework unless later proven necessary
```

Then add explicit **NON-GOALS**:

```text
NG1 No workflow engine
NG2 No new orchestration platform
NG3 No new authorization system
NG4 No second task system
NG5 No workflow-specific database
NG6 No workflow-specific UI
NG7 No persisted workflow Intent
NG8 No domain-specific logic in Conductor
```

This becomes our first design gate.

---

# Step 2 — Build the state/transition model

Before implementation, we need one authoritative diagram/table.

Something like:

```text
DISCOVER
   ↓
FORMULATE
   ↓
ASSIGN / CLAIM
   ↓
WORK
   │
   ├── ordinary → REVIEW
   │
   └── consequential proposal
                 ↓
              SOLVENT
                 ↓
              EXECUTOR
                 ↓
               RESULT
                 ↓
             INTERPRET
                 ↓
               REVIEW
                 ↓
             NEXT WORK
```

Then add the human overlay:

```text
HUMAN
 ↕
EVERY STATE
```

For each transition we record:

```text
current state
actor
action
required evidence
destination
who may intervene
what intervention is allowed
```

This is the **Workflow Contract**.

No implementation yet.

---

# Step 3 — Create the Role/Boundary Matrix

This should be a table, not prose.

Example:

| Concern          | Agent      | Conductor | Solvent     | Executor         | Human                                      |
| ---------------- | ---------- | --------- | ----------- | ---------------- | ------------------------------------------ |
| Reasoning        | Owner      | —         | —           | —                | Can intervene                              |
| Task lifecycle   | —          | Owner     | —           | —                | Can intervene                              |
| Authority        | —          | —         | Owner       | Enforces         | May exercise through Solvent               |
| External effect  | Requests   | —         | Authorizes  | Owner            | May stop/authorize through proper boundary |
| Domain truth     | Interprets | —         | —           | —                | Reviewer/verifier                          |
| Activity history | Produces   | Records   | Own records | Produces outcome | Can inspect                                |

We'll deliberately look for contradictions here before touching code.

---

# Step 4 — Define the conformance tests

Not a huge test suite.

Start with **10–15 tests**.

### Protocol mechanics

Examples:

```text
T01 ordinary work does not require Solvent
T02 consequential request cannot bypass authority
T03 denied authorization cannot execute
T04 authorization does not imply execution
T05 execution failure is distinct from authorization failure
T06 long-running execution does not block workflow globally
T07 human can pause workflow
T08 human can reject work
T09 human can intervene before consequence
T10 human can intervene during execution when supported
T11 no workflow Intent is persisted
T12 Conductor does not become an authority engine
```

These are deterministic.

### Agent behavior

Then perhaps:

```text
T13 GPT behavioral adherence
T14 Claude behavioral adherence
T15 scripted-client adherence
```

Run each multiple times and measure adherence rather than pretending stochastic behavior is deterministic conformance.

---

# Step 5 — Build only 3 BM-IST scenarios initially

Don't create the whole BM-IST project.

Pick three scenarios that deliberately exercise different parts of the architecture.

### Scenario A — Ordinary research

```text
Agent discovers research question
→ Conductor task
→ Agent performs local work
→ artifact
→ review
```

No Solvent.

### Scenario B — Consequential shared computation

```text
Agent proposes expensive/shared computation
→ authority boundary
→ Solvent
→ Executor
→ result
→ interpretation
```

This proves the consequential branch.

### Scenario C — Human intervention

```text
Agent working
→ human intervenes
→ task redirected/rejected/paused
→ agent resumes
→ review
```

This proves semi-autonomous operation.

Then we can add a mid-flight execution case.

---

# Step 6 — Run the first POC

Only after the above exists:

```text
Existing Conductor
        +
Existing Solvent
        +
One agent
        +
One executor
        +
BM-IST scenario A/B/C
```

Then ask:

> **Where did the existing architecture fail to express the required workflow?**

Only **that evidence** justifies a new primitive, adapter, API change, or other implementation.

This is our Growth Gate applied to Loop Engineering itself.

---

# The actual work sequence

I recommend we execute in this order:

```text
1. Requirements v0.1
        ↓
2. Workflow Contract v0.1
        ↓
3. Role/Boundary Matrix
        ↓
4. 10–15 Conformance Tests
        ↓
5. 3 BM-IST Scenarios
        ↓
6. Run POC
        ↓
7. Fix only demonstrated gaps
```

No more broad architectural essays until something actually fails.

## First deliverable

**I suggest we create Deliverable #1 now: `Loop Engineering Requirements v0.1`.**

It should be **one or two pages maximum**, with the requirements, non-goals, and acceptance criteria. That gives us the foundation against which every later design decision can be tested.
