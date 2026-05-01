# Design Spec: Arda CRM & BPM Ecosystem (Production-Ready)

**Status**: Draft (Updated with Expert Review & Advanced BPM Features)
**Date**: 2026-05-01
**Topic**: CRM Service (Kotlin), Camunda 8 (BPM), Media Service (Go), BPM Bridge (Advanced Operations)

## 1. Overview
Hệ sinh thái này cung cấp khả năng quản lý quy trình nghiệp vụ (BPM) cấp độ Enterprise và lưu trữ tài liệu (Media). Thiết kế này kế thừa và nâng cấp các tính năng từ EPAS (SLA, Custom Inbox, Monitoring) lên nền tảng Camunda 8.

## 2. Architecture & Patterns

### 2.1 BPM Bridge Service (New Component)

Để đáp ứng các yêu cầu về Dashboard và SLA nâng cao mà Camunda 8 chuẩn chưa hỗ trợ linh hoạt, chúng ta bổ sung **BPM Bridge Service**:

- **Aggregation**: Gom dữ liệu từ Zeebe (vận hành) và Database nghiệp vụ.
- **SLA Management**: Tính toán thời gian hết hạn (Overdue) dựa trên Business Calendar.
- **Custom Task API**: Cung cấp API cho Frontend hiển thị Inbox/Outbox tùy biến, hỗ trợ Batch Approval.

### 2.2 Communication Model (Zeebe Pull Model)

- **Workers**: CRM và các utility services poll jobs từ Zeebe Broker.
- **Idempotency**: Đảm bảo các worker xử lý trùng lặp an toàn bằng cách kiểm tra trạng thái record trong DB.

### 2.3 Transactional Consistency (Outbox Pattern)

- **Outbox**: Đảm bảo atomicity giữa lưu trữ nghiệp vụ và kích hoạt workflow.

## 3. Detailed Design

### 3.1 Media Service (Go)

- **Storage**: SeaweedFS (S3-compatible).
- **Features**: MIME validation, Image resizing, Presigned URL, File Lifecycle Management.

### 3.2 CRM Service (Kotlin)

- **Features**: Quản lý Customer (Basic, Identity, Address, Relation).
- **Integration**: `@JobWorker` cho các bước nghiệp vụ.

### 3.3 Advanced BPM Features (Inspired by EPAS)

- **SLA Process**: Định nghĩa và theo dõi SLA cho từng bước trong quy trình.
- **Process Monitor**: Giao diện theo dõi trạng thái chi tiết của từng Instance (Audit log, Variable changes, Incidents).
- **Decision Rules**: Sử dụng Camunda DMN cho các quy trình ra quyết định tự động.

## 4. Infrastructure Scaling (Prod vs Dev)

| Thành phần | Dev Mode (K3s/Local) | Production Mode |
| --- | --- | --- |
| **Zeebe Broker** | 1 Node | 3 Nodes (Cluster) |
| **Search Engine** | Elasticsearch (mini) | Elasticsearch (Full Cluster) |
| **BPM Bridge** | Single Instance | Multi-instance (H.A) |

## 5. Security & Compliance
- **Audit Logging**: Ghi log chi tiết "Ai - Làm gì - Khi nào" cho cả dữ liệu nghiệp vụ và thay đổi workflow.
- **Media Security**: Audit log cho hành động xem file nhạy cảm.

## 6. Workflow Example: Advanced Customer Onboarding
1. **Frontend** -> **Media Service**: Upload CCCD.
2. **Frontend** -> **CRM Service**: Submit Form.
3. **CRM Service**: Save + Outbox.
4. **Zeebe**: Start Process.
5. **BPM Bridge**: Bắt đầu tính toán SLA cho bước Phê duyệt.
6. **Admin**: Xem danh sách task qua **Custom Inbox (via Bridge)**.
7. **Admin**: Phê duyệt (Batch or Single).
8. **Zeebe**: Cập nhật kết quả qua CRM Worker.
