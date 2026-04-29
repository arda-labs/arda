package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/arda-labs/arda/arda-be-go/pkg/middleware"
	"github.com/arda-labs/arda/arda-be-go/services/iam-service/internal/biz"
	"github.com/jackc/pgx/v5"

	"github.com/redis/go-redis/v9"
)

type permissionRepo struct {
	data *Data
}

func NewPermissionRepo(data *Data) biz.PermissionRepo {
	return &permissionRepo{data: data}
}

func (r *permissionRepo) ResolveUserID(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", nil
	}

	var id string
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT u.id::text
			 FROM users u
			 LEFT JOIN tenant_users tu
			   ON tu.user_id = u.id
			  AND tu.tenant_id = $2
			  AND tu.deleted_at IS NULL
			 WHERE (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1) OR lower(tu.username) = lower($1))
			   AND u.deleted_at IS NULL
			 LIMIT 1`, userID, tenantID,
		).Scan(&id)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	return id, err
}

func (r *permissionRepo) IsPlatformAdmin(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, nil
	}

	var exists bool
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`SELECT EXISTS(
				SELECT 1
				FROM platform_admins pa
				JOIN users u ON u.id = pa.user_id
				WHERE pa.revoked_at IS NULL
				  AND u.deleted_at IS NULL
				  AND (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1))
			)`, userID,
		).Scan(&exists)
	})
	return exists, err
}

func (r *permissionRepo) CheckByRole(ctx context.Context, userID, tenantID, resource, action string) (bool, error) {
	var exists bool
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`WITH RECURSIVE actor AS (
				SELECT u.id
				FROM users u
				LEFT JOIN tenant_users tu
				  ON tu.user_id = u.id
				 AND tu.tenant_id = $2
				 AND tu.deleted_at IS NULL
				WHERE (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1) OR lower(tu.username) = lower($1))
				  AND u.deleted_at IS NULL
				LIMIT 1
			), role_tree AS (
				SELECT role_id FROM user_roles WHERE user_id = (SELECT id FROM actor) AND tenant_id = $2
				UNION
				SELECT rh.parent_role_id FROM role_hierarchy rh
				JOIN role_tree rt ON rh.child_role_id = rt.role_id
			)
			SELECT EXISTS(
				SELECT 1 FROM role_tree rt
				JOIN role_permissions rp ON rt.role_id = rp.role_id
				JOIN permissions p ON rp.permission_id = p.id
				WHERE p.resource = $3 AND p.action = $4
			)`, userID, tenantID, resource, action,
		).Scan(&exists)
	})
	return exists, err
}

func (r *permissionRepo) GetResourceOverride(ctx context.Context, userID, tenantID, resource, action, resourceID string) (*biz.ResourcePermission, error) {
	rp := &biz.ResourcePermission{}
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT rp.id, rp.user_id, rp.tenant_id, rp.resource, rp.action, rp.resource_id, rp.allowed, rp.created_at
			 FROM resource_permissions rp
			 JOIN users u ON u.id = rp.user_id
			 LEFT JOIN tenant_users tu
			   ON tu.user_id = u.id
			  AND tu.tenant_id = $2
			  AND tu.deleted_at IS NULL
			 WHERE (u.id::text = $1 OR u.external_id = $1 OR lower(u.email) = lower($1) OR lower(tu.username) = lower($1))
			   AND rp.tenant_id = $2
			   AND rp.resource = $3
			   AND rp.action = $4
			   AND rp.resource_id = $5
			   AND u.deleted_at IS NULL`,
			userID, tenantID, resource, action, resourceID,
		).Scan(&rp.ID, &rp.UserID, &rp.TenantID, &rp.Resource, &rp.Action, &rp.ResourceID, &rp.Allowed, &rp.CreatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if rp.ID == "" {
		return nil, nil
	}
	return rp, err
}

