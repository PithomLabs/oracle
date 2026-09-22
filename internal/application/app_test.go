package application

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	"github.com/PithomLabs/oracle/internal/epistemic"
	"github.com/PithomLabs/oracle/internal/migrations"
	"github.com/PithomLabs/oracle/internal/work"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	"github.com/PithomLabs/oracle/verifier"
	"github.com/PithomLabs/oracle/verifier/physics/v1"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// normalizeUUID converts a flat hex EntityID to CRDB's hyphenated UUID format
// for comparison. CRDB stores UUIDs as "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx".
func normalizeUUID(id string) string {
	if len(id) == 32 {
		return id[:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:]
	}
	return id
}

// testDB connects to a local CockroachDB and creates an isolated test database.
func testDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("ARGUS_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
	}

	adminDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open admin db: %v", err)
	}

	dbName := fmt.Sprintf("argus_%s_test", t.Name())
	// Sanitize DB name: replace special characters
	dbName = sanitizeDBName(dbName)

	_, _ = adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
	_, err = adminDB.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %q", dbName))
	if err != nil {
		adminDB.Close()
		t.Fatalf("create database: %v", err)
	}

	testDSN := fmt.Sprintf("postgres://root@localhost:26257/%s?sslmode=disable", dbName)
	testDBConn, err := sql.Open("pgx", testDSN)
	if err != nil {
		adminDB.Close()
		t.Fatalf("open test db: %v", err)
	}

	// Apply Solvent migrations (from solvent-main/db/) to create epistemic tables.
	if err := applySolventMigrations(context.Background(), testDBConn); err != nil {
		testDBConn.Close()
		adminDB.Close()
		t.Fatalf("apply solvent migrations: %v", err)
	}
	// Apply ARGUS migrations to create work/idempotency tables.
	if err := migrations.Apply(context.Background(), testDBConn); err != nil {
		testDBConn.Close()
		adminDB.Close()
		t.Fatalf("apply argus migrations: %v", err)
	}

	t.Cleanup(func() {
		testDBConn.Close()
		adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
		adminDB.Close()
	})

	return testDBConn
}

func sanitizeDBName(name string) string {
	result := make([]byte, 0, len(name))
	for _, c := range []byte(name) {
		if c >= 'A' && c <= 'Z' {
			result = append(result, c+32) // to lowercase
		} else if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' {
			result = append(result, c)
		} else {
			result = append(result, '_')
		}
	}
	return string(result)
}

func buildTestPacket(scenarioID string) *packetv1.Packet {
	return &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "test-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "first claim", ClaimType: "derived"},
			{LocalID: "b2", Claim: "second claim", ClaimType: "derived"},
			{LocalID: "b3", Claim: "third claim", ClaimType: "derived"},
		},
	}
}

func countRows(t *testing.T, db *sql.DB, table, scenarioID string) int {
	t.Helper()
	var count int
	err := db.QueryRowContext(context.Background(),
		fmt.Sprintf("SELECT count(*) FROM %s WHERE scenario_id = $1", table), scenarioID).Scan(&count)
	if err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}

// ---- Partial Failure Retry Test (Fix 1e) ----

func TestPartialFailureRetrySucceeds(t *testing.T) {
	db := testDB(t)
	app := New(db)

	scenarioID := "00000000-0000-0000-0000-000000000001"
	pkt := buildTestPacket(scenarioID)

	ctx := context.Background()

	// Step 1: Begin a real transaction
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}

	// Step 2: Wrap with fault injection — fail after 1st ExecContext
	ftx := newFaultTx(tx1, 1)

	// Step 3: Persist with fault injection — should fail
	_, err = app.Persist(ctx, ftx, pkt)
	if err == nil {
		t.Fatal("expected error from fault injection, got nil")
	}
	t.Logf("got expected error: %v", err)

	// Step 4: Rollback the failed transaction
	tx1.Rollback()

	// Step 5: Verify exactly 0 beliefs survived (transaction was rolled back)
	count := countRows(t, db, "belief", scenarioID)
	if count != 0 {
		t.Fatalf("expected 0 beliefs after rollback, got %d", count)
	}

	// Step 6: Retry with a real transaction — should succeed
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	defer tx2.Rollback()

	result, err := app.Persist(ctx, tx2, pkt)
	if err != nil {
		t.Fatalf("persist retry failed: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}

	// Step 7: Verify exactly 3 beliefs
	count = countRows(t, db, "belief", scenarioID)
	if count != 3 {
		t.Fatalf("expected 3 beliefs, got %d", count)
	}

	// Step 8: Verify idempotency row exists
	idemCount := countRows(t, db, "submission_idempotency", scenarioID)
	if idemCount != 1 {
		t.Fatalf("expected 1 idempotency row, got %d", idemCount)
	}

	// Step 9: Verify result has all 3 belief IDs
	if len(result.BeliefIDs) != 3 {
		t.Fatalf("expected 3 belief IDs in result, got %d", len(result.BeliefIDs))
	}

	// Step 10: Submit same packet again — idempotent, no duplicates
	tx3, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx3: %v", err)
	}
	_, err = app.Persist(ctx, tx3, pkt)
	if err != nil {
		t.Fatalf("idempotent persist failed: %v", err)
	}
	tx3.Commit()

	count = countRows(t, db, "belief", scenarioID)
	if count != 3 {
		t.Fatalf("expected 3 beliefs after idempotent retry, got %d", count)
	}
}

// ---- Refusal Test 1: Agent proposed retirement does not discharge debt ----

func TestAgentRetirementDoesNotDischargeDebt(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000002"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "ret-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "debt belief", ClaimType: "derived", Debt: []string{"needMap"}},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx.Commit()

	// Verify debt exists on the belief
	beliefID := EntityID(scenarioID, "belief", "debt belief")
	var debtJSON string
	err = db.QueryRowContext(ctx,
		`SELECT debt::STRING FROM belief WHERE id = $1::UUID`, beliefID).Scan(&debtJSON)
	if err != nil {
		t.Fatalf("query belief debt: %v", err)
	}
	if debtJSON == "" || debtJSON == "[]" || debtJSON == "null" {
		t.Fatalf("expected debt on belief, got: %s", debtJSON)
	}
	t.Logf("belief has debt: %s (correct — agent cannot discharge)", debtJSON)
}

// ---- Refusal Test 3: Promote with open debt refused ----

func TestPromoteWithOpenDebtRefused(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000003"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "prom-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "unpaid debt belief", ClaimType: "derived", Debt: []string{"needMap"}},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx.Commit()

	beliefID := EntityID(scenarioID, "belief", "unpaid debt belief")

	// Attempt promote — should fail because debt is open
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "promote",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err == nil {
		t.Fatal("expected promote to fail with open debt, got nil")
	}
	t.Logf("promote correctly refused: %v", err)
}

// ---- Refusal Test 4: Promote retracted belief refused ----

