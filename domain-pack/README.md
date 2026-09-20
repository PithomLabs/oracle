# Domain Pack

A Domain Pack is a declarative, validated, versioned artifact that defines
domain-specific semantics for the ARGUS trust verification substrate.

## What a Domain Pack Declares

- **Claim types** — the kinds of claims the domain supports
- **Evidence classes** — the kinds of evidence the domain accepts
- **Debt vocabulary** — the obligation identifiers the domain uses
- **Initial debt** — the obligations assigned to new claims
- **Retirement rules** — how evidence retires specific obligations
- **Falsifiers** — the attack vocabulary for adversarial review
- **Consequential actions** — actions that require promoted belief state
- **Human-gated transitions** — transitions that require human adjudication

## What BM-IST Declares

The BM-IST Domain Pack (`bmist/v1`) defines:

- 3 claim types: derived, accommodated, postulated
- 2 evidence classes: reproducible_artifact, operator_asserted
- 6 debt items: needMap, needInvariant, needToyCheck, needNullModel, needObstruction, needFaithfulnessReview
- 4 falsifiers: counterexample, contradiction, missing_evidence, alternative_explanation
- 1 consequential action: publish_claim (requires promoted, gates faithfulness_review)
- 3 human-gated transitions: faithfulness_review, scope_clarification, obstruction_assessment

## What the Registry Does

The `PackRegistry` is an in-memory, startup-loaded registry that:

- Loads packs from the `domain-pack/` directory tree
- Validates packs before registration
- Keys packs by `pack_id@version`
- Provides deterministic lookup via `Get(packID, version)`

## What the Registry Does NOT Do

- The registry does NOT interpret pack semantics
- The registry does NOT enforce debt retirement
- The registry does NOT promote beliefs
- The registry does NOT make policy decisions
- The registry does NOT persist data across restarts

## Relationship to Other Components

### Solvent

Solvent is domain-agnostic. It does not interpret debt strings, evidence
classes, or pack semantics. Solvent enforces structural invariants
(empty-vs-non-empty debt for promotion) but not domain-specific meaning.

### Coordinator

The Coordinator consumes Domain Packs. It uses the pack's debt vocabulary,
retirement rules, and human-gated transitions to compile research packets
into Solvent beliefs and evidence.

### EBP (Evidence-Based Planning)

Domain Packs are the configuration layer for EBP. They declare the
domain-specific obligation vocabulary that EBP uses to track research debt.

## Domain Pack Semantics

Domain Pack semantics are declarative. They are configuration + schema,
not an execution engine. The Pack does not contain:

- Verifier logic
- Coordinator logic
- Solvent API calls
- Agent behavior
- Policy decisions

The Pack is data. The Coordinator and Solvent are code. The Pack tells
the code what domain-specific vocabulary to use; it does not tell the
code how to reason.

## Agents Do Not Gain Authority from the Pack

A Domain Pack does not authorize agents, grant permissions, or elevate
privileges. Agents produce research packets; the Coordinator compiles
them; Solvent records the epistemic state. The Pack defines vocabulary,
not authority.
