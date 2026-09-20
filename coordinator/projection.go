package coordinator

import (
	"fmt"
	"sync"
)

// ProjectionRecord represents a failed Conductor projection.
type ProjectionRecord struct {
	BeliefID    string `json:"belief_id"`
	TaskData    string `json:"task_data"`
	Error       string `json:"error"`
	RetryCount  int    `json:"retry_count"`
	ProjectionID string `json:"projection_id"`
}

// ProjectionQueue is an in-memory queue for failed Conductor projections.
type ProjectionQueue struct {
	mu      sync.RWMutex
	records []ProjectionRecord
}

// NewProjectionQueue creates a new empty queue.
func NewProjectionQueue() *ProjectionQueue {
	return &ProjectionQueue{
		records: make([]ProjectionRecord, 0),
	}
}

// Enqueue adds a failed projection to the queue.
func (q *ProjectionQueue) Enqueue(record ProjectionRecord) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.records = append(q.records, record)
}

// RetryPending retries all pending projections.
func (q *ProjectionQueue) RetryPending() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	failed := make([]ProjectionRecord, 0)
	for _, record := range q.records {
		if err := retryProjection(record); err != nil {
			record.RetryCount++
			record.Error = err.Error()
			failed = append(failed, record)
		}
	}
	q.records = failed
	return nil
}

// retryProjection attempts to retry a single projection.
func retryProjection(record ProjectionRecord) error {
	// For POC, we just log the retry attempt
	fmt.Printf("retrying projection %s (attempt %d)\n", record.ProjectionID, record.RetryCount+1)
	return nil
}

// PendingCount returns the number of pending projections.
func (q *ProjectionQueue) PendingCount() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.records)
}
