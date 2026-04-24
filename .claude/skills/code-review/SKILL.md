---
name: code-review
description: Hỗ trợ review code và suggest improvements
disable-model-invocation: false

---
# Code Review Skill

Mục đích: Hỗ trợ review code và suggest improvements cho Go, Java, và TypeScript code trong dự án Arda.

## 🎯 Phạm vi

- Review code patterns
- Suggest best practices
- Identify code smells
- Suggest performance improvements
- Suggest security improvements
- Generate review report

## 📦 Code Review Checklist

### Go Code Review

```markdown
## Go Code Review Checklist

### Code Style
- [ ] Uses proper naming conventions (PascalCase for exports, camelCase for internal)
- [ ] Follows effective Go guidelines
- [ ] Uses proper package structure
- [ ] Has proper file organization

### Error Handling
- [ ] Errors are checked and handled
- [ ] Uses fmt.Errorf for error wrapping
- [ ] Errors have context (who, what, where)
- [ ] No panic in production code

### Concurrency
- [ ] Proper mutex usage (if needed)
- [ ] No race conditions
- [ ] Proper channel usage
- [ ] Context cancellation handled

### Performance
- [ ] No unnecessary allocations
- [ ] Proper string building
- [ ] Efficient database queries
- [ ] Appropriate use of goroutines

### Security
- [ ] No SQL injection
- [ ] Input validation
- [ ] Proper authentication/authorization
- [ ] No hardcoded secrets

### Testing
- [ ] Has unit tests
- [ ] Tests cover happy paths and error cases
- [ ] No flaky tests
- [ ] Test naming is clear

### Documentation
- [ ] Package-level documentation
- [ ] Exported functions have comments
- [ ] Complex code has explanations
- [ ] Examples in comments (if needed)
```

### Java Code Review

```markdown
## Java Code Review Checklist

### Code Style
- [ ] Follows Google Java Style Guide
- [ ] Proper use of generics
- [ ] Proper exception handling
- [ ] Consistent formatting

### Design Patterns
- [ ] Appropriate use of design patterns
- [ ] SOLID principles followed
- [ ] Proper separation of concerns
- [ ] DRY (Don't Repeat Yourself)

### Spring Boot
- [ ] Proper use of @Service, @Repository, @Controller
- [ ] Transaction boundaries correct
- [ ] Proper dependency injection
- [ ] No @Autowired on fields

### Performance
- [ ] Efficient database queries
- [ ] Proper use of caching
- [ ] No N+1 queries
- [ ] Proper lazy loading

### Security
- [ ] Input validation
- [ ] Proper authentication/authorization
- [ ] No SQL injection
- [ ] Secure serialization

### Testing
- [ ] Has unit tests
- [ ] Tests use JUnit 5
- [ ] Proper mocking (Mockito)
- [ ] Tests are deterministic

### Documentation
- [ ] Javadoc for public APIs
- [ ] Complex logic explained
- [ ] README with examples
- [ ] API documentation
```

### TypeScript/Angular Code Review

```markdown
## TypeScript/Angular Code Review Checklist

### Code Style
- [ ] Follows Angular style guide
- [ ] Proper use of TypeScript
- [ ] Consistent formatting (Prettier)
- [ ] Proper file naming

### Angular Specific
- [ ] Uses Angular conventions
- [ ] Proper use of RxJS
- [ ] No memory leaks (unsubscribe)
- [ ] Proper change detection strategy

### Component Design
- [ ] Single responsibility
- [ ] Proper input/output
- [ ] No logic in templates
- [ ] Proper use of services

### Performance
- [ ] Uses OnPush change detection
- [ ] Uses trackBy for lists
- [ ] No unnecessary re-renders
- [ ] Proper lazy loading

### Security
- [ ] Input validation
- [ ] XSS prevention
- [ ] Proper authentication
- [ ] No template injection

### Testing
- [ ] Has unit tests
- [ ] Tests use TestBed properly
- [ ] No test-specific code in production
- [ ] Tests are deterministic

### Documentation
- [ ] README with examples
- [ ] API documentation
- [ ] Complex code explained
- [ ] Comments where needed
```

## 📦 Common Code Smells

### Code Smells Detection

#### Go Code Smells

