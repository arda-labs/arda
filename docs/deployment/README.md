# Deployment — Hướng dẫn Triển khai

> Hướng dẫn chi tiết triển khai Arda Platform trên K3s
> CI/CD pipeline, monitoring, và troubleshooting

---

## 📋 Overview

Arda Platform được triển khai trên **K3s** (Lightweight Kubernetes) với **ArgoCD** cho GitOps. Deployment được tự động hóa qua GitHub Actions CI/CD pipeline.

### Deployment Flow

```
Git Push → GitHub Actions → Docker Build & Push → ArgoCD Sync → K3s Update
```

---

## 🚀 Prerequisites

### System Requirements
- Ubuntu 24.04 LTS
- 32GB RAM minimum
- 100GB disk space minimum
- 4 CPU cores minimum

### Software Required
```bash
# K3s
curl -sfL https://get.k3s.io | sh -

# kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

# Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# ArgoCD CLI
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/download/v2.9.0/argocd-linux-amd64
chmod +x argocd-linux-amd64
sudo mv argocd-linux-amd64 /usr/local/bin/argocd

# Docker (for local builds)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
```

---

## 📁 Repository Structure

```
arda/
├── .github/
│   └── workflows/
│       ├── backend-go.yml      # Go services CI/CD
│       ├── backend-java.yml    # Java services CI/CD
│       └── frontend.yml        # Frontend CI/CD
│
├── arda-infra/                 # Infrastructure manifests
│   ├── apps/                   # ArgoCD applications
│   │   ├── database/
│   │   │   ├── postgresql/
│   │   │   └── redis/
│   │   ├── messaging/
│   │   │   └── redpanda/
│   │   ├── bpm/
│   │   │   └── camunda/
│   │   ├── storage/
│   │   │   └── garage/
│   │   ├── gateway/
│   │   │   └── apisix/
│   │   ├── identity/
│   │   │   └── zitadel/
│   │   ├── ingress/
│   │   │   └── cloudflared/
│   │   ├── iam-service/
│   │   ├── crm-service/
│   │   ├── accounting-service/
│   │   ├── loan-service/
│   │   ├── mfe-shell/
│   │   ├── mfe-common/
│   │   └── ...
│   ├── bootstrap/
│   │   └── root-app.yaml      # ArgoCD root application
│   └── infrastructure/
│       ├── namespaces.yaml
│       ├── storageclass.yaml
│       └── network-policies.yaml
│
├── arda-be/                    # Go services
│   ├── crm-service/
│   ├── hrm-service/
│   ├── notification-service/
│   ├── system-config-service/
│   └── bpm-service/
│
├── arda-core/                  # Java services
│   ├── services/
│   │   ├── accounting/
│   │   ├── loan/
│   │   └── deposit/
│   └── libs/
│
└── arda-mfe/                   # Frontend
    ├── apps/
    │   ├── shell/
    │   ├── common/
    │   ├── accounting/
    │   ├── loan/
    │   └── crm/
    └── libs/
```

---

## 🔧 Initial Setup

### 1. Install K3s

```bash
# Install K3s with custom configuration
curl -sfL https://get.k3s.io | sh -s - \
  --write-kubeconfig-mode 644 \
  --disable traefik \
  --disable servicelb \
  --node-name arda-server \
  --kube-apiserver-arg service-node-port-range=30000-32767

# Verify installation
sudo k3s kubectl get nodes
sudo k3s kubectl get pods --all-namespaces
```

### 2. Create Namespaces

```bash
# Create namespaces
kubectl apply -f arda-infra/infrastructure/namespaces.yaml

# namespaces.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: gateway
---
apiVersion: v1
kind: Namespace
metadata:
  name: database
---
apiVersion: v1
kind: Namespace
metadata:
  name: messaging
---
apiVersion: v1
kind: Namespace
metadata:
  name: bpm
---
apiVersion: v1
kind: Namespace
metadata:
  name: storage
---
apiVersion: v1
kind: Namespace
metadata:
  name: identity
---
apiVersion: v1
kind: Namespace
metadata:
  name: infra
---
apiVersion: v1
kind: Namespace
metadata:
  name: argocd
---
apiVersion: v1
kind: Namespace
metadata:
  name: arda-dev
```

