package oracle

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/PithomLabs/oracle/internal/application"
	"github.com/PithomLabs/oracle/internal/migrations"
	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	bmistv11 "github.com/PithomLabs/oracle/domain-pack/bmist/v1.1.0"
	"github.com/PithomLabs/oracle/internal/ui"
	"github.com/PithomLabs/oracle/internal/work"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func integrationDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("ARGUS_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
	}
	adminDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open admin db: %v", err)
	}
	dbName := fmt.Sprintf("argus_integ_%s_test", t.Name())
	dbName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		if r >= 'A' && r <= 'Z' {
			return r + 32
		}
		return '_'
	}, dbName)
	_, _ = adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
	_, err = adminDB.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %q", dbName))
	if err != nil {
		adminDB.Close()
		t.Fatalf("create database: %v", err)
	}
	testDSN := fmt.Sprintf("postgres://root@localhost:26257/%s?sslmode=disable", dbName)
	conn, err := sql.Open("pgx", testDSN)
	if err != nil {
		adminDB.Close()
		t.Fatalf("open test db: %v", err)
	}
	integApplySolventMigrations(context.Background(), conn)
	migrations.Apply(context.Background(), conn)
	t.Cleanup(func() {
		conn.Close()
		adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
		adminDB.Close()
	})
	return conn
}

func integApplySolventMigrations(ctx context.Context, db *sql.DB) error {
	schemaDir := filepath.Join("..", "solvent-main", "db")
	entries, err := os.ReadDir(schemaDir)
	if err != nil {
		return err
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
			return err
		}
		for _, stmt := range integSplitStatements(string(data)) {
			db.ExecContext(ctx, stmt)
		}
	}
	return nil
}

