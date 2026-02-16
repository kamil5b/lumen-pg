# Test Runners - Lumen-PG

This directory contains test runners and fixtures for Lumen-PG following a strict TDD approach.

## Testing Strategy

### 1. Unit Tests
Unit tests verify individual components in isolation using mocks/stubs.
- Test domain models (Transaction, Session, Cursor, Query Splitter, WHERE Validator)
- Test usecase logic with mocked repositories
- Test business rules and validation

### 2. Integration Tests
Integration tests verify component interactions with real PostgreSQL database using testcontainers.
- Test repository implementations against real PostgreSQL
- Test data persistence and retrieval
- Test transaction and session management

### 3. E2E Tests
End-to-end tests verify complete user flows through HTTP API with HTMX responses.
- Test complete authentication flows
- Test data exploration and query execution
- Test multi-user isolation
- Test security features

---

## Project Structure

```
lumen-pg/
├─ internal/
│  ├─ domain/
│  │  ├─ types.go              (DatabaseMetadata, User, Session, QueryResult, etc.)
│  │  ├─ constants.go
│  │  └─ dto/
│  ├─ interfaces/
│  │  ├─ repository/           (Connection, Database, Metadata, Session, Transaction, etc.)
│  │  ├─ usecase/              (Authentication, DataView, Query, Transaction, ERD, etc.)
│  │  ├─ handler/              (Login, MainView, QueryEditor, ERDViewer, etc.)
│  │  └─ middleware/           (Authentication, Authorization, Validation, etc.)
│  └─ testrunners/
│     ├─ repository/           (repository_runner_test.go)
│     ├─ usecase/              (usecase_runner_test.go)
│     ├─ handler/                  (end-to-end test runners)
│     └─ mocks/                (generated mock files)
├─ go.mod
├─ go.sum
└─ README.md
```

---

## TDD Workflow - Strict Order

Follow this phase-by-phase TDD approach for development:

### Phase 1: Core Domain Models (Unit Tests)

Define and test domain models **first**, before any interfaces.

**Models to test:**
- Transaction domain (state management)
- Session domain (lifecycle and validation)
- Cursor domain (pagination state)
- Query Splitter (parsing multiple queries)
- WHERE Validator (injection prevention)

**Example:**
```go
// internal/testrunners/domain_runner_test.go
func TestTransactionDomain(t *testing.T) {
    // Test transaction state transitions
    // Test timeout expiration
    // Test edit buffering
}
```

---

### Phase 2: Repository Interfaces (Unit Tests)

Define interfaces **before** implementation.

**Interfaces to define:**
- `ConnectionRepository` - PostgreSQL connection management
- `MetadataRepository` - database/schema/table metadata
- `QueryRepository` - SQL query execution
- `SessionRepository` - user session persistence
- `TransactionRepository` - transaction state management
- `RBACRepository` - role-based access control
- `EncryptionRepository` - password/data encryption
- `CacheRepository` - caching layer
- `ClockRepository` - time operations
- `LoggerRepository` - logging

**Example:**
```go
// internal/testrunners/repository/repository_runner_test.go
type RepositoryConstructor func(db *sql.DB) repository.MetadataRepository

func MetadataRepositoryRunner(t *testing.T, constructor RepositoryConstructor) {
    t.Helper()
    // Use testcontainers for PostgreSQL
    // Define test cases for metadata retrieval
}
```

---

### Phase 3: Usecase Interfaces & Unit Tests

Define usecase interfaces and test with **mocked repositories**.

**Usecases to implement:**
- `AuthenticationUsecase` - login, logout, session validation
- `MetadataUsecase` - fetch and cache database metadata
- `DataViewUsecase` - load table data with filtering/pagination
- `QueryUsecase` - execute SQL queries with splitting
- `TransactionUsecase` - manage buffered edits/deletes
- `ERDUsecase` - generate entity relationship diagrams
- `RBACUsecase` - role-based access control logic
- `SecurityUsecase` - SQL injection prevention, encryption

