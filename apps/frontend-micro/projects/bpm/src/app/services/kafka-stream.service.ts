import { Injectable, signal } from '@angular/core';

export interface ArdaEvent {
  id: string;
  type: string;
  source: string;
  timestamp: string;
  data: any;
  traceId: string;
}

@Injectable({
  providedIn: 'root'
})
export class KafkaStreamService {
  private socket?: WebSocket;

  // Real-time events stream using Signals
  events = signal<ArdaEvent[]>([]);

  // Latest event received
  latestEvent = signal<ArdaEvent | null>(null);

  constructor() {
    // In a real environment, this would connect to the Go Gateway or BPM service
    // this.connect('ws://localhost:8000/ws/events');

    // Simulate events for demo purposes until infrastructure is up
    this.startSimulation();
  }

  connect(url: string) {
    this.socket = new WebSocket(url);
    this.socket.onmessage = (event) => {
      const ardaEvent: ArdaEvent = JSON.parse(event.data);
      this.pushEvent(ardaEvent);
    };
  }

  private pushEvent(event: ArdaEvent) {
    this.events.update(prev => [event, ...prev.slice(0, 49)]); // Keep last 50
    this.latestEvent.set(event);
  }

  private startSimulation() {
    const mockEvents = [
      { type: 'CUSTOMER_CREATED', source: 'crm-service', data: { id: 'C100', name: 'Nguyễn Văn A' } },
      { type: 'LOAN_APPLICATION_SUBMITTED', source: 'loan-service', data: { id: 'L999', amount: 500000000 } },
      { type: 'WORKFLOW_STARTED', source: 'bpm-service', data: { processId: 'wf-123', status: 'ACTIVE' } },
      { type: 'NOTIFICATION_SENT', source: 'notification-service', data: { to: 'user@example.com', template: 'welcome' } }
    ];

    setInterval(() => {
      const mock = mockEvents[Math.floor(Math.random() * mockEvents.length)];
      const event: ArdaEvent = {
        id: Math.random().toString(36).substring(7),
        type: mock.type,
        source: mock.source,
        timestamp: new Date().toISOString(),
        data: mock.data,
        traceId: 'tr-' + Math.random().toString(36).substring(7)
      };
      this.pushEvent(event);
    }, 15000); // Every 15 seconds for realism
  }
}
