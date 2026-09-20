---
title: "BM + IST + Asymptotic Safety"
subtitle: "Tripartite Synthesis v6.1 — Minimal, Typed, Adversarial TOE Candidate"
author: "Consolidated research architecture"
date: "20 September 2026"
geometry: "margin=1in"
fontsize: 10.5pt
linestretch: 1.08
---

# Status and Design Charter

**Status:** research program, not an established theory.

This document consolidates the strongest defensible material from the prior v5.1 tripartite synthesis and the newly attached DeepSeek, Z, and Gemini reviews. It is deliberately conservative about what has actually been derived. The purpose is to define the **smallest architecture worth trying to construct**, together with explicit theorem targets, numerical targets, empirical predictions, and death conditions.

> **Optimal design:** the design is at its optimum when no remaining feature can be removed without losing a load-bearing physical target, a necessary mathematical bridge, or a decisive falsification path.

The governing rule is:

$$
\boxed{\text{proposal}\;\to\;\text{typing}\;\to\;\text{formalization}\;\to\;\text{red-team}\;\to\;\text{construction/no-go}\;\to\;\text{numerics}\;\to\;\text{experiment}}
$$

A clean no-go is a successful research outcome.

---

# 0. Executive Thesis

The synthesis is **not** three theories placed beside one another. The intended candidate is one microscopic dynamical system whose effective descriptions separate by scale:

$$
\boxed{
(X,F)
\xrightarrow{\;\text{coarse-graining / RG}\;}
(\Gamma_k,\mu_k)
\xrightarrow{\;k\to \mathrm{IR}\;}
(\rho,S,\psi)
\xrightarrow{\;\text{recovery}\;}
\text{QM/BM} \;\text{and eventually QFT+GR+SM}
}
$$

with the roles:

- **IST** supplies the candidate microscopic ontology: finite/discrete, arithmetically structured, deterministic dynamics.
- **Asymptotic Safety** supplies the scale-flow discipline: theory space, effective action, RG flow, fixed points, relevant directions, universality and scheme tests.
- **Bohmian Mechanics** supplies the infrared quantum target: definite configurations, phase-guided dynamics, equivariance, contextual/nonlocal correlations, and a clear measurement ontology.
- **Valentini's H-theorem** supplies a possible relaxation mechanism, but not the definition of the equilibrium measure itself.

The decisive architectural change from earlier versions is this:

> **BM and AS are regime-specific machinery, not independent foundational ontologies; IST is only a candidate microscopic ontology. None of the three is allowed to smuggle in the desired IR result.**

This is a candidate architecture, not a claim that the architecture exists.

---

# 1. The Evaluation Framework Comes First

A TOE candidate should be evaluated before its narrative is trusted. The following framework is inherited from the project-wide chaos/complexity, Blind Men and the Elephant, and AI-evals analyses.

## 1.1 Nature's hard invariants

The original invariant list is useful, but it must be typed correctly. Not every item occurs in every physical system; some are universal constraints and some are conditional dynamical phenomena.

| ID | Invariant / structural demand | Discipline | Core question for this TOE |
|---|---|---|---|
| I1 | Sensitive dependence where chaotic dynamics occur | Chaos | Does the substrate reproduce bounded instability without making everything chaotic? |
| I2 | Physical invariant/equilibrium measure | Chaos | What measure is physical, and why? |
| I3 | Attractor structure where present | Dynamics | Which stable statistical states exist? |
| I4 | Structural stability / robustness | Dynamics | Do predictions survive admissible perturbations? |
| I5 | Universality classes | RG | Do microscopic variants converge to common IR behavior? |
| I6 | Fractal / multifractal structure where present | Chaos | What object is fractal, at what scale? |
| I7 | Bifurcation structure where parameters vary | Dynamics | Are qualitative transitions predicted rather than fitted? |
| I8 | Emergence | Complexity | How do macroscopic autonomous laws arise? |
| I9 | Renormalization / scale dependence | Complexity/QFT | How are descriptions at different scales related? |
| I10 | Criticality where physically relevant | Complexity | Why is a critical or near-critical regime present, if it is? |
| I11 | Computational structure | Information | What can the substrate compute and what is observable? |
| I12 | Information accounting | QM/QG | Is microscopic information conserved, transferred, or only coarse-grained away? |
| I13 | Entropy/information bounds | QG | Does the substrate respect known entropy bounds? |
| I14 | Thermodynamic arrow | Statistical physics | Can irreversibility emerge without being inserted? |

**Design correction:** positive Lyapunov exponents, attractors, bifurcations and global criticality are not universal properties of every physical system. The theory must support them where they are observed and must not falsely demand them everywhere.

## 1.2 Blind Men and the Elephant principles

1. **P1 — Partial correctness:** mature approaches may capture real structures without being complete.
2. **P2 — Surface incompatibility can hide deep compatibility:** apparently conflicting vocabularies may be descriptions at different levels.
3. **P3 — Deep compatibility must be demonstrated, not declared.**
4. **P4 — No totalizing claim:** the synthesis is itself provisional.
5. **P5 — Every foundational claim needs a projection to observables.**
6. **P6 — Structure outranks isolated vocabulary.**
7. **P7 — There is no privileged "view from nowhere"; the consistency of projections is what matters.

Operationally:

> Count what a theory actually touches correctly; penalize unsupported extrapolation; reward independent convergence.

## 1.3 AI-eval analogue for TOE candidates

| Eval | Requirement | Pass criterion |
|---|---|---|
| E1 | Benchmark reproduction | Known physics recovered within declared tolerance |
| E2 | Held-out prediction | At least one pre-specified prediction not used in construction |
| E3 | Adversarial/no-go audit | Bell, Haag, Coleman-Mandula, Weinberg-Witten and relevant no-go assumptions explicitly checked |
| E4 | Regression | New established observations do not require arbitrary reconstruction |
| E5 | Calibration | Every claim has an honest status label |
| E6 | Distribution shift | Domain shifts are either predicted or explicitly bounded |
| E7 | Ablation | Removing a claimed load-bearing component produces the expected degradation |
| E8 | Contamination | Construction and evaluation data are separated; no post-hoc knobs |
| E9 | Interpretability | Independent readers can verify the core derivations |
| E10 | Robustness | Predictions are stable or parameters are dynamically fixed |
| E11 | Falsifiability | A concrete experiment or observation can refute the claim |
| E12 | Kill conditions | Major claims have named failure conditions in advance |

## 1.4 Coherence criteria

| Criterion | Meaning |
|---|---|
| C1 | Internal consistency |
| C2 | Parsimony |
| C3 | Unification |
| C4 | Naturalness / absence of unexplained fine-tuning |
| C5 | Independent cross-checks |
| C6 | Honest accounting of debt |
| C7 | Survival under adversarial attack |

## 1.5 The five-clause "holds water" test

A candidate is **eval-sound** only when all five are present:

1. cheap gates run untuned;
2. exact statistics are invariants rather than fitted outputs;
3. red-team and ablation survive;
4. held-out **shape** predictions survive;
5. named death conditions are accepted in advance.

