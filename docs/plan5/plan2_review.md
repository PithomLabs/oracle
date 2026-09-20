## prompt

consolidate reviews below where valid, we need to be disciplined following the less is more design constraint (design is final until there is nothing absolutely essential to subtract) 

## claude

**Two things plan2.md gets right that are worth crediting specifically, not just approving in passing.**

Section 5's cascade deferral is exactly the decision I flagged as needed two turns ago — EBP has no `belief_edge` equivalent, and routing "T2 depended on T1" through an ordinary Conductor re-review task rather than building an EBP-native cascade is the correct, evidence-gated call. Section 6 is a genuine improvement *on my own* prior proposal: I had suggested one Skill per debt item (`map-builder`, `toy-checker`, `null-model-hunter`); plan2.md correctly rejects that as premature specialization and consolidates to a single `ebp-bmist` Skill that chooses its own tools internally. That's the same discipline that killed workflow.md's five-primitive framework several turns back, applied here to catch a mistake in my own design rather than someone else's — worth naming as a real correction, not just a stylistic preference.

**The substantive gap: Beads' `Gate` primitive is never mentioned, and it's arguably the single most relevant Beads feature to this specific pilot.** Section 9's "what to keep" list is READY frontier, atomic claim, parent/child, discovered-from — but beads.md's own workflow section describes a `Gate` as parking a step "until something external happens: a human sign-off, a timer, or a GitHub run or PR." That's not a loose analogy to what section 7 is trying to build for `needFaithfulnessReview` — it's the same mechanism, already specified, already native to the system plan2.md is auditing. As written, the plan builds the faithfulness checkpoint as a bespoke Pizza Bot Action/HITL wire-up, entirely separate from whatever the Conductor audit (item 1) turns up — but if Conductor ends up adopting Beads-style gates for other reasons, the faithfulness gate should very plausibly route through that same primitive rather than being a parallel, Pizza-Bot-only mechanism. This is worth adding to the audit's scope explicitly, and worth sequencing the audit *before* finalizing item 7's design rather than treating them as independent mandatory items.

**A second Beads omission, and it's not cosmetic — it would eliminate hand-written work the plan currently specifies manually.** Section 8's task list ("retire needMap, retire needInvariant, run toy check, investigate obstruction, perform faithfulness review") is, almost verbatim, the shape of a Beads *formula*: a declared-once TOML DAG of steps, cooked into a proto, poured into a molecule of real, dependency-linked beads on demand. EBP's own canonical six-item debt pipeline is exactly the kind of repeatable structure formulas exist for. Plan2.md hand-specifies this task list per idea rather than noticing it's describing a template Beads already has a name for. Worth folding "Workflows / formula-proto-molecule" into the audit's feature list alongside READY and discovered-from — it's not a "later" feature, it's the thing that would make item 9 (seeding one BM-IST claim) reusable for the second and third claims without re-deriving the task breakdown each time.

**An ambiguity plan2.md inherits from plan.md without resolving: "borrow the concepts" versus "adopt the tool."** The original question was "are the Beads ideas already integrated into Conductor?" — but there are two very different projects hiding behind "integrate Beads": reimplementing READY-frontier logic and discovered-from edges natively in Conductor's existing SQLite schema, versus actually running Beads/Dolt as Conductor's backing store or as a peer system. These have wildly different costs — the first is a modest query and an edge-type addition; the second means adopting Dolt's git-based sync model as infrastructure. Section 9's audit template (PRESENT/PARTIAL/MISSING per feature, checked against Conductor) reads as though the former is assumed, but nothing in plan2.md says so explicitly, and the user's own framing of the question left it open. Worth stating outright: the audit should be scoped to "does Conductor need to grow these concepts natively," not "should Conductor be replaced or backed by Beads" — and if that's not already the shared understanding, it needs to become one before the coding agent starts, since the two audits would look for entirely different things.

**A related procedural gap: what governs whether an audit finding is allowed to become a change?** Section 14's hard-stop list includes "no Conductor redesign," but if the audit comes back "READY frontier: MISSING," implementing it *is* a Conductor change — the plan doesn't say what threshold separates a permitted, evidence-justified addition from a forbidden redesign. Conductor's own AGENTS.md already has a Growth Gate for exactly this question (could this be an adapter or service instead of a core primitive?). The plan should say explicitly that any audit-recommended Conductor change routes through that existing gate rather than being pre-approved by "the audit found a gap" or pre-forbidden by the hard-stop list — right now both readings are consistent with what's written, which isn't a place this plan should leave ambiguous given how much weight the rest of the ecosystem puts on that exact gate.

