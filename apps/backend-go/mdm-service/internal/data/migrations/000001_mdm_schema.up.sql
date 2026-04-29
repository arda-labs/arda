CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS administrative_units (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code           TEXT NOT NULL,
    name           TEXT NOT NULL,
    full_name      TEXT NOT NULL DEFAULT '',
    short_name     TEXT NOT NULL DEFAULT '',
    level          TEXT NOT NULL,
    unit_type      TEXT NOT NULL,
    parent_id      UUID REFERENCES administrative_units(id),
    path           TEXT NOT NULL DEFAULT '',
    sort_order     INTEGER NOT NULL DEFAULT 0,
    latitude       DOUBLE PRECISION,
    longitude      DOUBLE PRECISION,
    status         TEXT NOT NULL DEFAULT 'ACTIVE',
    effective_from DATE,
    effective_to   DATE,
    source         TEXT NOT NULL DEFAULT '',
    metadata       JSONB NOT NULL DEFAULT '{}',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_administrative_units_code_active
    ON administrative_units (code)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_administrative_units_parent
    ON administrative_units (parent_id)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_administrative_units_level
    ON administrative_units (level, status)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS administrative_unit_mappings (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    old_unit_id    UUID REFERENCES administrative_units(id),
    new_unit_id    UUID REFERENCES administrative_units(id),
    mapping_type   TEXT NOT NULL,
    effective_date DATE,
    note           TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS area_types (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    allow_hierarchy BOOLEAN NOT NULL DEFAULT true,
    status          TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_area_types_code_active
    ON area_types (code)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS areas (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    area_type_id    UUID NOT NULL REFERENCES area_types(id),
    parent_id       UUID REFERENCES areas(id),
    code            TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    manager_user_id TEXT NOT NULL DEFAULT '',
    status          TEXT NOT NULL DEFAULT 'ACTIVE',
    effective_from  DATE,
    effective_to    DATE,
    metadata        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_areas_type_code_active
    ON areas (area_type_id, code)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_areas_parent
    ON areas (parent_id)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS area_administrative_units (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    area_id                UUID NOT NULL REFERENCES areas(id) ON DELETE CASCADE,
    administrative_unit_id UUID NOT NULL REFERENCES administrative_units(id),
    scope_type             TEXT NOT NULL DEFAULT 'INCLUDE',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (area_id, administrative_unit_id)
);

CREATE TABLE IF NOT EXISTS code_sets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code        TEXT NOT NULL,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    is_system   BOOLEAN NOT NULL DEFAULT false,
    status      TEXT NOT NULL DEFAULT 'ACTIVE',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_code_sets_code_active
    ON code_sets (code)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS code_items (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code_set_id    UUID NOT NULL REFERENCES code_sets(id),
    code           TEXT NOT NULL,
    name           TEXT NOT NULL,
    value          TEXT NOT NULL DEFAULT '',
    parent_id      UUID REFERENCES code_items(id),
    sort_order     INTEGER NOT NULL DEFAULT 0,
    color          TEXT NOT NULL DEFAULT '',
    icon           TEXT NOT NULL DEFAULT '',
    metadata       JSONB NOT NULL DEFAULT '{}',
    is_default     BOOLEAN NOT NULL DEFAULT false,
    is_system      BOOLEAN NOT NULL DEFAULT false,
    status         TEXT NOT NULL DEFAULT 'ACTIVE',
    effective_from DATE,
    effective_to   DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at     TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_code_items_set_code_active
    ON code_items (code_set_id, code)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_code_items_set
    ON code_items (code_set_id, sort_order)
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS system_parameters (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key                  TEXT NOT NULL,
    name                 TEXT NOT NULL,
    group_code           TEXT NOT NULL DEFAULT '',
    value_type           TEXT NOT NULL DEFAULT 'STRING',
    value_text           TEXT NOT NULL DEFAULT '',
    value_number         DOUBLE PRECISION,
    value_boolean        BOOLEAN,
    value_json           JSONB NOT NULL DEFAULT '{}',
    default_value        TEXT NOT NULL DEFAULT '',
    is_secret            BOOLEAN NOT NULL DEFAULT false,
    is_editable          BOOLEAN NOT NULL DEFAULT true,
    is_system            BOOLEAN NOT NULL DEFAULT false,
    validation_rule      JSONB NOT NULL DEFAULT '{}',
    description          TEXT NOT NULL DEFAULT '',
    status               TEXT NOT NULL DEFAULT 'ACTIVE',
    updated_by           TEXT NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at           TIMESTAMPTZ
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_system_parameters_key_active
    ON system_parameters (key)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_system_parameters_group
    ON system_parameters (group_code, status)
    WHERE deleted_at IS NULL;
