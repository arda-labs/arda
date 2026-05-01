# CRM Service (Kotlin) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Chuyển đổi CRM sang Kotlin, xây dựng nghiệp vụ khách hàng và tích hợp Zeebe Workers.

**Architecture:** Spring Boot 3.4 (Reactive) + Zeebe SDK. Sử dụng Outbox pattern cho tin cậy.

**Tech Stack:** Kotlin 2.1, Spring Boot, R2DBC, Spring Zeebe.

---

### Task 1: Initialize Kotlin CRM Project

**Files:**
- Create: `arda/apps/backend-java/crm-service/build.gradle.kts`
- Create: `arda/apps/backend-java/crm-service/src/main/kotlin/io/arda/crm/CrmApplication.kt`

- [ ] **Step 1: Create build.gradle.kts with dependencies**

```kotlin
plugins {
    kotlin("jvm") version "2.1.0"
    id("org.springframework.boot") version "3.4.0"
}
dependencies {
    implementation("org.springframework.boot:spring-boot-starter-webflux")
    implementation("io.camunda:spring-zeebe-starter:8.4.0")
    implementation("org.springframework.boot:spring-boot-starter-data-r2dbc")
}
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-java/crm-service/
git commit -m "feat(crm): init kotlin spring boot project"
```

---

### Task 2: Implement Customer Domain & R2DBC Repository

**Files:**
- Create: `arda/apps/backend-java/crm-service/src/main/kotlin/io/arda/crm/domain/Customer.kt`

- [ ] **Step 1: Define Customer entity**

```kotlin
@Table("customers")
data class Customer(
    @Id val id: UUID?,
    val customerCode: String,
    val name: String,
    val status: String, // PENDING, ACTIVE, REJECTED
    val cccdFileId: UUID?
)
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-java/crm-service/src/main/kotlin/io/arda/crm/domain/
git commit -m "feat(crm): add customer domain model"
```

---

### Task 3: Implement Camunda Job Worker

**Files:**
- Create: `arda/apps/backend-java/crm-service/src/main/kotlin/io/arda/crm/worker/CustomerWorker.kt`

- [ ] **Step 1: Create Worker for customer activation**

```kotlin
@Component
class CustomerWorker(private val repository: CustomerRepository) {

    @JobWorker(type = "activate-customer")
    suspend fun activateCustomer(job: ActivatedJob) {
        val customerId = job.variablesAsMap["customerId"] as String
        // Logic to update DB status to ACTIVE
    }
}
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-java/crm-service/src/main/kotlin/io/arda/crm/worker/
git commit -m "feat(crm): implement zeebe job worker for activation"
```