### 3. Install ArgoCD

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
kubectl wait --for=condition=available \
  deployment/argocd-server \
  -n argocd \
  --timeout=300s

# Get ArgoCD admin password
kubectl get secret argocd-initial-admin-secret \
  -n argocd \
  -o jsonpath="{.data.password}" | base64 -d

# Port forward to access ArgoCD UI
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Open http://localhost:8080
# Username: admin
# Password: <from above command>
```

### 4. Configure ArgoCD Repository

```bash
# Add GitHub repository
argocd repo add https://github.com.arda_labss/arda-infra.git \
  --username <github-username> \
  --password <github-pat> \
  --insecure-skip-server-verify

# Create root application
kubectl apply -f arda-infra/bootstrap/root-app.yaml

# root-app.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: arda-root
  namespace: argocd
  finalizers:
  - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: https://github.com.arda_labss/arda-infra.git
    targetRevision: main
    path: apps
  destination:
    server: https://kubernetes.default.svc
    namespace: argocd
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

---

## 🐳 Container Registry

### GitHub Container Registry (GHCR)

```yaml
# .github/workflows/docker-build.yml
name: Docker Build & Push

on:
  push:
    branches: [ main ]
    paths:
      - 'arda-be/crm-service/**'
      - 'arda-be/iam-service/**'

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
    - name: Checkout
      uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: ./arda-be/crm-service
        push: true
        tags: |
          ghcr.io/arda-labs/crm-service:latest
          ghcr.io/arda-labs/crm-service:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max
```

---

## 📦 Deployment Manifests

### Example: Go Service Deployment

```yaml
# arda-infra/apps/crm-service/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crm-service
  namespace: arda-dev
  labels:
    app: crm-service
    version: v1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: crm-service
  template:
    metadata:
      labels:
        app: crm-service
        version: v1
    spec:
      serviceAccountName: crm-service
      initContainers:
      - name: config-init
        image: alpine:3.20
        command:
        - sh
        - -c
        - apk add --no-cache gettext -q && envsubst < /config-template/config.yaml > /config-rendered/config.yaml
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: crm-service-secrets
              key: database-url
        - name: REDIS_ADDR
          valueFrom:
            secretKeyRef:
              name: crm-service-secrets
              key: redis-addr
        - name: REDPANDA_BROKERS
          valueFrom:
            configMapKeyRef:
              name: crm-service-config
              key: redpanda-brokers
        volumeMounts:
        - name: config-template
          mountPath: /config-template
        - name: config-rendered
          mountPath: /config-rendered
      containers:
      - name: crm-service
        image: ghcr.io/arda-labs/crm-service:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8000
          protocol: TCP
        - name: grpc
          containerPort: 9000
          protocol: TCP
        env:
        - name: SERVICE_NAME
          value: "crm-service"
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        args:
        - -conf
        - /data/conf
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
        volumeMounts:
        - name: config-rendered
          mountPath: /data/conf
      volumes:
      - name: config-template
        configMap:
          name: crm-service-config
      - name: config-rendered
        emptyDir: {}
      affinity:
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - crm-service
              topologyKey: kubernetes.io/hostname
---
apiVersion: v1
kind: Service
metadata:
  name: crm-service
  namespace: arda-dev
  labels:
    app: crm-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: 8000
    protocol: TCP
  - name: grpc
    port: 9090
    targetPort: 9000
    protocol: TCP
  selector:
    app: crm-service
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: crm-service
  namespace: arda-dev
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: crm-service-hpa
  namespace: arda-dev
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: crm-service
  minReplicas: 2
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

### Example: Java Service Deployment (GraalVM)

```yaml
# arda-infra/apps/accounting-service/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: accounting-service
  namespace: arda-dev
  labels:
    app: accounting-service
    version: v1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: accounting-service
  template:
    metadata:
      labels:
        app: accounting-service
        version: v1
    spec:
      initContainers:
      - name: config-init
        image: alpine:3.20
        command:
        - sh
        - -c
        - apk add --no-cache gettext -q && envsubst < /config-template/application.yml > /config-rendered/application.yml
        env:
        - name: POSTGRES_HOST
          value: "postgresql.database.svc.cluster.local"
        - name: POSTGRES_DB
          valueFrom:
            secretKeyRef:
              name: accounting-service-secrets
              key: postgres-db
        - name: POSTGRES_USER
          valueFrom:
            secretKeyRef:
              name: accounting-service-secrets
              key: postgres-user
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: accounting-service-secrets
              key: postgres-password
        volumeMounts:
        - name: config-template
          mountPath: /config-template
        - name: config-rendered
          mountPath: /config-rendered
      containers:
      - name: accounting-service
        image: ghcr.io/arda-labs/accounting-service:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 8080
          protocol: TCP
        - name: grpc
          containerPort: 9090
          protocol: TCP
        env:
        - name: SPRING_PROFILES_ACTIVE
          value: "production"
        - name: JAVA_OPTS
          value: "-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0"
        - name: MANAGEMENT_ENDPOINTS_WEB_EXPOSURE_INCLUDE
          value: "health,info,metrics,prometheus"
        args:
        - "--spring.config.location=/config/application.yml"
        livenessProbe:
          httpGet:
            path: /actuator/health/liveness
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /actuator/health/readiness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        volumeMounts:
        - name: config-rendered
          mountPath: /config
      volumes:
      - name: config-template
        configMap:
          name: accounting-service-config
      - name: config-rendered
        emptyDir: {}
