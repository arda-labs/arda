import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { AdministrativeUnit, ListOptions, PageResponse } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class AdministrativeUnitService {
  private http = inject(HttpClient);

  list(options: ListOptions = {}): Observable<PageResponse<AdministrativeUnit>> {
    return this.http.get<any>('/api/v1/mdm/administrative-units', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.units ?? []).map((item: any) => this.toModel(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  listProvinces(options: ListOptions = {}): Observable<AdministrativeUnit[]> {
    return this.http.get<any>('/api/v1/mdm/provinces', { params: buildParams(options) }).pipe(
      map(resp => (resp.units ?? []).map((item: any) => this.toModel(item))),
    );
  }

  listWards(provinceId: string, options: ListOptions = {}): Observable<AdministrativeUnit[]> {
    return this.http.get<any>(`/api/v1/mdm/provinces/${encodeURIComponent(provinceId)}/wards`, { params: buildParams(options) }).pipe(
      map(resp => (resp.units ?? []).map((item: any) => this.toModel(item))),
    );
  }

  create(unit: Partial<AdministrativeUnit>): Observable<AdministrativeUnit> {
    return this.http.post<any>('/api/v1/mdm/administrative-units', { unit: this.fromModel(unit) }).pipe(
      map(resp => this.toModel(resp)),
    );
  }

  update(id: string, unit: Partial<AdministrativeUnit>): Observable<AdministrativeUnit> {
    return this.http.put<any>(`/api/v1/mdm/administrative-units/${encodeURIComponent(id)}`, { unit: this.fromModel(unit) }).pipe(
      map(resp => this.toModel(resp)),
    );
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/administrative-units/${encodeURIComponent(id)}`);
  }

  private toModel(item: any): AdministrativeUnit {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      fullName: item.full_name ?? item.fullName ?? '',
      shortName: item.short_name ?? item.shortName ?? '',
      level: item.level ?? '',
      unitType: item.unit_type ?? item.unitType ?? '',
      parentId: item.parent_id ?? item.parentId ?? '',
      path: item.path ?? '',
      sortOrder: Number(item.sort_order ?? item.sortOrder ?? 0),
      latitude: Number(item.latitude ?? 0),
      longitude: Number(item.longitude ?? 0),
      status: item.status ?? 'ACTIVE',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
      source: item.source ?? '',
      metadataJson: item.metadata_json ?? item.metadataJson ?? '{}',
    };
  }

  private fromModel(item: Partial<AdministrativeUnit>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      full_name: item.fullName,
      short_name: item.shortName,
      level: item.level,
      unit_type: item.unitType,
      parent_id: item.parentId,
      path: item.path,
      sort_order: item.sortOrder,
      latitude: item.latitude,
      longitude: item.longitude,
      status: item.status,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
      source: item.source,
      metadata_json: item.metadataJson || '{}',
    };
  }
}

