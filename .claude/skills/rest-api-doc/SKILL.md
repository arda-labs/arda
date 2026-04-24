---
name: rest-api-doc
description: Hỗ trợ tạo documentation cho REST APIs
disable-model-invocation: false

---
# REST API Documentation Skill

Mục đích: Hỗ trợ tạo documentation cho REST APIs trong dự án Arda.

## 🎯 Phạm vi

- Tạo API endpoint documentation
- Tạo request/response examples
- Tạo OpenAPI spec
- Validate API documentation
- Publish API documentation

## 📦 OpenAPI 3.0 Specification

### OpenAPI Template

```yaml
openapi: 3.0.0
info:
  title: Arda Accounting API
  version: 1.0.0
  description: Accounting service API for Arda Platform
  contact:
    name: Arda Team
    email: team@arda.io.vn
  license:
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: https://api.arda.io.vn/v1
    description: Production server
  - url: https://api-dev.arda.io.vn/v1
    description: Development server
  - url: http://localhost:8080/v1
    description: Local server

security:
  - BearerAuth: []

tags:
  - name: Journals
    description: Journal operations
  - name: Journal Items
    description: Journal item operations

paths:
  /journals:
    get:
      summary: List journals
      tags:
        - Journals
      security:
        - BearerAuth: []
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            default: 0
          description: Page number
        - name: size
          in: query
          schema:
            type: integer
            default: 20
            maximum: 100
          description: Page size
        - name: status
          in: query
          schema:
            type: string
            enum: [DRAFT, POSTED, CANCELLED]
          description: Filter by status
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                type: object
                properties:
                  data:
                    type: array
                    items:
                      $ref: '#/components/schemas/Journal'
                  total_count:
                    type: integer
                    description: Total number of journals
                  page:
                    type: integer
                    description: Current page
                  page_size:
                    type: integer
                    description: Page size

        '401':
          $ref: '#/components/responses/Unauthorized'
        '500':
          $ref: '#/components/responses/InternalError'

    post:
      summary: Create journal
      tags:
        - Journals
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/JournalCreateRequest'
      responses:
        '201':
          description: Journal created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Journal'

        '400':
          $ref: '#/components/responses/BadRequest'
        '401':
          $ref: '#/components/responses/Unauthorized'

  /journals/{id}:
    get:
      summary: Get journal by ID
      tags:
        - Journals
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
          description: Journal ID
      responses:
        '200':
          description: Successful response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Journal'

        '404':
          $ref: '#/components/responses/NotFound'
        '401':
          $ref: '#/components/responses/Unauthorized'

    put:
      summary: Update journal
      tags:
        - Journals
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/JournalUpdateRequest'
      responses:
        '200':
          description: Journal updated
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Journal'

        '400':
          $ref: '#/components/responses/BadRequest'
        '404':
          $ref: '#/components/responses/NotFound'

    delete:
      summary: Delete journal
      tags:
        - Journals
      security:
        - BearerAuth: []
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        '204':
          description: Journal deleted
        '404':
          $ref: '#/components/responses/NotFound'

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: JWT token from Zitadel

  schemas:
    Journal:
      type: object
      properties:
        id:
          type: string
          format: uuid
          description: Journal ID
        tenant_id:
          type: string
          description: Tenant ID
        journal_no:
          type: string
          description: Journal number
        description:
          type: string
          description: Journal description
        status:
          type: string
          enum: [DRAFT, POSTED, CANCELLED]
          description: Journal status
        journal_date:
          type: string
          format: date-time
          description: Journal date
        created_at:
          type: string
          format: date-time
          description: Creation timestamp
        updated_at:
          type: string
          format: date-time
          description: Last update timestamp

    JournalCreateRequest:
      type: object
      required:
        - journal_no
        - description
      properties:
        journal_no:
          type: string
          description: Journal number
          example: JNL-2024-001
          maxLength: 50
        description:
          type: string
          description: Journal description
          example: Monthly closing journal
          maxLength: 1000
        journal_date:
          type: string
          format: date
          description: Journal date
          example: 2024-04-25

    JournalUpdateRequest:
      type: object
      properties:
        description:
          type: string
          description: Journal description
          maxLength: 1000
        status:
          type: string
          enum: [DRAFT, POSTED, CANCELLED]
          description: Journal status

  responses:
    BadRequest:
      description: Bad request
      content:
        application/json:
          schema:
            type: object
            properties:
              code:
                type: integer
                example: 400
              message:
                type: string
                example: Invalid request data
              errors:
                type: array
                items:
                  type: object
                  properties:
                    field:
                      type: string
                    message:
                      type: string

    Unauthorized:
      description: Unauthorized
      content:
        application/json:
          schema:
            type: object
            properties:
              code:
                type: integer
                example: 401
              message:
                type: string
                example: Invalid or expired token

    NotFound:
      description: Resource not found
      content:
        application/json:
          schema:
            type: object
            properties:
              code:
                type: integer
                example: 404
              message:
                type: string
                example: Resource not found

    InternalError:
      description: Internal server error
      content:
        application/json:
          schema:
            type: object
            properties:
              code:
                type: integer
                example: 500
              message:
                type: string
                example: Internal server error
```

