import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, computed, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { AuthService } from './auth.service';

export interface InAppNotification {
  id: string;
  deliveryId: string;
  recipientType: string;
  recipientId: string;
  title: string;
  body: string;
  dataJson: string;
  status: string;
  readAt: string;
  createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class NotificationInboxService {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private readonly baseUrl = '/api/v1/notifications/in-app';
  private streamRecipientId = '';
  private streamAbort: AbortController | null = null;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  readonly items = signal<InAppNotification[]>([]);
  readonly unreadCount = signal(0);
  readonly isLoading = signal(false);
  readonly streamConnected = signal(false);
  readonly hasUnread = computed(() => this.unreadCount() > 0);

  async refresh(recipientId: string): Promise<void> {
    if (!recipientId) return;
    this.isLoading.set(true);
    try {
      await Promise.all([
        this.refreshUnreadCount(recipientId),
        this.refreshItems(recipientId),
      ]);
    } finally {
      this.isLoading.set(false);
    }
  }

  async refreshUnreadCount(recipientId: string): Promise<void> {
    if (!recipientId) return;
    const params = new HttpParams()
      .set('recipient_type', 'USER')
      .set('recipient_id', recipientId);
    const resp = await firstValueFrom(this.http.get<any>(`${this.baseUrl}/unread-count`, { params }));
    this.unreadCount.set(resp.count ?? 0);
  }

  async refreshItems(recipientId: string): Promise<void> {
    if (!recipientId) return;
    const params = new HttpParams()
      .set('recipient_type', 'USER')
      .set('recipient_id', recipientId)
      .set('page_size', '10');
    const resp = await firstValueFrom(this.http.get<any>(this.baseUrl, { params }));
    this.items.set((resp.notifications ?? []).map((item: any) => this.toNotification(item)));
  }

  async markRead(item: InAppNotification): Promise<void> {
    await firstValueFrom(this.http.post(`${this.baseUrl}/${encodeURIComponent(item.id)}/read`, { actor: 'SHELL' }));
    this.items.update(items => items.map(current => current.id === item.id ? { ...current, status: 'READ' } : current));
    this.unreadCount.update(count => item.status === 'UNREAD' ? Math.max(0, count - 1) : count);
  }

  async markAllRead(recipientId: string): Promise<void> {
    if (!recipientId) return;
    await firstValueFrom(this.http.post(`${this.baseUrl}/read-all`, { recipient_type: 'USER', recipient_id: recipientId, actor: 'SHELL' }));
    this.items.update(items => items.map(item => ({ ...item, status: 'READ' })));
    this.unreadCount.set(0);
  }

  startRealtime(recipientId: string): void {
    if (!recipientId || this.streamRecipientId === recipientId) return;
    this.stopRealtime();
    this.streamRecipientId = recipientId;
    this.openStream(recipientId);
  }

  stopRealtime(): void {
    this.streamRecipientId = '';
    this.streamConnected.set(false);
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.streamAbort) {
      this.streamAbort.abort();
      this.streamAbort = null;
    }
  }

  clear(): void {
    this.stopRealtime();
    this.items.set([]);
    this.unreadCount.set(0);
  }

  private async openStream(recipientId: string): Promise<void> {
    if (typeof window === 'undefined') return;
    const controller = new AbortController();
    this.streamAbort = controller;
    try {
      const token = await this.authService.getAccessToken();
      const params = new URLSearchParams({ recipient_type: 'USER', recipient_id: recipientId });
      const response = await fetch(`${this.baseUrl}/stream?${params.toString()}`, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        signal: controller.signal,
      });
      if (!response.ok || !response.body) {
        throw new Error(`SSE failed: ${response.status}`);
      }
      this.streamConnected.set(true);
      await this.consumeStream(response.body, recipientId, controller.signal);
    } catch {
      if (!controller.signal.aborted && this.streamRecipientId === recipientId) {
        this.streamConnected.set(false);
        this.reconnectTimer = setTimeout(() => this.openStream(recipientId), 5000);
      }
    }
  }

  private async consumeStream(body: ReadableStream<Uint8Array>, recipientId: string, signal: AbortSignal): Promise<void> {
    const reader = body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    while (!signal.aborted) {
      const { done, value } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      const chunks = buffer.split('\n\n');
      buffer = chunks.pop() ?? '';
      for (const chunk of chunks) {
        this.handleStreamChunk(chunk, recipientId);
      }
    }
    if (!signal.aborted) {
      throw new Error('SSE closed');
    }
  }

  private handleStreamChunk(chunk: string, recipientId: string): void {
    const eventLine = chunk.split('\n').find(line => line.startsWith('event:'));
    const event = eventLine?.replace('event:', '').trim() ?? '';
    if (!event || event === 'heartbeat') return;
    this.refreshUnreadCount(recipientId).catch(() => undefined);
    if (event === 'notification.created') {
      this.refreshItems(recipientId).catch(() => undefined);
    }
  }

  private toNotification(item: any): InAppNotification {
    return {
      id: item.id ?? '',
      deliveryId: item.delivery_id ?? item.deliveryId ?? '',
      recipientType: item.recipient_type ?? item.recipientType ?? 'USER',
      recipientId: item.recipient_id ?? item.recipientId ?? '',
      title: item.title ?? '',
      body: item.body ?? '',
      dataJson: item.data_json ?? item.dataJson ?? '{}',
      status: item.status ?? 'UNREAD',
      readAt: item.read_at ?? item.readAt ?? '',
      createdAt: item.created_at ?? item.createdAt ?? '',
    };
  }
}
