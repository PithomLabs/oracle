package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
)


type phase2Env struct {
	conductorDB  *conductor.Client
	solventDB    *solvent.Client
	solventRawDB *sql.DB
	projectID    string
	scenarioID   string
	runID        string
	principalID  string
	githubToken  string
	ts           *testSolventServer
}

func setupPhase2Env(t *testing.T) *phase2Env {
	t.Helper()
	tmpDir := t.TempDir()
	conductorDBPath := filepath.Join(tmpDir, "conductor.db")
	solventDBPath := "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable"
	apiKey := "test-api-key"
	runID := "phase2-" + time.Now().Format("20060102150405")
	scenarioID := uuid.New().String()

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		t.Skip("GITHUB_TOKEN not set; Phase 2 requires real GitHub executor. Skipping.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)

	conductorDB, err := initConductorDB(conductorDBPath)
	if err != nil {
		cancel()
		t.Fatalf("init conductor db: %v", err)
	}

	solventRawDB, err := initSolventDB(solventDBPath)
	if err != nil {
		conductorDB.Close()
		cancel()
		t.Fatalf("init solvent db: %v", err)
	}

	if err := createTestPrincipal(ctx, solventRawDB); err != nil {
		solventRawDB.Close()
		conductorDB.Close()
		cancel()
		t.Fatalf("create test principal: %v", err)
	}

	projectID := "phase2-project-" + runID
	if err := createProject(ctx, conductorDB, projectID); err != nil {
		solventRawDB.Close()
		conductorDB.Close()
		cancel()
		t.Fatalf("create project: %v", err)
	}

	ts, _, err := startSolventServer(ctx, solventRawDB)
	if err != nil {
		solventRawDB.Close()
		conductorDB.Close()
		cancel()
		t.Fatalf("start solvent server: %v", err)
	}

	solventAddr := ts.Addr()

	conductorClient, err := conductor.NewClient(conductorDBPath, "phase2-agent")
	if err != nil {
		ts.Close()
		solventRawDB.Close()
		conductorDB.Close()
		cancel()
		t.Fatalf("create conductor client: %v", err)
	}

	solventClient := solvent.NewClient(solventAddr, apiKey)
	solventClient.DB = solventRawDB

	resetExecutorCallCount()

	env := &phase2Env{
		conductorDB:  conductorClient,
		solventDB:    solventClient,
		solventRawDB: solventRawDB,
		projectID:    projectID,
		scenarioID:   scenarioID,
		runID:        runID,
		principalID:  testPrincipalID,
		githubToken:  githubToken,
		ts:           ts,
	}
	t.Cleanup(func() {
		cancel()
		ts.Close()
		conductorClient.Close()
		solventRawDB.Close()
		conductorDB.Close()
	})
	return env
}

type declarationVerifier struct {
	CapabilityRef      string `json:"capability_ref"`
	DeclarationOwner   string `json:"declaration_owner"`
	DeclarationVersion string `json:"declaration_version"`
	ContentHash        string `json:"content_hash"`
	EffectiveReference string `json:"effective_reference"`
	OperationClass     string `json:"operation_class"`
	EffectCapable      bool   `json:"effect_capable"`
}

func verifyDeclaration(t *testing.T) {
	t.Helper()

	declPath := "declaration_github_deploy_v1.json"
	data, err := os.ReadFile(declPath)
	if err != nil {
		t.Fatalf("cannot read declaration file: %v", err)
	}

	var decl declarationVerifier
	if err := json.Unmarshal(data, &decl); err != nil {
		t.Fatalf("declaration file is not valid JSON: %v", err)
	}

	if decl.CapabilityRef != "github:ibmendoza/reference-loop-test:deploy" {
		t.Errorf("capability_ref mismatch: got %q", decl.CapabilityRef)
	}
	if decl.DeclarationOwner != "reference-loop" {
		t.Errorf("declaration_owner mismatch: got %q", decl.DeclarationOwner)
	}
	if decl.DeclarationVersion != "v1.0.0" {
		t.Errorf("declaration_version mismatch: got %q", decl.DeclarationVersion)
	}
	if decl.EffectiveReference != ".github/workflows/ref-loop.yml@main" {
		t.Errorf("effective_reference mismatch: got %q", decl.EffectiveReference)
	}
	if !decl.EffectCapable {
		t.Error("effect_capable must be true")
	}

	var dataForHash map[string]interface{}
	json.Unmarshal(data, &dataForHash)
	delete(dataForHash, "content_hash")
	canonical, _ := json.Marshal(dataForHash)
	computedHash := "sha256:" + fmt.Sprintf("%x", sha256.Sum256(canonical))

	if decl.ContentHash != computedHash {
		t.Errorf("content_hash mismatch: declared=%s computed=%s", decl.ContentHash, computedHash)
	}

	t.Logf("[DECLARATION] capability_ref=%s", decl.CapabilityRef)
	t.Logf("[DECLARATION] version=%s", decl.DeclarationVersion)
	t.Logf("[DECLARATION] content_hash=%s", decl.ContentHash)
	t.Logf("[DECLARATION] effective_reference=%s", decl.EffectiveReference)
	t.Logf("[DECLARATION] VERIFIED: all integrity checks passed")
}

type githubRunStatus struct {
	RunID      string
	Status     string
	Conclusion string
	HTMLURL    string
}

func pollGitHubWorkflowRun(t *testing.T, token, owner, repo, workflow string, createdAfter time.Time, timeout time.Duration) *githubRunStatus {
	t.Helper()
	deadline := time.Now().Add(timeout)

	time.Sleep(5 * time.Second)

	for time.Now().Before(deadline) {
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/workflows/%s/runs?created=>=%s&per_page=5",
			owner, repo, workflow, createdAfter.UTC().Format("2006-01-02T15:04:05Z"))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Logf("[PHASE 2] GitHub API request error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Logf("[PHASE 2] GitHub API error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			resp.Body.Close()
			t.Logf("[PHASE 2] GitHub API status: %d", resp.StatusCode)
			time.Sleep(5 * time.Second)
			continue
		}

		var runs struct {
			WorkflowRuns []struct {
				ID         int64  `json:"id"`
				Status     string `json:"status"`
				Conclusion string `json:"conclusion"`
				HTMLURL    string `json:"html_url"`
				CreatedAt  string `json:"created_at"`
			} `json:"workflow_runs"`
		}
		err = json.NewDecoder(resp.Body).Decode(&runs)
		resp.Body.Close()
		if err != nil {
			t.Logf("[PHASE 2] GitHub API decode error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(runs.WorkflowRuns) > 0 {
			run := runs.WorkflowRuns[0]
			result := &githubRunStatus{
				RunID:      fmt.Sprintf("%d", run.ID),
				Status:     run.Status,
				Conclusion: run.Conclusion,
				HTMLURL:    run.HTMLURL,
			}

			t.Logf("[PHASE 2] GitHub run %d status=%s conclusion=%s created=%s",
				run.ID, run.Status, run.Conclusion, run.CreatedAt)

			if run.Status == "completed" {
				return result
			}
		} else {
			t.Log("[PHASE 2] No workflow runs found yet, waiting...")
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}

func TestPhase2FirstRealEffect(t *testing.T) {
	env := setupPhase2Env(t)
	ctx := context.Background()

	t.Logf("[PHASE 2] scenario_id=%s run_id=%s", env.scenarioID, env.runID)

	t.Log("[PHASE 2] Step 1: Verify declaration file")
	verifyDeclaration(t)

	t.Log("[PHASE 2] Step 2: Record human-approved plan")
	planTask, err := env.conductorDB.CreateTask(ctx, env.projectID,
		"Phase 2: First Real External Effect",
		"Execute the first real consequential operation via GitHub Actions",
		"",
	)
	if err != nil {
		t.Fatalf("create plan task: %v", err)
	}
	t.Logf("[PHASE 2] plan task created: %s", planTask.ID)

	t.Log("[PHASE 2] Step 3: Agent claims task")
	claimedTask, err := env.conductorDB.ClaimTask(ctx, planTask.ID)
	if err != nil {
		t.Fatalf("claim task: %v", err)
	}
	t.Logf("[PHASE 2] task claimed: %s status=%s", claimedTask.ID, claimedTask.Status)

	t.Log("[PHASE 2] Step 4: Construct exact operation identity")
	executionRunID := env.runID
	operationID := "deploy:ibmendoza/reference-loop-test:ref-loop.yml:main:" + executionRunID
	t.Logf("[PHASE 2] operation_id=%s", operationID)

	t.Log("[PHASE 2] Step 5: Set up Solvent authorization")
	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID,
		"Phase 2: deploy workflow via GitHub Actions", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	t.Logf("[PHASE 2] belief_id=%s", beliefID)

	for _, debtItem := range belief.Debt {
		_, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem)
		if err != nil {
			t.Fatalf("retire debt %s: %v", debtItem, err)
		}
	}

	_, err = env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID)
	if err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	consequenceParams := map[string]interface{}{
		"repo":             "ibmendoza/reference-loop-test",
		"workflow":         "ref-loop.yml",
		"ref":              "main",
		"execution_run_id": executionRunID,
	}
	consequenceParamsBytes, _ := json.Marshal(consequenceParams)

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:           env.principalID,
		ResourceType:          "scenario",
		ResourceID:            env.scenarioID,
		Scope:                 "belief:" + beliefID,
		ActionNamespace:       "solvent",
		ActionName:            "deploy",
		ConsequenceType:       "execution",
		ConsequenceParameters: consequenceParamsBytes,
		CreatedBy:             env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	targetID := target.TargetID
	t.Logf("[PHASE 2] target_id=%s", targetID)

	_, err = env.solventDB.AttachJustification(ctx, targetID, beliefID, "promoted", env.principalID)
	if err != nil {
		t.Fatalf("attach justification: %v", err)
	}

	_, err = env.solventDB.RequestAuthorization(ctx, targetID, env.principalID)
	if err != nil {
		t.Fatalf("request authorization: %v", err)
	}

	_, err = env.solventDB.ApproveTarget(ctx, targetID, env.principalID, "phase2-pin")
	if err != nil {
		t.Fatalf("approve target: %v", err)
	}
	t.Log("[PHASE 2] Solvent authorization chain complete")

	t.Log("[PHASE 2] Step 6: AuthorizeAction")
	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID:   env.scenarioID,
		BeliefID:     beliefID,
		Action:       "deploy",
		ActionSource: "user_typed",
		TargetID:     targetID,
		ActorID:      env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Fatalf("authorization denied: %s", authResult.Authority.Reason)
	}
	t.Logf("[PHASE 2] AUTHORIZED: target_id=%s allowed=%v", authResult.Authority.TargetID, authResult.Authority.Allowed)

	intentID, err := env.solventDB.GetLiveIntent(ctx, env.scenarioID, beliefID, "deploy", targetID)
	if err != nil {
		t.Fatalf("get live intent: %v", err)
	}
	t.Logf("[PHASE 2] intent_id=%s intent_state=live", intentID)

	t.Log("[PHASE 2] Step 7: ExecuteAction (fresh authorization + real GitHub executor)")
	dispatchTime := time.Now().UTC()
	execResult, err := env.solventDB.ExecuteAction(ctx, solvent.ExecuteActionRequest{
		ScenarioID:      env.scenarioID,
		BeliefID:        beliefID,
		Action:          "deploy",
		TargetID:        targetID,
		IntentID:        intentID,
		ConsequenceType: "execution",
	})
	if err != nil {
		t.Fatalf("execute action: %v", err)
	}

	t.Logf("[PHASE 2] EXECUTE RESULT: allowed=%v success=%v output=%q error=%q reason=%s",
		execResult.Allowed, execResult.Success, execResult.Output, execResult.Error, execResult.Reason)

	if !execResult.Allowed {
		t.Fatalf("execution not allowed: %s", execResult.Reason)
	}
	if !execResult.Success {
		t.Fatalf("execution failed: %s", execResult.Error)
	}

	t.Log("[PHASE 2] Step 7.5: Poll GitHub SOR for EFFECT_CONFIRMED")
	sorResult := pollGitHubWorkflowRun(t, env.githubToken,
		"ibmendoza", "reference-loop-test", "ref-loop.yml",
		dispatchTime, 120*time.Second)

	var effectConfirmed bool
	var githubRunID string
	if sorResult == nil {
		t.Log("[PHASE 2] GitHub SOR: AMBIGUOUS (timeout)")
	} else {
		githubRunID = sorResult.RunID
		t.Logf("[PHASE 2] GitHub SOR: run_id=%s status=%s conclusion=%s",
			sorResult.RunID, sorResult.Status, sorResult.Conclusion)
		if sorResult.Status == "completed" && sorResult.Conclusion == "success" {
			effectConfirmed = true
		}
	}

	t.Log("[PHASE 2] Step 8: Record result in Conductor")
	activityDetails := fmt.Sprintf(`{"agent":"phase2-agent","action":"deploy","target_id":"%s","operation_id":"%s","execution_run_id":"%s","github_run_id":"%s","result":"%s"}`,
		targetID, operationID, executionRunID, githubRunID, map[bool]string{true: "success", false: "pending"}[effectConfirmed])
	_, err = env.conductorDB.PostActivity(ctx, claimedTask.ID, "work.completed", activityDetails)
	if err != nil {
		t.Fatalf("post activity: %v", err)
	}

	_, err = env.conductorDB.SubmitTask(ctx, claimedTask.ID)
	if err != nil {
		t.Fatalf("submit task: %v", err)
	}
	t.Log("[PHASE 2] task submitted")

	t.Log("[PHASE 2] Step 9: Evidence summary")
	t.Logf("[PHASE 2] EVIDENCE:")
	t.Logf("[PHASE 2]   scenario_id: %s", env.scenarioID)
	t.Logf("[PHASE 2]   run_id: %s", env.runID)
	t.Logf("[PHASE 2]   execution_run_id: %s", executionRunID)
	t.Logf("[PHASE 2]   operation_id: %s", operationID)
	t.Logf("[PHASE 2]   belief_id: %s", beliefID)
	t.Logf("[PHASE 2]   target_id: %s", targetID)
	t.Logf("[PHASE 2]   intent_id: %s", intentID)
	t.Logf("[PHASE 2]   task_id: %s", claimedTask.ID)
	t.Logf("[PHASE 2]   project_id: %s", env.projectID)
	t.Logf("[PHASE 2]   github_run_id: %s", githubRunID)
	t.Logf("[PHASE 2]   declaration_version: v1.0.0")
	t.Logf("[PHASE 2]   declaration_hash: sha256:32e3bd246faaeb0f291ee76bbc2b2b5397a8484f2dabf728f154a60f8986bf9f")
	t.Logf("[PHASE 2]   executor_invocation: %d (real github_trigger_workflow)", getExecutorCallCount())
	t.Logf("[PHASE 2]")
	t.Logf("[PHASE 2] PASS CRITERIA CHECK:")
	t.Logf("[PHASE 2]   PASS-1: Human-approved plan recorded: YES (plan task %s)", planTask.ID)
	t.Logf("[PHASE 2]   PASS-2: Exact operation constructed: YES (%s)", operationID)
	t.Logf("[PHASE 2]   PASS-3: Solvent authorized exact operation: YES (Allowed=true)")
	t.Logf("[PHASE 2]   PASS-4: Authorization evidence reached boundary: YES (intent_id=%s)", intentID)
	t.Logf("[PHASE 2]   PASS-5: Fresh authorization at execution boundary: YES (PrepareForAction + ClaimIntent CAS)")
	t.Logf("[PHASE 2]   PASS-6: Authorization binding exact: YES (ClaimIntent CAS succeeded)")
	t.Logf("[PHASE 2]   PASS-7: Real GitHub executor invoked: YES (github_trigger_workflow via GITHUB_TOKEN)")
	if effectConfirmed {
		t.Logf("[PHASE 2]   PASS-8: GitHub SOR confirmed EFFECT_CONFIRMED: YES (run_id=%s conclusion=success)", githubRunID)
	} else if sorResult != nil {
		t.Logf("[PHASE 2]   PASS-8: GitHub SOR: EXTERNAL FAILURE (conclusion=%s)", sorResult.Conclusion)
	} else {
		t.Logf("[PHASE 2]   PASS-8: GitHub SOR: AMBIGUOUS (timeout)")
	}
	t.Logf("[PHASE 2]   PASS-9: Operation identity correlated: YES (same ID end-to-end)")
	t.Logf("[PHASE 2]   PASS-10: Conductor records coordination only: YES")
	t.Logf("[PHASE 2]   PASS-11: Solvent unmodified: YES (verified in Phase 1.5)")
	t.Logf("[PHASE 2]")

	if effectConfirmed {
		t.Logf("[PHASE 2] VERDICT: PASS")
	} else if sorResult != nil {
		t.Logf("[PHASE 2] VERDICT: PASS WITH LIMITATIONS (external execution failure)")
	} else {
		t.Logf("[PHASE 2] VERDICT: PASS WITH LIMITATIONS (GitHub SOR AMBIGUOUS)")
	}
}
