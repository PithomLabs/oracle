# Implementation Plan: Physics Verifier POC + Domain-Agnostic Trust UI

**Version:** 1.0
**Date:** 2026-09-15
**Status:** PLAN — NOT IMPLEMENTED

---

## Executive Decision

Prove that BM-IST can run as the **first Domain Pack** on a domain-agnostic trust-verification substrate. The generic trust substrate governs BM-IST research without physics-specific logic entering Solvent or Conductor.

**Locked architecture:**

```
    Domain Pack (BM-IST)
        ↓
    Trust Verifier / Coordinator
        ↓
    Solvent = epistemic + authority state
    Conductor = operational workflow
        ↓
    Executor / Ferryman
        ↓
    External SOR = effect truth
```

**Technology decisions (locked):**

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Kernel freeze | Plan 11.1 fully applied, new freeze hash | Prerequisite satisfied |
| Trust UI location | New sibling repo `oracle/trust-ui/` | Clean separation, no Conductor coupling |
| Trust UI tech | Go HTTP + server-rendered HTML | Consistent with Solvent/Conductor, no JS framework |
| Physics Verifier | Symbolic algebra in Go | Deterministic, native, reproducible |
| EBP Packet Contract | JSON schema + Go validation code | Canonical interchange + typed validation |
| End-to-end agents | LLM WORK + deterministic ADVERSARIAL stub | Real research + reproducible attacks |

---

## Phase Dependency Graph

```
Phase 0 (Reconnaissance)
    ↓
Phase 1 (Freeze Reconciliation)
    ↓
Phase 2 (Domain Pack Spec) ←──── Phase 5 (Corpus Manifest)
    ↓                              ↑
Phase 3 (EBP Packet Contract)
    ↓
Phase 4 (Physics Verifier Runner) ← Phase 5
    ↓
Phase 6 (Deterministic Coordinator) ← Phases 2, 3, 4
    ↓
Phase 7 (Human Adjudication) ← Phase 6
    ↓
Phase 8 (Trust UI Information Architecture) ← Phase 2
    ↓
Phase 9 (Trust UI Implementation) ← Phases 6, 7, 8
    ↓
Phase 10 (Work + Adversarial Run) ← Phases 6, 7, 9
    ↓
Phase 11 (Consequential State Demo) ← Phases 9, 10
    ↓
Phase 12 (Domain-Neutrality Proof) ← Phase 11
    ↓
Phase 13 (Adversarial System Review) ← Phase 12
    ↓
Phase 14 (POC Evidence Package) ← Phase 13
    ↓
Phase 15 (Final Acceptance) ← Phase 14
```

**Parallel opportunities:**
- Phases 2+5 can run in parallel (Domain Pack spec + Corpus manifest)
- Phases 4+5 can overlap (Verifier runner starts while corpus is being pinned)
- Phase 8 (UI info architecture) can start as soon as Phase 2 lands
- Phase 9 (UI implementation) can begin partial work during Phase 6

---

## Phase 0 — Repository Reconnaissance

**Objective:** Establish the exact current state of all repositories as the baseline for planning.

**Preconditions:** None.

**Repository areas inspected:**

| Area | Repo | Key Files | Status |
|------|------|-----------|--------|
| Solvent kernel | solvent-main | `kernel/kernel.go`, `kernel/authority.go` | Frozen at Plan 11.1 |
| Solvent schema | solvent-main | `db/001_schema.sql` through `db/010_debt_opaque.sql` | 10 migrations, debt opaque |
| Solvent API | solvent-main | `api/types.go`, `api/belief.go`, `api/target.go` | 25+ endpoints |
| Solvent services | solvent-main | `service/ledger/`, `service/authority/`, `service/executor/` | Full authority lifecycle |
| Wizard domain | solvent-main | `internal/belief/debt.go`, `internal/belief/mapping.go` | Deployment-review vocabulary |
| Conductor core | conductor | `internal/domain/task.go`, `internal/store/task_repo.go` | Domain-agnostic workflow |
| Conductor governance | conductor | `internal/governance/`, `internal/adapter/solvent/` | Solvent adapter exists |
| Conductor UI | conductor | `internal/web/server.go`, `internal/web/templates/layout.html` | Kanban board, read-only |
| Reference loop | oracle/reference-loop | `main.go`, `agent/agent.go`, `contract.json` | Phase 2 complete |
| EBP seed | oracle/docs | `BMIST_POC_SEED.md` | BM-IST claim documented |
| Existing plans | oracle/plan | `roadmap.md`, `myth_symbolism.md` | Architecture locked |

**Key findings for planning:**

1. Solvent has no physics/domain vocabulary in kernel — `db/010_debt_opaque.sql` enforces opacity.
2. Conductor's `GovernanceReader` interface already provides read-only Solvent integration.
3. The reference-loop already demonstrates the full authority lifecycle (belief → promote → target → authorize → execute).
4. BM-IST EBP debt vocabulary is documented in `BMIST_POC_SEED.md` (needMap, needInvariant, needToyCheck, needNullModel, needObstruction, needFaithfulnessReview).
5. No Coordinator code exists in any repository — this is the primary new component.
6. No Trust UI exists — the existing Conductor UI is Kanban-only.

**Acceptance criteria:** This phase is informational. No deliverable required beyond this plan.

**Rollback:** N/A — read-only reconnaissance.

---

## Phase 1 — Freeze Reconciliation

**Objective:** Establish an immutable baseline from which all subsequent work proceeds.

**Preconditions:** Plan 11.1 fully applied and committed in solvent-main.

**Repository areas to inspect:**

| File | Check |
|------|-------|
| `solvent-main/db/010_debt_opaque.sql` | EXISTS, applied in all bootstrap paths |
| `solvent-main/internal/belief/debt.go` | `wizardDebt` is single source of truth |
| `solvent-main/kernel/kernel.go` | `FullDebt` constant REMOVED |
| `solvent-main/api/types.go` | `EnterBeliefRequest.Debt` field EXISTS |
| `solvent-main/cmd/solvent-mcp/tools.go` | Vocabulary guard REMOVED |

**Files/components likely to change:** None — this is verification only.

**Data/API impact:** None.

**Tests:**
```bash
cd /home/chaschel/Documents/go/solvent-main
git rev-parse HEAD              # Record freeze hash
git status --short              # Must be clean
go test ./...                   # Full suite passes
grep -r "FullDebt" kernel/      # Must return empty
```

**Acceptance criteria:**
1. Freeze hash recorded in this plan (see below).
2. Working tree clean.
3. All tests pass.
4. No `FullDebt` references in kernel package.
5. `debt SET DEFAULT ARRAY[]::TEXT[]` verified in migration chain.

