# Deployment Strategy

Updated: 2026-04-30

Arda uses application CI in `arda` and GitOps deployment state in
`arda-infra`.

## Flow

```text
push to arda/main
  -> GitHub Actions detects affected app/service
  -> build image and push to GHCR
  -> reusable GitOps workflow updates arda-infra/apps/<service>/overlays/dev
  -> ArgoCD syncs manifests to K3s
```

## Active CI

| Workflow | Builds |
| --- | --- |
| `ci-mfe.yml` | `mfe-shell`, `mfe-iam`, `mfe-mdm`, `mfe-ntf` |
| `ci-go.yml` | `iam-service`, `mdm-service`, `notification-service` |
| `ci-java.yml` | Placeholder; currently points to future `accounting` path |
| `gitops-update.yml` | Updates image tags in `arda-infra` |

## Image Naming

| App/service | Image |
| --- | --- |
| Shell MFE | `ghcr.io/arda-labs/mfe-shell:<tag>` |
| IAM MFE | `ghcr.io/arda-labs/mfe-iam:<tag>` |
| MDM MFE | `ghcr.io/arda-labs/mfe-mdm:<tag>` |
| NTF MFE | `ghcr.io/arda-labs/mfe-ntf:<tag>` |
| IAM service | `ghcr.io/arda-labs/iam-service:<tag>` |
| MDM service | `ghcr.io/arda-labs/mdm-service:<tag>` |
| Notification service | `ghcr.io/arda-labs/notification-service:<tag>` |

## Runtime Manifests

All runtime Kubernetes changes belong in `arda-infra`. Application code should
not hand-edit live cluster state.

Current app overlays:

```text
arda-infra/apps/iam-service/overlays/dev
arda-infra/apps/mdm-service/overlays/dev
arda-infra/apps/notification-service/overlays/dev
arda-infra/apps/mfe-shell/overlays/dev
arda-infra/apps/mfe-iam/overlays/dev
arda-infra/apps/mfe-mdm/overlays/dev
arda-infra/apps/mfe-ntf/overlays/dev
```

## Rules

- Keep image tags immutable for deploys; `latest` is useful for manual checks
  but not enough for auditability.
- Add resource requests/limits before treating a workload as production-ready.
- Route browser and API traffic through APISIX.
- Keep local-only host/CORS behavior in overlays or local APISIX config, not in
  production-oriented base manifests.
- Use secrets for database URLs and credentials; do not commit plain secrets.
