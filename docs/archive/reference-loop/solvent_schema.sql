-- Solvent SQLite Schema (adapted from db/001_schema.sql, 005_authority_mvp.sql, 007_service_tables.sql, 008_executing_state.sql, 009_exact_authority_binding.sql)
-- Using modernc.org/sqlite compatible syntax

-- Core belief tables (from 001_schema.sql)
CREATE TABLE IF NOT EXISTS belief (
    id          TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    scenario_id TEXT NOT NULL,
    claim       TEXT NOT NULL,
    claim_type  TEXT NOT NULL CHECK (claim_type IN ('derived','accommodated','postulated')),
    status      TEXT NOT NULL DEFAULT 'entered'
                CHECK (status IN ('entered','promoted','retracted')),
    debt        TEXT NOT NULL DEFAULT '["needProvenanceCheck","needContradictionSweep","needBlastRadius","needRollbackPlan","needVersionPin","needOperatorSignoff"]',
    final_truth INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT promoted_is_debt_free
        CHECK (status <> 'promoted'
               OR (coalesce(json_array_length(debt),0) = 0 AND final_truth = 0)),
    UNIQUE(id, status)
);

CREATE TABLE IF NOT EXISTS belief_edge (
    parent_id TEXT NOT NULL REFERENCES belief(id),
    child_id  TEXT NOT NULL REFERENCES belief(id),
    kind      TEXT NOT NULL DEFAULT 'derives' CHECK (kind IN ('derives','contradicts')),
    filed_at  TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (parent_id, child_id),
    CHECK (parent_id <> child_id)
);
CREATE INDEX IF NOT EXISTS belief_edge_child ON belief_edge (child_id);

