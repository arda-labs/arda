import { Component, ChangeDetectionStrategy, signal, inject, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { DatePickerModule } from 'primeng/datepicker';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { Router } from '@angular/router';
import { DefinitionService } from '../../services/definition.service';
import { InstanceService } from '../../services/instance.service';
import { ProcessDefinition, InstanceSummary } from '../../models/bpm.models';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule, InputTextModule, ButtonModule, TableModule, DatePickerModule, SelectModule, TagModule, IconFieldModule, InputIconModule],
  templateUrl: './search.html',
  styleUrl: './search.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchComponent {
  private definitionService = inject(DefinitionService);
  private instanceService = inject(InstanceService);
  private router = inject(Router);

  definitions = signal<ProcessDefinition[]>([]);
  keyword = signal('');
  selectedDefId = signal<string>('');
  selectedStatus = signal<string>('');
  dateRange = signal<Date[] | null>(null);
  results = signal<InstanceSummary[]>([]);
  loading = signal(false);

  statuses = [
    { label: 'Tất cả trạng thái', value: '' },
    { label: 'Đang xử lý', value: 'ACTIVE' },
    { label: 'Hoàn thành', value: 'COMPLETED' },
    { label: 'Lỗi', value: 'FAILED' },
    { label: 'Đã hủy', value: 'CANCELLED' },
  ];

  constructor() {
    this.definitionService.list({ pageSize: 100 }).subscribe(r => this.definitions.set(r?.items ?? []));
  }

  getStatusSeverity(status: string) {
    switch (status) {
      case 'COMPLETED': return 'success';
      case 'ACTIVE': return 'info';
      case 'FAILED': return 'danger';
      case 'CANCELLED': return 'secondary';
      default: return 'warn';
    }
  }

  onSearch() {
    this.loading.set(true);
    this.instanceService.list({
      keyword: this.keyword() || undefined,
      processDefinitionId: this.selectedDefId() || undefined,
      status: this.selectedStatus() || undefined,
      pageSize: 50,
    }).subscribe({
      next: (r) => this.results.set(r?.items ?? []),
      error: () => this.results.set([]),
      complete: () => this.loading.set(false),
    });
  }

  viewInstance(id: string) {
    this.router.navigate(['/bpm/monitor'], { queryParams: { instanceId: id } });
  }
}
