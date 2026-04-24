# Skills Proposal — Các skill chuyên biệt cho Arda

> Đề xuất các skill chuyên biệt để hỗ trợ phát triển dự án Arda
> Cập nhật: 2026-04-24

---

## 📋 Overview

Dựa trên kiến trúc Arda (Frontend: Angular 21 + Nx, Backend: Go + Kratos + Java + Gradle, Infra: K3s + ArgoCD), tôi đề xuất các skill chuyên biệt theo từng nhóm nghiệp vụ.

---

## 🎨 Frontend Development Skills

### 1. Angular Development Skill
**Mục đích**: Hỗ trợ phát triển Angular applications, components, và MFEs

**Chức năng chính**:
- Tạo Angular component mới
- Sửa Angular component hiện có
- Tạo Angular service mới
- Sửa Angular service
- Sửa Angular module
- Setup routing
- Setup forms validation
- Implement Angular pipes
- Setup dependency injection

**Tech stack**: Angular 21, TypeScript 5.9

### 2. Nx Workspace Skill
**Mục đích**: Hỗ trợ phát triển trong Nx monorepo

**Chức năng chính**:
- Tạo ứng dụng/library mới trong workspace
- Build/tự động build affected apps
- Tạo shared library
- Xem dependency graph
- Xử lý NX cache
- Migrate ứng dụng sang workspace
- Setup Module Federation

**Tech stack**: Nx 22, Rspack, Analog.js

### 3. PrimeNG Components Skill
**Mục đích**: Hỗ trợ tạo và sửa PrimeNG components

**Chức năng chính**:
- Tạo PrimeNG component mới
- Sửa PrimeNG component hiện có
- Setup PrimeNG theming
- Configure PrimeNG dialogs/tables/trees
- Setup PrimeNG icons
- Implement custom PrimeNG templates

**Tech stack**: PrimeNG 21, Tailwind CSS

### 4. Module Federation Skill
**Mục đích**: Hỗ trợ cấu hình Module Federation cho MFEs

**Chức năng chính**:
- Setup webpack.config.js cho host app
- Setup webpack.config.js cho remote app
- Configure shared dependencies
- Troubleshoot Module Federation issues
- Setup lazy loading
- Configure micro-frontend routing

**Tech stack**: Angular 21, @angular-architects/module-federation

---

## 🛠️ Backend Go Skills

### 5. Kratos Framework Skill
**Mục đích**: Hỗ trợ phát triển microservices với Kratos framework

**Chức năng chính**:
- Tạo Kratos service mới
- Setup Kratos server (HTTP/gRPC)
- Tạo Kratos proto definitions
- Generate gRPC code từ proto
- Setup Kratos middleware
- Configure Kratos configuration
- Setup Wire dependency injection

**Tech stack**: Kratos v2.9, Go 1.25+, gRPC-go

### 6. gRPC & Protobuf Skill
**Mục đích**: Hỗ trợ tạo gRPC services và proto definitions

**Chức năng chính**:
- Tạo .proto file mới
- Generate gRPC Go code
- Generate gRPC TypeScript code
- Validate proto definitions
- Setup gRPC client
- Troubleshoot gRPC issues
- Implement gRPC interceptors

**Tech stack**: Protocol Buffers, gRPC-go, Protoc-gen-go, ts-proto

### 7. PostgreSQL with pgx Skill
**Mục đích**: Hỗ trợ làm việc với PostgreSQL database sử dụng pgx driver

**Chức năng chính**:
- Viết SQL queries với pgx
- Execute SQL statements
- Tạo database migration
- Setup database connection pool
- Implement transactions
- Optimize SQL queries
- Setup prepared statements

**Tech stack**: pgx v5, PostgreSQL 16+

### 8. Go Redis Integration Skill
**Mục đích**: Hỗ trợ tích hợp Redis vào Go services

**Chức năng chính**:
- Setup Redis connection
- Get/Set cache values
- Implement cache patterns
- Setup Redis pub/sub
- Implement distributed locking
- Optimize cache usage

**Tech stack**: go-redis v9, Redis

### 9. Go Middleware Skill
**Mục đích**: Hỗ trợ tạo middleware cho Kratos services

**Chức năng chính**:
- Tạo auth middleware (forward auth)
- Tạo logging middleware
- Tạo recovery middleware
- Tạo request ID middleware
- Tạo tenant middleware
- Setup middleware chain

**Tech stack**: Kratos middleware, Go 1.25+

---

## ☕ Backend Java Skills

### 10. Spring Boot Development Skill
**Mục đích**: Hỗ trợ phát triển Spring Boot applications

