import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';

export interface MenuAdminItem {
  id: string;
  tenantId: string;
  parentId: string;
  name: string;
  slug: string;
  icon: string;
  route: string;
  sortOrder: number;
  enabled: boolean;
  permissionSlug: string;
}

export interface MenuAdminPayload {
  tenantId?: string;
  parentId: string;
  name: string;
  slug: string;
  icon: string;
  route: string;
  sortOrder: number;
  enabled: boolean;
  permissionSlug: string;
}

@Injectable({
  providedIn: 'root',
})
export class MenuAdmin {
  private http = inject(HttpClient);

  listMenus(tenantId: string): Observable<MenuAdminItem[]> {
    return this.http.get<{ menus: unknown[] }>(`/api/v1/menus?tenant_id=${encodeURIComponent(tenantId)}`).pipe(
      map(resp => (resp.menus ?? []).map(item => this.fromBackend(item))),
    );
  }

  createMenu(payload: MenuAdminPayload, tenantId: string): Observable<MenuAdminItem> {
    return this.http.post<unknown>('/api/v1/menus', this.toBackend({ ...payload, tenantId })).pipe(
      map(item => this.fromBackend(item)),
    );
  }

  updateMenu(id: string, payload: MenuAdminPayload): Observable<MenuAdminItem> {
    return this.http.put<unknown>(`/api/v1/menus/${id}`, this.toBackend(payload)).pipe(
      map(item => this.fromBackend(item)),
    );
  }

  deleteMenu(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/menus/${id}`);
  }

  private fromBackend(raw: unknown): MenuAdminItem {
    const item = raw as Record<string, unknown>;
    return {
      id: String(item['id'] ?? ''),
      tenantId: String(item['tenant_id'] ?? ''),
      parentId: String(item['parent_id'] ?? ''),
      name: String(item['name'] ?? ''),
      slug: String(item['slug'] ?? ''),
      icon: String(item['icon'] ?? ''),
      route: String(item['route'] ?? ''),
      sortOrder: Number(item['sort_order'] ?? 0),
      enabled: Boolean(item['enabled']),
      permissionSlug: String(item['permission_slug'] ?? ''),
    };
  }

  private toBackend(payload: MenuAdminPayload): Record<string, unknown> {
    return {
      tenant_id: payload.tenantId,
      parent_id: payload.parentId,
      name: payload.name,
      slug: payload.slug,
      icon: payload.icon,
      route: payload.route,
      sort_order: payload.sortOrder,
      enabled: payload.enabled,
      permission_slug: payload.permissionSlug,
    };
  }
}
