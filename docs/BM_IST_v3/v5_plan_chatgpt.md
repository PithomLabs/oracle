# BM–IST Synthesis v5 — Implementation Plan

The implementation should be treated as a **research program with hard gates**, not as an attempt to write the whole theory at once.

The governing principle is:

$$
\boxed{
\text{formalize}
\rightarrow
\text{construct}
\rightarrow
\text{prove}
\rightarrow
\text{compute}
\rightarrow
\text{test}
\rightarrow
\text{kill or advance}
}
$$

The v5 architecture already identifies the central ordering: **G0 → arithmetic substrate → scaling/measure → probability → projection → composition → amplitude/phase → guidance/contextuality → gravity**. 

My implementation plan is to turn that architecture into a sequence of **research packages, machine-checkable artifacts, explicit gates, and stop conditions**.

---

# 1. First principle: do not implement the whole theory

The biggest implementation mistake would be to start coding a giant “BM–IST simulator” before the mathematical dependencies are settled.

The v5 document contains several things that are currently only theorem targets or speculative shapes: the projection \(\pi\), T-Prob-1, T-Born, the tensor-product calculus, the Wood–Spekkens confinement mechanism, and the various proposed forms of \(\pi\). 

So the implementation should be organized as:

```text
Research Core
│
├── Formal definitions
├── Theorem ledger
├── Counterexample / no-go engine
├── Arithmetic dynamics engine
├── Measure / scaling analysis
├── Projection experiments
├── Many-body composition tests
├── Elephant-scan computation
└── Empirical-handle pipeline
```

Not:

```text
"Build BM-IST"
    ↓
one large simulation
```

---

# 2. Establish the BM–IST formal core

Before any substantive computation, freeze a minimal mathematical vocabulary.

Create a formal specification containing:

$$
\mathcal I
=
\text{substrate state space}
$$

$$
F:\mathcal I\rightarrow\mathcal I
=
\text{discrete dynamics}
$$

$$
\mu
=
\text{candidate invariant measure}
$$

$$
\pi
=
\text{projection/reconstruction structure}
$$

$$
\psi
=
\sqrt{\rho}\,e^{iS/\kappa}
$$

and:

$$
Q
=
-\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt{\rho}}{\sqrt{\rho}}.
$$

The key is to define each object independently.

Do not allow:

```text
π := "whatever produces ψ"
```

or:

```text
μ := distribution that gives Born
```

Those are circular definitions.

### Deliverable

`BMIST-FORMAL-CORE-v0.1.md`

Containing:

* objects;
* domains/codomains;
* topology/metric;
* discrete time;
* measure;
* equivalence relation;
* parameters;
* desired effective observables;
* explicit assumptions.

---

# 3. Build the theorem / claim provenance ledger

The v5 document explicitly calls for a provenance/verification ledger distinguishing `[EXISTING]`, `[THEOREM-SHAPED]`, `[SPECULATIVE-SHAPE]`, `[ASSERTED, NOT DERIVED]`, and `[SELF-CORRECTED]`. 

Implement that as an actual artifact rather than a documentation convention.

For every claim record:

```text
claim_id
statement
mathematical_dependencies
source
status
proof_status
counterexample_status
computational_test
experimental_handle
failure_condition
last_verified
```

This matters because v5 explicitly treats provenance as a competitive asset and warns against upgrading repeated speculation into established fact. 

### Deliverable

`claims.yaml`

plus:

`CLAIMS.md`

generated from it.

This should become the authoritative status layer for the research program.

---

# 4. Implement Gate G0 first

This is the first hard gate because it determines the entire architecture.

The consolidation identifies:

* **L1 Orbit Rigidity**
* **L2 One-Parameter Triviality**

and concludes that a countable/discrete defined-state ontology cannot simultaneously carry literal nontrivial continuous unitary evolution. 

Do not merely cite this.

Create a formal proof notebook.

## G0 implementation

### G0.1

Formalize the substrate as a discrete/countable object.

### G0.2

Formalize the continuous unitary group:

$$
U(t)=e^{-iHt/\kappa}.
$$

### G0.3

State exactly what “acts on the substrate” means.

### G0.4

Produce the impossibility proof.

### G0.5

Create a positive control example:

* continuous state space;
* continuous unitary flow;
* no contradiction.

### G0.6

Create a discrete-time successful example:

$$
x_{n+1}=F(x_n)
$$

with an effective continuum observable.

### G0 output

One of:

