package coordinator

import (
	"testing"
	"time"
)

func TestIdempotencyCacheNew(t *testing.T) {
	cache := NewIdempotencyCache()
	if cache == nil {
		t.Fatal("cache is nil")
	}
}

func TestIdempotencyCacheBeginNew(t *testing.T) {
	cache := NewIdempotencyCache()
	result, state := cache.Begin("hash1")
	if result != nil {
		t.Error("result should be nil for new entry")
	}
	if state != StateNew {
		t.Errorf("expected StateNew, got %v", state)
	}
}

func TestIdempotencyCacheBeginInFlight(t *testing.T) {
	cache := NewIdempotencyCache()
	cache.Begin("hash1")
	result, state := cache.Begin("hash1")
	if result != nil {
		t.Error("result should be nil for in-flight entry")
	}
	if state != StateInFlight {
		t.Errorf("expected StateInFlight, got %v", state)
	}
}

func TestIdempotencyCacheComplete(t *testing.T) {
	cache := NewIdempotencyCache()
	cache.Begin("hash1")

	completion := &CompilationResult{
		PacketID:  "test-001",
		Result:    "compiled",
		CreatedAt: time.Now(),
	}
	cache.Complete("hash1", completion)

	result, state := cache.Begin("hash1")
	if result == nil {
		t.Error("result should not be nil for completed entry")
	}
	if state != StateCompleted {
		t.Errorf("expected StateCompleted, got %v", state)
	}
	if result.PacketID != "test-001" {
		t.Errorf("wrong packet_id: %s", result.PacketID)
	}
}

func TestIdempotencyCacheFail(t *testing.T) {
	cache := NewIdempotencyCache()
	cache.Begin("hash1")
	cache.Fail("hash1")

	result, state := cache.Begin("hash1")
	if result != nil {
		t.Error("result should be nil for failed entry")
	}
	if state != StateInFlight {
		t.Errorf("expected StateInFlight for retry, got %v", state)
	}
}

func TestIdempotencyCacheWait(t *testing.T) {
	cache := NewIdempotencyCache()
	cache.Begin("hash2")

	go func() {
		time.Sleep(10 * time.Millisecond)
		completion := &CompilationResult{PacketID: "test-002"}
		cache.Complete("hash2", completion)
	}()

	result := cache.Wait("hash2")
	if result == nil {
		t.Error("result should not be nil after wait")
	}
	if result.PacketID != "test-002" {
		t.Errorf("wrong packet_id: %s", result.PacketID)
	}
}

func TestComputeCanonicalHash(t *testing.T) {
	pkt := map[string]interface{}{
		"packet_id": "test-001",
		"role":      "work",
		"beliefs":   []interface{}{},
	}

	hash1, err := ComputeCanonicalHash(pkt)
	if err != nil {
		t.Fatalf("ComputeCanonicalHash failed: %v", err)
	}

	hash2, err := ComputeCanonicalHash(pkt)
	if err != nil {
		t.Fatalf("ComputeCanonicalHash failed: %v", err)
	}

	if hash1 != hash2 {
		t.Error("canonical hash not deterministic")
	}
}

func TestComputeCanonicalHashDifferentContent(t *testing.T) {
	pkt1 := map[string]interface{}{
		"packet_id": "test-001",
		"role":      "work",
	}
	pkt2 := map[string]interface{}{
		"packet_id": "test-002",
		"role":      "adversarial",
	}

	hash1, _ := ComputeCanonicalHash(pkt1)
	hash2, _ := ComputeCanonicalHash(pkt2)

	if hash1 == hash2 {
		t.Error("different content should produce different hashes")
	}
}

func TestComputeCanonicalHashExcludesPacketID(t *testing.T) {
	pkt1 := map[string]interface{}{
		"packet_id": "test-001",
		"role":      "work",
	}
	pkt2 := map[string]interface{}{
		"packet_id": "test-002",
		"role":      "work",
	}

	hash1, _ := ComputeCanonicalHash(pkt1)
	hash2, _ := ComputeCanonicalHash(pkt2)

	if hash1 != hash2 {
		t.Error("canonical hash should exclude packet_id")
	}
}
