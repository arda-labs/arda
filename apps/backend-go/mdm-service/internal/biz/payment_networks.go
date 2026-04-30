package biz

import (
	"context"
	"strings"
	"time"
)

type BankBranch struct {
	ID              string
	InstitutionCode string
	Code            string
	Name            string
	BranchType      string
	Address         string
	ProvinceCode    string
	Phone           string
	SwiftCode       string
	NapasCode       string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PaymentNetwork struct {
	ID                 string
	Code               string
	Name               string
	NetworkType        string
	ClearingMethod     string
	SettlementCurrency string
	Operator           string
	Availability       string
	Description        string
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (uc *MdmUsecase) ListBankBranches(ctx context.Context, filter PageFilter) ([]*BankBranch, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListBankBranches(ctx, filter)
}

func (uc *MdmUsecase) CreateBankBranch(ctx context.Context, item *BankBranch) (*BankBranch, error) {
	normalizeBankBranch(item)
	return uc.repo.CreateBankBranch(ctx, item)
}

func (uc *MdmUsecase) UpdateBankBranch(ctx context.Context, item *BankBranch) (*BankBranch, error) {
	normalizeBankBranch(item)
	return uc.repo.UpdateBankBranch(ctx, item)
}

func (uc *MdmUsecase) DeleteBankBranch(ctx context.Context, id string) error {
	return uc.repo.DeleteBankBranch(ctx, id)
}

func (uc *MdmUsecase) ListPaymentNetworks(ctx context.Context, filter PageFilter) ([]*PaymentNetwork, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListPaymentNetworks(ctx, filter)
}

func (uc *MdmUsecase) CreatePaymentNetwork(ctx context.Context, item *PaymentNetwork) (*PaymentNetwork, error) {
	normalizePaymentNetwork(item)
	return uc.repo.CreatePaymentNetwork(ctx, item)
}

func (uc *MdmUsecase) UpdatePaymentNetwork(ctx context.Context, item *PaymentNetwork) (*PaymentNetwork, error) {
	normalizePaymentNetwork(item)
	return uc.repo.UpdatePaymentNetwork(ctx, item)
}

func (uc *MdmUsecase) DeletePaymentNetwork(ctx context.Context, id string) error {
	return uc.repo.DeletePaymentNetwork(ctx, id)
}

func normalizeBankBranch(item *BankBranch) {
	item.InstitutionCode = upperDefault(item.InstitutionCode, "")
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.BranchType = upperDefault(item.BranchType, "BRANCH")
	item.Address = strings.TrimSpace(item.Address)
	item.ProvinceCode = strings.TrimSpace(item.ProvinceCode)
	item.Phone = strings.TrimSpace(item.Phone)
	item.SwiftCode = upperDefault(item.SwiftCode, "")
	item.NapasCode = strings.TrimSpace(item.NapasCode)
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizePaymentNetwork(item *PaymentNetwork) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.NetworkType = upperDefault(item.NetworkType, "DOMESTIC")
	item.ClearingMethod = upperDefault(item.ClearingMethod, "")
	item.SettlementCurrency = upperDefault(item.SettlementCurrency, "VND")
	item.Operator = strings.TrimSpace(item.Operator)
	item.Availability = upperDefault(item.Availability, "24X7")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
}