```text
PASS: substrate must be discrete-time / emergent continuum
```

or

```text
FAIL: proposed ontology still requires primitive continuous quantum flow
```

Nothing downstream should proceed while G0 is unresolved.

---

# 5. Build the arithmetic substrate

Once G0 is fixed, construct the smallest concrete deterministic substrate.

I would start with:

$$
X = \prod_{p\in P} X_p
$$

rather than immediately attempting a full physical invariant set.

For the initial computational program, use a finite product such as:

$$
X=\mathbb Z_2\times\mathbb Z_3
$$

because this is already the proposed basis of the elephant scan. 

Define:

$$
F_\lambda:X\rightarrow X.
$$

The parameters should include:

* coupling \(\lambda\);
* depth/resolution \(L\);
* any explicitly justified arithmetic parameters.

The model must be deterministic.

No stochastic noise gets added just to make the probability distribution look right.

---

# 6. Introduce canonical-height structure

The next layer is intrinsic discreteness.

The v5 architecture proposes canonical heights because they can provide:

* intrinsic arithmetic structure;
* discrete bounded-height sets;
* dynamically selected preperiodic loci;
* native discrete-time dynamics. 

The implementation should therefore define:

$$
F:X\rightarrow X
$$

and a height:

$$
\hat h:X\rightarrow \mathbb R_{\ge0}
$$

with explicit questions:

1. Does \(\hat h\) exist?
2. Is it invariant or monotonic under \(F\)?
3. What is the zero-height locus?
4. Is the zero-height/preperiodic set the proposed physical skeleton?
5. Does it contain sufficient dynamical richness?
6. Is it countable where required?

### Important

Do not assert:

```text
height zero = physical state
```

until the relationship is derived.

### Deliverable

`ARITHMETIC-SUBSTRATE-v0.1`

with proofs, computations, and counterexamples.

---

# 7. Run the Scaling Audit before building the Born mechanism

This should be the first major computational research package.

The v5 document specifically identifies the Bernoulli-convolution/Pisot issue and says the actionable task is to audit **BM–IST's own scaling constants**, not Palmer's. 

The audit should compute:

$$
\lambda_i
$$

for the actual proposed invariant-set construction and determine:

* arithmetic class;
* overlap structure;
* absolute continuity/singularity;
* Fourier decay;
* dimension;
* density regularity;
* Fisher information.

## Output

For every scaling regime:

```text
parameter region
measure class
absolute continuity?
L² density?
Fourier decay?
Fisher finite?
```

### Critical decision

If the invariant measure is generically singular in the physical regime, stop the current probability branch and investigate whether the substrate itself must change.

Do **not** compensate by smoothing the output numerically.

---

# 8. Formalize T-Prob-1

The v5 target is approximately:

$$
\text{multiplicative independence}
+
\text{non-Pisot}
+
\text{transversality}
\Rightarrow
\rho\in L^2
$$

with enough Fourier control to make Fisher information finite and permit convergence toward Born statistics. 

This should become a theorem program.

## Implementation sequence

### T1

State exact hypotheses.

### T2

Build toy systems satisfying them.

### T3

Prove the easy implications first.

### T4

Numerically test representative cases.

### T5

Identify where existing theorems do not directly compose.

### T6

Prove or disprove the remaining bridge.

The target is not “show numerically that the histogram looks Gaussian/Born-like.”

It is:

$$
\boxed{
\text{prove the analytic regularity needed by the quantum layer}
}
$$

---

# 9. Build a probability laboratory

Create a standalone library for comparing candidate measure mechanisms.

It should implement at least:

### Mechanism A

Natural/SRB invariant measure.

### Mechanism B

p-adic Haar measure.

### Mechanism C

Entropy-saturation selection.

### Mechanism D

Small-height/arithmetic equidistribution.

### Mechanism E

Palmer finite-\(L\) microcanonical combinatorics.

These mechanisms must never be silently merged.

The v5 synthesis explicitly separates Palmer's finite-\(L\) combinatorics from the much harder continuum bridge. 

### Deliverable

`BORN-MEASURE-LAB`

with the output:

```text
candidate measure
    ↓
substrate distribution
    ↓
projection
    ↓
configuration density
    ↓
comparison against |ψ|²
```

---

# 10. Attack T-Born

After the measure laboratory exists, attack the strongest theorem-shaped Born route.

The target is essentially:

