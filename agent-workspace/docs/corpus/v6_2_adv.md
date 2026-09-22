---
title: "BM + IST + Asymptotic Safety — v6.2"
subtitle: "Adversarially Revised Tripartite Research Program"
author: "Working Research Charter"
date: "2026-09-20"
---

# Status

**Document class:** research program, not a constructed physical theory.

**Current claim:** the program is a disciplined candidate framework for investigating whether a discrete arithmetic substrate can yield the observed relativistic quantum regime through a controlled scale-dependent effective description.

**It is not yet justified to claim:** a completed Theory of Everything, a proven invariant set, a proven Born derivation, a constructed dFRG, a derived Schrödinger equation, or a recovered Standard Model + GR.

**Design rule:** *optimal design is reached when no remaining feature can be removed without losing a load-bearing constraint, prediction, or falsifiability condition.*

**Constitution:** Elephant Bridge Protocol (EBP) v2.1 — **Ideas Enter Free. Promotion Costs Debt.**

---

# 0. Executive Thesis

The program starts from an uncomfortable premise:

> The three source theories may all contain genuine physical insight while also carrying assumptions that cannot safely be imported into a common foundational theory.

Bohmian Mechanics (BM), Invariant Set Theory (IST), and Asymptotic Safety (AS) therefore undergo the same treatment. None is sacred.

The program retains only a small architectural core:

1. a candidate microscopic state space $X$ with a discrete/arithmetic structure;
2. a candidate deterministic microscopic update $F:X\to X$;
3. a candidate controlled coarse-graining structure;
4. an observable projection into the effective physics we actually measure.

Everything else is a **derived target, theorem obligation, or empirical question**.

The intended architecture is therefore not:

$$
\text{IST} + \text{AS} + \text{BM}
$$

as three stacked theories.

It is:

$$
X
\;\xrightarrow{\mathrm{microscopic\ dynamics}}
\text{coarse-grained structures}
\;\xrightarrow{\mathrm{candidate\ RG/effective\ flow}}
\text{IR quantum description}
\;\xrightarrow{\mathrm{recovery}}
\text{known physics}.
$$

The phrase **"one flow" is deliberately no longer used as an established architectural fact**. There is currently no constructed flow. The correct status is **one candidate flow architecture whose existence is itself a central theorem target**.

The central question is correspondingly sharp:

$$
\boxed{
\text{Can one explicit microscopic dynamical system generate the structures of quantum theory without assuming them?}
}
$$

The research program lives or dies on that question.

---

# 1. Why the Three Foundational Theories Must Be Reconsidered

## 1.1 The general principle

A theory that is successful in a regime need not be fundamental in that regime's deepest possible sense.

A low-energy theory may correctly describe experiments while encoding assumptions that become illegitimate when its domain of validity is extended beyond the scale at which those assumptions were justified.

The same logic applies to BM, IST, and AS.

The program therefore asks, for every foundational statement:

> Is this an observed regularity, an established mathematical theorem, a model-dependent assumption, a convenient representation, or an unconstructed bridge?

Those categories must never be silently collapsed.

## 1.2 Bohmian Mechanics: what must not simply be assumed

Standard BM gives a precise and valuable ontology: configurations are definite, and their evolution is tied to a wavefunction.

But the tripartite program cannot simply promote all of the following to microscopic axioms:

- literal continuum configuration space;
- exact unitary Schrödinger evolution at every scale;
- the Born distribution as a fundamental equilibrium postulate;
- fixed particle number as the fundamental matter ontology;
- an inserted preferred foliation;
- Markovian, memoryless effective dynamics;
- the textbook first-order guidance law as the microscopic law;
- the quantum potential as primitive.

A discrete substrate makes the first two especially delicate. If the substrate is genuinely countable/discrete, continuous unitary Schrödinger evolution cannot simply be its primitive dynamics. The continuum must be an **IR target**.

The relativistic problem also makes fixed-$N$ BM a poor universal starting point. Bell-type quantum field theories therefore become the downstream BM benchmark rather than fixed-$N$ nonrelativistic BM.

## 1.3 Invariant Set Theory: what must not simply be assumed

IST contributes the most distinctive ontology to the program, but that is precisely why it receives the strongest scrutiny.

The following cannot be treated as conclusions before the substrate is constructed:

- the existence of a physically realized invariant set $I_U$;
- a positive-codimension fractal geometry;
- a particular $p$-adic metric;
- a particular prime;
- a hand-set information capacity $L$;
- a hand-set update/collapse rate;
- a rational/irrational definability cut;
- measurement-independence violation as a separate postulate;
- a particular chaos/attractor structure.

The research program retains one foundational commitment more openly than the earlier versions did:

> **The incumbent substrate hypothesis is discrete and arithmetic.**

This is the one deliberate retained commitment. It is not declared true. It is declared **the hypothesis under test**.

Even "arithmetic" is typed rather than overcommitted: the permitted search class may include finitely generated algebraic structures, number-field actions, adelic constructions, and ultrametric completions. The particular arithmetic object must be determined by the smallest successful construction, not chosen to fit a desired result.

## 1.4 Asymptotic Safety: what must not simply be assumed

AS contributes mature RG discipline, but the tripartite program must not turn AS vocabulary into metaphysics.

The following remain open:

- existence of the relevant fixed point for the actual substrate;
- physical meaning of its fixed-point coordinates;
- relation between theory-space critical surfaces and state-space invariant sets;
- regulator/scheme independence under the required truncations;
- compatibility with the emergent matter sector;
- Euclidean-to-Lorentzian continuation;
- whether the continuum description remains fundamental or is itself emergent.

The Wetterich/effective-average-action framework is therefore adopted as a **tool family**, not a sacred equation.

If an exact discrete coarse-graining construction turns out to be better represented by a transfer operator, real-space RG, operator-algebraic flow, or another equivalent structure, the program must use the structure that actually works.

---

# 2. EBP v2.1: The Constitutional Layer

## 2.1 The governing rule

EBP v2.1 has one central asymmetry:

> **Imagination is cheap. Promotion is expensive.**

An idea may enter without proof, map, invariant, toy model, null model, or formalization.

But promotion requires its obligations to be explicit and sufficiently retired.

The minimum rule is:

```text
Ideas enter free.
Promotion costs debt.
Debt does not kill.
Debt is forever payable.
New evidence creates new debt.
Promotion never means truth.
```

This is not project management. It is an epistemic safeguard.

## 2.2 Why EBP is necessary for a TOE program

A TOE program is unusually vulnerable to motivated reasoning because:

