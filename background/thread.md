# DeepWiki Q&A with Code Context for Repository: PithomLabs/solvent
## Q1
using article below, how would design a Solvent-ready agent harness that uses Solvent primitives and MCP such that the agent harness uses a DAG of beliefs, intents etc (same Solvent primitives and data structures), the goal is to enable a sort of oracle verifier regardless of domain and industry: ## prompt

given attached how do you propose to implement DAG-based lessons and insights to real-world software engineering, what principles are applicable here, that DAG model is good only as a mental model for humans (and applicable to AI agents sort of), that's because software is inherently code-based and hopelessly entangled with one another, how do you reconcile the complexity of software code vs the utopia of DAG for mental models, how do you propose to evaluate DAG as review gates then, as in "you do not review the code, you review the gates." in the grand scheme of things in the world of software engineering


## answer

Yes. I think the key reconciliation is this:

**The DAG should not attempt to model the software. It should model the reasoning, obligations, dependencies, and evidence around the software.**

That distinction makes the idea much less utopian.

The attached piece's strongest lesson is not really “use DAGs.” It is that **externalized, shared state plus explicit dependencies and mechanical verification allowed otherwise-capable agents to coordinate without losing project state**. Prove2Me did this by making theorem statements, dependencies, descriptions, and verification state explicit rather than leaving them inside agent context. 

For software engineering, I would generalize that into something like:

> **Code is the implementation graph.
> The DAG is the reasoning-and-evidence graph.**

## 1. Why software cannot actually be a clean DAG

Your objection is exactly right.

Real software is full of:

* cyclic dependencies
* shared mutable state
* implicit contracts
* runtime behavior
* temporal coupling
* distributed effects
* emergent behavior
* configuration dependencies
* backward compatibility constraints
* tests that depend on fixtures that depend on infrastructure that depends on configuration
* code whose meaning depends on code elsewhere

A graph of the actual system is therefore not necessarily acyclic, and even when you can technically construct a dependency DAG, it can become so enormous that it stops being useful.

So I would **not** make “DAG-based software engineering” mean:

> “Represent the entire codebase as a DAG.”

That is the wrong abstraction.

Instead:

> **Represent the work required to establish confidence in a change as a DAG.**

That graph can remain relatively small even when the underlying software is monstrously complicated.

---

# 2. The DAG is a map of *claims*, not code

Suppose an engineer says:

> “I changed pagination so numeric navigation cannot skip a page.”

The code may touch 17 files and interact with rendering, retries, state settlement, persistence, browser execution, and parsing.

Trying to make one DAG represent all of that code is horrible.

Instead, create a DAG of **claims that must become true**:

```text
                    Change
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
   Pagination invariant     Retry invariant
          │                       │
          ▼                       ▼
   Advancement correctness   Retry termination
          │                       │
          └───────────┬───────────┘
                      ▼
                Integration
                      │
                      ▼
                  Regression
                      │
                      ▼
                 Release gate
```

The code can remain tangled underneath.

The DAG says:

> “These are the things that must be demonstrated before we consider the change trustworthy.”

That is dramatically more scalable.

---

# 3. So I would actually use **multiple graphs**

One DAG is still too simplistic.

I'd separate at least four structures:

### A. Implementation graph

The actual software.

```text
packages → modules → functions → state → runtime
```

Messy. Cyclic. Large.

### B. Obligation DAG

What must be true.

```text
requirement
   ↓
invariant
   ↓
behavior
   ↓
acceptance condition
```

Mostly acyclic.

### C. Evidence graph

How do we know it is true?

```text
invariant
   ↓
unit test
   ↓
integration test
   ↓
production telemetry
   ↓
observed stability
```

Also mostly acyclic.

### D. Knowledge graph

What did we learn and how does it relate to previous knowledge?

```text
lesson A ──refines──> lesson B
lesson C ──contradicts──> lesson D
lesson E ──applies-to──> invariant X
lesson F ──derived-from──> incident Y
```

This is **not necessarily a DAG**.

And that's important.

The attached piece is strongest when it talks about externalized state, searchability, and reuse. 

So I'd actually say:

> **Use a DAG for execution and verification. Use a graph for knowledge.**

That is a much more faithful model of software engineering.

---

# 4. What is a “lesson” in this system?

I would strongly avoid lessons like:

> “Always write tests.”

That is too generic and almost useless to an agent.

Instead, lessons should be **compressed experience attached to a concrete invariant**.

For example:

```text
INCIDENT
Numeric pagination skipped page 3.

        ↓

OBSERVATION
Page advancement was validated using fallback
content rather than the primary pagination identity.

        ↓

ROOT CAUSE
Settlement accepted structurally valid but
semantically different page state.

        ↓

INVARIANT
Page advancement must preserve numeric identity
within the current pagination scope.

        ↓

REGRESSION
Given pages 1,2,3,4...
page 2 → page 3 must settle on primary ID 3.

        ↓

LESSON
Never infer pagination advancement from generic
structural similarity when an authoritative numeric
identity exists.
```

Now the lesson is valuable.

It is not merely “knowledge.”

It has:

**failure → cause → invariant → gate → evidence**

That is exactly the sort of knowledge that can survive from one agent run to another.

---

# 5. I would make the DAG a DAG of **proof obligations**

This is probably the most powerful formulation.

Every meaningful engineering change creates a collection of obligations.

For example:

```text
             Feature
                │
       ┌────────┼────────┐
       ▼        ▼        ▼
   Contract   Safety   Compatibility
       │        │        │
       ▼        ▼        ▼
    Tests     Tests    Tests
       │        │        │
       └────────┼────────┘
                ▼
             Review
                │
                ▼
            Deployment
                │
                ▼
            Telemetry
```

Each node asks:

> **What must be proven?**

Not:

> What code should be written?

That's a huge conceptual difference.

---

# 6. This also changes what an “AI skill” should be

The article says:

> agents should pull work from the graph rather than relying on prompt/context memory. 

I'd take that further.

An agent shouldn't receive:

> “Fix pagination.”

It should receive:

```text
TASK NODE

Claim:
Numeric pagination cannot advance to an
incorrect page identity.

Prerequisites:
- pagination state model understood
- primary page identity defined

Required evidence:
- regression test
- invariant test
- integration test

Forbidden shortcuts:
- fallback identity cannot satisfy primary identity
- mutation cannot bypass settlement

Done when:
all gates PASS
```

Now the agent has a **bounded proof obligation**.

The model can choose the implementation.

The graph determines what “done” means.

That's much more powerful than sophisticated prompting.

---

# 7. This is where “you don't review the code, review the gates” becomes interesting

I would actually modify that statement.

Taken literally, it is wrong.

You **sometimes absolutely must review code**.

Security-sensitive code, concurrency primitives, novel algorithms, cryptography, infrastructure changes, etc. can contain defects that tests fail to expose.

So I would phrase it as:

> **Review the gates first; review the code where the gates cannot establish sufficient confidence.**

That is a much stronger engineering principle.

The objective isn't:

> eliminate human code review.

It is:

> **move human attention toward the places where automated evidence is weakest.**

---

# 8. Then what exactly is a review gate?

A gate is a predicate that must evaluate to true.

For example:

```text
G1 Contract gate
      ↓
Does implementation satisfy interface contract?

G2 Invariant gate
      ↓
Does behavior preserve architectural invariant?

G3 Regression gate
      ↓
Does historical failure remain impossible?

G4 Integration gate
      ↓
Does behavior remain valid across boundaries?

G5 Operational gate
      ↓
Does production telemetry remain within expected limits?
```

Now “review” becomes evaluation of the **gate system**.

Ask:

> Is the gate capable of detecting the failure we care about?

That is much deeper than:

> Does the diff look reasonable?

---

# 9. Gates themselves need to be evaluated

This is the part I think is missing from a lot of “AI writes code” discussions.

A green test suite does **not** necessarily mean a safe system.

It might mean your tests are bad.

Therefore the DAG should include **meta-gates**.

For example:

```text
             Change
                │
                ▼
         Required invariant
                │
        ┌───────┴────────┐
        ▼                ▼
   Test exists       Static check
        │                │
        └───────┬────────┘
                ▼
          Evidence passes
                │
                ▼
       ┌───────────────────┐
       │ Is evidence       │
       │ actually capable  │
       │ of falsifying the │
       │ claim?             │
       └─────────┬─────────┘
                 ▼
             Gate valid
```

That last node is incredibly important.

Otherwise agents learn:

> “Make the tests green.”

Instead of:

> “Produce evidence capable of disproving the claim.”

---

# 10. This gives you a hierarchy of confidence

I would model engineering confidence something like:

```text
              Human judgment
                    ▲
                    │
             ┌──────┴──────┐
             │  Gate       │
             │ adequacy    │
             └──────┬──────┘
                    │
             Evidence validity
                    ▲
                    │
             Automated gates
                    ▲
                    │
               Tests/checks
                    ▲
                    │
                 Code
```

Notice what happened.

**Code moves to the bottom of the hierarchy.**

Not because code doesn't matter.

Because code is only one component of the evidence chain.

---

# 11. The really interesting metric becomes “gate strength”

The article makes an important argument:

> autonomy should increase with the strength of the verification mechanism. 

I think this can become an actual engineering metric.

For every task:

```text
             CLAIM
               │
               ▼
        ┌─────────────┐
        │ Verification│
        │ strength    │
        └──────┬──────┘
               │
        ┌──────┴───────────┐
        ▼                  ▼
      Strong             Weak
        │                  │
        ▼                  ▼
   High autonomy       Human review
```

For example:

| Change                  | Gate strength | Agent autonomy |
| ----------------------- | ------------- | -------------- |
| Formatting              | Very high     | Very high      |
| CRUD endpoint           | High          | High           |
| Schema migration        | Medium        | Moderate       |
| distributed concurrency | Low           | Low            |
| security boundary       | Very low      | Very low       |

So the question becomes less:

> “How smart is the model?”

and more:

> **“How much of this work can the system mechanically falsify?”**

That's a much better basis for agent autonomy.

---

# 12. DAGs also provide something surprisingly valuable: **parallelism**

This is another direct lesson from Prove2Me. Multiple agents can work simultaneously when dependencies are explicit. 

Suppose:

```text
                   API contract
                  /     |      \
                 /      |       \
                ▼       ▼        ▼
           implementation  tests  docs
                │          │
                └────┬─────┘
                     ▼
                 integration
                     │
                     ▼
                  release
```

Three agents can work concurrently without inventing their own notion of project state.

That is much more interesting than merely “using a DAG.”

The DAG becomes a **coordination protocol**.

---

# 13. But there is one crucial modification: DAG nodes need state

A node should not simply be:

```text
"Implement pagination"
```

It should be something closer to:

```text
NODE

ID: INV-PAG-07

Claim:
Numeric page identity must survive settlement.

Depends on:
INV-STATE-02

Produces:
implementation
regression test

Verifier:
go test ./...

Evidence:
test/bug106_test.go

Status:
PASS

Confidence:
HIGH

Historical failures:
bug-106
bug-119
```

Now the DAG becomes a living engineering memory.

A future agent can discover:

> “I've seen this failure before.”

That directly addresses the context-loss problem described in the article. 

---

# 14. And this solves the “entangled code” problem surprisingly well

Imagine this giant dependency mess:

```text
              ┌─────────────┐
              │   Browser   │
              └──────┬──────┘
                     │
             ┌───────▼────────┐
             │ Render Engine  │
             └───────┬────────┘
                     │
        ┌────────────┼─────────────┐
        ▼            ▼             ▼
     Parser       State        Pagination
        │            │             │
        └────────────┼─────────────┘
                     ▼
                  Output
```

Instead of trying to simplify this topology, put a second topology over it:

```text
Requirement
    ↓
Invariant
    ↓
Boundary contract
    ↓
Failure mode
    ↓
Regression
    ↓
Integration evidence
    ↓
Release
```

The first graph tells you:

> **how the machine works.**

The second tells you:

> **what you need to prove about the machine.**

That's the reconciliation.

---

# 15. The DAG is therefore an abstraction layer above code

I would describe the architecture as:

```text
                 HUMAN / AGENT REASONING
                         │
                    ┌────▼────┐
                    │   DAG   │
                    │ Claims  │
                    │Evidence │
                    │ Gates   │
                    └────┬────┘
                         │
                   maps onto
                         │
                    ┌────▼────┐
                    │  CODE   │
                    │ SYSTEM  │
                    └────┬────┘
                         │
                  produces behavior
                         │
                    ┌────▼────┐
                    │ REAL    │
                    │ WORLD   │
                    └─────────┘
```

