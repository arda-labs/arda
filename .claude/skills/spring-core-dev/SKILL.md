---
name: spring-core-dev
description: Hỗ trợ phát triển Java Core Banking services với Spring Boot, R2DBC, và GraalVM Native Image
disable-model-invocation: false
---

# Spring Core Development Skill

Mục đích: Hỗ trợ phát triển Java microservices (Core Banking Services) trong thư mục `arda-core/`.

## 🎯 Phạm vi

- **Framework**: Spring Boot 3.x (Kotlin/Java)
- **Reactive DB**: R2DBC với PostgreSQL (Flyway migrations)
- **Native**: GraalVM Native Image (Optimized for RAM)
- **Build**: Gradle Multi-project
- **Patterns**: DDD, Reactive Programming (Reactor/Coroutines)

## 📦 Project Structure (Arda Core)

```
arda-core/
├── services/[service-name]/
│   └── src/main/kotlin/com.arda_labs.arda/[service]/
│       ├── controller/    # Reactive REST Controllers
│       ├── service/       # Business Logic
│       ├── repository/    # R2DBC Repositories
│       ├── domain/        # Entities & Value Objects
│       └── config/        # Spring/Native Configuration
└── libs/                  # Shared modules (common, database)
```

## 🛠️ Key Patterns

### 1. Reactive R2DBC
Sử dụng `Flux` và `Mono` cho không gian non-blocking.
```kotlin
@Repository
interface JournalRepository : R2dbcRepository<Journal, UUID> {
    @Query("SELECT * FROM journals WHERE tenant_id = :tenantId")
    fun findByTenantId(tenantId: String): Flux<Journal>
}
```

### 2. GraalVM Native Image Optimization
- Tránh reflection động nếu không cần thiết.
- Sử dụng `@RegisterReflectionForMapping` hoặc cấu hình `reflect-config.json`.
- Monitor RAM usage (Target: < 400MB per service).

### 3. Gradle Build
Sử dụng Kotlin DSL cho build scripts.

## 🎯 Usage
- `/spring-core-dev "Tạo reactive repository cho [Entity]"`
- `/spring-core-dev "Optimize service [name] cho GraalVM native build"`
- `/spring-core-dev "Implement transaction với R2DBC"`
