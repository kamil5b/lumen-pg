# Lumen-PG - TDD Structure

This repository implements a TDD (Test-Driven Development) structure for a Go application following best practices.

## Project Structure

```
myapp/
├─ internal/
│  ├─ domain/              # Domain models
│  │  └─ user.go
│  ├─ interfaces/          # Interface definitions
│  │  ├─ repository.go
│  │  ├─ service.go
│  │  └─ handler.go
│  ├─ implementations/     # Concrete implementations
│  │  ├─ mocks/            # Auto-generated mocks
│  │  │  ├─ user_repository_mock.go
│  │  │  └─ user_service_mock.go
│  │  ├─ repository/
│  │  ├─ service/
│  │  └─ handler/
│  └─ testrunners/         # Test runner specs
│     ├─ repository_runner_test.go
│     ├─ service_runner_test.go
│     └─ handler_runner_test.go
```

## TDD Flow

### Step 1 - Domain Models
Domain models are defined first in `internal/domain/`. These are plain Go structs representing business entities.

### Step 2 - Interfaces
Interfaces are defined in `internal/interfaces/` before implementations. This allows for dependency injection and testing with mocks.

### Step 3 - Generate Mocks
Before writing tests, generate mocks using `mockgen`:

```bash
# Generate repository mock
mockgen -source=internal/interfaces/repository.go \
  -destination=internal/implementations/mocks/user_repository_mock.go \
  -package=mocks

# Generate service mock
mockgen -source=internal/interfaces/service.go \
  -destination=internal/implementations/mocks/user_service_mock.go \
  -package=mocks
```

### Step 4 - Test Runners
Test runners in `internal/testrunners/` define the expected behavior:

- **Repository Runner**: Uses testcontainers for real PostgreSQL testing
- **Service Runner**: Uses gomock for mocked repository testing
- **Handler Runner**: Uses gomock for mocked service testing

### Step 5 - Implementations
After test runners are ready, implement the actual functionality:

```go
// Example usage in implementation tests
func TestPostgresRepository(t *testing.T) {
    UserRepositoryRunner(t, implementations.NewPostgresUserRepository)
}

func TestUserService(t *testing.T) {
    UserServiceRunner(t, implementations.NewUserService)
}

func TestUserHandler(t *testing.T) {
    UserHandlerRunner(t, implementations.NewUserHandler)
}
```

## Running Tests

```bash
# Build all packages
go build ./...

# Compile test runners
go test -c ./internal/testrunners

# When implementations are ready, run tests
go test ./...
```

## Dependencies

- **github.com/go-chi/chi/v5** - HTTP router
- **github.com/stretchr/testify** - Testing assertions
- **github.com/testcontainers/testcontainers-go** - Container-based integration tests
- **go.uber.org/mock** - Mock generation
- **github.com/lib/pq** - PostgreSQL driver

## Notes

- Mocks are auto-generated - do not edit them manually
- Test runners define the contract that implementations must satisfy
- Integration tests use real PostgreSQL via testcontainers
- Unit tests use mocks for isolated testing
