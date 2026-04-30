# Notification Service

Updated: 2026-04-30

Status: Implementation started. `notification-service` now has template APIs,
durable notification requests, delivery queue, retry API, DB-backed in-app
delivery, provider configuration records, shell inbox polling, and an SSE
stream for realtime in-app updates. External provider adapters, preferences,
and Kafka consumers are still roadmap.

## Goal

`notification-service` will own user and customer notifications across the
platform. It should be designed as a governed delivery service, not as a helper
library inside IAM, MDM, CRM, or core banking services.

The service should support:

- In-app notifications for users inside Arda.
- Email delivery.
- SMS delivery where provider integration is available.
- Zalo OA/ZNS delivery for Vietnam customer communication.
- Webhook delivery for system-to-system callbacks.
- Template, preference, queue, retry, audit, and provider failover management.

## Service Boundary

Notification service owns:

- Message templates and template versions.
- Notification requests and delivery jobs.
- Recipient resolution snapshot at send time.
- Channel preferences and quiet-hour rules.
- Provider configuration references and routing rules.
- Delivery state, retry state, provider response, and audit trail.

Notification service does not own:

- Business decisions about whether a transaction should notify a customer.
- Customer master data. It may cache recipient snapshots but source records
  remain in CRM/IAM/core banking.
- Authentication or role/menu authorization. That remains in IAM.
- MDM catalogs. Notification may reference MDM code sets for channel, status,
  language, country, or holiday/business calendar behavior.

## Proposed Runtime Shape

| Item | Decision |
| --- | --- |
| Service name | `notification-service` |
| Backend stack | Go/Kratos, same service layout as IAM and MDM |
| Database | Dedicated `notification` PostgreSQL database |
| API prefix | Native `/v1/notifications/*`, gateway `/api/v1/notifications/*` |
| Frontend MFE | Dedicated `ntf` MFE plus shell notification bell |
| Async mechanism | Redpanda/Kafka for cross-service events; DB queue remains source of truth |
| Primary consumers | IAM, CRM, payments, loan, accounting, workflow/BPM |

## Target Event Architecture

```text
Domain services
  -> Redpanda/Kafka topic: arda.notification.events.v1
  -> notification-service consumer
  -> notification_requests + notification_deliveries + in_app_notifications
  -> SSE stream to shell bell
  -> polling fallback and periodic reconciliation
```

Kafka carries domain events and supports decoupling, replay, and scalable
consumers. PostgreSQL remains the durable notification source of truth. SSE is
the primary browser realtime path for in-app notifications, while polling is the
fallback and periodic sync mechanism.

## Core Flows

### Event To Notification

1. Domain service decides a business event requires notification.
2. Domain service calls `CreateNotification` or writes an outbox event.
3. Notification service validates template/channel/policy.
4. Service creates one notification request and one or more delivery jobs.
5. Worker claims queued jobs and delivers through the channel adapter.
6. Provider result is recorded with delivery status and response metadata.
7. Failed jobs retry with backoff until max attempts, then move to dead-letter.

### Direct Send

Use direct send only for low-risk internal notifications or admin tools.
Financial/customer-facing messages should go through durable queue jobs so they
survive service restarts and have full audit history.

### Template Lifecycle

1. Draft a template with channel-specific content.
2. Validate placeholders and sample payload.
3. Approve the template version.
4. Send jobs always reference a concrete template version.
5. Editing a template creates a new version rather than mutating historical
   delivery context.

## Data Model

Initial tables:

| Table | Purpose |
| --- | --- |
| `notification_templates` | Template header: code, name, category, status, default channel |
| `notification_template_versions` | Versioned content per channel/language with schema and approval status |
| `notification_requests` | Implemented. Business request: source service, event type, correlation id, recipient, payload snapshot |
| `notification_deliveries` | Implemented. Per-channel delivery job with status, attempt count, provider result |
| `in_app_notifications` | Implemented. Internal inbox records created from delivered `IN_APP` jobs |
| `notification_preferences` | User/customer channel preference, opt-in/opt-out, quiet hours |
| `notification_provider_configs` | Implemented. Provider code, channel, priority, rate limit, status, non-secret metadata |
| `notification_audit_logs` | Template changes, approval actions, resend/manual override actions |

Recommended delivery statuses:

```text
QUEUED, CLAIMED, RENDERED, SENT, DELIVERED, FAILED, RETRYING, DEAD_LETTER, CANCELLED
```

Recommended channels:

```text
IN_APP, EMAIL, SMS, ZALO_OA, ZALO_ZNS, WEBHOOK
```

## API Design

External service-facing APIs:

| Method | Purpose |
| --- | --- |
| `POST /v1/notifications/requests` | Create durable notification request |
| `GET /v1/notifications/requests/{id}` | Inspect request and deliveries |
| `POST /v1/notifications/requests/{id}/cancel` | Cancel queued deliveries |
| `GET /v1/notifications/deliveries` | List delivery queue jobs |
| `POST /v1/notifications/deliveries/{id}/retry` | Manual retry for failed delivery |
| `POST /v1/notifications/deliveries/run-once` | Admin/local hook to process one delivery batch |
| `GET /v1/notifications/in-app` | List in-app inbox records by recipient |
| `GET /v1/notifications/in-app/unread-count` | Count unread in-app inbox records |
| `GET /v1/notifications/in-app/stream` | SSE stream for realtime in-app updates |
| `POST /v1/notifications/in-app/{id}/read` | Mark an in-app notification as read |
| `POST /v1/notifications/in-app/read-all` | Mark all recipient in-app notifications as read |

Admin APIs:

| Method | Purpose |
| --- | --- |
| `GET /v1/notifications/templates` | List templates |
| `POST /v1/notifications/templates` | Create template |
| `PUT /v1/notifications/templates/{id}` | Update template header |
| `POST /v1/notifications/templates/{id}/versions` | Create new version |
| `POST /v1/notifications/template-versions/{id}/approve` | Approve version |
| `GET /v1/notifications/provider-configs` | List providers |
| `PUT /v1/notifications/provider-configs/{code}` | Create or update a non-secret provider config |
| `PUT /v1/notifications/preferences/{subject_id}` | Update preferences |

Frontend pages for v1:

- Template management.
- Notification queue/delivery monitor.
- Provider config status.
- Preference lookup by user/customer.

## Provider Adapter Contract

Each adapter should expose a small internal interface:

```text
Send(ctx, delivery, renderedMessage) -> provider_message_id, status, raw_response
ValidateTemplate(ctx, templateVersion) -> validation result
Health(ctx) -> provider availability
```

Provider specifics stay inside adapters:

- SMTP or transactional email provider.
- SMS provider.
- Zalo OA.
- Zalo ZNS.
- Webhook endpoint with signing.
- In-app delivery store.

Secrets such as SMTP password, Zalo app secret, or API keys should live in
Kubernetes secrets/runtime config, not in database rows. Database rows should
store provider code, priority, rate limit, and non-secret options.

## Financial Platform Requirements

Notification in a banking platform needs stricter behavior than a generic SaaS
mailer:

- Idempotency by `source_service + event_type + correlation_id + recipient`.
- Full audit trail for customer-facing sends.
- Template version pinning for legal and compliance traceability.
- Quiet-hour and opt-out handling, with override for mandatory security alerts.
- Rate limiting per provider and per recipient.
- Retry/backoff with dead-letter inspection.
- PII minimization in logs; payload snapshots should avoid storing secrets,
  credentials, card numbers, or account balances unless explicitly required.
- Maker/checker approval for customer-facing template versions.
- Delivery evidence export for compliance and dispute handling.

## Relationship With MDM

Notification should consume MDM reference data rather than duplicate it:

- Business calendars for quiet-hour/business-day send windows.
- Languages/countries/code sets for template localization.
- Channel and status code sets where shared semantics matter.
- System parameters for global defaults such as retention windows or retry
  limits, if those are not notification-owned.

MDM should not send messages. It only provides reference data.

## Implementation Plan

Phase 1: service skeleton and durable queue

- Done: create `apps/backend-go/notification-service`.
- Done: add dedicated database config and first migrations.
- Done: implement template CRUD and template-version approval.
- Done: implement request creation, delivery queue, retry API, and worker loop.
- Done: support real `IN_APP` delivery by storing inbox records in database.
- Done: add provider configuration registry with non-secret metadata.
- Next: support `EMAIL` provider adapter and runtime secret loading.

Phase 2: operations UI

- Add notification pages to frontend.
- Template list/detail/version approval.
- Delivery queue monitor with retry/cancel.
- Provider health/status view.

Phase 3: Vietnam banking channels

- Add Zalo OA/ZNS adapters.
- Add SMS adapter.
- Add preference and quiet-hour enforcement.
- Add exportable delivery audit.

Phase 4: platform event integration

- Standardize outbox event envelope.
- Add Redpanda/Kafka consumer for `arda.notification.events.v1`.
- Integrate IAM security events.
- Integrate CRM customer lifecycle events.
- Integrate future payments/loan/accounting events.

## Open Decisions

- Whether notification service starts as its own MFE or pages inside shell/system
  until operational scope grows.
- Which email/SMS/Zalo providers are approved for local/dev/prod.
- Whether Redpanda/Kafka is introduced before high-volume transactional
  notifications.
- Retention period for payload snapshots and delivery logs.
