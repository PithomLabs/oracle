package packetv1

import (
	"encoding/hex"
	"fmt"
	"strings"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
)

// Validate performs full structural and semantic validation of a packet
// against the pack registry. Validation is fail-closed.
func Validate(pkt *Packet, registry *domainpack.PackRegistry) error {
	if pkt == nil {
		return fmt.Errorf("packet is nil")
	}

	// 1. Schema version check
	if pkt.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported schema version: %s", pkt.SchemaVersion)
	}

	// 2. packet_id non-empty
	if pkt.PacketID == "" {
		return fmt.Errorf("packet_id is required")
	}

	// 3. role validation
	if pkt.Role != RoleWork && pkt.Role != RoleAdversarial {
		return fmt.Errorf("invalid role: %s (must be %q or %q)", pkt.Role, RoleWork, RoleAdversarial)
	}

	// 3.5. agent identity validation
	if err := validateAgent(pkt); err != nil {
		return err
	}

	// 4. pack_ref resolves through registry
	packID, packVersion := parsePackRef(pkt.PackRef)
	_, err := registry.Get(packID, packVersion)
	if err != nil {
		return fmt.Errorf("pack_ref resolution failed: %w", err)
	}

	// 5. Validate beliefs
	if err := validateBeliefs(pkt); err != nil {
		return err
	}

	// 6. Validate evidence
	if err := validateEvidence(pkt); err != nil {
		return err
	}

	// 7. Validate all references
	if err := validateReferences(pkt); err != nil {
		return err
	}

	// 8. Validate edges
	if err := validateEdges(pkt); err != nil {
		return err
	}

	// 9. Validate tasks
	if err := validateTasks(pkt); err != nil {
		return err
	}

	return nil
}


// validateAgent checks agent identity fields are present and consistent.
func validateAgent(pkt *Packet) error {
	if pkt.Agent.ID == "" {
		return fmt.Errorf("agent.id is required")
	}
	if pkt.Agent.Role == "" {
		return fmt.Errorf("agent.role is required")
	}
	if pkt.Agent.Harness == "" {
		return fmt.Errorf("agent.harness is required")
	}
	if pkt.Agent.Model == "" {
		return fmt.Errorf("agent.model is required")
	}
	if pkt.Agent.Role != pkt.Role {
		return fmt.Errorf("agent.role %q must equal packet.role %q", pkt.Agent.Role, pkt.Role)
	}
	return nil
}

// parsePackRef splits "bmist-v1" into packID and version.
func parsePackRef(ref string) (string, string) {
	parts := strings.SplitN(ref, "-", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return ref, ""
}

// validateBeliefs checks belief structure and uniqueness.
func validateBeliefs(pkt *Packet) error {
	seen := make(map[string]bool)
	for i, b := range pkt.Beliefs {
		if b.LocalID == "" {
			return fmt.Errorf("belief[%d]: local_id is required", i)
		}
		if seen[b.LocalID] {
			return fmt.Errorf("belief[%d]: duplicate local_id %q", i, b.LocalID)
		}
		seen[b.LocalID] = true

		if b.Claim == "" {
			return fmt.Errorf("belief[%d]: claim is required", i)
		}
		if b.ClaimType == "" {
			return fmt.Errorf("belief[%d]: claim_type is required", i)
		}
	}
	return nil
}

// validateEvidence checks evidence structure and content hashes.
func validateEvidence(pkt *Packet) error {
	seen := make(map[string]bool)
	for i, e := range pkt.Evidence {
		if e.LocalID == "" {
			return fmt.Errorf("evidence[%d]: local_id is required", i)
		}
		if seen[e.LocalID] {
			return fmt.Errorf("evidence[%d]: duplicate local_id %q", i, e.LocalID)
		}
		seen[e.LocalID] = true

		if e.BeliefRef == "" {
			return fmt.Errorf("evidence[%d]: belief_ref is required", i)
		}
		if e.ContentSHA256 == "" {
			return fmt.Errorf("evidence[%d]: content_sha256 is required", i)
		}
		if _, err := hex.DecodeString(e.ContentSHA256); err != nil || len(e.ContentSHA256) != 64 {
			return fmt.Errorf("evidence[%d]: invalid content_sha256 format", i)
		}
		if e.ProvenanceClass == "" {
			return fmt.Errorf("evidence[%d]: provenance_class is required", i)
		}
	}
	return nil
}

// validateReferences checks that all local/canonical references are well-formed.
func validateReferences(pkt *Packet) error {
	// Build local ID set
	localIDs := make(map[string]bool)
	for _, b := range pkt.Beliefs {
		localIDs[b.LocalID] = true
	}
	for _, e := range pkt.Evidence {
		localIDs[e.LocalID] = true
	}
	for _, edge := range pkt.Edges {
		localIDs[edge.LocalID] = true
	}
	for _, t := range pkt.Tasks {
		localIDs[t.LocalID] = true
	}

	// Validate evidence belief_refs
	for i, e := range pkt.Evidence {
		if _, err := ParseReference(e.BeliefRef); err != nil {
			return fmt.Errorf("evidence[%d].belief_ref: %w", i, err)
		}
	}

	// Validate edge references
	for i, edge := range pkt.Edges {
		if _, err := ParseReference(edge.FromRef); err != nil {
			return fmt.Errorf("edge[%d].from_ref: %w", i, err)
		}
		if _, err := ParseReference(edge.ToRef); err != nil {
			return fmt.Errorf("edge[%d].to_ref: %w", i, err)
		}
	}

	return nil
}

// validateEdges checks edge structure and kind values.
func validateEdges(pkt *Packet) error {
	seen := make(map[string]bool)
	for i, e := range pkt.Edges {
		if e.LocalID == "" {
			return fmt.Errorf("edge[%d]: local_id is required", i)
		}
		if seen[e.LocalID] {
			return fmt.Errorf("edge[%d]: duplicate local_id %q", i, e.LocalID)
		}
		seen[e.LocalID] = true

		if e.Kind != EdgeDerives && e.Kind != EdgeContradicts {
			return fmt.Errorf("edge[%d]: invalid kind %q (must be %q or %q)", i, e.Kind, EdgeDerives, EdgeContradicts)
		}
	}
	return nil
}

// validateTasks checks task structure.
func validateTasks(pkt *Packet) error {
	seen := make(map[string]bool)
	for i, t := range pkt.Tasks {
		if t.LocalID == "" {
			return fmt.Errorf("task[%d]: local_id is required", i)
		}
		if seen[t.LocalID] {
			return fmt.Errorf("task[%d]: duplicate local_id %q", i, t.LocalID)
		}
		seen[t.LocalID] = true

		if t.Title == "" {
			return fmt.Errorf("task[%d]: title is required", i)
		}
	}
	return nil
}
