# ARGUS Execution Plan — Implementation-Specific

**Version:** 1.0
**Date:** 2026-09-15
**Status:** PLAN — NOT IMPLEMENTED
**Source:** v1.1 review-consolidated plan + locked decisions

---

## Locked Implementation Decisions

| Decision | Choice |
|----------|--------|
| Go modules | Two: `oracle/go.mod` + `trust-ui/go.mod`. `reference-loop/` untouched. |
| Coordinator deployment | Library in `oracle/coordinator/`. Thin HTTP wrapper in `oracle/coordinator/http/`. |
| Trust UI boundary | Calls Coordinator HTTP API, never the Go package directly. |
| Verifier scope | Hybrid: symbolic steps 1–2, hardcoded steps 3–4. |
| Projection journal | In-memory retry queue. No crash recovery. |
| Scan allowlist | `domain-pack/**`, `testdata/**`, `packet/v1/examples/**` |
| Human principal | `ARGUS_OPERATOR_PRINCIPAL_ID` env var. Fail closed if missing. |
| reference-loop | Parallel, untouched. Known-good historical reference. |

---

## Module Layout

```
oracle/                          ← go.mod (module github.com/PithomLabs/oracle)
├── domain-pack/bmist/v1/        ← Phase 2
├── packet/v1/                   ← Phase 3
├── verifier/physics/v1/         ← Phase 4
├── corpus/                      ← Phase 5
├── coordinator/                 ← Phase 6 (library)
│   └── http/                    ← Phase 6 (HTTP wrapper)
├── run/                         ← Phase 10
├── neutrality/                  ← Phase 12
├── adversarial/                 ← Phase 13
├── evidence/                    ← Phase 14
├── plan/                        ← Existing plans + Phase 1 reconciliation
├── docs/                        ← Existing docs
├── reference-loop/              ← UNCHANGED (separate go.mod)
└── testdata/                    ← Shared test fixtures

trust-ui/                        ← go.mod (module github.com/PithomLabs/trust-ui)
├── main.go
├── server/
├── handler/
├── client/
├── render/
├── templates/
├── docs/                        ← Phase 8 IA docs
└── testdata/
```

---

## Phase 0 — Repository Reconnaissance

**Objective:** Verify current state. No code changes.

**Deliverable:** Confirmation that all repos are at expected state.

**Verification:**
```bash
# Solvent
cd /home/chaschel/Documents/go/solvent-main && git rev-parse HEAD && git status --short && go test ./...

# Conductor
cd /home/chaschel/Documents/go/conductor && git rev-parse HEAD && git status --short && go test ./...

# Oracle
cd /home/chaschel/Documents/go/oracle && git rev-parse HEAD && git status --short
```

**Acceptance:** All tests pass, working trees clean.

---

## Phase 1 — Freeze Reconciliation

**Objective:** Record the immutable baseline. Create reconciliation record.

**Deliverables:**

| File | Content |
|------|---------|
| `oracle/plan/freeze_reconciliation.md` | Strategy-level record of frozen artifacts, supersessions, amendments, deferrals |

**Freeze reconciliation record must document:**
1. Prior freeze hash: `7602699` (superseded by Plan 11.1)
2. Current freeze: Plan 11.1 applied, new hash recorded here
3. Which designs are superseded (old `FullDebt` constant, old wizard vocabulary coupling)
4. Which designs are amended (debt now opaque, `wizardDebt` in `internal/belief/debt.go`)
5. Which capabilities remain deferred (DSSE/SLSA, second Domain Pack, etc.)
6. Schema state: migrations 001–010, `debt SET DEFAULT ARRAY[]::TEXT[]`

**Verification:**
```bash
cd /home/chaschel/Documents/go/solvent-main
git rev-parse HEAD              # Record as NEW_FREEZE_HASH
git status --short              # Must be empty
go test ./...                   # Full suite
grep -r "FullDebt" kernel/      # Must return empty
```

**Acceptance:** Freeze hash recorded, tree clean, tests pass, `FullDebt` absent from kernel.

**Dependencies:** Plan 11.1 complete (confirmed).

---

## Phase 2 — Domain Pack Specification

**Objective:** Define the BM-IST Domain Pack as a JSON descriptor with Go types and validator.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/go.mod` | Go module definition |
| `oracle/domain-pack/bmist/v1/pack.json` | Domain Pack descriptor |
| `oracle/domain-pack/bmist/v1/types.go` | Go types mirroring JSON schema |
| `oracle/domain-pack/bmist/v1/validate.go` | Structural validator |
| `oracle/domain-pack/bmist/v1/validate_test.go` | Positive + negative tests |
| `oracle/domain-pack/README.md` | Pack contract explanation |
| `oracle/testdata/packs/valid/` | Valid test pack fixtures |
| `oracle/testdata/packs/invalid/` | Invalid test pack fixtures |

**`go.mod` initial content:**
```
module github.com/PithomLabs/oracle

