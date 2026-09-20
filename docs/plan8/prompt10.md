Perform a FINAL ADVERSARIAL CODE REVIEW of the completed ARGUS POC development currently in the repository.

Do NOT assume the acceptance report is correct. Treat the actual code, Taskfile, migrations, tests, and runtime behavior as the source of truth.

Review against:
- frozen ARGUS POC architecture
- CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION
- Solvent authority boundary
- agent vs human authority separation
- domain-pack isolation
- verifier trust binding
- entity-level idempotency
- separate Solvent/work transactions
- migration ownership
- refusal-first acceptance philosophy

==================================================
1. REVIEW SCOPE
==================================================

Adversarially inspect:

- cmd/argus
- internal/application
- internal/epistemic
- internal/work
- internal/mcp
- internal/ui
- verifier
- domain-pack
- migrations
- Taskfile.yml
- .golangci.yml
- all tests added during Phase 0/1
- Solvent integration points
- process startup/shutdown behavior
- CockroachDB lifecycle
- HTTP server lifecycle
- environment/config handling

Do not merely rerun existing tests. Search for defects the tests do not cover.

==================================================
2. REQUIRED ATTACKS
==================================================

A. AUTHORITY / SECURITY

Attempt to determine whether an agent can indirectly:

- discharge debt
- promote beliefs
- retract beliefs
- authorize actions
- influence human attribution
- bypass the MCP two-tool boundary
- exploit HTTP/UI routes
- forge verifier identity
- submit an unauthorized verifier artifact
- exploit body-supplied principal fields
- bypass Origin/Host/CSRF checks
- reach Solvent kernel mutation methods through wrapper packages

Verify that RetireDebt is genuinely inaccessible from MCP/UI/application paths except the intended internal human discharge path.

Do not accept comments as enforcement.

B. EPISTEMIC INTEGRITY

Check:

- belief/evidence/edge scenario containment
- deterministic entity IDs
- duplicate packet behavior
- partial packet retry
- idempotency under restart
- conflicting packets with same content identity
- edge duplication
- evidence duplication
- whether ON CONFLICT can silently collapse semantically distinct entities
- whether proposed retirement metadata can alter authority
- promotion after retraction
- reopening lineage correctness
- dead-end linkage correctness

C. TRANSACTION / FAILURE SEMANTICS

The architecture intentionally does NOT claim global atomicity.

Verify that the implementation is honest about that.

Inject failures between:

- belief persistence
- evidence persistence
- edge persistence
- work-task persistence
- proposed-retirement persistence
- audit persistence
- idempotency completion

Verify retries converge without duplication or false success.

Check whether any failure can produce:

- false success response
- incomplete idempotency marker
- missing audit
- duplicated evidence
- orphaned task
- orphaned belief
- false refusal after a successful mutation

D. VERIFIER TRUST

Verify:

- VerifierID
- VerifierVersion
- VerifierHash
- InputHash
- InputSpec
- Tolerance
- canonical artifact hash
- verifier allowlist / minimum-version enforcement
- artifact input-to-claim binding
- timeout
- verifier error vs refuted result distinction
- CLI and in-process verifier use the same implementation
- no direct DB write bypass through `argus verify`
- no agent-selected arbitrary verifier

Attempt to tamper with each trust-binding field and ensure the appropriate validation fails.

E. DOMAIN-PACK ISOLATION

Verify that:

- core packages contain no BM-IST vocabulary
- BM-IST is injected rather than hard-coded
- PackDefinition is generic
- retirement rules remain domain-owned
- verifier specs remain domain-owned
- the dependency/import graph actually enforces the boundary

Run/inspect depguard AND the go/packages boundary tests.

Do not treat one as sufficient if the other is broken.

==================================================
3. CRITICAL DEVELOPER-ENVIRONMENT ATTACK
==================================================

The current `task dev` hardcodes ports and has already demonstrated a real failure:

- CockroachDB SQL: :26260
- CockroachDB HTTP: default :8080
- ARGUS HTTP: :8080

Therefore CockroachDB can fail before ARGUS starts because its HTTP admin port conflicts with ARGUS.

This MUST be treated as a defect in the developer experience.

Design and, if the review is authorized to patch implementation, implement a robust port-preflight mechanism.

Required behavior:

1. Before starting ANY process, run a port preflight.

2. Check ALL required TCP ports:
   - CockroachDB SQL
   - CockroachDB HTTP/admin
   - ARGUS HTTP/UI

3. If a preferred port is occupied:
   - automatically select the next available port from a deterministic local range
   - do not kill the existing process
   - do not silently reuse the occupied port

