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
import { Currency, FxRate, FxRateSource } from '../../models/mdm.models';
import { CurrencyFxService } from '../../services/currency-fx.service';
import { statusSeverity } from '../../services/mdm-http';

type FxTab = 'currencies' | 'sources' | 'rates';

@Component({
  selector: 'app-currency-fx-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, ConfirmDialog, DatePicker, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './currency-fx.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CurrencyFxPage {
  private service = inject(CurrencyFxService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly activeTab = signal<FxTab>('currencies');
  readonly keyword = signal('');
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);
  readonly selectedCurrency = signal<Currency | null>(null);
  readonly selectedSource = signal<FxRateSource | null>(null);
  readonly selectedRate = signal<FxRate | null>(null);

  readonly currencyResource = rxResource({ stream: () => this.service.listCurrencies({ pageSize: 500 }) });
  readonly sourceResource = rxResource({ stream: () => this.service.listSources({ pageSize: 500 }) });
  readonly rateResource = rxResource({ stream: () => this.service.listRates({ pageSize: 500 }) });
  readonly currencies = computed(() => this.currencyResource.value()?.items ?? []);
  readonly sources = computed(() => this.sourceResource.value()?.items ?? []);
  readonly rates = computed(() => this.rateResource.value()?.items ?? []);
  readonly approvedRates = computed(() => this.rates().filter(rate => rate.approvalStatus === 'APPROVED').length);
  readonly filteredCurrencies = computed(() => this.filterRows(this.currencies(), item => [item.code, item.name, item.countryCode]));
  readonly filteredSources = computed(() => this.filterRows(this.sources(), item => [item.code, item.name, item.sourceType]));
  readonly filteredRates = computed(() => this.filterRows(this.rates(), item => [item.baseCurrency, item.quoteCurrency, item.sourceCode, item.rateDate]));

  readonly statusOptions = [{ label: 'Hiệu lực', value: 'ACTIVE' }, { label: 'Ngưng hiệu lực', value: 'INACTIVE' }];
  readonly sourceTypeOptions = [{ label: 'Central bank', value: 'CENTRAL_BANK' }, { label: 'Internal', value: 'INTERNAL' }, { label: 'Manual', value: 'MANUAL' }, { label: 'Partner', value: 'PARTNER' }];

  readonly currencyForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    numericCode: new FormControl('', { nonNullable: true }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    minorUnit: new FormControl(0, { nonNullable: true }),
    symbol: new FormControl('', { nonNullable: true }),
    countryCode: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });
  readonly sourceForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    sourceType: new FormControl('MANUAL', { nonNullable: true }),
    priority: new FormControl(100, { nonNullable: true }),
    timezone: new FormControl('Asia/Ho_Chi_Minh', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });
  readonly rateForm = new FormGroup({
    baseCurrency: new FormControl('USD', { nonNullable: true, validators: [Validators.required] }),
    quoteCurrency: new FormControl('VND', { nonNullable: true, validators: [Validators.required] }),
    sourceCode: new FormControl('TREASURY', { nonNullable: true, validators: [Validators.required] }),
    rateDate: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    buyRate: new FormControl(0, { nonNullable: true }),
    sellRate: new FormControl(0, { nonNullable: true }),
    midRate: new FormControl(0, { nonNullable: true }),
    approvalStatus: new FormControl('DRAFT', { nonNullable: true }),
    version: new FormControl(1, { nonNullable: true }),
    approvedBy: new FormControl('', { nonNullable: true }),
    changeNote: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  setTab(tab: FxTab): void { this.activeTab.set(tab); }
  reload(): void { this.currencyResource.reload(); this.sourceResource.reload(); this.rateResource.reload(); }

  openCurrency(item?: Currency): void { this.activeTab.set('currencies'); this.selectedCurrency.set(item ?? null); this.currencyForm.reset(item ? { ...item } : { id: '', code: '', numericCode: '', name: '', minorUnit: 0, symbol: '', countryCode: '', status: 'ACTIVE' } as Currency); this.dialogVisible.set(true); }
  openSource(item?: FxRateSource): void { this.activeTab.set('sources'); this.selectedSource.set(item ?? null); this.sourceForm.reset(item ? { ...item } : { id: '', code: '', name: '', sourceType: 'MANUAL', priority: 100, timezone: 'Asia/Ho_Chi_Minh', description: '', status: 'ACTIVE' } as FxRateSource); this.dialogVisible.set(true); }
  openRate(item?: FxRate): void { this.activeTab.set('rates'); this.selectedRate.set(item ?? null); this.rateForm.reset(item ? { ...item } : { id: '', baseCurrency: 'USD', quoteCurrency: 'VND', sourceCode: 'TREASURY', rateDate: '', buyRate: 0, sellRate: 0, midRate: 0, approvalStatus: 'DRAFT', version: 1, approvedBy: '', changeNote: '', status: 'ACTIVE' } as FxRate); this.dialogVisible.set(true); }

  save(): void {
    if (this.activeTab() === 'currencies') return this.saveCurrency();
    if (this.activeTab() === 'sources') return this.saveSource();
    return this.saveRate();
  }

  deleteCurrency(item: Currency): void { this.confirmDelete(item.name, () => this.service.deleteCurrency(item.id), () => this.currencyResource.reload()); }
  deleteSource(item: FxRateSource): void { this.confirmDelete(item.name, () => this.service.deleteSource(item.id), () => this.sourceResource.reload()); }
  deleteRate(item: FxRate): void { this.confirmDelete(`${item.baseCurrency}/${item.quoteCurrency}`, () => this.service.deleteRate(item.id), () => this.rateResource.reload()); }
  approveRate(item: FxRate): void { this.runSave(this.service.approveRate(item.id), () => this.rateResource.reload()); }

  statusSeverity = statusSeverity;

  private saveCurrency(): void { if (this.currencyForm.invalid) return; const selected = this.selectedCurrency(); this.runSave(selected ? this.service.updateCurrency(selected.id, this.currencyForm.getRawValue()) : this.service.createCurrency(this.currencyForm.getRawValue()), () => this.currencyResource.reload()); }
  private saveSource(): void { if (this.sourceForm.invalid) return; const selected = this.selectedSource(); this.runSave(selected ? this.service.updateSource(selected.id, this.sourceForm.getRawValue()) : this.service.createSource(this.sourceForm.getRawValue()), () => this.sourceResource.reload()); }
  private saveRate(): void { if (this.rateForm.invalid) return; const selected = this.selectedRate(); this.runSave(selected ? this.service.updateRate(selected.id, this.rateForm.getRawValue()) : this.service.createRate(this.rateForm.getRawValue()), () => this.rateResource.reload()); }

  private runSave<T>(request: Observable<T>, reload: () => void): void {
    this.isSaving.set(true);
    request.subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu dữ liệu' }); this.dialogVisible.set(false); reload(); this.isSaving.set(false); }, error: () => { this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' }); this.isSaving.set(false); } });
  }
  private confirmDelete(name: string, request: () => Observable<void>, reload: () => void): void {
    this.confirmationService.confirm({ header: 'Xác nhận xóa', message: `Xóa "${name}"?`, icon: 'pi pi-exclamation-triangle', acceptLabel: 'Xóa', rejectLabel: 'Hủy', acceptButtonStyleClass: 'p-button-danger', accept: () => request().subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa dữ liệu' }); reload(); }, error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }) }) });
  }
  private filterRows<T>(rows: T[], values: (item: T) => string[]): T[] { const keyword = this.keyword().trim().toLocaleLowerCase('vi'); return rows.filter(item => !keyword || values(item).some(value => value.toLocaleLowerCase('vi').includes(keyword))); }
}
