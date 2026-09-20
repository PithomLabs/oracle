## prompt

let us pivot the poc seed to attached file, this is a more organized and more advanced BM-IST Synthesis research program, I would like phase 8 to be implemented with a thin first dry run


## answer

Yes. This is a better POC seed. **BM–IST Synthesis v5.1 should become the canonical research corpus seed**, replacing the earlier `BMIST_POC_SEED.md`. It is explicitly organized as a standing research reference and executable program, with epistemic status tags, gates, work packages, dependency order, verification ledger, and concrete execution outputs. 

The important change is that **Phase 8 should not try to expose the whole 47-section program**. The first dry run should prove that ARGUS can ingest and adjudicate this richer research program without making the UI or Coordinator domain-specific.

## Phase 8 — Thin Dry Run

### Objective

Exercise the entire trust surface against **one real BM–IST research slice**, end to end:

```text
BM–IST v5.1
   ↓
corpus / packet
   ↓
Oracle
   ↓
Coordinator
   ↓
Solvent
   ↓
Trust UI
   ↓
human decision
```

The dry run should demonstrate:

> **Can ARGUS represent one serious research claim, its provenance, dependencies, unresolved obligations, verification evidence, adversarial challenge, and consequential human decision without confusing research material with authority?**

That is much more valuable than building a broad UI mockup.

## Use one narrow research slice

I recommend making the dry-run subject **Gate G0**.

It is unusually well suited because v5.1 already distinguishes it from later speculative material and gives it a precise failure condition:

> countable substrate + literal continuous unitary evolution → structural incompatibility. 

It also has explicit supporting lemmas L1 and L2, a proof sketch, status classification, and a downstream consequence.

So the first run can contain approximately:

```text
Claim
  G0: discrete substrate is incompatible with literal continuous unitary evolution

Supporting claims
  L1: orbit rigidity
  L2: one-parameter triviality

Evidence
  formal proof sketch / cited mathematical source
  BM–IST v5.1 corpus artifact

Dependencies
  L1 → G0
  L2 → G0

Debt
  formal theorem-ledger entry
  independent verification
  counterexample status
  reviewer confirmation
```

The exact debt should still be expressed through the existing BM-IST pack vocabulary rather than inventing new Solvent semantics.

## Thin UI

Only build enough of the Phase 8 surface to perform the run.

### 1. Briefing

Show:

```text
Program: BM–IST Synthesis v5.1
Research slice: Gate G0
Status: THEOREM-LEVEL / pending verification
Current debt: N
Blocking decisions: N
```

The source document explicitly defines status tags and the governing method `proposal → formalization → adversarial audit → theorem/no-go → survivor specification → experiment`. 

### 2. Oracle

Show one belief graph:

```text
L1 ──┐
     ├──> G0
L2 ──┘
```

along with evidence and provenance.

### 3. Decision Room

Make exactly two consequential paths necessary for the dry run:

```text
RETIRE_DEBT
PROMOTE
```

Do not build the full six-decision UX yet.

The expected behavior:

```text
G0 entered
   ↓
debt exists
   ↓
human reviews evidence
   ↓
RETIRE_DEBT
   ↓
remaining debt = 0
   ↓
PROMOTE
   ↓
Solvent authoritative state changes
```

### 4. Refusal

Deliberately attempt:

```text
PROMOTE while debt remains
```

and show the refusal.

This gives you the most important trust demonstration: **ARGUS refuses to cross the consequential boundary merely because the workflow wants to proceed.**

### 5. Chronicle

Show that the history can be reconstructed entirely from Solvent:

```text
belief entered
→ debt retirement
→ promotion
→ refusal(s)
```

No UI-owned decision history.

---

# Pivot the corpus first

Phase 5 should be revised from:

```text
bmist-poc-seed-v1
  BMIST_POC_SEED.md
  payoff_qwen.md
```

to something conceptually like:

```text
bmist-synthesis-v5.1
  bm_ist_synthesis_v5_1.md
```

Keep the Qwen artifact as an **agent-derived secondary artifact**, not the canonical research source. The new document itself explicitly requires provenance for load-bearing claims, including source, status, dependencies, counterexample status, empirical handle, failure condition, and reviewer. 

Do **not** turn every section into a separate belief during this first run. The corpus remains a source artifact; the packet extracts a small number of actionable claims from it.

## What Phase 8 proves

The dry run should produce this evidence chain:

```text
BM–IST v5.1 source
        ↓
G0 claim extracted
        ↓
L1/L2 dependency edges
        ↓
evidence attached
        ↓
debt visible
        ↓
adversarial refusal demonstrated
        ↓
human retires debt
        ↓
G0 promoted
        ↓
complete audit trail
```

That would validate the **architecture**, rather than merely demonstrate that the HTML works.

## What stays out

Do not implement:

* all 47 sections
* full research-program navigation
* automated theorem proving
* automatic claim extraction
* scientific judgment by LLM
* Gate C
* cosmology
* empirical fitting
* multi-agent research orchestration
* a new BM–IST-specific database schema

The v5.1 program explicitly says the real work is to **build the projection rather than merely describe it**, and that imported mathematics must create falsifiable constraints.  

### Recommended Phase 8 acceptance gate

Phase 8 passes when a human can enter ARGUS, inspect **G0**, follow **L1/L2 → G0**, inspect its evidence and outstanding debt, trigger a refusal while debt exists, perform the required adjudication, promote the claim, and reconstruct the complete history — **with no new persistence layer and no bypass of Coordinator/Solvent authority**.

That gives us a genuinely thin first dry run while using the substantially stronger v5.1 research program as the real substrate.


## prompt

give me a rough overview of how the workflow UI goes in phase 8, why it is necessary, how OpenCode is being used as the agent harness as work agent and as adversarial agent, I want a holistic view


## answer

The holistic picture is this:

**Phase 8 is where ARGUS becomes understandable to a human researcher.** The UI is not the authority layer and not the research engine. It is the **control surface through which a human can see what the agent work produced, what the adversarial agent attacked, what remains unresolved, and where a consequential decision is being requested.**

The BM–IST v5.1 document is a good fit because it already defines a disciplined research lifecycle:

`proposal → formalization → adversarial audit → theorem/no-go → survivor specification → experiment`, with a decisive null result treated as a valid outcome. 

## The whole system

Think of Phase 8 as creating this loop:

```text
                 ┌──────────────────────┐
                 │   BM–IST Research    │
                 │      Program         │
                 └──────────┬───────────┘
                            │
                            ▼
                  ┌─────────────────┐
                  │   Work Agent    │
                  │   OpenCode      │
                  └────────┬────────┘
                           │
                    produces packet
                           │
                           ▼
                  ┌─────────────────┐
                  │   Coordinator   │
                  └────────┬────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │     Solvent     │
                  │ epistemic       │
                  │ authority       │
                  └────────┬────────┘
                           │
                ┌──────────┴──────────┐
                │                     │
                ▼                     ▼
        ┌───────────────┐     ┌───────────────┐
        │   Trust UI    │     │ Adversarial    │
        │ human sees    │     │ Agent/OpenCode │
        │ current state │     │ attacks claim  │
        └───────┬───────┘     └───────┬────────┘
                │                     │
                └──────────┬──────────┘
                           ▼
                    Human adjudication
                           │
                           ▼
                    Solvent transition
                           │
                     ┌─────┴─────┐
                     ▼           ▼
                  PROMOTE      REFUSE /
                  / RETIRE     REOPEN /
                  / etc.       RETRACT
```

