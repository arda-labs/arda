# Architecture — Tổng quan Kiến trúc

> Tổng quan kiến trúc hệ thống Arda Platform
> Cập nhật: 2026-04-24

---

## 📋 Overview

Arda là nền tảng microservices tài chính với kiến trúc 3-layer:
1. **Frontend Layer** — Angular 21 + Nx + MFE
2. **Backend Layer** — Go (Kratos) + Java (Spring Boot + Gradle)
3. **Infrastructure Layer** — K3s + APISIX + Zitadel + PostgreSQL

---

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Shell   │  │  Common  │  │  MFE Apps│  │  Admin   │    │
│  │  (Host)  │  │  (Lib)   │  │ (Domain) │  │  (Panel) │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                  API Gateway Layer                      │
│                Apache APISIX (Self-hosted)           │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                 Identity & Access Layer                  │
│                  Zitadel (Self-hosted)                │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                        │
│  ┌──────────────────┐  ┌──────────────────┐           │
│  │  Core Banking    │  │   Operational    │           │
│  │  (Java + Gradle) │  │   (Go + Kratos)  │           │
│  │                  │  │                  │           │
│  │  • Accounting    │  │  • CRM & Member  │           │
│  │  • Loan          │  │  • HRM           │           │
│  │  • Deposit       │  │  • Notification  │           │
│  │                  │  │  • System Config │           │
│  │                  │  │  • BPM Engine    │           │
│  └──────────────────┘  └──────────────────┘           │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌──────────────────────┼──────────────────────┐
        ▼                      ▼                      ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│   Messaging   │    │   Process     │    │   Storage     │
│   Redpanda    │    │   Camunda 7   │    │   Garage S3   │
└───────────────┘    └───────────────┘    └───────────────┘
        │                      │                      │
        └──────────────────────┼──────────────────────┘
                               ▼
                    ┌──────────────┐
                    │   Data Layer     │
                    │   PostgreSQL 16+ │
                    │   + Redis        │
                    └───────────────┘
```

---

## 🏗️ Frontend Architecture

### Monorepo Structure (Nx)

```
arda-mfe/
├── apps/
│   ├── shell/                 # Host app
│   ├── common/                # Common MFE
│   ├── accounting/            # Accounting MFE
│   ├── loan/                  # Loan MFE
│   ├── crm/                   # CRM MFE
│   ├── hrm/                   # HRM MFE
│   ├── bpm/                   # BPM MFE
│   └── admin/                 # Admin MFE
│
└── libs/
    ├── ui/                    # Shared UI components
    ├── auth/                  # Auth utilities
    └── shared/                # Shared utilities
```

### Technology Stack

- **Framework**: Angular 21
- **Monorepo**: Nx 22
- **Build**: Rspack 1.6
- **SSR**: Analog.js 2.x
- **UI**: PrimeNG 21
- **Auth**: angular-auth-oidc-client

---

## 🛠️ Backend Go Architecture (Kratos)

### Monorepo Structure

```
arda-be/
├── api/                     # Shared .proto definitions (Source of Truth)
│   ├── crm/v1/
│   ├── hrm/v1/
│   ├── notification/v1/
│   ├── system-config/v1/
│   └── bpm/v1/
│
├── pkg/                     # Shared libraries
│   ├── auth/                  # Zitadel integration
│   ├── database/              # PostgreSQL helpers
│   ├── redis/                 # Redis client wrapper
│   └── middleware/            # gRPC/HTTP middleware
│
├── crm-service/             # CRM & Member Service
├── hrm-service/             # HRM Service
├── notification-service/      # Notification Service
├── system-config-service/     # System Config Service
└── bpm-service/             # BPM Engine Wrapper
```

### Technology Stack

- **Language**: Go 1.25+
- **Framework**: Kratos v2.9
- **DB Driver**: pgx v5
- **Cache**: go-redis v9
- **gRPC**: grpc-go v1.74
- **DI**: Wire v0.6

---

## ☕ Backend Java Architecture (Spring Boot + Gradle)

### Multi-project Structure

```
arda-core/
├── services/
│   ├── accounting/            # Accounting Service
│   ├── loan/                  # Loan Service
│   ├── deposit/               # Deposit Service
│   └── treasury/              # Treasury Service
│
└── libs/
    ├── shared-core/           # Core utilities
    ├── grpc-client/           # gRPC clients
    ├── security/              # Security context
    └── accounting-core/       # Accounting domain logic