$$
\text{arithmetic dynamics}
\rightarrow
\text{small-height / near-defined states}
\rightarrow
\text{equidistribution}
\rightarrow
\rho_{\rm arch}=|\psi|^2.
$$

The v5 document labels T-Born as one of its deepest new theorem targets. 

The implementation should produce three independent objects:

1. substrate equilibrium measure;
2. effective configuration projection;
3. wavefunction whose squared modulus is compared with that projection.

### Hard rule

Do not define \(\psi\) from \(\rho\) by:

$$
\psi=\sqrt{\rho}
$$

and then claim Born has been derived.

The phase problem remains.

---

# 11. Build the phase laboratory

The program needs an explicit phase object.

Candidate construction:

$$
\mu_L
=
\{e^{2\pi i k/L}\}.
$$

Test whether the finite/cyclotomic structure can generate a continuous effective phase field.

The Wallstrom requirement becomes:

$$
\oint \nabla S\cdot dl
=
2\pi n\kappa.
$$

Implementation steps:

1. define discrete phase states;
2. define phase transitions;
3. define loop transport;
4. calculate holonomy;
5. derive circulation quantization;
6. test the continuum limit.

If quantization has to be imposed manually, record that as a failure of this branch.

---

# 12. Solve the Projection Problem as a separate package

This should be the central engineering artifact of the entire research program.

Do not build one \(\pi\) and assume it is correct.

Build **three competing implementations**.

---

## 12.1 Projection A — \(\mathcal E\circ\mathcal R\circ\mathcal H\)

Implement:

$$
\pi_A=\mathcal E\circ\mathcal R\circ\mathcal H.
$$

### H

Map substrate dynamics to discrete causal/hyperbolic structure.

### R

Apply transfer/RG/coarse-graining operators.

### E

Decode effective amplitude.

The v5 consolidation already lays out the six-stage target for this branch. 

---

## 12.2 Projection B — Adelic connection

Implement a mathematical prototype for:

$$
\pi_B = \nabla_{\rm adelic}
$$

with:

* base;
* fibers;
* connection;
* holonomy;
* curvature;
* geodesics.

The first objective is not to prove the physical interpretation.

It is simply to determine whether the structure is mathematically well-defined.

---

## 12.3 Projection C — Sheaf / elephant scan

Implement:

$$
\pi_C
=
\text{local data}
+
\text{gluing/obstruction structure}.
$$

Then run the explicit scan family.

The v5 consolidation considers this the most operationalized proposal. 

---

# 13. Establish a projection competition protocol

All three projection branches should be subjected to the same tests.

| Test                 | Requirement                                    |
| -------------------- | ---------------------------------------------- |
| Determinism          | same substrate state → same effective object   |
| Measure preservation | required invariant measure survives projection |
| Locality             | effective local structure where required       |
| Regularity           | enough smoothness for differential operators   |
| Composition          | compatible with tensor products                |
| Phase                | produces meaningful phase                      |
| Linearity            | effective superposition                        |
| Dynamics             | effective Schrödinger evolution                |
| Guidance             | Bohmian velocity                               |
| Empirical            | produces a distinct prediction                 |

The important implementation decision is:

> **Do not merge branches merely because they use similar vocabulary.**

They should compete.

---

# 14. Build the amplitude realization layer

Once the projection candidates exist, attack the Independent Realization Problem:

$$
A[I_U,\Phi_t,\mu,\pi]\rightarrow\psi.
$$

The v5 specification calls for S1–S10, including:

* deterministic realization;
* Born modulus;
* phase transport;
* exclusion compliance;
* local second-order Schrödinger generator;
* locality;
* linearity;
* Wallstrom quantization;
* universal \(\kappa\);
* Galilean covariance;
* equivariance.

Implementation-wise, create a test harness:

```text
AmplitudeCandidate
    ├── S1
    ├── S2
    ├── ...
    └── S10
```

Every candidate receives a machine-readable result:

```text
PASS
FAIL
UNKNOWN
BLOCKED
```

with evidence.

---

# 15. Gate C comes before serious Bell work

The v5 consolidation explicitly identifies the emergent tensor-product calculus as a missing first-class obligation. 

Implement this before building elaborate Bell simulations.

Start with two weakly coupled systems:

$$
A,\;B.
$$

Require:

$$
\psi_{AB}
\rightarrow
\psi_A\otimes\psi_B
$$

as coupling

$$
g\rightarrow0.
$$

Measure the factorization error:

$$
\epsilon(g)
=
\|\psi_{AB}-\psi_A\otimes\psi_B\|.
$$

