import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { BankBranch, ListOptions, PageResponse, PaymentNetwork } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class PaymentNetworksService {
  private http = inject(HttpClient);

  listBranches(options: ListOptions = {}): Observable<PageResponse<BankBranch>> {
    return this.http.get<any>('/api/v1/mdm/bank-branches', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.bank_branches ?? resp.bankBranches ?? []).map((item: any) => this.toBranch(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createBranch(item: Partial<BankBranch>): Observable<BankBranch> {
    return this.http.post<any>('/api/v1/mdm/bank-branches', { bank_branch: this.fromBranch(item) }).pipe(map(resp => this.toBranch(resp)));
  }
  updateBranch(id: string, item: Partial<BankBranch>): Observable<BankBranch> {
    return this.http.put<any>(`/api/v1/mdm/bank-branches/${encodeURIComponent(id)}`, { bank_branch: this.fromBranch(item) }).pipe(map(resp => this.toBranch(resp)));
  }
  deleteBranch(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/bank-branches/${encodeURIComponent(id)}`);
  }

  listNetworks(options: ListOptions = {}): Observable<PageResponse<PaymentNetwork>> {
    return this.http.get<any>('/api/v1/mdm/payment-networks', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.payment_networks ?? resp.paymentNetworks ?? []).map((item: any) => this.toNetwork(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createNetwork(item: Partial<PaymentNetwork>): Observable<PaymentNetwork> {
    return this.http.post<any>('/api/v1/mdm/payment-networks', { payment_network: this.fromNetwork(item) }).pipe(map(resp => this.toNetwork(resp)));
  }
  updateNetwork(id: string, item: Partial<PaymentNetwork>): Observable<PaymentNetwork> {
    return this.http.put<any>(`/api/v1/mdm/payment-networks/${encodeURIComponent(id)}`, { payment_network: this.fromNetwork(item) }).pipe(map(resp => this.toNetwork(resp)));
  }
  deleteNetwork(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/payment-networks/${encodeURIComponent(id)}`);
  }

  private toBranch(item: any): BankBranch {
    return { id: item.id ?? '', institutionCode: item.institution_code ?? item.institutionCode ?? '', code: item.code ?? '', name: item.name ?? '', branchType: item.branch_type ?? item.branchType ?? 'BRANCH', address: item.address ?? '', provinceCode: item.province_code ?? item.provinceCode ?? '', phone: item.phone ?? '', swiftCode: item.swift_code ?? item.swiftCode ?? '', napasCode: item.napas_code ?? item.napasCode ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromBranch(item: Partial<BankBranch>): any {
    return { id: item.id, institution_code: item.institutionCode, code: item.code, name: item.name, branch_type: item.branchType, address: item.address, province_code: item.provinceCode, phone: item.phone, swift_code: item.swiftCode, napas_code: item.napasCode, status: item.status };
  }
  private toNetwork(item: any): PaymentNetwork {
    return { id: item.id ?? '', code: item.code ?? '', name: item.name ?? '', networkType: item.network_type ?? item.networkType ?? 'DOMESTIC', clearingMethod: item.clearing_method ?? item.clearingMethod ?? '', settlementCurrency: item.settlement_currency ?? item.settlementCurrency ?? 'VND', operator: item.operator ?? '', availability: item.availability ?? '24X7', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromNetwork(item: Partial<PaymentNetwork>): any {
    return { id: item.id, code: item.code, name: item.name, network_type: item.networkType, clearing_method: item.clearingMethod, settlement_currency: item.settlementCurrency, operator: item.operator, availability: item.availability, description: item.description, status: item.status };
  }
}
