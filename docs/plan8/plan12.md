# Plan 12: ARGUS Documentation + Repository Cleanup

**Status:** APPROVED — ready for implementation
**Date:** 2026-09-20
**Scope:** Documentation correction + repository cleanup only. No new features.

---

## 0. Prerequisites

Before implementing, inspect the current source code at each step.
The current code is the authority. Documentation follows code, never the reverse.

---

## 1. File Classification

### DELETE

| File/Directory | Reason |
|----------------|--------|
| `trust-ui/` | Superseded by `internal/ui/`. Own `go.mod` is dead code. |
| `session-dryrun.md` | Session log of previous debugging work. No ongoing value. |
| `.gut` | Gut CLI profile metadata. Delete only if repository-local and confirmed unused by project workflow. |

### ARCHIVE (move to `docs/archive/`)

| Source | Destination | Reason |
|--------|-------------|--------|
| `reference-loop/` | `docs/archive/reference-loop/` | Legacy integration test harness referencing old multi-service architecture. Verify no active dependency first. |

### KEEP IN PLACE, MARK HISTORICAL (in documentation only)

| Directory | Reason |
|-----------|--------|
| `.opencode/plans/` (28 files) | Historical planning artifacts. AGENTS.md will note they are historical. |
| `docs/plan8/` (50 files) | Historical planning artifacts. AGENTS.md will note they are historical. |
| `docs/dryrun/` (9 files) | Historical planning artifacts. |
| `plan/` (38 files) | Phase implementation plans. Historical. |
| `background/` (30 files) | Research and design documents. Historical. |
| `docs/BM_IST_AS_v4/`, `docs/BM_IST_v3/` | Historical research documents. |
| `docs/deep_research/` | Historical research documents. |
| `docs/claims/` | Historical claims documents. |
| `docs/pivot/` | Historical pivot documents. |

### KEEP AS ACTIVE

| File/Directory | Reason |
|----------------|--------|
| `AGENTS.md` | Rewrite as AI agent operating contract. |
| `README.md` | Rewrite as human developer guide. |
| `cmd/argus/` | Current CLI implementation. |
| `internal/` | Current application code. |
| `mcp/adapter/` | Current MCP adapter. |
| `coordinator/` | Current coordinator library. |
| `packet/` | Current EBP packet types. |
| `domain-pack/` | Current domain pack registry. |
| `seed/` | Current seed logic. |
| `verifier/` | Current verification registry. |
| `prompts/` | Current agent prompts. |
| `docs/corpus/` | Current static research corpus (7 files). |
| `Taskfile.yml` | Current developer task runner. |
| `.golangci.yml` | Current linter config. |
| `integration_test.go` | Current integration test. |

---

## 2. Reference Check Before Deletion

### 2a. trust-ui/ reference check

Search the entire repository for:
- `trust-ui` in imports
- `trust-ui` in go.mod references
- Any import path containing `trust-ui`

Expected result: no active imports. The actual Trust UI is `internal/ui/`.

### 2b. reference-loop/ reference check

Search the entire repository for:
- `reference-loop` in imports
- `reference-loop` in go.mod references
- Any CI/GitHub Actions referencing reference-loop

Expected result: the reference-loop has its own `go.mod` (separate module) and is not imported by the main module.

### 2c. .gut reference check

Search the repository for:
- `.gut` referenced in Taskfile.yml, Makefile, CI configs, or any script
- Any tooling that reads `.gut`

Expected result: no active reference. It is repository-local Gut CLI metadata with no project workflow dependency.

### 2d. Verify after deletion/archival

```bash
go build ./...
go test ./...
go vet ./...
```

Then start `argus serve` and verify:
- `/ui/insights` renders (HTTP 200, seeded task + belief visible)
- `/ui/debts` renders (HTTP 200, seeded belief debt visible)
- `/health` returns OK

---

## 3. AGENTS.md — Full Rewrite

**Goal:** Rewrite as the AI agent operating contract.

Replace the entire file. The current content describes a multi-service architecture (Solvent as separate service, Conductor as separate service, Coordinator as orchestrator between services) that no longer reflects the single-binary implementation.

### New Structure

