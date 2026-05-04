import { Component, ElementRef, OnInit, ViewChild, ChangeDetectionStrategy, signal, inject, effect, resource, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TimelineModule } from 'primeng/timeline';
import { CardModule } from 'primeng/card';
import { TagModule } from 'primeng/tag';
import { ToggleButtonModule } from 'primeng/togglebutton';
import { SelectModule } from 'primeng/select';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { SkeletonModule } from 'primeng/skeleton';
import BpmnViewer from 'bpmn-js';
import { DefinitionService } from '../../services/definition.service';
import { InstanceService } from '../../services/instance.service';
import { HeatmapStep, InstanceSummary, ProcessDefinition, ProcessEvent } from '../../models/bpm.models';

@Component({
  selector: 'app-monitor',
  standalone: true,
  imports: [CommonModule, FormsModule, TimelineModule, CardModule, TagModule, ToggleButtonModule, SelectModule, TableModule, ButtonModule, SkeletonModule],
  templateUrl: './monitor.html',
  styleUrl: './monitor.css',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MonitorComponent implements OnInit {
  @ViewChild('bpmnContainer', { static: true }) private bpmnEl!: ElementRef;
  private viewer: any;
  private definitionService = inject(DefinitionService);
  private instanceService = inject(InstanceService);
  private route = inject(ActivatedRoute);
  private router = inject(Router);

  definitions = signal<ProcessDefinition[]>([]);
  selectedDef = signal<ProcessDefinition | null>(null);
  instances = signal<InstanceSummary[]>([]);
  selectedInstance = signal<InstanceSummary | null>(null);
  bpmnXml = signal('');
  events = signal<ProcessEvent[]>([]);
  heatmapData = signal<HeatmapStep[]>([]);
  showHeatmap = signal(false);
  loadingDef = signal(true);
  loadingInstances = signal(false);
  loadingDiagram = signal(false);

  constructor() {
    effect(() => {
      if (this.bpmnXml() && this.bpmnEl.nativeElement.children.length === 0) {
        this.renderDiagram();
      }
    });
  }

  async ngOnInit() {
    await this.loadDefinitions();
    const instanceId = this.route.snapshot.queryParamMap.get('instanceId');
    if (instanceId) {
      this.selectInstanceById(instanceId);
    }
  }

  private async loadDefinitions() {
    this.loadingDef.set(true);
    try {
      const res = await this.definitionService.list({ pageSize: 100 }).toPromise();
      this.definitions.set(res?.items ?? []);
      if (res?.items?.length && !this.selectedDef()) {
        this.selectedDef.set(res.items[0]);
        this.loadInstances(res.items[0].id);
      }
    } finally {
      this.loadingDef.set(false);
    }
  }

  onDefChange(def: ProcessDefinition) {
    this.selectedDef.set(def);
    this.selectedInstance.set(null);
    this.bpmnXml.set('');
    this.events.set([]);
    this.heatmapData.set([]);
    this.showHeatmap.set(false);
    if (def) this.loadInstances(def.id);
  }

  private async loadInstances(defId: string) {
    this.loadingInstances.set(true);
    try {
      const res = await this.instanceService.list({ processDefinitionId: defId, pageSize: 50 }).toPromise();
      this.instances.set(res?.items ?? []);
    } finally {
      this.loadingInstances.set(false);
    }
  }

  private async selectInstanceById(id: string) {
    try {
      const detail = await this.instanceService.getById(id).toPromise();
      if (!detail) return;
      this.selectedInstance.set({
        id: detail.id,
        zeebeInstanceKey: detail.zeebeInstanceKey,
        processDefinitionId: detail.processDefinitionId,
        status: detail.status,
        currentStep: detail.currentStep,
        assignedAgent: detail.assignedAgent,
        slaStatus: detail.slaStatus,
        createdAt: detail.createdAt,
        completedAt: detail.completedAt,
      });
      this.bpmnXml.set(detail.bpmnXml);
      this.loadEvents(id);
      this.loadHeatmap(detail.processDefinitionId);
      // Select the matching definition
      const def = this.definitions().find(d => d.id === detail.processDefinitionId);
      if (def) this.selectedDef.set(def);
      // Re-render with highlighting after view init
      setTimeout(() => this.renderDiagram(detail.activeElementIds, detail.completedElementIds), 100);
    } catch { /* not found */ }
  }

  onSelectInstance(inst: InstanceSummary) {
    this.selectedInstance.set(inst);
    this.loadingDiagram.set(true);
    this.instanceService.getById(inst.id).subscribe({
      next: (detail) => {
        this.bpmnXml.set(detail.bpmnXml);
        this.loadEvents(inst.id);
        this.loadHeatmap(detail.processDefinitionId);
        setTimeout(() => this.renderDiagram(detail.activeElementIds, detail.completedElementIds), 100);
      },
      error: () => this.loadingDiagram.set(false),
      complete: () => this.loadingDiagram.set(false),
    });
  }

  private loadEvents(instanceId: string) {
    this.instanceService.getEvents(instanceId).subscribe((res) => this.events.set(res.items));
  }

  private loadHeatmap(defId: string) {
    this.instanceService.getHeatmap(defId).subscribe((steps) => this.heatmapData.set(steps));
  }

  private renderDiagram(activeIds?: string[], completedIds?: string[]) {
    const el = this.bpmnEl.nativeElement;
    el.innerHTML = ''; // Clear previous diagram
    if (!this.bpmnXml()) return;

    this.viewer = new BpmnViewer({ container: el });
    this.viewer.importXML(this.bpmnXml()).then(() => {
      const canvas = this.viewer.get('canvas');
      canvas.zoom('fit-viewport');
      if (activeIds?.length) activeIds.forEach((id) => canvas.addMarker(id, 'bpmn-active'));
      if (completedIds?.length) completedIds.forEach((id) => canvas.addMarker(id, 'bpmn-completed'));
    });
  }

  toggleHeatmap() {
    this.showHeatmap.update((v) => !v);
    if (!this.viewer) return;
    const overlays = this.viewer.get('overlays');
    const canvas = this.viewer.get('canvas');

    if (this.showHeatmap()) {
      this.heatmapData().forEach((item) => {
        canvas.addMarker(item.elementId, `heatmap-${item.severity}`);
        overlays.add(item.elementId, {
          position: { bottom: 0, right: 0 },
          html: `<div class="heatmap-badge badge-${item.severity}">${item.instanceCount}</div>`,
        });
      });
    } else {
      this.heatmapData().forEach((item) => canvas.removeMarker(item.elementId, `heatmap-${item.severity}`));
      overlays.clear();
    }
  }

  getStatusSeverity(s: string) {
    switch (s) {
      case 'ACTIVE': return 'info';
      case 'COMPLETED': return 'success';
      case 'FAILED': return 'danger';
      case 'CANCELLED': return 'secondary';
      default: return 'warn';
    }
  }

  getSlaSeverity(s: string) {
    switch (s) {
      case 'ON_TRACK': return 'success';
      case 'WARNING': return 'warn';
      case 'BREACHED': return 'danger';
      default: return 'info';
    }
  }

  getEventIcon(type: string): string {
    if (type.includes('CREATED')) return 'pi pi-plus';
    if (type.includes('COMPLETED')) return 'pi pi-check';
    if (type.includes('ACTIVATED')) return 'pi pi-arrow-right';
    if (type.includes('ERROR')) return 'pi pi-exclamation-triangle';
    return 'pi pi-bell';
  }

  getEventColor(source: string): string {
    switch (source) {
      case 'crm-service': return '#3B82F6';
      case 'loan-service': return '#F59E0B';
      case 'bpm-service': return '#10B981';
      case 'zeebe': return '#8B5CF6';
      default: return '#6366F1';
    }
  }
}
