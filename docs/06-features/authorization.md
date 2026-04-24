# Authorization Architecture — Kiến trúc Phân quyền

> Fine-Grained Access Control (FGAC) cho hệ thống tài chính
> RBAC/ABAC/ReBAC với Maker-Checker pattern

---

## 📋 Overview

Hệ thống phân quyền Arda sử dụng mô hình lai giữa **RBAC** (Role-Based Access Control), **ABAC** (Attribute-Based Access Control), và **ReBAC** (Relationship-Based Access Control) để cung cấp kiểm soát quyền hạn chi tiết theo tiêu chuẩn tài chính/ngân hàng.

### Key Features
- **Multi-tenant** — Phân quyền theo workspace/tenant
- **Hierarchical Roles** — Global, Tenant, Resource levels
- **Maker-Checker** — Tách biệt người lập và người duyệt
- **Audit Trail** — Log tất cả actions
- **Forward Auth** — Enforcement tại API Gateway

---

## 🏗️ Architecture Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (MFE)                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  UI Components│  │  Auth Guard │  │ Permission  │     │
│  │  (Hide/Show) │  │  (Route)    │  │  Check      │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   API Gateway (APISIX)                      │
│  ┌──────────────┐                                           │
│  │  Forward Auth│  ──►  iam-service (Check Permission)     │
│  │  Plugin      │                                           │
│  └──────────────┘                                           │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              IAM Service (Go Kratos)                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   AuthN      │  │   AuthZ      │  │   Maker-     │     │
│  │  (Zitadel)   │  │  (RBAC/ABAC) │  │   Checker    │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              PostgreSQL (Policy Storage)                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Users      │  │   Roles      │  │ Permissions  │     │
│  │   Tenants    │  │   Members    │  │   Policies   │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

---

## 📋 Permission Model

### Permission Hierarchy

```
Global Role (System-wide)
├── Super Admin
│   └── Full access to all tenants and resources
│
└── System Admin
    └── Full access to system configuration

Tenant Role (Per-tenant)
├── Tenant Admin
│   └── Full access within tenant
├── Manager
│   └── Manage team members and resources
├── Staff
│   └── Standard access to business operations
└── Viewer
    └── Read-only access

Resource Permission (Per-resource)
├── Owner — Full control
├── Editor — Can edit but not delete
├── Contributor — Can add but not edit
└── Viewer — Read-only
```

### Permission Format

Permissions follow the format: `<resource>:<action>[:<subresource>]`

| Resource | Actions | Examples |
|----------|---------|----------|
| `accounting` | `read`, `write`, `post`, `reverse` | `accounting:read`, `accounting:post` |
| `loan` | `read`, `create`, `approve`, `disburse` | `loan:create`, `loan:approve` |
| `crm` | `read`, `create`, `update`, `delete` | `crm:read`, `crm:update` |
| `admin` | `read`, `manage`, `audit` | `admin:manage` |
| `report` | `read`, `export`, `schedule` | `report:export` |

---

## 🗄️ Database Schema

### Core Tables

