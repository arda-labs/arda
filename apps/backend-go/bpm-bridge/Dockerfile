FROM golang:1.26-bookworm AS builder
WORKDIR /src
COPY libs/go/pkg ./libs/go/pkg
COPY apps/backend-go/bpm-bridge ./apps/backend-go/bpm-bridge
WORKDIR /src/apps/backend-go/bpm-bridge
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build -o /src/bin/bpm-bridge .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /src/bin/bpm-bridge /app/bpm-bridge
USER nonroot:nonroot
ENTRYPOINT ["/app/bpm-bridge"]
