# Arda Monorepo Guide

Updated: 2026-04-30

This file is a quick orientation guide for contributors and coding agents.
The source of truth for runtime manifests is the sibling repo `arda-infra`.

## Repository Layout

```text
arda/
├── apps/
│   ├── frontend-micro/
│   │   ├── angular.json
│   │   └── projects/
│   │       ├── shell/
│   │       ├── iam/
│   │       ├── mdm/
│   │       └── core/
│   ├── backend-go/
│   │   ├── go.work
│   │   ├── iam-service/
│   │   ├── mdm-service/
│   │   └── crm-service/
│   └── backend-java/
│       └── accounting_tmp/
├── libs/
│   └── go/pkg/
├── docs/
└── .github/workflows/
```

## Current Modules

| Module | Status |
| --- | --- |
| `apps/frontend-micro/projects/shell` | Active host MFE |
| `apps/frontend-micro/projects/iam` | Active IAM remote MFE |
| `apps/frontend-micro/projects/mdm` | Active MDM remote MFE |
| `apps/backend-go/iam-service` | Active Go service |
| `apps/backend-go/mdm-service` | Active Go service |
| `apps/backend-go/crm-service` | Skeleton / roadmap |
| `apps/backend-java/accounting_tmp` | Prototype |

## Frontend

- Workspace: `apps/frontend-micro`
- Framework: Angular 21 with `@angular-architects/native-federation`
- Projects: `shell`, `iam`, `mdm`, `core`
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
- Run a service:
  ```powershell
  cd apps\backend-go\iam-service
  kratos run
  ```
- Test:
  ```powershell
  go test ./...
  ```

## Backend Java

- Workspace: `apps/backend-java`
- Current prototype: `accounting_tmp`
- CI still expects a future `apps/backend-java/accounting` module, so Java
  pipeline and docs should be aligned before treating it as deployable.

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
