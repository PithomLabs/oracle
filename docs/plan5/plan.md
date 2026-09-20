## prompt

read attached files, there is a new open source project from AWS https://github.com/pizza-bot-app/pizza-bot that I think we can use in conjunction with conductor and solvent, it is MCP-ready, tell me how we can make use of it in this BM-IST pilot, also in order that we have a common grammar and thus avoid vocabulary drift we must speak the Elephant Bridge Protocol v2.1, tell me your implementation plan   



Yes. I think **Pizza Bot is a very good Agent-side implementation for the BM-IST pilot**, precisely because we do not need to turn it into Conductor.

It already provides the pieces we need on the **Agent side**: a stateful DeepAgents/LangGraph runtime, skills, subagents, persistent runs/checkpoints, human-in-the-loop capability, and MCP integration.  

The key is to use it as **the Agent**, while Conductor remains the durable project control plane.

# BM-IST + Pizza Bot + Conductor + Solvent

## 1. Architecture

```text
Human
  ↓
Pizza Bot Agent
  │
  ├── PLAN / WORK
  ├── EBP v2.1 reasoning
  ├── subagents
  ├── research/tools
  │
  └── MCP
       ↓
    Conductor
       │
       └── durable BM-IST state/history

Consequential operation
       ↓
    Solvent
       ↓
    Executor
       ↓
 External SOR
```

Pizza Bot remains the **Agent runtime**. Its internal LangGraph/DeepAgents machinery stays inside the Agent boundary. Its repository explicitly isolates the production graph runtime and exposes extension through skills and MCP.  

## 2. Common Grammar: EBP v2.1

For the BM-IST pilot, use **EBP vocabulary at the intellectual-work layer**:

```text
Idea
Claim
Debt
Retire Debt
Obstruction
Null Model
Toy Check
Faithfulness Review
Promote
Dormant
```

That matches EBP's minimal pipeline:

```text
Capture → Clarify → Debt → Test → Promote
```

and its central rule:

> Ideas enter free. Promotion costs debt. 

This should become the **BM-IST research grammar**, not another technical protocol.

## 3. Do Not Mix Grammars

Keep three vocabularies distinct:

| Layer                | Grammar                        |
| -------------------- | ------------------------------ |
| Agent/research       | **EBP v2.1**                   |
| project coordination | **Conductor**                  |
| consequential action | **Solvent / Loop Engineering** |

So:

```text
EBP:
    Idea has Debt

Conductor:
    Work item is READY / CLAIMED / DONE

Solvent:
    Operation is AUTHORIZED / DENIED / UNKNOWN
```

Do **not** rename Conductor concepts into EBP terminology merely for consistency.

EBP is the common research language, not a replacement for Conductor's coordination semantics.

## 4. BM-IST Project State

Conductor becomes the durable memory of the synthesis:

```text
BM-IST
├── Ideas
├── Claims
├── Debt
├── Obstructions
├── Null models
├── Toy checks
├── Faithfulness reviews
├── Promoted artifacts
├── Dormant ideas
├── current branch/direction
├── active work
├── dependencies
└── history
```

The raw corpus remains external.

This fits EBP particularly well because EBP itself says accounting should preserve the work rather than become the work. 

## 5. Pizza Bot Skill

Create one project skill:

```text
skills/ebp-bmist/SKILL.md
```

Its job is to teach the Pizza Bot Agent:

```text
PLAN:
    understand BM-IST state
    formulate Ideas/Claims
    identify Debt
    propose work

WORK:
    retire Debt
    create evidence
    file Obstructions
    test Null Models
    perform Toy Checks
    conduct Faithfulness Review
    Promote when eligible
```

The Pizza Bot repo already supports skills and project/user extension through MCP/skills. 

Do **not** build a new EBP runtime.

## 6. Conductor MCP

Expose Conductor as the Agent's project-state interface.

The Agent should be able to ask:

```text
project state
current Ideas
active Debt
READY work
claim work
record result
record decision
```

The critical operation is:

```text
GET_PROJECT_STATE
```

A fresh Pizza Bot session should be able to ask:

> What is the current BM-IST synthesis state?

and immediately recover the project without inheriting another session's context.

## 7. How EBP Maps to Conductor

Keep the mapping tiny:

```text
EBP Idea
    ↔ Conductor project artifact/reference

EBP Debt item
    ↔ Conductor work item

EBP Debt retirement
    ↔ completed work + evidence

EBP Obstruction
    ↔ discovered/recorded issue

EBP Promotion
    ↔ accepted/promoted research artifact
```