**Example:**
```go
// internal/testrunners/usecase/usecase_runner_test.go
type UsecaseConstructor func(repo repository.MetadataRepository) usecase.DataViewUsecase

func DataViewUsecaseRunner(t *testing.T, constructor UsecaseConstructor) {
    t.Helper()
    
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockMetadataRepository(ctrl)
    
    uc := constructor(mockRepo)
    
    t.Run("LoadTableData with filters", func(t *testing.T) {
        mockRepo.EXPECT().
            GetTableMetadata(gomock.Any(), gomock.Any()).
            Return(&domain.TableMetadata{}, nil)
        
        result, err := uc.LoadTableData(context.Background(), &domain.TableDataParams{})
        require.NoError(t, err)
        require.NotNil(t, result)
    })
}
```

---

### Phase 4: Repository Implementations (Integration Tests)

Implement repositories and test against **real PostgreSQL** using testcontainers.

**Example:**
```go
// internal/testrunners/repository/postgres_runner_test.go
func TestPostgresMetadataRepository(t *testing.T) {
    repository.MetadataRepositoryRunner(t, NewPostgresMetadataRepository)
}

func TestPostgresConnectionRepository(t *testing.T) {
    repository.ConnectionRepositoryRunner(t, NewPostgresConnectionRepository)
}

func TestPostgresQueryRepository(t *testing.T) {
    repository.QueryRepositoryRunner(t, NewPostgresQueryRepository)
}
```

**Key integration test areas:**
- Connect to real PostgreSQL
- Load real database metadata
- Execute real queries with permission checks
- Test transaction isolation
- Test session persistence
- Test concurrent access

---

### Phase 5: Web Layer (Unit + Integration Tests)

Test handlers and middleware with **mocked usecases**.

**Handler tests:**
- Login handler with authentication usecase mocked
- Main view handler with data view usecase mocked
- Query editor handler with query usecase mocked
- ERD viewer handler with ERD usecase mocked
- Transaction handler with transaction usecase mocked

**Middleware tests:**
- Authentication middleware (session validation)
- Authorization middleware (permission checking)
- Validation middleware (input sanitization)
- Error handling middleware (error responses)

**Example:**
```go
// internal/testrunners/handler/handler_runner_test.go
type HandlerConstructor func(uc usecase.QueryUsecase) handler.QueryEditorHandler

func QueryEditorHandlerRunner(t *testing.T, constructor HandlerConstructor) {
    t.Helper()
    
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockUsecase := mocks.NewMockQueryUsecase(ctrl)
    
    h := constructor(mockUsecase)
    
    t.Run("HandleExecuteQuery returns results", func(t *testing.T) {
        mockUsecase.EXPECT().
            ExecuteQuery(gomock.Any(), gomock.Any()).
            Return(&domain.QueryResult{}, nil)
        
        req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
        rec := httptest.NewRecorder()
        
        h.HandleExecuteQuery(rec, req)
        
        require.Equal(t, http.StatusOK, rec.Code)
    })
}
```

---

### Phase 6: E2E Tests

Test complete user flows end-to-end.

**Story-based E2E tests:**

1. **Story 1: Setup & Configuration**
   - Connection string validation
   - PostgreSQL connection probe
   - Metadata initialization

2. **Story 2: Authentication & Identity**
   - Login form validation
   - Session creation and validation
   - Logout flow
   - Protected route access

3. **Story 3: ERD Viewer**
   - ERD generation from schema
   - Zoom and pan controls
   - Table navigation

4. **Story 4: Manual Query Editor**
   - Single/multiple query execution
   - Offset pagination
   - Query result display
   - Error handling

5. **Story 5: Main View & Data Interaction**
   - Table data loading with filters/sorting
   - Cursor pagination
   - Transaction management
   - Cell editing and row operations
   - Foreign key navigation

6. **Story 6: Isolation**
   - Multi-user session isolation
   - Transaction isolation
   - Permission-based resource isolation

7. **Story 7: Security & Best Practices**
   - SQL injection prevention
   - Password encryption
   - Cookie security
   - Session timeout

