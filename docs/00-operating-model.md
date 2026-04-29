# Arda Operating Model

Updated: 2026-04-30

This document is the source of truth for how Arda organizes code, infrastructure, local development, and deployment.

## Current GitHub Repositories

GitHub organization: `arda-labs`

| Repository | Role | Current status |
| --- | --- | --- |
| `arda-labs/arda` | Application monorepo: frontend, backend, shared libraries, docs, app CI | Active |
| `arda-labs/arda-infra` | GitOps and Kubernetes runtime configuration | Active |
| `arda-labs/.github` | Organization profile and GitHub metadata | Active |

There are no active separate GitHub repositories for frontend, Go backend, Java backend, or individual services.

## Ownership Boundaries

`arda` owns application source code:

- `apps/frontend-micro`: Angular micro frontends.
- `apps/backend-go`: Go services.
- `apps/backend-java`: Java services and prototypes.
- `libs/go`: shared Go packages.
- `libs/java`: shared Java packages.
- `docs`: product, architecture, and engineering documentation.
- `.github/workflows`: CI that builds application artifacts and updates GitOps image tags.

`arda-infra` owns runtime state:

- ArgoCD projects and applications.
- Kubernetes manifests and Kustomize overlays.
- APISIX gateway installation and routes.
- Cloudflared ingress.
- Runtime namespaces, service exposure, and deploy-time environment configuration.

Application repositories must not be the source of truth for Kubernetes runtime state. Infra manifests must not contain application implementation code.

## Environment Model

Arda uses three environment names:

| Environment | Purpose | Entry point |
| --- | --- | --- |
| `local` | Developer workstation using the same gateway contract as deploy | APISIX via SSH tunnel or local APISIX |
| `dev` | Shared integration environment on the cluster | APISIX in Kubernetes |
| `prod` | Production runtime | APISIX in Kubernetes |

The important rule is that auth, routing, CORS, tenant headers, and permission behavior should be tested through APISIX. Direct service calls are acceptable for isolated unit or service debugging, but they are not production-parity checks.

## Gateway And Auth Contract

APISIX is the edge contract for browser and external traffic.

Expected responsibilities:

- Route `/api/*` to backend services.
- Route MFE assets such as `/mfe-iam/*` and `/mfe-mdm/*`.
- Normalize path rewrites before traffic reaches services.
- Enforce or delegate token verification.
- Forward trusted identity, tenant, membership, and permission context to services after verification.
- Keep CORS behavior consistent between local, dev, and prod.

Zitadel is the source of truth for:

- User authentication.
- OIDC/OAuth2 login flows.
- Hosted login UI and client registration.
- Issued identity tokens and claims.

IAM remains the source of truth for:

- Login and identity integration.
- Tenant membership.
- Roles and permissions.
- Menu and capability authorization.
- Forward-auth decisions if APISIX delegates authorization to IAM.

## Local Development Rule

For frontend or integration work, use APISIX locally:

```powershell
ssh -N -L 9080:127.0.0.1:32459 hoan@thinkcenter
```

Then point local frontend runtime config to:

```js
window.__env.apiUrl = 'http://localhost:9080/api';
window.__env.mfeIamUrl = 'http://localhost:9080/mfe-iam';
window.__env.mfeMdmUrl = 'http://localhost:9080/mfe-mdm';
```

This keeps local browser traffic on the same route shape as deployed traffic.

## CI/CD Contract

Application CI in `arda` should:

- Detect changed frontend or service modules.
- Build and test only affected modules when possible.
- Publish immutable images to GHCR using short Git SHA tags.
- Update `arda-infra` GitOps overlays with the new image tag.

GitOps in `arda-infra` should:

- Own which image tag is running in each environment.
- Own route, host, CORS, config, secret reference, and scaling behavior.
- Avoid drift between live cluster state and repository manifests.

## Current Technical Debt

These items block a clean operating model:

- `mdm-service` now replaces the old `common-service` skeleton; shared master data should use MDM naming, routes, and ownership.
- `libs/go/pkg` references `github.com/shopspring/decimal` without a complete module dependency setup.
- Java CI expects `apps/backend-java/accounting`, but current code uses `accounting_tmp`.
- APISIX routes currently cover routing, but auth plugin or forward-auth behavior is not yet formalized.
- Local APISIX host/CORS support is kept in `arda-infra` dev overlays so the shared APISIX on `thinkcenter` can serve local tunnel traffic without putting workstation behavior in production-oriented base manifests.
- `docs/superpowers/plans/2026-04-25-monorepo-migration.md` describes an older migration state and should not be used as the current operating plan.

## Decision

Do not split repositories yet.

The current codebase is still stabilizing service boundaries, CI, gateway auth, and deploy overlays. Splitting now would multiply coordination cost without solving the root problems.

The near-term target is:

1. Keep `arda` as the application monorepo.
2. Keep `arda-infra` as the only GitOps/runtime repository.
3. Clean module boundaries inside `arda`.
4. Clean environment overlays inside `arda-infra`.
5. Revisit repository splitting only after CI, APISIX auth, and deploy contracts are stable.
