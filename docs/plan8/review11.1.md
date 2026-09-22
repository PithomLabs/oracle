## prompt

what can you make of attached files (adv_review8.md, review11.md), my goal is to ultimately freeze the Argus harness so we can focus more on productive research work instead of more accounting


## answer

What the attached files make of — the pattern is more interesting than any single defect, and it says something precise about where the harness is in its maturation.

**The files show a project that has crossed from "can it work?" to "does it hold together at the seams?"** — Plan 14 repaired the provenance spine (FK ordering, DB-backed persistence tests: all five PASS, atomicity/idempotency/provenance preservation verified at the write layer), and in doing so introduced a textbook-class regression: `scanBelief` column-count mismatch breaking the *read* path, so packets now persist perfectly and then vanish from `get_context`. The system can write its own history but not read it back — which for ARGUS is P0 by definition, since the entire agent protocol is `persist → reconstruct`. The review caught it because the integration test *interpreted a scanner error as empty data* — the failure was present, tested, and misread. That last detail is the most valuable artifact in the file: a test that passes-for-the-wrong-reason is more dangerous than a failing one, and the review's demand (a DB-backed `Persist → Commit → GetContext → assert equality` test) closes exactly that class.

**On your actual goal — freezing the harness — these files are the final pre-freeze evidence, and they support the freeze with three structural observations:**

**1. The defect classes are converging on exhaustion.** Tally the post-architecture defects across the recent reviews: FK ordering (write-path composition), scanner mismatch ×2 (read-path projection), dropped field at the HTTP layer (F4), untested production path (F4's root cause), pre-existing unenforced authority test (§13, correctly ruled unrelated). Every one is a *wiring* defect between layers — no boundary violation, no authority leak, no doctrine violation. The MCP capability boundary passed again (§3, §19: two tools, zero mutation surface). When a system's failures are all projection-arithmetic rather than architecture, it's telling you the architecture is done and the remaining work is mechanical. That's the empirical signature you wait for before freezing.

**2. The test-quality lesson generalizes into the freeze criteria.** The core finding — "DB-backed tests proved rows exist but never called GetContext" — is the read-path dual of the earlier F4 (mock-path/production-path divergence). Together they give you the freeze-ready test doctrine: **for every frozen capability, at least one test must traverse the full production path — UI→Coordinator→persist→read-back for writes; MCP→RCP→projection for reads — asserting end-state, not existence.** Bake that into the freeze decision record, because it's the rule that catches the *next* regression class too, and it converts "the harness is frozen" from an aspiration into a checklist.

**3. The plan's review response is itself freeze-quality behavior.** Note what review11.md does: accepts the P0/P1/P2 verdicts, scopes the fix to exactly three files plus one test, explicitly fences the packet-ID replay question as a documented limitation *under the bookkeeping freeze* (correctly — TestMCP's ambiguity finding is real but unproven-necessary), and pre-declares what won't be touched. That's the disposition-log discipline operating live. The fix prompt is correctly scoped; approve it as written, with one addition: **require the `TestAgentCannotPromote` pre-existing failure to be dispositioned in the same pass** — not necessarily fixed, but *decided* (either implement the authority check the test demands, or delete/reclassify the test with a recorded reason). A known-failing test left ambiguous at freeze time becomes the first post-freeze dispute.

**On the freeze itself — what remains before you can declare it:**

```text
FREEZE GATE — Argus Harness v1.0
[ ] Plan 14 read-path fixes merged (scanBelief, ListByProject, GetAllBeliefs)
[ ] End-to-end GetContext regression test green (persist→read-back equality)
[ ] TestAgentCannotPromote dispositioned (implement check OR remove test, recorded)
[ ] Full suite green: go vet/build/test across all modules
[ ] Manual verification: fresh agent get_context sees persisted beliefs; 
    Insights/Debts render with provenance
[ ] MCP capability boundary test green (two tools, zero mutation)
[ ] Known-limitations register written (packet-ID replay, snapshot ids, 
    coverage persistence, epistemic-kind, attestation — each with trigger)
[ ] Freeze decision record: version, date, what's frozen, what's deferred-with-trigger
[ ] Annotated tag
```

Everything on that list is hours. And the honest reading of the two files is that you're closer than the FAIL verdict suggests: the write path is verified sound, the read path has three one-line fixes and a test, and the failure taxonomy has gone entirely mechanical.

**The strategic point, stated once: the freeze is justified precisely because the accounting is now load-bearing and complete enough.** The provenance spine (origin_packet_id everywhere), the capability boundary, RCP, the EBP gate map, the attributed discharge path — the harness now answers the four load-bearing questions (reconstruct state, determine trust, prevent unauthorized transition, reproduce history) with tested machinery. Every further refinement is deferrable-with-trigger under the bookkeeping freeze you already declared. What is *not* deferrable is the thing the harness exists for: the work→adversarial→human cycle on the v6.3 program — the Kronecker-spectrum test, the F_L construction, the tame-factor search. The research program's Phase 1 (§28) is sitting there fully specified, waiting for the harness to stop being the project.

So: approve the fix prompt as written (plus the TestAgentCannotPromote disposition), run the freeze gate checklist, tag it, write the record — and then the next version number anyone should care about is v6.4, produced by agents doing physics, not by reviews reviewing plumbing. The harness's remaining purpose after freeze is to be boring. Demand nothing more of it, and it will hold.



## prompt

what are the features that must be retained vs what must be deferred using our less is more design constraint, we know we reached the optimal design when there is nothing more to subtract for the absolutely essential features only before the freeze


## answer

