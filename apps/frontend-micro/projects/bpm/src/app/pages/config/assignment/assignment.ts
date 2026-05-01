import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SelectModule } from 'primeng/select';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { ToggleSwitchModule } from 'primeng/toggleswitch';
import { CardModule } from 'primeng/card';
import { InputNumberModule } from 'primeng/inputnumber';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';

@Component({
  selector: 'app-assignment',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    SelectModule,
    ButtonModule,
    TableModule,
    ToggleSwitchModule,
    CardModule,
    InputNumberModule,
    TagModule,
    ToastModule
  ],
  providers: [MessageService],
  templateUrl: './assignment.html',
  styleUrl: './assignment.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class AssignmentComponent {
  strategies = [
    { label: 'Chia đều (Round Robin)', value: 'ROUND_ROBIN' },
    { label: 'Cân bằng tải (Load Balance)', value: 'LOAD_BALANCE' },
    { label: 'Ưu tiên trọng số (Weighted)', value: 'WEIGHTED' },
    { label: 'Thủ công (Manual)', value: 'MANUAL' }
  ];

  modules = [
    { label: 'Quản lý khách hàng (CRM)', value: 'CRM' },
    { label: 'Thẩm định khoản vay (LOAN)', value: 'LOAN' },
    { label: 'Quản trị nhân sự (HRM)', value: 'HRM' }
  ];

  selectedModule = signal('CRM');
  selectedStrategy = signal('ROUND_ROBIN');

  rules = signal([
    { id: 'A01', name: 'Nhân viên A', department: 'CRM', weight: 1, active: true, taskCount: 5 },
    { id: 'A02', name: 'Nhân viên B', department: 'CRM', weight: 1, active: true, taskCount: 3 },
    { id: 'A03', name: 'Nhân viên C', department: 'CRM', weight: 2, active: false, taskCount: 0 }
  ]);

  constructor(private messageService: MessageService) {}

  onSave() {
    this.messageService.add({
      severity: 'success',
      summary: 'Thành công',
      detail: `Đã cập nhật cấu hình chia bài cho module ${this.selectedModule()}`
    });
  }

  getLoadSeverity(count: number) {
    if (count > 10) return 'danger';
    if (count > 5) return 'warn';
    return 'success';
  }
}
