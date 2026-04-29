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
  authority: getAuthAuthority(),
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

function getAuthAuthority(): string {
  try {
    return (window as any).__env?.authAuthority ?? 'https://auth.arda.io.vn';
  } catch {
    return 'https://auth.arda.io.vn';
  }
}

function getClientId(): string {
  try {
    return (window as any).__env?.authClientId ?? '370596460112183382';
  } catch {
    return '370596460112183382';
  }
}

export function getApiUrl(): string {
  try {
    return (window as any).__env?.apiUrl ?? 'http://localhost:8000';
  } catch {
    return 'http://localhost:8000';
  }
}

export function getApiPath(): string {
  try {
    return normalizeApiPath((window as any).__env?.apiPath ?? '/v1');
  } catch {
    return '/v1';
  }
}

function normalizeApiPath(path: string): string {
  const trimmed = (path || '/v1').trim();
  if (!trimmed || trimmed === '/') {
    return '';
  }
  return `/${trimmed.replace(/^\/+|\/+$/g, '')}`;
}
