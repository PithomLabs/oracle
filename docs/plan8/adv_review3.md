## deepseek

# Adversarial Review: PIVOT_POC_IMPLEMENTATION_PLAN.md (Final)

## Short verdict

The plan has resolved nearly every issue from the previous three rounds. The Phase 0 gate is coherent, the idempotency design is sound, the domain pack contract is concrete, the refusal tests are executable, and the simplification ledger has verification methods. This is the first version that reads as implementation-ready in spirit.

But the final revision introduces two **new contradictions** that were not present before, and it leaves four **residual ambiguities** that will surface during implementation. None is fatal. All are fixable.

---

## 1. New contradictions introduced by this revision

### 1.1 Schema ownership contradiction

The plan states (§1.5):

> Solvent owns all epistemic/authority schema migrations. ARGUS consumes Solvent's migration via `solventmigrations.Apply(db)`.
> ARGUS owns only work-specific and idempotency schema.

And then:

> `-- Proposed retirement (ARGUS-owned, new migration)`
> `ALTER TABLE belief ADD COLUMN proposed_retirement JSONB;`

This is a contradiction. `belief` is a Solvent-owned table. ARGUS is altering it. The plan claims migration ownership is cleanly separated but then ALTERs Solvent's table.

**Consequences:**

- Solvent's next migration cannot safely alter `belief` without knowing whether ARGUS has already added `proposed_retirement`.
- The migration order (Solvent first, then ARGUS) is now load-bearing but not stated.
- If Solvent adds a migration that recreates `belief`, ARGUS's column is dropped.

**Concrete fix:** Either (a) move `proposed_retirement` to an ARGUS-owned sidecar table (`belief_proposed_retirement(belief_id, payload)`) with a FK to `belief`, or (b) explicitly document that ARGUS's ALTER is a deliberate violation of ownership, with a stability commitment from Solvent.

### 1.2 Idempotency status column with one allowed value

§1.4a schema:

```sql
status STRING NOT NULL DEFAULT 'completed'
    CHECK (status IN ('completed'))
```

A CHECK constraint that allows only one value is a no-op. The column carries no information. Either:

- The status column is meant to have multiple states (`pending`, `committed`, `failed`) but the schema was reduced without updating the plan, or
- The column is vestigial from the previous design and should be removed.

**Concrete fix:** Either remove the column or specify the state machine.

---

## 2. Residual ambiguities that will surface during implementation

### 2.1 "Return cached result" has no cache

§1.4a step 6:

> Write `submission_idempotency` row LAST with status `COMPLETED` — its presence certifies completion.
> Duplicate submission of completed packet → idempotency row exists, return cached result.

Where is the "cached result" stored? The `submission_idempotency` table has `content_hash`, `scenario_id`, `packet_id`, `status`, `created_at`. It does not store the response. To return the same result, the application must either:

- Reconstruct it from the persisted entities (read belief, evidence, task by packet_id).
- Store the response JSON in a new column.

The plan says "return cached result" as if a cache exists. It does not.

**Concrete fix:** Specify either (a) reconstruct response from persisted entities keyed by packet_id, or (b) add a `response JSONB` column to `submission_idempotency`.

### 2.2 Entity-level idempotency requires deterministic entity IDs

§1.4a:

> Persist beliefs — each belief write uses `INSERT ... ON CONFLICT (scenario_id, claim_hash) DO NOTHING`

For this to work, `claim_hash` must be deterministic from packet content. And the belief ID must be deterministic too — otherwise, on retry, the belief is written with a *different* ID, and the ON CONFLICT fires on `(scenario_id, claim_hash)`, silently discarding the new row while leaving the old row with a different ID. Downstream references to the "belief ID" will be broken.

The plan does not say whether belief IDs are:
- Content-derived (deterministic hash).
- Random UUIDs (nondeterministic on retry).
- Provided by the packet (agent-controlled, may collide).

The same question applies to evidence, edges, and tasks.

