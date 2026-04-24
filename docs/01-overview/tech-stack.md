# Tech Stack — Công nghệ Sử dụng

> Danh sách công nghệ sử dụng trong Arda Platform
> Cập nhật: 2026-04-24

---

## 📋 Overview

Arda sử dụng công nghệ mới nhất, tối ưu cho hiệu năng và tiết kiệm tài nguyên (32GB RAM).

---

## 🏗️ Infrastructure

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| OS | Ubuntu | 24.04 LTS | Base OS |
| Container | K3s | Latest | Lightweight K8s |
| Gateway | Apache APISIX | Latest | API Gateway |
| Identity | Zitadel | Latest | IdP/OAuth2 |
| Database | PostgreSQL | 16+ | Primary DB |
| Cache | Redis | Latest | Session cache |
| Messaging | Redpanda | Latest | Event broker |
| BPM | Camunda | 7.x | Workflow engine |
| Storage | Garage S3 | Latest | Object storage |
| Ingress | Cloudflared | Latest | Tunnel to Cloudflare |
| GitOps | ArgoCD | Latest | CD automation |

---

## 🎨 Frontend

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | TypeScript | 5.9 | Type-safe JS |
| Framework | Angular | 21.x | UI framework |
| Monorepo | Nx | 22.x | Monorepo toolkit |
| Build | Rspack | 1.6 | Fast bundler |
| SSR | Analog.js | 2.x | Angular SSR |
| UI | PrimeNG | 21.x | UI component library |
| Auth | angular-auth-oidc-client | 21.x | OIDC client |

---

## 🛠️ Backend Go (Kratos)

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.25+ | Runtime |
| Framework | Kratos | v2.9 | Microservice framework |
| DB Driver | pgx | v5 | PostgreSQL driver |
| Cache | go-redis | v9 | Redis client |
| gRPC | grpc-go | v1.74 | RPC framework |
| DI | Wire | v0.6 | Dependency injection |

---

## ☕ Backend Java (Spring Boot + Gradle)

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Java | 21 | Runtime |
| Framework | Spring Boot | 3.x | Application framework |
| Build | Gradle | 8.x | Build tool |
| Multi-project | Gradle | 8.x | Multi-project structure |
| Native | GraalVM CE | 23+ | Native image compilation |
| DB Driver | r2dbc-postgresql | Latest | Reactive DB driver |

---

## 🔒 Security & Auth

| Component | Technology | Purpose |
|-----------|-----------|---------|
| IdP | Zitadel | OAuth2/OIDC provider |
| JWT | Zitadel JWT | Access tokens |
| Encryption | TLS 1.3 | Transport layer |
| RBAC/ABAC/ReBAC | Custom impl | Fine-grained access control |
| Maker-Checker | Custom impl | Approval workflows |

---

## 📊 Monitoring & Observability

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Metrics | Prometheus | Metrics collection |
| Tracing | OpenTelemetry | Distributed tracing |
| Logging | Structured logging | Log aggregation |
| Dashboards | Grafana | Visualization |
| Alerts | Alertmanager | Alerting |

---

## 🔄 Event-Driven

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Message Broker | Redpanda | Event streaming |
| Outbox Pattern | Custom impl | Transactional events |
| Idempotency | Redis | Idempotent consumers |

---

## 🧪 Testing

| Component | Technology | Purpose |
|-----------|-----------|---------|
| Unit Tests | JUnit / Go test | Unit testing |
| Integration Tests | Testcontainers | Integration testing |
| E2E Tests | Playwright / Cypress | End-to-end testing |

---

## 📦 CI/CD

| Component | Technology | Purpose |
|-----------|-----------|---------|
| CI | GitHub Actions | Continuous integration |
| CD | ArgoCD | Continuous deployment |
| Container Registry | GitHub Container Registry | Docker image storage |

---

*Last Updated: 2026-04-24*
