import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
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
import { AdministrativeUnit } from '../../models/mdm.models';
import { AdministrativeUnitService } from '../../services/administrative-unit.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-administrative-units-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    Button,
    ConfirmDialog,
    Dialog,
    InputText,
    Select,
    TableModule,
    Tag,
    Textarea,
    Toast,
    Tooltip,
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
  readonly isSaving = signal(false);

  readonly unitsResource = rxResource({
    stream: () => this.unitService.list({ pageSize: 500 }),
  });
  readonly provincesResource = rxResource({
    stream: () => this.unitService.listProvinces({ pageSize: 200 }),
  });

  readonly units = computed(() => this.unitsResource.value()?.items ?? []);
  readonly provinces = computed(() => this.provincesResource.value() ?? []);

  readonly parentOptions = computed(() => [
    { label: 'Không có parent', value: '' },
    ...this.units()
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
    this.selectedUnit.set(unit ?? null);
    this.form.reset(unit ? { ...unit } : {
      code: '', name: '', fullName: '', shortName: '', level: 'WARD', unitType: 'XA', parentId: '', path: '',
      sortOrder: 0, latitude: 0, longitude: 0, status: 'ACTIVE', effectiveFrom: '', effectiveTo: '', source: '', metadataJson: '{}',
    });
    this.dialogVisible.set(true);
  }

  save(): void {
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

  reload(): void {
    this.unitsResource.reload();
    this.provincesResource.reload();
  }

  unitName(id: string): string {
    const unit = this.units().find(item => item.id === id);
    return unit ? `${unit.fullName || unit.name} (${unit.code})` : '';
  }

  statusSeverity = statusSeverity;

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
}

