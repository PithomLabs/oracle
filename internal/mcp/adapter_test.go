package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/PithomLabs/oracle/internal/application"
	"github.com/PithomLabs/oracle/internal/migrations"
	solventmigrations "github.com/PithomLabs/oracle/internal/solventmigrations"
	packetv1 "github.com/PithomLabs/oracle/packet/v1"
	domainpack "github.com/PithomLabs/oracle/domain-pack"
	bmistv1 "github.com/PithomLabs/oracle/domain-pack/bmist/v1"
	_ "github.com/jackc/pgx/v5/stdlib"
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

// ---- MCP-path validation: malformed packets rejected before Persist ----

func newTestAdapter(t *testing.T) *Adapter {
	t.Helper()
	app := application.New(nil) // nil DB — Persist is never reached
	registry := domainpack.NewRegistry()
	pack := &bmistv1.Pack{
		PackID:          "bmist",
		Version:         "1.0.0",
		ClaimTypes:      []string{"derived", "accommodated", "postulated"},
		EvidenceClasses: []string{"reproducible_artifact", "operator_asserted"},
		DebtVocabulary:  []string{"needMap", "needInvariant", "needToyCheck", "needNullModel", "needObstruction", "needFaithfulnessReview"},
		Falsifiers:              []string{"counterexample", "contradiction"},
		HumanGatedTransitions:   []string{"faithfulness_review", "scope_clarification", "obstruction_assessment"},
		ConsequentialActions:    []bmistv1.ConsequentialAction{{Action: "publish_claim", Requires: "promoted", Gates: []string{"faithfulness_review"}}},
	}
	if err := registry.Register(pack); err != nil {
		t.Fatalf("failed to register test pack: %v", err)
	}
	app.SetPackRegistry(registry)
	return NewAdapter(app)
}

func validMCPPacket() map[string]interface{} {
	return map[string]interface{}{
		"schema_version": packetv1.SchemaVersion,
		"role":           "work",
		"packet_id":      "test-mcp-001",
		"pack_ref":       "bmist@1.0.0",
		"scenario_id":    "00000000-0000-0000-0000-000000000001",
		"agent": map[string]interface{}{
			"id":      "test-agent",
			"role":    "work",
			"harness": "opencode",
			"model":   "test-model",
		},
		"beliefs": []map[string]interface{}{
			{"local_id": "b1", "claim": "test claim", "claim_type": "derived"},
		},
		"evidence": []map[string]interface{}{
			{
				"local_id":          "e1",
				"belief_ref":        "local:b1",
				"provenance_class":  "reproducible_artifact",
				"content_sha256":    "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
			},
		},
	}
}

// fast validation-path unit test using nil DB.
// Not production integration coverage.
// The DB-backed MCP tests (in app_test.go) are the production-path proof.
func TestMCPValidationRejectsMalformed_FastUnit(t *testing.T) {
	adapter := newTestAdapter(t)

	tests := []struct {
		name    string
		pkt     map[string]interface{}
		wantErr string
	}{
		{
			name: "missing agent.id",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				agent := p["agent"].(map[string]interface{})
				agent["id"] = ""
				return p
			}(),
			wantErr: "agent.id is required",
		},
		{
			name: "role mismatch",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				p["role"] = "adversarial"
				return p
			}(),
			wantErr: "agent.role",
		},
		{
			name: "duplicate belief local_id",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				p["beliefs"] = []map[string]interface{}{
					{"local_id": "b1", "claim": "first", "claim_type": "derived"},
					{"local_id": "b1", "claim": "second", "claim_type": "derived"},
				}
				return p
			}(),
			wantErr: "duplicate local_id",
		},
		{
			name: "invalid edge kind",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				p["edges"] = []map[string]interface{}{
					{"local_id": "ed1", "from_ref": "local:b1", "to_ref": "local:b1", "kind": "invalid_kind"},
				}
				return p
			}(),
			wantErr: "invalid kind",
		},
		{
			name: "invalid claim_type",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				p["beliefs"] = []map[string]interface{}{
					{"local_id": "b1", "claim": "test", "claim_type": "pizza"},
				}
				return p
			}(),
			wantErr: "invalid claim_type",
		},
		{
			name: "cross-type local_id collision",
			pkt: func() map[string]interface{} {
				p := validMCPPacket()
				p["evidence"] = []map[string]interface{}{
					{
						"local_id":          "b1",
						"belief_ref":        "local:b1",
						"provenance_class":  "reproducible_artifact",
						"content_sha256":    "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
					},
				}
				return p
			}(),
			wantErr: "duplicate local_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := json.Marshal(tt.pkt)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			_, err = adapter.HandleTool(context.Background(), ToolSubmitPacket, args)
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func toolNames(tools []Tool) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.Name
	}
	return names
}