**A structural risk in the target architecture diagram itself (section 1) — two parallel writes where the correspondence between them is only conventional.** The diagram shows the Pizza Bot Agent writing to the EBP Ledger and to Conductor as two independent arrows. Section 8's later prose implies something stronger — a Conductor task like "retire needMap for T2" that corresponds to a specific, live EBP debt item — but nothing in the plan enforces that correspondence structurally. If a Conductor task gets marked accepted without the corresponding `ebp_retire` call actually happening, the two stores silently disagree, and nothing catches it — precisely the "evidence has an owner" failure the Reference Loop document spent an entire section warning against, now reappearing one layer up in a system built partly in response to that document. The concrete fix: Conductor task metadata should carry the EBP idea-id and debt-item name as a structured reference, not just a similarly-worded title, and ideally the "submit/accept" transition on the Conductor side should verify (or trigger) the actual `ebp_retire` call rather than trusting the two to stay in sync by convention.

**One detail in the seed-state item worth getting right before it's run.** Item 9 says "seed one real BM-IST Idea/Claim and its current EBP debt" — but Theorem 1 (the candidate named in section 10) isn't a fresh idea with all six debts open. Prior rounds in this investigation already produced a verified derivation, found the Hall-Reginatto precedent (folding into `needMap`'s retirement), and left `needObstruction` explicitly open (nobody has attempted the $S$-dependent construction Z's own writeup flagged as the scope boundary). The seed step should encode that actual state, not reset Theorem 1 to a naive full-debt starting point — otherwise the first real run of this pilot would silently discard work already done, which is a strange way to open a demonstration whose whole point is showing the system coordinates real research state faithfully.


## Z

Evaluation against the frozen design and the conversation record. One verification boundary first: the EBP 2.1 specification and Pizza Bot's internals (DeepAgents/LangGraph runtime, Action queue) have never appeared in this conversation — I'm evaluating the summary as given and flagging what rests on it. With that bound:

**Verdict: the architecture mapping is sound and this POC is the right next move — it's the first plan that exercises every role in the frozen design, including the two never yet exercised (Domain, real SOR). But it contains one false premise that invalidates its sequencing, and the 10-item essential list omits the gates that stand directly between here and that premise.** Fixable; nothing architectural.

## Critical

**C1 — "Reuses the real Phase 2 path we already proved" is false against this record.** No real external effect has ever been proven in this conversation: Phase 1 ended with a RecordingFunc (fake) executor and an honestly-scoped PASS; the Phase 1.5 battery exists as a *plan* — no results have been produced; Phase 2 has never run. Yet §11 and item 8 describe the POC's consequential publication (step 12–13: Solvent → Executor → GitHub SOR) as reuse of a proven path, and the essential list contains no Phase 1.5 execution, no Phase 2 pre-registration, no capability declaration. As written, this plan smuggles the project's first real external effect in as step 12 of an implementation checklist — bypassing the frozen design's own §27 gate. Either Phase 2 has since run (then produce the report) or it hasn't (then correct the claim). The correct framing: **the POC's step 12 *is* the Phase 2 experiment** and inherits §27's preconditions in full.

**C2 — The §27 preconditions are all unmet and mostly unlisted.** Before the first real effect, the frozen design requires: pinned declaration (version + hash) — *no capability declaration exists for any operation; "first real capability declarations" is §23's open item, and a publication operation needs its own declaration including operation-identity definition, trust basis, validity model, replay behavior, outcome model*; pre-registered pass/fail/falsification criteria — *standing open item, never discharged*; verified effect-enforcement path — *Phase 1.5's mismatch rejection (Test A) is precisely the safety predicate for a real publish, and it hasn't run*; verified external SOR — *the GitHub repo must be named and pinned*. "No new Solvent work" is correct — the kernel stays frozen — but the declaration is an *integration artifact*, not Solvent work, and it is missing from the essential list entirely.

## High

**H1 — The EBP↔Conductor seam has no ownership fence, and it creates a dual-record drift risk.** The mapping (Debt ↔ work item, retirement ↔ work + evidence) leaves unanswered: when a debt item's epistemic status and its Conductor task lifecycle disagree, which is authoritative? This is the history-fence problem in a new costume. The rule consistent with the design: **EBP owns epistemic status (debt exists, promotion derivation); Conductor owns coordination (claims, lifecycle, acceptance records); retirement occurs only after Conductor-recorded acceptance, with EBP storing a reference to that decision — never re-deriving it.** One paragraph; without it, two systems track one work item and drift is guaranteed by Week 3.

**H2 — The essential list has no branch for the audit outcome.** Item 1 audits Conductor for READY/atomic-claim/provenance and forbids modification — correct — but the POC workflow *depends on the READY frontier*. If the audit returns MISSING, minimal remediation becomes mandatory before item 5 (discover READY work) can function. The list treats the audit as standalone; it's a gate with two outcomes, one of which changes the plan.

