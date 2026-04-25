import { LogLevel } from 'angular-auth-oidc-client';

export interface ArdaAuthConfig {
  authority: string;
  redirectUrl: string;
  postLogoutRedirectUri: string;
  clientId: string;
  scope: string;
  responseType: string;
  silentRenew: boolean;
  useRefreshToken: boolean;
  autoUserInfo: boolean;
  renewTimeBeforeTokenExpiresInSeconds: number;
  logLevel: LogLevel;
  secureRoutes: string[];
}

export const getAuthConfig = (): ArdaAuthConfig => ({
  authority: 'https://auth.arda.io.vn',
  redirectUrl: window.location.origin + '/auth/callback',
  postLogoutRedirectUri: window.location.origin,
  clientId: getClientId(),
  scope: 'openid profile email offline_access',
  responseType: 'code',
  silentRenew: false,
  useRefreshToken: true,
  autoUserInfo: true,
  renewTimeBeforeTokenExpiresInSeconds: 30,
  logLevel: LogLevel.Warn,
  secureRoutes: [getApiUrl()],
});

function getClientId(): string {
  try {
    return (window as any).__env?.authClientId ?? 'ZITADEL_SPA_CLIENT_ID';
  } catch {
    return 'ZITADEL_SPA_CLIENT_ID';
  }
}

function getApiUrl(): string {
  try {
    return (window as any).__env?.apiUrl ?? 'http://localhost:8000/api';
  } catch {
    return 'http://localhost:8000/api';
  }
}
