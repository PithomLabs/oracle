package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/PithomLabs/oracle/internal/application"
)

// ---- Refusal Test 2: MCP exposes exactly 2 tools ----

func TestMCPToolsCount(t *testing.T) {
	app := application.New(nil) // nil DB is fine — we're only testing tool listing
	adapter := NewAdapter(app)

	tools := adapter.ListTools()
	if len(tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(tools))
	}

	names := make(map[string]bool)
	for _, tool := range tools {
		names[tool.Name] = true
	}

	if !names[ToolGetContext] {
		t.Error("missing tool: argus.get_context")
	}
	if !names[ToolSubmitPacket] {
		t.Error("missing tool: argus.submit_packet")
	}

	// Verify no forbidden tools exist
	forbiddenTools := []string{"argus.promote", "argus.retract", "argus.discharge", "argus.authorize"}
	for _, ft := range forbiddenTools {
		if names[ft] {
			t.Errorf("forbidden tool found: %s", ft)
		}
	}

	t.Logf("MCP exposes exactly %d tools: %v", len(tools), toolNames(tools))
}

// ---- Refusal Test 5: Agent cannot retract ----

func TestAgentCannotRetract(t *testing.T) {
	app := application.New(nil)
	adapter := NewAdapter(app)

	// Verify no retract tool exists
	tools := adapter.ListTools()
	for _, tool := range tools {
		if tool.Name == "argus.retract" {
			t.Fatal("forbidden tool 'argus.retract' found in MCP tools")
		}
	}

	// HandleTool with unknown tool should return error
	_, err := adapter.HandleTool(context.Background(), "argus.retract", json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected error for unknown tool 'argus.retract', got nil")
	}
	t.Logf("agent retract correctly rejected: %v", err)
}

func toolNames(tools []Tool) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.Name
	}
	return names
}
