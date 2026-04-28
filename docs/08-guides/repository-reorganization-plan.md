# Repository Reorganization Plan

Updated: 2026-04-28

## Goal

Make development, infrastructure, and deployment responsibilities explicit before doing any large repository split.

The current GitHub organization already has a simple shape:

- `arda-labs/arda`
- `arda-labs/arda-infra`
- `arda-labs/.github`

The problem is not repository count. The problem is unclear ownership boundaries inside the repositories and inconsistent local/dev/prod contracts.

## Target Shape

### Keep

- `arda`: application monorepo.
- `arda-infra`: GitOps/runtime monorepo.
- `.github`: organization-level metadata only.

### Inside `arda`

```text
apps/
  frontend-micro/
  backend-go/
  backend-java/
libs/
  go/
  java/
docs/
.github/workflows/
```

Rules:

- Application code lives under `apps`.
- Shared source code lives under `libs`.
- Docs explain architecture and workflows, not cluster state.
- CI builds images and hands deployment intent to GitOps.
- No Kubernetes runtime manifests should be added here unless they are examples.

### Inside `arda-infra`

```text
argocd/
bootstrap/
apps/
  gateway/
  ingress/
  iam-service/
  mfe-shell/
  mfe-iam/
  mfe-common/
  base/
scripts/
```

Rules:

- `apps/*/base` contains shared Kubernetes resources.
- Environment-specific behavior belongs in overlays such as `overlays/local`, `overlays/dev`, and `overlays/prod`.
- APISIX routes, hosts, CORS, auth plugin configuration, and path rewrites belong here.
- The live cluster should be reproducible from this repo.

## Phase 0: Freeze The Current Reality

Status: in progress.

Actions:

- Document the current GitHub repositories.
- Document the source-of-truth boundary between `arda` and `arda-infra`.
- Mark the old monorepo migration plan as outdated.
- Keep local APISIX development documented.

Acceptance criteria:

- A new developer can tell where app code, infra code, docs, and gateway rules live.
- No one needs to infer the deploy model from scattered YAML files.

## Phase 1: Clean CI And Build Boundaries

Actions:

- Align Go CI with the Go version declared by the modules.
- Fix or remove `apps/backend-go/common-service` if it remains only a template.
- Fix `libs/go` module dependencies or move shared packages into a clean module boundary.
- Rename `apps/backend-java/accounting_tmp` to `apps/backend-java/accounting`, or mark Java as experimental and disable broken Java deploy CI.
- Ensure every CI job produces image tags that match an existing app path in `arda-infra`.

Acceptance criteria:

- `ci-go.yml`, `ci-java.yml`, and `ci-mfe.yml` either pass or are intentionally scoped to active code only.
- No CI workflow points to missing directories.
- Every image update has a matching GitOps overlay target.

## Phase 2: Normalize Environments

Actions:

- Define exactly three environment names: `local`, `dev`, `prod`.
- Add or clean Kustomize overlays in `arda-infra` for each active app.
- Keep localhost hosts and local-only CORS in the active dev overlays used by the shared APISIX on `thinkcenter`.
- Keep production base manifests free of workstation-specific behavior.
- Document how to switch frontend runtime config between local APISIX and deployed APISIX.

Acceptance criteria:

- A developer can run local frontend through APISIX without editing production manifests.
- Dev and prod differ only through overlays and secrets, not through undocumented manual cluster changes.
- `kubectl` live state and `arda-infra` manifests do not drift for APISIX service type, routes, and image tags.

## Phase 3: Formalize APISIX Auth

Actions:

- Keep Zitadel as the external IdP.
- Decide the production auth pattern:
  - APISIX verifies JWT directly, or
  - APISIX delegates to IAM through forward-auth.
- Standardize forwarded headers from APISIX to services.
- Make services trust identity headers only from APISIX/internal network.
- Add route-level auth policy for `/api/*`.
- Keep public MFE asset routes separate from protected API routes.

Acceptance criteria:

- Local, dev, and prod use the same auth decision path.
- Services do not implement different token behavior for local and deploy.
- Permission failures are observable as clear `401` or `403` responses.

## Phase 4: Decide Whether To Split Repositories

Do this only after phases 1-3 are complete.

Split only if there is a concrete ownership reason, such as separate teams, release cadences, compliance boundaries, or CI scale problems.

Potential future split:

- `arda-frontend`: Angular MFE workspace.
- `arda-backend-go`: Go services and Go shared libs.
- `arda-backend-java`: Java services and Java shared libs.
- `arda-infra`: unchanged.
- `arda-docs`: optional, only if docs need independent lifecycle.

Do not split just to make folders look cleaner. A premature split will make APISIX auth, image tag updates, cross-service changes, and docs coordination harder.

## Immediate Backlog

1. Fix Go version mismatch in CI.
2. Decide whether `common-service` is active or a template.
3. Fix Java accounting directory mismatch.
4. Add APISIX auth or forward-auth design and manifests.
5. Update `docs/01-overview/infra-status.md` after checking live cluster state.
6. Remove or archive outdated migration docs under `docs/superpowers` once this plan is accepted.