Eval-soundness is not proof of truth. It is proof that the research program is disciplined enough to deserve further testing.

## 1.6 T1 / T2 / T3 prediction tiers

- **T1:** parameter-free or near-parameter-free nulls and cheap exclusion tests.
- **T2:** predictions with a tightly bounded parameter budget and a prescribed functional form.
- **T3:** detailed structural fingerprints such as Diophantine, spectral, modular or colored-noise shape.

The synthesis must never use a T3 resemblance to compensate for a failed T1 gate.

---

# 2. Typed Mathematical Spaces — The Non-Negotiable Discipline

The largest source of confusion in earlier drafts was the silent identification of objects that live in different spaces. The revised synthesis separates them explicitly.

## 2.1 Microscopic state space

$$
X = \text{microscopic configurations},
\qquad
F:X\to X
$$

where \(F\) is the candidate deterministic update map, or its continuous-time analogue if one is proven legitimate.

## 2.2 Measure space

$$
\mathcal M(X)=\text{probability measures on }X.
$$

The physical measure, if it exists, belongs here. It is not the same object as the microscopic state space or an attractor in theory space.

## 2.3 Theory space

$$
\mathcal T=\{\Gamma\;\text{effective descriptions / actions / couplings}\}.
$$

The RG flow acts here:

$$
\mathcal R_s:\mathcal T\to\mathcal T,
\qquad
\frac{\partial \Gamma_s}{\partial s}=\beta[\Gamma_s],
$$

where \(s\) is a coarse-graining parameter. The relation between \(s\) and a conventional momentum scale \(k\) is a representation choice until derived.

## 2.4 Observable space

$$
\mathcal O=\text{experimentally accessible quantities}.
$$

The final projection is

$$
\Pi:\; X \;\text{or}\; \mathcal M(X) \longrightarrow \mathcal O.
$$

## 2.5 The invariant-set / fixed-point bridge

The proposed connection is **not**

$$
I_U = W^s(\Gamma_*)
$$

because these objects are typed differently.

The lawful target is something like

$$
\Phi:X\to\mathcal T,
\qquad
I_U \stackrel{?}{=} \Phi^{-1}\!\big(W^s(\Gamma_*)\big),
$$

or another explicitly constructed correspondence.

This is a theorem target, not an assumption.

---

# 3. Minimal Primitive Charter

Ruthless design asks how little must be assumed.

## Primitive P0 — Deterministic microscopic dynamics

There exists a microscopic state space \(X\) with a deterministic update law \(F\), together with an explicitly defined arithmetic/discrete structure.

**Status:** [CORE ASSUMPTION / PROGRAM CHARTER].

Everything else must be shown to emerge from it or be rejected.

## Primitive P1 — Controlled coarse-graining / RG structure

There exists a mathematically explicit map from microscopic descriptions to scale-dependent effective descriptions, sufficient to define an RG or Wilsonian flow.

**Status:** [THEOREM TARGET / CONSTRUCTION GATE].

The tripartite program does not assume that a discrete/arithmetic FRG exists merely because continuum FRG does.

## Primitive P2 — Observable projection

There exists a physically meaningful projection from microscopic/effective variables to experimentally accessible quantities.

**Status:** [THEOREM TARGET].

## Everything else is a requirement, not an axiom

Born statistics, smooth spacetime, p-adicity, the value of \(L\), the existence and location of a fixed point, the Schrödinger generator, the guidance law, Kramers behavior, spectral-dimension running, the gauge sector, Lorentz recovery, and the value of \(\kappa\) are **not** granted primitive status.

This is the strongest subtraction made in v6.1.

---

# 4. Bohmian Mechanics — First-Principles Re-foundation

## 4.1 Load-bearing BM content

Standard BM contributes several targets worth preserving:

- definite configurations/beables;
- phase-guided motion;
- a wave-function-based effective dynamics;
- equivariance / quantum equilibrium;
- contextual, nonlocal correlation structure compatible with no signaling;
- a transparent measurement ontology.

## 4.2 What is demoted

| Standard BM assumption | Tripartite treatment | Status |
|---|---|---|
| Continuous configuration space is primitive | Derived IR target | [DERIVED-TARGET] |
| Schrödinger evolution is fundamental | IR effective equation | [DERIVED-TARGET] |
| Guidance equation is fundamental | IR limit of microscopic dynamics | [THEOREM TARGET] |
| \(\rho=|\Psi|^2\) is an initial postulate | Equilibrium-measure gate | [GATE] |
| Wave function must be a primitive physical field | Effective/nomological/archimedean shadow until proven otherwise | [OPEN] |
| Preferred foliation inserted by hand | Derive if required; test equilibrium observability | [OPEN] |
| Fixed particle number is fundamental | Minimal fixed-N sector first; Bell-type/QFT extension later | [STAGED] |
| Quantum potential is primitive | Derived from effective \(\rho,S\) or effective action | [THEOREM TARGET] |

## 4.3 Asymptotic Bohmianity

The useful DeepSeek insight survives in a stricter form:

> Bohmian mechanics may be the infrared fixed description of a more general microscopic dynamics.

We do **not** assume a particular running guidance law or an explicit correction exponent until the microscopic model produces one.

The target remains

$$
v \longrightarrow \frac{\nabla S}{m}
\qquad (k\to k_{\mathrm{IR}}).
$$

Likewise, asymptotic unitarity is a target:

$$
\mathcal L_{k>k_\xi}
\to
\mathcal L_{k\ll k_\xi}^{\mathrm{Schr}}.
$$

A Lindblad or generalized Langevin bridge is an optional construction route, not a founding axiom.

## 4.4 Kramers / generalized Langevin correction

The three-AI review correctly noticed that a purely overdamped first-order flow may be too restrictive. The minimal retained statement is therefore the generalized second-order target

$$
M\ddot q(t)
=
-\nabla U(q,t)
-\int_0^t K(t-s)\dot q(s)\,ds
+\eta(t),
$$

where both the memory kernel \(K\) and colored noise \(\eta\) must be derived from tracing hidden degrees of freedom.

The familiar Kramers equation appears only when a controlled Markovian approximation exists.

Thus:

$$
\text{generalized memory dynamics}
\to
\text{Kramers limit}
\to
\text{overdamped BM guidance}
$$

is a staged target, not an added ontology.

## 4.5 Quantum potential target

With

$$
\psi=\sqrt{\rho}\,e^{iS/\kappa},
$$

the Bohmian quantum potential is

$$
\boxed{
Q[\rho]=-
\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt{\rho}}{\sqrt{\rho}}
}.
$$

The synthesis must derive this structure from the IR generator or effective action. A mathematical analogy to Fisher information or an arithmetic height is not enough.

## 4.6 Configuration space is itself a gate

A common mistake is to derive trajectories while silently assuming the continuum configuration space on which the wavefunction lives.

The real target is therefore

$$
X
\xrightarrow{\Phi,\,\Pi}
\mathcal C_{\mathrm{IR}}
$$

where \(\mathcal C_{\mathrm{IR}}\) is a smooth configuration space with the regularity needed for the Schrödinger/Bohmian equations.

