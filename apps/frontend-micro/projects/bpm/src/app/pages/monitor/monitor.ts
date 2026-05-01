import { Component, ElementRef, OnInit, ViewChild, ChangeDetectionStrategy, signal, inject, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TimelineModule } from 'primeng/timeline';
import { CardModule } from 'primeng/card';
import { TagModule } from 'primeng/tag';
import { ToggleButtonModule } from 'primeng/togglebutton';
import { FormsModule } from '@angular/forms';
import BpmnViewer from 'bpmn-js';
import { KafkaStreamService } from '../../services/kafka-stream.service';

const INITIAL_DIAGRAM = `<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL" xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI" xmlns:dc="http://www.omg.org/spec/DD/20100524/DC" xmlns:di="http://www.omg.org/spec/DD/20100524/DI" id="Definitions_1" targetNamespace="http://bpmn.io/schema/bpmn">
  <bpmn:process id="Process_1" isExecutable="false">
    <bpmn:startEvent id="StartEvent_1" name="Bắt đầu">
      <bpmn:outgoing>Flow_1</bpmn:outgoing>
    </bpmn:startEvent>
    <bpmn:task id="Task_1" name="Nhập liệu hồ sơ">
      <bpmn:incoming>Flow_1</bpmn:incoming>
      <bpmn:outgoing>Flow_2</bpmn:outgoing>
    </bpmn:task>
    <bpmn:sequenceFlow id="Flow_1" sourceRef="StartEvent_1" targetRef="Task_1" />
    <bpmn:exclusiveGateway id="Gateway_1" name="Duyệt hồ sơ?">
      <bpmn:incoming>Flow_2</bpmn:incoming>
      <bpmn:outgoing>Flow_3</bpmn:outgoing>
      <bpmn:outgoing>Flow_4</bpmn:outgoing>
    </bpmn:exclusiveGateway>
    <bpmn:sequenceFlow id="Flow_2" sourceRef="Task_1" targetRef="Gateway_1" />
    <bpmn:task id="Task_2" name="Phê duyệt">
      <bpmn:incoming>Flow_3</bpmn:incoming>
      <bpmn:outgoing>Flow_5</bpmn:outgoing>
    </bpmn:task>
    <bpmn:sequenceFlow id="Flow_3" name="Đồng ý" sourceRef="Gateway_1" targetRef="Task_2" />
    <bpmn:endEvent id="EndEvent_1" name="Kết thúc">
      <bpmn:incoming>Flow_5</bpmn:incoming>
    </bpmn:endEvent>
    <bpmn:sequenceFlow id="Flow_5" sourceRef="Task_2" targetRef="EndEvent_1" />
    <bpmn:task id="Task_3" name="Yêu cầu sửa đổi">
      <bpmn:incoming>Flow_4</bpmn:incoming>
      <bpmn:outgoing>Flow_6</bpmn:outgoing>
    </bpmn:task>
    <bpmn:sequenceFlow id="Flow_4" name="Từ chối" sourceRef="Gateway_1" targetRef="Task_3" />
    <bpmn:sequenceFlow id="Flow_6" sourceRef="Task_3" targetRef="Task_1" />
  </bpmn:process>
  <bpmndi:BPMNDiagram id="BPMNDiagram_1">
    <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="Process_1">
      <bpmndi:BPMNShape id="_BPMNShape_StartEvent_2" bpmnElement="StartEvent_1">
        <dc:Bounds x="173" y="102" width="36" height="36" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Task_1_di" bpmnElement="Task_1">
        <dc:Bounds x="260" y="80" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Gateway_1_di" bpmnElement="Gateway_1" isMarkerVisible="true">
        <dc:Bounds x="415" y="95" width="50" height="50" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Task_2_di" bpmnElement="Task_2">
        <dc:Bounds x="520" y="80" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="EndEvent_1_di" bpmnElement="EndEvent_1">
        <dc:Bounds x="682" y="102" width="36" height="36" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNShape id="Task_3_di" bpmnElement="Task_3">
        <dc:Bounds x="390" y="210" width="100" height="80" />
      </bpmndi:BPMNShape>
      <bpmndi:BPMNEdge id="Flow_1_di" bpmnElement="Flow_1">
        <di:waypoint x="209" y="120" />
        <di:waypoint x="260" y="120" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_2_di" bpmnElement="Flow_2">
        <di:waypoint x="360" y="120" />
        <di:waypoint x="415" y="120" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_3_di" bpmnElement="Flow_3">
        <di:waypoint x="465" y="120" />
        <di:waypoint x="520" y="120" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_5_di" bpmnElement="Flow_5">
        <di:waypoint x="620" y="120" />
        <di:waypoint x="682" y="120" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_4_di" bpmnElement="Flow_4">
        <di:waypoint x="440" y="145" />
        <di:waypoint x="440" y="210" />
      </bpmndi:BPMNEdge>
      <bpmndi:BPMNEdge id="Flow_6_di" bpmnElement="Flow_6">
        <di:waypoint x="390" y="250" />
        <di:waypoint x="310" y="250" />
        <di:waypoint x="310" y="160" />
      </bpmndi:BPMNEdge>
    </bpmndi:BPMNPlane>
  </bpmndi:BPMNDiagram>
</bpmn:definitions>`;

