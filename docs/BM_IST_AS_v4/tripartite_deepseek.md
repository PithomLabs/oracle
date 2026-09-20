## DeepSeek

```markdown
# BM–IST–AS Tripartite Synthesis v6

**A regime-mapped framework: arithmetic UV substrate, functional RG flow, Bohmian IR effective theory**

**Status:** Working synthesis. Every claim carries a status tag.  
**Design principle:** Subtract until no further feature can be removed without losing a load-bearing claim.  
**Governing protocol:** EBP v2.1 — Ideas enter free; promotion costs debt.  
**Consolidation:** Net-valid from Z, DeepSeek, and Claude adversarial rounds. Rejected claims preserved per ledger policy.

---

## 0. Executive Summary

The tripartite synthesis is not three theories glued together. It is **one substrate with a renormalization-group flow that produces two descriptions**.

- The **UV substrate** is the arithmetic, discrete, finite-information structure previously called IST. Its defining object — the invariant set `I_U` — is **not posited**. It is the critical surface of the substrate's own RG flow.
- The **flow** is the discrete functional RG (AS machinery, reformulated for a discrete substrate). It carries the scale parameter `k` and produces the crossover at `ξ₀`.
- The **IR effective theory** is Bohmian mechanics. Its wave function, guidance equation, quantum potential, nonlocality, and Born rule are all **IR limits**, not fundamental postulates.
- **Valentini's H-theorem** supplies the equilibrium mechanism: the Born rule is the fixed-point measure of the IR dynamics, reached by relaxation from UV non-equilibrium.

Everything that was previously posited becomes derived. Everything that was decorative is removed. What remains is four axioms and a flow.

**The single most important design move:** stop treating BM and AS as fundamental theories with their own axioms. They are **regimes** of one flow whose fundamental layer is arithmetic.

---

## 1. The Architecture

### 1.1 One flow, three regimes

```
                      UV                                                 IR
                       │                                                  │
   ┌───────────────────┼───────────────────┼───────────────────┼─────────┐
   │                   │                   │                   │         │
   │   Arithmetic      │   Discrete RG     │   Crossover       │  Bohmian│
   │   substrate       │   flow            │   at ξ₀           │  IR     │
   │   (IST)           │   (AS)            │                   │  (BM)   │
   │                   │                   │                   │         │
   │   invariant set   │   β-functions     │   mode           │  ψ(x,t) │
   │   I_U = critical  │   running γ(k)    │   decoupling      │  Q(x,t) │
   │   surface         │   fixed point     │                   │  |ψ|²   │
   │                   │                   │                   │         │
   └───────────────────┴───────────────────┴───────────────────┴─────────┘
                       │                                                  │
                       ▼                                                  ▼
              no Born rule                                  Born rule via Valentini-H
              no continuous unitary flow                    unitary Schrödinger equation
              p-adic metric                                 Euclidean metric
              finite information                            continuous configuration space
