# Documentation Audit And Proposals

Updated: 2026-04-30

This audit records the documentation cleanup after the move from the old
`common` concept to MDM and after confirming the current frontend/infra
architecture.

## What Was Corrected

| Area | Correction |
| --- | --- |
| Frontend architecture | Replaced old Nx/Rspack/Analog description with current Angular CLI + Native Federation structure |
| Frontend modules | Documented `shell`, `iam`, `mdm`, and `core` under `apps/frontend-micro/projects` |
| MFE routing | Standardized shell routes `/iam/*`, `/mdm/*` and asset routes `/mfe-iam/*`, `/mfe-mdm/*` |
| Backend Go | Documented active `iam-service` and `mdm-service`; marked `crm-service` as skeleton |
| Backend Java | Marked `accounting_tmp` as prototype and documented CI path mismatch |
| Infrastructure | Replaced old `arda-dev` assumptions with current `arda-apps` namespace |
| Deployment | Replaced large future-state examples with current GitOps flow and active manifests |
| Organization profile | Updated `.github/profile/README.md` to reflect current repos and modules |

## Current Source-Of-Truth Docs

- `docs/00-operating-model.md`
- `docs/01-overview/architecture.md`
- `docs/01-overview/tech-stack.md`
- `docs/04-frontend/architecture.md`
- `docs/02-backend-go/architecture.md`
- `docs/05-infrastructure/README.md`
- `docs/06-features/mdm.md`

## Docs To Treat As Roadmap

Feature docs such as accounting, loan, CRM, HRM, BPM, notification, and system
config describe target scope unless they explicitly mention current
implementation.

`docs/superpowers` is historical migration-planning material and should not be
used as the current operating model.

## Proposals

1. Add a status block to every feature doc with one of: `Implemented`,
   `Prototype`, `Skeleton`, or `Roadmap`.
2. Rename `apps/backend-java/accounting_tmp` to `accounting` only when the Java
   service boundary is ready; otherwise retarget or disable Java CI.
3. Decide whether `crm-service` should remain in `go.work` while it is still a
   skeleton.
4. Add a lightweight docs link checker in CI for relative Markdown links.
5. Add a release checklist that updates docs whenever an app route, gateway
   route, service port, namespace, or image name changes.
6. Keep `.github/profile/README.md` high-level and link to `arda/docs` for
   detailed architecture instead of duplicating deep implementation details.
