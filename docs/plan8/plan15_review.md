## chatgpt

Plan 15 is **appropriately scoped and is the right final repair pass before the freeze gate**. It fixes the confirmed read-path defects, corrects the misleading UI label, disposes the stale authority test at the correct boundary, and—most importantly—adds the missing `WRITE → COMMIT → READ-BACK` test.  

I would approve it with **one small refinement**.

### The plan is correct

The three read-path repairs are exactly the right ones:

```text
scanBelief
ListByProject
GetAllBeliefs
```

The proposed `scanBelief` fix directly addresses the P0 that caused `GetContext` to become effectively empty after successful persistence. 

The `ListByProject` correction is also necessary because its scan contract is currently inconsistent. 

The `GetAllBeliefs` correction completes the provenance projection without inventing another subsystem. 

### The authority-test disposition is exactly right

Renaming:

```text
TestAgentCannotPromote
```

to:

```text
TestMCPDoesNotExposeAuthorityTools
```

is preferable to adding authorization semantics to `SubmitDecision` just to make an incorrectly scoped test pass. The plan correctly preserves the actual boundary:

```text
Agent
  → get_context
  → submit_packet

Human/control path
  → authority decisions
```



### The end-to-end test is the most important part

This is the test I would treat as the **freeze-critical invariant**:

```text
Persist
→ Commit
→ GetContext
→ equality
```

The plan explicitly checks beliefs, provenance, evidence, edges, and availability. 

That closes the exact class of defect we just encountered: a system can no longer declare victory merely because rows exist in SQL while the agent-facing read path is broken.

### One refinement

For `TestGetContextReadsBackPersistedState`, I would add one explicit assertion:

```text
snapshot.Task.ID == taskID
```

and verify that the task returned by `GetContext` is the same task whose scenario/project supplied the snapshot.

That protects against the other known weak point: the current design uses the task/project relationship to determine the scenario context.

So add:

```text
- snapshot.Task.ID == taskID
- snapshot.Task.ProjectID == expected project/scenario
```

No new feature—just make the freeze test stronger.

### I would also preserve the bookkeeping freeze exactly as written

The plan explicitly excludes:

```text
packet replay policy
context snapshots
review coverage
epistemic kind
attestation
graph UI
corpus/vector search
new MCP tools
additional authority layers
```

and keeps research methodology unfrozen. 

That is exactly the right ending.

## Disposition

**Approve Plan 15.**

Add the small `GetContext` task identity assertion, implement it, and then run the freeze gate:

```text
MCP validation
      ↓
packet persistence
      ↓
atomicity
      ↓
provenance
      ↓
Commit
      ↓
GetContext
      ↓
same state reconstructed
      ↓
Insights / Debts
```

If that is green, **freeze the ARGUS harness** and stop adding accounting capabilities unless a concrete Work/Adversarial cycle failure triggers one. The next engineering effort should then be the actual BM–IST–AS research.



## Z

Plan 15 review — the freeze-gate plan itself. Verdict: **approvable with two findings, one of which is a genuine architectural decision the plan has made implicitly and must make explicitly before the freeze record is written.** The rest is verification detail. Given that this is the plan that closes the harness, the standard is "nothing ambiguous survives into the freeze," and there are two ambiguities.

## Finding 1 — HIGH: §3's reclassification quietly resolves the earlier review's F3 differently than the record says, and the freeze record must name it

Trace the history: the code review (adv_review8) found `TestAgentCannotPromote` failing because `SubmitDecision` doesn't enforce caller identity, and ruled it "pre-existing, unrelated." The plan now **decides the test was wrong** — the authority boundary lives at the MCP tool surface, not the application layer — and replaces it with `TestMCPDoesNotExposeAuthorityTools`. That's a defensible call, and arguably the correct one: the boundary-by-configuration doctrine has been the project's stated position since v0.4, and adding auth to `SubmitDecision` to satisfy a misplaced test is exactly the scope-creep-under-freeze the plan cites.

But notice what got softer between documents. The earlier F3 finding (unauthenticated `/decisions` endpoint) was graded **Critical** — "any network client that can reach the Coordinator can invoke submitDecision directly," including the curl example. The plan's reclassification treats that as resolved by *tool-surface absence* — but tool-surface absence constrains **agents**, not network clients. A curl to `/decisions` needs no MCP tools; it needs only network reach. So the two findings rest on different threat models: F3 (network actor) and TestAgentCannotPromote (agent process). The plan resolves the *agent* half correctly. The *network* half — can something reach `/decisions` on the wire? — is answered nowhere in this plan except implicitly by "POC, localhost." That may be perfectly acceptable for the dry run, but it is precisely the kind of thing that must be a **written line in the freeze record**, not an ambient assumption:

