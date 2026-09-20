# 1. Executive Summary

I have calibrated this against the current BM-IST synthesis rather than the earlier, broader formulation. The present synthesis has a very specific frontier: most analyzed routes are excluded; the residual problem is the **Independent Realization Problem**—a deterministic, non-conditional amplitude realization with a local second-order unitary Schrödinger evolution, together with the universal action scale \(\kappa\).  

That scientific content should **not** be embedded into either Conductor or Solvent.

The clean Pithom Labs division is:

$$
\boxed{
\text{Agent}=\textbf{AGENCY}
\qquad
\text{Conductor}=\textbf{COORDINATION}
\qquad
\text{Solvent}=\textbf{AUTHORITY}
}
$$

with the fourth boundary:

$$
\boxed{\text{External system}=\textbf{EFFECT}}
$$

and the hard invariant:

$$
\boxed{
\text{CAPABILITY}\neq
\text{WORK}\neq
\text{AUTHORITY}\neq
\text{EXECUTION}.
}
$$

There is one important addition to that formulation:

$$
\boxed{\text{DOMAIN MEANING is not owned by any of the three infrastructure layers.}}
$$

The agent can reason about domain meaning and produce candidate interpretations; the domain application/verifier determines their semantic status. This follows the existing architecture's separation between domain semantics and Solvent's authority substrate. 

The resulting model is therefore:

```text
                 DOMAIN / BM-IST
                 meaning & truth
                       │
                       ▼
                 AI / AGENTS
                    AGENCY
                       │
                 proposed work
                       ▼
                 CONDUCTOR
                  COORDINATION
                       │
             consequential intent
                       ▼
                  SOLVENT
                   AUTHORITY
                       │
                    allowed
                       ▼
                   EXECUTOR
                    EFFECT
```

But those are **responsibility boundaries**, not necessarily a serial request path. Agent ↔ Conductor is continuous; Agent/Executor can meet Solvent at the consequence boundary. 

---

# 2. Three-Way Responsibility Matrix

The most useful way to classify BM-IST is by **kind of work**, not by its existing research taxonomy.

