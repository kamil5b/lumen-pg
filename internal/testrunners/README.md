# Lumen-PG Test Runners

This directory contains test runners following the test-first, layered architecture pattern.

## Overview

Test runners are reusable test suites that verify implementations conform to interfaces. Each runner:
- Takes a constructor function that creates an implementation
- Runs a comprehensive test suite against that implementation
- Uses real databases (testcontainers) for repositories
- Uses mocks (gomock) for services and handlers

## Directory Structure

```
internal/
├─ domain/              # Pure business types (no dependencies)
├─ interfaces/          # Interface contracts for all layers
├─ testrunners/         # Test runners (this directory)
└─ implementations/     # Concrete implementations (to be created)
   ├─ mocks/            # Generated mocks
   ├─ repository/       # PostgreSQL repository implementations
   ├─ service/          # Service layer implementations
   └─ handler/          # HTTP handler implementations
```

## Generating Mocks

Before running service and handler tests, generate mocks:

```bash
# Install mockgen if not already installed
go install go.uber.org/mock/mockgen@latest

# Generate repository mocks
mockgen -source=internal/interfaces/repository.go \
    -destination=internal/implementations/mocks/mock_repository.go \
    -package=mocks

# Generate service mocks
mockgen -source=internal/interfaces/service.go \
    -destination=internal/implementations/mocks/mock_service.go \
    -package=mocks

# Generate handler mocks (if needed)
mockgen -source=internal/interfaces/handler.go \
    -destination=internal/implementations/mocks/mock_handler.go \
    -package=mocks
```

## Test Runner Pattern

Each runner follows this signature:
```go
func Runner(t *testing.T, constructor func(dependencies) Interface)
```

### Repository Runners

Repository runners use **testcontainers** to test against real PostgreSQL:

```go
func TestMyConnectionRepository(t *testing.T) {
    constructor := func(connStr string) interfaces.ConnectionRepository {
        return implementations.NewMyConnectionRepository(connStr)
    }
    testrunners.ConnectionRepositoryRunner(t, constructor)
}
```

Available repository runners:
- `ConnectionRepositoryRunner` - Tests connection management
- `MetadataRepositoryRunner` - Tests metadata loading
- `QueryRepositoryRunner` - Tests query execution
- `TransactionRepositoryRunner` - Tests transaction handling

### Service Runners

Service runners use **gomock** to mock dependencies:

```go
func TestMyAuthService(t *testing.T) {
    constructor := func(
        connRepo interfaces.ConnectionRepository, 
        metadataRepo interfaces.MetadataRepository,
    ) interfaces.AuthService {
        return implementations.NewMyAuthService(connRepo, metadataRepo)
    }
    testrunners.AuthServiceRunner(t, constructor)
}
```

Available service runners:
- `AuthServiceRunner` - Tests authentication and session management
- `MetadataServiceRunner` - Tests metadata operations
- `QueryServiceRunner` - Tests query execution and validation
- `DataExplorerServiceRunner` - Tests main view data operations
- `TransactionServiceRunner` - Tests transaction management

### Handler Runners

Handler runners use **gomock** to mock services:

```go
func TestMyAuthHandler(t *testing.T) {
    constructor := func(authService interfaces.AuthService) interfaces.AuthHandler {
        return implementations.NewMyAuthHandler(authService)
    }
    testrunners.AuthHandlerRunner(t, constructor)
}
```

Available handler runners:
- `AuthHandlerRunner` - Tests authentication endpoints
- `MainViewHandlerRunner` - Tests main view endpoints
- `QueryEditorHandlerRunner` - Tests query editor endpoints
- `ERDViewerHandlerRunner` - Tests ERD viewer endpoints
- `TransactionHandlerRunner` - Tests transaction endpoints

## Implementation Workflow

Follow this TDD workflow:

1. **Write failing tests** - Tests are already written in test runners
2. **Create implementation** - Implement the interface
3. **Bind implementation to test** - Create test file that calls the runner
4. **Run tests** - Watch them fail
5. **Fix implementation** - Iterate until tests pass

### Example: Implementing a Repository

```go
// Step 1: Create implementation
// internal/implementations/repository/postgres_connection_repository.go
package repository

import (
    "context"
    "database/sql"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type PostgresConnectionRepository struct {
    superAdminConnStr string
}

func NewPostgresConnectionRepository(connStr string) interfaces.ConnectionRepository {
    return &PostgresConnectionRepository{
        superAdminConnStr: connStr,
    }
}

func (r *PostgresConnectionRepository) ValidateConnectionString(connStr string) error {
    // Implementation here
    return nil
}

// ... implement other methods

// Step 2: Bind to test runner
// internal/implementations/repository/postgres_connection_repository_test.go
package repository_test

import (
    "testing"
    "github.com/kamil5b/lumen-pg/internal/implementations/repository"
    "github.com/kamil5b/lumen-pg/internal/testrunners"
)

func TestPostgresConnectionRepository(t *testing.T) {
    constructor := func(connStr string) interfaces.ConnectionRepository {
        return repository.NewPostgresConnectionRepository(connStr)
    }
    testrunners.ConnectionRepositoryRunner(t, constructor)
}
```

## Running Tests

```bash
# Run all tests
go test ./... -v

# Run only repository tests
go test ./internal/testrunners -v -run Repository

# Run only service tests  
go test ./internal/testrunners -v -run Service

# Run only handler tests
go test ./internal/testrunners -v -run Handler

# Run specific test runner
go test ./internal/testrunners -v -run ConnectionRepository

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Data Setup

Repository tests automatically create test schemas and data using testcontainers. Examples:

### Users and Posts Schema
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200),
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Key Testing Principles

1. **Test-First**: Write tests before implementation
2. **Interface-Driven**: All layers depend on interfaces, not implementations
3. **Isolation**: Repository tests use real DB, service/handler tests use mocks
4. **Reusability**: Same test runner works for multiple implementations
5. **Comprehensive**: Runners cover happy paths, error cases, edge cases

## Benefits

✅ Standardized testing across all layers  
✅ Easy to swap implementations  
✅ Forces interface-driven design  
✅ Catches integration issues early (with testcontainers)  
✅ Fast unit tests (with mocks)  
✅ Clear separation of concerns  

## Next Steps

1. Generate mocks: `mockgen` commands above
2. Implement repository layer (test against real DB)
3. Implement service layer (test with mocked repos)
4. Implement handler layer (test with mocked services)
5. Create E2E tests (optional, for full user flows)

## Notes

- Repository tests require Docker (for testcontainers)
- Service and handler tests require generated mocks
- All tests marked with `t.Skip()` until mocks are generated
- Follow the existing pattern when adding new test runners