```text
Freeze scope note: The Coordinator decision endpoint is authenticated by
network position (localhost/POC binding), not by request-level auth.
Agent authority is enforced by MCP tool-surface absence (tested).
Human authority is enforced by operator-principal injection (tested).
Network-level exposure of /decisions is a known limitation; trigger for
revisit: any deployment beyond loopback, or any second actor.
```

Without that paragraph, the freeze record silently downgrades a Critical finding to a footnote. With it, the decision is explicit, the earlier review is honored, and the trigger discipline continues. One paragraph; the difference between a freeze and a hope.

## Finding 2 — MEDIUM: the read-back test should also assert the *negative* that made this bug class dangerous

§4's test asserts the positive path thoroughly. Add two assertions that specifically pin the failure mode this plan is fixing:

- **Content equality, not just count**: assert `claim` text matches the packet's claim exactly (the scanner bug produced *empty* snapshots, but the same class of bug can produce *truncated* or *misordered* fields that count-based assertions pass).
- **A second, distinct packet**: persist two packets (two agents, e.g., one work, one adversarial) and assert each belief carries its *own* `OriginPacketID`. The provenance spine's actual invariant is per-object attribution across multiple origins — a single-packet test can pass with `origin_packet_id` accidentally uniform (e.g., a join bug copying one ID everywhere). Two origins, asserted distinct, pins it.

Cheap to add to the same test function; it's the difference between "the test would have caught the P0" and "the test would have caught the P0 *and* its nearest cousins."

## Verification notes

- **§1B column position**: the fix text says "after `governance_ref`, before `reopened_from_task_id`" and claims that matches `GetByID`/`ListAll`. Verify against the actual `GetByID` SELECT before merging — if the established order differs by one position, the same mismatch bug is reintroduced in the fix. Thirty-second check; worth stating in the plan as a verification step rather than an assumption.
- **§2's rename target**: confirm `insights.html:7` is the only occurrence — "Agent" as a column header may also appear in `debts.html` or the dashboard aggregate; the sweep should be template-wide, grep "Agent" in `templates/`.
- **§6**: good resolution of the subtraction-list item — renaming with a scoping comment rather than deleting preserves the cheap test while killing the false confidence. This is the bookkeeping freeze working as designed.
- **§10 checklist**: add two rows — "manual smoke performed per §7 and results recorded" and "known-limitations register written (per the freeze-gate requirements from the triage: replay, snapshots, coverage, kind, attestation, network exposure)". The freeze decision record is not complete without the limitations register; the checklist currently ends one artifact short.
- **§11's last line** ("Research methodology and EBP scientific obligations are NOT frozen") — correct and important; consider promoting it into the freeze record itself as the boundary sentence, since it's the clause that prevents the freeze from being misread as freezing the *research*.

## On the freeze recommendation format

The `[to be determined after verification]` slots are right — this plan correctly refuses to pre-declare PASS. One addition to the freeze record template, carried from the earlier triage: the record should include the **known-limitations register** (replay, snapshots, coverage, kind, attestation, network exposure — each with its trigger) and the **retained/deferred tier table** as of freeze. Those two artifacts are what make the freeze auditable in a year: a future contributor asking "why isn't X in Argus?" gets a citation, not an archaeology project.

## Verdict

**Approve. Execute §1–§7 as written, add the network-exposure freeze-scope note (F1), strengthen the read-back test with content-equality and dual-origin assertions (F2), extend the checklist with the limitations register, and verify the column-order claim in §1B against `GetByID` before merging.** Then the freeze recommendation can be filled in honestly — and given that this plan's defect list is three scanner mismatches and a mislabeled column, the harness is, as predicted, at its floor. Fill the slots, tag it, and the next adversarial review this project should run is of a physics claim.


## consolidated review

The review is largely **net-valid**, and it improves Plan 15 in the right way. I would incorporate its findings, but keep the additions tightly scoped so the freeze itself does not become another accounting project.

The combined conclusion is:

> **Plan 15 is the final mechanical repair pass before freeze. It should fix the read path, prove write→read reconstruction, explicitly record the authority threat-model boundary, and disposition the known failing test. After that, freeze.**

The attached Plan 15 already makes the critical repair and establishes the correct freeze invariant: `WRITE → COMMIT → READ-BACK`. 

### What from the review I would add

**1. Explicitly record the authority threat model.**

I agree with the review's distinction.

There are two different questions:

```text
Agent authority
    → what can an AI agent invoke?
    → MCP tool surface

Network authority
    → what can an arbitrary network client invoke?
    → HTTP/network exposure
```

