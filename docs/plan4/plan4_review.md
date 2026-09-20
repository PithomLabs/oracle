Consistency note, same as every round: this consolidates my review, so I've re-attacked the text rather than graded the merge. Result up front: **this is the closest round to lockable — but the corrected definition in §2 contains a grain error that contradicts the document's own §3 and §12, and if locked as written it breaks on the first real coding agent.** One critical repair, one integrity gap in the restoration, one process debt — then this is done.

## Disposition of prior findings

| Prior item | Status | Notes |
|---|---|---|
| C1 definition regression ("intentionally," "authoritative state") | **Half repaired** | Intent is gone, external effect defined. New defect below — see C1. |
| C2 four dropped rules | ✓ Restored | §6A–D land all four, in near-final form. Residue on A: see H2. |
| H1 class-grain plan scope | ✓ Landed | §4 is exactly the right compromise. Residue: see M2. |
| H2 enforcement bearer | ✓ Landed | §8 assigns it to deployment/conformance. Correct. |
| Tactical/structural advisory | ✓ Landed | §9, with correct containment reasoning. |
| Disposition log (demanded "before the next consolidation") | **Not created** | Escalated — see H1. |

## Critical

**C1 — §2's definition classifies at capability grain; §3 and §12 classify at operation grain. The document disagrees with itself, and both readings of the disagreement are bad.** §2: an operation is consequential "when its **declared capability** is capable of producing an external effect." §3: the branch occurs "when the Agent actually reaches an **effect-capable operation**." §12: "the declaration determines whether an **operation** is consequential." For a universal capability — shell, python, file tools — the capability as a whole is trivially capable of external effects. So capability-grain classification forces a fork with no good side: declare shell effect-capable and *every command* — every `ls`, every file write — demands Solvent, making Work Mode unusable and plan approval ceremonial; declare shell non-effect and §8's mediation requirement never attaches to it, so the raw-external-call bypass reopens. The stable point was settled rounds ago in the matrix: **declarations cover the operation *or operation class*** — shell declares local-file operations non-effect, network/push/deploy classes effect-capable. The consolidation's definition dropped the class grain. Lock this instead:

> An operation is consequential when the operation class it instantiates is declared effect-capable — capable of producing an observable state change outside the workflow — regardless of agent intent, workflow phase, or invoking client.

That sentence agrees with §3, §12, and the matrix simultaneously. This error would otherwise be discovered empirically, mid-run, by an agent whose first `git push`-adjacent shell call either sues for authorization it shouldn't need or skips a boundary it should hit. One sentence now.

## High

**H1 — The disposition log still doesn't exist, and this round proves both its necessity and its absence.** The four restorations in §6 happened because a reviewer enumerated them from memory. That works once; it is not a mechanism — the same four were silently dropped exactly one round ago, and the stale-claim × intent-expiry scenario has now been dropped twice. Minimum viable form, three columns: *finding → disposition → one-line reason*, appended to every consolidation, with provenance (which review raised it). Note §7's rebuttal illustrates the need: "the first review proposes claim TTLs, bounded plan scopes, proof tokens" — against this conversation's record, claim-TTLs were *explicitly deferred by me* (Phase 3, evidence-driven) and proof tokens were never proposed; the resolution reached (defer both, adopt bounded scope) is correct, but the attribution is checkable only if the log exists. If "the first review" is the other reviewer's, the log is how that gets arbitrated instead of asserted.

**H2 — Declaration resolution lost its integrity pin in the simplified restoration.** §6A: "the owner of the declaration serves it; Conductor merely records the reference." The prior fix was reference **+ version + content hash**, consumer verifies. The hash was silently simplified away — and classification is enforcement-critical: the declaration is what decides whether Solvent gets engaged at all, so a declaration source that can change content under a stable version string is a mutable security boundary. Restore the pin. Also tighten "cannot become actionable" to a named check point — fail-closed at CREATE, or exclusion from READY — so dangling references can't populate the frontier.

## Medium

- **M1 — The §11 diagram regresses to Conductor-in-the-consequence-path, third recurrence.** The arrow sequence reads Agent → Conductor → Solvent → Executor, while the text correctly says Conductor never routes on consequence. Draw the agreed two-plane picture: Agent with two edges — work protocol to Conductor, consequence path to Solvent — Conductor observing results. Diagrams get implemented; this one has now been wrong in three documents.
- **M2 — Plan-scope classes must resolve to declaration-defined class identifiers.** §4's `DeployWorkflow` scope and the declaration's class taxonomy must be the same vocabulary, resolved with the same discipline as `capability_ref` — otherwise approved-class vs. actual-class comparison (§5's detection story) degenerates into string matching. One line in the plan model.
- **M3 — Stale-claim × intent-expiry remains unowned.** Not TTL machinery — a Phase 1.5 scenario: submit after intent validity expires → which rejection class? Second silent carry; it joins the log H1 creates.

## What survives attack

The two-mode lock, the approval≠authority invariant, the no-policy-router position, operation-class plan scope with runtime exact identity, the four restored fences, the tactical/structural containment argument, and the deferrals (TTLs, tokens, auto-enforcement) — all correct, all consistent with the lineage. The final §12 sentence is right. The §2 definition is the only sentence in the document that isn't.

## Routed order

1. **C1** — replace the §2 definition with the operation-class formulation; then the lock sentence is clean.
2. **H2** — version+hash pinning and a named resolution check point in §6A.
3. **H1** — create the disposition log before the *next* consolidation, not this one's aftermath.
4. M1–M4 as line edits.

With C1 landed, I have no remaining design objections to the workflow layer. The open risk concentrates where it has always concentrated — declaration content, identity comparison semantics, and the enforcement boundary in real agent runtimes — and those are exactly what Phase 1.5 and Phase 2 exist to attack. This is the round where the design stops changing and starts being tested.