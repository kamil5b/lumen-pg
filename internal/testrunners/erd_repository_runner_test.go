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

// ERDRepoConstructor creates an ERD repository with database connection
type ERDRepoConstructor func(db *sql.DB) repository.MetadataRepository

// ERDRepositoryRunner runs integration tests for ERD repository (Story 3)
func ERDRepositoryRunner(t *testing.T, constructor ERDRepoConstructor) {
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

	// Create test schema with tables and relationships
	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			title VARCHAR(200) NOT NULL,
			content TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id),
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS likes (
			id SERIAL PRIMARY KEY,
			post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(post_id, user_id)
		)
	`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("IT-S3-01: ERD from Real Schema", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.NotEmpty(t, metadata.Databases)

		// Find our testdb
		var testdbFound bool
		for _, dbMeta := range metadata.Databases {
			if dbMeta.Name == "testdb" {
				testdbFound = true
				assert.NotEmpty(t, dbMeta.Schemas)

				// Check for tables
				for _, schema := range dbMeta.Schemas {
					if schema.Name == "public" {
						assert.GreaterOrEqual(t, len(schema.Tables), 4) // users, posts, comments, likes
						break
					}
				}
				break
			}
		}
		assert.True(t, testdbFound, "testdb should be present in metadata")
	})

	t.Run("IT-S3-02: Complex Relationships", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx)

		require.NoError(t, err)
		assert.NotNil(t, metadata)

		// Find testdb and verify relationships
		for _, db := range metadata.Databases {
			if db.Name == "testdb" {
				for _, schema := range db.Schemas {
					if schema.Name == "public" {
						// Find posts table and verify foreign key to users
						for _, table := range schema.Tables {
							if table.TableName == "posts" {
								assert.NotEmpty(t, table.ForeignKeys)
								assert.Equal(t, "users", table.ForeignKeys[0].ReferencedTableName)
								assert.Equal(t, "user_id", table.ForeignKeys[0].ColumnName)
							}

							// Find comments and verify multiple foreign keys
							if table.TableName == "comments" {
								assert.GreaterOrEqual(t, len(table.ForeignKeys), 2)
								var hasPostFK, hasUserFK bool
								for _, fk := range table.ForeignKeys {
									if fk.ReferencedTableName == "posts" {
										hasPostFK = true
									}
									if fk.ReferencedTableName == "users" {
										hasUserFK = true
									}
								}
								assert.True(t, hasPostFK, "comments should have foreign key to posts")
								assert.True(t, hasUserFK, "comments should have foreign key to users")
							}

							// Find likes and verify unique constraint and foreign keys
							if table.TableName == "likes" {
								assert.GreaterOrEqual(t, len(table.ForeignKeys), 2)
							}
						}
					}
				}
			}
		}
	})
}
