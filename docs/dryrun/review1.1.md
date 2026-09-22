I agree with most of that review, and it strengthens the picture from my earlier assessment. The dry run has demonstrated the core loop; what remains is **provenance completeness and adversarial-review discipline**, not another broad architecture change.

One important refinement: I would **not immediately enforce “adversarial targets must have open debt” or “target must never be retracted/promoted.”** An adversarial agent may legitimately challenge a promoted claim, and a historical/retracted claim may still be worth analyzing. The safer invariant to audit first is: **an adversarial contradiction should target an existing canonical belief in the same scenario, not a local belief created inside the same packet**, and its provenance should make the review lineage explicit. That keeps the gate semantic rather than inventing a new status policy.

The current context audit already establishes that `get_context` now provides the full scenario belief surface, and the edge projection completes the graph view.  

## Consolidated assessment

### Highest priority: provenance spine

I agree that `origin_packet_id` is preferable to a vague `created_by_submission_id`.

The model should eventually be:

```text
packet_submission
    |
    +---- creates ----> belief
    +---- creates ----> evidence
    +---- creates ----> edge
    +---- creates ----> task
```

with the important property that the packet submission is immutable provenance, not authority.

For retries/idempotency, anchoring lineage to `packet_id` is cleaner because one logical packet has one identity. The audit should determine whether that relationship already exists somewhere transiently and whether it is actually persisted.

### Highest-value behavioral finding: EBP category error

This is the strongest adversarial finding.

The adversarial agent treated:

```text
entered + eight open debts
```

as an EBP violation.

That is not consistent with the governing doctrine:

```text
Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
```

So the adversarial agent demonstrated a real **protocol-calibration failure**. I agree that this belongs in the shared constitution, not just the adversarial prompt:

> Do not call an unresolved obligation an EBP violation unless a specific active EBP or Domain Pack rule is actually violated.

This is exactly the sort of failure the adversarial branch should surface.

### Coverage

Also agree with the two-stage approach:

**Now:** require the adversarial agent to report a disposition for every visible belief:

```text
B1 — challenged
B2 — no material objection
B3 — unresolved / insufficient basis
...
```

**Later:** persist review coverage only if repeated dry runs prove that prose-level coverage is inadequate.

Do not build a review-coverage subsystem yet.

### Human lineage

The current graph is now available to the adversarial agent, but the human still sees mostly a flat belief list.

The next useful projection is therefore:

```text
belief
  → origin packet
  → agent / harness / model
  → relation(s)
```

A graph visualization is unnecessary. Tables are sufficient.

### Proposed tasks

I agree this deserves elevated priority.

The adversarial agent reported two tasks, but the dashboard did not visibly show them. This is exactly the kind of "accepted but disappeared" failure that must be ruled out before trusting the packet lifecycle.

The audit must determine whether they:

```text
were not persisted
or
were persisted but not projected
or
were outside the dashboard scope
or
were only narrative text and never submitted
```

That is not cosmetic. It is ingestion-integrity.

### Declared identity

Agreed: acceptable for the POC.

Treat:

```text
agent_id / harness / model
```

as **declared provenance**, not cryptographically verified runtime identity.

If later you add attestation, it should be additive:

```text
declared identity
+
verification/attestation metadata
```

not a replacement of the declared fields.

### Snapshot consistency

Keep this deferred, but record it as a named limitation:

```text
adversarial packet reviewed context snapshot Sx
```

is eventually useful for replayability. No implementation now.

---

# Revised audit prompt

I would now use this as the next coding-agent audit. It incorporates the valid points from both reviews without prematurely turning them into features.