This extends the original Projection Problem rather than solving it by changing notation.

---

# 5. Invariant Set Theory — First-Principles Re-foundation

## 5.1 Defensible IST core

The minimum defensible content is:

1. fine-grained reality may be intrinsically discrete / finite-information;
2. physical realizability may be a proper subset of mathematical possibility;
3. continuous observed physics may emerge from a deterministic structured substrate.

The familiar IST mechanisms are hypotheses about how to implement these three statements, not guaranteed deductions from them.

## 5.2 Information capacity \(L\)

Replace a hand-set bit count with a functional of the microscopic flow:

$$
L \stackrel{?}{=} C[\Gamma_k],
$$

where \(C\) is an information/entropy monotone if one exists for the substrate.

The use of a c/a-like RG quantity is a research direction, not a license to call \(L\) derived before a map is built.

**Kill condition:** no regulator-stable relation without fitted parameters.

## 5.3 Definability via arithmetic dynamics

Replace a coordinate-dependent rational/irrational rule by an intrinsic dynamical candidate:

$$
\hat h(x)=0
\quad\text{for the dynamically selected preperiodic / low-complexity sector}
$$

when the substrate map is of a class for which canonical heights exist.

The required sequence is:

$$
F
\to
\hat h
\to
\text{preperiodic locus}
\to
\text{physical admissibility}
$$

The Niven content becomes a faithfulness test, not a theorem imported as ontology.

## 5.4 p-adicity is derived, not declared

The candidate chain is

$$
F
\to
\text{Markov / symbolic partition}
\to
\text{ultrametric}
\to
p\text{-adic or adelic structure if selected}.
$$

A Markov partition is not guaranteed for an arbitrary dynamical system; its existence must be demonstrated for the actual substrate.

### Prime-universality gate

Test admissible prime choices and/or branching structures. Outcomes may include:

- prime independence at the IR level;
- unique dynamic prime selection;
- genuinely distinct prime sectors;
- full adelic necessity.

A chosen prime is never justified merely by convenience.

## 5.5 The fixed update rate is removed

No literal one-bit-per-Planck-time law is granted.

Instead test

$$
\Gamma_{\mathrm{mix}}
\stackrel{?}{\sim}
\Gamma_{\mathrm{RG}}
$$

only after both quantities are independently defined.

## 5.6 Fractality becomes a typed observable

"Fractal" must say what is fractal:

- a microscopic set;
- the support of a measure;
- density texture;
- a critical surface;
- a spectral dimension;
- a multifractal spectrum.

A useful regime hypothesis is

$$
\text{IR: regular/full-dimensional observable measure},
\qquad
\text{UV: arithmetic fine structure possibly multifractal},
$$

but the texture type is a computation, not a postulate.

## 5.7 Measurement independence

The synthesis **does not** get to win Bell by merely renaming measurement independence failure.

The target is:

$$
\rho(\lambda|a,b)\neq\rho(\lambda)
$$

only if that dependence is actually generated by the microscopic construction and survives the Wood-Spekkens audit without enabling signaling.

This is a gate, not a free explanatory label.

---

# 6. Asymptotic Safety — Adopt the Flow, Not the Dogma

## 6.1 What standard AS genuinely contributes

The durable methodological contribution is:

- effective-average-action methods;
- functional RG flow;
- fixed-point analysis;
- critical exponents and relevant directions;
- universality;
- regulator/truncation diagnostics;
- parameter counting through the UV critical surface.

## 6.2 Discrete functional RG (dFRG)

The synthesis proposes a substrate-level implementation:

$$
\boxed{\text{dFRG} = \text{AS methodology implemented on the microscopic substrate}}
$$

This is a research program. Continuum FRG does not automatically prove that the required discrete/arithmetic flow exists.

## 6.3 Fixed point versus invariant set

Let the RG flow be

$$
\partial_s\Gamma_s=\beta[\Gamma_s]
$$

with a fixed point

$$
\beta[\Gamma_*]=0.
$$

The tripartite target is a correspondence between state-space structure and theory-space critical structure, not an identification by fiat.

One possible target is

$$
I_U=\Phi^{-1}(W^s(\Gamma_*)),
$$

where \(W^s(\Gamma_*)\) is the stable/critical manifold of the fixed point.

The bridge is a major theorem target.

## 6.4 Regulator and scheme independence

A central correction from the AI reviews is that AS does **not** bypass scheme dependence.

Intermediate couplings and truncation estimates can move with the regulator. The physical target is regulator/scheme stability of **universal observables**.

The required gate is therefore:

$$
\boxed{
\text{derived observable stable under legitimate regulator/truncation changes}
}
$$

A quantity that moves materially under the allowed scheme family is not a physical prediction merely because one truncation produced a clean number.

## 6.5 Fixed-point existence is a gate

The synthesis cannot presuppose a non-Gaussian fixed point.

Gate F0:

1. construct the flow;
2. demonstrate a candidate fixed point;
3. test convergence across truncations/regulators;
4. compute the stability matrix;
5. identify relevant directions;
6. test universality.

## 6.6 Relevant directions and parameter economy

At a fixed point,

$$
\delta g_i(s)
\sim
\sum_a c_a e^{-\theta_a s}v_a,
$$

with relevant directions characterized by the appropriate positive critical exponents under the chosen convention.

The physical parameter count should emerge from this structure rather than from an arbitrary parameter table.

A core aspiration is **low effective codimension**, not a numerically assumed value.

## 6.7 Spectral dimension

Use the diffusion definition

$$
d_s(\sigma)=-2\frac{d\ln P(\sigma)}{d\ln\sigma},
$$

where \(P(\sigma)\) is the return probability.

The Alexander-Orbach-style relation

$$
d_s=\frac{2d_H}{d_w}
$$

is allowed only when the required diffusion/fractal scaling hypotheses are actually satisfied.

The tripartite prediction is **not** "\(d_s=2\) because other approaches do it." It is a held-out gate asking whether

$$
d_s^{\mathrm{IR}}\approx4,
\qquad
d_s^{\mathrm{UV}}\approx2
$$

or some different flow emerges from the substrate.

## 6.8 Metric emergence

The continuum metric is treated as an IR effective structure unless the microscopic construction forces otherwise.

The target is

$$
X,F,\mathcal R_s
\to
g_{\mu\nu}^{\mathrm{IR}}
$$

with diffeomorphism invariance, causal structure, and the equivalence principle recovered in the tested regime.

The metric is not removed by rhetoric; it is removed from the list of microscopic primitives.

## 6.9 What we do not adopt from the Gemini review

The proposed local scale field

$$
k^2(x)=\alpha\left|\Box R/R\right|+\beta R_{\rm curv}
$$

is excluded from the core. It introduces additional coefficients \(\alpha,\beta\), depends on the already-unrecovered amplitude field, and converts an RG bookkeeping parameter into an asserted local observable without a derivation.

That is the opposite of the subtraction rule.

---

# 7. The Unified Flow Architecture

## 7.1 One microscopic source, multiple effective descriptions

The desired chain is

