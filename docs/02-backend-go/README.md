# Backend Go Guide

Updated: 2026-04-30

This guide describes the current Go services in `apps/backend-go`.

## Current Status

| Service | Status | Notes |
| --- | --- | --- |
| `iam-service` | Active | Identity integration, tenants, roles, menus, permissions |
| `mdm-service` | Active | Master data, geography, code catalogs, system parameters |
| `notification-service` | Started | Template registry implemented; delivery queue/provider adapters next |
| `crm-service` | Skeleton | Keep as roadmap until implementation starts |

## Common Commands

Run tests:

```powershell
cd apps\backend-go\iam-service
go test ./...

cd ..\mdm-service
go test ./...

cd ..\notification-service
go test ./...
```

Run locally:

```powershell
cd apps\backend-go\iam-service
kratos run

cd ..\mdm-service
kratos run

cd ..\notification-service
go run ./cmd/notification-service -conf ./configs
```

Build a container from the repo root:

```powershell
docker build -f apps/backend-go/iam-service/Dockerfile -t ghcr.io/arda-labs/iam-service:dev .
docker build -f apps/backend-go/mdm-service/Dockerfile -t ghcr.io/arda-labs/mdm-service:dev .
docker build -f apps/backend-go/notification-service/Dockerfile -t ghcr.io/arda-labs/notification-service:dev .
```

## Service Ports

| Service | HTTP | gRPC |
| --- | --- | --- |
| `iam-service` | `8000` | `9000` |
| `mdm-service` | `8001` | `9001` |
| `notification-service` | `8002` | `9002` |

## API Contract

Service-native paths start with `/v1`. APISIX exposes them under `/api/v1`.

| Gateway path | Native path | Service |
| --- | --- | --- |
| `/api/v1/*` | `/v1/*` | IAM |
| `/api/v1/mdm/*` | `/v1/mdm/*` | MDM |

## Adding A Go Service

1. Keep the service name short and domain-owned, for example `crm-service`.
2. Add it to `apps/backend-go/go.work`.
3. Keep the Kratos layer structure: `biz`, `data`, `service`, `server`, `conf`.
4. Put schema migrations under `internal/data/migrations`.
5. Add Dockerfile and GitHub Actions detection.
6. Add matching manifests in `arda-infra/apps/<service>/overlays/dev`.
7. Add APISIX routes in infra, not in the app repo.

## Next Planned Service: Notification

`notification-service` has started with the same Kratos layout as IAM and MDM.
Part 1 implements the template registry and uses a dedicated PostgreSQL
database. Durable delivery queue and provider adapters are next. The feature
design is documented in
[Notification Service](../06-features/notification.md).

## Current Gaps

- Java CI still references a future `accounting` module; it does not affect Go
  services but should be cleaned before adding cross-service workflows.
- `crm-service` should either be implemented or removed from active workspace
  scans if it causes CI noise.
- Shared Go libraries should keep module dependencies explicit to avoid hidden
  `go.work` coupling.
