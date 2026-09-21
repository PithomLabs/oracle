-- Packet-level agent provenance (append-only).
-- Records which agent/harness/model produced each submitted research packet.
-- Distinct from conductor_task.current_agent (task claiming).

CREATE TABLE IF NOT EXISTS packet_submission (
    packet_id      STRING PRIMARY KEY,
    scenario_id    UUID NOT NULL,
    task_id        UUID,
    agent_id       STRING NOT NULL,
    role           STRING NOT NULL,
    harness        STRING NOT NULL,
    model          STRING NOT NULL,
    content_sha256 STRING NOT NULL,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
