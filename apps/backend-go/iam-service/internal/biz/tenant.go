package biz

import (
	"context"
	"time"
)

const (
	TenantDeploymentShared    = "SHARED"
	TenantDeploymentDedicated = "DEDICATED"

	TenantAuthShared    = "SHARED_AUTH"
	TenantAuthDedicated = "DEDICATED_AUTH"
)

type Tenant struct {
	ID             string
	Name           string
	Slug           string
	OwnerID        string
	DeploymentMode string
	AuthMode       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TenantRepo interface {
	Create(ctx context.Context, tenant *Tenant) (*Tenant, error)
	GetByID(ctx context.Context, id string) (*Tenant, error)
	GetByIDs(ctx context.Context, ids []string) ([]*Tenant, error)
	ListAll(ctx context.Context) ([]*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
	Update(ctx context.Context, tenant *Tenant) (*Tenant, error)
	SoftDelete(ctx context.Context, id string) error
}

type TenantUsecase struct {
	repo        TenantRepo
	tenantUsers *TenantUserUsecase
}

func NewTenantUsecase(repo TenantRepo, tenantUsers *TenantUserUsecase) *TenantUsecase {
	return &TenantUsecase{repo: repo, tenantUsers: tenantUsers}
}

func (uc *TenantUsecase) GetTenantsByIDs(ctx context.Context, ids []string) ([]*Tenant, error) {
	return uc.repo.GetByIDs(ctx, ids)
}

func (uc *TenantUsecase) ListTenants(ctx context.Context) ([]*Tenant, error) {
	return uc.repo.ListAll(ctx)
}

func (uc *TenantUsecase) CreateTenant(ctx context.Context, name, slug, ownerID, deploymentMode, authMode string) (*Tenant, error) {
	t, err := uc.repo.Create(ctx, &Tenant{
		Name:           name,
		Slug:           slug,
		OwnerID:        ownerID,
		DeploymentMode: normalizeDeploymentMode(deploymentMode),
		AuthMode:       normalizeAuthMode(authMode),
	})
	if err != nil {
		return nil, err
	}
	_, _ = uc.tenantUsers.AddTenantUser(ctx, ownerID, t.ID, "", "", "owner")
	return t, nil
}

func (uc *TenantUsecase) GetTenant(ctx context.Context, id string) (*Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *TenantUsecase) UpdateTenant(ctx context.Context, id, name, slug, deploymentMode, authMode string) (*Tenant, error) {
	t, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	t.Name = name
	t.Slug = slug
	t.DeploymentMode = normalizeDeploymentMode(deploymentMode)
	t.AuthMode = normalizeAuthMode(authMode)
	return uc.repo.Update(ctx, t)
}

func (uc *TenantUsecase) DeleteTenant(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}

func normalizeDeploymentMode(mode string) string {
	if mode == TenantDeploymentDedicated {
		return TenantDeploymentDedicated
	}
	return TenantDeploymentShared
}

func normalizeAuthMode(mode string) string {
	if mode == TenantAuthDedicated {
		return TenantAuthDedicated
	}
	return TenantAuthShared
}
