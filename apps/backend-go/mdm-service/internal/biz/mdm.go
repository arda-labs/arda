package biz

import (
	"context"
	stderrors "errors"
	"strings"
)

var (
	ErrNotFound        = stderrors.New("MDM record not found")
	ErrReadOnly        = stderrors.New("MDM record is read-only")
	ErrInvalidArgument = stderrors.New("invalid MDM request")
)

type PageFilter struct {
	Status    string
	Keyword   string
	PageSize  int
	PageToken string
}

type MdmRepo interface {
	ListAdministrativeUnits(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error)
	GetAdministrativeUnit(ctx context.Context, id string) (*AdministrativeUnit, error)
	CreateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error)
	UpdateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error)
	DeleteAdministrativeUnit(ctx context.Context, id string) error
	ReplaceAdministrativeUnits(ctx context.Context, units []*AdministrativeUnit) error

	ListAreaTypes(ctx context.Context, filter PageFilter) ([]*AreaType, string, error)
	GetAreaType(ctx context.Context, id string) (*AreaType, error)
	CreateAreaType(ctx context.Context, areaType *AreaType) (*AreaType, error)
	UpdateAreaType(ctx context.Context, areaType *AreaType) (*AreaType, error)
	DeleteAreaType(ctx context.Context, id string) error

	ListAreas(ctx context.Context, filter AreaFilter) ([]*Area, string, error)
	GetArea(ctx context.Context, id string) (*Area, error)
	CreateArea(ctx context.Context, area *Area) (*Area, error)
	UpdateArea(ctx context.Context, area *Area) (*Area, error)
	DeleteArea(ctx context.Context, id string) error
	AssignAreaAdministrativeUnit(ctx context.Context, item *AreaAdministrativeUnit) (*AreaAdministrativeUnit, error)
	ListAreaAdministrativeUnits(ctx context.Context, areaID string) ([]*AreaAdministrativeUnit, error)
	RemoveAreaAdministrativeUnit(ctx context.Context, areaID, administrativeUnitID string) error

	ListCodeSets(ctx context.Context, filter PageFilter) ([]*CodeSet, string, error)
	GetCodeSet(ctx context.Context, id string) (*CodeSet, error)
	GetCodeSetByCode(ctx context.Context, code string) (*CodeSet, error)
	CreateCodeSet(ctx context.Context, codeSet *CodeSet) (*CodeSet, error)
	UpdateCodeSet(ctx context.Context, codeSet *CodeSet) (*CodeSet, error)
	DeleteCodeSet(ctx context.Context, id string) error

	ListCodeItems(ctx context.Context, filter CodeItemFilter) ([]*CodeItem, string, error)
	GetCodeItem(ctx context.Context, id string) (*CodeItem, error)
	CreateCodeItem(ctx context.Context, item *CodeItem) (*CodeItem, error)
	UpdateCodeItem(ctx context.Context, item *CodeItem) (*CodeItem, error)
	DeleteCodeItem(ctx context.Context, id string) error

	ListSystemParameters(ctx context.Context, filter SystemParameterFilter) ([]*SystemParameter, string, error)
	GetSystemParameter(ctx context.Context, key string) (*SystemParameter, error)
	CreateSystemParameter(ctx context.Context, param *SystemParameter) (*SystemParameter, error)
	UpdateSystemParameter(ctx context.Context, param *SystemParameter) (*SystemParameter, error)
	DeleteSystemParameter(ctx context.Context, key string) error

	ListCreditInstitutions(ctx context.Context, filter PageFilter) ([]*CreditInstitution, string, error)
	GetCreditInstitution(ctx context.Context, id string) (*CreditInstitution, error)
	CreateCreditInstitution(ctx context.Context, item *CreditInstitution) (*CreditInstitution, error)
	UpdateCreditInstitution(ctx context.Context, item *CreditInstitution) (*CreditInstitution, error)
	DeleteCreditInstitution(ctx context.Context, id string) error

	ListBusinessCalendars(ctx context.Context, filter PageFilter) ([]*BusinessCalendar, string, error)
	GetBusinessCalendar(ctx context.Context, id string) (*BusinessCalendar, error)
	GetBusinessCalendarByCode(ctx context.Context, code string) (*BusinessCalendar, error)
	CreateBusinessCalendar(ctx context.Context, item *BusinessCalendar) (*BusinessCalendar, error)
	UpdateBusinessCalendar(ctx context.Context, item *BusinessCalendar) (*BusinessCalendar, error)
	DeleteBusinessCalendar(ctx context.Context, id string) error

	ListWorkingHours(ctx context.Context, calendarID string) ([]*WorkingHour, error)
	CreateWorkingHour(ctx context.Context, item *WorkingHour) (*WorkingHour, error)
	UpdateWorkingHour(ctx context.Context, item *WorkingHour) (*WorkingHour, error)
	DeleteWorkingHour(ctx context.Context, id string) error

	ListCalendarExceptions(ctx context.Context, filter CalendarExceptionFilter) ([]*CalendarException, error)
	CreateCalendarException(ctx context.Context, item *CalendarException) (*CalendarException, error)
	UpdateCalendarException(ctx context.Context, item *CalendarException) (*CalendarException, error)
	DeleteCalendarException(ctx context.Context, id string) error

	ListFeeSchedules(ctx context.Context, filter PageFilter) ([]*FeeSchedule, string, error)
	GetFeeSchedule(ctx context.Context, id string) (*FeeSchedule, error)
	CreateFeeSchedule(ctx context.Context, item *FeeSchedule) (*FeeSchedule, error)
	UpdateFeeSchedule(ctx context.Context, item *FeeSchedule) (*FeeSchedule, error)
	DeleteFeeSchedule(ctx context.Context, id string) error
	ApproveFeeSchedule(ctx context.Context, id, actor, note string) (*FeeSchedule, error)

	ListTaxRules(ctx context.Context, filter PageFilter) ([]*TaxRule, string, error)
	GetTaxRule(ctx context.Context, id string) (*TaxRule, error)
	CreateTaxRule(ctx context.Context, item *TaxRule) (*TaxRule, error)
	UpdateTaxRule(ctx context.Context, item *TaxRule) (*TaxRule, error)
	DeleteTaxRule(ctx context.Context, id string) error
	ApproveTaxRule(ctx context.Context, id, actor, note string) (*TaxRule, error)

	ListStandardLimits(ctx context.Context, filter PageFilter) ([]*StandardLimit, string, error)
	GetStandardLimit(ctx context.Context, id string) (*StandardLimit, error)
	CreateStandardLimit(ctx context.Context, item *StandardLimit) (*StandardLimit, error)
	UpdateStandardLimit(ctx context.Context, item *StandardLimit) (*StandardLimit, error)
	DeleteStandardLimit(ctx context.Context, id string) error
	ApproveStandardLimit(ctx context.Context, id, actor, note string) (*StandardLimit, error)

	ListCurrencies(ctx context.Context, filter PageFilter) ([]*Currency, string, error)
	CreateCurrency(ctx context.Context, item *Currency) (*Currency, error)
	UpdateCurrency(ctx context.Context, item *Currency) (*Currency, error)
	DeleteCurrency(ctx context.Context, id string) error

	ListFxRateSources(ctx context.Context, filter PageFilter) ([]*FxRateSource, string, error)
	CreateFxRateSource(ctx context.Context, item *FxRateSource) (*FxRateSource, error)
	UpdateFxRateSource(ctx context.Context, item *FxRateSource) (*FxRateSource, error)
	DeleteFxRateSource(ctx context.Context, id string) error

	ListFxRates(ctx context.Context, filter PageFilter) ([]*FxRate, string, error)
	CreateFxRate(ctx context.Context, item *FxRate) (*FxRate, error)
	UpdateFxRate(ctx context.Context, item *FxRate) (*FxRate, error)
	DeleteFxRate(ctx context.Context, id string) error
	ApproveFxRate(ctx context.Context, id, actor, note string) (*FxRate, error)

	ListBankingProducts(ctx context.Context, filter PageFilter) ([]*BankingProduct, string, error)
	CreateBankingProduct(ctx context.Context, item *BankingProduct) (*BankingProduct, error)
	UpdateBankingProduct(ctx context.Context, item *BankingProduct) (*BankingProduct, error)
	DeleteBankingProduct(ctx context.Context, id string) error

	ListServiceChannels(ctx context.Context, filter PageFilter) ([]*ServiceChannel, string, error)
	CreateServiceChannel(ctx context.Context, item *ServiceChannel) (*ServiceChannel, error)
	UpdateServiceChannel(ctx context.Context, item *ServiceChannel) (*ServiceChannel, error)
	DeleteServiceChannel(ctx context.Context, id string) error

	ListProductChannelRules(ctx context.Context, filter PageFilter) ([]*ProductChannelRule, string, error)
	CreateProductChannelRule(ctx context.Context, item *ProductChannelRule) (*ProductChannelRule, error)
	UpdateProductChannelRule(ctx context.Context, item *ProductChannelRule) (*ProductChannelRule, error)
	DeleteProductChannelRule(ctx context.Context, id string) error
}

type MdmUsecase struct {
	repo MdmRepo
}

func NewMdmUsecase(repo MdmRepo) *MdmUsecase {
	return &MdmUsecase{repo: repo}
}

func normalizePageFilter(filter *PageFilter) {
	filter.Status = strings.ToUpper(strings.TrimSpace(filter.Status))
	filter.Keyword = strings.TrimSpace(filter.Keyword)
}

func upperDefault(value, fallback string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	return value
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func jsonDefault(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}"
	}
	return value
}
