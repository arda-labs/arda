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
      description: r.description
    }));
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
}