**Chức năng chính**:
- Tạo Spring Boot application mới
- Tạo controller mới
- Tạo service mới
- Tạo repository mới
- Setup Spring Data JPA
- Configure Spring Boot properties
- Setup dependency injection
- Implement REST endpoints

**Tech stack**: Spring Boot 3.x, Java 21, Spring Data JPA

### 11. Gradle Build System Skill
**Mục đích**: Hỗ trợ build và quản lý dependencies với Gradle multi-project

**Chức năng chính**:
- Tạo Gradle module mới
- Thêm dependency mới
- Update dependency version
- Xem dependency tree
- Resolve dependency conflicts
- Setup Gradle tasks
- Optimize build time
- Setup multi-project build

**Tech stack**: Gradle 8.x, Gradle wrapper

### 12. GraalVM Native Image Skill
**Mục đích**: Hỗ trợ biên dịch Java service sang GraalVM native image

**Chức năng chính**:
- Setup GraalVM CE 23+
- Build native image với Gradle
- Configure native image build
- Setup GraalVM tracing
- Fix native image build errors
- Optimize native image size
- Test native image

**Tech stack**: GraalVM CE 23+, Java 21, Gradle Native Image plugin

### 13. R2DBC PostgreSQL Skill
**Mục đích**: Hỗ trợ reactive database access với R2DBC và PostgreSQL

**Chức năng chính**:
- Setup R2DBC connection
- Implement reactive repositories
- Write reactive queries
- Setup connection pool
- Implement transaction handling
- Optimize reactive queries

**Tech stack**: R2DBC PostgreSQL, Spring Boot 3.x, Project Reactor

---

## 🌐 DevOps/Infrastructure Skills

### 14. K3s/Kubernetes Skill
**Mục đích**: Hỗ trợ deploy và quản lý applications trên K3s

**Chức năng chính**:
- Tạo Kubernetes deployment manifest
- Tạo Kubernetes service manifest
- Tạo Kubernetes config map
- Tạo Kubernetes secret
- Setup namespace
- Deploy application
- Scale application
- Rollback deployment

**Tech stack**: kubectl, K3s, YAML

### 15. ArgoCD Operations Skill
**Mục đích**: Hỗ trợ quản lý ArgoCD applications

**Chức năng chính**:
- Tạo ArgoCD application manifest
- Sync ArgoCD application
- Rollback application
- View ArgoCD application history
- View ArgoCD application resources
- Configure ArgoCD sync policy

**Tech stack**: ArgoCD, argocd CLI

### 16. Docker/Container Management Skill
**Mục đích**: Hỗ trợ build và quản lý Docker containers

**Chức năng chính**:
- Tạo Dockerfile cho service
- Build Docker image
- Push Docker image
- Optimize Dockerfile size
- Setup Docker compose
- Troubleshoot Docker issues

**Tech stack**: Docker, Docker Compose

### 17. APISIX Gateway Skill
**Mục đích**: Hỗ trợ cấu hình APISIX gateway

**Chức năng chính**:
- Tạo APISIX route
- Tạo APISIX upstream
- Cấu hình APISIX plugins
- Setup APISIX forward auth
- Setup APISIX rate limiting
- Setup APISIX circuit breaker

**Tech stack**: APISIX, APISIX Admin API

### 18. Zitadel/Auth0 Configuration Skill
**Mục đích**: Hỗ trợ cấu hình Zitadel identity provider

**Chức năng chính**:
- Tạo Zitadel project
- Tạo Zitadel application
- Cấu hình OIDC clients
- Tạo Zitadel users
- Cấu hình Zitadel roles
- Cấu hình Zitadel permissions

**Tech stack**: Zitadel, OIDC, JWT

### 19. Redpanda Skill
**Mục đích**: Hỗ trợ setup và quản lý Redpanda message broker

**Chức năng chính**:
- Tạo Redpanda topic
- Tạo Redpanda consumer
- Tạo Redpanda producer
- Setup Redpanda security
- Monitor Redpanda topics
- Troubleshoot Redpanda issues

**Tech stack**: Redpanda, rpk CLI

### 20. Camunda 7 Workflow Skill
**Mục đích**: Hỗ trợ tạo và quản lý Camunda workflows

**Chức năng chính**:
- Tạo BPMN 2.0 workflow
- Deploy workflow lên Camunda
- Bắt đầu workflow instance
- Complete user task
- Assign user task
- Xem workflow history

**Tech stack**: Camunda 7, BPMN 2.0

---

## 📊 Database Skills

### 21. PostgreSQL Skill
**Mục đích**: Hỗ trợ thiết kế và tối ưu hóa PostgreSQL database

