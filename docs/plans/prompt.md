## prompt

I have BM-IST Synthesis on another thread, give prompt to it such that it delineates which work belong to which (e.g. delegate beliefs and other research items to AI/coding/research agents, which items are for Conductor and Solvent), the goal is division of labor as high level perspective




We are now using the BM-IST Synthesis as a concrete workload/POC to validate the three-way Pithom Labs architecture:

    AI / Coding / Research Agents = AGENCY
    Conductor = COORDINATION
    Solvent = AUTHORITY

with the hard invariant:

    CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

This is NOT a request to redesign BM-IST, Conductor, or Solvent.

The goal is to perform a HIGH-LEVEL DIVISION-OF-LABOR ANALYSIS of the existing BM-IST Synthesis.

Determine, for every major kind of work in the BM-IST program, which responsibility belongs to:

1. AI / coding / research agents
2. Conductor
3. Solvent
4. External executor / external system

The purpose is to make the boundary of responsibility obvious before implementation expands.

==================================================
CORE ARCHITECTURAL DEFINITIONS
==================================================

AI / Coding / Research Agents
--------------------------------
AGENCY.

Agents may:
- reason
- interpret domain material
- generate hypotheses
- formulate claims
- construct proofs
- search for counterexamples
- analyze literature
- write code
- run ordinary tools
- design experiments
- produce artifacts/results
- propose new work
- report findings
- request consequential actions

Agents do NOT gain authority merely because they are capable of performing an action.

Conductor
--------------------------------
COORDINATION.

Conductor owns:
- project
- task
- dependency
- assignment
- work lifecycle
- progress
- review workflow
- project activity/history

Conductor does NOT:
- determine scientific truth
- own domain ontology
- decide domain semantics
- authorize consequential actions
- execute consequential actions
- become a research engine
- become an agent framework
- become a policy engine

Solvent
--------------------------------
AUTHORITY.

Solvent answers:

    "May this exact consequential action happen against this exact state?"

Solvent owns:
- authority
- exact target/state binding
- authorization
- revocation
- intent
- claim
- consequence authorization boundary

Solvent does NOT:
- manage projects
- assign tasks
- manage agents
- manage dependencies
- decide scientific meaning
- become a research engine
- determine whether a task is complete

External executor / external system
--------------------------------
EFFECT.

The executor actually performs the external consequence.

Examples may include:
- GPU/compute execution
- external tool invocation
- deployment
- publication
- shared infrastructure mutation
- other real-world side effects

Important:

    AUTHORIZE ≠ EXECUTE

==================================================
TASK
==================================================

Read the entire BM-IST Synthesis and decompose its work at the highest useful level.

Do NOT simply repeat the existing BM-IST taxonomy.

Instead identify the actual categories of work an autonomous engineering/research system must perform.

For each major BM-IST activity, determine:

A. What the agent actually does
B. What Conductor needs to coordinate
C. When/why Solvent becomes relevant
D. What an external executor actually performs
E. What remains outside all three

Examples of BM-IST activity categories to examine include, but are not limited to:

- understanding the existing theory/model
- ingesting literature
- formulating hypotheses
- defining claims
- decomposing research questions
- constructing mathematical arguments
- writing formal proofs
- finding counterexamples
- designing experiments
- writing experiment code
- running ordinary/local computations
- running expensive/shared computations
- analyzing results
- comparing competing explanations
- identifying unresolved obligations
- semantic review
- revising hypotheses
- creating new research tasks
- reviewing another agent's work
- accepting/rejecting work
- publishing results
- modifying shared infrastructure
- modifying project code
- producing durable artifacts
- deciding whether an action is consequential
- requesting a consequential action
- recording what happened

Do not assume every activity above requires Solvent.

==================================================
IMPORTANT DISTINCTION
==================================================

Separate these concepts explicitly:

    DOMAIN MEANING
    WORK COORDINATION
    AUTHORITY
    EFFECT

