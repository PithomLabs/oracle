package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/PithomLabs/oracle/internal/mcp"
)

// jsonrpcRequest represents a JSON-RPC 2.0 request.
type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *interface{}    `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// jsonrpcResponse represents a JSON-RPC 2.0 response.
type jsonrpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      *interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

// rpcError represents a JSON-RPC 2.0 error object.
type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP server info returned during initialize.
var serverInfo = map[string]interface{}{
	"name":    "argus",
	"version": "0.1.0",
}

// MCP capabilities.
var capabilities = map[string]interface{}{
	"tools": map[string]interface{}{},
}

// runMCPStdio reads JSON-RPC messages from stdin and writes responses to stdout.
func runMCPStdio(adapter *mcp.Adapter) {
	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size for large packets (default 64KB, max 1MB).
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var req jsonrpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			log.Printf("mcp: invalid JSON-RPC: %v", err)
			continue
		}

		resp := handleMCPRequest(context.Background(), adapter, &req)

		// Notifications (no id) get no response.
		if req.ID == nil {
			continue
		}

		data, err := json.Marshal(resp)
		if err != nil {
			log.Printf("mcp: failed to marshal response: %v", err)
			continue
		}
		data = append(data, '\n')
		if _, err := os.Stdout.Write(data); err != nil {
			log.Printf("mcp: failed to write response: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		log.Printf("mcp: stdin read error: %v", err)
	}
}

// handleMCPRequest routes a JSON-RPC request to the appropriate handler.
func handleMCPRequest(ctx context.Context, adapter *mcp.Adapter, req *jsonrpcRequest) *jsonrpcResponse {
	resp := &jsonrpcResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	switch req.Method {
	case "initialize":
		resp.Result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    capabilities,
			"serverInfo":      serverInfo,
		}

	case "notifications/initialized":
		// Notification — no response needed (handled by caller skipping nil ID).

	case "tools/list":
		tools := adapter.ListTools()
		resp.Result = map[string]interface{}{
			"tools": tools,
		}

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			resp.Error = &rpcError{Code: -32602, Message: fmt.Sprintf("invalid params: %v", err)}
			return resp
		}

		result, err := adapter.HandleTool(ctx, params.Name, params.Arguments)
		if err != nil {
			resp.Result = map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": fmt.Sprintf("error: %v", err)},
				},
				"isError": true,
			}
			return resp
		}

		// Marshal result to JSON text content.
		resultJSON, err := json.Marshal(result)
		if err != nil {
			resp.Error = &rpcError{Code: -32603, Message: fmt.Sprintf("failed to marshal result: %v", err)}
			return resp
		}

		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{
				{"type": "text", "text": string(resultJSON)},
			},
		}

	default:
		resp.Error = &rpcError{Code: -32601, Message: fmt.Sprintf("method not found: %s", req.Method)}
	}

	return resp
}
