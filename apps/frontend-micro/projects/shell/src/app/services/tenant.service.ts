import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map, tap, of, switchMap, catchError, finalize } from 'rxjs';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { TenantProvider } from '@arda/core';

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
  private oidc = inject(OidcSecurityService);

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

  fetchTenants(): Observable<Tenant[]> {
    return this.http.get<TenantMembershipResponse>(`/api/v1/me/tenants`).pipe(
      map(resp => (resp.memberships ?? []).map(m => ({
        id: m.tenantId,
        name: m.tenantName,
        slug: m.tenantSlug,
        role: m.role
      })))
    );
  }

  loadTenants(): Observable<void> {
    return this.oidc.isAuthenticated$.pipe(
      switchMap(({ isAuthenticated }) => {
        if (!isAuthenticated) {
          console.log('TenantService: Not authenticated, skipping loadTenants');
          return of(undefined);
        }

        this.isLoading.set(true);
        console.log('TenantService: Loading tenants from /api/v1/me/tenants');

        return this.fetchTenants().pipe(
          tap(loaded => {
            console.log('TenantService: Received tenants', loaded);
            this.tenants.set(loaded);

            // Nếu tenant đã lưu không còn trong list thì chọn cái đầu tiên
            if (this.selectedTenantId() && !loaded.some(t => t.id === this.selectedTenantId())) {
              this.selectTenant(loaded[0]?.id ?? '');
            } else if (!this.selectedTenantId() && loaded.length > 0) {
              this.selectTenant(loaded[0].id);
            }
          }),
          map(() => undefined),
          catchError(err => {
            console.error('Failed to load tenants', err);
            this.tenants.set([]);
            return of(undefined);
          }),
          finalize(() => this.isLoading.set(false))
        );
      })
    );
  }

  selectTenant(id: string): void {
    this.selectedTenantId.set(id);
    try {
      localStorage.setItem(TenantService.STORAGE_KEY, id);
    } catch {
      /* SSR / restricted */
    }
  }

  createTenant(name: string, slug: string): Observable<Tenant> {
    return this.http.post<Tenant>(`/api/v1/tenants`, { name, slug }).pipe(
      tap(() => this.loadTenants().subscribe())
    );
  }

  private loadSaved(): string {
    try {
      return localStorage.getItem(TenantService.STORAGE_KEY) ?? '';
    } catch {
      return '';
    }
  }
}
