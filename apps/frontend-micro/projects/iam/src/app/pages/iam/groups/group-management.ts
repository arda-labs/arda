import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormControl, FormGroup, Validators } from '@angular/forms';
import { rxResource } from '@angular/core/rxjs-interop';
import { of, forkJoin } from 'rxjs';
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
import { Avatar } from 'primeng/avatar';
import { Tooltip } from 'primeng/tooltip';
import { UserService, Group, User, Role } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

@Component({
  selector: 'app-group-management',
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
    Avatar,
    Tooltip,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './group-management.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GroupManagement {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  groupsResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listGroups(params.params);
    }
  });

  isSaving = signal(false);
  isLoadingDetails = signal(false);

  // Dialog states
  displayDialog = signal(false);
  displayMembersDialog = signal(false);
  displayRolesDialog = signal(false);

  // Group form
  selectedGroup = signal<Group | null>(null);
  groupForm = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(3)] }),
    description: new FormControl('', { nonNullable: true })
  });

  // Members management
  groupMembers = signal<User[]>([]);
  allUsers = signal<User[]>([]);
  selectedUserIds = signal<string[]>([]);

  // Roles management
  groupRoles = signal<Role[]>([]);
  allRoles = signal<Role[]>([]);
  selectedRoleIds = signal<string[]>([]);

  openNew() {
    this.selectedGroup.set(null);
    this.groupForm.reset({ name: '', description: '' });
    this.displayDialog.set(true);
  }

  editGroup(group: Group) {
    this.selectedGroup.set(group);
    this.groupForm.patchValue({
      name: group.name,
      description: group.description ?? ''
    });
    this.displayDialog.set(true);
  }

  saveGroup() {
    if (this.groupForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { name, description } = this.groupForm.getRawValue();
    this.isSaving.set(true);
    const group = this.selectedGroup();

    const request$ = group
      ? this.userService.updateGroup(group.id, { name, description })
      : this.userService.createGroup({ name, description }, tenantId);

    request$.subscribe({
      next: () => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: group ? 'Đã cập nhật nhóm' : 'Đã tạo nhóm mới'
        });
        this.displayDialog.set(false);
        this.groupsResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu nhóm' });
        this.isSaving.set(false);
      }
    });
  }

  deleteGroup(group: Group) {
    this.confirmationService.confirm({
      message: `Bạn có chắc chắn muốn xóa nhóm "${group.name}" không?`,
      header: 'Xác nhận xóa',
      icon: 'pi pi-exclamation-triangle',
      accept: () => {
        this.userService.deleteGroup(group.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa nhóm' });
            this.groupsResource.reload();
          },
          error: () => {
            this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa nhóm' });
          }
        });
      },
    });
  }

  // Members Management
  openMembersDialog(group: Group) {
    this.selectedGroup.set(group);
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoadingDetails.set(true);
    forkJoin([
      this.userService.listGroupMembers(group.id),
      this.userService.listUsers(tenantId),
    ]).subscribe({
      next: ([members, users]) => {
        this.groupMembers.set(members);
        this.allUsers.set(users);
        this.selectedUserIds.set(members.map(m => m.id));
        this.displayMembersDialog.set(true);
        this.isLoadingDetails.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải thông tin thành viên' });
        this.isLoadingDetails.set(false);
      }
    });
  }

  saveMembers() {
    const group = this.selectedGroup();
    if (!group) return;

    this.isSaving.set(true);
    const currentIds = this.groupMembers().map(m => m.id);
    const newIds = this.selectedUserIds();

    const toAdd = newIds.filter(id => !currentIds.includes(id));
    const toRemove = currentIds.filter(id => !newIds.includes(id));

    const tasks = [
      ...toAdd.map(id => this.userService.addGroupMember(group.id, id)),
      ...toRemove.map(id => this.userService.removeGroupMember(group.id, id)),
    ];

    if (tasks.length === 0) {
      this.displayMembersDialog.set(false);
      this.isSaving.set(false);
      return;
    }

    forkJoin(tasks).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật thành viên nhóm' });
        this.displayMembersDialog.set(false);
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể cập nhật thành viên' });
        this.isSaving.set(false);
      }
    });
  }

  // Roles Management
  openRolesDialog(group: Group) {
    this.selectedGroup.set(group);
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoadingDetails.set(true);
    forkJoin([
      this.userService.listGroupRoles(group.id),
      this.userService.listRoles(tenantId),
    ]).subscribe({
      next: ([gRoles, allRoles]) => {
        this.groupRoles.set(gRoles);
        this.allRoles.set(allRoles);
        this.selectedRoleIds.set(gRoles.map(r => r.id));
        this.displayRolesDialog.set(true);
        this.isLoadingDetails.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải thông tin vai trò' });
        this.isLoadingDetails.set(false);
      }
    });
  }

  saveRoles() {
    const group = this.selectedGroup();
    if (!group) return;

    this.isSaving.set(true);
    const currentIds = this.groupRoles().map(r => r.id);
    const newIds = this.selectedRoleIds();

    const toAdd = newIds.filter(id => !currentIds.includes(id));
    const toRemove = currentIds.filter(id => !newIds.includes(id));

    const tasks = [
      ...toAdd.map(id => this.userService.assignGroupRole(group.id, id)),
      ...toRemove.map(id => this.userService.revokeGroupRole(group.id, id)),
    ];

    if (tasks.length === 0) {
      this.displayRolesDialog.set(false);
      this.isSaving.set(false);
      return;
    }

    forkJoin(tasks).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật vai trò nhóm' });
        this.displayRolesDialog.set(false);
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể cập nhật vai trò' });
        this.isSaving.set(false);
      }
    });
  }

  getInitials(name: string): string {
    return name?.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2) || '??';
  }
}