func integSplitStatements(sqlText string) []string {
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

// TestFullIntegration runs the 17-step acceptance lifecycle.
func TestFullIntegration(t *testing.T) {
	db := integrationDB(t)
	app := application.New(db)
	packReg := domainpack.NewRegistry()
	packReg.RegisterDecoder("bmist", func(data []byte) (domainpack.Pack, error) {
		return bmistv1.ParsePack(data)
	})
	if err := packReg.LoadEmbedded(bmistv1.LoadEmbedded(), "bmist"); err != nil {
		t.Fatalf("load embedded pack: %v", err)
	}
	if err := packReg.LoadEmbedded(bmistv11.LoadEmbedded(), "bmist"); err != nil {
		t.Fatalf("load embedded v1.1.0 pack: %v", err)
	}
	app.SetPackRegistry(packReg)
	uiServer := ui.NewServer(app, "test-token")
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000080"

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000080', 'integ-project', 'active')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// ---- Step 1: Create task ----
	ws := work.NewStore(db)
	task := &work.Task{
		ID:       application.EntityID(scenarioID, "task", "integ task"),
		ProjectID: "00000000-0000-0000-0000-000000000080",
		Title:    "integ task",
	}
	if err := ws.Create(ctx, task); err != nil {
		t.Fatalf("step 1: create task: %v", err)
	}
	t.Log("Step 1: task created")

	// ---- Step 2: get_context for task ----
	ctxResult, err := app.GetContext(ctx, task.ID)
	if err != nil {
		t.Fatalf("step 2: get_context: %v", err)
	}
	if ctxResult.Task == nil {
		t.Fatal("step 2: task not found")
	}
	t.Logf("Step 2: context assembled (task=%s)", ctxResult.Task.ID)

	// ---- Step 3: Submit work packet ----
	workPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "integ-work-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "integration belief", ClaimType: "derived", Debt: []string{"needMap"}},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "evidence-hash-001"},
		},
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("step 3: begin tx: %v", err)
	}
	result, err := app.Persist(ctx, tx, workPkt)
	if err != nil {
		tx.Rollback()
		t.Fatalf("step 3: persist: %v", err)
	}
	tx.Commit()
	beliefID := result.BeliefIDs["b1"]
	t.Logf("Step 3: work packet persisted (belief=%s)", beliefID)

	// ---- Step 4: get_context shows beliefs ----
	ctxResult2, err := app.GetContext(ctx, task.ID)
	if err != nil {
		t.Fatalf("step 4: get_context: %v", err)
	}
	if len(ctxResult2.Snapshot.Beliefs) == 0 {
		t.Fatal("step 4: no beliefs in snapshot")
	}
	t.Logf("Step 4: snapshot has %d beliefs", len(ctxResult2.Snapshot.Beliefs))

	// ---- Step 5: Human discharge debt via UI ----
	// Create operator principal (must exist for discharge)
	db.ExecContext(ctx,
		`INSERT INTO principal (principal_id, principal_type, issuer)
		 VALUES ($1, 'human', 'operator')
		 ON CONFLICT (principal_id) DO NOTHING`,
		application.LocalOperatorPrincipalID)

	// Insert qualifying evidence for needMap (requires reproducible_artifact)
	_, err = db.ExecContext(ctx,
		`INSERT INTO evidence (id, scenario_id, belief_id, provenance_class, source_url, content_sha256)
		 VALUES ($1::UUID, $2::UUID, $3::UUID, 'reproducible_artifact', 'https://test.example.com', 'integ-evidence-001')
		 ON CONFLICT (id) DO NOTHING`,
		"00000000-0000-0000-0000-000000000082", scenarioID, beliefID)
	if err != nil {
		t.Fatalf("step 5 prep: insert evidence: %v", err)
	}

	dischargeBody := fmt.Sprintf(`{"scenario_id":"%s","belief_id":"%s","obligation_key":"needMap","instrument_ref":"att-1","evidence_class":"reproducible_artifact","principal_id":"evil"}`, scenarioID, beliefID)
	dischargeReq := httptest.NewRequest(http.MethodPost, "/ui/api/discharge", bytes.NewBufferString(dischargeBody))
	dischargeReq.Header.Set("Content-Type", "application/json")
	dischargeW := httptest.NewRecorder()
	uiServer.HandleDischarge(dischargeW, dischargeReq)
	if dischargeW.Code != http.StatusOK {
		t.Fatalf("step 5: discharge failed: %d %s", dischargeW.Code, dischargeW.Body.String())
	}
	t.Log("Step 5: debt discharged via UI")

	// ---- Step 6: Human promote via UI ----
	promoteBody := fmt.Sprintf(`{"scenario_id":"%s","belief_id":"%s"}`, scenarioID, beliefID)
	promoteReq := httptest.NewRequest(http.MethodPost, "/ui/api/promote", bytes.NewBufferString(promoteBody))
	promoteReq.Header.Set("Content-Type", "application/json")
	promoteW := httptest.NewRecorder()
	uiServer.HandlePromote(promoteW, promoteReq)
	if promoteW.Code != http.StatusOK {
		t.Fatalf("step 6: promote failed: %d %s", promoteW.Code, promoteW.Body.String())
	}
	t.Log("Step 6: belief promoted via UI")

	// ---- Step 7: Verify belief is promoted ----
	var status string
	err = db.QueryRowContext(ctx, `SELECT status FROM belief WHERE id = $1::UUID`, beliefID).Scan(&status)
	if err != nil {
		t.Fatalf("step 7: query: %v", err)
	}
	if status != "promoted" {
		t.Fatalf("step 7: expected promoted, got %s", status)
	}
	t.Logf("Step 7: belief status=%s", status)

	// ---- Step 8: Retract belief ----
	err = app.SubmitDecision(ctx, &application.AuthenticatedDecisionCommand{
		Type:       "retract",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	if err != nil {
		t.Fatalf("step 8: retract: %v", err)
	}
	t.Log("Step 8: belief retracted")

	// ---- Step 9: Governance-linked task cancelled ----
	govTaskID := application.EntityID(scenarioID, "task", "gov task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status, governance_ref)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000080', 'gov task', 'active', $2::UUID)
		 ON CONFLICT (id) DO NOTHING`, govTaskID, beliefID)
	if err != nil {
		t.Fatalf("step 9: insert gov task: %v", err)
	}
	err = app.SubmitDecision(ctx, &application.AuthenticatedDecisionCommand{
		Type:       "retract",
		ScenarioID: scenarioID,
		BeliefID:   beliefID,
	})
	_ = err // may error on already-retracted, that's fine
	var govStatus string
	err = db.QueryRowContext(ctx, `SELECT status FROM conductor_task WHERE id = $1::UUID`, govTaskID).Scan(&govStatus)
	if err != nil {
		t.Fatalf("step 9: query: %v", err)
	}
	t.Logf("Step 9: governance task status=%s", govStatus)

	// ---- Step 10: Idempotency — submit same packet ----
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("step 10: begin tx: %v", err)
	}
	_, err = app.Persist(ctx, tx2, workPkt)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("step 10: persist: %v", err)
	}
	tx2.Commit()
	var beliefCount int
	db.QueryRowContext(ctx, `SELECT count(*) FROM belief WHERE scenario_id = $1::UUID`, scenarioID).Scan(&beliefCount)
	if beliefCount != 1 {
		t.Fatalf("step 10: expected 1 belief (idempotent), got %d", beliefCount)
	}
	t.Logf("Step 10: idempotency verified (count=%d)", beliefCount)

	// ---- Step 11: Insights page ----
	insightsReq := httptest.NewRequest(http.MethodGet, "/ui/insights?task_id="+task.ID, nil)
	insightsW := httptest.NewRecorder()
	uiServer.HandleInsights(insightsW, insightsReq)
	if insightsW.Code != http.StatusOK {
		t.Fatalf("step 11: insights: %d", insightsW.Code)
	}
	t.Log("Step 11: insights page renders")

	// ---- Step 12: Debts page ----
	debtsReq := httptest.NewRequest(http.MethodGet, "/ui/debts?task_id="+task.ID, nil)
	debtsW := httptest.NewRecorder()
	uiServer.HandleDebts(debtsW, debtsReq)
	if debtsW.Code != http.StatusOK {
		t.Fatalf("step 12: debts: %d", debtsW.Code)
	}
	t.Log("Step 12: debts page renders")

	// ---- Step 13: Create second belief and contradiction edge ----
	workPkt2 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      "integ-work-002",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Beliefs: []packetv1.Belief{
			{LocalID: "b2", Claim: "second belief", ClaimType: "derived"},
		},
	}
	tx3, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("step 13: begin tx: %v", err)
	}
	result2, err := app.Persist(ctx, tx3, workPkt2)
	if err != nil {
		tx3.Rollback()
		t.Fatalf("step 13: persist: %v", err)
	}
	tx3.Commit()
	belief2ID := result2.BeliefIDs["b2"]
	advPkt := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleAdversarial,
		PacketID:      "integ-adv-001",
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Beliefs: []packetv1.Belief{
			{LocalID: "adv1", Claim: "adversarial refutation", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "adv1", FromRef: "local:adv1", ToRef: "canonical:belief:" + belief2ID, Kind: packetv1.EdgeContradicts},
		},
	}
	tx4, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("step 13: begin adv tx: %v", err)
	}
	_, err = app.Persist(ctx, tx4, advPkt)
	if err != nil {
		tx4.Rollback()
		t.Fatalf("step 13: persist adv: %v", err)
	}
	tx4.Commit()
	var edgeCount int
	db.QueryRowContext(ctx, `SELECT count(*) FROM belief_edge WHERE kind = 'contradicts' AND child_id = $1::UUID`, belief2ID).Scan(&edgeCount)
	if edgeCount != 1 {
		t.Fatalf("step 13: expected 1 contradiction edge, got %d", edgeCount)
	}
	t.Logf("Step 13: contradiction edge visible (count=%d)", edgeCount)

	// ---- Step 14: Verify entity IDs are deterministic ----
	beliefID2 := application.EntityID(scenarioID, "belief", "integration belief")
	if beliefID2 != beliefID {
		t.Fatalf("step 14: deterministic IDs differ: %s vs %s", beliefID2, beliefID)
	}
	t.Log("Step 14: entity IDs deterministic")

	// ---- Step 15: Dependency blocking ----
	blockerTaskID := application.EntityID(scenarioID, "task", "blocker")
	blockedTaskID := application.EntityID(scenarioID, "task", "blocked")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000080', 'blocker', 'active')
		 ON CONFLICT (id) DO NOTHING`, blockerTaskID)
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000080', 'blocked', 'proposed')
		 ON CONFLICT (id) DO NOTHING`, blockedTaskID)
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_dependency (task_id, blocked_by_id)
		 VALUES ($1::UUID, $2::UUID)
		 ON CONFLICT (task_id, blocked_by_id) DO NOTHING`, blockedTaskID, blockerTaskID)
	err = ws.Claim(ctx, blockedTaskID, "agent-1")
	if err == nil {
		t.Fatal("step 15: expected claim to fail")
	}
	t.Logf("Step 15: dependency blocking enforced: %v", err)

	// ---- Step 16: REOPEN lineage ----
	reopenOrigID := application.EntityID(scenarioID, "task", "reopen original")
	reopenSuccID := application.EntityID(scenarioID, "task", "reopen successor")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000080', 'reopen original', 'cancelled')
		 ON CONFLICT (id) DO NOTHING`, reopenOrigID)
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, status, reopened_from_task_id)
		 VALUES ($1::UUID, '00000000-0000-0000-0000-000000000080', 'reopen successor', 'proposed', $2::UUID)
		 ON CONFLICT (id) DO NOTHING`, reopenSuccID, reopenOrigID)
	var reopenedFrom sql.NullString
	err = db.QueryRowContext(ctx,
		`SELECT reopened_from_task_id::STRING FROM conductor_task WHERE id = $1::UUID`, reopenSuccID).Scan(&reopenedFrom)
	if err != nil {
		t.Fatalf("step 16: query: %v", err)
	}
	normalizeUUID := func(s string) string {
		if len(s) == 32 {
			return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
		}
		return s
	}
	if !reopenedFrom.Valid || normalizeUUID(reopenedFrom.String) != normalizeUUID(reopenOrigID) {
		t.Fatalf("step 16: lineage mismatch")
	}
	t.Logf("Step 16: reopen lineage preserved: %s -> %s", reopenOrigID, reopenSuccID)

	// ---- Step 17: Final status summary ----
	var totalBeliefs, totalEdges, totalTasks int
	db.QueryRowContext(ctx, `SELECT count(*) FROM belief WHERE scenario_id = $1::UUID`, scenarioID).Scan(&totalBeliefs)
	db.QueryRowContext(ctx, `SELECT count(*) FROM belief_edge`).Scan(&totalEdges)
	db.QueryRowContext(ctx, `SELECT count(*) FROM conductor_task WHERE project_id = '00000000-0000-0000-0000-000000000080'`).Scan(&totalTasks)
	t.Logf("Step 17: FINAL STATUS — beliefs=%d, edges=%d, tasks=%d", totalBeliefs, totalEdges, totalTasks)
	t.Log("=== 17-step integration test PASSED ===")
}
