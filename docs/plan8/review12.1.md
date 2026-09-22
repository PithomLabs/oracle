## Z

Disposition of Review 9 against its own charter — this is the pre-freeze gate check, so the standard is whether its verdicts are earned, its gaps are correctly triaged, and nothing is missing that the freeze record needs. Verdict: **the review is sound, the verdict is correctly calibrated, and there are three things to fold into the freeze record — including one it left open that shouldn't stay open.**

## What the review earned

**The P1 finding is the right catch, correctly graded, and the most important one for the freeze.** The mock-path/production-path lesson (F4 of the code-review round, the round-one scanner bug) was applied one layer deeper: the end-to-end test calls `app.GetContext` directly, not the MCP stdio transport — so the *agent-facing* contract (JSON marshal of `Context`/`Snapshot` through `tools/call`) has zero coverage. `omitempty` silently dropping a field in the marshaled response is exactly the failure class that evaded the last two rounds. Correct P1: it doesn't block the code (which is verified sound at the application layer) but it does block the *freeze claim*, because the freeze contract says "agent-facing verification" and this is the only untested hop in the agent's path. The prescribed fix is right and small: one DB-backed test through `adapter.HandleTool("argus.get_context")` with JSON unmarshal and exact-field assertions.

**The nullable-`origin_packet_id` check (§4B) is the second-best catch.** The reviews had been testing packet-created objects exclusively; seed/human-created beliefs with NULL origins were an untested branch of the scanner. The review flagged it in the charter and confirmed the implementation handles it (`sql.NullString`). This closes the seed-versus-packet provenance distinction the dry run's Insights question originally exposed.

**The structural finding (§15/P2) is a genuine freeze-record item:** the UI test DB fix introduced *three copies* of solvent-migration application logic across test files — drift-prone by construction, the exact pattern that produced the original scanner bug. The prescription (centralize on `solventmigrations.Apply`) is right, and the review correctly grades it P2: non-blocking, but it should be *recorded* at freeze even if not fixed, because it's the next divergence waiting to happen.

**Bookkeeping check: NO, with justification** — correctly discharged. The review checked the freeze criteria against observed failures rather than vibe, and none triggered.

## Three things to fold into the freeze before it closes

**1. The P1 shouldn't be closed by the manual smoke test alone — the review half-says this and should say it fully.** The final paragraph proposes the manual smoke "should close the verification gap." But a manual smoke is exactly what the project's own test-quality doctrine exists to replace: it's a test that can't be re-run, can't regress-guard, and passes the wrong-reason trap (a human eyeballing JSON output won't catch an `omitempty` drop the way an assertion will). The manual smoke is *sufficient to proceed* only if the automated MCP-transport test is added alongside it — which is the review's own P1 prescription. So the freeze gate's checklist item should read: **"P1 MCP-transport test added AND manual smoke recorded"** — both, not either. This is a one-line amendment to the freeze gate; without it, the freeze inherits the same gap class (production-path-untested) that caused the last three review rounds.

**2. The dual-adapter finding (P2, legacy `mcp/adapter/` package) deserves one decision, not just a comment.** Two MCP adapters coexisting is precisely the condition that produced the original F4 (wrong-path tests passing against the wrong adapter). The review's smallest fix (deprecation comment) is acceptable under the freeze — but the freeze record should name *which* adapter is production (`internal/mcp`, via `cmd/argus/main.go`) explicitly, so the legacy package's existence is a documented landmine rather than an unexploded one. Trigger: any contributor importing it, or any test referencing it.

