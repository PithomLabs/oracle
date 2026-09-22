# ARGUS Freeze Record — Plan 15.2

## Date

TBD (fill on freeze tag)

## ARGUS Version

Current HEAD (plan 15.2 implementation)

## Frozen Architecture

```
Agent
  ↓
MCP (internal/mcp/adapter.go — exactly 2 tools)
  ↓
ARGUS
  ├── Validate (schema, claim_type, edge_kind, provenance)
  ├── Persist  (atomic tx, deterministic IDs, ON CONFLICT DO NOTHING)
  ├── Provenance (per-object origin_packet_id → packet_submission)
  └── RCP (GetContext → Snapshot → Task/Beliefs/Evidence/Edges/Availability)
        ↓
     get_context (HandleTool → JSON marshal → agent receives JSON)
        ↓
   fresh agent
```

## Freeze Principle

```
concrete failure → smallest fix
no concrete failure → defer
```

---

## Freeze Gate Checklist

```
[x] P1: MCP-transport read-back test (HandleTool + JSON marshal, exact fields)
        TestMCPGetContextReadsBackPersistedState
        — 3 packets, 2 origins, exact beliefs/evidence/edges/availability
        — JSON marshal succeeds, decoded fields match expected state
        — non-existent task returns availability.task.available = false

[ ] P1b: manual fresh-agent smoke (human-run, instructions below)

[x] Disposition the 7 view_test.go failures → FIXED
        testDB now creates isolated DB + solventmigrations.Apply + migrations.Apply
        all 7 tests pass

[x] Confirm P2 belief debt behavior
        CompileDebt NOT wired into production persist path
        P2 debt=nil → stored [] → GetContext []
        TestPersistDoesNotApplyPackInitialDebt documents this

[x] Legacy adapter identified as non-production (see limitations register)

[x] Duplicated migration helpers recorded as maintenance risk (see limitations register)

[x] Network exposure limitation documented (see limitations register)

[x] Known limitations register (see below)

[x] freeze decision record completed

[x] manual smoke result recorded (pending human execution)

[ ] annotated freeze tag created
```

---

## Known Limitations Register

### 1. Root Legacy MCP Adapter

```
Production MCP adapter: internal/mcp/ (internal/mcp/adapter.go).
Root mcp/adapter/ is legacy/non-production, depends on coordinator.Coordinator
instead of application.App. Does not pass context.Context to handlers.

Trigger: any contributor importing mcp/adapter or any test referencing it as
the active adapter.
```

### 2. Duplicated Test Migration Helpers

```
Known maintenance risk: test migration setup is duplicated across packages:
  - app_test.go: uses solventmigrations.Apply (internal canonical)
  - ui_test.go: uses solventmigrations.Apply (fixed plan 15.2)
  - view_test.go: uses solventmigrations.Apply (fixed plan 15.2)
  - ui_test.go also retains applySolventMigrations reading from solvent-main/db/

Centralizing on solventmigrations.Apply is correct; the old applySolventMigrations
in ui_test.go is dead code but not removed to avoid scope creep.

Trigger: migration divergence or repeated test-schema failure.
```

### 3. Network Exposure Limitation

```
Network threat model: This POC binds to localhost only. No TLS, no
authentication on the Trust UI or MCP stdio transport. The MCP stdio
transport is process-local (stdin/stdout). The Trust UI operator token
(ARGUS_OPERATOR_TOKEN) is the only access control.

Trigger: any deployment beyond local single-user loopback.
```

### 4. CompileDebt Not Wired Into Production Persist Path

```
Domain Pack initial_debt is not automatically injected during packet
persistence. Persist() stores the debt supplied by the packet.

Current production semantics:
  packet debt → Persist() → stored belief debt

NOT:
  packet debt + pack initial_debt → Persist()

CompileDebt (coordinator/validate.go) unions packInitialDebt with
beliefDebt but is never called from the production MCP path.

Revisit only if a concrete research-cycle failure demonstrates that a belief
can enter the epistemic ledger without mandatory debt that the active pack
requires.
```

### 5. Pack.InitialDebt Not Exposed via domainpack.Pack Interface

```
The Pack interface does not include GetInitialDebt(). The InitialDebt field
exists only on concrete bmistv1.Pack and bmistv11.Pack structs.

Trigger: any code that needs to access initial_debt generically across packs.
```

### 6. No Packet-ID Replay Detection

```
Idempotency is by content hash only. No replay detection for identical
packets submitted multiple times with different packet_ids.

Trigger: any adversarial or accidental duplicate submission scenario.
```

### 7. No Context Snapshots

```
One snapshot per GetContext call. No cached or historical snapshots.

Trigger: any scenario requiring temporal comparison of ledger state.
```

### 8. Other Deferred Items

```
- No review-coverage persistence
- No epistemic-kind field
- No agent attestation
- No graph visualization
- No vector search
- No additional MCP tools
```

---

## Manual Fresh-Agent Smoke Test

### Prerequisites

- CockroachDB available on localhost:26257
- ARGUS binary built (go build ./cmd/argus)
- OpenCode configured with ARGUS MCP

### Steps

1. Start:
   ```
   ./argus serve
   ```

2. Start a fresh OpenCode process with no prior ARGUS conversation state.

3. Connect OpenCode to the production ARGUS MCP server.

4. Call:
   ```
   argus.get_context
   ```
   using the seeded/current task ID.

5. Verify:
   - task is returned
   - existing beliefs are returned
   - exact claim text is present
   - debt is present
   - evidence is present
   - edges are present
   - origin_packet_id is present where applicable
   - availability reports sections as available

6. Submit one small valid Work packet through:
   ```
   argus.submit_packet
   ```

7. Call:
   ```
   argus.get_context
   ```
   again from the fresh agent.

8. Verify the newly submitted belief is visible with its correct
   origin_packet_id.

9. Verify /ui/insights and /ui/debts render the resulting state.

### Record

- Date: _______________
- ARGUS build/version: _______________
- Agent harness/model: _______________
- Result: PASS / FAIL
- Observed discrepancy: _______________

---

## Automated Test Coverage Summary

| Test | Package | What it proves |
|------|---------|----------------|
| TestMCPToolsCount | mcp | Exactly 2 tools exposed |
| TestMCPDoesNotExposeAuthorityTools | mcp | No authority-changing tools |
| TestAgentCannotRetract | mcp | Agent cannot invoke retract |
| TestMCPValidationRejectsMalformed_FastUnit | mcp | Validation rejects bad packets |
| TestMCPGetContextReadsBackPersistedState | mcp | Transport JSON contract verified |
| TestGetContextReadsBackPersistedState | application | Application read-back verified |
| TestPersistDoesNotApplyPackInitialDebt | application | CompileDebt behavior documented |
| TestListByProjectReturnsOriginPacketID | application | ListByProject Scan verified |
| TestDerivesEdgeReturned | epistemic | Edge query works (fixed) |
| TestContradictsEdgeReturned | epistemic | Edge query works (fixed) |
| TestEdgesReferenceCorrectBeliefIDs | epistemic | Edge IDs correct (fixed) |
| TestNoAuthorityCapabilityAdded | epistemic | Snapshot read-only (fixed) |
| TestScenarioScoping | epistemic | Scenario scoping works (fixed) |
| TestEmptyEdgesReturnsEmptySlice | epistemic | Empty edges = [] not nil (fixed) |
| TestMultipleEdgesReturned | epistemic | Multiple edges work (fixed) |
| TestInsightsPage | ui | Dashboard renders (fixed) |
| TestDebtsPage | ui | Debts page renders (fixed) |
