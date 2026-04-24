# Arda Infrastructure Status — Real-world Deployment

> Cập nhật thực trạng deployment trên server **thinkcenter**
> Date: 2026-04-24

---

## 🖥️ Server Hardware

**Server**: thinkcenter
- **CPU**: Intel Xeon E-2286M (8 cores, 16 threads, 2.4GHz base, up to 5.0GHz turbo)
- **RAM**: 16GB (planned upgrade to 32GB)
- **OS**: Ubuntu 24.04.4 LTS
- **Kernel**: 6.8.0-110-generic
- **IP**: 192.168.100.5
- **Container Runtime**: containerd 2.2.2

**Note**: Server đang chạy 16 logical cores (8 cores physical với hyperthreading)

---

## 🏗️ K3s Cluster

```
NAME        STATUS   ROLES           AGE   VERSION
thinkcenter Ready    control-plane   5d    v1.34.6+k3s1
```

- **K3s Version**: v1.34.6+k3s1
- **Kubectl Version**: v1.34.1
- **Node IP**: 192.168.100.5
- **Single-node cluster** (control-plane + worker)

---

## 📊 Resource Usage (Current)

```
Node: thinkcenter
- CPU: 497m (3% of 12 cores)
- RAM: 3858Mi (24% of 16GB)

Pods in arda-dev:
- iam-service: 1m CPU, 8Mi RAM
- mfe-common: 1m CPU, 13Mi RAM
- mfe-shell: 1m CPU, 13Mi RAM
- zitadel: 15m CPU, 111Mi RAM
- zitadel-login: 10m CPU, 106Mi RAM
- zitadel-postgresql: 41m CPU, 88Mi RAM
```

**Observation**: Resource usage rất thấp (3% CPU, 24% RAM), còn nhiều dư để thêm services.

---

## 🗂️ Namespaces Structure

```
NAME              STATUS   AGE
default           Active   5d
arda-dev          Active   2d23h     ← Dev environment
arda-prod         Active   4d3h      ← Prod environment (empty)
argocd            Active   4d8h      ← ArgoCD
gateway           Active   4d3h      ← APISIX Gateway
infra             Active   4d3h      ← Infrastructure (Cloudflared)
kube-system       Active   5d        ← K8s system
kube-public       Active   5d
kube-node-lease   Active   5d
```

**Observation**: `arda-prod` namespace đang rỗng, chỉ có `arda-dev` đang active.

---

## 🚀 Running Services

### arda-dev Namespace

| Service | Status | Image | Replicas | Ports |
|---------|--------|-------|----------|-------|
| iam-service | Running | ghcr.io/arda-labs/iam-service:de92316 | 1 | 8000, 9000 |
| mfe-common | Running | ghcr.io/arda-labs/mfe-common:1018fe9 | 1 | 80 |
| mfe-shell | Running | ghcr.io/arda-labs/mfe-shell:1018fe9 | 1 | 80 |
| zitadel | Running | ghcr.io/zitadel/zitadel:v4.13.0 | 1 | 8080 |
| zitadel-login | Running | ghcr.io/zitadel/zitadel-login:v4.13.0 | 1 | 3000 |
| zitadel-postgresql | Running | postgresql (StatefulSet) | 1 | 5432 |

### gateway Namespace

| Service | Status | Image | Replicas | Ports |
|---------|--------|-------|----------|-------|
| apisix | Running | apache/apisix:3.16.0-ubuntu | 1 | 9080, 9180, 9443 |
| apisix-etcd | Running | bitnami/etcd | 3 | 2379, 2380 |
| apisix-ingress-controller | Running | apache/apisix-ingress-controller:2.0.1 | 1 | 8080 |

### argocd Namespace

| Service | Status | Image | Replicas | Ports |
|---------|--------|-------|----------|-------|
| argocd-server | Running | quay.io/argoproj/argocd:v3.3.7 | 1 | 80, 443 |
| argocd-application-controller | Running | quay.io/argoproj/argocd:v3.3.7 | 1 | 8082 |
| argocd-repo-server | Running | quay.io/argoproj/argocd:v3.3.7 | 1 | 8081, 8084 |
| argocd-redis | Running | redis:8.2.3-alpine | 1 | 6379 |

