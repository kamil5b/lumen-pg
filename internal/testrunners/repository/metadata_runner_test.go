package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// MetadataRepositoryConstructor is a function type that creates a MetadataRepository
type MetadataRepositoryConstructor func(db *sql.DB) repository.MetadataRepository

// MetadataRepositoryRunner runs all metadata repository tests against an implementation
func MetadataRepositoryRunner(t *testing.T, constructor MetadataRepositoryConstructor) {
	t.Helper()

	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer db.Close()

	err = db.PingContext(ctx)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("StoreMetadata and GetMetadata", func(t *testing.T) {
		metadata := &domain.DatabaseMetadata{
			Name: "testdb",
			Schemas: []domain.SchemaMetadata{
				{
					Name: "public",
					Tables: []domain.TableMetadata{
						{
							Name: "users",
							Columns: []domain.ColumnMetadata{
								{
									Name:       "id",
									DataType:   "integer",
									IsNullable: false,
									IsPrimary:  true,
								},
								{
									Name:       "name",
									DataType:   "text",
									IsNullable: true,
									IsPrimary:  false,
								},
							},
							PrimaryKeys: []string{"id"},
							ForeignKeys: []domain.ForeignKeyMetadata{},
						},
					},
				},
			},
		}

		err := repo.StoreMetadata(ctx, metadata)
		require.NoError(t, err)

		retrieved, err := repo.GetMetadata(ctx, "testdb")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, metadata.Name, retrieved.Name)
		require.Len(t, retrieved.Schemas, 1)
		require.Equal(t, "public", retrieved.Schemas[0].Name)
	})

	t.Run("GetMetadata returns error for non-existent database", func(t *testing.T) {
		_, err := repo.GetMetadata(ctx, "non_existent_db")
		require.Error(t, err)
	})

	t.Run("StoreRoleMetadata and GetRoleMetadata", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "test_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public"},
			AccessibleTables: []domain.AccessibleTable{
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "users",
					HasSelect: true,
					HasInsert: true,
					HasUpdate: true,
					HasDelete: false,
				},
			},
		}

		err := repo.StoreRoleMetadata(ctx, "test_role", roleMetadata)
		require.NoError(t, err)

		retrieved, err := repo.GetRoleMetadata(ctx, "test_role")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.Equal(t, "test_role", retrieved.Name)
		require.Len(t, retrieved.AccessibleDatabases, 1)
		require.Len(t, retrieved.AccessibleTables, 1)
	})

	t.Run("GetRoleMetadata returns error for non-existent role", func(t *testing.T) {
		_, err := repo.GetRoleMetadata(ctx, "non_existent_role")
		require.Error(t, err)
	})

	t.Run("StoreAllRolesMetadata and GetAllRolesMetadata", func(t *testing.T) {
		roles := map[string]*domain.RoleMetadata{
			"role1": {
				Name:                "role1",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables:    []domain.AccessibleTable{},
			},
			"role2": {
				Name:                "role2",
				AccessibleDatabases: []string{"testdb"},
				AccessibleSchemas:   []string{"public"},
				AccessibleTables:    []domain.AccessibleTable{},
			},
		}

		err := repo.StoreAllRolesMetadata(ctx, roles)
		require.NoError(t, err)

		retrieved, err := repo.GetAllRolesMetadata(ctx)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		require.GreaterOrEqual(t, len(retrieved), 2)
	})

	t.Run("InvalidateMetadata", func(t *testing.T) {
		metadata := &domain.DatabaseMetadata{
			Name:    "testdb_invalidate",
			Schemas: []domain.SchemaMetadata{},
		}

		err := repo.StoreMetadata(ctx, metadata)
		require.NoError(t, err)

		err = repo.InvalidateMetadata(ctx, "testdb_invalidate")
		require.NoError(t, err)

		_, err = repo.GetMetadata(ctx, "testdb_invalidate")
		require.Error(t, err)
	})

	t.Run("InvalidateRoleMetadata", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "role_to_invalidate",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{},
			AccessibleTables:    []domain.AccessibleTable{},
		}

		err := repo.StoreRoleMetadata(ctx, "role_to_invalidate", roleMetadata)
		require.NoError(t, err)

		err = repo.InvalidateRoleMetadata(ctx, "role_to_invalidate")
		require.NoError(t, err)

		_, err = repo.GetRoleMetadata(ctx, "role_to_invalidate")
		require.Error(t, err)
	})

	t.Run("InvalidateAllMetadata", func(t *testing.T) {
		metadata := &domain.DatabaseMetadata{
			Name:    "testdb_all_invalidate",
			Schemas: []domain.SchemaMetadata{},
		}

		err := repo.StoreMetadata(ctx, metadata)
		require.NoError(t, err)

		err = repo.InvalidateAllMetadata(ctx)
		require.NoError(t, err)

		_, err = repo.GetMetadata(ctx, "testdb_all_invalidate")
		require.Error(t, err)
	})

	t.Run("GetAccessibleDatabases", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "db_access_role",
			AccessibleDatabases: []string{"db1", "db2", "db3"},
			AccessibleSchemas:   []string{},
			AccessibleTables:    []domain.AccessibleTable{},
		}

		err := repo.StoreRoleMetadata(ctx, "db_access_role", roleMetadata)
		require.NoError(t, err)

		databases, err := repo.GetAccessibleDatabases(ctx, "db_access_role")
		require.NoError(t, err)
		require.NotNil(t, databases)
		require.GreaterOrEqual(t, len(databases), 3)
	})

	t.Run("GetAccessibleSchemas", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "schema_access_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public", "private"},
			AccessibleTables:    []domain.AccessibleTable{},
		}

		err := repo.StoreRoleMetadata(ctx, "schema_access_role", roleMetadata)
		require.NoError(t, err)

		schemas, err := repo.GetAccessibleSchemas(ctx, "schema_access_role", "testdb")
		require.NoError(t, err)
		require.NotNil(t, schemas)
		require.GreaterOrEqual(t, len(schemas), 2)
	})

	t.Run("GetAccessibleTables", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "table_access_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public"},
			AccessibleTables: []domain.AccessibleTable{
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "users",
					HasSelect: true,
					HasInsert: true,
					HasUpdate: true,
					HasDelete: false,
				},
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "posts",
					HasSelect: true,
					HasInsert: false,
					HasUpdate: false,
					HasDelete: false,
				},
			},
		}

		err := repo.StoreRoleMetadata(ctx, "table_access_role", roleMetadata)
		require.NoError(t, err)

		tables, err := repo.GetAccessibleTables(ctx, "table_access_role", "testdb", "public")
		require.NoError(t, err)
		require.NotNil(t, tables)
		require.GreaterOrEqual(t, len(tables), 2)
	})

	t.Run("IsTableAccessible returns true for accessible table", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "table_check_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public"},
			AccessibleTables: []domain.AccessibleTable{
				{
					Database:  "testdb",
					Schema:    "public",
					Name:      "accessible_table",
					HasSelect: true,
					HasInsert: true,
					HasUpdate: true,
					HasDelete: false,
				},
			},
		}

		err := repo.StoreRoleMetadata(ctx, "table_check_role", roleMetadata)
		require.NoError(t, err)

		accessible, err := repo.IsTableAccessible(ctx, "table_check_role", "testdb", "public", "accessible_table")
		require.NoError(t, err)
		require.True(t, accessible)
	})

	t.Run("IsTableAccessible returns false for inaccessible table", func(t *testing.T) {
		roleMetadata := &domain.RoleMetadata{
			Name:                "no_access_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public"},
			AccessibleTables:    []domain.AccessibleTable{},
		}

		err := repo.StoreRoleMetadata(ctx, "no_access_role", roleMetadata)
		require.NoError(t, err)

		accessible, err := repo.IsTableAccessible(ctx, "no_access_role", "testdb", "public", "inaccessible_table")
		require.NoError(t, err)
		require.False(t, accessible)
	})

	t.Run("GetTablePermissions", func(t *testing.T) {
		expectedPerms := &domain.AccessibleTable{
			Database:  "testdb",
			Schema:    "public",
			Name:      "perm_table",
			HasSelect: true,
			HasInsert: true,
			HasUpdate: false,
			HasDelete: false,
		}

		roleMetadata := &domain.RoleMetadata{
			Name:                "perm_role",
			AccessibleDatabases: []string{"testdb"},
			AccessibleSchemas:   []string{"public"},
			AccessibleTables: []domain.AccessibleTable{
				*expectedPerms,
			},
		}

		err := repo.StoreRoleMetadata(ctx, "perm_role", roleMetadata)
		require.NoError(t, err)

		perms, err := repo.GetTablePermissions(ctx, "perm_role", "testdb", "public", "perm_table")
		require.NoError(t, err)
		require.NotNil(t, perms)
		require.Equal(t, expectedPerms.HasSelect, perms.HasSelect)
		require.Equal(t, expectedPerms.HasInsert, perms.HasInsert)
		require.Equal(t, expectedPerms.HasUpdate, perms.HasUpdate)
		require.Equal(t, expectedPerms.HasDelete, perms.HasDelete)
	})

	t.Run("GetTablePermissions returns error for non-existent table", func(t *testing.T) {
		_, err := repo.GetTablePermissions(ctx, "nonexistent_role", "testdb", "public", "nonexistent_table")
		require.Error(t, err)
	})
}
