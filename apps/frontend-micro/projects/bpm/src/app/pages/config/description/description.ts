import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';

@Component({
  selector: 'app-description',
  standalone: true,
  imports: [CommonModule, FormsModule, TableModule, ButtonModule, InputTextModule, ToastModule, SelectModule, TagModule],
  providers: [MessageService],
  templateUrl: './description.html',
  styleUrl: './description.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DescriptionComponent {
  modules = [
    { label: 'Dùng chung (General)', value: 'GENERAL' },
    { label: 'Khách hàng (CRM)', value: 'CRM' },
    { label: 'Khoản vay (LOAN)', value: 'LOAN' },
    { label: 'Nhân sự (HRM)', value: 'HRM' }
  ];

  templates = signal([
    { code: 'REJ_01', content: 'Hồ sơ thiếu thông tin định danh (CCCD/Hộ chiếu)', module: 'GENERAL' },
    { code: 'REJ_02', content: 'Hình ảnh đính kèm mờ, không rõ nét, không thể hậu kiểm', module: 'GENERAL' },
    { code: 'REJ_03', content: 'Thu nhập hàng tháng không đủ điều kiện tối thiểu để vay', module: 'LOAN' },
    { code: 'APP_01', content: 'Hồ sơ đầy đủ điều kiện, đồng ý phê duyệt', module: 'GENERAL' },
    { code: 'REQ_01', content: 'Yêu cầu bổ sung sao kê lương 3 tháng gần nhất', module: 'LOAN' }
  ]);

  constructor(private messageService: MessageService) {}

  addTemplate() {
    this.templates.update(prev => [
      { code: 'NEW_CODE', content: '', module: 'GENERAL' },
      ...prev
    ]);
  }

  deleteTemplate(index: number) {
    this.templates.update(prev => prev.filter((_, i) => i !== index));
  }

  onSave() {
    this.messageService.add({
      severity: 'success',
      summary: 'Thành công',
      detail: 'Cấu trúc diễn giải đã được cập nhật hệ thống.'
    });
  }
}
