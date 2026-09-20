The matrix lands most of the routed fixes — happens-before ordering, declaration lifecycle, enforcement-path-preserving sandboxes, trust-basis declarations, identity-function consistency — but it commits two errors of the kind this document is *uniquely* positioned to make: it fails its own acceptance criterion on roughly ten table cells, and its Ordinary Work row answers the document's own adversarial question #3 with "yes." One reference-integrity note first: §1 cites "Workflow Specification **v0.3**" — the document I last reviewed self-labeled v0.2 and was cited for a versioning violation. Either the spec was incremented (resolving that finding) or this cites a version that doesn't exist. Confirm which; everything below assumes the former.

## Disposition of carried findings

| Prior # | Finding | Status | Notes |
|---|---|---|---|
| H1-prior | Spec versioning self-violation | **Likely resolved** | Matrix cites v0.3; verify the increment actually happened. |
| H3-prior | Sandbox bypass → certification theater | **Landed** | §12's "MUST NOT bypass authorization verification / binding / fail-closed" is exactly the fix. |
| H2-prior | Declaration content retrievability | **Landed** | §11: owner, version, effective date, retrievable content. Residual: no authority-owner audit *right* over declarations — carried 4th time. |
| H4-prior | Undefined trust assumption | **Landed** | §5.4/§11 require declared trust basis. Nudge: the phrasing "authenticate the evidence *to* the authority owner" is directionally ambiguous — declare whether evidence reaches the Executor via authority-controlled channel vs. client-supplied artifact. |
| M3-prior | Identity-function consistency | **Landed** | §5.3 requires the same declared identity definition at proposal/authorization/execution. |
| M-B-prior | Ordering model | **Landed** | §9: causal/happens-before, not timestamps. Residual: state that HB edges derive from the correlation chain (unlinked records are unordered) or the harness still can't compute it. |
| M4-prior | G/H fixture support | **Not landed** | §15 says the harness "may provide fixtures" but doesn't enumerate the long-running Executor mode or Conductor-side non-blocking observation. Now the Test Matrix's problem. |
| M6-prior | Ambiguity-impossibility justification | **Not landed** | §11 requires reconciliation *mechanism* but accepts a declared "no ambiguous outcomes" without justification. |
| M-C/M-D/M-E/M-F/M10 | Conformance persistence, revision-bound owner, concurrency, client-internal state, cached-auth evidence marker | **Not landed** | All absent. Fourth carry for the concurrency item specifically — either adopt these or record an explicit documented rejection; silent re-carrying is drift. |

## Critical

**C1 — The Ordinary Work row licenses authority-free external effects.** §4, row 4: initiator Agent → receiver "**Tools / external services**," authority owner "**N/A**," effect owner "**Tool/service-specific**." Read against the architecture's own rules, this is the bypass the entire declaration/binding apparatus was built to close: a path that reaches external services, can produce effects, and has *no authorization boundary*. The agent-advisory classification was neutralized in v0.3 precisely because the *integration declaration* — not the caller — defines effect-capability; this row reintroduces caller-side routing, now blessed by the normative ownership reference. The document's §19 asks, "Can a client call ordinary execution when the operation is actually effect-capable?" — §4 answers *yes*, and supplies no authority owner for the result. A Test Matrix derived from this row will not test ordinary-path effects, and Scenario A ("no Solvent involvement") will certify the hole. The fix is one rule: *tools and external services used in ordinary work MUST be covered by an integration declaration stating they are not effect-capable; any effect-capable capability reached from ordinary work inherits the consequential boundary and its declaration obligations.* Note the deeper carried problem this exposes (M2 below): the series has *still* never defined "external effect" objectively, so "undeclared effect-capable behavior" — §16's own anti-pattern — is undetectable in principle by the harness. Four requests now.

## High

**H1 — The matrix fails its own §18.1 in at least ten cells.** Acceptance criterion 1: "Every responsibility has one authoritative owner." The tables are saturated with split ownership: §7 has "Solvent + integration declaration," "Solvent + integration contract," "Executor/integration"; §8 has "Conductor/integration," "Solvent/integration," "Integration/Solvent boundary," "Executor/integration" (twice); §4 has "Executor/external system" twice. Some duals encode real ambiguity that matters: for mid-flight revocation, *who owns the behavior* — Solvent (revokes) or the integration (checks)? Both act; ownership must be single, with the other party's role stated as "enforces declared model" or similar. A column fix (add "verified by / enforced at" separate from "authoritative owner") resolves all ten without losing information. As written, the one artifact whose entire purpose is disambiguation ships ambiguity as its format.

