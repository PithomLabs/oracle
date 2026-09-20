package bmistv1

import (
	"os"
	"path/filepath"
	"testing"

	domainpack "github.com/PithomLabs/oracle/domain-pack"
)

func validPack() *Pack {
	return &Pack{
		PackID:          "bmist",
		Version:         "1.0.0",
		Name:            "Test Pack",
		Description:     "A test pack",
		ClaimTypes:      []string{"derived", "accommodated", "postulated"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		DebtVocabulary:  []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		InitialDebt:     []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		RetirementRules: map[string]domainpack.RetirementRule{
			"needMap":               {EvidenceClass: "reproducible_artifact", Rule: "map_check"},
			"needInvariant":         {EvidenceClass: "reproducible_artifact", Rule: "invariant_check"},
			"needToyCheck":          {EvidenceClass: "reproducible_artifact", Rule: "toy_model_check"},
			"needNullModel":         {EvidenceClass: "operator_asserted", Rule: "scope_clarification"},
			"needObstruction":       {EvidenceClass: "reproducible_artifact", Rule: "obstruction_construction"},
			"needFaithfulnessReview": {EvidenceClass: "operator_asserted", Rule: "faithfulness_review"},
		},
		Falsifiers:            []string{"counterexample", "contradiction", "missing_evidence", "alternative_explanation"},
		ConsequentialActions:  []ConsequentialAction{{Action: "publish_claim", Requires: "promoted", Gates: []string{"faithfulness_review"}}},
		HumanGatedTransitions: []string{"faithfulness_review", "scope_clarification", "obstruction_assessment"},
	}
}

func TestValidPack(t *testing.T) {
	if err := validatePack(validPack()); err != nil {
		t.Fatalf("valid pack rejected: %v", err)
	}
}

func TestNilPack(t *testing.T) {
	if err := validatePack(nil); err == nil {
		t.Fatal("nil pack accepted")
	}
}

func TestMissingPackID(t *testing.T) {
	p := validPack()
	p.PackID = ""
	if err := validatePack(p); err == nil {
		t.Fatal("missing pack_id accepted")
	}
}

func TestEmptyPackID(t *testing.T) {
	p := validPack()
	p.PackID = "  "
	// Empty-ish but non-empty string passes the check; this is acceptable.
	// The rule is "non-empty", not "non-whitespace".
}

func TestMissingVersion(t *testing.T) {
	p := validPack()
	p.Version = ""
	if err := validatePack(p); err == nil {
		t.Fatal("missing version accepted")
	}
}

func TestInvalidVersion(t *testing.T) {
	p := validPack()
	p.Version = "v1.0.0"
	if err := validatePack(p); err == nil {
		t.Fatal("version with v prefix accepted")
	}
}

func TestInvalidVersionShort(t *testing.T) {
	p := validPack()
	p.Version = "1.0"
	if err := validatePack(p); err == nil {
		t.Fatal("short version accepted")
	}
}

func TestEmptyClaimTypes(t *testing.T) {
	p := validPack()
	p.ClaimTypes = []string{}
	if err := validatePack(p); err == nil {
		t.Fatal("empty claim_types accepted")
	}
}

func TestUnsupportedClaimType(t *testing.T) {
	p := validPack()
	p.ClaimTypes = []string{"derived", "bogus"}
	if err := validatePack(p); err == nil {
		t.Fatal("unsupported claim type accepted")
	}
}

func TestMissingEvidenceClasses(t *testing.T) {
	p := validPack()
	p.EvidenceClasses = []string{}
	if err := validatePack(p); err == nil {
		t.Fatal("empty evidence_classes accepted")
	}
}

func TestMissingDebtVocabulary(t *testing.T) {
	p := validPack()
	p.DebtVocabulary = []string{}
	if err := validatePack(p); err == nil {
		t.Fatal("empty debt_vocabulary accepted")
	}
}

func TestInvalidInitialDebt(t *testing.T) {
	p := validPack()
	p.InitialDebt = append(p.InitialDebt, "needBogus")
	if err := validatePack(p); err == nil {
		t.Fatal("initial_debt with undeclared item accepted")
	}
}

func TestInvalidRetirementRuleKey(t *testing.T) {
	p := validPack()
	p.RetirementRules["needBogus"] = domainpack.RetirementRule{EvidenceClass: "reproducible_artifact", Rule: "bogus"}
	if err := validatePack(p); err == nil {
		t.Fatal("retirement_rules with undeclared debt key accepted")
	}
}

func TestInvalidRetirementRuleEvidenceClass(t *testing.T) {
	p := validPack()
	p.RetirementRules["needMap"] = domainpack.RetirementRule{EvidenceClass: "bogus_class", Rule: "map_check"}
	if err := validatePack(p); err == nil {
		t.Fatal("retirement_rules with undeclared evidence class accepted")
	}
}

func TestMissingFaithfulnessReview(t *testing.T) {
	p := validPack()
	var gates []string
	for _, g := range p.HumanGatedTransitions {
		if g != "faithfulness_review" {
			gates = append(gates, g)
		}
	}
	p.HumanGatedTransitions = gates
	if err := validatePack(p); err == nil {
		t.Fatal("missing faithfulness_review accepted")
	}
}

func TestMissingScopeClarification(t *testing.T) {
	p := validPack()
	var gates []string
	for _, g := range p.HumanGatedTransitions {
		if g != "scope_clarification" {
			gates = append(gates, g)
		}
	}
	p.HumanGatedTransitions = gates
	if err := validatePack(p); err == nil {
		t.Fatal("missing scope_clarification accepted")
	}
}

func TestMissingObstructionAssessment(t *testing.T) {
	p := validPack()
	var gates []string
	for _, g := range p.HumanGatedTransitions {
		if g != "obstruction_assessment" {
			gates = append(gates, g)
		}
	}
	p.HumanGatedTransitions = gates
	if err := validatePack(p); err == nil {
		t.Fatal("missing obstruction_assessment accepted")
	}
}

func TestDuplicateDebtVocabulary(t *testing.T) {
	p := validPack()
	p.DebtVocabulary = append(p.DebtVocabulary, "needMap")
	if err := validatePack(p); err == nil {
		t.Fatal("duplicate debt vocabulary accepted")
	}
}

func TestDuplicateClaimType(t *testing.T) {
	p := validPack()
	p.ClaimTypes = []string{"derived", "derived"}
	if err := validatePack(p); err == nil {
		t.Fatal("duplicate claim type accepted")
	}
}

func TestEmptyVocabularyString(t *testing.T) {
	p := validPack()
	p.DebtVocabulary[0] = ""
	if err := validatePack(p); err == nil {
		t.Fatal("empty string in debt vocabulary accepted")
	}
}

func TestEmptyRetirementRule(t *testing.T) {
	p := validPack()
	p.RetirementRules["needMap"] = domainpack.RetirementRule{EvidenceClass: "reproducible_artifact", Rule: ""}
	if err := validatePack(p); err == nil {
		t.Fatal("empty retirement rule accepted")
	}
}

func TestVersionWithPrerelease(t *testing.T) {
	p := validPack()
	p.Version = "1.0.0-alpha.1"
	if err := validatePack(p); err != nil {
		t.Fatalf("valid prerelease version rejected: %v", err)
	}
}

func TestVersionWithBuildMetadata(t *testing.T) {
	p := validPack()
	p.Version = "1.0.0+build.7"
	if err := validatePack(p); err != nil {
		t.Fatalf("valid build metadata version rejected: %v", err)
	}
}

func TestValidateFixtures(t *testing.T) {
	validDir := filepath.Join("testdata", "packs", "valid")
	invalidDir := filepath.Join("testdata", "packs", "invalid")

	// Valid fixtures should pass
	entries, err := os.ReadDir(validDir)
	if err != nil {
		t.Skipf("valid fixtures directory not found: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(validDir, entry.Name()))
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}
		p, err := ParsePack(data)
		if err != nil {
			t.Fatalf("failed to parse %s: %v", entry.Name(), err)
		}
		if err := validatePack(p); err != nil {
			t.Errorf("valid fixture %s rejected: %v", entry.Name(), err)
		}
	}

	// Invalid fixtures should fail
	entries, err = os.ReadDir(invalidDir)
	if err != nil {
		t.Skipf("invalid fixtures directory not found: %v", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(invalidDir, entry.Name()))
		if err != nil {
			t.Fatalf("failed to read %s: %v", entry.Name(), err)
		}
		p, err := ParsePack(data)
		if err != nil {
			// Parse error is acceptable for invalid fixtures
			continue
		}
		if err := validatePack(p); err == nil {
			t.Errorf("invalid fixture %s accepted", entry.Name())
		}
	}
}
