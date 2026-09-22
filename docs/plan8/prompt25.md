# ARGUS — FINAL SCOPED ADVERSARIAL REVIEW
## Freeze Candidate After Plan 15.2

You are performing the final adversarial code review of the CURRENT ARGUS
repository immediately before the human fresh-agent smoke test and freeze.

DO NOT modify code.

DO NOT redesign ARGUS.

DO NOT add features.

DO NOT reopen the bookkeeping freeze.

DO NOT propose deferred capabilities unless a concrete P0/P1 integrity failure
makes one necessary.

The current implementation reports:

- `go vet ./...` clean
- `go build ./...` clean
- `go test ./...` = 253 passed, 0 failed, 2 skipped
- results stable across two runs
- automated freeze-gate items complete
- manual fresh-agent smoke test is the only remaining human action

Recent additional fixes in the same cycle:

1. `internal/work/store.go`
   - GetByID/ListAll/ListByProject now handle nullable task description via
     `sql.NullString`.

2. Integration test:
   - cross-type `local_id` collision test was fixed so the adversarial edge
     no longer reuses a belief local_id.

3. Integration test DB setup:
   - switched from external `solvent-main/db/` migration reading to
     internal `solventmigrations.Apply`.

4. `GetSnapshot`:
   - edge slice now initializes to `[]EdgeView{}` rather than nil.

Review the ACTUAL code, not the implementation report.

==================================================
1. PRIMARY FREEZE QUESTION
==================================================

Answer:

> Is there any remaining concrete P0/P1 defect that would make the ARGUS
> harness unsafe or incorrect for the intended Work → Adversarial → Human
> research cycle?

The desired answer is:

    NO

Do not turn P2/P3 maintainability observations into freeze blockers.

==================================================
2. FREEZE-CRITICAL WRITE → READ LOOP
==================================================

Trace the actual production path:

    argus.submit_packet
        ->
    validation
        ->
    transaction
        ->
    Persist
        ->
    Commit
        ->
    argus.get_context
        ->
    GetContext
        ->
    GetSnapshot
        ->
    JSON response to agent

Verify:

- valid packet commits
- malformed packet is rejected before persistence
- transaction is atomic
- provenance is preserved
- committed beliefs/evidence/edges/tasks are reconstructable
- read errors remain distinguishable from empty research state
- no in-process cache is masking DB state
- MCP JSON serialization preserves the fields already verified at the
  application layer

The intended invariant remains:

    UNKNOWN != EMPTY

==================================================
3. REVIEW THE FREEZE-CRITICAL MCP TRANSPORT TEST
==================================================

Inspect:

    TestMCPGetContextReadsBackPersistedState

Verify it actually:

1. creates real DB state
2. uses the production MCP adapter
3. calls `argus.get_context`
4. marshals the returned result with `json.Marshal`
5. unmarshals it
6. asserts exact fields

Verify assertions cover:

- task.id
- task.project_id
- beliefs
- exact claims
- exact IDs
- status
- debt
- origin_packet_id
- evidence
- edges
- availability

Verify the test does NOT merely prove that the application-layer
`GetContext()` works.

==================================================
4. NULL / SEED PATH REVIEW
==================================================

Review the `origin_packet_id` read path for BOTH:

A. packet-created object:
   origin_packet_id != NULL

B. seed/human-created object:
   origin_packet_id == NULL

Verify all scanners and views handle nullable origin_packet_id correctly.

Specifically inspect:

- scanBelief
- GetSnapshot
- GetAllBeliefs
- task readers
- evidence readers
- edge readers

Look for:

- NULL scan failures
- accidental empty-string substitution
- pointer/null confusion
- JSON `omitempty` behavior that changes semantics

==================================================
5. NEW NULLABLE DESCRIPTION FIX

Review the recent `sql.NullString` changes for task descriptions.

Inspect:

- GetByID
- ListAll
- ListByProject
- GetTasksByGovernanceRef
- corresponding Scan logic
- struct field type/semantics

Verify:

A. description present
B. description NULL

both work.

Confirm the change does NOT alter existing task semantics.

Do not turn this into a broader nullable-field refactor.

==================================================
6. LOCAL_ID COLLISION INTEGRITY

Review the corrected integration test and the actual packet validator/Persist
logic.

The invariant is:

    local:<id>

has one unambiguous target within one packet.

Verify global uniqueness across:

- beliefs
- evidence
- edges
- tasks

Verify cross-type collision is rejected.

Specifically inspect the test that previously reused a belief local_id for an
adversarial edge.

Ensure the test now proves a REAL collision and would fail if global
local_id uniqueness were removed.

Do not add new persistence tables.

==================================================
7. EDGE VALIDATION + EDGE SNAPSHOT

Review:

- edge endpoint validation
- same-scenario validation
- self-edge rejection
- canonical contradicts target
- `strings.EqualFold` UUID comparison
- EdgeSnapshot initialization

Verify:

- empty edge set serializes consistently as `[]`
- non-empty edge set preserves exact fields
- case differences in textual UUID representation do not weaken scenario
  isolation
