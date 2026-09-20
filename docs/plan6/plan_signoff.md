Agreed. The consolidation is sound, but the four dropped RCP dispositions and the new OpenCode tool-surface issue must be restored before freezing the baseline.

## Final Phase 8 baseline

### Critical: OpenCode tool surface

This becomes an explicit **agent-capability boundary**:

```text
WORK AGENT MCP SURFACE
  argus.get_context
  argus.submit_packet

ADVERSARIAL AGENT MCP SURFACE
  argus.get_context
  argus.submit_packet

NOT AVAILABLE
  solvent.retire_debt
  solvent.promote
  solvent.approve
  solvent.retract
  conductor.write/*
  direct database access
```

OpenCode agents therefore cannot mutate authority or operational state directly. Their only write path is an EBP packet submitted through ARGUS.

Phase 8 acceptance must include:

> Attempted direct Solvent mutation from a work/adversarial agent fails because the capability is absent from that agent's tool surface.

This is stronger than relying on prompt instructions.

---

# RCP v1 baseline

### 1. Endpoint

```http
GET /v1/context/{task_id}
```

Coordinator-owned, read-only.

MCP exposes this to OpenCode.

API remains the canonical integration contract.

### 2. Scope

RCP v1 returns the **full scenario projection**, not a sophisticated task-anchored graph traversal.

The response is assembled from the two existing ledgers:

```text
Conductor
  task
  dependencies
  activity

Solvent
  beliefs
  evidence
  edges
  debt
  intents
  audit history
```

The task-to-epistemic relationship is therefore intentionally broad in v1.1 rather than inventing a new inclusion algorithm.

### 3. Cross-ledger consistency

RCP is explicitly **eventually consistent** across Conductor and Solvent.

Every returned fact should identify its source:

```json
{
  "source": "conductor"
}
```

or:

```json
{
  "source": "solvent"
}
```

Projection gaps are therefore visible as lag, not silently synthesized state.

### 4. Artifacts

An evidence hash alone is insufficient for a fresh agent.

RCP must expose a resolvable artifact reference where available, for example:

```json
{
  "content_sha256": "...",
  "artifact_ref": "..."
}
```

The protocol does not duplicate artifact storage; it provides the existing resolution path.

---

# Dead ends

Use a structural definition, not free-form text.

For the Phase 8 POC:

> **Dead end = rejected Conductor task whose governance/decision record points to a retracted or contradicted belief.**

Do **not** use `cancelled` for the semantic dead-end counter.

Therefore Insights can derive:

```text
Dead Ends = count(tasks satisfying dead-end predicate)
```

No separate `dead_end` database field is required.

This is much better than trusting an agent-written `"status": "dead_end"` string.

---

# Debt retirement

I agree this needs one explicit decision before implementation.

For Phase 8, I would choose:

> **Coordinator mechanically validates that the evidence offered for debt retirement satisfies the Domain Pack retirement rule before invoking Solvent.**

So:

```text
Human
  ↓
select debt
  ↓
select/offers evidence
  ↓
Coordinator
  ↓
Pack retirement-rule validation
  ↓
Solvent RETIRE_DEBT
```

The Coordinator does not determine whether the research is scientifically correct. It only checks the mechanical contract already declared by the Pack.

That gives the UI's **PROPOSED DISCHARGE** language an actual enforcement path.

---

# Insights priority

Do not invent `HIGH / MEDIUM / LOW` as a new metadata field.

For the POC, define it as:

> **Manual task ordering in Conductor, rendered as ordered Kanban columns/cards.**

Thus the UI may display:

```text
ATTENTION

1. G0 formalization
2. Scaling audit
3. Canonical-height substrate
```

without pretending Conductor has a numeric priority model.

Later, a mechanical priority derivation can be added if there is a real need.

---

# UI

The sidebar should now be frozen as:

```text
ARGUS

Insights
Debts
```

### Insights

Derived from existing state:

```text
Program summary
Open claims
Open debt
Active work
Adversarial findings
Promoted claims
Refused transitions
Dead-end tasks
```

Then the ordered Kanban:

```text
ATTENTION
────────────────
G0
Scaling Audit
Canonical Height

WORKING
────────────────
...

BLOCKED
────────────────
...
```

The exact column labels can remain minimal; the important point is that they are **views over Conductor state**, not another workflow model.

### Debts

This is the human adjudication surface.

For each debt:

```text
Claim
Debt item
Why it exists
Evidence already available
Adversarial findings
Required evidence class/rule
Proposed discharge
Human decision
History
```

The action is:

```text
RETIRE DEBT
```

not:

```text
AI VERIFIED
```

The latter should never appear anywhere in the UI.

---

# Correct authority path

The architecture diagram should be corrected to make the write boundary unambiguous:

```text
                    HUMAN
                      │
                 Trust UI
                      │
                      ▼
                Coordinator
                 /        \
                ▼          ▼
          Conductor      Solvent
                ▲          ▲
                │          │
               RCP         │
                ▲          │
                │          │
          MCP / HTTP       │
                ▲          │
          OpenCode Agents ─┘
```

But the bottom arrow should be understood as **packet submission to Coordinator**, not direct Solvent access.

A cleaner representation is:

```text
OpenCode
   │
   ├── MCP: get_context
   │
   └── MCP: submit_packet
            │
            ▼
       Coordinator
        ┌────┴────┐
        ▼         ▼
   Conductor   Solvent
```

The UI follows:

```text
Human
  ↓
Trust UI
  ↓
Coordinator
  ↓
Solvent
```

No direct UI→Solvent mutation.

---

# Naming

Freeze one canonical protocol name:

**RCP — Research Context Protocol**

The envelope becomes:

```json
{
  "protocol": "RCP/v1"
}
```

rather than `argus-context/v1`.

ARGUS is the system implementing RCP; RCP is the protocol.

---

# Final frozen workflow

```text
1. HUMAN
   creates/selects Conductor task

2. WORK OPENCODE
   RCP → reconstruct current state
   performs work
   submits EBP packet

3. ARGUS
   validates packet
   updates Conductor/Solvent through Coordinator

4. ADVERSARIAL OPENCODE
   fresh process
   RCP → reconstruct same canonical state
   attacks unresolved claims
   submits adversarial EBP packet

5. ARGUS
   records findings/evidence/debt

6. HUMAN
   Insights → understands situation
   Debts → examines obligations
   discharges debt

7. COORDINATOR
   checks applicable Pack retirement rule

8. SOLVENT
   authoritatively retires debt

9. HUMAN
   requests promotion

10. SOLVENT
    decides whether promotion gate passes
```

This now preserves the complete chain:

```text
EBP
  ↓
RCP
  ↓
OpenCode
  ↓
EBP packet
  ↓
Coordinator
  ↓
Conductor + Solvent
  ↓
human debt discharge
  ↓
Solvent promotion gate
```

**Verdict: freeze this as the Phase 8 implementation baseline.** The remaining work is implementation, not another architectural expansion.
