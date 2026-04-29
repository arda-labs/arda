# Frontend Architecture

Updated: 2026-04-30

The current frontend is an Angular CLI workspace using Native Federation. It is
not an Nx workspace, and it does not currently use Analog SSR or Rspack.

## Workspace

```text
apps/frontend-micro/
├── angular.json
├── package.json
├── Dockerfile
├── nginx.conf
└── projects/
    ├── shell/
    ├── iam/
    ├── mdm/
    └── core/
```

| Project | Type | Default dev port | Role |
| --- | --- | --- | --- |
| `shell` | Host app | `3000` | Login callback, layout, workspace selection, remote loading |
| `iam` | Remote MFE | `3002` | IAM screens |
| `mdm` | Remote MFE | `3001` | Master-data screens |
| `core` | Library | N/A | Shared Angular library placeholder |

## Federation

The shell uses `@angular-architects/native-federation` and builds its manifest
at runtime from `window.__env`.

```ts
const federationManifest = {
  iam: `${env.mfeIamUrl || 'http://localhost:3002'}/remoteEntry.json`,
  mdm: `${env.mfeMdmUrl || 'http://localhost:3001'}/remoteEntry.json`
};
```

Each remote exposes `./Routes` from its `federation.config.js`. The shell then
loads them lazily:

| Shell route | Remote |
| --- | --- |
| `/iam/*` | `iam` |
| `/mdm/*` | `mdm` |

In deployed or APISIX-local mode, remotes are served through:

| Asset route | Service |
| --- | --- |
| `/mfe-iam/*` | `mfe-iam` |
| `/mfe-mdm/*` | `mfe-mdm` |

## Runtime Configuration

`projects/shell/public/env.js` is loaded before Angular. For APISIX-local
development, use:

```js
window.__env.apiUrl = 'http://localhost:9080/api';
window.__env.apiPath = '/v1';
window.__env.mfeIamUrl = 'http://localhost:9080/mfe-iam';
window.__env.mfeMdmUrl = 'http://localhost:9080/mfe-mdm';
```

This keeps browser traffic aligned with deployed route shapes.

## UI And State

- Angular 21 standalone components.
- PrimeNG 21 for enterprise controls.
- Tailwind CSS 4 for layout and utility styling.
- Signals for local state where they simplify the component.
- Reactive forms for create/update dialogs and filter forms.
- `p-table` is the default grid until requirements justify AG Grid Enterprise.

## Current Feature Ownership

Shell owns:

- public landing route `/`;
- login and auth callback routes;
- tenant/workspace selection and creation;
- shared layout, header, sidebar, and error pages;
- loading remote MFEs.

IAM remote owns:

- IAM management screens and tenant-scoped IAM operations.

MDM remote owns:

- administrative geography;
- area types and areas;
- code sets and code items;
- system parameters;
- banking reference-data view.

## Development

```powershell
cd apps\frontend-micro
npm install
npx ng serve shell
npx ng serve iam
npx ng serve mdm
```

Build:

```powershell
npx ng build shell
npx ng build iam
npx ng build mdm
```

Docker images are built from the shared `apps/frontend-micro/Dockerfile` using
targets such as `mfe-shell-runtime`, `mfe-iam-runtime`, and `mfe-mdm-runtime`.

## Design Rules

- Keep the shell focused on navigation, auth, tenancy, and composition.
- Keep domain screens in their remote MFE.
- Prefer APISIX-local checks over direct service calls when testing browser
  integration, auth headers, CORS, and path rewrites.
- Do not introduce a shared UI abstraction until at least two MFEs need it and
  the interaction contract is stable.

See also:

- [Tenant Creation](./tenant-creation.md)
- [UI/UX Data Grid Strategy](./ui-ux-data-grid.md)
