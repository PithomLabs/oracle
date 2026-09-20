## deepseek

# Adversarial Review: The ARGUS POC Simplification Pivot

## Short verdict

The pivot is directionally right and overdue. The retrospective correctly diagnoses that the POC's complexity has drifted from the thesis it is supposed to demonstrate. The prompt is well-structured and includes the right hard constraints.

But the pivot commits the same error it accuses the old architecture of: **it decides before it inspects.** The retrospective has already concluded that Coordinator and Conductor should be removed, and the prompt frames the plan-writer's task as verifying that conclusion rather than testing it. The central claim — that collapsing processes preserves semantics — is asserted, not demonstrated. And the plan-writing prompt is silent on the enforcement mechanisms that would keep the collapsed process from becoming a monolith in the bad sense.

I will organize this around the specific claims that need to be tested.

---

## 1. What the documents get right

- The diagnosis is correct: the current POC has accumulated infrastructure complexity (cross-service reference validation, artifact handoff, auth propagation, readiness orchestration) that is orthogonal to the architectural thesis.
- The separation of conceptual from deployment boundaries is correct in principle.
- The critical design question in Section 14 of the prompt — "What architectural property would be lost if Coordinator and Conductor stopped being separate services?" — is exactly the right question.
- The hard constraints (no microservices, no message brokers, no second database, no autonomous promotion) are appropriate.
- "Delete-first" as a default is correct.
- The end-to-end demonstration in Section 8 is a coherent thesis statement.
- The portability framing in Section 9 correctly distinguishes the generic core from the domain pack.
- The developer experience target (`git clone && task dev`) is correct.

The pivot is a genuine simplification in spirit. The concerns below are about execution, not direction.

---

## 2. The central unexamined assumption

**The retrospective assumes that deployment separation is the source of the POC's complexity. It may be a symptom, not the cause.**

Consider the concrete problems the retrospective cites:

- Cross-service reference validation
- Verifier artifact handoff between processes
- Auth propagation across services
- RCP projecting from two backends

Each of these could exist in a single-process design if the APIs between the (now in-process) components are badly designed. A monolith with the same reference-validation logic, the same artifact handoff pattern (now in-process file writes), and the same auth propagation (now in-process context objects) is not simpler — it is the same complexity with less visibility.

**The question the plan must answer:** For each "complexity" line item, is the complexity from (a) deployment separation, (b) bad API design, or (c) inherent domain complexity? Collapsing into one process fixes (a) but not (b) or (c). If most items are (b) or (c), the pivot will not reduce complexity; it will relocate it.

**Concrete test:** Before committing to the pivot, classify each P0/P1 item in the current plan by this taxonomy. If the majority are (a), the pivot is justified. If not, the pivot's premise is wrong.

---

## 3. Package boundaries in a single Go process are honor-system boundaries

The retrospective says:

> "The important thing is that the packages cannot casually reach across each other's responsibilities."

This is a social convention, not a technical enforcement. In Go, `internal/work` can import `internal/epistemic` if a developer writes that import. Nothing prevents it. The `internal/` convention prevents *external* code from importing, but does not prevent *sibling* packages from importing each other.

**Consequence:** The pivot trades "expensive but enforced service boundaries" for "cheap but honor-system package boundaries." This is not necessarily bad — but the plan must acknowledge the trade.

**What the plan needs:** A concrete mechanism to detect boundary violations. Options:

- An import-lint rule (e.g., `depguard` or a custom `golangci-lint` rule) that forbids certain imports.
- A package boundary test that runs in CI.
- A documented ownership table with a named reviewer for each boundary.

Without one of these, "the packages cannot casually reach across" is aspirational.

---

## 4. The pivot decides before it inspects

The prompt says (Section 4):

> "Treat the following as the baseline hypothesis for the new plan: SOLVENT — KEEP conceptually. CONDUCTOR — REMOVE as a standalone service. COORDINATOR — REMOVE as a standalone service."

This is a decision, not a finding. The reconnaissance phase (Section 6) is supposed to inform the decision, but the decision has already been made.

**Consequence:** The plan-writer will discover Conductor's actual functionality, then be told to remove it. If Conductor turns out to have features that are actually load-bearing (e.g., dependency-driven task scheduling that the POC's dead-end/lineage semantics depend on), the plan must either (a) reintroduce those features in the new work subsystem or (b) accept their loss. The prompt does not tell the plan-writer how to make this trade.

**What the plan needs:** A prior statement that the removal decision is provisional and will be revised if reconnaissance reveals load-bearing Conductor semantics.

