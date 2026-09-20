package coordinator

import (
	"fmt"
	"strings"

	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// SubmitPacket persists a compiled packet to Solvent and Conductor.
// Persistence order: beliefs → edges → evidence → tasks.
// No distributed transactions. Partial state is observable on failure.
func (c *Coordinator) SubmitPacket(pkt *packetv1.Packet) (*CompilationResult, error) {
	// 1. Compile the packet (validation, idempotency, canonical hash)
	result, err := c.compile(pkt)
	if err != nil {
		return nil, fmt.Errorf("compilation failed: %w", err)
	}

	// 2. Persist beliefs to Solvent
	beliefIDMap := make(map[string]string) // localID → solvent belief ID
	for _, b := range pkt.Beliefs {
		debt := b.Debt
		if debt == nil {
			debt = []string{}
		}
		solventID, err := c.solventClient.CreateBelief("", b.Claim, pkt.ScenarioID)
		if err != nil {
			return nil, fmt.Errorf("failed to persist belief %s: %w", b.LocalID, err)
		}
		beliefIDMap[b.LocalID] = solventID
	}

	// 3. Persist edges to Solvent
	for _, e := range pkt.Edges {
		fromID := resolveRef(e.FromRef, beliefIDMap)
		toID := resolveRef(e.ToRef, beliefIDMap)
		if fromID == "" || toID == "" {
			return nil, fmt.Errorf("edge references unresolved: from=%s to=%s", e.FromRef, e.ToRef)
		}
		if err := c.solventClient.CreateEdge(fromID, toID, e.Kind); err != nil {
			return nil, fmt.Errorf("failed to persist edge %s: %w", e.LocalID, err)
		}
	}

	// 4. Persist evidence to Solvent
	for _, ev := range pkt.Evidence {
		beliefID := resolveRef(ev.BeliefRef, beliefIDMap)
		if beliefID == "" {
			return nil, fmt.Errorf("evidence references unresolved belief: %s", ev.BeliefRef)
		}
		if _, err := c.solventClient.CreateEvidence(beliefID, ev.ProvenanceClass, ev.ContentSHA256); err != nil {
			return nil, fmt.Errorf("failed to persist evidence %s: %w", ev.LocalID, err)
		}
	}

	// 5. Persist tasks to Conductor
	for _, t := range pkt.Tasks {
		governanceRef := t.GovernanceRef
		if governanceRef == "" {
			governanceRef = fmt.Sprintf(`{"provider":"solvent","reference_id":"%s"}`, pkt.ScenarioID)
		}
		if _, err := c.conductorClient.CreateTask(pkt.ProjectRef, t.Title, t.Description, governanceRef); err != nil {
			return nil, fmt.Errorf("failed to persist task %s: %w", t.LocalID, err)
		}
	}

	return result, nil
}

// resolveRef resolves a packet reference to a Solvent belief ID.
// local:* references are resolved via the provided map.
// canonical:belief:* references are used as-is.
func resolveRef(ref string, beliefIDMap map[string]string) string {
	if strings.HasPrefix(ref, packetv1.RefPrefixLocal) {
		localID := strings.TrimPrefix(ref, packetv1.RefPrefixLocal)
		return beliefIDMap[localID]
	}
	if strings.HasPrefix(ref, packetv1.RefPrefixCanonical) {
		return strings.TrimPrefix(ref, packetv1.RefPrefixCanonical)
	}
	return ""
}
