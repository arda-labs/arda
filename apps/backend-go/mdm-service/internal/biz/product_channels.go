package biz

import (
	"context"
	"strings"
	"time"
)

type BankingProduct struct {
	ID              string
	Code            string
	Name            string
	ProductType     string
	Category        string
	CustomerSegment string
	Currency        string
	EffectiveFrom   string
	EffectiveTo     string
	Description     string
	Status          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ServiceChannel struct {
	ID           string
	Code         string
	Name         string
	ChannelType  string
	Availability string
	Timezone     string
	Description  string
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProductChannelRule struct {
	ID               string
	ProductCode      string
	ChannelCode      string
	TransactionType  string
	Enabled          bool
	Priority         int
	FeeScheduleCode  string
	LimitProfileCode string
	EffectiveFrom    string
	EffectiveTo      string
	Description      string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (uc *MdmUsecase) ListBankingProducts(ctx context.Context, filter PageFilter) ([]*BankingProduct, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListBankingProducts(ctx, filter)
}

func (uc *MdmUsecase) CreateBankingProduct(ctx context.Context, item *BankingProduct) (*BankingProduct, error) {
	normalizeBankingProduct(item)
	return uc.repo.CreateBankingProduct(ctx, item)
}

func (uc *MdmUsecase) UpdateBankingProduct(ctx context.Context, item *BankingProduct) (*BankingProduct, error) {
	normalizeBankingProduct(item)
	return uc.repo.UpdateBankingProduct(ctx, item)
}

func (uc *MdmUsecase) DeleteBankingProduct(ctx context.Context, id string) error {
	return uc.repo.DeleteBankingProduct(ctx, id)
}

func (uc *MdmUsecase) ListServiceChannels(ctx context.Context, filter PageFilter) ([]*ServiceChannel, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListServiceChannels(ctx, filter)
}

func (uc *MdmUsecase) CreateServiceChannel(ctx context.Context, item *ServiceChannel) (*ServiceChannel, error) {
	normalizeServiceChannel(item)
	return uc.repo.CreateServiceChannel(ctx, item)
}

func (uc *MdmUsecase) UpdateServiceChannel(ctx context.Context, item *ServiceChannel) (*ServiceChannel, error) {
	normalizeServiceChannel(item)
	return uc.repo.UpdateServiceChannel(ctx, item)
}

func (uc *MdmUsecase) DeleteServiceChannel(ctx context.Context, id string) error {
	return uc.repo.DeleteServiceChannel(ctx, id)
}

func (uc *MdmUsecase) ListProductChannelRules(ctx context.Context, filter PageFilter) ([]*ProductChannelRule, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListProductChannelRules(ctx, filter)
}

func (uc *MdmUsecase) CreateProductChannelRule(ctx context.Context, item *ProductChannelRule) (*ProductChannelRule, error) {
	normalizeProductChannelRule(item)
	return uc.repo.CreateProductChannelRule(ctx, item)
}

func (uc *MdmUsecase) UpdateProductChannelRule(ctx context.Context, item *ProductChannelRule) (*ProductChannelRule, error) {
	normalizeProductChannelRule(item)
	return uc.repo.UpdateProductChannelRule(ctx, item)
}

func (uc *MdmUsecase) DeleteProductChannelRule(ctx context.Context, id string) error {
	return uc.repo.DeleteProductChannelRule(ctx, id)
}

func normalizeBankingProduct(item *BankingProduct) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.ProductType = upperDefault(item.ProductType, "ACCOUNT")
	item.Category = upperDefault(item.Category, "")
	item.CustomerSegment = upperDefault(item.CustomerSegment, "")
	item.Currency = upperDefault(item.Currency, "")
	item.EffectiveFrom = strings.TrimSpace(item.EffectiveFrom)
	item.EffectiveTo = strings.TrimSpace(item.EffectiveTo)
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizeServiceChannel(item *ServiceChannel) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.ChannelType = upperDefault(item.ChannelType, "DIGITAL")
	item.Availability = upperDefault(item.Availability, "24X7")
	item.Timezone = defaultString(item.Timezone, "Asia/Ho_Chi_Minh")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizeProductChannelRule(item *ProductChannelRule) {
	item.ProductCode = upperDefault(item.ProductCode, "")
	item.ChannelCode = upperDefault(item.ChannelCode, "")
	item.TransactionType = upperDefault(item.TransactionType, "")
	item.FeeScheduleCode = upperDefault(item.FeeScheduleCode, "")
	item.LimitProfileCode = upperDefault(item.LimitProfileCode, "")
	item.EffectiveFrom = strings.TrimSpace(item.EffectiveFrom)
	item.EffectiveTo = strings.TrimSpace(item.EffectiveTo)
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
	if item.Priority <= 0 {
		item.Priority = 100
	}
}
