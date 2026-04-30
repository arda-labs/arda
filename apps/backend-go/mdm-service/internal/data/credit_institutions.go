package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListCreditInstitutions(ctx context.Context, filter biz.PageFilter) ([]*biz.CreditInstitution, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, `
		SELECT id::text, code, name, short_name, address, phone, email, license_number,
		       COALESCE(to_char(issued_date, 'YYYY-MM-DD'), ''), tax_code, website, note, status, created_at, updated_at
		FROM credit_institutions
		WHERE deleted_at IS NULL
		  AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR short_name ILIKE '%' || $2 || '%'
		           OR tax_code ILIKE '%' || $2 || '%' OR license_number ILIKE '%' || $2 || '%')
		ORDER BY code ASC, name ASC
		LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	list, err := scanCreditInstitutions(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) GetCreditInstitution(ctx context.Context, id string) (*biz.CreditInstitution, error) {
	row := r.data.db.Pool.QueryRow(ctx, `
		SELECT id::text, code, name, short_name, address, phone, email, license_number,
		       COALESCE(to_char(issued_date, 'YYYY-MM-DD'), ''), tax_code, website, note, status, created_at, updated_at
		FROM credit_institutions
		WHERE id = $1 AND deleted_at IS NULL`, id)
	return scanCreditInstitution(row)
}

func (r *MdmRepo) CreateCreditInstitution(ctx context.Context, item *biz.CreditInstitution) (*biz.CreditInstitution, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO credit_institutions (
			code, name, short_name, address, phone, email, license_number, issued_date,
			tax_code, website, note, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, '')::date, $9, $10, $11, $12)
		RETURNING id::text`,
		item.Code, item.Name, item.ShortName, item.Address, item.Phone, item.Email, item.LicenseNumber,
		item.IssuedDate, item.TaxCode, item.Website, item.Note, item.Status,
	).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.GetCreditInstitution(ctx, item.ID)
}

func (r *MdmRepo) UpdateCreditInstitution(ctx context.Context, item *biz.CreditInstitution) (*biz.CreditInstitution, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE credit_institutions
		SET code = $2, name = $3, short_name = $4, address = $5, phone = $6, email = $7,
		    license_number = $8, issued_date = NULLIF($9, '')::date, tax_code = $10,
		    website = $11, note = $12, status = $13, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.ShortName, item.Address, item.Phone, item.Email,
		item.LicenseNumber, item.IssuedDate, item.TaxCode, item.Website, item.Note, item.Status,
	)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.GetCreditInstitution(ctx, item.ID)
}

func (r *MdmRepo) DeleteCreditInstitution(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "credit_institutions", id)
}

func scanCreditInstitutions(rows pgx.Rows) ([]*biz.CreditInstitution, error) {
	var list []*biz.CreditInstitution
	for rows.Next() {
		item, err := scanCreditInstitution(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanCreditInstitution(row pgx.Row) (*biz.CreditInstitution, error) {
	item := &biz.CreditInstitution{}
	err := row.Scan(
		&item.ID, &item.Code, &item.Name, &item.ShortName, &item.Address, &item.Phone, &item.Email,
		&item.LicenseNumber, &item.IssuedDate, &item.TaxCode, &item.Website, &item.Note, &item.Status,
		&item.CreatedAt, &item.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
