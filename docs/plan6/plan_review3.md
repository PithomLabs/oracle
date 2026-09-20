I reviewed the revised plan. It is substantially stronger and the edge Growth Gate is now correctly incorporated: the plan adds only `POST /v1/beliefs/{id}/edges`, keeps it out of MCP, and routes packet edge creation through Coordinator → Solvent. 

I would **not implement yet**. Two critical workflow contradictions remain, plus two smaller fixes.

## C1 — The dry-run sequence is logically impossible as written

The plan says:

```text
adversarial contradicts G0
→ human retracts G0
→ G0 becomes dead end
→ human discharges remaining debt
→ human requests promotion
```

But once G0 is retracted, **G0 cannot then be promoted**. The plan itself defines retraction as a Solvent lifecycle transition and says the original remains retracted during REOPEN. 

This is the biggest remaining problem.

### Fix

Split the dry run into two explicit branches:

```text
BRANCH A — Adversarial / Dead-End

G0 entered
→ adversarial contradiction
→ human RETRACT
→ RetractCascade
→ cancelled task
→ Dead End visible
```

Then separately:

```text
BRANCH B — Surviving / Promotion

L1 or another surviving candidate
→ evidence
→ human discharges debt
→ promotion attempted
→ Solvent permits/refuses
```

Or use:

```text
G0
→ contradiction
→ RETRACT
→ REOPEN
→ successor belief via derives
→ discharge successor debt
→ promote successor
```

I recommend **Branch A + Branch B** for the first dry run. It is simpler and proves both behaviors without making REOPEN part of the critical path.

The acceptance criteria should no longer say that the same retracted belief is later promoted.

---

## C2 — Where does Conductor cancellation actually happen?

The plan says:

> `Solvent RetractCascade → Conductor task becomes terminal cancelled` 

But the listed Conductor REST API does **not show a cancellation endpoint**. The plan explicitly shows `cancelled` as a terminal state, but the exposed operations listed are claim, submit, accept, reject, blocker, resolve, and release. 

That means the dry-run path currently depends on a transition mechanism that the plan has not established.

This needs an explicit decision before coding:

```text
A. Existing Conductor cancellation interface actually exists
   → use it and document it.

B. It does not exist
   → add the minimal Conductor cancellation endpoint
     as another explicit Growth Gate exception.

C. Do not exercise automatic task cancellation in Phase 8
   → derive dead-end only from an already-cancelled task,
     and remove "retraction causes cancellation" from the claim.
```

Given the goal of exercising the full loop, **B is probably the cleanest if the endpoint genuinely does not exist**. But the coding agent must verify the actual repository first rather than assume.

Do not mutate the Conductor database directly.

---

## H1 — Operator-asserted debt needs one precise rule

The plan now correctly refuses to treat agent output as automatically `operator_asserted`. 

But the retirement flow still says:

```text
human selects debt
→ identify evidence
→ compare evidence class
→ discharge
```

For BM–IST:

```text
needNullModel           → operator_asserted
needFaithfulnessReview  → operator_asserted
```

Those may be satisfied by **human attestation itself**, not necessarily an evidence row.

Freeze this:

```text
reproducible_artifact debt
  → requires matching evidence object

operator_asserted debt
  → requires explicit human attestation
  → /v1/discharge records operator identity
```

That keeps the system minimal and makes the existing `DischargeRequest` actually useful. The plan already has the attributed discharge API and operator identity.  

---

## H2 — `available` should be consistent across the entire RCP envelope

The plan correctly fixes `UNKNOWN ≠ EMPTY`. 

But the response examples only put `available` under `epistemic`, while the text says:

> each section includes `available` boolean. 

Make that explicit:

```text
task:
  available
  data

dependencies:
  available
  data

epistemic:
  available
  data

activity:
  available
  data
```

Same deterministic envelope every time.

---

## Everything else is now in good shape

The edge endpoint is correctly scoped as a Growth Gate interface extension, not a schema change. 

The MCP boundary is clean:

```text
argus.get_context
argus.submit_packet
```

with no Solvent or Conductor mutation tools. 

The plan also correctly preserves:

```text
local:<id>
canonical:belief:<uuid>
```

as the packet reference grammar. 

Idempotency now includes scenario identity, which closes the cross-scenario deduplication problem. 

The Trust UI remains separate, minimal, and human-centered. 

And the test matrix is now appropriately broad, including edge creation, persistence, RCP degradation, human identity, adversarial reconstruction, retraction, dead-end derivation, UI, concurrency, and capability boundaries. 

### Final verdict

**Almost ready.**

Before implementation, patch:

```text
C1  Split retraction/dead-end and promotion into separate dry-run branches.
C2  Verify/add the actual Conductor cancellation path.
H1  Define operator_asserted debt as explicit human attestation.
H2  Make RCP availability metadata consistent per section.
```

After those four changes, I would stop architecture review and move to implementation.


## Z

Final-cut review, v2. Disposition first — this revision is an unusually faithful response to review pressure, and the record should show it: **every Critical and High finding from the previous round landed as mechanism, not prose.**

| Prior finding | Status |
|---|---|
| C1 unattributed discharge | ✓ Fixed — attributed `/v1/discharge` with `DischargedBy`, operator identity from server-side config, §24 reconciliation of the two-path split, bare-retire demoted |
| C2 governance_ref linkage | ✓ Fixed — scenario-anchored (`reference_id` = scenario), metadata dead-end derivation removed, §24 reconciles |
| C3 human identity | ✓ Fixed — §10, server-injected principal, fail-closed, browser-cannot-override, triple test |
| H1 adversarial findings not derivable | ✓ Fixed — contradicts-edge count (§13), adversarial packet now creates the actual edge (§9) |
| H2 error-swallowing | ✓ Fixed — `available`/`reason` metadata, UNKNOWN ≠ EMPTY promoted to named principle (§15, §16) |
| H3 retraction unexercised | ✓ Fixed — dry-run branch in §1.8/§12, dead-end machinery exercised (acceptance 8) |
| H4 rule-enforcement mismatch | ✓ Fixed — class matching enforced in Coordinator (§6), criteria 5 updated to match, matching tests added |
| M-tier | ✓ Landed — reference grammar frozen (§8), idempotency includes scenario (§18), cancelled reconciliation (§24), client-side ordering |

And the section this series has waited for: **§24 Freeze/Reconciliation Notes.** The plan now contains its own reconciliation document — state machine divergence recorded (with explicit supersession language), the edge endpoint Growth Gate exception named as such, the discharge-path decision, the UI module location, and — best of all — the agent-evidence provenance reconciliation ("Do not overload `operator_asserted` to mean 'agent acted on behalf of operator'"), which closes a hole no reviewer had explicitly flagged. The request count on the standalone reconciliation record drops accordingly: this §24 *is* it, for this plan's scope; the cross-plan version (Workflow Design v1.0, the research POC freeze, two rounds of supersessions) remains owed but is now a one-afternoon consolidation rather than an archaeology project.

## Residual findings — none architecture-class, a few worth fixing before or during implementation

**H1 — §6's code comment contradicts §6's own spec.** The pseudocode comment says "fail-safe: unavailable ≠ empty" while the code as sketched still assigns `beliefs, belErr := ...` with no visible handling, and the assembly step is a comment rather than the availability-stamping logic the tests require. The tests are right; the sketch should match them or a reader will implement the sketch. Also: `GetDebt` is called as a distinct client method in §6 but was never added to the §4 client-extension list, and the "new evidence creates new debt" flow (§11) names no mechanism — new *evidence* doesn't itself add debt in Solvent; a contradicts edge or human add-debt does. §11's sentence promises a behavior no component implements. Either specify the trigger (adversarial contradicts edge → coordinator offers `add_debt` decision to human) or reword to match EBP's actual operation (new evidence *justifies* new debt via human/pack action).

**H2 — The operator principal is configuration-trust, which is adequate for POC and should be labeled as such.** §10 solves C3 correctly for a single-operator dry run, but note what's deferred: any party who can edit the UI server's config *is* the operator — there's no authentication of humans, only of configuration. That's the right POC trade; it needs one sentence in §24 ("operator identity is configuration-trusted; multi-operator authentication is post-POC") so the record doesn't overclaim, and it closes the standing identity finding honestly rather than silently.