$$
(X,F)
\xrightarrow{\Phi}
\Gamma_s
\xrightarrow{\mathcal R_s}
\Gamma_{\mathrm{IR}}
\xrightarrow{\Pi}
(\psi,\rho,S,\mathcal C_{\mathrm{IR}}).
$$

Nothing in this diagram is allowed to collapse two spaces together merely for narrative convenience.

## 7.2 The internal-environment principle

If the microscopic theory is globally closed, the only environment available for a reduced subsystem is the hidden part of that same microscopic theory.

The intended trace structure is therefore

$$
\mathcal H_{\rm total}
=
\mathcal H_{\rm visible}\otimes\mathcal H_{\rm hidden},
$$

or its appropriate discrete analogue, with the hidden arithmetic degrees of freedom serving as the effective bath.

This is a design principle, not yet a constructed tensor factorization of the actual substrate.

## 7.3 Derived stochasticity

If a reduced stochastic law appears, it must be derived from deterministic hidden dynamics:

$$
F_{\rm micro}
\longrightarrow
K(t)
+\eta(t)
\longrightarrow
\text{effective reduced dynamics}.
$$

Noise and friction must obey a common fluctuation-dissipation structure where the assumptions of such a relation hold.

## 7.4 Three notions of flow

Keep distinct:

1. microscopic iteration time \(n\) or microscopic dynamics;
2. RG/coarse-graining scale \(s\) or \(k\);
3. emergent physical time \(t\).

No identification such as \(s=t\) is permitted without proof.

---

# 8. The Projection Problem — Rewritten in Full

The original projection problem is actually three linked problems:

### P0 — State projection

Can the microscopic state produce the correct effective configuration manifold?

$$
X\to \mathcal C_{\rm IR}.
$$

### P1 — Measure projection

Can the microscopic invariant/equilibrium measure project to a regular continuum measure?

$$
\Pi_\#\mu_*
=\rho(x)\,dx?
$$

### P2 — Dynamical projection

Does the projected generator converge to the Schrödinger/Bohmian generator?

$$
\mathcal L_s
\to
\mathcal L_{\rm Schr/BM}?
$$

This decomposition prevents a smooth density or a pretty continuum visualization from being mistaken for a derivation of dynamics.

---

# 9. Probability and the Born Rule — One Architecture, Separate Gates

## 9.1 B0 — Existence and uniqueness of the physical equilibrium measure

Need

$$
\exists!\,\mu_*.
$$

RDS methods such as strong Feller / irreducibility / asymptotic coupling can be used where their exact hypotheses are verified. The theorem class is imported; the substrate's membership is not assumed.

## 9.2 B1 — Projection regularity

Need

$$
\Pi_\#\mu_*=\rho(x)\,dx
$$

with enough regularity, ideally finite Fisher information,

$$
I_F[\rho]
=
\int \rho |\nabla\ln\rho|^2 dx
<\infty.
$$

If this fails in the required regime, the quantum kinetic structure cannot be established by the intended route.

## 9.3 B2 — Born identification

Need

$$
\boxed{\Pi_\#\mu_*=|\psi|^2}
$$

or a mathematically equivalent construction.

This is the actual Born gate.

## 9.4 B3 — Relaxation

Only after B2 is established does Valentini's H-theorem enter as the relaxation mechanism:

$$
H[P|\rho]
=
\int P\ln\frac{P}{\rho}\,dx,
$$

with an appropriate coarse-grained version satisfying

$$
\bar H(t)\le \bar H(0)
$$

under Valentini's hypotheses.

The H-theorem therefore supplies a **path to equilibrium**, not the definition of the destination.

## 9.5 T-Born and Valentini-H

Two mathematically distinct routes remain:

- **Arithmetic/equidistribution route:** a theorem target for equilibrium identification and convergence.
- **Valentini route:** a relaxation theorem for non-equilibrium initial data.

They should not be conflated, and neither is granted primacy before the relevant calculations.

A useful final architecture is:

$$
\boxed{
\text{equilibrium identification}
\;+
\text{relaxation mechanism}
\;+
\text{initial-condition accounting}
}
$$

## 9.6 The initial-condition obligation

A relaxation theorem does not explain why the universe begins away from equilibrium.

The program therefore tracks an explicit debt:

$$
\text{UV preparation}
\to
P_0\neq \rho_*
\to
P_t\to\rho_*.
$$

If no principled non-equilibrium initial condition exists, the H-theorem is incomplete as a cosmological explanation.

---

# 10. Linearity, Guidance and the Quantum Potential

## 10.1 Effective wavefunction

The target decomposition is

$$
\psi=\sqrt{\rho}\,e^{iS/\kappa}.
$$

No claim that \(\psi\) is a microscopic primitive is required.

## 10.2 Guidance target

The target guidance field is

$$
\boxed{v=\frac{\nabla S}{m}}
$$

in the appropriate nonrelativistic sector.

The derivation must show why the same \(S\) that enters the effective phase also generates the velocity field.

## 10.3 Schrödinger target

The target linear generator is

$$
i\kappa\partial_t\psi
=
-\frac{\kappa^2}{2m}\nabla^2\psi
+V\psi
$$

in the minimal nonrelativistic limit.

Linearity is a separate gate. A nonlinear microscopic flow can generate a linear effective theory near a fixed point, but that statement is not itself a derivation.

## 10.4 Quantum potential

Substitution of the polar form yields the target

$$
Q
=
-\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt\rho}{\sqrt\rho}.
$$

The program must derive it through one explicit route:

1. effective action / generator;
2. continuity-plus-phase system;
3. a mathematically controlled alternative.

Arakelov/height interpretations are kept in the workbench only.

---

# 11. Composition, Entanglement, Bell and Measurement

## 11.1 Gate C — emergent composition

The one-particle theory is insufficient.

Need an emergent composition rule with

$$
\mathcal H_{AB}
\cong
\mathcal H_A\otimes\mathcal H_B
$$

or the appropriate effective analogue in the constructed theory.

The gate must show:

1. weakly coupled subsystems factorize;
2. interactions generate entangled states;
3. Bell correlations emerge;
4. no-signaling survives.

## 11.2 Measurement independence is not renamed away

A Bell-capable synthesis must pass a Wood-Spekkens-style audit.

The question is not "does the theory violate measurement independence?" but:

> **What microscopic mechanism generates the setting dependence, and why does it not become usable signaling?**

## 11.3 No-signaling constraint

For any correction \(\delta\rho\), require the marginal correction to vanish:

$$
\int \delta\rho(a,b,\lambda)\,d\lambda_B=0
$$

in the appropriate notation, while joint correlators may change.

The fine structure is therefore allowed to be visible in **correlations**, not in locally controllable marginals.

## 11.4 Foliation

A preferred foliation is not to be inserted because standard BM uses one.

Potential targets include:

- an emergent causal/trace direction;
- a hidden foliation whose equilibrium imprint vanishes;
- a covariant Bell-type field-theoretic formulation.

No one route is promoted without proof.

---

# 12. QFT and Gravity Recovery Gate

A TOE candidate that reproduces only nonrelativistic BM has not solved the stated problem.

## QFT-R0 — Relativistic effective field theory