---
apiVersion: v1
kind: Service
metadata:
  name: accounting-service
  namespace: arda-dev
  labels:
    app: accounting-service
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: 8080
    protocol: TCP
  - name: grpc
    port: 9090
    targetPort: 9090
    protocol: TCP
  selector:
    app: accounting-service
```

### Example: Frontend Deployment

```yaml
# arda-infra/apps/mfe-shell/base/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mfe-shell
  namespace: arda-dev
  labels:
    app: mfe-shell
    version: v1
spec:
  replicas: 2
  selector:
    matchLabels:
      app: mfe-shell
  template:
    metadata:
      labels:
        app: mfe-shell
        version: v1
    spec:
      containers:
      - name: mfe-shell
        image: ghcr.io/arda-labs/mfe-shell:latest
        imagePullPolicy: Always
        ports:
        - name: http
          containerPort: 80
          protocol: TCP
        env:
        - name: API_BASE_URL
          value: "https://api.arda.io.vn"
        - name: AUTH_BASE_URL
          value: "https://auth.arda.io.vn"
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
---
apiVersion: v1
kind: Service
metadata:
  name: mfe-shell
  namespace: arda-dev
  labels:
    app: mfe-shell
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  selector:
    app: mfe-shell
```

---

## 🔄 CI/CD Pipelines

### Go Service Pipeline

```yaml
# .github/workflows/go-service.yml
name: Go Service CI/CD

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'arda-be/crm-service/**'
      - 'arda-be/pkg/**'
      - 'arda-be/api/**'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.25'

    - name: Download dependencies
      working-directory: ./arda-be/crm-service
      run: go mod download

    - name: Run tests
      working-directory: ./arda-be/crm-service
      run: go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        file: ./arda-be/crm-service/coverage.out

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: ./arda-be/crm-service
        push: true
        tags: |
          ghcr.io/arda-labs/crm-service:latest
          ghcr.io/arda-labs/crm-service:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Update ArgoCD
      uses: christian-korneck/update-argocd-action@v1
      env:
        AUTH_TOKEN: ${{ secrets.ARGOCD_AUTH_TOKEN }}
        ARGOCD_SERVER: argocd.arda.io.vn
        INSECURE_SKIP_TLS_VERIFY: "true"
        APP_NAME: crm-service-dev
        ARGOCD_NAMESPACE: argocd
