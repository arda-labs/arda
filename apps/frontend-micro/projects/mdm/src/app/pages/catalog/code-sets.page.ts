import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
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
import { CodeSet } from '../../models/mdm.models';
import { CodeCatalogService } from '../../services/code-catalog.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-code-sets-page',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './code-sets.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CodeSetsPage {
  private catalogService = inject(CodeCatalogService);
  private router = inject(Router);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selected = signal<CodeSet | null>(null);
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);

  readonly resource = rxResource({
    stream: () => this.catalogService.listSets({ pageSize: 500 }),
  });
  readonly codeSets = computed(() => this.resource.value() ?? []);

  readonly statusOptions = [
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Inactive', value: 'INACTIVE' },
  ];

  readonly form = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
    isSystem: new FormControl(false, { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  open(codeSet?: CodeSet): void {
    this.selected.set(codeSet ?? null);
    this.form.reset(codeSet ? { ...codeSet } : { code: '', name: '', description: '', isSystem: false, status: 'ACTIVE' });
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    const selected = this.selected();
    this.runSave(
      selected ? this.catalogService.updateSet(selected.id, this.form.getRawValue()) : this.catalogService.createSet(this.form.getRawValue()),
      selected ? 'Đã cập nhật danh mục' : 'Đã tạo danh mục',
    );
  }

  delete(codeSet: CodeSet): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa danh mục "${codeSet.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.catalogService.deleteSet(codeSet.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa danh mục' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  openItems(codeSet: CodeSet): void {
    void this.router.navigate(['/mdm/catalog/code-items'], { queryParams: { set: codeSet.code } });
  }

  reload(): void {
    this.resource.reload();
  }

  statusSeverity = statusSeverity;

  private runSave(request: Observable<CodeSet>, success: string): void {
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

