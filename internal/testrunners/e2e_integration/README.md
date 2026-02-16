# E2E Test Runners - Lumen-PG

This directory contains **Route-level End-to-End (E2E) Test Runners** that test complete HTTP request/response flows through the entire application stack.

## Overview

E2E test runners verify the **complete integration** of all application layers:
- HTTP routing and request handling
- Full middleware chain (authentication, authorization, logging, security, etc.)
- Handler implementations
- Use case orchestration
- Repository implementations
- Database interactions
- Session and cookie management
- Multi-user isolation

**Key Difference from Other Test Types:**
- **Handler Test Runners** (`internal/testrunners/handler/`): Test handlers in isolation with mocked use cases
- **Integration Tests**: Test repository implementations with real databases
- **E2E Test Runners** (this directory): Test complete routes with real HTTP stack and full middleware

## Test Structure

Each story from TEST_PLAN.md has a dedicated E2E test runner:

```
internal/testrunners/e2e/
├── README.md                           # This file
├── e2e_runner.go                       # Main orchestrator
├── story2_auth_e2e_runner.go          # Story 2: Authentication & Identity
├── story3_erd_e2e_runner.go           # Story 3: ERD Viewer
├── story4_query_editor_e2e_runner.go  # Story 4: Manual Query Editor
├── story5_main_view_e2e_runner.go     # Story 5: Main View & Data Interaction
├── story6_isolation_e2e_runner.go     # Story 6: Isolation
└── story7_security_e2e_runner.go      # Story 7: Security & Best Practices
```

## Story Mapping

### Story 2: Authentication & Identity (E2E-S2-01 to E2E-S2-06)
Tests complete login/logout flows with session management:
- Login form validation and connection probe
- Session cookie creation and validation
- Protected route access control
- Data explorer population after login
- Multi-step authentication flows

### Story 3: ERD Viewer (E2E-S3-01 to E2E-S3-04)
Tests ERD visualization routes:
- ERD viewer page access with authentication
- ERD zoom controls (in/out/reset)
- ERD pan/drag functionality
- Table click navigation from ERD to main view

### Story 4: Manual Query Editor (E2E-S4-01 to E2E-S4-06)
Tests query execution routes:
- Query editor page access
- Single and multiple query execution
- Query error display and handling
- Offset pagination with actual size display
- SQL syntax highlighting support

### Story 5: Main View & Data Interaction (E2E-S5-01 to E2E-S5-15a)
Tests main data view and transaction routes:
- Table data loading with cursor pagination
- WHERE clause filtering and column sorting
- Transaction management (start, edit, commit, rollback)
- Transaction timer and edit buffer display
- Foreign key and primary key navigation

### Story 6: Isolation (E2E-S6-01 to E2E-S6-03)
Tests multi-user isolation:
- Simultaneous users with different permissions
- Concurrent transaction isolation
- Session cookie isolation between users
- Independent query execution

### Story 7: Security & Best Practices (E2E-S7-01 to E2E-S7-06)
Tests security features:
- SQL injection prevention (WHERE bar, query editor)
- Cookie tampering detection
- Session timeout enforcement
- HTTPS-only cookies (when enabled)
- HttpOnly and SameSite cookie attributes
- Password encryption in cookies

## Usage

### Running All E2E Tests

```go
package main_test

import (
    "database/sql"
    "net/http"
    "testing"
    "lumen-pg/internal/testrunners/e2e_integration"
)

func TestE2ERoutes(t *testing.T) {
    // E2E runner automatically sets up testcontainer and database
    // You just provide a constructor function that builds your router
    e2e_integration.RunAllE2ETests(t, func(db *sql.DB) http.Handler {
        // Build your complete router with all middleware using the provided DB
        return setupRouter(db)
    })
}
```

### Running Specific Story Tests

```go
func TestAuthenticationE2E(t *testing.T) {
    e2e_integration.RunAuthenticationE2ETests(t, func(db *sql.DB) http.Handler {
        return setupRouter(db)
    })
}

func TestStory5E2E(t *testing.T) {
    e2e_integration.RunStory(t, func(db *sql.DB) http.Handler {
        return setupRouter(db)
    }, 5)
}
```

### Running Individual Test Cases (with pre-configured router)

```go
func TestLoginFlowE2E(t *testing.T) {
    // For individual test cases, you can use the old pattern
    // if you want more control over database setup
    router := setupRouter(db)
    e2e_integration.Story2AuthE2ERunner(t, router)
}
```

### Using Pre-Configured Router (Backward Compatibility)

