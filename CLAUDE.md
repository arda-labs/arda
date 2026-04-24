# Arda Monorepo Guide

> Quy trình phát triển và vận hành hệ thống Arda Platform
> Cập nhật: 2026-04-25

---

## 🛠️ Build & Development

### 🎨 Frontend (Angular MFE)
- **Cấu trúc**: `arda-mfe/` (Nx Monorepo)
- **Lệnh Build**: `npm run build` (trong arda-mfe) hoặc `nx build [app-name]`
- **Lệnh Test**: `nx test [app-name]`
- **Lệnh Lint**: `nx lint [app-name]`
- **Cấu hình**: `nx.json`, `package.json`

### 🛠️ Backend Go (Operational)
- **Cấu trúc**: `arda-be-go/` (Kratos services)
- **Lệnh Build**: `make build` (trong từng service)
- **Lệnh Generate API**: `make api`
- **Lệnh Wire DI**: `make generate` hoặc `wire ./...`
- **Chạy local**: `./bin/[service] -conf ./configs`

### ☕ Backend Java (Core Banking)
- **Cấu trúc**: `arda-be-java/` (Spring Boot + Gradle)
- **Lệnh Build**: `./gradlew build`
- **Lệnh Native Build**: `./gradlew nativeCompile`
- **Chạy local**: `./gradlew bootRun`

---

## 🏗️ Infrastructure & Deployment

### 🌐 Kubenetes (K3s)
- **Namespace**: `arda-dev` (Development), `arda-prod` (Production)
- **Context**: `thinkcenter` (thinkcenter)

### 🚀 GitOps (ArgoCD)
- **Repo**: `github.com.arda_labs/arda-infra`
- **Sync**: Tự động đồng bộ từ branch `main` của infra repo.

### 🌐 API Gateway (APISIX)

- **Domain chính**: `arda.io.vn`
- **Identity (Zitadel)**: `auth.arda.io.vn`
- **Admin API**: `10.43.150.100:9180`
- **Network Flow**: `User -> Cloudflare Tunnel -> APISIX -> Internal Services`

---

## 🌐 Network Architecture

| Hostname | Target Service | Path | Gateway |
| --- | --- | --- | --- |
| `arda.io.vn` | `mfe-shell` | `/*` | APISIX |
| `auth.arda.io.vn` | `zitadel` | `/*` | Traefik |
| `argocd.arda.io.vn` | `argocd-server` | `/*` | Cloudflare (Direct) |

**Lưu ý**: `auth.arda.io.vn` được tách riêng qua Traefik để đảm bảo tính độc lập cho hệ thống Identity, không đi qua APISIX Gateway.

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
kubectl logs -f deployment/iam-service -n arda-dev

# Kiểm tra resource usage
kubectl top pods -n arda-dev

# Sync ArgoCD app thủ công
argocd app sync iam-service-dev
```
