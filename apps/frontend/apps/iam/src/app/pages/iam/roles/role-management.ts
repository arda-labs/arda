import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { TextareaModule } from 'primeng/textarea';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { ToastModule } from 'primeng/toast';
import { MessageService, ConfirmationService } from 'primeng/api';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { MultiSelectModule } from 'primeng/multiselect';
import { SelectModule } from 'primeng/select';
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
    ButtonModule,
    InputTextModule,
    TextareaModule,
    TagModule,
    DialogModule,
    ToastModule,
    ConfirmDialogModule,
    MultiSelectModule,
    SelectModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './role-management.html',
})
export class RoleManagement implements OnInit {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);
  private fb = inject(FormBuilder);

  roles = signal<Role[]>([]);
  permissions = signal<string[]>([]);
  isLoading = signal(false);
  isSaving = signal(false);

  // Dialog state
  displayDialog = false;
  roleForm: FormGroup;
  selectedRole: Role | null = null;

  constructor() {
    this.roleForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: [''],
      permissions: [[]]
    });
  }

  ngOnInit(): void {
    this.loadData();
    this.loadPermissions();
  }

  async loadData() {
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const rolesData = await this.userService.listRoles(tenantId);
      this.roles.set(rolesData);
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải danh sách vai trò' });
    } finally {
      this.isLoading.set(false);
    }
  }

  async loadPermissions() {
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    try {
      const perms = await this.userService.listPermissions(tenantId);
      this.permissions.set(perms);
    } catch (err) {
      console.error('Failed to load permissions', err);
    }
  }

  openNew() {
    this.selectedRole = null;
    this.roleForm.reset({
      name: '',
      description: '',
      permissions: []
    });
    this.displayDialog = true;
  }

  editRole(role: Role) {
    this.selectedRole = role;
    this.roleForm.patchValue({
      name: role.name,
      description: role.description,
      permissions: role.permissions || []
    });
    this.displayDialog = true;
  }

  async saveRole() {
    if (this.roleForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    try {
      if (this.selectedRole) {
        await this.userService.updateRole(this.selectedRole.id, this.roleForm.value, tenantId);
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật vai trò' });
      } else {
        await this.userService.createRole(this.roleForm.value, tenantId);
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã tạo vai trò mới' });
      }
      this.displayDialog = false;
      this.loadData();
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu vai trò' });
    } finally {
      this.isSaving.set(false);
    }
  }

  deleteRole(role: Role) {
    this.confirmationService.confirm({
      message: `Bạn có chắc chắn muốn xóa vai trò "${role.name}" không?`,
      header: 'Xác nhận xóa',
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: async () => {
        const tenantId = this.tenantService.selectedTenantId();
        if (!tenantId) return;

        try {
          await this.userService.deleteRole(role.id, tenantId);
          this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa vai trò' });
          this.loadData();
        } catch (err) {
          this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa vai trò' });
        }
      }
    });
  }
}
