FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

WORKDIR /src

COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/notification-service ./apps/backend-go/notification-service

WORKDIR /src/apps/backend-go/notification-service

RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o /src/bin/notification-service \
    ./cmd/notification-service/

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/apps/backend-go/notification-service/configs /data/conf
COPY --from=builder /src/bin/notification-service /app/notification-service

EXPOSE 8002
EXPOSE 9002

USER nonroot:nonroot

ENTRYPOINT ["/app/notification-service"]
CMD ["-conf", "/data/conf"]
