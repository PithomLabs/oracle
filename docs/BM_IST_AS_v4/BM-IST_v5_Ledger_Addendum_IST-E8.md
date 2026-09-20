# BM–IST(-AS) v5 — Ledger Addendum: IST Foundational Revisions and E8 Attachments

**Governs:** BM–IST Synthesis v5.1, the v5 Implementation Plan, and the AS-template material under adversarial review two turns prior
**Adds:** six new `CLAIM_ID` entries (§1.1–1.4, §2.i–2.ii) plus one explicit rejection record, per the same ledger convention as the Valentini addendum
**Entry status:** all enter as debt. One entry (§2.ii, naive form) enters as a **recorded rejection**, not debt — per the standing rule from the Gemini adversarial review, a failed claim stays in the ledger with its failure mode on file rather than being deleted, so it isn't re-attempted the same way later.

---

## §1 — Invariant Set Theory: foundational revisions

### CLAIM_ID: IST-L-CFUNC-1

**Statement.** The substrate's finite information capacity \(L\) is not a hand-set constant (Palmer's \(L\sim10^{100}\), or the holographic \(L\sim10^{120}\) floated earlier as a better-motivated but still hand-set alternative). Instead, \(L\) is identified with a monotonic entanglement-entropy-type \(c\)-function evaluated along the substrate's AS-template RG trajectory, in the spirit of the Zamolodchikov \(c\)-theorem / \(a\)-theorem family of results.

**Function class.** Selection principle (what fixes the value of a parameter that was previously a free input). Per the typing discipline established after the Hecke-algebra episode: the check this class requires is whether the identification is *forced* by the flow's own structure, or merely *consistent with* a number chosen for other reasons.

**Status.** `[SPECULATIVE-SHAPE]` / Workbench rung **Postulates**. The general strategy ("capacity tracks a $c$-function") rests on real, existing monotonicity theorems; the specific claim ("this $c$-function value gives this bit count, for this substrate") is unconstructed.

**Open debt.**
- `needMap` — no explicit map from the substrate's RG trajectory to a computed \(L\) value exists yet.
- `needObstruction` — none stated prior to this entry (see Kill condition).

**Kill condition.** This claim is unattemptable until the AS-template's own substrate RG flow equation exists — flagged as step 1 of an unexecuted 10-step program in the AS-template adversarial review. If that flow equation is constructed and its associated \(c\)-function is computed but does not track \(L\) in any regime, or tracks it only after additional free parameters are introduced to force agreement, the identification fails and \(L\) reverts to a to-be-determined status rather than reintroducing a hand-set number.

**Dependency / phase placement.** Blocked on the AS-template flow-equation construction (currently unconstructed — see prior review, Finding 2). Not assignable to a specific Implementation Plan phase until that prerequisite exists; treat as parked pending Phase 3 (π construction) groundwork.

---

### CLAIM_ID: IST-DEFCUT-HEIGHT-1