**H3 — Approval and acceptance decisions must land in Conductor's durable record.** §7 routes the faithfulness gate through Pizza Bot's Action queue — a natural fit, but if the human approves in Pizza Bot and Conductor never records the decision, the durable history that lets Agent X resume has a hole exactly at the gate that matters most. Pizza Bot's queue can remain the *mechanism*; the *decision record* (attribution, plan/task linkage) must be written to Conductor. Same for the plan approval itself (§5).

## Medium

- **The POC's role map is actually a strength — make it explicit.** Pizza Bot=Agent, EBP ledger+faithfulness gate=**Domain** (promotion *is* domain acceptance), GitHub=SOR. That's the first plan where Domain and real SOR have owners. State it in the POC doc; it's also what makes this a genuine domain-agnosticism test rather than a demo.
- **Frozen-Solvent conformance delta must be recorded as known deviation.** Test D's documentation already concedes UNKNOWN=DENIED collapse and no age policy. The POC will run on a kernel that doesn't implement the design's authorization taxonomy — record it in the delta table (owed since the Phase 1.5 review) so the POC report can't be misread as conformance evidence.
- **Same for Conductor's declaration gate** (§8.1 not implemented — Phase 1.5 plan admitted it). Environment delta, or a minimal resolution check becomes an essential item.
- **Plan scope must include the publication operation class.** The plan the human approves should carry `PublishArtifact` (with target scope) at class grain, and the consequential step must record plan linkage — otherwise the POC violates its own frozen §5.1 on its flagship operation.
- **General debt work needs acceptance authority, not just faithfulness debt.** §7 covers `needFaithfulnessReview`; the other work items (retire needMap, toy check) also need `verification_ref` pointing at their decider (tests, checker, human) — §14.2 applies to all SUBMITs.
- **Make domain-agnosticism inspectable.** "Conceptually substitute the domain" is a narrative test. The checkable property: zero domain-specific fields in the Conductor schema, plus a protocol-vocabulary diff after the second domain arrives. One sentence in the POC's success criteria.
- **The two subtractions (no per-debt Skills, no cascade) are correct** — both match recorded deferrals, and the Agent-mediated re-review task pattern for dependencies is the right containment.

## Direct answers

**Beads integration status:** agree — unknown, and your instinct to have the coding agent investigate is right. The audit-first, no-changes-during-audit approach is exactly correct; just add H2's branch so the audit's output connects to a decision.

**Is the consolidation sound?** Yes, with C1's correction. The sequencing error is the consolidation's only structural flaw; everything else survives.

## Corrected essential list (delta to the ten)

Insert before item 8, renumbering accordingly:

```text
8a. Run Phase 1.5 battery (or attach its results) — mismatch
    rejection is the safety predicate for the real publish
8b. Capability declaration for PublishArtifact: owner, version,
    hash, class, identity definition, trust basis, validity,
    replay, outcome model — integration artifact, mandatory
8c. Phase 2 pre-registration: pass/fail/falsification criteria,
    pinned components, pinned declaration, named/pinned GitHub SOR,
    verified enforcement path — BEFORE step 12
8d. EBP↔Conductor seam fence: EBP owns epistemic status,
    Conductor owns coordination, retirement follows recorded
    acceptance, references not duplication
8e. Decisions land in Conductor: plan approval and faithfulness
    gate decisions recorded with attribution, whatever queue
    presents them
8f. Audit branch: if READY/claim gaps are MISSING, minimal
    remediation becomes mandatory before POC start
8g. Deviation table: frozen-Solvent taxonomy delta + Conductor
    declaration gate, recorded as environment deltas
```

With those seven inserted, the plan becomes what it wants to be: the first workload that exercises the complete architecture — every role, both planes, one real consequence, one new domain — under the discipline the design froze precisely so that this moment would be boring. It's one honest sentence away ("this will be the *first* real external effect, gated accordingly") and one declaration away from being executable.


## deepseek

## Evaluation of the plan

**Overall verdict:** The plan is directionally strong and unusually disciplined. It correctly separates concerns, defers premature machinery, and makes the Conductor audit the first task. But it is **not yet implementation-ready**. The mandatory list is mostly right, yet it is missing a few absolutely necessary contract-first items. Without those, the POC risks building integrations before the interfaces and state ownership are pinned down.

---

## What the plan gets right

