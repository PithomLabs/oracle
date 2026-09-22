# BM–IST–AS Attack Roadmap
### Prioritized implementation plan, by category — companion to the Gap-Filling Synthesis and v6.3.1

**Rule for using this doc:** work top category first. Do not start a lower category before its blocking dependency is either closed or explicitly killed. Every item ends in a checkable artifact, not a discussion.

---

## Priority 0 — Process Freezes (before any construction)

*No dependency. Cheapest tier. Must close before Tier 1 to avoid contamination.*

| # | Action | Artifact | Exit criteria |
|---|---|---|---|
| 0.1 | Freeze substrate-class search budget | Written list: 4 pre-registered classes (deg-2 arithmetic \(p{=}2\); \(p{=}3\); multi-place \(p{=}2{\times}3\); non-arithmetic control) | No 5th class added without written justification |
| 0.2 | Freeze null-model roster | List: standard BM typicality, BM+Valentini w/o arithmetic layer, non-arithmetic discrete substrate, plain-AS control, CDT, causal-set growth, GFT condensate | Roster locked before Tier 1 results exist |
| 0.3 | Freeze no-signaling residual family | Formal statement: \(\Delta_A(x_A;a,b)=\Delta_A(x_A;a)\), setting-indexed | Pre-registered before any composition work (Tier 5) |
| 0.4 | Resolve Valentini fork (theory-level) | Written position: is sustained lab non-equilibrium claimed preparable, yes/no, with reasoning | Position frozen before Born-chain work (Tier 2) |
| 0.5 | Appoint independent reviewer | Named different-model-family or human reviewer | Reviewer confirmed, not self-assigned |

---

## Priority 1 — Substrate Construction (Gates G0–G2)

*Blocks everything. Highest-leverage tier: nothing downstream is real without this.*

| # | Action | Gate | Artifact | Donor to consult |
|---|---|---|---|---|
| 1.1 | Construct \(F_L\): deg-2 arithmetic map, explicit \(p{=}2\) local realization | G1 | Frozen tuple \((X_L,L,\text{encoding},F_L,\text{conventions})\) — **A-freeze rule applies after this** | Wolfram/causal-set tooling (code only, not ontology) |
| 1.2 | Run canonical-height applicability test on \(F_L\) | G2 (half) | Pass/fail: does \(\hat h_F=0\) locus exist and behave as required | Call–Silverman |
| 1.3 | Run Markov-partition applicability test on the **same** \(F_L\) | G2 (half) | Pass/fail: explicit partition, branching number reported | Bowen / symbolic dynamics |
| 1.4 | Joint-satisfiability check: does 1.2 and 1.3 hold on the *same* map | G2 (joint) | Single verdict — compatible / incompatible | Check known coexistence examples (e.g. doubling-type maps) first — cheap, may resolve before new work |
| 1.5 | Build non-arithmetic control substrate in parallel | Null model | Comparable-complexity Wolfram-class object | Wolfram tooling |
| 1.6 | **Kill check:** no faithful \(F_L\) in budget → arithmetic branch stops | G1 kill | Stop/continue decision, recorded | — |

---

## Priority 2 — Amplitude Gate (G3)

*The one gate with zero donor support anywhere in the QG survey. Scientifically the highest-value, hardest tier. Start diagnostics as soon as 1.1 exists — do not wait for 1.4.*

| # | Action | Artifact | Kill condition |
|---|---|---|---|
| 2.1 | T1: confirm invariant measure / transfer-Koopman objects exist on \(F_L\) | Existence report | Fails → gate closed, no rescue |
| 2.2 | T2: cheap diagnostic — Kronecker/point-spectrum computation on smallest admissible approximants | Spectrum table, finite-\(L\) | Weak mixing + no tame factor + no alternative → **kill route** |
| 2.3 | T2 (continued): projective/infinite-limit consistency test — does a nontrivial factor survive \(L\to\infty\) | Convergence/divergence verdict | Factor vanishes as \(L\to\infty\) → kill route |
| 2.4 | T3: test whether surviving factor supports \(U(1)\) structure via \(\chi_p(x)=e^{2\pi i\{x\}_p}\) | Pass/fail on character route | No \(U(1)\) structure found → route fails, no substitute currently known |
| 2.5 | T4: phase + modulus closure into a complex amplitude generator | Constructed generator or documented failure | — |
| 2.6 | Unify: check whether \(\chi_p\), Markov coding, and canonical-height admissibility are one factor in disguise | Single unification report | — |

