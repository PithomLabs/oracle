## chatgpt

Plan 15.1 is **very close to the freeze candidate**. I would approve it with **one substantive correction and two small verification refinements**.

The plan correctly incorporates the earlier review: it fixes the three confirmed read-path defects, makes the `Agent → Claimant` distinction explicit, reclassifies the authority test at the MCP boundary, adds dual-origin read-back coverage, and freezes all non-triggered accounting work.   

## The one substantive correction

The freeze-critical test currently says:

```text
Persist
→ Commit
→ app.GetContext()
```

That is good, but the actual invariant you care about is the **agent-facing RCP boundary**:

```text
MCP get_context
→ production adapter
→ application GetContext
→ reconstructed state
```

The earlier P0 existed precisely because the application read path and the agent-facing context contract diverged. So I would not make the freeze gate depend solely on a direct `app.GetContext()` test.

Keep the application-level test—it is useful—but add the final assertion through the **real MCP tool**.

The freeze-critical path should be:

```text
write:
MCP submit_packet
    ↓
Persist
    ↓
Commit

read:
MCP get_context
    ↓
GetContext
    ↓
snapshot
```

That makes the actual invariant:

> A fresh agent can submit research and then reconstruct exactly that committed research through the same MCP interface it will use in production.

The plan already has real DB-backed persistence tests, so this is a small extension rather than new infrastructure. 

### Revised test structure

I would make `TestGetContextReadsBackPersistedState` do:

```text
1. persist P1
2. persist P2
3. persist P3
4. commit
5. call production MCP `argus.get_context`
6. decode returned snapshot
7. assert exact beliefs/evidence/edges/origins
```

You can keep the direct `app.GetContext()` assertion as an implementation-level test if it is useful, but the **MCP call is the freeze gate**.

---

## Small refinement 2: authority-test consolidation

The plan removes:

```text
TestAgentCannotPromote
TestAgentCannotDischargeDebt
TestAgentCannotRetractBelief
```

and replaces them with `TestMCPDoesNotExposeAuthorityTools`. 

That is conceptually correct **provided the replacement test asserts exactly two tools**, not merely absence of three names.

The invariant should be:

```text
MCP tools == {
  argus.get_context,
  argus.submit_packet
}
```

That is stronger and simpler.

No new tests for each nonexistent authority tool are required.

---

## Small refinement 3: the manual smoke test should explicitly test a fresh process

The freeze claim is about fresh-agent reconstruction.

So the manual step should be:

```text
argus serve
→ fresh OpenCode process
→ configured MCP
→ argus.get_context
```

not simply calling the endpoint from the same development process that created the packet.

That gives you the closest practical proof that the interface is genuinely reconstructable without hidden in-process state.

---

# What I would freeze

The freeze policy in Plan 15.1 is now exactly where it should be:

```text id="qhrdzt"
concrete failure → smallest fix
no concrete failure → defer
```

with deferred items explicitly named and trigger-based. 

I would **not add anything else**.

In particular, don't add:

```text
epistemic_kind
review coverage database
context snapshot IDs
agent attestation
packet replay machinery
graph visualization
vector search
new MCP tools
```

The plan explicitly keeps those deferred. 

That is consistent with your locked "accounting must never become the work" rule.

## Revised final instruction to the coding agent

```text id="bkrkpd"
Proceed with Plan 15.1 with ONE substantive change and two small verification
refinements.

SUBSTANTIVE CHANGE:

The freeze-critical read-back test must exercise the REAL production MCP
interface, not only app.GetContext() directly.

Required sequence:

    production MCP argus.submit_packet
        →
    real persistence + Commit
        →
    production MCP argus.get_context
        →
    returned snapshot
        →
    exact state assertions

Keep the direct application-level GetContext test only if useful as a lower
level regression test, but the MCP read-back path is the freeze gate.

The assertions must verify:
- task identity
- availability
- exact belief count
- exact claim text
- exact belief IDs
- exact origin_packet_id per belief
- two distinct packet origins
- exact evidence fields
- exact edges
- no empty/truncated snapshot

REFINEMENT 2:

`TestMCPDoesNotExposeAuthorityTools` should assert that the production MCP
tool set is EXACTLY:

    argus.get_context
    argus.submit_packet

and nothing else.

Do not add application-level authorization to SubmitDecision.

REFINEMENT 3:

Manual smoke verification must use a FRESH agent process after the packet
has been committed:

    argus serve
    fresh OpenCode process
    argus.get_context

Record the result.

Keep all other Plan 15.1 scope and bookkeeping-freeze constraints unchanged.

Do not implement deferred items.

After implementation run:

    go vet ./...
    go build ./...
    go test ./...

Then perform the manual fresh-agent smoke test.

Final output must state:

    READY TO FREEZE
or
    NOT READY TO FREEZE

with only concrete blockers if NOT READY.
```

