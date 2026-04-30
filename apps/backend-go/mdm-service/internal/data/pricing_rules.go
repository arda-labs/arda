package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListFeeSchedules(ctx context.Context, filter biz.PageFilter) ([]*biz.FeeSchedule, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, fee_type, calculation_method, currency, fixed_amount::float8,
		       rate_percent::float8, min_amount::float8, max_amount::float8, channel, product_code,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM fee_schedules
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%'
		       OR fee_type ILIKE '%' || $2 || '%' OR channel ILIKE '%' || $2 || '%' OR product_code ILIKE '%' || $2 || '%')
		ORDER BY code ASC, id ASC
		LIMIT $3 OFFSET $4`,
		filter.Status, filter.Keyword, page.Limit+1, page.Offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanFeeSchedules(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetFeeSchedule(ctx context.Context, id string) (*biz.FeeSchedule, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, fee_type, calculation_method, currency, fixed_amount::float8,
		       rate_percent::float8, min_amount::float8, max_amount::float8, channel, product_code,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM fee_schedules
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanFeeSchedule(row)
}

func (r *MdmRepo) CreateFeeSchedule(ctx context.Context, item *biz.FeeSchedule) (*biz.FeeSchedule, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO fee_schedules (
			code, name, fee_type, calculation_method, currency, fixed_amount, rate_percent,
			min_amount, max_amount, channel, product_code, effective_from, effective_to, description, status,
			approval_status, version, approved_by, change_note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NULLIF($12, '')::date, NULLIF($13, '')::date, $14, $15, $16, $17, $18, $19)
		RETURNING id::text`,
		item.Code, item.Name, item.FeeType, item.CalculationMethod, item.Currency, item.FixedAmount,
		item.RatePercent, item.MinAmount, item.MaxAmount, item.Channel, item.ProductCode,
		item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetFeeSchedule(ctx, item.ID)
}

func (r *MdmRepo) UpdateFeeSchedule(ctx context.Context, item *biz.FeeSchedule) (*biz.FeeSchedule, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE fee_schedules
		SET code = $2, name = $3, fee_type = $4, calculation_method = $5, currency = $6,
		    fixed_amount = $7, rate_percent = $8, min_amount = $9, max_amount = $10,
		    channel = $11, product_code = $12, effective_from = NULLIF($13, '')::date,
		    effective_to = NULLIF($14, '')::date, description = $15, status = $16,
		    approval_status = $17, version = $18, approved_by = $19, change_note = $20, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.FeeType, item.CalculationMethod, item.Currency,
		item.FixedAmount, item.RatePercent, item.MinAmount, item.MaxAmount, item.Channel,
		item.ProductCode, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetFeeSchedule(ctx, item.ID)
}

func (r *MdmRepo) DeleteFeeSchedule(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "fee_schedules", id)
}

