package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestContractPinPresent(t *testing.T) {
	path := "contract.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("contract.json does not exist at %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read contract.json: %v", err)
	}

	var contract map[string]interface{}
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatalf("contract.json is not valid JSON: %v", err)
	}

	required := []string{
		"workflow_spec_version",
		"role_boundary_matrix_version",
		"operation_identity_rule",
		"declaration_owner",
		"declaration_version",
		"pass_criteria",
		"sandbox_rule",
		"failure_taxonomy",
		"finding_arbiter",
		"bmist_status",
	}

	for _, key := range required {
		val, ok := contract[key]
		if !ok {
			t.Errorf("contract.json missing required field: %s", key)
			continue
		}
		if isBlank(val) {
			t.Errorf("contract.json field %s is blank", key)
		}
	}
}

func isBlank(v interface{}) bool {
	switch vv := v.(type) {
	case string:
		return vv == ""
	case []interface{}:
		return len(vv) == 0
	case map[string]interface{}:
		return len(vv) == 0
	default:
		return false
	}
}
