package ui

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

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	bmistv11 "github.com/PithomLabs/oracle/domain-pack/bmist/v1.1.0"
	"github.com/PithomLabs/oracle/internal/application"
	"github.com/PithomLabs/oracle/internal/migrations"
	solventmigrations "github.com/PithomLabs/oracle/internal/solventmigrations"
	"github.com/PithomLabs/oracle/internal/work"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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
	dbName := fmt.Sprintf("argus_ui_%s_test", t.Name())
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
	solventmigrations.Apply(context.Background(), conn)
	migrations.Apply(context.Background(), conn)
	t.Cleanup(func() {
		conn.Close()
		adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
		adminDB.Close()
	})
	return conn
}

func applySolventMigrations(ctx context.Context, db *sql.DB) error {
	schemaDir := filepath.Join("..", "..", "..", "solvent-main", "db")
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
		for _, stmt := range splitStatements(string(data)) {
			db.ExecContext(ctx, stmt)
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

func setupTestServer(t *testing.T) (*Server, *sql.DB) {
	t.Helper()
	db := testDB(t)
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
	server := NewServer(app, "test-token")
	return server, db
}

// ---- Acceptance Test 6: Human discharge debt via UI ----

func TestDischargeEndpoint(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Create the operator principal (must match handler's hardcoded UUID)
	_, err := db.ExecContext(ctx,
		`INSERT INTO principal (principal_id, principal_type, issuer)
		 VALUES ($1, 'human', 'operator')
		 ON CONFLICT (principal_id) DO NOTHING`,
		application.LocalOperatorPrincipalID)
	if err != nil {
		t.Fatalf("create principal: %v", err)
	}

	// Create a belief with debt
	scenarioID := "00000000-0000-0000-0000-000000000020"
	beliefID := "00000000-0000-0000-0000-000000000021"
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, status)
		 VALUES ($1::UUID, $2::UUID, 'discharge test', 'derived', '{"needMap"}', 'entered')
		 ON CONFLICT (id) DO NOTHING`, beliefID, scenarioID)
	if err != nil {
		t.Fatalf("insert belief: %v", err)
	}

	// Insert qualifying evidence (needMap requires reproducible_artifact)
	_, err = db.ExecContext(ctx,
		`INSERT INTO evidence (id, scenario_id, belief_id, provenance_class, source_url, content_sha256)
		 VALUES ($1::UUID, $2::UUID, $3::UUID, 'reproducible_artifact', 'https://test.example.com', 'test-sha-001')
		 ON CONFLICT (id) DO NOTHING`,
		"00000000-0000-0000-0000-000000000022", scenarioID, beliefID)
	if err != nil {
		t.Fatalf("insert evidence: %v", err)
	}

	// Build discharge request — body principal_id is ignored by server (uses canonical constant)
	body := fmt.Sprintf(`{"scenario_id":"%s","belief_id":"%s","obligation_key":"needMap","instrument_ref":"attestation-1","evidence_class":"reproducible_artifact","principal_id":"evil"}`, scenarioID, beliefID)
	req := httptest.NewRequest(http.MethodPost, "/ui/api/discharge", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.HandleDischarge(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	t.Log("discharge endpoint returned 200 OK")
}

// ---- Acceptance Test: Promote endpoint ----

func TestPromoteEndpoint(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Create a belief with no debt
	scenarioID := "00000000-0000-0000-0000-000000000022"
	beliefID := "00000000-0000-0000-0000-000000000023"
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, status)
		 VALUES ($1::UUID, $2::UUID, 'promote test', 'derived', '{}', 'entered')
		 ON CONFLICT (id) DO NOTHING`, beliefID, scenarioID)
	if err != nil {
		t.Fatalf("insert belief: %v", err)
	}

	// Build promote request
	body := fmt.Sprintf(`{"scenario_id":"%s","belief_id":"%s"}`, scenarioID, beliefID)
	req := httptest.NewRequest(http.MethodPost, "/ui/api/promote", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.HandlePromote(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify belief was promoted
	var status string
	err = db.QueryRowContext(ctx,
		`SELECT status FROM belief WHERE id = $1::UUID`, beliefID).Scan(&status)
	if err != nil {
		t.Fatalf("query belief: %v", err)
	}
	if status != "promoted" {
		t.Fatalf("expected status 'promoted', got '%s'", status)
	}
	t.Logf("promote endpoint succeeded: status=%s", status)
}

// ---- Acceptance Test: Insights page renders ----

func TestInsightsPage(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Create a project, task, and belief for dashboard
	db.ExecContext(ctx, `INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000050', 'test', 'active') ON CONFLICT (id) DO NOTHING`)
	taskID := application.EntityID("00000000-0000-0000-0000-000000000051", "task", "test")
	ws := work.NewStore(db)
	ws.Create(ctx, &work.Task{ID: taskID, ProjectID: "00000000-0000-0000-0000-000000000050", Title: "test-insight"})

	// Insert a belief for the dashboard
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, status)
		 VALUES ('00000000-0000-0000-0000-000000000060', '00000000-0000-0000-0000-000000000050', 'test belief claim', 'derived', '{}', 'entered')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert belief: %v", err)
	}

	// Dashboard: no task_id param
	req := httptest.NewRequest(http.MethodGet, "/ui/insights", nil)
	w := httptest.NewRecorder()

	server.HandleInsights(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "Insights") {
		t.Fatal("insights page missing expected content")
	}
	if !strings.Contains(body, "test-insight") {
		t.Fatal("insights page missing seeded task")
	}
	if !strings.Contains(body, "test belief claim") {
		t.Fatal("insights page missing seeded belief")
	}
	t.Log("insights page renders correctly with dashboard")
}

// ---- Acceptance Test: Debts page renders ----

func TestDebtsPage(t *testing.T) {
	server, db := setupTestServer(t)
	ctx := context.Background()

	// Create a project, task, and belief with debt for dashboard
	db.ExecContext(ctx, `INSERT INTO conductor_project (id, name, status) VALUES ('00000000-0000-0000-0000-000000000052', 'test', 'active') ON CONFLICT (id) DO NOTHING`)
	taskID := application.EntityID("00000000-0000-0000-0000-000000000053", "task", "test")
	ws := work.NewStore(db)
	ws.Create(ctx, &work.Task{ID: taskID, ProjectID: "00000000-0000-0000-0000-000000000052", Title: "test"})

	// Insert a belief with debt
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, status)
		 VALUES ('00000000-0000-0000-0000-000000000061', '00000000-0000-0000-0000-000000000052', 'debt test claim', 'derived', '{"needMap"}', 'entered')
		 ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		t.Fatalf("insert belief: %v", err)
	}

	// Dashboard: no task_id param
	req := httptest.NewRequest(http.MethodGet, "/ui/debts", nil)
	w := httptest.NewRecorder()

	server.HandleDebts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	if !strings.Contains(body, "Debt") {
		t.Fatal("debts page missing expected content")
	}
	if !strings.Contains(body, "debt test claim") {
		t.Fatal("debts page missing seeded belief")
	}
	t.Log("debts page renders correctly with dashboard")
}

func TestNewServerEmptyTokenDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("NewServer panicked with empty token: %v", r)
		}
	}()

	conn := testDB(t)
	defer conn.Close()

	app := application.New(conn)
	// Empty token should NOT panic — just use default
	srv := NewServer(app, "")
	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
}

func TestNewServerWithToken(t *testing.T) {
	conn := testDB(t)
	defer conn.Close()

	app := application.New(conn)
	srv := NewServer(app, "test-token-123")
	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
}