The DAG isn't a replacement for software architecture.

It is the **control plane for reasoning about software architecture and change**.

---

# 16. And lessons should flow *backward* into the DAG

This is where I think your “lessons and insights” idea becomes genuinely interesting.

Production failure:

```text
Incident
   ↓
Root cause
   ↓
Missing invariant
   ↓
New DAG node
   ↓
New gate
   ↓
New regression
```

So engineering knowledge evolves like:

```text
     FAILURE
        │
        ▼
     LESSON
        │
        ▼
    INVARIANT
        │
        ▼
      GATE
        │
        ▼
     EVIDENCE
        │
        ▼
     AUTONOMY
```

The important thing is that **a lesson is not considered institutional knowledge until it has become enforceable**.

That's a very powerful principle.

---

# 17. It also changes the meaning of technical debt

Technical debt could become:

> **known claims for which we do not yet have adequate gates.**

For example:

```text
Known risk
    ↓
Invariant exists
    ↓
No reliable automated evidence
    ↓
Human review required
```

That's much more precise than saying:

> “This code is messy.”

You can actually see where the engineering system is weak.

---

# 18. The grander model

I think the most useful synthesis of the article is not:

> **“DAGs are the future of software engineering.”**

It's:

> **Software is too complex to be represented by a single clean graph. But the reasoning required to safely change software can often be represented as a DAG.**

And then:

> **The DAG becomes the coordination layer for humans and agents, while tests, static analysis, compilers, integration environments, and production telemetry become its verification machinery.**

That leads to an engineering loop:

```text
                 ┌───────────────┐
                 │    CHANGE     │
                 └───────┬───────┘
                         ▼
                 ┌───────────────┐
                 │    CLAIMS     │
                 └───────┬───────┘
                         ▼
                 ┌───────────────┐
                 │      DAG      │
                 │ dependencies  │
                 │ obligations   │
                 └───────┬───────┘
                         ▼
              ┌─────────────────────┐
              │   AGENTS / HUMANS   │
              └──────────┬──────────┘
                         ▼
                       CODE
                         ▼
              ┌─────────────────────┐
              │       GATES         │
              │ compiler            │
              │ tests               │
              │ contracts           │
              │ security            │
              │ integration         │
              └──────────┬──────────┘
                         ▼
                    EVIDENCE
                         ▼
                  ┌────────────┐
                  │  RELEASE   │
                  └─────┬──────┘
                        ▼
                  REAL-WORLD DATA
                        │
                        ▼
                     LESSON
                        │
                        └──────────────► DAG
```

**That is the architecture I would pursue.**

Not **DAG → code**.

But:

**DAG → obligations → agents → code → gates → evidence → lessons → DAG.**

And I think that is also the more defensible interpretation of the article's central lesson: the breakthrough was not that a DAG somehow made a 13-million-line proof simple; it made the **work and state around that complexity explicit enough for many agents to coordinate and verify it**. 

The provocative engineering thesis that falls out of this is:

> **The future of agentic software engineering may not be “AI that writes code.” It may be systems that turn software changes into explicit proof obligations, then let humans and agents compete to satisfy those obligations while mechanical gates decide whether the work is actually true.**
### Answer
Solvent already implements a close cousin of the architecture the article describes, and it does so as working, invariant-enforced code rather than a philosophical framework — which is useful because it tells you exactly which parts of that vision are already tractable engineering and which parts are still open design problems.

## What Solvent already is

Solvent is a "transactional belief ledger for autonomous agents" built on four tables — `belief`, `belief_edge`, `evidence`, `action_intent` — where retrieval only *proposes* candidates and the database is the sole authority that *decides* solvent:111-114 . That is precisely the "DAG models obligations and evidence, not the software itself" principle from the article, already reified as SQL rows.

**The obligation/claim DAG.** A `belief` enters at `status='entered'` carrying a fixed set of "debt" items — epistemic obligations that must be retired before the belief can be promoted:

```
needMap, needInvariant, needToyCheck, needNullModel, needObstruction, needFaithfulnessReview
``` solvent:335-347 

This is exactly the article's "obligation DAG" (requirement → invariant → behavior → acceptance condition) — generic enough to be domain-agnostic, but concrete enough to be mechanically checked.

**The evidence graph.** `evidence` rows attach to beliefs with a typed `provenance_class` (`external_feed`, `reproducible_artifact`, `live_scan`, `operator_asserted`) and a mandatory `content_sha256` solvent:53-55 . Debt is discharged either from a retrieval citation the agent actually selected, or from an operator-asserted artifact — never from unverified assertion solvent:37-50 .

**The dependency/knowledge graph.** `belief_edge` carries `derives` and `contradicts` edges between beliefs, and is explicitly *not necessarily acyclic* — the recursive traversal uses `UNION` rather than `UNION ALL` specifically because the schema forbids only self-loops, not longer cycles, so deduplication (not acyclicity) is what guarantees termination solvent:429-433 solvent:132-135 . This matches the article's observation that the "knowledge graph" layer is not guaranteed to be a DAG even when the "obligation graph" mostly is.

**Gates, not code review, as the review surface.** Three CockroachDB-enforced invariants carry the entire safety guarantee:

```sql
-- promoted_is_debt_free
CHECK (status <> 'promoted' OR debt = '{}')
-- gate
FOREIGN KEY (belief_id, belief_status) REFERENCES belief(id, status) ON UPDATE CASCADE
-- live_requires_promoted
CHECK (state <> 'live' OR belief_status = 'promoted')
``` solvent:122-132 

These are the literal instantiation of "you don't review the code, you review the gates" — the agent's own judgment about whether a belief is promotable is structurally irrelevant; the schema decides solvent:191-197 .

**The authority boundary the article calls for.** The explicit division of responsibility is: agent owns reasoning, MCP owns translation, kernel owns transaction discipline, CockroachDB owns invariants solvent:174-199 . The MCP server itself is deliberately thin — every tool handler is "unmarshal → kernel call → format," with no belief semantics or SQL policy logic in that layer solvent:1-4 . The six tools it exposes (`solvent_ledger`, `solvent_ingest_evidence`, `solvent_retire_debt`, `solvent_promote`, `solvent_authorize_action`, `solvent_falsify`) are the entire agent-facing surface solvent:88-223 .

## How this maps onto the article's proposal — and what's missing for a domain-agnostic oracle harness

The pieces that generalize immediately:

1. **Obligation DAG → domain-agnostic debt vocabulary.** `kernel.FullDebt` is currently six fixed strings, and there is already a `db/004_debt_vocabulary.sql` migration suggesting the vocabulary is meant to be swappable per scenario solvent:28-30 . A software-engineering harness would define its own debt vocabulary — e.g. `needRegressionTest`, `needInvariantProof`, `needSecurityReview`, `needIntegrationEvidence` — mirroring the article's "proof obligations" list, without touching the kernel.

2. **Action intents → the article's "release gate."** `solvent_authorize_action` inserts a `live` intent that cites a promoted belief as its warrant, and the composite FK refuses it (`23503 · gate`) unless that belief is promoted *right now* solvent:181-203 . This is a direct implementation of "agents pull bounded proof obligations, gates decide when work is done" — merges/deploys become `action_intent` rows gated on invariant-checked beliefs.

3. **Falsification → lessons flowing backward into the DAG.** `solvent_falsify` retracts a belief and cancels its dependent live intent in the same transaction solvent:205-223 , and `RetractCascade` generalizes this transitively across `belief_edge`, cancelling before retracting in one transaction (I-8) solvent:108-108 . This is the mechanical version of "incident → root cause → invariant → new gate," except Solvent only models the *retraction* half — a production incident becomes a falsifying evidence row that kills a belief and everything downstream, but there's no built-in step that automatically *mints a new debt item* from the incident. That's the gap: the article's "failure → lesson → invariant → gate" loop requires a harness-level convention (e.g. an MCP tool `solvent_propose_debt_item` or a review-gated `belief_edge` from an incident-belief to a "new debt vocabulary entry" belief) that Solvent doesn't yet have.

4. **The audit query as "review the gates."** `AuditLiveOnNonPromoted` is a standing invariant check — it must read zero in any committed state solvent:109-109 . In a general oracle-verifier harness, this is the thing a human or meta-agent actually reviews: not the diff, but whether the audit invariant holds and whether the debt vocabulary for this domain is *strong enough* to falsify the failure modes that matter (the article's "gate adequacy" layer, section 9 above). Solvent enforces gate *presence* mechanically; it does not yet mechanically evaluate gate *adequacy* — that remains a human/meta-agent judgment recorded as its own belief with its own debt.

5. **The measured-retrieval finding is the empirical argument for "review gates, not code."** Solvent measured that the canonical falsifying evidence for the etcd deployment claim ranks 573rd of 7,239 by embedding distance solvent:104-106  — i.e., an LLM/agent's semantic judgment about relevance is not a safe substitute for a schema-enforced gate. That is the concrete, non-hypothetical version of the article's claim that "a green test suite does not mean a safe system" — here it's "a plausible top-k retrieval does not mean a justified belief."

**What a domain-agnostic oracle-verifier harness would need beyond current Solvent:**
- A per-domain debt-vocabulary registration mechanism (swap-in obligation sets, not hardcoded like `kernel.FullDebt`).
- A "gate adequacy" belief type — a belief *about* a gate, so meta-review is itself first-class and falsifiable, not out-of-band commentary.
- An incident→new-obligation pipeline connecting `solvent_falsify` outcomes to debt-vocabulary evolution, closing the article's lesson-loop instead of stopping at retraction.
- Multiple scenario-scoped belief graphs coordinating concurrent agents the way the article's Prove2Me example does — Solvent already scopes `RetractCascade` and traversal by `scenario_id`, which is the primitive this needs, just not yet exposed as a multi-agent coordination protocol over MCP.
### Citations
**File:** README.md (L104-106)
```markdown
The canonical incident is worse. `#13766 · Inconsistent revision and data occurs`, the issue that
actually destroys the claim, sits at distance `0.594920` from it — a measured **rank 573 of 7,239**.
No top-k, threshold or reranker puts that in front of a judge.
```
**File:** README.md (L111-114)
```markdown
The response is not to tune the retrieval. It is to stop treating retrieval as an authority:

> Retrieval proposes candidates. It decides nothing. The database decides.

```
**File:** README.md (L122-132)
```markdown
```sql
-- 1. Debt must be empty before a belief can be promoted.
CHECK (status <> 'promoted' OR debt = '{}')          -- promoted_is_debt_free -> 23514

-- 2. An intent may only cite a belief that is promoted RIGHT NOW.
FOREIGN KEY (belief_id, belief_status)
  REFERENCES belief(id, status) ON UPDATE CASCADE     -- gate -> 23503

-- 3. And a live intent may not survive its belief losing that status.
CHECK (state <> 'live' OR belief_status = 'promoted') -- live_requires_promoted -> 23514
```
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L335-347)
```markdown
### 6.3 Debt Items

The six starting debt items represent epistemic obligations that must be retired
before promotion:

| Item | Meaning | Retired By |
|---|---|---|
| `needMap` | The vulnerability-to-version mapping is known | CVE advisory with affected/fixed versions |
| `needInvariant` | The behavioral invariant is verified | Maintainer confirmation of fix behavior |
| `needToyCheck` | A minimal reproduction exists | Test case or reproducer confirmed |
| `needNullModel` | The baseline (no-fix) behavior is characterized | Understanding of what happens without the fix |
| `needObstruction` | The obstruction path is identified | Root cause analysis confirmed |
| `needFaithfulnessReview` | Security review is complete | Security team sign-off or maintainer assertion |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L429-433)
```markdown
### 8.3 Traversal

The dependency graph is traversed via `belief_edge` edges. The cascade
(`RetractCascade`) follows `derives` edges transitively using a recursive CTE.
Traversal is scoped to `scenario_id` — edges leaving the scenario are not followed.
```
**File:** plans/pre/IMPLEMENTATION_CONTRACT.md (L53-55)
```markdown
- `evidence.provenance_class` ∈ {`external_feed`,`reproducible_artifact`,`live_scan`,
  `operator_asserted`}; `content_sha256` NOT NULL.
