# ARGUS POC Execution Plan v1.1

**Version:** 1.1 — Final Pre-Implementation Revision
**Date:** 2026-09-15
**Status:** PLAN — IMPLEMENTATION READY
**Source:** v1.0 execution plan + four blocking fixes + consistency pass

---

## Final Implementation Lock

| Decision | Choice |
|----------|--------|
| Go modules | Two: `oracle/go.mod` + `trust-ui/go.mod`. `reference-loop/` untouched. |
| Coordinator deployment | Library in `oracle/coordinator/`. Thin HTTP wrapper in `oracle/coordinator/http/`. |
| Trust UI boundary | Calls Coordinator HTTP API, never imports Coordinator Go package. |
| Verifier scope | Hybrid: symbolic steps 1–2, hardcoded steps 3–4. |
| Artifact registry | In-memory, owned by Coordinator/harness. Trusted-verifier path only. |
| Decision persistence | Solvent state + audit_activity = authoritative. No second DB. |
| Sphinx projection | Read-only assembly of existing Solvent authority state. No policy logic. |
| Idempotency lifetime | In-memory for process lifetime. Restart clears. Race-safe. |
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
├── verifier/                    ← Phase 4
│   ├── registry.go              ← Phase 4 (FIX 1: artifact registry)
│   ├── registry_test.go
│   └── physics/v1/
│       ├── verifier.go
│       ├── verifier_test.go
│       ├── artifact.go
│       ├── artifact_test.go
│       ├── proof.go
│       └── proof_test.go
├── corpus/                      ← Phase 5
├── coordinator/                 ← Phase 6 (library)
│   ├── coordinator.go
│   ├── compiler.go
│   ├── compiler_test.go
│   ├── validate.go
│   ├── validate_test.go
│   ├── idempotency.go          ← Phase 6 (FIX 4: in-memory lifetime)
│   ├── idempotency_test.go
│   ├── projection.go
│   ├── projection_test.go
│   ├── human.go                 ← Phase 6 (FIX 2: no persistence)
│   ├── human_test.go
│   ├── sphinx.go                ← Phase 6 (FIX 3: authorization projection)
│   ├── sphinx_test.go
│   ├── client.go
│   └── http/
│       ├── server.go
│       ├── handler.go
│       └── types.go
├── run/                         ← Phase 10
├── neutrality/                  ← Phase 12
├── adversarial/                 ← Phase 13
├── evidence/                    ← Phase 14
├── plan/                        ← Existing + Phase 1
├── docs/                        ← Existing
├── reference-loop/              ← UNCHANGED (separate go.mod)
└── testdata/

trust-ui/                        ← go.mod (module github.com/PithomLabs/trust-ui)
├── main.go
├── server/
├── handler/
├── client/
├── render/
├── templates/
├── docs/                        ← Phase 8
└── testdata/
```

---

## Phase 0 — Repository Reconnaissance

**Objective:** Verify current state. No code changes.

**Deliverable:** Confirmation that all repos are at expected state.

**Verification:**
```bash
cd /home/chaschel/Documents/go/solvent-main && git rev-parse HEAD && git status --short && go test ./...
cd /home/chaschel/Documents/go/conductor && git rev-parse HEAD && git status --short && go test ./...
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

**Validator rules:**
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

**Tests:**
- Valid pack loads and validates
- Missing `pack_id` → error
- Empty `claim_types` → error
- `initial_debt` item not in `debt_vocabulary` → error
- `retirement_rules` references undeclared evidence class → error
- `human_gated_transitions` missing required gate → error
- Duplicate debt vocabulary item → error

**Acceptance:** All tests pass. Pack matches BM-IST POC seed vocabulary exactly.

**Dependencies:** Phase 1.

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

**Key validation rules:**
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

**Dependencies:** Phase 2.

---

## Phase 4 — Physics Verifier Runner

**Objective:** Deterministic Go proof-obligation verifier for 4 bounded Fisher-rigidity identities, plus artifact registry.

### FIX 1 — Verifier Artifact Registry

**Files to create:**