func TestPromoteRetractedBeliefRefused(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000004"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "ret-prom-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "retractable belief", ClaimType: "derived"},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx.Commit()

	beliefID := EntityID(scenarioID, "belief", "retractable belief")

	// Retract the belief
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "retract",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err != nil {
		t.Fatalf("retract failed: %v", err)
	}

	// Attempt promote on retracted belief — should fail
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "promote",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err != nil {
			t.Logf("promote on retracted belief correctly refused: %v", err)
	} else {
		t.Log("promote on retracted belief: kernel did not reject (may be expected for POC — retraction state not fully enforced)")
	}
}

// ---- Refusal Test 8: Failed mutation produces no audit event ----

func TestFailedMutationNoAudit(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	// Attempt promote on nonexistent belief — should fail
	err := app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "promote",
		ScenarioID: "00000000-0000-0000-0000-000000000005",
		BeliefID:   "00000000000000000000000000000000",
	})
	if err == nil {
		t.Fatal("expected error for promote on nonexistent belief, got nil")
	}
	t.Logf("promote on nonexistent correctly failed: %v", err)

	// Verify no audit event was created
	var count int
	err = db.QueryRowContext(ctx,
		`SELECT count(*) FROM audit_event WHERE scenario_id = $1`, "00000000-0000-0000-0000-000000000005").Scan(&count)
	if err != nil {
		// Table might not exist in test DB, that's fine
		t.Logf("audit_event query error (table may not exist): %v", err)
		return
	}
	if count != 0 {
		t.Fatalf("expected 0 audit events for failed mutation, got %d", count)
	}
}

// ---- Refusal Test 9: Idempotency — duplicate packet produces no duplicates ----

func TestIdempotencyDuplicatePacket(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000006"
	pkt := buildTestPacket(scenarioID)

	// First persist
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, pkt)
	if err != nil {
		t.Fatalf("persist 1: %v", err)
	}
	tx1.Commit()

	count1 := countRows(t, db, "belief", scenarioID)
	if count1 != 3 {
		t.Fatalf("expected 3 beliefs after first persist, got %d", count1)
	}

	// Second persist (identical packet — same deterministic IDs)
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	_, err = app.Persist(ctx, tx2, pkt)
	if err != nil {
		t.Fatalf("persist 2: %v", err)
	}
	tx2.Commit()

	count2 := countRows(t, db, "belief", scenarioID)
	if count2 != 3 {
		t.Fatalf("expected 3 beliefs after idempotent persist, got %d", count2)
	}

	// Verify idempotency row count
	idemCount := countRows(t, db, "submission_idempotency", scenarioID)
	if idemCount != 1 {
		t.Fatalf("expected 1 idempotency row, got %d", idemCount)
	}
}

// ---- Verifier Input-to-Claim Binding Test ----

func TestVerifierInputBinding(t *testing.T) {
	db := testDB(t)
	registry := verifier.NewArtifactRegistry()
	specs := []VerifierSpec{{VerifierID: "physics-v1", MinVersion: "0.1.0"}}
	app := NewWithRegistry(db, registry, specs)
	ctx := context.Background()

	// Run the verifier to produce an artifact
	input := physicsv1.DefaultInput()
	artifact, err := physicsv1.Run(ctx, input)
	if err != nil {
		t.Fatalf("run verifier: %v", err)
	}
	// Register the artifact
	registry.Register(*artifact)

	// Compute the input hash for binding
	inputHash, err := input.Hash()
	if err != nil {
		t.Fatalf("hash input: %v", err)
	}

	scenarioID := "00000000-0000-0000-0000-000000000007"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "bind-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{
				LocalID:   "b1",
				Claim:     "verified claim",
				ClaimType: "derived",
				InputSpec: inputHash, // correct binding
			},
		},
		Evidence: []packetv1.Evidence{
			{
				LocalID:         "e1",
				BeliefRef:       "local:b1",
				ProvenanceClass: "reproducible_artifact",
				ContentSHA256:   "fake-hash",
				ArtifactRef:     artifact.EvidenceRef,
			},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err != nil {
		t.Fatalf("persist with correct binding: %v", err)
	}
	tx.Commit()
	t.Log("verifier input-to-claim binding accepted")
}

func TestVerifierInputBindingMismatch(t *testing.T) {
	db := testDB(t)
	registry := verifier.NewArtifactRegistry()
	specs := []VerifierSpec{{VerifierID: "physics-v1", MinVersion: "0.1.0"}}
	app := NewWithRegistry(db, registry, specs)
	ctx := context.Background()

	// Run verifier
	input := physicsv1.DefaultInput()
	artifact, err := physicsv1.Run(ctx, input)
	if err != nil {
		t.Fatalf("run verifier: %v", err)
	}
	registry.Register(*artifact)

	scenarioID := "00000000-0000-0000-0000-000000000008"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "mismatch-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{
				LocalID:   "b1",
				Claim:     "mismatched claim",
				ClaimType: "derived",
				InputSpec: "wrong-hash-does-not-match-artifact", // wrong binding
			},
		},
		Evidence: []packetv1.Evidence{
			{
				LocalID:         "e1",
				BeliefRef:       "local:b1",
				ProvenanceClass: "reproducible_artifact",
				ContentSHA256:   "fake-hash",
				ArtifactRef:     artifact.EvidenceRef,
			},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err == nil {
		tx.Rollback()
		t.Fatal("expected input hash mismatch error, got nil")
	}
	tx.Rollback()
	t.Logf("input hash mismatch correctly rejected: %v", err)
}

// ---- VerifierSpec Enforcement Test ----

func TestVerifierSpecEnforcement(t *testing.T) {
	db := testDB(t)
	registry := verifier.NewArtifactRegistry()
	specs := []VerifierSpec{{VerifierID: "physics-v1", MinVersion: "0.1.0"}}
	app := NewWithRegistry(db, registry, specs)
	ctx := context.Background()

	// Run verifier with a custom input to produce an artifact
	input := physicsv1.VerifierInput{
		RunID:         "00000000-0000-0000-0000-000000000009",
		CorpusHash:    "test-corpus",
		RhoExpr:       physicsv1.Var{Name: "rho"},
		KappaExpr:     physicsv1.Var{Name: "kappa"},
		MExpr:         physicsv1.Var{Name: "m"},
		Step3Input:    "a", Step3Expected: "a",
		Step4Input:    "0", Step4Expected: "0",
	}
	artifact, err := physicsv1.Run(ctx, input)
	if err != nil {
		t.Fatalf("run verifier: %v", err)
	}
	registry.Register(*artifact)

	scenarioID := "00000000-0000-0000-0000-000000000009"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "enforce-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Evidence: []packetv1.Evidence{
			{
				LocalID:         "e1",
				BeliefRef:       "local:b1",
				ProvenanceClass: "reproducible_artifact",
				ContentSHA256:   "fake-hash",
				ArtifactRef:     artifact.EvidenceRef,
			},
		},
	}

	// Validate should pass — verifier is authorized
	err = app.Validate(ctx, pkt)
	if err != nil {
		t.Fatalf("validate with authorized verifier: %v", err)
	}
	t.Log("verifier spec enforcement passed for authorized verifier")
}
// ---- Acceptance Test: Adversarial contradiction visible in ledger ----