---

## 5. Conductor's "smallest possible work model" is under-specified

The prompt says the work subsystem should preserve:

```
task
dependency
state
priority
activity
cancellation/dead-end linkage
```

But the current POC's dead-end/lineage semantics (Phase 8 dry run: "linked Conductor task cancelled") may depend on more than these six primitives. For example:

- How is a "dead end" distinguished from a "cancelled" task?
- How does REOPEN create new lineage without erasing history?
- Does dependency propagation exist? If so, what is the propagation rule?
- What is the claim/release semantics for tasks?
- What is the ready/blocked computation?

The prompt's list is a first cut. The plan-writer must extend it based on reconnaissance. But if the plan-writer is told "smallest possible," they may under-specify and lose semantics the POC needs.

**What the plan needs:** A reverse-engineering task that enumerates the actual work semantics used by the current POC before deciding what to keep.

---

## 6. The "one database" decision is under-justified

The prompt says "Prefer one CockroachDB instance/database." But the current architecture may already use one CockroachDB with multiple schemas, or multiple DBs, or one DB with prefix-separated tables. The prompt does not specify which.

**Consequence:** If the current design already has one DB, the pivot does not reduce DB count. If it has multiple, the pivot reduces DB count but the schema merge is a nontrivial migration.

**What the plan needs:** A clear statement of (a) current DB topology and (b) target schema organization. If Solvent and Conductor currently have separate schemas, does the pivot merge them or keep them as separate logical namespaces within one DB? The prompt is silent.

---

## 7. Trust UI and API in the same process: security model unspecified

The prompt says "HTTP/API + Trust UI should be served by the same ARGUS process." This is fine for a POC. But the current architecture separates:

- The API for agents (via MCP)
- The Trust UI for humans (via a browser)

If both are in the same process, then the process must distinguish:

- Agent requests (should only touch `get_context` and `submit_packet`)
- Human requests (may touch discharge, promotion, retraction)

The auth model for these two callers must be different. The prompt does not discuss this.

**What the plan needs:** A statement of how the collapsed process distinguishes agent calls from human calls, and what enforces the capability boundary at the process level.

---

## 8. MCP as a subcommand: what happens to the boundary?

The prompt says:

> "MCP should be another mode of the same binary."

If MCP is a *mode*, then when ARGUS runs in MCP mode, it does not run the HTTP API. Fine. But if MCP is a *subcommand of the same binary*, then both modes share the same compiled code, the same dependencies, and the same potential for a bug in one mode to affect the other.

The current architecture has MCP as a separate process that talks to Coordinator over HTTP. This is friction, but the friction enforces the boundary. Collapsing MCP into the ARGUS binary removes this enforcement.

**What the plan needs:** A statement of how the MCP mode is isolated from the API mode at the code level (e.g., separate `internal/mcp` package with no imports from `internal/api` or `internal/epistemic`), and a test that verifies MCP cannot reach privileged operations.

---

## 9. The verification subcommand: process model is unspecified

The prompt says the verifier should be a library, not a separate process. This is a reasonable simplification. But:

- If `argus verify` runs as a subcommand, it is a separate process invocation of the same binary. It shares the binary but not the runtime state.
- If verification runs *within* `argus serve`, it blocks the server.
- If verification runs in a goroutine within `argus serve`, it shares the process but requires concurrency-safe design.

The prompt says "Physics verification should be a Go library/package, with a CLI/subcommand for running the reference verifier." This is ambiguous about runtime invocation. It also does not address whether the verifier needs external dependencies (numerical libraries, theorem provers) that would now be ARGUS's dependencies.

**What the plan needs:** A clear statement of the verifier's runtime model (subcommand, in-process, async goroutine) and its dependency footprint.

---

## 10. The BM-IST domain pack is the POC's scientific content, not config

The prompt treats the domain pack as a config-like module. But the BM-IST domain pack contains:

- Retirement rules (what counts as `needFaithfulnessReview` satisfaction)
- Evidence classes (what counts as a valid verification artifact)
- Falsifiers (what would kill a claim)

These are scientific claims. They can be wrong. The POC's success depends on the domain pack being correct enough to produce the demonstration.

The prompt does not address:

- How the domain pack is validated.
- Who owns it.
- Whether it is code or data.
- Whether it has tests.
- What happens if it is wrong.

**What the plan needs:** A statement of the domain pack's form (Go package, data files, or a mix) and its validation regime.

---

