# Features: BPM Platform (Quy trình nghiệp vụ)

Dự án Arda cung cấp một nền tảng BPM toàn diện để quản lý, vận hành và giám sát quy trình doanh nghiệp.

## 📋 Overview

Hệ thống sử dụng **Camunda 8 (Zeebe)** làm orchestration engine, kết hợp với **BPM Service (Go)** đóng vai trò là Bridge để cung cấp các tính năng nâng cao cho người dùng cuối.

## 📊 Danh sách Chức năng

### 1. Giám sát & Vận hành (Monitoring)

- **BPMN 2.0 Viewer**: Hiển thị sơ đồ quy trình trực quan.
- **Heatmap**: Hiển thị mật độ hồ sơ tại từng bước quy trình (Danger/Warn/Success) để phát hiện nút thắt cổ chai.
- **Kafka Stream Timeline**: Theo dõi dòng sự kiện thời gian thực từ các microservices.
- **Audit Log**: Truy vết chi tiết "Ai - Làm gì - Khi nào".

### 2. Hộp thư Công việc (Worklist)

- **Inbound (Giao dịch đến)**:
  - Danh sách hồ sơ chờ xử lý.
  - Phê duyệt hàng loạt (Batch Approval).
  - Xem chi tiết hồ sơ, lịch sử và đính kèm.
- **Outbound (Giao dịch đi)**:
  - Theo dõi hồ sơ đã gửi đang nằm tại bước nào.
  - Thu hồi (Recall) hồ sơ nếu chưa được xử lý.

### 3. Cấu hình & Quản trị (Configuration)

- **Rule Chia Bài (Assignment Rules)**:
  - Round Robin (Chia đều).
  - Load Balance (Theo tải/backlog).
  - Trọng số (Weighted assignment).
- **Quản lý SLA**:
  - Thiết lập thời gian xử lý cho từng bước.
  - Cấu hình ngưỡng cảnh báo (%) trước khi quá hạn.
- **Cấu trúc Diễn giải (Description Templates)**:
  - Quản lý các lý do mẫu cho Phê duyệt/Từ chối theo từng Module.

### 4. Xử lý lỗi (Error Hospital)

- **Dead Letter Queue (DLQ)**: Tập trung các service task bị lỗi kỹ thuật.
- **Payload Editor**: Cho phép sửa dữ liệu JSON đầu vào để "cứu" task.
- **Manual Retry**: Chạy lại các bước lỗi sau khi đã khắc phục nguyên nhân.

## 🛠️ Tech Stack

- **Engine**: Camunda 8 (Zeebe).
- **Backend Bridge**: Go (Kratos) + gRPC.
- **Frontend**: Angular 21 + PrimeNG 19 + bpmn-js.
- **Messaging**: Kafka (Real-time stream).