### My freeze criterion

I would treat the harness as frozen when this statement becomes demonstrably true:

> **A fresh agent can submit a valid packet through `argus.submit_packet`, ARGUS can commit it atomically with provenance intact, and a fresh agent can retrieve the exact resulting epistemic state through `argus.get_context`.**

Everything else is secondary.

Plan 15.1 is therefore **approved after that one MCP read-back correction**. Once the freeze gate passes, I would stop reviewing ARGUS plumbing and move the adversarial attention back onto the BM–IST–AS claims.


## Z

Plan 15.1 review — disposition of the two routed findings, then the residue.

**F1 (network-exposure freeze note): Landed, correctly scoped.** §7 is the right paragraph — the threat-model line ("does not claim protection against an arbitrary network actor reaching an internal authority endpoint") plus the deployment trigger, explicitly not implementing network hardening unless a POC test demands it. The earlier Critical finding is now a recorded, triggered limitation rather than an ambient assumption. This is the correct resolution.

**F2 (read-back test strengthening): Landed, and beyond spec.** §4 didn't just add the two suggested assertions — it upgraded the whole test to three packets with two distinct agents, cross-packet edges, exact content equality on every field, and the explicit `P1 != P2` provenance assertion. That's the dual-origin pinning requested, plus edge reconstruction through GetContext, which the original didn't attempt. The invariants-proven list correctly names what each assertion guards.

**F3 (limitations register): Landed** — §8, all six items with triggers. **F4 (column-order verification): Landed** — §1B now shows the verified three-way comparison table rather than the assumption. **Checklist extensions: Landed** — §13 includes the smoke-test row and limitations register. **Checklist and deferred-list alignment:** complete and consistent with the triage tiers.

## Residual findings — two, both pre-execution checks rather than changes

**R1 — §3 removes three tests but the replacement coverage is asymmetric.** `TestAgentCannotPromote` is replaced by `TestMCPDoesNotExposeAuthorityTools` (a *new, stronger* test asserting the exact tool list and the absence of named authority tools). But `TestAgentCannotDischargeDebt` and `TestAgentCannotRetractBelief` are removed on the grounds that "the existing `TestMCPToolsCount` and `TestAgentCannotRetract` in `adapter_test.go` already cover these invariants." Verify that claim before deletion rather than trusting it: does `TestAgentCannotRetract` in adapter_test.go actually exist and assert tool-surface absence (not mock behavior)? If it's another application-layer misplaced test, removing all three while adding one replacement nets the suite *thinner* on authority-boundary coverage than before. Thirty-second grep; if the coverage claim fails, fold the discharge/retract names into the new `TestMCPDoesNotExposeAuthorityTools` absence-list (one more line each — the test already has the right shape for it).

**R2 — §4 step 7's edge packet has no evidence or debt-bearing belief.** P3 carries only an edge. Two small consequences: (a) if `Persist` requires evidence to reference an existing belief and validates belief-debt presence, an edge-only packet may be the *first* exercise of that validation path — fine if it passes, but if it fails, the freeze-critical test fails for a reason unrelated to the read path being tested; (b) the adversarial packet P2 declares a belief with no debt items specified — confirm the compiler correctly applies pack `initial_debt` on this path (that's the standing `ebpInitialDebt` obligation, and this is its cheapest test ride: assert P2's belief comes back with the full initial debt set, not empty). One assertion, and the freeze test now also covers the fail-open protection that was F1 of the code-review round. If edge-only packets turn out to be invalid per the packet schema, restructure P3 to carry the edge plus a minimal belief — same coverage, no spec fight.

## Everything else

The three fix code blocks are correct as written and consistent with the verified column order. §2's grep discipline (keeping the correct "Agent" header on the submissions table) is exactly the right scope. §6's rename-with-comment preserves the cheap test per the freeze doctrine. §9's boundary sentence ("research methodology and EBP scientific obligations remain unfrozen") is present and correctly placed. The checklist is now complete against the freeze-gate requirements. Nothing in the plan introduces accounting, tools, or authority changes; the deferred list and triggers are consistent across §7, §8, §9, and §14.

## Verdict

