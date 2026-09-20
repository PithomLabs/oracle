package solventmigrations

import (
	"context"
	"database/sql"
	"fmt"
)

// Apply runs the core Solvent epistemic schema.
// This is a minimal POC subset of the full Solvent migration set.
func Apply(ctx context.Context, db *sql.DB) error {
	for _, stmt := range schema {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("solvent migration: %w", err)
		}
	}
	return nil
}

var schema = []string{
	`CREATE TABLE IF NOT EXISTS belief (
		id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		scenario_id UUID NOT NULL,
		claim       STRING NOT NULL,
		claim_type  STRING NOT NULL DEFAULT 'derived',
		status      STRING NOT NULL DEFAULT 'entered'
		          CHECK (status IN ('entered','promoted','retracted')),
		debt        STRING[] NOT NULL DEFAULT '{}',
		final_truth BOOL NOT NULL DEFAULT FALSE,
		created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
		CONSTRAINT promoted_is_debt_free CHECK (status <> 'promoted' OR (coalesce(array_length(debt,1),0) = 0 AND NOT final_truth))
	)`,
	`CREATE TABLE IF NOT EXISTS principal (
		principal_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		principal_type TEXT NOT NULL CHECK (principal_type IN ('human','agent','workload','service')),
		issuer         TEXT NOT NULL,
		revoked_at     TIMESTAMPTZ NULL,
		created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE TABLE IF NOT EXISTS debt_discharge (
		discharge_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		belief_id      UUID NOT NULL REFERENCES belief(id),
		obligation_key TEXT NOT NULL,
		instrument_ref TEXT NOT NULL,
		discharged_by  UUID NOT NULL REFERENCES principal(principal_id),
		accepted_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
		UNIQUE(belief_id, obligation_key, instrument_ref)
	)`,
	`CREATE TABLE IF NOT EXISTS belief_edge (
		parent_id UUID NOT NULL REFERENCES belief(id),
		child_id  UUID NOT NULL REFERENCES belief(id),
		kind      STRING NOT NULL CHECK (kind IN ('derives','contradicts')),
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		PRIMARY KEY (parent_id, child_id)
	)`,
	`CREATE TABLE IF NOT EXISTS evidence (
		id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		scenario_id      UUID NOT NULL,
		belief_id        UUID NOT NULL REFERENCES belief(id),
		provenance_class STRING NOT NULL,
		source_url       STRING,
		content_sha256   STRING NOT NULL,
		ingested_at      TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE TABLE IF NOT EXISTS action_intent (
		id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		scenario_id UUID NOT NULL,
		belief_id  UUID NOT NULL REFERENCES belief(id),
		action     STRING NOT NULL,
		state      STRING NOT NULL DEFAULT 'live'
		          CHECK (state IN ('live','executed','revoked')),
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`,
	`CREATE INDEX IF NOT EXISTS belief_edge_child ON belief_edge (child_id)`,
	`CREATE INDEX IF NOT EXISTS live_intents ON action_intent (belief_id) WHERE state = 'live'`,
	`CREATE INDEX IF NOT EXISTS debt_discharge_belief ON debt_discharge (belief_id)`,
	`CREATE INDEX IF NOT EXISTS debt_discharge_discharged_by ON debt_discharge (discharged_by)`,
}
