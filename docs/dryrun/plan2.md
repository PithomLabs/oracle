# ARGUS BM-IST-AS v6.3 — Final Implementation Plan

## Locked Decisions

| Area | Decision |
|------|----------|
| Research corpus | Static repository files (docs/corpus/) |
| Corpus DB ingestion | **Removed** |
| Embeddings | **Removed** |
| Vector search | **Removed** |
| `get_context(query)` | **Removed** |
| MCP tools | Exactly 2: `argus.get_context`, `argus.submit_packet` |
| MCP transport | Dedicated `argus mcp` subcommand with stdio |
| Domain pack | **`bmist@1.1.0`** (new) |
| `bmist@1.0.0` | Immutable historical version, not modified |
| Seed | 1 project + 1 task + 1 belief + 8 debts |
| Debt source | Domain pack initial_debt, not hardcoded seed list |
| `governance_ref` | Resolve in Phase-0 before seed implementation |
| Provenance enforcement | Existing evidence/retirement machinery |
| Provenance test | One focused verification, no new subsystem |
| ADD_DEBT | Deferred |
| Solvent kernel | No changes |
| New service | None |
| New DB | None |

---

## 1. Audit — What Already Works

| # | Capability | Status | Evidence |
|---|-----------|--------|----------|
| 1 | Corpus storage (manifest) | READY | `corpus/manifest.go` |
| 2 | Corpus ingestion (manifest) | N/A | Removed from scope |
| 3 | Embedding generation | N/A | Removed from scope |
| 4 | Vector search | N/A | Removed from scope |
| 5 | Corpus provenance | N/A | Removed from scope |
| 6 | Corpus citations | N/A | Removed from scope |
| 7 | get_context | READY | `app.go:459-504` — returns task, deps, snapshot, availability |
| 8 | submit_packet | READY | `app.go:316-457` — compile, validate, persist |
| 9 | Insights | READY | `/ui/insights` renders tasks + beliefs |
| 10 | Debts | READY | `/ui/debts` renders debt status + promote/block |
| 11 | Retirement-rule enforcement | READY | Validate + SubmitDecision |
| 12 | Human discharge | READY | SubmitDecision("discharge") -> kernel.Discharge |
| 13 | Promotion gate | READY | kernel.Promote -> schema CHECK |
| 14 | Adversarial packet handling | READY | Same submit_packet path |
| 15 | Domain-pack loading | READY | registry.go loads embedded bmist@1.0.0 |
| 16 | Database bootstrap | READY | bootstrap.go starts CRDB, applies migrations |
| 17 | BM-IST pack | PARTIAL | 6 debt items; needs 8 |
| 18 | Developer startup | PARTIAL | No seed, no MCP stdio |

### Critical Gaps

1. **MCP stdio transport** — cmdMCP hangs with `select{}`
2. **First-run seed** — empty database
3. **Domain pack v6.3** — missing 2 debt classes
4. **Hardcoded debt validation** — app.go isValidDebtItem checks only 6 items
5. **governance_ref type mismatch** — SQL UUID vs coordinator JSON

---

## 2. Solvent Kernel Impact

**No kernel changes required.**

All gaps are in Oracle's application layer, MCP adapter, DB schema (additive), and domain pack data.

---

## 3. Architecture (unchanged)

```
OpenCode
   |
   | MCP stdio
   v
ARGUS — one Go binary
   |
   +-- application/orchestration
   +-- Solvent epistemic kernel
   +-- minimal work subsystem
   +-- domain pack (bmist@1.1.0)
   +-- MCP adapter (2 tools)
   +-- HTTP/Trust UI
   |
   v
CockroachDB
```

Static background files: `docs/corpus/*.md` — read by agents via repository filesystem access.

---

## 4. Phase 0 — Audit + governance_ref Resolution

### 4a. governance_ref (BLOCKING)

**Before writing seed code**, inspect:
- `conductor_task.governance_ref` SQL type (currently UUID)
- `App.Persist` insertion path (line ~434)
- `ConductorClient.CreateTask` serialization
- Coordinator's `SubmitPacket` governance_ref handling
- `SubmitDecision` retract path (cancelLinkedTasks)

Establish the canonical representation. Prefer existing application constructor over raw SQL. Make seed use the same path the application uses.

### 4b. Full audit