```

### Java Service Pipeline

```yaml
# .github/workflows/java-service.yml
name: Java Service CI/CD

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'arda-core/services/accounting/**'
      - 'arda-core/libs/**'
      - 'arda-core/build.gradle.kts'
      - 'arda-core/settings.gradle.kts'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up GraalVM
      uses: graalvm/setup-graalvm@v1
      with:
        java-version: '21'
        distribution: 'graalvm-community'
        github-token: ${{ secrets.GITHUB_TOKEN }}

    - name: Build with Gradle
      working-directory: ./arda-core
      run: |
        ./gradlew :services:accounting:build -x test

    - name: Run tests
      working-directory: ./arda-core
      run: ./gradlew :services:accounting:test

    - name: Generate test report
      uses: dorny/test-reporter@v1
      if: always()
      with:
        name: Test Results
        path: '**/build/test-results/test/TEST-*.xml'
        reporter: java-junit
        fail-on-error: true

  build-native:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Set up GraalVM
      uses: graalvm/setup-graalvm@v1
      with:
        java-version: '21'
        distribution: 'graalvm-community'
        github-token: ${{ secrets.GITHUB_TOKEN }}

    - name: Build native image
      working-directory: ./arda-core
      run: ./gradlew :services:accounting:nativeCompile

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: ./arda-core/services/accounting
        file: ./arda-core/services/accounting/Dockerfile
        push: true
        tags: |
          ghcr.io/arda-labs/accounting-service:latest
          ghcr.io/arda-labs/accounting-service:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Update ArgoCD
      uses: christian-korneck/update-argocd-action@v1
      env:
        AUTH_TOKEN: ${{ secrets.ARGOCD_AUTH_TOKEN }}
        ARGOCD_SERVER: argocd.arda.io.vn
        INSECURE_SKIP_TLS_VERIFY: "true"
        APP_NAME: accounting-service-dev
        ARGOCD_NAMESPACE: argocd
```

### Frontend Pipeline

```yaml
# .github/workflows/frontend.yml
name: Frontend CI/CD

on:
  push:
    branches: [ main, develop ]
    paths:
      - 'arda-mfe/**'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '20'

    - name: Install dependencies
      working-directory: ./arda-mfe
      run: npm ci

    - name: Run lint
      working-directory: ./arda-mfe
      run: npx nx run-many -t lint

    - name: Run tests
      working-directory: ./arda-mfe
      run: npx nx run-many -t test --coverage

    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        directory: ./arda-mfe/coverage

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    permissions:
      contents: read
      packages: write

    steps:
    - uses: actions/checkout@v4

    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '20'

    - name: Install dependencies
      working-directory: ./arda-mfe
      run: npm ci

    - name: Build shell app
      working-directory: ./arda-mfe
      run: npx nx build shell --configuration=production

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Log in to GHCR
      uses: docker/login-action@v3
      with:
        registry: ghcr.io
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: ./arda-mfe/apps/shell
        file: ./arda-mfe/apps/shell/Dockerfile
        push: true
        tags: |
          ghcr.io/arda-labs/mfe-shell:latest
          ghcr.io/arda-labs/mfe-shell:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Update ArgoCD
      uses: christian-korneck/update-argocd-action@v1
      env:
        AUTH_TOKEN: ${{ secrets.ARGOCD_AUTH_TOKEN }}
        ARGOCD_SERVER: argocd.arda.io.vn
        INSECURE_SKIP_TLS_VERIFY: "true"
        APP_NAME: mfe-shell-dev
        ARGOCD_NAMESPACE: argocd
```

---

## 📊 Monitoring

### Prometheus Configuration

```yaml
# apps/monitoring/prometheus/base/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
  namespace: monitoring
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
      evaluation_interval: 15s

    scrape_configs:
    - job_name: 'kubernetes-pods'
      kubernetes_sd_configs:
      - role: pod
      relabel_configs:
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
        action: keep
        regex: true
      - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
        action: replace
        target_label: __metrics_path__
        regex: (.+)
      - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
        action: replace
        regex: ([^:]+)(?::\d+)?;(\d+)
        replacement: $1:$2
        target_label: __address__
      - action: labelmap
        regex: __meta_kubernetes_pod_label_(.+)
      - source_labels: [__meta_kubernetes_namespace]
        action: replace
        target_label: kubernetes_namespace
      - source_labels: [__meta_kubernetes_pod_name]
        action: replace
        target_label: kubernetes_pod_name
