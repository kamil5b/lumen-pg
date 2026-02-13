package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/kamil5b/lumen-pg/internal/interfaces/repository"
)

// RBACRepositoryConstructor is a function type that creates an RBACRepository
type RBACRepositoryConstructor func(db *sql.DB) repository.RBACRepository

// RBACRepositoryRunner runs all RBAC repository tests against an implementation
func RBACRepositoryRunner(t *testing.T, constructor RBACRepositoryConstructor) {
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

	// Create test role
	_, err = db.ExecContext(ctx, `CREATE ROLE test_role LOGIN PASSWORD 'testpass'`)
	require.NoError(t, err)

	// Grant permissions
	_, err = db.ExecContext(ctx, `GRANT CONNECT ON DATABASE testdb TO test_role`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `GRANT USAGE ON SCHEMA public TO test_role`)
	require.NoError(t, err)

	// Create test table and grant permissions
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_table (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100)
	)`)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `GRANT SELECT ON test_table TO test_role`)
	require.NoError(t, err)

	repo := constructor(db)

	t.Run("GetUserRole returns role for user", func(t *testing.T) {
		role, err := repo.GetUserRole(ctx, "testuser")
		require.NoError(t, err)
		require.NotEmpty(t, role)
	})

	t.Run("GetUserRole returns error for non-existent user", func(t *testing.T) {
		_, err := repo.GetUserRole(ctx, "nonexistent_user_xyz")
		require.Error(t, err)
	})

	t.Run("GetAllRoles returns list of roles", func(t *testing.T) {
		roles, err := repo.GetAllRoles(ctx)
		require.NoError(t, err)
		require.NotNil(t, roles)
		require.Greater(t, len(roles), 0)
	})

	t.Run("GetRolePermissions returns permission set for role", func(t *testing.T) {
		perms, err := repo.GetRolePermissions(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.NotNil(t, perms)
	})

	t.Run("GetRolePermissions for non-existent role returns error", func(t *testing.T) {
		_, err := repo.GetRolePermissions(ctx, "nonexistent_role", "testdb", "public", "test_table")
		require.Error(t, err)
	})

	t.Run("HasSelectPermission returns true for granted permission", func(t *testing.T) {
		has, err := repo.HasSelectPermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("HasSelectPermission returns false for non-existent table", func(t *testing.T) {
		has, err := repo.HasSelectPermission(ctx, "test_role", "testdb", "public", "nonexistent_table")
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("HasInsertPermission returns correct value", func(t *testing.T) {
		has, err := repo.HasInsertPermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.Equal(t, false, has)
	})

	t.Run("HasUpdatePermission returns correct value", func(t *testing.T) {
		has, err := repo.HasUpdatePermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.Equal(t, false, has)
	})

	t.Run("HasDeletePermission returns correct value", func(t *testing.T) {
		has, err := repo.HasDeletePermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.Equal(t, false, has)
	})

	t.Run("HasDatabaseConnectPermission returns true for granted permission", func(t *testing.T) {
		has, err := repo.HasDatabaseConnectPermission(ctx, "test_role", "testdb")
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("HasDatabaseConnectPermission returns false for non-existent database", func(t *testing.T) {
		has, err := repo.HasDatabaseConnectPermission(ctx, "test_role", "nonexistent_db")
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("HasSchemaUsagePermission returns true for granted permission", func(t *testing.T) {
		has, err := repo.HasSchemaUsagePermission(ctx, "test_role", "testdb", "public")
		require.NoError(t, err)
		require.True(t, has)
	})

	t.Run("HasSchemaUsagePermission returns false for non-existent schema", func(t *testing.T) {
		has, err := repo.HasSchemaUsagePermission(ctx, "test_role", "testdb", "nonexistent_schema")
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("GetAccessibleDatabases returns databases accessible by role", func(t *testing.T) {
		databases, err := repo.GetAccessibleDatabases(ctx, "test_role")
		require.NoError(t, err)
		require.NotNil(t, databases)
		require.Greater(t, len(databases), 0)
	})

	t.Run("GetAccessibleDatabases returns empty for role with no permissions", func(t *testing.T) {
		_, err := repo.GetAccessibleDatabases(ctx, "nonexistent_role")
		require.Error(t, err)
	})

	t.Run("GetAccessibleSchemas returns schemas accessible by role", func(t *testing.T) {
		schemas, err := repo.GetAccessibleSchemas(ctx, "test_role", "testdb")
		require.NoError(t, err)
		require.NotNil(t, schemas)
		require.Greater(t, len(schemas), 0)
	})

	t.Run("GetAccessibleSchemas returns error for non-existent database", func(t *testing.T) {
		_, err := repo.GetAccessibleSchemas(ctx, "test_role", "nonexistent_db")
		require.Error(t, err)
	})

	t.Run("GetAccessibleTables returns tables accessible by role", func(t *testing.T) {
		tables, err := repo.GetAccessibleTables(ctx, "test_role", "testdb", "public")
		require.NoError(t, err)
		require.NotNil(t, tables)
	})

	t.Run("GetAccessibleTables returns error for non-existent schema", func(t *testing.T) {
		_, err := repo.GetAccessibleTables(ctx, "test_role", "testdb", "nonexistent_schema")
		require.Error(t, err)
	})

	t.Run("CanAccessTable returns true for accessible table", func(t *testing.T) {
		can, err := repo.CanAccessTable(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.True(t, can)
	})

	t.Run("CanAccessTable returns false for non-existent table", func(t *testing.T) {
		can, err := repo.CanAccessTable(ctx, "test_role", "testdb", "public", "nonexistent_table")
		require.NoError(t, err)
		require.False(t, can)
	})

	t.Run("GetRoleMetadata returns complete role metadata", func(t *testing.T) {
		metadata, err := repo.GetRoleMetadata(ctx, "test_role")
		require.NoError(t, err)
		require.NotNil(t, metadata)
		require.Equal(t, "test_role", metadata.Name)
	})

	t.Run("GetRoleMetadata returns error for non-existent role", func(t *testing.T) {
		_, err := repo.GetRoleMetadata(ctx, "nonexistent_role")
		require.Error(t, err)
	})

	t.Run("IsReadOnlyRole returns true for SELECT-only role", func(t *testing.T) {
		isReadOnly, err := repo.IsReadOnlyRole(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.True(t, isReadOnly)
	})

	t.Run("IsReadOnlyRole returns false for non-existent table", func(t *testing.T) {
		isReadOnly, err := repo.IsReadOnlyRole(ctx, "test_role", "testdb", "public", "nonexistent_table")
		require.NoError(t, err)
		require.False(t, isReadOnly)
	})

	t.Run("ValidateUserAccessToResource returns true for accessible resource", func(t *testing.T) {
		can, err := repo.ValidateUserAccessToResource(ctx, "test_role", "table", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.True(t, can)
	})

	t.Run("ValidateUserAccessToResource returns false for inaccessible resource", func(t *testing.T) {
		can, err := repo.ValidateUserAccessToResource(ctx, "test_role", "table", "testdb", "public", "nonexistent_table")
		require.NoError(t, err)
		require.False(t, can)
	})

	t.Run("Multiple permission checks for same role", func(t *testing.T) {
		select1, err1 := repo.HasSelectPermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err1)

		select2, err2 := repo.HasSelectPermission(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err2)

		require.Equal(t, select1, select2)
	})

	t.Run("Role accessibility changes reflect correctly", func(t *testing.T) {
		// Check initial access
		can1, err := repo.CanAccessTable(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.True(t, can1)

		// Should still be accessible
		can2, err := repo.CanAccessTable(ctx, "test_role", "testdb", "public", "test_table")
		require.NoError(t, err)
		require.Equal(t, can1, can2)
	})
}
