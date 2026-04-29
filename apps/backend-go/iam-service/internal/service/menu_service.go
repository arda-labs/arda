package service

import (
	"context"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type MenuService struct {
	menuUsecase *biz.MenuUsecase
	log         *log.Helper
}

func NewMenuService(menuUsecase *biz.MenuUsecase, logger log.Logger) *MenuService {
	return &MenuService{
		menuUsecase: menuUsecase,
		log:         log.NewHelper(logger),
	}
}

// ─── HTTP request/response types (plain Go, no proto dependency) ─────────────────

type GetMenuRequest struct{}

type GetMenuResponse struct {
	Items []MenuItem `json:"items"`
}

type MenuItem struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Icon           string     `json:"icon"`
	Route          string     `json:"route"`
	SortOrder      int        `json:"sort_order"`
	PermissionSlug string     `json:"permission_slug"`
	Children       []MenuItem `json:"children"`
}

type Menu struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenant_id"`
	ParentID       string `json:"parent_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Icon           string `json:"icon"`
	Route          string `json:"route"`
	SortOrder      int    `json:"sort_order"`
	Enabled        bool   `json:"enabled"`
	PermissionSlug string `json:"permission_slug"`
}

type ListMenusRequest struct {
	TenantID string `json:"tenant_id"`
}

type ListMenusResponse struct {
	Menus []Menu `json:"menus"`
}

type CreateMenuRequest struct {
	TenantID       string `json:"tenant_id"`
	ParentID       string `json:"parent_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Icon           string `json:"icon"`
	Route          string `json:"route"`
	SortOrder      int    `json:"sort_order"`
	Enabled        bool   `json:"enabled"`
	PermissionSlug string `json:"permission_slug"`
}

type UpdateMenuRequest struct {
	ID             string `json:"id"`
	ParentID       string `json:"parent_id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Icon           string `json:"icon"`
	Route          string `json:"route"`
	SortOrder      int    `json:"sort_order"`
	Enabled        bool   `json:"enabled"`
	PermissionSlug string `json:"permission_slug"`
}

type DeleteMenuRequest struct {
	ID string `json:"id"`
}

type DeleteMenuResponse struct{}

// ─── Service methods ─────────────────────────────────────────────────────────

func (s *MenuService) GetMenu(ctx context.Context, req *GetMenuRequest) (*GetMenuResponse, error) {
	tenantID := middleware.GetTenantID(ctx)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, errors.Forbidden("UNAUTHORIZED", "missing subject")
	}
	if tenantID == "" {
		return nil, errors.Forbidden("MISSING_TENANT", "tenant context is required")
	}

	menus, err := s.menuUsecase.GetUserMenu(ctx, userID, tenantID)
	if err != nil {
		return nil, err
	}

	tree := buildMenuTree(menus)
	return &GetMenuResponse{Items: tree}, nil
}

func (s *MenuService) ListMenus(ctx context.Context, req *ListMenusRequest) (*ListMenusResponse, error) {
	if req.TenantID == "" {
		req.TenantID = middleware.GetTenantID(ctx)
	}
	if req.TenantID == "" {
		return nil, errors.Forbidden("MISSING_TENANT", "tenant context is required")
	}
	if err := s.requireMenuPermission(ctx, req.TenantID, "read"); err != nil {
		return nil, err
	}

	menus, err := s.menuUsecase.ListMenus(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}

	var out []Menu
	for _, m := range menus {
		out = append(out, toMenu(m))
	}
	return &ListMenusResponse{Menus: out}, nil
}

func (s *MenuService) CreateMenu(ctx context.Context, req *CreateMenuRequest) (*Menu, error) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" && req.TenantID != "" {
		tenantID = req.TenantID
	}
	if tenantID == "" {
		return nil, errors.Forbidden("MISSING_TENANT", "tenant context is required")
	}
	if err := s.requireMenuPermission(ctx, tenantID, "create"); err != nil {
		return nil, err
	}

	m, err := s.menuUsecase.CreateMenu(ctx, tenantID, req.ParentID, req.Name, req.Slug, req.Icon, req.Route, req.SortOrder, req.Enabled, req.PermissionSlug)
	if err != nil {
		return nil, err
	}
	out := toMenu(m)
	return &out, nil
}

func (s *MenuService) UpdateMenu(ctx context.Context, req *UpdateMenuRequest) (*Menu, error) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.Forbidden("MISSING_TENANT", "tenant context is required")
	}
	if err := s.requireMenuPermission(ctx, tenantID, "update"); err != nil {
		return nil, err
	}

	m, err := s.menuUsecase.UpdateMenu(ctx, req.ID, req.ParentID, req.Name, req.Slug, req.Icon, req.Route, req.SortOrder, req.Enabled, req.PermissionSlug)
	if err != nil {
		return nil, err
	}
	out := toMenu(m)
	return &out, nil
}

func (s *MenuService) DeleteMenu(ctx context.Context, req *DeleteMenuRequest) (*DeleteMenuResponse, error) {
	tenantID := middleware.GetTenantID(ctx)
	if tenantID == "" {
		return nil, errors.Forbidden("MISSING_TENANT", "tenant context is required")
	}
	if err := s.requireMenuPermission(ctx, tenantID, "delete"); err != nil {
		return nil, err
	}

	if err := s.menuUsecase.DeleteMenu(ctx, req.ID); err != nil {
		return nil, err
	}
	return &DeleteMenuResponse{}, nil
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func (s *MenuService) requireMenuPermission(ctx context.Context, tenantID, action string) error {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return errors.Forbidden("UNAUTHORIZED", "missing subject")
	}

	allowed, err := s.menuUsecase.CheckMenuPermission(ctx, userID, tenantID, action)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.Forbidden("PERMISSION_DENIED", "menu:"+action+" permission is required")
	}
	return nil
}

func buildMenuTree(menus []*biz.Menu) []MenuItem {
	menuMap := make(map[string]*biz.Menu)
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	children := make(map[string][]*biz.Menu)
	var roots []*biz.Menu
	for _, m := range menus {
		if m.ParentID != "" && menuMap[m.ParentID] != nil {
			children[m.ParentID] = append(children[m.ParentID], m)
			continue
		}
		roots = append(roots, m)
	}

	var toItem func(*biz.Menu) MenuItem
	toItem = func(m *biz.Menu) MenuItem {
		item := MenuItem{
			ID:             m.ID,
			Name:           m.Name,
			Icon:           m.Icon,
			Route:          m.Route,
			SortOrder:      m.SortOrder,
			PermissionSlug: m.PermissionSlug,
			Children:       []MenuItem{},
		}
		for _, child := range children[m.ID] {
			item.Children = append(item.Children, toItem(child))
		}
		return item
	}

	out := make([]MenuItem, 0, len(roots))
	for _, root := range roots {
		out = append(out, toItem(root))
	}

	return out
}

func toMenu(m *biz.Menu) Menu {
	return Menu{
		ID:             m.ID,
		TenantID:       m.TenantID,
		ParentID:       m.ParentID,
		Name:           m.Name,
		Slug:           m.Slug,
		Icon:           m.Icon,
		Route:          m.Route,
		SortOrder:      m.SortOrder,
		Enabled:        m.Enabled,
		PermissionSlug: m.PermissionSlug,
	}
}