**H3 — Retraction cascades to "dependent Conductor tasks become terminal cancelled" — but no mechanism connects Solvent's cascade to Conductor's task state.** §12 asserts the transition; nothing in §4 or §6 implements it (no coordinator watcher, no explicit step in the retraction decision handler). Either the retraction decision handler explicitly cancels linked tasks after Solvent confirms (one more client call, scenario-anchored linkage per §24's reconciliation), or the dead-end derivation must handle "task still active, belief retracted" — which §12's own dead-end definition *excludes* ("active + retracted belief: task not terminal"). As written, the dry run's step 9 (task becomes cancelled) happens by assertion, not by code. This is the one remaining seam where the demo will visibly stall.

**H4 — Scenario-extraction ambiguity on the human-created task.** §8's bootstrap has the human creating the task; §5/§6 derive the scenario from `governance_ref.reference_id`. Fine — but the plan should state that the human task's governance_ref MUST be pre-populated with the scenario (`track-g0`) at creation, since write-once-ness means there's no second chance, and the work packet's own task entries (which duplicate the human task per prior round's note) need a stated relationship or the run produces two tasks for one objective. One paragraph: human task = umbrella (scenario-anchored), packet tasks = sub-work, both scenario-linked.

**M-tier:** cross-scenario idempotency test exists (good) but the dedup store itself needs scenario in its key, not just the hash — same content, two scenarios, two cache entries; the §14 artifact base directory is referenced but never defined (one path constant); `available` metadata on the *task* section is described in §16's error table but absent from the §5 response example — align; "deterministic response" test should pin ordering of activity arrays (merged from two stores — sort by what? specify per-source ordering with stable tiebreak, per the happens-before doctrine).

## Verdict

**Approved for implementation.** The remaining findings are implementation-guidance items (H1's sketch/handler alignment, H3's retraction-to-task step, H4's scenario-prepopulation note) — each resolvable in the code as written or by one paragraph. Nothing requires another plan revision; requiring one would be the paralysis pattern this project buried.

The closing observation for the record: **this plan is the review lineage compiled.** Capability boundaries as tool-surface configuration with negative tests; attribution flowing into the authoritative ledger; refusals demonstrated deliberately; dead ends derived structurally; UNKNOWN ≠ EMPTY as a named principle; enforcement decisions reconciled against the freeze in-document. Every mechanism in it exists because some earlier version of this project failed that check. Phase 8's verdict will now be delivered by the only reviewer left — the run itself — and for the first time in the program's history, everything it needs to judge has been built in front of witnesses. Build it.

## chatgpt

Agreed. This review lands cleanly, and I would **close Phase 8 planning and move to implementation**.

The revised plan now captures the substantive architecture correctly: RCP over Conductor + Solvent, agent capability restricted to `argus.get_context` / `argus.submit_packet`, human-attributed discharge through Solvent, structural adversarial challenges, explicit retraction/dead-end exercise, and Solvent as the promotion authority.  

### Implementation-time invariants

The remaining findings should be treated as **implementation instructions**, not another planning round:

**1. RCP assembly must implement the availability contract.**
Do not copy the pseudocode's discarded-error pattern. Every section must preserve `available=false` plus a reason rather than silently returning empty state. The plan already establishes `UNKNOWN ≠ EMPTY`.  

**2. `GetDebt` must either be implemented or derived from the existing belief debt array.**
There is no reason to create unnecessary Solvent surface area. Since debt already lives on beliefs, prefer deriving the RCP debt projection from `belief.debt[]` unless an existing `GetDebt` API is actually needed.

**3. New evidence does not magically mutate debt.**
The implementation should follow the EBP doctrine correctly: new evidence can justify new debt; the actual debt transition must occur through the existing human/Coordinator/Solvent mechanism. Do not add an implicit agent-driven `add_debt` side effect merely to satisfy prose.

**4. Operator identity is configuration-trusted for the POC.**
Record explicitly that this is the POC boundary: configuration supplies the operator principal; multi-user authentication is deferred.

**5. Retraction must explicitly cancel the linked Conductor task.**
This is the one seam I would treat as mandatory rather than optional. The retraction decision handler should:

```text
Solvent RetractCascade succeeds
        ↓
Coordinator identifies scenario-linked task(s)
        ↓
Coordinator invokes Conductor cancellation
        ↓
Insights derives Dead End
```

Otherwise the planned dead-end demonstration cannot actually occur.

**6. Human task creation must pre-populate the scenario governance reference.**
Because `governance_ref` is immutable, the human-created umbrella task must be created with the correct scenario from the beginning. Packet-created tasks should be explicitly treated as subordinate work items, not accidental duplicates.

**7. Idempotency key must contain scenario identity.**
The implementation should ensure both the canonical hash **and the cache key** are scenario-scoped.

**8. Activity ordering must be deterministic.**
When merging Conductor and Solvent activity, define a stable ordering with timestamp plus deterministic tiebreaker. Otherwise `RCP/v1` is not truly reproducible.

### Final status

```text
Phase 8 planning: CLOSED
Phase 8 implementation: APPROVED
Architecture expansion: STOP
Next artifact: actual implementation + dry run
```

The plan has reached the point where further review would mostly be substituting hypothetical defects for observed ones. Its acceptance criteria now test the actual thesis of ARGUS: fresh-agent reconstruction, adversarial challenge, human debt discharge, structural dead ends, capability separation, and Solvent-controlled promotion. 

**Build it.**
