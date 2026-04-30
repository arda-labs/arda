ALTER TABLE fee_schedules
    ADD COLUMN IF NOT EXISTS approval_status TEXT NOT NULL DEFAULT 'DRAFT',
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS approved_by TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS change_note TEXT NOT NULL DEFAULT '';

ALTER TABLE tax_rules
    ADD COLUMN IF NOT EXISTS approval_status TEXT NOT NULL DEFAULT 'DRAFT',
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS approved_by TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS change_note TEXT NOT NULL DEFAULT '';

ALTER TABLE standard_limits
    ADD COLUMN IF NOT EXISTS approval_status TEXT NOT NULL DEFAULT 'DRAFT',
    ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS approved_by TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS change_note TEXT NOT NULL DEFAULT '';

UPDATE fee_schedules SET approval_status = 'APPROVED', approved_by = 'SYSTEM', approved_at = now()
WHERE status = 'ACTIVE' AND approval_status = 'DRAFT';

UPDATE tax_rules SET approval_status = 'APPROVED', approved_by = 'SYSTEM', approved_at = now()
WHERE status = 'ACTIVE' AND approval_status = 'DRAFT';

UPDATE standard_limits SET approval_status = 'APPROVED', approved_by = 'SYSTEM', approved_at = now()
WHERE status = 'ACTIVE' AND approval_status = 'DRAFT';

CREATE TABLE IF NOT EXISTS pricing_rule_audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rule_type TEXT NOT NULL,
    rule_id UUID NOT NULL,
    action TEXT NOT NULL,
    old_status TEXT NOT NULL DEFAULT '',
    new_status TEXT NOT NULL DEFAULT '',
    old_approval_status TEXT NOT NULL DEFAULT '',
    new_approval_status TEXT NOT NULL DEFAULT '',
    version INTEGER NOT NULL DEFAULT 1,
    actor TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS ix_pricing_rule_audit_logs_rule
    ON pricing_rule_audit_logs (rule_type, rule_id, created_at DESC);