### infra Namespace

| Service | Status | Image | Replicas | Ports |
|---------|--------|-------|----------|-------|
| cloudflared | Running | cloudflare/cloudflared:latest | 1 | - |

---

## 🔗 External Services (Bên ngoài K3s)

### PostgreSQL (Database)
- **Host**: thinkcenter (192.168.100.5)
- **Port**: 5432
- **User**: iam
- **Database**: iam
- **Connection**: `postgres://iam:iam@123@thinkcenter:5432/iam?sslmode=disable`
- **Status**: ✅ Running (external to K3s)

### Redis
- **Host**: thinkcenter (192.168.100.5)
- **Port**: 6379
- **Status**: ⚠️ Configured but not verified

### Kafka/Redpanda
- **Status**: ❌ Not installed yet

---

## 🌐 Ingress & Routes

### Cloudflare Tunnel

```
Domain                  Service                                    TLS
arda.io.vn           → apisix-gateway.gateway:80               noTLSVerify
auth.arda.io.vn      → traefik.kube-system.svc.cluster.local:80 noTLSVerify
argocd.arda.io.vn    → argocd-server.argocd.svc.cluster.local:443
```

### APISIX Routes (ApisixRoute)

| Name | Host | Path | Backend |
|------|------|------|---------|
| iam-service-v2-route | arda.io.vn | /api/v1/* | iam-service:80 |
| mfe-common-route | arda.io.vn | /common/* | mfe-common:80 |
| mfe-shell-route | arda.io.vn | /* | mfe-shell:80 |
| zitadel-route | auth.arda.io.vn | /* | Zitadel (via Traefik) |

### Traefik Ingress

| Name | Host | Service |
|------|------|---------|
| zitadel | auth.arda.io.vn | zitadel:8080 |
| zitadel-login | auth.arda.io.vn | zitadel-login:3000 |

---

## 📦 ArgoCD Applications

| Application | Repo | Path | Namespace | Sync Status |
|-------------|------|------|-----------|-------------|
| arda-root | github.com.arda_labss/arda-infra | argocd/apps | argocd | ✅ Synced |
| apisix | github.com.arda_labss/arda-infra | apps/gateway/apisix/overlays/dev | gateway | ✅ Synced |
| cloudflared | github.com.arda_labss/arda-infra | apps/ingress/cloudflared/overlays | infra | ✅ Synced |
| iam-service-dev | github.com.arda_labss/arda-infra | apps/iam-service/overlays/dev | arda-dev | ✅ Synced |
| mfe-common-dev | github.com.arda_labss/arda-infra | apps/mfe-common/overlays/dev | arda-dev | ✅ Synced |
| mfe-shell-dev | github.com.arda_labss/arda-infra | apps/mfe-shell/overlays/dev | arda-dev | ⚠️ OutOfSync |

**Observation**: `mfe-shell-dev` đang OutOfSync, cần sync lại.

---

## ⚠️ Issues & Gaps

### 1. Resource Limits Not Set
**Problem**: Các deployments không có `requests` và `limits` cho CPU/RAM.

```yaml
resources: {}  # Empty - needs to be set
```

**Impact**: Không có resource isolation, một pod có thể tiêu thụ toàn bộ RAM và gây OOM.

**Fix Needed**: Set appropriate resource limits for all deployments.

### 2. No Persistent Volumes
**Problem**: Không có PersistentVolume hoặc PersistentVolumeClaim.

**Impact**: Dữ liệu sẽ mất khi pod restart (ngoại trừ Zitadel PostgreSQL dùng StatefulSet).

**Fix Needed**: Set up PV/PVC cho databases, logs, và storage cần persistence.

### 3. Redis Not Configured Properly
**Problem**: REDIS_ADDR được set nhưng Redis service không có trong K3s (chỉ có external config).

**Status**: External Redis cần được verify.

### 4. Kafka/Redpanda Not Installed
**Problem**: Message broker chưa được cài đặt.

**Impact**: Không thể implement Event-Driven Architecture và Saga pattern.

**Fix Needed**: Install Redpanda for event messaging.

### 5. No Monitoring Stack
**Problem**: Không có Prometheus, Grafana, hay Alerting.

**Impact**: Không có visibility vào performance và issues.

**Fix Needed**: Install Prometheus + Grafana for monitoring.

### 6. ArgoCD Repo Structure Mismatch
**Problem**: ArgoCD apps track `github.com.arda_labss/arda-infra` nhưng local repo có thể khác.

**Action**: Verify repo structure sync.

---

## 📋 Services Structure (Actual)

### arda-infra/ Repository Structure

```
arda-infra/
├── argocd/apps/
│   ├── apisix.yaml
│   ├── arda-root.yaml
│   ├── cloudflared.yaml
│   ├── iam-service-dev.yaml
│   ├── mfe-common-dev.yaml
│   └── mfe-shell-dev.yaml
│
└── apps/
    ├── gateway/
    │   └── apisix/
    │       └── overlays/dev/
    │
    ├── iam-service/
    │   └── overlays/dev/
    │
    ├── mfe-common/
    │   └── overlays/dev/
    │
    └── mfe-shell/
        └── overlays/dev/
```

---

## 🎯 Immediate Action Items

### Priority 1 (Critical)
- [ ] Add resource limits to all deployments
- [ ] Verify external PostgreSQL is properly configured and backed up
- [ ] Verify external Redis is running and accessible
- [ ] Fix `mfe-shell-dev` OutOfSync issue

### Priority 2 (High)
- [ ] Install Redis cluster inside K3s (or verify external)
- [ ] Install Redpanda for event messaging
- [ ] Set up PersistentVolumes for data persistence
- [ ] Install monitoring stack (Prometheus + Grafana)

### Priority 3 (Medium)
- [ ] Implement HPA (Horizontal Pod Autoscaler)
- [ ] Set up log aggregation (Loki)
- [ ] Configure backup strategy for databases
- [ ] Set up alerting rules

---

## 💡 Recommendations for 16GB RAM (Current)

Given current hardware (16GB RAM), here's the recommended resource allocation:

```
Infrastructure:     8GB   (PostgreSQL external 2GB, Redis external 1GB, Zitadel 1GB, K3s 1GB, System 3GB)
Running Services:   1GB   (Current: ~300MB used)
Reserved:           7GB   (Future services, buffer)
```

**Current utilization**: ~400MB/16GB = 2.5%

**Recommendation**: Can safely add:
- 5-10 Go services (~50-100MB each)
- 2-3 Java services (~200-400MB each with GraalVM)
- Redis cluster inside K3s
- Redpanda (single broker ~1GB)
- Monitoring stack (~1GB)

---

## 💡 Recommendations for 32GB RAM (Future)

When upgraded to 32GB RAM:

```
Infrastructure:     10GB  (PostgreSQL 4GB, Redis 1GB, Redpanda 2GB, Camunda 1GB, S3 0.5GB, System 1.5GB)
Core Services:      4GB   (4 Java services × ~400MB + overhead)
Operational:        2GB   (8 Go services × ~100MB + overhead)
Frontend:           1GB   (Nx build cache, dev server)
Monitoring:         2GB   (Prometheus, Grafana, Loki)
K3s/OS:             3GB   (K8s overhead, system)
Reserved:           10GB  (Future services, buffer)
```

---

## 🔧 Quick Commands for Troubleshooting

### Check all pods
```bash
kubectl get pods --all-namespaces -o wide
```

### Check resource usage
```bash
kubectl top nodes
kubectl top pods -n arda-dev
```

### Check services
```bash
kubectl get svc --all-namespaces
```

### Check ArgoCD status
```bash
argocd app list
argocd app get <app-name>
argocd app sync <app-name>
```

### Check APISIX routes
```bash
kubectl get apisixroute -A
```

### Check ingress
```bash
kubectl get ingress -A
```

### Check logs
```bash
kubectl logs -f <pod-name> -n <namespace>
kubectl logs -f deployment/<deployment-name> -n <namespace>
```

---