4. Reserve/choose all ports BEFORE starting services so the system does not discover conflicts halfway through startup.

5. Propagate dynamically selected ports everywhere:
   - Cockroach startup
   - Cockroach client DSN
   - ARGUS listen address
   - readiness probes
   - integration tests
   - generated config if any

6. Print the final resolved topology clearly, for example:

   CRDB SQL   :26260
   CRDB HTTP  :8082
   ARGUS HTTP :8080

7. Ensure a user already running another ARGUS/CRDB instance does not cause `task dev` to fail unnecessarily.

8. MCP remains stdio and requires no TCP port.

9. Add a dedicated preflight test that deliberately occupies the preferred ports and verifies fallback allocation.

10. Add a test ensuring the three selected ports are distinct.

11. Ensure cleanup releases only processes started by this invocation.

12. Ensure Ctrl-C / task interruption does not leave orphaned Cockroach or ARGUS processes.

13. Ensure stale PID/data state does not break the next invocation.

Prefer a small reusable port allocator/preflight helper over duplicated shell logic.

==================================================
4. TASKFILE ADVERSARIAL REVIEW
==================================================

Specifically inspect:

- task dev
- task fresh
- task build
- task test
- task verify
- task lint
- process cleanup
- signal handling
- readiness polling
- migration application
- connection-string propagation

Find races such as:

- readiness checking before server fully starts
- server dies after readiness succeeds
- child process survives Task termination
- stale data directory locks
- stale PID assumptions
- two concurrent `task dev` invocations
- selected port becomes occupied between preflight and bind

Where perfect race-free reservation is impossible, make the behavior explicit and retry startup with a newly selected port rather than failing mysteriously.

==================================================
5. TEST QUALITY REVIEW
==================================================

Audit whether existing tests genuinely test production behavior.

Look for:

- tests testing mocks rather than actual code paths
- tests asserting implementation details instead of invariants
- missing HTTP integration tests
- missing UI auth tests
- missing retry tests
- missing restart/idempotency tests
- missing port-conflict tests
- false-positive acceptance criteria
- tests that pass because code paths are never exercised

Run:

go test ./...
go test -race ./...
go vet ./...

Run available lint/boundary tooling.

If golangci-lint is incompatible with the current Go toolchain, report it precisely and verify whether dependency/version update or alternate supported invocation is the smallest solution. Do not mark lint fully PASS merely because boundary tests exist.

==================================================
6. CODE-LEVEL QUESTIONS
==================================================

Explicitly answer:

1. Can any untrusted agent input directly or indirectly mutate authority state?
2. Can any failed mutation report success?
3. Can any retry create duplicate logical state?
4. Can any retry silently discard a distinct logical entity?
5. Can a verifier artifact be fabricated or replayed incorrectly?
6. Can a packet reference another scenario?
7. Can UI attribution be forged?
8. Can old/retracted beliefs become promoted?
9. Can agent-authored edges cause authority transitions?
10. Can developer startup fail merely because another local process owns a preferred port?
11. Can task shutdown leak child processes?
12. Are Solvent schema ownership and migration ownership still intact?
13. Does the implementation introduce any third Solvent change?
14. Are the architectural package boundaries mechanically enforced?
15. Does the implementation still substantiate the thesis rather than merely making the demo run?

==================================================
7. OUTPUT FORMAT
==================================================

Produce:

# ARGUS FINAL ADVERSARIAL CODE REVIEW

## P0 — Blockers
Only genuine authority, correctness, data-integrity, security, or reproducibility failures.

## P1 — Important
Issues that should be fixed before Phase 9/productization.

## P2 — Hardening
Non-blocking improvements.

## Portability / Developer Experience
Specifically report port-conflict behavior, preflight behavior, fallback allocation, cleanup, and concurrent startup.

## Acceptance Audit
For each of the 25 acceptance criteria:
PASS / PARTIAL / FAIL
with exact evidence.

## Architecture Verdict
State one of:
- PASS — architecture preserved
- PASS WITH LIMITATIONS
- FAIL — architecture/invariant compromised

Do NOT recommend architectural expansion merely because a hardening feature is absent.

==================================================
8. IMPLEMENTATION RULE
==================================================

First perform the review.

If the port-resilience defect and any other P0/P1 defect can be fixed without changing the frozen architecture, include a concise implementation plan for those fixes.

Do not silently alter architecture, add services, add databases, or expand the Solvent change budget.

If a proposed fix would require a third Solvent change, STOP and report it as a Phase-0/governance violation rather than implementing it.