The critical point is that **OpenCode agents never become the authority**.

---

# 1. What the Work Agent does

OpenCode is the **agent harness**.

The Work Agent gets a bounded research task extracted from the BM–IST program.

For example:

```text
Task: Formalize Gate G0

Source:
BM–IST v5.1 §3

Objective:
Turn L1, L2 and G0 into a precise research packet.

Expected output:
- claims
- dependencies
- evidence
- proposed falsifiers
- unresolved debt
- reproducible artifacts
```

This matches the research program's emphasis on turning ideas into explicit mathematical objects and executable work packages rather than prose descriptions.  

OpenCode can:

* inspect the corpus
* write proofs or proof sketches
* create notebooks/scripts
* run tests
* calculate numerical results
* produce artifacts
* package its results into an EBP packet

But its result is still:

```text
candidate work
```

not:

```text
truth
```

That distinction is foundational.

---

# 2. Why the adversarial agent is separate

Then the same general harness is used as a **different agent role**:

```text
OpenCode
   ↓
Adversarial Agent
```

Its job is not to improve the work.

Its job is to **try to kill it**.

That maps directly to BM–IST's governing methodology: adversarial audit and theorem/no-go are explicit steps, and the document repeatedly says that null results are legitimate research outcomes. 

For G0, the adversarial agent might ask:

```text
Can the connectedness argument fail
because the substrate topology is not actually
the topology assumed by L1?

Is the "countable / totally disconnected"
premise explicitly established?

Does the operator group actually satisfy
the assumptions used in L2?

Is "continuous unitary evolution" being
interpreted literally on the substrate?
```

It can search for:

* counterexamples
* missing assumptions
* invalid inference
* unsupported status upgrades
* contradictory evidence
* alternative constructions
* failures of reproducibility

The important principle is:

> **The work agent tries to construct. The adversarial agent tries to destroy.**

That symmetry is one of the reasons the architecture is useful.

---

# 3. Where the Trust UI comes in

Without the UI, you have agents producing artifacts and backend state.

That is not enough.

A researcher needs to understand:

```text
What is being claimed?
↓
Why is it being claimed?
↓
What evidence supports it?
↓
What does the adversary say?
↓
What remains unresolved?
↓
What decision am I being asked to make?
↓
What happens if I make it?
```

So the UI is essentially a **human epistemic cockpit**.

Not a dashboard of everything.

A small number of views should expose the research state.

---

# 4. The Phase 8 dry-run UI

For the first run, I would make the workflow visually resemble a **research case file**.

### A. Briefing

The researcher opens ARGUS and sees:

```text
BM–IST Synthesis v5.1

Active Research Case
G0 — Discrete Substrate / Continuous Unitary Compatibility

Status:
THEOREM-LEVEL

Claims:
3

Evidence:
4

Open Debt:
2

Adversarial Findings:
1

Pending Human Decision:
PROMOTE
```

The interface should immediately communicate that this is **research under adjudication**, not a task-management screen.

---

### B. Oracle

Clicking G0 opens the epistemic graph:

```text
L1 ─────────┐
            ├──► G0
L2 ─────────┘
```

Then:

```text
G0
│
├── claim
├── status
├── evidence
├── dependencies
├── debt
├── falsifiers
└── adversarial findings
```

This is where the user asks:

> “What does ARGUS currently believe, and why?”

That is the Oracle role.

The v5.1 research architecture already calls for a provenance/verification ledger containing claim, source, status, dependencies, counterexample status, empirical handle, failure condition, reviewer, and notes. 

---

# 5. The work/adversarial split should be visible

A very useful UI concept is:

```text
             G0
          /      \
       WORK      ADVERSARY
        │           │
   constructed   challenged
        │           │
        └─────┬─────┘
              ▼
        Human Review
```

The researcher can see:

**Work Agent**

> “Here is the argument.”

**Adversarial Agent**

> “Here are the assumptions, attacks, counterexamples, and unresolved concerns.”

Then the human decides what survives.

This is much more meaningful than showing one AI-generated "answer."

---

# 6. Decision Room

This is where the architecture becomes consequential.

Suppose the Work Agent produced a clean argument and the adversarial Agent found no fatal issue.

The UI shows:

```text
PROMOTION REQUEST

Claim: G0

Current status:
entered

Open debt:
2

Adversarial status:
no blocking counterexample found

Evidence:
4 artifacts

Required action:
resolve remaining debt
```

The researcher can inspect the evidence and decide.

The UI then invokes:

```text
Trust UI
    ↓
Coordinator
    ↓
Solvent
```

Never:

```text
Trust UI
    ↓
SQL
```

The Coordinator submits the decision; Solvent determines whether the state transition is actually permitted.

That means the human can **request** a promotion, but cannot simply declare one.

---

# 7. The most important dry-run moment: refusal

I would deliberately design the first demonstration around this.

Start with:

```text
G0
Open debt = 2
```

Human presses:

```text
PROMOTE
```

The system responds:

```text
REFUSED

Reason:
Promotion gate failed.

Outstanding debt:
- needInvariant
- needFaithfulnessReview
```

This is the moment ARGUS proves its reason for existing.

The workflow system might say:

> Task completed.

The agent might say:

> Proof completed.

But Solvent says:

> **You are not yet entitled to promote this claim.**

That is the core architectural distinction between work and authority.

---

# 8. Then the adversarial loop

The Work Agent produces something.

The Adversarial Agent attacks it.

The human sees both.

Then potentially:

```text
REOPEN
```

or:

```text
RETRACT
```

or:

```text
RETIRE DEBT
```

or eventually:

```text
PROMOTE
```

This gives you a genuine iterative scientific loop:

```text
construct
   ↓
attack
   ↓
revise
   ↓
attack again
   ↓
human adjudication
   ↓
authority transition
```

The v5.1 program itself is structured around this kind of progressive evidence ladder:

`definition → lemma → construction → theorem → numerical validation → experimental prediction → independent confirmation`. 

ARGUS is essentially giving that ladder an **operational substrate**.

---

# 9. Why OpenCode matters

OpenCode is not just "the coding AI."

In this architecture it becomes the **portable research-worker runtime**.

You can run the same harness in two roles:

```text
OpenCode
 ├── Work Agent
 │    ├── inspect corpus
 │    ├── formulate
 │    ├── implement
 │    ├── experiment
 │    └── produce evidence
 │
 └── Adversarial Agent
      ├── inspect work packet
      ├── attack assumptions
      ├── search counterexamples
      ├── reproduce results
      └── produce challenge packet
```

Then ARGUS sits above those agents and answers:

```text
What did they produce?
What evidence exists?
What conflicts?
What remains unresolved?
What is authoritative?
What is the human allowed to approve?
```

That is why the **agent harness remains separate from authority**.

---

# 10. The deeper architectural picture

