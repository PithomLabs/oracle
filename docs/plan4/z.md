## prompt

just to clarify about the reference loop workflow, when I say AI/coding agent we only refer to one agent (e.g. OpenCode, Antigravity, Codex) where it has its own mechanism to spawn subagents and it may also has Plan mode and Agent/Code mode, when I enter a prompt the task could be relatively complex where the AI agent could spawn on its own accord several subagents which we do not have control, the question now is not how the AI agent will decompose the task A to subtasks B, C and so on but how Conductor can identify and represent those subtasks and track their state akin to attached DAG.md without using DAG but using the Solvent mechanism. I want a protocol where the AI agent can do the subtasks as ordinary work vs subtasks that have consequential action and thus require Solvent. And we are just talking here of one AI agent but in my vision, there should be no stopping subsequent AI agents X, Y, Z and so on to join the endeavor and do the same workflow as the initial AI agent. That is, the Conductor basically orchestrates as the control plane of the entire state of the tasks using the mechanism of attached DAG.md where AI agents can pick up any task with the appropriate metadata such as its interface without stepping on each other, what do you think? how does it simplify or complicate in light of the latest development now? the goal is to clarify the protocol, semantics and workflow from Agent to Conductor to Solvent and the reference loop workflow itself, elaborate short but concise and organize your itemized elaboration around categories, the other goal is not to become another harness

## answer

Yes. With this clarification, I would **refine the architecture substantially**.

The core idea is no longer “the Agent decomposes Task A.” The Agent may already be an autonomous runtime such as OpenCode, Codex, or Antigravity that internally creates and schedules arbitrary subagents. **Conductor's job is to turn that emergent internal work into shared, persistent, externally visible work state.**

That is a much more interesting role for Conductor.

## 1. The right mental model

Think of the initial prompt as creating a **work objective**, not a fully decomposed task tree.

```text
Human
  ↓
AI Agent / coding runtime
  ↓
Task A
  │
  ├── internally spawns subagent
  ├── internally spawns subagent
  ├── internally plans/replans
  └── discovers additional work
              │
              ▼
          Conductor
      shared work state
```

The Agent's internal orchestration is **opaque and autonomous**.

Conductor does not need to know:

> Why did OpenCode create these three subagents?

It needs to know:

> What work now exists, what depends on what, who is doing it, and what state is it in?

That is a much cleaner boundary.

---

# 2. Conductor becomes the external control plane

This is the key insight.

Conductor should be the **control plane of work state**, not the control plane of the Agent's internal reasoning.

So:

```text
Agent runtime
    = private execution / reasoning plane

Conductor
    = shared work control plane
```

Conductor maintains a persistent representation such as:

```text
Task A
 ├── Task B
 │    ├── B1
 │    └── B2
 ├── Task C
 └── Task D
       └── D1
```

But I would **not call this “the DAG” as the protocol**.

Internally it can be represented by dependency edges, but semantically it is a **work graph with state**.

That distinction matters because DAG.md itself correctly argues that software is too entangled to be represented by one clean DAG; the useful graph is the graph of obligations, dependencies, evidence, and coordination around the software. 

---

# 3. The Agent creates work; Conductor materializes it

This is the protocol I think you're converging toward:

```text
Agent
  ↓
"Here is a new work item I discovered"
  ↓
Conductor
  ↓
persistent Task node
```

Conductor then assigns metadata such as:

```text
task_id
parent_task_id
dependencies
status
assignee
interface
inputs
outputs
operation_type
risk/consequence class
evidence requirements
```

The Agent does not need to ask Conductor:

> “How should I decompose this?”

It says:

> “These are the work items I have decided exist.”

Conductor makes them durable and coordinateable.

---

# 4. Ordinary work vs consequential work

This is where Solvent fits beautifully.

Not every task needs authority.

There should be at least two semantic classes:

```text
ORDINARY WORK

Task
  ↓
Claim
  ↓
Execute
  ↓
Report result
```

versus:

```text
CONSEQUENTIAL WORK

Task
  ↓
Propose exact operation
  ↓
Solvent authorization
  ↓
Executor
  ↓
Effect / result
```

So a task in Conductor could carry something like:

```text
kind = ordinary
```

or:

```text
kind = consequential
operation = deploy:repo:workflow:ref:run_id
authority = solvent
```