**This tier has no fallback donor. If it fails, the program does not get to borrow a replacement — record the failure as the program's own, not import a rescue.**

---

## Priority 3 — Measure & Born Chain (G4, G6, B0–B3)

*Depends on 1.x. Can begin once \(F_L\) is frozen; does not need G3 to finish first.*

| # | Action | Artifact |
|---|---|---|
| 3.1 | B0: prove existence/uniqueness of \(\mu_*\) | Existence-uniqueness theorem or obstruction |
| 3.2 | B1: regularity at observable scale \(\epsilon\) — finite Fisher information | Computed \(I_F[\rho_\epsilon]\) |
| 3.3 | B2: Born identification at IR observable algebra; report finite-resolution texture \(\delta_{\xi,\epsilon}\) separately, never conflated with B2 itself | Two numbers: IR-exact statement + measured/predicted texture |
| 3.4 | No-signaling safety check against 0.3's frozen family | Pass/fail per setting pair |
| 3.5 | B3: Valentini relaxation as path only — confirm destination is set by 3.1–3.3, not by relaxation dynamics | Written confirmation, no new mechanism introduced |

---

## Priority 4 — Flow Construction (dFRG, \(\Phi\), fixed points)

*Depends on 1.x + partial 3.x. This is where AS actually does work rather than lending vocabulary.*

