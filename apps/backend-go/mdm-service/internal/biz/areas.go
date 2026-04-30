package biz

import (
	"context"
	"strings"
	"time"
)

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
