package data

import (
	"context"
	"errors"
	"time"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func zeroTimeNil(t time.Time) any {
	if t.IsZero() {
		return nil
	}
	return t
}

func (r *MdmRepo) ListCurrencies(ctx context.Context, filter biz.PageFilter) ([]*biz.Currency, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, numeric_code, name, minor_unit, symbol, country_code, status, created_at, updated_at
		FROM currencies
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR country_code ILIKE '%' || $2 || '%')
		ORDER BY code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanCurrencies(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateCurrency(ctx context.Context, item *biz.Currency) (*biz.Currency, error) {
	if err := r.data.db.Pool.QueryRow(ctx, `INSERT INTO currencies (code,numeric_code,name,minor_unit,symbol,country_code,status) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id::text`,
		item.Code, item.NumericCode, item.Name, item.MinorUnit, item.Symbol, item.CountryCode, item.Status).Scan(&item.ID); err != nil {
		return nil, err
	}
	return r.getCurrency(ctx, item.ID)
}

func (r *MdmRepo) UpdateCurrency(ctx context.Context, item *biz.Currency) (*biz.Currency, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `UPDATE currencies SET code=$2,numeric_code=$3,name=$4,minor_unit=$5,symbol=$6,country_code=$7,status=$8,updated_at=now() WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.NumericCode, item.Name, item.MinorUnit, item.Symbol, item.CountryCode, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getCurrency(ctx, item.ID)
}

func (r *MdmRepo) DeleteCurrency(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "currencies", id)
}

func (r *MdmRepo) ListFxRateSources(ctx context.Context, filter biz.PageFilter) ([]*biz.FxRateSource, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, source_type, priority, timezone, description, status, created_at, updated_at
		FROM fx_rate_sources
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR source_type ILIKE '%' || $2 || '%')
		ORDER BY priority ASC, code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanFxRateSources(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateFxRateSource(ctx context.Context, item *biz.FxRateSource) (*biz.FxRateSource, error) {
	if err := r.data.db.Pool.QueryRow(ctx, `INSERT INTO fx_rate_sources (code,name,source_type,priority,timezone,description,status) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id::text`,
		item.Code, item.Name, item.SourceType, item.Priority, item.Timezone, item.Description, item.Status).Scan(&item.ID); err != nil {
		return nil, err
	}
	return r.getFxRateSource(ctx, item.ID)
}

func (r *MdmRepo) UpdateFxRateSource(ctx context.Context, item *biz.FxRateSource) (*biz.FxRateSource, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `UPDATE fx_rate_sources SET code=$2,name=$3,source_type=$4,priority=$5,timezone=$6,description=$7,status=$8,updated_at=now() WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.SourceType, item.Priority, item.Timezone, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getFxRateSource(ctx, item.ID)
}

func (r *MdmRepo) DeleteFxRateSource(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "fx_rate_sources", id)
}

func (r *MdmRepo) ListFxRates(ctx context.Context, filter biz.PageFilter) ([]*biz.FxRate, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, fxRateSelect()+`
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR base_currency ILIKE '%' || $2 || '%' OR quote_currency ILIKE '%' || $2 || '%' OR source_code ILIKE '%' || $2 || '%')
		ORDER BY rate_date DESC, base_currency ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanFxRates(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateFxRate(ctx context.Context, item *biz.FxRate) (*biz.FxRate, error) {
	if err := r.data.db.Pool.QueryRow(ctx, `INSERT INTO fx_rates (base_currency,quote_currency,source_code,rate_date,effective_at,buy_rate,sell_rate,mid_rate,approval_status,version,approved_by,change_note,status) VALUES ($1,$2,$3,$4::date,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id::text`,
		item.BaseCurrency, item.QuoteCurrency, item.SourceCode, item.RateDate, zeroTimeNil(item.EffectiveAt), item.BuyRate, item.SellRate, item.MidRate, item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote, item.Status).Scan(&item.ID); err != nil {
		return nil, err
	}
	return r.getFxRate(ctx, item.ID)
}

func (r *MdmRepo) UpdateFxRate(ctx context.Context, item *biz.FxRate) (*biz.FxRate, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `UPDATE fx_rates SET base_currency=$2,quote_currency=$3,source_code=$4,rate_date=$5::date,effective_at=$6,buy_rate=$7,sell_rate=$8,mid_rate=$9,approval_status=$10,version=$11,approved_by=$12,change_note=$13,status=$14,updated_at=now() WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.BaseCurrency, item.QuoteCurrency, item.SourceCode, item.RateDate, zeroTimeNil(item.EffectiveAt), item.BuyRate, item.SellRate, item.MidRate, item.ApprovalStatus, item.Version, item.ApprovedBy, item.ChangeNote, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getFxRate(ctx, item.ID)
}

func (r *MdmRepo) DeleteFxRate(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "fx_rates", id)
}

func (r *MdmRepo) ApproveFxRate(ctx context.Context, id, actor, note string) (*biz.FxRate, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `UPDATE fx_rates SET approval_status='APPROVED',status='ACTIVE',approved_by=$2,approved_at=now(),change_note=$3,updated_at=now() WHERE id=$1 AND deleted_at IS NULL`, id, actor, note)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getFxRate(ctx, id)
}