- **Clean separation of roles:** EBP = epistemic grammar, Conductor = coordination/state, Solvent = authority, Executor = effect, Pizza Bot = agent. This is the right domain-agnostic skeleton.
- **EBP stays out of Conductor core.** Correct. Conductor should not understand `needMap` or any domain debt semantics.
- **Minimal EBP surface:** only `capture`, `status`, `retire`, `add-debt`, `promote`. No cascade, no scoring, no dashboard.
- **One Pizza Bot skill, not one per debt type.** Correct for POC.
- **Human faithfulness gate is mandatory and uses existing Pizza Bot Action/HITL.** Correct and important.
- **Solvent only for consequential external action.** Correct. The synthesis itself should not cross Solvent.
- **Beads handled honestly:** it does not assume `READY`, `discovered-from`, or parent/child are already in Conductor. The audit-first instruction is exactly right.
- **Hard stop rules are good.** They prevent scope creep into Conductor redesign, EBP cascade, new Solvent/Executor, etc.

---

## Critical gaps before starting

These are the things I would add to the “absolutely necessary” list.

### 1. Conductor audit must produce a feature/contract matrix, not just a checklist

The audit should report:

- Does Conductor already expose an MCP or API surface?
- What are the exact task/project/claim/dependency/history schemas?
- Is `READY` derived or stored?
- What is the claim/release semantics, including lease/TTL if any?
- What provenance types exist: parent/child, `discovered-from`, blocked-by?
- What durable evidence/artifact attachment exists?
- What is the source of truth for task state?

If Conductor MCP does not already exist, mandatory item #5 is not an integration task — it is a build task. That must be explicit.

### 2. Contract-first schemas and source-of-truth matrix

Before code, define:

- EBP `Idea` and `DebtItem` schemas.
- Conductor work item schema.
- The exact mapping table between them.
- Which system owns which field.
- What happens when EBP and Conductor disagree.

The plan says “EBP Debt → Conductor work,” but not who creates it, who updates it, or how retirement evidence flows back.

### 3. Lifecycle/state mapping

You need a precise mapping:

- EBP debt open ↔ Conductor task state?
- EBP debt retired ↔ Conductor task completed + evidence?
- EBP promotion ↔ Conductor artifact accepted?
- Obstruction ↔ blocked/blocked-by?

Without this, the Agent will improvise and the POC will not be reproducible.

### 4. Concurrency, idempotency, and claim safety

Even if you defer TTL/scheduler, you still need:

- atomic claim semantics,
- idempotent MCP calls,
- retry behavior,
- conflict handling if two agents or humans act at once.

This is mandatory before any real workflow touches the systems.

### 5. Human faithfulness gate contract

The plan says use Pizza Bot’s Action queue. That is right, but you must define:

- What exactly is proposed to the human?
- What evidence is attached?
- What are the approve/reject outcomes?
- How is the decision recorded in EBP and Conductor?
- Who is allowed to approve?

Otherwise the gate is a UI gesture, not a control.

### 6. Solvent/Executor/SOR path confirmation

The plan says reuse the already-proven path. Before starting, confirm:

- Solvent is still available and scoped for this.
- The existing Executor can perform the GitHub SOR action.
- Credentials and rollback/failure behavior are defined.
- The “one consequential publication” is actually a safe, reversible operation.

If any of that is not true, it becomes a mandatory pre-req.

### 7. BM-IST seed claim vetting

The plan correctly conditions this on the artifact being “sufficiently characterized.” That condition should be a hard gate:

- If the linear Madelung ⇒ Fisher structure claim is not already well-scoped with clear debt, choose another seed claim.
- Do not ingest the whole corpus.

### 8. POC acceptance criteria

Define success explicitly, e.g.:

- One EBP idea captured.
- One debt attached.
- One Conductor task created, claimed, worked, submitted.
- Human faithfulness decision recorded.
- EBP debt retired or explicitly left open.
- One promotion or explicit non-promotion.
- One Solvent-gated external action.
- Full history traceable across EBP + Conductor.

Without this, “POC works” is subjective.

---

## One design concern

The skill is named `ebp-bmist`. That is domain-specific. If the goal is to prove the workflow is domain-agnostic, I would rename it to something like `ebp-research` and treat BM-IST as configuration/content, not as the skill identity. Otherwise the “domain-agnostic” test is weaker than it looks.

Also, the domain-agnostic test is currently conceptual: “substitute BM-IST with another domain.” To actually prove it, you need either a second tiny dry-run domain through the same Conductor protocol, or at least a tabletop substitution exercise with the same skill and a different EBP content pack. No new Conductor features.

---

## Bottom line

The plan is a solid POC skeleton. It correctly makes the Beads/Conductor audit first, defers cascade and specialized skills, and keeps Solvent as the only consequential authority. But before implementation, add a **contract-first phase**: schemas, mapping, source-of-truth, lifecycle states, idempotency/claim safety, human-gate contract, Solvent path confirmation, seed-claim gate, and explicit acceptance criteria.

On Beads: I cannot confirm from the provided text whether `READY`, `discovered-from`, or parent/child are already integrated into Conductor. The plan is right not to assume. Let the coding agent audit first.
