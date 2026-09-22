# ARGUS - AI Agent Operating Contract

## 1. What ARGUS Is

ARGUS is a single-binary proof-of-concept trust-verification architecture for agent-driven research.

Central invariant:

    CAPABILITY != WORK != AUTHORITY != EXECUTION

An AI agent may perform substantial research work without being allowed to declare that work authoritative.

One ARGUS codebase / binary, one CockroachDB, one integrated local runtime.

## 2. Deployment Model

### argus serve

Local runtime + database bootstrap + Trust UI.

    ./argus serve
    # CockroachDB auto-started, DB created, migrations applied, seed inserted
    # Trust UI on http://localhost:8080
    # Health: GET /health

### argus mcp

MCP stdio adapter for an agent harness (separate invocation).

    ./argus mcp
    # JSON-RPC 2.0 over stdio
    # Exactly 2 tools: argus.get_context, argus.submit_packet

Same binary, separate process modes as required.

### Other subcommands

    ./argus verify --input <file>   Run verification
    ./argus migrate                  Apply migrations
    ./argus reset                    Drop all tables and re-apply

## 3. Architecture Invariants (non-negotiable)

- Agents produce work; they do not own authority
- Humans discharge epistemic debt
- Solvent: epistemic authority (subsystem, not separate service)
- Conductor: operational task state (subsystem, not separate service)
- EBP v2.1 governs debt and promotion
- MCP constrains agent capability by tool-surface, not prompt convention

## 4. Repository Structure

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
    domain-pack/           Domain-specific packs (bmist@1.1.0 current, bmist@1.0.0 historical)
    seed/                  POC seed data
    verifier/              Verification artifact registry
    docs/corpus/           Static research background (7 files, NOT a database)
    prompts/               Agent role cards (opencode-work.md, opencode-adversarial.md)

## 5. AI Agent Interface

### Transport

    argus mcp  (JSON-RPC 2.0 over stdio)

### Tools (exactly 2)

#### argus.get_context

Read-only. Returns the RCP/v1 context for a task.

Input:

    {
      "task_id": "<UUID>"    // required - UUID of the Conductor task
    }

Output:

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

Availability metadata: UNKNOWN != EMPTY. If a section is unavailable, the agent sees { "available": false, "reason": "..." } rather than an empty list.

#### argus.submit_packet

Write. Validates, compiles, and persists an EBP research packet.

Input (required fields):

    {
      "schema_version": "ebp-research-packet/v1",   // required
      "role": "work" | "adversarial",                // required
      "packet_id": "<string>",                        // required - unique per packet
      "pack_ref": "bmist@1.1.0",                     // required - domain pack reference
      "scenario_id": "<UUID>",                        // optional
      "beliefs": [                                    // required - array of beliefs
        {
          "local_id": "<string>",                     // unique within packet
          "claim": "<text>",                          // the claim
          "claim_type": "derived" | "accommodated" | "postulated",
          "debt": ["<obligation_key>", ...]           // optional - initial debt items
        }
      ],
      "evidence": [                                   // required - array of evidence
        {
          "local_id": "<string>",
          "belief_ref": "local:<belief_local_id>",
          "provenance_class": "external_feed" | "reproducible_artifact" | "live_scan" | "operator_asserted",
          "content_sha256": "<hex>",
          "source_url": "<url>"                       // optional
        }
      ],
      "edges": [                                      // optional - array of edges
        {
          "local_id": "<string>",
          "from_ref": "local:<id>" | "canonical:belief:<uuid>",
          "to_ref": "local:<id>" | "canonical:belief:<uuid>",
          "kind": "derives" | "contradicts"
        }
      ],
      "tasks": [                                      // optional - array of proposed tasks
        {
          "local_id": "<string>",
          "title": "<text>",
          "description": "<text>"
        }
      ]
    }

Reference prefixes:
- local:<id> = object inside the submitted packet
- canonical:belief:<uuid> = reference to an existing Solvent belief

Output:

    Persisted result with entity IDs (belief IDs, task IDs, etc.)

