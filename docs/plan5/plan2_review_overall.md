Yes. The reviews converge strongly. Applying the **less-is-more constraint as a hard design law**, several suggestions should be accepted, while several others should be explicitly rejected as premature machinery.

## Consolidated decision

### 1. KEEP — Conductor audit first

**Mandatory. No code changes during audit.**

Audit specifically for:

| Capability                    | Audit                       |
| ----------------------------- | --------------------------- |
| Derived `READY` frontier      | present / partial / missing |
| Atomic claim                  | present / partial / missing |
| Parent / child provenance     | present / partial / missing |
| `discovered-from` provenance  | present / partial / missing |
| Durable project state/history | present / partial / missing |
| Existing MCP/API surface      | present / partial / missing |
| Evidence/artifact references  | present / partial / missing |
| Current source of truth       | explicit                    |

This is stronger than the earlier audit because it also checks the actual integration surface. 

**Important:** this is a **Conductor capability audit**, not a Beads adoption project.

Do **not** introduce Beads/Dolt as infrastructure. We are evaluating whether useful Beads concepts need to exist in Conductor's existing architecture. 

---

### 2. KEEP — Audit outcome becomes a gate

The plan previously had a small ambiguity:

> audit says something is missing → then what?

Resolve it explicitly:

```text
Audit
  ├─ sufficient → continue
  └─ missing essential capability
         ↓
   Conductor Growth Gate
         ├─ adapter/service sufficient → use that
         └─ core change genuinely required → minimal change
```

A missing Beads-style feature does **not** automatically justify changing Conductor. The existing Growth Gate decides that. 

This preserves the freeze discipline.

---

### 3. KEEP — Contract/source-of-truth fence

This is genuinely necessary.

We need one compact ownership table before implementation:

| Concern                                         | Owner            |
| ----------------------------------------------- | ---------------- |
| Idea / debt / promotion eligibility             | **EBP**          |
| Work / claim / lifecycle / coordination history | **Conductor**    |
| Consequential authorization                     | **Solvent**      |
| External effect                                 | **Executor**     |
| Whether effect actually occurred                | **External SOR** |
| Agent context / conversation                    | **Pizza Bot**    |

Most importantly:

> EBP and Conductor must **reference each other**, not duplicate each other's authority.

EBP owns epistemic state.
Conductor owns coordination state. 

---

### 4. KEEP — Minimal EBP 2.1 implementation

Exactly the five canonical operations:

```text
capture
status
retire
add-debt
promote
```

No cascade engine.
No scoring.
No truth engine.
No workflow engine.

That remains one of the strongest parts of the plan. 

---

### 5. KEEP — One EBP↔Conductor mapping

The mapping should be structural, not merely textual.

For each EBP debt item, the corresponding Conductor work item should carry a **reference** such as:

```text
ebp_idea_id
ebp_debt_item
```

Not copied semantic state.

Likewise, EBP retirement/promotion should retain the relevant Conductor decision/history reference.

This prevents the dual-store drift problem identified in the review. 

---

### 6. KEEP — Human faithfulness gate, but use existing Pizza Bot HITL

Do not build a new approval subsystem.

Pizza Bot supplies the interaction mechanism.

But the **decision itself must become durable Conductor history**:

```text
proposal
→ Pizza Bot Action/HITL
→ human decision
→ Conductor records decision
→ EBP retirement proceeds
```

Same principle for plan approval.

The queue is the **mechanism**.
The decision record is the **durable fact**. 

---

### 7. KEEP — One Skill only

The reviewer's earlier proposal for one Skill per debt type is correctly rejected.

Use:

```text
ebp-bmist
```

for the POC.

It teaches the Agent the EBP vocabulary and integration patterns; it does not encode separate orchestration systems for each debt.

The consolidation from multiple specialized Skills to one Skill is a genuine simplification. 

I would **not rename it yet** to `ebp-research`. That is unnecessary churn before the POC. The domain-agnostic property must come from the underlying protocol and Conductor schema, not from a cosmetic Skill name.

---

### 8. REJECT FOR NOW — Beads `Gate`

Interesting, but **not essential yet**.

A human faithfulness checkpoint is already required, but introducing a generalized Beads `Gate` primitive into Conductor merely because Beads has one risks turning one concrete requirement into a new generic primitive.

Therefore:

```text
Faithfulness HITL = POC requirement
Generic Conductor Gate primitive = deferred
```

The Conductor audit may note whether such a concept already exists, but **do not add it merely for this POC**. The review correctly identifies its conceptual relevance, but relevance is not enough to overcome the less-is-more rule. 

---

### 9. REJECT FOR NOW — Beads formulas / proto / molecules

This is the clearest overreach.

Yes, the EBP debt pipeline resembles a reusable formula.

No, that does **not** mean Conductor should acquire a formula/molecule workflow engine.

For this POC:

