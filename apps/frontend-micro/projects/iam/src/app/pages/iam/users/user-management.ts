import { ChangeDetectionStrategy, Component, computed, effect, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TableModule } from 'primeng/table';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Dialog } from 'primeng/dialog';
import { Select } from 'primeng/select';
import { Toast } from 'primeng/toast';
import { ConfirmationService, MessageService } from 'primeng/api';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { Avatar } from 'primeng/avatar';
import { createPagedResource } from '@arda/core';
import { UserService, User } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';
import { AuthSettingsService, PasswordPolicy } from '../../../services/auth-settings.service';
import { ArdaDataTable } from '../../../shared/table/arda-data-table';

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
    ConfirmDialog,
    Avatar,
    ArdaDataTable,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './user-management.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserManagement {
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private authSettingsService = inject(AuthSettingsService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);
  private router = inject(Router);

  usersTable = createPagedResource<User, string>({
    defaultPageSize: 10,
    rowsPerPageOptions: [10, 20, 50],
    params: () => this.tenantService.selectedTenantId(),
    load: (tenantId, page) => this.userService.listUsersPage(tenantId, page),
  });

  isSaving = signal(false);
  deletingUserId = signal('');
  authSettingsLoading = signal(false);
  private loadedAuthSettingsTenantId = '';
  passwordValue = signal('');
  passwordPolicy = signal<PasswordPolicy>({
    minLength: 8,
    requireUppercase: true,
    requireLowercase: true,
    requireNumber: true,
    requireSymbol: true,
  });

  // Dialog state
  displayInviteDialog = signal(false);
  displayCreateDialog = signal(false);

  // Forms (Ổn định nhất 2026)
  inviteForm = new FormGroup({
    externalId: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    username: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    displayName: new FormControl('', { nonNullable: true }),
    role: new FormControl('MEMBER', { nonNullable: true, validators: [Validators.required] })
  });

  createForm = new FormGroup({
    username: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    email: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.email] }),
    displayName: new FormControl('', { nonNullable: true }),
    password: new FormControl('', { nonNullable: true, validators: [Validators.required] })
  });

  passwordRules = computed(() => {
    const password = this.passwordValue();
    const policy = this.passwordPolicy();
    return [
      {
        label: `Tối thiểu ${policy.minLength} ký tự`,
        valid: password.length >= policy.minLength,
      },
      {
        label: 'Có chữ hoa',
        valid: !policy.requireUppercase || /[A-Z]/.test(password),
      },
      {
        label: 'Có chữ thường',
        valid: !policy.requireLowercase || /[a-z]/.test(password),
      },
      {
        label: 'Có chữ số',
        valid: !policy.requireNumber || /\d/.test(password),
      },
      {
        label: 'Có ký tự đặc biệt',
        valid: !policy.requireSymbol || /[^A-Za-z0-9]/.test(password),
      },
    ];
  });

  constructor() {
    this.createForm.controls.password.valueChanges.subscribe(value => this.passwordValue.set(value ?? ''));

    effect(() => {
      const tenantId = this.tenantService.selectedTenantId();
      if (!tenantId || tenantId === this.loadedAuthSettingsTenantId) return;
      this.loadAuthSettings(tenantId);
    });
  }

  openInviteDialog() {
    this.inviteForm.reset({ externalId: '', username: '', displayName: '', role: 'MEMBER' });
    this.displayInviteDialog.set(true);
  }

  openCreateDialog() {
    const tenantId = this.tenantService.selectedTenantId();
    if (tenantId) {
      this.loadAuthSettings(tenantId);
    }
    this.createForm.reset({ username: '', email: '', displayName: '', password: '' });
    this.displayCreateDialog.set(true);
  }

  refreshUsers() {
    this.usersTable.refresh();
  }

  openUserDetail(user: User) {
    this.router.navigate(['/iam/users', user.id]);
  }

  sendInvite() {
    if (this.inviteForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { externalId, username, displayName, role } = this.inviteForm.getRawValue();
    this.isSaving.set(true);
    this.userService.inviteMember(externalId, username, displayName, role, tenantId).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gửi lời mời thành viên' });
        this.displayInviteDialog.set(false);
        this.refreshUsers();
        this.isSaving.set(false);
      },
      error: (err) => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: this.errorDetail(err, 'Không thể gửi lời mời') });
        this.isSaving.set(false);
      }
    });
  }

  createUser() {
    if (this.createForm.invalid) return;

    const tenantId = this.tenantService.selectedTenantId();
    if (!tenantId) return;

    const { username, email, displayName, password } = this.createForm.getRawValue();
    this.isSaving.set(true);
    this.userService.createUser(
      { username, email, display_name: displayName, password },
      tenantId,
    ).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã tạo người dùng mới' });
        this.displayCreateDialog.set(false);
        this.refreshUsers();
        this.isSaving.set(false);
      },
      error: (err) => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: this.errorDetail(err, 'Không thể tạo người dùng') });
        this.isSaving.set(false);
      }
    });
  }

  deleteUser(user: User) {
    this.confirmationService.confirm({
      message: `Xóa "${user.displayName || user.username || user.email}" khỏi workspace hiện tại?`,
      header: 'Xác nhận xóa',
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        const tenantId = this.tenantService.selectedTenantId();
        if (!tenantId || !user.id) return;

        this.deletingUserId.set(user.id);
        this.userService.removeMember(user.id, tenantId).subscribe({
          next: () => {
            this.messageService.add({
              severity: 'success',
              summary: 'Thành công',
              detail: 'Đã xóa người dùng khỏi workspace',
            });
            this.refreshUsers();
            this.deletingUserId.set('');
          },
          error: () => {
            this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa người dùng' });
            this.deletingUserId.set('');
          },
        });
      },
    });
  }

  getInitials(name: string): string {
    if (!name) return 'U';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().substring(0, 2);
  }

  private errorDetail(err: unknown, fallback: string): string {
    if (err && typeof err === 'object' && 'error' in err) {
      const body = (err as { error?: { message?: string; reason?: string } }).error;
      return body?.message || body?.reason || fallback;
    }
    return fallback;
  }

  private loadAuthSettings(tenantId: string) {
    this.authSettingsLoading.set(true);
    this.authSettingsService.getAuthSettings(tenantId).subscribe({
      next: settings => {
        this.loadedAuthSettingsTenantId = tenantId;
        this.passwordPolicy.set(settings.passwordPolicy);
        this.applyPasswordValidators(settings.passwordPolicy);
        this.authSettingsLoading.set(false);
      },
      error: () => {
        this.authSettingsLoading.set(false);
      },
    });
  }

  private applyPasswordValidators(policy: PasswordPolicy) {
    this.createForm.controls.password.setValidators([
      Validators.required,
      Validators.minLength(policy.minLength),
      this.passwordPolicyValidator(policy),
    ]);
    this.createForm.controls.password.updateValueAndValidity({ emitEvent: false });
  }

  private passwordPolicyValidator(policy: PasswordPolicy) {
    return (control: AbstractControl<string>) => {
      const value = control.value ?? '';
      if (!value) return null;
      const errors: Record<string, boolean> = {};
      if (policy.requireUppercase && !/[A-Z]/.test(value)) errors['uppercase'] = true;
      if (policy.requireLowercase && !/[a-z]/.test(value)) errors['lowercase'] = true;
      if (policy.requireNumber && !/\d/.test(value)) errors['number'] = true;
      if (policy.requireSymbol && !/[^A-Za-z0-9]/.test(value)) errors['symbol'] = true;
      return Object.keys(errors).length ? errors : null;
    };
  }
}
