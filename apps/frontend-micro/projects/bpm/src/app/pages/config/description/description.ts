import { Component, ChangeDetectionStrategy, signal, inject, computed, resource } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { DefinitionService } from '../../../services/definition.service';
import { TemplateService } from '../../../services/template.service';
import { InstanceService } from '../../../services/instance.service';
import { ProcessDefinition, Template, TemplateVariable, VariableSource, InstanceSummary } from '../../../models/bpm.models';

@Component({
  selector: 'app-description',
  standalone: true,
  imports: [CommonModule, FormsModule, TableModule, ButtonModule, InputTextModule, ToastModule, SelectModule, TagModule, DialogModule],
  providers: [MessageService],
  templateUrl: './description.html',
  styleUrl: './description.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DescriptionComponent {
  private definitionService = inject(DefinitionService);
  private templateService = inject(TemplateService);
  private instanceService = inject(InstanceService);
  private messageService = inject(MessageService);

  definitions = signal<ProcessDefinition[]>([]);
  templates = signal<Template[]>([]);
  variableSources = signal<VariableSource[]>([]);
  activeTab = signal<'editor' | 'preview'>('editor');

  selectedDefId = signal<string>('');
  loading = signal(false);

  // Template editor state
  editingTemplate = signal<Template | null>(null);
  showEditor = signal(false);

  // Variable mapping dialog
  showVarDialog = signal(false);
  editingVar = signal<TemplateVariable | null>(null);
  editingVarIndex = signal(-1);

  // Preview state
  previewTemplateId = signal<string>('');
  previewInstanceId = signal<string>('');
  instances = signal<InstanceSummary[]>([]);
  previewResult = signal<{ renderedText: string; resolvedVariables: any[] } | null>(null);
  loadingPreview = signal(false);

  editorTemplates = signal<Template[]>([]);

  constructor() {
    this.definitionService.list({ pageSize: 100 }).subscribe(r => this.definitions.set(r?.items ?? []));
    this.templateService.listVariableSources().subscribe(r => this.variableSources.set(r ?? []));
    this.loadTemplates();
  }

  async loadTemplates() {
    this.loading.set(true);
    try {
      const r = await this.templateService.list({ pageSize: 100 }).toPromise();
      this.templates.set(r?.items ?? []);
      this.editorTemplates.set(r?.items ?? []);
    } finally {
      this.loading.set(false);
    }
  }

  /* ───── Template CRUD ───── */

  newTemplate() {
    this.editingTemplate.set({
      id: '',
      processDefinitionId: '',
      name: '',
      templateText: '',
      module: '',
      variables: [],
      createdAt: '',
      updatedAt: '',
    });
    this.showEditor.set(true);
  }

  editTemplate(tpl: Template) {
    this.editingTemplate.set({ ...tpl, variables: [...(tpl.variables ?? [])] });
    this.showEditor.set(true);
  }

  async saveTemplate() {
    const tpl = this.editingTemplate();
    if (!tpl || !tpl.name || !tpl.templateText) {
      this.messageService.add({ severity: 'warn', summary: 'Thiếu dữ liệu', detail: 'Vui lòng nhập tên và nội dung mẫu.' });
      return;
    }
    try {
      if (tpl.id) {
        await this.templateService.update(tpl.id, { name: tpl.name, templateText: tpl.templateText, module: tpl.module, variables: tpl.variables }).toPromise();
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Mẫu đã được cập nhật.' });
      } else {
        if (!tpl.processDefinitionId) {
          this.messageService.add({ severity: 'warn', summary: 'Thiếu dữ liệu', detail: 'Vui lòng chọn quy trình.' });
          return;
        }
        await this.templateService.create({
          processDefinitionId: tpl.processDefinitionId,
          name: tpl.name,
          templateText: tpl.templateText,
          module: tpl.module,
          variables: tpl.variables,
        }).toPromise();
        this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Mẫu mới đã được tạo.' });
      }
      this.showEditor.set(false);
      await this.loadTemplates();
    } catch (e: any) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: e?.message ?? 'Không thể lưu mẫu.' });
    }
  }

  async deleteTemplate(id: string) {
    try {
      await this.templateService.delete(id).toPromise();
      this.messageService.add({ severity: 'success', summary: 'Thành công', detail: 'Mẫu đã được xóa.' });
      await this.loadTemplates();
    } catch (e: any) {
      this.messageService.add({ severity: 'error', summary: 'Lỗi', detail: e?.message ?? 'Không thể xóa mẫu.' });
    }
  }

  /* ───── Variable mapping ───── */

  addVariable(tpl: Template) {
    this.editingVar.set({ variableName: '', sourceType: 'INSTANCE_VAR', sourceField: '', fallbackValue: '' });
    this.editingVarIndex.set(-1);
    this.showVarDialog.set(true);
  }

  editVariable(tpl: Template, idx: number) {
    this.editingVar.set({ ...tpl.variables[idx] });
    this.editingVarIndex.set(idx);
    this.showVarDialog.set(true);
  }

  saveVariable() {
    const tpl = this.editingTemplate();
    const v = this.editingVar();
    if (!tpl || !v || !v.variableName) return;
    if (this.editingVarIndex() >= 0) {
      tpl.variables[this.editingVarIndex()] = v;
    } else {
      tpl.variables = [...tpl.variables, v];
    }
    this.editingTemplate.set({ ...tpl });
    this.showVarDialog.set(false);
  }

  removeVariable(tpl: Template, idx: number) {
    tpl.variables = tpl.variables.filter((_, i) => i !== idx);
    this.editingTemplate.set({ ...tpl });
  }

  getSourceName(type: string): string {
    return this.variableSources().find(s => s.type === type)?.name ?? type;
  }

  /* ───── Preview ───── */

  async onSelectPreviewTemplate(id: string) {
    this.previewTemplateId.set(id);
    this.previewResult.set(null);
    if (!id) return;
    // Load instances for the template's definition
    const tpl = this.templates().find(t => t.id === id);
    if (tpl) {
      const r = await this.instanceService.list({ processDefinitionId: tpl.processDefinitionId, pageSize: 50 }).toPromise();
      this.instances.set(r?.items ?? []);
    }
  }

  async onPreview() {
    if (!this.previewTemplateId() || !this.previewInstanceId()) {
      this.messageService.add({ severity: 'warn', summary: 'Thiếu dữ liệu', detail: 'Vui lòng chọn mẫu và phiên bản xử lý.' });
      return;
    }
    this.loadingPreview.set(true);
    try {
      const r = await this.templateService.render(this.previewTemplateId(), this.previewInstanceId()).toPromise();
      this.previewResult.set(r ?? null);
    } finally {
      this.loadingPreview.set(false);
    }
  }
}