- `belief.embedding` is nullable; the ledger is fully functional with every embedding NULL.
```
**File:** internal/wizard/discharge.go (L37-50)
```go
// Discharge retires one debt item, recording the artifact that justifies it.
//
// Two kinds of check, and the visible difference between them is a feature:
//
//   - The retrieval-backed pair can only be discharged by something the judge actually
//     retrieved. There is no text box; there is a citation or there is a refusal.
//   - The other four accept a typed artifact and store it as operator-asserted evidence,
//     which is what human judgement looks like on the record.
//
// Idempotent. evidence carries no uniqueness constraint — contrary to two planning
// documents that claim UNIQUE(belief_id, content_sha256) — so re-discharging an
// already-retired item would otherwise silently add a second row. The outstanding-debt
// guard, not the schema, is what prevents that, which is why it is checked first and
// why a test pins it.
```
**File:** plans/M1/plan1.md (L108-108)
```markdown
| `RetractCascade(ctx, scenarioID, rootID string) (retracted int, err error)` | **Scenario-scoped (D-032).** In ONE transaction: compute the transitive descendants of `rootID` within `scenarioID`; set all `live` intents on those beliefs to `cancelled`; then set those beliefs from `promoted` to `retracted`. Cancel strictly precedes retract (I-8). Returns the count of beliefs retracted. |
```
**File:** plans/M1/plan1.md (L109-109)
```markdown
| `AuditLiveOnNonPromoted(ctx, scenarioID string) (int, error)` | Invariant I-5 as a query, scoped to the scenario. Read path — may use `QueryRowContext` directly; I-7 governs writes. |
```
**File:** plans/M1/plan1.md (L132-135)
```markdown
`UNION`, never `UNION ALL`. The docs warn that a recursive subquery must eventually return no
rows; `belief_edge`'s CHECK forbids only self-loops (`parent_id <> child_id`), not longer
cycles, so **`UNION`'s deduplication is what guarantees termination** on a cyclic belief graph.
This is load-bearing, not stylistic.
```
**File:** SOLVENT_ENGINEERING_GUIDE.md (L174-199)
```markdown
The MCP layer and Go application are intentionally not the final authority.

The design boundary is:

```text
Agent
  │
  ▼
MCP server
  │
  ▼
Solvent kernel / transaction code
  │
  ▼
CockroachDB
```

The agent does not get to decide that a belief is promotable or an intent is safe.

The MCP server is deliberately thin: validate input, call one kernel function, format the result. It does not implement belief semantics or SQL policy checks.

The kernel owns transaction discipline and reports the database result.

CockroachDB is the final invariant boundary.

The repository explicitly defines this division of responsibility: the agent owns reasoning, the MCP layer owns translation, the kernel owns transaction discipline, and CockroachDB owns the invariants. fileciteturn219file8L592-L623
```
**File:** cmd/solvent-mcp/main.go (L1-4)
```go
// Command solvent-mcp is a stdio MCP server exposing the Solvent transactional
// belief ledger as six tools. The server is an adapter — it has no opinion about
// beliefs. Every tool handler is exactly three moves: unmarshal → kernel call →
// format.
```
**File:** cmd/solvent-mcp/main.go (L88-223)
```go
	server.AddTool(&mcp.Tool{
		Name:        "solvent_ledger",
		Description: "Read the current ledger for a scenario: beliefs with status and open debt, optionally their evidence, action intents with state, and the safety audit count. This is the only source of truth about current state. Call it before asserting any count, status, or identifier, and call it again after any mutation — never answer from memory of an earlier tool result, and never state a number you did not just read here.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario to query",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "Optional: filter to a single belief by UUID",
				},
				"include_evidence": map[string]any{
					"type":        "boolean",
					"description": "Include evidence rows (default false)",
				},
			},
			"required": []string{"scenario"},
		},
	}, toolHandler("solvent_ledger"))

	server.AddTool(&mcp.Tool{
		Name:        "solvent_ingest_evidence",
		Description: "Process the pinned evidence fixtures for a scenario through the full pipeline (normalize → derive → ledger). Idempotent: re-running creates no duplicate beliefs, evidence, or intents.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario to ingest evidence for",
				},
			},
			"required": []string{"scenario"},
		},
	}, toolHandler("solvent_ingest_evidence"))

	server.AddTool(&mcp.Tool{
		Name: "solvent_retire_debt",
		// The valid items are no longer restated in prose. They were, and the prose was
		// one of five hand-copies of kernel.FullDebt; the Phase 5 vocabulary rename had
		// to find every one of them. The enum below is generated from the kernel, so
		// this description can only describe behaviour, not enumerate values.
		Description: "Record that one review obligation on a belief has been discharged. debt_item must be one of the six items the database issued (see the enum). An unrecognised item is refused rather than silently ignored.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief",
				},
				"debt_item": map[string]any{
					"type": "string",
					// Generated, not transcribed. Advertising the vocabulary to the agent
					// is only half the job: this SDK's low-level AddTool does not validate
					// arguments against the schema, so the enum is documentation and
					// handleSolventRetireDebt does the refusing.
					"enum":        kernel.FullDebt,
					"description": "Debt item to retire. Must be one of the six the database issued.",
				},
			},
			"required": []string{"scenario", "belief_id", "debt_item"},
		},
	}, toolHandler("solvent_retire_debt"))

	server.AddTool(&mcp.Tool{
		Name:        "solvent_promote",
		Description: "Attempt to promote a belief to authorized status. The database refuses promotion while the belief carries any open debt item, returning constraint promoted_is_debt_free (SQLSTATE 23514). Call this whenever the user asks to promote a belief — do not pre-check the debt yourself and do not decline on the user's behalf. The database is the authority on whether promotion is permitted; your job is to attempt it and report what it says.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief to promote",
				},
			},
			"required": []string{"scenario", "belief_id"},
		},
	}, toolHandler("solvent_promote"))

	server.AddTool(&mcp.Tool{
		Name:        "solvent_authorize_action",
		Description: "Record a live intent to take a real-world action, citing a belief as its warrant. The database refuses unless the belief is currently promoted, returning constraint gate (SQLSTATE 23503). Call this when the user asks to authorize, deploy, or act on a belief. Do not pre-check the belief's status.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief to cite as warrant",
				},
				"action": map[string]any{
					"type":        "string",
					"description": "Description of the real-world action to authorize",
				},
			},
			"required": []string{"scenario", "belief_id", "action"},
		},
	}, toolHandler("solvent_authorize_action"))

	server.AddTool(&mcp.Tool{
		Name:        "solvent_falsify",
		Description: "Retract a belief that new evidence has falsified. Cancels that belief's dependent live intent in the same transaction. Retracts a single belief — this does not propagate across a belief graph. Obtain the belief's id from solvent_ledger immediately before calling.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief to retract (read from solvent_ledger)",
				},
			},
			"required": []string{"scenario", "belief_id"},
		},
	}, toolHandler("solvent_falsify"))
```
**File:** Taskfile.yml (L28-30)
```yaml
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/002_corpus.sql
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/003_wizard.sql
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/004_debt_vocabulary.sql
```
## Q2
now the question is, do we need DAG here since we already have Solvent, does DAG make it redundant or not necessary?
### Answer
Direct answer: **no, DAG doesn't make Solvent redundant, but it's also not what Solvent's `belief_edge` graph actually is** — and the distinction matters for your harness design.

## Solvent's graph is explicitly *not* a DAG

The one place the phrase "directed acyclic graph" appears in the repo is a planning document describing an *idealized narrative example* solvent:815-816 . The actual schema and kernel go out of their way to reject that assumption:

- `belief_edge`'s `CHECK` constraint forbids only self-loops (`parent_id <> child_id`), not longer cycles solvent:132-135 .
- `RetractCascade`'s recursive CTE uses `UNION` (not `UNION ALL`) specifically because deduplication, not acyclicity, is what guarantees termination on a cyclic belief graph — this is flagged as "load-bearing, not stylistic" solvent:111-135 .
- The design review prompt explicitly interrogates cycle-handling as an open risk: "the cascade follows `derives` edges transitively... What if `contradicts` edges form a cycle with `derives` edges?" solvent:54-56 , and there's a dedicated test `TestCyclicGraph` asserting termination on a 2-node cycle solvent:346-353 .

So Solvent's own engineering already embodies the article's point #2/#3: real dependency structures aren't clean DAGs, and the system is built to survive that rather than assume it away.

## What this means for "do we need a DAG layer"

**You don't need a separate DAG data structure bolted on top of Solvent — Solvent's belief graph already is the obligation/evidence graph the article calls for.** The article's four-graph decomposition maps onto Solvent's existing tables almost one-to-one:

| Article's layer | Solvent's implementation |
|---|---|
| Obligation graph | `belief.debt` array — the six starting items solvent:335-347  |
| Evidence graph | `evidence` table, `Supports` edges (`evidence.belief_id` FK) solvent:420-427  |
| Dependency/knowledge graph | `belief_edge` with `derives`/`contradicts` kinds — general graph, not DAG solvent:420-427  |
| Execution/release gate | `action_intent` gated by the composite FK solvent:122-132  |

Adding a "clean DAG" as a separate primitive on top would be **redundant at best, and actively wrong at worst** — it would reintroduce exactly the acyclicity assumption that Solvent's own schema and tests were built to *not* require. The article's own conclusion — "use a DAG for execution/verification, use a graph for knowledge" — is already satisfied by treating `belief_edge` as a general (possibly cyclic) graph, with the *acyclicity* only holding locally and provisionally wherever `derives` edges happen not to loop back, never assumed globally.

## Where a DAG-shaped abstraction is still useful — as a *view*, not a new table

There is one place where DAG-thinking earns its keep: **as a derived, read-time projection**, not a stored structure. Since `RetractCascade`'s traversal is scoped per `scenario_id` and computed via recursive CTE at call time solvent:429-433 , a harness UI or meta-agent could present the *currently reachable, currently non-cyclic subset* of a scenario's belief graph as a DAG for human/agent legibility — while the underlying `belief_edge` table stays a general graph that CockroachDB never assumes is acyclic. That's the article's "DAG is a mental model for humans/agents, not a claim about the underlying structure" principle, and Solvent already gives you the primitive (`GetDescendants`/`RetractCascade` traversal) to compute that view on demand solvent:254-263 .

**Bottom line:** Solvent doesn't need a DAG bolted on; it needs the DAG *removed from your mental model of it*. The engineering discipline that matters — debt must be empty before promotion, gate FK enforcement, cancel-before-retract in the cascade — is already invariant-enforced independent of whether the graph is acyclic. Building a harness on Solvent as-is (general graph + gates) is more correct than building one on an assumed-DAG abstraction, because the latter would break the moment two beliefs mutually derive each other or a `contradicts` edge closes a loop with a `derives` edge — a case the schema and tests already anticipate but a naive DAG library would not.
### Citations
**File:** plans/post/evidence_dataset.md (L815-816)
```markdown
The resulting belief topology forms a directed acyclic graph^1. Security evidence
(ev_ghsa_q8m4_xhhv_38mg) supports the derived vulnerability claim bel_v3527_vulnerable via
```
**File:** plans/M1/plan1.md (L111-135)
```markdown
## `RetractCascade` — shape (D-032 + D-033)

`belief_edge` carries **no `scenario_id` column**, so scoping the traversal means constraining
each hop through `belief`. Both statements carry the same CTE:

```sql
WITH RECURSIVE d(id) AS (
    SELECT b.id FROM belief b WHERE b.id = $1::UUID AND b.scenario_id = $2::UUID
  UNION
    SELECT e.child_id
      FROM belief_edge e
      JOIN d          ON e.parent_id = d.id
      JOIN belief cb  ON cb.id = e.child_id AND cb.scenario_id = $2::UUID
)
UPDATE action_intent SET state = 'cancelled'
 WHERE state = 'live' AND scenario_id = $2::UUID AND belief_id IN (SELECT id FROM d);
-- then, same CTE:
UPDATE belief SET status = 'retracted'
 WHERE status = 'promoted' AND scenario_id = $2::UUID AND id IN (SELECT id FROM d);
```

`UNION`, never `UNION ALL`. The docs warn that a recursive subquery must eventually return no
rows; `belief_edge`'s CHECK forbids only self-loops (`parent_id <> child_id`), not longer
cycles, so **`UNION`'s deduplication is what guarantees termination** on a cyclic belief graph.
This is load-bearing, not stylistic.
```
**File:** plans/PRD/DATA_PIPELINE_SPEC_REVIEW_PROMPT.md (L54-56)
```markdown
5. **Termination:** The cascade follows `derives` edges transitively. Does
   the spec define cycle handling? What if `contradicts` edges form a cycle
   with `derives` edges?
