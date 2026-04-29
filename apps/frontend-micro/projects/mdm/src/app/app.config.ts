import { ApplicationConfig, inject, provideAppInitializer, provideBrowserGlobalErrorListeners } from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideHttpClient, withInterceptors, HttpClient } from '@angular/common/http';
import { provideAuth, OidcSecurityService } from 'angular-auth-oidc-client';
import { MessageService } from 'primeng/api';
import { providePrimeNG } from 'primeng/config';
import { provideTranslateService, TranslateLoader } from '@ngx-translate/core';
import { firstValueFrom } from 'rxjs';
import {
  ArdaPreset,
  createArdaPreset,
  DEFAULT_THEME_SETTINGS,
  getAuthConfig,
  HttpLoaderFactory,
  LanguageService,
  PALETTES,
  RADIUS,
  SCALE,
  ThemeSettings,
  apiInterceptor,
} from '@arda/core';

import { routes } from './app.routes';

function getInitialPreset() {
  try {
    const stored = localStorage.getItem('arda-theme-settings');
    if (stored) {
      const settings: ThemeSettings = { ...DEFAULT_THEME_SETTINGS, ...JSON.parse(stored) };
      return createArdaPreset({
        palette: PALETTES[settings.palette as keyof typeof PALETTES],
        radius: settings.radius as keyof typeof RADIUS,
        scale: settings.scale as keyof typeof SCALE,
      });
    }
  } catch {
    /* ignore */
  }
  return ArdaPreset;
}

function initializeAuth(): Promise<void> {
  const oidc = inject(OidcSecurityService);
  return firstValueFrom(oidc.checkAuth()).then(() => undefined);
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
    MessageService,
    provideAppInitializer(initializeAuth),
    provideAppInitializer(initializeLanguage),
  ],
};