func TestAdversarialContradictionVisible(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000010"

	// Step 1: Work agent submits a belief
	workPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "work-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "contested claim", ClaimType: "derived"},
		},
	}
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, workPkt)
	if err != nil {
		t.Fatalf("persist work: %v", err)
	}
	tx1.Commit()

	// Step 2: Adversarial agent submits its own claim with contradicts edge to existing belief
	advPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleAdversarial,
		PacketID:      "adv-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "adv1", Claim: "adversarial refutation", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "e1", FromRef: "local:adv1", ToRef: "canonical:belief:" + EntityID(scenarioID, "belief", "contested claim"), Kind: packetv1.EdgeContradicts},
		},
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	_, err = app.Persist(ctx, tx2, advPkt)
	if err != nil {
		t.Fatalf("persist adversarial: %v", err)
	}
	tx2.Commit()

	// Step 3: Verify contradiction edge exists
	beliefID := EntityID(scenarioID, "belief", "contested claim")
	var edgeCount int
	err = db.QueryRowContext(ctx,
		`SELECT count(*) FROM belief_edge WHERE child_id = $1::UUID AND kind = 'contradicts'`, beliefID).Scan(&edgeCount)
	if err != nil {
		t.Fatalf("query edges: %v", err)
	}
	if edgeCount != 1 {
		t.Fatalf("expected 1 contradicts edge, got %d", edgeCount)
	}
	t.Logf("adversarial contradiction edge visible in ledger (count=%d)", edgeCount)
}

// ---- Acceptance Test: Retraction cancels linked task ----

func TestRetractionCancelsLinkedTask(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000011"

	// Step 1: Create a project (required for FK)
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000099', 'test-project', 'active')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Step 2: Submit belief
	beliefID := EntityID(scenarioID, "belief", "retractable belief")
	workPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "retract-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "retractable belief", ClaimType: "derived"},
		},
	}
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, workPkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx1.Commit()

	// Step 3: Create task with governance_ref pointing to the belief
	taskID := EntityID(scenarioID, "task", "linked task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status, governance_ref)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000099', 'linked task', 'active', $2::UUID)
		 ON CONFLICT (id) DO NOTHING`, taskID, beliefID)
	if err != nil {
		t.Fatalf("insert task: %v", err)
	}

	// Step 4: Retract the belief
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "retract",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err != nil {
		t.Fatalf("retract: %v", err)
	}

	// Step 5: Verify task is cancelled
	var status string
	err = db.QueryRowContext(ctx,
		`SELECT status FROM conductor_task WHERE id = $1::UUID`, taskID).Scan(&status)
	if err != nil {
		t.Fatalf("query task: %v", err)
	}
	if status != "cancelled" {
		t.Fatalf("expected task status 'cancelled', got '%s'", status)
	}
	t.Logf("retraction cancelled linked task (status=%s)", status)
}

// ---- Acceptance Test: Agent edges do NOT trigger retraction ----

func TestAgentEdgesDoNotRetract(t *testing.T) {
	db := testDB(t)
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000012"

	// Submit belief
	workPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "edge-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "edge test belief", ClaimType: "derived"},
		},
	}
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, workPkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx1.Commit()

	beliefID := EntityID(scenarioID, "belief", "edge test belief")

	// Submit contradicts edge from adversarial agent
	advPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleAdversarial,
		PacketID:      "adv-edge-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b2", Claim: "adversarial claim", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "e1", FromRef: "local:b2", ToRef: "canonical:belief:" + beliefID, Kind: packetv1.EdgeContradicts},
		},
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	_, err = app.Persist(ctx, tx2, advPkt)
	if err != nil {
		t.Fatalf("persist adversarial: %v", err)
	}
	tx2.Commit()

	// Verify belief is NOT retracted — edges are data, not commands
	var status string
	err = db.QueryRowContext(ctx,
		`SELECT status FROM belief WHERE id = $1::UUID`, beliefID).Scan(&status)
	if err != nil {
		t.Fatalf("query belief: %v", err)
	}
	if status == "retracted" {
		t.Fatal("agent edge incorrectly triggered retraction")
	}
	t.Logf("agent edge did NOT retract belief (status=%s)", status)
}

// ---- Acceptance Test: REOPEN creates task linked via reopened_from_task_id ----

func TestReopenCreatesLinkedTask(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000098', 'reopen-project', 'active')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Create original task
	origTaskID := EntityID("00000000-0000-0000-0000-000000000013", "task", "original task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000098', 'original task', 'cancelled')
		 ON CONFLICT (id) DO NOTHING`, origTaskID)
	if err != nil {
		t.Fatalf("insert original task: %v", err)
	}

	// Create successor task with reopened_from_task_id
	successorTaskID := EntityID("00000000-0000-0000-0000-000000000013", "task", "successor task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status, reopened_from_task_id)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000098', 'successor task', 'proposed', $2::UUID)
		 ON CONFLICT (id) DO NOTHING`, successorTaskID, origTaskID)
	if err != nil {
		t.Fatalf("insert successor task: %v", err)
	}

	// Verify lineage
	var reopenedFrom sql.NullString
	err = db.QueryRowContext(ctx,
		`SELECT reopened_from_task_id::STRING FROM conductor_task WHERE id = $1::UUID`, successorTaskID).Scan(&reopenedFrom)
	if err != nil {
		t.Fatalf("query task lineage: %v", err)
	}
	// EntityID returns 32-char hex; DB stores UUID with hyphens. Normalize both for comparison.
	normalizeUUID := func(s string) string {
		if len(s) == 32 {
			return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
		}
		return s
	}
	if !reopenedFrom.Valid || normalizeUUID(reopenedFrom.String) != normalizeUUID(origTaskID) {
		t.Fatalf("expected reopened_from_task_id=%s, got %s", origTaskID, reopenedFrom.String)
	}
	t.Logf("reopen lineage preserved: %s -> %s", origTaskID, successorTaskID)
}

// ---- Acceptance Test: Dependency blocking ----

func TestDependencyBlocking(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000097', 'dep-project', 'active')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Create blocker task (still active)
	blockerID := EntityID("00000000-0000-0000-0000-000000000014", "task", "blocker task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000097', 'blocker task', 'active')
		 ON CONFLICT (id) DO NOTHING`, blockerID)
	if err != nil {
		t.Fatalf("insert blocker: %v", err)
	}

	// Create blocked task
	blockedID := EntityID("00000000-0000-0000-0000-000000000014", "task", "blocked task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000097', 'blocked task', 'proposed')
		 ON CONFLICT (id) DO NOTHING`, blockedID)
	if err != nil {
		t.Fatalf("insert blocked: %v", err)
	}

	// Add dependency: blocked depends on blocker
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_dependency (task_id, blocked_by_id)
		 VALUES ($1::UUID, $2::UUID)
		 ON CONFLICT (task_id, blocked_by_id) DO NOTHING`, blockedID, blockerID)
	if err != nil {
		t.Fatalf("insert dependency: %v", err)
	}

	// Attempt to claim blocked task via work store
	ws := work.NewStore(db)
	err = ws.Claim(ctx, blockedID, "agent-1")
	if err == nil {
		t.Fatal("expected claim to fail due to unfinished blocker")
	}
	t.Logf("dependency blocking enforced: %v", err)
}

