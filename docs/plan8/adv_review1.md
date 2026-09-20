## Deepseek

# Adversarial Review: PIVOT_POC_IMPLEMENTATION_PLAN.md

## Short verdict

The plan is well-structured, disciplined, and directionally correct. It has absorbed the previous review's concerns: the Phase 0 gate is real (in principle), the complexity classification is present, the simplification ledger is concrete, the portability proof is formalized, and the amendments are visible.

But three issues remain:

1. **The Phase 0 gate is a formality, not a gate.** The "GO" decision is the conclusion; the reconnaissance findings are the justification written backwards.
2. **The package boundaries are asserted, not enforced.** Depguard and a compile-time test are mentioned but not specified enough to know whether they will actually work.
3. **The domain pack interface is under-specified.** The plan says the pack owns vocabulary and rules but does not define the API surface between pack and kernel. Without it, the domain-agnostic claim cannot be evaluated.

I will go through each and then list the missing pieces.

---

## 1. The Phase 0 gate is a formality

The plan says:

> Phase 0 must produce a GO / NO-GO / REVISE decision before any code changes.

But the reconnaissance findings are already written with the conclusion in mind. For example:

- The complexity classification table has 11 rows, and 8 are classified as "deployment." No alternative classifications are considered. The table is constructed to support the "GO" conclusion.
- The work-semantics inventory concludes "Nothing is discarded." This is the desired outcome, not a discovered one. What if something *were* discardable? The inventory would have said so, but it doesn't have a category for it.
- The Phase 0.7 output is "Decision: GO" with five supporting bullets, each of which is a confirmation of what the plan wants to be true.

**Failure mode:** If the reconnaissance were to reveal that, say, Conductor's dependency propagation is load-bearing and non-trivial, the plan would have to revise. But the plan never considers this possibility. It assumes the GO conclusion and works backwards.

**What a real gate would look like:**

- A list of specific findings that would trigger NO-GO.
- A list of specific findings that would trigger REVISE.
- An explicit statement of what was checked and what was not checked.
- A commitment to revisit if reconnaissance reveals something different.

**Concrete question the plan must answer:** What would make you say NO-GO? If no answer is possible, the gate is not a gate.

---

## 2. Package boundaries are not mechanically enforced

The plan says:

> Enforcement:
> 1. `depguard` rules in `.golangci.yml` for dev feedback
> 2. Compile-time negative test in `internal/boundary/boundary_test.go` that fails if forbidden imports exist

This is correct in spirit. But:

**2.1 The depguard rules are not specified.** What are the exact import paths that are forbidden? The plan lists general directions ("epistemic → may NOT import application, work, api, mcp, ui, domainpack") but does not translate them into depguard's configuration format. Depguard is easy to get wrong — it has a `deny` mode and an `allow` mode, and misconfiguration can silently allow forbidden imports.

**2.2 The compile-time test is not described.** A "compile-time negative test" in Go is usually implemented by writing a test that tries to import a forbidden package and expects the build to fail. But Go does not have a built-in "expect this import to fail" mechanism. Options:

- A test that uses `go/build` to check import paths (fragile).
- A test that runs `go list -deps` and checks the output (external tooling).
- A test that relies on a build tag and checks for a compile error (unusual).

The plan does not say which. Without a concrete implementation, the test is aspirational.

**2.3 No mention of CI enforcement.** The plan says depguard runs in `.golangci.yml`, but does not say it runs in CI. If it runs only locally, developers can bypass it. The plan should say `task test` includes `golangci-lint run`.

**2.4 The boundary list is incomplete.** The plan lists forbidden imports for `epistemic`, `work`, `verifier`, and `mcp`, but not for `application`, `domainpack`, `api`, `ui`, or `decision`. If those packages are also boundary-sensitive, they need rules too. If they aren't, the plan should say why.

**What the plan needs:** An appendix listing the exact depguard configuration, the exact implementation of the boundary test, and a statement that both run in CI.

---

## 3. The domain pack interface is under-specified

The plan says:

> Pack owns the vocabulary (claim vocabulary, evidence classes, retirement rules, falsifiers, domain verifier).

And:

> Pack interface:
>     GetPackID() string
>     GetVersion() string
>     Validate() error

This is a three-method interface. It cannot possibly capture what the plan says the pack owns. The pack must also:

