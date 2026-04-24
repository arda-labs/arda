# Infrastructure — Hạ tầng K3s

> Mô tả chi tiết các thành phần hạ tầng của Arda Platform
> Target: 32GB RAM, Ubuntu 24.04 LTS, K3s

---

## 📋 Overview

Hạ tầng Arda được thiết kế để vận hành trên tài nguyên giới hạn (32GB RAM) với các component tối ưu hóa:
- **K3s**: Lightweight Kubernetes thay thế K8s đầy đủ
- **APISIX**: API gateway thay vì Traefik/Nginx
- **Zitadel**: Self-hosted IdP thay vì Auth0/Keycloak
- **Redpanda**: Kafka-compatible thay vì Kafka (tiết kiệm 2GB RAM)
- **Camunda 7**: Với PostgreSQL thay vì Camunda 8 + Elasticsearch (tiết kiệm 12GB RAM)
- **Garage S3**: Self-hosted object storage

---

## 🏗️ Architecture Diagram

```
                    ┌─────────────────────┐
                    │   Cloudflare CDN    │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │   Cloudflared       │
                    │   (Ingress Tunnel)  │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │   K3s Cluster       │
                    │                     │
                    │  ┌──────────────┐   │
                    │  │  APISIX GW   │   │
                    │  └──────┬───────┘   │
                    │         │           │
                    │  ┌──────▼───────┐   │
                    │  │   Zitadel    │   │
                    │  └──────┬───────┘   │
                    │         │           │
                    │  ┌──────▼───────┐   │
                    │  │   Services   │   │
                    │  │  (Go + Java) │   │
                    │  └──────┬───────┘   │
                    │         │           │
                    │  ┌──────▼───────┐   │
                    │  │  PostgreSQL  │   │
                    │  │   Redis      │   │
                    │  │  Redpanda    │   │
                    │  │  Camunda 7   │   │
                    │  │  Garage S3   │   │
                    │  └──────────────┘   │
                    └─────────────────────┘
```

---

## 📦 Components

### 1. K3s (Kubernetes)

**Role**: Container orchestration platform

**Configuration**:
```yaml
# K3s installation
curl -sfL https://get.k3s.io | sh -

# Verify installation
sudo k3s kubectl get nodes

# Configure resource limits
# /etc/rancher/k3s/config.yaml
kubelet-arg:
  - "max-pods=150"
  - "pod-max-pids-per-container=4096"
```

**Resource Usage**: ~3GB

**Notes**:
- Sử dụng containerd thay vì Docker
- Disable Traefik (dùng APISIX thay thế)
- Enable metrics-server cho monitoring

---

### 2. APISIX Gateway

**Role**: API Gateway, Load Balancer, API Management