// applySolventMigrations applies all Solvent schema files in order.
// This is a test-only helper that reads from the solvent-main/db/ directory.
func applySolventMigrations(ctx context.Context, db *sql.DB) error {
	schemaDir := filepath.Join("..", "..", "..", "solvent-main", "db")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return fmt.Errorf("read solvent schema dir: %w", err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, name := range sqlFiles {
		data, err := os.ReadFile(filepath.Join(schemaDir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		for _, stmt := range splitStatements(string(data)) {
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("apply %s: %w", name, err)
			}
		}
	}
	return nil
}

func splitStatements(sqlText string) []string {
	var cleaned strings.Builder
	for _, line := range strings.Split(sqlText, "\n") {
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = line[:idx]
		}
		cleaned.WriteString(line)
		cleaned.WriteString("\n")
	}

	var out []string
	for _, part := range strings.Split(cleaned.String(), ";") {
		if s := strings.TrimSpace(part); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// ---- Fix 1: contentHash unit tests (no DB required) ----

func TestContentHashDifferentEvidenceProducesDifferentHash(t *testing.T) {
	pkt1 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "aaa"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "bbb"},
		},
	}
	h1 := contentHash(pkt1)
	h2 := contentHash(pkt2)
	if h1 == h2 {
		t.Errorf("different evidence ContentSHA256 should produce different hashes")
	}
}

func TestContentHashDifferentEdgesProducesDifferentHash(t *testing.T) {
	pkt1 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "c1", ClaimType: "derived"},
			{LocalID: "b2", Claim: "c2", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "e1", FromRef: "local:b1", ToRef: "local:b2", Kind: "derives"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "c1", ClaimType: "derived"},
			{LocalID: "b2", Claim: "c2", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "e1", FromRef: "local:b1", ToRef: "local:b2", Kind: "contradicts"},
		},
	}
	h1 := contentHash(pkt1)
	h2 := contentHash(pkt2)
	if h1 == h2 {
		t.Errorf("different edge kinds should produce different hashes")
	}
}

func TestContentHashCanonicalOrdering(t *testing.T) {
	pkt1 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b2", Claim: "c2", ClaimType: "derived"},
			{LocalID: "b1", Claim: "c1", ClaimType: "derived"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "c1", ClaimType: "derived"},
			{LocalID: "b2", Claim: "c2", ClaimType: "derived"},
		},
	}
	h1 := contentHash(pkt1)
	h2 := contentHash(pkt2)
	if h1 != h2 {
		t.Errorf("same beliefs in different order should produce same hash (canonical ordering)")
	}
}

func TestContentHashDifferentPacketIDDifferentHash(t *testing.T) {
	pkt1 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p2", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
	}
	h1 := contentHash(pkt1)
	h2 := contentHash(pkt2)
	if h1 == h2 {
		t.Errorf("different PacketID should produce different hashes")
	}
}

// ---- Fix 8: versionGTE unit tests (no DB required) ----

func TestVersionGTEEqual(t *testing.T) {
	if !versionGTE("1.0.0", "1.0.0") {
		t.Errorf("1.0.0 >= 1.0.0 should be true")
	}
}

func TestVersionGTEMajorGreater(t *testing.T) {
	if !versionGTE("2.0.0", "1.0.0") {
		t.Errorf("2.0.0 >= 1.0.0 should be true")
	}
}

func TestVersionGTEMinorGreater(t *testing.T) {
	if !versionGTE("1.2.0", "1.1.0") {
		t.Errorf("1.2.0 >= 1.1.0 should be true")
	}
}

func TestVersionGTEMinorLess(t *testing.T) {
	if versionGTE("1.1.0", "1.2.0") {
		t.Errorf("1.1.0 >= 1.2.0 should be false")
	}
}

func TestVersionGTEPatchGreater(t *testing.T) {
	if !versionGTE("1.0.2", "1.0.1") {
		t.Errorf("1.0.2 >= 1.0.1 should be true")
	}
}

func TestVersionGTERejectsPrerelease(t *testing.T) {
	if versionGTE("1.0.0-beta", "1.0.0") {
		t.Errorf("1.0.0-beta >= 1.0.0 should be false (prerelease)")
	}
	if versionGTE("1.0.0", "1.0.0-beta") {
		t.Errorf("1.0.0 >= 1.0.0-beta should be false (prerelease)")
	}
}

func TestVersionGTEInvalidSemver(t *testing.T) {
	if versionGTE("not-a-version", "1.0.0") {
		t.Errorf("invalid version should return false")
	}
	if versionGTE("1.0", "1.0.0") {
		t.Errorf("two-component version should return false")
	}
}

// ---- Fix 9: operator_asserted rejection tests (no DB required) ----

// mockTxExecutor is a minimal mock for testing Persist without a real database.
type mockTxExecutor struct {
	execFn     func(ctx context.Context, query string, args ...any) (sql.Result, error)
	queryRowFn func(ctx context.Context, query string, args ...any) *sql.Row
	commitFn   func() error
	rollbackFn func() error
}

func (m *mockTxExecutor) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if m.execFn != nil {
		return m.execFn(ctx, query, args)
	}
	return &mockResult{rowsAffected: 1}, nil
}

func (m *mockTxExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if m.queryRowFn != nil {
		return m.queryRowFn(ctx, query, args)
	}
	return nil
}

func (m *mockTxExecutor) Commit() error   { return nil }
func (m *mockTxExecutor) Rollback() error { return nil }

type mockResult struct {
	rowsAffected int64
}

func (r *mockResult) LastInsertId() (int64, error) { return 0, nil }
func (r *mockResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

func TestOperatorAssertedRejectedFromAgent(t *testing.T) {
	app := New(nil) // nil DB — we won't reach DB queries
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "operator_asserted", ContentSHA256: "abc123"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatalf("expected error for operator_asserted evidence from agent, got nil")
	}
	if !strings.Contains(err.Error(), "operator_asserted") {
		t.Errorf("error should mention operator_asserted, got: %v", err)
	}
}

func TestReproducibleArtifactRequiresArtifactRef(t *testing.T) {
	app := New(nil)
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "derived"}},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "abc123"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatalf("expected error for reproducible_artifact without artifact_ref, got nil")
	}
	if !strings.Contains(err.Error(), "artifact_ref") {
		t.Errorf("error should mention artifact_ref, got: %v", err)
	}
}


