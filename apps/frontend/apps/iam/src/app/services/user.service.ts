import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';

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

  async listUsers(tenantId: string): Promise<User[]> {
    const resp = await firstValueFrom(
      this.http.get<{ users: any[] }>(`/api/v1/users?tenant_id=${tenantId}`)
    );
    return (resp.users || []).map(u => ({
      id: u.id,
      email: u.email,
      displayName: u.display_name,
      createdAt: u.created_at
    }));
  }

  async listRoles(tenantId: string): Promise<Role[]> {
    const resp = await firstValueFrom(
      this.http.get<{ roles: any[] }>(`/api/v1/roles?tenant_id=${tenantId}`)
    );
    return (resp.roles || []).map(r => ({
      id: r.id,
      name: r.name,
      description: r.description,
      permissions: r.permissions
    }));
  }

  async createRole(role: Partial<Role>, tenantId: string): Promise<Role> {
    const resp = await firstValueFrom(
      this.http.post<any>('/api/v1/roles', {
        ...role,
        tenant_id: tenantId
      })
    );
    return {
      id: resp.id,
      name: resp.name,
      description: resp.description,
      permissions: resp.permissions
    };
  }

  async updateRole(id: string, role: Partial<Role>, tenantId: string): Promise<Role> {
    const resp = await firstValueFrom(
      this.http.put<any>(`/api/v1/roles/${id}`, {
        ...role,
        tenant_id: tenantId
      })
    );
    return {
      id: resp.id,
      name: resp.name,
      description: resp.description,
      permissions: resp.permissions
    };
  }

  async deleteRole(id: string, tenantId: string): Promise<void> {
    await firstValueFrom(
      this.http.delete(`/api/v1/roles/${id}?tenant_id=${tenantId}`)
    );
  }

  async assignRole(userId: string, roleId: string, tenantId: string): Promise<void> {
    await firstValueFrom(
      this.http.post('/api/v1/assign-role', {
        user_id: userId,
        role_id: roleId,
        tenant_id: tenantId
      })
    );
  }

  async inviteMember(externalId: string, role: string, tenantId: string): Promise<void> {
    await firstValueFrom(
      this.http.post(`/api/v1/tenants/${tenantId}/invite`, {
        external_id: externalId,
        role: role
      })
    );
  }

  async createUser(userData: any, tenantId: string): Promise<User> {
    const resp = await firstValueFrom(
      this.http.post<any>('/api/v1/users', {
        ...userData,
        tenant_id: tenantId
      })
    );
    return {
      id: resp.id,
      email: resp.email,
      displayName: resp.display_name,
      createdAt: resp.created_at
    };
  }

  async listPermissions(tenantId: string): Promise<string[]> {
    const resp = await firstValueFrom(
      this.http.get<{ permissions: string[] }>(`/api/v1/permissions?tenant_id=${tenantId}`)
    );
    return resp.permissions || [];
  }

  // Groups
  async listGroups(tenantId: string): Promise<Group[]> {
    const resp = await firstValueFrom(
      this.http.get<{ groups: any[] }>(`/api/v1/groups?tenant_id=${tenantId}`)
    );
    return (resp.groups || []).map(g => ({
      id: g.id,
      tenantId: g.tenant_id,
      name: g.name,
      description: g.description,
      createdAt: g.created_at,
      updatedAt: g.updated_at
    }));
  }

  async createGroup(group: Partial<Group>, tenantId: string): Promise<Group> {
    const resp = await firstValueFrom(
      this.http.post<any>('/api/v1/groups', {
        ...group,
        tenant_id: tenantId
      })
    );
    return {
      id: resp.id,
      tenantId: resp.tenant_id,
      name: resp.name,
      description: resp.description
    };
  }

  async updateGroup(id: string, group: Partial<Group>): Promise<Group> {
    const resp = await firstValueFrom(
      this.http.put<any>(`/api/v1/groups/${id}`, group)
    );
    return {
      id: resp.id,
      tenantId: resp.tenant_id,
      name: resp.name,
      description: resp.description
    };
  }

  async deleteGroup(id: string): Promise<void> {
    await firstValueFrom(this.http.delete(`/api/v1/groups/${id}`));
  }

  async listGroupMembers(groupId: string): Promise<User[]> {
    const resp = await firstValueFrom(
      this.http.get<{ users: any[] }>(`/api/v1/groups/${groupId}/members`)
    );
    return (resp.users || []).map(u => ({
      id: u.id,
      email: u.email,
      displayName: u.display_name,
      createdAt: u.created_at
    }));
  }

  async addGroupMember(groupId: string, userId: string): Promise<void> {
    await firstValueFrom(
      this.http.post(`/api/v1/groups/${groupId}/members`, { user_id: userId })
    );
  }

  async removeGroupMember(groupId: string, userId: string): Promise<void> {
    await firstValueFrom(
      this.http.delete(`/api/v1/groups/${groupId}/members/${userId}`)
    );
  }

  async listGroupRoles(groupId: string): Promise<Role[]> {
    const resp = await firstValueFrom(
      this.http.get<{ roles: any[] }>(`/api/v1/groups/${groupId}/roles`)
    );
    return (resp.roles || []).map(r => ({
      id: r.id,
      name: r.name,
      description: r.description,
      permissions: r.permissions,
      isSystem: r.is_system
    }));
  }

  async assignGroupRole(groupId: string, roleId: string): Promise<void> {
    await firstValueFrom(
      this.http.post(`/api/v1/groups/${groupId}/roles`, { role_id: roleId })
    );
  }

  async revokeGroupRole(groupId: string, roleId: string): Promise<void> {
    await firstValueFrom(
      this.http.delete(`/api/v1/groups/${groupId}/roles/${roleId}`)
    );
  }
}