| File | Purpose |
|------|---------|
| `oracle/verifier/registry.go` | In-memory artifact registry |
| `oracle/verifier/registry_test.go` | Registry tests |
| `oracle/verifier/physics/v1/verifier.go` | Verifier orchestrator |
| `oracle/verifier/physics/v1/verifier_test.go` | Deterministic output tests |
| `oracle/verifier/physics/v1/artifact.go` | Artifact types |
| `oracle/verifier/physics/v1/artifact_test.go` | Hash determinism tests |
| `oracle/verifier/physics/v1/proof.go` | 4 proof-obligation checks |
| `oracle/verifier/physics/v1/proof_test.go` | Proof-obligation tests |

**Artifact Registry design:**

```go
package verifier

type ArtifactRegistry struct {
    mu       sync.RWMutex
    artifacts map[string]VerificationArtifact // key = artifact.EvidenceRef
}

func NewArtifactRegistry() *ArtifactRegistry

// Register adds a trusted verifier artifact to the registry.
// Only artifacts produced by the trusted Physics Verifier runner may be registered.
func (r *ArtifactRegistry) Register(ctx context.Context, artifact VerificationArtifact) error

// Resolve retrieves an artifact by its evidence reference.
// Returns error if not found.
func (r *ArtifactRegistry) Resolve(ctx context.Context, ref string) (VerificationArtifact, error)

// VerifyHash checks that the artifact's hash matches the expected value.
// Returns error on mismatch or missing artifact.
func (r *ArtifactRegistry) VerifyHash(ref string, expectedSHA256 string) error
```

**Invariants:**
1. A packet cannot admit `reproducible_artifact` evidence merely because the packet contains a forged JSON artifact.
2. The Coordinator resolves artifact references through the registry, not from packet contents.
3. The resolved artifact hash must equal the packet's declared `content_sha256`.
4. The artifact must have been produced by the trusted verifier runner during the current POC process.
5. Unknown artifact reference → reject with typed error.
6. Hash mismatch → reject with typed error.
7. Registry contents are not persisted. In-memory only.
8. Restart persistence is OUT OF SCOPE.
9. No cryptographic signing, DSSE, SLSA, or database.

**Proof obligations:**

| Step | Type | Description |
|------|------|-------------|
| 1 | `machine_verified` | Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²) |
| 2 | `machine_verified` | Coefficient matching: a(ρ) = κ²/(8mρ) |
| 3 | `human_attested` | |∇ρ|² coefficient consistency |
| 4 | `human_attested` | b′(ρ) = 0 |

**Steps 1–2 (symbolic):**
- Lightweight symbolic expression tree: `Var`, `Const`, `Add`, `Mul`, `Pow`, `Div`, `Sqrt`, `Grad`, `Laplacian`
- Simplification: constant folding, power rules, basic algebraic identities
- No full CAS needed

**Steps 3–4 (hardcoded):**
- Expected values as constants
- Deterministic assertion: input matches expected
- Recorded as `human_attested` in artifact

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
- **Registry: forged artifact rejected** (packet contains JSON artifact but registry is empty → reject)
- **Registry: unknown artifact_ref rejected** (ref not in registry → reject)
- **Registry: hash mismatch rejected** (artifact in registry but packet declares different hash → reject)
- **Registry: valid trusted artifact accepted** (artifact registered by verifier, hash matches → admit)

**Acceptance:** Verifier runs without LLM. Artifact deterministic. Registry enforces trusted-artifact-only admission.

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

**Corpus rules:**
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

**Objective:** Packet compiler library + thin HTTP wrapper. No reasoning. No persistence.

### FIX 2 — Decision Record Persistence

**Locked design:**
- Solvent state + Solvent audit_activity = authoritative persisted decision history
- `DecisionRecord` is a request/result envelope for the UI/API, NOT a second system of record
- Every consequential decision results in the existing Solvent state transition
- The consequential action is represented in Solvent's existing audit/activity trail
- REFUSE uses the existing Solvent `refusal_log` mechanism (migration 003) where applicable
- The Coordinator may return a DecisionRecord to the UI, but does not persist a second copy
- Trust UI decision history is reconstructable from canonical Solvent state + audit/refusal records
- Do NOT create: coordinator.sqlite, decision database, decision table, UI decision store

**Solvent audit trail mapping:**