**Statement.** The rational/irrational (Niven's-theorem) definability cut is replaced: the set of counterfactually-available, physically definable states is the preperiodic locus (canonical height \(\hat h=0\)) of the substrate's own dynamical map, per Call–Silverman, rather than an externally imposed arithmetic rule applied to a fixed coordinate choice.

**Function class.** Existence theorem (does a well-defined, intrinsic, coordinate-free "defined set" exist for a degree-≥2 self-map) plus a fit obligation (does that set reproduce Niven's actual empirical content — the specific angles known to survive).

**Status.** `[SPECULATIVE-SHAPE]`. The existence of a preperiodic locus for suitable maps is general, citable dynamical-systems mathematics (Northcott finiteness; Call–Silverman); what's unconstructed is whether *this* substrate's specific map has one, and whether it reproduces Niven's results rather than a different set entirely.

**Note on origin.** This is not new to this turn — it formalizes, as a standing IST-foundation revision rather than a Gate-2 patch, machinery already flagged in the original four-file consolidation (arithmetic dynamics / canonical heights, §4.3) and the T-Prob-1 line of work. It's ledgered here explicitly because it's being adopted as a *replacement* for Palmer's own posited mechanism, not merely an alternative route to Gate 2.

**Open debt.**
- `needMap` — the substrate's actual dynamical map has not been shown to have a nonempty, correctly-structured preperiodic locus.
- `needFaithfulnessReview` — does the resulting set match Niven's known angles, or diverge from them? A divergence isn't automatically fatal (it could be a genuine generalization), but it needs to be checked and reported honestly rather than assumed to match.

**Kill condition.** If the substrate's map is shown to have degree \(<2\), or to fail the conditions Call–Silverman's theorem requires, no preperiodic locus of the needed kind exists and this replacement fails outright — Niven's original mechanism (with its coordinate-artificiality problem intact) would need to be retained or a third mechanism sought.

**Dependency / phase placement.** Phase 2 (Measure & Probability Theorems), alongside T-Prob-1 — same substrate-construction prerequisite, same phase.

---

### CLAIM_ID: IST-ALPHABET-MARKOV-1

**Statement.** The p-adic metric is not posited as an input. Instead, the substrate's symbolic coding alphabet — and whatever ultrametric it carries — is derived from a Markov partition of the substrate's own dynamics (standard for sufficiently hyperbolic/expanding systems). Whether the resulting partition happens to have branching arity \(p\) for some specific prime, or a different arity entirely, becomes a computed output rather than an assumed input.

**Function class.** Construction. The check required: does the substrate's dynamics actually admit a Markov partition at all (a real precondition, not automatic for every dynamical system), and if so, what is its branching number?

**Status.** `[SPECULATIVE-SHAPE]`. Markov partition existence for hyperbolic/expanding systems is standard, citable dynamical-systems theory; nothing has yet shown the *specific* substrate dynamics used elsewhere in this architecture qualifies, or computed its actual partition.

**Open debt.**
- `needToyCheck` — construct the partition explicitly for the substrate used in the Phase 1 scaling audit and report its branching number.
- `needObstruction` — none stated prior to this entry (see Kill condition).

**Kill condition.** If the substrate dynamics fails to be sufficiently hyperbolic/expanding for a Markov partition to exist by the standard theorems, this replacement fails and "p-adic" reverts to an unmotivated posit pending a different derivation strategy. If a partition exists but its branching number bears no simple relation to any structure already load-bearing elsewhere in the architecture (κ, L, etc.), that's a weaker but still noteworthy outcome worth reporting rather than silently discarding.

**Dependency / phase placement.** This is a **prerequisite check, not a downstream item** — the Phase 1 scaling audit already assumes a specific (dyadic/triadic, ℤ₂×ℤ₃-type) symbolic alphabet. This entry should be run either before Phase 1 or explicitly flagged as an independent cross-check running in parallel with it, since a mismatch between the assumed alphabet and the derived one would mean Phase 1's own results need re-reading.

---

### CLAIM_ID: IST-RATE-RGMATCH-1

**Statement.** The substrate's fixed 1-bit-per-Planck-time collapse/shift rate is replaced by the AS-template's RG coarse-graining rate at the relevant scale, rather than being posited as a literal constant.

**Function class.** Selection principle, structurally identical in kind to IST-L-CFUNC-1 — a previously free-standing rate is proposed to equal a rate computed elsewhere in the architecture.

**Status.** `[SPECULATIVE-SHAPE]`, carrying the same borrowed-rigor risk explicitly named when this was first raised: "these are the same *kind* of thing" is not "these are numerically the same," and nothing here shows the substrate's actual mixing rate under the shift map matches the RG flow's actual coarse-graining rate.

**Open debt.**
- `needMap` — no computation exists comparing the two rates.
- `needObstruction` — none stated prior to this entry (see Kill condition).

**Kill condition.** If, once both rates are independently computable (the shift map's own mixing rate, and the AS-template flow's coarse-graining rate), they disagree by more than a simple, derivable rescaling, the identification fails and the collapse rate reverts to unexplained-constant status — it should **not** be quietly retuned to force agreement, which would constitute the same contamination failure flagged against the AS-template comparison table in the earlier review.

**Dependency / phase placement.** Blocked on the same AS-template flow-equation prerequisite as IST-L-CFUNC-1. First formal ledger entry for this claim — it was informally flagged during the AS-template adversarial review but never previously assigned a `CLAIM_ID`.

---

## §2 — E8 attachments

### CLAIM_ID: E8-PACKING-INFO-1

**Statement.** IST's finite-information substrate (once \(L\) is derived per IST-L-CFUNC-1 rather than posited) is encoded using the E8 lattice's sphere-packing optimality (Viazovska 2016) and its associated error-correcting-code structure, rather than an arbitrary or unspecified encoding.

**Function class.** Analogy, not yet promoted to selection principle. The mathematical fact borrowed (E8 packing optimality) is `[EXISTING]` and carries no borrowed-rigor risk in itself; what's missing is the physical reason nature would use the *optimal* code rather than any other.

**Status.** `[SPECULATIVE-SHAPE]`. Structurally sound but incomplete — the missing selection principle is plausibly the same entropy-saturation machinery already on the ledger from the earlier probability-gate work (a system saturating its information bound naturally favoring an optimal code is a reasonable conjecture), which is worth noting as encouraging: this doesn't require a *new* import, just a connection to one already load-bearing.

**Open debt.**
- `needMap` — the entropy-saturation-to-optimal-code connection has not been constructed, only conjectured.
- `needObstruction` — none stated prior to this entry (see Kill condition).

**Kill condition.** If the entropy-saturation measure, once actually computed (Phase 2), does not favor optimally-packed codes over other codes of comparable rate, this attachment fails — E8's packing optimality would then be a mathematically true but physically irrelevant fact about the substrate, not a structural feature of it.

**Dependency / phase placement.** Phase 2, blocked on both entropy-saturation (already open debt from prior probability-gate work) and IST-L-CFUNC-1.

---

### CLAIM_ID: E8-CFT-CHIRALITY-1 — REJECTED (naive form), recorded per standing ledger policy

**Statement (as proposed, now rejected).** "The required chiral fermions emerge not from the E8 group itself, but from the representation theory of the CFT that has the E8 lattice as its geometric dual."

**Function class.** Existence/resolution claim — asserted that a specific obstruction (Distler–Garibaldi) is resolved by a specific mechanism.

**Status.** `[BORROWED]` — not merely unproven, but checked and very likely false as stated. The E8 lattice generates a holomorphic vertex operator algebra (even self-dual lattice construction), and a lattice VOA of this type has **exactly one irreducible module: itself** (Dong 1993; Frenkel–Lepowsky–Meurman). There is no representation-theoretic richness in this specific object for a three-generation chiral spectrum to come from. Separately, "geometric dual" was found to be doing undefined work — the actual relationship (lattice defines VOA via vertex-operator construction) is a construction, not a duality in any established technical sense here.

**Recorded failure mode.** Single-module rigidity of holomorphic lattice VOAs forecloses the proposed mechanism outright. Do not re-attempt this specific claim without first addressing why the single-module obstruction wouldn't apply — e.g., by proposing a non-holomorphic or twisted variant of the construction, which would be a materially different claim requiring its own entry.

---

### CLAIM_ID: E8-CFT-CHIRALITY-2 (narrower surviving form)

**Statement.** *If* the substrate's spectral dimension is shown to flow toward 2 in the UV (an unconstructed claim from the original AS-template proposal, not yet formally ledgered), *and* the resulting 2D UV description is of lattice/rational-CFT type rather than another class (e.g., Liouville-type) — then an E8-lattice CFT is a legitimate, narrowly-scoped candidate for the gauge sector at that fixed point, with chirality supplied the way the heterotic string actually supplies it: by a separate topological/compactification-type structure (roughly, an index-theorem-fixed generation count), not by the E8-CFT's own representation theory.

**Function class.** Construction, doubly conditional, explicitly scoped as narrower than the rejected form above.

**Status.** `[SPECULATIVE-SHAPE]`, and explicitly **not** a resolution of Distler–Garibaldi. It is an acknowledged abandonment of Lisi's central claim (gauge bosons and fermions unified in one E8 structure) in favor of the ordinary heterotic-style split (gauge sector from E8, matter/chirality from elsewhere). This distinction must be stated honestly wherever this claim is used — presenting it as "resolving" the obstruction rather than "declining to need it resolved" would repeat the exact conflation caught in the prior review.

**Open debt.**
- `needMap` — the "separate topological/compactification-type structure" supplying generation number has not been specified for this substrate at all; unlike the heterotic string, there is no compactification manifold on offer here, so even the analogy's source mechanism has no substrate-side counterpart yet.
- `needObstruction` — none stated prior to this entry (see Kill condition).
- `needFaithfulnessReview` — confirm any future write-up of this claim does not drift back toward the rejected form's "representation theory resolves it" framing.

**Kill condition.** If the spectral-dimension-flows-to-2 claim fails (itself still unconstructed and un-ledgered), this entire entry is moot and should be marked inactive rather than pursued further. If that claim holds but the resulting UV theory is shown to be non-lattice-type (e.g., Liouville-type), this specific E8 attachment fails while the broader 2D-UV-description claim survives.

**Dependency / phase placement.** Strictly downstream of the (currently unformalized) spectral-dimension gate. Recommend formalizing that gate as its own `CLAIM_ID` before further work on this entry — flagged here as a gap rather than addressed, since it falls outside this turn's requested scope (§1.1–1.4, §2.i–2.ii only).

---

## Standing note, not a claim

Lisi's E8 dimension (248) is the Lie algebra's generator count, not an empirical particle census — no counting of Standard Model fields under any convention lands on 248, and no entry in this addendum should ever be cited as explaining "how many particles exist." This is recorded here as a standing caution for future write-ups, not as tracked debt, since it isn't a claim the architecture makes or needs to defend.

---

## Ledger Summary Table

| CLAIM_ID | Function class | Status | Blocking debt | Phase / dependency |
|---|---|---|---|---|
| IST-L-CFUNC-1 | Selection principle | Postulates | `needMap` | Blocked on AS-template flow equation |
| IST-DEFCUT-HEIGHT-1 | Existence + fit | Postulates | `needMap`, `needFaithfulnessReview` | Phase 2, with T-Prob-1 |
| IST-ALPHABET-MARKOV-1 | Construction | Postulates | `needToyCheck` | Prerequisite to / parallel with Phase 1 |
| IST-RATE-RGMATCH-1 | Selection principle | Postulates | `needMap` | Blocked on AS-template flow equation |
| E8-PACKING-INFO-1 | Analogy | Postulates | `needMap` | Phase 2, with entropy-saturation |
| E8-CFT-CHIRALITY-1 | Existence/resolution | **Rejected** | — (closed) | N/A |
| E8-CFT-CHIRALITY-2 | Construction (doubly conditional) | Postulates | `needMap` ×2 | Downstream of un-ledgered spectral-dimension gate |

All open entries follow Gate R before any future promotion past `needToyCheck` or `needMap`, consistent with every prior addendum to this architecture.
