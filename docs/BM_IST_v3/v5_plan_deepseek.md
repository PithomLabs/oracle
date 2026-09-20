# Implementation Plan: BM–IST Synthesis v5

## Governing principle

The v5 document and its successors have produced a comprehensive research architecture. The adversarial reviews have identified one consistent failure mode: **the program describes π rather than constructing it.** Every new document adds a layer of mathematical naming — adelic connections, Hecke co-products, Bruhat–Tits kernels, entropy-saturation inequalities — without building a single explicit π.

This implementation plan is designed to break that pattern. It prioritizes **construction over description**, **toy models over general theories**, and **kill tests over confirmations**. The plan is deliberately small in its first phase, because the first phase is where the program has repeatedly failed.

---

## Phase 0 — Foundations (Weeks 1–2)

**Goal:** Formalize the constraints that any construction must satisfy. No mathematics beyond what is already in the v5 document.

### Deliverable 0.1: G0 Theorem Ledger Entry

State L1 and L2 as precise theorems with proofs or precise references.

- **L1 (Orbit Rigidity):** A strongly continuous unitary group cannot act nontrivially on a countable discrete set with a literal action.
- **L2 (One-Parameter Triviality):** A finite signed-permutation operator group cannot support a nontrivial continuous one-parameter unitary flow.

**Status:** These are almost certainly standard results in functional analysis and representation theory. The work is to find the precise statements and cite them, not to prove them from scratch.

**Kill condition:** If either L1 or L2 turns out to be false as stated, Gate G0 collapses and the substrate can carry literal continuous dynamics. This would be a major positive result — but it is unlikely.

### Deliverable 0.2: Actual Invariant Set Specification

Read Palmer's 2009 Invariant Set Postulate paper and the 2026 paper together. Determine:

- What is the generating dynamics of I_U?
- Is I_U specified as an attractor of a specific map, or is it defined abstractly?
- What is the relationship between I_U and the bit-string representation?
- Is the bit string derived from I_U, or is it a separate structure?

**Status:** This is a read-only task. It requires no new mathematics.

**Kill condition:** If I_U is not specified precisely enough to have a well-defined generating dynamics, then the entire program is built on an underdetermined object. This would be a fatal finding — but it is better to discover it now than after years of work.

### Deliverable 0.3: Prioritization Document

One page. Which of the following is the single most important next step?

- Build explicit π for 1 qubit.
- Run the scaling audit on the actual invariant set.
- Formalize T-Prob-1.
- Build the Hecke co-product.

**Decision rule:** The first item that can be completed with a concrete deliverable in ≤ 4 weeks wins. Everything else is deferred.

---

## Phase 1 — Explicit π for 1 Qubit (Weeks 3–8)

**Goal:** Construct one explicit projection map π for a single qubit. Not the general case. Not the many-body case. One qubit.

This is the smallest possible test of the entire program. If π cannot be built for one qubit, it cannot be built for anything.

### Deliverable 1.1: Domain Specification

Define the domain of π precisely.

**Candidate domain:** The set of bit strings of length L with a fixed number m of 1s, modulo the permutation ξ.

- Number of equivalence classes: L+1 (one for each m).
- Each class corresponds to a value of cos²(θ/2) = m/L.
- The permutation ξ is fixed at system creation.

**Question:** Is this the right domain? What about states that are not in a rational basis?

**Kill condition:** If the domain cannot be specified without circular reference to the codomain, π is not well-defined.

### Deliverable 1.2: Codomain Specification

Define the codomain of π.

**Candidate codomain:** The space of qubit states in ℂ², or equivalently the Bloch sphere S².

- A bit string with m ones maps to the state cos(θ/2)|1⟩ + sin(θ/2)|−1⟩ where cos²(θ/2) = m/L.
- The permutation ξ maps to the phase φ.

**Question:** Is the codomain continuous or discrete? If continuous, how does the discrete domain map onto it?

**Kill condition:** If the codomain cannot be specified without already assuming the Born rule, π is circular.

### Deliverable 1.3: Explicit Formula for π

Write down π: {bit strings} → {qubit states} as an explicit formula.

**Candidate:**

π(b, ξ) = (√(m/L), e^{2πi n/L} √(1 − m/L))

where m is the number of 1s in b, and n is determined by ξ.

**Question:** Does this formula satisfy the requirements? Does it preserve the Born rule? Is it continuous in any limit?

