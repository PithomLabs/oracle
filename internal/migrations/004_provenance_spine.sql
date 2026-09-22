-- Provenance spine: link every packet-created object back to its origin packet.
-- origin_packet_id is nullable: seed/human-created objects get NULL.

ALTER TABLE belief ADD COLUMN IF NOT EXISTS origin_packet_id STRING NULL
    REFERENCES packet_submission(packet_id);

ALTER TABLE evidence ADD COLUMN IF NOT EXISTS origin_packet_id STRING NULL
    REFERENCES packet_submission(packet_id);

ALTER TABLE conductor_task ADD COLUMN IF NOT EXISTS origin_packet_id STRING NULL
    REFERENCES packet_submission(packet_id);

-- Edge provenance: separate table because belief_edge has no scenario_id
-- and a minimal schema. PK matches belief_edge identity.
CREATE TABLE IF NOT EXISTS edge_provenance (
    parent_id        UUID NOT NULL,
    child_id         UUID NOT NULL,
    origin_packet_id STRING NOT NULL
        REFERENCES packet_submission(packet_id),
    PRIMARY KEY (parent_id, child_id)
);
