# Media Service (Go) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Xây dựng service quản lý file tập trung cho Arda bằng Go, tích hợp SeaweedFS S3 API.

**Architecture:** Service sử dụng `pgx` cho database và `aws-sdk-go-v2` cho S3. Hỗ trợ validation và presigned URL.

**Tech Stack:** Go 1.25, PostgreSQL, S3 SDK, Gin/Echo (HTTP).

---

### Task 1: Initialize Media Service project

**Files:**
- Create: `arda/apps/backend-go/media-service/go.mod`
- Create: `arda/apps/backend-go/media-service/main.go`

- [ ] **Step 1: Create go.mod**

```go
module github.com/arda-labs/arda/apps/backend-go/media-service

go 1.25
```

- [ ] **Step 2: Create basic main.go**

```go
package main

import "fmt"

func main() {
    fmt.Println("Media Service starting...")
}
```

- [ ] **Step 3: Commit**

```bash
git add arda/apps/backend-go/media-service/
git commit -m "feat(media): init media service project"
```

---

### Task 2: Implement File Metadata Model & DB Migration

**Files:**
- Create: `arda/apps/backend-go/media-service/internal/data/models.go`
- Create: `arda/apps/backend-go/media-service/migrations/000001_init_media.up.sql`

- [ ] **Step 1: Create migration script**

```sql
CREATE TABLE media_metadata (
    id UUID PRIMARY KEY,
    filename TEXT NOT NULL,
    content_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    s3_key TEXT NOT NULL,
    owner_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-go/media-service/migrations/
git commit -m "feat(media): add metadata table migration"
```

---

### Task 3: Implement S3 Client & Upload Logic

**Files:**
- Create: `arda/apps/backend-go/media-service/internal/storage/s3.go`

- [ ] **Step 1: Write S3 initialization with SeaweedFS endpoint**

```go
// Mock/Code snippet for implementation
func NewS3Client(endpoint string) *s3.Client {
    // Config with custom endpoint for SeaweedFS
}
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-go/media-service/internal/storage/
git commit -m "feat(media): implement s3 client for seaweedfs"
```
