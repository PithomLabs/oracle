This revision lands nearly everything from the last two rounds — totality, class-membership determinism, declaration lifecycle, scenarios G/H, the two-layer conformance program, the mode redefinition, and finally the Growth Gate smuggling criterion on its fifth request. But the honest adversarial result is that the attack didn't stop; it migrated from the architecture to the *certification process itself* — and the document now violates its own policies in two places, which matters more here than it would in any other spec, because this document's entire authority rests on procedural discipline. Findings:

## Disposition of prior findings

| Prior # | Finding | Status | Notes |
|---|---|---|---|
| C1-v0.2 | Identity looseness / class subsumption | **Landed** | §7.2 totality rule with the include-or-declare-immaterial disjunction is exactly right; §7.0's "membership explicitly declared and deterministically evaluated; Executor MUST NOT invent class semantics" closes the subsumption bypass. Residual at M3/M-H below. |
| H1-v0.2 | Missing long-running + persistence scenarios | **Half landed** | Scenarios G and H exist and Scenario F's explicit UNKNOWN≠DENIED check is good. But see M4 — the §16 fixture catalog cannot execute them. |
| H2-v0.2 | Declaration lifecycle | **Mostly landed** | §7.0 + §28: versioned, re-conformance-triggering, undeclared behavior is a finding. Residuals at H2 below. |
| H3-v0.2 | Agent Skill / Harness definitions | **Landed** | Adequate in §24/§25. |
| M-A | Cached-authorization under UNKNOWN | **Partial** | Age policy now mandatory (§8, §7.3). Evidence marker for executed-under-pre-existing-authorization still absent (M10). |
| M-B | Ordering model | **Not landed** | §15 still demands "relevant ordering" with no clock/happens-before basis. Second request pending. |
| M-C | Production drift | **Not landed** | §22.2 helps at integration level; point-in-time vs. persistent conformance still undecided. |
| M-D | Production revision bounds | **Not landed** | §9 unchanged; owner still unnamed. |
| M-E / M-F / M-G | Concurrency; client-internal state; §13↔§12 cross-ref | **Not landed** | All still silent. |
| M3-prior | Mode incoherence | **Landed** | §18's redefinition (autonomous = zero boundary human-action required; semi = at least one) is mechanism-based. Clean. |
| R1 | Growth Gate smuggling criterion | **Landed** | §28's final paragraph. Fifth request; accepted. |

## High

**H1 — The document violates its own versioning policy.** §27: "Normative semantic changes require a specification version increment." This revision, still self-labeled v0.2, adds normative semantic content: §18's mode definitions changed meaning, §7.2 introduced the totality obligation, §7.0/§28 added binding declaration requirements, §22 created a two-layer conformance program, §29 created a new gate. Any of these is a semantic change under the document's own rule. A specification whose change-discipline fails at authoring time has no standing to demand it of implementations. Increment to v0.3 or mark explicitly as a pre-publication working draft. (Related: "changing the major contract" in §27 is ambiguous under an all-0.x scheme — define what increments.)

**H2 — The declaration is the conformance oracle, but only its version is required observable.** §28: "A declaration version MUST be observable to the conformance harness." Version ≠ content. The declared-capability match test (intervention behavior, §23C), the §22.2 integration negatives, and Scenario H all need to read *what the declaration says* to have an oracle at all. One line: declaration content MUST be retrievable by the harness and by the authority owner. Which surfaces the second residual: the authorization surface is defined by these declarations, and Solvent/domain still have no stated audit right over them — the harness can see them, the authority that depends on them cannot.

**H3 — The §22.2 sandbox paradox can make integration conformance vacuous.** "Safe, reversible, sandboxed, or dry-run environment" — if the dry-run short-circuits *before* the enforcement point, then all six negative tests (missing/invalid/stale/mismatched/duplicate/long-running) pass without ever exercising the code that must reject in production. A dry-run that always "would have executed" proves nothing about fail-closed behavior. Required: enforcement-path-preserving neutralization — the effect is neutralized at or after the enforcement point, never before it. This is the certification-theater hole: the layer exists to catch fixture-vs-reality divergence, and as written it can be satisfied by a simulation that shares no code with production enforcement.