- Declare its debt item names.
- Declare its evidence class names.
- Map debt items to evidence classes.
- Specify retirement rules that the kernel can execute.
- Provide a domain verifier.
- Provide vocabulary for claims, evidence, and debt descriptions.

The three-method interface is a stub, not a specification. The plan says "the plan must list the specific interfaces the domain pack must satisfy, so the argument has concrete form" — but the interfaces listed do not satisfy that requirement.

**3.1 The retirement rule format is undefined.** The plan says retirement rules map debt items to evidence classes. But what form do these rules take? Are they:

- Declarative data (JSON, YAML)?
- Go functions?
- Rules in a DSL?
- Callbacks the kernel invokes?

Without this, the "domain-agnostic" claim cannot be evaluated. If the rules are Go code, then each domain requires a Go package. If they are data, the kernel must interpret them. The two have very different consequences for portability.

**3.2 The domain verifier interface is undefined.** The plan says the pack provides "a domain verifier" but does not say what interface it must implement. Is it the same as the physics verifier? Does the kernel call it? Does the human call it? Does the pack call it?

**3.3 The debt vocabulary is undefined.** The plan says "Debt remains opaque to Solvent." But Solvent's schema has a `debt_discharge` table, a `RetireDebt` operation, and a promotion gate that checks whether debt is empty. Solvent must know *something* about debt: that it exists, that it can be listed, that it can be discharged. The plan conflates "Solvent doesn't know what `needMap` means" with "Solvent doesn't need to know about debt at all." The former is true and important. The latter is not achievable given the promotion gate.

**What the plan needs:** A concrete interface specification for the domain pack, including:

- The debt item type (an opaque string? a struct? a typed value?).
- The evidence class type.
- The retirement rule type.
- The verifier interface.
- A worked example using the BM-IST pack.

Without this, the portability argument is aspirational.

---

## 4. The DB migration is treated as trivial

The plan says:

> CockroachDB migration from SQLite may have syntax differences (mitigated: Conductor schema is simple, CRDB is Postgres-compatible)

"Postgres-compatible" is not "identical to SQLite." SQLite has different type affinities, different handling of `AUTOINCREMENT` (it uses `AUTOINCREMENT`; CRDB uses `SERIAL` or `UUID`), different default value semantics, and different transaction behavior.

**4.1 The plan does not say what it verified.** The reconnaissance says Conductor has 4 tables and uses `modernc.org/sqlite`. It does not say what the table definitions look like. If any column uses `INTEGER PRIMARY KEY AUTOINCREMENT`, the migration requires a rewrite. If any table relies on SQLite's loose typing, the migration requires explicit types.

**4.2 The plan does not address the atomic compare-and-swap SQL.** The reconnaissance says Conductor's state transitions use `WHERE status='proposed' AND current_agent IS NULL`. This works in SQLite and CRDB, but the plan should verify that the concurrency semantics are preserved. SQLite serializes writes; CRDB uses optimistic concurrency with retries. The application layer must handle CRDB's retry semantics.

**4.3 Solvent already handles CRDB retries.** The plan says Solvent's writes go through `crdb.ExecuteTx`. If the work subsystem is separate, it needs the same retry handling. The plan does not say it will use the same helper.

**What the plan needs:** A concrete migration checklist listing each Conductor table, its SQLite definition, and its CRDB translation. If any column or constraint requires special handling, state it.

---

## 5. The auth model is thin

The plan says:

> MCP mode: only exposes 2 tools. No HTTP server started.
> HTTP mode: `/api/retire`, `/api/promote`, `/api/retract` require Origin/Host header check + confirmation token. Attribution derived from server-side principal, not request body.
> No full identity system. Minimal CSRF for browser actions.

This is a reasonable POC scope. But the plan does not say:

**5.1 Where the confirmation token comes from.** Is it a session cookie? A one-time token issued by the UI? A static secret? If the token is a static secret stored in the server, the agent cannot reach it via MCP (good), but a malicious browser could still use it if the token is exposed in the UI.

**5.2 What "server-side principal" means.** Is it a single default principal for all POC users? A configurable principal? If there is only one principal, "attribution" is nominal. If there are multiple, how are they distinguished?

