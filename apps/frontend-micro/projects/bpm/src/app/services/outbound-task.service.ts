import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { BpmTask } from '../models/bpm.models';

/**
 * Stub service for outbound (sent/pushed) tasks.
 *
 * TODO: Replace with Zeebe tasklist API when available.
 * Expected endpoints:
 *   GET  /api/v1/bpm/tasks?state=CREATED&assignedTo=external
 *   POST /api/v1/bpm/tasks/{id}/recall
 */
@Injectable({ providedIn: 'root' })
export class OutboundTaskService {
  list(): Observable<BpmTask[]> {
    return of(this.mockTasks);
  }

  recall(id: string): Observable<void> {
    console.log('[OutboundTaskService] recall', id);
    return of(undefined);
  }

  stats(): Observable<{ inProgress: number; completed: number; returned: number }> {
    return of({ inProgress: 12, completed: 450, returned: 3 });
  }

  private readonly mockTasks: BpmTask[] = [
    { id: 'O001', title: 'Gửi thông báo CQ Thuế - KH A', module: 'CRM', priority: 'HIGH', slaStatus: 'ON_TRACK', startTime: '2026-05-03T07:00:00', customerId: 'KH001', processId: 'PID001', status: 'COMPLETED', assignedTo: 'Tax Authority' },
    { id: 'O002', title: 'Yêu cầu bổ sung hồ sơ - KH B', module: 'LOAN', priority: 'MEDIUM', slaStatus: 'WARNING', startTime: '2026-05-02T15:30:00', customerId: 'KH002', processId: 'PID002', status: 'IN_PROGRESS', assignedTo: 'KH B' },
    { id: 'O003', title: 'Gửin hợp đồng cho KH C', module: 'CRM', priority: 'HIGH', slaStatus: 'ON_TRACK', startTime: '2026-05-03T10:00:00', customerId: 'KH003', processId: 'PID003', status: 'IN_PROGRESS', assignedTo: 'KH C' },
    { id: 'O004', title: 'Thông báo kết quả thẩm định - KH D', module: 'LOAN', priority: 'LOW', slaStatus: 'BREACHED', startTime: '2026-04-30T08:00:00', customerId: 'KH004', processId: 'PID004', status: 'RETURNED', assignedTo: 'KH D' },
  ];
}
