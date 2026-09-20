# PHASE 1 — REFERENCE LOOP REVIEW

## 1. What the Reference Loop actually proved
- Reproducible happy-path execution across Agent → Conductor → Solvent → Executor → Agent.
- The same operation identity is carried through task creation, authorization, execution, and continuation.
- Solvent is the sole authority boundary: `AuthorizeAction` must precede `ExecuteAction`, and the execution path requires a live `intent_id`.
- The evidence collector is separated from the critical path and runs after execution.
- `RecordingFunc`-style execution is explicitly labeled as non-external-effect proof, not `EFFECT_CONFIRMED`.

## 2. What it did not prove
- Real external effect occurrence. GitHub execution is simulated by a fake executor registered in-process.
- Authoritative effect confirmation from an external SOR. There is no verified GitHub workflow run result bound to the operation.
- Concurrency, cancellation, failure recovery, or multi-agent contention.
- Long-running state transitions or intent lifecycle beyond `live` → `executed`.
- Adversarial or negative scenarios beyond a single simulated happy path.

## 3. Verified boundary invariants
- `CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION` holds in implementation.
- `AUTHORIZE ≠ EXECUTE` is enforced by Solvent REST: execution requires `intent_id` created by authorization.
- Agent does not perform authority evaluation.
- Conductor does not perform authority evaluation.
- Executor does not perform authority evaluation.

## 4. Evidence ownership
- Conductor owns work lifecycle and activity records.
- Solvent owns authority, intent, and audit records.
- Executor owns execution attempt evidence via Solvent audit; it does not own effect confirmation.
- Agent owns correlation, operation identity, and continuation decisions.
- Harness owns post-run evidence aggregation; it does not create authoritative operational state.

## 5. Operation identity findings
- Deterministic operation identity format: `deploy:repo:workflow:ref:run_id`.
- `run_id` is immutable per run.
- `scenario_id`, `task_id`, and `intent_id` are correctly correlated.
- No silent drift observed between authorization and execution in the happy path.

## 6. Agent ↔ Conductor observations
- Interface used: MCP stdio via `tools/call` against the built `conductor` binary.
- Agent discovers, claims, updates, submits, and continues work through MCP tools.
- Agent does not access Conductor storage directly.
- Minimum future Agent Skill behavior needed: work discovery, claim, progress reporting, submission, continuation.

## 7. Conductor ↔ Solvent observations
- Conductor does not call Solvent directly in the current loop.
- Solvent interaction is Agent-mediated via REST.
- Conductor's governance provider wiring exists but is not exercised by the happy-path test.

## 8. Solvent ↔ Executor observations
- Solvent resolves the executor by action name from an internal registry.
- The current test registers a fake `github_trigger_workflow` executor in `startSolventServer`.
- Frozen Solvent REST API does not return `intent_id` from `AuthorizeAction`; the harness uses a DB fallback query.

## 9. Executor ↔ External SOR observations
- No real external effect is triggered in Phase 1.
- Effect confirmation is simulated and labeled as non-external-effect proof.

## 10. Genuine defects
- IMPLEMENTATION_DEFECT — minor: `GetLiveIntent` DB fallback exists because the frozen Solvent REST API omits `intent_id` from the authorization response. This is observable but does not break the happy path.

## 11. Specification defects
- None identified. The current contract adequately describes the happy path.

## 12. New security properties
- None newly discovered in Phase 1.

## 13. Protocol concepts empirically justified
- Agent ↔ Conductor: MCP tool-call boundary for work coordination.
- Solvent as sole authority boundary with atomic authorize+intent creation.
- Executor as post-authorization, non-authority execution layer.
- Operation identity including `run_id` for end-to-end correlation.

## 14. Protocol concepts NOT yet justified
- Cross-role intent/work/authorization/effect/outcome primitive.
- Formalized failure, cancellation, and revocation protocols.
- Adversarial and negative-path contracts.
- BM-IST validation mechanics.

## 15. Recommended next experiment
- Add one focused negative test: mismatched actor_id authorization denial, and one execution-without-authorization denial.
- Keep Solvent frozen.
- Do not expand architecture.

---

## FINAL VERDICT: PASS WITH LIMITATIONS

- Solvent HEAD: `7602699`
- Solvent modified during this phase: NO
- Happy path reproducible: YES
- Authorization boundary proven: YES
- Execution boundary proven: YES
- Operation identity proven: YES
- Evidence ownership proven: YES
- Agent decomposition responsibility preserved: YES
- Conductor coordination responsibility preserved: YES
- New protocol primitive justified: NO
- New security property discovered: NO
- Next recommended experiment: Negative-path authorization/execution denial tests before any adversarial suite or protocol formalization.
