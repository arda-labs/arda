import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { rxResource } from '@angular/core/rxjs-interop';
import { forkJoin, of } from 'rxjs';
import { Avatar } from 'primeng/avatar';
import { Button } from 'primeng/button';
import { Select } from 'primeng/select';
import { Tag } from 'primeng/tag';
import { Toast } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { Tabs, TabList, Tab, TabPanels, TabPanel } from 'primeng/tabs';
import { Group, Role, User, UserService, UserTenantAccess } from '../../../services/user.service';
import { TenantService } from '../../../services/tenant.service';

interface UserDetailData {
  user: User;
  roles: Role[];
  allRoles: Role[];
  tenantAccess: UserTenantAccess[];
  groups: Group[];
  effectivePermissions: string[];
}

@Component({
  selector: 'app-user-detail',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    Avatar,
    Button,
    Select,
    Tag,
    Toast,
    Tabs,
    TabList,
    Tab,
    TabPanels,
    TabPanel,
  ],
  providers: [MessageService],
  templateUrl: './user-detail.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserDetail {
  private route = inject(ActivatedRoute);
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);

  readonly userId = signal(this.route.snapshot.paramMap.get('id') ?? '');
  readonly selectedRoleControl = new FormControl('', { nonNullable: true, validators: [Validators.required] });
  readonly isSaving = signal(false);

  readonly detailResource = rxResource({
    params: () => ({
      userId: this.userId(),
      tenantId: this.tenantService.selectedTenantId(),
    }),
    stream: ({ params }) => {
      if (!params.userId || !params.tenantId) {
        return of(null);
      }

      return forkJoin({
        user: this.userService.getTenantUser(params.userId, params.tenantId),
        roles: this.userService.listUserRoles(params.userId, params.tenantId),
        allRoles: this.userService.listRoles(params.tenantId, { pageSize: 100 }),
        tenantAccess: this.userService.listUserTenantAccess(params.userId),
        groups: this.userService.listUserGroups(params.userId, params.tenantId),
        effectivePermissions: this.userService.getUserEffectivePermissions(params.userId, params.tenantId),
      });
    },
  });

  readonly detail = computed(() => this.detailResource.value() as UserDetailData | null | undefined);
  readonly user = computed(() => this.detail()?.user ?? null);
  readonly assignedRoles = computed(() => this.detail()?.roles ?? []);
  readonly tenantAccess = computed(() => this.detail()?.tenantAccess ?? []);
  readonly userGroups = computed(() => this.detail()?.groups ?? []);
  readonly effectivePermissions = computed(() => this.detail()?.effectivePermissions ?? []);
  readonly totalTenantPermissions = computed(() =>
    this.tenantAccess().reduce((total, tenant) => total + tenant.permissions.length, 0)
  );
  readonly availableRoles = computed(() => {
    const assigned = new Set(this.assignedRoles().map(role => role.id));
    return (this.detail()?.allRoles ?? []).filter(role => !assigned.has(role.id));
  });

  assignSelectedRole(): void {
    if (this.selectedRoleControl.invalid) return;

    const user = this.user();
    const tenantId = this.tenantService.selectedTenantId();
    const roleId = this.selectedRoleControl.value;
    if (!user || !tenantId || !roleId) return;

    this.isSaving.set(true);
    this.userService.assignRole(user.id, roleId, tenantId).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gán vai trò' });
        this.selectedRoleControl.reset('');
        this.detailResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gán vai trò' });
        this.isSaving.set(false);
      },
    });
  }

  revokeRole(role: Role): void {
    const user = this.user();
    const tenantId = this.tenantService.selectedTenantId();
    if (!user || !tenantId) return;

    this.isSaving.set(true);
    this.userService.revokeRole(user.id, role.id, tenantId).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gỡ vai trò' });
        this.detailResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gỡ vai trò' });
        this.isSaving.set(false);
      },
    });
  }

  getInitials(name: string | undefined): string {
    if (!name) return 'U';
    return name.split(' ').map(part => part[0]).join('').toUpperCase().substring(0, 2);
  }
}
