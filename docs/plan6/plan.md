Yes. With the EBP v2.1 constraint added, I would make one final architectural consolidation before Phase 8:

> **RCP is a thin agent-facing read protocol over existing Conductor + Solvent state; EBP is the governing epistemic protocol; Solvent remains the enforcement point for debt and promotion; the UI is primarily a human debt-discharge and situational-awareness surface.**

EBP v2.1 is explicit that ideas may enter with essentially no prerequisites, but promotion is blocked by unpaid debt, and humans remain responsible for discharging that debt.  

## 1. How OpenCode should query RCP

I would use **API as the canonical protocol and MCP as the OpenCode-facing adapter**.

That gives us the least duplication:

```text
OpenCode
   │
   │ MCP tool
   ▼
ARGUS MCP adapter
   │
   │ HTTP
   ▼
Coordinator
   │
   ├── Conductor API
   └── Solvent API
```

So RCP itself is **not another server with another database**.

The canonical operation is conceptually:

```http
GET /v1/context/{task_id}
```

The MCP surface simply exposes that operation to OpenCode as something like:

```text
argus.get_context(task_id)
```

This is the right split because:

**API**

* canonical integration contract
* usable by UI, scripts, other agents, CI
* easy to test independently

**MCP**

* ergonomic agent interface
* lets OpenCode discover/use the capability as a tool
* no need to teach the agent your HTTP mechanics

And importantly, RCP should call **existing Conductor/Solvent interfaces**, not their databases directly.

```text
RCP
  → Conductor API
  → Solvent API
```

Only extend Conductor or Solvent when an existing interface genuinely cannot provide something required by the dry run.

### What RCP returns

Keep it extremely small:

```json
{
  "protocol": "argus-context/v1",
  "task": {},
  "dependencies": [],
  "epistemic": {
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  },
  "activity": []
}
```

Everything here is **projection**, not new state.

This directly addresses the six implementation plans' recurring concern that the research program must become queryable and that prior failures and dependencies must remain visible rather than being recreated from memory. The Claude plan in particular emphasizes turning the 47-section program into a queryable ledger and retaining failed work rather than rediscovering it.  

---

# 2. OpenCode workflow, finally simplified

### Work Agent

```text
Human
  ↓
create/select Conductor task
  ↓
OpenCode
  ↓
argus.get_context(task)
  ↓
do research
  ↓
produce EBP packet
  ↓
ARGUS
  ↓
Solvent + Conductor
```

### Adversarial Agent

Fresh OpenCode process, no conversational inheritance:

```text
OpenCode
  ↓
argus.get_context(task)
  ↓
sees:
  task
  dependencies
  completed work
  evidence
  current debt
  prior challenges
  dead ends
  activity
  ↓
attack unresolved work
  ↓
produce adversarial EBP packet
  ↓
ARGUS
  ↓
Solvent
```

The crucial property is:

> **The adversarial agent does not reconstruct history from files or conversation. It reconstructs the authoritative current state through RCP.**

That is the interface boundary you were after.

---

# 3. EBP changes the meaning of the workflow

We should make the loop explicitly:

```text
IDEA
 ↓
ENTER FREE
 ↓
WORK
 ↓
EVIDENCE
 ↓
ADVERSARIAL ATTACK
 ↓
DEBT
 ↓
HUMAN DISCHARGE
 ↓
PROMOTION
```

Not:

```text
agent does work → system decides it's good
```

The agent can **produce material that could retire debt**, but it does not itself discharge the debt.

That follows EBP's distinction between entry and promotion, and its statement that debt blocks promotion rather than killing the idea. 

---

# 4. The UI should now be only two primary surfaces

I agree with your sidebar concept. I would make it deliberately small:

```text
ARGUS
────────────────
INSIGHTS
DEBTS
────────────────
```

No mythology-heavy navigation for the first dry run. The mythology can remain in the visual language, but the actual POC should make the operational workflow obvious.

## A. Insights

This is the **research situation room**.

It answers:

> What have we learned, what is happening, and where should human attention go?

Top-level summary:

```text
BM–IST Synthesis v5.1

Ideas / Claims        12
Active Tasks           5
Open Debt              9
Adversarial Findings   3
Promoted               2
Refused                4
Dead Ends              2
```

Then the **attention Kanban**:

```text
HIGH PRIORITY
─────────────
G0 — Formalize proof
Gate 2 — Scaling audit

MEDIUM
──────
Canonical-height substrate

LOW
────
Deferred cosmology branch
```

But one important correction from the attached compliance review:

### Priority belongs to Conductor, not Solvent.

So the Insights board should derive task priority from **Conductor operational state**.

Do not create a Solvent `priority` field.

Likewise:

> **Attention state is workflow state, not epistemic state.**

That resolves the attached H2 concern cleanly.

---

# 5. Debts

This is the more important screen.

It should answer:

> **What does this idea still owe before it can be promoted?**

Example:

```text
G0 — Discrete Substrate / Continuous Unitary Compatibility

STATUS
Unpromoted

OPEN DEBT
────────────────────────
✓ needMap
✓ needInvariant
○ needToyCheck
○ needNullModel
○ needObstruction
○ needFaithfulnessReview
```

Click `needObstruction`:

```text
WHY THIS DEBT EXISTS

Adversarial finding:
The current formalization has not yet demonstrated
that the assumptions of the cited theorem apply to
the actual substrate object.

EVIDENCE
A-017
E-031
E-044

PROPOSED DISCHARGE
Formalize the missing assumption and provide a
counterexample attempt.

SOURCE
OpenCode / adversarial run #2
```

Then the human decides:

```text
[ RETIRE DEBT ]
[ KEEP OPEN ]
```

The actual retirement invokes the existing Solvent pathway.

That is exactly where your previous hackathon UI work is useful: **reuse the proven human debt-discharge interaction rather than inventing another mechanism.**

---

# 6. The human is the debt authority

This needs to be extremely explicit in Phase 8.

```text
Agent
  ↓
proposes evidence / argument
  ↓
Human
  ↓
decides whether debt is actually discharged
  ↓
Solvent
  ↓
records authoritative transition
```

So the UI should never show:

> “AI verified debt.”

It should show:

> “Adversarial agent supplied evidence relevant to `needObstruction`.”

Then:

> **Human decision: discharge debt**

That is much closer to EBP v2.1's intended philosophy. EBP describes debt retirement as a useful move that the researcher performs and says promotion only becomes available once current debt is retired. 

---

# 7. Consolidating the attached compliance findings

I would keep the useful findings, but simplify their resolution.

### C1 — Partial debt retirement

**Do not add partial-debt semantics.**

The current Solvent representation is binary:

```text
debt present
debt absent
```

Keep it that way for the POC.

Record explicitly:

> EBP v2.1 describes partial progress informally, but ARGUS v1 implements debt-item retirement as atomic. Partial progress remains evidence/prose until the entire debt item is discharged.

That is actually consistent with the broader "less is more" rule.

### C2 — System-inferred debt retirement

**Reject automatic retirement inference.**

This is a good deliberate strengthening.

The attached EBP text does describe seamless inference, but our authority architecture requires that semantic adjudication not be silently performed by an LLM at the Solvent boundary. EBP itself says accounting must not become the work. 

So:

```text
Agent suggests
Human adjudicates
Solvent records
```

No automatic NLP-based retirement.

### H1 — REOPEN

Keep the existing architecture:

```text
retracted belief
   ↓
new belief
   ↓
derives edge
```

but make the new belief **inherit relevant lineage/context**, not blindly become a new free-standing research idea.

That preserves the EBP living-ledger concept that an idea can accrue new debt after new evidence. 

### H2 — Attention

Resolve through Conductor.

No new Solvent field.

### H3 — “Next smallest useful move”

Keep it **UI/agent guidance**, never Solvent state.

That keeps the kernel domain-neutral.

---

# 8. One final refinement: RCP should expose debt, not prescribe how to discharge it

This is important.

RCP can say:

```text
OPEN DEBT:
needObstruction
```

It can expose:

```text
supporting evidence
prior attempts
adversarial findings
related tasks
dead ends
```

But it should not say:

```text
NEXT MOVE = run toy check X
```

That belongs to the **BM–IST Domain Pack / agent reasoning / human workflow**.

Otherwise RCP quietly becomes an epistemic advisor.

The attached compliance assessment was correct to identify this as a boundary issue.

---

# Final Phase 8 architecture

I would now freeze Phase 8 around this:

```text
                    HUMAN
                      │
          ┌───────────┴───────────┐
          ▼                       ▼
      INSIGHTS                   DEBTS
          │                       │
          │                       ▼
          │                Human discharges
          │                       │
          └──────────┬────────────┘
                     ▼
                  SOLVENT
                     ▲
                     │
                COORDINATOR
               /             \
              ▼               ▼
        CONDUCTOR           SOLVENT
              ▲
              │
             RCP
              ▲
              │
       MCP / HTTP adapter
              ▲
              │
          OpenCode
        /            \
      WORK         ADVERSARY
```

And the protocol boundaries are:

```text
EBP
= epistemic operating law

Conductor
= operational work state

Solvent
= authoritative epistemic/debt/authorization state

RCP
= thin read projection

MCP
= agent-facing transport

API
= canonical integration interface

Trust UI
= human observation + adjudication
```

### Phase 8 dry-run acceptance

The first dry run should prove only this:

1. Human creates a BM–IST task in Conductor.
2. Work OpenCode retrieves context through MCP → RCP.
3. Work agent produces an EBP packet.
4. ARGUS records the resulting claims/evidence/debt.
5. Fresh adversarial OpenCode retrieves the same current state through RCP.
6. It attacks the unresolved work without rediscovering completed/dead-end work.
7. **Insights** shows the evolving research state and Conductor priorities.
8. **Debts** shows why each obligation remains open.
9. Human discharges at least one debt through the UI.
10. Solvent records the retirement.
11. Promotion is attempted and Solvent—not the UI or agent—decides whether the gate is satisfied.

That is thin enough to implement now, while exercising the actual architecture rather than building a decorative UI. EBP v2.1 itself calls for a minimum implementation centered on capture, status, retire, add-debt, and promote, while explicitly warning against dashboards and bureaucracy becoming the work. 

**I would consider this the final Phase 8 baseline.**