Recover an effective QFT with the correct degrees of freedom and particle-number-changing processes.

## QFT-R1 — Lorentz symmetry

Recover Lorentz invariance to tested precision in the IR.

## QFT-R2 — Gauge structure

Recover the observed gauge structure without adding unforced dimensions or fitted group embeddings.

## QFT-R3 — No-go audit

The synthesis must explicitly audit the assumptions of:

- Haag's theorem;
- Coleman-Mandula;
- Weinberg-Witten;
- any additional theorem triggered by the actual construction.

A UV cutoff is not an automatic escape hatch. The relevant assumptions must be identified, shown to fail or hold, and the correct IR assumptions recovered.

## QFT-R4 — Euclidean to Lorentzian continuation

If the core construction is Euclidean/FRG based, establish a viable continuation to real-time dynamics:

$$
\text{Euclidean effective theory}
\to
\text{Lorentzian unitary / causal theory}.
$$

Without this, a BM trajectory law in real time is not grounded.

## QFT-R5 — Matter compatibility

The gravitational UV fixed point must remain compatible with the emergent matter sector. This is a named dependency, not an assumption.

## Gravity scope discipline

The synthesis does **not** yet claim that:

- the cosmological constant is solved;
- dark matter is derived;
- modified gravity is derived;
- a gravitational acceleration scale is derived;
- the Standard Model gauge group is uniquely selected.

These are recovery targets.

---

# 13. Spectral Dimension and the UV/IR Geometry Question

The synthesis permits a dimensional-flow signature but does not hard-code its value.

Compute

$$
d_s(\sigma)=-2\frac{d\ln P(\sigma)}{d\ln\sigma}.
$$

Possible outcomes are discriminating:

- \(d_s\approx4\to2\): supports the intended dimensional-flow branch;
- flat \(d_s\approx4\): weakens the arithmetic-crossover picture;
- another dimension: must be incorporated honestly.

The popular \(d_s\approx2\) value is therefore an **empirical/theorem target**, not an axiom.

---

# 14. E8 — Mathematics Retained, Baggage Removed

## 14.1 Permanently rejected

The synthesis does **not** use Lisi-style E8 unification of all gauge bosons and fermions as a foundational mechanism.

The known chirality obstruction is not removed by simply moving the fermions into an E8-related CFT representation space.

## 14.2 The 248 clarification

\(248\) is the dimension of the E8 Lie algebra — the number of generators.

It is **not** an observed count of the particles in nature.

The synthesis must never claim that E8 explains "why there are 248 particles."

## 14.3 Surviving E8 role

A narrower, technically meaningful role survives:

> **E8 may organize an internal information code if an independent derivation selects an eight-dimensional coding sector and if an entropy/packing principle forces the optimal lattice.**

The mathematical coincidence cluster is attractive:

- even unimodularity;
- optimal sphere packing in eight dimensions;
- coding-theoretic structure;
- modular-form structure.

But none of these facts alone selects E8 physically.

## 14.4 Selection principle requirement

The gate is

$$
\text{entropy/information constraint}
\to
\text{unique optimal code}
\to
E_8
$$

not

$$
E_8
\to
\text{retrofit everything}.
$$

## 14.5 Optional UV-CFT module

Only if the spectral-dimension gate produces a genuine two-dimensional UV description, and that theory is demonstrably lattice/rational type, may an E8 lattice CFT become a downstream candidate.

Chirality would still require a separate mechanism. The naive "E8 CFT representation theory fixes chirality" route remains rejected.

## 14.6 Modular fingerprint

If E8 survives selection, the highest-value empirical role is not a particle count but a **held-out statistical fingerprint** in the predicted fine structure.

A modular-form-like pattern in the deviation spectrum is valuable precisely because it is difficult to obtain accidentally and can be tested without fitting its detailed amplitude after the fact.

---

# 15. What Is Explicitly Not in the Core

The following are deliberately excluded from the core unless a prior gate forces them back:

- an externally defined local \(k(x)\) built from \(\Box R/R\);
- Kaluza-Klein-like physical extra dimensions;
- a separate internal-state geometry that is not part of the established IST program;
- Poisson sprinkling as an automatic solution of the foliation problem;
- ER=EPR as an ontological identification;
- holography as an assumption rather than a derived relation;
- Arakelov-height quantum potential as the primary derivation;
- Hecke/Langlands machinery without a gate that needs it;
- "one universal flow" as a literal identity across every mathematical space;
- E8 as gauge-group/particle census;
- fixed \(\xi_0=\ell_P\) without computation;
- fixed \(d_s=2\) without computation;
- guaranteed Born equilibrium without B0-B2;
- any claim that a theorem class automatically applies to the constructed substrate.

The purpose of the workbench is to keep useful ideas without turning them into ontology.

---

# 16. Cross-Substrate Universality

This is one of the strongest additions to the AS-inspired architecture.

Construct microscopically distinct substrates

$$
(X_1,F_1),\quad (X_2,F_2),\quad (X_3,F_3)
$$

and ask whether they flow toward the same effective universality class:

$$
(X_i,F_i)
\xrightarrow{\mathcal R}
\Gamma_{\mathrm{IR}}^{(i)}
\quad\Rightarrow\quad
\Gamma_{\mathrm{IR}}^{(i)}\sim\Gamma_{\mathrm{IR}}^{*}.
$$

If many different microscopic realizations converge to the same BM/QFT infrared structure, the universality claim gains real content.

If one specially tuned substrate alone works, the theory's universality claim weakens sharply.

This is the preferred AS-style answer to the question "does the microscopic detail really matter?"

---

# 17. One Scale, Many Derived Scales

The design aims for one coarse-graining scale variable rather than unrelated phenomenological knobs.

Possible derived scales include:

| Quantity | Role | Status |
|---|---|---|
| \(k\) or \(s\) | RG / coarse-graining variable | Structural |
| \(\xi_0\) | bath-decoupling crossover | To be computed |
| \(\ell_P\) | gravitational scale | Comparison target, not identification |
| \(Q\) | arithmetic resolution relevant to Bell geometry | To be derived |
| \(\epsilon_0\) | observational coarse-graining | Instrument-defined |

An equality such as \(\xi_0=\ell_P\) is a **cross-check**, never an input.

---

# 18. Parameter Discipline

Common parameters must be frozen before empirical evaluation.

A useful shared set remains

$$
(c,\xi_0,\kappa)
$$

only if independent derivations show that these quantities are truly distinct and physical.

Rules:

1. no parameter may be introduced solely to rescue a failed prediction;
2. a parameter appearing in two sectors must be cross-validated independently;
3. regulator, truncation and prime choices cannot silently become hidden fitting parameters;
4. relevant-direction counting must ultimately explain how many genuinely free parameters remain.

The parameter goal is not "as few numbers as possible" in the abstract. It is **as few unexplained choices as possible**.

---

# 19. Empirical Program — Shape Before Amplitude

## 19.1 Mesoscopic interferometry

The high-value question is not merely whether a larger mass can interfere. It is whether the synthesis predicts a distinctive scaling law.

The target is a theory-specific visibility function

