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
    "testing"
    "lumen-pg/internal/testrunners/e2e"
)

func TestE2ERoutes(t *testing.T) {
    // Setup your complete router with all middleware
    router := setupRouter()
    
    // Run all E2E tests
    e2e.RunAllE2ETests(t, router)
}
```

### Running Specific Story Tests

```go
func TestAuthenticationE2E(t *testing.T) {
    router := setupRouter()
    e2e.RunAuthenticationE2ETests(t, router)
}

func TestStory5E2E(t *testing.T) {
    router := setupRouter()
    e2e.RunStory(t, router, 5)
}
```

### Running Individual Test Cases

```go
func TestLoginFlowE2E(t *testing.T) {
    router := setupRouter()
    e2e.Story2AuthE2ERunner(t, router)
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

E2E tests require a fully configured router with all middleware. Example:

```go
func setupRouter() http.Handler {
    mux := http.NewServeMux()
    
    // Create all handlers
    loginHandler := NewLoginHandler(authUC, setupUC, rbacUC)
    mainViewHandler := NewMainViewHandler(dataUC, metadataUC)
    // ... other handlers
    
    // Create middleware
    authMiddleware := NewAuthenticationMiddleware(sessionRepo)
    securityMiddleware := NewSecurityMiddleware()
    loggingMiddleware := NewLoggingMiddleware()
    // ... other middleware
    
    // Register routes with middleware chain
    mux.Handle("/login", loginHandler)
    mux.Handle("/logout", loginHandler)
    mux.Handle("/main", authMiddleware.RequireAuth(mainViewHandler))
    mux.Handle("/query-editor", authMiddleware.RequireAuth(queryHandler))
    // ... other routes
    
    // Wrap with global middleware
    handler := loggingMiddleware.LogRequest(
        securityMiddleware.SetSecurityHeaders(
            mux,
        ),
    )
    
    return handler
}
```

## Test Data

E2E tests expect certain test data to exist:
- Test users: `testuser`, `testuser1`, `testuser2`, `noaccessuser`
- Test databases: `testdb`, `testdb1`, `testdb2`
- Test tables: `users`, `posts`, `comments`, `products`, `orders`

See TEST_PLAN.md [L831-899] for complete test data setup.

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
**Solution**: Ensure proper test isolation and cleanup between tests

### Issue: Cookies not being sent in requests
**Solution**: Remember to add cookies from login response to subsequent requests

### Issue: Tests fail with "no rows in result set"
**Solution**: Ensure test database is properly seeded before running tests

### Issue: Transaction tests interfere with each other
**Solution**: Use separate test users or databases for concurrent transaction tests

## Integration with CI/CD

E2E tests should run after:
1. Unit tests pass
2. Integration tests pass
3. Test database is seeded

Example CI pipeline:
```yaml
- name: Run Unit Tests
  run: go test -short ./...
  
- name: Run Integration Tests
  run: go test -tags=integration ./...
  
- name: Setup Test Database
  run: ./scripts/setup-test-db.sh
  
- name: Run E2E Tests
  run: go test -v ./internal/testrunners/e2e/...
```

## References

- TEST_PLAN.md: Complete test specifications
- TDD Workflow: Phase 6: E2E Tests [L233-281]
- Handler Test Runners: `internal/testrunners/handler/`
- Integration Tests: Test with real PostgreSQL

## Notes

- E2E tests are slower than unit tests - run them selectively during development
- Use `-short` flag to skip E2E tests during rapid iteration
- E2E tests require a running PostgreSQL instance (can use testcontainers)
- Some tests may require specific PostgreSQL versions or extensions
- Cookie security tests may behave differently in test vs production environments