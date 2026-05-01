import { Injectable, signal, computed, inject, effect, untracked } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, firstValueFrom } from 'rxjs';
import { TenantService } from './tenant.service';
import { PermissionService } from './permission.service';
import { AuthService } from './auth.service';

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
  private authService = inject(AuthService);
  private tenantService = inject(TenantService);
  private permService = inject(PermissionService);

  private _menuItems = signal<MenuItem[]>([]);
  readonly isLoading = signal(false);

  readonly menuItems = computed(() => {
    const items = this._menuItems();

    // Default/Fallback menu for development and new modules
    const fallbackMenu: MenuItem[] = [
      { id: 'home', label: 'Trang chủ', icon: 'pi pi-home', routerLink: ['/home'], items: [] },
      {
        id: 'crm',
        label: 'Khách hàng (CRM)',
        icon: 'pi pi-users',
        routerLink: [],
        items: [
          { id: 'crm-register', label: 'Đăng ký mới', icon: 'pi pi-user-plus', routerLink: ['/crm/register/init'], items: [] },
          { id: 'crm-list', label: 'Danh sách khách hàng', icon: 'pi pi-list', routerLink: ['/crm/info/customer-list'], items: [] }
        ]
      },
      {
        id: 'bpm',
        label: 'Quy trình (BPM)',
        icon: 'pi pi-sitemap',
        routerLink: [],
        items: [
          { id: 'bpm-inbound', label: 'Giao dịch đến', icon: 'pi pi-download', routerLink: ['/bpm/inbound'], items: [] },
          { id: 'bpm-outbound', label: 'Giao dịch đi', icon: 'pi pi-upload', routerLink: ['/bpm/outbound'], items: [] },
          { id: 'bpm-monitor', label: 'Giám sát vận hành', icon: 'pi pi-chart-bar', routerLink: ['/bpm/monitor'], items: [] },
          { id: 'bpm-search', label: 'Tra cứu giao dịch', icon: 'pi pi-search', routerLink: ['/bpm/search'], items: [] },
          {
            id: 'bpm-config',
            label: 'Cấu hình hệ thống',
            icon: 'pi pi-cog',
            routerLink: [],
            items: [
              { id: 'bpm-cfg-assignment', label: 'Quy tắc chia bài', icon: 'pi pi-users', routerLink: ['/bpm/config/assignment'], items: [] },
              { id: 'bpm-cfg-sla', label: 'Cấu hình SLA', icon: 'pi pi-clock', routerLink: ['/bpm/config/sla'], items: [] },
              { id: 'bpm-cfg-desc', label: 'Cấu trúc diễn giải', icon: 'pi pi-comment', routerLink: ['/bpm/config/description'], items: [] }
            ]
          },
          { id: 'bpm-error', label: 'Xử lý lỗi (Hospital)', icon: 'pi pi-heart-fill', routerLink: ['/bpm/error-hospital'], items: [] }
        ]
      },
      {
        id: 'loan',
        label: 'Khoản vay (Loan)',
        icon: 'pi pi-money-bill',
        routerLink: [],
        items: [
          { id: 'loan-app', label: 'Hồ sơ vay', icon: 'pi pi-file', routerLink: ['/loan/application'], items: [] }
        ]
      },
      {
        id: 'hrm',
        label: 'Nhân sự (HRM)',
        icon: 'pi pi-id-card',
        routerLink: [],
        items: [
          { id: 'hrm-onboarding', label: 'Onboarding', icon: 'pi pi-user-plus', routerLink: ['/hrm/onboarding'], items: [] }
        ]
      },
      {
        id: 'iam',
        label: 'Quản trị hệ thống',
        icon: 'pi pi-shield',
        routerLink: [],
        items: [
          { id: 'iam-users', label: 'Người dùng', icon: 'pi pi-user', routerLink: ['/iam/users'], items: [] },
          { id: 'iam-roles', label: 'Vai trò & Quyền', icon: 'pi pi-key', routerLink: ['/iam/roles'], items: [] }
        ]
      }
    ];

    if (items.length === 0) return fallbackMenu;

    // Merge logic or prioritize backend items if available
    return items.length > 0 ? items : fallbackMenu;
  });

  constructor() {
    // Tự động reload khi đổi tenant (và đã login)
    effect(() => {
      const isAuthenticated = this.authService.isAuthenticated();
      const tenantId = this.tenantService.selectedTenantId();

      if (isAuthenticated && tenantId) {
        untracked(() => this.loadMenu());
      } else if (!isAuthenticated) {
        untracked(() => this._menuItems.set([]));
      }
    });
  }

  fetchMenu(): Observable<MenuResponse> {
    return this.http.get<MenuResponse>(`/api/v1/me/menu`);
  }

  async loadMenu(forceAuthenticated = false): Promise<void> {
    if (!forceAuthenticated && !this.authService.isAuthenticated()) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const resp = await firstValueFrom(this.fetchMenu());
      this._menuItems.set(this.mapBackendToFrontend(resp.items ?? []));
    } catch (err) {
      console.error('Failed to load menu', err);
      this._menuItems.set([]);
    } finally {
      this.isLoading.set(false);
    }
  }

  private mapBackendToFrontend(items: BackendMenuItem[]): MenuItem[] {
    return items.map(item => {
      let route = item.route ?? '';
      if (route === '/' || route === '/app') {
        route = '/home';
      } else if (route === '/profile' || route === '/app/profile') {
        route = '/iam/profile';
      } else if (route.startsWith('/app/')) {
        route = route.replace('/app/', '/');
      }

      return {
        id: item.id,
        label: item.name ?? '',
        icon: item.icon ?? '',
        routerLink: route ? [route] : [],
        items: item.children ? this.mapBackendToFrontend(item.children) : [],
      };
    });
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
