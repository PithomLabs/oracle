package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/PithomLabs/oracle/internal/application"
)

// Fixed IDs for the BM-IST-AS POC seed.
const (
	ProjectID = "00000000-0000-0000-0000-000000000100"
	TaskID    = "00000000-0000-0000-0000-000000000101"
	BeliefID  = "00000000-0000-0000-0000-000000000102"
)

var debtItems = []string{
	"needMap",
	"needInvariant",
	"needToyCheck",
	"needNullModel",
	"needObstruction",
	"needFaithfulnessReview",
	"needInitialCondition",
	"needRegularity",
}

// SeedIfEmpty inserts the BM-IST-AS POC seed if the DB has no beliefs
// for this scenario. It is idempotent: beliefs in other scenarios do not
// prevent creation.
func SeedIfEmpty(db *sql.DB) error {
	ctx := context.Background()

	var count int
	if err := db.QueryRowContext(ctx,
		"SELECT count(*) FROM belief WHERE scenario_id = $1", ProjectID).Scan(&count); err != nil {
		return fmt.Errorf("seed: check belief count: %w", err)
	}
	if count > 0 {
		log.Printf("seed: %d beliefs exist for scenario %s, skipping seed", count, ProjectID)
		return nil
	}

	log.Println("seed: inserting BM-IST-AS POC seed...")

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("seed: begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Create project.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO conductor_project (id, name, status) VALUES ($1, $2, 'active')
		 ON CONFLICT (id) DO NOTHING`,
		ProjectID, "BM-IST-AS")
	if err != nil {
		return fmt.Errorf("seed: insert project: %w", err)
	}

	// 2. Create belief with debt.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO belief (id, scenario_id, claim, claim_type, debt, status)
		 VALUES ($1, $2, $3, 'derived', $4, 'entered')
		 ON CONFLICT (id) DO NOTHING`,
		BeliefID, ProjectID,
		"G0 — countable substrate + literal continuous unitary evolution → structural incompatibility",
		debtItems)
	if err != nil {
		return fmt.Errorf("seed: insert belief: %w", err)
	}

	// 3. Create task with governance_ref pointing to the belief.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO conductor_task (id, project_id, title, description, status, priority, governance_ref, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, 'proposed', 'high', $5, now(), now())
		 ON CONFLICT (id) DO NOTHING`,
		TaskID, ProjectID, "Gate G0", "Verify G0 structural incompatibility claim", BeliefID)
	if err != nil {
		return fmt.Errorf("seed: insert task: %w", err)
	}

	// 4. Create operator principal for human decisions.
	_, err = tx.ExecContext(ctx,
		`INSERT INTO principal (principal_id, principal_type, issuer)
		 VALUES ($1, 'human', 'operator')
		 ON CONFLICT (principal_id) DO NOTHING`,
		application.LocalOperatorPrincipalID)
	if err != nil {
		return fmt.Errorf("seed: insert operator principal: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seed: commit: %w", err)
	}

	log.Printf("seed: BM-IST-AS POC seed inserted (project=%s, task=%s, belief=%s)", ProjectID, TaskID, BeliefID)
	return nil
}
