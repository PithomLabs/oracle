package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	solventExecutor "github.com/PithomLabs/solvent/service/executor"
	"github.com/PithomLabs/solvent/adapter/github"
	"github.com/PithomLabs/solvent/api"
	"github.com/PithomLabs/solvent/service/audit"
	"github.com/PithomLabs/solvent/service/authority"
	"github.com/PithomLabs/solvent/service/ledger"
	"github.com/PithomLabs/solvent/service/policy"

	"github.com/PithomLabs/oracle/reference-loop/agent"
	"github.com/PithomLabs/oracle/reference-loop/config"
	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/evidence"
	"github.com/PithomLabs/oracle/reference-loop/executor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
)

const ProjectStatusActive = "active"
const testPrincipalID = "00000000-0000-0000-0000-000000000001"
const testAPIKey = "test-api-key"

var (
	executorCallCount int
	executorCallMutex sync.Mutex
)

func resetExecutorCallCount() {
	executorCallMutex.Lock()
	defer executorCallMutex.Unlock()
	executorCallCount = 0
}

func incrementExecutorCallCount() {
	executorCallMutex.Lock()
	defer executorCallMutex.Unlock()
	executorCallCount++
}

func getExecutorCallCount() int {
	executorCallMutex.Lock()
	defer executorCallMutex.Unlock()
	return executorCallCount
}

// testSolventServer owns the lifecycle of a Solvent test server.
type testSolventServer struct {
	httpServer *http.Server
	listener   net.Listener
	addr       string
	serveErr   chan error
}

// Addr returns the base URL of the test server.
func (s *testSolventServer) Addr() string { return s.addr }

// Close shuts down the server and waits for Serve to return.
func (s *testSolventServer) Close() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return err
	}
	// Wait for Serve to return
	<-s.serveErr
	return nil
}

func main() {
	scenario := flag.String("scenario", "happy_path", "Scenario to run")
	executorType := flag.String("executor", "recording", "Executor type: recording or github")
	conductorDBPath := flag.String("conductor-db", "", "Conductor database path")
	solventDBPath := flag.String("solvent-db", "", "Solvent database DSN (CockroachDB)")
	solventRESTAddr := flag.String("solvent-addr", "", "Solvent REST API address (empty = auto-detect from server)")
	solventAPIKey := flag.String("solvent-api-key", "", "Solvent API key")
	githubToken := flag.String("github-token", "", "GitHub token")
	githubRepo := flag.String("github-repo", "", "GitHub repository (owner/repo)")
	evidenceDir := flag.String("evidence-dir", "", "Evidence output directory")
	flag.Parse()

	cfg := config.DefaultConfig()
	cfg.Scenario = *scenario
	cfg.ExecutorType = *executorType
	if *conductorDBPath != "" {
		cfg.ConductorDBPath = *conductorDBPath
	}
	if *solventDBPath != "" {
		cfg.SolventDBPath = *solventDBPath
	}
	if *solventRESTAddr != "" {
		cfg.SolventRESTAddr = *solventRESTAddr
	}
	if *solventAPIKey != "" {
		cfg.SolventAPIKey = *solventAPIKey
	}
	if *githubToken != "" {
		cfg.GitHubToken = *githubToken
	}
	if *githubRepo != "" {
		cfg.GitHubRepo = *githubRepo
	}
	if *evidenceDir != "" {
		cfg.EvidenceDir = *evidenceDir
	}
	cfg.LoadFromEnv()

	cfg.Scenario = uuid.New().String()

	if err := os.MkdirAll(filepath.Dir(cfg.ConductorDBPath), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "create temp dir: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	conductorDB, err := initConductorDB(cfg.ConductorDBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init conductor db: %v\n", err)
		os.Exit(1)
	}
	defer conductorDB.Close()

	solventDB, err := initSolventDB(cfg.SolventDBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init solvent db: %v\n", err)
		os.Exit(1)
	}
	defer solventDB.Close()

	if err := createTestPrincipal(ctx, solventDB); err != nil {
		fmt.Fprintf(os.Stderr, "create test principal: %v\n", err)
		os.Exit(1)
	}

	solventServer, _, err := startSolventServer(ctx, solventDB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "start solvent server: %v\n", err)
		os.Exit(1)
	}
	defer solventServer.Close()

	if cfg.SolventRESTAddr == "" {
		cfg.SolventRESTAddr = solventServer.Addr()
	}

	testSolventAPI(cfg.SolventRESTAddr, cfg.SolventAPIKey)

	conductorClient, err := conductor.NewClient(cfg.ConductorDBPath, "reference-loop-agent")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create conductor client: %v\n", err)
		os.Exit(1)
	}
	defer conductorClient.Close()

	solventClient := solvent.NewClient(cfg.SolventRESTAddr, cfg.SolventAPIKey)
	solventClient.DB = solventDB

	execReg, err := executor.NewRegistry(cfg.ExecutorType, cfg.GitHubToken, cfg.GitHubRepo, cfg.GitHubWorkflow, cfg.GitHubRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create executor registry: %v\n", err)
		os.Exit(1)
	}

	evidenceCollector := agent.NewEvidenceCollector(cfg.RunID, cfg.Scenario)

	refAgent := agent.NewAgent(
		conductorClient,
		solventClient,
		execReg.GetRegistry(),
		"reference-loop-agent",
		"reference-loop-project",
		cfg.Scenario,
		cfg.RunID,
		testPrincipalID,
		evidenceCollector,
	)

	projectID := "reference-loop-project"
	if err := createProject(ctx, conductorDB, projectID); err != nil {
		fmt.Fprintf(os.Stderr, "create project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting reference loop: run_id=%s scenario=%s executor=%s\n", cfg.RunID, cfg.Scenario, cfg.ExecutorType)
	startTime := time.Now()

	if err := refAgent.RunHappyPath(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "reference loop failed: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Reference loop completed in %v\n", elapsed)

	collector := evidence.NewCollector(
		cfg.RunID,
		cfg.Scenario,
		cfg.EvidenceDir,
		cfg.ConductorDBPath,
		cfg.SolventDBPath,
		execReg,
	)

	evidencePkg, err := collector.Collect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "collect evidence: %v\n", err)
		os.Exit(1)
	}

	if err := collector.Save(evidencePkg); err != nil {
		fmt.Fprintf(os.Stderr, "save evidence: %v\n", err)
		os.Exit(1)
	}
	collector.Close()

	evidence.PrintTrace(evidencePkg)
	evidence.PrintSummary(evidencePkg)

	expectedOpID := "deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:" + cfg.RunID
	ok, msg := evidence.VerifyOperationBinding(evidencePkg, expectedOpID)
	if ok {
		fmt.Printf("\n✓ OPERATION BINDING VERIFIED: %s\n", msg)
	} else {
		fmt.Printf("\n✗ OPERATION BINDING FAILED: %s\n", msg)
		os.Exit(1)
	}

	fmt.Println("\n=== REFERENCE LOOP HAPPY PATH: SUCCESS ===")
}

