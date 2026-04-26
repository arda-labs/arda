import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';

export interface User {
  id: string;
  email: string;
  displayName: string;
  createdAt: string;
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

@Injectable({ providedIn: 'root' })
export class UserService {
  private http = inject(HttpClient);

  listUsers(tenantId: string): Observable<User[]> {
    return this.http.get<{ users: any[] }>(`/api/v1/users?tenant_id=${tenantId}`).pipe(
      map(resp => (resp.users || []).map(u => ({
        id: u.id,
        email: u.email,
        displayName: u.display_name,
        createdAt: u.created_at
      })))
    );
  }

  listRoles(tenantId: string): Observable<Role[]> {
    return this.http.get<{ roles: any[] }>(`/api/v1/roles?tenant_id=${tenantId}`).pipe(
      map(resp => (resp.roles || []).map(r => ({
        id: r.id,
        name: r.name,
        description: r.description,
        permissions: r.permissions
      })))
    );
  }

  createRole(role: Partial<Role>, tenantId: string): Observable<Role> {
    return this.http.post<any>('/api/v1/roles', {
      ...role,
      tenant_id: tenantId
    }).pipe(
      map(resp => ({
        id: resp.id,
        name: resp.name,
        description: resp.description,
        permissions: resp.permissions
      }))
    );
  }

  updateRole(id: string, role: Partial<Role>, tenantId: string): Observable<Role> {
    return this.http.put<any>(`/api/v1/roles/${id}`, {
      ...role,
      tenant_id: tenantId
    }).pipe(
      map(resp => ({
        id: resp.id,
        name: resp.name,
        description: resp.description,
        permissions: resp.permissions
      }))
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

  inviteMember(externalId: string, role: string, tenantId: string): Observable<void> {
    return this.http.post<void>(`/api/v1/tenants/${tenantId}/invite`, {
      external_id: externalId,
      role: role
    });
  }

  createUser(userData: any, tenantId: string): Observable<User> {
    return this.http.post<any>('/api/v1/users', {
      ...userData,
      tenant_id: tenantId
    }).pipe(
      map(resp => ({
        id: resp.id,
        email: resp.email,
        displayName: resp.display_name,
        createdAt: resp.created_at
      }))
    );
  }

  listPermissions(tenantId: string): Observable<string[]> {
    return this.http.get<{ permissions: string[] }>(`/api/v1/permissions?tenant_id=${tenantId}`).pipe(
      map(resp => resp.permissions || [])
    );
  }

  // Groups
  listGroups(tenantId: string): Observable<Group[]> {
    return this.http.get<{ groups: any[] }>(`/api/v1/groups?tenant_id=${tenantId}`).pipe(
      map(resp => (resp.groups || []).map(g => ({
        id: g.id,
        tenantId: g.tenant_id,
        name: g.name,
        description: g.description,
        createdAt: g.created_at,
        updatedAt: g.updated_at
      })))
    );
  }

  createGroup(group: Partial<Group>, tenantId: string): Observable<Group> {
    return this.http.post<any>('/api/v1/groups', {
      ...group,
      tenant_id: tenantId
    }).pipe(
      map(resp => ({
        id: resp.id,
        tenantId: resp.tenant_id,
        name: resp.name,
        description: resp.description
      }))
    );
  }

  updateGroup(id: string, group: Partial<Group>): Observable<Group> {
    return this.http.put<any>(`/api/v1/groups/${id}`, group).pipe(
      map(resp => ({
        id: resp.id,
        tenantId: resp.tenant_id,
        name: resp.name,
        description: resp.description
      }))
    );
  }

  deleteGroup(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/groups/${id}`);
  }

  listGroupMembers(groupId: string): Observable<User[]> {
    return this.http.get<{ users: any[] }>(`/api/v1/groups/${groupId}/members`).pipe(
      map(resp => (resp.users || []).map(u => ({
        id: u.id,
        email: u.email,
        displayName: u.display_name,
        createdAt: u.created_at
      })))
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
      map(resp => (resp.roles || []).map(r => ({
        id: r.id,
        name: r.name,
        description: r.description,
        permissions: r.permissions,
        isSystem: r.is_system
      })))
    );
  }

  assignGroupRole(groupId: string, roleId: string): Observable<void> {
    return this.http.post<void>(`/api/v1/groups/${groupId}/roles`, { role_id: roleId });
  }

  revokeGroupRole(groupId: string, roleId: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/groups/${groupId}/roles/${roleId}`);
  }
}
