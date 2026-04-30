package biz

import (
	"context"
	"strings"
	"time"
)

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
