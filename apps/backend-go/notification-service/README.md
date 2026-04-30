# Notification Service

Notification service owns governed notification templates, template versions,
delivery requests, queue jobs, provider adapters, retry state, and delivery
audit.

## Current Scope

Part 1 implements the template registry:

- Notification template CRUD.
- Versioned template content per channel and language.
- Template version approval.
- Dedicated notification database migrations.

Part 2 implements durable in-app delivery:

- `notification_requests` with idempotency key.
- `notification_deliveries` queue jobs per requested channel.
- `in_app_notifications` inbox storage for the `IN_APP` channel.
- Retry API for failed/dead-letter jobs.
- Background worker loop that claims due jobs and stores in-app messages.

External provider adapters, preferences, quiet-hour policy, audit UI, and
operations dashboard are intentionally left for later parts.

## Main APIs

```text
POST /v1/notifications/requests
GET  /v1/notifications/requests/{id}
GET  /v1/notifications/deliveries
POST /v1/notifications/deliveries/{id}/retry
POST /v1/notifications/deliveries/run-once
GET  /v1/notifications/in-app?recipient_type=USER&recipient_id={id}
POST /v1/notifications/in-app/{id}/read
```

Minimal request example:

```json
{
  "request": {
    "sourceService": "IAM",
    "eventType": "SECURITY_LOGIN",
    "correlationId": "login-123",
    "templateCode": "IAM_SECURITY_LOGIN",
    "recipientType": "USER",
    "recipientId": "user-001",
    "channels": ["IN_APP"],
    "language": "vi",
    "payloadJson": "{\"login_time\":\"2026-04-30 22:00\",\"ip_address\":\"127.0.0.1\"}"
  }
}
```

## Local Run

Default local database:

```text
postgres://notification:notification%40123@thinkcenter:5432/notification?sslmode=disable
```

Ports:

- HTTP: `8002`
- gRPC: `9002`

Run:

```powershell
go run ./cmd/notification-service -conf ./configs
```

## Verify

```powershell
go test ./...
```

## Docker

```powershell
docker build -f apps/backend-go/notification-service/Dockerfile -t ghcr.io/arda-labs/notification-service:dev .
```
