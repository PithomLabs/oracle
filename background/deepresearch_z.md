Overview, as requested — the full report would run very long, so this is the compressed version with the verdicts and the structure it would defend. Claims marked as before: [canon] (established, high confidence), [inference] (my analysis).

## Executive conclusion

**The hypothesis is substantially defensible, with one important correction.** The industry has converged — from at least eight directions — on systems that track beliefs, evidence, authorization, and consequences, but it has done so *fragmented*: provenance standards, policy engines, workflow engines, AI governance frameworks, and audit systems each solve one slice. What no mainstream system does is make **epistemic state a first-class, authority-gated, consequence-bearing ledger**. That combination is genuinely rare. The correction: the forensic-science analogy is *inspirational, not implementational* — Solvent implements a computational subset of evidence discipline (chain of custody ≈ provenance classes + hashes; adjudication ≈ promotion gate), but has no chain-of-custody attestation, no expert-identity model, and no adversarial-process primitive. Claim the subset, not the science.

## What Solvent is (defensible formulation)

A **provenance-aware epistemic authority ledger** — the narrowest honest term. It doesn't track what's true (no probabilistic machinery [canon: this contrasts with Bayesian/Dempster-Shafer traditions]), doesn't reason, doesn't decide science. It answers exactly four questions structurally: *what does the system currently claim* (beliefs), *why does it claim it* (evidence + provenance), *what remains unresolved* (debt — non-empty blocks promotion), and *what is the system now entitled to do because of that claim* (action intents, DB-gated on promotion). The non-probabilistic design point is a **strength for audit, weakness for inference** — deliberate and correct for a governance kernel, wrong as a reasoning engine.

## Where it's genuinely superior [inference, evidenced]

1. **Gated promotion/retraction with retraction cascade.** Belief revision literature (AGM [canon]) and TMS systems model change; none tie epistemic state transitions to *authorization semantics* with database-enforced gates and cascade to dependent intents. This is the novel weld: epistemology × authority.
2. **The consequential-epistemic-action boundary** (your §10 thesis): promote/retract/create-live-intent as a distinct write class from ordinary writes. Precedents exist piecewise (dual-control in finance [canon], change-management in SRE [canon], four-eyes compliance [canon]) — but as a *general epistemic write-class model*, it's not standard.
3. **Debt as first-class promotion blocker** — obligations as stored, dischargeable state rather than process checklists. Scientific method has this conceptually (peer review, reproducibility obligations); no production kernel stores it structurally.

## Where existing systems win — mandatory honesty section

