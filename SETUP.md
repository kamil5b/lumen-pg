# Setup Guide - Lumen-PG Implementation

This guide explains how to proceed with implementing the Lumen-PG application after completing Steps 1-3 (Domain, Interfaces, and Test Runners).

## Current Status ✅

**Steps 1-3 Complete:**
- ✅ Domain layer defined (6 files)
- ✅ Interface layer defined (3 files)  
- ✅ Test runners written (3 test files)
- ✅ Dependencies installed (go.mod updated)

## What's Been Done

### 1. Domain Layer (`internal/domain/`)
Pure business types with no external dependencies:
- `user.go` - User and LoginInput types
- `session.go` - Session and SessionCookies types
- `metadata.go` - Database metadata types (tables, columns, relations, permissions)
- `query.go` - Query execution types (requests and results)
- `transaction.go` - Transaction state and operations
- `erd.go` - Entity-relationship diagram types

### 2. Interface Layer (`internal/interfaces/`)
Contracts for all layers:
- `repository.go` - 4 repository interfaces (Connection, Metadata, Query, Transaction)
- `service.go` - 5 service interfaces (Auth, Metadata, Query, DataExplorer, Transaction)
- `handler.go` - 5 handler interfaces (Auth, MainView, QueryEditor, ERDViewer, Transaction)

### 3. Test Runners (`internal/testrunners/`)
Test-first specifications:
- `repository_runner_test.go` - 4 repository test runners (use testcontainers)
- `service_runner_test.go` - 5 service test runners (use gomock)
- `handler_runner_test.go` - 5 handler test runners (use gomock)
- `README.md` - Comprehensive test runner documentation

## Next Steps 🚀

### Step 4: Generate Mocks

Before implementing services and handlers, generate mocks for testing:

```bash
# Install mockgen if not already installed
go install go.uber.org/mock/mockgen@latest

# Create mocks directory
mkdir -p internal/implementations/mocks

# Generate repository mocks (for service layer tests)
mockgen -source=internal/interfaces/repository.go \
    -destination=internal/implementations/mocks/mock_repository.go \
    -package=mocks

# Generate service mocks (for handler layer tests)
mockgen -source=internal/interfaces/service.go \
    -destination=internal/implementations/mocks/mock_service.go \
    -package=mocks
```

### Step 5: Implement Repository Layer

Create PostgreSQL implementations that pass the repository test runners:

```bash
# Create directory
mkdir -p internal/implementations/repository

# Files to create:
# - postgres_connection_repository.go
# - postgres_metadata_repository.go
# - postgres_query_repository.go
# - postgres_transaction_repository.go

# For each repository, create:
# 1. Implementation file: internal/implementations/repository/postgres_xxx_repository.go
# 2. Test file that binds to runner: internal/implementations/repository/postgres_xxx_repository_test.go
```

#### Example: Connection Repository

**Implementation** (`internal/implementations/repository/postgres_connection_repository.go`):
```go
package repository

import (
    "context"
    "database/sql"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
    _ "github.com/lib/pq"
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
    // TODO: Implement validation logic
    return nil
}

func (r *PostgresConnectionRepository) TestConnection(ctx context.Context, connStr string) error {
    // TODO: Implement connection test
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return err
    }
    defer db.Close()
    return db.PingContext(ctx)
}

// ... implement other methods
```

**Test Binding** (`internal/implementations/repository/postgres_connection_repository_test.go`):
```go
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

**Run the test:**
```bash
go test ./internal/implementations/repository -v -run TestPostgresConnectionRepository
```

Repeat this pattern for:
- `postgres_metadata_repository.go` / `_test.go`
- `postgres_query_repository.go` / `_test.go`
- `postgres_transaction_repository.go` / `_test.go`

### Step 6: Implement Service Layer

Create service implementations that use repositories:

```bash
# Create directory
mkdir -p internal/implementations/service

