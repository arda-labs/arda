# Docker & CI for Media & CRM Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enable automatic build and push of Docker images for Media Service, BPM Bridge, and CRM Service to GHCR.

**Architecture:** 
- Go services: Multi-stage Docker build using `golang:1.26-bookworm` and `distroless/static-debian12`.
- Java service: Multi-stage Docker build using Gradle and `eclipse-temurin:21-jre-alpine`.
- CI: GitHub Actions workflows updated to detect and build new services.

**Tech Stack:** Docker, Go, Java, GitHub Actions.

---

### Task 1: Create Dockerfile for media-service

**Files:**
- Create: `apps/backend-go/media-service/Dockerfile`

- [ ] **Step 1: Write Dockerfile for media-service**

```dockerfile
FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/media-service ./apps/backend-go/media-service
WORKDIR /src/apps/backend-go/media-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -o /src/bin/media-service .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /src/bin/media-service /app/media-service
USER nonroot:nonroot
ENTRYPOINT ["/app/media-service"]
```

- [ ] **Step 2: Commit Dockerfile**

```bash
git add apps/backend-go/media-service/Dockerfile
git commit -m "feat: add Dockerfile for media-service"
```

### Task 2: Create Dockerfile for bpm-bridge

**Files:**
- Create: `apps/backend-go/bpm-bridge/Dockerfile`

- [ ] **Step 1: Write Dockerfile for bpm-bridge**

```dockerfile
FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/bpm-bridge ./apps/backend-go/bpm-bridge
WORKDIR /src/apps/backend-go/bpm-bridge
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -o /src/bin/bpm-bridge .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /src/bin/bpm-bridge /app/bpm-bridge
USER nonroot:nonroot
ENTRYPOINT ["/app/bpm-bridge"]
```

- [ ] **Step 2: Commit Dockerfile**

```bash
git add apps/backend-go/bpm-bridge/Dockerfile
git commit -m "feat: add Dockerfile for bpm-bridge"
```

### Task 3: Create Dockerfile for crm-service

**Files:**
- Create: `apps/backend-java/crm-service/Dockerfile`

- [ ] **Step 1: Write Dockerfile for crm-service**

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app
COPY . .
RUN ./gradlew :apps:backend-java:crm-service:build -x test

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=builder /app/apps/backend-java/crm-service/build/libs/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
```

- [ ] **Step 2: Commit Dockerfile**

```bash
git add apps/backend-java/crm-service/Dockerfile
git commit -m "feat: add Dockerfile for crm-service"
```

### Task 4: Update GitHub Actions workflows

**Files:**
- Modify: `.github/workflows/ci-go.yml`
- Modify: `.github/workflows/ci-java.yml`

- [ ] **Step 1: Update ci-go.yml to include new services**

```yaml
# Add to filters in ci-go.yml
media-service: 'apps/backend-go/media-service/**'
bpm-bridge: 'apps/backend-go/bpm-bridge/**'
```

- [ ] **Step 2: Update ci-java.yml to include crm-service**

```yaml
# Add to filters in ci-java.yml
crm-service: 'apps/backend-java/crm-service/**'
```

- [ ] **Step 3: Commit updates**

```bash
git add .github/workflows/ci-go.yml .github/workflows/ci-java.yml
git commit -m "ci: add media-service, bpm-bridge, and crm-service to build pipelines"
```