But **Conductor does not understand whether the debt has actually been intellectually retired**.

It records the evidence and the decision.

That preserves Agent intelligence.

## 8. Plan Mode

Pizza Bot's Agent uses EBP as its reasoning grammar:

```text
Current state
    ↓
Ideas / Claims
    ↓
Debt inventory
    ↓
next smallest useful moves
    ↓
Plan
```

The Human sees the plan and approves/rejects it.

EBP explicitly emphasizes that entry is cheap and promotion is expensive, so the Agent should not require full formalization before an idea can enter the workbench. 

## 9. Work Mode

After approval:

```text
Pizza Bot
    ↓
Conductor DISCOVER
    ↓
READY work
    ↓
CLAIM
    ↓
perform EBP debt-paying work
    ↓
SUBMIT
```

Examples:

```text
"Define the invariant."
"Construct the candidate map."
"Find one null model."
"Run toy check."
"File obstruction."
"Perform faithfulness review."
```

These correspond directly to EBP's "next smallest useful move" principle. 

## 10. Where Solvent Enters

Most BM-IST work never touches Solvent.

```text
read
analyze
formalize
write
compare
test
simulate
```

is ordinary work.

Solvent only appears when the Agent performs an **external consequential operation**, for example:

```text
push promoted synthesis to GitHub
publish artifact
trigger external workflow
modify an external database
```

Then:

```text
Pizza Bot Agent
    ↓
exact operation
    ↓
Solvent
    ↓
Executor
    ↓
external SOR
```

That preserves the Loop Engineering boundary.

## 11. Pizza Bot's Own Memory

Pizza Bot has its own thread/checkpoint/memory mechanisms. 

Use them for:

```text
Agent continuity
conversation context
runtime checkpoints
```

but **not as the authoritative BM-IST project state**.

The hierarchy should be:

```text
Pizza Bot memory
    = Agent-local working context

Conductor
    = authoritative project coordination state/history

External artifacts
    = authoritative corpus/evidence where applicable
```

This distinction is critical.

## 12. Minimal Implementation

I would build only four things:

```text
1. Pizza Bot project skill:
   ebp-bmist

2. Conductor MCP integration:
   project/work/state interface

3. EBP ↔ Conductor minimal mapping

4. One consequential-action path:
   Pizza Bot → Solvent → existing Executor
```

Do **not** modify Pizza Bot's runtime.

Its existing extension model is precisely what we need: MCP servers and skills are already first-class extension points. 

## 13. First BM-IST Experiment

Start with one tiny synthesis:

```text
Idea:
    Boundary information capacity may unify area and entropy.

EBP Debt:
    needMap
    needInvariant
    needToyCheck
    needNullModel
    needObstruction
    needFaithfulnessReview
```

The Pizza Bot Agent creates the Idea and its debt.

Conductor persists it.

The Agent then works through the debt:

```text
needInvariant
    ↓
needMap
    ↓
needToyCheck
    ↓
needNullModel
    ↓
needObstruction
    ↓
needFaithfulnessReview
    ↓
Promote
```

Those are EBP semantics, not new Conductor workflow states. 

## 14. Why Pizza Bot Is a Good Fit

The pairing is unusually clean:

```text
Pizza Bot
    = intelligent Agent runtime

Conductor
    = durable project memory/control plane

EBP
    = research grammar

Solvent
    = consequential authority

Executor
    = external effect
```

Pizza Bot already supports skills, MCP, delegation, persistence, and HITL; we should **compose those capabilities rather than rebuild them**. 

## 15. Implementation Sequence

```text
1. Freeze BM-IST EBP vocabulary
        ↓
2. Create ebp-bmist Pizza Bot skill
        ↓
3. Connect Conductor MCP
        ↓
4. Define minimal EBP ↔ Conductor mapping
        ↓
5. Seed BM-IST project state
        ↓
6. Run one-agent PLAN
        ↓
7. Human approves
        ↓
8. Run WORK through Conductor
        ↓
9. Retire real EBP debt
        ↓
10. Promote first artifact
        ↓
11. Perform one consequential artifact publication
        ↓
12. Solvent authorization
        ↓
13. External SOR confirmation
```

### The design principle

**Pizza Bot should provide the intelligence. EBP should provide the grammar. Conductor should remember and coordinate. Solvent should guard consequences.**

That gives us a real BM-IST pilot without building another agent runtime, another workflow engine, or another research-management system.
