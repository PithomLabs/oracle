The attached review is largely **net-valid and consistent with the direction we've already established**. It adds one important governance step: don't declare the harness frozen merely because the P0 is fixed; prove the **write → read-back** loop and disposition the one known failing test first. The review correctly identifies that the latest regression is mechanical: persistence now works, while `GetContext` was broken by the read-side projection change. 

I would consolidate it this way.

## What I would retain

### 1. The freeze gate should be empirical

The strongest part of the attached review is the proposed freeze gate:

```text
persist
→ commit
→ GetContext
→ equality
```

rather than treating "rows exist" as proof that ARGUS works. The attached review correctly points out that the previous DB-backed tests proved persistence but did not prove the primary read path. 

That is now a **load-bearing acceptance criterion**.

The harness isn't ready to freeze until a fresh agent can:

```text
submit packet
    ↓
commit
    ↓
get_context
    ↓
see the same research state
```

### 2. Add the `TestAgentCannotPromote` disposition

I agree with this addition.

You don't necessarily need to "fix" the application-level behavior. The important thing is that the known failing test cannot remain ambiguous at freeze time. The attached review explicitly calls for either implementing the authority check or reclassifying/removing the stale test with a recorded reason. 

Given the current architecture, I would **not add another authority mechanism** just to satisfy a test. Instead, verify whether that test is asserting the wrong layer and then align the test with the actual invariant:

```text
agent-facing MCP
    → no promote/discharge/retract capability

human authority path
    → promotion/discharge/retraction
```

### 3. Keep the bookkeeping freeze

The attached review's retained/deferred distinction is consistent with our lock:

```text
load-bearing infrastructure → retain
unproven accounting improvements → defer
```

and specifically keeps coverage persistence, snapshots, epistemic-kind, attestation, graph visualization, corpus vector search, etc. deferred until concrete triggers occur. 

That is exactly the right direction.

### 4. Treat provenance as frozen, not endlessly enriched

The attached review's important conclusion is that the current provenance spine is already the minimum sufficient structure:

```text
packet
  ↓
belief / evidence / edge / task
  ↓
agent identity
```

The remaining missing UI projections are not a reason to invent another provenance system. 

I agree.

## What I would modify

### Don't delete the nil-DB MCP test automatically

The attached review proposes removing `TestMCPValidationRejectsMalformed` because it is superseded by the DB-backed test. 

I would **not make deletion mandatory**.

It can still be useful as a cheap unit test of the validation boundary. The problem is its name and interpretation, not necessarily its existence.

Better:

```text
TestMCPValidationRejectsMalformed
    = fast validation-path unit test

TestMCPRejectsMalformedWithRealDB
    = production-path integration test
```

The documentation/test naming should make it impossible to mistake the first for end-to-end coverage.

That preserves the useful cheap test without letting it provide false confidence.

### Do not immediately remove the `Agent` column either

I agree its current semantics are misleading, but the simplest fix is:

```text
Agent → Claimant
```

rather than dropping it.

That is a one-word semantic repair and doesn't add bookkeeping. The packet-origin information belongs elsewhere.

### Don't create a giant freeze manifest yet

The attached review proposes a fairly extensive freeze decision record. That's reasonable as a one-time artifact, but keep it **one page**. Don't turn the freeze itself into another accounting project.

---

# Revised implementation prompt

The previous read-path prompt should therefore be slightly expanded, not fundamentally changed:

```text
# ARGUS — Final Read-Path Repair + Freeze-Gate Verification

Fix ONLY the confirmed read-path regressions and perform the minimum
verification needed to determine whether ARGUS is ready to freeze.

DO NOT redesign ARGUS.
DO NOT add new MCP tools.
DO NOT add new accounting infrastructure.
DO NOT implement deferred features.
DO NOT change the authority model.

==================================================
1. FIX CONFIRMED READ-PATH BUGS
==================================================

A. `internal/epistemic/view.go:scanBelief`

GetSnapshot now selects:

    id
    claim
    claim_type
    status
    debt
    final_truth
    origin_packet_id

Update scanBelief() to scan all seven fields, including:

    OriginPacketID

B. `internal/work/store.go:ListByProject`

Ensure SELECT and Scan have exactly matching columns, including
origin_packet_id.

C. `internal/epistemic/view.go:GetAllBeliefs`

Include origin_packet_id in the SELECT so dashboard belief projections
preserve the provenance field.

==================================================
2. ADD THE MISSING END-TO-END READ-BACK TEST
==================================================

Add one DB-backed test that proves:

    persist valid packet
        ->
    commit
        ->
    GetContext(task_id)
        ->
    returned beliefs/evidence/edges match persisted state

Verify:

- belief count
- belief IDs
- claim text
- status
- debt
- origin_packet_id
- evidence
- edges
- availability metadata

The test must fail if:
- scanBelief becomes mismatched again
- GetSnapshot returns an empty snapshot because of a query/Scan error
- provenance disappears during projection

This is the freeze-critical invariant:

    WRITE -> COMMIT -> READ-BACK

==================================================
3. VERIFY THE TASK READ PATH
==================================================

Add or update a focused test for ListByProject proving:

- query succeeds
- origin_packet_id is returned correctly
- existing task fields remain correct

Do not introduce another abstraction.

==================================================
4. VERIFY THE DASHBOARD PROJECTION
==================================================

Verify GetAllBeliefs returns origin_packet_id.

Verify `/ui/insights` renders:

- beliefs
- origin information
- tasks

Do not add graph visualization.

Keep the existing task column but rename:

    Agent

to:

    Claimant

because `current_agent` means task claimant, not packet provenance.

Do not add another claimant/provenance subsystem.

==================================================
5. DISPOSITION THE KNOWN TestAgentCannotPromote FAILURE
==================================================

Inspect the failing test:

    TestAgentCannotPromote

Determine whether it is asserting an actual current architecture invariant.

Current intended boundary:

    MCP agent surface
        -> no promote/discharge/retract tools

    human authority surface
        -> promotion/discharge/retraction

Do NOT add a new authority layer merely to make the test pass.

Choose ONE:

A. The test is valid:
   implement the smallest missing authority check.

B. The test is asserting the wrong layer:
   modify/reclassify the test so it verifies the actual architectural
   boundary at the MCP/human surface.

Document the decision in the test or adjacent architecture note.

There must be no unexplained known-failing test at freeze time.

==================================================
6. KEEP THE BOOKKEEPING FREEZE
==================================================

Do NOT implement:

- review-coverage persistence
- context snapshot IDs
- epistemic-kind field
- agent attestation
- graph visualization
- corpus/vector search
- new MCP tools
- packet-ID replay detection
- additional provenance machinery

unless the current fixes uncover a concrete integrity failure that makes
one of these necessary.

The freeze rule is:

    concrete failure
        ->
    smallest fix

    no concrete failure
        ->
    defer

Research methodology and EBP scientific obligations are NOT frozen.

==================================================
7. REMOVE FALSE CONFIDENCE, NOT USEFUL TESTS
==================================================

Keep the existing nil-DB MCP validation test only if it is clearly scoped as
a fast validation/unit test.

Rename or document it so it cannot be mistaken for production integration
coverage.

Do NOT delete it merely for cosmetic simplification.

The DB-backed MCP tests remain the production-path proof.

==================================================
8. FULL VERIFICATION
==================================================

With CockroachDB running:

    go vet ./...
    go build ./...
    go test ./...

Then perform manual smoke verification:

1. Start:

    ./bin/argus serve

2. Submit one valid packet through the production MCP path.

3. Call:

    argus.get_context

4. Verify the new beliefs are visible.

5. Verify `/ui/insights`.

6. Verify `/ui/debts`.

7. Verify provenance is visible where currently supported.

==================================================
9. FREEZE GATE
==================================================

At the end, report:

### PASS / FAIL

Answer:

> Can a fresh agent submit research and then reconstruct exactly that
> committed research state through get_context?

### Freeze checklist

- [ ] packet validation on production MCP path
- [ ] packet persistence atomic
- [ ] provenance persisted
- [ ] valid packet reaches Commit
- [ ] malformed packet rejected before persistence
- [ ] exact retry idempotent
- [ ] GetContext read-back works
- [ ] Insights renders current state
- [ ] known TestAgentCannotPromote failure dispositioned
- [ ] full test suite green, or every exception explicitly dispositioned
- [ ] no new deferred accounting feature introduced

### Freeze recommendation

Return exactly one:

    READY TO FREEZE

or

    NOT READY TO FREEZE

If NOT READY, identify only the concrete blocker.

Do not propose new features.

The objective is to finish the current repair, verify the complete
write/read agent loop, and then freeze the harness so substantive research
can resume.
```

## My final read on the attached review

The review's main conclusion is sound: **you're very close to the freeze, but `GetContext` must be proven end-to-end before declaring victory**. The latest failures are overwhelmingly projection/wiring defects rather than architectural defects, which is exactly the situation where a freeze makes sense once the read-back invariant is green. 

The freeze gate should therefore be:

```text
VALIDATE
   ↓
PERSIST
   ↓
COMMIT
   ↓
GET_CONTEXT
   ↓
RECONSTRUCT
   ↓
HUMAN UI
```

Once that passes and `TestAgentCannotPromote` is properly dispositioned, **stop touching the harness unless a concrete research-cycle failure forces a change**. That is the cleanest application of "nothing left to subtract."