**Kill condition:** If the formula is not well-defined for all bit strings, or if it fails to reproduce the Born rule for simple cases, π fails at Step 1.

### Deliverable 1.4: Verification Against Known Cases

Test π against:

- The |0⟩ state (m = L, n = 0).
- The |1⟩ state (m = 0, n = 0).
- The |+⟩ state (m = L/2, n = 0).
- The |−⟩ state (m = L/2, n = L/2).
- A generic state with m/L = 1/3, n/L = 1/7.

**Kill condition:** If π fails any of these, the construction is wrong.

### Deliverable 1.5: Continuum Limit

Check whether π has a well-defined continuum limit as L → ∞.

**Question:** Does the discrete state converge to the continuous qubit state? At what rate?

**Kill condition:** If the limit does not exist, or if it exists but does not reproduce standard quantum mechanics, the construction fails at the continuum step.

### Phase 1 success criteria

- π is explicitly written.
- π is verified against 5 simple cases.
- π has a well-defined continuum limit.
- π preserves the Born rule in the limit.

If all four are satisfied, Phase 1 succeeds and the program has, for the first time, an explicit construction. If any fails, the program must return to Phase 0 and reconsider.

---

## Phase 2 — Scaling Audit (Weeks 6–14, parallel to Phase 1)

**Goal:** Determine whether the actual BM–IST invariant set is in a regular or singular arithmetic regime.

This is the Highest-Priority Immediate Task in the v5 document. It has been deferred repeatedly. It must be done.

### Deliverable 2.1: Invariant Set Specification

From Phase 0.2, specify the actual generating dynamics of I_U. If Palmer's papers do not give a unique map, specify a candidate map and justify the choice.

**Candidate:** The skew-product on Z₂ × Z₃ with coupling λ, as proposed in Section 24 of v5.

**Question:** Is this the actual map, or a toy? If a toy, what is the relationship to the actual map?

### Deliverable 2.2: Scaling Constants

Compute the scaling constants of the actual (or candidate) invariant set.

- What are the contraction ratios?
- Are they rational, algebraic, or transcendental?
- Are they multiplicatively independent?
- Are they Pisot or non-Pisot?

### Deliverable 2.3: Absolute Continuity Check

Determine whether the invariant measure is absolutely continuous with respect to Lebesgue.

**Method:** Compute the Fourier transform of the measure and check its decay rate. If |μ̂(ξ)| decays faster than |ξ|^{-1/2}, the measure has an L² density.

**Kill condition:** If the measure is singular, Fisher information diverges and the quantum potential blows up. The program dies at Gate 2.

### Deliverable 2.4: Numerical Elephant Scan

Implement the Z₂ × Z₃ scan family computationally. Vary λ across rational, Salem, and Pisot values. Compute:

- Local dimension of the measure.
- Fourier decay rate.
- Fisher information.

**Output:** A phase diagram showing which λ values give regular measures and which give singular measures.

**Kill condition:** If all λ values give singular measures, the program dies. If only fine-tuned λ values give regular measures, the program requires fine-tuning and dies at the No-Free-Lunch test.

### Phase 2 success criteria

- The actual or candidate invariant set is specified.
- Its scaling constants are computed.
- Absolute continuity is determined (yes or no).
- If yes, Fisher information is finite.
- If no, the program is killed.

---

## Phase 3 — Gate C: Two-Qubit Composition (Weeks 10–20)

**Goal:** Extend π to two qubits and verify that it satisfies Gate C.

This is the most lethal gate. If Gate C fails, the program cannot reproduce entanglement or Bell correlations.

### Deliverable 3.1: Two-Qubit Domain

Extend the domain from Phase 1 to two bit strings (b_A, b_B) with a shared permutation ξ.

**Question:** Is the shared ξ enough to generate entanglement? Or is a coupling parameter g needed?

### Deliverable 3.2: Tensor Product

Construct the tensor product π_A ⊗ π_B and compare it to π_AB.

**Target:** π_AB should factorize as π_A ⊗ π_B in the limit g → 0, with an error bound f(g) → 0.

**Kill condition:** If π_AB does not factorize, or if the error bound does not go to zero, Gate C fails.

### Deliverable 3.3: Entanglement Generation

Show that for g ≠ 0, π_AB generates entangled states.

**Method:** Compute the entanglement entropy of π_AB(b_A, b_B, ξ) for various g. Show that it is nonzero for g ≠ 0.