```
**File:** plans/PRD/plan2.md (L254-263)
```markdown
### 4.8 `internal/graph/graph.go` (~150 lines)

Belief edge management.

Functions:
- `AddEdge(ctx context.Context, db *sql.DB, parentID, childID, kind string) error`
- `GetDescendants(ctx context.Context, db *sql.DB, scenarioID, rootID string) ([]string, error)`
- `GetChildren(ctx context.Context, db *sql.DB, beliefID string) ([]string, error)`

Edge kinds: `derives`, `contradicts` (from schema CHECK).
```
**File:** plans/PRD/plan2.md (L346-353)
```markdown
### 5.6 `internal/graph/graph_test.go`

| Test | Asserts |
|---|---|
| `TestEdgeCreation` | Creates derives edge, queryable |
| `TestEdgeContradicts` | Creates contradicts edge |
| `TestGetDescendants` | 3-node chain → returns all descendants |
| `TestCyclicGraph` | 2-node cycle → terminates, returns 2 |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L335-347)
```markdown
### 6.3 Debt Items

The six starting debt items represent epistemic obligations that must be retired
before promotion:

| Item | Meaning | Retired By |
|---|---|---|
| `needMap` | The vulnerability-to-version mapping is known | CVE advisory with affected/fixed versions |
| `needInvariant` | The behavioral invariant is verified | Maintainer confirmation of fix behavior |
| `needToyCheck` | A minimal reproduction exists | Test case or reproducer confirmed |
| `needNullModel` | The baseline (no-fix) behavior is characterized | Understanding of what happens without the fix |
| `needObstruction` | The obstruction path is identified | Root cause analysis confirmed |
| `needFaithfulnessReview` | Security review is complete | Security team sign-off or maintainer assertion |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L420-427)
```markdown
### 8.2 Edges

| Edge Type | Table | Meaning |
|---|---|---|
| **Supports** | `evidence.belief_id` (FK) | This evidence supports this belief. |
| **Derives** | `belief_edge` (`kind='derives'`) | This belief depends on that belief. |
| **Contradicts** | `belief_edge` (`kind='contradicts'`) | This belief contradicts that belief. |
| **Governs** | `action_intent.belief_id` (FK) | This intent is governed by this belief. |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L429-433)
```markdown
### 8.3 Traversal

The dependency graph is traversed via `belief_edge` edges. The cascade
(`RetractCascade`) follows `derives` edges transitively using a recursive CTE.
Traversal is scoped to `scenario_id` — edges leaving the scenario are not followed.
```
**File:** README.md (L122-132)
```markdown
```sql
-- 1. Debt must be empty before a belief can be promoted.
CHECK (status <> 'promoted' OR debt = '{}')          -- promoted_is_debt_free -> 23514

-- 2. An intent may only cite a belief that is promoted RIGHT NOW.
FOREIGN KEY (belief_id, belief_status)
  REFERENCES belief(id, status) ON UPDATE CASCADE     -- gate -> 23503

-- 3. And a live intent may not survive its belief losing that status.
CHECK (state <> 'live' OR belief_status = 'promoted') -- live_requires_promoted -> 23514
```
```
## Q3
if you are to design a ticket management system like Jira, how would you do it on top of Solvent primitives? itemize your reply organized around categories
### Answer
Mapping Jira onto Solvent's four tables is a reasonably clean fit for the "core truth" layer, but Jira has several concerns (free-form workflow states, sprints, roles/permissions, notifications) that sit outside what `belief`, `belief_edge`, `evidence`, and `action_intent` model. Below is the mapping organized by category, followed by where the fit breaks and what you'd need to add.

## 1. Core entity mapping

| Jira concept | Solvent primitive | Notes |
|---|---|---|
| Ticket/Issue | `belief` row | The ticket's "claim" is the ticket's title/description — a claim asserted (`derived`/`accommodated`/`postulated`) rather than fixed fact solvent:95-95 . |
| Project | `scenario_id` | Everything in Solvent is scenario-scoped by construction; a Jira project maps naturally onto a scenario boundary, since traversal, cascades, and audits are all scoped there solvent:111-113 . |
| Epic → Story → Subtask | `belief_edge` with `kind='derives'` | A subtask "derives from" (depends on) its parent story; an epic is the ancestor belief whose descendants are its child tickets solvent:420-427 . |
| Blocks / Blocked-by | `belief_edge` with `kind='contradicts'`, or a second `derives` edge in the reverse direction | The schema already supports a general, possibly-cyclic graph — `belief_edge`'s CHECK only forbids self-loops, not longer cycles solvent:132-135  — so "A blocks B" and "B blocks A" (a real, if pathological, Jira state) doesn't break the schema the way it would break a naive DAG library. |
| Attachment / Comment / CI result | `evidence` row | Typed by `provenance_class` ∈ `{external_feed, reproducible_artifact, live_scan, operator_asserted}` with mandatory `content_sha256` solvent:53-55  — e.g. a CI run maps to `reproducible_artifact`, a manager's manual sign-off to `operator_asserted`. |
| Sprint/release action, "deploy," "close and notify" | `action_intent` | A live intent citing the ticket-belief as its warrant, gated by the composite FK so it can only exist while the ticket is in a "promoted" (done/approved) state solvent:122-132 . |

## 2. Ticket lifecycle (`entered → promoted → retracted`)

Jira's default states (`To Do → In Progress → Done`, plus `Reopened`) map onto Solvent's three-state belief lifecycle, but semantically shifted:

- `entered`: ticket filed, carrying full "debt" — i.e., a fixed set of Definition-of-Done obligations not yet satisfied. Never enters pre-satisfied — "a belief is never entered at promoted... Promotion is earned, not assumed" solvent:349-357 .
- `promoted`: all DoD debt retired, so the schema itself accepts the "Done" transition — `CHECK (status <> 'promoted' OR debt = '{}')` solvent:122-123 .
- `retracted`: reopened/invalidated (a regression is found, a requirement changed). This is where Jira's fuzzy "Reopened" state becomes precise: reopening a ticket that had live downstream actions (a deploy intent, a linked release) cascades.

## 3. "Definition of Done" as debt vocabulary

Instead of the security-specific `needMap`/`needInvariant`/etc. solvent:335-347 , a ticket-management debt vocabulary would be domain-specific items such as `needAcceptanceCriteria`, `needCodeReview`, `needQASignoff`, `needCIGreen`, `needDocsUpdated`. Each is retired only by an artifact of the matching kind: `RetireDebt` removes one item via `array_remove`, is idempotent when the item is already absent, and never changes ticket status by itself solvent:105-105 . This directly gives you Jira's "can't transition to Done until checklist items are satisfied" — except enforced by a CHECK constraint the agent cannot talk its way around, rather than a UI-level workflow rule.

## 4. Ticket relationships and cascades

- `RetractCascade` becomes "invalidate a ticket and everything that depends on it": in one transaction it cancels all `live` intents on the ticket and its transitive descendants, then flips those tickets from `promoted` back to `retracted`, with cancel strictly preceding retract (I-8) solvent:108-108 . For Jira this is the mechanical version of "reopening an epic should reopen/flag its dependent stories and pull back any release action that already fired."
- The traversal is a `WITH RECURSIVE` CTE using `UNION` (not `UNION ALL`) specifically so it terminates even if two tickets end up mutually blocking each other — a state Jira absolutely allows and that a naive DAG assumption would not survive solvent:116-135 .
- Cross-project isolation is structural: the cascade only follows edges within the same `scenario_id`, so a dependency edge accidentally drawn to another project's ticket can't retract tickets there solvent:108-108 .

## 5. Actions/workflow-triggered side effects

`action_intent` state machine (`live → cancelled | executed`) solvent:363-369  maps onto Jira automations: "when this ticket is marked Done, trigger deploy" or "page on-call." The gate is the same composite FK + CHECK pair used everywhere else — a live intent cannot exist unless its citing ticket is currently `promoted` solvent:388-397 . This is stronger than Jira's actual automation model, where a webhook can fire based on a state that's already stale by the time it executes; here the FK is checked at insert time against current status.

## 6. Agent-facing surface (MCP)

The same six-tool pattern generalizes directly: `solvent_ledger` (read current board state — "call it before asserting any count or status... never answer from memory" solvent:88-88 ), `solvent_retire_debt` (check off a DoD item), `solvent_promote` (attempt to close — refused with `23514`/`promoted_is_debt_free` if checklist items remain solvent:167-167 ), `solvent_authorize_action` (trigger a workflow action, refused with `23503`/`gate` if the ticket isn't actually done solvent:181-203 ), and `solvent_falsify`/`RetractCascade` (reopen). An AI PM-bot built on this MCP surface literally cannot mark a ticket "Done" or fire a release action against the agent's own unverified judgment — the database decides solvent:191-197 .

## 7. Where the fit breaks — what Jira needs that Solvent doesn't have

- **Free-text mutable fields.** Jira tickets get retitled, redescribed, reassigned, reprioritized constantly. `belief.claim` is treated as a fixed assertion being evidenced toward truth, not a mutable record — repeatedly editing the claim string undermines the "evidence supports this specific claim" semantics that make `evidence.belief_id` meaningful.
- **No user/role/permission model.** Solvent's boundary is agent → MCP → kernel → CockroachDB, with no notion of assignee, reporter, or team visibility solvent:174-188 . This would need an entirely separate layer.
- **No sprints/boards/time-boxing.** `scenario_id` gives project-level isolation but nothing analogous to a sprint boundary or a Kanban swimlane; you'd need an additional non-kernel table or a convention for scenario sub-scoping.
- **Debt vocabulary is currently fixed, not per-project.** `kernel.FullDebt` is six hardcoded strings solvent:81-84 ; a real ticket system needs per-project-configurable Definition-of-Done, which the kernel doesn't yet expose (the `db/004_debt_vocabulary.sql` migration suggests this is anticipated but not the current shipped mechanism) solvent:28-30 .
- **No notification/event bus.** Jira's comment threads, watchers, and email digests have no analog here — `evidence` rows are append-only provenance, not a notification stream.

**Bottom line:** the ticket *state machine and its enforcement* (can't close without DoD, reopening cascades, actions gated on "actually done") maps almost directly onto Solvent's existing invariants and would be a strict improvement over Jira's UI-level workflow rules, which an agent (or a careless human) can bypass. What you'd be building on top is everything Jira has that isn't about state integrity — identity, scheduling, free-text mutability, notification — none of which Solvent's kernel currently models or should model, per its own domain-independence boundary solvent:503-510 .
### Citations
**File:** plans/M1/plan0.md (L95-95)
```markdown
| `EnterBelief(ctx, scenarioID, claim string, ct ClaimType) (id string, err error)` | Inserts at the door: `status='entered'`, `debt` = the full starting debt, `final_truth` left at its default. Never gated. Returns the new id via `RETURNING id`. |
```
**File:** plans/M1/plan1.md (L81-84)
```markdown
var FullDebt = []string{
    "needMap", "needInvariant", "needToyCheck",
    "needNullModel", "needObstruction", "needFaithfulnessReview",
}
```
**File:** plans/M1/plan1.md (L105-105)
```markdown
| `RetireDebt(ctx, beliefID, item string) error` | Removes one debt item via `array_remove`. Idempotent when absent — zero rows changed is success. |
```
**File:** plans/M1/plan1.md (L108-108)
```markdown
| `RetractCascade(ctx, scenarioID, rootID string) (retracted int, err error)` | **Scenario-scoped (D-032).** In ONE transaction: compute the transitive descendants of `rootID` within `scenarioID`; set all `live` intents on those beliefs to `cancelled`; then set those beliefs from `promoted` to `retracted`. Cancel strictly precedes retract (I-8). Returns the count of beliefs retracted. |
```
**File:** plans/M1/plan1.md (L111-113)
```markdown
## `RetractCascade` — shape (D-032 + D-033)

`belief_edge` carries **no `scenario_id` column**, so scoping the traversal means constraining
```
**File:** plans/M1/plan1.md (L116-135)
```markdown
```sql
WITH RECURSIVE d(id) AS (
    SELECT b.id FROM belief b WHERE b.id = $1::UUID AND b.scenario_id = $2::UUID
  UNION
    SELECT e.child_id
      FROM belief_edge e
      JOIN d          ON e.parent_id = d.id
      JOIN belief cb  ON cb.id = e.child_id AND cb.scenario_id = $2::UUID
)
UPDATE action_intent SET state = 'cancelled'
 WHERE state = 'live' AND scenario_id = $2::UUID AND belief_id IN (SELECT id FROM d);
-- then, same CTE:
UPDATE belief SET status = 'retracted'
 WHERE status = 'promoted' AND scenario_id = $2::UUID AND id IN (SELECT id FROM d);
```

`UNION`, never `UNION ALL`. The docs warn that a recursive subquery must eventually return no
rows; `belief_edge`'s CHECK forbids only self-loops (`parent_id <> child_id`), not longer
cycles, so **`UNION`'s deduplication is what guarantees termination** on a cyclic belief graph.
This is load-bearing, not stylistic.
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L335-347)
```markdown
### 6.3 Debt Items

The six starting debt items represent epistemic obligations that must be retired
before promotion:

| Item | Meaning | Retired By |
|---|---|---|
| `needMap` | The vulnerability-to-version mapping is known | CVE advisory with affected/fixed versions |
| `needInvariant` | The behavioral invariant is verified | Maintainer confirmation of fix behavior |
| `needToyCheck` | A minimal reproduction exists | Test case or reproducer confirmed |
| `needNullModel` | The baseline (no-fix) behavior is characterized | Understanding of what happens without the fix |
| `needObstruction` | The obstruction path is identified | Root cause analysis confirmed |
| `needFaithfulnessReview` | Security review is complete | Security team sign-off or maintainer assertion |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L349-357)
```markdown
### 6.4 Entry Conditions

A belief enters the system when:
1. New evidence arrives that makes a claim relevant.
2. No existing belief in the same scenario covers the same claim (dedup check).
3. The claim is typed as `derived` (from evidence), `accommodated` (merged from existing), or `postulated` (hypothesized).

A belief is **never** entered at `promoted` — all beliefs enter at `entered` with full
debt. Promotion is earned, not assumed.
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L363-369)
```markdown
### 7.1 States

| State | Meaning | Allowed Transitions |
|---|---|---|
| `live` | Agent has proposed an action on a promoted belief. | → `cancelled` (belief retracted), → `executed` (action completed). |
| `cancelled` | The belief this intent rode on has been retracted. Terminal. | None. |
| `executed` | The action was carried out. Terminal (MVP does not model post-execution). | None. |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L388-397)
```markdown
### 7.3 The Gate

A `live` intent **cannot** exist unless its belief is currently `promoted`. This is
enforced by:
1. The composite FK: `(belief_id, belief_status) REFERENCES belief(id, status) ON UPDATE CASCADE`.
2. The CHECK constraint: `live_requires_promoted: state <> 'live' OR belief_status = 'promoted'`.

When a belief is retracted, `ON UPDATE CASCADE` propagates the new status into the
intent row, which then violates `live_requires_promoted` unless the intent was already
cancelled. This is I-4, proven on CockroachDB at M0.
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L420-427)
```markdown
### 8.2 Edges