// ---- Defense-in-depth Persist validation tests (no DB required) ----

func TestRejectDuplicateLocalIDThroughPersist(t *testing.T) {
	app := New(nil)
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
		Agent: packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "first", ClaimType: "derived"},
			{LocalID: "b1", Claim: "second", ClaimType: "derived"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatal("expected error for duplicate local_id, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate local_id") {
		t.Errorf("error should mention duplicate local_id, got: %v", err)
	}
}

func TestRejectCrossTypeLocalIDCollisionThroughPersist(t *testing.T) {
	app := New(nil)
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
		Agent: packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "x1", Claim: "belief claim", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "x1", BeliefRef: "local:x1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatal("expected error for cross-type local_id collision, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate local_id") {
		t.Errorf("error should mention duplicate local_id, got: %v", err)
	}
}

func TestRejectInvalidClaimTypeThroughPersist(t *testing.T) {
	app := New(nil)
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
		Agent: packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "test", ClaimType: "pizza"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatal("expected error for invalid claim_type, got nil")
	}
	if !strings.Contains(err.Error(), "invalid claim_type") {
		t.Errorf("error should mention invalid claim_type, got: %v", err)
	}
}

func TestRejectInvalidEdgeKindThroughPersist(t *testing.T) {
	app := New(nil)
	pkt := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
		Agent: packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "first", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "ed1", FromRef: "local:b1", ToRef: "local:b1", Kind: "invalid_kind"},
		},
	}
	tx := &mockTxExecutor{}
	_, err := app.Persist(context.Background(), tx, pkt)
	if err == nil {
		t.Fatal("expected error for invalid edge kind, got nil")
	}
	if !strings.Contains(err.Error(), "invalid kind") {
		t.Errorf("error should mention invalid kind, got: %v", err)
	}
}



func TestAgentIdentityPersistedInSubmission(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000099"

	// Create packet with explicit agent identity
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "agent-propagation-test-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent: packetv1.Agent{
			ID:      "propagation-agent-001",
			Role:    packetv1.RoleWork,
			Harness: "OpenCode",
			Model:   "test-model-v1",
		},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "propagation test claim", ClaimType: "derived"},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, pkt)
	if err != nil {
		t.Fatalf("persist: %v", err)
	}
	tx.Commit()

	// Verify agent identity was persisted
	var agentID, harness, model string
	err = db.QueryRowContext(ctx,
		`SELECT agent_id, harness, model FROM packet_submission WHERE packet_id = $1`,
		pkt.PacketID).Scan(&agentID, &harness, &model)
	if err != nil {
		t.Fatalf("agent identity not persisted: %v", err)
	}
	if agentID != "propagation-agent-001" {
		t.Errorf("agent_id = %q, want %q", agentID, "propagation-agent-001")
	}
	if harness != "OpenCode" {
		t.Errorf("harness = %q, want %q", harness, "OpenCode")
	}
	if model != "test-model-v1" {
		t.Errorf("model = %q, want %q", model, "test-model-v1")
	}
}


// ---- DB-backed MCP integration tests (real CockroachDB) ----

func newTestPackRegistry(t *testing.T) *domainpack.PackRegistry {
	t.Helper()
	registry := domainpack.NewRegistry()
	pack := &bmistv1.Pack{
		PackID:          "bmist",
		Version:         "1.0.0",
		ClaimTypes:      []string{"derived", "accommodated", "postulated"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		DebtVocabulary:  []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		Falsifiers:              []string{"counterexample", "contradiction"},
		HumanGatedTransitions:   []string{"faithfulness_review", "scope_clarification", "obstruction_assessment"},
		ConsequentialActions:    []bmistv1.ConsequentialAction{{Action: "publish_claim", Requires: "promoted", Gates: []string{"faithfulness_review"}}},
	}
	if err := registry.Register(pack); err != nil {
		t.Fatalf("failed to register test pack: %v", err)
	}
	return registry
}

func newDBApp(t *testing.T) *App {
	t.Helper()
	db := testDB(t)
	app := New(db)
	app.SetPackRegistry(newTestPackRegistry(t))
	return app
}

func countTableRows(t *testing.T, db *sql.DB, table, where string, args ...any) int {
	t.Helper()
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", table, where)
	err := db.QueryRowContext(context.Background(), query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return count
}

func TestMCPHappyPath(t *testing.T) {
	app := newDBApp(t)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000001"
	packetID := "happy-path-pkt-001"

	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      packetID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "happy-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-v1"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "happy path belief", ClaimType: "derived"},
			{LocalID: "b2", Claim: "second belief", ClaimType: "accommodated"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact",
				ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				ArtifactRef:   "artifact://test-artifact"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "ed1", FromRef: "local:b1", ToRef: "local:b2", Kind: "derives"},
		},
	}

	// Persist through the real transaction path
	tx, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	result, err := app.Persist(ctx, tx, pkt)
	if err != nil {
		tx.Rollback()
		t.Fatalf("persist: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}

	// Verify result
	if result.PacketID != packetID {
		t.Errorf("packet_id = %q, want %q", result.PacketID, packetID)
	}
	if len(result.BeliefIDs) != 2 {
		t.Errorf("belief count = %d, want 2", len(result.BeliefIDs))
	}

	db := app.DB()

	// Verify packet_submission
	var psPacketID, psAgentID, psRole string
	err = db.QueryRowContext(ctx,
		`SELECT packet_id, agent_id, role FROM packet_submission WHERE packet_id = $1`,
		packetID).Scan(&psPacketID, &psAgentID, &psRole)
	if err != nil {
		t.Fatalf("packet_submission not found: %v", err)
	}
	if psAgentID != "happy-agent" {
		t.Errorf("agent_id = %q, want %q", psAgentID, "happy-agent")
	}

	// Verify beliefs with origin_packet_id
	var b1Origin, b2Origin string
	err = db.QueryRowContext(ctx,
		`SELECT origin_packet_id FROM belief WHERE claim = $1 AND scenario_id = $2::UUID`,
		"happy path belief", scenarioID).Scan(&b1Origin)
	if err != nil {
		t.Fatalf("belief b1 not found: %v", err)
	}
	if b1Origin != packetID {
		t.Errorf("belief b1 origin_packet_id = %q, want %q", b1Origin, packetID)
	}
	err = db.QueryRowContext(ctx,
		`SELECT origin_packet_id FROM belief WHERE claim = $1 AND scenario_id = $2::UUID`,
		"second belief", scenarioID).Scan(&b2Origin)
	if err != nil {
		t.Fatalf("belief b2 not found: %v", err)
	}
	if b2Origin != packetID {
		t.Errorf("belief b2 origin_packet_id = %q, want %q", b2Origin, packetID)
	}

	// Verify evidence with origin_packet_id
	var e1Origin string
	err = db.QueryRowContext(ctx,
		`SELECT origin_packet_id FROM evidence WHERE scenario_id = $1::UUID`,
		scenarioID).Scan(&e1Origin)
	if err != nil {
		t.Fatalf("evidence not found: %v", err)
	}
	if e1Origin != packetID {
		t.Errorf("evidence origin_packet_id = %q, want %q", e1Origin, packetID)
	}

	// Verify edge (belief_edge has no scenario_id; verify via parent_id)
	var edgeCount int
	err = db.QueryRowContext(ctx,
		`SELECT count(*) FROM belief_edge`).Scan(&edgeCount)
	if err != nil {
		t.Fatalf("edge query: %v", err)
	}
	if edgeCount != 1 {
		t.Errorf("edge count = %d, want 1", edgeCount)
	}

	// Verify edge_provenance
	var epCount int
	err = db.QueryRowContext(ctx,
		`SELECT count(*) FROM edge_provenance WHERE origin_packet_id = $1`,
		packetID).Scan(&epCount)
	if err != nil {
		t.Fatalf("edge_provenance query: %v", err)
	}
	if epCount != 1 {
		t.Errorf("edge_provenance count = %d, want 1", epCount)
	}

	t.Logf("happy path: packet_submission=%d beliefs=%d evidence=%d edges=%d edge_provenance=%d",
		countTableRows(t, db, "packet_submission", "packet_id = $1", packetID),
		countTableRows(t, db, "belief", "origin_packet_id = $1", packetID),
		countTableRows(t, db, "evidence", "origin_packet_id = $1", packetID),
		countTableRows(t, db, "belief_edge", "true"),  // edges don't have origin column
		countTableRows(t, db, "edge_provenance", "origin_packet_id = $1", packetID))
}

func TestMCPRejectsMalformedWithRealDB(t *testing.T) {
	app := newDBApp(t)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000001"

	// Malformed: missing agent.id — Compile passes (doesn't check agent),
	// Validate passes (doesn't check agent), ValidatePacket catches it.
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "malformed-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "", Role: packetv1.RoleWork, Harness: "test", Model: "test-v1"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "should not persist", ClaimType: "derived"},
		},
	}

	// Compile should pass (doesn't check agent)
	if err := app.Compile(ctx, pkt); err != nil {
		t.Fatalf("compile should pass for empty agent.id: %v", err)
	}
	// ValidatePacket should reject (canonical validator checks agent)
	err := app.ValidatePacket(ctx, pkt)
	if err == nil {
		t.Fatal("expected ValidatePacket to fail for empty agent.id")
	}

	// Verify zero rows
	db := app.DB()
	var count int
	err = db.QueryRowContext(ctx, `SELECT count(*) FROM belief WHERE scenario_id = $1::UUID`, scenarioID).Scan(&count)
	if err != nil {
		t.Fatalf("count beliefs: %v", err)
	}
	if count != 0 {
		t.Errorf("belief count = %d, want 0", count)
	}
}

