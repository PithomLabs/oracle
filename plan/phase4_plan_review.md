I would **not implement this Phase 4 plan yet**. There is one important Go package-boundary error and two smaller design corrections.

### 1. Critical: `registerTrusted` cannot be called from `physics/v1`

The plan says:

```go
// oracle/verifier
func (r *ArtifactRegistry) registerTrusted(...)
```

but the runner is:

```go
// oracle/verifier/physics/v1
func Run(...)
```

`physics/v1` is a **different Go package**, even though it is a subdirectory. Therefore `registerTrusted` is inaccessible. The proposed Phase 4 implementation cannot compile as written. 

This also means the claimed compile-time trust boundary is not achievable with the current package layout.

### 2. Fix the artifact ownership model

The plan says the generic registry owns `VerificationArtifact`, while `physics/v1/artifact.go` also defines `VerificationArtifact`. That creates conflicting ownership. 

Use one canonical generic artifact type:

```text
oracle/verifier/
    artifact.go
    registry.go
```

Then:

```text
oracle/verifier/physics/v1/
    verifier.go
    proof.go
```

The physics verifier produces the generic `verifier.VerificationArtifact`.

### 3. Use capability passing rather than pretending package privacy solves it

The clean POC boundary is:

```text
Physics Verifier
    → trusted registration capability
    → ArtifactRegistry

Coordinator
    → ArtifactReader only
    → Resolve / VerifyHash
```

In other words, **Coordinator receives a read-only registry interface**, not the concrete registry with registration capability.

That is stronger architecturally than exposing `Register` and hoping callers behave.

Also remove the proposed test:

> “external package cannot call registerTrusted”

A non-compiling call cannot be a normal Go unit test. Make this an **API-surface/code-review acceptance criterion** instead. 

### Lock the correction

```text
Generic artifact type:
    oracle/verifier/artifact.go

Generic registry:
    oracle/verifier/registry.go

Physics verifier:
    oracle/verifier/physics/v1/

Coordinator:
    receives read-only ArtifactReader interface

Registration:
    available only through trusted verifier execution path

No duplicated VerificationArtifact types.
```

The rest of the Phase 4 scope—bounded symbolic steps 1–2, human-attested steps 3–4, deterministic hashing, and registry tests—is aligned with the frozen architecture. 

**Make this package-boundary correction before allowing the agent to implement Phase 4.**
