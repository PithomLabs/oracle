# BM–IST(-AS) v5 — Ledger Addendum: Valentini H-Theorem Entries

**Governs:** BM–IST Synthesis v5.1 (§14–17, §35) and the v5 Implementation Plan (§1, §12)
**Adds:** three new `CLAIM_ID` entries, per the ledger convention established in v5.1 and formalized in the Implementation Plan's §1 (status-tag translation) and §1.2 (function-class typing)
**Entry status:** all three enter as **debt**, not as promoted claims — per EBP discipline, they are added to the ledger now so they're tracked, not re-derived from scratch next time this thread returns to the probability gates

---

## CLAIM_ID: PROB-XI0-COARSE-1

**Statement.** The coarse-graining cell size used in Valentini-style dynamical relaxation (the averaging scale at which \(H=\int \rho\ln(\rho/|\psi|^2)\,dq\) is computed) is identified with **ξ₀**, BM-IST's existing correlation scale — rather than being introduced as an independent free parameter, which is how Valentini's own program has always carried it.

**Function class.** Selection principle / parameter-identification claim. This is *not* an existence claim (it doesn't assert relaxation happens) and *not* a construction (it doesn't build the flow) — it asserts that two quantities already appearing separately in the architecture are the same quantity. Per the Implementation Plan's §1.2 typing rule, the check this class requires is: **is the identification forced by the dynamics, or merely convenient?**

**Status.** `[SPECULATIVE-SHAPE]` / Workbench rung **Postulates**. Not `[BORROWED]` — this is a same-architecture parameter reuse, not an import from an unrelated field, so it doesn't carry the borrowed-rigor risk flagged against the AS-template and Hecke-algebra episodes. But it is unproven: nothing yet shows the physical process that fixes ξ₀ elsewhere in the architecture is the *same* physical process that would set a relaxation-averaging cell.

