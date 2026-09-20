# BM–IST Synthesis v5
## Consolidated Research Architecture for Bohmian Mechanics, Invariant-Set Theory, Arithmetic Dynamics, and the Projection Problem

**Version:** 5.1  
**Purpose:** Standing research reference and executable research program  
**Status:** Research blueprint; not a finished physical theory

### Changelog: v5.0 → v5.1

- §3.1/§3.2 — G0's supporting lemmas (L1, L2) now carry inline proof sketches; G0 was previously stated as a conclusion without its argument reproduced in this document.
- §36 — Added the direct corollary that Palmer's unmodified-Schrödinger assumption is not merely a weaker starting point but, under G0, a provably illegitimate move for any countable-beable program.
- §13/§14–17 — Added an explicit note on how the four probability routes relate to one another, and flagged the dropped six-level unification as an intentional omission, not an oversight.

---

## 0. How to Read This Document

This v5 consolidates the earlier BM–IST synthesis architecture with the latest cross-model consolidation. It preserves self-corrections, separates theorem-backed mathematics from conjecture and analogy, and turns the strongest ideas into explicit gates and work packages.

### Status tags

- `[EXISTING]` — published mathematics/formalism being imported without claiming it proves BM–IST.
- `[THEOREM-SHAPED]` — precise statement that can be proved or refuted but is not established here.
- `[SPECULATIVE-SHAPE]` — candidate architecture or conjecture not yet derived.
- `[ASSERTED, NOT DERIVED]` — previously proposed mechanism explicitly recognized as unsupported.
- `[SELF-CORRECTED]` — earlier claim subsequently narrowed or corrected.

### Governing method

> **Proposal → formalization → adversarial audit → theorem/no-go → survivor specification → experiment.**

A decisive null result is a successful research result.

---

# 1. Executive Summary

The BM–IST program is fundamentally a **Projection Problem**, not a catalog of desired properties.

The central object is some structured map of the form

\[
\pi:\mathcal I_{IST}\longrightarrow \mathcal B_{BM}
\]

that turns deterministic, intrinsically discrete invariant-set data into the effective Bohmian/quantum objects

\[
(Q,\rho,S,\psi).
\]

The effective layer must recover, rather than assume:

- Born statistics;
- Bohmian guidance;
- a universal action scale \(\kappa\);
- local second-order Schrödinger dynamics;
- exact superposition;
- tensor-product composition;
- entanglement and no-signaling;
- contextuality structure;
- the correct quantum potential;
- and eventually gravitational/cosmological phenomenology.

The strongest new constraint is **Gate G0**:

> A countable/discrete defined-state set cannot itself carry a nontrivial literal strongly-continuous unitary flow.

Therefore the continuum is not allowed to be primitive on the countable substrate. Quantum time, amplitude, and continuous dynamics must be **emergent**.

The current research architecture is:

\[
\boxed{
\text{Arithmetic substrate}
\to
\text{discrete dynamics}
\to
\text{invariant measure}
\to
\text{scaling regularity}
\to
\pi
\to
(\rho,S,\psi)
\to
\text{Born + guidance + linearity}
\to
\text{Gate C}
\to
\text{contextuality/Bell}
\to
\text{gravity}
}
\]

---

# 2. Core Diagnosis

## 2.1 Palmer's program is a program, not a finished derivation

The consolidated material treats Palmer's 2026 RaQM work as a valuable source of candidate ontology, finite-information mechanisms, and sharp pressure points, but not as a completed proof. The distinction between valid mathematics, proposed mechanisms, and deferred derivations must be preserved in BM–IST itself.

## 2.2 BM–IST's real task

The synthesis does not need more prose describing what \(\pi\) ought to do. It needs an explicit mathematical object: domain, codomain, equivalence relation, regularity, composition properties, and theorems or failure conditions.

---

# 3. Gate G0 — Discrete Substrate / Continuous Unitary Compatibility

## 3.1 L1 — Orbit Rigidity `[EXISTING / THEOREM-LEVEL]`

A strongly continuous unitary group cannot act nontrivially on a countable discrete/totally disconnected defined-state set when the action is interpreted literally on that set. The continuous orbit collapses.

**Proof sketch.** For fixed \(v\), the orbit map \(t\mapsto U_tv\) is continuous by strong continuity of \(t\mapsto U_t\). The parameter group (\(\mathbb R\), or any connected Lie group) is connected, so its continuous image is connected. But a countable, totally disconnected set has no connected subset larger than a single point. Hence the image \(\{U_tv:t\in\mathbb R\}\) is a single point: \(U_tv=v\) for all \(t\). The orbit is forced constant, not merely slow-moving.

## 3.2 L2 — One-Parameter Triviality `[EXISTING / THEOREM-LEVEL]`

A finite signed-permutation operator group cannot itself support a nontrivial continuous one-parameter unitary flow.

**Proof sketch.** A continuous one-parameter flow is, by definition, a continuous homomorphism \(\phi:\mathbb R\to G\) into the operator group \(G\). \(\mathbb R\) is connected, so \(\phi(\mathbb R)\) is a connected subgroup of \(G\). \(G\) is finite, hence discrete and totally disconnected, so its only connected subgroup is the trivial one, \(\{e\}\). Therefore \(\phi(t)=e\) for all \(t\): the flow is constant. This is the same connectedness argument as L1, applied to the operator group itself rather than to a single orbit — the two lemmas fail the construction from two different angles for the same underlying reason.