# Files to create:
# - auth_service.go
# - metadata_service.go
# - query_service.go
# - data_explorer_service.go
# - transaction_service.go
```

#### Example: Auth Service

**Implementation** (`internal/implementations/service/auth_service.go`):
```go
package service

import (
    "context"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type AuthService struct {
    connRepo     interfaces.ConnectionRepository
    metadataRepo interfaces.MetadataRepository
}

func NewAuthService(
    connRepo interfaces.ConnectionRepository,
    metadataRepo interfaces.MetadataRepository,
) interfaces.AuthService {
    return &AuthService{
        connRepo:     connRepo,
        metadataRepo: metadataRepo,
    }
}

func (s *AuthService) Login(ctx context.Context, input domain.LoginInput) (*domain.Session, error) {
    // TODO: Implement login logic with connection probe
    return nil, nil
}

// ... implement other methods
```

**Test Binding** (`internal/implementations/service/auth_service_test.go`):
```go
package service_test

import (
    "testing"
    "github.com/kamil5b/lumen-pg/internal/implementations/service"
    "github.com/kamil5b/lumen-pg/internal/testrunners"
)

func TestAuthService(t *testing.T) {
    constructor := func(
        connRepo interfaces.ConnectionRepository,
        metadataRepo interfaces.MetadataRepository,
    ) interfaces.AuthService {
        return service.NewAuthService(connRepo, metadataRepo)
    }
    testrunners.AuthServiceRunner(t, constructor)
}
```

**Run the test:**
```bash
go test ./internal/implementations/service -v -run TestAuthService
```

### Step 7: Implement Handler Layer

Create HTTP handlers that use services:

```bash
# Create directory
mkdir -p internal/implementations/handler

# Files to create:
# - auth_handler.go
# - main_view_handler.go
# - query_editor_handler.go
# - erd_viewer_handler.go
# - transaction_handler.go
```

### Step 8: Wire Everything Together

Create `cmd/lumen-pg/main.go`:

```go
package main

import (
    "flag"
    "log"
    "net/http"
    
    "github.com/go-chi/chi/v5"
    "github.com/kamil5b/lumen-pg/internal/implementations/repository"
    "github.com/kamil5b/lumen-pg/internal/implementations/service"
    "github.com/kamil5b/lumen-pg/internal/implementations/handler"
)

func main() {
    dbConnStr := flag.String("db", "", "PostgreSQL superadmin connection string")
    port := flag.String("port", "8080", "HTTP server port")
    flag.Parse()

    if *dbConnStr == "" {
        log.Fatal("Database connection string is required (-db flag)")
    }

    // Initialize repositories
    connRepo := repository.NewPostgresConnectionRepository(*dbConnStr)
    metadataRepo := repository.NewPostgresMetadataRepository()
    queryRepo := repository.NewPostgresQueryRepository()
    txRepo := repository.NewPostgresTransactionRepository()

    // Initialize services
    authService := service.NewAuthService(connRepo, metadataRepo)
    metadataService := service.NewMetadataService(connRepo, metadataRepo)
    queryService := service.NewQueryService(connRepo, queryRepo)
    dataExplorerService := service.NewDataExplorerService(connRepo, queryRepo)
    transactionService := service.NewTransactionService(connRepo, txRepo)

    // Initialize handlers
    authHandler := handler.NewAuthHandler(authService)
    mainViewHandler := handler.NewMainViewHandler(dataExplorerService, metadataService, transactionService)
    queryEditorHandler := handler.NewQueryEditorHandler(queryService)
    erdViewerHandler := handler.NewERDViewerHandler(metadataService)
    transactionHandler := handler.NewTransactionHandler(transactionService)

    // Setup router
    r := chi.NewRouter()
    
    authHandler.RegisterRoutes(r)
    mainViewHandler.RegisterRoutes(r)
    queryEditorHandler.RegisterRoutes(r)
    erdViewerHandler.RegisterRoutes(r)
    transactionHandler.RegisterRoutes(r)

    // Start server
    log.Printf("Starting server on port %s", *port)
    log.Fatal(http.ListenAndServe(":"+*port, r))
}
```

**Run the application:**
```bash
go run cmd/lumen-pg/main.go -db "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
```

## Testing Strategy

### Repository Tests (Integration Tests)
```bash
# Requires Docker for testcontainers
go test ./internal/implementations/repository -v
```

### Service Tests (Unit Tests with Mocks)
```bash
# Requires generated mocks
go test ./internal/implementations/service -v
```

### Handler Tests (Unit Tests with Mocks)
```bash
# Requires generated mocks
go test ./internal/implementations/handler -v
```

### All Tests
```bash
go test ./... -v
```

### Coverage Report
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```

## Implementation Tips

### 1. Start with Repository Layer
- Repositories use real PostgreSQL (via testcontainers in tests)
- Implement SQL queries carefully with parameterization
- Focus on metadata loading and role permissions

### 2. Move to Service Layer
- Services coordinate repositories
- Implement business logic here
- Use mocked repositories in tests

### 3. Finish with Handler Layer
- Handlers translate HTTP to/from services
- Use HTMX for dynamic updates
- Use mocked services in tests

### 4. Follow TDD
- Tests are already written in test runners
- Implement until tests pass
- Don't skip tests or mark them as passing without implementation

## Key Implementation Areas

Based on `TEST_PLAN.md` and `REQUIREMENT.md`:

### Priority 1: Core Infrastructure
1. Connection management and validation
2. Metadata loading with role permissions
3. Session management with encrypted cookies

### Priority 2: Authentication
1. Login with connection probing
2. Session creation and validation
3. Cookie management (username and encrypted password)

### Priority 3: Data Browsing
1. Table data loading with cursor pagination
2. WHERE clause filtering (SQL injection safe)
3. Column sorting
4. Hard limit of 1000 rows with total count display

### Priority 4: Query Execution
1. SQL query execution (SELECT, DDL, DML)
2. Multiple query support (semicolon-separated)
3. Result pagination (1000 row limit)
4. Parameterized queries for safety

### Priority 5: Transactions
1. Transaction start/commit/rollback
2. Operation buffering (INSERT, UPDATE, DELETE)
3. 1-minute timeout
4. Isolated per-user transactions

### Priority 6: ERD Viewer
1. Schema visualization
2. Table and relationship display
3. Interactive navigation

## Resources

- **Test Plan**: See `TEST_PLAN.md` for all test cases
- **Requirements**: See `REQUIREMENT.md` for user stories
- **Test Runners**: See `internal/testrunners/README.md` for testing details
- **Architecture**: See `README.md` for overall architecture

## Common Pitfalls to Avoid

1. ❌ **Don't skip test runners** - They ensure your implementation meets the spec
2. ❌ **Don't bypass interfaces** - Always code to interfaces, not implementations
3. ❌ **Don't forget SQL injection prevention** - Use parameterized queries everywhere
4. ❌ **Don't ignore role permissions** - RBAC is a core requirement
5. ❌ **Don't skip testcontainers** - Repository tests need real PostgreSQL
6. ❌ **Don't forget the 1000 row limit** - Hard limit for pagination in all views

## Success Criteria

Your implementation is complete when:
- ✅ All repository test runners pass
- ✅ All service test runners pass (after uncommenting skip markers)
- ✅ All handler test runners pass (after uncommenting skip markers)
- ✅ Application starts and connects to PostgreSQL
- ✅ All user stories from REQUIREMENT.md are implemented
- ✅ Code coverage is above 90% for domain logic

## Getting Help

- Read the test runners to understand expected behavior
- Check `TEST_PLAN.md` for detailed test cases
- Review `REQUIREMENT.md` for user story context
- Look at the interfaces for method signatures and return types

Good luck with the implementation! 🚀
