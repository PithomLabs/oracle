Yes. I would add this as a **small protocol layer shared by both agent roles**, rather than duplicating the rules in the Work and Adversarial prompts.

The agent identity should be explicit and supplied by the human/agent-launch configuration, not inferred by ARGUS. The identity is metadata for provenance and observability; it is **not authority**.

Use this prompt for the coding agent:

````text
# ARGUS — Agent Identity + Common Agent Protocol

Implement two tightly related improvements to the current ARGUS agent workflow:

1. Explicitly identify the AI agent/harness/model performing Work or
   Adversarial research, and expose that identity in ARGUS state/UI.
2. Create one shared protocol of agent Do's and Don'ts that both Work and
   Adversarial agents must follow.

Do not redesign ARGUS.
Do not add new MCP tools.
Do not change Solvent authority semantics.
Do not add a new database.
Do not infer agent identity from the model or process automatically when the
identity is available explicitly from the agent's startup instructions.

The current ARGUS implementation is authoritative. Audit first.

==================================================
1. AUDIT CURRENT AGENT METADATA PATH
==================================================

Inspect the current code for:

- packet/v1/types.go
- packet/v1/schema.json
- packet/v1/validate.go
- packet compiler/persistence
- internal/application/
- internal/work/
- internal/epistemic/
- internal/ui/
- existing `current_agent` task field
- existing packet `agent` object, if present
- existing agent-related UI rendering
- prompts/opencode-work.md
- prompts/opencode-adversarial.md

Determine:

A. Does the current EBP packet already support:

