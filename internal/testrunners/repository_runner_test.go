package testrunners

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/kamil5b/lumen-pg/internal/interfaces"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// ConnectionRepoConstructor is a function that creates a ConnectionRepository
type ConnectionRepoConstructor func(superAdminConnStr string) interfaces.ConnectionRepository

// ConnectionRepositoryRunner tests ConnectionRepository implementations
func ConnectionRepositoryRunner(t *testing.T, constructor ConnectionRepoConstructor) {
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

	repo := constructor(connStr)

	t.Run("ValidateConnectionString - valid", func(t *testing.T) {
		err := repo.ValidateConnectionString(connStr)
		assert.NoError(t, err)
	})

	t.Run("ValidateConnectionString - invalid", func(t *testing.T) {
		err := repo.ValidateConnectionString("invalid://connection")
		assert.Error(t, err)
	})

	t.Run("TestConnection - success", func(t *testing.T) {
		err := repo.TestConnection(ctx, connStr)
		assert.NoError(t, err)
	})

	t.Run("TestConnection - failure", func(t *testing.T) {
		badConnStr := "postgres://baduser:badpass@localhost:5432/baddb?sslmode=disable"
		err := repo.TestConnection(ctx, badConnStr)
		assert.Error(t, err)
	})

	t.Run("GetConnection - success", func(t *testing.T) {
		db, err := repo.GetConnection(ctx, "postgres", "postgres", "testdb")
		require.NoError(t, err)
		assert.NotNil(t, db)
		defer db.Close()

		// Test the connection works
		err = db.Ping()
		assert.NoError(t, err)
	})

	// UC-S2-03: Login Connection Probe
	t.Run("ProbeFirstAccessible - success", func(t *testing.T) {
		// Create a mock global metadata with accessible resources
		mockMetadata := &domain.GlobalMetadata{
			Databases: []domain.DatabaseMetadata{
				{Name: "testdb"},
			},
			RolePermissions: map[string]*domain.RolePermissions{
				"postgres": {
					RoleName:            "postgres",
					AccessibleDatabases: []string{"testdb"},
					AccessibleSchemas:   map[string][]string{"testdb": {"public"}},
					AccessibleTables: map[string][]domain.TableRef{
						"testdb": {{Schema: "public", Name: "test_table"}},
					},
				},
			},
		}

		firstDB, err := repo.ProbeFirstAccessible(ctx, "postgres", "postgres", mockMetadata)
		require.NoError(t, err)
		assert.NotEmpty(t, firstDB)
	})

	// UC-S2-04: Login Connection Probe Failure
	t.Run("ProbeFirstAccessible - no accessible resources", func(t *testing.T) {
		// Create metadata with no accessible resources for this role
		mockMetadata := &domain.GlobalMetadata{
			Databases: []domain.DatabaseMetadata{},
			RolePermissions: map[string]*domain.RolePermissions{
				"noaccess": {
					RoleName:            "noaccess",
					AccessibleDatabases: []string{},
					AccessibleSchemas:   map[string][]string{},
					AccessibleTables:    map[string][]domain.TableRef{},
				},
			},
		}

		_, err := repo.ProbeFirstAccessible(ctx, "noaccess", "password", mockMetadata)
		assert.Error(t, err, "Should return error when user has no accessible resources")
	})
}

// MetadataRepoConstructor is a function that creates a MetadataRepository
type MetadataRepoConstructor func() interfaces.MetadataRepository

