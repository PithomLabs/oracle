# Phase 8 Implementation Plan — ARGUS Trust Verification POC

## 1. Phase 8 Objective

Phase 8 is a **thin first dry run** of the full ARGUS trust-verification loop. It exercises the complete architecture against one real BM–IST research slice (Gate G0), proving that:

1. A human creates/selects a research task in Conductor.
2. A Work OpenCode agent reconstructs canonical context via RCP (MCP → Coordinator → Conductor + Solvent).
3. The Work agent performs bounded research and submits an EBP packet.
4. The Coordinator persists claims, evidence, edges, and tasks through existing Conductor/Solvent REST clients.
5. A fresh Adversarial OpenCode agent (no conversation history) reconstructs the same canonical state via RCP.
6. The Adversarial agent attacks unresolved work without rediscovering completed/dead-end work.
7. The Trust UI shows the evolving research state (Insights) and human debt-discharge surface (Debts).
8. A human discharges debt through the UI; the Coordinator validates the applicable Domain Pack retirement rule before invoking Solvent.
9. Solvent authoritatively records debt retirement.
10. A human requests promotion; Solvent — not the UI or agent — decides whether the gate passes.

**What this proves:** An autonomous agent can do substantial research work, another autonomous agent can attack that work, and a human can adjudicate the resulting epistemic state through a thin control surface — without either agent being able to declare its own work authoritative.

## 2. Current Architecture and Boundaries

### Frozen architectural boundaries

```
CAPABILITY ≠ WORK ≠ AUTHORITY ≠ EXECUTION

Agent / OpenCode    = agency / work
Conductor           = operational coordination and task state
Solvent             = epistemic authority, belief/evidence/debt ledger,
                      authorization and consequential promotion gates
Trust UI            = human control surface
EBP v2.1            = epistemic operating discipline
RCP                 = thin read-only context protocol for agents
```

### EBP v2.1 doctrine (mandatory)

- Ideas enter free.
- Promotion costs debt.
- Debt does not kill.
- Debt is forever payable.
- New evidence creates new debt.
- No final-truth claim may be promoted.
- Accounting must never become the work.
- Human beings discharge epistemic debt. Agents do NOT directly retire debt or promote claims.

### Data ownership

| Owner | State |
|-------|-------|
| Conductor | project, task, dependency, activity, task status, governance_ref, priority |
| Solvent | belief, belief_edge, evidence, debt, action_intent, audit_activity, promotion/retraction/authorization lifecycle |

No new persistence layer is created. RCP is a read-only projection, not a store.

### Protocol roles

| Protocol | Role |
|----------|------|
| EBP | Epistemic operating law |
| RCP | Thin agent-facing read protocol over Conductor + Solvent |
| MCP | Agent-facing transport/adapter |
| HTTP/API | Canonical integration contract |
| Trust UI | Human observation + adjudication |

## 3. Existing Interfaces Discovered in Conductor/Solvent

### Solvent REST API (30+ endpoints under `/v1/`)

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/v1/beliefs` | Enter a belief (ideas enter free) |
| GET | `/v1/beliefs/{id}` | Get belief |
| GET | `/v1/beliefs` | List beliefs (filter by status/claimType) |
| POST | `/v1/beliefs/{id}/debt/retire` | Retire one debt item |
| POST | `/v1/beliefs/{id}/promote` | Promote (DB refuses if debt remains) |
| POST | `/v1/beliefs/{id}/retract` | Retract with cascade |
| GET | `/v1/beliefs/{id}/explain` | Explain promotion/authorization readiness |
| GET | `/v1/beliefs/{id}/evidence` | List evidence for a belief |
| POST | `/v1/evidence` | Add evidence to a belief |
| GET | `/v1/evidence/{id}` | Get evidence |
| POST | `/v1/discharge` | Record debt discharge with attribution |
| GET | `/v1/activity` | Read audit activity entries |
| GET | `/v1/ledger` | Ledger summary (6 aggregate counts) |

**Key types:**
- `EnterBeliefRequest{ScenarioID, Claim, ClaimType, Debt []string}`
- `AddEvidenceRequest{ScenarioID, BeliefID, ProvenanceClass, SourceURL, ContentSHA256}`
- `RetireDebtRequest{DebtItem string}`
- `DischargeRequest{ScenarioID, BeliefID, ObligationKey, InstrumentRef, DischargedBy}`

**Solvent MCP tools (18 total):** 3 read (ledger, explain, activity), 5 mutation (retire_debt, promote, falsify, discharge, ingest_evidence), 8 authority lifecycle (create_principal, revoke_principal, create_target, attach_justification, request_authorization, approve, authorize, revoke_target), 2 execution (authorize_action, execute).

**Promotion gate:** DB CHECK constraint `promoted_is_debt_free` — `array_length(debt,1) > 0` OR `final_truth = true` → promotion refused (SQLSTATE 23514).

### Conductor REST API (19 endpoints under `/v1/`)

| Method | Path | Purpose |
|--------|------|---------|
| GET/POST | `/v1/projects` | List/create projects |
| GET/PUT/DELETE | `/v1/projects/{id}` | CRUD project |
| GET/POST | `/v1/projects/{id}/tasks` | List/create tasks |
| GET/PATCH | `/v1/tasks/{id}` | Get/update task (title/description/priority only) |
| POST | `/v1/tasks/{id}/claim` | Claim (proposed → active) |
| POST | `/v1/tasks/{id}/submit` | Submit (active → review) |
| POST | `/v1/tasks/{id}/accept` | Accept (review → accepted) |
| POST | `/v1/tasks/{id}/reject` | Reject (review → active) |
| POST | `/v1/tasks/{id}/blocker` | Report blocker (active → blocked) |
| POST | `/v1/tasks/{id}/resolve` | Resolve blocker (blocked → active) |
| POST | `/v1/tasks/{id}/release` | Release (active → proposed) |
| POST | `/v1/tasks/{id}/activity` | Post activity |
| GET | `/v1/tasks/{id}/governance` | Get governance state (read-only) |

**Key types:**
- `Task{ID, ProjectID, Title, Description, Status, Priority, CurrentAgent, GovernanceRef, CreatedAt, UpdatedAt}`
- `TaskUpdateFields{Title, Description, Priority}` — does NOT include status/current_agent/governance_ref
- `Activity{ID, TaskID, ActorType, ActorID, Action, Details, CreatedAt}`
- `Dependency{ID, TaskID, BlockedByID, CreatedAt}`

**Task state machine:**
```
proposed → active (claim) → review (submit) → accepted (accept) [terminal]
                ↓ blocked (report blocker) → active (resolve)
                ↓ cancelled [terminal]
