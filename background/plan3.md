## Consolidated Deliverables — Physics Verifier as First Domain Pack

### P0 — Prerequisites

1. **Plan 11.1 executed and re-frozen**

   * New kernel hash pinned.
   * Freeze reconciliation/decision record completed, including deferred H5 items.

### A. Domain-Pack Contract

2. **Physics Domain Pack spec**

   * Declarative only.
   * Claim/evidence types, EBP six debt vocabulary, discharge semantics, falsifiers/challenges.
   * Consequential-action declarations and gate map.
   * Human-gated transitions explicitly marked.

3. **EBP Research Packet v1**

   * `belief / evidence / debt / edge / task`.
   * Local vs canonical references.
   * Deterministic mapping to Solvent/Conductor ownership and enforcement.

### B. Deterministic Verification Substrate

4. **Physics Verifier Runner**

   * Symbolic/algebraic/toy-model verification.
   * Reproducible artifacts.
   * Hash-pinned verifier/version/input metadata.
   * No LLM authority.

5. **Corpus Manifest**

   * Hash-pinned BM-IST seed + Qwen artifact.
   * Provenance class and stable locator per artifact.

### C. Coordinator + Ledgers

6. **Go Coordinator**

   * Packet validation.
   * Mandatory `ebpInitialDebt` attachment + negative regression test.
   * Compile to Solvent + Conductor.
   * `packet_id` idempotency.
   * Solvent-before-Conductor ordering.
   * Generic Conductor projection via `governance_ref`.
   * Human decision endpoint using `operator_asserted` provenance.

7. **Human Adjudication Loop**

   * Retire/reopen/merge/retract through Coordinator → Solvent.
   * Branch pruning/retraction is consequential and must use Solvent action-intent semantics.
   * No agent may promote, retire debt, or adjudicate.

### D. End-to-End POC Run

8. **Dual-agent run**

   * One WORK agent + one ADVERSARIAL agent.
   * Fresh processes, distinct model families.
   * Pre-registered pass/fail/falsification criteria.

9. **POC acceptance evidence**

   * Malformed packet rejected.
   * Incomplete adversarial coverage detected.
   * Zero agent-caused protected-state mutations.
   * Belief state reconstructable by a human from Solvent/Conductor projections.
   * Promotion → authorized intent demonstrated.
   * Contradiction/retraction → intent invalidation demonstrated.

### E. Domain-Neutrality Proof

10. **Substrate conformance artifact**

* Zero physics/domain vocabulary in Solvent and Conductor.
* Pack interacts only through packet/API contracts.
* Automated conformance test proving the pack can run without substrate changes.

**Do not rely on grep alone; grep is a static guard, while the conformance test is the actual architectural proof.**

### F. Explicit Non-Goals

11. Defer:

* signed attestation / DSSE;
* W3C PROV export;
* second domain pack;
* probabilistic inference;
* pack-interface freeze;
* new Solvent/Conductor domain primitives.

Those belong to the later roadmap, not this POC.

### G. Success Claim

12. **Honest claim**

* One domain pack demonstrates **substrate compatibility**, not empirical domain agnosticism.
* The stronger domain-agnosticism claim is validated by **Domain Pack #2** without kernel changes.

### Sequence

```text
1 → 2–3 → 4–5–6 → 7 → 8–9 → 10 → 11
```

The most important addition from my side is **the deterministic Physics Verifier Runner**: otherwise the POC proves packet/workflow plumbing more than a generalized trust-verification capability.

