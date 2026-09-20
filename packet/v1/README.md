# EBP Research Packet v1

## Role

`schema.json` is the canonical interchange contract for the EBP Research Packet v1. It defines structural constraints for JSON serialization.

Go code in this package performs executable structural + semantic validation against this schema. No external JSON Schema engine is used.

## Reference Format

- `local:<id>` — packet-local object reference
- `canonical:belief:<uuid>` — canonical epistemic object reference (syntax validated here; state resolution deferred to Coordinator/Phase 6)

## Pack Reference

Pack references use string format `"bmist-v1"` (pack_id-version). The PackRegistry resolves them at validation time.

## Validation Order

1. Schema version check
2. packet_id non-empty
3. role ∈ {work, adversarial}
4. pack_ref resolves through registry
5. Belief structure (local_id unique, claim non-empty, claim_type from pack)
6. Evidence structure (local_id unique, belief_ref resolves, content_sha256 valid hex, provenance_class from pack)
7. Debt vocabulary membership
8. All references resolve
9. Edge structure (from_ref/to_ref resolve, kind ∈ {derives, contradicts})
10. Task structure (local_id unique, title non-empty, no forbidden status)

## Invariants

- Phase 3 validates reference syntax. Phase 6 resolves canonical state.
- Phase 3 never acquires a Solvent dependency.
- One packet belief is one atomic epistemic object.
- Tasks cannot carry authoritative operational state.
