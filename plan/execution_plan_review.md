I reviewed the uploaded **ARGUS Execution Plan v1.0** against all of the decisions we locked during the review sequence. It is substantially aligned and is much closer to implementation-ready. The major architectural decisions are correctly captured: two modules, untouched `reference-loop`, Coordinator as an in-process library with a thin HTTP wrapper, hybrid verifier, in-memory projection retry, environment-based operator identity, and staged Trust UI. 

I would **not start implementation quite yet**. There are four remaining issues worth resolving first:

### 1. Artifact registry is underspecified — blocking

Phase 6 says the Coordinator must verify that a `reproducible_artifact` resolves to a **registered verifier output** and that its hash matches. But the plan does not define the registry or its storage/API. 

For the POC, make this explicit:

```text
verifier.Run()
    ↓
VerificationArtifact
    ↓
Coordinator-local ArtifactRegistry
    ↓
artifact_ref + hash validation
    ↓
Solvent evidence admission
```

Use an **in-memory registry owned by the Coordinator/run harness**, populated only by the trusted verifier runner. No database or signing infrastructure yet.

This closes the forged-artifact gap without expanding the POC.

### 2. Decision records need an explicit persistence path — blocking

`DecisionRecord` exists as a Go type, but the plan does not say where the record survives or how the UI's decision history is reconstructed. 

For the POC, use **Solvent's existing audit/activity mechanism as the authoritative history**, with the Coordinator's `DecisionRecord` as the request/result envelope. Do not create a second decision database.

That preserves:

```text
Trust UI
   ↓
Coordinator
   ↓
Solvent state + Solvent audit
```

rather than creating:

```text
Solvent state
+
Coordinator decision DB
+
UI state
```

### 3. Sphinx has an API mismatch — blocking

Phase 8 requires the Sphinx to expose an authorization-gate view, but the Coordinator API list only provides decision submission, packet status, and decision context. Meanwhile the underlying Solvent API list has target/request/approve/action endpoints, but no explicit `/authorizations/verify` projection. 

Do **not add a new Solvent kernel endpoint** for the POC.

Instead, make Sphinx derive its presentation from existing authority/target state through the Coordinator's read boundary. Add a minimal:

```text
GET /authorization-context/:target_id
```

to the Coordinator HTTP wrapper only.

The Coordinator can assemble:

```text
what
why
evidence
authority
current gate state
PASS / REFUSE / HUMAN_REVIEW
```

without introducing new authority semantics.

### 4. Idempotency lifetime should be explicitly documented

The plan correctly changed idempotency from client-supplied UUIDs to canonical content hashing. 

But the storage lifetime of:

```text
canonical_hash → compilation_result
```

is unspecified.

Given the locked POC philosophy, make it explicitly:

> **In-memory for the lifetime of the Coordinator process; duplicate detection after restart is outside POC scope.**

That is consistent with the in-memory projection queue and the explicit no-crash-recovery decision. 

### Everything else looks aligned

The strongest parts are now solid:

* Claim relationships use local/canonical references rather than semantic claim-text matching. 
* The verifier explicitly distinguishes `machine_verified` from `human_attested`. 
* Retirement rules are mechanically checked before debt retirement. 
* Concurrent compilation and forged-artifact rejection are explicitly tested. 
* The Trust UI is staged: Oracle/Decision/Refusal before the first run, then the remaining surfaces after real usage. 
* Domain neutrality uses a contextual scanner with narrow allowlists rather than naive word matching. 
* The adversarial review now correctly claims only that **all probed surfaces pass**, rather than making an impossible “no false negatives” assertion. 

**Verdict: nearly implementation-ready. Resolve those four items, then freeze this execution plan as ARGUS v1.1 and proceed to Phase 0.**