The three systems have very different jobs:

```text
OpenCode
    = intelligence / agency

Conductor
    = operational coordination

Solvent
    = epistemic + authorization authority

Trust UI
    = human control surface
```

And the research program provides the domain semantics:

```text
BM–IST Domain Pack
    = what counts as a claim,
      evidence,
      debt,
      falsifier,
      retirement rule,
      consequential action
```

So the complete architecture becomes:

```text
                BM–IST v5.1
                     │
                     ▼
              Domain semantics
                     │
          ┌──────────┴──────────┐
          ▼                     ▼
     Work Agent           Adversarial Agent
      (OpenCode)               (OpenCode)
          │                     │
          └──────────┬──────────┘
                     ▼
                  Packets
                     │
                     ▼
                Coordinator
                     │
             ┌───────┴───────┐
             ▼               ▼
         Conductor        Solvent
        operational      authority
                     │
                     ▼
                 Trust UI
                     │
                     ▼
                  Human
                     │
                     ▼
            consequential decision
```

## The Phase 8 dry run is therefore not primarily a UI test

It is a **closed-loop proof of architecture**.

You are demonstrating that:

> **An autonomous agent can do substantial research work, another autonomous agent can attack that work, and a human can adjudicate the resulting epistemic state through a thin control surface—without either agent being able to declare its own work authoritative.**

That is the real reason Phase 8 is necessary. BM–IST v5.1 gives you a sufficiently mature research program to make that demonstration much more credible than the original small seed.


## prompt

just to clarify so that we are on the same page, this is how I envision the workflow: 1) first I run OpenCode as work agent, give prompt for initial task 2) OpenCode runs Argus, Argus saves the relevant objects per the EBP packet (remember: everything must speak a common grammar and vocabulary, protocol), 3) I run a separate OpenCode as adversarial agent, give it a prompt to do adversarial review of the work so far, is that workflow right? If so, how does OpenCode as adversarial agent reconstruct what has been done so it has the necessary context and information to do the actual adversarial task, does it have a protocol to do exactly that? The attached DAG.md comes to my mind again because the task should have a minimal description and status that states something like "this task (with metadata) is dead end", the goal is to have interface boundary so that the agent need not reinvent the wheel and not waste time for what could have been done already (minimal interface / implementation metadata protocol)


## answer

Yes — **that is the workflow I would lock in**. The important refinement is that the adversarial OpenCode instance should **not have to reconstruct the work by rereading the repository or guessing from prose**. ARGUS must expose a small, canonical **research-state protocol** that lets any agent reconstruct the current case deterministically.

That is exactly where your DAG idea becomes important. The DAG document argues that agents should receive bounded proof obligations with prerequisites, required evidence, forbidden shortcuts, and an explicit definition of done, rather than a vague task description.  It also argues that nodes need persistent state so a future agent can discover what has already happened instead of repeating it. 

## The workflow I think we should lock

```text
1. HUMAN
   ↓
   creates/assigns initial research task

2. OPENCODE — WORK AGENT
   ↓
   receives bounded task
   ↓
   inspects canonical ARGUS context
   ↓
   does research / coding / experiments
   ↓
   produces artifacts + findings
   ↓
   submits EBP packet

3. ARGUS
   ↓
   validates packet against common grammar
   ↓
   records claims, evidence, dependencies, debt,
   task state, artifacts, etc.
   ↓
   updates authoritative research state

4. OPENCODE — ADVERSARIAL AGENT
   ↓
   asks ARGUS for the current case context
   ↓
   reconstructs what has already been done
   ↓
   identifies unresolved obligations
   ↓
   attacks the existing work
   ↓
   produces adversarial EBP packet

5. ARGUS
   ↓
   records challenge/evidence/findings
   ↓
   updates epistemic state

6. HUMAN
   ↓
   reviews work + adversarial findings
   ↓
   adjudicates consequential transition
```

The key point is that **step 4 starts from ARGUS state, not from the adversarial agent's own memory.**

---

# The missing piece: a common reconstruction protocol

I would make this a first-class ARGUS protocol.

Call it something simple such as:

> **Research Context Protocol (RCP)**

It is not another database or research framework. It is the **read contract** between an agent and ARGUS.

An OpenCode agent should be able to say, conceptually:

```text
GET /v1/research-cases/{case}/context
```

and receive the minimum canonical state necessary to continue the work.

For example:

```json
{
  "protocol": "argus-context/v1",
  "case_id": "...",
  "program_ref": "...",
  "objective": "...",

  "task": {
    "id": "...",
    "claim": "...",
    "status": "in_review",
    "metadata": {
      "kind": "formalization",
      "scope": "G0"
    }
  },

  "dependencies": [...],

  "claims": [...],

  "evidence": [...],

  "artifacts": [...],

  "open_debt": [...],

  "findings": [...],

  "adversarial_history": [...],

  "decisions": [...],

  "next_obligations": [...]
}
```

The **agent should not need to understand Solvent tables, Conductor internals, or UI implementation**.

It speaks one protocol.

That fits your larger invariant:

```text
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
```

And now we add:

```text
AGENT CONTEXT ≠ AGENT MEMORY
```

The authoritative context lives outside the agent.

---

# The task node is the crucial interface boundary

Your DAG idea is exactly right here.

The adversarial agent should not receive:

> “Review the BM–IST work so far.”

That is too vague.

It should receive something like:

```text
TASK
ID: G0-ADVERSARIAL-001

OBJECTIVE
Attempt to falsify Gate G0.

TARGET
belief:G0

PREREQUISITES
L1
L2

CURRENT STATUS
in_review

WORK COMPLETED
W-001 formalization
W-002 proof artifact
W-003 supporting derivation

EVIDENCE
E-001
E-002
E-003

OPEN DEBT
needInvariant
needObstruction

KNOWN FAILURES
none

KNOWN ATTACKS
A-001

DONE WHEN
Each stated premise has been attacked or verified,
and every material counterexample is recorded.

DO NOT REDO
W-001
W-002
W-003
```

Now the adversarial OpenCode agent has a **bounded interface to reality**.

It does not have to rediscover the project.

---

# This is where “dead end” becomes extremely valuable

I agree strongly with your example:

> “this task (with metadata) is dead end”

That should be machine-readable state.

For example:

```text
status:
  proposed
  active
  blocked
  in_review
  accepted
  rejected
  dead_end
  superseded
  cancelled
```

But importantly:

**`dead_end` should not mean “nothing learned.”**

It means:

> this particular route no longer merits further work.

The node should retain:

```text
dead_end:
  reason
  evidence_refs
  discovered_at
  discovered_by
  superseded_by
```

So the next agent gets:

```text
DO NOT PURSUE:
T-017
Reason: violates G0
Evidence: E-044
Replacement direction: T-023
```

That is exactly the kind of institutional memory your DAG paper is arguing for: state must be externalized so future agents don't waste cycles rediscovering prior work. 

---

# The important distinction: EBP packet vs context

I would **not** make the EBP packet itself carry the entire history.

Instead:

```text
EBP Packet
= what this agent is submitting now

ARGUS Case State
= everything established so far
```

So:

