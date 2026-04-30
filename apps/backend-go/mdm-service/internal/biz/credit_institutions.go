package biz

import (
	"context"
	"strings"
	"time"
)

type CreditInstitution struct {
	ID            string
	Code          string
	Name          string
	ShortName     string
	Address       string
	Phone         string
	Email         string
	LicenseNumber string
	IssuedDate    string
	TaxCode       string
	Website       string
	Note          string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (uc *MdmUsecase) ListCreditInstitutions(ctx context.Context, filter PageFilter) ([]*CreditInstitution, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListCreditInstitutions(ctx, filter)
}

func (uc *MdmUsecase) GetCreditInstitution(ctx context.Context, id string) (*CreditInstitution, error) {
	return uc.repo.GetCreditInstitution(ctx, id)
}

func (uc *MdmUsecase) CreateCreditInstitution(ctx context.Context, item *CreditInstitution) (*CreditInstitution, error) {
	normalizeCreditInstitution(item)
	return uc.repo.CreateCreditInstitution(ctx, item)
}

func (uc *MdmUsecase) UpdateCreditInstitution(ctx context.Context, item *CreditInstitution) (*CreditInstitution, error) {
	normalizeCreditInstitution(item)
	return uc.repo.UpdateCreditInstitution(ctx, item)
}

func (uc *MdmUsecase) DeleteCreditInstitution(ctx context.Context, id string) error {
	return uc.repo.DeleteCreditInstitution(ctx, id)
}

func normalizeCreditInstitution(item *CreditInstitution) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.ShortName = strings.TrimSpace(item.ShortName)
	item.Address = strings.TrimSpace(item.Address)
	item.Phone = strings.TrimSpace(item.Phone)
	item.Email = strings.TrimSpace(item.Email)
	item.LicenseNumber = strings.TrimSpace(item.LicenseNumber)
	item.IssuedDate = strings.TrimSpace(item.IssuedDate)
	item.TaxCode = strings.TrimSpace(item.TaxCode)
	item.Website = strings.TrimSpace(item.Website)
	item.Note = strings.TrimSpace(item.Note)
	item.Status = upperDefault(item.Status, "ACTIVE")
}
