package biz

import (
	"context"
	"time"
)

type Role struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	IsSystem    bool
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID       string
	TenantID string
	Resource string
	Action   string
}

type RoleRepo interface {
	Create(ctx context.Context, role *Role) (*Role, error)
	GetByID(ctx context.Context, id string) (*Role, error)
	ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Role, string, error)
	Update(ctx context.Context, role *Role) (*Role, error)
	SoftDelete(ctx context.Context, id string) error
	AssignRole(ctx context.Context, userID, roleID, tenantID string) error
	RevokeRole(ctx context.Context, userID, roleID, tenantID string) error
	GetUserRoles(ctx context.Context, userID, tenantID string) ([]*Role, error)
	GetGroupRoles(ctx context.Context, userID, tenantID string) ([]*Role, error)
	GetRolePermissions(ctx context.Context, roleID string) ([]*Permission, error)
	SetRolePermissions(ctx context.Context, roleID string, permIDs []string) error
}

type RoleUsecase struct {
	repo  RoleRepo
	audit *AuditUsecase
	cache PermissionCache
}

func NewRoleUsecase(repo RoleRepo, audit *AuditUsecase, cache PermissionCache) *RoleUsecase {
	return &RoleUsecase{repo: repo, audit: audit, cache: cache}
}

func (uc *RoleUsecase) CreateRole(ctx context.Context, tenantID, name, desc string, permIDs []string) (*Role, error) {
	role, err := uc.repo.Create(ctx, &Role{TenantID: tenantID, Name: name, Description: desc})
	if err != nil {
		return nil, err
	}
	if len(permIDs) > 0 {
		if err := uc.repo.SetRolePermissions(ctx, role.ID, permIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (uc *RoleUsecase) GetRole(ctx context.Context, id string) (*Role, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *RoleUsecase) ListRoles(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Role, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, id, name, desc string, permIDs []string) (*Role, error) {
	role, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Name = name
	role.Description = desc
	role, err = uc.repo.Update(ctx, role)
	if err != nil {
		return nil, err
	}
	if err := uc.repo.SetRolePermissions(ctx, id, permIDs); err != nil {
		return nil, err
	}
	return role, nil
}

func (uc *RoleUsecase) DeleteRole(ctx context.Context, id string) error {
	role, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return ErrCannotDeleteSystemRole
	}
	return uc.repo.SoftDelete(ctx, id)
}

func (uc *RoleUsecase) AssignRole(ctx context.Context, userID, roleID, tenantID string, actorID string) error {
	if err := uc.repo.AssignRole(ctx, userID, roleID, tenantID); err != nil {
		return err
	}
	uc.cache.InvalidateUser(ctx, userID, tenantID)
	uc.audit.Log(ctx, actorID, tenantID, "role.assigned", "user_role", userID, nil)
	return nil
}

func (uc *RoleUsecase) RevokeRole(ctx context.Context, userID, roleID, tenantID string, actorID string) error {
	if err := uc.repo.RevokeRole(ctx, userID, roleID, tenantID); err != nil {
		return err
	}
	uc.cache.InvalidateUser(ctx, userID, tenantID)
	uc.audit.Log(ctx, actorID, tenantID, "role.revoked", "user_role", userID, nil)
	return nil
}
