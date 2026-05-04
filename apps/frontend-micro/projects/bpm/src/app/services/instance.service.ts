import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { InstanceSummary, InstanceDetail, ProcessEvent, HeatmapStep, PageResponse } from '../models/bpm.models';
import { buildParams, BpmListOptions, extractPageResponse, snakeToCamel } from './bpm-http';

const API = '/api/v1/bpm';

@Injectable({ providedIn: 'root' })
export class InstanceService {
  private http = inject(HttpClient);

  list(options?: BpmListOptions): Observable<PageResponse<InstanceSummary>> {
    return this.http
      .get<{ instances: InstanceSummary[]; next_page_token: string }>(`${API}/instances`, { params: buildParams(options) })
      .pipe(map(extractPageResponse<InstanceSummary>('instances')));
  }

  getById(id: string): Observable<InstanceDetail> {
    return this.http.get<any>(`${API}/instances/${id}`).pipe(map(r => snakeToCamel<InstanceDetail>(r)));
  }

  getEvents(id: string, pageSize?: number, pageToken?: string): Observable<PageResponse<ProcessEvent>> {
    const params = buildParams({ pageSize, pageToken });
    return this.http
      .get<{ events: ProcessEvent[]; next_page_token: string }>(`${API}/instances/${id}/events`, { params })
      .pipe(map(extractPageResponse<ProcessEvent>('events')));
  }

  getHeatmap(processDefinitionId: string): Observable<HeatmapStep[]> {
    return this.http
      .get<{ steps: HeatmapStep[] }>(`${API}/heatmap`, { params: buildParams({ processDefinitionId }) })
      .pipe(map(r => r.steps ?? []));
  }
}