go 1.25

require (
    github.com/google/uuid v1.6.0
)
```

**Validator rules (Phase 2 specific):**
1. `pack_id` required, non-empty
2. `version` required, semver format
3. `claim_types` required, non-empty array, must contain at least one of `derived`, `accommodated`, `postulated`
4. `debt_vocabulary` required, non-empty array
5. `initial_debt` required, subset of `debt_vocabulary`
6. `retirement_rules` keys must be subset of `debt_vocabulary`
7. `retirement_rules` values must reference declared `evidence_classes`
8. `human_gated_transitions` must contain the structural human-gate floor (faithfulness_review, scope_clarification, obstruction_assessment)
9. `falsifiers` required, non-empty array
10. No duplicate items in any array

**`pack.json` — exact content from v1.1 plan §Phase 2**

**Tests:**
- Valid pack loads and validates
- Missing `pack_id` → error
- Empty `claim_types` → error
- `initial_debt` item not in `debt_vocabulary` → error
- `retirement_rules` references undeclared evidence class → error
- `human_gated_transitions` missing required gate → error
- Duplicate debt vocabulary item → error

**Acceptance:** All tests pass. Pack matches BM-IST POC seed vocabulary exactly.

**Dependencies:** Phase 1 (freeze established).

---

## Phase 3 — EBP Research Packet Contract

**Objective:** Define the Research Packet v1 as JSON schema with Go validation.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/packet/v1/schema.json` | JSON Schema for interchange |
| `oracle/packet/v1/types.go` | Go types |
| `oracle/packet/v1/resolve.go` | Local/canonical reference resolver |
| `oracle/packet/v1/validate.go` | Validation logic |
| `oracle/packet/v1/validate_test.go` | Positive + negative tests |
| `oracle/packet/v1/examples/valid_packet.json` | Valid example |
| `oracle/packet/v1/examples/malformed_*.json` | Negative examples |

**Key validation rules (from v1.1):**
1. `packet_id` required, UUID
2. `pack_ref` required, must resolve to registered pack
3. Every object has `local_id`
4. `beliefs[].claim` required, non-empty
5. `beliefs[].claim_type` in pack's `claim_types`
6. `evidence[].content_sha256` required, well-formed hex
7. `evidence[].belief_ref` references a declared `local_id` or canonical `belief:<uuid>`
8. `edges[].from_ref`/`to_ref` use local IDs or canonical `belief:<uuid>` — **no free-text claim matching**
9. `edges[].kind` in `{derives, contradicts}`
10. Debt items in pack vocabulary
11. Tasks compiled only as `proposed`

**`resolve.go` — reference resolution:**
- `local:b1` → resolves to packet object with `local_id: "b1"`
- `canonical:belief:<uuid>` → resolves to Solvent belief by UUID
- Resolution is deterministic, no semantic matching
- Resolution failures return typed errors

**Tests:**
- Valid packet passes
- Missing `packet_id` → rejected
- Empty claim → rejected
- Unknown claim_type → rejected
- Free-text edge target → rejected
- `belief_ref` to nonexistent local_id → rejected
- Debt item not in pack vocabulary → rejected
- Canonical ref with invalid UUID → rejected

**Acceptance:** Schema valid, Go types compile, all tests pass.

**Dependencies:** Phase 2 (pack defines vocabulary).

---

## Phase 4 — Physics Verifier Runner

