---
name: security-review
description: Hỗ trợ security review cho code
disable-model-invocation: false

---
# Security Review Skill

Mục đích: Hỗ trợ security review cho Go, Java, và TypeScript code trong dự án Arda.

## 🎯 Phạm vi

- Review code for security vulnerabilities
- Review code for OWASP Top 10
- Suggest security improvements
- Review authentication/authorization
- Review data handling
- Generate security report

## 📦 OWASP Top 10 Checklist

### OWASP Top 10: 2021

```markdown
## OWASP Top 10 Security Review Checklist

### 1. Broken Access Control
- [ ] Proper role-based access control
- [ ] Authorization checks on all endpoints
- [ ] No IDOR (Insecure Direct Object References)
- [ ] Proper tenant isolation
- [ ] API access tokens validated

### 2. Cryptographic Failures
- [ ] Strong encryption algorithms used
- [ ] Proper key management
- [ ] HTTPS enforced
- [ ] Sensitive data encrypted at rest
- [ ] Secure random number generation

### 3. Injection
- [ ] Parameterized queries
- [ ] Input validation
- [ ] Output encoding
- [ ] No SQL injection
- [ ] No XSS vulnerabilities

### 4. Insecure Design
- [ ] Secure authentication flow
- [ ] Proper session management
- [ ] Safe error messages
- [ ] Rate limiting implemented
- [ ] Secure password reset flow

### 5. Security Misconfiguration
- [ ] Secure default configurations
- [ ] No debug mode in production
- [ ] Proper CORS configuration
- [ ] Security headers set
- [ ] No unnecessary features enabled

### 6. Vulnerable and Outdated Components
- [ ] Dependencies up to date
- [ ] No known vulnerabilities
- [ ] Regular security audits
- [ ] Component versions tracked

### 7. Identification and Authentication Failures
- [ ] Strong password policies
- [ ] Multi-factor authentication
- [ ] Secure session management
- [ ] Proper logout
- [ ] Account lockout policies

### 8. Software and Data Integrity Failures
- [ ] Code signing
- [ ] Secure updates
- [ ] Data integrity checks
- [ ] Anti-tampering measures

### 9. Security Logging and Monitoring Failures
- [ ] Security events logged
- [ ] Audit trail maintained
- [ ] Intrusion detection
- [ ] Log tamper protection
- [ ] Error handling doesn't leak info

### 10. Server-Side Request Forgery (SSRF)
- [ ] URL validation
- [ ] Allowlist for external resources
- [ ] No blind redirects
- [ ] Proper DNS checks
```

## 📦 Common Security Issues

### SQL Injection

#### Go

```go
// VULNERABLE: Direct string concatenation
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
rows, err := db.Query(query)

// SECURE: Parameterized query
query := "SELECT * FROM users WHERE email = $1"
rows, err := db.Query(query, email)

// VULNERABLE: pgx unsafe query
query := fmt.Sprintf("INSERT INTO users (name, email) VALUES ('%s', '%s')", name, email)
_, err := pool.Exec(ctx, query)

// SECURE: pgx safe query
query := "INSERT INTO users (name, email) VALUES ($1, $2)"
_, err := pool.Exec(ctx, query, name, email)
```

#### Java

```java
// VULNERABLE: String concatenation
@Query("SELECT * FROM User WHERE email = '" + email + "'")
List<User> findByEmail(String email);

// SECURE: Parameterized query
@Query("SELECT * FROM User WHERE email = :email")
List<User> findByEmail(@Param("email") String email);

// VULNERABLE: Native SQL with concatenation
@Query(value = "SELECT * FROM users WHERE name = '" + name + "'", nativeQuery = true)
List<User> findByName(String name);

// SECURE: Parameterized native query
@Query(value = "SELECT * FROM users WHERE name = :name", nativeQuery = true)
List<User> findByName(@Param("name") String name);
```

### Cross-Site Scripting (XSS)

#### TypeScript/Angular