active → proposed (release)
```

**Priority:** `low | medium | high | critical`, default `medium`. Stored but not used in ordering queries. POC-scoped manual assignment.

**Governance ref:** `TEXT` column on task. Write-once at creation, immutable after. JSON structure: `{provider, reference_id, metadata:{belief_id}}`.

### Oracle Coordinator

**Existing HTTP endpoints:**
| Method | Path | Purpose |
|--------|------|---------|
| POST | `/decisions` | Submit human decision (promote/retire_debt/retract/reopen/authorize/refuse) |
| GET | `/packets/{id}/status` | Get packet compilation status |
| GET | `/beliefs/{id}/decision-context` | Get decision context for a belief |
| GET | `/authorization-context/{id}` | Get Sphinx authorization projection |

**Existing methods:**
- `CompilePacket(pkt)` — validates, computes canonical hash, idempotency check, creates in-memory IDs
- `SubmitDecision(req)` — dispatches to Solvent via REST clients
- `PacketStatus(packetID)` — returns compilation status
- `DecisionContext(beliefID)` — returns decision context
- `AuthorizationContext(targetID)` — returns Sphinx projection

**Solvent client methods:** `CreateBelief`, `CreateEdge`, `CreateEvidence`, `PromoteBelief`, `RetireDebt`, `RetractBelief`, `ApproveTarget`, `GetBelief`

**Conductor client methods:** `CreateTask`

**Critical gap:** `CompilePacket` creates in-memory IDs but does NOT persist to Conductor/Solvent. Phase 8 must close this gap.

### BM-IST Domain Pack

```json
{
  "pack_id": "bmist",
  "version": "1.0.0",
  "debt_vocabulary": [
    "needMap", "needInvariant", "needToyCheck",
    "needNullModel", "needObstruction", "needFaithfulnessReview"
  ],
  "initial_debt": ["needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"],
  "evidence_classes": ["reproducible_artifact", "operator_asserted"],
  "retirement_rules": {
    "needMap": {"evidence_class": "reproducible_artifact", "rule": "map_check"},
    "needInvariant": {"evidence_class": "reproducible_artifact", "rule": "invariant_check"},
    "needToyCheck": {"evidence_class": "reproducible_artifact", "rule": "toy_model_check"},
    "needNullModel": {"evidence_class": "operator_asserted", "rule": "scope_clarification"},
    "needObstruction": {"evidence_class": "reproducible_artifact", "rule": "obstruction_construction"},
    "needFaithfulnessReview": {"evidence_class": "operator_asserted", "rule": "faithfulness_review"}
  }
}
```

### Gate G0 (dry-run subject)

From BM–IST Synthesis v5.1 §3:

- **L1 (Orbit Rigidity):** A strongly continuous unitary group cannot act nontrivially on a countable discrete/totally disconnected defined-state set. [THEOREM-LEVEL]
- **L2 (One-Parameter Triviality):** A finite signed-permutation operator group cannot itself support a nontrivial continuous one-parameter unitary flow. [THEOREM-LEVEL]
- **Gate G0:** countable substrate + literal continuous unitary evolution → structural incompatibility.
- **Required consequence:** discrete substrate time must be primitive; continuous Schrödinger evolution must be emergent.
- **Failure condition:** any construction that directly places literal continuous Schrödinger dynamics on the countable substrate fails G0.

Claims for dry run: G0 (main), L1, L2 (supporting). Dependencies: L1→G0, L2→G0.

## 4. Minimal Changes Required in Solvent/Conductor

### Solvent — No schema changes

The existing Solvent schema and API are sufficient for Phase 8:
- `EnterBelief` with initial debt → beliefs enter free with pack initial debt
- `AddEvidence` → agents produce evidence
- `RetireDebt` → human debt discharge
- `Promote` → Solvent enforces promotion gate
- `RetractCascade` → retraction with derives edge
- `Discharge` → attributed debt retirement with audit trail
- `AuditActivity` → full audit history

**One extension needed:** The Solvent client in the Coordinator (`coordinator/client.go`) currently lacks `ListBeliefs`, `ListEvidence`, `ListEdges`, `ListActivities`, and `GetLedger` methods. These must be added for the RCP context projection.

### Conductor — No schema changes

The existing Conductor schema and API are sufficient:
- `CreateTask` with governance_ref → task linked to Solvent belief
- `Claim`/`Submit`/`Accept`/`Reject` → task lifecycle
- `ListByProject`/`ListByStatus` → task queries for Insights
- `ListByTask`/`ListByProject` → activity queries
- `ListByTask`/`ListByBlockedBy` → dependency queries
- `priority` field → Insights attention ordering

**One extension needed:** The Conductor client in the Coordinator (`coordinator/client.go`) currently only has `CreateTask`. It needs `GetTask`, `ListTasks`, `ListDependencies`, `ListActivity`, and `GetGovernance` methods for the RCP projection.

### Summary of minimal changes

| Component | Change | Reason |
|-----------|--------|--------|
| `coordinator/client.go` | Add Solvent read methods: `ListBeliefs`, `ListEvidence`, `ListEdges`, `ListActivities`, `GetLedger` | RCP needs to read Solvent state |
| `coordinator/client.go` | Add Conductor read methods: `GetTask`, `ListTasks`, `ListDependencies`, `ListActivity`, `GetGovernance` | RCP needs to read Conductor state |
| `coordinator/compiler.go` | Add persistence stage: after validation, persist beliefs/evidence/edges/tasks through REST clients | submit_packet must actually persist |
| No Solvent schema changes | — | Existing schema is sufficient |
| No Conductor schema changes | — | Existing schema is sufficient |

## 5. RCP Protocol Design

### Protocol identity

```json
{
  "protocol": "RCP/v1"
}
```

### Canonical API contract

```http
GET /v1/context/{task_id}
```

### Conceptual response

```json
{
  "protocol": "RCP/v1",
  "task": {
    "id": "...",
    "project_id": "...",
    "title": "...",
    "description": "...",
    "status": "active",
    "priority": "high",
    "governance_ref": {
      "provider": "solvent",
      "reference_id": "track1",
      "metadata": {"belief_id": "..."}
    },
    "current_agent": "...",
    "created_at": "...",
    "updated_at": "..."
  },
  "dependencies": [
    {
      "task_id": "...",
      "blocked_by_id": "...",
      "blocked_by_status": "accepted",
      "source": "conductor"
    }
  ],
  "epistemic": {
    "beliefs": [
      {
        "id": "...",
        "claim": "...",
        "claim_type": "derived",
        "status": "entered",
        "debt": ["needMap", "needInvariant"],
        "final_truth": false,
        "source": "solvent"
      }
    ],
    "evidence": [
      {
        "id": "...",
        "belief_id": "...",
        "provenance_class": "reproducible_artifact",
        "source_url": "...",
        "content_sha256": "...",
        "artifact_ref": "...",
        "source": "solvent"
      }
    ],
    "edges": [
      {
        "parent_id": "...",
        "child_id": "...",
        "kind": "derives",
        "source": "solvent"
      }
    ],
    "debt": [
      {
        "belief_id": "...",
        "items": ["needMap", "needInvariant"],
        "source": "solvent"
      }
    ],
    "intents": [
      {
        "belief_id": "...",
        "action": "...",
        "state": "live",
        "source": "solvent"
      }
    ]
  },
  "activity": [
    {
      "source": "conductor",
      "type": "task.claimed",
      "actor_id": "...",
      "details": "...",
      "at": "..."
    },
    {
      "source": "solvent",
      "type": "belief_entered",
      "actor_id": "...",
      "details": "...",
      "at": "..."
    }
  ]
}
```

### Design constraints

1. **v1 scope:** Full scenario/task projection. No sophisticated task-anchored graph closure.
2. **Cross-ledger consistency:** Eventually consistent. Each fact tags its source (`conductor` or `solvent`). Gaps are projection lag, not synthesized state.
3. **Source attribution:** Every returned object identifies whether it came from Conductor or Solvent.
4. **Artifact resolution:** Evidence includes `artifact_ref` where available; hash-only is insufficient for fresh agents.
5. **No prescription:** RCP returns context. It does not prescribe scientific next steps. "Next smallest useful move" belongs to the Domain Pack / agent reasoning / UI guidance.
6. **No new schema:** RCP objects are projections of existing Conductor + Solvent objects, not new database entities.

### Inclusion algorithm (v1)

RCP v1 returns the **full scenario projection** for the given task's scenario. The task-to-epistemic relationship is broad: all beliefs, evidence, edges, debt, and intents in the task's scenario are returned. This avoids inventing a new graph traversal algorithm.

The task's `governance_ref` metadata provides the direct link to the governing belief. Dependencies show upstream task status. Activity merges Conductor and Solvent audit trails.

## 6. Coordinator RCP Endpoint Design

### New endpoint

```http
GET /v1/context/{task_id}
```

### Implementation

```go
// In coordinator/http/handler.go — new case in ServeHTTP:
case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/context/"):
    h.handleGetContext(w, r)