**Objective:** Deterministic Go proof-obligation verifier for 4 bounded Fisher-rigidity identities.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/verifier/physics/v1/verifier.go` | Verifier orchestrator |
| `oracle/verifier/physics/v1/verifier_test.go` | Deterministic output tests |
| `oracle/verifier/physics/v1/artifact.go` | Artifact types |
| `oracle/verifier/physics/v1/artifact_test.go` | Hash determinism tests |
| `oracle/verifier/physics/v1/proof.go` | 4 proof-obligation checks |
| `oracle/verifier/physics/v1/proof_test.go` | Proof-obligation tests |

**Proof obligations (from v1.1):**

| Step | Type | Description |
|------|------|-------------|
| 1 | `machine_verified` | Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²) |
| 2 | `machine_verified` | Coefficient matching: a(ρ) = κ²/(8mρ) |
| 3 | `human_attested` | |∇ρ|² coefficient consistency |
| 4 | `human_attested` | b′(ρ) = 0 |

**Implementation approach for step 1 (symbolic):**
- Represent ρ as a symbolic variable
- Compute √ρ, then Δ(√ρ)/√ρ using chain rule
- Simplify and compare to ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²)
- Use a lightweight symbolic expression tree (no full CAS needed)
- Expression types: `Var`, `Const`, `Add`, `Mul`, `Pow`, `Div`, `Sqrt`, `Grad`, `Laplacian`
- Simplification: constant folding, power rules, basic algebraic identities

**Implementation approach for step 2 (symbolic):**
- Input: Lagrangian L with general E_q[ρ] = a(ρ)|∇ρ|²
- Euler-Lagrange matching forces a(ρ) = κ²/(8mρ)
- Verify the coefficient identity holds symbolically

**Steps 3–4 (hardcoded):**
- Encode expected values as constants
- Deterministic assertion: input matches expected
- Record as `human_attested` in artifact

**Artifact structure:**
```go
type VerificationArtifact struct {
    RunID           string    `json:"run_id"`
    VerifierVersion string    `json:"verifier_version"`
    VerifierHash    string    `json:"verifier_hash"`
    ClaimHash       string    `json:"claim_hash"`
    InputHash       string    `json:"input_hash"`
    ArtifactHash    string    `json:"artifact_hash"`
    Timestamp       time.Time `json:"timestamp,omitempty"`
    Steps           []Step    `json:"steps"`
    Result          string    `json:"result"`
    EvidenceRef     string    `json:"evidence_ref"`
}

type Step struct {
    Index       int    `json:"index"`
    Description string `json:"description"`
    Input       string `json:"input"`
    Rule        string `json:"rule"`
    Output      string `json:"output"`
    Verified    bool   `json:"verified"`
    CheckType   string `json:"check_type"` // "machine_verified" | "human_attested"
}
```

**Hash determinism:**
- `ArtifactHash` = SHA256 of canonical JSON (excluding `ArtifactHash` and `Timestamp` fields)
- Same inputs → same hash across 100 runs (test this)
- `Timestamp` excluded from hash computation

**Tests:**
- Step 1 symbolic check produces deterministic output
- Step 2 symbolic check produces deterministic output
- Steps 3–4 hardcoded checks pass
- Artifact hash is reproducible across 100 runs
- Missing input → error
- Malformed expression → error
- `Result` is one of `confirmed`, `refuted`, `inconclusive`
- Each step records `check_type` correctly

**Acceptance:** Verifier runs without LLM. Artifact deterministic. 4 steps produce correct check_type.

**Dependencies:** Phases 2+3.

---

## Phase 5 — Corpus Manifest

**Objective:** Hash-pin BM-IST research corpus.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/corpus/manifest.json` | Corpus manifest |
| `oracle/corpus/manifest.go` | Go loader + validator |
| `oracle/corpus/manifest_test.go` | Hash verification tests |

**Corpus rules (from v1.1):**
- Include BM-IST seed and Qwen-derived artifact
- `agent_derived` is corpus provenance only, does not confer authority
- Qwen artifact is context, not source truth, unless separately admitted
- Only verifier-produced artifacts enter Solvent as `reproducible_artifact`

**Tests:**
- Manifest loads and validates
- All referenced files exist and hashes match
- Missing file detected
- Hash mismatch detected

**Acceptance:** All corpus files hash-pinned. Manifest loads in Go. Hashes verified.

**Dependencies:** Phase 2.

---

## Phase 6 — Deterministic Coordinator

**Objective:** Packet compiler library + thin HTTP wrapper.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/coordinator/coordinator.go` | Core compilation logic |
| `oracle/coordinator/compiler.go` | Packet → Solvent + Conductor compilation |
| `oracle/coordinator/compiler_test.go` | Compilation tests |
| `oracle/coordinator/validate.go` | Packet + debt validation |
| `oracle/coordinator/validate_test.go` | Validation tests |
| `oracle/coordinator/idempotency.go` | Canonical-content idempotency |
| `oracle/coordinator/idempotency_test.go` | Idempotency tests |
| `oracle/coordinator/projection.go` | In-memory retry queue |
| `oracle/coordinator/projection_test.go` | Projection repair tests |
| `oracle/coordinator/human.go` | Human decision logic |
| `oracle/coordinator/human_test.go` | Decision tests |
| `oracle/coordinator/client.go` | HTTP clients for Solvent + Conductor |
| `oracle/coordinator/http/server.go` | HTTP wrapper (3 endpoints) |
| `oracle/coordinator/http/handler.go` | HTTP handlers |
| `oracle/coordinator/http/types.go` | HTTP request/response types |

**Coordinator library API:**

```go
type Coordinator struct {
    solventClient   *SolventClient
    conductorClient *ConductorClient
    packRegistry    *PackRegistry
    projectionQueue *ProjectionQueue
    operatorID      string // from ARGUS_OPERATOR_PRINCIPAL_ID
}