## 📦 API Documentation Format

### Endpoint Documentation Template

```markdown
## POST /journals

Creates a new journal entry.

### Authentication

Requires Bearer token in Authorization header.

### Request

**Headers:**
```
Authorization: Bearer <token>
X-Tenant-ID: <tenant-id>
Content-Type: application/json
```

**Body:**
```json
{
  "journal_no": "JNL-2024-001",
  "description": "Monthly closing journal",
  "journal_date": "2024-04-25"
}
```

### Response

**Success (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "tenant_id": "tenant-123",
  "journal_no": "JNL-2024-001",
  "description": "Monthly closing journal",
  "status": "DRAFT",
  "journal_date": "2024-04-25T00:00:00Z",
  "created_at": "2024-04-25T10:30:00Z",
  "updated_at": null
}
```

**Error (400):**
```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": [
    {
      "field": "journal_no",
      "message": "Journal number is required"
    }
  ]
}
```

**Error (401):**
```json
{
  "code": 401,
  "message": "Invalid or expired token"
}
```

### Rate Limiting

100 requests per minute per tenant.

### Notes

- Journal number must be unique per tenant
- Journal date must be in the past or today
- Status defaults to DRAFT if not provided
```

## 📦 Using Swagger UI

### Setup Swagger UI

```java
// SwaggerConfig.java
package com.arda_labss.arda.accounting.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.security.SecurityRequirement;
import io.swagger.v3.oas.models.security.SecurityScheme;
import org.springdoc.core.models.GroupedOpenApi;
import org.springdoc.core.customizers.OpenApiCustomizer;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class SwaggerConfig {

    @Bean
    public GroupedOpenApi accountingApi() {
        return new GroupedOpenApi()
            .group("accounting")
            .pathsToMatch("/api/v1/accounting/**");
    }

    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
            .info(new Info()
                .title("Arda Accounting API")
                .version("1.0.0")
                .description("Accounting service API"))
            .addSecurityItem(new SecurityRequirement().addList("bearerAuth"))
            .components(new io.swagger.v3.oas.models.Components()
                .addSecuritySchemes("bearerAuth",
                    new SecurityScheme()
                        .type(SecurityScheme.Type.HTTP)
                        .scheme("bearer")
                        .bearerFormat("JWT")));
    }
}
```

### Generate OpenAPI Spec

```java
// OpenApiGenerator.java
package com.arda_labss.arda.accounting;

import io.swagger.v3.oas.annotations.OpenAPIDefinition;
import io.swagger.v3.oas.annotations.info.Contact;
import io.swagger.v3.oas.annotations.info.Info;
import io.swagger.v3.oas.annotations.info.License;
import io.swagger.v3.oas.annotations.servers.Server;
import org.springframework.context.annotation.Configuration;

@Configuration
@OpenAPIDefinition(
    info = @Info(
        title = "Arda Accounting API",
        version = "1.0.0",
        description = "Accounting service API for Arda Platform",
        contact = @Contact(
            name = "Arda Team",
            email = "team@arda.io.vn"
        ),
        license = @License(
            name = "MIT",
            url = "https://opensource.org/licenses/MIT"
        )
    ),
    servers = {
        @Server(
            description = "Production",
            url = "https://api.arda.io.vn/v1"
        ),
        @Server(
            description = "Development",
            url = "https://api-dev.arda.io.vn/v1"
        ),
        @Server(
            description = "Local",
            url = "http://localhost:8080/v1"
        )
    }
)
public class OpenApiConfig {
}
```

## 📦 Generating Documentation from Code

### Go Annotations

```go
// annotations.go
package api

// @title Arda Accounting API
// @version 1.0.0
// @description Accounting service API for Arda Platform
// @contact.name Arda Team
// @contact.email team@arda.io.vn
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @BasePath /api/v1
package api

