package biz

import (
	"context"
	"strings"
	"time"
)

type Currency struct {
	ID          string
	Code        string
	NumericCode string
	Name        string
	MinorUnit   int
	Symbol      string
	CountryCode string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FxRateSource struct {
	ID          string
	Code        string
	Name        string
	SourceType  string
	Priority    int
	Timezone    string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FxRate struct {
	ID             string
	BaseCurrency   string
	QuoteCurrency  string
	SourceCode     string
	RateDate       string
	EffectiveAt    time.Time
	BuyRate        float64
	SellRate       float64
	MidRate        float64
	ApprovalStatus string
	Version        int
	ApprovedBy     string
	ApprovedAt     time.Time
	ChangeNote     string
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (uc *MdmUsecase) ListCurrencies(ctx context.Context, filter PageFilter) ([]*Currency, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListCurrencies(ctx, filter)
}

func (uc *MdmUsecase) CreateCurrency(ctx context.Context, item *Currency) (*Currency, error) {
	normalizeCurrency(item)
	return uc.repo.CreateCurrency(ctx, item)
}

func (uc *MdmUsecase) UpdateCurrency(ctx context.Context, item *Currency) (*Currency, error) {
	normalizeCurrency(item)
	return uc.repo.UpdateCurrency(ctx, item)
}

func (uc *MdmUsecase) DeleteCurrency(ctx context.Context, id string) error {
	return uc.repo.DeleteCurrency(ctx, id)
}

func (uc *MdmUsecase) ListFxRateSources(ctx context.Context, filter PageFilter) ([]*FxRateSource, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListFxRateSources(ctx, filter)
}

func (uc *MdmUsecase) CreateFxRateSource(ctx context.Context, item *FxRateSource) (*FxRateSource, error) {
	normalizeFxRateSource(item)
	return uc.repo.CreateFxRateSource(ctx, item)
}

func (uc *MdmUsecase) UpdateFxRateSource(ctx context.Context, item *FxRateSource) (*FxRateSource, error) {
	normalizeFxRateSource(item)
	return uc.repo.UpdateFxRateSource(ctx, item)
}

func (uc *MdmUsecase) DeleteFxRateSource(ctx context.Context, id string) error {
	return uc.repo.DeleteFxRateSource(ctx, id)
}

func (uc *MdmUsecase) ListFxRates(ctx context.Context, filter PageFilter) ([]*FxRate, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListFxRates(ctx, filter)
}

func (uc *MdmUsecase) CreateFxRate(ctx context.Context, item *FxRate) (*FxRate, error) {
	normalizeFxRate(item)
	return uc.repo.CreateFxRate(ctx, item)
}

func (uc *MdmUsecase) UpdateFxRate(ctx context.Context, item *FxRate) (*FxRate, error) {
	normalizeFxRate(item)
	return uc.repo.UpdateFxRate(ctx, item)
}

func (uc *MdmUsecase) DeleteFxRate(ctx context.Context, id string) error {
	return uc.repo.DeleteFxRate(ctx, id)
}

func (uc *MdmUsecase) ApproveFxRate(ctx context.Context, id, actor, note string) (*FxRate, error) {
	return uc.repo.ApproveFxRate(ctx, strings.TrimSpace(id), approvalActor(actor), strings.TrimSpace(note))
}

func normalizeCurrency(item *Currency) {
	item.Code = upperDefault(item.Code, "")
	item.NumericCode = strings.TrimSpace(item.NumericCode)
	item.Name = strings.TrimSpace(item.Name)
	item.Symbol = strings.TrimSpace(item.Symbol)
	item.CountryCode = upperDefault(item.CountryCode, "")
	item.Status = upperDefault(item.Status, "ACTIVE")
}

func normalizeFxRateSource(item *FxRateSource) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.SourceType = upperDefault(item.SourceType, "MANUAL")
	item.Timezone = defaultString(item.Timezone, "Asia/Ho_Chi_Minh")
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
	if item.Priority <= 0 {
		item.Priority = 100
	}
}

func normalizeFxRate(item *FxRate) {
	item.BaseCurrency = upperDefault(item.BaseCurrency, "")
	item.QuoteCurrency = upperDefault(item.QuoteCurrency, "")
	item.SourceCode = upperDefault(item.SourceCode, "MANUAL")
	item.RateDate = strings.TrimSpace(item.RateDate)
	item.ApprovalStatus = upperDefault(item.ApprovalStatus, "DRAFT")
	item.ApprovedBy = strings.TrimSpace(item.ApprovedBy)
	item.ChangeNote = strings.TrimSpace(item.ChangeNote)
	item.Status = upperDefault(item.Status, "ACTIVE")
	if item.Version <= 0 {
		item.Version = 1
	}
}