func New(solventURL, conductorURL, operatorID string) *Coordinator

// Core compilation
func (c *Coordinator) CompilePacket(ctx context.Context, pkt *packet.ResolvedPacket) (*CompilationResult, error)

// Human decisions
func (c *Coordinator) SubmitDecision(ctx context.Context, req DecisionRequest) (*DecisionRecord, error)

// Query
func (c *Coordinator) PacketStatus(ctx context.Context, packetID string) (*PacketStatus, error)
func (c *Coordinator) DecisionContext(ctx context.Context, beliefID string) (*DecisionContext, error)

// Projection retry
func (c *Coordinator) RetryPending(ctx context.Context) error
```

**Compilation steps (from v1.1):**
1. Validate packet against schema
2. Resolve pack_ref → load Domain Pack
3. For each belief in packet:
   a. Compute debt = pack.initial_debt ∪ packet.belief.debt
   b. `POST /v1/beliefs` with claim + compiled debt
   c. For each evidence referencing this belief:
      - If `provenance_class == "reproducible_artifact"`: verify `artifact_ref` resolves to registered verifier output AND hash matches
      - `POST /v1/evidence`
4. For each edge: record locally (Conductor doesn't have edges)
5. For each task: `POST /v1/projects/{id}/tasks` with `status: "proposed"`, `governance_ref: belief_id`
6. If Conductor write fails: queue in ProjectionQueue for retry
7. Record compilation result

**Idempotency (from v1.1):**
- Canonical content hash = SHA256(sorted canonical JSON of packet, excluding `packet_id` and runtime metadata)
- Same canonical content under different `packet_id` → no-op, return existing result
- Store mapping: `canonical_hash → compilation_result`

**Projection queue:**
- In-memory `[]FailedProjection`
- Each entry: `{beliefID, taskData, error, retries}`
- `RetryPending()` re-attempts failed Conductor projections
- No crash recovery (documented POC limitation)

**Human decision logic:**
```go
type DecisionRequest struct {
    Type      string `json:"type"`      // PROMOTE, RETIRE_DEBT, RETRACT, REOPEN, AUTHORIZE, REFUSE
    BeliefID  string `json:"belief_id"`
    ScenarioID string `json:"scenario_id"`
    DebtItem  string `json:"debt_item,omitempty"` // for RETIRE_DEBT
    Reason    string `json:"reason,omitempty"`
}