## 3.3 Gate G0

\[
\boxed{
\text{countable substrate}
+
\text{literal continuous unitary evolution}
\Rightarrow
\text{structural incompatibility}
}
\]

**Required consequence:** discrete substrate time must be primitive; continuous Schrödinger evolution must be emergent.

**Failure condition:** any construction that directly places literal continuous Schrödinger dynamics on the countable substrate fails G0.

---

# 4. Ontological Mapping

The synthesis must distinguish:

- substrate state \(\xi\);
- configuration \(Q(t)\);
- density \(\rho(q,t)\);
- phase/action field \(S(q,t)\);
- effective amplitude \(\psi=\sqrt\rho e^{iS/\kappa}\).

## 4.1 \(\xi\) as trajectory-generating seed `[ASSERTED, NOT DERIVED]`

The “flipbook” interpretation—\(\xi\) generates the entire continuous trajectory—was explicitly self-criticized because a global hidden variable is not automatically type-compatible with a per-system configuration, and Palmer's collapse process is not Bohmian guidance.

Use it only as the question:

> What substrate object actually plays the role of \(Q_0\) and determines the trajectory?

## 4.2 \(\xi\) as physical global phase `[SPECULATIVE-SHAPE]`

A distinct proposal interprets \(\xi\) as the physical ontology behind quantum phase. This is a genuine candidate not previously present in the BM–IST design, and it connects naturally to a cyclotomic phase lattice, but no derivation yet establishes it.

## 4.3 Palmer's finite-information amplitude

The finite relation

\[
\cos^2(\theta/2)=m/L
\]

is a legitimate discrete combinatorial probability mechanism. It does **not** by itself provide a continuum amplitude field, dynamics, composition law, or guidance law. Keep finite-\(L\) combinatorics separate from the continuum bridge.

---

# 5. Parameters \(L\) and \(\kappa\)

## 5.1 Mass-entry mechanism

Palmer's candidate

\[
L=L(m,E)
\]

is relevant because BM guidance contains

\[
v=\frac{\nabla S}{m}.
\]

This gives BM–IST a concrete candidate route for mass to enter the effective theory. It remains a hypothesis, not a derivation.

## 5.2 Universal \(\kappa\) versus capacity \(L\)

Keep distinct:

\[
\boxed{\kappa=\text{universal action/phase scale}}
\]

and

\[
\boxed{L=L(m,E)=\text{system-dependent information/granularity capacity}.}
\]

A varying \(L\) does not require a mass-dependent \(\kappa\).

## 5.3 Holographic grounding of \(L\) `[SPECULATIVE-SHAPE]`

A possible grounding is a horizon-area/Planck-area bound such as

\[
L\sim \frac{A}{4\ell_P^2},
\]

but no physical identification should be made from a numerical resemblance alone.

---

# 6. Arithmetic / Number-Theoretic Grounding

The central lesson is to replace coordinate-stipulated rational/irrational cuts with intrinsic arithmetic structure.

## 6.1 p-adic / Diophantine resolution `[SPECULATIVE-SHAPE]`

Candidate principle: physical resolvability is finite resolution in an ultrametric/adelic structure rather than a binary rationality rule.

The operational requirement is an observable consequence distinguishing state-space p-adic structure from a purely archimedean description.

## 6.2 Scaling audit and the Pisot trap

An earlier claim incorrectly treated \(\lambda=1/2\) as a singular Bernoulli-convolution corner. The corrected result is:

- \(\lambda=1/2\) is the favorable dyadic/Lebesgue case;
- singular/non-decaying-Fourier danger instead arises at non-integer reciprocal-Pisot parameters such as \(1/\varphi\).

### Actionable result

Run the arithmetic/scaling audit on **BM–IST's own invariant set**, not Palmer's. This can decide whether the projected probability marginal is regular enough for finite Fisher information.

---

# 7. T-Prob-1 — Probability/Fisher Theorem Target

### `[THEOREM-SHAPED]`

Candidate theorem:

> If the invariant-set scaling constants are multiplicatively independent, non-Pisot, and satisfy an appropriate overlap-transversality condition, then the natural configuration marginal has an \(L^2\) density, hence finite Fisher information, with quantitative Fourier decay controlling convergence toward the effective Born distribution.

This is not established.

The mathematical toolkit may include Bernoulli-convolution regularity, transversality, local-entropy methods, and Fourier decay.

**Why it matters:** failure here could kill the program before any quantum-potential derivation.

---

# 8. Canonical Heights and Dynamically Selected Discreteness

Canonical heights are a candidate intrinsic replacement for coordinate-defined rationality.

Potential properties:

- coordinate-independent arithmetic definition;
- Northcott-type discreteness at bounded degree/height;
- dynamically selected zero-height/preperiodic structure under suitable maps;
- naturally discrete-time dynamics.

Candidate structure:

