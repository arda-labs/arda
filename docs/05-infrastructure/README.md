# Infrastructure Guide

Updated: 2026-04-30

Runtime infrastructure is owned by the sibling repository `arda-infra`.

## Current Runtime Model

| Component | Current state |
| --- | --- |
| Host | `thinkcenter` |
| Kubernetes | K3s |
| Gateway | Apache APISIX |
| Identity | Zitadel |
| GitOps | ArgoCD |
| App namespace | `arda-apps` |
| Gateway namespace | `gateway` |
| Identity namespace | `identity` |
| Infra namespace | `infra` |
| Database | PostgreSQL on `thinkcenter` |

There is one active shared cluster. Environment overlays are named `dev`, but
the namespace used by app workloads is `arda-apps`, not `arda-dev`.

## Infra Repository

```text
arda-infra/
├── argocd/
│   ├── apps/
│   └── projects/
├── apps/
│   ├── gateway/apisix/
│   ├── identity/zitadel/
│   ├── ingress/cloudflared/
│   ├── iam-service/
│   ├── mdm-service/
│   ├── mfe-shell/
│   ├── mfe-iam/
│   └── mfe-mdm/
├── infrastructure/
│   ├── namespaces.yaml
│   └── storageclass.yaml
├── local/apisix/
└── scripts/
```

## APISIX Routes

| Public path | Backend |
| --- | --- |
| `/*` | `mfe-shell` |
| `/mfe-iam/*` | `mfe-iam` |
| `/mfe-mdm/*` | `mfe-mdm` |
| `/api/v1/*` | `iam-service` |
| `/api/v1/mdm/*` | `mdm-service` |

API routes are rewritten from `/api/<path>` to `/<path>`.

## Local APISIX

Local integration checks should use `arda-infra/local/apisix`:

```powershell
cd D:\Github\arda-labs\arda-infra\local\apisix
docker compose up -d
```

Then open:

```text
http://localhost:9080
```

## Kustomize Checks

Before pushing infra changes:

```powershell
cd D:\Github\arda-labs\arda-infra
kubectl kustomize apps\iam-service\overlays\dev
kubectl kustomize apps\mdm-service\overlays\dev
kubectl kustomize apps\mfe-shell\overlays\dev
kubectl kustomize apps\mfe-iam\overlays\dev
kubectl kustomize apps\mfe-mdm\overlays\dev
```

## Secrets

Development bootstrap helpers are in `arda-infra/scripts`:

- `create-dev-secrets.sh`
- `bootstrap-dev-postgres.sql`
- `create-zitadel-secret.sh`

Do not commit real secrets. Use Kubernetes Secrets or sealed/external secret
management before production hardening.

## Roadmap Infrastructure

The following components are platform direction but are not all active runtime
dependencies yet:

- Redpanda for event streaming.
- Camunda for workflow orchestration.
- Garage S3 for object storage.
- Prometheus/Grafana/Loki/Alertmanager for observability.
- Production-grade secret management.