```go
func TestE2EWithExistingRouter(t *testing.T) {
    // If you need to use a pre-configured router
    // (e.g., for testing against a specific database setup)
    router := setupRouterWithCustomDB()
    e2e_integration.RunAllE2ETestsWithRouter(t, router)
}
```

## Test Patterns

### 1. Authentication Flow Pattern

```go
// Login and get session cookies
formData := url.Values{}
formData.Set("username", "testuser")
formData.Set("password", "testpass")

req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
rec := httptest.NewRecorder()
router.ServeHTTP(rec, req)

cookies := rec.Result().Cookies()

// Use cookies in subsequent requests
req = httptest.NewRequest(http.MethodGet, "/main", nil)
for _, cookie := range cookies {
    req.AddCookie(cookie)
}
router.ServeHTTP(httptest.NewRecorder(), req)
```

### 2. Multi-Request Flow Pattern

```go
// Test complete user journey
// 1. Login
// 2. Access data explorer
// 3. Select table
// 4. Start transaction
// 5. Edit data
// 6. Commit transaction
// 7. Verify results
```

### 3. Concurrent User Pattern

```go
var wg sync.WaitGroup
wg.Add(2)

// User 1 executes action
go func() {
    defer wg.Done()
    // User 1 operations with user1Cookies
}()

// User 2 executes action
go func() {
    defer wg.Done()
    // User 2 operations with user2Cookies
}()

wg.Wait()

// Verify isolation
```

## Test Assertions

### HTTP Status Codes
```go
assert.Equal(t, http.StatusOK, rec.Code)
assert.True(t, rec.Code == http.StatusFound || rec.Code == http.StatusUnauthorized)
```

### Response Body
```go
body := rec.Body.String()
assert.Contains(t, body, "expected content")
assert.NotContains(t, body, "sensitive data")
```

### Cookie Validation
```go
cookies := rec.Result().Cookies()
assert.NotEmpty(t, cookies)

for _, cookie := range cookies {
    if cookie.Name == "session_cookie" {
        assert.True(t, cookie.HttpOnly)
        assert.True(t, cookie.Secure) // If HTTPS
        assert.True(t, cookie.SameSite == http.SameSiteStrictMode)
    }
}
```

### JSON Response
```go
var response map[string]interface{}
err := json.Unmarshal(rec.Body.Bytes(), &response)
require.NoError(t, err)

assert.NotNil(t, response["data"])
assert.Equal(t, "success", response["status"])
```

## Router Setup

E2E tests use a **RouterConstructor** pattern. You provide a function that builds your complete router given a database connection. The E2E runner handles testcontainer setup and teardown automatically.

### RouterConstructor Pattern

```go
func setupRouter(db *sql.DB) http.Handler {
    // Create repositories with the provided database
    dbRepo := repository.NewDatabaseRepository(db)
    sessionRepo := repository.NewSessionRepository(db)
    metadataRepo := repository.NewMetadataRepository(db)
    
    // Create use cases
    authUC := usecase.NewAuthenticationUseCase(dbRepo, sessionRepo)
    dataUC := usecase.NewDataUseCase(dbRepo, metadataRepo)
    
    // Create handlers
    loginHandler := handler.NewLoginHandler(authUC)
    mainViewHandler := handler.NewMainViewHandler(dataUC)
    queryHandler := handler.NewQueryEditorHandler(dataUC)
    
    // Create middleware
    authMiddleware := middleware.NewAuthenticationMiddleware(sessionRepo)
    securityMiddleware := middleware.NewSecurityMiddleware()
    loggingMiddleware := middleware.NewLoggingMiddleware()
    
    // Setup router
    mux := http.NewServeMux()
    
    // Public routes
    mux.Handle("/login", loginHandler)
    mux.Handle("/logout", loginHandler)
    
    // Protected routes
    mux.Handle("/main", authMiddleware.RequireAuth(mainViewHandler))
    mux.Handle("/query-editor", authMiddleware.RequireAuth(queryHandler))
    mux.Handle("/erd-viewer", authMiddleware.RequireAuth(erdHandler))
    
    // API routes
    mux.Handle("/api/data-explorer", authMiddleware.RequireAuth(dataExplorerHandler))
    mux.Handle("/api/execute-query", authMiddleware.RequireAuth(queryHandler))
    mux.Handle("/api/transaction/start", authMiddleware.RequireAuth(transactionHandler))
    
    // Wrap with global middleware
    handler := loggingMiddleware.LogRequest(
        securityMiddleware.SetSecurityHeaders(mux),
    )
    
    return handler
}
```

