export interface PageResponse<T> {
  items: T[];
  nextPageToken: string;
}

export interface ListOptions {
  status?: string;
  channel?: string;
  keyword?: string;
  recipientId?: string;
  pageSize?: number;
  pageToken?: string;
}

export interface NotificationTemplate {
  id: string;
  code: string;
  name: string;
  category: string;
  defaultChannel: string;
  description: string;
  status: string;
}

export interface NotificationTemplateVersion {
  id: string;
  templateId: string;
  version: number;
  channel: string;
  language: string;
  subject: string;
  body: string;
  payloadSchemaJson: string;
  approvalStatus: string;
  approvedBy: string;
  changeNote: string;
  status: string;
}

export interface NotificationDelivery {
  id: string;
  requestId: string;
  templateVersionId: string;
  channel: string;
  recipientType: string;
  recipientId: string;
  recipientAddress: string;
  subject: string;
  body: string;
  status: string;
  attemptCount: number;
  maxAttempts: number;
  providerCode: string;
  providerMessageId: string;
  providerResponseJson: string;
  errorMessage: string;
  priority: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProviderConfig {
  id: string;
  code: string;
  channel: string;
  name: string;
  priority: number;
  rateLimitPerMinute: number;
  optionsJson: string;
  status: string;
}