```json
"agent": {
  "id": "...",
  "model": "...",
  "role": "work|adversarial"
}
````

B. Where is the agent identity currently persisted?

C. What does `conductor_task.current_agent` actually represent?

D. Why does `/ui/insights` currently render:

```text
Agent = <nil>
```

E. Is agent identity already accepted but simply not supplied by the current
prompts?

Do not add duplicate fields if an existing authoritative field already serves
the purpose.

==================================================
2. DEFINE THE AGENT IDENTITY CONTRACT
=====================================

Use one explicit identity structure.

Minimum conceptual fields:

```text
agent_id
role
harness
model
```

Example:

```json
{
  "agent_id": "work-001",
  "role": "work",
  "harness": "OpenCode",
  "model": "exact-runtime-model-name"
}
```

For adversarial:

```json
{
  "agent_id": "adversarial-001",
  "role": "adversarial",
  "harness": "OpenCode",
  "model": "exact-runtime-model-name"
}
```

The agent identity MUST be supplied explicitly by the agent startup
instructions/configuration.

Do not guess the model.

Do not infer the role.

Do not use the human operator principal as the agent identity.

Do not use agent identity as an authority credential.

Agent identity is provenance/observability metadata only.

==================================================
3. PACKET CONTRACT
==================

If the current EBP packet schema already has an `agent` object, use it.

If it does not, add the smallest compatible metadata field necessary.

Preferred:

```json
"agent": {
  "id": "work-001",
  "role": "work",
  "harness": "OpenCode",
  "model": "exact-model-name"
}
```

The packet compiler must preserve this metadata.

The persistence path must not silently discard it.

Do not create a second agent identity model elsewhere.

If the current packet schema already records role/model, extend rather than
duplicate.

Add validation:

* `agent.id` must be non-empty
* `agent.role` must be `work` or `adversarial`
* `agent.harness` must be non-empty
* `agent.model` must be non-empty

Do not make model identity part of any authority decision.

==================================================
4. TASK / UI ASSOCIATION
========================

The current `/ui/insights` task table contains:

```text
ID | Title | Status | Priority | Agent
```

The Agent column currently displays `<nil>`.

Fix this so that when the Work Agent or Adversarial Agent has actually
submitted work associated with the task, the UI displays the corresponding
agent identity.

First inspect what `current_agent` means in the existing task model.

Prefer the existing field if it is already intended for this purpose.

Do not introduce an unrelated new database column if the current model can
represent the information.

The displayed identity should be useful to a human.

A compact form is acceptable, for example:

```text
work-001 / OpenCode / model-name
```

or:

```text
OpenCode:model-name [work-001]
```

Use whichever fits the existing UI cleanly.

Do not display `<nil>` when valid submitted agent identity exists.

If no agent has acted on a task yet, an empty/`—` representation is preferable
to `<nil>`.

==================================================
5. IMPORTANT: DO NOT CONFUSE "CURRENT AGENT" WITH AUTHORITY
===========================================================

Agent identity does NOT mean:

* authorized actor
* human operator
* principal
* authority holder
* debt discharger
* promoter
* retractor

The existing authority model remains unchanged.

Human authority continues through the Trust UI/application/Solvent path.

Agents remain limited to:

```text
argus.get_context
argus.submit_packet
```

==================================================
6. COMMON AGENT PROTOCOL
========================

Create one shared protocol file:

```text
prompts/argus-agent-protocol.md
```

Do not create separate duplicated rule sets.

The Work and Adversarial prompts should reference this common protocol.

The protocol should establish the common rules for EVERY ARGUS research
agent.

Use this substance:

---

## ARGUS AGENT PROTOCOL

### Identity

Before doing research, the agent must know and state:

```text
agent_id: <explicit identifier>
role: work | adversarial
harness: <agent harness>
model: <exact runtime model>
```

Never invent or guess these values.

Use the declared identity consistently in submitted packets.

### Required first actions

1. Read the assigned task/instructions.
2. Call `argus.get_context`.
3. Reconstruct current ARGUS state.
4. Read the available research background under `docs/corpus/`.
5. Distinguish live ARGUS state from background research documents.

### Agents MUST

* use `argus.get_context` for live research state
* use `argus.submit_packet` for research-state mutations
* preserve existing belief IDs when referring to existing claims
* create new belief IDs for substantively new/revised propositions
* make evidence provenance explicit
* make uncertainty explicit
* preserve contradictions rather than silently resolving them
* distinguish hypotheses from established results
* challenge unsupported assumptions
* keep reasoning auditable
* use the existing EBP packet grammar
* respect the active Domain Pack
* treat validation failures as information about the protocol
* maintain the distinction between work and authority

### Agents MUST NOT

* access or modify the ARGUS database directly
* inspect database tables to reconstruct research state
* modify Solvent authority state directly
* promote beliefs
* discharge debt
* retract beliefs
* authorize consequential actions
* mutate Conductor state directly
* claim human attestation
* fabricate evidence
* classify an agent-generated assertion as `operator_asserted`
* treat reading a document as automatically making it authoritative evidence
* overwrite an existing belief's claim to "correct" it
* silently delete contradictory research
* treat absence of a counterexample as proof
* treat absence of evidence as evidence of truth
* bypass packet validation
* add new authority mechanisms

### Belief lifecycle rule

Existing beliefs are historical propositions.

Do not rewrite an existing belief's claim in place.

When research produces a materially different proposition:

```text
old belief
    +
new research
    ↓