**Concrete fix:** State that entity IDs are deterministic from content (e.g., UUIDv5 of scenario + claim_hash). Otherwise, the ON CONFLICT logic returns stale IDs on retry.

### 2.3 Verifier artifact path is still ambiguous

§2.3 says:

> Artifact persistence: artifacts are stored as evidence rows with `provenance_class = 'reproducible_artifact'`. The CLI (`argus verify`) does NOT write to DB directly. It produces a JSON artifact that enters via `submit_packet` (agent) or `argus verify --submit` (application layer, same validation path).

Two paths stated, not fully reconciled:

- **Path A:** Agent runs `argus verify`, gets JSON, includes it in a packet, submits packet via MCP. The application layer re-verifies? Or trusts the artifact?
- **Path B:** `argus verify --submit` calls the application layer directly, which validates and persists.

If Path B exists, why does Path A exist? What prevents the agent from submitting a fabricated artifact via Path A?

The plan says "Same library underneath" but does not say the application layer **re-runs** the verifier. If it does not re-run, then Path A trusts the artifact's claimed verifier_id — which the agent controls.

**Concrete fix:** State that the application layer re-runs the verifier during packet validation, ignoring the artifact's claimed results, or state that artifact trust is out of scope for the POC and the CLI `--submit` path is the only accepted path.

### 2.4 Adversarial-agent independence is weaker than claimed

§3.1:

> Adversarial-agent independence is procedural in this demo (separate MCP session, no shared state).

"No shared state" is false. Both agents read the same CockroachDB. Both write to the same DB. The adversarial agent can read the work agent's beliefs, evidence, and edges. It cannot read the work agent's conversation, but it can read the work agent's output. That is the intended design (the adversarial agent is supposed to see the work), but calling it "no shared state" misdescribes it.

**Concrete fix:** Say "no shared conversational memory; system state is shared by design."

---

## 3. Weaker claims than they should be

### 3.1 "Credential sets: 1" is questionable

The simplification ledger says:

| Credential sets | multiple | 1 | single env var |

But the POC actually has:

- One operator token (env var) for the UI.
- One DB credential (or a connection string) for CockroachDB.
- (Possibly) one MCP configuration secret if any.

That is at least two credential sets. "Single env var" is a specific implementation choice, not a count of credential sets. The ledger row is imprecise.

**Concrete fix:** Either count credential *types* (operator, DB) rather than "sets," or clarify what "1 credential set" means.

### 3.2 "task dev is non-destructive" conflicts with `cockroach demo`

§3.6 Taskfile:

```yaml
dev:
  cmds:
    - cockroach demo --no-example-database --listen-addr :26260 &
```

`cockroach demo` starts an in-memory cluster. Data is lost when the process exits. If `task dev` uses `cockroach demo`, then:

- Every `task dev` starts with a fresh, empty DB.
- The schema "already applied" claim is false — the DB is empty.
- `task dev` is destructive in effect (state does not persist).

**Concrete fix:** Either use `cockroach start-single-node --insecure` with a persistent data directory, or rename `task dev` to `task ephemeral` and adjust the description.

### 3.3 `argus reset` semantics still undefined

The plan mentions `argus reset` as a subcommand but does not say:

- What it drops (all tables? only ARGUS tables? only CRDB system tables?).
- Whether it re-applies migrations.
- Whether it prompts for confirmation.

**Concrete fix:** Specify `argus reset` as: drop ARGUS-owned tables, re-apply Solvent and ARGUS migrations, log the reset.

---

## 4. The one remaining architectural concern

### 4.1 The work subsystem's load-bearing status is weaker than claimed

The plan preserves all Conductor semantics and says:

> Nothing is discarded. Conductor's domain model is small and entirely load-bearing.

But the work-semantics inventory then says:

| Dependencies | Minimal | Simplify — keep type, defer propagation |
| Blocked/ready | Minimal | Simplify — basic blocking only |

If dependency propagation is deferred and blocked/ready is "basic only," then what do dependencies actually do? If nothing, why keep the `Dependency` type? The plan needs to state what "basic blocking" means concretely — e.g., "a task with an unfinished blocker cannot be claimed" — and confirm that this behavior is used by the POC.

