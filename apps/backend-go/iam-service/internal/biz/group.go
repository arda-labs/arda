package biz

import (
	"context"
	"time"
)

type Group struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GroupRepo interface {
	Create(ctx context.Context, g *Group) (*Group, error)
	GetByID(ctx context.Context, id string) (*Group, error)
	Update(ctx context.Context, g *Group) (*Group, error)
	Delete(ctx context.Context, id string) error
	ListByTenant(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Group, string, error)

	AddMember(ctx context.Context, groupID, userID string) error
	RemoveMember(ctx context.Context, groupID, userID string) error
	ListMembers(ctx context.Context, groupID string) ([]*User, error)

	AssignRole(ctx context.Context, groupID, roleID, tenantID string) error
	RevokeRole(ctx context.Context, groupID, roleID string) error
	ListRoles(ctx context.Context, groupID string) ([]*Role, error)
}

type GroupUsecase struct {
	repo  GroupRepo
	perms PermissionCache
	audit *AuditUsecase
}

func NewGroupUsecase(repo GroupRepo, perms PermissionCache, audit *AuditUsecase) *GroupUsecase {
	return &GroupUsecase{repo: repo, perms: perms, audit: audit}
}

func (uc *GroupUsecase) CreateGroup(ctx context.Context, tenantID, name, desc string) (*Group, error) {
	return uc.repo.Create(ctx, &Group{
		TenantID:    tenantID,
		Name:        name,
		Description: desc,
	})
}

func (uc *GroupUsecase) GetGroup(ctx context.Context, id string) (*Group, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *GroupUsecase) UpdateGroup(ctx context.Context, id, name, desc string) (*Group, error) {
	g, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	g.Name = name
	g.Description = desc
	return uc.repo.Update(ctx, g)
}

func (uc *GroupUsecase) DeleteGroup(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *GroupUsecase) ListGroups(ctx context.Context, tenantID string, pageSize int, cursor string) ([]*Group, string, error) {
	return uc.repo.ListByTenant(ctx, tenantID, pageSize, cursor)
}

func (uc *GroupUsecase) AddMember(ctx context.Context, groupID, userID string, actorID string) error {
	if err := uc.repo.AddMember(ctx, groupID, userID); err != nil {
		return err
	}
	g, _ := uc.repo.GetByID(ctx, groupID)
	if g != nil {
		uc.perms.InvalidateUser(ctx, userID, g.TenantID)
		uc.audit.Log(ctx, actorID, g.TenantID, "group.member_added", "group", groupID, map[string]any{"user_id": userID})
	}
	return nil
}

func (uc *GroupUsecase) RemoveMember(ctx context.Context, groupID, userID string, actorID string) error {
	if err := uc.repo.RemoveMember(ctx, groupID, userID); err != nil {
		return err
	}
	g, _ := uc.repo.GetByID(ctx, groupID)
	if g != nil {
		uc.perms.InvalidateUser(ctx, userID, g.TenantID)
		uc.audit.Log(ctx, actorID, g.TenantID, "group.member_removed", "group", groupID, map[string]any{"user_id": userID})
	}
	return nil
}

func (uc *GroupUsecase) ListMembers(ctx context.Context, groupID string) ([]*User, error) {
	return uc.repo.ListMembers(ctx, groupID)
}

func (uc *GroupUsecase) AssignRole(ctx context.Context, groupID, roleID string, actorID string) error {
	g, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if err := uc.repo.AssignRole(ctx, groupID, roleID, g.TenantID); err != nil {
		return err
	}

	// Invalidate cache for all members of the group
	members, _ := uc.repo.ListMembers(ctx, groupID)
	for _, m := range members {
		uc.perms.InvalidateUser(ctx, m.ID, g.TenantID)
	}

	uc.audit.Log(ctx, actorID, g.TenantID, "group.role_assigned", "group", groupID, map[string]any{"role_id": roleID})
	return nil
}

func (uc *GroupUsecase) RevokeRole(ctx context.Context, groupID, roleID string, actorID string) error {
	g, err := uc.repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if err := uc.repo.RevokeRole(ctx, groupID, roleID); err != nil {
		return err
	}

	// Invalidate cache for all members of the group
	members, _ := uc.repo.ListMembers(ctx, groupID)
	for _, m := range members {
		uc.perms.InvalidateUser(ctx, m.ID, g.TenantID)
	}

	uc.audit.Log(ctx, actorID, g.TenantID, "group.role_revoked", "group", groupID, map[string]any{"role_id": roleID})
	return nil
}

func (uc *GroupUsecase) ListRoles(ctx context.Context, groupID string) ([]*Role, error) {
	return uc.repo.ListRoles(ctx, groupID)
}
