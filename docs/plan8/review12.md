This final review is **largely a green signal**. It confirms that the important architecture is now sound:

* write path works atomically;
* provenance survives;
* exact retries are safe;
* `GetContext` works correctly at the application layer;
* `UNKNOWN != EMPTY` is preserved;
* the MCP surface is exactly two tools;
* no new bookkeeping capability is justified.   

The only **freeze-relevant gap** is the one remaining at the actual MCP transport boundary.

### The one thing I would still do

The review found:

```text
application GetContext
    ✓

MCP HandleTool("argus.get_context")
    not directly tested with DB + JSON result

stdio JSON-RPC
    not automated-tested
```

The application read path is verified, but an incorrectly marshaled `Context` could still theoretically make the agent-facing response wrong even though all application tests pass.  

Because **your freeze criterion is explicitly agent-facing reconstruction**, I would add exactly **one small DB-backed MCP `get_context` integration test**. This is not new accounting; it is a verification of an already-frozen capability.

Do **not** start fixing the legacy adapter, centralized test migration helpers, or the POC schema duplication now. Those are legitimate P2/P3 findings, but they do not block the current research loop. The review itself classifies them that way. 

### Give the coding agent this final, very narrow prompt

```text
# ARGUS — Final MCP Read-Transport Verification

Do NOT redesign ARGUS.
Do NOT add any new capability.
Do NOT add bookkeeping infrastructure.
Do NOT fix unrelated P2/P3 findings.

The final adversarial review found exactly one remaining freeze-relevant gap:

`TestGetContextReadsBackPersistedState` verifies application.GetContext()
directly, but no DB-backed test verifies the production MCP get_context
response through the actual agent-facing adapter/JSON path.

Add ONE focused integration test.

## Required test

Create:

    TestMCPGetContextReadsBackPersistedState

Use a real CockroachDB and the production MCP adapter.

Steps:

1. Persist the existing three-packet test state:
   - P1 work belief + evidence
   - P2 adversarial belief + evidence
   - P3 cross-packet contradicts edge

2. Commit the state.

3. Call the production MCP path:

    adapter.HandleTool(
        "argus.get_context",
        args containing task_id
    )

4. Verify the returned MCP result is successfully marshaled and decoded.

5. Assert exact:
   - task ID
   - project/scenario identity
   - belief count
   - belief IDs
   - exact claim text
   - status
   - debt
   - origin_packet_id for P1 and P2
   - evidence count and fields
   - edge count and parent/child/kind
   - availability metadata

6. Specifically verify that a GetSnapshot/read failure cannot appear as a
   successful empty snapshot.

The test must exercise the production MCP result serialization path, not
only app.GetContext().

## Important

Do NOT add stdio-process machinery if HandleTool already represents the
production MCP adapter contract used by cmd/argus/mcp_stdio.go.

Only add a full stdio test if the current architecture proves that important
serialization behavior occurs exclusively outside HandleTool.

## Verification

Run:

    go vet ./...
    go build ./...
    go test ./...

With CockroachDB running.

Do not address:
- legacy MCP adapter cleanup
- duplicated test migration helpers
- schema-source harmonization
- packet-ID replay detection
- context snapshots
- review coverage
- epistemic-kind
- agent attestation
- network hardening
- any new MCP tools

Final report:

    MCP read transport verified: YES/NO
    Ready for final manual smoke: YES/NO
    No new accounting capability added: YES
```

### Then freeze

After this test passes, I would **stop the code-review loop**.

The remaining P2/P3 items should be recorded as known limitations, not used to reopen the harness. The adversarial review explicitly found no reason to add new accounting infrastructure. 

The freeze state should be:

```text
ARGUS
  ├─ validate
  ├─ persist
  ├─ commit
  ├─ provenance
  ├─ get_context
  ├─ MCP capability boundary
  └─ human control surface
        ↓
      FROZEN
```

Then the next adversarial review should be directed at **G0 / Deliverable A / the BM–IST–AS reasoning**, not ARGUS.

That is now the highest-value use of the system.