// New handler method:
func (h *Handler) handleGetContext(w http.ResponseWriter, r *http.Request) {
    taskID := strings.TrimPrefix(r.URL.Path, "/context/")
    if taskID == "" {
        writeError(w, "task_id required", http.StatusBadRequest)
        return
    }
    ctx, err := h.coordinator.GetContext(r.Context(), taskID)
    if err != nil {
        writeError(w, err.Error(), http.StatusInternalServerError)
        return
    }
    writeJSON(w, ctx, http.StatusOK)
}
```

### Coordinator method

```go
// In coordinator/context.go (new file):
func (c *Coordinator) GetContext(ctx context.Context, taskID string) (*RCPContext, error) {
    // 1. Get task from Conductor
    task, err := c.conductorClient.GetTask(ctx, taskID)
    if err != nil {
        return nil, fmt.Errorf("conductor: task not found: %w", err)
    }

    // 2. Get dependencies from Conductor
    deps, _ := c.conductorClient.ListDependencies(ctx, taskID)

    // 3. Get governance state if present
    var governanceRef *GovernanceReference
    if task.GovernanceRef != nil {
        governanceRef = parseGovernanceRef(*task.GovernanceRef)
    }

    // 4. Get epistemic state from Solvent (by scenario from governance_ref)
    scenarioID := ""
    beliefID := ""
    if governanceRef != nil {
        scenarioID = governanceRef.ReferenceID
        beliefID = governanceRef.Metadata["belief_id"]
    }

    var beliefs, evidence, edges, intents []interface{}
    var debt []DebtSummary
    if scenarioID != "" {
        beliefs, _ = c.solventClient.ListBeliefs(ctx, scenarioID)
        evidence, _ = c.solventClient.ListEvidence(ctx, scenarioID)
        edges, _ = c.solventClient.ListEdges(ctx, scenarioID)
        intents, _ = c.solventClient.ListIntents(ctx, scenarioID)
        // Debt is derived from beliefs
        for _, b := range beliefs {
            if len(b.Debt) > 0 {
                debt = append(debt, DebtSummary{BeliefID: b.ID, Items: b.Debt, Source: "solvent"})
            }
        }
    }

    // 5. Get activity from both sources
    conductorActivity, _ := c.conductorClient.ListActivity(ctx, taskID)
    solventActivity, _ := c.solventClient.ListActivities(ctx, scenarioID, 50)
    activity := mergeActivity(conductorActivity, solventActivity)

    return &RCPContext{
        Protocol:     "RCP/v1",
        Task:         task,
        Dependencies: deps,
        Epistemic: EpistemicBundle{
            Beliefs:  beliefs,
            Evidence: evidence,
            Edges:    edges,
            Debt:     debt,
            Intents:  intents,
        },
        Activity: activity,
    }, nil
}
```

### New types

```go
// coordinator/context.go
type RCPContext struct {
    Protocol     string            `json:"protocol"`
    Task         *ConductorTask    `json:"task"`
    Dependencies []DependencyView  `json:"dependencies"`
    Epistemic    EpistemicBundle   `json:"epistemic"`
    Activity     []ActivityView    `json:"activity"`
}

type ConductorTask struct {
    ID            string              `json:"id"`
    ProjectID     string              `json:"project_id"`
    Title         string              `json:"title"`
    Description   string              `json:"description"`
    Status        string              `json:"status"`
    Priority      string              `json:"priority"`
    GovernanceRef *GovernanceReference `json:"governance_ref,omitempty"`
    CurrentAgent  *string             `json:"current_agent,omitempty"`
    CreatedAt     string              `json:"created_at"`
    UpdatedAt     string              `json:"updated_at"`
}

type DependencyView struct {
    TaskID           string `json:"task_id"`
    BlockedByID      string `json:"blocked_by_id"`
    BlockedByStatus  string `json:"blocked_by_status"`
    Source           string `json:"source"`
}

type EpistemicBundle struct {
    Beliefs  []BeliefView  `json:"beliefs"`
    Evidence []EvidenceView `json:"evidence"`
    Edges    []EdgeView    `json:"edges"`
    Debt     []DebtSummary `json:"debt"`
    Intents  []IntentView  `json:"intents"`
}

type BeliefView struct {
    ID         string   `json:"id"`
    Claim      string   `json:"claim"`
    ClaimType  string   `json:"claim_type"`
    Status     string   `json:"status"`
    Debt       []string `json:"debt"`
    FinalTruth bool     `json:"final_truth"`
    Source     string   `json:"source"`
}

type EvidenceView struct {
    ID              string `json:"id"`
    BeliefID        string `json:"belief_id"`
    ProvenanceClass string `json:"provenance_class"`
    SourceURL       string `json:"source_url"`
    ContentSHA256   string `json:"content_sha256"`
    ArtifactRef     string `json:"artifact_ref,omitempty"`
    Source          string `json:"source"`
}

type EdgeView struct {
    ParentID string `json:"parent_id"`
    ChildID  string `json:"child_id"`
    Kind     string `json:"kind"`
    Source   string `json:"source"`
}

type DebtSummary struct {
    BeliefID string   `json:"belief_id"`
    Items    []string `json:"items"`
    Source   string   `json:"source"`
}

type IntentView struct {
    BeliefID string `json:"belief_id"`
    Action   string `json:"action"`
    State    string `json:"state"`
    Source   string `json:"source"`
}

type ActivityView struct {
    Source  string      `json:"source"`
    Type    string      `json:"type"`
    ActorID string      `json:"actor_id"`
    Details interface{} `json:"details,omitempty"`
    At      string      `json:"at"`
}

type GovernanceReference struct {
    Provider    string                 `json:"provider"`
    ReferenceID string                 `json:"reference_id"`
    Metadata    map[string]interface{} `json:"metadata"`
}
```

### New HTTP types

```go
// coordinator/http/types.go — additions:
type ContextResponse = coordinator.RCPContext  // reuse directly
```

## 7. MCP Adapter/Tool Surface

### Architecture

```
OpenCode
   ↓ MCP (stdio)
ARGUS MCP adapter
   ↓ HTTP
Coordinator
   ↓
Conductor + Solvent
```

### MCP tools exposed

| Tool | Maps to | Description |
|------|---------|-------------|
| `argus.get_context` | `GET /v1/context/{task_id}` | Read canonical context |
| `argus.submit_packet` | `POST /v1/packets` | Submit EBP packet |

### MCP tools NOT exposed (capability boundary)

| Tool | Source | Reason absent |
|------|--------|--------------|
| `solvent_retire_debt` | Solvent MCP | Agent must not retire debt directly |
| `solvent_promote` | Solvent MCP | Agent must not promote directly |
| `solvent_falsify` | Solvent MCP | Agent must not retract directly |
| `solvent_discharge` | Solvent MCP | Agent must not discharge directly |
| `solvent_create_principal` | Solvent MCP | Agent must not create principals |
| `solvent_create_target` | Solvent MCP | Agent must not create authority targets |
| `solvent_approve` | Solvent MCP | Agent must not approve targets |
| `solvent_authorize_action` | Solvent MCP | Agent must not authorize actions |
| `solvent_execute` | Solvent MCP | Agent must not execute |
| `conductor_create_task` | Conductor MCP | Agent must not create tasks |
| `conductor_claim_task` | Conductor MCP | Agent must not claim tasks |
| `conductor_submit_task` | Conductor MCP | Agent must not submit tasks |
| `conductor_accept_task` | Conductor MCP | Agent must not accept tasks |
| `conductor_reject_task` | Conductor MCP | Agent must not reject tasks |
| All Conductor write tools | Conductor MCP | Agent must not write Conductor |

### MCP adapter implementation

New file: `cmd/argus-mcp/main.go`

```go
// argus-mcp is the MCP server for OpenCode agents.
// It exposes exactly two tools: get_context and submit_packet.
// It does NOT expose Solvent mutation tools or Conductor write tools.
func main() {
    coordinatorURL := os.Getenv("ARGUS_COORDINATOR_URL")
    if coordinatorURL == "" {
        coordinatorURL = "http://localhost:9090"
    }

    server := mcp.NewServer(coordinatorURL)
    server.Run() // stdio JSON
}
```

### MCP tool schemas

```json
{
  "tools": [
    {
      "name": "argus.get_context",
      "description": "Get the canonical RCP context for a task. Returns task, dependencies, beliefs, evidence, edges, debt, intents, and activity from Conductor and Solvent.",
      "inputSchema": {
        "type": "object",
        "properties": {
          "task_id": {
            "type": "string",
            "description": "The Conductor task ID"
          }
        },
        "required": ["task_id"]
      }
    },
    {
      "name": "argus.submit_packet",
      "description": "Submit an EBP research packet. The Coordinator validates and persists beliefs, evidence, edges, and tasks through Conductor and Solvent.",
      "inputSchema": {
        "type": "object",
        "properties": {
          "packet": {
            "type": "object",
            "description": "The EBP packet following ebp-research-packet/v1 schema"
          }
        },
        "required": ["packet"]
      }
    }
  ]
}
```

### Agent OpenCode configuration

Each agent role gets its own `.opencode/` configuration that connects ONLY to the ARGUS MCP adapter:

```json
{
  "mcpServers": {
    "argus": {
      "command": "bin/argus-mcp",
      "args": [],
      "env": {
        "ARGUS_COORDINATOR_URL": "http://localhost:9090"
      }
    }
  }
}
```

**Critical:** The agent's MCP configuration does NOT include `solvent-mcp` or `conductor -mode mcp`. The capability boundary is enforced by configuration, not prompt instructions.

## 8. Work Agent Workflow

### Bootstrap

```text
1. Human creates a Conductor task for Gate G0
   - project: BM-IST Synthesis v5.1
   - title: "Formalize Gate G0"
   - description: "Turn L1, L2 and G0 into a precise research packet"
   - governance_ref: {"provider":"solvent","reference_id":"track1","metadata":{}}
   - priority: high