The desired result is:

$$
\epsilon(g)\rightarrow0.
$$

Then introduce interaction and test whether entanglement genuinely appears.

This creates a concrete numerical theorem target instead of an interpretive claim.

---

# 16. Derive no-signaling

Once Gate C exists, construct a bipartite measurement model.

Require:

$$
P(a|x,y)
=
P(a|x)
$$

under the effective theory.

Then test Bell correlations.

The order matters:

```text
factorization
→ interaction
→ entanglement
→ Bell correlations
→ no-signaling
```

not:

```text
make Bell violation appear
```

---

# 17. Implement contextuality as a structural test

Create a contextuality layer using local measurement contexts.

Represent:

```text
context
measurement
outcome assignment
compatibility relation
```

and test for:

$$
\text{global section existence/nonexistence}.
$$

The goal is to determine whether contextuality emerges from the projection rather than being inserted into the measurement rules.

This is where the sheaf proposal becomes an actual mathematical implementation rather than metaphor.

---

# 18. Attack the Wood–Spekkens problem explicitly

The v5 document is very clear that this remains a real open obligation: measurement dependence must be explained, not renamed. 

Implement a model of:

$$
P(\lambda|x,y)
$$

or the appropriate exact-setting analogue.

Then ask:

1. Can the distribution be derived from the substrate?
2. Is it stable?
3. Does it reproduce Born statistics?
4. Is the dependence compressed into substrate structure rather than hand-tuned?
5. Does the mechanism survive perturbation?
6. Does it reproduce multiple experiments with one parameter set?

If not, the confinement branch fails.

---

# 19. Derive the quantum potential only after \(\rho\) and \(S\) exist independently

The exact target is:

$$
Q
=
-\frac{\kappa^2}{2m}
\frac{\nabla^2\sqrt{\rho}}{\sqrt{\rho}}.
$$

This should be implemented as a symbolic identity test.

Inputs:

$$
\rho,\;S,\;\kappa.
$$

Outputs:

$$
Q,\;v.
$$

The test succeeds only if the derivation follows from the substrate/projection construction.

The earlier “forbidden gaps → pressure gradient” story is explicitly marked as asserted rather than derived. 

So do not implement “pressure” as a placeholder for \(Q\).

---

# 20. Derive guidance

Once \(S\) exists:

$$
v=
\frac{\nabla S}{m}.
$$

Then test:

$$
\partial_t\rho+\nabla\cdot(\rho v)=0.
$$

The critical test is **equivariance**:

if the substrate starts in the physical measure, does the effective distribution remain:

$$
\rho_t=|\psi_t|^2?
$$

That closes the Born/guidance loop.

---

# 21. Universal-\(\kappa\) two-mass experiment

The v5 program specifically separates universal \(\kappa\) from mass-dependent \(L(m,E)\).

This should become an explicit computational/experimental test:

$$
m_1\neq m_2
$$

while fitting:

$$
\kappa_1,\kappa_2.
$$

Require:

$$
\kappa_1=\kappa_2
$$

within the theory's predicted uncertainty.

At the same time allow:

$$
L_1\neq L_2.
$$

This is the cleanest implementation of the proposed mass-entry distinction.

---

# 22. Build the Elephant Scan as the main computational microscope

The elephant scan should not be treated as a visualization project.

It should be a **research instrument**.

For every:

$$
(\lambda,L)
$$

produce:

```text
entropy
Lyapunov data where defined
invariant measure
Hausdorff/information dimension where computable
Fourier decay
density regularity
Fisher information
spectral data
phase winding
factorization error
effective generator error
Born error
```

Then classify each point:

```text
CHAOTIC
ORDERED
CRITICAL
SINGULAR
FISHER-INVALID
AMPLITUDE-INVALID
COMPOSITION-INVALID
PROMISING
```

The important output may be that there is **no viable region**.

The v5 architecture explicitly treats an empty critical surface as a decisive result. 

---

# 23. Add an automated “theory killer”

This is important enough to make explicit.

Every experimental candidate should have automatic rejection rules.

For example:

```text
if measure_singular:
    FAIL(PROBABILITY)

if fisher_inf == infinity:
    FAIL(FISHER)

if kappa_mass_dependent:
    FAIL(GATE_4)

if factorization_error does not vanish:
    FAIL(GATE_C)

if born_error does not converge:
    FAIL(T_BORN)

if generator_not_second_order_local:
    FAIL(IRP)

if no_global_phase:
    FAIL(PHASE)

if per-handle_parameter_fitting:
    FAIL(EMPIRICAL_DISCIPLINE)
```