- many components are mathematically sophisticated;
- the empirical distance from the deepest claims is large;
- analogies between mature theories are easy to overread;
- a successful sub-result can create unjustified confidence in unrelated bridges;
- the same researcher can construct, formalize, interpret, and "verify" a claim.

EBP prevents an attractive narrative from silently acquiring the status of a physical axiom.

## 2.3 Debt classes

Every serious bridge eventually owes some combination of:

- `needMap` — what is the domain, codomain, and translation rule?
- `needInvariant` — what survives the bridge?
- `needToyCheck` — what finite computation can fail early?
- `needNullModel` — could a cheaper theory explain the same phenomenon?
- `needObstruction` — what theorem, counterexample, or incompatibility threatens the claim?
- `needFaithfulnessReview` — does the formal statement actually encode the intended physical claim?

Two additional debt classes are promoted in v6.2:

- `needInitialCondition` — what initial/boundary preparation is required for a relaxation argument?
- `needRegularity` — are the smoothness, phase, nodal, and generator-domain assumptions sufficient for the claimed quantum limit?

## 2.4 Promotion governance

EBP must not become bureaucracy, but self-grading alone is too weak.

The minimum governance is:

**Owner:** constructs the claim and evidence.

**Independent reviewer:** checks the debt retirement and may file a one-sentence objection.

**External red-team:** runs at milestone gates and before any claim is presented as a promoted physical result.

The reviewer does not possess a veto by prestige. The reviewer files debt with a stated reason.

A promoted claim can be unpromoted.

This is the practical meaning of:

$$
\text{new evidence}\Rightarrow\text{new debt}.
$$

---

# 3. The TOE Evaluation Constitution

The program evaluates itself with the same framework it applies to rival TOE candidates.

## 3.1 Nature's structural invariants

The working invariant set is:

- I1: sensitive dependence where present;
- I2: invariant/equilibrium measures;
- I3: attractor or invariant-set structure where appropriate;
- I4: structural stability;
- I5: universality;
- I6: fractal/spectral dimensions where appropriate;
- I7: bifurcation structure where applicable;
- I8: emergence;
- I9: renormalization;
- I10: criticality;
- I11: computational structure;
- I12: information accounting/conservation;
- I13: entropy and information bounds;
- I14: thermodynamic irreversibility.

These are not all axioms of nature. Some are conditional dynamical phenomena. Their role is to expose whether the proposed microscopic structure has the mathematical behaviors required of a candidate world.

## 3.2 Blind Men and the Elephant

The metaphor is methodological, not decorative.

Each established theory may have touched a real structural feature while overextending its language.

The operating principles are:

- **partial correctness** — retain real mathematical successes;
- **surface incompatibility may hide structural compatibility**;
- **no totalizing claim**;
- **projection to observables is mandatory**;
- **structure outranks terminology**;
- **different descriptions must be translated before comparison**;
- **convergence of independent descriptions is stronger than rhetorical similarity**.

The "elephant" is not another theory. It is the unknown structure that makes the successful descriptions cohere.

## 3.3 AI-evaluation analogue

The TOE receives the same classes of stress test used for difficult learned systems:

| Eval | Question for the research program |
|---|---|
| E1 Benchmark reproduction | Does it recover established physics in declared limits? |
| E2 Held-out prediction | Does it predict an unused shape or scaling law? |
| E3 Adversarial/no-go | Does it survive known no-go theorems or explain exactly which assumptions fail? |
| E4 Regression | Does new evidence create repairable debt rather than ad hoc rescue? |
| E5 Calibration | Are theorem, target, conjecture, and analogy clearly separated? |
| E6 Distribution shift | Does it survive early-universe, black-hole, many-body, and relativistic regimes? |
| E7 Ablation | Does removing a claimed component actually damage a load-bearing prediction? |
| E8 Contamination | Were predictions frozen before the data were examined? |
| E9 Interpretability | Can an independent researcher reproduce the bridge? |
| E10 Robustness | Do results survive perturbation, regulator changes, and alternate discretizations? |
| E11 Falsifiability | Is there a feasible observation that could refute the claim? |
| E12 Kill conditions | Is there an explicit condition under which the claim is abandoned? |

## 3.4 Coherence criteria

The program also evaluates itself on:

- C1 internal consistency;
- C2 parsimony;
- C3 unification;
- C4 naturalness;
- C5 independent cross-checks;
- C6 honest accounting;
- C7 survival under attack.

No numerical score is assigned before the underlying gates have actually run.

---

# 4. Self-Audit: The Program Turns Its Own Weapon on Itself

Current status is intentionally poor.

| Area | Current status | What is missing |
|---|---|---|
| I2 invariant measure | OPEN | Construct and characterize $\mu_*$ |
| I5 universality | TARGET | Cross-substrate test |
| I6 spectral dimension | UNVERIFIED | Compute $d_s(k)$ |
| I8 emergence | OPEN | Actual continuum/projection construction |
| I9 RG | OPEN | Construct candidate coarse-graining flow |
| I12 information | PARTIAL | Exact accounting through tracing |
| I14 thermodynamics | TARGET | Derive arrow from effective coarse-graining |
| E1 benchmark recovery | UNVERIFIED | Recovery theorem stack |
| E2 held-out prediction | UNRUN | Freeze a derived shape first |
| E3 no-go audit | PARTIAL | Full Haag/Coleman-Mandula/Weinberg-Witten analysis |
| E7 ablation | DESIGNED, NOT RUN | Remove each load-bearing component |
| E8 contamination | PROTOCOL READY | Pre-register actual calculations |
| E11 falsifiability | TARGETS EXIST | Must derive non-arbitrary prediction form |
| E12 kill conditions | DESIGNED | Trigger them honestly |
| C1 consistency | IMPROVED, UNPROVEN | Bridge construction |
| C2 parsimony | DESIGNED | Relevant-direction count not computed |
| C5 cross-checks | DESIGNED | Independent computations not yet performed |
| C6 accounting | OPERATIONAL | Process property, not physics evidence |
| C7 attack survival | PARTIAL | External red-team cadence not yet run |

The program therefore makes a stronger and narrower statement:

> **The program is presently eval-designed, not eval-passed.**

That sentence is itself a required part of the program.

---

# 5. Typed Mathematical Spaces

The most important discipline of v6.2 is to keep mathematical objects in the spaces where they actually live.

Let

$$
X=\text{microscopic state space},
$$

$$
\mathcal M(X)=\text{probability measures on }X,
$$

$$
\mathcal T=\text{effective-action / theory space},
$$

$$
\mathcal O=\text{observable space}.
$$

An RG transformation acts on theory descriptions:

$$
\mathfrak R_s:\mathcal T\to\mathcal T.
$$

A microscopic update acts on states:

