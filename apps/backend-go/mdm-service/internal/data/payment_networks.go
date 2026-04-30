package data

import (
	"context"
	"errors"

	"github.com/arda-labs/arda/arda-be-go/pkg/pagination"
	"github.com/arda-labs/arda/arda-be-go/services/mdm-service/internal/biz"
	"github.com/jackc/pgx/v5"
)

func (r *MdmRepo) ListBankBranches(ctx context.Context, filter biz.PageFilter) ([]*biz.BankBranch, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, bankBranchSelect()+`
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR institution_code ILIKE '%' || $2 || '%' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%'
		       OR swift_code ILIKE '%' || $2 || '%' OR napas_code ILIKE '%' || $2 || '%' OR province_code ILIKE '%' || $2 || '%')
		ORDER BY institution_code ASC, code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanBankBranches(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreateBankBranch(ctx context.Context, item *biz.BankBranch) (*biz.BankBranch, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO bank_branches (institution_code, code, name, branch_type, address, province_code, phone, swift_code, napas_code, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id::text`,
		item.InstitutionCode, item.Code, item.Name, item.BranchType, item.Address, item.ProvinceCode, item.Phone, item.SwiftCode, item.NapasCode, item.Status).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getBankBranch(ctx, item.ID)
}

func (r *MdmRepo) UpdateBankBranch(ctx context.Context, item *biz.BankBranch) (*biz.BankBranch, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE bank_branches
		SET institution_code=$2, code=$3, name=$4, branch_type=$5, address=$6, province_code=$7,
		    phone=$8, swift_code=$9, napas_code=$10, status=$11, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.InstitutionCode, item.Code, item.Name, item.BranchType, item.Address, item.ProvinceCode, item.Phone, item.SwiftCode, item.NapasCode, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getBankBranch(ctx, item.ID)
}

func (r *MdmRepo) DeleteBankBranch(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "bank_branches", id)
}

func (r *MdmRepo) ListPaymentNetworks(ctx context.Context, filter biz.PageFilter) ([]*biz.PaymentNetwork, string, error) {
	page := pagination.Normalize(filter.PageSize, filter.PageToken)
	rows, err := r.data.db.Pool.Query(ctx, paymentNetworkSelect()+`
		WHERE deleted_at IS NULL AND ($1 = '' OR status = $1)
		  AND ($2 = '' OR code ILIKE '%' || $2 || '%' OR name ILIKE '%' || $2 || '%' OR network_type ILIKE '%' || $2 || '%'
		       OR clearing_method ILIKE '%' || $2 || '%' OR operator ILIKE '%' || $2 || '%')
		ORDER BY network_type ASC, code ASC LIMIT $3 OFFSET $4`, filter.Status, filter.Keyword, page.Limit+1, page.Offset)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	list, err := scanPaymentNetworks(rows)
	if err != nil {
		return nil, "", err
	}
	next := pagination.NextOffsetToken(len(list), page.Limit, page.Offset)
	if len(list) > page.Limit {
		list = list[:page.Limit]
	}
	return list, next, nil
}

func (r *MdmRepo) CreatePaymentNetwork(ctx context.Context, item *biz.PaymentNetwork) (*biz.PaymentNetwork, error) {
	err := r.data.db.Pool.QueryRow(ctx, `
		INSERT INTO payment_networks (code, name, network_type, clearing_method, settlement_currency, operator, availability, description, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id::text`,
		item.Code, item.Name, item.NetworkType, item.ClearingMethod, item.SettlementCurrency, item.Operator, item.Availability, item.Description, item.Status).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return r.getPaymentNetwork(ctx, item.ID)
}

func (r *MdmRepo) UpdatePaymentNetwork(ctx context.Context, item *biz.PaymentNetwork) (*biz.PaymentNetwork, error) {
	tag, err := r.data.db.Pool.Exec(ctx, `
		UPDATE payment_networks
		SET code=$2, name=$3, network_type=$4, clearing_method=$5, settlement_currency=$6,
		    operator=$7, availability=$8, description=$9, status=$10, updated_at=now()
		WHERE id=$1 AND deleted_at IS NULL`,
		item.ID, item.Code, item.Name, item.NetworkType, item.ClearingMethod, item.SettlementCurrency, item.Operator, item.Availability, item.Description, item.Status)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, biz.ErrNotFound
	}
	return r.getPaymentNetwork(ctx, item.ID)
}

func (r *MdmRepo) DeletePaymentNetwork(ctx context.Context, id string) error {
	return softDelete(ctx, r.data, "payment_networks", id)
}

func (r *MdmRepo) getBankBranch(ctx context.Context, id string) (*biz.BankBranch, error) {
	row := r.data.db.Pool.QueryRow(ctx, bankBranchSelect()+` WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanBankBranch(row)
}

func (r *MdmRepo) getPaymentNetwork(ctx context.Context, id string) (*biz.PaymentNetwork, error) {
	row := r.data.db.Pool.QueryRow(ctx, paymentNetworkSelect()+` WHERE id=$1 AND deleted_at IS NULL`, id)
	return scanPaymentNetwork(row)
}

func bankBranchSelect() string {
	return `SELECT id::text, institution_code, code, name, branch_type, address, province_code, phone, swift_code, napas_code, status, created_at, updated_at FROM bank_branches`
}

func paymentNetworkSelect() string {
	return `SELECT id::text, code, name, network_type, clearing_method, settlement_currency, operator, availability, description, status, created_at, updated_at FROM payment_networks`
}

func scanBankBranches(rows pgx.Rows) ([]*biz.BankBranch, error) {
	var list []*biz.BankBranch
	for rows.Next() {
		item, err := scanBankBranch(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanBankBranch(row pgx.Row) (*biz.BankBranch, error) {
	item := &biz.BankBranch{}
	err := row.Scan(&item.ID, &item.InstitutionCode, &item.Code, &item.Name, &item.BranchType, &item.Address, &item.ProvinceCode, &item.Phone, &item.SwiftCode, &item.NapasCode, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}

func scanPaymentNetworks(rows pgx.Rows) ([]*biz.PaymentNetwork, error) {
	var list []*biz.PaymentNetwork
	for rows.Next() {
		item, err := scanPaymentNetwork(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, rows.Err()
}

func scanPaymentNetwork(row pgx.Row) (*biz.PaymentNetwork, error) {
	item := &biz.PaymentNetwork{}
	err := row.Scan(&item.ID, &item.Code, &item.Name, &item.NetworkType, &item.ClearingMethod, &item.SettlementCurrency, &item.Operator, &item.Availability, &item.Description, &item.Status, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, biz.ErrNotFound
	}
	return item, err
}
