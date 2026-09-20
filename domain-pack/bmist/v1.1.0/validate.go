package bmistv11

import (
	"fmt"
	"regexp"
)

var semverRegex = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?(\+[a-zA-Z0-9.]+)?$`)

var requiredHumanGates = map[string]bool{
	"faithfulness_review":    true,
	"scope_clarification":    true,
	"obstruction_assessment": true,
}

var supportedClaimTypes = map[string]bool{
	"derived":      true,
	"accommodated": true,
	"postulated":   true,
}

// Validate checks the structural integrity of a Pack.
// It does NOT prove scientific correctness, evaluate physics, or decide
// whether debt is actually retired.
func validatePack(p *Pack) error {
	if p == nil {
		return fmt.Errorf("pack is nil")
	}

	// 1. pack_id non-empty
	if p.PackID == "" {
		return fmt.Errorf("pack_id is required and must be non-empty")
	}

	// 2. version valid SemVer 2.0.0 (no v prefix)
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	if !semverRegex.MatchString(p.Version) {
		return fmt.Errorf("version %q is not valid SemVer 2.0.0 (no v prefix allowed)", p.Version)
	}

	// 3. claim_types non-empty
	if len(p.ClaimTypes) == 0 {
		return fmt.Errorf("claim_types must be non-empty")
	}

	// 4. claim_types ⊆ {derived, accommodated, postulated}
	for _, ct := range p.ClaimTypes {
		if !supportedClaimTypes[ct] {
			return fmt.Errorf("unsupported claim type %q", ct)
		}
	}

	// 5. evidence_classes non-empty
	if len(p.EvidenceClasses) == 0 {
		return fmt.Errorf("evidence_classes must be non-empty")
	}

	// 6. debt_vocabulary non-empty
	if len(p.DebtVocabulary) == 0 {
		return fmt.Errorf("debt_vocabulary must be non-empty")
	}

	// 7. initial_debt ⊆ debt_vocabulary
	debtSet := make(map[string]bool, len(p.DebtVocabulary))
	for _, d := range p.DebtVocabulary {
		debtSet[d] = true
	}
	for _, d := range p.InitialDebt {
		if !debtSet[d] {
			return fmt.Errorf("initial_debt contains undeclared item %q", d)
		}
	}

	// 8. retirement_rules keys ⊆ debt_vocabulary
	for key := range p.RetirementRules {
		if !debtSet[key] {
			return fmt.Errorf("retirement_rules contains undeclared debt key %q", key)
		}
	}

	// 9. each retirement rule's evidence_class ∈ declared evidence_classes
	evidenceSet := make(map[string]bool, len(p.EvidenceClasses))
	for _, ec := range p.EvidenceClasses {
		evidenceSet[ec] = true
	}
	for key, rule := range p.RetirementRules {
		if !evidenceSet[rule.EvidenceClass] {
			return fmt.Errorf("retirement_rules[%q] references undeclared evidence class %q", key, rule.EvidenceClass)
		}
		if rule.Rule == "" {
			return fmt.Errorf("retirement_rules[%q] has empty rule", key)
		}
	}

	// 10. falsifiers non-empty
	if len(p.Falsifiers) == 0 {
		return fmt.Errorf("falsifiers must be non-empty")
	}

	// 11. consequential_actions structurally valid
	for i, ca := range p.ConsequentialActions {
		if ca.Action == "" {
			return fmt.Errorf("consequential_actions[%d] has empty action", i)
		}
		if ca.Requires == "" {
			return fmt.Errorf("consequential_actions[%d] has empty requires", i)
		}
	}

	// 12. human_gated_transitions ⊇ required floor
	gateSet := make(map[string]bool, len(p.HumanGatedTransitions))
	for _, g := range p.HumanGatedTransitions {
		gateSet[g] = true
	}
	for required := range requiredHumanGates {
		if !gateSet[required] {
			return fmt.Errorf("human_gated_transitions is missing required gate %q", required)
		}
	}

	// 13. no duplicate values in arrays
	if err := checkNoDuplicates("claim_types", p.ClaimTypes); err != nil {
		return err
	}
	if err := checkNoDuplicates("evidence_classes", p.EvidenceClasses); err != nil {
		return err
	}
	if err := checkNoDuplicates("debt_vocabulary", p.DebtVocabulary); err != nil {
		return err
	}
	if err := checkNoDuplicates("initial_debt", p.InitialDebt); err != nil {
		return err
	}
	if err := checkNoDuplicates("falsifiers", p.Falsifiers); err != nil {
		return err
	}
	if err := checkNoDuplicates("human_gated_transitions", p.HumanGatedTransitions); err != nil {
		return err
	}

	// 14. no empty strings in arrays or required identifiers
	if err := checkNoEmptyStrings("claim_types", p.ClaimTypes); err != nil {
		return err
	}
	if err := checkNoEmptyStrings("evidence_classes", p.EvidenceClasses); err != nil {
		return err
	}
	if err := checkNoEmptyStrings("debt_vocabulary", p.DebtVocabulary); err != nil {
		return err
	}
	if err := checkNoEmptyStrings("initial_debt", p.InitialDebt); err != nil {
		return err
	}
	if err := checkNoEmptyStrings("falsifiers", p.Falsifiers); err != nil {
		return err
	}
	if err := checkNoEmptyStrings("human_gated_transitions", p.HumanGatedTransitions); err != nil {
		return err
	}

	return nil
}

func checkNoDuplicates(field string, items []string) error {
	seen := make(map[string]bool, len(items))
	for _, item := range items {
		if seen[item] {
			return fmt.Errorf("%s contains duplicate value %q", field, item)
		}
		seen[item] = true
	}
	return nil
}

func checkNoEmptyStrings(field string, items []string) error {
	for i, item := range items {
		if item == "" {
			return fmt.Errorf("%s contains empty string at index %d", field, i)
		}
	}
	return nil
}
