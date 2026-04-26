import { Component, inject, computed } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AvatarModule } from 'primeng/avatar';
import { ButtonModule } from 'primeng/button';
import { TranslatePipe } from '@ngx-translate/core';
import { AuthService } from '../../services/auth.service';
import { TenantService } from '../../services/tenant.service';
import { LanguageService } from '@arda-mfe/shared-core';

@Component({
  selector: 'app-select-workspace',
  standalone: true,
  imports: [CommonModule, AvatarModule, ButtonModule, TranslatePipe],
  templateUrl: './select-workspace.html',
  styleUrl: './select-workspace.css',
})
export class SelectWorkspace {
  private router = inject(Router);
  private tenantService = inject(TenantService);
  private langService = inject(LanguageService);

  authService = inject(AuthService);
  tenants = this.tenantService.tenants;
  isLoading = this.tenantService.isLoading;
  currentLang = this.langService.currentLang;

  changeLang(lang: string) {
    this.langService.setLanguage(lang);
  }

  /** Lấy 2 chữ cái đầu của tên tenant */
  tenantInitials(name: string): string {
    const parts = name.trim().split(/\s+/);
    if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }

  select(tenantId: string): void {
    this.tenantService.selectTenant(tenantId);
    this.router.navigate(['/app']);
  }

  async createFirstWorkspace(): Promise<void> {
    try {
      const name = prompt('Nhập tên Workspace mới:');
      if (!name) return;
      const slug = name.toLowerCase().replace(/\s+/g, '-');
      await this.tenantService.createTenant(name, slug);
    } catch (err) {
      console.error('Failed to create tenant', err);
      alert('Không thể tạo workspace. Vui lòng kiểm tra console.');
    }
  }

  logout(): void {
    this.authService.logout();
  }
}
