package coordinator

import (
	"fmt"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// ValidatePacket performs structural and semantic validation.
func ValidatePacket(pkt *packetv1.Packet, registry *domainpack.PackRegistry) error {
	return packetv1.Validate(pkt, registry)
}

// ValidateDebtMembership checks that all debt items are in the Pack vocabulary.
func ValidateDebtMembership(pkt *packetv1.Packet, pack domainpack.Pack) error {
	for i, b := range pkt.Beliefs {
		for _, debt := range b.Debt {
			if !packHasDebt(pack, debt) {
				return fmt.Errorf("belief[%d]: unknown debt item %q", i, debt)
			}
		}
	}
	return nil
}

// packHasDebt checks if a debt item exists in the pack vocabulary.
func packHasDebt(pack domainpack.Pack, debt string) bool {
	for _, d := range pack.GetDebtVocabulary() {
		if d == debt {
			return true
		}
	}
	return false
}

// ValidateDebtItemMembership checks that a single debt item is in the Pack vocabulary.
func ValidateDebtItemMembership(debtItem string, pack domainpack.Pack) error {
	if !packHasDebt(pack, debtItem) {
		return fmt.Errorf("unknown debt item %q: not in Domain Pack vocabulary", debtItem)
	}
	return nil
}

// ValidateEvidenceClass checks that evidence provenance class is supported by the Pack.
func ValidateEvidenceClass(pkt *packetv1.Packet, pack domainpack.Pack) error {
	supported := make(map[string]bool)
	for _, c := range pack.GetEvidenceClasses() {
		supported[c] = true
	}
	for i, e := range pkt.Evidence {
		if !supported[e.ProvenanceClass] {
			return fmt.Errorf("evidence[%d]: unsupported provenance_class %q", i, e.ProvenanceClass)
		}
	}
	return nil
}

// CompileDebt computes the final debt for a belief.
func CompileDebt(beliefDebt []string, packInitialDebt []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0)
	for _, d := range packInitialDebt {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}
	for _, d := range beliefDebt {
		if !seen[d] {
			seen[d] = true
			result = append(result, d)
		}
	}
	return result
}

// getRetirementRule extracts the retirement rule from a pack.
func getRetirementRule(debtItem string, pack domainpack.Pack) *domainpack.RetirementRule {
	rules := pack.GetRetirementRules()
	if rules == nil {
		return nil
	}
	rule, exists := rules[debtItem]
	if !exists {
		return nil
	}
	return &rule
}

// ValidateRetirementRule checks that the retirement request satisfies the Pack rule.
func ValidateRetirementRule(debtItem string, evidenceClass string, pack domainpack.Pack) error {
	rule := getRetirementRule(debtItem, pack)
	if rule == nil {
		return fmt.Errorf("no retirement rule defined for debt item %q", debtItem)
	}
	if rule.EvidenceClass != evidenceClass {
		return fmt.Errorf("retirement rule mismatch: debt %q requires evidence class %q, got %q", debtItem, rule.EvidenceClass, evidenceClass)
	}
	return nil
}