@Component({
  selector: 'app-monitor',
  standalone: true,
  imports: [CommonModule, FormsModule, TimelineModule, CardModule, TagModule, ToggleButtonModule],
  templateUrl: './monitor.html',
  styleUrl: './monitor.css',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MonitorComponent implements OnInit {
  @ViewChild('ref', { static: true }) private el!: ElementRef;
  private viewer: any;
  private kafkaStream = inject(KafkaStreamService);

  showHeatmap = signal(false);

  // Computed-like behavior from Signal stream
  streamEvents = signal<any[]>([]);

  // Mock data for heatmap (number of active instances per step)
  heatmapData = [
    { stepId: 'Task_1', count: 45, severity: 'warn' },
    { stepId: 'Task_2', count: 120, severity: 'danger' },
    { stepId: 'Task_3', count: 5, severity: 'success' }
  ];

  constructor() {
    // React to new Kafka events
    effect(() => {
      const allEvents = this.kafkaStream.events();
      this.streamEvents.set(allEvents.map(e => ({
        status: e.type,
        date: new Date(e.timestamp).toLocaleTimeString(),
        icon: this.getIcon(e.type),
        color: this.getColor(e.source),
        description: `${e.source}: ${JSON.stringify(e.data)}`
      })));
    });
  }

  ngOnInit(): void {
    this.viewer = new BpmnViewer({
      container: this.el.nativeElement
    });

    this.viewer.importXML(INITIAL_DIAGRAM).then(() => {
      this.viewer.get('canvas').zoom('fit-viewport');
    });
  }

  toggleHeatmap() {
    this.showHeatmap.update(v => !v);
    const overlays = this.viewer.get('overlays');
    const canvas = this.viewer.get('canvas');

    if (this.showHeatmap()) {
      this.heatmapData.forEach(item => {
        canvas.addMarker(item.stepId, `heatmap-${item.severity}`);
        overlays.add(item.stepId, {
          position: { bottom: 0, right: 0 },
          html: `<div class="heatmap-badge badge-${item.severity}">${item.count}</div>`
        });
      });
    } else {
      this.heatmapData.forEach(item => {
        canvas.removeMarker(item.stepId, `heatmap-${item.severity}`);
      });
      overlays.clear();
    }
  }

  private getIcon(type: string): string {
    if (type.includes('CREATED')) return 'pi pi-plus';
    if (type.includes('SUBMITTED')) return 'pi pi-file-export';
    if (type.includes('WORKFLOW')) return 'pi pi-play';
    return 'pi pi-bell';
  }

  private getColor(source: string): string {
    switch (source) {
      case 'crm-service': return '#3B82F6';
      case 'loan-service': return '#F59E0B';
      case 'bpm-service': return '#10B981';
      default: return '#6366F1';
    }
  }
}
