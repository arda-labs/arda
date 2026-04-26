import { ApplicationConfig, provideBrowserGlobalErrorListeners, inject, provideAppInitializer } from '@angular/core';
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
  LanguageService,
  ThemeSettings,
  DEFAULT_THEME_SETTINGS,
  getAuthConfig,
  HttpLoaderFactory,
  apiInterceptor,
  TENANT_PROVIDER,
} from '@arda/core';

import { routes } from './app.routes';
import { TenantService } from './services/tenant.service';

function getInitialPreset() {
  try {
    const stored = localStorage.getItem('arda-theme-settings');
    if (stored) {
      const s: ThemeSettings = { ...DEFAULT_THEME_SETTINGS, ...JSON.parse(stored) };
      return createArdaPreset({
        palette: PALETTES[s.palette as keyof typeof PALETTES],
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
  return firstValueFrom(oidc.checkAuth()).then(() => { });
}

function initializeLanguage(): Promise<void> {
  const langService = inject(LanguageService);
  return langService.init();
}

export const appConfig: ApplicationConfig = {
  providers: [
    provideTranslateService({
      loader: {
        provide: TranslateLoader,
        useFactory: HttpLoaderFactory,
        deps: [HttpClient],
      },
      lang: 'vi',
      fallbackLang: 'vi',
    }),
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes, withComponentInputBinding()),
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
    provideAppInitializer(initializeAuth),
    provideAppInitializer(initializeLanguage),
  ],
};
