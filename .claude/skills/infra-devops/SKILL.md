---
name: infra-devops
description: Hỗ trợ vận hành K3s, ArgoCD, APISIX, Redpanda và monitoring cho hạ tầng Arda
disable-model-invocation: false
---

# Infrastructure & DevOps Skill

Mục đích: Hỗ trợ quản lý hạ tầng và CI/CD trong thư mục `arda-infra/`.

## 🎯 Phạm vi

- **Kubernetes**: K3s (Ubuntu 24.04)
- **GitOps**: ArgoCD (Apps, Bootstrap)
- **API Gateway**: Apache APISIX
- **Event Broker**: Redpanda (Kafka compatible)
- **Identity**: Zitadel
- **Workflow**: Camunda 7
- **Containers**: Docker optimization

## 📦 Project Structure (Arda Infra)

```
arda-infra/
├── apps/                  # ArgoCD application manifests
├── bootstrap/             # Cluster bootstrap resources
└── infrastructure/        # Core infra (DB, Redis, Redpanda)
```

## 🛠️ Key Patterns

### 1. GitOps with ArgoCD
Mọi thay đổi hạ tầng phải được thực hiện qua YAML trong `arda-infra/`.

### 2. Resource Optimization
Cấu hình Resource Limits/Requests chặt chẽ để phù hợp với giới hạn 32GB RAM.
- PostgreSQL: 4GB
- Redis: 1GB
- Redpanda: 2GB

### 3. API Gateway (APISIX)
Cấu hình Routes, Upstreams và Plugins (Forward Auth, Rate Limiting).

## 🎯 Usage
- `/infra-devops "Deploy service [name] qua ArgoCD"`
- `/infra-devops "Cấu hình APISIX route cho [service]"`
- `/infra-devops "Optimize resource limits cho container [name]"`