**Deployment**:
```yaml
# apps/gateway/apisix/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: apisix
  namespace: gateway
spec:
  replicas: 2
  selector:
    matchLabels:
      app: apisix
  template:
    spec:
      containers:
      - name: apisix
        image: apache/apisix:3.9.0-debian
        ports:
        - containerPort: 9080  # HTTP
        - containerPort: 9443  # HTTPS
        - containerPort: 9180  # Admin API
        env:
        - name: APISIX_STAND_ALONE
          value: "false"
        - name: ETCD_ADVERTISE_CLIENT_URLS
          value: "http://apisix-etcd:2379"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

**Key Plugins**:
- `jwt-auth`: JWT token validation
- `forward-auth`: External auth with IAM service
- `rate-limit`: Rate limiting per tenant
- `proxy-rewrite`: URL rewriting
- `cors`: CORS handling
- `prometheus`: Metrics export

**Resource Usage**: ~1GB (2 replicas)

**Routes**:
- `/api/*` → Backend services
- `/auth/*` → Zitadel
- `/mfe/*` → Frontend MFEs

---

### 3. Zitadel (Identity Provider)

**Role**: Authentication & Authorization (OIDC/OAuth2)

**Deployment**:
```yaml
# apps/identity/zitadel/base/deployment.yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: zitadel
  namespace: identity
spec:
  chart:
    spec:
      chart: zitadel
      sourceRef:
        kind: HelmRepository
        name: zitadel
  values:
    zitadel:
      configMapConfig:
        tls:
          enabled: false
      masterkey:
        existingSecret: zitadel-secrets
      replicas: 1
      database:
        postgres:
          existingSecret: zitadel-db-secrets
      resources:
        requests:
          memory: "512Mi"
          cpu: "250m"
        limits:
          memory: "1Gi"
          cpu: "500m"
```

**Resource Usage**: ~1GB (1 replica + DB)

**Key Features**:
- Multi-tenant support
- MFA (TOTP, SMS, Email)
- SSO (SAML, OIDC)
- User management
- Session management

**Clients**:
- `iam-service`: Backend service auth
- `shell-app`: Frontend host app
- `mfe-*`: Individual MFEs

---

### 4. PostgreSQL

**Role**: Primary database for all services

**Deployment**:
```yaml
# apps/database/postgresql/base/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgresql
  namespace: database
spec:
  serviceName: postgresql
  replicas: 1
  selector:
    matchLabels:
      app: postgresql
  template:
    spec:
      containers:
      - name: postgresql
        image: postgres:16-alpine
        env:
        - name: POSTGRES_DB
          value: "arda"
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: postgresql-secrets
              key: username
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgresql-secrets
              key: password
        resources:
          requests:
            memory: "2Gi"
            cpu: "500m"
          limits:
            memory: "4Gi"
            cpu: "1000m"
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: postgresql-pvc
```

**Configuration**:
```sql
-- postgresql.conf
shared_buffers = 1GB
effective_cache_size = 3GB
maintenance_work_mem = 256MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 2621kB
min_wal_size = 1GB
max_wal_size = 4GB
```

**Resource Usage**: ~4GB

**Databases**:
- `arda_iam`: IAM service
- `arda_crm`: CRM service
- `arda_accounting`: Accounting service
- `arda_loan`: Loan service
- `arda_deposit`: Deposit service
- `zitadel`: Zitadel DB

---

### 5. Redis

**Role**: Caching, Session storage

**Deployment**:
```yaml
# apps/database/redis/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: database
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        command:
        - redis-server
        - --requirepass
        - $(REDIS_PASSWORD)
        env:
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: redis-secrets
              key: password
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "250m"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        emptyDir: {}
```

**Resource Usage**: ~512MB

**Use Cases**:
- User sessions
- JWT token blacklist
- Rate limiting
- Cache permissions
- System configuration

---

### 6. Redpanda

**Role**: Message broker for event-driven architecture

**Deployment**:
```yaml
# apps/messaging/redpanda/base/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redpanda
  namespace: messaging
spec:
  serviceName: redpanda
  replicas: 1
  selector:
    matchLabels:
      app: redpanda
  template:
    spec:
      containers:
      - name: redpanda
        image: redpandadata/redpanda:v23.3.11
        ports:
        - containerPort: 9092  # Kafka API
        - containerPort: 9644  # Admin API
        command:
        - redpanda
        - start
        - --smp
        - "1"
        - --memory
        - "1G"
        - --reserve-memory
        - "0M"
        - --overprovisioned
        - --node-id
        - "0"
        - --kafka-addr
        - INTERNAL://0.0.0.0:9092,EXTERNAL://0.0.0.0:19092
        - --advertise-kafka-addr
        - INTERNAL://redpanda-0.redpanda.messaging.svc.cluster.local:9092,EXTERNAL://redpanda.messaging.svc.cluster.local:19092
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
        volumeMounts:
        - name: data
          mountPath: /var/lib/redpanda/data
      volumes:
      - name: data
        emptyDir: {}
```

**Resource Usage**: ~2GB

**Topics**:
- `arda.loan.events` — Loan domain events
- `arda.accounting.events` — Accounting events
- `arda.notification.events` — Notification events
- `arda.outbox.*` — Outbox pattern events

---

### 7. Camunda 7

**Role**: BPM Engine for workflow orchestration

**Deployment**:
```yaml
# apps/bpm/camunda/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: camunda
  namespace: bpm
spec:
  replicas: 1
  selector:
    matchLabels:
      app: camunda
  template:
    spec:
      containers:
      - name: camunda
        image: camunda/camunda-bpm-platform:7.20.0
        env:
        - name: DB_DRIVER
          value: "org.postgresql.Driver"
        - name: DB_URL
          value: "jdbc:postgresql://postgresql.database.svc.cluster.local:5432/arda_bpm"
        - name: DB_USERNAME
          valueFrom:
            secretKeyRef:
              name: camunda-secrets
              key: username
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: camunda-secrets
              key: password
        - name: CAMUNDA_BPM_DATABASE_TYPE
          value: "postgres"
        # Disable Elasticsearch
        - name: ELASTICSEARCH_ENABLED
          value: "false"
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
```

**Resource Usage**: ~1GB

**Key Features**:
- BPMN 2.0 workflow execution
- User task management
- Service task integration
- Process variables
- History (on PostgreSQL)

**Workflows**:
- Loan approval workflow
- Account opening workflow
- Transaction approval workflow

---

### 8. Garage S3

**Role**: Object storage for files, documents

**Deployment**:
```yaml
# apps/storage/garage/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: garage
  namespace: storage
spec:
  replicas: 1
  selector:
    matchLabels:
      app: garage
  template:
    spec:
      containers:
      - name: garage
        image: dxflrs/garage:v1.0.0
        env:
        - name: AWS_ACCESS_KEY_ID
          valueFrom:
            secretKeyRef:
              name: garage-secrets
              key: access-key
        - name: AWS_SECRET_ACCESS_KEY
          valueFrom:
            secretKeyRef:
              name: garage-secrets
              key: secret-key
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "250m"
        volumeMounts:
        - name: data
          mountPath: /data
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: garage-pvc
```

**Resource Usage**: ~512MB

**Buckets**:
- `arda-documents` — Customer documents
- `arda-reports` — Generated reports
- `arda-media` — Media files
- `arda-backups` — Backup files

---

### 9. Cloudflared

**Role**: Tunnel to Cloudflare for public access

**Deployment**:
```yaml
# apps/ingress/cloudflared/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudflared
  namespace: infra
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cloudflared
  template:
    spec:
      containers:
      - name: cloudflared
        image: cloudflare/cloudflared:2024.4.0
        args:
        - tunnel
        - --config
        - /etc/cloudflared/config/config.yaml
        - run
        env:
        - name: TUNNEL_TOKEN
          valueFrom:
            secretKeyRef:
              name: cloudflared-secrets
              key: token
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
```

**Resource Usage**: ~128MB

**Domains**:
- `arda.io.vn` → K3s cluster
- `api.arda.io.vn` → APISIX
- `auth.arda.io.vn` → Zitadel

---

## 📊 Resource Summary

| Component | RAM Request | RAM Limit | CPU Request | CPU Limit |
|-----------|-------------|-----------|-------------|-----------|
| K3s + OS | 3GB | 3GB | - | - |
| PostgreSQL | 2GB | 4GB | 500m | 1000m |
| Redis | 256MB | 512MB | 100m | 250m |
| Redpanda | 1GB | 2GB | 500m | 1000m |
| Camunda 7 | 512MB | 1GB | 250m | 500m |
| Garage S3 | 256MB | 512MB | 100m | 250m |
| APISIX (2x) | 512MB | 1GB | 500m | 1000m |
| Zitadel | 512MB | 1GB | 250m | 500m |
| Cloudflared | 64MB | 128MB | 50m | 100m |
| **Infrastructure Total** | **~9.5GB** | **~13.5GB** | **~2.3 cores** | **~4.6 cores** |

---

## 🚀 Deployment

### Prerequisites
```bash
# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Install helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### Deploy Infrastructure
```bash
# Add helm repositories
helm repo add apisix https://charts.apiseven.com
helm repo add zitadel https://charts.zitadel.com
helm repo update

# Apply all manifests
kubectl apply -f arda-infra/infrastructure/namespaces.yaml
kubectl apply -f arda-infra/infrastructure/storageclass.yaml

# Deploy via ArgoCD
kubectl apply -f arda-infra/bootstrap/root-app.yaml
```

### Verify Deployment
```bash
# Check all pods
kubectl get pods --all-namespaces

# Check ArgoCD applications
kubectl get applications -n argocd

# Access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443
# Open http://localhost:8080
# Username: admin
# Password: kubectl get pods -n argocd -l app.kubernetes.io/name=argocd-server -o name | cut -d'/' -f 2
```

---

## 🔧 Configuration

### Environment Variables

Create `arda-infra/.env`:
```bash
# Database
POSTGRES_USERNAME=arda
POSTGRES_PASSWORD=change_me
POSTGRES_HOST=postgresql.database.svc.cluster.local
POSTGRES_PORT=5432

# Redis
REDIS_PASSWORD=change_me
REDIS_HOST=redis.database.svc.cluster.local
REDIS_PORT=6379

# Redpanda
REDPANDA_BROKERS=redpanda-0.redpanda.messaging.svc.cluster.local:9092

# Zitadel
ZITADEL_ADMIN_PASSWORD=change_me
ZITADEL_LOGIN_CLIENT_PAT=change_me

# Cloudflare
CLOUDFLARE_TUNNEL_TOKEN=change_me
```

### Secrets Management

```bash
# Create secrets for PostgreSQL
kubectl create secret generic postgresql-secrets \
  --from-literal=username=arda \
  --from-literal=password=change_me \
  -n database

# Create secrets for Redis
kubectl create secret generic redis-secrets \
  --from-literal=password=change_me \
  -n database

# Create secrets for Zitadel
kubectl create secret generic zitadel-secrets \
  --from-literal=masterkey=change_me \
  -n identity
```

---

## 📈 Monitoring

### Metrics Collection
```yaml
# Enable metrics in each service
metrics:
  enabled: true
  port: 9090
  path: /metrics
```

### Monitoring Stack
- **Prometheus**: Metrics collection
- **Grafana**: Visualization
- **Loki**: Log aggregation
- **Alertmanager**: Alerting

---

*Last Updated: 2026-04-24*