```
# ARGUS — AI Agent Operating Contract

## 1. What ARGUS Is
- Single-binary trust-verification POC for agent-driven research
- Central invariant: CAPABILITY != WORK != AUTHORITY != EXECUTION
- One ARGUS codebase / binary, one CockroachDB, one integrated local runtime

## 2. Deployment Model
- argus serve
    = local runtime + DB bootstrap + Trust UI on http://localhost:8080
- argus mcp
    = MCP stdio adapter for an agent harness (separate invocation)
- Same binary, separate process modes as required

## 3. Architecture Invariants (non-negotiable)
- Agents produce work; they do not own authority
- Humans discharge epistemic debt
- Solvent: epistemic authority (subsystem, not separate service)
- Conductor: operational task state (subsystem, not separate service)
- EBP v2.1 governs debt and promotion
- MCP constrains agent capability by tool-surface, not prompt convention

## 4. Repository Structure
(verified listing of actual directories)

## 5. Build / Test / Run
### Build
  go build -o bin/argus ./cmd/argus/

### Run
  ./argus serve
  # CockroachDB auto-started, DB created, migrations applied, seed inserted
  # Trust UI on http://localhost:8080

### Test
  go test ./...
  go vet ./...

### Verify
  ./argus verify --input <file>

### Reset
  ./argus reset

## 6. AI Agent Interface
### Transport
  argus mcp  (JSON-RPC 2.0 over stdio)

### Tools (exactly 2)
  argus.get_context
  argus.submit_packet

### argus.get_context
- Input: { "task_id": "<UUID>" }  (required)
- Output: task, dependencies, beliefs, evidence, edges, debt, intents, activity
- Read-only. Does not mutate state.
- Returns availability metadata per section (UNKNOWN != EMPTY)

### argus.submit_packet
- Input: EBP packet
  - schema_version: string (required)
  - role: "work" | "adversarial" (required)
  - packet_id: string (required)
  - pack_ref: string (required)
  - scenario_id: string (optional)
  - beliefs: array (required)
  - evidence: array (required)
  - edges: array (optional)
  - tasks: array (optional)
- Output: persisted result with entity IDs
- Write. Validates, compiles, persists to Solvent + Conductor.
- Idempotent: same content + scenario deduplicates.
- Validation failure returns error, does not persist.

## 7. Capability Boundary
Agents do NOT receive:
- Direct database access
- Solvent mutation tools
- Conductor mutation tools
- Promotion tools
- Discharge tools
- Retraction tools
- Edge mutation tools

Agents can ONLY:
- Read context (get_context)
- Submit work (submit_packet)

Enforced by tool-surface configuration, not by prompt convention.

## 8. Where Authority Lives
- Epistemic authority: Solvent subsystem (internal/epistemic/)
- Operational state: Conductor subsystem (internal/work/)
- Orchestration: application layer (internal/application/)
- Human adjudication: Trust UI -> application -> Solvent
- Agent work: submission only (submit_packet), not authority

## 9. Agent Workflow
- Fresh Work Agent: get_context, research, submit_packet
- Fresh Adversarial Agent: get_context, reconstruct, challenge, submit_packet
- Agents reconstruct from system state, not conversation memory
- The adversarial agent intentionally challenges, not merely improves

## 10. Static Research Corpus
- docs/corpus/ contains 7 static research files
- NOT a database, NOT vector search, NOT embeddings
- The current OpenCode workflow gives the agent repository filesystem access
- ARGUS does not serve these documents through MCP
- Background material only; agents use get_context for live system state

## 11. Domain Pack
- Current: bmist@1.1.0
- Historical: bmist@1.0.0 (immutable, not current)

## 12. Trust UI
- /ui/insights — research dashboard (current research scope)
- /ui/debts — debt obligations and discharge
- Human control surface, never "AI verified"

## 13. Dry Run Status
- System integration: demonstrated (bootstrap, seed, get_context, submit_packet, UI)
- Actual AI research dry run: OUTSTANDING (next operational experiment)

## 14. Deferred Scope
- ADD_DEBT / automatic new-evidence-creates-debt
- Full multi-user authentication
- RCP graph traversal
- Live Agent + Adversarial Agent research run

## 15. Historical Planning
- .opencode/plans/, docs/plan8/, docs/dryrun/, plan/, background/
- All contain historical planning artifacts
- Do not treat as current implementation instructions
```

---

## 4. README.md — Full Rewrite

**Goal:** Rewrite for a human developer/operator.

### New Structure

```
# ARGUS

## What is ARGUS?
(brief, 3-4 sentences)

## Quick Start
  go build -o bin/argus ./cmd/argus/
  ./argus serve
  open http://localhost:8080

## Prerequisites
- Go 1.25+
- cockroach binary in PATH

## Architecture
- One ARGUS codebase / binary, one CockroachDB, one integrated local runtime
- argus serve = local runtime + DB bootstrap + Trust UI
- argus mcp = MCP stdio adapter for AI agent harness (separate invocation)

## CLI Commands
| Command | Purpose |
|---------|---------|
| argus serve | Start CockroachDB + DB + migrations + seed + Trust UI |
| argus mcp | Start MCP stdio adapter (for AI agents) |
| argus verify | Run verification |
| argus migrate | Apply migrations |
| argus reset | Drop all tables and re-apply |

## AI Agent Connection
  argus mcp  (JSON-RPC 2.0 over stdio)
  Tools: argus.get_context, argus.submit_packet

### argus.get_context
- Input: { "task_id": "<UUID>" }
- Output: task, dependencies, beliefs, evidence, edges, debt, intents, activity
- Read-only

### argus.submit_packet
- Input: EBP packet (schema_version, role, packet_id, pack_ref, beliefs, evidence, edges, tasks)
- Output: persisted result with entity IDs
- Write, validates, idempotent

## Trust UI
- /ui/insights — research dashboard (current research scope)
- /ui/debts — debt obligations and discharge
- Default login token: argus-local-operator

## Configuration
| Setting | Default | Override |
|---------|---------|----------|
| DB URL | postgres://root@localhost:26257/argus | --db or ARGUS_DB_URL |
| Listen | :8080 | --listen |
| Token | argus-local-operator | ARGUS_OPERATOR_TOKEN |

## Testing
  go test ./...
  go vet ./...

## Project Structure
(verified listing)

## What is NOT Implemented (Deferred)
- ADD_DEBT automatic loop
- Full multi-user auth
- Vector search / embeddings
- Actual AI research dry run

## License / EBP
- EBP v2.1
```