### Deliverable 3.4: No-Signaling

Verify that π_AB does not allow signaling.

**Method:** Show that the marginal distribution on b_A is independent of the choice of measurement on b_B, and vice versa.

**Kill condition:** If signaling is possible, the construction violates special relativity and is dead.

### Deliverable 3.5: Bell Correlation

For a Bell test with settings A, B, compute the correlation E(A, B) = ⟨π_AB | σ_A ⊗ σ_B | π_AB⟩.

**Target:** The correlation should violate Bell's inequality for suitable settings.

**Kill condition:** If the correlation satisfies Bell's inequality, the construction does not reproduce quantum mechanics.

### Phase 3 success criteria

- Two-qubit π is explicitly constructed.
- Tensor product factorization holds in the g → 0 limit.
- Entanglement is generated for g ≠ 0.
- No-signaling holds.
- Bell inequality is violated.

If all five are satisfied, Gate C is cleared and the program has a genuine many-body construction.

---

## Phase 4 — Confinement and Wood–Spekkens (Weeks 18–28)

**Goal:** Derive the distribution of exact measurement settings from the substrate, and show that it clears the Wood–Spekkens gate without fine-tuning.

### Deliverable 4.1: Setting Distribution

Derive P(θ_A, θ_B | λ) from the substrate.

**Method:** Use the Hecke algebra structure (if it exists) or an alternative algebraic mechanism to constrain the allowed setting pairs.

**Kill condition:** If P(θ_A, θ_B | λ) requires per-experiment adjustment, the construction is fine-tuned and fails.

### Deliverable 4.2: Wood–Spekkens Measure

Compute the Wood–Spekkens fine-tuning measure for the derived distribution.

**Target:** The measure should be zero or bounded, indicating no fine-tuning.

**Kill condition:** If the measure diverges or scales with system size, the construction is fine-tuned and fails.

### Deliverable 4.3: Bell-Exact Correlations

Show that the derived setting distribution yields Bell-exact correlations matching quantum mechanics.

**Target:** The correlations should match the singlet-state prediction E(A, B) = −A·B exactly.

**Kill condition:** If the correlations deviate from quantum mechanics, the construction fails.

### Phase 4 success criteria

- Setting distribution is derived from the substrate.
- Wood–Spekkens fine-tuning measure is zero or bounded.
- Bell-exact correlations are reproduced.

---

## Phase 5 — Quantum Potential and Guidance (Weeks 24–36)

**Goal:** Derive the Bohmian quantum potential and guidance equation from the substrate.

This is the bridge from IST's kinematics to BM's dynamics. It is where the Metric Chasm must be crossed.

### Deliverable 5.1: Density Derivation

Derive ρ(q) from the substrate measure μ.

**Method:** Push forward μ under π.

**Target:** ρ(q) is absolutely continuous with finite Fisher information.

### Deliverable 5.2: Phase Derivation

Derive S(q) from the substrate phase structure.

**Method:** Identify the phase with the holonomy of a connection on the substrate, or with the argument of the effective amplitude ψ = √ρ e^{iS/κ}.

### Deliverable 5.3: Quantum Potential

Compute Q = −(κ²/2m)(∇²√ρ)/√ρ from the derived ρ.

**Target:** Q should match the standard quantum potential.

**Kill condition:** If Q does not match, or if Q diverges, the construction fails.

### Deliverable 5.4: Guidance Equation

Derive v = ∇S/m from the derived S.

**Target:** The guidance equation should match BM exactly.

**Kill condition:** If the guidance equation does not match, or if it is not well-defined, the construction fails.

### Phase 5 success criteria

- ρ(q) derived from substrate measure.
- S(q) derived from substrate phase.
- Q matches standard quantum potential.
- Guidance equation matches BM.

---

## Phase 6 — Empirical Handles (Weeks 30–40)

**Goal:** Connect the construction to at least one falsifiable prediction.

### Deliverable 6.1: Finite-L Signatures

Derive the leading-order deviation from standard quantum mechanics due to finite L.

**Target:** A specific prediction for quantum-computing experiments.

### Deliverable 6.2: Parameter Cross-Check

Fix parameters (c, ξ₀, κ) using one experiment. Predict the outcome of a second experiment without refitting.

**Kill condition:** If the parameters cannot be fixed by one experiment and used to predict another, the construction is not falsifiable.