**Chức năng chính**:
- Viết tối ưu SQL queries
- Thiết kế database schema
- Tạo database indexes
- Thiết kế database views
- Tạo database migrations
- Optimize query performance
- Setup PostgreSQL configuration

**Tech stack**: PostgreSQL 16+, SQL

### 22. Database Migration Skill
**Mục đích**: Hỗ trợ tạo và chạy database migrations

**Chức năng chính**:
- Tạo migration script (SQL)
- Tạo migration rollback script
- Version control migrations
- Setup migration tool (Flyway/Liquibase)
- Run migrations
- Validate migrations

**Tech stack**: Flyway/Liquibase, PostgreSQL

---

## 🧪 Testing Skills

### 23. Go Unit Testing Skill
**Mục đích**: Hỗ trợ viết unit tests cho Go code

**Chức năng chính**:
- Tạo test file mới
- Setup test dependencies
- Write table-driven tests
- Mock dependencies
- Run tests
- Generate test coverage report

**Tech stack**: testing, testify, mock

### 24. Java Unit Testing Skill
**Mục đích**: Hỗ trợ viết unit tests cho Java code

**Chức năng chính**:
- Tạo test case mới
- Setup test dependencies (JUnit 5, Mockito)
- Write unit tests
- Mock dependencies
- Run tests
- Generate test coverage report

**Tech stack**: JUnit 5, Mockito, Spring Boot Test

### 25. Integration Testing Skill
**Mục đích**: Hỗ trợ viết integration tests

**Chức năng chính**:
- Setup testcontainers
- Tạo test database seed
- Tạo test Redis seed
- Setup Mock servers
- Write integration tests
- Run integration tests

**Tech stack**: Testcontainers, WireMock, Testcontainers-go

### 26. End-to-End Testing Skill
**Mục đích**: Hỗ trợ viết E2E tests

**Chức năng chính**:
- Setup Playwright
- Tạo E2E test scenarios
- Run E2E tests
- Generate E2E test reports
- Record E2E test videos

**Tech stack**: Playwright, Cypress

---

## 📚 Documentation Skills

### 27. gRPC Documentation Skill
**Mục đích**: Hỗ trợ tạo documentation cho gRPC services

**Chức năng chính**:
- Tạo .proto file documentation
- Tạo API documentation từ proto
- Generate OpenAPI spec
- Validate API documentation
- Publish API documentation

**Tech stack**: Protocol Buffers, gRPC, OpenAPI

### 28. REST API Documentation Skill
**Mục đích**: Hỗ trợ tạo documentation cho REST APIs

**Chức năng chính**:
- Tạo API endpoint documentation
- Tạo request/response examples
- Tạo OpenAPI spec
- Validate API documentation
- Publish API documentation

**Tech stack**: OpenAPI 3.0, Swagger

### 29. README Generation Skill
**Mục đích**: Hỗ trợ tạo README files cho projects

**Chức năng chính**:
- Tạo README cho service mới
- Tạo README cho library mới
- Cập nhật README hiện có
- Standardize README format
- Tạo README từ code structure

**Tech stack**: Markdown, TypeScript

---

## 🔍 Code Quality Skills

### 30. Code Review Skill
**Mục đích**: Hỗ trợ review code và suggest improvements

**Chức năng chính**:
- Review code patterns
- Suggest best practices
- Identify code smells
- Suggest performance improvements
- Suggest security improvements
- Generate review report

**Tech stack**: Go, Java, TypeScript

### 31. Security Review Skill
**Mục đích**: Hỗ trợ security review cho code

**Chức năng chính**:
- Review code for security vulnerabilities
- Review code for OWASP Top 10
- Suggest security improvements
- Review authentication/authorization
- Review data handling
- Generate security report

**Tech stack**: Go, Java, TypeScript

---

## 🔄 Migration Skills

### 32. EPAS to Arda Mapping Skill
**Mục đích**: Hỗ trợ mapping EPAS features sang Arda architecture

**Chức năng chính**:
- Mapping EPAS table sang Arda schema
- Mapping EPAS API sang Arda gRPC API
- Mapping EPAS component sang Arda component
- Generate migration plan
- Tạo data mapping document

**Tech stack**: EPAS codebase, Arda architecture

### 33. Data Migration Skill
**Mục đích**: Hỗ trợ tạo và chạy data migration scripts

**Chức năng chính**:
- Tạo data migration script (SQL)
- Tạo data validation script
- Tạo rollback script
- Run data migration
- Validate migrated data
- Generate migration report

**Tech stack**: SQL, PostgreSQL

