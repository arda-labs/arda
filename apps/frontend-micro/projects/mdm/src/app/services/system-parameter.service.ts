import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { ListOptions, SystemParameter } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class SystemParameterService {
  private http = inject(HttpClient);

  list(options: ListOptions = {}): Observable<SystemParameter[]> {
    return this.http.get<any>('/api/v1/mdm/system-parameters', { params: buildParams(options) }).pipe(
      map(resp => (resp.parameters ?? []).map((item: any) => this.toModel(item))),
    );
  }

  create(parameter: Partial<SystemParameter>): Observable<SystemParameter> {
    return this.http.post<any>('/api/v1/mdm/system-parameters', { parameter: this.fromModel(parameter) }).pipe(map(resp => this.toModel(resp)));
  }

  update(key: string, parameter: Partial<SystemParameter>): Observable<SystemParameter> {
    return this.http.put<any>(`/api/v1/mdm/system-parameters/${encodeURIComponent(key)}`, { parameter: this.fromModel(parameter) }).pipe(map(resp => this.toModel(resp)));
  }

  delete(key: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/system-parameters/${encodeURIComponent(key)}`);
  }

  private toModel(item: any): SystemParameter {
    return {
      id: item.id ?? '',
      key: item.key ?? '',
      name: item.name ?? '',
      groupCode: item.group_code ?? item.groupCode ?? '',
      valueType: item.value_type ?? item.valueType ?? 'STRING',
      valueText: item.value_text ?? item.valueText ?? '',
      valueNumber: Number(item.value_number ?? item.valueNumber ?? 0),
      valueBoolean: item.value_boolean ?? item.valueBoolean ?? false,
      valueJson: item.value_json ?? item.valueJson ?? '{}',
      defaultValue: item.default_value ?? item.defaultValue ?? '',
      isSecret: item.is_secret ?? item.isSecret ?? false,
      isEditable: item.is_editable ?? item.isEditable ?? true,
      isSystem: item.is_system ?? item.isSystem ?? false,
      validationRuleJson: item.validation_rule_json ?? item.validationRuleJson ?? '{}',
      description: item.description ?? '',
      status: item.status ?? 'ACTIVE',
      updatedBy: item.updated_by ?? item.updatedBy ?? '',
    };
  }

  private fromModel(item: Partial<SystemParameter>): any {
    return {
      id: item.id,
      key: item.key,
      name: item.name,
      group_code: item.groupCode,
      value_type: item.valueType,
      value_text: item.valueText,
      value_number: item.valueNumber,
      value_boolean: item.valueBoolean,
      value_json: item.valueJson || '{}',
      default_value: item.defaultValue,
      is_secret: item.isSecret,
      is_editable: item.isEditable,
      is_system: item.isSystem,
      validation_rule_json: item.validationRuleJson || '{}',
      description: item.description,
      status: item.status,
      updated_by: item.updatedBy,
    };
  }
}