### Deliverable 6.3: Null Results

If the construction predicts no deviation in a given regime, state it explicitly. Null predictions are first-class outputs.

### Phase 6 success criteria

- At least one falsifiable prediction is derived.
- Parameters are cross-checked across at least two experiments.
- Null results are recorded.

---

## Cross-cutting disciplines

### D1: Provenance ledger

Every claim in the implementation plan gets a ledger entry with:

- Claim ID
- Exact statement
- Source
- Status
- Dependencies
- Derivation location
- Kill condition
- Verification date
- Reviewer

### D2: No descriptions masquerading as constructions

If a deliverable says "the mechanism is X" without providing an explicit formula, it is not a construction. It is a description. It must be flagged as such.

### D3: Toy models first

Every general construction should first be tested on a toy model (1 qubit, L = 6 or 8, Z₂ × Z₃). If the toy model fails, the general construction is unlikely to succeed.

### D4: Kill tests are mandatory

Every deliverable must have an explicit kill condition. If the kill condition cannot be stated, the deliverable is not well-posed.

### D5: Null results are first-class outputs

A clean no-go is as valuable as a positive construction. The plan should be designed to produce decisive results, not to produce confirmations.

---

## What this plan deliberately excludes

- **Adelic fiber bundles.** Deferred until Phase 1 succeeds. If π cannot be built for one qubit, adelic connections are premature.
- **Hecke algebra co-products.** Deferred until Phase 3. The Hecke structure may be useful for Gate C, but not before.
- **Bruhat–Tits trees.** Deferred until Phase 5. The tree structure may be useful for the bulk-boundary kernel, but not before.
- **Furstenberg/Hochman rigidity.** Deferred until Phase 2 reveals the actual scaling regime. If the actual invariant set is regular, these tools are unnecessary. If it is singular, they may not apply anyway.
- **Arakelov height conjecture.** Deferred indefinitely. It is the most speculative import and depends on everything else.

---

## Risk assessment

**Risk 1: π cannot be built for one qubit.** Probability: moderate. Mitigation: Phase 1 is designed to discover this quickly.

**Risk 2: The actual invariant set is singular.** Probability: high. Mitigation: Phase 2 is designed to discover this decisively. If singular, the program dies at Gate 2 and this is a clean no-go.

**Risk 3: Gate C fails.** Probability: moderate. Mitigation: Phase 3 is the most consequential test. If it fails, the program cannot reproduce entanglement.

**Risk 4: Wood–Spekkens cannot be cleared.** Probability: high. Mitigation: Phase 4 is designed to be honest about this. If it fails, the program may survive as a nonlocal hidden-variable theory but not as a local one.

**Risk 5: The program produces only descriptions, not constructions.** Probability: moderate. Mitigation: the discipline of D2 and the emphasis on toy models.

---

## Timeline and milestones

| Phase | Weeks | Deliverable | Kill condition |
|---|---|---|---|
| 0 | 1–2 | G0 ledger, invariant set specification, prioritization | If I_U is underdetermined, program dies |
| 1 | 3–8 | Explicit π for 1 qubit | If π fails simple cases, program dies |
| 2 | 6–14 | Scaling audit | If measure is singular, program dies |
| 3 | 10–20 | Two-qubit π and Gate C | If Gate C fails, no entanglement |
| 4 | 18–28 | Wood–Spekkens confinement | If fine-tuning is required, program dies |
| 5 | 24–36 | Quantum potential and guidance | If Q or v do not match, program dies |
| 6 | 30–40 | Empirical handles | If no falsifiable prediction, program is not science |

**Total:** 40 weeks (≈ 10 months) to a decisive result.

---

## Success criteria for the program

**Outcome A (Constructive survival):** All phases succeed. The program has an explicit π, a two-qubit construction, a derived setting distribution, a derived quantum potential and guidance equation, and at least one falsifiable prediction.

**Outcome B (Partial survival):** Some phases succeed, some fail. The program has a partial construction with clearly identified gaps.

**Outcome C (Clean no-go):** A kill condition is triggered. The program has a decisive negative result. This is scientifically valuable.

**Outcome D (Stalled):** No phase produces a decisive result. The program has neither a construction nor a no-go. This is the failure mode the plan is designed to prevent.

---

## The one-sentence summary

**Build π for one qubit first. If it works, extend to two. If it doesn't, stop.**
