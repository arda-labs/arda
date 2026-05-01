import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { TextareaModule } from 'primeng/textarea';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { InputTextModule } from 'primeng/inputtext';

@Component({
  selector: 'app-error-hospital',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TableModule,
    ButtonModule,
    TagModule,
    DialogModule,
    TextareaModule,
    ToastModule,
    InputTextModule
  ],
  providers: [MessageService],
  templateUrl: './error-hospital.html',
  styleUrl: './error-hospital.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ErrorHospitalComponent {
  failedTasks = signal([
    { id: 'ERR-001', process: 'CRM Register', step: 'Finalize API', error: 'Connect Timeout', payload: '{\n  "customer_id": "C001",\n  "status": "APPROVED",\n  "retry_count": 0\n}', date: '2026-05-02 08:30', retryCount: 2 },
    { id: 'ERR-002', process: 'Loan Approval', step: 'Risk Scoring', error: 'Service Unavailable (503)', payload: '{\n  "app_id": "L-99",\n  "score_model": "V2"\n}', date: '2026-05-02 09:15', retryCount: 0 }
  ]);

  displayEdit = signal(false);
  selectedTask = signal<any>(null);

  constructor(private messageService: MessageService) {}

  onRetry(task: any) {
    this.messageService.add({
      severity: 'info',
      summary: 'Đang xử lý',
      detail: `Đang gửi yêu cầu retry cho task ${task.id}...`
    });

    setTimeout(() => {
      this.messageService.add({
        severity: 'success',
        summary: 'Thành công',
        detail: `Task ${task.id} đã được đưa lại vào hàng đợi xử lý.`
      });
    }, 1000);
  }

  showEdit(task: any) {
    this.selectedTask.set({...task});
    this.displayEdit.set(true);
  }

  saveAndRetry() {
    const task = this.selectedTask();
    this.messageService.add({
      severity: 'success',
      summary: 'Đã lưu & Retry',
      detail: `Payload mới đã được cập nhật cho task ${task.id}.`
    });
    this.displayEdit.set(false);
  }

  getRetrySeverity(count: number) {
    if (count > 5) return 'danger';
    if (count > 0) return 'warn';
    return 'info';
  }
}
