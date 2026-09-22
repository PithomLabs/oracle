# Plan: Research Provenance, Adversarial Targeting, and Submission Integrity Audit

Date: 2026-09-21
Status: Audit findings and implementation plan

---

## 1. BELIEF / OBJECT PROVENANCE

### Decision Table

| Object | Current origin link | Persistent? | Exact location | Recommended location |
|--------|---------------------|-------------|----------------|----------------------|
| belief | NONE | NO | belief table has no packet_id column | kernel-adjacent: add `origin_packet_id UUID` to belief table |
| evidence | NONE | NO | evidence table has no packet_id column | kernel-adjacent: add `origin_packet_id UUID` to evidence table |
| edge | NONE | NO | belief_edge has no packet_id or scenario_id | additive provenance table (belief_edge has no scenario_id, schema is minimal) |
| task | governance_ref (optional, points to belief, not packet) | NO | conductor_task has no packet_id | kernel-adjacent: add `origin_packet_id UUID` to conductor_task |

### Analysis

The `packet_submission` table records which agent/harness/model produced each packet, but has NO FK to belief/evidence/edge/task. The `Result` struct (app.go:753-759) returns entity IDs in-memory but they are never persisted as a link back to the packet.

The relationship is transient during Persist():
```
pkt.Beliefs → INSERT belief → result.BeliefIDs map (in-memory only)
pkt.Evidence → INSERT evidence → (no ID tracking)
pkt.Edges → INSERT belief_edge → (no ID tracking)
pkt.Tasks → INSERT conductor_task → (no ID tracking)
```

The `Result.BeliefIDs` map (local→canonical) exists only during the transaction and is discarded after commit.

### Recommendation

**SMALL ADDITIVE CHANGE** for all four object types:
- Add `origin_packet_id UUID` column to `belief`, `evidence`, and `conductor_task` tables
- For `belief_edge`: add an additive provenance table `edge_provenance(edge_parent_id, edge_child_id, origin_packet_id)` since belief_edge has no scenario_id and a minimal schema
- Populate during Persist() using the packet_id already available in the method

---

## 2. PACKET SUBMISSION IMMUTABILITY

### Verified

- ONE packet_id maps to ONE record: `packet_id STRING PRIMARY KEY` — PK enforces uniqueness
- Cannot be silently updated: no UPDATE path exists in application code
- Cannot be deleted through normal paths: no DELETE path exists in application code
- Idempotency: `ON CONFLICT (packet_id) DO NOTHING` — no duplicate provenance

### Normal lifecycle objects

- Beliefs: INSERT with `ON CONFLICT (id) DO NOTHING` — never deleted, may be retracted (status change only)
- Evidence: INSERT with `ON CONFLICT (id) DO NOTHING` — never deleted
- Edges: INSERT with `ON CONFLICT (parent_id, child_id) DO NOTHING` — never deleted
- Tasks: INSERT with `ON CONFLICT (id) DO NOTHING` — status transitions only, never deleted

### Destructive commands

- `argus reset` drops all tables and re-applies migrations — development only
- No normal research lifecycle command deletes objects

### Verdict: PASS

No normal-path deletion risk found. Objects are preserved and may be superseded/retracted.

---

## 3. ADVERSARIAL TARGET VALIDATION

### What the current system permits

| Target type | Allowed? | Mechanism |
|-------------|----------|-----------|
| canonical existing belief (same scenario) | YES | resolveRef strips prefix, FK validates existence |
| local belief created in same packet | YES | resolveRef finds it in result.BeliefIDs map |
| belief in another scenario | YES | No scenario_id check on belief_edge |
| nonexistent belief | NO (DB-level) | REFERENCES belief(id) FK constraint |
| retracted belief | YES | No status check |
| promoted belief | YES | No status check |

### Current validation gaps

1. **Same-packet self-challenge**: ALLOWED. An adversarial packet can create contradicts edges between its own new beliefs. No check prevents this.
2. **Cross-scenario targeting**: ALLOWED. No scenario_id check on edges.
3. **Self-referencing edge**: ALLOWED. No CHECK (parent_id <> child_id) constraint.
4. **MCP path bypasses validateEdges()**: The MCP adapter calls app.Compile() + app.Validate(), NOT packetv1.Validate(). The validateEdges() checks (local_id present, kind valid) are skipped entirely in the agent submission path.

### What should be rejected

Based on the existing architecture:
- **Same-packet self-challenge SHOULD be rejected**: The adversarial rule says "challenge existing work" — work from the same packet is not "existing"
- **Cross-scenario targeting SHOULD be rejected**: Edges should be within a scenario
- **Local refs for contradicts SHOULD be rejected**: A contradiction target should be an existing canonical belief, not a new belief created in the same submission

