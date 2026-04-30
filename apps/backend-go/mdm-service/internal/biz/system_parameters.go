package biz

import (
	"context"
	"time"
)

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

func normalizeSystemParameter(param *SystemParameter) {
	param.Key = upperDefault(param.Key, "")
	param.GroupCode = upperDefault(param.GroupCode, "")
	param.ValueType = upperDefault(param.ValueType, "STRING")
	param.Status = upperDefault(param.Status, "ACTIVE")
	param.ValueJSON = jsonDefault(param.ValueJSON)
	param.ValidationRuleJSON = jsonDefault(param.ValidationRuleJSON)
}
