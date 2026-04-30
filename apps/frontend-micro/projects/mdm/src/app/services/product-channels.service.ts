import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { BankingProduct, ListOptions, PageResponse, ProductChannelRule, ServiceChannel } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class ProductChannelsService {
  private http = inject(HttpClient);

  listProducts(options: ListOptions = {}): Observable<PageResponse<BankingProduct>> {
    return this.http.get<any>('/api/v1/mdm/banking-products', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.banking_products ?? resp.bankingProducts ?? []).map((item: any) => this.toProduct(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createProduct(item: Partial<BankingProduct>): Observable<BankingProduct> {
    return this.http.post<any>('/api/v1/mdm/banking-products', { banking_product: this.fromProduct(item) }).pipe(map(resp => this.toProduct(resp)));
  }
  updateProduct(id: string, item: Partial<BankingProduct>): Observable<BankingProduct> {
    return this.http.put<any>(`/api/v1/mdm/banking-products/${encodeURIComponent(id)}`, { banking_product: this.fromProduct(item) }).pipe(map(resp => this.toProduct(resp)));
  }
  deleteProduct(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/banking-products/${encodeURIComponent(id)}`);
  }

  listChannels(options: ListOptions = {}): Observable<PageResponse<ServiceChannel>> {
    return this.http.get<any>('/api/v1/mdm/service-channels', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.service_channels ?? resp.serviceChannels ?? []).map((item: any) => this.toChannel(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createChannel(item: Partial<ServiceChannel>): Observable<ServiceChannel> {
    return this.http.post<any>('/api/v1/mdm/service-channels', { service_channel: this.fromChannel(item) }).pipe(map(resp => this.toChannel(resp)));
  }
  updateChannel(id: string, item: Partial<ServiceChannel>): Observable<ServiceChannel> {
    return this.http.put<any>(`/api/v1/mdm/service-channels/${encodeURIComponent(id)}`, { service_channel: this.fromChannel(item) }).pipe(map(resp => this.toChannel(resp)));
  }
  deleteChannel(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/service-channels/${encodeURIComponent(id)}`);
  }

  listRules(options: ListOptions = {}): Observable<PageResponse<ProductChannelRule>> {
    return this.http.get<any>('/api/v1/mdm/product-channel-rules', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.product_channel_rules ?? resp.productChannelRules ?? []).map((item: any) => this.toRule(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createRule(item: Partial<ProductChannelRule>): Observable<ProductChannelRule> {
    return this.http.post<any>('/api/v1/mdm/product-channel-rules', { product_channel_rule: this.fromRule(item) }).pipe(map(resp => this.toRule(resp)));
  }
  updateRule(id: string, item: Partial<ProductChannelRule>): Observable<ProductChannelRule> {
    return this.http.put<any>(`/api/v1/mdm/product-channel-rules/${encodeURIComponent(id)}`, { product_channel_rule: this.fromRule(item) }).pipe(map(resp => this.toRule(resp)));
  }
  deleteRule(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/product-channel-rules/${encodeURIComponent(id)}`);
  }

  private toProduct(item: any): BankingProduct {
    return { id: item.id ?? '', code: item.code ?? '', name: item.name ?? '', productType: item.product_type ?? item.productType ?? 'ACCOUNT', category: item.category ?? '', customerSegment: item.customer_segment ?? item.customerSegment ?? '', currency: item.currency ?? '', effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '', effectiveTo: item.effective_to ?? item.effectiveTo ?? '', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromProduct(item: Partial<BankingProduct>): any {
    return { id: item.id, code: item.code, name: item.name, product_type: item.productType, category: item.category, customer_segment: item.customerSegment, currency: item.currency, effective_from: item.effectiveFrom, effective_to: item.effectiveTo, description: item.description, status: item.status };
  }
  private toChannel(item: any): ServiceChannel {
    return { id: item.id ?? '', code: item.code ?? '', name: item.name ?? '', channelType: item.channel_type ?? item.channelType ?? 'DIGITAL', availability: item.availability ?? '24X7', timezone: item.timezone ?? 'Asia/Ho_Chi_Minh', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromChannel(item: Partial<ServiceChannel>): any {
    return { id: item.id, code: item.code, name: item.name, channel_type: item.channelType, availability: item.availability, timezone: item.timezone, description: item.description, status: item.status };
  }
  private toRule(item: any): ProductChannelRule {
    return { id: item.id ?? '', productCode: item.product_code ?? item.productCode ?? '', channelCode: item.channel_code ?? item.channelCode ?? '', transactionType: item.transaction_type ?? item.transactionType ?? '', enabled: item.enabled ?? true, priority: item.priority ?? 100, feeScheduleCode: item.fee_schedule_code ?? item.feeScheduleCode ?? '', limitProfileCode: item.limit_profile_code ?? item.limitProfileCode ?? '', effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '', effectiveTo: item.effective_to ?? item.effectiveTo ?? '', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromRule(item: Partial<ProductChannelRule>): any {
    return { id: item.id, product_code: item.productCode, channel_code: item.channelCode, transaction_type: item.transactionType, enabled: item.enabled, priority: item.priority, fee_schedule_code: item.feeScheduleCode, limit_profile_code: item.limitProfileCode, effective_from: item.effectiveFrom, effective_to: item.effectiveTo, description: item.description, status: item.status };
  }
}
