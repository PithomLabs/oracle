I reviewed **ARGUS POC Execution Plan v1.1**. It incorporates the four requested fixes correctly and is much closer to an executable baseline. The final lock is explicit and internally coherent on modules, Coordinator boundary, verifier scope, artifact trust, decision persistence, Sphinx projection, and process-local idempotency. 

I would make **five final corrections before freezing it**, though.

### 1. Fix the edge compilation boundary

Phase 6 currently says:

> “For each edge: record locally (Conductor doesn't have edges)” 

That is wrong for the architecture.

`belief_edge` is part of Solvent's canonical epistemic state. Contradiction and retraction demonstrations depend on those edges being persisted in Solvent.

Change to:

```text
For each edge:
    resolve from_ref/to_ref to canonical belief IDs
    POST/create the corresponding Solvent belief_edge
    never store epistemic edges only in Coordinator memory
```

This is important enough to classify as **blocking**.

### 2. Define `PackRegistry`

Phase 3 requires:

```text
pack_ref → registered pack
```

and Phase 6 has:

```go
packRegistry *PackRegistry
```

but the execution plan never defines `PackRegistry` or its files.

Add:

```text
oracle/domain-pack/registry.go
oracle/domain-pack/registry_test.go
```

For the POC, make it an **in-memory registry loaded from `domain-pack/**` at startup**. No persistent registry.

### 3. Resolve the REFUSE/audit inconsistency

The plan says REFUSE uses `refusal_log` "where applicable," but later tests require:

> `REFUSE creates Solvent audit_activity entry with refusal: true` 

Those cannot both be assumed without verifying the actual Solvent API/schema.

Do not invent a new refusal persistence mechanism.

The implementation plan should say:

> Use the existing Solvent refusal/audit mechanism exactly as exposed by the frozen API. The Coordinator must not create a parallel refusal store.

Then Phase 0/1 should verify the exact supported representation before implementation.

### 4. Fix the Phase 9 → Phase 11 sequencing

Phase 9 explicitly says the Sphinx, Passage, Challenge, and Chronicle surfaces are built **after the first run**. 

But Phase 11 requires those surfaces for the consequential-state demonstration. 

The dependency should therefore be:

```text
Phase 10
    ↓
Phase 9.2 UI expansion
    ↓
Phase 11 consequential demo
```

Otherwise Phase 11 requires UI that the plan says does not yet exist.

### 5. Tighten ArtifactRegistry trust boundary

The registry is correctly in-memory and trusted-path-only. 

But the API currently exposes:

```go
Register(ctx, artifact)
```

as a normal method. Any code holding the registry can register arbitrary artifacts.

For the POC, make registration explicitly internal, for example:

```go
registerTrusted(...)
```

and expose registration only through the verifier runner/package boundary.

That makes the architecture match the stated invariant:

```text
ONLY trusted verifier output may enter ArtifactRegistry
```

rather than relying on caller discipline.

### Verdict

After those five changes, I would freeze **v1.1 as the implementation baseline**.

The core architecture is now strong:

```text
Agent
  ↓ packet
ARGUS Coordinator
  ↓
Solvent ← canonical epistemic/authority state
  ↓
Conductor ← operational projection
  ↓
Trust UI ← read/projection surface
```

with the verifier artifact registry and Solvent audit trail remaining the authoritative trust boundaries. 
