# Core Banking Services — Dịch vụ Ngân hàng Cốt lõi

> Microservices Java với GraalVM Native Image cho tối ưu RAM
> Phân rã theo Domain-Driven Design (DDD)

---

## 📋 Overview

Nhóm Core Banking Services xử lý các giao dịch tài chính yêu cầu tính nhất quán (Consistency) và ACID tuyệt đối. Tất cả services được biên dịch sang **GraalVM Native Image** để giảm 70-80% lượng RAM tiêu thụ.

### Services
1. **Accounting Service** — Hạch toán kế toán, bút toán kép, sổ cái
2. **Loan Service** — Quản lý khoản vay, giải ngân, thu hồi
3. **Deposit Service** — Tiền gửi TCTD, nguồn vốn
4. **Treasury Service** — Quản lý thanh khoản, ngoại hối

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Core Banking Layer                        │
│                    (Java + GraalVM)                          │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  Accounting  │  │     Loan     │  │   Deposit    │     │
│  │   Service    │  │   Service    │  │   Service    │     │
│  │              │  │              │  │              │     │
│  │ • COA Mgmt   │  │ • Products   │  │ • Deposits   │     │
│  │ • Double-    │  │ • Disburse-  │  │ • Capital    │     │
│  │   Entry      │  │   ment       │  │   Sources    │     │
│  │ • Ledger     │  │ • Repayment  │  │ • Interest   │     │
│  │ • Closing    │  │ • Collateral │  │              │     │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘     │
│         │                  │                  │             │
│         └──────────────────┼──────────────────┘             │
│                            ▼                                 │
│                    ┌──────────────┐                         │
│                    │   Outbox     │                         │
│                    │   Pattern    │                         │
│                    └──────┬───────┘                         │
└───────────────────────────┼─────────────────────────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │   Redpanda   │
                    └──────────────┘
```

---

## 📁 Monorepo Structure

```
arda-core/
├── services/
│   ├── accounting/
│   │   ├── src/main/java/com/arda/accounting/
│   │   │   ├── AccountingApplication.java
│   │   │   ├── domain/
│   │   │   │   ├── ChartOfAccount.java
│   │   │   │   ├── JournalEntry.java
│   │   │   │   ├── Ledger.java
│   │   │   │   └── AccountBalance.java
│   │   │   ├── application/
│   │   │   │   ├── AccountService.java
│   │   │   │   ├── JournalService.java
│   │   │   │   └── ClosingService.java
│   │   │   ├── infrastructure/
│   │   │   │   ├── repository/
│   │   │   │   ├── outbox/
│   │   │   │   └── grpc/
│   │   │   └── config/
│   │   ├── src/main/resources/
│   │   │   ├── application.yml
│   │   │   ├── schema.sql
│   │   │   └── native-image/
│   │   │       ├── reflect-config.json
│   │   │       └── resource-config.json
│   │   ├── build.gradle.kts
│   │   └── Dockerfile
│   │
│   ├── loan/
│   │   ├── src/main/java/com/arda/loan/
│   │   │   ├── LoanApplication.java
│   │   │   ├── domain/
│   │   │   │   ├── LoanProduct.java
│   │   │   │   ├── LoanContract.java
│   │   │   │   ├── Disbursement.java
│   │   │   │   ├── Repayment.java
│   │   │   │   ├── Provision.java
│   │   │   │   └── Collateral.java
│   │   │   ├── application/
│   │   │   │   ├── LoanService.java
│   │   │   │   ├── DisbursementService.java
│   │   │   │   └── ProvisionService.java
│   │   │   ├── infrastructure/
│   │   │   └── config/
│   │   ├── src/main/resources/
│   │   │   ├── application.yml
│   │   │   ├── schema.sql
│   │   │   └── workflows/
│   │   │       └── loan-disbursement.bpmn
│   │   ├── build.gradle.kts
│   │   └── Dockerfile
│   │
│   └── deposit/
│       └── ... (similar structure)
│
├── libs/
│   ├── shared-core/
│   │   ├── src/main/java/com/arda/core/
│   │   │   ├── tenant/
│   │   │   │   ├── TenantContext.java
│   │   │   │   └── TenantAware.java
│   │   │   ├── event/
│   │   │   │   ├── DomainEvent.java
│   │   │   │   └── EventPublisher.java
│   │   │   └── exception/
│   │   │       └── BusinessException.java
│   │   └── build.gradle.kts
│   │
│   ├── grpc-client/
│   │   ├── src/main/java/com/arda/grpc/
│   │   │   ├── AccountingClient.java
│   │   │   ├── LoanClient.java
│   │   │   └── DepositClient.java
│   │   └── build.gradle.kts
│   │
│   └── accounting-core/
│       ├── src/main/java/com/arda/accounting/core/
│       │   ├── DoubleEntryValidator.java
│       │   ├── LedgerService.java
│       │   └── ClosingCalculator.java
│       └── build.gradle.kts
│
├── build.gradle.kts
└── settings.gradle.kts
```

---

## 1. Accounting Service (Kế toán)

### Responsibilities
- Quản lý hệ thống tài khoản quy định (COA - Chart of Accounts)
- Hạch toán bút toán kép (Double-entry bookkeeping)
- Quản lý sổ cái (General Ledger)
- Kết chuyển thu nhập chi phí (TNCP) cuối ngày
- Tổng hợp số liệu báo cáo

### Domain Model

```java
// Chart of Account
@Entity
@Table(name = "chart_of_accounts")
public class ChartOfAccount {
    @Id
    private String id;

