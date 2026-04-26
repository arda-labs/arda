import {
  ApplicationConfig,
  provideAppInitializer,
  inject,
  provideBrowserGlobalErrorListeners,
} from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { providePrimeNG } from 'primeng/config';
import { provideHttpClient, withInterceptors, HttpClient } from '@angular/common/http';
import { provideAuth, OidcSecurityService } from 'angular-auth-oidc-client';
import { MessageService } from 'primeng/api';
import { provideTranslateService, TranslateLoader } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';
import {
  ArdaPreset,
  createArdaPreset,
  PALETTES,
  RADIUS,
  SCALE,
} from '@arda-mfe/shared-theme';
import {
  LanguageService,
  ThemeSettings,
  DEFAULT_THEME_SETTINGS,
  getAuthConfig,
  HttpLoaderFactory,
} from '@arda-mfe/shared-core';
import { appRoutes } from './app.routes';
import { apiInterceptor, TENANT_PROVIDER } from '@arda-mfe/shared-api';
import { TenantService } from './services/tenant.service';
import { PermissionService } from './services/permission.service';
import { MenuService } from './services/menu.service';

function getInitialPreset() {
  try {
    const stored = localStorage.getItem('arda-theme-settings');
    if (stored) {
      const s: ThemeSettings = { ...DEFAULT_THEME_SETTINGS, ...JSON.parse(stored) };
      return createArdaPreset({
        palette: PALETTES[s.palette],
        radius: s.radius as keyof typeof RADIUS,
        scale: s.scale as keyof typeof SCALE,
      });
    }
  } catch {
    /* ignore */
  }
  return ArdaPreset;
}

function initializeAuth(): Promise<void> {
  const oidc = inject(OidcSecurityService);
  return firstValueFrom(oidc.checkAuth()).then(() => { /* auth check complete */ });
}

function initializeLanguage(): Promise<void> {
  const langService = inject(LanguageService);
  return langService.init();
}

function initializeTenants(): Promise<void> {
  const tenantService = inject(TenantService);
  const oidc = inject(OidcSecurityService);

  return firstValueFrom(oidc.isAuthenticated$).then(({ isAuthenticated }) => {
    if (isAuthenticated) {
      return tenantService.loadTenants();
    }
    return Promise.resolve();
  });
}

function initializePermissions(): Promise<void> {
  const permService = inject(PermissionService);
  const oidc = inject(OidcSecurityService);

  return firstValueFrom(oidc.isAuthenticated$).then(({ isAuthenticated }) => {
    if (isAuthenticated) {
      return permService.loadPermissions();
    }
    return Promise.resolve();
  });
}

function initializeMenu(): Promise<void> {
  const menuService = inject(MenuService);
  const tenantService = inject(TenantService);
  const oidc = inject(OidcSecurityService);

  return firstValueFrom(oidc.isAuthenticated$).then(({ isAuthenticated }) => {
    if (isAuthenticated) {
      const tenantId = tenantService.selectedTenantId();
      if (tenantId) {
        return menuService.loadMenu();
      }
    }
    return Promise.resolve();
  });
}

// Trigger CI/CD: fix premature API calls and update GitOps workflow
export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(appRoutes, withComponentInputBinding()),
    provideAuth({ config: getAuthConfig() }),
    provideAnimationsAsync(),
    providePrimeNG({
      theme: {
        preset: getInitialPreset(),
        options: {
          darkModeSelector: '.dark',
        },
      },
    }),
    provideHttpClient(withInterceptors([apiInterceptor])),
    { provide: TENANT_PROVIDER, useExisting: TenantService },
    MessageService,
    provideTranslateService({
      loader: {
        provide: TranslateLoader,
        useFactory: HttpLoaderFactory,
        deps: [HttpClient],
      },
      lang: 'vi',
      fallbackLang: 'vi',
    }),
    provideAppInitializer(initializeAuth),
    provideAppInitializer(initializeLanguage),
    provideAppInitializer(initializeTenants),
    provideAppInitializer(initializePermissions),
    provideAppInitializer(initializeMenu),
  ],
};
// Trigger CI/CD refresh