Conductor does **not** decide whether something is authoritative.

It merely knows that this task has a consequential execution contract.

---

# 5. Solvent should gate the action, not manage the task

This is the crucial separation.

Conductor owns:

```text
task existence
task dependencies
task assignment
task lifecycle
task coordination
```

Solvent owns:

```text
exact action authorization
exact target/state binding
authority evidence
```

Therefore:

```text
Conductor:
    "Task 27 is ready."

Solvent:
    "Operation X is authorized."

Executor:
    "I performed X."
```

Conductor should never turn:

```text
status = ready
```

into:

```text
authorized = true
```

That preserves:

```text
WORK ≠ AUTHORITY
```

---

# 6. The important change to the Reference Loop

I would now depict the Reference Loop less as a mandatory linear pipeline and more as **two interacting planes**.

### Work plane

```text
Agent
  ↕
Conductor
```

The Agent discovers, creates, claims, updates, and completes work.

### Consequence plane

```text
Agent
  ↓
Solvent
  ↓
Executor
  ↓
External SOR
```

And Conductor observes coordination-relevant results:

```text
External result
      ↓
   Agent
      ↓
 Conductor
```

So the real architecture is closer to:

```text
                 ┌───────────────┐
                 │     Agent     │
                 │ reasoning /   │
                 │ subagents /   │
                 │ replanning    │
                 └──────┬─┬──────┘
                        │ │
                work    │ │ consequence
                        │ │
                        ▼ ▼
                 ┌─────────┐   ┌─────────┐
                 │Conductor│   │ Solvent │
                 │ control │   │authority│
                 │  plane  │   │boundary │
                 └─────────┘   └────┬────┘
                                     │
                                  Executor
                                     │
                                     ▼
                                External SOR
```

That is closer to the architecture we actually want.

---

# 7. This naturally supports Agent X, Y, Z later

This is where your vision becomes much stronger.

Conductor does not care whether work originated from:

```text
Agent A
```

or:

```text
Agent B
```

or:

```text
Agent X
Agent Y
Agent Z
```

It sees a shared work state:

```text
                Conductor
                    │
       ┌────────────┼────────────┐
       ▼            ▼            ▼
     Task B       Task C       Task D
       ▲            ▲            ▲
       │            │            │
    Agent X      Agent Y      Agent Z
```

Agents can pick up work according to:

```text
interface
dependencies
required inputs
required outputs
status
operation class
capability requirements
```

The Agent implementation becomes substitutable.

That is a very powerful consequence.

---

# 8. The interface metadata is important

Your idea that an Agent can pick up a task based on its **interface** is exactly right.

A task shouldn't just say:

```text
"fix authentication"
```

It should expose a machine-readable work contract such as:

```text
Task:
    T-1042

Interface:
    AuthProviderMigration

Inputs:
    provider_config

Outputs:
    migrated_provider

Dependencies:
    T-1038

Work type:
    ordinary

Verification:
    tests/auth/...

Status:
    READY
```

Or:

```text
Task:
    T-1051

Interface:
    GitHubWorkflowDeploy

Inputs:
    repository
    workflow
    ref

Work type:
    consequential

Authority:
    Solvent

Operation:
    deploy:repo:workflow:ref:run_id

Status:
    READY_FOR_AUTHORIZATION
```

Now an Agent doesn't have to understand the entire project.

It can ask:

> Which tasks have an interface I can perform?

That is exactly what lets later Agents X, Y, Z participate safely.

---

# 9. Claiming prevents agents from stepping on each other

Conductor becomes the synchronization mechanism.

For example:

```text
Task B = READY

Agent X → CLAIM(B)
Agent Y → CLAIM(B)
```

Only one succeeds.

The other sees:

```text
B = ACTIVE
owner = Agent X
```

Meanwhile:

```text
Task C = READY
```

can be claimed by Agent Y.

So you get:

```text
          Conductor
        shared state
       /     |      \
      B      C       D
      │      │       │
 Agent X  Agent Y  Agent Z
```

without requiring the Agents to coordinate directly.

That is exactly the kind of control-plane behavior Conductor should own.

---

# 10. This is where the DAG insight becomes useful

The DAG.md thesis maps very well onto this.

The software itself can be tangled:

```text
code
 ↕
state
 ↕
runtime
 ↕
dependencies
 ↕
side effects
```

