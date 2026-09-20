The Phase 3 plan is mostly aligned with the frozen ARGUS architecture, but I would **not freeze it exactly as written**. There are four corrections worth making before implementation.

### 1. Do not split debt into a second top-level `debts` object

The plan defines:

```go
Belief.Debt []string
```

and also:

```go
Debt {
    LocalID string
    Items   []string
}
```

This creates two representations of the same epistemic obligation and invites divergence.

Keep the locked thin packet model:

```json
{
  "beliefs": [
    {
      "local_id": "b1",
      "claim": "...",
      "claim_type": "derived",
      "debt": ["needMap"]
    }
  ]
}
```

Remove the top-level `debts` field and the `Debt` type.

This preserves the principle that one packet belief is one atomic epistemic object.

### 2. Keep `pack_ref` as the explicit canonical reference form

The plan changes:

```json
"pack_ref": "bmist-v1"
```

to:

```json
"pack_ref": {
  "pack_id": "bmist",
  "version": "1.0.0"
}
```

The object is technically clean, but it diverges from the previously frozen Research Packet contract and introduces an unnecessary special reference format.

Use:

```json
"pack_ref": "bmist-v1"
```

and let `PackRegistry` resolve it.

The existing `canonical:belief:<uuid>` reference syntax should remain reserved for canonical epistemic objects.

### 3. Be precise about “JSON Schema validation”

The plan says:

> “No external dependencies — JSON Schema validation done in Go code, not a library.” 

That is technically inaccurate.

Go's standard library does not provide a JSON Schema validator. What you are actually proposing is:

```text
schema.json
    = canonical interchange/schema artifact

Go validator
    = executable structural + semantic validation
```

That is perfectly reasonable for the POC, but call it that. Do **not** claim that Go is executing JSON Schema validation unless an actual JSON Schema engine is used.

### 4. `canonical:belief:<uuid>` should remain a reference, not a lookup

The plan correctly separates parsing from external state resolution, but make the boundary explicit:

```text
Phase 3:
    validate canonical reference syntax

Phase 6:
    resolve canonical belief against Solvent
```

The current wording supports this, but it should be a hard invariant so Phase 3 never acquires a Solvent dependency. 

### Everything else is solid

The strongest parts are correct:

* local/canonical references replace dangerous claim-text matching;
* packet roles remain `work|adversarial`;
* Domain Pack supplies claim/debt/evidence vocabulary;
* tasks cannot carry authoritative operational state;
* free-text edge references are rejected;
* validation is fail-closed;
* packet code remains domain-neutral. 

### Recommended lock

```text
PACK REF:
    string: "bmist-v1"

BELIEFS:
    debt[] inside belief

TOP-LEVEL DEBTS:
    remove

SCHEMA:
    schema.json is canonical interchange contract
    Go performs executable validation

CANONICAL REFERENCES:
    syntax validated in Phase 3
    state resolution deferred to Coordinator/Phase 6
```

With those four changes, the Phase 3 design is clean and implementation-ready.
