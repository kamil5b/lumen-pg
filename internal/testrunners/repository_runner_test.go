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

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
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
