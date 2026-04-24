import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { TenantService } from './tenant.service';
import { PermissionService } from './permission.service';

export interface MenuItem {
  id: string;
  label: string;
  icon: string;
  routerLink: string[];
  items: MenuItem[];
}

/** Raw response from backend before transformation */
interface BackendMenuItem {
  id: string;
  name: string;
  icon: string;
  route: string;
  sort_order: number;
  children: BackendMenuItem[];
}

interface MenuResponse {
  items: BackendMenuItem[];
}

@Injectable({ providedIn: 'root' })
export class MenuService {
  private http = inject(HttpClient);
  private tenantService = inject(TenantService);
  private permService = inject(PermissionService);

  private _menuItems = signal<MenuItem[]>([]);
  readonly isLoading = signal(false);

  readonly menuItems = computed(() => {
    const items = this._menuItems();
    if (items.length === 0) return [];

    // Once permissions are loaded, filter items
    const loaded = this.permService.permissionList().length > 0;
    if (!loaded) return items;

    return items
      .map(item => this.filterByPermission(item))
      .filter(Boolean) as MenuItem[];
  });

  async loadMenu(): Promise<void> {
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const resp = await firstValueFrom(
        this.http.get<MenuResponse>(`/api/v1/me/menu`)
      );
      this._menuItems.set(this.mapBackendToFrontend(resp.items ?? []));
    } catch (err) {
      console.error('Failed to load menu', err);
      this._menuItems.set([]);
    } finally {
      this.isLoading.set(false);
    }
  }

  private mapBackendToFrontend(items: BackendMenuItem[]): MenuItem[] {
    return items.map(item => ({
      id: item.id,
      label: item.name ?? '',
      icon: item.icon ?? '',
      routerLink: item.route ? ['/app' + item.route] : [],
      items: item.children ? this.mapBackendToFrontend(item.children) : [],
    }));
  }

  private filterByPermission(item: MenuItem): MenuItem | null {
    if (item.items.length > 0) {
      const filteredChildren = item.items
        .map(child => this.filterByPermission(child))
        .filter(Boolean) as MenuItem[];

      return filteredChildren.length > 0
        ? { ...item, items: filteredChildren }
        : null;
    }

    // Items without route are category headers — always visible
    if (!item.routerLink || item.routerLink.length === 0) {
      return item;
    }

    return item;
  }
}
