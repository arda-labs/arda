import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { Area, AreaAdministrativeUnit, AreaType, ListOptions } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class AreaService {
  private http = inject(HttpClient);

  listTypes(options: ListOptions = {}): Observable<AreaType[]> {
    return this.http.get<any>('/api/v1/mdm/area-types', { params: buildParams(options) }).pipe(
      map(resp => (resp.area_types ?? resp.areaTypes ?? []).map((item: any) => this.toAreaType(item))),
    );
  }

  createType(areaType: Partial<AreaType>): Observable<AreaType> {
    return this.http.post<any>('/api/v1/mdm/area-types', { area_type: this.fromAreaType(areaType) }).pipe(map(resp => this.toAreaType(resp)));
  }

  updateType(id: string, areaType: Partial<AreaType>): Observable<AreaType> {
    return this.http.put<any>(`/api/v1/mdm/area-types/${encodeURIComponent(id)}`, { area_type: this.fromAreaType(areaType) }).pipe(map(resp => this.toAreaType(resp)));
  }

  deleteType(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/area-types/${encodeURIComponent(id)}`);
  }

  listAreas(options: ListOptions = {}): Observable<Area[]> {
    return this.http.get<any>('/api/v1/mdm/areas', { params: buildParams(options) }).pipe(
      map(resp => (resp.areas ?? []).map((item: any) => this.toArea(item))),
    );
  }

  createArea(area: Partial<Area>): Observable<Area> {
    return this.http.post<any>('/api/v1/mdm/areas', { area: this.fromArea(area) }).pipe(map(resp => this.toArea(resp)));
  }

  updateArea(id: string, area: Partial<Area>): Observable<Area> {
    return this.http.put<any>(`/api/v1/mdm/areas/${encodeURIComponent(id)}`, { area: this.fromArea(area) }).pipe(map(resp => this.toArea(resp)));
  }

  deleteArea(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/areas/${encodeURIComponent(id)}`);
  }

  listAdministrativeUnits(areaId: string): Observable<AreaAdministrativeUnit[]> {
    return this.http.get<any>(`/api/v1/mdm/areas/${encodeURIComponent(areaId)}/administrative-units`).pipe(
      map(resp => (resp.items ?? []).map((item: any) => this.toAreaAdministrativeUnit(item))),
    );
  }

  assignAdministrativeUnit(areaId: string, administrativeUnitId: string, scopeType: string): Observable<AreaAdministrativeUnit> {
    return this.http.post<any>(`/api/v1/mdm/areas/${encodeURIComponent(areaId)}/administrative-units`, {
      administrative_unit_id: administrativeUnitId,
      scope_type: scopeType,
    }).pipe(map(resp => this.toAreaAdministrativeUnit(resp)));
  }

  removeAdministrativeUnit(areaId: string, administrativeUnitId: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/areas/${encodeURIComponent(areaId)}/administrative-units/${encodeURIComponent(administrativeUnitId)}`);
  }

  private toAreaType(item: any): AreaType {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      description: item.description ?? '',
      allowHierarchy: item.allow_hierarchy ?? item.allowHierarchy ?? true,
      status: item.status ?? 'ACTIVE',
    };
  }

  private fromAreaType(item: Partial<AreaType>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      description: item.description,
      allow_hierarchy: item.allowHierarchy,
      status: item.status,
    };
  }

  private toArea(item: any): Area {
    return {
      id: item.id ?? '',
      areaTypeId: item.area_type_id ?? item.areaTypeId ?? '',
      areaTypeCode: item.area_type_code ?? item.areaTypeCode ?? '',
      parentId: item.parent_id ?? item.parentId ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      description: item.description ?? '',
      managerUserId: item.manager_user_id ?? item.managerUserId ?? '',
      status: item.status ?? 'ACTIVE',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
      metadataJson: item.metadata_json ?? item.metadataJson ?? '{}',
    };
  }

  private fromArea(item: Partial<Area>): any {
    return {
      id: item.id,
      area_type_id: item.areaTypeId,
      parent_id: item.parentId,
      code: item.code,
      name: item.name,
      description: item.description,
      manager_user_id: item.managerUserId,
      status: item.status,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
      metadata_json: item.metadataJson || '{}',
    };
  }

  private toAreaAdministrativeUnit(item: any): AreaAdministrativeUnit {
    return {
      id: item.id ?? '',
      areaId: item.area_id ?? item.areaId ?? '',
      administrativeUnitId: item.administrative_unit_id ?? item.administrativeUnitId ?? '',
      scopeType: item.scope_type ?? item.scopeType ?? 'INCLUDE',
      createdAt: item.created_at ?? item.createdAt,
    };
  }
}

