package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jackc/pgx/v5"
)

type MdmRepo struct {
	data *Data
	log  *log.Helper
}

func NewMdmRepo(data *Data, logger log.Logger) biz.MdmRepo {
	return &MdmRepo{data: data, log: log.NewHelper(logger)}
}

func (r *MdmRepo) ListAdministrativeUnits(ctx context.Context, filter biz.AdministrativeUnitFilter) ([]*biz.AdministrativeUnit, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, full_name, short_name, level, unit_type,
		       COALESCE(parent_id::text, ''), path, sort_order,
		       COALESCE(latitude, 0), COALESCE(longitude, 0), status,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       source, metadata::text, created_at, updated_at
		FROM administrative_units
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR parent_id::text = $1)
		  AND ($2 = '' OR level = $2)
		  AND ($3 = '' OR status = $3)
		  AND ($4 = '' OR code ILIKE '%' || $4 || '%' OR name ILIKE '%' || $4 || '%' OR full_name ILIKE '%' || $4 || '%')
		ORDER BY sort_order ASC, name ASC, id ASC
		LIMIT $5 OFFSET $6`,
		filter.ParentID, filter.Level, filter.Status, filter.Keyword, page.Limit+1, page.Offset,
	)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanAdministrativeUnits(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetAdministrativeUnit(ctx context.Context, id string) (*biz.AdministrativeUnit, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, full_name, short_name, level, unit_type,
		       COALESCE(parent_id::text, ''), path, sort_order,
		       COALESCE(latitude, 0), COALESCE(longitude, 0), status,
		       COALESCE(to_char(effective_from, 'YYYY-MM-DD'), ''),
		       COALESCE(to_char(effective_to, 'YYYY-MM-DD'), ''),
		       source, metadata::text, created_at, updated_at
		FROM administrative_units
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanAdministrativeUnit(row)
}

func (r *MdmRepo) CreateAdministrativeUnit(ctx context.Context, unit *biz.AdministrativeUnit) (*biz.AdministrativeUnit, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO administrative_units (
			code, name, full_name, short_name, level, unit_type, parent_id, path, sort_order,
			latitude, longitude, status, effective_from, effective_to, source, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9,
		        NULLIF($10, 0), NULLIF($11, 0), $12, NULLIF($13, '')::date, NULLIF($14, '')::date,
		        $15, COALESCE(NULLIF($16, '')::jsonb, '{}'::jsonb))
		RETURNING id::text`,
		unit.Code, unit.Name, unit.FullName, unit.ShortName, unit.Level, unit.UnitType, unit.ParentID, unit.Path, unit.SortOrder,
		unit.Latitude, unit.Longitude, unit.Status, unit.EffectiveFrom, unit.EffectiveTo, unit.Source, unit.MetadataJSON,
	).Scan(&unit.ID)
	if err != nil {
		return nil, err
	}
	return r.GetAdministrativeUnit(ctx, unit.ID)
}

func (r *MdmRepo) UpdateAdministrativeUnit(ctx context.Context, unit *biz.AdministrativeUnit) (*biz.AdministrativeUnit, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE administrative_units
		SET code = $2, name = $3, full_name = $4, short_name = $5, level = $6, unit_type = $7,
		    parent_id = NULLIF($8, '')::uuid, path = $9, sort_order = $10,
		    latitude = NULLIF($11, 0), longitude = NULLIF($12, 0), status = $13,
		    effective_from = NULLIF($14, '')::date, effective_to = NULLIF($15, '')::date,
		    source = $16, metadata = COALESCE(NULLIF($17, '')::jsonb, '{}'::jsonb), updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		unit.ID, unit.Code, unit.Name, unit.FullName, unit.ShortName, unit.Level, unit.UnitType,
		unit.ParentID, unit.Path, unit.SortOrder, unit.Latitude, unit.Longitude, unit.Status,
		unit.EffectiveFrom, unit.EffectiveTo, unit.Source, unit.MetadataJSON,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetAdministrativeUnit(ctx, unit.ID)
}

func (r *MdmRepo) DeleteAdministrativeUnit(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "administrative_units", id)
}

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

func softDelete(ctx context.Context, data *Data, table, id string) error {
	tag, err := data.db.Pool.Exec(ctx,
		fmt.Sprintf(`UPDATE %s SET deleted_at = now(), status = 'DELETED', updated_at = now() WHERE id = $1 AND deleted_at IS NULL`, table),
		id,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return biz.ErrNotFound
	}
	return nil
}

func scanAdministrativeUnits(rows pgx.Rows) ([]*biz.AdministrativeUnit, error) {
	var list []*biz.AdministrativeUnit
	for rows.Next() {
		item := &biz.AdministrativeUnit{}
		if err := rows.Scan(
			&item.ID, &item.Code, &item.Name, &item.FullName, &item.ShortName, &item.Level, &item.UnitType,
			&item.ParentID, &item.Path, &item.SortOrder, &item.Latitude, &item.Longitude, &item.Status,
			&item.EffectiveFrom, &item.EffectiveTo, &item.Source, &item.MetadataJSON, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanAdministrativeUnit(row pgx.Row) (*biz.AdministrativeUnit, error) {
	item := &biz.AdministrativeUnit{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.FullName, &item.ShortName, &item.Level, &item.UnitType,
		&item.ParentID, &item.Path, &item.SortOrder, &item.Latitude, &item.Longitude, &item.Status,
		&item.EffectiveFrom, &item.EffectiveTo, &item.Source, &item.MetadataJSON, &item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
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
