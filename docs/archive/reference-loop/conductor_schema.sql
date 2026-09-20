-- Conductor SQLite Schema (adapted from migrations/001_initial.sql and 002_dependencies.sql)

CREATE TABLE IF NOT EXISTS conductor_project (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    status      TEXT NOT NULL DEFAULT 'active'
                CHECK (status IN ('active', 'completed', 'archived')),
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS conductor_task (
    id              TEXT PRIMARY KEY,
    project_id      TEXT NOT NULL REFERENCES conductor_project(id),
    title           TEXT NOT NULL,
    description     TEXT,
    status          TEXT NOT NULL DEFAULT 'proposed'
                    CHECK (status IN ('proposed', 'active', 'review', 'accepted', 'blocked', 'cancelled')),
    priority        TEXT NOT NULL DEFAULT 'medium'
                    CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    current_agent   TEXT,
    governance_ref  TEXT,
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_task_project ON conductor_task(project_id);
CREATE INDEX IF NOT EXISTS idx_task_status ON conductor_task(status);
CREATE INDEX IF NOT EXISTS idx_task_agent ON conductor_task(current_agent);

CREATE TABLE IF NOT EXISTS conductor_activity (
    id          TEXT PRIMARY KEY,
    task_id     TEXT NOT NULL REFERENCES conductor_task(id),
    actor_type  TEXT NOT NULL CHECK (actor_type IN ('human', 'agent', 'system')),
    actor_id    TEXT NOT NULL,
    action      TEXT NOT NULL,
    details     TEXT,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_activity_task ON conductor_activity(task_id);
CREATE INDEX IF NOT EXISTS idx_activity_created ON conductor_activity(task_id, created_at DESC);

CREATE TABLE IF NOT EXISTS conductor_dependency (
    id              TEXT PRIMARY KEY,
    task_id         TEXT NOT NULL REFERENCES conductor_task(id),
    blocked_by_id   TEXT NOT NULL REFERENCES conductor_task(id),
    created_at      TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(task_id, blocked_by_id),
    CHECK (task_id != blocked_by_id)
);

CREATE INDEX IF NOT EXISTS idx_dependency_task ON conductor_dependency(task_id);
CREATE INDEX IF NOT EXISTS idx_dependency_blocked_by ON conductor_dependency(blocked_by_id);
