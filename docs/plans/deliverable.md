The ultimate deliverable should be a **combination**, but not a giant framework. The durable thing is the **Loop Engineering contract**, and the code artifacts exist to prove and operationalize it.

I would structure the deliverable as four layers:

### 1. Loop Engineering Specification — the primary deliverable

This is the intellectual artifact.

```text
Loop Engineering Specification
├── roles
├── boundaries
├── workflow vocabulary
├── handoff semantics
├── consequential-action branch
├── human intervention
├── failure / ambiguity / replay semantics
├── correlation / evidence requirements
├── conformance requirements
└── tool/domain agnosticism
```

This is what makes the model portable across GPT, Claude, Temporal, plain Go, human operators, etc.

### 2. Conformance Test Suite — the proof

This is the most important code deliverable.

It should exercise the **existing**:

```text
Agent/client
    ↕
Conductor
    ↕
Solvent
    ↕
Executor
```

and verify the Loop Engineering contract.

The tests should be capable of being driven by different clients:

```text
GPT client
Claude client
scripted client
human/curl client
Temporal-driven client
```

without changing the contract.

This is what turns "tool-agnostic" from a claim into evidence.

### 3. Reference Agent Skill / Instruction Package — the usability layer

I **would** create an Agent Skill, but as a thin implementation of the specification rather than another framework.

Something like:

```text
loop-engineering/
├── SKILL.md
├── workflow rules
├── decision guidance
├── handoff conventions
├── proposal format
├── evidence expectations
└── examples
```

Its purpose is to teach an agent:

> "Here is how you participate in the Loop Engineering contract."

It should know when to:

```text
discover
formulate
claim
work
propose
request authority
execute through the appropriate boundary
observe
interpret
review
continue
```

But the skill must **not** become the authority mechanism. The agent skill tells the agent how to behave; Solvent remains authoritative.

### 4. Thin integration surfaces — only where needed

This is where MCP becomes useful.

I would **not build a new Loop Engineering MCP server**.

Instead, use the existing interfaces:

```text
Conductor MCP
Solvent MCP/API
Executor interface
```

and, where useful, provide a **thin adapter/conformance harness** that demonstrates how they compose.

So:

```text
                 Loop Engineering Spec
                           │
             ┌─────────────┴─────────────┐
             ↓                           ↓
      Agent Skill                 Conformance Suite
             │                           │
             ↓                           ↓
        Agent/client ─────────────── Conductor MCP/API
                                        │
                                        ↓
                                   Solvent MCP/API
                                        │
                                        ↓
                                     Executor
```

## What I would NOT build

I would explicitly avoid:

```text
❌ Loop Engineering MCP server
❌ Loop Engineering database
❌ workflow runtime
❌ workflow scheduler
❌ workflow UI
❌ new task system
❌ new authority engine
❌ new SDK that abstracts everything
```

Those would turn the specification into another infrastructure product.

## So what is the "product" of the BM-IST demo?

I would describe the ultimate deliverable as:

> **A portable Loop Engineering protocol, an agent skill for participating in it, and a conformance suite proving that the same Agent–Conductor–Solvent–Executor workflow works across domains and tooling.**

BM-IST is the **demonstration workload**, not the framework.

And the minimal Conductor UI remains the human observation surface; we don't need another UI.

### In practical repo terms

I would expect a final BM-IST deliverable roughly like:

```text
bm-ist/
├── specification/
│   ├── loop-engineering-requirements.md
│   ├── workflow-specification.md
│   ├── role-boundary-matrix.md
│   └── conformance-matrix.md
│
├── skill/
│   └── SKILL.md
│
├── conformance/
│   ├── scenarios/
│   ├── scripted-client/
│   ├── agent-client/
│   └── tests/
│
├── bmist/
│   ├── scenario-a/
│   ├── scenario-b/
│   ├── scenario-c/
│   ├── scenario-d/
│   └── scenario-e/
│
└── runbook/
    └── poc-runbook.md
```

The important hierarchy is:

```text
Specification
    ↓
Agent Skill        Conformance Suite
    ↓                    ↓
        BM-IST Demonstration
```

So **yes: Agent Skill + conformance harness + existing MCP/API surfaces** is the right combination. The MCP server itself should remain Conductor/Solvent's existing infrastructure, not become another Loop Engineering component.

The most valuable artifact is actually the **conformance suite**, because that is what demonstrates that the specification survives substitution of the agent, orchestration framework, and domain.
