# Arda — Modern Financial Microservice Platform

> Platform Microservices thế hệ mới chuyên biệt cho lĩnh vực Tài chính & Ngân hàng
> Design: Domain-Driven Design (DDD), Event-Driven Architecture, Resource-Optimized
> Infrastructure: K3s on Ubuntu 24.04 (Current: 16GB RAM → Target: 32GB RAM)

---

## 📋 Overview

**Arda** là nền tảng microservices hiện đại được thiết kế để vận hành hệ thống tài chính quy mô lớn trên hạ tầng tài nguyên giới hạn. Hệ thống áp dụng các công nghệ tiên tiến như **GraalVM Native Image**, **Redpanda** và **Camunda 7** để tối ưu hóa hiệu suất và giảm tiêu thụ RAM.

### Key Characteristics
- **Resource-Constrained Architecture**: Tối ưu cho 32GB RAM với headroom 40%
- **Domain-Driven Design**: Services được phân rã theo domain nghiệp vụ
- **Event-Driven**: Saga pattern với Outbox để đảm bảo consistency
- **Multi-Tenancy**: Column-based với PostgreSQL Row-Level Security
- **Fine-Grained Access Control**: RBAC/ABAC/ReBAC với Maker-Checker pattern

---

## 🗂️ Monorepo Structure

```
arda/
├── docs/                          # Documentation (reorganized)
│   ├── overview/            # Tổng quan kiến trúc & tech stack
│   ├── features/            # Danh sách chức năng
│   ├── frontend/            # Frontend architecture
│   ├── backend-go/          # Go backend (Kratos)
│   ├── backend-java/        # Java backend (Spring Boot + Gradle)
│   ├── guides/             # Hướng dẫn setup
│   └── migration/          # Migration từ EPAS
│
├── arda-infra/                    # Infrastructure as Code
│   ├── apps/                      # ArgoCD applications
│   ├── bootstrap/                 # ArgoCD bootstrap
│   └── infrastructure/            # K3s setup
│
├── arda-be/                       # Go Backend (Operational Services)
│   ├── api/                       # Shared .proto definitions
│   ├── pkg/                       # Shared libraries
│   ├── iam-service/               # Identity & Access Management
│   ├── crm-service/               # CRM & Member Management
│   ├── hrm-service/               # Human Resource Management
│   ├── notification-service/      # Notification Service
│   ├── system-config-service/     # System Configuration
│   └── bpm-service/               # BPM Engine Wrapper
│
├── arda-core/                     # Java Backend (Core Banking Services)
│   ├── services/                  # Core services
│   │   ├── accounting/            # Accounting Service
│   │   ├── loan/                  # Loan Service
│   │   ├── deposit/               # Deposit Service
│   │   └── treasury/              # Treasury Service
│   ├── libs/                      # Shared libraries
│   ├── build.gradle.kts           # Gradle build config
│   └── settings.gradle.kts        # Gradle settings
│
├── arda-mfe/                      # Frontend (Angular Nx Monorepo)
│   ├── apps/
│   │   ├── shell/                 # Host application
│   │   ├── common/                # Shared MFE
│   │   ├── accounting/            # Accounting MFE
│   │   ├── loan/                  # Loan MFE
│   │   ├── crm/                   # CRM MFE
│   │   ├── hrm/                   # HRM MFE
│   │   └── admin/                 # Admin Panel MFE
│   ├── libs/
│   │   ├── ui/                    # Shared UI components
│   │   ├── auth/                  # Auth utilities
│   │   └── shared/                # Shared utilities
│   ├── nx.json                    # Nx configuration
│   └── package.json               # Node dependencies
│
└── INFRA_STATUS.md              # Thực trạng deployment (thinkcenter)
```

---

## 🚀 Quick Start

### Prerequisites
- Ubuntu 24.04 LTS
- 32GB RAM minimum
- K3s installed
- kubectl configured
- Helm 3.x
- Go 1.25+
- Java 21 + GraalVM CE 23+
- Node.js 20+
- Docker