```sql
-- Tenants
CREATE TABLE tenants (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    owner_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    settings JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenants_slug ON tenants(slug);
CREATE INDEX idx_tenants_owner ON tenants(owner_id);

-- Users
CREATE TABLE users (
    id VARCHAR(36) PRIMARY KEY,
    external_id VARCHAR(100) UNIQUE, -- Zitadel user ID
    email VARCHAR(200) NOT NULL UNIQUE,
    display_name VARCHAR(200) NOT NULL,
    avatar_url VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_external ON users(external_id);

-- Memberships (User-Tenant relationship)
CREATE TABLE memberships (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    role VARCHAR(50) NOT NULL, -- TENANT_ADMIN, MANAGER, STAFF, VIEWER
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    invited_by VARCHAR(36),
    joined_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (invited_by) REFERENCES users(id),
    UNIQUE(user_id, tenant_id)
);

CREATE INDEX idx_memberships_user ON memberships(user_id);
CREATE INDEX idx_memberships_tenant ON memberships(tenant_id);

-- Roles
CREATE TABLE roles (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36), -- NULL for global roles
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, name)
);

CREATE INDEX idx_roles_tenant ON roles(tenant_id);

-- Role Permissions (Many-to-Many)
CREATE TABLE role_permissions (
    id VARCHAR(36) PRIMARY KEY,
    role_id VARCHAR(36) NOT NULL,
    permission VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE(role_id, permission)
);

CREATE INDEX idx_rp_role ON role_permissions(role_id);

-- User Roles (Many-to-Many)
CREATE TABLE user_roles (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    granted_by VARCHAR(36),
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (granted_by) REFERENCES users(id),
    UNIQUE(user_id, role_id, tenant_id)
);

CREATE INDEX idx_ur_user ON user_roles(user_id);
CREATE INDEX idx_ur_role ON user_roles(role_id);
CREATE INDEX idx_ur_tenant ON user_roles(tenant_id);

-- Resource Permissions (Fine-grained)
CREATE TABLE resource_permissions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    resource VARCHAR(100) NOT NULL, -- accounting, loan, crm
    action VARCHAR(100) NOT NULL, -- read, write, approve
    resource_id VARCHAR(36), -- Specific resource ID
    allowed BOOLEAN NOT NULL DEFAULT true,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE', -- ACTIVE, PENDING_APPROVAL, EXPIRED
    approved_by VARCHAR(36),
    approved_at TIMESTAMP,
    expires_at TIMESTAMP,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (approved_by) REFERENCES users(id),
    UNIQUE(user_id, tenant_id, resource, action, resource_id)
);

CREATE INDEX idx_rsrc_user ON resource_permissions(user_id);
CREATE INDEX idx_rsrc_tenant ON resource_permissions(tenant_id);
CREATE INDEX idx_rsrc_status ON resource_permissions(status);

-- Audit Logs
CREATE TABLE audit_logs (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36),
    user_id VARCHAR(36),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(36),
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    status_code INT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_tenant ON audit_logs(tenant_id);
CREATE INDEX idx_audit_user ON audit_logs(user_id);
CREATE INDEX idx_audit_created ON audit_logs(created_at);
CREATE INDEX idx_audit_action ON audit_logs(action);
```

---

## 🔐 Forward Auth Flow

### 1. Request Flow

```
1. User requests API endpoint
   ↓
2. APISIX Gateway receives request
   ↓
3. APISIX extracts JWT from Authorization header
   ↓
4. APISIX calls IAM Service ForwardAuth endpoint
   ↓
5. IAM Service validates JWT signature
   ↓
6. IAM Service extracts user_id and tenant_id from JWT
   ↓
7. IAM Service checks user permissions
   ↓
8. IAM Service returns ALLOW or DENY
   ↓
9. APISIX forwards or blocks request
```

### 2. ForwardAuth Implementation