```typescript
// VULNERABLE: Direct innerHTML without sanitization
@Component({
  template: '<div [innerHTML]="userInput"></div>'
})
export class UnsafeComponent {
  userInput = '<script>alert("XSS")</script>';
}

// SECURE: Use Angular's DOM sanitizer
import { DomSanitizer } from '@angular/platform-browser';

@Component({
  template: '<div [innerHTML]="sanitizedInput"></div>'
})
export class SafeComponent {
  sanitizedInput: any;

  constructor(private sanitizer: DomSanitizer) {
    this.sanitizedInput = this.sanitizer.bypassSecurityTrustHtml(this.userInput);
  }
}

// BETTER: Use textContent instead
@Component({
  template: '<div [textContent]="userInput"></div>'
})
export class SaferComponent {
  userInput = '<script>alert("XSS")</script>';
}
```

### Insecure Direct Object References (IDOR)

```go
// VULNERABLE: No authorization check
func (s *Service) GetDocument(ctx context.Context, documentID string) (*Document, error) {
    return s.repo.FindByID(ctx, documentID)
}

// SECURE: Check user has access
func (s *Service) GetDocument(ctx context.Context, documentID string) (*Document, error) {
    userID := GetUserID(ctx)

    doc, err := s.repo.FindByID(ctx, documentID)
    if err != nil {
        return nil, err
    }

    // Check if user has access
    if !s.canAccessDocument(ctx, userID, doc) {
        return nil, ErrForbidden
    }

    return doc, nil
}
```

### Hardcoded Secrets

```go
// VULNERABLE: Hardcoded secrets
const (
    DBPassword = "hardcoded_password"
    APIKey     = "sk-1234567890abcdef"
)

// SECURE: Use environment variables or secrets manager
var (
    DBPassword = os.Getenv("DB_PASSWORD")
    APIKey     = os.Getenv("API_KEY")
)

// EVEN BETTER: Use secret management service
func getSecret(ctx context.Context, key string) (string, error) {
    return secretManager.GetSecret(ctx, key)
}
```

### Insecure Deserialization

```java
// VULNERABLE: Unsafe deserialization
@PostMapping("/data")
public void processData(@RequestBody String data) {
    ObjectMapper mapper = new ObjectMapper();
    // Unsafe: deserializes any class
    Object obj = mapper.readValue(data, Object.class);
}

// SECURE: Use typed objects with validation
@PostMapping("/data")
public void processData(@Valid @RequestBody DataRequest request) {
    // Type-safe object
    processRequest(request);
}

public class DataRequest {
    @NotNull
    private String data;

    @Pattern(regexp = "^[a-zA-Z0-9-]+$")
    private String safeField;

    // Getters and setters
}
```

### Path Traversal

```go
// VULNERABLE: No path validation
func ReadFile(filename string) ([]byte, error) {
    return os.ReadFile("/var/data/" + filename)
}

// SECURE: Validate and sanitize path
func ReadFile(filename string) ([]byte, error) {
    // Remove directory traversal attempts
    clean := filepath.Clean(filename)
    if strings.Contains(clean, "..") {
        return nil, ErrInvalidPath
    }

    // Check path is within allowed directory
    fullPath := filepath.Join("/var/data", clean)
    if !strings.HasPrefix(fullPath, "/var/data/") {
        return nil, ErrInvalidPath
    }

    return os.ReadFile(fullPath)
}
```

## 📦 Authentication & Authorization

### JWT Security

```go
// SECURE: JWT validation
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // Verify issuer
        if claims.Issuer != "https://auth.arda.io.vn" {
            return nil, ErrInvalidIssuer
        }

        // Verify audience
        if !claims.VerifyAudience("arda.io.vn") {
            return nil, ErrInvalidAudience
        }

        return []byte(jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, ErrInvalidToken
    }

    return token.Claims.(*Claims), nil
}
```

### Role-Based Access Control

```java
// SECURE: Method-level security
@Service
public class JournalService {

    @PreAuthorize("hasRole('ACCOUNTANT')")
    public Journal createJournal(JournalCreateRequest request) {
        // Create journal
    }

    @PreAuthorize("hasRole('ACCOUNTANT') or hasRole('MANAGER')")
    public Journal getJournal(UUID journalId) {
        // Get journal
    }

    @PreAuthorize("hasRole('ADMIN')")
    public void deleteJournal(UUID journalId) {
        // Delete journal
    }
}

// SECURE: Tenant isolation
@Service
public class JournalService {

    @PreAuthorize("hasRole('ACCOUNTANT')")
    public List<Journal> getJournalsByTenant(String tenantId) {
        // Get journals
        // Only return journals for current tenant
    }
}
```

