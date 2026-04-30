import { HttpParams } from '@angular/common/http';
import { ListOptions } from '../models/notification.models';

export function buildParams(options: ListOptions = {}): HttpParams {
  let params = new HttpParams().set('page_size', String(options.pageSize ?? 200));
  const entries: Record<string, string | undefined> = {
    status: options.status,
    channel: options.channel,
    keyword: options.keyword,
    recipient_id: options.recipientId,
    page_token: options.pageToken,
  };
  for (const [key, value] of Object.entries(entries)) {
    if (value) params = params.set(key, value);
  }
  return params;
}

export function statusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
  switch (status) {
    case 'ACTIVE':
    case 'APPROVED':
    case 'DELIVERED':
      return 'success';
    case 'QUEUED':
    case 'CLAIMED':
    case 'RETRYING':
      return 'warn';
    case 'DEAD_LETTER':
    case 'FAILED':
    case 'DELETED':
      return 'danger';
    case 'INACTIVE':
    case 'DRAFT':
      return 'secondary';
    default:
      return 'info';
  }
}