// ---- Authority boundary: MCP exposes exactly 2 tools, no authority-changing tools ----
//
// Agent authority is bounded by MCP tool-surface exposure (tested here).
// Human/control authority uses the existing operator/control path (not tested here).
// This POC freeze does not claim protection against an arbitrary network
// actor that can directly reach an internal authority endpoint.
// Revisit before non-loopback or multi-actor deployment.

func TestMCPDoesNotExposeAuthorityTools(t *testing.T) {
	app := application.New(nil)
	adapter := NewAdapter(app)

	tools := adapter.ListTools()
	if len(tools) != 2 {
		t.Fatalf("expected exactly 2 MCP tools, got %d: %v", len(tools), toolNames(tools))
	}

	// Assert exact tool set
	expectedTools := map[string]bool{
		ToolGetContext:   false,
		ToolSubmitPacket: false,
	}
	for _, tool := range tools {
		if _, ok := expectedTools[tool.Name]; !ok {
			t.Errorf("unexpected tool in MCP surface: %s", tool.Name)
		}
		expectedTools[tool.Name] = true
	}
	for name, found := range expectedTools {
		if !found {
			t.Errorf("missing expected tool: %s", name)
		}
	}

	// Authority-changing tools must not be exposed
	forbiddenTools := []string{"argus.promote", "argus.retract", "argus.discharge", "argus.authorize"}
	for _, ft := range forbiddenTools {
		_, err := adapter.HandleTool(context.Background(), ft, json.RawMessage(`{}`))
		if err == nil {
			t.Errorf("authority tool %s should return error, got nil", ft)
		}
	}

	t.Logf("MCP authority boundary verified: exactly tools=%v, no authority-changing tools", toolNames(tools))
}

// ---- DB-backed MCP get_context transport test (freeze-critical) ----

func mcpTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("ARGUS_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
	}

	adminDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open admin db: %v", err)
	}

	dbName := fmt.Sprintf("mcp_%s_test", t.Name())
	_, _ = adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
	_, err = adminDB.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %q", dbName))
	if err != nil {
		adminDB.Close()
		t.Fatalf("create database: %v", err)
	}

	testDSN := fmt.Sprintf("postgres://root@localhost:26257/%s?sslmode=disable", dbName)
	conn, err := sql.Open("pgx", testDSN)
	if err != nil {
		adminDB.Close()
		t.Fatalf("open test db: %v", err)
	}

	if err := solventmigrations.Apply(context.Background(), conn); err != nil {
		conn.Close()
		adminDB.Close()
		t.Fatalf("apply solvent migrations: %v", err)
	}
	if err := migrations.Apply(context.Background(), conn); err != nil {
		conn.Close()
		adminDB.Close()
		t.Fatalf("apply oracle migrations: %v", err)
	}

	t.Cleanup(func() {
		conn.Close()
		adminDB.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %q CASCADE", dbName))
		adminDB.Close()
	})
	return conn
}

// normalizeUUID converts a flat hex EntityID to CRDB's hyphenated UUID format.
func mcpNormalizeUUID(id string) string {
	if len(id) == 32 {
		return id[:8] + "-" + id[8:12] + "-" + id[12:16] + "-" + id[16:20] + "-" + id[20:]
	}
	return id
}

