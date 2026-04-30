import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { CreditInstitution, ListOptions, PageResponse } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class CreditInstitutionService {
  private http = inject(HttpClient);

  list(options: ListOptions = {}): Observable<PageResponse<CreditInstitution>> {
    return this.http.get<any>('/api/v1/mdm/credit-institutions', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.credit_institutions ?? resp.creditInstitutions ?? []).map((item: any) => this.toModel(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  create(item: Partial<CreditInstitution>): Observable<CreditInstitution> {
    return this.http.post<any>('/api/v1/mdm/credit-institutions', { credit_institution: this.fromModel(item) }).pipe(
      map(resp => this.toModel(resp)),
    );
  }

  update(id: string, item: Partial<CreditInstitution>): Observable<CreditInstitution> {
    return this.http.put<any>(`/api/v1/mdm/credit-institutions/${encodeURIComponent(id)}`, { credit_institution: this.fromModel(item) }).pipe(
      map(resp => this.toModel(resp)),
    );
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/credit-institutions/${encodeURIComponent(id)}`);
  }

  private toModel(item: any): CreditInstitution {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      shortName: item.short_name ?? item.shortName ?? '',
      address: item.address ?? '',
      phone: item.phone ?? '',
      email: item.email ?? '',
      licenseNumber: item.license_number ?? item.licenseNumber ?? '',
      issuedDate: item.issued_date ?? item.issuedDate ?? '',
      taxCode: item.tax_code ?? item.taxCode ?? '',
      website: item.website ?? '',
      note: item.note ?? '',
      status: item.status ?? 'ACTIVE',
    };
  }

  private fromModel(item: Partial<CreditInstitution>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      short_name: item.shortName,
      address: item.address,
      phone: item.phone,
      email: item.email,
      license_number: item.licenseNumber,
      issued_date: item.issuedDate,
      tax_code: item.taxCode,
      website: item.website,
      note: item.note,
      status: item.status,
    };
  }
}
