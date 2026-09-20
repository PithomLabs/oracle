I consolidated the two reviews against the actual Phase 8 plan. The second review correctly exposes several **integration-seam defects** that my previous review did not catch. Most of the architecture is now solid; the remaining issues are concrete and should be patched before implementation.

The plan already has the important foundation: RCP is a thin projection over Conductor + Solvent, OpenCode gets only `argus.get_context` and `argus.submit_packet`, persistence is real rather than compilation-only, and the test suite explicitly covers reconstruction, capability absence, EBP gating, UI, regression, and concurrency.   

## Final disposition

### C1 — Human discharge must be attributed in Solvent

**Accept. Must fix.**

The plan discovered both Solvent paths:

```text
POST /v1/beliefs/{id}/debt/retire
POST /v1/discharge
```

and the latter explicitly accepts `DischargedBy` and `InstrumentRef`. 

The human EBP path should therefore use:

```text
Trust UI
  ↓
Coordinator
  ↓
Solvent /v1/discharge
  ↓
audit_activity.actor_id
```

with:

```text
DischargedBy = configured operator principal
InstrumentRef = Coordinator decision/reference
```

The bare `retire` endpoint can remain an internal/mechanical primitive, but **human adjudication must use the attributed discharge path**.

This is necessary because Solvent's `audit_activity` explicitly requires `actor_id`. 

---

### C2 — Task ↔ belief linkage is currently broken

**Accept. Must fix.**

This is the most important newly discovered structural issue.

The plan says `governance_ref` is immutable after task creation, while the human-created task starts with no `belief_id`. Yet dead-end derivation later expects:

```text
task.governance_ref.metadata["belief_id"]
```

That linkage cannot be added afterward.

Do **not** weaken Conductor's immutability.

For Phase 8 I recommend the cleaner solution:

> **Human creates the operational task against the existing scenario, not a future belief ID.**

So:

```text
governance_ref:
{
  provider: "solvent",
  reference_id: "track1"
}
```

Then the packet-created belief/task relationship is represented through the scenario and existing EBP/Solvent state.

However, because the scenario can contain multiple beliefs, **do not silently claim that scenario membership uniquely identifies one belief**. For the G0 dry run, explicitly constrain the scenario to one governing claim or use a dedicated scenario for G0.

The implementation plan should state:

> Phase 8 dry run uses a one-governing-belief scenario; task governance references the scenario, not a future immutable belief ID.

That avoids modifying Conductor and avoids duplicate task creation.

---

### C3 — Human identity is currently underspecified

**Accept. Must fix.**

The decision request itself does not need to accept an arbitrary client-supplied actor. That would be weaker.

Use the existing frozen pattern:

```text
ARGUS_OPERATOR_PRINCIPAL_ID
```

The Trust UI server injects the configured operator identity; the Coordinator validates it and uses it for the authoritative Solvent discharge.

```text
Browser
  ↓
Trust UI
  ↓
Coordinator
  ↓
operator principal
  ↓
Solvent discharge
  ↓
audit_activity.actor_id
```

Fail closed when the operator principal is absent.

This makes:

> **human authority a server-side capability/configuration boundary, not a claim made by the browser.**

---

# H1 — "Adversarial Findings" is currently not derivable

**Accept, but simplify the UI metric.**

The current plan says to count evidence whose provenance came from adversarial role, but packet role is not persisted into the listed Solvent evidence model. 

Do not add a new role column.

For Phase 8, change the statistic from:

```text
Adversarial Findings
```

to:

```text
Adversarial Challenges
```

and derive it from:

```text
contradicts edges
```

This is already structural Solvent state.

The adversarial packet should therefore actually create:

```text
canonical:belief:G0
        ↑
   contradicts
        ↑
adversarial finding
```

rather than the current adversarial example having `edges: []`.

That gives us a clean invariant:

> Adversarial challenge = structurally recorded contradiction relationship.

More sophisticated adversarial finding attribution can come later.

---

# H2 — RCP must distinguish absence from failure

**Accept. Must fix.**

The sample `GetContext` implementation ignores nearly every client error:

```go
deps, _ := ...
evidence, _ = ...
```

That is dangerous.

A Solvent outage must never look like:

```text
beliefs: []
```

because the agent would interpret that as:

> “No work has happened.”

Use deterministic response shape **plus availability metadata**.

For example:

```json
{
  "epistemic": {
    "available": false,
    "reason": "solvent_unavailable",
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  }
}
```

Same shape, different epistemic meaning.

The key rule:

```text
UNKNOWN ≠ EMPTY
```

This is essential for the fresh-agent reconstruction claim.

---

# H3 — Exercise the retraction/dead-end path

**Accept, but keep it thin.**

The current acceptance path mostly exercises:

```text
enter → evidence → debt → refusal → retirement → promotion
```

while the plan also implements:

```text
contradicts
retract cascade
dead-end derivation
```

The current adversarial example does not create an edge, so that path would be dead code during the actual dry run. 

Add one explicit adversarial branch:

```text
Work
 ↓
Adversarial contradiction
 ↓
contradicts edge
 ↓
Human RETRACT
 ↓
Solvent cascade
 ↓
task becomes terminal/cancelled
 ↓
Insights derives Dead End
```

The retraction action does not need to become a third UI surface. Use the existing Coordinator decision path as the consequential mechanism and expose the resulting state in Insights.

