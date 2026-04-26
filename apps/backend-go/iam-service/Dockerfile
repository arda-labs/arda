# Build stage: golang:1.25 (Debian bookworm)
FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git make && rm -rf /var/lib/apt/lists/*

# Chép toàn bộ backend workspace vào để hỗ trợ go.work
COPY . /src
WORKDIR /src

# Build iam-service từ root workspace
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o /src/bin/iam-service \
    ./iam-service/cmd/iam-service/

# Runtime stage: distroless/static
FROM gcr.io/distroless/static-debian12:nonroot

# Copy configs từ thư mục service tương ứng
COPY --from=builder /src/iam-service/configs /data/conf
COPY --from=builder /src/bin/iam-service /app/iam-service

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/iam-service"]
CMD ["-conf", "/data/conf"]
