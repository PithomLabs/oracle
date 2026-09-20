The latest reviews converge: **the architecture is done; only implementation-semantics corrections remain.** The strongest remaining finding is schema ownership: entity-level idempotency and proposed-retirement metadata must not mutate Solvent-owned tables. 

### Net-valid final corrections

1. **Move `proposed_retirement` to an ARGUS-owned side table**, e.g. `belief_retirement_proposal`, with `belief_id`, proposed debt items, proposer, and timestamp. Do not `ALTER Solvent.belief`.

2. **Verify existing Solvent uniqueness constraints before implementation.** If the required `(scenario_id, claim_hash)`, evidence, and edge uniqueness constraints do not exist, do not silently add them. That becomes a Solvent change-budget decision. 

3. **Make entity IDs deterministic** from scenario + canonical content, so retries return the same identifiers rather than silently discarding a newly generated UUID. 

4. **Remove the meaningless `status='completed'` CHECK** from `submission_idempotency`, or omit the status column entirely. Since the row is written only after successful completion, its existence can itself mean completed. 

5. **Define duplicate-result behavior:** reconstruct the result from persisted entities by `packet_id`, rather than calling it a nonexistent cache. 

6. **Verifier trust:** choose one trusted path. The cleanest POC rule is: `argus verify --submit` goes through the application validation path; packet-supplied artifacts are revalidated/re-run against the pack's allowed verifier and input binding. Do not trust agent-supplied verifier claims. 

7. **Fix developer-state semantics:** `cockroach demo` is ephemeral, so it contradicts “non-destructive `task dev`.” Use a persistent single-node CockroachDB/data directory; `task dev` should apply migrations idempotently, while `task fresh` resets.  

8. **Make `RetireDebt` kernel-internal in ARGUS.** The application calls `Discharge` exclusively for human debt discharge. No ARGUS package should expose a direct `RetireDebt` path. 

9. **Clarify edge rule:** no agent-created edge can cause an authority transition; the existing test should enforce that directly. No edge-status schema expansion unless code proves it necessary.

10. **Define remaining small semantics:** dependency = “task with unfinished blocker cannot become active”; specify `proposed_retirement` JSON shape; mark UI boundary; make `argus reset` explicit/destructive; correct “credential sets” metric; say agents share system state intentionally, but not conversational memory. 

11. **Add one governance rule:** any third Solvent change beyond the declared two-item budget automatically reopens Phase 0 as REVISE. 

12. **Defer unused authority/executor machinery** because the POC does not exercise external execution. 

## Prompt to coding agent

```text
Revise the final PIVOT_POC_IMPLEMENTATION_PLAN.md once more, but DO NOT implement code.

Apply only these final corrections:

1. Move proposed_retirement to an ARGUS-owned side table; do not alter Solvent-owned belief.
2. Verify Solvent entity uniqueness constraints during Phase 0; if absent, reopen the
   Solvent change budget rather than silently adding them.
3. Make entity IDs deterministic so ON CONFLICT retries preserve references.
4. Remove the meaningless one-value idempotency status CHECK; completion is represented
   by the completed submission row.
5. Define duplicate-submission result reconstruction from persisted entities.
6. Make verifier submission trust explicit: same application validation path, allowed
   VerifierSpec + input binding; never trust agent-claimed verifier results.
7. Replace ephemeral cockroach demo with persistent CRDB. `task dev` applies migrations
   idempotently and is non-destructive; `task fresh` is destructive.
8. ARGUS calls Discharge exclusively; RetireDebt is kernel-internal/unreachable from MCP.
9. Keep agent-created edges non-authoritative; no extra edge schema unless code proves needed.
10. Specify dependency blocking, proposed-retirement JSON shape, UI import boundary,
    reset semantics, credential metric, and shared-vs-conversational state wording.
11. Add a rule: any third Solvent change beyond the declared two-item budget reopens
    Phase 0 as REVISE.
12. Defer unused authority/executor machinery not exercised by the POC.

Do not add new services, databases, security platforms, or domain packs.

After these edits, FREEZE THE PLAN. Do not start implementation.
```

After this pass, I would stop revising the plan and move to **Phase 0/1 implementation checklist execution**, as the reviews now recommend. 