// MetadataRepositoryRunner tests MetadataRepository implementations
func MetadataRepositoryRunner(t *testing.T, constructor MetadataRepoConstructor) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container with test data
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

	// Create test schema and tables
	_, err = db.Exec(`
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id),
			title VARCHAR(200),
			content TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		);
	`)
	require.NoError(t, err)

	repo := constructor()

	// IT-S1-02: Load Real Database Metadata with User Accessible Resources
	t.Run("LoadGlobalMetadata - with role permissions", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.NotEmpty(t, metadata.Databases)
		// UC-S1-05: Should include role permissions
		assert.NotNil(t, metadata.RolePermissions, "Should load role permissions")
	})

	// UC-S1-06: In-Memory Metadata Storage - Per Role
	t.Run("LoadRolePermissions - cache accessible resources", func(t *testing.T) {
		rolePerms, err := repo.LoadRolePermissions(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, rolePerms)
		assert.NotEmpty(t, rolePerms, "Should have at least one role")

		// Verify structure includes accessible resources per role
		for roleName, perms := range rolePerms {
			assert.NotEmpty(t, roleName)
			assert.NotNil(t, perms)
			assert.NotNil(t, perms.AccessibleDatabases, "Role should have accessible databases")
			assert.NotNil(t, perms.AccessibleSchemas, "Role should have accessible schemas map")
			assert.NotNil(t, perms.AccessibleTables, "Role should have accessible tables map")
		}
	})

	t.Run("LoadGlobalMetadata", func(t *testing.T) {
		metadata, err := repo.LoadGlobalMetadata(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, metadata)
		assert.NotEmpty(t, metadata.Databases)
	})

	t.Run("LoadDatabaseMetadata", func(t *testing.T) {
		dbMeta, err := repo.LoadDatabaseMetadata(ctx, db, "testdb")
		require.NoError(t, err)
		assert.NotNil(t, dbMeta)
		assert.Equal(t, "testdb", dbMeta.Name)
		assert.NotEmpty(t, dbMeta.Tables)
	})

	t.Run("GetTableMetadata", func(t *testing.T) {
		tableMeta, err := repo.GetTableMetadata(ctx, db, "public", "users")
		require.NoError(t, err)
		assert.NotNil(t, tableMeta)
		assert.Equal(t, "users", tableMeta.Name)
		assert.NotEmpty(t, tableMeta.Columns)
	})

	// IT-S1-03: Load Real Relations and Role Access
	t.Run("GetTableMetadata - with foreign keys", func(t *testing.T) {
		tableMeta, err := repo.GetTableMetadata(ctx, db, "public", "posts")
		require.NoError(t, err)
		assert.NotNil(t, tableMeta)
		assert.Equal(t, "posts", tableMeta.Name)
		assert.NotEmpty(t, tableMeta.ForeignKeys, "posts table should have foreign key to users")
	})

	t.Run("GetERDData", func(t *testing.T) {
		erdData, err := repo.GetERDData(ctx, db, "testdb", "public")
		require.NoError(t, err)
		assert.NotNil(t, erdData)
		assert.NotEmpty(t, erdData.Tables)
		assert.NotEmpty(t, erdData.Relationships)
	})
}

// QueryRepoConstructor is a function that creates a QueryRepository
type QueryRepoConstructor func() interfaces.QueryRepository