```text
Work Agent
   ↓
EBP packet
   ↓
ARGUS
   ↓
canonical case state

Adversarial Agent
   ↑
canonical case-state projection
   ↓
EBP adversarial packet
```

That is cleaner.

The packet is the **message**.

The ARGUS case is the **shared protocol state**.

---

# And this means the common grammar becomes very important

Every participant speaks the same basic vocabulary:

```text
Program
Case
Task
Claim
Evidence
Artifact
Dependency
Debt
Finding
Decision
Status
Provenance
Falsifier
```

The agent can specialize its **behavior**, but not invent its own nouns.

Thus:

```text
Work Agent
  "I produced Claim C, Evidence E,
   Artifact A, and retired Debt D."

Adversarial Agent
  "I challenged Claim C,
   produced Finding F,
   and opened Debt D2."

Human
  "I accepted Finding F,
   retired D2,
   and authorized transition X."
```

Everything is interoperable.

That is much more important than making the agents themselves interoperable at the prompt level.

---

# The DAG then becomes the navigation layer

I would refine the earlier architecture slightly.

The DAG is **not the canonical database model**.

It is the relationship model exposed by ARGUS.

The agent sees:

```text
             G0
            /  \
          L1    L2
           \    /
            \  /
           Evidence
              |
          Verification
              |
        Adversarial Review
              |
          Human Decision
```

But beneath that, ARGUS stores the independently authoritative objects.

So:

```text
DAG
 = navigation / dependency / obligation view

Objects
 = claims, evidence, tasks, artifacts,
   decisions, debt, findings

Solvent
 = epistemic authority

Conductor
 = operational state
```

That prevents us from making the classic mistake of forcing the entire software/research world into one DAG.

Your attached DAG analysis makes exactly this distinction: the implementation itself may be tangled, while the **reasoning-and-evidence graph** can remain comparatively small and structured.  

---

# What the adversarial agent actually does at startup

I would make the adversarial OpenCode startup protocol almost mechanical:

```text
STEP 1
Identify case.

STEP 2
Fetch canonical context.

STEP 3
Read active task.

STEP 4
Read completed tasks and their terminal states.

STEP 5
Read claims and dependency graph.

STEP 6
Read evidence and artifact manifests.

STEP 7
Read open debt.

STEP 8
Read prior adversarial findings.

STEP 9
Construct attack plan only against unresolved obligations.

STEP 10
Produce adversarial EBP packet.
```

That is the **interface boundary** you are describing.

The adversarial agent doesn't ask:

> “What have the other agents been doing?”

It asks ARGUS:

> **“Give me the canonical state of research case X according to protocol version Y.”**

That is a major architectural improvement.

---

# The first Phase 8 dry run should therefore test exactly this

Not a big UI.

Use G0 from BM–IST v5.1 and perform:

```text
HUMAN
  ↓
create G0 task

WORK OPENCODE
  ↓
GET context
  ↓
work
  ↓
submit EBP

ARGUS
  ↓
validate + persist/project

ADVERSARIAL OPENCODE
  ↓
GET context
  ↓
observe:
  G0
  L1
  L2
  evidence
  debt
  previous work
  task status
  ──────────────
  no rediscovery required
  ↓
attack
  ↓
submit adversarial EBP

ARGUS
  ↓
record attack

HUMAN
  ↓
inspect Briefing / Oracle / Challenge
  ↓
decide
```

The **real Phase 8 acceptance criterion** should therefore be:

> A completely fresh OpenCode adversarial agent can start with no conversational history, query ARGUS through the common protocol, reconstruct the current research state, identify what has already been attempted, understand what remains unresolved, and perform a meaningful adversarial review without duplicating completed or terminal work.

That is a much stronger POC than merely proving that the Trust UI renders.

And I think this is the right place to formalize the principle from the DAG work:

> **Agents should not inherit context. They should reconstruct context from a canonical protocol.**

That is how you prevent context-window dependency from becoming the hidden failure mode of the whole ARGUS architecture.


## prompt

Let us vastly simplify the Research Context Protocol RCP with the following requirements and constraints: 1) it must be domain-agnostic 2) it must not reinvent the metadata context schema, it must be based on Conductor database and Solvent database, 3) read attached implementation plans from 6 AIs to see the entire perspective of what this really is and draft RCP based on those insights, the goal is to leverage Conductor and Solvent databases if that is plausible and not add unnecessary complexity


## answer

Yes. After reading the six plans together, I think we should **shrink RCP dramatically**.

The common insight across the plans is not that we need another research metadata system. It is that the research program needs an **external, queryable state surface** so an agent can discover what exists, what was attempted, what passed, what failed, and what remains without relying on conversation memory. One plan explicitly describes the purpose as turning the architecture into a “queryable ledger”; another makes adversarial review a standing pipeline stage; another emphasizes retaining negative results and dead ends rather than rediscovering them.  

That does **not** justify another database.

# RCP should be a read protocol, not a data model

I would define it this way:

> **RCP is the minimal read-only contract by which an agent reconstructs the current research context from Conductor + Solvent.**

Nothing is stored by RCP.

Nothing is duplicated by RCP.

No new RCP tables.

No new RCP graph.

No separate RCP metadata schema.

```text
                 ┌──────────────┐
                 │   OpenCode   │
                 └──────┬───────┘
                        │
                  RCP read request
                        │
                        ▼
                 ┌──────────────┐
                 │ Coordinator  │
                 └──────┬───────┘
                        │
              ┌─────────┴─────────┐
              ▼                   ▼
         Conductor             Solvent
       operational state    epistemic state
              │                   │
              └─────────┬─────────┘
                        ▼
                 RCP response
                        │
                        ▼
                   OpenCode
```

That is the whole idea.

---

# What the six plans collectively tell us

There are really four pieces of context that an agent needs.

### 1. What am I supposed to work on?

That already lives in **Conductor**:

```text
project
task
dependency
task status
activity
governance_ref
```

The plans repeatedly describe work as bounded packages with explicit dependencies, deliverables, gates and kill conditions. 

### 2. What does the system currently claim?

That belongs to **Solvent**:

```text
belief
belief_edge
evidence
debt
action_intent
audit_activity
```

The research plans emphasize stable claim identity, dependencies, provenance, counterexample status, failure conditions and verification history. 

### 3. What has actually happened?

Again, existing state:

```text
Conductor activity
Solvent audit_activity
```

No RCP history table.

### 4. What remains to be done?

This should be **derived**, not stored.

For example:

```text
task active
+ prerequisite task blocked
→ agent should not start downstream work
```

or:

```text
belief has debt
→ promotion is not currently available
```

or:

```text
task cancelled + activity says dead end
→ don't repeat this route
```

That is precisely what the plans are trying to accomplish with explicit gates, stop conditions and retained negative results. 

---

# RCP v0.1

I would make the protocol almost embarrassingly small.

```text
RCP Request

project_ref
task_ref
perspective
```

Where:

```text
perspective = work | adversarial
```

That's basically it.

The agent says:

> “Give me the current research context for task T-017.”

ARGUS resolves the rest.

---

# RCP Response

The response should contain only four sections.

