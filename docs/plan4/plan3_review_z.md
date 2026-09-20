Consistency note: this consolidates my last two answers, so I've re-attacked from zero rather than grading my own homework. The result: the two headline invariants survive attack, but the sentence you want to *lock* contains a definition that regresses the hardest-won finding of the entire specification lineage — and the consolidation silently dropped four routed fixes that the new schema cannot be built without.

## Disposition of prior conditions and routed items

| Prior item | Status | Notes |
|---|---|---|
| Approval ≠ authority (condition 1) | ✓ Locked | §2 states it as an invariant. Correct and load-bearing. |
| Plan enumerates operations (condition 2) | ✗ Rejected at the wrong grain | §7 declines *exact* enumeration — defensible — but misses the middle position. See H1. |
| Conductor stores plans, never gates (condition 3) | ✓ Landed | plan_id/version/decision recorded; no plan-aware routing. |
| Declaration resolution path (schema review H1) | **Silently dropped** | §3 retains `capability_ref`; nobody serves what it points to. |
| Acceptance ownership (H2) | **Silently dropped** | §9's verb list has no ACCEPT/REJECT owner. |
| History fences (H3) | **Silently dropped** | §9 adds "record history" with no shadow-ledger fence. |
| Attribution + task-text-is-data (H4/M6) | **Silently dropped** | Multi-agent §11 makes this worse, not better. |

## Critical

**C1 — The §4 definition regressions: "intentionally" and "authoritative state."** The matrix's settled rule: *"Any operation **capable of** producing an external effect is consequential regardless of workflow phase or invoking client."* The consolidation's proposed lock: *"anything whose successful completion **intentionally** mutates **authoritative state** outside the workflow."* Both substitutions reopen classification at the definitional level after it was closed mechanically:

- **"Intentionally"** makes consequence a mental property of the Agent again. The agent whose shell command has unintended external effects hasn't "intentionally mutated" anything — and an agent *claiming* lack of intent is unfalsifiable. This is the advisory-classification bug, reintroduced one level up from where it was killed. Classification must remain a property of the *capability declaration* (mechanical, at the boundary), never of agent purpose.
- **"Authoritative state"** silently narrows "external effect." Is a posted webhook, an email, a chat message, a third-party API write "authoritative state"? Ambiguous — and ambiguity in the definition of the thing Solvent gates is where bypasses are argued into existence later.

Fix is one sentence, restoring the lineage: *"An operation is consequential when it is declared capable of producing an external effect — an observable state change outside the workflow — regardless of intent, phase, or invoking client."* Declarations classify; Solvent gates; intent is irrelevant by construction. **Do not lock the final architecture sentence until this is repaired**, because locking a defective definition is how this project's regressions become permanent.

**C2 — The consolidation dropped the punch list that the schema depends on.** This document presents itself as the consolidation of both reviews, and four High findings with routed fixes vanished without rejection reasons — the fifth occurrence of the silent-carry pattern, now blocking rather than cosmetic:

1. `capability_ref` **cannot resolve** — no declaration resolution path, no fail-closed at CREATE. Tasks referencing unresolvable capabilities are creatable under this doc as written.
2. **Nobody accepts or rejects** — task submission has no deciding authority, so DONE is agent-asserted and the READY frontier inherits unverified completion.
3. **History has no fences** — "durable history" adjacent to the anti-ledger invariant, with no coordination-facts-only rule and no "decided ≠ authorized" line.
4. **Arbitrary agents joining** (§11) with no attribution, no creation policy, and agent-authored task text that future agents will read as instructions.

Restoring these isn't scope creep — each is a one-or-two-sentence rule already drafted in the prior reviews. But it's now undeniable: **the disposition log must exist before the next consolidation**, or every round will re-litigate dropped findings from memory.

## High

**H1 — Plan observability quietly requires the inventory §7 just rejected.** §7 keeps `plan_id ↔ operation_id` correlation so evidence can later answer "which plan produced this operation?" — and §10 gates tactical work "within approved scope." But correlation tags don't answer scope questions, and "approved scope" is undefined unless the plan carries **machine-readable scope**. The resolution isn't the rejected exact-identity inventory; it's the middle grain the lineage already built: **plans enumerate operation classes** (capability/operation-class level — "deploy workflow to repo Y"), exact identity forms at runtime from declaration + inputs + run_id. Human approves at class grain; Solvent authorizes at exact grain; drift becomes *detectable* by comparing operation class against plan scope through the evidence chain. Without this, §7's observability is a tag with nothing to check against and §10's scope test is circular. Note this also dissolves the plan↔operation conformance question without a Growth Gate — it's derivation, not new machinery.

**H2 — §6's MUST has no bearer.** "External-effect-capable tools MUST terminate at an enforcement boundary" — enforced by whom? For an OpenCode/Codex runtime with raw shell, the boundary exists only if the deployment builds it (gateway/proxy at the tool boundary — an implementation choice per the lineage, but a *mandatory* one here). Name the deployment operator as owner of boundary integrity, and make the operating rule explicit: an unmediated external path is a conformance finding, exactly parallel to undeclared effect-capable operations.

## Medium

- **Tactical/structural is agent-judged** (§10). Contained — misjudged "tactical" work still hits Solvent for external effects, and history makes it visible — but say so: the judgment is advisory; containment is Solvent + history + human visibility of activity.
- **Multi-agent containment and creation policy** remain unassigned. With plan scope machine-readable (H1), constrain claims to approved scope by construction; "who may CREATE outside any plan" stays a named deployment decision.
- **`capability_ref` is now role-less** in the schema. Classification moved to the boundary; the field survives only as interface/matching metadata. Say that, or apply the subtraction razor to it honestly.
- **Strawman entries inflate the rejection ledger.** Nobody proposed a "cryptographic Conductor↔Solvent ledger"; rejecting it isn't discipline, it's padding. The ledger should record what was actually on the table — same honesty standard as the matrix's disposition table.

## What survives attack

The two-mode lock is clean and mechanically distinguishable. The approval≠authority invariant is stated in exactly the right form. The no-policy-router position is correct and correctly *rejects* Conductor-side classification. Enforcement-at-the-boundary survives with an ownership question attached. The plan↔operation correlation, deferred prevention, and most rejections (formulas, Dolt, leases, READY_FOR_AUTHORIZATION) are all right. The final sentence is one definition away from being the correct lock.

## Routed order

1. **C1** — repair the definition; then lock the final sentence with "declared capable of producing an external effect" language.
2. **C2** — restore the four dropped rules (resolution path, acceptance owner, history fences, attribution) and create the disposition log.
3. **H1** — class-grain plan scope; exact identity stays at runtime.
4. **H2** — name the deployment operator as boundary-integrity owner.
5. Mediums as line edits.

The impasse resolution holds. What's left is refusal to let a consolidation be looser than the documents it consolidates — the design earned its smallness in seven rounds of review, and it keeps that smallness only if every subtraction is a *recorded decision* rather than a dropped thread.