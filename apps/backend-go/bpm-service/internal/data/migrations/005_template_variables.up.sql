CREATE TABLE IF NOT EXISTS template_variables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL REFERENCES process_templates(id) ON DELETE CASCADE,
    variable_name VARCHAR(100) NOT NULL,
    source_type VARCHAR(50) NOT NULL,
    source_field VARCHAR(255),
    resolver_config JSONB,
    fallback_value TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tv_name ON template_variables(template_id, variable_name);