```json
{
  "protocol": "argus-context/v1",

  "task": {},
  "dependencies": [],
  "epistemic": {},
  "activity": []
}
```

### `task`

Direct projection of the Conductor task:

```json
{
  "id": "...",
  "project_id": "...",
  "title": "...",
  "description": "...",
  "status": "...",
  "governance_ref": {}
}
```

No new task metadata model.

### `dependencies`

Direct projection of relevant Conductor dependencies:

```json
[
  {
    "task_id": "...",
    "relation": "...",
    "status": "accepted"
  }
]
```

This is the agent's minimal map of:

> what this task depends on and what has already happened upstream.

### `epistemic`

A projection of existing Solvent objects related to the task:

```json
{
  "beliefs": [],
  "evidence": [],
  "edges": [],
  "debt": [],
  "action_intents": []
}
```

Again, these are **existing Solvent objects**, not RCP objects.

### `activity`

A compact chronological projection:

```json
[
  {
    "source": "conductor",
    "activity": "...",
    "at": "..."
  },
  {
    "source": "solvent",
    "activity": "...",
    "at": "..."
  }
]
```

That's enough for the agent to reconstruct the operational history.

---

# And that's it

I would **not** put these into RCP v0.1:

```text
next_obligations
confidence
agent_summary
research_graph
lessons
recommendations
context_score
priority
semantic_tags
dead_end object
```

They are all tempting, but they start turning RCP into another ontology.

The system should derive what it can from the actual authoritative state.

---

# What about your “dead end” example?

This is the one place where I think your instinct is absolutely right, but we should avoid adding a new RCP status.

For example:

```text
Task T-017

status: cancelled
```

and its existing Conductor activity records:

```text
reason: dead_end
```

Then the RCP response naturally tells the next agent:

```text
T-017
status = cancelled

latest activity:
dead_end

evidence:
E-044

related belief:
B-031
```

The next agent doesn't need a special RCP `dead_end` schema.

It understands:

> **This route was attempted and terminated. Do not repeat it unless new evidence reopens it.**

That fits the plans' explicit requirement that rejected claims and negative results remain in the ledger with their failure mode rather than being deleted. 

And importantly, this preserves our frozen Conductor state machine rather than inventing a new one.

---

# The adversarial agent then becomes extremely simple

The fresh OpenCode instance starts with no prior conversation.

Its bootstrap prompt can essentially say:

```text
You are the adversarial research agent.

First obtain the ARGUS context for the assigned task
using RCP.

Treat the returned Conductor task and Solvent state
as the authoritative description of work completed so far.

Do not redo work whose task state and evidence show
that it has already been completed.

Do not repeat routes recorded as dead ends unless
new evidence materially changes their premises.

Attack the unresolved claims and obligations.

Submit your findings as an EBP packet.
```

That's dramatically cleaner than giving the agent a giant context bundle.

The protocol itself provides the context.

---

# EBP remains the common write grammar

This distinction is important:

```text
RCP
= READ grammar

EBP packet
= WRITE grammar
```

So:

```text
OpenCode
   │
   │ RCP
   ▼
ARGUS context
   │
   │ work
   ▼
OpenCode
   │
   │ EBP packet
   ▼
ARGUS
```

Then the adversarial agent:

```text
OpenCode
   │
   │ RCP
   ▼
ARGUS current state
   │
   │ adversarial work
   ▼
OpenCode
   │
   │ EBP packet
   ▼
ARGUS
```

That gives us the **common protocol grammar** without inventing yet another research schema.

---

# The deeper architecture becomes very clean

```text
CONDUCTOR
"What work exists,
what depends on what,
and what happened operationally?"

              +

SOLVENT
"What claims, evidence,
debt, authority and decisions exist?"

              ↓

             RCP

"Here is the minimum
current context for this task."

              ↓

            AGENT

"What can I contribute?"

              ↓

          EBP PACKET

"Here is what I produced."

              ↓

        COORDINATOR

"Submit / validate / project."

              ↓

           CONDUCTOR
           + SOLVENT
```

That is the architecture I would freeze.

The six plans strongly reinforce this direction: **formalize the work, retain provenance, make claims queryable, keep adversarial review explicit, preserve dead ends, and avoid parallel metadata machinery.** 

## One important correction to our earlier Phase 8 thinking

I would **not implement a general-purpose “Research Context Protocol” as a new subsystem**.

I would implement:

> **one Coordinator read endpoint that assembles the minimum context projection from Conductor + Solvent.**

Something conceptually like:

```text
GET /v1/context/{task_id}
```

The Coordinator already knows how to reach both systems. The endpoint simply joins the existing authoritative state and emits the RCP envelope.

That makes RCP almost free architecturally.

And it gives us the exact property you are after:

> **The agent does not need to remember the work. It only needs to know how to ask ARGUS what the current state of the work is.**

That is the interface boundary. It is also where the DAG insight becomes practical without turning ARGUS into a DAG database.


## prompt

I forgot to add in our workflow, as always we should be compliant with EBP 2.1, that is Ideas Enter Free. Promotion Costs Debt. It is up to the humans to discharge debts so we have to exercise Solvent in our workflow


## answer

Exactly. That changes one important aspect of the workflow: **EBP is not merely metadata attached to research; it is the epistemic operating discipline through which ARGUS exercises Solvent.**

EBP v2.1 makes the rule explicit:

> Ideas enter free. Promotion costs debt. Debt does not kill. Debt is forever payable. New evidence creates new debt. No final-truth claim may be promoted. 

So our workflow should be:

```text
HUMAN
  │
  │ initial task / idea
  ▼
WORK OPENCODE
  │
  │ RCP reads current state
  ▼
ARGUS
  │
  │ EBP packet
  ▼
SOLVENT
  │
  │ idea/claim/evidence/debt
  ▼
ADVERSARIAL OPENCODE
  │
  │ RCP reconstructs state
  │ attacks claims / finds objections
  ▼
ARGUS
  │
  │ EBP adversarial packet
  ▼
SOLVENT
  │
  ├── new evidence
  ├── new obstruction
  └── new debt
  ▼
HUMAN
  │
  │ discharges debt
  ▼
SOLVENT
  │
  │ promotion gate
  ▼
PROMOTED CANDIDATE
```

### The critical correction

**Agents do not retire epistemic debt. Humans do.**

An agent can produce the work needed to pay a debt:

```text
Work Agent:
  produces map
  produces invariant
  runs toy check
  produces null model
  identifies obstruction
  prepares faithfulness material
```

But the transition:

```text
needToyCheck → retired
needInvariant → retired
...
```

is ultimately a **human adjudication decision**, exercised through Solvent.

That is directly consistent with EBP's definition of promotion: promotion means current operational debt has been retired; it does **not** mean the idea has been proven true. 

## Therefore RCP needs even less

RCP should **not** try to tell the agent:

> “Here are the research conclusions.”

It should tell it:

> **“Here is the current work and epistemic state, including what is still unpaid.”**

So the minimal context naturally contains:

```text
CONDUCTOR
  task
  dependency
  status
  activity

SOLVENT
  beliefs
  evidence
  edges
  debt
  audit history
```

The agent then determines what useful work can be done against the current debt.

