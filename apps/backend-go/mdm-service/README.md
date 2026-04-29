# MDM Service

Master Data Management service for shared reference data used by IAM, CRM,
HRM, accounting, banking, payments, and reporting.

## Local Run

The default local config uses a dedicated MDM database:

```text
postgres://mdm:mdm%40123@thinkcenter:5432/mdm?sslmode=disable
```

Start the service from this directory:

```bash
kratos run
```

or from the service command:

```bash
go run ./cmd/mdm -conf ./configs
```

Local ports from `configs/config.yaml`:

- HTTP: `8001`
- gRPC: `9001`

## API

HTTP routes are registered under `/v1/mdm/*`. Gateway deployments rewrite
`/api/v1/mdm/*` to the service route.

Main resources:

- Administrative units: provinces/cities, wards/communes, and hierarchical
  administrative levels.
- Area types and areas.
- Code sets and code items for reusable reference lists.
- System parameters for runtime configuration.

## Database

Migrations are embedded and run during service startup from:

```text
internal/data/migrations
```

The seed migration includes administrative levels, status lists, system
parameter groups, currencies, countries, bank account types, payment methods,
payment channels, transaction types, customer segments, risk ratings, document
types, interest rate types, collateral types, and fee types.

## Generate

Regenerate protobuf outputs after editing proto files:

```bash
protoc --proto_path=. --proto_path=./third_party --go_out=paths=source_relative:. --go-http_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. api/mdm/v1/mdm.proto
protoc --proto_path=. --go_out=paths=source_relative:. internal/conf/conf.proto
```

Regenerate Wire DI from `cmd/mdm` when the provider graph changes:

```bash
wire
```

## Verify

```bash
go test ./...
```

## Docker

The production container reads `/data/conf/config.yaml`; Kubernetes renders it
from `arda-infra/apps/mdm-service/base/configs.yaml`.

```bash
docker build -t ghcr.io/arda-labs/mdm-service:dev .
```
