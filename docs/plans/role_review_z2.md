This revision is the most responsive yet — C1 landed via the two-row split plus the Boundary Rule, the oracle problem got the documented rejection I asked for instead of a silent carry, and the duplicate-numbering defect is gone. But strict-reading the document in its own defense changes the severity map: the attack that remains is no longer a licensed bypass — it's an undefined foundation class, a self-subordinating precedence clause, and one criterion violation created *by the fix for the last one*. One integrity caveat first: §1/§25 cite Workflow Specification **v0.3** as normative parent. The last spec I reviewed self-labeled v0.2 (twice, with different content). Either the increment happened off-session — resolving the versioning finding — or this matrix projects a parent whose current text I haven't seen. Parent-level findings below are provisional on that; send the v0.3 text if you want them confirmed.

## Disposition of prior findings

| Prior # | Finding | Status | Notes |
|---|---|---|---|
| C1 | Ordinary Work licensing authority-free effects | **Landed** | Two-row split, Boundary Rule ("any operation capable of producing an external effect is consequential regardless of phase or client"), Q1/Q5 added to §24. Strict reading: the non-effect row *requires* its receiver be declared — so undeclared tools can't legally serve ordinary work. Residue: see H1. |
| H1 | Own-criterion violation via dual-owner cells | **Landed, incompletely** | §8's owner/enforcement split is the right mechanism, applied to §5 and §8 — but §13's owner column still carries slash-duals in ≥5 cells. See M1. |
| H2 | Missing scenario evidence rows | **Landed** | §14 adds intervention, declaration, rejection, long-running progress, reconciliation rows. Residual: long-running progress owner is a dual ("Executor + Conductor"). |
| H3-prior | Rejection lost distinct outcome status | **Landed** | §12 six-category enum with bidirectional MUST-NOT-reinterpret rules; §13/§14 rows added. |
| M1 | Cross-document precedence | **Landed, overbroad** | See H2. |
| M2 | Effect-detection oracle (4th carry) | **Documented rejection** | §5's oracle-limitation sentence + §23.12 completeness-as-deployment-responsibility. Acceptable resolution; residue at M4. |
| M3 | Intra-participant boundaries | **Half-landed, created a contradiction** | §6 created — but see H3. |
| M4 | §22 substitution escape clause | **Not landed, 2nd carry** | "Where semantics actually differ" still has no adjudicator; point it at §16 triggers. |
| M5 | Long-running coordination evidence | **Landed** | §14 row exists. |
| Carried M-tier | Production observation persistence; revision-bound owner; concurrency; client-internal state; cached-auth marker | **Silently carried again** | Second consecutive silent carry for most. See M3. |

## High

**H1 — The non-effect declaration class is referenced everywhere and defined nowhere.** The entire ordinary path now rests on "Declared non-effect tool/service" — but §16 mandates declarations *only* for effect-capable integrations. A non-effect declaration has no minimum content, no owner, no version, and no change trigger: when a tool *gains* an effect-capable operation, §16's invalidation clause doesn't apply (it never declared), so nothing obligates re-declaration — the transition from non-effect to effect-capable is invisible to the conformance lifecycle. This is the same attack shape as last round's C1, one level down: not a licensed bypass (the row's own cells forbid undeclared receivers), but an undefined artifact class that the rule depends on. Fix: extend §16's obligations to all tool/integration declarations, with the non-effect variant naming owner, version, scope, and inheriting the change-notification duty.

**H2 — The precedence clause subordinates the matrix's own normative contributions.** §1/§25: on any difference, the Workflow Specification controls. Read strictly, that's a bug: the matrix's six-outcome enum "differs" from the parent's five (it adds *rejected before effect*); its UNKNOWN temporal policy, transport-metadata exclusions, and authority-precedence arithmetic (§11) are all refinements the parent doesn't contain. Under an unqualified precedence rule, every such refinement is subordinate — and a conformance dispute can cite the parent to void the matrix's tighter rule. The clause needs the missing distinction: the matrix may **refine** (subdivide, tighten, make explicit) parent semantics; it may not **contradict** them. One sentence. Related: §1 claims the matrix "is normative for role ownership and boundary separation" *and* is subordinate on conflict — workable, but state that "normative" means "normative within the projection scope."

**H3 — Formulate is declared a boundary (§5, §9) and a non-boundary (§6) in the same document.** §6 says Formulate is an internal workflow phase, "not external participant boundaries" — yet it remains a §5 boundary row and a §9 intervention row. Interpretation was handled correctly (removed from §5, listed in §6); Formulate got both treatments. Separately, Next Work is absent from §9 while §5 declares human intervention for it — the dedicated intervention matrix is incomplete against the standing "every boundary declares" rule. Both are one-line fixes, but a projection artifact whose Test Matrix will cite rows by number cannot ship a row whose existence two of its own sections dispute.

