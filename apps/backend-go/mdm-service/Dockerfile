FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/mdm-service ./apps/backend-go/mdm-service

WORKDIR /src/apps/backend-go/mdm-service

RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o /src/bin/mdm-service \
    ./cmd/mdm/

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/apps/backend-go/mdm-service/configs /data/conf
COPY --from=builder /src/bin/mdm-service /app/mdm-service

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/mdm-service"]
CMD ["-conf", "/data/conf"]