| Work                                          | Agent                     | Conductor                            | Solvent                                       | External Executor             | Why                                                                                          |
| --------------------------------------------- | ------------------------- | ------------------------------------ | --------------------------------------------- | ----------------------------- | -------------------------------------------------------------------------------------------- |
| Understand existing BM-IST                    | **PRIMARY**               | NOT RESPONSIBLE                      | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain reasoning                                                                             |
| Ingest/read literature                        | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Agent interprets sources                                                                     |
| Extract findings from literature              | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Semantic work                                                                                |
| Formulate hypothesis                          | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain reasoning                                                                             |
| Define scientific claim                       | **PRIMARY**               | NOT RESPONSIBLE                      | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain semantics                                                                             |
| Build proof/argument                          | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | Sometimes supporting          | Reasoning, formalization                                                                     |
| Find counterexample/attack                    | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | Sometimes supporting          | Adversarial reasoning                                                                        |
| Identify open obligation                      | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain determines what remains unresolved                                                    |
| Decompose work into tasks                     | SUPPORTING                | **PRIMARY**                          | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Coordination                                                                                 |
| Assign agent                                  | NOT RESPONSIBLE           | **PRIMARY**                          | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Work ownership                                                                               |
| Track dependency                              | SUPPORTING                | **PRIMARY**                          | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Work graph                                                                                   |
| Track progress                                | SUPPORTING                | **PRIMARY**                          | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Work state                                                                                   |
| Review work lifecycle                         | **PRIMARY** for substance | **PRIMARY** for workflow             | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Separate semantic review from workflow review                                                |
| Ordinary local computation                    | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | **PRIMARY** if external tool  | Routine agency/tool use                                                                      |
| Write experiment code                         | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Agent produces artifact                                                                      |
| Run local experiment                          | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | SUPPORTING                    | No consequential boundary necessarily exists                                                 |
| Run expensive/shared experiment               | **PRIMARY**               | **PRIMARY** for work state           | **PRIMARY** for consequential authorization   | **PRIMARY**                   | Resource/effect boundary                                                                     |
| Modify shared repository                      | **PRIMARY**               | SUPPORTING                           | **PRIMARY** when consequential                | **PRIMARY**                   | External mutation                                                                            |
| Modify production/shared infrastructure       | PRIMARY                   | SUPPORTING                           | **PRIMARY**                                   | **PRIMARY**                   | Clear consequence                                                                            |
| Analyze results                               | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | SUPPORTING                    | Domain interpretation                                                                        |
| Compare competing explanations                | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Scientific reasoning                                                                         |
| Revise hypothesis                             | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain meaning                                                                               |
| Create newly discovered work                  | **PRIMARY** proposes      | **PRIMARY** records/coordinates      | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Agency discovers; Conductor operationalizes                                                  |
| Accept/reject scientific conclusion           | **PRIMARY** proposes      | NOT RESPONSIBLE for truth            | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Domain/human judgment                                                                        |
| Accept/reject task completion                 | SUPPORTING                | **PRIMARY** workflow state           | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Completion is work state, not authority                                                      |
| Decide whether action is consequential        | SUPPORTING                | NOT RESPONSIBLE as policy            | **PRIMARY only as authority boundary**        | NOT RESPONSIBLE               | Classification/policy belongs above infrastructure; Solvent enforces the resulting authority |
| Request consequential action                  | **PRIMARY**               | SUPPORTING                           | RECEIVES                                      | NOT RESPONSIBLE               | Agent initiates intent                                                                       |
| Determine authorization                       | NOT RESPONSIBLE           | NOT RESPONSIBLE                      | **PRIMARY**                                   | NOT RESPONSIBLE               | Solvent's core job                                                                           |
| Execute authorized action                     | NOT RESPONSIBLE           | NOT RESPONSIBLE                      | NOT RESPONSIBLE                               | **PRIMARY**                   | Effect belongs outside Solvent                                                               |
| Record actual external outcome                | REPORTS                   | **PRIMARY** as activity/work history | May retain authority/execution facts it needs | **PRIMARY source of outcome** | Execution fact comes from executor                                                           |
| Publish result                                | PRIMARY                   | SUPPORTING                           | **PRIMARY only if publication is governed**   | **PRIMARY**                   | Publication is external consequence                                                          |
| Maintain scientific provenance                | PRIMARY/domain            | SUPPORTING                           | NOT RESPONSIBLE                               | SUPPORTING                    | Domain provenance                                                                            |
| Maintain authorization provenance             | NOT RESPONSIBLE           | OBSERVES                             | **PRIMARY**                                   | SUPPORTING                    | Authority provenance                                                                         |
| Produce durable research artifact             | **PRIMARY**               | SUPPORTING                           | NOT RESPONSIBLE                               | Storage system                | Artifact content is agent/domain output                                                      |
| Decide what artifact means                    | **PRIMARY/domain**        | NOT RESPONSIBLE                      | NOT RESPONSIBLE                               | NOT RESPONSIBLE               | Semantics                                                                                    |
| Decide whether artifact may cause consequence | SUPPORTING                | SUPPORTING                           | **PRIMARY**                                   | NOT RESPONSIBLE               | Authority question                                                                           |
| Actual mutation of world state                | NOT RESPONSIBLE           | NOT RESPONSIBLE                      | NOT RESPONSIBLE                               | **PRIMARY**                   | Effect                                                                                       |

The key distinction is between **scientific acceptance** and **workflow acceptance**. Conductor can mark a task `COMPLETED` because the assigned work was delivered; it must not thereby declare the scientific proposition true. The research/domain layer owns that semantic judgment. This is consistent with the existing separation of domain meaning, work state, and authority state. 

---

# 3. BM-IST Concept → Lowest Common Denominator Mapping

The BM-IST synthesis contains a large epistemic apparatus: axioms, definitions, hypotheses, claims, evidence, proofs, attacks, obligations, versions, statuses, Gate 3 formulations, and so forth. The mistake would be turning all of these into infrastructure entities.

Here is the reduction.

