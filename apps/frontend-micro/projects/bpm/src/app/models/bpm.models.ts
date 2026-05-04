export interface ProcessDefinition {
  id: string;
  processKey: string;
  name: string;
  description: string;
  category: string;
  module: string;
  version: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface InstanceSummary {
  id: string;
  zeebeInstanceKey: string;
  processDefinitionId: string;
  status: string;
  currentStep: string;
  assignedAgent: string;
  slaStatus: string;
  createdAt: string;
  completedAt?: string;
}

export interface InstanceDetail {
  id: string;
  zeebeInstanceKey: string;
  processDefinitionId: string;
  processName: string;
  status: string;
  currentStep: string;
  variables: Record<string, unknown>;
  assignedAgent: string;
  slaStatus: string;
  bpmnXml: string;
  activeElementIds: string[];
  completedElementIds: string[];
  createdAt: string;
  completedAt?: string;
}

export interface ProcessEvent {
  id: string;
  processInstanceId: string;
  eventType: string;
  source: string;
  data: Record<string, unknown>;
  timestamp: string;
}

export interface HeatmapStep {
  elementId: string;
  instanceCount: number;
  avgDurationSeconds: number;
  severity: string;
}

export interface TemplateVariable {
  id?: string;
  templateId?: string;
  variableName: string;
  sourceType: string;
  sourceField: string;
  fallbackValue: string;
}

export interface Template {
  id: string;
  processDefinitionId: string;
  name: string;
  templateText: string;
  module: string;
  variables: TemplateVariable[];
  createdAt: string;
  updatedAt: string;
}

export interface VariableSource {
  type: string;
  name: string;
  description: string;
}

export interface VariableMapping {
  id: string;
  templateId: string;
  variableName: string;
  sourceType: string;
  sourceField: string;
  fallbackValue: string;
}

export interface ResolvedVariable {
  variableName: string;
  resolvedValue: string;
  usedFallback: boolean;
  sourceType: string;
}

export interface RenderResult {
  renderedText: string;
  resolvedVariables: ResolvedVariable[];
}

export interface PageResponse<T> {
  items: T[];
  nextPageToken: string;
}

/* ───── Definitions & Deploy ───── */

export interface DefinitionDiagram {
  id: string;
  bpmnXml: string;
  processKey: string;
  version: number;
}

export interface DeployDefinitionRequest {
  bpmnXml: string;
  processKey: string;
  name: string;
  description?: string;
  category?: string;
  module?: string;
}

/* ───── SLA Config (stub) ───── */

export interface SlaConfig {
  id: string;
  processName: string;
  stepName: string;
  durationHours: number;
  unit: 'MINUTES' | 'HOURS' | 'DAYS';
  warningPercent: number;
}

/* ───── Assignment Rules (stub) ───── */

export interface AssignmentRule {
  id: string;
  module: string;
  strategy: 'ROUND_ROBIN' | 'LOAD_BALANCE' | 'WEIGHTED' | 'MANUAL';
  agents: AgentAssignment[];
}

export interface AgentAssignment {
  id: string;
  name: string;
  department: string;
  weight: number;
  active: boolean;
  taskCount: number;
}

/* ───── Error Hospital (stub) ───── */

export interface FailedTask {
  id: string;
  instanceId: string;
  processName: string;
  stepName: string;
  error: string;
  payload: string;
  retryCount: number;
  lastAttempt: string;
  status: 'PENDING' | 'RETRYING' | 'RESOLVED';
}

/* ───── Task (Inbound/Outbound stub) ───── */

export interface BpmTask {
  id: string;
  title: string;
  module: string;
  priority: 'HIGH' | 'MEDIUM' | 'LOW';
  slaStatus: 'ON_TRACK' | 'WARNING' | 'BREACHED';
  startTime: string;
  customerId: string;
  processId: string;
  status: 'IN_PROGRESS' | 'COMPLETED' | 'RETURNED' | 'FAILED';
  assignedTo?: string;
}