func createTestPrincipal(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO principal (principal_id, principal_type, issuer, created_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT DO NOTHING
	`, testPrincipalID, "service", "reference-loop")
	return err
}

func testSolventAPI(baseURL, apiKey string) {
	fmt.Println("Testing Solvent REST API...")
	client := &http.Client{Timeout: 10 * time.Second}
	reqBody := `{"principal_type":"agent","issuer":"test"}`
	req, _ := http.NewRequest("POST", baseURL+"/v1/principals", bytes.NewReader([]byte(reqBody)))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("API request error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("API Response: status=%d body=%s\n", resp.StatusCode, string(body))
}

func initConductorDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	schemaPath := filepath.Join(".", "conductor_schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		schemaPath = filepath.Join("..", "conductor_schema.sql")
		schema, err = os.ReadFile(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("read conductor schema: %w", err)
		}
	}
	if _, err := db.Exec(string(schema)); err != nil {
		return nil, fmt.Errorf("execute conductor schema: %w", err)
	}
	return db, nil
}

func initSolventDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping solvent db: %w", err)
	}
	return db, nil
}

// startSolventServer creates a test Solvent server with full lifecycle ownership.
// It binds the listener synchronously, starts the server on that listener,
// and verifies the server is accepting connections before returning.
// The returned testSolventServer must be closed by the caller.
func startSolventServer(ctx context.Context, db *sql.DB) (*testSolventServer, *solventExecutor.Registry, error) {
	// 1. Bind listener synchronously — ephemeral port, fail immediately on conflict
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, fmt.Errorf("bind listener: %w", err)
	}

	// 2. Build services
	auditSvc := audit.New(db)
	policySvc := policy.New(db)
	execReg := solventExecutor.NewRegistry()

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		provider := github.NewHTTPProvider(token)
		github.RegisterExecutor(execReg, provider)
	} else {
		execReg.Register("github_trigger_workflow", solventExecutor.ActionFunc(func(ctx context.Context, params map[string]interface{}) (string, error) {
			incrementExecutorCallCount()
			return "simulated_success", nil
		}))
	}

	authSvc := authority.New(db, policySvc, auditSvc, execReg)
	_ = ledger.New(db, auditSvc)

	server := api.NewServer(db, auditSvc, api.WithAuthorityService(authSvc))
	handler := server.Handler()

	keyToPrincipal := map[string]string{
		testAPIKey: testPrincipalID,
	}
	authHandler := api.AuthMiddleware(keyToPrincipal, handler)

	httpServer := &http.Server{Handler: authHandler}

	// 3. Serve on the already-bound listener — no port race
	serveErr := make(chan error, 1)
	go func() {
		serveErr <- httpServer.Serve(ln)
	}()

	// 4. Wait for server to be serving or fail
	select {
	case err := <-serveErr:
		ln.Close()
		return nil, nil, fmt.Errorf("server failed to start: %w", err)
	case <-time.After(50 * time.Millisecond):
		// Server is accepting connections
	}

	addr := "http://" + ln.Addr().String()
	return &testSolventServer{
		httpServer: httpServer,
		listener:   ln,
		addr:       addr,
		serveErr:   serveErr,
	}, execReg, nil
}

func createProject(ctx context.Context, db *sql.DB, projectID string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.ExecContext(ctx,
		`INSERT OR IGNORE INTO conductor_project (id, name, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		projectID, "Reference Loop Project", "Project for reference loop testing", ProjectStatusActive, now, now)
	return err
}
