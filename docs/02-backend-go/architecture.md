# Go Backend Architecture

Updated: 2026-04-30

The Go backend is a Kratos-based workspace for operational services. The
active services are IAM and MDM. CRM exists in the workspace but should be
treated as skeleton/roadmap until implemented.

## Workspace

```text
apps/backend-go/
├── go.work
├── iam-service/
├── mdm-service/
└── crm-service/

libs/go/pkg/
├── database/
├── middleware/
├── pagination/
└── redis/
```

`go.work` includes the shared Go package and each service module.

## Service Layout

Each Kratos service follows the same basic shape:

```text
<service>/
├── api/<domain>/v1/          # protobuf and generated HTTP/gRPC code
├── cmd/<service>/            # main, Wire providers
├── configs/config.yaml       # local config
├── internal/
│   ├── biz/                  # use cases and repository interfaces
│   ├── data/                 # repositories, migrations, DB access
│   ├── server/               # HTTP/gRPC server registration
│   ├── service/              # protobuf service implementation
│   └── conf/                 # generated config structs
└── Dockerfile
```

## Current Services

| Service | Status | Native route | Gateway route | Local DB |
| --- | --- | --- | --- | --- |
| `iam-service` | Active | `/v1/*` | `/api/v1/*` | `iam` |
| `mdm-service` | Active | `/v1/mdm/*` | `/api/v1/mdm/*` | `mdm` |
| `crm-service` | Skeleton | TBD | TBD | TBD |

The gateway rewrites `/api/<path>` to `/<path>` before traffic reaches a
service.

## Migrations

Services embed migrations under `internal/data/migrations` and run them during
startup.

Current important migrations:

- IAM seeds tenant/menu/permission data and MDM menu entries.
- MDM creates the master-data schema and seed lists for geography, codes,
  system parameters, and banking reference data.

## Configuration

Local config is YAML under each service's `configs` directory. In Kubernetes,
`arda-infra` renders config from ConfigMaps and Secrets before mounting it at
`/data/conf`.

MDM uses a dedicated database user:

```text
postgres://mdm:mdm%40123@thinkcenter:5432/mdm?sslmode=disable
```

IAM keeps its own database/user. Do not reuse the IAM DB user for new services
except for short-lived bootstrap work.

## Development

```powershell
cd apps\backend-go\iam-service
go test ./...
kratos run

cd ..\mdm-service
go test ./...
kratos run
```

Regenerate protobuf outputs after editing `.proto` files:

```powershell
protoc --proto_path=. --proto_path=./third_party --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/mdm/v1/mdm.proto
```

Regenerate Wire after changing provider graphs:

```powershell
cd cmd\mdm
wire
```

## Near-Term Work

- Formalize APISIX forward-auth integration with IAM.
- Decide whether `crm-service` should remain in the workspace before it is
  implemented.
- Add focused service tests around repository behavior and migration seeds.
- Keep service names three-letter where practical: `iam`, `mdm`, then future
  `crm`, `hrm`, etc.
