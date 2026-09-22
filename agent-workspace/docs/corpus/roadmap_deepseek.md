```markdown
# BM–IST–AS Tripartite: Concrete Attack Plan and Roadmap
## Prioritized Deliverables and Gates — High Level

**Purpose:** this is the operational roadmap. No prose. No new theory. Every item is a gate, a deliverable, a null model, or a stop condition. Prioritized by **information gain per unit effort** — the highest-priority item is the one that most cheaply kills or confirms the largest branch of the program.

---

## 0. Priority Order at a Glance

```
P0  Foundation      G0, G1           — blocking; nothing proceeds without these
P1  Amplitude       G2, G3           — hardest; unique to this program
P2  Measure         G4               — shared wall; where every rival stalled
P3  Projection      G5, G6           — shared wall; second where every rival stalled
P4  Quantum         G7, G8, G9       — recovery of BM and composition
P5  Recovery        G10              — long-horizon; QFT + GR + SM
X   Cross-cutting   nulls, spectral dimension, red-team — parallel throughout
```

**Rule:** no later priority may be worked on as a primary deliverable until the earlier priority has either passed, or produced a documented localized failure that has been logged and triaged.

---

## 1. P0 — Foundation (Blocking)

### G1 — Explicit microscopic map `F`

**Deliverable:** one frozen candidate `(X_L, F_L)` with `p = 2` local realization, declared encoding, declared `L`, declared boundary/initial conventions.

**Attack:**
- Survey the arithmetic-dynamics literature for degree-2 algebraic self-maps with explicit `p = 2` realizations.
- Pick the smallest candidate satisfying the pre-registered class constraints.
- Freeze it as Deliverable A.
- No modification of the frozen tuple after this point.

**Kill condition:** no candidate in the pre-registered class can be constructed faithfully.

### G0 — Discrete/continuous typing discipline

**Deliverable:** a typed-space ledger recording every object used in P0 and P1 as living in exactly one space: `X`, `M(X)`, `T`, or `O`.

**Attack:**
- Run every claim in the frozen substrate through the typed ledger.
- Reject any claim that conflates spaces (e.g., `I_U = W^s(Γ_*)`).
- Record the ledger as a versioned artifact.

**Kill condition:** the frozen substrate cannot be described without at least one illegal identification.

---

## 2. P1 — Amplitude and Joint Arithmetic (Hardest, Unique)

### G2 — Joint arithmetic compatibility

**Deliverable:** a single report on the frozen `F_L` stating whether canonical-height and Markov/symbolic structures coexist on the *same* map.

**Attack:**
- Run Call–Silverman canonical-height applicability on `F_L`.
- Run Bowen-type Markov-partition applicability on the *same* `F_L`.
- If both fail or conflict: run the compatibility library (a set of pre-registered simple arithmetic maps) to determine whether joint satisfiability is possible in the class at all.
- Classify outcome: joint success / height-only / symbolic-only / neither.

**Kill condition:** no map in the pre-registered library satisfies both simultaneously.

### G3 — Amplitude-factor existence (no donor anywhere)

**Deliverable:** a T1–T4 amplitude-factor report documenting the existence or non-existence of a tame/Kronecker factor on the frozen `F_L`.

**Attack (T1 → T4):**
- **T1:** compute transfer/Perron-Frobenius operator on `F_L`; establish invariant measure and spectral data.
- **T2:** run Kronecker/point-spectrum diagnostic on the smallest admissible approximants; test whether a nontrivial tame factor survives in the projective limit.
- **T3:** if T2 holds, test whether the factor supports `U(1)` phase structure — run `χ_p` closure tests.
- **T4:** test whether the phase + modulus close into a complex amplitude and generator.
- If T2 fails (limit is weakly mixing with no tame factor), test the alternative amplitude constructions: `χ_p`-only, character-twisted, or non-arithmetic symbolic control.

**Kill condition:** no tame factor exists in the limit, and no alternative amplitude construction survives. This is the program's single sharpest kill gate.

---

## 3. P2 — Measure Existence (Shared Wall)

### G4 — Unique physical measure `μ_*`

**Deliverable:** constructed, characterized `μ_*` with regularity sufficient for finite Fisher information at the observable scale.

**Attack (branch-specific theorem shelf):**
- **Deterministic branch:** Lasota–Yorke-type estimates; ergodic decomposition; unique ergodicity results under verified hypotheses.
- **Post-trace stochastic branch:** RDS methods — strong Feller / asymptotic strong Feller; invariant measure and mixing theorems under actual hypotheses.
- Cross-check: compare with CDT finite-size scaling results as a numerical benchmark.
- Do not cite stochastic RDS theorems as though they applied to the deterministic pre-trace map.

**Kill condition:** no unique physical `μ_*` exists, or it is too singular for the IR generator.

---

## 4. P3 — Projection and Born (Shared Wall)

### G5 — Observable projection `Π`

**Deliverable:** explicit construction of `Π_s` and the state-to-theory prototype `Φ_s = A_s ∘ C_s`.

**Attack:**
- Construct `C_s` (coarse-graining map) for the frozen substrate.
- Construct `A_s` (effective-action reconstruction) from a pre-registered finite family of correlators.
- Test the commutation condition `Φ_{s'}(C_{s→s'} x) ≈ R_{s→s'}(Φ_s(x))` with declared tolerance.
- Failure at any stage means the bridge is not faithful at that level.

**Kill condition:** no faithful `Π` or `Φ` exists; the substrate cannot project to the effective theory.

### G6 — Born identification `Π_IR# μ_* = |ψ|²`

**Deliverable:** proof that the projected measure equals the Born distribution at the IR observable scale, with declared residual structure `δ_{ξ,ε}` below observational resolution.

**Attack:**
- Compute `Π_ε# μ_*` at coarse-graining scale `ε`.
- Compare to `|ψ|²` at the same scale.
- Prove `δ_{ξ,ε} → 0` in the declared Born regime.
- Verify the no-signaling constraint on `δρ_{a,b}(x_A,x_B)` for all allowed setting families.

**Kill condition:** the residual does not vanish, or it permits controllable signaling.

---

## 5. P4 — Quantum Recovery (BM Layer)

### G7 — Effective dynamics → Schrödinger/BM limit

**Deliverable:** derivation of the Schrödinger generator from the reduced dynamics.

**Attack:**
- Derive the generalized Langevin form from the substrate.
- Show controlled Markovian limit `K(t-s) → 2γ δ(t-s)`.
- Show strong-friction limit `v → ∇S/m`.
- Show `iκ∂_tψ = (-κ²/2m)∇²ψ + Vψ` emerges as the IR limit.

**Kill condition:** the IR generator is not in the linear Schrödinger class.

### G8 — Guidance law and quantum potential `Q`

**Deliverable:** derivation of `v = ∇S/m` and `Q = -(κ²/2m)(∇²√ρ)/√ρ` from the effective action.

**Attack:**
- Derive `ρ` and `S` from the amplitude-factor (G3) and measure (G4) outputs.
- Derive `Q` from the effective action's kinetic term.
- Show the guidance field and the quantum potential come from the same effective object.

**Kill condition:** `Q` does not emerge from the effective action and must be inserted by hand.

### G9 — Composition, entanglement, Bell, no-signaling

**Deliverable:** two-subsystem construction with `H_AB ≅ H_A ⊗ H_B` (or equivalent), Bell correlations, and no-signaling.

**Attack:**
- Build two weakly coupled frozen substrates.
- Test factorization in the `g → 0` limit.
- Generate entanglement for `g ≠ 0`.
- Reproduce Bell correlations.
- Run double-counting audit: does the mechanism require both nonlocal guidance and measurement-setting dependence for the same correlations?
- Verify no-signaling for all allowed setting families.

**Kill condition:** composition fails, or signaling is possible.

---

## 6. P5 — Recovery to Known Physics (Long-Horizon)

### G10 — QFT + Lorentz + gauge + GR + SM

**Deliverable:** recovery ladder QFT-R0 through GR-R2 (defined in v6.3.1 §17).

**Attack:**
- QFT-R0: relativistic effective field theory constructed.
- QFT-R1: Lorentz recovery to required precision.
- QFT-R2: gauge structure and matter content recovered.
- QFT-R3: no-go audit — Haag, Coleman–Mandula, Weinberg–Witten.
- QFT-R4: Euclidean → Lorentzian continuation.
- QFT-R5: AS fixed point compatible with emergent matter.
- GR-R0–R2: diffeomorphism, equivalence principle, Einstein dynamics.

**Kill condition:** failure at any rung means the program is not yet a TOE candidate.

---

## 7. Cross-Cutting Workstreams (Parallel Throughout)

### X1 — Pre-registered null models

**Deliverable:** a versioned null-model roster that every retained prediction must beat.

**Roster:**
- Standard BM typicality (no arithmetic substrate).
- BM + Valentini without arithmetic layer.
- Non-arithmetic discrete substrate (symbolic/graph-dynamical control).
- Plain AS / direct effective-flow control.

**Rule:** no prediction is promoted unless it beats at least one null model on shape, not on fitted amplitude.

### X2 — Spectral dimension `d_s(k)`

**Deliverable:** computed `d_s(k)` for the frozen substrate.

**Attack:**
- Compute return probability `P(σ)` on the substrate.
- Compute `d_s(σ) = -2 d ln P / d ln σ`.
- Compare to the five-way convergence target `d_s^UV ≈ 2`, `d_s^IR ≈ 4`.

**Kill condition:** `d_s` flat at 4, or no crossover.

### X3 — Arrow of time / irreversibility

**Deliverable:** derivation of the thermodynamic arrow from the coarse-graining direction.

**Attack:**
- Show the RG flow direction is physical (not just bookkeeping).
- Derive the second law from the coarse-graining.
- Address the `needInitialCondition` debt.

**Kill condition:** irreversibility must be inserted by hand.

### X4 — Independent red-team

**Deliverable:** versioned external review at each milestone gate.

**Attack:**
- Independent reviewer must be a different model family and/or a named human.
- Reviews tied to specific debt classes.
- Milestone gates: after Phase 1, after Phase 2, after Phase 3, before any promoted physical claim.

**Kill condition:** red-team cannot be convened, or review produces only process artifacts.

### X5 — Prime universality

**Deliverable:** scan across `p = 2, 3, 5, 7, ...` with classification.

**Attack:**
- Run P0 and P1 for each prime.
- Classify outcome: prime-independent universality / unique selected prime / prime-dependent physics / adelic necessity.
- Reject any prime chosen because it makes the numbers work.

### X6 — Scheme independence (AS-specific)

**Deliverable:** regulator/truncation stability report for every derived physical constant.

**Attack:**
- Run every derived constant through multiple regulators.
- Run every derived constant through multiple truncations.
- A constant that moves materially is demoted to coordinate artifact.

**Kill condition:** all derived constants move under legitimate scheme changes.

---

## 8. Kill Gates and Stop Conditions (Ordered)

| Priority | Gate | Stop condition |
|---|---|---|
| P0 | G1 | No faithful `F` in pre-registered class |
| P0 | G0 | Illegal identification unavoidable |
| P1 | G2 | No joint satisfiability anywhere in library |
| P1 | G3 | No tame factor in limit, no alternative amplitude construction |
| P2 | G4 | No unique `μ_*`, or too singular |
| P3 | G5 | No faithful `Π` or `Φ` |
| P3 | G6 | Residual does not vanish, or permits signaling |
| P4 | G7 | IR generator not in linear Schrödinger class |
| P4 | G8 | `Q` must be inserted by hand |
| P4 | G9 | Composition fails, or signaling possible |
| P5 | G10 | Failure at any rung |
| X | X1 | Null model beats every retained prediction |
| X | X2 | Spectral dimension flat at 4 |
| X | X6 | All derived constants move under scheme change |
| X | Occam | Arithmetic layer produces no distinctive signature |

Each kill is **localized**, **documented**, and **versioned**. A kill is progress. A kill that is not logged is entropy.

---

## 9. Sequential Dependency Chart

```
P0: G1 ──► P1: G2 ──► P1: G3 ──► P2: G4 ──► P3: G5 ──► P3: G6
                                                             │
                                                             ▼
                                          P4: G7 ──► G8 ──► G9
                                                             │
                                                             ▼
                                                       P5: G10
                        │
                        ▼
                 X1–X6 cross-cutting (parallel from P0 onward)
```

**Rule:** P0 unblocks everything. P1 is the load-bearing branch — no P1, no program. P2 and P3 are the shared walls. P4 is the quantum layer. P5 is downstream. Cross-cutting X1–X6 runs in parallel.

---

## 10. First 90 Days — Concrete Calendar

### Days 1–30
- **G1:** select and freeze `(X_L, F_L)` at `p = 2`.
- **G0:** typed-space ledger on the frozen substrate.
- **G2 (partial):** canonical-height applicability report.
- **G2 (partial):** Markov/symbolic applicability report on same map.
- **G3 (T1–T2):** transfer/Koopman diagnostic and first tame-factor test.
- **X5:** prime-universality scan begins.
- **X1:** null-model roster frozen.

### Days 31–60
- **G2:** joint-satisfiability classification.
- **G3 (T3–T4):** phase closure test; complex amplitude closure test.
- **G4:** invariant measure candidate constructed.
- **X2:** first spectral-dimension computation.
- **X6:** first scheme-independence run.

### Days 61–90
- **G5:** first `Π` prototype on the one-qubit problem.
- **G6:** first Born identification test.
- **G3:** amplitude-factor verdict.
- **X4:** first external red-team review.
- **Stop/continue decision** based on G1–G3 outcome.

---

## 11. Decision Points

| Decision point | When | Trigger | Consequence |
|---|---|---|---|
| DP1 | After Day 30 | G1 outcome | If no faithful `F`: stop arithmetic branch |
| DP2 | After Day 60 | G2 + G3 outcome | If joint satisfiability fails: select height-only, symbolic-only, or abandon |
| DP3 | After Day 90 | G3 verdict | If no tame factor and no alternative: **primary kill** — program reduces to shared-wall problem |
| DP4 | After Day 180 | G4 + G5 outcome | If no unique measure or faithful projection: shared-wall failure |
| DP5 | After Day 365 | G6 outcome | If Born not derived: program is pre-quantum |
| DP6 | After Phase 4 | G7–G9 outcome | If quantum layer not recovered: not a TOE candidate |
| DP7 | After Phase 7 | G10 outcome | If QFT/GR recovery fails: not a TOE |

---

## 12. One-Page Summary

**The program's attack plan is:**

1. **Freeze one arithmetic substrate** (G1) with `p = 2`.
2. **Verify joint arithmetic compatibility** (G2) — canonical height + Markov on the same map.
3. **Construct the amplitude factor** (G3) — the program's unique burden; no donor.
4. **Construct the invariant measure** (G4) — the shared wall.
5. **Construct the projection** (G5) and **derive Born** (G6) — the shared wall.
6. **Recover the quantum layer** (G7–G9) — Schrödinger, guidance, `Q`, composition, Bell.
7. **Recover known physics** (G10) — QFT, Lorentz, gauge, GR.
8. **Run cross-cutting workstreams** (X1–X6) in parallel throughout.
9. **Kill cleanly** at any gate whose stop condition triggers.

**The single most important gate is G3.** It is the only gate with no donor anywhere in the surveyed landscape. It is the program's unique burden and its unique bet.

**The single most important discipline is X4.** Without independent red-team, the program is self-graded.

**The single most important rule is:** kill locally, log formally, proceed only on the survivors.

**The program's success criterion is not "the theory is proven."** It is: every gate either passes or fails with a documented, reproducible, localized outcome that reduces the hypothesis space.

---

## 13. What to Do Tomorrow

If only one thing can be done tomorrow:

> **Select and freeze the smallest pre-registered degree-2 algebraic self-map with `p = 2` local realization, and write down its tuple `(X_L, L, encoding, F_L, placewise action, boundary/initial conventions, allowed parameters)`.**

Nothing else matters until that object exists.
```