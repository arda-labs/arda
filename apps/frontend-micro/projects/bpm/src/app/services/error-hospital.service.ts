import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { FailedTask } from '../models/bpm.models';

/**
 * Stub service for error hospital / dead letter queue.
 *
 * TODO: Wire to backend when endpoints exist.
 * Expected endpoints:
 *   GET  /api/v1/bpm/error-hospital
 *   POST /api/v1/bpm/error-hospital/{id}/retry
 *   POST /api/v1/bpm/error-hospital/{id}/save-and-retry
 */
@Injectable({ providedIn: 'root' })
export class ErrorHospitalService {
  list(): Observable<FailedTask[]> {
    return of(this.mockTasks);
  }

  retry(id: string): Observable<void> {
    console.log('[ErrorHospitalService] retry', id);
    return of(undefined);
  }

  saveAndRetry(id: string, payload: Record<string, unknown>): Observable<void> {
    console.log('[ErrorHospitalService] saveAndRetry', id, payload);
    return of(undefined);
  }

  private readonly mockTasks: FailedTask[] = [
    {
      id: 'E001',
      instanceId: 'INST001',
      processName: 'Đăng ký KH mới',
      stepName: 'Xác thực thông tin',
      error: 'Timeout khi gọi API CIC',
      payload: JSON.stringify({ customerId: 'KH001', idNumber: '123456789', retryCount: 2 }, null, 2),
      retryCount: 2,
      lastAttempt: '2026-05-03T09:15:00',
      status: 'PENDING',
    },
    {
      id: 'E002',
      instanceId: 'INST005',
      processName: 'Giải ngân khoản vay',
      stepName: 'Chuyển tiền',
      error: 'Số dư tài khoản không đủ',
      payload: JSON.stringify({ loanId: 'LN001', amount: 500000000, account: '1234567890' }, null, 2),
      retryCount: 0,
      lastAttempt: '2026-05-03T10:00:00',
      status: 'PENDING',
    },
  ];
}