Reclassifying `TestAgentCannotPromote` to the MCP boundary is correct. 

But the freeze record should explicitly state the scope of that guarantee:

```text
Agent authority is bounded by MCP tool exposure.

Human authority is bounded by the existing operator/control path.

This POC freeze assumes the Coordinator authority endpoint is not exposed
to untrusted network actors. Revisit before non-loopback or multi-actor
deployment.
```

That is a **scope statement**, not a new feature.

**2. Strengthen the read-back test.**

My previous recommendation was to assert task identity. The attached review adds two useful assertions:

```text
content equality
+
multiple packet origins
```

I agree with both.

The test should prove:

```text
Packet P1 → B1
Packet P2 → B2

GetContext()
  → B1.origin_packet_id == P1
  → B2.origin_packet_id == P2
```

This is better than a single-packet test, because it tests that provenance is actually per-object and not accidentally copied from one packet context.

Also assert exact claim text, not merely counts.

**3. Verify `ListByProject` column order against an existing known-good query.**

Correct. Since the current defect class is precisely column/scan mismatch, the plan should explicitly say to compare the exact SELECT/Scan ordering with `GetByID` and `ListAll` before implementing. Plan 15 already recognizes the mismatch; this makes the repair less error-prone. 

**4. Verify "Agent" label globally, but don't expand the UI.**

A quick template-wide grep is appropriate.

The actual fix remains one word:

```text
Agent → Claimant
```

No new UI.

**5. Add the limitations register to the freeze record, but keep it tiny.**

I agree that the final freeze should say explicitly what is deferred and what would trigger reconsideration.

But don't build a new subsystem or elaborate ledger. A small markdown section is enough:

```text
Deferred:
- packet-ID replay detection — trigger: concrete integrity incident
- context snapshots — trigger: concurrent/replay dispute
- review coverage persistence — trigger: human cannot distinguish reviewed vs skipped
- epistemic kind — trigger: repeated cross-stratum misclassification
- agent attestation — trigger: declared identity insufficient to resolve dispute
- network authority hardening — trigger: non-loopback/multi-actor deployment
```

That is consistent with the bookkeeping freeze you locked in.

### One thing I would *not* add

I would not make "full retained/deferred tier table" a second formal artifact if the current architecture documents already contain the same information. A **single freeze record with a short retained/deferred section** is enough.

That keeps the freeze from becoming its own bureaucracy.

---

# Consolidated revised prompt

Use this version for the coding agent:

