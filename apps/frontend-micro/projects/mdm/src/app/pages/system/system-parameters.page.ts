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
import { SystemParameter } from '../../models/mdm.models';
import { statusSeverity } from '../../services/mdm-http';
import { SystemParameterService } from '../../services/system-parameter.service';

@Component({
  selector: 'app-system-parameters-page',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './system-parameters.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SystemParametersPage {
  private parameterService = inject(SystemParameterService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selected = signal<SystemParameter | null>(null);
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);

  readonly resource = rxResource({
    stream: () => this.parameterService.list({ pageSize: 500 }),
  });
  readonly parameters = computed(() => this.resource.value() ?? []);

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
  ];
  readonly valueTypeOptions = [
    { label: 'String', value: 'STRING' },
    { label: 'Number', value: 'NUMBER' },
    { label: 'Boolean', value: 'BOOLEAN' },
    { label: 'JSON', value: 'JSON' },
    { label: 'Date', value: 'DATE' },
  ];

  readonly form = new FormGroup({
    key: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    groupCode: new FormControl('', { nonNullable: true }),
    valueType: new FormControl('STRING', { nonNullable: true }),
    valueText: new FormControl('', { nonNullable: true }),
    valueNumber: new FormControl(0, { nonNullable: true }),
    valueBoolean: new FormControl(false, { nonNullable: true }),
    valueJson: new FormControl('{}', { nonNullable: true }),
    defaultValue: new FormControl('', { nonNullable: true }),
    isSecret: new FormControl(false, { nonNullable: true }),
    isEditable: new FormControl(true, { nonNullable: true }),
    isSystem: new FormControl(false, { nonNullable: true }),
    validationRuleJson: new FormControl('{}', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
    updatedBy: new FormControl('', { nonNullable: true }),
  });

  open(parameter?: SystemParameter): void {
    this.selected.set(parameter ?? null);
    this.form.reset(parameter ? { ...parameter } : {
      key: '', name: '', groupCode: '', valueType: 'STRING', valueText: '', valueNumber: 0, valueBoolean: false,
      valueJson: '{}', defaultValue: '', isSecret: false, isEditable: true, isSystem: false,
      validationRuleJson: '{}', description: '', status: 'ACTIVE', updatedBy: '',
    });
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const selected = this.selected();
    const raw = this.form.getRawValue();
    this.runSave(
      selected ? this.parameterService.update(selected.key, raw) : this.parameterService.create(raw),
      selected ? 'Đã cập nhật tham số' : 'Đã tạo tham số',
    );
  }

  delete(parameter: SystemParameter): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa tham số "${parameter.key}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.parameterService.delete(parameter.key).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa tham số' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  reload(): void {
    this.resource.reload();
  }

  statusSeverity = statusSeverity;

  private runSave(request: Observable<SystemParameter>, success: string): void {
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

