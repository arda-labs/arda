# Arda Platform: Technical Standards (Java & Go)

> Phiên bản: 1.0
> Cập nhật: 2026-04-25
> Tài liệu này quy định các tiêu chuẩn kỹ thuật để đảm bảo tính đồng nhất giữa hệ thống Backend Java (Core Banking) và Go (Operational Services).

---

## 1. Package & Project Structure

### Java (Spring Boot / Kotlin)
- **Root Package**: `arda.*` (Ví dụ: `arda.common`, `arda.accounting`).
- **Thư mục**: Rút gọn, không sử dụng `com.arda.labs` để giảm độ sâu.
- **Build Tool**: Gradle Kotlin DSL.

### Go (Kratos)
- **Module**: `github.com/arda-labs/arda-be-go`.
- **Cấu trúc**: Tuân thủ [Kratos Layout](https://go-kratos.dev/en/docs/getting-started/layout).

---

## 2. Communication Protocol (gRPC & REST)

### 2.1. API-First Design
- Mọi API giao tiếp giữa các service (Inter-service) **BẮT BUỘC** sử dụng gRPC.
- Sử dụng Protobuf làm Contract. File `.proto` đặt tại thư mục chung của hệ thống.

### 2.2. Error Modeling (Rich Error Model)
Tất cả các service phải trả về lỗi theo format chuẩn:
- **Code**: String code duy nhất (Ví dụ: `ERR_INTERNAL`, `ERR_INSUFFICIENT_BALANCE`).
- **Message**: Mô tả lỗi thân thiện.
- **Details**: Danh sách các lỗi chi tiết (field-level validation).

**Java Mapping:**
```kotlin
data class ApiError(val code: String, val message: String, val details: List<ErrorDetail>?)
```

**Go Mapping (Kratos):**
Sử dụng `errors.New()` với các `details` được bọc trong Protobuf messages.

---

## 3. Data Standardization

### 3.1. Money & Currency
Để tránh sai lệch làm tròn giữa các ngôn ngữ:
- **Lưu trữ/Tính toán**: Sử dụng `BigDecimal` (Java) và `shopspring/decimal` (Go).
- **Truyền nhận (API)**: BẮT BUỘC sử dụng kiểu `String`.
- **Quy định**: Luôn có 2 chữ số thập phân (Scale = 2) và làm tròn `HALF_UP`.

### 3.2. Date & Time
- **Định dạng**: RFC 3339 / ISO-8601 (Ví dụ: `2026-04-25T15:30:00Z`).
- **Timezone**: Luôn sử dụng UTC.
- **Java**: `java.time.Instant`.
- **Go**: `time.Time`.

### 3.3. Pagination
Cấu trúc Response phân trang đồng nhất:
```json
{
  "items": [...],
  "metadata": {
    "page": 0,
    "size": 10,
    "totalElements": 100,
    "totalPages": 10,
    "hasNext": true
  }
}
```

---

## 4. Observability & Context

### 4.1. Context Propagation
Sử dụng chuẩn **W3C TraceContext** (`traceparent` header).
- **Java**: Sử dụng `ArdaContext` tích hợp với Reactor Context và Micrometer.
- **Go**: Sử dụng `context.Context` mặc định của Go kết hợp với Kratos metadata/tracing middleware.

### 4.2. Baggage
Các thông tin định danh người dùng (`user_id`, `tenant_id`, `roles`) được truyền qua metadata thay vì giải mã JWT nhiều lần.

---

## 5. Coding Style
- **Java/Kotlin**: [Google Java Style](https://google.github.io/styleguide/javaguide.html).
- **Go**: [Uber Go Style Guide](https://github.com/uber-go/guide).
- **Git**: [Conventional Commits](https://www.conventionalcommits.org/).