// QueryRepositoryRunner tests QueryRepository implementations
func QueryRepositoryRunner(t *testing.T, constructor QueryRepoConstructor) {
	t.Helper()

	ctx := context.Background()

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

	// Create test table with data
	_, err = db.Exec(`
		CREATE TABLE test_data (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			value INTEGER
		);
		INSERT INTO test_data (name, value) VALUES 
			('item1', 10),
			('item2', 20),
			('item3', 30);
	`)
	require.NoError(t, err)

	repo := constructor()

	t.Run("ExecuteQuery - SELECT", func(t *testing.T) {
		req := domain.QueryRequest{
			SQL:    "SELECT * FROM test_data",
			Params: nil,
		}
		result, err := repo.ExecuteQuery(ctx, db, req)
		require.NoError(t, err)
		assert.True(t, result.IsSelect)
		assert.NotEmpty(t, result.Rows)
		assert.NotEmpty(t, result.Columns)
	})

	t.Run("ExecuteQuery - INSERT", func(t *testing.T) {
		req := domain.QueryRequest{
			SQL:    "INSERT INTO test_data (name, value) VALUES ($1, $2)",
			Params: []interface{}{"item4", 40},
		}
		result, err := repo.ExecuteQuery(ctx, db, req)
		require.NoError(t, err)
		assert.False(t, result.IsSelect)
		assert.Equal(t, int64(1), result.RowsAffected)
	})

	t.Run("ExecuteMultipleQueries", func(t *testing.T) {
		sql := "INSERT INTO test_data (name, value) VALUES ('item5', 50); SELECT * FROM test_data WHERE name = 'item5';"
		results, err := repo.ExecuteMultipleQueries(ctx, db, sql)
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("GetTableData - basic", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema: "public",
			Table:  "test_data",
			Limit:  50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.Rows)
		assert.GreaterOrEqual(t, result.TotalCount, int64(3))
	})

	t.Run("GetTableData - with WHERE clause", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema:      "public",
			Table:       "test_data",
			WhereClause: "value > 20",
			Limit:       50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("GetTableData - with sorting", func(t *testing.T) {
		req := domain.TableDataRequest{
			Schema:   "public",
			Table:    "test_data",
			OrderBy:  "value",
			OrderDir: "DESC",
			Limit:    50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	// UC-S4-03a, UC-S5-07: Query Result Actual Size Display
	t.Run("GetTableData - shows total count", func(t *testing.T) {
		// Insert many rows to test total count
		for i := 0; i < 100; i++ {
			_, err := db.Exec("INSERT INTO test_data (name, value) VALUES ($1, $2)",
				"bulk_item", i)
			require.NoError(t, err)
		}

		req := domain.TableDataRequest{
			Schema: "public",
			Table:  "test_data",
			Limit:  50,
		}
		result, err := repo.GetTableData(ctx, db, req)
		require.NoError(t, err)
		assert.NotNil(t, result)
		// Should show actual total count even when limited
		assert.Greater(t, result.TotalCount, int64(50), "TotalCount should show all rows, not just returned")
	})

	// UC-S4-03b, UC-S5-08: Query Result Limit Hard Cap
	t.Run("GetTableData - enforces 1000 row hard limit", func(t *testing.T) {
		// Insert more than 1000 rows
		for i := 0; i < 1500; i++ {
			_, err := db.Exec("INSERT INTO test_data (name, value) VALUES ($1, $2)",
				"limit_test", i)
			require.NoError(t, err)
		}

		req := domain.TableDataRequest{
			Schema: "public",
			Table:  "test_data",
			Limit:  50,
		}

		// Fetch multiple pages up to the hard limit
		totalFetched := 0
		cursor := ""
		for totalFetched < 1500 {
			req.Cursor = cursor
			result, err := repo.GetTableData(ctx, db, req)
			require.NoError(t, err)

			if len(result.Rows) == 0 || !result.HasMore {
				break
			}

			totalFetched += len(result.Rows)
			cursor = result.NextCursor

			// Should stop at 1000 rows maximum
			if totalFetched >= 1000 {
				assert.LessOrEqual(t, totalFetched, 1000, "Should not fetch more than 1000 rows")
				assert.False(t, result.HasMore, "HasMore should be false at 1000 row limit")
				break
			}
		}

		assert.LessOrEqual(t, totalFetched, 1000, "Total fetched rows should not exceed 1000")
	})

	// UC-S5-18: GetReferencingTables for Primary Key navigation
	t.Run("GetReferencingTables - returns child tables", func(t *testing.T) {
		// Create parent and child tables
		_, err := db.Exec(`
			CREATE TABLE parent_table (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100)
			);
			CREATE TABLE child_table (
				id SERIAL PRIMARY KEY,
				parent_id INTEGER REFERENCES parent_table(id),
				data VARCHAR(100)
			);
			INSERT INTO parent_table (name) VALUES ('parent1');
			INSERT INTO child_table (parent_id, data) VALUES (1, 'child1'), (1, 'child2');
		`)
		require.NoError(t, err)

		// Get referencing tables for parent_table PK
		refTables, err := repo.GetReferencingTables(ctx, db, "public", "parent_table", "id", 1)
		require.NoError(t, err)
		assert.NotNil(t, refTables)

		// Should show child_table with count
		childCount, exists := refTables["child_table"]
		assert.True(t, exists, "Should include child_table")
		assert.Equal(t, int64(2), childCount, "Should show 2 child rows")
	})
}

// TransactionRepoConstructor is a function that creates a TransactionRepository
type TransactionRepoConstructor func() interfaces.TransactionRepository

// TransactionRepositoryRunner tests TransactionRepository implementations
func TransactionRepositoryRunner(t *testing.T, constructor TransactionRepoConstructor) {
	t.Helper()

	ctx := context.Background()

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
	_, err = db.Exec(`
		CREATE TABLE txn_test (
			id SERIAL PRIMARY KEY,
			data VARCHAR(100)
		);
	`)
	require.NoError(t, err)

	repo := constructor()

	t.Run("BeginTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)
		assert.NotNil(t, tx)
		tx.Rollback()
	})

	t.Run("CommitTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		_, err = tx.Exec("INSERT INTO txn_test (data) VALUES ('test')")
		require.NoError(t, err)

		err = repo.CommitTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify data was committed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data = 'test'").Scan(&count)
		assert.Equal(t, 1, count)
	})

	t.Run("RollbackTransaction - success", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		_, err = tx.Exec("INSERT INTO txn_test (data) VALUES ('rollback_test')")
		require.NoError(t, err)

		err = repo.RollbackTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify data was not committed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data = 'rollback_test'").Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("ExecuteOperations - buffered operations", func(t *testing.T) {
		tx, err := repo.BeginTransaction(ctx, db)
		require.NoError(t, err)

		operations := []domain.Operation{
			{
				Type:   domain.OperationInsert,
				Schema: "public",
				Table:  "txn_test",
				Data: map[string]interface{}{
					"data": "op1",
				},
			},
			{
				Type:   domain.OperationInsert,
				Schema: "public",
				Table:  "txn_test",
				Data: map[string]interface{}{
					"data": "op2",
				},
			},
		}

		err = repo.ExecuteOperations(ctx, tx, operations)
		require.NoError(t, err)

		err = repo.CommitTransaction(ctx, tx)
		assert.NoError(t, err)

		// Verify operations were executed
		var count int
		db.QueryRow("SELECT COUNT(*) FROM txn_test WHERE data IN ('op1', 'op2')").Scan(&count)
		assert.Equal(t, 2, count)
	})
}