\[
x_{n+1}=F(x_n),\qquad \hat h(F(x))\text{ structured/monotone},
\]

with a low-height or preperiodic skeleton serving as the candidate physically defined set.

This is an imported mathematical mechanism, not itself a physical theory.

---

# 9. Galois Counterfactual Semantics

A richer alternative to “defined/undefined” is graded field compatibility:

\[
x\in K_1,\qquad y\in K_2.
\]

Counterfactual availability becomes a question of field extension and algebraic compatibility rather than a coordinate-specific rational cut.

Status: `[SPECULATIVE-SHAPE]`.

Potential empirical consequence: field-structured or log-periodic deviations rather than a flat list of excluded angles.

---

# 10. Cyclotomic Phase Lattice / Wallstrom

Let

\[
\mu_L=\{e^{2\pi i k/L}:k=0,\ldots,L-1\}.
\]

If the effective phase is a projection of a torsion structure, then winding quantization may become automatic.

Candidate target:

\[
\oint \nabla S\cdot d\ell =2\pi n\kappa.
\]

Status: `[LEMMA-SHAPED]`.

Required proof: construct the phase bundle, define the projection, control nodal regions, and prove the allowed holonomy.

---

# 11. Dynamical Bridge: Discrete to Continuous

The required bridge is

\[
\text{discrete deterministic map}
\rightarrow
\text{continuous conservative effective dynamics}.
\]

## 11.1 Decoherence-as-redistribution `[ASSERTED, NOT DERIVED]`

Environmental degrees of freedom can redistribute information while the total system remains deterministic, but this does not by itself produce definite Bohmian trajectories or a guidance law.

## 11.2 p-adic fractional transport `[SPECULATIVE-SHAPE]`

Candidate:

\[
\partial_t\rho=-D^\alpha\rho.
\]

Possible tools include Vladimirov-type operators and coarse-graining. The actual drift and guidance equation must be derived rather than inserted.

## 11.3 Arithmetic dynamics as native discrete time `[EXISTING mathematics imported as structural fix]`

Map-based arithmetic dynamics directly respects G0 by putting discrete time at the substrate level and leaving continuous time to emergence.

---

# 12. Renormalization / Transfer-Operator Bridge

Define a coarse-graining operator

\[
R:\mathcal I\to\mathcal I
\]

and seek an infrared equivalence class or critical object:

\[
R^n(x)\to [x]_{IR}.
\]

Then \(\pi\) can be an effective map from microscopic states to the IR class rather than a pointwise map.

Separately, let \(\mathcal L_\sigma\) be the transfer operator of the substrate map \(\sigma\), with invariant measure

\[
\mathcal L_\sigma\mu=\mu.
\]

A candidate Born structure is

\[
\rho(q)=\pi_\#\mu.
\]

Status for the full identification: `[SPECULATIVE-SHAPE]`.

---

# 13. Probability and the Born Rule

Separate three problems:

### P1 — Measure existence
Does the substrate have a canonical invariant/natural measure?

### P2 — Projection regularity
Does its pushforward possess a suitable density and finite information functionals?

### P3 — Born identification
Can one prove

\[
\rho=|\psi|^2?
\]

Do not collapse P1–P3 into a single “typicality” argument.

### How §14–§17 relate to P1–P3

The four routes below are **not** a required stack and **not** yet shown to be equivalent. Each targets a different subset of P1–P3, and more than one may turn out to fail:

| Route | Targets | Status relative to P1–P3 |
|---|---|---|
| §14 Haar-slaving | P1 | Candidate source of canonical measure only |
| §15 Entropy-saturation selection | P1, partially P3 | Candidate selection principle for *which* measure is physical |
| §16 Embedded-agent bootstrap | P3 | Candidate self-consistency argument for why the selected measure is the *observed* one |
| §17 T-Born (arithmetic equidistribution) | P2, P3 | The only route with an explicit convergence mechanism (equidistribution rate) |

An earlier version of this material proposed treating all four as facets of a single underlying measure, unified across "six levels." That unification claim is **deliberately dropped here**, not merely omitted by oversight: it was never more than a suggestive framing, and folding four independently-speculative routes into one narrative risked hiding exactly the kind of unproven collapse §13 warns against. Until at least one route clears P1–P3 on its own, treat them as competing candidates to be tested and discarded independently.

---

# 14. Haar-Slaving `[SPECULATIVE-SHAPE]`

p-adic spaces carry canonical Haar measure. A candidate strategy is to let the p-adic bath supply the normalization structure and transfer it to the archimedean effective layer through a precisely defined influence functional.

The target is an actual mapping

\[
\mu_p\to\rho_{\mathbb R}
\]

not a verbal claim that the p-adic sector “grounds probability.”

---

# 15. Entropy-Saturation Selection `[THEOREM-SHAPED]`

A candidate alternative to unexplained typicality is that the physical measure is the unique measure selected by an entropy equality/saturation condition in the relevant dynamical system.

The program must define:

- the dynamical system;
- the entropy inequality;
- the equality case;
- uniqueness;
- the map from the selected measure to \(|\psi|^2\).

---

# 16. Embedded-Agent Bootstrap `[THEOREM-SHAPED SCAFFOLD]`