But Conductor can maintain a cleaner work/proof graph:

```text
Task
  ↓
dependencies
  ↓
verification
  ↓
result
  ↓
unblocked downstream task
```

The attached DAG work explicitly describes the useful graph as one of claims, obligations, dependencies, and evidence rather than a graph of the code itself. 

So you're not actually abandoning the DAG insight.

You're **embedding it into Conductor's persistent work-state model without making DAG itself the product abstraction**.

---

# 11. The protocol becomes surprisingly small

I think the core protocol can eventually be expressed with only a few operations.

### Work protocol

```text
CREATE
DISCOVER
CLAIM
RELEASE
UPDATE
BLOCK
SUBMIT
ACCEPT
REJECT
```

### Consequential protocol

```text
PROPOSE
AUTHORIZE
EXECUTE
OBSERVE
```

With shared metadata:

```text
task_id
run_id
operation_id
dependencies
interface
state
evidence
```

That is much simpler than building a new orchestration runtime.

---

# 12. The biggest architectural danger

The danger is turning Conductor into:

```text
Agent supervisor
+ planner
+ scheduler
+ policy engine
+ workflow engine
+ execution broker
+ authority engine
```

That would destroy the elegance.

Conductor should **not care how the Agent internally decomposes work**.

It only cares about the resulting externally coordinateable state.

Likewise, Conductor should **not perform Solvent's authorization**.

So the clean boundary remains:

```text
Agent:
    creates/reasons about work

Conductor:
    coordinates work

Solvent:
    authorizes consequences

Executor:
    produces effects
```

---

# 13. What “not another harness” means

This is probably the most important constraint going forward.

The Reference Loop should not become a giant program that manually orchestrates every participant:

```text
reference-loop/
    manually call Agent
    manually call Conductor
    manually call Solvent
    manually call Executor
    manually reconstruct everything
```

That would prove only that **our harness knows how to drive the system**.

Instead, the Reference Loop should prove that the participants can operate through their **actual contracts**.

The permanent product architecture should be:

```text
Agent ↔ Conductor
Agent ↔ Solvent
Solvent ↔ Executor
Executor ↔ External SOR
```

with Conductor maintaining the shared work state.

The harness exists only to test those boundaries.

---

# 14. What this changes about the roadmap

I would now refine the roadmap to:

```text
PHASE 1
Reference Loop
    ✓ basic work
    ✓ authorization
    ✓ execution
    ✓ identity
    ✓ denial boundaries

        ↓

PHASE 1.5
Work-graph safety
    → operation mismatch
    → duplicate delivery
    → stale authority
    → verification unavailable
    → explicit task/interface semantics

        ↓

PHASE 2
Real consequential action
    → real external SOR
    → authoritative effect
    → ambiguous outcome

        ↓

PHASE 3
Multi-agent coordination
    → Agent X/Y/Z
    → shared task discovery
    → claim contention
    → dependency resolution
    → substitution

        ↓

PHASE 4
Formal Agent / Conductor protocols
    → Agent Skill
    → Conductor Skill
    → task/interface schema
    → consequential-action contract

        ↓

PHASE 5
Conformance / substitution

        ↓

PHASE 6
BM-IST / real domain workload
```

Notice that **multi-agent coordination is now an explicit validation phase**, rather than assuming the single Agent reference loop proves it.

---

# 15. The most important synthesis

I think your clarified vision gives the architecture a much sharper formulation:

> **The Agent controls intelligence. Conductor controls shared work state. Solvent controls authority. Executor controls effects.**

And:

> **Conductor does not decompose the work; it materializes, coordinates, and persists the work that Agents discover or create.**

And perhaps the most important one:

> **The Reference Loop should validate the protocol between these control boundaries, not become another control layer itself.**

That makes the architecture much simpler.

You don't need to control the internal subagent machinery of OpenCode/Codex/Antigravity.

You don't need Conductor to know how an Agent plans.

You don't need Solvent to know what the project task graph looks like.

You don't need Executor to know why the Agent wants the operation.

Each system only needs to expose the **state and contract required at its boundary**.

That is a much cleaner foundation for eventually allowing an arbitrary succession of Agent X, Y, Z, etc. to participate in the same endeavor without stepping on one another.



## prompt

