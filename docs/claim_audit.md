# ARGUS Invariant Audit — Belief Claim Immutability

**Date:** 2026-09-20
**Scope:** Complete write-path audit for belief claim text immutability
**Codebase:** oracle (github.com/PithomLabs/oracle) + solvent-main (github.com/PithomLabs/solvent)

---

## 1. Verdict

**PASS** — existing belief claim text is effectively immutable across all application paths.

---

## 2. Desired Invariant

For an existing belief:

```
belief ID
    +
original claim text
```

must be immutable.

Research evolution must occur by creating a NEW belief ID and connecting it
to the previous belief through existing graph semantics (`derives`,
`contradicts`, etc.), rather than overwriting the original claim.

---

## 3. Evidence

### 3.1 Application Layer — `internal/application/app.go`

```
internal/application/app.go:357-364
Persist(...)
INSERT INTO belief (id, scenario_id, claim, claim_type, debt)
VALUES ($1::UUID, $2::UUID, $3, $4, $5)
ON CONFLICT (id) DO NOTHING

PASS — deterministic UUID from EntityID(scenarioID, "belief", claim).
       Same claim → same UUID → DO NOTHING → no mutation.
       Different claim → different UUID → new row.
```

```
internal/application/app.go:87-90
EntityID(scenarioID, "belief", claim)
SHA256(scenarioID + ":" + "belief" + ":" + claim)[0:16]

PASS — claim is baked into the ID derivation. Same claim produces same ID.
```

### 3.2 Solvent Kernel — `solvent-main/kernel/sql.go`

Every SQL statement that touches the `belief` table:

| Statement | Operation | Column mutated |
|-----------|-----------|----------------|
| `sqlEnterBelief` | INSERT | (new row) |
| `sqlRetireDebt` | UPDATE belief SET debt = array_remove(...) | `debt` only |
| `sqlPromote` | UPDATE belief SET status = 'promoted' | `status` only |
| `sqlRetractCascadeRetract` | UPDATE belief SET status = 'retracted' | `status` only |
| `sqlDischargeRetireDebt` | UPDATE belief SET debt = array_remove(...) | `debt` only |
| `sqlEnsureBelief` | CTE find-or-create by (scenario_id, claim) | (no update) |

```
PASS — no SQL statement in the kernel touches the claim column.
```

### 3.3 Solvent Kernel Methods — `solvent-main/kernel/kernel.go`

| Method | What it does | Mutates claim? |
|--------|-------------|----------------|
| `EnterBelief` | INSERT new belief (gen_random_uuid) | No |
| `EnsureBelief` | find-or-create by claim match | No |
| `AddEvidence` | INSERT evidence row | No |
| `RetireDebt` | UPDATE debt array | No |
| `Promote` | UPDATE status | No |
| `RetractCascade` | UPDATE status + cancel intents | No |
| `IntentOnPromoted` | INSERT intent | No |

```
PASS — no kernel method accepts a claim parameter for an existing belief ID.
       EnsureBelief returns existing ID without modifying the claim.
```

### 3.4 Solvent Service Layer — `solvent-main/service/ledger/ledger.go`

Thin passthrough to kernel methods. No independent claim mutation logic.

```
PASS — service layer delegates all operations to kernel.
```

### 3.5 Solvent API Handlers — `solvent-main/api/belief.go`

| Endpoint | Method | Claim mutation? |
|----------|--------|----------------|
| `POST /v1/beliefs` | EnterBelief | No — INSERT only |
| `GET /v1/beliefs/{id}` | ReadBelief | No — read only |
| `GET /v1/beliefs` | ListBeliefs | No — read only |
| `POST /v1/beliefs/{id}/debt/retire` | RetireDebt | No — UPDATE debt |
| `POST /v1/beliefs/{id}/promote` | PromoteBelief | No — UPDATE status |
| `POST /v1/beliefs/{id}/retract` | RetractBelief | No — UPDATE status |

```
PASS — no PUT/PATCH endpoint for belief claim exists.
       No claim parameter in retire/promote/retract endpoints.
```

### 3.6 Coordinator Persist — `coordinator/persist.go`

```
coordinator/persist.go:13-69 — SubmitPacket
  Calls c.solventClient.CreateBelief("", b.Claim, pkt.ScenarioID) for each belief.
  CreateBelief sends POST to /v1/beliefs (enter, not update).

PASS — always creates new belief via INSERT, never updates existing.
```