$$
V=V(N,m,t,\ldots)
$$

with a pre-registered shape/exponent and a fixed nuisance-parameter policy.

No exponent may be fitted after inspecting the data.

## 19.2 Bell / Born deviations

The synthesis permits small arithmetic deviations in non-equilibrium or fine-structure regimes only if they survive no-signaling and are quantitatively derived.

The most valuable signature is functional shape rather than a single amplitude.

## 19.3 Noise spectrum

If the hidden arithmetic bath produces colored noise, the spectrum is a natural T2/T3 handle.

The test should distinguish the predicted kernel from phenomenological white-collapse noise such as a generic spontaneous-localization spectrum.

## 19.4 Relic non-equilibrium

Possible cosmological handles may include relic sectors, primordial gravitational modes, or other systems with long-lived insufficient relaxation.

The exact observable must be derived before any data are used.

## 19.5 Spectral dimension

The UV/IR running is an internal mathematical prediction and may also become a quantum-gravity signature.

## 19.6 E8 modular fingerprint

If E8 is independently selected, the modular fingerprint is a held-out structural prediction, not a decorative mathematical analogy.

---

# 20. The Elephant Scan — Six Mandatory Bridge Audits

### S1 — Define

Do \(X,F,\mu,\mathcal T,\mathcal R,\Pi\) exist mathematically?

### S2 — Dynamics

Is the microscopic dynamics explicit and nontrivial?

### S3 — Measure

Does a unique physically relevant invariant/equilibrium measure exist?

### S4 — Projection

Does a regular continuum configuration space and density emerge?

### S5 — Quantum structure

Do Born, Schrödinger, guidance and quantum potential all emerge from the same effective object?

### S6 — Composition

Does the construction support tensor products, entanglement, Bell correlations and no signaling?

### S7 — Recovery

Does the same architecture recover QFT, gauge structure, Lorentz invariance and GR in the tested regime?

A TOE candidate that passes S1-S4 but fails S5 is a mathematically interesting substrate, not yet a quantum theory.

---

# 21. Adversarial Gates and Kill Conditions

| ID | Gate | Kill condition |
|---|---|---|
| D0 | Discrete/continuous consistency | The substrate secretly assumes the target continuum dynamics |
| D1 | dFRG construction | No valid substrate-level coarse-graining/RG flow |
| D2 | Fixed point | No regulator-stable physically meaningful fixed point |
| D3 | Measure | No unique physical invariant/equilibrium measure |
| D4 | Regularity | Projected measure too singular for the intended quantum generator |
| D5 | Born | Projected equilibrium measure is not \(|\psi|^2\) |
| D6 | Linearity | IR generator is not the required linear Schrödinger class |
| D7 | Guidance | Velocity field does not reduce to \(\nabla S/m\) |
| D8 | Composition | No emergent tensor-product / entanglement structure |
| D9 | Safety | Signaling or unacceptable Lorentz violation in the claimed domain |
| D10 | Scheme | Derived universal numbers move materially under legitimate scheme changes |
| D11 | Universality | Different microscopic substrates require hidden retuning to land in the same IR class |
| D12 | QFT | No credible relativistic effective field theory |
| D13 | No-go audit | Required theorem assumptions are contradicted without a controlled mechanism |
| D14 | E→L | No viable Euclidean-to-Lorentzian continuation |
| D15 | Matter | AS fixed point incompatible with emergent matter sector |
| D16 | Parameter | Free knobs proliferate at each mismatch |
| D17 | E8 | E8 appears only by analogy, without a physical selection principle |
| D18 | Empirical | Pre-registered held-out shape prediction fails |
| D19 | Internal Occam | Arithmetic T3 fingerprints disappear while a simpler non-arithmetic theory explains the data |

**Occam guillotine:** if the arithmetic sector adds no distinctive successful prediction, the cheaper effective theory wins.

---

# 22. Research Program Ordered by Information Gain

## Phase 0 — Cheap kill gates

1. compile all existing interferometry/Bell constraints;
2. define the allowed parameter window;
3. test whether any claimed deviation survives present bounds without retuning.

**Gate:** if no parameter window exists, stop.

## Phase 1 — Substrate audit

1. make \(X\) and \(F\) explicit;
2. perform the scaling / Bernoulli-convolution regularity audit;
3. test the canonical-height replacement;
4. test for a Markov partition;
5. run prime-universality checks.

## Phase 2 — Measure and probability

1. construct \(\mu_*\);
2. test uniqueness/mixing hypotheses;
3. test absolute continuity and Fisher regularity;
4. test Born identification;
5. address initial conditions.

## Phase 3 — dFRG

1. construct the substrate coarse-graining operator;
2. construct the effective action / theory space;
3. search for fixed points;
4. compute critical exponents;
5. perform regulator/truncation sensitivity tests;
6. compare fixed-point data with state-space invariants.

## Phase 4 — IR quantum construction

1. derive configuration-space emergence;
2. derive the Schrödinger generator;
3. derive guidance;
4. derive the quantum potential;
5. test the asymptotic unitarity limit.

## Phase 5 — Composition and Bell

1. construct two weakly coupled subsystems;
2. derive tensor composition;
3. generate entanglement;
4. reproduce Bell correlations;
5. run Wood-Spekkens and no-signaling audits.

## Phase 6 — Relativistic/gravitational recovery

1. Euclidean→Lorentzian continuation;
2. Lorentz recovery;
3. QFT gauge structure;
4. Haag/Coleman-Mandula/Weinberg-Witten audit;
5. AS-with-matter compatibility;
6. GR / diffeomorphism / equivalence-principle recovery.

## Phase 7 — Phenomenology

Only after earlier structural gates pass:

- interferometry scaling;
- colored-noise spectrum;
- Diophantine texture;
- spectral dimension;
- cosmological/relic signatures;
- optional E8 modular fingerprint.

---

# 23. Live Claims Ledger

