import { ChangeDetectionStrategy, Component, computed, effect, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Button } from 'primeng/button';
import { Tag } from 'primeng/tag';
import { AuthSettings, AuthSettingsService } from '../../../services/auth-settings.service';
import { TenantService } from '../../../services/tenant.service';

@Component({
  selector: 'app-auth-settings',
  standalone: true,
  imports: [CommonModule, Button, Tag],
  templateUrl: './auth-settings.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthSettingsPage {
  private tenantService = inject(TenantService);
  private authSettingsService = inject(AuthSettingsService);

  settings = signal<AuthSettings | null>(null);
  loading = signal(false);
  error = signal('');

  selectedTenant = computed(() => this.tenantService.selectedTenant());
  selectedTenantName = computed(() => this.tenantService.selectedTenant()?.name || 'hiện tại');
  selectedTenantSlug = computed(() => this.tenantService.selectedTenant()?.slug || '');

  passwordRules = computed(() => {
    const policy = this.settings()?.passwordPolicy;
    if (!policy) return [];
    return [
      `Tối thiểu ${policy.minLength} ký tự`,
      policy.requireUppercase ? 'Bắt buộc chữ hoa' : 'Không bắt buộc chữ hoa',
      policy.requireLowercase ? 'Bắt buộc chữ thường' : 'Không bắt buộc chữ thường',
      policy.requireNumber ? 'Bắt buộc chữ số' : 'Không bắt buộc chữ số',
      policy.requireSymbol ? 'Bắt buộc ký tự đặc biệt' : 'Không bắt buộc ký tự đặc biệt',
    ];
  });

  constructor() {
    effect(() => {
      const tenantId = this.tenantService.selectedTenantId();
      if (tenantId) {
        this.load(tenantId);
      }
    });
  }

  refresh() {
    const tenantId = this.tenantService.selectedTenantId();
    if (tenantId) {
      this.load(tenantId);
    }
  }

  private load(tenantId: string) {
    this.loading.set(true);
    this.error.set('');
    this.authSettingsService.getAuthSettings(tenantId).subscribe({
      next: settings => {
        this.settings.set(settings);
        this.loading.set(false);
      },
      error: err => {
        this.error.set(err?.error?.message || 'Không thể tải cấu hình xác thực');
        this.loading.set(false);
      },
    });
  }
}