**H4 — The §8 refactor fixed criterion 1 and broke criterion 3.** Last round: dual-owner cells violated "every responsibility has one authoritative owner." The fix consolidated §5's columns into "Semantic/authority owner" + "Enforcement" — and in doing so **dropped the effect-owner tracking** that §23.3 requires ("every consequential boundary identifies the authority owner **and effect owner**"). Consequential rows now carry authority owner only (Consequential Proposal: "Solvent for authority"; Authorization: "Solvent"; Execution Eligibility: "Solvent"), and Execution/Result lost their "Solvent remains authority owner" / "Solvent for prior authority fact" annotations from v0.1. By the document's own criterion 3, at least four consequential boundaries now fail. Either restore explicit authority/effect-owner tags on consequential rows or reword §23.3 to match the new columns — but the criterion and the table must agree, and right now the criterion is the stricter of the two.

## Medium

- **M1 — §13 didn't get the §8 treatment.** Owner column still reads "Conductor/integration," "Solvent/integration," "Solvent/integration boundary," "Effect-capable integration / Executor boundary," "Executor/integration." Resolve each to its single semantic owner (e.g., Authorization-unknown → Solvent; the integration merely reports unavailability) and move the second party to an enforcement note. The criterion the refactor exists to satisfy is still violated in the failure table.
- **M2 — Consequential Proposal vs. Ordinary-effect-capable rows overlap without a stated relationship.** Same initiator, same owner, near-identical purpose text. Presumably the latter is the WORK-phase realization that *inherits* the former's arc — say so ("MUST route through the Consequential Proposal and Authorization boundaries"), or the Test Matrix will test them as independent paths.
- **M3 — Silent carries are now a process finding.** Concurrency/interference ownership, client-internal state, production observation persistence, revision-bound owner, cached-authorization evidence marker: fourth-to-fifth appearances, zero dispositions, while other findings get meticulous fixes. The oracle item shows this team documents rejections well *when it engages* — so add a disposition log (finding → adopted/rejected/deferred + one line). Without it, the Test Matrix authors can't tell dropped-by-intent from dropped-by-omission, and neither can reviewers.
- **M4 — §23.12's "treated as an explicit integration/deployment responsibility" names no owner.** You documented the oracle limitation honestly; finish the move by naming who holds the completeness responsibility (integration owner, with domain/Solvent audit right per the earlier recommendation), or the criterion is a framing verb, not an obligation.
- **M5 — Executor vs. system-of-record (provisional pending parent v0.3).** The matrix's "subject to" is the right projection, but if the parent still says "or," the disjunction survives upstream of the precedence rule and precedence can't fix it. Split the concern explicitly: SOR owns *whether the effect occurred*; Executor owns *attempt and report*.
- **M6 — §10's operator-controlled indefinite wait needs evidence.** An operator hold with no recorded decision is "wait" with a human name attached. One line: operator waits are interventions and generate §14 intervention evidence.

## §24 self-posed questions, answered against this document

| Q | Answer |
|---|---|
| 1. Ordinary tool produces external effect outside boundary? | Not lawfully — receiver must be declared non-effect (H1 gap is the declaration class, not the rule). |
| 2. Authorize X, execute Y? | No — §7.3/§7.4 chain is closed. |
| 3. Same identity definition at all three points? | Yes — §7.3, with declared immaterial fields. |
| 4. Integration bypass Solvent? | Not within declared set; undeclared = post-hoc finding only (documented limitation). |
| 5. Ordinary execution of effect-capable capability? | No by Boundary Rule — enforcement rests on declarations (H1). |
| 6. UNKNOWN as denial / indefinite undeclared wait? | No — §10; operator-wait evidence residual (M6). |
| 7. Ambiguous as success? | No — §12/§13. |
| 8. Rejection confused with failure? | No — §12, bidirectional rules. |
| 9. Human shadow authority? | No. |
| 10. Cancellation mutating authority? | No. |
| 11. Declaration drift invisible? | No — §16; non-effect class excepted (H1). |
| 12. Harness passing while bypassing enforcement? | No — §17. |
| 13. Hide undeclared effect-capable op? | Undetectable ex ante by documented design; §23.12 assigns responsibility without an owner (M4). |
| 14. Rename absorption? | Growth Gate criterion; gate-time only. |

Twelve of fourteen hold; the two residual "yeses" share the same root (declaration completeness) and are now *documented* rather than silent — that's the difference between a hole and an accepted risk, and only the missing owner assignment (M4) keeps it in the latter category.

## Verdict

The matrix is converging: the bypass classes from every prior round have real closures, and the two open "yeses" are declared limitations rather than latent defects. What remains is consolidation work — H1 (non-effect declaration lifecycle) is the only finding that guards against a live bypass path, and its fix is an extension of machinery that already exists. Routed order: (1) H2's refine-vs-contradict sentence — without it, every other matrix normativity claim is arguably void; (2) H1's §16 extension; (3) H3's Formulate/Next Work consistency pass; (4) H4 or M1 — pick one direction for the owner columns and apply it to §5, §13, and §23 together so criteria and tables agree in one edit; (5) the disposition log (M3) before the Conformance Test Matrix starts — that document will inherit every silent carry silently. And send the Workflow Specification v0.3 text: two findings here are provisional on a parent this conversation has never seen.