2. Human launches OpenCode as Work Agent
   - OpenCode reads: opencode.json → connects to argus MCP
   - Agent receives prompt: "Investigate Gate G0. Reconstruct context via argus.get_context."
```

### Workflow

```text
Work Agent
  ↓
argus.get_context(task_id)
  ↓
  Receives:
    task: {title: "Formalize Gate G0", status: "active", ...}
    dependencies: []
    epistemic: {beliefs: [], evidence: [], ...}
    activity: []
  ↓
  Agent reasoning:
    "No beliefs yet. I need to formulate G0, L1, L2 as claims."
  ↓
  Agent performs research:
    - Reads BM-IST v5.1 corpus artifact
    - Formalizes L1, L2, G0 claims
    - Produces evidence (proof sketches, citations)
    - Identifies dependencies (L1→G0, L2→G0)
    - Identifies open debt (from pack initial_debt)
  ↓
  Agent constructs EBP packet:
    {
      "schema_version": "ebp-research-packet/v1",
      "role": "work",
      "packet_id": "...",
      "pack_ref": "bmist-1.0.0",
      "task_ref": "local:task-g0",
      "agent": {"id": "work-agent", "role": "work"},
      "beliefs": [
        {"local_id": "L1", "claim": "Orbit Rigidity...", "claim_type": "derived"},
        {"local_id": "L2", "claim": "One-Parameter Triviality...", "claim_type": "derived"},
        {"local_id": "G0", "claim": "Countable substrate + continuous unitary → incompatibility", "claim_type": "derived", "debt": ["needMap","needInvariant","needToyCheck","needNullModel","needObstruction","needFaithfulnessReview"]}
      ],
      "edges": [
        {"from_ref": "local:L1", "to_ref": "local:G0", "kind": "derives"},
        {"from_ref": "local:L2", "to_ref": "local:G0", "kind": "derives"}
      ],
      "evidence": [...],
      "tasks": [
        {"local_id": "task-g0", "title": "Formalize Gate G0", "governance_ref": "{...}"}
      ]
    }
  ↓
  argus.submit_packet(packet)
  ↓
  Coordinator
    ↓ validates packet
    ↓ persists to Conductor + Solvent
    ↓ returns compilation result
  ↓
  Work Agent
    ↓
    "Packet submitted. G0, L1, L2 recorded. Debt items assigned."
```

## 9. Adversarial Agent Workflow

### Bootstrap

```text
1. Human launches a FRESH OpenCode process as Adversarial Agent
   - No conversational history
   - Same argus MCP configuration
   - Agent receives prompt:
     "You are the adversarial research agent.
      First obtain the ARGUS context for task {task_id} using argus.get_context.
      Treat the returned state as authoritative.
      Do not redo completed work.
      Attack unresolved claims.
      Submit findings as an EBP adversarial packet."
```

### Workflow

```text
Adversarial Agent (fresh process, no memory)
  ↓
argus.get_context(task_id)
  ↓
  Receives:
    task: {title: "Formalize Gate G0", status: "active"}
    epistemic: {
      beliefs: [L1(entered), L2(entered), G0(entered)],
      evidence: [...],
      edges: [L1→G0, L2→G0],
      debt: [G0: 6 items],
      intents: []
    }
    activity: [
      {source:"conductor", type:"task.claimed", ...},
      {source:"solvent", type:"belief_entered", ...},
      ...
    ]
  ↓
  Agent reasoning:
    "G0 is entered with 6 debt items. L1 and L2 support it.
     I must attack the assumptions and proof sketches."
  ↓
  Agent performs adversarial research:
    - Attacks L1: "Is the countable/discrete premise established?"
    - Attacks L2: "Does the operator group satisfy assumptions?"
    - Attacks G0: "Is 'continuous unitary' interpreted literally?"
    - Searches for counterexamples
    - Identifies missing assumptions
  ↓
  Agent constructs adversarial EBP packet:
    {
      "schema_version": "ebp-research-packet/v1",
      "role": "adversarial",
      "packet_id": "...",
      "pack_ref": "bmist-1.0.0",
      "task_ref": "local:task-g0",
      "agent": {"id": "adversarial-agent", "role": "adversarial"},
      "beliefs": [],
      "edges": [],
      "evidence": [
        {
          "local_id": "adv-evidence-1",
          "belief_ref": "canonical:belief:{G0-belief-id}",
          "provenance_class": "operator_asserted",
          "content_sha256": "...",
          "source_url": "adversarial-analysis/missing-premise"
        }
      ],
      "tasks": []
    }
  ↓
  argus.submit_packet(packet)
  ↓
  Coordinator
    ↓ validates adversarial packet
    ↓ persists evidence to Solvent
    ↓ returns result
  ↓
  Adversarial Agent
    ↓
    "Adversarial review complete. Evidence attached to G0."
```

## 10. EBP Packet Lifecycle

### Packet flow

```text
Agent constructs EBP packet (local IDs, claims, evidence, edges)
  ↓
argus.submit_packet → Coordinator
  ↓
1. Validate packet schema (packet/v1/Validate)
2. Resolve pack_ref (domain-pack registry)
3. Validate debt membership against pack vocabulary
4. Validate evidence classes against pack
5. Compute canonical hash (idempotency)
6. For each belief:
   a. Call Solvent POST /v1/beliefs with initial debt from pack
   b. Record canonical belief_id mapping
7. For each edge:
   a. Call Solvent (belief edge creation) or record for batch
8. For each evidence:
   a. Call Solvent POST /v1/evidence
9. For each task:
   a. Call Conductor POST /v1/projects/{id}/tasks with governance_ref
10. Return CompilationResult with persisted IDs
```

### Idempotency

The existing `IdempotencyCache` handles duplicate submissions:
- Canonical hash excludes `packet_id`, `run_id`, `scenario_id`
- If same content submitted twice: second returns cached result
- States: New → InFlight → Completed/Failed

### Debt initialization

When a belief is entered via packet compilation:
1. Coordinator reads `initial_debt` from the BM-IST pack
2. Passes the full initial debt array to `SolventClient.CreateBelief`
3. Solvent records belief with `debt = initial_debt`
4. Agent must later produce evidence to retire individual items

### Role semantics

| Role | Allowed beliefs | Allowed evidence | Allowed edges | Purpose |
|------|----------------|-----------------|---------------|---------|
| `work` | New claims | Supporting evidence | derives edges | Constructive research |
| `adversarial` | No new claims | Challenge/counterexample evidence | contradicts edges | Attack existing work |

Adversarial packets may reference existing beliefs via `canonical:belief:{uuid}` refs. They produce evidence that challenges existing claims, potentially opening new debt or creating contradicts edges.

## 11. Human Debt-Discharge Workflow

### UI flow

```text
Human opens Trust UI → Debts screen
  ↓
Sees list of open debt items:
  G0: needMap, needInvariant, needToyCheck, needNullModel, needObstruction, needFaithfulnessReview
  ↓
Clicks on needObstruction:
  WHY THIS DEBT EXISTS
    Adversarial finding: "The current formalization has not yet demonstrated
    that the assumptions of the cited theorem apply to the actual substrate object."
  EVIDENCE
    A-017, E-031, E-044
  APPLICABLE RETIREMENT RULE
    evidence_class: reproducible_artifact
    rule: obstruction_construction
  PROPOSED DISCHARGE
    Formalize the missing assumption and provide a counterexample attempt.
  SOURCE
    OpenCode / adversarial run #2
  ↓
Human decides:
  [ RETIRE DEBT ]  [ KEEP OPEN ]
  ↓
If RETIRE DEBT:
  Trust UI → POST /decisions {type: "RETIRE_DEBT", belief_id, scenario_id, debt_item: "needObstruction"}
  ↓
  Coordinator
    ↓ validates retirement rule:
      - checks that the debt item exists in the pack
      - checks that offered evidence class matches pack rule
      (mechanical validation, NOT scientific judgment)
    ↓ calls Solvent POST /v1/beliefs/{id}/debt/retire
    ↓
  Solvent
    ↓ removes "needObstruction" from belief.debt array
    ↓ records audit activity
    ↓
  Trust UI
    ↓
    Shows: "Debt retired by operator."