- missing endpoint remains rejected
- cross-scenario endpoint remains rejected
- self-edge remains rejected

Do not change edge semantics.

==================================================
8. PROVENANCE

Verify:

    packet_submission
        ->
    belief/evidence/task/edge_provenance

still preserves exact origin_packet_id.

Verify:

- packet_submission inserted first
- foreign keys intact
- retries do not overwrite origin
- duplicate deterministic entity IDs preserve original provenance
- edge provenance corresponds to the actual edge

Do not reopen packet-ID replay policy.

It remains deferred unless a concrete integrity failure is demonstrated.

==================================================
9. IDEMPOTENCY

Verify:

A. exact retry:
    same packet_id
    same content
    same identity

    -> no duplicates
    -> original provenance preserved

B. current documented packet-ID reuse limitation remains understood.

Do NOT redesign idempotency during this review.

==================================================
10. MCP AUTHORITY SURFACE

Verify the production MCP tool set is EXACTLY:

    argus.get_context
    argus.submit_packet

Nothing else.

Verify:

- no promote
- no retract
- no discharge
- no authorize

is exposed through the production MCP adapter.

Do not reintroduce application-layer authentication to SubmitDecision merely
to satisfy old tests.

The architecture remains:

    agent authority
        -> MCP capability surface

    human/control authority
        -> existing operator/control path

==================================================
11. MIGRATION OWNERSHIP / TEST DATABASE

Review the recent change from external migration files to
`internal/solventmigrations.Apply` in the integration test path.

Verify:

- tests use the same canonical Solvent migration owner expected by ARGUS
- no hidden test-only schema has been introduced
- no schema duplication was added
- provenance tables required by current ARGUS exist in the test DB
- production and test schema ownership remain coherent

Do NOT attempt to centralize all historical test helpers in this pass.

If duplication remains, classify it P2/P3 only.

==================================================
12. UUID NORMALIZATION

Review the recent `normalizedUUID` / UUID comparison changes.

Verify:

- production identity semantics are unchanged
- EntityID remains deterministic
- normalization is used only for representation/equality where intended
- different UUID values cannot become equal
- upper/lowercase representation differences are handled safely

Do NOT revisit UUIDv7 design.

==================================================
13. TEST QUALITY

The report says:

    253 passed
    0 failed
    2 skipped
    stable across two runs

Verify the claim.

Identify ONLY tests that could provide false confidence about a freeze-critical
path.

Pay particular attention to:

- tests using nil DB
- application-level tests substituting for MCP tests
- tests that assert counts but not content
- tests that don't validate provenance
- tests that can pass while JSON transport is broken

Do NOT demand every test become end-to-end.

The goal is one trustworthy test per load-bearing boundary.

==================================================
14. PRE-EXISTING / SKIPPED TESTS

Inspect the two skipped tests.

For each state:

- why skipped
- whether the skip is intentional
- whether it affects a freeze-critical invariant

No unexplained skipped test should be allowed to hide a P0/P1.

P2/P3 skips may remain if explicitly documented.

==================================================
15. BOOKKEEPING FREEZE

The following remain OUT OF SCOPE:

- packet-ID replay detection
- context snapshots
- review coverage persistence
- epistemic-kind
- agent attestation
- additional provenance UI
- graph visualization
- corpus/vector search
- new MCP tools
- additional authority layers
- network hardening
- CompileDebt integration into Persist

Do NOT recommend them.

Apply:

    concrete failure → smallest fix
    no concrete failure → defer

Research methodology and EBP scientific obligations remain completely
unfrozen.

==================================================
16. FINAL SMOKE READINESS

Determine whether the repository is ready for:

    ./argus serve

followed by:

    fresh OpenCode process
        ->
    argus.get_context
        ->
    argus.submit_packet
        ->
    argus.get_context

Do NOT perform the manual smoke test.

Determine code readiness for it.

==================================================
17. REQUIRED OUTPUT

Return:

### EXECUTIVE VERDICT

One of:

    READY FOR FINAL HUMAN SMOKE
    NOT READY — P0/P1 REMAINS
    NOT READY — VERIFICATION GAP ONLY

### WRITE PATH
Can a valid packet commit atomically?

### READ PATH
Can the committed state be reconstructed exactly through production
`argus.get_context`?

### MCP AUTHORITY
Is the production MCP surface exactly the two intended tools?

### PROVENANCE
Is origin_packet_id correct and preserved?

### RECENT FIXES
Explicitly assess:

- nullable description handling
- local_id collision test
- migration-source change
- EdgeSnapshot empty-slice initialization
- UUID normalization

### TEST QUALITY
Identify only meaningful false-confidence gaps.

### SKIPPED TESTS
Explain the two skips and whether they matter.

### BOOKKEEPING CHECK
Did this review find ANY justified reason to add new accounting
infrastructure?

Default should be:

    NO

### FINAL DECISION

Answer:

> Is ARGUS now boring enough to freeze after the human smoke test?

Do not modify code.
Do not propose architecture expansion.
Do not turn P2/P3 cleanup into blockers.
Prefer concrete runtime/code evidence over speculation.