type DecisionRecord struct {
    ID          string    `json:"id"`
    Type        string    `json:"type"`
    BeliefID    string    `json:"belief_id"`
    ScenarioID  string    `json:"scenario_id"`
    ActorID     string    `json:"actor_id"`     // from ARGUS_OPERATOR_PRINCIPAL_ID
    Evidence    string    `json:"evidence"`     // snapshot of current state
    Result      string    `json:"result"`       // executed | refused
    RefusalReason string `json:"refusal_reason,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Decision execution:**
- PROMOTE: `POST /v1/beliefs/{id}/promote?scenario_id=X`
- RETIRE_DEBT: `POST /v1/beliefs/{id}/debt/retire?scenario_id=X` with `{"debt_item": "..."}`
  - Before invoking: mechanically verify pack retirement rule (debt item + evidence class/type match)
- RETRACT: `POST /v1/beliefs/{id}/retract?scenario_id=X`
- REOPEN: `POST /v1/beliefs` with same claim + `derives` edge from retracted origin
- AUTHORIZE: `POST /v1/targets/{id}/approve`
- REFUSE: Log only, no Solvent mutation

**Every decision exposes:** current state, requested transition, evidence, open obligations, downstream consequences, acting principal, resulting state.

**HTTP wrapper endpoints:**

| Method | Path | Handler |
|--------|------|---------|
| POST | `/decisions` | Submit human decision |
| GET | `/packets/:id/status` | Packet compilation status |
| GET | `/beliefs/:id/decision-context` | Decision context for belief |

**Tests:**
- Valid packet compiles correctly
- Malformed packet rejected
- Free-text edge target rejected
- Empty-debt packet still receives initial_debt
- Forged verifier artifact rejected
- Retirement request fails when evidence class doesn't match pack rule
- Same canonical content under different packet IDs is idempotent
- Conductor failure after Solvent success → queued for retry
- Retry succeeds on second attempt
- Human decision recorded with correct actor_id
- Duplicate decision submissions deduplicated
- Concurrent same-packet compilation is race-safe
- Missing `ARGUS_OPERATOR_PRINCIPAL_ID` → startup fails

**Acceptance:** All tests pass. Coordinator has zero reasoning logic (grep for confidence, score, infer, reason).

**Dependencies:** Phases 2, 3, 4.

---

## Phase 7 — Human Adjudication

**Objective:** Minimum consequential decision path. Built into Phase 6 Coordinator.

**Note:** Phase 7 is implemented as part of Phase 6 (`coordinator/human.go`). This phase's deliverable is the test suite and verification that all 6 decision types work.

**Additional tests beyond Phase 6:**
- PROMOTE on belief with debt → blocked (Solvent schema gate returns Verdict)
- PROMOTE on belief with no debt → success
- RETRACT on promoted belief with live intent → cascade cancels intent first
- REFUSE → logged, no state change
- Decision without configured principal → rejected
- Decision deduplication works

**Acceptance:** All 6 decision types verified. Every consequential transition logged.

**Dependencies:** Phase 6.

---

## Phase 8 — Trust UI Information Architecture

**Objective:** Define IA, state projections, navigation, domain-pack rendering, read/write boundaries.

**Files to create:**

| File | Purpose |
|------|---------|
| `trust-ui/docs/information-architecture.md` | IA specification |
| `trust-ui/docs/state-projections.md` | Solvent → UI mapping |
| `trust-ui/docs/navigation.md` | Route structure |
| `trust-ui/docs/domain-pack-contract.md` | Dynamic pack rendering |
| `trust-ui/docs/empty-states.md` | Empty/unknown/error semantics |
| `trust-ui/docs/read-write-boundaries.md` | Read/write rules |

**Navigation structure (from v1.1):**
```
/ (BRIEFING SCROLL)
/oracle (ORACLE)
/oracle/:belief_id
/sphinx (SPHINX)
/sphinx/:riddle_id
/passage (PASSAGE LEDGER)
/passage/:belief_id
/challenge (CHALLENGE ROOM)
/challenge/:belief_id
/decision (DECISION ROOM)
/decision/:decision_id
/refusal (REFUSAL ARCHIVE)
/refusal/:refusal_id
/chronicle (CHRONICLE)
```

**Trust UI API consumption:**
- Solvent API: `/v1/beliefs`, `/v1/beliefs/:id`, `/v1/beliefs/:id/evidence`, `/v1/beliefs/:id/explain`, `/v1/ledger`, `/v1/activity`, `/v1/evidence`
- Coordinator HTTP API: `POST /decisions`, `GET /packets/:id/status`, `GET /beliefs/:id/decision-context`
- Conductor API: `/v1/tasks/:id`, `/v1/tasks/:id/activity` (read-only, non-authoritative)

**Acceptance:** IA complete with all 8 surfaces. State projections cover all Solvent data.

**Dependencies:** Phases 2, 6, 7.

---

## Phase 9 — Trust UI Implementation

**Objective:** Build 8 surfaces in Go HTTP + server-rendered HTML.

**Files to create:**

```
trust-ui/
├── go.mod
├── main.go
├── server/
│   ├── server.go
│   ├── routes.go
│   └── middleware.go
├── handler/
│   ├── briefing.go
│   ├── oracle.go
│   ├── sphinx.go
│   ├── passage.go
│   ├── challenge.go
│   ├── decision.go
│   ├── refusal.go
│   └── chronicle.go
├── client/
│   ├── solvent.go
│   ├── coordinator.go
│   └── conductor.go
├── render/
│   ├── pack.go
│   └── debt.go
├── templates/
│   ├── layout.html
│   ├── briefing.html
│   ├── oracle.html
│   ├── oracle_detail.html
│   ├── sphinx.html
│   ├── sphinx_detail.html
│   ├── passage.html
│   ├── passage_detail.html
│   ├── challenge.html
│   ├── challenge_detail.html
│   ├── decision.html
│   ├── decision_detail.html
│   ├── refusal.html
│   ├── refusal_detail.html
│   └── chronicle.html
└── testdata/
```

**`trust-ui/go.mod`:**
```
module github.com/PithomLabs/trust-ui

go 1.25
```

**Implementation order (from v1.1 §9.1 + §9.2):**

**9.1 Core control surface (before first run):**
1. `go.mod` + `main.go` (server skeleton)
2. `server/server.go` + `server/routes.go` + `server/middleware.go`
3. `client/solvent.go` (Solvent REST client)
4. `client/coordinator.go` (Coordinator HTTP client)
5. `render/pack.go` + `render/debt.go` (Domain Pack renderer)
6. `templates/layout.html` (base layout with mythology navigation)
7. `handler/briefing.go` + `templates/briefing.html` (Briefing Scroll)
8. `handler/oracle.go` + `templates/oracle.html` + `templates/oracle_detail.html`
9. `handler/decision.go` + `templates/decision.html` + `templates/decision_detail.html`
10. `handler/refusal.go` + `templates/refusal.html` + `templates/refusal_detail.html`

**9.2 Expansion (after first run):**
11. `handler/sphinx.go` + `templates/sphinx.html` + `templates/sphinx_detail.html`
12. `handler/passage.go` + `templates/passage.html` + `templates/passage_detail.html`
13. `handler/challenge.go` + `templates/challenge.html` + `templates/challenge_detail.html`
14. `handler/chronicle.go` + `templates/chronicle.html`
15. `client/conductor.go` (read-only Conductor client)

**Solvent client API (matching Solvent's actual endpoints):**
```go
type SolventClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *SolventClient) GetBelief(ctx, scenarioID, beliefID string) (*BeliefResponse, error)
func (c *SolventClient) ListBeliefs(ctx, scenarioID, status, claimType string, limit, offset int) (*BeliefListResponse, error)
func (c *SolventClient) GetEvidence(ctx, scenarioID, beliefID string) ([]EvidenceResponse, error)
func (c *SolventClient) ExplainBelief(ctx, scenarioID, beliefID string) (*ExplainResponse, error)
func (c *SolventClient) GetLedger(ctx, scenarioID string) (*LedgerSummaryResponse, error)
func (c *SolventClient) GetActivity(ctx, scenarioID string, limit int) (*ActivityResponse, error)
```

**Coordinator HTTP client:**
```go
type CoordinatorClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *CoordinatorClient) SubmitDecision(ctx, req DecisionRequest) (*DecisionRecord, error)
func (c *CoordinatorClient) PacketStatus(ctx, packetID string) (*PacketStatus, error)
func (c *CoordinatorClient) DecisionContext(ctx, beliefID string) (*DecisionContext, error)
```

**Domain Pack renderer (`render/pack.go`):**
- Loads pack descriptor at startup
- Renders debt item names from `pack.debt_vocabulary`
- Renders claim types from `pack.claim_types`
- Renders falsifier types from `pack.falsifiers`
- No hardcoded domain terms in template logic

**Critical constraints:**
- UI writes ONLY through Coordinator HTTP API (never direct DB, never direct Solvent writes for consequential transitions)
- No Conductor modifications
- No physics/domain logic in handler code
- Pack content rendered dynamically

**Tests:**
- Each handler renders with mock Solvent responses
- Empty state rendering for each surface
- Error state when Solvent unavailable
- Domain Pack rendering with pack descriptor
- No direct database access (compile-time check via interface)
- All routes accessible

**Acceptance:** All 8 surfaces render. Domain Pack terms dynamic. No Conductor modifications.

**Dependencies:** Phases 6, 7, 8.

---

## Phase 10 — Work + Adversarial Run

**Objective:** First complete BM-IST research cycle.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/run/work_agent.go` | LLM-backed WORK agent |
| `oracle/run/work_agent_test.go` | Work agent tests |
| `oracle/run/adversarial_stub.go` | Deterministic adversarial stub |
| `oracle/run/adversarial_stub_test.go` | Adversarial stub tests |
| `oracle/run/harness.go` | Run orchestration |
| `oracle/run/harness_test.go` | End-to-end test |
| `oracle/run/acceptance.go` | Pre-registered acceptance criteria |
| `oracle/run/acceptance_test.go` | Acceptance verification |