new belief ID
```

Connect the new belief using existing graph semantics such as:

```text
derives
contradicts
```

as appropriate.

The historical belief remains available for inspection.

### Adversarial rule

An adversarial agent must challenge existing work rather than silently replace
it.

A challenge should identify:

* the target claim
* the specific problem
* supporting reasoning/evidence
* uncertainty
* possible consequence

The adversarial agent must not adjudicate its own challenge.

### Human authority rule

Agents may propose evidence or debt retirement.

Agents do not decide whether the debt is discharged.

Agents do not promote claims.

Agents do not retract claims.

Those are human/authority-layer operations.

### Provenance rule

Background documents are background context.

A document being read does not automatically become authoritative evidence.

Use the actual evidence/provenance requirements of the current ARGUS packet
validator.

### Unknown rule

Never convert:

```text
UNKNOWN
```

into:

```text
EMPTY
```

or:

```text
TRUE
```

because information was unavailable.

---

## END ARGUS AGENT PROTOCOL

Keep the protocol concise. It is a shared behavioral contract, not a research
manual.

==================================================
7. UPDATE EXISTING ROLE PROMPTS
===============================

Do NOT create new Work or Adversarial prompt variants.

Update:

```text
prompts/opencode-work.md
prompts/opencode-adversarial.md
```

to reference:

```text
prompts/argus-agent-protocol.md
```

Each prompt should say:

```text
Read and follow prompts/argus-agent-protocol.md.
```

Then retain role-specific behavior.

The Work Agent adds:

* decomposition
* evidence gathering
* candidate claims
* useful research

The Adversarial Agent adds:

* independent challenge
* counterexamples
* assumption attack
* contradiction construction

Do not duplicate the common rules into both files.

==================================================
8. AGENT STARTUP IDENTITY BLOCK
===============================

Update the role prompts so the human can supply identity explicitly near the
top.

Example Work Agent invocation:

```text
AGENT_ID: work-001
AGENT_ROLE: work
AGENT_HARNESS: OpenCode
AGENT_MODEL: <exact runtime model>
```

Example Adversarial Agent:

```text
AGENT_ID: adversarial-001
AGENT_ROLE: adversarial
AGENT_HARNESS: OpenCode
AGENT_MODEL: <exact runtime model>
```

The prompt must tell the agent:

> Do not invent these values. Use the exact identity supplied by the operator.

The implementation should carry this metadata into `submit_packet`.

Do not make the system attempt to infer the model from OpenCode internals.

==================================================
9. UI EXPECTATION
=================

After a Work Agent submission, `/ui/insights` should show an agent identity
instead of:

```text
<nil>
```

Example:

```text
Agent
OpenCode:model-name [work-001]
```

After an adversarial submission, the resulting state/activity should allow the
human to distinguish:

```text
work agent
vs.
adversarial agent
```

Do not redesign the UI.

Do not add a new dashboard.

Modify only what is necessary to expose the metadata already being carried by
the work.

==================================================
10. TESTS
=========

Add focused regression tests.

### Identity propagation

Verify:

1. Work packet declares agent identity.
2. Submit succeeds.
3. Agent identity is persisted/preserved.
4. UI or dashboard projection returns the identity.

### Adversarial identity

Verify:

1. Adversarial packet declares different role/identity.
2. Submit succeeds.
3. The adversarial identity remains distinguishable from Work Agent identity.

### Missing identity

Verify the desired behavior if identity is required by the current packet
contract.

Prefer fail-closed validation over silently persisting anonymous work if the
protocol requires explicit identity.

### Authority isolation

Verify agent identity cannot be used as a human operator principal or as
evidence of authority.

### Existing behavior

All existing tests must remain green.

==================================================
11. DO NOT BROADEN SCOPE
========================

Do NOT:

* add MCP tools
* add a new agent service
* add a new database
* add generic authentication
* infer model identity
* create a separate provenance database
* create agent accounts/principals
* change Solvent authority semantics
* redesign the Trust UI
* change the EBP model beyond the minimum metadata needed
* add automatic debt adjudication

The objective is:

```text
explicit agent identity
        +
shared behavioral protocol
        +
existing EBP workflow
```

not a new identity-management subsystem.

==================================================
12. FINAL REPORT
================

Report:

1. Where agent identity already existed.
2. What was missing.
3. Exact files changed.
4. Exact packet/schema changes.
5. How identity reaches the UI.
6. Common protocol location.
7. How both role prompts reference it.
8. Tests added.
9. Full test results.
10. Example Work Agent identity.
11. Example Adversarial Agent identity.

Then demonstrate the final flow:

```text
Work Agent
  identity
    ↓
argus.submit_packet
    ↓
ARGUS
    ↓
Insights

Adversarial Agent
  different identity + role
    ↓
argus.submit_packet
    ↓
ARGUS
    ↓
Insights / epistemic graph
```

The architectural invariant remains:

```text
CAPABILITY != WORK != AUTHORITY != EXECUTION
```

Agent identity improves provenance and observability.
It must never become authority.

```
```