```

### Language rules

The UI must NEVER show:
- "AI verified debt"
- "Agent retired debt"
- "Automatically discharged"

The UI must show:
- "Adversarial agent supplied evidence relevant to this debt."
- "Human decision required."
- "Debt retired by operator."

### Coordinator retirement rule validation

```go
func (c *Coordinator) handleRetireDebt(req DecisionRequest, record *DecisionRecord) (*DecisionRecord, error) {
    // 1. Resolve pack from governance_ref
    pack, err := c.resolvePack(req.BeliefID)
    if err != nil {
        record.Result = "refused"
        record.RefusalReason = fmt.Sprintf("pack resolution failed: %v", err)
        return record, nil
    }

    // 2. Check that debt item exists in pack vocabulary
    if !packHasDebt(pack, req.DebtItem) {
        record.Result = "refused"
        record.RefusalReason = fmt.Sprintf("unknown debt item %q", req.DebtItem)
        return record, nil
    }

    // 3. Mechanical retirement rule validation
    //    (Coordinator does NOT decide scientific correctness)
    //    For Phase 8: if the debt item has a retirement rule,
    //    the human's decision is sufficient evidence that the rule is satisfied.
    //    The Coordinator records the retirement through Solvent.

    // 4. Call Solvent retire debt
    err = c.solventClient.RetireDebt(req.BeliefID, req.DebtItem, req.ScenarioID)
    if err != nil {
        record.Result = "refused"
        record.RefusalReason = err.Error()
        return record, nil
    }

    record.Result = "executed"
    return record, nil
}
```

**Phase 8 decision:** The Coordinator mechanically validates that the debt item exists in the pack vocabulary. For Phase 8, the human's explicit RETIRE_DEBT action is treated as sufficient evidence that the retirement rule is satisfied. The Coordinator does NOT infer scientific correctness. This is recorded as an explicit Phase 8 implementation decision.

## 12. Solvent Promotion/Refusal Workflow

### Promotion flow

```text
Human clicks [ PROMOTE ] on G0
  ↓
Trust UI → POST /decisions {type: "PROMOTE", belief_id: G0, scenario_id: track1}
  ↓
Coordinator
  ↓ calls Solvent POST /v1/beliefs/{id}/promote
  ↓
Solvent
  ↓ CHECK: array_length(debt,1) > 0 → SQLSTATE 23514
  ↓ CHECK: final_truth = true → SQLSTATE 23514
  ↓
  If debt remains:
    → Refused: "Promotion blocked: open debt or final-truth language"
    → Trust UI shows refusal
  ↓
  If debt empty AND final_truth=false:
    → Belief status changes to "promoted"
    → Audit activity recorded
    → Trust UI shows success
```

### Deliberate refusal demonstration

For the dry run, intentionally attempt promotion while debt exists:

```text
G0 has 6 debt items
Human clicks [ PROMOTE ]
Solvent refuses: "Promotion blocked: outstanding debt [needMap, needInvariant, ...]"
Trust UI shows:
  REFUSED
  Reason: Promotion gate failed.
  Outstanding debt: needMap, needInvariant, needToyCheck, needNullModel, needObstruction, needFaithfulnessReview
```

This is the moment ARGUS proves its reason for existing: **the workflow system might say "task completed" and the agent might say "proof completed," but Solvent says "you are not yet entitled to promote this claim."**

## 13. UI Architecture

### Technology pattern

Following the existing Solvent wizard pattern:
- Go `html/template` server-side rendering
- Vanilla HTML + inline CSS + vanilla JS
- No React, Vue, HTMX, or frontend framework
- JSON API for data operations
- Single-page-like navigation via JS

### Server location

New package: `oracle/ui/` or `oracle/cmd/argus-ui/`

### Server structure

```go
// ui/server.go
type Server struct {
    Addr           string
    CoordinatorURL string
    templates      *template.Template
    httpClient     *http.Client
}

func NewServer(addr, coordinatorURL string) *Server
func (s *Server) ListenAndServe() error
```

### Routes

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| GET | `/` | `handleDashboard` | Redirect to /insights |
| GET | `/insights` | `handleInsights` | Insights view |
| GET | `/debts` | `handleDebts` | Debts view |
| GET | `/debts/{belief_id}` | `handleDebtDetail` | Single debt detail |
| POST | `/api/decisions` | `handleAPIDecision` | Submit decision to Coordinator |
| GET | `/api/context/{task_id}` | `handleAPIContext` | Proxy to Coordinator RCP |
| GET | `/api/stats` | `handleAPIStats` | Aggregate statistics |

### Data flow

```
Trust UI → GET /api/context/{task_id} → Coordinator → Conductor + Solvent → response
Trust UI → POST /api/decisions → Coordinator → Solvent → response
```

Never: Trust UI → Solvent directly. Never: Trust UI → Conductor directly.

## 14. Insights View

### Layout

```
┌─────────────────────────────────────────────┐
│  ARGUS                                      │
│  ────────                                   │
│  Insights | Debts                           │
├─────────────────────────────────────────────┤
│                                             │
│  BM–IST Synthesis v5.1                      │
│                                             │
│  Open Claims        3                       │
│  Open Debt          6                       │
│  Active Tasks       2                       │
│  Adversarial Find.  1                       │
│  Promoted           0                       │
│  Refused            1                       │
│  Dead Ends          0                       │
│                                             │
│  ATTENTION                                  │
│  ──────────────                             │
│  1. [HIGH]   Formalize Gate G0              │
│  2. [MEDIUM] Scaling Audit                  │
│  3. [LOW]    Deferred cosmology branch      │
│                                             │
│  BELIEF GRAPH                               │
│  ──────────────                             │
│  L1 ────┐                                   │
│         ├──► G0 (entered, 6 debt)           │
│  L2 ────┘                                   │
│                                             │
│  RECENT ACTIVITY                            │
│  ──────────────                             │
│  [conductor] task.claimed — work-agent      │
│  [solvent]   belief_entered — G0            │
│  [solvent]   evidence_added — L1 proof      │
│  [solvent]   debt_retired — needMap         │
│                                             │
└─────────────────────────────────────────────┘
```

### Statistics derivation

| Statistic | Source | Query |
|-----------|--------|-------|
| Open Claims | Solvent | `GET /v1/beliefs?status=entered` count |
| Open Debt | Solvent | Sum of `len(debt)` across all entered beliefs |
| Active Tasks | Conductor | `GET /v1/tasks?status=active` count |
| Adversarial Findings | Solvent | Evidence count where `provenance_class` comes from adversarial role |
| Promoted | Solvent | `GET /v1/beliefs?status=promoted` count |
| Refused | Solvent | `GET /v1/activity?type=promotion_refused` count (from audit_activity) |
| Dead Ends | Conductor + Solvent | Structural: rejected tasks whose governance_ref points to retracted/contradicted beliefs |

### Attention ordering

```sql
-- Pseudo-query (Conductor):
SELECT * FROM conductor_task
WHERE project_id = ?
ORDER BY
  CASE priority
    WHEN 'critical' THEN 1
    WHEN 'high' THEN 2
    WHEN 'medium' THEN 3
    WHEN 'low' THEN 4
  END,
  created_at DESC
```

Priority is manually assigned in Conductor. This is POC-scoped. Documented as: "Priority is manually assigned and is POC-scoped. A mechanical priority derivation may replace this in a later phase."

### Dead-end derivation

```go
func isDeadEnd(task *ConductorTask, beliefs []BeliefView) bool {
    if task.Status != "cancelled" && task.Status != "rejected" {
        return false
    }
    if task.GovernanceRef == nil {
        return false
    }
    beliefID := task.GovernanceRef.Metadata["belief_id"]
    for _, b := range beliefs {
        if b.ID == beliefID && (b.Status == "retracted" || hasContradictedEdge(b, edges)) {
            return true
        }
    }
    return false
}
```

Dead end = rejected Conductor task whose governance_ref points to a retracted or contradicted belief. NOT merely "cancelled."

## 15. Debts View

### Layout

```
┌─────────────────────────────────────────────┐
│  ARGUS                                      │
│  ────────                                   │
│  Insights | Debts                           │
├─────────────────────────────────────────────┤
│                                             │
│  OPEN DEBT                                  │
│  ──────────                                 │
│                                             │
│  G0 — Discrete Substrate / Continuous       │
│       Unitary Compatibility                 │
│  Status: Unpromoted                         │
│                                             │
│  ○ needMap                                  │
│  ○ needInvariant                            │
│  ○ needToyCheck                             │
│  ○ needNullModel                            │
│  ○ needObstruction                          │
│  ○ needFaithfulnessReview                   │
│                                             │
│  ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─  │
│  Click on a debt item to expand:            │
│                                             │
│  WHY THIS DEBT EXISTS                       │
│  Adversarial finding: The current           │
│  formalization has not yet demonstrated      │
│  that the assumptions of the cited theorem  │
│  apply to the actual substrate object.      │
│                                             │
│  EVIDENCE                                   │
│  A-017 (reproducible_artifact)              │
│  E-031 (reproducible_artifact)              │
│  E-044 (operator_asserted)                  │
│                                             │
│  APPLICABLE RETIREMENT RULE                 │
│  evidence_class: reproducible_artifact      │
│  rule: obstruction_construction             │
│                                             │
│  PROPOSED DISCHARGE                         │
│  Formalize the missing assumption and       │
│  provide a counterexample attempt.          │
│                                             │
│  [ RETIRE DEBT ]     [ KEEP OPEN ]          │
│                                             │
└─────────────────────────────────────────────┘
```

### Data flow for debt detail

```text
1. GET /api/context/{task_id} → RCP context
2. For each belief with open debt:
   a. Find matching debt items
   b. Find evidence attached to that belief
   c. Find adversarial evidence (role=adversarial packets)
   d. Find retirement rule from pack
   e. Render debt card
