# Deployment Runbook

Updated: 2026-04-30

This is the short runbook for the current deploy model. The full source of
truth for manifests is `arda-infra`.

## Current Namespaces

| Namespace | Purpose |
| --- | --- |
| `argocd` | ArgoCD |
| `gateway` | APISIX |
| `identity` | Zitadel |
| `infra` | Cloudflared |
| `arda-apps` | Application workloads |

## Current Applications

ArgoCD app definitions are under `arda-infra/argocd/apps`.

| Application | Manifest path |
| --- | --- |
| `iam-service` | `apps/iam-service/overlays/dev` |
| `mdm-service` | `apps/mdm-service/overlays/dev` |
| `mfe-shell` | `apps/mfe-shell/overlays/dev` |
| `mfe-iam` | `apps/mfe-iam/overlays/dev` |
| `mfe-mdm` | `apps/mfe-mdm/overlays/dev` |
| `cloudflared` | `apps/ingress/cloudflared/overlays` |
| `zitadel-routes` | `apps/identity/zitadel/base` |

## Verify Manifests

```powershell
cd D:\Github\arda-labs\arda-infra
kubectl kustomize apps\iam-service\overlays\dev
kubectl kustomize apps\mdm-service\overlays\dev
kubectl kustomize apps\mfe-shell\overlays\dev
kubectl kustomize apps\mfe-iam\overlays\dev
kubectl kustomize apps\mfe-mdm\overlays\dev
```

## Verify Cluster

```bash
kubectl get pods -n arda-apps
kubectl get apisixroute -A
kubectl get applications -n argocd
```

## Local APISIX

Use local APISIX for browser integration tests:

```powershell
cd D:\Github\arda-labs\arda-infra\local\apisix
docker compose up -d
```

Expected local routes:

```text
http://localhost:9080
http://localhost:9080/api/v1/me
http://localhost:9080/api/v1/mdm/code-sets
http://localhost:9080/mfe-iam/remoteEntry.json
http://localhost:9080/mfe-mdm/remoteEntry.json
```

## Troubleshooting

```bash
kubectl describe pod <pod> -n arda-apps
kubectl logs -f deployment/<deployment> -n arda-apps
argocd app get <app>
argocd app sync <app>
```

## Known Gaps

- Java `accounting_tmp` is not currently deployed.
- APISIX forward-auth behavior is not yet fully formalized.
- Production secret management needs hardening before real production data.
- Monitoring stack is roadmap, not a current guaranteed runtime dependency.
