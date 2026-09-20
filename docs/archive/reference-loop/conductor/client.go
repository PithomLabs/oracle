package conductor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type MCPRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int64                  `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc,omitempty"`
	ID      int64       `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Client struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	mu      sync.Mutex
	nextID  int64
	closed  bool
	agentID string
	reader  *bufio.Reader
}

type Task struct {
	ID            string  `json:"id"`
	ProjectID     string  `json:"project_id"`
	Title         string  `json:"title"`
	Description   string  `json:"description,omitempty"`
	Status        string  `json:"status"`
	Priority      string  `json:"priority"`
	CurrentAgent  *string `json:"current_agent,omitempty"`
	GovernanceRef *string `json:"governance_ref,omitempty"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type Activity struct {
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	ActorType string `json:"actor_type"`
	ActorID   string `json:"actor_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	CreatedAt string `json:"created_at"`
}

func NewClient(dbPath, agentID string) (*Client, error) {
	cmd := exec.Command("/tmp/conductor", "-mode", "mcp", "-db", dbPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start conductor: %w", err)
	}
	return &Client{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		reader:  bufio.NewReader(stdout),
		agentID: agentID,
	}, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	c.stdin.Close()
	c.stdout.Close()
	_ = c.cmd.Process.Kill()
	_ = c.cmd.Wait()
	return nil
}

func (c *Client) callTool(ctx context.Context, toolName string, args map[string]interface{}) (map[string]interface{}, error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil, fmt.Errorf("client closed")
	}
	id := c.nextID
	c.nextID++
	c.mu.Unlock()

	req := MCPRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		},
	}
	if args == nil {
		req.Params["arguments"] = map[string]interface{}{}
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	data = append(data, '\n')

	if _, err := c.stdin.Write(data); err != nil {
		return nil, fmt.Errorf("write request: %w", err)
	}

	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var resp MCPResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("mcp error: %s", resp.Error)
	}
	if resp.Result == nil {
		return nil, fmt.Errorf("empty result")
	}
	result, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected result type: %T", resp.Result)
	}
	return result, nil
}

func (c *Client) CreateTask(ctx context.Context, projectID, title, description, governanceRef string) (*Task, error) {
	args := map[string]interface{}{
		"project_id": projectID,
		"title":      title,
	}
	if description != "" {
		args["description"] = description
	}
	if governanceRef != "" {
		args["governance_ref"] = governanceRef
	}
	result, err := c.callTool(ctx, "conductor_create_task", args)
	if err != nil {
		return nil, err
	}
	task, err := mapToTask(result)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (c *Client) ClaimTask(ctx context.Context, taskID string) (*Task, error) {
	result, err := c.callTool(ctx, "conductor_claim_task", map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return nil, err
	}
	task, err := mapToTask(result)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (c *Client) SubmitTask(ctx context.Context, taskID string) (*Task, error) {
	result, err := c.callTool(ctx, "conductor_submit_task", map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return nil, err
	}
	task, err := mapToTask(result)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (c *Client) PostActivity(ctx context.Context, taskID, action, details string) (*Activity, error) {
	result, err := c.callTool(ctx, "conductor_post_activity", map[string]interface{}{
		"task_id": taskID,
		"action":  action,
		"details": details,
	})
	if err != nil {
		return nil, err
	}
	activity, err := mapToActivity(result)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (c *Client) NextTask(ctx context.Context, projectID string) (*Task, error) {
	result, err := c.callTool(ctx, "conductor_next_task", map[string]interface{}{
		"project_id": projectID,
	})
	if err != nil {
		return nil, err
	}
	task, err := mapToTask(result)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (c *Client) GetTask(ctx context.Context, taskID string) (*Task, error) {
	result, err := c.callTool(ctx, "conductor_get_task", map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return nil, err
	}
	task, err := mapToTask(result)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func mapToTask(m map[string]interface{}) (*Task, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func mapToActivity(m map[string]interface{}) (*Activity, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var activity Activity
	if err := json.Unmarshal(data, &activity); err != nil {
		return nil, err
	}
	return &activity, nil
}

func (c *Client) GetGovernance(ctx context.Context, taskID string) (map[string]interface{}, error) {
	result, err := c.callTool(ctx, "conductor_get_governance", map[string]interface{}{
		"task_id": taskID,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
