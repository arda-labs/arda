import { Component, ChangeDetectionStrategy, signal, inject, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { InputTextModule } from 'primeng/inputtext';
import { ButtonModule } from 'primeng/button';
import { TableModule } from 'primeng/table';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { IconFieldModule } from 'primeng/iconfield';
import { InputIconModule } from 'primeng/inputicon';
import { SkeletonModule } from 'primeng/skeleton';
import { TooltipModule } from 'primeng/tooltip';
import { DefinitionService } from '../../services/definition.service';
import { ProcessDefinition } from '../../models/bpm.models';
import { createPagedResource } from '@arda/core';

@Component({
  selector: 'app-definitions',
  standalone: true,
  imports: [CommonModule, FormsModule, InputTextModule, ButtonModule, TableModule, SelectModule, TagModule, IconFieldModule, InputIconModule, SkeletonModule, TooltipModule, DatePipe],
  templateUrl: './definitions.html',
  styleUrl: './definitions.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DefinitionsComponent {
  private definitionService = inject(DefinitionService);
  private router = inject(Router);

  readonly keyword = signal('');
  readonly selectedModule = signal<string>('');
  readonly includeInactive = signal(false);

  readonly moduleOptions = [
    { label: 'Tất cả module', value: '' },
    { label: 'CRM', value: 'crm' },
    { label: 'BPM', value: 'bpm' },
    { label: 'Loan', value: 'loan' },
    { label: 'HRM', value: 'hrm' },
  ];

  readonly refreshTrigger = signal(0);

  readonly table = createPagedResource<ProcessDefinition, number>({
    defaultPageSize: 20,
    rowsPerPageOptions: [10, 20, 50],
    params: () => this.refreshTrigger(),
    load: (_, page) => this.definitionService.list({
      keyword: this.keyword() || undefined,
      module: this.selectedModule() || undefined,
      includeInactive: this.includeInactive() || undefined,
      pageSize: page.pageSize,
      pageToken: page.pageToken,
    }),
  });

  onSearch() {
    this.table.refresh();
  }

  viewDiagram(def: ProcessDefinition) {
    this.router.navigate(['/bpm/monitor'], { queryParams: { definitionId: def.id } });
  }

  goToDeploy() {
    this.router.navigate(['/bpm/deploy']);
  }

  getVersionSeverity(v: number) {
    return v === 1 ? 'info' : 'success';
  }
}