this is my design workflow, you propose but the next stage is to vastly subtract the features you made unless they are absolute essential, we know we have reached the optimal design until there is nothing more to subtract, got it? now I want this constraint, you are right in saying the Conductor may face the biggest architectural danger and I do not want that. I want Conductor to remain as is. The changes should be 1) the entire project state should be in the Conductor database (e.g. the entire corpus of BM-IST Synthesis, what has been planned, which branch to pursue or not, etc, the other way to think about it is that, when a separate AI/coding agent queries the corpus state say the entire project it will know at a glance the current state of the project and resume where the project left off), and 2) I want to delegate the intelligence to the Agent (it has Plan and Agent work mode) so it is up to the human user how to use and orchestrate those agents (that is the boundary), let us make this question and answer until we agree to a plan


## answer

Got it. I think this is the right constraint, and it actually makes the architecture **simpler**, not more complicated.

For this round, I would lock only these two principles:

### 1. Conductor remains a pure coordination/control-plane system

No planner.
No decomposition engine.
No LLM.
No autonomous scheduling intelligence.
No domain reasoning.
No Solvent semantics.

We do **not** turn Conductor into a smarter system.

Its job remains:

```text
persist project/work state
coordinate work
track dependencies
track ownership
track lifecycle
track activity
expose current state
```

### 2. The Conductor database becomes the project's durable state

The crucial shift is:

```text
Agent context = temporary
Conductor database = durable project memory
```

So for BM-IST Synthesis, another Agent should be able to enter later and ask essentially:

> "What is the state of this project?"

and obtain enough structured state to understand:

```text
what we are doing
what has been established
what has been attempted
what has been rejected
what branches exist
which branch is active
what remains unresolved
what work is pending
what work is blocked
what evidence exists
what should happen next
```

The Agent then supplies the intelligence:

```text
inspect state
↓
reason
↓
plan
↓
decompose
↓
create/update work
↓
execute
↓
replan
```

That means OpenCode/Codex/Antigravity can internally use:

```text
Plan mode
Agent mode
subagents
tools
reasoning
```

without Conductor knowing or caring how that intelligence works.

---

## The crucial subtraction

I would **not** add a new "Project Intelligence" subsystem to Conductor.

Instead:

```text
                ┌─────────────────────┐
                │        Agent        │
                │                     │
                │ Plan / Agent mode   │
                │ subagents           │
                │ reasoning           │
                │ decomposition       │
                └──────────┬──────────┘
                           │
                    creates/updates
                           │
                           ▼
                ┌─────────────────────┐
                │      Conductor      │
                │                     │
                │ durable project     │
                │ state + work graph  │
                │ + dependencies      │
                │ + activity          │
                └─────────────────────┘
```

That is almost the whole system.

Solvent only enters when a particular work item crosses the consequential-action boundary:

```text
ordinary work
    → Agent + Conductor

consequential work
    → Agent
        ↓
      Solvent
        ↓
      Executor
```

---

## One important question before we go further

There is one ambiguity we should settle because it determines whether we are actually subtracting or accidentally making Conductor into a document database.

When you say:

> "the entire project state should be in the Conductor database, e.g. the entire corpus of BM-IST Synthesis"

do you mean:

**A.** Conductor stores the **complete structured project state and knowledge representation**, while large artifacts/source documents can remain external and be referenced;

or

**B.** Conductor literally stores the **entire textual/document corpus itself** in its database?

I strongly suspect **A** is the optimal design, because otherwise we risk turning Conductor into a general-purpose content repository—which violates your "subtract until nothing more can be removed" principle.

Which one do you intend?



## prompt

Option A

## answer

A is the right choice.

That gives us a much tighter design principle:

> **Conductor stores the project's durable structured state, not the entire project's raw corpus.**

So for BM-IST Synthesis, Conductor might know:

```text
Project
├── objectives
├── current state
├── hypotheses / branches
├── decisions
├── tasks
├── dependencies
├── assignments
├── status
├── blockers
├── conclusions
├── evidence references
├── artifact references
└── activity history
```

while the actual large artifacts remain where they naturally belong:

```text
Conductor
   │
   ├── state
   ├── metadata
   ├── relationships
   ├── decisions
   └── references
         │
         ├── source documents
         ├── code repository
         ├── research corpus
         ├── datasets
         └── generated artifacts
```

The Agent then reconstructs the project's current state from Conductor and retrieves the referenced artifacts when necessary.

