**strict TDD order**

1. Define domain
2. Define interfaces
3. **Generate mocks** (with `gomock`)
4. Create test runners using those mocks/testcontainers
5. Implement later

Here’s the step-by-step guide following this flow.

---

# Go Project – TDD Step-by-Step with Mock Generation First

## Project Structure

```
myapp/
├─ internal/
│  ├─ domain/
│  │  └─ user.go
│  ├─ interfaces/
│  │  ├─ repository.go
│  │  ├─ service.go
│  │  └─ handler.go
│  ├─ implementations/
│  │  ├─ repository/
│  │  ├─ service/
│  │  └─ handler/
│  └─ testrunners/
│     ├─ repository_runner_test.go
│     ├─ service_runner_test.go
│     └─ handler_runner_test.go
```

---

# Step 1 — Define Domain

**`internal/domain/user.go`**

```go
package domain

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type CreateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}
```

---

# Step 2 — Define Interfaces

**Repository Interface (`interfaces/repository.go`)**

```go
package interfaces

import (
	"context"
	"myapp/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
}
```

**Service Interface (`interfaces/service.go`)**

```go
package interfaces

import (
	"context"
	"myapp/internal/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
}
```

**Handler Interface (`interfaces/handler.go`)**

```go
package interfaces

import "github.com/go-chi/chi/v5"

type UserHandler interface {
	RegisterRoutes(r chi.Router)
}
```

---

# Step 3 — Generate Mocks (Before Writing Test Runners)

Install gomock if you haven’t:

```bash
go install github.com/golang/mock/mockgen@latest
```

---

## 3A — Generate Repository Mock

```bash
mockgen -source=internal/interfaces/repository.go \
  -destination=internal/implementations/mocks/user_repository_mock.go \
  -package=mocks
```

---

## 3B — Generate Service Mock

```bash
mockgen -source=internal/interfaces/service.go \
  -destination=internal/implementations/mocks/user_service_mock.go \
  -package=mocks
```

> **Note:** Handler doesn’t need a mock unless you’re testing higher layers depending on handlers. Usually only Service is mocked for handler testing.

---

# Step 4 — Create Test Runners (Using Mocks/Testcontainers)

We now write the **specs/test runners** using the mocks generated.

---

## 4A — Repository Runner (Testcontainers)

```go
package testrunners

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	_ "github.com/lib/pq"

	"myapp/internal/domain"
	"myapp/internal/interfaces"
)

type RepoConstructor func(db *sql.DB) interfaces.UserRepository

func UserRepositoryRunner(t *testing.T, constructor RepoConstructor) {
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

	t.Run("Save and FindByID roundtrip", func(t *testing.T) {
		user := &domain.User{
			ID:        "123",
			Email:     "repo@test.com",
			Name:      "Repo",
			CreatedAt: time.Now(),
		}

		err := repo.Save(ctx, user)
		require.NoError(t, err)

		found, err := repo.FindByID(ctx, "123")
		require.NoError(t, err)

		assert.Equal(t, user.Email, found.Email)
	})
}
```

---

## 4B — Service Runner (Gomock)

```go
package testrunners

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"myapp/internal/interfaces"
	repomock "myapp/internal/implementations/mocks"
)

type ServiceConstructor func(repo interfaces.UserRepository) interfaces.UserService

func UserServiceRunner(t *testing.T, constructor ServiceConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repomock.NewMockUserRepository(ctrl)

	svc := constructor(mockRepo)

	t.Run("CreateUser generates ID and saves", func(t *testing.T) {
		input := interfaces.CreateUserInput{
			Email: "service@test.com",
			Name:  "Service",
		}

		mockRepo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		user, err := svc.CreateUser(context.Background(), input)

		require.NoError(t, err)
		assert.NotEmpty(t, user.ID)
		assert.Equal(t, input.Email, user.Email)
	})
}
```

---

## 4C — Handler Runner (Gomock Service)

```go
package testrunners

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	svcmock "myapp/internal/implementations/mocks"
	"myapp/internal/interfaces"
)

type HandlerConstructor func(svc interfaces.UserService) interfaces.UserHandler

func UserHandlerRunner(t *testing.T, constructor HandlerConstructor) {
	t.Helper()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := svcmock.NewMockUserService(ctrl)

	h := constructor(mockSvc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	t.Run("POST /users returns 201", func(t *testing.T) {
		input := interfaces.CreateUserInput{
			Email: "http@test.com",
			Name:  "HTTP",
		}

		mockSvc.EXPECT().CreateUser(gomock.Any(), input).Return(&interfaces.User{
			ID:    "123",
			Email: input.Email,
			Name:  input.Name,
		}, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		require.Equal(t, http.StatusCreated, rec.Code)
	})
}
```

---

# Step 5 — Implementation Later

After runners/specs are ready, implement:

```go
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

Now TDD workflow is complete:

1. Define domain
2. Define interfaces
3. **Generate mocks first**
4. Write failing test runners
5. Implement code to make them pass

---

✅ This approach ensures:

* Tests first → fail until implementation exists
* Mocks available for service/handler runners
* Testcontainers for repository integration
* Clean, layered architecture
* Runner pattern standard across all layers
