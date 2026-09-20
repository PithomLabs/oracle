## chatgpt

I reviewed the uploaded plan against the decisions we just locked. It correctly incorporates the major choices: **ARGUS**, separate Trust UI, symbolic algebra in Go, JSON Schema + Go validation, and **LLM WORK + deterministic ADVERSARIAL stub**. 

Before implementation, I would make **four corrections**:

1. **ARGUS naming is missing from the plan.**
   The project should explicitly be named **ARGUS — Trust Verification POC**, with Oracle, Sphinx, Styx, Ferryman, and Chronicler as the UX vocabulary.

2. **Trust UI repository boundary is ambiguous.**
   Phase 0 says `oracle/trust-ui/` as a sibling repo, while Phase 9 presents `trust-ui/` as a directory structure. Pick one. Based on the locked architecture, I recommend a **separate `trust-ui` repository**.

3. **Phase 4 says “Symbolic algebra in Go” but then proposes building a miniature CAS from scratch.**
   That is substantially larger than the intended POC. The plan should either explicitly authorize a small existing Go symbolic-algebra dependency or sharply constrain the custom CAS to only the exact operations needed for the Fisher-rigidity demonstration. The current scope includes expressions, differentiation, integration, simplification, and several derivation rules. 

4. **Phase 3 needs a packet registry mechanism.**
   It requires `pack_ref` to resolve to a registered Domain Pack, but the plan does not define where that registry lives or how resolution works. 

Everything else is structurally aligned, including the 15-phase sequence and the explicit rule that Conductor remains unmodified. 

**Recommendation: make those four edits, then freeze this as the implementation baseline.**



## gemini

The implementation plan establishes a clean conceptual separation for domain packs, but it contains critical security gaps, boundary ambiguities, missing fail-closed mechanisms, and a truncated phase specification.

## 1. Plan Completeness & Truncation Risks

* **Truncated Phase Specifications (Phases 10–15):** The phase dependency graph outlines 16 phases (Phases 0–15), but the text cuts off midway through Phase 9 (`templates/challenge_de...`). Technical specifications, test suites, acceptance criteria, and rollback strategies for Work/Adversarial Runs (Phase 10), Consequential State Demo (Phase 11), Domain-Neutrality Proof (Phase 12), and Evidence Package Acceptance (Phases 13–15) are entirely missing.
* **Premature UI Information Architecture (Phase 8):** Phase 8 lists Phase 2 as its sole dependency. However, the IA defines surfaces like Sphinx (action gates), Decision Rooms, and Refusal Archives that depend directly on Solvent authority target endpoints and Coordinator decision APIs defined in Phases 6 and 7.

## 2. Epistemic & Security Fail-Open Vulnerabilities

* **Unverified Proof Artifact Forgery:** In Phase 4, the verifier computes a `content_sha256` for evidence. If an adversarial agent generates a synthetic `VerificationArtifact` JSON containing `result: "confirmed"` and embeds it directly into an EBP packet, the Coordinator (which explicitly executes zero reasoning or evidence validation) will ingest it as valid `reproducible_artifact` evidence. Without isolated sandbox re-execution or cryptographic signature verification by the Coordinator, forged artifacts will automatically discharge debt.
* **Trivial Idempotency Bypasses:** Phase 6 keys idempotency strictly to `packet_id` (UUID). A rogue or malfunctioning agent loop can mutate minor metadata (such as adding whitespace or updating a timestamp) to receive a fresh `packet_id`, bypassing idempotency and flooding Solvent with redundant unpromoted beliefs and floating obligations.
* **Decision Room Re-submission Denial of Service:** Phase 7 defines human refusal as logging without state mutation. If an operator rejects an epistemic transition, nothing stops an agent from immediately re-submitting the exact same payload under new packet IDs, spamming human reviewers.

## 3. Boundary Inflation & State Inconsistency

* **Coordinator Role Inflation:** Phase 2 places `retirement_rules` in `pack.json`, while Phase 6 tasks the Coordinator with issuing `POST /v1/beliefs/{id}/debt/retire` upon ingesting evidence. Evaluating evidence-to-debt mapping logic transforms the Coordinator into a domain policy engine. If it blindly trusts evidence tags, safety fails; if it validates evidence claims, its non-reasoning boundary invariant is violated.
* **Split-Brain Between Solvent and Conductor:** Phase 6 enforces "Solvent-before-Conductor ordering" across independent HTTP APIs without a saga pattern or two-phase commit. If Solvent successfully registers a belief but Conductor fails during task creation (e.g., due to a 500 error or network timeout), Solvent will retain active debt for an operational task that does not exist in Conductor.