This keeps the research from drifting into confirmation bias.

---

# 24. Separate symbolic, numerical, and empirical layers

The implementation should have three clearly separate systems.

## Symbolic layer

For:

* definitions;
* algebra;
* differential identities;
* theorem statements;
* formal derivations.

## Numerical layer

For:

* invariant sets;
* transfer operators;
* finite-depth scans;
* density estimates;
* Fourier transforms;
* Fisher information;
* scaling experiments.

## Empirical layer

For:

* experimental constraints;
* cosmological data;
* interferometry;
* Bell bounds;
* PBHs;
* CMB;
* BBN.

Never let a numerical approximation be treated as a theorem, or an empirical fit as a derivation.

---

# 25. Empirical handle pipeline

Only after the quantum core survives should the program activate the cosmological layer.

The v5 document proposes common handles including mesoscopic interferometry, Bell/Born deviations, PBHs, radial-acceleration relations, cluster lensing, CMB, BBN, and GW propagation. 

Implement a parameter registry:

```text
c
ξ0
κ
t_ent
L(m,E)
a0
```

Every empirical model consumes the same registered parameters.

No experiment gets:

```text
best-fit parameters for this dataset
```

unless that parameter freedom is itself predicted.

---

# 26. Lock P1 vs P2 before cosmology

The p-adic sector must be classified as either:

### P1 — trace-relative

or:

### P2 — autonomous.

This is already identified as an architectural decision in v5. 

Do not begin interpreting PBH or cluster observations until this is frozen.

Otherwise the model can silently change its ontological interpretation depending on the dataset.

---

# 27. Cross-handle audit

When the theoretical model exists, run:

$$
(c,\xi_0,\kappa)
$$

through the entire empirical matrix.

The criterion is not:

> Can every observation be fit?

It is:

> **Can one fixed theory survive every independently constraining handle?**

This is one of the strongest methodological commitments of the v5 architecture. 

---

# 28. Recommended repository structure

I would implement the research program roughly as:

```text
bm-ist/
├── docs/
│   ├── formal-core.md
│   ├── architecture.md
│   ├── gates.md
│   ├── provenance.md
│   └── empirical-handles.md
│
├── claims/
│   ├── claims.yaml
│   └── verified/
│
├── math/
│   ├── g0/
│   ├── arithmetic/
│   ├── scaling/
│   ├── measure/
│   ├── projection/
│   ├── amplitude/
│   ├── phase/
│   ├── composition/
│   ├── contextuality/
│   └── guidance/
│
├── simulations/
│   ├── substrate/
│   ├── transfer/
│   ├── scaling/
│   ├── born/
│   ├── phase/
│   ├── gate_c/
│   └── elephant_scan/
│
├── data/
│   ├── generated/
│   ├── external/
│   └── results/
│
├── experiments/
│   ├── g0/
│   ├── irp/
│   ├── gate_c/
│   ├── born/
│   └── empirical/
│
└── reports/
    ├── scaling-audit/
    ├── t-prob-1/
    ├── t-born/
    ├── projection/
    ├── elephant-scan/
    └── cross-handle/
```

The important architectural principle is that **failed branches remain recorded**, rather than deleted.

---

# 29. Milestone sequence

## Milestone M0 — Formalization

Deliver:

* formal core;
* provenance ledger;
* parameter registry;
* gate definitions.

**Exit condition:** all primary objects have precise domains and codomains.

---

## Milestone M1 — G0

Deliver:

* formal proof;
* discrete-time substrate;
* continuous-limit strategy.

**Exit condition:** no primitive continuous-unitary contradiction.

---

## Milestone M2 — Arithmetic Substrate

Deliver:

* explicit \(F\);
* canonical-height candidate;
* invariant-set construction.

**Exit condition:** genuine deterministic substrate exists.

---

## Milestone M3 — Scaling Audit

Deliver:

* arithmetic regime;
* measure classification;
* Fisher status.

**Kill condition:** required regularity impossible in the candidate substrate.

---

## Milestone M4 — Probability

Deliver:

* invariant measure;
* T-Prob-1 result;
* T-Born attack;
* Palmer finite-\(L\) comparison.

**Kill condition:** no defensible route to \(\rho=|\psi|^2\).

---

## Milestone M5 — Projection

