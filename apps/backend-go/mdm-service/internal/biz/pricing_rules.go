package biz

import (
	"context"
	"strings"
	"time"
)

type FeeSchedule struct {
	ID                string
	Code              string
	Name              string
	FeeType           string
	CalculationMethod string
	Currency          string
	FixedAmount       float64
	RatePercent       float64
	MinAmount         float64
	MaxAmount         float64
	Channel           string
	ProductCode       string
	EffectiveFrom     string
	EffectiveTo       string
	Description       string
	Status            string
	ApprovalStatus    string
	Version           int
	ApprovedBy        string
	ApprovedAt        time.Time
	ChangeNote        string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TaxRule struct {
	ID             string
	Code           string
	Name           string
	TaxType        string
	RatePercent    float64
	Inclusive      bool
	Jurisdiction   string
	EffectiveFrom  string
	EffectiveTo    string
	Description    string
	Status         string
	ApprovalStatus string
	Version        int
	ApprovedBy     string
	ApprovedAt     time.Time
	ChangeNote     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type StandardLimit struct {
	ID             string
	Code           string
	Name           string
	LimitType      string
	SubjectType    string
	Currency       string
	MinAmount      float64
	PerTxnAmount   float64
	DailyAmount    float64
	MonthlyAmount  float64
	CountLimit     int
	Channel        string
	ProductCode    string
	EffectiveFrom  string
	EffectiveTo    string
	Description    string
	Status         string
	ApprovalStatus string
	Version        int
	ApprovedBy     string
	ApprovedAt     time.Time
	ChangeNote     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (uc *MdmUsecase) ListFeeSchedules(ctx context.Context, filter PageFilter) ([]*FeeSchedule, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListFeeSchedules(ctx, filter)
}

func (uc *MdmUsecase) GetFeeSchedule(ctx context.Context, id string) (*FeeSchedule, error) {
	return uc.repo.GetFeeSchedule(ctx, id)
}

func (uc *MdmUsecase) CreateFeeSchedule(ctx context.Context, item *FeeSchedule) (*FeeSchedule, error) {
	normalizeFeeSchedule(item)
	return uc.repo.CreateFeeSchedule(ctx, item)
}

func (uc *MdmUsecase) UpdateFeeSchedule(ctx context.Context, item *FeeSchedule) (*FeeSchedule, error) {
	normalizeFeeSchedule(item)
	return uc.repo.UpdateFeeSchedule(ctx, item)
}

func (uc *MdmUsecase) DeleteFeeSchedule(ctx context.Context, id string) error {
	return uc.repo.DeleteFeeSchedule(ctx, id)
}

func (uc *MdmUsecase) ApproveFeeSchedule(ctx context.Context, id, actor, note string) (*FeeSchedule, error) {
	return uc.repo.ApproveFeeSchedule(ctx, strings.TrimSpace(id), approvalActor(actor), strings.TrimSpace(note))
}

func (uc *MdmUsecase) ListTaxRules(ctx context.Context, filter PageFilter) ([]*TaxRule, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListTaxRules(ctx, filter)
}

func (uc *MdmUsecase) GetTaxRule(ctx context.Context, id string) (*TaxRule, error) {
	return uc.repo.GetTaxRule(ctx, id)
}

func (uc *MdmUsecase) CreateTaxRule(ctx context.Context, item *TaxRule) (*TaxRule, error) {
	normalizeTaxRule(item)
	return uc.repo.CreateTaxRule(ctx, item)
}

func (uc *MdmUsecase) UpdateTaxRule(ctx context.Context, item *TaxRule) (*TaxRule, error) {
	normalizeTaxRule(item)
	return uc.repo.UpdateTaxRule(ctx, item)
}

func (uc *MdmUsecase) DeleteTaxRule(ctx context.Context, id string) error {
	return uc.repo.DeleteTaxRule(ctx, id)
}

func (uc *MdmUsecase) ApproveTaxRule(ctx context.Context, id, actor, note string) (*TaxRule, error) {
	return uc.repo.ApproveTaxRule(ctx, strings.TrimSpace(id), approvalActor(actor), strings.TrimSpace(note))
}

func (uc *MdmUsecase) ListStandardLimits(ctx context.Context, filter PageFilter) ([]*StandardLimit, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListStandardLimits(ctx, filter)
}

func (uc *MdmUsecase) GetStandardLimit(ctx context.Context, id string) (*StandardLimit, error) {
	return uc.repo.GetStandardLimit(ctx, id)
}

func (uc *MdmUsecase) CreateStandardLimit(ctx context.Context, item *StandardLimit) (*StandardLimit, error) {
	normalizeStandardLimit(item)
	return uc.repo.CreateStandardLimit(ctx, item)
}

func (uc *MdmUsecase) UpdateStandardLimit(ctx context.Context, item *StandardLimit) (*StandardLimit, error) {
	normalizeStandardLimit(item)
	return uc.repo.UpdateStandardLimit(ctx, item)
}

func (uc *MdmUsecase) DeleteStandardLimit(ctx context.Context, id string) error {
	return uc.repo.DeleteStandardLimit(ctx, id)
}

func (uc *MdmUsecase) ApproveStandardLimit(ctx context.Context, id, actor, note string) (*StandardLimit, error) {
	return uc.repo.ApproveStandardLimit(ctx, strings.TrimSpace(id), approvalActor(actor), strings.TrimSpace(note))
}

func normalizeFeeSchedule(item *FeeSchedule) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.FeeType = upperDefault(item.FeeType, "SERVICE_FEE")
	item.CalculationMethod = upperDefault(item.CalculationMethod, "FIXED")
	item.Currency = upperDefault(item.Currency, "VND")
	item.Channel = upperDefault(item.Channel, "")
	item.ProductCode = upperDefault(item.ProductCode, "")
	item.EffectiveFrom = strings.TrimSpace(item.EffectiveFrom)
	item.EffectiveTo = strings.TrimSpace(item.EffectiveTo)
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
	item.ApprovalStatus = upperDefault(item.ApprovalStatus, "DRAFT")
	item.ApprovedBy = strings.TrimSpace(item.ApprovedBy)
	item.ChangeNote = strings.TrimSpace(item.ChangeNote)
	if item.Version <= 0 {
		item.Version = 1
	}
}

func normalizeTaxRule(item *TaxRule) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.TaxType = upperDefault(item.TaxType, "VAT")
	item.Jurisdiction = upperDefault(item.Jurisdiction, "VN")
	item.EffectiveFrom = strings.TrimSpace(item.EffectiveFrom)
	item.EffectiveTo = strings.TrimSpace(item.EffectiveTo)
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
	item.ApprovalStatus = upperDefault(item.ApprovalStatus, "DRAFT")
	item.ApprovedBy = strings.TrimSpace(item.ApprovedBy)
	item.ChangeNote = strings.TrimSpace(item.ChangeNote)
	if item.Version <= 0 {
		item.Version = 1
	}
}

func normalizeStandardLimit(item *StandardLimit) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.LimitType = upperDefault(item.LimitType, "TRANSACTION_AMOUNT")
	item.SubjectType = upperDefault(item.SubjectType, "CUSTOMER")
	item.Currency = upperDefault(item.Currency, "VND")
	item.Channel = upperDefault(item.Channel, "")
	item.ProductCode = upperDefault(item.ProductCode, "")
	item.EffectiveFrom = strings.TrimSpace(item.EffectiveFrom)
	item.EffectiveTo = strings.TrimSpace(item.EffectiveTo)
	item.Description = strings.TrimSpace(item.Description)
	item.Status = upperDefault(item.Status, "ACTIVE")
	item.ApprovalStatus = upperDefault(item.ApprovalStatus, "DRAFT")
	item.ApprovedBy = strings.TrimSpace(item.ApprovedBy)
	item.ChangeNote = strings.TrimSpace(item.ChangeNote)
	if item.Version <= 0 {
		item.Version = 1
	}
}

func approvalActor(actor string) string {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return "SYSTEM"
	}
	return actor
}
