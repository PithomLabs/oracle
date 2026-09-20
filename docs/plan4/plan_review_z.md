Consistency note first: this document consolidates two reviews, one of which was mine — so I'm grading a synthesis that incorporates my own findings. I've re-attacked from zero rather than confirming; most of the synthesis survives, but its *central rendering* contains a grain error that will be implemented if not corrected, and three ownership questions the consolidation smoothed over are still open. The two load-bearing items it kept (§11) are the reason this is close to buildable.

## Disposition against the prior review

| Prior item | Status | Notes |
|---|---|---|
| capability_ref replaces `kind` | ✓ Adopted | Correctly the headline correction. |
| Minimal nine-field schema | ✓ Adopted | Two missing fields — see M6. |
| READY computed, not stored | ✓ Adopted | "blocked" used in the formula but never defined — M3. |
| Atomic claim, provenance edges | ✓ Adopted | |
| operation_id derived, not stored | ✓ Adopted | |
| Identity equality + declaration coverage "should not postpone" | ✓ Named | But the *timing* is wrong — see M5. |
| Protocol verb correction (AUTHORIZE/EXECUTE are not Conductor verbs) | ✗ Silently dropped | Not repeated, not repudiated. |
| Deferred list (stale claims, priorities, discovery filters) | ✗ Silently dropped | Third time the disposition-log request applies. |

## Critical

**C1 — The task-level ordinary/consequential fork is a category error; classification lives at the operation level.** The doc's §7 renders two paths — "ORDINARY: claim → work → submit" and "CONSEQUENTIAL: claim → Solvent → Executor" — and §12 draws the fork *below Conductor*, branching on `capability_ref`. Three problems compound:

1. **The fork is drawn as Conductor's routing decision.** Conductor stores `capability_ref`; it must never branch workflow on it — that's the policy-router prohibition. The consequence plane is Agent→Solvent directly (the earlier transcript's two-plane diagram had this right; §12 regressed it).
2. **The binary doesn't survive real capabilities.** The most important coding-agent capability is shell — a universal capability that can produce external effects by accident or by instruction. Its declaration cannot honestly say "non-effect." The lineage's answer exists: declarations are per *operation or operation class*, so a shell declaration classifies operation classes (local-file: non-effect; network/external: effect-capable, evidence-checked, fail-closed at the shell boundary). Under that model, a single "ordinary" task routinely *contains* consequential moments — the task doesn't fork; operations do.
3. **Agent self-gating reappears at the wrong granularity.** "The Agent can see the declaration and understand: I need Solvent before execution" is fine as advisory behavior — but as the doc's *mechanism* ("the resulting paths are…"), it re-imports classification-by-agent, one level up. The invariant that survived seven review rounds is: enforcement is at the effect-capable boundary, fail-closed, regardless of what any task or agent believes.

The fix preserves the doc's elegance: **the fork is a projection, not a mechanism.** Tasks carry `capability_ref` as coordination/matching metadata (optionally with the constraint that a task's primary capability should be singular — a useful decomposition discipline worth stating as guidance); classification is always resolved per-operation from the declaration at the boundary. Task-level rendering may exist for human comprehension, labeled as derived. Enforcement never depended on the fork — which is exactly why C1 is a semantics bug, not a security hole. But a design document whose central diagram encodes the wrong grain *will* be implemented wrong.

## High

**H1 — Declaration resolution path is unassigned.** The whole design now hangs on `capability_ref → declaration`, and nobody serves that lookup. The matrix requires retrievable content and harness-observable versions but names no resolver. The options have opposite failure modes: if Conductor stores declaration *content*, it becomes a contract registry (drift, staleness, semantic-ownership creep); if the Agent fetches from integration owners, there's a discovery problem and no pinning. Minimal resolution consistent with the subtraction doctrine: **Conductor stores `capability_ref → declaration_id + version + content_hash`; content is served by the declaration owner; the consumer verifies the hash.** And fail closed at CREATE: a task whose `capability_ref` doesn't resolve to an effective declaration version is not creatable (or not READY) — otherwise the frontier hands out work whose contract doesn't exist. This decision also finishes §11's equality item properly: the declaration defines identity *construction and comparison*; resolution is how both reach every participant.

**H2 — ACCEPT/REJECT have no owner.** The work protocol retains ACCEPT/REJECT, and `verification_ref` points at "what establishes acceptability" — but who decides? If Conductor evaluates, it's a verifier (domain semantics inside infrastructure — prohibited). If the submitting agent self-accepts, then DONE is agent-asserted and every downstream READY task inherits unverified completion. The rule the lineage supports: **Conductor records acceptance decisions; it never generates them.** `verification_ref` names the deciding authority (test suite invoked by the submitting agent *with reported evidence*, human reviewer, domain verifier); ACCEPT arrives through that authority's boundary as a decision record with attribution. Note the chain: READY = deps DONE, so the frontier's integrity is exactly as strong as this unresolved answer.

