import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import {
  ListOptions,
  NotificationDelivery,
  NotificationTemplate,
  NotificationTemplateVersion,
  PageResponse,
  ProviderConfig,
} from '../models/notification.models';
import { buildParams } from './notification-http';

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private http = inject(HttpClient);
  private readonly baseUrl = '/api/v1/notifications';

  listTemplates(options: ListOptions = {}): Observable<PageResponse<NotificationTemplate>> {
    return this.http.get<any>(`${this.baseUrl}/templates`, { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.templates ?? []).map((item: any) => this.toTemplate(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  createTemplate(item: Partial<NotificationTemplate>): Observable<NotificationTemplate> {
    return this.http.post<any>(`${this.baseUrl}/templates`, { template: this.fromTemplate(item) }).pipe(map(resp => this.toTemplate(resp)));
  }

  updateTemplate(id: string, item: Partial<NotificationTemplate>): Observable<NotificationTemplate> {
    return this.http.put<any>(`${this.baseUrl}/templates/${encodeURIComponent(id)}`, { template: this.fromTemplate(item) }).pipe(map(resp => this.toTemplate(resp)));
  }

  listTemplateVersions(templateId: string, options: ListOptions = {}): Observable<NotificationTemplateVersion[]> {
    return this.http.get<any>(`${this.baseUrl}/templates/${encodeURIComponent(templateId)}/versions`, { params: buildParams(options) }).pipe(
      map(resp => (resp.versions ?? []).map((item: any) => this.toTemplateVersion(item))),
    );
  }

  createTemplateVersion(templateId: string, item: Partial<NotificationTemplateVersion>): Observable<NotificationTemplateVersion> {
    return this.http.post<any>(`${this.baseUrl}/templates/${encodeURIComponent(templateId)}/versions`, { version: this.fromTemplateVersion(item) }).pipe(map(resp => this.toTemplateVersion(resp)));
  }

  approveTemplateVersion(id: string, actor = 'UI'): Observable<NotificationTemplateVersion> {
    return this.http.post<any>(`${this.baseUrl}/template-versions/${encodeURIComponent(id)}/approve`, { actor }).pipe(map(resp => this.toTemplateVersion(resp)));
  }

  listDeliveries(options: ListOptions = {}): Observable<PageResponse<NotificationDelivery>> {
    return this.http.get<any>(`${this.baseUrl}/deliveries`, { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.deliveries ?? []).map((item: any) => this.toDelivery(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  retryDelivery(id: string): Observable<NotificationDelivery> {
    return this.http.post<any>(`${this.baseUrl}/deliveries/${encodeURIComponent(id)}/retry`, { actor: 'UI' }).pipe(map(resp => this.toDelivery(resp)));
  }

  runWorkerOnce(): Observable<{ processed: number; failed: number }> {
    return this.http.post<any>(`${this.baseUrl}/deliveries/run-once`, { worker_id: 'ui-manual', batch_size: 20 }).pipe(
      map(resp => ({ processed: resp.processed ?? 0, failed: resp.failed ?? 0 })),
    );
  }

  listProviderConfigs(options: ListOptions = {}): Observable<ProviderConfig[]> {
    return this.http.get<any>(`${this.baseUrl}/provider-configs`, { params: buildParams(options) }).pipe(
      map(resp => (resp.provider_configs ?? resp.providerConfigs ?? []).map((item: any) => this.toProviderConfig(item))),
    );
  }

  upsertProviderConfig(code: string, item: Partial<ProviderConfig>): Observable<ProviderConfig> {
    return this.http.put<any>(`${this.baseUrl}/provider-configs/${encodeURIComponent(code)}`, { provider_config: this.fromProviderConfig(item) }).pipe(map(resp => this.toProviderConfig(resp)));
  }

  private toTemplate(item: any): NotificationTemplate {
    return { id: item.id ?? '', code: item.code ?? '', name: item.name ?? '', category: item.category ?? '', defaultChannel: item.default_channel ?? item.defaultChannel ?? 'IN_APP', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }

  private fromTemplate(item: Partial<NotificationTemplate>): any {
    return { id: item.id, code: item.code, name: item.name, category: item.category, default_channel: item.defaultChannel, description: item.description, status: item.status };
  }

  private toTemplateVersion(item: any): NotificationTemplateVersion {
    return { id: item.id ?? '', templateId: item.template_id ?? item.templateId ?? '', version: item.version ?? 1, channel: item.channel ?? 'IN_APP', language: item.language ?? 'vi', subject: item.subject ?? '', body: item.body ?? '', payloadSchemaJson: item.payload_schema_json ?? item.payloadSchemaJson ?? '{}', approvalStatus: item.approval_status ?? item.approvalStatus ?? 'DRAFT', approvedBy: item.approved_by ?? item.approvedBy ?? '', changeNote: item.change_note ?? item.changeNote ?? '', status: item.status ?? 'ACTIVE' };
  }

  private fromTemplateVersion(item: Partial<NotificationTemplateVersion>): any {
    return { id: item.id, template_id: item.templateId, version: item.version, channel: item.channel, language: item.language, subject: item.subject, body: item.body, payload_schema_json: item.payloadSchemaJson, approval_status: item.approvalStatus, change_note: item.changeNote, status: item.status };
  }

  private toDelivery(item: any): NotificationDelivery {
    return { id: item.id ?? '', requestId: item.request_id ?? item.requestId ?? '', templateVersionId: item.template_version_id ?? item.templateVersionId ?? '', channel: item.channel ?? '', recipientType: item.recipient_type ?? item.recipientType ?? '', recipientId: item.recipient_id ?? item.recipientId ?? '', recipientAddress: item.recipient_address ?? item.recipientAddress ?? '', subject: item.subject ?? '', body: item.body ?? '', status: item.status ?? '', attemptCount: item.attempt_count ?? item.attemptCount ?? 0, maxAttempts: item.max_attempts ?? item.maxAttempts ?? 0, providerCode: item.provider_code ?? item.providerCode ?? '', providerMessageId: item.provider_message_id ?? item.providerMessageId ?? '', providerResponseJson: item.provider_response_json ?? item.providerResponseJson ?? '{}', errorMessage: item.error_message ?? item.errorMessage ?? '', priority: item.priority ?? 100, createdAt: item.created_at ?? item.createdAt ?? '', updatedAt: item.updated_at ?? item.updatedAt ?? '' };
  }

  private toProviderConfig(item: any): ProviderConfig {
    return { id: item.id ?? '', code: item.code ?? '', channel: item.channel ?? 'IN_APP', name: item.name ?? '', priority: item.priority ?? 100, rateLimitPerMinute: item.rate_limit_per_minute ?? item.rateLimitPerMinute ?? 0, optionsJson: item.options_json ?? item.optionsJson ?? '{}', status: item.status ?? 'ACTIVE' };
  }

  private fromProviderConfig(item: Partial<ProviderConfig>): any {
    return { id: item.id, code: item.code, channel: item.channel, name: item.name, priority: item.priority, rate_limit_per_minute: item.rateLimitPerMinute, options_json: item.optionsJson, status: item.status };
  }
}
