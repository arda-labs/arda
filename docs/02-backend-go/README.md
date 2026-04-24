# Operational Services — Dịch vụ Vận hành

> Microservices Go với framework Kratos cho tốc độ và hiệu năng cao
> Xử lý các nghiệp vụ vận hành, tích hợp linh hoạt

---

## 📋 Overview

Nhóm Operational Services ưu tiên tốc độ xử lý nhanh và tích hợp linh hoạt. Tất cả services được xây dựng bằng **Go** sử dụng framework **Kratos** với **gRPC** cho inter-service communication.

### Services
1. **CRM & Member Service** — Quản lý khách hàng, hội viên, chấm điểm tín dụng
2. **HRM Service** — Cơ cấu tổ chức, thông tin nhân sự
3. **Notification Service** — Thông báo Zalo, Email, Push
4. **System Config Service** — Tham số hệ thống, ngày làm việc, đa ngôn ngữ
5. **BPM Engine Service** — Wrapper điều phối Camunda

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Operational Layer                          │
│                      (Go + Kratos)                           │
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │   CRM    │  │   HRM    │  │ Notification│ │  System  │  │
│  │  Service │  │ Service  │  │  Service   │ │  Config  │  │
│  │          │  │          │  │            │ │ Service  │  │
│  │ • Customer│  │ • Org    │  │ • Zalo     │  │ • Params │  │
│  │ • Member  │  │ • Employee│  │ • Email    │  │ • Calendar│  │
│  │ • Credit  │  │ • Position│  │ • Template │  │ • i18n   │  │
│  │   Score  │  │          │  │ • Queue    │  │          │  │
│  └────┬─────┘  └────┬─────┘  └──────┬─────┘  └─────┬────┘  │
│       │             │                │               │       │
│       └─────────────┼────────────────┼───────────────┘       │
│                     ▼                ▼                       │
│              ┌──────────────┐  ┌──────────┐                  │
│              │  IAM Service │  │  Redpanda│                  │
│              │  (AuthZ)     │  │  Events  │                  │
│              └──────────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

---

## 📁 Monorepo Structure

```
arda-be/
├── api/                             # Shared .proto definitions
│   ├── crm/v1/
│   │   ├── customer.proto
│   │   ├── member.proto
│   │   └── credit_score.proto
│   ├── hrm/v1/
│   │   ├── organization.proto
│   │   ├── employee.proto
│   │   └── position.proto
│   ├── notification/v1/
│   │   ├── notification.proto
│   │   ├── template.proto
│   │   └── channel.proto
│   └── system-config/v1/
│       ├── parameter.proto
│       ├── calendar.proto
│       └── i18n.proto
│
├── pkg/                             # Shared libraries
│   ├── auth/
│   │   ├── jwt.go                   # JWT utilities
│   │   └── zitadel.go               # Zitadel integration
│   ├── database/
│   │   ├── postgres.go              # PostgreSQL helpers
│   │   ├── tenant.go                # Tenant context
│   │   └── repository.go            # Base repository
│   ├── redis/
│   │   ├── client.go                # Redis client
│   │   └── cache.go                 # Cache utilities
│   ├── middleware/
│   │   ├── auth.go                  # Auth middleware
│   │   ├── tenant.go                # Tenant middleware
│   │   └── logging.go               # Logging middleware
│   └── event/
│       ├── publisher.go             # Event publisher
│       └── subscriber.go            # Event subscriber
│
├── crm-service/                     # CRM & Member Service
│   ├── cmd/crm-service/
│   │   └── main.go
│   ├── internal/
│   │   ├── biz/
│   │   │   ├── customer.go
│   │   │   ├── member.go
│   │   │   ├── credit_score.go
│   │   │   └── biz.go
│   │   ├── data/
│   │   │   ├── data.go
│   │   │   ├── customer.go
│   │   │   ├── member.go
│   │   │   └── credit_score.go
│   │   ├── service/
│   │   │   ├── crm.go
│   │   │   └── service.go
│   │   ├── server/
│   │   │   ├── grpc.go
│   │   │   ├── http.go
│   │   │   └── server.go
│   │   └── conf/
│   │       └── conf.pb.go
│   ├── api/crm/v1/
│   │   ├── customer.proto
│   │   ├── member.proto
│   │   └── credit_score.proto
│   ├── configs/
│   │   └── config.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   └── Dockerfile
│
├── hrm-service/                     # HRM Service
│   └── ... (similar structure)
│
├── notification-service/            # Notification Service
│   └── ... (similar structure)
│
├── system-config-service/           # System Config Service
│   └── ... (similar structure)
│
└── bpm-service/                     # BPM Engine Service
    └── ... (similar structure)
```

