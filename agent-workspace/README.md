# ARGUS

**Module:** `github.com/PithomLabs/oracle`
**EBP Version:** v2.1

ARGUS is a single-binary proof-of-concept trust-verification architecture for agent-driven research.

```
CAPABILITY != WORK != AUTHORITY != EXECUTION
```

Agents produce work. Solvent records authority. Humans discharge epistemic debt. No component conflates these roles.

---

## Quick Start

```bash
git clone <repo-url> && cd oracle
go build -o bin/argus ./cmd/argus/
./argus serve
```

Open http://localhost:8080 — Trust UI loads with seeded data.

Login token: `argus-local-operator`

---

## Prerequisites

- Go 1.25+
- `cockroach` binary in PATH (or at `~/.local/bin/cockroach`)

---

## Architecture

One ARGUS codebase / binary, one CockroachDB, one integrated local runtime.

### Deployment Model

```
argus serve
    = local runtime + DB bootstrap + Trust UI on http://localhost:8080

argus mcp
    = MCP stdio adapter for AI agent harness (separate invocation)
```

Same binary, separate process modes as required.

### Component Map

```
                    BM-IST Research Program
                             |
                             v
                      Domain Pack / EBP
                             |
                  +----------+----------+
                  v                     v
            Work Agent            Adversarial Agent
                  |                     |
                  +----------+----------+
                             v
                        argus mcp
                             |
                    Application Layer
                   /                 \
                  v                   v
             Conductor             Solvent
          operational state    epistemic authority
                  ^                   ^
                  |                   |
            argus get_context    argus get_context
```

### Human Control Path

```
Human
  v
Trust UI (http://localhost:8080)
  v
Application Layer
  v
Solvent (epistemic authority)
```

---

## CLI Commands

| Command | Purpose |
|---------|---------|
| `argus serve` | Start CockroachDB + DB + migrations + seed + Trust UI |
| `argus mcp` | Start MCP stdio adapter (for AI agents) |
| `argus verify` | Run verification |
| `argus migrate` | Apply migrations |
| `argus reset` | Drop all tables and re-apply |

---

## AI Agent Connection

```
argus mcp
```

JSON-RPC 2.0 over stdio. Exactly 2 tools.

### argus.get_context

Read-only. Returns the RCP/v1 context for a task.

**Input:**

```json
{
  "task_id": "<UUID>"
}
```

**Output:**

```json
{
  "task": { "id", "project_id", "title", "description", "status", "priority", "current_agent", "governance_ref", "created_at", "updated_at" },
  "dependencies": [...],
  "snapshot": {
    "beliefs": [{ "id", "claim", "claim_type", "status", "debt", "final_truth" }],
    "evidence": [{ "belief_id", "source_url", "provenance_class", "content_sha256" }],
    "intents": [{ "belief_id", "action", "state" }],
    "audit_live_on_nonpromoted": <int>
  },
  "availability": {
    "task": { "available": true/false, "reason": "..." },
    "dependencies": { "available": true/false, "reason": "..." },
    "snapshot": { "available": true/false, "reason": "..." }
  }
}
```

### argus.submit_packet

Write. Validates, compiles, and persists an EBP research packet.

**Required fields:**

| Field | Type | Description |
|-------|------|-------------|
| `schema_version` | string | Must be `ebp-research-packet/v1` |
| `role` | string | `work` or `adversarial` |
| `packet_id` | string | Unique per packet |
| `pack_ref` | string | Domain pack reference (e.g. `bmist@1.1.0`) |
| `beliefs` | array | Array of belief objects |
| `evidence` | array | Array of evidence objects |

**Optional fields:** `scenario_id`, `edges`, `tasks`

**Belief object:** `{ "local_id", "claim", "claim_type", "debt" }`

**Evidence object:** `{ "local_id", "belief_ref", "provenance_class", "content_sha256", "source_url" }`

**Edge object:** `{ "local_id", "from_ref", "to_ref", "kind" }` where kind is `derives` or `contradicts`

**Reference prefixes:** `local:<id>` (intra-packet), `canonical:belief:<uuid>` (existing Solvent belief)

**Validation:** Returns error on failure, does not persist. `operator_asserted` evidence requires human attestation; agent-submitted packets cannot create it.

**Idempotency:** Same content + same scenario deduplicates via content hash.

---

## Agent Context Sources

The current AI-agent workflow uses two distinct context sources.

## Live ARGUS state

Use:

```text
argus.get_context
```

