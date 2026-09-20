Adopting EBP v2.1's grammar starting now, and the first thing worth doing with it is a direct reconciliation against the debt vocabulary I'd been using informally for many rounds — because if I don't retire that vocabulary explicitly, we get exactly the drift the user is asking me to prevent.

**Vocabulary reconciliation, done once, so it doesn't need doing again:**

| My earlier ad hoc term | EBP v2.1 canonical term | Note |
|---|---|---|
| `needCitationCheck` | folded into `needMap`'s retirement condition | EBP has no dedicated prior-art debt item; per its own anti-bureaucracy stance, citation work belongs inside stating the map's translation rule, not a new debt category |
| `needCounterexampleAttack` | `needObstruction` | Exact match — "state one obstruction, explain whether it applies, record the consequence" is the attack track, already canonical |
| `needHumanMathematicalReview` | `needFaithfulnessReview` | Exact match, and EBP's version is more precisely specified: yes/no/contested, with a mandatory one-sentence mismatch statement on contest — no silent vetoes |
| `needDependencyImpactAnalysis` | **no EBP equivalent** | This is the important one — see below |

**The load-bearing structural gap, and it needs a decision before anything gets built, not after.** EBP v2.1's `ideaPromoted` is a pure computed predicate (`debt.isEmpty && !containsFinalTruthClaim`), not a stored status field — which is elegant (there's no stale "promoted" flag that can desync from reality) but it also means EBP, as specified, has **no cross-idea dependency primitive at all**. Section 2.5 ("a promoted idea can become unpromoted if new evidence creates new debt") describes what happens to *that one idea* when debt is added to it — it says nothing about what happens to a *different* idea that cited it. This is a real, meaningful difference from Solvent, whose `belief_edge` + `RetractCascade` exists specifically to propagate exactly that consequence atomically. If T1 (Fisher rigidity) gets un-promoted because new debt appears, EBP's kernel gives no mechanism for automatically flagging every idea that leaned on T1's promoted status. That gap needs an explicit decision, not a silent assumption — I've made it Phase 0's first item below.

---

**Phase 0 — Stand up the EBP store as its own minimal service, outside Conductor and Solvent's cores.**

Implement section 24's five verbs exactly, nothing more: `ebp_capture`, `ebp_status`, `ebp_retire`, `ebp_add_debt`, `ebp_promote`, exposed as MCP tools, backed by a store whose schema is the Lean structures in section 4 verbatim (`Idea { owner, claim, source, born, debt : List DebtItem, containsFinalTruthClaim }`). `ideaPromoted` stays a query, never a cached flag, per the spec.

Decide the cascade gap now: for this pilot, I'd recommend deliberately **not** building an EBP-native cascade. Instead, any cross-idea dependency ("T2's citation of T1 needs re-review if T1 changes") becomes an ordinary Conductor work-dependency between the *tasks* that track those ideas — a human or Skill notices T1 changed, and the Conductor task graph is what surfaces "this downstream task should be re-opened." This is a real, named limitation, not a hidden one — consistent with growing the ecosystem only where evidence forces it rather than pre-building a Solvent-style cascade EBP doesn't ask for.

**Phase 1 — Pizza Bot as the Agent, one Skill per non-human-gated debt item.**

Pizza Bot's Skills map almost exactly onto EBP's debt taxonomy: a `map-builder` Skill retires `needMap`, a `null-model-hunter` retires `needNullModel`, an `obstruction-finder` retires `needObstruction`, a `toy-checker` retires `needToyCheck` (this one can genuinely shell out to a Lean or symbolic-math tool via its own MCP connection). Deliberately **no autonomous Skill for `needFaithfulnessReview`** — per section 6.6's own rule, that debt item retires only through Pizza Bot's native human-in-the-loop approval gate, which already exists and needs no new engineering. A Skill may propose that faithfulness holds; only a human's approval in Pizza Bot's Action queue actually calls `ebp_retire needFaithfulnessReview`.

This is the piece that most directly answers the original "not autonomous, minimal UI, prompt external AI" requirement from many rounds ago, and it's already built: Pizza Bot's Unread/Action split *is* the minimal UI — a promoted idea lands in Unread, a pending debt-retirement or faithfulness decision lands in Action as an approval request. Bring-your-own-model means Kimi, GPT, Gemini, and Claude can each be wired as different Skills or model configs, formalizing what this whole investigation has been doing by hand across thirty-plus rounds of manual relay.

**Phase 2 — Conductor coordinates who's retiring what, and in what order.**

Each open debt item becomes a Conductor task ("retire `needToyCheck` for T2"). Work-dependencies — not domain/epistemic dependencies — live here: "the `needToyCheck` task can't be claimed until the `needMap` task is accepted" is a Conductor fact; "T2's validity logically depends on T1" is domain content and stays out of Conductor entirely, exactly per the decomposition-versus-coordination boundary from writeup6.md. Pizza Bot's agent discovers and claims these tasks through Conductor's existing MCP interface — no new integration surface needed on Conductor's side.

**Phase 3 — Solvent touches only the genuinely consequential moments.**

Debt retirement, promotion, faithfulness review — none of this goes near Solvent. Solvent enters only for the small set of real-world-effecting actions this program actually has: authorizing GPU-hours for a toy-model run, publishing a write-up externally, merging a shared artifact other agents depend on. Reuse the already-frozen Solvent baseline from the Reference Loop work rather than standing up a second instance.

**Phase 4 — Run one real claim through the whole stack before building anything else.**

Theorem 1 (linear Madelung ⇒ Fisher structure) is the right first candidate: its debts are already mostly characterized from earlier rounds — the map and invariant are stated, the citation question (Hall-Reginatto) folds cleanly into `needMap`'s retirement, and `needObstruction` is open (nobody has yet tried the $S$-dependent construction that would break it). Push exactly this one idea through capture → clarify → debt → Skill-driven retirement attempts → human faithfulness gate → promotion, end to end, before scaling to anything larger — matching the reference-loop discipline already earned two documents ago: prove the smallest real thing before formalizing further.

**Phase 5 — Let the cascade question answer itself empirically.**

If Phase 4 (or a later claim like T2's diffusive-corner restriction) surfaces a real case where un-promoting one idea should have propagated to something citing it and didn't, that's the evidence needed to decide whether EBP gets a Solvent-style `belief_edge`/cascade extension — earned through an observed gap, not designed in ahead of it.