**5.3 What prevents an agent from crafting a malicious packet that triggers a human-like decision.** The plan says agents only have `get_context` and `submit_packet`. But `submit_packet` can include a `discharge` request. Does the application layer distinguish "packet requests discharge" from "human requests discharge"? If not, an agent can trigger retirement via a packet.

**What the plan needs:** A concrete specification of the confirmation token flow and the mechanism that prevents agent-triggered consequential actions.

---

## 6. The simplification ledger is unverified

The plan says:

> | Dimension | Before | After | Verified? |
> | Processes | ~7 | 1 | |
> | Ports | ~6 | 1 | |
> | Databases | 1-2 | 1 | |
> | Go modules | 3 | 1+import | |
> | Credential sets | multiple | 1 | |
> | Internal HTTP calls | many | 0 | |
> | Cross-process handoffs | yes | none | |

The "Verified?" column is empty. The plan says verification happens in Phase 3, but does not say how. Some of these are easy to verify (count processes), some are not:

- "Credential sets: multiple" — how many? "Multiple" is not a number.
- "Internal HTTP calls: many" — how many? A count is possible (grep for HTTP client usage), but the plan does not say it will do this.
- "Cross-process handoffs: yes" — which ones? A list would be more concrete.

**What the plan needs:** Concrete numeric baselines for each dimension, and a specific verification method for each.

---

## 7. The null domain-pack test proves less than it seems

The plan says:

> Add a null domain-pack CI compile test: the core builds against a minimal stub domain pack that satisfies the interface with empty/null implementations. This proves at compile time that the core does not import BM-IST vocabulary.

This tests one thing: that the core compiles without importing BM-IST packages. It does not test:

- That the core's behavior is correct without domain-specific rules.
- That a domain pack can be swapped at runtime.
- That the pack interface is sufficient to capture domain-specific semantics.
- That the retirement rules, evidence classes, and verifier interface are actually domain-agnostic.

The test is a compile-time check, not a portability proof. The plan should rename it and adjust its claim.

**What the plan needs:** The compile test is necessary but not sufficient. The plan should add:

- A behavioral test that runs the core against a null pack and verifies it does not crash.
- A specification of the pack interface (see Section 3 above).

---

## 8. The plan does not address domain pack complexity

The plan claims domain agnosticism by construction. But the BM-IST domain pack is not simple:

- 6 debt items (`needMap`, `needInvariant`, etc.).
- Multiple evidence classes.
- Retirement rules that depend on the pack's vocabulary.
- A physics verifier that is domain-specific.

If a domain pack is this complex, "domain agnosticism" is not a strong claim. The plan should acknowledge that a new domain pack is nontrivial work. The claim is not "any domain requires zero effort" but "any domain requires a new pack, not a rewrite of the core." That's a weaker claim, but it's the honest one.

**What the plan needs:** A statement of the pack's expected size and complexity, and a criterion for "this is a new domain pack, not a rewrite."

---

## 9. Missing pieces

The plan does not address:

**9.1 Rollback and failure recovery.** The migration model says "revert to main" if Phase 1+ fails. But if Phase 1 introduces schema changes to CockroachDB and the migration is partly applied, what happens? The plan should say whether Phase 1 uses a migration that is reversible.

**9.2 Data loss prevention.** The plan says "no live state to protect (POC, fresh DB each run)." But if a developer has been running the POC and has interesting scenarios, wiping the DB is destructive. The plan should say whether there is an option to preserve state.

**9.3 Concurrent developer use.** If two developers run `task dev` on the same machine, they both try to start CockroachDB on `:26260`. The plan does not address port allocation or per-developer DB isolation.

**9.4 Error handling for MCP disconnect.** If OpenCode disconnects mid-packet, what does the server do? The plan does not say.

**9.5 The `argus reset` semantics.** The plan lists it as a subcommand but does not say what it does. Drop all tables? Truncate? Log the reset?

**9.6 The `argus verify` interface.** The plan says it "runs physics verifier and produces artifact." What artifact? Where is it stored? How is it referenced in a packet?

**9.7 The relationship between `argus verify` and in-process verification.** The plan says the verifier is called both by the CLI and by the application layer during packet validation. If the CLI produces an artifact on disk and the in-process call produces an artifact in memory, are they the same artifact? Are they cross-comparable?

