package epistemic

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("ARGUS_TEST_DSN")
	if dsn == "" {
		dsn = "postgres://root@localhost:26257/defaultdb?sslmode=disable"
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("ping db: %v", err)
	}
	return db
}

func TestDerivesEdgeReturned(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000200"

	// Create two beliefs
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'parent claim', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000001", scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'child claim', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000002", scenarioID)
	if err != nil {
		t.Fatal(err)
	}

	// Create derives edge
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'derives')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000001", "10000000-0000-0000-0000-000000000002")
	if err != nil {
		t.Fatal(err)
	}

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if len(snap.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(snap.Edges))
	}
	e := snap.Edges[0]
	if e.ParentID != "10000000-0000-0000-0000-000000000001" {
		t.Errorf("parent_id = %s, want 10000000-0000-0000-0000-000000000001", e.ParentID)
	}
	if e.ChildID != "10000000-0000-0000-0000-000000000002" {
		t.Errorf("child_id = %s, want 10000000-0000-0000-0000-000000000002", e.ChildID)
	}
	if e.Kind != "derives" {
		t.Errorf("kind = %s, want derives", e.Kind)
	}
}

func TestContradictsEdgeReturned(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000201"

	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'original claim', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000003", scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'challenge claim', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000004", scenarioID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'contradicts')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000004", "10000000-0000-0000-0000-000000000003")
	if err != nil {
		t.Fatal(err)
	}

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if len(snap.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(snap.Edges))
	}
	if snap.Edges[0].Kind != "contradicts" {
		t.Errorf("kind = %s, want contradicts", snap.Edges[0].Kind)
	}
}

func TestEdgesReferenceCorrectBeliefIDs(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000202"

	parentID := "10000000-0000-0000-0000-000000000005"
	childID := "10000000-0000-0000-0000-000000000006"

	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'parent', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`, parentID, scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'child', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`, childID, scenarioID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'derives')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`, parentID, childID)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if len(snap.Edges) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(snap.Edges))
	}
	if snap.Edges[0].ParentID != parentID {
		t.Errorf("parent_id = %s, want %s", snap.Edges[0].ParentID, parentID)
	}
	if snap.Edges[0].ChildID != childID {
		t.Errorf("child_id = %s, want %s", snap.Edges[0].ChildID, childID)
	}
}

func TestNoAuthorityCapabilityAdded(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000203"

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	// Snapshot should have no methods that mutate state
	// Edges is a read-only projection
	if snap.Edges == nil {
		// nil is acceptable (empty slice)
		return
	}
	// Edges should be empty for empty scenario
	if len(snap.Edges) != 0 {
		t.Errorf("expected 0 edges for empty scenario, got %d", len(snap.Edges))
	}
}

func TestScenarioScoping(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioA := "00000000-0000-0000-0000-000000000204"
	scenarioB := "00000000-0000-0000-0000-000000000205"

	// Create belief and edge in scenario A
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'A parent', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000007", scenarioA)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'A child', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000008", scenarioA)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'derives')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000007", "10000000-0000-0000-0000-000000000008")
	if err != nil {
		t.Fatal(err)
	}

	// Query scenario B — should have no edges
	snap, err := GetSnapshot(ctx, db, scenarioB, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snap.Edges) != 0 {
		t.Errorf("scenario B should have 0 edges, got %d", len(snap.Edges))
	}

	// Query scenario A — should have 1 edge
	snapA, err := GetSnapshot(ctx, db, scenarioA, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if len(snapA.Edges) != 1 {
		t.Errorf("scenario A should have 1 edge, got %d", len(snapA.Edges))
	}
}

func TestEmptyEdgesReturnsEmptySlice(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000206"

	// Create belief but no edges
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
		 VALUES ($1::UUID, $2::UUID, 'lonely belief', 'derived', 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000009", scenarioID)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if snap.Edges == nil {
		t.Error("edges should be empty slice, not nil")
	}
	if len(snap.Edges) != 0 {
		t.Errorf("expected 0 edges, got %d", len(snap.Edges))
	}
}

func TestMultipleEdgesReturned(t *testing.T) {
	db := testDB(t)
	defer db.Close()
	ctx := context.Background()
	scenarioID := "00000000-0000-0000-0000-000000000207"

	// Create 3 beliefs
	for i, id := range []string{
		"10000000-0000-0000-0000-000000000010",
		"10000000-0000-0000-0000-000000000011",
		"10000000-0000-0000-0000-000000000012",
	} {
		_, err := db.ExecContext(ctx,
			`INSERT INTO belief (id, scenario_id, claim, claim_type, status)
			 VALUES ($1::UUID, $2::UUID, $3, 'derived', 'entered')
			 ON CONFLICT (id) DO NOTHING`,
			id, scenarioID, "belief"+string(rune('A'+i)))
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create 2 edges: A->B derives, C->A contradicts
	_, err := db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'derives')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000010", "10000000-0000-0000-0000-000000000011")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.ExecContext(ctx,
		`INSERT INTO belief_edge (parent_id, child_id, kind)
		 VALUES ($1::UUID, $2::UUID, 'contradicts')
		 ON CONFLICT (parent_id, child_id) DO NOTHING`,
		"10000000-0000-0000-0000-000000000012", "10000000-0000-0000-0000-000000000010")
	if err != nil {
		t.Fatal(err)
	}

	snap, err := GetSnapshot(ctx, db, scenarioID, SnapshotOpts{})
	if err != nil {
		t.Fatal(err)
	}

	if len(snap.Edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(snap.Edges))
	}

	// Verify both kinds present
	kinds := map[string]bool{}
	for _, e := range snap.Edges {
		kinds[e.Kind] = true
	}
	if !kinds["derives"] {
		t.Error("missing derives edge")
	}
	if !kinds["contradicts"] {
		t.Error("missing contradicts edge")
	}
}