```text
# ARGUS Audit — Research Provenance, Adversarial Targeting, and Submission Integrity

DO NOT MODIFY CODE.

This is an audit only.

The first Work + Adversarial dry run has now completed.

We have verified:

- fresh Work Agent reconstruction through `argus.get_context`
- all current scenario beliefs are visible through `get_context`
- evidence, intents, and belief edges are visible
- beliefs are immutable
- adversarial work creates new belief IDs
- adversarial challenges can target existing canonical belief IDs
- agent identity is now declared in submitted packets
- packet-level agent provenance is persisted
- `packet_submission` exists as append-only provenance

The next audit is to determine whether the research lineage is complete and
whether the adversarial workflow has any silent integrity gaps.

Do not redesign the system.
Do not add MCP tools.
Do not implement fixes yet.

==================================================
1. BELIEF / OBJECT PROVENANCE
==================================================

Determine whether each persisted research object can be traced to the exact
packet that created it.

Audit:

- beliefs
- evidence
- belief edges
- tasks created through packets

For each object type answer:

1. Is `origin_packet_id` or an equivalent persistent relationship stored?
2. If not, is the relationship only transient inside packet compilation?
3. Can the relationship be reconstructed deterministically without relying
   on timestamps, ordering, agent names, hashes, or inference?
4. Can a human trace:

       object → packet → agent identity

   reliably?

5. What is the smallest persistence change needed, if any?

Do NOT assume that an in-memory local→canonical mapping constitutes persistent
provenance.

Produce a decision table:

| Object | Current origin link | Persistent? | Exact location | Recommended location |
|--------|---------------------|-------------|----------------|----------------------|
| belief | ... | ... | ... | ... |
| evidence | ... | ... | ... | ... |
| edge | ... | ... | ... | ... |
| task | ... | ... | ... | ... |

The recommendation must distinguish:

- kernel-adjacent schema change
- additive provenance table
- coordinator-side relationship
- no change required

Do not choose an architecture until the actual code is inspected.

==================================================
2. PACKET SUBMISSION IMMUTABILITY
==================================================

Verify:

- one packet_id maps to one packet_submission record
- packet_submission cannot be silently updated
- packet_submission cannot be deleted through normal application paths
- retry/idempotency does not create duplicate provenance

Also determine whether any normal command can delete research objects.

Distinguish:

- normal research lifecycle
- development reset/destructive commands

Desired normal lifecycle:

    objects are not silently deleted
    they are preserved and may be superseded/retracted

Do not implement deletion protections unless the audit finds an actual
normal-path deletion risk.

==================================================
3. ADVERSARIAL TARGET VALIDATION
==================================================

Audit what the current packet compiler/validator permits for `contradicts`
edges in an adversarial packet.

Determine:

A. Can an adversarial packet target:
   - a canonical existing belief?
   - a local belief created in the same packet?
   - a belief in another scenario?
   - a nonexistent belief?
   - a retracted belief?
   - a promoted belief?

B. Which of these are currently rejected?

C. Which should be rejected based on the existing ARGUS architecture?

Important:

Do NOT assume that only unresolved/open beliefs may be challenged.

A promoted belief may legitimately be challenged as part of later research.
A historical/retracted belief may still be analyzed.

The minimum integrity requirement to assess is:

> An adversarial contradiction should resolve to an existing canonical belief
> in the same scenario, not merely to a new local object inside the same
> submission.

Determine whether the current implementation already guarantees this.

Also determine whether an adversarial packet can create an apparent
self-challenge through local references.

Do not implement a new rule until the audit establishes an actual gap.

==================================================
4. ADVERSARIAL REVIEW COVERAGE
==================================================

The adversarial agent currently sees all beliefs in the scenario.

Determine whether the protocol or packet requires the adversarial agent to
record a disposition for every visible belief.

The desired current behavior is:

    B1 — challenged
    B2 — no material objection
    B3 — insufficient basis
    B4 — challenged
    ...

This is a PROCESS requirement, not necessarily a persistence feature.

Determine:

- Is this already in the current protocol?
- Is it in the adversarial prompt?
- Is it absent?

Do NOT propose a review-coverage table yet.

If absent, recommend only the minimum protocol change required to make
coverage auditable in the adversarial agent's final report.

==================================================
5. EBP CATEGORY-ERROR CALIBRATION
==================================================

Audit the shared agent protocol for this distinction:

    unresolved debt
        ≠
    EBP violation

The first adversarial run incorrectly characterized:

    entered belief + 8 open debts

as an EBP violation.

Determine whether the shared protocol currently prevents this category error.

If not, recommend the smallest addition:

> Do not label an unresolved obligation an EBP violation unless a specific
> active EBP or Domain Pack rule is actually violated.

This applies to BOTH Work and Adversarial agents.

Do not implement code for this; this is primarily a protocol/documentation
audit.

==================================================
6. TASK SUBMISSION INTEGRITY
==================================================

The adversarial dry run reported proposed downstream tasks.

Verify whether those tasks actually exist in ARGUS state.

For each proposed task determine:

- Was it included in the submitted EBP packet?
- Was it accepted by validation?
- Was it persisted?
- If persisted, what is its ID?
- Does it have packet provenance?
- Is it visible through get_context?
- Is it visible in the Trust UI?
- Is it excluded because of dashboard scope?
- Was it merely stated in the agent's prose and never submitted?

This is a priority finding.

A submitted task that is silently accepted but invisible is an
ingestion/projection integrity failure.

==================================================
7. HUMAN TRACEABILITY
==================================================

Determine whether a human looking at the current Trust UI can answer:

For any belief:

1. What is the claim?
2. Who/what agent created it?
3. Which packet created it?
4. What evidence accompanied it?
5. What other beliefs does it derive from?
6. What beliefs contradict it?
7. Which adversarial packet challenged it?
8. What tasks resulted from the work?

If the answer is NO, identify exactly which links are missing.

Do not design a graph UI.

Prefer simple relational projections/tables.

==================================================
8. DECLARED VS VERIFIED AGENT IDENTITY
==================================================

Confirm that:

agent_id
harness
model

are currently declared provenance supplied by the operator/launch prompt.

Determine whether any code implies these values are cryptographically or
runtime verified.

If not verified, state:

> Agent identity is declared provenance, not attested runtime identity.

Do NOT add attestation.

==================================================
9. CONTEXT SNAPSHOT / REPLAYABILITY
==================================================

Determine whether `argus.get_context` or submitted packets record a stable
context snapshot/version/hash.

Do not implement it.

Classify:

PASS
PARTIAL
MISSING

Explain whether the absence affects the current single-agent dry run.

Record it as deferred research-system debt if useful.

==================================================
10. FINAL ASSESSMENT
==================================================

Produce:

### A. Provenance spine

Draw the actual current graph:

    packet_submission
          |
          +--> belief?
          +--> evidence?
          +--> edge?
          +--> task?

Mark every missing relationship.

### B. Adversarial integrity

State whether the current system prevents:

- attacking nonexistent beliefs
- cross-scenario targeting
- same-packet self-challenge
- ambiguous target references

Do not invent policy where the current architecture has not defined one.

### C. Task integrity

Explain exactly what happened to the two proposed adversarial tasks.

### D. Protocol quality

Identify any remaining behavioral ambiguity in:

- shared protocol
- Work prompt
- Adversarial prompt

### E. UI lineage

State exactly what a human can and cannot currently trace.

### F. Migration decision

For each missing provenance relationship, recommend:

KEEP AS-IS
SMALL ADDITIVE CHANGE
DEFER

Prefer the smallest solution.

No code changes.
No new MCP tools.
No new authority layer.
No architecture expansion.
```

### Bottom line

I would run **this audit before making another implementation change**.

The dry run has moved ARGUS into a new stage: the basic mechanism works, and now you're testing whether the system can preserve a **credible chain of provenance from agent → packet → claim → challenge → successor**.

That provenance spine is more important than adding more research features at this point.
