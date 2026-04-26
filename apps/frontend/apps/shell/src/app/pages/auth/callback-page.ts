import { Component, inject } from '@angular/core';
import { Router } from '@angular/router';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { ProgressSpinnerModule } from 'primeng/progressspinner';
import { TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-callback-page',
  standalone: true,
  imports: [ProgressSpinnerModule],
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
    // Chỉ subscribe vào trạng thái isAuthenticated hiện tại.
    this.oidc.isAuthenticated$.subscribe(async ({ isAuthenticated }) => {
      console.log('Callback: isAuthenticated', isAuthenticated);
      if (!isAuthenticated) {
        // Có thể oidc chưa load xong state từ storage/URL
        return;
      }

      // Load danh sách tenant của user
      await this.tenantService.loadTenants();
      const tenants = this.tenantService.tenants();

      if (tenants.length === 0) {
        // Chưa có workspace → vào app, hiển thị empty state
        this.router.navigate(['/home']);
      } else if (tenants.length === 1) {
        // Chỉ 1 tenant → chọn luôn, vào thẳng app
        this.tenantService.selectTenant(tenants[0].id);
        this.router.navigate(['/home']);
      } else {
        // Nhiều tenant → cho user chọn workspace
        this.router.navigate(['/select-workspace']);
      }
    });
  }
}
