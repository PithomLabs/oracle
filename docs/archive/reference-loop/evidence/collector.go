package evidence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	refExecutor "github.com/PithomLabs/oracle/reference-loop/executor"
)

type Collector struct {
	runID       string
	scenarioID  string
	evidenceDir string
	conductorDB *sql.DB
	solventDB   *sql.DB
	executorReg *refExecutor.Registry
	githubToken string
	githubRepo  string
}

func NewCollector(runID, scenarioID, evidenceDir, conductorDBPath, solventDBPath string, executorReg *refExecutor.Registry) *Collector {
	c := &Collector{
		runID:       runID,
		scenarioID:  scenarioID,
		evidenceDir: evidenceDir,
		executorReg: executorReg,
	}

	if conductorDBPath != "" {
		db, err := sql.Open("sqlite", conductorDBPath)
		if err == nil {
			c.conductorDB = db
		}
	}

	if solventDBPath != "" {
		db, err := sql.Open("pgx", solventDBPath)
		if err == nil {
			c.solventDB = db
		}
	}

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		c.githubToken = token
		c.githubRepo = os.Getenv("GITHUB_REPO")
	}

	return c
}

func (c *Collector) Collect(ctx context.Context) (*EvidencePackage, error) {
	pkg := &EvidencePackage{
		RunID:      c.runID,
		ScenarioID: c.scenarioID,
		StartTime:  time.Now(),
		Records:    []EvidenceRecord{},
	}

	if c.conductorDB != nil {
		if err := c.collectConductorEvidence(ctx, pkg); err != nil {
			return nil, fmt.Errorf("conductor evidence: %w", err)
		}
	}

	if c.solventDB != nil {
		if err := c.collectSolventEvidence(ctx, pkg); err != nil {
			return nil, fmt.Errorf("solvent evidence: %w", err)
		}
	}

	if c.executorReg != nil {
		c.collectExecutorEvidence(pkg)
	}

	if c.githubToken != "" && c.githubRepo != "" {
		if err := c.collectGitHubEvidence(ctx, pkg); err != nil {
			fmt.Printf("Warning: GitHub evidence collection failed: %v\n", err)
		}
	}

	pkg.EndTime = time.Now()
	return pkg, nil
}

func (c *Collector) collectConductorEvidence(ctx context.Context, pkg *EvidencePackage) error {
	rows, err := c.conductorDB.QueryContext(ctx, `
		SELECT id, task_id, actor_type, actor_id, action, details, created_at
		FROM conductor_activity
		ORDER BY created_at
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, taskID, actorType, actorID, action, details, createdAt string
		if err := rows.Scan(&id, &taskID, &actorType, &actorID, &action, &details, &createdAt); err != nil {
			continue
		}

		ts, _ := time.Parse(time.RFC3339, createdAt)
		raw, _ := json.Marshal(map[string]interface{}{
			"id": id, "task_id": taskID, "actor_type": actorType,
			"actor_id": actorID, "action": action, "details": details, "created_at": createdAt,
		})

		event := c.mapConductorAction(action)
		pkg.Records = append(pkg.Records, EvidenceRecord{
			Event:      event,
			Owner:      "conductor",
			Source:     "conductor.activity",
			Timestamp:  ts,
			RunID:      c.runID,
			ScenarioID: c.scenarioID,
			TaskID:     taskID,
			Raw:        raw,
		})
	}

	taskRows, err := c.conductorDB.QueryContext(ctx, `
		SELECT id, project_id, title, status, current_agent, governance_ref, created_at, updated_at
		FROM conductor_task
	`)
	if err != nil {
		return err
	}
	defer taskRows.Close()

	for taskRows.Next() {
		var id, projectID, title, status, currentAgent, governanceRef, createdAt, updatedAt string
		if err := taskRows.Scan(&id, &projectID, &title, &status, &currentAgent, &governanceRef, &createdAt, &updatedAt); err != nil {
			continue
		}

		ts, _ := time.Parse(time.RFC3339, updatedAt)
		raw, _ := json.Marshal(map[string]interface{}{
			"id": id, "project_id": projectID, "title": title, "status": status,
			"current_agent": currentAgent, "governance_ref": governanceRef,
			"created_at": createdAt, "updated_at": updatedAt,
		})

		event := "TASK_STATUS_" + status
		pkg.Records = append(pkg.Records, EvidenceRecord{
			Event:      event,
			Owner:      "conductor",
			Source:     "conductor.task",
			Timestamp:  ts,
			RunID:      c.runID,
			ScenarioID: c.scenarioID,
			TaskID:     id,
			Raw:        raw,
		})
	}

	return nil
}

func (c *Collector) mapConductorAction(action string) string {
	switch action {
	case "task.created":
		return "TASK_CREATED"
	case "task.claimed":
		return "TASK_CLAIMED"
	case "task.submitted":
		return "TASK_SUBMITTED"
	case "task.accepted":
		return "TASK_ACCEPTED"
	case "task.rejected":
		return "TASK_REJECTED"
	case "task.blocked":
		return "TASK_BLOCKED"
	case "task.unblocked":
		return "TASK_UNBLOCKED"
	case "work.completed":
		return "WORK_COMPLETED"
	default:
		return action
	}
}

func (c *Collector) collectSolventEvidence(ctx context.Context, pkg *EvidencePackage) error {
	rows, err := c.solventDB.QueryContext(ctx, `
		SELECT id, scenario_id, type, actor_id, subject_id, details, sqlstate, constraint_name, refusal, created_at
		FROM audit_activity
		ORDER BY created_at
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, scenarioID, actType, actorID, subjectID, details, sqlstate, constraintName string
		var refusal bool
		var createdAt string
		if err := rows.Scan(&id, &scenarioID, &actType, &actorID, &subjectID, &details, &sqlstate, &constraintName, &refusal, &createdAt); err != nil {
			continue
		}

		ts, _ := time.Parse(time.RFC3339, createdAt)
		raw, _ := json.Marshal(map[string]interface{}{
			"id": id, "scenario_id": scenarioID, "type": actType, "actor_id": actorID,
			"subject_id": subjectID, "details": details, "sqlstate": sqlstate,
			"constraint_name": constraintName, "refusal": refusal, "created_at": createdAt,
		})

		event := c.mapSolventActivity(actType)
		pkg.Records = append(pkg.Records, EvidenceRecord{
			Event:      event,
			Owner:      "solvent",
			Source:     "solvent.audit",
			Timestamp:  ts,
			RunID:      c.runID,
			ScenarioID: c.scenarioID,
			TaskID:     subjectID,
			Raw:        raw,
		})
	}

	intentRows, err := c.solventDB.QueryContext(ctx, `
		SELECT id, scenario_id, belief_id, belief_status, action, state, target_id, snapshot_id
		FROM action_intent
		ORDER BY id
	`)
	if err != nil {
		return err
	}
	defer intentRows.Close()

	for intentRows.Next() {
		var id, scenarioID, beliefID, beliefStatus, action, state, targetID, snapshotID string
		if err := intentRows.Scan(&id, &scenarioID, &beliefID, &beliefStatus, &action, &state, &targetID, &snapshotID); err != nil {
			continue
		}

		raw, _ := json.Marshal(map[string]interface{}{
			"id": id, "scenario_id": scenarioID, "belief_id": beliefID,
			"belief_status": beliefStatus, "action": action, "state": state,
			"target_id": targetID, "snapshot_id": snapshotID,
		})

		event := "INTENT_" + state
		pkg.Records = append(pkg.Records, EvidenceRecord{
			Event:          event,
			Owner:          "solvent",
			Source:         "solvent.intent",
			Timestamp:      time.Now(),
			RunID:          c.runID,
			ScenarioID:     c.scenarioID,
			IntentID:       id,
			DeclarationVer: snapshotID,
			OperationID:    action + ":" + targetID,
			Raw:            raw,
		})
	}

	return nil
}