**H4 — "Directly and verifiably available" is an undefined trust assumption.** §7.2 requires the evidence to "carry, or otherwise make directly and verifiably available, the authorized operation identity to the Executor," while §7.3 permits "Solvent validation callback" and explicitly declines to mandate cryptography. An integration whose verification callback is invoked *by the requesting client* is arguably conformant — and its binding check reduces to ask-and-believe. The spec need not mandate crypto, but it must require the declaration to state the *trust basis* of the verification path (who the Executor obtains the bound identity from, through what integrity boundary), so the strength of the binding check is a declared, inspectable property rather than an accident of wiring. §29.2's "carries or directly exposes" inherits the same looseness.

## Structural integrity (Medium, but prominent)

- **M1 — Duplicate section numbers.** Two sections numbered 22 (Conformance Execution Layers / Scenario Test Structure) and two numbered 28 (Declaration Lifecycle / Growth Gate). Classic insert-without-renumber. In a document about to generate a Test Matrix that cites sections by number, "§28" is already ambiguous. Renumber before anything downstream is derived.
- **M2 — Duplicated normative blocks (merge residue).** §6 and §7.0 both contain declaration-requirement lists, and §6's is now a strict subset of §7.0's — two formulations that will drift. §8 contains two overlapping taxonomy blocks (the new AUTHORIZED/DENIED/UNKNOWN summary *and* the old detailed DENIED/UNKNOWN text). Consolidate to single normative statements.

## Medium

- **M3 — Identity-function consistency.** The mismatch check is only sound if the operation identity at execution time is computed by the *same declared identity function* used at proposal/authorization time. One line closes a divergence-by-implementation bypass.
- **M4 — Scenarios G/H lack fixture support.** §16's TestExecutor has no long-running/controllable-completion mode, and there is no TestConductor — yet Scenario G's non-blocking assertion is a *Conductor-side* observable. Scenario H's "inspect the baseline implementation" capability appears nowhere in §25's harness contents. The scenarios can't be executed as specified.
- **M5 — §29 Readiness Gate has no process apparatus.** Eight conditions, no owner, no evidence format, no decision record — asymmetric with the Growth Gate's own owner/date discipline the document just fixed. And condition 8 is circular until the Test Matrix exists; state the sequence explicitly (derive harness and matrix → demonstrate gate → implement).
- **M6 — Ambiguity-impossibility declarations.** §13 requires a reconciliation path, but an integration can declare no operation ever produces ambiguous outcomes — a crash-consistency claim about external systems that deserves required justification, not silent acceptance. Affects §22.2's "where applicable" hedge.
- **M7–M9 — Carried, third request each:** ordering model (M-B); concurrency/interference ownership (M-E); production persistence of observation obligations (M-C). Plus M10: executed-under-pre-existing-authorization still needs an evidence marker under §8's cached-authorization path; and M-D: production revision-loop bound still has no owner.

## Low

Provider names in §19.1 still unmarked as illustrative; Agent Skill still not version-pinned to the spec revision; §13's "idempotent retry" still uncross-referenced to §12's replay policy; §22.2's "affected integrations" should read "all declared effect-capable integrations."

## Verdict

The architecture is now closed against every attack class raised so far — the binding chain (total identity → declared membership → carried evidence → executor comparison → fail-closed) has no remaining step that reduces to self-attestation *as a mechanism*. What remains attackable is the machinery that will certify it: an oracle the harness can't fully read (H2), an integration-conformance layer that can pass vacuously (H3), a verification-trust assumption that can silently degrade to ask-and-believe (H4), and two self-violations of the document's own procedural rules (H1, M1). The §29 derivation claim — that the Test Matrix, Harness, and scenarios can be derived from this document — is now close to true: close H2–H4, add the G/H fixtures, fix the numbering, and it holds. Routed order: H1 and M1 first (trivial, and integrity prerequisites), then H2, H3, H4, then the mediums as the Test Matrix is drafted.