**Approved for execution.** R1 is a verification step, R2 is either an assertion addition or a test-structure tweak — both fold into the plan as written without re-approval. When §13's boxes are ticked and the recommendation slot is filled, write the freeze decision record with the four artifacts it now has all the inputs for: the limitations register (§8), the threat-model note (§7), the retained/deferred tiers, and the reconciliation supersedes-list (§24 items from earlier plans — state machine, discharge path, UI module location). Then tag it.

And the closing note stands as predicted two rounds ago: the harness's defect classes went mechanical, this plan contains zero architecture, and the only remaining question in §13 is one the run itself answers. After the freeze, the next thing in this repository that deserves adversarial review is a belief about physics — entered free, in debt, waiting for its first discharge.


## consolidated review

The review is **net-valid** and I agree with its disposition. It adds two useful pre-execution checks, but neither changes the architecture or breaks the bookkeeping freeze.

The consolidated position is:

### R1 — verify authority-boundary test coverage before deleting the old tests

This is valid.

The plan is replacing three application-layer tests with the single canonical boundary test:

```text
TestMCPDoesNotExposeAuthorityTools
```

That is the right abstraction, but the claim that the existing `TestMCPToolsCount` and `TestAgentCannotRetract` already cover discharge/retract must be verified before deleting anything.

The clean invariant should be one test:

```text
production MCP tool set == {
    argus.get_context,
    argus.submit_packet,
}
```

and therefore:

```text
no promote
no retract
no discharge
```

So I would **not preserve three separate tests** merely for symmetry. First grep/inspect the existing tests. If coverage is insufficient, extend the one canonical tool-surface test.

That follows "less is more."

### R2 — strengthen the freeze test around P2's initial debt

This is also valid, with one distinction.

The proposed P3 edge-only packet is useful because it tests edge reconstruction separately, but the important thing is that the packet must be valid under the actual packet/pack semantics. The reviewer is right to check that rather than assume it.

More valuable is the P2 assertion:

```text
P2 belief
    → compiler/pack applies initial_debt
    → GetContext
    → debt == expected initial debt
```

That tests a real existing invariant and costs essentially nothing.

I would add that assertion.

I would **not automatically add a belief to P3** unless the actual validator requires it. Let the source code decide that.

### The rest of the review

I agree with the remainder.

The threat-model note is now explicit and appropriately limited. It doesn't attempt to solve network hardening during the freeze. The limitations register is also sufficient without becoming another subsystem. The read-back test is now strong enough to exercise:

```text
P1 → belief/evidence
P2 → different belief/evidence
P3 → edge
       ↓
GetContext
       ↓
exact reconstruction
```

which is exactly what we need. 

And the plan correctly keeps the bookkeeping freeze intact: no snapshots, coverage persistence, epistemic-kind field, attestation, graph UI, vector search, new MCP tools, or extra authority layers. 

## Final disposition

**Approve Plan 15.1 for execution with two small pre-execution checks:**

```text
R1:
Verify existing MCP authority-surface test coverage before deleting the
three old application-layer tests. Extend the single tool-surface test if
needed.

R2:
Verify P3's edge-only packet is valid under current packet semantics, and
add an assertion that P2's belief receives the expected pack initial debt.
```

No architectural changes are needed.

### One final implementation instruction

Give the coding agent:

```text
Proceed with Plan 15.1.

Before deleting/replacing the authority tests:

1. Inspect `TestMCPToolsCount` and `TestAgentCannotRetract` in
   internal/mcp/adapter_test.go.
2. Verify they test the production MCP tool surface, not application-layer
   behavior.
3. If discharge/retract absence is not fully covered, extend the single
   `TestMCPDoesNotExposeAuthorityTools` test so the exact MCP tool set is:

       argus.get_context
       argus.submit_packet

   and no authority-changing tools are exposed.
4. Do not add application-level authentication to SubmitDecision.

For `TestGetContextReadsBackPersistedState`:

5. Verify P3's edge-only packet is valid under the current packet/pack rules.
6. If it is valid, keep it.
7. Assert that P2's belief receives the expected Domain Pack initial debt
   after submission and read-back.
8. Do not add unrelated infrastructure.

Then execute Plan 15.1 exactly as specified.

The freeze gate remains:

    MCP submit
      → persist
      → commit
      → fresh MCP get_context
      → exact state reconstruction

Do not implement deferred bookkeeping features.
```

After that, the standard is very simple: **run it, inspect the actual freeze checklist, and freeze only if the evidence says `READY TO FREEZE`.** The next substantial adversarial review should then be of the research claim, not the harness.
