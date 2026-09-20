package coordinator

import (
	"testing"
)

func TestProjectionQueueNew(t *testing.T) {
	queue := NewProjectionQueue()
	if queue == nil {
		t.Fatal("queue is nil")
	}
}

func TestProjectionQueueEnqueue(t *testing.T) {
	queue := NewProjectionQueue()
	record := ProjectionRecord{
		BeliefID:    "belief-001",
		TaskData:    "task data",
		Error:       "test error",
		ProjectionID: "proj-001",
	}

	queue.Enqueue(record)

	if queue.PendingCount() != 1 {
		t.Errorf("expected 1 pending, got %d", queue.PendingCount())
	}
}

func TestProjectionQueueRetryPending(t *testing.T) {
	queue := NewProjectionQueue()
	record := ProjectionRecord{
		BeliefID:    "belief-001",
		TaskData:    "task data",
		Error:       "test error",
		ProjectionID: "proj-001",
	}

	queue.Enqueue(record)

	err := queue.RetryPending()
	if err != nil {
		t.Fatalf("RetryPending failed: %v", err)
	}

	// After retry, queue should be empty (for POC)
	if queue.PendingCount() != 0 {
		t.Errorf("expected 0 pending after retry, got %d", queue.PendingCount())
	}
}
