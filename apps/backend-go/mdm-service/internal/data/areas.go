package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListAreaTypes(ctx context.Context, filter biz.PageFilter) ([]*biz.AreaType, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, description, allow_hierarchy, status, created_at, updated_at
		FROM area_types
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%')
		ORDER BY code ASC
		LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var list []*biz.AreaType
	for rows.Next() {
		item := &biz.AreaType{}
		if err := rows.Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.AllowHierarchy, &item.Status, &item.CreatedAt, &item.UpdatedAt); err != nil {
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

func (r *MdmRepo) GetAreaType(ctx context.Context, id string) (*biz.AreaType, error) {
	item := &biz.AreaType{}
	err := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, description, allow_hierarchy, status, created_at, updated_at
		FROM area_types
		WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&item.ID, &item.Code, &item.Name, &item.Description, &item.AllowHierarchy, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func (r *MdmRepo) CreateAreaType(ctx context.Context, areaType *biz.AreaType) (*biz.AreaType, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO area_types (code, name, description, allow_hierarchy, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text`, areaType.Code, areaType.Name, areaType.Description, areaType.AllowHierarchy, areaType.Status).
		Scan(&areaType.ID)
	if err != nil {
		return nil, err
	}
	return r.GetAreaType(ctx, areaType.ID)
}

func (r *MdmRepo) UpdateAreaType(ctx context.Context, areaType *biz.AreaType) (*biz.AreaType, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE area_types
		SET code = $2, name = $3, description = $4, allow_hierarchy = $5, status = $6, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		areaType.ID, areaType.Code, areaType.Name, areaType.Description, areaType.AllowHierarchy, areaType.Status,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetAreaType(ctx, areaType.ID)
}

func (r *MdmRepo) DeleteAreaType(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "area_types", id)
}

func (r *MdmRepo) ListAreas(ctx context.Context, filter biz.AreaFilter) ([]*biz.Area, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT a.id::text, a.area_type_id::text, at.code, COALESCE(a.parent_id::text, ''), a.code, a.name,
		       a.description, a.manager_user_id, a.status,
		       COALESCE(to_char(a.effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(a.effective_to, 'YYYY-MM-DD'), ''),
		       a.metadata::text, a.created_at, a.updated_at
		FROM areas a
		JOIN area_types at ON at.id = a.area_type_id
		WHERE a.deleted_at IS NULL
		  AND ($1 = '' OR a.area_type_id::text = $1)
		  AND ($2 = '' OR a.parent_id::text = $2)
		  AND ($3 = '' OR a.status = $3)
		  AND ($4 = '' OR a.code ILIKE '%' || $4 || '%' OR a.name ILIKE '%' || $4 || '%')
		ORDER BY a.code ASC, a.name ASC
		LIMIT $5 OFFSET $6`, filter.AreaTypeID, filter.ParentID, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanAreas(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetArea(ctx context.Context, id string) (*biz.Area, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT a.id::text, a.area_type_id::text, at.code, COALESCE(a.parent_id::text, ''), a.code, a.name,
		       a.description, a.manager_user_id, a.status,
		       COALESCE(to_char(a.effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(a.effective_to, 'YYYY-MM-DD'), ''),
		       a.metadata::text, a.created_at, a.updated_at
		FROM areas a
		JOIN area_types at ON at.id = a.area_type_id
		WHERE a.id = $1 AND a.deleted_at IS NULL`, id)
	return scanArea(row)
}

func (r *MdmRepo) CreateArea(ctx context.Context, area *biz.Area) (*biz.Area, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO areas (
			area_type_id, parent_id, code, name, description, manager_user_id, status,
			effective_from, effective_to, metadata
		)
		VALUES ($1, NULLIF($2, '')::uuid, $3, $4, $5, $6, $7,
		        NULLIF($8, '')::date, NULLIF($9, '')::date, COALESCE(NULLIF($10, '')::jsonb, '{}'::jsonb))
		RETURNING id::text`,
		area.AreaTypeID, area.ParentID, area.Code, area.Name, area.Description, area.ManagerUserID, area.Status,
		area.EffectiveFrom, area.EffectiveTo, area.MetadataJSON,
	).Scan(&area.ID)
	if err != nil {
		return nil, err
	}
	return r.GetArea(ctx, area.ID)
}

func (r *MdmRepo) UpdateArea(ctx context.Context, area *biz.Area) (*biz.Area, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE areas
		SET area_type_id = $2, parent_id = NULLIF($3, '')::uuid, code = $4, name = $5,
		    description = $6, manager_user_id = $7, status = $8,
		    effective_from = NULLIF($9, '')::date, effective_to = NULLIF($10, '')::date,
		    metadata = COALESCE(NULLIF($11, '')::jsonb, '{}'::jsonb), updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		area.ID, area.AreaTypeID, area.ParentID, area.Code, area.Name, area.Description, area.ManagerUserID,
		area.Status, area.EffectiveFrom, area.EffectiveTo, area.MetadataJSON,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetArea(ctx, area.ID)
}

func (r *MdmRepo) DeleteArea(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "areas", id)
}

func (r *MdmRepo) AssignAreaAdministrativeUnit(ctx context.Context, item *biz.AreaAdministrativeUnit) (*biz.AreaAdministrativeUnit, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO area_administrative_units (area_id, administrative_unit_id, scope_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (area_id, administrative_unit_id)
		DO UPDATE SET scope_type = EXCLUDED.scope_type
		RETURNING id::text, area_id::text, administrative_unit_id::text, scope_type, created_at`,
		item.AreaID, item.AdministrativeUnitID, item.ScopeType,
	).Scan(&item.ID, &item.AreaID, &item.AdministrativeUnitID, &item.ScopeType, &item.CreatedAt)
	return item, err
}

func (r *MdmRepo) ListAreaAdministrativeUnits(ctx context.Context, areaID string) ([]*biz.AreaAdministrativeUnit, error) {
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, area_id::text, administrative_unit_id::text, scope_type, created_at
		FROM area_administrative_units
		WHERE area_id = $1
		ORDER BY created_at DESC`, areaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*biz.AreaAdministrativeUnit
	for rows.Next() {
		item := &biz.AreaAdministrativeUnit{}
		if err := rows.Scan(&item.ID, &item.AreaID, &item.AdministrativeUnitID, &item.ScopeType, &item.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func (r *MdmRepo) RemoveAreaAdministrativeUnit(ctx context.Context, areaID, administrativeUnitID string) error {
	tag, err := r.data.db.Pool.Exec(ctx,
		`DELETE FROM area_administrative_units WHERE area_id = $1 AND administrative_unit_id = $2`,
		areaID, administrativeUnitID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func scanAreas(rows pgx.Rows) ([]*biz.Area, error) {
	var list []*biz.Area
	for rows.Next() {
		item, err := scanArea(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanArea(row pgx.Row) (*biz.Area, error) {
	item := &biz.Area{}
	err := row.Scan(
		&item.ID, &item.AreaTypeID, &item.AreaTypeCode, &item.ParentID, &item.Code, &item.Name,
		&item.Description, &item.ManagerUserID, &item.Status, &item.EffectiveFrom, &item.EffectiveTo,
		&item.MetadataJSON, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