Key differences from current README:
- Remove "Solvent, Conductor run as separate services"
- Remove "Coordinator standalone server binary" references
- Remove stale verification evidence table
- Remove Reference Loop section
- Remove old multi-service architecture descriptions
- Fix bmist@1.0.0 reference to bmist@1.1.0
- Lead with working commands, not historical architecture
- Explicit `argus serve` vs `argus mcp` distinction
- MCP schemas documented inline

---

## 5. Prompts — Verify and Keep

Both `prompts/opencode-work.md` and `prompts/opencode-adversarial.md` are current and correct.

Verify by inspecting:
- They reference `argus.get_context` and `argus.submit_packet` (correct)
- They reference `docs/corpus/` files (correct)
- They reference `bmist@1.1.0` (correct)
- They don't reference obsolete architecture

No changes needed unless inspection reveals stale instructions.

---

## 6. Plan Execution Order

| Step | Action | Verify |
|------|--------|--------|
| 1 | Reference check: grep for `trust-ui` across repo | No active imports found |
| 2 | Reference check: grep for `reference-loop` across repo | No active imports found |
| 3 | Reference check: grep for `.gut` references in project workflow | No active references found |
| 4 | Delete `trust-ui/` | `go build ./...` passes |
| 5 | Archive `reference-loop/` to `docs/archive/reference-loop/` | `go build ./...` passes |
| 6 | Delete `session-dryrun.md` | no impact |
| 7 | Delete `.gut` (if confirmed unused) | no impact |
| 8 | Run `go build ./... && go test ./... && go vet ./...` | all pass |
| 9 | Start `argus serve`, verify `/ui/insights`, `/ui/debts`, `/health` | all render correctly |
| 10 | Rewrite `AGENTS.md` | inspection confirms accuracy |
| 11 | Rewrite `README.md` | inspection confirms accuracy |
| 12 | Verify prompts are current | no stale instructions |
| 13 | Final: `go build ./... && go test ./... && go vet ./...` | all pass |

---

## 7. What This Does NOT Change

- No new features, services, databases, or MCP tools
- No architecture changes
- No agent prompt creation
- No changes to actual Go code (except deleting dead trust-ui/ module)
- No changes to database schemas or migrations
- No changes to domain pack content

---

## 8. Acceptance Criteria

| AC | Criterion |
|----|-----------|
| AC1 | AGENTS.md accurately describes the current single-binary architecture |
| AC2 | README.md accurately describes the current architecture |
| AC3 | Neither document claims separate Solvent/Conductor services |
| AC4 | `argus serve` documented as primary local entrypoint |
| AC5 | `argus mcp` documented as primary AI-agent transport |
| AC6 | Two MCP tools documented with full input/output schemas |
| AC7 | Agent capability boundary explicitly documented |
| AC8 | Agent work vs human authority distinction documented |
| AC9 | Static 7-file research corpus documented correctly |
| AC10 | bmist@1.1.0 current, 1.0.0 historical |
| AC11 | Trust UI routes documented correctly |
| AC12 | Developer/test commands verified against repository |
| AC13 | Live AI dry run described as outstanding, not completed |
| AC14 | Stale material removed/archived/marked historical |
| AC15 | No live code deleted |
| AC16 | No new features added |
| AC17 | A fresh AI coding agent can understand the current architecture, development workflow, capability boundaries, and agent interface from AGENTS.md without consulting historical planning documents (source code still needed for exact implementation details) |
| AC18 | Human developer can run ARGUS from README.md alone |

---

## 9. Final Report (to be produced at completion)

1. AGENTS.md changes
2. README.md changes
3. Files deleted (trust-ui/, session-dryrun.md, .gut)
4. Files archived (reference-loop/)
5. Files retained but marked historical
6. Stale architectural claims removed
7. Actual AI-agent interfaces documented
8. Tests run and results
9. Any remaining ambiguity