func (c *Collector) mapSolventActivity(actType string) string {
	switch actType {
	case "authorization_granted":
		return "AUTHORIZED"
	case "authorization_denied":
		return "AUTHORIZATION_DENIED"
	case "adapter_invoked":
		return "EXECUTION_ATTEMPTED"
	case "executor_completed":
		return "EXECUTION_COMPLETED"
	case "executor_failed":
		return "EXECUTION_FAILED"
	case "intent_completion_failed":
		return "INTENT_COMPLETION_FAILED"
	case "reconciliation_completed":
		return "RECONCILIATION_COMPLETED"
	case "reconciliation_failed":
		return "RECONCILIATION_FAILED"
	default:
		return actType
	}
}

func (c *Collector) collectExecutorEvidence(pkg *EvidencePackage) {
	if c.executorReg.WasRecordingCalled() {
		params := c.executorReg.GetRecordingParams()
		raw, _ := json.Marshal(params)
		repo := ""
		workflow := ""
		ref := ""
		if r, ok := params["repo"].(string); ok {
			repo = r
		}
		if w, ok := params["workflow"].(string); ok {
			workflow = w
		}
		if f, ok := params["ref"].(string); ok {
			ref = f
		}
		opID := "deploy:" + repo + ":" + workflow + ":" + ref
		pkg.Records = append(pkg.Records, EvidenceRecord{
			Event:                "EXECUTION_ATTEMPTED",
			Owner:                "executor",
			Source:               "executor.recording",
			Timestamp:            time.Now(),
			RunID:                c.runID,
			ScenarioID:           c.scenarioID,
			OperationID:          opID,
			Raw:                  raw,
		})
	}
}

func (c *Collector) collectGitHubEvidence(ctx context.Context, pkg *EvidencePackage) error {
	if c.githubToken == "" || c.githubRepo == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/actions/workflows/ref-loop.yml/runs?event=workflow_dispatch&per_page=1", c.githubRepo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("build github request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.githubToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("github request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TotalCount int `json:"total_count"`
		WorkflowRuns []struct {
			ID        int    `json:"id"`
			Status    string `json:"status"`
			Conclusion string `json:"conclusion"`
			CreatedAt string `json:"created_at"`
		} `json:"workflow_runs"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("unmarshal github response: %w", err)
	}

	confirmed := result.TotalCount > 0 && result.WorkflowRuns[0].Conclusion == "success"
	raw, _ := json.Marshal(map[string]interface{}{
		"github_response": string(body),
		"confirmed":       confirmed,
		"run_count":       result.TotalCount,
	})

	record := EvidenceRecord{
		Event:       "EFFECT_CONFIRMED",
		Owner:       "external",
		Source:      "github.api",
		Timestamp:   time.Now(),
		RunID:       c.runID,
		ScenarioID:  c.scenarioID,
		OperationID: "deploy:" + c.githubRepo + ":ref-loop.yml:main",
		Raw:         raw,
	}
	if !confirmed {
		record.Event = "EFFECT_UNCONFIRMED"
	}
	pkg.Records = append(pkg.Records, record)
	return nil
}

func (c *Collector) Save(pkg *EvidencePackage) error {
	if err := os.MkdirAll(c.evidenceDir, 0755); err != nil {
		return err
	}
	file := filepath.Join(c.evidenceDir, "evidence.json")
	data, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func (c *Collector) Close() {
	if c.conductorDB != nil {
		c.conductorDB.Close()
	}
	if c.solventDB != nil {
		c.solventDB.Close()
	}
}
