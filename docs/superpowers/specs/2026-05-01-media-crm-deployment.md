# Media & CRM Deployment Design

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:writing-plans to create implementation plan.

**Goal:** Deploy the Media Service, BPM Bridge, CRM Service, and their infrastructure (SeaweedFS, Camunda 8) to k3s using ArgoCD, with external database integration.

**Architecture:** 
- **Infrastructure**: SeaweedFS and Camunda 8 deployed in the `infra` namespace.
- **Applications**: `media-service`, `bpm-bridge`, and `crm-service` deployed in the `arda-apps` namespace.
- **Networking**: Internal cluster communication via Service DNS.
- **Database**: External connectivity to ThinkCentre Postgres via Kubernetes `Endpoints` and `Service`.

**Tech Stack:** Kubernetes (k3s), ArgoCD, Kustomize, Docker.

---

## 1. Infrastructure (infra namespace)

### 1.1 SeaweedFS
- **Master Service**: Metadata management.
- **Volume Service**: File storage.
- **S3 API Service**: Provides S3-compatible interface for Media Service.

### 1.2 Camunda 8 (Zeebe)
- **Standalone Broker**: Resource-optimized for dev.
- **Gateway**: GRPC interface for Bridge and CRM workers.

### 1.3 External Postgres Integration
- **Manifests**: `Service` + `Endpoints` pointing to ThinkCentre IP on port 5432.
- **Alias**: Accessible as `postgres-external.infra.svc.cluster.local`.

## 2. Business Services (arda-apps namespace)

### 2.1 Media Service (Go)
- **Role**: Handles file uploads, metadata, and S3 persistence.
- **Integration**: Connects to SeaweedFS S3.

### 2.2 BPM Bridge (Go)
- **Role**: Custom BPM features (SLA, monitoring).
- **Integration**: Connects to Zeebe Gateway.

### 2.3 CRM Service (Kotlin)
- **Role**: Core CRM logic and BPM workers.
- **Integration**: Connects to Zeebe Gateway and Postgres.

## 3. ArgoCD Configuration

### 3.1 App Definitions
Each component will have a dedicated `Application` resource:
- `infra-seaweedfs`: `apps/base/seaweedfs`
- `infra-camunda8`: `apps/base/camunda8`
- `media-service`: `apps/media-service/overlays/dev`
- `bpm-bridge`: `apps/bpm-bridge/overlays/dev`
- `crm-service`: `apps/crm-service/overlays/dev`

### 3.2 Automated Sync
- Prune: Enabled
- Self-Heal: Enabled

## 4. Security & Networking
- **Isolation**: Services run in `arda-apps`, storage in `infra`.
- **Traffic**: No public ingress by default; cluster-internal only.
- **Cloudflare (Optional)**: If external access is needed later, a `Tunnel` can be mapped to internal service ports.

## 5. Development Workflow
- **Dockerfiles**: Created for all three services.
- **Local Deployment**: Use `k3s image import` for local images if a private registry is not available.
