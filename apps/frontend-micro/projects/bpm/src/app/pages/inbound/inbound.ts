import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { DialogModule } from 'primeng/dialog';
import { TimelineModule } from 'primeng/timeline';
import { TabsModule } from 'primeng/tabs';
import { TextareaModule } from 'primeng/textarea';
import { FormsModule } from '@angular/forms';
import { SelectModule } from 'primeng/select';
import { TooltipModule } from 'primeng/tooltip';

@Component({
  selector: 'app-inbound',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TableModule,
    ButtonModule,
    TagModule,
    ToastModule,
    DialogModule,
    TimelineModule,
    TabsModule,
    TextareaModule,
    SelectModule,
    TooltipModule
  ],
  providers: [MessageService],
  templateUrl: './inbound.html',
  styleUrl: './inbound.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class InboundComponent {
  tasks = signal([
    { id: 'T001', title: 'Duyệt hồ sơ CRM - Nguyễn Văn A', module: 'CRM', priority: 'HIGH', slaStatus: 'ON_TIME', startTime: '2026-05-01 09:00', customerId: 'C1001', processId: 'wf-991' },
    { id: 'T002', title: 'Thẩm định Khoản vay - Trần Thị B', module: 'LOAN', priority: 'MEDIUM', slaStatus: 'WARNING', startTime: '2026-05-01 10:30', customerId: 'C1002', processId: 'wf-992' },
    { id: 'T003', title: 'Duyệt hồ sơ HRM - Lê Văn C', module: 'HRM', priority: 'LOW', slaStatus: 'OVERDUE', startTime: '2026-04-30 14:00', customerId: 'C1003', processId: 'wf-993' },
    { id: 'T004', title: 'Duyệt hồ sơ CRM - Phạm Văn D', module: 'CRM', priority: 'MEDIUM', slaStatus: 'ON_TIME', startTime: '2026-05-02 08:00', customerId: 'C1004', processId: 'wf-994' }
  ]);

  selectedTasks = signal<any[]>([]);
  displayDetail = signal(false);
  activeTask = signal<any>(null);

  reasons = [
    { label: 'Hồ sơ đầy đủ điều kiện', value: 'OK' },
    { label: 'Thiếu CCCD/Hộ chiếu', value: 'MISSING_ID' },
    { label: 'Chứng minh thu nhập chưa rõ ràng', value: 'INCOME_UNCLEAR' },
    { label: 'Nghi ngờ gian lận', value: 'FRAUD_SUSPECT' }
  ];

  selectedReason = signal('');
  comment = signal('');

  taskHistory = [
    { status: 'Khởi tạo', date: '2026-05-01 09:00', user: 'System', icon: 'pi pi-plus', color: '#9C27B0' },
    { status: 'Gán người xử lý', date: '2026-05-01 09:05', user: 'Admin', icon: 'pi pi-user', color: '#673AB7' },
    { status: 'Đang xử lý', date: '2026-05-01 09:10', user: 'Nhân viên A', icon: 'pi pi-cog', color: '#FF9800' }
  ];

  constructor(private messageService: MessageService) {}

  onBulkApprove() {
    const count = this.selectedTasks().length;
    this.messageService.add({
      severity: 'info',
      summary: 'Đang xử lý hàng loạt',
      detail: `Đã gửi yêu cầu phê duyệt cho ${count} hồ sơ qua Kafka.`
    });

    setTimeout(() => {
      this.tasks.update(ts => ts.filter(t => !this.selectedTasks().includes(t)));
      this.selectedTasks.set([]);
      this.messageService.add({
        severity: 'success',
        summary: 'Thành công',
        detail: 'Tất cả hồ sơ đã được đưa vào luồng phê duyệt.'
      });
    }, 1500);
  }

  showTaskDetail(task: any) {
    this.activeTask.set(task);
    this.displayDetail.set(true);
  }

  onAction(action: string) {
    const task = this.activeTask();
    this.messageService.add({
      severity: action === 'APPROVE' ? 'success' : 'danger',
      summary: action === 'APPROVE' ? 'Đã phê duyệt' : 'Đã từ chối',
      detail: `Giao dịch ${task.id} đã được xử lý.`
    });

    this.tasks.update(ts => ts.filter(t => t.id !== task.id));
    this.displayDetail.set(false);
  }

  getSeverity(status: string) {
    switch (status) {
      case 'ON_TIME': return 'success';
      case 'WARNING': return 'warn';
      case 'OVERDUE': return 'danger';
      default: return 'info';
    }
  }

  getPrioritySeverity(priority: string) {
    switch (priority) {
      case 'HIGH': return 'danger';
      case 'MEDIUM': return 'warn';
      case 'LOW': return 'info';
      default: return 'secondary';
    }
  }
}
