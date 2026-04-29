import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Button } from 'primeng/button';
import { Dialog } from 'primeng/dialog';
import { InputText } from 'primeng/inputtext';
import { Tenant, TenantService } from '../../services/tenant.service';

@Component({
  selector: 'app-workspace-management',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, Button, Dialog, InputText],
  templateUrl: './workspace-management.html',
})
export class WorkspaceManagement {
  private tenantService = inject(TenantService);
  private router = inject(Router);

  readonly tenants = this.tenantService.tenants;
  readonly selectedTenantId = this.tenantService.selectedTenantId;
  readonly isLoading = this.tenantService.isLoading;
  readonly isDialogOpen = signal(false);
  readonly isSaving = signal(false);
  readonly deletingTenantId = signal('');
  readonly selectedTenant = this.tenantService.selectedTenant;
  readonly tenantCount = computed(() => this.tenants().length);
  readonly selectedTenantName = computed(() => this.selectedTenant()?.name ?? 'Chưa chọn');

  readonly form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(2)] }),
    slug: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.pattern(/^[a-z0-9]+(?:-[a-z0-9]+)*$/)] }),
  });

  openCreate(): void {
    this.form.reset({ name: '', slug: '' });
    this.isDialogOpen.set(true);
  }

  closeDialog(): void {
    if (this.isSaving()) return;
    this.isDialogOpen.set(false);
  }

  fillSlugFromName(): void {
    const slugControl = this.form.controls.slug;
    if (slugControl.dirty && slugControl.value) return;
    slugControl.setValue(this.slugify(this.form.controls.name.value));
  }

  selectTenant(tenant: Tenant): void {
    this.tenantService.selectTenant(tenant.id);
    this.router.navigate(['/home']);
  }

  createWorkspace(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const { name, slug } = this.form.getRawValue();
    this.isSaving.set(true);
    this.tenantService.createTenant(name.trim(), slug.trim()).subscribe({
      next: (tenant) => {
        this.isSaving.set(false);
        this.isDialogOpen.set(false);
        this.tenantService.selectTenant(tenant.id);
        this.router.navigate(['/home']);
      },
      error: (err) => {
        console.error('Failed to create workspace', err);
        this.isSaving.set(false);
      },
    });
  }

  deleteWorkspace(tenant: Tenant): void {
    const ok = confirm(`Xóa workspace "${tenant.name}"? Thao tác này sẽ ẩn tenant khỏi danh sách.`);
    if (!ok) return;

    this.deletingTenantId.set(tenant.id);
    this.tenantService.deleteTenant(tenant.id).subscribe({
      next: () => {
        this.deletingTenantId.set('');
      },
      error: (err) => {
        console.error('Failed to delete workspace', err);
        this.deletingTenantId.set('');
      },
    });
  }

  tenantInitials(name: string): string {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return 'WS';
    if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }

  private slugify(value: string): string {
    return value
      .toLowerCase()
      .normalize('NFD')
      .replace(/[\u0300-\u036f]/g, '')
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '');
  }
}
