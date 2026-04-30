DROP TABLE IF EXISTS pricing_rule_audit_logs;

ALTER TABLE standard_limits
    DROP COLUMN IF EXISTS change_note,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS approval_status;

ALTER TABLE tax_rules
    DROP COLUMN IF EXISTS change_note,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS approval_status;

ALTER TABLE fee_schedules
    DROP COLUMN IF EXISTS change_note,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS version,
    DROP COLUMN IF EXISTS approval_status;