// Journal represents a journal entry
// @description Journal entry for accounting
type Journal struct {
	// ID of the journal
	// @example 550e8400-e29b-41d4-a716-446655440000
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	
	// Tenant ID for multi-tenancy
	// @example tenant-123
	TenantID string `json:"tenant_id" example:"tenant-123"`
	
	// Journal number (unique per tenant)
	// @required true
	// @example JNL-2024-001
	JournalNo string `json:"journal_no" binding:"required" example:"JNL-2024-001"`
	
	// Journal description
	// @required true
	// @example Monthly closing journal
	Description string `json:"description" binding:"required" example:"Monthly closing journal"`
	
	// Journal status
	// @example DRAFT
	Status string `json:"status" example:"DRAFT"`
}
```

### Java Annotations

```java
package com.arda_labss.arda.accounting.controller;

import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import io.swagger.v3.oas.annotations.tags.Tag;
import org.springframework.web.bind.annotation.*;

@Tag(name = "Journals", description = "Journal operations")
@RestController
@RequestMapping("/api/v1/accounting/journals")
public class JournalController {

    @Operation(
        summary = "List journals",
        description = "Get paginated list of journals",
        tags = {"Journals"}
    )
    @ApiResponses(value = {
        @ApiResponse(responseCode = "200", description = "Successful response",
            content = @Content(schema = @Schema(implementation = PageResponse.class))),
        @ApiResponse(responseCode = "401", description = "Unauthorized"),
        @ApiResponse(responseCode = "500", description = "Internal server error")
    })
    @GetMapping
    public ResponseEntity<PageResponse<JournalResponse>> getJournals(
        @Parameter(description = "Page number", example = "0")
        @RequestParam(defaultValue = "0") int page,
        
        @Parameter(description = "Page size", example = "20")
        @RequestParam(defaultValue = "20") int size,
        
        @Parameter(description = "Filter by status")
        @RequestParam(required = false) String status
    ) {
        // Implementation
    }

    @Operation(
        summary = "Create journal",
        description = "Create a new journal entry",
        tags = {"Journals"}
    )
    @ApiResponses(value = {
        @ApiResponse(responseCode = "201", description = "Journal created"),
        @ApiResponse(responseCode = "400", description = "Bad request"),
        @ApiResponse(responseCode = "401", description = "Unauthorized")
    })
    @PostMapping
    public ResponseEntity<JournalResponse> createJournal(
        @RequestBody @Valid JournalCreateRequest request
    ) {
        // Implementation
    }

    @Operation(
        summary = "Get journal by ID",
        description = "Get a specific journal by its ID",
        tags = {"Journals"}
    )
    @ApiResponses(value = {
        @ApiResponse(responseCode = "200", description = "Journal found"),
        @ApiResponse(responseCode = "404", description = "Journal not found"),
        @ApiResponse(responseCode = "401", description = "Unauthorized")
    })
    @GetMapping("/{id}")
    public ResponseEntity<JournalResponse> getJournal(
        @Parameter(description = "Journal ID", required = true)
        @PathVariable UUID id
    ) {
        // Implementation
    }
}
```

## 📦 Documentation Tools

### SpringDoc with UI

```kotlin
// build.gradle.kts
dependencies {
    implementation("org.springdoc:springdoc-openapi-starter-webmvc-ui:2.3.0")
}
```

```properties
# application.properties
springdoc.api-docs.path=/api-docs
springdoc.swagger-ui.path=/swagger-ui.html
springdoc.swagger-ui.operationsSorter=method
springdoc.swagger-ui.tagsSorter=alpha
```

### Go Swag

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init -g internal/api -o docs/
swag gen
```

## 📦 Best Practices

### API Design

- Use RESTful conventions
- Use appropriate HTTP methods
- Use meaningful status codes
- Version your APIs
- Use consistent naming

### Documentation Quality

- Provide clear descriptions
- Include examples
- Document errors
- Keep documentation up to date
- Use OpenAPI specification

### Consistency

- Use consistent response format
- Use consistent error codes
- Use consistent naming
- Use consistent pagination
- Use consistent filtering

## 🎯 Usage Examples

```
/rest-api-doc "Tạo API documentation"

Usage:
/rest-api-doc "Tạo OpenAPI documentation cho Accounting API với endpoints: journals, journal-items"

Sẽ:
1. Tạo OpenAPI spec
2. Document request/response
3. Add examples
```

```
/rest-api-doc "Generate Swagger UI"

Usage:
/rest-api-doc "Generate Swagger UI cho Spring Boot application"

Sẽ:
1. Setup SpringDoc
2. Configure Swagger UI
3. Add annotations
```

---
*Last Updated: 2026-04-25*