### Benefits of RouterConstructor Pattern

1. **Self-Contained**: E2E runner manages testcontainer lifecycle
2. **Consistent**: All E2E tests use the same database setup
3. **Easy to Use**: No manual database setup required
4. **Clean**: Automatic resource cleanup
5. **Follows Repository Pattern**: Similar to repository test runners

## Test Data

E2E tests **automatically seed test data** via the `seedE2ETestData` function. The data includes:

### Database Schema
- **users**: Test users with id, name, email, created_at
- **posts**: Blog posts with foreign key to users
- **comments**: Comments with foreign keys to posts and users
- **products**: Product catalog with pricing and stock
- **orders**: Order history with foreign keys to users and products

### Seeded Users
- Alice (id=1, alice@example.com)
- Bob (id=2, bob@example.com)
- Charlie (id=3, charlie@example.com)
- David (id=4, david@example.com)

### Seeded Data
- 4 posts from various users
- 4 comments on different posts
- 4 products in catalog
- 4 orders from different users

### Test Credentials
E2E tests expect your authentication to work with:
- **testuser** / **testpass** (valid user)
- **testuser1** / **testpass1** (for multi-user tests)
- **testuser2** / **testpass2** (for multi-user tests)
- **noaccessuser** / **testpass** (user with no database permissions)

**Note**: The authentication logic must be implemented in your router/handlers.
The E2E runner only provides the database schema and data.

## Running Tests

### Run all E2E tests
```bash
go test -v ./internal/testrunners/e2e/...
```

### Run specific story
```bash
go test -v ./internal/testrunners/e2e/ -run "Story2"
go test -v ./internal/testrunners/e2e/ -run "Story5"
```

### Run specific test case
```bash
go test -v ./internal/testrunners/e2e/ -run "E2E-S2-01"
go test -v ./internal/testrunners/e2e/ -run "LoginFlowWithConnectionProbe"
```

### Run with coverage
```bash
go test -v -cover ./internal/testrunners/e2e/...
```

## Best Practices

1. **Test Complete Flows**: Test entire user journeys, not just single endpoints
2. **Use Real HTTP Stack**: Use `httptest.NewRecorder()` and `router.ServeHTTP()`
3. **Cookie Management**: Properly handle session cookies across requests
4. **Isolation**: Each test should be independent and not affect others
5. **Cleanup**: Clean up test data after tests complete
6. **Concurrent Testing**: Test multi-user scenarios with goroutines
7. **Security Validation**: Always verify security features (HttpOnly, Secure, etc.)
8. **Error Scenarios**: Test both success and failure paths

## Common Issues

### Issue: Tests pass individually but fail when run together
**Solution**: Each E2E runner function now creates its own testcontainer, ensuring isolation

### Issue: Cookies not being sent in requests
**Solution**: Remember to add cookies from login response to subsequent requests

### Issue: Tests fail with "no rows in result set"
**Solution**: E2E runner automatically seeds data via `seedE2ETestData`. If tests still fail, check your authentication logic

### Issue: Transaction tests interfere with each other
**Solution**: Each test suite gets a fresh testcontainer, preventing interference

### Issue: "connection refused" or database errors
**Solution**: Ensure Docker is running (testcontainers requires Docker)

## Integration with CI/CD

E2E tests now **automatically handle database setup** via testcontainers. No manual setup required!

Example CI pipeline:
```yaml
- name: Run Unit Tests
  run: go test -short ./...
  
- name: Run Integration Tests
  run: go test -tags=integration ./...
  
- name: Run E2E Tests
  run: go test -v ./internal/testrunners/e2e_integration/...
  
# testcontainers automatically starts/stops PostgreSQL during tests
```

### CI Requirements
- Docker daemon must be available
- Sufficient permissions to create containers
- Network access to pull postgres:15 image

## References

- TEST_PLAN.md: Complete test specifications
- TDD Workflow: Phase 6: E2E Tests [L233-281]
- Handler Test Runners: `internal/testrunners/handler/`
- Integration Tests: Test with real PostgreSQL

## Notes

- E2E tests are slower than unit tests (due to testcontainer startup) - run them selectively during development
- Use `-short` flag to skip E2E tests during rapid iteration
- E2E tests automatically manage PostgreSQL via testcontainers (Docker required)
- Tests use PostgreSQL 15 by default
- Each test suite creates a fresh testcontainer for complete isolation
- Cookie security tests may behave differently in test vs production environments
- The `RouterConstructor` pattern is inspired by repository test runners for consistency