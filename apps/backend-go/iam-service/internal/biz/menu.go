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
	SortOrder      int
	Enabled        bool
	PermissionSlug string
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
	permUC   *PermissionUsecase
	log      *log.Helper
}

func NewMenuUsecase(repo MenuRepo, permUC *PermissionUsecase, logger log.Logger) *MenuUsecase {
	return &MenuUsecase{
		menuRepo: repo,
		permUC:   permUC,
		log:      log.NewHelper(logger),
	}
}

// GetUserMenu returns the full menu tree for a user (filtered by their permissions).
func (uc *MenuUsecase) GetUserMenu(ctx context.Context, userID, tenantID string) ([]*Menu, error) {
	// 1. Get all enabled menus for the tenant
	allMenus, err := uc.menuRepo.GetByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// 2. Get user permissions
	userActivePerms, err := uc.permUC.GetUserPermissions(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	permMap := make(map[string]bool)
	for _, p := range userActivePerms {
		permMap[p.Resource+":"+p.Action] = true
	}

	// 3. Filter menus
	var filtered []*Menu
	for _, m := range allMenus {
		if !m.Enabled {
			continue
		}
		// If no permission required, allow it
		if m.PermissionSlug == "" {
			filtered = append(filtered, m)
			continue
		}
		// Check if user has the required permission
		if permMap[m.PermissionSlug] {
			filtered = append(filtered, m)
		}
	}

	return filtered, nil
}

// ListMenus returns flat list of menus for a tenant (admin view).
func (uc *MenuUsecase) ListMenus(ctx context.Context, tenantID string) ([]*Menu, error) {
	return uc.menuRepo.GetByTenant(ctx, tenantID)
}

func (uc *MenuUsecase) GetMenu(ctx context.Context, id string) (*Menu, error) {
	return uc.menuRepo.GetByID(ctx, id)
}

func (uc *MenuUsecase) CreateMenu(ctx context.Context, tenantID, parentID, name, slug, icon, route string, sortOrder int, enabled bool, permSlug string) (*Menu, error) {
	m := &Menu{
		ID:             "",
		TenantID:       tenantID,
		ParentID:       parentID,
		Name:           name,
		Slug:           slug,
		Icon:           icon,
		Route:          route,
		SortOrder:      sortOrder,
		Enabled:        enabled,
		PermissionSlug: permSlug,
	}
	return uc.menuRepo.Create(ctx, m)
}

func (uc *MenuUsecase) UpdateMenu(ctx context.Context, id, parentID, name, slug, icon, route string, sortOrder int, enabled bool, permSlug string) (*Menu, error) {
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
	existing.PermissionSlug = permSlug
	return uc.menuRepo.Update(ctx, existing)
}

func (uc *MenuUsecase) DeleteMenu(ctx context.Context, id string) error {
	return uc.menuRepo.Delete(ctx, id)
}
