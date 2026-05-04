CREATE TABLE IF NOT EXISTS process_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zeebe_instance_key BIGINT NOT NULL,
    process_definition_id UUID NOT NULL REFERENCES process_definitions(id),
    status VARCHAR(50) NOT NULL DEFAULT 'RUNNING',
    current_step VARCHAR(255),
    variables JSONB DEFAULT '{}',
    assigned_agent VARCHAR(255),
    sla_status VARCHAR(50) DEFAULT 'ON_TIME',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_pi_status ON process_instances(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_pi_definition ON process_instances(process_definition_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pi_created ON process_instances(created_at DESC) WHERE deleted_at IS NULL;