Treat observers as internal subsystems and define a self-consistency/fixed-point condition between physical measure and observer memory/self-location.

Candidate form:

\[
\mu_{phys}=\mathcal F(\mu_{phys}).
\]

A solution matters only if existence and uniqueness can be proved and the fixed point coincides with the physical measure selected by the dynamics.

---

# 17. T-Born — Arithmetic Equidistribution

### `[THEOREM-SHAPED]`

Candidate statement:

> Quantum equilibrium is the archimedean equilibrium measure of the substrate's arithmetic dynamics, approached by equidistribution of small-height or near-defined preparation states.

Desired chain:

\[
\text{arithmetic orbit}
\to
\text{small-height states}
\to
\text{equidistribution}
\to
\rho_{arch}=|\psi|^2.
\]

Potential inputs include Bilu-type and Brolin–Lyubich-type equidistribution results. This is a target, not a result.

---

# 18. Palmer's Microcanonical Mechanism

Uniform permutations of a fixed-composition bit string can give exact finite ratios such as \(m/L\).

This is a useful discrete probability construction.

It does **not** solve the continuum projection problem.

The correct division of labor is:

- finite-\(L\) combinatorics → discrete probability mechanism;
- arithmetic regularity → continuum measure;
- projection → effective amplitude/phase;
- quantum generator → dynamics.

---

# 19. Fractal Uncertainty Principle Route

Existing fractal-uncertainty results constrain simultaneous localization on fractal position and dual supports.

Candidate role:

\[
\text{fractal substrate}
\to
\text{uncertainty/delocalization bound}.
\]

Potential value: a geometry-first source of quantum-like spreading.

It is **not** by itself a derivation of the Bohmian quantum potential.

---

# 20. Quantum Potential and Nonlocality

The exact target is

\[
Q=-\frac{\kappa^2}{2m}\frac{\nabla^2\sqrt\rho}{\sqrt\rho}.
\]

The slogan “forbidden gaps create pressure” is explicitly `[ASSERTED, NOT DERIVED]` and cannot stand in for this equation.

Keep the objects distinct:

\[
\psi=\sqrt\rho e^{iS/\kappa},
\]

so that:

- \(\rho\) determines \(Q\);
- \(S\) determines guidance;
- both derive from the same effective amplitude.

## 20.1 Holographic bulk–boundary kernel `[SPECULATIVE-SHAPE]`

Candidate:

\[
\psi(x)=\int K(x_p,x)\,d\mu(x_p).
\]

Possible structures include Bruhat–Tits trees and p-adic bulk/boundary kernels. A valid derivation must define \(K\), prove regularity and normalization, and recover the effective generator.

## 20.2 Arakelov quantum-potential conjecture `[SPECULATIVE-SHAPE]`

Candidate conjecture: the quantum potential is the archimedean component of an arithmetic height functional. This is the deepest and most speculative import; it should be attacked only after measure and projection exist.

---

# 21. The Projection Problem Proper

The current best abstraction is not simply

\[
\pi:\xi\mapsto q.
\]

Instead, investigate a structured object such as

\[
\boxed{
\pi:(\mathcal I,\Phi_t,\mu,\mathcal P,\mathcal H)
\to
(\rho,S,\psi,Q)
}
\]

where the input may include invariant data, discrete dynamics, invariant measure, partition/cyclotomic structure, and effective-state data.

Three independent proposal families converge on a negative design constraint:

> **A viable \(\pi\) is unlikely to be a simple pointwise function.**

It may act on equivalence classes, measures, transfer operators, gluing data, connections, or critical surfaces.

---

# 22. Three Candidate Shapes of \(\pi\)

## Shape A — Composite operator

\[
\pi=\mathcal E\circ\mathcal R\circ\mathcal H.
\]

- \(\mathcal H\): holographic/causal embedding of substrate dynamics;
- \(\mathcal R\): transfer/RG coarse-graining;
- \(\mathcal E\): effective decoder.

Required program:

1. define the lattice/causal substrate;
2. define the shift dynamics;
3. define the transfer operator;
4. define RG;
5. prove the effective fixed point;
6. derive Schrödinger dynamics;
7. test any QEC/decoder interpretation separately.

## Shape B — Adelic connection `[SPECULATIVE-SHAPE]`

Base: an appropriate adelic product. Fibers: p-adic substrate and archimedean/BM effective space. \(\pi\) is a connection.

Candidate identifications:

- curvature ↔ Born/probability structure;
- geodesics ↔ guidance.

Target theorems:

1. connection existence;
2. curvature/Born correspondence;
3. geodesic/guidance correspondence.

Possible complete failure: the fiber-bundle construction is a category error for the actual physical state space.

## Shape C — Sheaf gluing + elephant scan

Use local mathematical descriptions as charts and treat contextuality as a possible obstruction to global gluing.

Two useful methodological licenses:

- **Takens embedding theorem `[EXISTING]`** — generic partial observations can reconstruct attractor structure;
- **Feigenbaum universality `[EXISTING]`** — a model family can stand in for an exact microscopic map if universality-class membership is established.

These license model testing, not arbitrary theory-gluing.

---

