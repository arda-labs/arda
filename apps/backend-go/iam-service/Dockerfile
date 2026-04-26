# Build stage: golang:1.26 (Debian bookworm)
FROM golang:1.26-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends git make && rm -rf /var/lib/apt/lists/*

WORKDIR /src

# Thiết lập GOTOOLCHAIN để tự động tải version nếu cần
ENV GOTOOLCHAIN=auto

# Docker context là thư mục gốc của monorepo (arda/)
# Chép file quản lý workspace và thư viện chung
COPY apps/backend-go/go.work apps/backend-go/go.work.sum ./apps/backend-go/
COPY libs/go/pkg ./libs/go/pkg

# Chép mã nguồn của iam-service
COPY apps/backend-go/iam-service ./apps/backend-go/iam-service

# Thiết lập môi trường để build đúng context workspace
WORKDIR /src/apps/backend-go

# Build iam-service sử dụng go.work đã được copy
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-X main.Version=$(git describe --tags --always 2>/dev/null || echo dev) -w -s" \
    -o /src/bin/iam-service \
    ./iam-service/cmd/iam-service/

# Runtime stage: distroless/static
FROM gcr.io/distroless/static-debian12:nonroot

# Copy configs từ đúng đường dẫn
COPY --from=builder /src/apps/backend-go/iam-service/configs /data/conf
COPY --from=builder /src/bin/iam-service /app/iam-service

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/app/iam-service"]
CMD ["-conf", "/data/conf"]