```text
EBP debt
→ one or a few normal Conductor work items
```

Repeatable work can be hand-created initially.

Only recurring operational pain can justify a future abstraction.

Thus the Beads formula concept belongs in the audit's **observational findings**, not the implementation baseline. 

---

### 10. KEEP — Atomic claim / idempotency

These are not optional Beads decoration. They protect actual correctness.

Minimum requirements:

```text
atomic claim
idempotent mutation semantics
safe retry
conflict rejection
```

Do **not** add claim TTL/heartbeat yet.

The review is right that claim safety matters now; it is wrong to infer that the full lease machinery must therefore be built now. 

---

### 11. KEEP — Existing Phase 2 gate

There is one important correction from the review that must be incorporated into the consolidated plan:

**Do not describe the consequential path as merely “reuse the proven path” unless the actual Phase 2 evidence is attached.**

The mandatory precondition for a real external effect remains:

```text
capability declaration
+ exact operation identity
+ pinned declaration/hash
+ pre-registered pass/fail criteria
+ verified enforcement
+ pinned external SOR
```

The review correctly insists that this cannot be silently skipped. 

Given our latest project state, the Phase 2 real GitHub execution **has in fact been completed**, so the consolidated plan should reference that actual evidence rather than treating it as hypothetical.

Therefore the POC should **reuse the proven mechanism**, not rebuild it.

---

### 12. KEEP — Capability declaration for the POC consequence

The publication operation needs its own declaration artifact.

At minimum:

```text
owner
version
content hash
operation class
exact identity definition
target scope
trust basis
validity
replay/idempotency
outcome model
```

This is an **integration artifact**, not a Solvent-kernel modification.

Correct to keep. 

---

### 13. KEEP — POC seed must reflect real BM-IST state

Do not fabricate a fresh six-debt Idea if the research already has partially retired debt.

The seed should represent the **actual current research state**.

That matters because this is supposed to demonstrate faithful coordination of real research work, not a synthetic workflow toy. 

---

### 14. KEEP — Domain-agnosticism as an architectural property

Do not add another domain just to make a slide say “domain agnostic.”

The stronger test is:

```text
Conductor schema contains zero BM-IST concepts
EBP protocol contains zero BM-IST concepts
Solvent remains generic
Only Skill/content/seed data are BM-IST-specific
```

A second-domain tabletop test can happen later, but it is **not required before beginning the first POC**.

The architecture must be generic; the first workload is allowed to be specific.

---

## Final pre-implementation baseline

After consolidation, the **absolutely necessary** work is:

### A. Conductor

1. Audit existing Conductor against:

   * READY
   * atomic claim
   * parent/child
   * discovered-from
   * durable history
   * MCP/API
   * evidence references
2. No implementation until audit completes.
3. Missing essential capability → existing Growth Gate.
4. Only minimal justified additions.

### B. Contracts

5. Define one source-of-truth/ownership matrix.
6. Define EBP↔Conductor reference mapping.
7. Define minimal lifecycle/mutation/idempotency semantics.

### C. EBP

8. Implement only:

   * capture
   * status
   * retire
   * add-debt
   * promote
9. Keep EBP semantics exactly aligned with v2.1.

### D. Pizza Bot

10. One `ebp-bmist` Skill.
11. Connect existing Conductor MCP.
12. Connect EBP MCP.
13. Use native Pizza Bot HITL for faithfulness review.

### E. Governance / consequence

14. Human plan approval recorded durably in Conductor.
15. Human faithfulness decision recorded durably in Conductor.
16. Use the already-proven Solvent → Executor → GitHub mechanism.
17. Provide the publication capability declaration and pinned operation identity.
18. Preserve the frozen Solvent kernel unchanged.

### F. POC

19. Seed **one real BM-IST Idea with its actual existing debt state**.
20. Drive one debt item through real work.
21. Exercise human faithfulness review.
22. Retire/promote according to actual EBP state.
23. Perform one consequential external action through the existing Phase 2 path.
24. Produce one end-to-end trace across:

```text
Pizza Bot
   ↓
EBP
   ↓
Conductor
   ↓
Solvent
   ↓
Executor
   ↓
External SOR
```

---

## Explicitly DEFERRED

These do **not** belong in the starting implementation:

```text
Beads/Dolt adoption
generic Gate primitive
formulas / proto / molecules
EBP cascade
belief graph
scheduler
claim TTL / heartbeat
per-debt Skills
Conductor policy engine
new Solvent features
new Executor
domain-specific Conductor fields
second-domain implementation
```

The governing rule is:

> **A missing abstraction is not an implementation requirement merely because another system has the abstraction. It becomes a requirement only when the current POC cannot work correctly without it.**

That is the cleanest synthesis of all three reviews and is consistent with the project's existing discipline: **freeze the architecture, prove the smallest real loop, and only add machinery when evidence makes subtraction impossible.**