```text
# ARGUS — Final Read-Path Repair + Freeze-Gate Verification

Fix the confirmed read-path regressions and perform the minimum verification
required to determine whether ARGUS is ready to freeze.

DO NOT redesign ARGUS.
DO NOT add new MCP tools.
DO NOT add new accounting infrastructure.
DO NOT implement deferred features.
DO NOT change the authority model.
DO NOT broaden scope.

==================================================
1. FIX CONFIRMED READ-PATH BUGS
==================================================

A. internal/epistemic/view.go:scanBelief

GetSnapshot SELECT returns:

    id
    claim
    claim_type
    status
    debt::STRING
    final_truth
    origin_packet_id

Update scanBelief() to scan all seven fields, including OriginPacketID.

B. internal/work/store.go:ListByProject

Before changing it, inspect GetByID and ListAll and verify their exact
SELECT/Scan column ordering.

Then make ListByProject SELECT and Scan match exactly, including
origin_packet_id.

C. internal/epistemic/view.go:GetAllBeliefs

Include origin_packet_id in the SELECT so BeliefView receives the persisted
provenance.

==================================================
2. UI SEMANTIC REPAIR
==================================================

Change the task table header:

    Agent

to:

    Claimant

because current_agent means task claimant, not packet provenance.

Before changing, grep the UI templates for other occurrences of the same
misleading task-column label.

Do not add new UI systems.

==================================================
3. DISPOSITION TestAgentCannotPromote
==================================================

Replace:

    TestAgentCannotPromote

with:

    TestMCPDoesNotExposeAuthorityTools

The test must verify the actual agent-facing MCP surface exposes exactly:

    argus.get_context
    argus.submit_packet

and does not expose:

    promote
    retract
    discharge

Do NOT add authentication to SubmitDecision merely to satisfy the old test.

Document the architectural decision briefly:

- agent authority is bounded by the MCP capability surface
- human/control authority exists outside that surface

==================================================
4. FREEZE-CRITICAL END-TO-END TEST
==================================================

Add:

    TestGetContextReadsBackPersistedState

This test must prove:

    WRITE -> COMMIT -> READ-BACK

Use a real DB.

Steps:

1. Create the required project.
2. Create a task for that project.
3. Persist a valid packet P1 containing beliefs, evidence, and edges.
4. Commit.
5. Call GetContext(taskID).
6. Verify:
   - snapshot task ID == taskID
   - task/project identity is correct
   - exact belief count
   - exact belief IDs
   - exact claim text
   - status
   - debt
   - origin_packet_id
   - evidence count and exact fields
   - edges and exact fields
   - availability metadata

Then strengthen the test with a SECOND packet P2 from a different agent.

Verify that:

    B_from_P1.origin_packet_id == P1
    B_from_P2.origin_packet_id == P2
    P1 != P2

Do not merely verify counts.

The test must prove that provenance is per-object and survives read-back.

==================================================
5. LISTBYPROJECT REGRESSION TEST
==================================================

Add:

    TestListByProjectReturnsOriginPacketID

Verify:

- query succeeds
- no Scan mismatch
- OriginPacketID is populated correctly
- existing task fields remain correct

==================================================
6. FAST MCP TEST CLARIFICATION
==================================================

Rename:

    TestMCPValidationRejectsMalformed

to:

    TestMCPValidationRejectsMalformed_FastUnit

Add a comment making clear:

    Fast validation-path unit test using nil DB.
    Not production integration coverage.

Keep it.

Do not delete a useful cheap test.

The DB-backed MCP tests remain the production-path proof.

==================================================
7. AUTHORITY THREAT-MODEL SCOPE NOTE
==================================================

Add a short note to the freeze decision record:

    Agent authority is bounded by MCP tool-surface exposure.

    Human/control authority uses the existing operator/control path.

    This POC freeze does not claim protection against an arbitrary network
    actor that can directly reach an internal authority endpoint.

    Revisit before non-loopback or multi-actor deployment.

Do not implement network hardening now unless an existing POC test proves it
is required for the current localhost threat model.

==================================================
8. FREEZE LIMITATIONS REGISTER
==================================================

Do NOT create a new database table.

Add a short section to the freeze record listing deferred items and their
triggers:

- packet-ID replay detection
  trigger: concrete integrity incident caused by reuse

- context snapshot IDs
  trigger: concurrent/replay dispute requiring exact reviewed state

- review-coverage persistence
  trigger: human cannot distinguish reviewed-clean from never-reviewed

- epistemic-kind field
  trigger: repeated cross-stratum misclassification

- agent attestation
  trigger: declared identity cannot resolve a provenance dispute

- network authority hardening
  trigger: deployment beyond loopback or multiple untrusted actors

These remain DEFERRED, not implemented.

==================================================
9. BOOKKEEPING FREEZE

The current freeze rule remains:

    concrete failure
        ->
    smallest fix

    no concrete failure
        ->
    defer

Do not add:

- context snapshots
- review-coverage persistence
- epistemic-kind fields
- agent attestation
- graph visualization
- corpus/vector search
- new MCP tools
- additional provenance infrastructure

Research methodology and EBP scientific obligations remain unfrozen.

==================================================
10. VERIFICATION
==================================================

Run:

    go vet ./...
    go build ./...
    go test ./...

With CockroachDB running, also perform the real DB-backed suite.

Then manually verify:

1. ./bin/argus serve
2. Submit one valid packet through production MCP.
3. Call argus.get_context.
4. Confirm newly persisted beliefs are returned exactly.
5. Verify /ui/insights.
6. Verify /ui/debts.

Record actual results.

==================================================
11. FREEZE GATE
==================================================

The freeze gate is:

- packet validation works on production MCP path
- packet persistence is atomic
- provenance persists
- valid packet reaches Commit
- malformed packet rejected before persistence
- exact retry is idempotent
- GetContext reconstructs the committed state exactly
- Insights renders current state
- TestAgentCannotPromote is properly dispositioned/reclassified
- full suite is green OR every remaining failure is explicitly dispositioned
- limitations register exists
- no new accounting capability was introduced

Final question:

> Can a fresh agent submit research and then reconstruct exactly that
> committed research state through get_context?

Return:

    READY TO FREEZE

or:

    NOT READY TO FREEZE

If NOT READY, identify only the concrete blocker.

Do not propose new features.

The purpose of this pass is to finish the mechanical repair, prove the
complete agent loop, and then freeze ARGUS so the project returns to
substantive BM–IST–AS research.
```

The key refinement is that **the freeze is now a proven property of the full loop**, not merely a statement that the database and MCP layers individually work:

```text
validate
→ persist
→ commit
→ reconstruct
→ inspect
```

Once that chain passes and the authority/network scope is explicitly documented, I agree with the attached review's broader conclusion: **freeze the harness.** The next adversarial review should be directed at the physics, not another round of ARGUS plumbing.
