import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { rxResource } from '@angular/core/rxjs-interop';
import { of } from 'rxjs';
import { TableModule } from 'primeng/table';
import { Button } from 'primeng/button';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { Tag } from 'primeng/tag';
import { Dialog } from 'primeng/dialog';
import { Toast } from 'primeng/toast';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { Tooltip } from 'primeng/tooltip';
import { ConfirmationService, MessageService } from 'primeng/api';
import { MenuAdmin, MenuAdminItem, MenuAdminPayload } from '../../../../services/menu-admin';
import { UserService } from '../../../../services/user.service';
import { TenantService } from '../../../../services/tenant.service';

interface MenuRow {
  menu: MenuAdminItem;
  level: number;
  parentName: string;
}

@Component({
  selector: 'app-menu-management',
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TableModule,
    Button,
    InputText,
    Select,
    Tag,
    Dialog,
    Toast,
    ConfirmDialog,
    Tooltip,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './menu-management.html',
  styleUrl: './menu-management.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MenuManagement {
  private menuAdmin = inject(MenuAdmin);
  private userService = inject(UserService);
  private tenantService = inject(TenantService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly menusResource = rxResource({
    params: () => ({} as Record<string, never>),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.menuAdmin.listMenus();
    },
  });

  readonly permissionsResource = rxResource({
    params: () => this.tenantService.selectedTenantId(),
    stream: (params) => {
      if (!params.params) return of([]);
      return this.userService.listPermissions(params.params);
    },
  });

  readonly menuRows = computed(() => this.flattenMenus(this.menusResource.value() ?? []));
  readonly totalMenus = computed(() => (this.menusResource.value() ?? []).length);
  readonly enabledMenus = computed(() => (this.menusResource.value() ?? []).filter(menu => menu.enabled).length);
  readonly permissionBoundMenus = computed(() => (this.menusResource.value() ?? []).filter(menu => menu.permissionSlug).length);
  readonly rootMenus = computed(() => (this.menusResource.value() ?? []).filter(menu => !menu.parentId).length);
  readonly parentOptions = computed(() => this.buildParentOptions(this.menusResource.value() ?? []));
  readonly permissionOptions = computed(() => [
    { label: 'Không yêu cầu quyền', value: '' },
    ...(this.permissionsResource.value() ?? []).map(permission => ({ label: permission, value: permission })),
  ]);

  readonly displayDialog = signal(false);
  readonly selectedMenu = signal<MenuAdminItem | null>(null);
  readonly isSaving = signal(false);

  readonly menuForm = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(2)] }),
    slug: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(2)] }),
    icon: new FormControl('pi pi-circle', { nonNullable: true }),
    route: new FormControl('', { nonNullable: true }),
    parentId: new FormControl('', { nonNullable: true }),
    sortOrder: new FormControl(0, { nonNullable: true }),
    enabled: new FormControl(true, { nonNullable: true }),
    permissionSlug: new FormControl('', { nonNullable: true }),
  });

  openNew(): void {
    this.selectedMenu.set(null);
    this.menuForm.reset({
      name: '',
      slug: '',
      icon: 'pi pi-circle',
      route: '',
      parentId: '',
      sortOrder: this.nextSortOrder(),
      enabled: true,
      permissionSlug: '',
    });
    this.displayDialog.set(true);
  }

  editMenu(menu: MenuAdminItem): void {
    this.selectedMenu.set(menu);
    this.menuForm.reset({
      name: menu.name,
      slug: menu.slug,
      icon: menu.icon || 'pi pi-circle',
      route: menu.route,
      parentId: menu.parentId,
      sortOrder: menu.sortOrder,
      enabled: menu.enabled,
      permissionSlug: menu.permissionSlug,
    });
    this.displayDialog.set(true);
  }

  saveMenu(): void {
    if (this.menuForm.invalid) {
      this.menuForm.markAllAsTouched();
      return;
    }

    const payload = this.formPayload();
    const selected = this.selectedMenu();
    const request$ = selected
      ? this.menuAdmin.updateMenu(selected.id, payload)
      : this.menuAdmin.createMenu(payload);

    this.isSaving.set(true);
    request$.subscribe({
      next: () => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: selected ? 'Đã cập nhật menu' : 'Đã tạo menu mới',
        });
        this.displayDialog.set(false);
        this.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu menu' });
        this.isSaving.set(false);
      },
    });
  }

  toggleEnabled(menu: MenuAdminItem): void {
    this.menuAdmin.updateMenu(menu.id, { ...this.toPayload(menu), enabled: !menu.enabled }).subscribe({
      next: () => {
        this.messageService.add({
          severity: 'success',
          summary: 'Thành công',
          detail: menu.enabled ? 'Đã tắt menu' : 'Đã bật menu',
        });
        this.reload();
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể đổi trạng thái menu' });
      },
    });
  }

  deleteMenu(menu: MenuAdminItem): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa menu "${menu.name}"? Các menu con cũng sẽ bị xóa.`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.menuAdmin.deleteMenu(menu.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa menu' });
            this.reload();
          },
          error: () => {
            this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa menu' });
          },
        });
      },
    });
  }

  reload(): void {
    this.menusResource.reload();
    this.permissionsResource.reload();
  }

  private formPayload(): MenuAdminPayload {
    const raw = this.menuForm.getRawValue();
    return {
      parentId: raw.parentId,
      name: raw.name.trim(),
      slug: raw.slug.trim(),
      icon: raw.icon.trim(),
      route: raw.route.trim(),
      sortOrder: Number(raw.sortOrder || 0),
      enabled: raw.enabled,
      permissionSlug: raw.permissionSlug,
    };
  }

  private toPayload(menu: MenuAdminItem): MenuAdminPayload {
    return {
      parentId: menu.parentId,
      name: menu.name,
      slug: menu.slug,
      icon: menu.icon,
      route: menu.route,
      sortOrder: menu.sortOrder,
      enabled: menu.enabled,
      permissionSlug: menu.permissionSlug,
    };
  }

  private nextSortOrder(): number {
    const menus = this.menusResource.value() ?? [];
    if (menus.length === 0) return 10;
    return Math.max(...menus.map(menu => menu.sortOrder)) + 1;
  }

  private buildParentOptions(menus: MenuAdminItem[]): { label: string; value: string }[] {
    const selectedId = this.selectedMenu()?.id ?? '';
    const blocked = new Set<string>();
    if (selectedId) {
      this.collectDescendantIds(selectedId, menus, blocked);
    }

    return [
      { label: 'Không có parent', value: '' },
      ...this.flattenMenus(menus)
        .filter(row => row.menu.id !== selectedId && !blocked.has(row.menu.id))
        .map(row => ({
          label: `${'  '.repeat(row.level)}${row.menu.name}`,
          value: row.menu.id,
        })),
    ];
  }

  private flattenMenus(menus: MenuAdminItem[]): MenuRow[] {
    const byId = new Map(menus.map(menu => [menu.id, menu]));
    const children = new Map<string, MenuAdminItem[]>();
    const roots: MenuAdminItem[] = [];

    for (const menu of menus) {
      if (menu.parentId && byId.has(menu.parentId)) {
        const list = children.get(menu.parentId) ?? [];
        list.push(menu);
        children.set(menu.parentId, list);
      } else {
        roots.push(menu);
      }
    }

    const sort = (items: MenuAdminItem[]) => items.sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name));
    const rows: MenuRow[] = [];
    const visit = (menu: MenuAdminItem, level: number) => {
      rows.push({
        menu,
        level,
        parentName: menu.parentId ? byId.get(menu.parentId)?.name ?? 'Không rõ' : '',
      });
      for (const child of sort(children.get(menu.id) ?? [])) {
        visit(child, level + 1);
      }
    };

    for (const root of sort(roots)) {
      visit(root, 0);
    }
    return rows;
  }

  private collectDescendantIds(menuId: string, menus: MenuAdminItem[], out: Set<string>): void {
    for (const menu of menus) {
      if (menu.parentId === menuId) {
        out.add(menu.id);
        this.collectDescendantIds(menu.id, menus, out);
      }
    }
  }
}
