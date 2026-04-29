import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { PageRequest, PageResponse } from '@arda/core';

export interface User {
  id: string; // Global user id, used by role/group APIs.
  tenantUserId: string;
  tenantId: string;
  username: string;
  email: string;
  displayName: string;
  role: string;
  status: string;
  createdAt: string;
  updatedAt?: string;
  roles?: string[]; // Sẽ map từ logic membership
}

export interface Role {
  id: string;
  name: string;
  description: string;
  permissions?: string[];
  isSystem?: boolean;
}

export interface Group {
  id: string;
  tenantId: string;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface UserTenantAccess {
  tenantId: string;
  tenantName: string;
  tenantSlug: string;
  username: string;
  displayName: string;
  status: string;
  roles: Role[];
  permissions: string[];
  deploymentMode: string;
  authMode: string;
}

interface UserListResponse {
  users?: any[];
  nextPageToken?: string;
  next_page_token?: string;
}

interface RoleListResponse {
  roles?: any[];
  nextPageToken?: string;
  next_page_token?: string;
}

interface GroupListResponse {
  groups?: any[];
  nextPageToken?: string;
  next_page_token?: string;
}

interface UserTenantAccessResponse {
  userId?: string;
  user_id?: string;
  tenants?: any[];
}

@Injectable({ providedIn: 'root' })
export class UserService {
  private http = inject(HttpClient);

  listUsers(tenantId: string, page: PageRequest = { pageSize: 100 }): Observable<User[]> {
    return this.listUsersPage(tenantId, page).pipe(map(resp => resp.items));
  }

  listUsersPage(tenantId: string, page: PageRequest = { pageSize: 20 }): Observable<PageResponse<User>> {
    return this.http.get<UserListResponse>('/api/v1/users', {
      params: this.tenantPageParams(tenantId, page),
    }).pipe(
      map(resp => ({
        items: (resp.users || []).map(u => this.toUser(u)),
        nextPageToken: resp.nextPageToken ?? resp.next_page_token ?? '',
      }))
    );
  }

  getTenantUser(userId: string, tenantId: string): Observable<User> {
    return this.http.get<any>(`/api/v1/users/${encodeURIComponent(userId)}/tenant-account`, {
      params: new HttpParams().set('tenant_id', tenantId),
    }).pipe(
      map(resp => this.toUser(resp))
    );
  }

  listUserTenantAccess(userId: string): Observable<UserTenantAccess[]> {
    return this.http.get<UserTenantAccessResponse>(`/api/v1/users/${encodeURIComponent(userId)}/tenant-access`).pipe(
      map(resp => (resp.tenants ?? []).map(item => this.toUserTenantAccess(item)))
    );
  }

  listRoles(tenantId: string, page: PageRequest = { pageSize: 100 }): Observable<Role[]> {
    return this.listRolesPage(tenantId, page).pipe(map(resp => resp.items));
  }

  listRolesPage(tenantId: string, page: PageRequest = { pageSize: 20 }): Observable<PageResponse<Role>> {
    return this.http.get<RoleListResponse>('/api/v1/roles', {
      params: this.tenantPageParams(tenantId, page),
    }).pipe(
      map(resp => ({
        items: (resp.roles || []).map(r => this.toRole(r)),
        nextPageToken: resp.nextPageToken ?? resp.next_page_token ?? '',
      }))
    );
  }

  createRole(role: Partial<Role>, tenantId: string): Observable<Role> {
    return this.http.post<any>('/api/v1/roles', {
      ...role,
      tenant_id: tenantId
    }).pipe(
      map(resp => this.toRole(resp))
    );
  }

  updateRole(id: string, role: Partial<Role>, tenantId: string): Observable<Role> {
    return this.http.put<any>(`/api/v1/roles/${id}`, {
      ...role,
      tenant_id: tenantId
    }).pipe(
      map(resp => this.toRole(resp))
    );
  }

