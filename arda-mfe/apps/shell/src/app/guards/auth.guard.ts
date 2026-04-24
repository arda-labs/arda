import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { OidcSecurityService } from 'angular-auth-oidc-client';
import { map, take } from 'rxjs/operators';
import { TenantService } from '../services/tenant.service';

export const authGuard: CanActivateFn = (route) => {
  const oidc = inject(OidcSecurityService);
  const router = inject(Router);
  const tenantService = inject(TenantService);

  return oidc.isAuthenticated$.pipe(
    map(({ isAuthenticated }) => {
      if (!isAuthenticated) {
        console.log('AuthGuard: Not authenticated, redirecting to login');
        router.navigate(['/login']);
        return false;
      }

      // Nếu đang vào /app mà chưa có tenant được chọn → /select-workspace
      // (trường hợp F5 trực tiếp vào /app sau khi đã login)
      const targetPath = route.routeConfig?.path;
      if (targetPath === 'app' || targetPath === '') {
        const tenants = tenantService.tenants();
        const selected = tenantService.selectedTenantId();

        // tenants đã load (> 0) nhưng không có selection → chọn workspace
        if (tenants.length > 1 && !selected) {
          router.navigate(['/select-workspace']);
          return false;
        }
      }

      return true;
    }),
  );
};