**WORK agent:**
- LLM-backed (configurable model endpoint)
- Packet-only interface
- Produces: claim + evidence + debt for ONE belief
- Cannot mutate protected state
- Output: validated `packet.ResolvedPacket`

**ADVERSARIAL stub:**
- Deterministic, no LLM
- Input: claim text + pack falsifier list
- Produces 4 attack packets (one per falsifier type):
  1. counterexample
  2. contradiction
  3. missing_evidence
  4. alternative_explanation
- Each attack is a packet with `edges[].kind = "contradicts"`
- Also produces an incomplete-coverage variant (only 2 of 4 attacks) for testing

**Run orchestration:**
1. WORK agent produces packet
2. Coordinator compiles → Solvent belief + evidence
3. Physics Verifier produces artifact
4. Coordinator admits verifier evidence (artifact_ref + hash match)
5. Human reviews, retires debt, promotes
6. ADVERSARIAL stub produces attack packets
7. Coordinator validates coverage, compiles contradiction edges
8. Human reviews contradiction, retracts

**Pre-registered acceptance criteria:**
1. Malformed packet rejected
2. Incomplete adversarial coverage detected (2/4 attacks flagged)
3. Complete adversarial coverage passes (4/4 attacks)
4. Zero agent-caused protected-state mutations
5. Human reconstructs belief state from Trust UI
6. Promotion gate works
7. Action intent requires promoted belief
8. Retraction invalidates dependent authority/intents