$$
F:X\to X.
$$

An observable construction acts toward measurable quantities:

$$
\Pi:X\longrightarrow\mathcal O
$$

or, for measures,

$$
\Pi_\#:\mathcal M(X)\to\mathcal M(\mathcal O).
$$

These are not interchangeable.

## 5.1 Fixed point versus invariant set

Let $\Gamma_*\in\mathcal T$ be a candidate fixed point,

$$
\mathfrak R_s(\Gamma_*)=\Gamma_*.
$$

Let $W^s(\Gamma_*)\subset\mathcal T$ be its stable/critical surface.

Do **not** write

$$
I_U=W^s(\Gamma_*).
$$

That is ill-typed.

Instead construct a correspondence

$$
\Phi:X\to\mathcal T
$$

and test the theorem target

$$
\boxed{
I_{\rm crit}\stackrel{?}{=}\Phi^{-1}\!\left(W^s(\Gamma_*)\right).
}
$$

Only after such a map exists can the old symbol $I_U$ be justified as a physical invariant set.

This is the exact point where the former theory-space/state-space conflation is removed.

---

# 6. Minimal Primitive Charter

The charter is deliberately small.

## P0 — Candidate microscopic substrate

There exists a deterministic microscopic dynamical system

$$
(X,F)
$$

with a specified class of discrete/arithmetic structure.

**Status:** program charter, not established physics.

## P1 — Candidate controlled coarse-graining

There exists a family of maps, kernels, partitions, or equivalent structures that define a physically meaningful notion of scale:

$$
\mathcal C_s:X\to X_s
$$

and an induced effective flow.

**Status:** theorem target.

## P2 — Candidate observable projection

There exists a physically faithful map from the microscopic/effective description to experimentally meaningful observables:

$$
\Pi_s.
$$

**Status:** theorem target.

Everything else must earn its way in.

This is the core subtraction.

---

# 7. First-Principles Re-foundation of Bohmian Mechanics

## 7.1 What survives

BM supplies the IR target because it gives an unusually explicit answer to the ontology of quantum events:

- definite configurations/outcomes;
- a dynamical state capable of producing interference;
- a guidance structure;
- a coherent account of measurement outcomes.

These are **targets of derivation**, not microscopic axioms.

## 7.2 What is demoted

The following are regime statements:

- continuous configuration space;
- Schrödinger evolution;
- the Born distribution;
- the textbook guidance law;
- the usual quantum potential;
- fixed particle number;
- Markovianity.

## 7.3 Generalized Langevin first; Kramers second

The reduced dynamics should initially be written in the broader memory-bearing form:

$$
\dot q = \frac{p}{m},
$$

$$
\dot p
=
-\nabla V(q)
-\int_0^t K(t-s)p(s)\,ds
+\eta(t).
$$

Only if a controlled limit gives

$$
K(t-s)\rightarrow 2\gamma\,\delta(t-s)
$$

should the Kramers equation be adopted.

This avoids assuming Markovianity merely because it is convenient.

## 7.4 Asymptotic Bohmianity

The target is

$$
v_k(x)=\frac{\nabla S}{m}+\delta v_k,
$$

with

$$
\delta v_k\to0
\qquad
(k\to0)
$$

in the appropriate IR scaling regime.

This is a **derived-regime hypothesis**.

## 7.5 Nonlocality: no premature rescue

The attractive claim

$$
\text{nonlocality}_{IR}
=
\text{artifact of UV locality}
$$

is **not** part of the core.

It may be tested later, but the program refuses to assume that the substrate is local merely to save relativistic aesthetics.

The IR target must still reproduce the nonlocal dependence required for entanglement while satisfying no-signaling.

## 7.6 Bell-type QFT as downstream benchmark

If matter number is itself emergent, fixed-$N$ BM cannot be the ultimate formulation.

Bell-type quantum field theory is therefore the benchmark for the relativistic matter sector:

- particle number may change;
- configurations may jump;
- the effective state remains capable of generating the appropriate quantum statistics.

Connection to the substrate is a separate theorem obligation.

## 7.7 Foliation

A preferred foliation may be:

- fundamental,
- emergent,
- gauge-like,
- or unnecessary.

The program chooses none in advance.

Candidate source mechanisms include the coarse-graining ordering, a dynamically selected slicing, or a derived structure of the microscopic map.

The only hard requirement is empirical and structural:

> The construction must reproduce relativistic covariance/no-signaling in the declared regime.

---

# 8. First-Principles Re-foundation of Invariant Set Theory

## 8.1 The retained IST commitment

The program retains:

> **A discrete arithmetic substrate is the incumbent hypothesis.**

Everything else is open.

## 8.2 Information capacity $L$

Do not hand-set $L$.

Candidate mechanism:

$$
L \sim C(\text{RG trajectory}),
$$

where $C$ is an appropriately defined monotone information quantity, analogous in role to a $c$-function.

The claim survives only if the construction fixes $L$ without post-hoc tuning.

## 8.3 Definability: canonical height

Suppose the actual microscopic map is an algebraic self-map of appropriate degree.

Then a canonical-height route becomes testable:

$$
\hat h_F(x)\ge 0,
$$

with the preperiodic locus characterized, under the appropriate hypotheses, by

$$
\hat h_F(x)=0.
$$

This replaces a coordinate-dependent rational/irrational cut with a dynamical criterion.

But this route is conditional on the actual $F$ satisfying the required algebraic hypotheses.

## 8.4 Markov-partition route

A Markov partition may generate symbolic dynamics and hence an ultrametric coding structure.

But a Markov partition does **not** follow merely from having an arbitrary discrete dynamical system.

Its existence must be tested.

## 8.5 Joint-satisfiability gate

The canonical-height and Markov-partition routes cannot be allowed to coexist merely because both are attractive.

The actual map $F$ must be tested for:

$$
\boxed{
\text{algebraic self-map}
\quad+\quad
\text{required degree}
\quad+\quad
\text{Markov partition / symbolic factor}
}
$$

If these conditions cannot be jointly realized, the program must choose among:

1. height route;
2. symbolic route;
3. a third mechanism;
4. abandonment of the proposed arithmetic definability mechanism.

This gate is mandatory.

## 8.6 Prime universality

The program must not hard-code $p=2$, $p=3$, or any other prime.

Run the same construction across

$$
p=2,3,5,7,\ldots
$$

and classify the outcome:

- prime-independent universality;
- a uniquely selected prime;
- prime-dependent physics;
- necessity of a multi-place/adelic product.

A prime that is chosen because it makes the numbers work is contamination.

## 8.7 Fractality relocates to observables

