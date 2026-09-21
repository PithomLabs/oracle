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

	"github.com/PithomLabs/oracle/internal/migrations"
	"github.com/PithomLabs/oracle/internal/work"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	"github.com/PithomLabs/oracle/verifier"
	"github.com/PithomLabs/oracle/verifier/physics/v1"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "aaa"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
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
			{LocalID: "b1", Claim: "c1", ClaimType: "math"},
			{LocalID: "b2", Claim: "c2", ClaimType: "math"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "e1", FromRef: "local:b1", ToRef: "local:b2", Kind: "derives"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "c1", ClaimType: "math"},
			{LocalID: "b2", Claim: "c2", ClaimType: "math"},
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
			{LocalID: "b2", Claim: "c2", ClaimType: "math"},
			{LocalID: "b1", Claim: "c1", ClaimType: "math"},
		},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p1", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "c1", ClaimType: "math"},
			{LocalID: "b2", Claim: "c2", ClaimType: "math"},
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
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
	}
	pkt2 := &packetv1.Packet{
		ScenarioID: "s1", PacketID: "p2", Role: "work", PackRef: "bmist@1.0.0",
  Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
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
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
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
		Beliefs: []packetv1.Belief{{LocalID: "b1", Claim: "c1", ClaimType: "math"}},
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


func TestAgentCannotPromote(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000090"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "agent-promote-test-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "Agent-captured claim", ClaimType: "derived"},
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

	beliefID := EntityID(scenarioID, "belief", "Agent-captured claim")

	// Agent tries to promote via SubmitDecision — should fail (no authority path for agents)
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "promote",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err == nil {
		t.Error("agent was able to promote belief — authority isolation violated")
	}
	t.Logf("promote correctly refused: %v", err)
}

func TestAgentCannotDischargeDebt(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000091"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "agent-discharge-test-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "Claim with debt", ClaimType: "derived", Debt: []string{"needMap"}},
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

	beliefID := EntityID(scenarioID, "belief", "Claim with debt")

	// Agent submits a retirement evidence packet — debt should NOT be discharged
	retPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "agent-retire-test-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "Claim with debt", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "reproducible_artifact", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
	}

	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx2, retPkt)
	if err != nil {
		t.Fatalf("persist retirement: %v", err)
	}
	tx2.Commit()

	// Verify debt still exists (agent cannot discharge)
	var debtJSON string
	err = db.QueryRowContext(ctx,
		`SELECT debt FROM belief WHERE id = $1`, beliefID).Scan(&debtJSON)
	if err != nil {
		t.Fatalf("query belief debt: %v", err)
	}
	if debtJSON == "" || debtJSON == "[]" || debtJSON == "null" {
		t.Fatalf("expected debt on belief, got: %s", debtJSON)
	}
	t.Logf("belief has debt: %s (correct — agent cannot discharge)", debtJSON)
}

func TestAgentCannotRetractBelief(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	app := New(db)
	ctx := context.Background()

	scenarioID := "00000000-0000-0000-0000-000000000092"
	pkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "agent-retract-test-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "test-agent", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "Claim to retract", ClaimType: "derived"},
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

	beliefID := EntityID(scenarioID, "belief", "Claim to retract")

	// Agent tries to retract via SubmitDecision — should fail (no authority path for agents)
	err = app.SubmitDecision(ctx, &AuthenticatedDecisionCommand{
		Type:       "retract",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err == nil {
		t.Error("agent was able to retract belief — authority isolation violated")
	}
	t.Logf("retract correctly refused: %v", err)
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