**9.8 The Trust UI's role in the decision flow.** The plan says humans use the Trust UI to discharge, promote, retract. But it does not say what the UI shows, what confirmation is required, or what happens if the action fails.

---

## 10. What the plan gets right

To be fair:

- The Phase 0 concept is correct, even if under-implemented.
- The complexity classification is the right framing.
- The component disposition table is concrete.
- The 4-phase structure is lean.
- The hard constraints are appropriate.
- The migration model is stated.
- The test disposal matrix is present.
- The acceptance criteria are specific and checkable.
- The simplification ledger is a concrete target.

The plan is a significant improvement over the previous state. The gaps are in enforcement and specification, not in direction.

---

## 11. Concrete questions the plan must answer before implementation

1. What would make Phase 0 say NO-GO or REVISE? State it explicitly.
2. What is the exact depguard configuration, and does it run in CI?
3. What is the concrete implementation of the boundary test?
4. What is the full interface the domain pack must satisfy?
5. What is the retirement rule format (data or code)?
6. What is the domain verifier interface?
7. How does the application layer distinguish agent-submitted discharge from human-submitted discharge?
8. Where does the confirmation token come from, and how is it protected?
9. What is the concrete baseline for each simplification ledger row?
10. What does the Conductor-to-CockroachDB migration look like, table by table?
11. What does `argus reset` do?
12. What does `argus verify` produce, and where is it stored?
13. What is the rollback criterion for Phase 1+ failures?
14. What happens when two developers run `task dev` on the same machine?
15. What is the expected size and complexity of a new domain pack?

---

## 12. Bottom line

The plan is the right shape. It has absorbed the previous review's concerns, added the reconnaissance gate, classified complexity, and produced a concrete simplification target. It is a serious attempt to reduce the POC's accidental complexity without losing the thesis.

But it is not yet implementation-ready. Three things are still aspirational rather than specified:

1. **The Phase 0 gate** — it exists in form but not in function. The GO decision is the conclusion of the plan, not the result of reconnaissance.
2. **The package boundaries** — they are listed but not enforced. Depguard and the compile-time test are mentioned but not specified.
3. **The domain pack interface** — it is three methods when it needs to be at least ten. The portability claim depends on an interface that has not been defined.

These are fixable. The next revision should add:

- An appendix with the exact depguard configuration.
- A concrete domain pack interface with a worked example.
- A real Phase 0 gate with stated NO-GO and REVISE conditions.
- A table-by-table Conductor migration plan.
- A concrete auth flow specification.

If those are added, the plan becomes implementation-ready. Without them, it is a good blueprint for a plan.


## Z

# COLLECTIVE ADVERSARIAL REVIEW — background.md + prompt.md + PIVOT_POC_IMPLEMENTATION_PLAN.md
## The pivot is well-executed; the amendment pipeline is lossy in exactly the wrong place

**Executive verdict:**

$$
\boxed{\text{The plan is structurally excellent and the GO decision is plausibly correct — but the review→amend→plan pipeline dropped both P0 security findings from the prior review, and the plan then introduced a new one of its own: the agent packet flow calls an authority-mutating kernel method. The plan demonstrates the system will be simpler. It does not yet demonstrate the system will refuse.}}
$$

Three summary points before the detail:

1. **Amendment fidelity ≈ 70%.** The structural amendments all landed faithfully (recon-gating, semantics inventory, import-direction rules, null-pack test, simplification ledger, test disposal). The two highest-severity findings of the prior review — **F1 (agent-authored edges)** and **F3 (verifier trust chain)** — did not survive into `prompt_amended.md`, and therefore not into the plan.
2. **One new P0 in the plan itself.** §2.4's agent packet flow includes `epistemic.RetireDebt [Solvent kernel, if applicable]`. Read literally, an agent packet triggers a debt-retirement authority transition — the exact violation of `CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION` that the thesis is named for.
3. **The acceptance suite is all positive-path.** Thirteen criteria, every one a "can do." A system whose entire thesis is *refusal of unauthorized transitions* is being accepted without a single refusal test.

---

## 1. WHAT THE COLLECTIVE GETS RIGHT (confirmed)

