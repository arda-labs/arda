import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { SelectModule } from 'primeng/select';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { AvatarModule } from 'primeng/avatar';
import { UserService, User, Role } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TableModule,
    ButtonModule,
    InputTextModule,
    TagModule,
    DialogModule,
    SelectModule,
    ToastModule,
    AvatarModule
  ],
  providers: [MessageService],
  templateUrl: './user-management.html',
})
export class UserManagement implements OnInit {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);

  users = signal<User[]>([]);
  roles = signal<Role[]>([]);
  isLoading = signal(false);
  
  // Dialog state
  displayRoleDialog = false;
  displayInviteDialog = false;
  displayCreateDialog = false;
  selectedUser: User | null = null;
  selectedRoleId: string = '';
  inviteExternalId: string = '';
  inviteRole: string = 'MEMBER';

  newUser = {
    email: '',
    display_name: '',
    password: ''
  };

  isSaving = signal(false);

  ngOnInit(): void {
    this.loadData();
  }

  async loadData() {
    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isLoading.set(true);
    try {
      const [usersData, rolesData] = await Promise.all([
        this.userService.listUsers(tenantId),
        this.userService.listRoles(tenantId)
      ]);
      this.users.set(usersData);
      this.roles.set(rolesData);
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tải danh sách người dùng' });
    } finally {
      this.isLoading.set(false);
    }
  }

  openRoleDialog(user: User) {
    this.selectedUser = user;
    this.displayRoleDialog = true;
  }

  openInviteDialog() {
    this.inviteExternalId = '';
    this.inviteRole = 'MEMBER';
    this.displayInviteDialog = true;
  }

  openCreateDialog() {
    this.newUser = { email: '', display_name: '', password: '' };
    this.displayCreateDialog = true;
  }

  async saveRoleAssignment() {
    if (!this.selectedUser || !this.selectedRoleId) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    try {
      await this.userService.assignRole(this.selectedUser.id, this.selectedRoleId, tenantId);
      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật vai trò' });
      this.displayRoleDialog = false;
      this.loadData(); // Reload list
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gán vai trò' });
    } finally {
      this.isSaving.set(false);
    }
  }

  async sendInvite() {
    if (!this.inviteExternalId) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    try {
      await this.userService.inviteMember(this.inviteExternalId, this.inviteRole, tenantId);
      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gửi lời mời thành viên' });
      this.displayInviteDialog = false;
      this.loadData();
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gửi lời mời' });
    } finally {
      this.isSaving.set(false);
    }
  }

  async createUser() {
    if (!this.newUser.email || !this.newUser.password) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    try {
      await this.userService.createUser(this.newUser, tenantId);
      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã tạo người dùng mới' });
      this.displayCreateDialog = false;
      this.loadData();
    } catch (err) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tạo người dùng' });
    } finally {
      this.isSaving.set(false);
    }
  }

  getInitials(name: string): string {
    if (!name) return 'U';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2);
  }
}
