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

// DataExplorerRepoConstructor creates a data explorer repository with database connection
type DataExplorerRepoConstructor func(db *sql.DB) repository.QueryRepository

// DataExplorerRepositoryRunner runs integration tests for data explorer repository (Story 5)
func DataExplorerRepositoryRunner(t *testing.T, constructor DataExplorerRepoConstructor) {
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

	// Create test tables with relationships
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100) NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			title VARCHAR(200) NOT NULL,
			content TEXT
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id),
			user_id INTEGER NOT NULL REFERENCES users(id),
			content TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	// Insert test data
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (username, email) VALUES
		('alice', 'alice@test.com'),
		('bob', 'bob@test.com'),
		('charlie', 'charlie@test.com'),
		('diana', 'diana@test.com'),
		('eve', 'eve@test.com')
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO posts (user_id, title, content) VALUES
		(1, 'Post 1', 'Content 1'),
		(1, 'Post 2', 'Content 2'),
		(2, 'Post 3', 'Content 3'),
		(3, 'Post 4', 'Content 4'),
		(1, 'Post 5', 'Content 5')
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO comments (post_id, user_id, content) VALUES
		(1, 2, 'Great post!'),
		(1, 3, 'Thanks for sharing!'),
		(2, 4, 'Interesting'),
		(3, 1, 'Nice'),
		(4, 2, 'Awesome')
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("IT-S5-01: Real Table Data Loading", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users LIMIT 50")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Len(t, result.Columns, 3) // id, username, email
		assert.GreaterOrEqual(t, len(result.Rows), 5)
	})

	t.Run("IT-S5-02: Real Cursor Pagination", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users LIMIT 2 OFFSET 0")

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 2)
		assert.GreaterOrEqual(t, result.TotalRows, int64(5))
	})

	t.Run("IT-S5-03: Real WHERE Filter", func(t *testing.T) {
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users WHERE id > $1 LIMIT 50", 2)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 3)
	})

	t.Run("IT-S5-04: Real Transaction Commit", func(t *testing.T) {
		// Start transaction
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)

		// Insert data within transaction
		_, err = tx.ExecContext(ctx, "INSERT INTO users (username, email) VALUES ($1, $2)", "frank", "frank@test.com")
		require.NoError(t, err)

		// Commit transaction
		err = tx.Commit()
		require.NoError(t, err)

		// Verify insert was committed
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM users WHERE username = $1", "frank")
		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 1)
	})

	t.Run("IT-S5-05: Real Transaction Rollback", func(t *testing.T) {
		countBefore, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM users")
		require.NoError(t, err)
		defer countBefore.Close()

		var before int
		if countBefore.Next() {
			countBefore.Scan(&before)
		}

		// Start transaction
		tx, err := db.BeginTx(ctx, nil)
		require.NoError(t, err)

		// Insert data within transaction
		_, err = tx.ExecContext(ctx, "INSERT INTO users (username, email) VALUES ($1, $2)", "grace", "grace@test.com")
		require.NoError(t, err)

		// Rollback transaction
		err = tx.Rollback()
		require.NoError(t, err)

		// Verify insert was rolled back
		countAfter, err := db.QueryContext(ctx, "SELECT COUNT(*) FROM users")
		require.NoError(t, err)
		defer countAfter.Close()

		var after int
		if countAfter.Next() {
			countAfter.Scan(&after)
		}

		assert.Equal(t, before, after, "Row count should be same after rollback")
	})

	t.Run("IT-S5-06: Real Foreign Key Navigation", func(t *testing.T) {
		// Get posts by user_id = 1
		result, err := repo.ExecuteQuery(ctx, "SELECT * FROM posts WHERE user_id = $1 LIMIT 50", 1)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 3) // alice has 3 posts
	})

	t.Run("IT-S5-07: Real Primary Key Navigation", func(t *testing.T) {
		// Get all posts referencing user id=1
		result, err := repo.ExecuteQuery(ctx, `
			SELECT COUNT(*) as count FROM posts WHERE user_id = $1
		`, 1)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, len(result.Rows), 1)
	})
}
