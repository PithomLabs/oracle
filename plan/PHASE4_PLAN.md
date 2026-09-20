# PHASE 4 PLAN — PHYSICS VERIFIER RUNNER

## Files to Create (8 total)

```
oracle/verifier/
├── registry.go
├── registry_test.go
└── physics/v1/
    ├── artifact.go
    ├── artifact_test.go
    ├── proof.go
    ├── proof_test.go
    ├── verifier.go
    └── verifier_test.go
```

## Step 1: `oracle/verifier/registry.go` — Artifact Registry

```go
package verifier

type ArtifactRegistry struct {
    mu        sync.RWMutex
    artifacts map[string]VerificationArtifact // key = EvidenceRef
}

func NewArtifactRegistry() *ArtifactRegistry

// registerTrusted — unexported. Only callable from within verifier/ package.
func (r *ArtifactRegistry) registerTrusted(artifact VerificationArtifact) error

// Resolve — public. Coordinator reads only.
func (r *ArtifactRegistry) Resolve(ctx context.Context, ref string) (VerificationArtifact, error)

// VerifyHash — public. Coordinator reads only.
func (r *ArtifactRegistry) VerifyHash(ref string, expectedSHA256 string) error
```

Key: `registerTrusted` is unexported. `VerificationArtifact` type must be defined here (generic), not in physics/v1.

## Step 2: `oracle/verifier/physics/v1/artifact.go` — Artifact Types

```go
package physicsv1

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

Hash determinism:
- `ArtifactHash` = SHA-256 of canonical JSON excluding `ArtifactHash` and `Timestamp`
- Deterministic serialization: sort keys, no whitespace, no Timestamp

## Step 3: `oracle/verifier/physics/v1/proof.go` — Symbolic Engine + 4 Proofs

**Symbolic expression tree (bounded):**
```go
type Expr interface { expr() }
type Var struct { Name string }
type Const struct { Value float64 }
type Add struct { Left, Right Expr }
type Mul struct { Left, Right Expr }
type Pow struct { Base Expr; Exp float64 }
type Div struct { Num, Denom Expr }
type Sqrt struct { Arg Expr }
type Grad struct { Arg Expr }
type Laplacian struct { Arg Expr }
```

**Required transformations:**
- Constant folding
- Power rules
- Basic algebraic identities
- Gradient/Laplacian handling for the specific identity

**4 proof obligations:**

| Step | CheckType | What |
|------|-----------|------|
| 1 | `machine_verified` | Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(\|∇ρ\|²/ρ²) — symbolic construction + simplification + equality check |
| 2 | `machine_verified` | a(ρ) = κ²/(8mρ) — coefficient matching |
| 3 | `human_attested` | \|∇ρ\|² coefficient consistency — expected-value assertion |
| 4 | `human_attested` | b′(ρ) = 0 — expected-value assertion |

## Step 4: `oracle/verifier/physics/v1/verifier.go` — Orchestrator

```go
package physicsv1

func Run(ctx context.Context, registry *verifier.ArtifactRegistry, input VerifierInput) (*VerificationArtifact, error)
```

Flow:
1. Build symbolic expressions from input
2. Run steps 1–4
3. Aggregate result
4. Compute ArtifactHash
5. Register via `registry.registerTrusted()`
6. Return artifact

## Step 5: Tests

**Registry tests (~10):**
- valid registration
- duplicate reference rejection
- unknown reference rejection
- correct hash accepted
- wrong hash rejected
- concurrent reads safe
- concurrent registrations safe

**Artifact tests (~6):**
- canonical hashing deterministic
- timestamp excluded from hash
- hash reproducible across 100 runs
- malformed artifact rejected
- invalid result rejected
- invalid check type rejected

**Proof tests (~8):**
- Step 1: known-good identity → verified
- Step 1: modified identity → refuted/inconclusive
- Step 2: known-good coefficient → verified
- Step 2: modified coefficient → refuted/inconclusive
- Step 3: expected attestation → pass
- Step 3: malformed → fail
- Step 4: expected attestation → pass
- Step 4: malformed → fail

**Aggregation tests (~3):**
- all satisfied → confirmed
- any refuted → refuted
- unresolved machine check → inconclusive

**Trust boundary tests (~2):**
- registration API not exported
- external package cannot call registerTrusted

## Verification

```bash
cd /home/chaschel/Documents/go/oracle
go test ./...
go test -race ./...
go vet ./...
```