Classify all 22 items from the prompt's audit list. Record in audit results before any code changes.

---

## 5. Phase 1 — Domain Pack bmist@1.1.0

### 5a. New pack at bmist@1.1.0

Create `domain-pack/bmist/v1/pack.json` updated content (or new version directory if the registry requires it):

```json
{
  "pack_id": "bmist",
  "version": "1.1.0",
  "name": "BM-IST-AS Research Domain Pack",
  "description": "BM + IST + Asymptotic Safety research methodology",
  "claim_types": ["derived", "accommodated", "postulated"],
  "evidence_classes": ["reproducible_artifact", "operator_asserted"],
  "debt_vocabulary": [
    "needMap", "needInvariant", "needToyCheck",
    "needNullModel", "needObstruction", "needFaithfulnessReview",
    "needInitialCondition", "needRegularity"
  ],
  "initial_debt": [
    "needMap", "needInvariant", "needToyCheck",
    "needNullModel", "needObstruction", "needFaithfulnessReview",
    "needInitialCondition", "needRegularity"
  ],
  "retirement_rules": {
    "needMap": {"evidence_class": "reproducible_artifact", "rule": "map_check"},
    "needInvariant": {"evidence_class": "reproducible_artifact", "rule": "invariant_check"},
    "needToyCheck": {"evidence_class": "reproducible_artifact", "rule": "toy_model_check"},
    "needNullModel": {"evidence_class": "operator_asserted", "rule": "scope_clarification"},
    "needObstruction": {"evidence_class": "reproducible_artifact", "rule": "obstruction_construction"},
    "needFaithfulnessReview": {"evidence_class": "operator_asserted", "rule": "faithfulness_review"},
    "needInitialCondition": {"evidence_class": "operator_asserted", "rule": "preparation_accounting"},
    "needRegularity": {"evidence_class": "reproducible_artifact", "rule": "regularity_audit"}
  },
  "falsifiers": ["counterexample", "contradiction", "missing_evidence", "alternative_explanation"],
  "consequential_actions": [
    {"action": "publish_claim", "requires": "promoted", "gates": ["faithfulness_review"]}
  ],
  "human_gated_transitions": ["faithfulness_review", "scope_clarification", "obstruction_assessment"],
  "verifier_specs": [{"verifier_id": "physics-v1", "min_version": "0.1.0"}]
}
```

### 5b. Preserve bmist@1.0.0

Do NOT modify `domain-pack/bmist/v1/pack.json` content for 1.0.0.

The registry must support loading both versions. If the current registry only loads one version per pack_id, add the ability to register `bmist@1.1.0` alongside `bmist@1.0.0`.

### 5c. Fix hardcoded validation in app.go

Replace `isValidDebtItem()` (lines ~705-717) with pack registry lookup:

```go
func (a *App) isValidDebtItem(item string, packRef string) bool {
    packID, version := parsePackRef(packRef)
    pack := a.packRegistry.Get(packID, version)
    if pack == nil {
        return false
    }
    for _, v := range pack.GetDebtVocabulary() {
        if v == item {
            return true
        }
    }
    return false
}
```

Same for `isValidEvidenceClass()`.

### 5d. Update hardcoded pack references

Update `scenarioPackMapping()` in app.go to return `bmist@1.1.0`.

**Files changed:**
- `domain-pack/bmist/v1/` (new version directory or updated pack.json)
- `internal/application/app.go` (hardcoded validation + pack mapping)

**Acceptance:** Pack validates. All 8 debt items accepted. Existing tests updated to reference 1.1.0.

---

## 6. Phase 2 — MCP Stdio Transport

**Files changed:**
- `cmd/argus/main.go` (rewrite `cmdMCP`)

Implement JSON-RPC framing over stdin/stdout:

1. Read newline-delimited JSON-RPC from stdin
2. Route `initialize` -> return protocol version + capabilities
3. Route `notifications/initialized` -> acknowledge
4. Route `tools/list` -> adapter.ListTools()
5. Route `tools/call` -> adapter.HandleTool(name, arguments)
6. Write JSON-RPC responses to stdout
7. Logs/debugging to stderr ONLY (never stdout)

~100-150 lines. No MCP SDK dependency. Use `encoding/json` + `bufio.Scanner`.

