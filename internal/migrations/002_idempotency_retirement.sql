-- ARGUS-owned idempotency and proposed retirement tables.

CREATE TABLE IF NOT EXISTS submission_idempotency (
    content_hash  STRING NOT NULL,
    scenario_id   STRING NOT NULL,
    packet_id     STRING NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(content_hash, scenario_id)
);

CREATE TABLE IF NOT EXISTS belief_retirement_proposal (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    belief_id       UUID NOT NULL,
    scenario_id     STRING NOT NULL,
    debt_items      JSONB NOT NULL,
    proposer        STRING NOT NULL,
    proposed_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(belief_id, scenario_id)
);