Nondegenerate stochastic smoothing tends to destroy literal positive-codimension support at coarse scales.

Therefore the safer target is:

- continuum/absolutely continuous behavior at observable scales;
- arithmetic fine structure at sub-resolution scales.

A measured texture may have multifractal scaling, but that must be computed.

It is not a primitive.

---

# 9. The Amplitude Problem Is a First-Class Gate

This is the major new elevation in v6.2.

Earlier versions correctly said "derive the wavefunction," but treated that obligation too generically.

The sharper conclusion is:

> **The difficult problem is not merely deriving Fisher information or writing down the Madelung equations. It is constructing a complex amplitude sector that is genuinely generated by the microscopic dynamics and closed under evolution.**

## 9.1 What is already known

There are established variational/uncertainty-based derivations in which Fisher-information or exact-uncertainty assumptions contribute to obtaining Schrödinger dynamics.

Examples include Reginatto's Fisher-information derivation and the Hall-Reginatto exact-uncertainty program.

These results matter because they prevent the research program from claiming novelty where the mathematics is already mature.

The correct question is therefore not:

> "Can Fisher information lead to Schrödinger?"

It is:

> "Can the microscopic arithmetic substrate produce the assumptions needed for a complex amplitude theory without inserting them?"

## 9.2 Native p-adic character route

The strongest current candidate for the phase interface is the additive character

$$
\chi_p(x)
=
e^{2\pi i\{x\}_p}.
$$

This is not promoted as the answer.

It is a **high-priority workbench hypothesis** because it gives a mathematically native route from p-adic structure to unit-circle phase.

The required tests are:

1. Is $\chi_p$ defined on the required state classes?
2. Does it respect composition?
3. Does it close under the microscopic update?
4. Can a nontrivial amplitude be built from characters?
5. Does the resulting complex sector have a physically meaningful norm?
6. Does the phase survive coarse-graining?
7. Can a Schrödinger-type generator emerge without inserting linearity by hand?

## 9.3 Chaos-tameness / factor ladder

The program must not assume that a chaotic substrate automatically carries the desired complex Hilbert structure.

A high-entropy or mixing system can have spectral properties that make naive "search inside the dynamics for a linear quantum sector" misleading.

Therefore the amplitude search is organized as a ladder:

**T1:** invariant measure and dynamical spectrum exist.

**T2:** the required spectral factor exists.

**T3:** the factor supports a nontrivial complex/unit-circle structure.

**T4:** that factor closes into the effective amplitude dynamics required for quantum theory.

If no such factor exists, the program must not "manufacture" one by redefining observables until it appears.

This is a potential kill condition.

---

# 10. The Projection Problem Becomes Three Problems Plus One Hard Core

## P0 — State emergence

Construct

$$
X_L
\longrightarrow
\mathcal C_{\rm eff}
$$

and show that the effective configuration space has the topology, dimensionality, and regularity required by the intended IR theory.

## P1 — Measure emergence

Construct