| Edge Type | Table | Meaning |
|---|---|---|
| **Supports** | `evidence.belief_id` (FK) | This evidence supports this belief. |
| **Derives** | `belief_edge` (`kind='derives'`) | This belief depends on that belief. |
| **Contradicts** | `belief_edge` (`kind='contradicts'`) | This belief contradicts that belief. |
| **Governs** | `action_intent.belief_id` (FK) | This intent is governed by this belief. |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L503-510)
```markdown
### Kernel-Generic (internal/kernel)

These concepts belong to the transactional belief ledger and are domain-agnostic:

- Belief, evidence, action_intent (the three tables + belief_edge)
- Belief lifecycle: entered → promoted → retracted
- Intent lifecycle: live → cancelled/executed
- Debt retirement and promotion gating
```
**File:** plans/pre/IMPLEMENTATION_CONTRACT.md (L53-55)
```markdown
- `evidence.provenance_class` ∈ {`external_feed`,`reproducible_artifact`,`live_scan`,
  `operator_asserted`}; `content_sha256` NOT NULL.
- `belief.embedding` is nullable; the ledger is fully functional with every embedding NULL.
```
**File:** README.md (L122-132)
```markdown
```sql
-- 1. Debt must be empty before a belief can be promoted.
CHECK (status <> 'promoted' OR debt = '{}')          -- promoted_is_debt_free -> 23514

-- 2. An intent may only cite a belief that is promoted RIGHT NOW.
FOREIGN KEY (belief_id, belief_status)
  REFERENCES belief(id, status) ON UPDATE CASCADE     -- gate -> 23503

-- 3. And a live intent may not survive its belief losing that status.
CHECK (state <> 'live' OR belief_status = 'promoted') -- live_requires_promoted -> 23514
```
```
**File:** cmd/solvent-mcp/main.go (L88-88)
```go
	server.AddTool(&mcp.Tool{
```
**File:** cmd/solvent-mcp/main.go (L167-167)
```go
				"scenario": map[string]any{
```
**File:** cmd/solvent-mcp/main.go (L181-203)
```go
	server.AddTool(&mcp.Tool{
		Name:        "solvent_authorize_action",
		Description: "Record a live intent to take a real-world action, citing a belief as its warrant. The database refuses unless the belief is currently promoted, returning constraint gate (SQLSTATE 23503). Call this when the user asks to authorize, deploy, or act on a belief. Do not pre-check the belief's status.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief to cite as warrant",
				},
				"action": map[string]any{
					"type":        "string",
					"description": "Description of the real-world action to authorize",
				},
			},
			"required": []string{"scenario", "belief_id", "action"},
		},
	}, toolHandler("solvent_authorize_action"))
```
**File:** SOLVENT_ENGINEERING_GUIDE.md (L174-188)
```markdown
The MCP layer and Go application are intentionally not the final authority.

The design boundary is:

```text
Agent
  │
  ▼
MCP server
  │
  ▼
Solvent kernel / transaction code
  │
  ▼
CockroachDB
```
**File:** SOLVENT_ENGINEERING_GUIDE.md (L191-197)
```markdown
The agent does not get to decide that a belief is promotable or an intent is safe.

The MCP server is deliberately thin: validate input, call one kernel function, format the result. It does not implement belief semantics or SQL policy checks.

The kernel owns transaction discipline and reports the database result.

CockroachDB is the final invariant boundary.
```
**File:** Taskfile.yml (L28-30)
```yaml
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/002_corpus.sql
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/003_wizard.sql
      - docker exec -i solvent-crdb cockroach sql --insecure --database=fable < db/004_debt_vocabulary.sql
```
## Q4
design an AI swarm of agents based on Solvent whose goal is to verify physics claims per below: # Physics Verifier + Solvent — Architecture Discussion So Far

## 1. Goal

- **Objective:** Explore a Solvent demo in which AI swarm agents work on physics verification while Solvent provides the durable authority and evidence discipline around consequential decisions.
- **Core idea:** AI agents can reason, propose, attack, verify, and generate evidence, but they should not silently become the authority.
- **Important boundary:** Physics/math meaning stays outside the Solvent kernel. Solvent should enforce generic structural facts without understanding the content of the physics.

---

## 2. The Basic Mental Model

- **Agent/model:** Does reasoning, planning, mathematical exploration, proof search, counterexample search, and proposal generation.
- **Harness/runtime:** Controls agents, tools, credentials, files, network, sandbox, identity, and activity.
- **Policy/verifier layer:** Decides domain-specific semantic questions such as whether a claim is applicable, whether an attestation is required, and whether a review obligation is human-gated.
- **Solvent kernel:** Enforces generic durable facts such as belief identity, current promotion state, debt lifecycle, authority lifecycle, exact authority binding, and atomic transitions.
- **CockroachDB:** Provides the durable database invariants and concurrency guarantees.
- **Executor:** Performs the actual consequential external action only after the Solvent checks succeed.

Conceptually:

```text
Physics / domain semantics
        ↓
Policy / verifier / human review
        ↓
Solvent kernel
        ↓
CockroachDB invariants
        ↓
Execution
```

Important nuance: this is primarily an **authority flow**, not necessarily a literal dependency stack.

---

## 3. Single Canonical Mapping from Physics DAG to Solvent

The methodological rule is:

> Every distinct DAG concept should have exactly one canonical Solvent representation.

This is intended to prevent agents from choosing different valid-looking representations for the same concept and silently losing meaning.

### Mapping

1. **Node / claim**
   - One `Belief`, scoped to a `Scenario`.
   - Axiom, definition, result, or hypothesis is still represented as a belief rather than separate kernel tables.

2. **Epistemic kind**
   - `Belief.claim_type`.
   - The design discussed a many-to-few mapping such as:
     - axiom → Postulated
     - externally established result → Accommodated
     - internally proved result → Derived

3. **Mathematical / verification state**
   - Current Solvent state is represented through belief status plus debt rather than adding a separate `proved` state.
   - Roughly:
     - open → `entered` with unresolved debt
     - proved → `promoted`
     - refuted → `retracted`

4. **Review obligations**
   - Represented as debt / obligations.
   - Important correction: the kernel should not own a domain-specific vocabulary of what each debt *means*.

5. **Supporting artifacts**
   - `Evidence` rows linked to beliefs.
   - Evidence does not automatically discharge debt.

6. **Downstream consequential reliance**
   - `Action Intent` cites a belief.
   - The kernel checks whether that cited belief is currently promoted.

7. **Version revision**
   - Retract old belief → enter new belief → link with a `supersedes` relationship.
   - The old version remains readable and immutable.

8. **Staleness**
   - Prefer a computed query over stored staleness flags where the state can already be derived from belief relationships.
   - This avoids duplicated state that can drift.

9. **Separate research programs**
   - Use separate `Scenario`s.

10. **Actor capabilities**
    - Explicitly define which actors/tools can enter beliefs, add evidence, retire obligations, promote, etc.
    - Do not rely on descriptions or convention alone.

---

## 4. The Restricted-Belief Problem

### The problem

A belief may be true only under a narrower set of assumptions.

Example:

```text
Claim Y is true under {A, B, C}
```

A downstream action may operate under:

```text
{A, B}
```

The missing semantic question is:

```text
Is {A, B} a valid narrowing of {A, B, C}?
```

A foreign key can prove:

```text
intent.belief_id = belief.id
```

but it cannot prove:

```text
citing_context ⊆ belief.validity_conditions
```

That is a semantic-compatibility problem, not an identity problem.

### Key architectural conclusion

- The **kernel must not learn** about:
  - assumptions A/B/C
  - logical implication
  - subset relationships
  - mathematical meaning
  - domain-specific applicability rules

- Those belong in the **physics/policy/verifier layer**.

---

## 5. The Attestation Pattern

A proposed solution is to create an independent attestation belief.

Instead of:

```text
Restricted claim Y
       ↓
Action Intent directly cites Y
```

use:

```text
Restricted claim Y
       ↓
Attestation belief A
       ↓
Action Intent cites A
```

The attestation would mean, in effect:

> “Context X satisfies the applicability conditions for claim Y.”

The attestation gets its own evidence/debt/promotion lifecycle.

### Why this is attractive

- The semantic subset test remains outside the kernel.
- The final input to Solvent is a normal promoted belief.
- The kernel does not need to understand the mathematics.

### The critical gap

The attestation is useless as a security mechanism if an Action Intent can simply cite the original restricted belief directly.

That created the central **bypass question**:

> What actually prevents direct citation of a restricted belief?

This cannot be answered by saying “agents should use an attestation.” That is policy convention, not enforcement.

---

## 6. Case 1 vs Case 2 for Attestation Enforcement

Two cases were identified.

### Case 1 — Service-mediated Action Intent creation

- All Action Intent creation goes through a policy/service boundary.
- That layer can enforce:
  - restricted belief → attestation required
  - direct citation → rejected
- The Solvent kernel remains generic.
- No kernel understanding of “restricted” is required.

This is the preferred architecture **unless the threat model requires untrusted direct kernel access**.

### Case 2 — Untrusted direct kernel access

If an untrusted actor can directly invoke the kernel's authority-changing capabilities:

- The kernel does not understand “restricted.”
- Therefore it cannot enforce “attestation required” without importing domain semantics.
- That would be the wrong architectural move.

The resulting prerequisite is:

> Determine whether any untrusted actor can bypass the policy/service boundary and directly exercise authority-changing kernel capabilities.

---

## 7. The Human-Only Debt Retirement Problem Is the Same Trust-Boundary Problem

A similar proposed rule was:

> `semantic_applicability_unresolved` may only be retired through a human-facing interface.

The correction was that “human-only” should be treated as an **actor-capability policy**, not as domain semantics in the kernel.

But it introduces the same bypass question:

```text
Agent
  ↓