---

## 1. CRM & Member Service

### Responsibilities
- Quản lý hồ sơ khách hàng (multi-tenant)
- Quản lý hội viên, cấp độ hội viên
- Chấm điểm tín dụng
- Tích hợp với IAM service

### Domain Models

```go
// Customer
type Customer struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    CustomerCode string   `json:"customer_code" db:"customer_code"`
    FullName    string    `json:"full_name" db:"full_name"`
    Email       string    `json:"email" db:"email"`
    Phone       string    `json:"phone" db:"phone"`
    DateOfBirth time.Time `json:"date_of_birth" db:"date_of_birth"`
    IDNumber    string    `json:"id_number" db:"id_number"`
    IDType      string    `json:"id_type" db:"id_type"`
    Address     string    `json:"address" db:"address"`
    CustomerType string   `json:"customer_type" db:"customer_type"` // INDIVIDUAL, CORPORATE
    Status      string    `json:"status" db:"status"` // ACTIVE, INACTIVE, BLOCKED
    CreditScore int       `json:"credit_score" db:"credit_score"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Member
type Member struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    CustomerID  string    `json:"customer_id" db:"customer_id"`
    MemberNumber string   `json:"member_number" db:"member_number"`
    MemberLevel string    `json:"member_level" db:"member_level"` // BRONZE, SILVER, GOLD, PLATINUM
    Points      int       `json:"points" db:"points"`
    JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
    Status      string    `json:"status" db:"status"` // ACTIVE, SUSPENDED
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Credit Score History
type CreditScoreHistory struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    CustomerID  string    `json:"customer_id" db:"customer_id"`
    Score       int       `json:"score" db:"score"`
    Factors     string    `json:"factors" db:"factors"` // JSON
    CalculatedAt time.Time `json:"calculated_at" db:"calculated_at"`
}
```

### Key Use Cases

```go
// biz/customer.go
type CustomerUsecase struct {
    repo  data.CustomerRepo
    log   *log.Helper
    iam   *IAMClient
    event *event.Publisher
}

func NewCustomerUsecase(repo data.CustomerRepo, logger log.Logger, iam *IAMClient, event *event.Publisher) *CustomerUsecase {
    return &CustomerUsecase{
        repo:  repo,
        log:   log.NewHelper(logger),
        iam:   iam,
        event: event,
    }
}

