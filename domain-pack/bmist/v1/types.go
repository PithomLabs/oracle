package bmistv1

import (
	"encoding/json"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
)

// Pack is the top-level Domain Pack descriptor.
type Pack struct {
	PackID                string                              `json:"pack_id"`
	Version               string                              `json:"version"`
	Name                  string                              `json:"name"`
	Description           string                              `json:"description"`
	ClaimTypes            []string                            `json:"claim_types"`
	EvidenceClasses       []string                            `json:"evidence_classes"`
	DebtVocabulary        []string                            `json:"debt_vocabulary"`
	InitialDebt           []string                            `json:"initial_debt"`
	RetirementRules       map[string]domainpack.RetirementRule `json:"retirement_rules"`
	Falsifiers            []string                            `json:"falsifiers"`
	ConsequentialActions  []ConsequentialAction               `json:"consequential_actions"`
	HumanGatedTransitions []string                            `json:"human_gated_transitions"`
	VerifierSpecs         []domainpack.VerifierSpec           `json:"verifier_specs"`
}

// GetPackID returns the pack identifier.
func (p *Pack) GetPackID() string { return p.PackID }

// GetVersion returns the pack version.
func (p *Pack) GetVersion() string { return p.Version }

// GetDebtVocabulary returns the pack's debt vocabulary.
func (p *Pack) GetDebtVocabulary() []string { return p.DebtVocabulary }

// GetEvidenceClasses returns the pack's supported evidence classes.
func (p *Pack) GetEvidenceClasses() []string { return p.EvidenceClasses }

// GetRetirementRules returns the pack's retirement rules.
func (p *Pack) GetRetirementRules() map[string]domainpack.RetirementRule {
	return p.RetirementRules
}

// GetVerifierSpecs returns the pack's authorized verifier specs.
func (p *Pack) GetVerifierSpecs() []domainpack.VerifierSpec { return p.VerifierSpecs }

// Validate implements the Validator interface.
func (p *Pack) Validate() error {
	return validatePack(p)
}

// ConsequentialAction declares an action that requires promoted belief state.
type ConsequentialAction struct {
	Action   string   `json:"action"`
	Requires string   `json:"requires"`
	Gates    []string `json:"gates"`
}

// ParsePack decodes JSON bytes into a Pack.
func ParsePack(data []byte) (*Pack, error) {
	var p Pack
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
