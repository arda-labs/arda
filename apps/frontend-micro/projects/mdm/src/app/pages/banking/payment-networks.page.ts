import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
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
import { BankBranch, PaymentNetwork } from '../../models/mdm.models';
import { statusSeverity } from '../../services/mdm-http';
import { PaymentNetworksService } from '../../services/payment-networks.service';

type PaymentTab = 'branches' | 'networks';

@Component({
  selector: 'app-payment-networks-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, ConfirmDialog, Dialog, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService, ConfirmationService],
  templateUrl: './payment-networks.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PaymentNetworksPage {
  private service = inject(PaymentNetworksService);
  private messageService = inject(MessageService);
  private confirmationService = inject(ConfirmationService);

  readonly activeTab = signal<PaymentTab>('branches');
  readonly keyword = signal('');
  readonly dialogVisible = signal(false);
  readonly isSaving = signal(false);
  readonly selectedBranch = signal<BankBranch | null>(null);
  readonly selectedNetwork = signal<PaymentNetwork | null>(null);
  readonly branchResource = rxResource({ stream: () => this.service.listBranches({ pageSize: 500 }) });
  readonly networkResource = rxResource({ stream: () => this.service.listNetworks({ pageSize: 500 }) });
  readonly branches = computed(() => this.branchResource.value()?.items ?? []);
  readonly networks = computed(() => this.networkResource.value()?.items ?? []);
  readonly activeBranches = computed(() => this.branches().filter(item => item.status === 'ACTIVE').length);
  readonly activeNetworks = computed(() => this.networks().filter(item => item.status === 'ACTIVE').length);
  readonly filteredBranches = computed(() => this.filterRows(this.branches(), item => [item.institutionCode, item.code, item.name, item.swiftCode, item.napasCode, item.provinceCode]));
  readonly filteredNetworks = computed(() => this.filterRows(this.networks(), item => [item.code, item.name, item.networkType, item.clearingMethod, item.operator]));

  readonly statusOptions = [{ label: 'Hiệu lực', value: 'ACTIVE' }, { label: 'Ngưng hiệu lực', value: 'INACTIVE' }];
  readonly branchTypeOptions = [{ label: 'Hội sở', value: 'HEAD_OFFICE' }, { label: 'Chi nhánh', value: 'BRANCH' }, { label: 'Phòng giao dịch', value: 'TRANSACTION_OFFICE' }];
  readonly networkTypeOptions = [{ label: 'Nội địa', value: 'DOMESTIC' }, { label: 'Quốc tế', value: 'INTERNATIONAL' }, { label: 'Nội bộ', value: 'INTERNAL' }];
  readonly availabilityOptions = [{ label: '24x7', value: '24X7' }, { label: 'Giờ làm việc', value: 'BUSINESS_HOURS' }, { label: 'Theo lịch riêng', value: 'CUSTOM_CALENDAR' }];

  readonly branchForm = new FormGroup({
    institutionCode: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    branchType: new FormControl('BRANCH', { nonNullable: true }),
    address: new FormControl('', { nonNullable: true }),
    provinceCode: new FormControl('', { nonNullable: true }),
    phone: new FormControl('', { nonNullable: true }),
    swiftCode: new FormControl('', { nonNullable: true }),
    napasCode: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });
  readonly networkForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    networkType: new FormControl('DOMESTIC', { nonNullable: true }),
    clearingMethod: new FormControl('', { nonNullable: true }),
    settlementCurrency: new FormControl('VND', { nonNullable: true }),
    operator: new FormControl('', { nonNullable: true }),
    availability: new FormControl('24X7', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  setTab(tab: PaymentTab): void { this.activeTab.set(tab); }
  reload(): void { this.branchResource.reload(); this.networkResource.reload(); }
  openBranch(item?: BankBranch): void { this.activeTab.set('branches'); this.selectedBranch.set(item ?? null); this.branchForm.reset(item ? { ...item } : { id: '', institutionCode: '', code: '', name: '', branchType: 'BRANCH', address: '', provinceCode: '', phone: '', swiftCode: '', napasCode: '', status: 'ACTIVE' } as BankBranch); this.dialogVisible.set(true); }
  openNetwork(item?: PaymentNetwork): void { this.activeTab.set('networks'); this.selectedNetwork.set(item ?? null); this.networkForm.reset(item ? { ...item } : { id: '', code: '', name: '', networkType: 'DOMESTIC', clearingMethod: '', settlementCurrency: 'VND', operator: '', availability: '24X7', description: '', status: 'ACTIVE' } as PaymentNetwork); this.dialogVisible.set(true); }
  save(): void { if (this.activeTab() === 'branches') return this.saveBranch(); return this.saveNetwork(); }
  deleteBranch(item: BankBranch): void { this.confirmDelete(item.name, () => this.service.deleteBranch(item.id), () => this.branchResource.reload()); }
  deleteNetwork(item: PaymentNetwork): void { this.confirmDelete(item.name, () => this.service.deleteNetwork(item.id), () => this.networkResource.reload()); }
  statusSeverity = statusSeverity;

  private saveBranch(): void { if (this.branchForm.invalid) return; const selected = this.selectedBranch(); this.runSave(selected ? this.service.updateBranch(selected.id, this.branchForm.getRawValue()) : this.service.createBranch(this.branchForm.getRawValue()), () => this.branchResource.reload()); }
  private saveNetwork(): void { if (this.networkForm.invalid) return; const selected = this.selectedNetwork(); this.runSave(selected ? this.service.updateNetwork(selected.id, this.networkForm.getRawValue()) : this.service.createNetwork(this.networkForm.getRawValue()), () => this.networkResource.reload()); }
  private runSave<T>(request: Observable<T>, reload: () => void): void { this.isSaving.set(true); request.subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu dữ liệu' }); this.dialogVisible.set(false); reload(); this.isSaving.set(false); }, error: () => { this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể lưu dữ liệu' }); this.isSaving.set(false); } }); }
  private confirmDelete(name: string, request: () => Observable<void>, reload: () => void): void { this.confirmationService.confirm({ header: 'Xác nhận xóa', message: `Xóa "${name}"?`, icon: 'pi pi-exclamation-triangle', acceptLabel: 'Xóa', rejectLabel: 'Hủy', acceptButtonStyleClass: 'p-button-danger', accept: () => request().subscribe({ next: () => { this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã xóa dữ liệu' }); reload(); }, error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể xóa dữ liệu' }) }) }); }
  private filterRows<T>(rows: T[], values: (item: T) => string[]): T[] { const keyword = this.keyword().trim().toLocaleLowerCase('vi'); return rows.filter(item => !keyword || values(item).some(value => value.toLocaleLowerCase('vi').includes(keyword))); }
}
