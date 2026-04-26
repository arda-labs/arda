# Build stage: golang:1.26 (Debian bookworm)
FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git make && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# 1. Chép thư viện dùng chung (Bắt buộc)
COPY libs/go/pkg ./libs/go/pkg

# 2. Chép mã nguồn của iam-service
COPY apps/backend-go/iam-service ./apps/backend-go/iam-service

# 3. Build độc lập (Tắt GOWORK)
# Chúng ta sẽ di chuyển vào thư mục service và build trực tiếp
WORKDIR /src/apps/backend-go/iam-service

# Tắt GOWORK để không bị phụ thuộc vào các service khác trong workspace
# Go sẽ tự tìm các dependency thông qua go.mod (nếu có replace directive)
RUN CGO_ENABLED=0 GOOS=linux GOWORK=off go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o /src/bin/iam-service \
    ./cmd/iam-service/

# Runtime stage: distroless/static
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /src/apps/backend-go/iam-service/configs /data/conf
COPY --from=builder /src/bin/iam-service /app/iam-service

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/iam-service"]
CMD ["-conf", "/data/conf"]