## 4. Technical Feasibility & Testing Blind Spots

* **Go Symbolic Algebra Engine Underestimation:** Phase 4 proposes a custom Go CAS to verify Fisher-rigidity derivations involving non-linear Lagrangian calculus and variational quantum potential derivatives. A lightweight tree-simplification engine will encounter combinatorial explosion or issue false refutations on valid non-trivial identities without a full Gröbner basis or Computer Algebra System.
* **Concurrency Blind Spots:** Test suites across all phases lack coverage for concurrent packet compilation targeting the same scenario, race conditions between human adjudication and automated debt retirement, or concurrent retraction cascades.

To secure this pipeline before implementation, the plan must require cryptographic validation of execution artifacts, implement a saga compensation mechanism for multi-repo writes, restrict idempotency hashing to canonical payload contents rather than client-supplied UUIDs, and detail the missing phases (10–15).


## Z

Final-cut plan review. Disposition first, because the credit is real: **pre-registered acceptance criteria exist at both the run (Phase 10) and program (Phase 15 matrix) level — the oldest standing finding in this project, dating to the Phase 1 unstated-PASS, is now structurally resolved.** Per-phase rollback/escape, idempotency + Solvent-before-Conductor ordering, mandatory `initial_debt` with its negative test, neutrality scan + adversarial self-review phases, Sphinx-as-view via the verify-endpoint projection, Riddle Registry correctly absent, and a non-goals list carrying every deferral faithfully — this is the most disciplined plan of the series. Your message cut off at "Please help me:" again (and the document truncates mid-risk-table); proceeding on the standing request.

Two Criticals, then the rest.

## Critical

**C1 — Edges target beliefs by claim text. That is semantic matching in the coordinator — the exact thing "compiler, not normalizer" prohibits.** Phase 3's `edges[].target_claim: "string — claim text of related belief"` means the coordinator must match agent-authored free text against existing beliefs to resolve a contradicts/derives edge. That is semantic repair — the hidden adjudication this architecture banned in the packet round (the C1 finding then: local references valid intra-packet, canonical IDs for existing state, never prose). Claim-text matching fails on paraphrase, drifts silently, and makes the coordinator decide when two statements are "the same claim." Fix: restore the two-form reference — `local_id` for in-packet objects (and give packet objects local_ids; currently evidence can't even reference its own belief except implicitly), canonical `belief:<uuid>` from `GET /state` for existing objects, resolved deterministically. Free-text targeting is rejected at validation.

**C2 — Phase 4's CAS scope is research-grade, and the plan doesn't survive contact with it.** Building an expression-tree CAS with `Diff`, `Integral`, and simplification rules capable of *deriving* the Fisher-rigidity result via Euler–Lagrange matching on *functionals* is a multi-month project — variational calculus with functional derivatives and integration by parts is where toy CASes die. The truncated risk table apparently knows this. Honest options: (a) **narrow the verifier to what's actually checkable** — the four steps verified independently in this conversation's record (the identity Δ√ρ/√ρ = ½Δρ/ρ − ¼|∇ρ|²/ρ²; coefficient matching forcing a(ρ)=κ²/8mρ; the automatic consistency of the |∇ρ|² term; b′(ρ)=0) are each *simple symbolic manipulations* — polynomial identities and coefficient comparison, well within a small deterministic Go checker; (b) pin an external CAS (version + hash + recorded output) and treat its output as evidence provenance, trading self-containment for capability. Either is defensible; the plan as written will either blow Phase 4 or ship a verifier whose "confirmed" covers less than its name implies. The artifact already has per-step `Verified` fields — use them to record precisely which steps were machine-checked versus human-attested. That honesty *is* the POC's credibility.

## High

**H1 — The pack declares retirement rules; nothing enforces them.** `retirement_rules` map each debt item to an evidence class — and then Phase 7 lets a human `RETIRE_DEBT` any item via direct API call, with no evidence-class check anywhere in Phase 6 or 7. So the run will retire `needObstruction` without any obstruction-construction artifact existing — the rules are decorative. This is the v0.4 retirement-gate question resurfacing after its deferral. Decide and record: either the coordinator validates "offered evidence class matches the pack's rule for that debt item" before executing retirement (mechanical, cheap, and makes the pack real), or the rules are explicitly advisory in this POC and retirement is human-discretionary — which is defensible for a POC but must be stated, because "human retires remaining debt items" in the Phase 10 sequence currently *looks* like rule-following while bypassing the declared rules.