Validation behavior:
- Provenance class must be one of: external_feed, reproducible_artifact, live_scan, operator_asserted
- operator_asserted evidence requires human attestation; agent-submitted packets cannot create it
- Pack ref must resolve to a registered domain pack
- Returns error on validation failure, does not persist

Idempotency:
- Same content + same scenario deduplicates
- Verified by content hash

## 6. Capability Boundary

Agents do NOT receive:
- Direct database access
- Solvent mutation tools
- Conductor mutation tools
- Promotion tools
- Discharge tools
- Retraction tools
- Edge mutation tools
- Discharge/promote/retract HTTP endpoints

Agents can ONLY:
- Read context (argus.get_context)
- Submit work (argus.submit_packet)

Enforced by tool-surface configuration, not by prompt convention. If an agent cannot see a privileged tool, it cannot invoke it.

## 7. Where Authority Lives

- Epistemic authority: Solvent subsystem (internal/epistemic/)
- Operational state: Conductor subsystem (internal/work/)
- Orchestration: application layer (internal/application/)
- Human adjudication: Trust UI -> application -> Solvent
- Agent work: submission only (submit_packet), not authority

## 8. Agent Workflow

### Work Agent

    human gives initial task
        |
    agent calls argus.get_context
        |
    reconstructs current state from system (not conversation memory)
        |
    performs bounded research
        |
    produces EBP packet
        |
    agent calls argus.submit_packet

### Adversarial Agent

    separate fresh OpenCode process (NO conversational memory of Work Agent)
        |
    agent calls argus.get_context
        |
    reconstructs what has already been done
        |
    identifies unresolved work/debt
        |
    attacks the current work
        |
    produces adversarial EBP packet
        |
    agent calls argus.submit_packet

Agents reconstruct context from system state, not from inherited conversation memory. The adversarial agent intentionally challenges, not merely improves.

## 9. Static Research Corpus

docs/corpus/ contains 7 static research files:
- writeup_v1.md - context overview
- v6_1.md - base corpus
- v6_2_adv.md - adversarial corpus
- v6_3.md - latest version
- adv_review.md - first adversarial review
- adv_review2.md - second adversarial review
- adv_review3.md - third adversarial review

NOT a database, NOT vector search, NOT embeddings.

The current OpenCode workflow gives the agent repository filesystem access. ARGUS does not serve these documents through MCP. Agents read them directly from the filesystem.

For live system state, agents use argus.get_context, not the corpus.

## 10. Domain Pack

Current: bmist@1.1.0 (in domain-pack/bmist/v1.1.0/)
Historical: bmist@1.0.0 (in domain-pack/bmist/v1/) - immutable, not current

The pack defines:
- Debt vocabulary (obligation keys)
- Retirement rules (evidence class requirements)
- Claim types
- Faithfulness review requirements

## 11. Trust UI

- /ui/insights - research dashboard (current research scope)
- /ui/debts - debt obligations and discharge
- /health - health check

The Trust UI is a human control surface. It must never imply "AI verified debt" or "Agent retired debt". The agent may supply evidence; the human adjudicates; Solvent records the authoritative transition.

Default login token: argus-local-operator (override with ARGUS_OPERATOR_TOKEN)

## 12. Dry Run Status

- System integration: demonstrated (bootstrap, seed, get_context, submit_packet, UI)
- Actual AI research dry run: OUTSTANDING (next operational experiment)

## 13. Deferred Scope

- ADD_DEBT / automatic new-evidence-creates-debt loop
- Full multi-user authentication
- RCP graph traversal (broad scenario projection only)
- Live Work Agent + Adversarial Agent research run

## 14. Historical Planning

The following directories contain historical planning artifacts from earlier development phases. Do not treat them as current implementation instructions:

- .opencode/plans/ (28 files) - historical planning
- docs/plan8/ (50 files) - historical planning
- docs/dryrun/ (9 files) - historical planning
- plan/ (38 files) - phase implementation plans
- background/ (30 files) - research and design documents
- docs/archive/reference-loop/ - legacy integration test harness (archived, not active)

The current code is the authority. Documentation follows code, never the reverse.
