package biz

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrMenuNotFound = errors.New("menu not found")
)

type Menu struct {
	ID        string
	TenantID  string
	ParentID  string
	Name      string
	Slug      string
	Icon      string
	Route     string
	SortOrder int
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MenuRepo interface {
	GetByTenant(ctx context.Context, tenantID string) ([]*Menu, error)
	GetByID(ctx context.Context, id string) (*Menu, error)
	Create(ctx context.Context, m *Menu) (*Menu, error)
	Update(ctx context.Context, m *Menu) (*Menu, error)
	Delete(ctx context.Context, id string) error
}

type MenuUsecase struct {
	menuRepo MenuRepo
	log      *log.Helper
}

func NewMenuUsecase(repo MenuRepo, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{menuRepo: repo, log: log.NewHelper(logger)}
}

// GetUserMenu returns the full menu tree for a user (filtered by their permissions).
// Currently returns all enabled menus for the tenant; permission filtering is done on the frontend.
func (uc *MenuUsecase) GetUserMenu(ctx context.Context, userID, tenantID string) ([]*Menu, error) {
	return uc.menuRepo.GetByTenant(ctx, tenantID)
}

// ListMenus returns flat list of menus for a tenant (admin view).
func (uc *MenuUsecase) ListMenus(ctx context.Context, tenantID string) ([]*Menu, error) {
	return uc.menuRepo.GetByTenant(ctx, tenantID)
}

func (uc *MenuUsecase) GetMenu(ctx context.Context, id string) (*Menu, error) {
	return uc.menuRepo.GetByID(ctx, id)
}

func (uc *MenuUsecase) CreateMenu(ctx context.Context, tenantID, parentID, name, slug, icon, route string, sortOrder int, enabled bool) (*Menu, error) {
	m := &Menu{
		ID:        "",
		TenantID:  tenantID,
		ParentID:  parentID,
		Name:      name,
		Slug:      slug,
		Icon:      icon,
		Route:     route,
		SortOrder: sortOrder,
		Enabled:   enabled,
	}
	return uc.menuRepo.Create(ctx, m)
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, id, parentID, name, slug, icon, route string, sortOrder int, enabled bool) (*Menu, error) {
	existing, err := uc.menuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existing.ParentID = parentID
	existing.Name = name
	existing.Slug = slug
	existing.Icon = icon
	existing.Route = route
	existing.SortOrder = sortOrder
	existing.Enabled = enabled
	return uc.menuRepo.Update(ctx, existing)
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id string) error {
	return uc.menuRepo.Delete(ctx, id)
}