**H2 — §9's evidence matrix can't support the tests the scenarios demand.** Scenario C requires intervention evidence (intervention point, human action, resulting owned state); §9 has no human-intervention row. Scenario F requires a reconciliation determination record; absent. §11 makes declaration version observable; §9 has no declaration-version fact. §7's non-blocking requirement is a Conductor-side observable with no evidence row. And §4 defines a Next Work boundary that §6's intervention matrix never declares — violating the standing rule that every defined boundary declares intervention capability or states none exists. Each is a one-row fix; together they mean the Test Matrix would have pass conditions with no specified evidence source.

**H3 — Fail-closed rejection lost its distinct outcome status.** Workflow spec §17 distinguished *execution rejection* (evidence missing/invalid/mismatched — **no effect attempted**) from *execution failure* (attempted, didn't complete). That distinction is load-bearing: it's the difference between outcome categories "not attempted" and "failed," and it's what §22.2's negative tests actually produce. §8 here has no rejection row — rejection causes are scattered across authorization-invalid/stale and operation-mismatch, both phrased as evidence problems rather than the outcome fact "rejection occurred, no effect." Add the row: *Execution rejection (fail-closed) | Executor/integration | no external effect attempted or produced | must not be reinterpreted as execution failure or authorization denial.*

## Medium

- **M1 — Cross-document normative conflict, live instance.** §2 makes the Executor authoritative for what occurred "*subject to* the external system of record." The workflow spec says "Executor **or** external system of record" — a disjunction. The matrix picked a side (SOR wins — operationally the right call for reconciliation) without the spec amending. Two normative documents now disagree, and neither states a precedence rule. Add "on conflict, the Workflow Specification controls" or amend the spec.
- **M2 — Effect-detection oracle, fourth carry.** §16's anti-patterns ("integration inventing undeclared effect-capable operations," "Executor inventing class semantics") and §19's Q2 are unenforceable without an objective definition of external effect or an equivalent discovery mechanism. Everything downstream depends on declarations being complete; nothing tests completeness. This is the deepest open item in the series.
- **M3 — Intra-participant "boundaries."** Formulate (Agent → Agent/client context) and Interpretation (Agent → Agent/domain context) aren't boundary crossings — both endpoints are the same participant. Keeping them as boundary rows dilutes what "boundary" means for a document the Test Matrix will cite row-by-row. Mark them internal phases or annotate.
- **M4 — §17's escape clause.** "Declaration changes only where semantics actually differ" — decided by whom, against what? Point it at §11's re-conformance trigger: differences in declared semantics = different boundary = re-conformance.
- **M5 — Evidence for coordination-during-long-running.** Scenario G's core assertion ("unrelated work progressed") has no named observable anywhere. One §9 row.

## Answering the matrix's own ten questions (§19)

| Q | Answer against this document |
|---|---|
| 1. Authorize X, execute Y? | **No** within declared paths (§5.3 chain is sound). **Yes** via the ordinary path — no binding is required there. |
| 2. Integration bypass Solvent? | Not within the declared set; **undetectable outside it** — no effect-detection oracle (M2). |
| 3. Ordinary execution of an effect-capable operation? | **Yes — the §4 Ordinary Work row licenses it** (C1). |
| 4. UNKNOWN silently treated as DENIED? | No; §8 row + Scenario F test exist. |
| 5. Ambiguous silently treated as success? | No; §8 row exists. |
| 6. Human shadow authority? | No; §6 limits column is consistent. |
| 7. Cancellation mutating Solvent authority? | No; §16 anti-pattern + §7 ownership. |
| 8. Declaration drift invalidating conformance invisibly? | No; §11 invalidation-until-re-conformance landed. |
| 9. Harness passing while bypassing enforcement? | No; §12 landed. |
| 10. Rename-absorption? | Growth Gate criterion landed; note it's enforced only at gate time — no standing review mechanism, acceptable for now. |

Two of ten self-posed questions are answered "yes" by the document's own tables. That ratio is the review.

## Low

Map the effect-owner value "Pending" to outcome category "not attempted" explicitly; add proposal *version* identity to the Consequential Proposal evidence cell (the V1/V2 discipline otherwise has no evidence home); §18's criteria lack the evidence/owner discipline the Growth Gate now has; reconciliation appears in §7 and §11 but not as a §4 boundary — fine if intentional, annotate.

## Verdict

This is the most responsive artifact in the series — six prior findings landed cleanly, several with better mechanisms than requested. But C1 is disqualifying for a normative ownership reference: it doesn't merely miss the ordinary-work hole, it *codifies* it, and the Test Matrix derived from §4 will inherit the bypass as an untested surface. H1 matters almost as much — a matrix that violates its own first acceptance criterion hands every future dispute a citation. Routed order: (1) C1's one-rule fix; (2) H1's owner-column refactor; (3) H2/H3 as row additions; (4) M1's precedence line; (5) M2 — the effect-detection question now has four carries and should get a documented decision either way before the Test Matrix is drafted, because that document cannot define its own scope without it.