**Rollback/escape:** If tests fail or working tree is dirty, STOP. Do not proceed to Phase 2 until clean.

**Dependencies:** Plan 11.1 completion (confirmed).

---

## Phase 2 — Domain Pack Specification

**Objective:** Define a declarative Physics Pack specification that Solvent/Conductor never interpret but that the Coordinator uses to compile packets into generic Solvent API calls.

**Preconditions:** Phase 1 complete (freeze established).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `oracle/docs/BMIST_POC_SEED.md` | BM-IST claim, debt vocabulary |
| `oracle/plan/roadmap.md` | Domain Pack abstraction design |
| `solvent-main/internal/belief/debt.go` | Wizard debt pattern to follow |
| `solvent-main/kernel/kernel.go` | Claim types, debt semantics |
| `solvent-main/db/001_schema.sql` | Schema constraints the Pack must respect |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/domain-pack/bmist/v1/pack.json` | Domain Pack descriptor (JSON) |
| `oracle/domain-pack/bmist/v1/types.go` | Go types mirroring the JSON schema |
| `oracle/domain-pack/bmist/v1/validate.go` | Pack conformance validator |
| `oracle/domain-pack/bmist/v1/validate_test.go` | Negative tests |
| `oracle/domain-pack/README.md` | Pack contract explanation |

**Domain Pack Descriptor (`pack.json`) structure:**

```json
{
  "pack_id": "bmist-v1",
  "version": "1.0.0",
  "name": "BM-IST Physics Verification",
  "description": "Fisher-rigidity derivation verification domain",
  "claim_types": ["derived", "accommodated", "postulated"],
  "evidence_classes": [
    "reproducible_artifact",
    "operator_asserted"
  ],
  "debt_vocabulary": [
    "needMap",
    "needInvariant",
    "needToyCheck",
    "needNullModel",
    "needObstruction",
    "needFaithfulnessReview"
  ],
  "initial_debt": [
    "needMap",
    "needInvariant",
    "needToyCheck",
    "needNullModel",
    "needObstruction",
    "needFaithfulnessReview"
  ],
  "retirement_rules": {
    "needMap": ["reproducible_artifact:map_check"],
    "needInvariant": ["reproducible_artifact:invariant_check"],
    "needToyCheck": ["reproducible_artifact:toy_model_check"],
    "needNullModel": ["operator_asserted:scope_clarification"],
    "needObstruction": ["reproducible_artifact:obstruction_construction"],
    "needFaithfulnessReview": ["operator_asserted:faithfulness_review"]
  },
  "falsifiers": [
    "counterexample",
    "contradiction",
    "missing_evidence",
    "alternative_explanation"
  ],
  "consequential_actions": [
    {
      "name": "publish_claim",
      "requires": "promoted",
      "gates": ["faithfulness_review"]
    }
  ],
  "human_gated_transitions": [
    "faithfulness_review",
    "scope_clarification",
    "obstruction_assessment"
  ]
}
```

**Data/API impact:** None on Solvent/Conductor. New package in oracle.

**Tests:**
- Valid pack descriptor loads and validates.
- Invalid pack (missing required fields) rejected.
- Pack with unknown debt vocabulary rejected by pack-level validation.
- Pack's initial_debt matches BM-IST POC seed (6 items).

**Acceptance criteria:**
1. `pack.json` is valid JSON Schema.
2. Go types compile and match JSON.
3. Validator rejects malformed packs.
4. BM-IST debt vocabulary matches `BMIST_POC_SEED.md` exactly.
5. No physics vocabulary in Solvent or Conductor repositories.

**Rollback/escape:** Delete `oracle/domain-pack/` directory. No external impact.

**Dependencies:** Phase 1 (freeze baseline established).

---

## Phase 3 — EBP Research Packet Contract

**Objective:** Define the EBP Research Packet v1 as a JSON schema with Go validation code. The packet is the unit of exchange between agents and the Coordinator.

**Preconditions:** Phase 2 complete (Domain Pack spec defines the vocabulary the packet uses).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `oracle/reference-loop/agent/agent.go` | Current agent→Solvent interaction pattern |
| `oracle/reference-loop/contract.json` | Existing operation contract format |
| `solvent-main/api/types.go` | Solvent API request/response types |
| `oracle/domain-pack/bmist/v1/pack.json` | Pack vocabulary the packet must reference |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/packet/v1/schema.json` | JSON Schema for Research Packet v1 |
| `oracle/packet/v1/types.go` | Go types |
| `oracle/packet/v1/validate.go` | Validation logic |
| `oracle/packet/v1/validate_test.go` | Positive + negative tests |
| `oracle/packet/v1/examples/valid_packet.json` | Example valid packet |
| `oracle/packet/v1/examples/malformed_*.json` | Negative examples |

**Packet v1 structure:**

```json
{
  "packet_id": "uuid",
  "pack_ref": "bmist-v1",
  "scenario_id": "uuid",
  "belief": {
    "claim": "string — the research claim",
    "claim_type": "derived|accommodated|postulated",
    "debt": ["string — optional additional debt items beyond initial"]
  },
  "evidence": [
    {
      "provenance_class": "reproducible_artifact|operator_asserted",
      "source_url": "string",
      "content_sha256": "string — hex, required",
      "artifact_type": "string — domain-specific, opaque to substrate"
    }
  ],
  "edges": [
    {
      "kind": "derives|contradicts",
      "target_claim": "string — claim text of related belief"
    }
  ],
  "tasks": [
    {
      "title": "string",
      "description": "string",
      "governance_ref": "string — optional"
    }
  ],
  "local_references": {
    "verifier_version": "string",
    "corpus_manifest_ref": "string"
  },
  "ownership": {
    "work_agent_id": "string",
    "adversarial_agent_id": "string",
    "human_reviewer_id": "string"
  }
}
```

**Validation rules:**
1. `packet_id` is required, UUID format.
2. `pack_ref` must resolve to a registered Domain Pack.
3. `belief.claim` is required, non-empty.
4. `belief.claim_type` must be in pack's `claim_types`.
5. `evidence[].content_sha256` is required, non-empty.
6. All `debt` items must be in pack's `debt_vocabulary` OR be empty (Coordinator unions `initial_debt`).
7. `edges[].kind` must be `derives` or `contradicts`.

**Data/API impact:** None on Solvent/Conductor. New package in oracle.

**Tests:**
- Valid packet passes validation.
- Missing `packet_id` → rejected.
- Empty `belief.claim` → rejected.
- Unknown `claim_type` → rejected.
- Empty `content_sha256` → rejected.
- Debt item not in pack vocabulary → rejected at packet level (Coordinator adds initial_debt).
- Packet with no evidence → passes (evidence can be added later).

