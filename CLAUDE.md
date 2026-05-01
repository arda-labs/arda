# Arda Monorepo Guide

Updated: 2026-05-02

This file is a quick orientation guide for contributors and coding agents.
The source of truth for runtime manifests is the sibling repo `arda-infra`.

## Repository Layout

```text
arda/
├── apps/
│   ├── frontend-micro/
│   │   ├── angular.json
│   │   └── projects/
│   │       ├── shell/           # Host (4200)
│   │       ├── iam/             # Remote (4201)
│   │       ├── mdm/             # Remote (4202)
│   │       ├── crm/             # Remote (4210)
│   │       └── core/            # Lib
│   ├── backend-go/
│   │   ├── go.work
│   │   ├── iam-service/         # HTTP: 8000, gRPC: 9000
│   │   ├── mdm-service/         # HTTP: 8001, gRPC: 9001
│   │   └── bpm-service/         # HTTP: 8003, gRPC: 9003
│   └── backend-java/
│       ├── crm-service/         # HTTP: 8010, gRPC: 9010
│       ├── hrm-service/         # HTTP: 8011, gRPC: 9011
│       └── loan-service/        # HTTP: 8012, gRPC: 9012
├── libs/
│   ├── go/pkg/
│   └── java/
├── docs/
└── .github/workflows/
```

## Current Modules

| Module | Status | Port (Dev) |
| --- | --- | --- |
| `projects/shell` | Active host MFE | 4200 |
| `projects/iam` | Active remote MFE | 4201 |
| `projects/mdm` | Active remote MFE | 4202 |
| `iam-service` (Go) | Active | 8000 / 9000 |
| `mdm-service` (Go) | Active | 8001 / 9001 |
| `crm-service` (Java) | Active | 8010 / 9010 |

## Frontend

- Workspace: `apps/frontend-micro`
- Framework: Angular 21 with `@angular-architects/native-federation`
- Port Mapping: Shell (4200), Go MFEs (4201-4209), Java MFEs (4210-4219)
- Runtime remote config: `projects/shell/public/env.js`

- Build:
  ```powershell
  cd apps\frontend-micro
  npx ng build shell
  npx ng build iam
  npx ng build mdm
  ```

## Backend Go

- Workspace: `apps/backend-go/go.work`
- Shared Go package: `libs/go/pkg`
- Port Mapping: HTTP (800x), gRPC (900x)
- Run a service:
  ```powershell
  cd apps\backend-go\iam-service
  kratos run
  ```

## Backend Java

- Workspace: `apps/backend-java`
- Standard: Java 25 (LTS) with Spring Boot 4.0.6
- Core Engine: Virtual Threads (Loom) for imperative high-concurrency
- Migration: Flyway for database schema management
- Port Mapping: HTTP (801x), gRPC (901x)
- Language: Pure Java (migrated from Kotlin)
- Camunda 8 (Zeebe) integration active in CRM.

## Infrastructure

- Repo: `arda-labs/arda-infra`
- Runtime namespace for app workloads: `arda-apps`
- Gateway namespace: `gateway`
- Identity namespace: `identity`
- Local APISIX config: `arda-infra/local/apisix`

## Coding Standards

- Go: keep Kratos layering (`biz`, `data`, `service`, `server`) and run `gofmt`.
- Frontend: use standalone Angular components, signals where useful, and
  project-local conventions before adding new abstractions.
- Docs: distinguish current implementation from roadmap. Do not describe
  planned services as running services.
- Git: use Conventional Commits.
