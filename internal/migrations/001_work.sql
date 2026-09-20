-- ARGUS-owned work tables (CRDB-compatible).
-- Solvent-owned tables are applied via solventmigrations.Apply(db).

CREATE TABLE IF NOT EXISTS conductor_project (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        STRING NOT NULL,
    description STRING,
    status      STRING NOT NULL DEFAULT 'active'
                CHECK (status IN ('active','completed','archived')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conductor_task (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id            UUID NOT NULL REFERENCES conductor_project(id),
    title                 STRING NOT NULL,
    description           STRING,
    status                STRING NOT NULL DEFAULT 'proposed'
                          CHECK (status IN ('proposed','active','review','accepted','blocked','cancelled')),
    priority              STRING NOT NULL DEFAULT 'medium'
                          CHECK (priority IN ('low','medium','high','critical')),
    current_agent         STRING,
    governance_ref        UUID,
    reopened_from_task_id UUID REFERENCES conductor_task(id),
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conductor_dependency (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id      UUID NOT NULL REFERENCES conductor_task(id),
    blocked_by_id UUID NOT NULL REFERENCES conductor_task(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(task_id, blocked_by_id)
);

CREATE TABLE IF NOT EXISTS conductor_activity (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID NOT NULL REFERENCES conductor_task(id),
    actor_type STRING NOT NULL CHECK (actor_type IN ('human','agent','system')),
    actor_id   STRING NOT NULL,
    action     STRING NOT NULL,
    details    STRING,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS conductor_task_project ON conductor_task (project_id);
CREATE INDEX IF NOT EXISTS conductor_task_status ON conductor_task (status);
CREATE INDEX IF NOT EXISTS conductor_task_governance ON conductor_task (governance_ref);
CREATE INDEX IF NOT EXISTS conductor_activity_task ON conductor_activity (task_id);
