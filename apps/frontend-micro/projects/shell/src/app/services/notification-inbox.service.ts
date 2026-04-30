import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, computed, inject, signal } from '@angular/core';
import { firstValueFrom } from 'rxjs';

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
  private readonly baseUrl = '/api/v1/notifications/in-app';

  readonly items = signal<InAppNotification[]>([]);
  readonly unreadCount = signal(0);
  readonly isLoading = signal(false);
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

  clear(): void {
    this.items.set([]);
    this.unreadCount.set(0);
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
