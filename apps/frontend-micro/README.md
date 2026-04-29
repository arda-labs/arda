# Frontend Micro

Angular 21 workspace for Arda micro-frontends.

## Projects

| Project | Type | Dev port | Purpose |
| --- | --- | --- | --- |
| `shell` | Host app | `3000` | Layout, auth callback, workspace UI, remote loading |
| `iam` | Remote MFE | `3002` | IAM screens |
| `mdm` | Remote MFE | `3001` | Master Data Management screens |
| `core` | Library | N/A | Shared Angular library placeholder |

The workspace uses Angular CLI and `@angular-architects/native-federation`.
It is not an Nx workspace.

## Install

```powershell
npm install
```

## Run

```powershell
npx ng serve shell
npx ng serve iam
npx ng serve mdm
```

For production-parity browser checks, open the app through local APISIX from
`arda-infra/local/apisix` instead of calling services directly:

```text
http://localhost:9080
```

## Build

```powershell
npx ng build shell
npx ng build iam
npx ng build mdm
```

## Runtime Federation

The shell reads remote URLs from `projects/shell/public/env.js`:

```js
window.__env.mfeIamUrl = 'http://localhost:9080/mfe-iam';
window.__env.mfeMdmUrl = 'http://localhost:9080/mfe-mdm';
```

Local direct defaults are:

- IAM remote: `http://localhost:3002/remoteEntry.json`
- MDM remote: `http://localhost:3001/remoteEntry.json`

## Docker

The shared Dockerfile builds individual runtime targets:

```powershell
docker build --target mfe-shell-runtime --build-arg PROJECT=shell -t ghcr.io/arda-labs/mfe-shell:dev .
docker build --target mfe-iam-runtime --build-arg PROJECT=iam -t ghcr.io/arda-labs/mfe-iam:dev .
docker build --target mfe-mdm-runtime --build-arg PROJECT=mdm -t ghcr.io/arda-labs/mfe-mdm:dev .
```