**3. The pre-existing `view_test.go` failures (§ P2, `defaultdb` / no migrations) conflict with the freeze checklist item "full suite green OR every remaining failure explicitly dispositioned."** These are known-failing tests in the tree at freeze time. The freeze gate can absorb them — but only via the explicit disposition path: each failing test either gets fixed (the review's own prescription: apply migrations in `testDB`) or gets a recorded reason for deferral. Seven red tests in the suite at freeze is the opposite of the "boring enough to freeze" standard the review's own final question invokes. Fold into the freeze gate as a checklist line: *disposition the seven `view_test.go` failures (fix or defer-with-reason) before the freeze tag.*

## One check the review should have run and didn't

The charter (§6) said "Do not assume empty debt unless actual Domain Pack/compiler behavior produces empty debt. Verify what the current implementation actually does." The findings never return to this — the P2 belief's debt expectation is asserted in the read-back test but the review doesn't confirm whether the compiler applies pack `initial_debt` on that path (the standing `ebpInitialDebt` obligation, and the fail-open protection from the code-review round's F1 family). This is the last unverified invariant in the freeze-critical test. One line for the freeze gate: *confirm P2's belief carries the pack initial debt set through Persist→GetContext* — or record it as untested-at-freeze with the fail-open trigger already documented.

## The verdict, and what it means for the sequence

**NOT READY — TEST/VERIFICATION GAP ONLY** is the correct call, and notably: it's the mildest verdict in the program's review history, and the first one where the *code* is fully exonerated. Every prior round found architecture or implementation defects; this round found only verification coverage. That's the empirical signature of the freeze being ripe — the code has stopped moving faster than the tests can catch it.

Freeze sequence, consolidated from this review plus the standing gate:

```text
FREEZE GATE — final items
[ ] P1: MCP-transport read-back test (HandleTool + JSON unmarshal, exact fields)
[ ] P1b: manual smoke recorded (both, per the amendment above)
[ ] Disposition the 7 view_test.go failures (fix via testDB migrations, or defer-with-reason)
[ ] Confirm/assert P2 belief carries pack initial_debt (or record as untested with trigger)
[ ] Legacy adapter: deprecation comment + production-adapter named in freeze record
[ ] Triple-duplicated migration helpers: record in freeze record as known divergence risk
[ ] Known-limitations register (§7/§8 items, triggers) — already drafted, include verbatim
[ ] Freeze decision record + annotated tag
```

Everything on that list is a day of work or less. Then the answer to the review's final question becomes an unqualified yes — and it's worth stating what the review's own charter implies: *"boring enough to freeze"* is the highest compliment this project's doctrine can pay a harness, because the same doctrine says the accounting must never become the work. ARGUS is now one test and a tag away from being exactly boring enough — at which point the project returns to its actual subject: whether an arithmetic floor can hold up the sky.


## consolidated

The review is **mostly net-valid**, and it sharpens the freeze gate appropriately. I would incorporate its three verification additions, but I would **not turn all of its P2 observations into pre-freeze implementation work**. That would violate the bookkeeping freeze we just locked.

## What I agree with

### 1. MCP transport test + manual smoke are both required

This is the strongest addition.

The current application-level `GetContext` test proves the application read path, but not the actual agent-facing serialization boundary. Review 9 explicitly identified that gap. 

Therefore the final gate should be:

```text
P1:
production MCP get_context test
    +
manual fresh-agent smoke
```

Not one or the other.

The automated test proves regression resistance; the manual smoke proves the real local wiring works.

### 2. The seven `view_test.go` failures must be dispositioned

I agree with this.

A freeze can tolerate a known P2 that is explicitly deferred, but **unknown/red tests in the final suite are not acceptable**.

The review says these failures are caused by the tests connecting to `defaultdb` without applying the Solvent schema. 

Because this is a relatively contained test-infrastructure defect, I would actually prefer:

```text
fix the tests
→ rerun
→ green suite
```

rather than carry seven red tests into the frozen baseline.

That is not "more accounting"; it is removing false confidence from the test suite.

### 3. Verify P2 initial debt

Also valid.

The plan deliberately said:

> don't assume empty debt; verify actual pack/compiler behavior.

The final test should establish:

```text
P2
→ packet compiler / domain pack
→ persisted belief
→ GetContext
→ exact expected debt
```

This is a **research/EPB invariant**, not bookkeeping, so verifying it is justified.

### 4. Legacy adapter: document, don't expand

The review's concern about two adapters is legitimate. It correctly identifies:

```text
internal/mcp/adapter.go
    = production

mcp/adapter/adapter.go
    = legacy
```

and the production binary uses the internal adapter. 

But I would **not require deletion or a new abstraction before the freeze**.

A short freeze-record line is enough:

```text
Production MCP adapter: internal/mcp.
Root mcp/adapter/ is legacy/non-production; revisit if a contributor imports
or tests it as the active adapter.
```

A deprecation comment is optional. Don't turn this into another cleanup project.

### 5. Duplicate migration helpers: record, don't fix

Valid finding, but **defer the implementation cleanup**.

The review correctly identifies three copies of test migration logic as drift risk. 

But centralizing them now is exactly the sort of "while we're here" work the bookkeeping freeze is designed to prevent.

Record it as:

```text
Known maintenance risk:
test migration setup is duplicated across test packages.
Trigger: migration divergence or repeated test-schema failure.
```

Then move on.

---

## One adjustment to the review

I would **not require all eight items in its proposed freeze list to be equal-status gates**.

Separate them:

### Actual freeze blockers

```text
[ ] MCP get_context transport test passes
[ ] Manual fresh-agent smoke passes
[ ] view_test.go failures are all fixed or explicitly dispositioned
[ ] P2 initial debt behavior is explicitly verified
```

### Freeze-record documentation

```text
[ ] network threat-model note
[ ] legacy adapter identified as non-production
[ ] duplicated migration helpers recorded as maintenance risk
[ ] limitations register
[ ] retained/deferred table
```

That preserves the "less is more" principle.

---

# Final freeze gate

I would now lock the final sequence as:

```text
                 WRITE
                   |
            MCP submit_packet
                   |
              Validate
                   |
                Persist
                   |
                Commit
                   |
                 READ
                   |
            MCP get_context
                   |
             exact state
                   |
          fresh-agent smoke
                   |
                FREEZE
```

And the final freeze condition is:

```text
No P0/P1 remains.
No unexplained failing test remains.
Agent-facing write/read path is proven.
Known P2/P3 limitations are documented with triggers.
No new accounting capability is introduced.
```

The review's central conclusion is therefore right: **Review 9 earned `NOT READY — TEST/VERIFICATION GAP ONLY`; after those four concrete verification items pass, the harness should be frozen.** The review itself found no justification for new accounting infrastructure. 

### Final instruction to the coding agent

```text
Complete only the final freeze-gate items:

1. Add DB-backed production MCP `argus.get_context` transport test:
   persist P1/P2/P3, call adapter.HandleTool("argus.get_context"),
   JSON-decode result, assert exact task/belief/evidence/edge/origin/
   availability state.

2. Fix or explicitly disposition the seven pre-existing
   `internal/epistemic/view_test.go` failures. Prefer fixing the test DB
   migration setup so the full suite is green.

3. Verify and assert the exact Domain Pack initial debt received by P2 through:
   packet → Persist → GetContext.

4. Record, do not implement:
   - root legacy MCP adapter
   - duplicated test migration helpers
   - network exposure limitation
   - other deferred limitations

Do not add new bookkeeping infrastructure.
Do not add MCP tools.
Do not change the authority model.

Then run:
    go vet ./...
    go build ./...
    go test ./...

Report the final freeze checklist and whether the harness is:

    READY TO FREEZE

or

    NOT READY TO FREEZE
```

After that, **stop**. The next iteration should be the BM–IST–AS research cycle, not another ARGUS review.
