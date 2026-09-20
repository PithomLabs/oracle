## prompt

for now, let us focus on one AI agent only and focus on how to apply the conductor + solvent workflow against the BM-IST synthesis, how do you plan to do that, reply in same itemized description around categories format, short but concise


## BM-IST Synthesis Pilot — One Agent

### 1. Objective

Use **one AI/coding/research Agent** to prove that Conductor + Solvent can manage a real synthesis project without turning Conductor into an intelligent workflow engine.

```text
Human
  ↓
Agent PLAN
  ↓
Human approval
  ↓
Agent WORK
  ↕
Conductor
  ↓
Solvent only when external consequence occurs
```

### 2. What Conductor Stores

Conductor becomes the durable project state:

```text
BM-IST project
├── current objective
├── active synthesis branch
├── planned tasks
├── dependencies
├── completed/rejected directions
├── decisions
├── conclusions
├── blockers
├── artifact references
└── activity/history
```

The raw BM-IST corpus remains in its owning files/repositories; Conductor stores **structured state and references**, not the entire corpus.

### 3. PLAN Mode

The single Agent examines the current Conductor state and corpus, then produces a plan such as:

```text
A. establish synthesis scope
B. extract claims
C. compare competing interpretations
D. identify unresolved contradictions
E. construct synthesis
F. test/criticize synthesis
G. produce final artifact
```

The Agent can internally use its own Plan mode, subagents, tools, etc.

**Conductor does not participate in the reasoning or plan iteration.**

Human:

```text
APPROVE
or
REJECT / request revision
```

### 4. WORK Mode

After approval, the Agent executes the approved work through Conductor:

```text
DISCOVER READY
→ CLAIM
→ WORK
→ SUBMIT
→ ACCEPT / REJECT
→ next READY work
```

For example:

```text
Extract claims
      ↓
Compare interpretations
      ↓
Resolve contradiction
      ↓
Draft synthesis
      ↓
Critique synthesis
      ↓
Revise
```

Conductor tracks where the project is; the Agent decides what the work means.

### 5. Ordinary BM-IST Work

Most synthesis work should be ordinary:

```text
read corpus
analyze
compare
write notes
generate hypotheses
draft
revise
run local tests
produce artifacts
```

No Solvent involvement.

This is important: **we should not burden intellectual work with authority machinery.**

### 6. Consequential BM-IST Work

Only when the Agent crosses the external-effect boundary:

```text
git push
publish a release
modify an external database
trigger an external workflow
send an external communication
change some authoritative external state
```

then:

```text
Agent
  ↓
exact operation
  ↓
Solvent
  ↓
Executor
  ↓
external SOR
```

The synthesis task itself is not consequential merely because it is important.

### 7. Human Boundary

Keep human intervention extremely simple:

```text
PLAN
  → human approves

WORK
  → Agent proceeds

structural change
  → back to PLAN
  → human approves revised plan
```

No additional approval framework.

### 8. What We Are Actually Testing

The BM-IST pilot should answer four questions:

```text
1. Can a new Agent reconstruct the project from Conductor alone?
2. Can the Agent continuously advance work through READY tasks?
3. Can ordinary synthesis remain completely outside Solvent?
4. Can an external consequential action cleanly cross the Solvent boundary?
```

### 9. First Concrete Pilot

Start very small:

```text
Project: BM-IST Synthesis

PLAN:
    establish synthesis objective
    ↓
    map existing corpus
    ↓
    create first synthesis branch
    ↓
    identify unresolved questions
    ↓
    produce synthesis draft

Human approves.

WORK:
    Agent performs the tasks through Conductor.
```

Then introduce **one deliberate consequential action**, such as committing/pushing the approved synthesis artifact, to prove the same project can transition from:

```text
intellectual work
→ consequential external action
```

without changing the architecture.

### 10. Success Criterion

The strongest test is:

> **End the Agent session, start a fresh Agent session, give it the Conductor project state, and see whether it can resume BM-IST synthesis without inheriting the previous context window.**

That is where the Conductor concept earns its value.



