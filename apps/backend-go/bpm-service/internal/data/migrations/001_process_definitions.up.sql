CREATE TABLE IF NOT EXISTS process_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    process_key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    module VARCHAR(50) NOT NULL DEFAULT 'GENERAL',
    bpmn_xml TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    zeebe_deployment_key BIGINT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_process_key_version UNIQUE (process_key, version)
);

CREATE INDEX idx_pd_active ON process_definitions(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_pd_module ON process_definitions(module) WHERE deleted_at IS NULL;
