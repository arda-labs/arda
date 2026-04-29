package biz

import (
	"context"
	stderrors "errors"
	"strings"
	"time"
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

type AdministrativeUnitFilter struct {
	ParentID string
	Level    string
	PageFilter
}

type AdministrativeUnit struct {
	ID            string
	Code          string
	Name          string
	FullName      string
	ShortName     string
	Level         string
	UnitType      string
	ParentID      string
	Path          string
	SortOrder     int
	Latitude      float64
	Longitude     float64
	Status        string
	EffectiveFrom string
	EffectiveTo   string
	Source        string
	MetadataJSON  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AdministrativeUnitNode struct {
	Unit     *AdministrativeUnit
	Children []*AdministrativeUnitNode
}

type AreaType struct {
	ID             string
	Code           string
	Name           string
	Description    string
	AllowHierarchy bool
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AreaFilter struct {
	AreaTypeID string
	ParentID   string
	PageFilter
}

type Area struct {
	ID            string
	AreaTypeID    string
	AreaTypeCode  string
	ParentID      string
	Code          string
	Name          string
	Description   string
	ManagerUserID string
	Status        string
	EffectiveFrom string
	EffectiveTo   string
	MetadataJSON  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AreaNode struct {
	Area     *Area
	Children []*AreaNode
}

type AreaAdministrativeUnit struct {
	ID                   string
	AreaID               string
	AdministrativeUnitID string
	ScopeType            string
	CreatedAt            time.Time
}

type CodeSet struct {
	ID          string
	Code        string
	Name        string
	Description string
	IsSystem    bool
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CodeItemFilter struct {
	CodeSetCode string
	PageFilter
}

type CodeItem struct {
	ID            string
	CodeSetID     string
	CodeSetCode   string
	Code          string
	Name          string
	Value         string
	ParentID      string
	SortOrder     int
	Color         string
	Icon          string
	MetadataJSON  string
	IsDefault     bool
	IsSystem      bool
	Status        string
	EffectiveFrom string
	EffectiveTo   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SystemParameterFilter struct {
	GroupCode string
	PageFilter
}

type SystemParameter struct {
	ID                 string
	Key                string
	Name               string
	GroupCode          string
	ValueType          string
	ValueText          string
	ValueNumber        float64
	ValueBoolean       bool
	ValueJSON          string
	DefaultValue       string
	IsSecret           bool
	IsEditable         bool
	IsSystem           bool
	ValidationRuleJSON string
	Description        string
	Status             string
	UpdatedBy          string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type MdmRepo interface {
	ListAdministrativeUnits(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error)
	GetAdministrativeUnit(ctx context.Context, id string) (*AdministrativeUnit, error)
	CreateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error)
	UpdateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error)
	DeleteAdministrativeUnit(ctx context.Context, id string) error

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
}

type MdmUsecase struct {
	repo MdmRepo
}

func NewMdmUsecase(repo MdmRepo) *MdmUsecase {
	return &MdmUsecase{repo: repo}
}

func (uc *MdmUsecase) ListAdministrativeUnits(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error) {
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListAdministrativeUnits(ctx, filter)
}

func (uc *MdmUsecase) ListAdministrativeUnitTree(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnitNode, error) {
	filter.PageSize = 1000
	filter.PageToken = ""
	units, _, err := uc.repo.ListAdministrativeUnits(ctx, filter)
	if err != nil {
		return nil, err
	}
	return buildAdministrativeUnitTree(units), nil
}

func (uc *MdmUsecase) ListProvinces(ctx context.Context, filter AdministrativeUnitFilter) ([]*AdministrativeUnit, string, error) {
	filter.Level = "PROVINCE"
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListAdministrativeUnits(ctx, filter)
}

func (uc *MdmUsecase) ListWards(ctx context.Context, provinceID string, filter PageFilter) ([]*AdministrativeUnit, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListAdministrativeUnits(ctx, AdministrativeUnitFilter{
		ParentID:   provinceID,
		Level:      "WARD",
		PageFilter: filter,
	})
}

func (uc *MdmUsecase) GetAdministrativeUnit(ctx context.Context, id string) (*AdministrativeUnit, error) {
	return uc.repo.GetAdministrativeUnit(ctx, id)
}

func (uc *MdmUsecase) CreateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error) {
	normalizeAdministrativeUnit(unit)
	return uc.repo.CreateAdministrativeUnit(ctx, unit)
}

func (uc *MdmUsecase) UpdateAdministrativeUnit(ctx context.Context, unit *AdministrativeUnit) (*AdministrativeUnit, error) {
	normalizeAdministrativeUnit(unit)
	return uc.repo.UpdateAdministrativeUnit(ctx, unit)
}

func (uc *MdmUsecase) DeleteAdministrativeUnit(ctx context.Context, id string) error {
	return uc.repo.DeleteAdministrativeUnit(ctx, id)
}

func (uc *MdmUsecase) ListAreaTypes(ctx context.Context, filter PageFilter) ([]*AreaType, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListAreaTypes(ctx, filter)
}

func (uc *MdmUsecase) GetAreaType(ctx context.Context, id string) (*AreaType, error) {
	return uc.repo.GetAreaType(ctx, id)
}

func (uc *MdmUsecase) CreateAreaType(ctx context.Context, areaType *AreaType) (*AreaType, error) {
	normalizeAreaType(areaType)
	return uc.repo.CreateAreaType(ctx, areaType)
}

func (uc *MdmUsecase) UpdateAreaType(ctx context.Context, areaType *AreaType) (*AreaType, error) {
	normalizeAreaType(areaType)
	return uc.repo.UpdateAreaType(ctx, areaType)
}

func (uc *MdmUsecase) DeleteAreaType(ctx context.Context, id string) error {
	return uc.repo.DeleteAreaType(ctx, id)
}

func (uc *MdmUsecase) ListAreas(ctx context.Context, filter AreaFilter) ([]*Area, string, error) {
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListAreas(ctx, filter)
}

func (uc *MdmUsecase) ListAreaTree(ctx context.Context, filter AreaFilter) ([]*AreaNode, error) {
	filter.PageSize = 1000
	filter.PageToken = ""
	areas, _, err := uc.repo.ListAreas(ctx, filter)
	if err != nil {
		return nil, err
	}
	return buildAreaTree(areas), nil
}

func (uc *MdmUsecase) GetArea(ctx context.Context, id string) (*Area, error) {
	return uc.repo.GetArea(ctx, id)
}

func (uc *MdmUsecase) CreateArea(ctx context.Context, area *Area) (*Area, error) {
	normalizeArea(area)
	return uc.repo.CreateArea(ctx, area)
}

func (uc *MdmUsecase) UpdateArea(ctx context.Context, area *Area) (*Area, error) {
	normalizeArea(area)
	return uc.repo.UpdateArea(ctx, area)
}

func (uc *MdmUsecase) DeleteArea(ctx context.Context, id string) error {
	return uc.repo.DeleteArea(ctx, id)
}

func (uc *MdmUsecase) AssignAreaAdministrativeUnit(ctx context.Context, item *AreaAdministrativeUnit) (*AreaAdministrativeUnit, error) {
	item.ScopeType = upperDefault(item.ScopeType, "INCLUDE")
	return uc.repo.AssignAreaAdministrativeUnit(ctx, item)
}

func (uc *MdmUsecase) ListAreaAdministrativeUnits(ctx context.Context, areaID string) ([]*AreaAdministrativeUnit, error) {
	return uc.repo.ListAreaAdministrativeUnits(ctx, areaID)
}

func (uc *MdmUsecase) RemoveAreaAdministrativeUnit(ctx context.Context, areaID, administrativeUnitID string) error {
	return uc.repo.RemoveAreaAdministrativeUnit(ctx, areaID, administrativeUnitID)
}

func (uc *MdmUsecase) ListCodeSets(ctx context.Context, filter PageFilter) ([]*CodeSet, string, error) {
	normalizePageFilter(&filter)
	return uc.repo.ListCodeSets(ctx, filter)
}

func (uc *MdmUsecase) GetCodeSet(ctx context.Context, id string) (*CodeSet, error) {
	return uc.repo.GetCodeSet(ctx, id)
}

func (uc *MdmUsecase) CreateCodeSet(ctx context.Context, codeSet *CodeSet) (*CodeSet, error) {
	normalizeCodeSet(codeSet)
	return uc.repo.CreateCodeSet(ctx, codeSet)
}

func (uc *MdmUsecase) UpdateCodeSet(ctx context.Context, codeSet *CodeSet) (*CodeSet, error) {
	normalizeCodeSet(codeSet)
	return uc.repo.UpdateCodeSet(ctx, codeSet)
}

func (uc *MdmUsecase) DeleteCodeSet(ctx context.Context, id string) error {
	return uc.repo.DeleteCodeSet(ctx, id)
}

func (uc *MdmUsecase) ListCodeItems(ctx context.Context, filter CodeItemFilter) ([]*CodeItem, string, error) {
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListCodeItems(ctx, filter)
}

func (uc *MdmUsecase) GetCodeItem(ctx context.Context, id string) (*CodeItem, error) {
	return uc.repo.GetCodeItem(ctx, id)
}

func (uc *MdmUsecase) CreateCodeItem(ctx context.Context, codeSetCode string, item *CodeItem) (*CodeItem, error) {
	codeSet, err := uc.repo.GetCodeSetByCode(ctx, codeSetCode)
	if err != nil {
		return nil, err
	}
	item.CodeSetID = codeSet.ID
	item.CodeSetCode = codeSet.Code
	normalizeCodeItem(item)
	return uc.repo.CreateCodeItem(ctx, item)
}

func (uc *MdmUsecase) UpdateCodeItem(ctx context.Context, item *CodeItem) (*CodeItem, error) {
	normalizeCodeItem(item)
	return uc.repo.UpdateCodeItem(ctx, item)
}

func (uc *MdmUsecase) DeleteCodeItem(ctx context.Context, id string) error {
	return uc.repo.DeleteCodeItem(ctx, id)
}

func (uc *MdmUsecase) ListSystemParameters(ctx context.Context, filter SystemParameterFilter) ([]*SystemParameter, string, error) {
	normalizePageFilter(&filter.PageFilter)
	return uc.repo.ListSystemParameters(ctx, filter)
}

func (uc *MdmUsecase) GetSystemParameter(ctx context.Context, key string) (*SystemParameter, error) {
	return uc.repo.GetSystemParameter(ctx, key)
}

func (uc *MdmUsecase) CreateSystemParameter(ctx context.Context, param *SystemParameter) (*SystemParameter, error) {
	normalizeSystemParameter(param)
	return uc.repo.CreateSystemParameter(ctx, param)
}

func (uc *MdmUsecase) UpdateSystemParameter(ctx context.Context, param *SystemParameter) (*SystemParameter, error) {
	existing, err := uc.repo.GetSystemParameter(ctx, param.Key)
	if err != nil {
		return nil, err
	}
	if !existing.IsEditable {
		return nil, ErrReadOnly
	}
	normalizeSystemParameter(param)
	return uc.repo.UpdateSystemParameter(ctx, param)
}

func (uc *MdmUsecase) DeleteSystemParameter(ctx context.Context, key string) error {
	existing, err := uc.repo.GetSystemParameter(ctx, key)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		return ErrReadOnly
	}
	return uc.repo.DeleteSystemParameter(ctx, key)
}

func normalizePageFilter(filter *PageFilter) {
	filter.Status = strings.ToUpper(strings.TrimSpace(filter.Status))
	filter.Keyword = strings.TrimSpace(filter.Keyword)
}

func normalizeAdministrativeUnit(unit *AdministrativeUnit) {
	unit.Code = strings.TrimSpace(unit.Code)
	unit.Name = strings.TrimSpace(unit.Name)
	unit.Level = upperDefault(unit.Level, "WARD")
	unit.UnitType = upperDefault(unit.UnitType, "XA")
	unit.Status = upperDefault(unit.Status, "ACTIVE")
	unit.MetadataJSON = jsonDefault(unit.MetadataJSON)
}

func normalizeAreaType(areaType *AreaType) {
	areaType.Code = upperDefault(areaType.Code, "")
	areaType.Name = strings.TrimSpace(areaType.Name)
	areaType.Status = upperDefault(areaType.Status, "ACTIVE")
}

func normalizeArea(area *Area) {
	area.Code = strings.TrimSpace(area.Code)
	area.Name = strings.TrimSpace(area.Name)
	area.Status = upperDefault(area.Status, "ACTIVE")
	area.MetadataJSON = jsonDefault(area.MetadataJSON)
}

func normalizeCodeSet(codeSet *CodeSet) {
	codeSet.Code = upperDefault(codeSet.Code, "")
	codeSet.Name = strings.TrimSpace(codeSet.Name)
	codeSet.Status = upperDefault(codeSet.Status, "ACTIVE")
}

func normalizeCodeItem(item *CodeItem) {
	item.Code = upperDefault(item.Code, "")
	item.Name = strings.TrimSpace(item.Name)
	item.Status = upperDefault(item.Status, "ACTIVE")
	item.MetadataJSON = jsonDefault(item.MetadataJSON)
}

func normalizeSystemParameter(param *SystemParameter) {
	param.Key = upperDefault(param.Key, "")
	param.GroupCode = upperDefault(param.GroupCode, "")
	param.ValueType = upperDefault(param.ValueType, "STRING")
	param.Status = upperDefault(param.Status, "ACTIVE")
	param.ValueJSON = jsonDefault(param.ValueJSON)
	param.ValidationRuleJSON = jsonDefault(param.ValidationRuleJSON)
}

func upperDefault(value, fallback string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
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

func buildAdministrativeUnitTree(units []*AdministrativeUnit) []*AdministrativeUnitNode {
	nodes := make(map[string]*AdministrativeUnitNode, len(units))
	var roots []*AdministrativeUnitNode
	for _, unit := range units {
		nodes[unit.ID] = &AdministrativeUnitNode{Unit: unit}
	}
	for _, unit := range units {
		node := nodes[unit.ID]
		if unit.ParentID != "" {
			if parent := nodes[unit.ParentID]; parent != nil {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}

func buildAreaTree(areas []*Area) []*AreaNode {
	nodes := make(map[string]*AreaNode, len(areas))
	var roots []*AreaNode
	for _, area := range areas {
		nodes[area.ID] = &AreaNode{Area: area}
	}
	for _, area := range areas {
		node := nodes[area.ID]
		if area.ParentID != "" {
			if parent := nodes[area.ParentID]; parent != nil {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots
}