Keep `argus mcp` as a separate subcommand. Do NOT merge with `argus serve` HTTP traffic.

**Acceptance:** `argus mcp` communicates over stdio. OpenCode can call both tools.

---

## 7. Phase 3 — First-Run Research Seed

**Files changed:**
- `seed/seed.go` (new)
- `cmd/argus/main.go` (call seed after migrations)

### Seed data

Insert only when the target scenario has zero beliefs (idempotent guard).

1. **Project:**
   ```
   name: "BM-IST-AS v6.3"
   description: "BM + IST + Asymptotic Safety research program"
   ```

2. **Root task:**
   ```
   title: "Can one explicit microscopic dynamical system generate the structures of quantum theory without assuming them?"
   description: "Central research question. Current baseline: v6.3. Status: research program, not established theory."
   status: active
   priority: high
   governance_ref: <belief_id> (resolved per Phase 0 findings)
   ```

3. **Incumbent hypothesis (belief):**
   ```
   claim: "A discrete arithmetic microscopic dynamical system is a candidate substrate from which structures associated with quantum theory may emerge without being assumed."
   claim_type: postulated
   debt: pack.GetInitialDebt()  // from bmist@1.1.0, all 8 items
   status: entered
   ```

   The claim passes the final-truth language check. It is a hypothesis under test.

### Seed implementation rules

- Use the existing application/persistence paths, not raw SQL (per governance_ref resolution)
- Debt comes from the loaded pack, not a hardcoded list
- Check if beliefs exist before seeding (idempotency)
- Wire into `cmdServe` after migrations and pack loading

**Acceptance:** Fresh `argus serve` shows one active task + one belief with 8 debts on /ui/insights. /ui/debts shows all 8 obligations.

---

## 8. Phase 4 — Role Cards

**New files:**
- `prompts/opencode-work.md`
- `prompts/opencode-adversarial.md`

### Work Agent

```markdown
# ARGUS Work Agent

You are a research agent working within the ARGUS trust-verification system.

## Your tools
- `argus.get_context` — Read current ARGUS research state
- `argus.submit_packet` — Submit research findings as an EBP packet

## You MUST NOT
- Write to the database directly
- Retire debt, promote claims, retract beliefs, or authorize actions
- Treat your own conclusions as authoritative
- Treat any document as authoritative Solvent evidence

## Workflow
1. Call `argus.get_context` with the task ID to reconstruct current ARGUS state
2. Read all seven research documents under docs/corpus/ as background context
3. Distinguish ARGUS state from background documents
4. Identify the current research frontier and open obligations
5. Decompose the central question into manageable candidate claims
6. Gather evidence, make logic chains explicit
7. Create derives relationships where appropriate
8. Submit an EBP packet with beliefs, evidence, edges, and tasks

## Background Documents
- v6_3.md — current research-program baseline (primary reference)
- v6_2_adv.md — superseded historical baseline
- v6_1.md — earlier historical baseline
- adv_review.md, adv_review2.md, adv_review3.md — adversarial/agent-derived material (not authoritative)
- writeup_v1.md — explanatory/supporting material

## EBP v2.1
- Ideas enter free
- Promotion costs debt
- Debt does not kill
- Debt is forever payable
- New evidence creates new debt
- Promotion never means truth
```

### Adversarial Agent

```markdown
# ARGUS Adversarial Agent

You are an adversarial review agent. You are a completely fresh process with no memory of any prior agent.

## Your tools
- `argus.get_context` — Read current ARGUS research state
- `argus.submit_packet` — Submit adversarial findings as an EBP packet

## You MUST NOT
- Write to the database directly
- Retire debt, promote claims, retract beliefs, or authorize actions
- Assume prior work is correct

## Workflow
1. Call `argus.get_context` to reconstruct ALL prior research state
2. Read the same seven research documents under docs/corpus/
3. Identify unresolved claims, open debt, missing evidence, hidden assumptions
4. Attack logic, evidence completeness, and debt coverage
5. Submit an adversarial packet with contradicts edges and/or new beliefs

## Attack Vectors
- Logical validity of reasoning chains
- Hidden assumptions not made explicit
- Missing evidence for claimed conclusions
- Debt completeness
- Theorem applicability
- Category errors
- Counterexamples
- Alternative explanations
- Unsupported leaps
```