func (r *MdmRepo) getCurrency(ctx context.Context, id string) (*biz.Currency, error) {
	row := r.data.db.Pool.QueryRow(ctx, `SELECT id::text, code, numeric_code, name, minor_unit, symbol, country_code, status, created_at, updated_at FROM currencies WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanCurrency(row)
}

func (r *MdmRepo) getFxRateSource(ctx context.Context, id string) (*biz.FxRateSource, error) {
	row := r.data.db.Pool.QueryRow(ctx, `SELECT id::text, code, name, source_type, priority, timezone, description, status, created_at, updated_at FROM fx_rate_sources WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanFxRateSource(row)
}

func (r *MdmRepo) getFxRate(ctx context.Context, id string) (*biz.FxRate, error) {
	row := r.data.db.Pool.QueryRow(ctx, fxRateSelect()+` WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanFxRate(row)
}

func fxRateSelect() string {
	return `SELECT id::text, base_currency, quote_currency, source_code, to_char(rate_date,'YYYY-MM-DD'), COALESCE(effective_at, '0001-01-01 00:00:00+00'::timestamptz), buy_rate::float8, sell_rate::float8, mid_rate::float8, approval_status, version, approved_by, COALESCE(approved_at, '0001-01-01 00:00:00+00'::timestamptz), change_note, status, created_at, updated_at FROM fx_rates`
}

func scanCurrencies(rows pgx.Rows) ([]*biz.Currency, error) {
	var list []*biz.Currency
	for rows.Next() {
		item, err := scanCurrency(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanCurrency(row pgx.Row) (*biz.Currency, error) {
	item := &biz.Currency{}
	err := row.Scan(&item.ID, &item.Code, &item.NumericCode, &item.Name, &item.MinorUnit, &item.Symbol, &item.CountryCode, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanFxRateSources(rows pgx.Rows) ([]*biz.FxRateSource, error) {
	var list []*biz.FxRateSource
	for rows.Next() {
		item, err := scanFxRateSource(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanFxRateSource(row pgx.Row) (*biz.FxRateSource, error) {
	item := &biz.FxRateSource{}
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.SourceType, &item.Priority, &item.Timezone, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanFxRates(rows pgx.Rows) ([]*biz.FxRate, error) {
	var list []*biz.FxRate
	for rows.Next() {
		item, err := scanFxRate(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanFxRate(row pgx.Row) (*biz.FxRate, error) {
	item := &biz.FxRate{}
	err := row.Scan(&item.ID, &item.BaseCurrency, &item.QuoteCurrency, &item.SourceCode, &item.RateDate, &item.EffectiveAt, &item.BuyRate, &item.SellRate, &item.MidRate, &item.ApprovalStatus, &item.Version, &item.ApprovedBy, &item.ApprovedAt, &item.ChangeNote, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
