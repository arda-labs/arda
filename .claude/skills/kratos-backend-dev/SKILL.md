---
name: kratos-backend-dev
description: Hỗ trợ phát triển Go microservices với Kratos framework, pgx, Redis, và gRPC cho dự án Arda
disable-model-invocation: false
---

# Kratos Backend Development Skill

Mục đích: Hỗ trợ phát triển Go microservices (Operational Services) trong thư mục `arda-be/`.

## 🎯 Phạm vi

- **Framework**: Kratos v2.9 (HTTP/gRPC, Middleware, Wire DI)
- **Database**: PostgreSQL với pgx v5 (Pool, Transactions, RLS)
- **Cache**: Redis với go-redis v9
- **Communication**: gRPC & Protobuf definitions
- **Middleware**: Auth, Tenant (RLS), Logging, Recovery
- **Testing**: Unit & Integration testing

## 📦 Project Structure (Arda BE)

```
arda-be/[service-name]/
├── api/v1/                # Proto definitions
├── cmd/[service-name]/    # Entry point & Wire injection
├── internal/
│   ├── biz/               # Business Logic (Domain)
│   ├── data/              # Data Access (Repo, DB, Redis)
│   ├── service/           # API Implementation (DTO mapping)
│   └── conf/              # Config types (from proto)
├── configs/               # YAML configuration
└── Makefile               # Build & Gen commands
```

## 🛠️ Key Patterns

### 1. Database with pgx v5
Ưu tiên sử dụng `pgxpool` và `Row-Level Security (RLS)`.
```go
// Example Find with RLS context
func (r *repo) FindByID(ctx context.Context, id string) (*biz.Entity, error) {
    tenantID := middleware.GetTenantID(ctx)
    // RLS context is usually handled by middleware or explicitly set in tx
    query := "SELECT id, name FROM table WHERE id = $1"
    // ... logic
}
```

### 2. Kratos Middleware
Sử dụng cho Multi-tenancy và Forward Auth.
- `Tenant Middleware`: Trích xuất `X-Tenant-ID` vào context.
- `Auth Middleware`: Forward auth tới `iam-service`.

### 3. Wire Dependency Injection
```go
// cmd/service/wire.go
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
```

## 🎯 Usage
- `/kratos-backend-dev "Tạo service mới [name]"`
- `/kratos-backend-dev "Implement CRUD cho domain [entity] với pgx"`
- `/kratos-backend-dev "Setup gRPC middleware cho auth và multi-tenancy"`
