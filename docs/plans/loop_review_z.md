Noted first: what you've pasted is the spec itself — `plan2.md`, i.e., the Loop Engineering Requirements v0.2 — not `plan2_review.md`. So this is the adversarial review of the spec that's been queued since your first message. If the review document exists separately, paste it afterward and I'll evaluate it against this baseline.

---

## What the spec gets right (so the attacks below are read in context)

The authorization/execution/coordination/domain-truth separation is sound distributed-systems discipline. The anti-infrastructure stance with a Growth Gate release valve is the correct default — most agentic-workflow specs fail by inventing a runtime. Declaring phases to be vocabulary rather than state (§5) avoids the classic state-machine leak. §9's four-way distinction (authorized ≠ started ≠ completed ≠ accepted) is exactly right. The honesty of §7's "intervention cannot retroactively invalidate an external effect" is rare. §12's non-duplication partitioning of test scope is correct. This is a good document. It is not yet a *conformance-ready* document. The findings:

| # | Finding | Severity | Sections |
|---|---|---|---|
| F1 | "Consequential" is undefined — classification circularity enables authority bypass | **Critical** | §8, §13, §14 |
| F2 | No observable protocol surface — the conformance suite has nothing to read | **Critical** | §1, §12, §13 |
| F3 | Domain is a role in §1 but has no owner in §4 — violating the spec's own acceptance criterion 1 | High | §1, §4, §16 |
| F4 | No correlation identifiers — separated records are unjoinable, so key invariants are unprovable post-hoc | High | §6, §10 |
| F5 | Intervention races and precedence undefined; no intervention operation vocabulary | Medium | §7, §14 |
| F6 | No failure taxonomy; ambiguous-outcome case unaddressed | Medium | §2, §9 |
| F7 | Replay/idempotency/revocation ownership unassigned at the authority↔execution boundary | Medium | §8, §9, §12 |
| F8 | Scenarios A–C lack actors, assertions, and pass conditions | Medium | §14 |
| F9 | The Growth Gate has no process | Medium | §17 |
| F10 | No versioning/compatibility policy or requirement stability marks | Medium | §11 |
| F11 | Statistical adherence tier is a dangling capability | Low | §13 |
| F12 | Naming collisions | Low | §1 |

---

## Critical

**F1 — The consequential-classification circularity.** The entire authority architecture hinges on "consequential," and the term is never defined (§8, §13, §14 all depend on it). Who classifies an action as consequential? If it's the Agent, then a buggy or adversarial Agent routes consequential work through the ordinary path and violates no testable invariant — §13's "consequential proposals cannot bypass the authority boundary" presupposes you already know what's consequential. §8's "Effect-capable execution paths MUST enforce *the applicable* Solvent authorization" gestures at the right fix (enforcement at the effect boundary, not agent-side labeling) but "applicable" is doing undefined work: applicability *is* the classification. The spec needs one added invariant: *any external effect lacking a corresponding authorization is a conformance failure regardless of Agent classification* — i.e., classification is advisory, enforcement is fail-closed at the Executor boundary (default-deny for effect-capable paths, with an explicit allowlist owned by Solvent or domain policy). Related and unstated: the adversary model. Invariants like "Agent proposals do not bypass consequential-action controls" (§10.12) are anti-bypass *goals* — against an accidental misclassification or a deliberate one? The spec should say, because F1 is only an attack in the adversarial reading, and the conformance suite's job differs accordingly.

**F2 — Conformance with no observable surface.** §1 promises observable and testable; §12 commits to testing cross-role handoffs; §13 lists test-sounding assertions ("denied authorization prevents the corresponding Executor action"). But the spec defines zero interaction artifacts: no message or event model, no handshake sequences, no ordering or idempotency semantics, no error vocabulary, no timeouts. An automated conformance test must *observe* something — where does the harness read "authorization granted" distinct from "execution started"? Nothing in the spec requires boundary events, emitted observables, or a defined observation point. Two defensible fixes: (a) specify a minimal abstract interface per role — the operations a conforming Conductor/Solvent/Executor must expose, plus a minimal observable event at each §6 boundary trigger (name, roles, correlation refs), keeping storage out of scope; or (b) state explicitly that v1 conformance is audit/checklist-based (documentation review plus scenario demonstration) and defer automated protocol conformance to a later version. The current text sits ambiguously between the two, which is the worst position — it reads like a test suite commitment but can't produce one. (Related: §6 declares the boundary catalog "non-authoritative" while §13 tests behavior "matches the declared capability" — fine if the catalog is the test oracle, but say so.)

## High

**F3 — Domain is a role without an owner, and the spec fails its own acceptance criterion.** §1 lists Domain as one of six roles. §4's ownership table has five columns — Domain is absent. §6 then attributes domain acceptance to "the domain verifier/human," and §16.1 requires "the role boundaries are unambiguous." By the spec's own criterion, it fails: Domain is simultaneously a first-class role (§1), a non-participant (§4), and a human attribution (§6). Concretely: in conformance tests, who is the counterparty for domain acceptance? A role with no defined interface can't appear in a handoff test. Either add a Domain column to §4 (even if its cell says "external to contract") and define its single boundary interaction, or remove Domain from the role list and reclassify it as a concern.

