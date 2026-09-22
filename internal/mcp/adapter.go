package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/PithomLabs/oracle/internal/application"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
)

// Adapter is the ARGUS MCP adapter exposing exactly two tools.
// It depends only on the application layer — no epistemic, work, or domain imports.
type Adapter struct {
	app *application.App
}

// NewAdapter creates a new MCP adapter.
func NewAdapter(app *application.App) *Adapter {
	return &Adapter{app: app}
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
		return a.handleGetContext(ctx, args)
	case ToolSubmitPacket:
		return a.handleSubmitPacket(ctx, args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// handleGetContext returns the RCP/v1 context for a given task.
func (a *Adapter) handleGetContext(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var req struct {
		TaskID string `json:"task_id"`
	}
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	if req.TaskID == "" {
		return nil, fmt.Errorf("task_id is required")
	}

	return a.app.GetContext(ctx, req.TaskID)
}

// handleSubmitPacket validates, compiles, and persists an EBP packet.
func (a *Adapter) handleSubmitPacket(ctx context.Context, args json.RawMessage) (interface{}, error) {
	var pkt packetv1.Packet
	if err := json.Unmarshal(args, &pkt); err != nil {
		return nil, fmt.Errorf("invalid packet: %w", err)
	}

	if err := a.app.Compile(ctx, &pkt); err != nil {
		return nil, fmt.Errorf("compile: %w", err)
	}
	if err := a.app.Validate(ctx, &pkt); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if err := a.app.ValidatePacket(ctx, &pkt); err != nil {
		return nil, fmt.Errorf("validate packet: %w", err)
	}

	// Begin transaction, persist, commit. Rollback on error.
	tx, err := a.app.DB().BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	result, err := a.app.Persist(ctx, tx, &pkt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return result, nil
}

// ListTools returns the two ARGUS tools.
func (a *Adapter) ListTools() []Tool {
	return []Tool{
		{
			Name:        ToolGetContext,
			Description: "Get the RCP/v1 context for a task: task, dependencies, epistemic state, and activity.",
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
			Description: "Submit an EBP research packet. Validates, persists beliefs/evidence/edges/tasks. Idempotent.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"schema_version": map[string]interface{}{"type": "string"},
					"role":           map[string]interface{}{"type": "string", "enum": []string{"work", "adversarial"}},
					"packet_id":      map[string]interface{}{"type": "string"},
					"pack_ref":       map[string]interface{}{"type": "string"},
					"scenario_id":    map[string]interface{}{"type": "string"},
					"beliefs":        map[string]interface{}{"type": "array"},
					"evidence":       map[string]interface{}{"type": "array"},
					"edges":          map[string]interface{}{"type": "array"},
					"tasks":          map[string]interface{}{"type": "array"},
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