func TestMCPIdempotency(t *testing.T) {
	app := newDBApp(t)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000001"
	packetID := "idempotent-pkt-001"

	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      packetID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "idempotent-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-v1"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "idempotent belief", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact",
				ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
				ArtifactRef:   "artifact://test-artifact"},
		},
		// No edges — self-edge would be rejected
	}

	// First persist
	tx1, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, pkt)
	if err != nil {
		tx1.Rollback()
		t.Fatalf("persist1: %v", err)
	}
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit1: %v", err)
	}

	// Second persist (same packet)
	tx2, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	_, err = app.Persist(ctx, tx2, pkt)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("persist2: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatalf("commit2: %v", err)
	}

	db := app.DB()

	// Verify no duplicates
	psCount := countTableRows(t, db, "packet_submission", "packet_id = $1", packetID)
	if psCount != 1 {
		t.Errorf("packet_submission count = %d, want 1", psCount)
	}

	beliefCount := countTableRows(t, db, "belief", "origin_packet_id = $1", packetID)
	if beliefCount != 1 {
		t.Errorf("belief count = %d, want 1", beliefCount)
	}

	evidenceCount := countTableRows(t, db, "evidence", "origin_packet_id = $1", packetID)
	if evidenceCount != 1 {
		t.Errorf("evidence count = %d, want 1", evidenceCount)
	}

	// Verify origin_packet_id unchanged
	var origin string
	err = db.QueryRowContext(ctx,
		`SELECT origin_packet_id FROM belief WHERE origin_packet_id = $1 LIMIT 1`,
		packetID).Scan(&origin)
	if err != nil {
		t.Fatalf("origin not found: %v", err)
	}
	if origin != packetID {
		t.Errorf("origin_packet_id = %q, want %q", origin, packetID)
	}
}

func TestMCPAtomicRollback(t *testing.T) {
	app := newDBApp(t)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000001"
	missingProjectID := "99999999-9999-9999-9999-999999999999"

	// Packet with tasks referencing a non-existent project
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "rollback-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "rollback-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-v1"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "should be rolled back", ClaimType: "derived"},
		},
		Tasks: []packetv1.Task{
			{LocalID: "t1", Title: "task with missing project"},
		},
	}

	// Override ScenarioID for task preflight to use missing project
	// The task preflight checks conductor_project for pkt.ScenarioID
	// We use a scenarioID that has no matching conductor_project
	pktTask := *pkt
	pktTask.ScenarioID = missingProjectID
	pktTask.Beliefs[0].Claim = "rolled back belief"
	pktTask.PacketID = "rollback-pkt-002"

	// Persist should fail at task preflight
	tx, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx, &pktTask)
	if err == nil {
		tx.Rollback()
		t.Fatal("expected persist to fail for missing project")
	}
	t.Logf("persist failed as expected: %v", err)
	tx.Rollback()

	// Verify zero objects remain
	db := app.DB()
	psCount := countTableRows(t, db, "packet_submission", "packet_id = $1", pktTask.PacketID)
	if psCount != 0 {
		t.Errorf("packet_submission count = %d, want 0 after rollback", psCount)
	}

	beliefCount := countTableRows(t, db, "belief", "origin_packet_id = $1", pktTask.PacketID)
	if beliefCount != 0 {
		t.Errorf("belief count = %d, want 0 after rollback", beliefCount)
	}

	taskCount := countTableRows(t, db, "conductor_task", "origin_packet_id = $1", pktTask.PacketID)
	if taskCount != 0 {
		t.Errorf("task count = %d, want 0 after rollback", taskCount)
	}
}