policy boundary
  ↓
human-gated retirement rule
```

What if an actor bypasses the policy layer and calls the kernel directly?

Then the “human-only” rule becomes decorative unless the trust boundary prevents the bypass.

### Therefore

The attestation-bypass problem and the human-only-retirement problem are really one problem:

> **Which actors are permitted to reach the authority-changing kernel capabilities, and what policy obligations must already be satisfied before they do?**

---

## 8. Debt Vocabulary: Important Correction

An earlier design treated the kernel's fixed debt enum as the canonical vocabulary.

That does **not** generalize cleanly.

### Correct separation

- **Kernel**
  - Stores generic obligation/debt identifiers.
  - Enforces debt lifecycle.
  - Enforces the generic promotion gate:
    - unresolved debt → promotion blocked
    - no unresolved debt → debt no longer blocks promotion
- **Policy/domain layer**
  - Defines what an obligation means.
  - Defines which obligations are relevant to a domain.
  - Defines when a particular obligation can be retired.

This avoids making the kernel domain-aware merely because a new domain needs a new obligation type.

### Important nuance

A generic debt key can still be opaque and structured.

The issue is not whether the kernel stores a string or identifier.

The issue is whether the kernel **understands the meaning** of that identifier.

---

## 9. Should Debt Have a Maximum Count?

The discussion rejected the idea of imposing an arbitrary semantic limit just for “discipline.”

Two separate concepts must be distinguished.

### Semantic limit

Example:

```text
A belief may have at most 10 obligations.
```

This would encode assumptions about how a good workflow should look.

That does not belong in the kernel merely because 100 feels excessive.

### Technical/resource limit

Example:

```text
Reject pathologically large obligation collections
because they create unacceptable request/transaction/resource cost.
```

That can be legitimate.

### Design rule

- Do **not** assume an array needs a limit before verifying the actual representation.
- First inspect whether debt is:
  - a JSON/JSONB array
  - SQL array
  - normalized rows
  - another representation
- If a limit is necessary, determine which layer should own it.
- Prefer a resource/input boundary where that is sufficient.
- Do not invent an arbitrary number such as 6, 10, or 100 without evidence.

The question is:

> What security or resource property does the number enforce?

Not:

> What number looks disciplined?

---

## 10. Trust Boundary: Five Distinct Attack Vectors

The trust-boundary review was expanded into five separate vectors because “can a caller reach the kernel?” is too compressed.

### Vector 1 — Direct kernel function calls

Problem:

- An untrusted actor might directly invoke authority-changing kernel functions.

Controls discussed:

- Runtime authentication/authorization at the process/API boundary.
- I-7 remains useful but is **not** runtime access control.

I-7 means:

> Developers should not accidentally create raw DB write paths outside the kernel.

It does **not** mean:

> An untrusted actor cannot invoke a legitimate kernel function.

Those are different properties.

### Vector 2 — Other internal authority-mutating APIs

Problem:

- A second internal API or tool may expose authority-changing behavior without the same policy controls.

Need to enumerate all such entry points.

Possible controls:

- Narrow authority-changing interfaces.
- Explicit capability surfaces.
- Runtime access control.
- Structural call-path discipline.

### Vector 3 — Lower-level service path bypassing policy

This was considered the most suspicious evolving-codebase risk.

Examples:

- future admin tool
- batch script
- debugging endpoint
- internal service
- migration utility

Each may accidentally bypass the policy layer.

The open design question is whether a capability/session-context mechanism is actually needed.

Important restraint:

> Do not add capability tokens preemptively.

First determine from the real system:

- how many authority-mutating entry points exist
- who can call them
- which are public/internal
- whether a bypass exists
- whether the permitted call graph can be made structurally obvious

Only then decide whether a capability mechanism is justified.

### Vector 4 — Direct CockroachDB writes

This is fundamentally a database-credentials problem.

The important question is:

> Which processes have write credentials to `belief`, `debt`, `action_intent`, and related authority state?

Ideal boundary:

```text
policy/service + kernel
       ↓
DB write credentials

reporting / analytics / ops tools
       ↓
no write access
```

If other processes have write credentials, they can bypass policy/kernel controls regardless of how elegant the application design is.

This is therefore a deployment configuration fact.

### Vector 5 — Deployment/configuration exposure

Examples:

- overly permissive IAM
- debug flag
- raw SQL console
- exposed internal endpoint
- unsafe runtime configuration

This is operational security.

No amount of clean physics-to-Solvent modeling automatically protects against it.

It still belongs on the explicit review checklist.

---

## 11. I-7 Is Not the Whole Trust Boundary

The architecture now distinguishes:

```text
I-7
→ build-time code-hygiene invariant

authentication/authorization
→ runtime actor control

DB credentials
→ database trust boundary

deployment configuration
→ operational trust boundary
```

A clean Solvent codebase can demonstrate I-7 and still be insecure if the deployment gives an attacker authority-changing access.

---

## 12. The Trusted Computing Base (TCB)

A more useful definition emerged:

> **A layer is part of the trusted computing base whenever bypassing or corrupting it can cause an otherwise unauthorized consequential action.**

This is stronger than:

> “The kernel is small, therefore trusted.”

Smallness is useful because it makes auditing easier, but **consequence** is the real criterion.

### Consequence

The physics policy/verifier layer can become part of the TCB because it may decide:

- whether a restricted belief can be relied upon
- whether an attestation is required
- whether an obligation can be retired
- whether a human review is required

Even though these are outside the kernel.

---

## 13. The Physics / Policy Layer Cannot Be Sloppy

Once the policy layer makes security-relevant decisions, it needs similar engineering discipline to the kernel.

Desired properties include:

- explicit rules
- versioning
- reviewability
- auditable decisions
- non-bypassable enforcement
- explicit actor capabilities
- clear provenance
- controlled changes

The kernel stays small by refusing domain semantics.

That does **not** mean the policy layer is allowed to be informal.

---

## 14. The Policy Versioning Question

A proposed policy model looked like:

```text
Policy P17
Attestation A42
A42 is valid under P17
```

This created a deeper question.

The statement:

```text
A42 is valid under P17
```

is itself a claim.

Under the general Solvent worldview, claims are beliefs with:

- their own evidence
- their own debt
- their own promotion
- their own retraction/supersession behavior

Therefore a policy version should not merely be a bare label attached to an attestation if the validity of the attestation depends on that policy.

### Desired conceptual behavior

```text
Policy P17
   ↓
used to evaluate attestation A42
```

If P17 is later retracted or superseded:

```text
P17 invalid
   ↓
A42 may become stale/invalid
   ↓
downstream reliance must be discoverable/rejected
```

This mirrors the existing identity-drift discipline used for ordinary beliefs.

---

## 15. The Kernel Need Not Understand Policy Meaning

Even if policy versions become tracked, the kernel does not need to know what P17 means.

A higher-level verifier/policy system can establish a fact such as:

```text
A42 was evaluated under P17
and the applicability test passed
```

while Solvent's kernel continues to enforce only generic facts such as:

```text
A42 is a promoted belief
```

The domain meaning remains outside the kernel.

---

## 16. A Subtle Layering Correction

The clean diagram:

```text
Physics/domain semantics
        ↓
Policy
        ↓
Kernel
```

is an **authority diagram**, not a pure dependency diagram.

The policy layer may itself depend on Solvent's ledger to track:

- policy beliefs
- policy attestations
- retractions
- supersession
- debt
- audit facts

So both statements can be true:

1. The kernel has no opinion about policy meaning.
2. The policy layer may rely on the kernel's durability and atomicity to keep its own security-relevant history coherent.

This distinction should be explicit in architecture documentation.

---

## 17. The Phase 0–5 Investigation Loop

The proposed sequencing was refined from a simple one-way order into an iterative loop.

### Phase 0 — Provisional model

Draft the debt/policy model enough to expose important future call signatures and security properties.

### Phase 1 — Trust-boundary investigation

Run the five-vector investigation against the provisional model.

For every authority-changing entry point, ask:

- Who can call it?
- Can it bypass policy?
- What does it accept?
- Which fields does it trust?
- What verifies those fields?

### Phase 2 — Refine debt/policy model

Resolve:

- vocabulary ownership
- representation
- auditability
- resource limits
- retirement semantics
- actor capabilities
- policy-version representation
- whether a policy version is itself represented as a belief

### Phase 3 — Trust-boundary re-check

Re-run the critical V3 checks against the refined model.

This matters because Phase 2 can change what the authority-mutating calls actually need to enforce.

### Phase 4 — Relationship/schema verification

Verify things such as:

- `belief_edge` structure
- edge typing/labels
- representation of `supersedes`
- attestation relationships

### Phase 5 — Freeze decision

Only after the model and trust boundary stabilize should the design be considered complete enough for implementation/freeze decisions.

---

## 18. The Bidirectional Dependency

The important correction is:

```text
A/V3 ↔ B
```

not simply:

```text
A/V3 → B
```

Why?

- V3 constrains what policy-mediated operations are allowed to look like.
- B can change what V3 itself needs to check.

Example:

If debt retirement becomes:

```text
RetireDebt(
    belief,
    obligation,
    policy_version,
    review_context,
    actor,
    evidence
)
```

then each new field creates another potential truthfulness problem.

The trust review must ask not only:

> Can an unauthorized actor call this?

but also:

> For each security-relevant field, what verifies that its content is truthful?

---

## 19. Every New Policy Field Is Another Place a Lie Can Enter

Example fields:

- `actor`
  - Is it truly a human or merely a script claiming to be human?

- `policy_version`
  - Is it the correct/current applicable policy rather than an old permissive one?

- `evidence`
  - Does the referenced artifact actually establish what retirement claims it establishes?

This is distinct from path authorization.

Therefore the trust-boundary review must examine both:

```text
Who can call?
```

and:

```text
Why should we believe the arguments?
```

---

## 20. Domain Semantics as Part of the TCB

The same consequence test should be applied one level above policy.

If changing domain semantics can cause an otherwise unauthorized consequential action, then the domain-semantics mechanism is also part of the TCB.

This does not mean domain semantics must move into Solvent.

It means the system should explicitly govern the mechanism carrying those semantics.

Examples:

- versioned policy repository
- reviewed verifier configuration
- controlled mathematical specification
- auditable semantic rules

The physical/semantic layer can remain outside the Solvent kernel while still being recognized as security-critical for the demo.

---

## 21. Canonical Trust Architecture

The emerging architecture is:

```text
                TRUSTED COMPUTING BASE

    ┌───────────────────────────────────────────┐
    │ Domain semantics                          │
    │                                           │
    │ What claims mean                          │
    │ What assumptions imply                    │
    │ What obligations mean                     │
    └───────────────────────────────────────────┘
                      ↓
    ┌───────────────────────────────────────────┐
    │ Policy / verifier / human review          │
    │                                           │
    │ Applicability                             │
    │ Attestation requirements                  │
    │ Human-gated obligations                   │
    │ Actor capabilities                        │
    │ Policy versions                           │
    └───────────────────────────────────────────┘
                      ↓
    ┌───────────────────────────────────────────┐
    │ Solvent kernel                            │
    │                                           │
    │ Identity                                  │
    │ Currency                                  │
    │ Debt lifecycle                            │
    │ Promotion gate                            │
    │ Authority binding                         │
    │ Atomic transitions                        │
    └───────────────────────────────────────────┘
                      ↓
                  Execution
```

This is primarily an authority decomposition.

The kernel may also serve as the durable substrate on which higher-level policy history is recorded.

---

## 22. What Must NOT Enter the Solvent Kernel

The negative boundary is intentionally explicit.

The kernel should never need code that inspects:

- what a physics assumption means
- why an obligation exists
- whether one condition implies another
- whether a context is mathematically valid
- what a domain-specific debt item means
- whether a citation is semantically applicable

Instead, those become higher-level facts.

Solvent receives something structurally simple:

```text
"This belief is promoted."
```

and enforces the consequences of that fact.

---

## 23. What Likely Belongs in the Physics Demo

Once the kernel is frozen, the physics-verifier work should live outside the kernel.

A likely structure is:

```text
demos/
  physics-verifier/
    README.md
    swarm/
    verifier/
    policy/
    attestations/
    scenarios/
    evidence/
    adapters/
    scripts/
