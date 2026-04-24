# Arda Monorepo Migration & Renaming Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Chuyển đổi cấu trúc project sang Monorepo duy nhất (`arda-labs/arda`) và đổi tên toàn bộ từ `arda-labs` sang `arda-labs`.

**Architecture:** 
1. Hợp nhất `arda-mfe`, `arda-be-go`, `arda-be-java` và `docs` vào một git root duy nhất tại `d:/Github/arda`.
2. Thay thế chuỗi văn bản hàng loạt và cập nhật cấu hình module (Go, Nx).
3. Sử dụng GitHub CLI để dọn dẹp các repo cũ và khởi tạo repo mới.

**Tech Stack:** Git, GitHub CLI (gh), Sed/Grep (cho renaming), Go Modules, Nx.

---

### Task 1: Dọn dẹp Local Git & Khởi tạo Monorepo

**Files:**
- Modify: `.gitignore`
- Delete: `**/.git` (thư mục ẩn)

- [ ] **Step 1: Xóa các thư mục .git con**
Run: `find . -mindepth 2 -name ".git" -type d -exec rm -rf {} +` (Cẩn thận: Chỉ chạy từ root `d:/Github/arda`)

- [ ] **Step 2: Cấu hình .gitignore tại root**
```bash
cat <<EOF > .gitignore
# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Language/Framework
node_modules/
dist/
bin/
vendor/
*.log
.nx/cache
EOF
```

- [ ] **Step 3: Khởi tạo Git Root**
Run: `git init && git add . && git commit -m "initial: prepare monorepo structure"`

---

### Task 2: Đổi tên hàng loạt (Renaming)

**Files:**
- Modify: Tất cả các file chứa "arda-labs"

- [ ] **Step 1: Chạy script thay thế chuỗi "arda-labs" thành "arda-labs"**
Run: `grep -rli 'arda-labs' . | xargs -i sed -i 's/arda-labs/arda-labs/g' {}`

- [ ] **Step 2: Kiểm tra lại các file quan trọng**
Run: `grep -r "arda-labs" .` (Kỳ vọng: Không còn kết quả)

- [ ] **Step 3: Cập nhật Go Modules**
Run: 
```bash
cd arda-be-go/services/iam-service && go mod edit -module github.com/arda-labs/arda/arda-be-go/services/iam-service && go mod tidy
cd ../common-service && go mod edit -module github.com/arda-labs/arda/arda-be-go/services/common-service && go mod tidy
```

- [ ] **Step 4: Commit thay đổi renaming**
Run: `git add . && git commit -m "refactor: rename arda-labs to arda-labs across project"`

---

### Task 3: Dọn dẹp GitHub & Đẩy Monorepo lên Cloud

**Files:**
- Action: GitHub API/CLI

- [ ] **Step 1: Xóa các repo cũ trên GitHub (arda-labs org)**
Run: 
```bash
gh repo delete arda-labs/arda-mfe --confirm
gh repo delete arda-labs/iam-service --confirm
gh repo delete arda-labs/common-service --confirm
```

- [ ] **Step 2: Tạo repo Monorepo mới**
Run: `gh repo create arda-labs/arda --public --source=. --remote=origin --push`

---

### Task 4: Cập nhật Infrastructure (arda-infra)

**Files:**
- Modify: `arda-infra/**/*.yaml`

- [ ] **Step 1: Sửa file root-app.yaml hoặc các Application config**
Thay thế `repoURL: https://github.com/arda-labs/iam-service` thành `repoURL: https://github.com/arda-labs/arda` và thêm `path: arda-be-go/services/iam-service`.

- [ ] **Step 2: Kiểm tra và Commit arda-infra**
Run: `cd arda-infra && git add . && git commit -m "chore: update source repo paths to arda monorepo"`