| Decision | Solvent mutation | Audit entry type | Refusal log |
|----------|-----------------|------------------|-------------|
| PROMOTE | `POST /v1/beliefs/{id}/promote` | `belief_promoted` | `promote` (on SQLSTATE 23514) |
| RETIRE_DEBT | `POST /v1/beliefs/{id}/debt/retire` | `debt_retired` | — |
| RETRACT | `POST /v1/beliefs/{id}/retract` | `belief_retracted` | — |
| REOPEN | `POST /v1/beliefs` + edge | `belief_entered` | — |
| AUTHORIZE | `POST /v1/targets/{id}/approve` | `authorization_granted` | — |
| REFUSE | No mutation | Audit entry with `refusal: true` | `retract_unsafe` or custom |

### FIX 3 — Sphinx Authorization Context

**New Coordinator method + HTTP endpoint:**

```go
// SphinxProjection assembles existing Solvent authority state into a
// UI-readable authorization/riddle context. It does NOT decide new policy.
// It is a read-only VIEW of canonical state.
func (c *Coordinator) AuthorizationContext(ctx context.Context, targetID string) (*SphinxProjection, error)

type SphinxProjection struct {
    Target            TargetInfo            `json:"target"`
    SupportingBelief  *BeliefInfo           `json:"supporting_belief,omitempty"`
    BeliefStatus      string                `json:"belief_status"`
    EvidenceSummary   []EvidenceInfo        `json:"evidence_summary"`
    OpenDebt          []string              `json:"open_debt"`
    AuthorityState    string                `json:"authority_state"`    // approved | pending | denied | revoked
    JustificationState string               `json:"justification_state"`
    RequestState      string                `json:"request_state"`
    ApprovalState     string                `json:"approval_state"`
    IntentState       string                `json:"intent_state"`      // live | executing | executed | cancelled | none
    CurrentResult     string                `json:"current_result"`    // PASS | REFUSE | HUMAN_REVIEW
    RefusalReason     string                `json:"refusal_reason,omitempty"`
}
```

**Result determination (read-only, no new policy):**
- If target is revoked → `REFUSE`
- If belief is retracted → `REFUSE`
- If belief is not promoted → `REFUSE`
- If no approval exists → `HUMAN_REVIEW`
- If approval exists + belief promoted + no revocation → `PASS`
- All logic reads existing Solvent state; no inference, no scoring

### Files to create:

| File | Purpose |
|------|---------|
| `oracle/coordinator/coordinator.go` | Core compilation logic |
| `oracle/coordinator/compiler.go` | Packet → Solvent + Conductor compilation |
| `oracle/coordinator/compiler_test.go` | Compilation tests |
| `oracle/coordinator/validate.go` | Packet + debt validation |
| `oracle/coordinator/validate_test.go` | Validation tests |
| `oracle/coordinator/idempotency.go` | In-memory canonical-content idempotency |
| `oracle/coordinator/idempotency_test.go` | Idempotency tests |
| `oracle/coordinator/projection.go` | In-memory retry queue |
| `oracle/coordinator/projection_test.go` | Projection repair tests |
| `oracle/coordinator/human.go` | Human decision logic (no persistence) |
| `oracle/coordinator/human_test.go` | Decision tests |
| `oracle/coordinator/sphinx.go` | Authorization context projection (FIX 3) |
| `oracle/coordinator/sphinx_test.go` | Sphinx projection tests |
| `oracle/coordinator/client.go` | HTTP clients for Solvent + Conductor |
| `oracle/coordinator/http/server.go` | HTTP wrapper |
| `oracle/coordinator/http/handler.go` | HTTP handlers |
| `oracle/coordinator/http/types.go` | HTTP request/response types |

### FIX 4 — Idempotency Lifetime

**Locked design:**
- In-memory for lifetime of Coordinator process
- Canonical content hash excludes `packet_id`, runtime metadata, timestamps
- Identical canonical content under different `packet_id` → duplicate, no-op
- Duplicate submission during same Coordinator process → no-op, return existing result
- Coordinator restart clears the mapping
- Cross-restart duplicate detection is OUT OF SCOPE
- Do NOT add: persistent idempotency DB, Redis, SQLite, distributed locking
- Concurrent submissions with same canonical hash must be race-safe