```

### Grafana Dashboards

```yaml
# apps/monitoring/grafana/base/dashboards/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-dashboards
  namespace: monitoring
data:
  arda-services.json: |
    {
      "dashboard": {
        "title": "Arda Services",
        "panels": [
          {
            "title": "Request Rate",
            "targets": [
              {
                "expr": "rate(http_requests_total[5m])"
              }
            ]
          },
          {
            "title": "Error Rate",
            "targets": [
              {
                "expr": "rate(http_requests_total{status=~\"5..\"}[5m]) / rate(http_requests_total[5m])"
              }
            ]
          },
          {
            "title": "Memory Usage",
            "targets": [
              {
                "expr": "container_memory_usage_bytes / container_spec_memory_limit_bytes"
              }
            ]
          },
          {
            "title": "CPU Usage",
            "targets": [
              {
                "expr": "rate(container_cpu_usage_seconds_total[5m])"
              }
            ]
          }
        ]
      }
    }
```

---

## 🔍 Troubleshooting

### Common Issues

#### 1. Pod stuck in Pending state

```bash
# Check pod events
kubectl describe pod <pod-name> -n <namespace>

# Check resource availability
kubectl top nodes
kubectl describe nodes

# Check resource quotas
kubectl get resourcequota -n <namespace>
```

#### 2. Pod crash looping

```bash
# Check pod logs
kubectl logs <pod-name> -n <namespace> --previous

# Check pod events
kubectl describe pod <pod-name> -n <namespace>

# Exec into pod (if running)
kubectl exec -it <pod-name> -n <namespace> -- sh
```

#### 3. Service not reachable

```bash
# Check service endpoints
kubectl get endpoints <service-name> -n <namespace>

# Check DNS resolution
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -- sh
nslookup <service-name>.<namespace>.svc.cluster.local

# Check network policies
kubectl get networkpolicies -n <namespace>
```

#### 4. ArgoCD sync failing

```bash
# Check ArgoCD application status
argocd app get <app-name>

# Check ArgoCD application logs
kubectl logs -n argocd deployment/argocd-application-controller

# Force sync
argocd app sync <app-name> --force

# Check application resources
argocd app resources <app-name>
```

### Health Check Endpoints

All services should expose health check endpoints:

```
GET /health/live    # Liveness probe
GET /health/ready   # Readiness probe
GET /metrics        # Prometheus metrics
```

---

## 🔄 Rollback

### ArgoCD Rollback

```bash
# View application history
argocd app history <app-name>

# Rollback to previous revision
argocd app rollback <app-name>

# Rollback to specific revision
argocd app rollback <app-name> --revision <revision-id>
```

### Manual Rollback

```bash
# Get deployment history
kubectl rollout history deployment/<deployment-name> -n <namespace>

# Rollback to previous version
kubectl rollout undo deployment/<deployment-name> -n <namespace>

# Rollback to specific revision
kubectl rollout undo deployment/<deployment-name> -n <namespace> --to-revision=<revision>
```

---

## 📝 Best Practices

### 1. Resource Management
- Always set `requests` and `limits` for all containers
- Use HPA for auto-scaling based on CPU/memory
- Monitor resource usage regularly

### 2. Security
- Use secrets for sensitive data
- Enable RBAC for service accounts
- Use network policies to restrict traffic
- Regularly scan images for vulnerabilities

### 3. Reliability
- Use multiple replicas for critical services
- Configure liveness and readiness probes
- Use pod anti-affinity for high availability
- Implement circuit breakers and retries

### 4. Observability
- Enable Prometheus metrics for all services
- Use structured logging
- Set up alerts for critical metrics
- Maintain audit logs

---

*Last Updated: 2026-04-24*
