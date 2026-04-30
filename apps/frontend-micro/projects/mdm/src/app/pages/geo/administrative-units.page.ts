import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { ConfirmationService, MessageService, TreeNode } from 'primeng/api';
import { Button } from 'primeng/button';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { DatePicker } from 'primeng/datepicker';
import { Dialog } from 'primeng/dialog';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { Textarea } from 'primeng/textarea';
import { Toast } from 'primeng/toast';
import { Tooltip } from 'primeng/tooltip';
import { TreeTableModule } from 'primeng/treetable';
import { AdministrativeUnit, AdministrativeUnitNode } from '../../models/mdm.models';
import { AdministrativeUnitService } from '../../services/administrative-unit.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-administrative-units-page',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    Button,
    ConfirmDialog,
    DatePicker,
    Dialog,
    InputText,
    Select,
    TableModule,
    Tag,
    Textarea,
    Toast,
    Tooltip,
    TreeTableModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './administrative-units.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdministrativeUnitsPage {
  private unitService = inject(AdministrativeUnitService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selectedUnit = signal<AdministrativeUnit | null>(null);
  readonly dialogVisible = signal(false);
  readonly isViewing = signal(false);
  readonly isSaving = signal(false);
  readonly keyword = signal('');
  readonly selectedProvinceId = signal('');
  readonly selectedLevel = signal('');
  readonly selectedStatus = signal('');
  readonly isSyncing = signal(false);

  readonly treeResource = rxResource({
    stream: () => this.unitService.listTree({ pageSize: 100 }),
  });

  readonly provinces = computed(() => (this.treeResource.value() ?? []).map(node => node.unit));
  readonly allUnits = computed(() => this.flattenNodes(this.treeResource.value() ?? []));
  readonly wards = computed(() => this.allUnits().filter(unit => unit.level === 'WARD'));
  readonly availableUnits = computed(() => this.allUnits());
  readonly unitNameMap = computed(() => {
    const map = new Map<string, string>();
    for (const unit of this.availableUnits()) {
      map.set(unit.id, `${unit.fullName || unit.name} (${unit.code})`);
    }
    return map;
  });
  readonly unitIndexMap = computed(() => {
    const map = new Map<string, string>();
    (this.treeResource.value() ?? []).forEach((node, provinceIndex) => {
      const parentIndex = String(provinceIndex + 1);
      map.set(node.unit.id, parentIndex);
      node.children.forEach((child, childIndex) => {
        map.set(child.unit.id, `${parentIndex}.${childIndex + 1}`);
      });
    });
    return map;
  });
  readonly activeCount = computed(() => this.allUnits().filter(unit => unit.status === 'ACTIVE').length);
  readonly treeNodes = computed<TreeNode<AdministrativeUnit>[]>(() => this.buildTreeNodes());
  readonly visibleUnitsCount = computed(() => this.treeNodes().reduce((total, node) => total + 1 + (node.children?.length ?? 0), 0));

  readonly parentOptions = computed(() => [
    { label: 'Không có parent', value: '' },
    ...this.availableUnits()
      .filter(unit => unit.id !== this.selectedUnit()?.id)
      .map(unit => ({ label: `${unit.fullName || unit.name} (${unit.code})`, value: unit.id })),
  ]);
  readonly provinceOptions = computed(() => [
    { label: 'Không có parent', value: '' },
    ...this.provinces().map(unit => ({ label: `${unit.name} (${unit.code})`, value: unit.id })),
  ]);

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
    { label: 'Merged', value: 'MERGED' },
  ];
  readonly filterStatusOptions = [
    { label: 'Tất cả trạng thái', value: '' },
    ...this.statusOptions,
  ];
  readonly filterLevelOptions = [
    { label: 'Tất cả cấp', value: '' },
    { label: 'Chỉ tỉnh/thành', value: 'PROVINCE' },
    { label: 'Có phường/xã', value: 'WARD' },
  ];
  readonly unitLevelOptions = [
    { label: 'Tỉnh/Thành phố', value: 'PROVINCE' },
    { label: 'Phường/Xã/Đặc khu', value: 'WARD' },
    { label: 'Quận/Huyện legacy', value: 'DISTRICT_LEGACY' },
  ];
  readonly unitTypeOptions = [
    { label: 'Tỉnh', value: 'TINH' },
    { label: 'Thành phố', value: 'THANH_PHO' },
    { label: 'Phường', value: 'PHUONG' },
    { label: 'Xã', value: 'XA' },
    { label: 'Đặc khu', value: 'DAC_KHU' },
    { label: 'Huyện legacy', value: 'HUYEN_LEGACY' },
  ];

  readonly form = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    fullName: new FormControl('', { nonNullable: true }),
    shortName: new FormControl('', { nonNullable: true }),
    level: new FormControl('WARD', { nonNullable: true, validators: [Validators.required] }),
    unitType: new FormControl('XA', { nonNullable: true, validators: [Validators.required] }),
    parentId: new FormControl('', { nonNullable: true }),
    path: new FormControl('', { nonNullable: true }),
    sortOrder: new FormControl(0, { nonNullable: true }),
    latitude: new FormControl(0, { nonNullable: true }),
    longitude: new FormControl(0, { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    source: new FormControl('', { nonNullable: true }),
    metadataJson: new FormControl('{}', { nonNullable: true }),
  });

  open(unit?: AdministrativeUnit): void {
    this.isViewing.set(false);
    this.form.enable();
    this.selectedUnit.set(unit ?? null);
    this.form.reset(unit ? { ...unit } : {
      code: '', name: '', fullName: '', shortName: '', level: 'WARD', unitType: 'XA', parentId: '', path: '',
      sortOrder: 0, latitude: 0, longitude: 0, status: 'ACTIVE', effectiveFrom: '', effectiveTo: '', source: '', metadataJson: '{}',
    });
    this.dialogVisible.set(true);
  }

  view(unit: AdministrativeUnit): void {
    this.selectedUnit.set(unit);
    this.isViewing.set(true);
    this.form.reset({ ...unit });
    this.form.disable();
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.isViewing()) {
      return;
    }
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const selected = this.selectedUnit();
    this.runSave(
      selected ? this.unitService.update(selected.id, this.form.getRawValue()) : this.unitService.create(this.form.getRawValue()),
      selected ? 'Đã cập nhật đơn vị hành chính' : 'Đã tạo đơn vị hành chính',
    );
  }

  delete(unit: AdministrativeUnit): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa đơn vị "${unit.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.unitService.delete(unit.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa đơn vị hành chính' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  syncFromAddressKit(): void {
    this.confirmationService.confirm({
      header: 'Đồng bộ dữ liệu hành chính',
      message: 'Thao tác này sẽ xóa toàn bộ tỉnh, phường/xã hiện có và nạp lại dữ liệu mới nhất từ CASSO AddressKit. Tiếp tục?',
      icon: 'pi pi-refresh',
      acceptLabel: 'Đồng bộ',
      rejectLabel: 'Hủy',
      accept: () => {
        this.isSyncing.set(true);
        this.unitService.syncFromAddressKit().subscribe({
          next: result => {
            this.messageService.add({
              severity: 'success',
              summary: 'Đã đồng bộ',
              detail: `Đã nạp ${result.provinceCount} tỉnh/thành và ${result.wardCount} phường/xã từ ${result.source}.`,
            });
            this.clearFilters();
            this.reload();
            this.isSyncing.set(false);
          },
          error: () => {
            this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể đồng bộ dữ liệu từ AddressKit' });
            this.isSyncing.set(false);
          },
        });
      },
    });
  }

  reload(): void {
    this.treeResource.reload();
  }

  clearFilters(): void {
    this.keyword.set('');
    this.selectedProvinceId.set('');
    this.selectedLevel.set('');
    this.selectedStatus.set('');
  }

  unitName(id: string): string {
    return this.unitNameMap().get(id) ?? '';
  }

  unitIndex(id: string): string {
    return this.unitIndexMap().get(id) ?? '';
  }

  statusSeverity = statusSeverity;

  unitTypeLabel(value: string): string {
    return this.unitTypeOptions.find(option => option.value === value)?.label ?? value;
  }

  levelLabel(value: string): string {
    return this.unitLevelOptions.find(option => option.value === value)?.label ?? value;
  }

  private runSave(request: Observable<AdministrativeUnit>, success: string): void {
    this.isSaving.set(true);
    request.subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: success });
        this.dialogVisible.set(false);
        this.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' });
        this.isSaving.set(false);
      },
    });
  }

  private buildTreeNodes(): TreeNode<AdministrativeUnit>[] {
    const keyword = this.normalize(this.keyword());
    const provinceId = this.selectedProvinceId();
    const level = this.selectedLevel();
    const status = this.selectedStatus();

    return (this.treeResource.value() ?? [])
      .filter(node => !provinceId || node.unit.id === provinceId)
      .map(node => {
        const province = node.unit;
        const children = node.children
          .map(child => child.unit)
          .filter(ward => level !== 'PROVINCE')
          .filter(ward => this.matchesStatus(ward, status))
          .filter(ward => this.matchesKeyword(ward, keyword))
          .sort((a, b) => a.sortOrder - b.sortOrder || a.name.localeCompare(b.name, 'vi'))
          .map(ward => ({ key: ward.id, data: ward, leaf: true }));
        return {
          key: province.id,
          data: province,
          expanded: Boolean(keyword || provinceId),
          children: level === 'PROVINCE' ? [] : children,
          leaf: level === 'PROVINCE',
        };
      })
      .filter(node => {
        if (level === 'WARD') {
          return Boolean(node.children?.length);
        }
        if (level === 'PROVINCE') {
          return this.matchesStatus(node.data!, status) && this.matchesKeyword(node.data!, keyword);
        }
        return this.matchesKeyword(node.data!, keyword) || Boolean(node.children?.length);
      })
      .sort((a, b) => (a.data?.sortOrder ?? 0) - (b.data?.sortOrder ?? 0) || (a.data?.name ?? '').localeCompare(b.data?.name ?? '', 'vi'));
  }

  private flattenNodes(nodes: AdministrativeUnitNode[]): AdministrativeUnit[] {
    return nodes.flatMap(node => [node.unit, ...this.flattenNodes(node.children)]);
  }

  private matchesStatus(unit: AdministrativeUnit, status: string): boolean {
    return !status || unit.status === status;
  }

  private matchesKeyword(unit: AdministrativeUnit, keyword: string): boolean {
    if (!keyword) {
      return true;
    }
    return [unit.code, unit.name, unit.fullName, unit.shortName, unit.unitType]
      .some(value => this.normalize(value).includes(keyword));
  }

  private normalize(value: string): string {
    return value.trim().toLocaleLowerCase('vi');
  }
}