func (r *MdmRepo) ApproveFeeSchedule(ctx context.Context, id, actor, note string) (*biz.FeeSchedule, error) {
	item, err := r.GetFeeSchedule(ctx, id)
	if err != nil {
		return nil, err
	}
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE fee_schedules
		SET approval_status = 'APPROVED', status = 'ACTIVE', approved_by = $2,
		    approved_at = now(), change_note = $3, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`, id, actor, note)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	if _, err := r.data.db.Pool.Exec(ctx, `
		INSERT INTO pricing_rule_audit_logs (
			rule_type, rule_id, action, old_status, new_status, old_approval_status,
			new_approval_status, version, actor, note
		)
		VALUES ('FEE_SCHEDULE', $1, 'APPROVE', $2, 'ACTIVE', $3, 'APPROVED', $4, $5, $6)`,
		id, item.Status, item.ApprovalStatus, item.Version, actor, note,
	); err != nil {
		return nil, err
	}
	return r.GetFeeSchedule(ctx, id)
}

func (r *MdmRepo) ListTaxRules(ctx context.Context, filter biz.PageFilter) ([]*biz.TaxRule, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, tax_type, rate_percent::float8, inclusive, jurisdiction,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM tax_rules
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%'
		       OR tax_type ILIKE '%' || $2 || '%' OR jurisdiction ILIKE '%' || $2 || '%')
		ORDER BY code ASC, id ASC
		LIMIT $3 OFFSET $4`,
		filter.Status, filter.Keyword, page.Limit+1, page.Offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanTaxRules(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetTaxRule(ctx context.Context, id string) (*biz.TaxRule, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, tax_type, rate_percent::float8, inclusive, jurisdiction,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM tax_rules
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanTaxRule(row)
}

func (r *MdmRepo) CreateTaxRule(ctx context.Context, item *biz.TaxRule) (*biz.TaxRule, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO tax_rules (
			code, name, tax_type, rate_percent, inclusive, jurisdiction,
			effective_from, effective_to, description, status, approval_status, version, approved_by, change_note
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::date, NULLIF($8, '')::date, $9, $10, $11, $12, $13, $14)
		RETURNING id::text`,
		item.Code, item.Name, item.TaxType, item.RatePercent, item.Inclusive, item.Jurisdiction,
		item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetTaxRule(ctx, item.ID)
}

func (r *MdmRepo) UpdateTaxRule(ctx context.Context, item *biz.TaxRule) (*biz.TaxRule, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE tax_rules
		SET code = $2, name = $3, tax_type = $4, rate_percent = $5, inclusive = $6,
		    jurisdiction = $7, effective_from = NULLIF($8, '')::date,
		    effective_to = NULLIF($9, '')::date, description = $10, status = $11,
		    approval_status = $12, version = $13, approved_by = $14, change_note = $15, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.TaxType, item.RatePercent, item.Inclusive,
		item.Jurisdiction, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetTaxRule(ctx, item.ID)
}

func (r *MdmRepo) DeleteTaxRule(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "tax_rules", id)
}

func (r *MdmRepo) ApproveTaxRule(ctx context.Context, id, actor, note string) (*biz.TaxRule, error) {
	item, err := r.GetTaxRule(ctx, id)
	if err != nil {
		return nil, err
	}
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE tax_rules
		SET approval_status = 'APPROVED', status = 'ACTIVE', approved_by = $2,
		    approved_at = now(), change_note = $3, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`, id, actor, note)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	if _, err := r.data.db.Pool.Exec(ctx, `
		INSERT INTO pricing_rule_audit_logs (
			rule_type, rule_id, action, old_status, new_status, old_approval_status,
			new_approval_status, version, actor, note
		)
		VALUES ('TAX_RULE', $1, 'APPROVE', $2, 'ACTIVE', $3, 'APPROVED', $4, $5, $6)`,
		id, item.Status, item.ApprovalStatus, item.Version, actor, note,
	); err != nil {
		return nil, err
	}
	return r.GetTaxRule(ctx, id)
}

func (r *MdmRepo) ListStandardLimits(ctx context.Context, filter biz.PageFilter) ([]*biz.StandardLimit, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, limit_type, subject_type, currency, min_amount::float8,
		       per_txn_amount::float8, daily_amount::float8, monthly_amount::float8, count_limit,
		       channel, product_code, COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''), description, status,
		       approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM standard_limits
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%'
		       OR limit_type ILIKE '%' || $2 || '%' OR subject_type ILIKE '%' || $2 || '%'
		       OR channel ILIKE '%' || $2 || '%' OR product_code ILIKE '%' || $2 || '%')
		ORDER BY code ASC, id ASC
		LIMIT $3 OFFSET $4`,
		filter.Status, filter.Keyword, page.Limit+1, page.Offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanStandardLimits(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetStandardLimit(ctx context.Context, id string) (*biz.StandardLimit, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, limit_type, subject_type, currency, min_amount::float8,
		       per_txn_amount::float8, daily_amount::float8, monthly_amount::float8, count_limit,
		       channel, product_code, COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''), description, status,
		       approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, created_at, updated_at
		FROM standard_limits
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanStandardLimit(row)
}

func (r *MdmRepo) CreateStandardLimit(ctx context.Context, item *biz.StandardLimit) (*biz.StandardLimit, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO standard_limits (
			code, name, limit_type, subject_type, currency, min_amount, per_txn_amount,
			daily_amount, monthly_amount, count_limit, channel, product_code,
			effective_from, effective_to, description, status, approval_status, version, approved_by, change_note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NULLIF($13, '')::date, NULLIF($14, '')::date, $15, $16, $17, $18, $19, $20)
		RETURNING id::text`,
		item.Code, item.Name, item.LimitType, item.SubjectType, item.Currency, item.MinAmount,
		item.PerTxnAmount, item.DailyAmount, item.MonthlyAmount, item.CountLimit, item.Channel,
		item.ProductCode, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetStandardLimit(ctx, item.ID)
}

