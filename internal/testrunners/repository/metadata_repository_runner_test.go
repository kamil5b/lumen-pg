package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	_ "github.com/lib/pq"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// MetadataRepoConstructor creates a metadata repository with database connection
type MetadataRepoConstructor func(db *sql.DB) repository.MetadataRepository

// MetadataRepositoryRunner runs integration tests for metadata repository (Story 1)
func MetadataRepositoryRunner(t *testing.T, constructor MetadataRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithInitScripts("testdata/init.sql"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	repo := constructor(db)

	t.Run("IT-S1-01: Connect to Real PostgreSQL", func(t *testing.T) {
		err := db.Ping()
		require.NoError(t, err)
	})

	t.Run("IT-S1-02: Load Real Database Metadata", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx)
		
		require.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.NotEmpty(t, metadata.Databases)
		
		// Should have testdb and postgres databases
		var hasTestDB bool
		for _, db := range metadata.Databases {
			if db.Name == "testdb" {
				hasTestDB = true
				assert.NotEmpty(t, db.Schemas)
			}
		}
		assert.True(t, hasTestDB, "testdb should be present in metadata")
	})

	t.Run("IT-S1-03: Load Real Relations and Role Access", func(t *testing.T) {
		// This test would verify foreign key relationships
		dbMetadata, err := repo.LoadDatabaseMetadata(ctx, "testdb")
		
		require.NoError(t, err)
		assert.NotNil(t, dbMetadata)
		
		// Find tables with foreign keys
		for _, schema := range dbMetadata.Schemas {
			for _, table := range schema.Tables {
				if len(table.ForeignKeys) > 0 {
					assert.NotEmpty(t, table.ForeignKeys[0].ReferencedTableName)
					assert.NotEmpty(t, table.ForeignKeys[0].ReferencedColumnName)
				}
			}
		}
	})

	t.Run("IT-S1-04: Cache Accessible Resources Per Role", func(t *testing.T) {
		roles, err := repo.LoadRoles(ctx)
		
		require.NoError(t, err)
		assert.NotEmpty(t, roles)
		assert.Contains(t, roles, "postgres") // Default postgres role should exist
		
		// Load permissions for postgres role
		roleMetadata, err := repo.LoadRolePermissions(ctx, "postgres")
		require.NoError(t, err)
		assert.NotNil(t, roleMetadata)
		assert.NotEmpty(t, roleMetadata.AccessibleDatabases)
	})
}