**Why it matters if true.** Valentini's program has been criticized for decades for putting its coarse-graining scale in by hand. If this identification holds, BM-IST removes one of relaxation theory's oldest free parameters rather than adding one of its own — a genuine parsimony gain, not a relabeling (contrast with PROB-XI0-COARSE-1's neighbor claims in the AS-template review, which mostly relabeled problems).

**Open debt.**
- `needMap` — the identification has not been constructed. Retiring this requires showing, from the actual IST shift-map dynamics (not by analogy), that Valentini's own coarse-graining prescription applied to that dynamics yields a scale equal to (or simply, derivably related to) ξ₀ as already fixed elsewhere in the architecture.
- `needObstruction` — no failure condition has been stated yet (see Kill condition below, which supplies one).

**Kill condition.** If a direct computation of the relaxation-averaging scale — applying Valentini's coarse-graining prescription to the actual substrate dynamics, not to a generic chaotic system — yields a scale that is not equal to, or simply related to, ξ₀ as independently fixed, the identification fails. ξ₀ reverts to its prior role only; ξ₀-as-coarse-graining-scale is not reintroduced without new evidence.

**Dependency / phase placement.** Phase 2 (Measure & Probability Theorems) of the Implementation Plan, and only after Phase 1's scaling audit has fixed ξ₀'s actual regime — attempting this before Phase 1 completes would mean checking an identification against a value that hasn't been computed yet.

---

## CLAIM_ID: ARCH-PROB-PRIMACY-1

**Statement.** This is an **architectural decision entry**, not a derivation claim — matching the P1-vs-P2 decision format already in the Implementation Plan (§Phase 5). The decision: **T-Born (arithmetic equidistribution, v5.1 §17) is designated the primary mechanism for observable Born-rule statistics in laboratory/terrestrial, thermalized regimes. Valentini-H (dynamical relaxation) is scoped specifically to regimes that decoupled before relaxation could complete** — i.e., early-universe relic sectors (see HANDLE-RELIC-NONEQ-1 below) — rather than treating the two as competing, parallel, unexamined candidates for the same gate.

**Function class.** Architectural decision. Not existence, not construction — this is a scope assignment, and per §1.2's discipline it still needs a stated basis and a stated reversal condition rather than being adopted by preference.

**Rationale for this default (reversible, not final).** T-Born's mechanism gives exponential convergence, which is the only one of the two that explains why terrestrial quantum statistics show *no* observable relaxation residue — matching what's actually been measured. Valentini's own numerics generically show power-law, sometimes incomplete, relaxation, and his own program's predicted signatures are specifically aimed at systems that decoupled early (exactly where "no residue observed yet" is not a settled fact, since no one has looked in the right place). The two mechanisms are not actually competing for the same explanatory job once split by regime — this scoping is a genuine division of labor, not a forced choice.

**Status.** `[THEOREM-SHAPED]` for the underlying claim each route makes about its own regime (T-Born's exponential convergence, Valentini's regime-dependent relaxation completeness are each real, attemptable, or already-published results); `[SPECULATIVE-SHAPE]` for the *division of labor itself*, since nothing has yet checked that the boundary between "thermalized" and "decoupled-early" regimes is drawn in the right place for BM-IST's specific substrate.

**Open debt.**
- `needFaithfulnessReview` — does this division of labor actually match each mechanism's own stated domain of applicability in its source literature, or has it been drawn to fit BM-IST's convenience? This is the same check that caught the Hecke-algebra category error and the unconditional-Furstenberg overclaim in earlier reviews — apply it here too, deliberately, before promotion.
- `needObstruction` — no obstruction has yet been run (see Kill condition).

**Kill condition.** If Phase 2's actual computation of T-Born's convergence rate (once T-Prob-1's substrate-specific analysis is done, not assumed by analogy to Bilu/Brolin–Lyubich in general) turns out not to be exponential for BM-IST's specific arithmetic regime, T-Born loses its claim to primacy in the thermalized regime and this entry must be revisited — it does not default to Valentini-H by elimination; the primacy question reopens from scratch. Separately, if HANDLE-RELIC-NONEQ-1 (below) returns a null result inconsistent with any relic-nonequilibrium signature at any scale Valentini's mechanism would predict one, that narrows (but does not by itself kill) Valentini-H's scoped role.

**Dependency / phase placement.** Sits between Phase 2 and Phase 6 of the Implementation Plan — it cannot be finally settled until Phase 2 produces T-Born's actual convergence rate, but the *scoping decision itself* (which regime each mechanism is responsible for) can be adopted now, provisionally, so Phase 6's handle-table work (which needs HANDLE-RELIC-NONEQ-1, see below) isn't blocked waiting on Phase 2.

---

## CLAIM_ID: HANDLE-RELIC-NONEQ-1

**Statement.** Add a new row to the empirical handle table (v5.1 §35 / Implementation Plan §12): **relic quantum nonequilibrium** — the prediction, already established in Valentini's own published program (not new to this architecture), that primordial gravitons and/or relic neutrinos which decoupled before dynamical relaxation completed could retain observable deviations from Born-rule statistics, with candidate signatures in CMB polarization and related early-universe observables.

**Function class.** Phenomenological fit / existing empirical handle. This is the cheapest and lowest-risk of the three entries: it imports an already-published, already-peer-reviewed prediction wholesale rather than constructing anything new. It does **not** claim BM-IST derives this prediction — only that BM-IST's architecture, if ARCH-PROB-PRIMACY-1 is adopted, has a natural place to attach it.

**Status.** `[EXISTING]` for the underlying Valentini-literature prediction itself; `[SPECULATIVE-SHAPE]` for its attachment to BM-IST specifically, since nothing yet ties BM-IST's own parameters (ξ₀, κ, the decoupling epoch t_ent already used elsewhere in the handle table) to Valentini's predicted signature strength or scale.

**Open debt.**
- `needMap` — has BM-IST's own parameter set actually been connected to this prediction's observable magnitude, or is this row currently just sitting adjacent to the rest of the table without a computed link?
- `needFaithfulnessReview` — is the import faithful to what Valentini's program actually predicts (regime, magnitude, observable), not a garbled or strengthened version of it?

**Kill condition / governing discipline.** This row is subject to the same discipline already governing every other row in the handle table: it must be checked against the **same fixed (c, ξ₀, κ) parameter set** used everywhere else in the architecture — no separate fitting for this row alone. If satisfying this row's predicted signature scale requires a different ξ₀ or κ than the rest of the table already uses, that is a genuine internal inconsistency (a C1-type failure, not just an unfavorable result) and must be reported as such, not quietly absorbed by re-tuning.

**Dependency / phase placement.** Phase 6 (conditional) of the Implementation Plan — like the rest of the handle table, this only becomes an active work item once Phase 5's P1/P2 decision is made; if P1 (trace-relative, no autonomous p-adic-sector content) is retained as the default, this row still applies (it concerns the archimedean/observable sector's own relaxation history, not the p-adic bath's ontological status), so it is **not blocked** by the P1/P2 decision the way the gravity/dark-matter rows are — worth noting explicitly, since it means this row can be pursued independently and earlier than the rest of Phase 6.

---

## Ledger Summary Table

| CLAIM_ID | Function class | Status | Blocking debt | Phase |
|---|---|---|---|---|
| PROB-XI0-COARSE-1 | Selection principle | Postulates | `needMap` | Phase 2, after Phase 1 |
| ARCH-PROB-PRIMACY-1 | Architectural decision | Postulates (division), Theorem-shaped (components) | `needFaithfulnessReview` | Phase 2–6, adoptable now provisionally |
| HANDLE-RELIC-NONEQ-1 | Phenomenological fit | Existing (import), Speculative (attachment) | `needMap` | Phase 6, unblocked by P1/P2 |

All three follow Gate R (Implementation Plan §1.3) before any future promotion past `needToyCheck`: no claim above should be treated as more than tracked debt until its stated kill condition has actually been run.
