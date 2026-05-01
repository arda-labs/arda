# Tech Stack

Updated: 2026-05-02

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
| UI | PrimeNG 19+, PrimeIcons |
| Styling | Tailwind CSS 4, `tailwindcss-primeui` |
| Auth client | `angular-auth-oidc-client` |
| Tests | Angular unit-test builder with Vitest/jsdom |
| Package manager | npm 11 |

Current frontend projects: `shell`, `iam`, `mdm`, `crm`, `bpm`, `loan`, `hrm`, and `core`.

### Backend Go

| Component | Current technology |
| --- | --- |
| Language | Go 1.26.x |
| Framework | Kratos v2.8+ |
| Workspace | `apps/backend-go/go.work` |
| Database driver | pgx v5 |
| Migration | `golang-migrate` |
| Dependency injection | Google Wire |
| Transport | HTTP/JSON and gRPC |
| Workflow | Zeebe Go Client (Camunda 8) |
| Messaging | segmentio/kafka-go |

Current active services: `iam-service`, `mdm-service`, and `bpm-service`.

### Backend Java

| Component | Current technology |
| --- | --- |
| Language | Java 25 (LTS) |
| Framework | Spring Boot 4.x / WebFlux |
| Build | Gradle Kotlin DSL |
| Messaging | Spring Kafka (Reactor) |
| gRPC | net.devh:grpc-server-spring-boot-starter |
| Workflow | Spring Zeebe Starter |

Current modules: `crm-service`, `loan-service`, `hrm-service`.

### Runtime And Infra

| Component | Current technology |
| --- | --- |
| Kubernetes | K3s |
| Gateway | Apache APISIX |
| Identity provider | Zitadel / Keycloak |
| Event Bus | Apache Kafka / Redpanda |
| Workflow Engine | Camunda 8 (Zeebe) |
| Database | PostgreSQL |

## Roadmap Stack

- Redpanda for distributed event streaming.
- Outbox pattern for transactional consistency.
- SeaweedFS / Garage for S3-compatible storage.
- Prometheus, Grafana, Loki for full observability.
- GraalVM Native Image for resource-constrained Java modules.
