package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

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

func (r *MdmRepo) ReplaceAdministrativeUnits(ctx context.Context, units []*biz.AdministrativeUnit) error {
	tx, err := r.data.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM area_administrative_units`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM administrative_unit_mappings`); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM administrative_units`); err != nil {
		return err
	}

	provinceIDs := make(map[string]string)
	for _, unit := range units {
		if unit.Level != "PROVINCE" {
			continue
		}
		var id string
		if err := tx.QueryRow(ctx, `
			INSERT INTO administrative_units (
				code, name, full_name, short_name, level, unit_type, parent_id, path, sort_order,
				latitude, longitude, status, effective_from, effective_to, source, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, NULL, $7, $8,
			        NULLIF($9, 0), NULLIF($10, 0), $11, NULLIF($12, '')::date, NULLIF($13, '')::date,
			        $14, COALESCE(NULLIF($15, '')::jsonb, '{}'::jsonb))
			RETURNING id::text`,
			unit.Code, unit.Name, unit.FullName, unit.ShortName, unit.Level, unit.UnitType,
			unit.Path, unit.SortOrder, unit.Latitude, unit.Longitude, unit.Status,
			unit.EffectiveFrom, unit.EffectiveTo, unit.Source, unit.MetadataJSON,
		).Scan(&id); err != nil {
			return err
		}
		provinceIDs[unit.Code] = id
	}

	for _, unit := range units {
		if unit.Level == "PROVINCE" {
			continue
		}
		parentID, ok := provinceIDs[unit.ParentID]
		if !ok {
			return fmt.Errorf("province code %s for administrative unit %s not found", unit.ParentID, unit.Code)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO administrative_units (
				code, name, full_name, short_name, level, unit_type, parent_id, path, sort_order,
				latitude, longitude, status, effective_from, effective_to, source, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7::uuid, $8, $9,
			        NULLIF($10, 0), NULLIF($11, 0), $12, NULLIF($13, '')::date, NULLIF($14, '')::date,
			        $15, COALESCE(NULLIF($16, '')::jsonb, '{}'::jsonb))`,
			unit.Code, unit.Name, unit.FullName, unit.ShortName, unit.Level, unit.UnitType, parentID,
			unit.Path, unit.SortOrder, unit.Latitude, unit.Longitude, unit.Status,
			unit.EffectiveFrom, unit.EffectiveTo, unit.Source, unit.MetadataJSON,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
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