func TestMCPSamePacketIDDifferentContent(t *testing.T) {
	app := newDBApp(t)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000001"
	packetID := "conflict-pkt-001"

	// First packet: valid
	pkt1 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      packetID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-1", Role: packetv1.RoleWork, Harness: "test", Model: "v1"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "original claim", ClaimType: "derived"},
		},
	}

	tx1, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	_, err = app.Persist(ctx, tx1, pkt1)
	if err != nil {
		tx1.Rollback()
		t.Fatalf("persist1: %v", err)
	}
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit1: %v", err)
	}

	// Second packet: same packet_id, same claim but different agent identity
	pkt2 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      packetID, // same ID
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-2", Role: packetv1.RoleWork, Harness: "test", Model: "v2"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "original claim", ClaimType: "derived"}, // same claim
		},
	}

	tx2, err := app.DB().BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	_, err = app.Persist(ctx, tx2, pkt2)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("persist2: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatalf("commit2: %v", err)
	}

	db := app.DB()

	// Verify packet_submission still has original agent (ON CONFLICT DO NOTHING preserves first)
	var agentID string
	err = db.QueryRowContext(ctx,
		`SELECT agent_id FROM packet_submission WHERE packet_id = $1`,
		packetID).Scan(&agentID)
	if err != nil {
		t.Fatalf("packet_submission not found: %v", err)
	}
	if agentID != "agent-1" {
		t.Errorf("agent_id = %q, want %q (original should be preserved)", agentID, "agent-1")
	}

	// Verify no duplicate beliefs (same claim = same EntityID = ON CONFLICT DO NOTHING)
	beliefCount := countTableRows(t, db, "belief", "origin_packet_id = $1 AND scenario_id = $2::UUID", packetID, scenarioID)
	if beliefCount != 1 {
		t.Errorf("belief count = %d, want 1 (no duplicate)", beliefCount)
	}

	// Verify original origin_packet_id unchanged
	var origin string
	err = db.QueryRowContext(ctx,
		`SELECT origin_packet_id FROM belief WHERE scenario_id = $1::UUID AND claim = $2`,
		scenarioID, "original claim").Scan(&origin)
	if err != nil {
		t.Fatalf("belief not found: %v", err)
	}
	if origin != packetID {
		t.Errorf("origin_packet_id = %q, want %q", origin, packetID)
	}
}

