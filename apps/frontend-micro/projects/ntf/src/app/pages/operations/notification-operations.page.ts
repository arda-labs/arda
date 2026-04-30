import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { rxResource } from '@angular/core/rxjs-interop';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Observable } from 'rxjs';
import { MessageService } from 'primeng/api';
import { Button } from 'primeng/button';
import { Dialog } from 'primeng/dialog';
import { InputNumber } from 'primeng/inputnumber';
import { InputText } from 'primeng/inputtext';
import { Select } from 'primeng/select';
import { TableModule } from 'primeng/table';
import { Tag } from 'primeng/tag';
import { Textarea } from 'primeng/textarea';
import { Toast } from 'primeng/toast';
import { Tooltip } from 'primeng/tooltip';
import { NotificationDelivery, NotificationTemplate, NotificationTemplateVersion, ProviderConfig } from '../../models/notification.models';
import { statusSeverity } from '../../services/notification-http';
import { NotificationService } from '../../services/notification.service';

type NotificationTab = 'templates' | 'queue' | 'providers';

@Component({
  selector: 'app-notification-operations-page',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, Button, Dialog, InputNumber, InputText, Select, TableModule, Tag, Textarea, Toast, Tooltip],
  providers: [MessageService],
  templateUrl: './notification-operations.page.html',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NotificationOperationsPage {
  private service = inject(NotificationService);
  private messageService = inject(MessageService);

  readonly activeTab = signal<NotificationTab>('templates');
  readonly keyword = signal('');
  readonly dialogVisible = signal(false);
  readonly versionDialogVisible = signal(false);
  readonly providerDialogVisible = signal(false);
  readonly selectedTemplate = signal<NotificationTemplate | null>(null);
  readonly selectedProvider = signal<ProviderConfig | null>(null);
  readonly selectedTemplateForVersions = signal<NotificationTemplate | null>(null);
  readonly templateVersions = signal<NotificationTemplateVersion[]>([]);
  readonly isSaving = signal(false);

  readonly templateResource = rxResource({ stream: () => this.service.listTemplates({ pageSize: 500 }) });
  readonly deliveryResource = rxResource({ stream: () => this.service.listDeliveries({ pageSize: 500 }) });
  readonly providerResource = rxResource({ stream: () => this.service.listProviderConfigs({ pageSize: 500 }) });

  readonly templates = computed(() => this.templateResource.value()?.items ?? []);
  readonly deliveries = computed(() => this.deliveryResource.value()?.items ?? []);
  readonly providers = computed(() => this.providerResource.value() ?? []);
  readonly queuedCount = computed(() => this.deliveries().filter(item => item.status === 'QUEUED').length);
  readonly retryingCount = computed(() => this.deliveries().filter(item => item.status === 'RETRYING').length);
  readonly deadLetterCount = computed(() => this.deliveries().filter(item => item.status === 'DEAD_LETTER').length);
  readonly deliveredCount = computed(() => this.deliveries().filter(item => item.status === 'DELIVERED').length);
  readonly filteredTemplates = computed(() => this.filterRows(this.templates(), item => [item.code, item.name, item.category, item.defaultChannel]));
  readonly filteredDeliveries = computed(() => this.filterRows(this.deliveries(), item => [item.subject, item.channel, item.status, item.recipientId, item.providerCode, item.errorMessage]));
  readonly filteredProviders = computed(() => this.filterRows(this.providers(), item => [item.code, item.name, item.channel, item.status]));

  readonly statusOptions = [{ label: 'Hiệu lực', value: 'ACTIVE' }, { label: 'Ngưng hiệu lực', value: 'INACTIVE' }];
  readonly channelOptions = [{ label: 'In-app', value: 'IN_APP' }, { label: 'Email', value: 'EMAIL' }, { label: 'SMS', value: 'SMS' }, { label: 'Zalo OA', value: 'ZALO_OA' }, { label: 'Zalo ZNS', value: 'ZALO_ZNS' }, { label: 'Webhook', value: 'WEBHOOK' }];
  readonly approvalOptions = [{ label: 'Draft', value: 'DRAFT' }, { label: 'Approved', value: 'APPROVED' }];

  readonly templateForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    category: new FormControl('', { nonNullable: true }),
    defaultChannel: new FormControl('IN_APP', { nonNullable: true }),
    description: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  readonly versionForm = new FormGroup({
    version: new FormControl(1, { nonNullable: true }),
    channel: new FormControl('IN_APP', { nonNullable: true }),
    language: new FormControl('vi', { nonNullable: true }),
    subject: new FormControl('', { nonNullable: true }),
    body: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    payloadSchemaJson: new FormControl('{}', { nonNullable: true }),
    approvalStatus: new FormControl('DRAFT', { nonNullable: true }),
    changeNote: new FormControl('', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  readonly providerForm = new FormGroup({
    code: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    channel: new FormControl('IN_APP', { nonNullable: true }),
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    priority: new FormControl(100, { nonNullable: true }),
    rateLimitPerMinute: new FormControl(0, { nonNullable: true }),
    optionsJson: new FormControl('{}', { nonNullable: true }),
    status: new FormControl('ACTIVE', { nonNullable: true }),
  });

  setTab(tab: NotificationTab): void { this.activeTab.set(tab); }
  reload(): void { this.templateResource.reload(); this.deliveryResource.reload(); this.providerResource.reload(); }
  statusSeverity = statusSeverity;

  openTemplate(item?: NotificationTemplate): void {
    this.selectedTemplate.set(item ?? null);
    this.templateForm.reset(item ? { ...item } : { code: '', name: '', category: '', defaultChannel: 'IN_APP', description: '', status: 'ACTIVE' });
    this.dialogVisible.set(true);
  }

  saveTemplate(): void {
    if (this.templateForm.invalid) return;
    const selected = this.selectedTemplate();
    this.runSave(selected ? this.service.updateTemplate(selected.id, this.templateForm.getRawValue()) : this.service.createTemplate(this.templateForm.getRawValue()), () => {
      this.dialogVisible.set(false);
      this.templateResource.reload();
    });
  }

  openVersions(item: NotificationTemplate): void {
    this.selectedTemplateForVersions.set(item);
    this.loadVersions(item.id);
    this.versionDialogVisible.set(true);
  }

  addVersion(): void {
    const template = this.selectedTemplateForVersions();
    if (!template || this.versionForm.invalid) return;
    this.runSave(this.service.createTemplateVersion(template.id, this.versionForm.getRawValue()), () => this.loadVersions(template.id));
  }

  approveVersion(item: NotificationTemplateVersion): void {
    const template = this.selectedTemplateForVersions();
    this.runSave(this.service.approveTemplateVersion(item.id), () => {
      if (template) this.loadVersions(template.id);
    });
  }

  retryDelivery(item: NotificationDelivery): void {
    this.runSave(this.service.retryDelivery(item.id), () => this.deliveryResource.reload());
  }

  runWorkerOnce(): void {
    this.runSave(this.service.runWorkerOnce(), result => {
      this.messageService.add({ severity: 'info', summary: 'Worker', detail: `Đã xử lý ${result.processed}, lỗi ${result.failed}` });
      this.deliveryResource.reload();
    }, false);
  }

  openProvider(item?: ProviderConfig): void {
    this.selectedProvider.set(item ?? null);
    this.providerForm.reset(item ? { ...item } : { code: '', channel: 'IN_APP', name: '', priority: 100, rateLimitPerMinute: 0, optionsJson: '{}', status: 'ACTIVE' });
    this.providerDialogVisible.set(true);
  }

  saveProvider(): void {
    if (this.providerForm.invalid) return;
    const raw = this.providerForm.getRawValue();
    this.runSave(this.service.upsertProviderConfig(raw.code, raw), () => {
      this.providerDialogVisible.set(false);
      this.providerResource.reload();
    });
  }

  private loadVersions(templateId: string): void {
    this.service.listTemplateVersions(templateId, { pageSize: 200 }).subscribe({
      next: items => this.templateVersions.set(items),
      error: () => this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không tải được phiên bản mẫu' }),
    });
  }

  private runSave<T>(request: Observable<T>, onSuccess: (result: T) => void, showDefaultMessage = true): void {
    this.isSaving.set(true);
    request.subscribe({
      next: result => {
        if (showDefaultMessage) this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Đã lưu dữ liệu' });
        onSuccess(result);
        this.isSaving.set(false);
      },
      error: () => {
        this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: 'Không thể thực hiện thao tác' });
        this.isSaving.set(false);
      },
    });
  }

  private filterRows<T>(rows: T[], values: (item: T) => string[]): T[] {
    const keyword = this.keyword().trim().toLocaleLowerCase('vi');
    return rows.filter(item => !keyword || values(item).some(value => value.toLocaleLowerCase('vi').includes(keyword)));
  }
}