| BM-IST Concept                  | Lowest Common Denominator                | Owner                                   | Why                                                         |
| ------------------------------- | ---------------------------------------- | --------------------------------------- | ----------------------------------------------------------- |
| Axiom                           | **Domain proposition**                   | Agent/domain                            | Meaning is domain-specific                                  |
| Definition                      | **Domain proposition**                   | Agent/domain                            | Not infrastructure                                          |
| Hypothesis                      | **Candidate proposition**                | Agent/domain                            | Agent reasons over it                                       |
| Claim                           | **Proposition/output**                   | Agent/domain                            | Scientific meaning                                          |
| Theorem                         | **Verified proposition**                 | Domain verifier                         | Truth, not work coordination                                |
| Evidence                        | **Supporting artifact**                  | Agent/domain                            | Conductor may record activity referring to it               |
| Proof                           | **Verification artifact**                | Agent/domain/external verifier          | Proof semantics are domain-specific                         |
| Attack                          | **Adversarial artifact**                 | Agent/domain                            | Counterexample reasoning                                    |
| Counterexample                  | **Evidence against claim**               | Agent/domain                            | Semantic significance                                       |
| Debt                            | **Unresolved obligation**                | Domain/application                      | Conductor may coordinate its work; Solvent need not know it |
| Dependency                      | **Work dependency or domain dependency** | Depends                                 | Important distinction below                                 |
| Version                         | **Identity of a revision**               | Domain/application / Conductor metadata | Don't collapse all version concepts into one                |
| Stale                           | **Domain state requiring reassessment**  | Domain/application                      | Not a Conductor authority state                             |
| Semantic review                 | **Review obligation**                    | Domain/application + human              | Conductor can coordinate it                                 |
| Proof-pending                   | **Work/status metadata**                 | Domain/application                      | Not Solvent                                                 |
| Proved                          | **Domain verification status**           | Domain verifier                         | Not authority                                               |
| Refuted                         | **Domain verification status**           | Domain verifier                         | Not task failure                                            |
| Restricted                      | **Domain qualification**                 | Domain verifier                         | Not infrastructure state                                    |
| Research target                 | **Work objective**                       | Conductor                               | Once operationalized                                        |
| Research question               | **Work objective**                       | Agent proposes / Conductor coordinates  | Meaning remains domain-side                                 |
| Experiment                      | **Work activity + possibly consequence** | Agent/Conductor                         | Execution only crosses boundary when consequential          |
| Research snapshot               | **Domain/work state snapshot**           | Domain/application                      | Do not equate it automatically with Solvent snapshot        |
| Gate 3                          | **Research objective/workstream**        | Conductor coordinates                   | Mathematical meaning remains domain-side                    |
| Independent Realization Problem | **Research objective**                   | Agent/domain + Conductor                | Conductor tracks it; doesn't understand it                  |
| \(\mathcal A\) nonlinear lift   | **Candidate mathematical artifact**      | Agent/domain                            | Core scientific content                                     |
| \(\kappa\) hypothesis           | **Domain claim/obligation**              | Agent/domain                            | Not a Solvent primitive                                     |
| PYV transfer test               | **Research task**                        | Conductor                               | Its mathematical content is domain-side                     |
| T2/T4                           | **Research hypothesis branches**         | Agent/domain                            | Conductor only coordinates work around them                 |
| Koopman spectrum calculation    | **Research task/activity**               | Agent + Conductor                       | Execution may be local or governed                          |
| Publication                     | **External consequence**                 | Solvent + executor when governed        | Not merely a research status                                |

The particularly important distinctions are:

$$
\boxed{\text{domain dependency}\neq\text{work dependency}}
$$

and

$$
\boxed{\text{research version}\neq\text{authority snapshot}}.
$$

Conductor needs to know that Task B depends on Task A. It does not need to know the mathematical semantics of the theorem connecting them.

---

# 4. Agent Responsibilities

The agent is where the **intellectual work** happens.

For the actual BM-IST program, that means the agent may:

```text
read BM-IST material
→ reconstruct the current argument
→ inspect prior results
→ identify contradictions
→ formulate T2/T4 hypotheses
→ derive proof obligations
→ construct mathematical arguments
→ search literature
→ write simulations
→ analyze spectra
→ attack assumptions
→ generate candidate constructions
→ report findings
→ propose next work
```

This includes the hard scientific content of the present synthesis.

For example, the synthesis says the current residual object is the nonlinear, non-conditional amplitude realization satisfying S1–S10. 

The agent is responsible for working on questions such as:

> Can \(\mathcal A[I_U,\Phi_t,\mu,\pi]\) exist?

> Can it satisfy the locality requirement?

> Can it provide unitary second-order evolution?

> Can the action unit \(\kappa\) emerge?

> Does a particular toy model actually satisfy any candidate structural hypotheses?

These are **agency/domain questions**, not Conductor questions.

### What the agent must not infer

The agent finding a proof does not make the claim authoritative.

The agent finding a plausible experiment does not make the experiment approved.

The agent being capable of modifying the repository does not mean the work is assigned.

