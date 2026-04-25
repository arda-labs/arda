# Arda Monorepo Guide

> Quy trình phát triển và vận hành hệ thống Arda Platform (Monorepo)
> Cập nhật: 2026-04-25

---

## 🏗️ Monorepo Structure

```text
arda/
├── apps/
│   ├── backend-go/      # Operational services (Kratos)
│   ├── backend-java/    # Core Banking services (Spring Boot)
│   └── frontend/        # Angular MFE (Nx Monorepo)
├── libs/
│   ├── go/              # Shared Go packages
│   └── java/            # Shared Java libraries
├── infra/               # Infrastructure & GitOps (arda-infra)
├── .github/             # Centralized CI/CD Workflows
└── docs/                # System Documentation
```

---

## 🛠️ Build & Development

### 🎨 Frontend (Angular MFE)

- **Cấu trúc**: `apps/frontend/`
- **Lệnh Build**: `npm run build` (trong apps/frontend) hoặc `nx build [app-name]`
- **Lệnh Test**: `nx test [app-name]`

### 🛠️ Backend Go (Operational)

- **Cấu trúc**: `apps/backend-go/`
- **Lệnh Build**: `make build` (trong từng service)
- **Lệnh Wire DI**: `wire ./...` (trong cmd)

### ☕ Backend Java (Core Banking)

- **Cấu trúc**: `apps/backend-java/`
- **Lệnh Build**: `./gradlew build` (từ root)
- **Lệnh chạy cụ thể**: `./gradlew :apps:backend-java:accounting:bootRun`

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
