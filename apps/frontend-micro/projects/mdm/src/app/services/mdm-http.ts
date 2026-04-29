import { HttpParams } from '@angular/common/http';
import { ListOptions } from '../models/mdm.models';

export function buildParams(options: ListOptions = {}): HttpParams {
  let params = new HttpParams().set('page_size', String(options.pageSize ?? 500));
  const entries: Record<string, string | undefined> = {
    status: options.status,
    keyword: options.keyword,
    page_token: options.pageToken,
    parent_id: options.parentId,
    level: options.level,
    area_type_id: options.areaTypeId,
    group_code: options.groupCode,
  };
  for (const [key, value] of Object.entries(entries)) {
    if (value) params = params.set(key, value);
  }
  return params;
}

export function statusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
  switch (status) {
    case 'ACTIVE': return 'success';
    case 'INACTIVE': return 'secondary';
    case 'MERGED': return 'info';
    case 'DELETED': return 'danger';
    default: return 'warn';
  }
}

