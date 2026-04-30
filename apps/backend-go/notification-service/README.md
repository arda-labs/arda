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

Delivery queue, provider adapters, in-app inbox, preferences, and retry worker
are intentionally left for later parts.

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
