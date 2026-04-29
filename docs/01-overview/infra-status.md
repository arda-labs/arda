# Infrastructure Status

Updated: 2026-04-30

This page distinguishes the current GitOps target from the last manually
captured live-cluster metrics. Refresh live values with `kubectl`/ArgoCD before
using this page during an incident.

## Current GitOps Target

| Area | Current target |
| --- | --- |
| Server | `thinkcenter` |
| Runtime | K3s |
| Gateway | APISIX in namespace `gateway` |
| Identity | Zitadel in namespace `identity` |
| App workloads | namespace `arda-apps` |
| Infra ingress | Cloudflared in namespace `infra` |
| GitOps | ArgoCD in namespace `argocd` |
| Database | PostgreSQL on `thinkcenter` |

## Current ArgoCD Applications

| Application | Path | Namespace |
| --- | --- | --- |
| `iam-service` | `apps/iam-service/overlays/dev` | `arda-apps` |
| `mdm-service` | `apps/mdm-service/overlays/dev` | `arda-apps` |
| `mfe-shell` | `apps/mfe-shell/overlays/dev` | `arda-apps` |
| `mfe-iam` | `apps/mfe-iam/overlays/dev` | `arda-apps` |
| `mfe-mdm` | `apps/mfe-mdm/overlays/dev` | `arda-apps` |
| `cloudflared` | `apps/ingress/cloudflared/overlays` | `infra` |
| `zitadel-routes` | `apps/identity/zitadel/base` | `identity` |

## Current Gateway Routes

| Route | Backend |
| --- | --- |
| `/*` | `mfe-shell` |
| `/mfe-iam/*` | `mfe-iam` |
| `/mfe-mdm/*` | `mfe-mdm` |
| `/api/v1/*` | `iam-service` |
| `/api/v1/mdm/*` | `mdm-service` |

## Last Live Capture

The last detailed live capture was taken on 2026-04-24, before the MDM rollout
and before the namespace/docs cleanup.

Known values from that capture:

- Server: `thinkcenter`
- OS: Ubuntu 24.04.4 LTS
- K3s: v1.34.6+k3s1
- APISIX: apache/apisix 3.16
- ArgoCD: v3.3.7
- Zitadel: v4.13.0
- PostgreSQL: external to K3s on `thinkcenter:5432`

Historical pod names such as `mfe-common` and namespace names such as
`arda-dev` should be treated as old snapshot data, not target architecture.

## External Services

| Service | Current note |
| --- | --- |
| PostgreSQL | Active on `thinkcenter`; IAM uses `iam`, MDM uses dedicated `mdm` |
| Redis | Configured in service configs; live availability still needs verification |
| Redpanda | Roadmap / not verified as active runtime |
| Monitoring stack | Roadmap / not verified as active runtime |

## Refresh Commands

```bash
kubectl get nodes -o wide
kubectl get pods -A -o wide
kubectl get svc -A
kubectl get apisixroute -A
kubectl get applications -n argocd
kubectl top nodes
kubectl top pods -n arda-apps
```

## Open Items

- Refresh live ArgoCD sync status after MDM rollout.
- Add resource requests/limits consistently.
- Verify external PostgreSQL backup policy.
- Verify Redis runtime status.
- Decide monitoring stack and log aggregation approach.
