The consolidated tripartite architecture follows as a downloadable file. Adjudication rule applied throughout: a feature survives only if it is an axiom, a regime statement forced by an axiom, or a debt item with a named claim and kill condition — everything else was subtracted, with the cuts recorded in §13.

```markdown
# BM + IST + AS Tripartite Synthesis — v6.0

**File:** `bm_ist_as_synthesis_v6_0.md` (supersedes `bm_ist_synthesis_v5_1.md` and all addenda)
**Design rule:** optimal design = no remaining feature to subtract. A feature earns its place only as (i) an axiom, (ii) a regime statement derived from an axiom, or (iii) debt with a named claim, dependency, and kill condition.
**Status tags:** [AXIOM] / [REGIME] / [DERIVED-TARGET] / [SPECULATIVE-SHAPE] / [REJECTED-RECORDED]
**Consolidates:** net-valid revisions from Z, DeepSeek, and Claude (IST/E8 audit, BM/AS first-principles audits, Valentini integration, eval framework). Disagreements adjudicated in §13.

---

## 0. Identity (one paragraph)

One deterministic arithmetic substrate (revised IST) whose only legal environment is itself — the internal-environment principle. Integrating out the non-archimedean places via a functional-RG trace (AS machinery) yields a stochastic, dissipative reduced flow on the archimedean shadow: drift = Bohmian guidance (Kramers form), noise = derived colored noise, invariant measure = |Ψ|²(1+γδρ_ξ). The flow's UV junction — where the bath coupling runs strong — is the asymptotic-safety fixed point of the gravitational sector. Chance is ignorance of the traced-out arithmetic; determinism holds at the adelic bottom only; the quantum formalism is the decoupled-bath limit. Three former theories, three regimes, one running scale.

---

## 1. Axioms (four; count unchanged from v5.1)

- **P1′ [AXIOM] — Adelic origin of dissipation.** Unitary adelic dynamics; environment = the non-archimedean places; EAA/FRG integration with running bath coupling γ(k); output: reduced archimedean flow = drift + colored noise + memory kernel, **Kramers form** (positions + momenta), decoupling scale ξ₀. Consistency condition: **Ostrowski product formula** constrains the β-function (cross-place balance is a ledger law, not an option).
- **R [AXIOM, gated] — Recovery.** Above ξ₀ the shadow limit reproduces SM + GR exactly. Gates: precision null (no adelic corrections beyond the correction field), diffeomorphism-invariance recovery, Euclidean→Lorentzian continuation.
- **P2′ [AXIOM] — Mixing, not hyperbolicity.** Dissipativity + nondegenerate noise ⇒ unique exponentially mixing invariant measure μ* on the random attractor I_U (Hairer–Mattingly / Crauel–Flandoli class; theorems imported, not hoped for). c = anomaly parameter of the measure's fine structure, computable two ways (Lyapunov spectrum; γ-flow critical exponents) — cross-validation is a consistency gate.
- **P3′ [AXIOM] — Coarse-grained independence.** Cross-place correlations exist only through the rational skeleton (ℚ = ∩ₚ ℚₚ), exponentially suppressed at apparatus separation. **The suppressed error term IS the T3 observable** — Bell safety and Diophantine texture are two scales of one statement.

**Consistency requirement on P3′:** the fine structure must respect no-signaling — first moment of δρ vanishes; deviations live only in joint correlators, never marginals.

---

## 2. The regime map (replaces every "fundamental" claim of the three components)

| Band | What holds | Owner |
|---|---|---|
| **UV** (k ≳ junction) | Arithmetic substrate; finite information; lattice-organized code; no Born rule; no continuum; no unitary Schrödinger | IST (+ E8-as-code) |
| **Junction** | γ-flow fixed point; trace begins; fixed-point scale = ξ₀ = ℓ_P (check) | AS machinery |
| **Bridge** | Kramers stochastic dynamics; colored noise; memory; Valentini transient; foliation-dependence lives here | revised BM |
| **IR** (k ≪ ξ₀) | Pure Bohmian guidance; unitary Schrödinger; smooth spacetime; Born exact above coarse-graining; KAM islands = classical world | BM + R |

**Blind-Men operationalization:** each component owns exactly the regime where its theorems are strongest. No overlap, no stitching.

---

## 3. Derived-regime statements (standing where postulates used to be)

- **D1 — Asymptotic Bohmianity.** The guidance law is the IR fixed point of a running guidance law; deviations scale like (k/k_ξ), the same regime carrying the T3 texture. Determinism is regime-true (thermodynamic-limit object). Completes v5.1's §1.2 program: EA-EOM "replacing" the guiding equation is finished as *derivation*.
- **D2 — Asymptotic unitarity.** Lindblad-type dissipative evolution above k_ξ; unitary Schrödinger below. Mechanism target (why linear): Gaussian-IR fixed point route [DERIVED-TARGET].
- **D3 — Smoothness above the crossover.** Spectral dimension runs 4→2 because geometry crosses from archimedean to ultrametric; the elephant datum (AS/LQG/causal sets/holography) becomes this theory's signature. Three-way discrimination: flat 4 (skeleton wrong), 2 (recurrence-marginal — observed), 4/3 (Alexander–Orbach tree).
- **D4 — The foliation is the trace direction.** BM's foliation, Palmer's atemporal set, and AS's RG time are one arrow. L4 upgraded: at μ* the flow is stationary — foliation-dependence exists only in transients, which are *the same object* as Born-rule deviations (decaying at the same spectral-gap rate). One experimental program: non-equilibrium searches = Lorentz-violation searches.
- **D5 — Born rule.** Destination structural (unique invariant measure = |Ψ|²(1+γδρ_ξ)); path Valentini (relaxation); rate = the generator's spectral gap (critical-exponent class, scheme-gated). Corrections: deviations from exact guidance and from Born statistics are the same transient.
- **D6 — Pointer stability.** Corollary II restated for the Kramers flow: macroscopic branches = stable stationary manifolds; empty branches stay empty by theorem, not fiat.
- **D7 — KAM coexistence [DERIVED-TARGET].** Attractor and stable islands coexist; the classical world is the islands. Named toy-model computation (closes the eval turn's uncovered invariant).
- **D8 — Nonlocality retained in the IR.** The conservative fork stands. DeepSeek's "nonlocality = IR artifact of UV locality" is **rejected**: it re-imports the superdeterminist route ablated in v5.1. The substrate's locality status is *not asserted*.

---

## 4. Component revisions, consolidated

### 4.1 Bohmian Mechanics (revised)

| v5.1 assumption | Disposition |
|---|---|
| First-order, positions-only guidance | **REGIME-STRUCTURED** → Kramers form (B2); momentum beables restored; overdamped = strong-friction limit |
| ρ = \|Ψ\|² postulate (equilibrium) | **SUBTRACTED** → unique invariant measure (B3); Valentini subsumed |
| Unitary Schrödinger, always | **REGIME-STRUCTURED** → asymptotic unitarity (D2) |
| Preferred foliation by hand | **DERIVED** → trace direction (D4) |
| Markovian, memoryless | **SUBTRACTED** → non-Markovian fundamental; memory kernel = skeleton's response function (B5) |
| Fixed-N configuration space | **OPERATIVE FORM** → Bell-type QFT (DGTZ) as the working formulation; fixed-N = sector description |
| Continuous configuration space | **DERIVED-TARGET** → emergent; extends the Projection Problem (debt PCS-1) |
| Classical limit by decoherence | **UPGRADED** → KAM obligation (D7) |
| Holistic contextuality of guidance | **PRESERVED, FIXED TARGET** — the emergent construction is answerable to it |
| Quantum potential as brute fact | **DERIVED-TARGET** → IR effective-action kinetic term (Arakelov-height program on ledger) |

### 4.2 Invariant Set Theory (revised)

| Palmer assumption | Disposition |
|---|---|
| Geometric fractality (positive codimension set) | **REGIME-STRUCTURED** (R1): above ξ₀ full-dimensional (P2′ forces it); fractality survives only as sub-ξ₀ measure texture; c = anomaly, not codimension. Named computation: is the texture multifractal or smooth-with-arithmetic-structure |
| Violated measurement independence | **SUBTRACTED** → P3′ coarse-grained independence; superdeterminism stays dead |
| Gravity as selector | **SUBTRACTED** (v5.1) — unchanged |
| Epistemic wave function | **SUBTRACTED** (v5.1) — Ψ_∞ = archimedean shadow of the adelic state (derived, per P1′) |
| p-adic metric posited | **DERIVED-TARGET** → Markov-partition alphabet (IST-ALPHABET-MARKOV-1); branching arity = output |
| Niven definability cut | **REPLACED-TARGET** → preperiodic locus (IST-DEFCUT-HEIGHT-1); Niven-fit obligation |
| L hand-set (~10¹⁰⁰/10¹²⁰) | **DERIVED-TARGET** → c-function along the RG trajectory (IST-L-CFUNC-1) |
| 1-bit/Planck-time rate | **DERIVED-TARGET** → RG coarse-graining rate (IST-RATE-RGMATCH-1); borrowed-rigor flagged |
| Causal partial order, atemporal block | **DEMOTED** → SR causality emergent above ξ₀; atemporal formulation on the parallel track |
| Chaotic-cosmology intuition | **UPGRADED** → theorem-class (random dynamical systems shelf) |
| Internal-environment principle | **NOW EXPLICIT** — the foundational justification of the whole trace; minimal dissipative completion of BM |

### 4.3 Asymptotic Safety (revised)

| AS assumption | Disposition |
|---|---|
| EAA / Wetterich machinery | **ADOPTED** — P1′'s toolkit (LPA, derivative expansion) |
| Non-Gaussian UV fixed point | **ADOPTED as hypothesis for the γ-flow**; scheme-independence gate mandatory; cross-validation with c is the consistency proof |
| Fixed point ↔ I_U as one object | **CONDITIONAL two-projections** (Z): projections until proven, identity never assumed; DeepSeek's "critical surface of one universal flow" downgraded to parallel track |
| Metric fundamental at all scales | **REGIME-STRUCTURED** (A3/D3): continuum is the shadow's IR description; the continuum/discrete contradiction is resolved by regime structuring, and named as the synthesis's single largest coherence risk |
| Interpretation-free minimalism | **REFUSED** (A2): a fixed point without a state-space shadow is a calculator, not a world |
| Background-field method | **LEGITIMATE IN THE SHADOW** (A4); diffeo recovery is a gate on R |
| Universality of exponents | **ARITHMETIC-ANCHORED** (A5): scheme-independence *and* Ostrowski balance |
| RG arrow as bookkeeping | **PHYSICAL** (A6): the trace direction is the thermodynamic arrow (unifies D4) |
| Euclidean evidence base | **HARD GATE**: E→L continuation with real-time unitarity is a prerequisite for BM trajectories to be well-defined (Claude §4.6) |
| AS-with-matter | **LIVE DEPENDENCY**: compatibility of the NGFP with the emergent matter sector is unaddressed — named debt, not assumed |
| Eval-silence | **REFUSED** — T1/T2/T3 rows, kill test, death conditions stay front and center |

---

## 5. Probability: one mechanism, one rate, two regimes (resolves the three-mechanism sprawl)

- **One destination:** μ* = |Ψ|²(1+γδρ_ξ) — structural selection.
- **One rate object:** the spectral gap of the relaxation generator (computable from kernel A1; critical-exponent class).
- **Two regime appearances:** thermalized regimes — T-Born (equidistribution; exponential convergence, no residue); transient/relic regimes — Valentini-H relaxation (power-law; incomplete for low-complexity/symmetric sectors — sector-dependence stated, not assumed away).
- **Typicality:** absorbed — it is the measure-theoretic restatement of the same invariant measure.
- **ξ₀ = Valentini's coarse-graining scale** — his standing free parameter fixed by our structure (the integration's best jam point).
- **Initial-condition obligation (debt IC-1):** relaxation needs an off-measure start inside the basin; the necessity theorem constrains how gross the start can be. Specify basin + transient typicality.
- **The fork (cheapest decisive experiment in the theory):** plain BM+Valentini says sustained laboratory non-equilibrium is preparable; the tripartite says the attractor forbids it. One tabletop experiment, two theories, one survivor.
- **Relic non-equilibrium:** inherited as a citable, falsifiable row — redirected to skeleton texture (small amplitude, computable shape), not gross CMB deviations.

---

## 6. E8 fold (regime-mapped, below the chirality wall)

- **The inversion:** not E8-as-unification (dead: Distler–Garibaldi; chirality is the wall and stays walled) but **E8-as-organization** — the lattice structuring the bath's internal code and the skeleton's fine structure, in the unique regime where packing optimality, error correction, and modularity coincide and nothing can be fitted (even unimodular uniqueness in d = 8).
- **Ledger:** E8-PACKING-INFO-1 (code selection via entropy saturation — blocked on Phase 2); E8-CFT-CHIRALITY-1 **REJECTED-RECORDED** (holomorphic lattice VOA single-module rigidity forecloses it; do not re-attempt without a non-holomorphic/twisted variant as a new claim); E8-CFT-CHIRALITY-2 (narrow surviving form: heterotic-style split — gauge sector from E8, chirality from a separate topological structure — doubly conditional on the spectral-dimension gate and lattice-type UV).
- **Standing caution:** 248 is a generator count, not a particle census; never cite E8 as "how many particles exist."
- **Regimes:** bath error-correction restatement of D6 (testable: lattice-organized noise fine structure); gauge-sector shape conjecture (residual symmetry of the traced bath = shadow's internal symmetries — enters only through the no-go audit); κ-dimension gate (E8 enters only if the internal γ-flow fixed point is 8-dimensional and scheme-independent); **modular fingerprint** (E₄-class q-coefficient statistics in the deviation spectrum — the held-out T3 shape and the fold's falsifier); three-way spectral dimension; arithmetic-group note (E₈(ℤ)).

---

## 7. The scale factory (one running scale k)

| Decoupling event | Scale | Status |
|---|---|---|
| Bath drops out | ξ₀ | P1′ output (computed, not assumed) |
| Gravity strong / junction | ℓ_P | γ-flow fixed point; **check ξ₀ = ℓ_P** |
| Skeleton resolves for Bell angles | Q | A2 map (kernel → denominator cutoff) |
| Observation coarse-grains | ε₀ | Instrument-defined |

**ξ₀–G relation:** the Ostrowski balance point is the candidate mechanism (feeds owed item; testable the moment the adelic Caldeira–Leggett closed form exists).

---

## 8. Gravity scope

- AS governs the gravitational sector's UV — which is the archimedean residue. This is Claude's conservative fix (§4.2) merged with Z's conditional two-projections: no "one flow generates everything" claim.
- T_frac and archimedean selection (v5.1 §5/§7) remain future work; prohibited from supporting core rows.
- Emergent-continuum route (Jacobson/Verlinde class) = parallel track option [SPECULATIVE-SHAPE]; if adopted later, NGFP is downstream of the substrate and residual background-dependence is expected rather than alarming.

---

## 9. Gates and kill conditions (consolidated)

| # | Gate / row | Tier | Kill condition |
|---|---|---|---|
| 1 | Parameter window: existing interferometry (~10⁴ amu) + Bell precision (10⁻⁴–10⁻⁵) | T1 (kill test) | No consistent window → dead now |
| 2 | No collapse-noise *shape* | T2 | Residual noise white/spontaneous (GRW/CSL form) instead of derived colored spectrum |
| 3 | Mesoscopic ceiling scaling | T2 | Ceiling scales with G/mass, not N |
| 4 | Diophantine angular statistics | T3 | Continued-fraction law absent |
| 5 | Modular fingerprint (E₈ fold) | T3 | E₄-class q-statistics absent |
| 6 | Spectral dimension running | Gate | Flat 4 → skeleton wrong |
| 7 | ξ₀ = ℓ_P | Consistency | Black-hole area quantum fails |
| 8 | Scheme-independence | Gate | Any derived constant moves with regulator |
| 9 | Clustering transfer | Theorem-target | No exponential decay → no mass-gap bridge |
| 10 | Kramers limit from the trace | Theorem-target | No friction-dominated limit → dissipation story wrong |
| 11 | Diffeo recovery above ξ₀ | Gate on R | Fails → shadow inconsistent |
| 12 | E→L continuation | Gate on R | Fails → BM trajectories undefined |
| 13 | Ostrowski balance in the closed form | Consistency | Violated → R3 dead |
| 14 | Sustained lab non-equilibrium | T1 fork | Prepared → structural selection dead |
| 15 | No-signaling (equilibrium AND transients) | Safety | Any marginal deviation → dead on arrival |
| 16 | Niven fit of the preperiodic locus | Faithfulness | Divergence unreported → contamination |
| 17 | AS-matter compatibility | Dependency | NGFP incompatible with emergent matter → junction broken |
| 18 | Internal Occam control | Meta | T3 nulls accumulate → AS-alone (cheaper theory) wins |

---

## 10. Debt ledger (new and carried; IDs per addendum convention)

| ID | Item | Blocks on |
|---|---|---|
| IST-L-CFUNC-1 | L as RG c-function | Flow-equation construction |
| IST-DEFCUT-HEIGHT-1 | Preperiodic locus replaces Niven cut + Niven fit | Phase 2 |
| IST-ALPHABET-MARKOV-1 | Alphabet from Markov partition (arity = output) | Prerequisite to Phase 1 |
| IST-RATE-RGMATCH-1 | Collapse rate = RG rate | Flow equation; rate-match computation |
| E8-PACKING-INFO-1 | E8 lattice as code (entropy-saturation selection) | Phase 2 |
| E8-CFT-CHIRALITY-2 | Narrow heterotic-style split | Spectral-dimension gate (needs own CLAIM_ID) |
| PCS-1 | Configuration space itself emergent (Projection Problem extended) | Phase 3 |
| IC-1 | Initial-condition obligation (basin, transient typicality) | Phase 2 |
| ADR-1 | Quantum potential via Arakelov heights | Phase 3 |
| D2-MECH | Linearity from Gaussian-IR fixed point | Phase 4 |
| GATE-EL | Euclidean→Lorentzian continuation | Phase 4 |
| GATE-ASM | AS-with-matter compatibility | Phase 4 |
| FOL-1 | D4 mechanism (foliation-inertness at μ*) | Phase 2 |
| carried v5.1 items 1–25 | Kernel A1, δρ_ξ, fractal equivariance, stationarity of jumps, N_c formula, ξ₀→Q map, consistency computation, toy model, T_frac, archimedean selection, lift problem, UV-lattice theorem | As per v5.1 phasing |

---

## 11. Program (ordered by cost and information gain)

0. **Kill test first** (weeks; assumption-tagged A1/A2): window scan; re-run on every new number.
1. **Substrate + scaling audit; Markov-partition check** (alphabet arity as output).
2. **Measure & probability:** entropy saturation, IC-1, FOL-1, ξ₀-as-coarse-graining check.
3. **π construction, one qubit** (PCS-1; holistic contextuality as the fixed target).
4. **EAA flow computation** (γ(k), Ostrowski check, D2 mechanism, E→L, AS-matter).
5. **Gates 4–9; cross-validation of c** (Lyapunov vs critical exponents; agreement = universality evidence, disagreement = inconsistency gate).
6. **Parallel tracks:** necessity theorem (boundary-value over cosmic history); emergent-continuum option; atemporal formulation.

---

## 12. Sharpened questions

- What is a fixed point *for* — at which scale does the shadow's self-consistency crystallize?
- Is the Born rule a universality-class phenomenon — |Ψ|² as the critical phenomenon of an arithmetic substrate?
- Is dimension arithmetic — the 4→2 running as high-resolution emergence of the rational skeleton?
- Is noise derived or fundamental — fluctuation-dissipation shape versus postulated spectrum?
- Why linearity — is the Gaussian-IR route real?
- Why these relevant directions — does the count of relevant directions equal the count of physical parameters (mass, γ, ξ₀, L)?

---

## 13. Subtraction record (design-efficiency audit; disagreements adjudicated)

**Subtracted this round:**
- ρ = |Ψ|² as postulate (→ D5); memorylessness (→ B5); overdamped first-order form (→ Kramers); hand-set L, Niven cut, posited p-adic metric, fixed collapse rate (→ four ledger claims); interpretation-free silence (refused); "one universal flow" boldness (→ conditional two-projections); metric fundamentality (→ regime-true); AS eval-silence (refused); three parallel Born mechanisms (→ one destination, one rate, two regimes).

**Rejected with reasons on file:**
- DeepSeek's "nonlocality = IR artifact of UV locality" — re-imports the ablated superdeterminist fork; the conservative fork stands (D8).
- DeepSeek's "I_U = critical surface of one universal flow" as assertion — downgraded to parallel track (continuum/discrete contradiction; Claude §4.1 honesty flag adopted).
- E8-CFT-CHIRALITY-1 (naive) — single-module rigidity; recorded rejection per ledger policy.
- "248 particles" framing — standing caution.

**Preserved deliberately:** conservative Bell fork; holistic contextuality; T1/T2/T3 discipline; dependency tagging; kill-test-first sequencing; Blind-Men partiality as methodology.

**Subtraction test:** every remaining item is an axiom (4), a derived-regime statement (8), a gated debt row, or an empirical handle. No feature remains whose removal would not change a prediction, a gate, or an axiom.

---

## 14. Status note

The tripartite is now one running theory with a regime map, not three tenants: **UV is arithmetically organized, the bridge is the Valentini transient on a Kramers flow, the IR is Bohmian and unitary.** Axiom count unchanged (4) across three first-principles audits; postulate count reduced by seven; named computations +9; kill conditions +12; one component-revision rejected with reasons on file. The crux remains P1′'s adelic flow computation — the long campaign — but every ladder now borrowed doubles as a bridge, and the cheapest gate (the kill test) still costs weeks. Eval-sound by construction remains available the moment Gate 1 runs; truth is earned only at the gates.
```

**Three adjudications worth flagging outside the file:** (1) DeepSeek's most ambitious claims (universal single flow, nonlocality-as-artifact) were cut as assertion-class or fork-violating — their defensible cores (regime structuring, Kramers, mass-as-eigenvalue) were retained; (2) Claude's hard-nosed honesty items (E→L gate, AS-matter dependency, continuum-contradiction naming) were all promoted to gates — they cost nothing and catch the largest risks; (3) the probability sprawl is now closed with one destination, one rate object, two regime appearances — the thread's oldest discipline (no unexamined parallel mechanisms) finally applied to its own largest family.