The agent reporting:

> “The experiment succeeded”

does not itself constitute the authoritative execution fact.

That preserves:

$$
\boxed{\text{agent output}\neq\text{truth}\neq\text{authority}.}
$$

---

# 5. Conductor Responsibilities

Conductor should take the BM-IST intellectual program and turn **work** into an operational structure.

For example, from the current Gate 3 program:

```text
Project
  BM-IST Synthesis

Task
  Investigate Independent Realization Problem

Tasks
  ├── characterize candidate amplitude realizations
  ├── attack non-conditional lift candidates
  ├── test T2 factor hypotheses
  ├── test T4 hierarchical hypothesis
  ├── analyze toy model
  ├── run spectral calculations
  ├── review results
  └── update research direction
```

The synthesis itself is clear that the residual program is centered on Gate 3 and the Independent Realization Problem, with supporting work around the action quantum. 

Conductor owns:

$$
\boxed{
\text{What work exists?}
}
$$

$$
\boxed{
\text{Who is doing it?}
}
$$

$$
\boxed{
\text{What must happen first?}
}
$$

$$
\boxed{
\text{What is active, blocked, reviewed, or complete?}
}
$$

It may record:

```text
Task 84
Assigned Agent X
Status ACTIVE
Depends on Task 71
Activity: spectral-analysis-result-17
```

But it must not conclude:

```text
Task 84 completed
therefore
BM-IST T4 is true
```

That crosses from **work state** into **domain truth**.

The existing Conductor responsibility set—project, task, dependency, assignment and activity—is therefore sufficient. 

---

# 6. Solvent Responsibilities

Solvent should enter **only when work becomes a consequential action**.

This is the most important boundary.

Running:

```text
python analyze_spectrum.py
```

on an agent's local machine is normally agency/tool work.

Launching a $5,000 GPU experiment against shared infrastructure is different.

Deleting a shared dataset is different.

Pushing to a protected production repository is different.

Publishing an external paper is potentially different.

Changing shared Conductor/Solvent infrastructure is different.

Solvent asks only:

$$
\boxed{
\text{May this exact consequence happen against this exact state?}
}
$$

It does not ask:

> Is T4 scientifically correct?

It does not ask:

> Is this a good research direction?

It does not ask:

> Has Agent X done enough work?

It does not understand PYV, Koopman theory, Fisher information, or BM-IST.

The Solvent boundary remains exact target/state binding, authorization, revocation, intent and claim. 

So:

```text
Agent:
    "Run experiment 42 on shared GPU cluster."

Conductor:
    "This is Task 84, assigned and dependency-complete."

Solvent:
    "Is experiment-42 against this exact target/state authorized?"

Executor:
    "GPU job actually ran."
```

This is the cleanest application of:

$$
\boxed{\text{AUTHORIZE}\neq\text{EXECUTE}.}
$$

---

# 7. External Executor Responsibilities

The executor is where **the consequence actually occurs**.

For BM-IST that may be:

```text
local process
GPU cluster
Lean/Coq farm
CAS
simulation environment
shared repository
CI/CD
publication system
cloud infrastructure
database
```

But the important distinction is not the technology.

It is:

$$
\boxed{\text{Does this operation alter state beyond the agent's ephemeral environment?}}
$$

If yes, it may be an external consequence.

For example:

### Local symbolic calculation

```text
Agent → Python/SymPy
```

No Solvent necessary.

### Shared expensive computation

```text
Agent
 → Conductor task
 → Solvent authorization
 → GPU executor
```

### Code modification in an ordinary working branch

May remain agent work, depending on the environment.

### Merge into protected repository

Potentially:

```text
Agent
 → Conductor
 → Solvent
 → Git executor
```

### Publication

Potentially:

```text
Agent produces manuscript
 → Conductor coordinates review
 → Solvent authorizes publication
 → publication executor
```

This keeps Solvent from becoming an execution platform. The existing architecture explicitly separates authorization from actual external execution. 

---

# 8. What Must Remain Ephemeral / Domain-Specific

This is where we prevent architecture creep.

## Ephemeral agent state

These normally do **not** need durable infrastructure representation:

```text
chain-of-thought / private reasoning
scratch calculations
temporary hypotheses
temporary code experiments
temporary tool output
intermediate algebra
local searches
temporary files
discarded candidate ideas
```

The agent may persist artifacts when useful, but persistence is not itself a reason to make something a Conductor or Solvent primitive.

## Domain-specific state

These should remain in the BM-IST research/application layer:

```text
axiom
definition
hypothesis
claim
theorem
proof
counterexample
attack
Fisher condition
Koopman spectrum
spectral type
T2
T4
PYV compatibility
Independent Realization Problem
S1–S10
Gate status
mathematical validity
scientific interpretation
```

The current synthesis makes the distinction especially important because its survivor map and S1–S10 specification contain substantial mathematical semantics. 

### Things that should not become infrastructure merely because they're durable

A PDF containing a proof is an artifact.

A Git commit containing an experiment is an artifact.

A numerical result is an artifact.

A paper is an artifact.

A hypothesis is domain state.

None becomes a Conductor primitive merely because it is persisted.

---

# 9. Example End-to-End BM-IST Workflow

Take the current high-value research question:

> Test whether the Phase-2 toy model exhibits the structural regime relevant to T4: hierarchical microscopic structure producing emergent coarse-grained behavior.

The current BM-IST material explicitly describes this sort of multi-resolution experiment—coarse Lyapunov behavior, fine-scale structure, hierarchical/limit-periodic signatures, comparison across resolutions. 

### Step 1 — Agent discovers the question

Agent reasons:

```text
The current model may distinguish
fine-scale hierarchy from coarse-scale chaos.

I propose a multi-resolution experiment.
```

This is **AGENCY**.

### Step 2 — Conductor operationalizes it

Conductor records:

```text
Project: BM-IST Phase 2

Task: Investigate T4

Dependencies:
    Phase-2 model ready
    measurement code ready

Assignment:
    Agent-X

Subtasks:
    coarse Lyapunov analysis
    fine spectral analysis
    hierarchy test
    cross-resolution comparison
    review
```

This is **COORDINATION**.

### Step 3 — Agent performs ordinary work

Agent writes code, executes local simulations, inspects results.

No Solvent is required merely because the work concerns advanced mathematics.

This is **AGENCY**.

### Step 4 — Agent discovers a consequential experiment

Agent determines:

```text
Need 500 GPU-hours on shared cluster.
```

Agent proposes:

```text
RUN_T4_MULTIRESOLUTION_ANALYSIS
```

This is still agency.

### Step 5 — Conductor checks work state

Conductor confirms:

```text
Task active
Dependencies satisfied
Assigned agent valid
```

Conductor does **not** authorize the GPU expenditure.

### Step 6 — Solvent handles authority

Solvent evaluates:

```text
exact target
exact state
requested consequence
current authority
revocation state
```

and either:

```text
AUTHORIZED
```

or:

```text
DENIED
```

### Step 7 — Executor acts

If authorized:

```text
GPU executor → runs experiment
```

That is the actual external effect.

### Step 8 — Result returns

```text
Executor result
      ↓
Agent interprets
      ↓
Artifact/result produced
      ↓
Conductor records activity
      ↓
Agent proposes next task
```

The agent might conclude:

```text
T4 evidence is suggestive but not sufficient.
```

That conclusion remains a **domain claim**.

Conductor records that the work occurred.

Solvent has no need to know what T4 means.

---

# 10. Boundary Violations to Avoid

These are the failure modes I would use as architectural tests.

### Conductor becomes a research engine

Bad:

```text
Conductor:
    determine whether Fisher rigidity holds
```

Correct:

```text
Agent/domain:
    determine whether Fisher rigidity holds

Conductor:
    coordinate the task investigating Fisher rigidity
```

The research engine is domain/application work, not infrastructure. 

### Conductor becomes an authorization gate

Bad:

```text
Task = approved
therefore
action = authorized
```

Wrong.

A task is work state.

Authority is a different state.

$$
\boxed{\text{WORK}\neq\text{AUTHORITY}}
$$

### Solvent becomes a scientific judge

Bad:

```text
Solvent checks whether T4 is scientifically valid.
```

Correct:

```text
domain verifier checks T4
Solvent checks whether
the resulting consequential action is authorized.
```

### Solvent becomes an execution engine

Bad:

```text
Solvent runs the GPU experiment.
```

Correct:

```text
Solvent authorizes.
GPU executor executes.
```

### Agent capability becomes work assignment

Bad:

```text
Agent can modify repository
→ therefore it is allowed to modify repository.
```

Wrong.

$$
\boxed{\text{CAPABILITY}\neq\text{WORK}}
$$

### Work assignment becomes authority

Bad:

```text
Task 84 is assigned
→ therefore Agent X may spend \$5,000 of shared compute.
```

Wrong.

### Activity becomes evidence

Bad:

```text
Agent reports "proved theorem"
→ Conductor activity proves theorem.
```

Activity says what workflow activity occurred; it does not establish mathematical truth.

### Evidence becomes authority

Bad:

```text
Agent found an approval in a document
→ authorized.
```

This is exactly the kind of confused-deputy/agentjacking failure Solvent is intended to prevent. 

### Research snapshot becomes authority snapshot

Bad:

```text
research-state-042 = Solvent authority snapshot
```

They may be linked, but they represent different kinds of state.

---

# 11. Architectural Conclusions

The BM-IST POC actually provides an excellent stress test because it contains nearly every type of work an autonomous engineering system will encounter:

$$
\text{reasoning}
\rightarrow
\text{planning}
\rightarrow
\text{coordination}
\rightarrow
\text{computation}
\rightarrow
\text{review}
\rightarrow
\text{consequential action}
\rightarrow
\text{external effect}.
$$

The division should therefore be frozen as:

## AI / Coding / Research Agent — **AGENCY**

Owns:

$$
\boxed{
\text{reason}
,\text{interpret}
,\text{propose}
,\text{investigate}
,\text{implement}
,\text{analyze}
,\text{report}
}
$$

It produces the domain work.

## Conductor — **COORDINATION**

Owns:

$$
\boxed{
\text{project}
,\text{task}
,\text{dependency}
,\text{assignment}
,\text{lifecycle}
,\text{activity}
}
$$

It makes work operationally coherent.

The current four-table Conductor design is therefore sufficient; nothing in the BM-IST synthesis creates a demonstrated need for additional Conductor primitives.

## Solvent — **AUTHORITY**

Owns:

$$
\boxed{
\text{intent}
,\text{target}
,\text{state binding}
,\text{authorization}
,\text{revocation}
,\text{claim}
}
$$

It answers exactly one question:

$$
\boxed{
\text{May this consequential action happen now, against this exact state?}
}
$$

The existing Solvent architecture was deliberately designed around that narrow authority function. 

## External Executor — **EFFECT**

Owns:

$$
\boxed{\text{actual external state change}.}
$$

It runs the computation, modifies the repository, consumes the GPU allocation, publishes the artifact, mutates infrastructure, or otherwise creates the external consequence.

And everything else belongs **outside these infrastructure boundaries**:

$$
\boxed{
\text{scientific truth}
,\text{domain semantics}
,\text{proof meaning}
,\text{hypothesis status}
,\text{epistemic interpretation}
,\text{policy meaning}
}
$$

Those remain in the agent/domain/application/verifier world.

The final Pithom Labs model is therefore:

```text
                 ┌─────────────────────────┐
                 │       DOMAIN WORK       │
                 │  BM-IST / Go / Security │
                 │  Science / Operations   │
                 └───────────┬─────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │      AGENT      │
                    │     AGENCY      │
                    │ reason / do     │
                    └───────┬────────┘
                            │
                       work state
                            │
                            ▼
                    ┌────────────────┐
                    │   CONDUCTOR    │
                    │  COORDINATION  │
                    │ work / assign  │
                    │ depend / track │
                    └───────┬────────┘
                            │
                   consequential intent
                            │
                            ▼
                    ┌────────────────┐
                    │    SOLVENT     │
                    │    AUTHORITY   │
                    │ bind / permit  │
                    │ revoke / claim │
                    └───────┬────────┘
                            │
                         authorized
                            │
                            ▼
                    ┌────────────────┐
                    │    EXECUTOR    │
                    │     EFFECT    │
                    └───────┬────────┘
                            │
                            ▼
                      EXTERNAL WORLD
```

The deepest architectural invariant is thus slightly richer than the three-word slogan:

$$
\boxed{
\textbf{CAPABILITY}
\neq
\textbf{WORK}
\neq
\textbf{AUTHORITY}
\neq
\textbf{EFFECT}
}
$$

while **domain meaning cuts across the system without becoming part of any of those infrastructure responsibilities**.

That is exactly what makes BM-IST a useful POC: it is complex enough to exercise agency, coordination, verification, expensive computation, artifact production, review, and consequential execution—yet the underlying responsibilities can remain cleanly separated without adding a BM-IST-specific primitive to either Conductor or Solvent.