**Tests:**
- WORK agent produces valid packet
- ADVERSARIAL stub produces valid attack packets
- Incomplete coverage variant triggers detection
- Forged verifier artifact rejected
- Full cycle completes
- No agent state mutations

**Acceptance:** All 8 pre-registered criteria pass.

**Dependencies:** Phases 6, 7, 9.

---

## Phase 11 — Consequential State Demonstration

**Objective:** Demonstrate every consequential transition visible in Trust UI.

**No new code.** Uses Phases 9+10 deliverables.

**Sequence A (promotion path):**
1. Claim entered → Oracle shows entered
2. Evidence added → Oracle shows evidence
3. Debt retired one by one → Passage Ledger shows ✓
4. Verifier artifact → Oracle shows artifact
5. Human retires remaining debt → Decision Room records
6. Promoted → Oracle shows promoted
7. Intent authorized → Sphinx shows PASS

**Sequence B (retraction path):**
1. Promoted → Oracle shows promoted
2. Contradiction → Challenge Room shows attack
3. Human retracts → Oracle shows retracted
4. Intent cancelled → Sphinx shows REFUSE
5. Dependent retracted → Oracle cascade visible

**Verification:** Screenshots/recordings of each step in Trust UI. Audit trail in Chronicle matches Solvent audit_activity.

**Acceptance:** Both sequences complete and visible in UI.

**Dependencies:** Phases 9, 10.

---

## Phase 12 — Domain-Neutrality Proof