**Acceptance criteria:**
1. JSON Schema published and valid.
2. Go types compile.
3. All positive and negative tests pass.
4. Packet structure maps cleanly to Solvent API calls (documented in comments).

**Rollback/escape:** Delete `oracle/packet/` directory. No external impact.

**Dependencies:** Phase 2 (Domain Pack defines vocabulary).

---

## Phase 4 — Physics Verifier Runner

**Objective:** Implement a deterministic Go-based symbolic algebra verifier for the Fisher-rigidity derivation. Produces structured verification artifacts with no LLM authority.

**Preconditions:** Phases 2+3 complete (pack and packet contracts defined).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `oracle/docs/BMIST_POC_SEED.md` | Exact claim to verify |
| `oracle/background/physics_verifier_solvent_discussion.md` | Prior discussion of verification approach |
| `oracle/reference-loop/evidence/types.go` | Existing evidence artifact format |
| `solvent-main/db/001_schema.sql` | `reproducible_artifact` provenance class |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/verifier/physics/v1/verifier.go` | Fisher-rigidity symbolic verifier |
| `oracle/verifier/physics/v1/verifier_test.go` | Deterministic output tests |
| `oracle/verifier/physics/v1/artifact.go` | Verification artifact types |
| `oracle/verifier/physics/v1/artifact_test.go` | Artifact hash determinism tests |
| `oracle/verifier/physics/v1/cas/expr.go` | Symbolic expression tree |
| `oracle/verifier/physics/v1/cas/simplify.go` | Simplification rules |
| `oracle/verifier/physics/v1/cas/cas_test.go` | CAS correctness tests |

**Verification artifact structure:**

```go
type VerificationArtifact struct {
    VerifierVersion string    `json:"verifier_version"`
    VerifierHash    string    `json:"verifier_hash"`    // SHA256 of verifier binary/source
    ClaimHash       string    `json:"claim_hash"`       // SHA256 of claim text
    InputHash       string    `json:"input_hash"`       // SHA256 of all inputs
    ArtifactHash    string    `json:"artifact_hash"`    // SHA256 of this artifact (excl. self)
    Timestamp       time.Time `json:"timestamp"`
    Steps           []Step    `json:"steps"`
    Result          string    `json:"result"`           // "confirmed" | "refuted" | "inconclusive"
    EvidenceRef     string    `json:"evidence_ref"`     // content_sha256 for Solvent evidence
}

type Step struct {
    Index       int    `json:"index"`
    Description string `json:"description"`
    Input       string `json:"input"`
    Rule        string `json:"rule"`        // e.g. "euler_lagrange_matching"
    Output      string `json:"output"`
    Verified    bool   `json:"verified"`
}
```

**CAS implementation scope (POC):**
- Symbolic expressions: `Expr = Const(float64) | Var(string) | Add(Expr,Expr) | Mul(Expr,Expr) | Pow(Expr,Expr) | Diff(Expr,Var) | Integral(Expr,Var)`
- Simplification: constant folding, power rules, product rule
- The Fisher-rigidity derivation steps:
  1. Start from Lagrangian L = T - V + E_q[ρ]
  2. Apply Euler-Lagrange to ψ = √ρ · exp(iS/κ)
  3. Derive that E_q must be local, ρ-only, C², rotationally invariant
  4. Show forced form: E_q = (κ²/8m) ∫(|∇ρ|²/ρ) dx
  5. Verify coefficient consistency

**Data/API impact:** None on Solvent/Conductor. New package in oracle.

**Tests:**
- CAS operations are deterministic (same inputs → same outputs, verified by hash).
- Verifier produces identical artifact for identical inputs.
- Verifier artifact hash is reproducible across runs.
- Known-good claim produces `result: "confirmed"`.
- Known-bad claim produces `result: "refuted"`.
- Edge cases: division by zero, empty expression, malformed input.

**Acceptance criteria:**
1. Verifier runs without LLM calls.
2. Artifact is deterministic (hash reproducible across 100 runs).
3. Fisher-rigidity derivation produces `confirmed` result.
4. Each step has verifiable input/output/rule.
5. Artifact contains all metadata fields.
6. `content_sha256` can be used as Solvent evidence `content_sha256`.

**Rollback/escape:** Delete `oracle/verifier/` directory. No external impact.

**Dependencies:** Phases 2+3 (pack and packet define the claim/vocabulary).

---

## Phase 5 — Corpus Manifest

**Objective:** Hash-pin the BM-IST research corpus so every piece of evidence has a stable, verifiable locator.

**Preconditions:** Phase 2 complete (pack defines evidence classes).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `oracle/docs/BMIST_POC_SEED.md` | Primary BM-IST artifact |
| `oracle/background/` | Research materials |
| `oracle/docs/` | Writeups, reviews |
| `oracle/reference-loop/PHASE2_EVIDENCE.md` | Prior evidence format |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/corpus/manifest.json` | Corpus manifest with hashes |
| `oracle/corpus/manifest.go` | Go loader + validator |
| `oracle/corpus/manifest_test.go` | Hash verification tests |

**Manifest structure:**

```json
{
  "manifest_version": "1.0.0",
  "created_at": "2026-09-15T00:00:00Z",
  "corpus_id": "bmist-poc-seed-v1",
  "artifacts": [
    {
      "id": "bmist-seed-001",
      "path": "docs/BMIST_POC_SEED.md",
      "sha256": "hex...",
      "provenance_class": "operator_asserted",
      "description": "BM-IST POC seed claim and debt status",
      "stable_locator": "oracle/docs/BMIST_POC_SEED.md"
    }
  ]
}
```

**Data/API impact:** None. New package in oracle.

**Tests:**
- Manifest loads and validates.
- All referenced files exist and hashes match.
- Missing file detected.
- Hash mismatch detected.
- Unknown artifact in manifest → warning.

**Acceptance criteria:**
1. All BM-IST corpus files are hash-pinned.
2. Manifest loads in Go.
3. Hash verification passes for all artifacts.
4. Stable locators are repository-relative paths.

**Rollback/escape:** Delete `oracle/corpus/` directory. No external impact.

**Dependencies:** Phase 2 (pack defines evidence classes).

---

## Phase 6 — Deterministic Coordinator

**Objective:** Implement the Coordinator as a deterministic, non-reasoning packet compiler. Transforms validated packets into Solvent API calls and Conductor task projections.

