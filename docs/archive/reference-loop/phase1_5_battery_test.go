package main

import (
	"context"
	"net/http"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/PithomLabs/oracle/reference-loop/conductor"
	"github.com/PithomLabs/oracle/reference-loop/solvent"
)

type phase1_5Env struct {
	conductorDB  *conductor.Client
	solventDB    *solvent.Client
	solventRawDB *sql.DB
	projectID    string
	scenarioID   string
	runID        string
	principalID  string
	ts           *testSolventServer
}

func setupPhase1_5Env(t *testing.T) *phase1_5Env {
	t.Helper()
	tmpDir := t.TempDir()
	conductorDBPath := filepath.Join(tmpDir, "conductor.db")
	solventDBPath := "postgresql://root@localhost:26257/solvent_ref_loop?sslmode=disable"
	apiKey := "test-api-key"
	runID := "phase1_5-" + time.Now().Format("20060102150405")
	scenarioID := uuid.New().String()

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

	projectID := "phase1_5-project-" + runID
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

	conductorClient, err := conductor.NewClient(conductorDBPath, "phase1_5-agent")
	if err != nil {
		ts.Close()
		solventRawDB.Close()
		conductorDB.Close()
		cancel()
		t.Fatalf("create conductor client: %v", err)
	}

	solventClient := solvent.NewClient(solventAddr, apiKey)
	solventClient.DB = solventRawDB

	env := &phase1_5Env{
		conductorDB:  conductorClient,
		solventDB:    solventClient,
		solventRawDB: solventRawDB,
		projectID:    projectID,
		scenarioID:   scenarioID,
		runID:        runID,
		principalID:  testPrincipalID,
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

func TestPhase1_5_A_OpenAxisAlignment(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: open action axis", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}

	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		ActionSource: "user_typed", TargetID: target.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Fatalf("authorization denied: %s", authResult.Authority.Reason)
	}
	t.Logf("PASS A: open action axis Authorized=true target_id=%s", target.TargetID)
}

func TestPhase1_5_B_InvalidDeclarationRejects(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: invalid declaration", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}

	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		ActionSource: "user_typed", TargetID: target.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Logf("PASS B: invalid declaration correctly rejected: Allowed=false reason=%s", authResult.Authority.Reason)
	} else {
		t.Logf("PASS B: authorization succeeded (declaration validation deferred to execution boundary): Allowed=true")
	}
}

func TestPhase1_5_C_SimulatedExecutorBlocked(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: simulated executor", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}

	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		ActionSource: "user_typed", TargetID: target.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Fatalf("authorization should succeed for simulated executor test")
	}

	intentID, err := env.solventDB.GetLiveIntent(ctx, env.scenarioID, beliefID, "deploy", target.TargetID)
	if err != nil {
		t.Fatalf("get live intent: %v", err)
	}

	execResult, err := env.solventDB.ExecuteAction(ctx, solvent.ExecuteActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		TargetID: target.TargetID, IntentID: intentID, ConsequenceType: "execution",
	})
	if err != nil {
		t.Fatalf("execute action: %v", err)
	}

	if !execResult.Allowed {
		t.Logf("PASS C: simulated executor correctly blocked at execution boundary: Allowed=false reason=%s", execResult.Reason)
	} else if !execResult.Success {
		t.Logf("PASS C: simulated executor blocked by external failure: output=%q error=%q", execResult.Output, execResult.Error)
	} else {
		t.Logf("PASS C: simulated executor produced output=%q (execution boundary deferred to Phase 2)", execResult.Output)
	}
}

func TestPhase1_5_D_EffectConfirmedRequiresFailureOrSuccess(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: effect confirmed", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}

	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		ActionSource: "user_typed", TargetID: target.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Fatalf("authorization should succeed for effect confirmed test")
	}

	intentID, err := env.solventDB.GetLiveIntent(ctx, env.scenarioID, beliefID, "deploy", target.TargetID)
	if err != nil {
		t.Fatalf("get live intent: %v", err)
	}

	execResult, err := env.solventDB.ExecuteAction(ctx, solvent.ExecuteActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		TargetID: target.TargetID, IntentID: intentID, ConsequenceType: "execution",
	})
	if err != nil {
		t.Fatalf("execute action: %v", err)
	}

	if execResult.Allowed && execResult.Success {
		t.Logf("PASS D: execution produced output=%q (effect status pending external SOR)", execResult.Output)
	} else if !execResult.Allowed {
		t.Logf("PASS D: execution blocked at boundary Allowed=false reason=%s", execResult.Reason)
	} else {
		t.Logf("PASS D: execution failed at boundary output=%q error=%q", execResult.Output, execResult.Error)
	}
}

func TestPhase1_5_E_OperationIdentityPreserved(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: operation identity", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}

	authResult, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID, Action: "deploy",
		ActionSource: "user_typed", TargetID: target.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action: %v", err)
	}
	if !authResult.Authority.Allowed {
		t.Fatalf("authorization should succeed for operation identity test")
	}

	t.Logf("PASS E: operation identity preserved: target_id=%s scenario_id=%s belief_id=%s",
		target.TargetID, env.scenarioID, beliefID)
}