3. Human clicks RETIRE DEBT:
   a. POST /api/decisions {type: "RETIRE_DEBT", belief_id, scenario_id, debt_item}
   b. Coordinator validates pack retirement rule
   c. Coordinator calls Solvent retire_debt
   d. UI refreshes
```

### Debt item expansion

Each debt item shows:
- **Why it exists:** Derived from the belief's context (adversarial findings, missing evidence)
- **Evidence:** All evidence attached to the belief (from Solvent)
- **Adversarial findings:** Evidence from adversarial packets
- **Retirement rule:** From BM-IST pack (evidence_class + rule)
- **Proposed discharge:** Agent-provided description of what would satisfy the debt
- **Prior history:** Activity entries showing prior attempts

## 16. Artifact/Evidence Retrieval

### Problem

Evidence hashes without a resolution path stall fresh agents at "verified, unobtainable."

### Solution

RCP evidence includes `artifact_ref` where available:

```json
{
  "id": "evidence-123",
  "belief_id": "belief-456",
  "provenance_class": "reproducible_artifact",
  "source_url": "file:///path/to/artifact.py",
  "content_sha256": "abc123...",
  "artifact_ref": "verifier:physics-v1:run-001"
}
```

### Resolution path

The `artifact_ref` maps to the existing `verifier.ArtifactReader` interface:

```go
type ArtifactReader interface {
    Resolve(ctx context.Context, ref string) (model.VerificationArtifact, error)
    VerifyHash(ref string, expectedSHA256 string) error
}
```

For Phase 8, artifact resolution is best-effort:
- If the artifact is in the verifier registry → resolved
- If the artifact is a file path → resolved via filesystem
- If neither → `artifact_ref` is present but unresolved; agent must work from the evidence metadata

### Artifact resolution endpoint (optional)

```http
GET /v1/artifacts/{ref}
```

Returns the `VerificationArtifact` if resolvable, 404 otherwise. This is a thin wrapper around `ArtifactReader.Resolve`.

## 17. Dead-End Derivation

### Structural definition

> Dead end = rejected Conductor task whose governance_ref points at a retracted or contradicted belief.

### Implementation

```go
func deriveDeadEnds(tasks []ConductorTask, beliefs []BeliefView, edges []EdgeView) []DeadEnd {
    var deadEnds []DeadEnd
    for _, task := range tasks {
        if task.Status != "cancelled" && task.Status != "rejected" {
            continue
        }
        if task.GovernanceRef == nil {
            continue
        }
        beliefID, ok := task.GovernanceRef.Metadata["belief_id"].(string)
        if !ok || beliefID == "" {
            continue
        }
        for _, b := range beliefs {
            if b.ID != beliefID {
                continue
            }
            if b.Status == "retracted" {
                deadEnds = append(deadEnds, DeadEnd{
                    TaskID:  task.ID,
                    BeliefID: beliefID,
                    Reason:  "belief retracted",
                    Task:    task,
                })
                break
            }
            // Check for contradicted edge
            for _, e := range edges {
                if e.ChildID == beliefID && e.Kind == "contradicts" {
                    deadEnds = append(deadEnds, DeadEnd{
                        TaskID:   task.ID,
                        BeliefID: beliefID,
                        Reason:   fmt.Sprintf("contradicted by belief %s", e.ParentID),
                        Task:     task,
                    })
                    break
                }
            }
        }
    }
    return deadEnds
}
```

### Constraints

- Do NOT use merely "cancelled" as dead-end semantic
- Do NOT invent a free-form `dead_end` field
- Do NOT trust agent-written dead-end prose
- Dead-end is derived from structural state: task status + governance_ref + belief status + edges

## 18. Cross-Ledger Consistency Handling

### Eventual consistency

RCP is explicitly **eventually consistent** across Conductor and Solvent. The two systems are independent stores with no shared transaction.

### Source tagging

Every returned fact identifies its source:

```json
{"source": "conductor"}
```
or:
```json
{"source": "solvent"}
```

### Projection lag

If Conductor records a task but Solvent has not yet received the corresponding belief, the RCP response shows:
- Task exists (source: conductor)
- Belief does not appear (source: solvent — empty)

This is projection lag, not an error. The agent sees the gap and can retry or wait.

### Response shape determinism

The RCP response shape is always the same structure, even when sections are empty:

```json
{
  "protocol": "RCP/v1",
  "task": {...},
  "dependencies": [],
  "epistemic": {
    "beliefs": [],
    "evidence": [],
    "edges": [],
    "debt": [],
    "intents": []
  },
  "activity": []
}
```

This ensures agents can always parse the response without conditional logic.

## 19. Security/Capability Boundary

### Principle

> The design must make the capability boundary enforceable, not merely a prompt instruction.

### Enforcement mechanism

**Tool-surface absence, not prompt instructions.**

The Work Agent and Adversarial Agent connect to the ARGUS MCP adapter, which exposes exactly two tools:
- `argus.get_context`
- `argus.submit_packet`

The Solvent MCP server (with 18 tools including `retire_debt`, `promote`, `approve`, etc.) and the Conductor MCP server (with 15 tools including `create_task`, `claim_task`, etc.) are **not configured** in the agent's MCP configuration.

### Configuration boundary

```json
// Work Agent: .opencode/config.json
{
  "mcpServers": {
    "argus": {
      "command": "bin/argus-mcp",
      "env": {"ARGUS_COORDINATOR_URL": "http://localhost:9090"}
    }
  }
  // NO solvent-mcp
  // NO conductor -mode mcp
}
```

```json
// Adversarial Agent: .opencode/config.json
{
  "mcpServers": {
    "argus": {
      "command": "bin/argus-mcp",
      "env": {"ARGUS_COORDINATOR_URL": "http://localhost:9090"}
    }
  }
  // NO solvent-mcp
  // NO conductor -mode mcp
}
```

### Negative test

Phase 8 must include a test proving:

> An agent attempting direct Solvent mutation cannot do so because the capability is absent from its tool surface.

```go
func TestAgentCannotMutateSolvent(t *testing.T) {
    // Create agent MCP server with only argus tools
    agentServer := argusmcp.NewServer(coordinatorURL)

    // Verify only argus tools are listed
    tools := agentServer.ListTools()
    require.Len(t, tools, 2)
    require.Equal(t, "argus.get_context", tools[0].Name)
    require.Equal(t, "argus.submit_packet", tools[1].Name)

    // Verify Solvent tools are NOT listed
    for _, tool := range tools {
        require.NotContains(t, tool.Name, "solvent_")
        require.NotContains(t, tool.Name, "conductor_")
    }

    // Verify attempt to call solvent tool fails
    _, err := agentServer.CallTool("solvent_retire_debt", map[string]interface{}{})
    require.Error(t, err)
    require.Contains(t, err.Error(), "unknown tool")
}
```

### What prevents bypass

| Attack vector | Defense |
|---------------|---------|
| Agent calls Solvent MCP directly | Solvent MCP not in agent's MCP config; tool not available |
| Agent calls Conductor MCP directly | Conductor MCP not in agent's MCP config; tool not available |
| Agent calls Solvent REST API | No HTTP client configured; agent only has MCP stdio |
| Agent calls Conductor REST API | No HTTP client configured; agent only has MCP stdio |
| Agent modifies its own MCP config | Config is read-only at startup; agent cannot reload |
| Prompt injection | Tool-surface absence is mechanical, not prompt-dependent |

## 20. Error/Refusal Semantics

### Packet validation errors

| Error | Source | Response |
|-------|--------|----------|
| Invalid schema version | packet/v1/Validate | 400 + error message |
| Unknown pack_ref | domain-pack registry | 400 + error message |
| Unknown debt item | ValidateDebtMembership | 400 + error message |
| Unsupported evidence class | ValidateEvidenceClass | 400 + error message |
| Bad reference format | packet/v1/Validate | 400 + error message |
| Duplicate local_id | packet/v1/Validate | 400 + error message |

### Decision errors

| Error | Source | Response |
|-------|--------|----------|
| Unknown debt item | Coordinator validate | refusal + "unknown debt item" |
| Debt already retired | Solvent RetireDeBT | "already discharged" (idempotent) |
| Promotion blocked (debt remains) | Solvent Promote | refusal + SQLSTATE 23514 |
| Promotion blocked (final_truth) | Solvent Promote | refusal + SQLSTATE 23514 |
| Belief not found | Solvent | 404 + error message |
| Invalid decision type | Coordinator | 400 + error message |

### RCP errors

| Error | Source | Response |
|-------|--------|----------|
| Task not found | Conductor | 404 + error message |
| Solvent unavailable | Solvent client | 503 + error message |
| Governance ref malformed | Coordinator parse | 400 + error message |

### UI error display

All errors are displayed as:
```html
<div class="verdict refuse">
  <span>REFUSED</span>
  <span class="d">{error detail from Coordinator/Solvent}</span>
