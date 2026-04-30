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

### Administrative Unit Sync

Vietnam province and ward data is no longer seeded from a static SQL dump. The
service syncs the current 2-level administrative catalog from CASSO AddressKit:

```text
https://production.cas.so/address-kit/latest/provinces
https://production.cas.so/address-kit/latest/communes
```

Run the sync through the internal MDM API:

```bash
curl -X POST http://localhost:8001/v1/mdm/administrative-units/sync-addresskit
```

Through the gateway, call:

```bash
curl -X POST http://localhost:8001/api/v1/mdm/administrative-units/sync-addresskit
```

The endpoint deletes current administrative unit data and loads the latest
AddressKit response in one database transaction. It also clears dependent area
assignments and administrative unit mappings because their foreign keys point to
the replaced unit rows.

The MDM frontend exposes the same operation on the `Tỉnh, phường/xã` page via
the `Đồng bộ AddressKit` button. The button calls the internal API above; the
browser does not call AddressKit directly.

## Database

Migrations are embedded and run during service startup from:

```text
internal/data/migrations
```

The seed migration includes shared reference lists such as administrative
levels, statuses, system parameter groups, currencies, countries, bank account
types, payment methods, payment channels, transaction types, customer segments,
risk ratings, document types, interest rate types, collateral types, and fee
types. Migration `000003` intentionally clears any historical static Vietnam
administrative unit seed; use the AddressKit sync endpoint to load that data.

## Generate

Regenerate protobuf outputs after editing proto files:

```bash
protoc --proto_path=./api --proto_path=./third_party --go_out=paths=source_relative:./api --go-http_out=paths=source_relative:./api --go-grpc_out=paths=source_relative:./api --openapi_out=fq_schema_naming=true,default_response=false:. ./api/mdm/v1/mdm.proto
protoc --proto_path=./internal --proto_path=./third_party --go_out=paths=source_relative:./internal internal/conf/conf.proto
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
