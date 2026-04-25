import { Injectable, signal, computed, inject, effect, untracked } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { AuthService } from './auth.service';
import { TenantService } from './tenant.service';

interface ListPermissionsResponse {
  permissions: string[]; // e.g. ["user:read", "role:create"]
}

@Injectable({ providedIn: 'root' })
export class PermissionService {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private tenantService = inject(TenantService);
  private apiUrl = (window as any).__env?.apiUrl ?? 'http://localhost:8000';

  private _permissions = signal<Set<string>>(new Set());
  readonly isLoading = signal(false);

  /** Danh sách permissions dạng "resource:action" */
  readonly permissionList = computed(() => [...this._permissions()]);

  constructor() {
    // Tự động reload khi đổi tenant (và đã login)
    effect(() => {
      const isAuthenticated = this.authService.isAuthenticated();
      const tenantId = this.tenantService.selectedTenantId();

      if (isAuthenticated && tenantId) {
        untracked(() => this.loadPermissions());
      } else if (!isAuthenticated) {
        untracked(() => this._permissions.set(new Set()));
      }
    });
  }

  async loadPermissions(): Promise<void> {
    if (!this.authService.isAuthenticated()) {
      return;
    }

    this.isLoading.set(true);
    try {
      const resp = await firstValueFrom(
        this.http.get<ListPermissionsResponse>(`${this.apiUrl}/v1/me/permissions`)
      );
      this._permissions.set(new Set(resp.permissions ?? []));
    } catch (err) {
      console.error('Failed to load permissions', err);
      this._permissions.set(new Set());
    } finally {
      this.isLoading.set(false);
    }
  }

  /**
   * Kiểm tra user có quyền không.
   * @param permission Chuỗi "resource:action", vd: "user:read"
   */
  hasPermission(permission: string): boolean {
    return this._permissions().has(permission);
  }

  /**
   * Kiểm tra có ÍT NHẤT một trong các quyền được cung cấp không.
   */
  hasAnyPermission(permissions: string[]): boolean {
    return permissions.some((p) => this._permissions().has(p));
  }
}