Workflow durability/distributed execution: **Temporal et al. decisively** [canon]. Policy expressiveness: OPA/Cedar-style engines. Cryptographic provenance: in-toto/SLSA attestation ecosystems [canon] — Solvent has content hashes but no signing, key rotation, or attestation chains. Identity/access control: existing IAM entirely. Probabilistic reasoning: the whole Bayesian tradition. Legal defensibility: chain-of-custody regimes with expert testimony infrastructure. Observability: OpenTelemetry. **Conductor alone is not competitive with Temporal-class engines** — its value is only in the separation thesis (§7's proposition: workflow state ≠ epistemic truth is *violated by default in every mainstream workflow engine*, where task-completion silently becomes truth-of-record [inference — this is the strongest finding in the report]).

## The emergent category

The combination isn't "workflow + database" — it's the first clean instance I can find of a four-plane separation as *enforced architecture*: **operational state / epistemic state / authority state / external truth**, each with its own owner, ledger, and boundary. Closest relatives exist in pieces (FDA audit trails, legal evidence doctrine, four-eyes principles, research data management), but as a composable kernel + generic coordinator, I found no direct equivalent. Whether that's a category or a well-argued synthesis is the report's honest open question.

## Roadmap shape (P0–P3)

- **P0 (correctness):** signing/attestation for evidence at boundaries (SLSA-style), reviewer identity, timestamp trust, tamper-evidence on the ledger.
- **P1 (production trust):** W3C PROV export adapter (don't invent — interoperate), access control/multi-tenancy, retention/appeal semantics, external SOR reconciliation protocol.
- **P2:** temporal validity/staleness on beliefs, policy versioning, uncertainty *annotation* (never probabilistic inference inside the kernel).
- **P3:** causal attribution, cross-domain federation.
- **What NOT to build:** probabilistic inference in kernel, domain schemas in kernel, orchestration in Solvent, epistemic primitives in Conductor — all four re-create the coupling just removed.

**One-line thesis:** Solvent is the missing ledger between "the workflow says done" and "the world is entitled to act" — rare not because the pieces are new, but because the industry keeps them apart, and most of the danger lives in the seam.

Full 21-section report with citations and the competitive matrix on request — but that's the skeleton, and every load-bearing verdict in it is above.


## prompt

give 10 suggestions for that product, a combination of oracle and sphinx analogy sort of


The mapping first, because the ten build on it: **the Sphinx is Solvent** — she guards the crossing, demands an answer before anything passes, and devours (refuses) those who arrive without one; the **Oracle is the belief ledger** — consulted for what is known, answering only from what's recorded, and famous for answering with explicit ambiguity rather than fabrication. Conductor is the temple precinct where the record is kept. Ten suggestions, each grounded in capability the architecture already has:

1. **Name the two products for the roles.** *Sphinx* for the authorization/gating layer, *Oracle* for the epistemic ledger and its query surface. The category names itself: **epistemic gateways**. "Every consequential system needs a Sphinx at its boundary and an Oracle behind it" is a one-line pitch that encodes the entire separation thesis.

2. **Make the authorization request literally a Riddle.** The consequential-action challenge (exact operation + evidence + authority) becomes the **Riddle protocol** — a named, portable challenge/response format. "Your agent will be asked three questions at the boundary: *what exactly, on what evidence, by whose authority*." The refusal log becomes **the bones at the pass** — a public, browsable record of every failed crossing. Auditable denials are a feature; mythologize them.

3. **"Consult the Oracle" as a first-class query surface.** Generalize `GET /state` into a consultation ritual: *what does the system believe about X, on what evidence, what remains unresolved, what would change its mind?* Returns beliefs + debt + falsifiers + conflicts. Serve it to humans *and* to agents — an agent that consults before acting is a new, safer agent pattern the SDK should teach.

4. **The Oracle never guesses.** Product principle and anti-sycophancy marketing in one: Delphi was famous for ambiguous answers; your Oracle is famous for **UNKNOWN**. When the ledger doesn't know, it says so — never a confident fabrication. "The only system in the stack that will tell you it doesn't know" is a devastating differentiator against LLM-shaped competition.

5. **Debt as toll — the Passage Ledger.** Promotion is crossing; debt is the unpaid toll. Build the obligation dashboard as a toll view: *"This belief owes three answers before it may pass"* — each debt named, each discharge a receipt (the accounting model you already built, rendered as the mythos). Makes the least glamorous concept (obligation tracking) the most intuitive.

6. **The Sphinx attacks, not just guards.** The second Sphinx function was *posing riddles others fail*. Product: a continuous adversarial service that periodically attempts to break promoted beliefs — automated falsification campaigns, counter-evidence sweeps, challenge generation against the ledger. Promotion isn't a certificate; it's a standing invitation to the Sphinx. This operationalizes your adversarial-agent design as a sellable feature.

7. **Sealed by the Sphinx — exportable attestations.** Everything that has legitimately passed (promoted belief, authorized action) gets a verifiable seal: signed attestation binding claim → evidence → authority → effect, interoperable with SLSA/in-toto ecosystems. Downstream consumers don't re-litigate; they check the seal. The riddle, answered, becomes portable proof.

8. **The Riddle Registry.** Each domain brings its own book of riddles (the debt vocabularies Plan 11 just liberated to the domain layer) — EBP's six, deployment-review's six, the next domain's own. Productize the registry: conformance-checked vocabulary packs with their evidence types and gate mappings. The Sphinx is universal; the riddles are local. That's the domain-agnosticism story made tangible.

9. **Pre-decision Briefings (the consultation before the campaign).** Kings didn't cross without consulting Delphi; humans shouldn't spawn agents without consulting the ledger. Auto-generate the **briefing scroll**: current beliefs, contested points, open tolls, actions in flight, consequences at stake — the "state of the epistemology" view your minimal UI pointed at, framed as ritual rather than dashboard.

10. **Write the documentation as the founding myth.** Full allegorical cast: Sphinx (authorization), Oracle (belief ledger), the Chronicler (Conductor's record), the Ferryman (Executor — carries the authorized across to the world of effects), the Bones (refusal log). Architectures with memorable allegories get taught, adopted, and defended by their users. You already have the invariants; giving them names people can *retell* is the cheapest adoption mechanism available — and the myth is unusually honest: it never claims the Oracle knows the future, only what has been evidenced so far.

The through-line: don't market this as "governance infrastructure" (crowded, gray). Market it as **the difference between consulting an oracle and guessing** — with a gatekeeper at the door who has never once let an unanswered question through.