### 3.7 Coordinator Human Decisions — `coordinator/human.go`

```
coordinator/human.go:256-274 — handleReopen
  Calls c.solventClient.CreateBelief(req.BeliefID, req.Reason, req.ScenarioID)
  Then creates derives edge old→new.

PASS — reopen creates a NEW belief, never modifies the original claim.
```

### 3.8 MCP Adapter — `internal/mcp/adapter.go`

```
internal/mcp/adapter.go:57-87 — handleSubmitPacket
  Calls app.Compile → app.Validate → app.Persist
  app.Persist uses INSERT ... ON CONFLICT DO NOTHING.

PASS — agent-submitted packets cannot mutate existing claims.
```

### 3.9 Seed Path — `seed/seed.go`

```
seed/seed.go:64-72
INSERT INTO belief (id, scenario_id, claim, ...)
VALUES ($1, $2, $3, ...)
ON CONFLICT (id) DO NOTHING

PASS — seed uses fixed UUID + DO NOTHING. Idempotent, no mutation.
```

### 3.10 Trust UI Handlers — `internal/ui/ui.go`

Trust UI exposes: discharge (debt), promote (status).

```
PASS — no claim mutation endpoint exists in the UI.
```

### 3.11 Schema-Level — `internal/solventmigrations/migrate.go`

```sql
belief.id          UUID PRIMARY KEY
belief.claim       STRING NOT NULL
belief.status      CHECK (status IN ('entered','promoted','retracted'))
belief.debt        STRING[] NOT NULL DEFAULT '{}'
promoted_is_debt_free CHECK (status <> 'promoted' OR (...))
```

```
PASS — PK prevents duplicate IDs. No DB-level immutability trigger on claim,
       but none is needed because no application path attempts UPDATE claim.
```

### 3.12 Retraction Path

```
RetractCascade: UPDATE belief SET status = 'retracted'
               + UPDATE action_intent SET state = 'cancelled'
Does NOT touch claim.

PASS — retraction changes status, not content.
```

---

## 4. Mutation Surface

> Is there any application path that can mutate the claim text of an existing belief ID?

**No.** Across all 7 write surfaces examined:

| Surface | Operation | Claim mutated? |
|---------|-----------|----------------|
| App.Persist | INSERT ON CONFLICT DO NOTHING | No — different claim = different UUID |
| Kernel.EnterBelief | INSERT (new UUID) | No — always new row |
| Kernel.EnsureBelief | find-or-create by claim | No — returns existing, never updates |
| Kernel.RetireDebt | UPDATE debt only | No |
| Kernel.Promote | UPDATE status only | No |
| Kernel.RetractCascade | UPDATE status only | No |
| UI discharge/promote | Delegates to kernel | No |
| Coordinator.Reopen | CreateBelief + edge | No — creates NEW belief |
| Seed | INSERT ON CONFLICT DO NOTHING | No |

The deterministic UUID mechanism (`EntityID`) provides an additional structural
guarantee: same claim content always produces the same ID, so
`ON CONFLICT (id) DO NOTHING` silently deduplicates. A resubmission with a
**different** claim for the "same" logical belief produces a **different** UUID,
creating a new row — it does not overwrite.

---

## 5. Lifecycle Assessment

> Does the implementation support: enter → challenge → new proposition → graph link → historical continuity?

**Yes.** The integration test (`TestAdversarialContradictionVisible`, `app_test.go:589-653`)
demonstrates exactly this:

```
B1 (contested claim)
  |
  | adversarial submission creates B2 with contradicts edge
  v
B2 (adversarial refutation) --contradicts--> B1
```

The `handleReopen` path (`coordinator/human.go:256-274`) additionally demonstrates:

```
old retracted belief
  | reopen creates new belief with derives edge
  v
new belief --derives--> old belief
```

B1.claim never becomes B3.claim. Research evolution creates new belief IDs
connected through graph semantics (`derives`, `contradicts`), preserving
historical continuity.

---

## 6. Idempotency Classification

| Scenario | Behavior | Classification |
|----------|----------|----------------|
| Same packet_id + same claim | `ON CONFLICT DO NOTHING` → returns existing ID | C — creates/returns new belief (idempotent) |
| Same packet_id + different claim | Different UUID → new row | C — creates new belief |
| Different packet_id + same claim | Same UUID → `DO NOTHING` | C — creates/returns existing belief (idempotent) |
| Different packet_id + different claim | Different UUID → new row | C — creates new belief |