For example, if the BM-IST system contains:

    hypothesis
    claim
    evidence
    proof
    attack
    debt
    dependency
    version
    status
    review

do NOT automatically map each one into Conductor or Solvent.

Determine whether each concept is:

1. something the agent/domain workload reasons about,
2. something Conductor needs to coordinate,
3. something Solvent needs to authorize,
4. something external infrastructure executes,
5. or something that does not need to become infrastructure at all.

==================================================
USE THE FOLLOWING ANALYTICAL MODEL
==================================================

For each major work category, produce a table:

| Work | Agent | Conductor | Solvent | External Executor | Why |
|------|-------|-----------|---------|-------------------|-----|

Use:
- PRIMARY
- SUPPORTING
- NOT RESPONSIBLE

where useful.

Then produce a second table:

| BM-IST Concept | Lowest Common Denominator | Owner | Why |
|----------------|---------------------------|-------|-----|

The "Lowest Common Denominator" column should deliberately reduce each concept to the smallest useful architectural responsibility.

Example:

    Hypothesis
      → domain reasoning/output
      → Agent/domain workload

    Investigate hypothesis
      → work item
      → Conductor

    Run expensive experiment against exact approved state
      → consequential action
      → Solvent

    Run experiment
      → external effect
      → Executor

Do NOT create infrastructure merely because a concept exists in BM-IST.

==================================================
KEY QUESTIONS
==================================================

Answer these explicitly:

1. Which BM-IST activities are purely agent/domain work?

2. Which activities need to become Conductor tasks?

3. Which activities should remain entirely inside the agent's own reasoning/tool environment?

4. Which actions constitute a meaningful boundary where Solvent should become relevant?

5. What makes an action "consequential" at the architecture level without making Conductor a policy engine?

6. Which information should Conductor record merely as project activity rather than treating it as authoritative evidence?

7. Which BM-IST concepts should NEVER be added to Conductor?

8. Which BM-IST concepts should NEVER be added to Solvent?

9. Which things should remain ephemeral agent state rather than durable system state?

10. Where does the actual external execution begin?

11. What should happen when an agent discovers new work while performing an existing task?

12. What should happen when an agent produces a proposed claim/evidence/result but it has not yet been reviewed?

13. What should happen when a task is coordinated by Conductor but the requested consequence is denied by Solvent?

14. What should happen when an external execution succeeds or fails?

15. What facts can each of the three systems legitimately claim to be authoritative about?

==================================================
BOUNDARY TEST
==================================================

For every proposed responsibility, perform this test:

Would the responsibility still make sense if the BM-IST domain were replaced with:

- ordinary Go software development
- cybersecurity research
- data engineering
- scientific simulation
- infrastructure operations

If not, it is probably domain/application logic rather than Conductor infrastructure.

Likewise:

Would Solvent need to understand the domain-specific meaning in order to answer:

    "May this exact consequential action happen?"

If not, do NOT put the domain semantics into Solvent.

==================================================
FINAL OUTPUT
==================================================

Produce these sections:

1. Executive Summary
2. Three-Way Responsibility Matrix
3. BM-IST Concept → Lowest Common Denominator Mapping
4. Agent Responsibilities
5. Conductor Responsibilities
6. Solvent Responsibilities
7. External Executor Responsibilities
8. What Must Remain Ephemeral / Domain-Specific
9. Example End-to-End BM-IST Workflow
10. Boundary Violations to Avoid
11. Architectural Conclusions

The final section must state explicitly:

- what belongs to Agent
- what belongs to Conductor
- what belongs to Solvent
- what belongs to external execution
- what should not become infrastructure at all

Do NOT propose new Conductor tables.
Do NOT propose new Solvent kernel primitives.
Do NOT redesign either architecture.
Do NOT turn BM-IST into a separate infrastructure layer.

The objective is only:

    FIND THE LOWEST COMMON DENOMINATORS
    AND DRAW THE CLEANEST DIVISION OF LABOR.

Use the existing BM-IST Synthesis as the source material, but reason from the
locked Pithom Labs architecture above.

