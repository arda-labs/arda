import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SelectModule } from 'primeng/select';
import { FormsModule } from '@angular/forms';
import { TenantService } from '../../../../services/tenant.service';

@Component({
  selector: 'app-tenant-switcher',
  standalone: true,
  imports: [CommonModule, SelectModule, FormsModule],
  templateUrl: './tenant-switcher.html',
  styleUrl: './tenant-switcher.css',
})
export class TenantSwitcher {
  private tenantService = inject(TenantService);

  tenants = this.tenantService.tenants;
  selectedId = this.tenantService.selectedTenantId;
  isLoading = this.tenantService.isLoading;

  onTenantChange(event: any) {
    this.tenantService.selectTenant(event.value);
  }
}
