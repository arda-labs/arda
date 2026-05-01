import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ButtonModule } from 'primeng/button';
import { TooltipModule } from 'primeng/tooltip';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { InputTextModule } from 'primeng/inputtext';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-outbound',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TableModule,
    TagModule,
    ButtonModule,
    TooltipModule,
    IconFieldModule,
    InputIconModule,
    InputTextModule
  ],
  templateUrl: './outbound.html',
  styleUrl: './outbound.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class OutboundComponent {
  searchQuery = signal('');

  tasks = signal([
    { id: 'T101', title: 'Duyệt hồ sơ CRM - Lê Văn C', target: 'Manager Approval', status: 'IN_PROGRESS', startTime: '2026-05-01 15:00', duration: '2h 15m' },
    { id: 'T102', title: 'Thẩm định Loan - Phạm Văn D', target: 'Risk Dept', status: 'COMPLETED', startTime: '2026-04-30 11:00', duration: '4h 30m' },
    { id: 'T103', title: 'Hồ sơ HRM - Nguyễn Thị E', target: 'Compliance', status: 'RETURNED', startTime: '2026-05-02 09:00', duration: '1h 20m' },
    { id: 'T104', title: 'Đăng ký Merchant - Shop ABC', target: 'Finalize API', status: 'FAILED', startTime: '2026-05-02 10:00', duration: '5m' }
  ]);

  getStatusSeverity(status: string) {
    switch (status) {
      case 'IN_PROGRESS': return 'info';
      case 'COMPLETED': return 'success';
      case 'RETURNED': return 'warn';
      case 'FAILED': return 'danger';
      default: return 'secondary';
    }
  }

  onRecall(task: any) {
    console.log('Recalling task:', task.id);
  }
}
