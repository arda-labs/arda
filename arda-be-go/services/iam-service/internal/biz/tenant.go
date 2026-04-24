package biz

import (
	"context"
	"time"
)

type Tenant struct {
	ID        string
	Name      string
	Slug      string
	OwnerID   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) (*Tenant, error)
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetByIDs(ctx context.Context, ids []string) ([]*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) (*Tenant, error)
	SoftDelete(ctx context.Context, id string) error
}

type TenantUsecase struct {
	repo    TenantRepo
	members *MembershipUsecase
}

func NewTenantUsecase(repo TenantRepo, members *MembershipUsecase) *TenantUsecase {
	return &TenantUsecase{repo: repo, members: members}
}

func (uc *TenantUsecase) GetTenantsByIDs(ctx context.Context, ids []string) ([]*Tenant, error) {
	return uc.repo.GetByIDs(ctx, ids)
}

func (uc *TenantUsecase) CreateTenant(ctx context.Context, name, slug, ownerID string) (*Tenant, error) {
	t, err := uc.repo.Create(ctx, &Tenant{Name: name, Slug: slug, OwnerID: ownerID})
	if err != nil {
		return nil, err
	}
	_, _ = uc.members.AddMembership(ctx, ownerID, t.ID, "owner")
	return t, nil
}

func (uc *TenantUsecase) GetTenant(ctx context.Context, id string) (*Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TenantUsecase) UpdateTenant(ctx context.Context, id, name, slug string) (*Tenant, error) {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Name = name
	t.Slug = slug
	return uc.repo.Update(ctx, t)
}

func (uc *TenantUsecase) DeleteTenant(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
