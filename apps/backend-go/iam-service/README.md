# IAM Service

Identity and Access Management service for Arda.

## Responsibilities

- Zitadel login integration.
- Current-user profile and auth settings.
- Tenant/workspace creation and membership.
- Users, roles, groups, permissions, and resource permissions.
- Menu management and tenant-aware menu output.
- Forward-auth endpoint for APISIX integration.

## Local Run

The default local config uses:

```text
postgres://iam:iam%40123@thinkcenter:5432/iam?sslmode=disable
```

Start the service from this directory:

```powershell
kratos run
```

or:

```powershell
go run ./cmd/iam-service -conf ./configs
```

Local ports:

- HTTP: `8000`
- gRPC: `9000`

## API

Native HTTP routes are registered under `/v1/*`. APISIX exposes them as
`/api/v1/*`.

Important routes:

- `/v1/auth/login`
- `/v1/auth/settings`
- `/v1/auth/forward`
- `/v1/me`
- `/v1/me/tenants`
- `/v1/me/menu`
- `/v1/users`
- `/v1/tenants`
- `/v1/roles`
- `/v1/groups`
- `/v1/permissions`

## Database

Migrations are embedded and run during startup from:

```text
internal/data/migrations
```

The current migrations create IAM schema, seed initial data, add fine-grained
access-control tables, tenant deployment boundaries, workspace cleanup, and MDM
menu entries.

## Generate

Regenerate protobuf outputs after editing proto files:

```powershell
protoc --proto_path=. --proto_path=./third_party --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/iam/v1/iam.proto
protoc --proto_path=. --proto_path=./third_party --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/iam/v1/menu.proto
protoc --proto_path=. --go_out=paths=source_relative:. internal/conf/conf.proto
```

Regenerate Wire DI from the command directory when providers change:

```powershell
cd cmd\iam-service
wire
```

## Verify

```powershell
go test ./...
```

## Docker

From the repository root:

```powershell
docker build -f apps/backend-go/iam-service/Dockerfile -t ghcr.io/arda-labs/iam-service:dev .
```
