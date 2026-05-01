# Arda

Arda is the application monorepo for a financial and banking platform. It
contains the frontend micro-frontends, Go operational services, Java
core-banking services, shared libraries, CI workflows, and documentation.

Updated: 2026-05-02

## Current Reality

The active repositories are:

| Repository | Owns |
| --- | --- |
| `arda-labs/arda` | Application code, shared libraries, docs, app CI |
| `arda-labs/arda-infra` | Kubernetes manifests, APISIX routes, ArgoCD apps, runtime config |
| `arda-labs/.github` | Organization profile and GitHub metadata |

## Repository Structure

```text
arda/
├── apps/
│   ├── frontend-micro/          # Angular CLI workspace with Native Federation
│   │   └── projects/
│   │       ├── shell/           # Host app (4200)
│   │       ├── iam/             # IAM remote (4201)
│   │       ├── mdm/             # MDM remote (4202)
│   │       ├── ntf/             # Notification remote (4204)
│   │       ├── crm/             # CRM remote (4210)
│   │       └── core/            # Shared Angular library
│   ├── backend-go/              # Go workspace for Kratos services
│   │   ├── iam-service/         # Identity (8000/9000)
│   │   ├── mdm-service/         # Master Data (8001/9001)
│   │   ├── media-service/       # Storage (8002/9002)
│   │   ├── bpm-service/         # Zeebe worker (8003/9003)
│   │   └── notification-service/# Delivery (8004/9004)
│   └── backend-java/            # Gradle workspace (Java 25 + Virtual Threads)
│       ├── crm-service/         # Customer mgmt (8010/9010)
│       ├── hrm-service/         # Human resources (8011/9011)
│       └── loan-service/        # Lending core (8012/9012)
├── libs/
│   ├── go/pkg/                  # Shared Go helpers
│   └── java/                    # Shared Java libs (Imperative)
├── docs/                        # Architecture, feature, and operating docs
└── .github/workflows/           # CI and GitOps update workflows
```

Runtime manifests live in the sibling repo `../arda-infra`.

## Implemented Modules

| Area | Status | Port (Dev) |
| --- | --- | --- |
| Shell MFE | Active | 4200 |
| IAM MFE | Active | 4201 |
| MDM MFE | Active | 4202 |
| IAM service | Active (Go) | 8000 / 9000 |
| MDM service | Active (Go) | 8001 / 9001 |
| CRM service | Active (Java 25) | 8010 / 9010 |

## Local Development

Standardized port mapping:

- Frontend: Shell (4200), Go-MFEs (4201-4209), Java-MFEs (4210-4219).
- Backend Go: HTTP (800x), gRPC (900x).
- Backend Java: HTTP (801x), gRPC (901x).

```powershell
cd D:\Github\arda-labs\arda-infra\local\apisix
docker compose up -d
```

Run frontend remotes from `apps/frontend-micro`:

```powershell
npm install
npx ng serve shell
npx ng serve iam
npx ng serve mdm
npx ng serve ntf
```

Run Go services from `apps/backend-go`:

```powershell
cd apps\backend-go\iam-service
kratos run

cd ..\mdm-service
kratos run

cd ..\notification-service
go run .\cmd\notification-service -conf .\configs
```

Open the shell through APISIX:

```text
http://localhost:9080
```

Main local gateway routes:

| Route | Target |
| --- | --- |
| `/api/v1/*` | IAM service, rewritten to `/v1/*` |
| `/api/v1/mdm/*` | MDM service, rewritten to `/v1/mdm/*` |
| `/api/v1/notifications/*` | Notification service, rewritten to `/v1/notifications/*` |
| `/mfe-iam/*` | IAM remote assets |
| `/mfe-mdm/*` | MDM remote assets |
| `/mfe-ntf/*` | NTF remote assets |
| `/*` | Shell app |

## Build And Test

Frontend:

```powershell
cd apps\frontend-micro
npx ng build shell
npx ng build iam
npx ng build mdm
npx ng build ntf
```

Go services:

```powershell
cd apps\backend-go\iam-service
go test ./...

cd ..\mdm-service
go test ./...

cd ..\notification-service
go test ./...
```

Infra manifests:

```powershell
cd ..\..\..\arda-infra
kubectl kustomize apps\iam-service\overlays\dev
kubectl kustomize apps\mdm-service\overlays\dev
kubectl kustomize apps\mfe-shell\overlays\dev
kubectl kustomize apps\mfe-iam\overlays\dev
kubectl kustomize apps\mfe-mdm\overlays\dev
```

## Documentation

- [Documentation Center](./docs/README.md)
- [Operating Model](./docs/00-operating-model.md)
- [Architecture Overview](./docs/01-overview/architecture.md)
- [Frontend Architecture](./docs/04-frontend/architecture.md)
- [Go Backend Architecture](./docs/02-backend-go/architecture.md)
- [MDM Features](./docs/06-features/mdm.md)
- [Notification Service Design](./docs/06-features/notification.md)
- [Documentation Audit](./docs/08-guides/documentation-audit.md)

## License

Copyright © 2026 Arda Labs. All rights reserved.