This fits EBP's minimum pipeline:

```text
Capture → Clarify → Debt → Test → Promote
```

where testing retires debt and promotion happens only afterward. 

## The first dry run should explicitly exercise Solvent

This is important. We should **not** make Phase 8 a workflow where everything is already debt-free and the UI merely displays state.

The dry run should intentionally do this:

### 1. Work agent

OpenCode receives:

> Investigate G0.

It produces an EBP work packet.

ARGUS records the candidate claim and evidence.

Solvent shows:

```text
G0
Alive: yes
Promoted: no

Debt:
☑ needMap
☑ needInvariant
☑ needToyCheck
☑ needNullModel
☑ needObstruction
☑ needFaithfulnessReview
```

EBP explicitly says a fresh idea may enter without a map, invariant, null model, obstruction, proof, or formalization. 

### 2. Adversarial agent

Fresh OpenCode instance queries RCP.

It discovers:

```text
G0
 ├─ work already performed
 ├─ evidence already produced
 ├─ current debt
 └─ previous dead ends / objections
```

It attacks the work.

Suppose it finds a missing premise.

ARGUS records the challenge.

Solvent adds:

```text
needObstruction
```

or preserves it.

**The idea remains alive.**

That is central to EBP: debt blocks promotion; debt does not kill. 

### 3. Human uses the Decision Room

Now the human examines the evidence and adversarial findings.

The human can:

```text
RETIRE_DEBT
REOPEN
RETRACT
REFUSE
```

and eventually:

```text
PROMOTE
```

The important distinction is:

```text
Agent produces evidence
        ≠
Agent discharges debt

Human adjudicates debt
        ↓
Solvent records transition
```

### 4. New evidence can reopen the claim

Suppose G0 is promoted and later the adversarial agent discovers a genuine obstruction.

EBP v2.1 explicitly allows:

```text
promoted
   ↓
new evidence
   ↓
new debt
   ↓
unpromoted
```

without deleting the idea. 

That means Solvent's existing ability to reintroduce debt is not a side feature. **It is the mechanism that makes EBP operational.**

## This gives each component a very clean role

```text
OpenCode Work Agent
    = produce candidate work

OpenCode Adversarial Agent
    = try to break candidate work

Conductor
    = remember operational work state

Solvent
    = remember epistemic state and enforce promotion

RCP
    = read-only bridge exposing existing state

EBP
    = governing research protocol

Human
    = adjudicate and discharge debt
```

And importantly:

> **RCP does not own the research state. EBP does not create another database. Conductor and Solvent remain the stores of record.**

That is the simplification we were looking for.

### The Phase 8 proof becomes much stronger

The first dry run should prove this exact loop:

```text
Idea
 ↓
Work
 ↓
EBP packet
 ↓
Solvent
 ↓
Adversarial reconstruction via RCP
 ↓
Challenge
 ↓
New / remaining debt
 ↓
Human review
 ↓
Debt retirement
 ↓
Solvent promotion gate
 ↓
Promoted candidate
```

That is a real demonstration of **EBP + OpenCode + ARGUS + Conductor + Solvent**, rather than merely a UI demonstration.

And it preserves EBP v2.1's most important anti-bureaucracy principle: accounting exists only to prevent accidental promotion; it must never become the work itself. 


## prompt

just to clarify, 1) how does OpenCode query RCP, is it via MCP, API, etc? again as much as possible we have to maximize what are available with both Conductor and Solvent interfaces via MCP and/or API and extend the latter 2 if absolutely necessary, 2) for our UI I like a sidebar with the following options: a) Insights where it summarizes the insights gained so far with stats showing which tasks in order of priorities from highest to lowest that need attention Kanban board style and b) a UI for humans to discharge debts showing insights why it needs to be discharged (remember that in our Solvent hackathon demo we already did this in the UI) and 3) consolidate attached where net valid before we finalize the UI, workflow and proceed to phase 8 dry run:  >> Compliance check of the RCP/Phase-8 workflow against EBP v2.1 as provided — and the result is mostly good news: **the architecture you built is already a faithful EBP implementation at the structural level, with one genuine mapping conflict and three doctrine tensions worth resolving explicitly rather than by drift.**

## Where the workflow already complies

The mapping is nearly one-to-one, and it's worth seeing how much of the frozen machinery was EBP all along:

| EBP v2.1 doctrine | Workflow mechanism | Status |
|---|---|---|
| Ideas enter free | Packet ingestion creates beliefs at `entered`, agents can propose anything, no entry gatekeeping | ✓ Compliant |
| Promotion costs debt | Kernel promotion gate (empty-debt + no-final-truth, DB-enforced) | ✓ Compliant — literally kernel-enforced |
| Debt does not kill | No debt-induced rejection path; beliefs persist | ✓ |
| Debt is forever payable | Dormant beliefs retain state indefinitely | ✓ |
| New evidence creates new debt | Adversarial packets → contradicts edges → `ebp_add_debt` → promotion suspended | ✓ Compliant |
| No final-truth claim promoted | Kernel-gated, human-gated | ✓ |
| Accounting never becomes the work | Coordinator does accounting deterministically; agents produce; humans adjudicate — the protocol's ideal division | ✓ Exceeds it |

The `ebpInitialDebt` = six items rule already implements §4.4's `fullDebt`; the gate map already reserves faithfulness to humans, which §6.6 demands. So the workflow is EBP-conformant structurally. The findings below are where the letter of v2.1 differs from what's built.

## Critical

**C1 — "Partial retirement" has no representation, and §10/§19.2 assume it.** EBP v2.1's seamless-accounting examples retire debts *partially* ("needMap partially retired" appears twice; retirement conditions are "minimal," implying gradations). The kernel's debt is `TEXT[]` — binary membership, no partial state. Your workflow has the same hole: a retirement decision either removes the item or doesn't. Three honest options: (a) **keep binary, document the fidelity loss** — the first EBP implementation to enforce the *stricter* reading, defensible for a POC; (b) re-submission cycles — a partial retirement becomes keep + re-add + later full retire (awkward, and addDebt's dedup in §4.8 means "partial" never persists in EBP's own Lean kernel either — note that EBP's own formal model has the same limitation); (c) snapshot-JSONB notation. Recommend (a) + record in the decision record: the implementation enforces full-item retirement only, partial credit is prose, not state. This also touches the pending debt-vocabulary decision from Plan 11.1's aftermath — same decision point.

**C2 — "System inference" (§10) is prohibited by your own doctrine — and rightly so.** EBP's §10 envisions the system *inferring* retirement from natural language ("the invariant is boundary distinguishability capacity I(B)" → needInvariant retired). Your architecture's entire epistemics — coordinator-doesn't-repair-meaning, advisory agent claims, human-gated retirement — exists to prevent exactly that inference. This is not a compliance failure; **it's a deliberate, documented strengthening of EBP**: retirement in your workflow always occurs via attributed human decision or declared mechanical criteria, never text-matching. Write that one sentence into the Phase 8 spec as an explicit deviation-with-rationale, so nobody later "fixes" it to match §10 and installs an LLM into the ledger boundary. Related: EBP's "partially retired" prose feeds the same LLM-inference temptation — same sentence covers it.

