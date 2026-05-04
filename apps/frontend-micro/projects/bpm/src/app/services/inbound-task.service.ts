import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { BpmTask } from '../models/bpm.models';

/**
 * Stub service for inbound (user task inbox).
 *
 * TODO: Replace with Zeebe tasklist API when available.
 * Expected endpoints:
 *   GET  /api/v1/bpm/tasks?state=CREATED&assignedTo=me
 *   POST /api/v1/bpm/tasks/{id}/claim
 *   POST /api/v1/bpm/tasks/{id}/complete
 *   POST /api/v1/bpm/tasks/{id}/fail
 */
@Injectable({ providedIn: 'root' })
export class InboundTaskService {
  list(): Observable<BpmTask[]> {
    return of(this.mockTasks);
  }

  approve(id: string): Observable<void> {
    console.log('[InboundTaskService] approve', id);
    return of(undefined);
  }

  reject(id: string): Observable<void> {
    console.log('[InboundTaskService] reject', id);
    return of(undefined);
  }

  private readonly mockTasks: BpmTask[] = [
    { id: 'T001', title: 'Phê duyệt đăng ký KH mới - Công ty A', module: 'CRM', priority: 'HIGH', slaStatus: 'ON_TRACK', startTime: '2026-05-03T08:30:00', customerId: 'KH001', processId: 'PID001', status: 'IN_PROGRESS' },
    { id: 'T002', title: 'Thẩm định hồ sơ vay - Nguyễn Văn B', module: 'LOAN', priority: 'HIGH', slaStatus: 'WARNING', startTime: '2026-05-02T14:00:00', customerId: 'KH002', processId: 'PID002', status: 'IN_PROGRESS' },
    { id: 'T003', title: 'Xác nhận thông tin nhân sự - Trần Thị C', module: 'HRM', priority: 'MEDIUM', slaStatus: 'ON_TRACK', startTime: '2026-05-03T09:15:00', customerId: 'KH003', processId: 'PID003', status: 'IN_PROGRESS' },
    { id: 'T004', title: 'Duyệt báo giá - KH Xuất nhập khẩu D', module: 'CRM', priority: 'LOW', slaStatus: 'ON_TRACK', startTime: '2026-05-01T10:00:00', customerId: 'KH004', processId: 'PID004', status: 'IN_PROGRESS' },
  ];
}
