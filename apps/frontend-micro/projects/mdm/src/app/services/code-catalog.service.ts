import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { CodeItem, CodeSet, ListOptions } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class CodeCatalogService {
  private http = inject(HttpClient);

  listSets(options: ListOptions = {}): Observable<CodeSet[]> {
    return this.http.get<any>('/api/v1/mdm/code-sets', { params: buildParams(options) }).pipe(
      map(resp => (resp.code_sets ?? resp.codeSets ?? []).map((item: any) => this.toCodeSet(item))),
    );
  }

  createSet(codeSet: Partial<CodeSet>): Observable<CodeSet> {
    return this.http.post<any>('/api/v1/mdm/code-sets', { code_set: this.fromCodeSet(codeSet) }).pipe(map(resp => this.toCodeSet(resp)));
  }

  updateSet(id: string, codeSet: Partial<CodeSet>): Observable<CodeSet> {
    return this.http.put<any>(`/api/v1/mdm/code-sets/${encodeURIComponent(id)}`, { code_set: this.fromCodeSet(codeSet) }).pipe(map(resp => this.toCodeSet(resp)));
  }

  deleteSet(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/code-sets/${encodeURIComponent(id)}`);
  }

  listItems(codeSetCode: string, options: ListOptions = {}): Observable<CodeItem[]> {
    return this.http.get<any>(`/api/v1/mdm/code-sets/${encodeURIComponent(codeSetCode)}/items`, { params: buildParams(options) }).pipe(
      map(resp => (resp.code_items ?? resp.codeItems ?? []).map((item: any) => this.toCodeItem(item))),
    );
  }

  createItem(codeSetCode: string, item: Partial<CodeItem>): Observable<CodeItem> {
    return this.http.post<any>(`/api/v1/mdm/code-sets/${encodeURIComponent(codeSetCode)}/items`, { code_item: this.fromCodeItem(item) }).pipe(map(resp => this.toCodeItem(resp)));
  }

  updateItem(id: string, item: Partial<CodeItem>): Observable<CodeItem> {
    return this.http.put<any>(`/api/v1/mdm/code-items/${encodeURIComponent(id)}`, { code_item: this.fromCodeItem(item) }).pipe(map(resp => this.toCodeItem(resp)));
  }

  deleteItem(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/code-items/${encodeURIComponent(id)}`);
  }

  private toCodeSet(item: any): CodeSet {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      description: item.description ?? '',
      isSystem: item.is_system ?? item.isSystem ?? false,
      status: item.status ?? 'ACTIVE',
    };
  }

  private fromCodeSet(item: Partial<CodeSet>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      description: item.description,
      is_system: item.isSystem,
      status: item.status,
    };
  }

  private toCodeItem(item: any): CodeItem {
    return {
      id: item.id ?? '',
      codeSetId: item.code_set_id ?? item.codeSetId ?? '',
      codeSetCode: item.code_set_code ?? item.codeSetCode ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      value: item.value ?? '',
      parentId: item.parent_id ?? item.parentId ?? '',
      sortOrder: Number(item.sort_order ?? item.sortOrder ?? 0),
      color: item.color ?? '',
      icon: item.icon ?? '',
      metadataJson: item.metadata_json ?? item.metadataJson ?? '{}',
      isDefault: item.is_default ?? item.isDefault ?? false,
      isSystem: item.is_system ?? item.isSystem ?? false,
      status: item.status ?? 'ACTIVE',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
    };
  }

  private fromCodeItem(item: Partial<CodeItem>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      value: item.value,
      parent_id: item.parentId,
      sort_order: item.sortOrder,
      color: item.color,
      icon: item.icon,
      metadata_json: item.metadataJson || '{}',
      is_default: item.isDefault,
      is_system: item.isSystem,
      status: item.status,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
    };
  }
}