# 23. Mixing–Unitarity Dilemma

A strongly mixing substrate gives natural statistics but may destroy coherent amplitude information. A highly ordered odometer-like substrate preserves spectral organization but may lack the mixing required for robust statistics.

Candidate resolution:

> the effective quantum structure may lie on a **critical surface** between chaotic/mixing and ordered regimes.

This is a hypothesis, not a theorem.

---

# 24. Elephant-Scan Protocol

Use an explicit family, e.g. a skew product on \(\mathbb Z_2\times\mathbb Z_3\), with coupling \(\lambda\) and finite depth \(L\).

### S1 — Define the map

\[
F_\lambda:X\to X.
\]

### S2 — Compute dynamical invariants

Entropy, spectral data, recurrence, dimension, Lyapunov-type quantities where meaningful.

### S3 — Compute the invariant measure

Determine whether a natural/SRB/ergodic measure exists and is unique in the relevant class.

### S4 — Define and test the projection

Construct candidate factor/observable maps.

### S5 — Test probability

Measure density regularity, Fisher information, Fourier decay, and limiting Born behavior.

### S6 — Test phase and composition

Test phase winding, weak-coupling factorization, and the beginning of the Gate C program.

### S7 — Test effective quantum dynamics

Look for a local second-order generator, effective linearity, conservation/unitarity, and Bohmian guidance.

Outputs should be numbers, invariants, theorem statements, counterexamples, or empty regions—not explanatory prose alone.

---

# 25. Independent Realization Problem (IRP)

The surviving amplitude target is

\[
\boxed{
A[I_U,\Phi_t,\mu,\pi]\to\psi
}
\]

without relying on already-excluded fixed-scale factor, Reynolds-stress, diffusive-homogenization, zero-noise invariant-drift, conditional-expectation, or incompatible exact linear-factor routes.

## S1 — Determinism

\(\psi\) must be fixed by substrate data.

## S2 — Born modulus and phase transport

\[
|\psi|^2=\rho,
\qquad
S=\kappa\arg\psi.
\]

## S3 — No-go compliance

Explicitly avoid already-killed constructions.

## S4 — Effective Schrödinger dynamics

\[
i\kappa\partial_t\psi=\hat H\psi
\]

with the correct local second-order structure.

## S5 — Effective locality

The substrate may be ultrametric; the effective generator must have the required local form.

## S6 — Exact linearity/superposition

\[
\psi=a\psi_1+b\psi_2.
\]

## S7 — Wallstrom/circulation quantization

## S8 — Universal \(\kappa\)

At least a two-mass test.

## S9 — Galilean covariance

## S10 — Equivariance/closure

The substrate-generated distribution must remain compatible with the effective evolution.

---

# 26. Gate 2 — Fisher Information

The Fisher structure may be considered only after the projected density exists:

\[
I_F[\rho]=\int\frac{|\nabla\rho|^2}{\rho}\,dq.
\]

Required dependency:

\[
\text{measure}
\to
\text{density}
\to
\text{finite Fisher information}
\to
Q.
\]

Never assume Fisher information because it produces the desired quantum-potential formula.

---

# 27. Gate 3 — Effective Linearity

The microscopic dynamics may be deterministic and nonlinear while the effective amplitude layer is linear. That emergence requires explanation.

Candidate mechanisms:

- nonlinear lift;
- RG fixed-point linearization;
- transfer-operator eigenspaces;
- sheaf/gluing structure;
- hidden linear representation;
- decoder/QEC route.

A viable mechanism must pass exact superposition, not merely reproduce a few solutions.

---

# 28. Gate 4 — Universal Action Scale

Require

\[
\kappa_{effective}=\text{same universal constant}
\]

for systems with different masses and energies, while allowing finite information capacity

\[
L=L(m,E).
\]

A mass-dependent \(\kappa\) fails the current program unless the entire effective quantum target is correspondingly revised.

---

# 29. Gate C — Emergent Tensor-Product Calculus

This is a first-class new obligation.

For weak coupling \(g\to0\):

\[
\psi_{AB}\approx\psi_A\otimes\psi_B,
\]

with an error bound

\[
\|\psi_{AB}-\psi_A\otimes\psi_B\|\le f(g),
\qquad f(g)\to0.
\]

For interaction \(g\neq0\):

- entanglement must emerge;
- probability must remain normalized;
- no-signaling must hold;
- Bell correlations must be reproduced.

Without Gate C, a genuine Bell calculation is impossible inside the framework.

---

# 30. Contextuality / Sheaf Structure

Treat measurement contexts as local charts and quantum contextuality as an obstruction to a global section where the mathematical conditions warrant it.

Candidate goal:

\[
\text{local IST assignments}
\xrightarrow{\pi}
\text{quantum contextual structure}.
\]

The sheaf formalism is a language for a real compatibility problem, not a substitute for constructing \(\pi\).

---

# 31. Locality, Bell, and the Wood–Spekkens Gate

## 31.1 Metric chasm

p-adic state-space distance and Euclidean physical-space distance are different objects. Physical noise cannot be declared p-adically stable without a proven cross-space map.

## 31.2 Arithmetic rigidity / ×2×3 route