    @Column(nullable = false, length = 20)
    private String accountCode;

    @Column(nullable = false)
    private String accountName;

    @Enumerated(EnumType.STRING)
    private AccountType accountType; // ASSET, LIABILITY, EQUITY, INCOME, EXPENSE

    @Enumerated(EnumType.STRING)
    private AccountCategory category; // CURRENT, NON_CURRENT

    @Column(nullable = false)
    private String tenantId;

    @Column(nullable = false)
    private BigDecimal debitBalance = BigDecimal.ZERO;

    @Column(nullable = false)
    private BigDecimal creditBalance = BigDecimal.ZERO;

    private Boolean isActive = true;
}

// Journal Entry
@Entity
@Table(name = "journal_entries")
public class JournalEntry {
    @Id
    private String id;

    @Column(nullable = false, unique = true, length = 50)
    private String voucherNumber;

    @Column(nullable = false)
    private LocalDate entryDate;

    @Column(nullable = false, length = 100)
    private String description;

    @Enumerated(EnumType.STRING)
    private EntryStatus status; // DRAFT, POSTED, REVERSED

    @Column(nullable = false)
    private String tenantId;

    @OneToMany(mappedBy = "journalEntry", cascade = CascadeType.ALL)
    private List<JournalEntryLine> lines = new ArrayList<>();

    private String createdBy;
    private LocalDateTime createdAt;
    private String postedBy;
    private LocalDateTime postedAt;
}

// Journal Entry Line
@Entity
@Table(name = "journal_entry_lines")
public class JournalEntryLine {
    @Id
    private String id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "journal_entry_id")
    private JournalEntry journalEntry;

    @Column(nullable = false)
    private String accountCode;

    @Enumerated(EnumType.STRING)
    private DebitCredit drcr; // DEBIT, CREDIT

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal amount;

    @Column(length = 500)
    private String description;
}

// General Ledger
@Entity
@Table(name = "general_ledger")
public class GeneralLedger {
    @Id
    private String id;

    @Column(nullable = false)
    private String accountCode;

    @Column(nullable = false)
    private LocalDate transactionDate;

    @Column(nullable = false)
    private String voucherNumber;

    @Column(nullable = false)
    private String voucherType;

    @Enumerated(EnumType.STRING)
    private DebitCredit drcr;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal amount;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal balance;

    @Column(nullable = false)
    private String tenantId;

    private LocalDateTime createdAt;
}
```

### Key Services

```java
@Service
@Transactional
public class AccountService {

    @Autowired
    private ChartOfAccountRepository coaRepository;

    @Autowired
    private JournalEntryRepository journalEntryRepository;

    @Autowired
    private GeneralLedgerRepository ledgerRepository;

    @Autowired
    private EventPublisher eventPublisher;

    /**
     * Tạo bút toán kép
     */
    public String createJournalEntry(CreateJournalEntryRequest request) {
        // Validate double-entry: total debits = total credits
        validateDoubleEntry(request.getLines());

        JournalEntry entry = new JournalEntry();
        entry.setVoucherNumber(generateVoucherNumber());
        entry.setEntryDate(request.getEntryDate());
        entry.setDescription(request.getDescription());
        entry.setTenantId(TenantContext.getCurrentTenantId());

        List<JournalEntryLine> lines = request.getLines().stream()
            .map(line -> {
                JournalEntryLine jel = new JournalEntryLine();
                jel.setAccountCode(line.getAccountCode());
                jel.setDrcr(line.getDrcr());
                jel.setAmount(line.getAmount());
                jel.setDescription(line.getDescription());
                jel.setJournalEntry(entry);
                return jel;
            })
            .collect(Collectors.toList());

        entry.setLines(lines);
        entry.setStatus(EntryStatus.DRAFT);

        journalEntryRepository.save(entry);

        // Publish event
        eventPublisher.publish(new JournalEntryCreatedEvent(
            entry.getId(),
            entry.getVoucherNumber(),
            TenantContext.getCurrentTenantId()
        ));

        return entry.getId();
    }

    /**
     * Post journal entry (ghi sổ)
     */
    public void postJournalEntry(String entryId) {
        JournalEntry entry = journalEntryRepository.findById(entryId)
            .orElseThrow(() -> new BusinessException("Journal entry not found"));

        if (entry.getStatus() != EntryStatus.DRAFT) {
            throw new BusinessException("Only DRAFT entries can be posted");
        }

        // Update ledger
        for (JournalEntryLine line : entry.getLines()) {
            GeneralLedger ledger = new GeneralLedger();
            ledger.setAccountCode(line.getAccountCode());
            ledger.setTransactionDate(entry.getEntryDate());
            ledger.setVoucherNumber(entry.getVoucherNumber());
            ledger.setDrcr(line.getDrcr());
            ledger.setAmount(line.getAmount());
            ledger.setTenantId(entry.getTenantId());

            // Calculate running balance
            BigDecimal currentBalance = getCurrentBalance(
                line.getAccountCode(),
                entry.getTenantId(),
                entry.getEntryDate()
            );
            ledger.setBalance(calculateNewBalance(currentBalance, line));

            ledgerRepository.save(ledger);

            // Update COA balance
            updateCOABalance(line.getAccountCode(), line.getDrcr(), line.getAmount());
        }

        entry.setStatus(EntryStatus.POSTED);
        entry.setPostedAt(LocalDateTime.now());

        // Publish event
        eventPublisher.publish(new JournalEntryPostedEvent(
            entry.getId(),
            entry.getVoucherNumber(),
            TenantContext.getCurrentTenantId()
        ));
    }