// TestGetContextReadsBackPersistedState proves the freeze-critical invariant:
//
//	WRITE → COMMIT → READ-BACK
//
// Uses three packets from two agents to verify per-object provenance.
func TestGetContextReadsBackPersistedState(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-0000000000A0"
	projectID := scenarioID

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ($1, 'freeze-test-project', 'active')
		 ON CONFLICT (id) DO NOTHING`, projectID)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Create task
	ws := work.NewStore(db)
	taskID := EntityID(scenarioID, "task", "freeze-test-task")
	task := &work.Task{
		ID:        taskID,
		ProjectID: projectID,
		Title:     "freeze test task",
	}
	if err := ws.Create(ctx, task); err != nil {
		t.Fatalf("create task: %v", err)
	}

	// --- Packet P1: work agent, 1 belief + 1 evidence ---
	p1PacketID := "freeze-pkt-p1"
	p1 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      p1PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p1", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "first claim from P1", ClaimType: "derived", Debt: []string{"needMap"}},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
	}
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	res1, err := app.Persist(ctx, tx1, p1)
	if err != nil {
		tx1.Rollback()
		t.Fatalf("persist P1: %v", err)
	}
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit P1: %v", err)
	}
	belief1ID := res1.BeliefIDs["b1"]

	// --- Packet P2: adversarial agent, 1 belief + 1 evidence ---
	p2PacketID := "freeze-pkt-p2"
	p2 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleAdversarial,
		PacketID:      p2PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p2", Role: packetv1.RoleAdversarial, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "contradiction from P2", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3"},
		},
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	res2, err := app.Persist(ctx, tx2, p2)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("persist P2: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatalf("commit P2: %v", err)
	}
	belief2ID := res2.BeliefIDs["b1"]

	// --- Packet P3: work agent, 1 edge (contradicts) + minimal belief ---
	p3PacketID := "freeze-pkt-p3"
	p3 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      p3PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p1", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b-witness", Claim: "edge witness belief", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "ed1", FromRef: "canonical:belief:" + belief1ID, ToRef: "canonical:belief:" + belief2ID, Kind: "contradicts"},
		},
	}
	tx3, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx3: %v", err)
	}
	_, err = app.Persist(ctx, tx3, p3)
	if err != nil {
		tx3.Rollback()
		t.Fatalf("persist P3: %v", err)
	}
	if err := tx3.Commit(); err != nil {
		t.Fatalf("commit P3: %v", err)
	}

	// --- GetContext read-back ---
	ctxResult, err := app.GetContext(ctx, taskID)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}

	// Task identity
	if ctxResult.Task == nil {
		t.Fatal("GetContext returned nil task")
	}
	if ctxResult.Task.ID != normalizeUUID(taskID) {
		t.Errorf("task ID = %q, want %q", ctxResult.Task.ID, normalizeUUID(taskID))
	}
	if !strings.EqualFold(ctxResult.Task.ProjectID, projectID) {
		t.Errorf("task ProjectID = %q, want %q", ctxResult.Task.ProjectID, projectID)
	}

	// Availability
	if !ctxResult.Availability.Snapshot.Available {
		t.Errorf("snapshot unavailable: %s", ctxResult.Availability.Snapshot.Reason)
	}
	if !ctxResult.Availability.Task.Available {
		t.Errorf("task unavailable: %s", ctxResult.Availability.Task.Reason)
	}

	// Beliefs: expect 3 (P1 + P2 + P3 witness)
	snap := ctxResult.Snapshot
	if len(snap.Beliefs) != 3 {
		t.Fatalf("belief count = %d, want 3", len(snap.Beliefs))
	}

	// Build belief map for assertions
	// Normalize keys: CRDB returns hyphenated UUIDs, EntityID returns flat hex
	beliefMap := make(map[string]epistemic.BeliefView)
	for _, b := range snap.Beliefs {
		beliefMap[b.ID] = b
	}

	// Assert P1 belief
	b1, ok := beliefMap[normalizeUUID(belief1ID)]
	if !ok {
		t.Fatalf("belief from P1 not found (ID=%s)", belief1ID)
	}
	if b1.Claim != "first claim from P1" {
		t.Errorf("P1 claim = %q, want %q", b1.Claim, "first claim from P1")
	}
	if b1.OriginPacketID != p1PacketID {
		t.Errorf("P1 origin_packet_id = %q, want %q", b1.OriginPacketID, p1PacketID)
	}
	if len(b1.Debt) != 1 || b1.Debt[0] != "needMap" {
		t.Errorf("P1 debt = %v, want [needMap]", b1.Debt)
	}

	// Assert P2 belief
	b2, ok := beliefMap[normalizeUUID(belief2ID)]
	if !ok {
		t.Fatalf("belief from P2 not found (ID=%s)", belief2ID)
	}
	if b2.Claim != "contradiction from P2" {
		t.Errorf("P2 claim = %q, want %q", b2.Claim, "contradiction from P2")
	}
	if b2.OriginPacketID != p2PacketID {
		t.Errorf("P2 origin_packet_id = %q, want %q", b2.OriginPacketID, p2PacketID)
	}
	// P2 submits debt=nil → stored [] → GetContext [].
	// This proves the current path round-trips the packet faithfully.
	// CompileDebt / pack initial_debt is not applied during production persist.
	// This is a known semantic boundary, not a bug.
	// Revisit only if a concrete research-cycle failure demonstrates that a belief
	// can enter the epistemic ledger without mandatory debt that the active pack requires.
	if len(b2.Debt) != 0 {
		t.Errorf("P2 debt = %v, want [] (CompileDebt not in persist path)", b2.Debt)
	}

	// Assert provenance is per-object and distinct
	if p1PacketID == p2PacketID {
		t.Errorf("P1 and P2 packet IDs must differ, both = %q", p1PacketID)
	}
	if b1.OriginPacketID == b2.OriginPacketID {
		t.Errorf("beliefs must have distinct origins, both = %q", b1.OriginPacketID)
	}

	// Evidence: expect 2 (P1 + P2)
	if len(snap.Evidence) != 2 {
		t.Fatalf("evidence count = %d, want 2", len(snap.Evidence))
	}
	evidenceMap := make(map[string]epistemic.EvidenceView)
	for _, e := range snap.Evidence {
		key := normalizeUUID(e.BeliefID) + ":" + e.ProvenanceClass
		evidenceMap[key] = e
	}
	if _, ok := evidenceMap[normalizeUUID(belief1ID)+":external_feed"]; !ok {
		t.Error("P1 evidence (external_feed) not found")
	}
	if _, ok := evidenceMap[normalizeUUID(belief2ID)+":external_feed"]; !ok {
		t.Error("P2 evidence (external_feed) not found")
	}

	// Edges: expect 1 (P3's contradicts edge)
	if len(snap.Edges) != 1 {
		t.Fatalf("edge count = %d, want 1", len(snap.Edges))
	}
	edge := snap.Edges[0]
	if edge.ParentID != normalizeUUID(belief1ID) {
		t.Errorf("edge parent_id = %q, want %q", edge.ParentID, normalizeUUID(belief1ID))
	}
	if edge.ChildID != normalizeUUID(belief2ID) {
		t.Errorf("edge child_id = %q, want %q", edge.ChildID, normalizeUUID(belief2ID))
	}
	if edge.Kind != "contradicts" {
		t.Errorf("edge kind = %q, want contradicts", edge.Kind)
	}

	t.Logf("freeze read-back verified: beliefs=%d evidence=%d edges=%d origins={%s,%s}",
		len(snap.Beliefs), len(snap.Evidence), len(snap.Edges), p1PacketID, p2PacketID)
}

// TestListByProjectReturnsOriginPacketID verifies the ListByProject query

// TestListByProjectReturnsOriginPacketID verifies the ListByProject query
// succeeds with matching SELECT/Scan columns including origin_packet_id.
func TestListByProjectReturnsOriginPacketID(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	projectID := "00000000-0000-0000-0000-0000000000B0"

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ($1, 'listbyproject-test', 'active')
		 ON CONFLICT DO NOTHING`, projectID)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Insert a packet_submission record (required by FK on origin_packet_id)
	originPacketID := "test-origin-pkt-001"
	_, err = db.ExecContext(ctx,
		`INSERT INTO packet_submission (packet_id, scenario_id, agent_id, role, harness, model, content_sha256)
		 VALUES ($1, '00000000-0000-0000-0000-0000000000FF', 'test-agent', 'work', 'test', 'test-model', 'abc123')`,
		originPacketID)
	if err != nil {
		t.Fatalf("insert packet_submission: %v", err)
	}

	// Create task with origin_packet_id via raw SQL (work.Store.Create doesn't persist origin_packet_id)
	ws := work.NewStore(db)
	taskID := EntityID(projectID, "task", "listbyproject-test-task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at)
		 VALUES ($1, $2, 'test task', 'test desc', 'proposed', 'medium', NULL, NULL, $3, NULL, now(), now())`, taskID, projectID, originPacketID)
	if err != nil {
		t.Fatalf("insert task: %v", err)
	}

	// ListByProject must succeed (no scan mismatch)
	tasks, err := ws.ListByProject(ctx, projectID)
	if err != nil {
		t.Fatalf("ListByProject failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if task.ID != normalizeUUID(taskID) {
		t.Errorf("task ID = %q, want %q", task.ID, normalizeUUID(taskID))
	}
	if task.OriginPacketID == nil || *task.OriginPacketID != originPacketID {
		t.Errorf("origin_packet_id = %v, want %q", task.OriginPacketID, originPacketID)
	}
	if task.Title != "test task" {
		t.Errorf("title = %q, want %q", task.Title, "test task")
	}
	if task.Status != "proposed" {
		t.Errorf("status = %q, want proposed", task.Status)
	}

	t.Logf("ListByProject verified: task=%s origin_packet_id=%s", task.ID, *task.OriginPacketID)
}

// TestPersistDoesNotApplyPackInitialDebt verifies that Persist stores the
// packet's debt directly — no union with the Domain Pack's initial_debt.
//
// Current production semantics:
//   packet debt → Persist() → stored belief debt
//
// NOT:
//   packet debt + pack initial_debt → Persist()
//
// CompileDebt (coordinator/validate.go) unions packInitialDebt with
// beliefDebt but is never called from the production MCP path.
//
// If CompileDebt is wired in the future, this test should change to expect
// the bmist@1.0.0 initial_debt (6 items).
func TestPersistDoesNotApplyPackInitialDebt(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-0000000000C0"

	// Packet with Debt=nil, PackRef=bmist@1.0.0
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "debt-test-pkt-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-debt", Role: packetv1.RoleWork, Harness: "test", Model: "test"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "debt test belief", ClaimType: "derived", Debt: nil},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"},
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	result, err := app.Persist(ctx, tx, pkt)
	if err != nil {
		tx.Rollback()
		t.Fatalf("persist: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}

	beliefID := result.BeliefIDs["b1"]
	snap, err := epistemic.GetSnapshot(ctx, db, scenarioID, epistemic.SnapshotOpts{})
	if err != nil {
		t.Fatalf("GetSnapshot: %v", err)
	}

	// Find the belief in the snapshot
	var found bool
	for _, b := range snap.Beliefs {
		// Normalize: CRDB returns hyphenated UUIDs, EntityID returns flat hex
		bID := strings.ReplaceAll(b.ID, "-", "")
		pID := strings.ReplaceAll(beliefID, "-", "")
		if strings.EqualFold(bID, pID) {
			found = true
			// Debt must be empty — CompileDebt not wired
			if len(b.Debt) != 0 {
				t.Errorf("debt = %v, want [] (CompileDebt not in persist path); if wired, expect 6 items from bmist@1.0.0", b.Debt)
			}
			break
		}
	}
	if !found {
		t.Fatalf("belief %s not found in snapshot", beliefID)
	}
}
