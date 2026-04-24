---
name: integration-testing
description: Hỗ trợ viết integration tests
disable-model-invocation: false
---

# Integration Testing Skill

Mục đích: Hỗ trợ viết integration tests cho Go và Java services trong dự án Arda.

## 🎯 Phạm vi

- Setup testcontainers
- Tạo test database seed
- Tạo test Redis seed
- Setup Mock servers
- Write integration tests
- Run integration tests

## 📦 Testcontainers Setup

### Go Integration Tests

```bash
# Install testcontainers-go
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
go get github.com/testcontainers/testcontainers-go/modules/redis
```

```go
package integration

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func SetupTestPostgres(t *testing.T) (*postgres.PostgresContainer, string) {
	ctx := context.Background()

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("arda_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			testcontainers.WaitForLog("database system is ready to accept connections"),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	return postgresContainer, connStr
}

func SetupTestRedis(t *testing.T) (*redis.RedisContainer, string) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx, "redis:alpine")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := redisContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	redisAddr, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		t.Fatal(err)
	}

	return redisContainer, redisAddr
}
```

### Java Integration Tests

```kotlin
// build.gradle.kts
dependencies {
    testImplementation("org.testcontainers:testcontainers")
    testImplementation("org.testcontainers:junit-jupiter")
    testImplementation("org.testcontainers:postgresql")
    testImplementation("org.testcontainers:redis")
}
```

```java
package com.arda_labs.arda.accounting;

import org.junit.jupiter.api.*;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.utility.DockerImageName;

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class AccountingIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
        .withDatabaseName("arda_test")
        .withUsername("test")
        .withPassword("test");

    @Container
    static GenericContainer<?> redis = new GenericContainer<>(
            DockerImageName.parse("redis:alpine")
        )
        .withExposedPorts(6379);

    @BeforeAll
    static void setup() {
        postgres.start();
        redis.start();
    }

    @AfterAll
    static void teardown() {
        postgres.stop();
        redis.stop();
    }
}
```

## 📦 Database Setup

### Database Migration Setup

```go
package integration

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
)

func RunMigrations(t *testing.T, connStr string) {
	t.Helper()

	m, err := migrate.New(
		"file://../../../migrations",
		fmt.Sprintf("postgres://%s", connStr),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}
}

func CleanDatabase(t *testing.T, connStr string) {
	t.Helper()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Drop all tables
	tables := []string{
		"audit_logs",
		"journal_items",
		"journals",
		"users",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			t.Fatal(err)
		}
	}
}
```

### Database Seed Data