---

## Setting Up Test Runners

### Step 1: Define Domain Models

```go
// internal/domain/types.go
type User struct {
    Username     string
    DatabaseName string
    SchemaName   string
    TableName    string
}

type Session struct {
    ID        string
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

// ... other domain types
```

### Step 2: Define Repository Interfaces

```go
// internal/interfaces/repository/metadata_repository.go
package repository

import (
    "context"
    "lumen-pg/internal/domain"
)

type MetadataRepository interface {
    GetDatabaseMetadata(ctx context.Context) (*domain.DatabaseMetadata, error)
    GetSchemaMetadata(ctx context.Context, schema string) (*domain.SchemaMetadata, error)
    GetTableMetadata(ctx context.Context, schema, table string) (*domain.TableMetadata, error)
}
```

### Step 3: Generate Mocks

Install gomock:
```bash
go install github.com/golang/mock/mockgen@latest
```

Generate mocks for all repository interfaces:
```bash
mockgen -source=internal/interfaces/repository/metadata_repository.go \
  -destination=internal/testrunners/mocks/mock_metadata_repository.go \
  -package=mocks

mockgen -source=internal/interfaces/repository/database_repository.go \
  -destination=internal/testrunners/mocks/mock_database_repository.go \
  -package=mocks

mockgen -source=internal/interfaces/repository/session_repository.go \
  -destination=internal/testrunners/mocks/mock_session_repository.go \
  -package=mocks

# Generate mocks for all usecase interfaces similarly
mockgen -source=internal/interfaces/usecase/authentication_usecase.go \
  -destination=internal/testrunners/mocks/mock_authentication_usecase.go \
  -package=mocks
```

### Step 4: Create Repository Test Runners

```go
// internal/testrunners/repository/repository_runner_test.go
package repository

import (
    "context"
    "database/sql"
    "testing"
    
    "github.com/stretchr/testify/require"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    _ "github.com/lib/pq"
    
    "lumen-pg/internal/domain"
    "lumen-pg/internal/interfaces/repository"
)

type RepositoryConstructor func(db *sql.DB) repository.MetadataRepository

func MetadataRepositoryRunner(t *testing.T, constructor RepositoryConstructor) {
    t.Helper()
    
    ctx := context.Background()
    
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("user"),
        postgres.WithPassword("pass"),
    )
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    connStr, err := container.ConnectionString(ctx)
    require.NoError(t, err)
    
    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)
    defer db.Close()
    
    repo := constructor(db)
    
    t.Run("GetDatabaseMetadata returns valid structure", func(t *testing.T) {
        metadata, err := repo.GetDatabaseMetadata(ctx)
        require.NoError(t, err)
        require.NotNil(t, metadata)
    })
    
    t.Run("GetSchemaMetadata returns schema", func(t *testing.T) {
        metadata, err := repo.GetSchemaMetadata(ctx, "public")
        require.NoError(t, err)
        require.NotNil(t, metadata)
    })
}
```

### Step 5: Create Usecase Test Runners

```go
// internal/testrunners/usecase/usecase_runner_test.go
package usecase

import (
    "context"
    "testing"
    
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
    
    "lumen-pg/internal/domain"
    "lumen-pg/internal/interfaces/repository"
    "lumen-pg/internal/interfaces/usecase"
    "lumen-pg/internal/testrunners/mocks"
)

type UsecaseConstructor func(repo repository.MetadataRepository) usecase.DataViewUsecase

func DataViewUsecaseRunner(t *testing.T, constructor UsecaseConstructor) {
    t.Helper()
    
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockMetadataRepository(ctrl)
    
    uc := constructor(mockRepo)
    
    t.Run("LoadTableData returns QueryResult", func(t *testing.T) {
        mockRepo.EXPECT().
            GetTableMetadata(gomock.Any(), gomock.Any(), gomock.Any()).
            Return(&domain.TableMetadata{}, nil)
        
        result, err := uc.LoadTableData(context.Background(), &domain.TableDataParams{
            Database: "testdb",
            Schema:   "public",
            Table:    "users",
            Limit:    10,
            Offset:   0,
        })
        
        require.NoError(t, err)
        require.NotNil(t, result)
        require.IsType(t, &domain.QueryResult{}, result)
    })
}
```

