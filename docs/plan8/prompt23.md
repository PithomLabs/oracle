No—not from the current prompt alone.

The adversarial prompt says:

> “Call `argus.get_context` with the task ID”
> “Reconstruct what has already been done from the system response”
> “Identify unresolved work, weak assumptions, or unsupported claims”
> “Use `contradicts` edges to link challenges to existing beliefs.”  

That establishes the mechanism, but it does **not** require:

> enumerate every belief returned by `get_context` and explicitly assess each one.

And importantly, the adversarial agent should **not** read the database itself. Its authoritative claim set is whatever ARGUS returns through `get_context`.

### The workflow I would expect

For your dry run, I would expect:

```text
Fresh Adversarial Agent
        |
        v
argus.get_context(task_id)
        |
        v
ARGUS returns authoritative research context
        |
        +--> existing beliefs
        +--> evidence
        +--> debt
        +--> edges
        +--> task/dependencies
        |
        v
Agent reads docs/corpus/*
        |
        v
Build an internal claim inventory
        |
        +--> B1: attack
        +--> B2: attack
        +--> B3: no material objection
        +--> B4: insufficient basis
        ...
        |
        v
Produce adversarial findings
        |
        v
Create new challenge beliefs where needed
        |
        v
contradicts edge → specific existing belief ID
        |
        v
argus.submit_packet
```

The critical point is that **the target remains the existing immutable belief ID**. The adversarial agent should not rewrite it. Your shared protocol already states that existing beliefs remain historical propositions and materially different propositions get new IDs. 

### “All claims in the database” is not quite the right goal

I would not make the adversarial agent attack literally every belief in the entire database.

The correct target is:

> **all relevant beliefs in the authoritative research context for the assigned task/scenario that are visible through `argus.get_context`.**

That preserves the black-box architecture.

So the important question is actually:

> **What exactly does the current `argus.get_context(task_id)` return?**

That is the part I would audit before you run the adversarial agent.

Because there are two possible implementations:

```text
A.
get_context(task)
    → all beliefs in current scenario
```

or:

```text
B.
get_context(task)
    → only beliefs directly associated with task/project
```

Those produce very different adversarial coverage.

Your current dry-run evidence showing multiple beliefs becoming visible through successive `get_context` calls is encouraging, but it doesn't by itself prove that **every relevant belief** is included.

## I would run this audit first

Give the coding agent this exact prompt:

```text id="6j3d7n"
# ARGUS Audit — Adversarial Context Coverage

Do NOT change code.

I am about to run the first live Adversarial Agent dry run and need to know
exactly what an adversarial agent can see through:

    argus.get_context(task_id)

The adversarial role is intended to challenge the current research state
without direct database access.

Answer this specific question:

> Does `argus.get_context(task_id)` return all existing relevant beliefs in
> the current research scenario, or only a task/project-scoped subset?

Trace the complete production path:

1. MCP `argus.get_context`
2. MCP adapter
3. application `GetContext`
4. epistemic snapshot/view functions
5. work/task lookup
6. scenario/project resolution
7. database queries
8. any filtering/traversal
9. final JSON returned to the agent

Do not infer from comments or README. Inspect the actual implementation.

Report:

### A. Belief coverage

Given a scenario containing:

B1
B2
B3
...
Bn

and an arbitrary task T in that scenario:

- which beliefs are returned?
- are all scenario beliefs returned?
- only project beliefs?
- only task-associated beliefs?
- only beliefs reachable through some edge?
- is there any status filtering?
- are retracted beliefs returned?
- are contradicted beliefs returned?
- are newly submitted adversarial beliefs returned?

### B. Exact SQL/filtering

Show the actual relevant query/filter logic.

State precisely what determines inclusion/exclusion.

### C. Adversarial visibility

Can a fresh adversarial agent, using ONLY:

    argus.get_context(task_id)

see every belief it is expected to challenge in the current research
scenario?

Do NOT assume the answer is yes.

### D. Dry-run example

Using the current seeded scenario and the actual implementation, state:

- number of beliefs currently persisted
- number returned by get_context
- IDs returned
- whether all are visible to the adversarial agent

### E. Coverage gap

If not all relevant beliefs are returned, identify exactly what is hidden
and why.

Do NOT propose a new MCP tool.

Do NOT propose direct database access.

Do NOT change implementation yet.

### F. Recommendation

Give one of:

PASS:
`get_context` provides the intended adversarial claim surface.

PARTIAL:
`get_context` provides a subset and the adversarial workflow must explicitly
operate within that subset.

FAIL:
The adversarial agent cannot reconstruct the relevant current research state
through `get_context`.

The goal is to establish the actual claim surface available to a fresh
Adversarial Agent before the live dry run.
```

### My expectation for the actual dry run

I would **not** expect the adversarial agent to blindly generate a contradiction for every belief.

A good adversarial run should produce something more like:

```text
Claim inventory

B1 — G0
    Attack: YES
    Reason: ...
    Severity: ...

B2 — Deliverable A substrate
    Attack: YES
    Reason: ...

B3 — encoding definition
    Attack: NO MATERIAL OBJECTION
    Reason: ...

B4 — boundary convention
    Attack: UNCERTAIN
    Reason: insufficient basis
```

Then only the actual challenges become new beliefs with:

```text
B_new --contradicts--> B_target
```

That is more valuable than artificially forcing one attack per claim.

And the adversarial agent should **never decide that a claim is "validated"** simply because it failed to find an objection. Your shared protocol already establishes that absence of a counterexample is not proof. 

So before you launch the adversarial run, I would do the audit above. It will tell you exactly what the agent's field of view is, which is the key prerequisite for interpreting the dry-run result.
