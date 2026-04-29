import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, map } from 'rxjs';

export interface PasswordPolicy {
  minLength: number;
  requireUppercase: boolean;
  requireLowercase: boolean;
  requireNumber: boolean;
  requireSymbol: boolean;
}

export interface LoginPolicy {
  passwordLoginEnabled: boolean;
  externalIdpEnabled: boolean;
  mfaRequired: boolean;
}

export interface AuthSettings {
  tenantId: string;
  authMode: string;
  provider: string;
  passwordPolicy: PasswordPolicy;
  loginPolicy: LoginPolicy;
}

@Injectable({ providedIn: 'root' })
export class AuthSettingsService {
  private http = inject(HttpClient);

  getAuthSettings(tenantId: string): Observable<AuthSettings> {
    return this.http.get<any>('/api/v1/auth/settings', {
      params: new HttpParams().set('tenant_id', tenantId),
    }).pipe(map(resp => this.toAuthSettings(resp)));
  }

  private toAuthSettings(resp: any): AuthSettings {
    const passwordPolicy = resp.password_policy ?? resp.passwordPolicy ?? {};
    const loginPolicy = resp.login_policy ?? resp.loginPolicy ?? {};

    return {
      tenantId: resp.tenant_id ?? resp.tenantId ?? '',
      authMode: resp.auth_mode ?? resp.authMode ?? 'SHARED_AUTH',
      provider: resp.provider ?? 'ZITADEL',
      passwordPolicy: {
        minLength: Number(passwordPolicy.min_length ?? passwordPolicy.minLength ?? 8),
        requireUppercase: Boolean(passwordPolicy.require_uppercase ?? passwordPolicy.requireUppercase),
        requireLowercase: Boolean(passwordPolicy.require_lowercase ?? passwordPolicy.requireLowercase),
        requireNumber: Boolean(passwordPolicy.require_number ?? passwordPolicy.requireNumber),
        requireSymbol: Boolean(passwordPolicy.require_symbol ?? passwordPolicy.requireSymbol),
      },
      loginPolicy: {
        passwordLoginEnabled: Boolean(loginPolicy.password_login_enabled ?? loginPolicy.passwordLoginEnabled ?? true),
        externalIdpEnabled: Boolean(loginPolicy.external_idp_enabled ?? loginPolicy.externalIdpEnabled),
        mfaRequired: Boolean(loginPolicy.mfa_required ?? loginPolicy.mfaRequired),
      },
    };
  }
}
