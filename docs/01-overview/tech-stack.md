# Tech Stack

Updated: 2026-04-30

This file separates technology currently present in the repo from platform
roadmap choices.

## Current Application Stack

### Frontend

| Component | Current technology |
| --- | --- |
| Framework | Angular 21 |
| Workspace | Angular CLI workspace |
| Federation | `@angular-architects/native-federation` |
| Bundler | Angular application builder / esbuild |
| UI | PrimeNG 21, PrimeIcons |
| Styling | Tailwind CSS 4, `tailwindcss-primeui` |
| Auth client | `angular-auth-oidc-client` |
| Tests | Angular unit-test builder with Vitest/jsdom |
| Package manager | npm 11 |

Current frontend projects: `shell`, `iam`, `mdm`, and `core`.

### Backend Go

| Component | Current technology |
| --- | --- |
| Language | Go 1.26.x |
| Framework | Kratos v2.9 |
| Workspace | `apps/backend-go/go.work` |
| Database driver | pgx v5 |
| Migration | `golang-migrate` with embedded migrations |
| Dependency injection | Google Wire |
| Transport | HTTP/JSON and gRPC generated from protobuf |

Current active services: `iam-service` and `mdm-service`.

### Backend Java

| Component | Current technology |
| --- | --- |
| Language | Kotlin / Java 21 target |
| Build | Gradle Kotlin DSL |
| Current module | `accounting_tmp` prototype |

The Java production stack is not finalized. Older docs mention Spring Boot,
R2DBC, and GraalVM Native Image as the target direction, not the current
production runtime.

### Runtime And Infra

| Component | Current technology |
| --- | --- |
| Kubernetes | K3s on `thinkcenter` |
| Gateway | Apache APISIX |
| Identity provider | Zitadel |
| GitOps | ArgoCD + Kustomize |
| Registry | GHCR |
| Public ingress | Cloudflared |
| Database | PostgreSQL on `thinkcenter` |
| Local gateway | Standalone APISIX under `arda-infra/local/apisix` |

## CI/CD

| Workflow | Purpose |
| --- | --- |
| `.github/workflows/ci-mfe.yml` | Detects and builds affected frontend apps: `shell`, `iam`, `mdm` |
| `.github/workflows/ci-go.yml` | Detects and builds Go services: `iam-service`, `mdm-service` |
| `.github/workflows/ci-java.yml` | Java pipeline placeholder; still expects future `accounting` module |
| `.github/workflows/gitops-update.yml` | Updates image tags in `arda-infra/apps/<service>/overlays/dev` |

## Roadmap Stack

These choices are platform direction, not all deployed today:

- Redpanda for event streaming.
- Outbox pattern for cross-service events.
- Camunda for workflow orchestration.
- Garage S3 for object storage.
- Prometheus, Grafana, Loki, and Alertmanager for observability.
- Java core services compiled to native images where the resource tradeoff is
  worth it.
