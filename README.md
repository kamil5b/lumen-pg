# Lumen-PG - PostgreSQL DBMS Web Client

A minimalist, lightweight, web-based DBMS client application built using Go and HTMX, designed to connect to a single PostgreSQL database instance.

## Project Status

**Current Phase**: Test-First Architecture Setup (Steps 1-3 Complete)

✅ Step 1: Domain Layer - Pure business types defined  
✅ Step 2: Interfaces Layer - Repository, Service, and Handler contracts defined  
✅ Step 3: Test Runners - Test-first specifications written  
⏳ Step 4: Implementation - To be done  

## Architecture

This project follows a **test-first, layered architecture** pattern:

```
internal/
├─ domain/          # Pure business types (no dependencies)
│  ├─ user.go
│  ├─ session.go
│  ├─ metadata.go
│  ├─ query.go
│  ├─ transaction.go
│  └─ erd.go
│
├─ interfaces/      # Interface contracts
│  ├─ repository.go # DB operations
│  ├─ service.go    # Business logic
│  └─ handler.go    # HTTP handlers
│
├─ implementations/ # Concrete implementations (TO BE CREATED)
│  ├─ mocks/        # Generated mocks (mockgen)
│  ├─ repository/   # PostgreSQL repositories
│  ├─ service/      # Service implementations
│  └─ handler/      # HTTP handlers
│
└─ testrunners/     # Test-first specifications
   ├─ repository_runner_test.go
   ├─ service_runner_test.go
   ├─ handler_runner_test.go
   └─ README.md
```

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Pure business types
   - No database, HTTP, or framework logic
   - Examples: User, Session, Metadata, Query, Transaction, ERD

2. **Interfaces Layer** (`internal/interfaces/`)
   - Contracts for all layers
   - Repository interfaces for data access
   - Service interfaces for business logic
   - Handler interfaces for HTTP endpoints

3. **Test Runners** (`internal/testrunners/`)
   - Test-first specifications
   - Repository runners use **testcontainers** (real PostgreSQL)
   - Service runners use **gomock** (mocked repositories)
   - Handler runners use **gomock** (mocked services)

4. **Implementations** (TO BE CREATED in `internal/implementations/`)
   - Concrete repository implementations (PostgreSQL)
   - Service layer implementations
   - HTTP handler implementations

## Features (Planned)

Based on requirements from `REQUIREMENT.md`:

### Core Functionality
- **Single PostgreSQL Instance Connection** - Connect via superadmin connection string
- **Role-Based Access Control (RBAC)** - Based on PostgreSQL roles and permissions
- **Metadata Caching** - In-memory cache of databases, schemas, tables, relations, and permissions
- **Multi-User Sessions** - Isolated sessions with encrypted cookies

### User Stories
1. **Setup & Configuration** - Initialize and cache metadata with role permissions
2. **Authentication & Identity** - Login with PostgreSQL credentials, connection probing
3. **ERD Viewer** - Visualize entity-relationship diagrams
4. **Manual Query Editor** - Execute SQL queries with pagination (1000 row limit)
5. **Main View & Data Interaction** - Browse tables, cursor pagination, transactions
6. **Isolation** - Multi-user support with isolated sessions
7. **Security** - Parameterized queries, encrypted cookies, session timeouts

## Getting Started

### Prerequisites

- Go 1.25+
- Docker (for testcontainers)
- PostgreSQL instance

### Installation

```bash
# Clone the repository
git clone https://github.com/kamil5b/lumen-pg.git
cd lumen-pg

# Install dependencies
go mod download

# Generate mocks (required for service/handler tests)
go install go.uber.org/mock/mockgen@latest
mockgen -source=internal/interfaces/repository.go -destination=internal/implementations/mocks/mock_repository.go -package=mocks
mockgen -source=internal/interfaces/service.go -destination=internal/implementations/mocks/mock_service.go -package=mocks
```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run repository tests (requires Docker)
go test ./internal/testrunners -v -run Repository

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Development Workflow (TDD)

This project uses **Test-Driven Development (TDD)**:

1. **Tests are already written** in `internal/testrunners/`
2. **Implement** the interfaces in `internal/implementations/`
3. **Bind** your implementation to the test runner
4. **Run tests** and iterate until they pass

### Example: Implementing a Repository

```go
// 1. Create implementation
// internal/implementations/repository/postgres_connection_repository.go
package repository

import (
    "context"
    "database/sql"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type PostgresConnectionRepository struct {
    superAdminConnStr string
}

func NewPostgresConnectionRepository(connStr string) interfaces.ConnectionRepository {
    return &PostgresConnectionRepository{superAdminConnStr: connStr}
}

func (r *PostgresConnectionRepository) ValidateConnectionString(connStr string) error {
    // Implementation
    return nil
}

// ... implement other methods

// 2. Bind to test runner
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

## Test Runner Pattern

Each runner follows this pattern:

```go
func Runner(t *testing.T, constructor func(dependencies) Interface)
```

- **Repository Runners**: Use testcontainers (real PostgreSQL)
- **Service Runners**: Use gomock (mocked repositories)
- **Handler Runners**: Use gomock (mocked services)

See `internal/testrunners/README.md` for detailed documentation.

## Project Structure Benefits

✅ **Test-First** - All tests written before implementation  
✅ **Fully Symmetric** - Interfaces for all layers  
✅ **Easy to Swap** - Interface-driven design  
✅ **Real DB Testing** - Repository tests use testcontainers  
✅ **Isolated Testing** - Service/handler tests use mocks  
✅ **Clean Composition** - Clear separation of concerns  

## Next Steps

1. ✅ Define domain layer (Done)
2. ✅ Define interfaces (Done)
3. ✅ Write test runners (Done)
4. ⏳ Implement repository layer
5. ⏳ Implement service layer
6. ⏳ Implement handler layer
7. ⏳ Create main.go and wire up dependencies
8. ⏳ Add E2E tests

## Documentation

- `REQUIREMENT.md` - Full requirements and user stories
- `TEST_PLAN.md` - Comprehensive test plan with all test cases
- `internal/testrunners/README.md` - Test runner documentation

## Testing Strategy

See `TEST_PLAN.md` for the complete testing strategy:

- **Unit Tests** - Domain logic and use cases (90% coverage goal)
- **Integration Tests** - Repository layer with real PostgreSQL
- **E2E Tests** - Complete user flows via HTTP

## Contributing

1. All implementations must pass the test runners
2. Follow the existing architecture patterns
3. Add tests for new functionality
4. Maintain interface-driven design

## License

[License information]

## Related Documents

- [REQUIREMENT.md](./REQUIREMENT.md) - Full requirements
- [TEST_PLAN.md](./TEST_PLAN.md) - Testing strategy and test cases
- [Test Runners README](./internal/testrunners/README.md) - Test runner documentation