---

## 📋 Skills Recommendation Priority

### Phase 1 (Foundation) — Priority: P0

1. [Angular Development](#1-angular-development-skill) — Frontend base
2. [Nx Workspace](#2-nx-workspace-skill) — Frontend monorepo
3. [PostgreSQL](#21-postgresql-skill) — Database cơ bản
4. [Docker/Container Management](#16-dockercontainer-management-skill) — Container base
5. [K3s/Kubernetes](#14-k3skubernetes-skill) — K3s deployment

### Phase 2 (Core Banking) — Priority: P0

6. [Spring Boot Development](#10-spring-boot-development-skill) — Java backend
7. [GraalVM Native Image](#12-graalvm-native-image-skill) — Native compilation
8. [R2DBC PostgreSQL](#13-r2dbc-postgresql-skill) — Reactive DB
9. [Database Migration](#22-database-migration-skill) — Data migration
10. [gRPC Documentation](#27-grpc-documentation-skill) — API docs

### Phase 3 (Operational) — Priority: P1

11. [Kratos Framework](#5-kratos-framework-skill) — Go backend
12. [gRPC & Protobuf](#6-grpc-protobuf-skill) — gRPC
13. [PostgreSQL with pgx](#7-postgresql-with-pgx-skill) — Go DB driver
14. [Go Redis Integration](#8-go-redis-integration-skill) — Go cache
15. [Go Middleware](#9-go-middleware-skill) — Go middleware

### Phase 4 (Infrastructure) — Priority: P1

16. [ArgoCD Operations](#15-argocd-operations-skill) — GitOps
17. [APISIX Gateway](#17-apisix-gateway-skill) — API gateway
18. [Zitadel/Auth0 Configuration](#18-zitadel-auth0-configuration-skill) — Identity
19. [Redpanda](#19-redpanda-skill) — Event broker
20. [Camunda 7 Workflow](#20-camunda-7-workflow-skill) — BPM

### Phase 5 (Quality & Testing) — Priority: P2

21. [Go Unit Testing](#23-go-unit-testing-skill) — Go tests
22. [Java Unit Testing](#24-java-unit-testing-skill) — Java tests
23. [Integration Testing](#25-integration-testing-skill) — Integration tests
24. [End-to-End Testing](#26-e2e-testing-skill) — E2E tests
25. [Code Review](#30-code-review-skill) — Code review
26. [Security Review](#31-security-review-skill) — Security review

### Phase 6 (Documentation & Migration) — Priority: P2

27. [PrimeNG Components](#3-primeng-components-skill) — UI components
28. [Module Federation](#4-module-federation-skill) — MFE
29. [REST API Documentation](#28-rest-api-documentation-skill) — REST docs
30. [README Generation](#29-readme-generation-skill) — README
31. [EPAS to Arda Mapping](#32-epas-to-arda-mapping-skill) — Mapping
32. [Data Migration](#33-data-migration-skill) — Data migration

---

## 🤔 How to Use Skills

### Khi cần hỗ trợ chuyên biệt

Chỉ cần gọi skill theo tên, ví dụ:

- `/angular-dev` — Phát triển Angular
- `/nx-workspace` — Quản lý workspace
- `/spring-dev` — Phát triển Spring Boot
- `/kratos-dev` — Phát triển Kratos
- `/k3s-deploy` — Deploy lên K3s
- `/argocd-sync` — Sync ArgoCD
- `/api-doc` — Tạo API documentation
- `/code-review` — Review code
- `/security-review` — Security review

### Example Invocations

```
Tạo Angular component mới:
/angular-dev "Tạo component DataTable với các tính năng: sort, filter, pagination"

Build service Kratos:
/kratos-dev "Tạo service CRM service với các endpoint: CreateCustomer, GetCustomer, ListCustomers"

Deploy lên K3s:
/k3s-deploy "Tạo deployment manifest cho CRM service"

Review code:
/code-review "Review file be_loan/src/main/java/arda-labs/com/loan/feature/loan_parameter/interest_process/"
```

---

## 📊 Summary

- **Tổng số skills**: 33 skills
- **Phases**: 6 phases (Foundation → Core Banking → Operational → Infrastructure → Quality → Documentation/Migration)
- **Phân theo domain**: Frontend (4), Backend Go (5), Backend Java (4), DevOps (7), Database (2), Testing (4), Documentation (3), Quality (2), Migration (2)
- **Priority**: P0 (5), P1 (7), P2 (21)

---

Bạn có muốn tôi tạo skill nào trước? Hoặc bạn có đề xuất về skills này?

*Last Updated: 2026-04-24*
