import { Injectable, signal, computed, inject } from '@angular/core';
import { Injector } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { Router } from '@angular/router';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { map } from 'rxjs/operators';
import { TenantService } from './tenant.service';

export interface AuthUser {
  id: string;
  email: string;
  name: string;
  tenant: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private oidc = inject(OidcSecurityService);
  private tenantService = inject(TenantService);
  private _injector = inject(Injector);

  readonly isAuthenticated = toSignal(
    this.oidc.isAuthenticated$.pipe(map(({ isAuthenticated }) => isAuthenticated)),
    { initialValue: false },
  );

  readonly currentUser = signal<AuthUser | null>(null);

  readonly userInitials = computed(() => {
    const name = this.currentUser()?.name;
    if (!name) return 'U';
    const parts = name.trim().split(/\s+/);
    if (parts.length === 1) return parts[0][0].toUpperCase();
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  });

  constructor() {
    this.oidc.userData$.subscribe(({ userData }) => {
      if (userData) {
        this.currentUser.set({
          id: userData.sub ?? '',
          email: userData.email ?? '',
          name: userData.name ?? userData.preferred_username ?? '',
          tenant: userData.tenant_id ?? '',
        });
      } else {
        this.currentUser.set(null);
      }
    });
  }

  login(): void {
    this.oidc.authorize();
  }

  logout(): void {
    this.oidc.logoff().subscribe(() => {
      const router = this._injector.get(Router);
      router.navigate(['/']);
    });
  }

  getAccessToken(): Promise<string> {
    return this.oidc.getAccessToken().toPromise().then(t => t ?? '');
  }
}
