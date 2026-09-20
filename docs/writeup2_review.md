Bias disclosure first, same as the pivot review: this report narrates the review process I conducted, and a reviewer grading his own influence has incentives in both directions — to declare victory on the lineage and to protect its open findings. Weigh accordingly; the findings below rest on the report's own text against the record.

Verdict up front: **Phase 1 is real, disciplined work, and the PASS is defensible — but only for a claim narrower than the one the document makes, and with the single most important negative test of the entire program still unrun.** The report inherits two unclosed items from the pivot review (pre-registered criteria, the punch list) and both are now demonstrably consequential rather than pedantic. One new structural finding — the Conductor↔Solvent edge in the architecture diagrams — needs immediate clarification because as drawn it contradicts the specification lineage's core discipline.

## Disposition of pivot-review findings

| Prior # | Finding | Status | Notes |
|---|---|---|---|
| F1 | Pre-registered Reference Loop pass criteria | **Not addressed** | "Verdict: PASS" issued with no stated criteria. Becomes C1. |
| F2 | Failure taxonomy + arbiter | Informally applied, unformalized | §16's routing instincts are right ("Conductor and coordination questions, not reasons to move authorization into Conductor"); the classification rubric and owner still don't exist. M-tier. |
| F3 | Contract version pinning | **Partial** | Solvent is pinned — genuinely good. Conductor, executor stub, and the spec version under test are unpinned; F6 (parent spec unseen) also stands. |
| F4 | Enforcement-path sandbox rule | N/A this phase | Fake executor, honestly labeled. Must govern Phase 2 — see routed actions. |
| F5 | BM-IST deferral ownerless | **Improved** | Phase 6 placement is now explicit and framed as a fit-test for the generic architecture. Still no named owner of the deferred risk. |
| B1 | Paralysis framing | **Resolved** | This report tells the honest version — diminishing returns, empiricism — no impossibility claim. |
| B2 | Close the punch list first | **Not closed — and now it bites** | See C2 and M1. |
| B3 | Decision record discipline | Not applied to the freeze restoration | "Forensic process identified…" — by whom, when, recorded where? |

## What genuinely landed

Credit where the report earns it: §10's invocation ≠ external effect distinction is exactly the discipline the lineage demanded; §12's not-yet-proven list is honest and unusually complete; §15's refusal of the five-primitive runtime is the Growth Gate's first real application, correctly returning NO; §5's persistence-non-uniformity lesson is correct interoperability doctrine; the identity string embedding `run_id` shows the right totality instinct (run-specific parameters in identity); and the frozen-baseline method is sound experimental design *in principle* — my attack on it is below, and it's about what got frozen, not the freezing.

## Critical

**C1 — The PASS is unfalsifiable as issued.** What observation would have produced FAIL? Nothing in the document answers that. The pivot review required pre-registered pass conditions before running precisely to prevent post-hoc verdicts; the experiment ran, and the verdict was assigned against unstated criteria. The two negative tests were *designed to pass* — a denial test that doesn't deny would be a bug in the test, not a discovery. This doesn't mean Phase 1 proved nothing; it means the verdict word is doing rhetorical work the evidence hasn't purchased. Minimum remedy: enumerate retroactively the criteria Phase 1 claims to meet, and — non-negotiably — pre-register Phase 2's criteria before a real effect is triggered, including what result would count as architecture-falsifying. An empirical program that grades itself without pre-registration is the same pathology as the specification loop it escaped, one level up.

**C2 — The operation-mismatch test is missing from Phase 1, and the roadmap defers it wrongly.** Check the negative-test inventory against the lineage: wrong *actor* (denied), missing *authorization* (denied). The X≠Y case — authorized for `deploy:...:main:run-1`, execution attempted for `deploy:...:main:run-2` — is **not in the evidence**. This is the single invariant the entire specification apparatus exists to enforce: exact operation binding was the critical finding of the first spec review, survived the class-subsumption and identity-totality attacks, and became the matrix's binding chain. What §10 actually demonstrates is *propagation* ("consistent operation identity" on the happy path, through one implementation passing one string) — propagation is not binding *verification*. The binding property is proven only by a mismatch being rejected. The fix is cheap — an afternoon against the RecordingFunc executor, no external risk — and the roadmap's placement of "wrong operation" in Phase 3, *after* real external effects, is a sequencing error: the mismatch rejection is precisely the safety property a real effect depends on. It must run before Phase 2, not after it.

## High

