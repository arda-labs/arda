package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListSystemParameters(ctx context.Context, filter biz.SystemParameterFilter) ([]*biz.SystemParameter, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, key, name, group_code, value_type, value_text,
		       COALESCE(value_number, 0), COALESCE(value_boolean, false), value_json::text,
		       default_value, is_secret, is_editable, is_system, validation_rule::text,
		       description, status, updated_by, created_at, updated_at
		FROM system_parameters
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR group_code = $1)
		  AND ($2 = '' OR status = $2)
		  AND ($3 = '' OR key ILIKE '%' || $3 || '%' OR name ILIKE '%' || $3 || '%')
		ORDER BY group_code ASC, key ASC
		LIMIT $4 OFFSET $5`, filter.GroupCode, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanSystemParameters(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetSystemParameter(ctx context.Context, key string) (*biz.SystemParameter, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, key, name, group_code, value_type, value_text,
		       COALESCE(value_number, 0), COALESCE(value_boolean, false), value_json::text,
		       default_value, is_secret, is_editable, is_system, validation_rule::text,
		       description, status, updated_by, created_at, updated_at
		FROM system_parameters
		WHERE key = $1 AND deleted_at IS NULL`, key)
	return scanSystemParameter(row)
}

func (r *MdmRepo) CreateSystemParameter(ctx context.Context, param *biz.SystemParameter) (*biz.SystemParameter, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO system_parameters (
			key, name, group_code, value_type, value_text, value_number, value_boolean,
			value_json, default_value, is_secret, is_editable, is_system, validation_rule,
			description, status, updated_by
		)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, 0), $7,
		        COALESCE(NULLIF($8, '')::jsonb, '{}'::jsonb), $9, $10, $11, $12,
		        COALESCE(NULLIF($13, '')::jsonb, '{}'::jsonb), $14, $15, $16)
		RETURNING id::text`,
		param.Key, param.Name, param.GroupCode, param.ValueType, param.ValueText, param.ValueNumber, param.ValueBoolean,
		param.ValueJSON, param.DefaultValue, param.IsSecret, param.IsEditable, param.IsSystem,
		param.ValidationRuleJSON, param.Description, param.Status, param.UpdatedBy,
	).Scan(&param.ID)
	if err != nil {
		return nil, err
	}
	return r.GetSystemParameter(ctx, param.Key)
}

func (r *MdmRepo) UpdateSystemParameter(ctx context.Context, param *biz.SystemParameter) (*biz.SystemParameter, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE system_parameters
		SET name = $2, group_code = $3, value_type = $4, value_text = $5,
		    value_number = NULLIF($6, 0), value_boolean = $7,
		    value_json = COALESCE(NULLIF($8, '')::jsonb, '{}'::jsonb),
		    default_value = $9, is_secret = $10, is_editable = $11, is_system = $12,
		    validation_rule = COALESCE(NULLIF($13, '')::jsonb, '{}'::jsonb),
		    description = $14, status = $15, updated_by = $16, updated_at = now()
		WHERE key = $1 AND deleted_at IS NULL`,
		param.Key, param.Name, param.GroupCode, param.ValueType, param.ValueText, param.ValueNumber, param.ValueBoolean,
		param.ValueJSON, param.DefaultValue, param.IsSecret, param.IsEditable, param.IsSystem,
		param.ValidationRuleJSON, param.Description, param.Status, param.UpdatedBy,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetSystemParameter(ctx, param.Key)
}

func (r *MdmRepo) DeleteSystemParameter(ctx context.Context, key string) error {
	tag, err := r.data.db.Pool.Exec(ctx,
		`UPDATE system_parameters SET deleted_at = now(), status = 'DELETED', updated_at = now() WHERE key = $1 AND deleted_at IS NULL`,
		key,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func scanSystemParameters(rows pgx.Rows) ([]*biz.SystemParameter, error) {
	var list []*biz.SystemParameter
	for rows.Next() {
		item, err := scanSystemParameter(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanSystemParameter(row pgx.Row) (*biz.SystemParameter, error) {
	item := &biz.SystemParameter{}
	err := row.Scan(
		&item.ID, &item.Key, &item.Name, &item.GroupCode, &item.ValueType, &item.ValueText,
		&item.ValueNumber, &item.ValueBoolean, &item.ValueJSON, &item.DefaultValue,
		&item.IsSecret, &item.IsEditable, &item.IsSystem, &item.ValidationRuleJSON,
		&item.Description, &item.Status, &item.UpdatedBy, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
