# BM–IST Synthesis v5 — Implementation Plan
### A Consolidated Research Architecture for Bohmian Mechanics, Invariant-Set Theory, Arithmetic Dynamics, and the Projection Problem

**Governs:** BM–IST Synthesis v5.1 (47 sections) and all future extensions
**Discipline:** Evidence-Based Physics (EBP) — ideas enter free; promotion requires debt retirement; nothing is "final truth"
**Status of this document:** Execution plan, not a theory. It tells you what to do next and in what order, not what the answer is.

---

## 0. Why This Plan Exists, and What It Is Not

v5.1 already contains a dependency graph (§40), a priority checklist (§41), and success criteria (§43). This document is not a restatement of those. It is a **execution-layer plan** that does three things v5.1 doesn't yet do on its own:

1. Maps every gate, section, and claim in v5.1 onto your existing EBP debt-item machinery (`needMap`, `needInvariant`, `needToyCheck`, `needNullModel`, `needObstruction`, `needFaithfulnessReview`), so "is this promotable" has a checkable answer instead of a status tag someone assigned by feel.
2. Names, as a standing discipline rather than a one-off finding, the specific failure mode the recent Gemini adversarial review just caught: **borrowed rigor** — a claim dressed in a real theorem's name but applied to an object that theorem doesn't govern. This is a distinct and more dangerous failure than "asserted, not derived," and it needs its own gate.
3. Sequences the work by what's *actually computable right now* versus what needs new theorems first — which is not the same ordering as "logical dependency," and conflating the two is why the program has spent more cycles describing shapes for π than testing any one of them.

---

## 1. Governing Discipline: Retrofitting EBP onto v5.1

### 1.1 Status-tag translation

v5.1 uses `[EXISTING]`, `[THEOREM-SHAPED]`, `[SPECULATIVE-SHAPE]`, `[ASSERTED, NOT DERIVED]`. Your Workbench uses the explanation ladder: **Derives / Explains / Accommodates / Postulates / Undeclared**. Going forward, every claim added to the architecture carries *both*, because they're not measuring the same thing — one measures how finished the math is, the other measures what kind of epistemic work the claim is doing:

| v5.1 tag | Workbench rung | Meaning when combined |
|---|---|---|
| `[EXISTING]` | Derives | A real theorem, correctly applied, with the derivation shown or citable in full |
| `[EXISTING]` | Accommodates | A real theorem, but fitted after the fact to match a target rather than predicting it |
| `[THEOREM-SHAPED]` | Postulates | A specific, falsifiable statement, not yet proven |
| `[SPECULATIVE-SHAPE]` | Postulates / Undeclared | A structural analogy with no proof obligation discharged |
| `[ASSERTED, NOT DERIVED]` | Undeclared | Should not be in the architecture at all until re-tagged one level down |

**New required tag — `[BORROWED]`.** Every claim in the recently-rejected Gemini expansion (Hecke-eigenvalue/Ramanujan–Petersson misapplication, unconditional Furstenberg rigidity, the adelic connection form on a profinite structure group) shares one signature: a correctly-named, correctly-stated theorem is invoked, but the object it's applied to does not satisfy that theorem's hypotheses. This is worse than `[SPECULATIVE-SHAPE]` because it reads as `[EXISTING]` on a skim. Tag it explicitly, and require **function-class typing before derivation** (§1.2) as the standing countermeasure.

### 1.2 Function-class typing — the check that would have caught the Hecke-algebra error

Before attempting to close any gate, classify the claim by what *kind* of mathematical object is doing the work, and verify the named theorem's hypotheses actually match:

- **Existence theorem** (e.g., "a canonical height exists on this variety") — check: does the object have the algebraic structure (a polarized dynamical system, a number field, etc.) the theorem requires?
- **Selection principle** (e.g., "entropy saturation selects the physical measure") — check: is the selection *forced* by the dynamics, or merely *consistent with* a measure chosen for other reasons?
- **Construction** (e.g., "π is a connection on a fiber bundle") — check: does the structure group actually admit the differential-geometric objects (Lie algebra, curvature form) the construction writes down?
- **Phenomenological fit** (e.g., "a₀ ≈ cH₀/2π") — check: is this a target to explain, or is it being silently treated as already explained?
- **Analogy** (e.g., "gaps act like a pressure gradient") — check: has it been promoted past analogy without a derivation?

**Rule:** no claim advances past `needMap` in the debt ledger until its function class is stated in one sentence and its hypotheses are checked against the object it's being applied to — not against a nearby object with the same name.

### 1.3 Gate R — mandatory adversarial review, made standing rather than incidental

The review just run against the Gemini file is not a one-off correction; it's a pipeline stage. **No claim is promoted to candidate-artifact status without first passing Gate R.** Concretely:

- Every new claim gets a `CLAIM_ID` (v5.1's own ledger convention is good — keep it).
- Before promotion, a dedicated adversarial pass checks: (a) function-class match (§1.2), (b) whether cited theorems' actual hypotheses hold, (c) whether the claim internally contradicts anything already in the architecture (the P¹(ℚₚ)-is-secretly-Archimedean self-contradiction in the Gemini file would have been caught here on a first pass, since the file's own opening answer already ruled it out).
- A claim that fails Gate R is not deleted — per EBP's own discipline, ideas enter free — it's returned to `needToyCheck`/`needObstruction` debt status with the specific failure recorded, so it isn't silently re-attempted with the same error six months from now.

---

## 2. Current State of the Architecture (Honest Inventory)

| Status | Item |
|---|---|
| **Closed** | G0 (§3) — orbit rigidity + one-parameter triviality, both with proof sketches as of v5.1 |
| **Executable now, unrun** | Scaling audit (§41 Task 1), Elephant-Scan protocol S1–S7 (§24) — fully specified, needs code, not new theory |
| **Theorem-shaped, unattempted** | T-Prob-1 (§7), T-Born (§17), Gate 2 (§26), Gate 3 (§27), Gate 4 (§28) |
| **Speculative-shape, standing** | Haar-slaving (§14), entropy-saturation *as originally scoped* (§15 — the Gemini "Global Saturation Theorem" variant is rejected, not this), embedded-agent bootstrap (§16), holographic bulk-boundary kernel (§20), Arakelov QP conjecture (§20.2) |
| **Attempted and failed Gate R** | Gate C via Iwahori–Hecke algebras / Langlands (the entire Gemini WP7/WP8 chain) — back to open, with the specific failure modes on record (§3 below) |
| **Open, no attempt yet** | Gate C (§29) from first principles — the hardest unresolved obligation in the architecture |
| **Architectural decision, not yet made** | P1 (trace-relative) vs. P2 (autonomous) status of the p-adic sector (§33) |
| **Deliberately deferred** | Empirical handle table (§35), gravity/dark-matter branch (§34) — real, but downstream of the above and a magnet for premature numerology |

This inventory is the actual starting point. Everything below is organized to move items down this table in an order that matches what's computable, not what feels foundational.

---

## 3. Phase Plan

### Phase 0 — Intake & Debt Audit *(short, do first, low effort)*

Run every existing claim in v5.1 through the §1.1/§1.2 retagging. Deliverable: a single table, one row per `CLAIM_ID`, with v5.1 tag, Workbench rung, function class, and debt status. This is bookkeeping, not research, but it's the thing that turns "47 sections of a document" into "a queryable ledger" — and it's what would let a tool (see §5) eventually help instead of you re-reading the whole document by hand every time a new candidate import shows up.

**Retire:** `needMap` for the whole architecture at once.

### Phase 1 — Foundational Computation *(highest priority, pure execution)*

This is the one phase that needs no new mathematics, only code, and it gates almost everything else:

1. **Run the scaling audit.** Implement the ℤ₂×ℤ₃ skew-product family from §24, scan the coupling parameter λ across the actual regimes that matter (rational/Salem corners vs. reciprocal-Pisot corners), and compute local dimension and Fourier-decay rate at each point — using the correctly-scoped tools (Solomyak, Hochman–Shmerkin *local entropy averages*, Varjú decay estimates), **not** an unconditional Furstenberg-rigidity shortcut. Output: a phase diagram, not a proof — this is a numerical map of where BM–IST's own construction actually sits, which no amount of literature review substitutes for.
2. **Run the Elephant-Scan protocol S1–S7** against whichever regime the scaling audit identifies as most physically motivated for the actual invariant-set construction (not Palmer's).

**Deliverable:** two pieces of code, one phase diagram, one short written verdict: which arithmetic regime BM–IST's own substrate lives in, and whether that regime is inside or outside the danger zone for Gate 2 (Fisher information finiteness).

**Retire:** `needToyCheck` for the measure-existence question (P1). This does **not** retire `needToyCheck` for Gate 2 or Gate 3 — a favorable phase diagram is evidence, not the theorem itself.

**This is the load-bearing recommendation of this plan:** stop generating new candidate mechanisms for Gate 2/Born-rule grounding (Haar-slaving, entropy-saturation, T-Born, the now-rejected Hecke/entropy variant) until Phase 1's actual numbers exist. Every one of those candidates is currently being evaluated against an unknown target — you don't yet know, from your own construction, which arithmetic regime you're even trying to prove absolute continuity *for*. That's backwards, and it's exactly how the program ends up importing increasingly baroque machinery (Langlands, Hecke algebras) to solve a problem whose actual numerical shape has never been computed.

### Phase 2 — Measure & Probability Theorems *(conditional on Phase 1's output)*

Once Phase 1 identifies the actual regime:

- If it's favorable (transversal/non-Pisot): attempt **T-Prob-1** (§7) as literally stated, using the real toolkit (Solomyak/Hochman–Shmerkin/Varjú) already scoped for it. This is a real, attemptable theorem, not a speculative shape — it just hasn't been tried yet because Phase 1 hadn't happened.
- If it's unfavorable (Pisot-trapped): this is not a dead end — under EBP, a clean negative result is promotable. Deliverable becomes a documented **no-go finding**: "the construction as currently specified in §[X] lands in a singular regime; here is the specific structural change needed" — which is more valuable than another speculative rescue mechanism, and is the honest analog of what the Elephant-Scan protocol was designed to produce even on a null result.
- Attempt **T-Born** (§17) only after P1/P2 (measure existence, projection regularity) are actually settled by the above — not in parallel with them. T-Born presupposes a measure to equidistribute toward; without Phase 1's output, "attempting" it is choosing a target and reasoning backward, which is exactly the `Accommodates` rung rather than `Derives`.

**Retire (if favorable branch):** `needInvariant`, partial `needObstruction` for P1–P2.
**Retire (if unfavorable branch):** `needObstruction` fully, with the obstruction documented as a citable negative result; reopens `needMap` for a redesigned substrate.

### Phase 3 — π Construction *(pick one shape, stop running three in parallel)*

v5.1 §22 carries three candidate shapes (A: holographic/RG/QEC composite; B: adelic connection; C: sheaf-gluing/elephant-scan). Running all three indefinitely in parallel is itself a discipline failure — it's how the Gemini thread ended up bolting Langlands machinery onto Shape B without ever testing whether Shape C (which already has an executable protocol) works first.

**Recommendation: designate Shape C as primary**, because it is the only one of the three with a fully specified, already-executed-in-Phase-1 numerical protocol. Shapes A and B remain in the ledger as `needToyCheck`-status backups, not abandoned, but not worked on in parallel with Shape C until Shape C either succeeds or produces a documented obstruction.

- Attempt the existence proof for the gluing/connection map against Phase 1/2's actual output (not a hypothetical regime).
- Any future import of outside machinery (a new TOE feature, a new number-theoretic tool) must first pass function-class typing (§1.2) against Shape C's specific objects — the attractor, its generating partition, its transfer operator — before being written into the ledger.

**Retire:** `needInvariant` for the projection map itself, contingent on Phase 2's output.

### Phase 4 — Gate C (Emergent Tensor Product) *(hardest open item — restart from a minimal toy model)*

The Hecke-algebra/Langlands attempt at Gate C failed Gate R outright (category error on which "Hecke algebra" was being invoked; unproven coproduct structure; linearity kill test never actually checked). This gate remains genuinely open. Recommendation, in order:

1. **Do not reach for heavy machinery first.** Build the smallest possible toy model — two coupled minimal cyclic systems at small L — and numerically test whether *any* natural composition rule on the bit-string substrate produces (a) linearity in the amplitude sense, (b) no-signaling, and (c) genuine entanglement, simultaneously. This is a computational experiment, not a literature search.
2. Only if a toy-model composition rule is found to satisfy all three properties at small L should its generalization be sought in existing algebraic machinery — at that point, the machinery is being fitted to a demonstrated phenomenon, not asserted to produce one.
3. If no such rule is found at small L, that is a citable negative result narrowing the space of viable constructions, exactly as intended by the No-Free-Lunch discipline (§42).

**Retire:** nothing yet — this phase's entire purpose is producing the first real toy-check evidence on Gate C, which doesn't currently exist despite the Gemini file's confident ledger entries.

### Phase 5 — Architectural Decision: P1 vs. P2 for the p-adic Sector

This is a decision, not a derivation, and it should be made explicitly and early rather than left implicit:

- **P1 (trace-relative):** the p-adic bath has no autonomous physical content; it exists only to generate dissipation via Feynman–Vernon tracing. Simpler, fewer downstream obligations, closes off any dark-matter reading.
- **P2 (autonomous):** the p-adic bath has independent dynamical/gravitational content. Opens the door to the dark-sector/MOND-flavored branch (§34) but adds real obligations (an actual gravity-coupling mechanism, "Pillar G," which is currently undefined).

**Recommendation:** default to P1 until Phases 1–4 are further along. P2 should only be adopted once there's a concrete reason to need it — not because the shadow-realm analogy is appealing. This single decision, made explicitly, prevents the program from accumulating empirical obligations (§34's presence-predictions) for a branch that hasn't been chosen yet.