    /**
     * Day-end closing (kết chuyển ngày)
     */
    public void performDayEndClosing(LocalDate closingDate) {
        // Get all income and expense accounts
        List<ChartOfAccount> incomeAccounts = coaRepository
            .findByTenantIdAndAccountTypeAndIsActive(
                TenantContext.getCurrentTenantId(),
                AccountType.INCOME,
                true
            );
        List<ChartOfAccount> expenseAccounts = coaRepository
            .findByTenantIdAndAccountTypeAndIsActive(
                TenantContext.getCurrentTenantId(),
                AccountType.EXPENSE,
                true
            );

        // Calculate totals
        BigDecimal totalIncome = calculateTotalBalance(incomeAccounts, closingDate);
        BigDecimal totalExpense = calculateTotalBalance(expenseAccounts, closingDate);

        // Create closing entry
        JournalEntry closingEntry = new JournalEntry();
        closingEntry.setVoucherNumber(generateVoucherNumber());
        closingEntry.setEntryDate(closingDate);
        closingEntry.setDescription("Day-end closing");
        closingEntry.setTenantId(TenantContext.getCurrentTenantId());

        // Add lines
        closingEntry.getLines().addAll(createClosingLines(
            incomeAccounts,
            expenseAccounts,
            closingDate
        ));

        journalEntryRepository.save(closingEntry);

        // Post closing entry
        postJournalEntry(closingEntry.getId());

        // Publish event
        eventPublisher.publish(new DayEndClosingCompletedEvent(
            closingDate,
            totalIncome,
            totalExpense,
            TenantContext.getCurrentTenantId()
        ));
    }
}
```

### Database Schema

```sql
-- Chart of Accounts
CREATE TABLE chart_of_accounts (
    id VARCHAR(36) PRIMARY KEY,
    account_code VARCHAR(20) NOT NULL UNIQUE,
    account_name VARCHAR(200) NOT NULL,
    account_type VARCHAR(20) NOT NULL, -- ASSET, LIABILITY, EQUITY, INCOME, EXPENSE
    account_category VARCHAR(20) NOT NULL, -- CURRENT, NON_CURRENT
    tenant_id VARCHAR(36) NOT NULL,
    debit_balance DECIMAL(20,4) NOT NULL DEFAULT 0,
    credit_balance DECIMAL(20,4) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_coa_tenant ON chart_of_accounts(tenant_id);

-- Journal Entries
CREATE TABLE journal_entries (
    id VARCHAR(36) PRIMARY KEY,
    voucher_number VARCHAR(50) NOT NULL UNIQUE,
    entry_date DATE NOT NULL,
    description VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL, -- DRAFT, POSTED, REVERSED
    tenant_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    posted_by VARCHAR(100),
    posted_at TIMESTAMP
);

CREATE INDEX idx_je_tenant ON journal_entries(tenant_id);
CREATE INDEX idx_je_date ON journal_entries(entry_date);
CREATE INDEX idx_je_status ON journal_entries(status);

-- Journal Entry Lines
CREATE TABLE journal_entry_lines (
    id VARCHAR(36) PRIMARY KEY,
    journal_entry_id VARCHAR(36) NOT NULL,
    account_code VARCHAR(20) NOT NULL,
    drcr VARCHAR(10) NOT NULL, -- DEBIT, CREDIT
    amount DECIMAL(20,4) NOT NULL,
    description VARCHAR(500),
    FOREIGN KEY (journal_entry_id) REFERENCES journal_entries(id)
);

CREATE INDEX idx_jel_entry ON journal_entry_lines(journal_entry_id);

-- General Ledger
CREATE TABLE general_ledger (
    id VARCHAR(36) PRIMARY KEY,
    account_code VARCHAR(20) NOT NULL,
    transaction_date DATE NOT NULL,
    voucher_number VARCHAR(50) NOT NULL,
    voucher_type VARCHAR(50) NOT NULL,
    drcr VARCHAR(10) NOT NULL,
    amount DECIMAL(20,4) NOT NULL,
    balance DECIMAL(20,4) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gl_account ON general_ledger(account_code);
CREATE INDEX idx_gl_date ON general_ledger(transaction_date);
CREATE INDEX idx_gl_tenant ON general_ledger(tenant_id);

-- Outbox Table for Event Publishing
CREATE TABLE accounting_outbox (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    event_data JSONB NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PUBLISHED, FAILED
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

CREATE INDEX idx_outbox_status ON accounting_outbox(status);
CREATE INDEX idx_outbox_created ON accounting_outbox(created_at);
```

### Resource Requirements

```yaml
# K8s Deployment
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### gRPC API

```protobuf
syntax = "proto3";

package accounting.v1;

service AccountingService {
  // COA Management
  rpc CreateCOA(CreateCOARequest) returns (COA);
  rpc GetCOA(GetCOARequest) returns (COA);
  rpc ListCOAs(ListCOAsRequest) returns (ListCOAsResponse);
  rpc UpdateCOA(UpdateCOARequest) returns (COA);
  rpc DeleteCOA(DeleteCOARequest) returns (DeleteCOAResponse);

  // Journal Entries
  rpc CreateJournalEntry(CreateJournalEntryRequest) returns (JournalEntry);
  rpc GetJournalEntry(GetJournalEntryRequest) returns (JournalEntry);
  rpc ListJournalEntries(ListJournalEntriesRequest) returns (ListJournalEntriesResponse);
  rpc PostJournalEntry(PostJournalEntryRequest) returns (JournalEntry);

  // Ledger
  rpc GetAccountBalance(GetAccountBalanceRequest) returns (AccountBalance);
  rpc GetLedger(GetLedgerRequest) returns (GetLedgerResponse);

  // Closing
  rpc PerformDayEndClosing(PerformDayEndClosingRequest) returns (PerformDayEndClosingResponse);
}

message COA {
  string id = 1;
  string account_code = 2;
  string account_name = 3;
  string account_type = 4;
  string account_category = 5;
  string debit_balance = 6;
  string credit_balance = 7;
  bool is_active = 8;
}

message CreateCOARequest {
  string account_code = 1;
  string account_name = 2;
  string account_type = 3;
  string account_category = 4;
}

message JournalEntry {
  string id = 1;
  string voucher_number = 2;
  string entry_date = 3;
  string description = 4;
  string status = 5;
  repeated JournalEntryLine lines = 6;
}

message JournalEntryLine {
  string account_code = 1;
  string drcr = 2;
  string amount = 3;
  string description = 4;
}

message CreateJournalEntryRequest {
  string entry_date = 1;
  string description = 2;
  repeated JournalEntryLine lines = 3;
}

message AccountBalance {
  string account_code = 1;
  string account_name = 2;
  string debit_balance = 3;
  string credit_balance = 4;
  string net_balance = 5;
}
```

---

## 2. Loan Service (Cho vay)

### Responsibilities
- Thiết lập sản phẩm cho vay, lãi suất tham chiếu
- Khởi tạo giải ngân, thu hồi nợ gốc/lãi
- Gia hạn nợ, chuyển nhóm nợ, tính dự phòng rủi ro
- Quản lý tài sản bảo đảm

### Domain Model

```java
// Loan Product
@Entity
@Table(name = "loan_products")
public class LoanProduct {
    @Id
    private String id;

    @Column(nullable = false, length = 50)
    private String productCode;

    @Column(nullable = false)
    private String productName;

    @Enumerated(EnumType.STRING)
    private LoanType loanType; // CONSUMER, BUSINESS, MORTGAGE, OVERDRAFT

    @Column(nullable = false, precision = 10, scale = 4)
    private BigDecimal minInterestRate;

    @Column(nullable = false, precision = 10, scale = 4)
    private BigDecimal maxInterestRate;

    @Column(nullable = false)
    private Integer minTermMonths;

    @Column(nullable = false)
    private Integer maxTermMonths;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal maxAmount;

    private Boolean isActive = true;

    @Column(nullable = false)
    private String tenantId;
}

// Loan Contract
@Entity
@Table(name = "loan_contracts")
public class LoanContract {
    @Id
    private String id;

    @Column(nullable = false, length = 50, unique = true)
    private String contractNumber;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "customer_id")
    private Customer customer;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "product_id")
    private LoanProduct product;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal principalAmount;

    @Column(nullable = false, precision = 10, scale = 4)
    private BigDecimal interestRate;

    @Column(nullable = false)
    private Integer termMonths;

    @Enumerated(EnumType.STRING)
    private RepaymentMethod repaymentMethod; // EQUAL_INSTALLMENTS, EQUAL_PRINCIPAL, BULLET

    @Column(nullable = false)
    private LocalDate disbursementDate;

    @Column(nullable = false)
    private LocalDate maturityDate;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal outstandingPrincipal = BigDecimal.ZERO;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal outstandingInterest = BigDecimal.ZERO;

    @Enumerated(EnumType.STRING)
    private LoanStatus status; // PENDING, ACTIVE, CLOSED, DEFAULTED

    @Enumerated(EnumType.STRING)
    private RiskGroup riskGroup; // GROUP_1, GROUP_2, GROUP_3, GROUP_4, GROUP_5

    @Column(nullable = false)
    private String tenantId;

    private String createdBy;
    private LocalDateTime createdAt;

    @OneToMany(mappedBy = "loanContract", cascade = CascadeType.ALL)
    private List<Disbursement> disbursements = new ArrayList<>();

    @OneToMany(mappedBy = "loanContract", cascade = CascadeType.ALL)
    private List<Repayment> repayments = new ArrayList<>();

    @OneToMany(mappedBy = "loanContract", cascade = CascadeType.ALL)
    private List<Collateral> collaterals = new ArrayList<>();
}

// Disbursement
@Entity
@Table(name = "disbursements")
public class Disbursement {
    @Id
    private String id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "loan_contract_id")
    private LoanContract loanContract;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal amount;

    @Column(nullable = false)
    private LocalDate disbursementDate;

    @Column(length = 50)
    private String accountNumber;

    @Column(length = 50)
    private String bankCode;

    @Enumerated(EnumType.STRING)
    private DisbursementStatus status; // PENDING, COMPLETED, FAILED

    private String referenceNumber;
    private String failureReason;

    private LocalDateTime createdAt;
}

// Repayment Schedule
@Entity
@Table(name = "repayment_schedules")
public class RepaymentSchedule {
    @Id
    private String id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "loan_contract_id")
    private LoanContract loanContract;

    @Column(nullable = false)
    private Integer installmentNumber;

    @Column(nullable = false)
    private LocalDate dueDate;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal principalAmount;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal interestAmount;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal totalAmount;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal principalPaid = BigDecimal.ZERO;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal interestPaid = BigDecimal.ZERO;

    @Enumerated(EnumType.STRING)
    private ScheduleStatus status; // PENDING, PARTIAL, PAID, OVERDUE

    private LocalDateTime paidAt;
}

// Provision (Dự phòng rủi ro)
@Entity
@Table(name = "provisions")
public class Provision {
    @Id
    private String id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "loan_contract_id")
    private LoanContract loanContract;

    @Column(nullable = false, precision = 10, scale = 4)
    private BigDecimal provisionRate;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal provisionAmount;

    @Column(nullable = false)
    private LocalDate calculationDate;

    @Enumerated(EnumType.STRING)
    private RiskGroup riskGroup;

    @Column(nullable = false)
    private String tenantId;
}
```

### Key Services

```java
@Service
@Transactional
public class LoanService {

    @Autowired
    private LoanContractRepository loanContractRepository;

    @Autowired
    private DisbursementRepository disbursementRepository;

    @Autowired
    private RepaymentScheduleRepository scheduleRepository;

    @Autowired
    private AccountingClient accountingClient;

    @Autowired
    private EventPublisher eventPublisher;

    /**
     * Tạo hợp đồng vay
     */
    public String createLoanContract(CreateLoanContractRequest request) {
        // Validate product
        LoanProduct product = productRepository.findById(request.getProductId())
            .orElseThrow(() -> new BusinessException("Product not found"));

        // Validate limits
        if (request.getPrincipalAmount().compareTo(product.getMaxAmount()) > 0) {
            throw new BusinessException("Amount exceeds product limit");
        }

        // Create contract
        LoanContract contract = new LoanContract();
        contract.setContractNumber(generateContractNumber());
        contract.setCustomer(getCustomer(request.getCustomerId()));
        contract.setProduct(product);
        contract.setPrincipalAmount(request.getPrincipalAmount());
        contract.setInterestRate(request.getInterestRate());
        contract.setTermMonths(request.getTermMonths());
        contract.setRepaymentMethod(request.getRepaymentMethod());
        contract.setDisbursementDate(request.getDisbursementDate());
        contract.setMaturityDate(calculateMaturityDate(request.getDisbursementDate(), request.getTermMonths()));
        contract.setStatus(LoanStatus.PENDING);
        contract.setTenantId(TenantContext.getCurrentTenantId());

        // Generate repayment schedule
        List<RepaymentSchedule> schedules = generateRepaymentSchedule(contract);
        contract.setSchedules(schedules);

        loanContractRepository.save(contract);

        // Publish event
        eventPublisher.publish(new LoanContractCreatedEvent(
            contract.getId(),
            contract.getContractNumber(),
            contract.getPrincipalAmount(),
            TenantContext.getCurrentTenantId()
        ));

        return contract.getId();
    }

    /**
     * Giải ngân khoản vay
     */
    @Transactional
    public String processDisbursement(String contractId, DisbursementRequest request) {
        LoanContract contract = loanContractRepository.findById(contractId)
            .orElseThrow(() -> new BusinessException("Contract not found"));

        if (contract.getStatus() != LoanStatus.PENDING) {
            throw new BusinessException("Contract must be in PENDING status");
        }

        // Check total disbursements
        BigDecimal totalDisbursed = contract.getDisbursements().stream()
            .filter(d -> d.getStatus() == DisbursementStatus.COMPLETED)
            .map(Disbursement::getAmount)
            .reduce(BigDecimal.ZERO, BigDecimal::add);

        if (totalDisbursed.add(request.getAmount()).compareTo(contract.getPrincipalAmount()) > 0) {
            throw new BusinessException("Disbursement amount exceeds contract principal");
        }

        // Create disbursement
        Disbursement disbursement = new Disbursement();
        disbursement.setLoanContract(contract);
        disbursement.setAmount(request.getAmount());
        disbursement.setDisbursementDate(LocalDate.now());
        disbursement.setAccountNumber(request.getAccountNumber());
        disbursement.setBankCode(request.getBankCode());
        disbursement.setStatus(DisbursementStatus.PENDING);

        disbursementRepository.save(disbursement);

        // TODO: Call external banking system to transfer money

        // Update contract status if this is the first disbursement
        if (contract.getStatus() == LoanStatus.PENDING) {
            contract.setStatus(LoanStatus.ACTIVE);
            contract.setOutstandingPrincipal(request.getAmount());
        } else {
            contract.setOutstandingPrincipal(
                contract.getOutstandingPrincipal().add(request.getAmount())
            );
        }

        disbursement.setStatus(DisbursementStatus.COMPLETED);
        disbursement.setReferenceNumber(generateReferenceNumber());

        // Create accounting entry via Outbox pattern
        createDisbursementAccountingEntry(disbursement);

        // Publish event
        eventPublisher.publish(new DisbursementCompletedEvent(
            disbursement.getId(),
            contract.getId(),
            contract.getContractNumber(),
            disbursement.getAmount(),
            TenantContext.getCurrentTenantId()
        ));

        return disbursement.getId();
    }

    /**
     * Thu hồi nợ
     */
    @Transactional
    public void processRepayment(String contractId, RepaymentRequest request) {
        LoanContract contract = loanContractRepository.findById(contractId)
            .orElseThrow(() -> new BusinessException("Contract not found"));

        if (contract.getStatus() != LoanStatus.ACTIVE) {
            throw new BusinessException("Contract must be in ACTIVE status");
        }

        // Get pending schedules
        List<RepaymentSchedule> pendingSchedules = scheduleRepository
            .findByLoanContractIdAndStatusOrderByDueDateAsc(
                contractId,
                ScheduleStatus.PENDING
            );

        if (pendingSchedules.isEmpty()) {
            throw new BusinessException("No pending schedules found");
        }

        BigDecimal remainingAmount = request.getAmount();

        // Apply to schedules (FIFO)
        for (RepaymentSchedule schedule : pendingSchedules) {
            if (remainingAmount.compareTo(BigDecimal.ZERO) <= 0) {
                break;
            }

            BigDecimal principalDue = schedule.getPrincipalAmount()
                .subtract(schedule.getPrincipalPaid());
            BigDecimal interestDue = schedule.getInterestAmount()
                .subtract(schedule.getInterestPaid());
            BigDecimal totalDue = principalDue.add(interestDue);

            BigDecimal payment = remainingAmount.min(totalDue);

            // Calculate allocation (interest first)
            BigDecimal interestPayment = payment.min(interestDue);
            BigDecimal principalPayment = payment.subtract(interestPayment);

            // Update schedule
            schedule.setInterestPaid(schedule.getInterestPaid().add(interestPayment));
            schedule.setPrincipalPaid(schedule.getPrincipalPaid().add(principalPayment));

            if (schedule.getPrincipalPaid().compareTo(schedule.getPrincipalAmount()) >= 0
                && schedule.getInterestPaid().compareTo(schedule.getInterestAmount()) >= 0) {
                schedule.setStatus(ScheduleStatus.PAID);
                schedule.setPaidAt(LocalDateTime.now());
            } else {
                schedule.setStatus(ScheduleStatus.PARTIAL);
            }

            remainingAmount = remainingAmount.subtract(payment);
        }

        // Update contract outstanding
        BigDecimal principalPaid = request.getAmount().subtract(remainingAmount);
        contract.setOutstandingPrincipal(
            contract.getOutstandingPrincipal().subtract(principalPaid)
        );

        // Create accounting entry
        createRepaymentAccountingEntry(contract, request.getAmount());

        // Check if fully paid
        if (contract.getOutstandingPrincipal().compareTo(BigDecimal.ZERO) <= 0) {
            contract.setStatus(LoanStatus.CLOSED);
        }

        // Publish event
        eventPublisher.publish(new RepaymentProcessedEvent(
            contract.getId(),
            contract.getContractNumber(),
            request.getAmount(),
            contract.getOutstandingPrincipal(),
            TenantContext.getCurrentTenantId()
        ));
    }

    /**
     * Tính dự phòng rủi ro
     */
    public void calculateProvisions(LocalDate asOfDate) {
        List<LoanContract> activeContracts = loanContractRepository
            .findByStatusAndTenantId(
                LoanStatus.ACTIVE,
                TenantContext.getCurrentTenantId()
            );

        for (LoanContract contract : activeContracts) {
            RiskGroup riskGroup = determineRiskGroup(contract);
            BigDecimal provisionRate = getProvisionRate(riskGroup);
            BigDecimal provisionAmount = contract.getOutstandingPrincipal()
                .multiply(provisionRate)
                .divide(new BigDecimal("100"));

            // Create or update provision
            Provision provision = provisionRepository
                .findByLoanContractIdAndCalculationDate(
                    contract.getId(),
                    asOfDate
                )
                .orElse(new Provision());

            provision.setLoanContract(contract);
            provision.setProvisionRate(provisionRate);
            provision.setProvisionAmount(provisionAmount);
            provision.setCalculationDate(asOfDate);
            provision.setRiskGroup(riskGroup);
            provision.setTenantId(TenantContext.getCurrentTenantId());

            provisionRepository.save(provision);
        }
    }
}
```

### Database Schema

```sql
-- Loan Products
CREATE TABLE loan_products (
    id VARCHAR(36) PRIMARY KEY,
    product_code VARCHAR(50) NOT NULL UNIQUE,
    product_name VARCHAR(200) NOT NULL,
    loan_type VARCHAR(50) NOT NULL,
    min_interest_rate DECIMAL(10,4) NOT NULL,
    max_interest_rate DECIMAL(10,4) NOT NULL,
    min_term_months INT NOT NULL,
    max_term_months INT NOT NULL,
    max_amount DECIMAL(20,4) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    tenant_id VARCHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_lp_tenant ON loan_products(tenant_id);

-- Loan Contracts
CREATE TABLE loan_contracts (
    id VARCHAR(36) PRIMARY KEY,
    contract_number VARCHAR(50) NOT NULL UNIQUE,
    customer_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    principal_amount DECIMAL(20,4) NOT NULL,
    interest_rate DECIMAL(10,4) NOT NULL,
    term_months INT NOT NULL,
    repayment_method VARCHAR(50) NOT NULL,
    disbursement_date DATE NOT NULL,
    maturity_date DATE NOT NULL,
    outstanding_principal DECIMAL(20,4) NOT NULL DEFAULT 0,
    outstanding_interest DECIMAL(20,4) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    risk_group VARCHAR(20),
    tenant_id VARCHAR(36) NOT NULL,
    created_by VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES loan_products(id)
);

CREATE INDEX idx_lc_tenant ON loan_contracts(tenant_id);
CREATE INDEX idx_lc_customer ON loan_contracts(customer_id);
CREATE INDEX idx_lc_status ON loan_contracts(status);

-- Disbursements
CREATE TABLE disbursements (
    id VARCHAR(36) PRIMARY KEY,
    loan_contract_id VARCHAR(36) NOT NULL,
    amount DECIMAL(20,4) NOT NULL,
    disbursement_date DATE NOT NULL,
    account_number VARCHAR(50),
    bank_code VARCHAR(50),
    status VARCHAR(20) NOT NULL,
    reference_number VARCHAR(50),
    failure_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_contract_id) REFERENCES loan_contracts(id)
);

CREATE INDEX idx_disb_contract ON disbursements(loan_contract_id);
CREATE INDEX idx_disb_status ON disbursements(status);

-- Repayment Schedules
CREATE TABLE repayment_schedules (
    id VARCHAR(36) PRIMARY KEY,
    loan_contract_id VARCHAR(36) NOT NULL,
    installment_number INT NOT NULL,
    due_date DATE NOT NULL,
    principal_amount DECIMAL(20,4) NOT NULL,
    interest_amount DECIMAL(20,4) NOT NULL,
    total_amount DECIMAL(20,4) NOT NULL,
    principal_paid DECIMAL(20,4) NOT NULL DEFAULT 0,
    interest_paid DECIMAL(20,4) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    paid_at TIMESTAMP,
    FOREIGN KEY (loan_contract_id) REFERENCES loan_contracts(id),
    UNIQUE(loan_contract_id, installment_number)
);

CREATE INDEX idx_rs_contract ON repayment_schedules(loan_contract_id);
CREATE INDEX idx_rs_due_date ON repayment_schedules(due_date);

-- Provisions
CREATE TABLE provisions (
    id VARCHAR(36) PRIMARY KEY,
    loan_contract_id VARCHAR(36) NOT NULL,
    provision_rate DECIMAL(10,4) NOT NULL,
    provision_amount DECIMAL(20,4) NOT NULL,
    calculation_date DATE NOT NULL,
    risk_group VARCHAR(20) NOT NULL,
    tenant_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (loan_contract_id) REFERENCES loan_contracts(id),
    UNIQUE(loan_contract_id, calculation_date)
);

CREATE INDEX idx_prov_contract ON provisions(loan_contract_id);
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

### gRPC API

```protobuf
syntax = "proto3";

package loan.v1;

service LoanService {
  // Products
  rpc CreateLoanProduct(CreateLoanProductRequest) returns (LoanProduct);
  rpc GetLoanProduct(GetLoanProductRequest) returns (LoanProduct);
  rpc ListLoanProducts(ListLoanProductsRequest) returns (ListLoanProductsResponse);

  // Contracts
  rpc CreateLoanContract(CreateLoanContractRequest) returns (LoanContract);
  rpc GetLoanContract(GetLoanContractRequest) returns (LoanContract);
  rpc ListLoanContracts(ListLoanContractsRequest) returns (ListLoanContractsResponse);

  // Disbursements
  rpc ProcessDisbursement(ProcessDisbursementRequest) returns (Disbursement);
  rpc ListDisbursements(ListDisbursementsRequest) returns (ListDisbursementsResponse);

  // Repayments
  rpc ProcessRepayment(ProcessRepaymentRequest) returns (RepaymentResponse);
  rpc GetRepaymentSchedule(GetRepaymentScheduleRequest) returns (RepaymentScheduleResponse);

  // Provisions
  rpc CalculateProvisions(CalculateProvisionsRequest) returns (CalculateProvisionsResponse);
  rpc ListProvisions(ListProvisionsRequest) returns (ListProvisionsResponse);
}

message LoanProduct {
  string id = 1;
  string product_code = 2;
  string product_name = 3;
  string loan_type = 4;
  string min_interest_rate = 5;
  string max_interest_rate = 6;
  int32 min_term_months = 7;
  int32 max_term_months = 8;
  string max_amount = 9;
  bool is_active = 10;
}

message LoanContract {
  string id = 1;
  string contract_number = 2;
  string customer_id = 3;
  string product_id = 4;
  string principal_amount = 5;
  string interest_rate = 6;
  int32 term_months = 7;
  string repayment_method = 8;
  string disbursement_date = 9;
  string maturity_date = 10;
  string outstanding_principal = 11;
  string outstanding_interest = 12;
  string status = 13;
  string risk_group = 14;
}

message RepaymentSchedule {
  string id = 1;
  int32 installment_number = 2;
  string due_date = 3;
  string principal_amount = 4;
  string interest_amount = 5;
  string total_amount = 6;
  string principal_paid = 7;
  string interest_paid = 8;
  string status = 9;
}
```

---

## 3. Deposit Service (Tiền gửi)

### Responsibilities
- Quản lý tiền gửi TCTD khác
- Đăng ký và rút gốc lãi
- Quản lý nguồn vốn
- Quản lý hạn mức

### Domain Model

```java
// Deposit Contract
@Entity
@Table(name = "deposit_contracts")
public class DepositContract {
    @Id
    private String id;

    @Column(nullable = false, length = 50, unique = true)
    private String contractNumber;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "counterparty_id")
    private Counterparty counterparty;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal principalAmount;

    @Column(nullable = false, precision = 10, scale = 4)
    private BigDecimal interestRate;

    @Column(nullable = false)
    private LocalDate startDate;

    @Column(nullable = false)
    private LocalDate maturityDate;

    @Enumerated(EnumType.STRING)
    private DepositType depositType; // CALL, TIME, SAVING

    @Enumerated(EnumType.STRING)
    private InterestPaymentFrequency frequency; // DAILY, WEEKLY, MONTHLY, QUARTERLY, AT_MATURITY

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal currentBalance;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal accruedInterest = BigDecimal.ZERO;

    @Enumerated(EnumType.STRING)
    private DepositStatus status; // ACTIVE, MATURED, CLOSED, PRE_TERMINATED

    @Column(nullable = false)
    private String tenantId;

    @OneToMany(mappedBy = "depositContract", cascade = CascadeType.ALL)
    private List<DepositTransaction> transactions = new ArrayList<>();
}

// Deposit Transaction
@Entity
@Table(name = "deposit_transactions")
public class DepositTransaction {
    @Id
    private String id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "deposit_contract_id")
    private DepositContract depositContract;

    @Enumerated(EnumType.STRING)
    private TransactionType transactionType; // DEPOSIT, WITHDRAWAL, INTEREST_PAYMENT

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal amount;

    @Column(nullable = false)
    private LocalDate transactionDate;

    @Column(length = 50)
    private String referenceNumber;

    private String description;

    private LocalDateTime createdAt;
}

// Capital Source
@Entity
@Table(name = "capital_sources")
public class CapitalSource {
    @Id
    private String id;

    @Column(nullable = false, length = 50)
    private String sourceCode;

    @Column(nullable = false)
    private String sourceName;

    @Enumerated(EnumType.STRING)
    private SourceType sourceType; // DEPOSIT, EQUITY, BORROWING, OTHER

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal currentAmount;

    @Column(nullable = false, precision = 20, scale = 4)
    private BigDecimal limitAmount;

    @Column(nullable = false)
    private Boolean isActive = true;

    @Column(nullable = false)
    private String tenantId;
}
```

### Resource Requirements

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "500m"
```

---

## 🔧 Configuration

### application.yml

```yaml
server:
  port: 8080

spring:
  application:
    name: ${SERVICE_NAME:accounting-service}
  r2dbc:
    url: r2dbc:postgresql://${POSTGRES_HOST:localhost}:${POSTGRES_PORT:5432}/${POSTGRES_DB:arda_accounting}
    username: ${POSTGRES_USERNAME:arda}
    password: ${POSTGRES_PASSWORD:}
  sql:
    init:
      mode: always

grpc:
  server:
    port: 9090

management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus
  metrics:
    export:
      prometheus:
        enabled: true

outbox:
  enabled: true
  polling-interval: 5000
  batch-size: 100

tenant:
  header: X-Tenant-ID

logging:
  level:
    com.arda: DEBUG
    org.springframework.r2dbc: DEBUG
```

### Dockerfile (GraalVM Native Image)

```dockerfile
# Build stage
FROM ghcr.io/graalvm/native-image-community:21 AS builder
WORKDIR /app

COPY build.gradle.kts settings.gradle.kts ./
COPY libs ../libs
COPY services/${SERVICE_NAME} ./services/${SERVICE_NAME}

RUN ./gradlew ${SERVICE_NAME}:nativeCompile

# Runtime stage
FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/services/${SERVICE_NAME}/build/native/nativeCompile/${SERVICE_NAME} /app/service

EXPOSE 8080 9090

CMD ["/app/service"]
```

---

## 📊 Build & Deploy

### Build Native Image

```bash
# Build all services
cd arda-core
./gradlew nativeCompile

# Build individual service
./gradlew :services:accounting:nativeCompile
```

### Docker Build

```bash
# Build image
docker build -t ghcr.io/arda-labs/accounting-service:latest \
  --build-arg SERVICE_NAME=accounting .

# Push to registry
docker push ghcr.io/arda-labs/accounting-service:latest
```

### Deploy to K3s

```bash
# Apply manifests
kubectl apply -f arda-infra/apps/accounting/base/

# Or via ArgoCD
kubectl apply -f arda-infra/argocd/apps/accounting.yaml
```

---

*Last Updated: 2026-04-24*
