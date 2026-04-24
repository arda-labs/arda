# Arda Labs

> **Where data flows with intelligence.**

**Arda Labs** xây dựng hệ thống quản trị doanh nghiệp multi-tenant (SaaS) trên nền tảng cloud-native:

- **Frontend Monorepo** (Angular 21 + Nx 22 + Module Federation + Rspack)
- **Backend Monorepo Go** (Kratos v2.9 + PostgreSQL 16 RLS + pgx v5)
- **Backend Monorepo Java** (Spring Boot 3.5 + R2DBC + GraalVM Native Image)
- **GitOps Infrastructure** (K3s + ArgoCD + APISIX Gateway + Zitadel IAM)

---

## Brand

| Thuộc tính      | Giá trị                                               |
| --------------- | ----------------------------------------------------- |
| **Name**        | ARDA Labs                                             |
| **Fonts**       | Space Grotesk (logo) · Inter (text)                   |
| **Colors**      | Cyan `#06B6D4` · Indigo `#6366F1` · Midnight `#0F172A` |
| **Style**       | Minimal · Futuristic · Data-oriented                  |
| **Slogan**      | _"Where data flows with intelligence."_               |

---

# Arda Platform Architecture

```
Browser → Cloudflare Tunnel → APISIX Gateway (Forward Auth) → [K8s Services]
                                       ├─ mfe-shell    (Angular Host)
                                       ├─ zitadel      (Identity Provider)
                                       ├─ arda-be-go   (Go Microservices)
                                       └─ arda-be-java (Java Core Banking)
```

## CI/CD Flow

```
git push (main) → GitHub Actions
  ├─ lint · test · build (Nx/Go/Gradle)
  ├─ Docker build + push (GHCR)
  └─ Update tag in arda-infra → ArgoCD auto-sync → K3s deploys
```

## Tech Stack

| Layer        | Technology                                                      |
| :----------- | :-------------------------------------------------------------- |
| **Frontend** | Angular 21 · Nx v22 · Module Federation · Rspack · Tailwind CSS |
| **Go BE**    | Kratos v2.9 · Go 1.25 Workspace · PostgreSQL 16 RLS (pgx v5)    |
| **Java BE**  | Spring Boot 3.5 · R2DBC · GraalVM Native Image · Java 21        |
| **Security** | APISIX (Forward Auth) · Zitadel (OIDC/SAML) · Cloudflare ZT     |
| **Infra**    | K3s · ArgoCD · Kustomize · APISIX · Cloudflare Tunnel           |

## Repositories — [arda-labs](https://github.com.arda_labss)

### [arda](https://github.com.arda_labss/arda)

The **Master Monorepo** containing all business logic, services, and shared libraries.

- `arda-mfe/`: Nx monorepo for Frontend Micro-frontends.
- `arda-be-go/`: Go workspace for Kratos microservices.
- `arda-be-java/`: Gradle multi-project for Spring Boot core services.

### [arda-infra](https://github.com.arda_labss/arda-infra)

GitOps-managed Kubernetes manifests and environment configurations.

### [arda-docs](https://github.com.arda_labss/arda-docs)

System architecture, development guides, and business requirement documentation.

---

## Quick Start

```bash
# 1. Clone the master monorepo
git clone https://github.com.arda_labss/arda
cd arda

# 2. Run Go services locally
cd arda-be-go
go work sync
go run services/iam-service/cmd/iam-service

# 3. Start Frontend
cd arda-mfe
npx nx serve shell
```

---

- Website: [arda.io.vn](https://arda.io.vn)
- Contact: **contact@arda.io.vn**
- Organization: [github.com.arda_labss](https://github.com.arda_labss)