## High

**H1 — `REOPEN` has no EBP counterpart, and the mapping I'd expect is wrong.** EBP §4.11 has `clearFinalTruthClaim`, and §9's living-ledger loop runs promote → new debt → alive → repaired → re-promoted — *no un-promote-and-reenter operation exists in EBP*. Your workflow's REOPEN (implemented as new `EnterBelief` with same claim) has no EBP name. Not a violation — an *extension* — but it must be recorded as such, with the `derived_from` edge to the retracted origin (the standing M-medium) to preserve lineage, or the living-ledger property ("new evidence creates new debt" on the *same* idea) silently breaks: retraction→reentry in your system creates a fresh free-entry idea rather than a debt-bearing continuation. EBP's model says the *idea* stays alive and accrues; your implementation restarts the ledger. Either accept and record (POC fidelity loss) or implement re-entry-with-debt-inheritance (kernel-adjacent, needs the edge). For G0's Phase 8 loop, the inheritance matters — the whole point of retracting G0 is that its *replacement* inherits the failure context.

**H2 — Attention/lane state (§8, §11) is unowned.** EBP has dormant/active + Inbox/Workbench/Promoted lanes; your system has Conductor lifecycle (OPEN/ACTIVE/BLOCKED/SUBMITTED/DONE/REJECTED) — operational, not attentional. A dormant idea in EBP terms maps to *no Conductor task exists* + belief at `entered` — which is fine, but RCP can't distinguish "dormant, deliberately parked" from "forgotten." Not worth new state — but the mapping table should say: **attention is workload management, owned by humans via Conductor task creation; it is not epistemic state.** One row in the mapping doc; prevents a future "add dormant flag to Solvent" request that would violate the kernel boundary.

**H3 — The 15.1–15.3 output format (next-move suggestions) belongs to the pack/agent layer, not the substrate.** §15's CLI prints "Next move: define the candidate invariant." If anyone implements that against Solvent, the ledger has acquired research-coaching semantics. Correct home: Domain Pack guidance documents and/or the RCP consumer's own reasoning — EBP's next-move doctrine is *advisory pedagogy for humans*, not kernel output. One sentence prevents scope creep.

## Medium

- **`needsOperatorSignoff`-class items:** EBP's six don't include a human-signoff item (faithfulness is the human gate); Plan 11.1's outcome — human obligations live in the gate map, not debt vocabulary — is actually *more* EBP-faithful than treating signoff as debt. Note it in the mapping doc.
- **Free entry, coerced debt:** EBP §2.1 says entry requires nothing. Your `ebpInitialDebt` auto-attaches six obligations at creation. This is stricter than EBP (§2.3 says debt doesn't kill, and EBP's freshIdea also attaches fullDebt — so actually §4.5 *conforms*; no issue, but worth one line in the mapping doc citing §4.5 rather than leaving it as an apparent deviation).
- **Owner field:** EBP's Idea has `owner`; your Solvent belief doesn't (provenance lives on evidence). Minor — provenance_class + decisions cover attribution; note in mapping.
- **Three-lane UI:** the mythology UI's Oracle/Passage views *are* Inbox/Workbench/Promoted projections — no work needed, just name the correspondence in the IA doc.

## The verdict

**The workflow is EBP v2.1-conformant at the structural level, is deliberately *stricter* than EBP in three places (no system inference, human-gated faithfulness, attributed retirement), and needs two recorded decisions (partial retirement, REOPEN lineage) plus one mapping document that makes all of this explicit.** The mapping document is the deliverable — one page: doctrine → mechanism → conformant/deviation-with-reason, covering the table above plus the C1/C2 decisions. It's also the artifact that lets you answer "are we compliant?" for every future reviewer with a citation instead of a shrug — and given that EBP v2.1's own §4.3 debt-opaqueness doctrine was just enforced into the kernel by Plan 11.1, recording *this* mapping completes a satisfying loop: the protocol document and the kernel now agree with each other, verified in both directions.

Nothing here blocks the Phase 8 dry run. The mapping doc and the two decision-record entries can land alongside it.


## answer

Yes. With the EBP v2.1 constraint added, I would make one final architectural consolidation before Phase 8:

> **RCP is a thin agent-facing read protocol over existing Conductor + Solvent state; EBP is the governing epistemic protocol; Solvent remains the enforcement point for debt and promotion; the UI is primarily a human debt-discharge and situational-awareness surface.**

EBP v2.1 is explicit that ideas may enter with essentially no prerequisites, but promotion is blocked by unpaid debt, and humans remain responsible for discharging that debt.  

## 1. How OpenCode should query RCP

I would use **API as the canonical protocol and MCP as the OpenCode-facing adapter**.

That gives us the least duplication:

```text
OpenCode
   │
   │ MCP tool
   ▼
ARGUS MCP adapter
   │
   │ HTTP
   ▼
Coordinator
   │
   ├── Conductor API
   └── Solvent API
```

So RCP itself is **not another server with another database**.

The canonical operation is conceptually:

```http
GET /v1/context/{task_id}
```

The MCP surface simply exposes that operation to OpenCode as something like:

```text
argus.get_context(task_id)
```

This is the right split because:

**API**

* canonical integration contract
* usable by UI, scripts, other agents, CI
* easy to test independently

**MCP**

* ergonomic agent interface
* lets OpenCode discover/use the capability as a tool
* no need to teach the agent your HTTP mechanics

And importantly, RCP should call **existing Conductor/Solvent interfaces**, not their databases directly.

```text
RCP
  → Conductor API
  → Solvent API
```

Only extend Conductor or Solvent when an existing interface genuinely cannot provide something required by the dry run.

### What RCP returns

Keep it extremely small:

```json
{
  "protocol": "argus-context/v1",
  "task": {},
  "dependencies": [],
  "epistemic": {
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  },
  "activity": []
}
```

Everything here is **projection**, not new state.

This directly addresses the six implementation plans' recurring concern that the research program must become queryable and that prior failures and dependencies must remain visible rather than being recreated from memory. The Claude plan in particular emphasizes turning the 47-section program into a queryable ledger and retaining failed work rather than rediscovering it.  

---

# 2. OpenCode workflow, finally simplified

### Work Agent

```text
Human
  ↓
create/select Conductor task
  ↓
OpenCode
  ↓
argus.get_context(task)
  ↓
do research
  ↓
produce EBP packet
  ↓
ARGUS
  ↓
Solvent + Conductor
```

### Adversarial Agent

Fresh OpenCode process, no conversational inheritance:

```text
OpenCode
  ↓
argus.get_context(task)
  ↓
sees:
  task
  dependencies
  completed work
  evidence
  current debt
  prior challenges
  dead ends
  activity
  ↓
attack unresolved work
  ↓
produce adversarial EBP packet
  ↓
ARGUS
  ↓
Solvent
```

The crucial property is:

> **The adversarial agent does not reconstruct history from files or conversation. It reconstructs the authoritative current state through RCP.**

That is the interface boundary you were after.

---

# 3. EBP changes the meaning of the workflow

We should make the loop explicitly:

