The two reports point to a strong generalization path: **turn the current trust-physics verifier into a domain-agnostic verification substrate, while keeping domain semantics outside the kernel**. The defensible abstraction is a provenance-aware epistemic authority ledger, not a scientific reasoner or truth engine. 

### 1. Core abstraction — generalize the verifier

* Replace physics-specific concepts with generic primitives: **Claim, Evidence, Provenance, Obligation, Relation, State, Authorization, Consequence**.
* Verifier answers: **what is claimed, why, what remains unresolved, and what consequential transitions are permitted**.
* Do not claim it determines truth. 

### 2. Trust model — preserve the four-plane separation

* **Conductor:** work / workflow state.
* **Solvent:** epistemic + authority state.
* **Executor:** effect execution.
* **External SOR:** effect truth.
* Make this boundary the invariant of the generalized system. 

### 3. Domain neutrality — move semantics into Domain Packs

Create a domain adapter/configuration layer defining:

* claim types;
* evidence classes;
* debt/obligation vocabulary;
* promotion rules;
* falsification/challenge rules;
* consequential actions.

The kernel should never know `needMap`, physics assumptions, legal evidence types, etc. This directly follows the successful debt-vocabulary decoupling. 

### 4. Generalized trust verifier

Implement a generic verification pipeline:

```text
Input evidence
→ provenance validation
→ claim construction
→ obligation/debt evaluation
→ consistency/challenge checks
→ promotion gate
→ authorization decision
→ consequential intent
```

The verifier produces **auditable state transitions**, not a confidence score.

### 5. Consequential epistemic actions

Formalize a distinct mutation class:

* ordinary operational write → Conductor;
* epistemic write → Solvent;
* **consequential epistemic write** → gated Solvent transition.

Examples: promote, retract, satisfy final blocking obligation, create/cancel live action intent. The research supports this as one of the strongest potentially novel aspects of the architecture. 

### 6. Provenance / forensic foundation

Generalize the current trust verifier around:

* content hashes;
* source identity;
* provenance class;
* evidence lineage;
* immutable/tamper-evident history;
* explicit human adjudication.

But don't claim to implement full forensic chain-of-custody; the current system lacks producer attestation and expert-identity semantics. 

### 7. Attestation boundary

Add P0 capabilities outside or at the edge of Solvent:

* signed evidence;
* producer identity;
* reviewer identity;
* trusted timestamps;
* attestation receipts.

This addresses the largest current gap: **provenance ≠ authenticity**. 

### 8. Domain-specific challenge engines

Keep the kernel generic, but allow each Domain Pack to supply:

* falsifiers;
* contradiction tests;
* required evidence;
* adversarial checks;
* acceptance criteria.

Thus BM-IST becomes **one trust-verification domain**, not the architecture itself.

### 9. Conductor integration

Conductor should receive **research/work obligations**, not epistemic meaning:

* run verification;
* gather evidence;
* execute challenges;
* obtain human review;
* retry blocked work.

It should never decide that a claim is proven. The research explicitly identifies workflow durability as an area where mature engines still outperform the current Conductor concept. 

### 10. Interoperability

Do not invent a proprietary universal provenance format.
Add adapters/exporters for:

* W3C PROV;
* OpenTelemetry;
* signed artifact/attestation ecosystems;
* existing policy/IAM systems;
* external SORs.

The research specifically recommends interoperability rather than trying to replace established trust infrastructure. 

### 11. Missing capabilities to add incrementally

**P0:** attestation, trusted identity, timestamps, tamper evidence.
**P1:** provenance export, access control, retention/appeal, SOR reconciliation.
**P2:** temporal validity, stale-belief detection, policy-version binding, uncertainty annotation.
**P3:** causal attribution and federation. 

### 12. Explicit non-goals

Do **not** put these into Solvent:

* probabilistic inference;
* scientific reasoning;
* domain schemas;
* workflow orchestration;
* source reputation;
* autonomous adjudication;
* agent reasoning.

Those belong in domain/application layers or external systems. The report explicitly identifies these as places where expanding the kernel would recreate the coupling the architecture is trying to remove. 

### Resulting system

```text
                DOMAIN PACK
     ┌──────────────────────────────┐
     │ claims / evidence / debts    │
     │ tests / falsifiers / policy  │
     └──────────────┬───────────────┘
                    │
                    ▼
             TRUST VERIFIER
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
     SOLVENT                 CONDUCTOR
 epistemic/authority          workflow
        │                       │
        └───────────┬───────────┘
                    ▼
                 EXECUTOR
                    │
                    ▼
               External SOR
```

The strategic shift is therefore: **BM-IST is no longer the product architecture; it becomes the first Domain Pack running on the generalized trust substrate.** That is consistent with the strongest conclusion in the research: the defensible value is the missing layer between *“workflow says done”* and *“the system is entitled to act.”* 