**Implementation:**
```go
type IdempotencyCache struct {
    mu      sync.RWMutex
    entries map[string]*CompilationResult // key = canonical content hash
}

func (c *IdempotencyCache) Check(canonicalHash string) (*CompilationResult, bool)
func (c *IdempotencyCache) Store(canonicalHash string, result *CompilationResult)
```

**Coordinator library API:**

```go
type Coordinator struct {
    solventClient    *SolventClient
    conductorClient  *ConductorClient
    packRegistry     *PackRegistry
    artifactRegistry *verifier.ArtifactRegistry
    projectionQueue  *ProjectionQueue
    idempotencyCache *IdempotencyCache
    operatorID       string // from ARGUS_OPERATOR_PRINCIPAL_ID
}

func New(cfg Config) *Coordinator

// Core compilation
func (c *Coordinator) CompilePacket(ctx context.Context, pkt *packet.ResolvedPacket) (*CompilationResult, error)

// Human decisions (FIX 2: no persistence, writes to Solvent only)
func (c *Coordinator) SubmitDecision(ctx context.Context, req DecisionRequest) (*DecisionRecord, error)

// Query
func (c *Coordinator) PacketStatus(ctx context.Context, packetID string) (*PacketStatus, error)
func (c *Coordinator) DecisionContext(ctx context.Context, beliefID string) (*DecisionContext, error)

// Sphinx projection (FIX 3: read-only assembly of Solvent authority state)
func (c *Coordinator) AuthorizationContext(ctx context.Context, targetID string) (*SphinxProjection, error)

// Projection retry
func (c *Coordinator) RetryPending(ctx context.Context) error
```

**Compilation steps:**
1. Validate packet against schema
2. Resolve pack_ref → load Domain Pack
3. Check idempotency cache (canonical content hash) → return cached result if hit
4. For each belief in packet:
   a. Compute debt = pack.initial_debt ∪ packet.belief.debt
   b. `POST /v1/beliefs` with claim + compiled debt
   c. For each evidence referencing this belief:
      - If `provenance_class == "reproducible_artifact"`: resolve `artifact_ref` through ArtifactRegistry, verify hash matches `content_sha256`
      - `POST /v1/evidence`
