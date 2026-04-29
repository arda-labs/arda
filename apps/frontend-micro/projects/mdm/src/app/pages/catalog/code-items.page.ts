import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { rxResource, takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
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
import { CodeItem } from '../../models/mdm.models';
import { CodeCatalogService } from '../../services/code-catalog.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-code-items-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './code-items.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CodeItemsPage {
  private catalogService = inject(CodeCatalogService);
  private route = inject(ActivatedRoute);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selectedCodeSetCode = signal('');
  readonly selected = signal<CodeItem | null>(null);
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);

  readonly codeSetsResource = rxResource({
    stream: () => this.catalogService.listSets({ pageSize: 500 }),
  });
  readonly codeItemsResource = rxResource({
    params: () => this.selectedCodeSetCode(),
    stream: ({ params }) => params ? this.catalogService.listItems(params, { pageSize: 500 }) : of([]),
  });

  readonly codeSets = computed(() => this.codeSetsResource.value() ?? []);
  readonly codeItems = computed(() => this.codeItemsResource.value() ?? []);
  readonly codeSetOptions = computed(() => this.codeSets().map(item => ({ label: `${item.name} (${item.code})`, value: item.code })));
  readonly selectedCodeSetName = computed(() => this.codeSets().find(item => item.code === this.selectedCodeSetCode())?.name ?? '');

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
  ];

  readonly form = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    value: new FormControl('', { nonNullable: true }),
    parentId: new FormControl('', { nonNullable: true }),
    sortOrder: new FormControl(0, { nonNullable: true }),
    color: new FormControl('', { nonNullable: true }),
    icon: new FormControl('', { nonNullable: true }),
    metadataJson: new FormControl('{}', { nonNullable: true }),
    isDefault: new FormControl(false, { nonNullable: true }),
    isSystem: new FormControl(false, { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
  });

  constructor() {
    this.route.queryParamMap.pipe(takeUntilDestroyed()).subscribe(params => {
      const setCode = params.get('set');
      if (setCode) {
        this.selectedCodeSetCode.set(setCode);
      }
    });
  }

  selectCodeSet(code: string): void {
    this.selectedCodeSetCode.set(code);
    this.codeItemsResource.reload();
  }

  open(item?: CodeItem): void {
    if (!this.selectedCodeSetCode()) {
      this.messageService.add({ severity: 'warn', summary: 'Thiếu danh mục', detail: 'Chọn một danh mục trước khi thêm item' });
      return;
    }
    this.selected.set(item ?? null);
    this.form.reset(item ? { ...item } : {
      code: '', name: '', value: '', parentId: '', sortOrder: 0, color: '', icon: '', metadataJson: '{}',
      isDefault: false, isSystem: false, status: 'ACTIVE', effectiveFrom: '', effectiveTo: '',
    });
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const selected = this.selected();
    this.runSave(
      selected ? this.catalogService.updateItem(selected.id, this.form.getRawValue()) : this.catalogService.createItem(this.selectedCodeSetCode(), this.form.getRawValue()),
      selected ? 'Đã cập nhật item' : 'Đã tạo item',
    );
  }

  delete(item: CodeItem): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa item "${item.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.catalogService.deleteItem(item.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa item' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  reload(): void {
    this.codeSetsResource.reload();
    this.codeItemsResource.reload();
  }

  statusSeverity = statusSeverity;

  private runSave(request: Observable<CodeItem>, success: string): void {
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