### Installation

```bash
# 1. Clone repository
git clone https://github.com.arda_labs/arda.git
cd arda

# 2. Deploy infrastructure
cd arda-infra
kubectl apply -f infrastructure/

# 3. Build and deploy services
./scripts/deploy-all.sh

# 4. Access applications
# Frontend: https://arda.io.vn
# API Gateway: https://api.arda.io.vn
# Zitadel: https://auth.arda.io.vn
```

---

## 📊 Technology Stack

### Infrastructure
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

### Backend (Go)
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.25+ | Runtime |
| Framework | Kratos | v2.9 | Microservice framework |
| DB Driver | pgx | v5 | PostgreSQL driver |
| Cache | go-redis | v9 | Redis client |
| gRPC | grpc-go | v1.74 | RPC framework |
| DI | Wire | v0.6 | Dependency injection |

### Backend (Java)
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Java | 21 | Runtime |
| Framework | Spring Boot | 3.x | Application framework |
| Build | Gradle | 8.x | Build tool |
| Native | GraalVM CE | 23+ | Native image compilation |
| DB Driver | r2dbc-postgresql | Latest | Reactive DB driver |

### Frontend
| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | TypeScript | 5.9 | Type-safe JS |
| Framework | Angular | 21.x | UI framework |
| Build | Nx | 22.x | Monorepo toolkit |
| Bundle | Rspack | 1.6 | Fast bundler |
| SSR | Analog.js | 2.x | Angular SSR |
| UI | PrimeNG | 21.x | UI component library |
| Auth | angular-auth-oidc-client | 21.x | OIDC client |

---

## 🎯 Resource Allocation (32GB RAM)

```
Infrastructure:     10GB  (PostgreSQL 4GB, Redis 1GB, Redpanda 2GB, Camunda 1GB, S3 0.5GB, System 1.5GB)
Core Services:      2GB   (4 services × ~400MB + overhead)
Operational:        1GB   (5 services × ~100MB + overhead)
Frontend:           1GB   (Nx build cache, dev server)
K3s/OS:             3GB   (K8s overhead, system)
Reserved:           15GB  (Future services, buffer)
```

---

## 📖 Documentation

- [Documentation Hub](./docs/README.md) — Trung tâm tài liệu
- [Architecture](./docs/01-overview/architecture.md) — Tổng quan kiến trúc
- [Tech Stack](./docs/01-overview/tech-stack.md) — Công nghệ sử dụng
- [Infrastructure Status](./docs/01-overview/infra-status.md) — Thực trạng deployment
- [EPAS Migration Analysis](./docs/07-migration/epas-to-arda.md) — Migration từ EPAS

### Features Documentation
- [Accounting Features](./docs/features/accounting.md) — Chức năng Kế toán
- [Loan Features](./docs/features/loan.md) — Chức năng Cho vay
- [CRM Features](./docs/features/crm.md) — Chức năng Khách hàng
- [HRM Features](./docs/features/hrm.md) — Chức năng Nhân sự
- [BPM Features](./docs/features/bpm.md) — Chức năng Quy trình
- [Notification Features](./docs/features/notification.md) — Chức năng Thông báo
- [System Config Features](./docs/features/system-config.md) — Chức năng Cấu hình

### Technical Documentation
- [Frontend Architecture](./docs/frontend/architecture.md) — Chi tiết frontend
- [Go Architecture](./docs/backend-go/architecture.md) — Chi tiết Go backend
- [Java Architecture](./docs/backend-java/architecture.md) — Chi tiết Java backend
- [Infrastructure](./docs/infrastructure/README.md) — Chi tiết hạ tầng
- [Authorization](./docs/services/authorization.md) — Kiến trúc phân quyền
- [Deployment Guide](./docs/deployment/README.md) — Hướng dẫn triển khai

---

## 📜 License

Copyright © 2026 Arda Labs. All rights reserved.

---

*Last Updated: 2026-04-24*
