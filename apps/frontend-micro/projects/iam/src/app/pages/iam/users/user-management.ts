import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormControl, FormGroup, Validators } from '@angular/forms';
import { rxResource } from '@angular/core/rxjs-interop';
import { of } from 'rxjs';
import { TableModule } from 'primeng/table';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Dialog } from 'primeng/dialog';
import { Select } from 'primeng/select';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { Avatar } from 'primeng/avatar';
import { UserService, User } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TableModule,
    Button,
    InputText,
    Dialog,
    Select,
    Toast,
    Avatar,
  ],
  providers: [MessageService],
  templateUrl: './user-management.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserManagement {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);

  usersResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listUsers(params.params);
    }
  });

  rolesResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listRoles(params.params);
    }
  });

  isSaving = signal(false);

  // Dialog state
  displayRoleDialog = signal(false);
  displayInviteDialog = signal(false);
  displayCreateDialog = signal(false);
  selectedUser = signal<User | null>(null);

  // Forms (Ổn định nhất 2026)
  selectedRoleControl = new FormControl('', { nonNullable: true, validators: [Validators.required] });

  inviteForm = new FormGroup({
    externalId: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    role: new FormControl('MEMBER', { nonNullable: true, validators: [Validators.required] })
  });

  createForm = new FormGroup({
    email: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.email] }),
    displayName: new FormControl('', { nonNullable: true }),
    password: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(6)] })
  });

  openRoleDialog(user: User) {
    this.selectedUser.set(user);
    this.selectedRoleControl.reset();
    this.displayRoleDialog.set(true);
  }

  openInviteDialog() {
    this.inviteForm.reset({ externalId: '', role: 'MEMBER' });
    this.displayInviteDialog.set(true);
  }

  openCreateDialog() {
    this.createForm.reset({ email: '', displayName: '', password: '' });
    this.displayCreateDialog.set(true);
  }

  saveRoleAssignment() {
    if (this.selectedRoleControl.invalid) return;
    const user = this.selectedUser();
    const roleId = this.selectedRoleControl.value;
    if (!user || !roleId) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    this.isSaving.set(true);
    this.userService.assignRole(user.id, roleId, tenantId).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã cập nhật vai trò' });
        this.displayRoleDialog.set(false);
        this.usersResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gán vai trò' });
        this.isSaving.set(false);
      }
    });
  }

  sendInvite() {
    if (this.inviteForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { externalId, role } = this.inviteForm.getRawValue();
    this.isSaving.set(true);
    this.userService.inviteMember(externalId, role, tenantId).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gửi lời mời thành viên' });
        this.displayInviteDialog.set(false);
        this.usersResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gửi lời mời' });
        this.isSaving.set(false);
      }
    });
  }

  createUser() {
    if (this.createForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { email, displayName, password } = this.createForm.getRawValue();
    this.isSaving.set(true);
    this.userService.createUser(
      { email, display_name: displayName, password },
      tenantId,
    ).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã tạo người dùng mới' });
        this.displayCreateDialog.set(false);
        this.usersResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể tạo người dùng' });
        this.isSaving.set(false);
      }
    });
  }

  getInitials(name: string): string {
    if (!name) return 'U';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2);
  }
}
