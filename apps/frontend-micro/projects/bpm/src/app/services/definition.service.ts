import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ProcessDefinition, DefinitionDiagram } from '../models/bpm.models';
import { DeployDefinitionRequest } from '../models/bpm.models';
import { buildParams, BpmListOptions, extractPageResponse, snakeToCamel } from './bpm-http';
import { PageResponse } from '../models/bpm.models';

const API = '/api/v1/bpm';

@Injectable({ providedIn: 'root' })
export class DefinitionService {
  private http = inject(HttpClient);

  list(options?: BpmListOptions): Observable<PageResponse<ProcessDefinition>> {
    return this.http
      .get<{ definitions: ProcessDefinition[]; next_page_token: string }>(`${API}/definitions`, { params: buildParams(options) })
      .pipe(map(extractPageResponse<ProcessDefinition>('definitions')));
  }

  getById(id: string): Observable<ProcessDefinition> {
    return this.http.get<ProcessDefinition>(`${API}/definitions/${id}`);
  }

  getDiagram(id: string): Observable<DefinitionDiagram> {
    return this.http
      .get<any>(`${API}/definitions/${id}/diagram`)
      .pipe(map(r => snakeToCamel<DefinitionDiagram>(r)));
  }

  deploy(req: DeployDefinitionRequest): Observable<ProcessDefinition> {
    return this.http.post<ProcessDefinition>(`${API}/definitions/deploy`, {
      bpmn_xml: req.bpmnXml,
      process_key: req.processKey,
      name: req.name,
      description: req.description || '',
      category: req.category || '',
      module: req.module || '',
    });
  }
}
