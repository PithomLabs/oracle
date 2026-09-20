package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	coordinator "github.com/PithomLabs/oracle/coordinator"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// Adapter is the ARGUS MCP adapter exposing exactly two tools.
type Adapter struct {
	coordinator *coordinator.Coordinator
}

// NewAdapter creates a new MCP adapter.
func NewAdapter(c *coordinator.Coordinator) *Adapter {
	return &Adapter{coordinator: c}
}

// Tool names
const (
	ToolGetContext   = "argus.get_context"
	ToolSubmitPacket = "argus.submit_packet"
)

// HandleTool dispatches a tool call to the appropriate handler.
func (a *Adapter) HandleTool(ctx context.Context, name string, args json.RawMessage) (interface{}, error) {
	switch name {
	case ToolGetContext:
		return a.handleGetContext(args)
	case ToolSubmitPacket:
		return a.handleSubmitPacket(args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// handleGetContext returns the RCP/v1 context for a given task.
func (a *Adapter) handleGetContext(args json.RawMessage) (interface{}, error) {
	var req struct {
		TaskID string `json:"task_id"`
	}
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if req.TaskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	ctx, err := a.coordinator.GetContext(req.TaskID)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// handleSubmitPacket persists an EBP packet to Solvent and Conductor.
func (a *Adapter) handleSubmitPacket(args json.RawMessage) (interface{}, error) {
	var pkt packetv1.Packet
	if err := json.Unmarshal(args, &pkt); err != nil {
		return nil, fmt.Errorf("invalid packet: %w", err)
	}

	result, err := a.coordinator.SubmitPacket(&pkt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListTools returns the two ARGUS tools.
func (a *Adapter) ListTools() []Tool {
	return []Tool{
		{
			Name:        ToolGetContext,
			Description: "Get the RCP/v1 context for a task: task, dependencies, epistemic state (beliefs, evidence, edges, debt, intents), and activity. Returns full scenario projection with availability metadata.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"task_id": map[string]interface{}{
						"type":        "string",
						"description": "UUID of the Conductor task",
					},
				},
				"required": []string{"task_id"},
			},
		},
		{
			Name:        ToolSubmitPacket,
			Description: "Submit an EBP research packet. Validates, persists beliefs/evidence/edges/tasks to Solvent and Conductor. Idempotent: same content + same scenario deduplicates.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"schema_version": map[string]interface{}{
						"type": "string",
					},
					"role": map[string]interface{}{
						"type": "string",
						"enum": []string{"work", "adversarial"},
					},
					"packet_id": map[string]interface{}{
						"type": "string",
					},
					"pack_ref": map[string]interface{}{
						"type": "string",
					},
					"scenario_id": map[string]interface{}{
						"type": "string",
					},
					"beliefs": map[string]interface{}{
						"type": "array",
					},
					"evidence": map[string]interface{}{
						"type": "array",
					},
					"edges": map[string]interface{}{
						"type": "array",
					},
					"tasks": map[string]interface{}{
						"type": "array",
					},
				},
				"required": []string{"schema_version", "role", "packet_id", "pack_ref", "beliefs", "evidence"},
			},
		},
	}
}

// Tool represents an MCP tool definition.
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}
