# BM–IST v5 — IMPLEMENTATION PLAN
## From Standing Reference to Running Program

**Purpose.** Convert v5 from a blueprint into an executable program, in the correct order: repair the record first (per the Round-14 audit), stand up the verification infrastructure second, then execute fastest-decisive-first. This plan supersedes v5 §41/§46 (whose ordering was audited as mis-prioritized) and incorporates all ledger deltas through **#117**.

**Governing rule carried forward:**

$$\boxed{\text{Every phase ends in artifacts — theorems, code, datasets, or kill criteria. Never prose alone.}}$$

---

## 0. PLAN ARCHITECTURE AT A GLANCE

| Phase | Name | Duration | Decisive question it answers |
|---|---|---|---|
| P0 | Record repair + infrastructure | Weeks 1–3 | Is the record itself true? |
| P1 | Fast decisive layer (Papers + DT1/DT2) | Weeks 2–10 | Is the fatal scenario a theorem? |
| P2 | The elephant scan (merged WP2/WP11) | Months 2–6 | Does the critical regime exist? |
| P3 | Probability theorems | Months 4–12 | Can the marginal be regular *without* fine-tuning? |
| P4 | π and the amplitude | Months 6–18 | Can 𝒜 be constructed or excluded? |
| P5 | Gates C/4/contextuality → gravity | Months 12+ | Only if P4 converges |