</div>
```

Success is displayed as:
```html
<div class="verdict commit">
  <span>COMMITTED</span>
  <span class="d">{success detail}</span>
</div>
```

## 21. Observability/Audit Requirements

### Audit trail

Every consequential operation produces audit entries in Solvent's `audit_activity` table:

| Operation | Audit Type | Source |
|-----------|-----------|--------|
| Belief entered | `belief_entered` | Solvent |
| Evidence added | `evidence_added` | Solvent |
| Debt retired | `debt_retired` | Solvent |
| Belief promoted | `belief_promoted` | Solvent |
| Belief retracted | `belief_retracted` | Solvent |
| Promotion refused | (SQLSTATE logged) | Solvent |

Conductor records activity for task lifecycle:

| Operation | Activity Action | Source |
|-----------|----------------|--------|
| Task created | (implicit from create) | Conductor |
| Task claimed | `task.claimed` | Conductor |
| Task submitted | `task.submitted` | Conductor |
| Task accepted | `task.accepted` | Conductor |
| Task rejected | `task.rejected` | Conductor |
| Task blocked | `task.blocked` | Conductor |
| Task unblocked | `task.unblocked` | Conductor |

### RCP activity merge

RCP merges both activity streams, tagged by source, ordered by timestamp:

```json
[
  {"source": "conductor", "type": "task.claimed", "at": "2026-09-15T10:00:00Z"},
  {"source": "solvent", "type": "belief_entered", "at": "2026-09-15T10:01:00Z"},
  {"source": "solvent", "type": "evidence_added", "at": "2026-09-15T10:02:00Z"}
]
```

### Coordinator logging

The Coordinator logs:
- Packet compilation: packet_id, canonical_hash, result
- Decision submission: decision_type, belief_id, result, refusal_reason
- RCP context requests: task_id, response_size

### Solvent audit_activity schema

Already exists:
```sql
CREATE TABLE audit_activity (
    id UUID PRIMARY KEY,
    scenario_id UUID NOT NULL,
    type TEXT NOT NULL,
    actor_id TEXT NOT NULL,
    subject_id TEXT,
    details JSONB,
    sqlstate TEXT,
    constraint_name TEXT,
    refusal BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 22. Files/Modules to Create or Modify

### New files in oracle/

| File | Purpose |
|------|---------|
| `coordinator/context.go` | RCP context assembly (GetContext method) |
| `coordinator/context_test.go` | Tests for RCP context |
| `coordinator/client_conductor_read.go` | Conductor read methods: GetTask, ListTasks, ListDependencies, ListActivity, GetGovernance |
| `coordinator/client_solvent_read.go` | Solvent read methods: ListBeliefs, ListEvidence, ListEdges, ListActivities, GetLedger |
| `coordinator/compiler_persist.go` | Persistence stage for CompilePacket (persist through REST clients) |
| `coordinator/compiler_persist_test.go` | Tests for packet persistence |
| `coordinator/validate_retirement.go` | Retirement rule validation for human decisions |
| `coordinator/validate_retirement_test.go` | Tests for retirement validation |
| `cmd/argus-mcp/main.go` | ARGUS MCP adapter (get_context, submit_packet only) |
| `cmd/argus-mcp/mcp.go` | MCP tool implementation |
| `cmd/argus-mcp/mcp_test.go` | MCP tool tests |
| `cmd/argus-ui/main.go` | Trust UI server entry point |
| `ui/server.go` | Trust UI HTTP server |
| `ui/handlers.go` | Trust UI route handlers |
| `ui/templates/insights.html` | Insights view template |
| `ui/templates/debts.html` | Debts view template |
| `ui/templates/layout.html` | Shared layout template |
| `ui/static/` (if needed) | Minimal static assets |

### Modified files in oracle/

| File | Change |
|------|--------|
| `coordinator/http/handler.go` | Add `GET /v1/context/{task_id}` route |
| `coordinator/http/types.go` | Add RCP context response type |
| `coordinator/http/server.go` | No change needed |
| `coordinator/client.go` | Add read methods for Conductor and Solvent |
| `coordinator/mock.go` | Add mock methods for new client interfaces |
| `coordinator/compiler.go` | Add persistence stage after validation |
| `coordinator/human.go` | Add retirement rule validation |

### Modified files in solvent-main/

**None.** Solvent is frozen for Phase 8. All reads go through existing REST API.

### Modified files in conductor/

**None.** Conductor is frozen for Phase 8. All reads go through existing REST API.

## 23. Test Strategy

### Test levels

1. **Unit tests** — each new function/method in isolation
2. **Integration tests** — Coordinator against real Conductor + Solvent (or mocks)
3. **Contract tests** — HTTP API endpoint tests
4. **MCP tests** — argus-mcp tool surface tests
5. **Negative tests** — capability boundary enforcement
6. **End-to-end tests** — full dry-run workflow

### Test infrastructure

- Existing `MockSolventClient` and `MockConductorClient` in `coordinator/mock.go`
- Existing Conductor in-memory DB via `store.OpenTest()`
- Existing Solvent test setup via CockroachDB test containers or mocks
- Existing packet examples in `packet/v1/examples/`

### Test files

| Test File | Covers |
|-----------|--------|
| `coordinator/context_test.go` | RCP context assembly |
| `coordinator/compiler_persist_test.go` | Packet persistence |
| `coordinator/validate_retirement_test.go` | Retirement rule validation |
| `cmd/argus-mcp/mcp_test.go` | MCP tool surface |
| `ui/handlers_test.go` | Trust UI handlers |

## 24. Detailed Test Matrix

### RCP Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| R1 | Valid context retrieval | Known task_id with beliefs/evidence/debt | Full RCP response with all sections populated |
| R2 | Unknown task | Non-existent task_id | 404 error |
| R3 | No epistemic state yet | Task created but no beliefs entered | RCP with empty epistemic sections |
| R4 | Task with dependencies | Task blocked by another task | Dependencies section populated with upstream status |
| R5 | Task with existing beliefs | Task linked to beliefs with debt | Epistemic section populated with beliefs, debt |
| R6 | Source attribution | All objects present | Every object has `"source": "conductor"` or `"source": "solvent"` |
| R7 | Eventual consistency | Conductor has task, Solvent has no belief | Response shows task but empty epistemic |
| R8 | Artifact reference resolution | Evidence with artifact_ref | artifact_ref included in response |
| R9 | Deterministic response shape | Any valid task | Response always has protocol, task, dependencies, epistemic, activity |

### MCP Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| M1 | get_context exposed | List tools | argus.get_context present |
| M2 | submit_packet exposed | List tools | argus.submit_packet present |
| M3 | Solvent tools NOT exposed | List tools | No solvent_* tools |
| M4 | Conductor tools NOT exposed | List tools | No conductor_* tools |
| M5 | Negative: direct Solvent mutation | Call solvent_retire_debt | Error: unknown tool |
| M6 | Negative: direct Conductor mutation | Call conductor_create_task | Error: unknown tool |
| M7 | get_context returns valid RCP | Known task_id | Full RCP response |
| M8 | submit_packet accepts valid packet | Valid EBP packet | CompilationResult |

### Agent Reconstruction Tests

| # | Test | Setup | Expected |
|---|------|-------|----------|
| A1 | Fresh Work Agent context | No prior history | Receives full context via get_context |
| A2 | Fresh Adversarial Agent context | Work packet already submitted | Receives context with beliefs, evidence, debt |
| A3 | Completed work visible | Work agent completed | Adversarial agent sees completed work in context |
| A4 | Previous adversarial findings visible | Adversarial packet submitted | Next agent sees prior adversarial evidence |
| A5 | Dead-end work identifiable | Task cancelled + belief retracted | Context shows dead-end (structural derivation) |
| A6 | Unresolved debt visible | Belief with open debt | Context shows debt items |
| A7 | Artifact resolvable | Evidence with artifact_ref | Agent can resolve artifact via ref |

### EBP Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| E1 | Ideas enter without promotion | EnterBelief with initial_debt | Belief created at "entered" with debt |
| E2 | Initial debt appears | Belief created | All 6 pack initial_debt items present |
| E3 | Debt blocks promotion | Promote with open debt | SQLSTATE 23514, refusal |
| E4 | Human discharge retires debt | RetireDebt after human decision | Debt item removed from array |
| E5 | Agent cannot retire debt | Agent calls retire_debt | Tool not available |
| E6 | New evidence creates debt | Adversarial evidence added | New debt items if contradicts edge created |
| E7 | Final-truth claim cannot be promoted | Belief with final_truth=true | Promotion refused |
| E8 | Promotion Solvent-controlled | Promote with empty debt | Status changes to "promoted" |

### Retirement Rule Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| RR1 | Compliant evidence class succeeds | retire_debt with matching class | Debt retired |
| RR2 | Wrong evidence class refused | retire_debt with wrong class | Refusal |
| RR3 | Missing retirement rule handled | Debt item without rule | Coordinator handles explicitly |
| RR4 | Scientific correctness NOT inferred | Coordinator validates only mechanical rule | Coordinator does not judge science |

### REOPEN Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| RE1 | Retracted belief creates successor | REOPEN decision | New belief created, derives edge from retracted |
| RE2 | Derives lineage recorded | REOPEN | Edge from old belief to new belief |
| RE3 | Original remains in history | REOPEN | Retracted belief still in Solvent |

### Dead-End Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| DE1 | Structural derivation works | Cancelled task + retracted belief | Dead end detected |
| DE2 | Cancelled task alone is NOT dead end | Cancelled task + active belief | Not dead end |
| DE3 | Rejected task + contradicted belief | Rejected task + contradicts edge | Dead end detected |

### UI Tests

| # | Test | Input | Expected |
|---|------|-------|----------|
| U1 | Insights statistics render | Existing state | Statistics computed from Conductor + Solvent |
| U2 | Attention ordering uses Conductor priority | Tasks with different priorities | Ordered by priority, then created_at |
| U3 | Debts shows why debt exists | Belief with adversarial evidence | Evidence shown in debt detail |
| U4 | Debts shows evidence/adversarial findings | Belief with evidence | Evidence list shown |
| U5 | UI never presents "AI verified debt" | Any screen | No "AI verified" text |
| U6 | Human discharge invokes Coordinator path | RETIRE DEBT click | POST to /api/decisions → Coordinator → Solvent |
| U7 | UI does not write directly to Solvent | Any UI action | All writes go through Coordinator |

### Regression Tests

| # | Test | Expected |
|---|------|----------|
| RG1 | Existing coordinator tests | All pass |
| RG2 | Existing Solvent tests | All pass |
| RG3 | Existing Conductor tests | All pass |
| RG4 | No modifications to reference-loop | Unmodified |
| RG5 | No modifications to domain-pack validation | Unmodified |
| RG6 | No modifications to packet validation | Unmodified |
| RG7 | No modifications to verifier | Unmodified |
| RG8 | No modifications to corpus | Unmodified |

### Concurrency Tests

| # | Test | Expected |
|---|------|----------|
| C1 | RCP retrieval race-safe | Concurrent reads return consistent shape |
| C2 | Idempotency preserved | Duplicate submit_packet returns same result |
| C3 | Concurrent debt retirement | Second retirement of same item is idempotent |

## 25. Acceptance Criteria

### Architecture proof

1. Human creates a BM–IST task in Conductor with governance_ref.
2. Work OpenCode retrieves context through MCP → RCP → Conductor + Solvent.
3. Work agent produces an EBP packet and submits via argus.submit_packet.
4. Coordinator persists claims, evidence, edges through Solvent; tasks through Conductor.
5. Fresh adversarial OpenCode (no conversation history) retrieves the same canonical state through RCP.
6. Adversarial agent attacks unresolved work without rediscovering completed/dead-end work.
7. Insights shows the evolving research state and Conductor priorities.
8. Debts shows why each obligation remains open with evidence and retirement rules.
9. Human discharges at least one debt through the UI via Coordinator.
10. Solvent records the retirement authoritatively.
11. Promotion is attempted while debt remains; Solvent refuses.
12. All remaining debt retired; promotion succeeds; Solvent records the transition.
13. Complete audit trail reconstructible from Solvent + Conductor activity.

### Specific capability tests

14. Fresh adversarial agent can reconstruct full context with zero prior conversation.
15. Agent cannot directly mutate Solvent (tool-surface absence, not prompt instruction).
16. Agent cannot directly write to Conductor (tool-surface absence).
17. Dead-end tasks are structurally derivable (rejected task + retracted/contradicted belief).
18. RCP response is deterministic and includes source attribution.
19. Cross-ledger consistency gaps are visible as projection lag, not errors.
20. Debt items are atomic (present/retired) — no partial retirement.

### Non-regression

21. All existing coordinator tests remain green.
22. All existing Solvent tests remain green.
23. All existing Conductor tests remain green.
24. No modifications to unrelated packages unless explicitly justified.

## 26. Explicit Non-Goals / Deferred Work

### Not implemented in Phase 8

- New research database or generalized research graph database
- Task priority model in Solvent (uses existing Conductor field)
- Partial debt semantics (debt items are atomic)
- Automatic debt retirement from natural language
- Autonomous promotion
- Autonomous scientific adjudication
- General lesson/knowledge graph
- Full BM-IST research dashboard
- All mythology/navigation screens
- Direct OpenCode access to Solvent/Conductor writes
- Crash-recoverable persistent projection system
- Broad RCP graph traversal beyond v1 scope
- Sophisticated task-anchored graph closure
- RCP as a separate persistence layer
- New Solvent schema changes
- New Conductor schema changes
- React/Vue/HTMX frontend
- Multi-agent orchestration
- Automated theorem proving
- Automatic claim extraction
- Scientific judgment by LLM

### Deferred to later phases

- Full research-program navigation
- Multi-slice BM-IST exploration (beyond Gate G0)
- Empirical fitting
- Cosmology/dark-matter branch
- General RCP v2 with graph closure
- Mechanical priority derivation
- Real-time agent collaboration
- Production deployment

## 27. Risks and Remaining Decisions

### Risks

| Risk | Mitigation |
|------|------------|
| Coordinator client methods may not cover all needed reads | Inspect Solvent/Conductor REST APIs; add methods as needed |
| RCP assembly may be slow for large scenarios | Phase 8 is thin dry run; optimize later |
| MCP adapter may need OpenCode-specific configuration | Test with actual OpenCode process |
| Solvent REST client may lack batch operations | One-by-one persistence for Phase 8 |
| Cross-ledger consistency gaps during dry run | Document as expected behavior |
| Trust UI template rendering may need debugging | Start with minimal templates, iterate |

### Remaining decisions

1. **Trust UI hosting:** Same server as Coordinator, or separate? Recommendation: separate `argus-ui` binary, same pattern as Conductor's web mode.
2. **Scenario ID for dry run:** Use `track1` (existing hardcoded scenario) or create new? Recommendation: use `track1` for simplicity.
3. **Governance ref metadata:** What exactly goes in `metadata.belief_id`? The actual Solvent belief ID after packet persistence.
4. **Artifact resolution for Phase 8:** Best-effort only; agents work from evidence metadata if artifact unavailable.
5. **Adversarial packet evidence class:** Adversarial evidence uses `operator_asserted` (the agent's analysis is operator-asserted in the sense that the agent is the operator's tool).

---

**This plan will be reviewed before any Phase 8 coding begins.**
