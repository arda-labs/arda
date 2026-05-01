# Java Backend Architecture

Updated: 2026-05-02

The Java backend provides core-banking domains and complex business processes. It has been fully migrated to **Java 25 (LTS)** and **Spring Boot 4.0.6**, leveraging **Virtual Threads (Project Loom)** for a high-performance imperative model.

## Current Structure

```text
apps/backend-java/
├── build.gradle.kts      # Global Spring Boot 4 and Java 25 configuration
├── settings.gradle.kts   # Foojay JDK resolver
├── crm-service/          # 8010/9010 - Camunda 8, Flyway
├── hrm-service/          # 8011/9011
└── loan-service/         # 8012/9012
```

## Technology Stack

- **Language**: Pure Java 25 (LTS).
- **Framework**: Spring Boot 4.0.6.
- **Concurrency**: Virtual Threads enabled (`spring.threads.virtual.enabled=true`).
- **Data Access**: Spring Data JPA with Hibernate (Imperative).
- **Migration**: Flyway DB for versioned schema management.
- **Process Engine**: Camunda 8 (Zeebe) integration.

## Standards

- **Imperative Model**: Migrated away from Reactive (WebFlux) to simplify business logic while maintaining high concurrency via Virtual Threads.
- **Data Structures**: Use **Java Records** for DTOs and immutable context.
- **Error Handling**: Standardized `ApiResponse<T>` and `ArdaContext` propagation across thread boundaries.

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