**Retire:** `needMap` for the p-adic sector's ontological status, whichever way it's decided.

### Phase 6 (conditional) — Empirical/Gravitational Branch

Only opens if Phase 5 selects P2. Work: pin down the single (c, ξ₀, κ) parameter set and run it against every row of the handle table (§35) with no per-row fitting — the "1919-eclipse standard" already named in v5.1. Do not begin this phase to chase a₀ or the 5:1 dark-to-baryon ratio as isolated numerology; both are explicitly flagged in the architecture's own cautions (§37) as prior failure patterns (Kolmogorov-length/Higgs-mass, E8/Lie-group) to avoid repeating.

---

## 4. Recurring Discipline (applies at every phase, not a one-time step)

1. **Gate R on every promotion, no exceptions** — including future outputs from any model, including this one. The Gemini file passed its own internal "kill test" ledger checks by construction (it wrote both the claim and the test); an external adversarial pass is the only thing that actually catches borrowed rigor, and it has to be applied uniformly, not just when something already looks suspicious.
2. **Function-class typing before derivation**, every time a new mathematical tool is imported (§1.2).
3. **No parallel speculative branches on the same gate.** Phase 3's "pick one shape" rule generalizes: when multiple candidate mechanisms target the same gate (as currently happens at Gate 2's probability question), designate one as primary and the rest as backups in `needToyCheck` status, rather than developing all of them simultaneously — this is precisely how the program acquired an unreviewable backlog of ledger entries in a single conversation.
4. **Negative results are deliverables, not failures.** Every phase above states what a null result looks like and treats it as promotable. This is not a consolation prize — under the Elephant-Scan protocol's own original design, a documented empty critical surface was always meant to be as citable as a positive result.
5. **Provenance ledger stays append-only and dated.** Rejected claims (the current Gemini WP7/WP8 chain) stay in the ledger with their failure mode recorded — don't delete them, since a future model or a future you may re-derive the same dead end without that record.

