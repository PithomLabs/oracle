The fixes from the last round mostly landed — and the honest adversarial finding is that the attack didn't disappear, it migrated. v0.1's critical hole was "binding doesn't exist." v0.2 builds binding, and the hole is now "binding exists but is under-defined," which is more dangerous because it *looks* closed. Full disposition, then the fresh attacks.

## Disposition of prior findings

| Prior # | Finding | Status | Notes |
|---|---|---|---|
| C1 | Self-authored fail-closed; no declaration obligation | **Half landed** | §6 declarations now exist and are normative. But see High-2: the classification basis is still self-authored, and the escape hatches survive. |
| C2 | No proposal–execution binding | **Landed, underspecified** | §7.2 is the right mechanism. See Critical-1 for how it reopens. |
| C3 | No fault-injecting conformance actors | **Landed** | §16's TestSolvent/TestExecutor catalogs are genuinely good — "authorize X while execution requests Y" is exactly the right fixture. See High-1 for what the fixture list still can't test. |
| H1 | Denial ≡ unavailable; unbounded revision | **Mostly landed** | §8's taxonomy with "MUST NOT reinterpret unknown as denial" is correct. §9 bounds *tests* only — see Medium-4. |
| H2 | Observability asserted, not obligated | **Landed** | §14/§15 MUST with test-adapter escape. Residuals below. |
| H3 | Agent Skill / Harness undefined | **Landed** | §24/§25 definitions are adequate. |
| M1–M7 | Medium tier | Mixed | M2, M6 fixed (§6/§11, §10.1). M5 largely dissolved by per-boundary operation contracts. M1, M3, M4, M7 persist — see below. |
| R1 | Growth Gate adapter-smuggling criterion | **Declined, fourth request** | §28 unchanged. Still one line. |

## Critical

**C1-v0.2 — The binding check is only as strong as the operation identity, and the identity is deliberately loose.** §7.2 requires authorization bound to "the exact operation," but the identity is a SHOULD ("canonical operation identity **or equivalent stable description**") — and §6 simultaneously permits declarations over "the operation **or operation class**." Two attack chains follow directly:

1. **Class subsumption.** Solvent authorizes class C. The Executor maps operation Y into C by its own internal mapping. X≠Y check passes — Y *is* what was authorized, if membership says so. But membership semantics are defined nowhere: not in the declaration requirements (§6 lists no membership predicate), not in §7.2, not in §16's fixtures. A broad class plus a loose executor-side mapping is a class-shaped bypass of the exact-binding invariant, and the conformance suite's mismatch fixture (authorize X, execute Y, distinct literals) cannot catch it because the failure mode is *legitimate-looking subsumption*.
2. **Identity totality.** The identity must contain "relevant target/action parameters" — relevant to what, decided by whom? Any material parameter omitted from the identity (quantity, target scope, environment, flags) allows proposal and execution to agree on identity while diverging on the omitted field. Binding passes; effect differs.

Fixes, none requiring infrastructure: (a) the operation identity MUST be total over effect-relevant parameters — a parameter either appears in the identity or is declared immaterial; anything else makes it a different operation; (b) if class binding is permitted, membership semantics MUST be part of the §6 declaration and the executor's membership check MUST use the declared predicate; (c) add a class-subsumption negative test to §16 and §23B (authorize class C, execute an operation *outside* C that a buggy mapping would include). Absent these, §7.2 is a lock whose keyhole shape is chosen by the thing being locked out.

## High

**H1-v0.2 — The §23 scenario set silently under-delivers Requirements §16.** Requirements v0.3 §16 made "long-running external work does not block unrelated Conductor progress" a *minimum conformance assertion*. It has no scenario in §23 (A–F cover other assertions) and no fixture to run it with — §16 has no long-running/held TestExecutor mode and no TestConductor at all, though the non-blocking property is a *Conductor-side* observable. Likewise "workflow phases do not become persisted workflow state" is a Requirements minimum with no §23 test and no oracle strategy. A Test Matrix derived from §23 will ship a suite that passes while failing its parent document. Add: a long-running scenario (held executor, unrelated work progresses, completion reconciles), a blocking-Conductor fixture, and a persistence-location inspection test (which needs the declaration hook in High-2).

