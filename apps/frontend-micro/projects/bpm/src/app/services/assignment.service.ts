import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { AssignmentRule } from '../models/bpm.models';

/**
 * Stub service for task assignment rules.
 *
 * TODO: Wire to backend when assignment endpoints exist.
 * Expected endpoints:
 *   GET  /api/v1/bpm/assignments?module={module}
 *   PUT  /api/v1/bpm/assignments/{id}
 *   POST /api/v1/bpm/assignments
 */
@Injectable({ providedIn: 'root' })
export class AssignmentService {
  getRule(module: string): Observable<AssignmentRule | null> {
    return of(this.mockRules.find(r => r.module === module) ?? null);
  }

  updateRule(rule: AssignmentRule): Observable<void> {
    console.log('[AssignmentService] updateRule', rule);
    return of(undefined);
  }

  listModules(): Observable<string[]> {
    return of(['CRM', 'LOAN', 'HRM']);
  }

  listStrategies(): Observable<string[]> {
    return of(['ROUND_ROBIN', 'LOAD_BALANCE', 'WEIGHTED', 'MANUAL']);
  }

  private readonly mockRules: AssignmentRule[] = [
    {
      id: 'R001', module: 'CRM', strategy: 'ROUND_ROBIN',
      agents: [
        { id: 'U001', name: 'Nguyễn Văn A', department: 'CSKH', weight: 1, active: true, taskCount: 5 },
        { id: 'U002', name: 'Trần Thị B', department: 'CSKH', weight: 1, active: true, taskCount: 3 },
      ],
    },
    {
      id: 'R002', module: 'LOAN', strategy: 'WEIGHTED',
      agents: [
        { id: 'U003', name: 'Lê Văn C', department: 'Thẩm định', weight: 2, active: true, taskCount: 8 },
      ],
    },
    {
      id: 'R003', module: 'HRM', strategy: 'MANUAL',
      agents: [
        { id: 'U004', name: 'Phạm Thị D', department: 'Nhân sự', weight: 1, active: true, taskCount: 2 },
      ],
    },
  ];
}