**H2 — The pack validator should enforce the structural human-gate floor.** `human_gated_transitions` is pack data — meaning a future pack could omit faithfulness. The frozen principle: *which transitions may be consequential is domain vocabulary; which require a human is structural.* The Phase 2 validator should require the structural minimum (promotion, final-truth-claim, faithfulness-review-class transitions) to be present and human-gated, with packs adding beyond it. One validation rule.

**H3 — UI: eight surfaces pre-run vs. the three-screen extraction principle.** The C1 finding from the UI review (surfaces earned by usage) is only partially honored — the plan's argument (Phase 11 needs transitions visible live) is legitimate, but eight full surfaces before first use re-imports the pre-pivot pattern. Middle path that preserves both: Phase 9.1 builds Oracle + Decision + Refusals (+ Briefing as aggregate), the run proceeds, Phase 11 expands to Challenge/Passage/Chronicle/Sphinx-detail *from what the run shows humans actually consulting*. Or accept the deviation and record why. Related concrete gap: **Oracle's acceptance criteria promise "history," and Passage promises per-debt receipts — but the projection table maps neither** (no belief-transition log exists; receipts would need decision-record joins). Map them to `audit_activity` + coordinator decision records explicitly, or descope those fields to "recorded from adoption forward."

**H4 — The freeze reconciliation record: fourth request, now blocking-adjacent.** Phase 1 records the kernel hash but the strategy-level record still doesn't exist, and this plan *adds* supersessions: Conductor UI demoted as the system's face, packet v1 evolved from the frozen research-POC design (edge model changed), Workflow Design v1.0's UI requirement superseded. One document: which freezes bind, which are amended, which deferred — or "frozen" stops meaning anything at exactly the moment the project has the most frozen artifacts.

**H5 — Incomplete-attack coverage detection is never exercised.** The stub produces all four attack types always, so the `INCOMPLETE_ATTACK` validation path has no test. Add the negative case: a seeded under-attacking packet → flag fires. Otherwise Phase 13's coverage claim is tested only in its passing direction.

## Medium

- **REOPEN = new `EnterBelief` with same claim** — creates an unrelated sibling, severing lineage from the retracted belief. Require a `derived_from` edge to the retracted origin (the edge machinery exists) or accept and document the lineage break.
- **Task-creation scope governance dropped**: Phase 6 compiles every packet task; the prior rule (auto-create only within approved scope/debt-derived; otherwise human proposal) isn't restated. One sentence in Phase 6.
- **Human principal mechanism unspecified**: `middleware.go` says "Auth, logging" in one line; Phase 7 requires a human principal but nothing establishes one. POC-scale answer: config-declared operator identity injected into decision records, stated in the plan — otherwise the test "decision without principal → rejected" tests against nothing.
- **Phase 12 scan false positives**: "phase" and "density" are common English/CS words; a naive grep fails. Context-window matching or an explicit allowlist; otherwise the neutrality report arrives pre-discredited.
- **Phase 13's "no false negatives"** is unfalsifiable as acceptance. Reword: "all probed surfaces PASS; methodology and coverage documented."
- **Corpus manifest example shows one artifact** — the run depends on the Qwen PDF too; include it with its provenance class set honestly (`agent_derived` artifact, not source evidence — the M5 finding from the packet round).

## Low

grep-for-"reason" as the no-reasoning check is theater; the architectural check is no LLM dependency in `go.mod` — keep both, weight the latter. Repository sprawl: `oracle/` now spawns eight subpackages plus a sibling repo — one line on repo strategy. Risk table truncation noted; if the omitted rows cover the CAS risk, C2 is at least on the plan's radar. Evidence-without-belief-link is fine for single-belief packets — note the constraint so multi-belief packets don't get improvised later.

## Routed order

C1 (reference forms) and C2 (verifier scope decision) before Phase 3/4 code — both reshape schemas. H1 (retirement enforcement decision) before Phase 7. H2 (validator floor) into Phase 2. H3 (staged UI or recorded deviation) before Phase 9. H4 (reconciliation record) in parallel with Phase 1 — it's the one deliverable Phase 1 is missing. H5 and mediums as test/spec lines.

The structural achievement here is worth naming: **the project's review history is now encoded as plan structure** — pre-registration, rollback per phase, neutrality as falsifiable claim, refusal and unknown as first-class, agents as packet-producers only. The findings above are the residue of one genuine design decision (claim-text edges) and one scope miscalculation (CAS), plus the standing process debt. Fix the edges, right-size the verifier, write the reconciliation record — and this plan is executable as written.
