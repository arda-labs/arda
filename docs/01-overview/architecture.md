# Architecture Overview

Updated: 2026-04-30

Arda is currently a compact application monorepo backed by a separate GitOps
repository. The long-term target is a financial microservice platform, but the
current implementation should be described narrowly and accurately.

## Current System Shape

```text
Browser
  |
  v
APISIX Gateway
  |-- /*              -> mfe-shell
  |-- /mfe-iam/*      -> mfe-iam remote assets
  |-- /mfe-mdm/*      -> mfe-mdm remote assets
  |-- /api/v1/*       -> iam-service
  |-- /api/v1/mdm/*   -> mdm-service
  |
  v
Go/Kratos and Java/Spring services
  |-- iam-service (Go)
  |-- mdm-service (Go)
  |-- crm-service (Java 25, SB 4, Camunda 8)
  |
  v
PostgreSQL on thinkcenter
```

Zitadel provides OIDC authentication. IAM owns tenant membership, roles,
permissions, menus, and forward-auth policy once APISIX auth is formalized.

## Repository Boundaries

| Repository | Role |
| --- | --- |
| `arda` | Application source code, shared libs, CI, docs |
| `arda-infra` | Kubernetes manifests, ArgoCD apps, APISIX routes, runtime config |
| `.github` | Organization profile and metadata |

Do not put runtime Kubernetes state into `arda`. Do not put application
implementation code into `arda-infra`.

## Frontend Layer

The frontend is an Angular CLI workspace, not an Nx workspace.

```text
apps/frontend-micro/
├── angular.json
├── Dockerfile
├── nginx.conf
└── projects/
    ├── shell/       # host app, layout, auth callback, workspace UI
    ├── iam/         # remote MFE
    ├── mdm/         # remote MFE
    └── core/        # shared Angular library
```

The shell initializes Native Federation from runtime `env.js`:

```js
window.__env.mfeIamUrl = 'http://localhost:9080/mfe-iam';
window.__env.mfeMdmUrl = 'http://localhost:9080/mfe-mdm';
```

Shell routes:

| Shell route | Owner |
| --- | --- |
| `/home` | Shell |
| `/settings` | Shell |
| `/workspaces` | Shell |
| `/iam/*` | IAM remote |
| `/mdm/*` | MDM remote |

## Backend Go Layer

```text
apps/backend-go/
├── go.work
├── iam-service/
├── mdm-service/
└── crm-service/
```

Active services:

| Service | Status | HTTP | gRPC | Database |
| --- | --- | --- | --- | --- |
| `iam-service` | Active | `8000` | `9000` | `iam` |
| `mdm-service` | Active | `8001` | `9001` | `mdm` |
| `notification-service` | Started | `8002` | `9002` | `notification` |
| `crm-service` | Active (Java) | `8003` | N/A | `crm` |

Service-native HTTP routes use `/v1/*`; APISIX exposes them as `/api/v1/*`.

## Backend Java Layer

```text
apps/backend-java/
├── build.gradle.kts
├── settings.gradle.kts
└── accounting_tmp/
```

`accounting_tmp` is a prototype, not a deployable production service yet. Java
CI and docs should be aligned when it is renamed to the final `accounting`
module.

## Infrastructure Layer

Runtime manifests are in `arda-infra`:

```text
arda-infra/
├── argocd/
├── apps/
│   ├── gateway/apisix/
│   ├── identity/zitadel/
│   ├── iam-service/
│   ├── mdm-service/
│   ├── mfe-shell/
│   ├── mfe-iam/
│   └── mfe-mdm/
├── infrastructure/
└── local/apisix/
```

Namespaces currently used:

| Namespace | Purpose |
| --- | --- |
| `argocd` | ArgoCD |
| `gateway` | APISIX |
| `identity` | Zitadel |
| `infra` | Cloudflared |
| `arda-apps` | Application workloads |

## Current Versus Roadmap

Implemented now:

- Angular shell, IAM remote, MDM remote.
- IAM and MDM Go services.
- MDM schema, seed data, API, and UI.
- GitOps manifests for IAM, MDM, shell, IAM MFE, and MDM MFE.
- Local standalone APISIX for workstation integration checks.

Roadmap or prototype:

- CRM, HRM, BPM, loan, deposit, treasury services.
- Notification delivery queue, provider adapters, preferences, and operations UI.
- Accounting Java production service.
- Redpanda/outbox production flow.
- Camunda, Garage S3, Prometheus/Grafana runtime stack.
- Formal APISIX forward-auth plugin wiring.

## Related Docs

- [Frontend Architecture](../04-frontend/architecture.md)
- [Go Backend Architecture](../02-backend-go/architecture.md)
- [Java Backend Architecture](../03-backend-java/architecture.md)
- [Infrastructure Guide](../05-infrastructure/README.md)
- [MDM](../06-features/mdm.md)
