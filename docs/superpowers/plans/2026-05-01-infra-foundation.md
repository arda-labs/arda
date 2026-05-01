# Arda BPM & Storage Infrastructure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Thiết lập hạ tầng Camunda 8 và SeaweedFS trên cụm Kubernetes của Arda (K3s/Local) với cấu hình tối ưu tài nguyên.

**Architecture:** Sử dụng Kustomize để quản lý các manifest Kubernetes trong repo `arda-infra`. Camunda 8 sẽ chạy ở chế độ Standalone (1 Broker) cho môi trường Dev.

**Tech Stack:** Kubernetes (K3s), Zeebe, SeaweedFS, Elasticsearch, Kustomize.

---

### Task 1: Setup SeaweedFS Infrastructure

**Files:**
- Create: `arda-infra/apps/base/seaweedfs/kustomization.yaml`
- Create: `arda-infra/apps/base/seaweedfs/master-deployment.yaml`
- Create: `arda-infra/apps/base/seaweedfs/volume-deployment.yaml`
- Create: `arda-infra/apps/base/seaweedfs/service.yaml`

- [ ] **Step 1: Create SeaweedFS Master deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: seaweedfs-master
  namespace: arda-apps
spec:
  replicas: 1
  selector:
    matchLabels:
      app: seaweedfs-master
  template:
    metadata:
      labels:
        app: seaweedfs-master
    spec:
      containers:
      - name: master
        image: chrislusf/seaweedfs:latest
        args: ["master", "-ip=seaweedfs-master", "-port=9333"]
        ports:
        - containerPort: 9333
        - containerPort: 19333
```

- [ ] **Step 2: Create SeaweedFS Volume deployment**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: seaweedfs-volume
  namespace: arda-apps
spec:
  replicas: 1
  selector:
    matchLabels:
      app: seaweedfs-volume
  template:
    metadata:
      labels:
        app: seaweedfs-volume
    spec:
      containers:
      - name: volume
        image: chrislusf/seaweedfs:latest
        args: ["volume", "-mserver=seaweedfs-master:9333", "-port=8080", "-s3"]
        ports:
        - containerPort: 8080
        - containerPort: 18080
        - containerPort: 8333 # S3 port
```

- [ ] **Step 3: Create Service and Kustomization**

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: seaweedfs-s3
  namespace: arda-apps
spec:
  ports:
  - port: 8333
    targetPort: 8333
    name: s3
  selector:
    app: seaweedfs-volume
```

- [ ] **Step 4: Verify manifests**

Run: `kubectl kustomize arda-infra/apps/base/seaweedfs`
Expected: Hợp lệ, không có lỗi cú pháp.

- [ ] **Step 5: Commit**

```bash
git add arda-infra/apps/base/seaweedfs/
git commit -m "infra: add seaweedfs base manifests"
```

---

### Task 2: Setup Camunda 8 (Zeebe Standalone)

**Files:**
- Create: `arda-infra/apps/base/camunda8/kustomization.yaml`
- Create: `arda-infra/apps/base/camunda8/zeebe-deployment.yaml`
- Create: `arda-infra/apps/base/camunda8/service.yaml`

- [ ] **Step 1: Create Zeebe Deployment with resource limits**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zeebe
  namespace: arda-apps
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: zeebe
        image: camunda/zeebe:8.4.0
        env:
        - name: ZEEBE_BROKER_DATA_DISKUSAGECOMMANDWATERMARK
          value: "0.99"
        - name: ZEEBE_BROKER_DATA_DISKUSAGEREPLICATIONWATERMARK
          value: "0.99"
        resources:
          limits:
            memory: 1Gi
          requests:
            memory: 512Mi
```

- [ ] **Step 2: Create Zeebe Service**

```yaml
apiVersion: v1
kind: Service
metadata:
  name: zeebe-gateway
  namespace: arda-apps
spec:
  ports:
  - port: 26500
    name: gateway
  selector:
    app: zeebe
```

- [ ] **Step 3: Commit**

```bash
git add arda-infra/apps/base/camunda8/
git commit -m "infra: add zeebe standalone manifests"
```

---

### Task 3: Setup Databases for CRM & Media

**Files:**
- Modify: `arda-infra/scripts/bootstrap-dev-postgres.sql`

- [ ] **Step 1: Add CRM and Media databases to bootstrap script**

```sql
-- Add to end of file
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'crm') THEN
    CREATE ROLE crm LOGIN PASSWORD 'crm@123';
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'media') THEN
    CREATE ROLE media LOGIN PASSWORD 'media@123';
  END IF;
END
$$;

SELECT 'CREATE DATABASE crm OWNER crm' WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'crm')\gexec
SELECT 'CREATE DATABASE media OWNER media' WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'media')\gexec
```

- [ ] **Step 2: Commit**

```bash
git add arda-infra/scripts/bootstrap-dev-postgres.sql
git commit -m "infra: add crm and media database bootstrap"
```
