package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupPostgresContainer creates a real PostgreSQL container for integration testing.
func setupPostgresContainer(t *testing.T) (string, func()) {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb1"),
		postgres.WithUsername("superuser"),
		postgres.WithPassword("superpassword"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	cleanup := func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return connStr, cleanup
}

// === Story 1: Setup & Configuration - Integration Tests ===

// IT-S1-01: Connect to Real PostgreSQL
func TestIntegration_ConnectToRealPostgreSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	connStr, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubConnectionRepository()
	err := repo.TestConnection(context.Background(), connStr)
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should connect successfully")
}

// IT-S1-02: Load Real Database Metadata with User Accessible Resources
func TestIntegration_LoadRealDatabaseMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	metadata, err := repo.LoadAllMetadata(context.Background())
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should load metadata")
	assert.Nil(t, metadata)
}

// IT-S1-03: Load Real Relations and Role Access
func TestIntegration_LoadRealRelationsAndRoleAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	fks, err := repo.LoadForeignKeys(context.Background(), "testdb1", "public", "posts")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should load foreign keys")
	assert.Nil(t, fks)
}

// IT-S1-04: Cache Accessible Resources Per Role
func TestIntegration_CacheAccessibleResourcesPerRole(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	roles, err := repo.LoadRoles(context.Background())
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should load roles")
	assert.Nil(t, roles)
}

// === Story 2: Authentication - Integration Tests ===

// IT-S2-01: Real PostgreSQL Connection Probe
func TestIntegration_RealConnectionProbe(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubConnectionRepository()
	session, err := repo.ProbeConnection(context.Background(), "superuser", "superpassword", "localhost", 5432, "disable", []string{"testdb1"})
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should probe connection")
	assert.Nil(t, session)
}

// IT-S2-02: Real PostgreSQL Connection Probe Failure
func TestIntegration_RealConnectionProbeFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubConnectionRepository()
	session, err := repo.ProbeConnection(context.Background(), "noaccess_user", "password", "localhost", 5432, "disable", []string{})
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should fail probe")
	assert.Nil(t, session)
}

// IT-S2-03: Real Role-Based Resource Access
func TestIntegration_RealRoleBasedResourceAccess(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()

	for _, roleName := range []string{"admin", "editor", "viewer"} {
		role, err := repo.GetAccessibleResources(context.Background(), roleName)
		assert.ErrorIs(t, err, domain.ErrNotImplemented, fmt.Sprintf("stub should return ErrNotImplemented for role %s", roleName))
		assert.Nil(t, role)
	}
}

// === Story 3: ERD Viewer - Integration Tests ===

// IT-S3-01: ERD from Real Schema
func TestIntegration_ERDFromRealSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	erd, err := repo.GenerateERDData(context.Background(), "testdb1", "public")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should generate ERD")
	assert.Nil(t, erd)
}

// IT-S3-02: Complex Relationships
func TestIntegration_ComplexRelationships(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	erd, err := repo.GenerateERDData(context.Background(), "testdb1", "public")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should show all relationships")
	assert.Nil(t, erd)
}

// === Story 4: Manual Query Editor - Integration Tests ===

// IT-S4-01: Real SELECT Query
func TestIntegration_RealSelectQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	result, err := repo.ExecuteQuery(context.Background(), "testdb1", "SELECT * FROM users")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should execute query")
	assert.Nil(t, result)
}

// IT-S4-02: Real DDL Query
func TestIntegration_RealDDLQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	result, err := repo.ExecuteQuery(context.Background(), "testdb1", "CREATE TABLE test_table (id SERIAL PRIMARY KEY)")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should create table")
	assert.Nil(t, result)
}

// IT-S4-03: Real DML Query
func TestIntegration_RealDMLQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	result, err := repo.ExecuteQuery(context.Background(), "testdb1", "INSERT INTO users (username, email) VALUES ('test', 'test@example.com')")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should insert data")
	assert.Nil(t, result)
}

// IT-S4-04: Query with Permission Denied
func TestIntegration_QueryWithPermissionDenied(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	result, err := repo.ExecuteQuery(context.Background(), "testdb1", "SELECT * FROM restricted_table")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should return permission denied")
	assert.Nil(t, result)
}

// === Story 5: Main View & Data Interaction - Integration Tests ===