That gives the dry run a real adversarial-to-dead-end demonstration without expanding the UI.

---

# H4 — Retirement-rule enforcement needs to be genuinely mechanical

**Accept. Must fix.**

There is currently a contradiction in the plan.

The workflow says:

> Coordinator checks offered evidence class against the Pack rule. 

But the implementation text says that a human click is sufficient and does not actually require the evidence class. That is exactly the mismatch the adversarial review identified.

Use the actual mechanical rule:

```text
Debt item
     ↓
Pack retirement_rule
     ↓
required evidence_class
     ↓
candidate evidence
     ↓
MATCH → allowed
MISMATCH → refusal
```

No scientific judgment occurs here.

That makes the Domain Pack's existing retirement rules meaningful rather than decorative. 

---

# H5 — Verify edge API and missing `ListIntents`

**Accept.**

The plan already lists `CreateEdge` in the Solvent client, but the REST API discovery table shown in the plan does not establish the endpoint. 

Before implementation:

```text
verify:
POST /v1/beliefs/{parent_id}/edges
```

If it exists, use it.

If it does not exist, **stop and invoke the existing Growth Gate**, rather than introducing direct SQL.

Also add the missing:

```text
ListIntents
```

to the Solvent client because RCP explicitly reads intents.

---

# H6 — Packet reference grammar must remain canonical

**Accept.**

Freeze the existing rule:

```text
local:<id>
```

for references inside the submitted packet.

```text
canonical:belief:<uuid>
```

for references to already-existing Solvent beliefs.

Do not allow the examples to drift between conventions.

The work and adversarial examples currently use both forms. 

---

# H7 — Idempotency must include scenario identity

**Accept.**

The current plan explicitly excludes `scenario_id` from the canonical idempotency hash. 

That is unsafe for epistemic isolation.

Two otherwise identical packets in distinct scenarios must not deduplicate.

Change the rule to:

```text
canonical hash =
canonical packet content
+ scenario identity
```

while continuing to exclude:

```text
packet_id
run_id
timestamps
runtime metadata
```

This preserves the Phase 6 idempotency invariant while respecting scenario boundaries.

---

# H8 — Conductor `cancelled` vs `rejected`

**Accept as reconciliation, not a new architecture change.**

The plan currently treats both as dead-end candidates, but the actual Conductor state machine makes `cancelled` terminal while `rejected` returns to active. 

Therefore Phase 8 should define:

```text
Dead End =
terminal cancelled task
+
governance-linked governing belief is retracted/contradicted
```

Do not use `rejected` as terminal state.

Add this to the existing freeze-reconciliation record.

---

# H9 — SQL in Insights is conceptual only

**Accept.**

The plan correctly prohibits direct DB access.

Change the wording from:

```sql
SELECT ...
ORDER BY ...
```

to:

> Client-side ordering equivalent to `critical > high > medium > low`, then `created_at`.

The source remains Conductor's existing `priority` field. No new schema.

---

# H10 — UI module violates the previously frozen two-module boundary

This is one additional issue I would add to the consolidated review.

The plan says:

```text
oracle/ui/
oracle/cmd/argus-ui/
```

but our frozen architecture has **Trust UI as a separate Go module**:

```text
oracle/go.mod
trust-ui/go.mod
```

and the UI communicates with Coordinator over HTTP.

So the implementation plan should not put the Trust UI inside `oracle/`.

Use:

```text
trust-ui/
  go.mod
  server.go
  handlers.go
  templates/
```

with:

```text
Trust UI → Coordinator HTTP
```

This preserves the existing architectural boundary rather than accidentally collapsing it during implementation.

---

# One more important provenance correction

The plan currently labels adversarial-agent evidence as `operator_asserted` in its remaining-decision section. That should **not** be accepted.

An autonomous agent's output is not automatically an operator assertion.

The clean Phase 8 model is:

```text
Agent analysis
    ↓
candidate evidence
    ↓
human adoption/attestation OR reproducible artifact
    ↓
eligible for authoritative debt retirement
```

Do not overload `operator_asserted` to mean "generated by an agent acting on behalf of an operator."

---

# Final consolidated baseline

After merging both reviews, I would freeze these corrections:

```text
C1  Human discharge uses attributed /v1/discharge.
C2  Task governance anchors to scenario for Phase 8; no mutable future belief_id.
C3  Operator identity injected server-side and fail-closed.

H1  "Adversarial Challenges" derives from contradicts edges.
H2  RCP distinguishes unavailable from empty.
H3  Dry run exercises contradiction → retraction → dead-end.
H4  Retirement evidence class is mechanically validated.
H5  Verify edge endpoint; add ListIntents.
H6  Freeze local:/canonical: reference grammar.
H7  Include scenario identity in idempotency hash.
H8  Dead-end uses actual terminal Conductor state: cancelled.
H9  Insights ordering is client-side derivation, not SQL.
H10 Trust UI remains separate trust-ui Go module.
H11 Agent evidence is not automatically operator_asserted.
```

The foundational verdict remains positive. The plan has already converged on the right shape: **two existing ledgers, one thin RCP projection, two MCP capabilities, two human UI surfaces, EBP-governed debt, and Solvent as the final authority**. The remaining issues are almost entirely about making the implementation match the promises precisely. 

**I would not broaden the architecture further. Patch these items into the plan, then implement Phase 8.**