**H1 — The Conductor↔Solvent "Authorization Protocol" edge contradicts the lineage.** §14 and the §20 diagram route the consequential path Agent → Conductor → *(Authorization Protocol)* → Solvent. The canonical chain in the specification series has always been Agent/client → authority boundary directly, with Conductor explicitly excluded: it "does not classify the semantic payload," does not authorize, and the anti-pattern list polices precisely this absorption. So either (a) the diagram is conceptual shorthand and the implementation has Conductor as a transport-only relay — in which case annotate it and show evidence that Conductor holds no authority-relevant state for the consequential path; or (b) Phase 1's implementation actually routes authorization through Conductor as a participant — in which case that's an undeclared architectural change that should have gone through the Growth Gate. The report must say which. This is not pedantry: if Conductor is in the authorization path, everything the matrix says about Conductor's non-role is now false in the reference implementation, and every future conformance scenario inherits the ambiguity.

**H2 — What was frozen has never been checked against what was specified.** Two related exposures:

- **Conformance delta of the frozen baseline.** Solvent was frozen at `7602699` — presumably predating most of the specification series. Does it implement the DENIED vs. UNKNOWN/UNAVAILABLE distinction? The temporal policy? Declared revocation semantics? The report doesn't ask. If the frozen kernel lacks required semantics, then "Reference Loop adapts to frozen Solvent" doesn't freeze stability — it freezes *non-conformance*, and the rule forbidding Solvent to adapt guarantees the gap is permanent until someone re-cuts the freeze. Required: a conformance-delta assessment of `7602699` against the spec, then either a documented re-freeze at a conformant commit or an explicit accepted-deviations record.
- **The declaration regime is absent from the experiment.** §11 references "the declared execution model," but no capability declaration for the deploy operation appears anywhere in the Phase 1 evidence — no trust basis, no validity model, no idempotency behavior, no effect classification. Phase 1 validated the enforcement plumbing while skipping the declaration lifecycle the matrix spent three versions building. That's acceptable for a stub-executor phase; it is **not** acceptable for Phase 2 — running a real GitHub effect against an undeclared operation is anti-pattern 15 ("using a declaration that is not the currently effective declared version," one worse: using none). The governing declaration must exist and be pinned before the first real effect.

## Medium

- **M1 — Identity comparison semantics (the punch list's H1) is now load-bearing.** The happy path can't distinguish comparison rules — one implementation passes one string. The moment substitution testing begins (Phase 5), executor A comparing canonical bytes and executor B comparing parsed fields will diverge, and the report's "cross-implementation operation identity mismatch" scenario will find it the hard way. Define the comparison rule in the declaration now, while changing it costs a sentence.
- **M2 — Exactly-once is a property of the test, not yet of the system.** A thread-safe counter in a single controlled invocation proves the call happened once *in that run*; it says nothing about redelivery, retry, or duplicate handling. §12/§16.5 acknowledge this honestly — but §10's framing ("much stronger execution-path invariant") oversells slightly. The system-property version requires a duplicate-delivery test; cheap; belongs with C2.
- **M3 — Rejection-cause coverage is 2 of 5.** The spec's fail-closed causes: absent, invalid, stale/expired, mismatched, verification-unavailable. Phase 1 covers absent and (actor-)invalid. Stale and verification-unavailable matter especially given the frozen Solvent's unknown temporal semantics (H2). Add to the pre-Phase-2 battery.
- **M4 — Scope the PASS honestly.** "Survived contact with an actual end-to-end implementation" — with a RecordingFunc executor, three roles are real and the fourth is a stub, and the external boundary is the architecture's entire *point*. §17's phrasing ("complete for the current controlled experiment") is the honest one; the executive summary should match it.
- **M5 — Partial pinning.** Pin the Conductor build, the executor stub, and — critically — the spec version the loop claims to implement. F3/F6 residue; the Workflow Specification v0.3 remains unseen in this conversation, four requests running.

## Low

Pre-register the attribution measurement (both Phase 1 failures were induced, so diagnosis was trivially known; the claim "failure classification is a lookup" is still unmeasured). "Phase 1 complete" in the header vs. §17's qualifier — pick one. The freeze-restoration decision needs its owner/date record. Phase 6's BM-IST deferral needs a named risk owner.

## Routed actions, in order

1. **C2** — run the mismatch test (and M2's duplicate-delivery, M3's stale/unavailable) against the fake executor now. Hours of work; they gate Phase 2.
2. **H1** — clarify the Conductor↔Solvent edge: annotate as transport-only with supporting evidence, or escalate to the Growth Gate.
3. **H2a** — conformance-delta assessment of frozen Solvent; re-freeze or record accepted deviations.
4. **H2b** — produce and pin the capability declaration for the deploy operation before any real effect.
5. **C1** — enumerate Phase 1's criteria retroactively; pre-register Phase 2's (including falsification conditions and the attribution measurement) before running it.

The report's own closing framing is right: the architecture has become falsifiable. But falsifiability only pays when the verdicts are criterion-bound and the central negative case is actually exercised. Run the mismatch test, fix the diagram, audit the freeze — then Phase 2's real effect will be worth the risk it takes.