| CLAIM_ID | Claim | Status | Dependency |
|---|---|---|---|
| TRI-ARCH-1 | One substrate + one controlled coarse-graining infrastructure can produce distinct IR regimes | [ARCHITECTURE] | P1/P2 |
| TRI-RG-1 | A discrete/arithmetic RG flow exists | [SPECULATIVE-SHAPE] | Phase 3 |
| IST-L-1 | Capacity \(L\) is a derived flow functional | [SPECULATIVE-SHAPE] | dFRG |
| IST-H-1 | Canonical height selects a physical/preperiodic sector | [SPECULATIVE-SHAPE] | explicit \(F\) |
| IST-M-1 | Symbolic/Markov structure derives the ultrametric | [SPECULATIVE-SHAPE] | explicit \(F\) |
| IST-P-1 | A dynamically selected prime or adelic structure exists | [SPECULATIVE-SHAPE] | IST-M-1 |
| AS-FP-1 | A physically meaningful fixed point exists | [THEOREM/NUMERICAL TARGET] | dFRG |
| AS-REL-1 | Relevant directions explain parameter count | [DERIVED-TARGET] | AS-FP-1 |
| AS-SD-1 | Spectral dimension runs between IR and UV values | [HELD-OUT GATE] | dFRG |
| BM-CFG-1 | Smooth IR configuration space emerges | [THEOREM TARGET] | projection |
| BM-LIN-1 | Schrödinger linearity is an IR fixed behavior | [DERIVED-TARGET] | Phase 4 |
| BM-GUID-1 | \(v=\nabla S/m\) is recovered | [THEOREM TARGET] | Phase 4 |
| BM-Q-1 | Standard Bohmian quantum potential is recovered | [THEOREM TARGET] | Phase 4 |
| BORN-1 | Unique equilibrium measure projects to \(|\psi|^2\) | [THEOREM TARGET] | Phase 2 |
| VAL-H-1 | Non-equilibrium relaxes under the relevant hypotheses | [EXISTING TOOL / TARGET] | mixing |
| COMP-1 | Tensor-product composition emerges | [THEOREM TARGET] | Phase 5 |
| BELL-1 | Bell correlations emerge without signaling | [THEOREM TARGET] | COMP-1 |
| QFT-1 | Relativistic QFT emerges | [RECOVERY TARGET] | Phase 6 |
| GR-1 | General-relativistic IR behavior emerges | [RECOVERY TARGET] | Phase 6 |
| E8-1 | E8 is selected as a coding structure | [SPECULATIVE-SHAPE] | information-selection gate |
| E8-2 | E8 modular fingerprint appears | [HELD-OUT TARGET] | E8-1 |

---

# 24. Adjudication of the Three Attached AI Reviews

## 24.1 DeepSeek — net valid

**Kept:**

- the regime architecture: arithmetic UV, RG bridge, Bohmian IR;
- asymptotic Bohmianity as a target;
- asymptotic unitarity as a companion target;
- generalized/Kramers dynamics as a possible intermediate effective description;
- RDS mixing as a lower-cost mathematical route than demanding uniform hyperbolicity everywhere;
- the information-capacity-to-flow strategy;
- canonical heights as an intrinsic replacement candidate for coordinate rationality cuts;
- Markov-partition derivation of an ultrametric;
- the initial-condition debt for Valentini relaxation;
- explicit parameter/recovery/evaluation discipline.

**Demoted or corrected:**

- "\(I_U\) is the critical surface" became a typed correspondence target;
- "Kramers is the fundamental reduced law" became an optional effective limit;
- "T-Born is primary" became an undecided empirical/theorem gate;
- exact Born-correction formulas were not promoted without construction.

**Reason:** the architecture is useful; several strongest statements were one inferential step ahead of their evidence.

## 24.2 Z — net valid

**Kept:**

- ruthless subtraction rule;
- dFRG as the AS transplant;
- regulator/scheme-independence gate;
- spectral-dimension gate;
- cross-substrate universality;
- prime-universality test;
- E8 inversion from gauge unification to information organization;
- modular fingerprints;
- kill-test-first sequencing;
- internal Occam control.

**Demoted:**

- Kramers as a core axiom;
- equality of \(\xi_0\) and \(\ell_P\) as a computed check rather than a statement;
- exact one-rate claims;
- fixed-point / invariant-set identification.

## 24.3 Gemini — heavily pruned

The attached Gemini proposal uses the term "Internal State Theory" and introduces a Kaluza-Klein-style bundle geometry, local \(k(x)\), and geodesic flow. Those are **not** the minimal Invariant Set Theory program developed in this synthesis.

Rejected from the core:

- replacing BM configuration space with a finite-dimensional local internal bundle;
- treating the wavefunction as a frame-bundle section without deriving the construction;
- the local scale definition \(k^2\sim|\Box R/R|+R\);
- claims that the fixed point automatically removes singular forces;
- automatic gauge-coupling unification;
- automatic cosmological-constant relaxation.

The reason is not that bundle geometry is mathematically illegitimate. It is that these additions introduce substantial new structure **before** the projection, flow, measure and recovery gates have been solved. They therefore fail the subtraction rule.

The useful conceptual residue is narrower:

> internal geometry, if later forced by the gauge/QFT recovery gate, should be introduced as an emergent structure of the effective theory rather than as an independent microscopic ontology.

---

# 25. Evals as a Standing Constitution for the TOE Candidate

The evaluation framework is not an appendix to be consulted after construction. It is the constitution governing construction.

## 25.1 Construction/evaluation separation

Every empirical claim must be tagged at conception time:

$$
\text{input} \neq \text{training/fit datum} \neq \text{hold-out test}.
$$

No prediction is allowed to become a moving target after the data are examined.

## 25.2 Ablation protocol

For each core component ask:

- Remove arithmetic structure. What prediction disappears?
- Remove the RG flow. What bridge disappears?
- Remove BM. Which IR measurement/ontology target disappears?
- Remove the H-theorem. What relaxation prediction disappears?
- Remove E8. Which held-out fingerprint disappears, if any?

If removing a component changes no prediction and closes no logical gap, it does not belong in the core.

## 25.3 Cross-validation protocol

Key quantities should have independent routes where practical:

- \(c\): dynamical and RG routes;
- \(\xi_0\): bath-decoupling and independent physical-scale routes;
- invariant measure: dynamical and arithmetic routes;
- spectral dimension: direct diffusion and RG scaling routes.

Agreement is evidence. Disagreement is a scientific result.

## 25.4 Distribution-shift protocol

The theory must be tested outside the regime in which the construction was tuned:

- early universe;
- strong gravity;
- many-body sectors;
- relativistic regimes;
- quantum-information experiments.

A theory is not repaired by silently narrowing its domain after failure.

---

# 26. What Would Count as a Real Breakthrough?

### Breakthrough A — Actual dFRG

A complete microscopic coarse-graining operator and a closed flow equation are constructed.

### Breakthrough B — Actual physical measure

A unique invariant/equilibrium measure is constructed and its regularity is proved.

### Breakthrough C — Actual projection

The smooth continuum configuration space and density are derived.

### Breakthrough D — Born theorem

$$
\Pi_\#\mu_*=|\psi|^2
$$

is derived, not assumed.

### Breakthrough E — BM theorem

The Schrödinger generator and guidance law emerge in the same IR regime.

### Breakthrough F — Composition theorem

Tensor structure, entanglement and Bell correlations emerge with no signaling.

### Breakthrough G — QFT/GR recovery

The relativistic and gravitational sectors pass independent recovery and no-go audits.

A sequence of these is much more meaningful than a dozen suggestive analogies.

---

# 27. Hard Failure Conditions

The candidate must publish these before trying to argue its way around them.

1. The microscopic model secretly assumes the desired continuum theory.
2. The RG flow cannot be constructed or is mathematically ill-typed.
3. No regulator-stable fixed point exists.
4. The physical measure is nonexistent or too singular.
5. The Born identification fails.
6. The projected IR generator is not Schrödinger/BM-like.
7. The guidance law fails.
8. Composition fails to produce the quantum tensor structure.
9. No-signaling fails.
10. Lorentz or diffeomorphism recovery fails in the claimed regime.
11. Relevant-direction parameter counting does not reduce explanatory freedom.
12. Microscopic variants do not show universality without retuning.
13. The empirical T1/T2/T3 shape predictions fail.
14. E8 is needed only by analogy.
15. The arithmetic sector produces no successful unique empirical role and a cheaper theory explains all observed effects.