5. For each edge: record locally (Conductor doesn't have edges)
6. For each task: `POST /v1/projects/{id}/tasks` with `status: "proposed"`, `governance_ref: belief_id`
7. If Conductor write fails: queue in ProjectionQueue for retry
8. Store result in idempotency cache
9. Return CompilationResult

**Human decision logic (FIX 2: Solvent-only persistence):**
```go
type DecisionRequest struct {
    Type       string `json:"type"`       // PROMOTE, RETIRE_DEBT, RETRACT, REOPEN, AUTHORIZE, REFUSE
    BeliefID   string `json:"belief_id"`
    ScenarioID string `json:"scenario_id"`
    DebtItem   string `json:"debt_item,omitempty"` // for RETIRE_DEBT
    Reason     string `json:"reason,omitempty"`
}

type DecisionRecord struct {
    ID            string    `json:"id"`
    Type          string    `json:"type"`
    BeliefID      string    `json:"belief_id"`
    ScenarioID    string    `json:"scenario_id"`
    ActorID       string    `json:"actor_id"`       // from ARGUS_OPERATOR_PRINCIPAL_ID
    Evidence      string    `json:"evidence"`       // snapshot of current Solvent state
    Result        string    `json:"result"`         // executed | refused
    RefusalReason string    `json:"refusal_reason,omitempty"`
    SolventAuditID string   `json:"solvent_audit_id,omitempty"` // audit_activity entry ID
    CreatedAt     time.Time `json:"created_at"`
}
```

**Decision execution (FIX 2):**
- PROMOTE: `POST /v1/beliefs/{id}/promote?scenario_id=X` → Solvent writes `belief_promoted` audit entry. On refusal (SQLSTATE 23514), Solvent writes to `refusal_log` with statement `promote`.
- RETIRE_DEBT: `POST /v1/beliefs/{id}/debt/retire?scenario_id=X` → Solvent writes `debt_retired` audit entry.
  - Before invoking: mechanically verify pack retirement rule (debt item + evidence class/type match).
- RETRACT: `POST /v1/beliefs/{id}/retract?scenario_id=X` → Solvent writes `belief_retracted` audit entry.
- REOPEN: `POST /v1/beliefs` with same claim + `derives` edge from retracted origin → Solvent writes `belief_entered` audit entry.
- AUTHORIZE: `POST /v1/targets/{id}/approve` → Solvent writes `authorization_granted` audit entry.
- REFUSE: No Solvent mutation. Coordinator returns refusal DecisionRecord with `result: "refused"`.

**Every decision exposes:** current state, requested transition, evidence, open obligations, downstream consequences, acting principal, resulting state.

**DecisionRecord is NOT persisted by the Coordinator.** It is a transient response envelope. The authoritative record is:
- Solvent `audit_activity` table (for executed decisions)
- Solvent `refusal_log` table (for refused decisions, where applicable)
- Trust UI reconstructs decision history from these Solvent sources

**HTTP wrapper endpoints:**

| Method | Path | Handler | Notes |
|--------|------|---------|-------|
| POST | `/decisions` | Submit human decision | Returns DecisionRecord |
| GET | `/packets/:id/status` | Packet compilation status | |
| GET | `/beliefs/:id/decision-context` | Decision context for belief | |
| GET | `/authorization-context/:target_id` | Sphinx projection (FIX 3) | Read-only, no policy |

**Tests (expanded for all four fixes):**

*FIX 1 — Artifact registry:*
- Forged verifier artifact → rejected
- Unknown artifact_ref → rejected
- Correct artifact + correct hash → admitted
- Hash mismatch → rejected

*FIX 2 — Decision persistence:*
- PROMOTE creates Solvent audit_activity entry with type `belief_promoted`
- RETIRE_DEBT creates Solvent audit_activity entry with type `debt_retired`
- RETRACT creates Solvent audit_activity entry with type `belief_retracted`
- REFUSE creates Solvent audit_activity entry with `refusal: true`
- Coordinator restart → decision history reconstructable from Solvent audit_activity
- No coordinator.sqlite or decision table exists (verified by test)

*FIX 3 — Sphinx projection:*
- Approved authority state → PASS projection
- Missing required state → HUMAN_REVIEW or REFUSE per Solvent state
- Retracted supporting belief → REFUSE
- No new authority decision logic in Coordinator (verified by code inspection test)

*FIX 4 — Idempotency:*
- Same packet twice → one compilation
- Same content + different packet_id → one compilation
- Concurrent duplicate submissions → one compilation (race test)
- Coordinator restart → idempotency cache empty

*Existing:*
- Valid packet compiles correctly
- Malformed packet rejected
- Free-text edge target rejected
- Empty-debt packet still receives initial_debt
- Retirement request fails when evidence class doesn't match pack rule
- Conductor failure after Solvent success → queued for retry
- Retry succeeds on second attempt
- Missing `ARGUS_OPERATOR_PRINCIPAL_ID` → startup fails

**Acceptance:** All tests pass. Coordinator has zero reasoning logic. No decision persistence. Sphinx is read-only projection. Idempotency is process-local.

**Dependencies:** Phases 2, 3, 4.

---

## Phase 7 — Human Adjudication

**Objective:** Minimum consequential decision path. Built into Phase 6 Coordinator.

**Note:** Phase 7 is implemented as part of Phase 6 (`coordinator/human.go`). This phase's deliverable is the test suite and verification that all 6 decision types work through Solvent-only persistence.

**Additional tests beyond Phase 6:**
- PROMOTE on belief with debt → blocked (Solvent schema gate returns Verdict, refusal logged to `refusal_log`)
- PROMOTE on belief with no debt → success, audit entry present
- RETRACT on promoted belief with live intent → cascade cancels intent first
- REFUSE → audit entry with `refusal: true`, no Solvent state change
- Decision without configured principal → rejected
- Decision deduplication works (same decision submitted twice → deduplicated)
- **Decision history reconstructable from Solvent audit_activity after Coordinator restart**

**Acceptance:** All 6 decision types verified through Solvent-only persistence. Decision history is canonical Solvent state.

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

**Navigation structure:**
```
/ (BRIEFING SCROLL)
/oracle (ORACLE)
/oracle/:belief_id
/sphinx (SPHINX)
/sphinx/:target_id
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
- Coordinator HTTP API: `POST /decisions`, `GET /packets/:id/status`, `GET /beliefs/:id/decision-context`, `GET /authorization-context/:target_id`
- Conductor API: `/v1/tasks/:id`, `/v1/tasks/:id/activity` (read-only, non-authoritative)

**State projections (FIX 2 — decision history from Solvent):**

| UI Concept | Solvent Source | Mapping |
|------------|---------------|---------|
| Decision history | `audit_activity` (type in `belief_promoted`, `debt_retired`, `belief_retracted`, `belief_entered`) | Reconstructable from Solvent |
| Refusal history | `refusal_log` table + `audit_activity` where `refusal: true` | First-class Solvent data |
| Sphinx/Authorization | Coordinator `GET /authorization-context/:target_id` → reads Solvent target/belief/intent state | Read-only projection |

**Acceptance:** IA complete with all 8 surfaces. State projections cover all Solvent data. Decision history reconstructable from Solvent.

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

**Implementation order:**

**9.1 Core control surface (before first run):**
1. `go.mod` + `main.go`
2. `server/server.go` + `server/routes.go` + `server/middleware.go`
3. `client/solvent.go` (Solvent REST client)
4. `client/coordinator.go` (Coordinator HTTP client, includes `AuthorizationContext` method for FIX 3)
5. `render/pack.go` + `render/debt.go`
6. `templates/layout.html`
7. `handler/briefing.go` + `templates/briefing.html`
8. `handler/oracle.go` + `templates/oracle.html` + `templates/oracle_detail.html`
9. `handler/decision.go` + `templates/decision.html` + `templates/decision_detail.html`
10. `handler/refusal.go` + `templates/refusal.html` + `templates/refusal_detail.html`

**9.2 Expansion (after first run):**
11. `handler/sphinx.go` + `templates/sphinx.html` + `templates/sphinx_detail.html`
12. `handler/passage.go` + `templates/passage.html` + `templates/passage_detail.html`
13. `handler/challenge.go` + `templates/challenge.html` + `templates/challenge_detail.html`
14. `handler/chronicle.go` + `templates/chronicle.html`
15. `client/conductor.go`

**Solvent client API:**
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

**Coordinator HTTP client (FIX 3 added):**
```go
type CoordinatorClient struct {
    baseURL    string
    httpClient *http.Client
}

func (c *CoordinatorClient) SubmitDecision(ctx, req DecisionRequest) (*DecisionRecord, error)
func (c *CoordinatorClient) PacketStatus(ctx, packetID string) (*PacketStatus, error)
func (c *CoordinatorClient) DecisionContext(ctx, beliefID string) (*DecisionContext, error)
func (c *CoordinatorClient) AuthorizationContext(ctx, targetID string) (*SphinxProjection, error) // FIX 3
```

**Decision history reconstruction (FIX 2):**
- Trust UI fetches `GET /v1/activity?type=belief_promoted&type=debt_retired&type=belief_retracted&type=belief_entered` from Solvent
- Trust UI fetches `GET /v1/activity?refusal=true` for refusal display
- Trust UI does NOT call a separate decision-history endpoint
- Decision records shown in Decision Room are reconstructed from these Solvent audit entries

**Critical constraints:**
- UI writes ONLY through Coordinator HTTP API (never direct DB, never direct Solvent writes)
- No Conductor modifications
- No physics/domain logic in handler code
- Pack content rendered dynamically
- Decision history comes from Solvent audit_activity, not a separate store

**Tests:**
- Each handler renders with mock Solvent responses
- Empty state rendering for each surface
- Error state when Solvent unavailable
- Domain Pack rendering with pack descriptor
- No direct database access (compile-time check via interface)
- All routes accessible
- Sphinx handler calls Coordinator `AuthorizationContext` and renders projection

**Acceptance:** All 8 surfaces render. Domain Pack terms dynamic. No Conductor modifications. Decision history from Solvent.

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
3. Physics Verifier produces artifact → registers in ArtifactRegistry
4. Coordinator admits verifier evidence (artifact_ref resolved through registry, hash validated)
5. Human reviews, retires debt, promotes
6. ADVERSARIAL stub produces attack packets
7. Coordinator validates coverage, compiles contradiction edges
8. Human reviews contradiction, retracts

**Pre-registered acceptance criteria:**
1. Malformed packet rejected
2. Incomplete adversarial coverage detected (2/4 attacks flagged)
3. Complete adversarial coverage passes (4/4 attacks)
4. Zero agent-caused protected-state mutations
5. Human reconstructs belief state from Trust UI (decision history from Solvent audit_activity)
6. Promotion gate works
7. Action intent requires promoted belief
8. Retraction invalidates dependent authority/intents
9. **Forged verifier artifact rejected by registry**
10. **Valid trusted artifact admitted by registry**

**Tests:**
- WORK agent produces valid packet
- ADVERSARIAL stub produces valid attack packets
- Incomplete coverage variant triggers detection
- Forged verifier artifact rejected
- Full cycle completes
- No agent state mutations
- **Decision history reconstructable from Solvent after full cycle**

**Acceptance:** All pre-registered criteria pass.

**Dependencies:** Phases 6, 7, 9.

---

## Phase 11 — Consequential State Demonstration

**Objective:** Demonstrate every consequential transition visible in Trust UI.

**No new code.** Uses Phases 9+10 deliverables.

**Sequence A (promotion path):**
1. Claim entered → Oracle shows entered
2. Evidence added → Oracle shows evidence
3. Debt retired one by one → Passage Ledger shows ✓
4. Verifier artifact → Oracle shows artifact (registry admitted)
5. Human retires remaining debt → Decision Room records (from Solvent audit)
6. Promoted → Oracle shows promoted
7. Intent authorized → Sphinx shows PASS (from authorization-context projection)

**Sequence B (retraction path):**
1. Promoted → Oracle shows promoted
2. Contradiction → Challenge Room shows attack
3. Human retracts → Oracle shows retracted
4. Intent cancelled → Sphinx shows REFUSE
5. Dependent retracted → Oracle cascade visible

**Verification:**
- Screenshots/recordings of each step in Trust UI
- Decision history in Decision Room matches Solvent audit_activity entries
- Refusal evidence in Refusal Archive matches Solvent refusal_log
- Sphinx projection matches Solvent authority state
- Chronicle matches Solvent activity timeline

**Acceptance:** Both sequences complete and visible in UI. Decision history canonical from Solvent.

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

**9 attack surfaces:**
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
| `oracle/evidence/artifact_registry/` | **FIX 1:** Registry admission evidence |
| `oracle/evidence/solvent_traces/` | State traces |
| `oracle/evidence/conductor_projections/` | Task projections |
| `oracle/evidence/ui_screenshots/` | UI flows |
| `oracle/evidence/adversarial_results/` | Review results |
| `oracle/evidence/human_decisions/` | Decision records (from Solvent audit) |
| `oracle/evidence/refusal_evidence/` | Refusal records (from Solvent refusal_log) |
| `oracle/evidence/consequential_traces/` | Transition traces |
| `oracle/evidence/failure_cases/` | Failures |
| `oracle/evidence/architecture_report.md` | Architecture report |
| `oracle/evidence/forged_artifact_rejection/` | **FIX 1:** Forged artifact test evidence |
| `oracle/evidence/decision_reconstruction/` | **FIX 2:** Decision history from Solvent audit |
| `oracle/evidence/sphinx_projection/` | **FIX 3:** Authorization-context response |
| `oracle/evidence/idempotency_results/` | **FIX 4:** Duplicate + concurrent test results |
| `oracle/evidence/restart_limitation/` | **FIX 4:** Restart clears idempotency cache |

**Acceptance:** Every phase represented. Evidence verifiable. Report cites file paths. All four fixes documented.

**Dependencies:** Phases 10–13.

---

## Phase 15 — Final Acceptance

**Objective:** PASS / PARTIAL PASS / FAIL verdict.

**Acceptance matrix:**

| Criterion | PASS condition |
|-----------|---------------|
| Domain Pack | pack.json exists, validates |
| Packet contract | Valid + malformed tested |
| Verifier | Deterministic, no LLM, artifact produced |
| Artifact registry | Forged rejected, unknown rejected, hash mismatch rejected, valid accepted |
| Corpus | All files hash-pinned |
| Coordinator | Compiles packet → Solvent + Conductor |
| Mandatory debt | Negative regression passes |
| Human adjudication | 6 decision types work |
| Decision persistence | Recoverable from Solvent audit_activity, no second DB |
| Trust UI | 8 surfaces render |
| Sphinx projection | Authorization-context works, no new policy logic |
| Idempotency | Canonical duplicate suppressed, concurrent race-safe, restart clears |
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
| 4 | 7 | Symbolic algebra + hardcoded checks + registry |
| 5 | 2 | Manifest + hash verification |
| 6 | 9 | Library + HTTP + idempotency + projection + human + sphinx |
| 7 | 2 | Tests (logic in Phase 6) |
| 8 | 2 | Documentation |
| 9 | 9 | 8 surfaces + clients + renderer + decision reconstruction |
| 10 | 5 | Agents + harness + acceptance |
| 11 | 2 | Demo + screenshots |
| 12 | 3 | Scanner + conformance |
| 13 | 2 | Automated review |
| 14 | 3 | Evidence collection (expanded for 4 fixes) |
| 15 | 1 | Verdict |

---

## Key API Contracts

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
GET  /v1/targets/{id}               → reads authority target (FIX 3)
```

### Conductor endpoints the Coordinator calls:

```
POST /v1/projects/{id}/tasks        → creates task (status: proposed)
GET  /v1/tasks/{id}                 → reads task
GET  /v1/tasks/{id}/activity        → reads task activity
```

### Coordinator HTTP endpoints (for Trust UI):

```
POST /decisions                          → submit human decision
GET  /packets/:id/status                 → packet compilation status
GET  /beliefs/:id/decision-context       → decision context for belief
GET  /authorization-context/:target_id   → Sphinx projection (FIX 3)
```

---

## Risks

| Risk | Mitigation |
|------|------------|
| Bounded verifier insufficient | Explicit machine_verified vs human_attested |
| Solvent API breakage | Pin to frozen API contract |
| Conductor projection failure | In-memory retry queue |
| LLM agent invalid packets | Coordinator rejects |
| Forged verifier artifact | ArtifactRegistry enforces trusted path only |
| Concurrent compilation | Per-operation locking + race tests |
| Domain terms leak | Context-aware scan + allowlist |
| **Artifact registry lost on restart** | **Accepted POC limitation. Re-register artifacts on harness restart.** |
| **Coordinator restart loses idempotency cache** | **Accepted POC limitation. Documented. Cross-restart duplicates not detected.** |
| **Solvent audit projection insufficient for a UI field** | **Document exact missing canonical source; do not invent persistence** |
| **Authorization projection accidentally gains policy logic** | **Keep Sphinx projection read-only; test state-to-view mapping** |

---

## POC Limitations (explicit)

1. Artifact registry is in-memory. Restart requires re-registration of verifier artifacts.
2. Idempotency cache is process-local. Restart clears. Cross-restart duplicate detection is out of scope.
3. Projection journal is in-memory. Crash during Conductor projection requires manual retry.
4. DecisionRecord is a transient envelope, not persisted. History reconstructed from Solvent.
5. No cryptographic attestation (DSSE/SLSA). Trusted path relies on process-level trust.
6. Second Domain Pack not implemented. Only conformance test with mock pack.
7. No persistent infrastructure introduced (no new databases, no Redis, no SQLite).

---

## Definition of Done

The POC is DONE when:

1. BM-IST runs as a Domain Pack on the generic trust substrate.
2. Solvent contains zero physics vocabulary.
3. Conductor is unmodified.
4. Trust UI is a separate service consuming canonical state.
5. Coordinator is deterministic and performs only declared structural/artifact checks.
6. Artifact registry enforces trusted-verifier-only admission.
7. Decision history is reconstructable from Solvent audit_activity (no second DB).
8. Sphinx projection is read-only assembly of existing Solvent authority state.
9. Idempotency is process-local and race-safe.
10. WORK agent produces candidate research material, not authority.
11. ADVERSARIAL stub produces deterministic attacks with tested incomplete-coverage path.
12. Human controls all consequential epistemic transitions.
13. Every consequential transition is visible in the Trust UI.
14. Domain-neutrality is proven by context-aware static scan + conformance test.
15. Adversarial review passes all defined/probed surfaces honestly.
16. Evidence package is complete and verifiable.
17. Final acceptance verdict is recorded with rationale.