```go
// iam-service/internal/biz/auth.go
func (uc *AuthUsecase) ForwardAuth(ctx context.Context, method, path, token string) (bool, string, string, error) {
    // 1. Extract token
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }
    if token == "" {
        return false, "", "", fmt.Errorf("missing token")
    }

    // 2. Validate JWT
    set, err := uc.jwks.getSet()
    if err != nil {
        return false, "", "", err
    }

    tok, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
        kid, ok := t.Header["kid"].(string)
        if !ok {
            return nil, fmt.Errorf("missing kid in JWT header")
        }
        key, found := set.LookupKeyID(kid)
        if !found {
            return nil, fmt.Errorf("key %s not found in JWKS", kid)
        }
        var raw interface{}
        if err := key.Raw(&raw); err != nil {
            return nil, err
        }
        return raw, nil
    })
    if err != nil || !tok.Valid {
        return false, "", "", fmt.Errorf("invalid token")
    }

    // 3. Extract claims
    claims, ok := tok.Claims.(jwt.MapClaims)
    if !ok {
        return false, "", "", fmt.Errorf("invalid claims")
    }

    userID := stringClaimVal(claims["sub"])
    tenantID := stringClaimVal(claims["tenant_id"])

    // 4. Map method+path to resource+action
    resource, action := uc.mapPathToAction(method, path)

    // 5. Check permissions
    allowed, source, err := uc.perms.CheckPermission(ctx, userID, tenantID, resource, action, "")
    if err != nil {
        return false, userID, tenantID, err
    }

    if !allowed {
        uc.log.Warnf("Permission denied: user=%s tenant=%s resource=%s action=%s",
            userID, tenantID, resource, action)
        return false, userID, tenantID, nil
    }

    // 6. Log access
    uc.logAccess(ctx, userID, tenantID, method, path, resource, action, allowed)

    return true, userID, tenantID, nil
}

// Map HTTP method+path to resource+action
func (uc *AuthUsecase) mapPathToAction(method, path string) (string, string) {
    rules := []struct {
        method  string
        pattern *regexp.Regexp
        resource string
        action   string
    }{
        // Accounting
        {http.MethodGet, regexp.MustCompile(`^/v1/accounting/coa`), "accounting", "read"},
        {http.MethodPost, regexp.MustCompile(`^/v1/accounting/coa$`), "accounting", "create"},
        {http.MethodPost, regexp.MustCompile(`^/v1/accounting/journal`), "accounting", "write"},
        {http.MethodPost, regexp.MustCompile(`^/v1/accounting/journal/[^/]+/post`), "accounting", "post"},

        // Loan
        {http.MethodGet, regexp.MustCompile(`^/v1/loan/contracts`), "loan", "read"},
        {http.MethodPost, regexp.MustCompile(`^/v1/loan/contracts$`), "loan", "create"},
        {http.MethodPost, regexp.MustCompile(`^/v1/loan/contracts/[^/]+/approve`), "loan", "approve"},
        {http.MethodPost, regexp.MustCompile(`^/v1/loan/contracts/[^/]+/disburse`), "loan", "disburse"},

        // CRM
        {http.MethodGet, regexp.MustCompile(`^/v1/crm/customers`), "crm", "read"},
        {http.MethodPost, regexp.MustCompile(`^/v1/crm/customers$`), "crm", "create"},
        {http.MethodPut, regexp.MustCompile(`^/v1/crm/customers/`), "crm", "update"},
        {http.MethodDelete, regexp.MustCompile(`^/v1/crm/customers/`), "crm", "delete"},

        // Admin
        {http.MethodGet, regexp.MustCompile(`^/v1/admin/users`), "admin", "read"},
        {http.MethodPost, regexp.MustCompile(`^/v1/admin/users$`), "admin", "manage"},
        {http.MethodGet, regexp.MustCompile(`^/v1/admin/roles`), "admin", "read"},
        {http.MethodPost, regexp.MustCompile(`^/v1/admin/roles$`), "admin", "manage"},
    }

    for _, rule := range rules {
        if rule.method == method && rule.pattern.MatchString(path) {
            return rule.resource, rule.action
        }
    }

    // Default: public access
    return "public", "access"
}
```

### 3. Permission Check Logic

```go
// iam-service/internal/biz/permission.go
func (uc *PermissionUsecase) CheckPermission(ctx context.Context, userID, tenantID, resource, action, resourceID string) (bool, string, error) {
    // 1. Check cache first
    cacheKey := fmt.Sprintf("perm:%s:%s:%s:%s:%s", userID, tenantID, resource, action, resourceID)
    if cached, found := uc.cache.Get(ctx, cacheKey); found {
        return cached.(bool), "cache", nil
    }

    // 2. Check if user is super admin
    isSuperAdmin, err := uc.isSuperAdmin(ctx, userID)
    if err != nil {
        return false, "", err
    }
    if isSuperAdmin {
        uc.cache.Set(ctx, cacheKey, true, 5*time.Minute)
        return true, "super_admin", nil
    }

    // 3. Check user's roles in tenant
    roles, err := uc.getUserRoles(ctx, userID, tenantID)
    if err != nil {
        return false, "", err
    }

    // 4. Check permissions from roles
    for _, role := range roles {
        hasPermission, err := uc.roleHasPermission(ctx, role.ID, resource, action)
        if err != nil {
            continue
        }
        if hasPermission {
            uc.cache.Set(ctx, cacheKey, true, 5*time.Minute)
            return true, fmt.Sprintf("role:%s", role.Name), nil
        }
    }

    // 5. Check resource-specific permissions
    if resourceID != "" {
        hasResourcePerm, err := uc.hasResourcePermission(ctx, userID, tenantID, resource, action, resourceID)
        if err != nil {
            return false, "", err
        }
        if hasResourcePerm {
            uc.cache.Set(ctx, cacheKey, true, 5*time.Minute)
            return true, "resource_permission", nil
        }
    }

    // 6. Permission denied
    uc.cache.Set(ctx, cacheKey, false, 5*time.Minute)
    return false, "denied", nil
}

func (uc *PermissionUsecase) isSuperAdmin(ctx context.Context, userID string) (bool, error) {
    // Check if user has SUPER_ADMIN role
    count, err := uc.urRepo.CountByUserIDAndRoleName(ctx, userID, "SUPER_ADMIN")
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (uc *PermissionUsecase) getUserRoles(ctx context.Context, userID, tenantID string) ([]*Role, error) {
    // Get all roles for user in tenant
    return uc.urRepo.FindByUserIDAndTenantID(ctx, userID, tenantID)
}

func (uc *PermissionUsecase) roleHasPermission(ctx context.Context, roleID, resource, action string) (bool, error) {
    // Check if role has the permission
    permission := fmt.Sprintf("%s:%s", resource, action)
    count, err := uc.rpRepo.CountByRoleIDAndPermission(ctx, roleID, permission)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}

func (uc *PermissionUsecase) hasResourcePermission(ctx context.Context, userID, tenantID, resource, action, resourceID string) (bool, error) {
    // Check if user has specific resource permission
    perm, err := uc.rsrcRepo.Find(ctx, userID, tenantID, resource, action, resourceID)
    if err != nil {
        return false, err
    }

    // Check if permission is active and not expired
    if perm.Status != "ACTIVE" {
        return false, nil
    }
    if perm.ExpiresAt != nil && perm.ExpiresAt.Before(time.Now()) {
        return false, nil
    }

    return perm.Allowed, nil
}
```

