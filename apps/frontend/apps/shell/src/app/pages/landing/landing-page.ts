import { Component, inject } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { ButtonModule } from 'primeng/button';
import { SelectModule } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { TenantService } from '../../services/tenant.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-landing-page',
  standalone: true,
  imports: [CommonModule, RouterModule, ButtonModule, SelectModule, FormsModule],
  templateUrl: './landing-page.html',
  styleUrl: './landing-page.css',
})
export class LandingPage {
  private tenantService = inject(TenantService);
  authService = inject(AuthService);
  tenants = this.tenantService.tenants;
  selectedTenantId = this.tenantService.selectedTenantId;

  features = [
    {
      icon: 'pi pi-users',
      title: 'CRM & Sales',
      desc: 'Quản lý khách hàng, pipeline bán hàng và lead scoring tự động.',
    },
    {
      icon: 'pi pi-id-card',
      title: 'HRM',
      desc: 'Nhân sự, chấm công, bảng lương tích hợp trên một nền tảng.',
    },
    {
      icon: 'pi pi-chart-bar',
      title: 'Finance',
      desc: 'Kế toán, công nợ, báo cáo tài chính real-time.',
    },
    {
      icon: 'pi pi-objects-column',
      title: 'Micro Frontend',
      desc: 'Kiến trúc module hóa — mở rộng không giới hạn.',
    },
  ];

  onTenantChange(tenantId: string): void {
    this.tenantService.selectTenant(tenantId);
  }
}
