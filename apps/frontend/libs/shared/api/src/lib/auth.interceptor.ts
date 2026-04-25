import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject, InjectionToken } from '@angular/core';
import { Router } from '@angular/router';
import { MessageService } from 'primeng/api';
import { TranslateService } from '@ngx-translate/core';
import { catchError } from 'rxjs/operators';
import { throwError } from 'rxjs';
import { getEnv } from '@arda-mfe/shared-config';

export interface TenantProvider {
  getTenantId(): string;
}

export const TENANT_PROVIDER = new InjectionToken<TenantProvider>('TENANT_PROVIDER');

function generateRequestId(): string {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
}

export const apiInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const messageService = inject(MessageService);
  const translate = inject(TranslateService);
  const tenantProvider = inject(TENANT_PROVIDER, { optional: true });
  const requestId = generateRequestId();
  const env = getEnv();

  const isApiRequest = req.url.includes('/api/');
  const headers: Record<string, string> = {};

  if (isApiRequest) {
    headers['X-Request-ID'] = requestId;
    headers['Accept-Language'] = translate.currentLang || 'vi';

    if (tenantProvider) {
      const tenantId = tenantProvider.getTenantId();
      if (tenantId) {
        headers['X-Tenant-ID'] = tenantId;
      }
    }
  }

  const authReq = Object.keys(headers).length ? req.clone({ setHeaders: headers }) : req;

  return next(authReq).pipe(
    catchError((error: HttpErrorResponse) => {
      let errorMsg = translate.instant('COMMON.ERROR.SYSTEM');
      let errorSummary = translate.instant('COMMON.ERROR.SYSTEM');

      if (error.error && typeof error.error === 'object') {
        const beError = error.error;
        if (beError.message) {
          errorMsg = beError.message;
        }
        if (beError.reason) {
          errorSummary = `Lỗi: ${beError.reason}`;
        }
      }

      if (error.status === 401) {
        messageService.add({
          severity: 'warn',
          summary: translate.instant('COMMON.ERROR.UNAUTHORIZED'),
          detail: 'Vui lòng đăng nhập lại.',
          life: 3000
        });
        router.navigate(['/login']);
      } else if (error.status === 403) {
        messageService.add({
          severity: 'error',
          summary: translate.instant('COMMON.ERROR.FORBIDDEN'),
          detail: errorMsg || 'Bạn không có quyền thực hiện hành động này.',
          life: 5000
        });
        router.navigate(['/403']);
      } else if (error.status === 404 && req.url.includes('/api/')) {
        messageService.add({
          severity: 'error',
          summary: translate.instant('COMMON.ERROR.NOT_FOUND'),
          detail: `Tài nguyên không tồn tại. (Request ID: ${requestId})`,
          life: 5000
        });
      } else if (error.status >= 500 || error.status === 0) {
        messageService.add({
          severity: 'error',
          summary: errorSummary,
          detail: `${errorMsg} (ID: ${requestId})`,
          life: 7000
        });
      }

      console.error(`[API Error] Request ID: ${requestId}`, error);
      return throwError(() => error);
    })
  );
};

/** @deprecated Use `apiInterceptor` instead. Kept for backward compat re-exports. */
export const authInterceptor = apiInterceptor;