| # | Action | Artifact |
|---|---|---|
| 4.1 | Build prototype \(\Phi_s=\mathcal A_s\circ\mathcal C_s\): explicit coarse-graining + effective-action reconstruction | Constructed map, not slogan |
| 4.2 | Test RG-commutation: \(\Phi_{s'}(\mathcal C_{s\to s'}x)\simeq\mathfrak R_{s\to s'}(\Phi_s(x))\) | Tolerance or exact-theorem verdict |
| 4.3 | Search for fixed point \(\Gamma_*\) of the constructed flow | Found / not found, with truncation details logged |
| 4.4 | Regulator/scheme cross-validation | Table: quantity vs. scheme — flag anything scheme-dependent as artifact, not result |
| 4.5 | Compute \(d_s(k)\), compare to the 5-way convergence (AS/LQG/causal-sets/CDT ≈2 in UV) | Number, plotted against convergence target — convergence with others is context, not confirmation |

---

## Priority 5 — Guidance Law & Quantum Potential (G7–G8)

*Depends on 3.x + 4.x.*

| # | Action | Artifact |
|---|---|---|
| 5.1 | Construct reduced Langevin/Kramers dynamics from traced-out bath (requires internal-environment gate §3.5 licensed first) | Explicit \(K(t-s)\), \(\eta(t)\) |
| 5.2 | Take strong-friction limit; check convergence to \(v=\nabla S/m\) | Convergence proof or documented gap, tagged `needRegularity` |
| 5.3 | Derive \(Q=-\kappa^2/2m\cdot\nabla^2\sqrt\rho/\sqrt\rho\) from the effective action — do not insert by hand | Derivation or explicit non-derivation |
| 5.4 | Reuse Reginatto/Hall–Reginatto machinery once a continuum configuration space is actually reached (not before) | Applied derivation, citation-linked |
| 5.5 | Asymptotic-unitarity check (Lindblad → Schrödinger as \(k\to0\)) | Pass/fail |

---

## Priority 6 — Composition, Entanglement, Bell (G9)

*Design work can start early (cheap); construction depends on 5.x.*

| # | Action | Artifact |
|---|---|---|
| 6.1 | **Double-counting audit** (design-time, do this now): does the mechanism need both nonlocal guidance *and* setting-dependence for the same correlations? | Written audit — simplify away redundancy before construction begins |
| 6.2 | Construct \(\mathcal H_{AB}\simeq\mathcal H_A\otimes\mathcal H_B\) or equivalent from two weakly-coupled subsystems | Constructed structure or obstruction |
| 6.3 | Consult tensor-network/QEC literature for locality-as-code-property framing | `needMap`: transfer boundary-code theorems to tame-factor setting |
| 6.4 | Derive entanglement generation + Bell statistics + no-signaling jointly | Three-part verdict against 0.3's frozen family |

---

## Priority 7 — Relativistic & Gravitational Recovery (G10)

*Long-horizon. Do not front-load effort here before Priority 1–6 are further along.*

| # | Action | Artifact |
|---|---|---|
| 7.1 | QFT-R0: relativistic EFT recovery | Recovery theorem or obstruction |
| 7.2 | QFT-R1: Lorentz symmetry to tested precision | Bound comparison |
| 7.3 | Gauge structure recovery | Recovery theorem or obstruction |
| 7.4 | GR recovery, treated as downstream/residue (per Jacobson/entropic/holographic/Verlinde precedent) — precedent licenses the *strategy* only | Recovery theorem or obstruction |

---

## Priority 8 — Governance & Empirical Program (runs continuously, all tiers)

*Not a phase — a parallel track from Priority 0 onward.*

| # | Action | Cadence |
|---|---|---|
| 8.1 | Independent-reviewer sign-off on every promoted claim | Per claim |
| 8.2 | External red-team review | At every milestone gate (end of each Priority tier) |
| 8.3 | Subtraction record update (log every dropped/demoted branch, with reason) | Per tier |
| 8.4 | Pre-registration: "shape before amplitude" — freeze predicted functional forms before comparing to data | Before any T2/T3 empirical claim (below) |
| 8.5 | EBP ledger maintenance: `CLAIM_ID`, debt class, status, per new claim | Continuous |

---

## Priority 9 — Empirical Handle Table (shape before amplitude)

*Depends on Priority 3–5 producing actual predicted shapes. Do not fit to data before that.*

| # | Action | Artifact |
|---|---|---|
| 9.1 | T1: compile constraints, declare formula families, pre-register — before any comparison | Frozen pipeline document |
| 9.2 | T2: derive shapes — mesoscopic interferometry scaling, colored-noise spectrum, spectral-dimension flow, relic-nonequilibrium signature | Functional forms, not fits |
| 9.3 | T3: high-information fingerprints (modular statistics, prime dependence, cross-substrate universality) | Only after 9.1–9.2 exist |

---

## Priority 10 — Deferred / Stress-Test Only (do not actively work)

*Explicitly downstream. These are stress tests on a finished chain, not construction targets. Revisit only after Priority 1–3 close.*

- Black hole information paradox connections
- Cosmological constant / \(\Lambda\) connections
- Preferred foliation / relativity reconciliation
- Inflation & primordial power spectra
- Galactic rotation curves / "dark matter as \(Q\)"
- Problem of Time (WDW / Page–Wootters)
- E8 attachments (packing, CFT) — both already gated behind Priority 2 and Priority 4/5 respectively; do not revisit until those close

---

## One-page summary (the whole roadmap in order)

```
0. Freeze process (budget, null models, no-signaling family, Valentini fork, reviewer)
1. Construct + freeze F_L; run height + Markov + joint-satisfiability tests
2. Attack the amplitude gate (T1-T4) — no donor exists; start diagnostics as soon as F_L exists
3. Prove/construct mu_* and the Born chain (B0-B3), no-signaling check
4. Build the flow (Phi, RG-commutation, fixed point, d_s(k))
5. Derive guidance law and Q from the flow (not by hand)
6. Audit double-counting, then construct composition/Bell/no-signaling
7. Recover relativistic QFT + gauge + GR (long horizon)
8. Governance runs throughout: review, red-team, subtraction record, ledger
9. Only once 3-5 yield real shapes: pre-register and compare to data
10. Six thought-experiment domains + E8: stress tests, not construction work — leave parked
```