func TestMCPGetContextReadsBackPersistedState(t *testing.T) {
	db := mcpTestDB(t)
	ctx := context.Background()
	app := application.New(db)
	adapter := NewAdapter(app)

	scenarioID := "00000000-0000-0000-0000-0000000000A0"
	projectID := scenarioID

	// Create project
	_, err := db.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ($1, 'mcp-transport-test', 'active')
		 ON CONFLICT DO NOTHING`, projectID)
	if err != nil {
		t.Fatalf("insert project: %v", err)
	}

	// Create task
	taskID := application.EntityID(scenarioID, "task", "mcp-transport-test-task")
	_, err = db.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, description, status, priority, created_at, updated_at)
		 VALUES ($1, $2, 'mcp transport test', 'test desc', 'proposed', 'medium', now(), now())
		 ON CONFLICT DO NOTHING`, taskID, projectID)
	if err != nil {
		t.Fatalf("insert task: %v", err)
	}

	// --- P1: work agent, 1 belief + 1 evidence ---
	p1PacketID := "mcp-pkt-p1"
	p1 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      p1PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p1", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "transport claim P1", ClaimType: "derived", Debt: []string{"needMap"}},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"},
		},
	}
	tx1, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx1: %v", err)
	}
	res1, err := app.Persist(ctx, tx1, p1)
	if err != nil {
		tx1.Rollback()
		t.Fatalf("persist P1: %v", err)
	}
	if err := tx1.Commit(); err != nil {
		t.Fatalf("commit P1: %v", err)
	}
	belief1ID := res1.BeliefIDs["b1"]

	// --- P2: adversarial agent, 1 belief + 1 evidence ---
	p2PacketID := "mcp-pkt-p2"
	p2 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleAdversarial,
		PacketID:      p2PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p2", Role: packetv1.RoleAdversarial, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b1", Claim: "transport challenge P2", ClaimType: "derived"},
		},
		Evidence: []packetv1.Evidence{
			{LocalID: "e1", BeliefRef: "local:b1", ProvenanceClass: "external_feed", ContentSHA256: "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3"},
		},
	}
	tx2, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx2: %v", err)
	}
	res2, err := app.Persist(ctx, tx2, p2)
	if err != nil {
		tx2.Rollback()
		t.Fatalf("persist P2: %v", err)
	}
	if err := tx2.Commit(); err != nil {
		t.Fatalf("commit P2: %v", err)
	}
	belief2ID := res2.BeliefIDs["b1"]

	// --- P3: work agent, 1 edge (contradicts) + witness belief ---
	p3PacketID := "mcp-pkt-p3"
	p3 := &packetv1.Packet{
		SchemaVersion: packetv1.SchemaVersion,
		Role:          packetv1.RoleWork,
		PacketID:      p3PacketID,
		PackRef:       "bmist@1.0.0",
		ScenarioID:    scenarioID,
		Agent:         packetv1.Agent{ID: "agent-p1", Role: packetv1.RoleWork, Harness: "test", Model: "test-model"},
		Beliefs: []packetv1.Belief{
			{LocalID: "b-witness", Claim: "edge witness", ClaimType: "derived"},
		},
		Edges: []packetv1.Edge{
			{LocalID: "ed1", FromRef: "canonical:belief:" + belief1ID, ToRef: "canonical:belief:" + belief2ID, Kind: "contradicts"},
		},
	}
	tx3, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx3: %v", err)
	}
	_, err = app.Persist(ctx, tx3, p3)
	if err != nil {
		tx3.Rollback()
		t.Fatalf("persist P3: %v", err)
	}
	if err := tx3.Commit(); err != nil {
		t.Fatalf("commit P3: %v", err)
	}

	// === CALL THE PRODUCTION MCP PATH ===
	argsJSON, err := json.Marshal(map[string]string{"task_id": taskID})
	if err != nil {
		t.Fatalf("marshal args: %v", err)
	}
	result, err := adapter.HandleTool(ctx, ToolGetContext, json.RawMessage(argsJSON))
	if err != nil {
		t.Fatalf("HandleTool get_context: %v", err)
	}

	// Exercise the JSON serialization path (same as cmd/argus/mcp_stdio.go)
	raw, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal get_context result: %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("marshal produced empty output")
	}

	// Decode into generic map for field assertions
	var decoded map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal get_context result: %v", err)
	}

	// Assert task
	task, ok := decoded["task"].(map[string]interface{})
	if !ok {
		t.Fatal("task not found or wrong type in result")
	}
	if !strings.EqualFold(task["id"].(string), mcpNormalizeUUID(taskID)) {
		t.Errorf("task.id = %v, want %v", task["id"], mcpNormalizeUUID(taskID))
	}
	if !strings.EqualFold(task["project_id"].(string), projectID) {
		t.Errorf("task.project_id = %v, want %v", task["project_id"], projectID)
	}

	// Assert availability
	avail, ok := decoded["availability"].(map[string]interface{})
	if !ok {
		t.Fatal("availability not found or wrong type")
	}
	taskAvail, ok := avail["task"].(map[string]interface{})
	if !ok {
		t.Fatal("availability.task not found")
	}
	if taskAvail["available"] != true {
		t.Errorf("availability.task.available = %v, want true", taskAvail["available"])
	}
	snapAvail, ok := avail["snapshot"].(map[string]interface{})
	if !ok {
		t.Fatal("availability.snapshot not found")
	}
	if snapAvail["available"] != true {
		t.Errorf("availability.snapshot.available = %v, want true", snapAvail["available"])
	}

	// Assert snapshot
	snapshot, ok := decoded["snapshot"].(map[string]interface{})
	if !ok {
		t.Fatal("snapshot not found or wrong type")
	}

	// Assert beliefs (expect 3: P1 + P2 + P3 witness)
	beliefs, ok := snapshot["beliefs"].([]interface{})
	if !ok {
		t.Fatal("beliefs not found or wrong type")
	}
	if len(beliefs) != 3 {
		t.Fatalf("belief count = %d, want 3", len(beliefs))
	}

	beliefMap := make(map[string]map[string]interface{})
	for _, b := range beliefs {
		bm, ok := b.(map[string]interface{})
		if !ok {
			t.Fatal("belief entry wrong type")
		}
		beliefMap[bm["claim"].(string)] = bm
	}

	// P1 belief
	b1, ok := beliefMap["transport claim P1"]
	if !ok {
		t.Fatal("P1 belief not found by claim")
	}
	if b1["claim_type"] != "derived" {
		t.Errorf("P1 claim_type = %v, want derived", b1["claim_type"])
	}
	if b1["origin_packet_id"] != p1PacketID {
		t.Errorf("P1 origin_packet_id = %v, want %v", b1["origin_packet_id"], p1PacketID)
	}
	debt1, ok := b1["debt"].([]interface{})
	if !ok {
		t.Fatal("P1 debt not found or wrong type")
	}
	if len(debt1) != 1 || debt1[0].(string) != "needMap" {
		t.Errorf("P1 debt = %v, want [needMap]", debt1)
	}

	// P2 belief
	b2, ok := beliefMap["transport challenge P2"]
	if !ok {
		t.Fatal("P2 belief not found by claim")
	}
	if b2["origin_packet_id"] != p2PacketID {
		t.Errorf("P2 origin_packet_id = %v, want %v", b2["origin_packet_id"], p2PacketID)
	}
	// P2 submitted with no explicit debt → empty array (CompileDebt not wired)
	debt2, ok := b2["debt"].([]interface{})
	if !ok {
		t.Fatal("P2 debt not found or wrong type")
	}
	if len(debt2) != 0 {
		t.Errorf("P2 debt = %v, want [] (CompileDebt not in persist path)", debt2)
	}

	// Assert evidence (expect 2)
	evidence, ok := snapshot["evidence"].([]interface{})
	if !ok {
		t.Fatal("evidence not found or wrong type")
	}
	if len(evidence) != 2 {
		t.Fatalf("evidence count = %d, want 2", len(evidence))
	}

	// Assert edges (expect 1)
	edges, ok := snapshot["edges"].([]interface{})
	if !ok {
		t.Fatal("edges not found or wrong type")
	}
	if len(edges) != 1 {
		t.Fatalf("edge count = %d, want 1", len(edges))
	}
	edge, ok := edges[0].(map[string]interface{})
	if !ok {
		t.Fatal("edge entry wrong type")
	}
	if !strings.EqualFold(edge["parent_id"].(string), mcpNormalizeUUID(belief1ID)) {
		t.Errorf("edge parent_id = %v, want %v", edge["parent_id"], mcpNormalizeUUID(belief1ID))
	}
	if !strings.EqualFold(edge["child_id"].(string), mcpNormalizeUUID(belief2ID)) {
		t.Errorf("edge child_id = %v, want %v", edge["child_id"], mcpNormalizeUUID(belief2ID))
	}
	if edge["kind"] != "contradicts" {
		t.Errorf("edge kind = %v, want contradicts", edge["kind"])
	}

	t.Logf("MCP transport verified: beliefs=%d evidence=%d edges=%d origins={%s,%s}",
		len(beliefs), len(evidence), len(edges), p1PacketID, p2PacketID)

	// === Verify non-existent task returns availability false ===
	badArgs, _ := json.Marshal(map[string]string{"task_id": "00000000-0000-0000-0000-0000000000FF"})
	badResult, err := adapter.HandleTool(ctx, ToolGetContext, json.RawMessage(badArgs))
	if err != nil {
		t.Fatalf("HandleTool get_context for bad task: %v", err)
	}
	badRaw, err := json.Marshal(badResult)
	if err != nil {
		t.Fatalf("marshal bad result: %v", err)
	}
	var badDecoded map[string]interface{}
	if err := json.Unmarshal(badRaw, &badDecoded); err != nil {
		t.Fatalf("unmarshal bad result: %v", err)
	}
	badAvail, ok := badDecoded["availability"].(map[string]interface{})
	if !ok {
		t.Fatal("bad availability not found")
	}
	badTaskAvail, ok := badAvail["task"].(map[string]interface{})
	if !ok {
		t.Fatal("bad availability.task not found")
	}
	if badTaskAvail["available"] != false {
		t.Errorf("non-existent task: availability.task.available = %v, want false", badTaskAvail["available"])
	}
}