  deleteRole(id: string, tenantId: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/roles/${id}?tenant_id=${tenantId}`);
  }

  assignRole(userId: string, roleId: string, tenantId: string): Observable<void> {
    return this.http.post<void>('/api/v1/assign-role', {
      user_id: userId,
      role_id: roleId,
      tenant_id: tenantId
    });
  }

  listUserRoles(userId: string, tenantId: string): Observable<Role[]> {
    return this.http.get<RoleListResponse>(`/api/v1/users/${encodeURIComponent(userId)}/roles`, {
      params: new HttpParams().set('tenant_id', tenantId),
    }).pipe(
      map(resp => (resp.roles || []).map(r => this.toRole(r)))
    );
  }

  revokeRole(userId: string, roleId: string, tenantId: string): Observable<void> {
    return this.http.post<void>('/api/v1/revoke-role', {
      user_id: userId,
      role_id: roleId,
      tenant_id: tenantId
    });
  }

  inviteMember(externalId: string, username: string, displayName: string, role: string, tenantId: string): Observable<User> {
    return this.http.post<any>(`/api/v1/tenants/${tenantId}/invite`, {
      external_id: externalId,
      username,
      display_name: displayName,
      role,
    }).pipe(
      map(resp => this.toUser(resp))
    );
  }

  removeMember(userId: string, tenantId: string): Observable<void> {
    return this.http.delete<void>(
      `/api/v1/tenants/${encodeURIComponent(tenantId)}/members/${encodeURIComponent(userId)}`,
    );
  }

  createUser(userData: any, tenantId: string): Observable<User> {
    return this.http.post<any>('/api/v1/users', {
      ...userData,
      tenant_id: tenantId
    }).pipe(
      map(resp => this.toUser(resp))
    );
  }

  listPermissions(tenantId: string): Observable<string[]> {
    return this.http.get<{ permissions: string[] }>(`/api/v1/permissions?tenant_id=${tenantId}`).pipe(
      map(resp => resp.permissions || [])
    );
  }

  // Groups
  listGroups(tenantId: string, page: PageRequest = { pageSize: 100 }): Observable<Group[]> {
    return this.listGroupsPage(tenantId, page).pipe(map(resp => resp.items));
  }

  listGroupsPage(tenantId: string, page: PageRequest = { pageSize: 20 }): Observable<PageResponse<Group>> {
    return this.http.get<GroupListResponse>('/api/v1/groups', {
      params: this.tenantPageParams(tenantId, page),
    }).pipe(
      map(resp => ({
        items: (resp.groups || []).map(g => this.toGroup(g)),
        nextPageToken: resp.nextPageToken ?? resp.next_page_token ?? '',
      }))
    );
  }

  createGroup(group: Partial<Group>, tenantId: string): Observable<Group> {
    return this.http.post<any>('/api/v1/groups', {
      ...group,
      tenant_id: tenantId
    }).pipe(
      map(resp => this.toGroup(resp))
    );
  }

  updateGroup(id: string, group: Partial<Group>): Observable<Group> {
    return this.http.put<any>(`/api/v1/groups/${id}`, group).pipe(
      map(resp => this.toGroup(resp))
    );
  }

  deleteGroup(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/groups/${id}`);
  }

  listGroupMembers(groupId: string): Observable<User[]> {
    return this.http.get<{ users: any[] }>(`/api/v1/groups/${groupId}/members`).pipe(
      map(resp => (resp.users || []).map(u => this.toUser(u)))
    );
  }

  addGroupMember(groupId: string, userId: string): Observable<void> {
    return this.http.post<void>(`/api/v1/groups/${groupId}/members`, { user_id: userId });
  }

  removeGroupMember(groupId: string, userId: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/groups/${groupId}/members/${userId}`);
  }

  listGroupRoles(groupId: string): Observable<Role[]> {
    return this.http.get<{ roles: any[] }>(`/api/v1/groups/${groupId}/roles`).pipe(
      map(resp => (resp.roles || []).map(r => this.toRole(r)))
    );
  }

  assignGroupRole(groupId: string, roleId: string): Observable<void> {
    return this.http.post<void>(`/api/v1/groups/${groupId}/roles`, { role_id: roleId });
  }

  revokeGroupRole(groupId: string, roleId: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/groups/${groupId}/roles/${roleId}`);
  }

  private tenantPageParams(tenantId: string, page: PageRequest): HttpParams {
    let params = new HttpParams()
      .set('tenant_id', tenantId)
      .set('page_size', String(page.pageSize ?? 20));

    if (page.pageToken) {
      params = params.set('page_token', page.pageToken);
    }

    return params;
  }

  private toUser(u: any): User {
    return {
      id: u.user_id ?? u.userId ?? u.id,
      tenantUserId: u.id,
      tenantId: u.tenant_id ?? u.tenantId,
      username: u.username,
      email: u.email,
      displayName: u.display_name ?? u.displayName,
      role: u.role,
      status: u.status,
      createdAt: u.created_at ?? u.createdAt,
      updatedAt: u.updated_at ?? u.updatedAt,
    };
  }

  private toRole(r: any): Role {
    return {
      id: r.id,
      name: r.name,
      description: r.description,
      permissions: r.permissions,
      isSystem: r.is_system ?? r.isSystem,
    };
  }

  private toGroup(g: any): Group {
    return {
      id: g.id,
      tenantId: g.tenant_id ?? g.tenantId,
      name: g.name,
      description: g.description,
      createdAt: g.created_at ?? g.createdAt,
      updatedAt: g.updated_at ?? g.updatedAt,
    };
  }

  private toUserTenantAccess(item: any): UserTenantAccess {
    return {
      tenantId: item.tenant_id ?? item.tenantId,
      tenantName: item.tenant_name ?? item.tenantName,
      tenantSlug: item.tenant_slug ?? item.tenantSlug,
      username: item.username,
      displayName: item.display_name ?? item.displayName,
      status: item.status,
      roles: (item.roles ?? []).map((role: any) => this.toRole(role)),
      permissions: item.permissions ?? [],
      deploymentMode: item.deployment_mode ?? item.deploymentMode,
      authMode: item.auth_mode ?? item.authMode,
    };
  }
}