- **The complexity classification table** (deployment / API / inherent-domain) is the correct instrument, and its verdict — majority deployment complexity — justifies the pivot on evidence-shaped grounds.
- **Reconnaissance quality is high.** Specific claims (22-method Contract, 10 migrations / 20+ tables, Conductor's 4 types / 4 tables, SQLite via pure-Go driver) are the right granularity, and "code is the source of truth" was honored.
- **Phase discipline:** exactly 4 phases, GO/NO-GO gate, rollback criteria, delete-first honored in the disposition table (Solvent API, wizard, binaries, reference-loop all deleted without sentiment).
- **The two subtle calls are right:** Solvent imported as a module (kernel boundary enforced by module separation, standalone-product option preserved); Conductor's domain model extracted rather than deleted (the inventory showed it small and load-bearing — the honest outcome, even though "Nothing is discarded" overstates slightly: dependency *propagation* is deferred, per its own table).
- **§14's two answers are now commitments**, correctly worded: nothing architectural is lost by collapsing Coordinator/Conductor; Solvent survives because it is the only subsystem whose state is authoritative *and* the candidate reusable product.

---

## 2. AMENDMENT-FIDELITY AUDIT

| Prior finding | Amended prompt | Plan | Residual gap |
|---|---|---|---|
| F9 Coordinator rule-accretion → import-direction + mechanical enforcement | ✓ §3 | ✓ §0.5, §1.6 | None material |
| F5 dual-write saga → single-DB transaction | ✓ implicit | §2.4 "single CRDB transaction where possible" | **Broken by N1/N2 below** (RetireDebt ambiguity; tx-participation not designed) |
| F13 readiness orchestration → `task dev` | ✓ §15 | ✓ §3.5 | **Regressed to `sleep 2`** (N5) |
| F6 RetractCascade coupling | dissolved | dissolved in-process | None |
| F2 Trust UI auth/attribution | ✓ §4.5 (new) | ✓ §0.6, §1.7 | **Model defined, zero tests** (N6) |
| F3 verifier trust chain (script-hash pre-registration, tolerance, env pinning) | ✗ **dropped** | ✗ absent | **Open P0** (N13) — "trust-chain-bound" appears in the amendment as a phrase; the plan never defines the chain |
| F1 agent-authored edges PROPOSED/ACTIVE gate | ✗ **dropped** | ✗ absent | **Open P0** — demo steps 9→14 route retraction through an unconfirmed agent-authored contradiction |
| F10 evidence typing (PROOF_OBLIGATION vs COMPUTATION) | partial (pack owns evidence classes) | implicit in pack | Validation-per-class and the L1/L2 human-only rule unstated; artifact schema absent |
| Kernel extraction import rule + interface location | ✓ §3 | ✓ §0.5 | **Pack *interface location* undecided** (N4) — the keystone of the portability proof |
| Null-pack portability test | ✓ §9 | ✓ §3.2 | Mechanics hand-wavy; see N-minor |
| Complexity ledger | ✓ §16 | ✓ §3.3 | One honesty error (N8) |
| F8 agent freshness | explicitly NON-GOAL | deferred | Acceptable for POC — but then the demo's "adversarial independence" is procedural; say so in the demo notes |

**The pattern is legible: the pipeline preserved every finding about *structure* and dropped every finding about *authority enforcement*.** Those are the ones the demo exists to prove.

---

## 3. NEW P0 FINDINGS (introduced by the plan)

### P0-1 — `RetireDebt` in the agent packet path (§2.4)

The agent flow reads:

```
persist(packet) → EnterBelief → AddEvidence → RetireDebt [if applicable] → CreateTask → audit.Log
```

while the human flow (§2.6) uses `Discharge` for `RETIRE_DEBT`. Either (a) the agent path performs an authority transition — thesis violation, or (b) `RetireDebt` and `Discharge` are different operations — in which case **the plan never defines the difference**, and the kernel Contract (which has both methods) carries the ambiguity into the collapsed system. There is no third reading.

**Resolution (insert as Phase-0/1 decision):** packets may *propose* debt retirement only; packet-borne retirement enters at status `PROPOSED`, becomes effective solely via human `Discharge`; retraction/promotion preconditions check only *confirmed* state. If the kernel lacks the status column, that is a **Solvent module change** — see P1-1. And the demo must include the negative step: *an agent retirement attempt held at PROPOSED, visible in the UI as awaiting human confirmation.*

### P0-2 — Verifier script provenance is undefined (§2.3)

In-process verification fires "when a packet includes verification evidence." On *whose* script? If the agent supplies the script, the agent chooses what to verify — evidence laundering returns through the new front door, and the collapse removed the process boundary that used to make this at least visible. The amendment's phrase "trust-chain-bound" was never given content.

**Resolution:** pre-registered script hash per obligation (BM-IST pack data), declared numeric tolerance in the artifact schema (the natural checks are floating-point — two runs can disagree), environment pinning, and a kill rule: packet evidence whose script hash does not match the pack's registration is rejected at `compile.go`, not accepted as unknown. Note explicitly: **content hashing is not "cryptographic attestation infrastructure"** — do not let Hard-Constraint #7 be misread into deleting the trust chain. The prior review made this clarification; the plan's silence suggests it was misread exactly this way.

### P0-3 — Positive-only acceptance; auth untested

Acceptance criteria 1–13 are all "the machine can." The thesis is "the machine refuses." Missing negative assertions, each one-line testable: discharge with open debt → refused; promote after retraction → refused; promote with non-empty debt → refused (the gate is stated in #6 but only its success is tested); retraction against unconfirmed edge → refused; cross-origin POST to a consequential endpoint → rejected; body-supplied `discharged_by` → ignored, token attribution prevails; failed mutation → **no success audit row** (Solvent's own no-false-audit doctrine, made into a test). Add a corresponding ADD row to TEST DISPOSAL: `internal/api/auth_test.go`. Untested security is decorative security — and this is the surface the demo's Branch B depends on.

---

## 4. P1 FINDINGS

**P1-1 — Hidden Solvent-module dependency.** The GO decision assumes thin wrappers, zero kernel changes. But (i) P0-1's PROPOSED status likely needs a column; (ii) §2.4's single-transaction requirement conflicts with `crdb.ExecuteTx`'s internal transaction management — an outer application transaction cannot wrap kernel calls that manage their own tx. **Decide in Phase 0, not Phase 2:** either (a) in-process two-transaction saga with compensation and explicit orphan-marking (compliant with the no-distributed-transactions constraint; no kernel change), or (b) kernel gains tx-participation signatures (module version bump — a cross-repo change inside a "simplification" pivot, which must be budgeted, not discovered).

**P1-2 — Pack interface location undecided.** §3.2's null-pack test only works if `application` depends on a *generic* Pack interface defined **in the core**, with the BM-IST pack injected as an implementation. If the interface lives in `domainpack/bmist`, the portability proof is void. One sentence in §0.5 fixes it; without it, the plan's most important test tests nothing.

**P1-3 — Dedup cache is the wrong mechanism.** §1.4 carries over `idempotency.go` as a "dedup cache." In one process against one DB, retries after restart defeat an in-memory cache and duplicate beliefs (inflated evidential support — a corrupted record). Replace with a **content-hash unique index** in CRDB. This is Solvent's own §22 doctrine — structural truth over remembered state — applied to its own retry path.

**P1-4 — Artifact/evidence schema absent.** Where do verification artifacts live now that the file-based registry is gone? Presumably evidence rows with payload/hash/tolerance columns. One paragraph in §1.5. Also: evidence *class* column (typed validation per class) — currently only implied by the pack's evidence-class list.

**P1-5 — Ledger honesty: "Processes: 1".** `argus serve` and `argus mcp` are separate process instances (stdio is per-client; MCP cannot live inside the serve process). Correct claim: **1 binary, 2 process instances + OpenCode, 1 DB, 0 internal HTTP.** The simplification thesis survives intact; the ledger should not contain a claim acceptance #12 will falsify on a technicality.

**P1-6 — Unused authority surface retained.** `authority.Service` + `executor.Registry` are imported (KEEP as module) but the 17-step demo never executes an external action. Delete-first asks: DEFER the import. Keeping unexercised authorization machinery in a POC whose acceptance suite has no execution path is scope retention contradicting the plan's own doctrine — and it widens the audit surface P0-3 is supposed to cover.

---

## 5. P2 / MINOR

- **`sleep 2` in `task dev`** is the F13 anti-pattern reborn — Cockroach demo startup is variable; poll the health endpoint or a ready-file. The pivot removed readiness *orchestration*; it should not have replaced it with a race.
- **Scenario containment** (background's P1-3): with one DB it is a FK + a check in `compile.go` — name it explicitly or it silently evaporates in the collapse.
- **TRUNCATED ≠ COMPLETE** survives the pivot for projections (§2.5 context assembly). One marker field.
- **Null-pack test mechanics:** a Go test cannot literally "fail compilation"; implement as a `go/packages` import-graph assertion (core package imports must contain no `bmist` path) alongside depguard. Same intent, executable spec.
- **Phase 0 linchpin re-verification:** the GO decision rests on three claims ("kernel/services have zero HTTP imports" above all). Add acceptance criterion #0: a committed `go list -deps`-based script whose output is the evidence. Cheap insurance on the plan's single epistemic load-bearing wall.
- **Amended prompt §9 requirement not executed:** "explain exactly which parts would remain unchanged if BM-IST were replaced" — the null-pack test covers it mechanically; the plan should still contain the three-sentence answer (unchanged: beliefs/evidence/debt/edges/lifecycle/authority gates/MCP boundary/UI; per-domain: vocabulary, retirement rules, evidence classes, verifier, scenario semantics).
- **Demo step 8:** state where the artifact lands (evidence row per P1-4) and that its tolerance/hash are recorded — otherwise step 8 is unverifiable theater.

---

## 6. INSERTION-READY AMENDMENTS (A1–A10)

| # | Amend | Plan location |
|---|---|---|
| A1 | Packet-borne debt retirement = `PROPOSED` only; effective solely via human `Discharge`; define RetireDebt-vs-Discharge semantics in the data model; add negative demo step | §2.4, §1.5 |
| A2 | Verifier trust chain: pre-registered script hash per obligation (pack data), tolerance + env-pin fields in artifact schema, hash-mismatch rejection at compile.go; record the "hashing ≠ attestation" clarification | §2.3, §1.5, Hard Constraints note |
| A3 | Edge PROPOSED/ACTIVE gate: agent-authored edges non-actionable until human-confirmed; retraction preconditions check confirmed state only; UI labels agent-proposed edges | §2.6, §8 demo steps 9–14 |
| A4 | Negative-path acceptance criteria (refusal suite) + `internal/api/auth_test.go` added to TEST DISPOSAL/ADD | Acceptance criteria, §3 |
| A5 | Phase-0 sub-decision: saga-vs-kernel-tx strategy for packet persistence, with the Solvent-module-change budget if (b) | §0.7, §2.4 |
| A6 | Content-hash unique index replaces in-memory dedup cache | §1.4 |
| A7 | Pack interface defined in core (`internal/application` or a core contract package); BM-IST pack implements it; depguard rule: application may import the interface, never the bmist package | §0.5, §3.2 |
| A8 | `task dev` readiness by health polling; `argus reset` gated behind explicit confirmation flag | §3.5 |
| A9 | Ledger corrected to "1 binary / 2 process instances (+OpenCode)"; verification method noted per row | §3.3, §16 |
| A10 | Acceptance #0: committed script re-verifying the three Phase-0 linchpin claims | Acceptance criteria |

---

## 7. ON background.md

Superseded and correctly absorbed. Its diagnosis (infrastructure harder to reason about than the hypothesis), its kill question, and its component verdict are all carried into the plan faithfully; its one prediction worth logging — "the current fixes plan itself is a useful warning sign" — is now quantified by the complexity ledger. The orphaned final question (domain-agnosticity carryover) was answered in the prior review and is now implemented as §3.2 + A7. Nothing further needed from it.

---

## 8. CLOSE — THE PIPELINE ASSESSMENT

The review→amend→plan loop is **working on structure and lossy on enforcement**: every deployment-complexity finding dissolved cleanly (the pivot's promise kept — internal HTTP in the demo path: many → 0; the simplification is real), while the two findings about *who may cause which transition* fell out between the prior review and `prompt_amended.md`, and the plan then added a third at the exact same joint. That is not a criticism of the pivot — it is a measurement of the pipeline, and the fix is mechanical: **A1–A4 must enter the plan before Phase 1, because they are the difference between a demo that shows the machine running and a demo that shows the thesis holding.**

One sentence for the record, in the program's own idiom: the collapse removed the walls between rooms but the plan must not remove the doors' locks — and right now §2.4 shows a door standing open into the one room (authority) this entire architecture exists to guard.