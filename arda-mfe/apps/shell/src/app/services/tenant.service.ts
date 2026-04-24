import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { getEnv } from '@arda-mfe/shared-config';
import { TenantProvider } from '@arda-mfe/shared-api';

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  role: string;
}

interface TenantMembershipResponse {
  memberships: {
    tenantId: string;
    tenantName: string;
    tenantSlug: string;
    role: string;
  }[];
}

@Injectable({ providedIn: 'root' })
export class TenantService implements TenantProvider {
  private static STORAGE_KEY = 'arda-selected-tenant';
  private http = inject(HttpClient);
  private apiUrl = getEnv().apiUrl;

  readonly tenants = signal<Tenant[]>([]);
  readonly selectedTenantId = signal<string>(this.loadSaved());
  readonly isLoading = signal(false);

  readonly selectedTenant = computed(() =>
    this.tenants().find((t) => t.id === this.selectedTenantId()) ?? this.tenants()[0] ?? null,
  );

  /**
   * Implement TenantProvider interface for authInterceptor
   */
  getTenantId(): string {
    return this.selectedTenantId();
  }

  async loadTenants(): Promise<void> {
    this.isLoading.set(true);
    try {
      console.log('TenantService: Loading tenants from', `${this.apiUrl}/v1/me/tenants`);
      const resp = await firstValueFrom(
        this.http.get<TenantMembershipResponse>(`${this.apiUrl}/v1/me/tenants`)
      );
      console.log('TenantService: Received response', resp);

      const loaded = (resp.memberships ?? []).map(m => ({
        id: m.tenantId,
        name: m.tenantName,
        slug: m.tenantSlug,
        role: m.role
      }));

      this.tenants.set(loaded);

      // Nếu tenant đã lưu không còn trong list thì chọn cái đầu tiên
      if (this.selectedTenantId() && !loaded.some(t => t.id === this.selectedTenantId())) {
        this.selectTenant(loaded[0]?.id ?? '');
      } else if (!this.selectedTenantId() && loaded.length > 0) {
        this.selectTenant(loaded[0].id);
      }
    } catch (err) {
      console.error('Failed to load tenants', err);
      this.tenants.set([]);
    } finally {
      this.isLoading.set(false);
    }
  }

  selectTenant(id: string): void {
    this.selectedTenantId.set(id);
    try {
      localStorage.setItem(TenantService.STORAGE_KEY, id);
    } catch {
      /* SSR / restricted */
    }
  }

  async createTenant(name: string, slug: string): Promise<Tenant> {
    const t = await firstValueFrom(
      this.http.post<Tenant>(`${this.apiUrl}/v1/tenants`, { name, slug })
    );
    await this.loadTenants();
    return t;
  }

  private loadSaved(): string {
    try {
      return localStorage.getItem(TenantService.STORAGE_KEY) ?? '';
    } catch {
      return '';
    }
  }
}
