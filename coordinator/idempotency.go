package coordinator

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// ReservationState represents the state of an idempotency reservation.
type ReservationState int

const (
	StateNew       ReservationState = iota
	StateInFlight
	StateCompleted
	StateFailed
)

// CompilationResult is the result of a packet compilation.
type CompilationResult struct {
	PacketID    string            `json:"packet_id"`
	BeliefIDs   map[string]string `json:"belief_ids"`
	EvidenceIDs []string          `json:"evidence_ids"`
	EdgeCount   int               `json:"edge_count"`
	TaskCount   int               `json:"task_count"`
	Result      string            `json:"result"`
	CreatedAt   time.Time         `json:"created_at"`
}

type reservation struct {
	state  ReservationState
	result *CompilationResult
	owner  chan struct{}
}

// IdempotencyCache is a process-local canonical-content idempotency cache.
type IdempotencyCache struct {
	mu      sync.RWMutex
	entries map[string]*reservation
}

// NewIdempotencyCache creates a new empty cache.
func NewIdempotencyCache() *IdempotencyCache {
	return &IdempotencyCache{
		entries: make(map[string]*reservation),
	}
}

// Begin is an atomic check-and-reserve operation.
func (c *IdempotencyCache) Begin(canonicalHash string) (*CompilationResult, ReservationState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if r, exists := c.entries[canonicalHash]; exists {
		switch r.state {
		case StateCompleted:
			return r.result, StateCompleted
		case StateInFlight:
			return nil, StateInFlight
		case StateFailed:
			r.state = StateInFlight
			r.owner = make(chan struct{})
			return nil, StateInFlight
		}
	}

	r := &reservation{
		state: StateInFlight,
		owner: make(chan struct{}),
	}
	c.entries[canonicalHash] = r
	return nil, StateNew
}

// Complete marks a reservation as completed with the result.
func (c *IdempotencyCache) Complete(canonicalHash string, result *CompilationResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if r, exists := c.entries[canonicalHash]; exists {
		r.state = StateCompleted
		r.result = result
		close(r.owner)
	}
}

// Fail marks a reservation as failed.
func (c *IdempotencyCache) Fail(canonicalHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if r, exists := c.entries[canonicalHash]; exists {
		r.state = StateFailed
		close(r.owner)
	}
}

// Wait blocks until the reservation completes or fails.
func (c *IdempotencyCache) Wait(canonicalHash string) *CompilationResult {
	c.mu.RLock()
	r, exists := c.entries[canonicalHash]
	c.mu.RUnlock()

	if !exists {
		return nil
	}

	<-r.owner
	return r.result
}

// ComputeCanonicalHash computes a deterministic SHA-256 hash of packet content.
func ComputeCanonicalHash(pkt interface{}) (string, error) {
	data, err := json.Marshal(pkt)
	if err != nil {
		return "", fmt.Errorf("failed to marshal packet: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return "", fmt.Errorf("failed to parse packet: %w", err)
	}

	delete(raw, "packet_id")
	delete(raw, "run_id")

	canonical, err := json.Marshal(raw)
	if err != nil {
		return "", fmt.Errorf("failed to marshal canonical form: %w", err)
	}

	hash := sha256.Sum256(canonical)
	return fmt.Sprintf("%x", hash), nil
}
