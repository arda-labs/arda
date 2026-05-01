# Architecture Overview

Updated: 2026-05-02

Arda is currently a compact application monorepo backed by a separate GitOps
repository. The long-term target is a financial microservice platform.

## Current System Shape

```text
Browser
  |
  v
APISIX Gateway (9080)
  |-- /*              -> mfe-shell (4200)
  |-- /mfe-iam/*      -> mfe-iam (4201)
  |-- /mfe-mdm/*      -> mfe-mdm (4202)
  |-- /api/v1/*       -> iam-service (8000)
  |-- /api/v1/mdm/*   -> mdm-service (8001)
  |
  v
Backend Services
  |-- Go: IAM (8000), MDM (8001), Media (8002), BPM (8003), NTF (8004)
  |-- Java: CRM (8010), HRM (8011), Loan (8012)
  |
  v
PostgreSQL on thinkcenter
```

Zitadel provides OIDC authentication. IAM owns tenant membership and roles.

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
window.__env.mfeIamUrl = 'http://localhost:4201';
window.__env.mfeMdmUrl = 'http://localhost:4202';
```

## Backend Go Layer

```text
apps/backend-go/
├── go.work
├── iam-service/         # 8000
├── mdm-service/         # 8001
├── media-service/       # 8002
├── bpm-service/         # 8003
└── notification-service/# 8004
```

## Backend Java Layer

```text
apps/backend-java/
├── build.gradle.kts
├── settings.gradle.kts
├── crm-service/         # 8010
├── hrm-service/         # 8011
└── loan-service/        # 8012
```

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
