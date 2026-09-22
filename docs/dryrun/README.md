## prompt

if I copy to a separate folder the argus binary, docs/corpus, AGENTS.md and README.md, would those files suffice OpenCode, give it prompt and start as work agent, or are there other files needed to copy? the goal is OpenCode and other AI agents do not need to see the codebase in order to perform its work


## answer

Almost. **The codebase itself can be completely absent**, which is exactly the direction I would recommend.

For the agent-only workspace, I would use this:

```text
argus-agent/
├── AGENTS.md
├── README.md
├── TASK.md                 # strongly recommended
└── docs/
    └── corpus/
        ├── v6_3.1.md
        ├── v6_2_adv.md
        ├── v6_1.md
        ├── adv_review.md
        ├── adv_review2.md
        ├── adv_review3.md
        └── writeup_v1.1.md
```

### Why `TASK.md` is the one missing piece

Your current `argus.get_context` requires:

```json
{
  "task_id": "<UUID>"
}
```

So a fresh agent needs to know which task it is supposed to work on. 

I would make `TASK.md` extremely small:

```markdown
# ARGUS Work Task

Role: work

Task ID:
00000000-0000-0000-0000-000000000101

Research question:
Can one explicit microscopic dynamical system generate the structures
of quantum theory without assuming them?

Instructions:
1. Call `argus.get_context` first using the task ID above.
2. Read all files under `docs/corpus/` as background research.
3. Treat ARGUS state as the source of current research state.
4. Do not inspect or require the ARGUS source code.
5. Perform bounded research.
6. Submit results through `argus.submit_packet`.
7. Do not attempt promotion, discharge, retraction, or other authority actions.
```

You could also simply put the task ID in your startup prompt instead of creating `TASK.md`. I prefer the file because it makes the agent workspace self-describing.

### The binary is not really an agent-context file

You can put:

```text
argus-agent/
├── argus
├── AGENTS.md
├── README.md
├── TASK.md
└── docs/corpus/
```

and it will work conceptually, but **OpenCode does not need to inspect the binary**.

The cleanest separation is actually:

```text
ARGUS RUNTIME
    /path/to/argus
    + CockroachDB
    + database

AGENT WORKSPACE
    AGENTS.md
    README.md
    TASK.md
    docs/corpus/*
```

OpenCode sees only the second.

Its MCP configuration launches the first:

```text
/path/to/argus mcp
```

That is the architectural separation I would aim for.

### There is one runtime dependency you must not overlook

Your `argus serve` currently starts CockroachDB externally. Your actual runtime log showed:

```text
/home/chaschel/.local/bin/cockroach
```

So the **runtime environment**, not the agent workspace, needs the `cockroach` executable if `argus serve` is responsible for starting the database.

Therefore:

```text
Agent workspace
    AGENTS.md
    README.md
    TASK.md
    docs/corpus/*

Runtime
    argus
    cockroach
    database state
```

The agent does not need Go, the repository, `go.mod`, source files, migrations, `internal/`, `cmd/`, etc.

### The startup flow becomes very clean

You start ARGUS separately:

```bash
/path/to/argus serve
```

Then start OpenCode in the tiny agent workspace and give it the Work Agent prompt.

OpenCode has:

```text
AGENTS.md
TASK.md
docs/corpus/*
```

and through MCP:

```text
argus.get_context
argus.submit_packet
```

So its world is effectively:

```text
             OpenCode
                 |
       ┌─────────┴─────────┐
       |                   |
 filesystem             MCP
       |                   |
       v                   v
 docs/corpus/*       argus.get_context
 AGENTS.md            argus.submit_packet
 TASK.md                   |
       |                   v
       |                 ARGUS
       |                   |
       └────── research ───┘
```

This is actually **better than giving the agent the codebase**.

The agent has no reason to understand:

```text
internal/
cmd/
migrations/
coordinator/
epistemic/
work/
```

It only needs the **capability boundary** and the **research context**.

Your current AGENTS document explicitly establishes this separation: live state comes from `argus.get_context`, while the seven research documents are read directly from the repository filesystem. 

And the intended MCP surface is exactly two tools. 

## My recommended final agent package

I would make the distributable agent workspace:

```text
argus-workspace/
├── AGENTS.md
├── README.md
├── TASK.md
└── docs/
    └── corpus/
        ├── v6_3.md
        ├── v6_2_adv.md
        ├── v6_1.md
        ├── adv_review.md
        ├── adv_review2.md
        ├── adv_review3.md
        └── writeup_v1.md
```

And keep the executable outside it:

```text
argus-runtime/
└── argus
```

with CockroachDB available to the runtime.

That gives you a genuinely interesting property:

> **An AI research agent can work against ARGUS without having access to ARGUS's implementation.**

That is a much stronger demonstration of the capability/authority boundary than simply hiding the code informally.

### Startup prompt for OpenCode

Use this:

```text
You are the ARGUS Work Agent.

You are working as a fresh research agent with no prior conversational
memory.

Your task is defined in TASK.md.

You do NOT have access to the ARGUS source code, and you do not need it.
Do not attempt to obtain or reconstruct the implementation.

Your only ARGUS system interface is MCP:

- argus.get_context
- argus.submit_packet

Your workflow:

1. Read TASK.md.
2. Call argus.get_context using the task_id in TASK.md.
3. Read all seven documents under docs/corpus/ as background research.
4. Treat the ARGUS context returned by get_context as the authoritative
   current research/work state.
5. Identify the current research frontier and unresolved obligations.
6. Perform bounded research.
7. Produce candidate claims, evidence, reasoning, and relationships as
   appropriate.
8. Submit your work through argus.submit_packet.
9. Do not attempt to promote, discharge, retract, or otherwise exercise
   authority.
10. Do not inspect the ARGUS implementation or database.

Remember:

CAPABILITY != WORK != AUTHORITY != EXECUTION

You produce research work.
ARGUS records and validates it.
Human beings adjudicate epistemic debt.
Solvent controls authority.

Begin by reading TASK.md and calling argus.get_context.
```

That is the setup I would use for the actual experiment. It tests the interesting proposition directly: **can an agent operate as a useful research worker against ARGUS as a black-box authority/work substrate?**
