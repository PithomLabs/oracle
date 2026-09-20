package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/PithomLabs/oracle/reference-loop/agent"
	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
)

func TestPositiveExecutionBoundary(t *testing.T) {
	tmpDir := t.TempDir()
	conductorDBPath := filepath.Join(tmpDir, "conductor.db")
	solventDBPath := "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable"
	apiKey := "test-api-key"
	runID := "test-run-pos-" + time.Now().Format("20060102150405")
	scenarioID := uuid.New().String()
	principalID := testPrincipalID

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resetExecutorCallCount()
	defer resetExecutorCallCount()

	conductorDB, err := initConductorDB(conductorDBPath)
	if err != nil {
		t.Fatalf("init conductor db: %v", err)
	}
	defer conductorDB.Close()

	projectID := "reference-loop-project"
	if err := createProject(ctx, conductorDB, projectID); err != nil {
		t.Fatalf("create project: %v", err)
	}

	solventDB, err := initSolventDB(solventDBPath)
	if err != nil {
		t.Fatalf("init solvent db: %v", err)
	}
	defer solventDB.Close()

	if err := createTestPrincipal(ctx, solventDB); err != nil {
		t.Fatalf("create test principal: %v", err)
	}

	ts, execReg, err := startSolventServer(ctx, solventDB)
	if err != nil {
		t.Fatalf("start solvent server: %v", err)
	}
	defer ts.Close()

	solventAddr := ts.Addr()

	conductorClient, err := conductor.NewClient(conductorDBPath, "reference-loop-agent")
	if err != nil {
		t.Fatalf("create conductor client: %v", err)
	}
	defer conductorClient.Close()

	solventClient := solvent.NewClient(solventAddr, apiKey)
	solventClient.DB = solventDB

	evidenceCollector := agent.NewEvidenceCollector(runID, scenarioID)

	refAgent := agent.NewAgent(
		conductorClient,
		solventClient,
		execReg,
		"reference-loop-agent",
		"reference-loop-project",
		scenarioID,
		runID,
		principalID,
		evidenceCollector,
	)

	if err := refAgent.RunHappyPath(ctx); err != nil {
		t.Fatalf("happy path failed: %v", err)
	}

	records := evidenceCollector.GetRecords()

	authorizedCount := 0
	executionAttemptedCount := 0
	effectConfirmedCount := 0
	expectedOpID := "deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:" + runID
	var authorizedIntentID string
	var executionIntentID string

	for _, r := range records {
		switch r.Event {
		case "AUTHORIZED":
			authorizedCount++
			if r.OperationID != expectedOpID {
				t.Fatalf("AUTHORIZED operation_id mismatch: got %s, want %s", r.OperationID, expectedOpID)
			}
			authorizedIntentID = r.IntentID
		case "EXECUTION_ATTEMPTED":
			executionAttemptedCount++
			if r.OperationID != expectedOpID {
				t.Fatalf("EXECUTION_ATTEMPTED operation_id mismatch: got %s, want %s", r.OperationID, expectedOpID)
			}
			executionIntentID = r.IntentID
		case "EFFECT_CONFIRMED":
			effectConfirmedCount++
			if r.OperationID != expectedOpID {
				t.Fatalf("EFFECT_CONFIRMED operation_id mismatch: got %s, want %s", r.OperationID, expectedOpID)
			}
		}
	}

	if authorizedCount != 1 {
		t.Fatalf("expected exactly 1 AUTHORIZED event, got %d", authorizedCount)
	}
	if executionAttemptedCount != 1 {
		t.Fatalf("expected exactly 1 EXECUTION_ATTEMPTED event, got %d", executionAttemptedCount)
	}
	if effectConfirmedCount != 1 {
		t.Fatalf("expected exactly 1 EFFECT_CONFIRMED event, got %d", effectConfirmedCount)
	}

	if authorizedIntentID == "" {
		t.Fatal("AUTHORIZED event should have intent_id")
	}
	if executionIntentID == "" {
		t.Fatal("EXECUTION_ATTEMPTED event should have intent_id")
	}
	if authorizedIntentID != executionIntentID {
		t.Fatalf("intent_id mismatch between AUTHORIZED and EXECUTION_ATTEMPTED: %s != %s", authorizedIntentID, executionIntentID)
	}

	executorCalls := getExecutorCallCount()
	if os.Getenv("GITHUB_TOKEN") == "" && executorCalls != 1 {
		t.Fatalf("expected executor call count == 1 (simulated), got %d", executorCalls)
	}

	t.Logf("Positive execution boundary proven:")
	t.Logf("  Operation ID: %s", expectedOpID)
	t.Logf("  Intent ID: %s", authorizedIntentID)
	t.Logf("  AUTHORIZED events: %d", authorizedCount)
	t.Logf("  EXECUTION_ATTEMPTED events: %d", executionAttemptedCount)
	t.Logf("  EFFECT_CONFIRMED events: %d", effectConfirmedCount)
	t.Logf("  Executor invocation count: %d (real GitHub: %v)", executorCalls, os.Getenv("GITHUB_TOKEN") != "")
	t.Logf("  Canonical operation identity preserved: yes")
}
