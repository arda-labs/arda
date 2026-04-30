import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { ConfirmationService, MessageService } from 'primeng/api';
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
import { CreditInstitution } from '../../models/mdm.models';
import { CreditInstitutionService } from '../../services/credit-institution.service';
import { statusSeverity } from '../../services/mdm-http';

@Component({
  selector: 'app-credit-institutions-page',
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
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './credit-institutions.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CreditInstitutionsPage {
  private creditInstitutionService = inject(CreditInstitutionService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly selected = signal<CreditInstitution | null>(null);
  readonly dialogVisible = signal(false);
  readonly isViewing = signal(false);
  readonly isSaving = signal(false);
  readonly keyword = signal('');
  readonly selectedStatus = signal('');

  readonly resource = rxResource({
    stream: () => this.creditInstitutionService.list({ pageSize: 500 }),
  });
  readonly institutions = computed(() => this.resource.value()?.items ?? []);
  readonly filteredInstitutions = computed(() => {
    const keyword = this.normalize(this.keyword());
    const status = this.selectedStatus();
    return this.institutions()
      .filter(item => !status || item.status === status)
      .filter(item => !keyword || [
        item.code, item.name, item.shortName, item.address, item.phone,
        item.email, item.licenseNumber, item.taxCode, item.website,
      ].some(value => this.normalize(value).includes(keyword)));
  });

  readonly activeCount = computed(() => this.institutions().filter(item => item.status === 'ACTIVE').length);
  readonly inactiveCount = computed(() => this.institutions().filter(item => item.status === 'INACTIVE').length);

  readonly statusOptions = [
    { label: 'Hiệu lực', value: 'ACTIVE' },
    { label: 'Ngưng hiệu lực', value: 'INACTIVE' },
  ];
  readonly filterStatusOptions = [
    { label: 'Tất cả trạng thái', value: '' },
    ...this.statusOptions,
  ];

  readonly form = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    shortName: new FormControl('', { nonNullable: true }),
    address: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    phone: new FormControl('', { nonNullable: true }),
    email: new FormControl('', { nonNullable: true, validators: [Validators.email] }),
    licenseNumber: new FormControl('', { nonNullable: true }),
    issuedDate: new FormControl('', { nonNullable: true }),
    taxCode: new FormControl('', { nonNullable: true }),
    website: new FormControl('', { nonNullable: true }),
    note: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  open(item?: CreditInstitution): void {
    this.isViewing.set(false);
    this.form.enable();
    this.selected.set(item ?? null);
    this.form.reset(item ? { ...item } : this.emptyFormValue());
    this.dialogVisible.set(true);
  }

  view(item: CreditInstitution): void {
    this.selected.set(item);
    this.isViewing.set(true);
    this.form.reset({ ...item });
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
    const selected = this.selected();
    this.runSave(
      selected ? this.creditInstitutionService.update(selected.id, this.form.getRawValue()) : this.creditInstitutionService.create(this.form.getRawValue()),
      selected ? 'Đã cập nhật tổ chức tín dụng' : 'Đã tạo tổ chức tín dụng',
    );
  }

  delete(item: CreditInstitution): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa tổ chức "${item.name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.creditInstitutionService.delete(item.id).subscribe({
          next: () => {
            this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa tổ chức tín dụng' });
            this.reload();
          },
          error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }),
        });
      },
    });
  }

  clearFilters(): void {
    this.keyword.set('');
    this.selectedStatus.set('');
  }

  reload(): void {
    this.resource.reload();
  }

  statusSeverity = statusSeverity;

  private runSave(request: Observable<CreditInstitution>, success: string): void {
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

  private emptyFormValue(): CreditInstitution {
    return {
      id: '',
      code: '',
      name: '',
      shortName: '',
      address: '',
      phone: '',
      email: '',
      licenseNumber: '',
      issuedDate: '',
      taxCode: '',
      website: '',
      note: '',
      status: 'ACTIVE',
    };
  }

  private normalize(value: string): string {
    return value.trim().toLocaleLowerCase('vi');
  }
}
