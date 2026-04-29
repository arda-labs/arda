package biz

import (
	"context"
	"strings"
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
	repo     TenantUserRepo
	roleRepo RoleRepo
	cache    PermissionCache
}

func NewTenantUserUsecase(repo TenantUserRepo, roleRepo RoleRepo, cache PermissionCache) *TenantUserUsecase {
	return &TenantUserUsecase{repo: repo, roleRepo: roleRepo, cache: cache}
}

func (uc *TenantUserUsecase) ListByUser(ctx context.Context, userID string) ([]*TenantUser, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *TenantUserUsecase) GetTenantUser(ctx context.Context, userID, tenantID string) (*TenantUser, error) {
	return uc.repo.GetByUserAndTenant(ctx, userID, tenantID)
}

func (uc *TenantUserUsecase) AddTenantUser(ctx context.Context, userID, tenantID, username, displayName, role string) (*TenantUser, error) {
	tenantUser, err := uc.repo.Create(ctx, &TenantUser{
		UserID:      userID,
		TenantID:    tenantID,
		Username:    username,
		DisplayName: displayName,
		Role:        role,
		Status:      "ACTIVE",
	})
	if err != nil {
		return nil, err
	}
	if err := uc.assignTenantRole(ctx, userID, tenantID, role); err != nil {
		return nil, err
	}
	return tenantUser, nil
}

func (uc *TenantUserUsecase) InviteTenantUser(ctx context.Context, tenantID, userID, username, displayName, role string) (*TenantUser, error) {
	existing, err := uc.repo.GetByUserAndTenant(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		if err := uc.assignTenantRole(ctx, userID, tenantID, role); err != nil {
			return nil, err
		}
		return existing, nil
	}
	return uc.AddTenantUser(ctx, userID, tenantID, username, displayName, role)
}

func (uc *TenantUserUsecase) ListTenantUsers(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*TenantUser, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}

func (uc *TenantUserUsecase) RemoveTenantUser(ctx context.Context, userID, tenantID string) error {
	if err := uc.repo.SoftDelete(ctx, userID, tenantID); err != nil {
		return err
	}
	if uc.cache != nil {
		uc.cache.InvalidateUser(ctx, userID, tenantID)
	}
	return nil
}

func (uc *TenantUserUsecase) assignTenantRole(ctx context.Context, userID, tenantID, roleName string) error {
	roleName = strings.ToLower(strings.TrimSpace(roleName))
	if roleName == "" {
		roleName = "member"
	}

	role, err := uc.roleRepo.GetByName(ctx, tenantID, roleName)
	if err != nil {
		return err
	}
	if role == nil && roleName != "member" {
		role, err = uc.roleRepo.GetByName(ctx, tenantID, "member")
		if err != nil {
			return err
		}
	}
	if role == nil {
		return nil
	}
	if err := uc.roleRepo.AssignRole(ctx, userID, role.ID, tenantID); err != nil {
		return err
	}
	if uc.cache != nil {
		uc.cache.InvalidateUser(ctx, userID, tenantID)
	}
	return nil
}
