This breaks the impasse, and it's worth saying plainly why: **the impasse existed because one mechanism — task classification — was being asked to answer three different questions.** Plan|work splits them, and each lands on the mechanism that actually owns it:

```text
Is this operation effect-capable?     → declaration   (mechanical, at the boundary, fail-closed)
Should we do this body of work?       → human         (plan grain, before work starts)
Is this exact operation allowed now?  → Solvent       (operation grain, at execution time)
```

Two gates, two grains, two times — so they don't compete, they compose. The human gate is coarse and early; Solvent is fine and late. A plan approved Monday can be unsafe Friday — that's why both exist. And the ordinary/consequential fork dies rather than gets repaired, because nothing routes on it anymore: work mode flows claim → work → submit, and the only branch is at the boundary.

Three conditions before I'd call it locked:

**1. Plan approval is consent, not authority.** This must be written down in one sentence, or the project's most-policed invariant ("human presence ≠ second authority path") collapses: an agent that treats an approved plan as permission has made the human an out-of-band authority engine and Solvent into ceremony. Approved plan = license to enter work mode. Every consequential operation still passes Solvent. Conductor records the approval as a decision record — attribution, plan version — which is coordination history, not authority fact.

**2. Plans must enumerate their consequential operations as identity templates.** Otherwise the misbinding attack returns one level up: human approves prose, agent executes an operation that passes Solvent but was never contemplated. The lineage already solved this exact shape between authorization and execution via exact binding — apply it between plan and operation. Practical consequence, and it's a feature: **humans approve an operation inventory, not a narrative.** Revisions are new plan versions (the existing V1/V2 discipline).

**3. Conductor stores plans and records decisions; it never gates on them.** No plan-aware routing, no `READY_FOR_AUTHORIZATION`-style plan states — that's workflow-phase-as-persisted-state again. The entire Conductor delta is: plan registry (id, version, operation templates, decision records) and `plan_id` on tasks. That's it.

One honest limitation to record rather than hide: **plan↔operation conformance is detectable, not preventable, in v1.** With plan_id + operation_id on every consequential outcome, drift is visible to audit — but nothing automatically blocks an out-of-plan operation, because preventing it would mean teaching Solvent about plans, and Solvent stays frozen at `7602699`. Record it as an accepted limitation with an owner; if it bites in practice, plan-scoped Solvent policy is a Growth Gate question, not a silent patch.

What this kills — the subtraction ledger:

- The task-level fork and the `kind` field — deleted, dissolving the category error rather than patching it
- The classification *role* of `capability_ref` — it survives, if at all, as optional interface/matching metadata
- The autonomous/semi-autonomous mode matrix — replaced by two mechanically distinct modes, which also closes the mode-incoherence finding carried since the workflow spec
- The unowned revision bound — the human is the iteration controller; that question is answered, not deferred
- Acceptance ownership — largely closed by default: the approver is the top-level acceptance authority, delegable per plan
- The per-boundary human-gate catalog — collapsed to one mandatory gate plus always-available boundary interventions
- Most of the declaration's content — policy migrates to the plan gate; declarations shrink to mechanical classification, operation identity, and validity. Note they don't die: "anything with an external effect" still needs a definition of external effect, and the per-capability declaration is that definition. Keep the machinery, lose the weight.

What survives untouched: Solvent (frozen, Phase 1 evidence stands), Executor, SOR semantics, the evidence chain, the Reference Loop results. The mismatch test and stale/unavailable battery still gate Phase 2; add one scenario — *operation outside approved plan, detected via the evidence chain*.

And the multi-agent vision gets stronger for free: with `plan_id` on tasks, an arbitrary Agent X joining later can only claim work inside approved scope — containment by construction. The residual question ("who may CREATE tasks outside any plan?") becomes a named deployment policy instead of an accident.

One-sentence form: **the human approves what will be attempted, Solvent approves what is being attempted, the boundary enforces that only the second ever matters at execution time — and the task graph is just bookkeeping underneath.** Lock those three conditions and this is the cleanest version of the architecture the project has produced.