# Java Backend Architecture

Updated: 2026-04-30

The Java/Kotlin backend is currently a prototype area, not a deployed
production service layer.

## Current Structure

```text
apps/backend-java/
├── build.gradle.kts
├── settings.gradle.kts
└── accounting_tmp/
    ├── build.gradle.kts
    ├── Dockerfile
    └── src/main/kotlin/arda/accounting/
```

`accounting_tmp` contains early accounting domain models such as `Account` and
`Journal`. It should not be documented as the final accounting service until
the module name, API boundary, persistence model, and CI workflow are aligned.

## Target Direction

The intended Java layer is for core-banking domains that need stronger
transactional modeling than the operational Go services:

- accounting;
- loan;
- deposit;
- treasury.

The target stack may include Spring Boot, Gradle, PostgreSQL, and native-image
optimization, but those are roadmap choices until implemented in the repo.

## Current Gap

`.github/workflows/ci-java.yml` detects `apps/backend-java/accounting/**`, while
the current code is under `apps/backend-java/accounting_tmp/**`. Before Java is
treated as deployable, choose one of these paths:

1. Rename `accounting_tmp` to `accounting` and make the CI workflow real.
2. Keep `accounting_tmp` as a prototype and disable or retarget Java CI.

## Design Rule

Do not place shared reference lists or platform parameters in the Java core
services. Cross-domain reference data belongs in MDM unless a value is owned by
a specific product domain.
