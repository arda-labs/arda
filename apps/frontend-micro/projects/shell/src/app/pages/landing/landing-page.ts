import { Component, inject } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { Button } from 'primeng/button';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { TenantService } from '../../services/tenant.service';
import { AuthService } from '../../services/auth.service';
import { ThemeService, LanguageService } from '@arda/core';

@Component({
  selector: 'app-landing-page',
  standalone: true,
  imports: [CommonModule, RouterModule, Button, Select, FormsModule, TranslateModule],
  templateUrl: './landing-page.html',
  styleUrl: './landing-page.css',
})
export class LandingPage {
  private tenantService = inject(TenantService);
  private themeService = inject(ThemeService);
  private langService = inject(LanguageService);

  authService = inject(AuthService);
  tenants = this.tenantService.tenants;
  selectedTenantId = this.tenantService.selectedTenantId;
  currentLang = this.langService.currentLang;
  settings = this.themeService.settings;

  langOptions = [
    { label: 'Tiếng Việt', value: 'vi' },
    { label: 'English', value: 'en' },
  ];

  features = [
    {
      icon: 'pi pi-users',
      title: 'PAGES.LANDING.FEATURES.CRM.TITLE',
      desc: 'PAGES.LANDING.FEATURES.CRM.DESC',
    },
    {
      icon: 'pi pi-id-card',
      title: 'PAGES.LANDING.FEATURES.HRM.TITLE',
      desc: 'PAGES.LANDING.FEATURES.HRM.DESC',
    },
    {
      icon: 'pi pi-chart-bar',
      title: 'PAGES.LANDING.FEATURES.FINANCE.TITLE',
      desc: 'PAGES.LANDING.FEATURES.FINANCE.DESC',
    },
    {
      icon: 'pi pi-objects-column',
      title: 'PAGES.LANDING.FEATURES.MFE.TITLE',
      desc: 'PAGES.LANDING.FEATURES.MFE.DESC',
    },
  ];

  onTenantChange(tenantId: string): void {
    this.tenantService.selectTenant(tenantId);
  }

  onLanguageChange(lang: string): void {
    this.langService.setLanguage(lang);
  }

  toggleDarkMode(): void {
    this.themeService.updateSetting('darkMode', !this.settings().darkMode);
  }
}