---

## 🔄 Maker-Checker Pattern

### Overview

Maker-Checker pattern tách biệt người tạo/xử lý (Maker) và người duyệt (Checker) để đảm bảo kiểm soát nội dung.

### Implementation

```go
// iam-service/internal/biz/maker_checker.go
type MakerCheckerUsecase struct {
    rsrcRepo data.ResourcePermissionRepo
    userRepo data.UserRepo
    log      *log.Helper
}

// GrantResourcePermission requires approval
func (uc *MakerCheckerUsecase) GrantResourcePermission(ctx context.Context, req *GrantResourcePermissionRequest) (*ResourcePermission, error) {
    // Check if action requires approval
    if uc.requiresApproval(req.Action) {
        // Create pending approval
        perm := &ResourcePermission{
            ID:         uuid.New().String(),
            UserID:     req.UserID,
            TenantID:   req.TenantID,
            Resource:   req.Resource,
            Action:     req.Action,
            ResourceID: req.ResourceID,
            Allowed:    req.Allowed,
            Status:     "PENDING_APPROVAL",
            CreatedAt:  time.Now(),
        }

        if err := uc.rsrcRepo.Create(ctx, perm); err != nil {
            return nil, err
        }

        // Notify checkers
        uc.notifyCheckers(ctx, perm)

        return perm, nil
    }

    // No approval needed, create directly
    perm := &ResourcePermission{
        ID:         uuid.New().String(),
        UserID:     req.UserID,
        TenantID:   req.TenantID,
        Resource:   req.Resource,
        Action:     req.Action,
        ResourceID: req.ResourceID,
        Allowed:    req.Allowed,
        Status:     "ACTIVE",
        CreatedAt:  time.Now(),
    }

    if err := uc.rsrcRepo.Create(ctx, perm); err != nil {
        return nil, err
    }

    return perm, nil
}

// ApprovePermission approves a pending permission
func (uc *MakerCheckerUsecase) ApprovePermission(ctx context.Context, permissionID, checkerID string) (*ResourcePermission, error) {
    perm, err := uc.rsrcRepo.FindByID(ctx, permissionID)
    if err != nil {
        return nil, err
    }

    if perm.Status != "PENDING_APPROVAL" {
        return nil, fmt.Errorf("permission is not pending approval")
    }

    // Check if checker can approve (different from maker)
    if perm.UserID == checkerID {
        return nil, fmt.Errorf("cannot approve own permission request")
    }

    // Update permission
    perm.Status = "ACTIVE"
    perm.ApprovedBy = checkerID
    perm.ApprovedAt = time.Now()

    if err := uc.rsrcRepo.Update(ctx, perm); err != nil {
        return nil, err
    }

    // Clear permission cache
    uc.clearPermissionCache(ctx, perm.UserID, perm.TenantID)

    return perm, nil
}

// RejectPermission rejects a pending permission
func (uc *MakerCheckerUsecase) RejectPermission(ctx context.Context, permissionID, checkerID, reason string) (*ResourcePermission, error) {
    perm, err := uc.rsrcRepo.FindByID(ctx, permissionID)
    if err != nil {
        return nil, err
    }

    if perm.Status != "PENDING_APPROVAL" {
        return nil, fmt.Errorf("permission is not pending approval")
    }

    perm.Status = "REJECTED"
    perm.Reason = reason

    if err := uc.rsrcRepo.Update(ctx, perm); err != nil {
        return nil, err
    }

    return perm, nil
}

// ListPendingApprovals lists all pending approvals for a tenant
func (uc *MakerCheckerUsecase) ListPendingApprovals(ctx context.Context, tenantID string) ([]*ResourcePermission, error) {
    return uc.rsrcRepo.FindByTenantAndStatus(ctx, tenantID, "PENDING_APPROVAL")
}

// requiresApproval checks if action requires approval
func (uc *MakerCheckerUsecase) requiresApproval(action string) bool {
    // High-risk actions require approval
    requiringApproval := []string{
        "loan:approve",
        "loan:disburse",
        "accounting:post",
        "accounting:reverse",
        "admin:manage",
    }

    for _, a := range requiringApproval {
        if action == a {
            return true
        }
    }
    return false
}

// notifyCheckers sends notifications to checkers
func (uc *MakerCheckerUsecase) notifyCheckers(ctx context.Context, perm *ResourcePermission) {
    // Get users with approval permissions in tenant
    checkers, err := uc.userRepo.FindByPermission(ctx, perm.TenantID, "admin:approve")
    if err != nil {
        uc.log.Errorf("Failed to get checkers: %v", err)
        return
    }

    // Send notification to each checker
    for _, checker := range checkers {
        // TODO: Send via notification service
        uc.log.Infof("Notify checker %s about pending approval %s", checker.ID, perm.ID)
    }
}

func (uc *MakerCheckerUsecase) clearPermissionCache(ctx context.Context, userID, tenantID string) {
    // Clear all permission caches for user
    // This will be implemented with cache invalidation
}
```

