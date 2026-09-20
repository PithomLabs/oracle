## deepseek

# Adversarial Review: PIVOT_POC_IMPLEMENTATION_PLAN.md (Revised)

## Short verdict

The revised plan is materially better than the previous version. It has absorbed the strongest concerns: real Phase 0 gate with stated triggers, Solvent public projection contract, agent packets never performing discharge, verifier trust binding, DB-enforced idempotency, typed PackDefinition, table-by-table Conductor→CRDB migration, and refusal test suite. The remaining issues are narrower and more concrete than before.

But four issues remain unresolved, and one is a **logic error** that invalidates the Phase 0 gate. I will organize the review around those.

---

## 1. The Phase 0 gate contradicts itself

The plan states (§0.8):

> REVISE if: Semantics can be preserved, but only with a bounded Solvent module change (e.g., adding public projection methods).

Then in §0.9, answering question 6:

> Is the Solvent module boundary solvable? YES — narrow public projection contract (GetSnapshot, ExplainSnapshot).

And §0.3:

> Decision: import Solvent as a Go module. Solvent exposes a narrow public projection contract for ARGUS (see §1.2a).

By the plan's own definition, this is **REVISE**, not **GO**. A change to the Solvent repository is required — specifically, `internal/view.GetSnapshot` and `internal/view.ExplainSnapshot` must become public, along with their types (`Snapshot`, `Belief`, `Evidence`, `Intent`, `ExplainResult`, `BeliefExplain`). This is a Solvent module change. The plan classifies it as "YES, solvable" and then declares "GO."

**Two interpretations, both problematic:**

1. **If the plan means "the change is small enough to not count as REVISE,"** then the REVISE trigger is misframed. The trigger should be about scope or risk, not about whether a change is needed. As written, the trigger is auto-triggered by the plan's own decision.

2. **If the plan means to ignore its own trigger because the change is backward-compatible,** then the gate is not a gate — it is a formality that the plan can override at will.

**Concrete fix:** Either (a) reclassify §0.9 as REVISE with a plan to revise Solvent's public API in coordination, or (b) redefine the REVISE trigger to exclude backward-compatible additions. The current framing is incoherent.

---

## 2. Idempotency semantics are broken

§1.4a defines the DB-enforced idempotency flow:

1. Compute content hash
2. INSERT into `submission_idempotency` (fails on duplicate)
3. Persist beliefs, evidence, edges, tasks
4. On any failure after step 2, the idempotency row remains — caller must retry with same packet_id

**The failure mode is not handled.** Consider:

- Step 2 succeeds (idempotency row inserted).
- Step 3 partially succeeds: belief inserted, evidence insert fails.
- Caller retries with the same packet_id.

On retry:
- Step 2 fails with duplicate — the caller knows this packet was already submitted.
- But the caller has no way to resume step 3. The belief is already there; the evidence is not. The caller cannot re-run step 3 safely (would double-insert the belief) and cannot skip it (evidence is missing).

**The plan says:** "Partial failures are explicit, retryable, and idempotent." This is asserted, not constructed. The construction requires one of:

- A state field on the idempotency row (`pending` / `committed` / `partial`) plus per-subsystem operation tracking.
- Per-operation idempotency keys (each belief/evidence/task has its own content hash uniqueness constraint).
- A single transaction that covers all persistence (forbidden by §2.4's separate-transaction model).

None of these is specified. Without one, the idempotency guarantee is false.

**Concrete fix:** Specify the state machine for the idempotency row and the retry semantics. If the plan cannot accept a distributed transaction, then each sub-operation must be independently idempotent (per-operation content hash), and the retry must resume from the last successful sub-operation.

---

## 3. The domain pack contract is still thin

§1.2b defines:

```go
type PackDefinition struct {
    ClaimTypes      []ClaimType
    DebtVocabulary  []DebtItem
    EvidenceClasses []EvidenceClass
    RetirementRules []RetirementRule
    Falsifiers      []Falsifier
    Verifiers       []VerifierSpec
}
```

But the plan does not define:

- **`RetirementRule`'s structure.** Is it a predicate (debt item → boolean)? A mapping (debt item → list of required evidence classes)? A callback function the kernel invokes? Each has very different implications for portability. If rules are callbacks, the pack must provide Go code, and each new domain requires a Go package. If rules are declarative data, the kernel must interpret them, and the plan must specify the DSL.

- **`ClaimType`'s structure.** Is it an enum? A string? Does it carry semantics the core needs to evaluate? If the core evaluates claim types, it must understand them. If not, why are they typed?

- **`EvidenceClass`'s structure.** Same problem. If evidence classes are opaque strings, the core cannot verify evidence against rules. If they carry semantics, the core must interpret them.

- **`Falsifier`'s structure.** What is a falsifier in operational terms? Is it a rule the core evaluates? A document the human reads? A signal the verifier emits?

- **`VerifierSpec`'s structure.** The VerifierSpec is said to select "allowed verifier identities." Is this a whitelist of `(verifier_id, verifier_version)` pairs? Does it include parameters? Does the core enforce it?

- **`Pack.Validate()`'s semantics.** What does it validate? Structural consistency? That every debt item has a retirement rule? That every evidence class is reachable? The plan does not say.

**Why this matters:** The portability claim (§9) rests entirely on the pack contract being rich enough to capture domain-specific semantics without the core needing to know them. If the contract is too thin, domains cannot be expressed; if too thick, the core becomes domain-aware. The current contract is a list of 6 field names with no field schemas. That is not enough to evaluate the claim.

**Concrete fix:** Provide a fully worked example for the BM-IST pack: what are its `ClaimTypes`, `DebtVocabulary`, `EvidenceClasses`, `RetirementRules`, `Falsifiers`, `Verifiers` — and what does `Validate()` check for it?

---

## 4. The refusal test suite is under-specified

§3.2 lists:

```
agent → discharge           REFUSED
agent → promote             REFUSED
open debt → promote         REFUSED
retracted belief → promote  REFUSED
unconfirmed contradiction → retract REFUSED
bad Origin/Host → REFUSED
body-supplied principal → IGNORED
failed mutation → no success audit event
```

This is a list of desired properties, not tests. Each line requires specification:

- **"agent → discharge"** — via what path? MCP doesn't expose discharge. HTTP isn't reachable by agents. Does the test simulate an agent crafting a packet with a `proposed_retirement` field? If so, what is the expected behavior — the proposal is recorded as metadata, not executed? That's not "REFUSED," that's "accepted as a proposal." The test's expected outcome is unclear.

- **"bad Origin/Host → REFUSED"** — against which endpoint? What's the check? Is it applied to read endpoints too, or only consequential ones?

- **"body-supplied principal → IGNORED"** — what does "ignored" mean? The request succeeds but uses the server-side principal? Or the request fails because the body contains an unexpected field? Both are reasonable, but they're different tests.

- **"failed mutation → no success audit event"** — what does a failed mutation look like? A rejected promotion? A mid-transaction failure? The test needs a concrete trigger.

**Concrete fix:** For each line, specify: entry point (MCP / HTTP / direct function call), input, expected outcome, and assertion. Without this, the refusal suite cannot be implemented unambiguously.

---

## 5. Additional unresolved issues

### 5.1 The Solvent public projection contract is broader than acknowledged

§1.2a commits Solvent to exposing:

- `ledger.Service.GetSnapshot(ctx, scenarioID, opts) (*view.Snapshot, error)`
- `ledger.Service.ExplainSnapshot(scenario, scenarioID, snap) (*view.ExplainResult, error)`
- Types: `view.Snapshot`, `view.SnapshotOpts`, `view.Belief`, `view.Evidence`, `view.Intent`, `view.ExplainResult`, `view.BeliefExplain`

These were internal types. Making them public means committing to their stability. If Solvent wants to change the internal representation of `Belief` later, it can't without a breaking change. The plan calls this "narrow" and "backward-compatible." It is narrow, but it is not cost-free — it is a stability commitment.

**Concrete fix:** State whether the public types are intended to be stable (SemVer-protected) or are documented as "unstable, may change without notice." For a POC, the latter is acceptable but should be explicit.

### 5.2 The "single configured operator credential" is undefined

§0.6 shows the flow:

```
single configured operator credential
        ↓
authenticated server-side principal (from env/config, not request body)
        ↓
HttpOnly session cookie / minimal CSRF token
        ↓
consequential browser action
```

But the plan does not say:

- How the credential is generated (static secret? random at startup?).
- How the browser obtains the session cookie (login endpoint? auto-issue on first visit?).
- How the session is validated (HMAC? DB lookup? in-memory map?).
- What happens on logout or session expiry.
- How the credential is rotated.

For a POC, a static secret with an HMAC-signed cookie is fine. But it must be stated. Otherwise, the "auth boundary" is a diagram, not an implementation.

### 5.3 The verifier trust binding needs enforcement specification

§2.3 defines the `VerificationArtifact` schema. But:

- **Where is the artifact stored?** In the evidence table? In a separate verifier_artifacts table? On disk?
- **Who writes it?** The CLI (`argus verify`) or the in-process call? Both? How do they coordinate?
- **How does the human verify the artifact is from an allowed verifier?** The PackDefinition lists `Verifiers`, but the plan does not say the application layer checks the artifact's `verifier_id` against this list at submission time. If it doesn't, the trust binding is decorative.
- **What if `argus verify` produces an artifact whose verifier_id is not in the pack's VerifierSpec?** Is the artifact rejected? Recorded as warning? Silently accepted?

**Concrete fix:** Specify where artifacts are stored, who writes them, and the enforcement path for `VerifierSpec`.

### 5.4 The cross-subsystem dead-end linkage is unspecified

§0.2 preserves dead-end semantics "via governance_ref." But the plan does not say:

- How the work subsystem learns that a belief was retracted.
- How it finds the tasks linked to that belief.
- Whether the linkage is via `governance_ref` on the task, via a separate link table, or via the application layer's in-memory logic.

§2.6 says `work.CancelTask(...) [if retraction causes dead end — own work transaction]`. But the condition "if retraction causes dead end" is exactly what needs specification. How does the code know?

**Concrete fix:** Specify the lookup. For example: "On retraction, the application layer queries `conductor_task WHERE governance_ref = belief_id` and cancels each matching task in its own transaction."

### 5.5 The BM-IST pack's retirement rules are not shown

§1.2b says `RetirementRules` map debt items to evidence classes and rules. But the plan does not show a single BM-IST retirement rule. Without at least one worked rule, the pack contract cannot be validated. The portability argument requires knowing that the pack can express real rules.

**Concrete fix:** Include a table showing the BM-IST pack's `DebtVocabulary`, `EvidenceClasses`, and one `RetirementRule` per debt item.

### 5.6 REOPEN is in the walkthrough but not in the work model

§3.1 step 16: "REOPEN creates new lineage."

§0.2 work-semantics inventory: "Lineage / REOPEN | Yes | Preserve — creates new task linked to old."

But the CRDB schema in §1.5 has no column or table for "reopened_from" or "lineage_parent." `conductor_task` has `governance_ref` but no task lineage.

**Concrete fix:** Add a `reopened_from_task_id STRING` column or a separate `conductor_task_lineage` table, and specify how the new task is linked to the old.

### 5.7 The "GO" decision is not gated on a review

§0.7 defines triggers. §0.9 produces the decision. But the plan does not say:

- Who reviews the Phase 0 output.
- What happens if the reviewer disagrees with the GO decision.
- Whether Phase 1 can proceed without a sign-off.

For a "gate," there must be a checkpoint. If Phase 1 begins immediately after §0.9, the gate is procedural, not substantive.

### 5.8 Test disposal table has misclassified entries

§3.5 and the Test Disposal table:

- "DELETE | Solvent API handler tests (in solvent repo, not imported) | N/A"
- "DELETE | Conductor HTTP/MCP/Web tests (in conductor repo, not imported) | N/A"

These are not deletions. The tests still exist in their respective repositories. The plan's intent is that they are not imported. This should be classified as "NOT APPLICABLE" or "OUT OF SCOPE," not "DELETE."

Minor, but it demonstrates that the disposal table was written without care.

### 5.9 The simplification ledger's "Verified?" column remains empty

§3.4 still has an empty "Verified?" column. The plan says verification happens in Phase 3 but does not say how. Some entries are hard to verify:

- "Binaries: ~5 → 1" — how is this verified? `ls cmd/`?
- "Credential sets: multiple → 1" — what counts as a "credential set"?
- "Internal HTTP calls: many → 0" — how is "many" counted before?

**Concrete fix:** Specify the verification method for each row. Examples:
- Binaries: count of `main.go` files in `cmd/`.
- Ports: count of `net.Listen` calls.
- Internal HTTP calls: `grep -r "http.Client" internal/` returns nothing.

### 5.10 The plan does not address what happens if Solvent's public API change is rejected

The plan assumes Solvent will accept the `GetSnapshot` / `ExplainSnapshot` public API addition. What if the Solvent maintainers decline? Options:

- Fork Solvent (which breaks the "import as module" decision).
- Use a small in-process shim that copies the internal view logic (which duplicates code).
- Change the architecture to avoid needing the projection.

None is discussed. For a POC this may be acceptable, but the plan should say what happens on rejection.

---

## 6. What the plan gets right

To be clear about what's improved:

- The complexity classification table is present and grounded.
- The component disposition table is concrete.
- The Phase 0 gate is attempted with NO-GO/REVISE triggers.
- The Solvent module boundary is identified with a specific proposal.
- The Conductor→CRDB migration has a table-by-table mapping and CRDB retry handling.
- The idempotency mechanism is DB-enforced, not in-memory.
- Agent packets never perform discharge/promote/retract.
- Agent-authored edges are non-authoritative.
- Verifier trust binding is added.
- The refusal test suite is a first-class acceptance criterion.
- The BM-IST import isolation test is renamed honestly.
- Taskfile health polling is included.
- Migration ownership is stated.

This is a serious plan. The remaining issues are narrower than before.

---

## 7. Bottom line

The revised plan is close to implementation-ready. Four issues must be resolved first:

1. **The Phase 0 gate contradicts itself.** §0.9 says GO; §0.8's REVISE trigger fires on the Solvent public API change the plan requires. Fix: either reclassify as REVISE or redefine the trigger.

2. **Idempotency retry semantics are broken.** Partial failures leave the idempotency row set but do not resume partial work on retry. Fix: specify the idempotency state machine, or per-operation idempotency keys.

3. **The domain pack contract is thin.** Six field names, no field schemas. Fix: provide a fully worked BM-IST PackDefinition, including one retirement rule per debt item and the semantics of `Validate()`.

4. **The refusal test suite is under-specified.** Eight desired properties, no test specifications. Fix: for each line, specify entry point, input, expected outcome, assertion.

Plus the smaller issues in §5: public API stability commitment, auth mechanism details, verifier artifact storage and enforcement, dead-end linkage, REOPEN schema, GO-decision review, test disposal misclassifications, ledger verification method, and Solvent rejection fallback.

None of these are fatal. All are fixable in one more revision. The plan has moved from "aspirational" to "concrete," and the remaining gaps are the kind that must be closed before implementation begins, not the kind that can be discovered during it.


## Z

# ADVERSARIAL REVIEW — PIVOT PLAN, REVISED (12 AMENDMENTS)
## Amendment fidelity is now high; the revision fixed the authority holes and opened two new ones — one of which deadlocks the retry path it was written to enable

**Executive verdict:**

$$
\boxed{\text{10 of 12 amendments landed cleanly; this is a genuinely executable plan one amendment-pass away from approval. But §1.4a's idempotency design as written makes the partial-failure retry it promises impossible (unique constraint = "attempted" conflated with "completed"), and the plan now contains three internal contradictions it must resolve before Phase 1 — the same class of defect the amendment pipeline exists to catch.}}
$$

---

## 1. AMENDMENT-FIDELITY AUDIT

| # | Amendment | Status | Notes |
|---|---|---|---|
| 1 | Real Phase-0 gate | ✅ Landed | §0.8 triggers are concrete; but see **N-gate** (GO granted despite firing the REVISE example) |
| 2 | Solvent boundary via public projection contract | ✅ Landed | §1.2a is the right shape (public methods, internals stay internal) |
| 3 | Packets never discharge | ✅ Landed | RetireDebt removed from agent flow; `proposed_retirement` as metadata is a valid alternative to PROPOSED-status — but persistence location + UI surface unspecified (**N6**) |
| 4 | Edges non-authoritative | ⚠️ **Weakened + inconsistent** | §2.4a demotes the gate to application-layer convention ("should be sufficient if the application layer does not act") — then §3.2's refusal suite assumes a confirmation concept §2.4a declined to build. See **N2** |
| 5 | Verifier trust binding | ⚠️ **Partial — tolerance dropped** | Schema has id/version/input_hash/artifact_hash but **no tolerance field**, despite tolerance being in the amendment's own title. Numeric re-verification remains non-adjudicable (**N4**) |
| 6 | Separate-transaction model | ✅ Landed | Explicit non-atomicity, retryable, idempotent — correct POC call; but the idempotency implementation is broken (**N1**) |
| 7 | DB-enforced idempotency | ❌ **Implemented with a deadlock** | See **N1** — headline finding |
| 8 | Typed PackDefinition | ✅ Landed | Interface in core (`internal/domainpack`), pack implements — A7 honored; rule *schema* unspecified (**N15**) |
| 9 | CRDB migration with type mapping | ✅ Landed | Exemplary detail; two concrete defects: **N3** (go:embed cross-module impossible) and **N5** (retry wrapper incomplete) |
| 10 | Auth + refusal suite | ✅ Landed | §0.6 + §3.2 + acceptance #15; but test vectors underspecified (**N7**) and one test has no referent (**N2**) |
| 11 | Health polling + honest metrics | ✅ Landed | Polling loop present; "2 process instances + OpenCode" is the honest count |
| 12 | Migration ownership | ✅ Landed | Clean split; acceptance #17; but the ownership mechanism needs the **N3** fix, making it Solvent change #2 |

**Pattern across rounds, for the record:** round 1 dropped the enforcement findings entirely; round 2 restored them as amendments; round 3 implemented them but introduced implementation-level defects. The pipeline is converging — the defect class is migrating from architecture to detail, which is what progress looks like.

---

## 2. NEW FINDINGS

### N1 — P0: The idempotency design deadlocks exactly the retry it exists to serve (§1.4a)

The written sequence:

```
2. INSERT into submission_idempotency (fails on duplicate = already submitted)
3. Persist beliefs, evidence, edges, tasks
4. On any failure after step 2, the idempotency row remains — caller must retry with same packet_id
```

The retry re-executes step 2 with the same packet → same `content_hash` → **UNIQUE constraint violation → "already submitted."** But the packet was never fully persisted. The caller is permanently locked out of completing a partial packet: the design conflates *attempted* with *completed*, and step 4's instruction is unfulfillable as specified. Acceptance #16 only tests duplicate-after-*success*, so the bug ships undetected.

**Fix (structural truth, per Solvent's own §22 doctrine):** move the guarantee to the entities, not the attempt — content-hash unique constraints on belief/evidence rows per scenario (retry re-execution becomes naturally idempotent at entity level), and write the `submission_idempotency` row **last**, as a fast-path response cache: its presence then certifies *completion*, and a duplicate INSERT legitimately means "fully done, return cached result." Add acceptance **#16b: partial-failure retry completes the packet with no duplicate state** — this test would have caught N1, which is the strongest argument for adding it.

### N2 — P1: §2.4a and §3.2 contradict each other on edge confirmation

§2.4a: PROPOSED/ACTIVE edge persistence introduced *only if* inspection shows agent edges can trigger transitions; otherwise the existing edge model suffices under application-layer discipline. §3.2 refusal suite: `unconfirmed contradiction → retract REFUSED` — **"unconfirmed" has no referent** in a model with no confirmation status. Either (a) the edge model needs a confirmation mechanism (a schema question — and per §1.5, epistemic schema is Solvent-owned, so this is a *third* Solvent module change), or (b) rewrite the test to what the model actually enforces: `agent → retract via any path REFUSED` (agent HTTP surface: none; MCP: two tools) and `human retraction citing agent-authored contradiction → ALLOWED` (which is exactly demo Branch A). Resolve before Phase 1 — the answer determines whether a Solvent schema change enters the budget.

### N3 — P1: `go:embed` cannot embed another module's files (§1.5)

The migration-ownership mechanism — "ARGUS consumes Solvent's migration bundle via `go:embed` of Solvent's migration SQL" — is **impossible as written**: Go's embed directive cannot match files outside the package's module, and Solvent is a separate module. Fix: Solvent exports its own `//go:embed migrations/*.sql` as a public `embed.FS` (or, better, a `solvent.ApplyMigrations(db)` function); ARGUS calls it. Consequence: this is a **second bounded Solvent module change** alongside the projection contract — update §0.9's change ledger, because the plan currently claims the projection contract is the only one. Which surfaces:

### N-gate — P2: the GO decision fired its own REVISE trigger without saying so

§0.8: bounded Solvent module change → **REVISE**. §0.9: GO, with "a bounded public API addition" acknowledged. Both can't be right unless the intent is *REVISE → incorporate → re-gate → GO*, which is what actually happened. State it explicitly ("gate fired REVISE on the projection contract; resolved by §1.2a; now also resolving the migrations-FS addition from N3") — the gate's integrity is the entire point of amendment 1, and a gate that grants GO while its trigger condition is visibly true teaches the next reader to ignore gates.

### N4 — P1: tolerance dropped from VerificationArtifact

The amendment title promised it; the schema lacks it. The program's first real artifacts (shear-counterexample reproduction, G2 Madelung cancellation) are floating-point — two runs *will* disagree at the last bits, and without a declared tolerance the refusal rule "artifact re-run within declared tolerance → pass; outside → flagged" (prior review, and the program's own verifier-POC design) has no operand. Add `Tolerance string` (or per-step tolerances) to the schema and the comparison rule to acceptance #3. Also add: **input provenance binding** — `input_hash` is recorded, but nothing constrains *whose* inputs; the claim (belief) should pin the expected input spec via `VerifierSpec`, or the agent chooses inputs that make a weak check pass.

### N5 — P2: the CRDB retry wrapper is incomplete (§1.5)

It retries only commit-phase failures (`isRetryable` on `Commit()` error); a 40001 surfacing from a statement *inside* `fn` returns without retry — but CockroachDB surfaces serialization conflicts on any statement. Standard pattern: retryable error anywhere ⇒ re-execute the whole closure. Also: `for { ... }` with `continue` on retryable is **unbounded** — add max attempts + backoff. Affects `Transition()`, `TransitionWithAgent()`, `Claim()` as noted.

### N6 — P1: `proposed_retirement` is unanchored

"Recorded as metadata for human review" — in which table/column, and shown where in the UI? If it isn't rendered on the Debts view as "agent ⟨id⟩ proposes retiring ⟨debt item⟩," it is dead data and demo step 12 loses its agent-proposal context — the exact human-visibility property that made amendment 3 worth making. One paragraph in §1.5 + one line in §2.2 fixes it.

### N7 — P2: refusal-suite test vectors are underspecified

By design there is *no agent path* to discharge/promote/retract (MCP exposes two tools; "agent endpoints (if any HTTP) are read-only"). So what does the test `agent → discharge REFUSED` actually invoke? Recommend committing to the cleaner contract — **no agent HTTP surface exists in the POC** (agents are MCP-only) — and specifying vectors: (i) `submit_packet` containing a forged decision-type → refused at compile; (ii) HTTP mutation call presenting agent-class principal → refused by caller-class check; (iii) the N1 partial-retry case. As written, two of the eight refusal tests are near-vacuous.

### N15 — P1: Discharge vs. RetireDebt is still undefined

The kernel Contract has both `RetireDebt` and `Discharge` (22 methods, per recon). Amendment 3 banned the agent path but the plan never defines which method the *human* flow's semantics are, what the other one is for, and whether it should be deprecated in the ARGUS-visible surface. §2.6 uses `Discharge`; acceptance #6 says "debt retired." The first developer to touch `human.go` will guess. One Phase-0 paragraph: define the two, mark one ARGUS-forbidden (depguard-able if they're distinct methods), or plan the kernel deprecation.

### N8 — P2: verifier failure semantics during compile

`compile.go` invokes the verifier when a packet includes verification evidence. Specify: verifier *error/crash* ⇒ artifact recorded `inconclusive`, packet proceeds (epistemically, a failed check is evidence too — but an infrastructure failure must not silently read as "refuted"); verifier *refuted* ⇒ evidence recorded with `result: refuted`, packet proceeds, human sees it. Also add a timeout — CAS numerics in the synchronous compile path.

### N-art — P1 (residual, sharpened): artifact persistence path still unspecified

Prior P1-4, still open and now sharper: where does a `VerificationArtifact` *live*? The flow removes the file-based registry but §2.3's CLI (`argus verify`) produces an artifact — into what? Demo step 8→12 requires it attached as evidence for discharge. Specify: evidence rows carry `class`, `artifact_json/hash/tolerance` columns (per class typing, prior F10 — still unaddressed), and **the CLI does not write to the DB directly**; the artifact enters via the same validated paths (agent packet or human attach), or via an explicit `argus verify --submit` that goes through the application layer with identical validation. One path, same rules — otherwise the CLI is an unvalidated mutation surface outside the packet discipline.

---

## 3. RESIDUAL PRIOR FINDINGS (unchanged, for the ledger)

| Item | Status |
|---|---|
| Unused authority surface (`authority.Service`, `executor.Registry` imported, never exercised) | Still open — DEFER the import; demo has no execution path |
| Scenario containment | Still unmentioned — with one DB it is a FK + check in `compile.go`; name it or it evaporates |
| TRUNCATED ≠ COMPLETE for context assembly | Still unmentioned — one marker field |
| Adversarial-agent independence | Correctly a non-goal now — but add one line to §3.1: independence is procedural in this demo |
| Isolation-test mechanics | Still hand-wavy ("if they reference BM-IST types, compilation fails" — a Go test cannot fail compilation); implement as a `go/packages` import-graph assertion + depguard, name the tool |
| "Nothing is discarded" (Conductor) | Still overstates — its own table defers dependency propagation |
| Boundary matrix | Incomplete: `application`'s allowed imports unstated (should be: epistemic, work, domainpack-*interface*, verifier, packet); add the row |
| Edge persistence in flow diagram | §2.4's flow omits edges that §1.4a step 3 includes — reconcile |
| `governance_ref` | Now a plain string where a real FK is possible (one DB!) — structural-truth upgrade opportunity |
| `task dev` runs `reset` unconditionally | Acceptable for POC; document that demo seed data dies on restart |

---

## 4. WHAT IS NOW GENUINELY GOOD

Worth saying plainly, because three rounds in the trend matters: the recon is specific enough to be falsifiable (22-method Contract, 10 migrations, `store_test.go:542`); the transaction model decision is the honest one (no fake global atomicity); migration ownership resolves a real ambiguity cleanly; the type-mapping table including the pre-existing `Release()` non-transactional bug is the kind of find that only comes from actually reading code; the refusal suite exists at all, as acceptance criteria; and the complexity ledger now says "2 process instances" like an engineer rather than "1 process" like a brochure. The plan has crossed from architecture to engineering. The findings above are the last mile.

---

## 5. INSERTION-READY FIXES (B-series)

| # | Fix | Location | Blocks |
|---|---|---|---|
| B1 | Entity-level content-hash unique constraints; idempotency row written last as response cache; acceptance #16b (partial-retry completes) | §1.4a, §1.5, acceptance | **P0** |
| B2 | Resolve edge-confirmation: either status column (Solvent change #3, REVISE-gate it) or rewrite the refusal test to the enforceable contract | §2.4a, §3.2 | P1 |
| B3 | Migrations via Solvent-exported `embed.FS`/`ApplyMigrations`; update Solvent-change ledger to two items | §1.5, §0.9 | P1 |
| B4 | Add `Tolerance` + input-provenance binding to VerificationArtifact; re-run comparison rule into acceptance #3 | §2.3 | P1 |
| B5 | Define Discharge-vs-RetireDebt; mark one ARGUS-forbidden | Phase 0, §1.2 | P1 |
| B6 | Anchor `proposed_retirement` (column + Debts-view rendering) | §1.5, §2.2 | P1 |
| B7 | Fix retry wrapper (closure re-execution, max attempts); commit to no-agent-HTTP and specify refusal vectors; verifier error≠refuted + timeout; artifact persistence path; `application` boundary row; scenario containment; truncation marker; isolation test via go/packages; FK for `governance_ref` | various | P2 |

---

## 6. CLOSE

The revised plan has earned its GO on structure: every thesis invariant now has a named mechanism, and the enforcement findings from two rounds of review are implemented rather than acknowledged. What remains is the last defect class — *the plan's own mechanisms contradicting each other or themselves* — concentrated in four places: **N1** (the retry deadlock, which the acceptance suite would only catch if B1's #16b exists), **N2** (the edge-confirmation contradiction), **N3** (the impossible embed), and **N4** (the dropped tolerance). All four are paragraph-sized fixes. Insert B1–B6, re-run the gate with the corrected Solvent-change ledger, and this plan is ready for Phase 1 — with the honest note for the record that the most valuable reviewer of this document has been, at every round, the previous version of itself.