That gives us an important separation:

```text
RAW KNOWLEDGE / ARTIFACTS
        ≠
PROJECT STATE
```

And it keeps Conductor from becoming a document/blob/content-management system.

### I would now formulate the boundary this way

**Conductor is the durable memory of the project state.**

**Agent is the intelligence operating on that state.**

So a new Agent should be able to arrive later and get something like:

```text
PROJECT: BM-IST Synthesis

OBJECTIVE:
...

CURRENT STATE:
...

ACTIVE BRANCH:
...

ABANDONED BRANCHES:
...

ESTABLISHED:
...

UNRESOLVED:
...

BLOCKED:
...

ACTIVE TASKS:
...

COMPLETED TASKS:
...

DEPENDENCIES:
...

DECISIONS:
...

EVIDENCE:
...

ARTIFACT REFERENCES:
...

RECENT ACTIVITY:
...
```

Then the Agent decides what to do next.

That means the human's choice of **Plan mode, Agent mode, subagents, sequencing, parallelism, etc. remains entirely outside Conductor**.

The boundary becomes exceptionally clean:

```text
Human
  ↓
Agent runtime
  ↓
Conductor
  ↓
durable project state
```

and for consequences:

```text
Agent
  ↓
Solvent
  ↓
Executor
```

### One consequence we should settle next

There is now a key semantic question:

**Does Conductor represent only the current project state, or does it also retain the project's historical decision/state transitions so that a new Agent can understand how the current state was reached?**

I think this matters because your requirement is not merely:

> "Tell the next Agent what is true."

It is closer to:

> "Let the next Agent resume the project intelligently without inheriting the previous Agent's hidden context."

That suggests we need **current state + durable history**, but still without turning Conductor into a general knowledge engine.

That is the next question I would resolve before adding anything else.



## prompt

**current state + durable history**


## answer

Agreed. I would lock that in as the next design principle:

> **Conductor maintains both current project state and durable project history.**

That gives us a very clean distinction:

### Current state

What is true **now**.

```text
active branch
active tasks
blocked tasks
dependencies
current decisions
current hypotheses
current project status
current evidence references
```

A new Agent should be able to query this first and immediately understand where the project stands.

### Durable history

How the project **got here**.

```text
task transitions
decisions
branch changes
accepted/rejected approaches
blockers
results
important observations
activity
```

This matters because two projects can have the same current state but require completely different next actions depending on how they arrived there.

For example:

```text
CURRENT STATE

Branch: B
Status: active
```

is insufficient.

History may show:

```text
A → B
A rejected because X
C abandoned because Y
B chosen after evidence Z
```

That gives the next Agent the context necessary to avoid repeating failed reasoning.

### But keep the subtraction discipline

I would **not** introduce a separate "knowledge graph," "memory engine," "lesson engine," or "project intelligence layer."

The minimal architecture remains:

```text
                 AGENT
                   │
             reasoning / plan
                   │
                   ▼
              CONDUCTOR
          ┌────────┴────────┐
          │                 │
     current state      durable history
          │                 │
          └────────┬────────┘
                   │
              project memory
```

The Agent interprets that state/history.

Conductor merely makes it **durable, queryable, and coordinateable**.

That preserves the fundamental boundary:

```text
CONDUCTOR = STATE
AGENT     = INTELLIGENCE
```

And consequently:

```text
Conductor does not decide:
    what the project means
    what branch is best
    how to decompose work
    what to plan next

Agent decides those things.

Conductor records them once they become
project state/history.
```

### This also clarifies the DAG idea

I would now stop calling it a DAG internally.

The underlying Conductor model should simply be:

```text
Project
  ├── state
  ├── work
  ├── dependencies
  ├── decisions
  ├── history
  └── references
```

Some of those relationships may form a DAG, some may be cyclic, and some are temporal. The **DAG is a useful projection for an Agent**, not the fundamental database model.

That is another valuable subtraction.

### And the multi-Agent vision follows naturally

Agent X can leave:

```text
current state
+
durable history
```

Agent Y can later enter and resume.

Then Agent Z can continue.

They don't need to inherit each other's context windows.

They inherit the **project's durable state**, while retaining independent intelligence.

That is, I think, the central idea we should now build the protocol around:

> **Conductor is the shared durable state machine of the project; Agents are interchangeable intelligence operating on that state.**