## 11. The portability "test" is an argument, not a test

The prompt says (Section 9):

> "Create a portability test in the plan."

But the plan is a document. A "test in the plan" is an argument, not an executable test. The plan can describe what *would* change if the domain pack were replaced. It cannot demonstrate that the change would be small.

**What the plan needs:** A clear statement that portability is argued by construction, not demonstrated by execution. And a list of the specific interfaces the domain pack must satisfy, so that the argument has a concrete form.

---

## 12. The "current implementation" may not match the briefing

The prompt says:

> "Code is the source of truth for implementation details."

This is correct. But it means the plan cannot be written without the reconnaissance phase. And yet the plan is expected to contain sections (target architecture, component disposition, phased implementation) that depend on reconnaissance.

**The logical dependency:** The plan's other sections cannot be filled in until Section 6 (reconnaissance findings) is complete. The prompt should have made this explicit — the plan is a *two-phase* document: reconnaissance findings first, then the plan based on those findings.

**What the plan needs:** A clear statement that the plan is contingent on reconnaissance, and that if reconnaissance reveals contradictions with the briefing, the plan will revise the briefing's assumptions.

---

## 13. Migration risk is unaddressed

The pivot touches every part of the codebase. The prompt does not address:

- Whether the pivot is a branch, a new repo, or an in-place rewrite.
- How the current POC continues to function during the pivot.
- Whether there is a rollback path.
- Whether the pivot is atomic or incremental.

For a POC, this may not matter. But if the current POC has live state or is being used for demonstrations, a botched pivot could lose work.

**What the plan needs:** A one-sentence statement of the migration model (in-place, branch, or new repo) and a rollback criterion.

---

## 14. Missing: disposal of current tests

The prompt says:

> "Which current tests encode real architectural requirements versus deployment mechanics."

But it does not say what to do with tests that encode *neither* — e.g., tests that verify the current service topology that will no longer exist.

**What the plan needs:** A disposal classification:

- Preserve (tests that encode real requirements)
- Delete (tests for removed services)
- Rewrite (tests for moved code)
- Add (new tests for new boundaries)

---

## 15. Missing: what "argus reset" does

The prompt says "argus reset (only if useful)." But reset semantics are tricky:

- Drop all data?
- Preserve history?
- Preserve scenarios?
- Log the reset?

For a POC, this is minor. But it should be specified because it is the first command a developer will run.

---

## 16. The strongest concern: the prompt conflates two kinds of simplification

The retrospective says:

> "This does not throw away the architecture we've developed. It removes deployment machinery while preserving the invariants."

This is the correct framing. But the prompt's phrasing does not distinguish between:

- **Deployment simplification** (fewer processes, ports, credentials, startup steps)
- **Architectural simplification** (fewer concepts, fewer invariants, fewer failure modes)

The pivot achieves (a). It does not necessarily achieve (b). The retrospective is honest about this; the prompt is not.

**What the plan needs:** A clear statement that the pivot simplifies deployment, not architecture. If architectural simplification is also a goal, it must be specified separately.

---

## 17. Concerns with the plan-writing prompt itself

Beyond the pivot's content, the prompt has structural issues:

**17.1 Too many required sections.** Eighteen sections in the plan. Some are redundant (Section 4 and Section 6 both require component analysis; Section 8 and Section 9 both require demonstration scenarios). The plan should prioritize depth on the load-bearing sections (reconnaissance findings, component disposition, package boundaries) and reduce or drop the others.

**17.2 Phase design is over-engineered.** Twelve phases. For a POC pivot, this is too many. Phases 0–11 imply many weeks of work. The plan should be designed to reach a decision point quickly (e.g., by Phase 3) and continue only if the decision is favorable.

**17.3 The "hard constraints" list is good but incomplete.** It forbids new services but does not forbid new packages, new interfaces, or new abstractions. It should also forbid adding complexity under the guise of simplification (e.g., a "clean architecture" refactor that introduces more layers).

**17.4 No priority ordering.** The prompt lists all requirements equally. It should prioritize: what must be true for the pivot to succeed vs. what is nice to have.

**17.5 The success criteria are unfalsifiable.** "materially simpler than the current multi-service architecture" — this is a judgment, not a test. The criteria should be measurable (e.g., "fewer than N processes, fewer than M ports, zero HTTP boundaries within ARGUS").

---

## 18. Concrete questions the plan must answer before implementation

