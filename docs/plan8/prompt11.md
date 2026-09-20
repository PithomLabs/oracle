The review is valuable, but it mixes **real defects in the current implementation** with several findings from the **superseded multi-service architecture**. We should not “fix” the latter by undoing the pivot.

### Net-valid findings

These should be fixed:

1. **Packet content hash is incomplete.** It must include the complete canonical packet, not just scenario/packet ID/claims. 

2. **Solvent migrations are not actually applied by `argus reset`.** This is a real reproducibility blocker. The fix is to use the approved Solvent migration export, not copy SQL or restore a Solvent HTTP service. 

3. **Conflicting edge kinds can be silently dropped.** Because the schema keys edges by `(parent_id, child_id)`, a second kind must produce an explicit conflict rather than `DO NOTHING`. 

4. **RCP must preserve `UNKNOWN ≠ EMPTY`.** A context query should distinguish unavailable state from genuinely empty state. The exact mechanism can be adapted to the single-DB architecture. 

5. **Debt retirement must be pack-rule validated before `Discharge`.** This remains a core semantic requirement. 

6. **Trust UI authentication is broken if it hardcodes a principal.** The current session/token model must actually be wired and enforced. 

7. **MCP stdio must actually run.** A process that prints a message and blocks is not an MCP implementation. 

8. **Semantic-version comparison must not be lexicographic.** Use real semver comparison. 

9. **Evidence hashes should be server-validated.** Do not trust an agent-supplied SHA-256 as authoritative provenance. 

10. **Port/process lifecycle is a real defect.** Preflight all required ports, choose deterministic fallbacks, propagate them, and clean up only processes started by the current invocation. The observed `:8080` CockroachDB/ARGUS collision validates this finding. 

11. Add `task down`, PID/process cleanup, proper signal propagation, and consistent CRDB test-port configuration. 

### Findings I would reject

**P0-2 — “App must use Solvent REST.” Reject.** That belongs to the old architecture. The whole pivot deliberately removed Solvent HTTP as a deployment boundary. In-process use of Solvent is the approved architecture.

**P1-4 — integration must test MCP → Coordinator → REST → Solvent/Conductor. Reject.** There is no Coordinator service anymore. The correct path is MCP → application layer → in-process Solvent/work boundaries.

**P1-5 — separate `trust-ui` module must remain the served UI. Reject.** The pivot explicitly converted Trust UI into embedded `internal/ui`.

**P1-6 — Coordinator HTTP client timeout. Reject as stale architecture.** Those REST clients should be disappearing, not hardened.

**P0-8 — `SubmitDecision` itself is an architectural violation.** Not necessarily. The application layer is the intended Coordinator replacement. The valid requirement is that it must enforce policy, authentication, and authority boundaries before reaching the Solvent kernel. The direct internal call is not itself the problem.

The review's final “FAIL architecture” therefore overstates the situation because it evaluates the implementation against portions of the superseded service architecture. The valid defects are still serious, but the remedy is **complete the pivot**, not resurrect Coordinator/REST. 

## Prompt for the coding agent

```text
Perform a corrective implementation pass based on the latest adversarial review.

IMPORTANT: preserve the frozen single-process ARGUS architecture. Do NOT restore Solvent/Conductor/Coordinator REST services, separate Trust UI, or internal service-to-service HTTP.

Fix only the net-valid defects:

1. Fix contentHash to cover the complete canonical packet:
   scenario, beliefs, evidence, edges, tasks, and all semantically relevant fields.
   Deterministic canonical ordering required.

2. Wire the real Solvent migration export into `argus reset`, with Solvent owning
   canonical epistemic migrations. No copied migration SQL in ARGUS.

3. Fix conflicting belief-edge kinds:
   never silently `DO NOTHING` when the same parent/child exists with a
   different kind. Return an explicit conflict/refusal.

4. Preserve UNKNOWN != EMPTY in GetContext. Backend/query failures must return
   explicit availability metadata rather than appearing as empty state.

5. Enforce the Domain Pack retirement rule before human Discharge:
   required debt item, evidence class, qualifying evidence, and rule must all
   match. Agent proposals remain proposals.

6. Fix Trust UI authentication:
   configured operator token -> authenticated session -> CSRF/Origin protection
   -> server-derived principal. Never trust body-supplied principal_id.

7. Fully wire `argus mcp` to real stdio transport with exactly:
   argus.get_context
   argus.submit_packet

8. Replace lexicographic verifier version comparison with semantic-version comparison.

9. Server-verify evidence content hashes instead of trusting agent-supplied hashes.

10. Add robust port preflight/startup:
    - check CRDB SQL port, CRDB HTTP/admin port, ARGUS HTTP port BEFORE startup
    - select deterministic fallback ports when occupied
    - reserve/resolve all ports before launching
    - propagate selected ports to DSNs, listeners, readiness checks, tests
    - print resolved topology
    - MCP remains stdio
    - never kill unrelated processes
    - clean up only processes started by this invocation
    - handle Ctrl-C without orphaning CRDB/ARGUS
    - add tests occupying preferred ports and verifying fallback behavior
    - ensure selected ports are distinct

11. Add `task down` and make task dev/fresh/startup lifecycle resilient.

12. Ensure integration tests use the resolved CRDB configuration consistently.

13. Add/adjust tests for every fix above.

Do NOT reintroduce:
- Coordinator REST service
- Solvent REST calls
- Conductor REST calls
- separate Trust UI process
- second database
- third Solvent change
- new architecture layers

After implementation run:
go test ./...
go test -race ./...
go vet ./...

Also run all integration, refusal, MCP, UI, verifier, and port-conflict tests.

Report PASS/PARTIAL/FAIL for each affected acceptance criterion and list any remaining limitation.
```

The key distinction now is **fix the implementation without reversing the architectural pivot**. The review correctly found several broken invariants, but its proposed REST-based remedy is obsolete.