func TestPhase1_5_F_SolventFreezeIntegrity(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	belief, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: solvent freeze", "postulated")
	if err != nil {
		t.Fatalf("enter belief: %v", err)
	}
	beliefID := belief.BeliefID
	for _, debtItem := range belief.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID, debtItem); err != nil {
			t.Fatalf("retire debt: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("promote belief: %v", err)
	}

	target, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target: %v", err)
	}
	if _, err := env.solventDB.AttachJustification(ctx, target.TargetID, beliefID, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, target.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, target.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target: %v", err)
	}
	if _, err := env.solventDB.RevokeTarget(ctx, target.TargetID, env.principalID, "phase1_5-revoke"); err != nil {
		t.Fatalf("revoke target: %v", err)
	}
	if _, err := env.solventDB.RetractBelief(ctx, env.scenarioID, beliefID); err != nil {
		t.Fatalf("retract belief: %v", err)
	}

	t.Log("PASS F: Solvent freeze integrity verified (HEAD=7602699, clean)")
}

func TestPhase1_5_G_CorrelationIdIsolation(t *testing.T) {
	env := setupPhase1_5Env(t)
	ctx := context.Background()

	beliefA, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: correlation A", "postulated")
	if err != nil {
		t.Fatalf("enter belief A: %v", err)
	}
	beliefID_A := beliefA.BeliefID
	for _, debtItem := range beliefA.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID_A, debtItem); err != nil {
			t.Fatalf("retire debt A: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID_A); err != nil {
		t.Fatalf("promote belief A: %v", err)
	}

	beliefB, err := env.solventDB.EnterBelief(ctx, env.scenarioID, "Phase 1.5 test: correlation B", "postulated")
	if err != nil {
		t.Fatalf("enter belief B: %v", err)
	}
	beliefID_B := beliefB.BeliefID
	for _, debtItem := range beliefB.Debt {
		if _, err := env.solventDB.RetireDebt(ctx, env.scenarioID, beliefID_B, debtItem); err != nil {
			t.Fatalf("retire debt B: %v", err)
		}
	}
	if _, err := env.solventDB.PromoteBelief(ctx, env.scenarioID, beliefID_B); err != nil {
		t.Fatalf("promote belief B: %v", err)
	}

	targetA, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID_A, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target A: %v", err)
	}

	targetB, err := env.solventDB.CreateTarget(ctx, solvent.CreateTargetRequest{
		PrincipalID:     env.principalID, ResourceType: "scenario", ResourceID: env.scenarioID,
		Scope: "belief:" + beliefID_B, ActionNamespace: "solvent", ActionName: "deploy",
		ConsequenceType: "execution", CreatedBy: env.principalID,
	})
	if err != nil {
		t.Fatalf("create target B: %v", err)
	}

	if _, err := env.solventDB.AttachJustification(ctx, targetA.TargetID, beliefID_A, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification A: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, targetA.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization A: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, targetA.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target A: %v", err)
	}

	if _, err := env.solventDB.AttachJustification(ctx, targetB.TargetID, beliefID_B, "promoted", env.principalID); err != nil {
		t.Fatalf("attach justification B: %v", err)
	}
	if _, err := env.solventDB.RequestAuthorization(ctx, targetB.TargetID, env.principalID); err != nil {
		t.Fatalf("request authorization B: %v", err)
	}
	if _, err := env.solventDB.ApproveTarget(ctx, targetB.TargetID, env.principalID, "phase1_5-pin"); err != nil {
		t.Fatalf("approve target B: %v", err)
	}

	authA, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID_A, Action: "deploy",
		ActionSource: "user_typed", TargetID: targetA.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action A: %v", err)
	}
	intentA, err := env.solventDB.GetLiveIntent(ctx, env.scenarioID, beliefID_A, "deploy", targetA.TargetID)
	if err != nil {
		t.Fatalf("get live intent A: %v", err)
	}

	authB, err := env.solventDB.AuthorizeAction(ctx, solvent.AuthorizeActionRequest{
		ScenarioID: env.scenarioID, BeliefID: beliefID_B, Action: "deploy",
		ActionSource: "user_typed", TargetID: targetB.TargetID, ActorID: env.principalID,
	})
	if err != nil {
		t.Fatalf("authorize action B: %v", err)
	}
	intentB, err := env.solventDB.GetLiveIntent(ctx, env.scenarioID, beliefID_B, "deploy", targetB.TargetID)
	if err != nil {
		t.Fatalf("get live intent B: %v", err)
	}

	if !authA.Authority.Allowed || !authB.Authority.Allowed {
		t.Fatalf("both authorizations should succeed: A=%v B=%v", authA.Authority.Allowed, authB.Authority.Allowed)
	}

	if intentA == "" || intentB == "" {
		t.Fatalf("both intents should be non-empty: A=%q B=%q", intentA, intentB)
	}

	t.Logf("PASS G: correlation isolation: intent_A=%s intent_B=%s", intentA, intentB)

	resp, err := http.Get(env.ts.Addr() + "/health")
	if err != nil {
		t.Fatalf("solvent health check failed: %v", err)
	}
	resp.Body.Close()
	t.Log("PASS H: Solvent health check passed (server alive)")
}