Furstenberg-style rigidity is a possible confinement mechanism, but the full ×2×3 problem is open. Use only proven special cases as proven results.

## 31.3 Galois counterfactuals

Field-structured counterfactual availability is a candidate generalization of narrow rational-angle exclusion.

## 31.4 Wood–Spekkens target

The real requirement is to explain the **distribution** of exact measurement settings, not merely rename measurement dependence.

Target:

\[
\text{measurement dependence}
\to
\text{derived exact-setting distribution}
\to
\text{Born-exact Bell correlations}
\]

without per-experiment hidden tuning.

---

# 32. TOE Import Matrix

External theories enter only as constraints or tools.

| Candidate | Imported feature | BM–IST role | Failure condition |
|---|---|---|---|
| String / AdS-CFT | bulk-boundary structure | candidate projection/decoder | independent descriptions disagree |
| LQG | discrete spectra | substrate/operator constraints | incompatible effective spectrum |
| CDT / Wolfram | rewrite rules | generative discrete substrate | primitive closed-form continuum required |
| Causal sets | order + counting | discrete probability | non-Born marginal |
| 't Hooft CA | information loss/emergence | deterministic substrate | loss not Born-compatible |
| Jacobson / Verlinde | entanglement/thermodynamic gravity | scale-fixing candidate | wrong mass/scale dependence |
| ER=EPR | connectivity/correlation | geometry of correlations | contextuality unaffected |
| QEC / HaPPY | logical/physical encoding | decoder interpretation | correction threshold failure |
| Everett/decoherence | pointer selection | basis emergence | no preferred basis |
| Noncommutative geometry | spectral triples | algebraic reconstruction | no physical identification |
| Topos theory | sheaves/geometric morphisms | gluing/contextuality | no useful global structure |

**Rule:** borrowed mathematics earns a place only if it creates a falsifiable constraint.

---

# 33. P1 vs P2 — The p-adic Sector

### P1 — Trace-relative

The p-adic sector affects ordinary observables only through coupling to the visible sector.

### P2 — Autonomous

The p-adic sector possesses independent effective content/stress-energy and evolves as an autonomous sector.

This choice controls whether the cosmological interpretation is closer to a modified-gravity/trace-relative picture or a true dark-sector picture. It must be settled before fitting astrophysical data.

---

# 34. Candidate Positive / Presence Predictions

At least one non-null, nontrivial prediction is required.

## 34.1 Gravity-only shadow sector `[SPECULATIVE]`

The p-adic bath could carry gravitationally active but archimedean-invisible content.

**Not derived:** any particular dark-to-baryonic numerical ratio.

## 34.2 Acceleration scale `[TARGET]`

A possible target is

\[
a_0\sim\frac{cH_0}{2\pi}.
\]

This is a candidate target, not a result.

## 34.3 Arithmetic Born deviations `[SPECULATIVE]`

Potential signatures near special arithmetic corners may be field-structured or log-periodic.

## 34.4 Finite-L quantum-computing signatures `[IMPORTED HYPOTHESIS]`

A mass/energy-dependent granularity law may create a mesoscopic test window.

## 34.5 PBH abundance/mass function `[CANDIDATE HANDLE]`

Can test whether the p-adic sector has autonomous cosmological content, provided the same parameter set is used elsewhere.

---

# 35. Empirical Handle Table

The governing discipline is one fixed parameter set, e.g.

\[
(c,\xi_0,\kappa),
\]

across all observations unless scale dependence is derived.

| Handle | Discriminates | Parameter | Status |
|---|---|---|---|
| Mesoscopic interferometry ceiling | core finite-granularity mechanism | \(c,\xi_0\) | pending |
| Bell/Born-deviation bounds | projection/confinement | \(c\) | pending |
| PBH mass function | autonomous p-adic content | \(\xi_0\) | research |
| Radial-acceleration scatter | trace-relative vs autonomous | \(a_0\) mechanism | research |
| Cluster baryons/lensing | modified gravity vs dark sector | tracing rule | research |
| CMB peaks/heights | decoupling/partition | \(t_{ent}\) | research |
| BBN expansion rate | early coupling | \(t_{ent}\) | research |
| GW propagation | modified-propagation branches | — | constrained |

**No per-row fitting. Ever.**

---

# 36. Comparative Audit: BM–IST and Palmer

This is intentionally not a winner-take-all comparison.

## BM–IST's current strengths

- explicit Projection Problem;
- continuous quantum dynamics treated as a target rather than a primitive assumption;
- explicit Bohmian trajectory/guidance target;
- explicit fine-tuning and confinement obligations;
- broader many-body obligations;
- **a direct G0 corollary against Palmer's own construction.** Palmer's paper states that Schrödinger evolution is not modified — i.e., it takes literal continuous unitary dynamics as an axiom sitting alongside a countable bit-string beable. Under G0 (§3), this is not merely a weaker starting assumption than BM–IST's; it is the exact configuration L1/L2 rule out. Any countable-beable program that also assumes unmodified continuous Schrödinger evolution is, by G0, internally inconsistent — Palmer's construction included, unless continuous evolution is demoted from axiom to emergent target on his side as well. This is a corollary of G0, not a rhetorical comparison.