### Recommendation

**SMALL ADDITIVE CHANGE**: Add validation in Persist() or a new validateEdgeTargets() step:
1. For contradicts edges: reject local: references (target must be canonical)
2. For all edges: verify parent_id and child_id belong to the same scenario as the packet
3. Optionally: reject self-edges (parent_id = child_id)

---

## 4. ADVERSARIAL REVIEW COVERAGE

### Current state

NONE of the three protocols require recording a disposition for every visible belief:
- `argus-agent-protocol.md` — no coverage requirement
- `opencode-adversarial.md` — no coverage requirement
- `opencode-work.md` — no coverage requirement

The adversarial rule says "challenge existing work" but does not require addressing every belief.

### Recommendation

**SMALL PROTOCOL CHANGE**: Add to `argus-agent-protocol.md`:

```
## Adversarial coverage rule

For every belief visible through get_context, the adversarial agent must
record a disposition in its final report:

- challenged (with reason)
- no material objection (with reason)
- insufficient basis to assess (with reason)

Do not silently skip beliefs.
```

This is a protocol-only change, no persistence needed yet.

---

## 5. EBP CATEGORY-ERROR CALIBRATION

### Current state

The shared protocol does NOT prevent this category error. There is no rule distinguishing:
- unresolved debt
- EBP violation
- process violation
- false claim

### Recommendation

**SMALL PROTOCOL CHANGE**: Add to `argus-agent-protocol.md`:

```
## EBP calibration rule

Do not label an unresolved obligation an EBP violation unless a specific
active EBP or Domain Pack rule is actually violated.

Open debt alone is not a violation when the claim is merely entered.
```

---

## 6. TASK SUBMISSION INTEGRITY

### Current state

Tasks ARE persisted to `conductor_task` via Persist() (app.go:464-477). However:

1. **project_id is hardcoded to empty string `''`** — this will fail at the FK constraint level (project_id UUID NOT NULL REFERENCES conductor_project(id))
2. **No origin_packet_id** — no link back to the creating packet
3. **Dashboard fetches ALL tasks globally** — no scenario/project filtering
4. **Trust UI shows tasks but no provenance** — no packet origin, no agent identity

### Why the adversarial tasks may not appear

The INSERT uses `project_id = ''` which violates the FK constraint. The task INSERT likely fails silently (the error is not propagated to the agent). This is an ingestion-integrity bug.

### Recommendation

**SMALL FIX**: 
1. Fix the project_id to use pkt.ScenarioID instead of empty string
2. Add origin_packet_id to conductor_task
3. Filter dashboard tasks by project_id

---

## 7. HUMAN TRACEABILITY

### Current state

| Question | Answer | Gap |
|----------|--------|-----|
| What is the claim? | YES — beliefs table | — |
| Who/what agent created it? | NO — no origin_packet_id on belief | MISSING |
| Which packet created it? | NO — no origin_packet_id on belief | MISSING |
| What evidence accompanied it? | PARTIAL — evidence has belief_id but no packet origin | PARTIAL |
| What beliefs does it derive from? | YES — edges now visible in get_context | — |
| What beliefs contradict it? | YES — edges now visible in get_context | — |
| Which adversarial packet challenged it? | NO — no origin on edges | MISSING |
| What tasks resulted? | NO — tasks have no packet origin | MISSING |

### Missing links

The human CANNOT trace:
- belief → packet → agent identity (missing origin_packet_id)
- edge → packet (missing origin on edges)
- task → packet (missing origin_packet_id)

### Recommendation

Add origin_packet_id to belief, evidence, and conductor_task tables. Add edge_provenance table. Update Trust UI to show origin information.

---

## 8. DECLARED VS VERIFIED AGENT IDENTITY

### Current state

Agent identity (agent_id, harness, model) is declared in the packet by the operator/launch prompt. No code verifies these values cryptographically or at runtime.

The packet_submission table stores these values as-is from the packet.

### Verdict

> Agent identity is declared provenance, not attested runtime identity.

This is acceptable for the POC. No change needed.

---

## 9. CONTEXT SNAPSHOT / REPLAYABILITY

### Current state

`get_context` returns a point-in-time reconstruction. No snapshot ID, version, or hash is recorded. The `content_sha256` in packet_submission is a hash of the packet content, not the context reviewed.

### Classification: MISSING

### Impact on current dry run

Not significant for single-agent sequential runs. The adversarial agent sees the state at the time of its get_context call.