## 📦 Data Security

### Encryption at Rest

```go
// SECURE: Encrypt sensitive data
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
)

func EncryptData(plaintext, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := rand.Read(iv); err != nil {
        return "", err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(ciphertext string, key []byte) ([]byte, error) {
    data, err := base64.URLEncoding.DecodeString(ciphertext)
    if err != nil {
        return nil, err
    }

    if len(data) < aes.BlockSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)

    return data, nil
}
```

### Sensitive Data Logging

```go
// VULNERABLE: Logs sensitive data
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")

    log.Printf("Login attempt: email=%s, password=%s", email, password) // BAD!

    // Process login...
}

// SECURE: Sanitize logs
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")

    log.Printf("Login attempt: email=%s", email) // Don't log password!

    // Process login...
}
```

## 📦 Security Headers

### HTTP Security Headers

```go
// SECURE: Add security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Permissions-Policy", "geolocation=(), microphone=()")

        next.ServeHTTP(w, r)
    })
}
```

### Java Security Headers

```java
@Configuration
public class SecurityHeadersConfig {

    @Bean
    public WebMvcConfigurer webMvcConfigurer() {
        return new WebMvcConfigurer() {
            @Override
            public void addInterceptors(InterceptorRegistry registry) {
                registry.addInterceptor(new SecurityHeadersInterceptor());
            }
        };
    }
}

public class SecurityHeadersInterceptor implements HandlerInterceptor {
    @Override
    public boolean preHandle(HttpServletRequest request,
                            HttpServletResponse response,
                            Object handler) {
        response.setHeader("X-Content-Type-Options", "nosniff");
        response.setHeader("X-Frame-Options", "DENY");
        response.setHeader("X-XSS-Protection", "1; mode=block");
        response.setHeader("Strict-Transport-Security",
                          "max-age=31536000; includeSubDomains");
        response.setHeader("Content-Security-Policy",
                          "default-src 'self'");
        response.setHeader("Referrer-Policy",
                          "strict-origin-when-cross-origin");

        return true;
    }
}
```

## 📦 Security Report Template

```markdown
# Security Review Report

**Project:** Arda Platform
**Date:** 2026-04-25
**Reviewer:** Security Team

## Executive Summary

- Total issues found: X
- Critical: Y
- High: Z
- Medium: A
- Low: B

## Critical Issues

### Issue: [Title]
**Severity:** Critical
**CVSS Score:** X.X
**Location:** file:line

**Description:**
[Explain the vulnerability]

**Impact:**
[What could happen]

**Recommendation:**
[How to fix]

**Proof of Concept:**
```go
// Code snippet demonstrating vulnerability
```

## High Severity Issues

[Similar format for high issues]

## Medium Severity Issues

[Similar format for medium issues]

## Low Severity Issues

[Similar format for low issues]

## Recommendations

1. [Recommendation 1]
2. [Recommendation 2]
3. [Recommendation 3]

## References

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [OWASP ASVS](https://owasp.org/www-project-application-security-verification-standard/)
- [CWE Top 25](https://cwe.mitre.org/top25/)
```

## 🎯 Usage Examples

```
/security-review "Review Go service"

Usage:
/security-review "Review security cho iam-service"

Sẽ:
1. Check OWASP Top 10
2. Review auth/authorization
3. Generate security report
```

```
/security-review "Review Java service"

Usage:
/security-review "Review security cho accounting-service"

Sẽ:
1. Check for vulnerabilities
2. Review Spring Security
3. Suggest improvements
```

```
/security-review "Review Angular app"

Usage:
/security-review "Review security cho arda-mfe shell app"

Sẽ:
1. Check XSS vulnerabilities
2. Review auth implementation
3. Review API calls security
```

## 📦 Best Practices

### Security Development

- Follow security by design principles
- Use secure defaults
- Implement defense in depth
- Regular security reviews
- Keep dependencies updated

### Input Validation

- Validate all inputs
- Use whitelist approach
- Sanitize data
- Type checking
- Length validation

### Error Handling

- Don't leak sensitive info
- Use generic error messages
- Log security events
- Handle errors gracefully
- Don't expose stack traces

---

*Last Updated: 2026-04-25*
