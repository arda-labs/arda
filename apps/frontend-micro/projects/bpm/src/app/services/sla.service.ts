import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { SlaConfig } from '../models/bpm.models';

/**
 * Stub service for SLA configuration.
 *
 * TODO: Wire to backend when SLA config endpoints exist.
 * Expected endpoints:
 *   GET  /api/v1/bpm/sla-configs
 *   PUT  /api/v1/bpm/sla-configs/{id}
 *   POST /api/v1/bpm/sla-configs
 */
@Injectable({ providedIn: 'root' })
export class SlaService {
  list(): Observable<SlaConfig[]> {
    return of(this.mockConfigs);
  }

  update(configs: SlaConfig[]): Observable<void> {
    console.log('[SlaService] update', configs);
    return of(undefined);
  }

  private readonly mockConfigs: SlaConfig[] = [
    { id: 'S001', processName: 'Đăng ký KH mới', stepName: 'Tiếp nhận hồ sơ', durationHours: 24, unit: 'HOURS', warningPercent: 80 },
    { id: 'S002', processName: 'Đăng ký KH mới', stepName: 'Thẩm định', durationHours: 72, unit: 'HOURS', warningPercent: 75 },
    { id: 'S003', processName: 'Giải ngân khoản vay', stepName: 'Phê duyệt', durationHours: 48, unit: 'HOURS', warningPercent: 85 },
    { id: 'S004', processName: 'Thay đổi hạn mức', stepName: 'Xét duyệt', durationHours: 5, unit: 'DAYS', warningPercent: 70 },
  ];
}