**Concrete fix:** Add a sentence to §0.2: "Dependencies enforce a single rule: a task with an unfinished blocker cannot transition to `active`. Propagation (transitive blocking) is deferred."

If even this rule is not used, then `Dependency` should be removed.

---

## 5. Small but concrete issues

### 5.1 `proposed_retirement` is JSONB, but the plan doesn't say what schema

§1.5:

> When agent submits packet with proposed retirement, store as metadata on the belief.

What is the shape of the JSON? Is it a list of debt items? A map from debt item to proposed evidence? Without a schema, the Trust UI cannot render it reliably.

### 5.2 Discharge's replay protection conflicts with re-discharge after retraction

§0.8 says:

> `Discharge(ctx, scenario, beliefID, obligationKey, instrumentRef, dischargedBy)` — higher-level method that records attributed discharge with replay protection (UNIQUE constraint on belief+obligation+instrument).

If a belief is retracted and then re-entered (REOPEN), and the human discharges the same obligation again, the UNIQUE constraint blocks the second discharge. Is that intended? If so, how does the human re-authorize the new lineage?

### 5.3 Boundary matrix doesn't cover `ui`

The boundary matrix lists `epistemic`, `work`, `verifier`, `mcp`, `application` but not `ui`. Does `ui` have import restrictions? Can `ui` import `epistemic` directly, or must it go through `application`?

### 5.4 The Phase 0 gate's "GO granted" is self-authorized

§0.10 says:

> The gate fired REVISE on the Solvent module boundary. Resolved by the 2 bounded additions within the approved change budget. No material scope expansion. No thesis invariant weakened. **GO granted.**

Who granted it? The plan-writer? The plan-writer is the same entity that produced the Solvent change budget. There is no external reviewer. For the gate to be a real gate, either:

- An external reviewer must sign off, or
- The plan acknowledges that the gate is self-assessed and that Phase 1 may proceed without external approval.

Neither is stated.

### 5.5 `Discharge` vs `RetireDebt` boundary enforcement

§0.8 says:

> In ARGUS, `RetireDebt` is internal to the application layer. `Discharge` is the human-facing API. The boundary enforcement test verifies MCP cannot reach `RetireDebt`.

But `RetireDebt` is a method on the Solvent kernel, imported into `internal/epistemic`. The MCP package imports `internal/application`, which imports `internal/epistemic`. So MCP has a transitive path to `RetireDebt` unless `internal/epistemic` exposes only `Discharge` publicly and keeps `RetireDebt` unexported.

The plan does not say whether `RetireDebt` is exported or unexported in ARGUS's re-export.

**Concrete fix:** State that ARGUS's `internal/epistemic` package exposes only `Discharge` (and other human-facing methods), keeping kernel-level mutation methods unexported.

---

## 6. What the plan gets right

To be fair:

- The Phase 0 gate is now coherent: it fired REVISE, the issue was resolved within a stated budget, and GO was granted.
- The Solvent change budget is explicit: two bounded additions.
- Idempotency is entity-level, not insert-first.
- The domain pack contract has field schemas and a worked BM-IST example.
- The refusal test suite has entry points, inputs, expected outcomes, and assertions.
- The dead-end linkage has a concrete SQL lookup.
- REOPEN has a schema representation.
- The verifier trust binding has a canonical hash specification.
- Auth has a concrete flow (token → session → CSRF → server principal).
- The simplification ledger has verification methods.
- Test disposal is properly classified (including NOT APPLICABLE for out-of-scope repos).
- The Solvent rejection fallback is stated (copy ~200 lines).

This is a serious, well-considered plan.

---

## 7. Bottom line

The final plan is close to implementation-ready. Two contradictions must be resolved before Phase 1 begins:

1. **Schema ownership contradiction** — ARGUS ALTERs Solvent's `belief` table despite the stated ownership rule. Either move to a sidecar table or document the exception.
2. **Idempotency status column is a no-op** — a one-value CHECK constraint carries no information. Remove or specify.

