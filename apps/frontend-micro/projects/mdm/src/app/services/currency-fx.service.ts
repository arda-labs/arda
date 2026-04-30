import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, map } from 'rxjs';
import { Currency, FxRate, FxRateSource, ListOptions, PageResponse } from '../models/mdm.models';
import { buildParams } from './mdm-http';

@Injectable({ providedIn: 'root' })
export class CurrencyFxService {
  private http = inject(HttpClient);

  listCurrencies(options: ListOptions = {}): Observable<PageResponse<Currency>> {
    return this.http.get<any>('/api/v1/mdm/currencies', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.currencies ?? []).map((item: any) => this.toCurrency(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createCurrency(item: Partial<Currency>): Observable<Currency> {
    return this.http.post<any>('/api/v1/mdm/currencies', { currency: this.fromCurrency(item) }).pipe(map(resp => this.toCurrency(resp)));
  }
  updateCurrency(id: string, item: Partial<Currency>): Observable<Currency> {
    return this.http.put<any>(`/api/v1/mdm/currencies/${encodeURIComponent(id)}`, { currency: this.fromCurrency(item) }).pipe(map(resp => this.toCurrency(resp)));
  }
  deleteCurrency(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/currencies/${encodeURIComponent(id)}`);
  }

  listSources(options: ListOptions = {}): Observable<PageResponse<FxRateSource>> {
    return this.http.get<any>('/api/v1/mdm/fx-rate-sources', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.sources ?? []).map((item: any) => this.toSource(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createSource(item: Partial<FxRateSource>): Observable<FxRateSource> {
    return this.http.post<any>('/api/v1/mdm/fx-rate-sources', { source: this.fromSource(item) }).pipe(map(resp => this.toSource(resp)));
  }
  updateSource(id: string, item: Partial<FxRateSource>): Observable<FxRateSource> {
    return this.http.put<any>(`/api/v1/mdm/fx-rate-sources/${encodeURIComponent(id)}`, { source: this.fromSource(item) }).pipe(map(resp => this.toSource(resp)));
  }
  deleteSource(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/fx-rate-sources/${encodeURIComponent(id)}`);
  }

  listRates(options: ListOptions = {}): Observable<PageResponse<FxRate>> {
    return this.http.get<any>('/api/v1/mdm/fx-rates', { params: buildParams(options) }).pipe(map(resp => ({ items: (resp.fx_rates ?? resp.fxRates ?? []).map((item: any) => this.toRate(item)), nextPageToken: resp.next_page_token ?? resp.nextPageToken ?? '' })));
  }
  createRate(item: Partial<FxRate>): Observable<FxRate> {
    return this.http.post<any>('/api/v1/mdm/fx-rates', { fx_rate: this.fromRate(item) }).pipe(map(resp => this.toRate(resp)));
  }
  updateRate(id: string, item: Partial<FxRate>): Observable<FxRate> {
    return this.http.put<any>(`/api/v1/mdm/fx-rates/${encodeURIComponent(id)}`, { fx_rate: this.fromRate(item) }).pipe(map(resp => this.toRate(resp)));
  }
  deleteRate(id: string): Observable<void> {
    return this.http.delete<void>(`/api/v1/mdm/fx-rates/${encodeURIComponent(id)}`);
  }
  approveRate(id: string): Observable<FxRate> {
    return this.http.post<any>(`/api/v1/mdm/fx-rates/${encodeURIComponent(id)}/approve`, { actor: 'SYSTEM', note: 'Duyệt từ màn hình MDM' }).pipe(map(resp => this.toRate(resp)));
  }

  private toCurrency(item: any): Currency {
    return { id: item.id ?? '', code: item.code ?? '', numericCode: item.numeric_code ?? item.numericCode ?? '', name: item.name ?? '', minorUnit: item.minor_unit ?? item.minorUnit ?? 0, symbol: item.symbol ?? '', countryCode: item.country_code ?? item.countryCode ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromCurrency(item: Partial<Currency>): any {
    return { id: item.id, code: item.code, numeric_code: item.numericCode, name: item.name, minor_unit: item.minorUnit, symbol: item.symbol, country_code: item.countryCode, status: item.status };
  }
  private toSource(item: any): FxRateSource {
    return { id: item.id ?? '', code: item.code ?? '', name: item.name ?? '', sourceType: item.source_type ?? item.sourceType ?? 'MANUAL', priority: item.priority ?? 100, timezone: item.timezone ?? 'Asia/Ho_Chi_Minh', description: item.description ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromSource(item: Partial<FxRateSource>): any {
    return { id: item.id, code: item.code, name: item.name, source_type: item.sourceType, priority: item.priority, timezone: item.timezone, description: item.description, status: item.status };
  }
  private toRate(item: any): FxRate {
    return { id: item.id ?? '', baseCurrency: item.base_currency ?? item.baseCurrency ?? '', quoteCurrency: item.quote_currency ?? item.quoteCurrency ?? '', sourceCode: item.source_code ?? item.sourceCode ?? '', rateDate: item.rate_date ?? item.rateDate ?? '', buyRate: item.buy_rate ?? item.buyRate ?? 0, sellRate: item.sell_rate ?? item.sellRate ?? 0, midRate: item.mid_rate ?? item.midRate ?? 0, approvalStatus: item.approval_status ?? item.approvalStatus ?? 'DRAFT', version: item.version ?? 1, approvedBy: item.approved_by ?? item.approvedBy ?? '', changeNote: item.change_note ?? item.changeNote ?? '', status: item.status ?? 'ACTIVE' };
  }
  private fromRate(item: Partial<FxRate>): any {
    return { id: item.id, base_currency: item.baseCurrency, quote_currency: item.quoteCurrency, source_code: item.sourceCode, rate_date: item.rateDate, buy_rate: item.buyRate, sell_rate: item.sellRate, mid_rate: item.midRate, approval_status: item.approvalStatus, version: item.version, approved_by: item.approvedBy, change_note: item.changeNote, status: item.status };
  }
}