### Step 6: Create Handler E2E Test Runners

```go
// internal/testrunners/handler/handler_runner_test.go
package handler

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
    
    "lumen-pg/internal/domain"
    "lumen-pg/internal/interfaces/handler"
    "lumen-pg/internal/interfaces/usecase"
    "lumen-pg/internal/testrunners/mocks"
)

type HandlerConstructor func(uc usecase.QueryUsecase) handler.QueryEditorHandler

func QueryEditorHandlerRunner(t *testing.T, constructor HandlerConstructor) {
    t.Helper()
    
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockUsecase := mocks.NewMockQueryUsecase(ctrl)
    
    h := constructor(mockUsecase)
    
    t.Run("HandleExecuteQuery returns 200 with results", func(t *testing.T) {
        expectedResult := &domain.QueryResult{
            Columns:  []string{"id", "name"},
            Rows:     []map[string]interface{}{},
            RowCount: 0,
        }
        
        mockUsecase.EXPECT().
            ExecuteQuery(gomock.Any(), gomock.Any()).
            Return(expectedResult, nil)
        
        req := httptest.NewRequest(http.MethodPost, "/api/query", nil)
        rec := httptest.NewRecorder()
        
        h.HandleExecuteQuery(rec, req)
        
        require.Equal(t, http.StatusOK, rec.Code)
    })
}
```

### Step 7: Wire Everything Together

**Implementation files implement runner functions:**

```go
// internal/implementations/repository/postgres_metadata.go
package repository

import (
    "database/sql"
    "lumen-pg/internal/interfaces/repository"
    "lumen-pg/internal/testrunners"
)

type PostgresMetadataRepository struct {
    db *sql.DB
}

func NewPostgresMetadataRepository(db *sql.DB) repository.MetadataRepository {
    return &PostgresMetadataRepository{db: db}
}

// ... implement interface methods
```

**Test file calls runner:**

```go
// internal/implementations/repository/postgres_metadata_test.go
package repository

import (
    "testing"
    "lumen-pg/internal/testrunners/repository"
)

func TestPostgresMetadataRepository(t *testing.T) {
    repository.MetadataRepositoryRunner(t, NewPostgresMetadataRepository)
}
```

---

## Running Tests

### Run All Tests
```bash
go test ./...
```

### Run Unit Tests Only
```bash
go test ./internal/testrunners/... -v
```

### Run Integration Tests Only
```bash
go test ./internal/testrunners/repository/... -v -tags=integration
```

### Run with Coverage
```bash
go test ./... -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Specific Story Tests
```bash
# Story 1: Setup & Configuration
go test ./internal/testrunners/domain/... -v -run "Setup"

# Story 2: Authentication
go test ./internal/testrunners/usecase/... -v -run "Auth"

# Story 4: Query Editor
go test ./internal/testrunners/usecase/... -v -run "Query"

# Story 5: Main View
go test ./internal/testrunners/usecase/... -v -run "DataView"

# E2E Tests
go test ./internal/testrunners/handler/... -v
```

---

## Best Practices

✅ **DO:**
- Define domain models first
- Define interfaces before implementation
- Generate mocks with gomock
- Use testcontainers for integration tests
- Test business logic with mocked dependencies
- Test implementations against real databases
- Write test runners as reusable patterns

❌ **DON'T:**
- Test implementation details
- Mock across layer boundaries
- Skip integration tests
- Hardcode test data
- Test unrelated concerns together

---

## Notes

- All tests follow the TDD workflow: define → test → implement
- Test runners are reusable patterns for testing similar components
- Repository tests use testcontainers for real PostgreSQL
- Usecase tests use gomock for mocked repositories
- Handler tests use gomock for mocked usecases
- E2E tests verify complete user flows