```go
// Bad: Too long function
func ProcessRequest(req *Request) (*Response, error) {
    // 100+ lines of code
}

// Good: Break down into smaller functions
func ProcessRequest(req *Request) (*Response, error) {
    if err := validateRequest(req); err != nil {
        return nil, err
    }

    data, err := fetchRequestData(req)
    if err != nil {
        return nil, err
    }

    return processData(data)
}

// Bad: Deep nesting
func processData(data []Data) error {
    for _, item := range data {
        if item.Type == "A" {
            for _, sub := range item.SubItems {
                if sub.Status == "active" {
                    // process
                }
            }
        }
    }
}

// Good: Early returns and helper functions
func processData(data []Data) error {
    for _, item := range data {
        if item.Type == "A" {
            if err := processTypeA(item); err != nil {
                return err
            }
        }
    }
    return nil
}

// Bad: Inconsistent error handling
func CreateUser(req *CreateUserRequest) (*User, error) {
    user := &User{Name: req.Name}
    // Save user - no error check
    repo.Save(user)
    return user, nil
}

// Good: Proper error handling
func CreateUser(req *CreateUserRequest) (*User, error) {
    if err := validateUserRequest(req); err != nil {
        return nil, fmt.Errorf("invalid request: %w", err)
    }

    user := &User{Name: req.Name}
    if err := repo.Save(user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }

    return user, nil
}
```

#### Java Code Smells

```java
// Bad: Too many parameters
public void createUser(String name, String email, String phone,
                     String address, String city, String state,
                     String zip, String country, ...) {
    // Create user
}

// Good: Use parameter object
public class CreateUserRequest {
    private String name;
    private String email;
    private String phone;
    // Getters and setters
}

public void createUser(CreateUserRequest request) {
    // Create user
}

// Bad: God class
class UserService {
    public void createUser() { }
    public void updateUser() { }
    public void deleteUser() { }
    public void sendEmail() { }
    public void sendSMS() { }
    public void generateReport() { }
    // 50+ more methods...
}

// Good: Single responsibility
class UserService {
    public void createUser() { }
    public void updateUser() { }
    public void deleteUser() { }
}

class NotificationService {
    public void sendEmail() { }
    public void sendSMS() { }
}

// Bad: Magic numbers
public class DiscountCalculator {
    public double calculateDiscount(double price) {
        return price * 0.1; // What is 0.1?
    }
}

// Good: Use constants
public class DiscountCalculator {
    private static final double STANDARD_DISCOUNT = 0.1;
    private static final double PREMIUM_DISCOUNT = 0.15;

    public double calculateDiscount(double price, CustomerType type) {
        return price * getDiscountRate(type);
    }
}
```

#### TypeScript/Angular Code Smells

```typescript
// Bad: Logic in template
@Component({
  selector: 'app-user',
  template: `
    <div *ngIf="user.age > 18 && user.age < 65 && user.hasLicense">
      {{ user.name }}
    </div>
  `
})
export class UserComponent {
  user: any;
}

// Good: Logic in component
@Component({
  selector: 'app-user',
  template: `
    <div *ngIf="canDrive">
      {{ user.name }}
    </div>
  `
})
export class UserComponent {
  user: any;

  get canDrive(): boolean {
    return this.isAdult && this.hasLicense;
  }

  get isAdult(): boolean {
    return this.user.age >= 18 && this.user.age < 65;
  }

  get hasLicense(): boolean {
    return this.user.hasLicense;
  }
}

// Bad: No unsubscribe
@Component({
  selector: 'app-user',
  template: `<div>{{ data }}</div>`
})
export class UserComponent implements OnInit {
  data: any;

  ngOnInit() {
    this.userService.getData().subscribe(data => {
      this.data = data;
    });
  }
}

// Good: Proper cleanup
@Component({
  selector: 'app-user',
  template: `<div>{{ data }}</div>`
})
export class UserComponent implements OnInit, OnDestroy {
  data: any;
  private destroy$ = new Subject<void>();

  ngOnInit() {
    this.userService.getData()
      .pipe(takeUntil(this.destroy$))
      .subscribe(data => {
        this.data = data;
      });
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }
}

// Bad: Using any
@Component({
  selector: 'app-user',
  template: `<div>{{ user.name }}</div>`
})
export class UserComponent {
  user: any; // What is user?

  ngOnInit() {
    this.userService.getUser().subscribe(data => {
      this.user = data;
    });
  }
}

// Good: Proper typing
interface User {
  id: string;
  name: string;
  email: string;
}

@Component({
  selector: 'app-user',
  template: `<div>{{ user.name }}</div>`
})
export class UserComponent {
  user: User;

  ngOnInit() {
    this.userService.getUser().subscribe(data => {
      this.user = data;
    });
  }
}
```

