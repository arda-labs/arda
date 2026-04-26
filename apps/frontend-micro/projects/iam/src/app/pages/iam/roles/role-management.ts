import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormControl, FormGroup, Validators } from '@angular/forms';
import { rxResource } from '@angular/core/rxjs-interop';
import { of } from 'rxjs';
import { TableModule } from 'primeng/table';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Textarea } from 'primeng/textarea';
import { Tag } from 'primeng/tag';
import { Dialog } from 'primeng/dialog';
import { Toast } from 'primeng/toast';
import { MessageService, ConfirmationService } from 'primeng/api';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { MultiSelect } from 'primeng/multiselect';
import { UserService, Role } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

@Component({
  selector: 'app-role-management',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TableModule,
    Button,
    InputText,
    Textarea,
    Tag,
    Dialog,
    Toast,
    ConfirmDialog,
    MultiSelect,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './role-management.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoleManagement {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  rolesResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listRoles(params.params);
    }
  });

  permissionsResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listPermissions(params.params);
    }
  });

  isSaving = signal(false);

  // Dialog state
  displayDialog = signal(false);
  selectedRole = signal<Role | null>(null);
  selectedPermissions = signal<string[]>([]);

  // Role form
  roleForm = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(3)] }),
    description: new FormControl('', { nonNullable: true })
  });

  openNew() {
    this.selectedRole.set(null);
    this.roleForm.reset({ name: '', description: '' });
    this.selectedPermissions.set([]);
    this.displayDialog.set(true);
  }

  editRole(role: Role) {
    this.selectedRole.set(role);
    this.roleForm.patchValue({
      name: role.name,
      description: role.description ?? ''
    });
    this.selectedPermissions.set(role.permissions ?? []);
    this.displayDialog.set(true);
  }

  saveRole() {
    if (this.roleForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { name, description } = this.roleForm.getRawValue();
    const perms = this.selectedPermissions();
    this.isSaving.set(true);
    const role = this.selectedRole();

    const request$ = role
      ? this.userService.updateRole(role.id, { name, description, permissions: perms }, tenantId)
      : this.userService.createRole({ name, description, permissions: perms }, tenantId);

    request$.subscribe({
      next: () => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: role ? 'Đã cập nhật vai trò' : 'Đã tạo vai trò mới'
        });
        this.displayDialog.set(false);
        this.rolesResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu vai trò' });
        this.isSaving.set(false);
      }
    });
  }

  deleteRole(role: Role) {
    this.confirmationService.confirm({
      message: `Bạn có chắc chắn muốn xóa vai trò "${role.name}" không?`,
      header: 'Xác nhận xóa',
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        const tenantId = this.tenantService.selectedTenantId();
        if (!tenantId) return;

        this.userService.deleteRole(role.id, tenantId).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa vai trò' });
            this.rolesResource.reload();
          },
          error: () => {
            this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa vai trò' });
          }
        });
      },
    });
  }
}
