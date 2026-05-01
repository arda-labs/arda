# BPM Bridge & Advanced Features Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Xây dựng BPM Bridge service để hỗ trợ SLA, Custom Dashboard và Monitoring cho dự án Arda.

**Architecture:** Bridge service kết nối Zeebe (vận hành) và CRM DB (nghiệp vụ).

**Tech Stack:** Go (hoặc Kotlin), Zeebe Client, Redis (cho SLA caching).

---

### Task 1: Initialize BPM Bridge Project

**Files:**
- Create: `arda/apps/backend-go/bpm-bridge/go.mod`

- [ ] **Step 1: Init project**

```bash
cd arda/apps/backend-go/bpm-bridge
go mod init github.com/arda-labs/arda/apps/backend-go/bpm-bridge
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-go/bpm-bridge/
git commit -m "feat(bpm-bridge): init bpm bridge project"
```

---

### Task 2: Implement SLA Tracking Logic

**Files:**
- Create: `arda/apps/backend-go/bpm-bridge/internal/sla/sla_manager.go`

- [ ] **Step 1: Implement SLA calculator**

```go
func CalculateSLAExpiry(startTime time.Time, durationHours int) time.Time {
    // Logic accounting for business hours/holidays
}
```

- [ ] **Step 2: Commit**

```bash
git add arda/apps/backend-go/bpm-bridge/internal/sla/
git commit -m "feat(bpm-bridge): add sla calculation logic"
```
