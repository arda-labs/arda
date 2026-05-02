import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { ProgressSpinner } from 'primeng/progressspinner';
import { TenantService } from '../../services/tenant.service';
import { EMPTY, from, map, of, switchMap, take, firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-callback-page',
  standalone: true,
  imports: [ProgressSpinner],
  template: `
    <div class="min-h-screen flex flex-col items-center justify-center gap-4 bg-surface-50 dark:bg-surface-950">
      <p-progress-spinner styleClass="w-12! h-12!" />
      <p class="text-surface-400 text-sm">Đang xác thực...</p>
    </div>
  `,
})
export class CallbackPage {
  private oidc = inject(OidcSecurityService);
  private router = inject(Router);
  private tenantService = inject(TenantService);

  constructor() {
    this.handleCallback();
  }

  private handleCallback(): void {
    this.oidc.isAuthenticated$.pipe(
      take(1),
      switchMap(({ isAuthenticated }) => {
        if (!isAuthenticated) {
          this.router.navigate(['/login']);
          return EMPTY;
        }
        return this.tenantService.loadTenants();
      }),
    ).subscribe(() => this.navigateAfterAuth());
  }

  private navigateAfterAuth(): void {
    const tenants = this.tenantService.tenants();

    if (tenants.length === 0) {
      this.router.navigate(['/home']);
    } else if (tenants.length === 1) {
      // selectTenant sẽ trigger effect trong PermissionService/MenuService tự động load
      this.tenantService.selectTenant(tenants[0].id);
      this.router.navigate(['/home']);
    } else {
      this.router.navigate(['/select-workspace']);
    }
  }
}