## Palmer's current contributions

- concrete finite-information ontology;
- explicit bit-string/permutation structure;
- finite \(m/L\) probability mechanism;
- candidate \(L(m,E)\) mass-entry mechanism;
- candidate physical phase ontology;
- proven narrow counterfactual-exclusion tool via Niven's theorem.

## Productive relation

\[
\boxed{
\text{Palmer candidate ontology}
\to
\text{BM–IST exclusion machinery}
\to
\text{derive or kill}
}
\]

The most informative experiment may be whether the finite-information ontology can survive BM–IST's full realization tests.

---

# 37. Cross-Cutting Cautions

1. **Analogy is not derivation.** Flipbooks, pressure-gradient stories, p-adic stability stories, and “gaps produce Q” are not proofs.
2. **State-space and physical-space metrics are not interchangeable.**
3. **Density is not typicality.**
4. **Numerical coincidence is not evidence.**
5. **Describing \(\pi\) is not building \(\pi\).**
6. **External mathematics is not a physical identification.**
7. **Self-corrections must remain visible.**
8. **Credit and critique should be separated.**

---

# 38. Provenance / Verification Ledger

Every load-bearing claim gets a ledger entry:

| Field | Required |
|---|---|
| Claim ID | stable identifier |
| Claim | exact statement |
| Source | primary theorem/paper/derivation |
| Status | epistemic tag |
| Dependencies | prior results |
| Derivation location | equation/notebook/page |
| Counterexample status | known/unknown/tested |
| Empirical handle | if applicable |
| Failure condition | exact kill test |
| Last verification | date |
| Reviewer | human/checker |
| Notes | corrections/self-corrections |

Provenance is part of the scientific method, not administration.

---

# 39. Work Packages

## WP0 — Formal core

Define substrate, state variables, discrete dynamics, spaces, measures, parameters, and epistemic status.

## WP1 — G0 / discrete-time substrate

Formalize the discrete-time primitive and remove literal continuous-unitary dynamics from the substrate.

## WP2 — Scaling audit

Classify the actual BM–IST invariant-set scaling; determine regularity, Fourier behavior, and Fisher finiteness.

## WP3 — Arithmetic dynamics / height

Construct a candidate arithmetic map and canonical-height skeleton.

## WP4 — Probability

Run T-Prob-1, Haar-slaving, entropy selection, embedded-agent bootstrap, T-Born, and finite-\(L\) comparison.

## WP5 — Projection

Build and compare Shapes A, B, and C.

## WP6 — Amplitude / phase

Attack IRP S1–S10.

## WP7 — Gate C

Construct emergent tensor-product composition, interaction, entanglement, and no-signaling.

## WP8 — Contextuality / Wood–Spekkens

Formalize local-global obstruction and exact-setting distribution.

## WP9 — Quantum potential / guidance

Derive \(Q\) and Bohmian guidance from the same effective amplitude/phase object.

## WP10 — Gravity / cosmology

Only after the quantum core survives; resolve P1 vs P2, \(a_0\), CMB, BBN, PBH, and propagation constraints.

## WP11 — Elephant scan

Implement the explicit \((\lambda,L)\) scan family and preserve null regions as first-class output.

---

# 40. Recommended Dependency Order

\[
\boxed{
G0
\to
\text{Arithmetic substrate}
\to
\text{Scaling audit}
\to
\text{Measure}
\to
\text{T-Born / probability}
\to
\pi
\to
\text{Gate C}
\to
\text{Amplitude}
\to
\text{Phase}
\to
\text{Guidance}
\to
\text{Contextuality / Wood–Spekkens}
\to
\text{Gravity}
}
\]

The elephant scan is the computational engine that can operate alongside the analytic program once its map and observables are fixed.

---

# 41. Highest-Priority Immediate Tasks

### 1. Scaling Audit

Determine whether BM–IST's actual invariant set is in a regular/transversal regime or a singular arithmetic regime.

### 2. T-Prob-1

Convert the scaling result into a rigorous probability/Fisher theorem attempt.

### 3. Canonical-height substrate

Build the intrinsic discrete dynamical model.

### 4. Explicit \(\pi\)

Construct at least one candidate projection object.

### 5. Gate C

Do not proceed to serious Bell phenomenology without composition.

### 6. T-Born

Attempt the arithmetic-equidistribution route to Born equilibrium.

### 7. Wood–Spekkens gate

Attack measurement-dependence quantitatively.

### 8. Quantum potential

Attempt the exact \(Q\) derivation only after \(\rho\) and \(S\) exist.

### 9. Gravity

Delay dark-sector/cosmology fitting until the quantum core survives.

---

# 42. No-Free-Lunch Tests

Every mechanism must answer:

1. **Where did the information come from?**
2. **Where did the phase come from?**
3. **Where did \(\kappa\) come from?**
4. **Where did linearity come from?**
5. **Where did tensor products come from?**
6. **Where did Born normalization come from?**
7. **Where did Bell nonlocality/correlation structure come from?**
8. **Where did the exact quantum potential come from?**
9. **Where does mass enter?**
10. **What experiment can kill the mechanism?**