No path overwrites an existing belief's claim.

---

## 7. Test Gap

**No explicit same-ID claim-immutability regression test exists.**

The idempotency tests (`TestIdempotencyDuplicatePacket`, `TestFullIntegration`
step 10) verify that resubmitting the same packet produces no duplicates, but
they do not explicitly verify that a resubmission with a **different claim but
same packet_id** does not overwrite.

The deterministic UUID mechanism makes such an overwrite structurally impossible
(different claim → different UUID → new row), but an explicit regression test
would document this invariant.

### Smallest test to add:

```go
func TestSameBeliefIDDifferentClaimCreatesNewRow(t *testing.T) {
    db := testDB(t)
    app := New(db)
    ctx := context.Background()
    scenarioID := "00000000-0000-0000-0000-000000000099"

    // 1. Insert B1 with claim C1
    pkt1 := &packetv1.Packet{
        SchemaVersion: packetv1.SchemaVersion, Role: packetv1.RoleWork,
        PacketID: "pkt-001", PackRef: "bmist@1.0.0", ScenarioID: scenarioID,
        Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "original claim", ClaimType: "derived"}},
    }
    tx1, _ := db.BeginTx(ctx, nil)
    r1, _ := app.Persist(ctx, tx1, pkt1)
    tx1.Commit()
    b1ID := r1.BeliefIDs["b1"]

    // 2. Verify B1.claim == "original claim"
    var claim1 string
    db.QueryRowContext(ctx, `SELECT claim FROM belief WHERE id = $1::UUID`, b1ID).Scan(&claim1)
    if claim1 != "original claim" {
        t.Fatalf("B1 claim should be 'original claim', got %q", claim1)
    }

    // 3. Attempt to persist with different packet_id but DIFFERENT claim
    //    This produces a DIFFERENT belief ID (deterministic UUID changes).
    pkt2 := &packetv1.Packet{
        SchemaVersion: packetv1.SchemaVersion, Role: packetv1.RoleWork,
        PacketID: "pkt-002", PackRef: "bmist@1.0.0", ScenarioID: scenarioID,
        Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "mutated claim", ClaimType: "derived"}},
    }
    tx2, _ := db.BeginTx(ctx, nil)
    r2, _ := app.Persist(ctx, tx2, pkt2)
    tx2.Commit()
    b2ID := r2.BeliefIDs["b1"]

    // 4. B1 ID != B2 ID (different claim → different UUID)
    if b1ID == b2ID {
        t.Fatalf("different claims must produce different IDs, both got %s", b1ID)
    }

    // 5. Re-read B1 — claim unchanged
    db.QueryRowContext(ctx, `SELECT claim FROM belief WHERE id = $1::UUID`, b1ID).Scan(&claim1)
    if claim1 != "original claim" {
        t.Fatalf("B1 claim mutated! Expected 'original claim', got %q", claim1)
    }

    // 6. No silent overwrite — exactly 2 beliefs exist
    var count int
    db.QueryRowContext(ctx, `SELECT count(*) FROM belief WHERE scenario_id = $1::UUID`, scenarioID).Scan(&count)
    if count != 2 {
        t.Fatalf("expected 2 beliefs (original + new), got %d", count)
    }
}
```

---

## 8. Summary

The belief claim immutability invariant holds through a combination of:

1. **Deterministic UUID derivation** — claim text is an input to the belief ID,
   so different claims always produce different IDs.

2. **INSERT-only persistence** — `App.Persist` uses `INSERT ... ON CONFLICT (id)
   DO NOTHING`, which silently deduplicates identical claims and creates new
   rows for different claims.

3. **No UPDATE claim SQL** — every `UPDATE belief SET ...` in the entire
   Solvent kernel and application layer targets only `debt` or `status`,
   never `claim`.

4. **No API endpoint for claim mutation** — the Solvent REST API exposes
   create/read/promote/retract/retire-debt, but no update-claim endpoint.

5. **Graph-based research evolution** — adversarial challenges and reopen
   operations create NEW belief IDs connected through `derives`/`contradicts`
   edges, preserving historical continuity.

No database-level immutability constraint (trigger, immutable column, audit
column) exists on the claim column, but none is needed because no application
path attempts to mutate it.