### Recommendation

DEFER — record as named research-system debt.

---

## 10. FINAL ASSESSMENT

### A. Provenance spine

```
packet_submission
      |
      +--> belief?      MISSING (no origin_packet_id)
      +--> evidence?    MISSING (no origin_packet_id)
      +--> edge?        MISSING (no origin on edges)
      +--> task?        MISSING (no origin_packet_id, broken project_id)
```

All four relationships are missing.

### B. Adversarial integrity

| Check | Prevented? |
|-------|------------|
| attacking nonexistent beliefs | YES (DB FK) |
| cross-scenario targeting | NO |
| same-packet self-challenge | NO |
| ambiguous target references | PARTIAL (local refs resolve, canonical refs not validated) |

### C. Task integrity

The two proposed adversarial tasks likely failed to persist because:
- `project_id = ''` violates the FK constraint on conductor_task
- The INSERT fails silently (error not propagated)
- Tasks were in the packet but never materialized in the database

This is an ingestion-integrity bug: the agent submitted tasks, validation accepted them, but persistence failed.

### D. Protocol quality

Remaining ambiguities:
1. No coverage requirement for adversarial review
2. No EBP calibration rule (debt ≠ violation)
3. No explicit rule that contradicts targets should be canonical, not local

### E. UI lineage

A human CAN see:
- Beliefs with claim, status, debt
- Edges between beliefs
- Agent submissions (packet, agent, harness, model, role, time)

A human CANNOT see:
- Which packet created a belief
- Which agent created a belief (via packet)
- Which packet created an edge
- Which packet created a task
- Task origin/provenance

### F. Migration decision

| Change | Decision | Size |
|--------|----------|------|
| origin_packet_id on belief | SMALL ADDITIVE CHANGE | migration + Persist() update |
| origin_packet_id on evidence | SMALL ADDITIVE CHANGE | migration + Persist() update |
| origin_packet_id on conductor_task | SMALL ADDITIVE CHANGE | migration + Persist() update |
| edge_provenance table | SMALL ADDITIVE CHANGE | new migration + Persist() update |
| Fix task project_id bug | SMALL FIX | Persist() line 468 |
| Adversarial coverage protocol | SMALL PROTOCOL CHANGE | argus-agent-protocol.md |
| EBP calibration protocol | SMALL PROTOCOL CHANGE | argus-agent-protocol.md |
| Edge target validation | SMALL ADDITIVE CHANGE | Persist() validation |
| Context snapshot | DEFER | — |

---

## Implementation Plan

### Phase 1: Provenance spine (migration + Persist)

1. Create migration `004_provenance_spine.sql`:
   - `ALTER TABLE belief ADD COLUMN origin_packet_id STRING`
   - `ALTER TABLE evidence ADD COLUMN origin_packet_id STRING`
   - `ALTER TABLE conductor_task ADD COLUMN origin_packet_id STRING`
   - `CREATE TABLE edge_provenance (parent_id UUID, child_id UUID, origin_packet_id STRING, PRIMARY KEY(parent_id, child_id))`

2. Update `internal/application/app.go` Persist():
   - Add origin_packet_id to belief INSERT
   - Add origin_packet_id to evidence INSERT
   - Add origin_packet_id to conductor_task INSERT
   - Insert into edge_provenance after edge INSERT
   - Fix task project_id to use pkt.ScenarioID

### Phase 2: Edge target validation

3. Add validation in Persist() before edge INSERT:
   - For contradicts edges: reject local: references
   - For all edges: verify parent/child beliefs are in the same scenario

### Phase 3: Protocol updates

4. Update `prompts/argus-agent-protocol.md`:
   - Add adversarial coverage rule (disposition per belief)
   - Add EBP calibration rule (debt ≠ violation)

### Phase 4: UI provenance

5. Update `internal/epistemic/view.go`:
   - Add OriginPacketID to BeliefView
   - Query origin_packet_id in GetSnapshot

6. Update `internal/ui/templates/insights.html`:
   - Add Origin column to beliefs table
   - Add Origin column to tasks table

### Phase 5: Tests

7. Add tests for:
   - Provenance persistence (belief → packet link)
   - Edge target validation (reject local contradicts, reject cross-scenario)
   - Task project_id fix
   - Protocol coverage rules

### Phase 6: Verify

8. Run go vet, go test, verify compilation

---

## Estimated size

- ~50 lines of migration SQL
- ~30 lines of Persist() changes
- ~20 lines of edge validation
- ~10 lines of protocol text
- ~15 lines of UI template changes
- ~100 lines of tests

Total: ~225 lines across ~8 files
