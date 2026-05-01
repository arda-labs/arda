import { Component, ChangeDetectionStrategy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { DatePickerModule } from 'primeng/datepicker';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    InputTextModule,
    ButtonModule,
    TableModule,
    DatePickerModule,
    SelectModule,
    TagModule,
    IconFieldModule,
    InputIconModule
  ],
  templateUrl: './search.html',
  styleUrl: './search.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class SearchComponent {
  searchQuery = signal('');
  dateRange = signal(null);

  processTypes = [
    { label: 'Tất cả quy trình', value: 'ALL' },
    { label: 'Đăng ký khách hàng (CRM)', value: 'CRM' },
    { label: 'Thẩm định khoản vay (LOAN)', value: 'LOAN' },
    { label: 'Quản trị nhân sự (HRM)', value: 'HRM' }
  ];

  selectedProcess = signal('ALL');

  results = signal([
    { id: 'P-00123', process: 'CRM Register', customer: 'Nguyễn Văn A', status: 'COMPLETED', completedDate: '2026-05-01 14:20' },
    { id: 'P-00124', process: 'Loan App', customer: 'Trần Thị B', status: 'RUNNING', completedDate: '-' },
    { id: 'P-00125', process: 'CRM Register', customer: 'Lê Văn C', status: 'FAILED', completedDate: '2026-05-02 08:30' },
    { id: 'P-00126', process: 'HRM Onboarding', customer: 'Phạm Văn D', status: 'COMPLETED', completedDate: '2026-05-02 10:15' }
  ]);

  getStatusSeverity(status: string) {
    switch (status) {
      case 'COMPLETED': return 'success';
      case 'RUNNING': return 'info';
      case 'FAILED': return 'danger';
      case 'CANCELLED': return 'secondary';
      default: return 'warn';
    }
  }

  onSearch() {
    console.log('Searching for:', this.searchQuery(), this.selectedProcess(), this.dateRange());
  }

  onExport() {
    console.log('Exporting results to Excel/CSV...');
  }
}