func (r *MdmRepo) UpdateStandardLimit(ctx context.Context, item *biz.StandardLimit) (*biz.StandardLimit, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE standard_limits
		SET code = $2, name = $3, limit_type = $4, subject_type = $5, currency = $6,
		    min_amount = $7, per_txn_amount = $8, daily_amount = $9, monthly_amount = $10,
		    count_limit = $11, channel = $12, product_code = $13,
		    effective_from = NULLIF($14, '')::date, effective_to = NULLIF($15, '')::date,
		    description = $16, status = $17, approval_status = $18, version = $19,
		    approved_by = $20, change_note = $21, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.LimitType, item.SubjectType, item.Currency,
		item.MinAmount, item.PerTxnAmount, item.DailyAmount, item.MonthlyAmount, item.CountLimit,
		item.Channel, item.ProductCode, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
		item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetStandardLimit(ctx, item.ID)
}

func (r *MdmRepo) DeleteStandardLimit(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "standard_limits", id)
}

func (r *MdmRepo) ApproveStandardLimit(ctx context.Context, id, actor, note string) (*biz.StandardLimit, error) {
	item, err := r.GetStandardLimit(ctx, id)
	if err != nil {
		return nil, err
	}
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE standard_limits
		SET approval_status = 'APPROVED', status = 'ACTIVE', approved_by = $2,
		    approved_at = now(), change_note = $3, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`, id, actor, note)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	if _, err := r.data.db.Pool.Exec(ctx, `
		INSERT INTO pricing_rule_audit_logs (
			rule_type, rule_id, action, old_status, new_status, old_approval_status,
			new_approval_status, version, actor, note
		)
		VALUES ('STANDARD_LIMIT', $1, 'APPROVE', $2, 'ACTIVE', $3, 'APPROVED', $4, $5, $6)`,
		id, item.Status, item.ApprovalStatus, item.Version, actor, note,
	); err != nil {
		return nil, err
	}
	return r.GetStandardLimit(ctx, id)
}

func scanFeeSchedules(rows pgx.Rows) ([]*biz.FeeSchedule, error) {
	var list []*biz.FeeSchedule
	for rows.Next() {
		item, err := scanFeeSchedule(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanFeeSchedule(row pgx.Row) (*biz.FeeSchedule, error) {
	item := &biz.FeeSchedule{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.FeeType, &item.CalculationMethod, &item.Currency,
		&item.FixedAmount, &item.RatePercent, &item.MinAmount, &item.MaxAmount, &item.Channel,
		&item.ProductCode, &item.EffectiveFrom, &item.EffectiveTo, &item.Description,
		&item.Status, &item.ApprovalStatus, &item.Version, &item.ApprovedBy, &item.ApprovedAt,
		&item.ChangeNote, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanTaxRules(rows pgx.Rows) ([]*biz.TaxRule, error) {
	var list []*biz.TaxRule
	for rows.Next() {
		item, err := scanTaxRule(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanTaxRule(row pgx.Row) (*biz.TaxRule, error) {
	item := &biz.TaxRule{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.TaxType, &item.RatePercent, &item.Inclusive,
		&item.Jurisdiction, &item.EffectiveFrom, &item.EffectiveTo, &item.Description,
		&item.Status, &item.ApprovalStatus, &item.Version, &item.ApprovedBy, &item.ApprovedAt,
		&item.ChangeNote, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanStandardLimits(rows pgx.Rows) ([]*biz.StandardLimit, error) {
	var list []*biz.StandardLimit
	for rows.Next() {
		item, err := scanStandardLimit(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanStandardLimit(row pgx.Row) (*biz.StandardLimit, error) {
	item := &biz.StandardLimit{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.LimitType, &item.SubjectType, &item.Currency,
		&item.MinAmount, &item.PerTxnAmount, &item.DailyAmount, &item.MonthlyAmount,
		&item.CountLimit, &item.Channel, &item.ProductCode, &item.EffectiveFrom, &item.EffectiveTo,
		&item.Description, &item.Status, &item.ApprovalStatus, &item.Version, &item.ApprovedBy,
		&item.ApprovedAt, &item.ChangeNote, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