// IT-S5-01: Real Table Data Loading
func TestIntegration_RealTableDataLoading(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	cursor := domain.NewCursor()
	page, err := repo.LoadTableData(context.Background(), "testdb1", "public", "users", cursor, "id", "ASC", "")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should load data")
	assert.Nil(t, page)
}

// IT-S5-02: Real Cursor Pagination
func TestIntegration_RealCursorPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	cursor := domain.NewCursor()
	page, err := repo.LoadTableData(context.Background(), "testdb1", "public", "users", cursor, "id", "ASC", "")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should paginate")
	assert.Nil(t, page)
}

// IT-S5-03: Real WHERE Filter
func TestIntegration_RealWhereFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	cursor := domain.NewCursor()
	page, err := repo.LoadTableData(context.Background(), "testdb1", "public", "users", cursor, "id", "ASC", "id > 5")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should filter")
	assert.Nil(t, page)
}

// IT-S5-04: Real Transaction Commit
func TestIntegration_RealTransactionCommit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	ops := []domain.BufferedOperation{
		{
			Type:    domain.OpInsert,
			Table:   "users",
			Schema:  "public",
			RowData: map[string]interface{}{"username": "commituser", "email": "commit@test.com"},
		},
	}
	err := repo.ExecuteTransaction(context.Background(), "testdb1", ops)
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should commit")
}

// IT-S5-05: Real Transaction Rollback
func TestIntegration_RealTransactionRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	ops := []domain.BufferedOperation{
		{
			Type:    domain.OpInsert,
			Table:   "users",
			Schema:  "public",
			RowData: map[string]interface{}{"username": "rollbackuser"},
		},
	}
	err := repo.ExecuteTransaction(context.Background(), "testdb1", ops)
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should handle rollback scenario")
}

// IT-S5-06: Real Foreign Key Navigation
func TestIntegration_RealForeignKeyNavigation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	cursor := domain.NewCursor()
	page, err := repo.LoadTableData(context.Background(), "testdb1", "public", "posts", cursor, "id", "ASC", "user_id = 1")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should load FK data")
	assert.Nil(t, page)
}

// IT-S5-07: Real Primary Key Navigation
func TestIntegration_RealPrimaryKeyNavigation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	tables, err := repo.GetReferencingTables(context.Background(), "testdb1", "public", "users", "id", 1)
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should find referencing tables")
	assert.Nil(t, tables)
}

// === Story 6: Isolation - Integration Tests ===

// IT-S6-01: Real Multi-User Connection
func TestIntegration_RealMultiUserConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubConnectionRepository()
	// User A
	sessionA, err := repo.ProbeConnection(context.Background(), "user_readonly", "password", "localhost", 5432, "disable", []string{"testdb1"})
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, sessionA)

	// User B
	sessionB, err := repo.ProbeConnection(context.Background(), "user_limited", "password", "localhost", 5432, "disable", []string{"testdb1"})
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, sessionB)
}

// IT-S6-02: Real Permission Isolation
func TestIntegration_RealPermissionIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubMetadataRepository()
	// User A has access
	roleA, err := repo.GetAccessibleResources(context.Background(), "user_readonly")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, roleA)

	// User B has limited access
	roleB, err := repo.GetAccessibleResources(context.Background(), "user_limited")
	assert.ErrorIs(t, err, domain.ErrNotImplemented)
	assert.Nil(t, roleB)
}

// IT-S6-03: Real Transaction Isolation
func TestIntegration_RealTransactionIsolation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	// User A starts a transaction
	opsA := []domain.BufferedOperation{
		{Type: domain.OpUpdate, Table: "users", Schema: "public"},
	}
	err := repo.ExecuteTransaction(context.Background(), "testdb1", opsA)
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should isolate transactions")
}

// === Story 7: Security - Integration Tests ===

// IT-S7-01: Real SQL Injection Test
func TestIntegration_RealSQLInjection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	_, cleanup := setupPostgresContainer(t)
	defer cleanup()

	repo := repository.NewStubQueryRepository()
	// Attempt SQL injection through query
	result, err := repo.ExecuteQuery(context.Background(), "testdb1", "SELECT * FROM users WHERE id = 1; DROP TABLE users;--")
	assert.ErrorIs(t, err, domain.ErrNotImplemented, "stub should return ErrNotImplemented; real implementation should prevent injection")
	assert.Nil(t, result)
}
