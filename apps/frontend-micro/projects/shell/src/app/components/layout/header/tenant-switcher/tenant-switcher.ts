import { Component, inject, ChangeDetectionStrategy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Select } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { TenantService } from '../../../../services/tenant.service';

@Component({
  selector: 'app-tenant-switcher',
  standalone: true,
  imports: [CommonModule, Select, FormsModule],
  templateUrl: './tenant-switcher.html',
  styleUrl: './tenant-switcher.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TenantSwitcher {
  private tenantService = inject(TenantService);

  tenants = this.tenantService.tenants;
  selectedId = this.tenantService.selectedTenantId;
  isLoading = this.tenantService.isLoading;

  onTenantChange(event: { value: string }) {
    this.tenantService.selectTenant(event.value);
  }
}