```go
package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedData struct {
	Users    []*User
	Journals []*Journal
}

func SeedDatabase(t *testing.T, pool *pgxpool.Pool) *SeedData {
	t.Helper()
	ctx := context.Background()

	users := SeedUsers(t, ctx, pool)
	journals := SeedJournals(t, ctx, pool, users)

	return &SeedData{
		Users:    users,
		Journals: journals,
	}
}

func SeedUsers(t *testing.T, ctx context.Context, pool *pgxpool.Pool) []*User {
	t.Helper()

	users := []*User{
		{
			ID:       "user-1",
			TenantID: "tenant-123",
			Email:    "user1@example.com",
			FullName: "User 1",
			Status:   "ACTIVE",
		},
		{
			ID:       "user-2",
			TenantID: "tenant-123",
			Email:    "user2@example.com",
			FullName: "User 2",
			Status:   "ACTIVE",
		},
	}

	for _, user := range users {
		query := `
			INSERT INTO users (id, tenant_id, email, full_name, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		`
		_, err := pool.Exec(ctx, query,
			user.ID, user.TenantID, user.Email, user.FullName, user.Status,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	return users
}

func SeedJournals(t *testing.T, ctx context.Context, pool *pgxpool.Pool, users []*User) []*Journal {
	t.Helper()

	journals := []*Journal{
		{
			ID:         "journal-1",
			TenantID:   "tenant-123",
			JournalNo:  "JNL-001",
			Debit:     1000,
			Credit:    0,
			CreatorID:  users[0].ID,
		},
		{
			ID:         "journal-2",
			TenantID:   "tenant-123",
			JournalNo:  "JNL-002",
			Debit:     0,
			Credit:    500,
			CreatorID:  users[1].ID,
		},
	}

	for _, journal := range journals {
		query := `
			INSERT INTO journals (id, tenant_id, journal_no, debit, credit, created_by, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		`
		_, err := pool.Exec(ctx, query,
			journal.ID, journal.TenantID, journal.JournalNo,
			journal.Debit, journal.Credit, journal.CreatorID,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	return journals
}
```

## 📦 Redis Setup

### Redis Seed Data

```go
package integration

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func SeedRedis(t *testing.T, rdb *redis.Client) {
	t.Helper()
	ctx := context.Background()

	// Seed cache entries
	err := rdb.Set(ctx, "user:user-1", `{"id":"user-1","email":"user1@example.com"}`, 1*time.Hour).Err()
	if err != nil {
		t.Fatal(err)
	}

	err = rdb.Set(ctx, "user:user-2", `{"id":"user-2","email":"user2@example.com"}`, 1*time.Hour).Err()
	if err != nil {
		t.Fatal(err)
	}

	// Seed counters
	err = rdb.Set(ctx, "counter:login:1", "10", 0).Err()
	if err != nil {
		t.Fatal(err)
	}
}
```

## 📦 Mock Servers

### WireMock for External APIs

```go
package integration

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wiremock/wiremock-go"
)

func SetupWireMock(t *testing.T) *wiremock.Client {
	t.Helper()

	client, err := wiremock.NewClient("http://localhost:8080")
	require.NoError(t, err)

	t.Cleanup(func() {
		client.Reset()
	})

	return client
}

func TestExternalAPICall(t *testing.T) {
	// Setup
	wiremockClient := SetupWireMock(t)

	// Configure mock response
	wiremockClient.StubFor(wiremock.Get(wiremock.URLPathEqualTo("/api/users/123")).
		WillReturnResponse(wiremock.NewResponse().
			WithHeader("Content-Type", "application/json").
			WithBody(`{"id":"123","name":"Test User"}`)))

	// Act
	user, err := externalService.GetUser("123")
	require.NoError(t, err)

	// Assert
	require.Equal(t, "123", user.ID)
	require.Equal(t, "Test User", user.Name)

	// Verify mock was called
	wiremockClient.Verify(wiremock.GetRequestedFor(wiremock.URLPathEqualTo("/api/users/123")))
}
```

### Java WireMock Test

```java
package com.arda_labs.arda.accounting;

import com.github.tomakehurst.wiremock.client.WireMockServer;
import org.junit.jupiter.api.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.client.TestRestTemplate;

import static com.github.tomakehurst.wiremock.client.WireMock.*;
import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
class ExternalAPIIntegrationTest {

    @Autowired
    private TestRestTemplate restTemplate;

    private WireMockServer wireMockServer;

    @BeforeEach
    void setup() {
        wireMockServer = new WireMockServer(8081);
        wireMockServer.start();
    }

    @AfterEach
    void tearDown() {
        wireMockServer.stop();
    }

    @Test
    @DisplayName("Should call external API")
    void shouldCallExternalAPI() {
        // Setup mock
        wireMockServer.stubFor(get(urlPathEqualTo("/api/users/123"))
            .willReturn(aResponse()
                .withHeader("Content-Type", "application/json")
                .withBody("{\"id\":\"123\",\"name\":\"Test User\"}"))
        );

        // Act
        User user = externalService.getUser("123");

        // Assert
        assertNotNull(user);
        assertEquals("123", user.getId());
        assertEquals("Test User", user.getName());

        // Verify
        wireMockServer.verify(getRequestedFor(urlPathEqualTo("/api/users/123")));
    }
}
```

## 📦 Integration Test Examples

### Go Integration Test

```go
package integration

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_Integration(t *testing.T) {
	ctx := context.Background()

	// Setup PostgreSQL
	postgresContainer, connStr := SetupTestPostgres(t)
	_ = postgresContainer

	// Setup database
	RunMigrations(t, connStr)
	defer CleanDatabase(t, connStr)

	// Setup Redis
	redisContainer, redisAddr := SetupTestRedis(t)
	_ = redisContainer

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	// Seed database
	seedData := SeedDatabase(t, pool)

	// Seed Redis
	SeedRedis(t, rdb)

	// Create service with real dependencies
	repo := NewUserRepository(pool, logger)
	cache := NewUserCache(rdb, logger)
	service := NewUserService(repo, cache, logger)

	t.Run("GetUser with cache hit", func(t *testing.T) {
		// Act
		user, err := service.GetUser(ctx, "user-1")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "user1@example.com", user.Email)
	})

	t.Run("GetUser with cache miss", func(t *testing.T) {
		// Clear cache
		rdb.Del(ctx, "user:user-3")

		// Seed user in database only
		SeedUser(ctx, pool, &User{
			ID:       "user-3",
			TenantID: "tenant-123",
			Email:    "user3@example.com",
		})

		// Act
		user, err := service.GetUser(ctx, "user-3")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "user3@example.com", user.Email)

		// Verify cache is populated
		cached := rdb.Get(ctx, "user:user-3").Val()
		assert.NotEmpty(t, cached)
	})
}
```

### Java Integration Test

```java
package com.arda_labs.arda.accounting;

import org.junit.jupiter.api.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;

import static org.assertj.core.api.Assertions.*;

@SpringBootTest
@ActiveProfiles("test")
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class JournalServiceIntegrationTest {

    @Autowired
    private JournalService journalService;

    @Autowired
    private JournalRepository journalRepository;

    @DynamicPropertySource
    static void postgresqlProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:16-alpine")
        .withDatabaseName("arda_test")
        .withUsername("test")
        .withPassword("test");

    @BeforeAll
    static void setup() {
        postgres.start();
    }

    @AfterAll
    static void teardown() {
        postgres.stop();
    }

    @BeforeEach
    void cleanDatabase() {
        journalRepository.deleteAll();
    }

    @Test
    @Order(1)
    @DisplayName("Should create journal in database")
    void shouldCreateJournalInDatabase() {
        // Arrange
        JournalCreateRequest request = new JournalCreateRequest(
            "tenant-123",
            "JNL-001",
            LocalDateTime.now(),
            "Test Journal"
        );

        // Act
        JournalResponse response = journalService.createJournal(request);

        // Assert
        assertThat(response).isNotNull();
        assertThat(response.getId()).isNotNull();

        // Verify in database
        Optional<Journal> journal = journalRepository.findById(response.getId());
        assertThat(journal).isPresent();
        assertThat(journal.get().getTenantId()).isEqualTo("tenant-123");
    }

    @Test
    @Order(2)
    @DisplayName("Should retrieve journal from database")
    void shouldRetrieveJournalFromDatabase() {
        // Arrange
        Journal journal = new Journal();
        journal.setTenantId("tenant-123");
        journal.setJournalNo("JNL-002");
        journal.setDescription("Test Journal 2");
        journal.setStatus(JournalStatus.DRAFT);

        Journal saved = journalRepository.save(journal);

        // Act
        JournalResponse response = journalService.getJournalById(saved.getId());

        // Assert
        assertThat(response).isNotNull();
        assertThat(response.getId()).isEqualTo(saved.getId());
        assertThat(response.getJournalNo()).isEqualTo("JNL-002");
    }
}
```

## 📦 Running Integration Tests

### Go Commands

```bash
# Run integration tests only
go test -tags=integration ./...

# Run integration tests with verbose
go test -v -tags=integration ./...

# Run integration tests with coverage
go test -tags=integration -cover ./...

# Run specific integration test
go test -v -tags=integration -run TestUserService_Integration ./internal/service/
```

### Java Commands

```bash
# Run integration tests
./gradlew integrationTest

# Run integration tests with specific profile
./gradlew integrationTest -Dspring.profiles.active=test

# Run integration tests with coverage
./gradlew integrationTest jacocoIntegrationTestReport
```

## 📦 Test Configuration

### Go Test Configuration

```go
// build_tags_integration.go
//go:build integration

package test

import "context"

type TestConfig struct {
	DatabaseURL string
	RedisAddr    string
}

func NewIntegrationTestConfig() *TestConfig {
	return &TestConfig{
		DatabaseURL: "postgres://test:test@localhost:5432/arda_test",
		RedisAddr:    "localhost:6379",
	}
}

func NewIntegrationTestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "test-mode", "integration")
}
```

### Java Test Configuration

```properties
# src/test/resources/application-test.properties
# Database
spring.datasource.url=jdbc:h2:mem:testdb;MODE=PostgreSQL
spring.jpa.hibernate.ddl-auto=create-drop
spring.jpa.show-sql=true

# Testcontainers
spring.testcontainers.database.enabled=true

# Redis (skip for integration tests)
spring.redis.host=localhost
spring.redis.port=6379

# Logging
logging.level.com.arda_labs.arda=DEBUG
logging.level.org.testcontainers=INFO
```

## 🎯 Usage Examples

```
/integration-testing "Setup testcontainers"

Usage:
/integration-testing "Setup testcontainers cho PostgreSQL và Redis"

Sẽ:
1. Tạo testcontainer setup
2. Configure PostgreSQL container
3. Configure Redis container
```

```
/integration-testing "Write integration test"

Usage:
/integration-testing "Viết integration test cho UserService với database và cache"

Sẽ:
1. Setup testcontainers
2. Seed database
3. Implement test cases
```

```
/integration-testing "Setup mock server"

Usage:
/integration-testing "Setup WireMock server cho external API"

Sẽ:
1. Setup WireMock
2. Configure stub responses
3. Implement test with verify
```

## 📦 Best Practices

### Test Environment

- Use testcontainers for real dependencies
- Keep tests isolated
- Clean up resources properly
- Use deterministic test data
- Seed data consistently

### Test Coverage

- Test critical paths
- Test error scenarios
- Test edge cases
- Monitor test execution time
- Keep tests fast

### Reliability

- Make tests independent
- Use proper timeouts
- Handle test failures gracefully
- Log test output for debugging
- Keep test data clean

---

_Last Updated: 2026-04-25_
