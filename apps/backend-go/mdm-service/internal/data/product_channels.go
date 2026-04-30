package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListBankingProducts(ctx context.Context, filter biz.PageFilter) ([]*biz.BankingProduct, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, product_type, category, customer_segment, currency,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, created_at, updated_at
		FROM banking_products
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%'
		       OR product_type ILIKE '%' || $2 || '%' OR category ILIKE '%' || $2 || '%' OR customer_segment ILIKE '%' || $2 || '%')
		ORDER BY product_type ASC, code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanBankingProducts(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateBankingProduct(ctx context.Context, item *biz.BankingProduct) (*biz.BankingProduct, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO banking_products (code, name, product_type, category, customer_segment, currency, effective_from, effective_to, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::date, NULLIF($8, '')::date, $9, $10)
		RETURNING id::text`,
		item.Code, item.Name, item.ProductType, item.Category, item.CustomerSegment, item.Currency, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getBankingProduct(ctx, item.ID)
}

func (r *MdmRepo) UpdateBankingProduct(ctx context.Context, item *biz.BankingProduct) (*biz.BankingProduct, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE banking_products
		SET code=$2, name=$3, product_type=$4, category=$5, customer_segment=$6, currency=$7,
		    effective_from=NULLIF($8, '')::date, effective_to=NULLIF($9, '')::date,
		    description=$10, status=$11, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.ProductType, item.Category, item.CustomerSegment, item.Currency, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getBankingProduct(ctx, item.ID)
}

func (r *MdmRepo) DeleteBankingProduct(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "banking_products", id)
}

func (r *MdmRepo) ListServiceChannels(ctx context.Context, filter biz.PageFilter) ([]*biz.ServiceChannel, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, channel_type, availability, timezone, description, status, created_at, updated_at
		FROM service_channels
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR channel_type ILIKE '%' || $2 || '%' OR availability ILIKE '%' || $2 || '%')
		ORDER BY channel_type ASC, code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanServiceChannels(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateServiceChannel(ctx context.Context, item *biz.ServiceChannel) (*biz.ServiceChannel, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO service_channels (code, name, channel_type, availability, timezone, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id::text`,
		item.Code, item.Name, item.ChannelType, item.Availability, item.Timezone, item.Description, item.Status,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getServiceChannel(ctx, item.ID)
}

func (r *MdmRepo) UpdateServiceChannel(ctx context.Context, item *biz.ServiceChannel) (*biz.ServiceChannel, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE service_channels
		SET code=$2, name=$3, channel_type=$4, availability=$5, timezone=$6, description=$7, status=$8, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.ChannelType, item.Availability, item.Timezone, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getServiceChannel(ctx, item.ID)
}

func (r *MdmRepo) DeleteServiceChannel(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "service_channels", id)
}

func (r *MdmRepo) ListProductChannelRules(ctx context.Context, filter biz.PageFilter) ([]*biz.ProductChannelRule, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, productChannelRuleSelect()+`
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR product_code ILIKE '%' || $2 || '%' OR channel_code ILIKE '%' || $2 || '%'
		       OR transaction_type ILIKE '%' || $2 || '%' OR fee_schedule_code ILIKE '%' || $2 || '%' OR limit_profile_code ILIKE '%' || $2 || '%')
		ORDER BY product_code ASC, priority ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanProductChannelRules(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateProductChannelRule(ctx context.Context, item *biz.ProductChannelRule) (*biz.ProductChannelRule, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO product_channel_rules (product_code, channel_code, transaction_type, enabled, priority, fee_schedule_code, limit_profile_code, effective_from, effective_to, description, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '')::date, NULLIF($9, '')::date, $10, $11)
		RETURNING id::text`,
		item.ProductCode, item.ChannelCode, item.TransactionType, item.Enabled, item.Priority, item.FeeScheduleCode, item.LimitProfileCode, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getProductChannelRule(ctx, item.ID)
}

func (r *MdmRepo) UpdateProductChannelRule(ctx context.Context, item *biz.ProductChannelRule) (*biz.ProductChannelRule, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE product_channel_rules
		SET product_code=$2, channel_code=$3, transaction_type=$4, enabled=$5, priority=$6,
		    fee_schedule_code=$7, limit_profile_code=$8, effective_from=NULLIF($9, '')::date,
		    effective_to=NULLIF($10, '')::date, description=$11, status=$12, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.ProductCode, item.ChannelCode, item.TransactionType, item.Enabled, item.Priority, item.FeeScheduleCode, item.LimitProfileCode, item.EffectiveFrom, item.EffectiveTo, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getProductChannelRule(ctx, item.ID)
}

func (r *MdmRepo) DeleteProductChannelRule(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "product_channel_rules", id)
}

func (r *MdmRepo) getBankingProduct(ctx context.Context, id string) (*biz.BankingProduct, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, product_type, category, customer_segment, currency,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       description, status, created_at, updated_at
		FROM banking_products WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanBankingProduct(row)
}

func (r *MdmRepo) getServiceChannel(ctx context.Context, id string) (*biz.ServiceChannel, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, channel_type, availability, timezone, description, status, created_at, updated_at
		FROM service_channels WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanServiceChannel(row)
}

func (r *MdmRepo) getProductChannelRule(ctx context.Context, id string) (*biz.ProductChannelRule, error) {
	row := r.data.db.Pool.QueryRow(ctx, productChannelRuleSelect()+` WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanProductChannelRule(row)
}

func productChannelRuleSelect() string {
	return `SELECT id::text, product_code, channel_code, transaction_type, enabled, priority, fee_schedule_code, limit_profile_code, COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''), COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''), description, status, created_at, updated_at FROM product_channel_rules`
}

func scanBankingProducts(rows pgx.Rows) ([]*biz.BankingProduct, error) {
	var list []*biz.BankingProduct
	for rows.Next() {
		item, err := scanBankingProduct(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanBankingProduct(row pgx.Row) (*biz.BankingProduct, error) {
	item := &biz.BankingProduct{}
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.ProductType, &item.Category, &item.CustomerSegment, &item.Currency, &item.EffectiveFrom, &item.EffectiveTo, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanServiceChannels(rows pgx.Rows) ([]*biz.ServiceChannel, error) {
	var list []*biz.ServiceChannel
	for rows.Next() {
		item, err := scanServiceChannel(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanServiceChannel(row pgx.Row) (*biz.ServiceChannel, error) {
	item := &biz.ServiceChannel{}
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.ChannelType, &item.Availability, &item.Timezone, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanProductChannelRules(rows pgx.Rows) ([]*biz.ProductChannelRule, error) {
	var list []*biz.ProductChannelRule
	for rows.Next() {
		item, err := scanProductChannelRule(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanProductChannelRule(row pgx.Row) (*biz.ProductChannelRule, error) {
	item := &biz.ProductChannelRule{}
	err := row.Scan(&item.ID, &item.ProductCode, &item.ChannelCode, &item.TransactionType, &item.Enabled, &item.Priority, &item.FeeScheduleCode, &item.LimitProfileCode, &item.EffectiveFrom, &item.EffectiveTo, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
