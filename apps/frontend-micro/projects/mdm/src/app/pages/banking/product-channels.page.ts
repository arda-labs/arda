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
import { BankingProduct, ProductChannelRule, ServiceChannel } from '../../models/mdm.models';
import { statusSeverity } from '../../services/mdm-http';
import { ProductChannelsService } from '../../services/product-channels.service';

type ProductChannelTab = 'products' | 'channels' | 'rules';

@Component({
  selector: 'app-product-channels-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, ConfirmDialog, DatePicker, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './product-channels.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProductChannelsPage {
  private service = inject(ProductChannelsService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly activeTab = signal<ProductChannelTab>('products');
  readonly keyword = signal('');
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);
  readonly selectedProduct = signal<BankingProduct | null>(null);
  readonly selectedChannel = signal<ServiceChannel | null>(null);
  readonly selectedRule = signal<ProductChannelRule | null>(null);

  readonly productResource = rxResource({ stream: () => this.service.listProducts({ pageSize: 500 }) });
  readonly channelResource = rxResource({ stream: () => this.service.listChannels({ pageSize: 500 }) });
  readonly ruleResource = rxResource({ stream: () => this.service.listRules({ pageSize: 500 }) });
  readonly products = computed(() => this.productResource.value()?.items ?? []);
  readonly channels = computed(() => this.channelResource.value()?.items ?? []);
  readonly rules = computed(() => this.ruleResource.value()?.items ?? []);
  readonly enabledRules = computed(() => this.rules().filter(item => item.enabled).length);
  readonly filteredProducts = computed(() => this.filterRows(this.products(), item => [item.code, item.name, item.productType, item.category, item.customerSegment]));
  readonly filteredChannels = computed(() => this.filterRows(this.channels(), item => [item.code, item.name, item.channelType, item.availability]));
  readonly filteredRules = computed(() => this.filterRows(this.rules(), item => [item.productCode, item.channelCode, item.transactionType, item.feeScheduleCode, item.limitProfileCode]));

  readonly statusOptions = [{ label: 'Hiệu lực', value: 'ACTIVE' }, { label: 'Ngưng hiệu lực', value: 'INACTIVE' }];
  readonly productTypeOptions = [{ label: 'Account', value: 'ACCOUNT' }, { label: 'Deposit', value: 'DEPOSIT' }, { label: 'Loan', value: 'LOAN' }, { label: 'Card', value: 'CARD' }, { label: 'Payment', value: 'PAYMENT' }];
  readonly channelTypeOptions = [{ label: 'Branch', value: 'BRANCH' }, { label: 'Digital', value: 'DIGITAL' }, { label: 'Self service', value: 'SELF_SERVICE' }, { label: 'Partner', value: 'PARTNER' }];
  readonly availabilityOptions = [{ label: '24x7', value: '24X7' }, { label: 'Giờ làm việc', value: 'BUSINESS_HOURS' }, { label: 'Theo lịch riêng', value: 'CUSTOM_CALENDAR' }];

  readonly productForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    productType: new FormControl('ACCOUNT', { nonNullable: true }),
    category: new FormControl('', { nonNullable: true }),
    customerSegment: new FormControl('', { nonNullable: true }),
    currency: new FormControl('VND', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });
  readonly channelForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    channelType: new FormControl('DIGITAL', { nonNullable: true }),
    availability: new FormControl('24X7', { nonNullable: true }),
    timezone: new FormControl('Asia/Ho_Chi_Minh', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });
  readonly ruleForm = new FormGroup({
    productCode: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    channelCode: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    transactionType: new FormControl('', { nonNullable: true }),
    enabled: new FormControl(true, { nonNullable: true }),
    priority: new FormControl(100, { nonNullable: true }),
    feeScheduleCode: new FormControl('', { nonNullable: true }),
    limitProfileCode: new FormControl('', { nonNullable: true }),
    effectiveFrom: new FormControl('', { nonNullable: true }),
    effectiveTo: new FormControl('', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  setTab(tab: ProductChannelTab): void { this.activeTab.set(tab); }
  reload(): void { this.productResource.reload(); this.channelResource.reload(); this.ruleResource.reload(); }
  openProduct(item?: BankingProduct): void { this.activeTab.set('products'); this.selectedProduct.set(item ?? null); this.productForm.reset(item ? { ...item } : { id: '', code: '', name: '', productType: 'ACCOUNT', category: '', customerSegment: '', currency: 'VND', effectiveFrom: '', effectiveTo: '', description: '', status: 'ACTIVE' } as BankingProduct); this.dialogVisible.set(true); }
  openChannel(item?: ServiceChannel): void { this.activeTab.set('channels'); this.selectedChannel.set(item ?? null); this.channelForm.reset(item ? { ...item } : { id: '', code: '', name: '', channelType: 'DIGITAL', availability: '24X7', timezone: 'Asia/Ho_Chi_Minh', description: '', status: 'ACTIVE' } as ServiceChannel); this.dialogVisible.set(true); }
  openRule(item?: ProductChannelRule): void { this.activeTab.set('rules'); this.selectedRule.set(item ?? null); this.ruleForm.reset(item ? { ...item } : { id: '', productCode: '', channelCode: '', transactionType: '', enabled: true, priority: 100, feeScheduleCode: '', limitProfileCode: '', effectiveFrom: '', effectiveTo: '', description: '', status: 'ACTIVE' } as ProductChannelRule); this.dialogVisible.set(true); }
  save(): void { if (this.activeTab() === 'products') return this.saveProduct(); if (this.activeTab() === 'channels') return this.saveChannel(); return this.saveRule(); }
  deleteProduct(item: BankingProduct): void { this.confirmDelete(item.name, () => this.service.deleteProduct(item.id), () => this.productResource.reload()); }
  deleteChannel(item: ServiceChannel): void { this.confirmDelete(item.name, () => this.service.deleteChannel(item.id), () => this.channelResource.reload()); }
  deleteRule(item: ProductChannelRule): void { this.confirmDelete(`${item.productCode}/${item.channelCode}`, () => this.service.deleteRule(item.id), () => this.ruleResource.reload()); }
  statusSeverity = statusSeverity;

  private saveProduct(): void { if (this.productForm.invalid) return; const selected = this.selectedProduct(); this.runSave(selected ? this.service.updateProduct(selected.id, this.productForm.getRawValue()) : this.service.createProduct(this.productForm.getRawValue()), () => this.productResource.reload()); }
  private saveChannel(): void { if (this.channelForm.invalid) return; const selected = this.selectedChannel(); this.runSave(selected ? this.service.updateChannel(selected.id, this.channelForm.getRawValue()) : this.service.createChannel(this.channelForm.getRawValue()), () => this.channelResource.reload()); }
  private saveRule(): void { if (this.ruleForm.invalid) return; const selected = this.selectedRule(); this.runSave(selected ? this.service.updateRule(selected.id, this.ruleForm.getRawValue()) : this.service.createRule(this.ruleForm.getRawValue()), () => this.ruleResource.reload()); }
  private runSave<T>(request: Observable<T>, reload: () => void): void { this.isSaving.set(true); request.subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu dữ liệu' }); this.dialogVisible.set(false); reload(); this.isSaving.set(false); }, error: () => { this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' }); this.isSaving.set(false); } }); }
  private confirmDelete(name: string, request: () => Observable<void>, reload: () => void): void { this.confirmationService.confirm({ header: 'Xác nhận xóa', message: `Xóa "${name}"?`, icon: 'pi pi-exclamation-triangle', acceptLabel: 'Xóa', rejectLabel: 'Hủy', acceptButtonStyleClass: 'p-button-danger', accept: () => request().subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa dữ liệu' }); reload(); }, error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }) }) }); }
  private filterRows<T>(rows: T[], values: (item: T) => string[]): T[] { const keyword = this.keyword().trim().toLocaleLowerCase('vi'); return rows.filter(item => !keyword || values(item).some(value => value.toLocaleLowerCase('vi').includes(keyword))); }
}