**H2-v0.2 — Declarations are now the linchpin, and they have no lifecycle.** §6 fixed the missing obligation and immediately became the single point of failure:

1. **No objective definition of "external effect."** Declarations self-classify. An under-declared integration is undetectable *by the suite, by construction* — the suite tests exactly the declared set. The v0.3-review recommendation survives unimplemented: any discovered undeclared effect-capable behavior SHALL be recorded as a conformance-scope finding.
2. **No versioning or re-declaration trigger.** §27 versions the spec; nothing versions declarations or forces re-declaration when an integration adds an operation. Declarations rot, and the suite keeps certifying against last year's surface.
3. **No review hook.** The declaration defines the authorization surface — policy-adjacent by any reasonable reading, notwithstanding §6's "not a policy engine" disclaimer. Solvent or the domain authority has no stated audit right over the artifact that decides what they get asked to authorize.

One paragraph closes all three: declarations are versioned specification artifacts, re-declaration is required on capability change, discovered undeclared effect-capable behavior is a scope finding, and the authority owner MAY audit declarations.

## Medium

- **M-A — Outage-window execution is now legitimate by construction.** §8's UNKNOWN clause ("do not execute *unless an independently valid existing authorization already covers the requested operation*") plus reusable authorization (§6/§11 model 1) means a client that waits for Solvent unavailability and executes from cached authorization is fully conformant. Each step defensible; the aggregate is a cached-authority execution strategy. Require this path to be *declared* in the validity model and *evidenced* in §14 (executed-under-pre-existing-authorization marker), so at least the behavior is visible and countable.
- **M-B — "Deterministic inspection of … relevant ordering" with no order model.** §15 demands ordering evidence across heterogeneous record stores (Conductor API, Solvent records, executor records, external systems) with no common clock. Declare what suffices: happens-before reconstruction via the §14 correlation chain, or logical timestamps on cross-role records. Otherwise every ordering assertion in the Test Matrix is unfoundable.
- **M-C — Conformance obligations bind "a conformance deployment."** Production deployments have no observation obligations, so post-certification drift is invisible — certificate decay. State whether §14/§15 obligations persist in production or conformance is explicitly point-in-time.
- **M-D — Production revision bounds are unowned.** §9 bounds tests and honestly disclaims a production retry counter — good — but then nobody owns the production bound. Name the owner (Solvent policy, Conductor coordination, or explicitly domain policy) even if the answer is "out of scope."
- **M-E — Concurrency still silent.** Two concurrent consequential proposals with interacting effects: interference ownership (Solvent? Domain?) unstated. One paragraph.
- **M-F — Client-internal state vs. invariant 6.** §19.2 blesses Temporal-wrapped clients; such a client persists phase-like state internally. One scoping line: client-internal representations are not workflow-authoritative state, preserving the invariant.
- **M-G — §13's "idempotent retry" reconciliation *is* §12 duplicate delivery.** Cross-reference them so an ambiguous-outcome retry unambiguously follows the declared replay policy.

## Low

- §5.4's SHOULD ("ordinary work … without Solvent") vs. Scenario A's expectation — mark the no-Solvent check advisory-level or promote the SHOULD.
- §19.1 names GPT/Claude as examples — Requirements §11 bans provider names from *normative* definitions; add "illustrative" to be airtight.
- §24: the Agent Skill should be version-pinned to the spec revision it derives from; "examples" rot.
- M3 persists: §18's autonomous mode still admits policy-injected human gates, keeping the mode distinction mechanism-less.
- R1: the Growth Gate still lacks the smuggling criterion ("does the change relocate forbidden functionality into an existing role's component?"). Fourth request.

## Verdict

This is the strongest version of the document — C2, C3, H2, H3, M2, and M6 all landed with real mechanisms, not gestures. But the two load-bearing new mechanisms (declarations, binding) each carry an unresolved definitional hole, and the scenario list quietly drops two Requirements minimums. None of the remaining findings require infrastructure; all are definitional obligations the spec is already philosophically equipped to state. Routed order: (1) identity totality + class-membership semantics (Critical), (2) declaration lifecycle paragraph (High-2), (3) missing scenarios + fixtures (High-1), (4) the Medium tier as one-line edits. At that point the §29 claim — "complete enough to derive" the Test Matrix and Harness — would actually be true rather than aspirational.