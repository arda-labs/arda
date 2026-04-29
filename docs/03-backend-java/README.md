# Backend Java Guide

Updated: 2026-04-30

The Java area is currently a prototype workspace.

## Current Module

```text
apps/backend-java/accounting_tmp
```

This module is useful for exploring accounting domain modeling, but it is not
yet part of the deployed Arda runtime.

## Commands

From `apps/backend-java`:

```powershell
.\gradlew projects
.\gradlew build
```

Run module-specific Gradle tasks only after checking the actual task names:

```powershell
.\gradlew tasks
```

## Before Productionizing

- Rename `accounting_tmp` to the final module name, likely `accounting`.
- Align `.github/workflows/ci-java.yml` with the final module path.
- Decide the HTTP/gRPC API surface.
- Add migrations and persistence strategy.
- Add manifests in `arda-infra` only after the service is runnable.
- Add docs under `docs/06-features/accounting.md` that distinguish delivered
  behavior from roadmap behavior.