```text
IDEA
 ↓
ENTER FREE
 ↓
WORK
 ↓
EVIDENCE
 ↓
ADVERSARIAL ATTACK
 ↓
DEBT
 ↓
HUMAN DISCHARGE
 ↓
PROMOTION
```

Not:

```text
agent does work → system decides it's good
```

The agent can **produce material that could retire debt**, but it does not itself discharge the debt.

That follows EBP's distinction between entry and promotion, and its statement that debt blocks promotion rather than killing the idea. 

---

# 4. The UI should now be only two primary surfaces

I agree with your sidebar concept. I would make it deliberately small:

```text
ARGUS
────────────────
INSIGHTS
DEBTS
────────────────
```

No mythology-heavy navigation for the first dry run. The mythology can remain in the visual language, but the actual POC should make the operational workflow obvious.

## A. Insights

This is the **research situation room**.

It answers:

> What have we learned, what is happening, and where should human attention go?

Top-level summary:

```text
BM–IST Synthesis v5.1

Ideas / Claims        12
Active Tasks           5
Open Debt              9
Adversarial Findings   3
Promoted               2
Refused                4
Dead Ends              2
```

Then the **attention Kanban**:

```text
HIGH PRIORITY
─────────────
G0 — Formalize proof
Gate 2 — Scaling audit

MEDIUM
──────
Canonical-height substrate

LOW
────
Deferred cosmology branch
```

But one important correction from the attached compliance review:

### Priority belongs to Conductor, not Solvent.

So the Insights board should derive task priority from **Conductor operational state**.

Do not create a Solvent `priority` field.

Likewise:

> **Attention state is workflow state, not epistemic state.**

That resolves the attached H2 concern cleanly.

---

# 5. Debts

This is the more important screen.

It should answer:

> **What does this idea still owe before it can be promoted?**

Example:

```text
G0 — Discrete Substrate / Continuous Unitary Compatibility

STATUS
Unpromoted

OPEN DEBT
────────────────────────
✓ needMap
✓ needInvariant
○ needToyCheck
○ needNullModel
○ needObstruction
○ needFaithfulnessReview
```

Click `needObstruction`:

```text
WHY THIS DEBT EXISTS

Adversarial finding:
The current formalization has not yet demonstrated
that the assumptions of the cited theorem apply to
the actual substrate object.

EVIDENCE
A-017
E-031
E-044

PROPOSED DISCHARGE
Formalize the missing assumption and provide a
counterexample attempt.

SOURCE
OpenCode / adversarial run #2
```

Then the human decides:

```text
[ RETIRE DEBT ]
[ KEEP OPEN ]
```

The actual retirement invokes the existing Solvent pathway.

That is exactly where your previous hackathon UI work is useful: **reuse the proven human debt-discharge interaction rather than inventing another mechanism.**

---

# 6. The human is the debt authority

This needs to be extremely explicit in Phase 8.

```text
Agent
  ↓
proposes evidence / argument
  ↓
Human
  ↓
decides whether debt is actually discharged
  ↓
Solvent
  ↓
records authoritative transition
```

So the UI should never show:

> “AI verified debt.”

It should show:

> “Adversarial agent supplied evidence relevant to `needObstruction`.”

Then:

> **Human decision: discharge debt**

That is much closer to EBP v2.1's intended philosophy. EBP describes debt retirement as a useful move that the researcher performs and says promotion only becomes available once current debt is retired. 

---

# 7. Consolidating the attached compliance findings

I would keep the useful findings, but simplify their resolution.

### C1 — Partial debt retirement

**Do not add partial-debt semantics.**

The current Solvent representation is binary:

```text
debt present
debt absent
```

Keep it that way for the POC.

Record explicitly:

> EBP v2.1 describes partial progress informally, but ARGUS v1 implements debt-item retirement as atomic. Partial progress remains evidence/prose until the entire debt item is discharged.

That is actually consistent with the broader "less is more" rule.

### C2 — System-inferred debt retirement

**Reject automatic retirement inference.**

This is a good deliberate strengthening.

The attached EBP text does describe seamless inference, but our authority architecture requires that semantic adjudication not be silently performed by an LLM at the Solvent boundary. EBP itself says accounting must not become the work. 

So:

```text
Agent suggests
Human adjudicates
Solvent records
```

No automatic NLP-based retirement.

### H1 — REOPEN

Keep the existing architecture:

```text
retracted belief
   ↓
new belief
   ↓
derives edge
```

but make the new belief **inherit relevant lineage/context**, not blindly become a new free-standing research idea.

That preserves the EBP living-ledger concept that an idea can accrue new debt after new evidence. 

### H2 — Attention

Resolve through Conductor.

No new Solvent field.

### H3 — “Next smallest useful move”

Keep it **UI/agent guidance**, never Solvent state.

That keeps the kernel domain-neutral.

---

# 8. One final refinement: RCP should expose debt, not prescribe how to discharge it

This is important.

RCP can say:

```text
OPEN DEBT:
needObstruction
```

It can expose:

```text
supporting evidence
prior attempts
adversarial findings
related tasks
dead ends
```

But it should not say:

```text
NEXT MOVE = run toy check X
```

That belongs to the **BM–IST Domain Pack / agent reasoning / human workflow**.

Otherwise RCP quietly becomes an epistemic advisor.

The attached compliance assessment was correct to identify this as a boundary issue.

---

# Final Phase 8 architecture

I would now freeze Phase 8 around this:

```text
                    HUMAN
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
      INSIGHTS                   DEBTS
          │                       │
          │                       ▼
          │                Human discharges
          │                       │
          └──────────┬────────────┘
                     ▼
                  SOLVENT
                     ▲
                     │
                COORDINATOR
               /             \
              ▼               ▼
        CONDUCTOR           SOLVENT
              ▲
              │
             RCP
              ▲
              │
       MCP / HTTP adapter
              ▲
              │
          OpenCode
        /            \
      WORK         ADVERSARY
```

And the protocol boundaries are:

```text
EBP
= epistemic operating law

Conductor
= operational work state

Solvent
= authoritative epistemic/debt/authorization state

RCP
= thin read projection

MCP
= agent-facing transport

API
= canonical integration interface

Trust UI
= human observation + adjudication
```

### Phase 8 dry-run acceptance

The first dry run should prove only this:

1. Human creates a BM–IST task in Conductor.
2. Work OpenCode retrieves context through MCP → RCP.
3. Work agent produces an EBP packet.
4. ARGUS records the resulting claims/evidence/debt.
5. Fresh adversarial OpenCode retrieves the same current state through RCP.
6. It attacks the unresolved work without rediscovering completed/dead-end work.
7. **Insights** shows the evolving research state and Conductor priorities.
8. **Debts** shows why each obligation remains open.
9. Human discharges at least one debt through the UI.
10. Solvent records the retirement.
11. Promotion is attempted and Solvent—not the UI or agent—decides whether the gate is satisfied.

That is thin enough to implement now, while exercising the actual architecture rather than building a decorative UI. EBP v2.1 itself calls for a minimum implementation centered on capture, status, retire, add-debt, and promote, while explicitly warning against dashboards and bureaucracy becoming the work. 

**I would consider this the final Phase 8 baseline.**