---

## 🎨 Frontend Integration

### Permission Guard

```typescript
// libs/auth/src/lib/guards/permission.guard.ts
import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { ArdaAuthService } from '../services/auth.service';
import { PermissionService } from '../services/permission.service';

export const permissionGuard = (permission: string): CanActivateFn => {
  return async (route, state) => {
    const authService = inject(ArdaAuthService);
    const router = inject(Router);
    const permService = inject(PermissionService);

    if (!authService.isAuthenticated) {
      router.navigate(['/login'], { queryParams: { returnUrl: state.url } });
      return false;
    }

    const hasPermission = await permService.hasPermission(permission);

    if (!hasPermission) {
      router.navigate(['/unauthorized']);
      return false;
    }

    return true;
  };
};
```

### Permission Service

```typescript
// libs/auth/src/lib/services/permission.service.ts
import { Injectable } from '@angular/core';
import { ArdaAuthService } from './auth.service';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class PermissionService {
  private userPermissions: string[] = [];
  private permissionsLoaded = false;

  constructor(
    private auth: ArdaAuthService,
    private http: HttpClient,
  ) {}

  loadUserPermissions(): Observable<string[]> {
    return this.http.get<{ permissions: string[] }>('/api/v1/me/permissions', {
      headers: {
        Authorization: `Bearer ${this.auth.accessToken}`,
      },
    }).pipe(
      map(response => {
        this.userPermissions = response.permissions || [];
        this.permissionsLoaded = true;
        return this.userPermissions;
      }),
      catchError(() => {
        this.userPermissions = [];
        this.permissionsLoaded = true;
        return of([]);
      })
    );
  }

  hasPermission(permission: string): boolean {
    if (!this.permissionsLoaded) {
      return false;
    }

    // Check exact permission
    if (this.userPermissions.includes(permission)) {
      return true;
    }

    // Check wildcard permission (e.g., "accounting:*")
    const [resource] = permission.split(':');
    const wildcard = `${resource}:*`;
    if (this.userPermissions.includes(wildcard)) {
      return true;
    }

    return false;
  }

  hasAnyPermission(permissions: string[]): boolean {
    return permissions.some(p => this.hasPermission(p));
  }

  hasAllPermissions(permissions: string[]): boolean {
    return permissions.every(p => this.hasPermission(p));
  }

  getUserPermissions(): string[] {
    return this.userPermissions;
  }
}
```

