package coordinator

import (
	"fmt"
	"strings"

	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// parsePackRef splits "bmist@1.1.0" into packID and version.
func parsePackRef(ref string) (string, string) {
	parts := strings.SplitN(ref, "@", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, ""
}

// compile is the core packet compilation logic.
func (c *Coordinator) compile(pkt *packetv1.Packet) (*CompilationResult, error) {
	// 1. Validate packet
	if err := ValidatePacket(pkt, c.packRegistry); err != nil {
		return nil, fmt.Errorf("packet validation failed: %w", err)
	}

	// 2. Resolve pack_ref
	packID, packVersion := parsePackRef(pkt.PackRef)
	_, err := c.packRegistry.Get(packID, packVersion)
	if err != nil {
		return nil, fmt.Errorf("pack resolution failed: %w", err)
	}

	// 3. Compute canonical hash
	canonicalHash, err := ComputeCanonicalHash(pkt)
	if err != nil {
		return nil, fmt.Errorf("canonical hash computation failed: %w", err)
	}

	// 4. Check idempotency cache
	existing, state := c.idempotencyCache.Begin(canonicalHash)
	switch state {
	case StateCompleted:
		return existing, nil
	case StateInFlight:
		result := c.idempotencyCache.Wait(canonicalHash)
		if result != nil {
			return result, nil
		}
		return nil, fmt.Errorf("compilation failed for concurrent request")
	}

	// 5. Compile beliefs
	beliefIDs := make(map[string]string)
	for i, b := range pkt.Beliefs {
		beliefID := fmt.Sprintf("belief-%s-%d", pkt.PacketID, i)
		beliefIDs[b.LocalID] = beliefID
	}

	// 6. Compile evidence
	evidenceIDs := make([]string, 0)
	for _, e := range pkt.Evidence {
		evidenceID := fmt.Sprintf("evidence-%s", e.LocalID)
		evidenceIDs = append(evidenceIDs, evidenceID)
	}

	// 7. Compile edges
	edgeCount := len(pkt.Edges)

	// 8. Compile tasks
	taskCount := len(pkt.Tasks)

	// 9. Build result
	result := &CompilationResult{
		PacketID:    pkt.PacketID,
		BeliefIDs:   beliefIDs,
		EvidenceIDs: evidenceIDs,
		EdgeCount:   edgeCount,
		TaskCount:   taskCount,
		Result:      "compiled",
	}

	// 10. Store in cache
	c.idempotencyCache.Complete(canonicalHash, result)

	return result, nil
}
