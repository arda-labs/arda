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
import { AreaType } from '../../models/mdm.models';
import { AreaService } from '../../services/area.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-area-types-page',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './area-types.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AreaTypesPage {
  private areaService = inject(AreaService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selected = signal<AreaType | null>(null);
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);

  readonly resource = rxResource({
    stream: () => this.areaService.listTypes({ pageSize: 500 }),
  });
  readonly areaTypes = computed(() => this.resource.value() ?? []);

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
  ];

  readonly form = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
    allowHierarchy: new FormControl(true, { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  open(areaType?: AreaType): void {
    this.selected.set(areaType ?? null);
    this.form.reset(areaType ? { ...areaType } : { code: '', name: '', description: '', allowHierarchy: true, status: 'ACTIVE' });
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const selected = this.selected();
    this.runSave(
      selected ? this.areaService.updateType(selected.id, this.form.getRawValue()) : this.areaService.createType(this.form.getRawValue()),
      selected ? 'Đã cập nhật loại khu vực' : 'Đã tạo loại khu vực',
    );
  }

  delete(areaType: AreaType): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa loại khu vực "${areaType.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.areaService.deleteType(areaType.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa loại khu vực' });
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

  private runSave(request: Observable<AreaType>, success: string): void {
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