```

The demo should show how a swarm/harness becomes **Solvent-ready** without changing Solvent's kernel.

The swarm can perform:

- theorem exploration
- proof attempts
- counterexample search
- numerical checks
- adversarial attacks
- evidence collection
- semantic applicability checks
- attestation generation

Solvent remains the authority boundary.

---

## 24. Important Test Philosophy for the Physics Verifier

The same “prove it structurally” discipline should apply to the demo.

Do not rely on:

> “Agents are instructed to cite the attestation.”

Prefer tests that prove:

```text
restricted belief + no attestation
    → intent rejected

restricted belief + valid promoted attestation
    → intent allowed

attestation based on retracted policy/claim
    → reliance becomes invalid/discoverable

human-only obligation + agent caller
    → rejected

human review + valid policy context
    → obligation can be retired

T1-bound intent + T2 authority
    → rejected
```

The final physics demo should make the same principle visible at the application level:

> **The swarm can propose; the verifier can reason; the policy layer can interpret; Solvent decides whether the resulting authority state is structurally valid.**

---

## 25. Current Open Questions Before Implementation

1. **Trust boundary**
   - Can any untrusted actor directly reach the kernel?
   - Can any internal service bypass the policy layer?
   - What processes have direct DB write credentials?
   - What deployment/configuration paths expose authority-changing capabilities?

2. **Debt representation**
   - Is debt an array, rows, or another structure?
   - Does it actually need a technical cardinality bound?
   - Which layer should enforce that bound?

3. **Debt vocabulary**
   - Can domain vocabulary live entirely outside the kernel?
   - Are current fixed debt enums actually kernel semantics or merely one deployment's policy vocabulary?

4. **Attestation enforcement**
   - Is Action Intent creation always service-mediated?
   - If so, can the policy layer reliably enforce restricted-belief → attestation-required?

5. **Human-gated obligations**
   - Is the human-only rule enforced by the same trusted policy boundary?

6. **Policy versioning**
   - Should a policy version itself be a Belief?
   - If a policy version is retracted, how does invalidity propagate to attestations?

7. **Relationship model**
   - Does `belief_edge` support typed relationships?
   - Is `supersedes` distinguishable from dependency/attestation relationships?

8. **Argument truthfulness**
   - For every security-relevant argument added to policy operations, what verifies its content rather than merely its presence?

---

## 26. Current Architectural Position

The strongest conclusions reached so far are:

- Solvent should remain a **small trusted authority core**, not a physics reasoning engine.
- The physics verifier should not cause new kernel features simply because semantic verification is difficult.
- Restricted-claim applicability is a **semantic policy problem**, not a kernel mathematics problem.
- Attestation only works if the policy/service boundary prevents direct bypass.
- Human-only retirement has the same trust-boundary dependency.
- Debt vocabulary should be **policy-owned**, while the kernel owns generic debt lifecycle/invariants.
- A debt-count cap should only exist if there is a concrete technical/resource reason; an arbitrary discipline number does not belong in the kernel.
- I-7 is a useful build-time invariant but is not runtime authentication or access control.
- The trust boundary has five distinct attack vectors and needs five distinct answers.
- The policy/verifier layer may be part of the TCB even though it is outside the kernel.
- Domain semantics may also be part of the TCB when their corruption can change authorization outcomes.
- Policy versioning should be treated with the same identity, provenance, revision, and stale-propagation discipline as other security-relevant claims.
- The Phase 0–5 process should iterate because the trust-boundary design and debt/policy design constrain each other.
- The kernel-freeze rule remains intact: after freeze, new capabilities should default to service, adapter, executor, deployment, policy, demo, or documentation layers.

---

## 27. The Core Principle

The emerging design can be reduced to one sentence:

> **The physics system is allowed to understand meaning; the policy layer is allowed to decide semantic applicability; Solvent is allowed to enforce durable structural authority—but the security of the whole system depends on preventing untrusted actors from bypassing the layers above the kernel that carry meaning.**

That is the main architectural lesson driving the current physics-verifier exploration.
### Answer
Below is a concrete swarm architecture, mapped onto Solvent's actual primitives (not idealized ones), organized by the categories in the document you attached.

## 1. Agent roles in the swarm

| Agent | Function | Solvent interaction |
|---|---|---|
| **Explorer agents** | Proof search, counterexample search, numerical checks | Read-only against `solvent_ledger`; propose claims but cannot call `EnterBelief` directly — routed through a policy/service layer |
| **Attacker agents** | Adversarial search for restricted-condition violations, edge-case exploitation | Read promoted beliefs; write `evidence` rows of `provenance_class='reproducible_artifact'` when they find a counterexample solvent:53-55  |
| **Verifier/policy agent(s)** | Semantic applicability checks, attestation evaluation, debt-vocabulary interpretation | Not a Solvent primitive at all — this sits entirely in your policy layer, calling kernel functions on the swarm's behalf |
| **Human reviewer** | Discharges human-gated debt items | The only actor permitted to retire certain debt items — enforcement is entirely a caller-identity problem, discussed below |

None of these roles exist as a Solvent concept. Solvent has exactly one notion of "actor": whoever calls the kernel function. This is precisely the gap your document's §3.10 and §10 identify — there is currently no actor-capability primitive in the codebase.

## 2. Canonical mapping (validated against the actual schema, not the idealized one)

Your document's mapping is largely correct but two points need correction against what's actually shipped:

- **`claim_type`** is exactly `{derived, accommodated, postulated}` solvent:86-91  — your axiom→Postulated / external-result→Accommodated / internal-proof→Derived mapping fits this enum exactly with no schema change needed.
- **`belief_edge.kind`** is currently constrained to exactly `{derives, contradicts}` by a schema CHECK solvent:420-427 . Your document proposes a `supersedes` relationship for version revision (§3.7) and an attestation relationship (§5). **Neither exists yet.** This is a real, first-class design decision, not a detail: either (a) extend the CHECK to add `supersedes`/`attests` kinds, or (b) encode supersession as a `derives` edge plus an off-kernel convention that the physics-policy layer interprets. Per `AGENTS.md`'s freeze rule — "the architecture is frozen... No further tables without explicit approval" solvent:52-57  — even a new *edge kind* (not a new table) should go through the same scrutiny, because it's a kernel schema change, and the kernel is what's frozen, not the demo layer built on it.
- **Debt is a `TEXT[]` SQL array**, not JSONB, confirmed by the encoding assumption chain (C9/C15) and the `array_remove` mechanism solvent:140-140 solvent:169-204 . This answers your §9 question directly: there is no cardinality limit in the schema today, `array_remove` is O(n) and unindexed, and nothing currently guards against an unbounded debt array — if you need a resource bound for the physics demo, it belongs at the policy/service layer that constructs `EnterBelief` calls, not in the kernel.

## 3. The restricted-belief / attestation problem — what Solvent gives you today

Solvent's actual enforcement mechanism for "downstream reliance requires current promotion" is the composite FK gate: `action_intent(belief_id, belief_status) REFERENCES belief(id, status)` solvent:122-132 . This FK proves identity + currency. It cannot and does not prove semantic subset validity — your document's diagnosis is correct, and the codebase confirms it: nothing in the kernel or schema reasons about assumption sets.

The attestation pattern you propose (`Action Intent` cites `Attestation belief` instead of `Y` directly) is mechanically identical to how `solvent_authorize_action` already works — it's just another belief with its own debt/evidence/promotion lifecycle solvent:181-203 . **The bypass question is real and unsolved by the current codebase**: nothing prevents an `action_intent` from citing the restricted belief `Y` directly instead of the attestation belief. The only place this could be enforced today is exactly where your document lands — a service boundary in front of `IntentOnPromoted`, since the kernel function itself takes a bare `beliefID` with no notion of "this belief requires an attestation before citation" solvent:107-107 .

## 4. Trust boundary — what's actually enforced today vs. what's assumed

This is the most important thing the codebase confirms about your document's §6/§10/§11 analysis. Solvent currently has exactly one mechanically-enforced boundary and it is **not** an actor/authentication boundary:

- **I-7 is a static grep, not runtime access control.** `scripts/check_i7.sh` and the Taskfile's I-7 MCP boundary check grep for `.Exec`/`.ExecContext`/`.Begin` calls and raw SQL text (`INSERT INTO`, `UPDATE ... SET`, etc.) in `cmd/solvent-mcp` and `internal/view` at build/CI time solvent:251-266 . This proves developers didn't accidentally write a raw DB path — it says nothing about whether an untrusted network caller can invoke `solvent_authorize_action` legitimately. Your document's distinction between "I-7 = build-time hygiene" and "authentication = runtime actor control" is exactly right and is *not yet closed* by anything in the repo.
- **Session/scenario isolation exists only at the wizard demo layer**, via cookies — `TestW17_TwoClientsGetSeparateScenarios` proves two HTTP clients get distinct `scenario_id`s and can't see each other's screen state solvent:184-208 , and `TestW32` proves a new visitor can't perturb the read-only ledger counts solvent:19-29 . This is cookie-based session separation for a demo UI, not an actor-capability or authorization system. It does not generalize to "human-only debt retirement" or "agent vs. human caller" distinctions — there is no identity concept here beyond a browser cookie.
- **The MCP server has no authentication layer at all.** It's a stdio server that trusts whatever calls it solvent:1-4 . For a swarm with multiple agent types (explorer, attacker, human-review-gate), this means **Vector 1 and Vector 2 from your document are currently wide open** in the shipped code — anything that can speak the MCP protocol can call any of the six tools with any arguments, including `solvent_promote` and `solvent_authorize_action`.

**Conclusion for the swarm design: you must build the actor-capability/policy-service layer from scratch.** It does not exist in Solvent today in any form — not even a stub. This confirms your document's own conclusion (§10, "Determine whether any untrusted actor can bypass the policy/service boundary") is not a hypothetical to investigate later; it's a hard prerequisite because nothing currently answers it.

## 5. Concrete swarm harness design

Given the above, here's how I'd structure it:

**A. Kernel (unmodified Solvent, per the freeze rule)** — `belief`, `belief_edge`, `evidence`, `action_intent`, exactly as shipped solvent:54-57 . Physics semantics never enter here.

**B. Policy/service layer (new, sits in front of MCP or replaces the MCP server for this demo)**
- Owns a domain debt vocabulary distinct from `kernel.FullDebt` (which is etcd-specific and explicitly "must never appear in the kernel" per the domain-independence table) solvent:490-510 . Physics obligations (`needApplicabilityProof`, `needAttestation`, `needPeerCheck`) are policy-owned strings passed into `EnterBelief`'s debt argument — the kernel doesn't need to know what they mean, only that they're non-empty.
- Owns actor identity and enforces: (a) attestation-required routing for restricted beliefs, (b) human-only retirement gates. Since nothing in the kernel does this, it must be a real authentication/authorization layer wrapping the existing thin MCP handlers — not decoration on top of them.
- Owns policy-version tracking. Given `belief_edge` currently only supports `derives`/`contradicts`, representing "attestation A42 evaluated under policy P17" as its own belief with a `derives` edge to P17 is achievable *today* without a schema change — cleaner than inventing a new edge kind, and consistent with your document's §14-15 conclusion that policy versions should themselves be beliefs.

**C. Swarm agents (explorer/attacker/verifier)** — MCP or direct policy-layer clients, never given raw kernel/DB credentials. This is Vector 4 in your document: **whoever holds CockroachDB write credentials is the actual trust boundary**, and the swarm processes must not be among them regardless of how the policy layer is designed.

**D. Execution** — only fires on `solvent_authorize_action` success, i.e., only after the FK gate accepts a promoted (attestation) belief citation solvent:181-203 .

## 6. What must be built vs. what's reused

| Component | Status in Solvent today |
|---|---|
| Belief/evidence/intent lifecycle, promotion gate, cascade | Fully built, tested (I-1 through I-8), reusable as-is solvent:102-117  |
| Scenario-scoped isolation for parallel research programs | Built and tested at the kernel level solvent:314-343 , though scoped isolation at the *agent/session* level exists only in the demo wizard, not generally |
| Domain-agnostic, swappable debt vocabulary | Not built — `kernel.FullDebt` is a fixed six-item Go slice; policy-layer vocabulary is entirely new work |
| Actor identity / authentication / capability enforcement | **Not built at all** — this is the single largest gap between your document's architecture and the current repo |
| `supersedes` / attestation-typed edges | Not built — `belief_edge.kind` CHECK currently allows only `derives`/`contradicts` |
| Policy-version-as-belief with stale-propagation | Not built, but achievable without kernel changes using existing `derives` edges |
| Gate-adequacy / meta-review beliefs | Not built — would itself be "just another belief" but no convention exists yet |

The practical upshot for a physics-verifier demo: Solvent gives you a fully proven, invariant-enforced *substrate* for claims/evidence/actions, but essentially the entire "who is allowed to do what, and is the attestation actually being enforced rather than merely conventional" layer — which is where your document correctly identifies the real security work lives — has to be built as new code sitting in front of the kernel, since nothing resembling it exists in the shipped repository today.
### Citations
**File:** plans/pre/IMPLEMENTATION_CONTRACT.md (L53-55)
```markdown
- `evidence.provenance_class` ∈ {`external_feed`,`reproducible_artifact`,`live_scan`,
  `operator_asserted`}; `content_sha256` NOT NULL.