Four residual ambiguities will surface during implementation:

3. **Cached result mechanism** for idempotent packet re-submission is unspecified.
4. **Entity IDs must be deterministic** for ON CONFLICT to return stable IDs on retry.
5. **Verifier artifact path** — Path A vs Path B not reconciled; re-run vs trust decision unclear.
6. **Work subsystem's dependency semantics** — "basic blocking only" needs a concrete rule.

Plus smaller issues: `argus reset` semantics, `proposed_retirement` JSON schema, `Discharge` replay protection after retraction, `ui` boundary, Phase 0 self-authorization, `RetireDebt` export visibility, "task dev is non-destructive" vs `cockroach demo`, "credential sets: 1" imprecision.

None of these is fatal. All are fixable in a short addendum. The plan has moved from "a research program" to "an implementation plan." The next step is to resolve these and begin Phase 0.

**One observation:** the plan's own gate discipline says a GO requires no material scope expansion. The Solvent change budget is now 2 changes. If a third change becomes necessary during Phase 1, the gate should fire again — and the plan does not say what happens then. It should.


## Z

# ADVERSARIAL REVIEW — PIVOT PLAN, "FINAL" REVISION
## The header's claim is now true — remaining defects are implementation semantics, not architecture. But the entity-level idempotency fix quietly violated the migration-ownership model, and the narrowed gate can no longer catch that violation.

**Executive verdict:**

$$
\boxed{\text{B1–B7 landed at high fidelity; the plan has converged and should be frozen as a plan. One new defect cluster remains — ARGUS's new entity-level writes reach into Solvent-owned schema (constraints on belief/evidence/belief\_edge; ALTER TABLE belief), breaking the 2-change budget and the ownership rule — and §0.9's REVISE triggers were narrowed in exactly the way that can no longer detect this. Fix by decision, not by another plan revision: this is the last round where plan-editing beats a checklist.}}
$$

---

## 1. B-FIX FIDELITY AUDIT

| Fix | Status | Evidence / residue |
|---|---|---|
| B1 entity-level idempotency + #16b | ✅ landed, **new defect F1** | §1.4a rewritten correctly (row-last-as-certificate, ON CONFLICT resume, #16b present) — but see §2 |
| B2 edge confirmation | ✅ landed (option b) | No schema change; refusal test #5 rewritten to the enforceable contract; behavior covered by acceptance #10; residual: guarantee is test-enforced, not structure-enforced — acceptable, document |
| B3 migration export FS, ledger = 2 | ✅ landed, **blemish F2** | §1.2b correct in mechanism; fallback text is Change-1's fallback misplaced under Change-2 |
| B4 tolerance + input binding | ⚠️ partial | `Tolerance` in schema ✅, re-run rule in acceptance #3 ✅, rule text ties tolerance to needFaithfulnessReview ✅; **input-spec binding not implemented** — `VerifierSpec` still {ID, MinVersion} only; agent still chooses inputs |
| B5 Discharge vs RetireDebt | ⚠️ landed, **new ambiguity F3** | §0.8 defines both — but the layering sentence creates a possible unattributed bypass path |
| B6 proposed_retirement anchored | ✅ landed, **new defect F4** | Column + Debts-view rendering present — but persistence violates ownership (see §2) |
| B7 misc | ✅ mostly landed | Retry wrapper fixed (closure-wide, max 5) ✅; no-agent-HTTP committed, vectors executable ✅; verifier error≠refuted + timeout ✅; artifact persistence path + `--submit` ✅; application boundary row ✅; scenario containment in compile ✅; truncation flag ✅; go/packages isolation test ✅; governance_ref FK ✅; dev/fresh split ✅ — with **F8** |

**Gate integrity:** §0.10 now states the gate fired REVISE and was resolved within budget before GO. ✅ Closed — the round-2 finding is properly discharged.

---

## 2. NEW FINDING CLUSTER: THE SCHEMA-OWNERSHIP COLLISION (the round's headline)