Deliver:

* at least one explicit \(\pi\);
* competition between A/B/C branches.

**Kill condition:** no stable projection object can be defined.

---

## Milestone M6 — Quantum Core

Deliver:

* \(\psi\);
* \(S\);
* \(\rho\);
* \(\kappa\);
* local generator.

**Kill condition:** IRP failure.

---

## Milestone M7 — Gate C

Deliver:

* tensor-product approximation;
* interaction;
* entanglement;
* no-signaling.

**Kill condition:** multi-system composition cannot emerge.

---

## Milestone M8 — Guidance / Contextuality

Deliver:

* \(Q\);
* Bohm velocity;
* equivariance;
* contextuality;
* Wood–Spekkens analysis.

**Kill condition:** guidance or confinement requires unexplained external structure.

---

## Milestone M9 — Empirical Layer

Deliver:

* fixed parameter set;
* positive prediction;
* cross-handle analysis.

**Kill condition:** per-dataset tuning or inconsistent predictions.

---

# 30. What I would implement first

The first actual execution sequence should be only this:

```text
1. Freeze BM–IST Formal Core
2. Create Claims/Provenance Ledger
3. Formalize G0
4. Implement minimal Z2 × Z3 discrete substrate
5. Implement λ and finite-depth L
6. Compute invariant measures
7. Run Scaling/Pisot/Fourier/Fisher Audit
8. Build first Elephant Scan
9. Attempt T-Prob-1 on the actual scan family
10. Only then implement candidate π
```

This ordering is deliberate.

The first major question is **not**:

> “Can we derive the Schrödinger equation?”

It is:

> **“Does the proposed deterministic substrate even possess the analytic and measure-theoretic structure required for a quantum projection?”**

If the answer is no, the rest of the program should be redesigned before more machinery is added.

---

# 31. Research decision tree

The whole program can then be reduced to this:

```text
                  ┌─────────────┐
                  │   G0 PASS?  │
                  └──────┬──────┘
                         │
                       NO → STOP
                         │
                        YES
                         ↓
                ┌─────────────────┐
                │ Arithmetic      │
                │ substrate exists│
                └────────┬────────┘
                         │
                       NO → STOP
                         │
                        YES
                         ↓
                ┌─────────────────┐
                │ Scaling /       │
                │ measure PASS?   │
                └────────┬────────┘
                         │
                       NO → REVISE
                         │
                        YES
                         ↓
                ┌─────────────────┐
                │ T-Prob / T-Born │
                └────────┬────────┘
                         │
                       FAIL → REVISE/KILL
                         │
                       PASS
                         ↓
                ┌─────────────────┐
                │ Explicit π      │
                └────────┬────────┘
                         │
                       FAIL → KILL BRANCH
                         │
                       PASS
                         ↓
                ┌─────────────────┐
                │ IRP S1–S10     │
                └────────┬────────┘
                         │
                       FAIL → NO-GO
                         │
                       PASS
                         ↓
                ┌─────────────────┐
                │ Gate C          │
                └────────┬────────┘
                         │
                       FAIL → NO-GO
                         │
                       PASS
                         ↓
                ┌─────────────────┐
                │ Guidance /      │
                │ contextuality   │
                └────────┬────────┘
                         │
                       FAIL → NO-GO
                         │
                       PASS
                         ↓
                ┌─────────────────┐
                │ Fixed-parameter │
                │ empirical tests │
                └────────┬────────┘
                         │
                 ┌───────┴────────┐
               FAIL               PASS
                 │                  │
              NO-GO            SURVIVING
                             BM–IST THEORY
```

---

# 32. The central implementation philosophy

The v5 material's most important methodological insight is that several attractive mechanisms were proposed and later self-corrected because they were analogies masquerading as derivations. The document explicitly identifies the pattern “propose, formalize, audit, retract or narrow, and only then treat the survivor as usable.” 

So I would make that the actual execution protocol:

$$
\boxed{
\textbf{Every idea enters as a hypothesis.}
}
$$

Then:

$$
\boxed{
\textbf{Every hypothesis receives a test.}
}
$$

Then:

$$
\boxed{
\textbf{Every test has a kill condition.}
}
$$

And finally:

$$
\boxed{
\textbf{A result is promoted only when its evidentiary status changes.}
}
$$

That gives BM–IST v5 something much more useful than another conceptual synthesis: **a research machine capable of producing either a constructive theory or a rigorous sequence of no-go results.**

