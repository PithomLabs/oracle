package adapter_test

import (
	"encoding/json"
	"testing"

	coordinator "github.com/PithomLabs/oracle/coordinator"
	"github.com/PithomLabs/oracle/mcp/adapter"
)

func TestListTools_TwoARGUSToolsOnly(t *testing.T) {
	solvent := &coordinator.MockSolventClient{}
	conductor := &coordinator.MockConductorClient{}
	c := coordinator.NewWithMockClients(solvent, conductor, "test-operator")
	a := adapter.NewAdapter(c)

	tools := a.ListTools()

	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}

	names := make(map[string]bool)
	for _, tool := range tools {
		names[tool.Name] = true
	}

	if !names["argus.get_context"] {
		t.Error("argus.get_context not found")
	}
	if !names["argus.submit_packet"] {
		t.Error("argus.submit_packet not found")
	}
}

func TestHandleTool_GetContext(t *testing.T) {
	solvent := &coordinator.MockSolventClient{}
	conductor := &coordinator.MockConductorClient{}
	c := coordinator.NewWithMockClients(solvent, conductor, "test-operator")
	a := adapter.NewAdapter(c)

	args, _ := json.Marshal(map[string]string{"task_id": "test-task"})
	result, err := a.HandleTool(nil, "argus.get_context", args)
	if err != nil {
		t.Fatalf("HandleTool failed: %v", err)
	}

	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestHandleTool_UnknownTool(t *testing.T) {
	solvent := &coordinator.MockSolventClient{}
	conductor := &coordinator.MockConductorClient{}
	c := coordinator.NewWithMockClients(solvent, conductor, "test-operator")
	a := adapter.NewAdapter(c)

	_, err := a.HandleTool(nil, "unknown.tool", nil)
	if err == nil {
		t.Error("expected error for unknown tool")
	}
}

func TestHandleTool_MissingTaskID(t *testing.T) {
	solvent := &coordinator.MockSolventClient{}
	conductor := &coordinator.MockConductorClient{}
	c := coordinator.NewWithMockClients(solvent, conductor, "test-operator")
	a := adapter.NewAdapter(c)

	args, _ := json.Marshal(map[string]string{})
	_, err := a.HandleTool(nil, "argus.get_context", args)
	if err == nil {
		t.Error("expected error for missing task_id")
	}
}