**H3 — Durable history is one fence away from the shadow ledger.** The doc mandates "entire project state + history" — the user's constraint, correctly adopted — but never fences it against invariant 11 and matrix §17. Two rules prevent the collision: (a) history entries are **coordination facts with attribution** — who reported/decided/created what, when, referencing externally-owned evidence — never authoritative duplicates of authorization, execution, or effect facts; (b) "decided ≠ authorized" — a decision record ("branch B selected," "we will deploy X") is project memory, not a Solvent fact; expect someone to conflate them within weeks. Also say the quiet part: append-only history with a current-state projection is the right shape, and it is *Conductor's coordination record*, not a workflow event store — one sentence in the design, or the anti-ledger invariant and this feature will be argued about forever.

**H4 — Untrusted peers are now part of the threat model, and the doc doesn't notice.** The pitch is arbitrary successive agents joining a shared work graph. That makes `objective` text — agent-authored natural language that future agents will read and act on — an injection surface, and task *creation* a spam/poisoning vector. The security boundary survives (declaration governs classification; Solvent gates effects — an injected "and then deploy" still hits the authority boundary), but the design should state: task text is data, never instruction; every task and decision record carries creator attribution; and creation policy (who may CREATE) is an explicit deployment concern rather than an accident.

## Medium

- **M1 — Diagram semantics.** §12's fork-under-Conductor is C1's visual form; redraw with the Agent on both planes (work plane and consequence plane), Conductor observing results.
- **M2 — Lifecycle never enumerated.** "Only Conductor's ordinary lifecycle" needs its closed set (OPEN, ACTIVE, BLOCKED, SUBMITTED, DONE, REJECTED is sufficient) and its fence: no authorization/execution states, ever.
- **M3 — Blocked is undefined and duplicated.** READY's "not blocked" needs a blocker-reference semantic (task or external event), and `dependencies` (§2 field) vs `blocks` (§5 edge) are the same thing twice — keep one.
- **M4 — Project memory should get a Growth Gate record.** Durable project state/history is a genuinely new Conductor capability set — justified by the "agent can't resume" limitation, which the Reference Loop can evidence. Locking it by consolidation without a gate record violates the project's own decision discipline (the same finding as the pivot's B3).
- **M5 — §11's timing is off.** Identity construction/comparison semantics aren't a "before substitution testing" item anymore; the task schema makes them a *day-one* requirement — arbitrary agents construct operations from the declaration. Same for declaration coverage: unresolved `capability_ref` breaks CREATE, not Phase 5.
- **M6 — Attribution fields missing.** `created_by/created_at` on tasks and decisions (H4 depends on them); consider `claimed_at` in the claim record.
- **M7 — Stale claim × authorization validity is now near-term.** A long coding-agent task holds a claim for hours; intents have validity windows. The claim-expiry semantic was deferred to Phase 3, but the *interaction* (submit-after-intent-expiry → which rejection class?) belongs in Phase 1.5's negative battery.

## What survives attack

The two headline corrections are right and load-bearing: `kind`-field removal correctly relocates classification to the declaration, and operation-identity derivation from declaration+inputs+run_id kills the dual-definition drift. The Beads rejections are all correctly reasoned — formulas/molecules and Dolt sync would each have violated the subtraction doctrine, and declining hash IDs despite their elegance is the right call (semantic IDs are load-bearing in the correlation chain). "Conductor stores and coordinates; Agent supplies intelligence; declarations classify; Solvent authorizes; Executor effects" is a correct final rule — it's only §7's rendering and §12's diagram that encode it wrong.

## Routed order

1. **C1** — re-grain §7/§12: operation-level classification, task-level fork as derived projection, Conductor never routes on it.
2. **H1** — declaration resolution: reference+version+hash in Conductor, owner-served content, fail-closed at CREATE.
3. **H2** — acceptance ownership rule; verification_ref names the decider.
4. **H3** — the two history fences (coordination facts only; decided ≠ authorized).
5. **H4/M6** — attribution fields + task-text-is-data rule.
6. Mediums as line edits; restore the dropped protocol-verb correction and the deferred list — with a disposition log this time, so "silently dropped" stops being a recurring finding.

One process observation to close: the doc's own razor — "subtract until nothing more can be removed" — is the Growth Gate's inverse form, and it's being applied more faithfully here than anywhere else in the project. Apply it once more and it will cut §7 down to one sentence: *the boundary classifies the operation; the task merely references the capability.* That sentence is the design.