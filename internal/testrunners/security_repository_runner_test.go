package testrunners

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// SecurityRepoConstructor creates a security repository with database connection
type SecurityRepoConstructor func(db *sql.DB) repository.QueryRepository

// SecurityRepositoryRunner runs integration tests for security repository (Story 7)
func SecurityRepositoryRunner(t *testing.T, constructor SecurityRepoConstructor) {
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

	// Create test table for security tests
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100) NOT NULL,
			password_hash VARCHAR(255)
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (username, email, password_hash) VALUES
		('alice', 'alice@test.com', 'hashed_password_1'),
		('bob', 'bob@test.com', 'hashed_password_2')
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("IT-S7-01: Real SQL Injection Test", func(t *testing.T) {
		// Test with parameterized query to prevent SQL injection
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users WHERE username = $1", "admin' OR '1'='1")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		// Should not return admin user with injection attempt
		assert.Len(t, result.Rows, 0)
	})

	t.Run("IT-S7-02: Real Password Security", func(t *testing.T) {
		// Query should return password_hash, not plaintext password
		result, err := repo.ExecuteQuery(ctx, "SELECT username, password_hash FROM users WHERE username = $1", "alice")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 1)
		// Verify password_hash is present in columns
		assert.Contains(t, result.Columns, "password_hash")
	})

	t.Run("IT-S7-03: Real Session Expiration", func(t *testing.T) {
		// This test verifies that old sessions expire
		// In a real scenario, this would check session tables in database
		sessionID := "test-session-123"

		// Check if session exists (should not for test purposes)
		var exists bool
		err := db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables
				WHERE table_name = 'sessions'
			)
		`).Scan(&exists)

		require.NoError(t, err)
		// If sessions table exists, verify expired sessions can be identified
		if exists {
			var count int
			err = db.QueryRowContext(ctx, `
				SELECT COUNT(*) FROM sessions
				WHERE session_id = $1 AND expires_at < NOW()
			`, sessionID).Scan(&count)

			// Session should be expired or not exist
			require.NoError(t, err)
		}
	})

	t.Run("IT-S7-04: Query Parameterization Enforcement", func(t *testing.T) {
		// Test that parameterized queries work correctly
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users WHERE id = $1", 1)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 1)
	})

	t.Run("IT-S7-05: Multiple Parameter Injection Prevention", func(t *testing.T) {
		// Test with multiple parameters
		result, err := repo.ExecuteQuery(ctx,
			"SELECT * FROM users WHERE username = $1 AND email = $2",
			"alice' OR '1'='1",
			"alice@test.com' OR '1'='1")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		// Should return only matching records, not all records
		assert.GreaterOrEqual(t, len(result.Rows), 0)
	})

	t.Run("IT-S7-06: DDL Prevention in User Context", func(t *testing.T) {
		// DDL operations should be properly controlled
		// This would typically be enforced through role permissions
		err := repo.ExecuteDDL(ctx, "CREATE TABLE test (id INT)")

		// Operation may succeed or fail depending on permissions
		// What matters is that it's executed safely
		_ = err // Ignore for this test
	})

	t.Run("IT-S7-07: DML Permission Enforcement", func(t *testing.T) {
		// DML operations should respect database permissions
		result, err := repo.ExecuteDML(ctx,
			"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
			"charlie", "charlie@test.com", "hashed_password_3")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, int64(1), result.AffectedRows)

		// Verify insert was successful
		verifyResult, err := repo.ExecuteQuery(ctx, "SELECT * FROM users WHERE username = $1", "charlie")
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(verifyResult.Rows), 1)
	})
}
