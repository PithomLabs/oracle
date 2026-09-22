This adversarial review is valuable, and I would **take its main conclusion seriously: the current production MCP path is not yet ready for another research dry run**. The most important finding is not provenance now; it is the validation boundary.

The review says the production MCP path does:

```text
JSON
 ↓
Compile
 ↓
App.Validate
 ↓
Persist
```

while the stronger `packet/v1.Validate()` is never invoked. That leaves several packet-grammar invariants unenforced at the actual agent boundary. 

### What I agree with strongly

**1. The MCP validation bypass is the primary blocker.**

The review identifies concrete examples an agent could submit: missing agent identity, role mismatch, duplicate `local_id`, invalid claim type, invalid edge kind, malformed references, and fabricated provenance classes. 

That directly violates the principle we established:

```text
agent-facing boundary
    =
fail closed
```

Fix this before another dry run.

**2. The `@` vs `-` pack-reference bug must be fixed at the same time.**

The review found that `packet/v1.Validate()` parses `bmist-1.0.0` while the real wire format is `bmist@1.1.0`. So simply wiring that validator into MCP would break legitimate packets. 

That is a good adversarial finding.

**3. Duplicate local IDs are a real corruption risk.**

This one is especially important:

```text
b1 → claim A
b1 → claim B
```

could create two persisted beliefs while the internal map points `local:b1` at only the second one. 

That violates the packet's internal referential integrity. Reject duplicates.

**4. Agent-role mismatch must fail closed.**

The new provenance mechanism is only useful if:

```text
packet.role == agent.role
```

is actually enforced on the real MCP path. The review correctly identifies that the persistence layer currently accepts inconsistent provenance if structural validation is bypassed. 

### Findings I would *not* accept blindly

There are a couple of places where the adversarial review itself needs reconciliation.

**Edge validation is internally inconsistent in the report.**

One section says cross-scenario edges, local self-challenges, etc. are allowed; later it says the runtime `Persist()` path already performs self-edge/cross-scenario/local-contradicts checks.  

So before implementing another edge fix, have the agent reproduce each case through the actual production MCP path:

```text
same-scenario canonical
cross-scenario canonical
local contradicts
self-edge
invalid edge kind
```

The final answer should be based on actual behavior, not the review's conflicting descriptions.

**The FK-ordering finding is also a hypothesis, not yet a demonstrated failure.**

The review notices that objects may reference `packet_submission` before it is inserted, but the current tests are passing. 

Don't redesign it on speculation. Run a focused transaction test that proves whether the FK is immediate. If it works under CockroachDB, document why; if not, move `packet_submission` insertion earlier.

### Task/project handling

The task finding is important, but I would avoid auto-creating projects. The current architecture seems to be using:

```text
scenario_id == project_id
```

as a convention. The cleaner POC rule is:

> A packet that creates tasks must reference a scenario that already has the corresponding project.

Fail early rather than silently provisioning another object.

The review correctly identifies that a task persistence failure must not produce partial packet success. 

Your atomic transaction is otherwise a strong foundation; the review confirms the production MCP path currently uses one transaction for `Persist()`. 

### Things I would defer

These matter, but they are not the next blocker:

`contentHash` excluding agent identity. I would actually **keep it this way for now**. Identical research content from two agents should deduplicate as research state, while `packet_submission` still preserves the fact that two agents submitted it. That's a legitimate semantic choice, not necessarily a defect. The review itself describes it as a decision point rather than a proven failure. 

Evidence/edge provenance missing from UI is real, but it is a projection gap rather than an authority failure. 

The old duplicate `mcp/adapter/adapter.go` should eventually be removed if confirmed unused. It's a maintainability hazard because future developers can modify the wrong adapter. 

`context_snapshot_id` remains a useful future feature, but not a reason to stop this POC. 

## Recommended next pass

I would give the coding agent this **minimal corrective prompt**, rather than asking it to implement the whole review blindly:

```text
Perform a focused remediation of the adversarial code-review findings.

Do not broaden architecture.

FIRST, reconcile the review against the actual current code. Reproduce the
following through the production `argus.submit_packet` MCP path:

1. missing agent.id
2. missing agent.harness
3. agent.role != packet.role
4. duplicate belief local_id
5. duplicate evidence local_id
6. duplicate edge local_id
7. invalid claim_type
8. invalid edge kind
9. malformed belief references
10. local contradicts target
11. cross-scenario edge
12. self-edge

Do not rely on the previous review's claims where they conflict internally.

Then implement only the confirmed issues.

Required:

A. Production MCP validation
Wire the canonical packet structural validator into the actual MCP submission
path. Make sure it uses the real `@` pack-ref grammar and `bmist@1.1.0`.

B. Defense in depth
If structural validation can be bypassed by another current application
entrypoint, ensure the application layer still rejects malformed references,
invalid edge kinds, and duplicate local IDs.

C. Atomicity
Confirm that any failure while persisting beliefs, evidence, edges,
edge_provenance, tasks, or packet_submission rolls back the entire packet.

D. Task integrity
For packet-created tasks:
- use the existing scenario/project contract
- fail clearly if the required project does not exist
- never silently lose a task
- do not auto-create unrelated projects

E. Agent provenance
Require:
agent.id
agent.role
agent.harness
agent.model

and require:

agent.role == packet.role

F. Do not implement yet:
- context snapshots
- new MCP tools
- agent authentication/attestation
- content-hash redesign
- new review-coverage persistence
- graph UI
- new authority mechanisms

Add production-path regression tests for every confirmed failure.

Run:
go vet ./...
go test ./...
go build ./...

Then report:
1. Which review findings were confirmed.
2. Which were false/duplicate/internally inconsistent.
3. Exact files changed.
4. Exact invariants now enforced.
5. Full test results.

Do not claim readiness for another dry run until the production MCP path
itself has been tested with malformed packets.
```

The key takeaway from `adv_review6` is that the system has progressed from **“does the research loop work?”** to **“can the agent-facing boundary be trusted under adversarial input?”** The latter is the right thing to fix now. 

Once that pass is green, the next Work → Adversarial → corrective Work cycle will be much more meaningful because you can trust that a malformed packet cannot silently contaminate the research ledger.