**Objective:** Prove BM-IST is data, not substrate logic.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/neutrality/scan.go` | Context-aware vocabulary scanner |
| `oracle/neutrality/scan_test.go` | Scan tests |
| `oracle/neutrality/conformance.go` | Alternative pack conformance |
| `oracle/neutrality/conformance_test.go` | Conformance tests |

**Scanner design:**
- Path-aware: knows which directories are allowlisted
- Allowlist: `domain-pack/**`, `testdata/**`, `packet/v1/examples/**`
- Scan targets: all production Go code in oracle/, solvent-main/, conductor/, trust-ui/
- Match: exact identifiers and bounded phrases (not generic words)
- BM-IST scan terms: `needMap`, `needInvariant`, `needToyCheck`, `needNullModel`, `needObstruction`, `needFaithfulnessReview`, `Fisher-rigidity`, `FisherInformation`, `HamiltonianSchrödinger`
- Exclusions recorded in report

**Conformance test:**
- Create mock pack: `legal-contract-v1` with different debt vocabulary
- Load into same Coordinator, verify compilation works
- Verify Trust UI renders with alternative pack
- Verify no code changes needed

**Acceptance:** Zero violations in substrate. Conformance test passes with alternative pack. Report produced.

**Dependencies:** Phase 11.

---

## Phase 13 — Adversarial System Review

**Objective:** Attempt to falsify "BM-IST is a Domain Pack, not a kernel customization."

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/adversarial/system_review.go` | Automated adversarial checks |
| `oracle/adversarial/system_review_test.go` | Review tests |
| `oracle/adversarial/REVIEW_REPORT.md` | Human-readable report |

**9 attack surfaces (from v1.1):**
1. Physics leakage in Solvent → grep
2. UI semantic leakage → grep handlers
3. Hidden domain state → schema inspection
4. Agent authority → agents produce packets only
5. Duplicated state → Solvent is sole truth
6. Direct UI database writes → HTTP-only
7. Conductor modifications → diff from baseline
8. Coordinator reasoning → grep for inference logic
9. Incorrect authority projections → trace authority flow

**Acceptance:** All probed surfaces PASS. Report documents methodology and coverage.

**Dependencies:** Phase 12.

---

## Phase 14 — POC Evidence Package

**Objective:** Complete evidence documentation.

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/evidence/manifest.md` | Index |
| `oracle/evidence/corpus_manifest.json` | From Phase 5 |
| `oracle/evidence/packet_examples/` | Valid + malformed |
| `oracle/evidence/verifier_artifacts/` | Verifier output |
| `oracle/evidence/solvent_traces/` | State traces |
| `oracle/evidence/conductor_projections/` | Task projections |
| `oracle/evidence/ui_screenshots/` | UI flows |
| `oracle/evidence/adversarial_results/` | Review results |
| `oracle/evidence/human_decisions/` | Decision records |
| `oracle/evidence/refusal_evidence/` | Refusal records |
| `oracle/evidence/consequential_traces/` | Transition traces |
| `oracle/evidence/failure_cases/` | Failures |
| `oracle/evidence/architecture_report.md` | Architecture report |

**Acceptance:** Every phase represented. Evidence verifiable. Report cites file paths.

**Dependencies:** Phases 10–13.

---

## Phase 15 — Final Acceptance

**Objective:** PASS / PARTIAL PASS / FAIL verdict.

**Acceptance matrix (from v1.1):**

| Criterion | PASS condition |
|-----------|---------------|
| Domain Pack | pack.json exists, validates |
| Packet contract | Valid + malformed tested |
| Verifier | Deterministic, no LLM, artifact produced |
| Corpus | All files hash-pinned |
| Coordinator | Compiles packet → Solvent + Conductor |
| Mandatory debt | Negative regression passes |
| Human adjudication | 6 decision types work |
| Trust UI | 8 surfaces render |
| End-to-end | Complete cycle with agents |
| Consequential demo | Both sequences visible |
| Domain neutrality | Zero violations |
| Adversarial review | All surfaces PASS |
| Evidence package | Complete and verifiable |

**Dependencies:** Phase 14.

---

## Critical Path

```
0 → 1 → 2 → 3 → 4 → 6 → 7 → 8 → 9 → 10 → 11 → 12 → 13 → 14 → 15
```

**Parallel:**
- Phase 5 runs parallel with Phase 2–4

**Estimated phase effort (relative):**

| Phase | Effort | Notes |
|-------|--------|-------|
| 0 | 1 | Read-only |
| 1 | 1 | Verification + record |
| 2 | 3 | JSON + Go types + validator |
| 3 | 4 | JSON schema + resolver + validation |
| 4 | 6 | Symbolic algebra + hardcoded checks |
| 5 | 2 | Manifest + hash verification |
| 6 | 8 | Core library + HTTP wrapper + idempotency |
| 7 | 2 | Tests (logic in Phase 6) |
| 8 | 2 | Documentation |
| 9 | 8 | 8 surfaces + clients + renderer |
| 10 | 5 | Agents + harness + acceptance |
| 11 | 2 | Demo + screenshots |
| 12 | 3 | Scanner + conformance |
| 13 | 2 | Automated review |
| 14 | 2 | Evidence collection |
| 15 | 1 | Verdict |

---

## Key API Contracts (for implementer reference)

### Solvent endpoints the Coordinator calls:

```
POST /v1/beliefs                    → creates belief
POST /v1/evidence                   → adds evidence
POST /v1/beliefs/{id}/debt/retire   → retires one debt item
POST /v1/beliefs/{id}/promote       → promotes (may return Verdict refusal)
POST /v1/beliefs/{id}/retract       → retracts + cascade
POST /v1/targets                    → creates authority target
POST /v1/targets/{id}/justifications → attaches justification
POST /v1/targets/{id}/request       → requests authorization
POST /v1/targets/{id}/approve       → approves (activates)
POST /v1/authorizations/action      → creates live intent
POST /v1/authorizations/execute     → executes authorized action
GET  /v1/beliefs/{id}               → reads belief
GET  /v1/beliefs/{id}/explain       → reads belief + edges + intents
GET  /v1/beliefs/{id}/evidence      → reads evidence for belief
GET  /v1/ledger                     → reads ledger summary
GET  /v1/activity                   → reads audit activity
```

### Conductor endpoints the Coordinator calls:

```
POST /v1/projects/{id}/tasks        → creates task (status: proposed)
GET  /v1/tasks/{id}                 → reads task
GET  /v1/tasks/{id}/activity        → reads task activity
```

### Coordinator HTTP endpoints (for Trust UI):

```
POST /decisions                     → submit human decision
GET  /packets/:id/status            → packet compilation status
GET  /beliefs/:id/decision-context  → decision context for belief
```

---

## Risks (from v1.1, with mitigations)

| Risk | Mitigation |
|------|------------|
| Bounded verifier insufficient | Explicit machine_verified vs human_attested |
| Solvent API breakage | Pin to frozen API contract |
| Conductor projection failure | In-memory retry queue |
| LLM agent invalid packets | Coordinator rejects |
| Forged verifier artifact | Only trusted-verifier artifacts admitted |
| Concurrent compilation | Per-operation locking + race tests |
| Domain terms leak | Context-aware scan + allowlist |
