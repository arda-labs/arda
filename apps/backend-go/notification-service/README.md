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

Part 3 implements provider configuration records:

- `notification_provider_configs` stores provider code, channel, priority,
  rate limit, status, and non-secret options.
- `IN_APP_STORE` is seeded as the active in-app provider.
- API can list and upsert provider configs without storing secrets in DB.

External provider adapters, runtime secrets, preferences, quiet-hour policy,
audit UI, and operations dashboard are intentionally left for later parts.

## Main APIs

```text
POST /v1/notifications/requests
GET  /v1/notifications/requests/{id}
GET  /v1/notifications/deliveries
POST /v1/notifications/deliveries/{id}/retry
POST /v1/notifications/deliveries/run-once
GET  /v1/notifications/in-app?recipient_type=USER&recipient_id={id}
POST /v1/notifications/in-app/{id}/read
GET  /v1/notifications/provider-configs
PUT  /v1/notifications/provider-configs/{code}
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

Create the local/runtime database before starting the service. In the infra repo:

```bash
psql -h thinkcenter -U postgres -f scripts/bootstrap-dev-postgres.sql
```

The service also accepts `DATABASE_URL`, so a workstation PostgreSQL can be used
without editing `configs/config.yaml`:

```powershell
$env:DATABASE_URL='postgres://notification:notification%40123@localhost:5432/notification?sslmode=disable'
go run .\cmd\notification-service -conf .\configs
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

The production container reads `/data/conf/config.yaml`; Kubernetes renders it
from `arda-infra/apps/notification-service/base/configs.yaml`.

```powershell
docker build -f apps/backend-go/notification-service/Dockerfile -t ghcr.io/arda-labs/notification-service:dev .
```
