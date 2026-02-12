# Lumen-PG - PostgreSQL DBMS Web Application

A minimalist, lightweight, web-based DBMS client application built using Go and HTMX, designed to connect to a single PostgreSQL database instance.

## Overview

This application provides a simple interface to:
- Explore database schemas with role-based access control
- View and edit table data with transaction support
- Execute manual SQL queries
- Visualize entity-relationship diagrams (ERDs)

## Features

### Story 1: Setup & Configuration
- Validates PostgreSQL connection strings
- Caches metadata (databases, schemas, tables, columns, relationships)
- Implements role-based access control (RBAC)
- Maps accessible resources per PostgreSQL role

### Story 2: Authentication & Identity
- PostgreSQL credential-based login
- Session management with encrypted cookies
- Connection probe to verify user has accessible resources
- Data Explorer sidebar with role-aware table listing

### Story 3: ERD Viewer
- Dynamic entity-relationship diagrams
- Zoom and pan controls
- Table details on click

### Story 4: Manual Query Editor
- SQL syntax highlighting
- Multi-query execution (semicolon-separated)
- Pagination for large result sets (1000 row hard limit)
- DDL/DML query support with feedback

### Story 5: Main View & Data Interaction
- Table data browsing with WHERE filters
- Cursor-based pagination (50 rows/page, 1000 max)
- Sortable columns
- **Transaction Mode**:
  - Inline cell editing
  - Buffered operations
  - 1-minute timeout with commit/rollback
- **Read-Only Mode**:
  - Foreign key navigation
  - Primary key reference viewing

### Story 6: Isolation
- Multi-user support with isolated sessions
- Transaction isolation per user
- Role-based permission enforcement

### Story 7: Security
- Parameterized queries (SQL injection prevention)
- Encrypted password cookies
- Session timeouts
- HTTPS support

## Project Structure

```
lumen-pg/
├── internal/
│   ├── domain/              # Business entities
│   │   ├── connection.go    # PostgreSQL connection config
│   │   ├── metadata.go      # Database schema metadata
│   │   ├── role.go          # RBAC role & permissions
│   │   ├── session.go       # User session management
│   │   ├── query.go         # Query execution results
│   │   └── transaction.go   # Transaction state
│   │
│   ├── interfaces/          # Interface contracts
│   │   ├── connection.go    # Connection operations
│   │   ├── metadata.go      # Metadata loading
│   │   ├── query.go         # Query execution
│   │   ├── session.go       # Session management
│   │   ├── transaction.go   # Transaction handling
│   │   └── handler.go       # HTTP handlers
│   │
│   ├── implementations/     # Concrete implementations
│   │   ├── mocks/           # Auto-generated mocks
│   │   ├── repository/      # PostgreSQL repositories
│   │   ├── usecase/         # Business logic use cases
│   │   └── handler/         # HTTP handlers
│   │
│   └── testrunners/         # Test specifications (TDD)
│       ├── *_usecase_runner_test.go      # Unit test specs
│       └── *_repository_runner_test.go   # Integration test specs
│
├── REQUIREMENT.md           # Full requirements specification
├── TEST_PLAN.md            # Comprehensive test plan
└── README.md               # This file
```

## TDD Approach

This project follows Test-Driven Development:

1. **Domain Models** defined first (`internal/domain/`)
2. **Interfaces** define contracts (`internal/interfaces/`)
3. **Mocks** generated from interfaces (`internal/implementations/mocks/`)
4. **Test Runners** specify expected behavior (`internal/testrunners/`)
5. **Implementations** created to pass test runners

### Test Runner Pattern

Test runners define the contract that implementations must satisfy:

```go
// Example: Metadata use case test
func TestMetadataUseCase(t *testing.T) {
    testrunners.MetadataUseCaseRunner(t, implementations.NewMetadataUseCase)
}

// Example: Query repository integration test
func TestQueryRepository(t *testing.T) {
    testrunners.QueryRepositoryRunner(t, implementations.NewPostgresQueryRepository)
}
```

## Running Tests

```bash
# Build all packages
go build ./...

# Run unit tests (with mocks)
go test ./internal/testrunners -run UseCase

# Run integration tests (with testcontainers)
go test ./internal/testrunners -run Repository

# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover
```

## Dependencies

- **github.com/go-chi/chi/v5** - HTTP router
- **github.com/lib/pq** - PostgreSQL driver
- **github.com/stretchr/testify** - Testing assertions
- **github.com/testcontainers/testcontainers-go** - Integration testing with real PostgreSQL
- **go.uber.org/mock** - Mock generation

## Development

### Regenerating Mocks

After modifying interfaces:

```bash
# Generate all mocks
mockgen -source=internal/interfaces/connection.go -destination=internal/implementations/mocks/connection_mock.go -package=mocks
mockgen -source=internal/interfaces/metadata.go -destination=internal/implementations/mocks/metadata_mock.go -package=mocks
mockgen -source=internal/interfaces/query.go -destination=internal/implementations/mocks/query_mock.go -package=mocks
mockgen -source=internal/interfaces/session.go -destination=internal/implementations/mocks/session_mock.go -package=mocks
mockgen -source=internal/interfaces/transaction.go -destination=internal/implementations/mocks/transaction_mock.go -package=mocks
```

### Next Steps

1. Implement use cases in `internal/implementations/usecase/`
2. Implement repositories in `internal/implementations/repository/`
3. Implement HTTP handlers in `internal/implementations/handler/`
4. Create main application entry point
5. Add HTMX templates for UI

## Key Constraints

- **Pagination Limits**: 
  - Query results: 1000 rows max, shows total count
  - Main view: 50 rows/page, 1000 rows max
- **Transaction Timeout**: 1 minute
- **Single Database Instance**: Application connects to one PostgreSQL instance
- **Stateless**: No persistent storage (metadata cached in-memory)

## License

See LICENSE file for details.
