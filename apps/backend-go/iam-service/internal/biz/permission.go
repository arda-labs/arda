package biz

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCannotDeleteSystemRole = errors.New("cannot delete system role")
)

type ResourcePermission struct {
	ID         string
	UserID     string
	TenantID   string
	Resource   string
	Action     string
	ResourceID string
	Allowed    bool
	CreatedAt  time.Time
}

type PermissionRepo interface {
	CheckByRole(ctx context.Context, userID, tenantID, resource, action string) (bool, error)
	GetResourceOverride(ctx context.Context, userID, tenantID, resource, action, resourceID string) (*ResourcePermission, error)
	GetResourcePermission(ctx context.Context, id string) (*ResourcePermission, error)
	GrantResourcePermission(ctx context.Context, rp *ResourcePermission) (*ResourcePermission, error)
	UpdateResourceStatus(ctx context.Context, id string, status string, checkerID string) error
	RevokeResourcePermission(ctx context.Context, id string) error
	ListByTenant(ctx context.Context, tenantID string, roleID string) ([]*Permission, error)
}

type CachedPermission struct {
	Allowed bool
	Source  string
}

type PermissionCache interface {
	Get(ctx context.Context, userID, tenantID, resource, action, resourceID string) (*CachedPermission, bool)
	Set(ctx context.Context, userID, tenantID, resource, action, resourceID string, allowed bool, source string)
	InvalidateUser(ctx context.Context, userID, tenantID string)
}

type PermissionUsecase struct {
	roleRepo  RoleRepo
	permRepo  PermissionRepo
	cache     PermissionCache
	audit     *AuditUsecase
}

func NewPermissionUsecase(roleRepo RoleRepo, permRepo PermissionRepo, cache PermissionCache, audit *AuditUsecase) *PermissionUsecase {
	return &PermissionUsecase{roleRepo: roleRepo, permRepo: permRepo, cache: cache, audit: audit}
}

func (uc *PermissionUsecase) CheckPermission(ctx context.Context, userID, tenantID, resource, action, resourceID string) (bool, string, error) {
	cached, ok := uc.cache.Get(ctx, userID, tenantID, resource, action, resourceID)
	if ok {
		return cached.Allowed, cached.Source, nil
	}

	if resourceID != "" {
		override, err := uc.permRepo.GetResourceOverride(ctx, userID, tenantID, resource, action, resourceID)
		if err != nil {
			return false, "", fmt.Errorf("checking resource override: %w", err)
		}
		if override != nil {
			uc.cache.Set(ctx, userID, tenantID, resource, action, resourceID, override.Allowed, "resource_override")
			return override.Allowed, "resource_override", nil
		}
	}

	allowed, err := uc.permRepo.CheckByRole(ctx, userID, tenantID, resource, action)
	if err != nil {
		return false, "", fmt.Errorf("checking role permission: %w", err)
	}
	if allowed {
		uc.cache.Set(ctx, userID, tenantID, resource, action, resourceID, true, "role")
		return true, "role", nil
	}

	uc.cache.Set(ctx, userID, tenantID, resource, action, resourceID, false, "denied")
	return false, "denied", nil
}

func (uc *PermissionUsecase) ListPermissions(ctx context.Context, tenantID, roleID string) ([]*Permission, error) {
	return uc.permRepo.ListByTenant(ctx, tenantID, roleID)
}

func (uc *PermissionUsecase) GrantResourcePermission(ctx context.Context, rp *ResourcePermission, actorID string) (*ResourcePermission, error) {
	result, err := uc.permRepo.GrantResourcePermission(ctx, rp)
	if err != nil {
		return nil, err
	}
	uc.cache.InvalidateUser(ctx, rp.UserID, rp.TenantID)
	uc.audit.Log(ctx, actorID, rp.TenantID, "permission.granted", "resource_permission", rp.UserID, nil)
	return result, nil
}

func (uc *PermissionUsecase) ApprovePermission(ctx context.Context, id, checkerID string) error {
	rp, err := uc.permRepo.GetResourcePermission(ctx, id)
	if err != nil {
		return err
	}
	if rp == nil {
		return errors.New("permission request not found")
	}
	if rp.UserID == checkerID {
		return errors.New("maker cannot be checker (maker-checker violation)")
	}

	err = uc.permRepo.UpdateResourceStatus(ctx, id, "active", checkerID)
	if err != nil {
		return err
	}

	uc.cache.InvalidateUser(ctx, rp.UserID, rp.TenantID)
	return nil
}

func (uc *PermissionUsecase) RejectPermission(ctx context.Context, id, checkerID string) error {
	return uc.permRepo.UpdateResourceStatus(ctx, id, "rejected", checkerID)
}

func (uc *PermissionUsecase) GetResourcePermission(ctx context.Context, id string) (*ResourcePermission, error) {
	return uc.permRepo.GetResourcePermission(ctx, id)
}

func (uc *PermissionUsecase) RevokeResourcePermission(ctx context.Context, id, actorID string) error {
	return uc.permRepo.RevokeResourcePermission(ctx, id)
}

func (uc *PermissionUsecase) GetUserPermissions(ctx context.Context, userID, tenantID string) ([]*Permission, error) {
	roles, err := uc.roleRepo.GetUserRoles(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var perms []*Permission
	for _, role := range roles {
		p, err := uc.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			continue
		}
		for _, perm := range p {
			key := perm.Resource + ":" + perm.Action
			if !seen[key] {
				seen[key] = true
				perms = append(perms, perm)
			}
		}
	}
	return perms, nil
}
