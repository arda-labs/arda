import { Injectable, signal, computed, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map, tap, of, switchMap, catchError, finalize, take } from 'rxjs';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { TenantProvider } from '@arda/core';

export interface Tenant {
  id: string;
  name: string;
  slug: string;
  role: string;
  deploymentMode?: TenantDeploymentMode;
  authMode?: TenantAuthMode;
}

export type TenantDeploymentMode = 'SHARED' | 'DEDICATED';
export type TenantAuthMode = 'SHARED_AUTH' | 'DEDICATED_AUTH';

interface TenantMembershipResponse {
  memberships: {
    tenantId: string;
    tenantName: string;
    tenantSlug: string;
    role: string;
    deploymentMode?: TenantDeploymentMode;
    authMode?: TenantAuthMode;
  }[];
}

@Injectable({ providedIn: 'root' })
export class TenantService implements TenantProvider {
  private static STORAGE_KEY = 'arda-selected-tenant';
  private static CHANGE_EVENT = 'arda-selected-tenant-change';
  private http = inject(HttpClient);
  private oidc = inject(OidcSecurityService);

  readonly tenants = signal<Tenant[]>([]);
  readonly selectedTenantId = signal<string>(this.loadSaved());
  readonly isLoading = signal(false);

  readonly selectedTenant = computed(() =>
    this.tenants().find((t) => t.id === this.selectedTenantId()) ?? this.tenants()[0] ?? null,
  );

  constructor() {
    this.listenForExternalTenantChanges();
  }

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
        role: m.role,
        deploymentMode: m.deploymentMode,
        authMode: m.authMode,
      })))
    );
  }

  loadTenants(): Observable<void> {
    return this.oidc.isAuthenticated$.pipe(
      take(1),
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
      window.dispatchEvent(new CustomEvent<string>(TenantService.CHANGE_EVENT, { detail: id }));
    } catch {
      /* SSR / restricted */
    }
  }

  createTenant(
    name: string,
    slug: string,
    options: { deploymentMode?: TenantDeploymentMode; authMode?: TenantAuthMode } = {},
  ): Observable<Tenant> {
    return this.http.post<any>(`/api/v1/tenants`, {
      name,
      slug,
      deployment_mode: options.deploymentMode ?? 'SHARED',
      auth_mode: options.authMode ?? 'SHARED_AUTH',
    }).pipe(
      map(resp => this.toTenant(resp)),
      tap(created => {
        this.selectTenant(created.id);
        this.loadTenants().subscribe();
      })
    );
  }

  deleteTenant(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/tenants/${encodeURIComponent(id)}`).pipe(
      switchMap(() => this.fetchTenants()),
      tap(loaded => {
        this.tenants.set(loaded);
        if (this.selectedTenantId() === id || !loaded.some(t => t.id === this.selectedTenantId())) {
          this.selectTenant(loaded[0]?.id ?? '');
        }
      }),
      map(() => undefined)
    );
  }

  private toTenant(t: any): Tenant {
    return {
      id: t.id,
      name: t.name,
      slug: t.slug,
      role: t.role ?? 'owner',
      deploymentMode: t.deployment_mode ?? t.deploymentMode,
      authMode: t.auth_mode ?? t.authMode,
    };
  }

  private loadSaved(): string {
    try {
      return localStorage.getItem(TenantService.STORAGE_KEY) ?? '';
    } catch {
      return '';
    }
  }

  private listenForExternalTenantChanges(): void {
    if (typeof window === 'undefined') return;

    window.addEventListener('storage', (event) => {
      if (event.key === TenantService.STORAGE_KEY) {
        this.selectedTenantId.set(event.newValue ?? '');
      }
    });

    window.addEventListener(TenantService.CHANGE_EVENT, (event) => {
      const id = (event as CustomEvent<string>).detail;
      if (typeof id === 'string' && id !== this.selectedTenantId()) {
        this.selectedTenantId.set(id);
      }
    });
  }
}