$$
\mu_*
\longrightarrow
\Pi_{\epsilon\#}\mu_*.
$$

Show existence, uniqueness, absolute continuity in the observable regime, and sufficient regularity.

## P2 — Dynamics emergence

Show

$$
\mathcal L_s
\longrightarrow
\mathcal L_{\rm QM/BM}
$$

with the correct generator, not merely the correct stationary density.

## P3 — Amplitude-factor emergence

Construct the complex amplitude sector:

$$
X
\;\to\;
\mathcal H_{\rm eff}
$$

or an equivalent mathematical object, without assuming that Hilbert space already exists.

This is now the central bottleneck.

---

# 11. Born Probability: Exact Where It Means Exact

The earlier correction-field conflict is resolved by making the projection scale explicit.

## 11.1 B0 — Equilibrium measure

$$
\exists!\,\mu_*.
$$

## 11.2 B1 — Regularity

At an observable coarse-graining scale $\epsilon$,

$$
\Pi_{\epsilon\#}\mu_*
=
\rho_\epsilon(x)\,dx
$$

with sufficient regularity, ideally

$$
I_F[\rho_\epsilon]
=
\int \rho_\epsilon |\nabla\ln\rho_\epsilon|^2\,dx
<\infty.
$$

## 11.3 B2 — Born identification at the IR observable algebra

The exact target is

$$
\boxed{
\Pi_{\rm IR\#}\mu_*
=
|\psi|^2\,dx
}
$$

where $\Pi_{\rm IR}$ is explicitly defined as the IR observable/projection map.

At finite resolution, the program allows a correction term:

$$
\Pi_{\epsilon\#}\mu_*
=
|\psi|^2
\left(1+\delta_{\xi,\epsilon}\right)\,dx,
$$

with the requirement that

$$
\delta_{\xi,\epsilon}\to0
$$

in the declared Born regime and that any residual structure remain below the observational resolution or satisfy a separately declared bound.

Thus "Born exact" and "arithmetic texture" are not contradictory statements because they are statements about different projections/scales.

## 11.4 No-signaling constraint

Any residual correction must satisfy the safety condition that marginal statistics remain unchanged:

$$
\int \delta\rho(x_A,x_B)\,dx_B = 0
$$

and similarly for the other subsystem, while nontrivial structure may survive in joint correlators.

If a correction permits controllable signaling, the theory fails regardless of its other virtues.

## 11.5 B3 — Relaxation

Valentini's coarse-grained H-theorem supplies a relaxation mechanism under its hypotheses:

$$
H[P|\rho]
=
\int P\ln\frac{P}{\rho}\,dx.
$$

The program requires an appropriate coarse-graining and a demonstrable dynamical mixing regime.

Valentini's original work explicitly treats $P=|\Psi|^2$ as an equilibrium reached statistically rather than assumed as a fundamental postulate.

The key methodological distinction remains:

$$
\boxed{
\text{Born destination}
\neq
\text{relaxation path}.
}
$$

## 11.6 Typicality

Standard BM typicality is not a fourth independent theory of probability.

It is treated as an interpretive/measure-theoretic benchmark:

> If the derived equilibrium measure is the Born measure, does the resulting typicality structure coincide with the usual BM statement?

If yes, the mechanisms have converged.

If no, the disagreement becomes a diagnostic rather than another parallel story.

---

# 12. Linearity, Guidance and the Quantum Potential

## 12.1 Polar decomposition

The target is

$$
\psi=\sqrt{\rho}\,e^{iS/\kappa}.
$$

The program must derive both amplitude and phase structure.

## 12.2 Schrödinger generator

The nonrelativistic target is

$$
i\kappa\,\partial_t\psi
=
-\frac{\kappa^2}{2m}\nabla^2\psi
+
V\psi.
$$

Linearity is its own gate.

The claim

$$
\text{"nonlinear microscopics often flow to linear IR behavior"}
$$

is not enough.

The actual generator must be constructed.

## 12.3 Guidance

The target is

$$
v=\frac{\nabla S}{m}.
$$

The program must show why the same phase field $S$ that defines the effective amplitude produces the trajectory field.

## 12.4 Quantum potential

The target is

$$
Q
=
-\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt{\rho}}{\sqrt{\rho}}.
$$

Possible derivations:

1. effective-action route;
2. generator/Madelung route;
3. a justified arithmetic/height route.

No route is privileged until it works.

---

# 13. The Internal-Environment Principle Becomes a Theorem Gate

The attractive statement

$$
\text{hidden arithmetic degrees of freedom}
\longrightarrow
\text{effective noise}
$$

requires an actual reduction.

The program therefore adds:

## Tensor-Factorization / Conditional-Expectation Gate

Show that the microscopic algebra or state structure admits a decomposition sufficient to define a reduced state.

The strongest familiar realization would be

$$
\mathcal H_{\rm total}
\simeq
\mathcal H_{\rm visible}\otimes\mathcal H_{\rm hidden},
$$

but the program does **not** assume Hilbert-space tensor factorization at the microscopic level.

A weaker algebraic conditional-expectation structure is acceptable if it produces the same operational trace.

This gate is mandatory because otherwise "the internal bath" is only a metaphor.

---

# 14. Asymptotic Safety as Discipline, Not Dogma

## 14.1 What is retained

AS contributes:

- effective-action flow;
- relevant/irrelevant directions;
- fixed-point analysis;
- universality;
- regulator/scheme discipline;
- spectral-dimension diagnostics;
- parameter reduction.

## 14.2 Candidate dFRG

The discrete flow may take the abstract form

$$
\mathfrak R_s:\mathcal T\to\mathcal T,
$$

or be represented through an effective action

$$
\Gamma_s.
$$

A Wetterich-like equation is one candidate realization:

$$
\partial_s\Gamma_s
=
\frac12
{\rm Tr}
\left[
(\Gamma_s^{(2)}+\mathcal R_s)^{-1}
\partial_s\mathcal R_s
\right],
$$

but the discrete substrate determines what "trace", "mode", "regulator", and "scale" actually mean.

## 14.3 Scheme independence

Intermediate couplings may depend on the regulator.

Universal observables and critical quantities must demonstrate stability.

Therefore the program requires:

$$
\text{multiple regulators}
+
\text{multiple truncations}
+
\text{convergence study}.
$$

A number that moves materially under legitimate changes is not promoted as physical.

## 14.4 Fixed point

First establish:

$$
\Gamma_*.
$$

Then:

$$
W^s(\Gamma_*).
$$

Then test whether a state-space counterpart exists through $\Phi$.

No object is identified merely because its name is "fixed point" or "invariant set".

## 14.5 Relevant directions

If the fixed point has $N_{\rm rel}$ relevant directions, the program should explain why those directions correspond to the independent physical parameters.

Parameter proliferation is treated as a structural warning.

## 14.6 Spectral dimension

The program tests, rather than assumes:

$$
d_s(k)
$$

with an intended comparison such as

$$
d_s^{\rm IR}\approx4,
\qquad
d_s^{\rm UV}\approx2.
$$

The relation

$$
d_s=\frac{2d_H}{d_w}
$$

is conditional and may be used only when the required diffusion/fractal scaling exists.

A different result is data, not a failure "to be repaired".

## 14.7 Local scale problem

The rejected local rule

$$
k=k(x)
$$

is not reinstated as an axiom.

Instead the program records the genuine question as debt:

> Can one global RG scale parameter describe strongly inhomogeneous regimes, or is a locally adaptive coarse-graining required?

A local scale is permitted only if derived from the flow and shown to preserve recovery and scheme-independence.

---

# 15. The Candidate Common Architecture

The honest architecture is now:

$$
(X,F)
$$

at the microscopic level;

$$
\mathcal C_s,\;\mathfrak R_s
$$

as candidate coarse-graining/RG structure;

$$
\Gamma_s,\;\mu_s
$$

as effective dynamical and measure objects;

$$
\Phi:X\to\mathcal T,
\qquad
\Pi_\epsilon:X\to\mathcal O
$$

as correspondence maps;

and finally

$$
(\rho,S,\psi)
$$

as effective IR objects.

The central chain is therefore a **research target**:

$$
\boxed{
X
\overset{F}{\longrightarrow}
X
\overset{\mathcal C_s,\mathfrak R_s}{\longrightarrow}
(\Gamma_s,\mu_s)
\overset{\Pi}{\longrightarrow}
(\rho,S,\psi)
\longrightarrow
{\rm QM/BM}
\longrightarrow
{\rm QFT+GR}
}
$$

No arrow is assumed to exist merely because the diagram is aesthetically coherent.

---

# 16. Composition, Entanglement and Bell

## 16.1 Gate C

Construct two weakly coupled subsystems and show that their effective theory admits the appropriate composition law.

The desired target is

$$
\mathcal H_{AB}
\simeq
\mathcal H_A\otimes\mathcal H_B
$$

or a demonstrably equivalent algebraic structure.

## 16.2 Entanglement

Show that interactions can generate nonseparable states.

## 16.3 Bell

Reproduce the correct Bell correlations.

But do not stop there.

A complete Bell audit must track both:

- nonlocal dependence in the effective guidance;
- measurement-setting dependence or its exact replacement mechanism.

The program must not rename measurement dependence "contextuality" and declare the issue solved.

## 16.4 Double-counting audit

Explicit question:

> Does the constructed mechanism independently require both nonlocal guidance and measurement-setting dependence to explain the same correlations?

If yes, the theory is carrying redundant mechanisms and must be simplified.

---

# 17. QFT and Gravity Recovery

The program does not yet claim a TOE.

It claims a **candidate pre-geometric substrate and IR quantum theory with downstream TOE obligations**.

The recovery ladder is:

### QFT-R0
Construct a relativistic effective field theory.

### QFT-R1
Recover Lorentz symmetry to required precision.

### QFT-R2
Recover gauge structure and matter content.

### QFT-R3
Audit relevant no-go theorems:

- Haag;
- Coleman-Mandula;
- Weinberg-Witten;
- and any additional theorem whose assumptions overlap the proposed construction.

For each theorem:

1. list assumptions;
2. identify which hold;
3. identify which fail;
4. show why they fail;
5. show which structure re-emerges in the IR.

"UV cutoff" is not an explanation by itself.

### QFT-R4
Euclidean-to-Lorentzian continuation.

### QFT-R5
Compatibility of the AS fixed point with the emergent matter sector.

### GR-R0
Diffeomorphism recovery.

### GR-R1
Equivalence-principle recovery.

### GR-R2
Low-energy Einstein dynamics.

Failure at any of these means the program is not yet a TOE.

---

# 18. E8: Keep the Mathematics, Drop the Cargo

E8 is not a fourth pillar.

## Permanently rejected

- "248 particles in nature";
- naive E8 Lie-algebra embedding of all fields;
- the claim that an E8 lattice CFT's ordinary representation theory automatically resolves chirality.

## Surviving role

E8 is an optional module only if a physical selection principle selects an eight-dimensional internal code space.

Possible use:

- optimal information packing;
- error correction;
- modular structure;
- a held-out spectral fingerprint.

The sequence must be:

$$
\text{derive dimension}
\to
\text{derive code-selection criterion}
\to
\text{test E8}.
$$

Never:

$$
\text{notice E8}
\to
\text{fit physics}.
$$

If dimension 8 is not independently derived, E8 remains a mathematical comparison object.

---

# 19. Cross-Substrate and Prime Universality

A robust substrate theory should not depend on one hand-picked implementation.

Construct, where feasible, microscopically distinct substrates

$$
X_1,\;X_2,\;X_3,\ldots
$$

and test whether they flow into the same effective universality class.

Success would support structural, rather than implementation-specific, physics.

Failure is informative:

- one substrate succeeds uniquely;
- several succeed with distinct IR classes;
- all fail;
- or the alleged universal structure is an artifact.

The same logic applies to prime choices.

Prime universality is therefore a mandatory robustness study, not an optional embellishment.

---

# 20. Empirical Program: Shape Before Amplitude

No parameter fitting is allowed before the functional form is derived.

## T1 — Constraint compilation / pre-registration

Before a new formula exists:

- compile existing experimental bounds;
- freeze analysis pipelines;
- record which future formula families would be considered;
- prohibit post-hoc choice of the family after looking at data.

## T2 — Derived-shape tests

Examples:

- mesoscopic interferometry scaling;
- colored-noise spectrum;
- angular/Diophantine structure;
- spectral-dimension flow;
- relic non-equilibrium signatures.

The required object is a shape, exponent, or functional law.

## T3 — High-information fingerprints

Examples:

- modular statistics if E8 is independently selected;
- prime-independence/dependence;
- cross-substrate universality;
- scale-dependent departures from exact Born behavior.

A sign or threshold is not enough when a rival predicts a different functional shape.

---

# 21. Research Program: Information Gain Before Prestige

The corrected sequence fixes the earlier phase-ordering problem.

## Phase 0 — Constraint and contamination setup

Deliverables:

- compiled bounds;
- versioned preregistration templates;
- explicit status labels;
- external-review plan.

No claim is made that a physical kill test can be run before a derived prediction exists.

## Phase 1 — The first decisive microscopic computation

**Start here.**

Construct the smallest nontrivial explicit microscopic map

$$
F_L:X_L\to X_L.
$$

Use the actual declared substrate rules.

Immediately test:

1. determinism;
2. finite information bookkeeping;
3. algebraic/self-map conditions needed for canonical height;
4. existence or failure of a Markov partition;
5. symbolic entropy/spectral properties;
6. prime dependence;
7. compatibility with the native additive-character route.

This single computation can kill or redirect multiple branches.

### Phase-1 kill

If $F_L$ cannot be constructed faithfully, stop the arithmetic branch.

If $F_L$ exists but fails the conditions required by the canonical-height and symbolic routes, do not silently replace it with a different map.

## Phase 2 — Tameness and measure

Construct:

- invariant/equilibrium measure;
- ergodicity/mixing diagnostics;
- spectral decomposition;
- factor search;
- absolute-continuity/Fisher audit;
- initial-condition analysis.

### Required new artifact

A **T1–T4 amplitude-factor report** documenting whether the microscopic dynamics actually contains a viable complex sector.

## Phase 3 — Projection prototype

For the one-qubit problem, construct a concrete prototype rather than another desiderata list.

Candidate form:

$$
\pi_L:X_L
\to
\mathbb C^2
$$

built from explicitly declared microscopic observables/characters.

The first prototype must answer:

- what is the domain?
- what is the codomain?
- what structure is preserved?
- what is discarded?
- what topology/measure is induced?
- how is phase represented?
- how does $F_L$ act after projection?

This is the heart of the program.

## Phase 4 — Candidate discrete RG / dFRG

Construct:

$$
\mathcal C_s,\quad
\mathfrak R_s,\quad
\Gamma_s.
$$

Then test:

- fixed-point existence;
- critical surface;
- relevant directions;
- regulator independence;
- truncation convergence;
- spectral dimension;
- correspondence to state-space structures.

Only after these are built should the phrase "AS completes the UV" be used.

## Phase 5 — IR quantum derivation

Derive, in order:

1. configuration space;
2. amplitude;
3. phase;
4. generator;
5. continuity equation;
6. guidance law;
7. quantum potential;
8. Born projection;
9. asymptotic unitarity.

The order matters.

A correct Born density without the correct generator is not quantum mechanics.

A correct generator without an emergent amplitude is not a microscopic derivation.

## Phase 6 — Composition and Bell

Build two subsystems.

Test:

- factorization;
- entanglement generation;
- Bell statistics;
- no-signaling;
- setting-dependence audit;
- foliation behavior.

## Phase 7 — Relativistic and gravitational recovery

Only now:

- E→L;
- Lorentz symmetry;
- QFT;
- gauge structure;
- matter;
- no-go audit;
- GR;
- equivalence principle;
- cosmological sector.

## Phase 8 — Phenomenology

Only after the structural derivations exist:

- interferometry;
- Born-deviation constraints;
- noise spectra;
- cosmological relics;
- spectral dimension;
- optional E8 fingerprint.

---

# 22. The Program's Hardest Gates

The following are load-bearing:

### G0 — Discrete/continuous consistency

The continuum quantum theory must emerge; it may not be silently assumed.

### G1 — Explicit microscopic map

No vague "shift" or "symbolic dynamics" without a defined $F$.

### G2 — Joint arithmetic compatibility

Canonical height and symbolic/Markov machinery must actually apply to the same $F$, or one is removed.

### G3 — Amplitude-factor existence

The substrate must generate a genuine complex factor/representation suitable for a quantum amplitude.

### G4 — Measure

$$
\exists!\mu_*
$$

with the regularity needed for the IR theory.

### G5 — Projection

Construct $\Pi$ explicitly.

### G6 — Born

$$
\Pi_{\rm IR\#}\mu_*=|\psi|^2dx.
$$

### G7 — Dynamics

$$
\mathcal L_s\to\mathcal L_{\rm Schr/BM}.
$$

### G8 — Guidance and $Q$

$$
v=\frac{\nabla S}{m},
\qquad
Q=-\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt\rho}{\sqrt\rho}.
$$

### G9 — Composition

Entanglement and Bell.

### G10 — Recovery

QFT + Lorentz + gauge + GR.

No later gate can compensate for a failed earlier load-bearing gate.

---

# 23. Symmetric Adjudication of External Ideas

The rule applied to DeepSeek, Z, Gemini, Claude, or this program itself is identical:

## Retain in core if

- it removes an assumption;
- closes an existing gap;
- yields a cheaper decisive test;
- supplies established mathematics directly relevant to an existing gate;
- or converts a vague claim into a typed construction.

## Retain in workbench if

- mathematically interesting;
- physically plausible;
- but introduces a new dependency or cannot yet be connected to the core.

## Reject for the current program if

- it adds ontology before existing bridges are solved;
- it fits a parameter after seeing the target;
- it conflates objects from different mathematical spaces;
- it claims a no-go theorem has been evaded without identifying which hypothesis fails;
- or it duplicates an existing mechanism without new predictive leverage.

This is why the Gemini local-$k(x)$ and bundle-geodesic constructions are deferred, while p-adic characters are elevated to the amplitude workbench, and why "nonlocality is just an IR artifact of UV locality" remains unpromoted.

The same rule would be applied to any future contribution from this model.

---

# 24. Symmetric Occam Guillotine

The earlier guillotine targeted the arithmetic layer. It is now symmetric.

### D20 — Arithmetic redundancy

If arithmetic structure produces no distinctive successful observable signature, and a simpler non-arithmetic substrate explains the same data, remove the arithmetic layer.

### D21 — AS redundancy

If dFRG adds no regulator-stable quantitative content beyond what a simpler effective-theory construction already supplies, demote AS from foundational architecture to mathematical technique.

### D22 — BM redundancy

If the IR construction reproduces the observed quantum regime without a distinct Bohmian trajectory ontology, the Bohmian layer must be demoted accordingly.

### D23 — E8 redundancy

If E8's role cannot be selected independently of the desired result, remove it.

### D24 — EBP redundancy

If EBP becomes paperwork rather than an accelerator of falsification and debt retirement, simplify the protocol itself.

No component receives immunity from Occam.

---

# 25. Subtraction Record: What v6.2 Adds, Restores, or Reclassifies

## Restored at higher resolution

- Reginatto/Hall-Reginatto precedent;
- dedicated amplitude-factor gate;
- p-adic additive-character workbench;
- chaos-tameness/factor ladder;
- joint canonical-height/Markov compatibility gate;
- self-audit;
- Tensor-Factorization gate;
- symmetric AS redundancy kill condition;
- external red-team governance;
- explicit projection scale in the Born statement.

## Reinstated as explicit debt rather than architecture

- local/adaptive RG scale question;
- possible boundary-value/atemporal route;
- clustering/mass-gap transfer if it becomes necessary;
- relic non-equilibrium program;
- correction-field phenomenology.

## Explicitly rejected

- "one flow" as a fact before the flow is constructed;
- $I_U=W^s(\Gamma_*)$ as a literal identity;
- nonlocality-as-IR-artifact as an established conclusion;
- E8 as a particle census;
- E8 as an automatic chirality solution;
- post-hoc prime selection;
- local $k(x)$ as a primitive;
- any automatic identification of a theorem-class with applicability to the actual substrate.

---

# 26. Live Claims Ledger

| Claim | Status | Immediate dependency |
|---|---|---|
| Arithmetic substrate exists | INCUMBENT HYPOTHESIS | Explicit $F_L$ |
| Candidate microscopic map $F$ exists | OPEN | Phase 1 |
| Canonical-height mechanism | TARGET | $F$ algebraic/self-map |
| Markov-partition mechanism | TARGET | dynamical hypotheses |
| Prime universality | TARGET | family of $F_p$ or equivalent |
| Native complex phase via $\chi_p$ | WORKBENCH / HIGH PRIORITY | closure + factor test |
| Unique physical measure | TARGET | ergodicity/mixing |
| Born projection | TARGET | $\Pi_{\rm IR}$ |
| dFRG exists | TARGET | coarse-graining construction |
| Fixed point exists | TARGET | dFRG |
| $I_{\rm crit}=\Phi^{-1}(W^s)$ | TARGET | $\Phi$ + RG |
| Continuum configuration space emerges | TARGET | P0 |
| Schrödinger generator emerges | TARGET | amplitude + generator |
| Guidance emerges | TARGET | phase/generator |
| Quantum potential emerges | TARGET | generator/IR action |
| Composition emerges | TARGET | two-subsystem construction |
| Bell reproduced | TARGET | composition + guidance |
| Lorentz/QFT recovery | DOWNSTREAM TARGET | Phase 7 |
| GR recovery | DOWNSTREAM TARGET | Phase 7 |
| E8 physical role | OPTIONAL WORKBENCH | independent dimension/code selection |
| Arithmetic unique necessity | NOT CLAIMED | Occam gate |

---

# 27. Hard Failure Conditions

The following are explicit kill conditions.

1. No faithful microscopic map $F$.
2. The arithmetic class required by the program is mathematically incompatible with the actual $F$.
3. Canonical-height and Markov-partition requirements cannot coexist and no defensible replacement exists.
4. No viable complex factor or amplitude representation emerges.
5. The invariant/equilibrium measure is absent or unusably singular.
6. The IR projection cannot produce a continuous configuration manifold.
7. The Born measure fails at the defined IR observable scale.
8. The generator is not in the required linear quantum class.
9. Guidance cannot be recovered.
10. The quantum potential cannot be recovered.
11. Composition fails.
12. Bell correlations fail.
13. No-signaling fails.
14. Euclidean-to-Lorentzian continuation fails.
15. Required relativistic/gauge recovery fails.
16. Derived physical quantities are regulator/truncation artifacts.
17. Increasingly different substrates require hidden retuning to reach the same IR class.
18. Every purported arithmetic signature disappears while a simpler non-arithmetic model explains the observations.
19. AS contributes no distinctive quantitative or mathematical leverage.
20. E8 survives only by retrofitting.
21. The EBP process becomes bureaucratic overhead rather than debt-retiring work.

The rule is not "defend the theory."

The rule is:

$$
\boxed{\text{find out exactly where it breaks.}}
$$

---

# 28. The First 90 Days

## First 30 days

**Deliverable A:** explicit smallest nontrivial $F_L$.

**Deliverable B:** canonical-height applicability report.

**Deliverable C:** Markov-partition applicability report.

**Deliverable D:** first prime-universality scan.

**Deliverable E:** p-adic-character closure test.

The expected outcome is not success.

It is **branch elimination**.

## Days 31–60

If the microscopic map survives:

- construct invariant measure candidate;
- run mixing/spectral diagnostics;
- build the first factor search;
- attempt the one-qubit projection prototype;
- derive or falsify the first coarse-graining operator.

## Days 61–90

If the factor and measure survive:

- construct the first candidate effective flow;
- search for fixed points;
- begin regulator/truncation cross-validation;
- freeze the first held-out prediction shape;
- run the first external red-team review.

If any gate fails, the program branches or stops rather than adding more machinery.

---

# 29. Expected Shape of a Successful Theory

A successful outcome would not look like a giant list of new entities.

The desired shape is the opposite.

At the microscopic level:

$$
\boxed{
\text{small deterministic discrete/arithmetic dynamics}
}
$$

At intermediate scales:

$$
\boxed{
\text{coarse-graining + emergent measure + effective flow}
}
$$

At the IR:

$$
\boxed{
\text{smooth spacetime + complex amplitude + quantum dynamics}
}
$$

At macroscopic scales:

$$
\boxed{
\text{Born statistics + definite events + classical stability}
}
$$

At relativistic scales:

$$
\boxed{
\text{QFT + Lorentz symmetry + gauge structure + GR}
}
$$

The deepest success would be **compression**:

$$
\text{many observed laws}
\quad\longleftarrow\quad
\text{few microscopic principles}.
$$

A successful TOE candidate should therefore become simpler as its derivations improve.

If every unresolved mismatch creates another axiom, it is not converging.

---

# 30. What the "Elephant" Means

The Blind Men and the Elephant analogy gives the research program its epistemic humility.

BM may have touched the dynamics of definite quantum events.

IST may have touched discreteness, arithmetic structure, and the possibility that continuum quantum behavior is not fundamental.

AS may have touched universality, scale flow, and the idea that microscopic details can disappear while robust macroscopic laws remain.

None of these descriptions is entitled to become the elephant.

The elephant is whatever structure remains after:

- successful mathematics is retained;
- assumptions are removed;
- rival explanations are tested;
- translations between mathematical spaces are constructed;
- predictions survive held-out tests;
- and failed bridges are actually abandoned.

The program therefore does not ask:

> "How can we make BM, IST, and AS true together?"

It asks:

> **"What survives when all three are allowed to be wrong?"**

That question is more severe.

It is also more likely to reveal the common structure, if there is one.

---

# 31. Final Research Doctrine

The complete program can be compressed to twelve rules.

1. **No sacred cows.**
2. **Retain successful mathematics, not inherited metaphysics.**
3. **Separate state, measure, theory, and observable spaces.**
4. **Treat BM, IST, and AS as sources of targets and tools, not axiomatic packages.**
5. **Construct $F$ before theorizing about its consequences.**
6. **Construct the projection before claiming emergence.**
7. **Construct the amplitude factor before claiming quantum dynamics.**
8. **Separate equilibrium identification from relaxation.**
9. **Require regulator-, substrate-, and prime-robustness where appropriate.**
10. **Pre-register held-out prediction shapes before inspecting data.**
11. **Apply the Occam guillotine symmetrically to every component, including AS and EBP.**
12. **Let failure reduce the theory rather than enlarge it.**

The governing equation of the research program is therefore methodological:

$$
\boxed{
\text{Idea}
\rightarrow
\text{Debt}
\rightarrow
\text{Construction}
\rightarrow
\text{Adversarial Test}
\rightarrow
\text{Promotion or Subtraction}.
}
$$

And its governing epistemic rule is even simpler:

> **Ideas enter free. Promotion costs debt.**

---

# Appendix A — Reference Anchors

The following are methodological anchors, not endorsements of the tripartite program.

- M. Reginatto, *Derivation of the equations of nonrelativistic quantum mechanics using the principle of minimum Fisher information*, Phys. Rev. A 58, 1775 (1998).
- M. J. W. Hall and M. Reginatto, *Schrödinger equation from an exact uncertainty principle*, J. Phys. A 35, 3289–3303 (2002).
- A. Valentini, *Signal-locality, uncertainty, and the subquantum H-theorem. I*, Phys. Lett. A 156, 5–11 (1991).
- M. Hairer and J. C. Mattingly, *Ergodicity of the 2D Navier–Stokes equations with degenerate stochastic forcing*, Ann. Math. 164, 993–1032 (2006).
- G. S. Call and J. H. Silverman, *Canonical heights on varieties with morphisms*, Compositio Math. 89 (1993).
- E8 sphere-packing results and modular-form methods, with Viazovska's eight-dimensional solution as a benchmark for what genuine arithmetic optimality certificates look like.
- Established Bell-type Bohmian QFT constructions as the downstream relativistic benchmark.

---

# Appendix B — What This Version Corrects

Relative to v6.1, this revision specifically:

- converts "one flow" from an architectural claim into a flow-construction target;
- makes the self-audit explicit;
- names the smallest decisive microscopic computation;
- elevates amplitude-factor construction to a load-bearing gate;
- restores the Reginatto/Hall-Reginatto prior art;
- adds the p-adic additive-character route;
- adds the chaos-tameness/factor ladder;
- adds the canonical-height/Markov joint-satisfiability test;
- adds the Tensor-Factorization gate;
- resolves the Born-exact/fine-texture tension by defining the projection scale;
- restores a symmetric AS Occam kill condition;
- records the local-scale problem as debt rather than adding local $k(x)$;
- adds external red-team governance;
- restores a formal subtraction record;
- narrows the scope from "completed TOE" to "candidate pre-geometric substrate + IR quantum program with downstream TOE obligations";
- and applies the no-sacred-cows rule to this program's own preferred architecture.

---

# Final Statement

The strongest version of this research program is not the one that contains the most impressive mathematics.

It is the one that contains the **fewest assumptions capable of producing the most severe tests**.

At present, the central scientific object is still missing:

$$
\boxed{
F
\;\longrightarrow\;
\mu_*
\;\longrightarrow\;
\Pi
\;\longrightarrow\;
\psi
\;\longrightarrow\;
{\rm QM/BM}
}
$$

Construct those arrows or abandon the claim.

If the arrows can be built, the tripartite synthesis becomes a genuine candidate for deeper unification.

If they cannot, the program should end cleanly at the point of failure.

Either result is progress.

That is the purpose of the Elephant Bridge Protocol.
