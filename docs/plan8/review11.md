This review found another **real P0**, and it is importantly different from the previous one: **the write path is now working, but the primary read path is broken.**

The review verified that Plan 14 successfully fixed persistence: valid packets commit, rollback works, exact retries are idempotent, and provenance survives.  

But adding `origin_packet_id` to the belief query without updating `scanBelief` created a seven-column/seven-destination mismatch. `GetContext` therefore cannot reconstruct the research state and ends up exposing an empty snapshot. 

That is **P0 for ARGUS**, because your whole agent protocol depends on:

```text
submit_packet
    ↓
persistent state
    ↓
get_context
    ↓
next agent reconstructs state
```

The system currently breaks at the third step.

### The three fixes I would make

**1. Fix `scanBelief` — mandatory P0.**

Add `origin_packet_id` to the scan destinations exactly as the review specifies. 

**2. Fix `ListByProject` — P1.**

The query and scan have diverged by one column. This is another symptom of adding a field to the model without updating every projection. 

**3. Fix `GetAllBeliefs` — P2, but do it in the same pass.**

Otherwise the dashboard still cannot show belief provenance despite the field existing. 

And crucially, add the missing end-to-end test:

```text
Persist valid packet
        ↓
Commit
        ↓
GetContext(task_id)
        ↓
assert beliefs == persisted beliefs
assert origin_packet_id == packet_id
assert edges/evidence visible
```

The review explicitly identifies this as the missing proof. 

### The lesson is important

This is exactly why the full-stack test matters.

The previous DB-backed tests proved:

```text
database rows exist
```

but not:

```text
ARGUS can reconstruct those rows through its authoritative context API
```

For this system, **persistence without reconstructability is not enough**.

The true invariant is:

```text id="9fs8e7"
packet submitted
   ↓
state persisted
   ↓
state reconstructable
   ↓
fresh agent can reason from it
```

So I agree with the review's readiness decision: **not ready yet**. 

### One thing I would deliberately defer

The review revisits the same-packet-ID/different-content ambiguity and notes that it currently silently preserves the first submission. 

I would **not expand scope to solve that now**.

Your bookkeeping freeze applies. Unless the next real work cycle demonstrates a concrete integrity problem from packet-ID reuse, leave it as a documented limitation.

### Coding-agent prompt

```text
# ARGUS — P0 Read-Path Repair

Fix ONLY the confirmed read-path regressions from the latest adversarial review.

DO NOT redesign ARGUS.
DO NOT add new MCP tools.
DO NOT add new bookkeeping infrastructure.

Confirmed issues:

1. P0: `internal/epistemic/view.go:scanBelief`
   GetSnapshot SELECT now returns 7 columns including origin_packet_id,
   but scanBelief scans only 6.

2. P1: `internal/work/store.go:ListByProject`
   SELECT returns 11 columns but Scan expects 12 including
   origin_packet_id.

3. P2: `internal/epistemic/view.go:GetAllBeliefs`
   SELECT omits origin_packet_id although BeliefView contains it.

Implement:

A. Fix scanBelief:
   scan:
     ID
     Claim
     ClaimType
     Status
     Debt
     FinalTruth
     OriginPacketID

B. Fix ListByProject SELECT to include origin_packet_id in the same position
   expected by Scan.

C. Fix GetAllBeliefs SELECT to include origin_packet_id.

D. Add one DB-backed end-to-end regression test:

   1. Persist a valid packet.
   2. Commit.
   3. Call GetContext(task_id).
   4. Verify the persisted beliefs are returned.
   5. Verify origin_packet_id is populated correctly.
   6. Verify evidence/edges are also reconstructable.
   7. Verify availability remains correct.

E. Run:
   go vet ./...
   go build ./...
   go test ./...

F. Start ARGUS with CockroachDB and manually verify:
   /ui/insights
   /ui/debts
   argus.get_context

Do NOT address:
- packet-id replay policy
- context snapshots
- review coverage persistence
- epistemic kind
- agent attestation
- additional provenance UI
- new graph features

Those remain deferred under the bookkeeping freeze.

Report:
1. exact files changed
2. exact read-path defects fixed
3. end-to-end GetContext test result
4. full test results
5. whether a fresh agent now sees the same beliefs that were persisted
```

This is still **load-bearing repair**, not bookkeeping creep: without `get_context` reconstructing the state, the Work → Adversarial → Work lifecycle you built cannot function. After this fix passes, I would again freeze the infrastructure and return to the research cycle.
