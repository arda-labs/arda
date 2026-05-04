import { HttpParams } from '@angular/common/http';

export interface BpmListOptions {
  pageSize?: number;
  pageToken?: string;
  keyword?: string;
  status?: string;
  module?: string;
  processDefinitionId?: string;
  includeInactive?: boolean;
  fromDate?: string;
  toDate?: string;
}

export function buildParams(options: BpmListOptions = {}): HttpParams {
  let params = new HttpParams().set('page_size', String(options.pageSize ?? 100));
  const entries: Record<string, string | undefined> = {
    page_token: options.pageToken,
    keyword: options.keyword,
    status: options.status,
    module: options.module,
    process_definition_id: options.processDefinitionId,
    include_inactive: options.includeInactive != null ? String(options.includeInactive) : undefined,
    from_date: options.fromDate,
    to_date: options.toDate,
  };
  for (const [key, value] of Object.entries(entries)) {
    if (value != null && value !== '') params = params.set(key, value);
  }
  return params;
}

export function snakeToCamel<T>(obj: unknown): T {
  if (obj === null || obj === undefined) return obj as T;
  if (Array.isArray(obj)) return obj.map(snakeToCamel) as T;
  if (typeof obj === 'object') {
    const out: Record<string, unknown> = {};
    for (const [key, val] of Object.entries(obj)) {
      const camel = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
      out[camel] = snakeToCamel(val as Record<string, unknown>);
    }
    return out as T;
  }
  return obj as T;
}

export function extractPageResponse<T>(itemsKey: string) {
  return (r: Record<string, unknown>) => ({
    items: (r[itemsKey] ?? []) as T[],
    nextPageToken: (r['next_page_token'] ?? r['nextPageToken'] ?? '') as string,
  });
}