---

# 28. Final One-Page Architecture

## Ontology

$$
\boxed{(X,F)}
$$

One candidate microscopic deterministic substrate.

## Scale flow

$$
\boxed{\mathcal R_s: \mathcal T\to\mathcal T}
$$

One coarse-graining / RG infrastructure.

## State-to-theory bridge

$$
\boxed{\Phi:X\to\mathcal T}
$$

Constructed, not assumed.

## Measure

$$
\boxed{\mu_*\in\mathcal M(X)}
$$

Unique physical equilibrium, if it exists.

## Projection

$$
\boxed{\Pi_\#\mu_* = |\psi|^2}
$$

Born identification target.

## Quantum state

$$
\boxed{\psi=\sqrt\rho\,e^{iS/\kappa}}
$$

## Guidance

$$
\boxed{v=\nabla S/m}
$$

## Quantum potential

$$
\boxed{Q=-\frac{\kappa^2}{2m}\frac{\nabla^2\sqrt\rho}{\sqrt\rho}}
$$

## Effective composition

$$
\boxed{\mathcal H_{AB}\simeq\mathcal H_A\otimes\mathcal H_B}
$$

only after construction.

## Regimes

$$
\boxed{
\text{IST candidate ontology}
\;\to\;
\text{AS scale-flow machinery}
\;\to\;
\text{BM/QM infrared target}
}
$$

## Probability

$$
\boxed{
\text{Born identification}
\neq
\text{relaxation mechanism}
}
$$

## E8

$$
\boxed{
\text{optional information-code module, not particle census and not chirality rescue}
}
$$

## Scientific standard

$$
\boxed{
\text{truth is earned only at the gates}
}
$$

---

# 29. The Core Questions After All Revisions

The synthesis has now become smaller while the scientific questions become sharper:

1. **What is the microscopic map \(F\)?**
2. **What coarse-graining operator turns \(F\) into a controlled RG flow?**
3. **What exactly survives the flow: a fixed point, an invariant set, or a relation between them?**
4. **Why does a regular continuum configuration space emerge?**
5. **Why is the equilibrium measure Born?**
6. **Why is the Schrödinger generator linear?**
7. **Why does the same phase generate Bohmian guidance?**
8. **How does tensor-product composition emerge?**
9. **How do Bell correlations coexist with no signaling?**
10. **How does Lorentzian QFT emerge from the effective description?**
11. **How does general relativity emerge, and what is the status of gravity at the fixed point?**
12. **Which microscopic details are irrelevant, and which leave held-out fingerprints?**
13. **Is arithmetic genuinely load-bearing, or does the theory survive without it?**
14. **If E8 is selected, what selects it and what unique fingerprint follows?**

These are substantially more valuable than adding more components.

---

# 30. Closing Rule

> **Subtract until the model breaks. Then restore exactly one necessary feature.**

A feature may enter the core only when it:

1. closes a known theorem-level gap;
2. generates a genuinely new held-out prediction; or
3. removes more independent ontology than it introduces.

Everything else belongs in the workbench.

The intended candidate is therefore best represented as

$$
\boxed{
(X,F)
\xrightarrow{\;\mathcal R_s,\Phi\;}
(\Gamma_s,\mu_s)
\xrightarrow{\;\Pi\;}
(\psi,\rho,S)
\xrightarrow{\;\mathrm{IR}\;}
\mathrm{BM/QM}
\xrightarrow{\;\mathrm{recovery}\;}
\mathrm{QFT+GR+SM}
}
$$

with the uncompromising caveat:

> **This diagram is a research target, not a result.**

The program succeeds scientifically if it closes the bridges. It also succeeds scientifically if it proves that one of those bridges cannot exist.

---

# Appendix A — Secondary Workbench

These ideas remain available but are deliberately outside the core.

## A.1 Ostrowski product formula

For a nonzero rational,

$$
\prod_v |x|_v=1.
$$

This is a powerful arithmetic identity and may constrain a genuinely adelic flow. It is retained as a consistency test, not an axiom.

## A.2 Galois counterfactual semantics

Field-extension compatibility may generalize a binary rational/irrational admissibility cut into a graded structure. It enters only after the basic definability map exists.

## A.3 Cyclotomic phase lattice

A torsion phase structure may automatically enforce

$$
\oint \nabla S\cdot d\ell=2\pi n\kappa,
$$

providing a possible Wallstrom route. It is promoted only if the phase structure is derived rather than inserted.

## A.4 Arakelov/height route to \(Q\)

An archimedean component of an arithmetic height could, in principle, relate to the quantum potential. This is a fallback program, not the core derivation.

## A.5 Fractal uncertainty / Dirichlet-form route

Useful for probing projection regularity and Fisher information.

## A.6 Holographic / information-bound route

Relevant to the finite-information question, but numerical similarity to holographic bounds is not a derivation.

---

# Appendix B — Net-Valid Consolidation Summary

### DeepSeek contributions retained

- regime-structured BM/IST/AS architecture;
- asymptotic Bohmianity and asymptotic unitarity as targets;
- RDS/mixing route;
- generalized/Kramers effective dynamics;
- canonical-height definability;
- Markov-partition route to ultrametric structure;
- initial-condition debt for Valentini relaxation;
- sharpened death conditions.

### Z contributions retained

- dFRG as an explicit program;
- scheme-independence gate;
- spectral-dimension gate;
- cross-substrate universality;
- prime-universality gate;
- E8 inversion to an information/coding role;
- modular fingerprints;
- kill-test-first sequencing;
- internal Occam control.

### Gemini contributions retained

Only the high-level design lesson that effective internal geometry, if required by later gauge/QFT recovery, should be treated as **emergent** rather than introduced as a new microscopic ontology.

The bundle/local-\(k(x)\)/geodesic proposals are not admitted to the core because they add substantial structure before the existing bridge problems are solved.

---

# Appendix C — Source Discipline

This revision is deliberately **source-constrained**. The attached reviews are used to identify candidate revisions and adjudicate architectural choices. Established mathematics is not upgraded to theorem status merely because an AI review describes it confidently.

The synthesis therefore distinguishes:

- existing mathematics;
- theorem-shaped targets;
- speculative mechanisms;
- numerical targets;
- empirical claims;
- rejected claims preserved in the ledger.

That distinction is part of the theory's design, not merely editorial hygiene.

---

# Final Statement

The strongest current form of the BM + IST + Asymptotic Safety candidate is not a larger theory. It is a **smaller, more explicit, more heavily gated architecture**:

$$
\boxed{
\text{arithmetic deterministic substrate}
\;+
\text{controlled RG coarse-graining}
\;+
\text{derived infrared quantum/Bohmian regime}
}
$$

Everything else earns its place only by surviving the evaluation framework.

That is the tripartite synthesis at design optimality as presently understood: **one substrate, one flow, one projection problem, one quantum target, one recovery program, and no feature without a job.**
