import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Template, TemplateVariable, VariableSource, RenderResult, ResolvedVariable, PageResponse } from '../models/bpm.models';
import { buildParams, BpmListOptions, extractPageResponse } from './bpm-http';

const API = '/api/v1/bpm';

@Injectable({ providedIn: 'root' })
export class TemplateService {
  private http = inject(HttpClient);

  list(options?: BpmListOptions): Observable<PageResponse<Template>> {
    return this.http
      .get<{ templates: Template[]; next_page_token: string }>(`${API}/templates`, { params: buildParams(options) })
      .pipe(map(extractPageResponse<Template>('templates')));
  }

  getById(id: string): Observable<Template> {
    return this.http.get<Template>(`${API}/templates/${id}`);
  }

  create(data: {
    processDefinitionId: string;
    name: string;
    templateText: string;
    module: string;
    variables: TemplateVariable[];
  }): Observable<Template> {
    return this.http.post<Template>(`${API}/templates`, {
      process_definition_id: data.processDefinitionId,
      name: data.name,
      template_text: data.templateText,
      module: data.module,
      variables: data.variables,
    });
  }

  update(id: string, data: { name: string; templateText: string; module: string; variables: TemplateVariable[] }): Observable<Template> {
    return this.http.put<Template>(`${API}/templates/${id}`, {
      name: data.name,
      template_text: data.templateText,
      module: data.module,
      variables: data.variables,
    });
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${API}/templates/${id}`);
  }

  render(id: string, instanceId: string): Observable<RenderResult> {
    return this.http
      .get<{ rendered_text: string; resolved_variables: ResolvedVariable[] }>(`${API}/templates/${id}/render`, {
        params: { instance_id: instanceId },
      })
      .pipe(map(r => ({
        renderedText: r.rendered_text ?? '',
        resolvedVariables: r.resolved_variables ?? [],
      })));
  }

  listVariableSources(): Observable<VariableSource[]> {
    return this.http
      .get<{ sources: VariableSource[] }>(`${API}/templates/variables/sources`)
      .pipe(map(r => r.sources ?? []));
  }
}
