# Arda Monorepo Guide

> Quy trình phát triển và vận hành hệ thống Arda Platform (Monorepo)
> Cập nhật: 2026-04-25

---

## 🏗️ Monorepo Structure

```text
arda/
├── apps/
│   ├── backend-go/      # Operational services (Kratos Monorepo with Go Workspace)
│   │   ├── go.work      # Go workspace configuration
│   │   ├── common-service/
│   │   ├── crm-service/
│   │   └── iam-service/
│   ├── backend-java/    # Core Banking services (Gradle Multi-project Monorepo)
│   │   ├── build.gradle.kts
│   │   ├── settings.gradle.kts
│   │   └── accounting/
│   └── frontend/        # Angular MFE (Nx Monorepo)
├── libs/
│   ├── go/              # Shared Go packages
│   │   └── pkg/         # Common Go libraries (database, redis, middleware)
│   └── java/            # Shared Java libraries
│       ├── common/
│       ├── database/
│       └── ...
├── infra/               # Infrastructure & GitOps (arda-infra)
├── .github/             # Centralized CI/CD Workflows
├── docs/                # System Documentation
└── scripts/             # Utility scripts & Dev configs
```

---

## 🛠️ Build & Development

### 🎨 Frontend (Angular MFE)

- **Cấu trúc**: `apps/frontend-micro/`
- **Lệnh Build**: `npm run build` (trong apps/frontend-micro) hoặc `nx build [app-name]`

### 🛠️ Backend Go (Operational)

- **Cấu trúc**: `apps/backend-go/`
- **Quản lý Workspace**: `go.work` gom tất cả services và `libs/go/pkg`.
- **Lệnh Run**: `cd apps/backend-go/[service] && kratos run` hoặc `make run`.
- **Lệnh Wire DI**: `wire ./...` (trong thư mục cmd của service).

### ☕ Backend Java (Core Banking)

- **Cấu trúc**: `apps/backend-java/`
- **Quản lý Build**: Gradle monorepo độc lập trong `apps/backend-java`.
- **Lệnh Build**: `./gradlew build` (từ thư mục apps/backend-java).
- **Lệnh chạy cụ thể**: `./gradlew :accounting:bootRun`

---

## 🏗️ Infrastructure & Deployment

### 🌐 Kubenetes (K3s)

- **Namespace**: `arda-dev` (Development), `arda-prod` (Production)
- **Context**: `thinkcenter`

### 🚀 GitOps (ArgoCD)

- **Repo**: `github.com/arda-labs/arda-infra`
- **Sync**: Tự động đồng bộ từ branch `main` của infra repo.

---

## 🛡️ Coding Standards

- **Go**: Follow [Uber Go Style Guide](https://github.com/uber-go/guide).
- **Java/Kotlin**: Follow [Google Java Style](https://google.github.io/styleguide/javaguide.html).
- **Frontend**: Follow [Angular Style Guide](https://angular.io/guide/styleguide).
- **Git**: [Conventional Commits](https://www.conventionalcommits.org/).

---

## 📖 Useful Commands

```bash
# Xem logs service trong K8s
kubectl logs -f deployment/accounting-service -n arda-dev

# Kiểm tra Redpanda (Kafka) trong cụm
kubectl get pods -n arda-dev -l app=redpanda
```