The next question I'd resolve is narrower: **what is the minimum set of things that must become durable Conductor state/history for a new Agent to genuinely resume a project without inheriting the previous Agent's hidden context?**



## prompt

read attached beads.md, there is an existing GitHub repo called Beads, I want to adopt its style but not its entire monstrous workflow, only cherrypick what are the best for our workflow, does it have anything net new that would vastly improve our workflow?


Yes. After reading the Beads material, I think there are **a few genuinely excellent ideas worth cherry-picking**, but most of Beads should be rejected for our architecture.

The important distinction is:

> **Adopt Beads' coordination semantics, not Beads' workflow machinery.**

Beads' core insight is very close to what we have converged on: persistent structured work state survives the Agent, and the Agent asks for the current claimable frontier rather than reconstructing project state from memory. 

## 1. The one Beads idea I think is genuinely net-new and highly valuable

### `ready` = the claimable frontier

Beads computes:

> open work + no open blockers + not already claimed/deferred/gated → **ready**

Agents don't scan the entire graph; they ask:

```text
What can I safely work on right now?
```

and claim atomically. 

This fits our Conductor vision **extremely well**.

I would add exactly this concept to Conductor:

```text
READY = derived state
```

Not a new stored state.

Conceptually:

```text
                  Project state
                       │
                dependencies
                       │
                       ▼
                ┌─────────────┐
                │   READY     │
                │  frontier   │
                └─────────────┘
                   /   |   \
                  B    C    D
```

Then Agent X/Y/Z asks:

```text
GET_READY_WORK
```

and receives only tasks it can legitimately take.

This is probably the **single most valuable Beads contribution to our design**.

It also reinforces the idea that Conductor is a control plane, not an intelligent planner.

---

# 2. Atomic claim is worth keeping

Beads':

```text
bd ready --claim
```

means discovering and claiming work can be one atomic operation. 

Our Conductor already has atomic claim semantics, so this is not conceptually new, but Beads reinforces that the Agent protocol should be:

```text
discover ready work
        ↓
claim atomically
        ↓
work
```

rather than:

```text
discover
↓
think about it for a while
↓
try to claim
↓
oops somebody else got it
```

So I would make **atomic frontier acquisition** an explicit Agent↔Conductor protocol capability.

---

# 3. `discovered-from` is worth cherry-picking

This is subtle and very useful.

Beads distinguishes:

```text
blocks
parent-child
discovered-from
related
```

and specifically says `discovered-from` is provenance and does **not** affect readiness. 

I think that is valuable for us.

Example:

```text
Task A: investigate BM-IST synthesis

    discovers →

Task X: possible alternative theoretical branch
```

X is not necessarily a dependency of A.

It is:

```text
X discovered-from A
```

That gives the next Agent historical context without contaminating scheduling semantics.

This fits our:

> current state + durable history

principle extremely well.

---

# 4. Parent-child is useful, but don't build an "epic system"

Beads has parent-child / epic-subtask structure. 

We should probably support the **relationship**, because the Agent may create:

```text
A
├── B
├── C
└── D
```

But we don't need Beads' richer epic workflow.

The meaning should remain simple:

```text
parent
child
```

not:

```text
Epic engine
Milestone engine
Workflow engine
Template hierarchy
```

That would violate your subtraction principle.

---

# 5. Hash IDs are interesting, but I would NOT adopt them

Beads uses content-derived hash IDs so independent agents/branches can create work without renumbering collisions. 

That's clever.

But I don't think it gives us enough benefit to justify changing Conductor's identity model.

We already have stable IDs and explicit correlation identifiers.

For our system:

```text
task_id
run_id
operation_id
intent_id
```

have semantic meaning.

I would keep the existing identifier approach.

**Do not import Beads' hash-ID machinery merely because it is elegant.**

---

# 6. Formulas → molecules → wisps: reject

This is where Beads becomes much heavier.

Beads has:

```text
formula
    ↓
proto
    ↓
molecule
    ↓
wisp
    ↓
gate
```

for repeatable multi-step workflows. 

This is precisely the kind of machinery we have been trying to avoid.

Why?

Because it starts turning Conductor into:

```text
workflow engine
+
template engine
+
instantiation engine
+
runtime orchestration engine
```

Our Agent already has planning and orchestration intelligence.

Therefore:

> **Do not import Beads' molecule/formula/wisp architecture.**

An Agent can create ordinary Conductor tasks when it needs them.

---

# 7. Gates: keep the concept, not the subsystem

Beads' gate concept is interesting:

> a step can wait for an external event such as human sign-off, timer, or GitHub result. 

We probably need the **semantic capability**, because our reference workflow eventually needs things like:

```text
waiting for external result
waiting for human decision
waiting for authorization
```

But we should not build a "Gate Engine."

Instead, represent the minimum state:

```text
blocked / waiting
```

plus a reference to what it's waiting for.

That's enough.

Otherwise we are recreating Beads' workflow machinery inside Conductor.

---

# 8. Dolt synchronization: reject

Beads puts the whole work database into version-controlled Dolt and syncs through Git remotes. 

This is interesting, but it conflicts with the direction we've already chosen:

```text
Conductor = authoritative shared project state
```

We don't need:

```text
Agent clone
   ↕
Git
   ↕
Dolt
   ↕
another clone
```

We already have a central coordination service/API.

Multi-Agent participation should happen through Conductor:

```text
Agent X ─┐
Agent Y ─┼──→ Conductor
Agent Z ─┘
```

That is simpler.

---

# 9. Beads reinforces our "project memory" concept

This is where Beads is particularly relevant.

Its central motivation is:

> coding agents lose memory when a session ends. Persistent structured work survives the agent. 

That maps almost perfectly onto what we have now decided:

```text
Agent memory
    = ephemeral

Conductor
    = durable project state + durable history
```

But we should go one level beyond Beads.

Beads primarily persists work/issues.

Our Conductor should persist:

```text
work
+ dependencies
+ decisions
+ project state
+ activity/history
+ artifact references
+ provenance
```

The raw corpus remains external.

So the Agent can arrive later and reconstruct the project without inheriting the previous Agent's context.

---

# 10. The minimal Beads-inspired Conductor

If I aggressively subtract, I think we should steal only **four ideas**:

### Keep

```text
1. READY frontier
2. Atomic claim
3. Parent/child relationship
4. Discovered-from provenance
```

And maybe a fifth very small semantic:

```text
5. Waiting/blocked-for-external-event
```

But even #5 should probably be represented using existing lifecycle semantics rather than a new subsystem.

Everything else is optional or actively undesirable at this stage.

---

# 11. The resulting workflow becomes very clean

I think this is now the simplest model:

```text
Human
  ↓
AI Agent
  │
  │ plans / reasons / spawns subagents
  │
  │ creates work discovered during execution
  ▼
Conductor
  │
  ├── project state
  ├── durable history
  ├── dependencies
  ├── provenance
  ├── readiness
  └── atomic claims
       │
       ├──── Agent X
       ├──── Agent Y
       └──── Agent Z
```

For ordinary work:

```text
Agent
  ↓
claim READY task
  ↓
do work
  ↓
submit result
  ↓
Conductor updates state
  ↓
new READY frontier
```

For consequential work:

```text
Agent
  ↓
claim READY task
  ↓
construct exact operation
  ↓
Solvent
  ↓
authorization
  ↓
Executor
  ↓
effect/result
  ↓
Conductor records coordination state
```

And crucially:

> **Conductor never decides what the Agent should think about next.**

It only determines which work is **eligible and unclaimed** based on the durable state.

---

# 12. What Beads teaches us about our design

The biggest lesson isn't a Beads feature.

It is this:

> **A control plane becomes powerful when it can derive the next legal work frontier from durable state, rather than when it becomes more intelligent.**

That is exactly the direction you want.

The Agent remains intelligent.

Conductor remains simple.

Solvent remains authoritative.

Executor remains consequential.

And the project itself becomes persistent enough that:

```text
Agent A leaves
    ↓
Agent B arrives
    ↓
queries Conductor
    ↓
sees current state + history + READY frontier
    ↓
continues
```

That is a substantial improvement without turning Conductor into another harness.

### My recommendation

I would **not modify Conductor yet**.

First lock the minimal Beads-inspired semantic additions:

```text
READY frontier
atomic claim
parent-child
discovered-from
```

Then ask the harder question:

> **What is the smallest representation of a Conductor task that lets an arbitrary Agent understand what it can do, what it must produce, and whether the work is ordinary or consequential?**

That is the next design question I'd settle before writing code.




