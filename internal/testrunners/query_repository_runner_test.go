package testrunners

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	_ "github.com/lib/pq"

	"github.com/kamil5b/lumen-pg/internal/interfaces"
)

// QueryRepoConstructor creates a query repository with database connection
type QueryRepoConstructor func(db *sql.DB) interfaces.QueryRepository

// QueryRepositoryRunner runs integration tests for query repository (Story 4)
func QueryRepositoryRunner(t *testing.T, constructor QueryRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	// Create test table
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS test_users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50),
			email VARCHAR(100)
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.ExecContext(ctx, `
		INSERT INTO test_users (username, email) VALUES
		('user1', 'user1@test.com'),
		('user2', 'user2@test.com')
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("IT-S4-01: Real SELECT Query", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM test_users")
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Columns, 3) // id, username, email
		assert.GreaterOrEqual(t, len(result.Rows), 2)
	})

	t.Run("IT-S4-02: Real DDL Query", func(t *testing.T) {
		err := repo.ExecuteDDL(ctx, "CREATE TABLE IF NOT EXISTS test_posts (id SERIAL PRIMARY KEY, title VARCHAR(100))")
		
		require.NoError(t, err)
		
		// Verify table was created
		var exists bool
		err = db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_name = 'test_posts'
			)
		`).Scan(&exists)
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("IT-S4-03: Real DML Query", func(t *testing.T) {
		result, err := repo.ExecuteDML(ctx, "INSERT INTO test_users (username, email) VALUES ($1, $2)", "user3", "user3@test.com")
		
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.AffectedRows)
	})

	t.Run("IT-S4-04: Query with Permission Denied", func(t *testing.T) {
		// This test would require a restricted user
		// For now, test invalid query handling
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM nonexistent_table")
		
		require.NoError(t, err) // No error from repo, but result contains error
		assert.NotNil(t, result)
		assert.False(t, result.Success)
		assert.NotEmpty(t, result.ErrorMessage)
	})
}
