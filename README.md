# Arda

Arda is the application monorepo for a financial and banking platform. It
contains the frontend micro-frontends, Go operational services, Java/Kotlin
core-banking prototypes, shared libraries, CI workflows, and documentation.

Updated: 2026-04-30

## Current Reality

The active repositories are:

| Repository | Owns |
| --- | --- |
| `arda-labs/arda` | Application code, shared libraries, docs, app CI |
| `arda-labs/arda-infra` | Kubernetes manifests, APISIX routes, ArgoCD apps, runtime config |
| `arda-labs/.github` | Organization profile and GitHub metadata |

There are no active split repos such as `arda-mfe`, `arda-be`, or
`arda-core`. Those names appear only in older planning documents.

## Repository Structure

```text
arda/
├── apps/
│   ├── frontend-micro/          # Angular CLI workspace with Native Federation
│   │   └── projects/
│   │       ├── shell/           # Host app, layout, auth callback, workspace UI
│   │       ├── iam/             # IAM remote MFE
│   │       ├── mdm/             # MDM remote MFE
│   │       ├── ntf/             # Notification operations remote MFE
│   │       └── core/            # Shared Angular library
│   ├── backend-go/              # Go workspace for Kratos services
│   │   ├── iam-service/         # Identity, tenants, menus, permissions
│   │   ├── mdm-service/         # Master Data Management
│   │   ├── notification-service/# Notification templates and delivery roadmap
│   │   └── crm-service/         # Skeleton / roadmap service
│   └── backend-java/            # Gradle workspace
│       └── accounting_tmp/      # Accounting prototype
├── libs/
│   └── go/pkg/                  # Shared Go helpers
├── docs/                        # Architecture, feature, and operating docs
└── .github/workflows/           # CI and GitOps update workflows
```

Runtime manifests live in the sibling repo `../arda-infra`.

## Implemented Modules

| Area | Status |
| --- | --- |
| Shell MFE | Active, runs on port `3000`, loads remotes from runtime `env.js` |
| IAM MFE | Active remote, runs on port `3002`, route `/app/iam/*` |
| MDM MFE | Active remote, runs on port `3001`, route `/app/mdm/*` |
| NTF MFE | Active remote, runs on port `3003`, route `/app/ntf/*` |
| IAM service | Active Go/Kratos service, default HTTP `8000`, gRPC `9000` |
| MDM service | Active Go/Kratos service, default HTTP `8001`, gRPC `9001` |
| Notification service | Active Go/Kratos service; templates, delivery queue, in-app inbox, provider config |
| CRM service | Present in Go workspace as skeleton/roadmap |
| Accounting Java | Prototype under `apps/backend-java/accounting_tmp` |

## Local Development

Use APISIX for integration checks so local traffic has the same path shape as
deployed traffic.

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
