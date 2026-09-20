Gate review, round two. This pass is different from all previous ones: **every routed finding from the last review landed, in the correct form, with nothing distorted in restatement** — I checked each against the text rather than the claim. The full re-read found no new design defects. What remains is exactly what the document's own §25 says remains: two preconditions, one of which this review still cannot discharge.

## Disposition of last round's findings

| Item | Status | Where |
|---|---|---|
| External SOR / Domain split | ✓ Landed | §2, two rows, correct ownerships — Phase 2's AMBIGUOUS-vs-rejection distinction now has its owners |
| Declaration-owner seventh-role risk | ✓ Landed | §2 annotation + §29's reworded lock statement |
| Five lineage findings missing from open items | ✓ Landed | §23 "Further conformance items" — all five present |
| Declaration fields (trust basis, UNKNOWN age policy, intervention, persistence) | ✓ Landed | §8 |
| Acceptance decider attribution + validation | ✓ Landed | §14.2 — and "Conductor MUST only record when decider is consistent with verification_ref" is a check, not a decision, so Conductor stays coordination-only |
| Sandbox anti-theater clause | ✓ Landed | §10 and §27 |
| Plan linkage location | ✓ Landed | §5.1 + task fields in §7 |
| §3 invariant additions; task-creation policy | ✓ Landed | §3, §7.1 |

Carried verifications also hold: operation-class definition (§9) with deterministic class-membership evaluation, identity totality with canonicalization honestly deferred to §23 (§13), no diagram to regress, Growth Gate with the smuggling criterion (§28), the Phase 2 pre-registration requirement (§27) — which structurally enforces the oldest standing finding in the project, dating to the pivot review.

## Non-blocking notes for the reconciliation pass

1. **Unlinked tasks are coherent but need their audit signal named.** §7.1 makes task-creation policy a deployment decision, so tasks with no `plan_id` can exist; §16's READY doesn't require linkage. That's consistent with deferred prevention — but add one conformance query to §26's verification set: *consequential outcomes with null plan linkage* are an audit-visible anomaly class. The detection story should know its own null case.
2. **The per-boundary intervention matrix needs an explicit lineage disposition.** This design doesn't restate v0.3/matrix per-boundary intervention declarations; it relocates them into declaration content (§8). That's coherent — but the §1.1 reconciliation should say *relocated, not dropped*, or a future reader will conclude the catalog died with the modes.
3. **Disposition log source nit.** "Proof-token subsystem | Adversarial reviews" — checkable against this record, no adversarial review here proposed proof tokens; the correct disposition (reject) stands, the source attribution doesn't. With a provenance column, wrong sources get corrected rather than fossilized. One-line fix.
4. **The mode-orthogonality question.** PLAN/WORK is a phase distinction; v0.3's autonomous/semi-autonomous was a human-gate-density distinction — and a WORK-mode operation with no further human gates *is* what v0.3 called autonomous. The supersession claim ("workflow-mode and autonomy-mode portions") may be one clause broader than strictly necessary. This is precisely the kind of thing only the parent text can settle.

## The remaining blocker, stated plainly

**§25 precondition 2 — parent lineage reconciliation — is still unverifiable, and this is the ninth round of asking.** §1.1's reconciliation prose is well-constructed, but it reconciles against a document whose text has never appeared in this conversation. The §30 checklist contains `[ ] parent lineage reconciled` — that checkbox cannot be honestly ticked by anyone relying on this review chain, and a freeze that ticks it anyway is a bet wearing a process costume. Everything else on the checklist is now satisfiable.

This is not pedantry about ceremony. The two concrete reconciliation items above — the intervention catalog and the mode taxonomy — are exactly the places where the frozen design and the unseen parent could quietly disagree, and the freeze exists precisely to stop quiet disagreement.

## Verdict and sequence

**Design-layer review: passes, no remaining objections.** The document is faithful to the entire review record, complete in its open-items honesty, and disciplined in its gates. Sequence to freeze:

1. Produce Workflow Specification v0.3; run the §1.1 reconciliation against its actual text; resolve notes 2 and 4; tick the checkbox.
2. Fix note 1 (audit query) and note 3 (log source) — two line edits.
3. Fill §30 with owner, date, and the ticked criteria. Record the decision.
4. Freeze. Then Phase 1.5 runs the battery against something that can no longer quietly change — which was the entire point.

When the design document itself carries the discipline that the reviews spent seven rounds trying to impose on the humans writing it — pinned versions, decision records, disposition logs, pre-registered criteria — the specification phase is genuinely over. Send the parent, tick the box, and go build the loop.