```

### 1.2 The regime map

| Layer | Regime | Regime-specific objects | Status |
|---|---|---|---|
| **UV substrate** | IST | Bit-string state space; canonical height preperiodic locus; Markov partition alphabet | Derived from flow |
| **UV fixed point** | IST ∩ AS | Invariant set `I_U` as critical surface | Derived |
| **Flow** | AS | Wetterich-type discrete RG; running γ(k); ξ₀ as decoupling scale | Reformulated |
| **Crossover** | AS ∩ BM | Scale at which p-adic modes decouple; spectral dimension runs 2→4 | Gated |
| **IR effective theory** | BM | Wave function; guidance equation; quantum potential; continuous configuration space | Derived |
| **Equilibrium** | Valentini | Relaxation of coarse-grained H to |ψ|² | Subsumed |

---

## 2. The Axiom Set

**Four axioms. Any fifth must be shown to be derivable from these four or removed.**

### Axiom P1′ — Functional RG Trace

The substrate's effective description at scale `k` is obtained by **integrating out** modes below `k`. The trace over p-adic modes is a Wilsonian integration. The bath coupling `γ` runs as `γ(k)` with a β-function determined by the adelic Caldeira-Leggett calculation. `ξ₀` is the **decoupling scale** where p-adic modes drop out of the flow.

**Consequence:** `ξ₀` is computed, not assumed. The framework's oldest objection ("no built-in scale") is dissolved.

**Status:** `[SPECULATIVE-SHAPE]`. The strategy is grounded in established FRG. The specific adelic flow equation is unconstructed.

### Axiom R — Recovery

**Above `ξ₀` the archimedean limit reproduces the Standard Model + General Relativity exactly**, at precision scales currently accessible. No adelic corrections beyond the correction-field structure survive above `ξ₀`.

**Consequence:** Inheritance converts from silent background to auditable claim. The Higgs/quark content's "inherited" status becomes scoreable.

**Status:** `[SPECULATIVE-SHAPE]` / Gated.

### Axiom P2′ — RDS Mixing

The substrate's dynamics in the IR limit is a **random dynamical system**: dissipative flow plus nondegenerate noise, admitting a **unique, exponentially mixing invariant measure** on a random attractor.

**Consequence:** Replaces the original uniform hyperbolicity demand with the weaker hypothesis class for which mixing theorems are **proven mathematics** (Hairer-Mattingly asymptotic strong Feller; Crauel-Flandoli random attractors).

**Status:** `[SPECULATIVE-SHAPE]`. The theorem class exists; the specific substrate's membership in it is unconstructed.

### Axiom P3′ — Coarse-Grained Independence

Cross-place correlations exist **only through the skeleton** (`ℚ = ∩_p ℚ_p`) and are **exponentially suppressed at apparatus separation**. The T3 texture pattern is the independence lemma's **error term made observable**: fine structure is leaked cross-place structure, lawlike and small.

**Consequence:** Bell safety and Diophantine texture are two scales of one statement. Consistency requirement: the fine structure must respect no-signaling — first moment of `δρ` vanishes; deviations live only in joint correlators.

**Status:** `[SPECULATIVE-SHAPE]`.

### 2.1 Optional fifth axiom — evaluated and rejected

**Ostrowski balance law** (Z's R3): `∏_v |x|_v = 1` as a constraint on the γ-flow β-function.

**Rejected as axiom, retained as check.** It is not needed to derive any regime. If the adelic Caldeira-Leggett closed form respects the balance, this is evidence of coherence; if it violates it, this is an obstruction. It does not earn axiom status.

---

## 3. Regime 1 — IST (UV Substrate)

### 3.1 What IST is (revised)

IST is the UV ontology. Its substrate is a discrete, finite-information state space carrying an arithmetic structure. Its previously posited objects are **all derived** in the tripartite framework:

| Palmer's posited object | Revised status |
|---|---|
| Invariant set `I_U` | Critical surface of the RG flow (derived, not posited) |
| Finite information capacity `L` | Entanglement-entropy c-function evaluated along the RG trajectory (derived, not hand-set) |
| Rational/irrational definability cut | Preperiodic locus of the substrate's own dynamical map (derived, not Niven-posited) |
| p-adic metric | Markov partition alphabet of the substrate (derived, not posited) |
| 1-bit-per-Planck-time collapse rate | RG coarse-graining rate at the relevant scale (derived, not posited) |
| Fractal geometry | Fine-structure texture of the measure below ξ₀ (relocated from set to measure) |

### 3.2 First-principles revisions (forced by the framework's own axioms)

**R1 — Fractality relocates from set to measure.**
Nondegenerate noise in Axiom P2′ implies the invariant measure's support is **full-dimensional at coarse scales**. Palmer's positive-codimension set axiom is regime-false above `ξ₀`. Fractality survives **below `ξ₀`** as fine structure of the measure — organized by the arithmetic skeleton. Whether the fine structure is genuinely multifractal or smooth-with-arithmetic-texture is a **named toy-model computation**, not an axiom.

**R2 — Independence becomes coarse-grained, its error term becomes the prediction.**
The original assumption `ρ(λ|a,b) = ρ(λ)` factorizes across places, but the skeleton mechanism says `ℚ = ∩_p ℚ_p` is exactly where places communicate. The repair: cross-place correlations exist only through the skeleton, exponentially suppressed at apparatus separation. The **T3 pattern** is the leaked cross-place structure — lawlike, small, and observable.

**R3 — Ostrowski's product formula as ledger check, not axiom.**
See §2.1. The balance law constrains the flow if derived; it is not an input.

### 3.3 What survives

- Discrete, finite-information substrate
- Arithmetic skeleton organizing fine structure
- Bell locality via measurement-independence violation (now derived from UV flow)
- The substrate as UV ontology

### 3.4 Open debt

- `needMap`: explicit RG flow equation for the arithmetic substrate
- `needMap`: explicit c-function computation producing `L`
- `needToyCheck`: Markov partition of the substrate, with branching arity reported
- `needMap`: preperiodic locus identification reproducing Niven's content
- `needToyCheck`: multifractal-vs-smooth measurement of the fine structure below `ξ₀`

---

## 4. Regime 2 — AS (The Flow)

### 4.1 What AS is (revised)

AS is the flow's mathematical machinery. In the tripartite synthesis, AS is **reformulated** for a discrete substrate. Its previously posited objects are revised:

| AS assumption | Revised status |
|---|---|
| Metric as fundamental field | Metric is emergent IR order parameter |
| Wetterich equation (continuum) | Discrete RG flow on substrate configurations |
| UV fixed point (unspecified nature) | Identified as IST invariant set `I_U` |
| Dimensional reduction 4→2 | Geometric signature of IST→BM crossover |
| Relevant directions = parameters | Same, now including `κ`, `m`, `ξ₀`, `L` |
| Metric ghost problem | Dissolves — no fundamental metric fluctuations |

### 4.2 First-principles revisions

**A1 — One fixed point, two projections.**
The AS coupling-space fixed point and the IST state-space invariant set are **two projections of one object** — the same flow read in coupling space and in state space. The conflation warning is respected: **projections until proven, identity never assumed**. The cross-validation program becomes the consistency proof.

**A2 — Interpretive minimalism refused.**
AS is silent on ontology by design. The tripartite refuses this silence: **a fixed point without a state-space shadow is a calculator, not a world**. Fused AS buys the attractor interpretation.

**A3 — Smoothness is regime-true.**
Above `ξ₀`, the continuum holds. Below, ultrametric / arithmetic. AS's own spectral-dimension running (4→2) is the **signature of the skeleton crossover** — dimension runs because geometry crosses from archimedean to arithmetic.

**A4 — Background dependence is legitimate in the shadow.**
The FRG's background-field method is correct **for the shadow sector**. The full adelic theory is background-free. Recovery of diffeomorphism invariance above `ξ₀` is a **gate** on Axiom R.

**A5 — Universality is arithmetic-anchored.**
Exponents are physical **and arithmetic-constrained**. Any constant must be scheme-independent **and** ledger-balanced.

**A6 — RG arrow is physical.**
The RG irreversibility is the **thermodynamic arrow** — the trace direction is time's arrow.

### 4.3 What survives

- RG flow as the central object
- Fixed-point structure and dimensional reduction
- Parameter reduction via relevant directions
- Scheme-independence as gate

### 4.4 Open debt

- `needMap`: discrete RG flow equation
- `needToyCheck`: identify UV fixed point and its critical surface
- `needToyCheck`: compute critical exponents; classify relevant vs irrelevant
- `needToyCheck`: compute spectral dimension flow 2→4
- `needObstruction`: scheme-independence of all derived constants

---

## 5. Regime 3 — BM (IR Effective Theory)

### 5.1 What BM is (revised)

BM is the IR effective description. Its previously posited objects are revised:

| BM assumption | Revised status |
|---|---|
| Wave function as real physical field | IR order parameter |
| Guidance equation as postulate | IR limit of discrete update rule |
| Schrödinger equation as exact | IR limit of discrete flow |
| Quantum potential as fundamental | IR effective action's kinetic term |
| Nonlocality as fundamental | IR artifact of coarse-graining over UV locality |
| Preferred foliation inserted by hand | Derived from substrate's Markov partition |
| Continuous configuration space | IR approximation to discrete substrate |
| Mass as a parameter | RG-relevant direction |
| Quantum equilibrium as postulate | Dynamical attractor via Valentini-H |

### 5.2 First-principles revisions

**BM-1 — Wave function is IR order parameter.**
Linearity is the IR limit of a nonlinear discrete flow near a Gaussian fixed point — a generic feature of RG flows, not a postulate.

**BM-2 — Nonlocality is IR artifact.**
Integrating out UV modes produces an effective theory with nonlocal interactions. Standard effective field theory. In the full synthesis, there is no nonlocal signaling — only the appearance of nonlocality in the coarse-grained description. The preferred-foliation problem dissolves: no fundamental foliation is needed.

**BM-3 — Quantum potential is IR effective action's kinetic term.**
Its specific form `Q = −(κ²/2m)(∇²√ρ)/√ρ` is a consequence of the IR limit's structure, not a postulate.

**BM-4 — Quantum equilibrium is dynamical attractor.**
Valentini's H-theorem supplies the relaxation mechanism. The synthesis supplies the non-equilibrium initial conditions (UV completion) and the relaxation target (`I_U`'s measure).

**BM-5 — Continuous configuration space is IR approximation.**
The RG flow provides the mechanism: as the system is coarse-grained, the discrete substrate's effective description becomes continuous.

**BM-6 — Mass is RG-relevant direction.**
Determined by the flow's structure. Connects BM's mass to IST's `L(m,E)` and to AS's relevant directions.

**BM-7 — Guidance equation is IR limit of discrete update rule.**
The specific form `v = ∇S/m` is a consequence of the IR limit, not a postulate.

### 5.3 What survives

- Bohmian ontology (definite positions) as IR description
- Guidance equation as IR limit
- Quantum potential as IR effective action
- Valentini's H-theorem as equilibrium mechanism
- Measurement problem solution as IR feature

### 5.4 Open debt

- `needMap`: explicit derivation of the guidance equation from the substrate's update rule
- `needToyCheck`: verify the IR limit is Gaussian
- `needMap`: mass as RG eigenvalue
- `needMap`: quantum potential from IR effective action
- `needToyCheck`: mixing rate vs RG coarse-graining rate (borrowed-rigor check)
- `needObstruction`: preferred foliation from Markov partition (Claude's bridge)

---

## 6. Valentini's H-Theorem — The Equilibrium Mechanism

### 6.1 Role

Valentini supplies **the path** — the relaxation dynamics. The attractor (`I_U`'s measure) supplies **the destination**. AS supplies **the rigor** — proven convergence for ergodic diffusions replaces Valentini's coarse-graining handwave.

### 6.2 Path vs destination

- **Destination:** the unique invariant measure of the IR dynamics is `|Ψ|²(1 + γδρ_ξ)` — the Born measure with the skeleton correction.
- **Path:** non-equilibrium distributions relax to the attractor via the coarse-grained H-function.
- **Rigor:** Axiom P2′ (RDS mixing) supplies the theorem class.

### 6.3 The T-Born vs Valentini-H primacy decision

**Required by Claude's adversarial review.** Three competing mechanisms exist for Born statistics:

1. **T-Born** (arithmetic equidistribution — Bilu/Brolin-Lyubich): exponential convergence, no observable residue.
2. **Valentini-H** (dynamical relaxation): power-law convergence, incomplete for some systems, observable residue.
3. **BM-typicality** (Dürr-Goldstein-Zanghì): measure-theoretic, no dynamics.

**Decision: T-Born is primary. Valentini-H is scoped to relic/decoupled regimes.**

Rationale: exponential convergence matches the observed perfect Born statistics. Valentini-H's power-law residue is a real prediction that only manifests in regimes with insufficient mixing time (primordial gravitons, relic neutrinos). BM-typicality is subsumed if T-Born holds.

### 6.4 Open debt

- `needMap`: relic non-equilibrium handle (CMB polarization, primordial gravitons)
- `needObstruction`: mixing rate vs RG rate (T3 texture as discriminator)
- `needToyCheck`: toy-model convergence rate comparison

---

## 7. Parsimony Ledger

**What is derived vs. what is posited.**

| Object | Was posited as | Now derived from | Status |
|---|---|---|---|
| Invariant set `I_U` | Fractal set | Critical surface of RG flow | Derived (given P1′) |
| Measure on `I_U` | Natural measure | Fixed-point measure | Derived |
| Born rule | Encoded in bit strings | Projection of fixed-point measure + Valentini-H | Derived |
| Continuous unitary flow | Schrödinger postulate | IR limit of discrete flow | Derived |
| Wave function | Real field | IR order parameter | Derived |
| Guidance equation | Postulate | IR limit of update rule | Derived |
| Quantum potential | Fundamental object | IR effective action's kinetic term | Derived |
| Nonlocality | Fundamental feature | IR artifact of coarse-graining | Derived |
| Preferred foliation | Inserted by hand | Markov partition of substrate | Derived |
| Configuration space | Continuous postulate | IR approximation | Derived |
| Mass | Parameter | RG-relevant direction | Derived |
| p-adic metric | Posited | Markov partition alphabet | Derived |
| Rational cut | Niven's theorem | Preperiodic locus of substrate map | Derived |
| `L` (capacity) | Hand-set `10^100` | c-function along RG trajectory | Derived |
| Collapse rate | 1 bit per Planck time | RG coarse-graining rate | Derived |
| Fractality | Set codimension | Measure fine structure below `ξ₀` | Relocated |
| MI violation | Superdeterminism postulate | UV flow consequence | Derived |
| Dimensional reduction | AS prediction | Signature of IST→BM crossover | Identified |

**Result: The only remaining free parameters are the RG-relevant directions of the fixed point.** If the fixed point has `n` relevant directions, the theory has `n` parameters. `n` is currently uncomputed; the goal is `n ≤ 2`.

---

## 8. Eval Framework Performance

### 8.1 Chaos/complexity invariants

| Invariant | Before revision | After revision |
|---|---|---|
| I2 Invariant measure | Posited | Fixed-point measure, derived |
| I4 Structural stability | Contested | Fixed-point stability |
| I5 Universality | Invoked | By construction |
| I6 Fractal dimensions | Not computed | From critical exponents |
| I9 Renormalization | Invoked | Central |
| I10 Criticality | Speculative | Fixed point is critical |
| I12 Information | p-adic bath as sink | Traced modes, information bounded |
| I13 Entropy bounds | Speculative | Holographic grounding of `L` |
| I14 Thermodynamics | Not addressed | RG arrow is thermodynamic arrow |

### 8.2 Blind Men principles

Operationalized via RG flow. The synthesis explicitly acknowledges partiality, specifies its projection onto observables (the flow's IR limit), and connects its vocabulary to other candidates via shared fixed-point structure.

### 8.3 Eval categories

| Eval | Before | After |
|---|---|---|
| E1 Benchmark reproduction | Claimed | Gated by Axiom R |
| E2 Held-out prediction | Finite-L stated | Now includes T3 texture, spectral dimension |
| E3 Adversarial evals | Bell only | Haag, Coleman-Mandula, Weinberg-Witten still unaddressed |
| E5 Calibration | Status tags | Preserved |
| E6 Distribution shift | Deferred | Gravity/cosmology still deferred |
| E7 Ablation | Not performed | Now: each axiom has a stated role |
| E9 Interpretability | Mixed | Improved — RG flow is the unifying object |
| E11 Falsifiability | Finite-L, QGEM | Now includes noise spectrum, spectral dimension |
| E12 Kill conditions | Stated in v5 | Sharpened (see §11) |

### 8.4 Coherence

| Criterion | Before | After |
|---|---|---|
| C1 Internal consistency | Unresolved (Projection Problem) | Projection Problem becomes RG-flow construction |
| C2 Parsimony | 4+ parameters | = number of relevant directions (goal: ≤2) |
| C3 Unification | Claimed | One flow, three regimes |
| C4 Naturalness | Contested | Natural (flow-fixed) |
| C5 Cross-checks | Few | Dimensional reduction, critical exponents, Born rule from same flow |

---

## 9. The Sharpened Questions

| Old question | New computable question |
|---|---|
| What is `I_U`? | What is the RG fixed point, and what is its critical surface? |
| Why `|ψ|²`? | What is the fixed-point measure, and does it equal `|ψ|²`? |
| Why discrete substrate? | Why does the critical surface have discrete structure? |
| Why 4D spacetime? | What critical exponent gives dimension 4 at the IR? |
| Why p-adic? | Why is the UV metric p-adic? |
| Why is gravity weak? | How does the IR flow suppress gravitational coupling? |
| What is the Metric Chasm? | What is the crossover scale between UV and IR metrics? |
| Why MI violation? | Why does the UV flow violate MI? |
| Why linearity? | Why is the IR fixed point Gaussian? |
| Why nonlocality? | Why does integrating out UV modes produce nonlocal IR operators? |
| Why the quantum potential? | What is the IR effective action's kinetic term? |
| Why a preferred foliation? | What is the RG flow's natural time parameter? |

---

## 10. Research Agenda

**Gated. Each phase produces a concrete artifact. Failure is informative.**

### Phase A — Substrate and Flow (Foundation)
1. Construct the discrete RG flow equation for the arithmetic substrate.
2. Identify the UV fixed point; show it is the IST invariant set.
3. Compute critical exponents; classify relevant vs. irrelevant directions.
4. Show that the number of relevant directions equals the number of physical parameters.

**Gate A:** No flow equation → no tripartite synthesis.

### Phase B — Measure and Probability (T-Prob-1)
1. Compute the fixed-point measure.
2. Show absolute continuity and finite Fisher information.
3. Derive the Born rule as the projection of the fixed-point measure.
4. Decide T-Born vs Valentini-H primacy with concrete convergence rates.

**Gate B:** Singular measure → quantum potential diverges → synthesis dead.

### Phase C — IR Effective Theory (BM derivation)
1. Derive the IR effective action; show its kinetic term is the quantum potential.
2. Show the IR limit is the Schrödinger equation.
3. Derive the guidance equation from the substrate's update rule.
4. Show the IR dynamics is mixing (Axiom P2′).

**Gate C:** IR limit not Gaussian → linearity not derived → synthesis wounded.

### Phase D — Empirical Cross-Checks
1. Compute spectral dimension flow; check 2→4.
2. Compute noise spectrum; check for p-adic bath signature.
3. Compute relic non-equilibrium handle; check against CMB bounds.
4. Run Gate C (tensor product / entanglement from two-qubit construction).

**Gate D:** Spectral dimension flat at 4 → skeleton geometry wrong.

### Phase E — Confinement and Bell
1. Derive exact-setting distribution from substrate.
2. Clear Wood-Spekkens gate.
3. Reproduce Bell-exact correlations.

**Gate E:** Fine-tuning required → synthesis fails.

---

## 11. Death Conditions

**Cheap to trigger. Honest.**

| ID | Condition | Consequence |
|---|---|---|
| D1 | No discrete RG flow equation exists with the required properties | Tripartite collapse |
| D2 | Fixed-point measure is singular (not absolutely continuous) | Quantum potential diverges; IR effective theory undefined |
| D3 | IR fixed point is non-Gaussian | Linearity not derived; Schrödinger equation not IR limit |
| D4 | Mixing theorem fails for the substrate | Valentini-H inapplicable; Born rule not derived |
| D5 | Spectral dimension flat at 4 | Skeleton geometry wrong; no UV/IR crossover |
| D6 | Scheme-dependence of any derived constant | That constant is an artifact, not a prediction |
| D7 | Regulator-dependent T3 texture | Statistical predictions dead |
| D8 | Sustained non-equilibrium detected in laboratory | Structural selection falsified |
| D9 | Relic non-equilibrium detected inconsistent with T-Born primary | Probability-route decision wrong |
| D10 | Diffeomorphism invariance fails to recover above `ξ₀` | Axiom R violated; shadow inconsistent |
| D11 | Internal fixed point moves with truncation scheme | E8 attachment (if used) dead |
| D12 | Modular fingerprint absent from deviation spectrum | Lattice organization dead |

**The Occam guillotine:** If the substrate produces a clean 4D Gaussian fixed point with no arithmetic fine structure, the arithmetic layer was unnecessary and cheaper alternatives win.

---

## 12. Live Claims Ledger

| CLAIM_ID | Statement | Status | Blocking debt |
|---|---|---|---|
| TRI-ARCH-1 | One flow, three regimes | Architecture | — |
| TRI-AXIOM-P1′ | FRG trace as primary mechanism | Speculative-shape | Flow equation |
| TRI-AXIOM-R | Recovery axiom (SM+GR above ξ₀) | Speculative-shape | Gated |
| TRI-AXIOM-P2′ | RDS mixing (weakened from hyperbolicity) | Speculative-shape | Substrate verification |
| TRI-AXIOM-P3′ | Coarse-grained independence | Speculative-shape | T3 texture computation |
| IST-IU-CRIT-1 | I_U = critical surface of RG flow | Derived (given P1′) | Flow equation |
| IST-L-CFUNC-1 | L from c-function along RG trajectory | Speculative-shape | needMap |
| IST-DEFCUT-HEIGHT-1 | Preperiodic locus replaces Niven cut | Speculative-shape | needMap, needFaithfulnessReview |
| IST-ALPHABET-MARKOV-1 | p-adic from Markov partition | Speculative-shape | needToyCheck |
| IST-RATE-RGMATCH-1 | Collapse rate = RG coarse-graining rate | Speculative-shape | needMap |
| IST-MEASURE-FRACTAL-1 | Fractality relocates from set to measure | Speculative-shape | needToyCheck |
| AS-ONE-FP-2PROJ-1 | One fixed point, two projections | Speculative-shape | Cross-validation |
| AS-SPECTRAL-DIM-1 | Spectral dimension flows 2→4 | Speculative-shape | needToyCheck |
| AS-SCHEME-INDEP-1 | All constants scheme-independent | Speculative-shape | needObstruction |
| BM-WAVE-IR-1 | Wave function is IR order parameter | Derived (given flow) | needMap |
| BM-NONLOCAL-IR-1 | Nonlocality is IR artifact | Derived (given flow) | needMap |
| BM-GUIDANCE-IR-1 | Guidance equation is IR limit | Derived (given flow) | needMap |
| BM-FOLIATION-MARKOV-1 | Foliation from Markov partition | Speculative-shape | needMap |
| BM-QE-ATTRACTOR-1 | Quantum equilibrium as attractor | Derived (given Valentini) | needToyCheck |
| VAL-H-PATH-1 | Valentini supplies path, attractor supplies destination | Architecture | — |
| VAL-TBORN-PRIMACY-1 | T-Born primary, Valentini-H relic-scoped | Decision | needMap |
| VAL-RELIC-HANDLE-1 | Relic non-equilibrium handle | Speculative-shape | needToyCheck |
| E8-PACKING-INFO-1 | E8 as optimal packing for bath's internal code | Speculative-shape | needMap |
| E8-CFT-CHIRALITY-1 | "CFT representation theory resolves chirality" | **REJECTED** | Holomorphic VOA has one module |
| E8-CFT-CHIRALITY-2 | Narrower heterotic-style chirality | Speculative-shape | Doubly conditional |
| OSTROWSKI-BALANCE-1 | Product formula as ledger check | Check | needMap |

---

## 13. What Was Subtracted

**Design discipline record. Every removal loses a feature but preserves a load-bearing claim.**

| Removed | Rationale |
|---|---|
| Palmer's posited `I_U` | Replaced by derived critical surface |
| Palmer's hand-set `L` | Replaced by c-function derivation |
| Palmer's Niven-theorem rational cut | Replaced by preperiodic locus |
| Palmer's posited p-adic metric | Replaced by Markov partition alphabet |
| Palmer's fixed 1-bit/Planck-time rate | Replaced by RG coarse-graining rate |
| Palmer's positive-codimension set | Replaced by measure fine structure |
| AS's fundamental metric | Replaced by emergent order parameter |
| AS's continuum Wetterich equation | Replaced by discrete RG flow |
| AS's interpretive minimalism | Refused — fixed point needs state-space shadow |
| BM's postulated wave function | Replaced by IR order parameter |
| BM's postulated guidance equation | Replaced by IR limit of update rule |
| BM's fundamental nonlocality | Replaced by IR artifact |
| BM's hand-inserted foliation | Replaced by Markov partition |
| BM's continuous configuration space | Replaced by IR approximation |
| BM's mass as parameter | Replaced by RG-relevant direction |
| Hecke-algebra machinery | Only needed for Gate C; toy model hasn't run; not load-bearing |
| Langlands machinery | Decorative; no gate requires it |
| Arakelov-height conjecture | Speculative; no gate requires it |
| Multiple Born mechanisms run in parallel | Decision made: T-Born primary, Valentini-H relic |
| E8-as-unification | Rejected; only E8-as-packing survives |
| 248-as-particle-count framing | Retired; Lie algebra dimension, not census |
| Adelic product formula as axiom | Demoted to check |
| AS/IST as separate fundamental theories | Collapsed to regimes of one flow |

**Result:** The document that remains is the minimum needed to make the tripartite thesis testable.

---

## 14. Honest Open Problems

These are the things the synthesis **does not yet resolve**. They are tracked, not hidden.

### 14.1 The flow equation does not exist

Axiom P1′ requires a discrete RG flow equation. It has not been constructed. Everything downstream is conditional on it.

### 14.2 The measure problem is unresolved

Axiom P2′ provides the theorem class, but the specific substrate's membership in it is unconstructed. If the fixed-point measure is singular, the quantum potential diverges and the IR effective theory is undefined.

### 14.3 Three no-go theorems unaddressed

Haag, Coleman-Mandula, and Weinberg-Witten have not been checked against the tripartite framework. They must be addressed before the synthesis can claim QFT-level consistency.

### 14.4 T-Born vs Valentini-H primacy is a decision, not a derivation

The decision that T-Born is primary is a design choice. It is not forced by the framework. If relic non-equilibrium is detected, the decision reverses and the synthesis must explain why T-Born fails in decoupled regimes.

### 14.5 The borrowed-rigor risks are real

- **Mixing rate vs RG rate** (Claude): "these are the same kind of thing" is not "these are numerically the same." The claim is tagged and gated.
- **Foliation from Markov partition** (Claude): attractive, but requires the partition to exist and have the right causal properties. Tagged and gated.
- **Wave function as IR order parameter** (all three): stated as derived, but the derivation depends on the flow equation that does not exist.

### 14.6 The initial-condition problem

Valentini's H-theorem requires a non-equilibrium initial state to relax from. The synthesis does not yet specify why the substrate should have started away from its own equilibrium measure. This is a new, previously unopened debt item.

### 14.7 The gauge sector and generations

The tripartite synthesis has no account of the gauge sector or three generations. This is the largest gap relative to its ambitions. The E8 attachment is doubly conditional and does not resolve it.

### 14.8 The scale Λ

The cosmological constant is unexplained. The synthesis inherits this problem from every other candidate.

### 14.9 Gate C (tensor products)

The emergent tensor-product structure for many-body systems has not been constructed. Without it, Bell calculations cannot be performed inside the framework.

---

## 15. The One-Sentence Summary

**The BM–IST–AS tripartite synthesis is one arithmetic substrate with a discrete RG flow that produces Bohmian mechanics in the IR, with the Born rule emerging as the fixed-point measure via Valentini-type relaxation — four axioms, three regimes, one flow, and death conditions cheap enough to trigger within weeks.**

---

## Appendix A — Optional E8 Attachment

**Not load-bearing. Retained only because two claims have not been definitively closed.**

### A.1 Live claim: E8 as optimal packing for the bath's internal code

**Status:** `[SPECULATIVE-SHAPE]`, analogy only.

**Missing selection principle:** Why would nature use the *optimal* code? Plausibly the entropy-saturation machinery already on the ledger from the probability-gate work. If a system saturating its information bound naturally favors an optimal code, the connection is available; it does not require a new import.

**Kill condition:** If the entropy-saturation measure does not favor optimally-packed codes over other codes of comparable rate, this attachment fails.

### A.2 Doubly conditional: E8 lattice CFT as UV candidate

**Status:** `[SPECULATIVE-SHAPE]`, strictly downstream of the spectral-dimension gate.

**Conditions:** (1) spectral dimension flows to 2; (2) the resulting 2D theory is lattice/rational type rather than Liouville-type.

**What it is not:** a resolution of Distler-Garibaldi. Chirality would come from heterotic-style topological structure, not from the E8 CFT's own representation theory.

### A.3 Rejected claim: E8-CFT chirality via representation theory

**Status:** **REJECTED**, preserved per ledger policy.

**Reason:** The E8 lattice generates a holomorphic vertex operator algebra with exactly one irreducible module: itself. There is no representation-theoretic richness for a three-generation chiral spectrum to come from.

**Standing caution:** E8's dimension (248) is the Lie algebra's generator count, not a particle census. It should never be cited as explaining "how many particles exist."

---

## Appendix B — The Design Principle

**"You have reached design optimality when there are no more features to subtract."**

This document's structure reflects that principle:

- **Four axioms** — any fifth must be derived from the four.
- **Three regimes** — no fourth regime is needed.
- **One flow** — the flow is the primary object.
- **Named computations** — every open item is a computable target.
- **Cheap death conditions** — the synthesis can be killed in weeks, not decades.
- **No decoration** — Hecke, Langlands, Arakelov, and E8-as-unification are removed unless a gate requires them.

**The remaining document is the minimum needed to make the tripartite thesis testable.**

---

*End of synthesis v6.*
```

---

Save this as `BM_IST_AS_Tripartite_v6.md`. The document is designed to be read once, end to end, as a working specification. Every claim carries a status tag. Every open problem is named. The design principle — subtract until no further feature can be removed — has been applied throughout, with the subtraction record in §13 as its audit trail.