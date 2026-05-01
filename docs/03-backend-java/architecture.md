# Java Backend Architecture

Updated: 2026-05-01

The Java backend provides core-banking domains and complex business processes. It has been fully migrated to **Java 25 (LTS)** and **Spring Boot 4.0.6**.

## Current Structure

```text
apps/backend-java/
├── build.gradle.kts      # Global Spring Boot 4 and Java 25 configuration
├── settings.gradle.kts   # Foojay JDK resolver
└── crm-service/          # Active CRM service with Camunda 8
    ├── build.gradle.kts
    ├── Dockerfile
    └── src/main/java/io/arda/crm/
```

Shared logic is maintained in:

```text
libs/java/
├── common/               # Records for ArdaContext, ApiResponse, ErrorCode
├── database/             # BaseEntity and R2DBC support
├── grpc-client/          # Context propagation interceptors
├── messaging/            # CloudEvents and Kafka producers
└── security/             # WebFlux security filters and Gateway header trust
```

## Standards

- **Language**: Pure Java 25 (Kotlin has been removed to resolve JDK 25 toolchain issues).
- **Framework**: Spring Boot 4.0.6 (latest milestones).
- **Process Engine**: Camunda 8 (Zeebe) using the Spring Boot SDK.
- **Data Access**: Spring Data R2DBC for reactive persistence.
- **Data Structures**: Prefer **Java Records** for DTOs, API responses, and immutable context.

## Deployment

The Java pipeline is configured to build and containerize services using Gradle and Docker.
The `crm-service` is the first active production-grade Java module in the workspace.

## Near-term Modules

The intended Java layer is for core-banking and process-heavy domains:

- **CRM**: Customer lifecycle and onboarding via Camunda 8.
- **Accounting**: Double-entry ledger and financial reporting.
- **Loan**: Credit lifecycle and repayment schedules.
- **Deposit**: Account lifecycle and interest calculations.

## Configuration

Services use YAML-based configuration, standardized to work across local development and Kubernetes (ArgoCD/ConfigMaps).

Database user for CRM:

```text
postgres://crm:crm%40123@thinkcenter:5432/crm?sslmode=disable
```
