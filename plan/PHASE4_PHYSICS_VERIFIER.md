# PHASE 4 PHYSICS VERIFIER — COMPLETE

## Status: PASS

## Verifier Scope

Bounded BM-IST physics verifier implementing exactly 4 proof obligations.

## Proof Obligations

| Step | CheckType | Description |
|------|-----------|-------------|
| 1 | `machine_verified` | Δ√ρ/√ρ = ½(Δρ/ρ) − ¼(|∇ρ|²/ρ²) |
| 2 | `machine_verified` | a(ρ) = κ²/(8mρ) — coefficient matching |
| 3 | `human_attested` | |∇ρ|² coefficient consistency |
| 4 | `human_attested` | b′(ρ) = 0 |

## Symbolic Approach

Bounded symbolic expression tree with:
- Expr types: Var, Const, Add, Mul, Pow, Div, Sqrt, Grad, Laplacian
- Simplification: constant folding, power rules, basic algebraic identities
- No full CAS — bounded to specific proof obligations

## Human-Attested Boundary

Steps 3–4 are explicitly `human_attested`. The artifact preserves per-step distinction:
- `confirmed` result does NOT imply all steps were machine-derived
- Each step records its own `CheckType`

## Artifact Structure

```go
type VerificationArtifact struct {
    RunID, VerifierVersion, VerifierHash, ClaimHash, InputHash, ArtifactHash string
    Timestamp time.Time
    Steps     []Step
    Result    string  // "confirmed" | "refuted" | "inconclusive"
    EvidenceRef string
}
```

## Hashing Approach

- `ArtifactHash` = SHA-256 of canonical JSON excluding `ArtifactHash` and `Timestamp`
- Deterministic serialization: sorted keys, no whitespace
- Same inputs → identical hash across 100+ runs

## ArtifactRegistry Design

```go
type ArtifactReader interface {
    Resolve(ctx, ref) (VerificationArtifact, error)
    VerifyHash(ref, expectedSHA256) error
}

type ArtifactRegistry struct { ... }
func NewArtifactRegistry() *ArtifactRegistry
func (r *ArtifactRegistry) registerTrusted(artifact) error  // unexported
func (r *ArtifactRegistry) AsReader() ArtifactReader
```

## Trust Boundary

- Neutral artifact model: `verifier/model/artifact.go`
- `registerTrusted` is unexported — only callable from within `verifier/` package
- `RunPhysicsVerifier()` is the ONLY trusted registration path
- Coordinator receives `ArtifactReader` interface (read-only)
- Physics verifier produces artifact; runner registers it
- No external package can register artifacts

## Package Structure (acyclic)

```
verifier/model/       ← neutral artifact types
verifier/registry.go  ← ArtifactRegistry, registerTrusted
verifier/runner.go    ← RunPhysicsVerifier (trusted path)
verifier/physics/v1/  ← proof logic, imports model only
```

Dependency flow:
```
physics/v1 → verifier/model
verifier/registry → verifier/model
verifier/runner → physics/v1, verifier/model, registerTrusted()
```

## Test Results

```
83 passed in 6 packages
go test -race ./... passed
go vet ./... passed
```

## Repositories Changed

- `oracle/verifier/model/` — neutral artifact types
- `oracle/verifier/` — registry, runner
- `oracle/verifier/physics/v1/` — proof logic

## Repositories Unchanged

- `solvent-main` — NO changes
- `conductor` — NO changes
- `oracle/reference-loop` — NO changes
- `oracle/domain-pack/` — NO changes
- `oracle/packet/` — NO changes

## Limitations

1. Symbolic verification is bounded — does not implement full CAS
2. Steps 1–2 produce deterministic output but full symbolic equality checking is simplified
3. Steps 3–4 rely on expected-value comparison (human-attested by design)
4. Registry is in-memory, non-persistent

## Deviations

1. `Equals` uses string comparison for simplicity in POC (not mathematically general)
2. Physics verifier returns artifact; runner in `verifier/` package registers it (avoids import cycle)

## Acceptance

- [x] Four proof obligations implemented
- [x] Steps 1–2 perform real symbolic checks
- [x] Steps 3–4 explicitly human_attested
- [x] Artifact output is deterministic
- [x] Artifact hash is reproducible
- [x] Registry is process-local and in-memory
- [x] Trusted registration is package-private
- [x] `RunPhysicsVerifier()` is sole trusted path
- [x] Forged/unknown/wrong-hash artifacts are rejected
- [x] No LLM calls in verifier
- [x] No external CAS introduced
- [x] Generic registry contains no physics vocabulary
- [x] go test ./... passes
- [x] go test -race ./... passes
- [x] go vet ./... passes
- [x] Solvent unchanged
- [x] Conductor unchanged
- [x] reference-loop unchanged
- [x] No Phase 5+ implementation exists
