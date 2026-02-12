# TDD Implementation Guide

This document provides step-by-step instructions for implementing the TDD structure that has been set up.

## What's Already Done

✅ **Project Structure**: Complete directory layout  
✅ **Domain Models**: User entity defined  
✅ **Interfaces**: UserRepository, UserService, UserHandler  
✅ **Mocks**: Auto-generated mocks for testing  
✅ **Test Runners**: Specs that define expected behavior  

## Next Steps: Implementing the Actual Code

### Phase 1: Repository Implementation

Create `internal/implementations/repository/postgres_user_repository.go`:

```go
package repository

import (
    "context"
    "database/sql"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type PostgresUserRepository struct {
    db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) interfaces.UserRepository {
    return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    // TODO: Implement database save
    // 1. Create table if not exists
    // 2. INSERT or UPDATE user
    return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    // TODO: Implement database query
    // 1. SELECT user by ID
    // 2. Scan into domain.User
    return nil, nil
}
```

Then create a test file `internal/implementations/repository/postgres_user_repository_test.go`:

```go
package repository

import (
    "testing"
    "github.com/kamil5b/lumen-pg/internal/testrunners"
)

func TestPostgresUserRepository(t *testing.T) {
    testrunners.UserRepositoryRunner(t, NewPostgresUserRepository)
}
```

Run the test:
```bash
go test ./internal/implementations/repository -v
```

The test will fail initially. Implement the methods until the test passes.

### Phase 2: Service Implementation

Create `internal/implementations/service/user_service.go`:

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "time"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type UserService struct {
    repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) interfaces.UserService {
    return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
    // TODO: Implement user creation logic
    // 1. Generate unique ID
    // 2. Create User entity
    // 3. Save via repository
    user := &domain.User{
        ID:        uuid.New().String(),
        Email:     input.Email,
        Name:      input.Name,
        CreatedAt: time.Now(),
    }
    
    err := s.repo.Save(ctx, user)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
    // TODO: Implement user retrieval
    return s.repo.FindByID(ctx, id)
}
```

Then create a test file `internal/implementations/service/user_service_test.go`:

```go
package service

import (
    "testing"
    "github.com/kamil5b/lumen-pg/internal/testrunners"
)

func TestUserService(t *testing.T) {
    testrunners.UserServiceRunner(t, NewUserService)
}
```

Run the test:
```bash
go test ./internal/implementations/service -v
```

### Phase 3: Handler Implementation

Create `internal/implementations/handler/user_handler.go`:

```go
package handler

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/kamil5b/lumen-pg/internal/domain"
    "github.com/kamil5b/lumen-pg/internal/interfaces"
)

type UserHandler struct {
    service interfaces.UserService
}

func NewUserHandler(service interfaces.UserService) interfaces.UserHandler {
    return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
    r.Post("/users", h.CreateUser)
    r.Get("/users/{id}", h.GetUser)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement HTTP handler
    // 1. Parse request body
    // 2. Call service
    // 3. Return response
    var input domain.CreateUserInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    user, err := h.service.CreateUser(r.Context(), input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
    // TODO: Implement
}
```

Then create a test file `internal/implementations/handler/user_handler_test.go`:

```go
package handler

import (
    "testing"
    "github.com/kamil5b/lumen-pg/internal/testrunners"
)

func TestUserHandler(t *testing.T) {
    testrunners.UserHandlerRunner(t, NewUserHandler)
}
```

Run the test:
```bash
go test ./internal/implementations/handler -v
```

## Running All Tests

Once all implementations are complete:

```bash
# Run all tests
go test ./... -v

# Run with coverage
go test ./... -cover

# Run integration tests only (requires Docker for testcontainers)
go test ./internal/implementations/repository -v
```

## Key Principles

1. **Write tests first**: Test runners define the contract
2. **Implement to pass tests**: Write minimal code to make tests pass
3. **Refactor**: Improve code while keeping tests green
4. **Isolated testing**: Each layer tested independently
5. **Integration tests**: Repository layer uses real PostgreSQL

## Adding New Features

To add new features following the same pattern:

1. Define domain models in `internal/domain/`
2. Define interfaces in `internal/interfaces/`
3. Generate mocks: `mockgen -source=... -destination=...`
4. Write test runners in `internal/testrunners/`
5. Implement in `internal/implementations/`
6. Write tests that call the test runners

## Regenerating Mocks

If you modify interfaces, regenerate mocks:

```bash
# Repository mock
mockgen -source=internal/interfaces/repository.go \
  -destination=internal/implementations/mocks/user_repository_mock.go \
  -package=mocks

# Service mock
mockgen -source=internal/interfaces/service.go \
  -destination=internal/implementations/mocks/user_service_mock.go \
  -package=mocks
```

## Dependencies

All required dependencies are in `go.mod`:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/stretchr/testify` - Test assertions
- `github.com/testcontainers/testcontainers-go` - Integration tests
- `go.uber.org/mock` - Mock generation
- `github.com/lib/pq` - PostgreSQL driver

## Notes

- Test runners are the "spec" - they define what implementations must do
- Mocks are generated - never edit them manually
- Repository tests use real PostgreSQL via Docker containers
- Service and handler tests use mocks for fast, isolated tests
