import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable, of } from 'rxjs';
import { ConfirmationService, MessageService } from 'primeng/api';
import { Button } from 'primeng/button';
import { ConfirmDialog } from 'primeng/confirmdialog';
import { Dialog } from 'primeng/dialog';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { Textarea } from 'primeng/textarea';
import { Toast } from 'primeng/toast';
import { Tooltip } from 'primeng/tooltip';
import { AdministrativeUnit, Area, AreaAdministrativeUnit } from '../../models/mdm.models';
import { AdministrativeUnitService } from '../../services/administrative-unit.service';
import { AreaService } from '../../services/area.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-areas-page',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './areas.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AreasPage {
  private areaService = inject(AreaService);
  private unitService = inject(AdministrativeUnitService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selectedArea = signal<Area | null>(null);
  readonly selectedAreaForUnits = signal<Area | null>(null);
  readonly areaDialogVisible = signal(false);
  readonly areaUnitsVisible = signal(false);
  readonly isSaving = signal(false);

  readonly areaTypesResource = rxResource({
    stream: () => this.areaService.listTypes({ pageSize: 200 }),
  });
  readonly areasResource = rxResource({
    stream: () => this.areaService.listAreas({ pageSize: 500 }),
  });
  readonly unitsResource = rxResource({
    stream: () => this.unitService.list({ pageSize: 500 }),
  });
  readonly areaUnitsResource = rxResource({
    params: () => this.selectedAreaForUnits()?.id ?? '',
    stream: ({ params }) => params ? this.areaService.listAdministrativeUnits(params) : of([]),
  });

  readonly areaTypes = computed(() => this.areaTypesResource.value() ?? []);
  readonly areas = computed(() => this.areasResource.value() ?? []);
  readonly units = computed(() => this.unitsResource.value()?.items ?? []);
  readonly areaUnits = computed(() => this.areaUnitsResource.value() ?? []);

  readonly areaTypeOptions = computed(() => this.areaTypes().map(item => ({ label: `${item.name} (${item.code})`, value: item.id })));
  readonly areaParentOptions = computed(() => [
    { label: 'Không có parent', value: '' },
    ...this.areas()
      .filter(area => area.id !== this.selectedArea()?.id)
      .map(area => ({ label: `${area.name} (${area.code})`, value: area.id })),
  ]);
  readonly administrativeUnitOptions = computed(() => this.units().map(unit => ({ label: `${unit.fullName || unit.name} (${unit.code})`, value: unit.id })));

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
  ];
  readonly scopeOptions = [
    { label: 'Include', value: 'INCLUDE' },
    { label: 'Exclude', value: 'EXCLUDE' },
  ];

  readonly areaForm = new FormGroup({
    areaTypeId: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    parentId: new FormControl('', { nonNullable: true }),
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
    managerUserId: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    metadataJson: new FormControl('{}', { nonNullable: true }),
  });

  readonly areaUnitForm = new FormGroup({
    administrativeUnitId: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    scopeType: new FormControl('INCLUDE', { nonNullable: true }),
  });

  open(area?: Area): void {
    this.selectedArea.set(area ?? null);
    this.areaForm.reset(area ? { ...area } : {
      areaTypeId: this.areaTypes()[0]?.id ?? '',
      parentId: '',
      code: '',
      name: '',
      description: '',
      managerUserId: '',
      status: 'ACTIVE',
      effectiveFrom: '',
      effectiveTo: '',
      metadataJson: '{}',
    });
    this.areaDialogVisible.set(true);
  }

  save(): void {
    if (this.areaForm.invalid) {
      this.areaForm.markAllAsTouched();
      return;
    }
    const selected = this.selectedArea();
    this.runSave(
      selected ? this.areaService.updateArea(selected.id, this.areaForm.getRawValue()) : this.areaService.createArea(this.areaForm.getRawValue()),
      selected ? 'Đã cập nhật khu vực' : 'Đã tạo khu vực',
    );
  }

  delete(area: Area): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa khu vực "${area.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.areaService.deleteArea(area.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa khu vực' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  openAreaUnits(area: Area): void {
    this.selectedAreaForUnits.set(area);
    this.areaUnitForm.reset({ administrativeUnitId: '', scopeType: 'INCLUDE' });
    this.areaUnitsVisible.set(true);
    this.areaUnitsResource.reload();
  }

  assignAreaUnit(): void {
    if (this.areaUnitForm.invalid) {
      this.areaUnitForm.markAllAsTouched();
      return;
    }
    const area = this.selectedAreaForUnits();
    if (!area) return;
    const raw = this.areaUnitForm.getRawValue();
    this.isSaving.set(true);
    this.areaService.assignAdministrativeUnit(area.id, raw.administrativeUnitId, raw.scopeType).subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gán địa bàn' });
        this.areaUnitForm.reset({ administrativeUnitId: '', scopeType: 'INCLUDE' });
        this.areaUnitsResource.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gán địa bàn' });
        this.isSaving.set(false);
      },
    });
  }

  removeAreaUnit(item: AreaAdministrativeUnit): void {
    const area = this.selectedAreaForUnits();
    if (!area) return;
    this.confirmationService.confirm({
      header: 'Xác nhận gỡ',
      message: `Gỡ "${this.unitName(item.administrativeUnitId)}" khỏi khu vực?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Gỡ',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.areaService.removeAdministrativeUnit(area.id, item.administrativeUnitId).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã gỡ địa bàn' });
            this.areaUnitsResource.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể gỡ địa bàn' }),
        });
      },
    });
  }

  reload(): void {
    this.areasResource.reload();
    this.areaTypesResource.reload();
    this.unitsResource.reload();
    this.areaUnitsResource.reload();
  }

  areaTypeName(id: string): string {
    return this.areaTypes().find(item => item.id === id)?.name ?? '';
  }

  unitName(id: string): string {
    const unit: AdministrativeUnit | undefined = this.units().find(item => item.id === id);
    return unit ? `${unit.fullName || unit.name} (${unit.code})` : id;
  }

  statusSeverity = statusSeverity;

  private runSave(request: Observable<Area>, success: string): void {
    this.isSaving.set(true);
    request.subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: success });
        this.areaDialogVisible.set(false);
        this.reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' });
        this.isSaving.set(false);
      },
    });
  }
}