CREATE TABLE IF NOT EXISTS evidence (
    id                 TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    scenario_id        TEXT NOT NULL,
    belief_id          TEXT NOT NULL REFERENCES belief(id),
    provenance_class   TEXT NOT NULL CHECK (provenance_class IN
        ('external_feed','reproducible_artifact','live_scan','operator_asserted')),
    source_url         TEXT,
    snapshot           TEXT,
    content_sha256     TEXT NOT NULL,
    source_observed_at TEXT,
    ingested_at        TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS action_intent (
    id            TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    scenario_id   TEXT NOT NULL,
    belief_id     TEXT NOT NULL,
    belief_status TEXT NOT NULL DEFAULT 'promoted',
    action        TEXT NOT NULL,
    state         TEXT NOT NULL DEFAULT 'live'
                  CHECK (state IN ('live','cancelled','executing','executed')),
    target_id     TEXT,
    snapshot_id   TEXT,
    CONSTRAINT live_requires_promoted CHECK (state <> 'live' OR belief_status = 'promoted'),
    CONSTRAINT intent_authority_binding_fk
        FOREIGN KEY (target_id, snapshot_id)
        REFERENCES target_snapshot(target_id, snapshot_id)
        ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS live_intents ON action_intent (belief_id) WHERE state = 'live';
CREATE INDEX IF NOT EXISTS executing_intents ON action_intent (scenario_id) WHERE state = 'executing';
CREATE UNIQUE INDEX IF NOT EXISTS live_intent_per_snapshot
    ON action_intent (target_id, snapshot_id)
    WHERE state = 'live' AND target_id IS NOT NULL AND snapshot_id IS NOT NULL;

-- Authority tables (from 005_authority_mvp.sql)
-- Note: SQLite doesn't support gen_random_uuid(), so we use DEFAULT (lower(hex(randomblob(16))))
CREATE TABLE IF NOT EXISTS principal (
    principal_id   TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    principal_type TEXT NOT NULL CHECK (principal_type IN ('human','agent','workload','service')),
    issuer         TEXT NOT NULL,
    revoked_at     TEXT,
    created_at     TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS authority_target (
    target_id               TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    principal_id            TEXT NOT NULL REFERENCES principal(principal_id),
    resource_type           TEXT NOT NULL CHECK (resource_type <> ''),
    resource_id             TEXT NOT NULL CHECK (resource_id <> ''),
    scope                   TEXT NOT NULL CHECK (scope <> ''),
    action_namespace        TEXT NOT NULL CHECK (action_namespace <> ''),
    action_name             TEXT NOT NULL CHECK (action_name <> ''),
    consequence_type        TEXT NOT NULL CHECK (consequence_type <> ''),
    consequence_parameters  TEXT NOT NULL,
    created_by              TEXT NOT NULL REFERENCES principal(principal_id),
    created_at              TEXT NOT NULL DEFAULT (datetime('now')),
    requested_by            TEXT REFERENCES principal(principal_id),
    requested_at            TEXT,
    pinned_request_hash     TEXT
);

CREATE TABLE IF NOT EXISTS target_snapshot (
    snapshot_id            TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    target_id              TEXT NOT NULL REFERENCES authority_target(target_id),
    principal_id           TEXT NOT NULL REFERENCES principal(principal_id),
    resource_type          TEXT NOT NULL,
    resource_id            TEXT NOT NULL,
    scope                  TEXT NOT NULL,
    action_namespace       TEXT NOT NULL,
    action_name            TEXT NOT NULL,
    consequence_type       TEXT NOT NULL,
    consequence_parameters TEXT NOT NULL,
    justification_set      TEXT NOT NULL,
    approver_principal_id  TEXT NOT NULL REFERENCES principal(principal_id),
    approved_at            TEXT NOT NULL DEFAULT (datetime('now')),
    snapshot_hash          TEXT NOT NULL,
    UNIQUE(target_id, snapshot_id)
);

CREATE TABLE IF NOT EXISTS target_activation (
    activation_id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    target_id     TEXT NOT NULL REFERENCES authority_target(target_id),
    snapshot_id   TEXT NOT NULL,
    activated_at  TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(target_id),
    FOREIGN KEY (target_id, snapshot_id)
        REFERENCES target_snapshot(target_id, snapshot_id)
);

CREATE TABLE IF NOT EXISTS target_revocation (
    target_id   TEXT PRIMARY KEY REFERENCES authority_target(target_id),
    revoked_at  TEXT NOT NULL DEFAULT (datetime('now')),
    revoked_by  TEXT NOT NULL REFERENCES principal(principal_id),
    reason      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS justification (
    justification_id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    target_id        TEXT NOT NULL REFERENCES authority_target(target_id),
    belief_id        TEXT NOT NULL,
    belief_status    TEXT NOT NULL,
    attached_by      TEXT NOT NULL REFERENCES principal(principal_id),
    attached_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(target_id, belief_id, belief_status),
    FOREIGN KEY (belief_id, belief_status)
        REFERENCES belief(id, status)
);

CREATE TABLE IF NOT EXISTS debt_discharge (
    discharge_id   TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    belief_id      TEXT NOT NULL REFERENCES belief(id),
    obligation_key TEXT NOT NULL,
    instrument_ref TEXT NOT NULL,
    discharged_by  TEXT NOT NULL REFERENCES principal(principal_id),
    accepted_at    TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(belief_id, obligation_key, instrument_ref)
);

-- Service tables (from 007_service_tables.sql)
CREATE TABLE IF NOT EXISTS workflow_token (
    id            TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    scenario_id   TEXT NOT NULL,
    belief_id     TEXT NOT NULL,
    action_type   TEXT NOT NULL,
    state         TEXT NOT NULL DEFAULT 'pending'
                  CHECK (state IN ('pending','prepared','executing','completed','failed','expired')),
    payload       TEXT NOT NULL DEFAULT '{}',
    created_at    TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at    TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at    TEXT,
    completed_at  TEXT,
    failure_error TEXT
);
CREATE INDEX IF NOT EXISTS workflow_token_scenario ON workflow_token (scenario_id);
CREATE INDEX IF NOT EXISTS workflow_token_state ON workflow_token (state) WHERE state IN ('pending','prepared');

CREATE TABLE IF NOT EXISTS policy_tool (
    id               TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name             TEXT NOT NULL UNIQUE,
    class            TEXT NOT NULL CHECK (class IN ('read','mutate','orchestrate')),
    required_beliefs TEXT NOT NULL DEFAULT '[]',
    description      TEXT NOT NULL DEFAULT '',
    created_at       TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS policy_actor (
    id         TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    name       TEXT NOT NULL UNIQUE,
    roles      TEXT NOT NULL DEFAULT '[]',
    max_class  TEXT NOT NULL CHECK (max_class IN ('read','mutate','orchestrate')),
    active     INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS audit_activity (
    id              TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    scenario_id     TEXT NOT NULL,
    type            TEXT NOT NULL,
    actor_id        TEXT,
    subject_id      TEXT,
    details         TEXT,
    sqlstate        TEXT,
    constraint_name TEXT,
    refusal         INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS audit_activity_scenario ON audit_activity (scenario_id);
CREATE INDEX IF NOT EXISTS audit_activity_refusal ON audit_activity (scenario_id) WHERE refusal = 1;

-- Authority indexes
CREATE INDEX IF NOT EXISTS authority_target_principal ON authority_target (principal_id);
CREATE INDEX IF NOT EXISTS justification_target ON justification (target_id);
CREATE INDEX IF NOT EXISTS justification_belief ON justification (belief_id);
CREATE INDEX IF NOT EXISTS debt_discharge_belief ON debt_discharge (belief_id);
CREATE INDEX IF NOT EXISTS debt_discharge_discharged_by ON debt_discharge (discharged_by);
