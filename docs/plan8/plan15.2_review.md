Plan 15.2 is **good and remains within the bookkeeping freeze**. I would approve it, with **three small changes before execution**. The substantive design is now settled; these are freeze-gate hygiene.

### 1. The MCP transport test is correctly scoped

The plan's test does exactly what the remaining gap requires:

```text
HandleTool
→ json.Marshal(result)
→ json.Unmarshal(...)
→ assert decoded fields
```

It therefore tests the actual `Context` JSON contract without creating a separate stdio test harness. The plan explicitly distinguishes this from the application-layer `GetContext` test. 

I would keep it exactly this way.

One useful detail: the test should assert that `json.Marshal(result)` itself returns no error. That sounds obvious, but it makes the transport-contract assertion explicit.

### 2. The `view_test.go` fix is appropriate

Using `solventmigrations.Apply` and an isolated test database is the right way to eliminate the seven false-red tests. It also reinforces the existing ownership rule that Solvent owns the Solvent schema. 

I would **not** centralize the other duplicated helpers during this freeze. The plan correctly leaves that as a documented maintenance risk. 

### 3. The P2 debt decision is correctly frozen as current behavior

This is important: the plan does **not** silently change production semantics to make the test prettier.

It explicitly establishes:

```text
packet debt
   ↓
Persist
   ↓
stored debt
```

and confirms `CompileDebt` is not currently invoked in the production path. 

That is the right decision under the bookkeeping freeze.

The test should therefore prove:

```text
P2 debt=nil
→ stored []
→ GetContext []
```

while documenting that this is **current behavior**, not an EBP rule.

---

## Three small additions

### A. Add `json.Marshal` error assertion

In `TestMCPGetContextReadsBackPersistedState`:

```go
raw, err := json.Marshal(result)
if err != nil {
    t.Fatalf("marshal get_context result: %v", err)
}
```

Then decode and assert the fields. The current plan implicitly assumes this, but making it explicit improves the test's purpose.

### B. Add a final "no unexplained failures" freeze checkbox

The current checklist has the important individual checks, but I would add:

```text
[ ] No unexplained test failures remain; every remaining failure is fixed
    or explicitly dispositioned in the freeze record
```

The plan says the oracle integration test *may* remain a pre-existing timeout. 

That's fine, but the freeze should never rely on prose like "may still fail." It should finish with an actual disposition.

### C. Add the freeze decision/tag to the checklist

The plan is explicitly called the final freeze-gate, so the final artifact should be:

```text
[ ] freeze decision record completed
[ ] limitations register completed
[ ] manual smoke result recorded
[ ] annotated freeze tag created
```

The plan already contains the limitations and manual-smoke material; this simply closes the loop.  

This is **not** new bookkeeping capability. It's the final record of a decision you're already making.

---

# I would approve Plan 15.2

The frozen architecture remains:

```text
Agent
  ↓
MCP
  ↓
ARGUS
  ├── validation
  ├── persistence
  ├── provenance
  └── RCP
        ↓
     get_context
        ↓
   fresh agent
```

And the freeze principle remains:

```text
concrete failure → smallest fix
no concrete failure → defer
```

The deferred list remains appropriately bounded: replay detection, snapshots, review coverage, epistemic kind, attestation, extra provenance UI, graph visualization, vector search, additional MCP tools, and network hardening. 

### Final instruction to the coding agent

```text
Proceed with Plan 15.2.

Three final refinements:

1. In TestMCPGetContextReadsBackPersistedState, explicitly assert that
   json.Marshal(result) succeeds before unmarshalling and asserting fields.

2. Add a final freeze-gate checkbox:

   [ ] No unexplained test failures remain; every remaining failure is
       fixed or explicitly dispositioned in freeze-record.md

3. Add final freeze artifacts to the checklist:

   [ ] freeze decision record completed
   [ ] limitations register completed
   [ ] manual smoke result recorded
   [ ] annotated freeze tag created

Everything else in Plan 15.2 remains unchanged.

Do not:
- add CompileDebt to Persist
- add new MCP tools
- add accounting infrastructure
- fix deferred P2/P3 items unless an actual freeze blocker emerges
- modify the authority model

After implementation run:

    go vet ./...
    go build ./...
    go test ./...

Then perform the human-run fresh OpenCode smoke test exactly as specified
in the freeze record.

Only after all freeze-gate items are actually checked may the final verdict
be READY TO FREEZE.
```

At this point, I would **not commission another broad adversarial architecture review after Plan 15.2** unless this execution uncovers a new P0/P1. The whole purpose of the freeze is to stop the harness from becoming the project.