---

## 5. Tooling Note: `ebp-paper-evaluator`

The existing Go tooling is the natural home for automating the Phase 0 retagging and the Gate R checklist — but per your own prior technical evaluation of it, it currently has a non-functional security model, a likely-broken evidence-span verification system, and headline scores dominated by schema compliance rather than content quality. **Do not route Gate R decisions through it as-is.** Two options, not mutually exclusive:
- Treat the tooling repair itself as a small, separate, well-scoped side task (its own debt ledger), decoupled from the physics program's critical path.
- In the meantime, run Gate R manually using the checklist in §1.3 — it's short enough to apply by hand to each new `CLAIM_ID` without the tool.

---

## 6. Risk Register

Named explicitly because each has already occurred once in this program and will recur if not actively watched for:

| Risk | Where it already happened | Standing countermeasure |
|---|---|---|
| Borrowed rigor (real theorem, wrong object) | Gemini's Hecke-eigenvalue/Ramanujan–Petersson claim; unconditional Furstenberg rigidity | Function-class typing (§1.2) before any derivation attempt |
| Self-contradiction across turns going unnoticed | P¹(ℚₚ) declared "Archimedean" after the same file's own opening answer ruled this out | Gate R explicitly checks new claims against the full standing ledger, not just against the immediately preceding turn |
| Ledger inflation (format without verification) | A dozen `LEDGER ENTRY` blocks generated in one Gemini session, none executed | No claim reaches candidate-artifact status without a completed toy check, not just a stated kill test |
| Premature numerology | Prior Kolmogorov-length/Higgs-mass and E8/Lie-group conflations (self-named in the source material) | Phase 6 gated behind an explicit P1/P2 decision (§5) and the 1919-eclipse no-per-row-fitting rule |
| Parallel-shape sprawl | Three π shapes and four probability routes developed simultaneously with no resolution criterion | "One primary, rest as backups" rule (§4.3) |
| Analogy creep back into speculative-shape claims | The original four-source consolidation's own flagged self-corrections (flipbook, pressure-gradient, p-adic Cantor stability) | Explanation-ladder tagging (§1.1) catches `Postulates`/`Undeclared` claims before they're treated as settled |

---

## 7. Success and Stopping Criteria

- **Local success (per phase):** each phase above has a stated deliverable that is either a proof, a numerical result, or a documented negative result — never another round of prose description.
- **Program-level success:** T-Prob-1 or its negative-result equivalent, a chosen and partially-verified π shape, and a Gate C toy-model result (positive or negative) together constitute a publishable, EBP-promotable core — independent of whether the gravitational/dark-matter branch (Phase 6) is ever opened.
- **Honest stopping condition:** if Phase 1's scaling audit places the actual construction in a singular regime *and* Phase 4's toy model finds no composition rule satisfying linearity, no-signaling, and entanglement simultaneously, the program has produced two citable no-go results rather than a completed theory — which, under EBP's own discipline, is a legitimate and fundable outcome, not an abandoned project.

---

## Closing Note

The single change this plan makes to how the program has been running is sequencing: compute first, theorize second. Every phase above that involves new mathematics is gated behind a phase that involves only running code against an already-fully-specified protocol. The Gemini episode is useful precisely because it's a clean demonstration of what happens when that order is inverted — a great deal of correctly-named mathematics arrives before any of the numbers it's supposedly explaining have actually been computed.