```

### Technology Stack

- **Language**: Java 21
- **Framework**: Spring Boot 3.x
- **Build**: Gradle 8.x (Multi-project)
- **Native**: GraalVM CE 23+
- **DB Driver**: r2dbc-postgresql
- **Monitoring**: Micrometer + Prometheus

---

## 🌐 API Gateway Architecture

### APISIX Configuration

```
arda.io.vn          → APISIX Gateway
├── /api/v1/iam/*         → iam-service
├── /api/v1/crm/*         → crm-service
├── /api/v1/hrm/*         → hrm-service
├── /api/v1/accounting/*   → accounting-service
├── /api/v1/loan/*         → loan-service
├── /api/v1/bpm/*         → bpm-service
├── /api/v1/config/*       → system-config-service
├── /api/v1/notification/* → notification-service
├── /common/*              → mfe-common
├── /loan/*                → mfe-loan
├── /crm/*                 → mfe-crm
├── /accounting/*           → mfe-accounting
├── /bpm/*                 → mfe-bpm
├── /hrm/*                 → mfe-hrm
└── /*                     → mfe-shell
```

---

## 🔐 Identity & Access Architecture

### Zitadel Integration

```
User → Zitadel (auth.arda.io.vn)
     ↓
  JWT Token (access_token)
     ↓
APISIX → iam-service (Forward Auth)
     ↓
  Permission Check
     ↓
  Allow/Deny Request
```

### Permission Model

- **RBAC**: Role-based access
- **ABAC**: Attribute-based access
- **ReBAC**: Relationship-based access
- **Maker-Checker**: Tách biệt người lập và người duyệt

---

## 📊 Data Architecture

### Database Schema

```
PostgreSQL Schemas:
├── arda_iam          # IAM service
├── arda_crm          # CRM service
├── arda_hrm          # HRM service
├── arda_accounting   # Accounting service
├── arda_loan         # Loan service
└── zitadel           # Zitadel (separate)
```

### Multi-tenancy

- **Column-based**: `tenant_id` column ở tất cả tables
- **RLS**: Row-Level Security cho isolation
- **Tenant Context**: Propagated qua gRPC metadata

---

## 🔄 Event-Driven Architecture

### Message Flow

```
Service → Outbox Table (Local Transaction)
     ↓
Message Relay Worker
     ↓
Redpanda (Topic: arda.<service>.events)
     ↓
Consumer Service (Idempotent)
```

### Outbox Pattern

1. Service ghi event vào `outbox` table
2. Local transaction đảm bảo consistency
3. Message Relay worker đọc `outbox` và đẩy vào Redpanda
4. Consumer xử lý event

---

## 📈 Scalability Strategy

### Horizontal Scaling

- **Frontend**: Kubernetes Horizontal Pod Autoscaler (HPA)
- **Backend Go**: HPA dựa trên CPU/Memory
- **Backend Java**: HPA dựa trên CPU/Memory
- **Database**: Connection pooling (PgBouncer)

### Resource Allocation (32GB RAM Target)

```
Infrastructure:     10GB  (PostgreSQL 4GB, Redis 1GB, Redpanda 2GB, Camunda 1GB, S3 0.5GB, System 1.5GB)
Core Services:      2GB   (4 services × ~400MB + overhead)
Operational:        1GB   (5 services × ~100MB + overhead)
Frontend:           1GB   (Nx build cache, dev server)
K3s/OS:             3GB   (K8s overhead, system)
Reserved:           15GB  (Future services, buffer)
```

---

## 📚 Related Documentation

- [Tech Stack](tech-stack.md) — Công nghệ chi tiết
- [Frontend Architecture](../frontend/architecture.md) — Chi tiết frontend
- [Go Architecture](../backend-go/architecture.md) — Chi tiết Go backend
- [Java Architecture](../backend-java/architecture.md) — Chi tiết Java backend
- [Infrastructure Architecture](../../INFRA_STATUS.md) — Chi tiết hạ tầng thực tế

---

*Last Updated: 2026-04-24*