If an answer is “it is assumed,” that layer is not derived.

---

# 43. Success Criteria

## Outcome A — Constructive survival

An explicit substrate and projection survive G0, scaling, measure, T-Prob-1, T-Born, IRP, Gate C, contextuality, Wood–Spekkens, and the empirical cross-checks.

## Outcome B — Partial survival

A mathematically rigorous subset survives and isolates the remaining unresolved layer.

## Outcome C — Clean no-go

A theorem shows that no object satisfying an explicit set of constraints exists.

All three are scientifically useful outcomes.

---

# 44. Evidence Ladder

\[
\text{definition}
\to
\text{lemma}
\to
\text{construction}
\to
\text{theorem}
\to
\text{numerical validation}
\to
\text{experimental prediction}
\to
\text{independent confirmation}.
\]

The synthesis must not skip levels.

---

# 45. Final Consolidated Thesis

The strongest form of BM–IST is not:

> “The invariant set might explain quantum mechanics because it is discrete, fractal, and p-adic.”

It is:

> **Construct a deterministic, intrinsically discrete dynamical substrate whose invariant arithmetic structure yields a mathematically regular measure and an explicit projection into an effective complex amplitude/phase field; prove that this projection yields Born statistics, a universal action scale, linear second-order quantum dynamics, Bohmian guidance, tensor-product composition, contextuality, and Bell correlations; then determine whether the same fixed parameters produce gravitational and cosmological phenomenology.**

The core equations/constraints are:

\[
\boxed{
\text{countable substrate}
\neq
\text{literal continuous unitary flow}
}
\]

\[
\boxed{\text{density}\neq\text{typicality}}
\]

\[
\boxed{\text{finite-}L\text{ combinatorics}\neq\text{continuum Born derivation}}
\]

\[
\boxed{\text{borrowed mathematics}\neq\text{physical derivation}}
\]

and the central open object remains:

\[
\boxed{
\pi:
\text{IST invariant data}
\rightarrow
\text{effective BM amplitude/phase/configuration structure}
}
\]

The most consequential missing many-body obligation is:

\[
\boxed{\text{Gate C: emergent tensor-product calculus}}
\]

The empirical discipline is:

\[
\boxed{(c,\xi_0,\kappa)\text{ survives every handle without per-row fitting}.}
\]

The governing rule is:

> **Build the projection. Do not merely describe what it should look like.**

---

# 46. Immediate Execution Checklist

## Mathematical

- [ ] Formalize G0 in the theorem ledger.
- [ ] Specify BM–IST's actual invariant-set scaling.
- [ ] Run the Pisot/transversality/regularity audit.
- [ ] Determine absolute continuity or singularity.
- [ ] Test finite Fisher information.
- [ ] Define the arithmetic dynamical map.
- [ ] Construct a candidate canonical height.
- [ ] Define and characterize the invariant measure.
- [ ] Attack T-Prob-1.
- [ ] Attack T-Born.

## Projection

- [ ] Formalize Shape A.
- [ ] Formalize Shape B.
- [ ] Formalize Shape C.
- [ ] Determine whether the shapes are compatible, complementary, or mutually exclusive.
- [ ] Produce at least one explicit \(\pi\).

## Quantum

- [ ] Solve IRP S1–S10.
- [ ] Derive \(\rho\).
- [ ] Derive \(S\).
- [ ] Derive \(\psi\).
- [ ] Derive effective linearity.
- [ ] Derive universal \(\kappa\).
- [ ] Derive Wallstrom quantization.
- [ ] Derive \(Q\).
- [ ] Derive guidance.
- [ ] Establish equivariance.

## Many-body

- [ ] Gate C factorization.
- [ ] Interaction/entanglement.
- [ ] No-signaling.
- [ ] Bell correlations.
- [ ] Contextuality compatibility.

## Confinement

- [ ] Exact-setting model.
- [ ] Wood–Spekkens analysis.
- [ ] Arithmetic-rigidity analysis.
- [ ] Remove unsupported p-adic stability claims.

## Computation

- [ ] Define the \(\mathbb Z_2\times\mathbb Z_3\) scan family.
- [ ] Define \(\lambda\).
- [ ] Define finite depth \(L\).
- [ ] Implement S1–S7.
- [ ] Preserve failed regions and null results.
- [ ] Make outputs reproducible.

## Empirical

- [ ] Fix common parameters.
- [ ] Derive finite-\(L\) signatures.
- [ ] Derive candidate Born deviations.
- [ ] Resolve P1 vs P2 before cosmological fitting.
- [ ] Run the full parameter cross-check.

---

# 47. Closing Principle

The strongest synthesis is not the one with the most theories attached to it. It is the one in which every attractive imported idea is forced through the same sequence:

\[
\boxed{
\text{Can it be defined?}
\rightarrow
\text{Can it be constructed?}
\rightarrow
\text{Can it be proved?}
\rightarrow
\text{Can it be measured?}
\rightarrow
\text{Can it be killed?}
}
\]

A future result may therefore end in either:

\[
\boxed{\text{SURVIVES}}
\]

or

\[
\boxed{\text{NO-GO}}
\]

with equal scientific status.
