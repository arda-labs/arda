# UI/UX Va Data Grid Strategy

> Cap nhat: 2026-04-29
> Stack hien tai: Angular 21, PrimeNG 21, Tailwind CSS 4.

## Huong UI Cho He Thong Enterprise

Arda la operational platform, khong phai landing page. UI nen toi uu cho nguoi dung lam viec hang ngay:

- Mat do thong tin cao nhung co nhom ro rang.
- Toolbar co filter/search/action o dung ngu canh.
- Table la first-class surface, khong boc nhieu lop card.
- Detail nen dung drawer hoac split view de user khong bi mat ngu canh danh sach.
- Action nguy hiem phai co confirm, ly do, va neu high-risk thi chuyen sang maker-checker.
- Badge/status/risk level phai co ngon ngu mau sac nhat quan.
- Moi man hinh quan tri can co empty state, loading state, error state, permission-denied state.

## IAM UI De Xuat

Man hinh `Users`:

- List ben trai/table chinh: username, display name, tenant role summary, group count, status, last login, risk flags.
- Detail panel/tabs: Profile, Roles, Groups, Effective Permissions, Exceptions, Audit.
- Gan role/group tu detail panel, khong mo nhieu dialog roi roi.
- Hien "Effective Permissions" de admin hieu user co quyen tu dau: direct role, group role, hierarchy, exception, policy.

Man hinh `Roles`:

- Table role: name, code, system/custom, risk level, permission count, assigned users/groups, approval required.
- Detail tabs: Overview, Permission Matrix, Users, Groups, Hierarchy, Approval, Audit.
- Permission Matrix gom theo module/resource/action, co search va bulk toggle.
- Role high-risk khong update truc tiep; tao access request.

Man hinh `Groups`:

- Table group: name, member count, role count, owner, status.
- Detail tabs: Members, Roles, Effective Permissions, Audit.
- Giai thich ro user nhan quyen tu group nao.

Man hinh `Access Requests`:

- Inbox cho checker.
- Queue co filter theo status/risk/module/request type.
- Detail hien before/after diff, maker, reason, affected users, policy impact.

## Data Grid Decision

Khuyen nghi giai doan 1: tiep tuc dung PrimeNG `p-table`, nhung khong dung truc tiep lung tung o tung man hinh. Nen tao mot internal `DataGrid` convention/wrapper theo query contract chung.

Ly do:

- PrimeNG da co san trong stack, hop voi theme hien tai.
- `p-table` ho tro pagination, sorting, filtering, lazy loading, virtual scroll, column resize, frozen columns, CSV export va stateful table.
- Cac man hinh IAM giai doan dau chu yeu la CRUD/management, chua can pivot/grouping/excel-like editing nang.
- Neu sau nay can grid rat lon, server-side row model, row grouping, pivot, Excel export nang, AG Grid Enterprise la ung vien manh nhung can license.

Khong nen chon TanStack Table lam grid mac dinh cho admin console luc nay. TanStack rat tot neu muon headless/custom UI, nhung voi Angular + PrimeNG, team se phai tu lam nhieu phan interaction, accessibility polish va visual consistency.

## Khi Nao Dung Thu Vien Nao

| Nhu cau | Lua chon |
| --- | --- |
| CRUD table, IAM, config, danh muc | PrimeNG `p-table` thong qua internal `DataGrid` convention |
| Tree/menu/role hierarchy | PrimeNG `p-tree`, `p-treeTable` hoac custom split view |
| Bang cuc lon can server-side grouping/pivot/master-detail | AG Grid Enterprise |
| Bang custom UI cuc cao, khong muon bi rang buoc visual component | TanStack Table |

## Internal DataGrid Contract

Moi list API nen chap nhan query chung:

```json
{
  "page": 1,
  "page_size": 25,
  "sort": [{ "field": "created_at", "direction": "desc" }],
  "filters": [
    { "field": "status", "op": "eq", "value": "active" },
    { "field": "keyword", "op": "contains", "value": "admin" }
  ],
  "columns": ["username", "display_name", "status"]
}
```

Response chung:

```json
{
  "items": [],
  "total": 0,
  "page": 1,
  "page_size": 25,
  "has_next": false
}
```

Frontend grid state nen gom:

- pagination
- sorting
- column filters
- global search
- column visibility
- density
- selected rows
- persisted state key theo tenant + route

## UI Tokens Can Chuan Hoa

- Density: `comfortable`, `compact`.
- Table row height: 44-48px cho compact enterprise.
- Status severity: success/info/warn/danger/secondary.
- Risk severity: low/medium/high/critical.
- Action layout: primary action o toolbar, row actions dung icon button co tooltip.
- Dialog chi dung cho create/simple edit; detail-heavy flow dung drawer/split panel.

## Tai Lieu Tham Khao

- PrimeNG Table: https://primeng.org/table
- AG Grid Angular: https://www.ag-grid.com/angular-data-grid/getting-started/
- AG Grid Community vs Enterprise: https://www.ag-grid.com/angular-data-grid/community-vs-enterprise/
- TanStack Angular Table: https://tanstack.com/table/latest/docs/framework/angular
