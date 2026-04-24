# Build stage: golang:1.25 (Debian bookworm) - Go stdlib không có CVE
# Alpine chưa có Go 1.25 image, dùng bookworm cho builder
FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git make && rm -rf /var/lib/apt/lists/*

COPY . /src
WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o ./bin/iam-service \
    ./cmd/iam-service/

# Runtime stage: distroless/static - không có shell, không có OS packages
# Zero CVE từ OS layer, chỉ có binary và ca-certificates
FROM gcr.io/distroless/static-debian12:nonroot

# configs chứa template với ${placeholders}
# K8s: bị override bởi volume mount (initContainer envsubst)
# Local Docker: mount config thật qua -v hoặc truyền env vars
COPY --from=builder /src/configs /data/conf
COPY --from=builder /src/bin/iam-service /app/iam-service

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/iam-service"]
CMD ["-conf", "/data/conf"]
