# Core Angular Library

`core` is the shared Angular library placeholder inside
`apps/frontend-micro`.

Use it only for stable cross-MFE utilities or UI primitives. Do not move
domain-specific IAM or MDM behavior here.

## Build

```powershell
cd apps\frontend-micro
npx ng build core
```

## Guidance

- Prefer project-local code until at least two MFEs need the same behavior.
- Keep library APIs small and stable.
- Avoid adding shell-only concerns, auth-flow logic, or domain CRUD services to
  this library.
