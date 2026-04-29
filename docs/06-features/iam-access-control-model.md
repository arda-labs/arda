# IAM Access Control Model

> Cap nhat: 2026-04-29
> Pham vi: multi-tenant IAM, RBAC/ABAC, maker-checker, audit, UI quan tri quyen.
> Status: Active design, partially implemented in `iam-service` and IAM MFE.

## Muc Tieu

Arda IAM can tach ro 3 lop:

1. Platform administration: `super_admin` quan tri he thong va tenant, khong phu thuoc tenant role.
2. Tenant administration: `admin` va cac role quan tri trong tung tenant.
3. Business authorization: role, group, policy va resource exception dieu khien nghiep vu.

Thiet ke uu tien cach cac he thong tai chinh/ngan hang hay lam: it gan quyen truc tiep cho user, moi thay doi quyen nhay cam co maker-checker, tat ca co audit trail va co ky review dinh ky.

## Identity Va Tenant Boundary

`super_admin` la platform identity. Nguoi nay duoc luu trong `platform_admins` va co bypass o backend. Khong nen tao role `super_admin` trong moi tenant vi se gay nham lan ve pham vi quyen.

`tenant_admin` la tenant role. User co role nay chi full quyen trong tenant dang chon, khong nhin thay tenant khac neu khong co membership.

Moi API can xac dinh ro:

- `actor_user_id`: user noi bo trong bang `users`.
- `tenant_id`: tenant dang thao tac, lay tu selected tenant/header/context.
- `resource`, `action`, `resource_id`: doi tuong can check quyen.

## Target Schema

Model hien tai da co cac bang nen tang: `users`, `tenants`, `tenant_users`, `roles`, `permissions`, `role_permissions`, `user_roles`, `groups`, `group_members`, `group_roles`, `role_hierarchy`, `policies`, `resource_permissions`, `platform_admins`.

`users` la global identity. `tenant_users` la tenant account va chua username/display name/status theo tung tenant. Khong dung `users.username` lam username nghiep vu vi se gay xung dot giua cac tenant.

Huong chuan hoa tiep theo:

| Bang | Vai tro |
| --- | --- |
| `permission_catalog` | Catalog global cua cac permission code, vi du `iam.user.read`, `loan.contract.approve`. |
| `tenant_permission_entitlements` | Tenant nao duoc bat permission/module nao. |
| `roles` | Role theo tenant, co `code`, `risk_level`, `approval_required`, `status`. |
| `role_permissions` | Role co permission nao, them `effect` de ho tro allow/deny. |
| `role_hierarchy` | Role cha ke thua role con, chi dung cho role template/system role neu UI giai thich duoc. |
| `groups` | Don vi gan quyen theo phong ban/chuc nang. |
| `group_roles` | Gan role cho group, user ke thua quyen qua group. |
| `user_roles` | Direct role assignment, nen co `effective_from`, `expires_at`, `assigned_by`, `reason`, `status`. |
| `policies` | ABAC rule, vi du han muc phe duyet, maker khac checker, gio giao dich. |
| `resource_permissions` | Exception tren tai nguyen cu the, co expiry va maker-checker. |
| `access_requests` | Yeu cau thay doi quyen/role/group/policy. |
| `access_request_items` | Chi tiet thay doi trong mot request. |
| `access_review_campaigns` | Dot review quyen dinh ky. |
| `access_review_items` | Tung user/role/group can review trong campaign. |
| `audit_logs` | Lich su thay doi va truy cap. |

## Permission Code

Permission code nen di theo format dot-case:

```text
<module>.<resource>.<action>
```

Vi du:

- `iam.user.read`
- `iam.user.create`
- `iam.role.assign`
- `iam.policy.approve`
- `loan.contract.approve`
- `accounting.journal.post`

Trong giai do chuyen tiep, backend van co the chap nhan format cu `resource:action`. Migration moi bo sung `permissions.code` de map dan sang format moi ma khong pha forward-auth hien tai.

## Permission Evaluation

Thu tu check quyen nen giu on dinh:

1. Neu actor la `platform_admins` active: allow, source `super_admin`.
2. Kiem tra membership trong tenant.
3. Kiem tra deny exception active tren `resource_permissions`.
4. Kiem tra allow exception active tren `resource_permissions`.
5. Tong hop role truc tiep tu `user_roles`.
6. Tong hop role ke thua tu `group_roles`.
7. Mo rong role theo `role_hierarchy`.
8. Kiem tra `role_permissions`.
9. Ap dung `policies` ABAC neu permission/action co rule.
10. Deny mac dinh.

Deny exception nen uu tien hon allow role vi day la cach thuong dung de tam khoa quyen nhay cam ma khong phai sua role.

## Maker-Checker

Cac thao tac sau nen vao `access_requests` thay vi ghi thang:

- Gan/xoa role co `risk_level` tu `high` tro len.
- Sua permission cua role system/high-risk.
- Tao/sua policy ABAC.
- Gan resource exception.
- Cap quyen co `expires_at` trong thoi gian dai.
- Thay doi tenant admin.

Luot duyet phai dam bao:

- Maker va checker khac nhau.
- Checker co permission duyet phu hop.
- Trang thai request ro rang: `draft`, `pending`, `approved`, `rejected`, `cancelled`, `expired`.
- Khi approved moi apply thay doi vao bang dich.
- Audit ghi ca before/after state.

## UI Quan Tri

Nen tach thanh 2 console.

Platform Console danh cho `super_admin`:

- Tenants
- Platform Admins
- Permission Catalog
- Role Templates
- System Audit

Tenant Console danh cho tenant admin:

- Users
- Groups
- Roles
- Policies
- Access Requests
- Access Reviews
- Audit Logs

Man hinh user nen co tabs:

- Profile
- Tenants
- Roles
- Groups
- Effective Permissions
- Exceptions
- Audit

Man hinh role nen la trung tam quan tri:

- Overview
- Permission Matrix
- Users
- Groups
- Hierarchy
- Approval Settings
- Audit

## Roadmap Trien Khai

Giai doan 1:

- Bo sung schema governance additive.
- Giu permission check hien tai hoat dong.
- UI user/role/group chuyen sang layout ro hon: list view + detail panel/tabs.
- Chuan hoa table query contract: pagination, sort, filter, column visibility.

Giai doan 2:

- Them API effective permissions.
- Them API access requests va approve/reject.
- Them role permission matrix.
- Them audit view cho IAM.

Giai doan 3:

- Policy editor co validation.
- Access review campaign.
- Segregation of Duties rule.
- Tenant permission entitlements/module licensing.