Freeze-triage by the subtraction test. The method: a feature is **RETAINED** only if its removal breaks one of the four load-bearing capabilities (reconstruct state / verify trust / prevent unauthorized transition / reproduce history) or violates a frozen doctrine invariant. Everything else is **DEFERRED** — and the discipline is that deferral needs a trigger, not just a shrug. Here's the ledger.

## TIER 1 — ABSOLUTE CORE (nothing more to subtract; removal breaks the system)

| Feature | Why it survives the razor |
|---|---|
| **Solvent kernel as-is** (belief/evidence/debt, promotion gate, retraction cascade) | Frozen doctrine — the authority substrate. Not Argus's to modify. |
| **MCP two-tool surface** (`get_context`, `submit_packet`) | The entire authority model. Removing the boundary = agents mutate state. Removing the tools = no loop. |
| **Packet validation + Coordinator persistence** (single tx, Solvent-before-Conductor, idempotency on packet_id) | One of these gone = trust gone or state forked. |
| **RCP read projection + availability metadata** (`UNKNOWN ≠ EMPTY`) | The reconstruction contract. The scanner bug proved this is load-bearing, not convenience. |
| **Provenance spine** (`origin_packet_id` on belief/evidence/edge/task, `packet_submission` first) | Just demonstrated load-bearing: without it, "who challenged this?" is unanswerable. Minimum sufficient. |
| **EBP gate map** (human: promote/faithfulness/final-truth/add-debt; agents: candidates only) | Frozen doctrine, enforced by tool-surface + coordinator routing. |
| **Attributed discharge path** (`/v1/discharge`, operator principal, fail-closed) | Acceptance criterion 4; the human-authority record. |
| **Retirement-rule class check** (mechanical, coordinator-side) | The one enforcement that keeps pack rules from being decorative. |
| **End-to-end read-back test** (persist→commit→GetContext→equality) | The invariant that just failed; the freeze is contingent on it. |
| **Two-surface UI** (Insights, Debts) + dead-end derivation | The human control surface — the *minimum* adjudication view. Deriving dead-ends structurally is retained; nothing more. |
| **Contradicts/derives edges** | The lineage. Without edges, Insights is a list, not research. |

## TIER 2 — RETAINED-BUT-FROZEN-AS-IS (known limits, documented, not improved)

These stay in exactly their current imperfect form:

- **Declared agent identity** (harness/model as config) — attestation deferred, trigger: first provenance dispute an operator can't resolve.
- **Packet-ID replay = silent idempotent** — documented limitation, trigger: first concrete integrity incident from reuse.
- **Hardcoded scenario→pack mapping** — POC-scoped, trigger: second pack or multi-scenario contention.
- **Operator = configuration trust** — trigger: second operator or external exposure.
- **RCP scenario-projection scope** (no task-anchored closure) — trigger: scenario size makes context unusable.
- **Conductor lifecycle divergence** (cancelled/rejected semantics vs frozen doc) — recorded in §24; reconcile on next doc revision, not by code.

## TIER 3 — DEFERRED (explicitly not in the freeze; each with trigger)

| Feature | Trigger to revisit |
|---|---|
| Review-coverage persistence (per-belief dispositions) | Second cycle where a human can't tell "reviewed-clean" from "never-looked" |
| Context snapshot IDs | First cross-agent race where review-vs-submitted state matters to a decision |
| Epistemic-kind field (object/process/governance) | Second cross-stratum misclassification by an agent |
| Agent attestation (verified runtime identity) | First dispute unresolvable by declared identity |
| Priority derivation (mechanical) | Manual ordering demonstrably misdirects attention |
| Graph visualization / lineage UI beyond tables | Relational tables demonstrably fail a human adjudication |
| Corpus vector search / embeddings | An agent demonstrably fails a task for lack of semantic retrieval |
| Multi-operator auth / mTLS | External exposure or operator contention |
| Packet-ID replay detection | See Tier 2 |
| Any new MCP tool | A demonstrated loop failure traceable to tool absence |

## The subtraction verdict — what the razor actually cuts now

Applying "remove anything whose deletion breaks nothing load-bearing," **three things in the current codebase are candidates for *removal before* freeze, not deferral:**

1. **`TestMCPValidationRejectsMalformed` (nil-DB test)** — proves the validator is called, nothing else; the real-DB test supersedes it. Deleting it removes false confidence, costs nothing. (The one true subtraction available.)
2. **The `Agent = —` task column on Insights** — mislabeled semantics (claimant ≠ provenance) actively teach the wrong model. Relabel to `Claimant` or drop until origin display lands. Column deletion is subtraction.
3. **Any dead schema-path/wizard remnants from pre-Plan-14 state** — the sweep you'll run anyway; if the grep finds nothing, note that the system has genuinely reached nothing-left-to-subtract on this axis.

And the honest negative finding: **everything else currently in the harness survives the razor.** That's the mark that the freeze is ripe — the last three review rounds produced no feature *additions* to argue about, only wiring repairs. An architecture where the subtraction test yields only test-hygiene and one mislabeled column is an architecture at its floor.

## The freeze statement, in the doctrine's own grammar

> **Retained:** the minimum machinery that makes reconstruct, verify, gate, and reproduce true. **Frozen:** all of Tier 1 as-is, warts documented. **Deferred:** Tier 3, each with trigger. **Subtracted:** the nil-DB test and the mislabeled column. **After freeze:** the harness may receive bug fixes to Tier 1 and trigger-based Tier 3 items; it may receive no new capabilities. The next features this project ships are physics claims, debts, and gates — entered free, promoted under debt, adjudicated by humans.

One line of calibration to end on: the razor has now cut everything it can reach *in the harness*. The un-subtracted remainder — v6.3's 21 kill conditions, the substrate-class budget, the unrun Kronecker test — is where subtraction becomes research. Freeze the tool; the next thing that needs to shrink is the theory.