This is the source for current task, work, epistemic, evidence, debt, and availability state.

## Static research background

The repository contains seven background documents under:

```text
docs/corpus/
```

They are not:

- an authority store
- a database corpus
- vector search
- embeddings
- a second epistemic ledger
- an MCP search service

The current local OpenCode workflow gives the agent repository filesystem access, so the agent reads these files directly. ARGUS does not serve them through MCP.

This distinction is deliberate:

```text
Live state        -> argus.get_context
Static background -> docs/corpus/*
```

Do not add retrieval infrastructure merely to make these seven files discoverable.

# Canonical Agent Workflow

## Work Agent

```text
Fresh Work Agent
      |
      v
argus.get_context(task_id)
      |
      v
reconstruct current ARGUS state
      |
      v
read docs/corpus/*
      |
      v
perform bounded research
      |
      v
construct EBP packet
      |
      v
argus.submit_packet
      |
      v
persisted ARGUS state
```

## Adversarial Agent

A separate fresh agent process is used:

```text
Fresh Adversarial Agent
      |
      v
argus.get_context(task_id)
      |
      v
reconstruct current ARGUS state
      |
      v
read the same research background
      |
      v
challenge claims / evidence / debt / logic
      |
      v
construct adversarial EBP packet
      |
      v
argus.submit_packet
```

Agents reconstruct context from system state, not inherited conversation memory.

# Trust UI

Embedded in `argus serve` at http://localhost:8080

| Route | Purpose |
|-------|---------|
| `/ui` | Redirects to `/ui/insights` |
| `/ui/insights` | Research dashboard (current research scope) |
| `/ui/debts` | Debt obligations and discharge |
| `/ui/api/login` | POST — authenticate with operator token |
| `/ui/api/discharge` | POST — discharge debt (authenticated) |
| `/ui/api/promote` | POST — promote belief (authenticated) |
| `/health` | Health check |

The Trust UI is a human control surface. It must never imply "AI verified debt" or "Agent retired debt". The agent may supply evidence; the human adjudicates; Solvent records the authoritative transition.

---

## Configuration

| Setting | Default | Override |
|---------|---------|----------|
| DB URL | `postgres://root@localhost:26257/argus?sslmode=disable` | `--db <url>` or `ARGUS_DB_URL` |
| Listen | `:8080` | `--listen <addr>` |
| Operator token | `argus-local-operator` | `ARGUS_OPERATOR_TOKEN` |

---

## Testing

```bash
go test ./...
go vet ./...
```

---

## Project Structure

```
oracle/
  cmd/argus/             CLI entry point (serve, mcp, verify, migrate, reset)
  internal/
    application/         App layer (orchestrates everything)
    epistemic/           Solvent subsystem (beliefs, evidence, edges, debt)
    work/                Conductor subsystem (tasks, dependencies, lifecycle)
    ui/                  Trust UI (embedded HTML templates + HTTP handlers)
    mcp/                 MCP adapter (2 tools only)
    migrations/          Database migrations (work + idempotency)
    boundary/            Capability boundary enforcement
  coordinator/           Packet compiler + orchestration logic (Go library)
  packet/v1/             EBP packet types and validation
  domain-pack/           Domain packs (bmist@1.1.0 current, bmist@1.0.0 historical)
  seed/                  POC seed data
  verifier/              Verification artifact registry
  prompts/               Agent role cards (opencode-work.md, opencode-adversarial.md)
  docs/corpus/           Static research background (7 files, NOT a database)
  docs/archive/          Historical artifacts (reference-loop, old plans)
```

---

## What is NOT Implemented (Deferred)

| Item | Status |
|------|--------|
| `ADD_DEBT` / automatic new-evidence-creates-debt | Deferred |
| Full multi-user authentication | Deferred |
| RCP graph traversal | Broad scenario projection only |
| Actual AI research dry run | Outstanding |
| Vector search / embeddings | Not planned |

---

## EBP v2.1

> Ideas enter free.
> Promotion costs debt.
> Debt does not kill.
> Debt is forever payable.
> New evidence creates new debt.
> No final-truth claim may be promoted.
> Accounting must never become the work.

---

## Historical Planning

The following directories contain historical planning artifacts from earlier development phases. Do not treat them as current implementation instructions:

- `.opencode/plans/` 
- `docs/plan8/` 
- `docs/dryrun/` 
- `plan/` 
- `background/`
- `docs/archive/reference-loop/`

The current code is the authority. Documentation follows code, never the reverse.
