# BM-IST POC Seed

## Exact claim

> Under the assumptions that (a) the coarse-grained dynamics of a density ρ
> and phase S admits a Hamiltonian formulation, (b) a complex field
> ψ = √ρ · exp(iS/κ) is defined from this pair, and (c) the induced evolution
> of ψ is required to be a linear, local Schrödinger-type equation
> iκ∂ₜψ = [-κ²/(2m)Δ + V]ψ — then the quantum-correction energy functional
> added to the classical Hamilton–Jacobi energy is forced to take the
> Fisher-information form:
>
> E_q[ρ] = (κ²/8m) I_F[ρ] = (κ²/8m)∫(|∇ρ|²/ρ)dx = (κ²/2m)∫|∇√ρ|²dx
>
> under the further restriction that E_q is local, depends only on ρ and ∇ρ
> (no higher derivatives, no explicit S-dependence), C², and rotationally
> invariant.

This is referred to in the source material as "Theorem 1" / the
linear-amplitude → Fisher-rigidity result. It is **not** the same claim as the
separately-discussed diffusive-corner/Fisher-monotonicity restriction ("T2")
— the two should not be conflated in Conductor task naming.

## Source / location in the corpus

Documented across this research thread's multi-round Fisher-selection
investigation (the DR2/DR3 exchange reviewing Claude/Kimi/Gemini/Z proposals
for whether deterministic IST can produce the quantum amplitude layer), and
in the subsequent turn where the derivation's algebra was independently
hand-verified and its literature precedent was checked. No standalone
file-based corpus exists yet; this conversation is currently the corpus.

## Current evidence already established

- **Map stated**: domain (Hamiltonian (ρ,S) systems + the ψ ansatz demanding
  linear evolution), codomain (local correction functionals E_q[ρ]), and
  translation rule (Euler–Lagrange matching) are all explicitly given.
- **Invariant stated**: the specific forced functional form and coefficient
  (κ²/8m) are stated precisely enough to be checked, not left as a loose
  proportionality.
- **Independent algebraic verification performed**: the Euler–Lagrange
  matching was redone by hand in this thread, confirming the identity
  Δ√ρ/√ρ = ½Δρ/ρ − ¼|∇ρ|²/ρ², confirming the forced coefficient a(ρ) =
  κ²/(8mρ), confirming the |∇ρ|² term's coefficient is automatically
  consistent (not a second free condition), and confirming b′(ρ) = 0 is
  required. This is analytic/symbolic verification of the stated derivation
  — not a separately constructed numerical toy-model instance.
- **Literature precedent identified and checked**: this result is not novel
  in isolation — it closely tracks Reginatto (1998, *Phys. Rev. A* 58, 1775)
  and Hall & Reginatto (2002, *J. Phys. A* 35, 3289), confirmed via search in
  this thread. The IST-native question (does *this* substrate produce the
  linear-amplitude precondition in (c)) remains separately open and is out of
  scope for this claim.

## Current unresolved debt

- **Uniqueness is claimed but not yet defended or disclaimed.** The result
  is stated as forcing the *unique* local correction under the given
  hypotheses; no rival explanation has been ruled out, and no explicit
  disclaimer narrowing the uniqueness claim has been recorded.
- **A specific, already-identified obstruction has been stated but not yet
  assessed.** The hypothesis set explicitly excludes S-dependent corrections.
  Later work in this thread identified concrete candidate mechanisms
  (cocycle-twisted lifts, projective-representation constructions) as
  plausible sources of exactly this kind of S-dependence — but nobody has
  actually constructed such a candidate and checked whether it breaks the
  stated hypotheses.
- **No formal faithfulness check has been run.** The specific EBP question
  ("does this formalization match the intended physical claim?") has not
  been posed and answered by a reviewer.

## EBP v2.1 debt status

**Retired:**
- `needMap`
- `needInvariant`
- `needToyCheck` (retired via the analytic verification above — flagged
  explicitly as symbolic/analytic in nature; a reviewer may still request an
  independent numerical instantiation if that is judged insufficient)

**Open:**
- `needNullModel`
- `needObstruction`
- `needFaithfulnessReview`

## Smallest next useful research action

Retire `needNullModel` via EBP's own Option B: restate the claim's scope
explicitly — "the unique local correction *within the stated hypothesis
class*," not "the only possible quantum correction of any kind." This is a
scope clarification, not new research, and is the cheapest legitimate move
available under EBP v2.1 §6.4.

This should precede `needFaithfulnessReview`, since faithfulness review
should be asked against the claim's final, precisely-scoped wording rather
than re-run after a later scope change.

`needObstruction` is the deeper, substantive research item remaining
regardless of the POC's outcome: construct one candidate S-dependent E_q and
determine whether it evades the stated hypotheses.

## Why this claim is suitable for the POC

- It already exists in the corpus with real, checkable prior work — nothing
  here is invented for the demonstration.
- It has genuine, specific, already-characterized remaining debt (not a
  strawman gap), including one item with a nearly-free retirement path and
  one item that is real, open research.
- It maps onto EBP v2.1's six canonical debt items with no need for
  domain-specific extensions.
- It is narrow enough — one derivation, one scope, three open items — to
  carry through capture → clarify → debt-retirement → human faithfulness
  gate → promotion within a single demonstrable pilot pass, without pulling
  in the larger, still-open amplitude-layer research program it sits inside.
