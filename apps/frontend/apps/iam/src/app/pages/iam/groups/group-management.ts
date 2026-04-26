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
import { AvatarModule } from 'primeng/avatar';
import { TooltipModule } from 'primeng/tooltip';
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
    ButtonModule,
    InputTextModule,
    TextareaModule,
    TagModule,
    DialogModule,
    ToastModule,
    ConfirmDialogModule,
    MultiSelectModule,
    SelectModule,
    AvatarModule,
    TooltipModule
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './group-management.html',
})
export class GroupManagement implements OnInit {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);
  private fb = inject(FormBuilder);

  groups = signal<Group[]>([]);
  isLoading = signal(false);
  isSaving = signal(false);

  // Dialog states
  displayDialog = false;
  displayMembersDialog = false;
  displayRolesDialog = false;

  groupForm: FormGroup;
  selectedGroup: Group | null = null;

  // Members management
  groupMembers = signal<User[]>([]);
  allUsers = signal<User[]>([]);
  selectedUserIds = signal<string[]>([]);

  // Roles management
  groupRoles = signal<Role[]>([]);
  allRoles = signal<Role[]>([]);
  selectedRoleIds = signal<string[]>([]);

  constructor() {
    this.groupForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: [''],
    });
  }

  ngOnInit(): void {
    this.loadData();
  }

  async loadData() {
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const data = await this.userService.listGroups(tenantId);
      this.groups.set(data);
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải danh sách nhóm' });
    } finally {
      this.isLoading.set(false);
    }
  }

  openNew() {
    this.selectedGroup = null;
    this.groupForm.reset();
    this.displayDialog = true;
  }

  editGroup(group: Group) {
    this.selectedGroup = group;
    this.groupForm.patchValue({
      name: group.name,
      description: group.description
    });
    this.displayDialog = true;
  }

  async saveGroup() {
    if (this.groupForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    try {
      if (this.selectedGroup) {
        await this.userService.updateGroup(this.selectedGroup.id, this.groupForm.value);
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật nhóm' });
      } else {
        await this.userService.createGroup(this.groupForm.value, tenantId);
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã tạo nhóm mới' });
      }
      this.displayDialog = false;
      this.loadData();
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu nhóm' });
    } finally {
      this.isSaving.set(false);
    }
  }

  deleteGroup(group: Group) {
    this.confirmationService.confirm({
      message: `Bạn có chắc chắn muốn xóa nhóm "${group.name}" không?`,
      header: 'Xác nhận xóa',
      icon: 'pi pi-exclamation-triangle',
      accept: async () => {
        try {
          await this.userService.deleteGroup(group.id);
          this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa nhóm' });
          this.loadData();
        } catch (err) {
          this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa nhóm' });
        }
      }
    });
  }

  // Members Management
  async openMembersDialog(group: Group) {
    this.selectedGroup = group;
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const [members, users] = await Promise.all([
        this.userService.listGroupMembers(group.id),
        this.userService.listUsers(tenantId)
      ]);
      this.groupMembers.set(members);
      this.allUsers.set(users);
      this.selectedUserIds.set(members.map(m => m.id));
      this.displayMembersDialog = true;
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải thông tin thành viên' });
    } finally {
      this.isLoading.set(false);
    }
  }

  async saveMembers() {
    if (!this.selectedGroup) return;
    this.isSaving.set(true);
    try {
      const currentIds = this.groupMembers().map(m => m.id);
      const newIds = this.selectedUserIds();

      const toAdd = newIds.filter(id => !currentIds.includes(id));
      const toRemove = currentIds.filter(id => !newIds.includes(id));

      await Promise.all([
        ...toAdd.map(id => this.userService.addGroupMember(this.selectedGroup!.id, id)),
        ...toRemove.map(id => this.userService.removeGroupMember(this.selectedGroup!.id, id))
      ]);

      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật thành viên nhóm' });
      this.displayMembersDialog = false;
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể cập nhật thành viên' });
    } finally {
      this.isSaving.set(false);
    }
  }

  // Roles Management
  async openRolesDialog(group: Group) {
    this.selectedGroup = group;
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const [gRoles, allRoles] = await Promise.all([
        this.userService.listGroupRoles(group.id),
        this.userService.listRoles(tenantId)
      ]);
      this.groupRoles.set(gRoles);
      this.allRoles.set(allRoles);
      this.selectedRoleIds.set(gRoles.map(r => r.id));
      this.displayRolesDialog = true;
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải thông tin vai trò' });
    } finally {
      this.isLoading.set(false);
    }
  }

  async saveRoles() {
    if (!this.selectedGroup) return;
    this.isSaving.set(true);
    try {
      const currentIds = this.groupRoles().map(r => r.id);
      const newIds = this.selectedRoleIds();

      const toAdd = newIds.filter(id => !currentIds.includes(id));
      const toRemove = currentIds.filter(id => !newIds.includes(id));

      await Promise.all([
        ...toAdd.map(id => this.userService.assignGroupRole(this.selectedGroup!.id, id)),
        ...toRemove.map(id => this.userService.revokeGroupRole(this.selectedGroup!.id, id))
      ]);

      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật vai trò nhóm' });
      this.displayRolesDialog = false;
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể cập nhật vai trò' });
    } finally {
      this.isSaving.set(false);
    }
  }

  getInitials(name: string): string {
    return name?.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2) || '??';
  }
}