The two best fixes of this revision — entity-level idempotency and anchored proposed-retirement — both require writing to **Solvent-owned schema**, while §1.5 and amendment 12 declare "Solvent owns all epistemic/authority schema migrations" and the Solvent change budget is exactly 2.

- **F1 (P0):** `ON CONFLICT (scenario_id, claim_hash)` / `(scenario_id, content_sha256)` / `(scenario_id, from_id, to_id, kind)` each require a matching unique index on `belief`, `evidence`, `belief_edge` — Solvent tables. The recon does not record these constraints as existing. If absent, the first integration test fails loudly ("no unique constraint matching ON CONFLICT specification") — good failure mode, unbudgeted fix. Three resolutions, decide in Phase 0: (a) verify the constraints exist (then record them as recon evidence); (b) budget **Solvent change #3** (add unique indexes — note the semantic decision this carries: claim-hash dedup collapses two same-hash beliefs in one scenario; almost certainly the desired idempotency semantic, but it is Solvent's data-model decision, not ARGUS's); (c) side-table approach ARGUS owns. Also the `submission_idempotency.status` CHECK allowing only `'completed'` is a vestigial column from the abandoned design — drop it or mark reserved.
- **F4 (P0):** `ALTER TABLE belief ADD COLUMN proposed_retirement JSONB` is listed under "ARGUS-owned, new migration" while altering a Solvent-owned table — a direct self-contradiction of §1.5's own ownership rule. Clean fix: **ARGUS-owned side table** `proposed_retirement(belief_id, items JSONB, proposed_by, created_at)` — no Solvent change, ownership clean, and it carries attribution (currently absent: whose proposal?).
- **F9 (P1, governance):** §0.9's REVISE triggers were narrowed from the prior revision (which fired on "schema migration or transaction boundaries require a design change") to only material-scope/invariant/redesign items. That narrowing is what allows F1/F4 to slip through the gate ungoverned. Add one trigger: **"any Solvent change beyond the declared 2-item budget."** Without it, the budget is advisory — and the budget is the only thing keeping the pivot from re-growing a fourth integration surface.

These three are one finding in three costumes: **the collapse moved the ownership boundary into the schema, and the plan's governance hasn't followed it there yet.**

---

## 3. OTHER NEW/RESIDUAL FINDINGS

- **F3 (P1):** §0.8's layering is self-contradictory: `RetireDebt` is "called by the application layer's discharge path" while `Discharge` "is the method the human flow calls" — both kernel methods, one with replay protection (UNIQUE on belief+obligation+instrument), one without. If any ARGUS path calls `RetireDebt` directly, that path **bypasses replay protection and attribution**. One sentence fixes it: *ARGUS calls `Discharge` exclusively; `RetireDebt` is kernel-internal; depguard forbids its use outside the epistemic wrapper.* Until written, the definition documents the trap it was meant to remove.
- **F2 (P2):** §1.2b's fallback ("copy ~200 lines of view logic") is Change-1's fallback pasted under Change-2. Change-2's actual fallback: vendor Solvent's `.sql` files with a recorded drift-check test. One-line correction.
- **F5 (P2):** refusal-suite vectors need scoping. #7 (body-principal ignored) doesn't specify belief state or instrument — `Discharge` requires `(beliefID, obligationKey, instrumentRef)`; as written the test could 409 for unrelated reasons and the attribution assertion never executes. #8's audit query needs `WHERE belief_id = X` or cross-test pollution breaks it. #9 needs a named fault-injection mechanism (test hook in persist.go, or a forced constraint violation) — "simulate failure" is not yet executable.
- **F6 (P2):** `argus verify --submit` enters via "the same validation path" — good — but whose principal? One sentence: operator principal (human-class) or a distinct `verifier/system` class; it must also pass VerifierSpec (it does, per "same path").
- **F7 (P2):** composition root unstated in the boundary matrix — `cmd/argus/main.go` must import both `internal/application` and `domain-pack/bmist` to inject the pack. Declare main.go the exempt composition root, or depguard will flag the one wiring file that must exist.
- **F8 (P2):** `task dev` "assumes schema already applied" — this breaks the clone→`task dev` promise (§15) with a first-run trap. Better: since all migrations are `CREATE TABLE IF NOT EXISTS`, `task dev` should run `solventmigrations.Apply` + ARGUS migrations idempotently — migration application ≠ reset; nothing destructive needed.
- **F12 (P2, security):** `GET /login` validates the token — a secret in a query string lands in access logs and browser history. Use POST + body. Trivial, but it's on the adjudication path.
- **F10 (P2, third notice):** unused authority surface (`authority.Service`, `executor.Registry`) is still KEEP-imported with zero demo or acceptance coverage. The 17-step narrative has no execution path; the POC exercises none of it. DEFER the import. This is the last surviving repeat finding from round 1 — retire it.
- **Cosmetics:** §2.4's flow still omits the edge-persist step that §1.4a includes; §0.2 still says "Nothing is discarded" beside a table that defers dependency propagation; `:belief_id` named-param style vs `$N` elsewhere; SemVer-minor + "unstable, may change" is a mild tension in §1.2a (harmless).

