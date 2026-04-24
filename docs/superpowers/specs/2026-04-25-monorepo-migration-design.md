# Spec: Arda Monorepo Migration & Renaming

- **Topic**: Monorepo Migration & Organization Renaming
- **Date**: 2026-04-25
- **Author**: Claude (Antigravity)
- **Status**: Draft

## 1. Mục tiêu
Hợp nhất các repo riêng lẻ của dự án Arda vào một Monorepo duy nhất trên GitHub tổ chức `arda-labs`, đồng thời chuyển đổi toàn bộ tham chiếu từ `arda-labs` sang `arda-labs`.

## 2. Cấu trúc đích (Arda Monorepo)
Repo: `arda-labs/arda` (Public)
```text
arda/
├── docs/             # Tài liệu hệ thống (Giữ nguyên và bổ sung)
├── arda-mfe/         # Frontend (Nx Monorepo)
├── arda-be-go/       # Backend Go (Kratos)
├── arda-be-java/     # Backend Java (Spring Boot)
├── scripts/          # Script vận hành chung
├── .github/          # CI/CD Workflows
└── README.md         # Project entry point
```

## 3. Các bước thực hiện chi tiết

### Giai đoạn 1: Chuẩn bị Local & Đổi tên hàng loạt
1.  **Xóa lịch sử cũ**: Tìm và xóa tất cả các thư mục `.git` bên trong các thư mục con (`arda-mfe`, `arda-be-go`, v.v.).
2.  **Replace String**: Quét toàn bộ project và thay thế:
    - `arda-labs` -> `arda-labs`
    - `github.com/arda-labs` -> `github.com/arda-labs/arda` (cho các import Go)
3.  **Cập nhật cấu hình**:
    - Sửa `go.mod` trong các service Go để khớp với đường dẫn mới trong monorepo.
    - Cập nhật các file cấu hình Nx trong `arda-mfe`.

### Giai đoạn 2: Quản lý GitHub Repositories
1.  **Xóa repo cũ trên GitHub**:
    - `arda-labs/arda-mfe`
    - `arda-labs/iam-service`
    - `arda-labs/common-service`
2.  **Khởi tạo Monorepo mới**:
    - `git init` tại `d:/Github/arda`.
    - `git add .` (Loại trừ các file không cần thiết qua `.gitignore`).
    - `git commit -m "initial: merge all sub-repos into arda monorepo and rename to arda-labs"`
    - `gh repo create arda-labs/arda --public --source=. --remote=origin --push`

### Giai đoạn 3: Cập nhật Infrastructure
1.  Sửa các file cấu hình trong `arda-infra` (vốn đã là một repo riêng) để trỏ đến source code mới trong repo `arda`.
2.  Cập nhật đường dẫn Docker Image nếu có thay đổi.

## 4. Rủi ro và Giải pháp
- **Rủi ro**: Lỗi import trong Go do thay đổi module path.
- **Giải pháp**: Chạy `go mod tidy` cho từng service sau khi đổi tên.
- **Rủi ro**: CI/CD cũ bị hỏng.
- **Giải pháp**: Tạm thời disable các workflow cũ và migrate dần sang cấu trúc monorepo.

## 5. Tiêu chí hoàn thành
- [ ] Không còn từ khóa `arda-labs` trong toàn bộ codebase.
- [ ] Repo `arda-labs/arda` có đầy đủ code của MFE, BE Go, BE Java và Docs.
- [ ] Các service Go có thể build local thành công với module path mới.
- [ ] GitHub Org `arda-labs` đã được dọn dẹp sạch sẽ các repo cũ.
