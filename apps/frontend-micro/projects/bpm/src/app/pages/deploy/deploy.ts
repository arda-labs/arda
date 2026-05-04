import { Component, ChangeDetectionStrategy, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { CardModule } from 'primeng/card';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { TextareaModule } from 'primeng/textarea';
import { SelectModule } from 'primeng/select';
import { MessageModule } from 'primeng/message';
import { DefinitionService } from '../../services/definition.service';

@Component({
  selector: 'app-deploy',
  standalone: true,
  imports: [CommonModule, FormsModule, CardModule, ButtonModule, InputTextModule, TextareaModule, SelectModule, MessageModule],
  templateUrl: './deploy.html',
  styleUrl: './deploy.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DeployComponent {
  private definitionService = inject(DefinitionService);
  private router = inject(Router);

  readonly selectedFile = signal<File | null>(null);
  readonly bpmnXml = signal('');
  readonly processKey = signal('');
  readonly name = signal('');
  readonly description = signal('');
  readonly category = signal('');
  readonly selectedModule = signal<string>('');
  readonly isDeploying = signal(false);
  readonly deployError = signal<string | null>(null);

  readonly moduleOptions = [
    { label: '-- Chọn module --', value: '' },
    { label: 'CRM', value: 'crm' },
    { label: 'BPM', value: 'bpm' },
    { label: 'Loan', value: 'loan' },
    { label: 'HRM', value: 'hrm' },
  ];

  readonly canDeploy = () => !!this.selectedFile() && !!this.processKey() && !!this.name() && !this.isDeploying();

  onFileSelected(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files?.length) {
      this.readFile(input.files[0]);
    }
  }

  onFileDrop(event: DragEvent) {
    event.preventDefault();
    const file = event.dataTransfer?.files?.[0];
    if (file) {
      this.readFile(file);
    }
  }

  private readFile(file: File) {
    if (!file.name.endsWith('.bpmn') && !file.name.endsWith('.xml')) {
      this.deployError.set('Chỉ hỗ trợ file .bpmn hoặc .xml');
      return;
    }
    this.deployError.set(null);
    this.selectedFile.set(file);
    const reader = new FileReader();
    reader.onload = () => {
      this.bpmnXml.set(reader.result as string);
    };
    reader.readAsText(file);
  }

  deploy() {
    if (!this.canDeploy()) return;

    this.isDeploying.set(true);
    this.deployError.set(null);

    this.definitionService.deploy({
      bpmnXml: this.bpmnXml(),
      processKey: this.processKey(),
      name: this.name(),
      description: this.description() || undefined,
      category: this.category() || undefined,
      module: this.selectedModule() || undefined,
    }).subscribe({
      next: () => {
        this.isDeploying.set(false);
        this.router.navigate(['/bpm/definitions']);
      },
      error: (err) => {
        this.isDeploying.set(false);
        this.deployError.set(err?.error?.message || err?.message || 'Triển khai thất bại');
      },
    });
  }

  goBack() {
    this.router.navigate(['/bpm/definitions']);
  }
}