1. For each "complexity" line item in the current plan, is the complexity from deployment, API design, or domain?
2. What mechanism enforces package boundaries in the single-process ARGUS?
3. What is the actual work semantics the current POC depends on? Which are preserved and which are dropped?
4. What is the current DB topology, and what is the target?
5. How does the collapsed process distinguish agent calls from human calls?
6. How is MCP mode isolated from API mode at the code level?
7. What is the verifier's runtime model and dependency footprint?
8. What is the domain pack's form, ownership, and validation regime?
9. Is the pivot in-place, on a branch, or a new repo? Is there a rollback?
10. What is the disposal plan for current tests?
11. What does `argus reset` actually do?
12. What is the portability argument, concretely — which interfaces must the domain pack satisfy?

---

## 19. Bottom line

The pivot is the right move. The current architecture has drifted from the thesis, and the retrospective is honest about it. But the pivot is not as simple as it presents itself.

The strongest concern is that the pivot may relocate complexity rather than remove it. The single-process ARGUS will have fewer processes, ports, and HTTP calls. But it will still have cross-package coupling, reference validation, artifact handling, and auth propagation — just inside one binary. If those are not addressed at the design level, the pivot will produce a monolith with the same complexity as the current distributed system, but less visibility.

The second concern is that the pivot decides before it inspects. The plan-writer is asked to remove Coordinator and Conductor before determining what they actually do. The reconnaissance phase may reveal that some of their semantics are load-bearing, in which case the plan must either preserve them or accept their loss — and the prompt does not tell the plan-writer how to make that trade.

The third concern is that the plan-writing prompt is over-specified and under-prioritized. Eighteen sections and twelve phases is too many for a POC pivot. The prompt should be reduced to the load-bearing questions and the reconnaissance-first sequence.

**Recommendation:** Before writing the plan, answer one question: **What does the reconnaissance phase need to find for the pivot to proceed, and what would make the plan-writer recommend *against* the pivot?** If the reconnaissance cannot produce a "no" answer, it is not a real test. If it can, the plan should state what that "no" looks like. This is the same discipline the POC is supposed to enforce on research claims — applied to the pivot itself.


## z

# REVIEW OF THE PIVOT PLAN (background.md + prompt.md)
## The pivot is correct and overdue — but the brief has one dangerous silence and one underspecified keystone

**Executive verdict:**

