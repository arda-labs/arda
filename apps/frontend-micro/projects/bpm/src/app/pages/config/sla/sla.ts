import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputNumberModule } from 'primeng/inputnumber';
import { FormsModule } from '@angular/forms';
import { ToastModule } from 'primeng/toast';
import { TagModule } from 'primeng/tag';
import { MessageService } from 'primeng/api';
import { SelectModule } from 'primeng/select';

@Component({
  selector: 'app-sla',
  standalone: true,
  imports: [CommonModule, TableModule, ButtonModule, InputNumberModule, FormsModule, ToastModule, TagModule, SelectModule],
  providers: [MessageService],
  templateUrl: './sla.html',
  styleUrl: './sla.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SlaComponent {
  units = [
    { label: 'Phút', value: 'MINUTES' },
    { label: 'Giờ', value: 'HOURS' },
    { label: 'Ngày', value: 'DAYS' }
  ];

  slaConfigs = signal([
    { process: 'CRM Register', step: 'Khởi tạo hồ sơ', duration: 2, unit: 'HOURS', warningPercent: 80 },
    { process: 'CRM Register', step: 'Phê duyệt quản lý', duration: 4, unit: 'HOURS', warningPercent: 75 },
    { process: 'Loan App', step: 'Thẩm định rủi ro', duration: 1, unit: 'DAYS', warningPercent: 90 },
    { process: 'HRM Onboarding', step: 'Cấp phát thiết bị', duration: 4, unit: 'HOURS', warningPercent: 50 }
  ]);

  constructor(private messageService: MessageService) {}

  onSave() {
    this.messageService.add({
      severity: 'success',
      summary: 'Thành công',
      detail: 'Cấu hình SLA đã được lưu và áp dụng cho các instance mới.'
    });
  }
}