## 📦 Performance Suggestions

### Go Performance

```go
// Bad: String concatenation in loop
func concatStrings(items []string) string {
    var result string
    for _, item := range items {
        result += item + "," // Inefficient
    }
    return result
}

// Good: Use strings.Builder
func concatStrings(items []string) string {
    var sb strings.Builder
    for _, item := range items {
        sb.WriteString(item)
        sb.WriteString(",")
    }
    return sb.String()
}

// Bad: No context for long operations
func processData() error {
    time.Sleep(10 * time.Second)
    return nil
}

// Good: Use context for cancellation
func processData(ctx context.Context) error {
    select {
    case <-time.After(10 * time.Second):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### Java Performance

```java
// Bad: String concatenation in loop
public String concatenateStrings(List<String> items) {
    String result = "";
    for (String item : items) {
        result += item + ","; // Inefficient
    }
    return result;
}

// Good: Use StringBuilder
public String concatenateStrings(List<String> items) {
    StringBuilder sb = new StringBuilder();
    for (String item : items) {
        sb.append(item).append(",");
    }
    return sb.toString();
}

// Bad: No lazy loading
@Entity
public class User {
    @OneToMany(fetch = FetchType.EAGER)
    private List<Order> orders;
}

// Good: Use lazy loading
@Entity
public class User {
    @OneToMany(fetch = FetchType.LAZY)
    private List<Order> orders;
}
```

## 📦 Security Suggestions

### Common Security Issues

```go
// Bad: SQL injection possible
query := fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID)
rows, err := db.Query(query)

// Good: Use parameterized queries
query := "SELECT * FROM users WHERE id = $1"
rows, err := db.Query(query, userID)

// Bad: No input validation
func CreateUser(name, email string) (*User, error) {
    user := &User{Name: name, Email: email}
    return repo.Save(user)
}

// Good: Validate input
func CreateUser(name, email string) (*User, error) {
    if !isValidEmail(email) {
        return nil, fmt.Errorf("invalid email format")
    }

    if len(name) == 0 || len(name) > 255 {
        return nil, fmt.Errorf("invalid name length")
    }

    user := &User{Name: name, Email: email}
    return repo.Save(user)
}
```

```java
// Bad: SQL injection
@Query("SELECT * FROM User WHERE email = '" + email + "'")
List<User> findByEmail(String email);

// Good: Parameterized query
@Query("SELECT * FROM User WHERE email = :email")
List<User> findByEmail(@Param("email") String email);

// Bad: No input validation
@PostMapping("/users")
public User createUser(@RequestBody User user) {
    return userRepository.save(user);
}

// Good: Validate input
@PostMapping("/users")
public User createUser(@Valid @RequestBody User user) {
    return userRepository.save(user);
}
```

## 🎯 Usage Examples

```
/code-review "Review Go service"

Usage:
/code-review "Review file internal/service/user_service.go"

Sẽ:
1. Phân tích code style
2. Check error handling
3. Suggest improvements
```

```
/code-review "Review Java service"

Usage:
/code-review "Review file services/src/main/java/.../JournalService.java"

Sẽ:
1. Review design patterns
2. Check Spring Boot usage
3. Suggest performance improvements
```

```
/code-review "Review Angular component"

Usage:
/code-review "Review component apps/accounting/src/app/journals/journal-list.component.ts"

Sẽ:
1. Review Angular conventions
2. Check RxJS usage
3. Suggest improvements
```

## 📦 Best Practices

### Review Guidelines

- Be constructive and helpful
- Provide specific examples
- Explain why changes are needed
- Suggest alternatives when possible
- Respect the author's choices

### Feedback Format

```
## Issue: [Brief description]

**Severity:** [Critical/Major/Minor]

**Location:** [File:line]

**Problem:** [Explain the issue]

**Suggestion:** [Proposed fix]

**Example:**
```go
// Current code
func foo() {
    // code
}

// Suggested code
func foo() {
    // improved code
}
```
```

---

*Last Updated: 2026-04-25*