**Preconditions:** Phases 2, 3, 4 complete (pack, packet, verifier).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `solvent-main/api/types.go` | Solvent API request types for compilation |
| `solvent-main/api/api.go` | Solvent API route registration |
| `solvent-main/kernel/kernel.go` | Kernel functions the Coordinator calls |
| `conductor/internal/api/routes.go` | Conductor API for task creation |
| `conductor/internal/governance/reader.go` | GovernanceReader interface |
| `conductor/internal/adapter/solvent/adapter.go` | Solvent adapter pattern |
| `oracle/reference-loop/agent/agent.go` | Current agent→Solvent flow (to be replaced) |
| `oracle/reference-loop/solvent/client.go` | Solvent REST client |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/coordinator/coordinator.go` | Core compilation logic |
| `oracle/coordinator/compiler.go` | Packet → Solvent/Conductor compilation |
| `oracle/coordinator/compiler_test.go` | Compilation tests |
| `oracle/coordinator/validate.go` | Packet + debt validation |
| `oracle/coordinator/validate_test.go` | Validation tests |
| `oracle/coordinator/idempotency.go` | packet_id idempotency logic |
| `oracle/coordinator/idempotency_test.go` | Idempotency tests |
| `oracle/coordinator/human.go` | Human decision endpoint |
| `oracle/coordinator/human_test.go` | Human decision tests |
| `oracle/coordinator/client.go` | HTTP client for Solvent + Conductor |

**Coordinator responsibilities:**

1. **Packet validation:** Validate against packet schema, reject malformed.
2. **Debt attachment:** Union `pack.initial_debt` + `packet.belief.debt` → Solvent belief debt.
3. **Negative regression test:** If packet has no debt, Coordinator MUST still attach `initial_debt`. This prevents fail-open.
4. **Solvent compilation:**
   - `POST /v1/beliefs` with validated claim + compiled debt
   - `POST /v1/evidence` for each evidence item
   - `POST /v1/beliefs/{id}/debt/retire` for verified debt items
5. **Conductor compilation:**
   - `POST /v1/projects/{id}/tasks` for each task in packet
   - Set `governance_ref` to Solvent belief ID
6. **Idempotency:** `packet_id` used as idempotency key. Duplicate packets are no-ops.
7. **Ordering:** Solvent writes happen BEFORE Conductor writes (Solvent-before-Conductor ordering).
8. **Human decision endpoint:** `POST /v1/human/decisions` for consequential transitions.
9. **NO reasoning:** Coordinator does not infer truth, score confidence, auto-promote, auto-retire, or adjudicate.

**Coordinator responsibilities explicitly excluded:**
- Scientific reasoning
- Truth inference
- Confidence scoring
- Auto-promotion
- Auto-debt-retirement
- Autonomous adjudication
- Policy evaluation

**Data/API impact:**
- New coordinator service (standalone or in-process).
- Calls existing Solvent + Conductor APIs.
- No new Solvent/Conductor endpoints required.

**Tests:**
- Valid packet compiles to correct Solvent + Conductor API calls.
- Malformed packet rejected (missing packet_id, empty claim, etc.).
- Negative regression: packet with no debt → compiled belief still has initial_debt.
- Idempotency: same packet_id submitted twice → second is no-op.
- Ordering: Solvent belief exists before Conductor task references it.
- Human decision: transition recorded, not auto-executed.
- Conductor task has correct governance_ref.

**Acceptance criteria:**
1. Coordinator compiles packet → Solvent belief + evidence + Conductor task.
2. Mandatory initial_debt enforced (negative regression test passes).
3. packet_id idempotency works.
4. Solvent-before-Conductor ordering verified.
5. governance_ref correctly projected.
6. Human decision endpoint records (not executes).
7. Zero reasoning logic in Coordinator code (grep for confidence, score, infer, reason).

**Rollback/escape:** Delete `oracle/coordinator/` directory. No impact on Solvent/Conductor.

**Dependencies:** Phases 2, 3, 4 (pack, packet, verifier define what Coordinator compiles).

---

## Phase 7 — Human Adjudication

**Objective:** Implement the minimum consequential decision path for human review of epistemic transitions.

**Preconditions:** Phase 6 complete (Coordinator can compile packets).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `solvent-main/api/belief.go` | Existing Promote/Retract endpoints |
| `solvent-main/api/target.go` | Authority target lifecycle |
| `solvent-main/kernel/authority.go` | AuthorizeAndCreateIntent |
| `oracle/coordinator/human.go` | Human decision endpoint from Phase 6 |

**Files/components to create/modify:**

| Path | Content |
|------|---------|
| `oracle/coordinator/human.go` | Extend with decision types |
| `oracle/coordinator/human_test.go` | Decision workflow tests |

**Decision types:**

| Decision | Solvent API Call | Effect |
|----------|-----------------|--------|
| PROMOTE | `POST /v1/beliefs/{id}/promote` | Belief → promoted |
| RETIRE_DEBT | `POST /v1/beliefs/{id}/debt/retire` | Remove one debt item |
| RETRACT | `POST /v1/beliefs/{id}/retract` | Retract + cascade |
| REOPEN | `EnterBelief` with same claim | New entered belief |
| AUTHORIZE | `POST /v1/targets/{id}/approve` | Activate authority |
| REFUSE | Log refusal, no Solvent mutation | Record refusal |

**Every consequential transition must expose:**
- Current state (from Solvent ledger)
- Requested transition
- Evidence supporting the transition
- Open obligations
- Downstream consequences
- Acting human/principal
- Resulting Solvent state

**Data/API impact:** Uses existing Solvent endpoints. No new Solvent code.

**Tests:**
- PROMOTE on belief with debt → blocked (schema gate).
- PROMOTE on belief with no debt → success.
- RETRACT on promoted belief with live intent → cascade cancels intent first.
- REFUSE → refusal logged, no state change.
- Decision without human principal → rejected.

**Acceptance criteria:**
1. All 6 decision types work through existing Solvent API.
2. Every consequential transition logged with full context.
3. No silent state mutations.
4. Human principal required for consequential decisions.

**Rollback/escape:** Delete human decision code. Solvent endpoints remain unchanged.

**Dependencies:** Phase 6 (Coordinator compilation path exists).

---

## Phase 8 — Trust UI Information Architecture

**Objective:** Define the Trust UI's information model, navigation, state projections, domain-pack rendering contract, and read/write boundaries. No visual polish yet.

**Preconditions:** Phase 2 complete (Domain Pack spec defines what UI renders).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `conductor/internal/web/templates/layout.html` | Existing Conductor UI patterns |
| `solvent-main/demo/cloud/web/templates/` | Solvent demo UI patterns |
| `oracle/plan/myth_symbolism.md` | Mythology UX terminology |
| `oracle/plan/myth_ui.md` | UI metaphor design |
| `solvent-main/internal/view/view.go` | Read-only projections |

**Files/components to create:**

| Path | Content |
|------|---------|
| `trust-ui/docs/information-architecture.md` | IA specification |
| `trust-ui/docs/state-projections.md` | Solvent → UI state mapping |
| `trust-ui/docs/navigation.md` | Page/route structure |
| `trust-ui/docs/domain-pack-contract.md` | How UI renders pack-specific content |
| `trust-ui/docs/empty-states.md` | Unknown/empty/error semantics |
| `trust-ui/docs/read-write-boundaries.md` | What UI reads vs writes |

**Information Model:**

The Trust UI consumes ONLY:
1. Solvent API (`/v1/beliefs`, `/v1/evidence`, `/v1/ledger`, `/v1/activity`, etc.)
2. Coordinator API (packet status, compilation results)
3. Conductor API (task status, activity — read-only, non-authoritative)

The Trust UI NEVER:
1. Writes directly to Solvent database
2. Writes directly to Conductor database
3. Embeds physics/domain logic
4. Creates a second source of truth

**Navigation Structure:**

```
/ (BRIEFING SCROLL — landing page)
├── /oracle (ORACLE — epistemic query)
│   └── /oracle/:belief_id (belief detail)
├── /sphinx (SPHINX — action gate)
│   └── /sphinx/:riddle_id (riddle detail)
├── /passage (PASSAGE LEDGER — debt/obligations)
│   └── /passage/:belief_id (belief passage detail)
├── /challenge (CHALLENGE ROOM — adversarial verification)
│   └── /challenge/:belief_id (challenge detail)
├── /decision (DECISION ROOM — human consequence control)
│   └── /decision/:decision_id (decision detail)
├── /refusal (REFUSAL ARCHIVE — "the bones")
│   └── /refusal/:refusal_id (refusal detail)
└── /chronicle (CHRONICLE — historical activity)
```

**State Projections (Solvent → UI):**

| UI Concept | Solvent Source | Mapping |
|------------|---------------|---------|
| Oracle belief | `GET /v1/beliefs/:id` | Direct mapping |
| Belief evidence | `GET /v1/beliefs/:id/evidence` | Direct mapping |
| Belief debt | `belief.debt[]` | Render as passage requirements |
| Belief edges | `belief_edge` (via explain) | Derives/contradicts relations |
| Riddle (Sphinx) | `POST /v1/authorizations/verify` | Authorization gate result |
| Passage status | Debt items vs retirement evidence | Visual progress |
| Challenge | Adversarial agent output + falsifiers | Attack surface |
| Decision | `POST /v1/human/decisions` | Human action form |
| Refusal | `refusal_log` + audit_activity where refusal=true | First-class audit |
| Briefing | Aggregation of above | Dashboard projection |
| Chronicle | `GET /v1/activity` | Historical timeline |

**Domain-Pack Rendering Contract:**

The UI receives the Domain Pack descriptor and renders:
- Debt item names from `pack.debt_vocabulary` (not hardcoded)
- Claim types from `pack.claim_types`
- Falsifier types from `pack.falsifiers`
- Human-gated transitions from `pack.human_gated_transitions`

The UI kernel contains NO physics terms. Pack content is rendered dynamically.

**Empty/Unknown/Error States:**

| State | Rendering |
|-------|-----------|
| Unknown belief | "The Oracle has no record of this claim" |
| Empty debt | "No outstanding obligations" |
| No evidence | "No evidence has been presented" |
| No challenges | "No challenges have been raised" |
| API error | "The Oracle cannot answer this question right now" |
| Pack not loaded | "Domain Pack descriptor not available" |

**Data/API impact:** None on Solvent/Conductor. New Trust UI reads existing APIs.

**Tests:**
- IA document is internally consistent.
- All UI concepts map to Solvent API responses.
- No physics vocabulary in IA documents outside pack references.
- Empty states defined for every surface.

**Acceptance criteria:**
1. IA document complete with all 8 surfaces.
2. State projection table covers all Solvent data the UI needs.
3. Domain-pack rendering contract specifies dynamic rendering.
4. Empty/unknown/error states defined.
5. Read/write boundaries explicit (UI reads Solvent, writes only through Coordinator).
6. No Conductor modification required.

**Rollback/escape:** Delete `trust-ui/docs/` directory. No external impact.

**Dependencies:** Phase 2 (Domain Pack spec).

---

## Phase 9 — Trust UI Implementation

**Objective:** Build the six primary surfaces plus Briefing/Chronicle views. Go HTTP server with server-rendered HTML templates.

**Preconditions:** Phases 6, 7, 8 complete (Coordinator APIs exist, IA defined).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `conductor/internal/web/server.go` | Go HTTP server pattern |
| `conductor/internal/web/templates/layout.html` | Template structure |
| `solvent-main/demo/cloud/web/main.go` | Solvent web server pattern |
| `solvent-main/demo/cloud/web/templates/*.html` | Template patterns |

**Repository structure:**

```
trust-ui/
├── main.go                    # Entry point
├── go.mod                     # Go module
├── server/
│   ├── server.go              # HTTP server
│   ├── routes.go              # Route registration
│   └── middleware.go          # Auth, logging
├── handler/
│   ├── briefing.go            # Briefing Scroll handler
│   ├── oracle.go              # Oracle handler
│   ├── sphinx.go              # Sphinx/Riddle handler
│   ├── passage.go             # Passage Ledger handler
│   ├── challenge.go           # Challenge Room handler
│   ├── decision.go            # Decision Room handler
│   ├── refusal.go             # Refusal Archive handler
│   └── chronicle.go           # Chronicle handler
├── client/
│   ├── solvent.go             # Solvent API client
│   ├── coordinator.go         # Coordinator API client
│   └── conductor.go           # Conductor API client (read-only)
├── render/
│   ├── pack.go                # Domain Pack renderer
│   └── debt.go                # Debt visualization
├── templates/
│   ├── layout.html            # Base layout
│   ├── briefing.html          # Briefing Scroll
│   ├── oracle.html            # Oracle surface
│   ├── oracle_detail.html     # Oracle belief detail
│   ├── sphinx.html            # Sphinx/Riddle surface
│   ├── sphinx_detail.html     # Riddle detail
│   ├── passage.html           # Passage Ledger
│   ├── passage_detail.html    # Belief passage detail
│   ├── challenge.html         # Challenge Room
│   ├── challenge_detail.html  # Challenge detail
│   ├── decision.html          # Decision Room
│   ├── decision_detail.html   # Decision detail
│   ├── refusal.html           # Refusal Archive
│   ├── refusal_detail.html    # Refusal detail
│   └── chronicle.html         # Chronicle
├── docs/                      # IA docs from Phase 8
└── testdata/                  # Test fixtures
```

**Implementation order within Phase 9:**

1. Server skeleton + routing (server.go, routes.go)
2. Solvent client (client/solvent.go) — wraps existing Solvent REST API
3. Coordinator client (client/coordinator.go) — wraps Coordinator API
4. Briefing Scroll (handler/briefing.go + templates/briefing.html)
5. Oracle (handler/oracle.go + templates/oracle.html, oracle_detail.html)
6. Passage Ledger (handler/passage.go + templates/passage.html, passage_detail.html)
7. Sphinx/Riddle (handler/sphinx.go + templates/sphinx.html, sphinx_detail.html)
8. Decision Room (handler/decision.go + templates/decision.html, decision_detail.html)
9. Refusal Archive (handler/refusal.go + templates/refusal.html, refusal_detail.html)
10. Challenge Room (handler/challenge.go + templates/challenge.html, challenge_detail.html)
11. Chronicle (handler/chronicle.go + templates/chronicle.html)
12. Domain Pack renderer (render/pack.go) — dynamic debt/claim rendering

**Critical constraints:**
- UI writes ONLY through Coordinator/Solvent APIs (never direct DB).
- No Conductor modifications.
- No physics/domain logic in handler code.
- Domain Pack content rendered via pack descriptor, not hardcoded.

**Data/API impact:** New Trust UI service. Reads Solvent + Conductor + Coordinator APIs.

**Tests:**
- Each handler renders correctly with mock Solvent responses.
- Empty state rendering for each surface.
- Error state rendering when Solvent is unavailable.
- Domain Pack rendering works with pack descriptor.
- No direct database access from handler code (compile-time check).
- Route coverage: all surfaces accessible.

**Acceptance criteria:**
1. All 8 surfaces (6 primary + Briefing + Chronicle) render.
2. Each surface consumes Solvent API correctly.
3. Briefing Scroll shows current beliefs, promoted vs entered, open debts, active challenges, pending decisions, active intents, consequences.
4. Oracle shows claim, status, evidence, provenance, debt, challenges, relations, history, unknown.
5. Sphinx shows what/why/evidence/authority/result (PASS/REFUSE/HUMAN_REVIEW).
6. Passage Ledger shows debt items with ✓/○ status.
7. Challenge Room shows attacks and outcomes.
8. Decision Room shows consequential transitions with full context.
9. Refusal Archive shows refusal evidence as first-class audit.
10. Chronicle shows correlated historical activity.
11. Domain Pack terms rendered dynamically, not hardcoded.
12. No Conductor code modified.

**Rollback/escape:** Delete `trust-ui/` directory. No impact on Solvent/Conductor.

**Dependencies:** Phases 6, 7, 8.

---

## Phase 10 — Work + Adversarial Run

**Objective:** Execute the first complete BM-IST research cycle with one LLM-backed WORK agent and one deterministic ADVERSARIAL stub.

**Preconditions:** Phases 6, 7, 9 complete (Coordinator, human adjudication, Trust UI).

**Repository areas to inspect:**

| File | Relevance |
|------|-----------|
| `oracle/reference-loop/agent/agent.go` | Existing agent pattern |
| `oracle/reference-loop/main.go` | Reference loop harness |
| `oracle/domain-pack/bmist/v1/pack.json` | BM-IST pack |
| `oracle/packet/v1/examples/valid_packet.json` | Valid packet example |
| `oracle/verifier/physics/v1/verifier.go` | Physics verifier |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/run/work_agent.go` | LLM-backed WORK agent |
| `oracle/run/work_agent_test.go` | Work agent tests |
| `oracle/run/adversarial_stub.go` | Deterministic adversarial stub |
| `oracle/run/adversarial_stub_test.go` | Adversarial stub tests |
| `oracle/run/harness.go` | Run harness (orchestration) |
| `oracle/run/harness_test.go` | End-to-end run test |
| `oracle/run/acceptance.go` | Pre-registered acceptance criteria |
| `oracle/run/acceptance_test.go` | Acceptance verification |

**WORK Agent constraints:**
- Packet-only interface (produces Research Packet, no direct Solvent/Conductor calls).
- LLM-backed for research generation.
- Bounded scope: produces claim + evidence + debt for ONE belief.
- Cannot mutate protected state (no direct Solvent writes).
- Output is a validated Research Packet.

**ADVERSARIAL Stub constraints:**
- Deterministic, reproducible.
- Produces attack patterns against a given claim.
- Packet-only interface.
- Attack types: counterexample, contradiction, missing_evidence, alternative_explanation.
- Output is a validated Research Packet with `edges[].kind = "contradicts"`.

**Run orchestration:**
1. WORK agent produces packet for BM-IST Fisher-rigidity claim.
2. Coordinator compiles packet → Solvent belief + evidence.
3. Physics Verifier produces deterministic artifact.
4. Human adjudicator reviews belief + evidence + verifier artifact.
5. Human retires remaining debt items.
6. Human promotes belief.
7. ADVERSARIAL stub produces attack packet.
8. Coordinator compiles attack → Solvent contradiction edge.
9. Human reviews contradiction.
10. Human retracts belief (cascade cancels dependent intents).

**Data/API impact:** Uses existing Solvent + Conductor + Coordinator APIs. No new endpoints.

**Tests:**
- WORK agent produces valid packet.
- ADVERSARIAL stub produces valid attack packet.
- Coordinator compiles both without error.
- Physics verifier produces deterministic artifact.
- Full cycle completes: claim → evidence → debt → promotion → attack → retraction.
- No agent-caused protected-state mutations.

**Acceptance criteria (pre-registered):**
1. Malformed packet rejected by Coordinator.
2. Incomplete adversarial coverage detected (stub produces all 4 attack types).
3. Zero agent-caused protected-state mutations (agents only produce packets).
4. Human can reconstruct belief state from Trust UI projections.
5. Promotion gate works (debt must be retired before promotion).
6. Action intent requires promoted belief.
7. Retraction invalidates dependent authority/intents.

**Rollback/escape:** Delete `oracle/run/` directory. No impact on Solvent/Conductor.

**Dependencies:** Phases 6, 7, 9.

---

## Phase 11 — Consequential State Demonstration

**Objective:** Demonstrate every consequential epistemic transition visible in the Trust UI.

**Preconditions:** Phases 9, 10 complete (Trust UI + first run complete).

**Transitions to demonstrate:**

**Sequence A — Promotion path:**
```
claim → evidence → debt → verification → human decision → promotion → authorized intent
```

| Step | Solvent State | Trust UI Surface |
|------|--------------|------------------|
| Claim entered | belief.status = entered, debt = [6 items] | Oracle: shows entered belief |
| Evidence added | evidence rows exist | Oracle: shows evidence |
| Debt retired | debt items removed one by one | Passage Ledger: ✓ marks |
| Verified | verifier artifact produced | Oracle: shows evidence with artifact |
| Human decision | human retires remaining debt | Decision Room: decision recorded |
| Promoted | belief.status = promoted, debt = [] | Oracle: status = promoted |
| Intent authorized | action_intent.state = live | Sphinx: PASS |

**Sequence B — Retraction path:**
```
promoted belief → contradiction → retraction → dependent invalidation → intent invalidation
```

| Step | Solvent State | Trust UI Surface |
|------|--------------|------------------|
| Promoted | belief.status = promoted | Oracle: promoted |
| Contradiction | belief_edge.kind = contradicts | Challenge Room: attack found |
| Human retracts | belief.status = retracted | Oracle: retracted |
| Intent cancelled | action_intent.state = cancelled | Sphinx: REFUSE |
| Dependent retracted | descendant beliefs retracted | Oracle: cascade visible |

**Every transition must be visible in the Trust UI at the time it occurs.**

**Data/API impact:** No new endpoints. Uses existing Solvent + Trust UI.

**Tests:**
- Step-by-step: each transition produces correct Solvent state.
- Trust UI shows correct state at each step.
- No silent state mutations.
- Full audit trail in Chronicle.

**Acceptance criteria:**
1. Sequence A completes: claim → promotion → intent, all visible in UI.
2. Sequence B completes: promoted → retracted → intent cancelled, all visible in UI.
3. Every transition has a Trust UI surface that reflects it.
4. Audit trail in Chronicle matches Solvent audit_activity.

**Rollback/escape:** N/A — demonstration, no code changes.

**Dependencies:** Phases 9, 10.

---

## Phase 12 — Domain-Neutrality Proof

**Objective:** Prove BM-IST is data/configuration, not substrate logic. Test that replacing BM-IST pack data does not require Solvent/Conductor/UI kernel changes.

**Preconditions:** Phase 11 complete (full cycle demonstrated).

**Repository areas to inspect:**

| Area | Check |
|------|-------|
| Solvent kernel | Zero physics vocabulary |
| Conductor core | Zero physics vocabulary |
| Trust UI handlers | Zero physics vocabulary outside pack rendering |
| Coordinator | Zero physics vocabulary outside pack compilation |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/neutrality/scan.go` | Automated vocabulary scan |
| `oracle/neutrality/scan_test.go` | Scan tests |
| `oracle/neutrality/conformance.go` | Pack conformance test |
| `oracle/neutrality/conformance_test.go` | Conformance tests |

**Neutrality artifact structure:**

```go
type NeutralityReport struct {
    SolventScan   ScanResult `json:"solvent_scan"`
    ConductorScan ScanResult `json:"conductor_scan"`
    UIScan        ScanResult `json:"ui_scan"`
    CoordinatorScan ScanResult `json:"coordinator_scan"`
    PackConformance bool     `json:"pack_conformance"`
}

type ScanResult struct {
    RepositoryPath string   `json:"repository_path"`
    DomainTermsFound []string `json:"domain_terms_found"`
    Violations      []string `json:"violations"`
    Pass            bool     `json:"pass"`
}
```

**Scan terms (BM-IST specific, should NOT appear in substrate):**
- "Fisher", "rigidity", "Hamiltonian", "Schrödinger", "quantum", "density", "phase", "cocycle", "projective"
- "needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"

**Conformance test:**
- Load a DIFFERENT hypothetical pack (e.g., "legal-contract-v1") with different debt vocabulary.
- Verify the same Coordinator, Solvent, Conductor, and Trust UI code works without modification.
- This is a compile-time + runtime check, not a real second Domain Pack.

**Data/API impact:** None. Static analysis + conformance test.

**Tests:**
- Scan finds zero domain terms in Solvent kernel.
- Scan finds zero domain terms in Conductor core.
- Scan finds zero domain terms in Trust UI handlers.
- Scan finds zero domain terms in Coordinator (outside pack references).
- Conformance test passes with alternative pack descriptor.
- No hardcoded debt vocabulary in UI rendering code.

**Acceptance criteria:**
1. Solvent: zero violations.
2. Conductor: zero violations.
3. Trust UI: zero violations (pack rendering is dynamic).
4. Coordinator: zero violations (pack compilation is generic).
5. Conformance test passes with alternative pack.
6. Static scan report produced as artifact.

**Rollback/escape:** N/A — analysis only.

**Dependencies:** Phase 11 (full cycle demonstrated).

---

## Phase 13 — Adversarial System Review

**Objective:** Attempt to falsify the claim "BM-IST is a Domain Pack, not a kernel customization."

**Preconditions:** Phase 12 complete (neutrality proof exists).

**Attack surfaces to probe:**

| Attack | Method | Expected |
|--------|--------|----------|
| Physics leakage in Solvent | Grep for domain terms | None found |
| UI semantic leakage | Grep UI handlers for domain terms | None found |
| Hidden domain state | Schema inspection — no physics columns | None found |
| Agent authority | Agents produce packets, not state mutations | Confirmed |
| Duplicated state | Solvent is sole epistemic truth source | Confirmed |
| Direct UI database writes | UI calls only HTTP APIs | Confirmed |
| Conductor modifications | Conductor code unchanged from baseline | Confirmed |
| Coordinator reasoning | Coordinator has no inference logic | Confirmed |
| Incorrect authority projections | Authority always flows through Solvent | Confirmed |

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/adversarial/system_review.go` | Automated adversarial checks |
| `oracle/adversarial/system_review_test.go` | Review tests |
| `oracle/adversarial/REVIEW_REPORT.md` | Human-readable report |

**Data/API impact:** None. Analysis only.

**Tests:**
- Each attack surface produces a clear PASS/FAIL.
- No ambiguous results.
- Report is human-readable.

**Acceptance criteria:**
1. All 9 attack surfaces produce PASS.
2. No false negatives (every real violation detected).
3. Report documents methodology and findings.
4. Any NOT_YET_PROVEN items are explicitly flagged.

**Rollback/escape:** N/A — analysis only.

**Dependencies:** Phase 12.

---

## Phase 14 — POC Evidence Package

**Objective:** Produce the complete evidence package documenting the POC.

**Preconditions:** Phases 10-13 complete.

**Files/components to create:**

| Path | Content |
|------|---------|
| `oracle/evidence/manifest.md` | Evidence package index |
| `oracle/evidence/corpus_manifest.json` | Corpus manifest from Phase 5 |
| `oracle/evidence/packet_examples/` | Valid + malformed packet examples |
| `oracle/evidence/verifier_artifacts/` | Physics verifier output artifacts |
| `oracle/evidence/solvent_traces/` | Solvent state traces at each step |
| `oracle/evidence/conductor_projections/` | Conductor task projections |
| `oracle/evidence/ui_screenshots/` | Trust UI screenshots/flows |
| `oracle/evidence/adversarial_results/` | Adversarial review results |
| `oracle/evidence/human_decisions/` | Human decision records |
| `oracle/evidence/refusal_evidence/` | Refusal archive records |
| `oracle/evidence/consequential_traces/` | State transition traces |
| `oracle/evidence/failure_cases/` | Documented failure cases |
| `oracle/evidence/architecture_report.md` | Architecture report |

**Data/API impact:** None. Documentation only.

**Acceptance criteria:**
1. Every deliverable from Phases 0-13 is represented.
2. Evidence is verifiable (hashes match, traces are consistent).
3. Architecture report cites specific file paths and line numbers.
4. Failure cases are documented honestly.

**Rollback/escape:** N/A — documentation only.

**Dependencies:** Phases 10-13.

---

## Phase 15 — Final Acceptance

**Objective:** Decide PASS, PARTIAL PASS, or FAIL based on explicit acceptance criteria.

**Preconditions:** Phase 14 complete.

**Acceptance matrix:**

| Criterion | PASS condition | Evidence |
|-----------|---------------|----------|
| Domain Pack spec | BM-IST pack.json exists, validates | Phase 2 artifact |
| Packet contract | Valid + malformed packets tested | Phase 3 tests |
| Physics Verifier | Deterministic, no LLM, produces artifact | Phase 4 tests |
| Corpus manifest | All files hash-pinned | Phase 5 tests |
| Coordinator | Compiles packet → Solvent + Conductor | Phase 6 tests |
| Mandatory debt | Negative regression test passes | Phase 6 test |
| Human adjudication | All 6 decision types work | Phase 7 tests |
| Trust UI | All 8 surfaces render | Phase 9 tests |
| End-to-end run | Complete cycle with LLM agent + stub | Phase 10 tests |
| Consequential demo | Both sequences visible in UI | Phase 11 tests |
| Domain neutrality | Zero violations in substrate | Phase 12 scan |
| Adversarial review | 9/9 attack surfaces PASS | Phase 13 report |
| Evidence package | Complete and verifiable | Phase 14 artifacts |

**Verdict definitions:**

| Verdict | Meaning |
|---------|---------|
| PASS | All criteria met, honest claim warranted |
| PARTIAL PASS | Some criteria met, honest partial claim warranted |
| FAIL | Critical criteria not met, claim not warranted |

**Decision process:**
1. Review each criterion's evidence.
2. Check for any NOT_YET_PROVEN items.
3. Apply adversarial review findings.
4. Make honest assessment.
5. Record decision with rationale.

---

## Trust UI Information Architecture Summary

### Mythology → Technical Mapping

| Mythology Term | Technical Component | Primary Question |
|---------------|-------------------|-----------------|
| ORACLE | Belief query surface | What does the system know? |
| SPHINX | Authorization gate | May this request pass? |
| FERRYMAN | Execution boundary | What actually happened? |
| CHRONICLER | Activity timeline | What happened along the way? |
| CONDUCTOR | Workflow coordination | Who is doing the work? |

### UI Design Principles (Locked)

1. Oracle never guesses.
2. Sphinx never silently passes.
3. Refusal is first-class evidence.
4. Unknown is a valid state.
5. Every consequential state change is visible.
6. Human authority is explicit.
7. The UI presents canonical state; it does not invent state.
8. Domain Pack semantics are rendered dynamically.
9. Conductor remains unchanged.
10. The mythology is a semantic aid, not a replacement for technical labels.

---

## Explicit Non-Goals (Locked)

These are NOT implemented in this POC:

- Modifications to Conductor core or UI
- Cryptographic DSSE/SLSA attestation
- W3C PROV export
- Second Domain Pack (only conformance test with mock pack)
- Probabilistic inference
- Source reputation engine
- Ontology/knowledge graph engine
- Workflow primitives in Solvent
- Epistemic primitives in Conductor
- Embedded LLM reasoning inside Solvent
- Embedded LLM reasoning inside Coordinator
- Autonomous promotion
- Autonomous adjudication
- Pack-interface freeze

---

## Risks and Failure Modes

| Risk | Impact | Mitigation |
|------|--------|------------|
| Go CAS library insufficient for Fisher-rigidity | Phase 4 blocked | Use simplified symbolic checks, document limitation |
| Solvent API changes break Coordinator | Phase 6 blocked | Pin to frozen Solvent API contract |
| Trust UI can't render all Solvent state | Phase 9 incomplete | Document gaps, implement core surfaces first |
| LLM agent produces invalid packets | Phase 10 degraded | Coordinator rejects, agent retries |
| Adversarial stub too simplistic | Phase 13 weak | Document as POC limitation |
| Domain terms leak into substrate | Phase 12 fail | Grep-based CI gate |
| Human adjudication too slow for demo | Phase 11 delayed | Pre-script decisions for demo |

---

## Parallel vs Sequential Work

**Can run in parallel:**
- Phase 2 (Domain Pack) + Phase 5 (Corpus Manifest)
- Phase 4 (Verifier) partially overlaps with Phase 5
- Phase 8 (UI IA) starts as soon as Phase 2 lands
- Phase 9 (UI impl) begins partial work during Phase 6

**Must be sequential:**
- Phase 0 → Phase 1 → Phase 2 → Phase 3 → Phase 6 (chain)
- Phase 6 → Phase 7 → Phase 10 (chain)
- Phase 8 → Phase 9 (IA before implementation)
- Phase 10 → Phase 11 → Phase 12 → Phase 13 → Phase 14 → Phase 15 (chain)

**Critical path:** 0 → 1 → 2 → 3 → 6 → 7 → 10 → 11 → 12 → 13 → 14 → 15

---

## Definition of Done

The POC is DONE when:

1. BM-IST runs as a Domain Pack on the generic trust substrate.
2. Solvent contains zero physics vocabulary.
3. Conductor is unmodified.
4. Trust UI is a separate service consuming canonical state.
5. Coordinator is deterministic and non-reasoning.
6. WORK agent produces candidate research material, not authority.
7. ADVERSARIAL stub produces deterministic attacks.
8. Human controls all consequential epistemic transitions.
9. Every transition is visible in the Trust UI.
10. Domain-neutrality is proven by static scan + conformance test.
11. Adversarial review finds no violations.
12. Evidence package is complete and verifiable.
13. Final acceptance verdict is recorded with rationale.
