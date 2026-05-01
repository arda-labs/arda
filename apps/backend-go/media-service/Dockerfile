FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/media-service ./apps/backend-go/media-service
WORKDIR /src/apps/backend-go/media-service
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -o /src/bin/media-service .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /src/bin/media-service /app/media-service
COPY --from=builder /src/apps/backend-go/media-service/configs /app/configs
USER nonroot:nonroot
ENTRYPOINT ["/app/media-service", "-conf", "/app/configs"]