func (r *permissionRepo) GrantResourcePermission(ctx context.Context, rp *biz.ResourcePermission) (*biz.ResourcePermission, error) {
	err := r.data.ExecInTenant(ctx, rp.TenantID, func(ctx context.Context, tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO resource_permissions (user_id, tenant_id, resource, action, resource_id, allowed)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (user_id, tenant_id, resource, action, resource_id)
			 DO UPDATE SET allowed = EXCLUDED.allowed
			 RETURNING id, created_at`,
			rp.UserID, rp.TenantID, rp.Resource, rp.Action, rp.ResourceID, rp.Allowed,
		).Scan(&rp.ID, &rp.CreatedAt)
	})
	return rp, err
}

func (r *permissionRepo) RevokeResourcePermission(ctx context.Context, id string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `DELETE FROM resource_permissions WHERE id = $1`, id)
		return err
	})
}

func (r *permissionRepo) GetResourcePermission(ctx context.Context, id string) (*biz.ResourcePermission, error) {
	rp := &biz.ResourcePermission{}
	tenantID := middleware.GetTenantID(ctx)
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		err := tx.QueryRow(ctx,
			`SELECT id, user_id, tenant_id, resource, action, resource_id, allowed, created_at FROM resource_permissions WHERE id = $1`, id,
		).Scan(&rp.ID, &rp.UserID, &rp.TenantID, &rp.Resource, &rp.Action, &rp.ResourceID, &rp.Allowed, &rp.CreatedAt)
		if err == pgx.ErrNoRows {
			return nil
		}
		return err
	})
	if rp.ID == "" {
		return nil, nil
	}
	return rp, err
}

func (r *permissionRepo) UpdateResourceStatus(ctx context.Context, id string, status string, checkerID string) error {
	tenantID := middleware.GetTenantID(ctx)
	return r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		_, err := tx.Exec(ctx,
			`UPDATE resource_permissions SET status = $2, checker_id = $3, updated_at = now() WHERE id = $1`,
			id, status, checkerID,
		)
		return err
	})
}

func (r *permissionRepo) ListByTenant(ctx context.Context, tenantID string, roleID string) ([]*biz.Permission, error) {
	var perms []*biz.Permission
	err := r.data.ExecInTenant(ctx, tenantID, func(ctx context.Context, tx pgx.Tx) error {
		q := `SELECT id, tenant_id, resource, action FROM permissions WHERE tenant_id = $1`
		args := []interface{}{tenantID}
		if roleID != "" {
			q = `SELECT p.id, p.tenant_id, p.resource, p.action FROM permissions p
				 JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = $1`
			args = []interface{}{roleID}
		}
		rows, err := tx.Query(ctx, q, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			p := &biz.Permission{}
			if err := rows.Scan(&p.ID, &p.TenantID, &p.Resource, &p.Action); err != nil {
				return err
			}
			perms = append(perms, p)
		}
		return rows.Err()
	})
	return perms, err
}

// Redis cache

type permissionCache struct {
	rdb *redis.Client
}

func NewPermissionCache(data *Data) biz.PermissionCache {
	return &permissionCache{rdb: data.rdb.Client}
}

func cacheKey(userID, tenantID, resource, action, resourceID string) string {
	return fmt.Sprintf("perm:%s:%s:%s:%s:%s", userID, tenantID, resource, action, resourceID)
}

func (c *permissionCache) Get(ctx context.Context, userID, tenantID, resource, action, resourceID string) (*biz.CachedPermission, bool) {
	val, err := c.rdb.Get(ctx, cacheKey(userID, tenantID, resource, action, resourceID)).Result()
	if err != nil {
		return nil, false
	}
	allowedText, source, ok := strings.Cut(val, ":")
	if !ok {
		return nil, false
	}
	allowed := allowedText == "true"
	return &biz.CachedPermission{Allowed: allowed, Source: source}, true
}

func (c *permissionCache) Set(ctx context.Context, userID, tenantID, resource, action, resourceID string, allowed bool, source string) {
	val := fmt.Sprintf("%v:%s", allowed, source)
	c.rdb.Set(ctx, cacheKey(userID, tenantID, resource, action, resourceID), val, 5*time.Minute)
}

func (c *permissionCache) InvalidateUser(ctx context.Context, userID, tenantID string) {
	pattern := fmt.Sprintf("perm:%s:%s:*", userID, tenantID)
	var cursor uint64
	for {
		keys, next, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return
		}
		if len(keys) > 0 {
			c.rdb.Del(ctx, keys...)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
}