- `belief.embedding` is nullable; the ledger is fully functional with every embedding NULL.
```
**File:** plans/pre/IMPLEMENTATION_CONTRACT.md (L102-117)
```markdown
## 5. Invariants That Must Never Be Violated

- **I-1** No belief with `status='promoted'` has non-empty `debt`. (Schema CHECK.)
- **I-2** No belief with `status='promoted'` has `final_truth=true`. (Schema CHECK.)
- **I-3** No `live` action_intent references a belief that is not currently `promoted`.
  (Composite FK + `live_requires_promoted` CHECK.)
- **I-4** A belief carrying a live intent cannot transition to `retracted` without that intent
  first being cancelled. (Enforced by ON UPDATE CASCADE propagating status into the intent row,
  which then violates `live_requires_promoted` unless already cancelled.)
- **I-5** `AuditLiveOnNonPromoted` returns 0 in every committed state.
- **I-6** All invariant tests pass with every `embedding` column NULL. Vectors are an
  optimization, never part of belief semantics.
- **I-7** Every kernel write goes through `crdb.ExecuteTx`. No raw `db.Exec`/`db.Query` writes.
- **I-8** `RetractCascade` is a single transaction. Cancel-before-retract ordering is mandatory.

Any change that would require weakening or removing any of I-1..I-8 is out of scope: STOP (§9).
```
**File:** plans/M1/plan1.md (L86-91)
```markdown
type ClaimType string
const (
    Derived      ClaimType = "derived"
    Accommodated ClaimType = "accommodated"
    Postulated   ClaimType = "postulated"
)
```
**File:** plans/M1/plan1.md (L107-107)
```markdown
| `IntentOnPromoted(ctx, scenarioID, beliefID, action string) error` | Inserts a `live` intent. On SQLSTATE **23503** returns an error satisfying `errors.Is(err, ErrActionOnUnpromoted)`. |
```
**File:** plans/M1/plan1.md (L140-140)
```markdown

```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L420-427)
```markdown
### 8.2 Edges

| Edge Type | Table | Meaning |
|---|---|---|
| **Supports** | `evidence.belief_id` (FK) | This evidence supports this belief. |
| **Derives** | `belief_edge` (`kind='derives'`) | This belief depends on that belief. |
| **Contradicts** | `belief_edge` (`kind='contradicts'`) | This belief contradicts that belief. |
| **Governs** | `action_intent.belief_id` (FK) | This intent is governed by this belief. |
```
**File:** plans/PRD/DATA_PIPELINE_SPEC.md (L490-510)
```markdown
### etcd-Specific (EvidenceFeed implementations)

These concepts belong to the etcd domain and must never appear in the kernel:

- CVE IDs, CISA KEV catalog
- etcd version strings, release tags
- GitHub issue/PR event types
- Maintainer identity and comment semantics
- "vulnerable", "patched", "exploited"
- `notify_oncall`, `defer_patch`, `page_incident_comminator`
- Version comparison logic (semver)
- The specific six debt items (though the mechanism of debt is generic)

### Kernel-Generic (internal/kernel)

These concepts belong to the transactional belief ledger and are domain-agnostic:

- Belief, evidence, action_intent (the three tables + belief_edge)
- Belief lifecycle: entered → promoted → retracted
- Intent lifecycle: live → cancelled/executed
- Debt retirement and promotion gating
```
**File:** AGENTS.md (L52-57)
```markdown
The architecture is frozen. Do not redesign.

The ledger is four tables — `belief`, `belief_edge`, `evidence`, `action_intent`. The corpus and
wizard layer adds `corpus_issue`, `belief_corpus_citation`, and `refusal_log`. No further tables
without explicit approval. The kernel remains domain-agnostic: evidence sources are `EvidenceFeed`
implementations, and changing the source must not require kernel changes.
```
**File:** internal/wizard/discharge.go (L169-204)
```go
// retrievalDebtsRetired counts how many of the two retrieval checks are already
// discharged, which is how many citations are already spoken for.
//
// It reads the debt array rather than the evidence table because a citation-backed
// discharge writes no evidence row -- the citation IS the record. The same unnest
// checks() uses, so the count and the receipts can never disagree.
func (s *Server) retrievalDebtsRetired(ctx context.Context, beliefID string) (int, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT unnest(debt) FROM belief WHERE id = $1::UUID`, beliefID)
	if err != nil {
		return 0, fmt.Errorf("wizard: read debt for retrieval count: %w", err)
	}
	defer func() { _ = rows.Close() }()

	outstanding := map[string]bool{}
	for rows.Next() {
		var d sql.NullString
		if err := rows.Scan(&d); err != nil {
			return 0, fmt.Errorf("wizard: scan debt for retrieval count: %w", err)
		}
		if d.Valid {
			outstanding[d.String] = true
		}
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	n := 0
	for item := range retrievalChecks {
		if !outstanding[item] {
			n++
		}
	}
	return n, nil
}
```
**File:** README.md (L122-132)
```markdown
```sql
-- 1. Debt must be empty before a belief can be promoted.
CHECK (status <> 'promoted' OR debt = '{}')          -- promoted_is_debt_free -> 23514

-- 2. An intent may only cite a belief that is promoted RIGHT NOW.
FOREIGN KEY (belief_id, belief_status)
  REFERENCES belief(id, status) ON UPDATE CASCADE     -- gate -> 23503

-- 3. And a live intent may not survive its belief losing that status.
CHECK (state <> 'live' OR belief_status = 'promoted') -- live_requires_promoted -> 23514
```
```
**File:** cmd/solvent-mcp/main.go (L1-4)
```go
// Command solvent-mcp is a stdio MCP server exposing the Solvent transactional
// belief ledger as six tools. The server is an adapter — it has no opinion about
// beliefs. Every tool handler is exactly three moves: unmarshal → kernel call →
// format.
```
**File:** cmd/solvent-mcp/main.go (L181-203)
```go
	server.AddTool(&mcp.Tool{
		Name:        "solvent_authorize_action",
		Description: "Record a live intent to take a real-world action, citing a belief as its warrant. The database refuses unless the belief is currently promoted, returning constraint gate (SQLSTATE 23503). Call this when the user asks to authorize, deploy, or act on a belief. Do not pre-check the belief's status.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"scenario": map[string]any{
					"type":        "string",
					"enum":        []string{"track1", "track2"},
					"description": "Scenario the belief belongs to",
				},
				"belief_id": map[string]any{
					"type":        "string",
					"description": "UUID of the belief to cite as warrant",
				},
				"action": map[string]any{
					"type":        "string",
					"description": "Description of the real-world action to authorize",
				},
			},
			"required": []string{"scenario", "belief_id", "action"},
		},
	}, toolHandler("solvent_authorize_action"))
```
**File:** Taskfile.yml (L251-266)
```yaml
      # I-7 MCP boundary — no direct writes in cmd/solvent-mcp or internal/view
      - |
        for d in cmd/solvent-mcp internal/view; do
          [ -d "$d" ] || { echo "FAIL: $d missing"; exit 1; }
        done
        # Stage 1: reject write/transaction entry points
        if grep -rnE --include='*.go' '\.(Exec|ExecContext|Prepare|PrepareContext|Begin|BeginTx)\(' cmd/solvent-mcp internal/view; then
          echo "FAIL: direct write/transaction call in MCP/view"
          exit 1
        fi
        # Stage 2: reject write SQL text (case-insensitive, SQL-shaped)
        if grep -rniE --include='*.go' '(INSERT[[:space:]]+INTO|UPDATE[[:space:]]+[a-z_]+[[:space:]]+SET|DELETE[[:space:]]+FROM|CREATE[[:space:]]+(TABLE|INDEX|DATABASE)|DROP[[:space:]]+(TABLE|INDEX|DATABASE))' cmd/solvent-mcp internal/view; then
          echo "FAIL: write SQL text in MCP/view"
          exit 1
        fi
        echo "I-7 MCP boundary: PASS"
```
**File:** internal/wizard/http_test.go (L184-208)
```go
func TestW17_TwoClientsGetSeparateScenarios(t *testing.T) {
	s := newServer(t, &fakeEmbedder{})
	a, b := newClient(t, s), newClient(t, s)

	// Both act, so both hold a real scenario id. If only one did, the other's id would
	// be the empty string and the inequality below would hold for the wrong reason.
	a.do(http.MethodPost, wizard.Prefix+"/api/promote", "{}")
	b.do(http.MethodPost, wizard.Prefix+"/api/search", `{"query":"anything"}`)

	_, sa := a.do(http.MethodGet, wizard.Prefix+"/api/state", "")
	_, sb := b.do(http.MethodGet, wizard.Prefix+"/api/state", "")
	if sa["scenario_id"] == sb["scenario_id"] {
		t.Fatal("two clients share a scenario")
	}
	if sa["screen"] != wizard.ScreenDischarge {
		t.Errorf("client a screen = %v, want DISCHARGE", sa["screen"])
	}
	if sb["screen"] != wizard.ScreenAsk {
		t.Errorf("client b screen = %v, want ASK — a's refusal leaked", sb["screen"])
	}
	if sa["scenario_id"] == "" || sb["scenario_id"] == "" {
		t.Fatalf("a scenario id is empty; both judges must have acted: a=%v b=%v",
			sa["scenario_id"], sb["scenario_id"])
	}
}
```
**File:** internal/wizard/ledger_isolation_test.go (L19-29)
```go
// --- W-32: a wizard visitor cannot move the ledger's numbers -------------------
//
// This is the precondition for ever mounting the wizard at /. The landing page and
// /ledger render five counts; if any of them aggregated across scenarios, every
// judge who opened the demo would inflate the numbers a later judge reads.
//
// The queries below are copied verbatim from getLedgerCounts (handlers.go:25-53).
// Asserting only that they don't move would pass even if the filter were removed
// and the test simply got lucky, so the unfiltered count is measured alongside:
// it MUST rise. That is what proves the scenario filter is doing the work.
func TestW32_WizardDoesNotDisturbLedgerCounts(t *testing.T) {
```
**File:** internal/corpus/corpus_test.go (L314-343)
```go
func TestC06_ScenarioIsolation(t *testing.T) {
	ctx := context.Background()

	// An exact duplicate of the nearest planted document, in a different scenario.
	if _, _, err := corpus.Insert(ctx, shared, corpus.Issue{
		ScenarioID: scenarioB, IssueNumber: 9001, Title: planted[0].title, State: "closed",
		URL: "https://x/b-9001", ContentSHA256: "iso-b", Embedding: corpus.HashEmbed(planted[0].title),
	}); err != nil {
		t.Fatal(err)
	}

	hits, err := corpus.Search(ctx, shared, scenarioB, corpus.HashEmbed(plantedQuery), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("scenario B returned %d hits, want exactly its own 1", len(hits))
	}

	// The same issue_number exists in both scenarios, which is only possible
	// because identity is (scenario_id, issue_number) rather than issue_number.
	var n int
	if err := shared.QueryRowContext(ctx,
		`SELECT count(*) FROM corpus_issue WHERE issue_number = 9001`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("issue 9001 present %d times across scenarios, want 2", n)
	}
}
```
