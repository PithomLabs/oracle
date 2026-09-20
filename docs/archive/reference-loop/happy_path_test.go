package main

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/PithomLabs/oracle/reference-loop/agent"
	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/executor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
)

func TestHappyPathEvidenceSequence(t *testing.T) {
	tmpDir := t.TempDir()
	conductorDBPath := filepath.Join(tmpDir, "conductor.db")
	solventDBPath := "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable"
	apiKey := "test-api-key"
	runID := "test-run-" + time.Now().Format("20060102150405")
	scenarioID := uuid.New().String()
	principalID := testPrincipalID

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

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

	ts, _, err := startSolventServer(ctx, solventDB)
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

	execReg, err := executor.NewRegistry("recording", "", "ibmendoza/reference-loop-test", "ref-loop.yml", "main")
	if err != nil {
		t.Fatalf("create executor registry: %v", err)
	}

	evidenceCollector := agent.NewEvidenceCollector(runID, scenarioID)

	refAgent := agent.NewAgent(
		conductorClient,
		solventClient,
		execReg.GetRegistry(),
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
	events := []string{"TASK_CREATED", "TASK_CLAIMED", "PROPOSAL_CREATED", "AUTHORIZED", "EXECUTION_ATTEMPTED", "EFFECT_CONFIRMED", "RESULT_OBSERVED", "CONDUCTOR_UPDATED", "AGENT_CONTINUED"}
	for _, event := range events {
		found := false
		for _, r := range records {
			if r.Event == event {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing event: %s", event)
		}
	}

	data, _ := json.MarshalIndent(records, "", "  ")
	_ = data

	t.Logf("Evidence records: %d", len(records))
	for _, r := range records {
		t.Logf("  [%s] %s | owner=%s | source=%s | task=%s | intent=%s", r.Timestamp.Format("15:04:05"), r.Event, r.Owner, r.Source, r.TaskID, r.IntentID)
	}
}

func TestHappyPathOperationBinding(t *testing.T) {
	tmpDir := t.TempDir()
	conductorDBPath := filepath.Join(tmpDir, "conductor.db")
	solventDBPath := "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable"
	apiKey := "test-api-key"
	runID := "test-run-binding-" + time.Now().Format("20060102150405")
	scenarioID := uuid.New().String()
	principalID := testPrincipalID

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

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

	ts, _, err := startSolventServer(ctx, solventDB)
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

	execReg, err := executor.NewRegistry("recording", "", "ibmendoza/reference-loop-test", "ref-loop.yml", "main")
	if err != nil {
		t.Fatalf("create executor registry: %v", err)
	}

	evidenceCollector := agent.NewEvidenceCollector(runID, scenarioID)

	refAgent := agent.NewAgent(
		conductorClient,
		solventClient,
		execReg.GetRegistry(),
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
	expectedOpID := "deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:" + runID
	var ops []string
	for _, r := range records {
		if r.OperationID != "" {
			ops = append(ops, r.OperationID)
		}
	}

	if len(ops) == 0 {
		t.Fatal("no operation IDs found in evidence")
	}

	first := ops[0]
	for _, op := range ops[1:] {
		if op != first {
			t.Fatalf("operation ID mismatch: %s != %s", first, op)
		}
	}

	if first != expectedOpID {
		t.Fatalf("expected operation ID %s, got %s", expectedOpID, first)
	}
}