Both prompts explicitly state: "You have repository filesystem access. Read the seven research documents directly from docs/corpus/. Do not attempt to obtain them through ARGUS."

---

## 9. Phase 5 — End-to-End Dry Run

### Verification sequence

1. **Clean startup:** `argus serve` from clean state, no manual setup
2. **Seed visible:** /ui/insights shows 1 task + 1 belief; /ui/debts shows 8 obligations
3. **Work Agent:** Fresh OpenCode -> get_context -> read 7 files -> research -> submit_packet
4. **State projection:** Claims/evidence/edges/tasks visible in Insights + Debts
5. **Adversarial Agent:** Fresh OpenCode -> get_context -> read same 7 files -> attack -> submit_packet
6. **Challenge visible:** Contradiction/challenge appears in Insights
7. **Human discharge:** Authenticated discharge via Trust UI, pack-checked
8. **Mismatched discharge refused:** Attempt discharge with wrong evidence class -> refused
9. **Promotion gate:** Attempt promotion with remaining debt -> refused by Solvent

### Acceptance criteria

| AC | Description |
|----|-------------|
| AC1 | `argus serve` starts from clean local state |
| AC2 | Fresh startup creates minimal research seed |
| AC3 | 8 debt classes accepted by bmist@1.1.0 |
| AC4 | Debt validation from pack registry, not hardcoded list |
| AC5 | `argus mcp` communicates over stdio |
| AC6 | OpenCode can call `argus.get_context` |
| AC7 | OpenCode can call `argus.submit_packet` |
| AC8 | Fresh Work Agent reconstructs state via get_context |
| AC9 | Work Agent reads 7 documents from repository filesystem |
| AC10 | No corpus tables, embeddings, or vector indexes used |
| AC11 | No query parameter on get_context |
| AC12 | Exactly 2 MCP tools exist |
| AC13 | Documents never treated as authoritative Solvent state |
| AC14 | AI-derived reviews cannot independently discharge debt |
| AC15 | Work Agent packets persist correctly |
| AC16 | Adversarial packets persist correctly |
| AC17 | Contradictions visible in state/UI |
| AC18 | Human discharge authenticated and rule-checked |
| AC19 | Promotion with unresolved debt blocked |
| AC20 | Existing foundation tests pass |
| AC21 | No physics vocabulary in core packages |
| AC22 | No new persistence layer |
| AC23 | No new service |
| AC24 | No Solvent kernel changes |

---

## 10. Implementation Order

```
Phase 0  governance_ref resolution + full audit    [BLOCKING]
Phase 1  Domain pack bmist@1.1.0 + fix validation  [depends on 0]
Phase 2  MCP stdio transport                        [independent]
Phase 3  First-run seed                             [depends on 0, 1]
Phase 4  Role cards                                 [independent]
Phase 5  End-to-end dry run                         [depends on all]
```

---

## 11. Files Changed Summary

| File | Action | Phase |
|------|--------|-------|
| `domain-pack/bmist/v1/` (1.1.0 version) | New or edit | 1 |
| `internal/application/app.go` | Edit (validation, pack mapping) | 1 |
| `cmd/argus/main.go` | Edit (MCP stdio, seed call) | 2, 3 |
| `seed/seed.go` | New | 3 |
| `prompts/opencode-work.md` | New | 4 |
| `prompts/opencode-adversarial.md` | New | 4 |

**Zero Solvent kernel changes.**
**Zero new services.**
**Zero new databases.**
**Zero new MCP tools.**
**No corpus/ingestion/embedding infrastructure.**

---

## 12. What NOT to Build

- Corpus ingestion
- Embeddings
- Vector search
- get_context(query)
- argus.search_corpus
- ADD_DEBT
- Docker Compose
- New services or databases
- New MCP tools
- Prebuilt research graphs
- Workflow engines
- New UI
- Second authority mechanism
- Provenance subsystem

---

## 13. Final Principle

Before adding any component, ask:

> What concrete failure in the current dry run does this solve?

If the answer is "none yet", do not add it.

The POC is:

```
static research files
    +
authoritative ARGUS context
    +
two MCP tools
    +
EBP packet submission
    +
human adjudication
    +
adversarial replay
```

That is the experiment. Only after it works should retrieval infrastructure be reconsidered.
