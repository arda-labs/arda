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
import { FeeSchedule, StandardLimit, TaxRule } from '../../models/mdm.models';
import { statusSeverity } from '../../services/mdm-http';
import { PricingRuleService } from '../../services/pricing-rule.service';

type PricingTab = 'fees' | 'taxes' | 'limits';

@Component({
  selector: 'app-pricing-rules-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, ConfirmDialog, DatePicker, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './pricing-rules.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PricingRulesPage {
  private service = inject(PricingRuleService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly activeTab = signal<PricingTab>('fees');
  readonly keyword = signal('');
  readonly selectedStatus = signal('');
  readonly isSaving = signal(false);
  readonly dialogVisible = signal(false);
  readonly selectedFee = signal<FeeSchedule | null>(null);
  readonly selectedTax = signal<TaxRule | null>(null);
  readonly selectedLimit = signal<StandardLimit | null>(null);

  readonly feeResource = rxResource({ stream: () => this.service.listFeeSchedules({ pageSize: 500 }) });
  readonly taxResource = rxResource({ stream: () => this.service.listTaxRules({ pageSize: 500 }) });
  readonly limitResource = rxResource({ stream: () => this.service.listStandardLimits({ pageSize: 500 }) });

  readonly fees = computed(() => this.feeResource.value()?.items ?? []);
  readonly taxes = computed(() => this.taxResource.value()?.items ?? []);
  readonly limits = computed(() => this.limitResource.value()?.items ?? []);
  readonly filteredFees = computed(() => this.filterRows(this.fees(), item => [item.code, item.name, item.feeType, item.channel, item.productCode, item.currency, item.status]));
  readonly filteredTaxes = computed(() => this.filterRows(this.taxes(), item => [item.code, item.name, item.taxType, item.jurisdiction, item.status]));
  readonly filteredLimits = computed(() => this.filterRows(this.limits(), item => [item.code, item.name, item.limitType, item.subjectType, item.channel, item.productCode, item.currency, item.status]));
  readonly activeCount = computed(() => this.fees().filter(item => item.status === 'ACTIVE').length + this.taxes().filter(item => item.status === 'ACTIVE').length + this.limits().filter(item => item.status === 'ACTIVE').length);

  readonly statusOptions = [
    { label: 'Hiệu lực', value: 'ACTIVE' },
    { label: 'Ngưng hiệu lực', value: 'INACTIVE' },
  ];
  readonly filterStatusOptions = [{ label: 'Tất cả trạng thái', value: '' }, ...this.statusOptions];
  readonly feeTypeOptions = [
    { label: 'Phí chuyển khoản', value: 'TRANSFER_FEE' },
    { label: 'Phí rút tiền', value: 'WITHDRAWAL_FEE' },
    { label: 'Phí dịch vụ', value: 'SERVICE_FEE' },
    { label: 'Phí thẻ', value: 'CARD_FEE' },
  ];
  readonly calculationMethodOptions = [
    { label: 'Cố định', value: 'FIXED' },
    { label: 'Theo tỷ lệ', value: 'PERCENTAGE' },
    { label: 'Cố định + tỷ lệ', value: 'HYBRID' },
  ];
  readonly taxTypeOptions = [
    { label: 'VAT', value: 'VAT' },
    { label: 'Withholding', value: 'WITHHOLDING' },
  ];
  readonly limitTypeOptions = [
    { label: 'Hạn mức chuyển khoản', value: 'TRANSFER_AMOUNT' },
    { label: 'Hạn mức rút tiền', value: 'WITHDRAWAL_AMOUNT' },
    { label: 'Hạn mức giao dịch', value: 'TRANSACTION_AMOUNT' },
  ];
  readonly subjectTypeOptions = [
    { label: 'Khách hàng cá nhân', value: 'RETAIL_CUSTOMER' },
    { label: 'Khách hàng doanh nghiệp', value: 'CORPORATE_CUSTOMER' },
    { label: 'Thẻ', value: 'CARD' },
    { label: 'Tài khoản', value: 'ACCOUNT' },
  ];

  readonly feeForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    feeType: new FormControl('SERVICE_FEE', { nonNullable: true }),
    calculationMethod: new FormControl('FIXED', { nonNullable: true }),
    currency: new FormControl('VND', { nonNullable: true }),
    fixedAmount: new FormControl(0, { nonNullable: true }),
    ratePercent: new FormControl(0, { nonNullable: true }),
    minAmount: new FormControl(0, { nonNullable: true }),
    maxAmount: new FormControl(0, { nonNullable: true }),
    channel: new FormControl('', { nonNullable: true }),
    productCode: new FormControl('', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  readonly taxForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    taxType: new FormControl('VAT', { nonNullable: true }),
    ratePercent: new FormControl(10, { nonNullable: true }),
    inclusive: new FormControl(false, { nonNullable: true }),
    jurisdiction: new FormControl('VN', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  readonly limitForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    limitType: new FormControl('TRANSACTION_AMOUNT', { nonNullable: true }),
    subjectType: new FormControl('CUSTOMER', { nonNullable: true }),
    currency: new FormControl('VND', { nonNullable: true }),
    minAmount: new FormControl(0, { nonNullable: true }),
    perTxnAmount: new FormControl(0, { nonNullable: true }),
    dailyAmount: new FormControl(0, { nonNullable: true }),
    monthlyAmount: new FormControl(0, { nonNullable: true }),
    countLimit: new FormControl(0, { nonNullable: true }),
    channel: new FormControl('', { nonNullable: true }),
    productCode: new FormControl('', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  setTab(tab: PricingTab): void {
    this.activeTab.set(tab);
  }

  openFee(item?: FeeSchedule): void {
    this.activeTab.set('fees');
    this.selectedFee.set(item ?? null);
    this.feeForm.reset(item ? { ...item } : this.emptyFee());
    this.dialogVisible.set(true);
  }

  openTax(item?: TaxRule): void {
    this.activeTab.set('taxes');
    this.selectedTax.set(item ?? null);
    this.taxForm.reset(item ? { ...item } : this.emptyTax());
    this.dialogVisible.set(true);
  }

  openLimit(item?: StandardLimit): void {
    this.activeTab.set('limits');
    this.selectedLimit.set(item ?? null);
    this.limitForm.reset(item ? { ...item } : this.emptyLimit());
    this.dialogVisible.set(true);
  }

  save(): void {
    if (this.activeTab() === 'fees') this.saveFee();
    if (this.activeTab() === 'taxes') this.saveTax();
    if (this.activeTab() === 'limits') this.saveLimit();
  }

  deleteFee(item: FeeSchedule): void {
    this.confirmDelete(item.name, () => this.service.deleteFeeSchedule(item.id), () => this.feeResource.reload());
  }

  deleteTax(item: TaxRule): void {
    this.confirmDelete(item.name, () => this.service.deleteTaxRule(item.id), () => this.taxResource.reload());
  }

  deleteLimit(item: StandardLimit): void {
    this.confirmDelete(item.name, () => this.service.deleteStandardLimit(item.id), () => this.limitResource.reload());
  }

  reload(): void {
    this.feeResource.reload();
    this.taxResource.reload();
    this.limitResource.reload();
  }

  clearFilters(): void {
    this.keyword.set('');
    this.selectedStatus.set('');
  }

  formatMoney(value: number, currency = 'VND'): string {
    return new Intl.NumberFormat('vi-VN', { maximumFractionDigits: 2 }).format(value) + ` ${currency}`;
  }

  statusSeverity = statusSeverity;

  private saveFee(): void {
    if (this.feeForm.invalid) return this.feeForm.markAllAsTouched();
    const selected = this.selectedFee();
    this.runSave(selected ? this.service.updateFeeSchedule(selected.id, this.feeForm.getRawValue()) : this.service.createFeeSchedule(this.feeForm.getRawValue()), () => this.feeResource.reload());
  }

  private saveTax(): void {
    if (this.taxForm.invalid) return this.taxForm.markAllAsTouched();
    const selected = this.selectedTax();
    this.runSave(selected ? this.service.updateTaxRule(selected.id, this.taxForm.getRawValue()) : this.service.createTaxRule(this.taxForm.getRawValue()), () => this.taxResource.reload());
  }

  private saveLimit(): void {
    if (this.limitForm.invalid) return this.limitForm.markAllAsTouched();
    const selected = this.selectedLimit();
    this.runSave(selected ? this.service.updateStandardLimit(selected.id, this.limitForm.getRawValue()) : this.service.createStandardLimit(this.limitForm.getRawValue()), () => this.limitResource.reload());
  }

  private runSave<T>(request: Observable<T>, reload: () => void): void {
    this.isSaving.set(true);
    request.subscribe({
      next: () => {
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu cấu hình' });
        this.dialogVisible.set(false);
        reload();
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu cấu hình' });
        this.isSaving.set(false);
      },
    });
  }

  private confirmDelete(name: string, request: () => Observable<void>, reload: () => void): void {
    this.confirmationService.confirm({
      header: 'Xác nhận xóa',
      message: `Xóa cấu hình "${name}"?`,
      icon: 'pi pi-exclamation-triangle',
      acceptLabel: 'Xóa',
      rejectLabel: 'Hủy',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => request().subscribe({
        next: () => {
          this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa cấu hình' });
          reload();
        },
        error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa cấu hình' }),
      }),
    });
  }

  private filterRows<T>(rows: T[], values: (item: T) => string[]): T[] {
    const keyword = this.keyword().trim().toLocaleLowerCase('vi');
    const status = this.selectedStatus();
    return rows
      .filter((item: any) => !status || item.status === status)
      .filter(item => !keyword || values(item).some(value => value.toLocaleLowerCase('vi').includes(keyword)));
  }

  private emptyFee(): FeeSchedule {
    return { id: '', code: '', name: '', feeType: 'SERVICE_FEE', calculationMethod: 'FIXED', currency: 'VND', fixedAmount: 0, ratePercent: 0, minAmount: 0, maxAmount: 0, channel: '', productCode: '', effectiveFrom: '', effectiveTo: '', description: '', status: 'ACTIVE' };
  }

  private emptyTax(): TaxRule {
    return { id: '', code: '', name: '', taxType: 'VAT', ratePercent: 10, inclusive: false, jurisdiction: 'VN', effectiveFrom: '', effectiveTo: '', description: '', status: 'ACTIVE' };
  }

  private emptyLimit(): StandardLimit {
    return { id: '', code: '', name: '', limitType: 'TRANSACTION_AMOUNT', subjectType: 'CUSTOMER', currency: 'VND', minAmount: 0, perTxnAmount: 0, dailyAmount: 0, monthlyAmount: 0, countLimit: 0, channel: '', productCode: '', effectiveFrom: '', effectiveTo: '', description: '', status: 'ACTIVE' };
  }
}
