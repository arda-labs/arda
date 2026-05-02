import { ChangeDetectionStrategy, Component, inject, signal, computed } from '@angular/core';
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
import { InputGroup } from 'primeng/inputgroup';
import { ScrollPanel } from 'primeng/scrollpanel';
import { UserService, Role } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

interface PermissionGroup {
  resource: string;
  label: string;
  permissions: string[];
  selected: boolean;
}

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
    InputGroup,
    ScrollPanel,
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

  // Permission filtering
  searchQuery = signal('');

  filteredGroups = computed(() => {
    const allPerms = this.permissionsResource.value() ?? [];
    const selected = this.selectedPermissions();
    const q = this.searchQuery().toLowerCase();

    // Group by resource
    const grouped = new Map<string, string[]>();
    for (const perm of allPerms) {
      const colon = perm.indexOf(':');
      const resource = colon > 0 ? perm.substring(0, colon) : 'other';
      if (!grouped.has(resource)) grouped.set(resource, []);
      grouped.get(resource)!.push(perm);
    }

    const groups: PermissionGroup[] = [];
    for (const [resource, perms] of grouped) {
      const filtered = q ? perms.filter(p => p.toLowerCase().includes(q)) : perms;
      if (filtered.length === 0) continue;
      groups.push({
        resource,
        label: this.resourceLabel(resource),
        permissions: filtered,
        selected: filtered.every(p => selected.has(p)),
      });
    }
    return groups.sort((a, b) => a.label.localeCompare(b.label));
  });

  // Selected permissions as a Set for O(1) lookup
  selectedPermissions = signal<Set<string>>(new Set());

  // Role form
  roleForm = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(3)] }),
    description: new FormControl('', { nonNullable: true })
  });

  private resourceLabel(resource: string): string {
    const labels: Record<string, string> = {
      user: 'Người dùng',
      role: 'Vai trò',
      'user-group': 'Nhóm người dùng',
      menu: 'Menu',
      tenant: 'Workspace',
      dashboard: 'Dashboard',
      system: 'Hệ thống',
      member: 'Thành viên',
      permission: 'Quyền hạn',
      approval: 'Phê duyệt',
      audit: 'Nhật ký',
      settings: 'Cấu hình',
      lead: 'Khách hàng tiềm năng',
      deal: 'Cơ hội',
      contact: 'Liên hệ',
      employee: 'Nhân viên',
      payroll: 'Lương',
      attendance: 'Chấm công',
      invoice: 'Hóa đơn',
      expense: 'Chi phí',
      report: 'Báo cáo',
      quote: 'Báo giá',
      crm: 'CRM',
      hrm: 'HRM',
      finance: 'Tài chính',
      mdm: 'MDM',
      ntf: 'Thông báo',
      me: 'Cá nhân',
      public: 'Công khai',
      iam: 'IAM',
      '': 'Khác',
    };
    return labels[resource] || resource.charAt(0).toUpperCase() + resource.slice(1);
  }

  openNew() {
    this.selectedRole.set(null);
    this.roleForm.reset({ name: '', description: '' });
    this.selectedPermissions.set(new Set());
    this.searchQuery.set('');
    this.displayDialog.set(true);
  }

  editRole(role: Role) {
    this.selectedRole.set(role);
    this.roleForm.patchValue({
      name: role.name,
      description: role.description ?? ''
    });
    this.selectedPermissions.set(new Set(role.permissions ?? []));
    this.searchQuery.set('');
    this.displayDialog.set(true);
  }

  togglePermission(perm: string) {
    this.selectedPermissions.update(s => {
      const next = new Set(s);
      if (next.has(perm)) next.delete(perm);
      else next.add(perm);
      return next;
    });
  }

  toggleGroup(group: PermissionGroup) {
    const allSelected = group.permissions.every(p => this.selectedPermissions().has(p));
    this.selectedPermissions.update(s => {
      const next = new Set(s);
      for (const perm of group.permissions) {
        if (allSelected) next.delete(perm);
        else next.add(perm);
      }
      return next;
    });
  }

  saveRole() {
    if (this.roleForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { name, description } = this.roleForm.getRawValue();
    const perms = [...this.selectedPermissions()];
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
