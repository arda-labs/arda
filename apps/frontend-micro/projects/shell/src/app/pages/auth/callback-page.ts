import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { ProgressSpinner } from 'primeng/progressspinner';
import { TenantService } from '../../services/tenant.service';
import { filter, switchMap } from 'rxjs';

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
      filter(({ isAuthenticated }) => isAuthenticated),
      switchMap(() => this.tenantService.loadTenants())
    ).subscribe(() => {
      const tenants = this.tenantService.tenants();

      if (tenants.length === 0) {
        this.router.navigate(['/home']);
      } else if (tenants.length === 1) {
        this.tenantService.selectTenant(tenants[0].id);
        this.router.navigate(['/home']);
      } else {
        this.router.navigate(['/select-workspace']);
      }
    });
  }
}
