package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListCodeSets(ctx context.Context, filter biz.PageFilter) ([]*biz.CodeSet, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, description, is_system, status, created_at, updated_at
		FROM code_sets
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%')
		ORDER BY code ASC
		LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var list []*biz.CodeSet
	for rows.Next() {
		item := &biz.CodeSet{}
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.IsSystem, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, "", err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetCodeSet(ctx context.Context, id string) (*biz.CodeSet, error) {
	return r.getCodeSet(ctx, `id = $1`, id)
}

func (r *MdmRepo) GetCodeSetByCode(ctx context.Context, code string) (*biz.CodeSet, error) {
	return r.getCodeSet(ctx, `code = $1`, code)
}

func (r *MdmRepo) getCodeSet(ctx context.Context, predicate string, arg string) (*biz.CodeSet, error) {
	item := &biz.CodeSet{}
	err := r.data.db.Pool.QueryRow(ctx, fmt.Sprintf(`
		SELECT id::text, code, name, description, is_system, status, created_at, updated_at
		FROM code_sets
		WHERE %s AND deleted_at IS NULL`, predicate), arg).
		Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.IsSystem, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func (r *MdmRepo) CreateCodeSet(ctx context.Context, codeSet *biz.CodeSet) (*biz.CodeSet, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO code_sets (code, name, description, is_system, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text`, codeSet.Code, codeSet.Name, codeSet.Description, codeSet.IsSystem, codeSet.Status).
		Scan(&codeSet.ID)
	if err != nil {
		return nil, err
	}
	return r.GetCodeSet(ctx, codeSet.ID)
}

func (r *MdmRepo) UpdateCodeSet(ctx context.Context, codeSet *biz.CodeSet) (*biz.CodeSet, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE code_sets
		SET code = $2, name = $3, description = $4, is_system = $5, status = $6, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		codeSet.ID, codeSet.Code, codeSet.Name, codeSet.Description, codeSet.IsSystem, codeSet.Status,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetCodeSet(ctx, codeSet.ID)
}

func (r *MdmRepo) DeleteCodeSet(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "code_sets", id)
}

func (r *MdmRepo) ListCodeItems(ctx context.Context, filter biz.CodeItemFilter) ([]*biz.CodeItem, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT ci.id::text, ci.code_set_id::text, cs.code, ci.code, ci.name, ci.value,
		       COALESCE(ci.parent_id::text, ''), ci.sort_order, ci.color, ci.icon,
		       ci.metadata::text, ci.is_default, ci.is_system, ci.status,
		       COALESCE(to_char(ci.effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(ci.effective_to, 'YYYY-MM-DD'), ''),
		       ci.created_at, ci.updated_at
		FROM code_items ci
		JOIN code_sets cs ON cs.id = ci.code_set_id
		WHERE ci.deleted_at IS NULL
		  AND cs.code = $1
		  AND ($2 = '' OR ci.status = $2)
		  AND ($3 = '' OR ci.code ILIKE '%' || $3 || '%' OR ci.name ILIKE '%' || $3 || '%')
		ORDER BY ci.sort_order ASC, ci.code ASC
		LIMIT $4 OFFSET $5`, filter.CodeSetCode, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanCodeItems(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetCodeItem(ctx context.Context, id string) (*biz.CodeItem, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT ci.id::text, ci.code_set_id::text, cs.code, ci.code, ci.name, ci.value,
		       COALESCE(ci.parent_id::text, ''), ci.sort_order, ci.color, ci.icon,
		       ci.metadata::text, ci.is_default, ci.is_system, ci.status,
		       COALESCE(to_char(ci.effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(ci.effective_to, 'YYYY-MM-DD'), ''),
		       ci.created_at, ci.updated_at
		FROM code_items ci
		JOIN code_sets cs ON cs.id = ci.code_set_id
		WHERE ci.id = $1 AND ci.deleted_at IS NULL`, id)
	return scanCodeItem(row)
}

func (r *MdmRepo) CreateCodeItem(ctx context.Context, item *biz.CodeItem) (*biz.CodeItem, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO code_items (
			code_set_id, code, name, value, parent_id, sort_order, color, icon, metadata,
			is_default, is_system, status, effective_from, effective_to
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, $6, $7, $8, COALESCE(NULLIF($9, '')::jsonb, '{}'::jsonb),
		        $10, $11, $12, NULLIF($13, '')::date, NULLIF($14, '')::date)
		RETURNING id::text`,
		item.CodeSetID, item.Code, item.Name, item.Value, item.ParentID, item.SortOrder, item.Color, item.Icon, item.MetadataJSON,
		item.IsDefault, item.IsSystem, item.Status, item.EffectiveFrom, item.EffectiveTo,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetCodeItem(ctx, item.ID)
}

func (r *MdmRepo) UpdateCodeItem(ctx context.Context, item *biz.CodeItem) (*biz.CodeItem, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE code_items
		SET code = $2, name = $3, value = $4, parent_id = NULLIF($5, '')::uuid,
		    sort_order = $6, color = $7, icon = $8,
		    metadata = COALESCE(NULLIF($9, '')::jsonb, '{}'::jsonb),
		    is_default = $10, is_system = $11, status = $12,
		    effective_from = NULLIF($13, '')::date, effective_to = NULLIF($14, '')::date,
		    updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.Value, item.ParentID, item.SortOrder, item.Color, item.Icon,
		item.MetadataJSON, item.IsDefault, item.IsSystem, item.Status, item.EffectiveFrom, item.EffectiveTo,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetCodeItem(ctx, item.ID)
}

func (r *MdmRepo) DeleteCodeItem(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "code_items", id)
}

func scanCodeItems(rows pgx.Rows) ([]*biz.CodeItem, error) {
	var list []*biz.CodeItem
	for rows.Next() {
		item, err := scanCodeItem(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanCodeItem(row pgx.Row) (*biz.CodeItem, error) {
	item := &biz.CodeItem{}
	err := row.Scan(
		&item.ID, &item.CodeSetID, &item.CodeSetCode, &item.Code, &item.Name, &item.Value,
		&item.ParentID, &item.SortOrder, &item.Color, &item.Icon, &item.MetadataJSON, &item.IsDefault,
		&item.IsSystem, &item.Status, &item.EffectiveFrom, &item.EffectiveTo, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
