package work

// Task represents a unit of work that agents can claim and execute.
type Task struct {
	ID                string  `json:"id"`
	ProjectID         string  `json:"project_id"`
	Title             string  `json:"title"`
	Description       string  `json:"description,omitempty"`
	Status            string  `json:"status"`
	Priority          string  `json:"priority"`
	CurrentAgent      *string `json:"current_agent,omitempty"`
	GovernanceRef     *string `json:"governance_ref,omitempty"`
	ReopenedFromTaskID *string `json:"reopened_from_task_id,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// Task status constants
const (
	TaskStatusProposed  = "proposed"
	TaskStatusActive    = "active"
	TaskStatusReview    = "review"
	TaskStatusAccepted  = "accepted"
	TaskStatusBlocked   = "blocked"
	TaskStatusCancelled = "cancelled"
)

// Task priority constants
const (
	TaskPriorityLow      = "low"
	TaskPriorityMedium   = "medium"
	TaskPriorityHigh     = "high"
	TaskPriorityCritical = "critical"
)

// IsTerminal returns true if the task status is terminal.
func (t *Task) IsTerminal() bool {
	return t.Status == TaskStatusAccepted || t.Status == TaskStatusCancelled
}

// CanTransitionTo checks if a transition from current status to target is allowed.
func (t *Task) CanTransitionTo(target string) bool {
	allowed := map[string][]string{
		TaskStatusProposed: {TaskStatusActive, TaskStatusBlocked, TaskStatusCancelled},
		TaskStatusActive:   {TaskStatusReview, TaskStatusBlocked, TaskStatusCancelled, TaskStatusProposed},
		TaskStatusReview:   {TaskStatusAccepted, TaskStatusActive, TaskStatusCancelled},
		TaskStatusBlocked:  {TaskStatusActive, TaskStatusCancelled},
	}
	transitions, ok := allowed[t.Status]
	if !ok {
		return false
	}
	for _, valid := range transitions {
		if valid == target {
			return true
		}
	}
	return false
}

// Project represents a project that groups related work items.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Project status constants
const (
	ProjectStatusActive    = "active"
	ProjectStatusCompleted = "completed"
	ProjectStatusArchived  = "archived"
)

// Dependency represents a blocking relationship between tasks.
type Dependency struct {
	ID          string `json:"id"`
	TaskID      string `json:"task_id"`
	BlockedByID string `json:"blocked_by_id"`
	CreatedAt   string `json:"created_at"`
}

// Activity represents an immutable event log entry for task mutations and agent actions.
type Activity struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`
	Action    string `json:"action"`
	Details   string `json:"details,omitempty"`
	CreatedAt string `json:"created_at"`
}

// Actor type constants
const (
	ActorTypeHuman  = "human"
	ActorTypeAgent  = "agent"
	ActorTypeSystem = "system"
)

// Action constants
const (
	ActionTaskClaimed   = "task.claimed"
	ActionTaskSubmitted = "task.submitted"
	ActionTaskAccepted  = "task.accepted"
	ActionTaskRejected  = "task.rejected"
	ActionTaskCancelled = "task.cancelled"
	ActionTaskBlocked   = "task.blocked"
	ActionTaskUnblocked = "task.unblocked"
	ActionTaskUpdated   = "task.updated"
	ActionTaskReleased  = "task.released"
)