// CreateCustomer tạo khách hàng mới
func (uc *CustomerUsecase) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*Customer, error) {
    // Generate customer code
    customerCode := uc.generateCustomerCode(ctx, req.TenantID)

    customer := &Customer{
        ID:           uuid.New().String(),
        TenantID:     req.TenantID,
        CustomerCode: customerCode,
        FullName:     req.FullName,
        Email:        req.Email,
        Phone:        req.Phone,
        DateOfBirth:  req.DateOfBirth,
        IDNumber:     req.IDNumber,
        IDType:       req.IDType,
        Address:      req.Address,
        CustomerType: req.CustomerType,
        Status:       "ACTIVE",
        CreditScore:  500, // Default score
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    // Save to database
    if err := uc.repo.Create(ctx, customer); err != nil {
        return nil, err
    }

    // Create user in Zitadel
    if err := uc.iam.CreateUser(ctx, &iam.CreateUserRequest{
        Email:    req.Email,
        Username: req.Email,
        Name:     req.FullName,
    }); err != nil {
        uc.log.Errorf("Failed to create user in Zitadel: %v", err)
    }

    // Publish event
    uc.event.Publish(ctx, &CustomerCreatedEvent{
        CustomerID:   customer.ID,
        CustomerCode: customerCode,
        TenantID:     req.TenantID,
    })

    return customer, nil
}

// GetCustomer lấy thông tin khách hàng
func (uc *CustomerUsecase) GetCustomer(ctx context.Context, tenantID, customerID string) (*Customer, error) {
    customer, err := uc.repo.Get(ctx, tenantID, customerID)
    if err != nil {
        if errors.Is(err, data.ErrNotFound) {
            return nil, errors.New("customer not found")
        }
        return nil, err
    }
    return customer, nil
}

// ListCustomers liệt kê khách hàng
func (uc *CustomerUsecase) ListCustomers(ctx context.Context, req *ListCustomersRequest) (*ListCustomersResponse, error) {
    customers, total, err := uc.repo.List(ctx, req.TenantID, req.Page, req.PageSize)
    if err != nil {
        return nil, err
    }

    return &ListCustomersResponse{
        Customers: customers,
        Total:     total,
        Page:      req.Page,
        PageSize:  req.PageSize,
    }, nil
}

// UpdateCustomer cập nhật thông tin khách hàng
func (uc *CustomerUsecase) UpdateCustomer(ctx context.Context, tenantID, customerID string, req *UpdateCustomerRequest) (*Customer, error) {
    customer, err := uc.repo.Get(ctx, tenantID, customerID)
    if err != nil {
        return nil, err
    }

    if req.FullName != "" {
        customer.FullName = req.FullName
    }
    if req.Phone != "" {
        customer.Phone = req.Phone
    }
    if req.Address != "" {
        customer.Address = req.Address
    }
    customer.UpdatedAt = time.Now()

    if err := uc.repo.Update(ctx, customer); err != nil {
        return nil, err
    }

    return customer, nil
}

// biz/credit_score.go
type CreditScoreUsecase struct {
    repo     data.CreditScoreRepo
    customer *CustomerUsecase
    loan     *LoanClient
    log      *log.Helper
}

// CalculateCreditScore tính điểm tín dụng
func (uc *CreditScoreUsecase) CalculateCreditScore(ctx context.Context, tenantID, customerID string) (*CreditScoreResult, error) {
    customer, err := uc.customer.GetCustomer(ctx, tenantID, customerID)
    if err != nil {
        return nil, err
    }

    factors := make(map[string]int)

    // Factor 1: Payment History (40%)
    paymentScore := uc.calculatePaymentHistory(ctx, customerID)
    factors["payment_history"] = paymentScore

    // Factor 2: Credit Utilization (30%)
    utilizationScore := uc.calculateCreditUtilization(ctx, customerID)
    factors["credit_utilization"] = utilizationScore

    // Factor 3: Length of Credit History (15%)
    historyScore := uc.calculateCreditHistoryLength(ctx, customerID)
    factors["credit_history_length"] = historyScore

    // Factor 4: Types of Credit (10%)
    typesScore := uc.calculateCreditTypes(ctx, customerID)
    factors["credit_types"] = typesScore

    // Factor 5: New Credit (5%)
    newCreditScore := uc.calculateNewCredit(ctx, customerID)
    factors["new_credit"] = newCreditScore

    // Calculate weighted score
    totalScore := paymentScore*40/100 +
        utilizationScore*30/100 +
        historyScore*15/100 +
        typesScore*10/100 +
        newCreditScore*5/100

    // Update customer credit score
    customer.CreditScore = totalScore
    customer.UpdatedAt = time.Now()
    if err := uc.customer.repo.Update(ctx, customer); err != nil {
        return nil, err
    }

    // Save history
    history := &CreditScoreHistory{
        ID:          uuid.New().String(),
        TenantID:    tenantID,
        CustomerID:  customerID,
        Score:       totalScore,
        Factors:     uc.encodeFactors(factors),
        CalculatedAt: time.Now(),
    }
    if err := uc.repo.CreateHistory(ctx, history); err != nil {
        uc.log.Errorf("Failed to save credit score history: %v", err)
    }

    return &CreditScoreResult{
        Score:      totalScore,
        Factors:    factors,
        Rating:     uc.getRating(totalScore),
        CalculatedAt: time.Now(),
    }, nil
}

func (uc *CreditScoreUsecase) calculatePaymentHistory(ctx context.Context, customerID string) int {
    // Get loan contracts and check repayment history
    // Returns score 0-100
    return 85
}

func (uc *CreditScoreUsecase) calculateCreditUtilization(ctx context.Context, customerID string) int {
    // Calculate total credit used / total credit limit
    // Returns score 0-100
    return 70
}

func (uc *CreditScoreUsecase) calculateCreditHistoryLength(ctx context.Context, customerID string) int {
    // Check how long customer has had credit
    // Returns score 0-100
    return 60
}

func (uc *CreditScoreUsecase) calculateCreditTypes(ctx context.Context, customerID string) int {
    // Check variety of credit types
    // Returns score 0-100
    return 80
}

func (uc *CreditScoreUsecase) calculateNewCredit(ctx context.Context, customerID string) int {
    // Check recent credit inquiries
    // Returns score 0-100
    return 90
}

func (uc *CreditScoreUsecase) getRating(score int) string {
    switch {
    case score >= 800:
        return "EXCELLENT"
    case score >= 700:
        return "GOOD"
    case score >= 600:
        return "FAIR"
    default:
        return "POOR"
    }
}
```

### Database Schema

```sql
-- Customers
CREATE TABLE customers (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    customer_code VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    id_number VARCHAR(50),
    id_type VARCHAR(20), -- ID_CARD, PASSPORT, BUSINESS_LICENSE
    address TEXT,
    customer_type VARCHAR(20) NOT NULL, -- INDIVIDUAL, CORPORATE
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    credit_score INT NOT NULL DEFAULT 500,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customers_tenant ON customers(tenant_id);
CREATE INDEX idx_customers_code ON customers(customer_code);
CREATE INDEX idx_customers_email ON customers(email);

-- Members
CREATE TABLE members (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    member_number VARCHAR(50) NOT NULL UNIQUE,
    member_level VARCHAR(20) NOT NULL DEFAULT 'BRONZE', -- BRONZE, SILVER, GOLD, PLATINUM
    points INT NOT NULL DEFAULT 0,
    joined_at DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE INDEX idx_members_tenant ON members(tenant_id);
CREATE INDEX idx_members_customer ON members(customer_id);

-- Credit Score History
CREATE TABLE credit_score_history (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    score INT NOT NULL,
    factors JSONB,
    calculated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);

CREATE INDEX idx_csh_customer ON credit_score_history(customer_id);
CREATE INDEX idx_csh_date ON credit_score_history(calculated_at);
```

### gRPC API

```protobuf
syntax = "proto3";

package crm.v1;

option go_package = "crm-service/api/crm/v1;v1";

service CRMService {
  // Customers
  rpc CreateCustomer(CreateCustomerRequest) returns (Customer);
  rpc GetCustomer(GetCustomerRequest) returns (Customer);
  rpc ListCustomers(ListCustomersRequest) returns (ListCustomersResponse);
  rpc UpdateCustomer(UpdateCustomerRequest) returns (Customer);

  // Members
  rpc CreateMember(CreateMemberRequest) returns (Member);
  rpc GetMember(GetMemberRequest) returns (Member);
  rpc ListMembers(ListMembersRequest) returns (ListMembersResponse);
  rpc UpdateMemberPoints(UpdateMemberPointsRequest) returns (Member);

  // Credit Score
  rpc CalculateCreditScore(CalculateCreditScoreRequest) returns (CreditScoreResult);
  rpc GetCreditScoreHistory(GetCreditScoreHistoryRequest) returns (CreditScoreHistoryResponse);
}

message Customer {
  string id = 1;
  string tenant_id = 2;
  string customer_code = 3;
  string full_name = 4;
  string email = 5;
  string phone = 6;
  string date_of_birth = 7;
  string id_number = 8;
  string id_type = 9;
  string address = 10;
  string customer_type = 11;
  string status = 12;
  int32 credit_score = 13;
  string created_at = 14;
  string updated_at = 15;
}

message CreateCustomerRequest {
  string tenant_id = 1;
  string full_name = 2;
  string email = 3;
  string phone = 4;
  string date_of_birth = 5;
  string id_number = 6;
  string id_type = 7;
  string address = 8;
  string customer_type = 9;
}

message ListCustomersRequest {
  string tenant_id = 1;
  int32 page = 2;
  int32 page_size = 3;
  string status = 4;
  string customer_type = 5;
}

message ListCustomersResponse {
  repeated Customer customers = 1;
  int64 total = 2;
  int32 page = 3;
  int32 page_size = 4;
}

message CreditScoreResult {
  int32 score = 1;
  map<string, int32> factors = 2;
  string rating = 3;
  string calculated_at = 4;
}
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"
```

---

## 2. HRM Service

### Responsibilities
- Quản lý cơ cấu tổ chức (Phòng ban, Chức danh)
- Quản lý thông tin nhân sự
- Quản lý lương thưởng

### Domain Models

```go
// Organization
type Organization struct {
    ID        string    `json:"id" db:"id"`
    TenantID  string    `json:"tenant_id" db:"tenant_id"`
    ParentID  *string   `json:"parent_id" db:"parent_id"`
    Code      string    `json:"code" db:"code"`
    Name      string    `json:"name" db:"name"`
    Level     int       `json:"level" db:"level"` // 1=Company, 2=Department, 3=Team
    Type      string    `json:"type" db:"type"` // HEADQUARTER, BRANCH, DEPARTMENT, TEAM
    ManagerID *string   `json:"manager_id" db:"manager_id"`
    Status    string    `json:"status" db:"status"` // ACTIVE, INACTIVE
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Employee
type Employee struct {
    ID             string    `json:"id" db:"id"`
    TenantID       string    `json:"tenant_id" db:"tenant_id"`
    EmployeeNumber string    `json:"employee_number" db:"employee_number"`
    FullName       string    `json:"full_name" db:"full_name"`
    Email          string    `json:"email" db:"email"`
    Phone          string    `json:"phone" db:"phone"`
    OrganizationID string    `json:"organization_id" db:"organization_id"`
    PositionID     string    `json:"position_id" db:"position_id"`
    ManagerID      *string   `json:"manager_id" db:"manager_id"`
    JoinDate       time.Time `json:"join_date" db:"join_date"`
    Status         string    `json:"status" db:"status"` // ACTIVE, RESIGNED, TERMINATED
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Position
type Position struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    Code        string    `json:"code" db:"code"`
    Name        string    `json:"name" db:"name"`
    Level       int       `json:"level" db:"level"` // 1=CEO, 2=Director, 3=Manager, 4=Staff
    Description string    `json:"description" db:"description"`
    Status      string    `json:"status" db:"status"` // ACTIVE, INACTIVE
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
```

### Database Schema

```sql
-- Organizations
CREATE TABLE organizations (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    parent_id VARCHAR(36),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    level INT NOT NULL,
    type VARCHAR(50) NOT NULL,
    manager_id VARCHAR(36),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES organizations(id)
);

CREATE INDEX idx_org_tenant ON organizations(tenant_id);
CREATE INDEX idx_org_parent ON organizations(parent_id);

-- Employees
CREATE TABLE employees (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    employee_number VARCHAR(50) NOT NULL UNIQUE,
    full_name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    phone VARCHAR(20),
    organization_id VARCHAR(36) NOT NULL,
    position_id VARCHAR(36) NOT NULL,
    manager_id VARCHAR(36),
    join_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE INDEX idx_emp_tenant ON employees(tenant_id);
CREATE INDEX idx_emp_org ON employees(organization_id);

-- Positions
CREATE TABLE positions (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    level INT NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pos_tenant ON positions(tenant_id);
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"
```

---

## 3. Notification Service

### Responsibilities
- Gửi thông báo qua Zalo (ZOA, ZNS)
- Gửi Email
- Quản lý mẫu thông báo
- Queue-based delivery

### Domain Models

```go
// NotificationTemplate
type NotificationTemplate struct {
    ID        string    `json:"id" db:"id"`
    TenantID  string    `json:"tenant_id" db:"tenant_id"`
    Code      string    `json:"code" db:"code"`
    Name      string    `json:"name" db:"name"`
    Channel   string    `json:"channel" db:"channel"` // ZALO, EMAIL, PUSH
    Subject   string    `json:"subject" db:"subject"`
    Body      string    `json:"body" db:"body"`
    Variables string    `json:"variables" db:"variables"` // JSON array
    Status    string    `json:"status" db:"status"` // ACTIVE, INACTIVE
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Notification
type Notification struct {
    ID         string    `json:"id" db:"id"`
    TenantID   string    `json:"tenant_id" db:"tenant_id"`
    Recipient  string    `json:"recipient" db:"recipient"` // Phone, Email, User ID
    Channel    string    `json:"channel" db:"channel"` // ZALO, EMAIL, PUSH
    TemplateID string    `json:"template_id" db:"template_id"`
    Subject    string    `json:"subject" db:"subject"`
    Body       string    `json:"body" db:"body"`
    Status     string    `json:"status" db:"status"` // PENDING, SENT, FAILED
    SentAt     *time.Time `json:"sent_at" db:"sent_at"`
    ErrorMsg   *string   `json:"error_msg" db:"error_msg"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
```

### Zalo Integration

```go
// data/zalo.go
type ZaloClient struct {
    client *http.Client
    config *conf.Zalo
}

type ZaloOARequest struct {
    Message struct {
        Text string `json:"text"`
    } `json:"message"`
    Recipient struct {
        UserId string `json:"user_id"`
    } `json:"recipient"`
}

type ZaloZNSRequest struct {
    TemplateID string                 `json:"template_id"`
    Phone      string                 `json:"phone"`
    Params     map[string]interface{} `json:"params"`
}

func (c *ZaloClient) SendOA(ctx context.Context, userID, message string) error {
    req := &ZaloOARequest{}
    req.Message.Text = message
    req.Recipient.UserId = userID

    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST",
        c.config.OAEndpoint+"/message/send",
        bytes.NewBuffer(body))
    httpReq.Header.Set("access_token", c.config.OAAccessToken)
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(httpReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("Zalo OA API error: %s", resp.Status)
    }

    return nil
}

func (c *ZaloClient) SendZNS(ctx context.Context, phone, templateID string, params map[string]interface{}) error {
    req := &ZaloZNSRequest{
        TemplateID: templateID,
        Phone:      phone,
        Params:     params,
    }

    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST",
        c.config.ZNSEndpoint+"/message/template/send",
        bytes.NewBuffer(body))
    httpReq.Header.Set("access_token", c.config.ZNSAccessToken)
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(httpReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("Zalo ZNS API error: %s", resp.Status)
    }

    return nil
}
```

### Database Schema

```sql
-- Notification Templates
CREATE TABLE notification_templates (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(200) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    subject VARCHAR(500),
    body TEXT NOT NULL,
    variables JSONB,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_nt_tenant ON notification_templates(tenant_id);

-- Notifications
CREATE TABLE notifications (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    recipient VARCHAR(200) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    template_id VARCHAR(36),
    subject VARCHAR(500),
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    sent_at TIMESTAMP,
    error_msg TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES notification_templates(id)
);

CREATE INDEX idx_notif_tenant ON notifications(tenant_id);
CREATE INDEX idx_notif_status ON notifications(status);
CREATE INDEX idx_notif_created ON notifications(created_at);
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"
```

---

## 4. System Config Service

### Responsibilities
- Quản lý tham số hệ thống
- Quản lý ngày làm việc, ngày nghỉ lễ
- Quản lý đa ngôn ngữ

### Domain Models

```go
// SystemParameter
type SystemParameter struct {
    ID        string    `json:"id" db:"id"`
    TenantID  string    `json:"tenant_id" db:"tenant_id"`
    Code      string    `json:"code" db:"code"`
    Name      string    `json:"name" db:"name"`
    Value     string    `json:"value" db:"value"`
    ValueType string    `json:"value_type" db:"value_type"` // STRING, NUMBER, BOOLEAN, JSON
    Category  string    `json:"category" db:"category"`
    IsPublic  bool      `json:"is_public" db:"is_public"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// WorkingDay
type WorkingDay struct {
    Date      time.Time `json:"date" db:"date"`
    TenantID  string    `json:"tenant_id" db:"tenant_id"`
    IsWorking bool      `json:"is_working" db:"is_working"`
    IsHoliday bool      `json:"is_holiday" db:"is_holiday"`
    Note      string    `json:"note" db:"note"`
}

// Translation
type Translation struct {
    ID        string    `json:"id" db:"id"`
    TenantID  string    `json:"tenant_id" db:"tenant_id"`
    Key       string    `json:"key" db:"key"`
    Language  string    `json:"language" db:"language"`
    Value     string    `json:"value" db:"value"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### Database Schema

```sql
-- System Parameters
CREATE TABLE system_parameters (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    value TEXT NOT NULL,
    value_type VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, code)
);

CREATE INDEX idx_sp_tenant ON system_parameters(tenant_id);
CREATE INDEX idx_sp_category ON system_parameters(category);

-- Working Days
CREATE TABLE working_days (
    date DATE NOT NULL,
    tenant_id VARCHAR(36),
    is_working BOOLEAN NOT NULL DEFAULT true,
    is_holiday BOOLEAN NOT NULL DEFAULT false,
    note VARCHAR(500),
    PRIMARY KEY (date, tenant_id)
);

-- Translations
CREATE TABLE translations (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36),
    key VARCHAR(200) NOT NULL,
    language VARCHAR(10) NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, key, language)
);

CREATE INDEX idx_trans_tenant ON translations(tenant_id);
CREATE INDEX idx_trans_key ON translations(key);
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"
```

---

## 5. BPM Engine Service

### Responsibilities
- Wrapper gRPC cho Camunda 7
- Điều phối workflow
- Event publishing

### Domain Models

```go
// WorkflowDefinition
type WorkflowDefinition struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    Key         string    `json:"key" db:"key"`
    Name        string    `json:"name" db:"name"`
    Version     int       `json:"version" db:"version"`
    BPMNXML     string    `json:"bpmn_xml" db:"bpmn_xml"`
    Status      string    `json:"status" db:"status"` // ACTIVE, INACTIVE
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WorkflowInstance
type WorkflowInstance struct {
    ID          string    `json:"id" db:"id"`
    TenantID    string    `json:"tenant_id" db:"tenant_id"`
    DefinitionID string  `json:"definition_id" db:"definition_id"`
    BusinessKey string    `json:"business_key" db:"business_key"`
    Status      string    `json:"status" db:"status"` // RUNNING, SUSPENDED, COMPLETED, TERMINATED
    Variables   string    `json:"variables" db:"variables"` // JSON
    StartedAt   time.Time `json:"started_at" db:"started_at"`
    CompletedAt *time.Time `json:"completed_at" db:"completed_at"`
}

// UserTask
type UserTask struct {
    ID           string    `json:"id" db:"id"`
    TenantID     string    `json:"tenant_id" db:"tenant_id"`
    InstanceID   string    `json:"instance_id" db:"instance_id"`
    TaskKey      string    `json:"task_key" db:"task_key"`
    TaskName     string    `json:"task_name" db:"task_name"`
    Assignee     *string   `json:"assignee" db:"assignee"`
    CandidateGroups string  `json:"candidate_groups" db:"candidate_groups"`
    Status       string    `json:"status" db:"status"` // CREATED, ASSIGNED, COMPLETED, CANCELLED
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
    CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`
}
```

### Camunda Integration

```go
// data/camunda.go
type CamundaClient struct {
    client *http.Client
    baseURL string
}

func (c *CamundaClient) StartProcess(ctx context.Context, key string, variables map[string]interface{}) (string, error) {
    req := map[string]interface{}{
        "processDefinitionKey": key,
        "variables":            variables,
    }

    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST",
        c.baseURL+"/process-instance",
        bytes.NewBuffer(body))
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 201 {
        return "", fmt.Errorf("Camunda API error: %s", resp.Status)
    }

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)

    return result["id"].(string), nil
}

func (c *CamundaClient) CompleteTask(ctx context.Context, taskID string, variables map[string]interface{}) error {
    req := map[string]interface{}{
        "variables": variables,
    }

    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequest("POST",
        c.baseURL+"/task/"+taskID+"/complete",
        bytes.NewBuffer(body))
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(httpReq)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 204 {
        return fmt.Errorf("Camunda API error: %s", resp.Status)
    }

    return nil
}
```

### Database Schema

```sql
-- Workflow Definitions
CREATE TABLE workflow_definitions (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    key VARCHAR(100) NOT NULL,
    name VARCHAR(200) NOT NULL,
    version INT NOT NULL,
    bpmn_xml TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, key, version)
);

CREATE INDEX idx_wd_tenant ON workflow_definitions(tenant_id);

-- Workflow Instances
CREATE TABLE workflow_instances (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    definition_id VARCHAR(36) NOT NULL,
    business_key VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
    variables JSONB,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    FOREIGN KEY (definition_id) REFERENCES workflow_definitions(id)
);

CREATE INDEX idx_wi_tenant ON workflow_instances(tenant_id);
CREATE INDEX idx_wi_status ON workflow_instances(status);

-- User Tasks
CREATE TABLE user_tasks (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    instance_id VARCHAR(36) NOT NULL,
    task_key VARCHAR(100) NOT NULL,
    task_name VARCHAR(200) NOT NULL,
    assignee VARCHAR(36),
    candidate_groups VARCHAR(500),
    status VARCHAR(20) NOT NULL DEFAULT 'CREATED',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    FOREIGN KEY (instance_id) REFERENCES workflow_instances(id)
);

CREATE INDEX idx_ut_tenant ON user_tasks(tenant_id);
CREATE INDEX idx_ut_instance ON user_tasks(instance_id);
CREATE INDEX idx_ut_assignee ON user_tasks(assignee);
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"
```

---

## 🔧 Configuration

### config.yaml

```yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 10s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 10s

data:
  database:
    driver: postgres
    source: ${DATABASE_URL}
  redis:
    addr: ${REDIS_ADDR}
    read_timeout: 0.2s
    write_timeout: 0.2s
    db: 0

auth:
  jwt:
    jwks_endpoint: https://auth.arda.io.vn/oauth/v2/keys
    issuer: https://auth.arda.io.vn
    audience: crm-service

tenant:
  header: X-Tenant-ID

event:
  enabled: true
  brokers: ${REDPANDA_BROKERS}
  topic_prefix: arda.crm
```

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o service ./cmd/crm-service

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/service /app/service

EXPOSE 8000 9000

CMD ["/app/service", "-conf", "/data/conf"]
```

---

## 📊 Build & Deploy

### Build

```bash
# Build specific service
cd crm-service
go build -o bin/service ./cmd/crm-service

# Build all services
make build
```

### Docker Build

```bash
# Build image
docker build -t ghcr.io/arda-labs/crm-service:latest .

# Push to registry
docker push ghcr.io/arda-labs/crm-service:latest
```

### Deploy to K3s

```bash
# Apply manifests
kubectl apply -f arda-infra/apps/crm/base/

# Or via ArgoCD
kubectl apply -f arda-infra/argocd/apps/crm.yaml
```

---

*Last Updated: 2026-04-24*