Phases overlap deliberately (P3 begins on P2's best regimes before P2 completes). Every phase has a **stop/go gate (G-n)** and a **pre-registered kill criterion**.

---

## PHASE 0 — RECORD REPAIR AND INFRASTRUCTURE (WEEKS 1–3)

### 0.1 Produce v5.1 (mandatory, before any execution)

Apply all Round-14 repairs as a tracked diff:

| Repair | Action |
|---|---|
| C1 (#91) | Restore the closed-gate list: Theorem G2, Proposition L, separability lock, no-free-lunch trichotomy, shear counterexample — tagged `[PROVED IN PROGRAM; write-up owed]`. §20 must cite G2 for the coefficient lock |
| C2/C3 (#92–93) | Merge WP2≡WP11 (scaling audit = λ-classification of F_λ); restore the Mixing–Unitarity Dilemma as the standing fatal scenario with DT1–DT3 chartered |
| C4 (#94) | Re-order §40 chain: G0 → substrate (F_λ) → scan+scaling → **dilemma-resolved regime** → T-Prob-1 → T-Born → amplitude spec (S1–S4) → π → Gate C → phase/κ → guidance → contextuality → gravity |
| M1 (#95) | Downgrade #79: cyclotomic–Wallstrom is picture-dependent (point-limit ⟹ dense circulation; character-limit requires selection principle = Gate-4 content) |
| M2 (#96) | Restate T-Prob-1 with explicit Besov/Sobolev rung (√ρ ∈ H¹); separate static (T-Prob-1) from dynamic (T-Born) |
| M3 (#97) | Complete singular-corner list: {λ = 1/N, N ≥ 3} ∪ {non-integer reciprocal-Pisot}; note the scan family contains the N = 3 corner by construction |
| M4 (#98) | §11.2: supply the α → 2 locality-emergence argument or delete the item |
| D1–D10 (#99–105) | Holographic-L kill-test or removal; L1/L2 re-tagged; scan S2 criterion + cohomology-vs-λ reinstated; spectrally-active typed constraint; "unlikely" → "excluded for typed families"; Shape B/T-Born discriminator; standing verdict line; tag hygiene; ledger propagation rule |

Header of v5.1 carries the **live verdict line**: *C-not-D, dated, with current fatal scenario = Mixing–Unitarity Dilemma.*

### 0.2 Stand up the ledger and verifier POC (tick-driven, non-autonomous)

- **Ledger:** JSON store per §38 schema + propagation rule (#D9 fix): a `FALSIFIED` node auto-cascades `CHALLENGED` along its `kills`-edges.
- **Verifier POC tick-1 = Batch 1 (self-audit):** CAS verification of the program's own five computation-checkable claims:
  - C1: δI_F/δρ = −4Δ√ρ/√ρ (variational finite-difference on random smooth ρ)
  - C2: heat-flow Fisher decay identity
  - C3: **shear counterexample reproduction** (Fisher growth under b = (y,0) + small ν) — the flagship self-falsification
  - C4: Theorem G2's Madelung cancellation (one-step evolution of random (ρ,S) vs. Schrödinger step)
  - C5: separability identity on product densities
- **Batch 2:** citation audit of the ~25 anchors (Hudson 1974; Haussmann–Pardoux 1986; Guerra–Morato 1983; Hall–Reginatto 2002; Wallstrom 1994; Melbourne–Stuart 2011; JKO 1998; Hegerfeldt 1974; Rudolph 1995/Johnson; Bilu 1997; Hochman–Shmerkin; Varjú; Anashin; Abramsky–Brandenburger; Takens 1981; …), one content-anchoring quote each.

> **Gate G-1 (stop/go):** if C3 or C4 fails CAS reproduction, the record is wrong at a load-bearing joint — halt all downstream work and repair the mathematics first. The verifier's first mission is auditing its own authors, by design.

---

## PHASE 1 — THE FAST DECISIVE LAYER (WEEKS 2–10, OVERLAPPING P0)

The program's history (shear counterexample, L1/L2: one-week lemmas that changed the record) dictates this ordering. Four artifacts, all write-up-or-short-proof:

### 1.1 Paper A — "Conditional Fisher Rigidity" (target: weeks 2–6)
Contents: Theorem G2 with full proof; Proposition L; separability lemma; escape-hatch classification (T3: nonlocal, higher-derivative, S-dependent, internal-DOF); the no-free-lunch trichotomy. **This converts the program's zero published artifacts into one, and closes audit finding C1 permanently.**

### 1.2 Technical note — "G0: countable beables cannot carry continuous unitary flow" (weeks 3–5)
L1 + L2 written as a formal note with the generalized statement (#73: any countable-beable program, including the discrete-Hilbert genre). Sets the board for the field: unmodified-Schrödinger-on-discrete is dead everywhere.

### 1.3 DT1 — Bernoulli-horn contraction theorem (weeks 4–10)
*On the dyadic shift with Haar measure, every amplitude functional carrying the marginal weights and measurable w.r.t. the natural filtration evolves contractively — no unitary channel exists.* Tools: tail-boundary triviality, Pinsker factor, extension of no-free-lunch beyond conditional expectations.

### 1.4 DT2 — Odometer-horn context theorem (weeks 6–10)
*For odometer-type (uniquely ergodic, singular-continuous) substrates, preparation-context correlations persist without decay — Born context-independence fails formally.*

### 1.5 Micro-lemma (#116): Ẑ-holonomy character picture
One page: profinite structure group ⟹ totally disconnected holonomy ⟹ emergent periods in 2πκℤ — the only viable Wallstrom branch per the #95 downgrade. Establishes whether the selection principle is separable from the holonomy construction.

> **Gate G-2 (stop/go, week 10):**
> - *DT1 + DT2 both hold* → the dilemma is theorem-shaped on both horns; the survivor hunt (P2) proceeds as the sole decisive object, with sharpened constraints. **Expected outcome.**
> - *DT1 fails constructively* (a unitary channel exists on the shift) → that counterexample is the program's cornerstone-in-embryo; abandon the scan, attack the construction immediately.
> - *Either horn blocks the write-up* (hidden assumptions surface) → the dilemma downgrades from fatal-scenario to open-problem; record honestly; P2 still runs.

---

## PHASE 2 — THE ELEPHANT SCAN (MONTHS 2–6)

Merged WP2/WP11. The scan is the program's computational engine and its scaling audit, now with a real object.

### 2.1 The family
Skew product on ℤ₂ × ℤ₃ (multiplicatively independent scalings — the transversality requirement), coupling λ interpolating shift-horn ↔ odometer-horn, finite depth L as second knob. Discrete time (G0-forced).

### 2.2 The seven steps, with pre-registered outputs

| Step | Computation | Output artifact |
|---|---|---|
| S1 | Implement F_λ; verify ergodicity classes across λ | Phase diagram v0 |
| S2 | **Round-14 reinstated criterion:** locate zones of mixing ∧ singular-continuous coexistence (the Banach zone) | Spectral-type-vs-λ curve |
| S3 | Equilibrium measure; equidistribution of near-defined states | Born-marginal fidelity vs. λ (no fitting) |
| S4 | Candidate projections; marginal regularity classification per λ | **AC vs singular zones map — this IS the scaling audit (#97's N=3 corner expected inside the family)** |
| S5 | Fourier decay measurement; Fisher-finiteness numerics (with the H¹ rung, not just L²) | Regularity-rung table |
| S6 | Contextuality obstruction class vs. λ (Abramsky-style, finite approximations) | Cohomology-class-vs-λ curve [stretch goal] |
| S7 | κ read-off from cyclotomic unit; two-mass check | κ(λ, L) table |

Deliverable: **Paper C-preprint — "The λ-phase diagram of arithmetic skew products"** + reproducible dataset.

> **Gate G-3 (stop/go, month 6):**
> - *A zone with (regular marginal ∧ context-washing ∧ habitat-compatible spectra) exists* → P3 targets that zone; the program has its substrate candidate.
> - *No such zone* → the dilemma is now **computed** no-go evidence, not conjecture: verdict D-track with the strongest artifact the program can produce. This is a successful scientific outcome per v5 §43.
> - *Ambiguous zones* → iterate the family (third prime; coupling renormalization) for one bounded cycle, then decide.

---

## PHASE 3 — PROBABILITY THEOREMS ON THE BEST REGIME (MONTHS 4–12)

- **T-Prob-1 (restated, #96):** on the scan's best λ-zone, attack the full regularity chain: transversality ⇒ quantitative Fourier decay ⇒ density in stated Besov/Sobolev class ⇒ √ρ ∈ H¹ ⇒ finite Fisher. Rudolph–Johnson route (#115) with its three chartered obligations: (i) pushforward joint-invariance under ×2×3, (ii) physical meaning of h_μ > 0 on the substrate, (iii) the Fisher rung beyond Lebesgue.
- **T-Born (static/dynamic separated):** equidistribution of near-defined preparations toward the archimedean equilibrium measure; numerics from S3 generalized into a theorem attempt.
- **Anti-pattern discipline (from Round 15):** every convergence claim enters the ledger as `ASSUMED` until it names a theorem. No Deligne/Sato–Tate invocations outside their domain; no inserted couplings.

> **Gate G-4:** T-Prob-1 proved on one zone → the Measure Gap closes constructively; π work becomes properly posed. Blocked → identify the exact missing rung; that rung becomes the program's named frontier.

---

## PHASE 4 — π AND THE AMPLITUDE (MONTHS 6–18; THE OPEN FRONTIER)

- **Spec-first (#94):** fix the amplitude specification (S1–S4: determinism; |ψ|² = ρ and S = κ arg ψ; no-go compliance; Stone-linear output with generator −(κ²/2m)Δ + V) *before* constructing π. Proposition L guarantees everything downstream of S4-form.
- **Shape formalization (WP5):** Shape A (composite operator), B (adelic connection — with the #113 product-formula consistency obligation), C (sheaf/obstruction) formalized; compatibility/complementarity determined; **one explicit π candidate required** per v5 §41.4.
- **Wallstrom via Ẑ-holonomy (#116)** as the leading phase route; cyclotomic character-limit selection principle as Gate-4 content.
- **Kill discipline:** every candidate amplitude channel is tested against the typed exclusion set (not conditional average, not factor, not linear channel, not diffusive corner, spectrally active) *before* investment.

**Honesty clause, pre-registered:** this phase may not converge in 18 months. At month 18: formal verdict review against outcomes A/B/C/D (v5 §43), with the honest limbo — "habitat nonempty, channel unconstructed" — recorded as verdict C with an address, not failure.

---

## PHASE 5+ — DOWNSTREAM GATES (CONDITIONAL ON P4)

Only after an amplitude candidate survives P4's typed tests:
- **Gate C** (WP7): emergent tensor product with error bound f(g) → 0; weak-coupling factorization; interaction ⟹ entanglement; no-signaling.
- **Gate 4** (WP6/WP8): κ universality, two-mass test, Wood–Speckens exact-setting distribution (as *derived* distribution, never by-hand measurement dependence — the #111 lesson).
- **WP9–WP10:** Q derivation (G2 machinery now applies), guidance, equivariance closure; then and only then P1/P2 fork, a₀ target, CMB/BBN/PBH handles under the no-per-row-fitting rule.

---

## GOVERNANCE: ROLES, CADENCE, PUBLICATION

| Role | Who/what | Function |
|---|---|---|
| Arbiter & scheduler | You | Ticks, constraint register, CHALLENGED arbitration, courier for external packets |
| Deterministic core | Sequencer + oracle code | One pass per tick; verdict taxonomy; ledger propagation |
| Checker-CAS / Checker-CITE | POC agents | Exit-code evidence; source verification (Batch 2 ongoing) |
| Attacker | POC agent + **external AIs via packets only** | Adversarial passes; §37.9: all external output auto-tagged `ASSUMED` pending audit |
| Human authorship | You (+ collaborators) | All theorem statements and proofs; Papers A/B/C |

**Publication milestones:** Paper A (month ~2), G0 note (month ~3), Scan paper (month ~6), T-Prob-1 (when proved — the credibility cornerstone).

**Resources:** CAS environment (sympy/Python — trivial); scan compute (modest: 2-adic × 3-adic symbolic dynamics is laptop-to-small-cluster scale); the scarce resource is theorem-proving attention, which is why the plan front-loads write-up and short proofs.

---

## RISK REGISTER

| Risk | Likelihood | Mitigation |
|---|---|---|
| Record corruption by enthusiastic formalism (the Round-15 failure mode) | High | §37.9 rule; verifier POC as gatekeeper; named-object-density flag |
| Scan finds only ambiguous zones | Medium | Bounded iteration (one family extension), then honest verdict |
| DT1/DT2 harder than shaped | Medium | They are write-up-shaped on the canonical cases; if genuinely hard, the difficulty itself is diagnostic — record it |
| Scope creep into cosmology/dark sector | Medium | WP10 hard-gated behind P4 convergence (v5 §41.9) |
| Single-contributor bottleneck | Real | Papers A + G0 note front-loaded precisely because they need no research risk |

---

## LEDGER DELTA (#118–123)

| # | Item | Status |
|---|---|---|
| 118 | v5.1 repair set (#91–105) as tracked diff | Charter |
| 119 | G-1…G-5 gate structure with pre-registered kill criteria | Charter |
| 120 | Scan S4 regularity map = the scaling audit, now with object | Charter (implements #97) |
| 121 | Paper A as first publication artifact | Charter |
| 122 | Phase-4 honesty clause: 18-month verdict review pre-registered | Charter |
| 123 | External-AI packet protocol (§37.9 operationalized) | Charter |

---

## THE FIRST 30 DAYS, CONCRETELY

1. **Week 1:** v5.1 diff; ledger schema live; CAS sandbox up; C1–C2 checks run.
2. **Week 2:** C3 (shear) + C4 (G2) checks — *Gate G-1 verdict*; Batch 2 citation audit begins.
3. **Week 3:** G0 note drafted; DT1 statement finalized with exact hypotheses.
4. **Week 4:** F_λ implemented (S1); ergodicity spot-checks; Paper A introduction + G2 section complete.

---

## CLOSING STATEMENT

This plan embodies the program's central lesson across fifteen rounds: **the decisive moves are cheap and the expensive moves are only worth funding after the cheap ones fire.** Write up what is proved (Papers A, G0); prove what is provable in weeks (DT1/DT2, the Ẑ-lemma); compute what is computable (the scan, which is simultaneously the scaling audit the program never had an object for); and only then spend the deep research budget (T-Prob-1, π) against the exact regime the scan selects. Every phase is designed so that its failure mode is itself a first-class result — because in this program, as the frozen record states and this plan operationalizes:

$$\boxed{\text{A decisive null is a successful research result — and the audit procedure, run relentlessly, is the product.}}$$