**F4 — Records that can't be joined can't prove the invariants.** The no-workflow-ledger design is coherent *if* each role's records are individually authoritative and jointly reconstructable: Solvent holds authorization records, Executor holds execution records, Conductor holds coordination records. But nothing requires correlation identifiers to propagate across boundaries. Without a shared join key, "coordination activity ≠ proof of external execution" (§10.3) and "authorization precedes execution" (§8) are unauditable in practice — you cannot demonstrate ordering across record stores that cannot be correlated. This is a one-line normative fix (require request/authorization/execution identifiers propagated at every cross-role boundary) that adds no state and no runtime, and it's a prerequisite for any F2 observation model.

## Medium

**F5 — Intervention semantics have races and no vocabulary.** §7 defines *where* intervention is possible but not *what* an intervention is. Scenario C offers three — redirect, reject, revise — with materially different downstream requirements. Worse, precedence is undefined: human A approves at the authority boundary while human B simultaneously intervenes at the Conductor boundary to cancel — which wins, and does the answer differ pre-execution vs. in-flight? And is cancel a stop-proposing, an abort-in-flight, or a drain-and-stop? A minimal normative operation enum (pause / cancel-before-effect / cancel-in-flight-where-supported / redirect) plus a precedence rule for concurrent authority decisions would close this.

**F6 — The ambiguous-outcome case is the dangerous one, and it's unaddressed.** §2 scopes "failure, denial, and recovery expectations," and §13 separates execution failure from authorization failure. Missing: *partial* failure (authorized, started, failed mid-effect) and *ambiguous* failure (crash between effect start and result — the effect may or may not have occurred). The ambiguous case is precisely where a no-ledger design hurts most: if the Executor's record is lost, no one can reconstruct whether the consequential effect happened, and §7 correctly says no intervention can fix it retroactively. The spec should assign "who is the source of truth for whether an effect occurred" — by invariant decomposition it must be the Executor — and say so in one paragraph.

**F7 — Replay, idempotency, revocation.** §9's distinctions imply asynchronous crossing of the authority→execution boundary, which implies at-least-once delivery risk: one authorization, two executions. Component-internal correctness is properly out of scope (§12), but cross-boundary replay is an *interaction* problem and therefore squarely in scope. Similarly: does an authorization persist indefinitely until executed? Can it be revoked between authorization and execution start, and who checks? Assign ownership (likely: one-time-use authorization references, Executor-side dedup, Solvent-side revocation checked at execution start) or explicitly delegate with a named requirement to a component suite.

**F8 — Scenarios aren't testable yet.** All three §14 scenarios are one-line sketches with no setup, actor typing, assertions, or pass conditions. Actor typing matters given §11: is the Agent in each scenario the real stochastic agent or a scripted client? For deterministic conformance (per §13's own split), it must be scripted — say so. Scenario A's "review" has no acceptance semantics, and its "unless a consequential action is actually proposed" inherits F1 entirely. Each scenario needs: setup, actors, steps, observable assertions, pass condition. Cheap fix, high payoff — §16.9's acceptance hinges on them.

**F9 — The escape hatch is unprocessesed.** The Growth Gate is the single mechanism by which this deliberately rigid spec may evolve, and it has no owner, no convening process, no evidence format, no decision record. A gate this consequential needs: who convenes it, what the POC evidence looks like, where the decision is recorded. Also note the structural risk the Gate guards against incompletely: blanket SHALL-NOTs tend to force adapter accretion *inside* existing components, smuggling the forbidden workflow engine into Conductor one adapter at a time. The Gate's fourth criterion ("does not duplicate another role's authority or state") partially covers this; an explicit "net new code must not reside inside Conductor/Solvent/Executor to evade this spec" criterion would close it.

**F10 — An unchangeable contract with no change policy.** §11 declares the cross-role contract may not change; v0.2 is the second version. That contradiction needs a compatibility policy: are new boundaries additive? Must conforming clients tolerate unrecognized boundary declarations, or fail? Which requirements are stable versus provisional at v0.2? Certifying conformance against a spec where half the requirements might move is premature certification — mark maturity per section.

## Low

**F11 —** §13's statistical-adherence clause is directionally right (it matches the deterministic/stochastic split) but "MAY be measured" with no threshold, consequence, or gating status is a dangling capability. One sentence: *behavioral adherence metrics are advisory and non-gating in suite v1.*

**F12 —** "Solvent" and "Conductor" collide badly with common vocabulary — "solvent" especially, in a project called *Synthesis*, where solvent is a real term of art. If these are inherited BM-IST system names, the spec inherits them and should say so with a terminology-provenance note; if they're choosable, choose again.

**F13 —** §16.9's "existing systems" is never enumerated in the document. Fine if this doc travels with other BM-IST material, but the acceptance criterion references systems the document doesn't identify, even by internal reference.

---

## Verdict

The architecture is right; the specification is not yet testable. F1 and F2 are blocking: without a defined consequential-classification rule and a minimal observable surface, a conformance suite literally cannot be written — §13's assertions have no oracle. F3 and F4 are blocking for credibility: the spec currently fails its own acceptance criterion 1, and its core invariants are unprovable without correlation requirements. The rest are cheap, high-leverage revisions that don't threaten the document's excellent instincts — notably, every fix above adds vocabulary or metadata discipline, not infrastructure, so they're all consistent with the Growth Gate philosophy the spec already champions.

If `plan2_review.md` gets pasted next, I'll evaluate it against this baseline — whether it caught F1/F2, and whether its objections survive contact with the actual text.