---

## 4. WHAT IS NOW GENUINELY CLOSED (convergence evidence)

Worth recording, because four rounds of review should end with the scorecard: refusal suite is executable with entry/input/expected/assertion columns; the simplification ledger has verification methods per row and honest counts (2 instances + OpenCode); the gate fired and resolved visibly; artifact trust binding exists with canonical-hash documentation and a no-DB-direct-write CLI rule; REOPEN lineage has a real column and FK; dead-end lookup is specified SQL; retry semantics are closure-wide with bounded attempts; truncation and scenario containment survived into the flows. **Every structural finding from rounds 1–3 is implemented or explicitly deferred.** The defect class has completed its migration: architecture → enforcement → mechanism → semantics. That trajectory is the pipeline working.

---

## 5. RECOMMENDATION: FREEZE THE PLAN, CONVERT TO A PHASE-0/1 CHECKLIST

Another plan revision has negative expected value — the remaining items are decisions and test-spec details, not document structure. Convert to:

**Phase-0 decision list (blocking, one sitting):**
1. F1: verify the three unique constraints exist in Solvent's schema (recon evidence) → else choose change-#3 vs side-table, and set the claim-hash-dedup semantic with Solvent.
2. F4: proposed_retirement → ARGUS-owned side table with attribution.
3. F3: one-sentence Discharge-exclusivity rule + depguard.
4. F9: add the budget-overrun REVISE trigger.
5. F10: DEFER the authority/executor import.

**Implementation-order rule (TDD-shaped, cheap insurance):** implement the refusal suite and #16b **before** `persist.go` and `human.go`. Every defect found in this plan's last two rounds (N1's deadlock, F1's constraint absence) is exactly the class those tests catch — writing the tests first makes the remaining risk self-detecting rather than review-detecting.

**Small-fix batch (no review needed):** F2, F5 scoping, F6 principal sentence, F7 exemption, F8 idempotent dev migrations, F12 POST login, cosmetics.

---

## 6. CLOSE

The plan's own header claims the remaining defects are implementation semantics, not architectural problems — **this review confirms the claim, with one exception worth naming: the schema-ownership collision (F1/F4/F9) is architectural in the exact sense this program cares about — an authority boundary (who owns the epistemic schema) that the implementation crossed without the governance noticing.** It is also the only finding here that the narrowed gate could not have caught, which is why F9's trigger matters more than its size suggests.

For the record, the pipeline retrospective: four rounds, ~40 findings total, defect class migrating from *is the pivot right?* → *is enforcement real?* → *do the mechanisms contradict themselves?* → *do the details hold?* — each round's findings smaller, later, and closer to code. That is what a converging review process looks like. The plan is done. The checklist above is the last artifact the document owes; after that, the repository is the source of truth, exactly as the prompt originally demanded — and the first thing the code should prove is the refusal suite, because a system named for what it refuses should ship with its refusals tested before its flows.