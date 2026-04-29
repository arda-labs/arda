package biz

import (
	"context"
	"time"
)

type TenantUser struct {
	ID          string
	UserID      string
	TenantID    string
	Username    string
	Email       string
	DisplayName string
	Role        string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TenantUserRepo interface {
	Create(ctx context.Context, tenantUser *TenantUser) (*TenantUser, error)
	ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*TenantUser, string, error)
	ListByUser(ctx context.Context, userID string) ([]*TenantUser, error)
	GetByUserAndTenant(ctx context.Context, userID, tenantID string) (*TenantUser, error)
	SoftDelete(ctx context.Context, userID, tenantID string) error
}

type TenantUserUsecase struct {
	repo TenantUserRepo
}

func NewTenantUserUsecase(repo TenantUserRepo) *TenantUserUsecase {
	return &TenantUserUsecase{repo: repo}
}

func (uc *TenantUserUsecase) ListByUser(ctx context.Context, userID string) ([]*TenantUser, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *TenantUserUsecase) GetTenantUser(ctx context.Context, userID, tenantID string) (*TenantUser, error) {
	return uc.repo.GetByUserAndTenant(ctx, userID, tenantID)
}

func (uc *TenantUserUsecase) AddTenantUser(ctx context.Context, userID, tenantID, username, displayName, role string) (*TenantUser, error) {
	return uc.repo.Create(ctx, &TenantUser{
		UserID:      userID,
		TenantID:    tenantID,
		Username:    username,
		DisplayName: displayName,
		Role:        role,
		Status:      "ACTIVE",
	})
}

func (uc *TenantUserUsecase) InviteTenantUser(ctx context.Context, tenantID, userID, username, displayName, role string) (*TenantUser, error) {
	existing, err := uc.repo.GetByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return uc.AddTenantUser(ctx, userID, tenantID, username, displayName, role)
}

func (uc *TenantUserUsecase) ListTenantUsers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*TenantUser, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}

func (uc *TenantUserUsecase) RemoveTenantUser(ctx context.Context, userID, tenantID string) error {
	return uc.repo.SoftDelete(ctx, userID, tenantID)
}
