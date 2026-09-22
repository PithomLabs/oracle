package work

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Store handles task persistence with CockroachDB.
type Store struct {
	db *sql.DB
}

// NewStore creates a new work Store.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Errors
var (
	ErrNotFound          = fmt.Errorf("record not found")
	ErrInvalidTransition = fmt.Errorf("invalid lifecycle transition")
	ErrClaimFailed       = fmt.Errorf("claim failed")
	ErrReleaseFailed     = fmt.Errorf("release failed: task not in active state or not assigned to caller")
	ErrBlockedByDependency = fmt.Errorf("task has unfinished blocker")
)

// TaskUpdateFields defines which fields can be updated via generic PATCH.
type TaskUpdateFields struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

func generateID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("generate uuid v7: %w", err)
	}
	return id.String(), nil
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Create inserts a new task in proposed state.
func (s *Store) Create(ctx context.Context, task *Task) error {
	if task.ID == "" {
		id, err := generateID()
		if err != nil {
			return err
		}
		task.ID = id
	}
	if task.ProjectID == "" {
		return fmt.Errorf("task project_id is required")
	}
	if task.Title == "" {
		return fmt.Errorf("task title is required")
	}

	task.Status = TaskStatusProposed
	task.CurrentAgent = nil
	if task.Priority == "" {
		task.Priority = TaskPriorityMedium
	}
	now := nowRFC3339()
	task.CreatedAt = now
	task.UpdatedAt = now

	_, err := s.db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, description, status, priority, current_agent, governance_ref, reopened_from_task_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		task.ID, task.ProjectID, task.Title, task.Description, task.Status, task.Priority,
		task.CurrentAgent, task.GovernanceRef, task.ReopenedFromTaskID, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

// GetByID retrieves a task by ID.
func (s *Store) GetByID(ctx context.Context, id string) (*Task, error) {
	task := &Task{}
	var desc sql.NullString
	err := s.db.QueryRowContext(ctx,
		`SELECT id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
		 FROM conductor_task WHERE id = $1`, id).Scan(
		&task.ID, &task.ProjectID, &task.Title, &desc, &task.Status,
		&task.Priority, &task.CurrentAgent, &task.GovernanceRef, &task.OriginPacketID,
		&task.ReopenedFromTaskID, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get task: %w", err)
	}
	if desc.Valid {
		task.Description = desc.String
	}
	return task, nil
}

// Update updates ordinary task fields (NOT status, current_agent, or governance_ref).
func (s *Store) Update(ctx context.Context, taskID string, fields TaskUpdateFields) error {
	if taskID == "" {
		return fmt.Errorf("task ID is required")
	}

	setClauses := []string{}
	args := []interface{}{}

	if fields.Title != nil {
		setClauses = append(setClauses, "title = $1")
		args = append(args, *fields.Title)
	}
	if fields.Description != nil {
		setClauses = append(setClauses, "description = $1")
		args = append(args, *fields.Description)
	}
	if fields.Priority != nil {
		setClauses = append(setClauses, "priority = $1")
		args = append(args, *fields.Priority)
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at = $1")
	args = append(args, nowRFC3339())
	args = append(args, taskID)

	// Rebuild with positional $N placeholders
	clauses := make([]string, len(setClauses))
	for i, c := range setClauses {
		clauses[i] = strings.Replace(c, "$1", fmt.Sprintf("$%d", i+1), 1)
	}
	argsForQuery := make([]interface{}, len(args))
	copy(argsForQuery, args)

	query := fmt.Sprintf("UPDATE conductor_task SET %s WHERE id = $%d", strings.Join(clauses, ", "), len(clauses)+1)

	result, err := s.db.ExecContext(ctx, query, argsForQuery...)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// validateTransition checks if a transition from one status to another is allowed.
func validateTransition(fromStatus, toStatus string) bool {
	allowed := map[string][]string{
		TaskStatusProposed: {TaskStatusActive, TaskStatusBlocked, TaskStatusCancelled},
		TaskStatusActive:   {TaskStatusReview, TaskStatusBlocked, TaskStatusCancelled, TaskStatusProposed},
		TaskStatusReview:   {TaskStatusAccepted, TaskStatusActive, TaskStatusCancelled},
		TaskStatusBlocked:  {TaskStatusActive, TaskStatusCancelled},
	}
	transitions, ok := allowed[fromStatus]
	if !ok {
		return false
	}
	for _, valid := range transitions {
		if valid == toStatus {
			return true
		}
	}
	return false
}

// Transition performs a lifecycle transition atomically with activity event.
func (s *Store) Transition(ctx context.Context, taskID string, fromStatus string, toStatus string, actorType string, actorID string, action string, details string) error {
	if !validateTransition(fromStatus, toStatus) {
		return ErrInvalidTransition
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var currentStatus string
	err = tx.QueryRowContext(ctx, "SELECT status FROM conductor_task WHERE id = $1", taskID).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("get task status: %w", err)
	}
	if currentStatus != fromStatus {
		return ErrInvalidTransition
	}

	_, err = tx.ExecContext(ctx,
		"UPDATE conductor_task SET status = $1, updated_at = now() WHERE id = $2",
		toStatus, taskID)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}

	activityID, err := generateID()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO conductor_activity (id, task_id, actor_type, actor_id, action, details, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, now())`,
		activityID, taskID, actorType, actorID, action, details)
	if err != nil {
		return fmt.Errorf("insert activity: %w", err)
	}

	return tx.Commit()
}

// Claim atomically claims a task (proposed → active) after checking dependency blocking.
func (s *Store) Claim(ctx context.Context, taskID string, agentID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// Check dependency blocking: task with unfinished blocker cannot become active.
	var hasBlocker bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM conductor_dependency d
			JOIN conductor_task t ON t.id = d.blocked_by_id
			WHERE d.task_id = $1 AND t.status NOT IN ('accepted', 'cancelled')
		)`, taskID).Scan(&hasBlocker)
	if err != nil {
		return fmt.Errorf("check blockers: %w", err)
	}
	if hasBlocker {
		return ErrBlockedByDependency
	}

	result, err := tx.ExecContext(ctx,
		`UPDATE conductor_task
		 SET current_agent = $1, status = 'active', updated_at = now()
		 WHERE id = $2 AND status = 'proposed' AND current_agent IS NULL`,
		agentID, taskID)
	if err != nil {
		return fmt.Errorf("claim task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrClaimFailed
	}

	activityID, err := generateID()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT INTO conductor_activity (id, task_id, actor_type, actor_id, action, details, created_at)
		 VALUES ($1, $2, 'agent', $3, 'task.claimed', '{}', now())`,
		activityID, taskID, agentID)
	if err != nil {
		return fmt.Errorf("insert activity: %w", err)
	}

	return tx.Commit()
}

// Release releases a task back to proposed state.
func (s *Store) Release(ctx context.Context, taskID string, agentID string) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE conductor_task
		 SET current_agent = NULL, status = 'proposed', updated_at = now()
		 WHERE id = $1 AND current_agent = $2 AND status = 'active'`,
		taskID, agentID)
	if err != nil {
		return fmt.Errorf("release task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrReleaseFailed
	}

	activityID, err := generateID()
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO conductor_activity (id, task_id, actor_type, actor_id, action, details, created_at)
		 VALUES ($1, $2, 'agent', $3, 'task.released', '{}', now())`,
		activityID, taskID, agentID)
	if err != nil {
		return fmt.Errorf("insert activity: %w", err)
	}

	return nil
}

// ListAll returns all tasks across all projects.
func (s *Store) ListAll(ctx context.Context) ([]*Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
		 FROM conductor_task ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		var desc sql.NullString
		if err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &desc, &task.Status,
			&task.Priority, &task.CurrentAgent, &task.GovernanceRef, &task.OriginPacketID,
			&task.ReopenedFromTaskID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		if desc.Valid {
			task.Description = desc.String
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// ListByProject returns all tasks for a project.
func (s *Store) ListByProject(ctx context.Context, projectID string) ([]*Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, project_id, title, description, status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
		 FROM conductor_task WHERE project_id = $1 ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		var desc sql.NullString
		if err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &desc, &task.Status,
			&task.Priority, &task.CurrentAgent, &task.GovernanceRef, &task.OriginPacketID,
			&task.ReopenedFromTaskID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		if desc.Valid {
			task.Description = desc.String
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// GetTasksByGovernanceRef finds tasks linked to a belief via governance_ref.
func (s *Store) GetTasksByGovernanceRef(ctx context.Context, beliefID string) ([]*Task, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, project_id, title, COALESCE(description, ''), status, priority, current_agent, governance_ref, origin_packet_id, reopened_from_task_id, created_at, updated_at
		 FROM conductor_task WHERE governance_ref = $1::UUID ORDER BY created_at DESC`, beliefID)
	if err != nil {
		return nil, fmt.Errorf("list tasks by governance_ref: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		var desc sql.NullString
		if err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Title, &desc, &task.Status,
			&task.Priority, &task.CurrentAgent, &task.GovernanceRef, &task.OriginPacketID,
			&task.ReopenedFromTaskID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		if desc.Valid {
			task.Description = desc.String
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

// AddDependency adds a blocking dependency between tasks.
func (s *Store) AddDependency(ctx context.Context, taskID, blockedByID string) error {
	id, err := generateID()
	if err != nil {
		return err
	}
	dep := &Dependency{
		ID:          id,
		TaskID:      taskID,
		BlockedByID: blockedByID,
		CreatedAt:   nowRFC3339(),
	}
	_, err = s.db.ExecContext(ctx,
		`INSERT INTO conductor_dependency (id, task_id, blocked_by_id, created_at)
		 VALUES ($1, $2, $3, $4)`,
		dep.ID, dep.TaskID, dep.BlockedByID, dep.CreatedAt)
	if err != nil {
		return fmt.Errorf("add dependency: %w", err)
	}
	return nil
}
