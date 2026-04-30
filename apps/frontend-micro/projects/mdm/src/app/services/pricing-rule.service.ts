import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { FeeSchedule, ListOptions, PageResponse, StandardLimit, TaxRule } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class PricingRuleService {
  private http = inject(HttpClient);

  listFeeSchedules(options: ListOptions = {}): Observable<PageResponse<FeeSchedule>> {
    return this.http.get<any>('/api/v1/mdm/fee-schedules', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.fee_schedules ?? resp.feeSchedules ?? []).map((item: any) => this.toFeeSchedule(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  createFeeSchedule(item: Partial<FeeSchedule>): Observable<FeeSchedule> {
    return this.http.post<any>('/api/v1/mdm/fee-schedules', { fee_schedule: this.fromFeeSchedule(item) }).pipe(map(resp => this.toFeeSchedule(resp)));
  }

  updateFeeSchedule(id: string, item: Partial<FeeSchedule>): Observable<FeeSchedule> {
    return this.http.put<any>(`/api/v1/mdm/fee-schedules/${encodeURIComponent(id)}`, { fee_schedule: this.fromFeeSchedule(item) }).pipe(map(resp => this.toFeeSchedule(resp)));
  }

  deleteFeeSchedule(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/fee-schedules/${encodeURIComponent(id)}`);
  }

  approveFeeSchedule(id: string, note = ''): Observable<FeeSchedule> {
    return this.http.post<any>(`/api/v1/mdm/fee-schedules/${encodeURIComponent(id)}/approve`, { actor: 'SYSTEM', note }).pipe(map(resp => this.toFeeSchedule(resp)));
  }

  listTaxRules(options: ListOptions = {}): Observable<PageResponse<TaxRule>> {
    return this.http.get<any>('/api/v1/mdm/tax-rules', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.tax_rules ?? resp.taxRules ?? []).map((item: any) => this.toTaxRule(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  createTaxRule(item: Partial<TaxRule>): Observable<TaxRule> {
    return this.http.post<any>('/api/v1/mdm/tax-rules', { tax_rule: this.fromTaxRule(item) }).pipe(map(resp => this.toTaxRule(resp)));
  }

  updateTaxRule(id: string, item: Partial<TaxRule>): Observable<TaxRule> {
    return this.http.put<any>(`/api/v1/mdm/tax-rules/${encodeURIComponent(id)}`, { tax_rule: this.fromTaxRule(item) }).pipe(map(resp => this.toTaxRule(resp)));
  }

  deleteTaxRule(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/tax-rules/${encodeURIComponent(id)}`);
  }

  approveTaxRule(id: string, note = ''): Observable<TaxRule> {
    return this.http.post<any>(`/api/v1/mdm/tax-rules/${encodeURIComponent(id)}/approve`, { actor: 'SYSTEM', note }).pipe(map(resp => this.toTaxRule(resp)));
  }

  listStandardLimits(options: ListOptions = {}): Observable<PageResponse<StandardLimit>> {
    return this.http.get<any>('/api/v1/mdm/standard-limits', { params: buildParams(options) }).pipe(
      map(resp => ({
        items: (resp.standard_limits ?? resp.standardLimits ?? []).map((item: any) => this.toStandardLimit(item)),
        nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '',
      })),
    );
  }

  createStandardLimit(item: Partial<StandardLimit>): Observable<StandardLimit> {
    return this.http.post<any>('/api/v1/mdm/standard-limits', { standard_limit: this.fromStandardLimit(item) }).pipe(map(resp => this.toStandardLimit(resp)));
  }

  updateStandardLimit(id: string, item: Partial<StandardLimit>): Observable<StandardLimit> {
    return this.http.put<any>(`/api/v1/mdm/standard-limits/${encodeURIComponent(id)}`, { standard_limit: this.fromStandardLimit(item) }).pipe(map(resp => this.toStandardLimit(resp)));
  }

  deleteStandardLimit(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/standard-limits/${encodeURIComponent(id)}`);
  }

  approveStandardLimit(id: string, note = ''): Observable<StandardLimit> {
    return this.http.post<any>(`/api/v1/mdm/standard-limits/${encodeURIComponent(id)}/approve`, { actor: 'SYSTEM', note }).pipe(map(resp => this.toStandardLimit(resp)));
  }

  private toFeeSchedule(item: any): FeeSchedule {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      feeType: item.fee_type ?? item.feeType ?? 'SERVICE_FEE',
      calculationMethod: item.calculation_method ?? item.calculationMethod ?? 'FIXED',
      currency: item.currency ?? 'VND',
      fixedAmount: item.fixed_amount ?? item.fixedAmount ?? 0,
      ratePercent: item.rate_percent ?? item.ratePercent ?? 0,
      minAmount: item.min_amount ?? item.minAmount ?? 0,
      maxAmount: item.max_amount ?? item.maxAmount ?? 0,
      channel: item.channel ?? '',
      productCode: item.product_code ?? item.productCode ?? '',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
      description: item.description ?? '',
      status: item.status ?? 'ACTIVE',
      approvalStatus: item.approval_status ?? item.approvalStatus ?? 'DRAFT',
      version: item.version ?? 1,
      approvedBy: item.approved_by ?? item.approvedBy ?? '',
      changeNote: item.change_note ?? item.changeNote ?? '',
    };
  }

  private fromFeeSchedule(item: Partial<FeeSchedule>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      fee_type: item.feeType,
      calculation_method: item.calculationMethod,
      currency: item.currency,
      fixed_amount: item.fixedAmount,
      rate_percent: item.ratePercent,
      min_amount: item.minAmount,
      max_amount: item.maxAmount,
      channel: item.channel,
      product_code: item.productCode,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
      description: item.description,
      status: item.status,
      approval_status: item.approvalStatus,
      version: item.version,
      approved_by: item.approvedBy,
      change_note: item.changeNote,
    };
  }

  private toTaxRule(item: any): TaxRule {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      taxType: item.tax_type ?? item.taxType ?? 'VAT',
      ratePercent: item.rate_percent ?? item.ratePercent ?? 0,
      inclusive: item.inclusive ?? false,
      jurisdiction: item.jurisdiction ?? 'VN',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
      description: item.description ?? '',
      status: item.status ?? 'ACTIVE',
      approvalStatus: item.approval_status ?? item.approvalStatus ?? 'DRAFT',
      version: item.version ?? 1,
      approvedBy: item.approved_by ?? item.approvedBy ?? '',
      changeNote: item.change_note ?? item.changeNote ?? '',
    };
  }

  private fromTaxRule(item: Partial<TaxRule>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      tax_type: item.taxType,
      rate_percent: item.ratePercent,
      inclusive: item.inclusive,
      jurisdiction: item.jurisdiction,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
      description: item.description,
      status: item.status,
      approval_status: item.approvalStatus,
      version: item.version,
      approved_by: item.approvedBy,
      change_note: item.changeNote,
    };
  }

  private toStandardLimit(item: any): StandardLimit {
    return {
      id: item.id ?? '',
      code: item.code ?? '',
      name: item.name ?? '',
      limitType: item.limit_type ?? item.limitType ?? 'TRANSACTION_AMOUNT',
      subjectType: item.subject_type ?? item.subjectType ?? 'CUSTOMER',
      currency: item.currency ?? 'VND',
      minAmount: item.min_amount ?? item.minAmount ?? 0,
      perTxnAmount: item.per_txn_amount ?? item.perTxnAmount ?? 0,
      dailyAmount: item.daily_amount ?? item.dailyAmount ?? 0,
      monthlyAmount: item.monthly_amount ?? item.monthlyAmount ?? 0,
      countLimit: item.count_limit ?? item.countLimit ?? 0,
      channel: item.channel ?? '',
      productCode: item.product_code ?? item.productCode ?? '',
      effectiveFrom: item.effective_from ?? item.effectiveFrom ?? '',
      effectiveTo: item.effective_to ?? item.effectiveTo ?? '',
      description: item.description ?? '',
      status: item.status ?? 'ACTIVE',
      approvalStatus: item.approval_status ?? item.approvalStatus ?? 'DRAFT',
      version: item.version ?? 1,
      approvedBy: item.approved_by ?? item.approvedBy ?? '',
      changeNote: item.change_note ?? item.changeNote ?? '',
    };
  }

  private fromStandardLimit(item: Partial<StandardLimit>): any {
    return {
      id: item.id,
      code: item.code,
      name: item.name,
      limit_type: item.limitType,
      subject_type: item.subjectType,
      currency: item.currency,
      min_amount: item.minAmount,
      per_txn_amount: item.perTxnAmount,
      daily_amount: item.dailyAmount,
      monthly_amount: item.monthlyAmount,
      count_limit: item.countLimit,
      channel: item.channel,
      product_code: item.productCode,
      effective_from: item.effectiveFrom,
      effective_to: item.effectiveTo,
      description: item.description,
      status: item.status,
      approval_status: item.approvalStatus,
      version: item.version,
      approved_by: item.approvedBy,
      change_note: item.changeNote,
    };
  }
}
