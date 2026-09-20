The attached reviews converge with my prior review. The pivot is still sound, but **the plan needs one final correction pass before implementation**. The strongest net-valid findings are the idempotency retry deadlock, the Phase 0 gate contradiction, the impossible cross-module `go:embed`, the missing verifier tolerance, and several under-specified enforcement details.  

### Consolidated net-valid fixes

1. **Fix Phase 0 gate semantics.** A bounded, backward-compatible Solvent API/migration-package change should be explicitly within the planned change budget; only material scope expansion or weakened invariants should trigger REVISE. Then state the resulting GO decision clearly. 

2. **Fix idempotency.** The current “insert idempotency row first, leave it on failure” design deadlocks retries. Use entity-level deterministic uniqueness for belief/evidence/task/edge writes, and write the submission record only after successful persistence as a completed-response cache; add a partial-failure retry test.  

3. **Fix migration delivery.** `go:embed` cannot embed another module's files. Solvent should export a small migration `embed.FS` or `ApplyMigrations()` API; ARGUS owns only work/idempotency migrations. 

4. **Keep agent edges non-authoritative without inventing unnecessary edge state.** Replace the contradictory “unconfirmed edge” test with an enforceable invariant: agent-created edges cannot themselves trigger retract/promote/authorize; only human decision paths can. Add schema state only if repository inspection proves it is necessary. 

5. **Complete verifier trust binding.** Add `tolerance` and explicit verifier/input binding, define canonical artifact hashing, and specify where artifacts persist and how the application checks the pack's allowed `VerifierSpec`. Do not add full cryptographic attestation.  

6. **Define PackDefinition concretely.** Keep the data-driven typed structure you chose, but specify the schemas/semantics of `ClaimType`, `EvidenceClass`, `DebtItem`, `RetirementRule`, `Falsifier`, `VerifierSpec`, plus one worked BM-IST example per debt item. 

7. **Clarify authority semantics.** Define `Discharge` vs `RetireDebt`, anchor `proposed_retirement` in persistent state/UI, and ensure packets cannot invoke human transitions. 

8. **Specify the refusal tests.** Each negative test needs entry point, input, expected outcome, and assertion. Add partial-retry and failed-mutation/no-success-audit coverage. 

9. **Close the smaller structural gaps.** Explicitly define retraction→dead-end task lookup, REOPEN lineage storage, scenario containment, projection `TRUNCATED != COMPLETE`, the application import boundary, and objective verification methods for the simplification ledger. 

10. **Do not re-expand scope.** Do not add agent containers, full identity infrastructure, cryptographic attestation, a second domain, or a generalized architecture framework.

## Prompt to coding agent

```text
Revise `docs/plan8/PIVOT_POC_IMPLEMENTATION_PLAN.md` one final time. DO NOT implement code.

Incorporate the following net-valid corrections from the latest adversarial reviews:

1. Fix Phase 0 gate semantics: bounded backward-compatible Solvent changes
   (public projections + migration export) are within the approved change budget;
   REVISE only for material scope expansion or weakened thesis invariants. Re-state
   the final GO decision accordingly.

2. Replace the broken idempotency flow. Do not leave an "attempted" row that
   blocks retry. Use DB-enforced entity-level deterministic uniqueness for
   belief/evidence/edge/task writes, and write submission completion metadata only
   after successful persistence. Add a partial-failure retry acceptance test.

3. Replace impossible cross-module `go:embed` migration handling with a small
   Solvent-owned exported migration FS or `ApplyMigrations()` API. Solvent owns
   epistemic/authority migrations; ARGUS owns only work/idempotency migrations.

4. Resolve edge semantics without unnecessary schema expansion: agent-created
   edges are proposals/evidence only and can never themselves trigger
   retract/promote/authorize. Rewrite the refusal test accordingly unless
   repository inspection proves explicit edge state is required.

5. Complete verifier trust binding: add tolerance, canonical hash rules, verifier
   input binding, artifact storage path, and explicit `VerifierSpec` enforcement.
   Same verifier library for CLI and in-process paths; no agent-selected verifier.

6. Fully specify the typed `PackDefinition` contract, including field semantics
   and one worked BM-IST example for every debt item. Keep the core generic.

7. Define `Discharge` vs `RetireDebt`, anchor `proposed_retirement` in persistent
   state/UI, and ensure agent packets cannot invoke consequential human actions.

8. Turn the refusal suite into executable specifications: entry point, input,
   expected result, assertion. Include agent-forged decision packet, auth/CSRF,
   open-debt promotion, retracted-belief promotion, failed mutation/no success
   audit, and partial-retry cases.

9. Resolve remaining concrete gaps: retraction→dead-end task lookup,
   REOPEN lineage storage, scenario containment, TRUNCATED != COMPLETE,
   application boundary imports, and objective before/after ledger verification.

10. Preserve the four-phase structure and keep the POC lean. Do not introduce
    new services, databases, generalized frameworks, full identity systems,
    cryptographic attestation, agent containers, or a second real domain.

Then overwrite the plan with the reconciled version and update its acceptance
criteria so every stated guarantee has a corresponding test or measurable check.

Deliver only the revised plan and a concise change summary. Do not modify code.
```

The key point is that this should be the **last plan revision**, not the beginning of another architecture-review loop. The reviews now describe a plan that is very close to implementation-ready; the remaining defects are concrete implementation semantics rather than a problem with the pivot itself. 