$$
\boxed{\text{Adopt the pivot. It is the right call, three weeks late — and my own ARGUS review independently confirms it: of 15 findings, the pivot dissolves 5 outright and simplifies 4 more. But the brief never mentions authentication (the demo's crown jewel), and it under-specifies the Solvent kernel-extraction path (the difference between "simplify a POC" and "strangle the long-term authority-substrate thesis"). Fix those two and execute.}}
$$

---

## 1. WHAT THE PIVOT GETS RIGHT — CONFIRMED, NOT MERELY AGREED

The central diagnosis — **deployment separation was mistaken for evidence of conceptual separation** — is correct, and it is the same principle Solvent's own docs state (*"prefer structural truth over labels"*, §22) turned on the infrastructure itself. Three independent confirmations:

1. **My prior review's failure pattern matches.** The findings that consumed the most review attention (F5 dual-write saga, F6 RetractCascade coupling, F13 readiness orchestration, F3's cross-process artifact handoff) are all *accidental complexity of the topology*, not of the thesis. The brief's own signal (P0-2 artifact registries, P1-3 cross-service reference validation) is the same signature.
2. **The kill question is the right one** ("what architectural claim becomes untestable if Coordinator/Conductor stop being services?"). The honest answer, which the plan must state verbatim: **nothing architectural — only operational isolation (independent deployment, ownership, trust-scoping), which no POC needs.**
3. **The component verdict table is right**, including the two subtle calls most people get wrong: keeping the Domain Pack (it *is* the portability thesis, not deployment machinery) and keeping the epistemic model untouched while collapsing everything around it.

---

## 2. FINDINGS DISSOLVED BY THE PIVOT (dividends to bank explicitly)

| Prior finding | Disposition under one-process/one-DB |
|---|---|
| F5 — dual-write packet saga, idempotency | **Dissolved**: packet persistence = one DB transaction |
| F6 — RetractCascade service coupling | **Dissolved**: in-process follow-through with the same transaction boundary |
| F13 — readiness/PID/credential orchestration | **Dissolved** by `task dev` |
| F3-part — cross-process artifact handoff | **Dissolved**; verifier invoked in-process during packet validation |
| F12 — cross-scenario references across services | **Simplified** to same-DB foreign keys + scenario containment |
| F7 — partial multi-service failure | **Mostly dissolved** (one DB down = 503, never empty context); TRUNCATED≠COMPLETE still applies to projections |
| F9 — Coordinator rule accretion | **Transformed, not dissolved**: rules become declarative pack config; enforce at package level (see §3) |

One new dividend worth stating in the plan: **the promote-gate TOCTOU dies too.** "Debt empty → promote" becomes a single serializable transaction (CRDB gives this for free at single node), closing writeup §9's race class at the POC scale.

---

## 3. FINDINGS THAT SURVIVE THE PIVOT — THE BRIEF'S GAP

Package boundaries do not fix data-model and policy failures. The pivot plan must carry these forward or the demo will prove the wrong thing:

| Surviving finding | What the pivot plan must add |
|---|---|
| **F1 — agent-authored edges** (worst one) | Packet-borne `contradicts`/`derives` edges enter as `PROPOSED`; retraction legal only from `ACTIVE`; UI labels "proposed by agent ⟨id⟩, awaiting human confirmation." Demo steps 9→14 must route through this gate explicitly |
| **F2 — Trust UI auth/attribution** | **THE GAP: the brief mentions authentication zero times.** Same attack surface survives at `:8080` (DNS rebinding, CSRF, body-supplied `discharged_by`). Required: token-derived attribution only, per-action confirmation token, Origin/Host allowlist, step-up auth on discharge/retraction. Without this, the §8 demo proves nothing about authority — curl can discharge every debt |
| **F3-core — verifier trust chain** | Pre-registered script hash per obligation; declared numeric tolerance in the artifact schema (the G-1 checks are floating-point — two runs can disagree); env pinning; invocation by the application layer, not by agents. Clarify §13: **content hashing is not "cryptographic attestation infrastructure"** — don't let that constraint delete the trust chain |
| **F4 — overstated mutation invariant** | Restate precisely everywhere: *agents create proposals; humans discharge, confirm, promote, retract.* Acceptance tests: discharge-with-open-debt refused; promote-after-retract refused; failed mutation produces no success audit row (writeup §14's no-false-audit, made a test) |
| **F8 — agent freshness** | Container per agent, no shared mounts, MCP as the only channel; record image hash in packet provenance so independence is auditable post hoc |
| **F10 — evidence typing** | Becomes *more* important: the pack supplies evidence classes (BM-IST: `PROOF_OBLIGATION` vs `COMPUTATION` vs `CITATION`); the core supplies the extension mechanism and **never hardcodes types**. L1/L2 are human-discharge-only obligations — the verifier can only supply numeric sub-artifacts |

---

## 4. THE UNDERSPECIFIED KEYSTONE: THE SOLVENT KERNEL-EXTRACTION PATH

The brief says "preserve epistemic/authority semantics and canonical data model" — necessary, not sufficient. The ecosystem thesis (writeup §18/26: Solvent as the reusable authority substrate) lives or dies on one rule the plan must state mechanically:

> **Import-direction rule:** `internal/epistemic` is the innermost layer. It imports nothing from `application`, `work`, `api`, `mcp`, `ui`, or any domain pack. It owns its schema. Domain vocabulary is forbidden inside it (debt stays opaque — writeup §2.3, kernel-freeze §2.2).

Two decisions the plan must make explicit:

1. **Repo strategy:** import Solvent as a Go module from its own repository (preserves the standalone-product option; module boundary *enforces* the freeze) vs. copy into the monorepo (simpler POC; costs the extraction path). Decision criterion: does the epistemic package compile with zero ARGUS-specific imports? If the HTTP coupling is too deep to extract cheaply, copy now behind the same interface and extract later — but write the interface first, because the interface *is* the kernel.
2. **Mechanical boundary enforcement:** every §3 invariant maps to tooling, not comments — Go module separation plus `depguard`-style import lint plus one negative compile test per invariant (`work` cannot promote; `epistemic` cannot execute tasks; `verifier` cannot authorize; `decision` cannot invent truth). This is writeup §22 applied to code structure: the pivot's whole point is structural truth, so boundaries must be structural.

Note the pleasing symmetry: **the §9 portability test, applied to Solvent itself, is the kernel-extraction proof.** One test, two theses.

---

## 5. THE ORPHANED QUESTION — ANSWERED (background.md's cut-off prompt)

*Will domain-agnostic features carry over if the market repurposes the physics POC?* **Yes — if and only if the core never imports pack vocabulary.** Carries over unchanged: beliefs, evidence, debt, edges, promotion/retraction state machines, authority gates, the two-tool MCP boundary, UI projections, audit. Per-domain: vocabulary, retirement-rule config, evidence classes, the verifier implementation, scenario semantics. The proof is not prose — it is a **compile-time stub-pack test in CI**: the core builds against a null domain pack. That stub is a *test fixture*, not a second domain pack — clarify §13's constraint accordingly, or the most valuable portability evidence gets forbidden by accident. The failure mode to name and guard: retirement rules drifting into the core "just this once."

---

## 6. SPECIFIC COMMENTS ON THE BRIEF ITSELF

- **§6 (inspect real code):** right instinct; require the output format — a code-inventory table with the disposition column (KEEP/CONVERT/DELETE/DEFER) filled from the actual tree, and a count of "code that exists only because of the topology." That count is the pivot's ROI evidence for §16.
- **§2 (`argus verify` vs in-process verification):** both appear; state the division of labor explicitly — CLI subcommand = development/reproduction runs; in-process invocation = trust-chain-bound verification during packet validation. Same library underneath, or you will grow two divergent verification paths.
- **§8 demo:** steps 9–14 must be annotated with the F1 gate and F2 attribution requirements; step 8 requires pre-registered scripts + tolerance in the schema. Step 12's "where permitted by the BM-IST pack" is exactly where `PROOF_OBLIGATION` typing does its work.
- **§14:** the two answers should be written into the plan as commitments, not questions. Second answer, in full: *Solvent is worth preserving as a distinct conceptual subsystem because it is the only subsystem whose state is authoritative rather than operational — the authority model is the thesis; everything else is machinery — and because it is the candidate reusable product, hence the kernel discipline.*
- **§5 (CRDB):** a conscious exception to delete-first — keep it, single node, for schema reuse and migration-cost avoidance; note that swapping to embedded Postgres/SQLite remains a recorded fallback if operational overhead bites, not a POC work item.
- **§11 plan-format additions:** (19) security & adjudication path; (20) invariants→mechanical-enforcement map; (21) artifact/evidence schema; (22) auth model. The current 18 headings would produce a plan that is silent on exactly the layer the demo exists to demonstrate.
- **Phase design:** sound; one amendment — Phase 0's inventory drives everything, so gate Phase 2 (topology collapse) behind a completed Phase 3–5 interface skeleton where possible (strangler order: build inside, cut over, delete last). Phase 11 deletion should include the inter-service HTTP clients and their tests — deletion of *tests that encode deployment mechanics* (brief §6.9) is where most of the simplification is felt.

---

## 7. COMPLEXITY LEDGER FOR §16 (fill with repo-verified counts)

| Dimension | Current | Pivoted |
|---|---|---|
| Processes | ~7 (+ OpenCode) | 1 (+ OpenCode) |
| Ports | ~6 | 1 (HTTP+UI; MCP is stdio) |
| Databases | 1 | 1 |
| Repositories | 2–3 | 1 (+ Solvent module import) |
| Credential sets | N | 1 |
| Startup steps | orchestrated readiness chain | `task dev` |
| Internal HTTP calls in demo critical path | many | **0** |
| Cross-process artifact handoff | yes | none |

---

## 8. BM-IST PROGRAM CONTINUITY (why this pivot serves the research, not just the codebase)

The pivot does not change the research program's needs — it makes them cheaper to satisfy. The verifier's first pre-registered scripts remain the **G-1 self-audit set** (shear-counterexample reproduction; G2 Madelung cancellation), now with an in-process trust chain instead of a cross-process bootstrap — strictly simpler. The demo slice (G0/L1/L2) aligns with implementation-plan Phases 0–1, whose honest debts are exactly "write the note" — i.e., `PROOF_OBLIGATION` evidence discharged by a human, which is the demo Branch B in miniature. And §37.9 (external AI output auto-tagged ASSUMED pending audit) is now structurally enforced rather than procedural: agent packets are proposals, human adjudication is the audit pass, and the named-object/citation checker slots in as a packet-validation step producing typed evidence, never authority.

**Bottom line:** execute the pivot as briefed, with three amendments written into the plan before Phase 0 — the auth/adjudication section (F2), the edge PROPOSED/ACTIVE gate (F1) and evidence typing (F10) as acceptance criteria, and the kernel-extraction import rule with the stub-pack portability test. The pivot's own kill question answers itself; the plan's job is to make sure that what survives the collapse is not just smaller, but still *demonstrably* the thesis — and the thesis was never the number of services.