### Permission Directive

```typescript
// libs/auth/src/lib/directives/permission.directive.ts
import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { PermissionService } from '../services/permission.service';

@Directive({
  selector: '[ardaPermission]',
  standalone: true,
})
export class PermissionDirective {
  @Input() ardaPermission: string | string[] = '';

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private permissionService: PermissionService,
  ) {
    this.updateView();
  }

  private updateView() {
    this.viewContainer.clear();

    let hasPermission = false;

    if (typeof this.ardaPermission === 'string') {
      hasPermission = this.permissionService.hasPermission(this.ardaPermission);
    } else {
      hasPermission = this.permissionService.hasAnyPermission(this.ardaPermission);
    }

    if (hasPermission) {
      this.viewContainer.createEmbeddedView(this.templateRef);
    }
  }
}
```

### Usage in Template

```html
<!-- Hide/Show based on permission -->
<button
  ardaPermission="loan:approve"
  pButton
  label="Approve"
  (click)="approveLoan()"></button>

<!-- Multiple permissions -->
<div *ardaPermission="['admin:manage', 'admin:read']">
  Admin content
</div>

<!-- Else clause -->
<button
  *ardaPermission="loan:disburse; else cannotDisburse"
  pButton
  label="Disburse">
</button>

<ng-template #cannotDisburse>
  <p-button label="Disburse" disabled></p-button>
</ng-template>
```

---

## 📊 APISIX Configuration

### Forward Auth Plugin

```yaml
# apps/gateway/apisix/base/plugins/forward-auth.yaml
apiVersion: apisix.apache.org/v2
kind: ApisixPluginConfig
metadata:
  name: forward-auth
  namespace: gateway
spec:
  plugins:
  - name: forward-auth
    enable: true
    config:
      uri: http://iam-service.arda-dev.svc.cluster.local/v1/auth/forward
      request_headers:
        - "X-Forwarded-Method: $request_method"
        - "X-Forwarded-Path: $request_uri"
        - "X-Forwarded-Token: $http_authorization"
      keepalive: true
      keepalive_timeout: 60s
      keepalive_pool: 5
```

### Route Configuration

```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: accounting-routes
  namespace: gateway
spec:
  http:
  - name: accounting-coa
    match:
      paths:
      - /api/v1/accounting/coa*
      methods:
      - GET
      - POST
      - PUT
      - DELETE
    plugins:
    - name: forward-auth
      enable: true
    backends:
    - serviceName: accounting-service
      servicePort: 80
```

---

## 🧪 Testing

### Permission Check Tests

```go
// iam-service/internal/biz/permission_test.go
func TestPermissionUsecase_CheckPermission(t *testing.T) {
    tests := []struct {
        name          string
        userID        string
        tenantID      string
        resource      string
        action        string
        expected      bool
        expectedSource string
    }{
        {
            name:     "Super admin has all permissions",
            userID:   "super-admin-id",
            tenantID: "tenant-1",
            resource: "loan",
            action:   "approve",
            expected: true,
            expectedSource: "super_admin",
        },
        {
            name:     "User with role has permission",
            userID:   "user-id",
            tenantID: "tenant-1",
            resource: "accounting",
            action:   "read",
            expected: true,
            expectedSource: "role:accountant",
        },
        {
            name:     "User without permission denied",
            userID:   "user-id",
            tenantID: "tenant-1",
            resource: "admin",
            action:   "manage",
            expected: false,
            expectedSource: "denied",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup test
            // ...

            // Execute
            allowed, source, err := uc.CheckPermission(ctx, tt.userID, tt.tenantID, tt.resource, tt.action, "")

            // Assert
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, allowed)
            assert.Equal(t, tt.expectedSource, source)
        })
    }
}
```

---

*Last Updated: 2026-04-24*
